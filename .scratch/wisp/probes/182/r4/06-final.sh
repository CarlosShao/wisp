# 182 r4 leg — script 06: final gap-closers
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "== chongshu wording in ledger/AGENTS/handoff =="
grep -rn "充数" docs/reports/pending-and-issues.md AGENTS.md docs/reports/orchestrator-handoff-20261013.md 2>/dev/null | head -4 || echo "(none)"
echo "== docs/ working-tree status (tracked mods expected: none before FACE write) =="
git status --porcelain -- docs/
echo "(empty = clean)"
echo "== r3 evidence file last commit =="
git log --oneline -1 -- docs/evidence/s1/182-census-r3.md
echo "== ticket 182 last commit =="
git log --oneline -1 -- .scratch/wisp/issues/182-*.md
exit 0
