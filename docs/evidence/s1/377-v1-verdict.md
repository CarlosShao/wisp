# 377-v1 终裁：〔丁〕的**记账面够格**，`AC#2` **能翻** —— 附 4 处须更正的件面/注释读数（非实现者裁腿）

- 腿 `377-v1`（**非实现者**）。被裁对象 `377-r1`＝`5cac443a5`。`git show --numstat 5cac443a5` 现量 **6 文件 `+617/−17`**＝`bridge.go 37/9`·`approval_always.go 26/6`·`config_reload.go 5/2`·`booked_mode_377_test.go 325/0`·`permission_mode_effective_377_test.go 102/0`·`377-r1-impl.md 122/0`。
- ★**读数基线**＝`5cac443a5` 的**代码面**。读数期间 HEAD 前进到 `32afe01fb`（编排者 A1240 的台账提交，`git diff --stat 5cac443a5..HEAD`＝**只有** `docs/reports/pending-and-issues.md +64`）⇒ 代码面**逐字同**，本笔读数不因 HEAD 前进而失效；台账那件本笔 ⛔ 碰。
- 现树＝HEAD：写面 `git status --porcelain -- docs/evidence/s1 internal/tools cmd/wisp internal/perm` **空**（脏 `.scratch/**` 与 `$null`／`-`／`-dummy` 是别人的，**未碰**）。
- 材料现读：票面 **106 行**（含「`AC#1` 裁形（改后）」节）、普查 **160 行**、r1 件 **122 行** —— 三件行数与派单一致。
- 尺的射程（落笔前五问）：①对象＝**`tools: call` 行的生产 formatter 输出**（`bridge.go:1201` 那条 `b.log`），⛔ 不是 `Decision.Mode` 字段；②形状＝**整行逐字**；③数的是**代码**（`git grep` 分层，注释单列）；④每个数都与总数对过；⑤谁算＝`go1.27.1 windows/amd64` 现跑。
- 读数环境：`git archive HEAD` **全量**导出（17716 文件，含 `scripts/`／`frontend/`／`docs/`／`.scratch/`）到仓外 `D:\tmp\wisp\20261024\377v1\v1leg377export\headtree377`；改前树＝`5eddc6223` 同法。**导出树要满足 252 的 8.3 前提：树根**和**父目录都要有别名** —— 我先放 `…\377v1\head`（全 8.3 合规名）得 **8 红**（全 `paths_shortname_252*`、逐条自报「前置条件缺失」）；只把树根改成长名仍剩 **1 红**（`_r1_test.go:166` 要的是**父目录**别名）；两级都换到 `…\v1leg377export\headtree377` 才**全绿** ⇒ **那 8＋1 枚是尺缺件、⛔ HEAD 缺陷**（r1 报 8＋1；`../../scripts` 那枚因我导出全量而本来不红）。
- 探针**全在仓外**（`probes/zzz377v1_{tools,runseam}_test.go`，`-overlay` 注入）；变异 M1…M5＋`m1`/`m2` 亦全在仓外。**工作树零字节突变窗口。**

## ① 十一项逐项判（每项另起尺）

