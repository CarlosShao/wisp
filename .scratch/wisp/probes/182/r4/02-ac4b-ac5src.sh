# 182 r4 leg — script 02: git_readonly_exec location, guikou rule source, Q-49, AGENTS 1.2
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
echo "== panel git family files =="
ls internal/panel/ | grep -i "git"
echo "== runGitReadOnly across panel package =="
grep -rn "runGitReadOnly" internal/panel/ --include=*.go | grep -v _test.go
echo "== exec.Command in git.go vs git_readonly_exec.go (0 in git.go expected; positive control in ro exec) =="
grep -c "exec.Command" internal/panel/git.go || true
grep -c "exec.Command" internal/panel/git_readonly_exec.go || true
echo "== guikou rule source (README + ledger) =="
grep -rn "归口" .scratch/wisp/issues/README.md 2>/dev/null | head -4
grep -rn "归口" docs/reports/pending-and-issues.md 2>/dev/null | head -4
echo "== Q-49 family in ledger (line numbers only + first lines) =="
grep -n "Q-49" docs/reports/pending-and-issues.md | head -8
echo "== AGENTS.md 1.2 panel-L2-allow ban line =="
grep -n "由面板侧来源" AGENTS.md
echo "== ticket 182 guikou sentence (AC#4 source) =="
grep -n "归口" .scratch/wisp/issues/182-*.md
exit 0
