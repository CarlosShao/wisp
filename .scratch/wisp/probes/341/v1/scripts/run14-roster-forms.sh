#!/usr/bin/env bash
# 341-v1: roster ruler, three forms side by side, to separate A1030-style mixing from a real 7th leg commit.
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
echo "[F1] --grep='票 341 落地腿 341-r1·' (the dispatch's literal string, middle dot U+00B7)"
git log --all --grep='票 341 落地腿 341-r1·' --format='   %h %s' | cut -c1-72
echo "   count = $(git log --all --grep='票 341 落地腿 341-r1·' --format='%H' | wc -l)"
echo
echo "[F2] anchored subject regex '票 341 落地腿 341-r1·' (^ start = subject line only)"
git log --all --grep='^票 341 落地腿 341-r1·' --format='   %h %s' | cut -c1-72
echo "   count = $(git log --all --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l)"
echo
echo "[F3] the byte that F1 actually matched (context bytes around 341-r1) for the flagged commit"
git log -1 --format='%s' 95b1015d | grep -oP '票 341.{0,3}终裁腿 341-v1\x{00b7}' | od -c | head -3
echo "   => does 95b1015d's subject contain the substring 落地腿 ? $(git log -1 --format='%s' 95b1015d | grep -c '落地腿' || true)"
echo "   => its BODY (message below the subject) mention of the leg prefix:"
git log -1 --format='%b' 95b1015d | grep -c '票 341 落地腿 341-r1·' || true
echo
echo "[F4] MY roster with the same anchored form (want 8)"
git log --all --grep='^票 341 终裁腿 341-v1·' --format='   %h %ad %s' --date=format:'%H:%M' | cut -c1-70
echo "   count = $(git log --all --grep='^票 341 终裁腿 341-v1·' --format='%H' | wc -l)"
echo
echo "[F5] untracked-index state (the 3 entries I do NOT own)"
git diff --cached --name-only | sed 's/^/   /'
