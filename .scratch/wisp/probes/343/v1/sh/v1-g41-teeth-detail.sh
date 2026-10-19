#!/bin/sh
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
R=.scratch/wisp/probes/343/v1/logs/teeth/ac2-sweep-runner.log
O=/d/tmp/wisp/20261019/343v1/logs/final
mkdir -p $O
{
echo "== header line (anchor/pristine/tree = MINE) =="
grep -E '^### AC#2 sweep' $R
echo "== every pristine/restore hash check in order (11 shots + final) =="
grep -E '^  hash=' $R
echo "== 落刀前后命中数 (per mutation) =="
grep -E 'needle_before|branch_removals_on_disk|changed_lines' $R
echo "== red subtest names per mutation (the 8 teeth) =="
grep -E 'RED +--- FAIL: TestAnchorFallbackTeethPerAnchor336/' $R
echo "== final per-mutation counts =="
grep -E '^MU-1' $R
echo "== TEETH336 operand roster in the UNMUTATED positive control =="
grep -E 'TEETH336 operand=' $R | head -8
echo "== runner rc + worktree dirtiness =="
grep -E 'ac2_sweep_rc=|shared_worktree_dirty' $R
echo DONE
} > $O/teeth-detail.log 2>&1
cat $O/teeth-detail.log
