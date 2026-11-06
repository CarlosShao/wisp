# 182 r4 leg — script 05: AC#6 contract axis + AC#7 scheduling
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "== AC#6: working-tree status on contract/production axes (expect only known untracked test files) =="
git status --porcelain -- docs/PLAN.md docs/specs/ internal/ cmd/ frontend/ design/
echo "-- end status --"
echo "== AC#6: ledger untouched in working tree =="
git status --porcelain -- docs/reports/pending-and-issues.md
echo "(empty above = clean)"
echo "== AC#6: FACE file fresh? =="
ls docs/evidence/s1/182-ac4-ac7.md 2>/dev/null || echo "NOT-EXIST (fresh create)"
echo "== AC#6: my leg's write surface this session (probes 182/r4) =="
ls .scratch/wisp/probes/182/r4/
echo "== AC#7: chain ticket files + done-suffix =="
for t in 179 176 177 175; do ls .scratch/wisp/issues/${t}-*.md | sed 's/.*\///'; done
echo "== AC#7: does any chain ticket depend on ticket 182? =="
for f in .scratch/wisp/issues/179-*.md .scratch/wisp/issues/176-*.md .scratch/wisp/issues/177-*.md .scratch/wisp/issues/175-*.md; do
  echo "-- $f --"
  grep -nE "\b182\b" "$f" | head -5 || echo "(no 182 mention)"
done
echo "== AC#7: 182 ticket own scheduling sentence =="
grep -n "179" .scratch/wisp/issues/182-*.md | head -3
exit 0
