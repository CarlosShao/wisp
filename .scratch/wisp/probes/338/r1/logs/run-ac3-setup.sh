#!/bin/sh
# 票 338 落地腿 338-r1 · AC#3 台面搭建（仓外导出树，改前⇄改后两枚，同一配方）
# ★落点＝D:\tmp\wisp\20261018\338r1\tree-*（⛔ 盘符根、⛔ C 卷 Temp）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
ROOT=/d/tmp/wisp/20261018/338r1
BEFORE_SHA=0b1d1f0b     # 我 AC#0 那笔＝AC#1 归位之前
AFTER_SHA=87521cb4      # 我 AC#1 归位之后
mkdir -p "$ROOT"

for tag in before after; do
  case $tag in
    before) rev=$BEFORE_SHA ;;
    after)  rev=$AFTER_SHA ;;
  esac
  # ★临时件只建不删：每次跑用唯一目录名（时刻后缀），不复用、不清空既有树
  T="$ROOT/tree-$tag-$(date +%H%M%S)-$tag"
  mkdir -p "$T"
  echo "$tag=$T" >> "$ROOT/tree-paths.txt"
  echo "===== tree-$tag  rev=$rev ====="
  git archive "$rev" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$T"
  git archive "$rev" docs | tar -x -C "$T"
  # 手工补三件（⛔ 跟踪件，不在 archive 里）
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend"
  cp -r build/. "$T/build/" 2>/dev/null; echo "cp_build_rc=$?"
  cp -r third_party/sherpa-onnx/. "$T/third_party/sherpa-onnx/" 2>/dev/null; echo "cp_sherpa_rc=$?"
  if [ -d frontend/dist ]; then mkdir -p "$T/frontend/dist"; cp -r frontend/dist/. "$T/frontend/dist/" 2>/dev/null; echo "cp_dist_rc=$?"; else echo "cp_dist=NO-SUCH-DIR"; fi
  echo "go_files=$(find "$T" -name '*.go' | wc -l)"
  echo "dll_count=$(find "$T" -name '*.dll' | wc -l)"
done

echo '===== 两枚树的差异（应只有 tools 那枚） ====='
TB=$(grep '^before=' "$ROOT/tree-paths.txt" | tail -1 | cut -d= -f2)
TA=$(grep '^after=' "$ROOT/tree-paths.txt" | tail -1 | cut -d= -f2)
diff -rq "$TB" "$TA" 2>/dev/null | head -20
echo '===== 短名别名核验（D 卷 8.3） ====='
cmd //c "dir /x D:\\tmp\\wisp\\20261018\\338r1" 2>&1 | sed -n '1,20p'
