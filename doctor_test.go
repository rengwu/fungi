package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHeadlessDoctor(t *testing.T) {
	root := filepath.Join(t.TempDir(), "state with spaces")
	cmd := exec.Command(testBinary, "doctor")
	cmd.Env = append(os.Environ(), "PATH=", "DISPLAY=", "FUNGI_STATE_DIR="+root)
	output, err := cmd.Output()
	var report struct {
		Protocol  int
		Engine    string
		Available bool
		StateDir  string
	}
	if err != nil || json.Unmarshal(output, &report) != nil || report.Protocol != 2 || report.Engine != "headless" || !report.Available || report.StateDir != root {
		t.Fatalf("headless doctor: %s (%v)", output, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("doctor left state behind: %v %v", entries, err)
	}
	blocked := filepath.Join(t.TempDir(), "file")
	if err = os.WriteFile(blocked, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FUNGI_STATE_DIR", blocked)
	var out bytes.Buffer
	if code := run([]string{"doctor"}, &out); code != 1 || !json.Valid(out.Bytes()) {
		t.Fatalf("unusable storage: %d %s", code, &out)
	}
}
