# 普查件 267-a1 — 票 267 的 AC#0 四问（只读腿，⛔ 零 `go` 命令）

> 腿：`267-a1`（只读普查）｜对象票：`.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish.md`
> ⚠ **本件由编排者代落盘**（09:0x）：派单选错了腿型——`Explore` 会话**没有 Write/Edit 工具**，腿按派单要求"落 `census.md`"物理上做不到，于是它把全文直接交回。⇒ 正文逐字是腿的产出，编排者只加这一句抬头与本节标注。
> ★ **编排者复尺（同发 09:0x 现跑，七条全对）**：①`grep -rn "WarningLead" cmd tools internal --include=*.go | grep -v internal/agent/approval | grep -v _test.go`＝**0 命中**；②`internal/config/schema.go` 内 `ConfirmTimeoutSec int` 的 `default:"300"`、`L1WindowSec int` 的 `default:"2"` 逐字在盘；③`queue.go:84-90` 的兜底**逐字**是 `if timeout <= 0 { DefaultApprovalTimeout }` 与 `if warnBefore <= 0 || warnBefore >= timeout { DefaultApprovalWarning }`；④`gate.go:143-151` 的 L1 三档 switch（`<=0→DefaultL1Window`／`<Min→Min`／`>Max→Max`）逐字在盘；⑤`[hotkey]`＋`[app]` 两节里数值键枚数＝**0**；⑥`validate.go` 中断数值值域的只有 `ball.size`/`ball.opacity_idle`/`cost.alert_threshold` **三枚**；⑦票面 §现量 3 那句 `resident_approval_windows.go:310` **行号确实已漂**——`:308-312` 现在落在注释块里，那枚算式真身在 **`:460`**（`run.go:616` 未漂）。⇒ 腿的第 4 节里"票面行号已漂"这一处**顶回我票面成立**，更正另记台账，原句不抹。

---

## §0 起手锚

`2026-10-05T08:59:56+0800`｜HEAD `21bec8a1`（`dev`）｜工作树脏（他人 `design/**` 删除、`.scratch/**` 未跟踪一片）。下面每节行号只对写下这一刻的树负责。

## §1 AC#0① WarningLead 的来源与恒定性 → **生产恒为 30s**

- 声明：`internal/agent/approval/gate.go:33-35`（`Options.WarningLead time.Duration`）；唯一进路 `gate.go:158` → `NewQueue(o.ApprovalTimeout, o.WarningLead, ...)`；落地 `queue.go:84,88-90,98`；读面 `queue.go:129`。
- 尺：`grep -rn "WarningLead" cmd tools internal --include=*.go | grep -v internal/agent/approval | grep -v _test.go` ⇒ **0 命中**（0 命中＝好消息）。生产里**没有任何一枚生产者传过它**，两枚装配点 `run.go:612-619`、`resident_approval_windows.go:368-375` 都只交 UI/Channels/Window/ApprovalTimeout/Logf/Grants ⇒ `Options.WarningLead` 恒零值 ⇒ `queue.go:88` 兜成 `DefaultApprovalWarning = 30s`（`queue.go:109`）。
- 唯一 `NewQueue` 调用者是 `gate.go:158`（其余命中全在 `internal/agent/approval/*_test.go`：100ms/25s/30s/300s 都是包内仪器）。
- ⚠ 关键读码：`queue.go:88` 的 `warnBefore >= timeout` 那一支**救不了本票形状**——判的是入参（0），命中后重新赋 30s；于是 `timeout=10s` 时 `q.timeout=10s, q.warn=30s` ⇒ `gate.go:528` `lead=-20s` ⇒ `warn` 保持 nil。C18 提示消失**不是**因为有人传了大 lead，而是因为兜底常量本身就没跟着 timeout 缩。

## §2 AC#0② 裸键名册

**题面三节逐枚清账（尺：`grep -n "int \`toml\|int64 \`toml\|float64 \`toml" internal/config/schema.go` ⇒ 全文件 40 枚，按节切分）**：`[hotkey]`（`schema.go:179-184`）与 `[app]`（`:137-150`）**一枚数值键都没有**（全 string/bool）⇒ 这两节裸键 **0 枚**。`[risk]` 数值键共 **2 枚**，即票面那两枚：

| 键 | schema 行 | 消费点 | 钳位 |
|---|---|---|---|
| `risk.confirm_timeout_sec` | `:450` default 300 | `run.go:616`、`resident_approval_windows.go:460` | 仅下界 `queue.go:85-87`（`<=0→300s`）；**无上界** |
| `risk.l1_window_sec` | `:453` default 2 | `run.go:615`、`resident_approval_windows.go:459` | `gate.go:143-151` 双向钳 `[2s,3s]`（带常量 `queue.go:115-122`） |

⇒ **AC#0② 严格答题数＝1 枚**（`confirm_timeout_sec`）。两枚键的全部生产读者就那 4 行，别无他处。

