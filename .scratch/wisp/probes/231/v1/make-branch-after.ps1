$ErrorActionPreference = "Stop"
$repo = "D:/work/workspace/projects plans/Wisp"
$src = Join-Path $repo "cmd/wisp/config_reload.go"
$dst = Join-Path $repo ".scratch/wisp/probes/231/v1/mutation/config_reload_branch_after_prefix.go"
$lines = [IO.File]::ReadAllLines($src)
$i = -1; $j = -1
for ($k = 0; $k -lt $lines.Length; $k++) {
    if ($i -lt 0 -and $lines[$k].Contains('case strings.Contains(d, "was written by a newer build"):')) { $i = $k }
    elseif ($i -ge 0 -and $lines[$k].Contains('case strings.HasPrefix(d, "config.toml:"):')) { $j = $k; break }
}
if ($i -lt 0 -or $j -lt 0) { throw "needles not found: i=$i j=$j" }
$block = $lines[$i..($j - 1)]
# rest = file without the newer-build block
$rest = @()
if ($i -gt 0) { $rest += $lines[0..($i - 1)] }
$rest += $lines[$j..($lines.Length - 1)]
# locate the prefix case inside $rest, then step over its 3-line return body
$jj = -1
for ($k = 0; $k -lt $rest.Length; $k++) {
    if ($rest[$k].Contains('case strings.HasPrefix(d, "config.toml:"):')) { $jj = $k; break }
}
if ($jj -lt 0) { throw "prefix case missing after removal" }
if (-not $rest[$jj + 1].Contains('cause=invalid')) { throw "unexpected shape after prefix case: " + $rest[$jj + 1] }
$out = @()
$out += $rest[0..($jj + 3)]
$out += $block
$out += $rest[($jj + 4)..($rest.Length - 1)]
$joined = [string]::Join("`n", $out)
# self-check: the newer-build case must now sit AFTER the prefix case
$p = $joined.IndexOf('case strings.HasPrefix(d, "config.toml:"):')
$n = $joined.IndexOf('case strings.Contains(d, "was written by a newer build"):')
if ($n -lt $p) { throw "mutation failed: branch still before the prefix case (n=$n p=$p)" }
if ($p -lt 0 -or $n -lt 0) { throw "mutation lost a case (n=$n p=$p)" }
[IO.File]::WriteAllText($dst, $joined)
Write-Output ("wrote " + $dst + " lines=" + $out.Length + " bytes=" + (Get-Item $dst).Length)
