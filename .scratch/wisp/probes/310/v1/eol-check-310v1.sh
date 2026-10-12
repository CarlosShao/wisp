#!/usr/bin/env bash
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
D=.scratch/wisp/probes/310/v1/logs
: > "$D/30-blob-eol.txt"
say() { echo "$@" >> "$D/30-blob-eol.txt"; }
for f in cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go; do
  w=$(tr -cd '\r' < "$f" | wc -c | tr -d ' ')
  b=$(git cat-file blob "HEAD:$f" | tr -cd '\r' | wc -c | tr -d ' ')
  a=$(tr -cd '\r' < "/d/tmp/310v1-tree.UcNNmC/$f" | wc -c | tr -d ' ')
  say "file=$f CR_bytes worktree=$w head_blob=$b git_archive_tree=$a"
done
say "gitattributes_content:"
head -12 .gitattributes >> "$D/30-blob-eol.txt"
say "ruler_selfcheck_this_leg_worktree_files_CR=$(tr -cd '\r' < .scratch/wisp/probes/310/r1/90-final.md | wc -c | tr -d ' ')"
say "eol_done=1"
echo "=== ticket bullet commit 3122306e paths ===" >> "$D/30-blob-eol.txt"
git show --name-only --format= 3122306e >> "$D/30-blob-eol.txt" 2>&1
echo "=== r1 zero-byte classification: does the citing md line carry a count ===" >> "$D/30-blob-eol.txt"
grep -n -F -e '10b-new.txt' -e '10b2-new.txt' -e '10b2-root.txt' -e '10c-run.run.txt' -e '10b2-shipped-root.txt' .scratch/wisp/probes/310/r1/*.md >> "$D/30-blob-eol.txt" 2>&1
echo "=== the two files the leg wrote INSIDE the D: dedicated root ===" >> "$D/30-blob-eol.txt"
wc -c /d/tmp/wisp-test-gotmp/shapec-115000-c2_redirected.run.txt /d/tmp/wisp-test-gotmp/shapec-115000-c2_redirected.run.txt.err >> "$D/30-blob-eol.txt" 2>&1
