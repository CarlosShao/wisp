#!/usr/bin/env bash
# 票 216 AC#6 非实现者终裁腿 216-v3 的突变重放台件（裁＝两枚仪器有没有牙）。
# 三发全部打在【仓外导出树】的物理字节上（⛔ 改跟踪件、⛔ -overlay）：
#   M1  internal/risk/rules_shell.go  strings.EqualFold(argv0, entry) -> strings.Contains(...)
#   M2  internal/risk/blacklist.go    B 档循环里插"base 含 已允许 就 return ClassNone"
#   R1  internal/risk/rules_shell.go  shellAllowlisted 整枚改名（声明＋调用＋注释＝3 处同改）
# 落点全部用【内容锚】（函数体内那枚 token），⛔ 绝对行号。
# 每发：LANDS 自证 -> 当场跑读数 -> finally 还原 -> 还原哈希对回锚上 blob。
# 台面＝TMP/TEMP 指 D 卷（t.TempDir() 落主目录祖先链之外＝那枚 C 卷特红不在场，
# 突变红才可归因到突变本身）；另附一发 M2 打在 C 卷默认台面上，看那枚特红会不会冒充它。
set -u
REPO="D:/work/workspace/projects plans/Wisp"
HERE="$(cd "$(dirname "$0")/.." && pwd)"  # = .scratch/wisp/probes/216/v3（读数写进 $HERE/logs）
ANCHOR=7a94c3bd
ROOT=/d/W216V3MUT
GOTMP=/d/W216V3GOTMP   # 只建不导出：导出了 go test 会把 t.TempDir() 挪走（216-r2 那两发就是这么落在 D 卷的）
mkdir -p "$GOTMP" /d/W216V3TMP "$HERE/logs"
echo "GOTMPDIR_in_env=[${GOTMPDIR:-<unset>}]"
export GOPATH='D:\work\base\gopath'
EIGHT="go.mod go.sum internal cmd tools scripts frontend models"
SEL216='TestAllowSideFunctionsCarryNoSubstringOperator|TestSubstringOperatorCounter216IsNotBlind|TestFileBaseNameClaimChangesNoClassOrReason|TestArgvZeroClaimChangesNoShellVerdict|TestHostClaimChangesNoNetworkVerdict|TestOverrideKeySpelledWithClaimDoesNotUnlock'

counts () { printf 'RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s\n' \
  "$(grep -c '^=== RUN' "$1" || true)" "$(grep -c '^--- PASS' "$1" || true)" \
  "$(grep -c '^--- FAIL' "$1" || true)" "$(grep -c '^--- SKIP' "$1" || true)"; }

echo "########## 0 建树（锚 $ANCHOR，八件套＋仓根 docs/＋手工三件，落 D 卷）"
rm -rf "$ROOT"; mkdir -p "$ROOT"
( cd "$REPO" && git archive "$ANCHOR" $EIGHT ) | tar -x -C "$ROOT"
( cd "$REPO" && git archive "$ANCHOR" docs ) | tar -x -C "$ROOT"
mkdir -p "$ROOT/build" "$ROOT/third_party/sherpa-onnx" "$ROOT/frontend/dist"
cp "$REPO"/build/*.dll "$ROOT/build/" 2>/dev/null
cp "$REPO"/third_party/sherpa-onnx/*.dll "$ROOT/third_party/sherpa-onnx/" 2>/dev/null
cp -r "$REPO"/frontend/dist/. "$ROOT/frontend/dist/" 2>/dev/null
echo "go_files=$(find "$ROOT/internal" "$ROOT/cmd" "$ROOT/tools" -name '*.go' | wc -l) dll_build=$(ls "$ROOT/build"/*.dll 2>/dev/null | wc -l) dll_sherpa=$(ls "$ROOT/third_party/sherpa-onnx"/*.dll 2>/dev/null | wc -l) dist=$(ls "$ROOT/frontend/dist" 2>/dev/null | wc -l) docs=$([ -d "$ROOT/docs" ] && echo yes || echo no)"
( cd "$ROOT" && go build ./internal/risk/; echo "go_build_rc=$?" )

echo "########## 1 起手哈希（导出树字节 ⇄ 锚上 blob；不同＝建树带过转换，还原基准按导出树那枚记）"
for f in rules_shell.go blacklist.go; do
  echo "$f tree=$(git hash-object "$ROOT/internal/risk/$f") blob=$(cd "$REPO" && git rev-parse "$ANCHOR:internal/risk/$f")"
done

mutate () { # $1=M1|M2|R1 $2=文件名 -> 打印 LANDS
  python - "$ROOT/internal/risk/$2" "$1" <<'PY'
import io, sys
path, which = sys.argv[1], sys.argv[2]
data = io.open(path, encoding="utf-8", newline="").read()
if which == "M1":
    old = "if strings.EqualFold(argv0, entry) || strings.EqualFold(base, entry) {"
    new = "if strings.Contains(argv0, entry) || strings.EqualFold(base, entry) { // MUTATION-216V3-M1"
    n = data.count(old); out = data.replace(old, new)
elif which == "M2":
    old = "\t\tbase := strings.ToLower(baseName(s))\n"
    new = ("\t\tbase := strings.ToLower(baseName(s))\n"
           "\t\tif strings.Contains(base, \"\u5df2\u5141\u8bb8\") { // MUTATION-216V3-M2\n"
           "\t\t\treturn ClassNone, \"\"\n\t\t}\n")
    n = data.count(old); out = data.replace(old, new)
elif which == "R1":
    old, new = "shellAllowlisted", "shellAllowlistedRenamed216V3"
    n = data.count(old); out = data.replace(old, new)
else:
    raise SystemExit("unknown mutation " + which)
if which == "M1" and n != 1: print("LANDS:0 %s count=%d" % (which, n)); raise SystemExit(3)
if which == "M2" and n != 1: print("LANDS:0 %s count=%d" % (which, n)); raise SystemExit(3)
if which == "R1" and n != 3: print("LANDS:0 %s count=%d" % (which, n)); raise SystemExit(3)
io.open(path, "w", encoding="utf-8", newline="").write(out)
print("LANDS: %s applied count=%d marker=%d" % (which, n, out.count("MUTATION-216V3")))
PY
}

run216 () { # $1=输出日志 $2=台面标签 $3=TMP 值（空＝默认）
  if [ -z "$3" ]; then
    ( cd "$ROOT" && go test ./internal/risk/ -count=1 -timeout 20m -v -run "$SEL216" ) >"$1" 2>"$1.err"
  else
    ( cd "$ROOT" && env TMP="$3" TEMP="$3" go test ./internal/risk/ -count=1 -timeout 20m -v -run "$SEL216" ) >"$1" 2>"$1.err"
  fi
  echo "rc=$?" >>"$1"
  echo "  [$2] $(counts "$1") red=$(grep '^--- FAIL' "$1" | sed 's/ (.*//' | tr -s ' ' | paste -sd'; ' -)"
}

