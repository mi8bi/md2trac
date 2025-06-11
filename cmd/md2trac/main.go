package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mi8bi/md2trac/internal/convert"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: md2trac input.md [output.wiki]")
		fmt.Println("  If output file is not specified, it will be input filename with .wiki extension")
		return
	}

	infile := os.Args[1]
	var outfile string

	if len(os.Args) >= 3 {
		outfile = os.Args[2]
	} else {
		// 出力ファイル名が指定されていない場合、入力ファイル名から生成
		ext := filepath.Ext(infile)
		base := strings.TrimSuffix(infile, ext)
		outfile = base + ".wiki"
	}

	// 入力ファイルの存在確認
	if _, err := os.Stat(infile); os.IsNotExist(err) {
fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", infile)
os.Exit(1)
		return
	}

	// 入力ファイルの読み込み
	data, err := os.ReadFile(infile)
	if err != nil {
		fmt.Printf("Error reading input file '%s': %v\n", infile, err)
		return
	}

	// Markdown → Trac Wiki変換
	fmt.Printf("Converting '%s' to '%s'...\n", infile, outfile)
	result := convert.MdToTrac(string(data))

	// 出力ディレクトリの作成（必要に応じて）
	outDir := filepath.Dir(outfile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fmt.Printf("Error creating output directory '%s': %v\n", outDir, err)
			return
		}
	}

	// 出力ファイルの書き込み
	err = os.WriteFile(outfile, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Error writing output file '%s': %v\n", outfile, err)
		return
	}

	fmt.Println("Conversion completed successfully!")

	// 統計情報の表示
	inputLines := strings.Count(string(data), "\n") + 1
	outputLines := strings.Count(result, "\n") + 1
	fmt.Printf("Input: %d lines, Output: %d lines\n", inputLines, outputLines)
}
