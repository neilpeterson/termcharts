---
description: Mark current task as complete and update project status
allowed-tools: Read, Edit, Bash(git:*)
---

# Task Completion Workflow

You've been asked to mark a task as complete and update the project status.

## Steps to Complete:

1. **Confirm completion**: Ask the user what task was just completed (or infer from recent conversation)

2. **Update status.md**:
   - Read `docs/status.md`
   - Move the completed task from "IN PROGRESS" to "COMPLETED"
   - Add completion notes with what was accomplished
   - Update the "Recent Changes" section with a summary
   - Add any newly discovered tasks to the backlog

3. **Summarize**:
   - Provide a brief summary of what was accomplished
   - Highlight any important outcomes or learnings
   - List any new tasks that were added to the backlog

4. **Suggest next steps**:
   - Review remaining tasks in status.md
   - Recommend what to work on next based on priority
   - Ask the user to confirm the next task

## Important Notes:

- Always read status.md first to understand current state
- Be specific in completion notes (what was done, any issues encountered)
- Keep "Recent Changes" concise (1-2 sentence summary)
- Don't mark tasks complete if there are known issues or incomplete work
