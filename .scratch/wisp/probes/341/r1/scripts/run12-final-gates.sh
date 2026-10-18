#!/bin/bash
# 341-r1 delivery side: gates (shared worktree) + per-commit numstat (必答 v) + red-roster diff.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
unset GOTMPDIR TMPDIR
out=.scratch/wisp/probes/341/r1/logs/11-gates-final.txt
mystack='ticket 341|341-r1|writeguard_unnailed_341'
{
  echo "HEAD at measurement = $(git rev-parse --short HEAD)"
  echo
  echo "[gate 1] gofumpt -l on MY file (empty + rc=0 = clean)"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/config/writeguard_unnailed_341_test.go
  echo "rc=$?"
  echo
  echo "[gate 1b] gofumpt -l on the whole package: the pre-existing CRLF inventory, named only"
  /d/work/base/gopath/bin/gofumpt.exe -l internal/config
  echo "rc=$?   (⛔ of these are mine: my file is the only one I wrote and it is clean above)"
  echo
  echo "[gate 2] go vet ./internal/config/"
  go vet ./internal/config/ 2>&1
  echo "rc=$?"
  echo
  echo "[gate 3] sh scripts/d22scan.sh  (last line is the verdict; the self-test inside prints seeded violations)"
  sh scripts/d22scan.sh 2>&1 | tail -3
  echo "rc=${PIPESTATUS[0]}"
  echo
  echo "[gate 4] sh scripts/check-path-length-budget.sh"
  sh scripts/check-path-length-budget.sh 2>&1 | tail -4
  echo "rc=${PIPESTATUS[0]}"
  echo
  echo "[必答 v] git diff --numstat per commit of this leg: production .go added/deleted MUST be 0/0"
  for c in $(git log --format='%H' --grep='票 341 落地腿 341-r1' --reverse); do
    echo "  --- commit $(git log -1 --format='%h %s' "$c")"
    git diff --numstat "$c^" "$c" 2>/dev/null | awk '
      /\.go$/ && !/_test\.go$/ { print "      PROD  " $0 }
      /_test\.go$/            { print "      TEST  " $0 }
      !/\.go$/                { print "      OTHER " $0 }'
    prod=$(git diff --numstat "$c^" "$c" 2>/dev/null | grep -E '\.go$' | grep -vE '_test\.go$' | wc -l)
    echo "      production .go paths changed = $prod  (want 0)"
  done
  echo
  echo "[write-face ruler] every path touched by my commits, union:"
  for c in $(git log --format='%H' --grep='票 341 落地腿 341-r1'); do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u
  echo
  echo "[越界尺] paths outside internal/config test face + probes/341/r1 + ticket 341 = "
  for c in $(git log --format='%H' --grep='票 341 落地腿 341-r1'); do git show --name-only --format= "$c"; done | sed '/^$/d' | sort -u \
    | grep -vcE '^(internal/config/.*_test\.go$|\.scratch/wisp/probes/341/r1/|\.scratch/wisp/issues/341-)' || true
  echo
  echo "[non-comment byte check on the three forbidden files] blob id anchor-vs-HEAD (must be identical all through)"
  A=98dbf20a35027a9290f52649d8ab2634646cb14e
  for p in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
    echo "    $p anchor=$(git rev-parse "$A:$p") HEAD=$(git rev-parse "HEAD:$p") same=$([ "$(git rev-parse "$A:$p")" = "$(git rev-parse "HEAD:$p")" ] && echo YES || echo NO)"
  done
  echo
  echo "[ticket box ruler] must still be 7 unchecked / 0 checked (⛔ 翻框)"
  t=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
  echo "  unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true) lines=$(wc -l < "$t")"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
