package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_listfiles"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	
	// Create some test files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("test2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file3.txt"), []byte("test3"), 0644)
	
	// Change to a different directory
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	// Test the current implementation
	fmt.Println("Testing from /tmp directory, listing", tmpDir)
	
	// Simulate the bug: filepath.Rel(".", path) when we're in /tmp
	testPath := filepath.Join(tmpDir, "file1.txt")
	rel, _ := filepath.Rel(".", testPath)
	fmt.Printf("Current bug: filepath.Rel(\".\", %q) = %q\n", testPath, rel)
	
	// What it should be: relative to the directory being listed
	relCorrect, _ := filepath.Rel(tmpDir, testPath)
	fmt.Printf("Correct: filepath.Rel(%q, %q) = %q\n", tmpDir, testPath, relCorrect)
	
	// Test with subdirectory
	testPath2 := filepath.Join(tmpDir, "subdir", "file3.txt")
	rel2, _ := filepath.Rel(".", testPath2)
	fmt.Printf("Current bug: filepath.Rel(\".\", %q) = %q\n", testPath2, rel2)
	
	relCorrect2, _ := filepath.Rel(tmpDir, testPath2)
	fmt.Printf("Correct: filepath.Rel(%q, %q) = %q\n", tmpDir, testPath2, relCorrect2)
}