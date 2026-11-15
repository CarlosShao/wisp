# 214 — **加号菜单第二类："往这次对话里加东西"**（选文件／附件／@引用）——附件那台机器**已经建好但一根线都没接**（410 行、生产零调用者）

- Status: **待派（owner 09-28 22:1x 口头接单，原话见票 213 第一节）**。这一票是**接线＋补最小缺口**，不是从零做——**这是今天查出来的最便宜的一枚**。
- 与票 213 的关系：213 管"命令"，214 管"加内容"，215 管"技能/插件/MCP"。三票共用那一根**入向线**（前置见 §禁区）。

> ⚠ **09-29 12:1x 补一格（腿 `213-c1`＋编排者自己 `grep` 复核）**：**配置层今天零枚附件键**——`grep -rn -i attachment internal/config/`（剥测试）**空结果**。⇒ 本票上面那句"四枚键早就在"指的是**快照字段**（`internal/panel/composer.go:238-241`），**不是 `config.toml` 里能调的东西**；`maxAttachmentBytes` 与 `acceptedAttachmentMimes` 今天**来自代码里的默认**（`AcceptedMIMETypes()`），**没有任何用户可以调的入口**。
> ⇒ **落点定案（我自己定，不摆成待答项）**：本票**只用现成默认**，接线那一步**不新增配置键**——加键属 `D36` 那棵配置树（契约面），要动先落 `A##`。**"用户能不能自己改上限／白名单"这一格本票明确不做**，在此具名留给后续票。⚠ 连带一条别家现成形状：**"限制是写死的"这件事界面上要说得出**（三家都把"这次生效的是什么"画出来了，见 `docs/reports/survey-2026-09-29-command-catalog-across-harnesses.md` ③ 末行；我们 `instructions_200.go:26-40` 那五枚态＋`:59 Dropped` 就是同一格的现成写法）。

## 现量（起手逐条复算）

> ⚠ **09-28 22:3x 更正本票起手两格**（我原来写错了，是调研腿 `survey-plusmenu` 顶回来、我自己复跑确认的）：
> ① **快照里"附件那四枚键早就有了"**——`internal/panel/composer.go:238-241` 就是 `attachments`／`acceptedAttachmentMimes`／`maxAttachmentBytes`／`attachmentError`，
>    所以本票**不加任何新键**、也就**不会打红那两把"加字段必两侧同批移动"的尺**（我原来在 §交付 2 里写的"要能进快照"是**接读者**，不是加字段）；
> ② 真正断的那一根线在 `internal/panel/pump.go:226`：`NewComposerState(mode, workspace, nil, maxAttachment)`——**附件那枚实参被写死成 `nil`**。
>    ⇒ 缺的是**读者与实参**，字段与键都在。落点若写成"加字段"就是白改。

| 事实 | 读数 | 尺 |
|---|---|---|
| 附件受理器已存在 | **410 行**：`AttachmentBroker.Ingest`、MIME 嗅探 `matchesISOBaseMedia`、白名单 `AcceptedMIMETypes`、大小上限、文件名守卫 `NameGuard`、`ArtifactSink` | `internal/panel/attachments.go:59-282` |
| **谁调用它** | **生产零调用者**（只有自己的测试件） | `grep -rn --include=*.go "panel\.Attachment\|Attachments(" internal/ cmd/ \| grep -v internal/panel/attachments` ⇒ 空 |
| 快照侧的附件四枚键 | **已在**（值恒缺：`pump.go:226` 传 `nil`） | `internal/panel/composer.go:238-241`＋`internal/panel/pump.go:226` |
| 落点在哪 | 有 `ArtifactSink` 接口，但生产装配根有没有真给它一个 sink **未量** | 起手现读 `cmd/wisp/run.go` 的装配段 |
| 选文件/选目录 | 别家做成两枚分包（浏览版＋原生版） | `packages/client/ui-directory-picker-browse`／`-native`（见 `docs/reports/survey-2026-09-28-dsh-ui-packages.md`） |
| 我们这边的路径权威 | 一切文件系统决策必须走 `risk.PathResolver`（在它外面用 `filepath.Clean/Abs` 是 `AGENTS.md` §1.2 硬禁） | `PLAN.md:1288` 一带 |

