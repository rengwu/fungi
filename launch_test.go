package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestHelperLaunchAndReopen(t *testing.T) {
	root := t.TempDir()
	t.Setenv("FUNGI_STATE_DIR", filepath.Join(root, "sessions"))
	input := filepath.Join(root, "input with spaces.json")
	data := []byte(`{"settings":{"width":900},"items":[true,null,12]}`)
	if err := os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	helperArgs := []string{"page with spaces.html", "--width", "900", "$(literal)"}
	args := append([]string{"open", "--exec", testHelper, "--input", input, "--delivery", "manual", "--"}, helperArgs...)
	l, mode, err := parseLaunch(args)
	if err != nil || mode != "manual" || l.Protocol != 2 || l.Exec != testHelper || !reflect.DeepEqual(l.Args, helperArgs) || !bytes.Equal(l.Input, data) {
		t.Fatalf("helper launch: %+v (%v)", l, err)
	}
	s, _ := sessionStore(l.Session)
	if err = s.init(l); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(input, []byte(`null`), 0600); err != nil {
		t.Fatal(err)
	}
	reopened, _, err := parseLaunch([]string{"reopen", l.Session, "--delivery", "manual"})
	if err != nil || reopened.ReopenedFrom != l.Session || reopened.Session == l.Session || !bytes.Equal(reopened.Input, data) || !reflect.DeepEqual(reopened.Args, helperArgs) {
		t.Fatalf("reopen lost helper arguments or snapshotted input: %+v (%v)", reopened, err)
	}
}

func TestEngineRejectsPresentationCommands(t *testing.T) {
	root := filepath.Join(t.TempDir(), "sessions")
	t.Setenv("FUNGI_STATE_DIR", root)
	for _, args := range [][]string{
		{"open"}, {"open", "index.html"}, {"open", "form.json"},
		{"open", "--exec", testHelper, "index.html"},
		{"open", "--exec", testHelper, "--width", "600"},
		{"open", "--exec", testHelper, "--height", "600"},
		{"open", "--exec", ""}, {"open", "--exec"},
		{"open", "--exec", testHelper, "--exec", testHelper},
		{"reopen", newID(), "--width", "600"},
	} {
		var out bytes.Buffer
		if code := run(args, &out); code != 2 || !json.Valid(out.Bytes()) {
			t.Fatalf("%v: exit %d %s", args, code, &out)
		}
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("invalid launch created state: %v", err)
	}
}

func TestOtherProtocolSessionsArePreserved(t *testing.T) {
	t.Setenv("FUNGI_STATE_DIR", t.TempDir())
	for _, version := range []int{1, 3} {
		id := newID()
		s, _ := sessionStore(id)
		if err := s.init(launch{Protocol: version, Session: id, Exec: testHelper, Input: json.RawMessage(`null`)}); err != nil {
			t.Fatal(err)
		}
		saved := []byte(`{"protocol":1,"kind":"result","outcome":"closed","artifacts":[]}`)
		if err := os.WriteFile(s.path("result.json"), saved, 0600); err != nil {
			t.Fatal(err)
		}
		for _, action := range []string{"collect", "reopen"} {
			var out bytes.Buffer
			if code := run([]string{action, id}, &out); code != 1 || !bytes.Contains(out.Bytes(), []byte("unsupported_protocol")) {
				t.Fatalf("%s silently accepted protocol %d: %s", action, version, &out)
			}
		}
		got, _ := os.ReadFile(s.path("result.json"))
		if !bytes.Equal(got, saved) {
			t.Fatal("incompatible result was modified")
		}
	}
}
