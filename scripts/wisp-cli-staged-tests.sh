#!/usr/bin/env bash
# wisp-cli-staged-tests.sh -- the windows staging carrier for "go test ./cmd/wisp/".
#
# WHY THIS EXISTS (ticket 409's finding, extended to cmd/wisp by the round-87
# whole-package reading): scripts/wisp-cli-tests.sh only puts
# third_party/sherpa-onnx on PATH, and PATH is the WRONG lever for the
# second-level dependency. Windows resolves the whole import tree at PROCESS
# LOAD time, and its search order is
#
#     directory of the executable -> System32 -> ... -> PATH (last)
#
# so a host whose C:\Windows\System32 carries an OS-provided onnxruntime.dll
# (measured on this machine: 1.17.260311, 10,572,960 bytes) binds THAT copy no
# matter what PATH says. The linked sherpa-onnx C API needs 1.28, so every
# voice-backed test in cmd/wisp reports
#
#     onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API
#     requires; a second onnxruntime.dll is shadowing the one on PATH
#
# and either fails or silently degrades. The 2026-11-12 whole-package run read
# 39 such reds out of 52 - the single largest block of red on this host, and
# none of it is a product defect.
#
# internal/speech already solved exactly this in ticket 409 by staging the
# pinned DLLs NEXT TO the test binary and running that binary from the package
# directory. The executable's own directory is the loader's FIRST search
# location, so it fires before System32 - the only lever that works at load
# time. This script is the same carrier for cmd/wisp, and it deliberately reuses
# scripts/stage-sherpa-dlls.sh (the deps.toml-pinned list) instead of naming the
# DLLs again.
#
# WHAT IT IS NOT: it does not relax any assertion, does not change the package
# set, and does not replace scripts/wisp-cli-tests.sh's guards. portable-tests.sh
# still owns the four numbers, the SKIP rule and the ticket-98 load-failure
# GUARD; this script only changes HOW the binary is launched.
#
# Usage: scripts/wisp-cli-staged-tests.sh [extra -test.* flags...]
#   e.g. scripts/wisp-cli-staged-tests.sh -test.run 'Test15e'
# NOTE ON THE FLAG SPELLING: the binary is run DIRECTLY (not through `go test`),
# so its flags already carry the `test.` prefix - `-test.run`, `-test.v`. Passing
# the bare `-run` shape that `go test` accepts fails here with
# "flag provided but not defined: -test", which is the caller's error and is
# reported as such rather than silently swallowed.
# Output: the test log is kept at $WISP_CLI_STAGE_DIR/wisp-cli-tests.log.
# Exit codes: the test binary's exit code, or 1/2 on a named setup failure.

set -u

root=$(cd "$(dirname "$0")/.." && pwd)
stage=${WISP_CLI_STAGE_DIR:-"$root/.scratch/wisp/stage-cli"}

goos=$(go env GOOS 2>/dev/null || true)
if [ "$goos" != windows ]; then
	echo "wisp-cli-staged-tests.sh: this carrier exists for the windows loader problem" >&2
	echo "wisp-cli-staged-tests.sh: (System32 onnxruntime shadowing), and GOOS=$goos does not have it." >&2
	echo "wisp-cli-staged-tests.sh: refusing to run a diluted scope on the wrong platform." >&2
	exit 2
fi

if ! command -v go >/dev/null 2>&1; then
	echo "wisp-cli-staged-tests.sh: go is not on PATH; nothing was run" >&2
	exit 2
fi

mkdir -p "$stage" || exit 2
abs_stage=$(cd "$stage" && pwd) || exit 2
log="$abs_stage/wisp-cli-tests.log"

# Build first. A build failure is a setup failure, not a verdict on the package.
echo "wisp-cli-staged-tests.sh: building test binary (go test -c ./cmd/wisp/)"
if ! go test -c -o "$abs_stage/wisp.test.exe" ./cmd/wisp/; then
	echo "wisp-cli-staged-tests.sh: go test -c failed; nothing was run" >&2
	exit 1
fi

# Stage the pinned DLLs next to the binary. The shared script reads deps.toml, so
# a pin added there is required here automatically and a missing file is loud.
echo "wisp-cli-staged-tests.sh: staging pinned sherpa DLLs next to the test binary"
if ! bash "$root/scripts/stage-sherpa-dlls.sh" "$abs_stage"; then
	echo "wisp-cli-staged-tests.sh: staging the pinned DLLs failed; nothing was run" >&2
	exit 1
fi

# Verify the staging actually put an onnxruntime next to the binary before
# claiming the shape works: a silently empty stage directory would reproduce the
# very shadowing this exists to defeat, and the tests would look "red for other
# reasons" instead of "not measured".
if [ ! -f "$abs_stage/onnxruntime.dll" ]; then
	echo "wisp-cli-staged-tests.sh: no onnxruntime.dll next to the test binary after staging." >&2
	echo "wisp-cli-staged-tests.sh: the loader would fall back to System32 and the run would not" >&2
	echo "wisp-cli-staged-tests.sh: measure the product. Treat this as a setup failure." >&2
	exit 1
fi

# Run from the package directory: the same cwd contract "go test" itself uses,
# and what cmd/wisp's fixture paths are written against.
echo "wisp-cli-staged-tests.sh: running the staged binary from cmd/wisp"
(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.count=1 -test.timeout=90m -test.v "$@") 2>&1 | tee "$log"
exit "${PIPESTATUS[0]}"