## 这一票要交付什么

1. **把受理器接上**：装配根给它一枚真 sink＋一条真入口（界面/命令行都行，但**要具名写是哪一种**）。
   ⚠ 接的时候三件不许省：**大小上限从配置读**（不许写死）、**文件名守卫要过 C26 路径判定**（不许自己 `filepath` 拼）、**落盘位置必须是它该落的那棵树**（票 174 那一族：artifacts 根不在授权根里就是洞）。
2. **加进去的东西要能被"看见"**：这次送了什么进去（文件名、类型、字节数、有没有被拒、为什么拒）要能进快照——**这是 owner 截图里"添加上下文/上传附件"那一格的真内容**。
   ⚠ 不许只报"成功"：被拒的附件要带**人话原因**（`A408` 那一形：文案在、控件不在＝反面教材；反过来"控件在、状态永远成功"也是撒谎）。
3. **选文件这一格**：桌面版走系统选择器还是自绘列表，**由界面那支定**；本票只保证 Go 侧拿到路径之后**判定链完整**（授权根、改写账户、来源戳）。
4. **@引用**：引用工作区内文件/目录 → 走同一台受理器，**不许开第二条读文件的路**。

## 判据

- [x] **AC#1**：受理器有**生产调用者**（现跑 `grep` 点数：改前 0 枚 ⇒ 改后 ≥1 枚，且那枚调用者在 `cmd/wisp` 的装配路径上，不在测试里）。
- [x] **AC#2**：投一枚超限文件 ⇒ 拒、**且快照里能看到"被拒＋为什么"**；正控＝同文件放到限内 ⇒ 收。
- [x] **AC#3**：投一枚伪装文件（扩展名说自己是图片、字节不是）⇒ 按嗅探结果判，**不按扩展名判**（这条尺要现读 `matchesISOBaseMedia` 的语义再写死）。
- [x] **AC#4**：落盘位置在授权根内 ⇒ 越界那一形**必红**（票 174/175 那族判据同族）。
- [x] **AC#5**：来源戳（C25 现成名）盖上 ⇒ 附件内容进上下文时被当作**不可信外部文本**处理（票 200 那节"说明文件≠授权"同一条纪律）。
- [x] **AC#6**：三件"没接之前不许说已具备"——界面点不了、没有选择器、没有 @引用，**逐条具名留在票上**。

## 禁区

不新造第二条文件读入口（一切走现成受理器＋`risk.PathResolver`）；不在 `risk.PathResolver` 外做路径决策；不新增 C17 方法名；
`frontend/**`、`design/**` 零写；`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/**` 一字不动；
只 commit 不 push、显式 pathspec、禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；仓内不删东西。
— 2026-10-09 12:0x 编排者腿 boxes-1 搬框：本节原为散句 bullet，现改成 - [ ] 形以进全勾尺；⛔ 判据文字一字未改、⛔ 未翻任何框。

## Progress log

- 2026-11-07 17:5x 腿 214-attach-wire：接线落地——`ComposerDispatch.Attachment` 在 `cmd/wisp/panel_inbound.go` 挂上真受理器（`panel_attachment_214.go`：现成 broker＋真 memory.Store artifacts 树、懒开库、上限用编译默认零新增配置键）；`pump.go` 加 `Attachments`/`AttachmentReason` 两枚 reader＋run.go 装配根两枚 method value＋resident boot 绑 ledger（票 186 holder 同款）；145 钉子按新事实改名改断言为 `TestComposerAttachmentsTravelFromTheIngestDoorToTheWire`（未删断言，理由见 docs/evidence/s1/214-attachment-wiring.md）。AC#1 census 0→2；cmd/wisp 门禁子集全绿；internal/panel 余 2 红＝既有红（契约 files 键 page 侧未声明／C21 被 ball 在飞件牵红），归属与读数见证据件 §3。AC 勾未翻。
- 2026-11-07 18:5x 腿 214-ac5-stamp-feasibility（只读取证，零产码）：AC#5 盖戳三选一裁定＝**乙**——C25 盖戳宿主＝消息门把附件文本喂进 agent 的那一跳；现量 `OutgoingMessage.ForAgent` 生产调用者 0、`MessageRequestHandler` 产码实现 0、`panel_inbound.go:352` Message 仍 nil（按名拒绝）；且票 35 现有 AC 八枚框全是传输层、无一格认领 handler 本体 ⇒ 工单池今天没有 owner，须编排者立票。丙排除（AC#5 逐字在票面 ：41，契约没记错）；甲排除三理由（对话入口超射程／Mark 是 scope 键控、ingest 时无 task scope 可挂／本波 cmd/wisp＝188·internal/panel＝187）。正面凭据：盖戳引擎已在装配根（`run.go:330 c25`）＋附件读者与 ForAgent 渲染 214 已备好。全链引证与未来落点见 docs/evidence/s1/214-ac5-stamp-feasibility.md。门禁：零 Go 改动，无对象可扫。
---

