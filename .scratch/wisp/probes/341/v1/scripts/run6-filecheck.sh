#!/usr/bin/env bash
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
echo "=== files ==="; ls -l .scratch/wisp/probes/341/v1/logs .scratch/wisp/probes/341/v1/scripts
echo "=== zero-byte count ==="; find .scratch/wisp/probes/341/v1 -type f -size 0 | wc -l
echo "=== total files ==="; find .scratch/wisp/probes/341/v1 -type f | wc -l
echo "=== banned char U+00A7 hits in my write面 (ruler self-exclusion noted) ==="
grep -rlP '\x{00A7}' .scratch/wisp/probes/341/v1 2>/dev/null | wc -l
echo "positive control for that ruler (same ruler on docs/reports/HANDOVER.md):"
grep -rlP '\x{00A7}' docs/reports/HANDOVER.md 2>/dev/null | wc -l
echo "=== .out extension count (want 0) ==="; find .scratch/wisp/probes/341/v1 -name "*.out" | wc -l
