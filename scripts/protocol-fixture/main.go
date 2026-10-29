// Deliberately broken helpers for cross-platform protocol regression tests.
// This fixture is not shipped in the portable distribution.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "protocol fixture: unknown mode")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "update-session":
		// A bounded synthetic session controlled by a test-owned file. This is
		// never a human-facing waiter and is not part of the shipped helper API.
		if len(os.Args) != 3 {
			panic("update-session needs a gate file")
		}
		dec, enc := json.NewDecoder(os.Stdin), json.NewEncoder(os.Stdout)
		var launch struct {
			Session   string `json:"session"`
			Draft     any    `json:"draft"`
			OutputDir string `json:"outputDir"`
		}
		if err := dec.Decode(&launch); err != nil {
			panic(err)
		}
		data := launch.Draft
		if data == nil {
			data = map[string]string{"notes": "saved before update"}
		}
		enc.Encode(map[string]any{"type": "draft", "id": "draft", "data": data})
		var ack struct {
			OK bool `json:"ok"`
		}
		if err := dec.Decode(&ack); err != nil || !ack.OK {
			panic("draft failed")
		}
		enc.Encode(map[string]any{"type": "ready"})
		if err := os.WriteFile(os.Args[2]+".ready", []byte(launch.Session), 0600); err != nil {
			panic(err)
		}
		deadline := time.Now().Add(60 * time.Second)
		for {
			gate, err := os.ReadFile(os.Args[2])
			if err == nil {
				if err = os.WriteFile(filepath.Join(launch.OutputDir, "note.txt"), []byte("artifact survives update"), 0600); err != nil {
					panic(err)
				}
				enc.Encode(map[string]any{"type": "result", "id": "result", "outcome": strings.TrimSpace(string(gate)), "data": data, "artifacts": []any{map[string]string{"key": "note", "name": "note.txt", "path": "note.txt", "mediaType": "text/plain"}}})
				if err = dec.Decode(&ack); err != nil || !ack.OK {
					panic("submission failed")
				}
				return
			}
			if !os.IsNotExist(err) || time.Now().After(deadline) {
				panic("fixture gate unavailable")
			}
			time.Sleep(20 * time.Millisecond)
		}
	case "queue":
		f, err := os.OpenFile(os.Getenv("FUNGI_TEST_NOTICE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(f, strings.Join(os.Args[1:], "\n"))
		f.Close()
		os.Exit(7) // Receiver accepted bytes but did not acknowledge success.
	case "exit":
		os.Exit(1)
	case "early-result":
		fmt.Println(`{"type":"ready"}`)
		fmt.Println(`{"type":"result","outcome":"submitted","data":{"drained":true}}`)
		os.Exit(7)
	case "never-read":
		time.Sleep(30 * time.Second)
	case "malformed", "closed-stdout":
		var launch any
		if err := json.NewDecoder(os.Stdin).Decode(&launch); err != nil {
			panic(err)
		}
		fmt.Println(`{"type":"ready"}`)
		if os.Args[1] == "malformed" {
			fmt.Println("not json")
		} else {
			os.Stdout.Close()
		}
		// Longer than the test deadline, so a pass proves owner termination
		// rather than the fixture eventually exiting by itself.
		time.Sleep(30 * time.Second)
	default:
		fmt.Fprintln(os.Stderr, "protocol fixture: unknown mode")
		os.Exit(2)
	}
}
