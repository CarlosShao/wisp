#!/bin/sh
# r3-sample-deleted.sh - ticket 171 r3, AC#9-2: the front-facing sample for "删样本不删腿"
# (the half of ticket 158/161's symptom line "谁摘掉一枚腿或它的样本" that AC#8-1 never got a
# positive fixture for - 171-r2's own next= item 2, and 171-v2 table §14 next= item 2).
#
# ⚠ 票 178（31b5cd36，"四堆不分"修复）之后本钉按 **piles 口径**重推（171s，10-14）：一条 piles 味
#   腿的『实测』＝【② 只开】那一堆的**文件**枚数；【③ 只关】逐枚点名但不进判据（它关的可以是别处
#   开的那一枚）。旧版翻的是一行 `want_n 1 -> 2`＋`want G6 ring` 那个"混数"（正＋反）口径的形状，
#   该形状在 178 后不存在 ⇒ 旧版在钉处 rc=2（逐字 `no 'want G6' declaration - no verdict`）。
#   重推的做法：在**副本**里把 G6 那一腿的登记补成合成树里说实话的两行——
#     want G6 quiet piles -> want G6 ring piles  （真树诚实登记才是 quiet+0；这棵树有两枚只开样本）
#     want_n 0            -> want_n 2            （基线＝② 堆文件枚数＝2）
#   仓里那把尺一字节没动（这一发只在 $work 的副本上翻）。旧读数 logs/samples-D*.txt 原样保留，
#   重推后的读数落 logs/samples2$tag-D*.txt（SAMPLE_PIN_TAG 可加后缀）。
#
# THE FOUR SHAPES THIS ANSWERS, per leg, with the leg still declared and still run:
#   D2  two registered opening samples    -> 实测 = 基线 = 2 枚           (the honest state)
#   D1  one of them deleted, leg intact   -> 实测 1 < 基线 2 = SHR 行 + 单列一张【基线过期】表,
#                                           聚合退码**不动**             (this is the "不响" answer)
#   D1c one opening + one CLOSE-ONLY file -> ② 堆 1 枚／③ 堆 1 枚 ⇒ 实测 = 1 枚（SHR）＝
#                                           ③ 只关没被计进那一列（票 178 的新口径本体）
#   D0  the last sample deleted           -> 实测 0 while the leg says ring = BAD(空/非空那一味),
#                                           聚合退码**跟着动**           (this is the "响" answer)
# 牙（自造反向变异，只在副本上）：`SAMPLE_PIN_MUT=noempty sh 本脚本` 把『空/非空』那一味从副本里
#   挪走 ⇒ 摘光样本不再 BAD ⇒ 本钉转红并逐字点名 D0；不带该变量＝绿。两态读数分开落盘。
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
tag="${SAMPLE_PIN_TAG:-}"
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

# ---- a ruler that tells the truth about THIS tree's G6 leg (two lines, in a copy) ----
# 票 178 piles 口径：合成树里两枚样本各是一枚【只开】⇒ ② 堆＝2 枚文件，诚实登记＝ring＋基线 2。
# 真树的登记是 quiet＋want_n 0（那条射程两堆都空）＝这一发要翻的两行；翻在哪一行先现读、不记号。
wantline=$(grep -n '^want G6 quiet piles$' "$GATE" | cut -d: -f1 | head -1)
[ -n "$wantline" ] || { echo "r3-sample-deleted.sh: no 'want G6 quiet piles' declaration - the leg moved (票 178 重推的形), no verdict" >&2; exit 2; }
nline=$((wantline + 1))
bas=$(sed -n "${nline}p" "$GATE")
case $bas in
"want_n 0") ;;
*) echo "r3-sample-deleted.sh: the line after 'want G6 quiet piles' is '$bas', not 'want_n 0' - the leg moved, no verdict" >&2; exit 2 ;;
esac
gate_two="$work/gate-wantn2.sh"
awk -v wl="$wantline" -v nl="$nline" '
    NR == wl { print "want G6 ring piles"; next }
    NR == nl { print "want_n 2"; next }
    { print }' "$GATE" >"$gate_two"
