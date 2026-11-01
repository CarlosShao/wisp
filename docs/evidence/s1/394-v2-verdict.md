# 票 394 终裁腿 `394-v2` —— risk 层 mark 观察缝（第三形）终裁（非实现者）

> **尺**：实现件 `27bda97e0`（7 枚 +675/-4）。读数全部取自仓外导出树 `D:\tmp\wisp\20261101\394v2\exp`（`git archive 27bda97e0` 纯 commit 态：无 .git、无共享树触碰、在飞腿零归因）。共享树全程只读（开工核 `git status` tracked 零脏，HEAD=`ba0299245`）。staged 配方＝382-w-r2/394-r1 runner 配方（no-tag 二进制＋三 DLL 同目录＋五 `WISP_*` env＋`-test.timeout 30m`＋cwd=包目录；Git Bash 跑，rc 重定向现取，每跑 `unset GOTMPDIR`）。⛔ push。

## 0 一句话判语
**第三形实现正确、生产真挂、四 mark 点全覆盖、读数逐位复现、反向两发有牙、零 394 归因红 ⇒ 票面五格全勾，可收口（带三条已定性残余）。** R382-1 的「生产 RecordTaintInheritance 调用方」缺口就此闭合：mark 时累积归一化内容的义务由缝在 risk 内履行（382-w-r2 停手②收回）。

## 1 选型核（清单①）——亚形 B 三条理由在码上全部成立
1. **接线顺序**：`ProvOptions`（provenance.go:194-234）是纯构造期配置、无任何观察者面，`NewProvenance`(:298) 定格；而生产闸只在 `attachProvenance`（resident_voice_windows.go:615）才存在，引擎更早（run.go:868）已建——亚形 A 必须在组合根加晚绑定管道；亚形 B 的 `NewTTSGate(prov, cfg)` 天然持有引擎，382-w-r2 已落接线（:615 原样）零改动生效（§4 读数即证）。
2. **未导出面**：markseam_394.go 全部新标识符未导出（`taintMarkObserver`/`addMarkObserver`/`notifyTaintMark`/`observeTaintMark`）；provenance.go 只加未导出字段 `markObs`＋未导出调用；全 diff 零新导出符号 ⇒「契约面只在 risk」逐字达成。
3. **internal/tools 零触碰**：commit 7 枚无 tools 路径；tools 最后落库 79041f071 先于本票；commit 树 tools.Options（bridge.go:32 起）无 mark 观察钩（`git grep OnMark`=0）；4 处 mark 点 0 字节。
锁纪律码核：通知锁内拷贝、锁外调用，观察者 g.mu 不与引擎锁倒序；零观察者一次 len() 读零分配；`go test -race ./internal/risk -run Test394` rc=0。判：AC#1 勾。

## 2 四 mark 点覆盖核（清单②）——汇流口主张坐实
**关键核点**：通知只加在 `MarkWithHostPath`（provenance.go:652，落账后、p.mu 外），而三处生产点写的是裸 `Mark`——**`Mark`＝`MarkWithHostPath(...,"")` 薄包装（provenance.go:483-485），四点全部流经汇流口**，主张成立。
- bridge.go:687 `MarkWithHostPath`／subagent_197.go:512 `Mark`／projctx.go:349 `Mark`——同一台引擎 rt.c25（run.go:868/869 唯一引擎；:896/:917/:1374 交接），闸借同一台（resident_windows.go:579 `prov=src.run.c25` 到 :615）＝三处直喂。
- panel_assets.go:245：`detector()`(:235-248) 现建私有引擎（`ProvOptions{NoProbe:true}`）；`wisp panel-assets`＝main.go:112 派发的一次性 CLI 诊断面，进程无播报路径无闸。**判维持：缝照常触发（同汇流口、同归一化通知）但不硬造消费方**——硬造被消费方为零的闸正是 382-w-r2 停手②拒绝的假接线；停手条款①未触发（缺的是该进程本就没有的播报，不是观察面够不着）。

