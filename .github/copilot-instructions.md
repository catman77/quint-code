# GitHub Copilot Instructions for Quint Code

You are working with **Quint Code** — a structured reasoning framework for AI coding tools implementing the **First Principles Framework (FPF)** methodology.

## Core Philosophy

**Make decisions visible and auditable.** Quint Code turns black-box AI reasoning into transparent, evidence-backed trails that persist across sessions.

## Quick Context

- **Technology**: Go-based MCP (Model Context Protocol) server
- **Architecture**: Functional Core + Imperative Shell pattern
- **Database**: SQLite via sqlc (type-safe SQL queries)
- **Knowledge Base**: `.quint/` directory with L0/L1/L2 assurance layers
- **Primary Use**: Structured hypothesis → verification → validation → decision cycle

## FPF Methodology Overview

### Three Modes of Inference

1. **Abduction** — Generate competing hypotheses (don't anchor on first idea)
2. **Deduction** — Verify logic and constraints (does it make sense?)
3. **Induction** — Test empirically (does it work in reality?)

Then: **Audit** for bias → **Decide** → **Document** rationale

### Knowledge Layers (Epistemic Status)

| Layer | Meaning | How Reached |
|-------|---------|-------------|
| **L0** | Conjecture (unverified) | `quint_propose` |
| **L1** | Substantiated (logically verified) | `quint_verify` PASS |
| **L2** | Corroborated (empirically validated) | `quint_test` PASS |
| **invalid** | Falsified | Any FAIL verdict |

### State Machine

```
IDLE → ABDUCTION → DEDUCTION → INDUCTION → DECISION → IDLE
       (/q1)        (/q2)        (/q3)        (/q4→/q5)
```

## Available MCP Tools

When user types slash commands, these are exposed via MCP:

| Command | MCP Tool | Phase | Purpose |
|---------|----------|-------|---------|
| `/q0-init` | `quint_init` | Setup | Initialize `.quint/` structure |
| `/q1-hypothesize` | `quint_propose` | Abduction | Generate hypothesis → L0 |
| `/q1-add` | `quint_propose` | Abduction | Add user hypothesis → L0 |
| `/q2-verify` | `quint_verify` | Deduction | Logical check → L1 |
| `/q3-validate` | `quint_test` | Induction | Empirical test → L2 |
| `/q4-audit` | `quint_audit` | Audit | Calculate R_eff, check bias |
| `/q5-decide` | `quint_decide` | Decision | Create DRR (Design Rationale Record) |
| `/q-status` | `quint_status` | — | Show current phase |
| `/q-query` | Search `.quint/` | — | Query knowledge base |
| `/q-decay` | `quint_check_decay` | — | Check evidence freshness |
| `/q-actualize` | `quint_actualize` | — | Reconcile KB with code changes |

## Key Concepts to Understand

### Holons
Knowledge units stored in `.quint/`. Each has:
- **Layer**: L0/L1/L2/invalid
- **Kind**: `system` (code/architecture) or `episteme` (process/docs)
- **Scope (G)**: Conditions where claim applies
- **R_eff**: Trust score (0-1), calculated not estimated

### WLNK (Weakest Link Principle)
System reliability = min(component reliabilities), never average.

### Relations (Declared During Creation)
- **ComponentOf/ConstituentOf**: Created via `depends_on` in `quint_propose`
  - Affects R_eff: child ≤ parent
- **MemberOf**: Created via `decision_context` in `quint_propose`
  - Groups alternatives, doesn't affect R_eff

### Congruence Levels (CL)
Evidence transfer quality:
- **CL3**: Same context (internal test) — no penalty
- **CL2**: Similar context (related project) — small penalty
- **CL1**: Different context (external docs) — large penalty

### Transformer Mandate
**Critical**: Systems cannot transform themselves. AI generates options + evidence, human decides. Never make architectural decisions autonomously.

## Code Patterns

### Database Layer (`src/mcp/db/`)
- Uses **sqlc** for type-safe SQL
- Schema: `schema.sql`
- Queries: `query.sql` → generated Go code in `query.sql.go`
- Transactions: Use `Store.ExecTx` for atomic operations

### FSM Layer (`src/mcp/internal/fpf/fsm.go`)
- State machine enforcement
- Phase transitions must follow rules
- State persisted in `.quint/state.json`

### Tools Layer (`src/mcp/internal/fpf/tools.go`)
- Business logic for each MCP tool
- Precondition checks (see `preconditions.go`)
- File system operations (markdown files in `.quint/knowledge/`)

### Server Layer (`src/mcp/internal/fpf/server.go`)
- JSON-RPC 2.0 protocol handler
- stdio-based communication
- Tool invocation router

## Testing Patterns

Prefer: **Integration > Unit**

- E2E tests check full MCP tool flow
- Integration tests verify DB + FSM + Tools together
- Unit tests only for complex pure functions (e.g., assurance calculator)

Example test structure:
```go
func TestToolName(t *testing.T) {
    // Setup: temp DB, FSM state
    store := setupTestDB(t)
    fsm, _ := fpf.LoadState("test", store.GetRawDB())
    tools := fpf.NewTools(fsm, tempDir, store)
    
    // Act: call tool
    result, err := tools.SomeMethod(...)
    
    // Assert: check result and state changes
    require.NoError(t, err)
    assert.Contains(t, result, "expected")
    assert.Equal(t, fpf.PhaseExpected, fsm.State.Phase)
}
```

## Architecture Guidelines

### Functional Core, Imperative Shell
- **Core**: Pure business logic (FSM transitions, calculations, validations)
- **Shell**: I/O operations (DB, file system, stdio)
- Keep them separate — core never directly performs I/O

### Error Handling
- Return errors explicitly: `func DoThing() (string, error)`
- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Fail fast for programmer errors
- Handle gracefully for user/external errors

### File System Conventions
```
.quint/
├── knowledge/
│   ├── L0/          # Unverified hypotheses
│   ├── L1/          # Logically verified
│   ├── L2/          # Empirically validated
│   └── invalid/     # Falsified claims
├── evidence/        # Supporting data
├── decisions/       # DRRs
├── sessions/        # Reasoning session logs
└── quint.db         # SQLite database
```

## Common Code Review Checkpoints

When reviewing/generating code:

1. **Preconditions**: Does this tool check FSM phase?
2. **State Transition**: Is `fsm.State.Phase` updated correctly?
3. **Persistence**: Is state saved via `fsm.SaveState("default")`?
4. **Error Wrapping**: Are errors wrapped with context?
5. **Layer Placement**: Does the file go to correct L0/L1/L2 folder?
6. **WLNK**: If dependencies exist, is R_eff calculated correctly?
7. **Audit Trail**: Are tool calls logged for transparency?

## Command Examples

### Generate Hypothesis
```go
tools.ProposeHypothesis(
    title: "Use Redis for caching",
    content: "Implement Redis cache for user sessions",
    scope: "Auth service, >1000 req/s",
    kind: "system",
    rationale: `{"anomaly":"Slow login", "approach":"Redis", "alternatives_rejected":"Memcached"}`,
    decision_context: "caching-strategy",  // Groups alternatives
    depends_on: ["session-storage"],       // Dependency for WLNK
    dependency_cl: 3,                      // Same context
)
```

### Verify Logic
```go
tools.VerifyHypothesis(
    hypothesis_id: "h-001",
    checks_json: `{"preconditions":["Redis available"],"contradictions":[],"constraints":["Must persist >1hr"]}`,
    verdict: "PASS",  // PASS|FAIL|REFINE
)
```

### Validate Empirically
```go
tools.ManageEvidence(
    phase: PhaseInduction,
    operation: "add",
    hypothesis_id: "h-001",
    evidence_type: "internal",
    content: "Load test: 5000 req/s, 50ms p99",
    verdict: "PASS",
    assurance_level: "L2",
    source_agent: "test-runner",
    external_url: "",
)
```

## Integration Points

### For GitHub Copilot Chat
When user asks "How do I...":
1. Check if `.quint/` exists → suggest using FPF workflow
2. For decisions: recommend `/q1-hypothesize` → `/q2-verify` → `/q3-validate`
3. For queries: suggest `/q-query <search term>`
4. For stale evidence: suggest `/q-decay`

### For Code Completion
Context-aware suggestions:
- In `tools.go`: Suggest MCP tool function signatures
- In `schema.sql`: Suggest table schema patterns
- In `query.sql`: Suggest sqlc query patterns
- In tests: Suggest test setup boilerplate

## Useful References

- **FPF Theory**: See `docs/advanced.md`
- **Architecture**: See `docs/architecture.md`
- **Workflow Examples**: See `docs/workflow_example/`
- **Agent Instructions**: See `CLAUDE.md` (comprehensive reference)

## Critical Reminders

1. **Never auto-commit** — only commit when explicitly asked
2. **Follow existing patterns** — check neighboring files before introducing new patterns
3. **Test contracts, not internals** — test through public interfaces
4. **Keep functions focused** — extract complex logic into named functions
5. **Documentation in context** — prefer self-documenting code, comment only WHY not WHAT
6. **Transformer Mandate** — generate options, human decides

---

**When in doubt**: Check `CLAUDE.md` for comprehensive guidelines. This file is optimized for GitHub Copilot integration specifically.
