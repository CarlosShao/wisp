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
#
# ===================================================================
# TICKET 111 FIXTURES LEG (2026-11-13) - WHAT CHANGED AND WHY
#
# The CI run that dispatched this change (43389487057 / test-windows step 7)
# died before running one test: the hosted runner has no
# third_party/spike-models (the whole third_party/ tree is git-ignored and
# arrived there only via actions/cache + build.ps1, which fetch the SHERPA DLLS
# but never the voice MODELS), so the voice_env defaults below all missed, and
# the old code answered with a setup-failure exit 1. portable-tests.sh then
# printed its four numbers as RUN=0 and GUARD B judged "package produced no
# top-level result line" - the package-disappeared-is-red alarm (ticket 111
# AC#2/AC#8) doing exactly its job on a run that was a setup failure. The alarm
# is correct; the setup now matches what it assumes.
#
# ROUTE A (taken, plus route B as the named fallback): a new CI step
# (`scripts/fetch-voice-fixtures.ps1`, wired next to stage-sherpa-dlls in
# .github/workflows/ci.yml) downloads the pinned voice fixtures onto the
# runner, so the voice tests RUN there instead of being absent from the
# verdict. The pins and URLs are read from build/models/manifest.json - the
# same single source `wisp models ensure` uses - and every file is sha256
# verified, so the CI shape cannot drift from the manifest (AGENTS.md 1.1's
# "deps.toml-pinned list" discipline, extended to the models).
#
# ROUTE B (the fallback INSIDE this carrier, not a replacement): if a target
# tree is still absent after the fetch step did its thing - a runner with no
# network, a manifest URL gone dead, or a developer machine that never fetched
# - the carrier no longer exits before running anything. It builds and runs
# the package, skipping exactly the fixture-dependent CASES by NAME
# (-test.skip '^Test...|Test...$', full-name anchored) and everything else
# RUNS. The name list is UNIONED from two measurements, both recorded in the
# evidence file: (a) 29 cases whose own body calls a fixture helper, and (b)
# four more the no-fixture reading exposed through shared helpers that reach
# the fixtures INDIRECTLY (Test394R1Seam x2 through the production plan
# resolver, TestSLO159Case14's tense assertions through the 395r1 file, and
# TestAC4Focus... below). Over-skipping leans toward the alarm: a stale entry
# removes coverage and shows up as a shorter RUN count; under-skipping turns
# a setup gap into a red test - the shape this fallback exists to prevent.
# It is NOT a silent skip: the roster is printed here on every fallback run,
# the carrier verifies the pattern still matches the binary's own roster
# (go test -list -run anchored against the compiled binary) and REFUSES to
# ship a pattern that matches zero names - the green-no-op mistake ticket 71
# AC#3 exists to prevent cannot recur silently. portable-tests.sh's
# unaccounted-SKIP rule stays untouched (a ledger row cannot be written from
# outside its file, and the -test.skip mechanism makes the four numbers
# clean), so the names and counts are booked on ticket 111 and in
# docs/evidence/s1/111-cmdwisp-cli-fixtures.md, and the fallback is visible
# in the step log. GUARD B never sees a package without a result line from
# this carrier again.
#
# THE PACKAGE RESULT LINE (new, and the reason GUARD B works at all here):
# a raw test binary run never prints Go's `ok <import path> <duration>` line,
# which is what portable-tests.sh's GUARD B table greps for in mode=cli (the
# speech tier got a speech-shaped GUARD B for exactly this reason; the cli
# tier did not, and portable-tests.sh is not this leg's write surface). The
# carrier therefore prints the line itself, honestly derived: from the
# binary's own measured wall duration and its REAL exit code - `ok` only on
# exit 0, `FAIL` on anything else. There is no path that prints `ok` after a
# setup refusal: every refusal below exits before the run line exists.
#
# NO GUARD WAS TOUCHED: stage-sherpa-dlls refusal, the onnxruntime.dll
# check, runtests.sh's rules, portable-tests.sh's GUARD A/B/C/D, ticket 98's
# load-failure guard - all byte-identical. What changed is only what happens
# AFTER the DLLs are staged, when the voice MODELS are missing.
#
# ===================================================================
# TICKET 111 WAVE 8B (2026-11-13) - THE ROSTER BECOMES COMPUTED, AND A MISSED
# NAME BECOMES ITS OWN VERDICT INSTEAD OF A RED TEST.
#
# The 2026-11-13 evening CI run (43413718734, test-windows step 7) reded one
# case the 41-name roster never listed: Test41a2r1WakeEndToEndArmedToListening
# reached the fixtures through a SECOND-LEVEL wrapper (voice41a2r1Start, whose
# own body calls voice15eRepoFile / voice41a2r1KwsDir / voice15eInstallAsr),
# and the roster was built by grepping helper names against TEST FUNCTION
# BODIES - a case whose body only names the wrapper is invisible to that
# grep, so the fallback skipped everything except exactly the case that then
# failed on the missing tree. The roster also over-skipped one case:
# Test41a2r1GuardRowsWalkWithTheProductionFact never touches the fixtures at
# all (it drives the D43 walk over a stub KWS engine; measured in the
# no-fixture reading: PASS without any fixture on disk), and sat in the
# roster only because its file belongs to the voice family.
#
# THREE ROUTES, ONE NAMED PLACE EACH (the roster is no longer a hand-typed
# 41-name regex; it is computed from the three lists below):
#   (a) fixture HELPER and WRAPPER names - the leaf helpers the bodies call
#       (voice15eVadModel / voice15eSpikeModels / voice41a2r1KwsDir /
#       voiceWirePuncModel / voiceWireTtsFiles / voice15eRepoFile /
#       voice15eInstallAsr) PLUS the second-level constructors that reach the
#       fixtures on the caller's behalf (voice41a2r1Start / voice15eStartLeg /
#       voice15eArm / voiceWireStart / ttsInheritLeg / newVoice26pRig /
#       voice26pRenderTtsWav / voice26pSpikeAsr). The wrapper half is what the
#       CI red exposed; the leaf half is the original route.
#   (b) fixture PATH LITERALS - the tokens only the third_party voice trees
#       contain (spike-models, model-fetch, matcha-icefall, silero_vad,
#       sherpa-onnx-punct, sherpa-onnx-kws, vocos, asr-encoder, test_wavs),
#       so a case (or a wrapper written tomorrow) that reaches the fixtures
#       by spelling a path is caught even with no helper in sight.
#   (c) the four INDIRECT dependents the no-fixture reading exposed and no
#       source grep can see (markseam_394 x2, SLO159, AC4Focus) - kept as the
#       hand-maintained list, reason still in the evidence file.
# The union is evaluated per test function over cmd/wisp/*_test.go on every
# run, so a NEW case that references any listed name lands in the roster the
# first time it compiles - the CI miss above cannot recur through this route.
#
# AND THE GUARD FOR THE ROUTE THIS CANNOT SEE: a case the three routes all
# miss still runs in fallback mode and fails on the missing fixtures - but it
# now fails AS A ROSTER GAP, not as a broken package. After the run, in
# fallback mode only, the carrier re-reads its own log and singles out every
# --- FAIL whose message carries a fixture-absence signature ("cannot find
# third_party" from voice15eRepoFile's walk-up, or "WISP_SPIKE_MODELS is not
# set" from voice15eInstallAsrFiles), prints
# "依赖夹具但不在名册上：<name>" for each, and exits non-zero as FIXTURE-GAP -
# deliberately NOT folded into the plain FAIL count of a verdict, because a
# case that failed on absent fixtures is a roster problem (fix the lists
# above) and not an assertion that broke. GUARD B's职权 is untouched: the
# package result line is still printed first, every run, exactly as before.
#
# ONE MORE SEMANTIC CHANGE, MEASURED AGAINST THE SAME EVENING: voice_env used
# to treat ANY caller-set value as present ("caller values win"). A caller
# that points WISP_SPIKE_MODELS at a directory that does not exist got a full
# run of 30 voice reds and no roster. A caller value is now honoured only if
# the target exists on disk; a caller-set-but-absent target is collected into
# voice_missing exactly like an unset one, so the named fallback (not a wall
# of missing-model reds) is the answer. An existing-but-wrong target still
# wins as before - judging CONTENT was never this script's job, and the
# incomplete-checkpoint reds (voiceTtsFiles.complete / kws.complete) remain
# the tests' own honest failures.
# ===================================================================

