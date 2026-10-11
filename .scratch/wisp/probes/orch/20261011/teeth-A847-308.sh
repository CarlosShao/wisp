#!/bin/sh
# A847 承重复跑（编排者自己跑票 308 AC#2 那枚"拔掉守卫⇒回到 panic"凭据）
# 规矩：先跑正控（原样必须绿）→ 突变（必须红且逐字 panic）→ 立刻还原并 cmp＋sha256 证明逐字节一致
# 台面：此刻两枚腿都交完、`git status -- internal cmd`＝0 行（无人在飞），突变只打 internal/risk 的测试件
cd "$(git rev-parse --show-toplevel)" || exit 9
F=internal/risk/pathresolver_volumeless_spelling_308_test.go
D=.scratch/wisp/probes/orch/20261011
L=$D/logs/orch-A847-308-teeth.txt
BK=$D/logs/A847-308-guard-backup.go.txt          # 从一开始就落自家 logs/，⛔ 仓库根、⛔ 叫 .go
SHA_BEFORE=$(sha256sum < "$F" | cut -d' ' -f1)
cp "$F" "$BK"

restore() {
  cp "$BK" "$F"
  SHA_AFTER=$(sha256sum < "$F" | cut -d' ' -f1)
  cmp -s "$BK" "$F" && echo "RESTORE cmp: identical" || echo "RESTORE cmp: DIFFER"
  echo "RESTORE sha_before=$SHA_BEFORE sha_after=$SHA_AFTER"
  echo "RESTORE worktree_lines=$(git status --porcelain -- internal cmd | grep -c '')"
}
trap restore EXIT

R() { # $1=tag
  GOFLAGS= go test ./internal/risk/ -count=1 -v -run TestVolumeInvariantSpellingsVolumeLess308 > "$L.$1.raw" 2>&1
  RC=$?
  { echo "--- $1 rc=$RC"
    echo "    panic_line=$(grep -m1 'panic: runtime error: slice bounds out of range' "$L.$1.raw" || echo ABSENT)"
    echo "    topPASS=$(grep -c '^--- PASS' "$L.$1.raw") topFAIL=$(grep -c '^--- FAIL' "$L.$1.raw") skip=$(grep -c -- '--- SKIP' "$L.$1.raw")"
    echo "    pkg_eval=$(grep -c -- '--- PASS' "$L.$1.raw") 终态=$(grep -m1 -E '^ok |^FAIL|^panic' "$L.$1.raw")"
  } >> "$L"
  return $RC
}

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)"
  echo "sha_before=$SHA_BEFORE (与 308-r1 记的 f2c01307… 对拉＝${SHA_BEFORE:0:8})"
} >> "$L"

R pristine
# 突变＝拔掉 len(vol)==0 那三行守卫，让 vol[:1] 裸跑（⛔ 按行号删＝行号会漂；按内容锚删）
sed '/if len(vol) == 0 {/,+2d' "$BK" > "$F"
echo "landed_check_guard_lines_left=$(grep -c 'if len(vol) == 0' "$F") delta_vs_backup=$(diff "$BK" "$F" | grep -c '^[<>]')" >> "$L"
R pulled
restore
trap - EXIT
R restored
cp "$L.pristine.raw" "$L.raw" 2>/dev/null
cat "$L"
