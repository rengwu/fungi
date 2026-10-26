package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOfflineHelpAndDiscovery(t *testing.T) {
	// TestMain builds a standalone executable outside the checkout. Run it
	// elsewhere, with no renderer or manuals beside it and no tools on PATH.
	directory := t.TempDir()
	invoke := func(args ...string) ([]byte, error) {
		cmd := exec.Command(testBinary, args...)
		cmd.Dir = directory
		cmd.Env = append(os.Environ(), "PATH=", "FUNGI_STATE_DIR="+directory)
		return cmd.CombinedOutput()
	}
	for _, topic := range manualTopics() {
		got, err := invoke("help", topic)
		body, _ := manual.ReadFile("docs/manual/" + topic + ".md")
		title := strings.TrimPrefix(strings.SplitN(string(body), "\n", 2)[0], "# ")
		if err != nil || !strings.HasPrefix(string(got), title+"\n") {
			t.Fatalf("offline help %s: %v\n%s", topic, err, got)
		}
		if strings.Contains(string(got), "```") || strings.Contains(string(got), "\n| ") || strings.Contains(string(got), "\n## ") {
			t.Fatalf("raw Markdown in help %s: %s", topic, got)
		}
		// Commands, JSON and source code remain copyable, including long lines
		// and shell operators. Do not compare to the formatter's own output.
		for n, block := range strings.Split(string(body), "```") {
			if n%2 == 0 {
				continue
			}
			_, code, _ := strings.Cut(block, "\n")
			for _, line := range strings.Split(strings.TrimSuffix(code, "\n"), "\n") {
				if !strings.Contains(string(got), "  "+line+"\n") {
					t.Fatalf("help %s lost example line %q", topic, line)
				}
			}
		}
	}
	for _, args := range [][]string{nil, {"--help"}, {"help"}} {
		got, err := invoke(args...)
		if err != nil || !strings.Contains(string(got), "Common starting points") {
			t.Fatalf("entry %v: %v\n%s", args, err, got)
		}
	}
	for _, args := range [][]string{{"help", "missing"}, {"help", "html", "extra"}, {"help", "../cli"}} {
		got, err := invoke(args...)
		if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 2 || !json.Valid(got) {
			t.Fatalf("invalid help %v: %v\n%s", args, err, got)
		}
	}
	got, err := invoke("describe")
	var index struct{ Topics []string }
	if err != nil || json.Unmarshal(got, &index) != nil {
		t.Fatalf("describe: %v %s", err, got)
	}
	for _, topic := range index.Topics {
		got, err = invoke("describe", topic)
		var description struct{ Reference []string }
		if err != nil || json.Unmarshal(got, &description) != nil || len(description.Reference) != 3 {
			t.Fatalf("describe %s: %v %s", topic, err, got)
		}
		if _, err = invoke(description.Reference[1:]...); err != nil {
			t.Fatal(err)
		}
	}
	got, err = invoke("describe", "examples")
	var discovery struct{ Capability struct{ Items []example } }
	if err != nil || json.Unmarshal(got, &discovery) != nil || len(discovery.Capability.Items) == 0 {
		t.Fatalf("examples: %v %s", err, got)
	}
	for _, item := range discovery.Capability.Items {
		if item.Available || len(item.Argv) != 0 || !filepath.IsAbs(item.Source) {
			t.Fatalf("core-only discovery promises an unavailable example: %+v", item)
		}
	}
}
