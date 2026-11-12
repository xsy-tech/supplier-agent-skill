package installer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstalledInventoryScript(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 unavailable")
	}
	home := t.TempDir()
	if err := InstallSkill("codex", home, "", false); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(home, ".codex", "skills", "supplier-open-platform", "scripts", "fabric_inventory_excel.py")
	output := filepath.Join(home, "template.xlsx")
	if data, err := exec.Command(python, script, "template", "--output", output).CombinedOutput(); err != nil {
		t.Fatalf("installed template: %v %s", err, data)
	}
	if data, err := exec.Command(python, script, "template", "--output", output).CombinedOutput(); err == nil {
		t.Fatalf("overwrote template: %s", data)
	}
}

func TestInstalledInventoryWorkflowUsesSeparateWriteToolsAndKeys(t *testing.T) {
	home := t.TempDir()
	if err := InstallSkill("codex", home, "", false); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(home, ".codex", "skills", "supplier-open-platform")
	for _, name := range []string{"SKILL.md", "references/fabric-inventory.md"} {
		content, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, required := range []string{"commit_fabric_inventory_import", "commit_fabric_inventory_adjustments", "idempotency_key"} {
			if !strings.Contains(string(content), required) {
				t.Errorf("%s omits %s", name, required)
			}
		}
		if strings.Contains(string(content), "commit_fabric_inventory_changes") || strings.Contains(string(content), "提交无幂等键") {
			t.Errorf("%s still instructs obsolete shared tool or keyless commit", name)
		}
	}
}
