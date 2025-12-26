package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m0n0x41d/quint-code/db"

	"github.com/spf13/cobra"
)

var (
	initClaude  bool
	initCursor  bool
	initGemini  bool
	initCodex   bool
	initCopilot bool
	initAll     bool
	initLocal   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize FPF project structure and MCP configuration",
	Long: `Initialize a new Quint Code project in the current directory.

This command creates:
  - .quint/ directory structure (knowledge base, evidence, decisions)
  - MCP configuration for selected AI tools
  - Slash commands (global by default, or local with --local)

Examples:
  quint-code init              # Claude, global commands (~/.claude/commands/)
  quint-code init --local      # Claude, local commands (.claude/commands/)
  quint-code init --all        # All tools, global commands
  quint-code init --cursor     # Cursor only
  quint-code init --copilot    # GitHub Copilot (VS Code)`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initClaude, "claude", false, "Configure for Claude Code")
	initCmd.Flags().BoolVar(&initCursor, "cursor", false, "Configure for Cursor")
	initCmd.Flags().BoolVar(&initGemini, "gemini", false, "Configure for Gemini CLI")
	initCmd.Flags().BoolVar(&initCodex, "codex", false, "Configure for Codex CLI")
	initCmd.Flags().BoolVar(&initCopilot, "copilot", false, "Configure for GitHub Copilot (VS Code)")
	initCmd.Flags().BoolVar(&initAll, "all", false, "Configure for all supported tools")
	initCmd.Flags().BoolVar(&initLocal, "local", false, "Install commands in project directory instead of global")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	quintDir := filepath.Join(cwd, ".quint")
	dbPath := filepath.Join(quintDir, "quint.db")

	_, quintExists := os.Stat(quintDir)
	_, dbExists := os.Stat(dbPath)

	fmt.Println("Initializing Quint Code project...")

	if err := createDirectoryStructure(quintDir); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	if os.IsNotExist(quintExists) {
		fmt.Println("  ✓ Created .quint/ directory structure")
	} else {
		fmt.Println("  ✓ .quint/ directory structure OK")
	}

	if err := initializeDatabase(quintDir); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	if os.IsNotExist(dbExists) {
		fmt.Println("  ✓ Initialized database")
	} else {
		fmt.Println("  ✓ Database OK")
	}

	binaryPath, err := getBinaryPath()
	if err != nil {
		fmt.Printf("  ⚠ Could not determine binary path: %v\n", err)
		binaryPath = "quint-code"
	}

	if initAll {
		initClaude, initCursor, initGemini, initCodex, initCopilot = true, true, true, true, true
	}

	if !initClaude && !initCursor && !initGemini && !initCodex && !initCopilot {
		initClaude = true
	}

	if initClaude {
		if err := configureMCPClaude(cwd, binaryPath); err != nil {
			fmt.Printf("  ⚠ Failed to configure Claude Code MCP: %v\n", err)
		} else {
			fmt.Println("  ✓ Configured MCP for Claude Code (.mcp.json)")
		}
		if destPath, count, err := installCommands(cwd, "claude", initLocal); err != nil {
			fmt.Printf("  ⚠ Failed to install Claude commands: %v\n", err)
		} else {
			fmt.Printf("  ✓ Installed %d slash commands (%s)\n", count, destPath)
		}
	}

	if initCursor {
		if err := configureMCPCursor(cwd, binaryPath); err != nil {
			fmt.Printf("  ⚠ Failed to configure Cursor MCP: %v\n", err)
		} else {
			fmt.Println("  ✓ Configured MCP for Cursor (.cursor/mcp.json)")
			fmt.Println("    Note: Make sure quint-code MCP is enabled in Cursor settings")
		}
		if destPath, count, err := installCommands(cwd, "cursor", initLocal); err != nil {
			fmt.Printf("  ⚠ Failed to install Cursor commands: %v\n", err)
		} else {
			fmt.Printf("  ✓ Installed %d slash commands (%s)\n", count, destPath)
		}
	}

	if initGemini {
		if err := configureMCPGemini(cwd, binaryPath); err != nil {
			fmt.Printf("  ⚠ Failed to configure Gemini CLI MCP: %v\n", err)
		} else {
			fmt.Printf("  ✓ Configured MCP for Gemini CLI (project: %s)\n", cwd)
		}
		if destPath, count, err := installCommands(cwd, "gemini", initLocal); err != nil {
			fmt.Printf("  ⚠ Failed to install Gemini commands: %v\n", err)
		} else {
			fmt.Printf("  ✓ Installed %d slash commands (%s)\n", count, destPath)
		}
	}

	if initCodex {
		if err := configureMCPCodex(cwd, binaryPath); err != nil {
			fmt.Printf("  ⚠ Failed to configure Codex CLI MCP: %v\n", err)
		} else {
			fmt.Printf("  ✓ Configured MCP for Codex CLI (project: %s)\n", cwd)
		}
		// Codex only supports global prompts
		if destPath, count, err := installCommands(cwd, "codex", false); err != nil {
			fmt.Printf("  ⚠ Failed to install Codex prompts: %v\n", err)
		} else {
			fmt.Printf("  ✓ Installed %d prompts (%s)\n", count, destPath)
			fmt.Println("    Note: Use /prompts:q0-init to invoke")
		}
	}

	if initCopilot {
		if err := configureMCPCopilot(cwd, binaryPath); err != nil {
			fmt.Printf("  ⚠ Failed to configure GitHub Copilot: %v\n", err)
		} else {
			fmt.Println("  ✓ Configured MCP for GitHub Copilot (.vscode/settings.json)")
			fmt.Println("    Note: Make sure GitHub Copilot extension is installed and enabled")
		}
	}

	fmt.Println("\nInitialization complete! Run /q0-init to start.")
	return nil
}

