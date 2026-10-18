#!/bin/bash
# 341-r1: precision pass on two readings that my own first census left ambiguous.
#   (a) P3's second anchor hit 1 - name WHAT that hit is, so "P3 = nothing" is not resting on a fuzzy count.
#   (b) AC#3's overclaim sentence - the ticket says "sec 1 与 sec 5 那句"; measure how many times it actually occurs and where.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
A=98dbf20a35027a9290f52649d8ab2634646cb14e
H=$(git rev-parse --short HEAD)
out=.scratch/wisp/probes/341/r1/logs/13-precision.txt
{
  echo "HEAD at measurement = $H ; census anchor = $A (short $(git rev-parse --short $A))"
  echo
  echo "(a) the single hit of the P3 'multi-line block containing [ball]' anchor:"
  git grep -nE '"\[ball\]\\n' "$A" -- 'internal/config/*_test.go' 'cmd/wisp/*_test.go'
  echo "  => is it a DELETION of a section or something else? Read the surrounding lines:"
  git show "$A:cmd/wisp/config_reload_223_test.go" | sed -n '348,356p'
  echo "  (that is r.plant(...) = a test WRITING a [ball] section, not a guarded-write delete shape)"
  echo "  so P3 as 'no resident case deletes a whole section' stands: anchor 1 (a Replace whose target starts"
  echo "  with the header) = 0 hits, and its same-form positive control (P1's size-line delete) = 1 hit."
  echo
  echo "(b) AC#3: occurrences of the overclaiming phrase in the ticket-226 verdict piece:"
  for f in docs/evidence/s1/226-config-write-no-clobber-r1.md docs/evidence/s1/226-config-write-no-clobber-v2.md; do
    echo "  ---- $f ----"
    echo "    count(不白重排)  = $(grep -c '不白重排' "$f" || true)"
    echo "    count(不白丢注释) = $(grep -c '不白丢注释' "$f" || true)"
    grep -n '不白' "$f" | sed 's/^/      /'
  done
  echo
  echo "  which section is the hit in? headers of r1.md with line numbers:"
  grep -n '^## ' docs/evidence/s1/226-config-write-no-clobber-r1.md | sed 's/^/      /'
  echo
  echo "  the sentence's own predicate, verbatim (the case it describes), at the current HEAD:"
  grep -n 'func TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' internal/config/writeguard_226_test.go
  echo
  echo "  is that case GREEN today on the bench (baseline tree, FAIL_TOP=0 run)? name it from the log:"
  grep -E '^--- (PASS|FAIL): TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' \
    .scratch/wisp/probes/341/r1/logs/05-color-Adflt-base.txt .scratch/wisp/probes/341/r1/logs/05-color-Bdtemp-base.txt
  echo
  echo "  what the production log sentence says NOW (227's narrowing), verbatim:"
  sed -n '152,157p' internal/config/writeguard.go
  echo
  echo "  ticket 341's own words for the target of AC#3 (its 现量 row 3 + the AC#3 box):"
  grep -n 're-laid in canonical form' .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md | sed 's/^/      /'
  grep -n '不白丢' .scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md | sed 's/^/      /'
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
