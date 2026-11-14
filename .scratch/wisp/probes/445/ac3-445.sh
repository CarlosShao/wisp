#!/bin/sh
# Ticket 445 AC#3 -- "the tooth is still there", measured on an OUT-OF-REPO
# EXPORT TREE. The main repository is never mutated here: every mutation is
# applied to $TREE (a copy of the working tree with third_party junctioned back
# to the repo, so the fixtures and the pinned DLLs are the real ones).
#
# WHAT IT SHOWS, in order:
#   01 baseline      the repaired case is green on the exported copy (same bytes
#                    as the repo), so the readings below are about the mutation.
#   02 m1            TriggerInterrupt REALLY not called (the click's gesture keeps
#                    running, the executor is still attached, the call is gone)
#                    => the repaired case must red, naming the MUTANT line.
#   03 restored      the same bytes back (sha256 printed) => green again.
#   04 m2            the walk still happens but the D43 edge is never recorded
#                    => the repaired case's WAIT must red with its own named line
#                    and numbers. A "wait longer" fix cannot pass this: the edge
#                    never lands.
#   05 restore       both files back to their pristine bytes, sha256 verified.
#
# Usage: sh .scratch/wisp/probes/445/ac3-445.sh [tree]
set -u

tree=${1:-/d/tmp/wisp/20261114/445-tree}
here=$(cd "$(dirname "$0")" && pwd)
repo=$(cd "$here/../../../.." && pwd)
logs="$here/logs"
[ -d "$tree" ] || { echo "ac3-445: no export tree at $tree" >&2; exit 2; }
[ -f "$tree/go.mod" ] || { echo "ac3-445: $tree has no go.mod" >&2; exit 2; }
mkdir -p "$logs"

# The repaired instruments, byte-for-byte from the repo, so the exported copy
# judges the same code the commit ships.
for f in resident_voice_interrupt_390r1_windows_test.go \
	resident_voice_interrupt_381r2_windows_test.go \
	resident_voice_wire_r1_windows_test.go \
	resident_voice_acting_395r1_windows_test.go \
	voice_e2e_demo_16r1_windows_test.go; do
	cp "$repo/cmd/wisp/$f" "$tree/cmd/wisp/$f" || exit 2
done

run() { # run <label>
	label=$1
	log="$logs/$label.log"
	(cd "$tree" && WISP_445_ROUNDS=20 scripts/wisp-cli-staged-tests.sh \
		-test.run '^Test390r1BallClickPlayingInterruptAbortsAndWalks27$') \
		>"$log" 2>&1
	rc=$?
	echo "== $label rc=$rc log=$log"
	grep -E 'MUTANT|the D43 walk never landed|the click.s stop walked|CLICK-27|^--- (PASS|FAIL)' "$log" | sed 's/^/   /'
	return $rc
}

# The system `python` on this host is a uv trampoline that cannot spawn a child
# from inside a bash tool call ("uv trampoline failed to spawn Python child
# process" - measured on the first run of this script), so the bundled
# interpreter is named explicitly. A mutation that did not apply must STOP the
# run: a green "mutant" is the worst possible reading.
PY=${PY445:-/c/Users/swq/.dsh/dsh-runtimes/dsh-primary-runtime/dependencies/python/python.exe}

echo "== production sha256 in the export tree before any mutation"
sha256sum "$tree/cmd/wisp/resident_ball_windows.go" "$tree/cmd/wisp/resident_voice_windows.go" | tee "$logs/00-pristine-sha.txt"

run 01-baseline-green

"$PY" "$here/mutate-445.py" "$tree" m1 | tee "$logs/02-m1-mutation.txt" || {
	echo "ac3-445: the m1 mutation did not apply; nothing was measured" >&2; exit 4; }
run 02-m1-trigger-never-called

"$PY" "$here/mutate-445.py" "$tree" restore | tee "$logs/03-restore-sha.txt" || exit 4
run 03-restored-green

"$PY" "$here/mutate-445.py" "$tree" m2 | tee "$logs/04-m2-mutation.txt" || {
	echo "ac3-445: the m2 mutation did not apply; nothing was measured" >&2; exit 4; }
run 04-m2-edge-never-recorded

"$PY" "$here/mutate-445.py" "$tree" restore | tee "$logs/05-restore-sha.txt" || exit 4
sha256sum "$tree/cmd/wisp/resident_ball_windows.go" "$tree/cmd/wisp/resident_voice_windows.go" | tee "$logs/05-restored-sha.txt"

echo "== AC#3 done: compare 00-pristine-sha.txt with 05-restored-sha.txt"
