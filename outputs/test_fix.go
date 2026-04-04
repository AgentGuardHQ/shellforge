package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFilesCurrent(dir string, ext string) ([]string, error) {
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
		// CURRENT BUG: uses "." (cwd) not dir
		rel, _ := filepath.Rel(".", path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func listFilesFixed(dir string, ext string) ([]string, error) {
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
		// FIXED: use dir as base
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
	// Create test directory structure
	cwd, _ := os.Getwd()
	testDir := filepath.Join(cwd, "test_fix_dir")
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.py"), []byte(""), 0644)
	
	fmt.Println("Testing listFiles fix")
	fmt.Println("Test directory:", testDir)
	fmt.Println("Current directory:", cwd)
	
	// Test current implementation
	fmt.Println("\n=== CURRENT (BUGGY) IMPLEMENTATION ===")
	files, err := listFilesCurrent(testDir, "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, f := range files {
			fmt.Println("  ", f)
		}
	}
	
	// Test fixed implementation
	fmt.Println("\n=== FIXED IMPLEMENTATION ===")
	files, err = listFilesFixed(testDir, "")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		for _, f := range files {
			fmt.Println("  ", f)
		}
	}
	
	// Cleanup
	os.RemoveAll(testDir)
}