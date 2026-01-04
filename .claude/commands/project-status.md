---
description: Display current project status and recent changes
allowed-tools: Read, Glob
---

# Project Status Overview

You've been asked to provide the current project status.

## Steps:

1. **Read status file**: Read `docs/status.md` to get the current project state

2. **Display organized summary**:
   - **Current Status**: What's currently IN PROGRESS
   - **Recent Changes**: Latest updates and completions
   - **Upcoming Tasks**: What's in the backlog
   - **Blockers**: Any items marked as blocked or needing attention

3. **Provide context**:
   - Highlight which tasks are high priority
   - Note any tasks that have been in progress for a while
   - Identify dependencies between tasks if relevant

4. **Ask next steps**:
   - Ask the user what they'd like to work on next
   - Offer to start a new task from the backlog
   - Suggest `/plan` if planning a new feature

## Output Format:

Present the information in a clear, scannable format:
- Use headers and bullet points
- Highlight important items
- Keep it concise but informative
- Make it easy to see what needs attention
- Mention `/project-progress` for visual progress chart

## Important Notes:

- Always start by reading docs/status.md
- If status.md doesn't exist, let user know and offer to create it
- Don't make assumptions about task priority - check the file
