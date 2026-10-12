#!/usr/bin/env bash
# 310-v1 起手名册 ① 四条精确名 tasklist 各一发（正典 gate-roster-304.md §1 的 bash 拼法）
# 每条 0 读数同行带正控（svchost.exe 必非零）+ stderr 字节数。空输出不当 0。
set -u
OUT="${1:?usage: $0 <out-file>}"
: > "$OUT"
for n in wisp.exe balldebug.exe mockllm.exe wisp.test.exe svchost.exe; do
  f=".scratch/wisp/probes/310/v1/logs/00-tl-$n.txt"
  e=".scratch/wisp/probes/310/v1/logs/00-tl-$n.err"
  tasklist //FI "IMAGENAME eq $n" > "$f" 2> "$e"
  rc=$?
  hits=$(grep -c -- "$n" "$f" || true)
  bytes=$(wc -c < "$f" || true)
  errb=$(wc -c < "$e" || true)
  pc=$(tasklist //FI "IMAGENAME eq svchost.exe" 2>/dev/null | grep -c -- "svchost.exe" || true)
  echo "exact_name=$n rc=$rc hit_lines=$hits stdout_bytes=$bytes stderr_bytes=$errb positive_control_svchost=$pc" >> "$OUT"
  if [ "$hits" = "0" ]; then echo "  name=$n empty above = 0 枚（hit_lines 由 grep -c 自证，stderr_bytes 同行，空输出不当 0）" >> "$OUT"; fi
done
echo "tl_done=1" >> "$OUT"
