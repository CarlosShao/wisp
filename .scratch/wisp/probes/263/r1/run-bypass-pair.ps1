# Bypass demonstration for ticket 263 AC#4, second half.
#
# One PowerShell session runs TWO mutations of the gate back to back:
#   1. clean-copy        -> exits 0, and in doing so LEAVES an exit code behind
#                           in this session's automatic variable.
#   2. plant-global-used -> reads $global:LASTEXITCODE (a spelling the word nail
#                           does not match) where the state's code should come
#                           from the sampler.
# Expected and measured: the nail stays green and the second run reports a
# normal-looking exit code it never earned, i.e. the bypass does not merely
# avoid the nail, it can borrow another program's number.
$ErrorActionPreference = 'Continue'
$here = $PSScriptRoot
& (Join-Path $here 'run-nail-mutations.ps1') -Subset smoke -SecondsPerState 1 -Only clean-copy
& (Join-Path $here 'run-nail-mutations.ps1') -Subset smoke -SecondsPerState 1 -Only plant-global-used
