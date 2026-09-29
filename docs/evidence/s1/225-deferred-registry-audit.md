# 票 225 — `DEFERRED(...)` 代码标记 ↔ `SPEC-12 §5` 登记表：双向对账（只读审计腿）

- 锚点 `git rev-parse --short HEAD` = **`49abb593`**（分支 `dev`）。
- 落笔时刻：**2026-09-29 12:07 +0800**（本节下面每一个时刻都各跑一次 `date`，不许互相推）。
- 本腿性质：**只读**。零编译（没有跑过 `go test`/`go build`/`go vet`/任何构建），零写入除本文件；
  `docs/PLAN.md`／`docs/specs/**`／源码／工单／台账**一字未动**；`frontend/**`／`design/**` **零读零引零转述**。
- 与票 223 的关系：同源不同射程——223 管接线，225 管账目；本文件**不做合并**（票 225 禁区末条）。

---

## 0. 我跑过的尺（原样可复跑）与我自己的分母

| 尺（逐字） | 我这把的读数 |
|---|---|
| `grep -rn "DEFERRED(" --include=*.go internal/ cmd/ tools/ \| grep -v _test \| wc -l` | **26 处命中** |
| `grep -rho "DEFERRED([^)]*)" --include=*.go internal cmd tools \| grep -v _test \| sort -u \| wc -l` | **22 枚去重标记** |
| `grep -rn "DEFERRED(" --include=*.go internal/ cmd/ tools/ \| grep -v _test \| grep -o "DEFERRED([^)]*)" \| sort \| uniq -c` | `C25-loop-wiring` 3、`D28-1` 3、其余 20 枚各 1（3+3+20＝26，闭合） |
| `grep -rn "DEFERRED(" --include=*_test.go internal cmd tools` | **4 处**（全在 `internal/risk/*_test.go`，都是**注释里引用**别的标记、不是新标记）⇒ 本表分母不含它们 |
| `sed -n '64,92p' docs/specs/SPEC-12-roadmap-governance.md \| grep -c "^\| DEFERRED"` | **12 行** |
| 同上段 `grep -c "^\| RESERVED"` / `grep -c "^\| REJECTED"` / `grep -c "^\| ~~DEFERRED~~"` | **7 / 6 / 2** |

**分母与票面（锚 `36b46125`）的读数完全一致：26 处／22 枚／12＋7＋6＋2 行。无差异要说明。**

口径三条，先摆明（不然下面的"枚"是糊的）：

1. **"一枚"＝去重后的标记名**（`DEFERRED(<名字>)` 括号内整串算名字，大小写敏感、`sort -u` 去重）。
   同名标记在多枚文件重复算**一枚**（`C25-loop-wiring` 三处、`D28-1` 三处），表一里逐处仍给行号。
2. **`tools/` 与 `cmd/` 都计入**（票面那把尺含 `cmd/ tools/`）；本腿实测 `tools/**` 里**零枚**标记，
   `cmd/` 里 **1 枚**（`cmd/wisp/panel_assets.go:164`）。所以"取舍"这一支差异源今天为空。
3. **§5 的"一行"＝表体里以 `| DEFERRED`／`| RESERVED`／`| REJECTED`／`| ~~DEFERRED~~` 起头的行**，
   不含表头行与 `|---|` 分隔行（票 150 AC#1 要为"29 vs 47"收的那个口径，本腿用的是这一把，可复算）。

### 结论摘要（四行）

- **表一方向**：22 枚去重标记里，**只有 3 枚**能在 §5 指到行（严格同名的**只有 1 枚**：`kws-veto, B1`）；
  **12 枚判〔该摘／该改写〕；另有 14 枚共用"这件事排在票 NN 上"的排期指针措辞**（两拨**部分重叠**，
  逐字 `implemented by ticket NN. This ticket only freezes the package boundary` 只命中 **9 枚**，
  其余 5 枚因 Go 注释**折行／换动词**躲过了那把尺，两把尺都在 §1.2 第 1 条），
  按 §5 的类型语义（`DEFERRED`＝**有计划不做**）它们**根本不是推迟项**，是"计划内还没做／已做"的指针；
  **3 枚是真欠账而 §5 无行（该补一行，须人批）**；**2 枚（`D28-1`／`D11-3`）已由票 150＋`Q-55` 立案并立为 `blocked`**，本腿不重裁。
- **表二方向（本票主要价值）**：§5 那 **12 行 DEFERRED 里只有 2 行**在代码里有对应标记
  （`macOS 平台层`、`快捷键路径语音否决（B1）`）⇒ **10 行登记了、代码里零标记**；
  **7 行 RESERVED 里只有 1 行**有标记（`Linux 支持`，与 macOS 共用那枚）⇒ **6 行零标记**。
  最刺眼的一格：**`剪贴板历史（D34）`／`doc.read xlsx/OCR（D34）`／`system.eject（D34）`三行**——
  行名里**带着 D 号**、代码里**一枚标记都没有**，正是票 150 说的"表上有人、代码里查无此人"的反向。
- **表三**：`internal/config/doc.go:14` ⇒ **改标记（摘）＋不为它补 §5 行**，欠账本身归票 223；
  `internal/watchdog/doc.go:18` ⇒ **改标记（措辞是假的：票 42 无 `-done`、`internal/watchdog/` 今天只有 `doc.go` 一枚文件）**，
  **也不要补 §5 行**（票 42 是排期内、不是"计划不做"）。两枚都**不动 `docs/specs/**`**；改源码注释不属"改契约"，但**摘掉标记会被读者读成"这件事结了"**，
  所以两枚都必须与票 223／票 42 的落点一起核销，见 §3。
- **AC#4**：**只出形状，本腿没上任何仪器**（`tools/` 今天只有 `d22scan`／`mockllm`／`signmodels` 三枚程序，
  `grep -rln "DEFERRED" tools/` ＝**0 命中**，与票面一致）。要做常驻检要先摆给 owner，理由见 §4。

---

## 1. 表一：代码 → §5 登记表（AC#1，22 枚去重标记／26 处出现逐条）

> 三态用票面那三个词。**⚠ 本腿加了一枚口径注记（`〔摘〕`／`〔改写〕`）**：票面三态里没有"用错类型"这一支，
> 而 22 枚里最大的一族（12 枚）恰好落在这里。我的处理：**按"标记本身过期该摘"计**，但在"凭什么"那列写清它是
> "自称由 NN 实现"的占位，**不是"这件事已做完"**。这一格是我加的裁定，依据与后果见 §6 第 1 条。

