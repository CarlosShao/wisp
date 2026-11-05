#!/bin/bash
# 339-r2 提交：两处（add / commit）都带**逐枚手写**的 pathspec，⛔ 从共享 index 反推；index.lock 重试。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
MSG=".scratch/wisp/probes/339/r2/msg-commit.txt"
TICKET=".scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md"
PATHS=(
  internal/risk/anchor_miss_reason_334_test.go
  internal/risk/anchor_miss_reason_334_windows_test.go
  internal/risk/anchor_miss_reason_339_other_test.go
  docs/evidence/s1/339-r2-impl.md
  .scratch/wisp/probes/339/r2
  "$TICKET"
)
retry() { # $@ = command
  local i rc
  for i in 1 2 3 4 5 6 7 8; do
    "$@"; rc=$?
    [ "$rc" -eq 0 ] && return 0
    if [ -f .git/index.lock ]; then echo "index.lock present, retry $i"; sleep 2; else return $rc; fi
  done
  return $rc
}
retry git add -- "${PATHS[@]}"; echo "ADD_RC=$?"
echo "== staged (only these must appear) =="
git diff --cached --name-only
echo "== commit =="
retry git commit -F "$MSG" -- "${PATHS[@]}"; echo "COMMIT_RC=$?"
echo "== HEAD =="
git log -1 --stat --oneline | head -12
echo "== post: worktree clean for these paths? =="
git status --porcelain -- "${PATHS[@]}"
