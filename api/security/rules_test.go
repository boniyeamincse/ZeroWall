package security

import (
	"os"
	"strings"
	"testing"
)

func TestRuleManagement(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rules_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ruleContent := "pass in on em0 proto tcp from any to any port 80 sid:10001;\npass in on em0 proto tcp from any to any port 443 sid:10002;"
	ruleFile := tmpDir + "/test.rules"
	if err := os.WriteFile(ruleFile, []byte(ruleContent), 0644); err != nil {
		t.Fatal(err)
	}

	mgr := RulesManager{RulesDir: tmpDir}

	// Test Disable
	if err := mgr.DisableRule("10001"); err != nil {
		t.Errorf("DisableRule failed: %v", err)
	}

	content, _ := os.ReadFile(ruleFile)
	if !strings.Contains(string(content), "# pass in on em0 proto tcp from any to any port 80 sid:10001;") {
		t.Errorf("Rule 10001 not disabled correctly")
	}

	// Test Enable
	if err := mgr.EnableRule("10001"); err != nil {
		t.Errorf("EnableRule failed: %v", err)
	}

	content, _ = os.ReadFile(ruleFile)
	if strings.Contains(string(content), "# pass in on em0 proto tcp from any to any port 80 sid:10001;") {
		t.Errorf("Rule 10001 not enabled correctly")
	}
}
