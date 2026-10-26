package main

import (
	"strings"
	"unicode/utf8"
)

// Format the controlled manual syntax as plain terminal text. This is not a
// general Markdown renderer: our manuals use headings, paragraphs, two-column
// tables and fenced examples. Examples stay verbatim, including long commands.
// Plain output also reads cleanly in agent captures and redirected files.
func formatHelp(markdown string) string {
	lines := strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n")
	var blocks []string
	for i := 0; i < len(lines); {
		line := strings.TrimSpace(lines[i])
		switch {
		case line == "":
			i++
		case strings.HasPrefix(line, "```"):
			i++
			var code []string
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				code = append(code, "  "+lines[i])
				i++
			}
			blocks = append(blocks, strings.Join(code, "\n"))
			if i < len(lines) {
				i++
			}
		case strings.HasPrefix(line, "#"):
			blocks = append(blocks, helpInline(strings.TrimLeft(line, "# ")))
			i++
		case strings.HasPrefix(line, "|"):
			var rows [][2]string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
				cells := strings.SplitN(strings.Trim(strings.TrimSpace(lines[i]), "|"), "|", 2)
				if len(cells) == 2 && strings.Trim(cells[0], " :-") != "" {
					rows = append(rows, [2]string{helpInline(strings.TrimSpace(cells[0])), helpInline(strings.TrimSpace(cells[1]))})
				}
				i++
			}
			blocks = append(blocks, helpTable(rows))
		default:
			var paragraph []string
			for i < len(lines) {
				line = strings.TrimSpace(lines[i])
				if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "|") || strings.HasPrefix(line, "```") {
					break
				}
				paragraph = append(paragraph, line)
				i++
			}
			blocks = append(blocks, strings.Join(helpWrap(helpInline(strings.Join(paragraph, " ")), 80), "\n"))
		}
	}
	return strings.Join(blocks, "\n\n") + "\n"
}

func helpInline(text string) string {
	return strings.NewReplacer("`", "", "**", "").Replace(text)
}

func helpWrap(text string, width int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		if line != "" && utf8.RuneCountInString(line)+1+utf8.RuneCountInString(word) > width {
			lines = append(lines, line)
			line = ""
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	return append(lines, line)
}

func helpTable(rows [][2]string) string {
	width := 0
	for _, row := range rows {
		width = max(width, utf8.RuneCountInString(row[0]))
	}
	width = min(width, 36)
	var lines []string
	for _, row := range rows {
		left, right := helpWrap(row[0], width), helpWrap(row[1], 80-width-4)
		for i := 0; i < max(len(left), len(right)); i++ {
			a, b := "", ""
			if i < len(left) {
				a = left[i]
			}
			if i < len(right) {
				b = right[i]
			}
			padding := strings.Repeat(" ", max(0, width-utf8.RuneCountInString(a))+2)
			lines = append(lines, strings.TrimRight("  "+a+padding+b, " "))
		}
	}
	return strings.Join(lines, "\n")
}
