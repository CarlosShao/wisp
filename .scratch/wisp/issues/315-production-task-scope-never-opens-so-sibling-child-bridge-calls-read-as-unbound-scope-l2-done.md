# 315 — 生产里"任务开始"不开 C25 作用域：子代理／兄弟任务的首批走桥调用会被判 unbound-scope L2，落成假审批卡或超时硬拒（来路＝票 237 的独立真形；**先活体复量再定真伪**）

- **Status**：**done**（2026-10-14 2026-10-14 11:51 编排者收口：AC#0 复现（`315-a1` 活体＋编排者亲跑）＋AC#1 修（`315-r1` 选案ⓐ 带四处设计原文引证；修前 31.2s 卡拒→修后 0.18s 放行；Case C 同修；钉子去门控常驻）＋`315-v1` 六格全成立 ⇒ 双勾改名 `-done`。候选残余（不阻塞）：并发兄弟形未单独量（机制上按任务 admission 同名修复、登记候选）。账 `A915`）。⚠ 原 Status 逐字留档（已过期、不删）：「- **Status**：**ready-for-agent（第一格＝活体复量；⛔ 不许先修）**（编排者 2026-10-14 2026-10-14 10:18 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）。」
- 来源＝非实现者验收腿 `237-v1`（件 `docs/evidence/s1/ticket237-237-v1-accept-v1.md`，`be0934b7`）：它在裁票 237 AC#3 时撞出并修好 **harness 侧**同形（harness `AdmitTask` 补 `OpenTask/CloseTask` 孩子才上桥），随后按令核了生产面静态链。
- 关联：票 151（-done，"关"那半——`admitTask` revoke 挂 `CloseTask`）· 票 154（-done，RESERVED：宿主自带 id 那一形没人关）· 票 237（-done，测试面）。⚠ 本票**只裁"开"这一侧的生产形状**，⛔ 不重开 151/154。

## 现量（`237-v1` 逐行核过；引用前先现读——树会漂）

1. 非测试 `OpenTask(` 全树仅 `internal/tools/bridge.go:649`（`b.mark` 里懒开，**仅敏感源结果收口时**）；定义 `:847`。
2. `cmd/wisp/run.go:953-962` 的 admitTask 只 `CloseTask(:960)`；`internal/agent/loop.go:373-379`、子环路 `internal/tools/subagent_197.go:311-320` 均**不开**。
3. 触发条件（`internal/risk/provenance.go:694-711` ＋ Detector `:858-869`）：**调用方 scope 未注册＋别的 scope 有 mark ⇒ 每次走桥判 unbound-scope R4/L2**。
4. 首发诱因：首枚子代理结论把 mark 落父名（`subagent_197.go:487`）后，兄弟孩子（工具面＝父桥 `:296`，dispatch 带自己 task id `loop.go:675-676`）即中。
5. 生产 L2 落真审批门（`run.go:612` approval.New）＝假审批卡／超时拒；`run.go:1052` 项目说明 scope（`ProjectInstructionsEnabled` 默认 true）整跑期间构成"别处有 mark"。
⚠ 交叉读数（151/154 时代）：154 记"唯一非测试 `OpenTask` 调用者＝`bridge.go:559`（mark 里）"——同一处、行号漂了；"懒开**仅敏感源收口**"这一窄化是 237-v1 的新读数，**以现读为准**。

## 要建什么

- [x] **AC#0 活体复量（先跑这格；量不到就地撤本票）**：一发真跑组合根／CLI 接缝的用例或最小驾驶脚本，触发"子代理的兄弟孩子工具调用"形；读数＝`unbound-scope` R4/L2／审批卡／拒绝句是否出现（命令＋时刻＋rc＋红句逐字；⛔ 不许用 mock 冒充）。若**量不到** ⇒ 撤票（口令「**315 撤**」）＋台账写明量法。 〔10-14 2026-10-14 11:51 编排者翻勾：`315-v1` 判"成立"（六格全成立、停手项零）；裁决表 `docs/evidence/s1/ticket315-315-v1-accept-v1.md`（`301db209`）；账 `A915`〕
- [x] **AC#1（按 AC#0 的读数裁）**：或**修**（三案各自代价表：ⓐ`admitTask` 补 OpenTask／ⓑ子环路自开／ⓒ桥懒开扩到每次 dispatch）或**登记**（判"有意为之"须给理由链）。⛔ 改判据方向／放宽断言一律禁。 〔10-14 2026-10-14 11:51 编排者翻勾：`315-v1` 判"成立"（六格全成立、停手项零）；裁决表 `docs/evidence/s1/ticket315-315-v1-accept-v1.md`（`301db209`）；账 `A915`〕

