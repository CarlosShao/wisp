#!/bin/sh
# r2-roster-diff.sh - ticket 171 r2, the "did this leg's RANGE shrink?" evidence, same anchor both sides.
#
# WHY IT EXISTS
# AC#7/AC#8 touch the shared want/want_n/legs machinery, so the standing red line of this ticket
# (AC#6-2, repeated in AC#8's wording) applies: no reading may be made green by dropping a leg,
# a pathspec, or a sample out of range. The honest answer is a set difference, not a sentence:
# run the PRE-FIX ruler (git show $BASE:...gate-clauses.sh) and the working-copy ruler against the
# SAME anchor, then diff every leg's roster both ways.
#
# WHAT A ROSTER IS PER LEG
#   run legs (G1 G1b G2 G3 G4)  -> the git grep lines themselves (path:line:text) - 命中行数 的本体
#   pair legs (G5.. G7..)       -> the "#   UNPAIRED ..." + "#   UNPAIRED-REV ..." name lists
#
# HARD EXIT RULES: rc=0 both directions empty (no leg names anything new or lost anything);
# rc=1 a name moved (a finding to report, never a knob to relax); rc=2 the instrument broke.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r2-roster-diff.sh: no .git above $here" >&2; exit 2; }
    d=$p
done
root=$d
gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
BASE="${BASE:-f1b99a70c3dfdad04f34094306b5942a06355f52}"
anchor="${ANCHOR:-$(git -C "$root" rev-parse HEAD)}"
logdir="$here/logs"
mkdir -p "$logdir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r2-roster.XXXXXX")
git -C "$root" show "$BASE:$gate_rel" >"$work/gate-pre.sh"
[ -s "$work/gate-pre.sh" ] || { echo "r2-roster-diff.sh: no gate at $BASE - no verdict" >&2; exit 2; }

rc_pre=0
(cd "$root" && sh "$work/gate-pre.sh" "$anchor") >"$logdir/gate-pre-at-anchor.txt" 2>&1 || rc_pre=$?
rc_post=0
(cd "$root" && sh "$gate_rel" "$anchor") >"$logdir/gate-post-at-anchor.txt" 2>&1 || rc_post=$?
printf 'r2-roster-diff.sh: anchor %s  pre-fix ruler (git show %s) rc=%s  working-copy ruler rc=%s\n' \
    "$(git -C "$root" rev-parse --short HEAD)" "$BASE" "$rc_pre" "$rc_post"

sec() { # <prefix> <log>
    awk -v want="$1" '
        /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
        ins { print }
    ' "$2"
}

roster() { # <prefix> <log> -> the leg's name list, normalised
    sec "$1" "$2" | grep -E '^#   UNPAIRED|^#   DISCARD-HANDLE' ||
        sec "$1" "$2" | grep -E '^[^#].*:[0-9]+:' || true
}

fail=0
sum_pre=0
sum_post=0
# ⚠ 腿号 → 段名前缀不是一一对得上的：每一枚 pair 腿的正文里先有一行 echo "## G5 正控（…）"，
#   紧接着 pair() 自己打 "## G5-正控 …"。awk 的段界是「下一个 ## 结束」，所以拿 echo 那行当
#   前缀只能切到一行标题、名册恒空＝**一枚恒真比较**（本程第一版就死在这里，六枚对照腿全报 0/0）。
#   下面的前缀逐枚取 pair()/run() 自己打出的那一行标题头。
for leg in 'G1|## G1 生产码' 'G1b|## G1b' 'G2|## G2 OpenTask' 'G3|## G3 Q-56' \
	'G4|## G4 生产组合根' 'G5|## G5 OpenScope' 'G5pos|## G5-正控' 'G5neg|## G5-负一负' \
	'G6|## G6 主尺' 'G6pos|## G6-正控' 'G6neg|## G6-负一负' \
	'G7|## G7 主尺' 'G7pos|## G7-正控' 'G7neg|## G7-负一负'; do
    id=${leg%%|*}
    pref=${leg#*|}
    roster "$pref" "$logdir/gate-pre-at-anchor.txt" | sort -u >"$work/pre-$id.txt"
    roster "$pref" "$logdir/gate-post-at-anchor.txt" | sort -u >"$work/post-$id.txt"
    added=$(comm -13 "$work/pre-$id.txt" "$work/post-$id.txt" | grep -c . || true)
    lost=$(comm -23 "$work/pre-$id.txt" "$work/post-$id.txt" | grep -c . || true)
    npre=$(grep -c . "$work/pre-$id.txt" || true)
    npost=$(grep -c . "$work/post-$id.txt" || true)
    sum_pre=$((sum_pre + npre))
    sum_post=$((sum_post + npost))
    printf '  腿=%-6s 改前名册=%-3s 改后名册=%-3s 新增点名=%s 退场=%s\n' "$id" "$npre" "$npost" "$added" "$lost"
    if [ "$added" != 0 ] || [ "$lost" != 0 ]; then
        printf '     新增/退场逐字：\n'
        comm -3 "$work/pre-$id.txt" "$work/post-$id.txt" | sed -e 's/^/       /'
        fail=1
    fi
done
printf '  合计名册行数：改前=%s 改后=%s   （两向 comm 全空＝射程一枚未动）\n' "$sum_pre" "$sum_post"
printf '  pathspec 本体是否被改：git diff %s HEAD -- %s 里含 ":!" 或 "$GO" 的增删行：\n' "$BASE" "$gate_rel"
git -C "$root" diff "$BASE" HEAD -- "$gate_rel" | grep -E '^[+-].*(:!|\$GO|\$GO2|pathspec)' | sed -e 's/^/       /' || printf '       （零枚——没有一条射程行被增删）\n'
if [ "$fail" = 0 ]; then
    printf 'r2-roster-diff.sh: GREEN - 十四枚腿的名册两向差集全空，一条射程行未动。\n'
    exit 0
fi
printf 'r2-roster-diff.sh: RED - 有名册位移，按票面那是必须上报的发现，不是可以放宽的旋钮。\n'
exit 1
