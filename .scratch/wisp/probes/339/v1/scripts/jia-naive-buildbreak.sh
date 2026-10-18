#!/bin/bash
# 339-v1 台件⑥：★naive 整档加 tag 会不会打坏 Linux 编译（AC#1 甲形代价读数，我自己复现）。
# 落点＝pre 树（e6fd0f50＝334 件无档、无 339 件）的克隆 naive 树，物理字节加一行 //go:build windows。
# 验命中＝落刀前该文件 build_tag_lines=0、落刀后=1；⛔ 跟踪件、⛔ -overlay。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
TREE="/d/tmp/wisp/20261018/339v1/naive/tree"
F="$TREE/internal/risk/anchor_miss_reason_334_test.go"
LOGDIR="$REPO/.scratch/wisp/probes/339/v1/logs"
mkdir -p "$LOGDIR"
OUT="$LOGDIR/v1-jia-naive-buildbreak.log"
BLOB=$(git -C "$REPO" rev-parse 'e6fd0f5051812a23c72db444d12d76b4d27285a7:internal/risk/anchor_miss_reason_334_test.go')
h() { git -C "$REPO" hash-object "$1"; }
{
  echo "# naive 形＝只给 anchor_miss_reason_334_test.go 加 //go:build windows，夹具不搬（⛔ 339 件在场）"
  echo "tree=$TREE blob_expect=$BLOB"
  echo "P0=$(h "$F")"
  [ "$(h "$F")" = "$BLOB" ] || { echo GUARD-ABORT-pristine; exit 2; }
  echo "build_tag_lines_before=$(grep -c '^//go:build' "$F" || true) (must be 0)"
  printf '%s\n%s\n' '//go:build windows' "$(cat "$F")" > "$F.naive" && mv "$F.naive" "$F"
  echo "build_tag_lines_after=$(grep -c '^//go:build' "$F" || true) (must be 1)"
  echo "first_line=$(head -1 "$F")"
  echo "== Linux 档 go vet =="
  MSYS_NO_PATHCONV=1 docker run --rm -v 'D:/tmp/wisp/20261018/339v1/naive/tree':/src \
    -v 'D:/work/base/gopath/pkg/mod':/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go vet ./internal/risk/ 2>&1; echo VET_RC=\$?; go test -count=1 -run XXX_NONE ./internal/risk/ 2>&1 | tail -6; echo TESTBUILD_RC=\$?"
  echo "== 还原 naive 树＝重取 blob（这棵树专用、不作正控）=="
  git -C "$REPO" show "e6fd0f50:internal/risk/anchor_miss_reason_334_test.go" > "$F"
  echo "R=$(h "$F") match=$([ "$(h "$F")" = "$BLOB" ] && echo YES || echo NO)"
} > "$OUT" 2>&1
cat "$OUT"
