# Initial Prompt: QR Frontend Hardening

**Copy and paste this into a new Cline agent instance:**

---

I need you to work on hardening the frontend QR scan functionality for non-OIDC interactive registration in the Headscale project.

**Your working directory is:**
```
/home/denny/Project/headscale-project/headscale-qr-frontend-harden
```

**Your branch is:** `feature/qr-frontend-harden`

**Your task:**
Harden the QR scan UI, error handling, and user selection in Headplane. Add comprehensive tests for edge cases including expired QR codes, invalid payloads, scan errors, and user selection validation.

**Instructions:**
Read the complete agent instructions at:
```
/home/denny/Project/headscale-project/headscale/AGENT_INSTRUCTIONS_QR_FRONTEND_HARDEN.md
```

**Start by:**
1. Changing your working directory to `/home/denny/Project/headscale-project/headscale-qr-frontend-harden`
2. Reading the agent instructions file
3. Reading the current implementation in `headplane/app/routes/machines/scan-qr.tsx`
4. Identifying edge cases not covered by tests
5. Adding tests and implementing hardening

**Important project rules:**
- Read `AGENTS.md` for project interaction rules
- Follow the existing code style and patterns
- Add tests before implementing changes (TDD)
- Use existing error/toast patterns
- Commit incrementally with clear messages

**When complete:**
Report which files were changed, which tests were added, and the results of running:
- `cd headplane && pnpm test:component`
- `cd headplane && pnpm exec oxlint app/routes/machines/scan-qr.tsx`
- `cd headplane && pnpm run build`

Begin by reading the agent instructions and confirming your understanding of the task.