## 3 留存姿势追认核（清单③）——复核成立
1. **upsert 码核**：`observeTaintMark` 每作用域一条 link（按 TaskScopeID 查找；TaintContent/ToolSources 平行追加、MarkedAt 更新末次）；`Test394SeamUpsertsOneLinkPerScope` 钉住（同 scope 两次 mark＝1 link 2 taints；显式 record 仍追加绝不并进缝的 link）。
2. **漏毒论证前提为真**：`checkAncestorScopes` 两轮扫描都受 `i < MaxAncestorDepth`（tts_gate.go:179/:195，Default=16）约束、链无淘汰 ⇒ 逐 mark 追加会让第 17 个 mark 起的生产毒永久落扫描窗外（保护性失聪）；upsert 把链增长压到 O(不同已 mark 作用域数)。382-w-v2「只收任务作用域」建议的否答（risk 内无任务作用域判别器／前缀过滤=新脆弱契约／不过滤是 fail-closed 且永闭作用域的链是它到达播报裁决的唯一通路——:193-203 全链 fallback 扫描坐实）核过成立。
3. **RecordTaintInheritance 逐字节未动**（tts_gate.go:272-296；diff hunk 只在 NewTTSGate 区，函数体无 hunk）；**AncestorInheritance 1→2＝已批姿势可见面**（同 scope 两次 prov.Mark 经缝 upsert 折 1 条＋显式 record 1 条＝2，注释钉明）。编排者已追认留存姿势＝本腿复核论证（非重开）。

## 4 生产真挂读数复跑（清单④）——逐位复现
rig＝382wr2 腿原样（ttsInheritLeg 走真 attachProvenance），纯 `prov.Mark` 喂链、零手动 record：
- `SEAM-PRODUCTION refused=1 passed=0 spoken=0 samples=0 slot_loads=0 links=1`（Test394R1 两案 PASS）——与实现者读数逐位同。
- **normalizeTaint 双钉**：staged 案全角原文 mark vs 半角播报（raw 对 raw 零共享窗，只有缝内存的归一化副本才可能拒）＋单元侧白盒断言链内容==normalizeTaint(raw)（Test394SeamRecordsNormalizedContentAtMarkTime 绿）。
- 382WR2 三案不退化：INHERIT refused=1（slot_loads=0 合成前拒）／CLEAN passed=1 spoken=1（真 matcha 合成，本程 samples=75286；实现者带 75326~75473，同形状抖动）／DIRECT refused=1。

## 5 AC#3 反向两发复刻（清单⑤）——牙是活的
突变在导出树文件级（备份→摘→跑→还原；还原后与原伴字节校验 0 差异；等效 -overlay 法，零触共享树）：
| 突变 | 摘什么 | 红相（本腿实测） |
|---|---|---|
| M1 | provenance.go:652 汇流口通知 | risk 层 rc=1 恰 5 红（4 缝测试＋AncestorInheritance；NotRecorded/Disabled 独立案绿）；staged 层 Test394R1 rc=1：`the seam recorded 0 links`＋`timed out waiting for the seam-fed gate to refuse`＋**污句真合成上设备 `runes=27 samples=81386`** |
| M2 | tts_gate.go NewTTSGate 自注册 | risk 层红相同形（同 5 红）；staged 层同形（0 links＋拒播超时＋`runes=27 samples=81377`） |
| 缝独立性 | 两发下跑 382WR2 | 三案全绿（INHERIT refused=1／CLEAN spoken=1／DIRECT refused=1）——显式 record 通路与缝互不依赖 |
| 还原核 | 突变前 pristine 二进制复跑 | risk rc=0＋staged **5/5 PASS**：喂数 links=1＋SEAM-PRODUCTION 逐位复现＋INHERIT/CLEAN(75363)/DIRECT |

