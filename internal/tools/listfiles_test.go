package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFiles_RelativePaths(t *testing.T) {
	// Create a test directory structure
	tmpDir := t.TempDir()
	
	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	
	// Create files
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("test"), 0644)
	
	// Change to a different directory to test the issue
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to parent directory
	parentDir := filepath.Dir(tmpDir)
	os.Chdir(parentDir)
	
	// Test listFiles from a different directory
	r := listFiles(map[string]string{
		"directory": tmpDir,
	}, 0)
	
	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	
	// Check that paths are relative to the listed directory, not cwd
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	
	expected := []string{"file1.txt", "subdir/", "subdir/file2.txt"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(lines), lines)
	}
	
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}

func TestListFiles_ExtensionFilter(t *testing.T) {
	tmpDir := t.TempDir()
	
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.md"), []byte("test"), 0644)
	
	r := listFiles(map[string]string{
		"directory": tmpDir,
		"extension": ".go",
	}, 0)
	
	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 1 || lines[0] != "file2.go" {
		t.Fatalf("expected only file2.go, got: %v", lines)
	}
}

func TestListFiles_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	
	r := listFiles(map[string]string{
		"directory": tmpDir,
	}, 0)
	
	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	
	if r.Output != "(empty directory)" {
		t.Fatalf("expected '(empty directory)', got %q", r.Output)
	}
}

func TestListFiles_SkipsHiddenAndSpecialDirs(t *testing.T) {
	tmpDir := t.TempDir()
	
	os.WriteFile(filepath.Join(tmpDir, "normal.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden.txt"), []byte("test"), 0644)
	
	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("test"), 0644)
	
	nodeModulesDir := filepath.Join(tmpDir, "node_modules")
	os.MkdirAll(nodeModulesDir, 0755)
	os.WriteFile(filepath.Join(nodeModulesDir, "package.json"), []byte("test"), 0644)
	
	r := listFiles(map[string]string{
		"directory": tmpDir,
	}, 0)
	
	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	
	// Should only have normal.txt
	if len(lines) != 1 || lines[0] != "normal.txt" {
		t.Fatalf("expected only normal.txt, got: %v", lines)
	}
}