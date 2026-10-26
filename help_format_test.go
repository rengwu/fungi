package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTerminalHelpFormatting(t *testing.T) {
	input := "# Fungi\n\nUse **local** pages\nand `fungi help`.\n\n## Commands\n\n| Command | Purpose |\n| --- | --- |\n| `open` | Show a page |\n| `collect` | Read a response |\n\n```sh\necho '`literal` **stars**' | cat\n  # keep code indentation\n```\n"
	want := "Fungi\n\nUse local pages and fungi help.\n\nCommands\n\n  Command  Purpose\n  open     Show a page\n  collect  Read a response\n\n  echo '`literal` **stars**' | cat\n    # keep code indentation\n"
	if got := formatHelp(input); got != want {
		t.Fatalf("terminal text:\n%s\nwant:\n%s", got, want)
	}
}

func TestTerminalHelpWrappingAndLongCommands(t *testing.T) {
	command := "curl https://example.test/" + strings.Repeat("a", 90) + " | sh && export PATH=\"$HOME/.fungi/bin:$PATH\""
	paragraph := strings.Repeat("A readable paragraph → ", 12)
	input := "# Overview\n\n" + paragraph + "\n\n| Mode | Meaning |\n| --- | --- |\n| `auto` | " + paragraph + " |\n\n```sh\n" + command + "\n```\n"
	output := formatHelp(input)
	if !strings.Contains(output, "  "+command+"\n") {
		t.Fatal("installation command was wrapped or modified")
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.TrimSpace(line) != command && utf8.RuneCountInString(line) > 80 {
			t.Fatalf("prose/table exceeds terminal width: %q", line)
		}
	}
}