func createDirectoryStructure(quintDir string) error {
	dirs := []string{
		"evidence",
		"decisions",
		"sessions",
		"knowledge/L0",
		"knowledge/L1",
		"knowledge/L2",
		"knowledge/invalid",
		"agents",
	}

	for _, d := range dirs {
		path := filepath.Join(quintDir, d)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		gitkeep := filepath.Join(path, ".gitkeep")
		if err := os.WriteFile(gitkeep, []byte(""), 0644); err != nil {
			return err
		}
	}
	return nil
}

func initializeDatabase(quintDir string) error {
	dbPath := filepath.Join(quintDir, "quint.db")
	database, err := db.NewStore(dbPath)
	if err != nil {
		return err
	}
	_ = database.Close()
	return nil
}

func getBinaryPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(exe)
}

type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Cwd     string            `json:"cwd,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Timeout int               `json:"timeout,omitempty"`
}

func mergeMCPConfig(configPath, binaryPath, projectRoot string, extraFields map[string]interface{}) error {
	var config MCPConfig

	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("existing config at %s is not valid JSON: %w", configPath, err)
		}
	}

	if config.MCPServers == nil {
		config.MCPServers = make(map[string]MCPServer)
	}

	server := MCPServer{
		Command: binaryPath,
		Args:    []string{"serve"},
		Cwd:     projectRoot,
		Env: map[string]string{
			"QUINT_PROJECT_ROOT": projectRoot,
		},
	}

	if timeout, ok := extraFields["timeout"].(int); ok {
		server.Timeout = timeout
	}

	config.MCPServers["quint-code"] = server

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func configureMCPClaude(projectRoot, binaryPath string) error {
	configPath := filepath.Join(projectRoot, ".mcp.json")
	return mergeMCPConfig(configPath, binaryPath, projectRoot, nil)
}

func configureMCPCursor(projectRoot, binaryPath string) error {
	configPath := filepath.Join(projectRoot, ".cursor", "mcp.json")
	return mergeMCPConfig(configPath, binaryPath, projectRoot, nil)
}

func configureMCPGemini(projectRoot, binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, ".gemini", "settings.json")
	return mergeMCPConfig(configPath, binaryPath, projectRoot, map[string]interface{}{
		"timeout": 30000,
	})
}

func configureMCPCodex(projectRoot, binaryPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, ".codex", "config.toml")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	existing := ""
	if data, err := os.ReadFile(configPath); err == nil {
		existing = string(data)
	}

	tomlSection := fmt.Sprintf(`
[mcp_servers.quint-code]
command = "%s"
args = ["serve"]
env = { QUINT_PROJECT_ROOT = "%s" }
`, binaryPath, projectRoot)

	if start := strings.Index(existing, "[mcp_servers.quint-code]"); start != -1 {
		end := len(existing)
		if nextSection := strings.Index(existing[start+1:], "\n["); nextSection != -1 {
			end = start + 1 + nextSection
		}
		existing = existing[:start] + existing[end:]
	}

	updated := strings.TrimRight(existing, "\n") + tomlSection

	return os.WriteFile(configPath, []byte(updated), 0644)
}

// removeJSONComments removes single-line comments from JSON (JSONC support)
func removeJSONComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	var cleaned []string
	for _, line := range lines {
		// Remove single-line comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, line)
		}
	}
	return []byte(strings.Join(cleaned, "\n"))
}

func configureMCPCopilot(projectRoot, binaryPath string) error {
	vscodeDir := filepath.Join(projectRoot, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return err
	}

	// Create/update settings.json
	settingsPath := filepath.Join(vscodeDir, "settings.json")
	var settings map[string]interface{}

	// Load existing settings if present
	if data, err := os.ReadFile(settingsPath); err == nil {
		// Remove comments for JSON parsing
		cleanedData := removeJSONComments(data)
		if err := json.Unmarshal(cleanedData, &settings); err != nil {
			return fmt.Errorf("existing settings.json is not valid JSON: %w", err)
		}
	}

	// Load base settings from template
	var baseSettings map[string]interface{}
	if err := json.Unmarshal([]byte(vscodeSettingsBase), &baseSettings); err != nil {
		return fmt.Errorf("failed to parse base settings template: %w", err)
	}

	// Merge settings (preserve existing, add new)
	if settings == nil {
		settings = baseSettings
	} else {
		for k, v := range baseSettings {
			if _, exists := settings[k]; !exists {
				settings[k] = v
			}
		}
	}

	// Add MCP configuration for VS Code
	mcpServers := make(map[string]interface{})
	if existing, ok := settings["mcp.servers"].(map[string]interface{}); ok {
		mcpServers = existing
	}

	mcpServers["quint-code"] = map[string]interface{}{
		"command": binaryPath,
		"args":    []string{"serve"},
		"env": map[string]string{
			"QUINT_PROJECT_ROOT": projectRoot,
		},
	}
	settings["mcp.servers"] = mcpServers

	// Write settings.json
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return err
	}

	// Create extensions.json
	extensionsPath := filepath.Join(vscodeDir, "extensions.json")
	if err := os.WriteFile(extensionsPath, []byte(vscodeExtensionsTemplate), 0644); err != nil {
		return err
	}

	// Create .github/copilot-instructions.md
	githubDir := filepath.Join(projectRoot, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return err
	}

	instructionsPath := filepath.Join(githubDir, "copilot-instructions.md")
	if err := os.WriteFile(instructionsPath, []byte(copilotInstructionsTemplate), 0644); err != nil {
		return err
	}

	return nil
}
