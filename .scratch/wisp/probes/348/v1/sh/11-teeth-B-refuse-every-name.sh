#!/bin/sh
# 342-v2 必答(iii) 发2＝反形"逢名就拒"（★落点与腿报不同处＝腿在 openPinnedLeaf 函数体首行插 EACCES；
# 本腿把 pinnedParent 的父钉调用整个换成"无条件拒绝"＝票 342/348 AC#3 点名要照红的那一族的极端形）。
# ★形状自查：⛔ 复合字面量当接收者（leafPin{}.m() 需要额外括号＝上一波那种"编译不过、四数全 0"的假读数），
#   这里换成包级函数 refuseEveryName342v2(dir)，返回值类型与 pinPlacement 一致 ⇒ 其余字节一字不动。
# 只打仓外拷贝树 mut-B 的物理字节；⛔ 改跟踪文件、⛔ -overlay；与发 1 串行。
set -u
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"; mkdir -p "$LOGS"
SRC="$ROOT/treeL"; MUT="$ROOT/mut-B"
T=internal/winsec/pinnedleaf_342_unix.go
[ -d "$MUT" ] || cp -a "$SRC" "$MUT"
PRISTINE=$(sha256sum "$SRC/$T" | cut -c1-16)
NOW=$(sha256sum "$MUT/$T" | cut -c1-16)
echo "PRISTINE_SHA16=$PRISTINE MUT_BEFORE=$NOW PRISTINE_MATCH_BEFORE=$([ "$PRISTINE" = "$NOW" ] && echo YES || echo NO)"
echo "OLD_ANCHOR_HITS_BEFORE=$(grep -c 'pinPlacement(dir)' "$MUT/$T")"
perl -0pi -e 's/\tpin, err := pinPlacement\(dir\)/\tpin, err := refuseEveryName342v2(dir) \/\/ MUTATION-B: the ancestor guard refuses every name/' "$MUT/$T"
echo "OLD_ANCHOR_HITS_AFTER=$(grep -c 'pinPlacement(dir)' "$MUT/$T" || true)"
echo "NEW_ANCHOR_HITS_AFTER=$(grep -c 'MUTATION-B: the ancestor guard refuses every name' "$MUT/$T")"
cat >> "$MUT/$T" <<'GOEOF'

// MUTATION-B: the anti-form both tickets' AC#3 demand a sensitive ruler answer for - the ancestor
// guard turned into "refuse every name", so a pre-existing name can never be overwritten.
func refuseEveryName342v2(dir string) (*os.File, error) {
	return nil, &os.PathError{Op: "open", Path: dir, Err: syscall.EACCES}
}
GOEOF
echo "HELPER_HITS=$(grep -c 'MUTATION-B: the anti-form' "$MUT/$T")"
echo "MUT_SHA16_AFTER=$(sha256sum "$MUT/$T" | cut -c1-16)"
docker run --rm -v "$MUT:/src" -v /d/work/base/gopath/pkg/mod:/gm:ro -v "$ROOT/gocache:/root/.cache/go-build" -w /src \
  -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
  golang:1.27 go vet ./internal/winsec/ > "$LOGS/iii-B-vet.txt" 2>&1
echo "VET_RC=$?"; head -5 "$LOGS/iii-B-vet.txt"
LABEL=iii-B-refuse-every-name
docker run --rm \
  -v "$MUT:/src" -v /d/work/base/gopath/pkg/mod:/gm:ro -v "$ROOT/gocache:/root/.cache/go-build" -w /src \
  -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
  golang:1.27 \
  go test ./internal/winsec/ -v -count=1 -timeout 40m -run 'TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1|342r1' > "$LOGS/$LABEL.raw.txt" 2>&1
echo "RC=$?"
RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
echo "FOUR RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK"
grep -E '^(    )?--- (FAIL|PASS|SKIP)' "$LOGS/$LABEL.raw.txt" || true
grep -E 'AC#2 RED' "$LOGS/$LABEL.raw.txt" | head -8
cp -f "$SRC/$T" "$MUT/$T"
echo "RESTORE_SHA16=$(sha256sum "$MUT/$T" | cut -c1-16) PRISTINE_MATCH_AFTER=$([ "$PRISTINE" = "$(sha256sum "$MUT/$T" | cut -c1-16)" ] && echo YES || echo NO)"