## 6 零回归复核（清单⑥）——无 394 归因红
- risk 全套 rc=0；`go test -race ./internal/risk -run Test394` rc=0。
- 六族 staged 合跑两遍 **33/2/0＋33/2/0**，红集逐枚归因：run1 的 Test41a2r1 两红＝**导出树缺真仓未 tracked 的 punct 模型**（third_party\model-fetch\x-punc\…，rig 按仓相对路径取料；junction 后同案绿——环境缺料非回归）；run2 的 Test384ResidentWatchdogEndToEnd＝A1329 既有合跑时窗形状（基线唯一红同案；run1 该案反绿＝负载摆动，两侧读数都在账）；run2 的 Test15eSpeechEndToEnd＝球投影时序断言（链路本身 handed=1 成功），**单跑复测 PASS**。
- `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 两次现取逐位同（六族 run1＋15e 单跑）。
- d22scan（exp 树一发含正控）：**clean rc=0**；正控 PASS=35 FAIL=0 SKIP=0（RUN=77）＝基线同值；**四分母 288/52/734/184 逐位一致**。逐点归因（git ls-tree 算术核）：bans#1-5 internal 287（382-w-r2 起值）+1＝markseam_394.go＝288；cmd 53 到 53（walker 面 52 零变；本腿 cmd 新增为 _test 件不入 production 面）；ban#8 internal 732（父树＝729 起点＋两腿之间他人落库 3 件）+2（markseam_394.go＋其 _test）=734；ban#8 cmd walker 面 183+1（本腿测试件）=184。注：exp 无 .git，扫描器 loud 提示 gitignore 未应用，design/ 计数 14（真仓 39 含未 tracked 设计件，非本票口径，internal/cmd 四数不受影响）；docs/ 不在任何 ban 作用域，未扫 docs/。
- tagged（webrtcap）形状未跑：本腿零触碰音频面（与实现者同口径）；26p (a)(b) 不在派单六族清单。

## 7 五格逐条判＋收口＋残余（清单⑦）
| 格 | 判 | 依据 |
|---|---|---|
| AC#1 观察缝落地 | 勾 | 第三形＝机主 A1362 批形状；亚形 B 理由三条码上成立；未导出面/零 tools 触碰逐字（§1） |
| AC#2 端到端 | 勾 | 任务污点→缝 accumulate→播报拒，拒在合成前（slot_loads=0）；生产真挂读数逐位复现（§4） |
| AC#3 反向判据 | 勾 | 两发红相精确、382WR2 独立绿、还原 0 差异＋5/5 复现（§5） |
| AC#4 零回归 | 勾 | risk 全套 rc=0＋race rc=0＋六族无 394 归因红＋四分母逐位一致＋d22scan clean（§6） |
| AC#5 纪律格 | 勾 | 读数全在仓外导出树；共享树零触碰零写入；票面/台账/契约/thresholds/golden 零触碰；⛔ push |

**收口判：可收口**（三条残余全部如实定性，均为既有性质或显式边界，非本票缺陷）：
1. **注册前 mark 不入链**：缝从 NewTTSGate 注册起观察，boot 早期（闸建立前）的 mark（如启动期 projctx 盖戳）不追溯——其毒仍在各自作用域内活着（与今日保护面同、无回退）；追溯需内容读出口，而 ScopeTaints 按安全设计只回元数据，不动。
2. **缝日志走 logf**（宿主注入汇，nil=drop）：与 mark 路径既有日志惯例一致；staged 下缝日志不可见是路由使然，链数断言才是缝触发判据。
3. **MaxAncestorDepth(16) 首-16 扫描帽**：upsert 后链长＝不同已 mark 作用域数，超长会话（>16 个不同作用域）仍触帽——382-w-r1 既有防御参数的既有性质，本票未动 CheckText 语义；如需演化归 382 契约面另开票。

## 8 裁腿自报（环境与偏差，读数解释力）
1. 读数树＝`git archive 27bda97e0` 纯 commit 态；与实现者「共享树合并态」读数不同源但同值，归因更干净。
2. 导出树补料两笔（真仓未 tracked 本地料，复制/junction 进 exp，不触共享树）：`build/models/manifest.json(+.minisig)`（WISP_MODELS_MANIFEST 解析需）＋`third_party/model-fetch` junction（41a2r1 rig 取 punct 模型）。第一遍六族 33/2/0 的两红即缺料所致——**建议后续 staged 腿的 runner 配方把这两笔列为显式前提**。
3. 突变＝导出树文件级（备份＋还原＋字节校验），与实现者 -overlay 等效；两法都零触共享树。
4. 仓外产物 `D:\tmp\wisp\20261101\394v2\**`（exp/bin/logs/cases/mut/run394v2.sh，只建不删）。
