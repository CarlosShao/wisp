#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/A1030.txt" >> docs/reports/pending-and-issues.md
printf '%s\n' "322-r2 收工确认（10 笔）＋★抓出收件尺盲区：严格前缀尺对"主题被编码打坏成逐字 ?"的那笔失明 ⇒ 我数 9、盘上真 10（82c4a1d4，正文同坏、ASCII 完好＝PowerShell 侧编码落空）｜复点尺＝逐枚哈希点名（台件 recount-322r2.sh）｜全仓扫今日仅此一枚 ⇒ 其余五枚腿枚数不受影响｜定式升级：腿报枚数与我的尺不符时，除"先怀疑我的尺"外必须再跑一把逐枚哈希点名尺；前缀尺只做收集、不做分母｜内容尺复认：去重 83 枚（逐笔相加 90＝5 枚件被多笔重复触及，算术自洽）、.go 1、越界 3 枚全部已入账（同一根因两次、腿两次自抓自报当场改）｜票 322 未勾仍 1（AC#4）⛔ 翻框 ⛔ -done｜★本波写腿恢复 3≥3（336-r1∥227-r1∥338-r1）＝§4.0en 的例外状态已解除（原句不抹）｜在飞 4、已收 5、CPU 48/51/57、cmd/wisp 车道空已放行 338-r1｜账 A1030｜推送已授权仍未执行（未推 428）" > "$T/40eo-msg.txt"
git add -- docs/reports/pending-and-issues.md "$T/A1030.txt" "$T/40eo-msg.txt" "$T/recount-322r2.sh"
echo "=== staged audit ==="
git diff --cached --name-only | grep -v -E '^(\.scratch/wisp/probes/orch/20261018/|docs/reports/pending-and-issues\.md)$' || echo "  only-mine"
git commit -F "$T/40eo-msg.txt" -- docs/reports/pending-and-issues.md "$T/A1030.txt" "$T/40eo-msg.txt" "$T/recount-322r2.sh"
echo "=== self-check ==="
git show --name-only --format='' HEAD | grep -v -E '^\.scratch/wisp/probes/orch/20261018/' | grep -v '^docs/reports/pending-and-issues.md$' || echo "  commit-clean: only mine"
git log -1 --format='%h'
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
