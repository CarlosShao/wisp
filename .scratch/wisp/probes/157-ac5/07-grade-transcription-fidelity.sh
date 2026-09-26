#!/usr/bin/env bash
# 157 AC#5 recheck — probe 07: transcription fidelity (does B.15 15.1's 档 column match its cited source?).
# This is the semantic-level ruler the dispatch asked for, applied to the GRADE column:
#   15.1 claims "档＝157 验收件 §3.1 逐字转记" and cites `89a35e9` 版 `:166`–`:170`.
#   Here the cited lines are taken at that anchor and compared row by row with what 15.1 wrote.
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 9
F=docs/evidence/s1/155-three-unjudged-cells-r1.md
ACC=docs/evidence/s1/157-record-level-cleanup-r1-accept-r1.md
ANCH=4f6c14c

echo "### W-A  where the AC rows of 157-accept 3.1 actually are (worktree copy vs the cited anchor 89a35e9)"
echo "  worktree copy:"
awk -F'|' '/^\| \*\*AC#[1-5]\*\*/{printf "   line %-4s %s grade=%s\n", NR, substr($2,1,11), substr($(NF-1),1,40)}' "$ACC"
echo "  same ruler on the anchor the transcription cites (89a35e9 = the acceptor's own commit):"
git cat-file blob 89a35e9e:"$ACC" | awk -F'|' '/^\| \*\*AC#[1-5]\*\*/{printf "   line %-4s %s grade=%s\n", NR, substr($2,1,11), substr($(NF-1),1,40)}'
printf '  blob equal? worktree=%s anchor=%s\n' "$(git hash-object "$ACC")" "$(git rev-parse 89a35e9e:"$ACC")"

echo "### W-B  15.1 grade column (as transcribed by r2) side by side with the source grade"
git cat-file blob "$ANCH:$F" | awk -F'|' '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#/{gsub(/\\[|]/,"Q"); printf "   TRANSCRIBED %s -> %s\n", substr($2,1,11), substr($4,1,40)}'
echo "   (source rows are the W-A lines above; compare one-to-one: AC#1..AC#4 成立, AC#5 不成立（未交）)"

echo "### W-C  does the AC#5 row of 15.1 keep its own verdict distinct from the source?"
git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f&&/^\| \*\*AC#5/' | cut -c1-330

echo "### W-D  the four sub-asks of AC#5, each with its reading (all from probes 01/02/03/05/06, this line restates them)"
printf '  1) 双向都在场        : 12.1 forward=%s rows / 12.2 reverse=%s rows / 15.1=%s / 15.2=%s / 15.3=%s\n' \
 "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -c '^| \*\*AC#')" \
 "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -cE '^\| [0-9]+ \|')" \
 "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.1/{f=1;next} f&&/^\*\*15\.2/{f=0} f' | grep -c '^| \*\*AC#')" \
 "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.2/{f=1;next} f&&/^\*\*15\.3/{f=0} f' | grep -cE '^\| [0-9] ')" \
 "$(git cat-file blob "$ANCH:$F" | awk '/^\*\*15\.3/{f=1;next} f&&/^\*\*15\.4/{f=0} f' | grep -cE '^\| [0-9]+ \|')"
printf '  2) 明写"未裁"        : 12.1 grade cells containing 未裁 = %s / 5 ; B.15 15.3 rows each naming who closes = 7 / 7 (probe 06 V-D)\n' \
 "$(git cat-file blob "$ANCH:$F" | awk -F'|' 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f&&/^\| \*\*AC#/{gsub(/\\[|]/,"Q"); if($(NF-1) ~ /未裁/) c++} END{print c+0}')"
printf '  3) 最小闭合          : 12.2 closure cells 127/137/74/56/103 bytes (probe 06 V-B) ; 15.3 76/50/101/55/35/355/47 bytes (V-C) ; empty-shell rows (cell<10 bytes) = 0 both tables\n'
printf '  4) 4 枚框 vs 5 枚框  : ticket face unchecked boxes=%s ; ticket AC#5 text says "4 枚框"=%s ; ticket 155 boxes=%s ; delivered forward rows=%s\n' \
 "$(grep -c '^- \[ \] \*\*AC#' .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md)" \
 "$(grep -c '4 枚框' .scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md)" \
 "$(grep -cE '^- \[[ x]\] \*\*AC#' .scratch/wisp/issues/155-ticket-153-leaves-three-cells-unjudged-because-my-attack-point-grid-and-the-ticket-ac-grid-are-not-the-same-coordinate-system-done.md)" \
 "$(git cat-file blob "$ANCH:$F" | awk 'index($0,"### B.12 ")==1{f=1;next} f&&/^### /{f=0} f' | grep -c '^| \*\*AC#')"
