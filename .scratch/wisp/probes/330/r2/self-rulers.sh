#!/usr/bin/env bash
# 330-r2 self-rulers: "zero production byte" over internal/agent, plus its
# positive control run on a COPY (a --local clone), never in the shared worktree.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
LOG=.scratch/wisp/probes/330/r2/logs
ANCHOR=feab45e3abee588ae30245c820a762d3b695b91f
{
  echo "### ruler 1: git diff --numstat \$ANCHOR HEAD -- internal/agent   (every line must be *_test.go or testdata/)"
  git diff --numstat "$ANCHOR" HEAD -- internal/agent
  echo "numstat_rc=$?"
  echo
  echo "### ruler 2 (the delivery ruler): name-only filtered, must be EMPTY and rc=1"
  git diff --name-only "$ANCHOR" HEAD -- internal/agent | grep -vE '(_test\.go|^internal/agent/testdata/)'
  echo "filtered_rc=$?  (1 = no line survived = zero production byte)"
  echo
  echo "### whole-tree view for the record (other legs are in flight in this worktree; their files are NOT mine)"
  git diff --name-only "$ANCHOR" HEAD | head -30
  echo
  echo "### positive control of ruler 2, run on a --local clone (a fresh dir under /d/tmp, never deleted, never in the shared worktree)"
  POS=$(mktemp -d /d/tmp/w330pos.XXXXXX)
  echo "clone_dir=$POS"
  git clone --local --quiet . "$POS" >/dev/null 2>&1
  echo "clone_rc=$?"
  cd "$POS" || exit 8
  python - <<'PY'
import io
p = "internal/agent/inject.go"
s = io.open(p, encoding="utf-8", newline="").read()
old = "bm25Select(query, candidates, MaxRetrievedTools, b.SectionToolIndex)"
assert s.count(old) == 1, s.count(old)
s = s.replace(old, "bm25Select(query, candidates, MaxRetrievedTools+1, b.SectionToolIndex)")
io.open(p, "w", encoding="utf-8", newline="").write(s)
print("positive-control edit applied to a NON-comment production byte:", p)
PY
  git add -- internal/agent/inject.go
  git commit -q -m "330r2 ruler positive control: one production byte changed in a copy" -- internal/agent/inject.go
  echo "control_commit_rc=$?"
  git diff --name-only "$ANCHOR" HEAD -- internal/agent | grep -vE '(_test\.go|^internal/agent/testdata/)'
  echo "control_filtered_rc=$?  (0 with a listed file = the ruler is armed)"
} 2>&1 | tee "$LOG/self-rulers.log"
