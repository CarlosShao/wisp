#!/bin/sh
# 327-r2 shot 2 + the "removal/ablation" control. Export-tree ONLY; finally-restore
# from the named commit tree; reports restored hash (want 8050f22d...).
set -u
BENCH=/d/tmp/w327r2/bench
LOG=/d/tmp/w327r2/raw
SCRUB=internal/panel/display_scrub_216.go
TESTF=internal/panel/display_scrub_327r2_test.go
ANCHOR=7535d83b
REPO="/d/work/workspace/projects plans/Wisp"
PRISTINE=8050f22d14203bba006cbfabcfb2a748081286d492691109c3a35a7a435ff8fd
export GOTMPDIR=/d/tmp/w327r2/gotmp
export PATH="$BENCH/third_party/sherpa-onnx:$BENCH/build:$PATH"

restore_scrub() { ( cd "$BENCH" && git -C "$REPO" show "$ANCHOR:$SCRUB" > "$SCRUB" ); }
restore_test()  { ( cd "$BENCH" && git -C "$REPO" show HEAD:$TESTF > "$TESTF" 2>/dev/null || cp "$REPO/$TESTF" "$TESTF" ); }
finally() {
  restore_scrub; restore_test
  h=$(sha256sum "$BENCH/$SCRUB" | cut -d' ' -f1)
  echo "FINALLY scrub_sha256=$h"; [ "$h" = "$PRISTINE" ] || echo "FINALLY_MISMATCH want=$PRISTINE"
}
trap finally EXIT

run_panel() {  # $1 tag
  out="${LOG}-$1.log"
  ( cd "$BENCH" && go test ./internal/panel/ -count=1 -timeout 30m -v ) > "$out" 2>&1
  echo "tag=$1 RUN=$(grep -c '^=== RUN' "$out" || true) PASS_TOP=$(grep -c '^--- PASS' "$out" || true) FAIL_TOP=$(grep -c '^--- FAIL' "$out" || true) SKIP=$(grep -c '^--- SKIP' "$out" || true)"
}

echo "###### SHOT 2: silently REMOVE the existing U+00AD from the product roster ######"
restore_scrub
( cd "$BENCH" && sed -i 's/^\tcase 0x00AD,/\tcase/' "$SCRUB" )
echo "--- mutated case line (U+00AD gone) ---"; sed -n '164,172p' "$BENCH/$SCRUB" | sed 's/^/    /'
run_panel shot2
echo "--- my nail red sentences ---"; grep -E 'display_scrub_327r2_test\.go:[0-9]+:' "${LOG}-shot2.log" | grep -iE 'no longer washed|cardinality|U\+00AD' | head -8
echo "--- shot2 vs shot1 red-roster diff ---"
diff <(grep '^--- FAIL' "${LOG}-shot1.log" | sed 's/ (.*//') <(grep '^--- FAIL' "${LOG}-shot2.log" | sed 's/ (.*//')

echo; echo "###### ABLATION CONTROL: make the cardinality nail's ACTUAL side a copy of the EXPECTED roster (self-referential), keep shot 2's product mutation ######"
( cd "$BENCH" && sed -i 's/^\tactual := productWashedRunes()/\tactual := map[rune]bool{}; for _, r := range ruledRoster327r2 { actual[r] = true }\/\/ ABLATED/' "$TESTF" )
grep -n 'ABLATED' "$BENCH/$TESTF"
run_panel ablation
echo "--- ablation cardinality-nail status (want: no longer red, i.e. product-blind copy-vs-copy) ---"
grep -E '^--- (PASS|FAIL): TestDisplayWashRosterCardinality' "${LOG}-ablation.log"
finally; echo "RESTORED_OK"
