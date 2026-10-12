#!/usr/bin/env bash
# Orchestrator's own re-runs for leg 311-r1b (ticket 311 AC#2 teeth + AC#4 gates).
# Nothing here mutates a tracked file: the strip control runs only in out-of-repo
# export trees made by `git archive HEAD`.
set -u
REPO="D:/work/workspace/projects plans/Wisp"
LOGROOT="$REPO/.scratch/wisp/probes/orch/20261012/logs"
TS=$(date +%Y%m%d-%H%M%S)
OUT="$LOGROOT/orch-311r1b-collect-$TS.txt"
PY="$REPO/.scratch/wisp/probes/orch/20261012/txt/strip311.py"
cd "$REPO" || exit 9

exec 3>"$OUT"
say() { printf '%s\n' "$*" >&3; }
say "orch-311r1b-collect TS=$TS head=$(git rev-parse --short HEAD) localappdata=$LOCALAPPDATA"

# ---- STEP 0: static shape of the four blocks, HEAD object layer -----------
say ""
say "=== STEP 0 static block census (provenance = HEAD) ==="
for f in scripts/winsec-tests.sh scripts/portable-tests.sh scripts/wisp-cli-tests.sh tools/d22scan/runtests.sh; do
  say "file=$f $(git grep -c -- 'ticket 311 AC#1 shape B' HEAD -- "$f" || echo 'HEAD:0') $(git grep -c -- '^# ---- TICKET 311 AC#1' HEAD -- "$f" || echo 'HEAD:0')"
done
say "worktree_dirty_lines_vs_HEAD_scripts_tools=$(git diff --numstat HEAD -- scripts tools | wc -l)"

XVOL="D:/tmp/orch-311collect-$TS"
mkdir -p "$XVOL" 2>/dev/null || true

# run_shot TAG ENVASSIGN ROOT CMD...
run_shot() {
  local tag="$1"; shift
  local envassign="$1"; shift
  local rootdir="$1"; shift
  local file="$LOGROOT/orch-311r1b-shot-$tag-$TS.txt"
  {
    printf 'tag=%s\nworkdir=%s\nenv_assignment=%s\ncommand_line=%s\n' "$tag" "$rootdir" "$envassign" "$*"
  } > "$file"
  local rc=0
  if [ -z "$envassign" ]; then
    ( cd "$rootdir" && "$@" ) >> "$file" 2>&1 || rc=$?
  else
    ( cd "$rootdir" && env "$envassign" "$@" ) >> "$file" 2>&1 || rc=$?
  fi
  {
    printf 'guard_lines=%s\ninconclusive_lines=%s\nrun_lines=%s\ntempdir_cleanup=%s\naccess_denied=%s\npass_top=%s\nfail_top=%s\n' \
      "$(grep -c -- 'AC#1 shape B: GUARD' "$file" || true)" \
      "$(grep -c -- 'AC#1 shape B: INCONCLUSIVE' "$file" || true)" \
      "$(grep -c -- '^=== RUN' "$file" || true)" \
      "$(grep -c -- 'TempDir RemoveAll cleanup' "$file" || true)" \
      "$(grep -c -- 'Access is denied' "$file" || true)" \
      "$(grep -c -- '^--- PASS' "$file" || true)" \
      "$(grep -c -- '^--- FAIL' "$file" || true)"
    printf 'rc=%s\n' "$rc"
  } >> "$file"
  say "shot=$tag rc=$rc bytes=$(wc -c < "$file") guard=$(grep -c -- 'AC#1 shape B: GUARD' "$file" || true) run=$(grep -c -- '^=== RUN' "$file" || true) passdash=$(grep -c -- '^--- PASS' "$file" || true) faildash=$(grep -c -- '^--- FAIL' "$file" || true) cleanup=$(grep -c -- 'TempDir RemoveAll cleanup' "$file" || true) summary=$(grep -E '^=== RUN=|^ *(--- )?(PASS|FAIL)' "$file" | tail -1 | tr -d '\r')"
}

say ""
say "=== STEP 1 (a) cross-volume refusal, offender root = $XVOL/* ==="
run_shot "A1-winsec"   "GOTMPDIR=$XVOL/a1" "$REPO" bash scripts/winsec-tests.sh ./internal/winsec/
run_shot "A2-portable" "GOTMPDIR=$XVOL/a2" "$REPO" bash scripts/portable-tests.sh ./internal/winsec/
run_shot "A3-cli"      "GOTMPDIR=$XVOL/a3" "$REPO" bash scripts/wisp-cli-tests.sh
run_shot "A4-runtests" "GOTMPDIR=$XVOL/a4" "$REPO" sh tools/d22scan/runtests.sh -C tools/d22scan ./...
run_shot "A5-runtests-gotmpdir-only" "GOTMPDIR=$XVOL/a5" "$REPO" sh tools/d22scan/runtests.sh -C tools/d22scan ./...

