# Ticket 442 AC#2: the same compiling probe in two out-of-repo export trees.
#   old = HEAD's ruler  (the one ticket 442 says is blind to A1b/A1c)
#   new = this working copy's ruler (ticket 442's fix)
# Each tree gets an identical copy of .scratch/wisp/probes/442/zzz_442_probe_windows_test.go.
# The log is what docs/evidence/s1/442-*.md quotes; nothing here touches the repo.
$ErrorActionPreference = 'Continue'
$repo = 'D:\work\workspace\projects plans\Wisp'
$root = 'C:\Temp\wisp-442-20261114'
$probe = "$repo\.scratch\wisp\probes\442\zzz_442_probe_windows_test.go"
$ruler = "$repo\cmd\wisp\panel_profile_leaf_shape_windows_test.go"

if (Test-Path "$root\old") { Remove-Item -Recurse -Force "$root\old" }
if (Test-Path "$root\new") { Remove-Item -Recurse -Force "$root\new" }
New-Item -ItemType Directory -Force $root, "$root\stage-old", "$root\stage-new" | Out-Null
Push-Location $repo
git archive HEAD -o "$root\src.zip"
"ARCHIVE_RC=$LASTEXITCODE"
Pop-Location
Expand-Archive -Path "$root\src.zip" -DestinationPath "$root\old" -Force
Expand-Archive -Path "$root\src.zip" -DestinationPath "$root\new" -Force
Copy-Item -Force $ruler "$root\new\cmd\wisp\panel_profile_leaf_shape_windows_test.go"
Copy-Item -Force $probe "$root\old\cmd\wisp\zzz_442_probe_windows_test.go"
Copy-Item -Force $probe "$root\new\cmd\wisp\zzz_442_probe_windows_test.go"
Copy-Item -Force "$repo\third_party\sherpa-onnx\*.dll" "$root\stage-old\"
Copy-Item -Force "$repo\third_party\sherpa-onnx\*.dll" "$root\stage-new\"

"REPO_RULER_SHA=$((Get-FileHash $ruler -Algorithm SHA256).Hash)"
"OLD_RULER_SHA=$((Get-FileHash "$root\old\cmd\wisp\panel_profile_leaf_shape_windows_test.go" -Algorithm SHA256).Hash)"
"NEW_RULER_SHA=$((Get-FileHash "$root\new\cmd\wisp\panel_profile_leaf_shape_windows_test.go" -Algorithm SHA256).Hash)"
"PROBE_SHA=$((Get-FileHash $probe -Algorithm SHA256).Hash)"
"OLD_FILES=$((Get-ChildItem "$root\old\cmd\wisp\*.go").Count) NEW_FILES=$((Get-ChildItem "$root\new\cmd\wisp\*.go").Count)"

$pairs = @(
	[pscustomobject]@{ tree = 'old'; stage = 'stage-old' },
	[pscustomobject]@{ tree = 'new'; stage = 'stage-new' }
)
foreach ($pair in $pairs) {
	$tree = $pair.tree
	$stage = $pair.stage
	Push-Location "$root\$tree"
	"=== build $tree ==="
	go test -c -o "$root\$stage\wisp.test.exe" ./cmd/wisp 2>&1
	"BUILD_RC_$tree=$LASTEXITCODE"
	Pop-Location
	Push-Location "$root\$tree\cmd\wisp"
	"=== run $tree (probe present) ==="
	& "$root\$stage\wisp.test.exe" '-test.run=^Test314' '-test.v' 2>&1
	"RUN_RC_$tree=$LASTEXITCODE"
	Pop-Location
}

# Restore: the probe leaves the new tree, the ruler is untouched, so the run must
# be green again and the ruler's bytes must be the repo's.
"=== restore new (remove the probe) ==="
Remove-Item -Force "$root\new\cmd\wisp\zzz_442_probe_windows_test.go"
"RESTORED_RULER_SHA=$((Get-FileHash "$root\new\cmd\wisp\panel_profile_leaf_shape_windows_test.go" -Algorithm SHA256).Hash)"
Push-Location "$root\new"
"=== rebuild new (probe restored) ==="
go test -c -o "$root\stage-new\wisp.test.exe" ./cmd/wisp 2>&1
"BUILD_RC_restored=$LASTEXITCODE"
Pop-Location
Push-Location "$root\new\cmd\wisp"
"=== run new (probe restored) ==="
& "$root\stage-new\wisp.test.exe" '-test.run=^Test314' '-test.v' 2>&1
"RUN_RC_restored=$LASTEXITCODE"
Pop-Location
"DoNE"
