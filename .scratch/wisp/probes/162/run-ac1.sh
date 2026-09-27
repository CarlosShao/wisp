#!/usr/bin/env bash
# 162-r1 台件：票 162 AC#1 的读数跑法（"整文件写回时漏抄一段"在未修码上到底有没有信号）。
#
# 这条命令就是读数本身：Go 测试里 t.Logf 打出来的那五行（(a)/(b)/(c) + 结论）
# 是本票前提成立与否的唯一凭据，不许手抄、不许事后编辑。
#
# 跑法：bash .scratch/wisp/probes/162/run-ac1.sh > .scratch/wisp/probes/162/ac1-reading.txt
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
go test -count=1 -run 'TestFSWriteSilentLossIsNotReported$' -v ./internal/tools/ 2>&1 \
  | grep -v 'level=INFO'
