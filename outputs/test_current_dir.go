package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func testCurrentDirectory() {
	// Create test structure in current directory
	testDir := "./test-current-dir"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "sub"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "sub", "file2.txt"), []byte("test"), 0644)
	
	// Test listing current directory
	fmt.Println("=== Test: Listing current directory ===")
	result := tools.ExecuteDirect("list_files", map[string]string{"directory": "."}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output (first few lines):")
		allLines := strings.Split(result.Output, "\n")
		for i := 0; i < 5 && i < len(allLines); i++ {
			fmt.Println(allLines[i])
		}
	}
	
	// Test listing test directory from within it
	fmt.Println("\n=== Test: Listing test directory from within it ===")
	os.Chdir(testDir)
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": "."}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	if result.Success {
		fmt.Println("Output:")
		fmt.Println(result.Output)
	}
	
	// Clean up and go back
	os.Chdir("..")
	os.RemoveAll(testDir)
}

func main() {
	testCurrentDirectory()
}