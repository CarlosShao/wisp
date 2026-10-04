#!/usr/bin/env bash
# 212-v3 plant battery: AC#4 positive control + the exemption's real range.
# Runs the REAL scanner binary against a throwaway root OUTSIDE the repo
# (D:/tmp/212v3tree) so no tracked production file is ever touched.
TREE=/d/tmp/212v3tree
LOGS="/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/212/v3/logs"
PLANT="$TREE/internal/probe/zz_plant.go"
cd "/d/work/workspace/projects plans/Wisp/tools/d22scan" || exit 9

run_case() {
  local name="$1"
  go run . -root "D:/tmp/212v3tree" > "$LOGS/plant-$name.txt" 2>&1
  local rc=$?
  echo "rc=$rc" >> "$LOGS/plant-$name.txt"
  local hits
  hits=$(grep -c "phantom-citation" "$LOGS/plant-$name.txt")
  local toks
  toks=$(grep -o 'comment cites repo path "[^"]*"' "$LOGS/plant-$name.txt" | sed 's/^/    /')
  echo "CASE $name rc=$rc phantom_lines=$hits" >> "$LOGS/plant-summary.txt"
  if [ -n "$toks" ]; then echo "$toks" >> "$LOGS/plant-summary.txt"; fi
  rm -f "$PLANT"
}

: > "$LOGS/plant-summary.txt"
echo "battery run at $(date '+%Y-%m-%d %H:%M:%S%z')" >> "$LOGS/plant-summary.txt"

# P1 AC#4 ring: fully spelled, no mark anywhere, page absent
printf 'package probe\n\n// See docs/evidence/s1/212-v3-plant-phantom.md for the readings.\nfunc plant() {}\n' > "$PLANT"
run_case p1-fullspelled-ring

# P2 class 3 by design: "..." in the middle of a half-spelled path
printf 'package probe\n\n// See docs/evidence/s1/212-v3-...-phantom.md for the readings.\nfunc plant() {}\n' > "$PLANT"
run_case p2-middleshrink-silent

# P3 shape 2: complete phantom + trailing U+2026
printf 'package probe\n\n// See docs/evidence/s1/212-v3-tail-ellipsis.md\xe2\x80\xa6 for the readings.\nfunc plant() {}\n' > "$PLANT"
run_case p3-tail-ellipsis-silent

# P4 DISGUISE HUNT: markdown bold around a complete phantom path (two star marks)
printf 'package probe\n\n// **docs/evidence/s1/212-v3-bold-phantom.md** is the reading.\nfunc plant() {}\n' > "$PLANT"
run_case p4-markdown-bold

# P5 shape 3: CJK glue + a following abbreviation swallows the complete phantom
printf 'package probe\n\n// docs/evidence/s1/212-v3-glued.md\xe8\xa7\x81docs/evidence/s1/...\nfunc plant() {}\n' > "$PLANT"
run_case p5-cjk-glue-mark-silent

# P6 boundary: CJK glue but NO mark at all
printf 'package probe\n\n// \xe8\xa7\x81docs/evidence/s1/212-v3-glued-nomark.md\nfunc plant() {}\n' > "$PLANT"
run_case p6-cjk-glue-nomark-ring

# P7 existence-driven: the same P1 citation with the page now present on disk
mkdir -p "$TREE/docs/evidence/s1"
printf '# present\n' > "$TREE/docs/evidence/s1/212-v3-plant-phantom.md"
printf 'package probe\n\n// See docs/evidence/s1/212-v3-plant-phantom.md for the readings.\nfunc plant() {}\n' > "$PLANT"
run_case p7-page-exists-silent
rm -f "$TREE/docs/evidence/s1/212-v3-plant-phantom.md"

# P8 boundary: abbreviation on the NEXT comment line does not exempt line one
printf 'package probe\n\n// See docs/evidence/s1/212-v3-crossline.md for the readings.\n// or its shorthand docs/evidence/s1/212-v3-... -md.\nfunc plant() {}\n' > "$PLANT"
run_case p8-crossline-boundary

# P9 DISGUISE HUNT (ascii dots): complete phantom + trailing "..." = "and so on"
printf 'package probe\n\n// See docs/evidence/s1/212-v3-tail-dots.md... for details.\nfunc plant() {}\n' > "$PLANT"
run_case p9-tail-dots-ascii

# P10 control: plant removed entirely (battery leaves no plant behind)
run_case p10-plant-removed

echo "final tree state: $(ls "$TREE/internal/probe" | tr '\n' ' ')" >> "$LOGS/plant-summary.txt"
