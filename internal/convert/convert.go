package convert

import (
	"fmt"
	"regexp"
)

func MdToTrac(input string) string {
	// Header conversion
	reHeader := regexp.MustCompile(`(?m)^(#{1,6})\s*(.+)$`)
	input = reHeader.ReplaceAllStringFunc(input, func(s string) string {
		m := reHeader.FindStringSubmatch(s)
		level := len(m[1])
		eq := ""
		for i := 0; i < level; i++ {
			eq += "="
		}
		return fmt.Sprintf("%s %s %s", eq, m[2], eq)
	})

	// Bold
	reBold := regexp.MustCompile(`\*\*(.*?)\*\*`)
	input = reBold.ReplaceAllString(input, "'''$1'''")

	// Italic
	reItalic := regexp.MustCompile(`\*(.*?)\*`)
	input = reItalic.ReplaceAllString(input, "''$1''")

	// Links
	reLink := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	input = reLink.ReplaceAllString(input, "[$2 $1]")

	// Images
	reImg := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	input = reImg.ReplaceAllString(input, "[[Image($2, $1)]]")

	// Code block (simple)
	reCode := regexp.MustCompile("(?s)```(.*?)```")
	input = reCode.ReplaceAllString(input, "{{{$1}}}")

	// Unordered list
	reUL := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
	input = reUL.ReplaceAllString(input, " * $1")

	// Ordered list
	reOL := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	input = reOL.ReplaceAllString(input, " 1. $1")

	return input
}
