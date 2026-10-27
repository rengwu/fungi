package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWindowDimensions(t *testing.T) {
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
	for name, body := range map[string]string{"index.html": "<!doctype html>", "form.json": `{"title":"Size check","fields":[]}`} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		name, source  string
		flags         []string
		width, height int
	}{
		{"defaults", "index.html", nil, 0, 0},
		{"both", "index.html", []string{"--width", "1100", "--height", "800"}, 1100, 800},
		{"width only", "index.html", []string{"--width", "700"}, 700, 0},
		{"height only", "index.html", []string{"--height", "600"}, 0, 600},
		{"form minimum", "form.json", []string{"--width", "420", "--height", "420"}, 420, 420},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"open", filepath.Join(root, tc.source), "--delivery", "manual"}, tc.flags...)
			l, _, err := parseLaunch(args)
			if err != nil || l.Width != tc.width || l.Height != tc.height {
				t.Fatalf("launch: %+v (%v)", l, err)
			}
			s, _ := sessionStore(l.Session)
			if err := s.init(l); err != nil {
				t.Fatal(err)
			}
			// An interrupted session is eligible for explicit reopen. Persist and
			// read the real launch record, including older records without sizes.
			for _, override := range []bool{false, true} {
				args := []string{"reopen", l.Session, "--delivery", "manual"}
				wantWidth := tc.width
				if override {
					args = append(args, "--width", "900")
					wantWidth = 900
				}
				reopened, _, err := parseLaunch(args)
				if err != nil || reopened.Width != wantWidth || reopened.Height != tc.height || reopened.ReopenedFrom != l.Session {
					t.Fatalf("reopen: %+v (%v)", reopened, err)
				}
			}
		})
	}
}

func TestInvalidWindowDimensions(t *testing.T) {
	root := filepath.Join(t.TempDir(), "sessions")
	t.Setenv("FUNGI_STATE_DIR", root)
	for _, flag := range []string{"--width", "--height"} {
		for _, value := range []string{"0", "-1", "419", "1.5", "wide", "", "2147483648", "999999999999999999999999"} {
			var out bytes.Buffer
			code := run([]string{"open", "index.html", flag, value}, &out)
			if code != 2 || !json.Valid(out.Bytes()) || !bytes.Contains(out.Bytes(), []byte(`"invalid_arguments"`)) {
				t.Fatalf("%s %q: exit %d %s", flag, value, code, &out)
			}
		}
		for _, args := range [][]string{
			{"open", "index.html", flag},
			{"open", "--exec", testHelper, flag, "600", "--delivery", "manual"},
		} {
			var out bytes.Buffer
			if code := run(args, &out); code != 2 {
				t.Fatalf("%v: exit %d %s", args, code, &out)
			}
		}
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("invalid arguments created session storage: %v", err)
	}
}
