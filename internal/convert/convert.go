package convert

import (
	"fmt"
	"regexp"
	"strings"
)

type codeBlock struct {
	placeholder string
	content     string
}

func MdToTrac(input string) string {
	// すべての改行をLFに統一
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")

	input = escapeMarkdownSpecials(input)
	codeBlocks, input := extractAndReplaceCodeBlocks(input)
	input = convertTables(input)
	input = convertImages(input)
	input = convertLinks(input)
	input = convertHeaders(input)
	input = restoreCodeBlocks(input, codeBlocks) // ★ここで復元
	input = convertTextFormatting(input)
	input = convertLists(input)
	input = convertBlockquotes(input)
	input = convertHorizontalRules(input)
	input = convertFootnotes(input)
	input = convertBadges(input)
	input = unescapeMarkdownSpecials(input)
	input = normalizeNewlines(input)
	return strings.TrimSpace(input)
}

func escapeMarkdownSpecials(input string) string {
	input = strings.ReplaceAll(input, `\\*`, "ESCAPED_ASTERISK")
	input = strings.ReplaceAll(input, `\\_`, "ESCAPED_UNDERSCORE")
	input = strings.ReplaceAll(input, `\\~`, "ESCAPED_TILDE")
	return input
}

func unescapeMarkdownSpecials(input string) string {
	input = strings.ReplaceAll(input, "ESCAPED_ASTERISK", "*")
	input = strings.ReplaceAll(input, "ESCAPED_UNDERSCORE", "_")
	input = strings.ReplaceAll(input, "ESCAPED_TILDE", "~")
	return input
}

func extractAndReplaceCodeBlocks(input string) ([]codeBlock, string) {
	var codeBlocks []codeBlock
	reCodeBlockAll := regexp.MustCompile("(?s)```([a-zA-Z0-9+#-]*)\\n(.*?)\\n?```")
	idx := 0
	input = reCodeBlockAll.ReplaceAllStringFunc(input, func(s string) string {
		m := reCodeBlockAll.FindStringSubmatch(s)
		lang := m[1]
		code := m[2]
		placeholder := fmt.Sprintf("[[[CODEBLOCK_PLACEHOLDER_%d]]]", idx)
		var content string
		if lang == "http" {
			content = "{{{\n#!text\n" + code + "\n}}}"
		} else if lang == "json" {
			content = "{{{\n#!javascript\n" + code + "\n}}}"
		} else if lang != "" {
			content = "{{{\n#!" + lang + "\n" + code + "\n}}}"
		} else {
			content = "{{{\n" + code + "\n}}}"
		}
		codeBlocks = append(codeBlocks, codeBlock{placeholder, content})
		idx++
		return placeholder
	})
	return codeBlocks, input
}

func restoreCodeBlocks(input string, codeBlocks []codeBlock) string {
	for _, cb := range codeBlocks {
		input = strings.ReplaceAll(input, cb.placeholder, cb.content)
	}
	return input
}

func convertTables(input string) string {
	reTable := regexp.MustCompile(`(?m)((?:\|[^\n]+\|\n?)+)`)
	return reTable.ReplaceAllStringFunc(input, func(table string) string {
		lines := regexp.MustCompile(`\r?\n`).Split(strings.TrimSpace(table), -1)
		var out string
		reHeaderSeparator := regexp.MustCompile(`^\|(?:\s*:?-+:?\s*\|)+\s*$`)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || reHeaderSeparator.MatchString(line) {
				continue // 区切り行をスキップ
			}
			line = strings.TrimPrefix(line, "|")
			line = strings.TrimSuffix(line, "|")
			cells := regexp.MustCompile(`\|`).Split(line, -1)
			for i := range cells {
				cells[i] = strings.TrimSpace(cells[i])
			}
			out += "|| " + strings.Join(cells, " || ") + " ||\n"
		}
		return strings.TrimRight(out, "\n")
	})
}

func convertImages(input string) string {
	reImgLink := regexp.MustCompile(`\[!\[(.*?)\]\((.*?)\)\]\((.*?)\)`)
	input = reImgLink.ReplaceAllString(input, "[$3 [[Image($2, $1)]]]")
	reImgSimple := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	return reImgSimple.ReplaceAllString(input, "[[Image($2, $1)]]")
}

func convertLinks(input string) string {
	reLink := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	return reLink.ReplaceAllString(input, "[$2 $1]")
}

