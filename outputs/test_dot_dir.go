package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func testDotDirectory() {
	// Create a simple test
	testDir := "/tmp/test-dot-dir"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0755)
	os.WriteFile(filepath.Join(testDir, "testfile.txt"), []byte("test"), 0644)
	
	// Change to test directory
	originalDir, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalDir)
	
	// Test listing "." from within the directory
	fmt.Println("=== Test: Listing '.' from within directory ===")
	result := tools.ExecuteDirect("list_files", map[string]string{"directory": "."}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output:")
		fmt.Println(result.Output)
	}
	
	// The output should be just "testfile.txt" (not "./testfile.txt" or "testDir/testfile.txt")
	if result.Output == "testfile.txt" {
		fmt.Println("✓ Correct: Path is relative to '.'")
	} else {
		fmt.Printf("✗ Incorrect: Expected 'testfile.txt', got %q\n", result.Output)
	}
	
	os.RemoveAll(testDir)
}

func main() {
	testDotDirectory()
}