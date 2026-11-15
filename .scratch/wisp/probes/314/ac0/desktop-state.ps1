# Desktop-state proxy for ticket 314 candidate (a): is the session locked/idle?
# input_desktop name ("Default" = unlocked interactive desktop; "Winlogon"/other = locked),
# LogonUI process count, and the current foreground window. Read-only.
$src = @'
using System;
using System.Runtime.InteropServices;
using System.Text;
public static class Desk314 {
  [DllImport("user32.dll", SetLastError=true)] public static extern IntPtr OpenInputDesktop(uint dwFlags, bool fInherit, uint dwDesiredAccess);
  [DllImport("user32.dll", SetLastError=true)] public static extern bool CloseDesktop(IntPtr hDesktop);
  [DllImport("user32.dll")] public static extern bool GetUserObjectInformation(IntPtr hObj, int nIndex, StringBuilder pvInfo, uint nLength, out uint lpnLengthNeeded);
  [DllImport("user32.dll")] public static extern IntPtr GetForegroundWindow();
  [DllImport("user32.dll", CharSet=CharSet.Unicode)] public static extern int GetWindowText(IntPtr hWnd, StringBuilder lpString, int nMaxCount);
  [DllImport("user32.dll")] public static extern int GetWindowThreadProcessId(IntPtr hWnd, out int pid);
  public static string InputDesktopName() {
    IntPtr h = OpenInputDesktop(0, false, 0x0001);
    if (h == IntPtr.Zero) return "ERR:" + Marshal.GetLastWin32Error();
    try {
      var sb = new StringBuilder(256); uint need;
      if (!GetUserObjectInformation(h, 2, sb, 256, out need)) return "ERR:" + Marshal.GetLastWin32Error();
      return sb.ToString();
    } finally { CloseDesktop(h); }
  }
  public static string Foreground() {
    IntPtr h = GetForegroundWindow();
    if (h == IntPtr.Zero) return "none";
    var sb = new StringBuilder(256); GetWindowText(h, sb, 256);
    int pid; GetWindowThreadProcessId(h, out pid);
    return string.Format("0x{0:X} pid={1} title={2}", h.ToInt64(), pid, sb.ToString());
  }
}
'@
Add-Type -TypeDefinition $src -ErrorAction Stop
"ts=$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
"input_desktop=$( [Desk314]::InputDesktopName() )"
"logonui_s1=$((Get-Process LogonUI -ErrorAction SilentlyContinue | Where-Object { $_.SessionId -eq 1 } | Measure-Object).Count)"
"foreground=$( [Desk314]::Foreground() )"
"cpu_total_pct=$((Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average)"
"msedgewebview2=$((Get-Process msedgewebview2 -ErrorAction SilentlyContinue | Measure-Object).Count)"
"wisp_processes=$((Get-Process wisp -ErrorAction SilentlyContinue | Measure-Object).Count)"
"wisp_test_processes=$((Get-Process wisp-winlive -ErrorAction SilentlyContinue | Measure-Object).Count)"
