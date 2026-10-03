#!/bin/sh
# r3-roster-diff.sh - ticket 171 r3, the "did AC#1 change the RANGE or only the NUMERATOR?"
# evidence, both rulers打在 the SAME anchor.
#
# WHY IT EXISTS
# AC#1 rewrites what counts as a call site for BOTH leg families. That is a numerator change,
# and a numerator change can hide a range change (a dropped pathspec also makes numbers get
# smaller, and this ticket's standing red line - AC#6-2, repeated in AC#8 - is that no reading
# may be made green by dropping a leg, a pathspec or a sample). The honest answer is a set
# difference with a reason per name, not a sentence:
#
#   check 1  名册本体（git grep 的逐行读数）改前 vs 改后 **逐字节相同**，十四枚腿全如此。
#            名册是射程的产物；它一字不动 ⇒ 射程一字未动。
#   check 2  pathspec 本体（含 ':!' 与 $GO 的行）在两版尺的 diff 里增删行数 = 0。
#   check 3  每一腿的点名句单（UNPAIRED / UNPAIRED-REV）允许位移，但**每一枚位移都必须被
#            同一腿自己打出的 COMMENT-ONLY 格里那一行解释**（同文件、被剔的是注释行）。
#            位移而无解释 = 射程动了 = 本程判失败，不是"下一程再核"。
#   check 4  run 腿打出两个数：# 名册行数＝（未动）与 # 命中行数＝（分子）。名册行数两版必须相同。
#
# HARD EXIT RULES: rc=0 all four checks hold / rc=1 a movement has no comment-only explanation
# (that would be a range change, a finding, never a knob) / rc=2 the instrument broke.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r3-roster-diff.sh: no .git above $here" >&2; exit 2; }
    d=$p
done
root=$d
gate_rel=.scratch/wisp/probes/154/gate-clauses.sh

# BASE = a commit whose tracked ruler is still the unmodified bytes. Pinned and integrity-checked
# against this leg's own snapshot (gate-prefix-f78d3cb.sh) so "pre-fix" is a byte claim, not a date.
BASE="${BASE:-f78d3cb183b041e6211d3955fdfb0f5e96098d6b}"
ANCHOR="${ANCHOR:-4f2a777b60770744adabbd92144d261812a025a0}"
logdir="$here/logs"
mkdir -p "$logdir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3-roster.XXXXXX")

git -C "$root" show "$BASE:$gate_rel" >"$work/gate-pre.sh"
[ -s "$work/gate-pre.sh" ] || { echo "r3-roster-diff.sh: no gate at $BASE - no verdict" >&2; exit 2; }
if ! cmp -s "$work/gate-pre.sh" "$here/gate-prefix-f78d3cb.sh"; then
    echo "r3-roster-diff.sh: git show $BASE of the gate is not the snapshot this leg saved - the pre-fix side is not the same ruler, no verdict" >&2
    exit 2
fi

rc_pre=0
(cd "$root" && sh "$work/gate-pre.sh" "$ANCHOR") >"$logdir/roster-gate-pre.txt" 2>&1 || rc_pre=$?
rc_post=0
(cd "$root" && sh "$gate_rel" "$ANCHOR") >"$logdir/roster-gate-post.txt" 2>&1 || rc_post=$?
printf 'r3-roster-diff.sh: anchor=%s  改前尺(git show %s) 退码=%s  改后尺(工作副本) 退码=%s\n' \
    "$ANCHOR" "$BASE" "$rc_pre" "$rc_post"

sec() {
    awk -v want="$1" '/^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 } ins { print }' "$2"
}
names() { # <段前缀> <log> -> 这一腿点名的文件（正向＋反向），逐枚一行
    sec "$1" "$2" | grep -E '^#   (UNPAIRED|UNPAIRED-REV|DISCARD-HANDLE|CLOSE-BY-GENERIC-CLOSE) ' |
        sed -e 's/^#   [A-Z-]* //' -e 's/ (.*//' -e 's/[[:space:]]*#.*$//' || true
}
comment_only() { # <段前缀> <log> -> 这一腿自己打出的被剔行里的文件名
    sec "$1" "$2" | grep '^#   COMMENT-ONLY ' |
        sed -e 's/^#   COMMENT-ONLY //' -e 's/^[^:]*://' -e 's/:[0-9][0-9]*:.*$//' | sort -u || true
}
grep_lines() { # <段前缀> <log> -> 名册本体（git grep 的逐行读数，不含 # 打的注释行）
    sec "$1" "$2" | grep -E '^[^#[:space:]].*:[0-9]+:' || true
}
raw_count() { sec "$1" "$2" | grep -m1 '^# 名册行数＝' | sed -e 's/^# 名册行数＝\([0-9]*\).*/\1/'; }
num_count() { sec "$1" "$2" | grep -m1 '^# 命中行数＝' | sed -e 's/^# 命中行数＝\([0-9]*\).*/\1/'; }

