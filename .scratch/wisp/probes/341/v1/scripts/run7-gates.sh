#!/usr/bin/env bash
# 341-v1 AC#4 gates: shared worktree only, GOTMPDIR/TMPDIR/TEMP/TMP unset first
# (ticket 311 AC#1 shape B GUARD reads an env'd D: temp root as rc=2 scope refusal).
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
unset GOTMPDIR TMPDIR TEMP TMP
echo "=== env after unset: GOTMPDIR=[$GOTMPDIR] TMPDIR=[$TMPDIR] TEMP=[$TEMP] TMP=[$TMP] ==="
echo "=== HEAD at gate time ==="; git log -1 --format='%H %s' | cut -c1-40
echo "=== internal/config worktree status (want clean) ==="; git status --porcelain -- internal/config cmd/wisp; echo "[empty above = clean]"

echo
echo "=== [1] gofumpt -l on the leg's new file (want empty, rc=0) ==="
"/d/work/base/gopath/bin/gofumpt.exe" -l internal/config/writeguard_unnailed_341_test.go
echo "GOFUMPPT_RC=$?"
echo "--- package-wide -l (informational: existing CRLF stock) ---"
"/d/work/base/gopath/bin/gofumpt.exe" -l internal/config | head -10
echo "GOFUMPPT_PKG_RC=$?"

echo
echo "=== [2] go vet ./internal/config/ (want rc=0) ==="
go vet ./internal/config/ 2>&1 | tail -10
echo "VET_RC=${PIPESTATUS[0]}"

echo
echo "=== [3] d22scan (want clean rc=0; scope needs design/ present) ==="
sh scripts/d22scan.sh 2>&1 | tail -6
echo "D22_RC=${PIPESTATUS[0]}"

echo
echo "=== [4] path-length budget (want rc=0 GREEN; needs .git) ==="
sh scripts/check-path-length-budget.sh 2>&1 | tail -8
echo "PATHLEN_RC=${PIPESTATUS[0]}"

echo
echo "=== AC#5 forbidden-zone blob identity, my own re-run ==="
for x in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
  printf "%s HEAD=%s\n" "$x" "$(git rev-parse "HEAD:$x")"
done
echo "--- start-anchor(98dbf20a, the leg's own) vs HEAD for the three ---"
for x in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
  a=$(git rev-parse "98dbf20a:$x" 2>/dev/null); b=$(git rev-parse "HEAD:$x")
  printf "%s leg_anchor=%s head=%s same=%s\n" "$x" "${a:0:12}" "${b:0:12}" "$([ "$a" = "$b" ] && echo YES || echo NO)"
done
echo "--- internal/config non-test .go changes in 6948eb08^..HEAD ---"
git log --name-status --format="COMMIT %h" 6948eb08^..HEAD -- internal/config | grep -E '^(COMMIT|[AMD])'
echo "--- .go files touched by the leg roster (union of 6 commits) ---"
git show --name-only --format='' 6948eb08 fd82f786 f710d104 920537ee 46c920c8 cef2a6c9 | sort -u | grep '\.go$' | sed 's/^/   /'
echo "--- t.Skip / Skip( in the leg's test file (want 0) ---"
grep -c 't\.Skip\|runtime\.GoMAXPROCS' internal/config/writeguard_unnailed_341_test.go || true
grep -c 'Skip(' internal/config/writeguard_unnailed_341_test.go || true
echo "--- AC box gauge (ticket tail) ---"
T=".scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md"
echo "unchecked=$(grep -c '^- \[ \] \*\*AC#' "$T" || true) checked=$(grep -c '^- \[x\]' "$T" || true) lines=$(wc -l < "$T")"
echo "Status line unchanged: $(grep -c '^- Status: \*\*已立，未派\*\*' "$T" || true)"
echo "title: $(head -1 "$T" | cut -c1-40)"
