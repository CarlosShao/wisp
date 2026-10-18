#!/bin/bash
# 227-v1 own re-runs of the leg's production-side mutation pair, DIFFERENT landing spots:
# MUA: diverged emptied (adoption-always branch)  -> leg claimed FAIL_TOP=4
# MUB: merge base switched from disk-copy to memory-snapshot -> leg claimed FAIL_TOP=9
# (MUB is the "补回行为" mutation that should redden AC#1 case's ① assertion => AC#1 teeth.)
set -u
root='D:/tmp/wisp/20261018/227v1'
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs'
export GOTMPDIR='D:/tmp/wisp/20261018/227v1/gotmp'
unset TMPDIR
WG="$root/deliv/internal/config/writeguard.go"
PRISTINE=$(sha256sum "$WG" | cut -d' ' -f1)

for tag in mua mub; do
  tree="$root/$tag"
  rm -rf "$tree"; cp -r "$root/deliv" "$tree"
  f="$tree/internal/config/writeguard.go"
  {
    echo "PRISTINE writeguard.go sha256=$PRISTINE (verified before landing: $(sha256sum "$f" | cut -d' ' -f1))"
  } > "$out/09-$tag.txt"
  if [ "$tag" = mua ]; then
    perl -0pi -e 's/diverged := diffKeyPaths\(m\.cur, base\)/diverged := []string{}/s' "$f"
    echo "mutation: diverged := diffKeyPaths(m.cur, base)  ==>  diverged := []string{}" >> "$out/09-$tag.txt"
  else
    perl -0pi -e 's/\tbase = deepCopyConfig\(disk\)/\tbase = deepCopyConfig(m.cur)/s' "$f"
    echo "mutation: base = deepCopyConfig(disk)  ==>  base = deepCopyConfig(m.cur)  (merge base back to memory snapshot)" >> "$out/09-$tag.txt"
  fi
  echo "landed sha256=$(sha256sum "$f" | cut -d' ' -f1) bytes=$(wc -c < "$f")" >> "$out/09-$tag.txt"
  cd "$tree" || exit 2
  go test ./internal/config/ -count=1 -timeout 30m -v > "$out/09-$tag-full.txt" 2>&1
  rc=$?
  RUN=$(grep -c '^=== RUN' "$out/09-$tag-full.txt" || true)
  PASSTOP=$(grep -c '^--- PASS' "$out/09-$tag-full.txt" || true)
  FAILTOP=$(grep -c '^--- FAIL' "$out/09-$tag-full.txt" || true)
  SKIP=$(grep -c '^--- SKIP' "$out/09-$tag-full.txt" || true)
  echo "$tag rc=$rc RUN=$RUN PASS_TOP=$PASSTOP FAIL_TOP=$FAILTOP SKIP=$SKIP" >> "$out/09-$tag.txt"
  echo "---- named FAIL roster ----" >> "$out/09-$tag.txt"
  grep '^--- FAIL' "$out/09-$tag-full.txt" >> "$out/09-$tag.txt" || echo NONE >> "$out/09-$tag.txt"
  # restore the pristine bytes into this throw tree and verify hash equality
  cp "$WG" "$f"
  aft=$(sha256sum "$f" | cut -d' ' -f1)
  echo "RESTORED sha256=$aft  match_pristine=$([ "$aft" = "$PRISTINE" ] && echo YES || echo NO)" >> "$out/09-$tag.txt"
  echo "---- decisive red lines for the AC#1 case (MUB wants ① at :43 red) ----" >> "$out/09-$tag.txt"
  grep -nE '^    writeguard_handedit_227_test\.go:[0-9]+:' "$out/09-$tag-full.txt" | head -8 >> "$out/09-$tag.txt"
  echo "$tag done rc=$rc"
done
