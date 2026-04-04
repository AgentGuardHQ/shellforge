package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	data, _ := ioutil.ReadFile("internal/tools/tools.go")
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.Contains(line, "rel, _ := filepath.Rel") {
			fmt.Printf("Line %d: %q\n", i+1, line)
			// Show previous line for context
			if i > 0 {
				fmt.Printf("Prev line: %q\n", lines[i-1])
			}
		}
	}
}