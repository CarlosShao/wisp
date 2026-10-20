#!/bin/sh
# 353-v1 R43 = 落笔前预检（工作树形；分节码位/CR/行数/未完成标记词/末字节）
S7=$(printf '\xc2\xa7')
for f in /d/tmp/wisp/20261020/353v1/msg-02.txt /d/tmp/wisp/20261020/353v1/msg-03.txt /d/tmp/wisp/20261020/353v1/msg-04.txt /d/tmp/wisp/20261020/353v1/msg-05.txt \
         '/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/353/v1/00-ruler-register.md' \
         '/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/353/v1/01-independent-readings.txt' \
         '/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/353/v1/02-ticket-draft-secC.md' \
         '/d/work/workspace/projects plans/Wisp/docs/evidence/s1/353-v1-verdict.md' ; do
  [ -f "$f" ] || { printf '%-70s ABSENT\n' "$f"; continue; }
  printf '%-58s SECT=%-3s CR=%-3s LINES=%-4s PH=%-3s LAST=%s\n' "$(basename "$f")" \
    "$(grep -o "$S7" "$f" | wc -l)" \
    "$(tr -cd '\r' < "$f" | wc -c)" \
    "$(wc -l < "$f")" \
    "$(grep -cE 'TODO|FIXME|XXX' "$f" || true)" \
    "$(tail -c 1 "$f" | od -An -tx1 | tr -d ' \n')"
done
echo "== 正控：自造 2 枚分节码位 + 2 枚 CR 用同尺读回 =="
printf "a${S7}b${S7}c\r" > /d/tmp/wisp/20261020/353v1/_ctrl.bin
printf '  CTRL SECT=%s CR=%s\n' "$(grep -o "$S7" /d/tmp/wisp/20261020/353v1/_ctrl.bin | wc -l)" "$(tr -cd '\r' < /d/tmp/wisp/20261020/353v1/_ctrl.bin | wc -c)"
