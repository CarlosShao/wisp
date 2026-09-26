#!/usr/bin/env bash
# 157 AC#5 recheck — probe 08: delivery shape and write-face self-proof (the 交件形状 section).
# X-K1  my commit roster (subject prefix) + the name-only paths of each of my commits + out-of-boundary check
# X-K2  forbidden faces: zero occurrences in MY name-only roster (positive control on the same ruler first)
# X-K3  append-only / untouched: cmp the three read faces against my entry anchor 4f6c14c
# X-K4  ticket faces and -done counts
# X-K5  unpushed proof + shared-index condition re-check
# X-K6  shape of my own deliverable: section list, line count, forbidden-token scan of my own file
# read-only, no code, no test runners
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md
T157=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md
OUT=docs/evidence/s1/157-ac5-recheck-r1.md
ENTRY=4f6c14c

echo "### X-K1  my commit roster and the paths each one touched"
mine=$(git log --pretty=%h --grep='evidence(157 AC#5')
printf '   count by my own subject prefix = %s\n' "$(printf '%s\n' "$mine" | grep -c .)"
for c in $mine; do
  echo "  commit $c $(git log -1 --pretty='%ad %s' --date=format:'%H:%M:%S' "$c" | cut -c1-60)"
  git show --name-only --pretty=tformat: "$c" | sed 's/^/     /'
done
echo "   out-of-boundary check (roster minus my two families):"
for c in $mine; do git show --name-only --pretty=tformat: "$c"; done | sort -u \
  | grep -vE '^docs/evidence/s1/157-ac5-recheck-r1\.md$|^\.scratch/wisp/probes/157-ac5/' | sed 's/^/     FOREIGN: /'
printf '   foreign paths in my roster = %s (grep rc=1 means zero)\n' "$(for c in $mine; do git show --name-only --pretty=tformat: "$c"; done | sort -u | grep -cvE '^docs/evidence/s1/157-ac5-recheck-r1\.md$|^\.scratch/wisp/probes/157-ac5/')"

echo "### X-K2  forbidden faces, each measured with the same roster ruler (and a positive control that it is alive)"
printf '   POSITIVE CONTROL: the same name-only ruler on b23c7f7 counts ^internal/ = %s paths\n' \
 "$(git show --name-only --pretty=tformat: b23c7f7 | grep -cE '^internal/')"
printf '   (a different object, do not cross-copy: git ls-tree -r --name-only b23c7f7 -- internal/ = %s paths)\n' \
 "$(git ls-tree -r --name-only b23c7f7 -- internal/ | wc -l)"
ROSTER=/tmp/157ac5-roster.txt   # outside the repo, create-only (never deleted)
: > "$ROSTER"
for c in $mine; do git show --name-only --pretty=tformat: "$c"; done | sort -u > "$ROSTER"
for pat in '^internal/' '^cmd/' '^docs/PLAN\.md$' '^docs/specs/' '^docs/reports/pending-and-issues\.md$' '^docs/reports/HANDOVER\.md$' '^docs/reports/injection-timeline\.md$' '^frontend/' '^design/' '^tools/d22scan/' '^scripts/' 'thresholds\.go' 'golden' 'allowlist\.txt' 'slo-check' '^\.scratch/wisp/issues/' '^docs/evidence/s1/152-' '^docs/evidence/s1/153-' '^docs/evidence/s1/154-' '^docs/evidence/s1/155-' '^docs/evidence/s1/156-' '^docs/evidence/s1/157-record-level'; do
  printf '   forbidden %-50s hits in my roster = %s\n' "$pat" "$(grep -cE "$pat" "$ROSTER")"
done

echo "### X-K3  append-only / zero-touch proof against my entry anchor"
for p in "$F" "$G" "$ACC" "$T157" docs/reports/pending-and-issues.md docs/reports/HANDOVER.md docs/reports/injection-timeline.md; do
  if git cat-file blob "$ENTRY:$p" | cmp -s - "$p"; then s=IDENTICAL; else s=CHANGED; fi
  printf '   %-58s vs %s => %s\n' "$(basename "$p")" "$ENTRY" "$s"
done
echo "   and the appendix B headings of the two delivered subjects still count the same:"
printf '   155件 ^### B\\. = %s | 153件 ^### B\\. = %s | ^## 附录 B in each = %s/%s\n' \
 "$(git cat-file blob "$ENTRY:$F" | grep -cE '^### B\.')" "$(git cat-file blob "$ENTRY:$G" | grep -cE '^### B\.')" \
 "$(git cat-file blob "$ENTRY:$F" | grep -c '^## 附录 B')" "$(git cat-file blob "$ENTRY:$G" | grep -c '^## 附录 B')"

echo "### X-K4  ticket faces: ticks, boxes, -done"
for t in .scratch/wisp/issues/155-*.md .scratch/wisp/issues/153-*.md "$T157"; do
  printf '   %-62s unchecked=%s checked=%s last-change=%s\n' "$(basename "$t" | cut -c1-58)" \
   "$(grep -c '^- \[ \]' "$t")" "$(grep -c '^- \[x\]' "$t")" "$(git log -1 --pretty='%h %ad' --date=format:'%H:%M:%S' -- "$t")"
done
printf '   -done files matching 157: %s\n' "$(ls .scratch/wisp/issues | grep -c '157.*-done')"

echo "### X-K5  unpushed proof and the shared-index condition"
printf '   ahead of upstream = %s | my commits inside @{u}..HEAD = %s\n' \
 "$(git rev-list --count @{u}..HEAD)" "$(git log --oneline @{u}..HEAD --grep='evidence(157 AC#5' | grep -c .)"
echo "   \$ git diff --cached --name-only  (must list nothing but my own, or nothing at all):"
git diff --cached --name-only | sed 's/^/     staged: /'
echo "   \$ git status --porcelain -- docs/reports/"
git status --porcelain -- docs/reports/ | sed 's/^/     /'

echo "### X-K6  my own deliverable shape"
printf '   lines=%s  sections:\n' "$(wc -l < "$OUT")"
grep -nE '^#{1,3} ' "$OUT" | cut -c1-90 | sed 's/^/     /'
echo "   banned emoji bands in MY OWN file (the repo's instrument band, strings not comments):"
for band in 'U+1F000-1FAFF' 'U+2200-22FF' 'U+2600-27BF' 'U+2B00-2BFF' 'U+FE0F'; do
  case "$band" in
    'U+1F000-1FAFF') rgx='[\x{1F000}-\x{1FAFF}]';; 'U+2200-22FF') rgx='[\x{2200}-\x{22FF}]';;
    'U+2600-27BF') rgx='[\x{2600}-\x{27BF}]';; 'U+2B00-2BFF') rgx='[\x{2B00}-\x{2BFF}]';;
    'U+FE0F') rgx='\x{FE0F}';;
  esac
  printf '     %-14s hits=%s\n' "$band" "$(LC_ALL=C.UTF-8 grep -cP "$rgx" "$OUT" 2>/dev/null || echo "n/a")"
done
echo "   my probe roster:"
ls -1 .scratch/wisp/probes/157-ac5 | sed 's/^/     /'
