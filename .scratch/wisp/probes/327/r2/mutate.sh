#!/bin/sh
# 327-r2 paired-mutation driver. Lands ONLY on the OUT-OF-REPO export tree
# (/d/tmp/w327r2/bench). Tracked bytes are NEVER touched. Every shot restores the
# scrub file from the named commit tree in a trap (finally) and reports the restored
# hash, which must read back 8050f22d... each time.
set -u
BENCH=/d/tmp/w327r2/bench
LOG=/d/tmp/w327r2/raw
SCRUB=internal/panel/display_scrub_216.go
ANCHOR=7535d83b
PRISTINE=8050f22d14203bba006cbfabcfb2a748081286d492691109c3a35a7a435ff8fd
export GOTMPDIR=/d/tmp/w327r2/gotmp
export PATH="$BENCH/third_party/sherpa-onnx:$BENCH/build:$PATH"

restore() {
  ( cd "$BENCH" && git -C "/d/work/workspace/projects plans/Wisp" \
      show "$ANCHOR:$SCRUB" > "$SCRUB" )
  h=$(sha256sum "$BENCH/$SCRUB" | cut -d' ' -f1)
  echo "RESTORE_OK sha256=$h"
  [ "$h" = "$PRISTINE" ] || echo "RESTORE_MISMATCH (want $PRISTINE)"
}
# finally: always restore, even if a shot dies mid-run.
trap restore EXIT

run_panel() {  # $1 = tag
  out="${LOG}-$1.log"
  ( cd "$BENCH" && go test ./internal/panel/ -count=1 -timeout 30m -v ) > "$out" 2>&1
  echo "tag=$1 rc=$?"
  echo "  RUN=$(grep -c '^=== RUN' "$out" || true) PASS_TOP=$(grep -c '^--- PASS' "$out" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true) SKIP=$(grep -c '^--- SKIP' "$out" || true)"
}

echo "###### SHOT 1: silently ADD U+206B to the product roster ######"
# content anchor, NOT an absolute line: insert 0x206B into the case list.
( cd "$BENCH" && sed -i 's/0x2060, 0xFEFF,/0x2060, 0x206B, 0xFEFF,/' "$SCRUB" )
echo "--- mutated line ---"; grep -n '0x206B' "$BENCH/$SCRUB"
run_panel shot1
echo "--- red roster ---"; grep '^--- FAIL' "${LOG}-shot1.log"
echo "--- my nail red sentences (file:line verbatim) ---"
grep -E 'display_scrub_327r2_test\.go:[0-9]+:' "${LOG}-shot1.log" | grep -iE 'never named|no longer washed|cardinality|U\+206B' | head -20
restore
