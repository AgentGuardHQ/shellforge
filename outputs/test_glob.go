package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func testGlob() {
	// Create test structure
	os.RemoveAll("/tmp/test-glob")
	os.MkdirAll("/tmp/test-glob/a/b", 0755)
	os.WriteFile("/tmp/test-glob/file1.txt", []byte("test"), 0644)
	os.WriteFile("/tmp/test-glob/a/file2.txt", []byte("test"), 0644)
	os.WriteFile("/tmp/test-glob/a/b/file3.txt", []byte("test"), 0644)
	
	// Change to home directory
	os.Chdir(os.Getenv("HOME"))
	
	fmt.Println("Testing glob from:", os.Getenv("HOME"))
	fmt.Println("Directory: /tmp/test-glob")
	fmt.Println("Pattern: **/*.txt")
	
	dir := "/tmp/test-glob"
	pattern := "**/*.txt"
	
	var matches []string
	if strings.Contains(pattern, "**") {
		suffix := strings.TrimPrefix(pattern, "**/")
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			name := info.Name()
			if info.IsDir() && (name == ".git" || name == "node_modules") {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				matched, _ := filepath.Match(suffix, name)
				if matched {
					matches = append(matches, path)
				}
			}
			if len(matches) > 200 {
				return fmt.Errorf("limit reached")
			}
			return nil
		})
	}
	
	fmt.Println("\nGlob results (full paths):")
	for _, m := range matches {
		fmt.Println(m)
	}
	
	// Make relative to dir
	fmt.Println("\nGlob results (relative to dir):")
	for _, m := range matches {
		rel, _ := filepath.Rel(dir, m)
		fmt.Println(rel)
	}
}

func main() {
	testGlob()
}