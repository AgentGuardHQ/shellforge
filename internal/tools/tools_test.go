package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── edit_file tests ──

func TestEditFile_BasicReplace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello world"), 0o644)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": "world",
		"new_text": "go",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "hello go" {
		t.Fatalf("expected 'hello go', got %q", string(data))
	}
}

func TestEditFile_NoMatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello world"), 0o644)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": "foobar",
		"new_text": "baz",
	}, 0)

	if r.Success {
		t.Fatal("expected failure for no match")
	}
	if r.Error != "no_match" {
		t.Fatalf("expected error 'no_match', got %q", r.Error)
	}
}

func TestEditFile_MultipleMatches(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("aaa bbb aaa"), 0o644)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": "aaa",
		"new_text": "ccc",
	}, 0)

	if r.Success {
		t.Fatal("expected failure for multiple matches")
	}
	if r.Error != "multiple_matches" {
		t.Fatalf("expected error 'multiple_matches', got %q", r.Error)
	}

	// Verify file is unchanged
	data, _ := os.ReadFile(path)
	if string(data) != "aaa bbb aaa" {
		t.Fatalf("file should be unchanged, got %q", string(data))
	}
}

func TestEditFile_FileNotFound(t *testing.T) {
	r := editFile(map[string]string{
		"path":     "/nonexistent/file.txt",
		"old_text": "foo",
		"new_text": "bar",
	}, 0)

	if r.Success {
		t.Fatal("expected failure for missing file")
	}
	if r.Error != "not_found" {
		t.Fatalf("expected error 'not_found', got %q", r.Error)
	}
}

func TestEditFile_PreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.sh")
	os.WriteFile(path, []byte("#!/bin/bash\necho hello"), 0o755)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": "echo hello",
		"new_text": "echo goodbye",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	info, _ := os.Stat(path)
	perm := info.Mode().Perm()
	if perm != 0o755 {
		t.Fatalf("expected permissions 0755, got %o", perm)
	}
}

func TestEditFile_MissingParams(t *testing.T) {
	r := editFile(map[string]string{
		"old_text": "foo",
		"new_text": "bar",
	}, 0)
	if r.Success {
		t.Fatal("expected failure for missing path")
	}

	r = editFile(map[string]string{
		"path":     "/tmp/test.txt",
		"new_text": "bar",
	}, 0)
	if r.Success {
		t.Fatal("expected failure for missing old_text")
	}
}

func TestEditFile_MultilineContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	content := "package main\n\nfunc old() {\n\treturn\n}\n"
	os.WriteFile(path, []byte(content), 0o644)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": "func old() {\n\treturn\n}",
		"new_text": "func newFunc() {\n\tfmt.Println(\"new\")\n}",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "func newFunc()") {
		t.Fatalf("expected new function name, got %q", string(data))
	}
}

