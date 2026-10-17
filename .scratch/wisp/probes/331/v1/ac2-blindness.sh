#!/bin/sh
# 331-v1 补尺：AC#2 那把零改动尺的三枚失明对照（打在仓外副本，跟踪字节零改动）。
# B1 空转试（⛔ 任何改动都没有）· B2 整段删掉注释 · B3 在整行注释里塞一枚 U+2713（d22scan ban #8 认得的形）
# 判语要的那一维：这把尺的"绿"能不能被一枚没干活的腿直接拿走。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
EXP=/d/V331V1/exp2
mkdir -p "$EXP"
cd "$REPO"
strip() { grep -Ev '^[[:space:]]*//' "$@"; }

git show 72b65c61:internal/observe/errors.go > "$EXP/base.go"
git show e09187e1:internal/observe/errors.go  > "$EXP/after.go"

run() {   # run <label> <A> <B> <预期>
  strip < "$2" > "$EXP/s-A-$1.txt"; strip < "$3" > "$EXP/s-B-$1.txt"
  diff "$EXP/s-A-$1.txt" "$EXP/s-B-$1.txt" > "$EXP/d-$1.txt"; rc=$?
  echo "[$1] diff-rc=$rc (预期 $4)  非空diff=$([ -s "$EXP/d-$1.txt" ] && echo yes || echo no)"
}

echo "### B1 空转试：锚 blob ⇄ 锚 blob 自己（一枚字节都没改的腿）"
run B1 "$EXP/base.go" "$EXP/base.go" "0＝绿 ⇒ 一把什么都没干的腿也照样过 AC#2 的这把尺"

echo "### B2 整段删注释：把 MessageKey 那七行实话注释全删掉（产码字节仍零改动）"
AS=$(grep -n '^// MessageKey returns' "$EXP/after.go" | cut -d: -f1)
AE=$(grep -n '^func (c ErrorClass) MessageKey' "$EXP/after.go" | cut -d: -f1)
{ sed -n "1,$((AS-1))p" "$EXP/after.go"; sed -n "${AE},\$p" "$EXP/after.go"; } > "$EXP/mut-nocomment.go"
echo "  删掉的行数=$((AE-AS))（after=$AS..$((AE-1))），副本比 after.go 少 $(( $(wc -l < "$EXP/after.go") - $(wc -l < "$EXP/mut-nocomment.go") )) 行"
run B2 "$EXP/base.go" "$EXP/mut-nocomment.go" "0＝绿 ⇒ 注释被整段抹掉这把尺也不响"

echo "### B3 在整行注释里塞一枚 U+2713（d22scan ban #8 的射程内形，本尺的射程外）"
BS=$(grep -n '^// MessageKey returns' "$EXP/base.go" | cut -d: -f1)
{ sed -n "1,${BS}p" "$EXP/base.go"; printf '%s\n' "// verified checkmark: $(printf '%b' '\342\234\223')"; sed -n "$((BS+1)),\$p" "$EXP/base.go"; } > "$EXP/mut-emoji-comment.go"
grep -n 'checkmark' "$EXP/mut-emoji-comment.go" | sed 's#^#  落点: #'
run B3 "$EXP/base.go" "$EXP/mut-emoji-comment.go" "0＝绿 ⇒ 注释里的 emoji/幽灵引用这类形，这把尺一概看不见（归 d22scan，不归根本尺）"

echo "### B4 牙尺（正控）：同一算子在非注释行必须响，证 B1/B2/B3 的绿是射程不是恒真"
sed 's/return "error." + string(c)/return "error." + string(c) + ""/' "$EXP/base.go" > "$EXP/mut-code.go"
cmp -s "$EXP/base.go" "$EXP/mut-code.go" && echo "  B4 前置失败：没改上" || echo "  B4 前置：非注释行确实改了"
run B4 "$EXP/base.go" "$EXP/mut-code.go" "1＝红（该响）"

echo "### B5 跟踪字节还原核对"
git status --porcelain -- internal/observe/errors.go | sed 's#^#  porcelain(应为空): #'
git show HEAD:internal/observe/errors.go | sha256sum | sed 's#^#  HEAD-blob: #'
sha256sum "$EXP/after.go" | sed 's#^#  e09187e1-blob: #'
