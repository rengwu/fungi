// Deliberately broken helpers for cross-platform protocol regression tests.
// This fixture is not shipped in the portable distribution.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(2)
	}
	switch os.Args[1] {
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
		os.Exit(2)
	}
}
