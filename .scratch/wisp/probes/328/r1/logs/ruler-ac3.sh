#!/usr/bin/env bash
# AC#3（硬框 ⓮）四面 diff 逐枚报 0：每枚自己的提交 × 每一面，逐格印行数与 rc。
# 起手锚 648eda04；本腿四枚提交＝b2d55342 1696d01b ffebb6d1 f705b082（共享树里别人也在这段落笔，
# 所以尺只能打在"我自己的每一笔"上，作差 whole-range 只当归属读数）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 1
MINE="b2d55342 1696d01b ffebb6d1 f705b082"
FACES="gate-mode:internal/risk/mode.go internal/agent/approval/gate.go internal/agent/approval/queue.go
panel-l2-source:internal/risk/provenance.go internal/panel/approval.go
d43-table:internal/statemachine
frozen-three:internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go"
echo "date=$(date '+%F %T %z')"
for c in $MINE; do
  echo "=== commit $c ($(git log -1 --format='%s' "$c" | cut -c1-40)) ==="
  echo "$FACES" | while IFS= read -r line; do
    name="${line%%:*}"; paths="${line#*:}"
    n=$(git diff --stat "$c^" "$c" -- $paths | wc -l)
    rc=$?
    echo "  face=$name diff_lines=$n rc=$rc paths=[$paths]"
  done
done
echo "=== whole-range attribution: which commits touched each face since 648eda04 ==="
echo "$FACES" | while IFS= read -r line; do
  name="${line%%:*}"; paths="${line#*:}"
  echo "-- face=$name"
  git log --oneline 648eda04..HEAD -- $paths | sed 's/^/   /' || echo "   (no commits)"
done
echo "=== worktree dirty check on those faces (must be empty) ==="
git status --porcelain -- internal/risk internal/statemachine internal/panel/approval.go internal/panel/tokens_fourway_test.go internal/panel/l2_grant_boundary_test.go internal/perm/ticket90_persist_test.go internal/agent/approval
echo "status_rc=$?"
echo "=== my four files vs HEAD (must be clean) ==="
git status --porcelain -- internal/panel/subagent_roster_197.go internal/panel/subagent_blocked_220_test.go internal/panel/subagent_corr_task_prefix_328_test.go internal/agent/corr_roster_join_328_test.go
echo "done rc=0"
