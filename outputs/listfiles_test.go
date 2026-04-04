package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFiles_Basic(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "b.go"), []byte(""), 0o644)
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(sub, "c.md"), []byte(""), 0o644)

	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to a different directory to test relative paths
	os.Chdir("/tmp")

	r := listFiles(map[string]string{
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	// Check that paths are relative to 'dir', not to cwd
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	expected := []string{"./", "a.txt", "b.go", "sub/", "sub/c.md"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %s", len(expected), len(lines), r.Output)
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line %d: got %q, want %q", i, line, expected[i])
		}
	}
}

func TestListFiles_FromWithinDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte(""), 0o644)

	// Save original working directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	// Change to the test directory
	os.Chdir(dir)

	r := listFiles(map[string]string{
		"directory": ".",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	// Should list files relative to current directory (.)
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 1 || lines[0] != "test.txt" {
		t.Fatalf("expected ['test.txt'], got %v", lines)
	}
}

func TestListFiles_ExtensionFilter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "b.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "c.go"), []byte(""), 0o644)

	r := listFiles(map[string]string{
		"directory": dir,
		"extension": ".go",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 .go files, got %d: %s", len(lines), r.Output)
	}
	for _, line := range lines {
		if !strings.HasSuffix(line, ".go") {
			t.Errorf("expected .go file, got %q", line)
		}
	}
}

func TestListFiles_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	r := listFiles(map[string]string{
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if r.Output != "(empty directory)" {
		t.Fatalf("expected '(empty directory)', got %q", r.Output)
	}
}

func TestListFiles_SkipsHiddenAndSpecialDirs(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "normal.txt"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, ".hidden"), []byte(""), 0o644)
	gitDir := filepath.Join(dir, ".git")
	os.MkdirAll(gitDir, 0o755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte(""), 0o644)
	nodeModules := filepath.Join(dir, "node_modules")
	os.MkdirAll(nodeModules, 0o755)
	os.WriteFile(filepath.Join(nodeModules, "package.json"), []byte(""), 0o644)

	r := listFiles(map[string]string{
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	// Should only show normal.txt
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 1 || lines[0] != "normal.txt" {
		t.Fatalf("expected only 'normal.txt', got %v", lines)
	}
}

func TestListFiles_SubdirectoryListing(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(sub, "another.md"), []byte(""), 0o644)

	r := listFiles(map[string]string{
		"directory": sub,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	// Should list files relative to subdir
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	expected := []string{"./", "file.txt", "another.md"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %s", len(expected), len(lines), r.Output)
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("line %d: got %q, want %q", i, line, expected[i])
		}
	}
}