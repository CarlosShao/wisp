#!/bin/sh
# 171-v1 台件⑥（AC#6 主攻③那一问）：来一发「句柄被丢掉 ＋ 同文件另有一句无关的 Close()」，
# 今天的尺会不会又静默？——同一枚丢弃树，两把尺，同一条射程。
set -eu
D=.scratch/wisp/probes/171/v1
A=3df749607a0da1e70579f5c3abb71003350a45d2
T=$(mktemp -d)
mkdir -p "$D/logs"
S=$(cd "$D/logs" && pwd)
git show 19513cce:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-old.sh"
git show 3df74960:.scratch/wisp/probes/154/gate-clauses.sh > "$T/gate-new.sh"
mkdir -p "$T/base"
git archive "$A" internal cmd | tar -x -C "$T/base"
mkdir -p "$T/base/internal/fake171v1d"
printf '%s\n' \
	'package fake171v1d' \
	'' \
	'import "wisp/internal/risk"' \
	'' \
	'func alsoClosesSomething(p *risk.Provenance, f io.Closer, id string) {' \
	'	_ = p.OpenScope(id)' \
	'	_ = f.Close()   // 与这族毫无关系的泛用 Close()' \
	'}' > "$T/base/internal/fake171v1d/d_discard_plus_generic_close.go"
( cd "$T/base" && git init -q && git add internal cmd && git -c user.name=v1 -c user.email=v1@invalid commit -q -m S-D )
for v in old new; do
	rc=0
	( cd "$T/base" && export GIT_DIR="$T/base/.git" GIT_WORK_TREE="$T/base" && sh "$T/gate-$v.sh" HEAD > "$S/s6-discard-plus-close-$v.txt" 2>&1 ) || rc=$?
	printf '   ruler=%s rc=%s 点名这一枚文件的腿=[%s]\n' "$v" "$rc" \
		"$(awk '/^## /{k=$2} /^#   UNPAIRED internal\/fake171v1d/{print k}' "$S/s6-discard-plus-close-$v.txt" | sort -u | tr '\n' ' ')"
done
echo "   新尺为它摊的批注："
grep -E 'DISCARD-HANDLE|CLOSE-BY-GENERIC-CLOSE' "$S/s6-discard-plus-close-new.txt" | grep fake171v1d | sed 's|^|     |' || echo "     （无批注）"
