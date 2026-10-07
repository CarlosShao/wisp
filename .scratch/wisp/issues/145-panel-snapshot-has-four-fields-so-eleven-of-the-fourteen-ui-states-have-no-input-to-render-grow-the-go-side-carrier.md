# 145 — 面板快照只有四个字段，导致十四种界面状态里十一种"没有输入可画"：把 Go 侧的 `Snapshot` 扩成真载体

- Status: ready-for-agent（**未派**，见末节"为什么现在按住"）
  > **Status 追加（09-26 10:5x，编排者；上面那行原样保留不删）**：**已派、已交件、已按非实现者验收表收**＝`review`→**部分结案**。
  > 档位＝**附条件（六枚全为记录级、无一枚要求改码）**；勾＝**AC#1/AC#3/AC#4/AC#5 已勾，AC#2/AC#6 维持未勾**（**未勾≠没干活**：AC#2 的量出来的答案是"落地集＝空"，见下面"收表"那节）。
  > `-done` **未加**——加它的前提是 AC#2/AC#6 有一枚可判的归宿（现挂在 `Q-51` 与我的 R1/R2 上）。
- 来源：前端会话 `docs/reports/frontend-session-log.md` §47（13:1x 十四态逐行盘点）＋编排者 13:3x 独立复量
- 关联：票 35（那句 "ticket 35 owns the pump" 就是本票的泵）、票 92（`Snapshot` 从 `approval_test.go` 搬进生产码的那一程）、`Q-50`（甲已落地，"当前屏"这一维今天仍无处可传）、`A241②`、`A234`、`task #79`

## 现量的形状（一行不省，两条命令可复算）

- `internal/panel/composer.go:43-49` 的 `Snapshot` **恰四字段**：`pending` / `results` / `composer` / `generatedAt`；TS 侧 `frontend/src/lib/panel.ts:129-137` **逐键对齐**。
- `PLAN.md:3473-3488` 那张表列出**十四种**界面状态。逐行判下来（前端会话 §47.1 给了逐行依据，我复量了标 ✅ 的三条）：**已有 1（行 7 L2 卡）／部分 2（行 2 流式、行 14 注入检出）／矛盾 3（行 6、8、10 里各一条）／缺 8**。
- ⚠ **"缺"的主因不是组件没写**：`loading-state.tsx`、`thinking.tsx`、`task-rows.tsx`、`stream-text.tsx` 这些**已在盘上**（`frontend/src/components` 共 17 枚），但**它们的输入在这份快照里根本不存在**——token 计数、工具调用条目与四值状态、推理过程、错误分类、重复计数、成本（token/金额）、取消原因、当前屏，**一枚字段都没有**。
- ⇒ 前端今天的处理是 `App.tsx:52` 的 `UnfedScreen`：`:82` 除 chat/approval 外**每屏渲染一句人话**，**不画假内容**（守 owner 的 P9"不得拿假件当真实字段"）。⇒ 所以"页面看起来没进展"是**这一格的下游症状**，不是界面层的进度。

## AC

- [x] **AC#1 先把"哪一行缺哪个字段"落成一张对照表**，不许笼统说"扩快照"。判据：`PLAN.md:3473-3488` **十四行逐行**给出ⓐ需要的字段名与类型（Go 侧）、ⓐᐟ它的**真来源**（哪个既有包已经算得出来：`observe`／`tools`／`risk`／`llmrecord`／`memory`…）、ⓑ若今天算不出来就**明写"无源"**并给归口（另开票／DEFERRED 五字段）。⇒ **ⓑ 那一支是本票的硬约束：宁缺毋造。** 一个字段如果只能靠填常量才"通"，它就不许进本票的落地集。
- [ ] **AC#2 落地有源那一集**：`Snapshot` 加字段＋构造点填真值（构造点是谁、几处，AC#1 的表里必须点名）；TS 侧 interface 的对齐**不写 `frontend/**`**，只在证据件里给出"应当长这样"的逐键清单，由前端会话自己落（归属见 `A102` 与项目账）。
- [x] **AC#3 "当前屏"这一维要一次定死两种角色**：ⓐGo→面板（快照里带"现在该显示哪一屏"）ⓑ面板→Go（`Q-50` 甲刚删掉的那条请求）。本票只做ⓐ；ⓑ**留给泵（票 35）**，并在 AC#1 的表里写明"ⓐ 单独存在时界面能不能真换屏"——如果不能，就写"这一格落地后界面仍点不动，缺的是 ⓑ"，**不许让下一位读成"换屏已通"**。
- [x] **AC#4 契约轴（本票最容易越的面）**：`Snapshot` 是 C17／D35 那一族的对外形状，⚠ **加字段＝扩契约**，所以：①`docs/PLAN.md`、`docs/specs/**` 一字不动；②`thresholds.go`／golden／`rules_gateway.go`／`allowlist.txt`／`internal/risk/**` 零命中；③**TS 侧对齐之前，Go 先加字段不会弄红任何门**——要用一枚用例证明"Go 加了字段而 `frontend/src/lib/panel.ts` 没跟上时，谁会响"（这把尺今天是否存在？不存在就登记，别假装有）。④`panel.NewApprovalCardView` 与渲染那行不许改。
- [x] **AC#5 门禁与名册**：`go test ./internal/panel/ ./cmd/wisp/` 改前改后各一次（四数之外**点名册差集**，⚠ `cmd/wisp` 不带 sherpa DLL 时改前就是红的——那是本机环境不是被验物，基线一律用 `scripts/wisp-cli-tests.sh` 或带 DLL 取，并在证据件写清用了哪支）；`gofmt -l` 空；`go vet ./internal/panel/ ./cmd/wisp/` 空；`sh scripts/d22scan.sh` rc=0（独立 module，别在根目录 `go vet ./tools/d22scan/`）。
- [ ] **AC#6 反向判据（防"加了个没人读的字段"）**：对每一枚新字段答一句——**哪一枚用例断言了它的值来自真来源**？答不出的字段＝装饰品，从落地集里拿掉并在表里写明。

## 为什么现在按住（具名，不写"等一等"）

`internal/panel/**` 此刻被**非实现者复判程**在读（门钉 r3 的对抗验收，表 `panel-l2-grant-nail-fix-r3-accept-r1.md`，它全程断言 `git status --porcelain -- internal/panel/` 为空）⇒ 我此刻往那枚目录写代码，会让它的"树干净"断言出现**假污染**读数。⇒ 放行条件＝那程交件；本票的 AC#1（只读盘点＋对照表）**可以现在就开**，因为它不写 `internal/panel/**`。

地界与规矩照 `README`（显式 pathspec、禁 `add -A`/`--amend`/`reset`/`stash`/`clean`、禁 push、临时件只建不删、变异落仓外副本、台账与勾归编排者、前提不成立就报回来）。

