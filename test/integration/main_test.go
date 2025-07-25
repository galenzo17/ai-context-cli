package integration

import (
	"os/exec"
	"testing"
)

func TestCLIVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/ai-context-cli/main.go", "version")
	output, err := cmd.Output()
	
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
	}
	
	expected := "ai-context-cli v0.1.0\n"
	if string(output) != expected {
		t.Errorf("Expected %q, got %q", expected, string(output))
	}
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("go", "run", "../../cmd/ai-context-cli/main.go", "help")
	output, err := cmd.Output()
	
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
	}
	
	if len(output) == 0 {
		t.Error("Expected help output to be non-empty")
	}
}