#!/bin/bash
# 227-v1 self-proof gauges, delivery side: empty-slot ruler with positive control + inventory.
cd '/d/work/workspace/projects plans/Wisp' || exit 2
out='.scratch/wisp/probes/227/v1/logs/15-final-roster.txt'
{
  echo "== empty-slot ruler on my verdict file (want 0 hits) =="
  v=docs/evidence/s1/227-v1-verdict.md
  pat='TODO|TBD|待填|占位|placeholder|FIXME|XXX'
  h=$(grep -cE "$pat" "$v" || true)
  echo "slot-ruler hits=$h  (want 0; rc of grep itself would be 1 => reported via || true)"
  echo "== positive controls (same shape, tokens known to EXIST) =="
  p1=$(grep -cE 'AC#|verdict' "$v" || true)
  echo "posctrl-1 (ruler shape over a pattern present in the same file) hits=$p1  (want >0)"
  p2=$(grep -cE 'schema_version' internal/config/manager_test.go || true)
  echo "posctrl-2 (ruler shape over another file, known token) hits=$p2  (want >0)"
  p3=$(grep -cE 'ZZZnotathing' "$v" || true)
  echo "negctrl (known-absent token) hits=$p3  (want 0)"
  echo "verdict wc: $(wc -l < "$v") lines / $(wc -c < "$v") bytes (cap 130 lines)"
  echo
  echo "== zero-byte check over my write surface =="
  z=$(find .scratch/wisp/probes/227/v1 -type f -size 0 | wc -l)
  echo "zero-byte files under probes/227/v1 = $z (want 0); files = $(find .scratch/wisp/probes/227/v1 -type f | wc -l)"
  echo
  echo "== ticket 227 AC box gauge (criteria-region scoped form; six boxes must stay unchecked) =="
  t=.scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
  echo "unchecked('- [ ] **AC#')=$(grep -c '^[[:space:]]*- \[ \] \*\*AC#' "$t" || true)  checked=[$(grep -c '^[[:space:]]*- \[x\] \*\*AC#' "$t" || true)]"
} > "$out" 2>&1
echo "written"
