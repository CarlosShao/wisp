#!/usr/bin/env bash
# 票 216 AC#6 的三发仓外突变台件（腿 216-r2）。全部打在仓外导出树 wisp216r2-mut 上，
# 逐发还原并自证还原哈希；共享工作树与提交树一字不动。
#   M1  rules_shell.go   strings.EqualFold(argv0, entry) -> strings.Contains(argv0, entry)
#      ＝AC#6 甲 点名的变异正控：同一把尺必须回 1
#   M2  blacklist.go     B 档循环里插一发"名字含『已允许』就 ClassNone（放行）"
#      ＝AC#6 乙 点名的牙证：hostile 臂必须红，且ⓐ 那三窗仍 0/0/0（两发落在不同处）
#   R1  rules_shell.go   整枚函数改名 shellAllowlisted -> shellAllowlistedRenamed216
#      ＝编排者必答题：算子钉对"函数名找不到"是红还是静默
set -u
REPO="D:/work/workspace/projects plans/Wisp"
HERE="$(cd "$(dirname "$0")" && pwd)"
ROOT=/d/tmp/leg216r2trees
MUT="$ROOT/wisp216r2-mut"
OUT="$HERE"
SHA_AFTER=671d3585
ANCHOR=72b65c61
GOTMP=/d/tmp/leg216r2gotmp
mkdir -p "$GOTMP"
export GOTMPDIR='D:\tmp\leg216r2gotmp'
export GOPATH='D:\work\base\gopath'
SEL='^TestAllowSideFunctionsCarryNoSubstringOperator$|^TestSubstringOperatorCounter216IsNotBlind$|^TestFileBaseNameClaimChangesNoClassOrReason$|^TestArgvZeroClaimChangesNoShellVerdict$|^TestHostClaimChangesNoNetworkVerdict$|^TestOverrideKeySpelledWithClaimDoesNotUnlock$'
SHELLGAUGE='strings\.(Contains|ContainsAny|ContainsRune|Index|IndexAny)\('

echo "########## 0 建树（突变台面＝已提交态 $SHA_AFTER）"
bash "$HERE/build-trees-216r2.sh" "$SHA_AFTER:wisp216r2-mut"

echo "########## 1 起手红名册作差（改前 $ANCHOR ⇄ 改后 $SHA_AFTER，整包 -v 读数）"
grep -E '^--- (FAIL|SKIP)' "$OUT/baseline-before-72b65c61.log" | sed 's/ (.*//' | sort > "$OUT/roster-before.log"
grep -E '^--- (FAIL|SKIP)' "$OUT/after-671d3585.log" | sed 's/ (.*//' | sort > "$OUT/roster-after.log"
echo "作差 comm -3（应为空）:"; comm -3 "$OUT/roster-before.log" "$OUT/roster-after.log"
echo "comm -3 行数=$(comm -3 "$OUT/roster-before.log" "$OUT/roster-after.log" | wc -l)"

mutate () { # $1=文件 $2=python 片段名
  python - "$MUT/internal/risk/$1" "$2" <<'PY'
import io, sys
path, which = sys.argv[1], sys.argv[2]
data = io.open(path, encoding="utf-8", newline="").read()
if which == "M1":
    old = "if strings.EqualFold(argv0, entry) || strings.EqualFold(base, entry) {"
    new = "if strings.Contains(argv0, entry) || strings.EqualFold(base, entry) { // MUTATION-216R2-M1"
    assert data.count(old) == 1, ("LANDS:0", data.count(old))
    out = data.replace(old, new)
elif which == "M2":
    old = "\t\tbase := strings.ToLower(baseName(s))\n"
    new = ("\t\tbase := strings.ToLower(baseName(s))\n"
           "\t\tif strings.Contains(base, \"\u5df2\u5141\u8bb8\") { // MUTATION-216R2-M2\n"
           "\t\t\treturn ClassNone, \"\"\n"
           "\t\t}\n")
    assert data.count(old) == 1, ("LANDS:0", data.count(old))
    out = data.replace(old, new)
elif which == "R1":
    old, new = "shellAllowlisted", "shellAllowlistedRenamed216"
    assert data.count(old) == 3, ("LANDS:0", data.count(old))
    out = data.replace(old, new)
else:
    raise SystemExit("unknown mutation " + which)
io.open(path, "w", encoding="utf-8", newline="").write(out)
print("LANDS: %s applied to %s, marker count=%d" % (which, path, out.count("MUTATION-216R2") if which != "R1" else 0))
PY
  grep -n "MUTATION-216R2\|Renamed216" "$MUT/internal/risk/$1" | head -4
}

