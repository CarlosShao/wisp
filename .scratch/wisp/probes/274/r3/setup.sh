#!/bin/bash
# 274-r3 bench setup: out-of-repo mirror trees (cp -al per top-level dir,
# NEVER .git, NEVER frontend/node_modules, NEVER the repo build/ dir).
set -uo pipefail
ROOT="/d/work/workspace/projects plans/Wisp"
BASE="/d/tmp/wisp274r3"
STAMP=$(date +%s)
mkdir -p "$BASE/bin"

for t in mut red3 fresh; do
  T="$BASE/tree-$t"
  if [ -e "$T" ]; then mv "$T" "$T.old-$STAMP"; fi
  mkdir -p "$T"
  for d in .github agent cmd docker docs internal models plans third_party tools; do
    cp -al "$ROOT/$d" "$T/" || echo "COPYDIR_FAIL $d"
  done
  for f in go.mod go.sum deps.toml .gitattributes .gitignore AGENTS.md README.md; do
    cp -a "$ROOT/$f" "$T/" || echo "COPYFILE_FAIL $f"
  done
  # frontend: everything except node_modules (excluded by discipline) and dist (rebuilt below)
  mkdir -p "$T/frontend"
  for fe in "$ROOT"/frontend/* "$ROOT"/frontend/.*; do
    fb=$(basename "$fe")
    case "$fb" in node_modules|dist|.|..) continue;; esac
    cp -al "$fe" "$T/frontend/"
  done
  mkdir -p "$T/frontend/node_modules"   # forces the 'node_modules EXISTS -> npm ci SKIPPED' branch
  # scripts: REAL copy, never a hardlink into the repo (tree build.ps1 is swapped per case)
  cp -a "$ROOT/scripts" "$T/scripts"
  mkdir -p "$T/build"                   # fresh build dir, repo build/ is never linked
  # dist: REAL copies of the 4 stale tracked bytes (not hardlinks -> runs cannot write the repo)
  D="$T/frontend/dist"
  mkdir -p "$D/assets"
  cp -a "$ROOT/frontend/dist/.gitkeep" "$D/"
  cp -a "$ROOT/frontend/dist/index.html" "$D/"
  cp -a "$ROOT/frontend/dist/assets/." "$D/assets/"
done

# THE MUTATION SHAPE: dist additionally holds one page file whose NAME CONTAINS SPACES.
cp -a "$ROOT/frontend/dist/index.html" "$BASE/tree-mut/frontend/dist/index with a space.html"

# stub npm: normalise the Write-authored .cmd to CRLF (no bash printf, no %1 width specifiers)
awk '{printf "%s\r\n", $0}' "$BASE/bin/npm.cmd" > "$BASE/bin/npm.cmd.cr" && mv "$BASE/bin/npm.cmd.cr" "$BASE/bin/npm.cmd"

echo "SETUP_OK"
for t in mut red3 fresh; do
  echo "--- tree-$t dist ---"
  ls -1 "$BASE/tree-$t/frontend/dist" "$BASE/tree-$t/frontend/dist/assets"
  echo "--- tree-$t scripts/build.ps1 md5 ---"
  md5sum "$BASE/tree-$t/scripts/build.ps1"
done
echo "--- repo dist untouched (must equal start.txt) ---"
(cd "$ROOT" && md5sum frontend/dist/.gitkeep frontend/dist/assets/* frontend/dist/index.html)
