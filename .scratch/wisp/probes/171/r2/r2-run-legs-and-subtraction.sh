#!/bin/sh
# r2-run-legs-and-subtraction.sh - ticket 171 r2 (写码位·仪器程), the readings AC#7 and AC#8 demand.
#
# WHAT IT ANSWERS (three cells, all on the SAME pinned pre-fix ruler vs the working-copy ruler)
#   AC#7  a run() leg (G1-G4) still only knows "空／非空": does the aggregate exit code move when
#         ONE leg's hits go 2 枚 -> 3 枚? 未修码上不许动，修完之后必须动。
#   AC#8① 整条腿被摘掉：未修码上退码不许动（腿数只打印、无断言），修完之后必须动。
#   AC#8② SHR（实测 < 基线）允许变好，但必须单列一张读得到的『基线过期』表；未修码上只有人眼看得见。
#
# WHY THE PRISTINE RULER IS TAKEN OUT OF A PINNED COMMIT
# "未修码上必须不响" is only a judgement if the reading comes from the pre-fix bytes, so the
# pristine gate is `git show $BASE:.scratch/wisp/probes/154/gate-clauses.sh` written into $TMPDIR
# and run side by side with the working copy. Same tree, same anchor, same pathspecs.
# BASE defaults to the anchor this leg started from (step 0 reading).
#
# WHY ONE SAMPLE SET LIVES IN A SYNTHETIC TREE OUTSIDE THE REPO
# run()/pair() read `git grep "$A"` where $A is a COMMIT pinned to internal/**/*.go + cmd/**/*.go.
# A sample dropped under .scratch/wisp/probes/171/** is invisible twice over. So AC#7's samples are
# committed into a throwaway repo under $TMPDIR - never inside this repo, no worktree, no checkout,
# no branch switch, no relaxed ':!' filter (precedent: 171-r1's r1-pair-shapes.sh, 161-r6's faketree.sh).
# Nothing under $TMPDIR is ever deleted by this script: states are ADDITIVE (state 2 adds a second
# sample file to the same tree) - 临时件只建不删.
#
# AC#8's two samples edit a COPY of the ruler in $TMPDIR (a leg deleted from a copy, a want_n raised
# on a copy) and run the copy with cwd = the repo root, so the READINGS come from the real committed
# tree while the REPO's own gate bytes are never touched.
#
# HARD EXIT RULES
#   rc=0  every assertion held
#   rc=1  an assertion did not hold -> the cell is decoration, or the ruler lies
#   rc=2  the instrument could not run -> no verdict, never read it as green
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r2-run-legs-and-subtraction.sh: no .git above $here - refusing to answer with a ruler that did not run" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
gate_fixed="$root/$gate_rel"
[ -f "$gate_fixed" ] || { echo "r2-run-legs-and-subtraction.sh: working-copy gate missing - no verdict" >&2; exit 2; }

# step-0 anchor of this leg (171-r2); the pre-fix ruler is the gate as committed there.
BASE="${BASE:-f1b99a70c3dfdad04f34094306b5942a06355f52}"
git -C "$root" cat-file -e "${BASE}^{commit}" 2>/dev/null || {
    echo "r2-run-legs-and-subtraction.sh: pinned pre-fix anchor $BASE is not a commit here - no verdict" >&2; exit 2; }

command -v git >/dev/null 2>&1 || { echo "no git" >&2; exit 2; }
command -v mktemp >/dev/null 2>&1 || { echo "no mktemp" >&2; exit 2; }

work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r2.XXXXXX")
tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r2-tree.XXXXXX")
logdir="$here/logs"
mkdir -p "$logdir"

gate_pristine="$work/gate-pristine.sh"
git -C "$root" show "$BASE:$gate_rel" >"$gate_pristine"

anchor=$(git -C "$root" rev-parse HEAD)

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

pre_fix_pass=0
if cmp -s "$gate_pristine" "$gate_fixed"; then
    pre_fix_pass=1
    say "!!  PRE-FIX PASS: the working-copy gate is byte-identical to git show $BASE - so every"
    say "!!  '修完之后' assertion below is EXPECTED to fail in this run. That failure is the reading"
    say "!!  the ticket demands (这一发今天不响); it is not this instrument misbehaving."
fi

sec() { # <prefix> <logfile> -> one leg's section, up to the next "## "
    awk -v want="$1" '
        /^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 }
        ins { print }
    ' "$2"
}

