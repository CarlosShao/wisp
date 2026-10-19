#!/usr/bin/env bash
# 342-v1 container driver, phase 2: deterministic rows (A1 by-name through a planted
# ancestor link, A3 the refusal-path unlink) + landed production rows on mut-jia/mut-yi
# + the full-package colours that answer AC#0's "which existing pins go red".
set -u
S=/d/tmp/wisp/20261018/342v1
P=$S/probe
L=$S/logs
mkdir -p "$L"

for form in jia yi; do
  cp "$P/v342v1-landed-prod.go.txt" "$S/mut-$form/internal/winsec/v342v1landed_test.go"
  echo "== landed probe hash mut-$form: $(sha256sum "$S/mut-$form/internal/winsec/v342v1landed_test.go" | cut -d' ' -f1)"
done

echo "### (1) deterministic A1 + A3, treeV, repeats 1..3"
for r in 1 2 3; do
  bash "$S/scripts/v342v1-docker.sh" treeV \
    "go test -count=1 -timeout 20m -v -run 'TestV342v1A1DeterministicAncestorFollow|TestV342v1A3RefusalPathUnlink' ./internal/winsec/" \
    "11-det-a1a3-r$r.log"
  grep -E 'V342V1 A1|V342V1 A3' "$L/11-det-a1a3-r$r.log" || true
done

echo "### (2) landed production rows, 1500 attempts, repeats 1..3 per form"
for form in jia yi; do
  for r in 1 2 3; do
    bash "$S/scripts/v342v1-docker.sh" "mut-$form" \
      "V342V1_ATTEMPTS=1500 go test -count=1 -timeout 40m -v -run TestV342v1PosixLandedProductionRow ./internal/winsec/" \
      "12-landed-$form-r$r.log"
    grep -E 'V342V1 VPROD-LANDED: attempts' "$L/12-landed-$form-r$r.log" || true
    grep -E 'V342V1 VPROD-LANDED leg=' "$L/12-landed-$form-r$r.log" | head -8 || true
    echo "winsec.go sha in mut-$form: $(sha256sum "$S/mut-$form/internal/winsec/winsec.go" | cut -d' ' -f1)"
  done
done

echo "### (3) full-package colours: baseline + three landing faces (CLEAN trees, no probe files)"
for form in jia yi yinofb; do
  bash "$S/scripts/v342v1-build-col.sh" "$form" > "$L/00-colbuild-$form.txt" 2>&1
  echo "colbuild $form rc=$? :"; tail -5 "$L/00-colbuild-$form.txt" || true
done
bash "$S/scripts/v342v1-docker.sh" treeC "go test -count=1 -timeout 40m -v ./internal/winsec/" "13-color-treeC.log"
for tree in treeC col-jia col-yi col-yinofb; do
  bash "$S/scripts/v342v1-docker.sh" "$tree" \
    "go test -count=1 -timeout 40m -v ./internal/winsec/" "13-color-$tree.log"
  grep -E '^--- FAIL' "$L/13-color-$tree.log" | sed 's/ ([0-9.]*s)//' > "$L/13-color-$tree-failnames.txt" || true
  echo "failnames $tree: $(wc -l < "$L/13-color-$tree-failnames.txt")"
  cat "$L/13-color-$tree-failnames.txt" || true
done
echo "PHASE2-DONE"
