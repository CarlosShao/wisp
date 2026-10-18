#!/usr/bin/env bash
# 336-v1 AC#5/AC#6 discipline rulers (judge-run, read-only on named commits + disk census).
# Section-sign needle built from bytes so the ruler itself never contains U+00A7.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
L="${LOGDIR:-/d/tmp/wisp/20261018/336v1/logs}"
SEC=$(printf '\302\247')
{
echo "### A. leg artefacts: tracked at e5021339 vs on disk now (HEAD drifts; tracked is authoritative)"
echo "tracked_at_e5021339=$(git ls-tree -r --name-only e5021339 -- .scratch/wisp/probes/336/r1 | wc -l)"
echo "on_disk_now=$(find .scratch/wisp/probes/336/r1 -type f | wc -l)"
echo "### B. zero-byte files ON DISK in the leg write face (the 6 the leg self-reported)"
find .scratch/wisp/probes/336/r1 -type f -size 0 -printf "%s\t%p\n" | sort -k2 > "$L/ac6-r1-zerobyte-disk.txt"
echo "zerobyte_disk_count=$(wc -l < "$L/ac6-r1-zerobyte-disk.txt")"; cat "$L/ac6-r1-zerobyte-disk.txt"
echo "### C. any of those TRACKED (that would be the real 0-byte=没交 violation)? want 0"
git ls-tree -r -l e5021339 -- .scratch/wisp/probes/336/r1 | awk '$5==0 {print}' > "$L/ac6-r1-zerobyte-tracked.txt"
echo "zerobyte_tracked_count=$(grep -c '' "$L/ac6-r1-zerobyte-tracked.txt" || true)"
echo "  positive control for the same awk ruler (size column present for all): total=$(git ls-tree -r -l e5021339 -- .scratch/wisp/probes/336/r1 | wc -l)"
echo "### D. section-sign census in the leg's report + logs (needle built from bytes)"
printf "sec_in_r1_report=%s\n" "$(grep -c "$SEC" .scratch/wisp/probes/336/r1/report.md || true)"
printf "sec_files_in_r1_logs=%s\n" "$(grep -rl "$SEC" .scratch/wisp/probes/336/r1/logs 2>/dev/null | wc -l || true)"
printf "  positive control: sec_files_found_when_grepping_a_token_that_exists=%s\n" "$(grep -rl "TEETH336" .scratch/wisp/probes/336/r1/logs 2>/dev/null | wc -l || true)"
echo "### E. section-sign in MY write face (want 0 files)"
find .scratch/wisp/probes/336/v1 -type f -print0 | xargs -0 grep -l "$SEC" 2>/dev/null > "$L/ac6-v1-section-sign.txt" || true
echo "my_sec_files=$(grep -c '' "$L/ac6-v1-section-sign.txt" || true)"
echo "### F. longest path in my write face (char count)"
find .scratch/wisp/probes/336/v1 -type f -printf "%p\n" | awk '{ print length($0) }' | sort -n | tail -1
echo "### G. AC-box flip census on ticket 336 (judgment area only, lines 20-26) — want unchecked 7 / checked 0"
echo "unchecked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[ \]' || true)"
echo "checked=$(sed -n '20,26p' .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | grep -c '^[[:space:]]*- \[x\]' || true)"
echo "### H. ticket name / Status line untouched by me: diff vs e5021339 for lines 1-2"
diff <(git show e5021339:.scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md | head -2) <(head -2 .scratch/wisp/issues/336-a-tier-fallback-branch-per-anchor-teeth-334-ac2-remainder.md) && echo "  lines1-2=IDENTICAL"
echo "ac6_rc=0"
} 2>&1 | tee "$L/ac6-discipline.log"
