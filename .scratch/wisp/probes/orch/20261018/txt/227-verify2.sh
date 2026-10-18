#!/bin/sh
# 编排者收件复跑 227-r1 第二步：① 写面尺（ruler 1）② 那枚决定性突变我自己重打（:393 的 fixture 删 schema_version ⇒ 必红）
set -u
REPO="/d/work/workspace/projects plans/Wisp"
W="/d/tmp/wisp/20261018/orch/227-verify"
ANCHOR=bcf90e06
cd "$REPO" || exit 9
{
echo "=== [A] ruler 1：逐枚 name-only 越界尺（严格 subject 前缀） ==="
CS=$(git log --all --grep='^票 227 落地腿 227-r1' --format='%h')
echo "commits=$(echo $CS | wc -w)"
for c in $CS; do git show --name-only --format='' "$c"; done | grep -v '^$' | sort -u > /tmp/c227.txt
echo "paths=$(wc -l < /tmp/c227.txt) go=$(grep -c '\.go$' /tmp/c227.txt || true)"
echo "outside=$(grep -vcE '^(internal/config/|\.scratch/wisp/probes/227/r1/|\.scratch/wisp/issues/227-)' /tmp/c227.txt || true)"
grep -vE '^(internal/config/|\.scratch/wisp/probes/227/r1/|\.scratch/wisp/issues/227-)' /tmp/c227.txt
echo "--- .go files touched ---"; grep '\.go$' /tmp/c227.txt
echo
echo "=== [B] 决定性突变我自己重打 ==="
cd "$W/tree" || exit 9
T=internal/config/writeguard_226_test.go
echo "schema_version 在该用例附近哪几行："
awk 'NR>=380 && NR<=412 && /schema_version/ {print NR": "$0}' "$T"
cp "$T" "$W/orig.txt"
# 内容锚：删掉该 fixture 里那行 schema_version（腿逐字＝删 `schema_version = 2`）
sed -i '/schema_version = 2/d' "$T"
echo "删掉行数=$(diff "$W/orig.txt" "$T" | grep -c '^<')  新sha256=$(sha256sum "$T" | cut -d' ' -f1)"
export GOTMPDIR="$W/gotmp" TMPDIR="$W/gotmp" TEMP="$W/gotmp" TMP="$W/gotmp"
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL)|ok  |FAIL|    writeguard)' | head -10
echo "MUT_RC=${PIPESTATUS[0]}"
cp "$W/orig.txt" "$T"
echo "还原 sha256=$(sha256sum "$T" | cut -d' ' -f1)  期望=49960da670dbe1350478a7f710900b52a5127dcfc24735250be1e4104cbef9b8"
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v 2>&1 | grep -E '^(--- (PASS|FAIL)|ok  |FAIL)' | head -4
echo "RESTORED_RC=${PIPESTATUS[0]}"
} 2>&1
