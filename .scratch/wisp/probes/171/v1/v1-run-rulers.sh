#!/bin/sh
# 171-v1（验收位台件）① 把改前／改后两把尺都从 commit 里抽到仓外，钉在同一枚锚点上各跑一遍。
# 用法：sh v1-run-rulers.sh <anchor-sha> [<anchor-sha> ...]
# 规矩：只读；尺的字节走 git show，不落仓内；读数落本目录 logs/。
set -eu
D=.scratch/wisp/probes/171/v1
T=$(mktemp -d)
mkdir -p "$D/logs"
OLD_SHA=19513cce   # 改前那版尺
NEW_SHA=3df74960   # 被验对象那版尺
git show "$OLD_SHA:.scratch/wisp/probes/154/gate-clauses.sh" > "$T/gate-old.sh"
git show "$NEW_SHA:.scratch/wisp/probes/154/gate-clauses.sh" > "$T/gate-new.sh"
echo "# 尺的字节出处：old=$OLD_SHA new=$NEW_SHA（均为 git show，非工作树）"
echo "# old sha1=$(git hash-object "$T/gate-old.sh") new sha1=$(git hash-object "$T/gate-new.sh")"
for A in "$@"; do
	git cat-file -t "$A" >/dev/null   # 复核是 commit（或可解引用对象）再用
	for v in old new; do
		sh "$T/gate-$v.sh" "$A" > "$D/logs/gate-$v-$A.txt" 2>&1 && rc=0 || rc=$?
		echo "== run ruler=$v anchor=$A rc=$rc out=$D/logs/gate-$v-$A.txt"
	done
	echo "== anchor=$A 声明/退码两行（old / new）"
	grep -E '^# 腿数＝|^# 聚合退码＝' "$D/logs/gate-old-$A.txt" | sed 's|^|  old |' || true
	grep -E '^# 腿数＝|^# 聚合退码＝' "$D/logs/gate-new-$A.txt" | sed 's|^|  new |' || true
	grep -E '^# (BAD|SHR)' "$D/logs/gate-new-$A.txt" | sed 's|^|  new-BAD/SHR |' || true
done
