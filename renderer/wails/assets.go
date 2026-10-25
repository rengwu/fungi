package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	builtin "fungi/renderer"
	"golang.org/x/net/html"
)

const contentPolicy = "default-src 'self' data: blob:; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; connect-src 'self' data: blob:; worker-src 'self' blob:; frame-src 'self'; object-src 'none'; base-uri 'none'; form-action 'none'"

type assets struct {
	root  *os.Root
	entry string
	form  bool
}

func newAssets(source string, form bool) (*assets, error) {
	root, err := os.OpenRoot(filepath.Dir(source))
	if err != nil {
		return nil, err
	}
	return &assets{root: root, entry: filepath.Base(source), form: form}, nil
}
func (a *assets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Security-Policy", contentPolicy)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != "GET" && r.Method != "HEAD" {
		w.WriteHeader(405)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "_fungi/bridge.js" {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(bridgeJS))
		return
	}
	entry := name == "" || name == "index.html"
	if entry {
		name = a.entry
	}
	var data []byte
	if a.form && entry {
		data, _ = builtin.Builtin.ReadFile("form.html")
	} else if a.form && name == "form.js" {
		data, _ = builtin.Builtin.ReadFile("form.js")
	} else {
		resolved, err := filepath.EvalSymlinks(filepath.Join(a.root.Name(), filepath.FromSlash(name)))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		rel, err := filepath.Rel(a.root.Name(), resolved)
		if err != nil || !filepath.IsLocal(rel) {
			w.WriteHeader(403)
			return
		}
		f, err := a.root.Open(filepath.FromSlash(name))
		if err != nil {
			w.WriteHeader(403)
			return
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil || !info.Mode().IsRegular() {
			w.WriteHeader(404)
			return
		}
		if !entry {
			// Windows may have no .mjs association. Module loading must not
			// depend on a user's registry or MIME database.
			t := mime.TypeByExtension(filepath.Ext(name))
			switch strings.ToLower(filepath.Ext(name)) {
			case ".js", ".mjs":
				t = "text/javascript; charset=utf-8"
			case ".wasm":
				t = "application/wasm"
			}
			if t != "" {
				w.Header().Set("Content-Type", t)
			}
			http.ServeContent(w, r, name, info.ModTime(), f)
			return
		}
		data, err = io.ReadAll(f)
		if err != nil {
			http.Error(w, "Unable to read entry", 500)
			return
		}
	}
	if entry {
		var err error
		data, err = injectBridge(data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/javascript")
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
}
func injectBridge(data []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var visit func(*html.Node) bool
	visit = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "head" {
			n.InsertBefore(&html.Node{Type: html.ElementNode, Data: "script", Attr: []html.Attribute{{Key: "src", Val: "/_fungi/bridge.js"}}}, n.FirstChild)
			return true
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if visit(child) {
				return true
			}
		}
		return false
	}
	if !visit(doc) {
		return nil, fmt.Errorf("entry has no head")
	}
	var out bytes.Buffer
	err = html.Render(&out, doc)
	return out.Bytes(), err
}
