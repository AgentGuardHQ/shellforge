package main

import (
	"fmt"
	"os"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func main() {
	// Create test directory structure
	testDir := "/tmp/listfiles_integration_test"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir+"/subdir", 0755)
	os.WriteFile(testDir+"/file1.txt", []byte("test1"), 0644)
	os.WriteFile(testDir+"/file2.go", []byte("package main"), 0644)
	os.WriteFile(testDir+"/subdir/file3.md", []byte("# Test"), 0644)
	
	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	
	fmt.Println("Test 1: List files from parent directory")
	fmt.Println("========================================")
	os.Chdir("/tmp")
	result := tools.ExecuteDirect("list_files", map[string]string{
		"directory": testDir,
	}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	fmt.Println()
	
	fmt.Println("Test 2: List files from within the directory")
	fmt.Println("============================================")
	os.Chdir(testDir)
	result = tools.ExecuteDirect("list_files", map[string]string{
		"directory": ".",
	}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	fmt.Println()
	
	fmt.Println("Test 3: List files with .go extension filter")
	fmt.Println("============================================")
	os.Chdir("/tmp")
	result = tools.ExecuteDirect("list_files", map[string]string{
		"directory": testDir,
		"extension": ".go",
	}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	fmt.Println()
	
	fmt.Println("Test 4: List files from subdirectory")
	fmt.Println("====================================")
	result = tools.ExecuteDirect("list_files", map[string]string{
		"directory": testDir + "/subdir",
	}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Cleanup
	os.RemoveAll(testDir)
}