#!/usr/bin/env bash
# ci-red-attribution.sh - ticket 346 AC#3: turn a CI red roster into a per-name
# "which ticket owns this" reading, so no leg has to guess whether a red is new.
#
# INPUT   one or more GitHub Actions run ids, OR pre-fetched logs via --log.
# OUTPUT  per unique top-level FAIL name: jobs / criteria-line owner / file owner /
#         prose mentions / ledger ids / REGISTERED or UNREGISTERED. Only rows with
#         no criteria, no named file and no ledger id are UNREGISTERED: those are
#         the only rows a new ticket may be opened for.
#
# Read-only bench piece. No CI step (.github/** = human approval), no standing
# assertion. It writes only inside a temp dir.
#
# Guards against the two failure modes this ticket was opened for:
#   1. gh prefixes every log line with "job<TAB>step<TAB>timestamp ", so the
#      ticket-face ruler (^--- FAIL at column 0) reads 0 hits on a red log.
#      This ruler strips the prefix first, then anchors ^---; indented subtests
#      and mid-line prose stay out of the denominator.
#   2. a zero-row roster is REFUSED (exit 5), never printed as "0 red names".
#      A ruler that measured nothing may not hand back a clean roster.
# The positive control runs before attribution: if a name we KNOW is registered
# cannot be found, the ruler exits 1 instead of attributing a roster it never read.
#
# Usage:  bash ci-red-attribution.sh 39101060753 39124122121
#         bash ci-red-attribution.sh --log a.log --log b.log
#         bash ci-red-attribution.sh --log a.log --issues DIR --ledger FILE
# Options: --issues DIR (default .scratch/wisp/issues)
#          --ledger FILE (default docs/reports/pending-and-issues.md)
#          --func-index FILE (optional "Name<TAB>path" map; default: git grep)
#          --posctl NAME (default TestC26RewriteAccountIsConsumedAtEverySecurityLeg)
# Run from the repository root. Exit: 0 ok / 1 positive control / 2 usage /
# 3 gh fetch / 4 unreadable log / 5 empty roster.
set -uo pipefail

ISSUES_DIR="${ISSUES_DIR:-.scratch/wisp/issues}"
LEDGER="${LEDGER:-docs/reports/pending-and-issues.md}"
POSCTL="${POSCTL:-TestC26RewriteAccountIsConsumedAtEverySecurityLeg}"
FUNC_INDEX="${FUNC_INDEX:-}"
logs=()
runids=()

while [ $# -gt 0 ]; do
  case "$1" in
    --log)
      shift
      [ $# -gt 0 ] || { echo "ci-red-attrib: --log needs a path" >&2; exit 2; }
      logs+=("$1"); shift ;;
    --issues)
      shift
      [ $# -gt 0 ] || { echo "ci-red-attrib: --issues needs a dir" >&2; exit 2; }
      ISSUES_DIR="$1"; shift ;;
    --ledger)
      shift
      [ $# -gt 0 ] || { echo "ci-red-attrib: --ledger needs a path" >&2; exit 2; }
      LEDGER="$1"; shift ;;
    --func-index)
      shift
      [ $# -gt 0 ] || { echo "ci-red-attrib: --func-index needs a path" >&2; exit 2; }
      FUNC_INDEX="$1"; shift ;;
    --posctl)
      shift
      [ $# -gt 0 ] || { echo "ci-red-attrib: --posctl needs a name" >&2; exit 2; }
      POSCTL="$1"; shift ;;
    -h | --help)
      sed -n '2,31p' "$0"; exit 0 ;;
    -*)
      echo "ci-red-attrib: unknown option $1" >&2; exit 2 ;;
    *)
      runids+=("$1"); shift ;;
  esac
done

