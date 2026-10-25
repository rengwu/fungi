package main

import (
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalAssetsAndRanges(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "skill")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{"index.html": "<!doctype html><title>fixture</title>", "sample.wav": "0123456789", "module.mjs": "export default 42;", "empty.txt": ""} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
	}
	root, _ = filepath.EvalSymlinks(root)
	handler, err := newAssets(filepath.Join(root, "index.html"), false)
	if err != nil {
		t.Fatal(err)
	}
	defer handler.root.Close()
	get := func(name, method, rangeHeader string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://wails/"+name, nil)
		if rangeHeader != "" {
			r.Header.Set("Range", rangeHeader)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	for _, c := range []struct {
		rangeHeader, body string
		status            int
	}{
		{"", "0123456789", 200}, {"bytes=2-5", "2345", 206}, {"bytes=7-", "789", 206}, {"bytes=-3", "789", 206}, {"bytes=99-", "", 416},
	} {
		w := get("sample.wav", "GET", c.rangeHeader)
		if w.Code != c.status || (c.status != 416 && w.Body.String() != c.body) {
			t.Fatalf("range %q: %d %q", c.rangeHeader, w.Code, w.Body.String())
		}
	}
	if w := get("sample.wav", "HEAD", ""); w.Code != 200 || w.Body.Len() != 0 || w.Header().Get("Content-Length") != "10" {
		t.Fatal("HEAD response", w)
	}
	if w := get("sample.wav", "POST", ""); w.Code != 405 {
		t.Fatal("POST accepted", w)
	}
	if w := get("module.mjs", "GET", ""); !strings.HasPrefix(w.Header().Get("Content-Type"), "text/javascript") {
		t.Fatal("module MIME", w)
	}
	if w := get("empty.txt", "GET", ""); w.Code != 200 || w.Body.Len() != 0 {
		t.Fatal("empty file", w)
	}
	if w := get("", "GET", ""); !strings.Contains(w.Body.String(), "/_fungi/bridge.js") {
		t.Fatal("bridge missing")
	}
	t.Run("outside asset root", func(t *testing.T) {
		private := filepath.Join(dir, "private.txt")
		if err := os.WriteFile(private, []byte("private"), 0600); err != nil {
			t.Fatal(err)
		}
		if w := get("%2e%2e%2fprivate.txt", "GET", ""); w.Code != 403 {
			t.Fatal("escaped path accepted", w.Code)
		}
		if err := os.Symlink(private, filepath.Join(root, "escape.txt")); err != nil {
			t.Skip("symlink unavailable:", err)
		}
		w := get("escape.txt", "GET", "")
		b, _ := io.ReadAll(w.Result().Body)
		if w.Code != 403 || strings.Contains(string(b), "private") {
			t.Fatal("symlink escape", w.Code)
		}
	})
}