run_sel () { # $1=日志名
  ( cd "$MUT" && go test ./internal/risk/ -count=1 -timeout 30m -v -run "$SEL" \
      > "$OUT/$1.log" 2>&1; echo "rc=$?" >> "$OUT/$1.log" )
  echo "  $1: topPASS=$(grep -c '^--- PASS' "$OUT/$1.log") topFAIL=$(grep -c '^--- FAIL' "$OUT/$1.log") subFAIL=$(grep -c '^    --- FAIL' "$OUT/$1.log") panic=$(grep -c '^panic' "$OUT/$1.log") rc=$(tail -1 "$OUT/$1.log")"
}

restore () { # $1=文件
  cp "$MUT/internal/risk/$1.orig216" "$MUT/internal/risk/$1"
  local got want
  got="$(git -C "$REPO" hash-object "$MUT/internal/risk/$1")"
  want="$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/$1")"
  echo "restore $1: hash-object=$got anchor-blob=$want MATCH=$([ "$got" = "$want" ] && echo yes || echo NO)"
  echo "  MUTATION 残留=$(grep -c 'MUTATION-216R2\|Renamed216' "$MUT/internal/risk/$1")"
}

echo "########## 2 突变台面起手一发（应全绿）"
run_sel muttree-clean-start

for f in rules_shell.go blacklist.go; do cp "$MUT/internal/risk/$f" "$MUT/internal/risk/$f.orig216"; done

echo "########## 3 M1＝rules_shell.go 里 EqualFold(argv0,entry) 换成 Contains"
mutate rules_shell.go M1
bash "$HERE/ac6a-present-state.sh" "$MUT" | tee "$OUT/ac6a-under-M1.log"
run_sel mut-M1
echo "----- M1 红句逐字（甲那把尺）"
grep -nE 'allow_substring_pin_216_test\.go:[0-9]+:|name_claim_two_arm_216_test\.go:[0-9]+:' "$OUT/mut-M1.log" | head -12
restore rules_shell.go

echo "########## 4 M2＝blacklist.go 的 B 档循环插『名字含已允许就放行』"
mutate blacklist.go M2
bash "$HERE/ac6a-present-state.sh" "$MUT" | tee "$OUT/ac6a-under-M2.log"
run_sel mut-M2
echo "----- M2 红句逐字（乙那臂）"
grep -nE 'name_claim_two_arm_216_test\.go:[0-9]+:' "$OUT/mut-M2.log" | head -12
restore blacklist.go

echo "########## 5 R1＝整枚 shellAllowlisted 改名（编排者必答题：红还是静默）"
mutate rules_shell.go R1
run_sel mut-R1-rename
echo "----- R1 红句逐字"
grep -nE 'allow_substring_pin_216_test\.go:[0-9]+:|build failed|undefined' "$OUT/mut-R1-rename.log" | head -8
restore rules_shell.go

echo "########## 6 还原后端面自证：突变台面再跑一发（应回到全绿）"
run_sel muttree-restored-final
echo "残留标记=$(grep -rc 'MUTATION-216R2\|Renamed216' "$MUT/internal/risk/" | grep -v ':0' | wc -l) 个文件带残留"
echo "共享工作树 internal/risk 状态："
git -C "$REPO" status --porcelain -- internal/risk cmd/wisp internal/panel internal/observe
echo "########## done"