func convertHeaders(input string) string {
	reHeader := regexp.MustCompile(`(?m)^(#{1,6})\s*(.+?)(?:\s*#+\s*)?$`)
	return reHeader.ReplaceAllStringFunc(input, func(s string) string {
		m := reHeader.FindStringSubmatch(s)
		level := len(m[1])
		title := strings.TrimSpace(m[2])
		eq := strings.Repeat("=", level)
		return fmt.Sprintf("%s %s %s", eq, title, eq)
	})
}

func convertTextFormatting(input string) string {
	reStrike := regexp.MustCompile(`~~(.*?)~~`)
	input = reStrike.ReplaceAllString(input, "~~$1~~")
	reBoldItalic1 := regexp.MustCompile(`\*\*\*(.*?)\*\*\*`)
	input = reBoldItalic1.ReplaceAllString(input, "'''''$1'''''")
	reBoldItalic2 := regexp.MustCompile(`___(.*?)___`)
	input = reBoldItalic2.ReplaceAllString(input, "'''''$1'''''")
	reBold1 := regexp.MustCompile(`\*\*(.*?)\*\*`)
	input = reBold1.ReplaceAllString(input, "'''$1'''")
	reItalic1 := regexp.MustCompile(`\*(.*?)\*`)
	input = reItalic1.ReplaceAllString(input, "''$1''")
	reItalic2 := regexp.MustCompile(`_(.*?)_`)
	input = reItalic2.ReplaceAllString(input, "''$1''")
	reInlineCode := regexp.MustCompile("`([^`]+)`")
	input = reInlineCode.ReplaceAllString(input, "`$1`")
	return input
}

func convertLists(input string) string {
	reCheckboxChecked := regexp.MustCompile(`(?m)^(\s*)[-*]\s+\[x\]\s+(.+)$`)
	input = reCheckboxChecked.ReplaceAllString(input, "$1 * [X] $2")
	reCheckboxUnchecked := regexp.MustCompile(`(?m)^(\s*)[-*]\s+\[\s\]\s+(.+)$`)
	input = reCheckboxUnchecked.ReplaceAllString(input, "$1 * [ ] $2")
	reSubUL := regexp.MustCompile(`(?m)^(\s+)[-*]\s+(.+)$`)
	input = reSubUL.ReplaceAllStringFunc(input, func(s string) string {
		matches := reSubUL.FindStringSubmatch(s)
		indent := matches[1]
		return indent + "*" + " " + matches[2]
	})
	reSubOL := regexp.MustCompile(`(?m)^(\s+)\d+\.\s+(.+)$`)
	input = reSubOL.ReplaceAllStringFunc(input, func(s string) string {
		matches := reSubOL.FindStringSubmatch(s)
		indent := matches[1]
		return indent + "1." + " " + matches[2]
	})
	reUL := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
	input = reUL.ReplaceAllString(input, " * $1")
	reOL := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	input = reOL.ReplaceAllString(input, " 1. $1")
	return input
}

func convertBlockquotes(input string) string {
reNestedQuote := regexp.MustCompile(`(?m)^>\s*>\s*(.*)$`)
input = reNestedQuote.ReplaceAllString(input, "  $1")

reQuote := regexp.MustCompile(`(?m)^>\s*(.*)$`)
input = reQuote.ReplaceAllString(input, " $1")
	return input
}

func convertHorizontalRules(input string) string {
	reHR1 := regexp.MustCompile(`(?m)^-{3,}\s*$`)
	input = reHR1.ReplaceAllString(input, "----")
	reHR2 := regexp.MustCompile(`(?m)^\*{3,}\s*$`)
	input = reHR2.ReplaceAllString(input, "----")
	return input
}

func convertFootnotes(input string) string {
	reFootnoteRef := regexp.MustCompile(`\[\^([^\]]+)\]`)
	input = reFootnoteRef.ReplaceAllString(input, "^$1^")
	reFootnoteDef := regexp.MustCompile(`(?m)^\[\^([^\]]+)\]:\s*(.+)$`)
	input = reFootnoteDef.ReplaceAllString(input, "[[FootNote($1,$2)]]")
	return input
}

func convertBadges(input string) string {
	reBadge := regexp.MustCompile(`\[!\[([^\]]*)\]\(([^)]+)\)\]\(([^)]+)\)`)
	return reBadge.ReplaceAllString(input, "[$3 [[Image($2, $1)]]]")
}

func normalizeNewlines(input string) string {
	reMultipleNewlines := regexp.MustCompile(`\n{3,}`)
	return reMultipleNewlines.ReplaceAllString(input, "\n\n")
}
