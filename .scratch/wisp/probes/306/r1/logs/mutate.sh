#!/bin/bash
# 306-r1 mutation driver (lives in the out-of-repo export tree, never in the repo).
# Four cycles: three nails x (mutate -> named case must go red -> cmp restore -> re-green).
# Each cycle also runs a "fixture moved away" whole-package control, which is the ticket's
# "changing the in-line value leaves the package green" half, re-run on this same bench.
set -u
TREE="$1"
cd "$TREE" || exit 9

cycle() {
  local name="$1" sedexpr="$2" lno="$3"
  cp pristine-wavinjector.go internal/audio/wavinjector.go
  sed -i "$sedexpr" internal/audio/wavinjector.go
  {
    echo "################ CYCLE $name ################"
    echo "--- mutated line $lno verbatim ---"
    sed -n "${lno}p" internal/audio/wavinjector.go
    echo "--- pristine line $lno verbatim (for the byte-for-byte pair) ---"
    sed -n "${lno}p" pristine-wavinjector.go
    echo "--- whole-file diff pristine vs mutated (must be exactly one hunk) ---"
    diff -u pristine-wavinjector.go internal/audio/wavinjector.go | sed -n '/^@@/,$p'
    echo "rc_diff=$?"
  } > "mut-$name-01-mutated.txt" 2>&1

  # 改前台面：把那枚新用例 mv 走（不删），产码已带突变
  mv internal/audio/wavinjector_extensible_float_306_test.go internal/audio/_306_fixture_moved_away.go.bak
  ls internal/audio/_306_fixture_moved_away.go.bak > /dev/null || echo "MV FAILED"
  GOFLAGS= go test ./internal/audio/ -count=1 -v > "mut-$name-02-nofixture-fullpkg.txt" 2>&1
  local rc_a=$?
  mv internal/audio/_306_fixture_moved_away.go.bak internal/audio/wavinjector_extensible_float_306_test.go
  {
    echo "rc_pkg_no_fixture=$rc_a"
    echo "top_fail=$(grep -cE '^--- FAIL: ' mut-$name-02-nofixture-fullpkg.txt)"
    echo "top_pass=$(grep -cE '^--- PASS: ' mut-$name-02-nofixture-fullpkg.txt)"
    echo "top_skip=$(grep -cE '^--- SKIP: ' mut-$name-02-nofixture-fullpkg.txt)"
    echo "sub_fail=$(grep -cE '^[[:space:]]+--- FAIL: ' mut-$name-02-nofixture-fullpkg.txt)"
    echo "tail:"
    tail -2 mut-$name-02-nofixture-fullpkg.txt
    echo "any mention of the mutated values in the FAIL roster (expect none):"
    grep -E '^--- FAIL' mut-$name-02-nofixture-fullpkg.txt || echo "NONE"
  } >> "mut-$name-02-nofixture-fullpkg.txt" 2>&1

  # 改后台面：夹具在场，定向那两枚指名用例
  GOFLAGS= go test ./internal/audio/ -count=1 -v -run 306 > "mut-$name-03-red.txt" 2>&1
  local rc_b=$?
  {
    echo "rc_named_cases_with_mutation=$rc_b"
    echo "--- verbatim red lines (top-level) ---"
    grep -E '^--- FAIL: ' "mut-$name-03-red.txt" || echo "NO TOP-LEVEL FAIL LINE"
    echo "--- verbatim failure text ---"
    grep -vE '^(=== RUN|--- PASS|PASS$|ok )' "mut-$name-03-red.txt"
  } >> "mut-$name-03-red.txt" 2>&1

  # cmp 逐字节还原
  cp pristine-wavinjector.go internal/audio/wavinjector.go
  cmp internal/audio/wavinjector.go pristine-wavinjector.go > "mut-$name-04-cmp.txt" 2>&1
  local rc_c=$?
  echo "rc_cmp_restore=$rc_c (0 = byte-for-byte identical to the pristine copy)" >> "mut-$name-04-cmp.txt"

  # 复绿
  GOFLAGS= go test ./internal/audio/ -count=1 -v -run 306 > "mut-$name-05-regreen.txt" 2>&1
  local rc_d=$?
  {
    echo "rc_regreen=$rc_d"
    grep -E '^--- (PASS|FAIL): ' "mut-$name-05-regreen.txt"
    tail -2 "mut-$name-05-regreen.txt"
  } >> "mut-$name-05-regreen.txt" 2>&1

  echo "CYCLE $name rc_no_fixture_fullpkg=$rc_a rc_red=$rc_b rc_cmp=$rc_c rc_regreen=$rc_d"
}

cycle jia-192-fffe-to-fffd  '192s/0xFFFE/0xFFFD/' 192
cycle yi-218-three-to-two   '218s/fmtTag == 3 /fmtTag == 2 /' 218
cycle bing-196-off24-to-26  '196s/data\[body+24 : body+26\]/data[body+26 : body+28]/' 196
cycle ding-196-literal-wording '196s/data\[body+24 :/data[body+26 :/' 196

echo "--- final safety: production file back to pristine, and identical to the source repo HEAD blob ---"
cd "$TREE" || exit 9
cp pristine-wavinjector.go internal/audio/wavinjector.go
cmp internal/audio/wavinjector.go pristine-wavinjector.go; echo "rc_final_cmp_pristine=$?"
cmp internal/audio/wavinjector.go "D:/work/workspace/projects plans/Wisp/internal/audio/wavinjector.go"; echo "rc_final_cmp_vs_real_worktree=$?"
