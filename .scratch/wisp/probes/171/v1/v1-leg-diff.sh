#!/bin/sh
# 171-v1 台件②：逐腿差集——同一枚锚点上把改前/改后两把尺的「命令行（含 pathspec 射程）」与
# 「UNPAIRED 名册」两样东西分别做 comm 两向。这是「一条 pathspec 都没减」的机器证，不是读它的话。
# 用法：sh v1-leg-diff.sh <anchor-sha>
set -eu
A="$1"
D=.scratch/wisp/probes/171/v1
mkdir -p "$D/legs"
ex() { # $1 gate 输出文件 $2 抽什么（cmd|unpaired）
	awk -v mode="$2" '
		/^## /{ leg=substr($0,4); gsub(/[ \t]+$/,"",leg); split(leg,p," "); key=p[1]; next }
		/^\$ git grep/ { if (mode=="cmd") print key"\t"$0; next }
		/^#   UNPAIRED/ { if (mode=="unpaired") print key"\t"$2":"$3; next }
	' "$1" | sort
}
for v in old new; do
	f="$D/logs/gate-$v-$A.txt"
	ex "$f" cmd      > "$D/legs/$v-$A-cmd.tsv"
	ex "$f" unpaired > "$D/legs/$v-$A-unpaired.tsv"
done
echo "###### ① 命令行（射程本体）逐腿对照：old vs new（同一锚 $A）"
for leg in G5 G5-正控 G5-负一负 G6 G6-正控 G6-负一负 G7 G7-正控 G7-负一负; do
	o=$(grep -E "^$leg	" "$D/legs/old-$A-cmd.tsv" | sed "s|^$leg	||" | tr '\n' '|')
	n=$(grep -E "^$leg	" "$D/legs/new-$A-cmd.tsv" | sed "s|^$leg	||" | tr '\n' '|')
	# 只比 pathspec 段（-- 之后），pattern 段本来就是本次改动处
	op=$(printf '%s' "$o" | sed 's|.* -- \(.*\)|\1|')
	np=$(printf '%s' "$n" | sed 's|.* -- \(.*\)|\1|')
	if [ "$op" = "$np" ]; then s='SAME-pathspec'; else s='DIFF-pathspec'; fi
	echo "-- $leg [$s]"
	echo "   old-cmd: $o"
	echo "   new-cmd: $n"
	[ "$s" = 'DIFF-pathspec' ] && { echo "   old-path: $op"; echo "   new-path: $np"; }
	true
done
echo
echo "###### ② 逐腿 UNPAIRED 名册差集（comm -23＝退场／comm -13＝新增）"
for leg in G5 G5-正控 G5-负一负 G6 G6-正控 G6-负一负 G7 G7-正控 G7-负一负; do
	grep -E "^$leg	" "$D/legs/old-$A-unpaired.tsv" | cut -f2 | sort -u > "$D/legs/t-old.txt"
	grep -E "^$leg	" "$D/legs/new-$A-unpaired.tsv" | cut -f2 | sort -u > "$D/legs/t-new.txt"
	gone=$(comm -23 "$D/legs/t-old.txt" "$D/legs/t-new.txt" | tr '\n' ' ')
	add=$(comm -13 "$D/legs/t-old.txt" "$D/legs/t-new.txt" | tr '\n' ' ')
	printf -- '-- %s 改前%s枚 改后%s枚 退场=[%s] 新增=[%s]\n' "$leg" \
		"$(wc -l < "$D/legs/t-old.txt" | tr -d ' ')" "$(wc -l < "$D/legs/t-new.txt" | tr -d ' ')" "$gone" "$add"
done
echo
echo "###### ③ 全门合计（九腿并集）退场/新增"
cut -f2 "$D/legs/old-$A-unpaired.tsv" | sort -u > "$D/legs/all-old.txt"
cut -f2 "$D/legs/new-$A-unpaired.tsv" | sort -u > "$D/legs/all-new.txt"
echo "退场(只在改前)："; comm -23 "$D/legs/all-old.txt" "$D/legs/all-new.txt" || true
echo "新增(只在改后)："; comm -13 "$D/legs/all-old.txt" "$D/legs/all-new.txt" || true
echo "改前枚数合计=$(wc -l < "$D/legs/old-$A-unpaired.tsv") 改后枚数合计=$(wc -l < "$D/legs/new-$A-unpaired.tsv")"
