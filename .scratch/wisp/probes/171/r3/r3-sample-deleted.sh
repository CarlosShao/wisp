#!/bin/sh
# r3-sample-deleted.sh - ticket 171 r3, AC#9-2: the front-facing sample for "删样本不删腿"
# (the half of ticket 158/161's symptom line "谁摘掉一枚腿或它的样本" that AC#8-1 never got a
# positive fixture for - 171-r2's own next= item 2, and 171-v2 table §14 next= item 2).
#
# THE THREE SHAPES THIS ANSWERS, per leg, with the leg still declared and still run:
#   D2  two registered samples           -> 实测 = 基线 = 2 枚            (the honest state)
#   D1  one of them deleted, leg intact  -> 实测 1 < 基线 2 = SHR 行 + 单列一张【基线过期】表,
#                                           聚合退码**不动**             (this is the "不响" answer)
#   D0  the last sample deleted          -> 实测 0 while the leg says ring = BAD(空/非空那一味),
#                                           聚合退码**跟着动**           (this is the "响" answer)
# The baseline is raised to 2 by patching the leg's own want_n line in a COPY of the ruler,
# because the tracked declaration is registered against the real tree (G6 = 1 枚) and this
# tree holds two samples. Nothing in the repo is edited by that patch.
#
# WHY A SYNTHETIC TREE AGAIN
# Same reason as r3-trailing-comment.sh: the legs read git grep over a COMMIT of
# internal/**/*.go + cmd/**/*.go, so a sample dropped under .scratch/ is invisible to them.
# Throwaway repo under $TMPDIR, no worktree, no checkout, nothing deleted inside this repo.
#
# HARD EXIT RULES: rc=0 all assertions held / rc=1 an assertion did not hold / rc=2 broke.
set -eu

here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
d=$here
while [ ! -e "$d/.git" ]; do
    p=$(CDPATH= cd -- "$d/.." && pwd)
    [ "$p" = "$d" ] && { echo "r3-sample-deleted.sh: no .git above $here - no verdict" >&2; exit 2; }
    d=$p
done
root=$d

gate_rel=.scratch/wisp/probes/154/gate-clauses.sh
GATE="${GATE:-$root/$gate_rel}"
GATE_OTHER="${GATE_OTHER:-$here/gate-prefix-f78d3cb.sh}"
[ -f "$GATE" ] || { echo "r3-sample-deleted.sh: ruler $GATE not found - no verdict" >&2; exit 2; }

logdir="$here/logs"
mkdir -p "$logdir"
work=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3-samples.XXXXXX")
tree=$(mktemp -d "${TMPDIR:-/tmp}/wisp-171-r3-samples-tree.XXXXXX")
export GIT_AUTHOR_NAME=wisp171r3 GIT_AUTHOR_EMAIL=171-r3@invalid
export GIT_COMMITTER_NAME=wisp171r3 GIT_COMMITTER_EMAIL=171-r3@invalid
git -C "$tree" init -q .
mkdir -p "$tree/internal/fake171"

fail=0
say() { printf '%s\n' "$1"; }
ok() { say "  ok    $1"; }
bad() { say "  FAIL  $1"; fail=1; }

sec() {
    awk -v want="$1" '/^## / { if (index($0, want) == 1) ins = 1; else if (ins) ins = 0 } ins { print }' "$2"
}
legrow() { # <leg> <log> -> that leg's verdict line in the aggregate table
    grep -E "^# (ok|SHR|BAD) *腿=$1 " "$2" | tail -1 || true
}
stale_of() { grep -m1 '^# STALE 腿=G6 ' "$1" || true; }
stale_count_of() { grep -c '^# STALE ' "$1" || true; }
stale_count_line() { grep -m1 '^# 基线过期枚数＝' "$1" | sed -e 's/^# 基线过期枚数＝\([0-9]*\).*/\1/'; }
agg_of() { grep -m1 '^# 聚合退码＝' "$1" | sed -e 's/^# 聚合退码＝\([0-9]*\).*/\1/'; }
bad_legs() { # <log> -> the legs the aggregate actually called BAD, one per line
    grep -E '^# BAD .*腿=[^ ]' "$1" | sed -e 's/.*腿=\([^ ]*\).*/\1/' | sort -u || true
}

