package security

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// RulesManager handles IDS ruleset updates
type RulesManager struct {
	RulesDir string
	Sources  []string
}

// DownloadLatestRules fetches rules from configured OISF/EtOpen sources
func (m RulesManager) DownloadLatestRules() error {
	if err := os.MkdirAll(m.RulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create rules directory: %v", err)
	}

	for _, source := range m.Sources {
		fmt.Printf("Downloading rules from %s...\n", source)
		resp, err := http.Get(source)
		if err != nil {
			return fmt.Errorf("failed to download from %s: %v", source, err)
		}
		defer resp.Body.Close()

		filename := "latest_rules.tar.gz" // Simplified for logic demonstration
		filePath := fmt.Sprintf("%s/%s", m.RulesDir, filename)
		out, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create rule file: %v", err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save rule file: %v", err)
		}
		
		fmt.Printf("Rules from %s downloaded and saved to %s.\n", source, filePath)
		// Logic to extract tar.gz/zip would follow here
	}
	return nil
}

// DisableRule comments out a specific rule by SID in files within RulesDir
func (m RulesManager) DisableRule(sid string) error {
	files, err := os.ReadDir(m.RulesDir)
	if err != nil {
		return err
	}

	target := fmt.Sprintf("sid:%s;", sid)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".rules") {
			continue
		}

		path := m.RulesDir + "/" + file.Name()
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		modified := false
		for i, line := range lines {
			if strings.Contains(line, target) && !strings.HasPrefix(strings.TrimSpace(line), "#") {
				lines[i] = "# " + line
				modified = true
			}
		}

		if modified {
			err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Disabled rule %s in %s\n", sid, file.Name())
		}
	}
	return nil
}

// EnableRule uncomments a specific rule by SID in files within RulesDir
func (m RulesManager) EnableRule(sid string) error {
	files, err := os.ReadDir(m.RulesDir)
	if err != nil {
		return err
	}

	target := fmt.Sprintf("sid:%s;", sid)
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".rules") {
			continue
		}

		path := m.RulesDir + "/" + file.Name()
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		modified := false
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(line, target) && strings.HasPrefix(trimmed, "#") {
				// Remove the leading # and any following space
				lines[i] = strings.TrimPrefix(trimmed, "#")
				lines[i] = strings.TrimSpace(lines[i])
				modified = true
			}
		}

		if modified {
			err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Enabled rule %s in %s\n", sid, file.Name())
		}
	}
	return nil
}
