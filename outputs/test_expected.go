package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_expected"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir", "deep"), 0755)
	
	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "deep", "file3.txt"), []byte("test"), 0644)
	
	// What would we expect?
	// If we're in /tmp/test_expected and we list "subdir", we should get:
	// - file2.txt (file in subdir)
	// - deep/ (subdirectory)
	// - deep/file3.txt (file in deep subdirectory)
	// But NOT "subdir/" or "subdir/deep/" - those paths are relative to the listed directory
	
	fmt.Println("Expected behavior when listing 'subdir':")
	fmt.Println("- file2.txt")
	fmt.Println("- deep/")
	fmt.Println("- deep/file3.txt")
	fmt.Println("\nNOT:")
	fmt.Println("- subdir/")
	fmt.Println("- subdir/deep/")
	fmt.Println("- subdir/deep/file3.txt (this is relative to cwd, not to subdir)")
	
	// Let's also check what happens with the extension filter
	fmt.Println("\n--- Testing actual listFiles function ---")
	
	// First, let me check if there are any other places in the code that might give us hints
	// about the expected behavior
}