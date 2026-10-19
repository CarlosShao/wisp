#!/bin/sh
# 独立名册尺 第二把：整腿 diff --name-status（端点由第 1 把尺自己导出的名册推得，⛔ 用派单里写的号）
cd "$(git rev-parse --show-toplevel)" || exit 9
OUT=/d/tmp/wisp/20261019/342v2
FIRST=$(head -1 "$OUT/roster-union.tsv" | cut -f1)
LAST=$(tail -1 "$OUT/roster-union.tsv" | cut -f1)
echo "ROSTER_FIRST=$FIRST"
echo "ROSTER_LAST=$LAST"
echo "ROSTER_LAST_SUBJECT=$(git log -1 --format=%s $LAST | cut -c1-60)"
PARENT=$(git rev-parse "$FIRST^")
echo "PARENT_OF_FIRST=$PARENT"
git log -1 --format='PARENT_SUBJECT=%s' "$PARENT" | cut -c1-90
echo '--- NAME-STATUS over PARENT..LAST ---'
git diff --name-status "$PARENT" "$LAST" > "$OUT/leg-diff-namestatus.txt"
wc -l < "$OUT/leg-diff-namestatus.txt" | sed 's/^/NAMESTATUS_LINES=/'
cat "$OUT/leg-diff-namestatus.txt"
echo '--- per-commit name-only ⊆ check (17) ---'
: > "$OUT/per-commit-nameonly.txt"
while IFS=$(printf '\t') read -r h s; do
  echo "COMMIT $h" >> "$OUT/per-commit-nameonly.txt"
  git show --name-only --format= "$h" >> "$OUT/per-commit-nameonly.txt"
done < "$OUT/roster-union.tsv"
grep -c '^COMMIT ' "$OUT/per-commit-nameonly.txt" | sed 's/^/PER_COMMIT_HEADS=/'
echo '--- unique files touched across the leg ---'
grep -v '^COMMIT ' "$OUT/per-commit-nameonly.txt" | sort -u > "$OUT/leg-unique-files.txt"
wc -l < "$OUT/leg-unique-files.txt" | sed 's/^/UNIQUE_FILES=/'
grep -c '\.go$' "$OUT/leg-unique-files.txt" | sed 's/^/UNIQUE_GO=/'
echo '--- added vs modified .go ---'
awk '$1=="A"{print $2}' "$OUT/leg-diff-namestatus.txt" | grep '\.go$' | sort > "$OUT/go-added.txt"
awk '$1=="M"{print $2}' "$OUT/leg-diff-namestatus.txt" | grep '\.go$' | sort > "$OUT/go-mod.txt"
wc -l < "$OUT/go-added.txt" | sed 's/^/GO_ADDED=/'
wc -l < "$OUT/go-mod.txt" | sed 's/^/GO_MODIFIED=/'
cat "$OUT/go-added.txt"; echo '---MOD---'; cat "$OUT/go-mod.txt"
echo '--- diffstat of the modified .go ---'
git diff --numstat "$PARENT" "$LAST" -- internal/winsec/winsec.go
