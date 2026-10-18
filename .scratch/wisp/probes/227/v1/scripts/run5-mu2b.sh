#!/bin/bash
# 227-v1 MU-2b: revert ALL THREE log literals to the pre-227 wording (production bytes),
# keep the leg's tests = "未改产码 ⇄ 只换它的用例". The leg reported exactly 2 red
# (文案 assertions at handedit_227_test.go:64 and :129).
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR='D:/tmp/wisp/20261018/227v1/gotmp'
unset TMPDIR
rm -rf "$root/mu2b"; cp -r "$root/deliv" "$root/mu2b"
f="$root/mu2b/internal/config/writeguard.go"
{
  echo "PRISTINE sha256 writeguard.go = $(sha256sum "$root/deliv/internal/config/writeguard.go" | cut -d' ' -f1)"
  echo "BEFORE(landed) sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")"
} > "$out/05-mu2b.txt"
# sentence 3+2: the two delete-shape sentences (content anchors, not line numbers)
perl -0pi -e 's/; a key line the\s*"\s*\+\s*\n\s*"operator deleted by hand comes back at its schema default value//gs' "$f"
perl -0pi -e 's/; a key line the operator deleted by\s*"\s*\+\s*\n\s*"hand comes back at its schema default value//gs' "$f"
# sentence 1: the narrowed nothing-written promise back to the old short sentence
perl -0pi -e 's/already carries this key; this write added nothing \(the\s*"\s*\+\s*\n\s*"file.s bytes and mtime stayed put only because its schema_version was already current;\s*"\s*\+\s*\n\s*"a file behind on schema_version had already been re-laid in canonical form by the\s*"\s*\+\s*\n\s*"re-read, so comments and hand-chosen key order are not preserved\)/already carries this key; nothing written/s' "$f"
echo "AFTER(landed) sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")" >> "$out/05-mu2b.txt"
echo "residue delete-shape token hits = $(grep -c 'comes back at its schema default value' "$f" || true)" >> "$out/05-mu2b.txt"
echo "residue schema-cond token hits = $(grep -c 'only because its schema_version was already current' "$f" || true)" >> "$out/05-mu2b.txt"
echo "old short sentence present = $(grep -c 'nothing written"' "$f" || true)" >> "$out/05-mu2b.txt"
cd "$root/mu2b" || exit 2
go test ./internal/config/ -count=1 -timeout 30m -v > "$out/05-mu2b-full.txt" 2>&1
rc=$?
RUN=$(grep -c '^=== RUN' "$out/05-mu2b-full.txt" || true)
PASSTOP=$(grep -c '^--- PASS' "$out/05-mu2b-full.txt" || true)
FAILTOP=$(grep -c '^--- FAIL' "$out/05-mu2b-full.txt" || true)
SKIP=$(grep -c '^--- SKIP' "$out/05-mu2b-full.txt" || true)
echo "MU2b rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP" >> "$out/05-mu2b.txt"
echo "---- named FAIL roster ----" >> "$out/05-mu2b.txt"
grep '^--- FAIL' "$out/05-mu2b-full.txt" >> "$out/05-mu2b.txt" || echo "NONE" >> "$out/05-mu2b.txt"
echo "---- red file:line verbatim ----" >> "$out/05-mu2b.txt"
grep -nE '^    (writeguard_handedit_227_test|writeguard_226_test|manager_test)\.go:[0-9]+:' "$out/05-mu2b-full.txt" | head -20 >> "$out/05-mu2b.txt"
cp "$root/deliv/internal/config/writeguard.go" "$f"
echo "RESTORED sha256 = $(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")" >> "$out/05-mu2b.txt"
echo "rc=$rc"
