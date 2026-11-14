#!/usr/bin/env pwsh
#Requires -Version 7.0
<#
.SYNOPSIS
  One real-machine session, captured once, for tickets 290 / 291 / 449 / 247.

.DESCRIPTION
  This is the machine half of .scratch/wisp/probes/runbook/real-machine-once.md.
  The human half is: press the mute hot key when asked, speak when asked, stay
  quiet when asked. Everything else is collected here.

  Hard rules this script obeys (they come from the dispatch, not from taste):

    * It never starts, stops, restarts or signals a product process. The ball is
      already running; it is READ (pid from run-marker.json, boot records from
      the persistent log sink), never touched. No Stop-Process, no Start-Process
      of wisp.exe, no hotkey synthesis.
    * It never writes config, thresholds, golden files, tickets or product code.
      The only writes are inside its own run directory.
    * It is idempotent and never overwrites: every run gets a fresh directory
      under <OutRoot>\<yyyyMMdd>\<label>[-HHmmss]\ and refuses to write a file
      that already exists there. Running it twice yields two directories.
    * It never invents a number. A reading that was not observed is written as
      the literal string ABSENT together with the reason it is absent; numeric
      fields are null, never 0.

  Where the readings come from:

    1. The product side (tokens, log lines, gate posture) - the resident's
       persistent JSONL sink under <DataDir>\logs\, by line number, anchored to
       the running process's own boot record.
    2. The numbers (quiet / speak / tone levels, the reference band, the
       before/after delta) - levelprobe\main.go, built and run from this
       directory. It opens the same default capture endpoint through the same
       internal/audio functions the product calls, and computes both the
       shipping ruler and the pre-449 ruler on the SAME frame. It is a measuring
       instrument, not the product, and it does not need the product to stop.

.PARAMETER OutRoot
  Root of the capture tree. Default D:\tmp\wisp.

.PARAMETER Label
  Directory name under <OutRoot>\<yyyyMMdd>\. Default 'realonce'.

.PARAMETER DataDir
  The running product's data root. Default %APPDATA%\wisp-dev (WISP_ENV=dev).

.PARAMETER StdoutLog
  Optional: the stdout/stderr file of the running ball, if somebody redirected
  it. When omitted the script records ABSENT for the fmt.Printf lines and uses
  the slog records instead (which carry the same facts).

.PARAMETER NoProbe
  Skip the level windows. Use only when no human is at the machine: the 291 /
  449 / 247 numeric cells then stay ABSENT.

.PARAMETER NoHotkey
  Skip the hot-key pauses. The 290 cells then stay ABSENT.

.PARAMETER NoPause
  Do not wait for Enter anywhere (non-interactive drill). Timestamps are still
  recorded, so an empty gesture window is visible as such.

.EXAMPLE
  pwsh -NoProfile -File .\capture-once.ps1

.EXAMPLE
  # dry drill with no human: proves the collector runs and marks the gaps
  pwsh -NoProfile -File .\capture-once.ps1 -NoPause -NoHotkey -NoProbe
