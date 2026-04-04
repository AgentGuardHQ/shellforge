package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func main() {
	// Create test directory structure
	testDir := "/tmp/test_listfiles_ext"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	
	// Create test files with different extensions
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file3.md"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file4.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file5.go"), []byte("test"), 0644)
	
	// Change to a different directory
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	// Test listFiles with .go extension filter
	fmt.Println("Testing listFiles with .go extension filter:")
	
	result := tools.ExecuteDirect("list_files", map[string]string{
		"directory": testDir,
		"extension": ".go",
	}, 10)
	
	if !result.Success {
		fmt.Printf("Error: %s\n", result.Error)
		return
	}
	
	fmt.Println("Output:")
	fmt.Println(result.Output)
	
	// Check that only .go files are listed
	expectedGoFiles := []string{"file2.go", "subdir/file5.go"}
	output := result.Output
	allCorrect := true
	
	lines := splitLines(output)
	for _, line := range lines {
		if line == "./" || line == "subdir/" {
			// directories are OK
			continue
		}
		found := false
		for _, expected := range expectedGoFiles {
			if line == expected {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Unexpected file: %q\n", line)
			allCorrect = false
		}
	}
	
	if allCorrect {
		fmt.Println("\n✓ Only .go files are listed correctly")
	}
	
	// Clean up
	os.RemoveAll(testDir)
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