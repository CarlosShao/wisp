#!/usr/bin/env bash
# 326-r1 第二批突变：M03 用 python 精确替换重做（第一发因 sed 把行尾反引号吃掉
# ⇒ [build failed]＝台面故障、⛔ 计入判据），M04 加命中率（真杀台件 ×4 轮＝32 枚杀点）。
set -u
TREE=/d/tmp/w326/after-86298e2
SRC="$TREE/internal/memory/dao_session.go"
OUT=/d/tmp/w326/mutations
mkdir -p "$OUT"
cd "$TREE" || exit 1
export GOTMPDIR='C:\Users\swq\AppData\Local\Temp\wisp326r1'
export PATH="D:\\work\\workspace\\projects plans\\Wisp\\third_party\\sherpa-onnx:$PATH"
PRISTINE="$OUT/dao_session.go.pristine"
DELIVERED=ec6ee530b190795093829eca07e6003bf866db5da3e64c47e6d0b0b1bcd1d2b2

restore() { cp "$PRISTINE" "$SRC"; echo "RESTORE sha256=$(sha256sum "$SRC" | cut -d' ' -f1) (want $DELIVERED)"; }

echo "=== 326-r1 第二批突变（$(date '+%F %T %z')）==="
cp "$SRC" "$PRISTINE"; sha256sum "$PRISTINE"

echo "#### M03b-to-created-at（刷新成母表自己的 created_at）"
python - "$SRC" <<'PY'
import io, sys
p = sys.argv[1]
s = io.open(p, encoding='utf-8').read()
a = "`UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID"
b = "`UPDATE session SET updated_at=created_at WHERE id=?`, t.SessionID"
assert s.count(a) == 1, s.count(a)
io.open(p, 'w', encoding='utf-8', newline='\n').write(s.replace(a, b))
print("M03b written")
PY
echo "sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
go test ./internal/memory/ -count=1 -v -run 'TestAppendSessionTurn(SurfacesTheConversation|RefreshUsesStoredClockFamily|CrashAtomicAcrossKill)' > "$OUT/M03b.memory.txt" 2>&1
echo "memory rc=$?"; grep -vE 'level=(INFO|WARN|ERROR)' "$OUT/M03b.memory.txt"
go test ./internal/session/ -count=1 -v -run 'TestSessionAppendTurnMovesItBackToTheFront' > "$OUT/M03b.session.txt" 2>&1
echo "session rc=$?"; grep -vE 'level=(INFO|WARN|ERROR)' "$OUT/M03b.session.txt"
restore

echo "#### M04b-outside-txn 命中率（真杀台件 4 轮＝32 枚杀点）"
python - "$SRC" <<'PY'
import io, sys
p = sys.argv[1]
s = io.open(p, encoding='utf-8').read()
inside = "\t\t_, err = tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID)\n\t\treturn err\n"
assert inside in s
s = s.replace(inside, "\t\treturn nil\n")
tail = "\tif err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: append session turn: %w\", err)\n\t}\n\treturn id, nil\n"
assert tail in s
moved = ("\tif err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: append session turn: %w\", err)\n\t}\n"
         "\t// M04 mutation: the refresh lives in its OWN transaction, after the turn's.\n"
         "\tif err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {\n"
         "\t\t_, err := tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID)\n"
         "\t\treturn err\n"
         "\t}); err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: touch session: %w\", err)\n\t}\n"
         "\treturn id, nil\n")
io.open(p, 'w', encoding='utf-8', newline='\n').write(s.replace(tail, moved))
print("M04b written")
PY
echo "sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
go test ./internal/memory/ -count=4 -v -run 'TestAppendSessionTurnCrashAtomicAcrossKill' > "$OUT/M04b.crash.txt" 2>&1
echo "crash rc=$?"
printf 'M04b 杀点总数 = '; grep -c '326 AC#3 kill@' "$OUT/M04b.crash.txt"
printf 'M04b HALF-STATE 命中 = '; grep -c '=> HALF-STATE' "$OUT/M04b.crash.txt"
grep -E '^--- (FAIL|PASS)|^FAIL|^ok ' "$OUT/M04b.crash.txt"
go test ./internal/session/ -count=1 -v -run 'TestSessionAppendTurnMovesItBackToTheFront' > "$OUT/M04b.session.txt" 2>&1
echo "session(AC#1 出口面) rc=$?"; grep -E '^--- (FAIL|PASS)|^FAIL|^ok ' "$OUT/M04b.session.txt"
restore

echo "#### 交付字节同轮数对照（⛔ 假红：同一台件跑 4 轮应 0 命中）"
echo "sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
go test ./internal/memory/ -count=4 -v -run 'TestAppendSessionTurnCrashAtomicAcrossKill' > "$OUT/M99.crash4.txt" 2>&1
echo "crash rc=$?"
printf '交付态杀点总数 = '; grep -c '326 AC#3 kill@' "$OUT/M99.crash4.txt"
printf '交付态 HALF-STATE 命中 = '; grep -c '=> HALF-STATE' "$OUT/M99.crash4.txt"; echo "grep rc=$? (1＝零命中)"
grep -E '^--- (FAIL|PASS)|^FAIL|^ok ' "$OUT/M99.crash4.txt"
cmp "$SRC" "$PRISTINE"; echo "cmp rc=$?"
echo "还原后 sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
git -C "D:/work/workspace/projects plans/Wisp" show 86298e24:internal/memory/dao_session.go | sha256sum