restore () { # $1=文件名 -> 重新从锚取该枚文件（⛔ 手抄），再核哈希
  ( cd "$REPO" && git archive "$ANCHOR" "internal/risk/$1" ) | tar -x -C "$ROOT"
  echo "  restored $1 hash=$(git hash-object "$ROOT/internal/risk/$1") marker_residual=$(grep -rc 'MUTATION-216V3' "$ROOT/internal/risk/$1" || true)"
}

echo "########## 2 基线（未突变，D 台面）"
run216 "$HERE/logs/mut-baseline-D.log" "D台面" 'D:\W216V3TMP'
( cd "$ROOT" && env TMP='D:\W216V3TMP' TEMP='D:\W216V3TMP' go test ./internal/risk/ -count=1 -timeout 20m -v ) >"$HERE/logs/mut-baseline-D-full.log" 2>&1; echo "rc=$?" >>"$HERE/logs/mut-baseline-D-full.log"
echo "  [D台面落点具名] $(grep -o 'D:[^" ]*W216V3TMP' "$HERE/logs/mut-baseline-D-full.log" | head -1) 整包 $(counts "$HERE/logs/mut-baseline-D-full.log")"
echo "########## 3 M1（ⓐ 点名的变异正控：算子尺该回 1；ⓑ 的 argv0 臂该同红）"
mutate M1 rules_shell.go && run216 "$HERE/logs/mut-M1-D.log" "D台面" 'D:\W216V3TMP'; restore rules_shell.go
echo "########## 4 M2（ⓑ 的牙证：hostile 臂该红；ⓐ 三窗该仍 0/0/0）"
mutate M2 blacklist.go && run216 "$HERE/logs/mut-M2-D.log" "D台面" 'D:\W216V3TMP'; restore blacklist.go
echo "########## 5 M2 打在 C 卷默认台面（看那枚特红会不会冒充 M2 的牙）"
mutate M2 blacklist.go && run216 "$HERE/logs/mut-M2-Cdefault.log" "C默认台面" ''; restore blacklist.go
echo "########## 6 R1 必答题（整枚改名 -> ⓐ 该红不静默，红句归位 allow_substring_pin_216_test.go:68）"
mutate R1 rules_shell.go && run216 "$HERE/logs/mut-R1-D.log" "D台面" 'D:\W216V3TMP'; restore rules_shell.go
echo "########## 7 还原后端面复跑（两把台面都该回到基线色）"
run216 "$HERE/logs/restored-D-final.log" "D台面" 'D:\W216V3TMP'
run216 "$HERE/logs/restored-C-final.log" "C默认台面" ''
echo "########## 8 收尾哈希（两枚被突变过的件都该等于锚上 blob）"
for f in rules_shell.go blacklist.go; do
  echo "$f tree=$(git hash-object "$ROOT/internal/risk/$f") blob=$(cd "$REPO" && git rev-parse "$ANCHOR:internal/risk/$f")"
done
echo "########## 9 共享工作树未被本腿动过（应空）"
( cd "$REPO" && git status --porcelain -- internal/risk cmd docs/PLAN.md )
echo "########## done"
