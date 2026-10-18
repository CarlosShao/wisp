#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
git add -- "$T/common-20261018.txt" "$T/338r1-amend-01.txt"
echo "=== staged audit (correct ruler: no trailing dollar anchor) ==="
git diff --cached --name-only | grep -vE '^\.scratch/wisp/probes/orch/20261018/' || echo "  only-mine"
git commit -m "公共块更正一枚帽的射程：每件<=160 行只管正文件（report/census/散件），不适用于 *_test.go 与产码 => 为凑行数压注释/删用例是禁手（注释＝本仓最贵的东西）；起因＝227-r1 12:1x 简报里写「正压到 160 行、只压注释」＝我这句话写得不够硬，当场叫停并记我一枚；166-r1 11:1x 也为这帽来问过、我当时没裁＝欠它一个答复（已补发）" -- "$T/common-20261018.txt"
git log -1 --format='%h %s' | cut -c1-70
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -vE '^\.scratch/wisp/probes/orch/20261018/' || echo "  commit-clean: only mine"
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