set -u

# Refuse the WSL relay before anything else. On this host `bash` on PATH resolves
# to C:\Windows\System32\bash.exe, which is the WSL forwarder and not a shell:
# it dies with "execvpe(/bin/bash) failed: No such file or directory" and runs
# ZERO tests. A leg that invoked it saw an empty log and concluded the carrier
# was broken (measured 2026-11-12 by the resident-ac1-verdict leg, which then
# spent time re-deriving the Git-bash path by hand). Detecting it here turns a
# silent zero-test run into a named setup failure.
#
# The test is what the relay cannot do, not what it is called: ask the shell to
# run a builtin and require the answer. Naming conventions differ across hosts
# (Git Bash, MSYS2, cygwin) and a path check would refuse valid ones.
if [ "$(printf 'ok')" != ok ]; then
	echo "wisp-cli-staged-tests.sh: this shell cannot even run printf; nothing was run" >&2
	exit 2
fi
case "${OSTYPE:-}" in
msys* | cygwin*) : ;;
*)
	echo "wisp-cli-staged-tests.sh: OSTYPE=${OSTYPE:-unset} is not an msys/cygwin shell." >&2
	echo "wisp-cli-staged-tests.sh: on this host 'bash' on PATH is the WSL relay" >&2
	echo "wisp-cli-staged-tests.sh: (C:\\Windows\\System32\\bash.exe), which runs no tests at all." >&2
	echo "wisp-cli-staged-tests.sh: run this through D:\\work\\soft\\Git\\bin\\bash.exe instead." >&2
	exit 2
	;;
