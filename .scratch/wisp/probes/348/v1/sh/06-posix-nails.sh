#!/bin/sh
# 342-v2 必答(i)：POSIX 两枚新常驻钉（342 截断支 ⇄ 348 删除支）＋票 348 确定性件，
# 基线树 treeB（winsec.go=67783ae0）⇄ 落地树 treeL（07300bff）各 4 发。
# ★容器 TMPDIR/GOTMPDIR ⛔ 绑 Windows 卷＝留容器内默认 /tmp（342-v1 定式：bind mount 上 os.SameFile 失真＝身份门被台面摘掉）。
# ★GOCACHE 指到 D: 上宿主目录（跨发复用编译缓存；⛔ 测试临时根）。
set -u
ROOT=/d/tmp/wisp/20261019/342v2
mkdir -p "$ROOT/logs" "$ROOT/gocache"
NAILSET='TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1|TestAC2POSIXRefusalPathWithASwappedAncestorNeverDeletesAForeignFile342r1|TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1'
for TREE in treeB treeL; do
  for SHOT in 1 2 3 4; do
    LABEL="i01-${TREE}-shot${SHOT}"
    docker run --rm \
      -v "$ROOT/$TREE:/src" \
      -v /d/work/base/gopath/pkg/mod:/gm:ro \
      -v "$ROOT/logs:/logs" \
      -v "$ROOT/gocache:/root/.cache/go-build" \
      -w /src \
      -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm \
      golang:1.27 \
      go test ./internal/winsec/ -v -count=1 -timeout 40m -run "$NAILSET" > "$ROOT/logs/$LABEL.raw.txt" 2>&1
    RC=$?
    RUN=$(grep -c '^=== RUN' "$ROOT/logs/$LABEL.raw.txt" || true)
    PT=$(grep -c '^--- PASS' "$ROOT/logs/$LABEL.raw.txt" || true)
    FT=$(grep -c '^--- FAIL' "$ROOT/logs/$LABEL.raw.txt" || true)
    SK=$(grep -c '^--- SKIP' "$ROOT/logs/$LABEL.raw.txt" || true)
    echo "$LABEL RC=$RC RUN=$RUN PASS_TOP=$PT FAIL_TOP=$FT SKIP=$SK" | tee -a "$ROOT/logs/four-counts-posix.tsv"
    grep -E 'attempts=[0-9]+' "$ROOT/logs/$LABEL.raw.txt" | sed -E 's/^ +[a-z0-9_]+\.go:[0-9]+: //' | head -20
    echo "-----"
  done
done
