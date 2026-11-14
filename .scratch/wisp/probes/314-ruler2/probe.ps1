# Ticket 314 post-audit repair: build/run the same probe file in two out-of-repo
# export trees - HEAD's ruler (before) and this working copy's ruler (after).
$ErrorActionPreference = 'Continue'
$repo = 'D:\work\workspace\projects plans\Wisp'
$root = 'C:\Temp\wisp-ruler2-20261114'

# Repair the "new" tree (the first copy attempt could not merge into a
# pre-created directory).
if (Test-Path "$root\new") { Remove-Item -Recurse -Force "$root\new" }
Expand-Archive -Path "$root\src.zip" -DestinationPath "$root\new" -Force
Copy-Item -Force "$repo\cmd\wisp\panel_profile_leaf_shape_windows_test.go" "$root\new\cmd\wisp\panel_profile_leaf_shape_windows_test.go"
Copy-Item -Force "$repo\.scratch\wisp\probes\314-ruler2\zzz_ruler2_probe_windows_test.go" "$root\new\cmd\wisp\zzz_ruler2_probe_windows_test.go"
Copy-Item -Force "$repo\.scratch\wisp\probes\314-ruler2\zzz_ruler2_probe_windows_test.go" "$root\old\cmd\wisp\zzz_ruler2_probe_windows_test.go"

# The export tree has no third_party (git-ignored): the DLLs must sit next to the
# test binary, which is why they are copied from the repo.
Copy-Item -Force "$repo\third_party\sherpa-onnx\*.dll" "$root\stage-old\"
Copy-Item -Force "$repo\third_party\sherpa-onnx\*.dll" "$root\stage-new\"

"OLD_SHA=$((Get-FileHash "$root\old\cmd\wisp\panel_profile_leaf_shape_windows_test.go" -Algorithm SHA256).Hash)"
"NEW_SHA=$((Get-FileHash "$root\new\cmd\wisp\panel_profile_leaf_shape_windows_test.go" -Algorithm SHA256).Hash)"
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
	"=== run $tree ==="
	& "$root\$stage\wisp.test.exe" '-test.run=^Test314' '-test.v' 2>&1
	"RUN_RC_$tree=$LASTEXITCODE"
	Pop-Location
}
"DONE"