leg_field() { # <prefix> <log> <line-prefix> -> the integer on that line, or NO-READING
    raw=$(sec "$1" "$2" | grep -m1 "^$3" | sed -e "s/^$3\([0-9][0-9]*\).*$/\1/" || true)
    case $raw in
    '' | *[!0-9]*) printf 'NO-READING' ;;
    *) printf '%s' "$raw" ;;
    esac
}

run_in() { # <script> <cwd> <anchor-arg|''> <log-name> -> echoes the gate's exit code, log under $logdir
    script=$1; cwd=$2; arg=$3
    log="$logdir/$4"
    rc=0
    if [ -n "$arg" ]; then
        (cd "$cwd" && sh "$script" "$arg") >"$log" 2>&1 || rc=$?
    else
        (cd "$cwd" && sh "$script") >"$log" 2>&1 || rc=$?
    fi
    case $rc in
    *[!0-9]*) echo "run_in: exit code '$rc' is not a number - the runner broke, no verdict" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

# ============================================================================
# PART 1 - AC#7: 同一条 run() 腿内命中从 2 枚变 3 枚（基线＝本程在真树现量的 G2＝2 行）
# ============================================================================
say ""
say "===== PART 1 (AC#7): one run() leg, 2 枚 then 3 枚. Does the exit code move? ====="

export GIT_AUTHOR_NAME=wisp171r2 GIT_AUTHOR_EMAIL=171-r2@invalid
export GIT_COMMITTER_NAME=wisp171r2 GIT_COMMITTER_EMAIL=171-r2@invalid
git -C "$tree" init -q .
mkdir -p "$tree/cmd/wisp"

# The sample lines are shaped like the real production hit (cmd/wisp/run.go 那一发 CloseTask), and each
# one carries BOTH word roots ON PURPOSE: the pair legs (G6/G6neg) then read this file as balanced in
# both directions, so between the two states EXACTLY ONE leg's count moves - G2's. The samples declare
# no helper funcs with those word roots either: `git grep -nEw` counts a `func (x) CloseTask(...)`
# definition line as a hit too, and this leg's 枚数 is 命中行数 (本程现量：那样会多两行).
# (A comment-shaped sample would have been cheaper but would prove nothing about call sites:
#  run() counts every roster line, pair() strips whole-line comments.)
cat >"$tree/cmd/wisp/g2_state_one.go" <<'GO'
package wisp

// 样本是尺的靶子，不求编译（尺吃的是 git grep 的行，不是 go build 的图）。
func one(ids [2]string) {
	_ = CloseTask(OpenTask(ids[0]))
	_ = CloseTask(OpenTask(ids[1]))
}
GO
git -C "$tree" add -- cmd/wisp/g2_state_one.go >/dev/null
git -C "$tree" commit -q -m "171-r2 AC#7 state1: exactly 2 hits inside leg G2" || {
    echo "commit failed in synthetic tree" >&2; exit 2; }
say "-- state1 committed at $(git -C "$tree" rev-parse --short HEAD) (synthetic tree, outside $root)"

