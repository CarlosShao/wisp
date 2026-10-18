#!/bin/bash
# 227-v1 probe run: mu3 = pristine delivery copy + my two probe cases (off-repo only).
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
rm -rf "$root/mu3"
cp -r "$root/deliv" "$root/mu3"
cp "$root/probe-src/zz_probe_227v1_test.go" "$root/mu3/internal/config/zz_probe_227v1_test.go"
export GOTMPDIR="$root/gotmp"
unset TMPDIR
cd "$root/mu3" || exit 2
{
  echo "probe = 227-v1 read-only probe tests (off-repo tree, never committed; not part of any AC roster)"
  echo "gofiles=$(find . -name '*.go' | wc -l)"
  echo "writeguard.go sha (want efe12661...) = $(sha256sum internal/config/writeguard.go | cut -c1-16)"
  echo "handedit test sha (want 9b430e3b...) = $(sha256sum internal/config/writeguard_handedit_227_test.go | cut -c1-16)"
} > "$out/11-probe-iv-branch-truth.txt"
go test ./internal/config/ -run 'TestProbe227v1' -count=1 -timeout 30m -v >> "$out/11-probe-iv-branch-truth.txt" 2>&1
echo "probe rc=$?" >> "$out/11-probe-iv-branch-truth.txt"
echo "---- leg named cases alongside (for branch comparison) ----" >> "$out/11-probe-iv-branch-truth.txt"
go test ./internal/config/ -run 'TestHandDeletedScalarLineComesBackAtSchemaDefault|TestDeletedLineRestorationIsNotClaimedAsAWrittenKey|TestHandDeletedMapEntryIsNotRestored' -count=1 -timeout 30m -v >> "$out/11-probe-iv-branch-truth.txt" 2>&1
echo "named rc=$?" >> "$out/11-probe-iv-branch-truth.txt"
echo "done"
