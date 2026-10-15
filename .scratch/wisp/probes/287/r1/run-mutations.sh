#!/bin/sh
# 287-r1 突变台件（票 287 AC#2 的"能红"判据）。
#
# 定式：物理改写 internal/tools/subagent_197.go ＋ sha256 还原，⛔ -overlay
# （本仓 2026-10-15 新定式：overlay 只换编译期输入，对运行时读盘的尺失明）。
# 三枚突变体＋还原前的基线副本一律留盘（临时件只建不删）。
#
# 每一发都自带 LANDS 自证（突变那一行必须真在盘上、且锚点唯一），因为行号锚会
# 腐烂：种坏之前用 grep 找内容锚，⛔ 只信行号。
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

if [ ! -f "$BAK" ]; then
    echo "NO-BACKUP: $BAK 不在盘上，⛔ 不许种坏" >&2
    exit 2
fi
BASE=$(sha256sum "$BAK" | cut -d' ' -f1)
echo "BASELINE-SHA $BASE  ($BAK)"

python - "$BAK" "$PROBE/mutants" <<'PY' || exit 3
import pathlib, sys
bak = pathlib.Path(sys.argv[1])
out = pathlib.Path(sys.argv[2])
text = bak.read_bytes()          # 逐字节：⛔ 文本模式（Windows 会把 LF 翻成 CRLF，
subs = {                         #  那样 diff 的 hunk 数就成了"整枚文件重写"的假账）
    "M1-disable-branch": (
        b"if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {",
        b"if rec, ok := t.d.Roster.Look(parentID); false && ok && rec.Kind == TaskKindSubagent {",
        b"false && ok && rec.Kind == TaskKindSubagent"),
    "M2-identity-to-corr": (
        b"parentID := TaskID(ctx)",
        b"parentID := CorrelationID(ctx)",
        b"parentID := CorrelationID(ctx)"),
    "M3-overbroad-any-row": (
        b"if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {",
        b"if rec, ok := t.d.Roster.Look(parentID); ok {",
        b"Look(parentID); ok {"),
}
for name, (a, b, lands) in subs.items():
    n = text.count(a)
    if n != 1:
        print("ANCHOR-NOT-UNIQUE %s hits=%d - refuse to seed" % (name, n))
        sys.exit(3)
    mutated = text.replace(a, b)
    if mutated.count(lands) != 1:
        print("LANDS-CHECK-FAILED %s" % name)
        sys.exit(3)
    (out / (name + ".go")).write_bytes(mutated)
    print("WROTE %s anchor_hits=%d lands=%s" % (name, n, lands))
PY

for M in M1-disable-branch M2-identity-to-corr M3-overbroad-any-row; do
    echo ""
    echo "############ $M ############"
    echo "-- 落地尺寸（vs 基线）：hunks=$(diff -u "$BAK" "$PROBE/mutants/$M.go" | grep -c '^@@') 改行数=$(diff -u "$BAK" "$PROBE/mutants/$M.go" | grep -c '^+[^+]') --"
    cp "$PROBE/mutants/$M.go" "$SRC"
    echo "-- 种后 sha256 --"; sha256sum "$SRC"
    echo "-- 种后现量（271 / 276）--"; sed -n '271p;276p' "$SRC"
    echo "-- 内容锚命中（LANDS）--"
    case "$M" in
        M1-disable-branch)      grep -c 'false && ok && rec.Kind == TaskKindSubagent' "$SRC" ;;
        M2-identity-to-corr)    grep -c 'parentID := CorrelationID(ctx)' "$SRC" ;;
        M3-overbroad-any-row)   grep -c 'Look(parentID); ok {' "$SRC" ;;
    esac
    go test ./internal/tools/ -count=1 -run 'Test287' -v > "$PROBE/287r1-$M.log" 2>&1
    rc=$?
    echo "-- go test ./internal/tools/ -count=1 -run Test287 -v  rc=$rc"
    grep -E '^(--- (PASS|FAIL|SKIP)|ok|FAIL)' "$PROBE/287r1-$M.log"
    echo "-- 红句逐字（文件:行 + 文案）--"
    grep -E 'subagent_depth_fixture_287_test.go:[0-9]+:' "$PROBE/287r1-$M.log"
    cp "$BAK" "$SRC"
    NOW=$(sha256sum "$SRC" | cut -d' ' -f1)
    echo "-- 还原 sha256 --  $NOW"
    if [ "$NOW" = "$BASE" ]; then echo "RESTORED-EQUAL-OK"; else echo "RESTORE-MISMATCH-FAIL"; exit 4; fi
    echo "-- 还原现量（271 / 276）--"; sed -n '271p;276p' "$SRC"
    echo "-- git status --porcelain -- internal --"; git status --porcelain -- internal
done

echo ""
echo "############ 台面读数（三枚用例，还原后）############"
go test ./internal/tools/ -count=1 -run 'Test287' -v > "$PROBE/287r1-00-targeted-green.log" 2>&1
echo "targeted-green rc=$?"
grep -E '^(--- (PASS|FAIL|SKIP)|ok|FAIL)' "$PROBE/287r1-00-targeted-green.log"
echo "-- 终盘 sha256 --"; sha256sum "$SRC"
echo "-- 终盘 git status --porcelain -- internal cmd --"; git status --porcelain -- internal cmd
echo "DONE"
