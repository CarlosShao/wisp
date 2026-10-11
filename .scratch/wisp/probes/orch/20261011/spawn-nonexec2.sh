#!/usr/bin/env bash
# Q-304-a1-I 重量：非 os/exec 的直调起子进程形状有没有（编排者自跑，只读，⛔ 跑 Go）
# ⚠ 前一发 orch-A852-spawn-nonexec.txt 里 `exec\.Command(` 那行 hits=0 是我的尺坏了：
#   ERE 的裸 `(` 开了一个没闭合的组 ⇒ grep 报错（被 2>/dev/null 吞掉）⇒ 管道只剩 wc 的 0 行＝假 0。这发把括号转义并留下 stderr。
set -u
cd "/d/work/workspace/projects plans/Wisp"
O=.scratch/wisp/probes/orch/20261011/logs/orch-A852-spawn-nonexec2.txt
: > "$O"
{
  echo "ts=$(date '+%F %T') HEAD=$(git rev-parse --short HEAD)  尺=grep -rn --include='*.go' -E（只读，⛔ 跑 Go）  射程=internal cmd tools scripts"
  echo "分母＝同一射程里的 .go 枚数："
  find internal cmd tools scripts -name '*.go' | wc -l | sed 's/^/  go_files=/'
} >> "$O"
for pat in 'os\.StartProcess' 'syscall\.StartProcess' 'CreateProcess' 'ShellExecute' 'syscall\.Exec' 'exec\.Command' 'exec\.CommandContext' '"build"' '"-o"'; do
  err=.scratch/wisp/probes/orch/20261011/logs/orch-A852-greperr.txt
  c=$(grep -rn --include='*.go' -E "$pat" internal cmd tools scripts 2>"$err" | wc -l)
  e=$(head -1 "$err" | cut -c1-80)
  echo "pattern=$pat hits=$c grep_stderr=\"${e:-none}\"" >> "$O"
done
echo "--- exec\.Command 逐枚落点前 10 行（分母＝上面 hits）:" >> "$O"
grep -rn --include='*.go' -E 'exec\.Command' internal cmd tools scripts 2>/dev/null | head -10 | cut -c1-150 >> "$O"
echo "--- 非 os/exec 的直调形（若 hits>0 才逐枚列）:" >> "$O"
grep -rn --include='*.go' -E 'os\.StartProcess|syscall\.StartProcess|CreateProcess|ShellExecute' internal cmd tools scripts 2>/dev/null | cut -c1-160 >> "$O"
echo "done_rc=$?" >> "$O"
wc -l < "$O" | sed 's/^/lines=/'