fail=0
say() { printf '%s\n' "$1"; }
# 腿号 -> 段名前缀。⚠ 一枚腿的正文里先有一行 echo "## G5 正控（…）"，紧接着 pair() 自己打
#   "## G5-正控 …"；awk 的段界是"下一个 ## 结束"，拿 echo 那行当前缀只能切到一行标题、名册恒空
#   ＝一枚恒真比较（票 171 r2 的自纠就死在这里，见 docs/evidence/s1/171-run-legs-and-subtraction-r2.md
#   的 next= 那条）。下面的前缀逐枚取 run()/pair() 自己打出的那一行标题头。
LEG_ALL='G1|## G1 生产码 G1b|## G1b G2|## G2 OpenTask G3|## G3 Q-56 G4|## G4 生产组合根 G5|## G5 OpenScope G5pos|## G5-正控 G5neg|## G5-负一负 G6|## G6 主尺 G6pos|## G6-正控 G6neg|## G6-负一负 G7|## G7 主尺 G7pos|## G7-正控 G7neg|## G7-负一负'
LEGS_WRITTEN=0
legs_of() { # <all|run> -> 每腿一行 "<腿号>|<段名前缀>"，前缀里的空格由 read -r 保留
    case $1 in
    all)
        printf '%s\n' 'G1|## G1 生产码' 'G1b|## G1b' 'G2|## G2 OpenTask' 'G3|## G3 Q-56' 'G4|## G4 生产组合根' \
            'G5|## G5 OpenScope' 'G5pos|## G5-正控' 'G5neg|## G5-负一负' \
            'G6|## G6 主尺' 'G6pos|## G6-正控' 'G6neg|## G6-负一负' \
            'G7|## G7 主尺' 'G7pos|## G7-正控' 'G7neg|## G7-负一负'
        ;;
    run)
        printf '%s\n' 'G1|## G1 生产码' 'G1b|## G1b' 'G2|## G2 OpenTask' 'G3|## G3 Q-56' 'G4|## G4 生产组合根'
        ;;
    *) echo "legs_of: unknown set $1" >&2; exit 2 ;;
    esac
}

legs_of all >"$work/legs-all.txt"
legs_of run >"$work/legs-run.txt"