| # | 文件:行 | 标记原文（去重名） | 它自称由谁实现 | 那枚票今天（判据＝有无 `-done` 后缀，现查 `.scratch/wisp/issues/`） | §5 对应行 | 结论三态 |
|---|---|---|---|---|---|---|
| 1 | `internal/agent/approval/approval.go:99` | `DEFERRED(kws-veto, B1)` | 否决词归票 41、面板按钮归票 37 | **41 无 `-done`／37 无 `-done`** ⇒ 都还开着 | **有**：§5 `| DEFERRED | 快捷键路径语音否决（B1）…`，标记**逐字抄了那一整行**（含五字段） | **〔§5 有对应行〕** 22 枚里唯一一枚"名字＋五字段"双向都咬合的 |
| 2 | `internal/agent/approval/doc.go:13` | `DEFERRED(queue)` | 最小门＝票 21（本包）、多任务路由那半＝票 48 | **21 无 `-done`／48 无 `-done`** | **无**。⚠ 标记自己逐字写了 `-> SPEC-12 §5 has no row for it` | **〔§5 没有该补一行〕**（诚实的欠账：作者当场承认没行；补不补须人批，见下"该补三枚"） |
| 3 | `internal/agent/compress.go:27` | `DEFERRED(D28-1)` ① | 把 `Compress` 调用移进 Warm-window hook（票 28） | **28 无 `-done`** | **无**（§5 段内 `D28` 命中 0 次） | **已由票 150 立案**（`blocked`，等 `Q-55`）⇒ **本腿不重裁**，不新增第 4 种结论 |
| 4 | `internal/agent/compress.go:53` | `DEFERRED(D28-1)` ② | 同上（叙述性引用，同一件事第二处） | 同上 | 同上 | 同上（票 150 计"出现次数"里的一枚） |
| 5 | `internal/agent/control.go:15` | `DEFERRED(D11-3)` | 逐字 `see the ticket report`（票内编号，不指任何 `-done` 件） | 票 10 **有 `-done`**；`D11(3)` 的 `force_tool/force_chat` TOML 规则表**没有专属票** | §5 只有 `D11` 的**一行 REJECTED**（`意图分类四级流水线 D11 过度设计已废`）⇒ **判不同件事**（见 §6 第 3 条依据） | **已由票 150 立案**（`blocked`，等 `Q-55`）⇒ 本腿不重裁 |
| 6 | `internal/agent/doc.go:14` | `DEFERRED(loop)` | 逐字 `implemented by ticket 10. This ticket only freezes the package boundary` | **10 有 `-done`**，且 `loop.go` 今天是真有实现（394 行那一处在跑） | **无** | **〔标记本身过期该摘〕**——已结案＋确已实现＝纯残留 |
| 7 | `internal/agent/loop.go:394` | `DEFERRED(D28-1)` ③ | 同 #3（同步兜底路径的注释） | 同 #3 | 同 #3 | 同 #3（票 150） |
| 8 | `internal/agent/scheduler/doc.go:13` | `DEFERRED(scheduler)` | `implemented by ticket 47. This ticket only freezes the package boundary` | **47 无 `-done`** | **无**。⚠ 最近的一行是 §5 `| REJECTED | 周期性调度器（cron）…`——**类型相反**（"别做"vs"由 47 做"） | **〔标记本身该摘／改写〕** ＋ 单独列为**危险形状**（见 §1.2 末条） |
| 9 | `internal/audio/doc.go:27` | `DEFERRED(playback)` | `TTS output lands with ticket 26` | **26 无 `-done`** | **无**。同段还提"AEC source 归票 26/59"，而 §5 有 AEC 那一行（`~~DEFERRED~~ 已采纳（Path C）`）——**但那句话不是 `DEFERRED(...)` 标记**，仪器按名字扫它会漏 | **〔标记本身该摘／改写〕**（排期内、不是"计划不做"） |
| 10 | `internal/config/doc.go:14` | `DEFERRED(schema/hot-reload/migration)` | `implemented by ticket 05. This ticket only freezes the package boundary` | **05 有 `-done`**（票面标题含 `hot reload`） | **无**（§5 没有任何热加载／配置迁移行） | **→ 表三 §3.1**（这枚是票 223 误判的源头） |
| 11 | `internal/llm/probe.go:22` | `DEFERRED(audio, ticket 09)` | C7 无 audio content part ⇒ 请求定义归语音级联票 15/60/61 | **09 有 `-done`**；**15／60／61 全无 `-done`** | **无**（§5 的 `云端 ASR/TTS provider` 那行是"已采纳"，讲的不是 C7 探针） | **〔§5 没有该补一行〕**——22 枚里**唯一一枚"生产里会响"的真推迟**（`RunProbe(audio)` 报 `OK=false` 且写明理由，不假过） |
| 12 | `internal/memory/doc.go:17` | `DEFERRED(sqlite/schema)` | 存储核心＝票 04、L1/L2 抽取流＝票 29 | **04 有 `-done`**；**29 无 `-done`** | **无** | **〔标记本身该摘／改写〕**——一枚标记**混装**"已落（04）＋未落（29）"两半，`DEFERRED` 这个词只对未落那半勉强成立 |
| 13 | `internal/observe/doc.go:22` | `DEFERRED(JSONL sink/redaction/diagnostics/CostMeter)` | 逐字 `implemented by ticket 08` | **08 有 `-done`**；但名字里四件事的两件另有票：**44（CostMeter C23）／45（diagnostics guards）都无 `-done`** | **无** | **〔标记本身该摘／改写〕**——自称"由 08 全实现"而名字含 44/45 的活 ⇒ 与 #10 同一类**措辞过度**（本腿只核到"名字含 4 件事、其中 2 件另有开着的票"，没逐行核 08 的交付面，见 §6 第 5 条） |

