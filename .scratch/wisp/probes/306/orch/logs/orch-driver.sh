#!/usr/bin/env bash
# 306 编排者复跑：牙／反形敏感度，七发，全部在仓外导出树里（⛔ 动共享工作树、⛔ 动 third_party）
# 尺口径：每发自己落一行 rc_=；颜色尺＝grep -c '^--- FAIL' 与 '^--- PASS'，射程＝整份件（含子测试行）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="$(cat "$REPO/.scratch/wisp/probes/306/orch/logs/TREE-PATH.txt")"
L="$REPO/.scratch/wisp/probes/306/orch/logs"
INJ="$TREE/internal/audio/wavinjector.go"
FIX="$TREE/internal/audio/wavinjector_extensible_float_306_test.go"
SHELF="/c/Users/swq/tmp/wisp-306-orch-shelf-$(date +%Y%m%d-%H%M)"
NAMES='^(TestParseWavExtensibleFloat306x|TestParseWavExtensibleFloat32306|TestWavInjectorExtensibleFloat32306)$'
NAMES='^(TestParseWavExtensibleFloat32306|TestWavInjectorExtensibleFloat32306)$'
mkdir -p "$SHELF"

inj_from_head() { ( cd "$REPO" && git show "HEAD:internal/audio/wavinjector.go" ) > "$INJ"; }
fix_from_head() { ( cd "$REPO" && git show "HEAD:internal/audio/wavinjector_extensible_float_306_test.go" ) > "$FIX"; }

# run <tag> <run-regex|ALL>
run() {
  local tag="$1" pat="$2"
  local f="$L/orch-$tag.txt"
  if [ "$pat" = ALL ]; then
    ( cd "$TREE" && GOFLAGS= go test ./internal/audio/ -count=1 -v ) > "$f" 2>&1
  else
    ( cd "$TREE" && GOFLAGS= go test ./internal/audio/ -count=1 -v -run "$pat" ) > "$f" 2>&1
  fi
  echo "rc_$tag=$?" >> "$f"
  printf '%s rc=%s FAIL=%s PASS=%s SKIP=%s\n' \
    "$tag" "$(grep -o "rc_$tag=[0-9]*" "$f" | cut -d= -f2)" \
    "$(grep -c '^--- FAIL' "$f")" "$(grep -c '^--- PASS' "$f")" "$(grep -c '^--- SKIP' "$f")"
}

echo "== 起手：工作树⛔ 被碰的自证（scoped porcelain）"
( cd "$REPO" && git status --porcelain -- internal cmd ) > "$L/orch-00-worktree-clean.txt" 2>&1
echo "rc_porcelain=$?" >> "$L/orch-00-worktree-clean.txt"
wc -l < "$L/orch-00-worktree-clean.txt"

echo "== 发 A 台面自检（pristine，两枚具名用例必须发出颜色）"
inj_from_head; fix_from_head
run A "$NAMES"

echo "== 发 B mut192（0xFFFE -> 0xFFFD，行内权威值取坏）"
perl -i -pe 's/fmtTag == 0xFFFE/fmtTag == 0xFFFD/' "$INJ"
grep -n 'fmtTag == 0xFFFD' "$INJ" >> "$L/orch-B-mut-landed.txt"
run B "$NAMES"
inj_from_head
cmp -s "$INJ" <( cd "$REPO" && git show HEAD:internal/audio/wavinjector.go ) && echo "restored_B=identical"

echo "== 发 C mut196 上下界同移（data[body+24:body+26] -> [body+26:body+28]）"
perl -i -pe 's/data\[body\+24 : body\+26\]/data[body+26 : body+28]/' "$INJ"
grep -n 'data\[body+26 : body+28\]' "$INJ" >> "$L/orch-C-mut-landed.txt"
run C "$NAMES"
inj_from_head

echo "== 发 D mut218（bits == 32 -> 16）"
perl -i -pe 's/case fmtTag == 3 && bits == 32:/case fmtTag == 3 \&\& bits == 16:/' "$INJ"
grep -n 'case fmtTag == 3 && bits == 16:' "$INJ" >> "$L/orch-D-mut-landed.txt"
run D "$NAMES"
inj_from_head

echo "== 发 E 无夹具正控（新用例 mv 走 + mut192 + 整包）"
mv "$FIX" "$SHELF/"
perl -i -pe 's/fmtTag == 0xFFFE/fmtTag == 0xFFFD/' "$INJ"
run E ALL
inj_from_head
mv "$SHELF/$(basename "$FIX")" "$FIX"

echo "== 发 F FLIP 两枚（把断言的等号倒过来：码没动、用例必须响）"
perl -i -pe 's/\tif rate != 48000 \{/\tif rate == 48000 {/' "$FIX"
grep -n 'if rate == 48000' "$FIX" >> "$L/orch-F1-flip-landed.txt"
run F1 "$NAMES"
fix_from_head
perl -i -pe 's/\tif len\(inj\.samples\) != len\(w306WantI16\) \{/\tif len(inj.samples) == len(w306WantI16) {/' "$FIX"
grep -n 'if len(inj.samples) == len(w306WantI16)' "$FIX" >> "$L/orch-F2-flip-landed.txt"
run F2 "$NAMES"
fix_from_head

echo "== 收局：两枚文件对 HEAD blob 逐字节还原自证"
for p in internal/audio/wavinjector.go internal/audio/wavinjector_extensible_float_306_test.go; do
  ( cd "$REPO" && git show "HEAD:$p" ) > /tmp/orch-blob-ref.txt 2>/dev/null
  cmp -s "$TREE/$p" /tmp/orch-blob-ref.txt && echo "cmp_$p=identical rc=0" || echo "cmp_$p=DIFFER rc=1"
done
run Z "$NAMES"
echo "== 完（临时件只建不删：$SHELF 与 $TREE 均留下）"
