package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

var testBinary string
var testHelper string
var testFixture string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "fungi-tests-")
	if err != nil {
		panic(err)
	}
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	testBinary = filepath.Join(dir, "fungi"+suffix)
	testHelper = filepath.Join(dir, "helper"+suffix)
	testFixture = filepath.Join(dir, "protocol-fixture"+suffix)
	for _, build := range [][]string{{"build", "-o", testBinary, "."}, {"build", "-o", testHelper, "./examples/helper"}, {"build", "-o", testFixture, "./scripts/protocol-fixture"}} {
		if out, err := exec.Command("go", build...).CombinedOutput(); err != nil {
			panic(string(out))
		}
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
func cli(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command(testBinary, args...)
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %v\n%s", args, err, out)
	}
	return out
}
func TestNativeRoundTripAndConcurrentIsolation(t *testing.T) {
	root := t.TempDir()
	var wg sync.WaitGroup
	ids := make(chan string, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := cli(t, root, "open", "--exec", testHelper, "--delivery", "wait")
			var r result
			if err := json.Unmarshal(b, &r); err != nil {
				t.Error(err)
				return
			}
			if r.Outcome != "submitted" || len(r.Artifacts) != 1 {
				t.Errorf("bad result %s", b)
				return
			}
			data, err := os.ReadFile(r.Artifacts[0].Path)
			if err != nil || !strings.Contains(string(data), "native helper") {
				t.Error("missing artifact", err)
			}
			again := cli(t, root, "collect", r.Session)
			if string(again) != string(b) {
				t.Error("collect differed from wait")
			}
			ids <- r.Session
		}()
	}
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatal("sessions collided")
		}
		seen[id] = true
	}
}
func TestManualOwnerSurvivesLauncherAndCollects(t *testing.T) {
	root := t.TempDir()
	gate := filepath.Join(t.TempDir(), "submit")
	defer os.WriteFile(gate, []byte("cancelled"), 0600)
	var opened opened
	b := cli(t, root, "open", "--exec", testFixture, "--delivery", "manual", "--", "update-session", gate)
	if err := json.Unmarshal(b, &opened); err != nil {
		t.Fatal(err)
	}
	pending := cli(t, root, "collect", opened.Session)
	if !strings.Contains(string(pending), `"pending"`) {
		t.Fatalf("expected pending: %s", pending)
	}
	s := &store{dir: filepath.Join(root, opened.Session)}
	// Release the synthetic helper only after the launcher has returned and
	// collection has observed pending. No scheduler-dependent fixed sleep.
	if err := os.WriteFile(gate, []byte("submitted"), 0600); err != nil {
		t.Fatal(err)
	}
	eventually(t, func() bool { _, err := s.readResult(); return err == nil })
	r, err := s.readResult()
	if err != nil {
		t.Fatal(err)
	}
	if r.Outcome != "submitted" {
		t.Fatal(r)
	}
	// A later tool invocation may not inherit launch-time storage overrides.
	// The receipt must locate its own result without relying on that context.
	if len(opened.Collect) != 3 || opened.CollectEnv["FUNGI_STATE_DIR"] != root {
		t.Fatalf("incomplete collection instructions: %+v", opened)
	}
	cmd := exec.Command(opened.Collect[0], opened.Collect[1:]...)
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+opened.CollectEnv["FUNGI_STATE_DIR"])
	collected, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("receipt collection failed: %v: %s", err, collected)
	}
	var fromReceipt result
	if err := json.Unmarshal(collected, &fromReceipt); err != nil || fromReceipt.Session != opened.Session || fromReceipt.Outcome != "submitted" {
		t.Fatalf("receipt collected wrong result: %s (%v)", collected, err)
	}
}
func TestInterruptedSessionReopenPreservesOldRecord(t *testing.T) {
	root := t.TempDir()
	t.Setenv("FUNGI_STATE_DIR", root)
	id := newID()
	s, _ := sessionStore(id)
	l := launch{Protocol: 1, Session: id, Exec: testHelper, Input: json.RawMessage(`{"seed":3}`), Delivery: delivery{Mode: "manual"}}
	if err := s.init(l); err != nil {
		t.Fatal(err)
	}
	if err := s.saveDraft(message{Data: json.RawMessage(`{"notes":"restore me"}`)}); err != nil {
		t.Fatal(err)
	}
	if _, err := collect(id); err == nil || asProblem(err).Code != "session_interrupted" {
		t.Fatal(err)
	}
	b := cli(t, root, "reopen", id, "--delivery", "wait")
	var r result
	if err := json.Unmarshal(b, &r); err != nil {
		t.Fatal(err)
	}
	if r.Session == id || !strings.Contains(string(r.Data), "restore me") {
		t.Fatal("reopen did not restore", string(b))
	}
	var next launch
	if err := readJSON(filepath.Join(root, r.Session, "launch.json"), &next); err != nil {
		t.Fatal(err)
	}
	if next.ReopenedFrom != id {
		t.Fatal("missing link")
	}
	if _, err := s.readResult(); !os.IsNotExist(err) {
		t.Fatal("old interrupted session changed")
	}
}
func TestWaitConnectionHasNoHumanTimeout(t *testing.T) {
	// The ten-second helper exceeds the owner's short status request timeout.
	// This is a bounded regression check, not evidence of indefinite host waits.
	root := t.TempDir()
	start := time.Now()
	cli(t, root, "open", "--exec", testHelper, "--delivery", "wait", "--", "wait")
	if time.Since(start) < 10*time.Second {
		t.Fatal("wait returned before helper finished")
	}
}

