#!/bin/sh
# 353-v1 R60 = 哪几枚已入库尺本体里含那三枚未完成标记词的字面 pattern（＝尺自身、⛔ 正文占位符；具名登记用）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
PHRE=$(printf 'T%soDO|FIX%sME|XXX' 'O' 'X')
echo "== 已入库 .sh（HEAD）里含该 pattern 的枚数与件名 =="
n=0
for p in $(git ls-tree --name-only -r HEAD .scratch/wisp/probes/353/v1/); do
  case "$p" in *.sh) ;; *) continue ;; esac
  if git cat-file blob "HEAD:$p" | grep -qE "$PHRE"; then echo "  RULERHASPTN $p"; n=$((n+1)); fi
done
echo "  枚数=$n"
echo "== 四枚交付件（.md/.txt）里该 pattern 命中（应全 0）=="
for p in docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1/00-ruler-register.md .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md .scratch/wisp/probes/353/v1/90-truncated-nonzero-rows.txt; do
  if git cat-file -e "HEAD:$p" 2>/dev/null; then
    printf '  %-52s PH=%s\n' "$(basename "$p")" "$(git cat-file blob "HEAD:$p" | grep -cE "$PHRE" || true)"
  else
    printf '  %-52s WT_PH=%s\n' "$(basename "$p")" "$(grep -cE "$PHRE" "$p" || true)"
  fi
done
echo "== 正控（同尺打一把已知含该 pattern 的 .sh＝应 1）=="
git cat-file blob "HEAD:.scratch/wisp/probes/353/v1/sh/30-body-rulers.sh" | grep -cE "$PHRE" || true