esac

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

# ---------------------------------------------------------------------------
# The voice fixtures. The binaries live on disk under third_party/ after
# scripts/fetch-voice-fixtures.ps1 (CI) or `wisp models ensure` (dev host); the
# only thing this script supplies is the environment that points at them.
# Measured 2026-11-12: without these five variables twenty voice tests fail
# with "install the streaming ASR into <tmp>: WISP_SPIKE_MODELS is not set" -
# which reads like a product defect and is not one.
#
# WHY THE DEFAULTS LIVE HERE: docs/evidence/s1/12r2-impl.md:21 records the full
# staging set as prose ("WISP_SILERO_MODEL/WISP_TTS_DIR/WISP_SPIKE_MODELS/
# WISP_KWS_DIR/WISP_PUNC_MODEL/WISP_15E_TMP plus binary-direct cwd"), and
# docs/evidence/s1/15a-r1-impl.md:115 gives the one line for WISP_SPIKE_MODELS.
# Prose cannot be executed, so every leg re-derived it and several got it wrong.
#
# TICKET 111 FIXTURES LEG, THE BEHAVIOUR CHANGE: a target that is absent no
# longer aborts the run. It is COLLECTED into voice_missing, and after all five
# are resolved the run proceeds in exactly one of two modes, both loud:
#
#   fixtures present  - all five env vars exported, the caller's values still
#                       win, nothing skipped, full denominator.
#   fixtures absent   - the fixture-dependent CASES are skipped BY NAME
#                       (-test.skip, full-name anchored, 33 names: 29 whose
#                       own body calls a fixture helper + 4 the no-fixture
#                       reading exposed as indirect dependents, see the
#                       evidence file), everything else runs. Refusing to run
#                       here is what turned a setup failure into "the package
#                       disappeared" on CI (run 43389487057): a package that
#                       runs a REDUCED, NAMED set is a different, honest
#                       verdict - and GUARD B, which rightly refuses a package
#                       that produced NO result, still audits every case this
#                       run does print. WISP_CLI_VOICE_FORCE_SETUP_FAILURE=1
#                       restores the old hard exit for legs that want it.
# ---------------------------------------------------------------------------
voice_missing=""

