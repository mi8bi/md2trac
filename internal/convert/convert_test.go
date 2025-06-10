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
			input:    "- Item 1\n  - Sub item 1\n  - Sub item 2\n- Item 2",
			expected: " * Item 1\n  * Sub item 1\n  * Sub item 2\n * Item 2",
		},
		{
			name:     "Checkboxes",
			input:    "- [x] Completed task\n- [ ] Incomplete task",
			expected: " * [X] Completed task\n * [ ] Incomplete task",
		},
		{
			name:     "Table",
			input:    "| Col1 | Col2 |\n|------|------|\n| Data1 | Data2 |",
			expected: "|| Col1 || Col2 ||\n|| Data1 || Data2 ||",
		},
		{
			name:     "Blockquote",
			input:    "> This is a quote\n> Second line",
			expected: " This is a quote\n Second line",
		},
		{
			name:     "Horizontal rule",
			input:    "---",
			expected: "----",
		},
		{
			name:     "Multiple formatting",
			input:    "# Title\n\nThis is **bold** and *italic* text with `code`.\n\n- List item\n- Another item",
			expected: "= Title =\n\nThis is '''bold''' and ''italic'' text with `code`.\n\n * List item\n * Another item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MdToTrac(tt.input)
			// 空白の違いを正規化して比較
			result = strings.TrimSpace(result)
			expected := strings.TrimSpace(tt.expected)
			
			if result != expected {
				t.Errorf("MdToTrac() = %q, want %q", result, expected)
			}
		})
	}
}

func TestComplexMarkdown(t *testing.T) {
	input := `# API Documentation

## Overview

This API provides **CRUD operations** for user management.

### Authentication

All endpoints require \`Authorization\` header:

\`\`\`http
GET /api/users HTTP/1.1
Authorization: Bearer token123
\`\`\`

### Endpoints

#### GET /users

Returns a list of users.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number |
| limit | integer | No | Items per page |

**Response:**

\`\`\`json
{
  "users": [
    {
      "id": 1,
      "name": "John Doe"
    }
  ]
}
\`\`\`

## Tasks

- [x] Implement authentication
- [ ] Add rate limiting
- [ ] Write tests

## Notes

> **Important:** Always use HTTPS in production.

---

For more information, visit [our website](https://example.com).`

	result := MdToTrac(input)
	
	// 基本的な変換が正しく行われているかチェック
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