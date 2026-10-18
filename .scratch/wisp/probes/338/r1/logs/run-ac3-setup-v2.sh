#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 台面 v2：把导出树再嵌深一层，让父目录名不合 8.3＝真能生别名
# （v1 的 8 枚红里有 1 枚 TestTicket252R1AlignmentAddsNoAuthorization 卡在
#   "D:\tmp\wisp\20261018\338r1 本身合 8.3＝无别名"＝台面落点形状，⛔ 判据响）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
ROOT=/d/tmp/wisp/20261018/338r1
BEFORE_SHA=0b1d1f0b
AFTER_SHA=87521cb4
# trees-at-v2-<时刻> ＝ 14 字符含连字符 ⇒ 不合 8.3 ⇒ 卷应生别名
NEST="$ROOT/trees-at-$(date +%H%M%S)"
mkdir -p "$NEST"
echo "NEST=$NEST"

for tag in before after; do
  case $tag in before) rev=$BEFORE_SHA ;; after) rev=$AFTER_SHA ;; esac
  T="$NEST/tree-$tag"
  mkdir -p "$T"
  git archive "$rev" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$T"
  git archive "$rev" docs | tar -x -C "$T"
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
  cp -r build/. "$T/build/"; cp -r third_party/sherpa-onnx/. "$T/third_party/sherpa-onnx/"
  [ -d frontend/dist ] && cp -r frontend/dist/. "$T/frontend/dist/"
  echo "tree-$tag go=$(find "$T" -name '*.go' | wc -l) dll=$(find "$T" -name '*.dll' | wc -l)"
done
echo '===== 别名核验：NEST 与 tree-* 两枚都该有短名 ====='
cmd //c "dir /x D:\\tmp\\wisp\\20261018\\338r1" 2>&1 | grep -iE "trees-at|blobch|wtest|TREE" | head -10
cmd //c "dir /x $(echo $NEST | sed 's|/d|D:|; s|/|\\|g')" 2>&1 | sed -n '5,14p'
echo "$NEST" > "$ROOT/nest-v2.txt"
