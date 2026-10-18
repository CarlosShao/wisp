#!/bin/bash
# 227-v1 (必答 i) gauge: SUFFICIENCY of the "+1 line above the two assertions" shape.
# Pristine committed bytes put the two red sentences at :404/:407 (my MU-1 reading).
# The leg reported :405/:408 on a file whose committed blob is byte-identical (29dbd968).
# This run measures whether inserting ONE line inside the function (content-anchored)
# reproduces exactly the leg's numbers while reddening the SAME two assertions.
# It does NOT decide what the leg actually did — that answer is the leg's to give.
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR="$root/gotmp"; unset TMPDIR
rm -rf "$root/mu4"; cp -r "$root/deliv" "$root/mu4"
MT="$root/mu4/internal/config/manager_test.go"
WT="$root/mu4/internal/config/writeguard_226_test.go"
{
  echo "PRISTINE manager_test.go sha (want 201dc786...) = $(sha256sum "$root/deliv/internal/config/manager_test.go" | cut -c1-16)"
  echo "PRISTINE writeguard_226 sha (want 49960da6...) = $(sha256sum "$root/deliv/internal/config/writeguard_226_test.go" | cut -c1-16)"
  echo "pristine red-sentence call-site lines (content anchor, git-tree bytes):"
  grep -n 'rewritten anyway:\|mtime moved on a write that wrote nothing' "$root/deliv/internal/config/writeguard_226_test.go"
} > "$out/13-one-line-gap.txt"
sed -i '/^schema_version = 2$/d' "$MT"
echo "after AC#2 mutation: manager_test bytes=$(wc -c < "$MT") (pristine 9677)" >> "$out/13-one-line-gap.txt"
sed -i '/^func TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt/,+1i\	// probe line for the one-row gap question (227-v1, off-repo only)' "$WT"
echo "writeguard_226 lines now=$(wc -l < "$WT") pristine=$(wc -l < "$root/deliv/internal/config/writeguard_226_test.go")" >> "$out/13-one-line-gap.txt"
echo "probe-line landed at:" >> "$out/13-one-line-gap.txt"
grep -n 'probe line for the one-row gap' "$WT" >> "$out/13-one-line-gap.txt" || echo "INSERT DID NOT LAND" >> "$out/13-one-line-gap.txt"
echo "red-sentence call-site lines AFTER insert:" >> "$out/13-one-line-gap.txt"
grep -n 'rewritten anyway:\|mtime moved on a write that wrote nothing' "$WT" >> "$out/13-one-line-gap.txt"
cd "$root/mu4" || exit 2
go test ./internal/config/ -run 'TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt' -count=1 -timeout 30m -v >> "$out/13-one-line-gap.txt" 2>&1
echo "targeted rc=$?" >> "$out/13-one-line-gap.txt"
cp "$root/deliv/internal/config/manager_test.go" "$MT"
cp "$root/deliv/internal/config/writeguard_226_test.go" "$WT"
{
  echo "RESTORED manager_test.go sha = $(sha256sum "$MT" | cut -c1-16) bytes=$(wc -c < "$MT")"
  echo "RESTORED writeguard_226 sha = $(sha256sum "$WT" | cut -c1-16) bytes=$(wc -c < "$WT")"
} >> "$out/13-one-line-gap.txt"
echo done
