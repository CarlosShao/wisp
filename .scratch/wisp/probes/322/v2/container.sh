#!/usr/bin/env bash
# 322-v2 终裁腿容器取证（golang:1.27 linux）——AC#4 的两维现跑：
#   1) internal/winsec 整包 -count=2 -v 四数（票面 AC#4 点名的那把尺）
#   2) AC#4 消费者面 Linux 档：internal/risk/secret/config/agent/memory -count=2 -v
# 台面＝仓外导出树 D:/tmp/wisp/20261018/322v2/tree.v2A（HEAD 字节）；⛔ 共享工作树跑 test。
# 挂载/缓存形状照 322-r2 container.ps1（GOMODCACHE 用宿主 pkg/mod 只读；⛔ 命名卷 ro 坑＝r1b 证据件第 6 节）。
set -u
REPO='/d/work/workspace/projects plans/Wisp'
ROOT_WIN='D:/tmp/wisp/20261018/322v2'
TREE_WIN="$ROOT_WIN/tree.v2A"
OUT="$(cygpath -u "$TREE_WIN")"
LOGS="$REPO/.scratch/wisp/probes/322/v2/logs"
MOD_WIN='D:/work/base/gopath/pkg/mod'
HEAD=$(git -C "$REPO" rev-parse HEAD)
SUM="$LOGS/container-summary.txt"
: > "$SUM"
say() { echo "$@" | tee -a "$SUM"; }

say "### 322-v2 container started $(date '+%F %T %z') HEAD=$HEAD"
say "TREE=$TREE_WIN  MODCACHE=$MOD_WIN"
say "GO-VERSION-IMAGE: $(docker run --rm golang:1.27 go version)"

four() {
  local b="$1"
  printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s' \
    "$(grep -c '^=== RUN' "$b" || true)" \
    "$(grep -c '^--- PASS' "$b" || true)" \
    "$(grep -c '^--- FAIL' "$b" || true)" \
    "$(grep -c '^--- SKIP' "$b" || true)"
}
fault() {
  local b="$1"
  printf 'c0000135=%s buildfailed=%s noteststorun=%s setupfailed=%s' \
    "$(grep -c '0xc0000135' "$b" || true)" \
    "$(grep -c 'build failed' "$b" || true)" \
    "$(grep -c 'no tests to run' "$b" || true)" \
    "$(grep -c 'setup failed' "$b" || true)"
}

crun() {
  local label="$1" pkgs="$2"
  local inner="export GOMODCACHE=/gomod GOTMPDIR=/tmp/gt-$label; mkdir -p \$GOTMPDIR; go test $pkgs -count=2 -timeout 40m -v > /work/$label.log 2>&1; rc=\$?; echo \"go-test-rc=\$rc\" >> /work/$label.log; exit \$rc"
  # ★台面自纠（第一发 rc=125＝零读数，⛔ 计入判据）：Git Bash 把 `-w /work` 做 MSYS 路径改写成
  #   `D:/work/soft/Git/work` ⇒ docker 报 "working directory ... is invalid"。修法＝只对 docker 那一发
  #   关路径转换（MSYS_NO_PATHCONV=1），cp/git 那些 MSYS 形路径照旧走原生转换。
  MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off -e GOFLAGS=-mod=mod \
    -v "$TREE_WIN:/work" -v probe322v2build:/root/.cache/go-build \
    -v "$MOD_WIN:/gomod:ro" -w /work golang:1.27 sh -c "$inner" 2>>"$SUM"
  local rc=$?
  local b="$OUT/$label.log"
  say "=== $label host-rc=$rc"
  if [ -f "$b" ]; then
    cp "$b" "$LOGS/$label.body.log"
    four "$LOGS/$label.body.log" | tee -a "$SUM"
    fault "$LOGS/$label.body.log" | tee -a "$SUM"
    grep -E '^(ok|\?|FAIL)\s' "$LOGS/$label.body.log" | tee -a "$SUM"
    grep '^--- FAIL' "$LOGS/$label.body.log" | tee -a "$SUM"
    grep -oE 'AC#1 paired reading on one screen: attempts=[0-9]+ flips=[0-9]+ flipErrs=[0-9]+ damagedWitnesses=[0-9]+ leafLegRefusals=[0-9]+[^ ]*' "$LOGS/$label.body.log" | tee -a "$SUM"
    tail -1 "$LOGS/$label.body.log" | tee -a "$SUM"
  else
    say "  NO BODY at $b 台面故障"
  fi
}

crun v2winsec './internal/winsec/'
crun v2cons './internal/risk/ ./internal/secret/ ./internal/config/ ./internal/agent/ ./internal/memory/'
say "### done $(date '+%F %T %z')"