| # | 文件:行 | 标记原文（去重名） | 它自称由谁实现 | 那枚票今天（判据＝有无 `-done` 后缀） | §5 对应行 | 结论三态 |
|---|---|---|---|---|---|---|
| 14 | `internal/panel/doc.go:16` | `DEFERRED(host/bridge)` | `implemented by ticket 33 (host), ticket 35 (bridge)` ＋ `only freezes the package boundary` | **33 无 `-done`／35 无 `-done`** | **无**（§5 `RESERVED | MCP client` 那行的"与 host bridge 收口冲突"是**另一件事**） | **〔标记本身该摘／改写〕** |
| 15 | `internal/plugin/doc.go:17` | `DEFERRED(manifest/goja)` | 票 50（Tier-1）／51（Tier-2 goja），本票只落 `DisposalScope` | **50 无 `-done`／51 无 `-done`** | **无**（§5 有 `插件 SDK+文档+示例`、`社区插件 registry` 两行，但那是**票 57** 的面，与 50/51 不同事，见 §2） | **〔标记本身该摘／改写〕** |
| 16 | `internal/proc/doc.go:22` | `DEFERRED(update/rollback)` | D41(b) staging/原子替换/回滚，归发行票 56 | **56 无 `-done`** | **无同名行**。最近邻＝§5 `| DEFERRED | 代码签名/包管理器分发`，**同一枚前置票 56** ⇒ 我判**不同件事**（签名/安装渠道 ≠ 自更新与回滚），依据见 §6 第 3 条 | **〔§5 没有该补一行〕**（自更新/回滚是独立能力面，今天**整仓零标记**指向 §5 任一行） |
| 17 | `internal/risk/doc.go:15` | `DEFERRED(R1-R9/PathResolver/Provenance)` | `implemented by tickets 17, 18 and 19` ＋ 边界冻结句 | **17／18／19 三枚全有 `-done`** | **无**（`R1–R9` 在 §5 不成行，它们在 `AGENTS.md §1.1` 是**禁止清单**、不是推迟项） | **〔标记本身过期该摘〕**——已结案＋确已实现；⚠ 但别把它和 #19 那枚 `C25-loop-wiring` 一起摘（后者才是没接完的那半） |
| 18 | `internal/risk/pathresolver_other.go:5` | `DEFERRED(macOS/Linux)` | SPEC-06 §4 的 realpath＋lstat 软链拒绝，非 Windows 平台 | 票 55（S8 macOS port）**无 `-done`** | **有（父项口径）**：§5 `| DEFERRED | macOS 平台层`＋§5 `| RESERVED | Linux 支持` | **〔§5 有对应行〕** ⚠ 判"同一件事"靠的是 §5 那行的**完成判据**（macOS 跑不起来＝这一枚的现状表现），**不是字符串同名**（见 §6 第 3 条） |
| 19 | `internal/risk/provenance.go:58` | `DEFERRED(C25-loop-wiring)` ① | 票 19 的欠账第 (1) 条；形状由票 160 落，**"Defer 到 `DisposalScope`"那半仍 OPEN** | **160 无 `-done`**；活跃追票 `151`/`158`/`175`/`176`（**均无 `-done`**） | **无**（最近邻 §5 `RESERVED | 精确 taint tracking` 是"粗粒度已够"的取舍，**不是这条接线债**） | **〔标记本身该摘／改写〕**——它在排期内且有票在追 ⇒ 按 §5 语义不该叫 `DEFERRED`；⚠ 这枚**内容是真的没接**（票 223 同族的另一根线），摘字不等于销账 |
| 20 | `internal/risk/syncdirs_other.go:17` | `DEFERRED(P12-macos)` | macOS 同步客户端探测，归票 55，AC 写进票 55 的 Constraints 补录 | **55 无 `-done`** | **有（父项口径）**：§5 `| DEFERRED | macOS 平台层`（与 #18 共用同一行） | **〔§5 有对应行〕**（同一行的第二枚标记 ⇒ 对任何"1:1"仪器都是**一对多**形状，见 §4） |
| 21 | `internal/session/doc.go:16` | `DEFERRED(SessionScope)` | `implemented by ticket 28. This ticket only freezes the package boundary` | **28 无 `-done`** | **无**（⚠ 别和票 224「会话级授权」混，那是授权档、不是这枚） | **〔标记本身该摘／改写〕** |
| 22 | `internal/speech/doc.go:19` | `DEFERRED(engines)` | 票 15（ASR＋CER harness）／26（TTS）／41（KWS）＋边界冻结句 | **15／26／41 全无 `-done`** | **无**（§5 的 AEC 那行与 `云端 ASR/TTS provider` 那行都是**另一件事**，后者已采纳） | **〔标记本身该摘／改写〕** |
| 23 | `internal/tools/bridge.go:687` | `DEFERRED(C25-loop-wiring)` ② | 同一枚名字的第二处（`OpenTask` 的形状由票 160 落） | 同 #19 | 同 #19 | 同 #19 |
| 24 | `internal/tools/doc.go:47` | `DEFERRED(tools)` | 票 20 起（20/22/23/24 加 web/system/doc 工具）＋边界冻结句 | **20／22／23／24 全无 `-done`** | **无**（⚠ 本票**禁区**指到 `docs/TOOLS.md` 不存在；`D34` 那张权威工具表在 `PLAN.md:2514`，§5 里只有 `doc.read` xlsx/OCR 与 `system.eject` 两枚**子工具**行） | **〔标记本身该摘／改写〕** |
| 25 | `internal/watchdog/doc.go:18` | `DEFERRED(watchdog loop/thresholds)` | `implemented by ticket 42. This ticket only freezes the package boundary` | **42 无 `-done`** | **无**（`grep -n "看门狗\|watchdog"` 全篇 SPEC-12 只命中 S6 那行**切片表**，§5 里零命中） | **→ 表三 §3.2** |
| 26 | `cmd/wisp/panel_assets.go:164` | `DEFERRED(C25-loop-wiring)` ③ | 同一枚名字第三处（CLI 探针故意不持有 closer） | 同 #19 | 同 #19 | 同 #19 ⚠ `cmd/` 里唯一一枚，摘它会连带改 `cmd/wisp` 的注释射程 |

### 1.1 表一的三态计数（口径：去重名 22 枚；括号内是按出现处 26 计的数）

| 结论 | 枚数 | 具名 |
|---|---|---|
| 〔§5 有对应行〕 | **3**（出现 3 处） | `kws-veto, B1`、`macOS/Linux`、`P12-macos`（后两枚共用 §5 的 `macOS 平台层` 那一行） |
| 〔§5 没有该补一行〕 | **3**（出现 3 处） | `queue`、`audio, ticket 09`、`update/rollback` |
| 〔标记本身过期该摘／该改写〕 | **12**（出现 **14** 处，占 26 处的 **54%**） | `loop`、`R1-R9/PathResolver/Provenance`、`sqlite/schema`、`JSONL sink/redaction/diagnostics/CostMeter`、`host/bridge`、`manifest/goja`、`SessionScope`、`engines`、`playback`、`tools`、`scheduler`、`C25-loop-wiring`(×3 处) |
| 已由票 150＋`Q-55` 立案，本腿不重裁 | **2**（出现 4 处） | `D28-1`(×3)、`D11-3`(×1) |
| 进表三单独结论 | **2**（出现 2 处） | `schema/hot-reload/migration`、`watchdog loop/thresholds` |

**枚数闭合**：3＋3＋12＋2＋2 ＝ **22 枚**＝去重分母 ✓
**处数闭合**：3＋3＋14＋4＋2 ＝ **26 处**＝命中分母 ✓
（26−22＝4 的全部来源只有两枚同名多现：`C25-loop-wiring` ×3、`D28-1` ×3，各多计 2 处；其余 20 枚各 1 处）

### 1.2 从表一抽出的三句硬话（都是现读）

