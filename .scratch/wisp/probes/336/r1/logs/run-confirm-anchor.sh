#!/usr/bin/env bash
# 336-r1 交件锚确认发：在【交件侧 HEAD 导出的树】上，同一枚名册跑三发串行——
# ① pristine（成对正控＝必绿）② 只摘 MU-1h 那一枚 :237 的回落支（⇒ 必红，红名册逐名）③ 还原再跑（必回绿）。
# ★首版我在 -c 的双引号里写 perl 表达式⇒反引号被 shell 当命令替换吃掉⇒落点命中=0＝那发的 rc=0 是【假读数】（台面故障，
#   已具名作废，见 logs/96-confirm-anchor.log 末节）；本版＝表达式落在本件单引号里，⛔ 依赖 echo。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?具名提交锚}"
TREE="${TREE:?仓外导出台根}"
LOGDIR="$REPO/.scratch/wisp/probes/336/r1/logs"
GOTMP="${GOTMP:-/d/tmp/wisp/20261018/336r1/gotmp-confirm-anchor}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
ROSTER='TestAnchorFallbackTeethPerAnchor336|TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"

shot () { # $1=日志名
  mkdir -p "$GOTMP/$1" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOTMPDIR="$GOTMP/$1" go test ./internal/risk/ -run "$ROSTER" -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^ *--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
}

echo "### 交件锚确认发 anchor=$ANCHOR tree=$TREE pristine=$PRISTINE"
git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
echo "  ① 起手哈希=$(git -C "$REPO" hash-object "$F") match=$([ "$(git -C "$REPO" hash-object "$F")" = "$PRISTINE" ] && echo YES || echo NO)"
echo "  teeth_on_tree=$(git -C "$REPO" hash-object "$TREE/internal/risk/anchor_fallback_teeth_336_windows_test.go")"
echo "  ① pristine 读数（成对正控＝必绿）:"; shot cf-01-pristine.log

git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
perl -0pi -e 's/case localAppDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(localAppDataS, `microsoft`, `credentials`\)\)\):/case localAppDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder):/' "$F"
echo "  落点逐字:"; diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
echo "  ② 落点命中枚数=$(grep -c 'anchorMatches(cand, anchorForms' "$F")（期望 1）hash=$(git -C "$REPO" hash-object "$F")"
echo "  ② 摘 MU-1h 那一枚的读数（⇒必红）:"; shot cf-02-MU1h.log

git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"
echo "  ③ 还原哈希=$(git -C "$REPO" hash-object "$F") match=$([ "$(git -C "$REPO" hash-object "$F")" = "$PRISTINE" ] && echo YES || echo NO)"
echo "  ③ 还原后复跑（必回绿）:"; shot cf-03-restored.log
echo "  末态哈希=$(git -C "$REPO" hash-object "$F")"
echo "  共享树 internal/cmd 脏行=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)"
echo "confirm_anchor_rc=0"
