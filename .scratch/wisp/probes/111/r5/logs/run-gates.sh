# 111-r5 local gate readings (no CI, no push). Writes a full log, prints only a summary.
# Rulers: bash -n / dash -n on the carriers my change depends on, the guarded step's OWN
# command line, and the YAML census before/after. NO go commands except the ones the
# measured step itself runs.
set -u
out=.scratch/wisp/probes/111/r5/logs/gates-local.txt
: > "$out"
say() { printf '%s\n' "$*" >> "$out"; }

say "== date =="
date "+%Y-%m-%d %H:%M %z" >> "$out"

say
say "== bash -n on the carrier scripts (rc=0 required) =="
for f in scripts/portable-tests.sh scripts/portable-tests-selftest.sh scripts/winsec-tests.sh \
         scripts/runtests.sh tools/d22scan/runtests.sh scripts/slo-freshness.sh \
         scripts/d22scan.sh scripts/check-path-length-budget.sh \
         .scratch/wisp/probes/161/r5/attrib.sh; do
    if bash -n "$f" 2>>"$out"; then r=0; else r=$?; fi
    say "bash -n rc=$r  $f"
done

say
say "== the guarded step's own command line, verbatim from ci.yml =="
cmd=$(grep -n 'attrib.sh --tracked-only' .github/workflows/ci.yml | head -1)
say "ci.yml line: $cmd"
sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only >> "$out" 2>&1
say "attrib.sh --tracked-only rc=$?"

say
say "== YAML census, same parser before/after (PyYAML) =="
for v in before after; do
    say "--- $v ---"
    head -9 ".scratch/wisp/probes/111/r5/logs/yaml-census-$v.txt" >> "$out"
done

say
say "== step-name census in ci.yml (what each guard-bearing step is) =="
grep -nE '^\s+if: ' .github/workflows/ci.yml | sed 's/^/  /' >> "$out"

wc_l=$(wc -l < "$out")
echo "gates-local.txt lines=$wc_l"
