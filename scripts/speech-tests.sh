#!/usr/bin/env bash
# speech-tests.sh -- the windows carrier for "go test ./internal/speech/".
#
# WHY A CARRIER AND NOT BARE GO TEST (ticket 409). internal/speech links the
# sherpa-onnx cgo import library; the Windows loader resolves the whole import
# tree (sherpa-onnx-c-api.dll, and its second-level dependency
# onnxruntime.dll) at PROCESS LOAD time, before any Go code -- TestMain, init,
# LoadLibrary -- can influence it. On a host whose C:\Windows\System32 carries
# an OS-provided onnxruntime.dll (measured here: 1.17.1, 10,572,960 bytes),
# the loader picks that copy over anything on PATH (System32 is searched
# before PATH), checkOnnxruntime sees 1.17.1 where the linked C API needs
# 1.28, and every model-backed test is skipped -- 47 skips measured
# 2026-11-08, which tools/d22scan/runtests.sh counts as red (SKIP is not a
# pass, ticket 71 AC#3). The one lever that fires before the System32 probe
# is the loader's FIRST search location: the directory of the executable. So
# this step builds the test binary with "go test -c", stages the three pinned
# DLLs next to it (scripts/stage-sherpa-dlls.sh), and runs the binary from
# the package directory -- the same cwd contract "go test" itself uses, and
# the same binary-plus-DLLs-colocated deployment shape scripts/build.ps1
# already produces for wisp.exe in build/. sileroModelPath in
# internal/speech/fixtures_test.go walks up from the working directory
# precisely so this shape works without extra configuration.
#
# This mirrors scripts/wisp-cli-tests.sh for cmd/wisp: same guard style, same
# deps.toml pin source, same refusal to relax. On non-windows GOOS this
# refuses loudly rather than running a diluted scope.
#
# Usage: scripts/speech-tests.sh [extra -test.* flags...]
#   e.g. scripts/speech-tests.sh -test.run 'TestVAD'
# Exit codes: the test binary's exit code, or 1/2 on a named setup failure.
# Output: the test log is kept at $WISP_SPEECH_STAGE_DIR/speech-tests.log, and
# the same -v stream is ALSO printed on stdout (teed, not swallowed): the
# full-package runner (scripts/portable-tests.sh --scope=speech, ticket 409
# wave99) counts the four numbers by scanning `=== RUN` / `--- PASS` /
# `--- FAIL` / `--- SKIP` lines in what this step prints - a summary line alone
# would be "the word appeared in the log" wearing four numbers, the exact shape
# that runner exists to refuse.

set -u

if [ "$(go env GOOS)" != windows ]; then
    echo "speech-tests.sh: GUARD - this is the windows carrier for internal/speech," \
        "GOOS=$(go env GOOS). On linux the cgo surface is excluded and the load shape" \
        "differs; that reading belongs to scripts/portable-tests.sh, not here." >&2
    exit 2
fi

root=$(cd "$(dirname "$0")/.." && pwd)
cd "$root"

stage=${WISP_SPEECH_STAGE_DIR:-"$root/.scratch/wisp/stage-speech"}
if [ ! -d "$stage" ]; then
    mkdir -p "$stage" || {
        echo "speech-tests.sh: GUARD - cannot create stage dir $stage" >&2
        exit 1
    }
fi

# Build the test binary exactly as "go test" would, but keep it: -o with a
# repo-relative path (the pwd -W PATH trap of wisp-cli-tests.sh does not apply
# to -o, but repo-relative keeps the command readable in logs on both hosts).
echo "speech-tests.sh: building test binary (go test -c ./internal/speech/)"
if ! go test -c -o "$stage/speech.test.exe" ./internal/speech/; then
    echo "speech-tests.sh: go test -c failed; nothing was run" >&2
    exit 1
fi

bash "$root/scripts/stage-sherpa-dlls.sh" "$stage" || exit 1

# Run from the package directory, the cwd "go test" uses, so package-relative
# fixtures resolve the same way. The binary is invoked by its absolute path.
# The -v stream is teed to stdout (see the Output note in the header): the
# full-package runner counts the four numbers off this stream, and the log
# file keeps the same bytes for the skip-reason reading below.
log="$stage/speech-tests.log"
abs_stage=$(cd "$stage" && pwd)
echo "speech-tests.sh: running speech.test.exe (cwd = internal/speech)"
(cd internal/speech && "$abs_stage/speech.test.exe" -test.count=1 -test.timeout=20m -test.v "$@") 2>&1 | tee "$log"
rc=${PIPESTATUS[0]}

passed=$(grep -c '^--- PASS' "$log" || true)
failed=$(grep -c '^--- FAIL' "$log" || true)
skipped=$(grep -c '^--- SKIP' "$log" || true)
ran=$(grep -c '^=== RUN' "$log" || true)
echo "speech-tests.sh: top-level: PASS=$passed FAIL=$failed SKIP=$skipped, === RUN=$ran"
if [ "$failed" -ne 0 ]; then
    echo "speech-tests.sh: failing tests:"
    grep '^--- FAIL' "$log"
fi
if [ "$skipped" -ne 0 ]; then
    echo "speech-tests.sh: note: $skipped skip(s) remain; check their reasons in $log" \
        "(a skip that names onnxruntime 1.17.1 means the shadow won again)"
fi
echo "speech-tests.sh: log kept at $log"
exit "$rc"
