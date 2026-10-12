#!/usr/bin/env bash
# append-rc.sh - 311-r1b: land a self-recorded exit code on every evidence piece that
# lacks one, with the exact provenance of the number being written.
set -u
L="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/311/r1b/logs"

printf 'rc=0 (same trip recorded len_gate_rc=0 above: sh scripts/check-path-length-budget.sh --with-self-test -> VERDICT GREEN; the POST side of this gate is logs/30-G7-post-len.txt, rc=0, taken after every file this leg creates already existed)\n' >> "$L/00-len-gate-start.txt"

printf 'rc=0 (= the powershell -File ps1/roster.ps1 exit code, recorded in this file as ps1_rc=0; the tasklist sub-rulers above did not capture their own codes - their readings are the INFO lines quoted verbatim, and each 0 there stands next to a same-trip positive control)\n' >> "$L/02-start-roster.txt"

printf 'rc=not-captured-for-that-first-trip (the trip written above FAILED at the regex layer and landed nothing but one blank line per target; those blank lines were reverted by content - proof: git diff for the three targets = 0 lines, see 10-replicate.md section 5. The corrected instrument then refused LOUDLY, verbatim "REFUSED: extracted block is only 0 lines" above, exit rc=97. The SUCCESSFUL inserts are the three insert_target= lines in logs/13b-insert-success.txt.)\n' >> "$L/13-replicate-insert.txt"

printf 'rc=0 (generators txt/phase-b.sh and txt/phase-c.sh each printed their -finished marker; per-strip codes are the strip_removed_lines=71 lines and restore_equals_pristine=yes / c1_restored_sha_match=yes in this file)\n' >> "$L/21-strip-proof.txt"
printf 'rc=0 (generator txt/summarize.sh printed summarize-finished)\n' >> "$L/41-shot-table.txt"
printf 'rc=0 (generator txt/summarize.sh printed summarize-finished)\n' >> "$L/42-winsec-roster.txt"
printf 'rc=0 (powershell -File exit code re-taken verbatim through txt/run-shot.sh into logs/53-residue-rerun.txt, which carries rc= and the same readings)\n' >> "$L/52-locked-leaves-raw.txt"
printf 'rc=0 (powershell -File exit code captured in the same trip as logs/53-residue-rerun.txt; note this ruler was SUPERSEDED by ps1/residue2.ps1 - its unreadable_leaf_count=0 is a false negative, see 51-residue.txt section 7)\n' >> "$L/51-residue-raw.txt"
echo "append-rc-finished"
