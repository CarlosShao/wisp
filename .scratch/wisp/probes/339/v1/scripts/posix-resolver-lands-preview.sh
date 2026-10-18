#!/bin/bash
# 339-v1 台件⑩：必答(iv) 的★现量分母＝"POSIX 真解析器落地那一刻，Linux 档整包红名册是谁"。
# 落点＝仓外导出树物理字节（⛔ 跟踪件、⛔ -overlay）：把 pathresolver_other.go 的 resolveHandle
# 从 ("",false) 桩换成 realpath 形（filepath.Abs + EvalSymlinks，SPEC-06 第4节 那句"realpath + lstat
# symlink reject"的最小可跑形）。定式＝落刀前旧锚＝1、落刀后旧锚＝0 且新锚＝1；跑前跑后各验 pristine。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ANCHOR=7437ef882fac5b427bf17c4bf590819c58d0b375
TREE=/d/tmp/wisp/20261018/339v1/realres/tree
PR="$TREE/internal/risk/pathresolver_other.go"
LOGDIR="$REPO/.scratch/wisp/probes/339/v1/logs"
mkdir -p "$LOGDIR"
OUT="$LOGDIR/v1-posix-resolver-lands-preview.log"
PRISTINE=$(git -C "$REPO" rev-parse "$ANCHOR:internal/risk/pathresolver_other.go")
h() { git -C "$REPO" hash-object "$1"; }
{
  echo "### 必答(iv) 分母尺：resolveHandle 落地成 realpath ⇒ Linux 整包谁红（逐名点名）"
  echo "anchor=$ANCHOR tree=$TREE"
  echo "P0=$(h "$PR") expect=$PRISTINE"
  [ "$(h "$PR")" = "$PRISTINE" ] || { echo GUARD-ABORT-pristine; exit 2; }
  echo "old_anchor_hits_before=$(grep -cF 'return "", false }' "$PR" || true) (must be 1)"
  cat > "$PR" <<'EOF'
//go:build !windows

package risk

import (
	"path/filepath"
)

// MUT339V1-REALRES: preview of what a landed POSIX resolver would do (realpath
// form per SPEC-06 第4节), used only to measure which instruments go red the
// moment pathForms.root stops being empty on this platform. NOT a deliverable.
func resolveHandle(p string) (string, bool) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", false
	}
	real, err2 := filepath.EvalSymlinks(abs)
	if err2 != nil {
		return "", false
	}
	return real, true
}

func reparseComponents(p string) []string { return nil }
EOF
  echo "MUT=$(h "$PR") old_anchor_hits_after=$(grep -cF 'func resolveHandle(p string) (string, bool) { return "", false }' "$PR" || true) (old must be 0) new_anchor_hits=$(grep -cF 'MUT339V1-REALRES' "$PR" || true) (new must be 1)"
  MSYS_NO_PATHCONV=1 docker run --rm -v 'D:/tmp/wisp/20261018/339v1/realres/tree':/src \
    -v 'D:/work/base/gopath/pkg/mod':/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -count=1 -timeout 30m -v 2>&1; echo RC=\$?"
  git -C "$REPO" show "$ANCHOR:internal/risk/pathresolver_other.go" > "$PR"
  echo "R=$(h "$PR") match=$([ "$(h "$PR")" = "$PRISTINE" ] && echo YES || echo NO)"
} > "$OUT" 2>&1
echo "wrote=$OUT bytes=$(wc -c < "$OUT")"
echo "RUN=$(grep -c '^=== RUN' "$OUT" || true) PASS_TOP=$(grep -c '^--- PASS' "$OUT" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$OUT" || true) SKIP=$(grep -c '^--- SKIP' "$OUT" || true)"
echo '=== red roster ==='
grep -E '^--- FAIL' "$OUT" || true
echo '=== 台面故障 ==='
grep -nE 'build failed|no test files|setup failed|creating work dir|0xc0000135' "$OUT" || echo none
