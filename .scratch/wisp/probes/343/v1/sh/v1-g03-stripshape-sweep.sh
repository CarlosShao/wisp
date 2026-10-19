#!/bin/sh
# 复核腿引用的剥注释哈希串（腿 report 写 97796acc…/265 行；派单写 b1610acf…/232 行）。
# 我自己不采信任何串 ⇒ 试若干把「形状明确」的尺，逐把报两侧行数＋两侧 sha256＋是否相同，
# 看腿那串能不能在任何一种合理形状下复现（＝是形状差，还是数字差）。
D=/d/tmp/wisp/20261019/343v1/ac2
LOG=/d/tmp/wisp/20261019/343v1/logs
A=$D/raw-start.go
B=$D/raw-head.go

run() { # $1=name  $2=command pipeline (as string evaluated by sh -c)
  ha=$(eval "$2" < "$A" | sha256sum | cut -d' ' -f1)
  la=$(eval "$2" < "$A" | wc -l)
  hb=$(eval "$2" < "$B" | sha256sum | cut -d' ' -f1)
  lb=$(eval "$2" < "$B" | wc -l)
  if [ "$ha" = "$hb" ]; then eq=SAME; else eq=DIFF; fi
  printf '%-6s lines=%s/%s shaA=%.8s shaB=%.8s full=%s EQ=%s\n' "$1" "$la" "$lb" "$ha" "$hb" "$ha" "$eq"
}

{
echo "== variant sweep: does any shape give 265 lines / 97796acc ? =="
run v1 "sed 's://.*::' | tr -d ' \t\r' | sed '/^\$//d'"
run v2 "sed 's://.*::' | tr -d ' \t\r'"
run v3 "sed 's://.*::' | tr -d '[:space:]'"
run v4 "sed '/^[[:space:]]*\/\//d' | sed 's://.*::' | tr -d ' \t\r' | sed '/^\$//d'"
run v5 "sed '/^[[:space:]]*\/\//d' | sed 's://.*::' | grep -v '^[[:space:]]*\$'"
run v6 "sed 's://.*::'"
run v7 "sed 's://.*::' | sed '/^[[:space:]]*\$/d'"
run v8 "grep -v '^[[:space:]]*//'"
run v9 "sed 's://.*::' | tr -s ' \t' ' ' | sed '/^ *\$/d'"
run v10 "sed 's:[[:space:]]*//.*::' | tr -d ' \t' | sed '/^\$/d'"
echo "== also: raw file line counts for reference =="
wc -l "$A" "$B"
} > $LOG/g03-stripshape-sweep.log 2>&1
echo DONE
