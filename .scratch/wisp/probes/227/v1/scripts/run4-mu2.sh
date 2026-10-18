#!/bin/bash
# 227-v1 MU-2: different landing spot from MU-1 (production bytes, not test bytes).
# Remove the two added delete-shape sentences in writeguard.go (equivalent to
# "未改产码 ⇄ 只换它的用例") => the leg's 文案 assertions must go red.
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR='D:/tmp/wisp/20261018/227v1/gotmp'
unset TMPDIR
rm -rf "$root/mu2"; cp -r "$root/deliv" "$root/mu2"
f="$root/mu2/internal/config/writeguard.go"
{
echo "PRISTINE sha256 writeguard.go = $(sha256sum "$root/deliv/internal/config/writeguard.go" | cut -d' ' -f1)"
echo "BEFORE(landed) sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")"
} > "$out/04-mu2.txt"
# content anchor (the sentence token), not a line number
perl -0pi -e 's/; a key line the\s*"\s*\+\s*\n\s*"operator deleted by hand comes back at its schema default value//gs' "$f"
perl -0pi -e 's/; a key line the operator deleted by\s*"\s*\+\s*\n\s*"hand comes back at its schema default value//gs' "$f"
echo "AFTER(landed) sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")" >> "$out/04-mu2.txt"
echo "residue hits of the sentence token in landed file (want 0): $(grep -c 'comes back at its schema default value' "$f" || true)" >> "$out/04-mu2.txt"
cd "$root/mu2" || exit 2
go test ./internal/config/ -count=1 -timeout 30m -v > "$out/04-mu2-full.txt" 2>&1
rc=$?
RUN=$(grep -c '^=== RUN' "$out/04-mu2-full.txt" || true)
PASSTOP=$(grep -c '^--- PASS' "$out/04-mu2-full.txt" || true)
FAILTOP=$(grep -c '^--- FAIL' "$out/04-mu2-full.txt" || true)
SKIP=$(grep -c '^--- SKIP' "$out/04-mu2-full.txt" || true)
echo "MU2 rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP" >> "$out/04-mu2.txt"
echo "---- named FAIL roster ----" >> "$out/04-mu2.txt"
grep '^--- FAIL' "$out/04-mu2-full.txt" >> "$out/04-mu2.txt" || echo "NONE" >> "$out/04-mu2.txt"
echo "---- red sentences verbatim (numbered lines of the full log) ----" >> "$out/04-mu2.txt"
grep -nE '^    (writeguard_handedit_227_test|writeguard_226_test|manager_test)\.go:[0-9]+:' "$out/04-mu2-full.txt" | head -20 >> "$out/04-mu2.txt"
cp "$root/deliv/internal/config/writeguard.go" "$f"
echo "RESTORED sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")" >> "$out/04-mu2.txt"
echo "rc=$rc"
