$p = Get-Process -Id 24568 -ErrorAction SilentlyContinue
if ($null -eq $p) { 'alive=no'; exit 0 }
"alive=yes"
"main_window_handle=$($p.MainWindowHandle)"
"main_window_title=[$($p.MainWindowTitle)]"
"windows_visible=$(([string]$p.MainWindowHandle) -ne '0')"
$sig = @'
[DllImport("user32.dll")] public static extern bool EnumWindows(EnumWindowsProc lpEnumFunc, IntPtr lParam);
public delegate bool EnumWindowsProc(IntPtr hWnd, IntPtr lParam);
[DllImport("user32.dll")] public static extern uint GetWindowThreadProcessId(IntPtr hWnd, out uint pid);
[DllImport("user32.dll")] public static extern int GetWindowTextLength(IntPtr hWnd);
[DllImport("user32.dll", CharSet=CharSet.Unicode)] public static extern int GetWindowText(IntPtr hWnd, System.Text.StringBuilder s, int n);
[DllImport("user32.dll", CharSet=CharSet.Unicode)] public static extern int GetClassName(IntPtr hWnd, System.Text.StringBuilder s, int n);
[DllImport("user32.dll")] public static extern bool IsWindowVisible(IntPtr hWnd);
[DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr hWnd, out RECT r);
public struct RECT { public int Left; public int Top; public int Right; public int Bottom; }
'@
Add-Type -MemberDefinition $sig -Name W -Namespace U
$targets = New-Object System.Collections.ArrayList
$cb = [U+W+EnumWindowsProc]{
    param($h, $l)
    $pid2 = 0
    [U+W]::GetWindowThreadProcessId($h, [ref]$pid2) | Out-Null
    if ($pid2 -eq 24568) {
        $sb = New-Object System.Text.StringBuilder 256
        [U+W]::GetWindowText($h, $sb, 256) | Out-Null
        $cn = New-Object System.Text.StringBuilder 256
        [U+W]::GetClassName($h, $cn, 256) | Out-Null
        $r = New-Object U+W+RECT
        [U+W]::GetWindowRect($h, [ref]$r) | Out-Null
        [void]$targets.Add("hwnd=$h visible=$([U+W]::IsWindowVisible($h)) title=[$($sb.ToString())] class=[$($cn.ToString())] rect=$($r.Left),$($r.Top),$($r.Right),$($r.Bottom)")
    }
    return $true
}
[U+W]::EnumWindows($cb, [IntPtr]::Zero) | Out-Null
"window_count=$($targets.Count)"
$targets | ForEach-Object { $_ }
