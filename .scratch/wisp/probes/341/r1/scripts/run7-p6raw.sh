#!/bin/bash
# 341-r1: raw slice of the two P6 cases: show failure headers + the llm/providers part of the dumped file.
set -u
out='/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/341/r1/logs'
log="$out/05-color-Adflt-deliv.txt"
{
  echo "[A] every failure line naming my file, with its line number"
  grep -nE 'writeguard_unnailed_341_test\.go:[0-9]+:' "$log" | cut -c1-200
  echo
  echo "[B] the dumped file bodies: only llm / provider / base_url / protocol lines around those failures"
  grep -nE "^\s*(\[llm|\[llm\.providers|base_url|protocol|api_key_ref|text_chain)" "$log" | sed -n '1,60p' | cut -c1-160
  echo
  echo "[C] the exact 30 lines following the first failure message (line numbers from grep -n above)"
} > "$out/07-p6-raw.txt" 2>&1
first=$(grep -nE 'writeguard_unnailed_341_test\.go:[0-9]+: the deleted line of an unpresetted' "$log" | head -1 | cut -d: -f1)
{
  echo "[C] context around line $first (body of the unknown-provider failure):"
  if [ -n "$first" ]; then sed -n "${first},$((first+40))p" "$log" | grep -E "llm|base_url|protocol|^\s*\[|test.go:[0-9]+" | cut -c1-140; fi
  echo
  echo "[D] context around the FIRST failure of the preset case:"
  s=$(grep -nE 'writeguard_unnailed_341_test\.go:[0-9]+: the deleted base_url line came back' "$log" | head -1 | cut -d: -f1)
  echo "start=$s"
  if [ -n "$s" ]; then sed -n "${s},$((s+40))p" "$log" | grep -E "llm|base_url|protocol|^\s*\[|test.go:[0-9]+" | cut -c1-140; fi
} >> "$out/07-p6-raw.txt" 2>&1
echo "written $(wc -c < "$out/07-p6-raw.txt") bytes"
