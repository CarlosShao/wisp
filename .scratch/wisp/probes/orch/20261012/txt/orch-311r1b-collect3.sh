#!/usr/bin/env bash
# Third pass, orchestrator's own: the AC#4 cells I had not measured myself.
#  T_TRUE_PRE = git archive 87951fd0  (the commit BEFORE leg 311-r1's block landed
#              -- no shape-B block anywhere, so this is the true "改前" side for the
#              winsec tier; 2567e3a9 already carries the salvaged winsec block)
#  WT         = shared worktree @HEAD (only for the "改后" side; the 改前 side is
#              never run on the shared tree -- see A863)
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGROOT="$REPO/.scratch/wisp/probes/orch/20261012/logs"
TS=$(date +%Y%m%d-%H%M%S)
OUT="$LOGROOT/orch-311r1b-collect3-$TS.txt"
cd "$REPO" || exit 9
exec 3>"$OUT"
say() { printf '%s\n' "$*" >&3; }
say "orch-311r1b-collect3 TS=$TS head=$(git rev-parse --short HEAD) true_pre=87951fd0"

T_PRE=$(mktemp -d); git archive 87951fd0 | tar -x -C "$T_PRE"
say "T_PRE=$T_PRE go_files=$(find "$T_PRE" -name '*.go' | wc -l) top=$(ls "$T_PRE" | wc -l) blocks=$(grep -rc -- '^# ---- TICKET 311 AC#1' "$T_PRE/scripts/winsec-tests.sh" 2>/dev/null || echo 0)"

counts() { # counts LABEL FILE
  printf 'label=%s run=%s pass_top=%s fail_top=%s skip_top=%s pass_incl_sub=%s\n' "$1" \
    "$(grep -c -- '^=== RUN' "$2" || true)" \
    "$(grep -c -- '^--- PASS' "$2" || true)" \
    "$(grep -c -- '^--- FAIL' "$2" || true)" \
    "$(grep -c -- '^--- SKIP' "$2" || true)" \
    "$(grep -c -- '^ *--- PASS' "$2" || true)" >&3
}
roster() { grep -E '^--- PASS' "$1" | sed 's/^--- PASS: //; s/ .*//; s/\r//' | sort -u; }

say ""
say "=== winsec tier, TRUE PRE side (T_PRE, default env), two shots ==="
for n in 1 2; do
  f="$LOGROOT/orch-311r1b-shot3-pre-winsec-run$n-$TS.txt"
  { printf 'tag=pre-winsec-run%s\nworkdir=%s\nenv_assignment=(default: neither GOTMPDIR nor WISP_TEST_GOTMP_ROOT set)\ncommand_line=bash scripts/winsec-tests.sh ./internal/winsec/\n' "$n" "$T_PRE"; } > "$f"
  ( cd "$T_PRE" && env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash scripts/winsec-tests.sh ./internal/winsec/ ) >> "$f" 2>&1
  printf 'rc_pre_%s=%s\n' "$n" "$?" >> "$f"
  counts "T_PRE_run$n" "$f"
  say "pre_run$n rc=$(tail -1 "$f" | cut -d= -f2) bytes=$(wc -c < "$f") summary=$(grep -E 'four numbers' "$f" | tail -1 | tr -d '\r' | cut -c1-120)"
done

say ""
say "=== winsec tier, POST side on the worktree @HEAD, two shots ==="
for n in 1 2; do
  f="$LOGROOT/orch-311r1b-shot3-post-winsec-run$n-$TS.txt"
  { printf 'tag=post-winsec-run%s\nworkdir=%s\nenv_assignment=(default)\ncommand_line=bash scripts/winsec-tests.sh ./internal/winsec/\n' "$n" "$REPO"; } > "$f"
  ( env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash scripts/winsec-tests.sh ./internal/winsec/ ) >> "$f" 2>&1
  printf 'rc_post_%s=%s\n' "$n" "$?" >> "$f"
  counts "WT_run$n" "$f"
  say "post_run$n rc=$(tail -1 "$f" | cut -d= -f2) guard=$(grep -c -- 'AC#1 shape B: GUARD' "$f" || true) inconclusive=$(grep -c -- 'AC#1 shape B: INCONCLUSIVE' "$f" || true) bytes=$(wc -c < "$f") summary=$(grep -E 'four numbers' "$f" | tail -1 | tr -d '\r' | cut -c1-120)"
