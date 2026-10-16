#!/usr/bin/env bash
# 287-v1b 仓外导出树：建一次、验六件套、记基线哈希与 CR 数、跑基线读数
# 用法：sh bench-setup.sh <anchor>   （⛔ 删任何已建的树）
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${1:?anchor required}"
STAMP="$(date +%Y%m%d-%H%M%S)"
TREE="/d/work/tmp/287v1b/tree-$STAMP"
OUT="$REPO/.scratch/wisp/probes/287/v1b"
mkdir -p "$TREE"
echo "TREE=$TREE"
echo "ANCHOR=$ANCHOR"
cd "$REPO" || exit 99
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts | tar -x -C "$TREE"
echo "archive_rc=$?"
for p in go.mod go.sum internal cmd tools scripts; do
  if [ -e "$TREE/$p" ]; then echo "piece OK  $p"; else echo "piece MISSING $p"; fi
done
echo "gofiles_total=$(find "$TREE" -name '*.go' | wc -l)"
echo "gofiles_internal=$(find "$TREE/internal" -name '*.go' | wc -l)"
echo "gofiles_cmd=$(find "$TREE/cmd" -name '*.go' | wc -l)"
echo "gofiles_tools=$(find "$TREE/tools" -name '*.go' | wc -l)"
for f in internal/tools/subagent_197.go internal/tools/task.go internal/agent/loop.go \
         internal/tools/bridge.go internal/tools/subagent_depth_fixture_287_test.go; do
  printf 'base %s sha256=%s CR=%s bytes=%s\n' "$f" \
    "$(sha256sum "$TREE/$f" | cut -d' ' -f1)" \
    "$(tr -cd '\r' <"$TREE/$f" | wc -c | tr -d ' ')" \
    "$(wc -c <"$TREE/$f" | tr -d ' ')"
done
echo "repo_worktree_matches=$(sha256sum "$REPO/internal/tools/subagent_197.go" | cut -d' ' -f1) vs $(sha256sum "$TREE/internal/tools/subagent_197.go" | cut -d' ' -f1)"
echo "$TREE" >"$OUT/.tree-path"
