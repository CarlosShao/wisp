#!/usr/bin/env bash
# 303/300 编排者自己的复验（收 300-r2 那一枚新用例）：⛔ 信交件转述，本腿＝编排者本人现跑。
# 台面＝仓外导出树（git archive <HEAD blob> | tar -x），⛔ 在共享工作树里做任何 checkout／改产码。
# 每次跑用一枚新目录名（定式：读数件⛔ 复用文件名）。
set -u

HEADSHA="$(git -C "D:/work/workspace/projects plans/Wisp" rev-parse HEAD)"
STAMP="$(date +%Y%m%d-%H%M%S)"
ROOT="/tmp/300r2-verify-${STAMP}"
OUT="$(dirname "$0")/logs/r1-verify-300r2-${STAMP}.txt"
mkdir -p "$(dirname "$OUT")" "$ROOT/repo"

REPOTXT="D:/work/workspace/projects plans/Wisp"
TESTFILE="internal/audio/wave_format_float_300_windows_test.go"
CONSTFILE="internal/audio/wasapi_windows.go"

exec >"$OUT" 2>&1

echo "# 编排者复验 300-r2（AC#6 的 C2 形）"
echo "captured_local=$(date '+%Y-%m-%d %H:%M:%S %z')"
echo "HEAD_sha=${HEADSHA}"
echo "exported_tree=${ROOT}/repo"
echo "artifact=${OUT}"
echo

echo "## 0. 起手尺：工作树⛔ 动过（导出树＝HEAD 对象层，与盘上脏无关）"
git -C "$REPOTXT" status --porcelain -- internal cmd docs scripts .github | sed 's/^/porcelain: /'
echo "porcelain_lines=$(git -C "$REPOTXT" status --porcelain -- internal cmd docs scripts .github | wc -l)"
echo

echo "## 1. 导出 HEAD 到仓外树"
git -C "$REPOTXT" archive "$HEADSHA" | tar -x -C "$ROOT/repo"
echo "archive_rc=$?"
echo "testfile_present=$(test -f "$ROOT/repo/$TESTFILE" && echo yes || echo no)"
echo "top_level_cases_in_pkg=$(grep -h '^func Test' "$ROOT"/repo/internal/audio/*_test.go | wc -l)"
echo "declare_line_before=$(grep -n 'waveFormatFloat' "$ROOT/repo/$CONSTFILE" | head -1)"
echo

echo "## 2. 导出树里那枚用例与 HEAD blob 逐字节同（cmp）"
git -C "$REPOTXT" show "${HEADSHA}:${TESTFILE}" > "$ROOT/headblob.go"
cmp "$ROOT/headblob.go" "$ROOT/repo/$TESTFILE" && echo "cmp_head_blob=IDENTICAL" || echo "cmp_head_blob=DIFF"
echo

echo "## 3. 门① go vet ./internal/audio/"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go vet ./internal/audio/ )
echo "rc_vet=$?"
echo

echo "## 4. 门② 定向用例 pristine"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 300s -v -run 'TestWaveFormatConstantsMatchMmregAuthority300' ./internal/audio/ )
echo "rc_targeted_pristine=$?"
echo

echo "## 5. 门③ 编排者补跑那一枚：GOFLAGS= go build ./...（腿⛔ 跑、我 17:5x 派单漏列，本程我销账）"
( cd "$ROOT/repo" && GOFLAGS= go build ./... )
echo "rc_build=$?"
echo

echo "## 6. 突变一发：waveFormatFloat 3 -> 1"
cp "$ROOT/repo/$CONSTFILE" "$ROOT/backup-const.go"
sed -i 's/waveFormatFloat[[:space:]]*=[[:space:]]*3\b/waveFormatFloat = 1/' "$ROOT/repo/$CONSTFILE"
echo "declare_line_mutated=$(grep -n 'waveFormatFloat =' "$ROOT/repo/$CONSTFILE" | head -1)"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 300s -v -run 'TestWaveFormatConstantsMatchMmregAuthority300' ./internal/audio/ )
echo "rc_mut_3to1=$?"
cp "$ROOT/backup-const.go" "$ROOT/repo/$CONSTFILE"
cmp "$ROOT/backup-const.go" "$ROOT/repo/$CONSTFILE" && echo "restore_after_3to1=IDENTICAL" || echo "restore_after_3to1=DIFF"
echo

