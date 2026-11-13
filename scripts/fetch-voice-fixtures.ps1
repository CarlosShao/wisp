# scripts/fetch-voice-fixtures.ps1 - fetch the pinned voice fixtures into
# third_party/ so cmd/wisp's voice tests can RUN where they are judged
# (ticket 111 fixtures leg, 2026-11-13).
#
# WHY THIS EXISTS. Run 43389487057 (test-windows step 7) died before running a
# single test: the hosted runner's third_party/ (git-ignored; restored by
# actions/cache + build.ps1) carries the SHERPA DLLs but never the voice
# MODELS, so scripts/wisp-cli-staged-tests.sh's five WISP_* defaults all
# missed, the old carrier answered with a setup-failure exit 1, and
# portable-tests.sh's GUARD B rightly red-flagged "the package printed no
# top-level result line". Route A of the dispatch: put the fixtures ON the
# runner, next to the DLL staging, so the voice family is really executed
# there instead of being absent from the verdict.
#
# SINGLE SOURCE OF TRUTH: build/models/manifest.json - the same pin ledger
# `wisp models ensure` walks and the fixtures' own tests read
# (WISP_MODELS_MANIFEST). Every URL and every SHA256 below is READ from it,
# never retyped: a pin changed there is required here automatically, and a
# file whose bytes do not match its pin fails this script loudly (exit 1)
# instead of reaching the tests as a silently wrong fixture.
#
# WHAT IT FETCHES (exactly the five targets the staged carrier's env points
# at, and nothing else):
#   kws-zipformer-wenetspeech-3.3M-2024-01-01  -> third_party/spike-models/sherpa-onnx-kws-...
#   vad-silero                                 -> third_party/spike-models/silero_vad.onnx
#   asr-streaming-paraformer-zh-en             -> third_party/spike-models/asr-{encoder,decoder}.int8.onnx + asr-tokens.txt
#   tts-matcha-zh-baker (archive + vocoder)    -> third_party/spike-models/matcha-icefall-zh-baker/ + vocos-22khz-univ.onnx
#   punc-ct-transformer-zh-en-vocab272727      -> third_party/model-fetch/x-punc/sherpa-onnx-punct-...-int8/
#
# The ASR DESTINATION names are the DEV spellings (asr-encoder.int8.onnx, ...),
# not the manifest's store spellings (encoder.int8.onnx): the carrier's env and
# the fixtures' own walk-up fallbacks name the dev tree, and
# cmd/wisp/resident_voice_windows.go:113-114 documents exactly this
# dev-vs-store name difference. The mapping is stated once in $Targets below.
#
# IDEMPOTENCE. Each target is verified BEFORE any download: when the
# destination already carries every pinned file at its pinned SHA256 (the
# actions/cache hit shape), the target is a no-op. When a file is absent or
# mismatched it is re-downloaded. So on a warm cache this step does nothing
# but hash-verify; on a cold or partial one it fetches exactly what is missing.
#
# EXTRACTION. The pinned .tar.bz2 archives are unpacked with the OS tar
# (bsdtar), the same instrument scripts/fetch-deps.ps1 uses for the sherpa
# archive on this repo's CI (12/12 green builds), and every member the
# manifest pins is hash-verified AFTER extraction, before anything is moved
# into place. A tarball whose bytes do not match its archive pin is never
# extracted; an archive whose members do not match their member pins is never
# installed.
#
# EXIT CODES: 0 = every target verified present (downloaded or cached);
# 1 = named failure (a URL dead after all its variants, an archive or member
# hash mismatch). A red here is a SETUP failure for the steps after it - which
# is loud and correct: the carrier's named fallback will still run the rest of
# cmd/wisp and book the voice files by name, while this step says WHY.
#
# Offline dev hosts: this script is NOT the dev path - `wisp models ensure`
# (the manifest's own consumer) remains the documented route; this script
# exists for CI and for a maintainer who wants the CI shape on a laptop.

$ErrorActionPreference = 'Stop'

$repo = Split-Path -Parent $PSScriptRoot
$manifestPath = Join-Path $repo 'build\models\manifest.json'
if (-not (Test-Path -LiteralPath $manifestPath)) {
    Write-Host "fetch-voice-fixtures: FATAL: manifest not found: $manifestPath"
    exit 1
}
$manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json