if [ ${#logs[@]} -eq 0 ] && [ ${#runids[@]} -eq 0 ]; then
  echo "ci-red-attrib: no run id and no --log; nothing to measure (not a green)" >&2
  exit 2
fi
[ -d "$ISSUES_DIR" ] || { echo "ci-red-attrib: issue pool not a directory: $ISSUES_DIR" >&2; exit 1; }
[ -f "$LEDGER" ] || { echo "ci-red-attrib: ledger not a file: $LEDGER" >&2; exit 1; }

tmpd="$(mktemp -d 2>/dev/null || echo "./.ci-red-attrib.$$")"
mkdir -p "$tmpd"
trap 'rm -rf "$tmpd"' EXIT

# Stage 1: fetch run ids. A gh failure is a reading, not a green.
for r in ${runids[@]+"${runids[@]}"}; do
  dst="$tmpd/run-$r.log"
  if ! command -v gh >/dev/null 2>&1; then
    echo "ci-red-attrib: gh not on PATH; pass --log instead" >&2; exit 3
  fi
  if ! gh run view "$r" --log-failed >"$dst" 2>"$tmpd/err-$r"; then
    echo "ci-red-attrib: gh run view $r FAILED -> take the run id off the ticket face," >&2
    echo "  read the roster there and MARK IT 未重取 (a reading, not a green)." >&2
    sed -n '1,3p' "$tmpd/err-$r" >&2
    exit 3
  fi
  logs+=("$dst")
done

for L in "${logs[@]}"; do
  if [ ! -s "$L" ]; then
    echo "ci-red-attrib: unreadable log $L (missing or zero bytes)" >&2
    exit 4
  fi
done

# Stage 2: strip the gh "job<TAB>step<TAB>timestamp " prefix, THEN anchor ^---.
# An indented subtest ("    --- FAIL:") and mid-line prose ("note: --- FAIL:")
# both fail ^--- after the strip, so they never enter the denominator.
for L in "${logs[@]}"; do
  awk -F '\t' 'NF >= 3 {
    payload = $3
    sub(/^[0-9-]+T[0-9:.]+Z /, "", payload)
    if (payload ~ /^--- (PASS|FAIL|SKIP): [A-Za-z0-9_.\/-]+( |$)/) {
      split(payload, a, " ")
      verdict = a[2]; sub(/:$/, "", verdict)
      print $1 "\t" verdict "\t" a[3]
    }
  }' "$L"
done >"$tmpd/all.tsv"

awk -F '\t' '$2 == "FAIL" { print $1 "\t" $3 }' "$tmpd/all.tsv" | sort -u >"$tmpd/roster.tsv"
rows="$(wc -l <"$tmpd/roster.tsv" | tr -d ' ')"
if [ "$rows" -eq 0 ]; then
  echo "ci-red-attrib: REFUSING to report a roster: 0 top-level FAIL rows in ${#logs[@]} log(s)." >&2
  echo "  A roster nobody measured is not a green. Check the log shape before reading anything into this." >&2
  exit 5
fi
cut -f2 "$tmpd/roster.tsv" | sort -u >"$tmpd/names.txt"

echo "ci-red-attrib: logs=${#logs[@]}  顶格 FAIL (job x name)=$rows"

# Stage 3: positive control BEFORE attribution. A ruler that cannot find a name
# we know is registered has no business reporting on the rest.
if ! grep -rqF -- "$POSCTL" "$ISSUES_DIR" "$LEDGER" 2>/dev/null; then
  echo "ci-red-attrib: POSITIVE CONTROL FAILED - $POSCTL is not in the pool or ledger;" >&2
  echo "  the ruler is blind (or the pool path is wrong). Refusing to attribute." >&2
  exit 1
fi
echo "ci-red-attrib: POSITIVE CONTROL ok ($POSCTL resolves to a ticket file or ledger entry)"

# Stage 4: one criteria index (file -> ticket id), then per-name lookups.
grep -HE '^- \[[ x]\]' "$ISSUES_DIR"/*-*.md 2>/dev/null |
  sed -E 's#^([^:]*):#\1\t#' |
  awk -F '\t' '{ p = $1; sub(/^.*\//, "", p); sub(/-.*/, "", p); print p "\t" $2 }' \
    >"$tmpd/crit.tsv"

echo "---- per-name attribution ----"
unreg=0
posctl_status=""
while IFS= read -r name; do
  [ -n "$name" ] || continue
  jobs="$(awk -F '\t' -v n="$name" '$2 == n { print $1 }' "$tmpd/roster.tsv" | sort -u | paste -sd, -)"
  own="$(awk -F '\t' -v n="$name" 'index($2, n) { print $1 }' "$tmpd/crit.tsv" | sort -u |
    sed 's/^0*//' | sort -un | paste -sd, -)"

  file=""
  if [ -n "$FUNC_INDEX" ]; then
    file="$(awk -F '\t' -v n="$name" '$1 == n { print $2; exit }' "$FUNC_INDEX" 2>/dev/null)"
  else
    file="$(git grep -l -E "^func ${name}(\(| )" -- '*_test.go' 2>/dev/null | head -1)"
  fi
  base=""
  own_file=""
  if [ -n "$file" ]; then
    base="$(basename "$file")"
    own_file="$(grep -rlF -- "$base" "$ISSUES_DIR" 2>/dev/null |
      sed -n 's#.*/\([0-9]\{2,4\}\)-[^/]*\.md#\1#p' | sed 's/^0*//' | sort -un | paste -sd, -)"
  fi

  ment="$(grep -rlF -- "$name" "$ISSUES_DIR" 2>/dev/null | wc -l | tr -d ' ')"
  led="$(grep -hF -- "$name" "$LEDGER" 2>/dev/null |
    grep -oE '\b[ARQ][0-9]{2,}\b' | sort -u | paste -sd, -)"

  if [ -n "$own" ] || [ -n "$own_file" ] || [ -n "$led" ]; then
    st=REGISTERED
  else
    st=UNREGISTERED
    unreg=$((unreg + 1))
  fi
  [ "$name" = "$POSCTL" ] && posctl_status="$st"

  printf 'name=%s\n' "$name"
  printf '  jobs=%s  file=%s\n' "${jobs:-?}" "${base:-未解析}"
  printf '  判据在册票=%s  具名文件在册票=%s  散文提及=%s枚票  台账=%s\n' \
    "${own:-—}" "${own_file:-—}" "${ment:-0}" "${led:-—}"
  printf '  => %s\n' "$st"
done <"$tmpd/names.txt"

if grep -qxF "$POSCTL" "$tmpd/names.txt" && [ "$posctl_status" != "REGISTERED" ]; then
  echo "ci-red-attrib: POSITIVE CONTROL FAILED - $POSCTL is on the roster but the resolver read it as $posctl_status" >&2
  exit 1
fi

echo "---- summary ----"
printf 'unique red names=%s  UNREGISTERED=%s\n' "$(wc -l <"$tmpd/names.txt" | tr -d ' ')" "$unreg"
if [ "$unreg" -eq 0 ]; then
  echo "NOTE: every red name resolves to a ticket; no new ticket may be opened from this roster."
else
  echo "NOTE: only the UNREGISTERED rows may open a new ticket."
fi
exit 0
