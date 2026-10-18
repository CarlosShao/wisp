#!/usr/bin/env bash
# 342-a1 delivery commits. ⛔ push. Every pathspec literal, in BOTH add and commit.
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 2
unset GOTMPDIR TMPDIR TEMP TMP
R=.scratch/wisp/probes/342/a1
S=/d/tmp/wisp/20261018/342a1
T=.scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md

LOGS=(
"$R/logs/00-h0-open.txt"
"$R/logs/01-h1-rows123.txt"
"$R/logs/02-h2-callers.txt"
"$R/logs/10-host-winsec-baseline-GOTMPDIR-D.keyrows.txt"
"$R/logs/11-host-winsec-baseline-defaultroot.keyrows.txt"
"$R/logs/12-host-probe-windows-rows-run1.log"
"$R/logs/12-host-probe-windows-rows-run2.log"
"$R/logs/13-host-win-forms-run1.log"
"$R/logs/13-host-win-forms-run2.log"
"$R/logs/14-host-win-forms-yi-fixed-and-debug.log"
"$R/logs/15-host-simlanding-jia-win.keyrows.txt"
"$R/logs/15-host-simlanding-yi-win.keyrows.txt"
"$R/logs/20-container-winsec-baseline-1.keyrows.txt"
"$R/logs/21-container-probe-A1A2-run1.log"
"$R/logs/22-container-probe-A1A2-run2.log"
"$R/logs/23-container-forms-run1.log"
"$R/logs/23-container-forms-run2.log"
"$R/logs/23-container-forms-run3.log"
"$R/logs/24-container-A3-unlink-through-ancestor.log"
"$R/logs/24-container-swing-nails-3x.log"
"$R/logs/25-container-simlanding-jia-linux.keyrows.txt"
"$R/logs/25-container-simlanding-yi-linux.keyrows.txt"
"$R/logs/25-driver-g3.summary.txt"
"$R/logs/26-container-simlanding-jia-v2.keyrows.txt"
"$R/logs/27-container-simlanding-yinofb.keyrows.txt"
"$R/logs/27-yinofb-failnames.txt"
"$R/logs/28-container-A2-landed-jia-linux.log"
"$R/logs/28-container-A2-landed-yi-linux.log"
"$R/logs/30-pristine-hashes-treeA-winsec.txt"
"$R/logs/31-landing-diffstats.txt"
"$R/logs/32-config-resample-mut-jia-linux.log"
"$R/logs/32-config-resample-mut-yi-linux.log"
"$R/logs/32-config-resample-treeB.log"
"$R/logs/33-container-crossvet-and-resample.keyrows.txt"
"$R/logs/40-roster-ruler.txt"
"$R/logs/50-bench-hygiene.txt"
"$R/logs/gate-1-gofumpt.txt"
"$R/logs/gate-2-d22scan.txt"
"$R/logs/gate-3-length.txt"
"$R/logs/gate-4-vet.txt"
)
SCRIPTS=("$R/scripts/build-mut.sh" "$R/scripts/driver-g3.sh" "$R/scripts/driver-g4.sh" "$R/scripts/commit-face.sh" "$R/scripts/append-bullet.sh")
PROBE=(
"$R/probe/census342a1_a3_unlink_test.go.txt"
"$R/probe/census342a1_forms_linux_test.go.txt"
"$R/probe/census342a1_posix_race_test.go.txt"
"$R/probe/census342a1_windows_forms_test.go.txt"
"$R/probe/census342a1_windows_rootdebug_test.go.txt"
"$R/probe/census342a1_windows_rows_test.go.txt"
"$R/probe/sim_jia_posix_linux.go.txt"
"$R/probe/sim_jia_win_windows.go.txt"
"$R/probe/sim_yi_posix_linux.go.txt"
"$R/probe/sim_yi_win_windows.go.txt"
)

echo "== shared-index precondition: the only staged things may be my own =="
git status --porcelain -- .scratch/wisp/probes/342 .scratch/wisp/issues | head -60

git add -- "$R/census.md"
git commit -q -F "$S/msg-1.txt" -- "$R/census.md"; echo "commit1 rc=$?"
git add -- "${LOGS[@]}"
git commit -q -F "$S/msg-2.txt" -- "${LOGS[@]}"; echo "commit2 rc=$?"
git add -- "${SCRIPTS[@]}"
git commit -q -F "$S/msg-3.txt" -- "${SCRIPTS[@]}"; echo "commit3 rc=$?"
git add -- "${PROBE[@]}"
git commit -q -F "$S/msg-4.txt" -- "${PROBE[@]}"; echo "commit4 rc=$?"
git add -- "$T"
git commit -q -F "$S/msg-5.txt" -- "$T"; echo "commit5 rc=$?"

echo "== post-commit: working tree must be clean on my faces =="
git status --porcelain -- internal cmd docs .scratch/wisp/probes/342 .scratch/wisp/issues | head -10
git log -5 --format='%h %s' | cat