function Get-Sha256([string]$Path) {
    (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

# Download the first URL that works, into $OutFile. No mirror prefix is added
# beyond what the manifest itself lists: the ASR entry carries hf-mirror and
# huggingface variants, the release-hosted entries carry their GitHub URLs, and
# the manifest's order is the ledger's own preference order.
function Invoke-DownloadFirst([string[]]$Urls, [string]$OutFile) {
    foreach ($u in $Urls) {
        try {
            Invoke-WebRequest -Uri $u -OutFile $OutFile -UseBasicParsing -TimeoutSec 1800
            return $true
        } catch {
            Write-Host ("fetch-voice-fixtures:   url failed: {0} ({1})" -f $u, $_.Exception.Message)
        }
    }
    return $false
}

# The five targets. urls/pins come from the manifest at run time; only the
# DESTINATION paths (and the ASR dev-name mapping) are written here, because
# destinations are this script's contract with the staged carrier's env.
$spike = Join-Path $repo 'third_party\spike-models'
$xpunc = Join-Path $repo 'third_party\model-fetch\x-punc'

$Targets = @(
    @{
        id       = 'vad-silero'
        kind     = 'file'
        dest     = (Join-Path $spike 'silero_vad.onnx')
    },
    @{
        id       = 'asr-streaming-paraformer-zh-en'
        kind     = 'files'
        files    = @(
            @{ member = 'encoder.int8.onnx'; dest = (Join-Path $spike 'asr-encoder.int8.onnx') },
            @{ member = 'decoder.int8.onnx'; dest = (Join-Path $spike 'asr-decoder.int8.onnx') },
            @{ member = 'tokens.txt';        dest = (Join-Path $spike 'asr-tokens.txt') }
        )
    },
    @{
        id       = 'tts-matcha-zh-baker'
        kind     = 'archive'
        entry    = 'matcha-icefall-zh-baker.tar.bz2'
        dest     = (Join-Path $spike 'matcha-icefall-zh-baker')
    },
    @{
        id       = 'tts-matcha-zh-baker'
        kind     = 'file'
        member   = 'vocos-22khz-univ.onnx'
        dest     = (Join-Path $spike 'vocos-22khz-univ.onnx')
    },
    @{
        id       = 'punc-ct-transformer-zh-en-vocab272727'
        kind     = 'archive'
        entry    = 'sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8.tar.bz2'
        dest     = (Join-Path $xpunc 'sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8')
    },
    @{
        id       = 'kws-zipformer-wenetspeech-3.3M-2024-01-01'
        kind     = 'archive'
        entry    = 'sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01.tar.bz2'
        dest     = (Join-Path $spike 'sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01')
    }
)

function Get-ModelEntry([string]$id) {
    foreach ($m in $manifest.models) { if ($m.id -eq $id) { return $m } }
    return $null
}

function Get-FileEntry($model, [string]$path) {
    if (-not $model.files) { return $null }
    foreach ($f in $model.files) { if ($f.path -eq $path) { return $f } }
    return $null
}

# Verify every (path, sha256) pair against the disk. $Pairs is a list of
# @{dest=...; sha256=...}. Returns $true only when every file exists AND
# matches. One mismatched byte fails the pair - never a size-only check.
function Test-PinnedPairs($Pairs) {
    foreach ($p in $Pairs) {
        if (-not (Test-Path -LiteralPath $p.dest -PathType Leaf)) { return $false }
        if ((Get-Sha256 $p.dest) -ne $p.sha256) { return $false }
    }
    return $true
}

$failed = 0

foreach ($t in $Targets) {
    $model = Get-ModelEntry $t.id
    if ($null -eq $model) {
        Write-Host "fetch-voice-fixtures: FATAL: manifest has no model id $($t.id)"
        exit 1
    }

    if ($t.kind -eq 'file') {
        # A single-file target: the vocoder is a model-level files entry (no
        # archive wrapper); the VAD is the model's own single artifact (its pin
        # sits on the model entry itself).
        $urls = $null; $sha = $null
        if ($t.member) {
            $fe = Get-FileEntry $model $t.member
            if ($null -eq $fe) { Write-Host "fetch-voice-fixtures: FATAL: $($t.id) has no files entry $($t.member)"; exit 1 }
            $urls = $fe.urls; $sha = $fe.sha256
        } else {
            $urls = $model.urls; $sha = $model.sha256
        }
        if (Test-PinnedPairs @(@{ dest = $t.dest; sha256 = $sha })) {
            Write-Host "fetch-voice-fixtures: $($t.id): ok (cached, sha256 verified)"
            continue
        }
        $tmp = "$t.dest.part"
        Write-Host "fetch-voice-fixtures: $($t.id): downloading -> $($t.dest)"
        if (-not (Invoke-DownloadFirst $urls $tmp)) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): every URL failed"
            $failed = 1; continue
        }
        if ((Get-Sha256 $tmp) -ne $sha) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): downloaded bytes do not match the manifest pin (sha256)"
            Remove-Item -LiteralPath $tmp -Force
            $failed = 1; continue
        }
        New-Item -ItemType Directory -Path (Split-Path -Parent $t.dest) -Force | Out-Null
        Move-Item -LiteralPath $tmp -Destination $t.dest -Force
        Write-Host "fetch-voice-fixtures: $($t.id): ok (downloaded, sha256 verified)"
        continue
    }

    if ($t.kind -eq 'files') {
        $pairs = @()
        foreach ($f in $t.files) {
            $fe = Get-FileEntry $model $f.member
            if ($null -eq $fe) { Write-Host "fetch-voice-fixtures: FATAL: $($t.id) has no files entry $($f.member)"; exit 1 }
            $pairs += @{ dest = $f.dest; sha256 = $fe.sha256; urls = $fe.urls }
        }
        if (Test-PinnedPairs $pairs) {
            Write-Host "fetch-voice-fixtures: $($t.id): ok (cached, sha256 verified)"
            continue
        }
        foreach ($p in $pairs) {
            if ((Test-Path -LiteralPath $p.dest -PathType Leaf) -and ((Get-Sha256 $p.dest) -eq $p.sha256)) {
                Write-Host "fetch-voice-fixtures:   $($p.dest | Split-Path -Leaf): ok (cached)"
                continue
            }
            $tmp = "$($p.dest).part"
            Write-Host "fetch-voice-fixtures:   downloading $($p.dest | Split-Path -Leaf)"
            if (-not (Invoke-DownloadFirst $p.urls $tmp)) {
                Write-Host "fetch-voice-fixtures: FATAL: $($t.id): every URL failed for $($p.dest | Split-Path -Leaf)"
                $failed = 1; continue
            }
            if ((Get-Sha256 $tmp) -ne $p.sha256) {
                Write-Host "fetch-voice-fixtures: FATAL: $($t.id): downloaded bytes do not match the manifest pin"
                Remove-Item -LiteralPath $tmp -Force
                $failed = 1; continue
            }
            New-Item -ItemType Directory -Path (Split-Path -Parent $p.dest) -Force | Out-Null
            Move-Item -LiteralPath $tmp -Destination $p.dest -Force
            Write-Host "fetch-voice-fixtures:   $($p.dest | Split-Path -Leaf): ok (downloaded, sha256 verified)"
        }
        continue
    }

    if ($t.kind -eq 'archive') {
        # The archive entry's member list IS the pin set for the extracted
        # tree. Verify the extracted destination against it first (cache hit
        # shape), then download -> verify archive -> extract -> verify members
        # -> move into place.
        $fe = Get-FileEntry $model $t.entry
        if ($null -eq $fe) { Write-Host "fetch-voice-fixtures: FATAL: $($t.id) has no archive entry $($t.entry)"; exit 1 }
        if (-not $fe.archive -or -not $fe.archive.files) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): archive entry $($t.entry) carries no pinned member list"
            exit 1
        }
        $memberPairs = @()
        foreach ($m in $fe.archive.files) {
            $memberPairs += @{ dest = (Join-Path $t.dest $m.path); sha256 = $m.sha256 }
        }
        if (Test-PinnedPairs $memberPairs) {
            Write-Host "fetch-voice-fixtures: $($t.id): ok (cached, all pinned members sha256 verified)"
            continue
        }
        $tmpRoot = Join-Path $env:TEMP ("wisp-voice-fetch-" + [guid]::NewGuid().ToString('N'))
        New-Item -ItemType Directory -Path $tmpRoot -Force | Out-Null
        $archiveTmp = Join-Path $tmpRoot $t.entry
        Write-Host "fetch-voice-fixtures: $($t.id): downloading archive ($($fe.size_bytes) bytes pinned)"
        if (-not (Invoke-DownloadFirst $fe.urls $archiveTmp)) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): every URL failed"
            Remove-Item -LiteralPath $tmpRoot -Recurse -Force
            $failed = 1; continue
        }
        if ((Get-Sha256 $archiveTmp) -ne $fe.sha256) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): archive bytes do not match the manifest pin"
            Remove-Item -LiteralPath $tmpRoot -Recurse -Force
            $failed = 1; continue
        }
        $ex = Join-Path $tmpRoot 'extracted'
        New-Item -ItemType Directory -Path $ex -Force | Out-Null
        & "$env:SystemRoot\System32\tar.exe" -xjf $archiveTmp -C $ex
        if ($LASTEXITCODE -ne 0) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): tar extraction failed (exit $LASTEXITCODE)"
            Remove-Item -LiteralPath $tmpRoot -Recurse -Force
            $failed = 1; continue
        }
        # Locate the extracted top directory (the archive is expected to carry
        # exactly one; more or fewer is a shape this script refuses to guess).
        $tops = @(Get-ChildItem -LiteralPath $ex -Directory)
        if ($tops.Count -ne 1) {
            Write-Host "fetch-voice-fixtures: FATAL: $($t.id): archive extracts to $($tops.Count) top-level directories, expected 1"
            Remove-Item -LiteralPath $tmpRoot -Recurse -Force
            $failed = 1; continue
        }
        $tree = Join-Path $tops[0].FullName ''
        $bad = $false
        foreach ($p in $memberPairs) {
            $rel = $p.dest.Substring($t.dest.Length).TrimStart('\', '/')
            $cand = Join-Path $tree $rel
            if (-not (Test-Path -LiteralPath $cand -PathType Leaf) -or (Get-Sha256 $cand) -ne $p.sha256) {
                Write-Host "fetch-voice-fixtures: FATAL: $($t.id): extracted member does not match its pin: $rel"
                $bad = $true
            }
        }
        if ($bad) {
            Remove-Item -LiteralPath $tmpRoot -Recurse -Force
            $failed = 1; continue
        }
        # Same cache-refresh pattern fetch-deps.ps1 uses for third_party\sherpa-onnx:
        # replace the destination tree wholesale with the verified one.
        if (Test-Path -LiteralPath $t.dest) { Remove-Item -LiteralPath $t.dest -Recurse -Force }
        New-Item -ItemType Directory -Path (Split-Path -Parent $t.dest) -Force | Out-Null
        Move-Item -LiteralPath $tree -Destination $t.dest
        Remove-Item -LiteralPath $tmpRoot -Recurse -Force
        Write-Host "fetch-voice-fixtures: $($t.id): ok (extracted, all pinned members sha256 verified)"
        continue
    }
}

if ($failed -ne 0) {
    Write-Host "fetch-voice-fixtures: FATAL: one or more targets failed (see the named lines above)"
    exit 1
}

Write-Host "fetch-voice-fixtures: all five targets verified:"
Write-Host "fetch-voice-fixtures:   WISP_SPIKE_MODELS=$spike"
Write-Host "fetch-voice-fixtures:   WISP_SILERO_MODEL=$(Join-Path $spike 'silero_vad.onnx')"
Write-Host "fetch-voice-fixtures:   WISP_TTS_DIR=$(Join-Path $spike 'matcha-icefall-zh-baker')"
Write-Host "fetch-voice-fixtures:   WISP_KWS_DIR=$(Join-Path $spike 'sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01')"
Write-Host "fetch-voice-fixtures:   WISP_PUNC_MODEL=$(Join-Path $xpunc 'sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8\model.int8.onnx')"
exit 0
