package main

import (
	"fmt"
	"os"

	"github.com/mi8bi/md2trac/internal/convert"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: md2trac input.md output.wiki")
		return
	}
	infile := os.Args[1]
	outfile := os.Args[2]

	data, err := os.ReadFile(infile)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	result := convert.MdToTrac(string(data))

	err = os.WriteFile(outfile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing output:", err)
	}
}
