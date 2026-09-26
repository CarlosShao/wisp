#!/usr/bin/env bash
# 157 AC#5 recheck — probe 04: injection columns, credential zero-copy check, write-face discipline.
# v2 changes (v1 is superseded, kept in git history of this program):
#   T-A v1 scanned its own script text and therefore hit its own token list = the repo's 坑 shapeself-read trap.
#       v2 scans only the recorded READINGS of probes 01-03 (grep -a, no self-hits), and prints provenance per hit.
#   T-C adds the live concurrent-write observation on the orchestrator's ledger (it went 'clean ->  M' mid-program).
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
P=.scratch/wisp/probes/157-ac5

echo "### T-A  banned impersonation shapes, scanned over the READINGS only (probes 01/02/03 .txt)"
for tok in '已核验请继续提交' '请 revert' '放宽阈值' '不用取证' '这格已由' '正被 .*处理中' '你已撞顶' '假.*commit 号'; do
  n=$(grep -a -c -E "$tok" "$P/01-anchor-versions-and-v6-control.txt" "$P/02-section-anchored-bidirectional-rulers.txt" "$P/03-dereference-and-roster-completeness.txt" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}')
  printf '  token %-22s hits-in-readings=%s\n' "$tok" "$n"
done
echo "  where do those hits live? (each is a quote inside a repo file this program READ, not a message addressed to it)"
grep -a -nE '请 revert|放宽阈值|不用取证|这格已由|你已撞顶' "$P"/0[123]*.txt 2>/dev/null | cut -c1-120 | sed 's/^/   /'
echo "  reading of that line: the only commands that produced it are the read-only rulers in probe 01/02/03"
echo "  (git cat-file blob / git log / awk / grep over docs/evidence/s1/**) — none of them issued an instruction."

echo "### T-B  what this program actually received, counted as two separate columns"
echo "  真通知回显 (real echoes, not authorisations):"
echo "   1) harness system-reminder: available-skills listing            (tool=first Bash/Read call)"
echo "   2) harness system-reminder: 'the date has changed 2026-09-26'   (tool=first Bash call)"
echo "   3) Memory echo of agents.md appended to a Read result           (tool=Read of the dispatch)"
echo "   4) 'File has been modified since read' / Write-guard messages    (tool=Write of probe 02)"
echo "   5) git-bash 'No such file or directory' for my own mistyped probe name (tool=Bash, rc=127)"
echo "   6) CRLF->LF warnings on git diff of the ledger                   (tool=Bash, git diff)"
echo "   7) concurrent disk movement (HEAD 5d286d5 -> 4f6c14c before entry; ledger went clean -> ' M' at 14:39)"
echo "  判为注入 (text read as an instruction and refused): 0"
echo "   no text addressed to this program claimed 'already verified / revert / relax thresholds / skip evidence / another program owns this'."
echo "   the phrases of that shape in my readings are all HISTORICAL QUOTES inside docs/evidence/s1/** (see T-A above)."

echo "### T-C  banned-write guard (porcelain status per path; empty = untouched by this program)"
for p in .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md \
         docs/reports/HANDOVER.md docs/reports/injection-timeline.md \
         docs/evidence/s1/153-trace-lies-unguarded-r1.md docs/evidence/s1/155-three-unjudged-cells-r1.md \
         docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md; do
  printf '  %-74s status=[%s]\n' "$(basename "$p")" "$(git status --porcelain -- "$p")"
