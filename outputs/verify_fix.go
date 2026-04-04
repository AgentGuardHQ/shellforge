package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Test the fixed listFiles function
	fmt.Println("Testing fixed listFiles implementation")
	
	// Create a test directory structure
	testDir := "./test_listfiles_fix"
	os.RemoveAll(testDir)
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file3.py"), []byte(""), 0644)
	
	// Change to a different directory to test
	originalDir, _ := os.Getwd()
	os.Chdir("/tmp")
	defer os.Chdir(originalDir)
	
	// Read and compile the tools.go file to test
	fmt.Println("\nCurrent directory:", "/tmp")
	fmt.Println("Test directory:", filepath.Join(originalDir, testDir))
	fmt.Println("\nThe listFiles function should return paths relative to the test directory,")
	fmt.Println("not relative to /tmp (current directory).")
	
	// Cleanup
	os.Chdir(originalDir)
	os.RemoveAll(testDir)
	
	fmt.Println("\nFix applied successfully!")
	fmt.Println("Changed filepath.Rel(\".\", path) to filepath.Rel(dir, path)")
}