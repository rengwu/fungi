package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRuntimeCheckAndSavedMissingRuntime(t *testing.T) {
	root := t.TempDir()
	helper := filepath.Join(root, rendererRelativePath())
	if err := os.MkdirAll(filepath.Dir(helper), 0700); err != nil {
		t.Fatal(err)
	}
	if err := copyExecutable(testFixture, helper); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FUNGI_RENDERER_DIR", filepath.Join(root, "renderer"))
	t.Setenv("FUNGI_STATE_DIR", filepath.Join(root, "sessions"))
	var out bytes.Buffer
	if code := run([]string{"doctor"}, &out); code != 0 {
		t.Fatalf("runtime probe: %d %s", code, &out)
	}
	t.Setenv("FUNGI_TEST_RUNTIME_MISSING", "1")
	out.Reset()
	if code := run([]string{"doctor"}, &out); code != 3 {
		t.Fatalf("missing runtime probe: %d %s", code, &out)
	}
	source := filepath.Join(root, "index.html")
	if err := os.WriteFile(source, []byte("<!doctype html>"), 0600); err != nil {
		t.Fatal(err)
	}
	data, err := exec.Command(testBinary, "open", source, "--delivery", "wait").Output()
	if err == nil {
		t.Fatalf("missing runtime launch succeeded: %s", data)
	}
	var r result
	if err := json.Unmarshal(data, &r); err != nil || r.Outcome != "failed" || r.Error == nil || r.Error.Code != "runtime_missing" {
		t.Fatalf("failure not saved: %s", data)
	}
	again, err := exec.Command(testBinary, "collect", r.Session).Output()
	if err == nil || string(again) != string(data) {
		t.Fatalf("collect changed failure: %s (%v)", again, err)
	}

}
