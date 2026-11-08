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
# Usage: scripts/stage-sherpa-dlls.sh DEST_DIR
# Exit codes: 0 staged; 1 setup failure (named); 2 usage error.

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

while IFS= read -r dll; do
    [ -n "$dll" ] || continue
    cp -f "$dll_dir/$dll" "$dest/$dll"
    if [ ! -f "$dest/$dll" ]; then
        echo "stage-sherpa-dlls: GUARD - copy of $dll to $dest did not land" >&2
        exit 1
    fi
    echo "stage-sherpa-dlls: staged $dll -> $dest/$dll"
done <<<"$need"

exit 0
