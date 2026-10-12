# Probe 311-a1 substitute ruler (dai-yong-chi): plant the fixture's own SDDL shapes
# onto directories THIS script creates inside the two dedicated roots, so a
# "after seal" DACL can be read on BOTH volumes. Uses Set-Acl (Win32
# SetNamedSecurityInfo), i.e. NOT wisp's code path - the label goes in the log.
# Creates directories; deletes nothing.
# Usage: powershell -NoProfile -File plant-probe.ps1 <rootA> <rootB>
$ErrorActionPreference = 'Continue'
$runner = 'D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)'
$oneLa  = 'D:(A;;FA;;;LA)'

function Plant($dir, $sddl) {
    if (Test-Path -LiteralPath $dir) {
        Write-Output ('SKIP_EXISTS=' + $dir)
        return
    }
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
    Write-Output ('CREATED=' + $dir)
    $before = $Error.Count
    $ds = New-Object System.Security.AccessControl.DirectorySecurity
    $ds.SetSecurityDescriptorSddlForm($sddl, [System.Security.AccessControl.AccessControlSections]::Access)
    Set-Acl -LiteralPath $dir -AclObject $ds
    Write-Output ('plant_threw=' + ($Error.Count -ne $before))
}

foreach ($r in $args) {
    Write-Output ('===== ROOT=' + $r)
    Plant (Join-Path $r 'plant2-runner') $runner
    Plant (Join-Path $r 'plant2-la') $oneLa
}
Write-Output ('ps_errors=' + ($Error.Count))
