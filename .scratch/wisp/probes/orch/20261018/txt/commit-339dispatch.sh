#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
cat "$T/339r1-head.txt" "$T/common-20261018.txt" "$T/339r1-specific.txt" > "$T/339r1-dispatch.txt"
echo "339r1-dispatch = $(wc -l < "$T/339r1-dispatch.txt") lines / $(wc -c < "$T/339r1-dispatch.txt") bytes"
mkdir -p "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/339/r1/logs"
git add -- "$T/339r1-head.txt" "$T/339r1-specific.txt" "$T/339r1-dispatch.txt" "$T/commit-A1034-35.sh" "$T/commit-A1034-35.sh"
echo "=== staged audit ==="
git diff --cached --name-only | grep -vE '^\.scratch/wisp/probes/orch/20261018/' || echo "  only-mine"
git commit -m "票 339 落地腿派单写档＝待派（⛔ 未启动、不计在飞）：head＋specific＋dispatch 共 3d26240c 形三枚，落点 D:\tmp\wisp\20261018\339r1\；按住理由具名＝在飞的 336-v1 正自己重跑 internal/risk 那把 8 发去牙扫描，本票落地腿同包跑两平台成对四数会互相洗脏（本包在册争用假红件 TestResolvePerCallBudget 要求安静复量＝票 228 AC#6 逐字在册）；释放条件＝336-v1 交件并收件完成⇒同时释放 internal/risk 测试面；★容器台件配方补一味进派单＝Windows 宿主路径进 docker -v 要 MSYS_NO_PATHCONV=1（否则 /src 被改写成 D:/src＝容器起不来）" -- "$T/339r1-head.txt" "$T/339r1-specific.txt" "$T/339r1-dispatch.txt" "$T/commit-A1034-35.sh"
git log -1 --format='%h'
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
