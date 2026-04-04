package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("../internal/tools/tools.go")
	if err != nil {
		panic(err)
	}
	
	content := string(data)
	lines := strings.Split(content, "\n")
	
	// Find the line with filepath.Rel
	for i, line := range lines {
		if strings.Contains(line, `filepath.Rel(".", path)`) {
			fmt.Printf("Line %d: %q\n", i+1, line)
			// Show context
			for j := i-2; j <= i+2; j++ {
				if j >= 0 && j < len(lines) {
					fmt.Printf("%4d: %q\n", j+1, lines[j])
				}
			}
			break
		}
	}
}