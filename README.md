# md2trac

A command-line tool to convert Markdown files to Trac wiki format.

[![Build Status](https://github.com/mi8bi/md2trac/actions/workflows/test.yml/badge.svg)](https://github.com/mi8bi/md2trac/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mi8bi/md2trac)](https://goreportcard.com/report/github.com/mi8bi/md2trac)
[![codecov](https://codecov.io/gh/mi8bi/md2trac/branch/main/graph/badge.svg)](https://codecov.io/gh/mi8bi/md2trac)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/mi8bi/md2trac)](https://github.com/mi8bi/md2trac/releases/latest)

## Features

- **Headers**: Convert `# Header` to `= Header =`
- **Text Formatting**: Bold (`**bold**` → `'''bold'''`), Italic (`*italic*` → `''italic''`), Strikethrough (`~~text~~`)
- **Links**: `[text](url)` → `[url text]`
- **Images**: `![alt](image.png)` → `[[Image(image.png, alt)]]`
- **Code Blocks**: Support for syntax highlighting with language detection
- **Tables**: Convert Markdown tables to Trac table format
- **Lists**: Ordered and unordered lists with nesting support
- **Checkboxes**: Task lists (`- [x] Done` → ` * [X] Done`)
- **Blockquotes**: Quote formatting
- **Horizontal Rules**: `---` → `----`
- **Footnotes**: Reference-style footnotes

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/mi8bi/md2trac/releases/latest):

- **Linux**: `md2trac_v*.*.*.linux_amd64.tar.gz`
- **macOS**: `md2trac_v*.*.*.darwin_amd64.tar.gz` 
- **Windows**: `md2trac_v*.*.*.windows_amd64.zip`

Extract the binary and add it to your PATH.

### From Source

```bash
git clone https://github.com/mi8bi/md2trac.git
cd md2trac
go build -o md2trac ./cmd/md2trac
```

## Usage

```bash
# Convert input.md to input.wiki
md2trac input.md

# Convert input.md to custom output file
md2trac input.md output.wiki

# Convert documentation.md to trac format
md2trac documentation.md docs/trac-format.wiki
```

### Command Line Options

```
Usage: md2trac input.md [output.wiki]
  If output file is not specified, it will be input filename with .wiki extension
```

## Conversion Examples

### Headers
```markdown
# Main Title
## Section Title  
### Subsection
```
↓
```wiki
= Main Title =
== Section Title ==
=== Subsection ===
```

### Text Formatting
```markdown
**Bold text**
*Italic text*
***Bold and italic***
~~Strikethrough~~
`inline code`
```
↓
```wiki
'''Bold text'''
''Italic text''
'''''Bold and italic'''''
~~Strikethrough~~
`inline code`
```

### Code Blocks
```markdown
```python
def hello_world():
    print("Hello, World!")
```
↓
```wiki
{{{
#!python
def hello_world():
    print("Hello, World!")
}}}
```

### Tables
```markdown
| Name | Age | City |
|------|-----|------|
| John | 30  | NYC  |
| Jane | 25  | LA   |
```
↓
```wiki
|| Name || Age || City ||
|| John || 30 || NYC ||
|| Jane || 25 || LA ||
```

### Lists
```markdown
1. First item
2. Second item
   - Nested item
   - Another nested item
3. Third item

- [x] Completed task
- [ ] Pending task
```
↓
```wiki
 1. First item
 1. Second item
   * Nested item
   * Another nested item
 1. Third item

 * [X] Completed task
 * [ ] Pending task
```

## Supported Languages

Code blocks support syntax highlighting for various languages:

- `python`, `javascript`, `java`, `c`, `cpp`, `go`
- `json` (converted to JavaScript highlighting)
- `http` (converted to plain text)
- `html`, `css`, `xml`
- And many more...

## Development

### Prerequisites

- Go 1.24.4 or later

### Building

```bash
# Clone the repository
git clone https://github.com/mi8bi/md2trac.git
cd md2trac

# Install dependencies
go mod tidy

# Build the binary
go build -o md2trac ./cmd/md2trac

# Run tests
go test ./...
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/convert -v
```

### Project Structure

```
md2trac/
├── cmd/md2trac/           # Main application entry point
│   └── main.go
├── internal/convert/       # Core conversion logic
│   ├── convert.go
│   └── convert_test.go
├── .github/               # GitHub Actions workflows
│   ├── workflows/
│   └── dependabot.yml
├── .goreleaser.*.yaml     # Release configuration
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Adding New Conversions

To add support for new Markdown elements:

1. Add the conversion logic to `internal/convert/convert.go`
2. Add corresponding tests to `internal/convert/convert_test.go`
3. Update this README with examples

## Known Limitations

- Nested blockquotes are simplified to single-level quotes
- Complex table formatting (alignment, colspan) is not supported
- Some advanced Markdown features may not have direct Trac equivalents

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [Releases](https://github.com/mi8bi/md2trac/releases) for version history and changes.

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/mi8bi/md2trac/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/mi8bi/md2trac/issues)
- 📖 **Documentation**: This README and inline code comments

## Related Projects

- [Trac](https://trac.edgewall.org/) - The project management and bug tracking system
- [CommonMark](https://commonmark.org/) - Markdown specification reference