package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFiles(params map[string]string) string {
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
	if len(files) == 0 {
		return "(empty directory)"
	}
	return strings.Join(files, "\n")
}

func main() {
	// Test current behavior
	fmt.Println("Current behavior test:")
	fmt.Println("======================")
	
	// Save original cwd
	originalCwd, _ := os.Getwd()
	
	// Create test directory structure
	testDir := "/tmp/test_listfiles"
	os.Chdir(testDir)
	
	// Test listing from current directory
	result := listFiles(map[string]string{"directory": "."})
	fmt.Println("listFiles from '.':")
	fmt.Println(result)
	fmt.Println()
	
	// Test listing from subdirectory
	result = listFiles(map[string]string{"directory": "subdir"})
	fmt.Println("listFiles from 'subdir':")
	fmt.Println(result)
	fmt.Println()
	
	// Change to a different directory and test
	os.Chdir("/tmp")
	result = listFiles(map[string]string{"directory": testDir})
	fmt.Println("listFiles from '/tmp' with directory='", testDir, "':")
	fmt.Println(result)
	fmt.Println()
	
	// Restore original cwd
	os.Chdir(originalCwd)
}