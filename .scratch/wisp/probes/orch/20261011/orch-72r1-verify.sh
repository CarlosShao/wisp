#!/bin/sh
# 编排者自己复跑 72-r1 的交件：两台面基线＋四发反形＋还原自证＋门禁
# 台件＝仓外导出树（git archive HEAD go.mod go.sum internal cmd）；⛔ 在共享工作树里跑改前那一发
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="/c/Users/swq/tmp/wisp-72r1-orch-tree"
ALIAS="/c/Users/swq/tmp/riskattrib1/RunnerTempAreaForAttribution1"
ALIAS_SHORT='C:\Users\swq\tmp\RISKAT~1\RUNNER~1'
D="$REPO/.scratch/wisp/probes/orch/20261011"
LOG="$D/logs/orch-72r1-verify.txt"
PRE=9c83598d   # 72-r1 起手锚的父（改动前那一版）
T3_FILES="internal/risk/pathresolver_anchor_spelling_windows_test.go internal/risk/pathresolver_junction_windows_test.go internal/risk/syncdirs_test.go"

{
echo "# orch-72r1-verify start $(date '+%F %T') HEAD=$(cd "$REPO" && git rev-parse --short HEAD)"

fresh_tree() {
  rm -rf "$TREE" 2>/dev/null
  mkdir -p "$TREE" || { echo "tree_mkdir_rc=$?"; exit 1; }
  (cd "$REPO" && git archive HEAD go.mod go.sum internal cmd | tar -x -C "$TREE")
  echo "archive_rc=$? tree_go=$(find "$TREE" -name '*.go' | wc -l)"
}

run_aliased() { # $1=phase label
  cd "$TREE" || exit 1
  mkdir -p "$ALIAS" 2>/dev/null
  s=$(date +%s)
  TEMP="$ALIAS_SHORT" TMP="$ALIAS_SHORT" TMPDIR="$ALIAS_SHORT" \
    go test ./internal/risk/ -count=1 -timeout 900s -v > "$LOG.$1" 2>&1
  rc=$?
  echo "--- phase=$1 aliased rc=$rc elapsed=$(( $(date +%s) - s ))s"
  echo "    RUN=$(grep -c '^=== RUN   Test' "$LOG.$1") PASS=$(grep -c '^--- PASS\|^    --- PASS' "$LOG.$1") FAIL=$(grep -c '^--- FAIL\|^    --- FAIL' "$LOG.$1") SKIP=$(grep -c '^--- SKIP\|^    --- SKIP' "$LOG.$1")"
  echo "    fail names:"
  grep -E '^\s*--- FAIL' "$LOG.$1" | sed 's/^/      /'
}

run_default() { # $1=label
  cd "$TREE" || exit 1
  s=$(date +%s)
  go test ./internal/risk/ -count=1 -timeout 900s -v > "$LOG.$1" 2>&1
  rc=$?
  echo "--- phase=$1 default-TEMP rc=$rc elapsed=$(( $(date +%s) - s ))s"
  echo "    RUN=$(grep -c '^=== RUN   Test' "$LOG.$1") PASS=$(grep -c '^--- PASS\|^    --- PASS' "$LOG.$1") FAIL=$(grep -c '^--- FAIL\|^    --- FAIL' "$LOG.$1") SKIP=$(grep -c '^--- SKIP\|^    --- SKIP' "$LOG.$1")"
  echo "    fail names:"
  grep -E '^\s*--- FAIL' "$LOG.$1" | sed 's/^/      /'
}

# ============ 1. HEAD（含 72-r1）两台面基线 ============
fresh_tree
echo "=== [1] HEAD 两台面（判据＝别名台面 FAIL 应为 0；本机台面颜色应与 A839 那次同形） ==="
run_aliased P1-aliased-head
run_default  P1-default-head

# ============ 2. 反形 M-a：三枚测试件退回 PRE（= 改动前那一版） ============
echo "=== [2] M-a 退回改前三枚测试件（判据＝别名台面应恰 7 枚红，名册＝票 72 甲族） ==="
for f in $T3_FILES; do (cd "$REPO" && git show "$PRE:$f") > "$TREE/$f"; echo "revert_rc=$? $f"; done
run_aliased P2-aliased-pristtest
for f in $T3_FILES; do (cd "$REPO" && git show "HEAD:$f") > "$TREE/$f"; done

# ============ 3. 反形 M-b：namesSameFile 退回旧的字符串同形 ============
echo "=== [3] M-b namesSameFile 换成 strings.EqualFold(a,b)（判据＝甲族 #1/#4/#5/#6 应转红＝那枚夹具确实是这四枚的承重件） ==="
sed -i 's|^\treturn fa.real != "" && fb.real != "" && strings.EqualFold(fa.real, fb.real)$|\t_ = fa; _ = fb; return strings.EqualFold(a, b)|' "$TREE/internal/risk/pathresolver_anchor_spelling_windows_test.go"
echo "mb_patch_lines=$(grep -c 'return strings.EqualFold(a, b)' "$TREE/internal/risk/pathresolver_anchor_spelling_windows_test.go")"
run_aliased P3-aliased-equalfold
sed -i 's|^\t_ = fa; _ = fb; return strings.EqualFold(a, b)$|\treturn fa.real != "" \&\& fb.real != "" \&\& strings.EqualFold(fa.real, fb.real)|' "$TREE/internal/risk/pathresolver_anchor_spelling_windows_test.go"
echo "mb_restore_lines=$(grep -c 'return fa.real != "" && fb.real != "" && strings.EqualFold(fa.real, fb.real)' "$TREE/internal/risk/pathresolver_anchor_spelling_windows_test.go")"

# ============ 4. 反形 M-c：spellings() 恒交两形（测「恰 1 形」那半句还剩不剩仪器） ============
echo "=== [4] M-c 单形那一支也恒交两形（判据＝若整包零红，则 AC#5 那句「只产一种比较形」确实零仪器；⚠ 若转红要看红句是不是大小写敏感那另一件事） ==="
sed -i 's|^\t\treturn \[\]string{f.raw}$|\t\treturn []string{f.raw, strings.ToUpper(f.raw)}|' "$TREE/internal/risk/pathresolver.go"
echo "mc_patch_lines=$(grep -c 'return \[\]string{f.raw, strings.ToUpper(f.raw)}' "$TREE/internal/risk/pathresolver.go")"
run_aliased P4-aliased-twofields
sed -i 's|^\t\treturn \[\]string{f.raw, strings.ToUpper(f.raw)}$|\t\treturn []string{f.raw}|' "$TREE/internal/risk/pathresolver.go"
echo "mc_restore_lines=$(grep -c 'return \[\]string{f.raw}$' "$TREE/internal/risk/pathresolver.go")"

# ============ 5. 反形 M-d：sync 那枚兜底守卫退恒真 ============
echo "=== [5] M-d syncdirs.go 兜底支退恒真（判据＝新那两枚负向钉必须红；旧正向全绿） ==="
sed -i 's|^\tif !s.complete && s.home != "" && aUnder(candForms, anchorForms(s.home)) {$|\tif true {|' "$TREE/internal/risk/syncdirs.go"
echo "md_patch_lines=$(grep -c '^	if true {$' "$TREE/internal/risk/syncdirs.go")"
run_aliased P5-aliased-tautology
echo "    新负向钉那两句在 M-d 里的红句（带文件:行）："
grep -E 'not-a-sync-tree|outside the profile and outside' "$LOG.P5-aliased-tautology" | sed 's/^/      /' | head -6

# ============ 6. 反形 M-e：恒真形 + 拔掉两枚新负向钉（证明这两枚是我加的那点仪器在响） ============
echo "=== [6] M-e 恒真形 + syncdirs_test.go 退回 PRE（判据＝红数应比 M-d 少，少的正是新钉那两枚） ==="
(cd "$REPO" && git show "$PRE:internal/risk/syncdirs_test.go") > "$TREE/internal/risk/syncdirs_test.go"
echo "me_revert_rc=$?"
run_aliased P6-aliased-taut-noNails
echo "    M-d ⇄ M-e 红名作差："
grep -E '^\s*--- FAIL' "$LOG.P5-aliased-tautology" | sed 's/^ *--- FAIL: //;s/ .*//' | sort > "$LOG.m-d.names"
grep -E '^\s*--- FAIL' "$LOG.P6-aliased-taut-noNails" | sed 's/^ *--- FAIL: //;s/ .*//' | sort > "$LOG.m-e.names"
echo "      only_in_Md=$(comm -23 "$LOG.m-d.names" "$LOG.m-e.names" | grep -c '')"
comm -23 "$LOG.m-d.names" "$LOG.m-e.names" | sed 's/^/        /'
echo "      only_in_Me=$(comm -13 "$LOG.m-d.names" "$LOG.m-e.names" | grep -c '')"
comm -13 "$LOG.m-d.names" "$LOG.m-e.names" | sed 's/^/        /'

# ============ 7. 还原自证：整树重导出 + sha256 对拉 ============
echo "=== [7] 还原自证（逐枚 sha256 前 16 位 ⇄ 仓内 HEAD blob） ==="
fresh_tree
for f in $T3_FILES internal/risk/pathresolver.go internal/risk/syncdirs.go internal/risk/blacklist.go; do
  t=$(sha256sum "$TREE/$f" | cut -c1-16)
  r=$(cd "$REPO" && git show "HEAD:$f" | sha256sum | cut -c1-16)
  [ "$t" = "$r" ] && echo "      SAME $t $f" || echo "      DIFF tree=$t repo=$r $f"
done
echo "=== 正控：还原后的树再跑一发别名（应仍 FAIL=0） ==="
run_aliased P7-aliased-restored

# ============ 8. 门禁（在仓内跑，⛔ 树里） ============
echo "=== [8] 门禁 ==="
cd "$REPO" || exit 1
GOFLAGS= go build ./... > "$D/logs/tmp-build-72r1.txt" 2>&1; echo "build_rc=$? (件 $(wc -c < "$D/logs/tmp-build-72r1.txt") 字节)"
go vet ./internal/risk/ > "$D/logs/tmp-vet-72r1.txt" 2>&1; echo "vet_rc=$? (件 $(wc -c < "$D/logs/tmp-vet-72r1.txt") 字节)"
sh scripts/d22scan.sh > "$D/logs/tmp-d22-72r1.txt" 2>&1; echo "d22scan_rc=$?"
echo "gofmt_wt_internal_risk=$(gofmt -l internal/risk | grep -c '')"
gofmt -l internal/risk | sed 's/^/      /'
for f in $T3_FILES; do (cd "$REPO" && git show "HEAD:$f") > /tmp/blobcheck.go; echo "      blob_gofmt $f = $(gofmt -l /tmp/blobcheck.go | grep -c '')"; done
echo "eol_三枚写面:"; git ls-files --eol $T3_FILES | sed 's/^/      /'
echo "worktree_dirty_internal_cmd=$(git status --porcelain -- internal cmd | grep -c '')"
echo "ticket72_boxes=$(grep -cE '^[[:space:]]*- \[[ x]\]' .scratch/wisp/issues/72-atble-classification-runner.md) unchecked=$(grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/72-atble-classification-runner.md)"
echo "# done $(date '+%F %T')"
} >> "$LOG" 2>&1