> **AC#1 交件后对本票票面的六处更正**（`docs/evidence/s1/145-snapshot-field-census-r1.md`，锚点 `88eab34`，只读程、零生产码改动。原句一律留着不抹）
>
> **① 标题里那枚"十一种"是我抄来的、且站不住。** 普查现量：面板侧**完全无输入**的是 **8 行**（1/3/4/5/10/11/12/13），算上行 8 的半枚输入＝**9**，再算上行 9（原生侧、面板契约根本够不着）＝**10**。**要凑到 11 必须把行 14 也算成"无输入"——而行 14 恰恰是 fed 的**（见 ④）。⇒ 下一位引这个数字请用 **8**，并把 ⑤ 那枚更硬的阻塞一起读。
>
> **② 我引的行号漂了一格**：`Snapshot` 的四字段是 `internal/panel/composer.go:44-49`（**不是 `:43-49`**），"ticket 35 owns the pump" 那句在 `:41`（不是 `:39-40`）。`sed -n '38,52p'` 可复算。
>
> **③ 两处"存在性"我写错了**：ⓐ 我在候选真来源里点了 `internal/llmrecord` —— **这枚包不存在**，仓里只有 `cmd/llmrecord/main.go`；ⓑ "`frontend/src/components` 共 17 枚组件"只在**递归**数法下成立（顶层 6 ＋ `ai-native/` 8 ＋ `ui/` 3）。⇒ 不抹，按现量收窄。
>
> **④ 前端 §47.2.3 那条"注入检出的来源送不到面板"是**反方向的错**（比真相悲观）。** 那条链今天**是通的**：`internal/risk/rules_gateway.go:112` → `internal/panel/approval.go:83`（`Reason` 原样带过来）→ `l2-approval-card.tsx:159`、挂在 `:208`。真正的缺陷是 `:216-220` 那**多出来的一句硬编码文案**。⇒ **这一格是前端自己能动的，不是在等 Go 加字段**，别拿它当"快照太瘦"的证据。
>
> **⑤ 阻塞比票面写的更强，而这一条我已独立复量过**（不只信普查）：`Snapshot` 不只是"字段少"——**它从来没有生产代码构造过**。现量（`grep -rn --include=*.go` 剔 `_test.go` 与函数定义本身）：`NewSnapshot(` **0 枚调用者**、`NewComposerState(` **0 枚**、`NewApprovalCardView(` 唯一一枚生产调用者是 `cmd/wisp/panel_assets.go:68` 那枚**命令行诊断分支**。⇒ **AC#2 的验收判据必须加一句**：不许用"字段加了、构造点在**测试**里填真值"结掉（`composer.go:41-43` 那行注释自己就在警告这种"只有测试能构造的视图模型，生产从不发"）。普查为此新立第四态 **「无生产者」**（源与导出的读法都在、没有任何东西去读它们），三态表达不了这一形。
>
> **⑥ AC#2 与 AC#4③ 结构上打架，这一格要我裁、普查程没有擅自选。** AC#4③ 问的那把尺**存在**：`internal/panel/composer_test.go:48`，双向对账在 `:73-78`；而它的方向是 **Go 先加字段、`frontend/src/lib/panel.ts` 没跟上 ⇒ 那把尺会红**。这正好撞 AC#2 自己那句"对不齐不写 `frontend/**`"。⇒ 三种读法列在普查件 §4，**编排者裁决见下方"裁决"节**（另：这把尺读的是**工作树**，所以一枚脏的 `panel.ts` 能让它对着别人的未提交改动报绿——这条也留给验收方打）。
>
> **⑦ 本节上面"为什么现在按住"那格的放行条件已满足**：门钉 r3 的复判程已交件（`panel-l2-grant-nail-fix-r3-accept-r1.md`），`internal/panel/**` 不再被复判程读。⇒ 本票 AC#2 的写码腿**可以派**，但必须先解 ⑥。

## 编排者 14:3x 裁决（对 ⑥，并把本票的落地次序重排）

**先说结论：AC#2 按住，本票拆成三段，只有一段的判据已经齐了。**

- **甲段（现在就能派，零扩契约）**：把 `Snapshot` **现有的那四枚字段**真送进面板——也就是票 35 那根泵。理由＝⑤：`NewSnapshot(` 生产码**零枚调用者**（我 14:3x 独立复量，不只信普查），所以今天加第十二枚字段仍然**没有任何地方能给它真值**，AC#6 那一问（"哪枚用例断言它的值来自真来源"）只能答"某枚测试里手搓的"⇒ **那正是本票自己要防的装饰品**。泵先落地，后面每一枚新字段才有一个能填真值的构造点。
- **乙段（不需要 Go，也不动 `frontend/**` 之外任何一处）**：④ 那格——行 14 的来源链今天**是通的**，缺陷只是 `l2-approval-card.tsx:216-220` 多出的那句硬编码。归前端会话自己，已并进它的 §10 队列。
- **丙段（＝原 AC#2 扩字段，按住，卡在 ⑥ 上）**：`composer_test.go:73-78` 那把尺是**双向**的 ⇒ "Go 先加、TS 后跟"与"TS 先加、Go 后跟"**两个顺序都会红**，而"不许放宽断言换绿"是本仓硬规矩。⇒ 剩下的唯一非放水解法＝**同一枚 commit 同时含 `internal/panel/**` 与 `frontend/src/lib/panel.ts`**。这件事**我不裁决**：谁有权写 `frontend/**` 是 owner 09-23 自己做的委派（档 (c) 真需人拍），已立 **`Q-51`** 摆给他，带推荐与人话后果。**`Q-51` 未答之前丙段不派**。
- **AC#3 的 ⓐ/"换屏"这一格**跟着丙段走（它本身就是一枚新字段＝扩契约），并按普查 §3 的硬答案原样登记：**只做 ⓐ 换不了屏**（ⓐ 让 Go 能强制指定屏幕，用户点击仍然无处可去；而且没有甲段的泵，连 ⓐ 也没有传输通道）；**ⓑ 不是泵能解决的**——把那条路由加回来就是重开 `C17` 的方法白名单（`panel.ts:234-235`），那在 AGENTS §2 的停手名单上。⇒ 本票任何一节都**不许**被读成"换屏已通"。

## 编排者 09-26 10:5x 收表（按非实现者验收表 `145-snapshot-fields-landed-r1-accept-r1.md` 定勾；本程不重跑它的读数）

**表本体**：588 行／§0–§12／8 枚 commit（`ec2247a`→`d5a2498`，我逐枚 `git cat-file -t`＝commit、逐枚 `git show --name-only` 数"非本程路径"＝**0**）＋ 11 枚探针件（`probes/145-accept/`）。
**档位＝附条件**（六枚条件 C1–C6 **全部记录级、无一枚要求改码**）；分界按本仓规矩＝**验收程造没造出来推翻它**，它明写"五发独立构造全指同一侧"。

