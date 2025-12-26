# Setup Instructions for Quint Code with GitHub Copilot

## Automatic Setup (Recommended)

Simply run one command in the root of the Quint Code project:

```bash
quint-code init --copilot
```

This will automatically create/update:
- ✅ `.vscode/settings.json` - Full VS Code configuration with GitHub Copilot + MCP
- ✅ `.vscode/extensions.json` - Recommended extensions list
- ✅ `.github/copilot-instructions.md` - Complete FPF methodology guide for Copilot

## What Gets Configured

### 1. VS Code Settings (`.vscode/settings.json`)

Automatically merges with existing settings:
- GitHub Copilot enabled for all file types
- Instruction files enabled (`github.copilot.chat.useInstructionFiles: true`)
- Go language server configuration
- Editor settings for Go and Markdown
- File associations and exclusions
- MCP server configuration

### 2. Extensions (`.vscode/extensions.json`)

Recommends essential extensions:
- `github.copilot` & `github.copilot-chat` - AI assistance
- `golang.go` - Go language support
- `mtxr.sqltools` + `mtxr.sqltools-driver-sqlite` - SQL tools
- `yzhang.markdown-all-in-one` - Markdown support
- `davidanson.vscode-markdownlint` - Markdown linting
- `redhat.vscode-yaml` & `tamasfe.even-better-toml` - Config files
- `eamodio.gitlens` - Git tools

### 3. Copilot Instructions (`.github/copilot-instructions.md`)

Complete guide including:
- FPF methodology (Abduction/Deduction/Induction)
- Knowledge layers (L0/L1/L2/invalid)
- MCP tools reference
- Key concepts (Holons, WLNK, Relations, Congruence Levels)
- Code patterns (Database, FSM, Tools, Server layers)
- Testing patterns
- Architecture guidelines
- Command examples
- Integration points

## Manual Setup (If Needed)

If you need to manually set up individual components:

### Settings Only
The settings are defined in `src/mcp/cmd/templates.go` as `vscodeSettingsBase`.

### Extensions Only
The extensions list is in `src/mcp/cmd/templates.go` as `vscodeExtensionsTemplate`.

### Instructions Only
The full instructions template is in `src/mcp/cmd/templates.go` as `copilotInstructionsTemplate`.

## Features

### JSONC Support
The init command now handles JSON files with comments (JSONC format) correctly, so existing VS Code settings won't break.

### Non-Destructive
- Existing settings are preserved and merged with new ones
- Files are only created if they don't exist or need updating
- MCP configuration is added without removing other settings

### Smart Merging
The command intelligently merges settings:
- Preserves existing custom settings
- Adds missing recommended settings
- Updates MCP server configuration
- Keeps file readable with proper indentation

## Rebuilding After Changes

If you modify templates in `templates.go`:

```bash
cd src/mcp
go build -o quint-code .
sudo cp quint-code /usr/local/bin/
```

Then re-run:
```bash
cd /path/to/quint-code
quint-code init --copilot
```

## Verification

After running `quint-code init --copilot`, verify:

```bash
# Check files exist
ls -la .vscode/settings.json
ls -la .vscode/extensions.json
ls -la .github/copilot-instructions.md

# Verify settings.json is valid JSON
jq . .vscode/settings.json > /dev/null && echo "✓ Valid JSON"

# Check file sizes
wc -l .github/copilot-instructions.md  # Should be ~250 lines
```

## Usage After Setup

1. Open project in VS Code: `code .`
2. Install recommended extensions (VS Code will prompt)
3. Reload window if needed (Ctrl+Shift+P → "Reload Window")
4. Open GitHub Copilot Chat (Ctrl+Shift+I or Cmd+Shift+I)
5. Test: Ask "What is the FPF methodology in this project?"

GitHub Copilot will now understand:
- Project structure and architecture
- FPF concepts and workflow
- Available MCP tools
- Code patterns and best practices
- Testing strategies

## Troubleshooting

### "Command not found"
```bash
# Install the binary
cd src/mcp
go build -o quint-code .
sudo cp quint-code /usr/local/bin/
sudo chmod +x /usr/local/bin/quint-code
```

### "Invalid JSON" error
The command now handles JSONC (JSON with comments) automatically. If you still see errors:
1. Check that `settings.json` doesn't have syntax errors
2. Backup and remove `.vscode/settings.json`
3. Re-run `quint-code init --copilot`

### Files not created
- Check permissions in project directory
- Ensure `.vscode/` and `.github/` directories are writable
- Check disk space

### GitHub Copilot not using instructions
1. Verify `.github/copilot-instructions.md` exists
2. Check VS Code setting: `"github.copilot.chat.useInstructionFiles": true`
3. Reload VS Code window
4. Restart VS Code completely

## Comparison: Before vs After

### Before
- Manual file creation required
- Settings had to be manually merged
- Risk of JSON syntax errors with comments
- Instructions file needed manual copying

### After
✅ One command setup: `quint-code init --copilot`  
✅ Automatic JSONC handling  
✅ Smart settings merging  
✅ Full instructions template embedded  
✅ Extensions recommendations included  
✅ Non-destructive updates  

## Developer Notes

### Template Location
All templates are in: `src/mcp/cmd/templates.go`

### Key Functions
- `configureMCPCopilot()` - Main setup function
- `removeJSONComments()` - JSONC parser helper
- Templates are Go string constants with proper escaping

### Adding New Templates
1. Add constant to `templates.go`
2. Update `configureMCPCopilot()` to use it
3. Rebuild and test

### Testing Changes
```bash
# Clean test
rm -rf /tmp/test-copilot
mkdir /tmp/test-copilot
cd /tmp/test-copilot
quint-code init --copilot

# Verify output
ls -la .vscode/ .github/
```

---

**Result**: GitHub Copilot integration is now fully automated! 🎉
