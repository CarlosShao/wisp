#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/227v1-head.txt" "$T/common-20261018.txt" "$T/227v1-specific.txt" > "$T/227v1-dispatch.txt"
echo "227v1-dispatch = $(wc -l < "$T/227v1-dispatch.txt") lines / $(wc -c < "$T/227v1-dispatch.txt") bytes"
mkdir -p "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/227/v1/logs"
cat "$T/A1036.txt" >> docs/reports/pending-and-issues.md
git add -- docs/reports/pending-and-issues.md "$T/A1036.txt" "$T/227v1-head.txt" "$T/227v1-specific.txt" "$T/227v1-dispatch.txt"
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md)' || echo "  only-mine"
git commit -m "★机主澄清落账：「乙」＝D 盘清理那一问（⛔ 放宽分类）＋legacy 那 46 枚保留 ⇒ 放宽分类仍未答＝按默认只做甲（判级链一字未动、无需返工）｜记我一枚：拿到「1. 推吧 2. 乙」时⛔ 先对齐就动手搬 46 枚目录——两问都有「乙」这个选项、歧义真实存在，这次碰巧对是运气不是判断 ⇒ 定式＝凡机主以序号／单字回多问清单且选项名在两问里都成立，落盘前先一问一答对齐（同形若落在不可逆动作上＝事故，比如把「乙」读成放宽判级去改 blacklist.go 档位链）｜票 227 裁决腿 227-v1 派单写档＝待派（⛔ 未启动不计在飞；含我那枚 1 行号差要它去问、⛔ 替我编原因）｜★CPU 现量 97/94＞70 ⇒ 不放任何新腿；两枚裁决腿都活着在写盘（336/v1 119 枚最新 13:26、322/v2 39 枚最新 13:25）；322-v2-verdict.md 已在盘上但未提交＝在飞字节 ⇒ ⛔ 读它当凭据、⛔ 代它提交（今天两枚越界代提的根因＝在腿还活着时碰它的文件）｜账 A1036｜推送仍未执行（凭据闸已过、fast-forward 尺 0、未推 474）" -- docs/reports/pending-and-issues.md "$T/A1036.txt" "$T/227v1-head.txt" "$T/227v1-specific.txt" "$T/227v1-dispatch.txt"
git log -1 --format='%h'
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -vE '^(\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md)' || echo "  commit-clean: only mine"
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
