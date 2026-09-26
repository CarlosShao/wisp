# 票 145 —— 快照字段落地（r1，实现程 · 代码那半）

> **本件的射程**：票 145 的 **AC#2／AC#3／AC#5／AC#6**，按 `.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md`
> 解冻的那一小块（`internal/panel/composer.go`／`internal/panel/pump.go`／`cmd/wisp/panel_pump.go`）执行。
> AC#1 的对照表是另一枚只读程交的（`docs/evidence/s1/145-snapshot-field-census-r1.md`，锚 `88eab34`），本件在它上面**复量并更正**。
>
> **一句话结论（先说，免得被读成"做完了"）**：**落地集＝空。** 不是"没找到有源的字段"，是
> **有源的字段全部卡在四把锁上，而那四把锁的钥匙一枚都不在本程的放开面里**。
> 本程把这四把锁逐枚**量出来**（含一次仓外全量闭合实验：把字段真落下去、真填上活对象、真加断言，看它能不能绿），
> 并把**验证过的闭合补丁**交到 `.scratch/wisp/probes/145/patch/`（`git apply --check` 在 `e24ae28` 上干净）。
> **票面六格本程一格不自勾。**

---

## 0. 锚点 · 地界 · 零改动自证

### 0.1 锚点（进场现读，非抄派单）

```
$ git rev-parse --short HEAD        # 进场 09:4x
720cae6
$ git rev-parse --short HEAD        # 写本件时（10:08）
e24ae28
$ git branch --show-current
dev
```

锚点从 `720cae6` 漂到 `e24ae28` 是**共享树里别家在提交**，不是漂移。`720cae6..e24ae28` 现量 **17 枚**：
`f58c5136` `5d46f241` `bec17272` `4e169760` `eed7c998` `63e228fc` `9a8766e6` `10e3585c` `c0778696`
`6550dc42` `81673ddc` `e9ce4198` `dd2b277b` `6de3d1c5` `5429c0d4` `da8b6677` `e24ae28d`
（票 153 证据件／票 151 交件与收表／票 152 交件与探针／台账 A273–A279／三枚验收派单存档）。

**它们有没有碰本程那三枚文件——现量，不推理**：

```
$ for h in $(git log --format="%h" 720cae6..HEAD); do git show --name-only --format= $h \
      | grep -E "internal/panel/|cmd/wisp/panel_pump"; done
（无输出）
$ git status --porcelain -- internal/panel/composer.go internal/panel/pump.go cmd/wisp/panel_pump.go
（无输出 = 工作树与 HEAD 一致，且无人在写）
```

⇒ 那 17 枚落的是 `cmd/wisp/run.go`／`slo_windows.go`／`internal/tools/bridge.go`／`docs/**`／`.scratch/**`，
**与本程的放开面零重叠**（编排者在 `A273②` 也是这么划的；他那句"不碰 151 的 `run.go`"本件要更正，见 §1 P3）。

### 0.2 本程写了什么、没写什么

| | 读数 |
|---|---|
| 生产码改动 | **零枚文件的逻辑改动**。只有 `internal/panel/composer.go` 的 `Snapshot` 头注释加了一段**指针**（把本件量到的四把锁钉在载体上），**不改行为、不加字段**（见 §7 那枚 commit） |
| 测试码改动 | **零**。放开面里没有 `*_test.go`，本程不假装能写（⇒ AC#6 的锁，见 §3.4） |
| `frontend/**`／`design/**` | **一个字没读、没写**。唯一例外＝按普查 §4.4 自己要求的取证，留了一行 `git status --porcelain -- frontend/src/lib/panel.ts`（读数：**空**，即那把尺读的是与 HEAD 一致的版本，见 §5.3） |
| 单独保留的三枚 | `tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene_test.go` — **零字节**，且本程没为它们中任何一枚的变绿做任何事 |
| 契约轴 | `docs/PLAN.md`／`docs/specs/**`／`internal/risk/**`／`thresholds.go`／golden／`allowlist.txt`／`rules_gateway.go` — 本程的全部操作是 `grep`／`sed -n` **读** |
| 台账／工单 | `docs/reports/pending-and-issues.md` 与 `.scratch/wisp/issues/145-…` **零字节**（勾与账归编排者） |
| push | **零**。三枚 commit 全在本地 `dev` |
| 临时件 | 仓外副本 `D:\tmp\wisp145-e1`（**只建不删**，闭合实验的现场）＋仓内原始读数 `.scratch/wisp/probes/145/`（E1–E4、基线、名册、补丁） |

### 0.3 尺（本件只认这几把）

| 尺 | 出处 | 用来判什么 |
|---|---|---|
| 载体 | `internal/panel/composer.go:44-49` | 有没有那一枚键 |
| 装配根 | `cmd/wisp/run.go:421-427`（`panel.NewSnapshotPump(panel.PumpSources{…})`） | 那一枚键**谁填** |
| 事件汇 | `cmd/wisp/run.go:751-785`（`consoleSink.Publish`） | `run.*`／`failures` 的数据今天落在哪 |
| 双向对账 | `internal/panel/composer_test.go:48`／`approval_test.go:105` | Go 加了键而 `panel.ts` 没跟上 ⇒ **谁响**（本件实跑，见 §2） |
| 键集钉 | `internal/panel/pump_test.go:111-124`／`:270-276` | 快照的 JSON 键集＝契约本身，第五枚键进来就红 |
| 屏词汇 | `frontend/src/lib/panel-views.ts`（已提交版）＋全树 grep | AC#3 ⓐ 的值有没有名字 |
| 单调时钟 | `internal/observe/clock.go`＋`internal/agent/approval/gate.go:257`＋AGENTS §1.2 | `remainingMs` 这类计时字段能不能落 |