| # | 项 | 判 | 我的读数（现跑） |
|---|---|---|---|
| 1 | 票面前提新鲜度 | **成立**；行号两处漂已记 | 票面点名的 `approval_always.go:151-156` 在 r1 基线 `72da06315` **逐字**＝`permissionMode()`、`:155`＝`rt.mgr.Config().PermissionMode()`；现树该函数已移到 `:168-176`（读 `rt.modes`）。`config_reload.go:282`（普查 #2）现树 `:285`。`store.go:158`／`:254-255` 现树**逐字未动**。`run.go:864` `approval.UnsilencedMode(…)` 逐字在。 |
| 2 | `AC#2` 那一行 | ★**能读出 `mode=ask_every_step`** | 见 ③（我自己装的形状，逐字） |
| 3 | 两枚"错答案"用例 | **真的是逐名断言**（无"只看能读出正确值"） | 见 ④：两个错值**各有可落地的错实现**（M3/M4 撞 `auto_approve`、M5 撞 `ask_high_risk`）|
| 4 | 正控（真 `wisp run` seam） | **没被改成"一律最严"** | 见 ⑤ P5：`mode=auto_approve` |
| 5 | "一次读、全链一致" | **删第二次读对**；"两次读会不一致"**成立** | 见 ⑥（P3＋窗口内有 FS 活）|
| 6 | 追加位置那条理由 | **结论对、证据数错**：连续子串读者**只有 1 处** | 见 ⑦ |
| 7 | `m1`/`m2` 有没有落地 | **零落地**＋两结论我自己复现 | 见 ⑧ |
| 8 | 反向判据 | **我自造 M1／M2 都红**；另加 M3/M4/M5 与 M0 | 见 ⑨ |
| 9 | 四数复跑 | 四包分开、⛔ 相减 | 见 ⑩ |
| 10 | 边界逐项 | **6 文件逐枚对上**；五处**零 diff** | 见 ⑪ |
| 11 | "宽一格"该不该在本票做 | **不该在本票做**；但 r1 的归类要收窄 | 见 ⑫ |

## ② `AC#2` 能翻／⛔ 翻 ＋ 缺什么

**能翻。** 票面判据（`:101-102`）是「造一发 `panel-inbound` 形状 ⇒ `tools: call` 那一行必须能读出 `mode=ask_every_step`」，且⛔ 只测"字段存在"、⛔ 只测 `Decision.Mode`。我对这两条各设一尺：
- **⛔ 字段存在**：`M1`（摘掉新加那一跳）把行尾还原成"没有 mode"⇒ 三枚**逐名红**；正确值断言用的是 `mode=ask_every_step` 整串，不是 `mode=` 存在。
- **⛔ 只测 `Decision.Mode`**：断言对象是 `b.log` 的**整行**（③）。`git grep -n "Decision\.Mode"` 在两枚新测试里 **4 处命中、全是注释/错误文案**（`booked_mode:6,7`、`perm_mode_eff:5`、`:71` 的 `t.Errorf` 措辞），**没有一处读／断言该字段**。
**缺什么（不挡翻，但要写明）**：今天**没有任何生产装配同时有"桥"与"nil 腿"**（`panel-inbound` 无桥，`wisp run` 恒有腿 `confirmModeSwitch`）⇒ 票面点名的形状是**合成形状**；它在本仓的**生产实例**是 375 那发 —— P6（⑤/⑫）：`wisp run` 无人可答、`config=auto_approve`，行上已是 `mode=ask_every_step`。⇒ 本票的价值不靠合成形状独撑。

## ③ ★★★`AC#2` 那一行（我自己的，探针 P1，逐字）

装配＝`config.toml(permission_mode=auto_approve)` ＋ 真 `config.NewManager` ＋ 真 `perm.New{Confirm:nil}` ＋ `Set(ask_high_risk)` ＋ 真 `BuiltinFSEntries`/`C26`/`Bridge`（⛔ 复用 r1 的 fixture，我自己写）。前提三格先钉死：`effective=ask_every_step`／`configured=auto_approve`；Set 后 `configured=ask_high_risk`、`effective` 仍 `ask_every_step`；`windows=1 approvals=0 executed=true`。

```
tools: call kind=success task=task-377 corr=corr-377 tool=fs.write risk=L1 decision=allow outcome=success rules_hit=[R1] in_allowlist_scope=true grant_id=0 refusal=none reason="R1: 工具声明为下界（L1）" mode=ask_every_step
```

`has_auto_approve=false has_ask_high_risk=false` ⇒ **既不是**文件原值、**也不是**切档后文件值。与 r1 §③ 报的逐字**相同**（含 `reason=` 那段中文、`grant_id=0`、`rules_hit=[R1]`）。

## ④ ★★★那两枚"错答案"用例的核

用例（`booked_mode_377_test.go:239-245`）先把 `mode=ask_every_step` 断真，再对 `{auto_approve, ask_high_risk}` 逐枚断"若出现 `mode=<错值>` 就报**错的这个名字**"。**核过，不是只看正确值。** 我为这两枚各造一个**可落地的错实现**（都只改接线，断言一字不改）：

