#!/bin/sh
# 票 338 落地腿 338-r1 · AC#2 四道门（共享工作树；★先 unset GOTMPDIR TMPDIR＝票 311 AC#1 shape B 卷闸）
set -u
cd "/d/work/workspace/projects plans/Wisp" || exit 9
unset GOTMPDIR TMPDIR
echo "GOTMPDIR=[${GOTMPDIR:-unset}] TMPDIR=[${TMPDIR:-unset}]  (unset 后)"
echo "HEAD=$(git rev-parse --short HEAD)"

echo '===== [1] sh scripts/d22scan.sh ====='
sh scripts/d22scan.sh > /d/tmp/wisp/20261018/338r1/ac2-d22scan.txt 2>&1
echo "rc_d22scan=$?"
tail -25 /d/tmp/wisp/20261018/338r1/ac2-d22scan.txt

echo '===== [2] sh scripts/check-path-length-budget.sh ====='
sh scripts/check-path-length-budget.sh > /d/tmp/wisp/20261018/338r1/ac2-len.txt 2>&1
echo "rc_len=$?"
tail -25 /d/tmp/wisp/20261018/338r1/ac2-len.txt

echo '===== [3] gofumpt -l internal cmd（交件侧复量，逐名） ====='
/d/work/base/gopath/bin/gofumpt.exe -l internal cmd | sed 's/\\/\//g' | sort
echo "rc_gofumpt_list=$?"

echo '===== [4] go vet ./internal/tools/ ====='
go vet ./internal/tools/ > /d/tmp/wisp/20261018/338r1/ac2-vet-tools.txt 2>&1
echo "rc_vet_tools=$?"
tail -30 /d/tmp/wisp/20261018/338r1/ac2-vet-tools.txt

echo '===== [5] go vet ./cmd/wisp/ ====='
go vet ./cmd/wisp/ > /d/tmp/wisp/20261018/338r1/ac2-vet-wisp.txt 2>&1
echo "rc_vet_wisp=$?"
tail -30 /d/tmp/wisp/20261018/338r1/ac2-vet-wisp.txt

echo '===== [6] CRLF 四枚 + panel 枚：门跑完后再验一字未动 ====='
for f in internal/agent/approval/pending_read.go internal/agent/tools.go internal/risk/provenance.go cmd/wisp/models.go cmd/wisp/panel_transport_35r2_test.go; do
  echo "$f | wt_sha=$(sha256sum < "$f" | cut -d' ' -f1)"
done
