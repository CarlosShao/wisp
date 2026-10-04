param(
    [Parameter(Mandatory)][string]$Path
)
# Parser-only check (nothing executes): counts the syntax errors PowerShell
# finds in a script file and prints each with its line and column.
$errors = $null
$tokens = $null
[void][System.Management.Automation.Language.Parser]::ParseFile($Path, [ref]$tokens, [ref]$errors)
Write-Host ("parse-file={0} syntax-errors={1} tokens={2}" -f $Path, $errors.Count, $tokens.Count)
foreach ($e in $errors) {
    Write-Host ("  {0} at line {1} col {2}: {3}" -f $e.ErrorId, $e.Extent.StartLineNumber, $e.Extent.StartColumnNumber, $e.Message)
}