say ""
say "== check 1 名册本体（射程的产物）两版逐字节相同"
while IFS= read -r pair_; do
    id=${pair_%|*}
    pref=${pair_#*|}
    grep_lines "$pref" "$logdir/roster-gate-pre.txt" >"$work/rpre-$id.txt"
    grep_lines "$pref" "$logdir/roster-gate-post.txt" >"$work/rpost-$id.txt"
    if cmp -s "$work/rpre-$id.txt" "$work/rpost-$id.txt"; then
        say "   ok    腿=$id 名册 $(wc -l <"$work/rpre-$id.txt") 行逐字节相同"
    else
        say "   FAIL  腿=$id 名册动了（射程被动了！）："
        diff "$work/rpre-$id.txt" "$work/rpost-$id.txt" | sed -e 's/^/         /'
        fail=1
    fi
done <"$work/legs-all.txt"

say ""
say "== check 2 pathspec 本体在 diff 里零枚增删"
# 两枚防空转的正控（本仓教训："零命中"要先怀疑尺，别直接读成"没动"）：
#   (a) 这一味词表在尺自己身上必须打得出东西——打不出＝正则是一枚哑弹，0 不算读数；
#   (b) 两版尺的 diff 本身必须非空——空 diff＝两把尺是同一份字节，什么都没比。
nrange_pat=$(grep -cE "(:!|\$GO|\$GO2)" "$gate_rel" || true)
ndiff_all=$(git -C "$root" diff "$BASE" -- "$gate_rel" | grep -c '^[+-][^+-]' || true)
say "   正控 (a) 尺自己含 ':!'／\$GO／\$GO2 的行数＝$nrange_pat（必须 >0，否则下面那个 0 是哑弹）"
say "   正控 (b) 两版尺的 diff 增删行数＝$ndiff_all（必须 >0，否则两把尺是同一份字节）"
{ [ "$nrange_pat" -gt 0 ] && [ "$ndiff_all" -gt 0 ]; } || { say "   尺是哑弹或 diff 是空的——不给出读数"; exit 2; }
nshrink=$(git -C "$root" diff "$BASE" -- "$gate_rel" | grep -E '^[+-].*(:!|\$GO|\$GO2)' | grep -vE '^[+-][[:space:]]*#' | grep -c . || true)
say "   含 ':!'／\$GO／\$GO2 的增删行（注释行除外）＝$nshrink"
git -C "$root" diff "$BASE" -- "$gate_rel" | grep -E '^[+-].*(:!|\$GO|\$GO2)' | sed -e 's/^/     /' || true
[ "$nshrink" = 0 ] || fail=1

say ""
say "== check 3 点名句单的每一枚位移都必须被同一腿的 COMMENT-ONLY 那一格解释"
while IFS= read -r pair_; do
    id=${pair_%|*}
    pref=${pair_#*|}
    names "$pref" "$logdir/roster-gate-pre.txt" | sort -u >"$work/npre-$id.txt"
    names "$pref" "$logdir/roster-gate-post.txt" | sort -u >"$work/npost-$id.txt"
    added=$(comm -13 "$work/npre-$id.txt" "$work/npost-$id.txt" || true)
    lost=$(comm -23 "$work/npre-$id.txt" "$work/npost-$id.txt" || true)
    comment_only "$pref" "$logdir/roster-gate-post.txt" >"$work/cmt-$id.txt"
    if [ -z "$added" ] && [ -z "$lost" ]; then
        say "   ok    腿=$id 点名未动（$(grep -c . "$work/npost-$id.txt" || true) 枚）"
        continue
    fi
    for f in $added $lost; do
        if grep -qxF "$f" "$work/cmt-$id.txt"; then
            say "   ok    腿=$id 位移有解释：$f（该文件在本腿的 COMMENT-ONLY 格里，被剔的是注释行）"
        else
            say "   FAIL  腿=$id 位移无解释：$f（本腿的 COMMENT-ONLY 格里没有这一枚＝射程动了）"
            fail=1
        fi
    done
done <"$work/legs-all.txt"

say ""
say "== check 4 run 腿：名册行数未动、命中行数（分子）动了"
while IFS= read -r pair_; do
    id=${pair_%|*}
    pref=${pair_#*|}
    rp=$(grep_lines "$pref" "$logdir/roster-gate-pre.txt" | grep -c . || true)
    rq=$(grep_lines "$pref" "$logdir/roster-gate-post.txt" | grep -c . || true)
    post_num=$(num_count "$pref" "$logdir/roster-gate-post.txt")
    post_raw=$(raw_count "$pref" "$logdir/roster-gate-post.txt")
    say "   腿=$id 改前名册行=$rp 改后名册行=$rq 改后分子=$post_num（改前那一格没有「名册行数」这一行，因为改前的分子就是名册行数本身）"
    [ "$rp" = "$rq" ] || { say "     FAIL 名册行数动了＝射程动了"; fail=1; }
    [ "$post_raw" = "$rq" ] || { say "     FAIL # 名册行数＝ 与实际名册行数不符"; fail=1; }
done <"$work/legs-run.txt"

say ""
if [ "$fail" = 0 ]; then
    say "r3-roster-diff.sh: GREEN - 射程零动、名册零动；点名句单的位移逐枚都有 COMMENT-ONLY 的解释。"
    exit 0
fi
say "r3-roster-diff.sh: RED - 有一枚位移解释不了，或射程/名册动过。按票面那是发现，不是旋钮。"
exit 1
