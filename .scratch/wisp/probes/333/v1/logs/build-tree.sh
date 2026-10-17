#!/bin/sh
# 333-v1 export-tree builder (read-only w.r.t. the repo; writes only into the given stem)
# usage: sh build-tree.sh <anchor> <stem>
set -e
ANCHOR="$1"
STEM="$2"
REPO="D:/work/workspace/projects plans/Wisp"
cd "$REPO"
mkdir -p "$STEM"
git archive "$ANCHOR" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$STEM"
git archive "$ANCHOR" docs | tar -x -C "$STEM"
cp -r build "$STEM"/
mkdir -p "$STEM"/third_party
cp -r third_party/sherpa-onnx "$STEM"/third_party/
cp -r frontend/dist "$STEM"/frontend/ 2>/dev/null || echo "WARN: frontend/dist copy"
# provenance markers
git rev-parse "$ANCHOR" > "$STEM"/.anchor-full
echo ".go files in tree = $(find "$STEM" -name '*.go' | wc -l)"
echo "internal/panel files = $(ls "$STEM"/internal/panel | wc -l)"
