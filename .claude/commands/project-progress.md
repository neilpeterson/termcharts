---
description: Display project progress as a pie chart
allowed-tools: Bash(grep:*), Bash(go run:*)
---

# Project Progress Chart

Display a colorized pie chart showing project completion progress.

## Steps:

1. **Count tasks**:
   - Count completed tasks: `grep -c '\[x\]' docs/status.md`
   - Count incomplete tasks: `grep -c '\[ \]' docs/status.md`

2. **Generate and display pie chart**:
   ```bash
   go run ./cmd/termcharts pie <completed> <incomplete> --labels "Completed,Remaining" --color --theme default
   ```

3. **Output the raw console result** - Do NOT wrap in markdown code blocks or try to reproduce the chart as text. Just let the Bash tool output display directly so colors are preserved.

## Important:

- Do NOT use markdown formatting around the chart output
- Do NOT copy/paste the chart into your response as text
- Simply run the Bash command and let the tool output speak for itself
- Add a brief one-line summary below: "X of Y tasks completed (Z%)"
