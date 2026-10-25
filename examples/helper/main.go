// A tiny native executable implementing the same protocol as the window renderer.
// Build: go build -o /tmp/fungi-helper ./examples/helper
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	var launch struct {
		Protocol  int    `json:"protocol"`
		OutputDir string `json:"outputDir"`
		Input     any    `json:"input"`
		Draft     any    `json:"draft"`
	}
	if err := dec.Decode(&launch); err != nil {
		panic(err)
	}
	if launch.Protocol != 1 {
		fmt.Fprintln(os.Stderr, "unsupported protocol")
		os.Exit(1)
	}
	enc.Encode(map[string]any{"type": "ready"})
	if len(os.Args) > 1 && os.Args[1] == "wait" {
		seconds := 10
		if len(os.Args) > 2 {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil || n < 1 || n > 600 {
				panic("fixture wait must be 1–600 seconds")
			}
			seconds = n
		}
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	b := []byte("A native helper can return files through Fungi.\n")
	if err := os.WriteFile(filepath.Join(launch.OutputDir, "note.txt"), b, 0600); err != nil {
		panic(err)
	}
	enc.Encode(map[string]any{"type": "draft", "id": "draft", "data": map[string]any{"notes": "Native draft"}})
	var ack struct {
		OK    bool `json:"ok"`
		Error any  `json:"error"`
	}
	if err := dec.Decode(&ack); err != nil || !ack.OK {
		panic(fmt.Sprint(err, ack.Error))
	}
	enc.Encode(map[string]any{"type": "result", "id": "result", "outcome": "submitted", "data": map[string]any{"notes": "Native response", "input": launch.Input, "restoredDraft": launch.Draft}, "artifacts": []any{map[string]any{"key": "note", "name": "note.txt", "path": "note.txt", "mediaType": "text/plain"}}})
	if err := dec.Decode(&ack); err != nil || !ack.OK {
		panic(fmt.Sprint(err, ack.Error))
	}
}
