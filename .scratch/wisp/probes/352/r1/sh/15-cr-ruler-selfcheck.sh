#!/bin/bash
# 尺⑮：CR 尺自校——同一枚字节流上三把尺对撞（证明 `grep -c $'\r'` 这台 shell 上是钝尺）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
probe() {
  f="$1"
  a=$(tr -cd '\r' < "$f" | wc -c)
  b=$(awk 'BEGIN{n=0} /\r/{n++} END{print n}' "$f")
  c=$(grep -c $'\r' "$f" || true)
  l=$(wc -l < "$f")
  printf '%s CRBYTES_tr=%s CRROWS_awk=%s CRROWS_grep=%s LINES=%s\n' "$(basename "$f")" "$a" "$b" "$c" "$l"
}
echo '== 已知纯 LF 件（正控：CR 必 0）=='
probe internal/winsec/winsec.go
probe .scratch/wisp/probes/352/r1/readings-ac3-discipline.md
echo '== 造一枚含 CR 的件（自门：三把尺都该报>0，否则尺钝）=='
D=/d/tmp/wisp/20261020/352r1cont/tmp
mkdir -p "$D"
printf 'line one\r\nline two\r\nline three\n' > "$D/crlf-probe.txt"
probe "$D/crlf-probe.txt"
echo '== 我本趟要入库的两枚新件（须 CRBYTES=0）=='
probe .scratch/wisp/probes/352/r1/sh/13-cont-roster.sh
probe /d/tmp/wisp/20261020/352r1cont/sh/14-claims-extract.sh
echo '== 孤儿件形制（上一腿遗产）=='
probe .scratch/wisp/probes/352/r1/logs/handoff-number-reconcile.txt
