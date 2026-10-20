#!/bin/sh
# 353-v1 R59 = 交件前最后一把：msg-09 件身 + 待提交两枚件的在位尺 + 主件引用的件名册尺（逐枚 cat-file -e 对 HEAD）
S7=$(printf '\xc2\xa7')
cd '/d/work/workspace/projects plans/Wisp' || exit 9
for f in /d/tmp/wisp/20261020/353v1/msg-09.txt docs/evidence/s1/353-v1-verdict.md .scratch/wisp/probes/353/v1/sh/58-final-precheck.sh; do
  printf '%-56s SECT=%s CR=%s LINES=%s BYTES=%s PH=%s\n' "$(basename "$f")" \
   "$(grep -o "$S7" "$f" | wc -l)" "$(tr -cd '\r' < "$f" | wc -c)" "$(wc -l < "$f")" "$(wc -c < "$f")" "$(grep -cE 'TODO|FIXME|XXX' "$f" || true)"
done
echo "== 主件正文引用的件名逐枚验在（HEAD 或 待提交）=="
for p in docs/evidence/s1/353-a1-census.md \
         .scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md \
         .scratch/wisp/probes/orch/20261020/txt/353v1-dispatch.txt \
         .scratch/wisp/probes/353/a1/logs/04-ac2-instrument-shape-draft.txt \
         .scratch/wisp/probes/353/r1/readings-ac2-baseline.md \
         .scratch/wisp/probes/353/r1/readings-ac2-landed.md \
         .scratch/wisp/probes/353/v1/00-ruler-register.md \
         .scratch/wisp/probes/353/v1/01-independent-readings.txt \
         .scratch/wisp/probes/353/v1/02-ticket-draft-secC.md \
         .scratch/wisp/probes/353/v1/90-truncated-nonzero-rows.txt \
         docs/evidence/s1/342-v3-verdict.md docs/evidence/s1/322-v1-acceptance.md \
         internal/winsec/nametopup_353_windows_test.go internal/winsec/pinnedleaf_342_windows.go \
         internal/winsec/ancestorwindow_342r2_windows_test.go internal/winsec/ancestorunlink_348_windows_test.go \
         internal/winsec/ancestorworld_342_348_windows_test.go internal/winsec/leafrace_342r2_windows_test.go \
         internal/winsec/openleaf_322_other_test.go internal/winsec/winsec_windows.go; do
  if git cat-file -e "HEAD:$p" 2>/dev/null; then echo "  AT_HEAD $p"; else echo "  !! ABSENT $p"; fi
done
