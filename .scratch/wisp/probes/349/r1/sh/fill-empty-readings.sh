#!/bin/sh
# 公共块「⛔ 0 字节文件（落笔前后 wc -c）」＝逐枚非 0。这些件是**空读数**（0 命中／无 stderr／无告警），
# 空本身就是那条读数 ⇒ 就地写明它是什么的读数＋正控指向哪一枚件，⛔ 删（临时件只建不删）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
R=.scratch/wisp/probes/349/r1
fix () { # $1=path $2=说明
  if [ -f "$1" ] && [ ! -s "$1" ]; then
    printf '%s\n' "$2" > "$1"
    printf 'filled=%s bytes=%s\n' "$1" "$(wc -c < "$1")"
  else
    printf 'skip(non-empty or absent)=%s bytes=%s\n' "$1" "$( [ -f "$1" ] && wc -c < "$1" || echo ABSENT )"
  fi
}
fix "$R/logs/census/ac0-base-full.err"   "EMPTY READING = gauge stderr 0 bytes = no PARSE_FAIL line (whole-package census, START side). Positive control for this gauge shape = logs/census/ac0-base-full.txt (non-zero, 11498 bytes)."
fix "$R/logs/census/ac0-base-prod.err"   "EMPTY READING = gauge stderr 0 bytes = no PARSE_FAIL line (three production files, START side). Positive control = logs/census/ac0-base-prodonly.txt (non-zero)."
fix "$R/logs/census/ac0-del-full.err"    "EMPTY READING = gauge stderr 0 bytes = no PARSE_FAIL line (whole-package census, DELIVERY side). Positive control = logs/census/ac0-del-full.txt (non-zero)."
fix "$R/logs/census/ac0-del-prod.err"    "EMPTY READING = gauge stderr 0 bytes = no PARSE_FAIL line (three production files, DELIVERY side). Positive control = logs/census/ac0-del-prodonly.txt (non-zero)."
fix "$R/logs/census/ac0-del-four-names.txt" "ZERO-HITS READING = DELIVERY side: none of the four ticket names (fires()/taskScopeID/disposalScope/bare Complete()) appears in the census output any more. This IS the wanted reading. Positive control = logs/census/ac0-base-four-names.txt (3 lines, START side)."
fix "$R/logs/gates/g1-before.txt"        "ZERO-LINES READING = gofumpt -l on the START-side blobs (LF, git object form) printed 0 lines = my byte-side start was already gofumpt-clean. Positive control = the same gauge on provenance.go WORKING COPY (CRLF 存量) which does print it: see logs/gates/gates-outer.log."
fix "$R/logs/gates/g1-after.txt"         "ZERO-LINES READING = gofumpt -l on the DELIVERY blobs (LF, git object form) printed 0 lines = my new bytes are gofumpt-clean. Positive control: see logs/gates/gates-outer.log (gofumpt -l on the working copies prints provenance.go = pre-existing CRLF, not my byte)."
fix "$R/logs/gates/g1b-gofumpt-startblobs.txt" "ZERO-LINES READING = gofumpt -l over the three START blobs (extracted to disk as LF) printed nothing = start side already clean; the CRLF working-copy hit is a pre-existing condition, not introduced by this leg."
fix "$R/logs/gates/g1c-gofumpt-delblobs.txt"   "ZERO-LINES READING = gofumpt -l over the three DELIVERY blobs (git object form, LF) printed nothing = my new bytes clean."
fix "$R/logs/gates/g4.txt"               "ZERO-BYTES READING = go vet ./internal/risk/ rc=0 with NO output (go vet prints nothing when clean). rc= is recorded in logs/gates/gates-outer.log; positive control for a real vet finding is not needed here because rc=0 and empty = the passing shape."
echo "== re-count zero-byte files (want none) =="
find "$R" -type f -size 0 | sort
echo "(nothing listed above = rule satisfied)"
