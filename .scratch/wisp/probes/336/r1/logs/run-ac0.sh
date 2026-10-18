#!/usr/bin/env bash
# 336-r1 AC#0 起手现量：把票面「现量」四发（M-V2a／M-V2b／M-V2c／M-V3）串行重放。
# 只打仓外导出树的物理字节（⛔ 改跟踪文件、⛔ -overlay）；逐发「从具名提交重取字节覆盖＋断言哈希」，
# 跑前跑后各验一次 pristine 哈希；两发串行＝同一个 bash 进程内顺序执行，⛔ 依赖 trap/finally。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ANCHOR="${ANCHOR:?给出具名提交锚}"
TREE="${TREE:?给出仓外导出台根}"
LOGDIR="${LOGDIR:?给出日志目录}"
F="$TREE/internal/risk/blacklist.go"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/blacklist.go")
GOTMPROOT="${GOTMPROOT:?给出 D 卷 GOTMPDIR 根}"
export GOPATH='D:\work\base\gopath'
export PATH="$TREE/third_party/sherpa-onnx:$TREE/build:$PATH"
# 名册＝票 334 那三枚在册仪器（＝334-v1 读的 10/3/0/0 那把），本腿⛔ 改它一字
ROSTER_LEG='TestAnchorMissReasonNamesCandidateAndMissingAnchor|TestAnchorMissKeepsNonOverridableClassAndNeedL2|TestFileBaseNameClaimChangesNoClassOrReason'

restore () { git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go" > "$F"; }
hash_of () { git -C "$REPO" hash-object "$1"; }
check_hash () { local h; h=$(hash_of "$F")
  echo "  hash=$h expected=$PRISTINE match=$([ "$h" = "$PRISTINE" ] && echo YES || echo NO)  [$1]"; }

run_leg () { # $1=日志名
  local gt="$GOTMPROOT/$1"; mkdir -p "$gt" || { echo "  mkdir_gotmp_rc=NONZERO"; return 1; }
  ( cd "$TREE" && GOTMPDIR="$gt" go test ./internal/risk/ -run "$ROSTER_LEG" \
      -count=1 -timeout 30m -v ) > "$LOGDIR/$1" 2>&1
  local rc=$?
  echo "  go_test_rc=$rc RUN=$(grep -c '^=== RUN' "$LOGDIR/$1" || true) PASS_TOP=$(grep -c '^--- PASS' "$LOGDIR/$1" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$LOGDIR/$1" || true) SKIP=$(grep -c '^--- SKIP' "$LOGDIR/$1" || true)"
  grep -E '^\s*--- FAIL: ' "$LOGDIR/$1" | sed 's/([0-9.]*s)//;s/^/  RED /' || true
}

mutate () { # $1=perl -0pi 表达式（内容锚，⛔ 绝对行号）
  restore
  perl -0pi -e "$1" "$F"
  echo "  落点逐字（⇄ 具名提交 blob 的行差）:"
  diff <(git -C "$REPO" show "$ANCHOR:internal/risk/blacklist.go") "$F" | sed 's/^/    /'
}

echo "### AC#0 起手：anchor=$ANCHOR pristine=$PRISTINE tree=$TREE"
restore; check_hash "restored-from-commit"
run_leg ac0-00-baseline-leg.log

echo "### M-V2a 只摘 %LOCALAPPDATA%\Microsoft\Credentials 那一枚 case 的回落支"
mutate 's/case localAppDataS != "" && aUnder\(cand, anchorForms\(anchorPath\(localAppDataS, `microsoft`, `credentials`\)\)\):/case localAppDataS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(localAppDataS, `microsoft`, `credentials`)), isUnder):/'
run_leg ac0-01-MV2a-leg.log
restore; check_hash after-M-V2a

echo "### M-V2b 只摘 ~/.ssh/** 那一枚 case 的回落支"
mutate 's/case homeS != "" && aUnder\(cand, anchorForms\(anchorPath\(homeS, `\.ssh`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.ssh`)), isUnder):/'
run_leg ac0-02-MV2b-leg.log
restore; check_hash after-M-V2b

echo "### M-V2c 只摘 ~/.kube/config 那一枚 case 的回落支"
mutate 's/case homeS != "" && aEq\(cand, anchorForms\(anchorPath\(homeS, `\.kube`, `config`\)\)\):/case homeS != "" \&\& anchorMatches(cand, anchorForms(anchorPath(homeS, `.kube`, `config`)), eqPath):/'
run_leg ac0-03-MV2c-leg.log
restore; check_hash after-M-V2c

echo "### M-V3 反形（变严）：回落支去掉祖先关系界＝任何候选都算命中"
mutate 's/\treturn isUnder\(cand\.root, anchor\.root\) \|\| isUnder\(anchor\.root, cand\.root\)/\treturn true/'
run_leg ac0-04-MV3-leg.log
restore; check_hash after-M-V3

echo "### 末态"
check_hash final-blacklist
echo "  shared_worktree_untouched=$(git -C "$REPO" status --porcelain --untracked-files=no -- internal cmd | wc -l) 行（必须 0）"
echo "ac0_runner_rc=0"
