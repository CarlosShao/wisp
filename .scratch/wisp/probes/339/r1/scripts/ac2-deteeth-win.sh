#!/bin/bash
# 339-r1 AC#2 去牙反控（Windows 台面＝牙齿新家，票 339 甲形把它从 Linux 搬来了）。
# 正控＝未突变 ⇒ 334 牙齿绿；去牙＝摘掉 aRelates 尾行回落支（同 334-v1 的 M-V1 形）⇒ 334 牙齿点名红。
# 突变只打在仓外导出树物理字节；逐发 git show 还原 + hash-object 自证 pristine。⛔ 信 echo。
# ⛔ 用行号当落点判据：内容锚＝"\treturn uncertainAnchorMiss(cand, anchor)"，现跑命中枚数必须是 1。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="/d/tmp/wisp/20261018/339r1/win/after"
BL="$TREE/internal/risk/blacklist.go"
ANCHOR=e6fd0f5051812a23c72db444d12d76b4d27285a7
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
RUN='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestAnchorFallbackTeethPerAnchor336'
mkdir -p "$LOGD"
# Windows-style temp roots: Git Bash would half-translate a /d/... TMPDIR into
# mixed separators and turn Go's own path prints into台面故障（今天踩过）。
export TEMP='D:\tmp\wisp\20261018\339r1\gt' TMP="$TEMP" TMPDIR="$TEMP" GOTMPDIR="$TEMP"
export GOFLAGS=-mod=mod GOPROXY=off CGO_ENABLED=0 GOCACHE='D:\tmp\wisp\20261018\339r1\gocache'
cd "$TREE" || exit 9

gh() { git -C "$REPO" hash-object "$1"; }
run_leg() { # $1=label
  local out="$LOGD/$1.log"
  go test ./internal/risk/ -run "$RUN" -count=1 -timeout 30m -v > "$out" 2>&1
  local rc=$?
  echo "[$1] RC=$rc RUN=$(grep -c '^=== RUN' "$out"||true) PASS=$(grep -c '^--- PASS' "$out"||true) FAIL=$(grep -c '^--- FAIL' "$out"||true) SKIP=$(grep -c '^--- SKIP' "$out"||true)"
  grep -E '^--- FAIL' "$out" || true
  grep -E '^    --- FAIL' "$out" || true
}

PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
echo "pristine_expect=$PRISTINE"
echo "P0 pristine_now=$(gh "$BL")"

echo "=== 正控（未突变树：牙齿必须在 Windows 绿，⛔ 恒红）==="
run_leg win-ac2-positive-unmutated

echo "=== 去牙一发：aRelates 回落支整支摘掉（M-V1 形）==="
hit=$(grep -cF "$(printf '\treturn uncertainAnchorMiss(cand, anchor)')" "$BL")
echo "anchor_hits=$hit (must be 1)"
[ "$hit" -eq 1 ] || { echo "GUARD: content anchor not unique, aborting"; exit 2; }
printf 'MUT_BEGIN %s\n' "$(gh "$BL")"
perl -0pi -e "s/\Q$(printf '\treturn uncertainAnchorMiss(cand, anchor)')\E/\treturn false \/\/ M-V1 ticket339 deteeth/" "$BL"
echo "mutated_hash=$(gh "$BL")"
run_leg win-ac2-deteeth-fallback

echo "=== 当场还原（git show 物理字节，⛔ 从记忆重写）==="
git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$BL"
R1=$(gh "$BL")
echo "restored_hash=$R1 match=$( [ "$R1" = "$PRISTINE" ] && echo YES || echo NO )"

echo "=== 成对正控二跑（还原后必须回绿）==="
run_leg win-ac2-positive-after-restore
R2=$(gh "$BL")
echo "final_hash=$R2 match=$( [ "$R2" = "$PRISTINE" ] && echo YES || echo NO )"