| 变异 | 错实现 | 行上印的 | 用例报的（逐字摘） |
|---|---|---|---|
| `M3` | `Execute` 里硬编码最松档（stale/cached 面） | `mode=auto_approve` | `the booked line reports config.toml's "auto_approve", not the effective档` |
| `M4` | 桥的 `ModeSource` 接**装配时快照** `mgr.Config()` | `mode=auto_approve` | 同上（**auto_approve** 被逐名点出）|
| `M5` | 桥的 `ModeSource` 接**现读**配置面（每次读文件） | `mode=ask_high_risk` | `the booked line reports config.toml's "ask_high_risk", not the effective档` |

★两枚错值的**触发条件不同**、都真实：`ask_high_risk` 由**现读**配置面撞出（M5，这是**最自然**的错写法）；`auto_approve` 由**装配时快照**撞出（M4 —— `mgr.Config()` 返回的指针在 `Set` 后仍是旧值，P4 实测同一发印 `mode=auto_approve`）。⇒ r1 那句"把两个错误答案逐名断言为红"**成立**（不是"能读出正确值"就算）。

## ⑤ ★★正控（我自己的，探针 P5）＋ 真 seam 的第二种差面（P6）

P5＝真 `wisp run` 装配（`runTextTask`＋回复管、`config=auto_approve`、一发 L1 写、exit=0）：

```
[audit] tools: call kind=success task=t377v1-attended corr=t377v1-attended-corr tool=fs.write risk=L1 decision=allow outcome=success rules_hit=[R1] in_allowlist_scope=true grant_id=0 refusal=none reason="R1: 工具声明为下界（L1）" mode=auto_approve
```

⇒ **没被改成"一律 ask_every_step"**；卡面 `rt.permissionMode()=auto_approve` 与行一致（该 seam 两面恒等）。
P6＝同配置、**无人可答**的 console 跑：`MODE-UNANSWERABLE=true`／`L1-UNANSWERABLE=true`／行上 `mode=ask_every_step`（`kind=refused`）／**卡面仍是 `auto_approve`** ⇒ ⑫ 那条边界**是真的**，且它在**生产 seam** 上就成立，不只在合成形状里。

## ⑥ ★★"一次读、全链一致"的判语（含"两次读真的可能不一致吗"）

- **结构核**：`mode := b.permissionMode()`（`:286`）→ 早退分支 `Decision{…, Mode: mode}`（`:290`）→ `dec.Mode = mode`（`:297`，C3 后）→ `mode.Screen(verdict)`（`:327`）与 `dec.Mode = mode`（`:328`）→ `book()` 印 `dec.Mode`（`:1206`）。**同一个变量**，⛔ 第二次采样 ⇒ "行上的值就是筛选用过的值"**逐字成立**。
- **"两次读真的可能不一致吗"＝会**：①P3 实测读面是**活的**（`first=auto_approve second=ask_every_step differ=true`，同一 `*perm.Store`，中间一次 `Set`）⇒ 不是"值被冻结"；②旧代码里两个采样点之间隔着 `Lookup`＋`checkCaps`＋`decodeArgs`＋`Assess`＋`factsFor`＋`readBlacklist`（**含 FS 活**），面板侧的 `Set` 是本进程另一条 goroutine 就能落的。⇒ **删第二次读是对的**（若 `book()` 再采一次，行可能与 `MODE-SILENCE` 说的不是一件事，而 `book()` 发生在工具跑完之后 —— r1 注释点的正是这个窗口）。
- **一处要写明（r1 未提）**：提读把**筛选用的采样点**从"Assess 之后"挪到"Assess 之前" ⇒ 若 `Set` 恰落在 `Assess` 期间，这一发会晚一格生效（收紧／放宽都晚一格）。这是"一次读"的代价，方向主流（准入时取一次）、量极微（单发），但**应当具名**，⛔ 当零成本。

## ⑦ ★★追加位置那条理由的核

