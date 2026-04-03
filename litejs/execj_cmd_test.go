package litejs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCmdNodeCallWithErrorRequiresJsPath(t *testing.T) {
	cmd := CmdNode{}

	_, err := cmd.CallWithError()
	if err == nil {
		t.Fatal("expected empty js path to return error")
	}
	if !strings.Contains(err.Error(), "js path is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCmdNodeCallWithErrorRunsNodeScript(t *testing.T) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed")
	}

	dir := t.TempDir()
	script := filepath.Join(dir, "echo.js")
	if err := os.WriteFile(script, []byte("console.log(process.argv.slice(2).join(','))\n"), 0o600); err != nil {
		t.Fatalf("write script: %v", err)
	}

	cmd := CmdNode{
		Node:   nodePath,
		JsPath: script,
	}
	output, err := cmd.CallWithError("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(string(output)); got != "a,b" {
		t.Fatalf("unexpected output: got %q want %q", got, "a,b")
	}
}
