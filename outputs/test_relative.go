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
	// Create test structure in current directory
	os.RemoveAll("./test-local")
	os.MkdirAll("./test-local/x/y", 0755)
	os.WriteFile("./test-local/file1.txt", []byte("test"), 0644)
	os.WriteFile("./test-local/x/file2.txt", []byte("test"), 0644)
	os.WriteFile("./test-local/x/y/file3.txt", []byte("test"), 0644)
	
	fmt.Println("Testing from current directory")
	fmt.Println("Listing directory: ./test-local")
	fmt.Println("\n=== CURRENT IMPLEMENTATION ===")
	result1 := listFilesCurrent(map[string]string{"directory": "./test-local"})
	fmt.Println(result1)
	
	fmt.Println("\n=== FIXED IMPLEMENTATION ===")
	result2 := listFilesFixed(map[string]string{"directory": "./test-local"})
	fmt.Println(result2)
	
	fmt.Println("\n=== WHAT WE EXPECT ===")
	fmt.Println("Paths should be relative to ./test-local")
	fmt.Println("x/")
	fmt.Println("x/y/")
	fmt.Println("x/y/file3.txt")
	fmt.Println("x/file2.txt")
	fmt.Println("file1.txt")
	
	// Clean up
	os.RemoveAll("./test-local")
}