### 票面 6 枚 AC ↔ 表逐格（双向；这是 `A282②` 那条新规矩第一次在票面上落实）

| AC | 勾 | 裁它的是表的哪一节 | 一句话凭据 |
|---|---|---|---|
| AC#1 对照表 | **[x]** | §1（前提①）＋§2 | 加键响的是**四枚**用例、且"爆炸半径与枚数无关"（它自造单键变体仍响同一批四枚） |
| AC#2 落地有源那一集 | **[ ] 维持未勾** | §2（E1b 升级为真装配根） | **量出来的空**：新段出门是 `depth:0 / windowMs:0 / null / 全零`，同包真值却是活的。**⇒ 未勾不是欠账，是这一格的结论本身**；表的硬约束：本票 AC#2 **不许**被下一位读成"等下次顺手补两枚键" |
| AC#3 两种角色 | **[x]** | §6.4（验收程独立打过两发）＋实现件 §4 | ⓐ 不落（Go 侧不知道那九枚屏的名字）；ⓑ **停手上报**——它重开 `C17` 白名单，属 `AGENTS §2` 未定义即停那一条 |
| AC#4 契约轴零字节 | **[x]** | §8 第 4 点（反扫八枚 commit 名册） | 路径并集＝`probes/145/*`＋它自己的证据件＋`internal/panel/{composer,pump}.go`；**禁写面一枚未碰**、三枚 `_test.go` 保留件零命中 |
| AC#5 门禁与名册 | **[x]** | §5.1 | 四数逐枚复现（panel `rc1/105/58/1` 两版同数、`cmd/wisp` `rc0/139/79/0`）、名册三层两向 `comm` **皆空**、**没有被吞的第二枚红** |
| AC#6 反向判据 | **[ ] 维持未勾** | §3.2 三发＋§3.3 | 措辞按表的 C1 写死：**"本放开面内无定义域；与 AC#2 落地集非空互为条件"**——**禁用"待补"**（"待补"会把一枚结构上没有对象的格子读成一件小欠账） |

### 我这一侧的六个决定（条件的闭合归属）

- **C1 已照办**（就是上面 AC#6 那一行措辞）。
- **C2／C3／C4 已按"编排者代记"追加到实现件末尾一节**（证据件只追加不改写；C4 的行号我本轮自己现量过：HEAD 上 `type Snapshot struct` 在 **`composer.go:57`**、`NewSnapshot` 在 **`:76`**，票面/普查/实现件 §0.3 引的 `:44-49`／`:79` 一律视为过期）。
- **C5 已办**：台账小标号 `A273②` → **`A273③`**（我那两份派单正文里同一处也错了，就地各补一节更正，不回填）。
- **C6 采纳为全仓规矩**：**复算任何一发仓外读数必须用 blob 精确树或真 checkout**；`git archive | tar -x` 在本仓**不是纯净树**——我本轮自己复核到机制：`.gitattributes:1` 是 `* text=auto` 且 `core.autocrlf=true`，而只对 `go/md/sh/ps1/bat/sse` 显式写了 `eol=lf`；实发一枚：`.scratch/wisp/probes/149/combos2/x-p1-d11a.log` **blob 3391 字节 vs 工作树 3436**。⇒ 已写进 `dispatches/README.md`，**以后派单不再指定 archive 取版**。

### R1／R2（我的决定：两个都不给，理由写在这）

表只判了"结构上可否满足"（结论：**两者均可满足**，R2 的边际成本是一行），**没替我选，我也不让它替我选**。
- **R1（开一枚 `internal/panel/*_test.go`）＝暂不给**。它量到一枚更该防的形状：把 `TestXxx` 写进**非** `_test.go`（放开面内的 `composer.go`）**编得进、`go vet` rc=0、但永不执行** ⇒ 为了"让 AC#6 有个对象"而开门，换来的是一枚**看起来绿、从不跑**的断言，而那枚仪器缺口今天没有尺抓（表 §3.2／§6.3）。**先登记缺口，再谈开门。**
- **R2（给 `cmd/wisp/run.go` 一行读数）＝不给**。两个理由：①`run.go` 此刻是票 151 验收程在判的脸（冲突判到文件级）；②表 §8 已把话说清——**就算给了 R2，没有 `Q-51` 照样落不了地**（B／C 两把锁对任何一枚键都响），所以它买不到今天的进展，只买一次跨票抢面。
- 代价我认：本票净产出停在"注释＋判词"。表对这一形状的判断是**"这正是它自己选的形状，不是失职"** ⇒ 我按这句结案为**部分结案**，`-done` 不加。

**留给 owner 的一句（不催）**：面板要不要显示"当前是哪一屏"（`Q-51` 那一族）仍是界面侧的未答项；它不解，AC#2 这一格就一直是空的——**但这是设计要等的东西，不是后端的欠账**。

---

### 09-28 10:2x 编排者追加（新格 **AC#2b**：前端 composer 那两枚空数组 `models`／`efforts` 的**源我今天量到了**——本票 AC#2 当年那句"落地集＝空"对这两枚字段**已过期**，原句不抹）

**来路**：前端那侧递给 owner 的"要拍板三件事"第 ② 件（"模型列表＋思考强度档位 — PromptBar 的 `models`／`efforts` 现在是空数组（所以那两个按钮不画）。需要快照上多两个字段：有哪些模型、当前模型支持哪几档思考"）。它把这件事摆成了**待拍板**；我按盘上现量判：**它不是拍板项，是本票 AC#2 那一格欠的载体**。台账 `A360`。

**现量（锚 `039efb47`，本程现跑，每把尺可复制）**

