#!/bin/sh
# 333-v1 AC#0 (delivery name roster, re-walked independently) + AC#8 (no loosening,
# no deleted case, no t.Skip; the 327-r2 nail flip direction; stale prose lines).
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 1
A=/d/W333V1-SuqCO8
P=/d/W333V2-hLzM9f
echo "### AC#0 : the six 333-r1 commits, per-commit name-only"
for c in 3f3dde46f dbaad267d 5f4f7945 9574e794 04f2e847 859026d2; do
  echo "--- $c $(git log -1 --format=%s $c | cut -c1-40)"
  git show --name-only --format= $c
done
echo "### AC#0 : union of paths touched, checked against the allowed surfaces"
git diff --name-only ae710c49 859026d2 -- internal/panel > /tmp/ac0-panel.txt
for c in 3f3dde46f dbaad267d 5f4f7945 9574e794 04f2e847 859026d2; do git show --name-only --format= $c; done | sort -u > /tmp/ac0-all.txt
cat /tmp/ac0-all.txt
echo "-- OFF-SURFACE hunt: any path not in the four allowed shapes --"
grep -v -E '^internal/panel/|^\.scratch/wisp/probes/333/r1/|^docs/evidence/s1/333-r1-|^\.scratch/wisp/issues/333-' /tmp/ac0-all.txt \
  | sed 's/^/OFF-SURFACE: /' ; echo "off-surface count = $(grep -v -E '^internal/panel/|^\.scratch/wisp/probes/333/r1/|^docs/evidence/s1/333-r1-|^\.scratch/wisp/issues/333-' /tmp/ac0-all.txt | wc -l)"
echo "-- display_scrub_216_test.go bytes changed across the six? --"
git diff --stat ae710c49 859026d2 -- internal/panel/display_scrub_216_test.go
echo "### AC#8 : diff-level rulers over internal/panel"
echo "-- added t.Skip / Run(. / Parallel : --"
git diff ae710c49 859026d2 -- internal/panel/ | grep -c '^+.*t\.Skip' || true
echo "-- removed test functions (deleted cases) --"
git diff ae710c49 859026d2 -- internal/panel/ | grep '^-' | grep -c 'func Test' || true
echo "-- added test functions --"
git diff ae710c49 859026d2 -- internal/panel/ | grep '^+' | grep -c 'func Test' || true
echo "-- removed assertion lines vs added assertion lines (whole panel pkg) --"
echo "removed t.Error/t.Fatal lines = $(git diff ae710c49 859026d2 -- internal/panel/ | grep '^-' | grep -cE 't\.(Errorf?|Fatalf?|Logf?)' || true)"
echo "added   t.Error/t.Fatal lines = $(git diff ae710c49 859026d2 -- internal/panel/ | grep '^+' | grep -cE 't\.(Errorf?|Fatalf?|Logf?)' || true)"
echo "### AC#8 : the 327-r2 nail flip, anchor vs delivery, side by side"
echo "--- ANCHOR nail 1 body ---"; awk '/^func TestDisplayWashRosterCardinality/,/^}$/' "$P/internal/panel/display_scrub_327r2_test.go"
echo "--- ANCHOR nail 2 body ---"; awk '/^func TestDisplayWashRosterExcludes/,/^}$/' "$P/internal/panel/display_scrub_327r2_test.go"
echo "### AC#8 : the len(cfRoster) != 15 assertion the leg says is untouched"
grep -rn 'len(cfRoster) != 15' "$A/internal/panel" "$P/internal/panel"
echo "### AC#8 : stale prose the leg deliberately left (delivery tree line numbers)"
sed -n '88,93p' "$A/internal/panel/display_scrub_216.go"
sed -n '190,197p' "$A/internal/panel/display_scrub_216.go"
echo "### AC#8 : the escape-not-delete sentence"
sed -n '70,76p' "$A/internal/panel/display_scrub_216.go"
