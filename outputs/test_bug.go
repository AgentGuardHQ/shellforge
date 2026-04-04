package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Create test directory structure
	testDir := "/tmp/test_listfiles_bug"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	
	// Create test files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.txt"), []byte("test"), 0644)
	
	// Change to a different directory
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	// Test the current implementation
	fmt.Println("Testing current implementation:")
	testListFiles(testDir)
	
	// Clean up
	os.RemoveAll(testDir)
}

func testListFiles(dir string) {
	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
		if len(files) > 200 {
			return fmt.Errorf("limit reached")
		}
		// Current buggy implementation
		rel, _ := filepath.Rel(".", path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	
	fmt.Println("Files found (relative to cwd):")
	for _, f := range files {
		fmt.Printf("  %q\n", f)
	}
	
	// Now test the fixed version
	var filesFixed []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
		if len(filesFixed) > 200 {
			return fmt.Errorf("limit reached")
		}
		// Fixed implementation
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			filesFixed = append(filesFixed, rel+"/")
		} else {
			filesFixed = append(filesFixed, rel)
		}
		return nil
	})
	
	fmt.Println("\nFiles found (relative to dir - fixed):")
	for _, f := range filesFixed {
		fmt.Printf("  %q\n", f)
	}
}