1. **"排期指针"那一族：两把尺数出两个数，都摆出来**（⚠ 我第一把尺差点把它写成"12 枚共用同一句话"，那是**过裁**，这里按现跑改成两把）
   - **尺 A（同行逐字）**：`grep -rn "implemented by ticket" --include=*.go internal cmd tools | grep -v _test` ＝ **11 命中**，
     减掉**不是标记**的两处包头散文（`internal/proc/doc.go:5`、`internal/statemachine/states.go:5`）＝ **9 枚**标记：
     `loop`、`scheduler`、`schema/hot-reload/migration`、`host/bridge`、`manifest/goja`、`R1-R9/PathResolver/Provenance`、
     `SessionScope`、`engines`、`watchdog loop/thresholds`。
   - **尺 B（含折行／换动词）**：再加 **5 枚**＝`JSONL sink/redaction/diagnostics/CostMeter`（`implemented by`＋换行＋`ticket 08`）、
     `sqlite/schema`（动词是 `landed with ticket 04`）、`tools`（`implemented from ticket 20 … onward`）、
     `playback`（`lands with ticket 26`）、`queue`（`landed in ticket 21 … belongs to ticket 48`）⇒ **同族共 14 枚**。
   ⇒ **结论：14 枚（去重名的 64%）说的是"这件事排在哪枚票上"，不是"这件事计划不做"**；
   而 §5 的类型语义（`SPEC-12:61`）逐字写的是 **`DEFERRED` = 有计划不做** ⇒ **按 §5 自己的定义，这一族今天就不该用 `DEFERRED` 这个词**，
   所以"§5 没有该补一行"对它们是**错的修法**（补 14 行会把"计划内"污染成"推迟"）。
   其中自称的票**全部已结案**的＝**4 枚**（`loop`＝10、`R1-R9/…`＝17/18/19、`schema/hot-reload/migration`＝05、`JSONL …`＝08），
   **混装已落＋未落两半**＝3 枚（`sqlite/schema` 04+29、`JSONL …` 08 但名字含 44/45、`queue` 21+48），**其余指着的票都还开着**。
   ⚠ **顺带一枚仪器教训**：尺 A 因为 **Go 注释折行**漏掉 5 枚（同族 14 里只有 9 被扫到）——**任何按行扫 `DEFERRED(` 的仪器都会少报三分之一**，见 §4.3 第 7 条。
2. **§5 第 94 行规定的那个形状，全仓零枚**：`docs/specs/SPEC-12:94` 逐字要求
   `// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点`；而 **`docs/DEFERRED.md` 不存在**（`ls docs/` ＝ `BUILD.md HOTKEYS.md PLAN.md PRECHECK.md SLO.md evidence reports specs tickets`），
   代码里**没有任何一枚**带 `docs/DEFERRED` 字样（`grep -rn "docs/DEFERRED" --include=*.go internal cmd tools \| grep -v _test \| wc -l` ＝ **0**）。
   ⇒ 那条"1:1"规矩**连它自己指定的锚点文件都没有**；`AGENTS.md §4` 末段已经声明这批交付物"截至锚点 `4e66817` 在仓里不存在"，与此吻合。
3. **一枚"类型相反"的形状**（本表最容易被忽略的一条）：代码里有 `DEFERRED(scheduler)` 自称"由票 47 实现"，
   而 §5 对调度器的正式态度是 **`REJECTED | 周期性调度器（cron）`**（＝有意取舍、**不得被"修复"**）。
   一次性 `reminder` 那半已在 §5 那行里写明是"窄例外"。⇒ 后来的人照 §5 反查会得出"这包不该存在"，
   照代码读会得出"这包有人排期"。**两个方向都不是票面三态能表达的**，我按〔该摘／改写〕计，并把裁量权交回。

## 2. 表二：§5 登记表 → 代码（AC#2，反向；本票主要价值）

时刻 **2026-09-29 12:11 +0800**（本节节头现跑 `date`）。

**先摆一枚口径（不做这条，下面这张表会被读歪）**：`AGENTS.md §1.1` 那句原文是
"**`DEFERRED(D-xx)` 代码标记**必须与 `SPEC-12 §5` 登记表 1:1 双向对得上"——**字面只要求 `DEFERRED` 一型**。
⇒ **`RESERVED`／`REJECTED` 那 13 行本来就不在"必须有标记"的射程内**；我这一把还量到：
**代码里 `RESERVED(...)` 形状零枚**（`grep -rn "RESERVED(" --include=*.go internal cmd tools | grep -v _test | wc -l` ＝ **0**），
只有两枚**散文式类型注**（`internal/tools/registry.go:55` 的 `KindMCP → "REJECTED: D13/16.9#7, interface slot only"`、
`internal/risk/taintmatch.go:15` 的 `// REJECTED (16.9#1) and must not be resurrected here`）。
⇒ **下面 RESERVED 那 7 行是我按票面（"12 行 DEFERRED＋7 行 RESERVED 逐枚问"）额外走的第二把尺，不是 §1.1 那句话要求的**。

### 2.1 §5 的 12 行 `DEFERRED`（行号＝`docs/specs/SPEC-12-roadmap-governance.md` 现读）

| §5 行 | 项（逐字） | 代码里有对应标记吗 | 最近的代码证据（现读） | 结论 |
|---|---|---|---|---|
| 66 | macOS 平台层 | **有** | `internal/risk/pathresolver_other.go:5` ＋ `internal/risk/syncdirs_other.go:17`（共用这一行） | 〔已证（现读）〕唯一有标记的 DEFERRED 行之一 |
| 67 | 代码签名/包管理器分发 | **无** | `tools/signmodels/main.go:1` 是 **C29 模型清单 minisign 签名**（票 14），**不是**应用二进制 SignPath/Authenticode；`Authenticode\|winget\|scoop` 在 Go 里**零命中** | 〔仅文档写了、代码没有〕登记了但代码零标记 |
| 68 | 插件 SDK+文档+示例 | **无** | `internal/plugin/doc.go:17` 是 `DEFERRED(manifest/goja)`＝票 50/51，另一件事 | 〔仅文档写了、代码没有〕 |
| 69 | 社区插件 registry | **无** | `internal/tools/registry.go` 有 `Kind` 位，但没有 registry 下载/索引的标记或注 | 〔仅文档写了、代码没有〕 |
| 70 | i18n/非中文 | **无** | `internal/agent/control.go:35` 一句散文 `expanding the vocabulary is an i18n/UX decision, ticket 58`（票 58 无 `-done`）；`internal/ball/renderer_windows.go:166` 逐字写死 `zh-CN` | 〔建了但没接〕硬编码 locale 在跑、i18n 无标记 |
| 71 | 无障碍 | **无** | `a11y`／`accessibility` 在 Go 里**零命中** | 〔仅文档写了、代码没有〕 |
| 73 | 快捷键路径语音否决（B1） | **有** | `internal/agent/approval/approval.go:99`（逐字引 §5 那一整行） | 〔已证（现读）〕全表**唯一**一枚同名同字的双向咬合 |
| 74 | 剪贴板历史（D34） | **无** | `internal/tools/capability.go:22` 有 `CapClipboard`（当次读写的能力位，与 §5"只有当次读写"吻合），**历史**零标记 | 〔建了但没接〕能力位在、欠的那半无标记 |
| 75 | `doc.read` xlsx/OCR（D34） | **无** | `xlsx` **零命中**、`OCR` 2 枚文件（本腿未逐行读那两处的语义） | 〔仅文档写了、代码没有〕 |
| 76 | `system.eject`（D34） | **无** | 票面尺 `grep -rli eject` 会命中 69 枚文件＝**`project`/`reject` 的子串假阳性**，具名工具零命中 | 〔仅文档写了、代码没有〕⚠ 这条我给的是"关键词尺不可信"的读数 |
| 77 | 完整错误文案体系 | **无** | `internal/observe/errors.go` 存在（分类在跑），用户可见文案那半无标记 | 〔建了但没接〕 |
| 79 | 竞品对比文档 | **无** | README 侧（本腿没读 README，票面 §5 的完成判据落在 README 第一段） | 〔仅文档写了、代码没有〕按设计也不该有代码标记 |

