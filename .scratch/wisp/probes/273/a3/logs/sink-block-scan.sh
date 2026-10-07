#!/bin/sh
# 273-a3 尺 A1（v2）：括号平衡扫描。对每一枚匹配行，从该行起累行直到圆括号闭合，
# 再只在**闭合的那一块**里判 `Sink:` 是否出现。
# 为什么要平衡：起手锚 §0.2 已证唯一真传 Sink 的那枚是跨 3 行字面量
# （internal/models/bridge_test.go:16 起 New(statemachine.Options{ ，:18 才出现 Sink:），
# 只看匹配行＝对正控失明。⛔ 本尺不数裸符号名，只认构造调用形状。
# usage: sink-block-scan.sh <awk-regex（用 [.] [(] 免转义）> <file...>
PAT="$1"; shift
awk -v pat="$PAT" '
function emit(f, l, b,   has, k) {
  has = (b ~ /Sink[[:space:]]*:/) ? "SINK=YES " : "SINK=no  "
  k = gsub(/\n/, "\n", b)
  printf "%s %s:%d  block_lines=%d\n", has, f, l, k
}
FNR==1 { active=0; depth=0; buf="" }
{
  if (!active) {
    if ($0 ~ pat) {
      active=1; startline=FNR; buf=$0 "\n"
      s=$0; depth += gsub(/\(/, "(", s) - gsub(/\)/, ")", s)
      if (depth <= 0) { emit(FILENAME, startline, buf); active=0; depth=0; buf="" }
    }
  } else {
    buf = buf $0 "\n"
    s=$0; depth += gsub(/\(/, "(", s) - gsub(/\)/, ")", s)
    if (depth <= 0) { emit(FILENAME, startline, buf); active=0; depth=0; buf="" }
  }
}
END { if (active) emit(FILENAME, startline, buf "<<UNCLOSED>>") }
' "$@"