say ""
say "=== STEP 2 export trees (git archive HEAD; worktree untouched) ==="
T_ALL=$(mktemp -d) || exit 8
T_ONE=$(mktemp -d) || exit 8
git archive HEAD | tar -x -C "$T_ALL"; say "archive_all_rc=$?"
git archive HEAD | tar -x -C "$T_ONE"; say "archive_one_rc=$?"
say "T_ALL=$T_ALL go_files=$(find "$T_ALL" -name '*.go' | wc -l) top=$(ls "$T_ALL" | wc -l)"
say "T_ONE=$T_ONE go_files=$(find "$T_ONE" -name '*.go' | wc -l) top=$(ls "$T_ONE" | wc -l)"
for f in scripts/winsec-tests.sh scripts/portable-tests.sh tools/d22scan/runtests.sh; do
  say "preop tree_all $f crlf=$(grep -c "$(printf '\r')" "$T_ALL/$f" || true) blocks=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/$f" || true)"
done
python "$PY" "$T_ALL/scripts/winsec-tests.sh" "$T_ALL/scripts/portable-tests.sh" "$T_ALL/scripts/wisp-cli-tests.sh" "$T_ALL/tools/d22scan/runtests.sh" 2>&1 | sed 's/^/strip_all: /' >&3
python "$PY" "$T_ONE/scripts/portable-tests.sh" 2>&1 | sed 's/^/strip_one: /' >&3
say "poststrip_T_ALL_blocks(winsec,portable,cli,runtests)=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/scripts/winsec-tests.sh" || true),$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/scripts/portable-tests.sh" || true),$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/scripts/wisp-cli-tests.sh" || true),$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ALL/tools/d22scan/runtests.sh" || true)"
say "poststrip_T_ONE_portable=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ONE/scripts/portable-tests.sh" || true) T_ONE_winsec_untouched=$(grep -c -- '^# ---- TICKET 311 AC#1' "$T_ONE/scripts/winsec-tests.sh" || true)"
say "worktree_winsec_still_clean=$(git diff --numstat HEAD -- scripts/winsec-tests.sh | wc -l) sha_worktree=$(sha256sum scripts/winsec-tests.sh | cut -c1-16) sha_T_ALL_stripped=$(sha256sum "$T_ALL/scripts/winsec-tests.sh" | cut -c1-16)"

say ""
say "=== STEP 3 (b) sensitivity control: guard removed => the colour must come back ==="
run_shot "B1-winsec-xvol-allblocks-stripped" "GOTMPDIR=$XVOL/b1" "$T_ALL" bash scripts/winsec-tests.sh ./internal/winsec/
run_shot "B2-portable-xvol-ownblock-stripped" "GOTMPDIR=$XVOL/b2" "$T_ONE" bash scripts/portable-tests.sh ./internal/winsec/

say ""
say "=== STEP 4 (c) shapes the guard must not break (present, on HEAD worktree) ==="
run_shot "C1-winsec-default-env" "" "$REPO" bash scripts/winsec-tests.sh ./internal/winsec/
run_shot "C2-winsec-samevol-other-dir" "WISP_TEST_GOTMP_ROOT=C:/Users/swq/AppData/Local/orch-311collect-$TS-c2" "$REPO" bash scripts/winsec-tests.sh ./internal/winsec/

say ""
say "=== STEP 5 AC#4 gates ==="
GOFLAGS= go build ./... > "$LOGROOT/orch-311r1b-build-$TS.txt" 2>&1; say "build_rc=$? bytes=$(wc -c < "$LOGROOT/orch-311r1b-build-$TS.txt")"
sh scripts/d22scan.sh > "$LOGROOT/orch-311r1b-d22scan-$TS.txt" 2>&1; say "d22scan_rc=$? tail=$(tail -3 "$LOGROOT/orch-311r1b-d22scan-$TS.txt" | tr -d '\r' | tr '\n' '|')"
GP=$(go env GOPATH)
"$GP/bin/gofumpt.exe" -l internal/winsec tools/d22scan scripts > "$LOGROOT/orch-311r1b-gofumpt-$TS.txt" 2>&1; say "gofumpt_rc=$? dirty=$(grep -c . "$LOGROOT/orch-311r1b-gofumpt-$TS.txt" || true) list=$(tr '\n' ',' < "$LOGROOT/orch-311r1b-gofumpt-$TS.txt")"
say "tasklist_wisp=$(tasklist //FI "IMAGENAME eq wisp.exe" | grep -c wisp || true) tasklist_balldebug=$(tasklist //FI "IMAGENAME eq balldebug.exe" | grep -c balldebug || true)"

say ""
say "=== STEP 6 per-commit over-reach ruler (leg 311-r1b's eight commits) ==="
for c in e4e787fb f8b1f446 10734fbd 0f888a4f 501af017 3b553a2f 07e7a1e0 8f350a9d; do
  say "commit=$c files=$(git show --name-only --format= "$c" | grep -c . || true) forbidden=$(git show --name-only --format= "$c" | grep -c -E '^(frontend/|design/|PLAN\.md|docs/specs/|internal/observe/thresholds\.go|tools/d22scan/allowlist\.txt|internal/panel/tokens_fourway_test\.go|internal/panel/l2_grant_boundary_test\.go|internal/perm/ticket90_persist_test\.go)' || true) nonmd=$(git show --name-only --format= "$c" | grep -c -v '\.md$' || true)"
  git show --name-only --format= "$c" | grep -v '\.md$' | sed 's/^/    nonmd: /' >&3
done
say "final_worktree_dirty_scripts_tools=$(git status --porcelain -- scripts tools | wc -l)"
say "ALL DONE"