- **真读者只有 1 处**：`ticket363_test.go:258` 传一整串 `"tool=fs.read risk=L0 decision=allow outcome=success"` 给 `t363LineWith`（内部 `strings.Contains(l, needle)`）⇒ **是连续子串读者**，中段插 token **会破**。
- **r1 点的第二处不是**：`fs_edit_ac5_gate_r4_test.go:206` 把 `"risk=L2"`／`"R2"`／`"in_allowlist_scope=false"` 交给 `r4AuditLine`（`:62-83`），那是**逐 token `Contains`**，插在中间**不破**。⇒ 件里"至少两处"应改成"一处"。
- **反方向我也查了**：全仓 `*_test.go` 里**没有**对 `tools: call` 行做 `HasSuffix`／`$` 锚定的读者（`git grep HasSuffix` 唯一命中是 `refusal_20r2_test.go:378` 的数文件名，`git grep 'tools: call.*\$'` 0 命中）⇒ **追加在行尾两头都安全**。
- ⇒ **结论（行尾、`reason=` 之后）对**；**理由的计数错**，且 `bridge.go:1195-1196` 注释把它归给 `ticket 173's shape filter` 属**误归**（`audit_allowlist_scope_173_test.go` 用的是单 token `Contains`；该 shape filter 在 `ticket363_test.go:258`）。

## ⑧ ★`m1`/`m2` 有没有落地

**零落地。** `git diff --numstat 5cac443a5~1 5cac443a5 -- internal/perm` **空**；`store.go:158`（`configured := s.configuredMode()`）与 `:254`（`from := s.configuredMode()`）**逐字未动**；`RequiresCard`（`direction.go:23`）零变化 ⇒ `R20/M4` 未动。
"只在核心结论面复跑"我给它**定量**：r1 复跑的是 `M2`（`Store.PermissionMode` 回配置档），**不是**普查的 `m1`／`m2` 原件（`:254`／`:158`）。我把原件也各造一发（仓外 overlay）：`m1` ⇒ `Set(ask_every_step)` 变 `result=no-change`、`config_after=auto_approve`（**严档请求被静默吞掉**，`from=ask_every_step to=ask_every_step`）；`m2` ⇒ `has_MODE-DOWNGRADE=false`（**整条消失**）。⇒ 普查那两条结论**今天我复现**，两处确实 ⛔ 动。

## ⑨ ★★反向判据（我的 M1／M2 ＋ 我加的 M3–M5 ＋ M0）

| 变异 | 改哪（全在仓外 `-overlay`） | rc | 红行 |
|---|---|---|---|
| **M1** | `book()` 摘掉 ` mode=%s` 与 `dec.Mode` 实参 | **1** | 3 枚逐名：`does not name the effective档 "ask_every_step"`；两发对照 `mode= is not "ask_every_step"`／`"auto_approve"`；`the pre-C3 refusal booked no effective档`。行尾**回到"没有 mode"**（逐字）|
| **M2** | `Store.PermissionMode` 回配置档 | **1** | `premise broken: effective档 = auto_approve`；`windows = 0, want 1`；行印 `mode=auto_approve`；`MODE-SILENCE present = true, want false` |
| **M3** | 硬编码最松档 | 1 | ④（`auto_approve` 逐名）|
| **M4** | 桥接装配时快照 | 1 | ④（`auto_approve` 逐名）|
| **M5** | 桥接现读配置面 | 1 | ④（`ask_high_risk` 逐名）|
| **M0**（改前树＋新测试） | 换回 `5eddc6223` | **1** | `internal/tools` 3 红（行尾无 `mode`）；`cmd/wisp` **1 红**（`the widening card would book auto_approve …`）＋**2 枚 PASS** |

★**"2 枚改前亦绿"＝刻意的正控，⛔ "没牙"**：M0 里 `…KeepsTheConfiguredModeWhenAChannelExists`（防"改成硬编码最严档"）与 `…WithoutAStoreFailsClosed`（防"无 store 时猜向静默"）都**必须**改前绿 —— 它们的职责是**反方向**；证的正是方向性（P5 独立复证）。**小缺口（建议，不挡翻）**：新 fail-closed 分支`rt.modes==nil 而 rt.mgr!=nil` 没有"改前红"用例（`bare := &agentRuntime{}` 两 nil，改前也绿），补 `&agentRuntime{mgr: mgr}` 一发即可钉住"删掉 `rt.mgr` 兜底"这一跳。

