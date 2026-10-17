#!/bin/sh
# 331-v1 台面构造：两棵仓外导出树（T1=锚 72b65c61，T2=T1+本票落地件 errors.go@e09187e1）。
# 规矩块「台面配方」八件套 + 仓根 docs/ + 手工三件；树放 D: 卷（8.3 短名）；GOTMPDIR 指 D: 先 mkdir。
set -e
REPO="D:/work/workspace/projects plans/Wisp"
ROOT=/d/V331V1
T1=$ROOT/T1
T2=$ROOT/T2
cd "$REPO"

rm_existing() { :; }   # 临时件只建不删：这里只在目录为空时构造，绝不 clean 仓内任何东西

mkdir -p "$T1" "$T2" "$ROOT/tmp"

for ref_dir in "72b65c61:$T1" "e09187e1:$T2"; do
  ref=${ref_dir%%:*}; dir=${ref_dir#*:}
  git archive "$ref" go.mod go.sum internal cmd tools scripts frontend models docs \
    | tar -x -C "$dir"
done

# 手工三件（⛔ 跟踪，所以 archive 取不到）
for dir in "$T1" "$T2"; do
  mkdir -p "$dir/build" "$dir/third_party/sherpa-onnx" "$dir/frontend/dist"
  cp -f build/*.dll "$dir/build/" 2>/dev/null || true
  cp -f third_party/sherpa-onnx/*.dll "$dir/third_party/sherpa-onnx/" 2>/dev/null || true
  cp -rf frontend/dist/. "$dir/frontend/dist/" 2>/dev/null || true
done

# 取出后先数 .go 证非空
echo "T1_gocount=$(find $T1 -name '*.go' | wc -l)"
echo "T2_gocount=$(find $T2 -name '*.go' | wc -l)"
# 证两棵树的差异面 = 只那一枚文件
diff -rq "$T1" "$T2" 2>/dev/null | sed 's#^#DIFFQ #' | head -20
echo "build_ok"
