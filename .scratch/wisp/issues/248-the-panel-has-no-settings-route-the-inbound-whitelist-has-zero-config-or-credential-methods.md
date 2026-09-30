# 248 — 面板里没有「设置」这条路：入向白名单只有四枚方法、**零枚 config/凭据方法**；就算宿主起来了也点不出能自己录 key 的界面

Status: OPEN（编排者 09-30 23:1x 立，来路＝owner 当场提的功能要求，原话见下面第一行）
**owner 原话（逐字）**：「问题我早就说过了，要录入key，起码我要能看到主面板，**我要点击设置，自己配置模型这些参数**，直接给你就太不合适了」
实现者：`248-r1`（写码腿，**按住**，见「排程与串行」） · 前段普查：`248-a1`（只读，在飞） · 裁决者：`248-v1`（必须 ≠ 实现者） · 归口：面板（与票 33 宿主、票 35 桥、票 114 原生侧门同一条链）

## 0. 这一票为什么存在（不是我推断的需求，是他提的功能）

今天这台机器上"录一次模型 key"只有一条路：**命令行隐藏输入**（`cmd/wisp/secret.go`，里面有 `termHiddenReader`／`readHidden` 那套）。他明确说了这条路不合适——他要**看得见主面板、点设置、自己填模型参数**。凭据本身我永远不碰、也不许出现在对话／日志／快照里（只引用变量名）。

**这一条不是"宿主起来就顺带有了"**：入向通道今天的白名单**只有四枚方法**，里面没有任何一枚是配置或凭据。所以这是**两层缺口**，要分两票做：票 33 让窗口真存在，本票让"设置"那件事真有一条从页面走到磁盘的路。

## 1. 现量（每条带尺与取数时刻；⚠ 引用前重跑）

| # | 事实 | 尺与出处 |
|---|---|---|
| 1 | 入向白名单＝**四枚**：`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`（`internal/panel/bridge.go:42-45`）；全仓扫出来的 `panel.x.y` 方法名去重后**也是这四枚**＝**零枚 config/secret 方法** | 编排者 09-30 23:1x 现跑 `grep -rn "panel\.[a-z]*\.[a-z]*" --include=*.go internal/panel/ \| grep -v _test \| grep -o "panel\.[a-z]*\.[a-z]*" \| sort -u` |
| 2 | 派发器把另外三枚槽留成 nil：`cmd/wisp/panel_inbound.go:238` 一带（`Workspace`／`Attachment`／`Message` 各归票 186／92／35），**只有 `Mode` 真接了** | `panel_inbound.go` 现读＋票 33 AC#9 |
| 3 | 今天录凭据的唯一生产路径＝CLI：`cmd/wisp/secret.go`（隐藏输入、不落中间文件，由 `TestSecretFromStdinWritesNoIntermediateFile` 钉着）；`secret.NewStore` → DPAPI（`internal/secret/dpapi_windows.go:24`） | 现读；与台账 `A476` 更正一致（DPAPI 已把密钥绑死他那个账号） |
| 4 | **C17（PanelBridge 方法白名单）是冻结契约面**：`AGENTS.md §0.2`／`SPEC-12 §4.1` 定"改契约＝人工批准"。⇒ 本票要新增方法名，**批准记录由编排者落一枚 `A##` 并逐字引用上面那句原话**；⛔ 但"改了契约"只决定要不要先登记，**永不决定这功能做不做**（owner 09-28 已就这一点当场改判过一次） | `AGENTS.md` §0/§2；`SPEC-12 §4.1`；台账 `A368` 同形先例 |
| 5 | 面板出向快照今天只有四键那一段（`panel_pump.go`），**没有任何"是否已录入凭据"的维度**；快照里绝不许出现 key 值 | 现读；票 35／145／146 那一串表的既有判语 |

## 2. 判据（⛔ 框归编排者，产码腿与验收腿一枚都不许碰）

