#!/usr/bin/env bash
# 336-r1 AC#2 去牙反控扫描：对名册里每一枚 aEq/aUnder 操作数，只摘「这一枚」的回落支
# （把它那一处 aEq/aUnder 换成不带回落支的 anchorMatches），⛔ 一次摘多枚。
# 只打仓外导出树物理字节；逐发「从具名提交重取字节覆盖＋跑前跑后各验一次 pristine 哈希」。
# 成对正控＝同一发在未突变树上必绿（teeth-01-pristine-dshape.log）。
# 两把尺并行：TEETH＝本腿新牙；LEG＝票 334 那三枚在册仪器（差集＝本牙补上的射程）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?具名提交锚}"
TREE="${TREE:?仓外导出台根}"
LOGDIR="$REPO/.scratch/wisp/probes/336/r1/logs"
GOTMPROOT="${GOTMPROOT:-/d/tmp/wisp/20261018/336r1/gotmp}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
ROSTER_TEETH='TestAnchorFallbackTeethPerAnchor336'
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
check_hash () { local h; h=$(git -C "$REPO" hash-object "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }
run () { # $1=名册 $2=日志名
  local gt="$GOTMPROOT/$2"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$1" -count=1 -timeout 30m -v ) > "$LOGDIR/$2" 2>&1
  local rc=$?
  echo "  rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$2" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$2" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$2" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$2" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$2" | sed 's/([0-9.]*s)//;s/^/    RED /' || true
}
mutate () { # $1=perl -0pi（内容锚，⛔ 行号）
  restore; perl -0pi -e "$1" "$F"
  echo "  落点逐字（⇄ 具名提交 blob）:"; diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
  echo "  落点命中枚数=$(grep -c 'anchorMatches(cand, anchorForms' "$F")（本发＝摘一枚，期望 1；摘两枚时期望 2）"
}

echo "### AC#2 扫描 anchor=$ANCHOR pristine=$PRISTINE tree=$TREE"
restore; check_hash start
echo "### 基线（未突变＝成对正控）"
echo "  TEETH:"; run "$ROSTER_TEETH" ac2-00-baseline-teeth.log
echo "  LEG:"; run "$ROSTER_LEG" ac2-00-baseline-leg.log

# 每一枚：把那一处 case 的 aEq/aUnder 换成 anchorMatches（同关系、同锚拼法，只去回落支）
declare -a MU=(
  'MU-1a git-credentials aEq|s/&& aEq\(cand, anchorForms\(anchorPath\(homeS, `\.git-credentials`\)\)\):/\&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.git-credentials`)), eqPath):/'
  'MU-1b wisp-config-under aUnder|s/case aUnder\(cand, anchorForms\(anchorPath\(appDataS, `wisp`\)\)\) &&/case anchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`)), isUnder) \&\&/'
  'MU-1c wisp-config-eq aEq|s/\t\taEq\(cand, anchorForms\(anchorPath\(appDataS, `wisp`, `config\.toml`\)\)\):/\t\tanchorMatches(cand, anchorForms(anchorPath(appDataS, `wisp`, `config.toml`)), eqPath):/'
  'MU-1d ssh aUnder|s/case homeS != "" && aUnder\(cand, anchorForms\(anchorPath\(homeS, `\.ssh`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder):/'
  'MU-1e aws aEq|s/case homeS != "" && aEq\(cand, anchorForms\(anchorPath\(homeS, `\.aws`, `credentials`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.aws`, `credentials`)), eqPath):/'
  'MU-1f kube aEq|s/case homeS != "" && aEq\(cand, anchorForms\(anchorPath\(homeS, `\.kube`, `config`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)), eqPath):/'
  'MU-1g dpapi-protect aUnder|s/case appDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(appDataS, `microsoft`, `protect`\)\)\):/case appDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(appDataS, `microsoft`, `protect`)), isUnder):/'
  'MU-1h local-credentials aUnder|s/case localAppDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(localAppDataS, `microsoft`, `credentials`\)\)\):/case localAppDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder):/'
)
for entry in "${MU[@]}"; do
  tag="${entry%%|*}"; expr="${entry#*|}"
  echo "### $tag"
  mutate "$expr"
  echo "  TEETH:"; run "$ROSTER_TEETH" "ac2-${tag%% *}-teeth.log"
  echo "  LEG:"; run "$ROSTER_LEG" "ac2-${tag%% *}-leg.log"
  restore; check_hash "after-${tag%% *}"
done
echo "### 末态"
check_hash final
echo "  shared_worktree_internal_cmd_dirty=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l)（期望 0）"
echo "ac2_sweep_rc=0"
