#!/usr/bin/env bash
# 编排者自己复跑票 327 的成对两发（⛔ 抄 327-v1 的数）
# 台面＝仓外导出树；两发各打一侧、还原逐发报哈希。
set -u
ROOT="/d/work/workspace/projects plans/Wisp"
ANCHOR="$(cd "$ROOT" && git rev-parse HEAD)"
TREE_ROOT="/d/work/workspace/wisp-327v1-orch"
TREE="$TREE_ROOT/repo"
GOTMP="$TREE_ROOT/gotmp"
LOG="$ROOT/.scratch/wisp/probes/orch/20261017/txt/orch-327v1-replay.log"
TARGET="internal/panel/display_scrub_216.go"

mkdir -p "$GOTMP"
TREE_ROOT="/d/work/workspace/wisp-327v1-orch-$(date +%H%M%S)"
TREE="$TREE_ROOT/repo"
GOTMP="$TREE_ROOT/gotmp"
mkdir -p "$TREE_ROOT" "$GOTMP" "$TREE"
exec >"$LOG" 2>&1
echo "anchor=$ANCHOR"
cd "$ROOT" || exit 1

git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$TREE"
echo "go_files=$(find "$TREE" -name '*.go' | wc -l)"
cd "$TREE" || exit 1
export GOTMPDIR="$GOTMP"
export GOPROXY=off
ORIG_SHA="$(sha256sum "$TARGET" | cut -d' ' -f1)"
echo "orig_sha256=$ORIG_SHA"

run_pair() {
  local name="$1"
  echo "### RUN $name"
  PATH="/d/work/base/gopath/bin:$PATH" go test ./internal/panel/ -count=1 -timeout 30m -v \
    -run 'TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired|TestNeedsDisplayWashNamesEveryWashedCfRune|TestScrubForDisplay|TestTheSnapshotBoundary' \
    > "$TREE_ROOT/$name.log" 2>&1
  echo "rc=$?"
  echo "RUN=$(grep -cE '^=== RUN' "$TREE_ROOT/$name.log")"
  echo "PASS_TOP=$(grep -cE '^--- PASS' "$TREE_ROOT/$name.log")"
  echo "FAIL_TOP=$(grep -cE '^--- FAIL' "$TREE_ROOT/$name.log")"
  echo "SKIP=$(grep -cE '^--- SKIP' "$TREE_ROOT/$name.log")"
  grep -E '^--- FAIL' "$TREE_ROOT/$name.log" | head -8
  grep -E 'display_scrub_216_test\.go:[0-9]+:' "$TREE_ROOT/$name.log" | head -4
}

restore() {
  cd "$TREE" || exit 1
  git -C "$ROOT" show "$ANCHOR:$TARGET" > "$TARGET"
  local now="$(sha256sum "$TARGET" | cut -d' ' -f1)"
  if [ "$now" = "$ORIG_SHA" ]; then echo "RESTORE_OK sha=$now"; else echo "RESTORE_FAIL now=$now want=$ORIG_SHA"; fi
}

# 发〇：基线（未突变）
run_pair baseline
restore

# 发一：只摘写入侧那 7 行支（Cf 落 default ⇒ 原样穿过）
python - "$TARGET" <<'PY'
import sys, io
p = sys.argv[1]
src = io.open(p, encoding='utf-8').read()
start = src.index('\t\tcase isDisplayWashCf(r):\n')
marker = '\t\t\ti += size\n'
end = src.index(marker, start) + len(marker)
block = src[start:end]
assert block.count('\n') == 7, 'LANDS_wrong_linecount=%d' % block.count('\n')
assert 'writeControlEscape(&b, r)' in block and 'Ticket 327' in block
io.open(p + '.mut1', 'w', encoding='utf-8', newline='').write(src[:start] + src[end:])
print('MUT1_LANDS lines=7')
PY
cp "$TARGET.mut1" "$TARGET"
run_pair m1_writer_branch_removed
restore

# 发二：只摘判据侧那枚 term（fast path 不再认 Cf）
python - "$TARGET" <<'PY'
import sys, io
p = sys.argv[1]
src = io.open(p, encoding='utf-8').read()
old = 'if (r == utf8.RuneError && size == 1) || unicode.IsControl(r) || isDisplayWashCf(r) {'
new = 'if (r == utf8.RuneError && size == 1) || unicode.IsControl(r) {'
assert src.count(old) == 1, 'LANDS_wrong_hits=%d' % src.count(old)
io.open(p, 'w', encoding='utf-8', newline='').write(src.replace(old, new))
print('MUT2_LANDS hits=1')
PY
run_pair m2_judge_term_removed
restore

# 收尾：工作树里那枚跟踪文件⛔ 被动过（本脚本只在仓外动字节）
echo "worktree_target_status=$(git -C "$ROOT" status --porcelain -- "$TARGET" | wc -l)"
echo "mutated_copies_left=$TREE_ROOT（只建不删）"
