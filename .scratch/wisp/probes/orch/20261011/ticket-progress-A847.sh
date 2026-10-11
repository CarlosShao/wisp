#!/bin/sh
# A847 票面进度尺（编排者现跑）
# 口径三条：① 名册＝git ls-tree HEAD（⛔ ls-files 读的是索引＝此刻两枚腿可能刚 stage 东西）
#          ② 框普查在 HEAD blob 层量（⛔ 工作树：两枚验收腿正在往票尾写 bullet）
#          ③ 框尺＝grep -cE '^[[:space:]]*- \[( |x)\]'（10-02 起含缩进项；跨换尺点⛔ 与旧数相减）
cd "$(git rev-parse --show-toplevel)" || exit 9
DIR=.scratch/wisp/issues
OUT=.scratch/wisp/probes/orch/20261011/logs/orch-A847-ticket-progress.txt
: > "$OUT.tmp"

git ls-tree -r HEAD --name-only -- "$DIR" > "$OUT.tmp.names"

total_all=$(grep -c '' "$OUT.tmp.names")
readme=$(grep -c 'README' "$OUT.tmp.names")
done_n=$(grep -c -- '-done\.md$' "$OUT.tmp.names")
open_n=$((total_all - readme - done_n))

checked_all=0
unchecked_all=0
checked_open=0
unchecked_open=0
open_zero_unchecked=0
done_still_unchecked=0

while IFS= read -r p; do
  case "$p" in
    *README*) continue ;;
  esac
  body=$(git show "HEAD:$p")
  c=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[x\]')
  u=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[ \]')
  checked_all=$((checked_all + c))
  unchecked_all=$((unchecked_all + u))
  case "$p" in
    *-done.md)
      if [ "$u" -gt 0 ]; then done_still_unchecked=$((done_still_unchecked + 1)); fi
      ;;
    *)
      checked_open=$((checked_open + c))
      unchecked_open=$((unchecked_open + u))
      if [ "$u" -eq 0 ]; then open_zero_unchecked=$((open_zero_unchecked + 1)); fi
      ;;
  esac
done < "$OUT.tmp.names"

all=$((checked_all + unchecked_all))
allopen=$((checked_open + unchecked_open))

{
  echo "date=$(date '+%Y-%m-%d %H:%M:%S %z')  HEAD=$(git rev-parse --short HEAD)  尺＝HEAD blob 层"
  echo "== 分母 =="
  echo "ticket_files_including_README=$total_all"
  echo "README=1  =>  tickets=$((total_all - 1))"
  echo "done=$done_n   open=$open_n"
  echo "== 框（判据格） =="
  echo "ALL tickets      : checked=$checked_all unchecked=$unchecked_all total=$all"
  echo "OPEN tickets only: checked=$checked_open unchecked=$unchecked_open total=$allopen"
  echo "== 百分比（保留一位小数；手算，标明算式） =="
  awk -v d="$done_n" -v o="$open_n" -v n=$((total_all - 1)) -v ca="$checked_all" -v ta="$all" \
      -v co="$checked_open" -v to="$allopen" -v uz="$open_zero_unchecked" -v dz="$done_still_unchecked" 'BEGIN{
      printf "P1 结案票占比 = %d/%d = %.1f%%\n", d, n, 100*d/n;
      printf "P2 格已勾占比（全库） = %d/%d = %.1f%%\n", ca, ta, 100*ca/ta;
      printf "P3 格已勾占比（只算开放票） = %d/%d = %.1f%%\n", co, to, 100*co/to;
      printf "P4 开放票里零未勾格（＝能收口） = %d/%d = %.1f%%\n", uz, o, 100*uz/o;
      printf "残  已 -done 却仍有未勾框（乙类） = %d/%d = %.1f%%\n", dz, d, 100*dz/d;
    }'
  echo "== 甲类名册（开放票、未勾格=0 => 只差收口动作） =="
  : > "$OUT.tmp.jia"
  : > "$OUT.tmp.yi"
} >> "$OUT" 2>&1

while IFS= read -r p; do
  case "$p" in *README*) continue ;; esac
  body=$(git show "HEAD:$p")
  c=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[x\]')
  u=$(printf '%s\n' "$body" | grep -cE '^[[:space:]]*- \[ \]')
  case "$p" in
    *-done.md) [ "$u" -gt 0 ] && printf '  %s  checked=%s unchecked=%s\n' "$(basename "$p")" "$c" "$u" >> "$OUT.tmp.yi" ;;
    *)         [ "$u" -eq 0 ] && printf '  %s  checked=%s\n' "$(basename "$p")" "$c" >> "$OUT.tmp.jia" ;;
  esac
done < "$OUT.tmp.names"

{ echo "-- 甲类（开放票、未勾格=0 ⇒ 只差收口动作）--"; cat "$OUT.tmp.jia";
  echo "-- 乙类（已 -done 却仍有未勾框）--"; cat "$OUT.tmp.yi"; } >> "$OUT" 2>&1
cat "$OUT"
