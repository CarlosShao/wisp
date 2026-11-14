# Ticket 314 post-audit repair: do the NEW resident controls have teeth?
# In an out-of-repo export tree (no probe file, so the real scan is green), break
# one matcher arm at a time and check the matching resident control goes red.
# The audit's complaint was exactly that breaking these arms left the table green.
$ErrorActionPreference = 'Continue'
$repo = 'D:\work\workspace\projects plans\Wisp'
$root = 'C:\Temp\wisp-ruler2-20261114'
$tree = "$root\teeth"
$stage = "$root\stage-teeth"
$rel = 'cmd\wisp\panel_profile_leaf_shape_windows_test.go'
$orig = Get-Content "$repo\$rel" -Raw -Encoding utf8

if (Test-Path $tree) { Remove-Item -Recurse -Force $tree }
Expand-Archive -Path "$root\src.zip" -DestinationPath $tree -Force
New-Item -ItemType Directory -Force $stage | Out-Null
Copy-Item -Force "$repo\$rel" "$tree\$rel"
Copy-Item -Force "$repo\third_party\sherpa-onnx\*.dll" "$stage\"

function Invoke-Teeth([string]$label, [string]$mutation, [string]$from, [string]$to) {
	$text = $orig
	if (-not $text.Contains($from)) {
		"=== $label === MUTATION_ANCHOR_MISSING"
		return
	}
	$text = $text.Replace($from, $to)
	Set-Content -Path "$tree\$rel" -Value $text -Encoding utf8 -NoNewline
	Push-Location $tree
	go test -c -o "$stage\wisp.test.exe" ./cmd/wisp 2>&1 | Out-Null
	$build = $LASTEXITCODE
	Pop-Location
	Push-Location "$tree\cmd\wisp"
	$out = & "$stage\wisp.test.exe" '-test.run=^Test314' '-test.v' 2>&1
	$run = $LASTEXITCODE
	Pop-Location
	$fails = $out | Where-Object { $_ -match 'FAIL: Test314PanelProfileShapeRulerRingsOnTheFixedNameShape/' }
	"=== $label === build_rc=$build run_rc=$run fails=$($fails.Count)"
	$fails | ForEach-Object { "    $_" }
	Copy-Item -Force "$repo\$rel" "$tree\$rel"
}

Invoke-Teeth -label 'T1 newPanelControl arm' `
	-mutation 'drop newPanelControl from the sink call set' `
	-from '"NewPanelManager", "newPanelControl":' `
	-to '"NewPanelManager":'

Invoke-Teeth -label 'T2 DataPath field-assignment arm' `
	-mutation 'make the selector-LHS arm never fire' `
	-from 'isSel && sel.Sel != nil && panelProfileSinkKey(sel.Sel.Name)' `
	-to 'isSel && sel.Sel != nil && false'

Invoke-Teeth -label 'T3 DataPath: key arm' `
	-mutation 'make the KeyValueExpr arm do nothing' `
	-from 'hits(expr.Value)' `
	-to '_ = expr.Value'

Invoke-Teeth -label 'T4 summary verdict' `
	-mutation 'count only unexempted shape hits again (the old line)' `
	-from 'clean := red == 0 && faults == 0' `
	-to 'clean := red == 0'

# Restore and prove the tree is green again.
Push-Location $tree
go test -c -o "$stage\wisp.test.exe" ./cmd/wisp 2>&1 | Out-Null
$build = $LASTEXITCODE
Pop-Location
Push-Location "$tree\cmd\wisp"
$out = & "$stage\wisp.test.exe" '-test.run=^Test314' '-test.v' 2>&1
$run = $LASTEXITCODE
Pop-Location
"=== RESTORED === build_rc=$build run_rc=$run"
$out | Where-Object { $_ -match '形状尺【|^(--- )?(FAIL|PASS): Test314' } | ForEach-Object { "    $_" }
"RESTORED_SHA=$((Get-FileHash "$tree\$rel" -Algorithm SHA256).Hash)"
"REPO_SHA=$((Get-FileHash "$repo\$rel" -Algorithm SHA256).Hash)"
"DONE"
