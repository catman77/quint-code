# Quick Start: GitHub Copilot Integration

## TL;DR

```bash
# Install Quint Code
curl -fsSL https://raw.githubusercontent.com/m0n0x41d/quint-code/main/install.sh | bash

# Initialize your project with GitHub Copilot support
cd /path/to/your/project
quint-code init --copilot

# Open in VS Code
code .

# In GitHub Copilot Chat, type:
/q0-init
```

## What Gets Configured

✅ `.quint/` directory structure (knowledge base)  
✅ `.vscode/settings.json` (MCP server + Copilot config)  
✅ `.github/copilot-instructions.md` (FPF context for Copilot)  
✅ `.vscode/extensions.json` (recommended extensions)

## First Commands

```
/q0-init                    # Initialize knowledge base
/q1-hypothesize "problem"   # Generate hypotheses
/q2-verify h-001            # Verify logic
/q3-validate h-001          # Test empirically
/q5-decide                  # Make decision
```

## Example Usage

**You**: "We need to add caching to our API"

**Copilot**: 
```
Generated 3 hypotheses (L0):
- h-001: Redis in-memory cache
- h-002: CDN edge caching  
- h-003: Database query cache

Use /q2-verify h-001 to check logic
```

**You**: "/q2-verify h-001"

**Copilot**:
```
✓ Verification PASS
  Promoted to L1 (Substantiated)
  
Use /q3-validate h-001 to test empirically
```

## Learn More

- [Full GitHub Copilot Guide](docs/github-copilot.md)
- [Main Documentation](README.md)
- [FPF Methodology](docs/advanced.md)

## Troubleshooting

**GitHub Copilot not using instructions?**
1. Check setting: `"github.copilot.chat.useInstructionFiles": true`
2. Verify `.github/copilot-instructions.md` exists
3. Reload VS Code (Ctrl+Shift+P → "Reload Window")

**MCP server not starting?**
```bash
# Test manually:
cd /path/to/project
quint-code serve
```

---

**Ready to go!** Open GitHub Copilot Chat and start with `/q0-init` 🚀