if cmp -s "$GATE" "$gate_two"; then
    echo "r3-sample-deleted.sh: raising the baseline changed not one byte - that is not a reading" >&2
    exit 2
fi
n=$(diff "$GATE" "$gate_two" | grep -c '^[<>]')
[ "$n" = 4 ] || { echo "r3-sample-deleted.sh: the two-line patch touched $n diff lines instead of four" >&2; exit 2; }
say "r3-sample-deleted.sh: ruler with G6 登记 ring＋基线 2（② 堆文件枚数） = $gate_two"
say "   (two lines changed: want G6 quiet piles -> want G6 ring piles, want_n 0 -> 2)"
if [ -n "${SAMPLE_PIN_MUT:-}" ]; then
    # ---- 反向变异：把被钉的行为挪走（只在 $work 的副本上，仓里零字节） ----
    case $SAMPLE_PIN_MUT in
    noempty)
        sed 's|elif \[ "\$bin" != "\$EXPECT" \]; then|elif false; then|' "$gate_two" >"$gate_two.mut"
        if cmp -s "$gate_two" "$gate_two.mut"; then
            echo "r3-sample-deleted.sh: the 空/非空 mutation changed not one byte - that is not a reading" >&2
            exit 2
        fi
        gate_two="$gate_two.mut"
        say "r3-sample-deleted.sh: MUTATION noempty active - the 空/非空 verdict moved out of the copy, D0 must go RED"
        ;;
    *)
        echo "r3-sample-deleted.sh: unknown SAMPLE_PIN_MUT '$SAMPLE_PIN_MUT' - no verdict" >&2
        exit 2
        ;;
    esac
fi
gate_two_other=''
if [ -f "$GATE_OTHER" ] && ! cmp -s "$GATE_OTHER" "$GATE"; then
    wantline=$(grep -n '^want G6 ring rev$' "$GATE_OTHER" | cut -d: -f1 | head -1)
    nline=$((wantline + 1))
    if [ -n "$wantline" ] && [ "$(sed -n "${nline}p" "$GATE_OTHER")" = "want_n 1" ]; then
        gate_two_other="$work/gate-wantn2-other.sh"
        awk -v ln="$nline" 'NR == ln { print "want_n 2"; next } { print }' "$GATE_OTHER" >"$gate_two_other"
        say "r3-sample-deleted.sh: the other ruler (pre-178 混数口径), want_n 1 -> 2 = $gate_two_other"
    fi
fi

run_one() { # <ruler> <tag> -> echoes aggregate rc, log at $logdir/samples2<tag>-<state>.txt
    log="$logdir/samples2$tag-$2.txt"
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
    if git -C "$tree" diff --cached --quiet -- internal/fake171; then
        say "-- state '$label' is byte-identical to the previous state at $(git -C "$tree" rev-parse --short HEAD) - nothing to commit, reusing it"
        return 0
    fi
    git -C "$tree" commit -q -m "171-r3: samples $label" || { echo "r3-sample-deleted.sh: commit failed" >&2; exit 2; }
    say "-- state '$label' committed at $(git -C "$tree" rev-parse --short HEAD) in $tree ($(ls "$tree"/internal/fake171/*.go 2>/dev/null | wc -l) sample file(s) present)"
}

say ""
say "===== D2: both registered samples present (the honest state) ====="
state D2 d_sample_a d_sample_b
rc_d2=$(run_one "$gate_two" D2)
say "   G6 row : $(legrow G6 "$logdir/samples2$tag-D2.txt")"
say "   STALE  : $(stale_of "$logdir/samples2$tag-D2.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples2$tag-D2.txt") (打印值=$(stale_count_line "$logdir/samples2$tag-D2.txt"))  聚合退码=$rc_d2"
say "   BAD 名册：$(bad_legs "$logdir/samples2$tag-D2.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples2$tag-D2.txt") in
*ok*基线=2枚*实测=2枚*) ok "D2: 实测 equals the registered baseline, no table row, no BAD" ;;
*) bad "D2: expected ok 基线=2枚 实测=2枚, got '$(legrow G6 "$logdir/samples2$tag-D2.txt")'" ;;
esac
[ "$(stale_count_of "$logdir/samples2$tag-D2.txt")" = 0 ] && ok "D2: the 基线过期 table is empty" || bad "D2: the table is not empty in the honest state"
bad_legs "$logdir/samples2$tag-D2.txt" | grep -qxF G6 && bad "D2: G6 counted BAD while it matched its own baseline" || ok "D2: G6 is not in the BAD list"
say ""

