---
description: Display current project status and recent changes
allowed-tools: Read, Glob, Bash(grep:*), Bash(go run:*)
---

# Project Status Overview

You've been asked to provide the current project status.

## Steps:

1. **Read status file**: Read `docs/status.md` to get the current project state

2. **Calculate and display progress chart**:
   - Count completed tasks: `grep -c '\[x\]' docs/status.md`
   - Count incomplete tasks: `grep -c '\[ \]' docs/status.md`
   - Generate a pie chart showing progress:
     ```bash
     go run ./cmd/termcharts pie <completed> <incomplete> --labels "Completed,Remaining" --color
     ```
   - Display the chart at the top of the status report

3. **Display organized summary**:
   - **Current Status**: What's currently IN PROGRESS
   - **Recent Changes**: Latest updates and completions
   - **Upcoming Tasks**: What's in the backlog
   - **Blockers**: Any items marked as blocked or needing attention

4. **Provide context**:
   - Highlight which tasks are high priority
   - Note any tasks that have been in progress for a while
   - Identify dependencies between tasks if relevant

5. **Ask next steps**:
   - Ask the user what they'd like to work on next
   - Offer to start a new task from the backlog
   - Suggest `/plan` if planning a new feature

## Output Format:

Present the information in a clear, scannable format:
- Start with the progress pie chart
- Use headers and bullet points
- Highlight important items
- Keep it concise but informative
- Make it easy to see what needs attention

## Important Notes:

- Always start by reading docs/status.md
- If status.md doesn't exist, let user know and offer to create it
- Don't make assumptions about task priority - check the file
- The pie chart visually shows project completion at a glance
