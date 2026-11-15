#!/usr/bin/env bash
# stage-sherpa-dlls.sh DEST_DIR -- copy the pinned sherpa-onnx DLLs next to a
# test binary.
#
# WHY THIS EXISTS (ticket 409). internal/speech links the sherpa-onnx cgo
# import library, so the Windows loader resolves sherpa-onnx-c-api.dll and its
# second-level dependency onnxruntime.dll at PROCESS LOAD time, before any Go
# code can run. The loader searches the directory of the executable FIRST,
# then System32, and only then PATH. A dev host that carries an OS-provided
# onnxruntime.dll in C:\Windows\System32 (measured here: 1.17.1, 10,572,960
# bytes) therefore wins over anything on PATH, and the tests see onnxruntime
# 1.17.1 where the linked C API needs 1.28 -- internal/speech's
# checkOnnxruntime then turns every model-backed test into a skip (47 of them,
# measured 2026-11-08), which tools/d22scan/runtests.sh counts as red.
# Pre-pending third_party/sherpa-onnx to PATH (scripts/wisp-cli-tests.sh, the
# ticket 98 recipe) rescues the FIRST-level import only; it cannot rescue the
# second-level dependency, because System32 is searched before PATH. Staging
# the DLLs into the executable's own directory is the one lever that fires
# before the loader's System32 probe, and it is the same deployment shape
# scripts/build.ps1 already produces for wisp.exe in build/ (the binary and
# the three DLLs colocated).
#
# The DLL names are read from deps.toml ([sherpa-onnx.dll."..."] sections),
# not copied here: a pin added there is required here automatically, the same
# contract scripts/wisp-cli-tests.sh already relies on.
#
# TICKET 455 AC#2' (2026-11-15): the copy is now VERIFIED, not just present.
# After each cp, the staged file's sha256 is computed and compared against the
# same section's sha256 pin in deps.toml; ANY mismatch (a wrong source file, a
# truncated copy, a stray System32-shaped substitute) is a named GUARD that
# exits 1 -- refusing to run -- instead of printing a warning and continuing.
# Before this gate the script could "stage" a wrong onnxruntime.dll and print
# "staged onnxruntime.dll -> ..." right above a run that then bound that wrong
# image; the pin is the whole point of staging, so a copy that does not match
# the pin defeats the staging silently. sha256sum (Git Bash / CI) is used, with
# shasum as the named fallback; neither present is its own GUARD, not a skip.
#
# Usage: scripts/stage-sherpa-dlls.sh DEST_DIR
# Exit codes: 0 staged and sha256-verified against deps.toml; 1 setup failure
# (named: missing tree, missing pin, copy did not land, sha256 MISMATCH);
# 2 usage error.

set -u

if [ "$#" -ne 1 ]; then
    echo "stage-sherpa-dlls: usage: stage-sherpa-dlls.sh DEST_DIR" >&2
    exit 2
fi

dest=$1
if [ -z "$dest" ]; then
    echo "stage-sherpa-dlls: DEST_DIR is empty" >&2
    exit 2
fi
if [ ! -d "$dest" ]; then
    echo "stage-sherpa-dlls: GUARD - DEST_DIR '$dest' does not exist; create it first" \
        "(the caller owns the directory's lifetime, this script does not mkdir)" >&2
    exit 1
fi

root=$(cd "$(dirname "$0")/.." && pwd)
dll_dir="$root/third_party/sherpa-onnx"
if [ ! -d "$dll_dir" ]; then
    echo "stage-sherpa-dlls: GUARD - $dll_dir does not exist. Run scripts/fetch-deps.ps1" \
        "first (the windows job stages it via actions/cache on the deps.toml key)." >&2
    exit 1
fi

need=$(sed -n 's/^\[sherpa-onnx\.dll\."\(.*\)"\]$/\1/p' "$root/deps.toml")
if [ -z "$need" ]; then
    echo "stage-sherpa-dlls: GUARD - no [sherpa-onnx.dll.*] sections found in deps.toml," \
        "so this script has nothing to stage. Empty instrument: red, not green." >&2
    exit 1
