#!/bin/sh
# 342-v2 台面：自己 git archive 两棵仓外导出树（基线 fde1df02 ⇄ 落地 9581a48c），⛔ 复用别腿的树。
# 两枚终裁＝同一枚 winsec.go 路由的两支损伤 ⇒ 单差异台＝只差 internal/winsec 的落地字节。
# 六枚新钉两棵都带（它们⛔ 引用任何落地专属符号＝见 05 号尺），所以基线侧能直接量到损伤。
# ⛔ 删任何东西：本脚本只建；突变⛔ 打在这两棵 pristine 树上（另建 mut-* 拷贝）。
set -e
ROOT=/d/tmp/wisp/20261019/342v2
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO"
BASE=$(git rev-parse fde1df02)
LAND=$(git rev-parse 9581a48c)
echo "BASE_ANCHOR=$BASE"
echo "LAND_ANCHOR=$LAND"
echo "BASE_winsec_blob=$(git rev-parse $BASE:internal/winsec/winsec.go)"
echo "LAND_winsec_blob=$(git rev-parse $LAND:internal/winsec/winsec.go)"
PIECES='go.mod go.sum internal cmd tools scripts frontend models docs'
mkdir -p "$ROOT/treeB" "$ROOT/treeL"
git archive --format=tar "$BASE" $PIECES | tar -x -C "$ROOT/treeB"
git archive --format=tar "$LAND" $PIECES | tar -x -C "$ROOT/treeL"
# 两棵都补同一批新钉（逐枚从落地锚 blob 取，⛔ 从工作树取）
for F in ancestoropen_342_unix_test.go ancestorunlink_348_unix_test.go ancestorworld_342_348_unix_test.go \
         ancestoropen_342_windows_test.go ancestorunlink_348_windows_test.go ancestorworld_342_348_windows_test.go; do
  git cat-file blob "$LAND:internal/winsec/$F" > "$ROOT/treeB/internal/winsec/$F"
  git cat-file blob "$LAND:internal/winsec/$F" > "$ROOT/treeL/internal/winsec/$F"
done
# 手工补三件（都⛔ 跟踪）：build/*.dll、third_party/sherpa-onnx/*.dll、frontend/dist
for T in treeB treeL; do
  mkdir -p "$ROOT/$T/build" "$ROOT/$T/third_party/sherpa-onnx" "$ROOT/$T/frontend/dist"
  cp -f "$REPO"/build/*.dll "$ROOT/$T/build/" 2>/dev/null || true
  cp -f "$REPO"/third_party/sherpa-onnx/*.dll "$ROOT/$T/third_party/sherpa-onnx/" 2>/dev/null || true
done
# 取出后先数 .go 证非空（"0" ⛔ 是读数）
for T in treeB treeL; do
  echo "$T GO_FILES=$(find "$ROOT/$T" -name '*.go' | wc -l) WINSEC_GO=$(find "$ROOT/$T/internal/winsec" -name '*.go' | wc -l)"
done
echo '--- pristine 清单（internal/winsec 全枚 sha256，两棵并排，差异＝落地面）---'
( cd "$ROOT/treeB/internal/winsec" && sha256sum *.go | sort ) > "$ROOT/pristine-B.txt"
( cd "$ROOT/treeL/internal/winsec" && sha256sum *.go | sort ) > "$ROOT/pristine-L.txt"
diff "$ROOT/pristine-B.txt" "$ROOT/pristine-L.txt" > "$ROOT/pristine-diff.txt" || true
echo "DIFFERENT_OR_ONLY_IN_ONE=$(( $(grep -c '^[<>]' "$ROOT/pristine-diff.txt" || true) ))"
cat "$ROOT/pristine-diff.txt"
