#!/bin/bash
# 341-r1: the carrier table the report needs (case -> file:line -> unmutated reading -> de-gum landing -> red line).
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
f=internal/config/writeguard_unnailed_341_test.go
out=.scratch/wisp/probes/341/r1/logs/15-carrier-table.txt
{
  echo "[my five cases: the exact '^func Test' line numbers in the committed file]"
  grep -n '^func Test' "$f"
  echo
  echo "[file shape] lines=$(wc -l < "$f") bytes=$(wc -c < "$f") CR=$(tr -dc '\r' < "$f" | wc -c) nonASCII=$(tr -dc '\200-\377' < "$f" | wc -c)"
  echo "  (for comparison the recipe file 227 used: writeguard_handedit_227_test.go lines=$(wc -l < internal/config/writeguard_handedit_227_test.go) nonASCII=$(tr -dc '\200-\377' < internal/config/writeguard_handedit_227_test.go | wc -c))"
  echo
  echo "[placeholder ruler: any TODO/FIXME/XXX/⛔-as-claim in my file = 0 wanted]"
  grep -cnE 'TODO|FIXME|XXX|placeholder|待补' "$f" || true
  echo
  echo "[unmutated readings, named, from the delivery tree runs (both temp-root surfaces)]"
  for s in Adflt Bdtemp; do
    for c in TestHandDeletedDefaultValuedLineIsRestoredInvisibly TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed TestHandDeletedProviderBaseUrlComesBackAsThePresetValue TestHandDeletedBaseUrlOfAnUnknownProviderStaysEmptyInTheFile; do
      printf '  %s/%s ' "$s" "$c"
      grep -E "^--- (PASS|FAIL): $c" ".scratch/wisp/probes/341/r1/logs/05-color-$s-deliv.txt" | sed 's/^/= /'
    done
  done
  echo
  echo "[the same five on the BASELINE tree (my file absent) = must be absent, which is the census reading]"
  grep -cE 'TestHandDeletedDefaultValuedLineIsRestoredInvisibly' .scratch/wisp/probes/341/r1/logs/05-color-Adflt-base.txt || true
  echo
  echo "[de-gum red lines, verbatim, per shot (the FIRST red sentence of the targeted case = the判点)]"
  for t in p2 p3 p5 p6; do
    echo "  ---- $t ----"
    grep -E 'writeguard_unnailed_341_test\.go:[0-9]+:' ".scratch/wisp/probes/341/r1/logs/05-mu-$t.txt" | head -3
  done
  echo
  echo "[differential reading for AC#2/必答(ii): under the P3 shot (merge base removed), does the P2 case go red?]"
  grep -cE '^--- FAIL: TestHandDeletedDefaultValuedLineIsRestoredInvisibly' .scratch/wisp/probes/341/r1/logs/05-mu-p3-full.txt || true
  echo "  0 above = P2 stays GREEN when the file-based merge base is removed, while its own shot reddens ONLY it (FAIL_TOP=1 in 05-mu-p2.txt)"
  echo
  echo "[必答(ii) production line numbers, at HEAD (= anchor, blob identical)]"
  grep -n 'func (m \*Manager) mergeWrite' internal/config/writeguard.go
  grep -n 'base = deepCopyConfig(disk)' internal/config/writeguard.go
  grep -n 'diverged := diffKeyPaths' internal/config/writeguard.go
  grep -n 'if len(diverged) == 0 {' internal/config/writeguard.go
  grep -n 'slog.Info("config: wrote merged change' internal/config/writeguard.go
  grep -n 'slog.Warn("config: wrote merged change' internal/config/writeguard.go
  grep -n 'func (m \*Manager) writeAllowedDirs' internal/config/allowdirs.go
  grep -n 'm.mergeWrite' internal/config/allowdirs.go internal/config/permmode.go
  echo
  echo "[blob identity of the committed file vs the tree I measured]"
  echo "  git blob      = $(git rev-parse HEAD:internal/config/writeguard_unnailed_341_test.go)"
  echo "  worktree sha  = $(sha256sum "$f" | cut -d' ' -f1)"
  echo "  bench tree    = $(sha256sum /d/tmp/wisp/20261018/341r1/deliv/internal/config/writeguard_unnailed_341_test.go | cut -d' ' -f1)"
  echo "  mu_p2 tree    = $(sha256sum /d/tmp/wisp/20261018/341r1/mu_p2/internal/config/writeguard_unnailed_341_test.go | cut -d' ' -f1)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
