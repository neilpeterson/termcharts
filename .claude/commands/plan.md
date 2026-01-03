---
description: Plan a new feature or task with comprehensive analysis
allowed-tools: Read, Glob, Grep, Edit, Bash(git:*)
---

# Comprehensive Feature Planning

You've been asked to plan a new feature or task. This requires thorough analysis and structured breakdown.

## Planning Process:

### 1. **Understand the Request**
   - Ask the user to describe what they want to build
   - Clarify scope, goals, and success criteria
   - Identify any constraints (time, compatibility, dependencies)
   - Ask about target users and use cases

### 2. **Gather Context**
   - Read `docs/status.md` to understand current project state
   - Search the codebase for related existing functionality
   - Identify files/modules that will be affected
   - Review architecture from CLAUDE.md or README
   - Check for similar patterns already implemented

### 3. **Technical Analysis**
   - **Architecture**: Determine where this fits in the codebase structure
   - **Dependencies**: Identify what this feature depends on and what depends on it
   - **Interfaces**: Define public APIs, data structures, or contracts
   - **Integration points**: How does this connect with existing code?
   - **Testing strategy**: What needs to be tested and how?

### 4. **Risk Assessment**
   Identify potential risks:
   - Breaking changes to existing functionality
   - Performance implications
   - Security considerations
   - Complexity hotspots
   - External dependencies or API limitations

### 5. **Break Down Into Subtasks**
   Create granular, actionable tasks:
   - **Aim for <1 hour per task** when possible
   - **Order by dependencies**: tasks that must happen first come first
   - **Mark complexity**: simple/moderate/complex
   - **Identify parallel work**: what can be done concurrently?

   For each subtask include:
   - Clear description of what needs to be done
   - Files that will be modified or created
   - Estimated complexity (simple/moderate/complex)
   - Dependencies on other subtasks
   - Acceptance criteria (how to verify it's done)

### 6. **Create Implementation Plan**
   Present a structured plan with:

   **Feature Overview:**
   - What: Brief description
   - Why: Purpose and benefits
   - Who: Target users or use cases

   **Architecture Changes:**
   - New files to create
   - Existing files to modify
   - New packages or dependencies needed
   - Data structures and interfaces

   **Task Breakdown:**
   - Phase 1: Foundation (data structures, core logic)
   - Phase 2: Implementation (main functionality)
   - Phase 3: Integration (wire it up, CLI commands if needed)
   - Phase 4: Testing & Polish (tests, examples, docs)

   **Acceptance Criteria:**
   - How to verify the feature works
   - Edge cases to handle
   - Performance benchmarks if applicable

   **Risks & Mitigations:**
   - Potential issues and how to handle them

### 7. **Get Confirmation**
   - Present the plan clearly with all sections
   - Ask if the scope and approach make sense
   - Offer to adjust based on feedback
   - Discuss any alternative approaches considered

### 8. **Update Status**
   Once confirmed:
   - Add feature to `docs/status.md` under appropriate section
   - Add all subtasks to the backlog
   - Mark the first task as the next item to work on
   - Commit the updated status.md if appropriate

### 9. **Ready to Start**
   - Summarize what will happen first
   - Ask if the user wants to start immediately
   - If yes, mark first task as IN PROGRESS and begin

## Planning Depth Options:

Adjust depth based on task complexity:

- **Quick Plan** (simple tasks): Steps 1, 5, 7, 8
- **Standard Plan** (most features): Steps 1-5, 7-9
- **Deep Plan** (complex features): All steps with detailed analysis

## Output Format:

Use clear, structured markdown:
```
# Feature: [Name]

## Overview
[What, why, who]

## Architecture
- Files to create: ...
- Files to modify: ...
- Dependencies: ...

## Task Breakdown

### Phase 1: Foundation
- [ ] Task 1 (simple) - description | Files: x.go, y.go | Depends: none
- [ ] Task 2 (moderate) - description | Files: z.go | Depends: Task 1

### Phase 2: Implementation
...

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Risks
- Risk 1 → Mitigation
- Risk 2 → Mitigation
```

## Best Practices:

- **Be thorough but concise**: Provide enough detail without overwhelming
- **Think about edge cases**: Empty data, large datasets, invalid input
- **Consider the user**: DevOps engineers want simplicity and reliability
- **Follow existing patterns**: Match the codebase style and architecture
- **Plan for testing**: Every feature needs tests
- **Documentation**: Note what docs need updating (README, godoc, examples)

## When to Use:

- Planning new chart types
- Adding CLI commands or data sources
- Significant refactoring
- New library features
- Performance optimizations
- Breaking changes

## Integration with Workflow:

This command integrates with your workflow:
- Updates `docs/status.md` automatically
- Creates tasks ready for `/done` to mark complete
- Works with `/project-status` to show planning progress
- Follows architecture defined in CLAUDE.md