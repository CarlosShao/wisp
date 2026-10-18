#!/bin/bash
# 341-r1 bench: build two export trees (baseline = my start anchor, delivery = baseline + my new test file).
# Recipe from the dispatch (public block "台面配方"): git archive of the eight named paths + docs/,
# on the D: volume (8.3 short names must exist), plus the three untracked pieces when present.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
A=98dbf20a35027a9290f52649d8ab2634646cb14e
root='D:/tmp/wisp/20261018/341r1'
out=.scratch/wisp/probes/341/r1/logs/03-trees.txt
mkdir -p "$root/gotmp" || exit 2
{
  echo "[0] 8.3 short name check on the D: volume (dispatch: 目录名要产生 8.3 短名)"
  cmd //c "dir /x D:\\tmp\\wisp\\20261018" | sed -n '1,20p'
  echo
  echo "[1] what exists in the repo for the three manual pieces (untracked):"
  for p in build third_party/sherpa-onnx frontend/dist; do
    if [ -d "$p" ]; then echo "    $p  DIR present  dllfiles=$(find "$p" -maxdepth 1 -name '*.dll' | wc -l)  bytes=$(du -sb "$p" | cut -f1)"; else echo "    $p  ABSENT"; fi
  done
  echo
  echo "[2] archive: git archive $A go.mod go.sum internal cmd tools scripts frontend models  (+ docs)"
  # ⛔ rm here: the dispatch says temporary pieces are built, never deleted. Fresh numbered dirs instead.
  for t in base deliv; do
    mkdir -p "$root/$t" || exit 2
    git archive "$A" go.mod go.sum internal cmd tools scripts frontend models | tar -x -C "$root/$t" || { echo "archive FAILED for $t"; exit 3; }
    git archive "$A" docs | tar -x -C "$root/$t" || echo "  (docs archive skipped: rc=$?)"
  done
  echo "  base = anchor as-is;  deliv = base + my new test file copied in from the worktree"
  cp internal/config/writeguard_unnailed_341_test.go "$root/deliv/internal/config/" || exit 4
  echo
  echo "[3] prove the trees are non-empty (a 0 count is a bench failure, not a reading):"
  for t in base deliv; do
    echo "    $t: go files = $(find "$root/$t" -name '*.go' | wc -l)   config pkg files = $(ls "$root/$t/internal/config" | wc -l)"
    echo "    $t: my file present = $([ -f "$root/$t/internal/config/writeguard_unnailed_341_test.go" ] && echo YES || echo NO)  sha256=$(sha256sum "$root/$t/internal/config/writeguard_unnailed_341_test.go" 2>/dev/null | cut -d' ' -f1)"
  done
  echo "  worktree sha256 = $(sha256sum internal/config/writeguard_unnailed_341_test.go | cut -d' ' -f1)"
  echo
  echo "[4] pristine hashes of the three production files (de-gum targets are verified against these):"
  for p in writeguard.go loader.go migrate.go; do
    echo "    internal/config/$p  base=$(sha256sum "$root/base/internal/config/$p" | cut -d' ' -f1)"
    echo "    internal/config/$p deliv=$(sha256sum "$root/deliv/internal/config/$p" | cut -d' ' -f1)"
  done
  echo
  echo "[5] anchor of each tree's writeguard_handedit_227_test.go (must be identical to the anchor blob)"
  git rev-parse "$A:internal/config/writeguard_handedit_227_test.go" | cut -c1-12
  echo "  base copy sha=$(sha256sum "$root/base/internal/config/writeguard_handedit_227_test.go" | cut -d' ' -f1)"
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes rc=$?"
