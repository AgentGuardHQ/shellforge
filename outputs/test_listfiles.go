package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Create a test directory structure
	testDir := "/tmp/test_listfiles"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	
	// Create some test files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.txt"), []byte("test"), 0644)
	
	// Change to a different directory
	os.Chdir("/tmp")
	
	// Simulate the current buggy behavior
	dir := testDir
	path := filepath.Join(testDir, "file1.txt")
	rel, _ := filepath.Rel(".", path)
	fmt.Printf("Current buggy behavior (relative to cwd): %s -> %s\n", path, rel)
	
	// What it should be
	relCorrect, _ := filepath.Rel(dir, path)
	fmt.Printf("Correct behavior (relative to dir): %s -> %s\n", path, relCorrect)
	
	// Test with subdirectory
	path2 := filepath.Join(testDir, "subdir", "file3.txt")
	rel2, _ := filepath.Rel(".", path2)
	fmt.Printf("Current buggy behavior (relative to cwd): %s -> %s\n", path2, rel2)
	
	rel2Correct, _ := filepath.Rel(dir, path2)
	fmt.Printf("Correct behavior (relative to dir): %s -> %s\n", path2, rel2Correct)
}