#!/usr/bin/env bash
# 322-v2 容器成对（fresh，我自己的两棵树）：head(v2B=c11163c8 pristine) ⇄ mine(v2A=HEAD 含 r2 那枚 .go)。
# 目的＝独立复认 AC#4 的"容器成对 142/86/0/0×两台"那一维，⛔ 只靠数腿的旧 log。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT_WIN='D:/tmp/wisp/20261018/322v2'
LOGS="$REPO/.scratch/wisp/probes/322/v2/logs"
MOD_WIN='D:/work/base/gopath/pkg/mod'
SUM="$LOGS/pair-fresh-summary.txt"
: > "$SUM"
say() { echo "$@" | tee -a "$SUM"; }
run1() {
  local label="$1" treewin="$2"
  local inner="export GOMODCACHE=/gomod GOTMPDIR=/tmp/gt-$label; mkdir -p \$GOTMPDIR; go test ./internal/winsec/ -count=2 -timeout 40m -v > /work/$label.log 2>&1; rc=\$?; echo \"go-test-rc=\$rc\" >> /work/$label.log; exit \$rc"
  MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
    -v "$treewin:/work" -v probe322v2build:/root/.cache/go-build \
    -v "$MOD_WIN:/gomod:ro" -w /work golang:1.27 sh -c "$inner" 2>>"$SUM"
  local rc=$?
  local b="$(cygpath -u "$treewin")/$label.log"
  say "=== $label host-rc=$rc"
  if [ -f "$b" ]; then
    cp "$b" "$LOGS/$label.body.log"
    local bl="$LOGS/$label.body.log"
    printf '%s ' "$label"; printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s ' \
      "$(grep -c '^=== RUN' "$bl"||true)" "$(grep -c '^--- PASS' "$bl"||true)" \
      "$(grep -c '^--- FAIL' "$bl"||true)" "$(grep -c '^--- SKIP' "$bl"||true)" | tee -a "$SUM"; echo
    printf 'fault: c0000135=%s buildfailed=%s noteststorun=%s setupfailed=%s\n' \
      "$(grep -c '0xc0000135' "$bl"||true)" "$(grep -c 'build failed' "$bl"||true)" \
      "$(grep -c 'no tests to run' "$bl"||true)" "$(grep -c 'setup failed' "$bl"||true)" | tee -a "$SUM"
    grep -E '^(ok|\?|FAIL)\s' "$bl" | tee -a "$SUM"
    grep -oE 'AC#1 paired reading on one screen: attempts=[0-9]+ flips=[0-9]+ flipErrs=[0-9]+ damagedWitnesses=[0-9]+ leafLegRefusals=[0-9]+' "$bl" | tee -a "$SUM"
  else say "  NO BODY = 台面故障"; fi
}
say "### 322-v2 fresh pair $(date '+%F %T %z')"
run1 v2head "$ROOT_WIN/tree.v2B"
run1 v2mine "$ROOT_WIN/tree.v2A"
say "### done $(date '+%F %T %z')"
