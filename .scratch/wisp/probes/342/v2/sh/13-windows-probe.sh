#!/bin/sh
# 342-v2 必答(ii) 决定性探针（Windows 本机档）：同一枚只读探针件打进【基线树拷贝】⇄【落地树拷贝】，
# 量的是"稳态祖先形到底答在哪一支＋叶子复验到底答没答过"——⛔ 靠 racer 撞运气（确定性一发）。
# ★探针件⛔ 入库：只存在于 D:\tmp 的 mut-W* 拷贝树里；两棵树各拷一份、跑完留原样（⛔ 删）。
set -u
ROOTW='D:\tmp\wisp\20261019\342v2'
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"
PROBE="$ROOT/zzprobe/probe-342v2-windows_test.go.txt"
for TREE in treeB treeL; do
  MUT="mut-W-$TREE"
  [ -d "$ROOT/$MUT" ] || cp -a "$ROOT/$TREE" "$ROOT/$MUT"
  cp -f "$PROBE" "$ROOT/$MUT/internal/winsec/probe_342v2_windows_test.go"
  LABEL="ii02-${MUT}-probe"
  DLLP=$(cygpath -w "$ROOTW\\$MUT\\build"); SHP=$(cygpath -w "$ROOTW\\$MUT\\third_party\\sherpa-onnx")
  ( cd "$ROOTW\\$MUT" && PATH="$SHP;$DLLP;$PATH" GOCACHE='D:\tmp\wisp\20261019\342v2\gocache-win' \
      go test ./internal/winsec/ -v -count=1 -timeout 20m -run 'TestP342V2' ) > "$LOGS/$LABEL.raw.txt" 2>&1
  RC=$?
  RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
  PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
  FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
  SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
  echo "$LABEL RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" | tee -a "$LOGS/four-counts-probe.tsv"
  grep -E 'W1a |W1b |W2 |W3 |^\s+--- (FAIL|SKIP|PASS)' "$LOGS/$LABEL.raw.txt" | head -30
  echo "-----"
done
