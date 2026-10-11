#!/usr/bin/env bash
# 票 304 AC#4 定死的那段闸门名册，用它当场跑一发（编排者自跑，只读）
set -u
cd "/d/work/workspace/projects plans/Wisp"
OUT=.scratch/wisp/probes/orch/20261011/logs/orch-A85x-gate-live.txt
: > "$OUT"
echo "ts=$(date '+%F %T')  HEAD=$(git rev-parse --short HEAD)  branch=$(git rev-parse --abbrev-ref HEAD)" >> "$OUT"
echo "ruler=tasklist //FI \"IMAGENAME eq <name>\" //FO CSV  (命中行数 = 去掉表头的行数)" >> "$OUT"
echo >> "$OUT"
for n in wisp.exe balldebug.exe mockllm.exe wisp.test.exe; do
  body=$(tasklist //FI "IMAGENAME eq $n" //FO CSV 2>/dev/null | tail -n +2); trc=$?
  raw=$(tasklist //FI "IMAGENAME eq $n" //FO CSV 2>&1 | head -1 | cut -c1-40)
  cnt=$(printf '%s' "$body" | grep -c . || true)
  [ -z "$body" ] && cnt=0
  echo "$n count=$cnt tasklist_rc=$trc first_line=\"$raw\"" >> "$OUT"
  if [ "$cnt" -gt 0 ]; then
    pids=$(printf '%s\n' "$body" | cut -d'"' -f2 | tr '\n' ' ')
    echo "  pids=$pids" >> "$OUT"
    for p in $pids; do
      cl=$(powershell.exe -NoProfile -Command "Get-CimInstance Win32_Process -Filter 'ProcessId=$p' | ForEach-Object { \$_.CommandLine }" 2>/dev/null | cut -c1-150 || true)
      pp=$(powershell.exe -NoProfile -Command "(Get-CimInstance Win32_Process -Filter 'ProcessId=$p').ParentProcessId" 2>/dev/null | tr -d '\r' || true)
      ppalive=$(powershell.exe -NoProfile -Command "if (Get-Process -Id $pp -ErrorAction SilentlyContinue) { 'parent-alive' } else { 'parent-gone' }" 2>/dev/null | tr -d '\r' || true)
      echo "  pid=$p parent=$pp $ppalive cmdline=$cl" >> "$OUT"
    done
  fi
done
echo >> "$OUT"
web=$(powershell.exe -NoProfile -Command "(Get-CimInstance Win32_Process -Filter \"Name='msedgewebview2.exe'\").Count" 2>/dev/null | tr -d '\r' || true)
echo "msedgewebview2.exe count=$web  （只看枚数，归属另按 --webview-exe-name= 逐树指，⛔ 由裸枚数生成停法账）" >> "$OUT"
tmp=$(powershell.exe -NoProfile -Command "(Get-ChildItem -Path \$env:TEMP -Directory | Where-Object { \$_.Name -like 'Test*' }).Count" 2>/dev/null | tr -d '\r' || true)
echo "TEMP Test* 目录 count=$tmp" >> "$OUT"
# gate_rc＝这一发里“尺本身没打到”的次数（0＝四枚 tasklist＋两枚 CimInstance 全部拿到读数；⛔ 拿空输出冒充 0 枚）
gate_rc=0
for v in "$web" "$tmp"; do case "$v" in ''|*[!0-9]*) gate_rc=$((gate_rc+1));; esac; done
echo "gate_rc=$gate_rc  （0＝每一把尺都取到了数字读数；非 0＝有尺空手，具名看上面哪一行缺数）" >> "$OUT"
echo "lines=$(wc -l < "$OUT") bytes=$(wc -c < "$OUT")"
