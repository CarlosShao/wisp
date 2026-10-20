#!/bin/sh
# 353-r1 交件前终量尺＝★随笔④动过字节 ⇒ 所有读数重取（⛔ 拿笔④之前的数顶替＝公共块读数纪律）。
# 台面重建＝两棵树都从**锚 a59f5630 的 git archive** 重新取出，再只喂仓外物理字节：
#   treeF-B（基线侧）＝产码四枚＝锚 blob ＋ 测试两枚＝HEAD blob（★两侧同一对测试字节⇒逐名作差只反映产码）
#   treeF-L（落地侧）＝产码四枚＋测试两枚＝HEAD blob
# ★⛔ 碰共享工作树字节；★⛔ 复用别腿的树；★串行。
set -u
R=/d/tmp/wisp/20261020/353r1
ROOTW='D:\tmp\wisp\20261020\353r1'
REPO='/d/work/workspace/projects plans/Wisp'
LOGS="$R/logs"
mkdir -p "$LOGS" "$R/base" "$R/incoming/treeL" "$R/incoming/treeB"
RUNPAT='TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1|TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1|TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2'
cd "$REPO" || exit 9

echo "=== V-0 取字节（一律走 git cat-file blob＝⛔ 把工作树当权威）==="
for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go \
         nametopup_353_windows_test.go ancestorwindow_342r2_windows_test.go; do
  git cat-file blob "HEAD:internal/winsec/$f" > "$R/incoming/treeL/$f"
  echo "  HEAD $f $(git rev-parse "HEAD:internal/winsec/$f" | cut -c1-12) $(wc -c < "$R/incoming/treeL/$f")B"
done
for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go; do
  git cat-file blob "a59f5630:internal/winsec/$f" > "$R/base/$f"
  echo "  BASE $f $(git rev-parse "a59f5630:internal/winsec/$f" | cut -c1-12)"
done

build_tree () {
  name="$1"; src="$2"
  T="$R/$name"
  rm -rf "$T" 2>/dev/null || true
  mkdir -p "$T"
  git archive --format=tar a59f5630 go.mod go.sum internal cmd tools scripts frontend models docs | tar -x -C "$T"
  mkdir -p "$T/build" "$T/third_party/sherpa-onnx" "$T/frontend/dist"
  cp -f build/*.dll "$T/build/" 2>/dev/null || true
  cp -f third_party/sherpa-onnx/*.dll "$T/third_party/sherpa-onnx/" 2>/dev/null || true
  for f in nametopup_353_windows_test.go ancestorwindow_342r2_windows_test.go; do
    cp -f "$R/incoming/treeL/$f" "$T/internal/winsec/$f"
  done
  if [ "$src" = base ]; then
    for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go; do
      cp -f "$R/base/$f" "$T/internal/winsec/$f"
    done
  else
    for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go winsec.go; do
      cp -f "$R/incoming/treeL/$f" "$T/internal/winsec/$f"
    done
  fi
  echo "  TREE $name GO=$(find "$T" -name '*.go' | wc -l) pinnedleaf=$(sha256sum "$T/internal/winsec/pinnedleaf_342_windows.go" | cut -d' ' -f1 | cut -c1-12) winsec=$(sha256sum "$T/internal/winsec/winsec.go" | cut -d' ' -f1 | cut -c1-12) test=$(sha256sum "$T/internal/winsec/nametopup_353_windows_test.go" | cut -d' ' -f1 | cut -c1-12)"
}
echo "=== V-1 建两棵终量树 ==="
build_tree finalB base
build_tree finalL head

run_pkg () {  # <tree> <label> <out> <pkg...>
  tree="$1"; label="$2"; out="$3"; shift 3
  TS=$(date +%Y%m%d-%H%M%S)-$$
  GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\353r1f-$TS"
  mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; return 7; }
  DLLP=$(cygpath -w "$ROOTW\\$tree\\build"); SHP=$(cygpath -w "$ROOTW\\$tree\\third_party\\sherpa-onnx")
  ( cd "$ROOTW\\$tree" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
      GOCACHE="$ROOTW\\gocache-win" go test -v -count=1 -timeout 60m "$@" ) > "$out" 2>&1
  RC=$?
  echo "$label tree=$tree RC=$RC RUN=$(grep -c '^=== RUN' "$out" || true) PASS_TOP=$(grep -c '^--- PASS' "$out" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true) SKIP=$(grep -c '^--- SKIP' "$out" || true) GOTMP=$GOTMPW" | tee -a "$LOGS/final-fourcounts.tsv"
}

echo "=== V-2 定向成对（★到达数分子）==="
run_pkg finalB v2dir "$LOGS/final-dir-B.raw.txt" -run "$RUNPAT" ./internal/winsec/
grep -E '353r1 AC#2 over|AC#1-one over|353r1 AC#2 RED:|353r1 AC#2 CONTROL' "$LOGS/final-dir-B.raw.txt" | sed 's/^/   B| /'
run_pkg finalL v2dir "$LOGS/final-dir-L.raw.txt" -run "$RUNPAT" ./internal/winsec/
grep -E '353r1 AC#2 over|AC#1-one over|353r1 AC#2 RED:|353r1 AC#2 CONTROL' "$LOGS/final-dir-L.raw.txt" | sed 's/^/   L| /'

echo "=== V-3 整包逐名矩阵（各 2 发＝摆动名册靠复量、⛔ 单发定论）==="
for shot in 1 2; do
  run_pkg finalB "v3fullB$shot" "$LOGS/final-full-B-$shot.raw.txt" ./internal/winsec/
  grep -E '^--- FAIL' "$LOGS/final-full-B-$shot.raw.txt" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/final-full-B-$shot.rednames.txt"
  run_pkg finalL "v3fullL$shot" "$LOGS/final-full-L-$shot.raw.txt" ./internal/winsec/
  grep -E '^--- FAIL' "$LOGS/final-full-L-$shot.raw.txt" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$LOGS/final-full-L-$shot.rednames.txt"
done
echo "  --- 逐名矩阵（两侧各 2 发）---"
cat "$LOGS"/final-full-B-*.rednames.txt "$LOGS"/final-full-L-*.rednames.txt | sort -u > "$LOGS/final-allnames.txt"
while read -r n; do
  [ -n "$n" ] || continue
  b=""; l=""
  for s in 1 2; do
    grep -qxF "$n" "$LOGS/final-full-B-$s.rednames.txt" && b="$b RED" || b="$b green"
    grep -qxF "$n" "$LOGS/final-full-L-$s.rednames.txt" && l="$l RED" || l="$l green"
  done
  echo "   $n |B:$b |L:$l"
done < "$LOGS/final-allnames.txt"
echo "  --- 新增红（★尺＝名集，⛔ 四数差）---"
cat "$LOGS"/final-full-L-*.rednames.txt | sort -u > "$LOGS/final-Lunion.txt"
cat "$LOGS"/final-full-B-*.rednames.txt | sort -u > "$LOGS/final-Bunion.txt"
comm -13 "$LOGS/final-Bunion.txt" "$LOGS/final-Lunion.txt" | sed 's/^/   ONLY_IN_LANDED: /' || true
echo "   NEW_RED_COUNT=$(comm -13 "$LOGS/final-Bunion.txt" "$LOGS/final-Lunion.txt" | grep -c . || true)"
echo "   DISAPPEARED_IN_LANDED=$(comm -23 "$LOGS/final-Bunion.txt" "$LOGS/final-Lunion.txt" | grep -c . || true)"
comm -23 "$LOGS/final-Bunion.txt" "$LOGS/final-Lunion.txt" | sed 's/^/   only-B: /'
