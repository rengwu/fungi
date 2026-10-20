package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

var testBinary string
var testHelper string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "fungi-tests-")
	if err != nil {
		panic(err)
	}
	testBinary = filepath.Join(dir, "fungi")
	testHelper = filepath.Join(dir, "helper")
	for _, build := range [][]string{{"build", "-o", testBinary, "."}, {"build", "-o", testHelper, "./examples/helper"}} {
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
	var opened opened
	b := cli(t, root, "open", "--exec", testHelper, "--delivery", "manual", "--", "wait")
	if err := json.Unmarshal(b, &opened); err != nil {
		t.Fatal(err)
	}
	pending := cli(t, root, "collect", opened.Session)
	if !strings.Contains(string(pending), `"pending"`) {
		t.Fatalf("expected pending: %s", pending)
	}
	s := &store{dir: filepath.Join(root, opened.Session)}
	// Bounded test-driver wait; the product uses no status-polling loop.
	time.Sleep(11 * time.Second)
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
	cmd := exec.Command(testBinary, "open", "--exec", testHelper, "--delivery", "wait", "--", "wait")
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	// This finite test delay is not a model-facing polling mechanism.
	time.Sleep(time.Second)
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatal("owner did not start", err)
	}
	id := entries[0].Name()
	if err = cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Wait()
	time.Sleep(11 * time.Second)
	s := &store{dir: filepath.Join(root, id)}
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
	cmd := exec.Command(testBinary, "open", "--exec", "/usr/bin/false", "--delivery", "wait")
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+root)
	b, err := cmd.Output()
	if err == nil {
		t.Fatal("failed helper exited successfully")
	}
	var r result
	if err = json.Unmarshal(b, &r); err != nil {
		t.Fatal("expected one JSON envelope", string(b))
	}
	if r.Outcome != "failed" || r.Error == nil {
		t.Fatal("missing failure result", string(b))
	}
	// A helper can terminate immediately after emitting its final message.
	// Even a nonzero exit cannot overwrite an already committed submission.
	script := `printf '%s\n' '{"type":"ready"}' '{"type":"result","outcome":"submitted","data":{"drained":true}}'
exit 7`
	b = cli(t, root, "open", "--exec", "/bin/sh", "--delivery", "wait", "--", "-c", script)
	if err = json.Unmarshal(b, &r); err != nil || r.Outcome != "submitted" {
		t.Fatal("result lost at process exit", string(b), err)
	}
}