fi

missing=""
while IFS= read -r dll; do
    [ -n "$dll" ] || continue
    if [ ! -f "$dll_dir/$dll" ]; then missing="$missing $dll"; fi
done <<<"$need"
if [ -n "$missing" ]; then
    echo "stage-sherpa-dlls: GUARD - pinned DLL(s) absent from $dll_dir:$missing" >&2
    echo "stage-sherpa-dlls: run scripts/fetch-deps.ps1; do not substitute copies from" >&2
    echo "stage-sherpa-dlls: System32 (that is the shadowed version this staging exists to beat)." >&2
    exit 1
fi

# Ticket 455 AC#2': sha256 of one file, whichever hasher this shell carries.
# sha256sum ships with Git Bash and the CI image; shasum is the perl-based
# fallback. Two parsing hazards are handled here, both measured 2026-11-15:
#   - GNU sha256sum PREFIXES THE WHOLE LINE WITH A BACKSLASH when the file
#     name contains a backslash (a Windows-style argument like D:\x\y does),
#     so the path is normalized to forward slashes before hashing and any
#     leading backslash is stripped from the answer anyway;
#   - a shell with neither hasher is a GUARD (exit 1), not a soft skip --
#     this gate must never degrade into "print a warning and continue".
stage_sha256() {
    hash_target=$(printf '%s' "$1" | tr '\\' '/')
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$hash_target" | sed 's/^\\//' | cut -d' ' -f1
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$hash_target" | sed 's/^\\//' | cut -d' ' -f1
    else
        echo "" >&2
        return 1
    fi
}

while IFS= read -r dll; do
    [ -n "$dll" ] || continue
    cp -f "$dll_dir/$dll" "$dest/$dll"
    if [ ! -f "$dest/$dll" ]; then
        echo "stage-sherpa-dlls: GUARD - copy of $dll to $dest did not land" >&2
        exit 1
    fi
    # Ticket 455 AC#2': the staged bytes must BE the pinned bytes. The pin's
    # sha256 is read from the SAME deps.toml section the name came from, so a
    # pin added there is verified here automatically.
    pin=$(awk -v dll="$dll" '
        $0 == "[sherpa-onnx.dll.\"" dll "\"]" { insec = 1; next }
        /^\[/ { insec = 0 }
        insec && $1 == "sha256" {
            gsub(/"/, "", $3); print $3; exit
        }
    ' "$root/deps.toml")
    if [ -z "$pin" ]; then
        echo "stage-sherpa-dlls: GUARD - no sha256 pin found in deps.toml for [sherpa-onnx.dll.\"$dll\"];" \
            "staging without verifying the bytes would let a wrong DLL defeat the whole point of staging" >&2
        exit 1
    fi
    got=$(stage_sha256 "$dest/$dll") || {
        echo "stage-sherpa-dlls: GUARD - neither sha256sum nor shasum is available; the ticket 455" >&2
        echo "stage-sherpa-dlls: AC#2' sha256 gate cannot run. Refusing to stage unverified (this is a" >&2
        echo "stage-sherpa-dlls: named failure, not a warning-and-continue)." >&2
        exit 1
    }
    if [ "$got" != "$pin" ]; then
        echo "stage-sherpa-dlls: GUARD - $dest/$dll sha256 MISMATCH: got $got, deps.toml pins $pin" >&2
        echo "stage-sherpa-dlls: the staged DLL is NOT the pinned image (wrong file, truncated copy, or a" >&2
        echo "stage-sherpa-dlls: substituted build). Refusing to run: re-run scripts/fetch-deps.ps1; do not" >&2
        echo "stage-sherpa-dlls: copy over it from System32 (that is the shadowed version this staging exists to beat)." >&2
        exit 1
    fi
    echo "stage-sherpa-dlls: staged $dll -> $dest/$dll (sha256 ok)"
done <<<"$need"

exit 0
