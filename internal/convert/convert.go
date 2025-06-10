package convert

import (
	"fmt"
	"regexp"
	"strings"
)

func MdToTrac(input string) string {
	// バックスラッシュエスケープの処理（最初に実行）
	input = strings.ReplaceAll(input, `\*`, "ESCAPED_ASTERISK")
	input = strings.ReplaceAll(input, `\_`, "ESCAPED_UNDERSCORE")
	input = strings.ReplaceAll(input, `\~`, "ESCAPED_TILDE")

	// Table: Markdown → Trac Wiki
	reTable := regexp.MustCompile(`(?m)((?:\|[^\n]+\|\n?)+)`)
	input = reTable.ReplaceAllStringFunc(input, func(table string) string {
		lines := regexp.MustCompile(`\r?\n`).Split(strings.TrimSpace(table), -1)
		var out string
		reHeaderSeparator := regexp.MustCompile(`^\|[\s:-]+\|$`)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// ヘッダー区切り行（|----|----|）はスキップ
			if reHeaderSeparator.MatchString(line) {
				continue
			}
			// 先頭・末尾の|を||に、セル区切り|を||に
			line = regexp.MustCompile(`^\|`).ReplaceAllString(line, "||")
			line = regexp.MustCompile(`\|$`).ReplaceAllString(line, "||")
			line = regexp.MustCompile(`\|`).ReplaceAllString(line, "||")
			out += line + "\n"
		}
		return out
	})

	// Image link: [![alt](src)](url) → [url [[Image(src, alt)]]]
	reImgLink := regexp.MustCompile(`\[\!\[(.*?)\]\((.*?)\)\]\((.*?)\)`)
	input = reImgLink.ReplaceAllString(input, "[$3 [[Image($2, $1)]]]")

	// Code block (with or without language) → {{{ ... }}}
	reCodeBlock := regexp.MustCompile("(?s)```[a-zA-Z0-9+#-]*\\n(.*?)\\n?```")
	input = reCodeBlock.ReplaceAllString(input, "{{{\n$1\n}}}")

	// インラインコード
	reInlineCode := regexp.MustCompile("`([^`]+)`")
	input = reInlineCode.ReplaceAllString(input, "`$1`")

	// Header conversion
	reHeader := regexp.MustCompile(`(?m)^(#{1,6})\s*(.+?)(?:\s*#+\s*)?$`)
	input = reHeader.ReplaceAllStringFunc(input, func(s string) string {
		m := reHeader.FindStringSubmatch(s)
		level := len(m[1])
		title := strings.TrimSpace(m[2])
		eq := strings.Repeat("=", level)
		return fmt.Sprintf("%s %s %s", eq, title, eq)
	})

	// 取り消し線
	reStrike := regexp.MustCompile(`~~(.*?)~~`)
	input = reStrike.ReplaceAllString(input, "~~$1~~")

	// 太字かつ斜体 (***text*** or ___text___)
	reBoldItalic1 := regexp.MustCompile(`\*\*\*(.*?)\*\*\*`)
	input = reBoldItalic1.ReplaceAllString(input, "'''''$1'''''")
	reBoldItalic2 := regexp.MustCompile(`___(.*?)___`)
	input = reBoldItalic2.ReplaceAllString(input, "'''''$1'''''")

	// Bold (**text** or __text__)
	reBold1 := regexp.MustCompile(`\*\*(.*?)\*\*`)
	input = reBold1.ReplaceAllString(input, "'''$1'''")
	reBold2 := regexp.MustCompile(`__(.*?)__`)
	input = reBold2.ReplaceAllString(input, "'''$1'''")

	// Italic (*text* or _text_)
	reItalic1 := regexp.MustCompile(`\*(.*?)\*`)
	input = reItalic1.ReplaceAllString(input, "''$1''")
	reItalic2 := regexp.MustCompile(`_(.*?)_`)
	input = reItalic2.ReplaceAllString(input, "''$1''")

	// Links
	reLink := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	input = reLink.ReplaceAllString(input, "[$2 $1]")

	// Images
	reImg := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	input = reImg.ReplaceAllString(input, "[[Image($2, $1)]]")

	// チェックリスト
	reCheckboxChecked := regexp.MustCompile(`(?m)^(\s*)[-*]\s+\[x\]\s+(.+)$`)
	input = reCheckboxChecked.ReplaceAllString(input, "$1 * [X] $2")
	reCheckboxUnchecked := regexp.MustCompile(`(?m)^(\s*)[-*]\s+\[\s\]\s+(.+)$`)
	input = reCheckboxUnchecked.ReplaceAllString(input, "$1 * [ ] $2")

	// ネストしたリストの処理
	reSubUL := regexp.MustCompile(`(?m)^(\s{2,})[-*]\s+(.+)$`)
	input = reSubUL.ReplaceAllStringFunc(input, func(s string) string {
		matches := reSubUL.FindStringSubmatch(s)
		indent := len(matches[1]) / 2
		bullet := strings.Repeat(" ", indent) + " *"
		return bullet + " " + matches[2]
	})

	reSubOL := regexp.MustCompile(`(?m)^(\s{2,})\d+\.\s+(.+)$`)
	input = reSubOL.ReplaceAllStringFunc(input, func(s string) string {
		matches := reSubOL.FindStringSubmatch(s)
		indent := len(matches[1]) / 2
		bullet := strings.Repeat(" ", indent) + " 1."
		return bullet + " " + matches[2]
	})

	// Unordered list (通常レベル)
	reUL := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
	input = reUL.ReplaceAllString(input, " * $1")

	// Ordered list (通常レベル)
	reOL := regexp.MustCompile(`(?m)^\d+\.\s+(.+)$`)
	input = reOL.ReplaceAllString(input, " 1. $1")

	// 引用の処理
	reQuote := regexp.MustCompile(`(?m)^>\s*(.*)$`)
	input = reQuote.ReplaceAllString(input, " $1")

	// ネストした引用の処理
	reNestedQuote := regexp.MustCompile(`(?m)^>\s*>\s*(.*)$`)
	input = reNestedQuote.ReplaceAllString(input, "  $1")

	// 水平線
	reHR1 := regexp.MustCompile(`(?m)^-{3,}\s*$`)
	input = reHR1.ReplaceAllString(input, "----")
	reHR2 := regexp.MustCompile(`(?m)^\*{3,}\s*$`)
	input = reHR2.ReplaceAllString(input, "----")

	// 脚注の処理
	reFootnoteRef := regexp.MustCompile(`\[\^([^\]]+)\]`)
	input = reFootnoteRef.ReplaceAllString(input, "^$1^")

	// 脚注定義の処理
	reFootnoteDef := regexp.MustCompile(`(?m)^\[\^([^\]]+)\]:\s*(.+)$`)
	input = reFootnoteDef.ReplaceAllString(input, "[[FootNote($1,$2)]]")

	// HTTPヘッダーブロックの処理（特殊なコードブロック形式）
	reHTTPBlock := regexp.MustCompile("(?s)```http\\n(.*?)\\n```")
	input = reHTTPBlock.ReplaceAllString(input, "{{{\n#!text\n$1\n}}}")

	// JSONコードブロックの処理
	reJSONBlock := regexp.MustCompile("(?s)```json\\n(.*?)\\n```")
	input = reJSONBlock.ReplaceAllString(input, "{{{\n#!javascript\n$1\n}}}")

	// バッジ（[![...](...)](...)形式）の処理
	reBadge := regexp.MustCompile(`\[\!\[([^\]]*)\]\(([^)]+)\)\]\(([^)]+)\)`)
	input = reBadge.ReplaceAllString(input, "[$3 [[Image($2, $1)]]]")

	// エスケープされた文字を元に戻す
	input = strings.ReplaceAll(input, "ESCAPED_ASTERISK", "*")
	input = strings.ReplaceAll(input, "ESCAPED_UNDERSCORE", "_")
	input = strings.ReplaceAll(input, "ESCAPED_TILDE", "~")

	// 複数の連続する空行を1つの空行に変換
	reMultipleNewlines := regexp.MustCompile(`\n{3,}`)
	input = reMultipleNewlines.ReplaceAllString(input, "\n\n")

	return strings.TrimSpace(input)
}
