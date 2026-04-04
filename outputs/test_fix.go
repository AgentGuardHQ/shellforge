package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func simulateCurrentListFiles(dir string) []string {
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
	return files
}

func simulateFixedListFiles(dir string) []string {
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
		// Fixed implementation
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	return files
}

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_listfiles_fix"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir", "deep"), 0755)
	
	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "deep", "file3.txt"), []byte("test"), 0644)
	
	// Test 1: List files from current directory (should work the same)
	fmt.Println("=== Test 1: Listing current directory ===")
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	
	current := simulateCurrentListFiles(".")
	fixed := simulateFixedListFiles(".")
	fmt.Printf("Current implementation: %v\n", current)
	fmt.Printf("Fixed implementation: %v\n", fixed)
	
	// Test 2: List files from subdirectory
	fmt.Println("\n=== Test 2: Listing subdirectory ===")
	current = simulateCurrentListFiles("subdir")
	fixed = simulateFixedListFiles("subdir")
	fmt.Printf("Current implementation: %v\n", current)
	fmt.Printf("Fixed implementation: %v\n", fixed)
	
	// Test 3: List files from deep subdirectory
	fmt.Println("\n=== Test 3: Listing deep subdirectory ===")
	current = simulateCurrentListFiles("subdir/deep")
	fixed = simulateFixedListFiles("subdir/deep")
	fmt.Printf("Current implementation: %v\n", current)
	fmt.Printf("Fixed implementation: %v\n", fixed)
	
	os.Chdir(originalDir)
	
	// Test 4: List files from absolute path
	fmt.Println("\n=== Test 4: Listing with absolute path ===")
	os.Chdir("/") // Change to root to make it clear
	
	current = simulateCurrentListFiles(tmpDir)
	fixed = simulateFixedListFiles(tmpDir)
	fmt.Printf("Current implementation: %v\n", current)
	fmt.Printf("Fixed implementation: %v\n", fixed)
}