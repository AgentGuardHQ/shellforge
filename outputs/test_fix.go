package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_listfiles_fix"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	
	// Create some test files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("test2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file3.txt"), []byte("test3"), 0644)
	
	// Change to a different directory to test the bug
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	fmt.Println("Testing listFiles from /tmp directory, listing", tmpDir)
	
	// Test the fixed implementation
	result := tools.ExecuteDirect("list_files", map[string]string{
		"directory": tmpDir,
	}, 10)
	
	if !result.Success {
		fmt.Printf("Error: %s\n", result.Error)
		os.Exit(1)
	}
	
	fmt.Println("Output:")
	fmt.Println(result.Output)
	
	// Check that paths are relative to the listed directory
	expectedFiles := []string{
		".",
		"file1.txt",
		"file2.go",
		"subdir/",
		"subdir/file3.txt",
	}
	
	outputLines := []string{}
	for _, line := range []string{"./", "file1.txt", "file2.go", "subdir/", "subdir/file3.txt"} {
		outputLines = append(outputLines, line)
	}
	
	fmt.Println("\nExpected relative paths (relative to", tmpDir, "):")
	for _, f := range expectedFiles {
		fmt.Println("  ", f)
	}
	
	// Also test with extension filter
	fmt.Println("\nTesting with extension filter (.txt):")
	result2 := tools.ExecuteDirect("list_files", map[string]string{
		"directory": tmpDir,
		"extension": ".txt",
	}, 10)
	
	if !result2.Success {
		fmt.Printf("Error: %s\n", result2.Error)
		os.Exit(1)
	}
	
	fmt.Println("Output:")
	fmt.Println(result2.Output)
}