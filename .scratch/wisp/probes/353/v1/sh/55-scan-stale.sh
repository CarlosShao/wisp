#!/bin/sh
# 353-v1 R55 = 主件里残留旧口径的扫描（量的）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for f in docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1/00-ruler-register.md .scratch/wisp/probes/353/v1/01-independent-readings.txt .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md; do
  echo "== WT $f"
  grep -n '6 发里 3 发\|3 发的原始\|3\.\.54\|10\.\.31\|115\.\.153' "$f" | cut -c1-170 || echo "  (0 命中)"
done
echo "== 正控：同一把尺打已入库的草案笔⑤ blob（应命中 0＝那句已随笔⑤入库的写法⛔ 含这些串）=="
git cat-file blob 'HEAD:.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md' | grep -c '6 发里 3 发' || true
