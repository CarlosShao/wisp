#!/bin/sh
# 把本腿仓外读数件搬进仓内证据面（只 cp⛔ 删＝临时件只建不删；仓外原件留在 D:\tmp\...\349r1\）。
# 逐枚 wc -c 证非 0、⛔ .out 扩展名。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
T=/d/tmp/wisp/20261019/349r1
R=.scratch/wisp/probes/349/r1
mkdir -p "$R/logs/ac1" "$R/logs/pair" "$R/logs/gates" "$R/logs/census" "$R/logs/teeth" "$R/logs/ac3"
cp -f "$T/logs/01-ac1-blacklist-fires.log" "$R/logs/ac1/"
cp -f "$T/logs/02-ac1-syncdirs-complete.log" "$R/logs/ac1/"
cp -f "$T/logs/03-ac1-provenance-two.log" "$R/logs/ac1/"
cp -f "$T/logs/04-ac1-provenance-disposal.log" "$R/logs/ac1/"
cp -f "$T/logs/00-make-trees.log" "$R/logs/" 2>/dev/null || true
cp -f "$T/logs/gates-outer.log" "$R/logs/gates/gates-outer.log" 2>/dev/null || true
cp -f "$T/logs/teeth/runner-outer.log" "$R/logs/teeth/runner-outer.log" 2>/dev/null || true
cp -f "$T/logs/ac3-outer.log" "$R/logs/ac3/ac3-outer.log" 2>/dev/null || true
cp -f "$T/logs/ac3-outer2.log" "$R/logs/ac3/ac3-outer2.log" 2>/dev/null || true
cp -f "$T/logs/ac3-wordlist-outer.log" "$R/logs/ac3/ac3-wordlist-outer.log" 2>/dev/null || true
cp -f "$T/logs/pair/base.log" "$R/logs/pair/win-base-raw.log" 2>/dev/null || true
cp -f "$T/logs/pair/del.log" "$R/logs/pair/win-del-raw.log" 2>/dev/null || true
cp -f "$T/logs/pair/red-base.txt" "$R/logs/pair/red-base.txt" 2>/dev/null || true
cp -f "$T/logs/pair/red-del.txt" "$R/logs/pair/red-del.txt" 2>/dev/null || true
echo "== non-zero proof, per file =="
find "$R" -type f | sort | while read -r f; do printf '%8s  %s\n' "$(wc -c < "$f")" "$f"; done
echo "== zero-byte (want none) =="
find "$R" -type f -size 0 | sort
echo "== .out extension (want none) =="
find "$R" -type f -name '*.out' | sort
