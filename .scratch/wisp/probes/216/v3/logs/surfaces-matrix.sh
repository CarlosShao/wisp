#!/usr/bin/env bash
# 216-v3 台面矩阵：同一枚命令、同一份码，只改 t.TempDir() 的落点，看那枚红跟的是
# "C 卷" 还是 "用户主目录的祖先链"。判据＝票 216 AC#6 那两枚仪器（⛔ 本腿不动任何断言）。
# 尺（颜色尺四数，四数不许混用、不许相减）：
#   RUN=      grep -c '^=== RUN'   （含子测试）
#   PASS_TOP= grep -c '^--- PASS'  （只顶层）
#   FAIL_TOP= grep -c '^--- FAIL'  （只顶层）
#   SKIP=     grep -c '^--- SKIP'
# 用法：bash surfaces-matrix.sh
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
LOGDIR=.scratch/wisp/probes/216/v3/logs
RISK216='TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock'

mkdir -p 'C:/Users/swq/wisp216v3-home' 'C:/Users/Public/wisp216v3-pub' 'D:/wisp-216v3-tmp'

counts () { # $1=日志
  printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s panic=%s rc=%s\n' \
    "$(grep -c '^=== RUN' "$1" || true)" "$(grep -c '^--- PASS' "$1" || true)" \
    "$(grep -c '^--- FAIL' "$1" || true)" "$(grep -c '^--- SKIP' "$1" || true)" \
    "$(grep -c '^panic' "$1" || true)" "$(grep -o 'rc=[0-9]*' "$1" | tail -1)"
}

surface () { # $1=名 $2=TMP 值（空＝默认台面）
  local name="$1" tmpv="$2" full="$LOGDIR/surface-$1-full.log" scoped="$LOGDIR/surface-$1-216only.log"
  if [ -z "$tmpv" ]; then
    go test ./internal/risk/ -count=1 -timeout 20m -v          >"$full"   2>"$full.err";   echo "rc=$?" >>"$full"
    go test ./internal/risk/ -count=1 -timeout 20m -v -run "$RISK216" >"$scoped" 2>"$scoped.err"; echo "rc=$?" >>"$scoped"
  else
    env TMP="$tmpv" TEMP="$tmpv" go test ./internal/risk/ -count=1 -timeout 20m -v          >"$full"   2>"$full.err";   echo "rc=$?" >>"$full"
    env TMP="$tmpv" TEMP="$tmpv" go test ./internal/risk/ -count=1 -timeout 20m -v -run "$RISK216" >"$scoped" 2>"$scoped.err"; echo "rc=$?" >>"$scoped"
  fi
  printf '%s\n' "### surface=$name TMP=${tmpv:-<默认 /tmp -> C:\\Users\\swq\\AppData\\Local\\Temp>}"
  printf '  full   %s  tempdir_seen=%s\n'   "$(counts "$full")"   "$(grep -o '^\s*[a-z_]*two_arm_216_test.go:10[0-9]: plain "[^"\\]*\\[^"\\]*' "$full" | head -1 | sed 's/.*plain "//')"
  printf '  216only %s\n' "$(counts "$scoped")"
  printf '  red_roster_full=%s\n'   "$(grep '^--- FAIL' "$full"   | tr -s ' ' | paste -sd' ; ' -)"
  printf '  red_roster_216=%s\n'    "$(grep '^--- FAIL' "$scoped" | tr -s ' ' | paste -sd' ; ' -)"
}

surface C-default          ''
surface D-redirect         'D:\wisp-216v3-tmp'
surface C-home-nonLAD      'C:\Users\swq\wisp216v3-home'
surface C-outside-home     'C:\Users\Public\wisp216v3-pub'
echo "### done"