func TestEditFile_EmptyNewText(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("keep this remove this keep that"), 0o644)

	r := editFile(map[string]string{
		"path":     path,
		"old_text": " remove this",
		"new_text": "",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "keep this keep that" {
		t.Fatalf("expected deletion, got %q", string(data))
	}
}

// ── glob tests ──

func TestGlob_SimplePattern(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "b.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte(""), 0o644)

	r := globFiles(map[string]string{
		"pattern":   "*.go",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	lines := strings.Split(r.Output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 matches, got %d: %s", len(lines), r.Output)
	}
}

func TestGlob_RecursivePattern(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(dir, "a.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(sub, "b.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(sub, "c.txt"), []byte(""), 0o644)

	r := globFiles(map[string]string{
		"pattern":   "**/*.go",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}

	lines := strings.Split(r.Output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 matches, got %d: %s", len(lines), r.Output)
	}
}

func TestGlob_NoMatches(t *testing.T) {
	dir := t.TempDir()

	r := globFiles(map[string]string{
		"pattern":   "*.xyz",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if r.Output != "No files matched" {
		t.Fatalf("expected 'No files matched', got %q", r.Output)
	}
}

func TestGlob_SkipsGitDir(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.MkdirAll(gitDir, 0o755)
	os.WriteFile(filepath.Join(gitDir, "config.go"), []byte(""), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(""), 0o644)

	r := globFiles(map[string]string{
		"pattern":   "**/*.go",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if strings.Contains(r.Output, ".git") {
		t.Fatalf("should not include .git files, got: %s", r.Output)
	}
	lines := strings.Split(r.Output, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 match, got %d: %s", len(lines), r.Output)
	}
}

func TestGlob_DefaultDirectory(t *testing.T) {
	r := globFiles(map[string]string{
		"pattern": "*.nonexistent_extension_xyz",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if r.Output != "No files matched" {
		t.Fatalf("expected 'No files matched', got %q", r.Output)
	}
}

// ── grep tests ──

func TestGrep_BasicMatch(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("package main\n\nfunc hello() {}\nfunc world() {}\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "func hello",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if !strings.Contains(r.Output, "func hello") {
		t.Fatalf("expected match for 'func hello', got: %s", r.Output)
	}
	if !strings.Contains(r.Output, ":3:") {
		t.Fatalf("expected line number 3, got: %s", r.Output)
	}
}

func TestGrep_RegexPattern(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("foo123\nbar456\nfoo789\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "foo\\d+",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 matches, got %d: %s", len(lines), r.Output)
	}
}

func TestGrep_InvalidRegex(t *testing.T) {
	r := grepFiles(map[string]string{
		"pattern":   "[invalid",
		"directory": ".",
	}, 0)

	if r.Success {
		t.Fatal("expected failure for invalid regex")
	}
	if r.Error != "bad_pattern" {
		t.Fatalf("expected error 'bad_pattern', got %q", r.Error)
	}
}

func TestGrep_NoMatches(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("hello world\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "nonexistent_string",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if r.Output != "No matches found" {
		t.Fatalf("expected 'No matches found', got %q", r.Output)
	}
}

func TestGrep_FileTypeFilter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("match this\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "test.py"), []byte("match this\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "match",
		"directory": dir,
		"file_type": "go",
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if !strings.Contains(r.Output, "test.go") {
		t.Fatalf("expected test.go match, got: %s", r.Output)
	}
	if strings.Contains(r.Output, "test.py") {
		t.Fatalf("should not include test.py, got: %s", r.Output)
	}
}

func TestGrep_SkipsGitAndNodeModules(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	nmDir := filepath.Join(dir, "node_modules")
	os.MkdirAll(gitDir, 0o755)
	os.MkdirAll(nmDir, 0o755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("match this\n"), 0o644)
	os.WriteFile(filepath.Join(nmDir, "pkg.js"), []byte("match this\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("match this\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "match",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if strings.Contains(r.Output, ".git") {
		t.Fatalf("should not include .git files, got: %s", r.Output)
	}
	if strings.Contains(r.Output, "node_modules") {
		t.Fatalf("should not include node_modules, got: %s", r.Output)
	}
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 match, got %d: %s", len(lines), r.Output)
	}
}

func TestGrep_RecursiveSearch(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub", "deep")
	os.MkdirAll(sub, 0o755)
	os.WriteFile(filepath.Join(dir, "a.go"), []byte("findme\n"), 0o644)
	os.WriteFile(filepath.Join(sub, "b.go"), []byte("findme too\n"), 0o644)

	r := grepFiles(map[string]string{
		"pattern":   "findme",
		"directory": dir,
	}, 0)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	lines := strings.Split(strings.TrimSpace(r.Output), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 matches, got %d: %s", len(lines), r.Output)
	}
}

// ── ExecuteDirect dispatch tests ──

func TestExecuteDirect_UnknownTool(t *testing.T) {
	r := ExecuteDirect("nonexistent_tool", map[string]string{}, 10)
	if r.Success {
		t.Fatal("expected failure for unknown tool")
	}
	if r.Error != "unknown_tool" {
		t.Fatalf("expected error 'unknown_tool', got %q", r.Error)
	}
}

func TestExecuteDirect_EditFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("old content"), 0o644)

	r := ExecuteDirect("edit_file", map[string]string{
		"path":     path,
		"old_text": "old",
		"new_text": "new",
	}, 10)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "new content" {
		t.Fatalf("expected 'new content', got %q", string(data))
	}
}

func TestExecuteDirect_Glob(t *testing.T) {
	r := ExecuteDirect("glob", map[string]string{
		"pattern": "*.nonexistent_xyz",
	}, 10)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
}

func TestExecuteDirect_Grep(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello world\n"), 0o644)

	r := ExecuteDirect("grep", map[string]string{
		"pattern":   "hello",
		"directory": dir,
	}, 10)

	if !r.Success {
		t.Fatalf("expected success, got error: %s", r.Error)
	}
	if !strings.Contains(r.Output, "hello world") {
		t.Fatalf("expected match, got: %s", r.Output)
	}
}

// ── list_files tests ──

func TestListFiles_RelativeToDirectory(t *testing.T) {
	// Create a test directory structure
	tmpDir := t.TempDir()
	
	// Create some files and subdirectories
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "file3.go"), []byte("test"), 0644)
	
	// Change to a different directory to test that paths are relative to dir, not cwd
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)
	
	// Change to /tmp (or another directory)
	os.Chdir("/tmp")
	
	// Test 1: List all files
	result := listFiles(map[string]string{"directory": tmpDir}, 0)
	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.Error)
	}
	
	output := result.Output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	// Check that we get the expected files
	expectedFiles := []string{"file1.txt", "subdir/", "subdir/file2.txt", "subdir/file3.go"}
	if len(lines) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d: %v", len(expectedFiles), len(lines), lines)
	}
	
	// Check that paths are relative to tmpDir, not cwd
	for _, line := range lines {
		if strings.Contains(line, "/tmp") {
			t.Errorf("Path should not contain absolute path /tmp: %s", line)
		}
		if strings.HasPrefix(line, "/") {
			t.Errorf("Path should not be absolute: %s", line)
		}
	}
	
	// Test 2: Filter by extension
	result2 := listFiles(map[string]string{"directory": tmpDir, "extension": ".go"}, 0)
	if !result2.Success {
		t.Fatalf("Expected success, got error: %s", result2.Error)
	}
	
	output2 := result2.Output
	lines2 := strings.Split(strings.TrimSpace(output2), "\n")
	if len(lines2) != 1 || !strings.Contains(lines2[0], "file3.go") {
		t.Errorf("Expected only file3.go, got: %v", lines2)
	}
	
	// Test 3: Empty directory
	emptyDir := t.TempDir()
	result3 := listFiles(map[string]string{"directory": emptyDir}, 0)
	if !result3.Success {
		t.Fatalf("Expected success, got error: %s", result3.Error)
	}
	if result3.Output != "(empty directory)" {
		t.Errorf("Expected '(empty directory)', got: %s", result3.Output)
	}
}

func TestListFiles_CurrentDirectory(t *testing.T) {
	// Test with "." as directory
	tmpDir := t.TempDir()
	
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Create a file in the current directory
	testFile := filepath.Join(tmpDir, "testfile.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	
	// Debug: check if file exists
	if _, err := os.Stat(testFile); err != nil {
		t.Fatalf("Test file doesn't exist: %v", err)
	}
	
	result := listFiles(map[string]string{"directory": "."}, 0)
	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.Error)
	}
	
	t.Logf("Output: %q", result.Output)
	
	if !strings.Contains(result.Output, "testfile.txt") {
		t.Errorf("Expected to find testfile.txt, got: %s", result.Output)
	}
}
