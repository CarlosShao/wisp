#!/bin/bash
# 341-r1 start gauges (dispatch 起手 1/2/3). Read-only: no writes outside probes/341/r1.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/341/r1/logs'
mkdir -p "$out" || exit 2
A=98dbf20a35027a9290f52649d8ab2634646cb14e   # my own start anchor (git log -1, first action)
t=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
{
  echo "[gauge 1] anchors, own readings (⛔ the dispatch's numbers)"
  echo "start HEAD (first action) = $A"
  echo "git log -1 now            = $(git log -1 --format='%H %s' | cut -c1-90)"
  echo "git status --porcelain -- internal/config cmd/wisp (first action) = <empty, zero lines>"
  echo "git status --porcelain -- internal/config cmd/wisp (now)          ="
  git status --porcelain -- internal/config cmd/wisp
  echo "[end status]"
  echo
  echo "  HEAD drift 98dbf20a..HEAD = $(git log --format='%H' 98dbf20a..HEAD | wc -l) commits, subjects:"
  git log --format='    %h %s' 98dbf20a..HEAD | cut -c1-120
  echo "  do the drifted commits touch my write face? (want: no lines)"
  git diff --name-only 98dbf20a..HEAD -- internal/config cmd/wisp
  echo "[end drift-paths]"
  echo
  echo "  blob equality check for the three forbidden-to-touch prod files + the shared fixtures,"
  echo "  my start anchor vs current HEAD (ruler: git rev-parse <ref>:<path> = blob id)"
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go \
           internal/config/manager_test.go internal/config/writeguard_226_test.go \
           internal/config/writeguard_handedit_227_test.go internal/config/defaults.go internal/config/parse.go; do
    echo "    $p  $A=$(git rev-parse $A:$p | cut -c1-12)  HEAD=$(git rev-parse HEAD:$p | cut -c1-12)  same=$([ "$(git rev-parse $A:$p)" = "$(git rev-parse HEAD:$p)" ] && echo YES || echo NO)"
  done
  echo
  echo "[gauge 2] ticket 现量 four rows, re-run against the NAMED commit tree ($A), ⛔ in-flight worktree"
  echo "  row1: ticket 227 pinned = P1/P4/P5' three cases in writeguard_handedit_227_test.go"
  echo "    carrier lines: :23 / :70 / :106  (尺 = git show $A:<file> | grep -n '^func Test')"
  git show "$A:internal/config/writeguard_handedit_227_test.go" | grep -n '^func Test'
  echo "    lines=$(git show "$A:internal/config/writeguard_handedit_227_test.go" | wc -l)  CR count=$(git show "$A:internal/config/writeguard_handedit_227_test.go" | grep -c $'\r' || true)"
  echo
  echo "  row2: P2/P3/P5/P6 四格 today unnailed. ⛔ copy the table - see gauge 3 (own census)."
  echo
  echo "  row3: the overclaim sentence ruler (dispatch AC#3): git grep -n 're-laid in canonical form' $A -- internal/config"
  git grep -n 're-laid in canonical form' "$A" -- internal/config | cut -c1-150
  echo "    (row3 line numbers quoted WITH this anchor = $A short)"
  echo
  echo "  row4: the A1033/票面 de-gum recipe = manager_test.go:38 const matrixBase; check that line NOW"
  git show "$A:internal/config/manager_test.go" | sed -n '37,39p'
  echo "    schema_version line inside matrixBase? $(git show "$A:internal/config/manager_test.go" | sed -n '37,56p' | grep -c 'schema_version = 2' || true) (want 1)"
  echo
  echo "[gauge 3] AC#0 SELF-ADDED ruler: per shape P1..P6, does ANY *_test.go assert THIS shape's RESULT?"
  echo "  scope = git grep -nE over the tree at $A, restricted to *_test.go under internal/config + cmd/wisp"
  echo "  分母 = 226-v1 b2 P1-P6 (⛔ 重造). 尺 = single-line content anchors, per shape, plus reading the carrier."
  echo
  echo "  -- P0 改值 (size 56->60)"
  git grep -nE 'Replace\(matrixBase, "size = 56", "size = 60' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P1 删一枚值非默认的标量行"
  git grep -nE '"size = 60[^\\]*\\n", ""' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P2 删一枚值本就是默认的标量行 (anchor: deleting the literal 'size = 56' line)"
  git grep -nE '"size = 56\\n", ""' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P3 删整节 (any test that deletes a whole [ball] header+body)"
  git grep -nE '\[ball\]' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P4 删 map 条目 (provider sub-table)"
  git grep -nE 'llm\.providers\.deepseek' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P5/P5' 删 schema_version 行"
  git grep -nE '"schema_version = 2\\n", ""' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "  -- P6 删 provider 的 base_url 行 (any test that deletes a base_url line)"
  git grep -nE 'base_url[^"]*\\n", ""' "$A" -- '*_test.go' | cut -c1-160
  echo "  P6 wider: every test mentioning base_url with a Replace/delete on it"
  git grep -nE 'Replace.*base_url' "$A" -- '*_test.go' | cut -c1-160
  echo
  echo "[gauge 4] ticket box ruler (criteria-region scoped): want 7 unchecked / 0 checked"
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true)  checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)"
  grep -n '^[[:space:]]*- \[' "$t" | cut -c1-70
  echo "  ticket lines=$(wc -l < "$t")"
  echo
  echo "[gauge 5] leg-identity triple key: my subject prefix on disk = ? commits (want 0)"
  echo "  ruler: git log --all --grep='票 341 落地腿 341-r1' --format='%h %ad'"
  git log --all --grep='票 341 落地腿 341-r1' --format='%h %ad' --date=format:'%m-%d %H:%M'
  echo "  count=$(git log --all --grep='票 341 落地腿 341-r1' --format='%H' | wc -l)"
  echo "  find internal/config -name '*.go' -newermt '-30 minutes' = [see gauge 6]"
} > "$out/00-start.txt" 2>&1
echo "written $(wc -c < "$out/00-start.txt") bytes"
