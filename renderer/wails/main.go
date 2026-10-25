// Experimental Wails renderer, invoked through Fungi's existing --exec protocol.
// This prototype is for trusted local pages; see README.md for isolation gaps.
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed bridge.js
var bridgeJS []byte

type attachment struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	MediaType string `json:"mediaType"`
	Bytes     []byte `json:"bytes,omitempty"`
	Path      string `json:"path,omitempty"`
}
type launch struct {
	Protocol       int             `json:"protocol"`
	Source         string          `json:"source"`
	OutputDir      string          `json:"outputDir"`
	Input          json.RawMessage `json:"input"`
	Draft          json.RawMessage `json:"draft"`
	DraftArtifacts []struct {
		Key       string `json:"key"`
		Name      string `json:"name"`
		MediaType string `json:"mediaType"`
		Path      string `json:"path"`
	} `json:"draftArtifacts"`
	Delivery json.RawMessage `json:"delivery"`
	Limits   struct {
		JSON     int `json:"jsonBytes"`
		Artifact int `json:"artifactBytes"`
		Total    int `json:"totalArtifactBytes"`
	} `json:"limits"`
}

// Bridge exposes only the operations needed by bridge.js. Wails' own runtime
// is also exposed by upstream; this is not the Electron isolation boundary.
type Bridge struct {
	mu       sync.Mutex
	launch   launch
	form     json.RawMessage
	dec      *json.Decoder
	ready    bool
	terminal bool
	sequence int
}

func (b *Bridge) signalReady() error {
	if b.ready {
		return nil
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"type": "ready"}); err != nil {
		return err
	}
	b.ready = true
	return nil
}

func (b *Bridge) Ready() (map[string]any, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	items := make([]attachment, 0, len(b.launch.DraftArtifacts))
	for _, item := range b.launch.DraftArtifacts {
		data, err := os.ReadFile(item.Path)
		if err != nil {
			return nil, err
		}
		items = append(items, attachment{Key: item.Key, Name: item.Name, MediaType: item.MediaType, Bytes: data})
	}
	if err := b.signalReady(); err != nil {
		return nil, err
	}
	return map[string]any{"input": b.launch.Input, "draft": b.launch.Draft, "attachments": items,
		"form": b.form, "limits": b.launch.Limits, "delivery": b.launch.Delivery}, nil
}

// request runs under mu, preserving write/ack order with native window close.
func (b *Bridge) request(value map[string]any) (json.RawMessage, error) {
	b.sequence++
	value["id"] = fmt.Sprint(b.sequence)
	line, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(line)+1 > b.launch.Limits.JSON {
		return nil, errors.New("response envelope exceeds JSON limit")
	}
	if _, err := os.Stdout.Write(append(line, '\n')); err != nil {
		return nil, err
	}
	var ack struct {
		ID     string          `json:"id"`
		OK     bool            `json:"ok"`
		Result json.RawMessage `json:"result"`
		Error  struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := b.dec.Decode(&ack); err != nil {
		return nil, err
	}
	if ack.ID != value["id"] {
		return nil, errors.New("unexpected owner acknowledgement")
	}
	if !ack.OK {
		return nil, errors.New(ack.Error.Message)
	}
	return ack.Result, nil
}

func (b *Bridge) Persist(kind string, data json.RawMessage, items []attachment) (json.RawMessage, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.terminal {
		return nil, errors.New("session already completed")
	}
	if kind != "draft" && kind != "submitted" && kind != "cancelled" {
		return nil, errors.New("invalid outcome")
	}
	if len(data) > b.launch.Limits.JSON {
		return nil, errors.New("response JSON exceeds limit")
	}
	files := make([]attachment, 0, len(items))
	defer func() {
		for _, item := range files {
			_ = os.Remove(filepath.Join(b.launch.OutputDir, item.Path))
		}
	}()
	keys := make(map[string]bool)
	total := 0
	for _, item := range items {
		if item.Key == "" || item.Name == "" || keys[item.Key] {
			return nil, errors.New("attachments need unique keys and names")
		}
		keys[item.Key] = true
		total += len(item.Bytes)
		if len(item.Bytes) > b.launch.Limits.Artifact || total > b.launch.Limits.Total {
			return nil, errors.New("attachment size limit exceeded")
		}
		f, err := os.CreateTemp(b.launch.OutputDir, "wails-")
		if err != nil {
			return nil, err
		}
		item.Path, item.Name = filepath.Base(f.Name()), filepath.Base(item.Name)
		files = append(files, attachment{Key: item.Key, Name: item.Name, MediaType: item.MediaType, Path: item.Path})
		_, writeErr := f.Write(item.Bytes)
		closeErr := f.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			return nil, err
		}
	}
	value := map[string]any{"type": "draft", "data": data, "artifacts": files}
	if kind != "draft" {
		value["type"], value["outcome"] = "result", kind
	}
	result, err := b.request(value)
	if err == nil && kind != "draft" {
		b.terminal = true
	}
	return result, err
}

func (b *Bridge) beforeClose(context.Context) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.terminal {
		return false
	}
	if _, err := b.request(map[string]any{"type": "result", "outcome": "closed"}); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot save close outcome:", err)
		return true
	}
	b.terminal = true
	return false
}

type stderrLogger struct{}

func (stderrLogger) Print(s string)   { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Trace(s string)   { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Debug(s string)   { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Info(s string)    { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Warning(s string) { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Error(s string)   { fmt.Fprintln(os.Stderr, s) }
func (stderrLogger) Fatal(s string)   { fmt.Fprintln(os.Stderr, s); os.Exit(1) }

func run() error {
	b := &Bridge{dec: json.NewDecoder(os.Stdin)}
	if err := b.dec.Decode(&b.launch); err != nil {
		return err
	}
	if b.launch.Protocol != 1 {
		return errors.New("unsupported protocol")
	}
	if len(os.Args) != 2 {
		return errors.New("usage: fungi open --exec fungi-wails -- PATH_TO_HTML_OR_JSON")
	}
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		return err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	b.launch.Source = path
	if filepath.Ext(path) == ".json" {
		b.form, err = os.ReadFile(path)
		if err != nil {
			return err
		}
		if !json.Valid(b.form) {
			return errors.New("invalid form JSON")
		}
	}
	handler, err := newAssets(path, len(b.form) > 0)
	if err != nil {
		return err
	}
	defer handler.root.Close()
	m := menu.NewMenu()
	m.Append(menu.AppMenu())
	m.Append(menu.EditMenu())
	return wails.Run(&options.App{
		Title: "Fungi · Wails preview", Width: 820, Height: 780, MinWidth: 420, MinHeight: 420,
		BackgroundColour: options.NewRGB(24, 25, 28), Menu: m,
		AssetServer: &assetserver.Options{Handler: handler}, Bind: []interface{}{b},
		Logger: stderrLogger{}, LogLevel: logger.ERROR, LogLevelProduction: logger.ERROR,
		OnBeforeClose: b.beforeClose,
		OnDomReady:    func(context.Context) { b.mu.Lock(); defer b.mu.Unlock(); _ = b.signalReady() },
		Windows:       &windows.Options{WebviewUserDataPath: filepath.Join(b.launch.OutputDir, "webview")},
	})
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
