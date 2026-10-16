#!/usr/bin/env bash
# 编排者复跑 317-v1 的成对读数（R2a/R2b 同形）：同一目录、只差分隔符
set -u
root=/c/tmp/wisp317orch
log=$root/logs
mkdir -p "$log" "$root/tree"
nat='C:\tmp\wisp317orch\pair'
psx='C:/tmp/wisp317orch/pair'
mkdir -p /c/tmp/wisp317orch/pair
cd "/d/work/workspace/projects plans/Wisp" || exit 9
git archive HEAD | tar -x -C "$root/tree"
n=$(find "$root/tree" -name '*.go' | wc -l)
echo "exported_tree_go_files=$n" | tee "$log/00-tree.txt"
for p in go.mod go.sum internal cmd tools scripts; do
  [ -e "$root/tree/$p" ] && echo "sixpiece_ok=$p" >> "$log/00-tree.txt" || echo "sixpiece_MISSING=$p" >> "$log/00-tree.txt"
done
cd "$root/tree" || exit 8
RE='TestSyncAncestorActableLegFailsClosedOnMovedAncestor116|TestSyncFirstActableLegStillRefusesHomeRewriteAtVerdict116|TestSyncAncestorLegReResolvesThroughJunctionToRealTree|TestSyncAncestorLegNotExemptedStaysFailClosed'
for form in nat psx; do
  if [ "$form" = nat ]; then export GOTMPDIR="$nat"; else export GOTMPDIR="$psx"; fi
  echo "form=$form GOTMPDIR=$(printf '%s' "$GOTMPDIR" | od -c | head -3 | tr '\n' ' ')" > "$log/10-$form.txt"
  GOFLAGS= go test ./internal/risk/ -count=1 -v -run "$RE" >> "$log/10-$form.txt" 2>&1
  echo "rc=$?" >> "$log/10-$form.txt"
  {
    printf 'form=%s RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s rc=%s\n' "$form" \
      "$(grep -c '^=== RUN' "$log/10-$form.txt" || true)" \
      "$(grep -c '^--- PASS' "$log/10-$form.txt" || true)" \
      "$(grep -c '^--- FAIL' "$log/10-$form.txt" || true)" \
      "$(grep -c '^--- SKIP' "$log/10-$form.txt" || true)" \
      "$(grep -c '^rc=' "$log/10-$form.txt" || true)"
  } >> "$log/01-summary.txt"
done
grep '^--- FAIL' "$log/10-nat.txt" "$log/10-psx.txt" > "$log/02-failnames.txt" 2>&1
echo "done" >> "$log/01-summary.txt"
