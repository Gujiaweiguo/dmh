---
description: Create a new OpenSpec change proposal
agent: plan
---

Create a new OpenSpec change proposal following these steps:

1. **Review current state**:
   - Run `openspec list` to see active changes
   - Run `openspec list --specs` to see existing specs
   - Read `@/openspec/project.md` for project context
   - Search existing requirements with `rg -n "Requirement:|Scenario:" openspec/specs`

2. **Choose a unique change-id**:
   - Use kebab-case, verb-led (add-, update-, remove-, refactor-)
   - Ensure uniqueness by checking existing changes

3. **Scaffold the proposal**:
   - Create `openspec/changes/<change-id>/proposal.md` with:
     * Why (problem/opportunity)
     * What changes (bullet list, mark **BREAKING** changes)
     * Impact (affected specs, affected code)
   - Create `openspec/changes/<change-id>/tasks.md` with ordered checklist
   - Create `openspec/changes/<change-id>/design.md` ONLY if needed:
     * Cross-cutting change or new architecture pattern
     * New external dependency or significant data model changes
     * Security, performance, or migration complexity
   - Create spec deltas in `openspec/changes/<change-id>/specs/<capability>/spec.md`
     * Use `## ADDED|MODIFIED|REMOVED Requirements`
     * Include at least one `#### Scenario:` per requirement
     * Reference related capabilities when relevant

4. **Validate strictly**:
   - Run `openspec validate <change-id> --strict --no-interactive`
   - Resolve every validation issue before presenting the proposal

**IMPORTANT**: Do NOT write any code during the proposal stage. Only create design documents. Implementation happens in /openspec-apply after approval.

Reference `@/openspec/AGENTS.md` for additional conventions and clarifications.
