#!/bin/sh
# 331-v1 AC#2 独立复跑（非实现者腿，自己重放，⛔ 背实现腿的读数）。
# 读数一律取自具名提交树的 blob（⛔ 读共享工作树在飞字节）。
# 突变只打在仓外副本 /d/V331V1/exp/**，跟踪字节全程零改动；临时件只建不删。
# 尺 = 滤掉整行注释（^[[:space:]]*//）后比 sha256 + diff（与实现腿 zero-change-check.sh 同一算子，本腿自己重写一遍）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
EXP=/d/V331V1/exp
mkdir -p "$EXP"
cd "$REPO"

strip() { grep -Ev '^[[:space:]]*//' "$@"; }

ruler() {                       # ruler <label> <pathA> <pathB>
  nA=$(basename "$2" .go); nB=$(basename "$3" .go)
  A="$EXP/strip-$nA.txt"; B="$EXP/strip-$nB.txt"; D="$EXP/diff-$nA-$nB.txt"
  strip < "$2" > "$A"; strip < "$3" > "$B"
  SA=$(sha256sum < "$A" | cut -d' ' -f1)
  SB=$(sha256sum < "$B" | cut -d' ' -f1)
  diff "$A" "$B" > "$D"; rc=$?
  [ -s "$D" ] && ne=yes || ne=no
  echo "[$1] shaA=$SA"
  echo "[$1] shaB=$SB"
  echo "[$1] same=$([ "$SA" = "$SB" ] && echo yes || echo no) diff-nonempty=$ne diff-rc=$rc"
  return $rc
}

echo "### A0 两侧 blob 全字节身份（未滤注释，作还原基线）"
git show 72b65c61:internal/observe/errors.go > "$EXP/base.go"
git show e09187e1:internal/observe/errors.go  > "$EXP/after.go"
sha256sum "$EXP/base.go" "$EXP/after.go" | sed 's#^#  #'

echo "### A1 主尺：锚 blob 72b65c61 ⇄ 落地 blob e09187e1，滤整行注释后（票面要 sha 同值 + diff 空 + rc=0）"
ruler A1 "$EXP/base.go" "$EXP/after.go"; echo "  A1_rc=$?"

echo "### A2 正控（必做，否则 A1 的空不算证据）：只在非注释行改一枚字面量"
echo "  落点尺＝内容锚（函数体内那枚 token），⛔ 绝对行号"
sed 's/return "error." + string(c)/return "errorX." + string(c)/' "$EXP/after.go" > "$EXP/mut-noncomment.go"
grep -n 'return "errorX' "$EXP/mut-noncomment.go" | sed 's#^#  落点: #'
ruler A2 "$EXP/base.go" "$EXP/mut-noncomment.go"; echo "  A2_rc=$? (期望 1)"

echo "### A3 必答实验·反向写：把落地件的注释块整块换回改前那三行谎（UI and logs embed the key）"
# 反向写＝取 after.go 的非注释骨架 ＋ base.go 的那三行注释块，拼成一枚"说假话但产码字节全同"的副本。
# 块边界用内容锚（首行 // MessageKey returns…、尾行 func (c ErrorClass) MessageKey）现算，⛔ 写死行号。
BS=$(grep -n '^// MessageKey returns'      "$EXP/base.go"  | cut -d: -f1)
BE=$(grep -n '^func (c ErrorClass) MessageKey' "$EXP/base.go"  | cut -d: -f1)
AS=$(grep -n '^// MessageKey returns'      "$EXP/after.go" | cut -d: -f1)
AE=$(grep -n '^func (c ErrorClass) MessageKey' "$EXP/after.go" | cut -d: -f1)
echo "  block-lines: base=$BS..$((BE-1))  after=$AS..$((AE-1))"
{ sed -n "1,$((AS-1))p" "$EXP/after.go"; sed -n "${BS},$((BE-1))p" "$EXP/base.go"; sed -n "${AE},\$p" "$EXP/after.go"; } > "$EXP/mut-lie.go"
echo "  反向写副本里那段注释（证谎确实被写回去了）："
sed -n "${AS},$((AE-1))p" "$EXP/mut-lie.go" | sed 's#^#    #'
echo "  反向写副本 vs after.go 的差（应只有注释块那几行）："
diff "$EXP/after.go" "$EXP/mut-lie.go" | sed 's#^#    #'
ruler A3 "$EXP/base.go" "$EXP/mut-lie.go"; echo "  A3_rc=$? (尺对注释内容失明 ⇒ 仍 0)"
echo "### A3b 决定性对偶：实话件 after.go ⇄ 说谎件 mut-lie.go，注释内容相反，尺却回同一串"
ruler A3b "$EXP/after.go" "$EXP/mut-lie.go"; echo "  A3b_rc=$? (期望 0＝这把尺分不开'实话'与'谎')"
cmp -s "$EXP/after.go" "$EXP/mut-lie.go" && echo "  A3b 前置校验失败：两副本全字节相同，实验没做成" || echo "  A3b 前置校验：两副本全字节【不等】（注释内容确实相反），尺却同串 ⇒ 失明成立"
sha256sum "$EXP/mut-lie.go" | sed 's#^#  A3_mutation-copy-hash(反向写副本全字节): #'
echo "  还原核对：after.go 全字节哈希（⛔ 突变没留在任何跟踪件上）"
sha256sum "$EXP/after.go" | sed 's#^#  A3-restore: #'

echo "### A4 失明方向的对偶：把整行注释再改一次（改向另一句真话），尺同样看不见"
sed 's#^// that deferred layer.*#// comment-only rewrite, alternate true sentence#' "$EXP/after.go" > "$EXP/mut-recomment.go"
cmp -s "$EXP/after.go" "$EXP/mut-recomment.go" && echo "  A4 全字节 cmp：相等（不该相等，说明 sed 没落）" || echo "  A4 全字节 cmp：不等（证注释确实被改过）"
ruler A4 "$EXP/after.go" "$EXP/mut-recomment.go"; echo "  A4_rc=$? (期望 0＝尺对整行注释改动不敏感)"

echo "### A5 差尺·逐行核改动枚数（票面 AC#2 末句：非整行注释改动枚数必为 0）"
git diff 72b65c61..e09187e1 -- internal/observe/errors.go \
  | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' > "$EXP/changed-lines.txt"
TOT=$(wc -l < "$EXP/changed-lines.txt")
CMT=$(grep -cE '^[+-][[:space:]]*//' "$EXP/changed-lines.txt" || true)
NONCMT=$(grep -vcE '^[+-][[:space:]]*//' "$EXP/changed-lines.txt" || true)
echo "  changed_total=$TOT  whole-comment=$CMT  non-whole-comment=$NONCMT (票面要求 0)"

echo "### A6 还原核对（跟踪字节全程零改动，⛔ 任何突变没落进仓）"
git diff --stat -- internal/observe/errors.go | sed 's#^#  worktree-vs-index: #'
git status --porcelain -- internal/observe/errors.go | sed 's#^#  porcelain: #'
sha256sum "$EXP/base.go" "$EXP/after.go" | sed 's#^#  restore-baseline: #'
