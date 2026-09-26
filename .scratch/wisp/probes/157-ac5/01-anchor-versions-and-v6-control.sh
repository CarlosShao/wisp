#!/usr/bin/env bash
# 157 AC#5 recheck — probe 01: anchor, version timeline, positive control on the V6 line-window ruler.
# read-only; no code, no test runners.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
G=docs/evidence/s1/153-trace-lies-unguarded-r1.md
ISS=.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md

echo "### P0. entry anchor / clock / branch state"
date '+  $ date => %Y-%m-%d %H:%M:%S %z'
echo "  \$ git rev-parse --short HEAD => $(git rev-parse --short HEAD)"
echo "  \$ git rev-parse HEAD         => $(git rev-parse HEAD)"
echo "  \$ git rev-parse --abbrev-ref HEAD => $(git rev-parse --abbrev-ref HEAD)"
echo "  \$ git status -sb | head -1   => $(git status -sb | head -1)"

echo "### P1. commit timeline of the five anchors named in the dispatch (commit author time)"
for a in ec815820 ef1c474f 89a35e9e 67c0177c 5d286d5 4f6c14c; do
  printf '  \$ git log -1 --pretty=%%ad --date=format:%%H:%%M:%%S %s => ' "$a"
  git log -1 --pretty=%ad --date=format:'%Y-%m-%d %H:%M:%S' "$a"
done
echo "  derived gaps (seconds):"
echo "  \$ python-free math via date -d:"
e=$(git log -1 --pretty=%at ec815820); l=$(git log -1 --pretty=%at ef1c474f); a=$(git log -1 --pretty=%at 89a35e9e)
echo "    ec81582->ef1c474 = $((l-e))s ; ef1c474->89a35e9 = $((a-l))s ; ec81582->89a35e9 = $((a-e))s"

echo "### P2. is ef1c474 an ancestor of the acceptor's own close-out commit 89a35e9 ?"
if git merge-base --is-ancestor ef1c474f 89a35e9e; then echo "  YES - B.12 was already on disk inside the acceptor's own close-out tree"; else echo "  NO"; fi
echo "  \$ git cat-file blob 89a35e9:$F | grep -cE '^### B\.12 ' => $(git cat-file blob 89a35e9e:$F | grep -cE '^### B\.12 ')"
echo "  \$ git cat-file blob 89a35e9:$F | wc -l                  => $(git cat-file blob 89a35e9e:$F | wc -l)"

echo "### P3. POSITIVE CONTROL - the V6 line-window ruler verbatim from probes/157-accept/run.sh:137"
echo "  ruler: git cat-file blob <A>:$F | sed -n '389,760p' | grep -c '对账\|双向'"
for a in ec815820 ef1c474f 89a35e9e 67c0177c 5d286d5 4f6c14c; do
  lines=$(git cat-file blob "$a:$F" | wc -l)
  win=$(git cat-file blob "$a:$F" | sed -n '389,760p' | grep -c '对账\|双向')
  full=$(git cat-file blob "$a:$F" | grep -c '对账\|双向')
  b12=$(git cat-file blob "$a:$F" | grep -cE '^### B\.12 ')
  b15=$(git cat-file blob "$a:$F" | grep -cE '^### B\.15 ')
  printf '  %-9s lines=%-4s V6window=%-3s fulltext=%-3s ^### B.12=%s ^### B.15=%s\n' "$a" "$lines" "$win" "$full" "$b12" "$b15"
done

echo "### P4. trailing-newline artifact (why a command-substitution round trip reads one line less)"
git cat-file blob 4f6c14c:$F | tail -c 8 | od -c | head -2 | sed 's/^/  od tail: /'
echo "  \$ git cat-file blob 4f6c14c:$F | wc -l                 => $(git cat-file blob 4f6c14c:$F | wc -l)"
cs=$(git cat-file blob 4f6c14c:$F); printf '%s\n' "$cs" | wc -l | sed 's/^/  \$ via \$()-roundtrip => /'

echo "### P5. line-window boundaries actually contain which sections (the rot mechanism)"
echo "  \$ git cat-file blob 4f6c14c:\$F | grep -n '^### B\.' | cut -c1-28"
git cat-file blob 4f6c14c:$F | grep -n '^### B\.' | cut -c1-28 | sed 's/^/    /'
echo "  lines 389 and 760 of the same version:"
git cat-file blob 4f6c14c:$F | sed -n '389p;760p' | cut -c1-70 | sed 's/^/    /'

echo "### P6. the 153 side at the same versions (the acceptor pinned 593,700p there)"
for a in ec815820 4f6c14c; do
  printf '  %-9s 153件 lines=%-4s V6-153window=%-3s fulltext=%-3s ^### B.9/10=%s\n' "$a" \
    "$(git cat-file blob "$a:$G" | wc -l)" \
    "$(git cat-file blob "$a:$G" | sed -n '593,700p' | grep -c '对账\|双向')" \
    "$(git cat-file blob "$a:$G" | grep -c '对账\|双向')" \
    "$(git cat-file blob "$a:$G" | grep -cE '^### B\.(9|10) ')"
done

echo "### P7. ticket face: AC box counts (anchored) and the '4 枚框' wording"
echo "  \$ grep -c '^- \[ \] \*\*AC#' <157票面> => $(grep -c '^- \[ \] \*\*AC#' "$ISS")"
echo "  \$ grep -c '^- \[x\]'                 => $(grep -c '^- \[x\]' "$ISS")"
echo "  \$ grep -c '4 枚框'                   => $(grep -c '4 枚框' "$ISS")"
echo "  \$ grep -c '5 枚框'                   => $(grep -c '5 枚框' "$ISS")"
echo "  verbatim AC#5 line:"
grep -n '^- \[ \] \*\*AC#5' "$ISS" | sed 's/^/    /'
echo "  \$ git log -1 --pretty='%h %ad' -- <157票面> => $(git log -1 --pretty='%h %ad' --date=format:'%Y-%m-%d %H:%M:%S' -- "$ISS")"
echo "  \$ git status --porcelain -- <157票面> (empty = untouched on disk too) => [$(git status --porcelain -- "$ISS")]"

echo "### P8. the '40 枚' negative conclusion in the acceptor's 3.2 row 5 - does it still read 0?"
for a in ec815820 89a35e9e 4f6c14c; do
  printf '  %-9s 155件 grep -c "40 枚"=%s | 153件=%s\n' "$a" \
    "$(git cat-file blob "$a:$F" | grep -c '40 枚')" \
    "$(git cat-file blob "$a:$G" | grep -c '40 枚')"
done
echo "### P9. the V7 per-section \$-command-line counts (acceptor read 5/7/3/3/4/6/6/6 at ec81582)"
for a in ec815820 4f6c14c; do
  echo "  version $a:"
  git cat-file blob "$a:$F" | awk '/^### B\./{sec=$0; c[sec]=0; n++; order[n]=sec} sec!=""{if(/^\$ /) c[sec]++} END{for(i=1;i<=n;i++) printf "    %-42s $-line=%d\n", substr(order[i],1,42), c[order[i]]}'
done
