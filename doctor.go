package main

import (
	"io"
	"os"
	"path/filepath"
)

// Exercise the engine's storage primitives without launching a helper.
func doctor(out io.Writer) error {
	root, err := dataRoot()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(root, 0700); err != nil {
		return err
	}
	dir, err := os.MkdirTemp(root, ".doctor-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	unlock, err := lockSession(filepath.Join(dir, "lock"))
	if err != nil {
		return err
	}
	defer unlock()
	if release, err := lockSession(filepath.Join(dir, "lock")); err == nil {
		release()
		return fail("storage_check_failed", "session storage does not enforce exclusive locks")
	} else if asProblem(err).Code != "session_owned" {
		return err
	}
	file := filepath.Join(dir, "check.json")
	if err = atomicJSON(file, map[string]int{"protocol": protocolVersion}); err != nil {
		return err
	}
	var saved struct {
		Protocol int `json:"protocol"`
	}
	if err = readJSON(file, &saved); err != nil {
		return err
	}
	if saved.Protocol != protocolVersion {
		return fail("storage_check_failed", "saved engine probe could not be read back")
	}
	return emit(out, map[string]any{"protocol": protocolVersion, "kind": "doctor", "engine": "headless", "available": true, "stateDir": root, "checks": []string{"atomicStorage", "exclusiveLock"}})
}
