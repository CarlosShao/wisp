#!/usr/bin/env bash
# 216-v3：第三把台面（★不是"C 卷⇄D 卷"，而是"主目录是谁"）。
# 造一枚假主目录 D:\W216V3FHOME（USERPROFILE 指它），把 t.TempDir() 落在它之下，
# 全部在 D 卷上——如果那枚红跟着"主目录祖先链"走而⛔ 跟着卷走，它应当在 D 卷上照样响。
# 第二发再在假主目录里补造 .aws（那枚锚点今天缺位才说出 "~/.aws/credentials"），
# 看 reason 会不会换到下一枚缺位锚 ~/.kube/config ⇒ 直接证"reason 说的是哪一枚锚"由盘上缺位决定。
# ⛔ 动跟踪件 ⛔ 动断言；台面＝仓外突变树 D:\W216V3MUT（起手已还原＝锚上 blob）。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
HERE="$(cd "$(dirname "$0")/.." && pwd)"
TREE=/d/W216V3MUT
FH=/d/W216V3FHOME
SEL='TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock'
counts () { printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s\n' \
  "$(grep -c '^=== RUN' "$1" || true)" "$(grep -c '^--- PASS' "$1" || true)" \
  "$(grep -c '^--- FAIL' "$1" || true)" "$(grep -c '^--- SKIP' "$1" || true)"; }

echo "########## 0 起手：树里两枚被突变件仍等于锚上 blob"
for f in rules_shell.go blacklist.go; do
  echo "$f tree=$(git hash-object "$TREE/internal/risk/$f") blob=$(cd "$REPO" && git rev-parse "7a94c3bd:internal/risk/$f")"
done
echo "########## 1 造假主目录（照真机摆：.ssh 与 .git-credentials 在场＝全展开锚；.aws/.kube 缺位）"
mkdir -p "$FH/tempdir" "$FH/.ssh"
printf "fake home fixture, not a real key
" > "$FH/.git-credentials"
ls -a "$FH" | paste -sd' ' -
echo "########## 2 发一：USERPROFILE=假主目录，TMP=假主目录\tempdir（全在 D 卷）"
( cd "$TREE" && env USERPROFILE='D:\W216V3FHOME' TMP='D:\W216V3FHOME\tempdir' TEMP='D:\W216V3FHOME\tempdir' \
  go test ./internal/risk/ -count=1 -timeout 20m -v -run "$SEL" ) >"$HERE/logs/fakehome-aws-absent.log" 2>&1
echo "rc=$?" >>"$HERE/logs/fakehome-aws-absent.log"
echo "  $(counts "$HERE/logs/fakehome-aws-absent.log") red=$(grep '^--- FAIL' "$HERE/logs/fakehome-aws-absent.log" | sed 's/ (.*//' | paste -sd';' -)"
echo "  红句里的锚 reason: $(grep -o 'A-list sensitive path (non-overridable): [^\\"]*' "$HERE/logs/fakehome-aws-absent.log" | sort -u | paste -sd'; ' -)"
echo "  落点具名: $(grep -o 'D:.W216V3FHOME.tempdir' "$HERE/logs/fakehome-aws-absent.log" | head -1)"
echo "########## 3 发二：同一棵假主目录里补造 .aws（那枚锚改成全展开，reason 该换到下一枚缺位锚）"
mkdir -p "$FH/.aws"
( cd "$TREE" && env USERPROFILE='D:\W216V3FHOME' TMP='D:\W216V3FHOME\tempdir' TEMP='D:\W216V3FHOME\tempdir' \
  go test ./internal/risk/ -count=1 -timeout 20m -v -run "$SEL" ) >"$HERE/logs/fakehome-aws-present.log" 2>&1
echo "rc=$?" >>"$HERE/logs/fakehome-aws-present.log"
echo "  $(counts "$HERE/logs/fakehome-aws-present.log") red=$(grep '^--- FAIL' "$HERE/logs/fakehome-aws-present.log" | sed 's/ (.*//' | paste -sd';' -)"
echo "  红句里的锚 reason: $(grep -o 'A-list sensitive path (non-overridable): [^\\"]*' "$HERE/logs/fakehome-aws-present.log" | sort -u | paste -sd'; ' -)"
echo "########## 4 收尾：跟踪字节未动（应空）、树内无 MUTATION 残留"
( cd "$REPO" && git status --porcelain -- internal/risk )
echo "mutation_residual=$(grep -rc 'MUTATION-216V3' "$TREE/internal/risk/rules_shell.go" "$TREE/internal/risk/blacklist.go" | paste -sd' ' -)"
echo "########## done"