voice_env() {
	name=$1
	want=$2
	what=$3
	cur=$(eval "printf '%s' \"\${$name:-}\"")
	if [ -n "$cur" ]; then
		# Wave 8b: a caller value is honoured only when the target exists.
		# WISP_SPIKE_MODELS=/nonexistent used to win silently and produce a
		# full run of 30 missing-model reds instead of the named fallback.
		if [ -e "$cur" ]; then
			return 0
		fi
		voice_missing="$voice_missing $name=$cur(caller-set but not on disk: $what)"
		echo "wisp-cli-staged-tests.sh: $name=$cur is SET but not on disk ($what)." >&2
		echo "wisp-cli-staged-tests.sh: -> named-fallback mode: the fixture-dependent cases will be" >&2
		echo "wisp-cli-staged-tests.sh:    skipped by name below, the rest of the package still runs." >&2
		return 0
	fi
	if [ -e "$want" ]; then
		export "$name=$want"
		return 0
	fi
	voice_missing="$voice_missing $name($what)"
	echo "wisp-cli-staged-tests.sh: $name defaults to $want, which is not on disk ($what)." >&2
	echo "wisp-cli-staged-tests.sh: -> named-fallback mode: the fixture-dependent cases will be" >&2
	echo "wisp-cli-staged-tests.sh:    skipped by name below, the rest of the package still runs." >&2
}

voice_env WISP_SPIKE_MODELS "$root/third_party/spike-models" "streaming ASR + VAD tree"
# WHY spike-models AND NOT model-fetch/x-matcha: the two trees are NOT
# interchangeable. voiceTtsFiles.complete() (resident_voice_engines_windows.go:85)
# demands Acoustic, Vocoder, Lexicon, Tokens AND DictDir - the x-matcha tree has
# the acoustic model and the lexicon but no vocoder and no dict directory, so a
# leg pointed there fails with "the matcha checkpoint is incomplete". The
# spike-models subdirectory is the shape the fixtures were written against
# (voiceWireTtsFiles' own fallback at resident_voice_wire_r1_windows_test.go:213
# names exactly this path). Measured 2026-11-12: pointing at x-matcha instead
# took the voice family from 20 failures to 31.
voice_env WISP_SILERO_MODEL "$root/third_party/spike-models/silero_vad.onnx" "silero VAD file"
voice_env WISP_TTS_DIR "$root/third_party/spike-models/matcha-icefall-zh-baker" "matcha TTS checkpoint"
voice_env WISP_KWS_DIR "$root/third_party/spike-models/sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01" "KWS checkpoint"
# WHY THIS IS A FILE AND NOT A DIRECTORY: voiceWirePunctModel
# (resident_voice_wire_r1_windows_test.go:198-205) returns WISP_PUNC_MODEL
# verbatim as the model path and only falls back to a file when it is unset -
# its own fallback names .../model.int8.onnx. Pointing at the x-punc directory
# instead made the leg log "标点级未接上（这一进程的 final 是原始假设，不带标点）"
# and TestWireR1PunctuationReachesTheAgent failed, measured 2026-11-12.
voice_env WISP_PUNC_MODEL "$root/third_party/model-fetch/x-punc/sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8/model.int8.onnx" "punctuation model file"

