#!/bin/sh
# 342-r2 consolidated gate + root-classification evidence. Everything in the output file is taken by the
# commands themselves (no hand-written numbers).
cd "/d/work/workspace/projects plans/Wisp" || exit 1
# Absolute OUT: the brace group below cd's into the log directory and brace groups run in the SAME shell,
# so a relative path made the trailing `wc -l` fail (rc=1) even though every number in the log was
# correct. 本腿自抓第二枚同类(第一枚=utf8-and-codepoints.sh),都修在入库前,⛔ 带着一枚说不清的 rc=1 交件。
OUT="$PWD/.scratch/wisp/probes/342/r2/logs/gates-and-roots.txt"
{
  echo "== HEAD at gate time =="
  # NOTE (本腿自抓＝散件里的一处真损): this line used to be `git log -1 --format='%H %s' | cut -c1-90`,
  # and `cut -c` counts BYTES, so it cut a CJK commit subject in the middle of a multi-byte character
  # and left an undecodable byte in the evidence file (iconv -f UTF-8 -t UTF-8 reports it; a grep that
  # tolerates invalid bytes does not). Nothing after this point byte-cuts multi-byte text.
  git log -1 --format='%H %h' | sed 's/^/HEAD_FULL_SHA=/'
  git log -1 --format='%s' | sed 's/^/HEAD_SUBJECT=/'
  echo '(subject printed whole: cut -c counts BYTES and would split a multi-byte character. And NO'
  echo ' backticks inside a double-quoted echo string - that is command substitution: it once leaked'
  echo ' a whole git log into this evidence file, caught by reading the generated head back.)'
  echo
  echo "== gate 1: gofumpt -l, my two literal paths, run before unset too (same answer either way) =="
  unset GOTMPDIR TMPDIR TEMP TMP
  /d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go
  echo "GOFUMPT_RC=$? (empty listing above = clean)"
  echo
  echo "== gate 2: d22scan =="
  sh scripts/d22scan.sh > /d/tmp/wisp/20261019/342r2/logs/d22scan-consol.log 2>&1
  echo "D22_RC=$?"
  tail -1 /d/tmp/wisp/20261019/342r2/logs/d22scan-consol.log
  echo
  echo "== gate 3: path-length budget =="
  sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261019/342r2/logs/lenbudget-consol.log 2>&1
  echo "LEN_RC=$?"
  tail -1 /d/tmp/wisp/20261019/342r2/logs/lenbudget-consol.log
  echo
  echo "== gate 4: go vet, native + cross (cross files are NOT paired against the native one) =="
  go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-native.log 2>&1
  echo "VET_NATIVE_RC=$?"
  wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-native.log | sed 's/^/VET_NATIVE_OUTPUT_BYTES=/'
  GOOS=linux go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-linux.log 2>&1
  echo "VET_LINUX_RC=$?"
  wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-linux.log | sed 's/^/VET_LINUX_OUTPUT_BYTES=/'
  GOOS=darwin go vet ./internal/winsec/ > /d/tmp/wisp/20261019/342r2/logs/vet-darwin.log 2>&1
  echo "VET_DARWIN_RC=$?"
  wc -c < /d/tmp/wisp/20261019/342r2/logs/vet-darwin.log | sed 's/^/VET_DARWIN_OUTPUT_BYTES=/'
  echo "(a 0-byte vet log with rc=0 = no diagnostics, recorded as such instead of being skipped)"
  echo
  echo "== temp-root actually used by each reading (measured off the logs themselves) =="
  cd /d/tmp/wisp/20261019/342r2/logs || exit 1
  for f in land-targeted-1 base-targeted-1 land-full-1 base-full-1 land-full-2 base-full-2 \
           land-full-3 base-full-3 landminus-full-1 land-det-r1 land-det-r2 land-det-r3 \
           base-det-r1 base-det-r2 base-det-r3 land-det-shortroot base-det-shortroot \
           leaf-base-r1 leaf-base-r2 leaf-base-r3 leaf-land-v2-r1 leaf-land-v2-r2 leaf-land-v2-r3 \
           mutA-targeted-r1 mutA-targeted-r2 mutB-targeted-r1; do
    [ -f "$f.log" ] || continue
    droot=$(grep -c '342r2.gotmp' "$f.log" || true)
    sroot=$(grep -c 'AppData.Local.Temp' "$f.log" || true)
    proot=$(grep -c 'C:.Users.swq.Test' "$f.log" || true)
    echo "$f : D-gotmp-rows=$droot  LOCALAPPDATA-temp-rows=$sroot  USERPROFILE-direct-rows=$proot"
  done
  echo
  echo "== junction-construction health (A1053 = long D roots break mklink /J) =="
  grep -c 'mklink /J' land-det-shortroot.log base-det-shortroot.log leaf-base-r1.log 2>/dev/null || true
  grep -h 'is not a reparse point' land-det-shortroot.log base-det-shortroot.log leaf-*.log 2>/dev/null | sed 's/^/REPARSE-FAULT: /' || echo "REPARSE-FAULT-ROWS=0 (no construction failure in any of my logs)"
} > "$OUT" 2>&1
echo "WROTE $OUT"
wc -l "$OUT"