done

say ""
say "=== roster intersection, two rulers side by side (top-level only vs incl. subtests) ==="
roster "$LOGROOT/orch-311r1b-shot3-pre-winsec-run1-$TS.txt"  > "$LOGROOT/orch-311r1b-roster-pre-$TS.txt"
roster "$LOGROOT/orch-311r1b-shot3-pre-winsec-run2-$TS.txt"  >> "$LOGROOT/orch-311r1b-roster-pre-$TS.txt"
roster "$LOGROOT/orch-311r1b-shot3-post-winsec-run1-$TS.txt" > "$LOGROOT/orch-311r1b-roster-post-$TS.txt"
roster "$LOGROOT/orch-311r1b-shot3-post-winsec-run2-$TS.txt" >> "$LOGROOT/orch-311r1b-roster-post-$TS.txt"
sort -u -o "$LOGROOT/orch-311r1b-roster-pre-$TS.txt" "$LOGROOT/orch-311r1b-roster-pre-$TS.txt"
sort -u -o "$LOGROOT/orch-311r1b-roster-post-$TS.txt" "$LOGROOT/orch-311r1b-roster-post-$TS.txt"
say "roster_pre_names=$(wc -l < "$LOGROOT/orch-311r1b-roster-pre-$TS.txt") roster_post_names=$(wc -l < "$LOGROOT/orch-311r1b-roster-post-$TS.txt")"
say "only_in_pre=$(comm -23 "$LOGROOT/orch-311r1b-roster-pre-$TS.txt" "$LOGROOT/orch-311r1b-roster-post-$TS.txt" | grep -c . || true) only_in_post=$(comm -13 "$LOGROOT/orch-311r1b-roster-pre-$TS.txt" "$LOGROOT/orch-311r1b-roster-post-$TS.txt" | grep -c . || true)"
comm -23 "$LOGROOT/orch-311r1b-roster-pre-$TS.txt" "$LOGROOT/orch-311r1b-roster-post-$TS.txt" | sed 's/^/  only_pre: /' >&3
comm -13 "$LOGROOT/orch-311r1b-roster-pre-$TS.txt" "$LOGROOT/orch-311r1b-roster-post-$TS.txt" | sed 's/^/  only_post: /' >&3

say ""
say "=== census totals line, verbatim, before vs after (ticket baseline string in the ticket) ==="
for t in "$T_PRE" "$REPO"; do
  f="$LOGROOT/orch-311r1b-census-$(basename "$t")-$TS.txt"
  ( cd "$t" && bash scripts/portable-tests.sh --scope=census ) > "$f" 2>&1
  say "census tree=$t rc=$? totals=$(grep -E 'packages=[0-9]+ with-zero-compiled-tests' "$f" | tr -d '\r' | tail -1)"
done

say ""
say "=== residue accounting created by my three passes (ticket 310/311 rule: 只建不删) ==="
say "dtmp_orch_dirs=$(ls -d D:/tmp/orch-311collect* 2>/dev/null | wc -l)"
say "dtmp_run_roots=$(find D:/tmp -maxdepth 2 -name 'winsec-202610121*' -o -maxdepth 2 -name 'winsec-20261012-1*' 2>/dev/null | wc -l)"
TEMPD="/c/Users/swq/AppData/Local/Temp"
say "temp_top_Test_count=$(ls -d "$TEMPD"/Test* 2>/dev/null | wc -l)"
say "temp_top_all_count=$(ls "$TEMPD" 2>/dev/null | wc -l)"
say "c_gotmp_root_run_dirs=$(ls -d C:/Users/swq/AppData/Local/wisp-test-gotmp/*-20261012-* 2>/dev/null | wc -l)"
say "d_gotmp_root_dirs=$(ls -d D:/tmp/wisp-test-gotmp/* 2>/dev/null | wc -l)"
say "ALL DONE3"