# ---------------------------------------------------------------------------
# THE FIXTURE-DEPENDENT CASE ROSTER (ticket 111 wave 8b, computed on every run).
# The old hand-typed 41-name regex is gone; the roster is now the UNION of
# three routes, each with its list in exactly one named place. A stale entry
# can no longer hide here either: every name in (a)/(b) is greped against the
# tree while computing, so a rename turns up as a shorter count instead of a
# silently matching regex.
#
# ROUTE (a) - fixture helper and wrapper NAMES a test body can reference.
#   The seven leaf helpers are the functions whose own bodies open the
#   fixture trees (they honour the WISP_* env vars above and fall back to
#   walking up to third_party/). The eight wrappers are the second-level
#   constructors and renderers that reach the fixtures ON THE CALLER'S
#   BEHALF - the route the 2026-11-13 evening CI red exposed
#   (Test41a2r1WakeEndToEndArmedToListening names only voice41a2r1Start,
#   which is why the helper-name-only roster missed it).
#   Members are exactly:
#     voice15eVadModel      - silero VAD file (leaf)
#     voice15eSpikeModels   - the ASR model tree (leaf)
#     voice15eRepoFile      - the walk-up locator every fallback spells (leaf)
#     voice15eInstallAsr    - hardlinks the ASR bytes into the store (leaf)
#     voice41a2r1KwsDir     - the KWS checkpoint directory (leaf)
#     voiceWirePuncModel    - the punctuation model file (leaf)
#     voiceWireTtsFiles     - the matcha TTS checkpoint (leaf)
#     voice41a2r1Start      - builds the full 41a2r1 wake leg (wrapper)
#     voice15eStartLeg      - runs a caller-built voicePlan (wrapper)
#     voice15eArm           - pushes the wake wav through the seam (wrapper)
#     voiceWireStart        - the wire-r1 leg constructor (wrapper)
#     ttsInheritLeg         - the 382wr2 leg constructor (wrapper)
#     newVoice26pRig        - the 26p/392 bargein rig constructor (wrapper)
#     voice26pRenderTtsWav  - renders wav via the real matcha checkpoint (wrapper)
#     voice26pSpikeAsr      - builds the real spike ASR transcriber (wrapper)
#   Deliberately NOT listed: voice15eWaitFor / voice15eWriteWav /
#   voice15eDataDir / voice15eWriteConfig (pure test utilities with no
#   fixture dependency - measured: the cases that use only those PASS with
#   no fixtures on disk) and voice41a2r1Walked / voice41a2r1Feed /
#   voice41a2r1PushSilence / voice41a2r1Wav* (operate on data the case
#   already has; the feed/wav helpers are reached only through the Start
#   constructor above). Over-skipping is the failure mode that silently
#   removes coverage, so the wrapper set is exact names, not a prefix.
#
# ROUTE (b) - fixture PATH LITERALS. These tokens exist only in (or are
#   spelled only for) the git-ignored voice fixture trees, so any body - a
#   test's or a helper's - that contains one reaches the fixtures by path:
#     spike-models, model-fetch, matcha-icefall, silero_vad,
#     sherpa-onnx-punct, sherpa-onnx-kws, vocos, asr-encoder, test_wavs
#   When a fixture tree moves or a new one is added, THIS list is the one
#   place to update (the dispatch that commissioned it names this property).
#
# ROUTE (c) - the four INDIRECT dependents no source grep can see, kept as
#   the hand-maintained list the original leg measured (reasons live in
#   docs/evidence/s1/111-cmdwisp-cli-fixtures.md 4):
#     Test394R1SeamRecordsProductionMark            - drives the production
#        voice plan resolver, which needs the fixture set
#     Test394R1SeamFedRefusalRefusesClosedTaskPoison - same resolver
#     TestSLO159Case14CommentKeepsItsTenseAndItsNumbers - asserts against the
#        395r1 file's logged readings, which are fixture products
#     TestAC4FocusReturnToPriorWindowGap33r5 - its hidden-window foreground
#        sample was only produced inside the voice run
# ---------------------------------------------------------------------------
fixture_helper_pattern='voice15eVadModel|voice15eSpikeModels|voice15eRepoFile|voice15eInstallAsr|voice41a2r1KwsDir|voiceWirePuncModel|voiceWireTtsFiles|voice41a2r1Start|voice15eStartLeg|voice15eArm|voiceWireStart|ttsInheritLeg|newVoice26pRig|voice26pRenderTtsWav|voice26pSpikeAsr'
fixture_literal_pattern='spike-models|model-fetch|matcha-icefall|silero_vad|sherpa-onnx-punct|sherpa-onnx-kws|vocos|asr-encoder|test_wavs'
fixture_roster_indirect="Test394R1SeamRecordsProductionMark Test394R1SeamFedRefusalRefusesClosedTaskPoison TestSLO159Case14CommentKeepsItsTenseAndItsNumbers TestAC4FocusReturnToPriorWindowGap33r5"

