#!/usr/bin/env bash
# Corrected sensitivity control for ticket 311 AC#2 (my first pass was an
# instrument failure: GOTMPDIR was pointed at a directory that did not exist,
# so `go` died before the gate could run -- see the voided shots B1/B2).
# Every offender root below is created BEFORE the shot, and the winsec-tier
# roots are built by the ticket-310 block itself (WISP_TEST_GOTMP_ROOT).
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGROOT="$REPO/.scratch/wisp/probes/orch/20261012/logs"
TS=$(date +%Y%m%d-%H%M%S)
OUT="$LOGROOT/orch-311r1b-collect2-$TS.txt"
cd "$REPO" || exit 9
exec 3>"$OUT"
say() { printf '%s\n' "$*" >&3; }

T_ALL=$(mktemp -d); T_ONE=$(mktemp -d)
git archive HEAD | tar -x -C "$T_ALL"; say "T_ALL=$T_ALL go_files=$(find "$T_ALL" -name '*.go' | wc -l)"
git archive HEAD | tar -x -C "$T_ONE"; say "T_ONE=$T_ONE go_files=$(find "$T_ONE" -name '*.go' | wc -l)"
python "$REPO/.scratch/wisp/probes/orch/20261012/txt/strip311.py" \
  "$T_ALL/scripts/winsec-tests.sh" "$T_ALL/scripts/portable-tests.sh" \
  "$T_ALL/scripts/wisp-cli-tests.sh" "$T_ALL/tools/d22scan/runtests.sh" 2>&1 | sed 's/^/strip_all: /' >&3
python "$REPO/.scratch/wisp/probes/orch/20261012/txt/strip311.py" "$T_ONE/scripts/portable-tests.sh" 2>&1 | sed 's/^/strip_one: /' >&3
say "blocks_after_T_ALL(winsec)=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/scripts/winsec-tests.sh" || true) blocks_after_T_ONE(winsec)=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ONE/scripts/winsec-tests.sh" || true) blocks_after_T_ONE(portable)=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ONE/scripts/portable-tests.sh" || true)"

# offender root that really exists, on volume D:
DEXIST="D:/tmp/orch-311collect2-$TS-exists"
mkdir -p "$DEXIST" && say "d_exist_created=$(ls -d "$DEXIST" | wc -l)"

run_shot() {
  local tag="$1"; shift
  local envassign="$1"; shift
  local rootdir="$1"; shift
  local file="$LOGROOT/orch-311r1b-shot2-$tag-$TS.txt"
  { printf 'tag=%s\nworkdir=%s\nenv_assignment=%s\ncommand_line=%s\n' "$tag" "$rootdir" "$envassign" "$*"; } > "$file"
  local rc=0
  if [ -z "$envassign" ]; then ( cd "$rootdir" && "$@" ) >> "$file" 2>&1 || rc=$?
  else ( cd "$rootdir" && env "$envassign" "$@" ) >> "$file" 2>&1 || rc=$?; fi
  { printf 'guard_lines=%s\ninconclusive=%s\nrun_lines=%s\ncleanup_fail=%s\naccess_denied=%s\nwork_dir_error=%s\n' \
      "$(grep -c -- 'AC#1 shape B: GUARD' "$file" || true)" \
      "$(grep -c -- 'AC#1 shape B: INCONCLUSIVE' "$file" || true)" \
      "$(grep -c -- '^=== RUN' "$file" || true)" \
      "$(grep -c -- 'TempDir RemoveAll cleanup' "$file" || true)" \
      "$(grep -c -- 'Access is denied' "$file" || true)" \
      "$(grep -c -- 'creating work dir' "$file" || true)"
    grep -E '^=== RUN=|^ok |^FAIL ' "$file" | tail -3 | sed 's/^/summary: /'
    printf 'rc=%s\n' "$rc"; } >> "$file"
  say "shot=$tag rc=$rc bytes=$(wc -c < "$file") guard=$(grep -c -- 'AC#1 shape B: GUARD' "$file" || true) run=$(grep -c -- '^=== RUN' "$file" || true) cleanup_fail=$(grep -c -- 'TempDir RemoveAll cleanup' "$file" || true) work_dir_error=$(grep -c -- 'creating work dir' "$file" || true) counts=$(grep -E '^=== RUN=' "$file" | tail -1 | tr -d '\r') verdictlines=$(grep -E '^(ok|FAIL)\s' "$file" | tr -d '\r' | tr '\n' '|')"
}

say ""
say "=== (a) re-run: portable refusal, offender root really exists ==="
run_shot "A2b-portable-xvol-existingdir" "GOTMPDIR=$DEXIST" "$REPO" bash scripts/portable-tests.sh ./internal/winsec/

say ""
say "=== (b) sensitivity control, roots built by the ticket-310 block on D: ==="
WROOT_ALL="D:/tmp/orch-311collect2-$TS-b1root"; mkdir -p "$WROOT_ALL"
run_shot "B1b-winsec-xvol-allblocks-stripped" "WISP_TEST_GOTMP_ROOT=$WROOT_ALL" "$T_ALL" bash scripts/winsec-tests.sh ./internal/winsec/
WROOT_ONE="D:/tmp/orch-311collect2-$TS-b2root"; mkdir -p "$WROOT_ONE"
run_shot "B2b-portable-xvol-ownblock-stripped" "WISP_TEST_GOTMP_ROOT=$WROOT_ONE" "$T_ONE" bash scripts/portable-tests.sh ./internal/winsec/
WROOT_P="D:/tmp/orch-311collect2-$TS-b3root"; mkdir -p "$WROOT_P"
run_shot "B3b-winsec-xvol-blocks-present" "WISP_TEST_GOTMP_ROOT=$WROOT_P" "$T_ONE" bash scripts/winsec-tests.sh ./internal/winsec/

say ""
say "=== (c) same roots, guard removed, DEFAULT (C:) location -- colour must stay 58/0 ==="
run_shot "B4b-winsec-default-allblocks-stripped" "" "$T_ALL" bash scripts/winsec-tests.sh ./internal/winsec/
say "ALL DONE2"
