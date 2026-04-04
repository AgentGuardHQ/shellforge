package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Create test directory structure
	testDir := "./test_listfiles_dir"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.py"), []byte(""), 0644)
	
	// Save current directory
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	
	// Change to a different directory to test
	os.Chdir("/tmp")
	
	// Simulate what listFiles does currently
	var files []string
	filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		// This is the problematic line - uses "." (cwd) not testDir
		rel, _ := filepath.Rel(".", path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	
	fmt.Println("Current implementation (relative to cwd which is /tmp):")
	for _, f := range files {
		fmt.Println("  ", f)
	}
	
	// Now test what it should do
	files = []string{}
	filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		// This is what it should do - use testDir as base
		rel, _ := filepath.Rel(testDir, path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	
	fmt.Println("\nShould be (relative to testDir):")
	for _, f := range files {
		fmt.Println("  ", f)
	}
	
	// Cleanup
	os.Chdir(cwd)
	os.RemoveAll(testDir)
}