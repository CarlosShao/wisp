#!/bin/bash
# 341-r1: verify each de-gum anchor hits EXACTLY the expected count before any shot lands.
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out=.scratch/wisp/probes/341/r1/logs/08-anchor-prescan.txt
{
  echo "[INFO branch literal bytes, writeguard.go lines 171-177]"
  sed -n '171,177p' internal/config/writeguard.go
  echo
  echo "[WARN branch literal bytes, lines 184-190]"
  sed -n '184,190p' internal/config/writeguard.go
  echo
  echo "[anchor hit counts]"
  printf 'P2 anchor  not-preserved; a key line the "+  = '; grep -cF 'are not preserved; a key line the "' internal/config/writeguard.go || true
  printf 'P2 clause  operator deleted by hand comes back at its schema default value  = '; grep -cF 'operator deleted by hand comes back at its schema default value' internal/config/writeguard.go || true
  printf 'P3 anchor  base = deepCopyConfig(disk)  = '; grep -cF 'base = deepCopyConfig(disk)' internal/config/writeguard.go || true
  printf 'P5 anchor  cp.SchemaVersion = SchemaVersionCurrent  = '; grep -cF 'cp.SchemaVersion = SchemaVersionCurrent' internal/config/loader.go || true
  printf 'P6 anchor  p.BaseURL = pre.BaseURL  = '; grep -cF 'p.BaseURL = pre.BaseURL' internal/config/loader.go || true
  echo
  echo "[the two branches' tail strings verbatim, to pick a UNIQUE INFO anchor]"
  grep -n 'not preserved' internal/config/writeguard.go
} > "$out" 2>&1
echo "written $(wc -c < "$out") bytes"
