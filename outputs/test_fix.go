package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFilesCurrent(params map[string]string) string {
	dir := params["directory"]
	if dir == "" {
		dir = params["path"]
	}
	if dir == "" {
		dir = params["dir"]
	}
	if dir == "" {
		dir = "."
	}
	ext := params["extension"]
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
		if ext != "" && filepath.Ext(name) != ext {
			return nil
		}
		rel, _ := filepath.Rel(".", path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	return strings.Join(files, "\n")
}

func listFilesFixed(params map[string]string) string {
	dir := params["directory"]
	if dir == "" {
		dir = params["path"]
	}
	if dir == "" {
		dir = params["dir"]
	}
	if dir == "" {
		dir = "."
	}
	ext := params["extension"]
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
		if ext != "" && filepath.Ext(name) != ext {
			return nil
		}
		// FIXED: Make path relative to dir, not to cwd
		rel, _ := filepath.Rel(dir, path)
		if rel == "." {
			// This is the directory itself
			return nil
		}
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	return strings.Join(files, "\n")
}

func main() {
	// Create test structure
	os.RemoveAll("/tmp/test-fix")
	os.MkdirAll("/tmp/test-fix/a/b", 0755)
	os.WriteFile("/tmp/test-fix/file1.txt", []byte("test"), 0644)
	os.WriteFile("/tmp/test-fix/a/file2.txt", []byte("test"), 0644)
	os.WriteFile("/tmp/test-fix/a/b/file3.txt", []byte("test"), 0644)
	
	// Change to home directory
	os.Chdir(os.Getenv("HOME"))
	
	fmt.Println("Testing from:", os.Getenv("HOME"))
	fmt.Println("Listing directory: /tmp/test-fix")
	fmt.Println("\n=== CURRENT IMPLEMENTATION ===")
	result1 := listFilesCurrent(map[string]string{"directory": "/tmp/test-fix"})
	fmt.Println(result1)
	
	fmt.Println("\n=== FIXED IMPLEMENTATION ===")
	result2 := listFilesFixed(map[string]string{"directory": "/tmp/test-fix"})
	fmt.Println(result2)
	
	fmt.Println("\n=== WHAT WE EXPECT ===")
	fmt.Println("Paths should be relative to /tmp/test-fix, not to home directory")
	fmt.Println("a/")
	fmt.Println("a/b/")
	fmt.Println("a/b/file3.txt")
	fmt.Println("a/file2.txt")
	fmt.Println("file1.txt")
}