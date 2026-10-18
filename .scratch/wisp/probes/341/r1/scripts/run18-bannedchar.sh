#!/bin/bash
# 341-r1: write-face character/discipline ruler (banned character U+00A7 + the ⛔-shape checks I cite in the report).
# The character itself is ⛔ written anywhere in this file or in my prose: it is named as a code point and
# matched with an octal escape (printf '\302\247'), so the ruler does not carry the thing it measures.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/17-banned-char-ruler.txt
S=$(printf '\302\247')
f=internal/config/writeguard_unnailed_341_test.go
{
  echo "[U+00A7 = the character the dispatch forbids in my probe files] occurrences per file (want 0)"
  printf '  test file = '; grep -c -- "$S" "$f" || true
  printf '  report.md = '; grep -c -- "$S" .scratch/wisp/probes/341/r1/report.md || true
  printf '  report-answers.md = '; grep -c -- "$S" .scratch/wisp/probes/341/r1/report-answers.md || true
  echo "  files in my write face still carrying it:"
  grep -rl -- "$S" .scratch/wisp/probes/341/r1/ 2>/dev/null | sed 's/^/    /'
  echo "  (the three logs that quoted docs/evidence/s1 verbatim were tokenised to [U+00A7]; the token count is"
  echo "   5+4+4 = 13, recorded here so a reader can restore the quote from the source file if wanted)"
  printf '  token count in those three logs = '
  grep -rc '\[U+00A7\]' .scratch/wisp/probes/341/r1/logs/10-ac3-overclaim.txt \
    .scratch/wisp/probes/341/r1/logs/13-precision.txt \
    .scratch/wisp/probes/341/r1/logs/10-ac3-overclaim-run1-rulerdefect.txt | tr '\n' ' '; echo
  echo
  echo "[the 3 non-ASCII BYTES my test file carries, by code point] (same count as the 227 recipe file = 3)"
  grep -oP '[^\x00-\x7F]' "$f" | sort | uniq -c | while read -r cnt ch; do
    printf '    x%s  hex=%s  name=U+%s\n' "$cnt" "$(printf '%s' "$ch" | od -An -tx1 | tr -d ' \n')" \
      "$(printf '%s' "$ch" | iconv -f utf-8 -to=UTF-32LE 2>/dev/null | od -An -tx4 | tr -d ' \n' | tr 'a-f' 'A-F')"
  done
  echo "  ⛔ of these are emoji; d22scan (ban #8) rc=0 on the tree with my file in it = 14-nonascii-precedent.txt"
  echo
  echo "[shapes the dispatch forbids inside my test file: want 0 each]"
  printf '  t.Skip            = '; grep -c 't\.Skip' "$f" || true
  printf '  TODO/FIXME/XXX    = '; grep -cE 'TODO|FIXME|XXX' "$f" || true
  printf '  new exported name = '; grep -cE '^func [A-Z]|^var [A-Z]|^type [A-Z]|^const [A-Z]' "$f" || true
  echo "    (the 5 lines matching '^func Test' are test functions, which the go test runner requires;"
  echo "     'new exported name' here counts package-level EXPORTS besides those - see next line)"
  printf '  non-Test exported decls = '; grep -E '^func [A-Z]|^var [A-Z]|^type [A-Z]|^const [A-Z]' "$f" | grep -vc '^func Test' || true
  printf '  panel. method names = '; grep -c 'panel\.' "$f" || true
  printf '  new JSON keys       = '; grep -cE '"[a-z_]+":' "$f" || true
  echo
  echo "[the cited line numbers in my file, for the report's table]"
  grep -n 'the INFO branch must state' "$f"
  grep -c 'kept_in_file_not_in_memory' "$f"
  grep -n 'strings.Contains(logText, "kept_in_file_not_in_memory")' "$f" | sed 's/^/    reverse-assertion line: /'
  grep -n 'the resurrected section carried' "$f"
  grep -n 'the write branch did not re-stamp' "$f"
  grep -n 'came back as something other than the built-in preset' "$f"
  echo
  echo "[report shape] lines / bytes / CR / zero-byte files in my write face"
  for g in .scratch/wisp/probes/341/r1/report.md .scratch/wisp/probes/341/r1/report-answers.md; do
    echo "    $g lines=$(wc -l < "$g") bytes=$(wc -c < "$g")"
  done
  echo "    zero-byte files = $(find .scratch/wisp/probes/341/r1 -type f -size 0 | wc -l)"
  echo "    total files     = $(find .scratch/wisp/probes/341/r1 -type f | wc -l)"
  echo
  echo "[d22scan one more time with the report files present]"
  sh scripts/d22scan.sh 2>&1 | tail -1
  echo "  rc=${PIPESTATUS[0]}"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
