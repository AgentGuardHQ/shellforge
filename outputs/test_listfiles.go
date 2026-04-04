package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_listfiles"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	
	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.txt"), []byte("test"), 0644)
	
	// Change to tmpDir to simulate being in that directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)
	
	// Now test what listFiles would do
	testPath := "subdir"
	
	// Simulate what listFiles does currently
	var files []string
	filepath.WalkDir(testPath, func(path string, d os.DirEntry, err error) error {
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
		
		// This is the buggy line - uses "." (cwd) instead of testPath
		rel, _ := filepath.Rel(".", path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	
	fmt.Printf("Current directory: %s\n", tmpDir)
	fmt.Printf("Listing directory: %s\n", testPath)
	fmt.Printf("Files found: %v\n", files)
	
	// What should happen: files should be relative to testPath
	// So for file in subdir/file2.txt, it should return "file2.txt" not "subdir/file2.txt"
}