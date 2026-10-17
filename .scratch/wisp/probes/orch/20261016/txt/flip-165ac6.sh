#!/usr/bin/env bash
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
T165=".scratch/wisp/issues/165-plan-first-then-act-a-brake-the-voice-entry-needs.md"
echo "== before =="
printf 'unchecked='; grep -cE '^[[:space:]]*- \[ \]' "$T165"
printf 'checked=';  grep -cE '^[[:space:]]*- \[x\]'  "$T165"
grep -c '^- \[ \] \*\*AC#6' "$T165" | sed 's/^/hits=/'
sed -i 's/^- \[ \] \*\*AC#6/- [x] **AC#6/' "$T165"
echo "== after =="
printf 'unchecked='; grep -cE '^[[:space:]]*- \[ \]' "$T165"
printf 'checked=';  grep -cE '^[[:space:]]*- \[x\]'  "$T165"
echo "== numstat (must be N/N and only box flips) =="
git diff --numstat -- "$T165"
git diff --unified=0 -- "$T165" | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)' | cut -c1-40
