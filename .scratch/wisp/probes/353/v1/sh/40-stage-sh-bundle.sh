#!/bin/sh
# 353-v1 落尺包＝把本腿全部尺脚本从仓外临时根复制进写面 .scratch/wisp/probes/353/v1/sh/（只建不删、⛔ 动任何 .go）
SRC=/d/tmp/wisp/20261020/353v1/sh
DST='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/353/v1/sh'
mkdir -p "$DST"
for f in 01-anchors.sh 02-blob-dump-diff.sh 03-diff-rest.sh 10-ac2-pair-recheck.sh 11-mutation-recheck.sh 12-ac3-color-recheck.sh 13-teeth-control.sh 14-ticket-face-check.sh 15-scope-rulers.sh 16-secC-recheck.sh 17-familyA-numerators.sh 18-landing-positions.sh 19-full-table.sh 20-misc-rulers.sh 21-misc-part2.sh 21b-misc.sh 22-ancestorwindow-and-swing.sh 23-mutation-trees.sh 24-secC-semantics.sh 25-crossbatch-newred.sh 26-alt-blank-and-mother-ticket.sh 27-precedent-322.sh 28-precedent-and-misc.sh 30-body-rulers.sh 31-find-sect.sh 32-trunc-side-count.sh; do
  if [ -f "$SRC/$f" ]; then cp -f "$SRC/$f" "$DST/$f"; echo "COPIED $f $(wc -c < "$DST/$f")B"; else echo "MISSING $f"; fi
done
echo "== 另附一枚读数原件（非 0 字节＝truncated 非零全表，尺 R32）=="
cp -f /d/tmp/wisp/20261020/353v1/trunc-rows.txt "$DST/../90-truncated-nonzero-rows.txt" && wc -c -l < "$DST/../90-truncated-nonzero-rows.txt"
