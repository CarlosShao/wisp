#!/usr/bin/env bash
# 341-v1: stage two export trees (delivery + base) outside the repo.
R="/d/work/workspace/projects plans/Wisp"
cd "$R" || exit 9
W=/d/tmp/wisp/20261018/341v1

mk() {
  anchor="$1"; dst="$2"
  mkdir -p "$dst"
  ( cd "$R" && git archive "$anchor" go.mod go.sum internal cmd tools scripts frontend models ) | tar -x -C "$dst"
  ( cd "$R" && git archive "$anchor" docs ) | tar -x -C "$dst"
  mkdir -p "$dst/build" "$dst/third_party/sherpa-onnx" "$dst/frontend"
  cp -f "$R/build/"*.dll "$dst/build/"
  cp -f "$R/third_party/sherpa-onnx/"*.dll "$dst/third_party/sherpa-onnx/"
  cp -rf "$R/frontend/dist" "$dst/frontend/"
  echo "=== TREE $dst anchor=$anchor ==="
  ( cd "$dst" && ls -d go.mod go.sum internal cmd tools scripts frontend models docs build third_party/sherpa-onnx frontend/dist 2>/dev/null | tr '\n' ' '; echo
    echo "GO_FILES_TOTAL=$(find . -name '*.go' | wc -l)"
    echo "CONFIG_GO=$(ls internal/config/*.go 2>/dev/null | wc -l)"
    echo "CONFIG_TEST_GO=$(ls internal/config/*_test.go 2>/dev/null | wc -l)"
    echo "EMPTY_GO=$(find . -name '*.go' -size 0 | wc -l)" )
  ( cd "$dst" && go list ./internal/config/ >/dev/null 2>&1; echo "GOLIST_RC=$?" )
}

mk "$ANCHOR_D" "$W/treeD"
mk "$ANCHOR_B" "$W/treeB"
echo "=== presence of 341 test file: treeD=$(ls $W/treeD/internal/config/writeguard_unnailed_341_test.go 2>/dev/null | wc -l) treeB=$(ls $W/treeB/internal/config/writeguard_unnailed_341_test.go 2>/dev/null | wc -l) ==="
echo "=== pristine hashes (sha256 first 12) ==="
for t in treeD treeB; do
  for f in internal/config/writeguard.go internal/config/loader.go internal/config/migrate.go; do
    printf "%s %s %s\n" "$t" "$f" "$(sha256sum "$W/$t/$f" | cut -c1-12)"
  done
done
