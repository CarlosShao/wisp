#!/bin/bash
# 341-r1 start gauges, second attempt, ASCII-only echoes.
# Why a new file: the first attempt (logs/00-start.txt) died mid-run -
#   defect 1: $'\r' inside a double-quoted string is NOT ANSI-C quoting, it opened an
#             unbalanced single quote -> "line 36: unexpected EOF while looking for `)'"
#             (so everything after gauge 2 row1 never ran);
#   defect 2: `cut -c1-90` applied to git subjects truncated a multi-byte CJK char,
#             making the log itself invalid UTF-8.
# Both are ruler defects, not readings; kept on disk as evidence (temp files are built, never deleted).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/341/r1/logs'
mkdir -p "$out" || exit 2
A=98dbf20a35027a9290f52649d8ab2634646cb14e
t=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
WG=internal/config/writeguard_handedit_227_test.go
MT=internal/config/manager_test.go
{
  echo "[gauge 1] anchors, own readings (NOT the dispatch numbers)"
  echo "start HEAD (my first action) = $A"
  echo "git log -1 at start = 98dbf20a35027a9290f52649d8ab2634646cb14e <subject: dispatch write档, see 00-start.txt line2 region>"
  echo "git log -1 now      ="
  git log -1 --format='%H %s'
  echo "git status --porcelain -- internal/config cmd/wisp  at start  = <empty, zero lines>"
  echo "git status --porcelain -- internal/config cmd/wisp  now       ="
  git status --porcelain -- internal/config cmd/wisp
  echo "[end status: zero lines above means nobody is writing my face]"
  echo
  echo "  drift 98dbf20a..HEAD = $(git log --format='%H' 98dbf20a..HEAD | wc -l) commits"
  git log --format='    %h %s' 98dbf20a..HEAD
  echo "  paths touched by drift inside my face + cmd/wisp (want: none):"
  git diff --name-only 98dbf20a..HEAD -- internal/config cmd/wisp
  echo "[end drift paths]"
  echo
  echo "  blob equality anchor-vs-HEAD, ruler = git rev-parse <ref>:<path>"
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go \
           internal/config/manager_test.go internal/config/writeguard_226_test.go \
           internal/config/writeguard_handedit_227_test.go internal/config/defaults.go internal/config/parse.go; do
    x=$(git rev-parse "$A:$p"); y=$(git rev-parse "HEAD:$p")
    echo "    $p same=$([ "$x" = "$y" ] && echo YES || echo NO)  ${x:0:12}"
  done
  echo
  echo "[gauge 2] ticket row1 re-run on the named commit tree (not the in-flight worktree)"
  echo "  ruler: git show \$A:$WG | grep -n '^func Test'"
  git show "$A:$WG" | grep -n '^func Test'
  echo "  lines=$(git show "$A:$WG" | wc -l)"
  printf '  CR bytes = %s\n' "$(git show "$A:$WG" | tr -dc '\r' | wc -c)"
  echo "  => carrier :23 (P1) / :70 (P4) / :106 (P5') present; ticket row1 says 195 lines/5 funcs/CR 0"
  echo
  echo "[gauge 2b] row3: the overclaim ruler, with THIS anchor's line numbers"
  git grep -n 're-laid in canonical form' "$A" -- internal/config
  echo
  echo "[gauge 2c] row4: de-gum recipe = manager_test.go:38 const matrixBase"
  git show "$A:$MT" | sed -n '36,42p'
  echo
  echo "[gauge 3] AC#0 self-added ruler: for each shape P1..P6, does ANY *_test.go assert it?"
  echo "  scope: git grep -nE over tree $A, only *_test.go under internal/config + cmd/wisp"
  echo "-- P0 changed value (size 56 -> 60)"
  git grep -nE 'Replace\(matrixBase, "size = 56", "size = 6[04]' "$A" -- '*_test.go'
  echo "-- P1 delete a NON-default scalar line"
  git grep -nE '"size = 60[^"]*\\n", ""' "$A" -- '*_test.go'
  echo "-- P2 delete a line whose value WAS the default"
  git grep -nE '"size = 56\\n", ""' "$A" -- '*_test.go'
  echo "-- P3 delete a whole section (anything removing the [ball] header)"
  git grep -nE '\[ball\]' "$A" -- '*_test.go'
  echo "-- P4 delete a map entry (provider sub-table)"
  git grep -nE '"\[llm\.providers\.deepseek\]' "$A" -- '*_test.go'
  echo "-- P5 / P5' delete schema_version line"
  git grep -nE '"schema_version = 2\\n", ""' "$A" -- '*_test.go'
  echo "-- P6 delete provider base_url line"
  git grep -nE 'base_url[^"]*\\n", ""' "$A" -- '*_test.go'
  echo "-- P6 wider (any Replace touching base_url)"
  git grep -nE 'Replace.*base_url' "$A" -- '*_test.go'
  echo "[end gauge 3: zero-match sections above ARE readings because the positive controls P0/P1/P5 hit]"
  echo
  echo "[gauge 4] ticket box ruler, scoped to criteria region"
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true)"
  echo "  checked  =$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  echo "  ticket lines=$(wc -l < "$t")"
  echo
  echo "[gauge 5] leg identity triple key: my subject prefix already on disk?"
  echo "  ruler: git log --all --grep=<prefix> ; prefix = 341-r1"
  git log --all --grep='341-r1' --format='    %h %ad %s' --date=format:'%m-%d %H:%M'
  echo "  count=$(git log --all --grep='341-r1' --format='%H' | wc -l)  (want 0 before my first commit)"
  echo
  echo "[gauge 6] write-face contention: find internal/config -name '*.go' -newermt '-40 minutes'"
  find internal/config -name '*.go' -newermt '-40 minutes' -print
  echo "[end gauge 6: empty = nobody landing on my face]"
} > "$out/01-start-gauges.txt" 2>&1
echo "written $(wc -c < "$out/01-start-gauges.txt") bytes lines $(wc -l < "$out/01-start-gauges.txt")"
