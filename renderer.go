package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func rendererRelativePath() string {
	if runtime.GOOS == "darwin" {
		return filepath.Join("renderer", "Fungi.app", "Contents", "MacOS", "fungi-renderer")
	}
	name := "fungi-renderer"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join("renderer", name)
}

func rendererCommand() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	name := filepath.Join(filepath.Dir(exe), rendererRelativePath())
	if dir := os.Getenv("FUNGI_RENDERER_DIR"); dir != "" {
		rel, _ := filepath.Rel("renderer", rendererRelativePath())
		name = filepath.Join(dir, rel)
	}
	name, err = filepath.Abs(name)
	if err != nil {
		return "", err
	}
	if info, err := os.Stat(name); err != nil || !info.Mode().IsRegular() {
		return "", fail("renderer_unavailable", "Fungi's window helper is missing; keep the complete platform package together. Development build: node scripts/package.mjs")
	}
	return name, nil
}

// This bounded, windowless check never installs a runtime or starts a user wait.
func checkRenderer(name string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, "--check-runtime")
	cmd.Env = childEnv()
	b, err := cmd.CombinedOutput()
	var reply struct {
		Protocol  int    `json:"protocol"`
		Available bool   `json:"available"`
		Runtime   string `json:"runtime"`
		Message   string `json:"message"`
	}
	valid := json.Unmarshal(b, &reply) == nil && reply.Protocol == 1
	if err == nil && valid && reply.Available {
		return map[string]any{"protocol": 1, "kind": "doctor", "renderer": "wails", "runtime": reply.Runtime, "available": true}, nil
	}
	text := strings.TrimSpace(string(b))
	if (valid && !reply.Available) || (runtime.GOOS == "linux" && strings.Contains(text, "error while loading shared libraries:")) {
		if valid {
			text = reply.Message
		}
		return nil, fail("runtime_missing", text+". Re-run the Fungi installer with --setup-runtime (Windows: -SetupRuntime), or install the system runtime yourself. Ordinary UI launches never install software.")
	}
	if ctx.Err() != nil {
		return nil, fail("renderer_check_failed", "window helper runtime check timed out")
	}
	return nil, fail("renderer_check_failed", fmt.Sprintf("window helper runtime check failed: %v; %s", err, text))
}
