package convert

import (
	"strings"
	"testing"
)

func TestMdToTrac(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Headers",
			input:    "# Header 1\n## Header 2\n### Header 3",
			expected: "= Header 1 =\n== Header 2 ==\n=== Header 3 ===",
		},
		{
			name:     "Bold text",
			input:    "This is **bold** text",
			expected: "This is '''bold''' text",
		},
		{
			name:     "Italic text",
			input:    "This is *italic* text",
			expected: "This is ''italic'' text",
		},
		{
			name:     "Bold and italic",
			input:    "This is ***bold and italic*** text",
			expected: "This is '''''bold and italic''''' text",
		},
		{
			name:     "Strikethrough",
			input:    "This is ~~strikethrough~~ text",
			expected: "This is ~~strikethrough~~ text",
		},
		{
			name:     "Links",
			input:    "[Google](https://google.com)",
			expected: "[https://google.com Google]",
		},
		{
			name:     "Images",
			input:    "![Alt text](image.png)",
			expected: "[[Image(image.png, Alt text)]]",
		},
		{
			name:     "Inline code",
			input:    "This is `inline code`",
			expected: "This is `inline code`",
		},
		{
			name:     "Code block",
			input:    "```javascript\nconsole.log('hello');\n```",
			expected: "{{{\n#!javascript\nconsole.log('hello');\n}}}",
		},
		{
			name:     "Code block (no language)",
			input:    "```\nfoo\nbar\n```",
			expected: "{{{\nfoo\nbar\n}}}",
		},
		{
			name:     "Code block (http)",
			input:    "```http\nGET /api HTTP/1.1\n```",
			expected: "{{{\n#!text\nGET /api HTTP/1.1\n}}}",
		},
		{
			name:     "Code block (json)",
			input:    "```json\n{\"a\":1}\n```",
			expected: "{{{\n#!javascript\n{\"a\":1}\n}}}",
		},
		{
			name:     "Multiple code blocks",
			input:    "```js\n1\n```\n\n```python\n2\n```",
			expected: "{{{\n#!js\n1\n}}}\n\n{{{\n#!python\n2\n}}}",
		},
		{
			name:     "Unordered list",
			input:    "- Item 1\n- Item 2\n- Item 3",
			expected: " * Item 1\n * Item 2\n * Item 3",
		},
		{
			name:     "Ordered list",
			input:    "1. Item 1\n2. Item 2\n3. Item 3",
			expected: " 1. Item 1\n 1. Item 2\n 1. Item 3",
		},
		{
			name:     "Nested list",
			input:    "  - Nested 1\n    - Nested 2",
			expected: "  * Nested 1\n    * Nested 2",
		},
		{
			name:     "Checkboxes",
			input:    "- [x] Checked\n- [ ] Unchecked",
			expected: " * [X] Checked\n * [ ] Unchecked",
		},
		{
			name:     "Table",
			input:    "| Col1 | Col2 |\n|------|------|\n| Data1 | Data2 |",
			expected: "|| Col1 || Col2 ||\n|| Data1 || Data2 ||",
		},
		{
			name:     "Blockquote",
			input:    "> Blockquote",
			expected: " Blockquote",
		},
		{
			name:     "Horizontal rule",
			input:    "---",
			expected: "----",
		},
		{
			name:     "Multiple formatting",
			input:    "**bold** and *italic* and `code`",
			expected: "'''bold''' and ''italic'' and `code`",
		},
		{
			name:     "Code block (javascript, hello world)",
			input:    "```javascript\nfunction hello() {\n    console.log(\"Hello, World!\");\n}\n\nhello();\n```",
			expected: "{{{\n#!javascript\nfunction hello() {\n    console.log(\"Hello, World!\");\n}\n\nhello();\n}}}",
		},
		{
			name:     "Code block (python, greet)",
			input:    "```python\ndef greet(name):\n    return f\"こんにちは、{name}さん！\"\n\nprint(greet(\"世界\"))\n```",
			expected: "{{{\n#!python\ndef greet(name):\n    return f\"こんにちは、{name}さん！\"\n\nprint(greet(\"世界\"))\n}}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MdToTrac(tt.input)
			result = strings.TrimSpace(result)
			expected := strings.TrimSpace(tt.expected)
			if result != expected {
				t.Errorf("MdToTrac() = %q, want %q", result, expected)
			}
		})
	}
}

func TestComplexMarkdown(t *testing.T) {
	input := "# API Documentation\n" +
		"\n" +
		"## Overview\n" +
		"\n" +
		"This API provides **CRUD operations** for user management.\n" +
		"\n" +
		"### Authentication\n" +
		"\n" +
		"All endpoints require `Authorization` header:\n" +
		"\n" +
		"```http\n" +
		"GET /api/users HTTP/1.1\n" +
		"Authorization: Bearer token123\n" +
		"```\n" +
		"\n" +
		"### Endpoints\n" +
		"\n" +
		"#### GET /users\n" +
		"\n" +
		"Returns a list of users.\n" +
		"\n" +
		"**Parameters:**\n" +
		"\n" +
		"| Parameter | Type | Required | Description |\n" +
		"|-----------|------|----------|-------------|\n" +
		"| page | integer | No | Page number |\n" +
		"| limit | integer | No | Items per page |\n" +
		"\n" +
		"**Response:**\n" +
		"\n" +
		"```json\n" +
		"{\n" +
		"  \"users\": [\n" +
		"    {\n" +
		"      \"id\": 1,\n" +
		"      \"name\": \"John Doe\"\n" +
		"    }\n" +
		"  ]\n" +
		"}\n" +
		"```\n" +
		"\n" +
		"## Tasks\n" +
		"\n" +
		"- [x] Implement authentication\n" +
		"- [ ] Add rate limiting\n" +
		"- [ ] Write tests\n" +
		"\n" +
		"## Notes\n" +
		"\n" +
		"> **Important:** Always use HTTPS in production.\n" +
		"\n" +
		"---\n" +
		"\n" +
		"For more information, visit [our website](https://example.com)."

	result := MdToTrac(input)
	if !strings.Contains(result, "= API Documentation =") {
		t.Error("Header conversion failed")
	}
	if !strings.Contains(result, "'''CRUD operations'''") {
		t.Error("Bold conversion failed")
	}
	if !strings.Contains(result, "|| Parameter || Type ||") {
		t.Error("Table conversion failed")
	}
	if !strings.Contains(result, " * [X] Implement authentication") {
		t.Error("Checkbox conversion failed")
	}
	if !strings.Contains(result, "[https://example.com our website]") {
		t.Error("Link conversion failed")
	}
}
