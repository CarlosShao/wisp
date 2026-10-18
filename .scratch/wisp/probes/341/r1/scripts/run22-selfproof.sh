#!/bin/bash
# 341-r1 delivery self-proof, three rulers, each run fresh for the file it describes (公共块 交件自证).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
T=.scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md
P=.scratch/wisp/probes/341/r1
out="$P/logs/20-selfproof-final.txt"
# ⛔ no commit inside this gauge: the delivery commit is made by its own script with literal pathspecs.
{
  echo "HEAD at delivery = $(git rev-parse --short HEAD) ; start anchor = 98dbf20a"
  echo
  echo "=== ruler 1: per-commit id + git show --name-only (must contain ONLY my write face) ==="
  for c in $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H'); do
    echo "--- $(git rev-parse --short "$c") $(git log -1 --format='%s' "$c")"
    git show --name-only --format= "$c" | sed '/^$/d' | sed 's/^/      /'
    bad=$(git show --name-only --format= "$c" | sed '/^$/d' \
      | grep -vcE '^(internal/config/[^/]*_test\.go|\.scratch/wisp/probes/341/r1/|\.scratch/wisp/issues/341-)' || true)
    echo "      OUTSIDE-SIDE count = $bad  (want 0)"
    git diff --numstat "$c^" "$c" | grep -E '\.go$' | grep -vE '_test\.go$' | sed 's/^/      PROD-CHANGED /'
  done
  echo "commit count = $(git log --perl-regexp --grep='^票 341 落地腿 341-r1·' --format='%H' | wc -l)"
  echo "  (枚数帽按公共块口径＝参考值、⛔ 判据；本腿按'每落一枚就 commit'落到这个数)"
  echo
  echo "=== ruler 2: absolute paths + wc -c + lines + placeholder ruler (0 hits wanted) ==="
  echo "  main test file: $(pwd)/internal/config/writeguard_unnailed_341_test.go"
  echo "    bytes=$(wc -c < internal/config/writeguard_unnailed_341_test.go) lines=$(wc -l < internal/config/writeguard_unnailed_341_test.go) CR=$(tr -dc '\r' < internal/config/writeguard_unnailed_341_test.go | wc -c) funcTest=$(grep -c '^func Test' internal/config/writeguard_unnailed_341_test.go)"
  echo "  my five files under $P:"
  for g in "$P/report.md" "$P/report-answers.md"; do
    echo "    $(pwd)/$g bytes=$(wc -c < "$g") lines=$(wc -l < "$g")"
  done
  echo "  file census: logs=$(ls -1 "$P/logs" | wc -l) scripts=$(ls -1 "$P/scripts" | wc -l) total=$(find "$P" -type f | wc -l) zero-byte=$(find "$P" -type f -size 0 | wc -l)"
  echo "  placeholder ruler (TODO|FIXME|XXX|待补) hits across my face = $(grep -rlE 'TODO|FIXME|XXX' "$P" 2>/dev/null | wc -l) files"
  echo "  banned character U+00A7 across my face = $(grep -rl -- "$(printf '\302\247')" "$P" internal/config/writeguard_unnailed_341_test.go 2>/dev/null | wc -l) files"
  echo "  .out extension misuse (dispatch forbids it) = $(find "$P" -name '*.out' | wc -l) files"
  echo
  echo "=== ruler 3: the decisive gauges, verbatim, with file:line on both sides ==="
  echo "[unmutated, both temp-root surfaces, named]"
  for s in Adflt Bdtemp; do
    grep -E '^--- (PASS|FAIL): Test(HandDeletedDefault|HandDeletedWhole|HandDeletedSchema|HandDeletedProvider|HandDeletedBaseUrlOfUnknown)' \
      "$P/logs/05-color-$s-deliv.txt" | sed "s|^|    $s |"
  done
  echo "[baseline vs delivery four numbers]"
  for s in Adflt Bdtemp; do for t in base deliv; do
    printf '    %s %s: %s\n' "$s" "$t" "$(grep '^rc=' "$P/logs/05-color-$s-$t-head.txt")"
  done; done
  echo "[four de-gum shots: red sentence + restore hash]"
  for t in p2 p3 p5 p6; do
    printf '    %s landing: %s\n' "$t" "$(grep 'after shot' "$P/logs/05-mu-$t.txt")"
    grep -E 'writeguard_unnailed_341_test\.go:[0-9]+:' "$P/logs/05-mu-$t.txt" | head -1 | sed 's/^/      red /'
    grep 'RESTORED' "$P/logs/05-mu-$t.txt" | sed 's/^/      /'
  done
  echo "[gates, shared worktree, after unset GOTMPDIR TMPDIR]"
  grep -E '^\[gate|^rc=' "$P/logs/11-gates-final.txt" | head -14 | sed 's/^/    /'
  echo
  echo "[ticket discipline]"
  echo "    lines=$(wc -l < "$T") (43 at start) ; my bullet occurrences=$(grep -c '落地腿 341-r1 交件' "$T" || true)"
  echo "    boxes unchecked=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$T" || true)"
  echo "    status line untouched=$(git diff 98dbf20a..HEAD -- "$T" | grep -c '^- ' || true) deleted lines (want 0)"
  echo "    forbidden prod blobs: writeguard=$(git rev-parse HEAD:internal/config/writeguard.go | cut -c1-12) loader=$(git rev-parse HEAD:internal/config/loader.go | cut -c1-12) migrate=$(git rev-parse HEAD:internal/config/migrate.go | cut -c1-12)"
  echo "      anchor values were 0ebc415b0565 / 9100ba92e220 / 78d85e22e951  => same=?"
  echo "    push status: $(git rev-list --count origin/dev..HEAD 2>/dev/null || echo '?') commits ahead of origin/dev (⛔ pushed by this leg)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
