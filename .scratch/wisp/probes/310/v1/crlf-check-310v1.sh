#!/usr/bin/env bash
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
T=/d/tmp/310v1-tree.UcNNmC
D=.scratch/wisp/probes/310/v1/logs
: > "$D/30-crlf-check2.txt"
say() { echo "$@" >> "$D/30-crlf-check2.txt"; }
for f in cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go; do
  wc_=$(tr -cd '\r' < "$f" | wc -c | tr -d ' ')
  ac_=$(tr -cd '\r' < "$T/$f" | wc -c | tr -d ' ')
  bc_=$(git show "HEAD:$f" | tr -cd '\r' | wc -c | tr -d ' ')
  inA=$( (cd "$T" && gofmt -l "$f") 2>&1 | head -1)
  inW=$(gofmt -l "$f" 2>&1 | head -1)
  say "file=$f worktree_CR_bytes=$wc_ archive_CR_bytes=$ac_ blob_CR_bytes=$bc_ gofmt_flags_worktree=[$inW] gofmt_flags_archive=[$inA]"
done
say "gitattributes_exists=$([ -f .gitattributes ] && echo yes || echo no)"
if [ -f .gitattributes ]; then say "gitattributes_eol_lines=$(grep -c -E 'eol|text' .gitattributes || true)"; grep -n -E 'eol|text' .gitattributes | head -6 > /dev/null; fi
say "core_autocrlf=$(git config core.autocrlf || echo unset) core_eol=$(git config core.eol || echo unset)"
say "gofmt_d_first8:"
gofmt -d internal/agent/tools.go 2>&1 | head -8 >> "$D/30-crlf-check2.txt"
say "crlf_check_done=1"
