#!/bin/sh
# 编排者翻框机械件（2026-10-16 11:0x）：把三枚非实现者验收表的判语落成票面勾。
# 只改 "- [ ]" → "- [x]"，逐枚指名 AC 标题（⛔ 批量替换、⛔ 动正文）。
set -u
cd "$(git rev-parse --show-toplevel)" || exit 1

I=.scratch/wisp/issues
T317=$(ls $I/317-*.md)
T325=$(ls $I/325-*.md)
T165=$(ls $I/165-*.md)
T287=$(ls $I/287-*.md)

flip() {
  f=$1
  pat=$2
  before=$(grep -c "^- \[ \] \*\*${pat}" "$f" || true)
  sed -i "s/^- \[ \] \*\*${pat}/- [x] **${pat}/" "$f"
  after=$(grep -c "^- \[x\] \*\*${pat}" "$f" || true)
  printf 'flip %s :: %s  before-unchecked=%s after-checked=%s\n' "$(basename "$f")" "$pat" "$before" "$after"
}

flip "$T317" 'AC#3 整条'
flip "$T325" 'AC#1 改形'
flip "$T325" 'AC#2 成对反控'
flip "$T325" 'AC#3 顺带交出票 317'
flip "$T325" 'AC#4 门禁'
flip "$T325" 'AC#5 账'
flip "$T165" 'AC#1 未修码读数'
flip "$T165" 'AC#2 复用形状'
flip "$T165" 'AC#3 批准只批这一批'
flip "$T165" 'AC#4 面板不许代答'
flip "$T165" 'AC#5 契约轴'
flip "$T287" 'AC#2 造那枚夹具'
flip "$T287" 'AC#3 顺手答一句'
flip "$T287" 'AC#4 边界与门禁'

echo "=== 票框尺（翻完后）==="
for f in "$T317" "$T325" "$T165" "$T287"; do
  printf '%s unchecked=%s checked=%s\n' "$(basename "$f")" \
    "$(grep -c '^[[:space:]]*- \[ \]' "$f" || true)" "$(grep -c '^[[:space:]]*- \[x\]' "$f" || true)"
done
