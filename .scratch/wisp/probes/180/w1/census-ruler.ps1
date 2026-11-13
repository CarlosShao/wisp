# Ticket 180 AC#4 census ruler (leg 180-w1).
# For every field name carrying a `default:` tag in internal/config/schema.go,
# count non-test production hits of `\.Name\b` under internal/ and cmd/,
# and dump every hit line for manual true-reader classification.
$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath (Join-Path $PSScriptRoot '..\..\..\..\..')

$names = @(
  'Language','Theme','Autostart','SingleInstance',
  'Size','Style','OpacityIdle','ClickThrough','HideOnFullscreen',
  'Cancel','WarmTimeoutSec','SettlingSec','ConversationIdleSec',
  'Enabled','VetoWords','Provider','Speed','EchoRef',
  'Punctuation','ConversationMode','InputDevice','SampleRate','HalfDuplex','MicMutedDefault',
  'Temperature','ThinkingIntensity','Max','BackoffMS','Loose','AllowMissingUsage',
  'Billing','TimeoutMS','RepeatThresholds','MaxRounds',
  'SteeringEnabled','ProjectInstructionsEnabled',
  'ConfirmTimeoutSec','L1WindowSec','ShellEnabled','AllowShellString','PermissionMode',
  'DeleteEnabled','Mode','BlockPrivateRanges',
  'RedactPaths','DiagnosticsOptIn','RetentionDays','KeepTranscript','KeepAudio',
  'L1Enabled','L1Max','L3RetentionDays',
  'Width','KeepAliveInSession','FontSize',
  'AlertThreshold','OverBudget','Tier2Enabled','VerifySignature',
  'SizeMB','Days','Level','SLOSampleIntervalSec'
)

$root = (Get-Location).Path
$files = Get-ChildItem -Path 'internal','cmd' -Recurse -Filter '*.go' |
  Where-Object { $_.Name -notmatch '_test\.go$' }
"files_scanned=$($files.Count)" | Write-Output

$summary = New-Object System.Collections.Generic.List[string]
$raw = New-Object System.Collections.Generic.List[string]

foreach ($n in $names) {
  $pattern = '\.' + $n + '\b'
  $hits = @($files | Select-String -Pattern $pattern)
  $summary.Add(("{0}`t{1}" -f $n, $hits.Count))
  $raw.Add(("=== {0} count={1}" -f $n, $hits.Count))
  foreach ($h in $hits) {
    $rel = $h.Path.Substring($root.Length + 1)
    $raw.Add(("{0}:{1}: {2}" -f $rel, $h.LineNumber, $h.Line.Trim()))
  }
}

$outDir = Join-Path $PSScriptRoot '.'
$summaryPath = Join-Path $outDir 'summary.tsv'
$rawPath = Join-Path $outDir 'raw-ruler.txt'
[System.IO.File]::WriteAllText($summaryPath, ($summary -join "`n") + "`n")
[System.IO.File]::WriteAllText($rawPath, ($raw -join "`n") + "`n")
"summary=$summaryPath"
"raw=$rawPath"
"names=$($names.Count)"