## ⑩ ★四数复跑（四包分开、⛔ 相减；改后＝`5cac443a5`，改前＝`5eddc6223`）

| 包 | 改后 rc／读数（`5cac443a5` 树） | 改前 rc／读数（`5eddc6223` 树；与 `72da06315` 代码面同） |
|---|---|---|
| `./internal/tools` | **0** `ok 20.167s` | **0** `ok 20.209s` |
| `./internal/agent/approval` | **0** `ok 0.661s` | **0** `ok 0.583s` |
| `./internal/perm`（未改、旁证） | **0** `ok 0.394s` | **0** `ok 0.465s` |
| `./cmd/wisp -run '<9 名roster>'`（带 `PATH=…/sherpa-onnx`、只定向） | **0** `ok 184.056s` | **0** `ok 184.564s` |

roster 与 r1 同：`TestAlwaysBranch|TestAC1AlwaysBranch|TestTicket223|TestTicket114r2|TestTicket101|Test375|Test376R1|Test377|TestReplyListenerReject`。⛔ 整包（会死在 `MARK299`）。
★**四包四对全 0，⛔ 相减**：改后并非一律更快（`internal/tools` 改后 20.167 ＜ 改前 20.209，`perm` 改后 0.394 ＜ 改前 0.465，`approval` 改后 0.661 ＞ 改前 0.583）—— 两侧同量级、**读数不可用于"更快／更慢"的结论**。r1 报的 `21.326/0.727/0.571/207.246` 与我的同量级、rc 相同。

## ⑪ ★边界逐项（6 文件逐枚 ＋ 五处零 diff）

- 提交面**恰** 6 文件（`git show --name-only`），与件 §⑩⑤ 列的点名**逐枚相同**，无第 7 件；`git diff --numstat` 三枚代码文件＝`37/9`／`26/6`／`5/2`（与件同）。
- `internal/perm/**`、`cmd/wisp/run.go`（含 `consoleApprovalUI.Prompt`）、`internal/agent/approval/{reach.go,gate.go}`、`internal/{winsec,memory,panel}/**`：`git diff --name-only 5cac443a5~1 5cac443a5 -- <这些>` **逐个空**。
- `R20/M4`：`perm.RequiresCard` 与 `Store.Set` 零字节（⑧）。

## ⑫ ★★"宽一格"那条边界：我判**不在本票做（对）**，但 r1 的归类要收窄

- **事实核**（P6，真 seam）：`wisp run` 无人可答时，**卡面** `rt.permissionMode()=auto_approve`，**桥面**（`approval.UnsilencedMode(rt.modes, answer, …)`，`run.go:864`）`=ask_every_step` ⇒ 差的是**两档**（`auto_approve`→`ask_every_step`），不是字面的"宽一格"。
- **归因**：r1 说"属 375 的射程"—— **半对**。**形状是 377 的**（同一件事两个读面不同源，正是本票立项句）；**缺的件是 375 的**（`answer` 谓词没进 `agentRuntime`）。收它＝把 `answer` 引进 runtime，或让卡也吃包装源 — 那是**结构改动**，且受益字段（`Decision.Mode`）**今天 0 读者**、`MODE-UNANSWERABLE` 已把信息落下（375）⇒ **不值得在〔丁〕这一格里做**。
- ⇒ **判：本票不做＝对**；但件里应写成"**同形缺陷的第二个面、缺 375 的件**"，并**登记为 375 的跟进项**（⛔ 写成"另一张票的形状所以本票不提"）。

## ⑬ 门禁 rc（**工作树**，现跑）