#>
[CmdletBinding()]
param(
    [string]$Repo = '',
    [string]$DataDir = (Join-Path $env:APPDATA 'wisp-dev'),
    [string]$OutRoot = 'D:\tmp\wisp',
    [string]$Label = 'realonce',
    [string]$Date = (Get-Date -Format 'yyyyMMdd'),
    [string]$StdoutLog = '',
    [int]$QuietSecs = 4,
    [int]$SpeakSecs = 5,
    [int]$ToneSecs = 4,
    [switch]$NoProbe,
    [switch]$NoHotkey,
    [switch]$NoPause
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
try { [Console]::OutputEncoding = [Text.UTF8Encoding]::new($false) } catch { }

# ---------------------------------------------------------------------------
# constants the readings are judged against - every one of them is a line of
# product code, quoted here so the report can name its source instead of
# remembering a number.
# ---------------------------------------------------------------------------
$GATE = 0.06            # internal/ball/liquid.go:30  SilenceLevelGate = 0.06
$HOLD_MS = 150          # internal/ball/liquid.go:31  SilenceHoldMs
$FRAME_MS = 32          # internal/audio/audio.go:53   FrameDuration (512 @ 16k)
$FRAME_HZ = 1000.0 / $FRAME_MS
$FULLSCALE = 32768.0    # internal/audio/level.go:67   LevelFullScale
$LSB = 1.0 / $FULLSCALE # internal/audio/level.go:72   LevelLSB
$BAND_LO = 0.0013       # ticket 449 AC#8's reference band (landing table, measured)
$BAND_HI = 0.0031       # ticket 449 AC#8's reference band (landing table, measured)
$PRE_FIX_QUIET = 0.368462  # ticket 247 AC#2 v1's quiet reading, pre-449
$PRE_FIX_DC_COUNTS = $PRE_FIX_QUIET * $FULLSCALE  # ~12073.5 counts of DC

$script:Absent = 'ABSENT'
$script:Notes = [System.Collections.Generic.List[string]]::new()
$script:Cells = [System.Collections.Generic.List[object]]::new()

# ---------------------------------------------------------------------------
# small helpers
# ---------------------------------------------------------------------------
function Resolve-Repo([string]$Given) {
    if ($Given) {
        if (-not (Test-Path (Join-Path $Given 'go.mod'))) { throw "Repo '$Given' has no go.mod" }
        return (Resolve-Path $Given).Path
    }
    $p = $PSScriptRoot
    foreach ($i in 1..8) {
        $p = Split-Path -Parent $p
        if ($p -and (Test-Path (Join-Path $p 'go.mod'))) { return $p }
    }
    throw "could not derive the repo root from $PSScriptRoot; pass -Repo"
}

function New-RunDir([string]$root, [string]$date, [string]$label) {
    $dayDir = Join-Path $root $date
    if (-not (Test-Path $dayDir)) { $null = New-Item -ItemType Directory -Path $dayDir -Force }
    $candidates = @($label) + @(1..50 | ForEach-Object { '{0}-{1:000}' -f $label, $_ })
    foreach ($c in $candidates) {
        $d = Join-Path $dayDir $c
        try {
            # New-Item without -Force throws when the directory exists: that is
            # the atomic "never reuse a capture" guard.
            $null = New-Item -ItemType Directory -Path $d -ErrorAction Stop
            return $d
        } catch {
            continue
        }
    }
    throw "could not allocate a fresh run directory under $dayDir"
}

function Write-Utf8([string]$path, [string]$text) {
    if (Test-Path $path) { throw "refusing to overwrite existing capture file: $path" }
    [IO.File]::WriteAllText($path, $text, [Text.UTF8Encoding]::new($false))
}

function Read-Sink([string]$path) {
    # The resident keeps the sink open for append, so a plain ReadAllLines is
    # refused ("being used by another process"). Open it with FileShare.ReadWrite:
    # this is a read of a live log, not a claim that the file is ours.
    $fs = $null; $sr = $null
    $all = [System.Collections.Generic.List[string]]::new()
    try {
        $fs = [IO.File]::Open($path, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::ReadWrite)
        $sr = [IO.StreamReader]::new($fs, [Text.UTF8Encoding]::new($false), $true)
        while ($null -ne ($l = $sr.ReadLine())) { $all.Add($l) }
    } finally {
        if ($sr) { $sr.Dispose() }
        if ($fs) { $fs.Dispose() }
    }
    # Bounded read: the sink is append-only and shared by every process using
    # this data root, so we keep the last N lines and say so when we truncate.
    $limit = 20000
    $truncated = $false
    $start = 0
    if ($all.Count -gt $limit) { $start = $all.Count - $limit; $truncated = $true }
    $rows = [System.Collections.Generic.List[object]]::new()
    for ($i = $start; $i -lt $all.Count; $i++) {
        $line = $all[$i]
        if ([string]::IsNullOrWhiteSpace($line)) { continue }
        try {
            $obj = $line | ConvertFrom-Json -AsHashtable
        } catch {
            continue
        }
        $t = $null
        try { $t = [datetimeoffset]::Parse([string]$obj['time']) } catch { }
        $rows.Add([pscustomobject]@{
                LineNo = $i + 1
                Time   = $t
                Msg    = [string]$obj['msg']
                Level  = [string]$obj['level']
                Raw    = $line
                Fields = $obj
            })
    }
    return [pscustomobject]@{ Rows = $rows; Truncated = $truncated; Total = $all.Count }
}

function Select-Hits($sink, [DateTimeOffset]$from, [DateTimeOffset]$to) {
    if ($null -eq $sink) { return @() }
    return @($sink.Rows | Where-Object {
            $_.Time -ne $null -and $_.Time -ge $from -and $_.Time -le $to
        })
}

function Add-Cell([string]$id, [string]$label, [string]$ticket, [string]$box,
    $value, [string]$status, [string]$source, [string]$note = '') {
    $script:Cells.Add([pscustomobject]@{
            id = $id; label = $label; ticket = $ticket; box = $box
            value = $value; status = $status; source = $source; note = $note
        })
}

function Fmt($v, [string]$format = '0.000000') {
    if ($null -eq $v) { return $script:Absent }
    return ([double]$v).ToString($format, [cultureinfo]::InvariantCulture)
}

function Wait-Human([string]$prompt) {
    Write-Host ''
    Write-Host $prompt -ForegroundColor Yellow
    if ($NoPause) {
        Write-Host '  (NoPause: not waiting)' -ForegroundColor DarkGray
        return [pscustomobject]@{ PromptAt = (Get-Date).ToUniversalTime(); EnterAt = (Get-Date).ToUniversalTime() }
    }
    $null = Read-Host '  按回车继续（Enter when done）'
    return [pscustomobject]@{ PromptAt = (Get-Date).ToUniversalTime(); EnterAt = (Get-Date).ToUniversalTime() }
}

# ---------------------------------------------------------------------------
# 0. preflight
# ---------------------------------------------------------------------------
$repoPath = Resolve-Repo $Repo
$runDir = New-RunDir $OutRoot $Date $Label
$startedAt = (Get-Date)

Write-Host "wisp real-machine once - capture run" -ForegroundColor Cyan
Write-Host "  repo     : $repoPath"
Write-Host "  data dir : $DataDir"
Write-Host "  run dir  : $runDir"

$facts = [ordered]@{
    tool            = 'capture-once.ps1'
    started_at      = $startedAt.ToString('o')
    host            = $env:COMPUTERNAME
    user            = $env:USERNAME
    pwsh            = $PSVersionTable.PSVersion.ToString()
    os              = [Environment]::OSVersion.VersionString
    repo            = $repoPath
    repo_head       = ''
    repo_dirty_files = $null
    data_dir        = $DataDir
    run_dir         = $runDir
    quiet_secs      = $QuietSecs
    speak_secs      = $SpeakSecs
    tone_secs       = $ToneSecs
}

Push-Location $repoPath
try {
    $facts.repo_head = (& git rev-parse HEAD 2>$null | Select-Object -First 1)
    $dirty = @(& git status --porcelain 2>$null)
    $facts.repo_dirty_files = $dirty.Count
} finally { Pop-Location }

# --- the running product: read the marker, never touch the process -----------
$markerPath = Join-Path $DataDir 'run-marker.json'
$marker = $null
if (Test-Path $markerPath) {
    try { $marker = Get-Content $markerPath -Raw -Encoding UTF8 | ConvertFrom-Json } catch { }
}
$procInfo = $null
if ($marker -and $marker.pid) {
    $p = Get-Process -Id ([int]$marker.pid) -ErrorAction SilentlyContinue
    if ($p) {
        $procInfo = [ordered]@{
            pid = $p.Id; name = $p.ProcessName; path = $p.Path
            start_time = $p.StartTime.ToString('o')
        }
    } else {
        $procInfo = [ordered]@{ pid = [int]$marker.pid; alive = $false; note = 'run-marker names a pid that is not running' }
    }
}
$facts.running_product = $procInfo
$facts.wisp_processes = @(Get-Process -Name 'wisp' -ErrorAction SilentlyContinue).Count
$facts.balldebug_processes = @(Get-Process -Name 'balldebug' -ErrorAction SilentlyContinue).Count
Write-Host ("  ball     : " + $(if ($procInfo) { "pid $($procInfo.pid) ($($procInfo.name))" } else { 'ABSENT - no run-marker.json' }))

# --- config facts, values only (never the hot-key string itself) ------------
$cfgPath = Join-Path $DataDir 'config.toml'
$cfgFacts = [ordered]@{ path = $cfgPath; present = (Test-Path $cfgPath) }
if ($cfgFacts.present) {
    $section = ''
    foreach ($line in [IO.File]::ReadAllLines($cfgPath)) {
        $trim = $line.Trim()
        if ($trim -match '^\[(.+)\]$') { $section = $matches[1]; continue }
        if ($trim -match '^([A-Za-z0-9_.]+)\s*=\s*(.*)$') {
            $k = $matches[1]; $v = $matches[2].Trim()
            switch ("$section/$k") {
                'voice/enabled' { $cfgFacts.voice_enabled = $v }
                'voice.wake_word/enabled' { $cfgFacts.wake_word_enabled = $v }
                'audio/mic_muted_default' { $cfgFacts.mic_muted_default = $v }
                'audio/input_device' { $cfgFacts.input_device = $v }
                'hotkey/mute' { $cfgFacts.hotkey_mute = $(if ($v -eq "''" -or $v -eq '""') { '<empty: compiled default>' } else { '<set: value not recorded>' }) }
            }
        }
    }
}
$facts.config = $cfgFacts
Write-Host ("  config   : mic_muted_default=$(if ($cfgFacts.Contains('mic_muted_default')) { $cfgFacts.mic_muted_default } else { 'ABSENT' }) voice.enabled=$(if ($cfgFacts.Contains('voice_enabled')) { $cfgFacts.voice_enabled } else { 'ABSENT' })")

# --- the persistent sink ----------------------------------------------------
$sinkPath = $null
$sinkCandidates = @(Get-ChildItem (Join-Path $DataDir 'logs') -Filter "wisp-$Date-*.jsonl" -ErrorAction SilentlyContinue |
        Sort-Object LastWriteTime -Descending)
if ($sinkCandidates.Count -gt 0) { $sinkPath = $sinkCandidates[0].FullName }
$facts.log_sink = [ordered]@{
    path = $sinkPath; exists = [bool]$sinkPath
    mtime = $(if ($sinkPath) { (Get-Item $sinkPath).LastWriteTime.ToString('o') } else { $null })
    bytes = $(if ($sinkPath) { (Get-Item $sinkPath).Length } else { $null })
}
$sink = $null
if ($sinkPath) {
    $sink = Read-Sink $sinkPath
    $facts.log_sink.lines = $sink.Total
    $facts.log_sink.truncated = $sink.Truncated
    Write-Host ("  sink     : $sinkPath ($($sink.Total) lines)")
} else {
    Write-Host '  sink     : ABSENT - no wisp-<date>-*.jsonl under the data root' -ForegroundColor Red
}

# --- optional stdout log ----------------------------------------------------
$stdoutFacts = [ordered]@{ provided = [bool]$StdoutLog; path = $StdoutLog; exists = $false }
if ($StdoutLog) {
    if (Test-Path $StdoutLog) {
        $stdoutFacts.exists = $true
        $stdoutFacts.bytes = (Get-Item $StdoutLog).Length
        $stdoutFacts.mtime = (Get-Item $StdoutLog).LastWriteTime.ToString('o')
    }
}
$facts.stdout_log = $stdoutFacts

# --- ticket 290 AC#3 shape 1 (free): the boot record of the running process --
$bootWindowStart = [DateTimeOffset]::MinValue
if ($procInfo -and $procInfo.Contains('start_time')) {
    $bootWindowStart = [DateTimeOffset]::Parse($procInfo.start_time).AddSeconds(-2)
}
$bootRows = @()
if ($sink -and $bootWindowStart -ne [DateTimeOffset]::MinValue) {
    $bootRows = @($sink.Rows | Where-Object {
            $_.Time -ne $null -and $_.Time -ge $bootWindowStart -and
            $_.Msg -match 'capture armed but muted at boot|capture leg running|capture leg not built|config unreadable at boot'
        })
}
if ($bootRows.Count -gt 0) {
    $b = $bootRows[-1]
    Add-Cell '290-1a' '默认档开机：门处于静音、设备未打开' '290' 'AC#3①' `
        ("{0} (config_source={1})" -f $b.Msg, $b.Fields['config_source']) 'READ' `
        ("{0}:{1} @ {2:o}" -f $sinkPath, $b.LineNo, $b.Time) `
        'free read: the current process booted with this posture; nothing was restarted'
} else {
    Add-Cell '290-1a' '默认档开机：门处于静音、设备未打开' '290' 'AC#3①' $script:Absent 'ABSENT' `
        'no boot record at/after the running pid start time in the sink' 'the process may predate this sink file'
}

# ---------------------------------------------------------------------------
# 1. hot key step (ticket 290 AC#3 shape 2): unmute, keep it open for the probe
# ---------------------------------------------------------------------------
$unmuteAt = $null
$remuteAt = $null
if ($NoHotkey) {
    Add-Cell '290-2a' '按一次静音热键 ⇒ 取消静音行' '290' 'AC#3②' $script:Absent 'ABSENT' '-NoHotkey' 'skipped on request'
    Add-Cell '290-2b' '日志出现 mute-hotkey 这一枚手势词' '290' 'AC#3②' $script:Absent 'ABSENT' '-NoHotkey' 'skipped on request'
} else {
    $w = Wait-Human @'
【第 1 步 / Step 1】请按一次「静音热键」（ticket 290）。
  效果：悬浮球那个麦克风门从静音位打开（日志会写一行 已取消静音）。
  不确定是哪个键：看球启动那一行里的 mute=… ，或右键托盘里的「静音」项
  （托盘那支走同一个执行者，但手势字段是 tray-mute，290 要的 mute-hotkey 得用热键）。
  按完再按回车。会保持开门状态，给下一步的真麦读数用。
'@
    $from = $w.PromptAt.AddSeconds(-1)
    $to = $w.EnterAt.AddSeconds(2)
    # Re-read the sink AFTER the hand: the preflight read happened before the key
    # was pressed, and the record this step is looking for did not exist then.
    if ($sinkPath) { $sink = Read-Sink $sinkPath }
    $hits = @(Select-Hits $sink $from $to | Where-Object { $_.Msg -match 'mute gesture' })
    $unmuteRow = $hits | Where-Object { $_.Fields['outcome'] -match '取消静音' } | Select-Object -First 1
    $unmuteAt = $w.EnterAt
    if ($unmuteRow) {
        $gesture = [string]$unmuteRow.Fields['gesture']
        Add-Cell '290-2a' '按一次静音热键 ⇒ 取消静音行' '290' 'AC#3②' `
            ([string]$unmuteRow.Fields['outcome']) 'READ' `
            ("{0}:{1} @ {2:o}" -f $sinkPath, $unmuteRow.LineNo, $unmuteRow.Time)
        if ($gesture -eq 'mute-hotkey') {
            Add-Cell '290-2b' '日志出现 mute-hotkey 这一枚手势词' '290' 'AC#3②' $gesture 'READ' `
                ("{0}:{1} (gesture field)" -f $sinkPath, $unmuteRow.LineNo)
        } else {
            Add-Cell '290-2b' '日志出现 mute-hotkey 这一枚手势词' '290' 'AC#3②' "gesture=$gesture" 'ABSENT' `
                ("{0}:{1}" -f $sinkPath, $unmuteRow.LineNo) 'the gate DID turn, but not through the hot key - ticket 290 AC#3② names mute-hotkey'
        }
    } else {
        Add-Cell '290-2a' '按一次静音热键 ⇒ 取消静音行' '290' 'AC#3②' $script:Absent 'ABSENT' `
            ("no 'mute gesture' record in {0:o}..{1:o} of {2}" -f $from, $to, $sinkPath) `
            'either the key never reached the ball, or the sink is not this data root'
        Add-Cell '290-2b' '日志出现 mute-hotkey 这一枚手势词' '290' 'AC#3②' $script:Absent 'ABSENT' 'same window' 'same reason'
    }
    $script:Notes.Add('290: 手势窗口 = ' + $from.ToString('o') + '..' + $to.ToString('o'))
}

# ---------------------------------------------------------------------------
# 2. the level windows (tickets 291 / 449 / 247)
# ---------------------------------------------------------------------------
$probe = $null
$probeDir = Join-Path $runDir 'probe'
$framesRows = @()
if ($NoProbe) {
    $script:Notes.Add('probe skipped (-NoProbe): 291 / 449 / 247 numeric cells stay ABSENT')
} else {
    $probeSrc = Join-Path $PSScriptRoot 'levelprobe\main.go'
    if (-not (Test-Path $probeSrc)) { throw "probe source missing: $probeSrc" }
    $probeExe = Join-Path $runDir 'levelprobe.exe'
    Write-Host ''
    Write-Host '正在编译电平探针（go build，约 2-5 秒）…' -ForegroundColor Cyan
    Push-Location $repoPath
    try {
        $env:GOFLAGS = ''
        & go build -o $probeExe $probeSrc 2>&1 | ForEach-Object { Write-Host "    $_" }
        if ($LASTEXITCODE -ne 0) { throw "go build of the probe failed (rc=$LASTEXITCODE)" }
    } finally { Pop-Location }

    $null = New-Item -ItemType Directory -Path $probeDir -ErrorAction Stop
    $probeLog = Join-Path $runDir 'levelprobe.console.log'
    if (Test-Path $probeLog) { throw "refusing to overwrite $probeLog" }
    Write-Host ''
    Write-Host '【第 2 步 / Step 2】接下来 17 秒，请照屏幕上的提示做：先安静 4 秒，再正常说 5 秒，再安静 4 秒，最后 4 秒是本程序自己放音（别说话）。' -ForegroundColor Yellow
    $probeArgs = @('-out', $probeDir, '-quiet', $QuietSecs, '-speak', $SpeakSecs, '-tone', $ToneSecs)
    & $probeExe @probeArgs 2>&1 | Tee-Object -FilePath $probeLog
    $probeRc = $LASTEXITCODE
    $facts.probe_exit_code = $probeRc
    $probeJson = Join-Path $probeDir 'levelprobe.json'
    if (Test-Path $probeJson) {
        $probe = Get-Content $probeJson -Raw -Encoding UTF8 | ConvertFrom-Json
    } else {
        $script:Notes.Add("probe wrote no levelprobe.json (rc=$probeRc)")
    }
    $framesPath = Join-Path $probeDir 'levelprobe.frames.jsonl'
    if (Test-Path $framesPath) {
        $framesRows = @(Get-Content $framesPath -Encoding UTF8 | Where-Object { $_.Trim() } | ForEach-Object { $_ | ConvertFrom-Json })
    }
    $facts.probe_frames = $framesRows.Count
}

$W = @{}
if ($probe) { foreach ($x in $probe.windows) { $W[$x.name] = $x } }

if ($W.Count -gt 0) {
    # ---- per-window cells: 291 AC#2 shapes 1/2/3 -------------------------
    $shape = @(
        @{ name = 'quiet-1'; id = '291-1'; label = '形①完全安静'; secs = $QuietSecs },
        @{ name = 'speak'; id = '291-2'; label = '形②正常音量说话'; secs = $SpeakSecs },
        @{ name = 'quiet-2'; id = '291-1b'; label = '形①的对照窗（第二段安静）'; secs = $QuietSecs },
        @{ name = 'tone'; id = '291-3'; label = '形③外放已知幅度音调'; secs = $ToneSecs }
    )
    foreach ($s in $shape) {
        $win = $W[$s.name]
        if ($null -eq $win) {
            Add-Cell $s.id $s.label '291' 'AC#2' $script:Absent 'ABSENT' 'probe json has no such window' ''
            continue
        }
        $n = $win.new_level_dc_stripped
        $cnt = [int]$win.samples
        if ($cnt -eq 0) {
            Add-Cell $s.id $s.label '291' 'AC#2' $script:Absent 'ABSENT' `
                ("levelprobe.json window '$($s.name)': 0 samples") 'window recorded nothing'
        } else {
            $val = "samples={0} mean={1} max={2} span={3}s levels_per_s={4}" -f `
                $cnt, (Fmt $n.mean), (Fmt $n.max), (Fmt $win.span_s '0.00'), (Fmt $win.levels_per_s '0.00')
            $status = 'READ'
            $note = ''
            if ([double]$win.span_s -lt $s.secs) { $note = "span shorter than requested $($s.secs)s" }
            Add-Cell $s.id $s.label '291' 'AC#2' $val $status `
                ("{0}\levelprobe.json .windows[{1}] + levelprobe.frames.jsonl" -f $probeDir, $s.name) $note
        }
    }

    # ---- 291 shape 4: synthetic, no microphone --------------------------
    if ($probe.synthetic) {
        $syn = ($probe.synthetic | ForEach-Object {
                '{0}: measured={1} expected={2} within={3}' -f $_.name, (Fmt $_.measured '0.############'), (Fmt $_.expected '0.############'), $_.within_tolerance
            }) -join ' | '
        $allWithin = @($probe.synthetic | Where-Object { -not $_.within_tolerance }).Count -eq 0
        Add-Cell '291-4' '形④合成帧（纯 DC + 已知正弦，⛔ 不据此判生产）' '291' 'AC#2' $syn `
            $(if ($allWithin) { 'READ' } else { 'READ-DISAGREES' }) `
            'levelprobe.json .synthetic (production LevelOfSamples, closed forms from level.go)' `
            $(if ($allWithin) { 'all within SineLevelTolerance=5e-5' } else { 'at least one case outside tolerance - do not overwrite, report it' })
    } else {
        Add-Cell '291-4' '形④合成帧（纯 DC + 已知正弦）' '291' 'AC#2' $script:Absent 'ABSENT' 'probe produced no synthetic block' ''
    }

    # ---- the shared quiet/speak numbers ---------------------------------
    $q = $W['quiet-1']; $sp = $W['speak']; $q2 = $W['quiet-2']; $tn = $W['tone']
    $qMean = if ($q) { $q.new_level_dc_stripped.mean } else { $null }
    $qMax = if ($q) { $q.new_level_dc_stripped.max } else { $null }
    $spMean = if ($sp) { $sp.new_level_dc_stripped.mean } else { $null }
    $spMax = if ($sp) { $sp.new_level_dc_stripped.max } else { $null }

    # ---- ticket 291 AC#2's verdict: can an EXISTING criterion separate them?
    if ($null -ne $qMean -and $null -ne $spMean) {
        $ratio = if ([double]$qMean -gt 0) { [double]$spMean / [double]$qMean } else { $null }
        $sepMean = if ($null -ne $ratio) { $ratio -gt 2 } else { $false }   # ticket 247 AC#2's own 2x gate
        $sepMax = ([double]$spMax -gt [double]$qMax)
        $gateSpeech = ([double]$spMax -ge $GATE)                              # liquid.go:272: raw >= gate => not silent
        $verdict = if ($sepMean -and $sepMax) { 'FILLED-SEPARABLE' } elseif ($sepMax) { 'FILLED-MAX-ONLY' } else { 'FILLED-NOT-SEPARABLE' }
        $val = "speak_mean/quiet_mean={0}x speak_max={1} quiet_max={2} 2x_mean_gate={3} max_gate={4} speak_max>=SilenceLevelGate={5}" -f `
            (Fmt $ratio '0.00'), (Fmt $spMax), (Fmt $qMax), $sepMean, $sepMax, $gateSpeech
        Add-Cell '291-5' '形②与形①能不能被既有判据分开' '291' 'AC#2' $val $verdict `
            'arithmetic over levelprobe.json windows; the criterion is ticket 247 AC#2''s existing 2x gate, not a new threshold' `
            $(if ($verdict -eq 'FILLED-NOT-SEPARABLE') { 'names itself unseparable - 291 AC#2 explicitly allows that answer; do not tune a threshold to pass' } else { '' })
    }

    # ---- ticket 449 AC#8 -------------------------------------------------
    if ($q -and [int]$q.samples -gt 0) {
        $gateVerdict = ([double]$qMean -lt $GATE)
        Add-Cell '449-a' '修后同一枚真麦的安静态 mean/max' '449' 'AC#8' `
            ("mean={0} max={1} (< SilenceLevelGate={2}: {3})" -f (Fmt $qMean), (Fmt $qMax), $GATE, $gateVerdict) 'READ' `
            ("{0}\levelprobe.json .windows[quiet-1]" -f $probeDir) `
            ("reference band from the landing table: {0}-{1} (measured there, not derivable from code)" -f $BAND_LO, $BAND_HI)
    }
    if ($framesRows.Count -gt 0) {
        $deltas = @($framesRows | ForEach-Object { [double]$_.new - [double]$_.old })
        $maxDelta = ($deltas | Measure-Object -Maximum).Maximum
        $violations = @($framesRows | Where-Object { [double]$_.new -gt [double]$_.old }).Count
        Add-Cell '449-b' '同一帧上 新尺−旧尺 的最大值（必须 ≤ 0）' '449' 'AC#8' `
            ("max(new-old)={0:E3} frames_with_new>old={1} of {2}" -f $maxDelta, $violations, $framesRows.Count) `
            $(if ($violations -eq 0) { 'READ' } else { 'READ-DISAGREES' }) `
            'levelprobe.frames.jsonl: both rulers computed on the same decoded frame' `
            'new <= old is provable from level.go (the mean minimises squared deviation); equality only when the frame mean is 0'
        $dcMean = ($framesRows | ForEach-Object { [double]$_.dc_counts } | Measure-Object -Average).Average
        $dcMax = ($framesRows | ForEach-Object { [math]::Abs([double]$_.dc_counts) } | Measure-Object -Maximum).Maximum
        $implied = [math]::Abs($dcMean) / $FULLSCALE
        Add-Cell '449-c' '今天这台机器的直流分量（决定 0.368 能不能复现）' '449' 'AC#8' `
            ("dc_mean={0:F1} counts dc_max_abs={1:F1} counts => pre-449 reading would be about {2:F6} from DC alone" -f $dcMean, $dcMax, $implied) 'READ' `
            'levelprobe.frames.jsonl .dc_counts (frame mean in int16 counts)' `
            ("ticket 247 v1's 0.368462 implies {0:F1} counts of DC; the fixture value 12345 counts reads 0.376740" -f $PRE_FIX_DC_COUNTS)
    }

    # ---- ticket 449 AC#8's gate direction, frame by frame ---------------
    if ($framesRows.Count -gt 0) {
        $best = 0; $run = 0; $bestPhase = ''
        foreach ($f in $framesRows) {
            if ($f.phase -match '^quiet') {
                if ([double]$f.new -lt $GATE) { $run++ } else { $run = 0 }
                if ($run -gt $best) { $best = $run; $bestPhase = [string]$f.phase }
            } else { $run = 0 }
        }
        $needFrames = [math]::Ceiling($HOLD_MS / $FRAME_MS)   # 150ms / 32ms = 5 frames
        $ok = ($best -ge $needFrames)
        Add-Cell '449-d' '安静窗里连续低于门限的帧数（≥5 帧 ⇒ 球判定“没人说话”）' '449' 'AC#8' `
            ("longest_run={0} frames ({1}{2} ms) in {3}; need >= {4} frames for SilenceHoldMs={5}ms" -f `
                $best, '', ($best * $FRAME_MS), $(if ($bestPhase) { $bestPhase } else { 'n/a' }), $needFrames, $HOLD_MS) `
            $(if ($ok) { 'READ' } else { 'READ-DISAGREES' }) `
            'levelprobe.frames.jsonl vs internal/ball/liquid.go:30-31 (SilenceLevelGate, SilenceHoldMs) and :272-276' `
            'this is the mechanism by which the ball decides nobody is speaking; it is arithmetic, not a picture of the ball'
    }

    # ---- ticket 247 AC#2 ------------------------------------------------
    if ($null -ne $qMean -and $null -ne $spMean) {
        $different = ([double]$spMean -ne [double]$qMean)
        Add-Cell '247-a' '说话那一发与安静那一发的电平不同（同一把尺、同一枚麦）' '247' 'AC#2' `
            ("quiet mean={0} max={1} | speak mean={2} max={3} | different={4}" -f (Fmt $qMean), (Fmt $qMax), (Fmt $spMean), (Fmt $spMax), $different) 'READ' `
            'levelprobe.json windows quiet-1 / speak' `
            "this is the seam's own number; the ball-side counter is 247-b and needs the process to stop"
    }
    if ($tn) {
        $epText = 'ABSENT'
        if ($probe.endpoints -and $probe.endpoints.capture) {
            $epText = "capture=$($probe.endpoints.capture) render=$($probe.endpoints.render)"
        } elseif ($probe.endpoints_error) {
            $epText = "ABSENT ($($probe.endpoints_error))"
        }
        $toneRms = 'ABSENT'
        if ($probe.tone -and $probe.tone.digital_rms_fullscale) {
            $toneRms = (Fmt $probe.tone.digital_rms_fullscale '0.000000')
        }
        Add-Cell '247-c' '出处：端点、帧率、放音幅度' '247' 'AC#2' `
            ("{0} levels_per_s={1} tone_digital_rms={2}" -f $epText, (Fmt $tn.levels_per_s '0.00'), $toneRms) 'READ' `
            'levelprobe.json .endpoints / .windows[tone] / .tone'
    }
}

# ---------------------------------------------------------------------------
# 3. put the machine back: press again -> muted (the shipped boot posture)
# ---------------------------------------------------------------------------
if (-not $NoHotkey -and -not $NoPause) {
    $w2 = Wait-Human @'
【第 3 步 / Step 3】再按一次同一个静音热键，把麦克风关回静音位（恢复出厂开机态，别让它开着）。
  按完再按回车。
'@
    $from2 = $w2.PromptAt.AddSeconds(-1)
    $to2 = $w2.EnterAt.AddSeconds(2)
    if ($sinkPath) { $sink = Read-Sink $sinkPath }
    $hits2 = @(Select-Hits $sink $from2 $to2 | Where-Object { $_.Msg -match 'mute gesture' })
    $muteRow = $hits2 | Where-Object { $_.Fields['outcome'] -match '已静音' } | Select-Object -First 1
    if ($muteRow) {
        Add-Cell '290-2d' '再按一次 ⇒ 已静音行（两个方向都动了）' '290' 'AC#3②' `
            ([string]$muteRow.Fields['outcome']) 'READ' `
            ("{0}:{1} @ {2:o}" -f $sinkPath, $muteRow.LineNo, $muteRow.Time)
    } else {
        Add-Cell '290-2d' '再按一次 ⇒ 已静音行（两个方向都动了）' '290' 'AC#3②' $script:Absent 'ABSENT' `
            ("no 'mute gesture' record in {0:o}..{1:o}" -f $from2, $to2) 'the gate may still be open - check the tray item by hand'
        $script:Notes.Add('WARNING: the mute gesture back was not observed; the microphone may still be open.')
    }
}

# ---------------------------------------------------------------------------
# 4. the two cells that this run CANNOT fill, said out loud
# ---------------------------------------------------------------------------
$levelsDeliveredRows = @()
if ($sinkPath) {
    # Last read of the sink: everything above may have appended to it.
    $sink = Read-Sink $sinkPath
    $levelsDeliveredRows = @($sink.Rows | Where-Object { $_.Msg -match 'levels_delivered' })
}
if ($levelsDeliveredRows.Count -gt 0) {
    $r = $levelsDeliveredRows[-1]
    Add-Cell '290-2c' '日志出现 levels_delivered' '290' 'AC#3②' $r.Msg 'READ-HISTORIC' `
        ("{0}:{1} @ {2:o}" -f $sinkPath, $r.LineNo, $r.Time) `
        'this line exists only in the exit record (resident_audio_windows.go:429); it is NOT from the running process - read the timestamp before citing it'
} else {
    Add-Cell '290-2c' '日志出现 levels_delivered' '290' 'AC#3②' $script:Absent 'ABSENT' `
        'resident_audio_windows.go:429 prints levels_delivered only from stop(); the process is still running and this runbook must not stop it' `
        'named gap: needs one orderly exit (tray Exit / Ctrl+C) in a later observation, not this one'
}
Add-Cell '247-b' '球侧 levels_delivered > 0（SetAudioLevel 真收到了数）' '247' 'AC#2' $script:Absent 'ABSENT' `
    'same reason as 290-2c: the counter is only printed by stop()' `
    'named gap, same later observation'

$facts.finished_at = (Get-Date).ToString('o')

# Every cell the runbook promises appears in the table, even when nothing could
# be read: "skipped" and "absent" must be visible, never a missing row.
$required = @(
    @{ id = '290-1a'; t = '290'; b = 'AC#3①'; l = '默认档开机：门处于静音、设备未打开' },
    @{ id = '290-2a'; t = '290'; b = 'AC#3②'; l = '按一次静音热键 ⇒ 取消静音行' },
    @{ id = '290-2b'; t = '290'; b = 'AC#3②'; l = '日志出现 mute-hotkey 这一枚手势词' },
    @{ id = '290-2c'; t = '290'; b = 'AC#3②'; l = '日志出现 levels_delivered' },
    @{ id = '290-2d'; t = '290'; b = 'AC#3②'; l = '再按一次 ⇒ 已静音行' },
    @{ id = '291-1'; t = '291'; b = 'AC#2'; l = '形①完全安静' },
    @{ id = '291-1b'; t = '291'; b = 'AC#2'; l = '形①的对照窗（第二段安静）' },
    @{ id = '291-2'; t = '291'; b = 'AC#2'; l = '形②正常音量说话' },
    @{ id = '291-3'; t = '291'; b = 'AC#2'; l = '形③外放已知幅度音调' },
    @{ id = '291-4'; t = '291'; b = 'AC#2'; l = '形④合成帧（纯 DC + 已知正弦）' },
    @{ id = '291-5'; t = '291'; b = 'AC#2'; l = '形②与形①能不能被既有判据分开' },
    @{ id = '449-a'; t = '449'; b = 'AC#8'; l = '修后同一枚真麦的安静态 mean/max' },
    @{ id = '449-b'; t = '449'; b = 'AC#8'; l = '同一帧上 新尺−旧尺 的最大值（必须 ≤ 0）' },
    @{ id = '449-c'; t = '449'; b = 'AC#8'; l = '今天这台机器的直流分量' },
    @{ id = '449-d'; t = '449'; b = 'AC#8'; l = '安静窗里连续低于门限的帧数' },
    @{ id = '247-a'; t = '247'; b = 'AC#2'; l = '说话那一发与安静那一发的电平不同' },
    @{ id = '247-b'; t = '247'; b = 'AC#2'; l = '球侧 levels_delivered > 0' },
    @{ id = '247-c'; t = '247'; b = 'AC#2'; l = '出处：端点、帧率、放音幅度' }
)
$present = @($script:Cells | ForEach-Object { $_.id })
foreach ($r in $required) {
    if ($present -notcontains $r.id) {
        Add-Cell $r.id $r.l $r.t $r.b $script:Absent 'ABSENT' 'not reached in this run' `
            'see the Notes section for which step was skipped'
    }
}

# ---------------------------------------------------------------------------
# 5. write the capture (never overwrite)
# ---------------------------------------------------------------------------
Write-Utf8 (Join-Path $runDir 'facts.json') ($facts | ConvertTo-Json -Depth 8)
$cellsObj = [ordered]@{
    tool         = 'capture-once.ps1'
    generated_at = (Get-Date).ToString('o')
    gate         = $GATE
    ball_hold_ms = $HOLD_MS
    frame_ms     = $FRAME_MS
    cells        = $script:Cells
    notes        = $script:Notes
}
Write-Utf8 (Join-Path $runDir 'cells.json') ($cellsObj | ConvertTo-Json -Depth 8)

$sb = [Text.StringBuilder]::new()
$null = $sb.AppendLine('# One real-machine session - captured readings')
$null = $sb.AppendLine('')
$null = $sb.AppendLine("run dir: ``$runDir``  ")
$null = $sb.AppendLine("started: $($facts.started_at)  finished: $($facts.finished_at)  ")
$null = $sb.AppendLine("repo HEAD: ``$($facts.repo_head)``  ball pid: $(if ($procInfo) { $procInfo.pid } else { 'ABSENT' })  ")
$null = $sb.AppendLine("log sink: ``$sinkPath``  ")
$null = $sb.AppendLine('')
$null = $sb.AppendLine('`ABSENT` is the literal answer "not observed"; it is never a zero.')
$null = $sb.AppendLine('')
$null = $sb.AppendLine('| cell | what it is | ticket / box | status | value | source |')
$null = $sb.AppendLine('|---|---|---|---|---|---|')
foreach ($c in $script:Cells) {
    $null = $sb.AppendLine(("| ``{0}`` | {1} | {2} {3} | {4} | {5} | {6} |" -f `
                $c.id, $c.label, $c.ticket, $c.box, $c.status, ($c.value -replace '\|', '\|'), $c.source))
}
$null = $sb.AppendLine('')
$null = $sb.AppendLine('## Notes')
foreach ($n in $script:Notes) { $null = $sb.AppendLine("- $n") }
Write-Utf8 (Join-Path $runDir 'summary.md') $sb.ToString()

Write-Host ''
Write-Host "captured -> $runDir" -ForegroundColor Green
$readCount = @($script:Cells | Where-Object { $_.status -like 'READ*' }).Count
$absentCount = @($script:Cells | Where-Object { $_.status -eq 'ABSENT' }).Count
Write-Host ("  cells: {0} READ, {1} ABSENT, {2} total" -f $readCount, $absentCount, $script:Cells.Count)
Write-Host '  summary.md / cells.json / facts.json' -ForegroundColor Green
Write-Host '  the evidence file to fill next: docs/evidence/s1/real-machine-once.md' -ForegroundColor Green

# Exit: 0 when at least one ticket's cells were read and nothing disagrees;
# 3 when a whole family is absent; 4 when a reading contradicts the code.
if (@($script:Cells | Where-Object { $_.status -eq 'READ-DISAGREES' }).Count -gt 0) { exit 4 }
if ($absentCount -ge 5) { exit 3 }
exit 0
