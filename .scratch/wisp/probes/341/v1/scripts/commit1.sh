#!/usr/bin/env bash
# 341-v1 commit helper: explicit literal pathspecs only (no index-derived forms).
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
git add \
  .scratch/wisp/probes/341/v1/logs/00-start-gauges.txt \
  .scratch/wisp/probes/341/v1/logs/01-ac0-census-and-landings.txt \
  .scratch/wisp/probes/341/v1/logs/03-color-base.txt \
  .scratch/wisp/probes/341/v1/logs/04-color-deliv.txt \
  .scratch/wisp/probes/341/v1/logs/05-controls.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p2.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p3.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p5.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p6.txt \
  .scratch/wisp/probes/341/v1/scripts/run1-trees.sh \
  .scratch/wisp/probes/341/v1/scripts/run2-color.sh \
  .scratch/wisp/probes/341/v1/scripts/run3-mu.sh \
  .scratch/wisp/probes/341/v1/scripts/run5-controls.sh
echo "ADD_RC=$?"
git commit -F 'D:\tmp\wisp\20261018\341v1\msg-1-ac1.txt' -- \
  .scratch/wisp/probes/341/v1/logs/00-start-gauges.txt \
  .scratch/wisp/probes/341/v1/logs/01-ac0-census-and-landings.txt \
  .scratch/wisp/probes/341/v1/logs/03-color-base.txt \
  .scratch/wisp/probes/341/v1/logs/04-color-deliv.txt \
  .scratch/wisp/probes/341/v1/logs/05-controls.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p2.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p3.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p5.txt \
  .scratch/wisp/probes/341/v1/logs/10-mu-p6.txt \
  .scratch/wisp/probes/341/v1/scripts/run1-trees.sh \
  .scratch/wisp/probes/341/v1/scripts/run2-color.sh \
  .scratch/wisp/probes/341/v1/scripts/run3-mu.sh \
  .scratch/wisp/probes/341/v1/scripts/run5-controls.sh
echo "COMMIT_RC=$?"
git log -1 --format='%H %s' | cut -c1-120
git show --name-only --format='COMMIT %h' HEAD | tail -20
echo "=== status of my write面 after commit ==="; git status --porcelain -- .scratch/wisp/probes/341/v1 docs/evidence/s1/341-v1-verdict.md
