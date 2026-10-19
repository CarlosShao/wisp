#!/bin/sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo '--- bf7f13bc parent (= my base anchor candidate) ---'
git log -1 --format='%H %ad %s' --date=iso bf7f13bc~1 | cut -c1-160
echo '--- HEAD ---'
git log -1 --format='%H %ad' --date=iso HEAD
echo '--- blobs: 7b633fa9 vs HEAD vs leg-last(ed8db1ea) ---'
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
  printf '%s BASE7b633fa9=%s HEAD=%s LEG=%s\n' "$f" \
    "$(git rev-parse 7b633fa9:$f)" "$(git rev-parse HEAD:$f)" "$(git rev-parse ed8db1ea:$f)"
done
echo '--- is 7b633fa9 ancestor of HEAD? ---'
git merge-base --is-ancestor 7b633fa9 HEAD && echo ANCYES || echo ANCNO
echo '--- 7b633fa9 subject ---'
git log -1 --format='%h %ad %s' --date=iso 7b633fa9 | cut -c1-120
echo '--- is bf7f13bc~1 == 7b633fa9? ---'
[ "$(git rev-parse bf7f13bc~1)" = "$(git rev-parse 7b633fa9)" ] && echo SAME || echo DIFFERENT
echo '--- CR count check: disk vs blob for provenance.go ---'
printf 'BLOB_CR=%s\n' "$(git cat-file blob HEAD:internal/risk/provenance.go | tr -cd '\r' | wc -c)"
printf 'BLOB_LF=%s\n' "$(git cat-file blob HEAD:internal/risk/provenance.go | tr -cd '\n' | wc -c)"
printf 'DISK_CR=%s\n' "$(tr -cd '\r' < internal/risk/provenance.go | wc -c)"
printf 'DISK_LF=%s\n' "$(tr -cd '\n' < internal/risk/provenance.go | wc -c)"
echo RC=$?
