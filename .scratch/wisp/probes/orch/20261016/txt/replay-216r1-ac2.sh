#!/usr/bin/env bash
# 编排者收件复跑（票 216 的 216-r1）：AC#2 成对那一发 ⇄ AC#3 去牙那一发，我自己重放
# ⛔ 输出路径带 orch- 前缀（本仓今天发生过：在飞腿复跑同名台件，把我已提交的读数件打成 46 字节）
set -u
REPO="D:/work/workspace/projects plans/Wisp"
ROOT=/d/tmp/orch-216r1-replay                       # ⛔ 树必须落 D: 卷（C: 的 mktemp 路径会被挂成空目录＝静默 rc=1 无日志）
OUT="$REPO/.scratch/wisp/probes/orch/20261016/logs"
PROBE="$REPO/.scratch/wisp/probes/216/r1/pair/probe_216_pair_test.go"
SEVEN="go.mod go.sum internal cmd tools scripts frontend"   # 七件套，少 frontend＝[setup failed] 假红
BEFORE_SHA=ca0e30fc                                   # 本腿第一笔的父（它自报误取过 HEAD~2，这里用 rev-parse 定死的同一枚）
AFTER_SHA=c9ba5254                                    # 本腿末笔
GOTMP='/c/Users/swq/AppData/Local/Temp/orch216replay'
mkdir -p "$GOTMP"
export GOTMPDIR='C:\Users\swq\AppData\Local\Temp\orch216replay'   # 原生反斜杠形（票 317 那一族）
export GOPATH='D:\work\base\gopath'
NEWTESTS='^TestScrubForDisplayEscapesEveryUnseenCharAndTouchesNothingElse$|^TestTheWashReturnsTheSameStringWhenThereIsNothingToWash$|^TestTheSnapshotBoundaryWashesBothRosterFieldsAndCopiesTheRows$|^TestAHostileAttachmentNameWalksTheRealBrokerIntoAWashedMenuRow$|^TestTheClaimInANameChangesNoVerdict$|^TestWorkspaceAbsentAndWorkspaceUnreadAreTwoSentences$|^TestTheFiveInstructionStatesAreStillTheDiskSpelling$'

rm -rf "$ROOT"; mkdir -p "$ROOT" "$OUT"

build_tree () {  # $1=sha $2=目录名
  local d="$ROOT/$2"
  mkdir -p "$d"
  ( cd "$REPO" && git archive "$1" $SEVEN ) | tar -x -C "$d"
  mkdir -p "$d/build" "$d/third_party/sherpa-onnx" "$d/frontend/dist"
  cp "$REPO"/build/*.dll "$d/build/" 2>/dev/null
  cp "$REPO"/third_party/sherpa-onnx/*.dll "$d/third_party/sherpa-onnx/" 2>/dev/null
  cp -r "$REPO"/frontend/dist/. "$d/frontend/dist/" 2>/dev/null
  cp "$PROBE" "$d/internal/panel/probe_216_pair_test.go"
  echo "tree $2 : go_files=$(find "$d/internal" -name '*.go' | wc -l) dll_build=$(ls "$d/build"/*.dll 2>/dev/null | wc -l) dll_sherpa=$(ls "$d/third_party/sherpa-onnx"/*.dll 2>/dev/null | wc -l) dist_entries=$(ls "$d/frontend/dist" | wc -l)"
}

echo "########## 1 建树（台面缺件手工补）"
build_tree "$BEFORE_SHA" tree-before
build_tree "$AFTER_SHA" tree-after

echo "########## 2 两枚 build 门"
( cd "$ROOT/tree-before" && go build ./internal/panel/ ; echo "build-before rc=$?" )
( cd "$ROOT/tree-after"  && go build ./internal/panel/ ; echo "build-after  rc=$?" )

echo "########## 3 AC#2 成对：同一枚探针用例、两棵树"
for t in tree-before tree-after; do
  ( cd "$ROOT/$t" && go test ./internal/panel/ -run '^TestProbe216HostileDisplayNameAtTheDisplayBoundary$' -count=1 -v -timeout 30m > "$OUT/orch-216-ac2-$t.txt" 2>&1; echo "rc=$?" >> "$OUT/orch-216-ac2-$t.txt" )
  echo "----- $t rc=$(tail -1 "$OUT/orch-216-ac2-$t.txt")"
  grep -nE 'probe_216_pair_test\.go:[0-9]+:|^--- (PASS|FAIL|SKIP)|no tests to run|setup failed|build failed|0xc0000135' "$OUT/orch-216-ac2-$t.txt" | head -8
done

echo "########## 4 AC#3 去牙：tree-after 里把 scrubForDisplay 首行改成直接 return s"
SCRUB="$ROOT/tree-after/internal/panel/display_scrub_216.go"
cp "$SCRUB" "$ROOT/scrub.orig"
git -C "$REPO" hash-object "$ROOT/scrub.orig" | sed 's/^/orig-copy sha = /'
git -C "$REPO" rev-parse "$AFTER_SHA":internal/panel/display_scrub_216.go | sed 's/^/committed sha = /'
python - "$SCRUB" <<'PY'
import io, sys
p = sys.argv[1]
src = io.open(p, encoding='utf-8').read().split('\n')
for i, line in enumerate(src):
    if line.startswith('func scrubForDisplay('):
        src.insert(i+1, '\treturn s // MUTATION-216R1-TOOTH')
        break
io.open(p, 'w', encoding='utf-8', newline='\n').write('\n'.join(src))
print('LANDS: inserted at line', i+2)
PY
grep -n 'MUTATION-216R1-TOOTH' "$SCRUB"
( cd "$ROOT/tree-after" && go test ./internal/panel/ -run "$NEWTESTS" -count=1 -v -timeout 30m > "$OUT/orch-216-ac3-teeth-off.txt" 2>&1; echo "rc=$?" >> "$OUT/orch-216-ac3-teeth-off.txt" )
echo "----- 去牙发：红名册"
grep -E '^--- FAIL|^--- PASS|^    --- FAIL' "$OUT/orch-216-ac3-teeth-off.txt" | head -20
echo "top PASS=$(grep -c '^--- PASS' "$OUT/orch-216-ac3-teeth-off.txt") top FAIL=$(grep -c '^--- FAIL' "$OUT/orch-216-ac3-teeth-off.txt")"

echo "########## 5 还原＋sha 自证"
cp "$ROOT/scrub.orig" "$SCRUB"
git -C "$REPO" hash-object "$SCRUB" | sed 's/^/restored sha = /'
( cd "$ROOT/tree-after" && go test ./internal/panel/ -run "$NEWTESTS" -count=1 -v -timeout 30m > "$OUT/orch-216-ac3-restored-green.txt" 2>&1; echo "rc=$?" >> "$OUT/orch-216-ac3-restored-green.txt" )
echo "----- 还原发：top PASS=$(grep -c '^--- PASS' "$OUT/orch-216-ac3-restored-green.txt") top FAIL=$(grep -c '^--- FAIL' "$OUT/orch-216-ac3-restored-green.txt")"
grep -c 'MUTATION-216R1-TOOTH' "$SCRUB" | sed 's/^/remaining MUTATION markers in tree copy = /'
git -C "$REPO" status --porcelain -- internal cmd | sed 's/^/SHARED-TREE-STATUS: /'
echo "########## done"