**附录（同形材料，供 §3/§5 用；`[risk]/[hotkey]/[app]` 之外，不扩题、只摆料）**——上界缺失的 6 形／下界有守卫的 8 形：

| 键 | schema 行 | 消费点 | 钳位在哪一层 |
|---|---|---|---|
| `llm.timeout_ms` | `:415` 60000 | `run.go:1085-1086` 裸乘进 `context.WithTimeout` | **全无**（0/负＝即刻过期；无上界） |
| `llm.retry.max` | `:306` 3 | `run.go:1227`→`internal/llm/retry.go:79-83` | 消费侧填 `==0→上限默认 / <0→0`；**无上界** |
| `llm.retry.backoff_ms` | `:307` | `run.go:1226`→`retry.go:84-85` | 下界 `<=0→1s`；无上界 |
| `agent.max_rounds` | `:431` | `internal/agent/guard.go:101-103`（`loop.go:223-224` 同形重复） | 下界；无上界 |
| `agent.token_budget` | `:433` | `guard.go:104-106` | 下界；无上界 |
| `agent.per_tool_timeout_ms` | `:436` | `run.go:763,1014`→`internal/tools/bridge.go:172-174` | 下界 `<=0→DefaultToolTimeout`；无上界 |
| `agent.loop_guard.repeat_thresholds` | `:425` | `guard.go:90-100` | 逐元素滤 `v>0`＋排序＋空则默认；无上界 |
| `llm.providers.*.rpm`/`tpm` | `:397-398` | `internal/llm/ratelimit.go:48-51` | `<=0→nil 桶（不限速）`；无上界 |
| `llm...context_window` | `:354` | `internal/agent/budgets.go:86-87` | **双向**：`<=0 || > ReferenceContextWindow → 128000` |
| `panel.width`/`height` | `:532,534` | `panel_resident_windows.go:207`→`panel_host_windows.go:245-258` | 下界 `>0` 才用；无上界（直塞 `uint`） |
| `panel.scale` | `:538` | 无生产读者〔仅读码〕 | — |
| `observe.roll.size_mb`/`days` | `:595-596` | 装配点未喂（`grep "RollSizeMB:"` 非测试 **0 命中**）；若喂则 `internal/observe/logging.go:182-187` 下界 | — |
| `[session]` 三枚／`privacy.retention_days`／`observe.slo_sample_interval_sec`／`llm...quota_*_micro`／`tts.speed`／`wake_word.thresholds` | `:188-190,513,606,367-368,218,201` | **零生产消费者**（命中只落在 `schema.go`/`manager.go` 层级拷贝/`*_test.go`）⇒ 属票 255「没读者」的地界，不是本票「有读者不问范围」 |

## §3 AC#0③ L1 钳位落在哪一层、可否复用

钳**不在 `internal/config`**、**不在 `cmd/wisp`**，而在**消费包自己的构造入口 `internal/agent/approval/gate.go:143-151`**（`approval.New` 的 `switch`，常量 `queue.go:115-122`）。`256-v1` 报的 99→3s／1→2s／0,-1→3s 与此实现**逐支对得上**（`:145` `<=0→DefaultL1Window(3s)`、`:147→Min(2s)`、`:149→Max(3s)`）——注意 `0/-1→3s` 走的是 `:145` 而非带内最大，读数相同、路径不同，别写成一回事。
超时不同路的原因：它绕开 `New` 的 switch、进 `NewQueue`（`queue.go:84-102`），那里只有 `timeout<=0` 一枚下界、没有上界分支，且 `warn` 的兜底不看 `timeout` 的带。
可复用性（只摆料）：同一层已有两枚同类件——`gate.go:143-151`（同函数）与 `internal/agent/budgets.go:86-87`（同形双向）；`guard.go:107` 是第三枚（`<=0 || > Max → 默认`）。⚠ 复用代价具名：这一层钳＝「配了不生效」，正是票 255 立案的形状；且 `resident_approval_windows.go:376-384` 的回执是**读对象现值**（`ra.gate.Queue().Timeout()`），钳在 `NewQueue` 里回执会自动跟着变，但 `unwired.go:121-122` 的指认释文与 `riskProvenance*` 三态（`:430-439,459-461`）会开始说谎。

## §4 AC#0④ 既有值域尺名册（自己拉的，未照抄票面）

