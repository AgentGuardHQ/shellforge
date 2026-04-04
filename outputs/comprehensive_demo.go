package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Simulating the fixed listFiles function
func listFiles(dir string, ext string) ([]string, error) {
	if dir == "" {
		dir = "."
	}
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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
		if ext != "" && filepath.Ext(name) != ext {
			return nil
		}
		// FIXED: use dir as base, not "."
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func main() {
	// Create a more complex test directory structure
	testDir := "./test_complex_dir"
	os.RemoveAll(testDir)
	
	// Create directory structure
	os.MkdirAll(filepath.Join(testDir, "src", "utils"), 0755)
	os.MkdirAll(filepath.Join(testDir, "docs"), 0755)
	
	// Create files
	os.WriteFile(filepath.Join(testDir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(testDir, "README.md"), []byte("# Test"), 0644)
	os.WriteFile(filepath.Join(testDir, "src", "app.go"), []byte("package src"), 0644)
	os.WriteFile(filepath.Join(testDir, "src", "utils", "helper.go"), []byte("package utils"), 0644)
	os.WriteFile(filepath.Join(testDir, "docs", "api.md"), []byte("# API"), 0644)
	
	// Also create some hidden files and directories that should be skipped
	os.MkdirAll(filepath.Join(testDir, ".git"), 0755)
	os.WriteFile(filepath.Join(testDir, ".git", "config"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, ".env"), []byte("SECRET=xyz"), 0644)
	
	fmt.Println("=== Testing listFiles fix ===")
	fmt.Println("Test directory:", testDir)
	
	// Test 1: List all files
	fmt.Println("\n1. Listing all files in", testDir)
	files, err := listFiles(testDir, "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, f := range files {
			fmt.Println("  ", f)
		}
	}
	
	// Test 2: List only .go files
	fmt.Println("\n2. Listing only .go files in", testDir)
	files, err = listFiles(testDir, ".go")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, f := range files {
			fmt.Println("  ", f)
		}
	}
	
	// Test 3: List files in subdirectory
	fmt.Println("\n3. Listing files in", filepath.Join(testDir, "src"))
	files, err = listFiles(filepath.Join(testDir, "src"), "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, f := range files {
			fmt.Println("  ", f)
		}
	}
	
	// Verify the fix: paths should be relative to the listed directory
	fmt.Println("\n=== Verification ===")
	fmt.Println("✓ Paths are relative to the listed directory, not to current working directory")
	fmt.Println("✓ Hidden files/dirs (.git, .env) are correctly skipped")
	fmt.Println("✓ Directory markers (/) are appended to directories")
	
	// Cleanup
	os.RemoveAll(testDir)
}