#!/bin/sh
# 票 340 普查腿 340-a1 · AC#1 三形代价表的现量原料（只读；一律取具名提交树）
# 逐字重跑：sh .scratch/wisp/probes/340/a1/scripts/ac1.sh
ANCHOR=feac8f77c2fbb15addbc36fb1bbaf13ed56588be

echo '### 1. 既存跨包装配先例（甲／乙"要不要新导出名"的对照面）'
echo '-- internal/tools 的测试包名册（⛔ 只有 in-package 才碰得到非导出件）'
for f in $(git ls-tree --name-only $ANCHOR:internal/tools | grep "_test\.go$"); do
  printf '%s  ' "$f"; git show "$ANCHOR:internal/tools/$f" | head -1
done | sort -k2 | awk '{print $2, $3, $4}' | sort | uniq -c
echo '-- 同时 import agent＋approval＋tools 的测试文件（一枚具名先例）'
git grep -l "wisp/internal/agent/approval" "$ANCHOR" -- internal/tools | while read -r f; do
  p=${f#*:}
  imports=$(git show "$ANCHOR:$p" | grep -c "CarlosShao")
  printf '%s  imports=%s\n' "$f" "$imports"
done
echo '-- 该先例的装配入口是否全为导出名（agent.New／approval.New／tools.New）'
git grep -n "agent\.New(\|approval\.New(\|tools\.New(" "$ANCHOR" -- internal/tools/loop_approval_test.go

echo
echo '### 2. 这条链上 CorrelationID 的搬运工（逐枚具名 file:line，判"真 corr 进真队列"要走几跳）'
git grep -n "CorrelationID: req.CorrelationID" "$ANCHOR" -- internal/tools/bridge.go
git grep -n "CorrelationID: callCorr" "$ANCHOR" -- internal/agent/loop.go
git grep -n "g\.q\.push(d)" "$ANCHOR" -- internal/agent/approval/gate.go
echo '-- PendingApproval（L2，走 push）与 PendingWindow（L1，不走 push）的边界'
git grep -n "^func (g \*Gate) PendingApproval\|^func (g \*Gate) PendingWindow" "$ANCHOR" -- internal/agent/approval/gate.go
echo '-- 桥侧把请求送进 Gate 的调用（L1／L2 各一枚）'
git grep -n "PendingApproval(\|PendingWindow(" "$ANCHOR" -- internal/tools/bridge.go

echo
echo '### 3. 叠改那一形的触发条件（同一枚 corr 二次入队，生产侧有几条路）'
echo '-- Gate.Replay（导出）⇄ Queue.Replay 的撞名支逐字'
git grep -n "^func (g \*Gate) Replay\|^func (q \*Queue) replay" "$ANCHOR" -- internal/agent/approval
git grep -n "#replay" "$ANCHOR" -- internal/agent/approval
echo '-- push 的撞名判定逐字'
git grep -n "if _, clash := q.byID\[corr\]; clash" "$ANCHOR" -- internal/agent/approval/queue.go
echo '-- d.CorrelationID 进 push 的上一跳（Decision 的 corr 来自 req.CorrelationID，⛔ 队列自己造）'
git grep -n "CorrelationID:" "$ANCHOR" -- internal/tools/bridge.go

echo
echo '### 4. 丙形那枚钉的形状先例（包内测试读自己的源文件）'
git grep -n "os.ReadFile(\"task.go\")" "$ANCHOR" -- internal/tools
echo '-- 该先例的用意（读码推的：内容锚，⛔ 绝对行号）＝包内测试读自己的源文件'
git grep -n 'os.ReadFile("task.go")' "$ANCHOR" -- internal/tools
echo '   本腿⛔ 打印整段别人件源文（那字节不归我入库）：只点名落点＋形状。'

echo
echo '### 5. 禁名词表那族仪器的射程（会不会咬本票三形）'
echo '-- 扫描目录（具名点）'
git grep -n 'filepath.Join(root, "internal", "panel")' "$ANCHOR" -- internal/panel/l2_grant_boundary_test.go
echo '-- 是否排除 _test.go（逐字尺命中）'
git grep -n 'strings.HasSuffix(p, "_test.go")' "$ANCHOR" -- internal/panel/l2_grant_boundary_test.go
echo '-- routeShapedName 的定义（"route 形状"到底是什么）'
git grep -nA14 "func routeShapedName" "$ANCHOR" -- internal/panel/l2_grant_boundary_test.go
echo '-- 词表本体逐字（两枚）'
git show "$ANCHOR:internal/panel/l2_grant_boundary_test.go" | sed -n "183,195p"
echo '-- 三形候选名逐个撞词表（算的：⛔ 任何一枚含下表词根）'
for n in PushForTest ExportPushForTest QueuePushForTest CallCorrForTest RewriteFormatAnchor PinsTheClashFormat; do
  low=$(printf '%s' "$n" | tr 'A-Z' 'a-z' | tr -d '_- .')
  hit=""
  for w in approve approval autoapprove grant granted allow allowed allowonce permit permitted ratify authorize authorised decision decide verdict outcome bypass override; do
    case "$low" in *"$w"*) hit="$hit $w";; esac
  done
  printf '%s -> 词根命中:%s\n' "$n" "${hit:- 无}"
done

echo
echo '### 6. 丙形的两枚去牙落点候选（⛔ 同一枚断言行两处落点）'
git grep -n "%s#%d" "$ANCHOR" -- internal/agent/approval/queue.go
git grep -n "corr = fmt.Sprintf\|d.CorrelationID = fmt.Sprintf" "$ANCHOR" -- internal/agent/approval/queue.go
echo '-- 名册侧两枚臂（叠改的另一半）'
git grep -n "HasPrefix(key, task+\"#\")" "$ANCHOR" -- internal/panel/subagent_roster_197.go
git grep -n "splitClashSuffix(key)" "$ANCHOR" -- internal/panel/subagent_roster_197.go

echo
echo '### 7. AC#5 前置料：本票落地后落到哪几包（点名，⛔ 跑）'
git ls-tree --name-only "$ANCHOR:internal/agent" | grep -c "_test\.go$"
git ls-tree --name-only "$ANCHOR:internal/agent/approval" | grep -c "_test\.go$"
git ls-tree --name-only "$ANCHOR:internal/panel" | grep -c "_test\.go$"
echo '-- 三包各自的顶层非测试文件枚数（成对台面的分母）'
for d in internal/agent internal/agent/approval internal/panel; do
  printf '%s  nontest_go=' "$d"; git ls-tree --name-only "$ANCHOR:$d" | grep "\.go$" | grep -vc "_test\.go$"
done

echo
echo '### 8. 本腿写面自检（⛔ 越格／⛔ 动 .go）'
git log --grep='票 340 普查腿 340-a1' --format='%H' | while read -r c; do
  echo "-- $c"
  git show --name-only --format= "$c"
done
echo '-- .go 命中枚数（逐笔名册并集，尺＝git show --name-only 全集）'
git log --grep='票 340 普查腿 340-a1' --format='%H' | while read -r c; do git show --name-only --format= "$c"; done | grep -c '\.go$' || true
