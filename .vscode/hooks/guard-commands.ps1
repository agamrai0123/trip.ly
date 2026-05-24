# .vscode/hooks/guard-commands.ps1
# Called by the PreToolUse hook before Copilot runs any shell command on Windows.
# Exit 1 to BLOCK the command; exit 0 to ALLOW it.

param([string]$Command)

$BlockedPatterns = @(
    "DROP DATABASE",
    "DROP SCHEMA",
    "dropdb",
    "rm -rf /",
    "rm -rf ~",
    "git push --force",
    "git push -f",
    "git push origin main",
    "git push origin master",
    "curl.*\| bash",
    "wget.*\| bash",
    "chmod 777",
    "truncate.*wanderplan"
)

foreach ($pattern in $BlockedPatterns) {
    if ($Command -match $pattern) {
        Write-Host "BLOCKED by WanderPlan hook: command matches forbidden pattern `"$pattern`""
        Write-Host "  Command: $Command"
        Write-Host "  If this is intentional, run it manually in the terminal."
        exit 1
    }
}

exit 0
