package main

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// The embedded and packaged manuals share one source. Help also works when
// only the core executable is available, without a webview or adjacent files.
//
//go:embed docs/manual/*.md
var manual embed.FS

func manualTopics() []string {
	entries, _ := manual.ReadDir("docs/manual")
	topics := make([]string, 0, len(entries))
	for _, entry := range entries {
		topics = append(topics, strings.TrimSuffix(entry.Name(), ".md"))
	}
	return topics
}

func help(args []string, out io.Writer) error {
	topic := "start"
	if len(args) > 1 {
		return fail("invalid_arguments", "help accepts at most one topic")
	}
	if len(args) == 1 {
		topic = args[0]
	}
	for _, name := range manualTopics() {
		if topic == name {
			body, err := manual.ReadFile("docs/manual/" + name + ".md")
			if err != nil {
				return err
			}
			_, err = out.Write(body)
			return err
		}
	}
	return fail("invalid_arguments", "unknown help topic; choose "+strings.Join(manualTopics(), ", "))
}

type example struct {
	Name      string   `json:"name"`
	Purpose   string   `json:"purpose"`
	Source    string   `json:"source"`
	Input     string   `json:"input,omitempty"`
	Available bool     `json:"available"`
	Argv      []string `json:"argv,omitempty"`
}

func examples(executable string) []example {
	items := []example{
		{Name: "reply", Purpose: "Parameterised HTML response, no draft", Source: "examples/reply/index.html", Input: "examples/reply/input.json"},
		{Name: "compare", Purpose: "Canvas comparison, PNG export and optional recovery", Source: "examples/compare/index.html"},
		{Name: "form", Purpose: "Experimental flat JSON form", Source: "examples/feedback.json"},
		{Name: "review", Purpose: "Animation timing, notes and image attachment", Source: "examples/review.html"},
	}
	for i := range items {
		item := &items[i]
		item.Source = filepath.Join(filepath.Dir(executable), filepath.FromSlash(item.Source))
		info, err := os.Stat(item.Source)
		item.Available = err == nil && info.Mode().IsRegular()
		if item.Input != "" {
			item.Input = filepath.Join(filepath.Dir(executable), filepath.FromSlash(item.Input))
			info, err = os.Stat(item.Input)
			item.Available = item.Available && err == nil && info.Mode().IsRegular()
		}
		if item.Available {
			item.Argv = []string{executable, "open", item.Source}
			if item.Input != "" {
				item.Argv = append(item.Argv, "--input", item.Input)
			}
			item.Argv = append(item.Argv, "--delivery", "auto")
		}
	}
	return items
}