say "===== D1: ONE sample deleted, the leg and its declaration untouched ====="
state D1 d_sample_a
rc_d1=$(run_one "$gate_two" D1)
say "   G6 row : $(legrow G6 "$logdir/samples2$tag-D1.txt")"
say "   STALE  : $(stale_of "$logdir/samples2$tag-D1.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples2$tag-D1.txt") (打印值=$(stale_count_line "$logdir/samples2$tag-D1.txt"))  聚合退码=$rc_d1"
say "   BAD 名册：$(bad_legs "$logdir/samples2$tag-D1.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples2$tag-D1.txt") in
*SHR*基线=2枚*实测=1枚*) ok "D1: the leg is flagged SHR (measured below baseline) - the deletion is readable" ;;
*) bad "D1: expected a SHR row for G6, got '$(legrow G6 "$logdir/samples2$tag-D1.txt")'" ;;
esac
[ -n "$(stale_of "$logdir/samples2$tag-D1.txt")" ] && ok "D1: the standalone 基线过期 table names G6" || bad "D1: no table row - the deletion is human-eyes-only"
if bad_legs "$logdir/samples2$tag-D1.txt" | grep -qxF G6; then
    bad "D1: G6 was counted into the exit code - SHR is not supposed to reach it"
else
    say "  ANSWER  D1 不响退码（SHR）：腿 G6 一枚样本被摘掉（② 堆 2 枚→1 枚）之后不在 BAD 名册里，"
    say "          退码里不含它那一味（$rc_d2 -> $rc_d1 的差额全部来自共用同一条射程的别腿，"
    say "          逐枚见上面两张 BAD 名册：$(bad_legs "$logdir/samples2$tag-D2.txt" | tr '\n' ' ') => $(bad_legs "$logdir/samples2$tag-D1.txt" | tr '\n' ' ')）；"
    say "          它只在【基线过期】那一格里看得见（表内 $(stale_count_of "$logdir/samples2$tag-D1.txt") 枚）。"
fi
say ""

say "===== D1c: one opening sample + one CLOSE-ONLY sample (③ 只关 must NOT reach the measured column) ====="
state D1c d_sample_a d_sample_c_close_only
rc_d1c=$(run_one "$gate_two" D1c)
say "   G6 row : $(legrow G6 "$logdir/samples2$tag-D1c.txt")"
say "   piles  : $(sec '## G6 主尺' "$logdir/samples2$tag-D1c.txt" | grep -E '^#   [②③] ' | tr '\n' ' ')"
say "   BAD 名册：$(bad_legs "$logdir/samples2$tag-D1c.txt" | tr '\n' ' ')  聚合退码=$rc_d1c"
case $(legrow G6 "$logdir/samples2$tag-D1c.txt") in
*SHR*基线=2枚*实测=1枚*) ok "D1c: the close-only file did not enter the measured column (票 178：③ 只关不计进这一列)" ;;
*) bad "D1c: expected a SHR row 基线=2枚 实测=1枚 for G6, got '$(legrow G6 "$logdir/samples2$tag-D1c.txt")'" ;;
esac
if sec '## G6 主尺' "$logdir/samples2$tag-D1c.txt" | grep -q '^#   ② 只开＝1 枚文件' &&
    sec '## G6 主尺' "$logdir/samples2$tag-D1c.txt" | grep -q '^#   ③ 只关＝1 枚文件'; then
    ok "D1c: the four-pile printout names ② 只开＝1 枚／③ 只关＝1 枚（实测只吃 ②）"
