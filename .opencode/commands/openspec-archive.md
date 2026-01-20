---
description: Archive a deployed OpenSpec change
agent: build
---

Archive the OpenSpec change following these steps:

1. **Identify the change-id**:
   - If the prompt includes a specific change ID, use it
   - If referenced loosely, run `openspec list` to find candidates and confirm
   - Otherwise, run `openspec list` and ask the user which change to archive
   - Stop if you cannot identify a single change ID

2. **Validate the change**:
   - Run `openspec list` or `openspec show <id>` to confirm it exists
   - Stop if the change is missing, already archived, or not ready

3. **Archive the change**:
   - Run `openspec archive <id> --yes` to move and update specs
   - Use `--skip-specs` only for tooling-only work

4. **Review results**:
   - Confirm target specs were updated
   - Confirm change is in `openspec/changes/archive/`

5. **Validate**:
   - Run `openspec validate --strict --no-interactive`
   - Inspect with `openspec show <id>` if anything looks off

Reference `@/openspec/AGENTS.md` for additional conventions and clarifications.
