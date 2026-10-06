$ErrorActionPreference = "Stop"
$repo = "D:/work/workspace/projects plans/Wisp"
$src = Join-Path $repo "cmd/wisp/config_reload.go"
$dst = Join-Path $repo ".scratch/wisp/probes/231/r1/mutation/config_reload_no_branch.go"
$text = [IO.File]::ReadAllText($src)
# Cut exactly the newer-build case: from the case line through the line before the
# strings.HasPrefix(d, "config.toml:") case.
$startNeedle = "		case strings.Contains(d, ""was written by a newer build""):"
$endNeedle = "		case strings.HasPrefix(d, ""config.toml:""):"
$i = $text.IndexOf($startNeedle)
$j = $text.IndexOf($endNeedle)
if ($i -lt 0 -or $j -lt 0 -or $j -le $i) { throw "mutation needles not found: i=$i j=$j" }
$out = $text.Substring(0, $i) + $text.Substring($j)
if ($out.Contains("was written by a newer build")) { throw "mutation did not remove the branch" }
[IO.File]::WriteAllText($dst, $out)
Write-Output ("wrote " + $dst + " bytes=" + (Get-Item $dst).Length)
