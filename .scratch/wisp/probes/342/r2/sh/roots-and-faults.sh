#!/bin/sh
# 342-r2 root classification + construction-fault scan, with patterns tolerant of the doubled
# backslashes Go prints inside quoted paths, and without hiding any grep exit code behind a pipeline.
cd /d/tmp/wisp/20261019/342r2/logs || exit 1
echo "== temp root per log (D-gotmp / LOCALAPPDATA\\Temp / USERPROFILE-direct) =="
for f in *.log; do
  d=$(grep -c '342r2.gotmp' "$f"); s=$(grep -c 'AppData.*Local.*Temp' "$f"); p=$(grep -c 'C:...Users...swq...Test' "$f")
  printf '%-28s D-gotmp=%-4s LOCALAPPDATA=%-4s USERPROFILE=%-4s\n' "$f" "$d" "$s" "$p"
done
echo
echo "== construction-fault scan (A1053 family = mklink /J failing, or a staged name that is not a
reparse point) =="
hits=0
if grep -n 'mklink /J' *.log > /d/tmp/wisp/20261019/342r2/logs/mklink-mentions.txt; then
  echo "mklink /J mentioned in $(wc -l < /d/tmp/wisp/20261019/342r2/logs/mklink-mentions.txt) rows (a junction
failure would name it here):"
  grep -n 'mklink /J .*: exit status\|mklink /J.*cannot create\|Cannot create a file when that file already exists' *.log
  echo "mklink-failures=$?"
fi
if grep -n 'is not a reparse point' *.log > /d/tmp/wisp/20261019/342r2/logs/reparse-faults.txt; then
  echo "REPARSE-FAULT rows below:"; cat /d/tmp/wisp/20261019/342r2/logs/reparse-faults.txt
else
  echo "REPARSE-FAULT-ROWS=0 across all logs"
fi
if grep -n 'cannot construct' *.log > /d/tmp/wisp/20261019/342r2/logs/skip-faults.txt; then
  echo "SKIP/CONSTRUCTION rows below:"; cat /d/tmp/wisp/20261019/342r2/logs/skip-faults.txt
else
  echo "SKIP-ROWS=0 (no instrument skipped a construction on this machine)"
fi
echo
echo "== SKIP count per paired full run (acceptance has to count these) =="
for f in base-full-2.log land-full-2.log base-full-3.log land-full-3.log landminus-full-1.log; do
  printf '%-24s SKIP=%s\n' "$f" "$(grep -c '^--- SKIP' "$f")"
done
