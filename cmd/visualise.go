package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type token struct {
	Ttype string `json:"type"`
	Value string
}

func longestTokenName(tokens *[]token) int {
	largest := 0
	for _, t := range *tokens {
		var length int = len(t.Ttype)
		if length > largest {
			largest = length
		}
	}
	return largest
}

func visualiseFile1(tokens *[]token) {
	tokenWidth := longestTokenName(tokens)
	formatGlobal := fmt.Sprintf("%%-%ds > %%s%%s", tokenWidth)
	formatContinuation := fmt.Sprintf("\n%%-%ds | ", tokenWidth)
	continuation := fmt.Sprintf(formatContinuation, " ")

	currentIndentation := 0
	for _, t := range *tokens {
		// print first line of this token
		lines := strings.Split(t.Value, "\n")
		indentStr := strings.Repeat(" ", currentIndentation)
		fmt.Printf(formatGlobal, t.Ttype, indentStr, lines[0])
		currentIndentation += len(lines[0])

		for _, line := range lines[1:] {
			fmt.Printf("%s%s", continuation, line)
			currentIndentation = len(line)
		}
		if strings.HasSuffix(t.Value, "\n") {
			currentIndentation = 0
		}

		fmt.Print("\n")
	}
}

func visualiseFile2(tokens *[]token) {
	currentLine := ""
	for _, t := range *tokens {
		// print first line of this token
		fmt.Printf("%s%s", currentLine, t.Value)
		lineCount := strings.Count(t.Value, "\n")
		newLinePos := strings.LastIndex(t.Value, "\n")
		

		if lineCount > 0 {
			fmt.Printf("\t\t <- (%d lines until here) %s\n", lineCount, t.Ttype)
			currentLine = t.Value[newLinePos+1:]
		} else {
			fmt.Printf("\t\t <- %s\n", t.Ttype)
			currentLine += t.Value
		}
	}
}

func VisualiseFile(reader io.Reader, mode string) {
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	var tokens []token

	err = json.Unmarshal(fileBytes, &tokens)
	if err != nil {
		log.Fatal(err)
	}
	if mode == "v1" {
		visualiseFile1(&tokens)
	} else if mode == "v2" {
		visualiseFile2(&tokens)
	} else {
		fmt.Printf("unknown format type: %s", mode)
	}
}

func VisualiseFileWithName(name string, mode string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	VisualiseFile(file, mode)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Please provide the format desired for printing (v1 or v2)\n")
		return
	}
	VisualiseFileWithName(os.Args[1], os.Args[2])
}
