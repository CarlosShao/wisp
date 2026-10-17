#!/bin/sh
# 333-v1 delivery self-check (three rulers), readings taken live, never hand-copied.
R="D:/work/workspace/projects plans/Wisp"
cd "$R" || exit 1
echo "### ruler 1 : my own commits, per-commit name-only"
for c in 4bdacd4a e3333c13 ea64b5e0; do
  echo "--- $c"; git show --name-only --format='%H %s' -s $c | cut -c1-60; git show --name-only --format= $c
done
echo "### off-surface hunt over my three commits"
for c in 4bdacd4a e3333c13 ea64b5e0; do git show --name-only --format= $c; done | sort -u | tee /tmp/v1-paths.txt
echo "off-surface = $(for c in 4bdacd4a e3333c13 ea64b5e0; do git show --name-only --format= $c; done | sort -u | grep -v -E '^docs/evidence/s1/333-v1-verdict.md$|^\.scratch/wisp/probes/333/v1/|^\.scratch/wisp/issues/333-' | wc -l)"
echo "### ruler 2 : the verdict artifact"
V=docs/evidence/s1/333-v1-verdict.md
echo "abs = $(pwd)/$V"
wc -c "$V"; wc -l "$V"
echo "placeholder ruler (word-boundary triple + the two CJK words):"
grep -c -E '\b(TODO|FIXME|XXX)\b|待补|占位' "$V" 2>/tmp/ph-err.txt || true
echo "stderr bytes of that grep = $(wc -c < /tmp/ph-err.txt)"; cat /tmp/ph-err.txt
echo "positive control, same ruler over a file that does contain the words:"
grep -rc -E '\b(TODO|FIXME|XXX)\b|待补|占位' docs/evidence/s1/333-r1-display-cf-category.md
echo "### ruler 3 : product bytes untouched by this leg, shared worktree"
echo "git status --porcelain -- internal/panel  =>"; git status --porcelain -- internal/panel; echo "(empty above = clean)"
echo "git status --porcelain -- internal cmd tools scripts docs/PLAN.md =>"; git status --porcelain -- internal cmd tools scripts docs/PLAN.md
echo "### ticket state: boxes this leg did NOT flip"
T=.scratch/wisp/issues/333-display-boundary-wash-cf-category-not-hand-roster.md
echo "unchecked = $(grep -c '^[[:space:]]*- \[ \]' $T || true)  checked = $(grep -c '^[[:space:]]*- \[x\]' $T || true)"
git diff --numstat ae710c49 HEAD -- "$T" | tail -1
echo "### export trees still on disk (temp artifacts: created, never deleted)"
ls -d /d/W333V1-SuqCO8 /d/W333V2-hLzM9f /d/W333V3-b3pHjd /d/wisp-333v1-tmp 2>&1
echo "### final byte identity of the mutation surface vs the delivery tree"
M=$(cat /d/333v1-tree-m.txt)
diff -r /d/W333V1-SuqCO8/internal/panel "$M/internal/panel" && echo "M == A (no mutation left behind)"
sha256sum /d/W333V1-SuqCO8/internal/panel/display_scrub_216.go "$M/internal/panel/display_scrub_216.go" | cut -c1-16
