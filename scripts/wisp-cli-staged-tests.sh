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

# The voice fixtures are NOT downloaded by this script and must not be: the bytes
# are already on disk under third_party/ (spike-models/ and model-fetch/), and the
# only thing that was missing was the environment that points at them. Measured
# 2026-11-12: without these five variables twenty voice tests fail with
# "install the streaming ASR into <tmp>: WISP_SPIKE_MODELS is not set" - which
# reads like a product defect and is not one. Each target is verified before the
# run so a genuinely absent tree is named as a setup failure instead of being
# mistaken for red tests. Variables already set by the caller win, so a leg can
# still point at a different staging tree.
#
# WHY THE DEFAULTS LIVE HERE: docs/evidence/s1/12r2-impl.md:21 records the full
# staging set as prose ("WISP_SILERO_MODEL/WISP_TTS_DIR/WISP_SPIKE_MODELS/
# WISP_KWS_DIR/WISP_PUNC_MODEL/WISP_15E_TMP plus binary-direct cwd"), and
# docs/evidence/s1/15a-r1-impl.md:115 gives the one line for WISP_SPIKE_MODELS.
# Prose cannot be executed, so every leg re-derived it and several got it wrong.
voice_env() {
	name=$1
	want=$2
	cur=$(eval "printf '%s' \"\${$name:-}\"")
	if [ -n "$cur" ]; then
		return 0
	fi
	if [ ! -e "$want" ]; then
		echo "wisp-cli-staged-tests.sh: $name defaults to $want, which is not on disk." >&2
		echo "wisp-cli-staged-tests.sh: the voice fixtures were never fetched on this host, so the" >&2
		echo "wisp-cli-staged-tests.sh: voice tests cannot be judged here. Treat this as a setup" >&2
		echo "wisp-cli-staged-tests.sh: failure, not as red tests." >&2
		exit 1
	fi
	export "$name=$want"
}

voice_env WISP_SPIKE_MODELS "$root/third_party/spike-models"
voice_env WISP_SILERO_MODEL "$root/third_party/spike-models/silero_vad.onnx"
# WHY spike-models AND NOT model-fetch/x-matcha: the two trees are NOT
# interchangeable. voiceTtsFiles.complete() (resident_voice_engines_windows.go:85)
# demands Acoustic, Vocoder, Lexicon, Tokens AND DictDir - the x-matcha tree has
# the acoustic model and the lexicon but no vocoder and no dict directory, so a
# leg pointed there fails with "the matcha checkpoint is incomplete". The
# spike-models subdirectory is the shape the fixtures were written against
# (voiceWireTtsFiles' own fallback at resident_voice_wire_r1_windows_test.go:213
# names exactly this path). Measured 2026-11-12: pointing at x-matcha instead
# took the voice family from 20 failures to 31.
voice_env WISP_TTS_DIR "$root/third_party/spike-models/matcha-icefall-zh-baker"
voice_env WISP_KWS_DIR "$root/third_party/spike-models/sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01"
voice_env WISP_PUNC_MODEL "$root/third_party/model-fetch/x-punc"
echo "wisp-cli-staged-tests.sh: voice fixture environment resolved (caller values win)"

# Run from the package directory: the same cwd contract "go test" itself uses,
# and what cmd/wisp's fixture paths are written against.
echo "wisp-cli-staged-tests.sh: running the staged binary from cmd/wisp"
(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.count=1 -test.timeout=90m -test.v "$@") 2>&1 | tee "$log"
exit "${PIPESTATUS[0]}"
