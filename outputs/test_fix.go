package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func main() {
	// Create test directory structure
	testDir := "/tmp/test_listfiles_fix"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	
	// Create test files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.txt"), []byte("test"), 0644)
	
	// Change to a different directory
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	// Test listFiles with the fixed implementation
	fmt.Println("Testing listFiles from /tmp on directory:", testDir)
	
	result := tools.ExecuteDirect("list_files", map[string]string{
		"directory": testDir,
	}, 10)
	
	if !result.Success {
		fmt.Printf("Error: %s\n", result.Error)
		return
	}
	
	fmt.Println("Output:")
	fmt.Println(result.Output)
	
	// Check that paths are relative to the directory, not cwd
	lines := []string{"./", "file1.txt", "file2.go", "subdir/", "subdir/file3.txt"}
	output := result.Output
	allFound := true
	for _, line := range lines {
		if !containsLine(output, line) {
			fmt.Printf("Missing expected line: %q\n", line)
			allFound = false
		}
	}
	
	if allFound {
		fmt.Println("\n✓ All paths are correctly relative to the directory parameter")
	} else {
		fmt.Println("\n✗ Some paths are missing or incorrect")
	}
	
	// Clean up
	os.RemoveAll(testDir)
}

func containsLine(output, line string) bool {
	lines := splitLines(output)
	for _, l := range lines {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}