echo "## 7. 突变二发：waveFormatFloat 3 -> 0"
sed -i 's/waveFormatFloat[[:space:]]*=[[:space:]]*3\b/waveFormatFloat = 0/' "$ROOT/repo/$CONSTFILE"
echo "declare_line_mutated=$(grep -n 'waveFormatFloat =' "$ROOT/repo/$CONSTFILE" | head -1)"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 300s -v -run 'TestWaveFormatConstantsMatchMmregAuthority300' ./internal/audio/ )
echo "rc_mut_3to0=$?"
cp "$ROOT/backup-const.go" "$ROOT/repo/$CONSTFILE"
cmp "$ROOT/backup-const.go" "$ROOT/repo/$CONSTFILE" && echo "restore_after_3to0=IDENTICAL" || echo "restore_after_3to0=DIFF"
echo

echo "## 8. 还原后再跑一发 pristine（成对两发的'绿'腿）"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 300s -run 'TestWaveFormatConstantsMatchMmregAuthority300' ./internal/audio/ )
echo "rc_restored=$?"
echo

echo "## 9. internal/audio 整包两发红名册（同一把尺：两发都带 -v；改前＝把那枚新用例从导出树里挪走）"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 600s -v ./internal/audio/ 2>&1 | tee "$ROOT/pkg-after.txt" )
echo "rc_pkg_after=$?"
mv "$ROOT/repo/$TESTFILE" "$ROOT/removed-300r2-test.go"
echo "moved_aside=$(basename "$ROOT/removed-300r2-test.go")"
( cd "$ROOT/repo" && GOFLAGS=-mod=mod go test -count=1 -timeout 600s -v ./internal/audio/ 2>&1 | tee "$ROOT/pkg-before.txt" )
echo "rc_pkg_before=$?"
grep -E '^--- (FAIL|SKIP): ' "$ROOT/pkg-after.txt" 2>/dev/null | sort > "$ROOT/names-after.txt"
grep -E '^--- (FAIL|SKIP): ' "$ROOT/pkg-before.txt" | sort > "$ROOT/names-before.txt"
echo "names_after_lines=$(wc -l < "$ROOT/names-after.txt")"
echo "names_before_lines=$(wc -l < "$ROOT/names-before.txt")"
echo "new_red_or_skip_comm_minus13=$(comm -13 "$ROOT/names-before.txt" "$ROOT/names-after.txt" | wc -l)"
echo "gone_comm_minus23=$(comm -23 "$ROOT/names-before.txt" "$ROOT/names-after.txt" | wc -l)"
echo "before_names_verbatim:"; cat "$ROOT/names-before.txt"
echo "after_names_verbatim:"; cat "$ROOT/names-after.txt"
echo
echo "pkg_before_summary=$(grep -E '^(ok|FAIL|---)' "$ROOT/pkg-before.txt" | grep -E '^(ok|FAIL)' | tail -1)"
echo "pkg_after_summary=$(tail -1 "$ROOT/pkg-after.txt")"
echo
echo "## 10. 格式两把并排（射程目录＝internal/audio；工作树一把＋HEAD blob 一把）"
echo "gofmt_worktree=$(gofmt -l internal/audio | wc -l)"
echo "gofmt_worktree_names: $(cd "$REPOTXT" && gofmt -l internal/audio)"
( cd "$ROOT/repo" && gofmt -l internal/audio ) | sed 's/^/gofmt_headblob: /'
echo
echo "## 11. d22scan（射程＝全仓，跑在真仓里；只读）"
( cd "$REPOTXT" && sh scripts/d22scan.sh )
echo "rc_d22scan=$?"
echo
echo "## 12. census totals 行（逐字对拉基线 35/7/7/0）"
( cd "$REPOTXT" && bash scripts/portable-tests.sh --scope=census ) 2>&1 | grep -E 'packages=' | sed 's/^/census: /'
echo "rc_census_pipeline=$?"
echo
echo "## 13. 台面清理（导出树在 /tmp，⛔ 入库；只建⛔ 删＝本件自建的临时树按名留着）"
echo "done=$(date '+%Y-%m-%d %H:%M:%S %z')"