| 尺 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 无输出 |
| `go vet ./internal/tools ./internal/perm ./internal/agent/approval ./cmd/wisp` | **0** | 无输出（导出树同：0／0） |
| `go run . -root <repo>`（`tools/d22scan`，独立 module） | **0** | `clean - no D22 ban violations`；`bans #1-5 internal/=250`＋`cmd/=42` ⇒ **292 枚生产 `.go`**（与 r1 同数）；`#8` 违例 0 |
| 格式**甲尺** `gofumpt.exe -l <本笔 5 枚 .go>` | **0** | 输出 **0 字节** |
| 格式**乙尺** `gofmt -l <同 5 枚>` | **0** | 输出 **0 字节** |
| EOL：5 枚字节统计 | — | `CR(0x0D)=0`、纯 LF（`bridge.go 61859B/1443LF`、`booked_mode 12392/325`、`approval_always 11067/215`、`config_reload 24875/444`、`perm_mode_eff 4635/102`）|
| ⛔ 反向对照（**尺的射程**） | 2 | `gofumpt -l .` **54 行、rc=2**：`^\.scratch[\\/]` 命中 **51** 行（故意写坏的探针）＋**3 行是生产 `.go`**（`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`，**纯 CRLF**：CR=389/131/225），`gofmt -l` 同样列这 3 枚。⇒ ①⛔ 用 `gofumpt -l .` 判 HEAD 格式／EOL；②我的"两尺空"**只对具名对象（本笔 5 枚）成立**；③那 3 枚是**改前就有的**（末次提交 `aeea4e5b7` 2026-10-22／`d03d166fe`／`5413f46d4` 2026-10-04，**都不在 `5cac443a5` 的 6 文件里**）⇒ **⛔ 记成 377-r1 的缺陷**；④它与 `.gitattributes` 的 `*.go text eol=lf` 冲突，属 **HEAD 级既有缺陷**（建议另立票）|

## ⑬b 还原核（**最后一发变异 M5 `18:12` 之后**量；本笔变异全在仓外，工作树零突变）

- 五枚 `.go` 的 sha256（`sha256sum`，64 位；前 12 位现量）＝`bridge.go 952380bc2c10`、`approval_always.go e55cad2739c5`、`config_reload.go c22708f87f9b`、`booked_mode_377_test.go 48ea5fccc2bc`、`permission_mode_effective_377_test.go fcc69adcbb32` ⇒ 与 r1 件 §⑧ 报的**逐枚同名同前 8 位**。
- `git diff --numstat 5cac443a5~1 5cac443a5`＝`26/6`、`5/2`、`37/9`（＋三枚 `新增`）⇒ 与件同。
- `git status --porcelain -- docs/evidence/s1 internal/tools cmd/wisp internal/perm internal/agent/approval`＝**只有** `?? docs/evidence/s1/377-v1-verdict.md`（本笔写面）。

## ⑭ 提交号 ＋ 逐笔越界

- 本笔唯一写面：`docs/evidence/s1/377-v1-verdict.md`（＋仓外 `D:\tmp\wisp\20261024\377v1\**`）。提交号＝本件所在 commit（主题逐字以 `票 377 终裁腿 377-v1` 开头）。⛔ 翻框／⛔ 改票名／`Status:`／`-done`／台账／⛔ push。
- **逐枚 `.go`：0**（`internal/**`、`cmd/**`、`scripts/**`、`.github/**` 零字节；探针 `.go` 与变异件全在仓外）。⛔ 碰 `.scratch/**` 与怪名残留。

## ★须更正的 4 处（不影响"能翻"，但件面要改）

1. `377-r1-impl.md` §⑤ 的"至少两处连续子串读者"⇒ **一处**（`fs_edit_ac5_gate_r4_test.go:206` 是逐 token）。结论（行尾追加）不变。
2. `bridge.go:1195-1196` 注释把 shape filter 归给 `ticket 173` ⇒ 实为 `ticket363_test.go:258`（`ticket 173` 自己的测试用单 token）。
3. `config_reload.go:57` 注释仍写「call rt.permissionMode(), **which reads Manager.Config()**」⇒ 本笔把该函数改读 `perm.Store`，这句话**已成假陈述**；同文件 `:278-284` 改了、`:57` 漏了（零行为、本笔射程内）。
4. §⑥ 那条"采样点前移"未在件／注释里具名。
