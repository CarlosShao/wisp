#!/bin/sh
# 342-v2 第 4 笔的提交尺件（★逐枚字面 pathspec、⛔ 从共享 index 反推、⛔ add -A、⛔ commit -- .、⛔ amend）
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
git add \
  .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md \
  .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md \
  .scratch/wisp/probes/342/v2/logs/bullet-342.txt \
  .scratch/wisp/probes/342/v2/logs/bullet-348.txt \
  .scratch/wisp/probes/342/v2/logs/logs-01-roster.txt \
  .scratch/wisp/probes/342/v2/logs/logs-02-legdiff.txt \
  .scratch/wisp/probes/342/v2/logs/logs-04-buildtrees.txt \
  .scratch/wisp/probes/342/v2/logs/logs-06-posix-nails.txt \
  .scratch/wisp/probes/342/v2/logs/logs-08-windows-nails.txt \
  .scratch/wisp/probes/342/v2/logs/logs-09-posix-fullpkg.txt \
  .scratch/wisp/probes/342/v2/logs/logs-12-windows-fullpkg.txt \
  .scratch/wisp/probes/342/v2/logs/x01-discipline.txt \
  .scratch/wisp/probes/342/v2/logs/x02-ticket-tails.txt \
  .scratch/wisp/probes/342/v2/logs/x03-banned-split.txt \
  .scratch/wisp/probes/342/v2/sh/20-commit-batch.sh \
  .scratch/wisp/probes/342/v2/sh/21-ticket-tails-append.sh \
  .scratch/wisp/probes/342/v2/sh/22-banned-codepoint-split.sh \
  .scratch/wisp/probes/348/v1/sh/21-ticket-tails-append.sh \
  .scratch/wisp/probes/348/v1/sh/22-banned-codepoint-split.sh
git commit -F 'D:\tmp\wisp\20261019\342v2\cm-04.txt' -- \
  .scratch/wisp/issues/342-privatefile-open-follows-a-swapped-ancestor-on-both-platforms-no-ticket-owns-it.md \
  .scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md \
  .scratch/wisp/probes/342/v2/logs/bullet-342.txt \
  .scratch/wisp/probes/342/v2/logs/bullet-348.txt \
  .scratch/wisp/probes/342/v2/logs/logs-01-roster.txt \
  .scratch/wisp/probes/342/v2/logs/logs-02-legdiff.txt \
  .scratch/wisp/probes/342/v2/logs/logs-04-buildtrees.txt \
  .scratch/wisp/probes/342/v2/logs/logs-06-posix-nails.txt \
  .scratch/wisp/probes/342/v2/logs/logs-08-windows-nails.txt \
  .scratch/wisp/probes/342/v2/logs/logs-09-posix-fullpkg.txt \
  .scratch/wisp/probes/342/v2/logs/logs-12-windows-fullpkg.txt \
  .scratch/wisp/probes/342/v2/logs/x01-discipline.txt \
  .scratch/wisp/probes/342/v2/logs/x02-ticket-tails.txt \
  .scratch/wisp/probes/342/v2/logs/x03-banned-split.txt \
  .scratch/wisp/probes/342/v2/sh/20-commit-batch.sh \
  .scratch/wisp/probes/342/v2/sh/21-ticket-tails-append.sh \
  .scratch/wisp/probes/342/v2/sh/22-banned-codepoint-split.sh \
  .scratch/wisp/probes/348/v1/sh/21-ticket-tails-append.sh \
  .scratch/wisp/probes/348/v1/sh/22-banned-codepoint-split.sh
echo "RC=$?"
