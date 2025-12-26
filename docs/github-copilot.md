# Using Quint Code with GitHub Copilot

This guide helps you get started with Quint Code in VS Code with GitHub Copilot.

## Prerequisites

1. **VS Code** with **GitHub Copilot** extension installed
2. **Quint Code** binary installed (see main README.md)
3. **Go** language support (recommended: `golang.go` extension)

## Quick Setup

### 1. Initialize Your Project

```bash
cd /path/to/your/project
quint-code init --copilot
```

This command will:
- Create `.quint/` directory structure
- Configure `.vscode/settings.json` with MCP server settings
- Create `.github/copilot-instructions.md` with FPF context
- Set up recommended VS Code extensions

### 2. Verify Configuration

Check that the following files were created/updated:

```
your-project/
├── .github/
│   └── copilot-instructions.md    # GitHub Copilot context
├── .vscode/
│   ├── settings.json               # MCP + Copilot config
│   └── extensions.json             # Recommended extensions
└── .quint/
    ├── knowledge/                  # L0/L1/L2 layers
    ├── evidence/
    ├── decisions/
    └── quint.db
```

### 3. Start the MCP Server

The MCP server should start automatically when you open the project in VS Code. If not, you can start it manually:

```bash
quint-code serve
```

The server runs in the background and communicates with GitHub Copilot via the Model Context Protocol.

## Using FPF with GitHub Copilot

### Basic Workflow

1. **Initialize Knowledge Base**
   ```
   Ask Copilot: "Initialize the FPF knowledge base"
   Or in chat: /q0-init
   ```

2. **Generate Hypotheses** (Abduction Phase)
   ```
   Ask: "Generate hypotheses for [your problem]"
   Or: /q1-hypothesize "problem description"
   ```

3. **Verify Logic** (Deduction Phase)
   ```
   Ask: "Verify hypothesis h-001"
   Or: /q2-verify h-001
   ```

4. **Validate Empirically** (Induction Phase)
   ```
   Ask: "Validate hypothesis h-001 with tests"
   Or: /q3-validate h-001
   ```

5. **Audit & Decide**
   ```
   Ask: "Audit hypothesis h-001"
   Ask: "Make decision with h-001 as winner"
   ```

### Querying the Knowledge Base

GitHub Copilot can access your `.quint/` knowledge base:

```
Ask: "What do we know about caching strategies?"
Ask: "Show me L2-validated decisions"
Ask: "Check for stale evidence"
```

### Example Conversation

**You:** "We need to add caching to our API. What are the options?"

**Copilot:** (uses `quint_init` and `quint_propose` tools)
```
Generated 3 hypotheses:
- h-001: Redis in-memory cache
- h-002: CDN edge caching
- h-003: Database query cache

Use /q2-verify h-001 to verify logic
```

**You:** "Verify h-001"

**Copilot:** (uses `quint_verify` tool)
```
Verification PASS:
✓ Preconditions: Redis available
✓ No contradictions
✓ Constraints satisfied

Hypothesis promoted to L1 (Substantiated)
Use /q3-validate h-001 to test empirically
```

## FPF Concepts GitHub Copilot Understands

### Knowledge Layers
- **L0**: Unverified hypotheses (conjectures)
- **L1**: Logically verified claims (substantiated)
- **L2**: Empirically tested claims (corroborated)
- **invalid**: Falsified claims (kept for learning)

### Relations
- **ComponentOf**: System dependencies (affects R_eff via WLNK)
- **MemberOf**: Decision alternatives (grouping only)
- **ConstituentOf**: Epistemic dependencies

### Key Principles
- **WLNK**: Weakest Link — R_eff = min(dependencies), never average
- **Congruence**: Evidence quality depends on context match
- **Transformer Mandate**: AI proposes, human decides

## Advanced Usage

### Custom Instructions

Edit `.github/copilot-instructions.md` to add project-specific context:

```markdown
## Project-Specific Context

- **Stack**: Go, PostgreSQL, Redis
- **Patterns**: Repository pattern, dependency injection
- **Constraints**: Must support 10k RPS, <100ms p99
```

### MCP Configuration

Customize MCP server settings in `.vscode/settings.json`:

```json
{
  "mcp.servers": {
    "quint-code": {
      "command": "/usr/local/bin/quint-code",
      "args": ["serve"],
      "env": {
        "QUINT_PROJECT_ROOT": "${workspaceFolder}"
      }
    }
  }
}
```

### Keyboard Shortcuts

Add to `.vscode/keybindings.json`:

```json
[
  {
    "key": "ctrl+shift+q",
    "command": "github.copilot.chat.open",
    "args": "/q-status"
  }
]
```

## Troubleshooting

### MCP Server Not Starting

1. Check if binary is in PATH:
   ```bash
   which quint-code
   ```

2. Test server manually:
   ```bash
   cd /path/to/project
   quint-code serve
   ```

3. Check VS Code output panel (View → Output → GitHub Copilot)

### GitHub Copilot Not Using Instructions

1. Verify `.github/copilot-instructions.md` exists
2. Check setting: `"github.copilot.chat.useInstructionFiles": true`
3. Reload VS Code window (Ctrl+Shift+P → "Reload Window")

### Knowledge Base Not Found

Make sure:
- `.quint/` directory exists
- You ran `quint-code init` or `/q0-init`
- `QUINT_PROJECT_ROOT` env variable points to correct path

## Best Practices

### 1. Start with Context
Always initialize the knowledge base before making decisions:
```
/q0-init
```

### 2. Use Structured Workflow
Follow the FPF cycle:
```
Abduction → Deduction → Induction → Audit → Decision
```

### 3. Keep Evidence Fresh
Regularly check for stale evidence:
```
Ask: "Check for stale evidence"
```

### 4. Document Decisions
Always create Design Rationale Records (DRRs):
```
/q5-decide
```

### 5. Query Before Creating
Check existing knowledge before proposing new hypotheses:
```
Ask: "What do we know about [topic]?"
```

## Resources

- **Main Documentation**: [../README.md](../README.md)
- **FPF Theory**: [docs/advanced.md](../docs/advanced.md)
- **Architecture**: [docs/architecture.md](../docs/architecture.md)
- **Workflow Examples**: [docs/workflow_example/](../docs/workflow_example/)
- **Agent Instructions**: [CLAUDE.md](../CLAUDE.md)

## Getting Help

If you encounter issues:

1. Check [GitHub Issues](https://github.com/m0n0x41d/quint-code/issues)
2. Read the [main documentation](../README.md)
3. Ask GitHub Copilot: "Explain how FPF works in this project"

---

**Remember**: GitHub Copilot is a powerful assistant, but YOU make the final decisions. The Transformer Mandate applies: AI generates options, humans decide.
