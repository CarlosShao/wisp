#!/bin/bash
# 339-v1 台件⑨：AC#5 禁区尺＋AC#6 纪律尺（裁决腿自己重打，⛔ 抄腿／编排者的数）。
# 全部＝量的、具名提交树 e6fd0f50..7437ef88（落码区间）与 e6fd0f50..HEAD（交件全程）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
LOG=".scratch/wisp/probes/339/v1/logs/v1-forbidden-and-discipline.txt"
PRE=e6fd0f5051812a23c72db444d12d76b4d27285a7
CODE=7437ef882fac5b427bf17c4bf590819c58d0b375
{
  echo "== 尺①禁区件在 e6fd0f50..HEAD 的被动次数（应全 0）=="
  for p in internal/risk/blacklist.go internal/risk/pathresolver.go internal/risk/pathresolver_other.go \
           internal/risk/pathresolver_windows.go internal/risk/anchor_fallback_teeth_336_windows_test.go \
           internal/risk/anchor_fallback_teeth_336_other_test.go internal/risk/name_claim_two_arm_216_test.go \
           internal/observe/thresholds.go scripts/portable-tests.sh .github/workflows/ci.yml go.mod go.sum .gitignore \
           internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go \
           docs/PLAN.md frontend design cmd/wisp internal/panel; do
    N=$(git log --oneline $PRE..HEAD -- "$p" | wc -l)
    echo "  touched=$N  $p"
  done
  echo
  echo "== 尺②三枚生产码件 blob 同串（⛔ 非注释字节＝整件未动）=="
  for p in internal/risk/blacklist.go internal/risk/pathresolver.go internal/risk/pathresolver_other.go; do
    A=$(git rev-parse "$PRE:$p"); B=$(git rev-parse "HEAD:$p")
    echo "  $p pre=$A head=$B same=$([ "$A" = "$B" ] && echo YES || echo NO)"
  done
  echo
  echo "== 尺③票 339 票面＝纯追加？numstat ＋ Status 行 ＋ 票名 rename =="
  echo "  numstat $PRE..HEAD:"; git diff --numstat $PRE HEAD -- .scratch/wisp/issues/ | sed 's/^/    /'
  echo "  票 339 文件 rename 检测:"; git log --follow --name-status --format='  commit=%h' $PRE..HEAD -- '.scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md' | sed 's/^/    /'
  echo "  Status 行现读（第 3 行前 90 字）:"; sed -n '3p' .scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md | cut -c1-90 | sed 's/^/    /'
  echo
  echo "== 尺④票 334 票面文件在 339-r1 期间被动枚数（应 0＝它没改别票）=="
  git log --oneline $PRE..HEAD --name-only -- '.scratch/wisp/issues/334-*' | sed 's/^/    /'
  echo "  (空＝零被动)"
  echo
  echo "== 尺⑤甲形落码那笔的 .go 面（恰 2 枚、⛔ 别家件）=="
  git show --name-only --format='  %h' $CODE | sed 's/^/    /'
  echo
  echo "== 尺⑥断言未放宽：want A／ClassA 那三句在 HEAD 逐字在场 =="
  grep -n 'want A\|ClassA\|non-overridable' internal/risk/anchor_miss_reason_334_test.go | sed 's/^/    /'
  echo "  t.Skip 命中（两枚仪器件，应 0）:"
  grep -c 't\.Skip' internal/risk/anchor_miss_reason_334_test.go internal/risk/anchor_miss_reason_339_other_test.go | sed 's/^/    /'
  echo
  echo "== 尺⑦新导出名尺（⛔ 新导出名＝两枚件里以大写字母开头的顶层 func/const/type 新增）=="
  git diff $PRE HEAD -- internal/risk/ | grep -E '^\+(func|const|type|var) [A-Z]' | sed 's/^/    /'
  echo "  (只许出现测试函数名 Test*，测试名不算生产导出 API＝同包内 336/334 先例同形)"
  echo
  echo "== 尺⑧绊线件是否真引用了它点名的那三枚锚（⛔ 注释里指名不存在的件）=="
  for t in pathresolver_other.go blacklist.go name_claim_two_arm_216_test.go anchor_fallback_teeth_336_other_test.go; do
    echo "    ref $t in 339_other=$(grep -c "$t" internal/risk/anchor_miss_reason_339_other_test.go || true)"
  done
} > "$LOG" 2>&1
cat "$LOG"
