#!/bin/sh
# 342-v3 gates (shared worktree ONLY) + frame ruler + blob identities. Run after `unset` per the dispatch.
cd "/d/work/workspace/projects plans/Wisp" || exit 1
unset GOTMPDIR TMPDIR TEMP TMP
OUT=/d/tmp/wisp/20261020/342v3/logs
{
echo "== gate 0: gofumpt on a .md = tool misuse, rc recorded, NOT read as a red =="
/d/work/base/gopath/bin/gofumpt.exe -l docs/evidence/s1/348-v1-verdict.md >/dev/null 2>&1
echo "gofumpt-on-md rc=$?"
echo "== gate 1: gofumpt -l over my write face (0 .go files today) =="
/d/work/base/gopath/bin/gofumpt.exe -l >/dev/null 2>&1
echo "(skipped: no-arg gofumpt walks the tree; my write face has no .go bytes today)"
echo "== gate 2: d22scan =="
sh scripts/d22scan.sh > "$OUT/gate-d22scan.txt" 2>&1
echo "d22scan rc=$?"
tail -4 "$OUT/gate-d22scan.txt"
echo "== gate 3: path length budget =="
sh scripts/check-path-length-budget.sh > "$OUT/gate-length.txt" 2>&1
echo "length rc=$?"
tail -4 "$OUT/gate-length.txt"
echo "== gate 4: go vet ./internal/winsec/ (local = windows) =="
go vet ./internal/winsec/ > "$OUT/gate-vet-local.txt" 2>&1
echo "vet-local rc=$?"
cat "$OUT/gate-vet-local.txt"
echo "== cross vet (registered as compiler coverage only, NOT paired red/green with local) =="
GOOS=linux go vet ./internal/winsec/ > "$OUT/gate-vet-linux.txt" 2>&1
echo "vet-linux rc=$?"
GOOS=darwin go vet ./internal/winsec/ > "$OUT/gate-vet-darwin.txt" 2>&1
echo "vet-darwin rc=$?"
echo "== box rulers (the AC frames as they stand NOW, before my legs) =="
for n in 342 348; do
  f=$(ls .scratch/wisp/issues/${n}-*.md)
  x=$(grep -c '^- \[x\] \*\*AC' "$f" || true)
  o=$(grep -c '^- \[ \] \*\*AC' "$f" || true)
  echo "$f  [x]=$x  [ ]=$o"
done
echo "== open frames by name =="
for n in 342 348; do
  f=$(ls .scratch/wisp/issues/${n}-*.md)
  echo "--- $f"
  grep '^- \[ \] \*\*AC' "$f" | cut -c1-28
done
echo "== winsec.go identity =="
git rev-parse HEAD:internal/winsec/winsec.go
echo "-- the :118-140 block sha16 (my ruler = sha1sum over sed -n '118,140p'; 342-v2 reported sha16 e8a2f33dea8cdeab over the same range: compare) --"
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '118,140p' | sha1sum | cut -c1-16
echo "-- what 342-v2's own number covers: its ruler was sed -n 118,140p too. Print first/last line of the range --"
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '118p;140p'
echo "== tail of the 342-v2 quote: is the ancestor-Windows sentence at :138 or :140? =="
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '134,141p' | cat -n
echo "== production callers unchanged =="
grep -c 'os.Remove(resolved)' internal/winsec/winsec.go || true
grep -rn 'unlinkPinnedLeaf(parent' internal/winsec/winsec.go
echo "== the seam 342-r2 leaned on exists at baseline too (so case one is not landing-only plumbing) =="
git cat-file blob fde1df02:internal/winsec/winsec_windows.go | grep -n 'var applyDescriptor'
git cat-file blob fde1df02:internal/winsec/winsec_windows.go | grep -n 'func sealHandle'
git cat-file blob HEAD:internal/winsec/winsec_windows.go | grep -n 'func sealHandle'
echo "== the two instrument files: line counts (ruler = wc -l over HEAD blobs) =="
for p in internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go; do
  echo "$p lines=$(git cat-file blob HEAD:$p | wc -l)"
done
echo "== POSIX-family liveness names inside the two Windows instruments (expect 0 = no borrowed green) =="
for p in internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/leafrace_342r2_windows_test.go; do
  c=$(git cat-file blob HEAD:$p | grep -c 'pinnedParentBeforeOpen\|provedLinkCycles' || true)
  echo "$p borrowed-POSIX-key-hits=$c"
done
echo "== where those two names actually live =="
grep -rln 'pinnedParentBeforeOpen\|provedLinkCycles' internal/winsec/ | head -10
echo "== the landing's SameFile face: hits inside the re-check roster across ALL my logs =="
grep -c 'SameFile face' "$OUT"/treeL-det-r*.log || true
} > "$OUT/gates-and-frames.txt" 2>&1
echo "GATES_RC=$?"
