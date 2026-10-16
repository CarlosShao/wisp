#!/bin/sh
# 120-v2, second control batch: the WALK pair, one face at a time.
# Why: tooth.sh control 3 neutered both faces of sameObjectPinned together, which
# proves the two pairs are distinguishable but not which case covers WHICH walk face.
# AC#3 must answer "is the shape 'the same object got swapped out' really constructed in
# this batch", so each walk face gets its own tooth here.
V=/d/tmp/probe120v2/head
F="$V/internal/winsec/winsec_other.go"
LOG=/d/tmp/probe120v2/v2-tooth2-driver.log
BAK="$F.bak-120v2-teeth-on"
: > "$LOG"
say() { echo "$@" >> "$LOG"; }
sha_of() { sha256sum "$1" | cut -d' ' -f1; }
CLEAN_SHA=$(sha_of "$F")

neuter() {
  local n="$1" out=/d/tmp/probe120v2/neuter2.out
  awk -v n="$n" 'NR==n{ if (substr($0,1,4) == "\tif ") { $0 = "\tif false && " substr($0,5) } else { print "NEUTER_FAILED line " n > "/dev/stderr"; exit 3 } } {print}' "$F" > "$out" && mv "$out" "$F"
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
restore() { cp "$BAK" "$F"; local now; now=$(sha_of "$F"); [ "$now" = "$CLEAN_SHA" ] && say "RESTORE MATCHES ($now)" || say "RESTORE MISMATCH clean=$CLEAN_SHA now=$now"; }

for spec in "219 walk-link-face" "222 walk-samefile-face"; do
  set -- $spec
  line="$1"; name="$2"
  say "--- control: neuter sameObjectPinned $name at line $line ---"
  say "before: $(sed -n "${line}p" "$F")"
  neuter "$line"
  say "after : $(sed -n "${line}p" "$F")"
  say "LANDS: false && lines: $(grep -n 'false &&' "$F" | tr '\n' ' ') count=$(grep -c 'false &&' "$F")"
  say "mutated sha256 = $(sha_of "$F")"
  run_pkg "v2-tooth-$name"
  restore
done
say "final sha256 = $(sha_of "$F") clean = $CLEAN_SHA"
