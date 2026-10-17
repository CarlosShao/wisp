#!/usr/bin/env bash
# 334-r1 AC#2 那把"改前⇄改后三值逐字相同"的尺。
# 同一枚只读探针（同一哈希）分别装进两棵树：
#   改前＝起手锚 03f7e7a6 的字节；改后＝起手锚＋本腿三枚字节（其余字节逐枚验过相等）。
# 台面一律＝CI 那把形。SHAPE334 行只带 class/allow/needL2 三值 ⇒ 两发必须逐字相等；
# REASON334 行带路径，先把临时根削掉再比 ⇒ 两发必须不相等（那正是本票要修的东西）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
BEFORE="${1:?改前树}"; AFTER="${2:?改后树}"
LOGS="$REPO/.scratch/wisp/probes/334/r1/logs"
export GOPATH='D:\work\base\gopath'
PROBE="$LOGS/probe-334r1_test.go.txt"
echo "probe_src_hash=$(git -C "$REPO" hash-object "$PROBE")"
for pair in "before:$BEFORE" "after:$AFTER"; do
    tag="${pair%%:*}"; root="${pair#*:}"
    cp "$PROBE" "$root/internal/risk/zz_probe_334r1_test.go"
    echo "$tag blacklist.go=$(git -C "$REPO" hash-object "$root/internal/risk/blacklist.go") probe=$(git -C "$REPO" hash-object "$root/internal/risk/zz_probe_334r1_test.go") anchor=$(git -C "$REPO" rev-parse 03f7e7a6:internal/risk/blacklist.go)"
done
gotmp="$LOCALAPPDATA\\wisp-test-gotmp\\portable-334r1-shapecmp-$(date +%H%M%S)"
mkdir -p "$(cygpath -u "${gotmp//\\//}")"; export GOTMPDIR="$gotmp"
for pair in "before:$BEFORE" "after:$AFTER"; do
    tag="${pair%%:*}"; root="${pair#*:}"
    ( cd "$root" && go test ./internal/risk/ -count=1 -timeout 30m -v -run 'TestProbe334' ) \
        >"$LOGS/ac2-shapecmp-$tag.log" 2>&1
    echo "$tag rc=$? RUN=$(grep -c '^=== RUN' "$LOGS/ac2-shapecmp-$tag.log" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGS/ac2-shapecmp-$tag.log" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGS/ac2-shapecmp-$tag.log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGS/ac2-shapecmp-$tag.log" || true)"
done
for tag in before after; do
    grep -o 'SHAPE334 [^"]*' "$LOGS/ac2-shapecmp-$tag.log" | sed 's/[[:space:]]*$//' > "$LOGS/ac2-shapecmp-$tag.shape"
done
echo "### 三值尺：改前 ⇄ 改后 SHAPE 行 diff（必须空）"
diff "$LOGS/ac2-shapecmp-before.shape" "$LOGS/ac2-shapecmp-after.shape" && echo "THREE_VALUE_IDENTICAL=YES lines=$(wc -l < "$LOGS/ac2-shapecmp-before.shape")"
echo "### 那句 reason：改前 ⇄ 改后（各三行，取规则尾巴）"
for tag in before after; do
    echo "-- $tag --"
    grep -o 'reason="A-list sensitive path (non-overridable): [^"]*' "$LOGS/ac2-shapecmp-$tag.log" | sed 's/\\//g' | cut -c1-150 | sort -u
done