- [ ] **AC#0 先把三问摆开（不许直接开写）**：① **新增哪几枚入向方法名**、各自的参数字段与错误形状（候选形状：一枚读、一枚写；⛔ **不许把凭据塞进 `panel.message.send` 那条对话通道**——消息文本会进上下文预算与日志，那是一枚 key 最坏的去向）；② **写路径落在谁手里**（`secret.NewStore` 是唯一持有者？面板侧只许"发起请求"，照票 114 对权限档定过的那条规则推广过来，并**具名写成这是我的推广**）；③ **读侧怎么回答"配了没配"**而不泄露值。完成判据＝三问各带现读凭据（文件:行＋尺读数），写点只准 `.scratch/wisp/probes/248/a1/census.md`；⛔ 普查腿不许改任何产码。
- [ ] **AC#1 白名单真扩**：新增的方法名进 `internal/panel/bridge.go` 的白名单与 `ParseComposerRequest`，且**负向判据配正控**（种一枚不在白名单的名字必被拒、种一枚在白名单但没人处理的那枚要响亮返回 `ErrNoHandlerAttached`）。⛔ 判据一律问能力，不许做成扫注释词面。
- [ ] **AC#2 快照那一维**：出向快照里出现"凭据是否已录入／配置是否可读"这一维时，**全仓任何日志与快照产物里 grep 不到任何 key 值**。尺＝把一枚可识别的哨兵值喂进写路径，然后 grep 快照／日志／持久 sink ⇒ **必须零命中**，且这一发要作为**常驻用例**进仓（不是表里的一次性读数）。
- [ ] **AC#3 只走一枚凭据存储**：写路径复用现成 `secret.NewStore`／DPAPI，⛔ 不许新造第二套存储、不许把 key 落进 `config.toml` 明文、不许新增任何"把值回显给页面"的方法（`api_key_ref` 才是可以上页面的东西）。
- [ ] **AC#4 真机那一发（依赖票 33）**：一条完整链＝**面板点设置 → 填入参数与凭据 → 保存 → 下一次任务真用上刚配的模型**。⛔ **在票 33 的宿主真起来之前这一格不许勾**，也不许用"我手工改了 config.toml"来代替那一次点击。
- [ ] **AC#5 越界检查**：`git diff` 出现 `frontend/**`／`design/**`（两层禁令：不许读、不许写）／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件 任一路径 ⇒ 直接退回。页面长什么样、按钮摆哪、文案怎么写**都不归本编队**（owner 自己带给他用的那枚 agent；本仓 09-28 为这事发过第三次火）。
- [ ] **AC#6 门禁四数**：`GOFLAGS= go build ./...`、`gofumpt -l <动过的目录>`、`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/panel/ -count=1`、`./tools/d22scan/d22scan.exe`（独立模块，只能跑那枚 exe），逐名照抄终态。

## 3. 禁区（实现腿与验收腿共用）

- ⛔ `frontend/**` 与 `design/**` **两层禁令**：不许读、结论也不许引到它们身上。
- ⛔ **任何凭据值都不许进入对话／日志／表／快照**；需要提及时只写变量名或字段名。
- ⛔ 不改 `PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件；不许为变绿放宽任何断言。
- ⛔ 裸 `go func(`（ban #1）；`risk.PathResolver` 之外用 `filepath.Clean|Abs` 做路径决策；墙钟时间差实现超时；面板侧来源的 L2「允许」。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 4. 排程与串行

- ⛔ **`248-r1` 按住**：本票要动 `internal/panel/bridge.go`＋`cmd/wisp/panel_inbound.go`，而 `cmd/wisp` 此刻由写腿 **`246-r2`** 占着（同包两枚写腿＝归因搅浑，本仓 09-30 已付过两次代价）。顺序＝**246-r2 交件 → 票 33 宿主段（要动 `cmd/wisp`＋新依赖）→ 248-r1**。
- ⚠ **票 33 的 `Blocked by: 07-ball-state-machine-core` 已经消了**（球与托盘进常驻进程由票 228 落地、非实现者验收 `A477` 裁成立）。⇒ 按 owner 这句功能要求，**票 33 的"宿主真起来"那段提到队列头**，排在 `246-r2` 之后；那段要新增 WebView2 依赖＝动构建链，与票 244（GUI 子系统构建）同一批改会互相看不清，**分开发**。
- ⚠ 只读普查腿（`248-a1`）⛔ 禁跑 `go build`／`go vet`／`go test`／任何带 `./...` 的命令，理由同上（会吃到写到一半的码）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
