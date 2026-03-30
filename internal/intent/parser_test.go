package intent

import (
	"testing"
)

// TestParse_JSONBlock verifies extraction from ```json ... ``` code blocks.
func TestParse_JSONBlock(t *testing.T) {
	content := "I'll write the file now.\n```json\n{\"tool\": \"write_file\", \"params\": {\"path\": \"main.go\", \"content\": \"package main\"}}\n```"
	a := Parse(content)
	if a == nil {
		t.Fatal("expected action, got nil")
	}
	if a.Tool != "write_file" {
		t.Errorf("Tool = %q, want %q", a.Tool, "write_file")
	}
	if a.Params["path"] != "main.go" {
		t.Errorf("params[path] = %q, want %q", a.Params["path"], "main.go")
	}
	if a.Source != SourceJSONBlock {
		t.Errorf("Source = %q, want %q", a.Source, SourceJSONBlock)
	}
}

// TestParse_XMLTag verifies extraction from <tool>...</tool> XML tags.
func TestParse_XMLTag(t *testing.T) {
	content := `Running the command now.<tool>{"tool": "run_shell", "params": {"command": "go test ./..."}}</tool>`
	a := Parse(content)
	if a == nil {
		t.Fatal("expected action, got nil")
	}
	if a.Tool != "run_shell" {
		t.Errorf("Tool = %q, want %q", a.Tool, "run_shell")
	}
	if a.Params["command"] != "go test ./..." {
		t.Errorf("params[command] = %q, want %q", a.Params["command"], "go test ./...")
	}
	if a.Source != SourceXMLTag {
		t.Errorf("Source = %q, want %q", a.Source, SourceXMLTag)
	}
}

// TestParse_FunctionCall verifies extraction from OpenAI function_call format.
func TestParse_FunctionCall(t *testing.T) {
	content := `{"name": "read_file", "arguments": "{\"path\": \"/etc/hosts\"}"}`
	a := Parse(content)
	if a == nil {
		t.Fatal("expected action, got nil")
	}
	if a.Tool != "read_file" {
		t.Errorf("Tool = %q, want %q", a.Tool, "read_file")
	}
	if a.Params["path"] != "/etc/hosts" {
		t.Errorf("params[path] = %q, want %q", a.Params["path"], "/etc/hosts")
	}
	if a.Source != SourceFunctionCall {
		t.Errorf("Source = %q, want %q", a.Source, SourceFunctionCall)
	}
}

// TestParse_BareJSON verifies extraction from inline JSON objects.
func TestParse_BareJSON(t *testing.T) {
	content := `Let me list the files: {"tool": "list_files", "directory": "/tmp"}`
	a := Parse(content)
	if a == nil {
		t.Fatal("expected action, got nil")
	}
	if a.Tool != "list_files" {
		t.Errorf("Tool = %q, want %q", a.Tool, "list_files")
	}
	if a.Source != SourceBareJSON {
		t.Errorf("Source = %q, want %q", a.Source, SourceBareJSON)
	}
}

// TestParse_NoAction verifies that plain prose returns nil.
func TestParse_NoAction(t *testing.T) {
	cases := []string{
		"I've finished the task. The code looks good.",
		"Based on the analysis, the bug is in line 42.",
		"",
		"Here is the answer: 42.",
	}
	for _, c := range cases {
		if a := Parse(c); a != nil {
			t.Errorf("Parse(%q) = %+v, want nil", c, a)
		}
	}
}

// TestParse_ToolAliases verifies that model-emitted aliases map to canonical names.
func TestParse_ToolAliases(t *testing.T) {
	cases := []struct {
		raw      string
		wantTool string
	}{
		{`{"tool": "Bash", "params": {"command": "ls"}}`, "run_shell"},
		{`{"tool": "Read", "params": {"path": "main.go"}}`, "read_file"},
		{`{"tool": "Write", "params": {"path": "out.go", "content": "x"}}`, "write_file"},
		{`{"tool": "Glob", "params": {"directory": "."}}`, "list_files"},
		{`{"tool": "Grep", "params": {"directory": ".", "pattern": "foo"}}`, "search_files"},
	}
	for _, tc := range cases {
		content := "```json\n" + tc.raw + "\n```"
		a := Parse(content)
		if a == nil {
			t.Errorf("Parse(%q): got nil, want tool %q", tc.raw, tc.wantTool)
			continue
		}
		if a.Tool != tc.wantTool {
			t.Errorf("Parse(%q): Tool = %q, want %q", tc.raw, a.Tool, tc.wantTool)
		}
	}
}

// TestParse_ParamAliases verifies that param aliases are normalized.
func TestParse_ParamAliases(t *testing.T) {
	content := "```json\n{\"tool\": \"write_file\", \"params\": {\"file_path\": \"main.go\", \"text\": \"hello\"}}\n```"
	a := Parse(content)
	if a == nil {
		t.Fatal("expected action, got nil")
	}
	if a.Params["path"] != "main.go" {
		t.Errorf("file_path should normalize to path, got %q", a.Params["path"])
	}
	if a.Params["content"] != "hello" {
		t.Errorf("text should normalize to content, got %q", a.Params["content"])
	}
}

// TestParse_UnknownTool verifies that unknown tool names return nil.
func TestParse_UnknownTool(t *testing.T) {
	content := "```json\n{\"tool\": \"do_something_weird\", \"params\": {}}\n```"
	a := Parse(content)
	if a != nil {
		t.Errorf("Parse with unknown tool: got %+v, want nil", a)
	}
}

// TestFlattenParams verifies numeric and bool conversions.
func TestFlattenParams(t *testing.T) {
	input := map[string]any{
		"name":    "foo",
		"count":   float64(42),
		"ratio":   float64(3.14),
		"enabled": true,
	}
	got := flattenParams(input)

	if got["name"] != "foo" {
		t.Errorf("name = %q, want %q", got["name"], "foo")
	}
	if got["count"] != "42" {
		t.Errorf("count = %q, want %q", got["count"], "42")
	}
	if got["ratio"] != "3.14" {
		t.Errorf("ratio = %q, want %q", got["ratio"], "3.14")
	}
	if got["enabled"] != "true" {
		t.Errorf("enabled = %q, want %q", got["enabled"], "true")
	}
}
