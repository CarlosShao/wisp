#!/usr/bin/env bash
# 139 AC#3 mutation battery. Each mutation is applied to the working tree,
# measured with a named -run, then reverted from probes/139/mut/*.orig.
# Nothing here relaxes an assertion or skips a test: a mutation that does not
# turn anything red is reported as such, which is the point of the exercise.
set -u
cd "$(git rev-parse --show-toplevel)" || exit 1
[ -d internal/agent ] || { echo "not at repo root: $PWD"; exit 1; }
P=.scratch/wisp/probes/139
OUT=$P/mutations.log
: > "$OUT"

restore() { cp "$P/mut/compress.go.orig" internal/agent/compress.go; cp "$P/mut/loop.go.orig" internal/agent/loop.go; }

verify_clean() { # proves the tree is back to what was committed-as-written
  cmp -s internal/agent/compress.go "$P/mut/compress.go.orig" || { echo "RESTORE FAILED compress.go" | tee -a "$OUT"; exit 1; }
  cmp -s internal/agent/loop.go "$P/mut/loop.go.orig" || { echo "RESTORE FAILED loop.go" | tee -a "$OUT"; exit 1; }
}

run_case() { # $1 label, $2 -run regex
  echo "" | tee -a "$OUT"
  echo "### $1  (go test -run '$2')" | tee -a "$OUT"
  echo '$ go build ./internal/agent/ ; echo rc=$?' | tee -a "$OUT"
  go build ./internal/agent/ 2>&1 | tail -3 | tee -a "$OUT"
  echo "build_rc=${PIPESTATUS[0]}" | tee -a "$OUT"
  echo '$ go test ./internal/agent/ -count=1 -v -run '"'$2'" | tee -a "$OUT"
  go test ./internal/agent/ -count=1 -v -run "$2" > "$P/mut/$3.log" 2>&1
  rc=$?
  echo "test_rc=$rc" | tee -a "$OUT"
  echo "--- red lines (--- FAIL only) ---" | tee -a "$OUT"
  grep -E '^\s*--- (FAIL|PASS):' "$P/mut/$3.log" | tee -a "$OUT"
  echo "--- failure bodies ---" | tee -a "$OUT"
  grep -vE '^\s*(===|--- (PASS|FAIL))' "$P/mut/$3.log" | grep -E 'compress_trace_test\.go|FAIL|ok ' | head -12 | tee -a "$OUT"
}

echo "# 139 mutation battery — $(date -Iminutes) HEAD=$(git rev-parse --short HEAD)" > "$OUT"

echo "### [M0] baseline: trace IN place, all four new tests" >> "$OUT"
run_case "[M0] unmutated" 'TestCompressionTrace' m0

# M1: remove the trace emission entirely (this is AC#3's named mutation, and it
# is behaviourally identical to the pre-fix code for these four tests).
restore
perl -0pi -e 's/\tif rep\.Ran \{.*?\n\t\}\n\treturn out, rep, nil/\treturn out, rep, nil/s' internal/agent/compress.go
echo "" >> "$OUT"; echo '$ diff of mutation M1 (trace emission deleted):' >> "$OUT"
diff -u "$P/mut/compress.go.orig" internal/agent/compress.go | head -40 >> "$OUT"
if grep -q 'c.log().Info' internal/agent/compress.go; then echo "M1 DID NOT LAND" | tee -a "$OUT"; restore; exit 1; fi
run_case "[M1] emission deleted (= pre-fix behaviour)" 'TestCompressionTrace' m1

# M2: keep the emission, drop the loop's wiring of the host logger.
restore
perl -0pi -e 's/NewCompressor\(b, opt\.Summarizer, WithLogger\(opt\.Logger\)\)/NewCompressor(b, opt.Summarizer)/' internal/agent/loop.go
echo "" >> "$OUT"; echo '$ diff of mutation M2 (loop no longer hands its logger down):' >> "$OUT"
diff -u "$P/mut/loop.go.orig" internal/agent/loop.go | head -20 >> "$OUT"
if grep -q 'WithLogger(opt.Logger)' internal/agent/loop.go; then echo "M2 DID NOT LAND" | tee -a "$OUT"; restore; exit 1; fi
run_case "[M2] wiring deleted" 'TestCompressionTrace' m2

# M3: make the emission unconditional - the tautological-judgement shape.
restore
perl -0pi -e 's/\tif rep\.Ran \{\n\t\t\/\/ The one trace/\tif true {\n\t\t\/\/ The one trace/' internal/agent/compress.go
echo "" >> "$OUT"; echo '$ diff of mutation M3 (emitted even when nothing folded):' >> "$OUT"
diff -u "$P/mut/compress.go.orig" internal/agent/compress.go | head -20 >> "$OUT"
if grep -q 'if rep.Ran {' internal/agent/compress.go; then echo "M3 DID NOT LAND" | tee -a "$OUT"; restore; exit 1; fi
run_case "[M3] unconditional emission" 'TestCompressionTrace' m3

# M4: leak one history part into the record - the AC#2 privacy rule.
restore
perl -0pi -e 's/"history_changed", rep\.TokensAfter != rep\.TokensBefore \|\| len\(out\) != len\(hist\)\)/"history_changed", rep.TokensAfter != rep.TokensBefore || len(out) != len(hist), "peek", fmt.Sprint(out[0].Content[0]))/ or die "M4 no match\n"' internal/agent/compress.go
echo "" >> "$OUT"; echo '$ diff of mutation M4 (a history part joined the record):' >> "$OUT"
diff -u "$P/mut/compress.go.orig" internal/agent/compress.go | head -20 >> "$OUT"
run_case "[M4] history content leaked into the record" 'TestCompressionTrace' m4

restore
verify_clean
echo "" >> "$OUT"
echo "### tree restored to pre-mutation state: $(git status --porcelain -- internal/agent | tr '\n' ' ')" >> "$OUT"
echo "ALL MUTATIONS DONE"
