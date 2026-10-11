#!/sh
# A846 门禁与暂存 第二发（改对我 a 发那两把被 core.quotePath 骗了的尺：件名含非 ASCII 时 git 会加引号＋八进制转义
# ⇒ ① grep -v '^\.scratch/' 前缀⛔ 匹配 ② [ -s "带引号的路径" ] 必然为假。正解＝git -c core.quotePath=false ＋ 逐枚 -z 读）
cd "$(git rev-parse --show-toplevel)" || exit 9
D=.scratch/wisp/probes/orch/20261011
LOG=$D/logs/orch-A846-gates-b.txt
T298=.scratch/wisp/issues/298-two-gofmt-dirty-test-files-at-head-and-one-crlf-fake-cell.md
T308=.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md
T294=.scratch/wisp/issues/294-no-production-shape-test-for-the-mute-attach-hop.md
T295=.scratch/wisp/issues/295-ra-started-bare-bool-written-off-thread.md
T296=.scratch/wisp/issues/296-hotkey-reload-source-skips-applyhotkeydefaults-keys-die-1s-after-boot.md
T309=$(ls .scratch/wisp/issues/309-*.md)

git add -- $D $T298 $T308 $T294 $T295 $T296 2>/dev/null
Q="git -c core.quotePath=false diff --cached --name-only"
$Q > $D/logs/staged-A846-b.txt
n_all=$(grep -c '' $D/logs/staged-A846-b.txt)
n_go=$(grep -c '[.]go$' $D/logs/staged-A846-b.txt)
n_out=$(grep -c '[.]out$' $D/logs/staged-A846-b.txt)
foreign=$(grep -v '^docs/reports/\|^\.scratch/wisp/issues/\|^\.scratch/wisp/probes/orch/20261011/' $D/logs/staged-A846-b.txt | grep -c '')
inflight=$(grep -cE '^internal/|^cmd/|^\.scratch/wisp/probes/(309|308)/r1/' $D/logs/staged-A846-b.txt)
zero=$(git diff --cached --name-only -z | tr '\0' '\n' | while IFS= read -r f; do [ -s "$f" ] || echo "$f"; done | grep -c '')
longname=$(grep -v '^docs/reports/\|^\.scratch/' $D/logs/staged-A846-b.txt | grep -c '')

{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)"
  echo "staged_total=$n_all staged_go=$n_go staged_dot_out=$n_out foreign_staged=$foreign inflight_staged=$inflight zero_byte_staged=$zero"
  echo "对照上一发（同一名册、⛔ 关 quotePath）＝foreign_staged=1 zero_byte_staged=1 ⇒ 两枚都是尺坏⛔ 码坏，本发改对"
  echo "numstat_ledger=$(git diff --cached --numstat -- docs/reports/pending-and-issues.md | tr '\n' ' ')"
  echo "numstat_handover=$(git diff --cached --numstat -- docs/reports/HANDOVER.md | tr '\n' ' ')"
  for t in "$T298" "$T308" "$T294" "$T295" "$T296"; do
    echo "numstat_$(basename "$t" | cut -c1-3)=$(git diff --cached --numstat -- "$t" | tr '\n' ' ')"
  done
  echo "prefix_ledger=$(git show HEAD:docs/reports/pending-and-issues.md | cmp -s - <(head -"$(git show HEAD:docs/reports/pending-and-issues.md | wc -l)" docs/reports/pending-and-issues.md) && echo IDENTICAL || echo DIFF)"
  echo "prefix_handover=$(git show HEAD:docs/reports/HANDOVER.md | cmp -s - <(head -"$(git show HEAD:docs/reports/HANDOVER.md | wc -l)" docs/reports/HANDOVER.md) && echo IDENTICAL || echo DIFF)"
  echo "box_298=$(grep -cE '^[[:space:]]*- \[x\]' $T298)/$(grep -cE '^[[:space:]]*- \[ \]' $T298) box_308=$(grep -cE '^[[:space:]]*- \[x\]' $T308)/$(grep -cE '^[[:space:]]*- \[ \]' $T308) box_309=$(grep -cE '^[[:space:]]*- \[x\]' $T309)/$(grep -cE '^[[:space:]]*- \[ \]' $T309) box_294=$(grep -cE '^[[:space:]]*- \[x\]' $T294)/$(grep -cE '^[[:space:]]*- \[ \]' $T294) box_295=$(grep -cE '^[[:space:]]*- \[x\]' $T295)/$(grep -cE '^[[:space:]]*- \[ \]' $T295) box_296=$(grep -cE '^[[:space:]]*- \[x\]' $T296)/$(grep -cE '^[[:space:]]*- \[ \]' $T296)"
  echo "bullet_308=$(grep -c '^- \[2026' $T308)  # 我那一枚＋代提的 308-r1 那一枚"
  echo "worktree_dirty_internal_cmd=$(git status --porcelain -- internal cmd | wc -l | tr -d ' ')"
  echo "ticket309_staged_by_me=$(grep -c "$(basename "$T309")" $D/logs/staged-A846-b.txt)  # 应为 0＝那枚腿的票面增量⛔ 由我提交"
} > $LOG 2>&1
cat $LOG