else
    bad "D1c: expected the G6 leg to print ② 只开＝1 枚文件 and ③ 只关＝1 枚文件"
fi
if bad_legs "$logdir/samples2$tag-D1c.txt" | grep -qxF G6; then
    bad "D1c: G6 was counted into the exit code while ③ 只关 was the only mover"
else
    say "  ANSWER  D1c 不响退码：多出来的那一枚是【③ 只关】（关的是别处开的）——它逐枚点名、不进"
    say "          实测那一列（SHR 表可见），聚合退码里不含它那一味。"
fi
say ""

say "===== D0: the LAST sample deleted (leg still declared, still run) ====="
state D0
rc_d0=$(run_one "$gate_two" D0)
say "   G6 row : $(legrow G6 "$logdir/samples2$tag-D0.txt")"
say "   表内枚数=$(stale_count_of "$logdir/samples2$tag-D0.txt") (打印值=$(stale_count_line "$logdir/samples2$tag-D0.txt"))  聚合退码=$rc_d0"
say "   BAD 名册：$(bad_legs "$logdir/samples2$tag-D0.txt" | tr '\n' ' ')"
case $(legrow G6 "$logdir/samples2$tag-D0.txt") in
*BAD*因=空/非空那一味*) ok "D0: the emptied leg is a BAD on the 空/非空 dimension (票 161 AC#6①那一味)" ;;
*) bad "D0: expected a BAD row naming 空/非空 for G6, got '$(legrow G6 "$logdir/samples2$tag-D0.txt")'" ;;
esac
if bad_legs "$logdir/samples2$tag-D0.txt" | grep -qxF G6; then
    say "  ANSWER  D0 响（BAD 空/非空那一味）：腿 G6 的样本被删空（腿还在跑、声明还在）⇒ 它进了 BAD 名册，"
    say "          聚合退码 $rc_d1 -> $rc_d0（差额的逐枚归属见上面两张名册，本态里共用射程的别腿也一起掉下来）。"
else
    bad "D0: the emptied leg never reached the BAD list - AC#8-1's other half is hollow"
fi
say ""

if [ -n "$gate_two_other" ] && [ -z "${SAMPLE_PIN_MUT:-}" ]; then
    say "===== the same three states against the OTHER ruler (pre-178 混数口径; cross-check only) ====="
    state X2 d_sample_a d_sample_b
    x2=$(run_one "$gate_two_other" X2)
    state X1 d_sample_a
    x1=$(run_one "$gate_two_other" X1)
    state X0
    x0=$(run_one "$gate_two_other" X0)
    say "   other ruler: D2 rc=$x2  D1 rc=$x1  D0 rc=$x0"
    say "   G6 rows: $(legrow G6 "$logdir/samples2$tag-X2.txt")"
    say "            $(legrow G6 "$logdir/samples2$tag-X1.txt")"
    say "            $(legrow G6 "$logdir/samples2$tag-X0.txt")"
    state D0
    say ""
fi

say "r3-sample-deleted.sh: what this answers for AC#9-2（171s 按票 178 piles 口径重推）- 摘掉一枚样本"
say "  （② 堆 2→1：退码不动、单列表 1 枚）／只多一枚【③ 只关】（不进实测）与"
say "  摘光一条腿的样本（② 堆 0、腿还自称 ring ⇒ BAD、退码跟着动）四形都有读数；"
say "  台件＝probes/171/r3/fixtures/d_sample_*.go.txt（+ d_sample_c_close_only.go.txt，171s 新加），"
say "  读数＝probes/171/r3/logs/samples2$tag-D[012c]*.txt（旧混数口径读数 samples-D*.txt 原样保留）。"
say ""
if [ "$fail" = 0 ]; then
    say "r3-sample-deleted.sh: GREEN."
    exit 0
fi
say "r3-sample-deleted.sh: RED."
exit 1
