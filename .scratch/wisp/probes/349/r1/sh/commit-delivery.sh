#!/bin/sh
# 349-r1 交件笔：全部证据件＋票尾 bullet 入库。
# ★公共块：add 与 commit 两处逐枚**字面** pathspec（本清单＝我手写的字面路径、⛔ 从 index／目录反推、
#   ⛔ xargs、⛔ 通配）；★新文件必须先 git add 才进得了 commit pathspec（编排者今天栽过两次）。
# 提交后现跑 git log -1＋git show --name-only＋git status 对回，⛔ 拿"命令没报错"当已落库。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 90
P1=.scratch/wisp/probes/349/r1
T=.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md
PATHS="
$T
$P1/report.md
$P1/report-sec-b.md
$P1/report-head-part.txt
$P1/report-sec-b-body.txt
$P1/tmp-plus.txt
$P1/msg/c1-blacklist.txt
$P1/msg/c2-syncdirs.txt
$P1/msg/c34-provenance.txt
$P1/msg/c6-delivery.txt
$P1/msg/ticket-bullet-349.md
$P1/sh/ac1-proof.sh
$P1/sh/append-bullet.sh
$P1/sh/build-tree.sh
$P1/sh/census-349r1-main.go.txt
$P1/sh/census-run.sh
$P1/sh/check-banned-in-added.sh
$P1/sh/collect-logs.sh
$P1/sh/commit-delivery.sh
$P1/sh/fill-empty-readings.sh
$P1/sh/final-hygiene.sh
$P1/sh/fix-1-blacklist-fires.sh
$P1/sh/fix-2-syncdirs-complete.sh
$P1/sh/fix-34-provenance.sh
$P1/sh/fix-mojibake.sh
$P1/sh/hygiene-scan.sh
$P1/sh/make-trees.sh
$P1/sh/msg-hygiene.sh
$P1/sh/run-ac3-wordlist.sh
$P1/sh/run-ac3.sh
$P1/sh/run-gates.sh
$P1/sh/run-pair-linux.sh
$P1/sh/run-pair-windows.sh
$P1/sh/run-teeth-336.sh
$P1/sh/scrub-banned.sh
$P1/sh/strip-sha.sh
$P1/logs/00-make-trees.log
$P1/logs/ac1/01-ac1-blacklist-fires.log
$P1/logs/ac1/02-ac1-syncdirs-complete.log
$P1/logs/ac1/03-ac1-provenance-two.log
$P1/logs/ac1/04-ac1-provenance-disposal.log
$P1/logs/ac3/ac3-misfire.log
$P1/logs/ac3/ac3-outer.log
$P1/logs/ac3/ac3-outer2.log
$P1/logs/ac3/ac3-wordlist-outer.log
$P1/logs/ac3/ac3-wordlist.log
$P1/logs/census/ac0-base-four-names.txt
$P1/logs/census/ac0-base-full.err
$P1/logs/census/ac0-base-full.txt
$P1/logs/census/ac0-base-prod.err
$P1/logs/census/ac0-base-prodonly.txt
$P1/logs/census/ac0-del-four-names.txt
$P1/logs/census/ac0-del-full.err
$P1/logs/census/ac0-del-full.txt
$P1/logs/census/ac0-del-prod.err
$P1/logs/census/ac0-del-prodonly.txt
$P1/logs/gates/g1-after.txt
$P1/logs/gates/g1-before.txt
$P1/logs/gates/g1-gofumpt-mine.txt
$P1/logs/gates/g1.txt
$P1/logs/gates/g1b-gofumpt-startblobs.txt
$P1/logs/gates/g1c-gofumpt-delblobs.txt
$P1/logs/gates/g2.txt
$P1/logs/gates/g3.txt
$P1/logs/gates/g4.txt
$P1/logs/gates/gates-outer.log
$P1/logs/pair/lin-base.log
$P1/logs/pair/lin-del.log
$P1/logs/pair/lin-pair-runner.log
$P1/logs/pair/lin-red-base.txt
$P1/logs/pair/lin-red-del.txt
$P1/logs/pair/red-base.txt
$P1/logs/pair/red-del.txt
$P1/logs/pair/win-base-raw.log
$P1/logs/pair/win-base.log
$P1/logs/pair/win-del-raw.log
$P1/logs/pair/win-del.log
$P1/logs/pair/win-pair-runner.log
$P1/logs/teeth/ac2-00-pristine-leg.log
$P1/logs/teeth/ac2-00-pristine-teeth.log
$P1/logs/teeth/ac2-MU-1a-leg.log
$P1/logs/teeth/ac2-MU-1a-teeth.log
$P1/logs/teeth/ac2-MU-1b-leg.log
$P1/logs/teeth/ac2-MU-1b-teeth.log
$P1/logs/teeth/ac2-MU-1c-leg.log
$P1/logs/teeth/ac2-MU-1c-teeth.log
$P1/logs/teeth/ac2-MU-1d-leg.log
$P1/logs/teeth/ac2-MU-1d-teeth.log
$P1/logs/teeth/ac2-MU-1e-leg.log
$P1/logs/teeth/ac2-MU-1e-teeth.log
$P1/logs/teeth/ac2-MU-1f-leg.log
$P1/logs/teeth/ac2-MU-1f-teeth.log
$P1/logs/teeth/ac2-MU-1g-leg.log
$P1/logs/teeth/ac2-MU-1g-teeth.log
$P1/logs/teeth/ac2-MU-1h-leg.log
$P1/logs/teeth/ac2-MU-1h-teeth.log
$P1/logs/teeth/ac2-sweep-runner.log
$P1/logs/teeth/runner-outer.log
"
# shellcheck disable=SC2086
git add $PATHS
echo "ADD_RC=$?"
cp -f "$P1/msg/c6-delivery.txt" /d/tmp/wisp/20261019/349r1/msg/c6-delivery.txt
# shellcheck disable=SC2086
git commit -F "D:\tmp\wisp\20261019\349r1\msg\c6-delivery.txt" -- $PATHS
echo "COMMIT_RC=$?"
echo "== verify =="
git log -1 --format='%h %s'
git show --name-only --format='COMMIT %h' HEAD | grep -vE "^$|^COMMIT" | sed 's/^/  /' | head -30
echo "  ...(truncated display; full count below)"
git show --name-only --format='' HEAD | grep -c .
echo "== out-of-scope check on every one of my commits =="
A=7b633fa9713b0514806c8b808ed4db91b1485cdc
for c in bf7f13bc 07f178df db2f9eab 7d56962b; do
  echo "-- $c --"
  git show --name-only --format='' "$c" | grep . | grep -vE '^\.scratch/wisp/probes/349/r1/|^internal/risk/(blacklist|provenance|syncdirs)\.go$' || echo "  in-scope: all"
done
echo "-- HEAD --"
git show --name-only --format='' HEAD | grep . | grep -vE '^\.scratch/wisp/probes/349/r1/|^\.scratch/wisp/issues/349-' || echo "  in-scope: all"
echo "== go files count at HEAD commit (want 0 in this final one) =="
git show --name-only --format='' HEAD | grep -c '\.go$' || true
echo "== status after =="
git status --porcelain -- internal/risk .scratch/wisp/probes/349 | head -5
echo "(empty = clean)"
