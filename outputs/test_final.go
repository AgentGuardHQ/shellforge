package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func testListFiles() {
	// Create test structure
	testDir := "/tmp/test-listfiles-fix"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "a", "b"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "a", "file3.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "a", "b", "file4.txt"), []byte("test"), 0644)
	
	// Save current directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Test 1: Absolute path from different directory
	fmt.Println("=== Test 1: Absolute path from different directory ===")
	os.Chdir("/tmp")
	result := tools.ExecuteDirect("list_files", map[string]string{"directory": testDir}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output:")
		fmt.Println(result.Output)
	}
	
	// Check that paths are relative to testDir
	expected := []string{"a/", "a/b/", "a/b/file4.txt", "a/file3.txt", "file1.txt", "file2.go"}
	lines := strings.Split(strings.TrimSpace(result.Output), "\n")
	fmt.Printf("\nExpected %d files, got %d\n", len(expected), len(lines))
	
	// Test 2: Relative path from parent directory
	fmt.Println("\n=== Test 2: Relative path from parent directory ===")
	os.Chdir("/tmp")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": "test-listfiles-fix"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output:")
		fmt.Println(result.Output)
	}
	
	// Test 3: With extension filter
	fmt.Println("\n=== Test 3: With extension filter (.txt) ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": testDir, "extension": ".txt"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output:")
		fmt.Println(result.Output)
	}
	
	// Test 4: Empty directory
	fmt.Println("\n=== Test 4: Empty directory ===")
	emptyDir := "/tmp/empty-test-dir"
	os.RemoveAll(emptyDir)
	os.MkdirAll(emptyDir, 0755)
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": emptyDir}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output: %q\n", result.Output)
	
	// Clean up
	os.RemoveAll(testDir)
	os.RemoveAll(emptyDir)
}

func main() {
	testListFiles()
}