# The union itself: per test function over the package's test files, a body
# matching ANY route-(a) name or route-(b) literal joins the roster. One awk
# invocation per file (POSIX awk only - index(), no dynamic regex), flushed
# across file boundaries, deduplicated with sort -u.
fixture_roster() {
	fixture_pattern=$1
	fixture_roster_out=""
	for fixture_file in "$root"/cmd/wisp/*_test.go; do
		[ -f "$fixture_file" ] || continue
		fixture_names=$(awk -v pat="$fixture_pattern" '
			BEGIN { n = split(pat, toks, "|") }
			FNR == 1 { if (name != "" && hit) print name; name = ""; hit = 0 }
			/^func Test[A-Za-z0-9_]+\(/ {
				if (name != "" && hit) print name
				name = $2; sub(/\(.*/, "", name); hit = 0
			}
			name != "" {
				for (i = 1; i <= n; i++) {
					if (toks[i] == "") continue
					if (index($0, toks[i])) { hit = 1; break }
				}
			}
			END { if (name != "" && hit) print name }
		' "$fixture_file")
		fixture_roster_out="$fixture_roster_out$fixture_names
"
	done
	printf '%s' "$fixture_roster_out" | sed '/^$/d' | sort -u
}

skip_args=()
if [ -n "$voice_missing" ]; then
	echo "wisp-cli-staged-tests.sh: missing fixture set:$voice_missing" >&2
	if [ -n "${WISP_CLI_VOICE_FORCE_SETUP_FAILURE:-}" ]; then
		echo "wisp-cli-staged-tests.sh: WISP_CLI_VOICE_FORCE_SETUP_FAILURE=1 -> the pre-2026-11-13" >&2
		echo "wisp-cli-staged-tests.sh: setup-failure exit, requested by the caller. Nothing was run." >&2
		exit 1
	fi
	# Compute the roster from the three named routes above. The computed count
	# is a READING (it moves as tests move); a computed roster of ZERO means
	# the scan broke (no test files found, awk missing) - refusing here is the
	# green-no-op guard, same shape as the -list check below.
	computed_roster=$(fixture_roster "$fixture_helper_pattern
$fixture_literal_pattern")
	computed_count=$(printf '%s\n' "$computed_roster" | grep -c . || true)
	if [ "$computed_count" -eq 0 ]; then
		echo "wisp-cli-staged-tests.sh: GUARD - the computed fixture roster is EMPTY (routes (a)+(b) found" >&2
		echo "wisp-cli-staged-tests.sh:   no test body in cmd/wisp). The scan broke; skipping nothing is the" >&2
		echo "wisp-cli-staged-tests.sh:   green no-op ticket 71 AC#3 forbids. Refusing to run." >&2
		exit 1
	fi
	indirect_count=$(printf '%s\n' $fixture_roster_indirect | grep -c . || true)
	# The skip expression = computed roster + the route-(c) names. The four (c)
	# cases are appended but NOT counted in computed_count: they carry their own
	# hand-maintained accounting (see the route-(c) list above and the printout).
	# NOTE: computed_count + indirect_count is deliberately NOT printed as a
	# sum - markseam_394 x2 sit in BOTH routes, so the honest number is the
	# -test.list verification below, not an arithmetic union of two lists.
	voice_case_pattern="^($(printf '%s\n' "$computed_roster" | paste -sd'|' -)|$(printf '%s' "$fixture_roster_indirect" | tr ' ' '|'))$"
	voice_case_count=$computed_count

	# The pattern MUST match the compiled binary's own roster before it is
	# allowed to skip anything: a -skip that matches zero tests is the
	# green-no-op shape ticket 71 AC#3 was written against. go test -list
	# evaluates the same expression (cheap, no test executes).
	listed=$(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.list "$voice_case_pattern" | grep -c . || true)
	if [ "$listed" -eq 0 ]; then
		echo "wisp-cli-staged-tests.sh: GUARD - the fixture-skip pattern matched ZERO names in the" >&2
		echo "wisp-cli-staged-tests.sh:   compiled test binary (go test -list). The voice cases were renamed" >&2
		echo "wisp-cli-staged-tests.sh:   or deleted, and this pattern is now a green no-op. Refusing to run:" >&2
		echo "wisp-cli-staged-tests.sh:   update the three fixture roster lists alongside the rename, in the" >&2
		echo "wisp-cli-staged-tests.sh:   same commit." >&2
		exit 1
	fi
	skip_args+=("-test.skip=$voice_case_pattern")
	echo "wisp-cli-staged-tests.sh: NAMED FALLBACK - fixtures absent, so $listed fixture-dependent cases" >&2
	echo "wisp-cli-staged-tests.sh:   (computed roster $voice_case_count from routes (a)+(b), plus the" >&2
	echo "wisp-cli-staged-tests.sh:   $indirect_count route-(c) indirect dependents, deduplicated against the" >&2
	echo "wisp-cli-staged-tests.sh:   compiled binary's own roster by go test -list = $listed) are" >&2
	echo "wisp-cli-staged-tests.sh:   skipped BY NAME this run; every other case in ./cmd/wisp/ still runs." >&2
	echo "wisp-cli-staged-tests.sh:   Roster routes: (a) fixture helper/wrapper names, (b) fixture path" >&2
	echo "wisp-cli-staged-tests.sh:   literals, (c) indirect dependents - all three lists live above, in" >&2
	echo "wisp-cli-staged-tests.sh:   THIS script, each in one named place with its rationale." >&2
	echo "wisp-cli-staged-tests.sh:   This is a NAMED skip, not a silent one (ticket 71 AC#3's discipline)." >&2
	echo "wisp-cli-staged-tests.sh:   Fetch the fixtures with" >&2
	echo "wisp-cli-staged-tests.sh:   powershell -File scripts/fetch-voice-fixtures.ps1 to restore the full run." >&2
else
	echo "wisp-cli-staged-tests.sh: voice fixture environment resolved (caller values win)"
fi

# Run from the package directory: the same cwd contract "go test" itself uses,
# and what cmd/wisp's fixture paths are written against.
echo "wisp-cli-staged-tests.sh: running the staged binary from cmd/wisp"
run_start=$(date +%s)
if [ "${#skip_args[@]}" -gt 0 ]; then
	(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.count=1 -test.timeout=90m -test.v "${skip_args[@]}" "$@") 2>&1 | tee "$log"
	bin_rc=${PIPESTATUS[0]}
else
	(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.count=1 -test.timeout=90m -test.v "$@") 2>&1 | tee "$log"
	bin_rc=${PIPESTATUS[0]}
fi
run_end=$(date +%s)
run_dur=$((run_end - run_start))
if [ "$run_dur" -lt 0 ]; then run_dur=0; fi

# The package result line portable-tests.sh's GUARD B (mode=cli table) greps
# for. A raw binary run produces no such line, so the carrier prints it,
# derived from the binary's own exit code and its measured wall duration:
# exit 0 -> `ok`, anything else -> `FAIL`. There is no path that prints `ok`
# after a named setup failure - every refusal above exits before this line
# exists. The import path is the cli tier's pinned one (cli_pin in
# scripts/portable-tests.sh); keep the two spellings in the same commit.
cli_pkg="github.com/CarlosShao/wisp/cmd/wisp"
if [ "$bin_rc" -eq 0 ]; then
	printf 'ok  %s\t%d.00s\n' "$cli_pkg" "$run_dur"
else
	printf 'FAIL  %s\t%d.00s\n' "$cli_pkg" "$run_dur"
fi

if [ -n "$voice_missing" ]; then
	# FIXTURE-GAP GUARD (ticket 111 wave 8b): in fallback mode, a --- FAIL whose
	# message carries a fixture-absence signature is a ROSTER GAP - a case the
	# three routes above all missed, red only because the fixtures are absent -
	# and it must be reported as its own verdict instead of being read as a
	# broken package. The count below is deliberately NOT folded into the plain
	# FAIL accounting anywhere: the four numbers stay honest (the case really
	# ran and really failed), and this block names the GAP on top of them.
	# Skipped cases print nothing, so every --- FAIL in this log is a case that
	# actually ran. Signatures (both measured 2026-11-13, no-fixture reading):
	#   "cannot find third_party"      - voice15eRepoFile's walk-up failure
	#   "WISP_SPIKE_MODELS is not set" - voice15eInstallAsrFiles' guard
	fixture_gap_names=""
	fixture_gap_count=0
	# Attribution shape: in go -v output a test's message lines sit between its
	# own "=== RUN <name>" and its "---- FAIL: <name>", so a signature line is
	# credited to the currently-running case and the verdict is taken from that
	# case's own FAIL line. A FAIL the awk cannot attribute (a subtest-shaped
	# red) simply stays a plain FAIL - the run is still red, only the GAP label
	# is withheld, so the guard's failure mode is under-labeling, never a green.
	while IFS= read -r fixture_gap_name; do
		[ -n "$fixture_gap_name" ] || continue
		fixture_gap_count=$((fixture_gap_count + 1))
		fixture_gap_names="$fixture_gap_names
	依赖夹具但不在名册上：$fixture_gap_name"
	done < <(awk '
		/^=== RUN/ { cur = $3; next }
		/^--- FAIL: / {
			if (cur == $3 && seen[cur]) gap[cur] = 1
			cur = ""
			next
		}
		/^--- PASS: / || /^--- SKIP: / { cur = ""; next }
		cur != "" && (index($0, "cannot find third_party") || index($0, "WISP_SPIKE_MODELS is not set")) {
			seen[cur] = 1
		}
		END { for (n in gap) print n }
	' "$log" | sort)
	if [ "$fixture_gap_count" -gt 0 ]; then
		{
			echo "wisp-cli-staged-tests.sh: FIXTURE-GAP - $fixture_gap_count failed case(s) failed on the fixture-absence" >&2
			echo "wisp-cli-staged-tests.sh:   signature, not on an assertion: this is a ROSTER GAP, not a package defect." >&2
			echo "wisp-cli-staged-tests.sh:   Add these to the route-(a)/(b)/(c) lists above in this same script:$fixture_gap_names" >&2
			echo "wisp-cli-staged-tests.sh:   FIXTURE-GAP carrier verdict: non-zero exit (1); do not read this run as a plain FAIL count." >&2
		} >&2
	fi
	echo "wisp-cli-staged-tests.sh: run completed in NAMED-FALLBACK mode (fixture-dependent cases skipped" >&2
	echo "wisp-cli-staged-tests.sh:   by name; missing:$voice_missing ). Binary exit code: $bin_rc" >&2
	if [ "$fixture_gap_count" -gt 0 ]; then
		# A missed roster name must stay red even if the binary somehow exits 0.
		exit 1
	fi
fi
exit "$bin_rc"