⚠ **方法学一条**（承普查 §7.2）：本仓所有 commit 的 author 同名，**"我碰了什么"只能按自己的 commit hash 集 `--name-only` 现量**；
`720cae6..HEAD` 的区间 diff 里 `cmd/wisp/run.go`／`internal/tools/bridge.go` 各 1 枚，**都不是本程的**。

---

## 1. 派单与普查给的前提，逐枚复量（不照抄）

| # | 前提（出处） | 现量读数 | 判 |
|---|---|---|---|
| P1 | "`Snapshot` 恰四字段，`composer.go:44-49`；生产构造点 `:79`"（093x 特批 §1） | `:44-49` 确为四枚；`NewSnapshot` 在 `:63`、其 `return` 在 `:79` | ✅ 成立（但"构造点"这词有第二层，见 P2） |
| P2 | 同上，被当成"改一处构造点就填得上真值" | 真值不来自 `NewSnapshot`，来自**装配根的读口**：`cmd/wisp/run.go:421-427` 那枚 `PumpSources` 字面量（全树 `NewSnapshotPump(` 生产调用者**唯一一枚**＝此处）。`pump.go:182`／`panel_pump.go` 只是它的下游 | ⚠ 成立但**不足**：按 P2 的字面做，加出来的键今天**没人填**（实测见 §2.3） |
| P3 | 编排者 `A273②`："四组两两不重叠……**不碰 151 的 `run.go`**" | 候选落地集 5 枚段里 **4 枚的读口必须出现在 `run.go:421-427`**（`run`／`tools`／`cost`／`failures`；其中 `run.usage` 还要 `internal/agent/loop.go` 先转发 `EvUsage`，那是票 153 的地界）。只有 `approval.*` 一枚能走已接线的 `Verdicts` 读口旁通（本程在仓外副本真做通了，§2.4） | ❌ **不成立**：AC#2 的"填真值"这一半结构性需要 `run.go`。⇒ **报回**，不硬改 |
| P4 | 093x 特批 §2：`加字段＝扩契约` 的人工批准"只到新增字段为止" | 本程**接受为授权**，并严格按字面用：不加字段＝本程不落键；**未**据此动任何现有字段语义、未删字段、未动 `C17` 白名单条目名/类型 | ✅ 授权在效（它解的是契约批准那把锁，没解 §3 那三把） |
| P5 | 普查 §4.4："Go 加键而 `panel.ts` 没跟上 ⇒ **那把尺**会红"（单数，指 `composer_test.go:73-78`） | 实跑（仓外副本 E1）：响的是 **4 枚**用例，不是 1 枚——`approval_test.go:129`、`composer_test.go:74`、`pump_test.go:124`、`pump_test.go:276`（后两枚是票 35 的**键集钉**，普查没计） | ⚠ 成立但**数少了**：闭合的爆炸半径是 4 枚测试＋1 枚前端文件（§2.1/§2.2） |
| P6 | 普查 §2.8 行 8："`approval` 那一组**全有源，且生产者今天已经在跑**"，含倒计时 `remainingMs` | 逐枚分：**`windowMs` 有源**（`gate.Window()`，实测 `2000ms`）；**`vetoChannels` 有源**（`ChannelRegistry.Statuses()`，实测四行含 B1 原句）；**`remainingMs` 无源**——`approval.EventTick`（`ui.go:64`）**声明了、全仓零生产者**，真发的三种 Event（`gate.go:274`／`:387`／`:510`）带的 `Remaining` 是**静态窗口长／预警提前量**，不是倒计时；倒计时本体是 `gate.go:257` 的 `g.clock.After(g.window)`，**没有"还剩多少"的读口** | ❌ **部分不成立**：按普查自己的 §0.4 判据①（"只有类型没有生产者，不算 ⓐᐟ 成立"），`remainingMs` 应从**甲组挪进乙组**，与 `thinkingMs`/`reasoningMs`/`durationMs` 同族（另开票装计时器） |
| P7 | 普查 §3.2：ⓐ"消费端今天就位，Go 加一枚 `view` 前端一行都不用动" | 已提交版 `panel-views.ts` 的 `currentView()` 确读 `snapshot.view`；而 Go 侧九枚屏 id 现量**零命中**（`grep -rn --include=*.go` 限 `internal/panel/`＋`cmd/wisp/`、剔 `_test.go` → 无输出） | ✅ 成立，且**加重**：ⓐ 落地的话值只能由宿主随便给（§4.1） |
| P8 | 派单共同纪律："判跑没跑到只认 `=== RUN` 枚数" | 本程每一发门禁都带 `=== RUN` 计数（§5）：仓内基线 `internal/panel` RUN=105、`cmd/wisp` RUN=139；仓外副本各发同样带计数 | ✅ 照做 |
| P9 | "同树三枚别家在跑" | 现量：票 151 **已交件**（`63e228fc` 收表）⇒ `run.go` 现在**没人写**，但**仍不在本程放开面**（地界由特批划，不由"没人抢"划） | ✅ 成立，附一条更新 |
| P10 | 记忆里那条 owner 排序信号"前端那儿就等后端完事儿了" | 与实测冲突的一面：**前端那一侧今天缺的不是这五枚键**——五枚键落下去也仍要 `panel.ts` 同批＋泵把包送到页（最后一公里在票 33/35 名下，`pump.go:15-23` 自证"无 Go→页通道"） | ⚠ 不当授权用；登记以免下一位以为"落了字段前端就动了" |

⇒ **两条要报回的**：P3（`A273②` 那句"不碰 run.go"对 AC#2 不成立）与 P6（`remainingMs` 无生产者）。
本件把它们写成读数与出处，**不据它们硬改，也不据它们扩权**。
