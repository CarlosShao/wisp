#!/bin/bash
# 341-r1: P4/P3 carrier re-measure with correct anchors (the first gauge's P4 anchor was a ruler defect = 0 hits).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
A=98dbf20a35027a9290f52649d8ab2634646cb14e
out=.scratch/wisp/probes/341/r1/logs/09-census-p4-p3-remeasure.txt
{
  echo "[why a re-run] 01-start-gauges.txt gauge 3 used the anchor '\"\\[llm\\.providers\\.deepseek\\]' for P4 and hit 0."
  echo "  That is a RULER defect (my escaping), ⛔ a reading: 227-v1 measured the same格 with a different anchor and got"
  echo "  18 mentions / 1 real assertion. Per the dispatch (报 0 命中前先跑正控), re-measured below with the"
  echo "  assertion-shape anchor, and P3 got the same treatment so its 0 is not my escaping either."
  echo
  echo "[P4 carrier = the test that asserts a deleted provider table is NOT restored]"
  git grep -nE 'Providers\["deepseek"\]; ok' "$A" -- 'internal/config/*_test.go'
  echo "  (that one line is the assertion; writeguard_handedit_227_test.go:92 inside the case starting :70)"
  echo
  echo "[P4 total mentions of the provider table in *_test.go under internal/config + cmd/wisp]"
  echo "  hits=$(git grep -nE 'llm\.providers\.deepseek' "$A" -- 'internal/config/*_test.go' 'cmd/wisp/*_test.go' | wc -l)"
  echo "  of which delete-shape carriers = 1 (above)"
  echo
  echo "[P3 re-measure: does ANY *_test.go delete the [ball] HEADER (not merely mention it)?]"
  echo "  anchor 1: a Replace whose target string starts with the header   hits=$(git grep -nE 'Replace\([^,]*, "\[ball\]' "$A" -- 'internal/config/*_test.go' 'cmd/wisp/*_test.go' | wc -l)"
  echo "  anchor 2: any deletion of a multi-line block containing [ball]   hits=$(git grep -nE '"\[ball\]\\n' "$A" -- 'internal/config/*_test.go' 'cmd/wisp/*_test.go' | wc -l)"
  echo "  positive control for the SAME ruler style (P1, known to exist):   hits=$(git grep -nE '"size = 60  # my size\\n", ""' "$A" -- 'internal/config/*_test.go' | wc -l)"
  echo "  => P3's 0 is a reading, because the identical ruler form returns 1 for P1."
  echo
  echo "[P6 re-measure with the same discipline]"
  echo "  anchor: deletion of a base_url line                  hits=$(git grep -nE 'base_url[^"]*\\n", ""' "$A" -- 'internal/config/*_test.go' 'cmd/wisp/*_test.go' | wc -l)"
  echo "  positive control (same form, P5'):                   hits=$(git grep -nE '"schema_version = 2\\n", ""' "$A" -- 'internal/config/*_test.go' | wc -l)"
  echo
  echo "[P5 write-branch re-measure] a test that deletes schema_version AND still writes (owned key differs):"
  echo "  candidates = both places deleting that line:"
  git grep -nE '"schema_version = 2\\n", ""' "$A" -- 'internal/config/*_test.go'
  echo "  :110 = writeguard_handedit_227_test.go's nothing-written case (its AddAllowedDir is a no-op, log proves it)"
  echo "  => the write branch (b2 P5) had no carrier before this leg."
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