**小结：12 行 DEFERRED ⇒ 只有 2 行（66、73）在代码里有 `DEFERRED(...)` 标记 ⇒ 10 行"登记了、代码里查无此人"。**
⚠ 其中 **67／68／69／71／79 五行本来就不该期待代码标记**（它们的前置是 S7/S8 的票或纯文档）；
**真正"代码就在手边却没标"的是 74（剪贴板历史）、75（doc.read xlsx/OCR）、76（system.eject）、77（错误文案）、70（i18n）** 五行——
它们的**相邻代码今天活着**（`CapClipboard`、`zh-CN` 硬编码、`observe/errors.go`），下一位读到的是"这儿没欠账"。

### 2.2 §5 的 7 行 `RESERVED`（票面额外要求的一把尺）

| §5 行 | 项（逐字） | 代码里有对应标记/占位吗 | 现读证据 | 结论 |
|---|---|---|---|---|
| 80 | MCP client | **有类型注、无标记** | `internal/tools/registry.go:55` `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}` | ⚠ **类型相反**：§5 那行是 **RESERVED**、代码里那格写的是 **REJECTED**（同一根 `16.9#7` 出处，两个类型词） |
| 81 | Codex guardian 双轴评分 | **有散文注、无标记** | `internal/risk/assessor.go:169` `SubAssessor is the pluggable assessor interface. RESERVED for guardian` | 〔建了但没接〕接口位在，评分器零实现 |
| 82 | 跨会话持久授权档（workspace/user） | **无** | 票 49 `49-session-grants-d45-2.md` **无 `-done`**；`internal/perm`/`approval` 里无 `D45-2` 标记 | 〔仅文档写了、代码没有〕 |
| 83 | `available_decisions`（服务端下发按钮） | **无** | `available_decisions` 在 Go 里**零命中** | 〔仅文档写了、代码没有〕 |
| 84 | embedding 语义记忆 | **无** | `embedding` 在 `internal/memory` **零命中**；§5 那行"前置＝`MemoryStore` 接口位"——**`type MemoryStore` 全仓零命中**（该包有的是 `Memory` 结构体，`internal/memory/models.go:42`） | 〔仅文档写了、代码没有〕⚠ §5 的"前置"那格叫不出代码里的名字 |
| 85 | Linux 支持 | **有（与 macOS 共用）** | `internal/risk/pathresolver_other.go:5` 那枚名字逐字是 `macOS/Linux`；`//go:build !windows` 文件另有票 78/82/113/119/120 在追 | 〔已证（现读）〕一枚标记喂两行（RESERVED＋DEFERRED），**非 1:1** |
| 86 | 精确 taint tracking | **有类型注、无标记** | `internal/risk/taintmatch.go:15` `// REJECTED (16.9#1) ...`（而 §5 这行是 RESERVED；REJECTED 的是"token 级"那一层） | 〔建了但没接〕粗粒度在跑、细粒度按取舍不接 |

**小结：7 行 RESERVED ⇒ `DEFERRED(...)` 形状零枚**（只有 4 处散文类型注，其中 2 处的**类型词与 §5 不一致**）。

### 2.3 那 2 行 `~~DEFERRED~~ → 已采纳` 要不要算进分母：**我判"不算必须有标记的分母，但要算表"**

- §5:72（AEC/barge-in 陪聊路径）：已采纳为 S4 必做 ⇒ **不是欠账**；但同一行逐字留着
  "仍 DEFERRED 的只有干活路径全双工（设计非债务）"，而**代码里没有任何标记指向它**——
  最接近的是 `internal/audio/doc.go:24-25` 的散文 `no AEC source (Path C AEC is tickets 26/59...)`（**不是 `DEFERRED(...)` 形状**，任何按名字扫的仪器都看不见它）。票 59 `59-p15-aec-spike.md` **无 `-done`**。
- §5:78（云端 ASR/TTS provider，级联 C9）：已随票 61 落地首个 provider ⇒ **不该再要求标记**；票 61 **无 `-done`**（本腿未核它的实现进度）。
- **裁定与理由（一句）**：已采纳＝"计划做/在做"，§5 的类型语义只给"计划不做"发标记义务 ⇒ **计入表、不计入"缺失"那栏**。
  ⚠ 但这 2 行留在表里有实用价值：**它们各自还有一句"仍缺的那半"没有代码痕迹**（干活路径全双工＝设计性缺失），
  这正是票面说的"最容易被后来的人当成没这回事"。

**6 行 REJECTED：本腿不判缺失。**理由：`REJECTED`＝"有意取舍、不得被修复"，给它配代码标记等于把"别做"再广播一遍，
且 §5:61-62 的类型语义没要求它进 1:1；唯一值得记的是 §5:91 那行与表一 #5（`DEFERRED(D11-3)`）的**语义撞车**，已写在 §1.2 第 3 条。

## 3. 表三：票面点名的那两枚确凿失配（AC#3）

时刻 **2026-09-29 12:12 +0800**（本节现跑 `date`）。

### 3.1 `internal/config/doc.go:14` — `DEFERRED(schema/hot-reload/migration): implemented by ticket 05. This ticket only freezes the package boundary.`

**现读事实（四条，都可复跑）**

1. 票 05 的文件名逐字是 `05-config-model-done.md` ⇒ **已结案**（`-done` 是防重领唯一键，我没看它里面的 Status 行）。
2. `Manager.CheckAndReload` **确有实现**：`internal/config/manager.go:117`（注释 `:113` 自称 `SPEC-03 sec 4.3, watchdog 1s tick`）。
3. **非测试调用者只有一枚**：`cmd/balldebug/main.go:243`；`cmd/wisp/**` 里 `grep -rn -i "reload" --include=*.go cmd/wisp | grep -v _test` ＝ **0 命中** ⇒ 生产二进制里没人轮询。
4. §5 里**没有**任何"热加载／配置迁移"行（表一 #10）；SPEC-03 有规格文本（`docs/specs/SPEC-03*.md:73 ### 4.2 热加载范围`）。

