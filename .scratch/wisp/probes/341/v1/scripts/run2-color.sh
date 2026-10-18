#!/usr/bin/env bash
# 341-v1 AC#4: paired colour runs on ./internal/config/ -- base tree (pre-341)
# vs delivery tree, on two temp-root tables (A = default C volume, B = D: pre-created).
# Gauges: RUN=^=== RUN / PASS_TOP=^--- PASS / FAIL_TOP=^--- FAIL / SKIP=^--- SKIP.
W=/d/tmp/wisp/20261018/341v1
LOGROOT="$1"   # repo probe logs dir (posix)
cd "$W" || exit 9

one() {
  tree="$1"; label="$2"; tmpmode="$3"
  d="$W/run-$label"; mkdir -p "$d/gotmp" "$d/gocache"
  ( cd "$W/$tree" || exit 9
    case "$tmpmode" in
      B) export GOTMPDIR="$(cygpath -w "$d/gotmp")" ;;
      *) unset GOTMPDIR TMPDIR TEMP TMP ;;
    esac
    export GOCACHE="$(cygpath -w "$d/gocache")"
    export PATH="$(cygpath -w "$PWD/third_party/sherpa-onnx");$(cygpath -w "$PWD/build");$PATH"
    out="$d/test.log"
    go test ./internal/config/ -count=1 -timeout 30m -v > "$out" 2>&1
    rc=$?
    echo "### tree=$tree label=$label tmp=$tmpmode rc=$rc"
    echo "RUN=$(grep -c '^=== RUN' $out || true)"
    echo "PASS_TOP=$(grep -c '^--- PASS' $out || true)"
    echo "FAIL_TOP=$(grep -c '^--- FAIL' $out || true)"
    echo "SKIP=$(grep -c '^--- SKIP' $out || true)"
    echo "--- FAIL roster ---"
    grep '^--- FAIL' "$out" | sed 's/^--- FAIL: //' | cut -d' ' -f1 | sort || true
    echo "--- tail 12 ---"; tail -12 "$out"
  )
}

mkdir -p "$LOGROOT" || exit 9
one treeB base-A   A | tee -a "$LOGROOT/03-color-base.txt"
one treeB base-B   B | tee -a "$LOGROOT/03-color-base.txt"
one treeD deliv-A  A | tee -a "$LOGROOT/04-color-deliv.txt"
one treeD deliv-B  B | tee -a "$LOGROOT/04-color-deliv.txt"
echo "DONE"
