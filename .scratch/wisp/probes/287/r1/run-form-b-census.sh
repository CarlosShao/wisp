#!/bin/sh
# 287-r1 形乙普查台件（票 287 AC#3 的代价表用读数，⛔ 不改产码＝物理改写＋sha256 还原）。
#
# 问的是："把结构那一半（名字闸／目录过滤）摘掉，深度那一支就会真在生产链上被走到"
# 这一形要付多少既有仪器。两发：
#   B1 = 只摘名字闸（孩子的请求打到桥 ⇒ :276 成为实际拦下者）
#   B2 = 名字闸＋目录过滤一起摘（孩子还看得见 task.spawn 这个名字）
# 跑整包（⛔ 只跑定向档：代价表问的就是"既有仪器红几枚"）。-timeout 兜住挂死，
# 挂死本身也是代价，具名登记。
set -u

repo="D:/work/workspace/projects plans/Wisp"
cd "$repo" || exit 2
SRC=internal/tools/subagent_197.go
PROBE=.scratch/wisp/probes/287/r1
BAK="$PROBE/subagent_197.go.pre-mutation.bak"
mkdir -p "$PROBE/mutants"

export PATH="$repo/third_party/sherpa-onnx:$PATH"
export GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\287r1gotmp'
mkdir -p /c/Users/swq/AppData/Local/wisp-test-gotmp/287r1gotmp || exit 2

BASE=$(sha256sum "$BAK" | cut -d' ' -f1)
echo "BASELINE-SHA $BASE"

python - "$BAK" "$PROBE/mutants" <<'PY' || exit 3
import pathlib, sys
bak = pathlib.Path(sys.argv[1])
out = pathlib.Path(sys.argv[2])
text = bak.read_bytes()
GATE = b"\tif req.Name == subagentHiddenTool {"
LIST = b"\t\tif info.Name == subagentHiddenTool {"
subs = {
    "B1-name-gate-off": ([(GATE, b"\tif false && req.Name == subagentHiddenTool {")],
                         b"if false && req.Name == subagentHiddenTool {"),
    "B2-name-gate-and-listing-off": (
        [(GATE, b"\tif false && req.Name == subagentHiddenTool {"),
         (LIST, b"\t\tif false && info.Name == subagentHiddenTool {")],
        b"if false && info.Name == subagentHiddenTool {"),
}
for name, (pairs, lands) in subs.items():
    mutated = text
    for a, b in pairs:
        if text.count(a) != 1:
            print("ANCHOR-NOT-UNIQUE %s hits=%d for %r - refuse to seed" % (name, text.count(a), a))
            sys.exit(3)
        mutated = mutated.replace(a, b)
    if mutated.count(lands) != 1:
        print("LANDS-CHECK-FAILED %s" % name)
        sys.exit(3)
    changed = sum(1 for x, y in zip(text.split(b"\n"), mutated.split(b"\n")) if x != y)
    (out / (name + ".go")).write_bytes(mutated)
    print("WROTE %s line_changes=%d" % (name, changed))
PY

for M in B1-name-gate-off B2-name-gate-and-listing-off; do
    echo ""
    echo "############ $M（整包普查）############"
    cp "$PROBE/mutants/$M.go" "$SRC"
    echo "-- 种后 sha256 --"; sha256sum "$SRC"
    echo "-- 内容锚逐枚现量 --"; grep -n "false && req.Name == subagentHiddenTool\|false && info.Name == subagentHiddenTool" "$SRC"
    go test ./internal/tools/ -count=1 -timeout 180s -v > "$PROBE/287r1-$M.log" 2>&1
    rc=$?
    echo "-- go test ./internal/tools/ -count=1 -timeout 180s -v  rc=$rc"
    echo "-- PASS=$(grep -c '^--- PASS' "$PROBE/287r1-$M.log") FAIL=$(grep -c '^--- FAIL' "$PROBE/287r1-$M.log") SKIP=$(grep -c '^--- SKIP' "$PROBE/287r1-$M.log" || true)"
    echo "-- 红名册逐枚 --"
    grep '^--- FAIL' "$PROBE/287r1-$M.log"
    echo "-- 287 本票三枚的颜色 --"
    grep -E '^(--- PASS|--- FAIL): Test287' "$PROBE/287r1-$M.log"
    echo "-- 挂死证据（有则具名）--"
    grep -E 'panic: test timed out|^FAIL$|\(timeout\)' "$PROBE/287r1-$M.log" | head -4
    cp "$BAK" "$SRC"
    NOW=$(sha256sum "$SRC" | cut -d' ' -f1)
    if [ "$NOW" = "$BASE" ]; then echo "RESTORED-EQUAL-OK $NOW"; else echo "RESTORE-MISMATCH-FAIL"; exit 4; fi
done

echo ""
echo "-- 终盘 sha256 --"; sha256sum "$SRC"
echo "-- 终盘 git status --porcelain -- internal cmd tools --"; git status --porcelain -- internal cmd tools
echo "DONE"
