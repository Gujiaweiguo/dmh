---
description: Implement an approved OpenSpec change
agent: build
---

Implement the OpenSpec change following these steps:

1. **Read the proposal**:
   - Read `@/openspec/changes/<change-id>/proposal.md` to understand scope
   - Read `@/openspec/changes/<change-id>/design.md` if present
   - Read `@/openspec/changes/<change-id>/tasks.md` to get the checklist

2. **Implement sequentially**:
   - Work through tasks in `tasks.md` one by one
   - Keep edits minimal and focused on the requested change
   - Confirm completion before moving to the next task
   - Make sure every item in `tasks.md` is actually finished

3. **Update checklist**:
   - After all work is done, mark every task as `- [x]` in `tasks.md`
   - Ensure the checklist accurately reflects reality

4. **Reference tools**:
   - Use `openspec list` to see active changes
   - Use `openspec show <item>` for additional context when needed

**IMPORTANT**: Only implement changes for an approved proposal. Do not start implementation until the proposal has been reviewed and approved.

Reference `@/openspec/AGENTS.md` for additional conventions and clarifications.
