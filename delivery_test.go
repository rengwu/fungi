package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCallbackUncertaintyDoesNotResendOrLoseResult(t *testing.T) {
	s := testStore(t)
	id := filepath.Base(s.dir)
	if _, err := s.commit(id, message{Outcome: "submitted", Data: json.RawMessage(`{"answer":"saved"}`)}); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "notice.txt")
	t.Setenv("FUNGI_TEST_NOTICE", log)
	l := launch{Session: id, Delivery: delivery{Mode: "callback", Adapter: "codex", Thread: "11111111-1111-1111-1111-111111111111", Command: testFixture}}
	notifyResult(s, l)
	first, err := os.ReadFile(log)
	if err != nil {
		t.Fatal(err)
	}
	notifyResult(s, l)
	second, _ := os.ReadFile(log)
	if string(first) != string(second) {
		t.Fatal("ambiguous callback retried")
	}
	if !strings.Contains(string(first), `"FUNGI_STATE_DIR"`) || !strings.Contains(string(first), l.Delivery.Thread) {
		t.Fatal("callback lost exact caller/storage")
	}
	var d delivery
	if err := readJSON(s.path("delivery.json"), &d); err != nil || d.State != "unknown" {
		t.Fatal(d, err)
	}
	if r, err := s.readResult(); err != nil || r.Outcome != "submitted" {
		t.Fatal("notification failure lost saved answer", err)
	}
}
