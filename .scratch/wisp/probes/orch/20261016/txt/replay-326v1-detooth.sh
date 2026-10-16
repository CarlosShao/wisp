#!/usr/bin/env bash
# 编排者自己重跑 326-v1 的主凭据：摘掉 dao_session.go:300 那行 touch ⇒ 两枚常驻判据必须红
# ⛔ 只在仓外导出树里做（共享工作树一字不动）
set -u
ROOT=/d/tmp/orch-326v1-replay
LOG=/d/work/workspace/"projects plans"/Wisp/.scratch/wisp/probes/orch/20261016/logs
cd "/d/work/workspace/projects plans/Wisp" || exit 9
rm -rf "$ROOT" 2>/dev/null
mkdir -p "$ROOT" || exit 9
git archive HEAD | tar -x -C "$ROOT" || { echo "ARCHIVE_FAIL rc=$?"; exit 8; }
echo "TREE=$ROOT"
echo "go_files=$(find "$ROOT" -name '*.go' | wc -l)"
cd "$ROOT" || exit 7
echo "--- 起手 sha256（摘之前）---"
sha256sum internal/memory/dao_session.go | tee "$LOG/orch-326v1-replay-sha-before.txt"
sed -n '300p' internal/memory/dao_session.go
cp internal/memory/dao_session.go "$ROOT/dao_session.go.orig"
sed -i '300d' internal/memory/dao_session.go
echo "--- 摘之后 299-302 ---"
sed -n '299,302p' internal/memory/dao_session.go
echo "--- 去牙那一发（两枚常驻判据）---"
go test ./internal/memory -run 'TestAppendSessionTurn326v1' -count=1 -timeout 30m -v > "$LOG/orch-326v1-de-tooth.txt" 2>&1
echo "rc=$?" | tee -a "$LOG/orch-326v1-de-tooth.txt"
grep -E '^--- (PASS|FAIL)|^(ok|FAIL)|326v1' "$LOG/orch-326v1-de-tooth.txt" | head -25
echo "--- 还原（仓外树自带原件，验逐字同形）---"
cp "$ROOT/dao_session.go.orig" internal/memory/dao_session.go
sha256sum internal/memory/dao_session.go | tee "$LOG/orch-326v1-replay-sha-after.txt"
echo "--- 还原后同一发（应绿）---"
go test ./internal/memory -run 'TestAppendSessionTurn326v1' -count=1 -timeout 30m > "$LOG/orch-326v1-restored.txt" 2>&1
echo "rc=$?" | tee -a "$LOG/orch-326v1-restored.txt"
tail -3 "$LOG/orch-326v1-restored.txt"
