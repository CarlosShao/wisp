#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 90
LOG=/d/tmp/wisp/20261019/343v1/logs
{
echo "== manual pieces presence (working tree, untracked) =="
for p in build third_party/sherpa-onnx frontend/dist docs scripts models internal cmd tools go.mod go.sum; do
  if [ -e "$p" ]; then
    n=$(find "$p" -type f 2>/dev/null | wc -l)
    printf 'EXISTS %s files=%s\n' "$p" "$n"
  else
    printf 'MISSING %s\n' "$p"
  fi
done
echo "== dll files =="
find build third_party/sherpa-onnx -maxdepth 2 -type f 2>/dev/null | sort
echo "== check-ignore =="
git check-ignore -v build third_party/sherpa-onnx frontend/dist || true
echo "== go/ast availability for AC#0 third gauge =="
/d/work/base/go/bin/go.exe version
echo "== 336 teeth test files present in HEAD? =="
git ls-files -- internal/risk | grep -c 'anchor_fallback_teeth_336' || true
git ls-files -- internal/risk | grep 'anchor_fallback_teeth_336' || true
echo "== risk package file roster in HEAD =="
git ls-tree --name-only HEAD -- internal/risk/ | wc -l
} > $LOG/g05-tree-prereq.log 2>&1
echo DONE
