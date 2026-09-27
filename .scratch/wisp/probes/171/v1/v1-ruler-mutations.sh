#!/bin/sh
# 171-v1 台件④：三发「尺自己被动过」的变异 + 两枚现量计数。
#  全部尺的字节走 git show 到仓外 $TMPDIR，**probes/154/** 与 probes/171/r1/** 一字节不碰**。
#   S1 样本 C（丢句柄＋同文件关别人的 id）：改前的尺该静默、改后的尺必须点名（AC#6③ 的那一发）
#   S2 SHR 那一味藏不藏东西：把 G7 主尺的射程悄悄调小（多两条 :!）⇒ 读数 3→1，退码动不动？
#   S3 空心那一发：摘掉 G5 的 want_n ⇒ pair() 必须拒绝出结论（预期 rc=4，且不出「干净」）
set -eu
D=.scratch/wisp/probes/171/v1
A=3df749607a0da1e70579f5c3abb71003350a45d2
T=$(mktemp -d)
mkdir -p "$D/logs" "$D/logs"
S=$(cd "$D/logs" && pwd)
git show 19513cce:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-old.sh"
git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-new.sh"

echo "###### S0 现量两枚计数（票面 §4 说 42、尺内注释说 40）"
for a in 19513cce 3df74960; do
	n=$(git grep -lEw 'Close\(\)' "$a" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' | wc -l)
	echo "   Close() 命中的非测试生产码文件数 @$a ＝ $n"
done
echo "   flip-declaration.sh 里与「腿数」有关的行："
grep -n '腿数\|legs' .scratch/wisp/probes/161/r6/flip-declaration.sh | sed 's|^|     |' || echo "     （无）"

echo
echo "###### S1 样本 C：丢弃树 base + c_discard.go（_ =OpenScope 丢句柄 ＋ 同文件 CloseScope(别人的 id)）"
rm -rf "$T/base" 2>/dev/null || true
mkdir -p "$T/base"
git archive "$A" internal cmd | tar -x -C "$T/base"
mkdir -p "$T/base/internal/fake171v1c"
printf '%s\n' \
	'package fake171v1c' \
	'' \
	'import "wisp/internal/risk"' \
	'' \
	'func leaky(p *risk.Provenance, l *risk.Provenance, id string) {' \
	'	_ = p.OpenScope(id)' \
	'	l.CloseScope("task-other")' \
	'}' > "$T/base/internal/fake171v1c/c_discard.go"
gofmt -w "$T/base/internal/fake171v1c/c_discard.go"
( cd "$T/base" && git init -q && git add internal cmd && git -c user.name=v1 -c user.email=v1@invalid commit -q -m S1 )
for v in old new; do
	rc=0
	( cd "$T/base" && export GIT_DIR="$T/base/.git" GIT_WORK_TREE="$T/base" && sh "$T/gate-$v.sh" HEAD > "$S/s1-cdiscard-$v.txt" 2>&1 ) || rc=$?
	echo "   ruler=$v rc=$rc 点名 fake171v1c 的腿=[$(awk '/^## /{k=$2} /^#   UNPAIRED internal\/fake171v1c/{print k}' "$S/s1-cdiscard-$v.txt" | sort -u | tr '\n' ' ')]"
done
echo "   新尺里那一行的批注："
grep -E 'DISCARD-HANDLE|note 丢句柄' "$S/s1-cdiscard-new.txt" | grep -i fake171 | sed 's|^|     |' || true
echo "   改前的尺对同一枚射程的 G5 名册："
awk '/^## G5 OpenScope/{f=1} /^## G5 正控/{f=0} f&&/^#   UNPAIRED|^# 未成对枚数/' "$S/s1-cdiscard-old.txt" | sed 's|^|     |'

echo
echo "###### S2 SHR 那一味：把 G7 主尺射程调小（追加两条 :!，want_n 仍是 3）⇒ 读数 3→1，退码动不动"
awk -v add="':!internal/observe/*' ':!internal/proc/*'" \
	'index($0,":!internal/plugin/disposal.go") && index($0,"GO2"){print $0 " " add; next} {print}' \
	"$T/gate-new.sh" > "$T/gate-shrunk.sh"
cmp -s "$T/gate-shrunk.sh" "$T/gate-new.sh" && { echo "    改射程没命中，射程没改成——这发不成立"; exit 9; }
diff "$T/gate-new.sh" "$T/gate-shrunk.sh" | sed 's|^|   diff |'
sh "$T/gate-shrunk.sh" "$A" > "$S/s2-shrunk-range.txt" 2>&1 && rc=0 || rc=$?
echo "   rc=$rc（基线：新尺射程未动时 rc=0）"
grep -E '^# (BAD|SHR)|^# 腿数＝|^# 聚合退码＝' "$S/s2-shrunk-range.txt" | sed 's|^|     |'

echo
echo "###### S3 空心那一发：摘掉 G5 的 want_n ⇒ 尺必须拒绝出结论"
awk 'BEGIN{n=0} /^want_n 1$/ && n==0 {n=1; next} {print}' "$T/gate-new.sh" > "$T/gate-non.txt"
echo "   摘掉的行数＝$(diff "$T/gate-new.sh" "$T/gate-non.txt" | grep -c '^<') 枚 want_n"
sh "$T/gate-non.txt" "$A" > "$S/s3-no-want-n.txt" 2>&1 && rc=0 || rc=$?
echo "   rc=$rc（预期 4＝pair 腿没登记基线直接死）"
echo "   末三行："
tail -3 "$S/s3-no-want-n.txt" | sed 's|^|     |'
echo "   它有没有给出「干净」结论（腿数／退码行）："
grep -E '^# 腿数＝|^# 聚合退码＝' "$S/s3-no-want-n.txt" | sed 's|^|     |' || echo "     （没有出结论——死在登记处）"
