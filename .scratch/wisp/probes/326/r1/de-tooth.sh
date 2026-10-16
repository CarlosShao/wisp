#!/usr/bin/env bash
# 326-r1 去牙反控＋反形突变：全部打在仓外导出树 /d/tmp/w326/after-86298e2 的
# 盘上字节（⛔ -overlay、⛔ 动共享工作树），每一发先记哈希、跑完还原再跑一发。
set -u
TREE=/d/tmp/w326/after-86298e2
SRC="$TREE/internal/memory/dao_session.go"
OUT=/d/tmp/w326/mutations
LOG="$OUT/all.txt"
mkdir -p "$OUT"
cd "$TREE" || exit 1
export GOTMPDIR='C:\Users\swq\AppData\Local\Temp\wisp326r1'
export PATH="D:\\work\\workspace\\projects plans\\Wisp\\third_party\\sherpa-onnx:$PATH"

PRISTINE="$OUT/dao_session.go.pristine"
cp "$SRC" "$PRISTINE"

restore() {
  cp "$PRISTINE" "$SRC"
  echo "RESTORE sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
}

run_pair() {
  tag="$1"
  echo "----------------------------------------------------------------"
  echo "#### $tag"
  echo "when=$(date '+%F %T %z') sha256(dao_session.go)=$(sha256sum "$SRC" | cut -d' ' -f1)"
  echo "\$ go test ./internal/memory/ -run 'TestAppendSessionTurn(SurfacesTheConversation|RefreshUsesStoredClockFamily|CrashAtomicAcrossKill)' -count=1 -v"
  go test ./internal/memory/ -count=1 -v \
    -run 'TestAppendSessionTurn(SurfacesTheConversation|RefreshUsesStoredClockFamily|CrashAtomicAcrossKill)' \
    > "$OUT/$tag.memory.txt" 2>&1
  echo "rc=$?"
  grep -vE 'level=(INFO|WARN|ERROR)' "$OUT/$tag.memory.txt"
  echo "\$ go test ./internal/session/ -run 'TestSessionAppendTurnMovesItBackToTheFront' -count=1 -v"
  go test ./internal/session/ -count=1 -v -run 'TestSessionAppendTurnMovesItBackToTheFront' \
    > "$OUT/$tag.session.txt" 2>&1
  echo "rc=$?"
  grep -vE 'level=(INFO|WARN|ERROR)' "$OUT/$tag.session.txt"
}

echo "=== 326-r1 去牙反控＋反形台件（$(date '+%F %T %z')）==="
echo "tree=$TREE  原始 sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
echo "正控（未突变）先跑一发，证明这棵树本身是绿的："
run_pair M00-delivered

# M01 去牙：物理摘掉那条 UPDATE（整行删除）
sed -i '/UPDATE session SET updated_at=? WHERE id=?/d' "$SRC"
grep -c 'UPDATE session SET updated_at' "$SRC"
run_pair M01-de-tooth-no-update
restore

# M02 反形：刷新成常量 0
sed -i 's/`UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt/`UPDATE session SET updated_at=? WHERE id=?`, 0/' "$SRC"
grep -c 'updated_at=? WHERE id=?., 0,' "$SRC"
run_pair M02-constant-zero
restore

# M03 反形：刷新成 created_at（等于永不前进）
sed -i 's/`UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt/`UPDATE session SET updated_at=created_at WHERE id=?/' "$SRC"
grep -c 'SET updated_at=created_at' "$SRC"
run_pair M03-to-created-at
restore

# M04 反形：刷新挪到事务外（另起一枚 s.write）
python - "$SRC" <<'PY'
import io, sys
p = sys.argv[1]
s = io.open(p, encoding='utf-8').read()
inside = "\t\t_, err = tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID)\n\t\treturn err\n"
after = "\t\treturn nil\n"
assert inside in s, "anchor not found"
s = s.replace(inside, after)
tail = "\tif err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: append session turn: %w\", err)\n\t}\n\treturn id, nil\n"
moved = ("\tif err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: append session turn: %w\", err)\n\t}\n"
         "\t// M04 mutation: the refresh lives in its OWN transaction, after the turn's.\n"
         "\tif err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {\n"
         "\t\t_, err := tx.ExecContext(ctx, `UPDATE session SET updated_at=? WHERE id=?`, t.CreatedAt, t.SessionID)\n"
         "\t\treturn err\n"
         "\t}); err != nil {\n\t\treturn 0, fmt.Errorf(\"memory: touch session: %w\", err)\n\t}\n"
         "\treturn id, nil\n")
assert tail in s, "tail not found"
s = s.replace(tail, moved)
io.open(p, 'w', encoding='utf-8', newline='\n').write(s)
print("M04 written")
PY
run_pair M04-outside-txn
restore

echo "----------------------------------------------------------------"
echo "=== 还原复跑（交付字节）==="
run_pair M99-restored
echo "cmp 对 pristine："; cmp "$SRC" "$PRISTINE"; echo "cmp rc=$? (0＝逐字节同)"
echo "还原后 sha256=$(sha256sum "$SRC" | cut -d' ' -f1)"
git -C "D:/work/workspace/projects plans/Wisp" show 86298e24:internal/memory/dao_session.go | sha256sum
echo "=== 全部突变结束 ==="
