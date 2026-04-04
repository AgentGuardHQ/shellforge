package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	data, _ := ioutil.ReadFile("internal/tools/tools.go")
	lines := strings.Split(string(data), "\n")
	for i := 200; i < 220; i++ {
		fmt.Printf("%3d: %q\n", i+1, lines[i])
	}
}