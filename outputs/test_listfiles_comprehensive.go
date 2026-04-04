package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listFiles(params map[string]string) string {
	dir := params["directory"]
	if dir == "" {
		dir = params["path"]
	}
	if dir == "" {
		dir = params["dir"]
	}
	if dir == "" {
		dir = "."
	}
	ext := params["extension"]
	var files []string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if len(files) > 200 {
			return fmt.Errorf("limit reached")
		}
		if ext != "" && filepath.Ext(name) != ext {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			files = append(files, rel+"/")
		} else {
			files = append(files, rel)
		}
		return nil
	})
	if len(files) == 0 {
		return "(empty directory)"
	}
	return strings.Join(files, "\n")
}

func main() {
	// Create a fresh test directory
	testDir := "/tmp/test_listfiles_fresh"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir+"/subdir", 0755)
	os.WriteFile(testDir+"/file1.txt", []byte("test"), 0644)
	os.WriteFile(testDir+"/file2.go", []byte("package main"), 0644)
	os.WriteFile(testDir+"/subdir/file3.md", []byte("# Test"), 0644)
	
	// Save original cwd
	originalCwd, _ := os.Getwd()
	
	fmt.Println("Test 1: From parent directory, list subdirectory")
	fmt.Println("================================================")
	os.Chdir("/tmp")
	result := listFiles(map[string]string{"directory": testDir})
	fmt.Println("Result:")
	fmt.Println(result)
	fmt.Println()
	
	fmt.Println("Test 2: From within the directory, list current directory")
	fmt.Println("========================================================")
	os.Chdir(testDir)
	result = listFiles(map[string]string{"directory": "."})
	fmt.Println("Result:")
	fmt.Println(result)
	fmt.Println()
	
	fmt.Println("Test 3: From parent, list subdirectory with extension filter")
	fmt.Println("===========================================================")
	os.Chdir("/tmp")
	result = listFiles(map[string]string{"directory": testDir, "extension": ".go"})
	fmt.Println("Result:")
	fmt.Println(result)
	fmt.Println()
	
	fmt.Println("Test 4: List a nested subdirectory")
	fmt.Println("==================================")
	result = listFiles(map[string]string{"directory": testDir + "/subdir"})
	fmt.Println("Result:")
	fmt.Println(result)
	fmt.Println()
	
	// Restore original cwd
	os.Chdir(originalCwd)
}