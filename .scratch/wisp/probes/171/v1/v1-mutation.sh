#!/bin/sh
# 171-v1 台件③：变异树三连（全在仓外 $TMPDIR 的丢弃 git 树里，本仓 internal/** 零字节）。
#   M1 unpaired-bridge ：把 internal/tools/bridge.go:706 `closeErr = scope.Close()` 整行摘掉
#                        ⇒「那枚文件只开不关」⇒ 新尺必须【重新点名】它（AC#6② 的判死发）。
#   M2 extra-leak      ：bridge.go 原样，另加一枚「接住句柄、全文件不关」的新文件
#                        ⇒ 同一条腿内未成对 1 枚 → 2 枚 ⇒ 新尺聚合退码必须动（AC#2 有牙那一发）。
#   M3 control         ：internal+cmd 逐字节＝锚点的拷贝 ⇒ 两把尺在此的读数＝基线。
# 尺的字节仍旧走 git show（old=19513cce／new=3df74960），锚点＝各枚丢弃树自己的 HEAD。
set -eu
D=.scratch/wisp/probes/171/v1
A=3df749607a0da1e70579f5c3abb71003350a45d2
T=$(mktemp -d)
mkdir -p "$D/logs"
S=$(cd "$D/logs" && pwd)   # 绝对路径：跑尺时要 cd 进丢弃树，相对路径会落到别的仓
git show 19513cce:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-old.sh"
git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-new.sh"

mk_tree() { # $1 目录名
	cp -r "$T/base" "$T/$1"
	git -C "$T/$1" init -q
	git -C "$T/$1" add internal cmd
	git -C "$T/$1" -c user.name=v1 -c user.email=v1@invalid commit -q -m "$1"
	echo "# tree $1 sha=$(git -C "$T/$1" rev-parse HEAD)"
}

rm -rf "$T/base" 2>/dev/null || true
mkdir -p "$T/base"
git archive "$A" internal cmd | tar -x -C "$T/base"
mk_tree control

# M1：摘掉 bridge.go 那一行 Close()
mkdir -p "$T/tmp1"
sed '/closeErr = scope\.Close()/d' "$T/base/internal/tools/bridge.go" > "$T/tmp1/bridge.go"
cmp -s "$T/tmp1/bridge.go" "$T/base/internal/tools/bridge.go" && { echo "MUTATION FAILED: 行没摘掉"; exit 9; }
grep -c 'Close()' "$T/tmp1/bridge.go" | sed 's|^|# M1 bridge.go 剩 Close() 次数＝|'
mk_tree unpaired-bridge
cp "$T/tmp1/bridge.go" "$T/unpaired-bridge/internal/tools/bridge.go"
git -C "$T/unpaired-bridge" add internal
git -C "$T/unpaired-bridge" -c user.name=v1 -c user.email=v1@invalid commit -q -m M1
git -C "$T/unpaired-bridge" show --stat --oneline HEAD | head -5

# M2：加一枚「接住句柄、全文件不关」的文件（形状照 bridge.go：赋值给字段/变量，从不 Close）
mk_tree extra-leak
mkdir -p "$T/extra-leak/internal/fake171v1"
printf '%s\n' \
	'package fake171v1' \
	'' \
	'import "wisp/internal/risk"' \
	'' \
	'type holder struct { scopes map[string]*risk.Scope }' \
	'' \
	'func (h *holder) leak(taskID string) {' \
	'	h.scopes[taskID] = (&risk.Provenance{}).OpenScope(taskID)' \
	'}' > "$T/extra-leak/internal/fake171v1/leak.go"
gofmt -w "$T/extra-leak/internal/fake171v1/leak.go"
gofmt -l "$T/extra-leak/internal/fake171v1/leak.go" | sed 's|^|# gofmt 未格式化: |' || true
git -C "$T/extra-leak" add internal/fake171v1
git -C "$T/extra-leak" -c user.name=v1 -c user.email=v1@invalid commit -q -m M2

for tree in control unpaired-bridge extra-leak; do
	for v in old new; do
		rc=0
		( cd "$T/$tree" && export GIT_DIR="$T/$tree/.git" GIT_WORK_TREE="$T/$tree" && sh "$T/gate-$v.sh" HEAD > "$S/mut-$tree-$v.txt" 2>&1 ) || rc=$?
		echo "== tree=$tree ruler=$v rc=$rc anchorsha=$(git -C "$T/$tree" rev-parse --short HEAD)"
	done
done
echo
echo "###### M1 判死读数：internal/tools/bridge.go 在不在 G5 主尺/正控的点名词里"
for tree in control unpaired-bridge; do
	for v in old new; do
		n=$(awk '/^## /{k=$2} /^#   UNPAIRED internal\/tools\/bridge\.go/{print k}' "$S/mut-$tree-$v.txt" | sort | tr '\n' ' ')
		echo "   tree=$tree ruler=$v 点名腿=[${n:-无}]"
	done
done
echo "   新尺在 M1 里另摊出来的批注（DISCARD-HANDLE / CLOSE-BY-GENERIC-CLOSE）："
grep -E 'DISCARD-HANDLE|CLOSE-BY-GENERIC-CLOSE' "$S/mut-unpaired-bridge-new.txt" | sed 's|^|     |' || echo "     （无）"
echo
echo "###### M2 有牙读数：G5 那一条腿的实测枚数 + 聚合退码（control → extra-leak）"
for v in old new; do
	for tree in control extra-leak; do
		printf '   ruler=%s tree=%-13s %s %s\n' "$v" "$tree" \
			"$(grep -E '^# (BAD|SHR|ok) +(腿=)?G5 ' "$S/mut-$tree-$v.txt" | head -1)" \
			"$(grep '^# 聚合退码＝' "$S/mut-$tree-$v.txt")"
	done
done
echo "   M2 里新尺逐枚点名（G5 主尺名册）："
awk '/^## G5 OpenScope/{f=1} /^## G5 正控/{f=0} f&&/^#   UNPAIRED|^# 未成对枚数/' "$S/mut-extra-leak-new.txt" | sed 's|^|     |'