1. **快照里没有这两维**：`grep -rn "type ComposerState struct" -A 26 internal/panel/composer.go` ⇒ 六枚字段＝mode／workspace／attachments／acceptedAttachmentMimes／maxAttachmentBytes／attachmentError（`internal/panel/composer.go:200-207`），**没有 models、没有 efforts**；`Snapshot` 四枚（`:57-62`）同样没有。⇒ 前端那两个数组为空**是 Go 没送**，与本票主结论同形（"宁缺毋造"是既有裁定，不是前端的缺陷）。
2. **模型清单有源**：`internal/config/schema.go:345`（注释逐字 `ModelSpec is llm.providers.<name>.models.<model-id>: one catalog entry`）＋ `:348`（`type ModelSpec struct`），带 `Display`（人看的名）与 `Capabilities`；尺 `grep -nE "type ModelSpec|// ModelSpec is|type Capabilities struct" internal/config/schema.go` ⇒ `:324`／`:345`／`:348`。
3. **"支持不支持思考"有权威位**：`internal/config/schema.go:324-332` 的 `Capabilities` 位集里**逐字有 `Thinking bool`**（还有 Text/Vision/AudioIn/AudioOut/Realtime/FC/Stream），且其注释写明"**声明可由人填，但必须由探测核实（票 11）**，探测结果本身住在 SQLite `provider_health`、这里刻意不留字段（存储分家）"。⇒ **档位的第一维（能不能思考）不是新数据，是既有位。**
4. **档位词表有源且已被校验**：`internal/config/schema.go:290-291` 逐字 `// ThinkingIntensity is one of off|low|medium|high (enum enforced).` ＋ `ThinkingIntensity string \`toml:"thinking_intensity" default:"off"\``；校验在 `internal/config/validate.go:211-233`（四枚 role 逐一验枚举）；适配器侧各家有各自映射（`internal/llm/openairesponses/request.go:62-63` 的 `effortLevels`＝low|medium|high；`internal/llm/anthropic/adapter.go:18` 映射到 `thinking` 参数）。⚠ **注意这里有一形真分歧要裁**：配置侧词表是**四档（含 `off`）**，OpenAI 兼容侧映射表只有**三档（无 `off`）**——`off` 在那些适配器上是"不发这个参数"还是"发不出去就报错"，本票 AC#2b 必须逐家量出来再画进快照，**不许把四档原样送出去当作每家都支持**。
5. **现成的读面**：`cmd/wisp/models.go:104 cmdModels`（`wisp models` 那条腿，票 131 给它补过钉）已经在"列目录里的模型"这件事上取过一次真值 ⇒ AC#2b 的正解形状是**复用它的取数**，不是新写一套。

- [ ] **AC#2b（09-28 追加，未做）**：把上面 1～5 落成快照里的两维——**① 当前可选模型清单**（provider＋model id＋`Display`，来源＝配置目录 ∩ `enabled=true`；⚠ 不许把未启用的也列进去）；**② 当前模型的思考档位集**（＝`Capabilities.Thinking` 决定"有没有这一维" ＋ 该适配器实际接受的档位词表决定"有哪几档"，两问分开答）。**判据要钉"值来自真源"**：拿一份两模型／两档位的假配置树喂进去，快照里读到的清单与那棵树逐字一致（与票 92 那格同形），⚠ 不许用"字段非空"充当判据。**本格只做"显示"**：任何"面板改档位／换模型"的**写回**都不在此格，它要 `internal/panel/bridge.go:42-45` 那四枚 `panel.*` 之外**新增方法＝C17 契约变更** ⇒ 已并入 `Q-64`（**默认不做**）。
  ⚠ 与 AC#2 的分工：AC#2 的结论"落地集＝空"**保留不抹**，它当时量的是"十四行表里那些字段有没有源"；本格是**新增两维**，所以**不算推翻 AC#2、也不算替它翻勾**。⚠ **判据要说准（我一度用错尺）**：`sed -n '3473,3488p' docs/PLAN.md | grep -nE "模型|档位|思考|monitor"` 现量**命中 2 行**——但那两行是**状态名"思考中（等 LLM 首 token）"**（`:4` 与 `:9` 那两行），**与"思考档位"同名不同物**；十四行表里**真正要求"模型清单／档位集"的行＝零枚**（尺：`sed -n '3473,3488p' docs/PLAN.md | grep -nE "模型|档位|effort|thinking_intensity"` ⇒ 逐枚点名后为 0）。⇒ **别拿"grep 到'思考'两字"当"这一维已被要求过"的证据**，这是"同名不同物"那一族在中文关键词上的复发。

## Progress log（append-only）

### 09-28 17:4x–18:0x　`145-r2` 写腿（派单 `.scratch/wisp/dispatches/2026-09-28-174x-impl-145-r2-snapshot-growth-with-nail-preapproval.md`；证据件 `docs/evidence/s1/145-snapshot-growth-r2.md`）

- **三格 `AC#2`／`AC#6`／`AC#2b` 一律未勾**（勾要非实现者表），本程一格未自勾、票面原文一字未改、未抹。
- **落的维数＝0；拿掉的维数＝14 枚候选逐枚点名**（证据件 §2 那张表：`composer.models`／`composer.efforts`／
  `approval.windowMs`／`approval.vetoChannels`／`approval.remainingMs`／`run.reasoning`／`run.usage`（实时）／
  `run.status`／`tools[]`／`cost`／`failures[]`／`failures[].humanText`／`view`／`approval.depth`——
  每枚都写着"真源在哪一行"与"AC#6 答不出的原因"）。**一行生产码未改。**
