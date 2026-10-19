#!/bin/sh
# 342-v2 必答(iii) 发1＝去牙（★落点与腿报不同处＝腿摘的是父钉那支 pinPlacement→os.Open；
# 本腿摘的是票 348 自己那支：unlinkPinnedLeaf 的 Unlinkat ⇒ 改回按名 unlink）。
# 只打仓外拷贝树 mut-A 的物理字节；⛔ 改跟踪文件、⛔ -overlay；两发串行。
set -u
ROOT=/d/tmp/wisp/20261019/342v2
LOGS="$ROOT/logs"; mkdir -p "$LOGS"
SRC="$ROOT/treeL"; MUT="$ROOT/mut-A"
T=internal/winsec/pinnedleaf_342_unix.go
[ -d "$MUT" ] || cp -a "$SRC" "$MUT"
PRISTINE=$(sha256sum "$SRC/$T" | cut -c1-16)
NOW=$(sha256sum "$MUT/$T" | cut -c1-16)
echo "PRISTINE_SHA16=$PRISTINE MUT_BEFORE=$NOW PRISTINE_MATCH_BEFORE=$([ "$PRISTINE" = "$NOW" ] && echo YES || echo NO)"
echo "OLD_ANCHOR_HITS_BEFORE=$(grep -c 'unix.Unlinkat(int(p.fd()), base, 0)' "$MUT/$T")"
perl -0pi -e 's/\t_ = unix\.Unlinkat\(int\(p\.fd\(\)\), base, 0\)/\t_ = syscall.Unlink(filepath.Join(p.name, base)) \/\/ MUTATION-A: back to by-name unlink/' "$MUT/$T"
echo "OLD_ANCHOR_HITS_AFTER=$(grep -c 'unix.Unlinkat(int(p.fd()), base, 0)' "$MUT/$T" || true)"
echo "NEW_ANCHOR_HITS_AFTER=$(grep -c 'MUTATION-A: back to by-name unlink' "$MUT/$T")"
echo "MUT_SHA16_AFTER=$(sha256sum "$MUT/$T" | cut -c1-16)"
LABEL=iii-A-unlink-byname
docker run --rm \
  -v "$MUT:/src" -v /d/work/base/gopath/pkg/mod:/gm:ro -v "$ROOT/gocache:/root/.cache/go-build" -w /src \
  -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
  golang:1.27 \
  go test ./internal/winsec/ -v -count=1 -timeout 40m -run '342r1|TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1' > "$LOGS/$LABEL.raw.txt" 2>&1
echo "RC=$?"
RUN=$(grep -c '^=== RUN' "$LOGS/$LABEL.raw.txt" || true)
PT=$(grep -c '^--- PASS' "$LOGS/$LABEL.raw.txt" || true)
FT=$(grep -c '^--- FAIL' "$LOGS/$LABEL.raw.txt" || true)
SK=$(grep -c '^--- SKIP' "$LOGS/$LABEL.raw.txt" || true)
echo "FOUR RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK"
grep -E '^--- (FAIL|PASS|SKIP)' "$LOGS/$LABEL.raw.txt" || true
grep -E 'AC#(1|2) RED' "$LOGS/$LABEL.raw.txt" | head -8
cp -f "$SRC/$T" "$MUT/$T"
echo "RESTORE_SHA16=$(sha256sum "$MUT/$T" | cut -c1-16) PRISTINE_MATCH_AFTER=$([ "$PRISTINE" = "$(sha256sum "$MUT/$T" | cut -c1-16)" ] && echo YES || echo NO)"
