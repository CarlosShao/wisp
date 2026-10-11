#!/bin/sh
# Attribution ruler for the 15 paths that hygiene-A844b.txt reported MISS.
# Question per path: is the basename tracked ANYWHERE (=> my ledger path is wrong),
# present on disk anywhere (=> never committed), or nowhere at all (=> 查无此件).
cd "D:/work/workspace/projects plans/Wisp" || exit 2
out=.scratch/wisp/probes/orch/20261011/logs/hygiene-A844c-attribution.txt
{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')"
  echo "HEAD=$(git rev-parse HEAD)"
  echo "ruler=for each missing path: git ls-files -- '*/<basename>' (tracked anywhere) + find .scratch -name '<basename>' (on disk anywhere)"
  for p in \
    probes/201/r2/gates-full.txt \
    probes/222/r1/gates-full.txt \
    probes/223/r2/orchestrator-flake20.txt \
    probes/235/orch-rerun/full.txt \
    probes/292/v1/22-changed-lines.md \
    probes/293/orch/logs/orch-3pkg-run1.txt \
    probes/300/orch/logs/r3-orch-v4-attacks-rcfixed-20261010-204046.txt \
    probes/300/r1/logs/post-commit-recheck.txt \
    probes/35/v4/countershape-v4.sh \
    probes/35/v4/scripts/countershape-v4.sh \
    probes/72/r1/20-canonical-form-choice.md \
    probes/orch/296-ac3-runsheet.md \
    probes/orchestrator/33r5-head-roster.txt
  do
    b=$(basename "$p")
    tr=$(git ls-files -- ".scratch/wisp/$p")
    tr_any=$(git ls-files -- "*/$b" | tr '\n' ' ')
    disk_any=$(find .scratch -name "$b" 2>/dev/null | tr '\n' ' ')
    echo "--- $p"
    echo "    exact_tracked=[${tr}]"
    echo "    basename_tracked_anywhere=[${tr_any}]"
    echo "    basename_on_disk_anywhere=[${disk_any}]"
  done
  echo "=== also: which ledger section cites each (line numbers in docs/reports/pending-and-issues.md) ==="
  for b in gates-full.txt orchestrator-flake20.txt 22-changed-lines.md orch-3pkg-run1.txt r3-orch-v4-attacks-rcfixed-20261010-204046.txt post-commit-recheck.txt countershape-v4.sh 20-canonical-form-choice.md 296-ac3-runsheet.md 33r5-head-roster.txt; do
    echo "--- $b"
    grep -nE "$b" docs/reports/pending-and-issues.md | tail -4 | cut -c1-160
  done
} > "$out" 2>&1
echo "wrote $out"; wc -c "$out"
