#!/bin/sh
# 353-v1 笔①＝逐枚字面 pathspec（每枚一行、写死全名；⛔ 目录形／add -A／.／从共享 index 反推）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P1=.scratch/wisp/probes/353/v1/90-truncated-nonzero-rows.txt
P2=.scratch/wisp/probes/353/v1/sh/01-anchors.sh
P3=.scratch/wisp/probes/353/v1/sh/02-blob-dump-diff.sh
P4=.scratch/wisp/probes/353/v1/sh/03-diff-rest.sh
P5=.scratch/wisp/probes/353/v1/sh/10-ac2-pair-recheck.sh
P6=.scratch/wisp/probes/353/v1/sh/11-mutation-recheck.sh
P7=.scratch/wisp/probes/353/v1/sh/12-ac3-color-recheck.sh
P8=.scratch/wisp/probes/353/v1/sh/13-teeth-control.sh
P9=.scratch/wisp/probes/353/v1/sh/14-ticket-face-check.sh
P10=.scratch/wisp/probes/353/v1/sh/15-scope-rulers.sh
P11=.scratch/wisp/probes/353/v1/sh/16-secC-recheck.sh
P12=.scratch/wisp/probes/353/v1/sh/17-familyA-numerators.sh
P13=.scratch/wisp/probes/353/v1/sh/18-landing-positions.sh
P14=.scratch/wisp/probes/353/v1/sh/19-full-table.sh
P15=.scratch/wisp/probes/353/v1/sh/20-misc-rulers.sh
P16=.scratch/wisp/probes/353/v1/sh/21-misc-part2.sh
P17=.scratch/wisp/probes/353/v1/sh/21b-misc.sh
P18=.scratch/wisp/probes/353/v1/sh/22-ancestorwindow-and-swing.sh
P19=.scratch/wisp/probes/353/v1/sh/23-mutation-trees.sh
P20=.scratch/wisp/probes/353/v1/sh/24-secC-semantics.sh
P21=.scratch/wisp/probes/353/v1/sh/25-crossbatch-newred.sh
P22=.scratch/wisp/probes/353/v1/sh/26-alt-blank-and-mother-ticket.sh
P23=.scratch/wisp/probes/353/v1/sh/27-precedent-322.sh
P24=.scratch/wisp/probes/353/v1/sh/28-precedent-and-misc.sh
P25=.scratch/wisp/probes/353/v1/sh/30-body-rulers.sh
P26=.scratch/wisp/probes/353/v1/sh/31-find-sect.sh
P27=.scratch/wisp/probes/353/v1/sh/32-trunc-side-count.sh
P28=.scratch/wisp/probes/353/v1/sh/40-stage-sh-bundle.sh
P29=.scratch/wisp/probes/353/v1/sh/41-body-final.sh
P30=.scratch/wisp/probes/353/v1/sh/51-commit-01.sh
git add -- "$P1" "$P2" "$P3" "$P4" "$P5" "$P6" "$P7" "$P8" "$P9" "$P10" "$P11" "$P12" "$P13" "$P14" "$P15" "$P16" "$P17" "$P18" "$P19" "$P20" "$P21" "$P22" "$P23" "$P24" "$P25" "$P26" "$P27" "$P28" "$P29" "$P30"
echo "ADD_RC=$?"
git commit -F 'D:\tmp\wisp\20261020\353v1\msg-01.txt' -- "$P1" "$P2" "$P3" "$P4" "$P5" "$P6" "$P7" "$P8" "$P9" "$P10" "$P11" "$P12" "$P13" "$P14" "$P15" "$P16" "$P17" "$P18" "$P19" "$P20" "$P21" "$P22" "$P23" "$P24" "$P25" "$P26" "$P27" "$P28" "$P29" "$P30" > /d/tmp/wisp/20261020/353v1/commit-out-01.txt 2>&1
echo "COMMIT_RC=$?"
tail -4 /d/tmp/wisp/20261020/353v1/commit-out-01.txt
echo "=== git log -1 ==="
git log -1 --format='%H %s'
echo "=== name-only（本笔）==="
git show --name-only --format='' HEAD | grep -c . | sed 's/^/  件数=/'
git show --name-only --format='' HEAD | grep -c '\.go$' | sed 's/^/  .go 枚数=/'
git show --name-only --format='' HEAD | grep -vcE '^\.scratch/wisp/probes/353/v1/' | sed 's/^/  越界枚数=/'
