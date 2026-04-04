package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Create test directory
	dir := "/tmp/debug_test"
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte(""), 0644)
	
	// Change to the directory
	originalWd, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalWd)
	
	// Walk the directory
	fmt.Println("Walking directory '.'")
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		fmt.Printf("  path=%q, name=%q, isDir=%v\n", path, name, d.IsDir())
		
		// Check if it would be filtered
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			fmt.Printf("    -> Would be filtered (starts with '.' or is special dir)\n")
			if d.IsDir() {
				fmt.Printf("    -> Would skip directory\n")
			}
		}
		return nil
	})
}