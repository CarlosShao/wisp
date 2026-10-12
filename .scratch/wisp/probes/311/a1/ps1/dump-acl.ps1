# Probe 311-a1: dump one directory's DACL, ACE by ACE (read-only).
# Usage: powershell -NoProfile -File dump-acl.ps1 [-Sddl] <path> [<path> ...]
$ErrorActionPreference = 'Continue'

Write-Output ('now=' + (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'))
Write-Output ('FileSystemRightsEnumNames=' + ([Enum]::GetNames([System.Security.AccessControl.FileSystemRights]) -join ','))
Write-Output ('Delete=' + [int][System.Security.AccessControl.FileSystemRights]::Delete)
foreach ($n in @('ListDirectory', 'Traverse', 'CreateFiles', 'CreateDirectories', 'Delete', 'ReadPermissions', 'ChangePermissions', 'TakeOwnership', 'Synchronize', 'FullControl', 'Modify')) {
    try { Write-Output ('maskbit ' + $n + '=' + [int]([System.Security.AccessControl.FileSystemRights]::$n)) }
    catch { Write-Output ('maskbit ' + $n + '=UNNAMED') }
}
# the .NET name holding FILE_DELETE_CHILD (0x40) is looked up by value, one name per line
foreach ($n in @('DeleteSubdirectoriesAndFiles', 'Write', 'Read', 'ReadAndExecute', 'Modify', 'FullControl', 'ExecuteFile', 'ListDirectory')) {
    Write-Output ('maskbit ' + $n + '=' + [int]([System.Security.AccessControl.FileSystemRights]::$n))
}
Write-Output ('Modify_has_0x40=' + (([int][System.Security.AccessControl.FileSystemRights]::Modify -band 64) -ne 0))
Write-Output ('FullControl_has_0x40=' + (([int][System.Security.AccessControl.FileSystemRights]::FullControl -band 64) -ne 0))
Write-Output ('0x1301BF_has_0x40=' + ((0x1301BF -band 64) -ne 0))
Write-Output ('0x1F01FF_has_0x40=' + ((0x1F01FF -band 64) -ne 0))

foreach ($p in $args) {
    Write-Output ('===== PATH=' + $p)
    if (-not (Test-Path -LiteralPath $p)) {
        Write-Output ('exists=NO')
        continue
    }
    $acl = Get-Acl -LiteralPath $p -ErrorAction SilentlyContinue
    if ($null -eq $acl) {
        Write-Output ('getacl_failed=1')
        continue
    }
    Write-Output ('owner=' + $acl.Owner)
    Write-Output ('dacl_protected_flag=' + $acl.AreAccessRulesProtected)
    Write-Output ('sddl=' + $acl.Sddl)
    $i = 0
    foreach ($r in $acl.Access) {
        $sid = ''
        try { $sid = $r.IdentityReference.Translate([System.Security.Principal.SecurityIdentifier]).Value }
        catch { $sid = 'UNRESOLVED(' + $r.IdentityReference.Value + ')' }
        $mask = ''
        try { $mask = [int]$r.FileSystemRights } catch { $mask = 'UNCASTABLE' }
        Write-Output ('ACE[' + $i + '] type=' + $r.AccessControlType + ' sid=' + $sid + ' who=' + $r.IdentityReference.Value +
            ' rights=' + $r.FileSystemRights + ' mask=0x' + ('{0:X}' -f $mask) +
            ' inh=' + $r.InheritanceFlags + ' prop=' + $r.PropagationFlags + ' isinherited=' + $r.IsInherited)
        $i = $i + 1
    }
    Write-Output ('ace_count=' + $i)
    Write-Output ('--- icacls_verbatim:')
    & $env:WINDIR\System32\icacls.exe $p 2>&1 | ForEach-Object { Write-Output ('ICACLS| ' + $_) }
}
Write-Output ('ps_errors=' + ($Error.Count))