## AC#1 收件（2026-11-15 13:5x，`b6bc1c61b`）：**满足**，翻

**尺1**：`b6bc1c61b` 在，7 件，**产码/票面/台账命中 0**（证据件＋probes）。
**尺2**：`docs/evidence/s1/214-ac1-receipt.md` **19830 B／226 行**；载具 `.scratch/wisp/probes/214/ac1/`（变异副本／overlay／三份日志／census）。
**尺3**：★这一格我出的错和纠正都记下来 ——
- **生产调用点我自己 grep 复核**：`Attachment: attachment` 产码**唯一**在 `cmd/wisp/panel_inbound.go:386`；`NewAttachmentBroker` 产码 3 处（`internal/panel/attachments.go:175/179` 定义侧 ＋ `cmd/wisp/panel_attachment_214.go:204` 真装配）；`newAttachmentIngestHandler214` 产码 4 处含定义与调用。
- ★**我先读到"改前不是 0"并准备退回这一格 —— 是我错**：我那条核命令写成 `$parent+':cmd/wisp/panel_inbound.go'`，而 `$parent` 是 `@(...)` 出来的**数组**，PowerShell 把 `+` 当**数组追加**，于是 `git show` 收到两个参数，读出来的根本不是父提交那棵树。**正确读法**：父提交 `25eb216` 里是 **`Attachment: nil,`（`:317`）**，注释逐字还写着「attachment = ticket 92, message = ticket 35, both still refused by name」；`fc76d5d7e` 的 diff 正是 `nil → attachment`。⇒ **改前 0 成立，腿对**。

### 口径交代（腿主动区分，我认）
实现件写 `0→2` 是**票面原尺**（只数受理器本体两处）；收件件给 `N=4`（含装配两站）。**两个数都对，尺不同**。⇒ 票面按原尺的口径看是 `0→2`，收件按"生产调用点"看得 4 处。⛔ 我不用一个数去否另一个。

### 会红的测试（这是 A1875 的第二件，有了）
`-overlay` 把 `:386` 的槽换回 nil（**盘上原件一字节没动**，md5 前后同为 `53243FD2…`，`git status --porcelain -- cmd/wisp` 空）⇒ `--- FAIL: Test214AttachmentAddEndToEndStoresBytesInArtifacts`，红句「该方法在名册内，但本机未接入处理器」；不摘接线同一枚 `PASS`。组跑 `-run '^Test214'` ＝ **4 红／4 绿**（超限那枚对 nil 门无鉴别力，腿**自己具名不用它当正控**）。

### 其余格（⛔ 未翻，照实）
AC#2–#6 逐格写「未裁」并具名交回：AC#2 的快照半边要另立收件腿／AC#5 已知未完成指向新票／AC#6 三件仍在。
另登记一枚**过期注释** `internal/panel/composer_dispatch.go:165-166`「Untouched by this slice」—— 腿**没拿它当证据**，也没改（对，注释不是证据）。

### 〔未跑包〕如实
`go test ./cmd/wisp` 全包未跑（该包归写腿 452），只跑 `-run` 单枚／单组，跑前 `git status --porcelain -- cmd/wisp` 为空。`go vet`/`gofumpt`/`d22scan` 未跑：本腿零 Go 改动，无对象可扫。
---

