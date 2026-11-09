# 2026-11-09 票411 门禁三件（gofumpt scripts/ ／CR=0 ／d22scan；只建不删）
cd "$(git rev-parse --show-toplevel)" || exit 9
echo "== gofumpt scripts/ =="
/d/work/base/gopath/bin/gofumpt.exe -l scripts/
echo "gofumpt-done（上面无输出＝0 命中）"
echo "== CR check =="
for f in .scratch/wisp/probes/154/gate-clauses.sh .scratch/wisp/issues/411-gate-clauses-argv-overflow.md docs/evidence/s1/gate-clauses-argv-overflow.md .scratch/wisp/probes/154/r20261108/*.sh .scratch/wisp/probes/154/r20261108/*.md; do
	cr=$(grep -c "$(printf '\r')" "$f" || true)
	echo "$f CR=$cr"
done
echo "== d22scan =="
bash scripts/d22scan.sh
echo "d22scan-rc=$?"
