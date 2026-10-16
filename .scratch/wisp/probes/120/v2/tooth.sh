#!/bin/sh
# Leg 120-v2 (adjudicator, read-only for product code).
# Dejaw controls for the two inline re-checks of sealHandle, plus a third control on
# the NEAR-IDENTICAL pair in sameObjectPinned, so that "which pair did I hit" is a
# reading and not an assumption.
#
# Platform: out-of-repo export tree built by
#   git archive 232f309e go.mod go.sum internal cmd tools scripts | tar -x -C /d/tmp/probe120v2/head
# Mutations are physical byte rewrites in that tree (no -overlay). The repository work
# tree is never touched.
V=/d/tmp/probe120v2/head
F="$V/internal/winsec/winsec_other.go"
LOG=/d/tmp/probe120v2/v2-tooth-driver.log
BAK="$F.bak-120v2-teeth-on"

: > "$LOG"
say() { echo "$@" >> "$LOG"; }

sha_of() { sha256sum "$1" | cut -d' ' -f1; }

CLEAN_SHA=$(sha_of "$F")
say "clean sha256 = $CLEAN_SHA"
cp "$F" "$BAK"
say "backup kept (never deleted) = $BAK"

# body_range NAME -> prints "start end" for a top-level func
body_range() {
  local s e
  s=$(grep -n "^func $1" "$F" | cut -d: -f1 | head -1)
  e=$(awk -v s="$s" 'NR>s && /^}/ {print NR; exit}' "$F")
  echo "$s $e"
}

# find_lines START END PATTERN -> line numbers in range matching the fixed string
find_lines() {
  awk -v s="$1" -v e="$2" -v p="$3" 'NR>=s && NR<=e && index($0,p) {print NR}' "$F"
}

# neuter LINEIDX -> rewrite that one line, inserting "false && " after the leading tab+if
neuter() {
  local n="$1" out=/d/tmp/probe120v2/neuter.out
  awk -v n="$n" 'BEGIN{OFS=""} NR==n{
      if (substr($0,1,5) != "\tif w") {
        # expect "\tif <cond> {"
      }
      if (substr($0,1,4) == "\tif ") { $0 = "\tif false && " substr($0,5) }
      else { print "NEUTER_FAILED: line " n " does not start with tab-if" > "/dev/stderr"; exit 3 }
    }
    {print}' "$F" > "$out" && mv "$out" "$F"
}

run_pkg() {
  local tag="$1" rc
  rc=$(MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -e GOPROXY=off \
      -v /d/work/base/gopath/pkg/mod:/gomod -e GOMODCACHE=/gomod \
      -v "$V":/work -v probe120v2build:/root/.cache/go-build -w /work golang:1.27 \
      sh -c "go test -count=1 -v ./internal/winsec/ > /work/$tag.log 2>&1; echo \$?" 2>>"$LOG")
  say "$tag go-test rc=$rc"
  say "$tag RUN=$(grep -c '^=== RUN' "$V/$tag.log") TOPPASS=$(grep -c '^--- PASS' "$V/$tag.log") TOPFAIL=$(grep -c '^--- FAIL' "$V/$tag.log") TOPSKIP=$(grep -c '^--- SKIP' "$V/$tag.log")"
  grep '^--- FAIL' "$V/$tag.log" | sed "s/^/$tag FAILED /" >> "$LOG"
}

restore() {
  cp "$BAK" "$F"
  local now; now=$(sha_of "$F")
  if [ "$now" = "$CLEAN_SHA" ]; then say "RESTORE MATCHES ($now)"; else say "RESTORE MISMATCH clean=$CLEAN_SHA now=$now"; fi
}

HR=$(body_range sealHandle)
SR=$(body_range sameObjectPinned)
say "sealHandle body range = $HR   sameObjectPinned body range = $SR"
say "whole-file occurrences of the link predicate: $(grep -n 'ModeSymlink != 0' "$F" | tr '\n' ' ')"
say "whole-file occurrences of the identity predicate: $(grep -n '!os.SameFile(pinned, now)' "$F" | tr '\n' ' ')"

# ---------- baseline (teeth on) ----------
run_pkg v2-teeth-on

# ---------- control 1: sealHandle link face only ----------
say "--- control 1: neuter the link face INSIDE sealHandle ---"
L1=$(find_lines $HR 'now.Mode()&os.ModeSymlink != 0')
say "sealHandle link-face line(s) in range: [$L1]  (count=$(echo $L1 | wc -w))"
say "before: $(sed -n "${L1}p" "$F")"
neuter "$L1"
say "after : $(sed -n "${L1}p" "$F")"
say "LANDS: false && in file at lines: $(grep -n 'false &&' "$F" | tr '\n' ' ')  count=$(grep -c 'false &&' "$F")"
say "mutated sha256 = $(sha_of "$F")"
run_pkg v2-tooth-handle-link
restore

# ---------- control 2: sealHandle SameFile face only ----------
say "--- control 2: neuter the SameFile face INSIDE sealHandle ---"
L2=$(find_lines $HR '!os.SameFile(pinned, now)')
say "sealHandle SameFile-face line(s) in range: [$L2]  (count=$(echo $L2 | wc -w))"
say "before: $(sed -n "${L2}p" "$F")"
neuter "$L2"
say "after : $(sed -n "${L2}p" "$F")"
say "LANDS: false && in file at lines: $(grep -n 'false &&' "$F" | tr '\n' ' ')  count=$(grep -c 'false &&' "$F")"
say "mutated sha256 = $(sha_of "$F")"
run_pkg v2-tooth-handle-samefile
restore

# ---------- control 3: the near-identical pair in sameObjectPinned ----------
say "--- control 3: neuter BOTH faces INSIDE sameObjectPinned (the other pair) ---"
M1=$(find_lines $SR 'now.Mode()&os.ModeSymlink != 0')
M2=$(find_lines $SR '!os.SameFile(pinned, now)')
say "sameObjectPinned faces in range: [$M1 $M2]  (count=$(( $(echo $M1 | wc -w) + $(echo $M2 | wc -w) )))"
for m in $M1 $M2; do
  say "before line $m: $(sed -n "${m}p" "$F")"
  neuter "$m"
  say "after  line $m: $(sed -n "${m}p" "$F")"
done
say "LANDS: false && in file at lines: $(grep -n 'false &&' "$F" | tr '\n' ' ')  count=$(grep -c 'false &&' "$F")"
say "mutated sha256 = $(sha_of "$F")"
run_pkg v2-tooth-walk-pair
restore

# ---------- re-confirm baseline colour after every restore ----------
run_pkg v2-restored-reconfirm
say "final file sha256 = $(sha_of "$F")   clean = $CLEAN_SHA"
say "driver done"
