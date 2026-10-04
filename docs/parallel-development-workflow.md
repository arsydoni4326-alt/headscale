# Parallel Development Workflow with Git Worktrees

Reusable workflow for executing parallel development tasks using Git Worktrees and multiple Cline agent instances.

## Overview

Execute multiple independent tasks simultaneously in isolated Git worktrees, each with its own feature branch.

## When to Use

**Use when:**
- Phase contains multiple independent tasks
- Minimal file overlap between tasks
- Tasks can be validated independently
- Up to 10 tasks can run safely

**Do NOT use when:**
- Tasks modify the same core files
- Sequential implementation required
- Tightly coupled tasks
- Database migrations (unless coordinated)

## Workflow

### 1. Planning (Required)
- Read ROADMAP.md and documentation
- Create implementation plan
- Identify parallelizable tasks (max 10)
- Document dependencies
- **GET APPROVAL** before proceeding

### 2. Worktree Creation

```bash
# Verify clean state
git status --porcelain
git branch --show-current
git rev-parse HEAD  # Record SHA

# Create worktrees (one directory above project)
cd /path/to/parent
git -C project worktree add -b feature/task-name ../project-task-name base-branch

# Verify
git worktree list
```

### 3. Agent Instructions

Create:
- `AGENT_INSTRUCTIONS.md` - Overview and general rules
- `AGENT_INSTRUCTIONS_TASK1.md` - Task-specific instructions
- `INITIAL_PROMPT.md` - Initial prompts for each agent

Each task instruction must include:
- Objective, worktree path, branch
- Scope (what to change and NOT change)
- Step-by-step instructions
- Validation checklist
- Exact commit message format
- Dependencies

### 4. Agent Execution

**Sequential start if dependencies:**
1. Primary agent completes first
2. Reports commit SHA
3. Dependent agents merge and start in parallel

**Isolation rules:**
- Work only in assigned worktree
- Never modify other worktrees
- Never merge to main yourself
- Stay within assigned scope

### 5. Integration Analysis

After all agents complete:
```bash
# Check for conflicts
git diff feature/task1..feature/task2 --name-only

# Identify shared files
git log feature/task1 --name-only --pretty=format: | sort -u > task1.txt
git log feature/task2 --name-only --pretty=format: | sort -u > task2.txt
comm -12 task1.txt task2.txt
```

### 6. Validation

Per-branch:
```bash
cd /worktree
make build && make test
```

### 7. Merge

Merge in dependency order:
```bash
git checkout main
git merge --no-ff feature/task-name -m "Description"
make clean && make build && make test
```

### 8. Cleanup

Only after approval:
```bash
git worktree remove /path/to/worktree
git branch -d feature/task-name  # Optional
```

## Directory Structure

```
parent/
├── project/           # Main worktree
├── project-task1/     # Task 1 worktree
├── project-task2/     # Task 2 worktree
```

## Branch Naming

```
feature/<short-descriptive-name>
```

## Common Issues

**Merge conflicts:** Identify files, determine precedence, request guidance
**Dependency not ready:** Wait, do not implement workaround
**Scope creep:** Stop, document, ask for approval
**Build failures:** Report exact error, do not patch

## Phase 14 Example

5 worktrees: core module rewrite (primary) + 4 parallel tasks (submodule, build, docs, nix).
See `AGENT_INSTRUCTIONS*.md` for complete example.

## References

- `AGENTS.md` - Agent behavior
- `ROADMAP.md` - Project roadmap
- `git help worktree`
