#!/bin/sh
# 171-v1 台件⑤：声明机制还连着退码吗（三发）＋「减法侧有没有牙」那一发。
# 一律在仓外 $TMPDIR 的尺副本上翻，**probes/154/gate-clauses.sh 一字节不改**。
set -eu
D=.scratch/wisp/probes/171/v1
A=3df749607a0da1e70579f5c3abb71003350a45d2
T=$(mktemp -d)
mkdir -p "$D/logs"
S=$(cd "$D/logs" && pwd)
git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate.sh"
cp "$T/gate.sh" "$T/f0.sh"                       # 未翻＝基线
grep -E "^want (G7pos|G6neg) " "$T/gate.sh" | sed "s|^|# 待翻的声明行: |"
awk 'NR==FNR{next}1' /dev/null /dev/null 2>/dev/null || true

flip() { # $1 输出名 $2.. 「LEG|旧值|新值」
	out="$1"; shift
	cp "$T/gate.sh" "$T/$out.sh"
	for spec in "$@"; do
		leg=$(printf '%s' "$spec" | cut -d'|' -f1)
		old=$(printf '%s' "$spec" | cut -d'|' -f2)
		new=$(printf '%s' "$spec" | cut -d'|' -f3)
		awk -v leg="$leg" -v old="$old" -v new="$new" '
			$1=="want" && $2==leg && $3==old { $3=new; print; next } { print }' \
			"$T/$out.sh" > "$T/$out.tmp"
		cmp -s "$T/$out.tmp" "$T/$out.sh" && { echo "   !! 翻 $leg 的声明没命中"; exit 9; }
		mv "$T/$out.tmp" "$T/$out.sh"
	done
}

flip f1 'G7pos|quiet|ring'
flip f2 'G7pos|quiet|ring' 'G6neg|ring|quiet'
# f3＝把 f2 翻回去（还原发）：在已翻的尺上再做一次反向替换
awk '$1=="want" && $2=="G7pos" && $3=="ring" {$3="quiet"}
     $1=="want" && $2=="G6neg" && $3=="quiet" {$3="ring"} {print}' "$T/f2.sh" > "$T/f3.sh"

echo "###### 声明机制三发（同一把新尺、同一枚锚 $A）"
for v in f0 f1 f2; do
	rc=0
	sh "$T/$v.sh" "$A" > "$S/flip-$v.txt" 2>&1 || rc=$?
	printf '   %s rc=%s  %s\n' "$v" "$rc" "$(grep -E '^# 腿数＝' "$S/flip-$v.txt")"
	grep -E '^# BAD|^# SHR' "$S/flip-$v.txt" | sed 's|^|        |' || true
done
echo "   f3＝f2 翻回原样（还原）：cmp 尺的字节 ⇒ rc 必须回 0"
cmp -s "$T/f3.sh" "$T/f0.sh" && echo "     还原后的尺与基线逐字节相同" || echo "     !! 还原不等于基线"
rc=0; sh "$T/f3.sh" "$A" > "$S/flip-f3.txt" 2>&1 || rc=$?
echo "     f3 rc=$rc"

echo
echo "###### 减法侧的第二形：整腿删掉（去掉 G7-正控 那一发 pair 调用）⇒ 腿数 14→13，退码动不动"
awk '/^pair "G7-正控/{skip=2} skip>0{skip--; next} {print}' "$T/gate.sh" > "$T/legs-dropped.sh"
grep -c '^want G7pos' "$T/legs-dropped.sh" | sed 's|^|     剩下的 want G7pos 行数＝|'
rc=0; sh "$T/legs-dropped.sh" "$A" > "$S/legs-dropped.txt" 2>&1 || rc=$?
echo "   rc=$rc  $(grep -E '^# 腿数＝' "$S/legs-dropped.txt" || echo '（没出腿数行）')"
