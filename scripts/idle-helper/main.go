// Development fixture: acknowledge a draft, idle without polling, then submit.
// Build separately; it is not a Fungi runtime dependency.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	dec, enc := json.NewDecoder(os.Stdin), json.NewEncoder(os.Stdout)
	var launch struct {
		Protocol  int    `json:"protocol"`
		Session   string `json:"session"`
		OutputDir string `json:"outputDir"`
		Input     struct {
			Seconds int    `json:"seconds"`
			Marker  string `json:"marker"`
		} `json:"input"`
		Draft any `json:"draft"`
	}
	if err := dec.Decode(&launch); err != nil {
		return err
	}
	if launch.Protocol != 2 || launch.Input.Seconds < 1 || launch.Input.Seconds > 86400 || launch.Input.Marker == "" {
		return fmt.Errorf("expected protocol 2, input seconds 1–86400 and a nonempty marker")
	}
	request := func(value any) error {
		if err := enc.Encode(value); err != nil {
			return err
		}
		var ack struct {
			OK    bool `json:"ok"`
			Error any  `json:"error"`
		}
		if err := dec.Decode(&ack); err != nil {
			return err
		}
		if !ack.OK {
			return fmt.Errorf("request rejected: %v", ack.Error)
		}
		return nil
	}
	start := time.Now()
	data := map[string]any{"marker": launch.Input.Marker, "startedAt": start.UTC(), "requestedSeconds": launch.Input.Seconds}
	if err := request(map[string]any{"type": "draft", "id": "draft", "data": data}); err != nil {
		return err
	}
	if err := enc.Encode(map[string]any{"type": "ready"}); err != nil {
		return err
	}
	// No status checks, heartbeats, output or model work during this timer.
	time.Sleep(time.Duration(launch.Input.Seconds) * time.Second)
	data["finishedAt"] = time.Now().UTC()
	data["elapsedSeconds"] = time.Since(start).Seconds()
	data["session"] = launch.Session
	data["restoredDraft"] = launch.Draft
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(launch.OutputDir, "idle.json"), append(b, '\n'), 0600); err != nil {
		return err
	}
	return request(map[string]any{
		"type": "result", "id": "result", "outcome": "submitted", "data": data,
		"artifacts": []any{map[string]any{"key": "timing", "name": "idle.json", "path": "idle.json", "mediaType": "application/json"}},
	})
}
