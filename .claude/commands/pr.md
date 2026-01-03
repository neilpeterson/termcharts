---
description: Create PR, watch tests, merge if passing, clean up
allowed-tools: Bash(git:*), Bash(gh:*), Bash(make:*), Read, Grep
argument-hint: [pr-title]
---

# Pull Request Creation and Management

Create a pull request, watch CI tests, and automatically merge if all tests pass.

## Workflow:

### 1. **Pre-flight Checks**
   - Verify you're on a feature branch (not main)
   - Check for uncommitted changes
   - Run local tests with `make test` to catch issues early
   - If local tests fail, STOP and report the errors to the user

### 2. **Review Changes**
   - Run `git status` to see current changes
   - Run `git diff --staged` and `git diff` to review changes
   - Run `git log main..HEAD --oneline` to see commits
   - Summarize the changes for the user

### 3. **Commit Changes (if needed)**
   - If there are uncommitted changes, stage and commit them
   - Follow the commit message format from CLAUDE.md:
     - Clear summary of changes
     - Include "🤖 Generated with [Claude Code](https://claude.com/claude-code)"
     - Include "Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

### 4. **Push to Remote**
   - Push the branch to origin with `git push` or `git push -u origin <branch>`
   - Verify push succeeded

### 5. **Create Pull Request**
   - Use `gh pr create` with a clear title and body
   - Title: Use argument if provided, otherwise infer from branch name or commits
   - Body should include:
     - ## Summary (bullet points of changes)
     - ## Test plan (checklist of what was tested)
     - "🤖 Generated with [Claude Code](https://claude.com/claude-code)"
   - Use HEREDOC for proper formatting

### 6. **Watch CI Tests**
   - Get the PR number from the create command
   - Run `gh pr checks <pr-number> --watch` to watch CI tests
   - Wait for all checks to complete
   - Parse the output to determine if all tests passed

### 7. **Decision Point: Merge or Not**

   **IF ALL TESTS PASS:**
   - Merge the PR with `gh pr merge <pr-number> --squash --delete-branch`
   - Switch to main branch with `git checkout main`
   - Delete local branch with `git branch -d <branch-name>`
   - Pull latest changes with `git pull`
   - Report success to user with PR URL

   **IF ANY TESTS FAIL:**
   - DO NOT merge the PR
   - DO NOT delete the local branch
   - Report which tests failed to the user
   - Provide the URLs to failed CI jobs
   - Ask the user if they want to:
     - Fix the issues and push new commits
     - Close the PR
     - Merge anyway (not recommended)

### 8. **Error Handling**

   - If any step fails, stop and report the error
   - Never force push unless explicitly requested
   - Never merge with failing tests unless user explicitly overrides
   - Keep the branch if tests fail so user can fix and re-push

## Important Notes:

- **Local tests MUST pass before creating PR**
- **CI tests MUST pass before merging** (unless user explicitly overrides)
- **Always watch tests** - don't blindly merge
- **Preserve branch on failure** - user needs it to fix issues
- **Use squash merge** - keeps main branch clean
- **Delete remote branch** - but only if merge succeeds
- **Clean up local branch** - but only if merge succeeds

## Example Usage:

```bash
# With PR title
/pr "Add new feature X"

# Without title (infers from commits)
/pr
```

## Integration with Workflow:

This command ensures:
- All tests pass before merging
- Main branch stays clean (squash merge)
- Branches are cleaned up on success
- Branches are preserved on failure for debugging
- User is informed of all steps and outcomes