**结论（一句话）**：**该改标记、不该补登记**——这枚标记的"由票 05 实现"在"函数存在"这个意义上是真的、
在"这件事结了"这个意义上是假的，而它**本来就不该出现在 `DEFERRED` 名下**（票 05 是排期内已做完的活，不是"计划不做"）。
⇒ 修法＝**摘掉/改写那一句**（改成不承诺"已实现"、也不冒充推迟项），**不往 §5 加行**；
真正没接的那半（谁在轮询）是**票 223 的射程**，它已经开着、且已有 `cmd/balldebug` 之外的第二处承重注释在指一个不存在的主机
（`internal/ball/hotkey_reload.go:29` 逐字 `the host's poll (CheckAndReload, driven by the watchdog tick)`——见 §3.2，那台"watchdog tick"今天在生产里零拉起者）。

**要不要人批准**：改的是**源码注释**、不是 `docs/specs/**`、不是 `PLAN.md`、不是 `C1–C32`/`D1–D47` ⇒
**按 `SPEC-12 §4.1` 的字面射程这不构成契约变更**（票 150 那程的编排者在 `Q-55` 的更正里已经把"乙＝补行才叫动契约"这句收回去了，
现读 `docs/reports/pending-and-issues.md:1094` 末段①）。
**但**它动的是 `AGENTS.md §1.1` 那条 1:1 规矩的**账目本体**，且**摘标记会被下一位读成"这件事全结了"** ⇒ 我的处置建议是：
**不由本腿执行、也不单独特批，而是并进票 223 的接线落点一起核销**（同一次改动里"接上轮询＋摘掉骗人的标记"，两件事一起变绿）。
**落点**：票 223（`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`，今天无 `-done`）；
台账那一行由编排者追加（现读最新是 `A427`，`docs/reports/pending-and-issues.md:9227`；`grep -c A428` ＝ **0** ⇒ 下一个空位＝**`A428`**，本腿不写）。

### 3.2 `internal/watchdog/doc.go:18` — `DEFERRED(watchdog loop/thresholds): implemented by ticket 42. This ticket only freezes the package boundary.`

**现读事实（四条）**

1. 票面文件名 `42-watchdog.md` ⇒ **无 `-done`＝还开着**（标题逐字：`# 42 — Watchdog: per-state thresholds, safe unloadables, settle enforcement, alerts`）。
2. ⚠ **`internal/watchdog/` 目录今天只有 `doc.go` 一枚文件**（`ls internal/watchdog/` ＝ `doc.go`）⇒ 这枚标记说的"loop/thresholds 由票 42 实现"是**未来式**，包里**连实现都没有**。
3. **全仓零枚生产文件 import 它**：`grep -rn "wisp/internal/watchdog" --include=*.go . | grep -v _test` ＝ **0 命中**；
   `watchdog` 这个词在 `internal/statemachine/events.go:22,56` 只以**事件名**出现（`EvWatchdogOverrun`/`EvWatchdogFatal`），事件有名字、没人发。
4. §5 里**没有**看门狗行：票面尺 `grep -n "看门狗\|watchdog" docs/specs/SPEC-12-roadmap-governance.md` 只命中 §2 的 S6 切片表（我复跑同结论）。

**结论（一句话）**：**该改标记、不该补登记**——理由和 §3.1 同形（票 42 在排期内＝"计划要做"，§5 的 `DEFERRED` 语义＝"计划不做"），
但这枚比 §3.1 **更不该被摘**：它是今天**唯一一处**在源码里喊出"watchdog 还没接进生产拓扑"的痕迹。

**"给 watchdog 建一圈"还是"`cmd/wisp` 自己起"——本腿不裁，但交出一枚能替它裁的现读事实**：
`internal/ball/hotkey_reload.go:23-29` 已经把**两条**都写进注释了（一条是 `go r.Run(ctx, time.Second)`＋宿主自起轮询、
一条是"or, on a host that already polls config every tick (**the wisp watchdog**), call `r.Check()` from that tick"），
且 `internal/config/manager.go:113` 也把 `CheckAndReload` 的节拍**逐字命名为 `watchdog 1s tick`**。
⇒ **规格与两处注释今天已经把"轮询的宿主＝watchdog"这件事写成前提了**；
选"`cmd/wisp` 自己起一圈"就要同时改掉那三处指认（`manager.go:113`／`hotkey_reload.go:29`／`watchdog/doc.go:18`），
选"给 watchdog 建一圈"则 `internal/watchdog` 那枚"只有 doc.go"的空包就成了票 42 的必然落点。
**这个取舍要不要人批准**：**要**——它不是注释级别，票面 `AGENTS.md §2` 的"未定义即停"清单虽不含它，但
`PLAN.md` 的 **D25/D38**（进程拓扑／并发线程模型）是决策面，**动拓扑＝动 D 号的落地解释** ⇒ 归**编排者＋owner**，
且这一问的**权威落点是票 223（不是本票）**；本票只把它记成"表三的第二枚结论的承重条件"。
**台账落点**：同上，`A428` 之后的下一枚（由编排者追加，本腿不写）。

## 4. AC#4：要不要上仪器 —— **只出形状，本腿零实现**

**先声明**：本腿**没有**新建任何 `tools/**` 程序、**没有**改 `tools/d22scan`、**没有**往 CI 里加任何一步。
现读佐证：`ls tools/` ＝ `d22scan mockllm signmodels`（三枚，与本腿开工前一致）、
`grep -rln "DEFERRED" tools/` ＝ **0 命中**（与票面一致：今天没有仪器读这条 1:1）。
⇒ **做不做、放哪、什么射程＝新增一台仪器＝新射程，必须先摆给 owner**；本文件只是把形状写清，**不替他拍**。

### 4.1 要做的话，扫什么

- **代码侧尺**（本腿用的这把，需先钉成唯一口径）：
  `grep -rn "DEFERRED(" --include=*.go internal/ cmd/ tools/ | grep -v _test` → 出现处；
  再去重成枚：`grep -rho "DEFERRED([^)]*)" ... | grep -v _test | sort -u`。
  **去重必须是尺的一部分**，否则 `C25-loop-wiring` 三处会被报成三枚欠账。
- **登记侧尺**：`sed -n '/^## 5\./,/^## 6\./p' docs/specs/SPEC-12-roadmap-governance.md` 里的表体行，
  且**先声明含不含**表头/`|---|`分隔行/`REJECTED` 行/`~~DEFERRED~~` 行——票 150 的"29 vs 47"就是没先钉这一句才打的架。

