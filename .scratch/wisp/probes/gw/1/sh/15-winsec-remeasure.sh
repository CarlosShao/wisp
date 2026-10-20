#!/bin/sh
# gw-1 味③ 撞红复量＝对甲（gwW=3a509e50 ⇄ gwL=HEAD）本机档那 2/2 逐名换位，公共块令＝先安静复量 3 发再归因。
# 摆动在册那几枚（票 342/348 的 racer 钉＋票 353 新钉）两侧都出现过＝⛔ 单发换位当"新增红"；
# ★并集作差＝判据；另加一把"两侧交替同发序"的尺（⛔ 先打完一侧再打另一侧＝负载漂移会冒充字节差）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
run_host() { # $1=tree $2=outfile
  T="$R/$1"
  TS=$(date +%Y%m%d-%H%M%S)-$$
  GOTMPW="${LOCALAPPDATA:-C:\\Users\\swq\\AppData\\Local}\\wisp-test-gotmp\\gw1-$TS"
  mkdir -p "$(cygpath -u "$GOTMPW")" || { echo "CANNOT_CREATE_GOTMP=$GOTMPW"; return 7; }
  TW=$(cygpath -w "$T")
  DLLP=$(cygpath -w "$T/build"); SHP=$(cygpath -w "$T/third_party/sherpa-onnx")
  ( cd "$TW" && PATH="$SHP;$DLLP;$PATH" GOTMPDIR="$GOTMPW" TMPDIR="$GOTMPW" \
      GOCACHE='D:\tmp\wisp\20261020\gw1\gocache-win' \
      go test -v -count=1 -timeout 30m ./internal/winsec/ ) > "$2" 2>&1
  return $?
}
count() { # $1=raw $2=label
  RUN=$(grep -c '^=== RUN' "$1" || true); PT=$(grep -c '^--- PASS' "$1" || true)
  FT=$(grep -c '^--- FAIL' "$1" || true); SK=$(grep -c '^--- SKIP' "$1" || true)
  SHAPE=OK; [ "$RUN" = 0 ] && SHAPE=VOID
  echo "$2 RC=$3 RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK SHAPE=$SHAPE" | tee -a "$LOGS/four-counts.tsv"
  grep -E '^--- FAIL' "$1" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)$//' | sort > "$1.rednames.txt"
  echo "   RED=$(grep -c . "$1.rednames.txt" || true)"; grep . "$1.rednames.txt" | sed 's|^|     |'
}
echo "== START anchor =="; git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
for S in 1 2 3; do
  echo "### 复量发 $S ###"
  for T in gwW gwL; do
    OUT="$LOGS/ws3-$S-$T.raw.txt"
    run_host "$T" "$OUT"; RC=$?
    count "$OUT" "ws3 s$S $T MODE=host" "$RC"
  done
  echo "  -- 单发作差 s$S --"
  comm -13 "$LOGS/ws3-$S-gwW.raw.txt.rednames.txt" "$LOGS/ws3-$S-gwL.raw.txt.rednames.txt" | sed 's|^|   NEWRED |'
  comm -23 "$LOGS/ws3-$S-gwW.raw.txt.rednames.txt" "$LOGS/ws3-$S-gwL.raw.txt.rednames.txt" | sed 's|^|   GONERED |'
done
echo
echo "== ★三发并集作差（＝判据那一把）=="
cat "$LOGS"/ws3-1-gwW.raw.txt.rednames.txt "$LOGS"/ws3-2-gwW.raw.txt.rednames.txt "$LOGS"/ws3-3-gwW.raw.txt.rednames.txt | sort -u > "$LOGS/ws3-union-W.txt"
cat "$LOGS"/ws3-1-gwL.raw.txt.rednames.txt "$LOGS"/ws3-2-gwL.raw.txt.rednames.txt "$LOGS"/ws3-3-gwL.raw.txt.rednames.txt | sort -u > "$LOGS/ws3-union-L.txt"
echo "UNION_W=$(grep -c . "$LOGS/ws3-union-W.txt" || true)"; grep . "$LOGS/ws3-union-W.txt" | sed 's|^|  W |'
echo "UNION_L=$(grep -c . "$LOGS/ws3-union-L.txt" || true)"; grep . "$LOGS/ws3-union-L.txt" | sed 's|^|  L |'
echo "NEW_RED_UNION=$(comm -13 "$LOGS/ws3-union-W.txt" "$LOGS/ws3-union-L.txt" | grep -c . || true)"
comm -13 "$LOGS/ws3-union-W.txt" "$LOGS/ws3-union-L.txt" | sed 's|^|  NEWRED_UNION |'
echo "GONE_RED_UNION=$(comm -23 "$LOGS/ws3-union-W.txt" "$LOGS/ws3-union-L.txt" | grep -c . || true)"
comm -23 "$LOGS/ws3-union-W.txt" "$LOGS/ws3-union-L.txt" | sed 's|^|  GONERED_UNION |'
echo "== END anchor =="; git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