- **构造点枚数（派单 `起手必做` 第 3 条那把尺按字面量是空的，具名更正）**：
  `grep -rn "PanelSnapshot{" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ **0 枚**——Go 侧那枚类型叫
  **`Snapshot`**（`internal/panel/composer.go:57`，四字段 `:58-61`），`PanelSnapshot` 是 TS 侧 interface 名
  （配对行在 `composer_test.go:60`）。按真身现量：`Snapshot{` 非测试命中 **5 枚**，其中面板的 3 枚**全不是装配点**
  （`composer.go:97` 是 `NewSnapshot` 自己的 return、`pump.go:228` 与 `cmd/wisp/panel_pump.go:248` 是零值错误支），
  另 2 枚同名不同物（`internal/perm/store.go:247`、`internal/ball/liquid.go:276`）。
  **真·生产构造链只有唯一一枚装配根＝`cmd/wisp/run.go:442`（读口行 `:443-448` 共六根线）→ `internal/panel/pump.go:207`
  → `internal/panel/composer.go:97`**；`NewSnapshotPump(` 的非测试调用者＝**1**。
- **字段枚数 起手→终态：`Snapshot` 4→4、`ComposerState` 7→7（本程未加字段）**。
- **`pump_test.go:124`／`:276` 那两枚具名解冻本程未使用**（两行一字未动）：没有第五枚键时把它们从"恰好四枚"
  换成"随结构体"＝买不到任何东西、还白丢一枚钉；派单自己那句"没有正控＝把门换成好看"在此更严。
  改法与正控形状（含 r1 §2.7／R6 那枚"无标签导出字段两把尺看不见"的仪器缺口要一起钉住）写在证据件 §4 R-2。
- **停手报回三枚**（证据件 §4，具名、可复核，本程未擅自办也未据此扩权）：
  **R-1** 派单前提"在生产构造点填真值"落不了地——`cmd/wisp/run.go:442` 不在这次写面（`A388` 逐枚点名）；这就是 r1 §8 的
  R2，只是今天**只剩它一把锁**：C 锁已由派单自己解开、B 锁已被判成"更红是预期后果"，而 D 锁（`run.go` 读口＋记录点）
  与 E 锁（**AC#6 的用例之家**，r1 的 R1）仍在。⇒ 本程按"答不出的字段直接不要"结为**落地集空**，不造填常量的字段。
  **R-3** `AC#2b` 两枚前提与现量不符：①票面点 5"复用 `cmd/wisp/models.go:104 cmdModels` 的取数"指向的是
  **签名下载清单**（`modelsList` `cmd/wisp/models.go:205-225` 打 `store.manifest.Models`，`:206-208` 注释逐字
  "not about this boot's [models] section"），**不是** `llm.providers.<name>.models.<id>` 配置目录（`schema.go:345,348,400`
  ∩ `:369-371 Enabled`）——复用复不出来；②票面点 4 要的"各家真接受的档位词表"**量出来了**：
  `openairesponses/request.go:63`＝三档、`anthropic/request.go:92`＝三档、`openaichat/adapter.go:14-17`＝**一档都不映射**
  （"intentionally NOT mapped…严格端点会 400"），而**这三张表全是包内私有 `var`、全仓无导出访问器**
  ⇒ 这一维缺的是**源**不是字段；另 `Capabilities.Thinking`（`schema.go:330`）是人填的声明位，
  核实位在 SQLite `provider_health`（读 DAO `internal/memory/dao_providerhealth.go:170,178` **零枚非测试调用者**），
  发现式录入把 capabilities 全填 false（`internal/llm/discover.go:116-121`）⇒ 不许把"未探测"画成"不支持"。
- **在册三枚红**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／
  `TestC21DesignTokensFourWayAgree`）：进场逐名复量一致、**未修、未放宽、未豁免**；`frontend/**`／`design/**` 零写面，
  TS 那份"应当长这样"逐键清单在证据件 §5（由 owner 自己带）。
- **门禁读数**（原始件 `.scratch/wisp/probes/145/r2/`）：`./internal/panel/` rc=1／`=== RUN` 137／PASS 81／FAIL 3／SKIP 0；
  `./cmd/wisp/` rc=0／RUN 149／PASS 89／FAIL 0／SKIP 0（DLL 进 `PATH` 直跑，未走 `scripts/wisp-cli-tests.sh`）；
  `gofumpt -l internal/panel/ cmd/wisp/` **空**；`sh scripts/d22scan.sh` **rc=0 clean**（bans #1-5 `internal/=211`＋`cmd/=24`、
  ban #6 `frontend/=85`、ban #7 `internal/tools/=21`、ban #8 `design/=39`／`frontend/=85`／`internal/=441`／`cmd/=47`）。
  ⚠ **欠读一枚照实登记**：`probes/154/gate-clauses.sh` 本程**未跑**（零码改动；证据件 §7 N3）。
- **Git**：只 commit、**未 push**；每次带显式 pathspec 逐枚点名（起手 `git diff --cached --name-only`＝空，全程未动别家暂存态）；
  禁面（`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`bridge.go:42-45`／
  三枚冻结测试件／`internal/risk/**`）零字节；仓内零删除命令。
- **离"那十四态能真画"还差什么（本程现量口径）**：卡点已不再是"哪一维没源"——普查甲组那 12 行的源大多还在，
  真正没源的仍是那 7 枚（`thinkingMs`／`reasoningMs`／`durationMs`／`humanText`／`fragment`／`IconClass`／`remainingMs`）
  加 `run.phase`（映射表不存在）与本轮新点名的**"各家档位词表"（私有 `var`、无导出）**；
  而**有源那批今天一律还差同一件事**：装配根多接一根读口线（R-1）＋一枚能落断言的用例之家（R-2）。
- 09-28 18:0x 编排者收 `145-r2`（四枚 commit，零产码、三枚停手上报）：14 枚候选逐枚点名拒收（含 `approval.depth`＝复述字段当装饰）＝**`AC#6` 第一次真生效**，`AC#2` 未落不勾。⚠ **真锁在我的排程里**：唯一生产装配根＝`cmd/wisp/run.go:442`（`panel.NewSnapshotPump`），而 `run.go` 是我给所有写腿划的禁区 ⇒ **`145-r3` 与票 174 `AC#2b` 合并成一程、同开这一行**（射程＝只到给 `PumpSources` 多接一根读口，不动 `:365`、不碰冻结文字）。另记我派单两把坏尺：用 `PanelSnapshot` grep 构造点必得 0（Go 侧叫 `Snapshot`）；我预解冻的 `pump_test.go:124/:276` 在无新键时不该改（顶回成立，作废）。账 `A393`。

### 09-28 18:1x–18:5x　`145-r3` 写腿（派单 `.scratch/wisp/dispatches/2026-09-28-181x-wave2-impl-145r3-and-188r2.md` §F；证据件 `docs/evidence/s1/145-snapshot-growth-r3.md`）

- **AC 框四枚一律未勾**（`AC#2`／`AC#6`／`AC#2b`＋票 174 `AC#2b`；勾要非实现者表），票面原文一字未改、未抹。
- **先更正票面一处错前提（`AC#2b` 点 5，具名更正、原句留着不抹）**：那一句"现成的读面 `cmd/wisp/models.go:104 cmdModels`……**复用它的取数**"**指错了对象**——`cmdModels` 走的是 `modelsList`（`cmd/wisp/models.go:205-225`，打 `store.manifest.Models`），**签名下载清单**，`:206-208` 注释逐字 "not about this boot's [models] section"；它**不是** `llm.providers.<name>.models.<id>` 配置目录（`internal/config/schema.go:345,348,400` ∩ `:369-371 Enabled`）。⇒ "复用复不出来"这条已由 `145-r2` 报过（其 §4 R-3①）、`A393` 采纳，本程按派单**没有**据此往快照里塞 `models`／`efforts` 任何一枚（那两维＝票 187 射程：三家档位词表全是包内私有 `var`、全仓无导出面）。
- **落的维数＝1（两枚键，嵌在既有 `composer` 段内，顶层键集未增）**：`composer.currentModel`＋`composer.modelKnown`。
  真源逐跳到行：`llm.Endpoint.Model`（`internal/llm/resolver.go:46`）→ 生产赋值 `cmd/wisp/run.go:308 rt.endpoint = ep`（`:314-315` 的 `rt.provs`/`rt.names` 同一枚 `ep`）→ 读口 `cmd/wisp/panel_pump.go currentModel()`（`gitView()` 之后）→ `internal/panel/pump.go PumpSources.Model`（`Git` 之后）→ 装配根 `cmd/wisp/run.go:442` 字面量 `Model: rt.currentModel,` → `pump.go Snapshot()` 填段。
  `AC#6` 答句＝`internal/panel/composer_test.go` 新增的 `TestComposerCurrentModelTravelsOnlyFromItsReader`（三臂：有读口逐字带出／无读口空＋`modelKnown=false`／读到空串两态不塌，另钉出口字节含 `"currentModel":"glm-5"`）。
- **拿掉的维数＝6 组**（每组的"答不出"写在证据件 §2 那张表）：`composer.models`／`.efforts`（派单具名排除＝票 187）、`pending[].taskId`／`.callId`（**真源齐全**：`internal/tools/bridge.go:331` → `approval/gate.go:550` → `pending_read.go:46`，但 `internal/panel/approval_test.go:105 TestApprovalCardViewJSONKeysMatchFrontendTypes` 那枚**今天绿着**的双向键集尺会被打红 ⇒ 按撞钉预检停手，见证据件 §4 N1）、`pending[].position`（＝队列下标 `+1` 的复述，同 `145-r2` 被拒的 `approval.depth` 一族）、`pending[].timeoutMs`／`windowMs`（源真且配置驱动 `run.go:379-380`，但 `approval/gate.go:416` 逐字"an L2 approval never opens an L1 window"⇒ 盖到 L2 卡上是造一形假相邻；`approval` 段本身＝顶层键，禁）、`results[].*` 与 `tools[]`／`cost`／`failures[]`／`view`（记录点 `cmd/wisp/run.go:800-816` 本程禁写、单位口径未裁、Go 侧无九枚屏 id）。
- **`run.go` 只动派单具名的两行**：`:365` `tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}`（＝票 174 `AC#2b` 的装配那一行；`rt.paths` 造于 `:330`，已喂 `FSDeps:346` 与桥 `:453`）＋`:442` 字面量多接一根 `Model:`。两枚都零新增注释、零逻辑；**`:651-670` 那批状态词写入者与 `:800-816` 记录点零字节**（票 196／票 153 射程）。⚠ 上一条编排者那句"不动 `:365`"由本轮派单具名改判（`A395`→§F），照派单执行。
- **`pump_test.go:124`／`:276`：未动，且仍绿**——本程顶层键集仍是 `composer,generatedAt,pending,results` 四枚（新键嵌在 `composer` 段里），派单"只有真落第五枚键解冻才算生效"这一条件**未触发** ⇒ 解冻不追用、正控无从属必要项。复算尺：`go test -count=1 -run 'TestSnapshotJSON|TestThePumpExitCarries|TestAPumpWithNoReaders' ./internal/panel/` ⇒ `ok 0.054s`。⚠ 本程**未跑**新维的变异正控（`-overlay` 禁用、仓外副本要带 `frontend/dist`、预算到顶），牙齿自证写在证据件 §4 N3 并标〔仅自述〕。
- **用例之家一枚裁量**：断言追加进 `internal/panel/composer_test.go`（非三枚冻结件、非 `pump_test.go` 那两行、只追加一枚 `TestXxx`、零删除），口径依 `A389` 写面的 `internal/panel/**`；若编排者判写面只到 `A388` 那三枚文件具名，这一枚请按**越界退回**处理而不是放宽断言（证据件 §4 N2）。
- **在册三枚红**（`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）：改前改后**逐名同册、名数未增**，未修、未放宽、未豁免；`frontend/**`／`design/**` 零写面，TS 那份"应当长这样"逐键清单在证据件 §5（含本程已发两枚＋仍未发的四枚，**由 owner 自己带**）。
- **门禁读数**（口径＝顶层 `--- FAIL` 行计数，与 `145-r2` 的 `-v` 口径不同，别当同尺复算）：改前 `./internal/panel/` rc=1／红 3（逐名同上）、`./cmd/wisp/` rc=0、`./internal/tools/` rc=0；改后三把**同数**；新增用例单跑 `ok 0.039s`；`go build ./...` 净、`go vet ./internal/panel/ ./cmd/wisp/` 净、`gofumpt -l internal/panel/ cmd/wisp/` 空；`sh scripts/d22scan.sh` **clean rc=0**（ban #8 `internal/=442`／`cmd/=47`、bans #1-5 `internal/=211`）；`probes/154/gate-clauses.sh` BAD 名册**仍只有 `G6neg`**（那格 `实测=3枚` 是 `188-r2` 已登记的涨法，非本程所加）；`flip-declaration.sh` 未跑（派单禁）。
- **Git**：只 commit、**未 push**；三枚产码／判据 commit `4db5f3f6`／`1ba16de0`／`0e5c0d4a`＋台件一枚，每枚带显式 pathspec；禁面（`go.mod`／`go.sum`／`thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／`bridge.go:42-45`／三枚冻结测试件）零字节；仓内零删除；被拒调用 0。
- **超预算具名**：硬顶 ≤35、实跑 ≈39（去向：写面与真源逐枚核 5＋选型 3＋门禁两遍 2＋台件 4；详见证据件 §8）。**未据此放宽任何断言。**
- **离"那十四态能真画"还差什么（本程现量口径）**：装配根那把锁（`A393` 的真锁）**已开并已用掉一根读口**；现在最贵的两把新锁是 ①`approval_test.go:105` 那枚**绿的**双向尺（它把 `pending[]`／`results[]`／顶层三族全钉在 `frontend/src/lib/panel.ts` 上，`Q-51` 不答就一枚都进不去）②记录点 `run.go:800-816`（`reasoning`／`usage`／`tools[]`／`failures[]` 四维的数据只在被打印的那一刻存在，泵读不到）。`AC#2` 的落地集今天**不为空**（1 维 2 键），但十四行表里那些行**仍没画全**。

## 标题级读数更正＋再框定（10-03 09:5x，只读腿 `145-c2`＝`.scratch/wisp/probes/145/c2/census.md`，锚 `5f9ff9d4`，六枚 pathspec commit；⛔ 原句不改，就地追加；落账 `A560`）

- **两枚标题级读数已过期**（腿的数，我按〔腿报，未复核〕挂着，翻任何勾之前要我自己现跑）：`Snapshot` 现在**6 枚直接字段**（`internal/panel/composer.go:57-92`）⛔ 不是标题那句"恰四字段"；`ComposerState` **11 枚**（`:235-270`）⛔ 不是 AC#2b 的"6 枚"也不是 r2 的"7→7"。顶层 6 枚**全部有生产者并已接**，唯一装配根 `cmd/wisp/run.go:699-726`，十枚读口全接线 ⇒ **推翻票面⑤"Snapshot 从来没有生产代码构造过"**（该论在 r1 锚点成立、后被票 35 的装配根取代＝**过期非错**）。
- ★**对我有用的一记再框定**：新增这 6 枚喂的是**"设置／名册"侧**，**无一枚喂那十四行工作态** ⇒ 票面①"缺 8~10"的方向**仍成立**；真卡点两类＝**跨-seam**（`cost`／`usage`／`reasoning`／`stuck` 有活源，但活在 `internal/agent/loop.go:367` 的 run 局部、泵读不到）＋**纯无源**（计时 elapsed／人话文案／命中片段／实时剩余／phase 映射）。恒空两枚＝`ComposerState.Attachments`（`internal/panel/pump.go:242` 硬写 nil）与 `attachmentError`（泵从不填）。
- **"宁缺毋造"这侧的现量**：现役字段**没发现**把"没测过"画成"空闲／零"（每个后加维都带未知守卫，三形已落＝伴生 `*Known`／保留 sentinel 枚举／非真态强制 Reason 非空）；残余两枚＝`attachmentError` 无 known 位、`StreamLog.TruncationFor` 未知塌零值（`pump.go:651-655` 自知）。
- 一枚**未复认的否证**留在册上不下结论：前人 r1 §2.6"`tool_call.started_at` 无写者"——本程不据它派腿，要派先按第 74 条规矩把那支的"作用面"跑一遍。
- ⚠ **与票 167 的关系**：`Q-51`（序号落哪枚 JSON 键）**已由 `167-c3` 答死＝必须新增键、不能复用**，判语与凭据在票 167 新节；本票那把双向尺（`approval_test.go:105`）的解冻条件随 167-r2 一起走，⛔ 两条腿不许各加各的。

## 机主 10-07 两句改口径（⛔ 原句不改，就地追加；落账 `A655`）

原话照录：①「**给你所有权限，你可以想看啥看啥，不用再问我了**」②「你肯定要跟前端那边核对好要对接的字段，**这个你随便对接**」。⇒ 对本票产生三处**实际**改动，其余一律不变：

1. **AC#2 里那句归属已失效**（原句：「TS 侧 interface 的对齐**不写 `frontend/**`**，只在证据件里给出'应当长这样'的逐键清单，**由前端会话自己落**」）。两重失效：那个会话 09-28 就废弃了（⛔ 仍不许跨会话转达，这条机主没收回），而今天他**放开了读**。**改后口径＝页面代码我可以读（只读）、由我出对照表、Go 侧字段该加哪几枚由这张表定**；⚠ **读≠写**：他给的是"想看啥看啥"＋"字段你随便对接"，⛔ **我没有把它读成"可以往那棵分支写代码"**，真要写页面文件是另一枚更大的动作，届时按规矩单独具名。⇒ 本票 AC#2 的"逐键清单由别人落"那一支**作废**，但**框我不翻**——它要等 `145-f1` 的表与非实现者读数。
2. **本票过去三年的半个洞：页面那一侧从没人量过。** `A560`（10-03）已经把我自己标题那句"恰四字段"更正成 **6 枚直接字段**，但**"页面消费哪几枚"这半张表始终空着**——所以"缺 8~10 枚"这个方向今天是**从 Go 侧推出来的，不是从页面侧量出来的**。⇒ 已派只读腿 `145-f1`（`.scratch/wisp/probes/145/f1/field-map.md`）去抽这半张表，方法具名：**只用 git 对象层读那枚分支的 blob**（`git show dsh/feat/frontend-p0-v2:<path>`，`frontend/src` 下 **85 枚**文件），⛔ **不进 `D:/wt/fe` 那棵工作树一个字节**（既有边界里唯一他明确收回的就是"不许看"，"不许进那棵树"我没有替他收回）。
3. **AC#6 的判据现在有真凭据可查了**：原句问"哪一枚用例断言了它的值来自真来源"。⚠ 别忘了本仓**已经有一枚双向尺**在把 Go 快照字段钉到 TS 侧（`145-c2` 报的 `approval_test.go:105`，它把 `pending[]`／`results[]`／顶层三族钉在 `frontend/src/lib/...` 的形状上）⇒ **加字段要先过它那关**，这条解冻条件是 167-r2 同批，⛔ 不许 `145-f1` 或后续写腿绕过它来"让门变绿"。

## 收 `145-f1`（页面字段名册；件 `.scratch/wisp/probes/145/f1/field-map.md`，**353 行／48,883 字节**；四枚 commit `53d46242`→`7c1ac218`→`a8d607b2`→`c1144180`；锚＝dev `a781fdc8`／页面 `16c2f038`（`frontend/src` 85 枚）；10-07 11:4x，落账 `A657`）

★**这半张空表填上了，而它推翻的是本票的标题级前提。** 四档枚数（腿的尺，逐档我都另起一把复认）：

| 档 | 枚数 | 我的独立复跑 |
|---|---|---|
| ①两边都有 | **38 枚叶子**（3 顶层＋10 卡＋3 chunk＋6 composer 段＋3 mode＋5 workspace＋8 attachment） | ⚠ 见下"挂载"一节——**真面板路径只 21 枚** |
| ②页面要、Go 没交 | 快照契约上**只有 1 枚**＝`view`；另 props-only 族 41 叶 | `git grep -E '\.view\b' … -- 'frontend/src'` ＝ **4 处读取**，而 Go 侧 `Snapshot` 六枚里没有它 ⇒ 与 `panel-views.ts:89-93` 那句"ticket 35's pump owns it"对上 |
| ③Go 交了、页面不读 | **8 枚具名 key ＋ 2 枚整段（29 叶）** | 我按"这句话在说什么"扫（不是只扫符号名）：同尺打 `frontend/src` 目录形式 ⇒ `.generatedAt`／`.tasks`／`.instructions`／`.git`／`.currentModel`／`.modelKnown`／`.credentialState`／`.credentialKnown`／`.rewritten`／`.artifact` **全 0**，正控 `.pending`＝5／`.results`＝4／`.composer`＝2 ⇒ 读数成立 |
| ④名字像、词汇不同源 | **3 族** | `TaskRow.id/title ↔ TaskRowView.taskId/label`；`GitBranchView.isRepo(bool) ↔ GitView.Kind＋Reason＋SwitchBlocked`；状态词汇**四套并存**（D43 名／`TaskStatus`／`SessionState`／`MonitorAgentRow.state`）⇒ ⚠ **这不是改字段名能对接上的**，直接改名会把 D43 判决当人话画、把"不许切"那条真判决吞掉 |

**1. ★本票"缺字段"这个方向里，有一族病因根本不在快照：挂载。** 腿报"17 枚叶子的读取点住在没人 import 的组件里"，我逐枚复跑（尺＝`git grep -l -E 'components/<名>"' <ref> -- 'frontend/src'`，剔组件自身）：
- `composer.tsx` → **0 处引用**；`chat-screen.tsx` → **0 处引用**；而 `result-stream.tsx` 的**唯一**引用者就是那枚已挂死的 `chat-screen.tsx`（⇒ 传递性死）；`approval-screen.tsx` 只被 `showcase.tsx` 挂，而 showcase 走 `main.tsx:104` 的 `?harness=2` 路由——**不是真面板路径**。
⇒ **我采纳，且这条改变派单口径**：`attachments`／`acceptedAttachmentMimes`／`maxAttachmentBytes`／`attachmentError`／`l2ConfirmNames`／`spelling`／`reparse`／`workspace.reason`／`done` ＋ `ComposerAttachment` 8 枚 ⇒ **⛔ 不许算进"Go 该加的字段"**。把它们当载体扩张派给 Go 侧，会加出**真数据＋假通路**。
- ⚠ 诚实标注：腿 §5 第 12 条自报这把挂载尺只打了 5 枚组件、**没跑 `tsc`/lint 证伪** ⇒ "17／21"这两枚数**只能当量级、不能当验收凭据**；真要派"改挂载"的活，验收方要补一枚挂载级尺（我写进下面的派单约束）。

**2. ★★最硬的一条不在字段层，而在"有没有耳朵／有没有推送"——而它不归本票，归票 35。** 腿说"最后一跳不存在"，我把它拆成四把尺复跑（⛔ 不是抄它的）：
- ⓐ`grep -rn --include=*.go -E 'PostWebMessage|EvaluateScript|CreateWebMessageAsJson' cmd/wisp internal/panel`（剔 test）＝**0**；正控＝同一把尺全仓只命中 **1 枚**文件 `.scratch/wisp/probes/33/p1/q2/main.go` ⇒ **今天"推给页面"这件事只活在一枚探针里，不在产码里**。
- ⓑ记账值的读者：`(rt *agentRuntime).lastPanelSnapshot()` 带括号数调用点＝**10 枚，全部在 `_test.go`，非 test 0 枚**（`cmd/wisp/instructions_200r2_test.go` 3／`panel_pump_test.go` 3／`subagent_*_197_test.go` 4）。⇒ 泵 marshal 出的字节**今天只被测试读过**，没有任何生产读者。
- ⓒ宿主出向能力面：`func (m *PanelManager) …` 导出方法**共 8 枚**＝`IsCreated/IsShown/LastColdMs/LastHotMs/Show/HotShow/Hide/Destroy`——**没有一枚是"把数据送进去"**；入向那扇 `Bind(panelDispatchBinding, …)`（`panel_host_windows.go:405`）是**页→Go**，方向相反。
- ⓓ页面耳朵：`git grep -E 'webview\.addEventListener|postMessage|onmessage'` 于 `frontend/src` ＝ **0**；正控＝同树 `addEventListener` 在 8 枚文件里命中（`App.tsx` 2 枚等）⇒ **尺有牙**。页面唯一的桥引用是 `lib/panel.ts:154 return window.wispBridge ?? window.chrome?.webview ?? null`，**只用于发起调用**。
⇒ **判语（我裁）**：面板"画不出真数据"的第一阻断**不是快照太瘦**，而是**出向那一跳（Go→页）与页面接收器两枚都不存在**。本票的字段扩张若先落地，结果是**更瘦的通路带更肥的字段**，仍画不出来。⇒ 字段名册作为**对接依据入册**，而"先做哪一枚"的排序交给**票 35**（`35-panel-bridge-c17.md`：它就是 C17 那座桥，含 Go→页事件推与 `panel.resync` 全量推）——⛔ **我不另开新票重复它**，只在它面上补一节现状对拉（机主原话"不要重复劳动哈"）。

**3. 我顺手替腿结掉它 §5 的一格（它没做的赋值穷举）**：`snap.Instructions =`／`snap.Tasks =` 的真身＝**`internal/panel/pump.go:304-327`**（两枚都包在 `if p.src.Instructions != nil`／`if p.src.Tasks != nil` 里），而装配根**两枚 reader 都给了**：`cmd/wisp/run.go:718 Instructions: rt.instructionBundle`、`:724 Tasks: rt.taskRosterState` ⇒ **"泵发出的包带不带这两段"＝带**（只要走 `wisp run` 那条泵）。
- ⇒ 但同一段里我量到一枚**腿没看见、也不在本票射程**的缺口：`PumpSources.L1Windows`（`pump.go:200`）在全仓**零装配点**（尺＝`grep -rn 'L1Windows:' cmd internal` 剔 test ＝**无输出**），于是 `pump.go:324` 那枚 `if p.src.L1Windows != nil` 恒假、`waits` 恒空 ⇒ **票 220 的 `blockedOnApproval` 第二来源今天是黑的**。这条**归队列里那枚待派腿**（task #297「给 L1Windows 供数据那一跳」），⛔ 不算进本票，我只在此具名登记它的现量。
- ⚠ 另记一枚账面过期：`cmd/wisp/run.go:695-697` 注释逐字写着"this tree carries no WebView2 host (tickets 33/35)"，而 `panel_host_windows.go:405` 今天真在 `Bind` ⇒ **注释与代码相反**，与票 274 §2① 同族（那条我另案处置，本票不动它一字）。

**4. 腿顶回我的两处，我复跑后都算它对（原句不抹）：** ①票面"恰四字段"→ 现量**六枚**（我 A655 §6 早自量得同一枚数，两把尺同读数）；②票面"`UnfedScreen`：除 chat/approval 外每屏渲染一句人话"→ **`UnfedScreen` 全树 0 命中**（正控 `RetiredView` 命中 `App.tsx:327`／`:495`），今天的形状是 **approval／tasks／ball 三枚被判"不是屏"**、其余落到 `Conversation`，而屏名册只有 **9 枚**（`panel-views.ts:26-35`）不是 14 枚。⇒ ⚠ **本票标题里那枚"十四态"是 `PLAN.md:3473-3488` 那张状态表的行数，不是页面屏数**——这两件事以前被我混着讲过，此处具名更正。

**5. 排程与归属（AC 框一枚不翻）**：
- **AC#2**：那半张空表现在满了，字段该加哪几枚**由这张表定**——但表里的 ② 只有 `view` 一枚真缺，而 ③／④ 两档说明**今天的问题不是"Go 少交"而是"页面没接、且词汇不同源"**。⇒ AC#2 的落地集**不因为表交了就能勾**；它等**非实现者**按 ①②③④ 四档逐档判"谁改哪一边"。
- **AC#6**：维持未勾，凭据现在有两枚（`approval_test.go:105` 双向尺 ＋ 本表 ③ 档的 0 读数尺）。
- 下一步不是本票的写腿，而是**票 35 的现状对拉**＋**票 274（构建带页面）**；本票的写腿排在两者之后。⛔ 派单给后续腿时必须带上这一句：**"不许把 ① 档 38 枚读成 38 枚通路"**（真面板路径只 21 枚，且最后一跳为 0）。