### 4.2 双向各怎么判（我的形状，不是我的决定）

- **正向（代码→表）**：枚名 → §5 行的映射**不能靠字符串**。22 枚里只有 **1 枚**（`kws-veto, B1`）在标记里逐字抄了 §5 那一行，
  其余 21 枚**没有任何机械可对齐鲁**（`macOS/Linux` vs `macOS 平台层`、`update/rollback` vs `代码签名/包管理器分发`）。
  ⇒ 正向要过，得先有一枚**人工维护的别名字典**（标记名 → §5 行），而**这枚字典本身就是一个新的真相源**，得有人管、有人批——
  **这是本票交给 owner 的核心一问，不是"写个正则"能了的事。**
- **反向（表→代码）**：只对 `| DEFERRED` 那 12 行要求"至少一枚标记或至少一处类型注"，
  **`REJECTED`／已采纳（`~~DEFERRED~~`）不该要求标记**（§2.3 的理由：给"别做"发标记义务等于把它再广播一遍）。
  按本腿现读，反向今天**红 10/12**，其中**5 行本就不该期待标记**（§2.1 小结）⇒ **一上仪器就先红五枚假账**，这是要提前跟 owner 讲清的代价。
- **承重问答**（照票 150 AC#3 那句原样答一次）：**摘掉这枚检，是否存在一发漂移从此打不红？**
  ⇒ **今天是"是"**：新加一枚 `DEFERRED(D99-1)` 不会有任何东西变红（零仪器读 `DEFERRED(`，见 §4 开头那把尺）。

### 4.3 会不会误报（我量到的七种真实形状）

1. **同一件事的两种拼法**＝两枚标记、一行登记（`macOS/Linux` ＋ `P12-macos` → §5:66）⇒ **1:1 这个措辞本身在这仓不成立**，得写成 **N:1 允许**。
2. **一枚标记喂两行**＝`macOS/Linux` 同时覆盖 §5:66（DEFERRED）与 §5:85（RESERVED）⇒ 反向按行查会重复计一次。
3. **测试文件里的"提及"**＝4 处（全在 `internal/risk/*_test.go`，逐字是在**引用别的标记**）⇒ 尺不带 `grep -v _test` 就把 26 变 30。
4. **`cmd/` 与 `tools/` 的取舍**＝`cmd/` 1 枚、`tools/` 0 枚 ⇒ 少扫一个目录就少一枚，本腿两把尺都给的是**含 `cmd` 含 `tools`** 的口径。
5. **关键词尺的假阳性**＝票面 §5:76 `system.eject` 用 `grep -rli eject` 会命中 69 枚文件（全是 `project`/`reject` 的子串）
   ⇒ 若把尺从 `DEFERRED(` 放宽到"提到 DEFERRED/RESERVED/REJECTED 就算"，会**大量误报**（同类教训）。
6. **票号 vs 决策号撞车**＝`D28-1`/`D11-3` 这种"D 号＋子号"在 `docs/specs/**`＋`PLAN.md` 里从未被定义（票 150 现量：全仓只有 `D45-1`/`D45-2` 两个先例），
   仪器若拿 `DEFERRED\(D[0-9]+\)` 当形状，会把**票内编号**误读成**决策号** ⇒ 这一族已经归 `Q-55`，**仪器不该替它判**。
7. **注释折行让"按行扫"的仪器少报三分之一**（本腿现跑出来的，不是推的）：同族 14 枚里只有 **9 枚**被
   `grep "implemented by ticket"` 扫到，`observe`/`memory`/`tools`/`audio`/`approval` 五枚因**换行或换了动词**躲过
   （例：`internal/observe/doc.go:22-23` 逐字是 `…CostMeter): implemented by` 然后**换行** `// ticket 08.`）。
   ⇒ 任何"标记名 ↔ §5 行"的常驻检若只看**标记名**不受影响（名字在第一行就闭合了），
   但**任何想连"它自称由谁实现"一起读的检**必须把整段注释当**一个单元**取（`//` 连续块），否则§这一列会成片为空。

---

## 5. 与既有账的关系（避免本文件被当成"第一次发现"）

- **票 150**（`150-four-deferred-markers-have-no-row-in-the-spec-12-registry-…md`，**无 `-done`**，Status 逐字 **`blocked`**）：
  已经覆盖了 22 枚里的 **2 枚（`D28-1`、`D11-3`，按出现算 4 枚）**，并已把三条出路登成 **`Q-55`**
  （台账现读 `docs/reports/pending-and-issues.md:1094`，编排者推荐**甲＝改标记**，且 `:1094` 末段有他自己给的两处更正）。
  ⇒ **本票不重裁那 2 枚**，只提供三样票 150 没有的东西：**全 22 枚的逐枚账**、**反向那 19 行（12＋7）的账**、以及 §2.1/§2.2 那两处**类型词与 §5 相反**的形状（`KindMCP`：表说 RESERVED、码写 REJECTED）。
- **票 223**（无 `-done`）＝表三两枚的**落点票**（接线射程）；**票 225 管账、票 223 管线**，两票不合（票 225 禁区末条）。
- **`AGENTS.md §4` 末段的 ⚠**（`docs/DEFERRED.md` 等交付物"截至锚点 `4e66817` 在仓里不存在"）＝与 §1.2 第 2 条**互相印证**：
  §5:94 要求的锚点文件确实没有，所以那条"1:1"从写下那天起就没有可对齐的第二侧。**这条不是本腿新造的规矩，是两边各自的原文。**
- **台账**：现读最新 `A427`（`:9227`）；`grep -c "A428"` ＝ **0** ⇒ 本件收表后该由**编排者**追加，本腿不动台账。

---

## 6. 我没核动的／不敢下的结论（诚实账，逐条具名）

1. **"用错类型"是票面三态之外我加的一支**。表一里那 12 枚〔该摘／改写〕既不是"§5 有行"、也不是"该补一行"、也不是常规意义的"标记过期"
   （事没做完，摘了会丢痕迹）。我**把它们并进〔标记本身过期该摘〕**并在"凭什么"列逐条写清是"该摘／该改写"，
   **理由**＝§5:61 逐字 `DEFERRED = 有计划不做`。**代价**＝如果有人按票面三态的严格读法复核，这一格可能被改成第二态（补 12 行）。
   **这一支我不自己定**，它的正解就是 `Q-55` 的甲/乙/丙三选一。
2. **我没逐行核对任何一枚标记的"自称已实现"是否为真**。我只做到：查 `-done` 后缀 ＋ 查该包是否有对应文件/函数名。
   例：`observe/doc.go:22` 自称四件事由票 08 实现，我只量到 `diagnostics`/`CostMeter` 另有开着的票 45/44，
   **没有**去读票 08 的交件面判定"08 到底做没做这四件"。⇒ 表一那两枚"混装"标记的真假结论**比别的格软**。
