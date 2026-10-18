#!/bin/bash
# 339-r1 AC#1 丙形现量：redlist 行的机械前提＝ledger 名字必须在“这个平台编译出的测试二进制”里
# 存在（portable-tests.sh 的 staleness GUARD：go test -list 点名，查不到＝exit 1）。
# 所以“先落甲、之后想补丙行”＝那一行在 linux 档直接 stale＝GUARD 红。本尺量这件事，⛔ 推理。
# 正控＝同一把尺打在甲形前的基线树（两名字必须读 1）；读 0＝尺坏（第一发就坏过：嵌套引号
# 把 ^$n$ 打成 ^'$n'$，具名作废＝本件上一版，新尺重打）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
OUT="$LOGD/bing-staleness-pair.log"
: > "$OUT"
probe_tree() { # $1=tree $2=标签
  echo "## label=$2 tree=$1" >> "$OUT"
  MSYS_NO_PATHCONV=1 docker run --rm -v "$1":/src -v "$(cd "$REPO" && go env GOMODCACHE)":/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c 'mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt
      for n in TestAnchorMissReasonNamesCandidateAndMissingAnchor TestAnchorMissKeepsNonOverridableClassAndNeedL2 TestAnchorMissHasNoWorldOnThisPlatform339; do
        c=$(go test -list "^$n$" ./internal/risk/ 2>/dev/null | grep -cx "$n" || true)
        echo "  name=$n in_linux_test_binary=$c"
      done' >> "$OUT" 2>&1
  echo "== $2 done rc=$?" >> "$OUT"
}
echo "# 丙形 staleness 成对尺：go test -list 点名，甲前（基线树，两枚必须在＝正控）⇄ 甲后（final 树，两枚随档下线）。" >> "$OUT"
probe_tree "/d/tmp/wisp/20261018/339r1/linux/tree" "before-jia"
probe_tree "/d/tmp/wisp/20261018/339r1/linux/final" "after-jia"
cat "$OUT"
