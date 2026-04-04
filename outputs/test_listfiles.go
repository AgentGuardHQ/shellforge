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
	return strings.Join(files, "\n")
}

func main() {
	// Change to /tmp directory
	os.Chdir("/tmp")
	
	// Test listing files in testdir from /tmp
	result := listFiles(map[string]string{"directory": "testdir"})
	fmt.Println("Current implementation (relative to cwd /tmp):")
	fmt.Println(result)
	
	// Now let's see what we expect
	fmt.Println("\nExpected (relative to testdir):")
	fmt.Println("testdir/")
	fmt.Println("testdir/file1.txt")
	fmt.Println("testdir/file2.go")
	fmt.Println("testdir/subdir/")
	fmt.Println("testdir/subdir/file3.txt")
}