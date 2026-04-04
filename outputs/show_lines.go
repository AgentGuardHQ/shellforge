package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	data, _ := ioutil.ReadFile("./internal/tools/tools.go")
	lines := strings.Split(string(data), "\n")
	
	// Get lines 210-220
	for i := 209; i < 220 && i < len(lines); i++ {
		fmt.Printf("%d: %q\n", i+1, lines[i])
	}
}