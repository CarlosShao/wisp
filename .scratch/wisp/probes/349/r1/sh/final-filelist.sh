#!/bin/sh
# 交件必报：逐枚绝对路径＋字节数＋行数（三枚产码文件＋两枚报告＋票面）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
T=.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md
{
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go \
         .scratch/wisp/probes/349/r1/report.md .scratch/wisp/probes/349/r1/report-sec-b.md "$T"; do
  printf '%9s bytes  %5s lines  D:/work/workspace/projects plans/Wisp/%s\n' "$(wc -c < "$f")" "$(wc -l < "$f")" "$f"
done
echo "-- blob sizes at delivery HEAD (the three production files) --"
for f in blacklist provenance syncdirs; do
  printf '%9s bytes  blob %s  internal/risk/%s.go\n' "$(git cat-file -s "HEAD:internal/risk/$f.go")" "$(git rev-parse "HEAD:internal/risk/$f.go")" "$f"
done
echo "-- evidence dir totals --"
printf 'files=%s nonzero=%s zero=%s out-extension=%s logs=%s sh=%s\n' \
  "$(find .scratch/wisp/probes/349 -type f | wc -l)" \
  "$(find .scratch/wisp/probes/349 -type f -size +0c | wc -l)" \
  "$(find .scratch/wisp/probes/349 -type f -size 0 | wc -l)" \
  "$(find .scratch/wisp/probes/349 -type f -name '*.out' | wc -l)" \
  "$(find .scratch/wisp/probes/349/r1/logs -type f | wc -l)" \
  "$(find .scratch/wisp/probes/349/r1/sh -type f | wc -l)"
echo "-- ticket box ruler --"
printf 'unchecked=%s checked=%s\n' "$(grep -c '^- \[ \] \*\*AC#' "$T")" "$(grep -c '^- \[x\]' "$T" || true)"
}
