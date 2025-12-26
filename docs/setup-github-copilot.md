# GitHub Copilot Setup for Quint Code Users

## For End Users (Using Quint Code in Your Projects)

If you want to use Quint Code with GitHub Copilot in your own projects:

### 1. Install Quint Code

```bash
curl -fsSL https://raw.githubusercontent.com/m0n0x41d/quint-code/main/install.sh | bash
```

### 2. Initialize Your Project

```bash
cd /path/to/your/project
quint-code init --copilot
```

This will:
- ✓ Create `.quint/` directory structure
- ✓ Configure `.vscode/settings.json` with MCP server
- ✓ Create `.github/copilot-instructions.md` with FPF context
- ✓ Recommend VS Code extensions

### 3. Install GitHub Copilot Extension

In VS Code:
1. Open Extensions (Ctrl+Shift+X)
2. Install "GitHub Copilot" and "GitHub Copilot Chat"
3. Sign in with your GitHub account

### 4. Start Using FPF

Open GitHub Copilot Chat and try:
```
/q0-init
```

See [docs/github-copilot.md](docs/github-copilot.md) for complete guide.

---

## For Quint Code Developers

If you're contributing to the Quint Code project itself:

### The `.github/copilot-instructions.md` File

This file provides GitHub Copilot with context about:
- FPF methodology and concepts
- Project architecture (Functional Core + Imperative Shell)
- Code patterns and testing strategies
- Common development workflows

**Location**: `.github/copilot-instructions.md`

### VS Code Settings

The `.vscode/settings.json` includes:
- GitHub Copilot configuration
- Go language server settings
- Editor preferences for Go and Markdown
- File associations

**Location**: `.vscode/settings.json`

### Recommended Extensions

The `.vscode/extensions.json` recommends:
- `github.copilot` - GitHub Copilot
- `github.copilot-chat` - GitHub Copilot Chat
- `golang.go` - Go language support
- `mtxr.sqltools` - SQL tools for schema.sql/query.sql
- And more...

**Location**: `.vscode/extensions.json`

### When Contributing Code

GitHub Copilot will:
1. **Understand FPF concepts** (L0/L1/L2, WLNK, R_eff, etc.)
2. **Follow project patterns** (functional core, error handling, testing)
3. **Suggest context-aware code** based on neighboring files
4. **Know the architecture** (MCP server, FSM, Tools layer, DB layer)

Example prompts that work well:
```
"Add a new MCP tool for exporting decisions to JSON"
"Write integration test for quint_verify with database"
"Implement congruence level calculation in assurance calculator"
```

### Testing GitHub Copilot Integration

After making changes to instructions or configuration:

1. Reload VS Code window
2. Open GitHub Copilot Chat
3. Ask: "What is the FPF methodology in this project?"
4. Verify it understands L0/L1/L2 layers and core concepts

### Updating Instructions

When adding new features or changing architecture:

1. Update `.github/copilot-instructions.md`
2. Update `CLAUDE.md` for consistency
3. Test with GitHub Copilot to ensure understanding

---

## Architecture Notes

### MCP Integration

GitHub Copilot can use MCP (Model Context Protocol) tools defined in `server.go`:
- `quint_init`
- `quint_propose`
- `quint_verify`
- `quint_test`
- `quint_audit`
- `quint_decide`
- And more...

The MCP server configuration is added to `.vscode/settings.json` under `mcp.servers`.

### File Locations

```
quint-code/
├── .github/
│   └── copilot-instructions.md    # GitHub Copilot context (auto-generated for users)
├── .vscode/
│   ├── settings.json               # MCP + Copilot config
│   └── extensions.json             # Recommended extensions
├── docs/
│   └── github-copilot.md          # User guide for GitHub Copilot
└── src/mcp/cmd/
    └── init.go                     # Contains configureMCPCopilot() function
```

### How It Works

1. User runs `quint-code init --copilot`
2. `configureMCPCopilot()` in `init.go`:
   - Creates/updates `.vscode/settings.json` with MCP config
   - Enables GitHub Copilot features
   - Creates `.github/copilot-instructions.md` if missing
3. GitHub Copilot reads instructions from `.github/copilot-instructions.md`
4. MCP server starts when VS Code opens the project
5. GitHub Copilot can invoke MCP tools via JSON-RPC

### Difference from Claude Code

| Feature | Claude Code | GitHub Copilot |
|---------|-------------|----------------|
| MCP Config | `.mcp.json` (root) | `.vscode/settings.json` |
| Instructions | `CLAUDE.md` (root) | `.github/copilot-instructions.md` |
| Slash Commands | `~/.claude/commands/*.md` | Native MCP tool invocation |
| Setup | `--claude` (default) | `--copilot` (explicit) |

Both can use the same MCP server (`quint-code serve`), but configuration differs.

---

## FAQ

**Q: Can I use both Claude Code and GitHub Copilot?**  
A: Yes! Run `quint-code init --all` to configure both.

**Q: Will this work with other IDEs?**  
A: Currently optimized for VS Code. Other IDEs would need similar configuration.

**Q: Do I need to restart the MCP server when code changes?**  
A: No, the server automatically reloads. But you may need to reload VS Code window.

**Q: Where are instructions stored?**  
A: For users: `.github/copilot-instructions.md` (in their project)  
For developers: Same file, but in quint-code repo itself.

---

## Contributing

When improving GitHub Copilot integration:

1. Test with real prompts in GitHub Copilot Chat
2. Ensure instructions are clear and concise
3. Update both `.github/copilot-instructions.md` and `CLAUDE.md`
4. Document new patterns in this file
5. Add examples to `docs/github-copilot.md`
