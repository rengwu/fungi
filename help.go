package main

import (
	"embed"
	"io"
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
			_, err = io.WriteString(out, formatHelp(string(body)))
			return err
		}
	}
	return fail("invalid_arguments", "unknown help topic; choose "+strings.Join(manualTopics(), ", "))
}
