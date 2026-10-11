#!/bin/sh
# A847 票面进度第二发：把逐枚 (checked, unchecked) 落成名册，再按形状分桶（⛔ 只算 u==0 会误报）
cd "$(git rev-parse --show-toplevel)" || exit 9
DIR=.scratch/wisp/issues
D=.scratch/wisp/probes/orch/20261011
RO=$D/logs/orch-A847-ticket-cells.txt
git ls-tree -r HEAD --name-only -- "$DIR" | grep -v README > "$RO.tmp.paths"
: > "$RO.tmp.cells"
while IFS= read -r p; do
  body=$(git show "HEAD:$p")
  c=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[x\]')
  u=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[ \]')
  case "$p" in *-done.md) st=done ;; *) st=open ;; esac
  printf '%s\t%s\t%s\t%s\n' "$st" "$c" "$u" "$(basename "$p")" >> "$RO.tmp.cells"
done < "$RO.tmp.paths"
sort -t"$(printf '\t')" -k1,1 -k3,3nr "$RO.tmp.cells" > "$RO"
{ echo "date=$(date '+%Y-%m-%d %H:%M:%S %z') HEAD=$(git rev-parse --short HEAD)  尺＝HEAD blob 层（⛔ 工作树：两枚验收腿在写票尾）"
  echo "roster_lines=$(grep -c '' "$RO")  (tab 分列：state  checked  unchecked  file)"
  echo "=== 分桶（开放票） ==="
  echo "open_total=$(awk -F'\t' '$1=="open"' "$RO" | grep -c '')"
  echo "open_零框（⛔ 判据格⛔ 算进度）=$(awk -F'\t' '$1=="open" && $2+$3==0' "$RO" | grep -c '')"
  echo "open_全勾待收口（u==0 且 c>0）=$(awk -F'\t' '$1=="open" && $3==0 && $2>0' "$RO" | grep -c '')"
  echo "open_一格未动（c==0 且 u>0）=$(awk -F'\t' '$1=="open" && $2==0 && $3>0' "$RO" | grep -c '')"
  echo "open_部分勾（c>0 且 u>0）=$(awk -F'\t' '$1=="open" && $2>0 && $3>0' "$RO" | grep -c '')"
  echo "=== 分桶（已结案票） ==="
  echo "done_total=$(awk -F'\t' '$1=="done"' "$RO" | grep -c '')"
  echo "done_零框=$(awk -F'\t' '$1=="done" && $2+$3==0' "$RO" | grep -c '')"
  echo "done_仍有未勾框（乙类）=$(awk -F'\t' '$1=="done" && $3>0' "$RO" | grep -c '')"
  echo "=== 格数合计（两把尺并排，⛔ 混用） ==="
  awk -F'\t' '{a+=$2;b+=$3} $1=="open"{ao+=$2;bo+=$3} END{printf "全库 checked=%d unchecked=%d total=%d => %.1f%%\n开放票 checked=%d unchecked=%d total=%d => %.1f%%\n",a,b,a+b,100*a/(a+b),ao,bo,ao+bo,100*ao/(ao+bo)}' "$RO"
  echo "=== 全勾待收口名册（若⛔ 0 枚就写空） ==="
  awk -F'\t' '$1=="open" && $3==0 && $2>0 {printf "  %s  checked=%s\n", $4, $2}' "$RO"
  echo "=== 开放票里未勾格最多的 12 枚（＝最重的在飞面） ==="
  awk -F'\t' '$1=="open" && $3>0 {printf "%4s  %s\n", $3, $4}' "$RO" | sort -rn | head -12
}