func TestDisconnectedWaiterDoesNotCancel(t *testing.T) {
	root := t.TempDir()
	gate := filepath.Join(t.TempDir(), "submit")
	cmd := exec.Command(testBinary, "open", "--exec", testFixture, "--delivery", "wait", "--", "update-session", gate)
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.WriteFile(gate, []byte("submitted"), 0600)
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}()
	// Wait for a real helper/owner round trip before disconnecting the caller.
	// A fixed sleep can kill the launcher before its owner has even started.
	var id []byte
	eventually(t, func() bool {
		var err error
		id, err = os.ReadFile(gate + ".ready")
		return err == nil && len(id) == 32
	})
	s := &store{dir: filepath.Join(root, string(id))}
	defer func() {
		if t.Failed() {
			log, _ := os.ReadFile(s.path("worker.log"))
			t.Logf("worker diagnostics: %s", log)
		}
	}()
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Wait()
	if _, err := s.readResult(); !os.IsNotExist(err) {
		t.Fatal("fixture completed before submission was released", err)
	}
	if err := os.WriteFile(gate, []byte("submitted"), 0600); err != nil {
		t.Fatal(err)
	}
	eventually(t, func() bool { _, err := s.readResult(); return err == nil })
	r, err := s.readResult()
	if err != nil {
		t.Fatal("waiter loss killed owner", err)
	}
	if r.Outcome != "submitted" {
		t.Fatal(r)
	}
}

func TestHelperExitAndStdoutDrain(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command(testBinary, "open", "--exec", testFixture, "--delivery", "wait", "--", "exit")
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
	b, err := cmd.Output()
	if err == nil {
		t.Fatal("failed helper exited successfully")
	}
	var r result
	if err = json.Unmarshal(b, &r); err != nil {
		t.Fatal("expected one JSON envelope", string(b))
	}
	if r.Outcome != "failed" || r.Error == nil || r.Error.Code != "helper_exited" {
		t.Fatal("missing failure result", string(b))
	}
	// A helper can terminate immediately after emitting its final message.
	// Even a nonzero exit cannot overwrite an already committed submission.
	b = cli(t, root, "open", "--exec", testFixture, "--delivery", "wait", "--", "early-result")
	if err = json.Unmarshal(b, &r); err != nil || r.Outcome != "submitted" {
		t.Fatal("result lost at process exit", string(b), err)
	}
}

func TestRendererStartupFailureIsSavedAndDiagnosable(t *testing.T) {
	root := t.TempDir()
	renderer := filepath.Join(root, "renderer")
	helper := filepath.Join(root, rendererRelativePath())
	if err := os.MkdirAll(filepath.Dir(helper), 0700); err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(testFixture)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(helper, bytes, 0700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "index.html")
	if err = os.WriteFile(source, []byte("<!doctype html>"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FUNGI_RENDERER_DIR", renderer)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	state := filepath.Join(root, "sessions")
	cmd := exec.CommandContext(ctx, testBinary, "open", source, "--delivery", "wait")
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+state)
	b, err := cmd.Output()
	var r result
	if ctx.Err() != nil || err == nil || json.Unmarshal(b, &r) != nil || r.Outcome != "failed" || r.Error == nil || r.Error.Code != "renderer_startup_failed" {
		t.Fatalf("expected one saved renderer startup failure: %s (%v)", b, err)
	}
	log := filepath.Join(state, r.Session, "worker.log")
	diagnostics, err := os.ReadFile(log)
	if err != nil || !strings.Contains(string(diagnostics), "protocol fixture: unknown mode") || !strings.Contains(r.Error.Message, log) {
		t.Fatalf("startup diagnostics unavailable: %+v, %s (%v)", r.Error, diagnostics, err)
	}
	cmd = exec.Command(testBinary, "collect", r.Session)
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+state)
	again, err := cmd.Output()
	if err == nil || string(again) != string(b) {
		t.Fatalf("collection must return the identical failed result: %s (%v)", again, err)
	}
}

func TestBrokenHelperCannotHoldOwnerOpen(t *testing.T) {
	for _, mode := range []string{"malformed", "closed-stdout"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, testBinary, "open", "--exec", testFixture, "--delivery", "wait", "--", mode)
			cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
			b, err := cmd.Output()
			if ctx.Err() != nil {
				t.Fatal("owner hung after broken helper output")
			}
			var r result
			if err == nil || json.Unmarshal(b, &r) != nil || r.Outcome != "failed" || r.Error == nil {
				t.Fatalf("expected one failed result, got %s (%v)", b, err)
			}
		})
	}
}

func TestHelperThatDoesNotReadLaunchIsBounded(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.json")
	if err := os.WriteFile(input, []byte(`"`+strings.Repeat("x", 200000)+`"`), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, testBinary, "open", "--exec", testFixture, "--input", input, "--delivery", "wait", "--", "never-read")
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+filepath.Join(root, "sessions"))
	b, err := cmd.Output()
	var r result
	if err == nil || json.Unmarshal(b, &r) != nil || r.Outcome != "failed" || r.Error == nil {
		t.Fatalf("expected a failed session, got %s (%v)", b, err)
	}
	if ctx.Err() != nil {
		t.Fatal("helper blocked launch until its own exit")
	}
}
