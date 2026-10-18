# 339-r1 AC#0 台件④：甲-form 的 Linux 分母尺（两棵导出树各跑整包，逐名点名）。
# 用法：bash jia-denom.sh
# 树 A＝HEAD 原样（⛔ probe 件）；树 B＝甲形＝334 文件拆成 world（无档）＋断言（//go:build windows）。
# 判据＝两棵树的 RUN/PASS/FAIL/SKIP 逐名作差 ⇒ 甲在 Linux 少跑哪几枚（分母变化，逐枚点名）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
ROOT="/d/tmp/wisp/20261018/339r1"
A="$ROOT/jiaA/tree"
B="$ROOT/jiaB/tree"
LOGD="$REPO/.scratch/wisp/probes/339/r1/logs"
cd "$REPO" || exit 9
ANCHOR=e6fd0f5051812a23c72db444d12d76b4d27285a7
for t in "$A" "$B"; do
  rm -rf "$t" 2>/dev/null
  mkdir -p "$t"
  git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models docs | tar -xf - -C "$t"
done
echo "A go=$(find "$A" -name '*.go' | wc -l) B go=$(find "$B" -name '*.go' | wc -l)"
# 树 B 手工改：334 文件的 fixture 区（1..112 行）无档留在 POSIX；114 行起（断言区）加 //go:build windows。
# 甲-form 落盘形状与此等价：world 件带档⛔ 动 336/216 字节，断言区搬进带档新件。
python3 - "$B" <<'PY' 2>/dev/null || {
  F="$1/internal/risk/anchor_miss_reason_334_test.go"
  head -112 "$F" > "$F.new"
  { printf '//go:build windows\n\npackage risk\n\nimport (\n\t"os"\n\t"path/filepath"\n\t"strconv"\n\t"strings"\n\t"testing"\n)\n\n'; sed -n '114,$p' "$F"; } >> "$F.new"
  mv "$F.new" "$F"
  echo "split done"
} || echo "python split done"
PY
MSYS_NO_PATHCONV=1
run_leg() { # $1=tree $2=label
  docker run --rm -v "$1":/src -v "$(go env GOMODCACHE)":/go/pkg/mod:ro -w /src \
    -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=0 -e HOME=/src -e GOMODCACHE=/go/pkg/mod \
    golang:1.27 bash -c "mkdir -p /tmp/gt; export TMPDIR=/tmp/gt GOTMPDIR=/tmp/gt; go test ./internal/risk/ -count=1 -timeout 30m -v > /tmp/o 2>&1; echo GO_TEST_RC=\$?; cat /tmp/o" > "$LOGD/$2.log" 2>&1
  echo "== $2: RUN=$(grep -c '^=== RUN' "$LOGD/$2.log" || true) PASS=$(grep -c '^--- PASS' "$LOGD/$2.log" || true) FAIL=$(grep -c '^--- FAIL' "$LOGD/$2.log" || true) SKIP=$(grep -c '^--- SKIP' "$LOGD/$2.log" || true) $(grep '^GO_TEST_RC=' "$LOGD/$2.log")"
  grep -E '^--- (FAIL|SKIP)' "$LOGD/$2.log" || true
  grep -nE 'build failed|setup failed|no test files|creating work dir' "$LOGD/$2.log" || true
}
run_leg "$A" jiaA-linux-baseline
run_leg "$B" jiaB-linux-jiaform
echo "=== 逐名作差（A PASS\B PASS = 甲形在 Linux 少跑的顶层名）==="
grep -E '^--- PASS' "$LOGD/jiaA-linux-baseline.log" | sed 's/^--- PASS: //; s/ .*//' | sort > "$LOGD/.a.pass"
grep -E '^--- PASS' "$LOGD/jiaB-linux-jiaform.log" | sed 's/^--- PASS: //; s/ .*//' | sort > "$LOGD/.b.pass"
comm -23 "$LOGD/.a.pass" "$LOGD/.b.pass"
echo "=== 反方向（B 有 A 无）==="
comm -13 "$LOGD/.a.pass" "$LOGD/.b.pass"