## AC#2–#6 收件（2026-11-15 15:4x，`f563a3c0e`）：**满足 3／部分 1／不满足 1**

**尺1**：`f563a3c0e` 在；件 `docs/evidence/s1/214-ac2-ac6-receipt.md` **35326 B／266 行**。
**尺2**：七发表、锚点日志、`residue.log`（五枚原件 md5）齐；`git status --porcelain --untracked-files=no -- cmd/wisp internal/panel` 空。
**尺3 我核的两件**：
- 它说那枚装配测试**读盘上源码字符串** ⇒ 我读了 `panel_attachment_214_test.go:284`：确实 `runtime.Caller(0)` ＋ `os.ReadFile("run.go")` ＋ `regexp.MustCompile("Attachments:\\s+rt\\.attachmentDraftRows,")` ⇒ **overlay 下不响，零牙成立**
- `OutgoingMessage{` 产码构造 **0**（我现跑）；`ForAgent(` 产码只有 `internal/panel/composer.go:495` 一处，**是定义不是调用者** ⇒ 腿的"调用者 0"对

### 逐格
| 格 | 判语 | 生产点 → 会红测试 | 牙种 |
|---|---|---|---|
| AC#2 | **满足** | 上限 `attachments.go:214-217`→1 红／入账 `panel_attachment_214.go:218`→**4 红**／理由 hop `pump.go:432-434`→2 红 | 行为牙 |
| AC#3 | **满足** | 嗅探 `attachments.go:240-251`、声明一致性 `:252-255` → 改成"信声明" **2 红** | 行为牙 |
| AC#4 | **满足（见下裁决）** | 裸名守卫 `attachments.go:207`→1 红；落点树有行为断言（`<dir>/artifacts` 恰 1 枚、内容哈希名） | 行为牙（拒的一半） |
| AC#5 | **满足（见下裁决）** | 移交**新票 453** | — |
| AC#6 | **满足** | 文件类判据，无测试牙；A1881 在锚点 `fc76d5d7e` 上三件在**产品页**都真 | — |

### ★我的两条裁决（⛔ 不是腿让我拍，是我该拍）
**① AC#4 按票面原意读 ⇒ 满足。** 票面问的是"落点是不是真的那棵 artifacts 树"，答案是**是**（名字守卫、`O_EXCL`、私有 DACL、递归配额都在，且恰 1 枚、内容哈希命名）。
腿指出的"C26 授权根不成立"是**另一个洞**：授权根只从 `cfg.FS.AllowedDirs` 来（`cmd/wisp/run.go:772-774`，默认空），artifacts 根从不入册 ⇒ **与票 174 同一个洞**。⛔ 把它并进 AC#4 会变成"214 要等 174 定案"，那是范围蔓延。⇒ **AC#4 与本洞分开记**：洞归 174（具名），本格按原意翻。

**② AC#5 立票移交 ⇒ 满足。** 已立 **票 453**（`453-the-attachment-goes-in-but-there-is-no-exit-…`），把"消息门 handler 本体 ＋ C25 附件源"接住。⛔ 不用散文交接（散文不可复核）。

### ⚠ 两件必须留在账上的新发现
1. **装配读口零牙**：`-overlay` 摘 `cmd/wisp/run.go:1197-1198` ⇒ `-run '^Test(214|65)'` **0 红**；摘 `panel_resident_pump.go:101-102` ⇒ **也 0 红**。⇒ **快照能力活着，但接线可被无声删除。**
2. **AC#6 票面两件已过期**：今天 `frontend/src/components/composer-row.tsx:291/:299-306` 有真选择器＋附件按钮 → `frontend/src/lib/panel.ts:718` → `panel.attachment.add`，构建产物（2026-11-11）含该路由；**只剩 @引用仍无**（`ComposerRow` 不传 `sources`）。⛔ 别拿"界面点不了"当今天现状。
3. **`.mov` 误拒**：诚实 QuickTime（`qt  `）声明 `video/quicktime` 被拒，理由"文件头表明它是 video/mp4" ⇒ 已立**票 454**。

### 本票 5/5 结案
