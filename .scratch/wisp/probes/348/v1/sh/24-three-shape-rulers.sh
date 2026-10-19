#!/bin/sh
# 342-v2 三枚剥注释／结构尺串（★我自己算的，⛔ 抄腿报；两把锚＝510a1519..9581a48c）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B51=510a1519; LAND=9581a48c
echo '### 尺 1＝winsec.go 那笔可执行行变化枚数（`git diff -U0` 里 `^[+-]` 逐枚判是否以 // 开头或为注释续行）'
LINESF=/d/tmp/wisp/20261019/342v2/w342v2-lines.txt
git diff -U0 "$B51" "$LAND" -- internal/winsec/winsec.go | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' > $LINESF
TOT=$(wc -l < $LINESF)
COMMENT=$(grep -cE '^[+-][[:space:]]*(//|\*)' $LINESF || true)
echo "CHANGED_LINES=$TOT COMMENT_LINES=$COMMENT EXECUTABLE_LINES=$((TOT-COMMENT))"
echo '--- 可执行行逐字（应只＝一次 open 换形＋两处删除换形＋父钉那三行）---'
grep -vE '^[+-][[:space:]]*(//|\*)' $LINESF
echo
echo '### 尺 2＝新增顶层导出标识符枚数（^func [A-Z] ⇄ ^var [A-Z] ⇄ ^type [A-Z] ⇄ ^const [A-Z]，只看落地 9 枚新件）'
NEWGO=$(git diff --name-status "$B51" "$LAND" -- internal/winsec/ | awk '$1=="A"{print $2}')
echo "$NEWGO" | wc -l | sed 's/^/NEW_GO_FILES=/'
EXP=0
for f in $NEWGO; do
  n=$(git cat-file blob "$LAND:$f" | grep -cE '^(func [A-Z]|var [A-Z]|type [A-Z]|const [A-Z])' || true)
  m=$(git cat-file blob "$LAND:$f" | grep -oE '^(func|var|type|const) [A-Z][A-Za-z0-9_]*' | tr '\n' ' ')
  [ -n "$m" ] && echo "  $f -> $m"
  EXP=$((EXP+n))
done
echo "EXPORTED_TOP_LEVEL_IDENTIFIERS=$EXP"
echo '--- 分栏：产码件（⛔ _test.go）的导出顶层标识符＝本票"⛔ 新导出名"的自证 ---'
PROD_EXP=0
for f in $NEWGO; do
  case "$f" in *_test.go) continue;; esac
  n=$(git cat-file blob "$LAND:$f" | grep -cE '^(func [A-Z]|var [A-Z]|type [A-Z]|const [A-Z])' || true)
  echo "  $f -> $n"
  PROD_EXP=$((PROD_EXP+n))
done
n=$(git cat-file blob "$LAND:internal/winsec/winsec.go" | grep -E '^(func [A-Z]|var [A-Z]|type [A-Z]|const [A-Z])' | wc -l)
o=$(git cat-file blob "$B51:internal/winsec/winsec.go" | grep -E '^(func [A-Z]|var [A-Z]|type [A-Z]|const [A-Z])' | wc -l)
echo "  winsec.go: landed=$n baseline=$o  (symdiff 逐名)"
diff <(git cat-file blob "$B51:internal/winsec/winsec.go" | grep -oE '^(func|var|type|const) [A-Z][A-Za-z0-9_]*' | sort) <(git cat-file blob "$LAND:internal/winsec/winsec.go" | grep -oE '^(func|var|type|const) [A-Z][A-Za-z0-9_]*' | sort) && echo "  winsec.go EXPORTED SYMBOL SET IDENTICAL"
echo "PROD_CODE_EXPORTED_TOP_LEVEL_IDENTIFIERS=$PROD_EXP  (须=0)"
echo
echo '### 尺 3＝每枚 .go 的行数＋sha256 前 16（落地树 blob ⇄ 我的仓外树同串自证）'
for f in $(git diff --name-only "$B51" "$LAND" -- internal/winsec/); do
  echo "  $(git cat-file blob "$LAND:$f" | wc -l) lines  $(git cat-file blob "$LAND:$f" | sha256sum | cut -c1-16)  $(basename $f)"
done
echo '--- 我的 treeL 物理字节逐枚对回（须与上面同串）---'
for f in $(git diff --name-only "$B51" "$LAND" -- internal/winsec/ | xargs -n1 basename); do
  echo "  $(wc -l < "/d/tmp/wisp/20261019/342v2/treeL/internal/winsec/$f") lines  $(sha256sum "/d/tmp/wisp/20261019/342v2/treeL/internal/winsec/$f" | cut -c1-16)  $f"
done