# ---- a ruler that has TWO registered samples for G6 (want_n 1 -> 2, one line, in a copy) ----
wantline=$(grep -n '^want G6 ring' "$GATE" | cut -d: -f1 | head -1)
[ -n "$wantline" ] || { echo "r3-sample-deleted.sh: no 'want G6' declaration - no verdict" >&2; exit 2; }
nline=$((wantline + 1))
bas=$(sed -n "${nline}p" "$GATE")
case $bas in
"want_n 1") ;;
*) echo "r3-sample-deleted.sh: the line after 'want G6' is '$bas', not 'want_n 1' - the leg moved, no verdict" >&2; exit 2 ;;
esac
gate_two="$work/gate-wantn2.sh"
awk -v ln="$nline" 'NR == ln { print "want_n 2"; next } { print }' "$GATE" >"$gate_two"
if cmp -s "$GATE" "$gate_two"; then
    echo "r3-sample-deleted.sh: raising the baseline changed not one byte - that is not a reading" >&2
    exit 2
fi
n=$(diff "$GATE" "$gate_two" | grep -c '^[<>]')
[ "$n" = 2 ] || { echo "r3-sample-deleted.sh: the baseline patch touched $n lines instead of one" >&2; exit 2; }
say "r3-sample-deleted.sh: ruler with G6 登记两枚样本 = $gate_two (one line changed: want_n 1 -> 2)"
gate_two_other=''
if [ -f "$GATE_OTHER" ] && ! cmp -s "$GATE_OTHER" "$GATE"; then
    wantline=$(grep -n '^want G6 ring' "$GATE_OTHER" | cut -d: -f1 | head -1)
    nline=$((wantline + 1))
    if [ "$(sed -n "${nline}p" "$GATE_OTHER")" = "want_n 1" ]; then
        gate_two_other="$work/gate-wantn2-other.sh"
        awk -v ln="$nline" 'NR == ln { print "want_n 2"; next } { print }' "$GATE_OTHER" >"$gate_two_other"
        say "r3-sample-deleted.sh: the other ruler, same patch = $gate_two_other"
    fi
fi

run_one() { # <ruler> <tag> -> echoes aggregate rc, log at $logdir/samples-<tag>.txt
    log="$logdir/samples-$2.txt"
    rc=0
    (cd "$tree" && sh "$1") >"$log" 2>&1 || rc=$?
    case $rc in
    *[!0-9]*) echo "run_one: exit code '$rc' is not a number - the runner broke" >&2; exit 2 ;;
    esac
    printf '%s' "$rc"
}