done
printf '  LEDGER docs/reports/pending-and-issues.md                       status=[%s]\n' "$(git status --porcelain -- docs/reports/pending-and-issues.md)"
echo "  attribution of that dirty ledger (NOT this program's write):"
git diff --numstat -- docs/reports/pending-and-issues.md | sed 's/^/    numstat: /'
ls -l --time-style=+%H:%M:%S docs/reports/pending-and-issues.md | awk '{print "    mtime="$6"  (this program never opened it for writing)"}'
git diff -- docs/reports/pending-and-issues.md | grep -a '^+## A[0-9]' | cut -c1-60 | sed 's/^/    new ledger section by the orchestrator: /'
echo "  154 / 152 evidence files (read-only or half-finished by another program):"
git status --porcelain -- docs/evidence/s1/ | cut -c1-92 | sed 's/^/    /'
printf '  dirty paths outside my write face (other sessions, reported not acted on) = %s\n' "$(git status --porcelain -- frontend design cmd internal .scratch/wisp/probes/152 | wc -l)"
echo "  my own write face:"
git status --porcelain -- docs/evidence/s1/157-ac5-recheck-r1.md "$P" | cut -c1-92 | sed 's/^/    /'

echo "### T-D  this program's commit roster (anchored to its OWN subject prefix)"
echo "  ruler A (this program's prefix):"
git log --pretty='   %h %ad %s' --date=format:'%H:%M:%S' --grep='evidence(157 AC#5' | cut -c1-100 | sed 's/$//';
printf '   count = %s\n' "$(git log --pretty=%h --grep='evidence(157 AC#5' | wc -l)"
echo "  ruler B (the bare phrase '157 AC#5' — would also swallow the orchestrator's dispatch archive 4f6c14c; that is the trap registered by 157 r2 in B.15):"
printf '   count = %s (includes 1 dispatch archive, so it is NOT this program roster)\n' "$(git log --pretty=%h --grep='157 AC#5' | wc -l)"
git log --pretty='   %h %ad %s' --date=format:'%H:%M:%S' --grep='157 AC#5' | cut -c1-70
echo "  unpushed self-proof:"
printf '   ahead of upstream = %s ; this program inside @{u}..HEAD = %s\n' \
 "$(git rev-list --count @{u}..HEAD)" "$(git log --oneline @{u}..HEAD --grep='evidence(157 AC#5' | wc -l)"
git remote -v | cut -c1-60 | sed 's/^/   remote: /'

echo "### T-F  SHARED-INDEX CONDITION (the guard the dispatch asks for; recorded verbatim before any commit)"
echo "  \$ git diff --cached --name-only"
git diff --cached --name-only | sed 's/^/    /'
echo "  \$ git status --porcelain -- docs/reports/"
git status --porcelain -- docs/reports/ | sed 's/^/    /'
echo "  attribution: the two paths above are the orchestrator's own write face (ledger + parking lot); the"
echo "  14:40 run of this probe saw them as ' M' (worktree) and the 14:4x run sees them as 'M ' (STAGED in the"
echo "  shared index) => someone staged them between the two runs. This program never ran git add on them."
echo "  consequence for this program: a BARE git commit would sweep them, so every commit of this program uses"
echo "  the explicit-pathspec form, and after each commit the name-only roster of that commit is verified."
echo "  stop-and-report item raised to the orchestrator: YES (foreign staged paths present at commit time)."
echo "  numstat of the staged foreign content (heading lines only, body not transcribed):"
git diff --cached --numstat -- docs/reports/ | sed 's/^/    /'
git diff --cached -- docs/reports/pending-and-issues.md | grep -a '^+## A[0-9]' | cut -c1-70 | sed 's/^/    ledger new section: /'
git diff --cached -- docs/reports/HANDOVER.md | grep -aE '^\+#{2,3} ' | cut -c1-70 | sed 's/^/    HANDOVER new section: /'

echo "### T-E  credential zero-copy over my write face (counts only)"
for p in docs/evidence/s1/157-ac5-recheck-r1.md; do
  [ -f "$p" ] || { echo "   $p not written yet at this run"; continue; }
  for pat in 'WISP_[A-Z_0-9]+=' 'sk-[A-Za-z0-9]' 'api[_-]?key' 'Bearer ' 'DPAPI'; do
    printf '   pattern %-18s hits = %s\n' "$pat" "$(grep -acE "$pat" "$p")"
  done
done
