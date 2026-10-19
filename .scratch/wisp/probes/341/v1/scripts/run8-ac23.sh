#!/usr/bin/env bash
# 341-v1 AC#2/AC#3/AC#5 measurement block.
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
T=/d/tmp/wisp/20261018/341v1/treeD

echo "########## AC#2(i) the quote-shape change: is canonical output really single-quote? ##########"
echo "-- MarshalCanonical string emission (production, my own read) --"
grep -n "MarshalCanonical" internal/config/parse.go | head -5
sed -n '150,175p' internal/config/parse.go
echo "-- the emitter for string values (which quote char) --"
grep -n "WriteByte(''')\|'%s'\|%q\|WriteString(\"'\")" internal/config/parse.go | head -20

echo
echo "-- live evidence from MY OWN unmutated delivery-tree run: the re-laid file's quoted lines --"
grep -nE "^ *(base_url|protocol|model|language|theme|monitor) = " "$T/internal/config/writeguard_unnailed_341_test.go" >/dev/null
echo "-- the P6 target assertion lines (post-fix, verbatim) --"
grep -n "base_url = '" internal/config/writeguard_unnailed_341_test.go
grep -n "protocol = '" internal/config/writeguard_unnailed_341_test.go
echo "-- any double-quote form left in a Contains(text,...) assertion in that file? --"
grep -nE 'Contains\(text, "[^"]*"' internal/config/writeguard_unnailed_341_test.go | sed 's/^/   /'

echo
echo "-- the leg's own PRE-FIX red lines (logs/06, logs/07) vs current assertion sites --"
grep -n "writeguard_unnailed_341_test.go:3" .scratch/wisp/probes/341/r1/logs/07-p6-raw.txt | head
echo "current file line numbers for the same three messages:"
grep -n "came back as something other than the built-in preset\|did not come back as an explicit empty string\|did not survive the re-lay" internal/config/writeguard_unnailed_341_test.go
echo "-- the comment block that occupies the 6-line delta --"
sed -n '312,320p' internal/config/writeguard_unnailed_341_test.go

echo
echo "########## AC#2(ii): t.Skip / assertion-shape scan on the leg's file ##########"
echo "t.Skip hits = $(grep -c 't\.Skip' internal/config/writeguard_unnailed_341_test.go || true)"
echo "top-level func decls (want only Test*):"
grep -nE '^func ' internal/config/writeguard_unnailed_341_test.go
echo "new exported names introduced (want 0 beyond the 5 Test funcs):"
grep -oE '^func [A-Z][A-Za-z0-9_]*' internal/config/writeguard_unnailed_341_test.go | sort
echo "helpers used from 227's file (shared fixtures, not new behaviour):"
grep -nE 'matrixBase|newTestManager|readText|readLoaded|statMod|captureSlog|testDir' internal/config/writeguard_unnailed_341_test.go | head -8
echo "where those helpers are defined:"
grep -rnE 'func (newTestManager|readLoaded|readText|statMod|captureSlog)\(' internal/config/*_test.go | sed 's/^/   /'
grep -rn 'matrixBase = \|matrixBase =' internal/config/*_test.go | head -3

echo
echo "########## AC#3 the 226 verdict file: my own whole-file census ##########"
for f in docs/evidence/s1/226-config-write-no-clobber-r1.md docs/evidence/s1/226-config-write-no-clobber-v2.md; do
  echo "-- $f --"
  printf "   wc -l = %s\n" "$(wc -l < "$f")"
  printf "   不白重排 hits = %s\n" "$(grep -c '不白重排' "$f" || true)"
  grep -n '不白重排' "$f" | sed 's/^/     /' | cut -c1-160
  printf "   不白丢注释 hits = %s\n" "$(grep -c '不白丢注释' "$f" || true)"
  grep -n '不白丢注释' "$f" | sed 's/^/     /' | cut -c1-160
done
echo "-- whole docs/evidence/s1 for these two words (broader denominator) --"
grep -rn '不白重排\|不白丢注释' docs/evidence/s1 | sed 's/^/   /' | cut -c1-120
echo "-- r1.md sec 1 region: what does :22 actually say --"
sed -n '22p' docs/evidence/s1/226-config-write-no-clobber-r1.md | cut -c1-400
echo "-- r1.md :111 verbatim (the overclaim sentence) --"
sed -n '111p' docs/evidence/s1/226-config-write-no-clobber-r1.md
echo "-- ticket 226 issues file (the append target) --"
ls .scratch/wisp/issues | grep '^226'

echo
echo "########## AC#3(ii) production self-admission: comments DO get dropped ##########"
grep -n 're-laid in canonical form' internal/config/writeguard.go
sed -n '152,156p' internal/config/writeguard.go
sed -n '172,176p' internal/config/writeguard.go
sed -n '184,190p' internal/config/writeguard.go
echo "-- live proof from my own run: the leg's P3 shot output contains no operator comment after re-lay --"
grep -c 'my size' "$T/internal/config/writeguard_unnailed_341_test.go"
