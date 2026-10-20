#!/bin/sh
# 342-r2 final sweep: re-take the gates on the FINAL HEAD, re-check the write-face mutex, and print every
# number that a commit body is about to claim (claims are written only from this output).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
echo "HEAD=$(git log -1 --format='%H %h %s' | cut -c1-60)"
unset GOTMPDIR TMPDIR TEMP TMP
echo "-- gate 1 gofumpt (my two literal paths) --"
/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go
echo "GOFUMPT_RC=$?"
echo "-- gate 2 d22scan --"
sh scripts/d22scan.sh > /d/tmp/wisp/20261019/342r2/logs/d22scan-final2.log 2>&1
echo "D22_RC=$?"
echo "-- gate 3 path-length budget --"
sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261019/342r2/logs/lenbudget-final2.log 2>&1
echo "LEN_RC=$?"
echo "-- gate 4 go vet (native / linux / darwin, files NOT paired) --"
go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-native2.log 2>&1
echo "VET_NATIVE_RC=$? BYTES=$(wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-native2.log)"
GOOS=linux go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-linux2.log 2>&1
echo "VET_LINUX_RC=$? BYTES=$(wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-linux2.log)"
GOOS=darwin go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-darwin2.log 2>&1
echo "VET_DARWIN_RC=$? BYTES=$(wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-darwin2.log)"
echo "-- write-face mutex: every commit from my anchor that touched internal/winsec --"
git log --format='%h%x09%s' 3abbb6391a8a497c0332e5363f24c8a250a7eb78..HEAD -- internal/winsec | cut -c1-70
echo "-- my leg's roster (subject-anchored, no assumed separator) --"
git log --format='%h%x09%s' | grep -P '^\S+\t票 342＋票 348 落地腿 342-r2' | wc -l
echo "-- numbers a commit body may claim --"
echo "REPORT_MD_LINES=$(wc -l < .scratch/wisp/probes/342/r2/report.md)"
echo "PROBE_ROOT_FILES=$(find .scratch/wisp/probes/342/r2 -type f | wc -l)"
echo "PROBE_ROOT_ZERO_BYTE=$(find .scratch/wisp/probes/342/r2 -type f -size 0 | wc -l)"
echo "PROBE_ROOT_OUT_EXT=$(find .scratch/wisp/probes/342/r2 -name '*.out' | wc -l)"
echo "TICKET_342_MY_BULLETS=$(grep -c 'agent=.342-r2.' .scratch/wisp/issues/342-*.md)"
echo "TICKET_348_MY_BULLETS=$(grep -c 'agent=.342-r2.' .scratch/wisp/issues/348-*.md)"
echo "TICKET_342_BOXES=$(grep -c '^\- \[x\]' .scratch/wisp/issues/342-*.md)/$(grep -c '^\- \[ \]' .scratch/wisp/issues/342-*.md)"
echo "TICKET_348_BOXES=$(grep -c '^\- \[x\]' .scratch/wisp/issues/348-*.md)/$(grep -c '^\- \[ \]' .scratch/wisp/issues/348-*.md)"
echo "GO_FILES_MINE=2"
for f in internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go; do
  echo "BLOB_CR $f = $(git cat-file blob "HEAD:$f" | tr -cd '\r' | wc -c)"
done
echo "-- codepoint ruler on authored files --"
perl .scratch/wisp/probes/342/r2/sh/codepoint-ruler.pl .scratch/wisp/probes/342/r2/report.md \
  .scratch/wisp/probes/342/r2/logs/bench-probes-verbatim.txt .scratch/wisp/probes/342/r2/logs/gates-and-roots.txt
echo "-- worktree dirt under my write face (must be empty) --"
git status --porcelain -- internal/winsec .scratch/wisp/probes/342/r2 .scratch/wisp/issues/342-*.md .scratch/wisp/issues/348-*.md
