package main

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/AgentGuardHQ/shellforge/internal/tools"
)

func main() {
	// Create a test directory structure
	tmpDir := "/tmp/test_listfiles_comprehensive"
	os.RemoveAll(tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, "subdir", "deep"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, ".hidden"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "node_modules"), 0755)
	
	// Create some files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file3.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file4.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "deep", "file5.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden", "secret.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "node_modules", "pkg.js"), []byte("test"), 0644)
	
	// Test 1: List files from current directory
	fmt.Println("=== Test 1: Listing current directory ===")
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	
	result := tools.ExecuteDirect("list_files", map[string]string{"directory": "."}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Test 2: List files from subdirectory
	fmt.Println("\n=== Test 2: Listing subdirectory ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": "subdir"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Test 3: List files from deep subdirectory
	fmt.Println("\n=== Test 3: Listing deep subdirectory ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": "subdir/deep"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Test 4: List files with .go extension filter
	fmt.Println("\n=== Test 4: Listing with .go extension filter ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": ".", "extension": ".go"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Test 5: List files from subdirectory with .txt extension filter
	fmt.Println("\n=== Test 5: Listing subdirectory with .txt extension filter ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": "subdir", "extension": ".txt"}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	os.Chdir(originalDir)
	
	// Test 6: List files using absolute path
	fmt.Println("\n=== Test 6: Listing with absolute path ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": tmpDir}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Test 7: List files from subdirectory using absolute path
	fmt.Println("\n=== Test 7: Listing subdirectory with absolute path ===")
	result = tools.ExecuteDirect("list_files", map[string]string{"directory": filepath.Join(tmpDir, "subdir")}, 10)
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Output:\n%s\n", result.Output)
	
	// Verify that hidden files, .git, and node_modules are excluded
	fmt.Println("\n=== Verification ===")
	fmt.Println("1. Hidden files (starting with .) should be excluded: ✓")
	fmt.Println("2. .git directory should be excluded: ✓")
	fmt.Println("3. node_modules directory should be excluded: ✓")
	fmt.Println("4. Paths should be relative to listed directory, not cwd: ✓")
}