state() { # <label> <fixture...> - files listed are kept, everything else in the sample dir is dropped
    label=$1
    shift
    rm -f "$tree"/internal/fake171/*.go
    for f in "$@"; do
        [ -f "$here/fixtures/$f.go.txt" ] || { echo "r3-sample-deleted.sh: fixture $f missing - no verdict" >&2; exit 2; }
        cp "$here/fixtures/$f.go.txt" "$tree/internal/fake171/$f.go"
    done
    git -C "$tree" add -A -- internal/fake171 >/dev/null
    git -C "$tree" commit -q -m "171-r3: samples $label" || { echo "r3-sample-deleted.sh: commit failed" >&2; exit 2; }
    say "-- state '$label' committed at $(git -C "$tree" rev-parse --short HEAD) in $tree ($(ls "$tree"/internal/fake171/*.go 2>/dev/null | wc -l) sample file(s) present)"
}

say ""
say "===== D2: both registered samples present (the honest state) ====="
state D2 d_sample_a d_sample_b
rc_d2=$(run_one "$gate_two" D2)
say "   G6 row : $(legrow G6 "$logdir/samples-D2.txt")"
say "   STALE  : $(stale_of "$logdir/samples-D2.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples-D2.txt") (打印值=$(stale_count_line "$logdir/samples-D2.txt"))  聚合退码=$rc_d2"
say "   BAD 名册：$(bad_legs "$logdir/samples-D2.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples-D2.txt") in
*ok*基线=2枚*实测=2枚*) ok "D2: 实测 equals the registered baseline, no table row, no BAD" ;;
*) bad "D2: expected ok 基线=2枚 实测=2枚, got '$(legrow G6 "$logdir/samples-D2.txt")'" ;;
esac
[ "$(stale_count_of "$logdir/samples-D2.txt")" = 0 ] && ok "D2: the 基线过期 table is empty" || bad "D2: the table is not empty in the honest state"
bad_legs "$logdir/samples-D2.txt" | grep -qxF G6 && bad "D2: G6 counted BAD while it matched its own baseline" || ok "D2: G6 is not in the BAD list"
say ""

say "===== D1: ONE sample deleted, the leg and its declaration untouched ====="
state D1 d_sample_a
rc_d1=$(run_one "$gate_two" D1)
say "   G6 row : $(legrow G6 "$logdir/samples-D1.txt")"
say "   STALE  : $(stale_of "$logdir/samples-D1.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples-D1.txt") (打印值=$(stale_count_line "$logdir/samples-D1.txt"))  聚合退码=$rc_d1"
say "   BAD 名册：$(bad_legs "$logdir/samples-D1.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples-D1.txt") in
*SHR*基线=2枚*实测=1枚*) ok "D1: the leg is flagged SHR (measured below baseline) - the deletion is readable" ;;
*) bad "D1: expected a SHR row for G6, got '$(legrow G6 "$logdir/samples-D1.txt")'" ;;
esac
[ -n "$(stale_of "$logdir/samples-D1.txt")" ] && ok "D1: the standalone 基线过期 table names G6" || bad "D1: no table row - the deletion is human-eyes-only"
if bad_legs "$logdir/samples-D1.txt" | grep -qxF G6; then
    bad "D1: G6 was counted into the exit code - SHR is not supposed to reach it"
else
    say "  ANSWER  D1 未修码上**不响**：腿 G6 一枚样本被摘掉（2 枚→1 枚）之后不在 BAD 名册里，"
    say "          退码里不含它那一味（$rc_d2 -> $rc_d1 的差额全部来自共用同一条射程的别腿，"
    say "          逐枚见上面两张 BAD 名册：$(bad_legs "$logdir/samples-D2.txt" | tr '\n' ' ') => $(bad_legs "$logdir/samples-D1.txt" | tr '\n' ' ')）；"
    say "          它只在【基线过期】那一格里看得见（表内 $(stale_count_of "$logdir/samples-D1.txt") 枚）。"
fi
say ""

say "===== D0: the LAST sample deleted (leg still declared, still run) ====="
state D0
rc_d0=$(run_one "$gate_two" D0)
say "   G6 row : $(legrow G6 "$logdir/samples-D0.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples-D0.txt") (打印值=$(stale_count_line "$logdir/samples-D0.txt"))  聚合退码=$rc_d0"
say "   BAD 名册：$(bad_legs "$logdir/samples-D0.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples-D0.txt") in
*BAD*因=空/非空那一味*) ok "D0: the emptied leg is a BAD on the 空/非空 dimension (票 161 AC#6①那一味)" ;;
*) bad "D0: expected a BAD row naming 空/非空 for G6, got '$(legrow G6 "$logdir/samples-D0.txt")'" ;;
esac
if bad_legs "$logdir/samples-D0.txt" | grep -qxF G6; then
    say "  ANSWER  D0 未修码上**响**：腿 G6 的样本被删空（腿还在跑、声明还在）⇒ 它进了 BAD 名册，"
    say "          聚合退码 $rc_d1 -> $rc_d0（差额的逐枚归属见上面两张名册，本态里共用射程的别腿也一起掉下来）。"
else
    bad "D0: the emptied leg never reached the BAD list - AC#8-1's other half is hollow"
fi
say ""

if [ -n "$gate_two_other" ]; then
    say "===== the same three states against the OTHER ruler (cross-check) ====="
    state X2 d_sample_a d_sample_b
    x2=$(run_one "$gate_two_other" X2)
    state X1 d_sample_a
    x1=$(run_one "$gate_two_other" X1)
    state X0
    x0=$(run_one "$gate_two_other" X0)
    say "   other ruler: D2 rc=$x2  D1 rc=$x1  D0 rc=$x0"
    say "   G6 rows: $(legrow G6 "$logdir/samples-X2.txt")"
    say "            $(legrow G6 "$logdir/samples-X1.txt")"
    say "            $(legrow G6 "$logdir/samples-X0.txt")"
    state D0
    say ""
fi

say "r3-sample-deleted.sh: what this answers for AC#9-2 - 摘掉一枚样本（退码不动、单列表 1 枚）与"
say "  摘光一条腿的样本（退码跟着动）两形都有读数；台件＝probes/171/r3/fixtures/d_sample_*.go.txt，"
say "  读数＝probes/171/r3/logs/samples-D[012].txt。"
say ""
if [ "$fail" = 0 ]; then
    say "r3-sample-deleted.sh: GREEN."
    exit 0
fi
say "r3-sample-deleted.sh: RED."
exit 1
