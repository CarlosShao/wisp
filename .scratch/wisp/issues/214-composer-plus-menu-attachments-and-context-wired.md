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

- **AC#1**：受理器有**生产调用者**（现跑 `grep` 点数：改前 0 枚 ⇒ 改后 ≥1 枚，且那枚调用者在 `cmd/wisp` 的装配路径上，不在测试里）。
- **AC#2**：投一枚超限文件 ⇒ 拒、**且快照里能看到"被拒＋为什么"**；正控＝同文件放到限内 ⇒ 收。
- **AC#3**：投一枚伪装文件（扩展名说自己是图片、字节不是）⇒ 按嗅探结果判，**不按扩展名判**（这条尺要现读 `matchesISOBaseMedia` 的语义再写死）。
- **AC#4**：落盘位置在授权根内 ⇒ 越界那一形**必红**（票 174/175 那族判据同族）。
- **AC#5**：来源戳（C25 现成名）盖上 ⇒ 附件内容进上下文时被当作**不可信外部文本**处理（票 200 那节"说明文件≠授权"同一条纪律）。
- **AC#6**：三件"没接之前不许说已具备"——界面点不了、没有选择器、没有 @引用，**逐条具名留在票上**。

## 禁区

不新造第二条文件读入口（一切走现成受理器＋`risk.PathResolver`）；不在 `risk.PathResolver` 外做路径决策；不新增 C17 方法名；
`frontend/**`、`design/**` 零写；`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/**` 一字不动；
只 commit 不 push、显式 pathspec、禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；仓内不删东西。