`internal/config/validate.go`（242 行）里断数值值域的**只有 3 枚**：`ball.size` `[44,72]`（`:58-61`）、`ball.opacity_idle` `[0,1]`（`:62-65`）、`cost.alert_threshold` `[0,1]`（`:122-125`）。`validateRisk`（`:109-114`）确实只断 `permission_mode` 枚举——票面 §现量 2 的读数**方向成立**，但它具名的 `resident_approval_windows.go:310` **行号已漂**（现树该算式在 `:459-460`，`:310` 落在注释块里），`run.go:616`/`gate.go:528`/`gate.go:562` 未漂。
测试尺：`internal/config/validate_test.go:155 TestValidateBallSizeRange` 是**唯一一枚**直接断配置数值域的既有尺（尺 `grep -rln "out of range" cmd internal tools --include=*_test.go` ⇒ 3 个文件，另两枚 `leg_sink_nail_131_windows_test.go:214`、`audio/resample_test.go:153` 与配置值域无关）。
另一枚**间接**断值域的尺，别漏：`cmd/wisp/resident_approval_risk_256_windows_test.go:263-266` 现跑断 `Window()` 必落 `[MinL1Window,MaxL1Window]`——它断的是**消费侧钳**而不是 schema 带；同文件 `:251-254` 只断 timeout 等于种子值（45s/90s 都在带内语义上＝没界）。
`tools/d22scan` ban 名册：`allowlist.txt` 首列去重＝`bare-goroutine`／`mirror-hash`／`pathresolver-bypass` 三枚，全是代码形状，**与配置值域无关**。`unwired.go:121-122` 是「有没有读者」的尺，不是「值域」的尺（并顺带确认票面 §排程那句过期指认成立：只写了 `run.go`，未写常驻腿）。

## §5 三形各自要动的文件面（现读复认，⛔ 不裁不推荐）

- **ⓐ 值域进门**：`internal/config/validate.go`（`:109-114` 扩 `validateRisk`，注册进 `:24-45` 的链）＋`validate_test.go` 新钉；连带释文面 `internal/config/unwired.go:121-122`；`internal/config/schema.go` 注释（`:449-453`）。⚠ 作用面：常驻腿走的是 `config.LoadFile`（`resident_approval_windows.go:452`），ⓐ 的红会同时打在常驻腿与 `wisp run` 上——两枚消费者都在。`internal/agent/approval` 可零改动。
- **ⓑ 消费侧钳位**：`internal/agent/approval/queue.go:84-102`（新增带常量在 `:104-123` 一块）或 `gate.go:143-151` 同层；若钳在 queue，则 `gate.go:528` 的 `lead>0` 判据与 `queue.go:88` 的 warn 兜底要一起重读；连带 `cmd/wisp/resident_approval_windows.go:376-384` 回执与 `unwired.go:121-122`、`cmd/wisp/config_readers_255.go` 的 claim 文本。既有反向控制钉 `resident_approval_risk_256_windows_test.go:182-269` 的 `wantTime` 期望值会被这一形改写。
- **ⓒ 只补可见性**：`gate.go:526-535`（`warn==nil` 那一支）＋审计落点（`run.go:617` 的 `rt.auditf` / `resident_approval_windows.go:479-483`）；⛔ 治不了 27h46m 那一支（上界仍无人管），且对 27h 一侧需要第二枚落点才响得了。
- 三形共同面：`docs/PLAN.md` C18 属冻结面（AC#3 禁改、`Q-77` 待拍），任何一形的文案若引用「300s/30s」都只能引为默认值而非契约常量。

## §6 判不动／量不到（具名＋归口）

1. 种子读数（99999⇒27h46m39s、10⇒提示消失、45s 仍在）**未复跑**——禁 go 命令，归 `256-v1`。我只能证明机制式子在树里逐字成立（§1、§3）。
2. `queue.go:88` 在 `timeout` 被配到 `< 30s` 时的净效果（warn 仍 30s 而非缩短）是**读码结论**，没有跑出来的 `warn` 值佐证。〔仅读码〕
3. `observe.roll.*`、`panel.scale`、`[session].*`、`privacy.retention_days`、`slo_sample_interval_sec`、`quota_*_micro`、`tts.speed`、`wake_word.thresholds` 的「零消费者」判定基于 `cmd internal tools` 的具名 grep 零命中；`scripts/`、CI 采样侧是否偷读 `slo_sample_interval_sec` 未穷尽（那是 `ci-phantom-1` 地界）。归口票 255。
4. 三形作用面**一枚都没跑**，票面 §要建什么 的〔待验〕状态不因本腿而改变。
5. `internal/agent/approval/fakes_test.go:261-269` 那种逐字段拷贝的仪器会不会吃掉新字段（`256-a2` P9 报过）——本腿没跑包内仪器，未复认。

## §7 交件判语

零 go 命令。跑过的命令名册（全部）：`date`；`git rev-parse --short HEAD`／`--abbrev-ref HEAD`；`git status --porcelain`；`ls`；`grep -rn/-n/-c/-o/-l/-A/-i`（含 `cut`、`sort -u`）；`wc -l`；`sed -n`（只读打印）；`head`/`tail`；管道与 `for` 循环仅串上述只读命令。**无 `go`/`go test`/`go build`/`go vet`、无任何写命令、无 `git add/commit`（本腿无写面，交件为本文）。**
票面、`docs/PLAN.md`、`docs/specs/**`、任何 AC 框、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`：**零改动**（`git status` 起手即脏、收尾同一集合，均非我所为）。
搜索根只用 `cmd internal tools docs scripts .scratch`，从未以 `.` 为根；`frontend/**`、`design/**` **一枚都没读**（起手 `git status` 输出里出现 `design/**` 行，是 porcelain 列表自带，非读取）。