rc1p=$(run_in "$gate_pristine" "$tree" '' "part1-state1-pristine.txt")
rc1f=$(run_in "$gate_fixed" "$tree" '' "part1-state1-fixed.txt")
h1p=$(leg_field '## G2 OpenTask' "$logdir/part1-state1-pristine.txt" '# 命中行数＝')
h1f=$(leg_field '## G2 OpenTask' "$logdir/part1-state1-fixed.txt" '# 命中行数＝')
true1=$(git -C "$tree" grep -nEw 'OpenTask|CloseTask' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/tools/bridge.go' | grep -c .)
say "   state1: leg G2 really measures ${true1} 行 | pristine prints 命中行数=${h1p} | fixed prints 命中行数=${h1f}"
say "   state1 aggregate rc: pristine=${rc1p} fixed=${rc1f}"

cat >"$tree/cmd/wisp/g2_state_two.go" <<'GO'
package wisp

// 第三行命中，同一条腿、同一射程、同一条声明。
func two(id string) {
	_ = CloseTask(OpenTask(id))
}
GO
git -C "$tree" add -- cmd/wisp/g2_state_two.go >/dev/null
git -C "$tree" commit -q -m "171-r2 AC#7 state2: 2 hits -> 3 hits inside leg G2" || {
    echo "commit failed in synthetic tree" >&2; exit 2; }
say "-- state2 committed at $(git -C "$tree" rev-parse --short HEAD) (same tree, nothing deleted)"

rc2p=$(run_in "$gate_pristine" "$tree" '' "part1-state2-pristine.txt")
rc2f=$(run_in "$gate_fixed" "$tree" '' "part1-state2-fixed.txt")
h2p=$(leg_field '## G2 OpenTask' "$logdir/part1-state2-pristine.txt" '# 命中行数＝')
h2f=$(leg_field '## G2 OpenTask' "$logdir/part1-state2-fixed.txt" '# 命中行数＝')
true2=$(git -C "$tree" grep -nEw 'OpenTask|CloseTask' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/tools/bridge.go' | grep -c .)
say "   state2: leg G2 really measures ${true2} 行 | pristine prints 命中行数=${h2p} | fixed prints 命中行数=${h2f}"
say "   state2 aggregate rc: pristine=${rc2p} fixed=${rc2f}"
say "   which legs the FIXED aggregate calls out in state2:"
grep -E '^# (BAD|SHR|STALE) ' "$logdir/part1-state2-fixed.txt" | sed -e 's/^/       /' || say "       (none)"

if [ "$true1" = 2 ] && [ "$true2" = 3 ]; then
    ok "the two states differ by exactly ONE hit inside leg G2 (${true1} -> ${true2}), nothing else in the tree moved"
else
    bad "the states measured ${true1}/${true2} rows, expected 2/3 - the two runs are not comparable"
fi
if [ "$rc1p" = "$rc2p" ]; then
    ok "UNMODIFIED CODE (AC#7 hole reproduced): 2 枚 -> 3 枚 moves the aggregate exit code NOT AT ALL (rc=${rc1p} both ways) - run() legs still only know 空/非空"
else
    bad "UNMODIFIED CODE already reacted (rc ${rc1p} -> ${rc2p}) - then AC#7 is decoration; report it, do not fix the judgement to match"
fi
if [ "$rc1f" != "$rc2f" ]; then
    ok "FIXED CODE: the same jump moves the aggregate exit code (${rc1f} -> ${rc2f}) - a second hit can no longer hide inside 非空"
else
    bad "FIXED CODE did not react (rc=${rc1f} both) - AC#7's 枚数 judgement is decoration (expected in the PRE-FIX PASS)"
fi
if [ "$rc2p" = $((rc2f - 1)) ] || { [ "$pre_fix_pass" = 1 ] && [ "$rc2p" = "$rc2f" ]; }; then
    say "   --     the fixed ruler's extra BAD is attributable to leg G2 alone (see the list above)"
fi

# ============================================================================
# PART 2 - AC#8①: 整条腿被摘掉（从尺的副本里删，仓里的字节不动）
# ============================================================================
say ""
say "===== PART 2 (AC#8-1): a whole leg deleted. Does the exit code move? ====="

drop_leg() { # <in> <out> <leg-id> -> drops the leg's 4 lines (want / want_n / pair / its arg line)
    awk -v head="want $3 " '
        index($0, head) == 1 && !skipping { skipping = 1; left = 4; next }
        skipping { left--; if (left <= 0) skipping = 0; next }
        { print }
    ' "$1" > "$2"
}

for which in pristine fixed; do
    src=$gate_pristine; [ "$which" = fixed ] && src=$gate_fixed
    drop_leg "$src" "$work/gate-nog6-$which.sh" G6
    if cmp -s "$src" "$work/gate-nog6-$which.sh"; then
        bad "dropping leg G6 changed not one byte of the $which ruler - that is not a reading"
    fi
    if grep -qE '^want G6 ' "$work/gate-nog6-$which.sh"; then
        bad "leg G6 survived the drop in the $which copy - the sample is a no-op"
    fi
done
rcp2=$(run_in "$work/gate-nog6-pristine.sh" "$root" "$anchor" "part2-nog6-pristine.txt")
rcf2=$(run_in "$work/gate-nog6-fixed.sh" "$root" "$anchor" "part2-nog6-fixed.txt")
say "   pristine: 摘掉 G6 -> rc=${rcp2}  $(grep -oE '腿数＝[0-9]+ 声明与实测不符＝[0-9]+' "$logdir/part2-nog6-pristine.txt" | tail -1)"
say "   fixed:    摘掉 G6 -> rc=${rcf2}  $(grep -oE '腿数＝[0-9]+ 声明与实测不符＝[0-9]+' "$logdir/part2-nog6-fixed.txt" | tail -1)"
say "   which legs the FIXED aggregate calls out:"
grep -E '^# BAD ' "$logdir/part2-nog6-fixed.txt" | sed -e 's/^/       /' || say "       (none)"
if [ "$rcp2" = 0 ]; then
    ok "UNMODIFIED CODE (AC#8-1 hole reproduced): a leg deleted outright, aggregate exit code still 0 - 腿数 only ever printed, nothing asserted"
else
    bad "UNMODIFIED CODE already rung on a deleted leg (rc=${rcp2}) - then AC#8-1 is decoration; report it"
fi
if [ "$rcf2" != 0 ]; then
    ok "FIXED CODE: deleting the same leg moves the exit code (0 -> ${rcf2})"
else
    bad "FIXED CODE stayed 0 with a leg missing - AC#8-1's 腿数 assertion is decoration (expected in the PRE-FIX PASS)"
fi

# ============================================================================
# PART 3 - AC#8②: SHR 允许变好，但必须单列一张『基线过期』表
# ============================================================================
say ""
say "===== PART 3 (AC#8-2): 实测 < 基线 - allowed to stay quiet, must print a readable table ====="

for which in pristine fixed; do
    src=$gate_pristine; [ "$which" = fixed ] && src=$gate_fixed
    sed 's/^want_n 8$/want_n 9/' "$src" >"$work/gate-stale-$which.sh"
    if cmp -s "$src" "$work/gate-stale-$which.sh"; then
        bad "raising G5neg's baseline 8 -> 9 changed no byte of the $which ruler - not a reading"
    fi
done
rcp3=$(run_in "$work/gate-stale-pristine.sh" "$root" "$anchor" "part3-stale-pristine.txt")
rcf3=$(run_in "$work/gate-stale-fixed.sh" "$root" "$anchor" "part3-stale-fixed.txt")
say "   pristine: 基线 8->9（实测 8）-> rc=${rcp3}  STALE 表行数=$(grep -c '^# STALE ' "$logdir/part3-stale-pristine.txt" || true)"
say "   fixed:    基线 8->9（实测 8）-> rc=${rcf3}  STALE 表行数=$(grep -c '^# STALE ' "$logdir/part3-stale-fixed.txt" || true)"
say "   fixed 那一份的表（原样）："
awk '/^## 基线过期/,/^# 基线过期枚数/' "$logdir/part3-stale-fixed.txt" | sed -e 's/^/       /' || say "       （没有表——这一格就是空的）"
if [ "$rcp3" = 0 ] && [ "$rcf3" = 0 ]; then
    ok "变好仍然不计进退码（两把尺都 rc=0）：这一格没有把门换成'为好消息响'的恒红形状（台账 A320／A321）"
else
    bad "rc moved on a SHR reading (pristine=${rcp3} fixed=${rcf3}) - 变好被当成言行不一＝新造的恒红"
fi
p_tbl=$(grep -c '^# STALE ' "$logdir/part3-stale-pristine.txt" || true)
f_tbl=$(grep -c '^# STALE ' "$logdir/part3-stale-fixed.txt" || true)
if [ "$p_tbl" = 0 ]; then
    say "   --     pristine 那一把只有一行行内注（'注=读数变好了'），单列的表一枚没有"
else
    say "   --     NOTE: pristine already prints a standalone table (${p_tbl} rows) - then AC#8-2's table is not the new thing; check what pinned it"
fi
if [ "$f_tbl" -ge 1 ] && grep -q '^# STALE 腿=G5neg ' "$logdir/part3-stale-fixed.txt"; then
    ok "FIXED CODE: 基线过期单列成一张表并点名了那一枚腿（G5neg）"
else
    bad "FIXED CODE printed no 基线过期 table naming the stale leg - AC#8-2 stays human-eyes-only (expected in the PRE-FIX PASS)"
fi

say ""
say "r2-run-legs-and-subtraction.sh: logs kept (created, never deleted): $logdir/part{1,2,3}-*.txt"
say "r2-run-legs-and-subtraction.sh: gate copies + synthetic tree left in place at $work and $tree"
say "r2-run-legs-and-subtraction.sh: this run wrote NOTHING inside $root (the repo's gate bytes were read, never edited)"
if [ "$fail" = 0 ] && [ "$pre_fix_pass" = 0 ]; then
    say "r2-run-legs-and-subtraction.sh: GREEN - run() legs now carry 枚数, a deleted leg moves the exit code, and a stale baseline is a table rather than a sentence."
    exit 0
fi
say "r2-run-legs-and-subtraction.sh: RED - at least one assertion above failed (PRE-FIX PASS expected: $pre_fix_pass)."
exit 1