## 禁区

⛔ 不动 C19/C25 语义本身；⛔ 三枚冻结件／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`docs/PLAN.md`／`docs/specs/**` 一字不动；⛔ `frontend/**`／`design/**` 零读零写零转述；⛔ 零 push；⛔ 不许重开 151/154 的框。

## 排程

写面≈`cmd/wisp`／`internal/tools`／`internal/agent`；⚠ 与票 236 AC#3 腿撞 `internal/tools` ⇒ **串行**；本票腿起手先跑包级互斥尺。

## Progress log (append-only, newest last)

- [2026-10-14 10:18 +08] 编排者立票（来路＝`237-v1` 生产面读数；静态链五条已核；AC#0 活体复量为门；撤票口令「315 撤」）。账 `A905`。
- [2026-10-14 10:53 +08] 探针腿 `315-a1`（只读复量；锚 `a5b2938c`）：**AC#0＝复现**。载具＝新建 `cmd/wisp/zz315a1probe_test.go`（golden SSE 接缝真跑组合根，模型自发的 `task.spawn`/孩子环路自发的 `fs.read`；⛔ mockllm 载具过不了此形——它只在非 nil ToolChoice 下发 tool_call）。命令 `WISP_315A1=1 GOFLAGS= go test ./cmd/wisp/ -run 'TestZZ315A1SiblingChildFirstBridgeCall' -count=1 -v`（10:47:20→10:47:59 +08）**rc=1**（该 run 自身 rc=0/completed）：兄弟孩子乙（同父、在飞 Thinking）的首批走桥 `fs.read`（白名单内，audit `in_allowlist_scope=true`）被判 `L2 rules=[R4]`，卡逐字 `[确认 L2 fs.read] R4: 包含来自 unbound-scope scope is not open (no OpenScope for this id, or its own handle closed it) 的内容`＋`[dismissed] 审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝`，行 `L2/timeout/error class=user_rejected`；对照腿（甲同路径首调用，尚无 mark）＝L0/allow/success；种子 mark＝甲结论落父名（`task.output origin=subagent:<甲>` 19 runes），乙 scope 全程 `was_open=false`。件＝`.scratch/wisp/probes/315/a1/`（evidence.md＋2 枚日志 5,925/17,615 B；对照腿 Case A rc=0 原点为真环路首发）。
- [2026-10-14 11:2x +08] 写腿 `315-r1`（实现者；起手锚 `5ed56c2f`）：**AC#1 落修＝选案 ⓐ**——`cmd/wisp/run.go` 的 `admitTask` 在任务开始处补 `rt.bridge.OpenTask(taskID)`，与同钩子 revoke 里那枚 `CloseTask(taskID)` 同 id/同所有者成对（桥持 handle；零新接口、零新生命周期）。原文依据＝`internal/risk/provenance.go:453` "composition root calls this at task start"＋`:706` "ensure OpenScope at task start"＋票 151 验收件 §4 表行①（AdmitTask 是唯一 per-task 边界）＋票 237 harness 同形（`internal/tools/subagent_222_test.go:331-333`）的生产侧对称形；ⓑ 子环路自开＝盖不住根且要新路由（`internal/agent`→C25），ⓒ 桥每次 dispatch 懒开＝会在评估器前登记、使 unbound fail-closed 分支不可达（改 C19/C25 判定语义）⇒ 均排除。**读数**：修前红（未修树，11:07:29→11:08:09）B＝rc=1、`--- FAIL … (31.22s)`（卡＋`ANSWER-EXPIRED … decision=timeout->reject after=31s`，乙 close `was_open=false`）→修后绿（11:23:29→11:23:37）A 0.15s／B 0.18s／C 0.12s，乙走桥行逐字 `risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true`、三条 close 行全 `was_open=true`；突变（`-overlay` 只摘这一行，工作树零突变）B 红 31.23s／C 红 31.15s。**钉子去门控**：删除 `WISP_315A1` 门控——旧门控在 cli scope 严格 runner 下实测 `SKIP=1` 致命红 rc=1（`runtests.sh … -run '^…SiblingChild…$'`，ledger 无此名；去门控后无 SKIP、不用碰 174-r5 的 `scripts/portable-tests.sh` 写面），默认跑三形全绿（0.45s 合计）。**自然涵盖**＝a1 未试的 project-instructions 种子形**活体复到**（新增 Case C：数据目录 `AGENTS.md` 被 projctx 盖戳进 `run:project-instructions`，修前 31.15s 红→修后 0.16s 绿，根的首发 `fs.read` L0/allow/success）；并发兄弟形未单独量＝登记候选。**门禁**（改后终树）：build／vet／`gofumpt v0.12.0`／`gofmt`／`d22scan` clean／长度门 `VERDICT GREEN` 59=59／`Test222 -race`（3 PASS）／151 两枚关侧钉／197 自批族（无 `-race`）／158 开侧钉 全 rc=0；⚠ 既有 `-race` 夹具红（`Test197SubagentSelfApproval…`，栈顶 `bytes.(*Buffer).Write()`⇄`auditf()` `run.go:933`，同族已在台账）在 overlay 摘本修的树上同红＝非本修引入。件＝`.scratch/wisp/probes/315/r1/`（`00-anchor-and-gates.md` 167 行＋logs 21 枚＋mutation 2 枚；入库名册见该件 §8）。⛔ 未翻框、⛔ 未动 C19/C25 语义/阈值/golden/`allowlist.txt`、⛔ 未 push。
- [2026-10-14 11:4x +08] 验收腿 `315-v1`（**非实现者**；起手锚 `b7a15297`，窗口 11:32–11:45 +08，判读两文件全窗 sha 不变）：**判语＝六格全成立**——①选案 ⓐ＝用既有机制非改 C25 语义（非注释产码增行恰 1 条＝`rt.bridge.OpenTask(taskID)`；四处引证 `provenance.go:453/:706`＋151 验收件 `:192`＋`subagent_222_test.go:331-333` 逐字核过；ⓑ/ⓒ 结构排除复核成立）；②读数独立复现：默认三形 rc=0 A 0.15s／B 0.20s／C 0.14s、`[确认 L2` 命中 0，**我自造反证**（仓外 overlay 摘那一行）B rc=1 31.22s／C rc=1 31.17s（`unbound-scope`＋`after=31s`，乙 `was_open=false`），交叉复跑他们的 `overlay-mut-noopen.json` B rc=1 31.21s 同红句；③Case C＝前置断言为 `t.Fatalf`（`tier=global` 清单行必须出现），修前 31.17s 红／修后 0.14s 绿；④去门控成立：strict runner 现跑 `PASS=3 FAIL=0 SKIP=0` rc=0，门控复造形 `SKIP is not a pass` rc=1，`scripts/portable-tests.sh` 零 315 名；⑤197 自批 `-race` 红＝当前树与 overlay 无本修树**同族同数**（各 7 warning；FAIL 13.99s／13.95s）⇒ 既有夹具红非本修引入（台账 `A865` 段有账）；⑥门禁现跑 build rc=0／`d22scan` clean／长度门 `VERDICT GREEN` 59=59／`Test222 -race` rc=0。件＝`docs/evidence/s1/ticket315-315-v1-accept-v1.md`＋`.scratch/wisp/probes/315/v1/`（锚＋logs 16 枚）。⛔ 未翻框、⛔ 未动 `Status:`/票名、⛔ 未 push。**收口建议：AC#0＋AC#1 双成立 ⇒ 票 315 可收**（翻框归编排者）；遗留仅"并发兄弟形未单独量"（写腿已登记候选，机制上不阻塞）。
