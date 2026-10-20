#!/bin/sh
# gw-1 味① 第 4 把的**权威仪器**＝本仓 CI 逐步执行的那枚（ci.yml:172-181 gofumpt ⇄ ci.yml:230 attrib.sh --tracked-only）。
# ★⛔ 我自己另拼一把"过滤"尺当判据＝改用仓里在册的仪器；我 sh/04 那两把只作**并排的第二把独立尺＋存量具名登记**。
# 三发＝① --self-test（仪器自称的牙：两个方向都得响）② --tracked-only（＝CI 真跑的那形＝判据 A）③ both（＝形 B＝归因尺）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
mkdir -p "$LOGS"
OUT="$LOGS/gates-attest.txt"
: > "$OUT"
{
  echo "== gates-attest START anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
  echo "WHEN=$(date '+%F %T %z')"
  unset GOTMPDIR TMPDIR TEMP TMP
  echo "env GOTMPDIR=[${GOTMPDIR-<unset>}] TMPDIR=[${TMPDIR-<unset>}] TEMP=[${TEMP-<unset>}] TMP=[${TMP-<unset>}]"
  echo "which gofumpt (attrib.sh 自寻形)=$(go env GOPATH)/bin/gofumpt.exe"
  echo
  echo "== 1) attrib.sh --self-test（仪器的牙，两方向）=="
  sh .scratch/wisp/probes/161/r5/attrib.sh --self-test > "$LOGS/attest-selftest.txt" 2>&1
  echo "SELFTEST_RC=$?"
  tail -12 "$LOGS/attest-selftest.txt" | sed 's|^|  ST |'
  echo
  echo "== 2) attrib.sh --tracked-only（＝CI 真跑那形＝判据 A）=="
  sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only > "$LOGS/attest-tracked.txt" 2>&1
  echo "TRACKED_ONLY_RC=$?"
  cat "$LOGS/attest-tracked.txt" | sed 's|^|  TR |'
  echo
  echo "== 3) attrib.sh both（＝形 B＝归因尺；含逐枚归因行）=="
  sh .scratch/wisp/probes/161/r5/attrib.sh > "$LOGS/attest-both.txt" 2>&1
  echo "BOTH_RC=$?"
  cat "$LOGS/attest-both.txt" | sed 's|^|  BT |'
  echo
  echo "== 4) 独立第二把（CI 逐字 172-181 那形）＝gofumpt -l . tools/d22scan tools/mockllm 的原始输出逐枚 =="
  echo "  （读数已在 logs/gate4-r3-raw.txt；此处只数归因＝⛔ 拿枚数顶名册）"
  /d/work/base/gopath/bin/gofumpt.exe -l . tools/d22scan tools/mockllm > "$LOGS/gate4-r3-raw.txt" 2>&1
  echo "RAW_RC=$?"
  echo "RAW_LINES=$(wc -l < "$LOGS/gate4-r3-raw.txt")"
  echo "-- 生产面（非 .scratch/**）逐名列＝存量具名登记（⛔ 洗、⛔ -w）--"
  tr '\\' '/' < "$LOGS/gate4-r3-raw.txt" | grep '\.go' | grep -v '^\.scratch/' | sed 's|^|  STOCK |'
  echo
  echo "== 5) 逐枚 CR 尺（对**入库 blob**量＝CI 看到的字节；工作树那形另列）=="
  for p in cmd/wisp/models.go internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go; do
    WT=$(tr -cd '\r' < "$p" | wc -c)
    BL=$(git cat-file blob "HEAD:$p" | tr -cd '\r' | wc -c)
    echo "  CRFILE $p worktree_cr=$WT blob_cr=$BL"
  done
  echo "  （正控＝同一条 tr 尺对自造 2 枚 CR 件读出 2＝尺⛔ 瞎）"
  printf 'a\rb\r\nc\n' > "$R/posctrl/cr2.txt"
  echo "  POSCTRL_CR2=$(tr -cd '\r' < "$R/posctrl/cr2.txt" | wc -c)"
  echo
  echo "== END anchor =="
  git rev-parse HEAD | sed 's|^|HEAD_FULL_END=|'
  echo "WHEN_END=$(date '+%F %T %z')"
} >> "$OUT" 2>&1
cat "$OUT"
