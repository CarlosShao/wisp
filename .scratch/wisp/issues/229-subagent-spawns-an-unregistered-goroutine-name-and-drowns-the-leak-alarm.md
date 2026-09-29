# 229 — 子代理收尾时注册了一枚**不在 D38 名册里的协程名**（`subagent-finish-*`）：每 spawn 一次就打一行"疑似泄漏"警告，还被算进 SLO 的协程总数；而"名字不在名册"本该是发现真泄漏时唯一的信号

- Status: **待派，小活**（队尾，排在票 228 之后或与票 227 同批）。⚠ **两支修法里只有一支不用问人**：把 spawn 名换成名册内的形状＝可直接做；往名册里加 `subagent-finish-`＝**动 D38b 那张表＝人工批准**，须先落 `A##` 并摆给 owner。⛔ 第三支"把那行警告调轻／删掉"＝**削弱仪器**，本票不许。
- 来源：只读普查腿 `228-a1`（`.scratch/wisp/probes/228/a1/cost-table.md`，69,717 字节／354 行／10 枚提交）第 **§8 末节**（`:336` 那一条）＋**编排者 09-29 16:0x 自己复跑三把尺复认**（台账 `A439`）。
- ⚠ **归属要说清**：这是**票 197 子代理落地那批产码**留下的，也就是**本编队自己写的码**——不是别人地界，不许"顺手当历史遗留放着"。
- ⚠ **这不是安全漏洞，照三行读**：①**现象在哪**＝本机运行日志与 `wisp slo` 的协程数读数；②**有没有本机被入侵的证据**＝**没有**；③**最坏后果是什么形状**＝**将来真有一枚跑飞的协程时，那行唯一的报警混在常规噪声里没人看得见**（是"报警器被自己的日常响声淹掉"，不是"有东西被绕过"）。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| spawn 点 | `internal/tools/subagent_197.go:372` `observe.Default.Spawn("subagent-finish-"+bg.ID, "tools", watchRoot, func(c context.Context) {…})`；`:371` 那枚 root 用**同一个名字** | 现读（编排者复跑） |
| D38 名册一共三张表 | `internal/observe/goroutine.go`：常驻表含 `ui-sta`／`audio-capture`／`hotkey-listener`／`db-writer`／`watchdog`／`log-flusher`；按需表＝`kws-infer`；每任务前缀表＝`agent-task-`／`tool-exec-`／`approval-waiter`；另有一族（`asr-infer`／`tts-infer`／`panel-host`／`retention-job`／`model-downloader`／`memory-extract`／`disposal-worker`） | 现读（`sed -n '30,120p'`） |
| ⛔ **`subagent-finish-` 逐枚全不匹配** | `grep -n "subagent" internal/observe/goroutine.go` ＝ **0 命中** | 编排者现跑 |
| 不匹配以后发生什么 | `ClassifyGoroutine` 归 `CategoryUnknown`（注释逐字 `NOT in the roster: leak symptom`）⇒ `:270-272` 逐字 `slog.Warn("goroutine outside the D38 roster (leak symptom)", "goroutine", name, "owner", owner)` ⇒ **每次 spawn 一行 WARN** | 现读 |
| 还被算进分母 | `RosterReport.Unknown` 的**生产消费者＝0**（只有那行 WARN）；但 SLO 那扇门数的是 registry **存活总数**（`thresholds.go:161-162` ← `sampler.go:365` `s.reg.Count()`）⇒ **名册外协程抬高分母**，拿"总数"去对 `PLAN.md:2838` 那句"常驻 6＋每任务 3"会得出"超了"的错误结论 | 〔腿报＋编排者复认方向；枚数请续腿自己现跑〕 |

## 后果

1. **唯一的泄漏信号被常规噪声淹掉**：`CategoryUnknown` 是"这枚协程不在冻结名册里＝跑飞了"的那一句报警。今天它每枚子代理收尾都响一次 ⇒ **真跑飞那天，没人能从日志里把它挑出来**。
2. **协程数读数被系统性抬高**：任何拿 SLO 协程总数对 D38 预算的人，都会把"我们自己漏登记的一枚名字"读成"超预算"。
3. **日志噪声**：按今天并发 8 的形状，一发任务派若干子代理就若干行 WARN。

## 判据（每格都要现跑读数；不许用"改了那行字符串"充当"这一类已清"）

- [ ] **AC#1 名册先现读一遍**：三张表逐枚抄出（带 `file:line`），并现跑 `grep -n "subagent" internal/observe/goroutine.go` 证"不匹配"。**不许引本票这张表**。
- [ ] **AC#2 两支修法选一，另一支具名判死**：**甲＝换 spawn 名**到名册内已有的形状（每任务前缀那一族最贴：收尾协程本来就是任务派生）；**乙＝往每任务前缀表加 `subagent-finish-`**——⚠ 那**是动 D38b 那张表**，须先落一条 `A##` 批准记录并摆给 owner，⛔ 不许 agent 自己加完再说"表里本来就该有"。
- [ ] **AC#3 正控：门不许被拆**。修法落地后**必须仍有一发判据**证明"真·不在名册的名字照样打 WARN／照样被归 `CategoryUnknown`"（种一枚假名字进去断它响）。⚠ 本票最容易犯的错就是**为了消噪把那行警告改轻**——那＝削弱仪器，验收直接判失败。
- [ ] **AC#4 全仓同类普查（这一格才是本票的主要价值）**：把**全仓每一枚 `observe.Default.Spawn(...)`／`NewRoot(...)` 的首参名字**逐枚对名册，交回"不在册的还有几枚、分别在哪枚 `file:line`"的**具名清单或'零枚'的读数**。⛔ 只修自己撞见的那一枚就宣称这类已清＝不合格。⚠ 顺带复核票 228 的 `228-a1` §8 里那句"口径只到 `cmd/wisp` 这一层"——**这一次要扫到全仓**。
- [ ] **AC#5 SLO 口径要说一句话**：名册外协程计入总数那一支（`sampler.go:365`）在修法落地后读数变化具名写清（改前／改后各一枚，安静期）。⛔ `thresholds.go`／golden 一字节不许动。
- [ ] **AC#6 整包终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1` 到终态＋**逐名比红名册**（历史在册 `internal/ball` 1＋`internal/panel` 4 属别人地界，不算本票新增、不许顺手修；`internal/risk TestResolvePerCallBudget` 是争用型假红，安静 `-count=3` 复量为准）。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字不动；`frontend/**`／`design/**` 零读零写零转述；不新增 C17 方法名；新增导出名要先落 `A##`；⛔ **不许用"把 WARN 降级成 DEBUG／删掉那行"**来交差；⚠ 与票 **228** 同撞 `cmd/wisp`＋常驻拓扑、与票 **197 续腿**同撞 `internal/tools/subagent_197.go` ⇒ **串行**，跑门前自证没有别的 `go test` 在飞。