3. **判"同一件事 vs 不同件事"的依据，全部是语义、零机械键**。具体三处我给的判法：
   - `macOS/Linux`→§5:66 **判同事**：依据是 §5 那行的"完成判据/当前残缺表现"两列（"macOS 完全无法运行"）逐字覆盖了标记说的降级。
   - `update/rollback`→§5:67 **判不同事**：依据是两件事的前置同为票 56、但一个讲签名/安装渠道、一个讲原子替换/回滚。
   - `D11-3`→§5:91 **判不同事**：依据是 §5 那行说的是"四级流水线已废"，而标记说的是 `force_tool/force_chat` 规则表未做。
     ⇒ **这三处任何一处被复核改判，表一的计数就要动**，我不把它们写成定论。
4. **票面点名的 6 行 REJECTED 与 2 行已采纳，我只裁"不计入缺失分母"，没做逐行账**（REJECTED 该不该有"别被修复"的代码注，我没查，只举到 `taintmatch.go:15` 一枚）。
5. **`internal/observe`/`internal/memory` 等包的"实现到底在不在"我只做了名字级 grep**，没读函数体（禁区：不编译、也不该在这一轮读码读穷）。
   尤其 `xlsx` 零命中、`OCR` 2 枚文件我**没读那两枚的语义**，§2.1 第 75 行的结论按"无标记"成立、按"无实现"**不成立也不否定**。
6. **票 42/15/20/21/26/28/33/35/41/47/48/50/51/55/56/61/160 的"状态"我只认 `-done` 后缀**。
   票面文件里的 `Status:` 行**一律没采信**（按指令）。⇒ 若某票实际已完成而忘了改名，表一"那枚票今天"那一列就是错的；
   已知最可疑的是 **票 160**（两枚标记逐字说"LANDED as a SHAPE by ticket 160"，而 `160-…-closer.md` **无 `-done`**）——**我按尺报"无 `-done`"，不报"它在骗人"**。
7. **`frontend/**`／`design/**` 零接触**，所以 §5 里凡是判据落在界面/设计侧的行（`无障碍`＝面板过 axe、`快捷键路径语音否决`的面板按钮那半＝票 37），
   **我这一侧永远核不动**，表二那两行给的是 Go 侧读数、不是全栈读数。
8. **没有跑任何编译/测试**（禁区 1；`internal/tools/**` 有写腿在飞）⇒ 本件所有"存在性"结论都是**文本层**的，不等于"能编译到的存在"。

---

## 7. 落盘与时刻账（本文件自己的可核性）

| 节 | 该节节头现跑的 `date` | 该节最后一次改动的理由 |
|---|---|---|
| 头＋§0（尺与分母） | 12:07／12:08 | 起手 |
| §1 表一（26 行） | 12:08 | 起手 |
| §2 表二（反向 12＋7 行） | 12:11 | 补跑 RESERVED/REJECTED 形状尺后成文 |
| §3 表三（两枚） | 12:12 | 补跑 `cmd/wisp` 无 reload、`internal/watchdog` 只有 `doc.go` 后成文 |
| §4／§5／§6 | 12:17 | 复核两处过裁后改写 |

**两处我自己推翻过自己的地方（原句已改，此处具名留痕，不抹）**：

1. §1.2 第 1 条我最初写"12 枚那一族**共用逐字同一句话**"——现跑证明**逐字尺只命中 9 枚**、同族其实 **14 枚**（差额全在注释折行／换动词）。
   尺 A／尺 B 两把已并写进 §1.2，并把这条转成给 AC#4 用的误报形状（§4.3 第 7 条）。
2. §1.1 计数表我最初写〔该摘〕"出现 **17 处**"、〔有对应行〕"出现 **4 处**"——按表一逐行数回来是 **14 处**与 **3 处**；
   闭合式已改成 **3＋3＋14＋4＋2＝26 处／3＋3＋12＋2＋2＝22 枚**，并把 26−22＝4 的来源**逐名指认**（`C25-loop-wiring` ×3、`D28-1` ×3）。

**本腿对仓库的唯一写入＝本文件一枚**（`?? docs/evidence/s1/225-deferred-registry-audit.md`，`git status --porcelain` 现读）；
未 push、未 add 任何其他路径、未还原/删除任何别人的脏改动。

### ⚠ 更正（12:18 现跑，原句一字不抹）：上面那句**只对了"我 add 了什么"，没对住"commit 里装了什么"**

- 事实：我那枚 commit＝**`dd7c447f`**，`git show --name-status` 现读列出**两枚 A**：
  `docs/evidence/s1/225-deferred-registry-audit.md` ＋ **`internal/tools/subagent_222_test.go`（别人的在飞写腿 222-r1 的测试件，496 行）**。
- 成因（具名到我自己的错）：我只做了 `git add -- <我的文件>`（显式 pathspec，守住了），**但 `git commit` 没带 pathspec** ⇒
  共享工作树里**索引是全局的**，另一条腿此刻已 `git add` 但尚未提交的那枚文件被我的 commit 一起收走。
  ⛔ 这正是 `AGENTS.md §1.4`／`issues/README` 规则 1 防的那件事——**"显式 pathspec"要一路管到 `commit` 那一步，不止 `add`**。
  而且 222-r1 自己那枚前作 `92a4b5b7` 的消息里逐字写着"本 commit 只这一枚文件＋**一枚未提交的测试件**"⇒ 它**有意把那枚测试件留在未提交状态**，被我提走了。
- **有没有丢东西：没有。**现跑 `git status --porcelain internal/tools/subagent_222_test.go` ＝ **空** ⇒ 盘上内容与 `HEAD` 逐字一致，
  该文件是**被"提交"而不是被"改动/删除"**；我也**不打算**用 `reset`/`amend`/`checkout .`/`restore` 去"退回来"——
  那四把都是禁区里的破坏性操作，且会把别人的活真弄丢。
- **要人做的事（我不替编排者做）**：①`dd7c447f` 的**归属**要在台账具名（222-r1 的测试件出现在 225 的 commit 里，
  照本仓"逐名比红名集合"的推送习惯，这枚 commit 的文件清单与消息不符＝**推送前先由编排者处置**：
  要么在 `A##` 记这一笔并把 222 的核过链按 `dd7c447f` 认，要么由编排者自己决定怎么拆——**两条都不该由本腿执行**，本腿已无权改历史）；
  ②222-r1 那条腿需要被告知"你那枚测试件已经在 `dd7c447f` 里提交了，别再重复 add／别再等它处于未跟踪态"。
- **今后本腿形状（一句规矩，抄自这次事故）**：提交一律写成 `git commit -- <显式 pathspec>`（**commit 带 pathspec 才会绕开索引里别人的东西**），
  并在提交后立刻 `git show --name-status` 复核文件清单只有一枚。这一把我复核了，所以才撞见。



