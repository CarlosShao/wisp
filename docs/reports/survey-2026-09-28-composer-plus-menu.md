# 调研：聊天框那个「加号 / 斜杠菜单」在五家怎么分层（2026-09-28）

- 调研腿：`survey-plusmenu`（只读，未改任何仓）
- 下单原话（owner 09-28 22:1x）：「聊天框基本都有加号，点这个加号，都会有菜单选项，有一些/+命令的，比如/可以调出skills，MCPs，还有一些命令，比如/goal, /spec，/compact之类的……看看这块怎么设计的，要怎么做」「聊天框这块的内容不少的」
- 三张截图**均已读到**（Read 直接出图，非转述）：
  - `56a03fe6-…png`＝**DeepSeek Harness 的加号菜单本体**（浅底：`添加` 组 文件file/目标goal/计划plan/反馈feedback，`指令` 组 压缩compact/权限permission/模型model/下载日志export；输入框占位「发消息或创建任务，/ 调用指令，@ 文件或对话」）。这张图与 DSH 代码里的 `SECTION_ROWS` 逐项对得上，见 §1.1。
  - `6321f484-…png`＝**深色界面**：左菜单 目标/计划(Shift+Tab)/工作区文件/插件/技能，右侧一枚「搜索技能」面板列出 `mcp-config`、`skill-creator`、`security-scan`、`qoder-*`…＋末行「管理技能」。**这枚不属于本次五家任一的前端可核对产物**（我在五家仓里没找到与之逐项吻合的菜单代码），只当作"技能可以做成二级面板＋可搜索＋带管理入口"的界面证据。
  - `b97c027b-…png`＝**深色加号菜单**：`# 添加上下文`／`使用 / 调用命令和技能`／`上传附件`，二级 `命令 >` 展开出 `Spec / Plan / Goal`，另有 `插件 >`、`技能 >`、`飞书文档 >`、`飞书通讯录 >`；底栏 `完全访问`＋`Agent`＋模型 `Seed-Code`。**"命令/插件/技能"三类各挂一个二级**这一形，是 owner 那句"点加号有菜单"的最直白答案。
- 硬规矩：五家克隆**没有 `.git`** ⇒ 全文不引提交历史。每条"他们有"带 `文件:行`；每条"我们没有"带我们仓的现量读数；量不到写〔未量〕。
- 三态标注：〔已证〕＝读到实现码；〔建了但没接〕＝有码无调用者；〔仅文档〕＝只读到 README/注释。

---

## 0. 我们的现量（起手复算，先摆数字）

| 事实 | 现量读数（真实命令输出） | 结论 |
|---|---|---|
| 斜杠命令目录／注册表 | `grep -rln --include=*.go -iE "slashcommand\|commandcatalog\|CommandRegistry\|\"/compact\"\|/goal" internal/ cmd/` ⇒ **0 行输出** | 全仓 0 命中〔已证：我们没〕 |
| 输入框载体 | `internal/panel/composer.go:235-257` 的 `ComposerState` **有 9 枚字段**（`json:` 标签行数推导式：`sed -n '235,257p' internal/panel/composer.go \| grep -c "json:"` ⇒ **9**）：mode / workspace / **attachments / acceptedAttachmentMimes / maxAttachmentBytes / attachmentError** / git / currentModel / modelKnown | **不是"只有档位＋模型两枚选择器"——附件那四枚字段已经在**，只是永远为空（见下一行）。仍然**没有"名册"这一层**〔已证〕 |
| 附件受理器 | `internal/panel/attachments.go` **410 行**；`grep -rn --include=*.go "panel\.Attachment\|Attachments(" internal/ cmd/ \| grep -v internal/panel/attachments` ⇒ **0 行输出（exit=1）**；`grep -rn --include=*.go "\.Ingest(\|AttachmentBroker" internal/ cmd/` ⇒ 命中**只在 `attachments.go` 与 `attachments_test.go`** | 建了但没接〔已证〕 |
| 附件为什么永远为空 | 生产装配点 `internal/panel/pump.go:226`：`composer := NewComposerState(mode, workspace, nil, maxAttachment)` —— **第三枚实参是字面量 `nil`**；`PumpSources`（`internal/panel/pump.go:111-141`）有 Verdicts/Mode/Workspace/Git/Model/Results/Instructions/Tasks 八类读者，**没有 Attachments 读者、没有 Commands 读者** | **缺的是生产者，不是字段**〔已证〕 |
| 技能/插件配置键 | `grep -c Skill internal/config/schema.go` ⇒ **0**；`grep -c Plugin internal/config/schema.go` ⇒ **8**（命中的是 `PluginsSection` / `PluginEntry` / `tier2_enabled` 一族，`schema.go:125-131, 551-575`） | **插件有键、技能连键都没有**〔已证〕——票 215 现量表写的"skill／plugin 字样命中"要按这两个数收紧 |
| "写了不管用"的旗 | `internal/config/unwired.go` **147 行**，`grep -c "path:"` ⇒ **6 枚响亮拒收键**；文件头三条件口径：`Such a key lies: the user writes one line and the program swallows it.` | **有，且比五家都硬**〔已证〕 |
| 界面→Go 那一跳 | 票 213/214/215 前置段写死：今天只有 `wisp panel-inbound` 命令行路（账 `A410`） | 三票终态都到不了"界面点得动" |

> 这一节是**对照的基准线**：五家每一层的"有"，都要能落到我们这张表的某一行上，否则就是空谈。

---

## 1. DeepSeek Harness（DSH）——五家里分层最清楚的一家，也是截图 1 的本体

DSH 把这件事切成**三层四包**，这是我们最值得抄的骨架：

```
纯核心（无 DOM、可单测）      进程内服务（状态机）           业务包（各自注册自己）
ui-input-trigger/src/core/  →  ui-input-trigger/src/client/  →  ui-commands（/）
  detect.ts   触发词检测          service.ts  名册登记簿          ui-skill      （/ 里的技能）
  menu.ts     菜单分组归约        controller.ts 每会话控制器       ui-reference  （@ 里的引用）
  contract.ts 类型契约                                            ui-conversation（+按钮、file 动作）
```

### 1.1 菜单分几类、分类写死还是后端给（问题 1）

**答案：分类写死在前端，成员来自后端目录＋前端贡献，两者按名字合并。**

- 两个组的名字与顺序**硬编码在前端**：`packages/client/ui-commands/src/client/presentation.ts:19-22`
  ```ts
  const SECTION_ROWS: Readonly<Record<MenuSection, readonly string[]>> = {
    add: ['file', 'goal', 'plan', 'feedback'],
    commands: ['compact', 'permission', 'model', 'export'],
  }
  ```
  注释写明口径：`Row names per section, highest usage first; rows outside both lists close the Commands section in catalog order`（`presentation.ts:18`）。**没被点名的行不会消失**，它们按后端目录顺序排在 `commands` 组末尾（`presentation.ts:83`）。⇒ 这一条与截图 1 逐字吻合（截图里就是这两组这八行）。
- 组的标题本身是**本地化键**，不是硬串：`presentation.ts:82-84` 用 `t('section.add')` / `t('section.commands')`。
- 每一组＝一个"触发源"（source），组名就是源名：`ui-input-trigger/src/types.ts:163-167`（`name` 是菜单组标签、每个 trigger 下唯一、重名直接 throw；`order` 决定组序；`showGroupTitle` 决定要不要画组标题行）。
- 源登记表是**运行期注册**而非编译期常量：`ui-input-trigger/src/client/service.ts:55-77` `registerSource()`，重名 `throw`；注册顺序＝菜单组顺序＝`matchSpace/matchEnter` 的轮询顺序（`service.ts:24`）。
- 全仓实际注册的源只有 **4 处**（推导式：`grep -rn --include=*.ts --include=*.tsx "registerSource(" packages/ apps/ | grep -v "test\|spec"` ⇒ 除契约/实现自身外 4 枚调用点）：
  - `/` 命令源：`ui-commands/src/client/service.ts:115`
  - `@` 引用源：`ui-reference/src/client/index.ts:158`
  - `/` 技能源：`ui-skill/src/client/index.ts:233`
  - （第 4 枚是 `ui-input-trigger/src/client/contract.ts:18` 的接口声明本身）
  ⇒ **口径注意**：DSH 的"分类"不是"一个源一类"。`/` 这一个源内部又分 Add/Commands 两组，技能源是**另一个** `/` 源。也就是说**同一个触发字符可以有多个源、多个组**。

### 1.2 `/` 与 `@` 的触发语义（问题 2）

**答案：检测在前端纯函数里做；命令"是什么、能不能跑"由后端目录说了算；两边各有一道闸。**

- **触发检测＝前端纯函数**，跑在浏览器主线程，零 DOM 零框架：`ui-input-trigger/src/core/detect.ts:48-76`。规则（`detect.ts:20-30`）：
  - 词边界：行首、空白后、标点后才算触发；前一个字符是字母/数字/下划线 ⇒ 不触发（`user@host` 不触发 `@`）。
  - `/` 的两条 URL 豁免：紧跟另一个 `/` 不算（`//` 的第二根斜杠）；前面是 `:` 且 `:` 前还有非空白 ⇒ 不算（`https:/…`）。注释具名"both pinned by tests"。
  - 三档闸门 `TriggerGuard`：`plain` 两枚都活 / `claimed` **`/` 全抑制、`@` 活** / `frozen` 都不活（`types.ts:236-240`）。⇒ "已经在补全菜单里打字"时 `/` 不再二次触发，这个状态机是显式的。
- **`@` 用另一套共享文法**（不是同一根正则）：`detect.ts:6, 50-60` 调 `activeAtToken`（来自 `dsh-file-reference/grammar`），支持**带引号可跨空白的 token**（`types.ts:55` `quoted`）。
- **谁解析命令**：分两跳，**必须分开说**——
  1. **前端只解析"这一行是不是一个已知的触发 token"**：`ui-commands/src/client/service.ts` 的 `matchSpace`（同步、只读热缓存，`service.ts:277-284`）与 `matchEnter`（可等待目录就绪，`service.ts:311-364`）。判定表就写在 `service.ts:276` 的注释里：`Decision table, menu column / space column / enter column`。
  2. **参数语义与执行在后端**：后端只做一个严格文法 `parseCommand`：`packages/interaction/commands/src/index.ts:125-132`，正则 `^\/([a-z][a-z0-9_-]*)(?=$|[\t\n\r ])`，名字**强制小写**、`rawInput` 是"名字之后逐字保留的原文（含分隔空白）"（`index.ts:84-85`）。**参数不再由框架解析**——每个 handler 自己读 `rawInput`（例：`/goal` 自己判 `clear|edit|pause|resume`，见 `packages/goal/command-goal/src/index.ts:191-197` 的 hint；`/permission` 自己判空串/未知档名）。
- **模糊匹配＝有序子序列打分，不是前缀包含**：`packages/client/ui-primitives/src/rank-by-name.ts:70-93`
  - 口径：查询串必须是候选名（或本地化标题）的**大小写无关有序子序列**（`rank-by-name.ts:1-8`）；前缀命中排最前，然后按对齐分，再按输入顺序（`:90-91`）。
  - 打分细节：行首与 `-`/`_` 边界 `+8`（`:19-21`），相邻续接 `+4`（`:47`），跳字符与前置字符扣分（`:43-48`）。
  - **两个搜索键**：`name` 与 `label` 都算（`:78`）⇒ 中文标题下仍能用英文命令名搜到（`ui-commands/README.md`：`a localized title stays findable by its command name`）。
- **解析失败怎么提示（这是 owner 最关心的一形）**：
  - 前端**不猜**：`matchEnter` 里 `desc === undefined ⇒ return undefined`（`service.ts:333`），交给默认出口。
  - 真正的"未知命令"由后端裁决：`CommandRuntime.execute` 里 `parseCommand` 不匹配或名字查不到 ⇒ **返回 `undefined`**（`packages/interaction/commands/src/index.ts:367-371`），客户端把它当**拒绝**抛错。
  - 前端 `execute()` 的注释把这条写成了契约：`An unmatched line reports an error outcome (the composer's immediate admission feedback)`（`ui-commands/src/client/service.ts:391-393`）；README 更硬：`a command line is never silently downgraded to a plain prompt`（`ui-commands/README.md` Summary 段）。
  - **参数级失败带可用值回显**：`/permission` 未知档名 ⇒ `unknown preset "${name}" (available: ${this.names.join(', ')})`（`packages/interaction/permission-presets/src/index.ts:269-271`）。⇒ 不只说"错了"，还说"有哪些是对的"。
- ⚠ **措辞纪律**：以上只证明**DSH 这么做**。minimax-code 与 Step-Code 的命令解析在 §3、§4 各自单独取证，形状不同。

### 1.3 命令名册从哪来（问题 3）

**答案：后端目录（RPC）＋前端贡献两张表按名字合并；有明确的"能力没落地就不进名册"机制，一共三种写法。**

- **后端目录**：`commands.list` 是 `@Remote` 方法，返回**按名字排序**的 `CommandDescriptor[]`：`packages/interaction/commands/src/index.ts:316-323`。描述符只有四个字段：`definitionId / name / description / input`（`packages/interaction/commands/src/types.ts:57-66`）⇒ **后端不给图标、不给分组、不给本地化标题**，那些全在前端。
- **前端合并**：`ui-commands/src/client/service.ts:223-247` `candidates()`：先铺后端行，再叠前端 `contributions`；**同名直接 throw**（`:228-230`），注释口径 `collides with a host command`——**不是覆盖、不是遮蔽，是响亮失败**。
- **缓存与失效**：每会话一份目录缓存 `CommandDirectory`，取之前**必须**该会话已被 retain 且首次 history open 成功，否则不发 RPC（`ui-commands/README.md` 实现内部段；码在 `ui-commands/src/client/service.ts:96-110`）。失效源三枚：`commands/change` 事件 → `invalidateAll()`（`service.ts:124`）、`agent-preset/selected` → `resetSession()`（`service.ts:127-128`，注释：换预设会改变"这个会话的 agent 实际注册了哪些命令"）、`connection/reset` → `resetConnected()`（`:129`）。
- **"能力没落地就不进名册"的三种写法**（这一条对我们最有用）：
  1. **组合期条件挂载**：`ctx.inject(['commands'], …)` —— 只有命令注册表被组合进来才注册这条命令。见 `packages/plan/plan-mode/src/index.ts:230-231`（注释 `The command child activates only when a command registry is composed.`）与 `packages/interaction/permission-presets/src/index.ts:255-256`。**能力不在 ⇒ 那一行根本不存在**，不是"存在但置灰"。
  2. **运行期能力位**：`/permission` 的可选档名列表里 `Auto` 只在"它的集成还活着"时才追加：`packages/interaction/permission-presets/src/index.ts:279-282`（`this.autoAdmit === undefined ? [] : [AUTO_PRESET]`）。
  3. **前端逐会话 `available()`**：`CommandContribution.available(session)` 每次候选枚举都重新调（`ui-commands/src/client/contract.ts:88-89` 注释 `called with a fresh projection per candidate pass`）。实例：`/model` 在"寻址到子 agent 的会话"里不可用（`ui-model-selection/src/client/index.ts:149`）；`file` 动作看 `inputHub.canPickFiles(sessionId)`（`ui-conversation/src/client/apply.ts:275`）。
- **防"冒充内置"这一形**（我们票 215 的 AC#2 同款）：内置命令的菜单长相（标题/说明/图标/本地化拼写）是**按 `definitionId` 身份**选的，不是按后端给的描述文本选的：`ui-commands/src/client/presentation.ts:41-48` 那张 `HOST_FACES` 表 + README 的 `changing a Host description cannot change that selection. Same-name overrides without the matching identity keep their own copy and receive no first-party aliases.`。⇒ **后端把 description 改成"我是 /model"也不会拿到内置那套脸。**
- **`+` 按钮与打字的 `/` 是同一个菜单**：`ui-conversation/src/client/apply.ts:505-516` 的 `toggleCommandMenu` 直接调 `inputTriggers.toggleSource('command', { trigger:'/', query:'', … })`；`inputTriggers` 服务不在 ⇒ `toggleCommandMenu = undefined` ⇒ 按钮 `disabled`（`ui-conversation/src/client/skeleton/InputBar.tsx:427`）。**按钮的可用性挂在服务可用性上，不是画死了再点了报错。**

### 1.4 技能/插件在菜单里显示什么（问题 4）

- 技能候选来自后端 `skills/list` Remote；**`modelInvocable: false` 的技能（只能人调、模型看不见）在描述前加"仅用户"标记**：`packages/client/ui-skill/README.md`（`a modelInvocable: false entry … wears the user-only marker as a description prefix in the active language`）。⇒ 这是"模型可见 vs 人可见"这一维在菜单里的显式表达。
- **技能行不显示开/关状态，也不显示版本**：候选只有 `name/label/description/icon/hint/section/value/drill` 七个展示字段（`ui-input-trigger/src/types.ts:49-73`，注释 `Pure display data — zero behavior declaration`）。
- **"配置 true 但实际没生效"这一形，DSH 的防法不是"在菜单里显示生效状态"，而是把生效判定挪到执行侧并且不回显**：
  - 菜单里选技能＝**只往草稿里塞字面量 `/name `**（`ui-skill/src/client/index.ts:220-226` `onPick` 返回 `{ text: '/'+name+' ' }`），**真正加载在宿主的 pre-step 边界**（`dsh-tool-skill`），README 具名：`Determinism lives host-side`。
  - 会话里出现的 `Skill` 工具行**只从冻结的 call/result 切片取名与状态**，绝不查当前目录：`ui-skill/README.md`（`The row derives its name, lifecycle, and body only from the frozen call/result slice … never from the current catalog, so replay stays stable when installed skills … change`），失败时`failures replace the name with the first error line`。⇒ **"生效没有"看的是会话里那一枚工具行有没有红，不是菜单里的勾。**
- **技能目录拉取失败的 UX**：**静默丢掉那一组**，只在控制台记日志：`ui-input-trigger/src/core/contract.ts:48`（`Source failure = silent group removal (log only; no error UI tier)`）＋ `menu.ts:122-129`（`source-failed` 分支 filter 掉该组）＋ `ui-skill/README.md`（`A failed skills/list call is logged and folded into a silent menu-group drop`）。⇒ **注意：这是 DSH 的选择，不是"业界都这样"；对我们票 213 的 AC#3"拒了要说为什么"来说，这一形恰恰是我们不该照抄的地方。**
- 插件侧：`ui-plugin-manager` / `ui-settings-plugins` / `ui-settings-plugin-inventory` 是**设置页**的地界，不进加号菜单（菜单里只有命令与技能两类；见 `SECTION_ROWS` 与 4 枚源）。〔已证：菜单侧无插件行；设置页内部本轮未展开读〕

### 1.5 附件/选文件（问题 5）

- **`file` 是一枚前端 `action` 型命令贡献**，不是后端命令：`ui-conversation/src/client/apply.ts:270-278`，`ui: { kind: 'action', run: session => inputHub.pickFiles(session.sessionId) }`。`ActionSpec` 的口径：`a bare invocation consumes the trigger token and runs one client-side callback. It submits nothing`（`ui-commands/src/client/contract.ts:56-60`）⇒ **动作型命令不受"这条命令收不收附件"的约束**（因为它是"加附件"本身）。
- 选文件用的是**浏览器原生 `<input type="file" multiple>`**（隐藏起来，按钮触发）：`ui-conversation/src/client/skeleton/InputBar.tsx:433-439`，`disabled={subagent !== null}` ⇒ **子 agent 会话里直接不给选**。
- **上传在客户端 Worker 里流式进行**，带进度与取消，换回一枚**不透明收据**（receipt）供后续 prompt 引用：`packages/client/file-upload/README.md`（`Callers can observe consumed bytes and cancel an active operation`；`receive an opaque receipt for a later prompt`）。收据在后端侧的类型是 `{ type:'file', receiptId: string }`（`packages/interaction/commands/src/types.ts:15-18`）。
- **附件与命令的关系判在"哪一层"——两层都判，后端是权威**：
  - 前端先拦（体验层）：只有声明了 `input.attachments` 的宿主命令才放行，其余路径抛本地化 `attachmentsUnsupported`，**渲染成一枚 transient toast，草稿与附件卡片原地保留**（`ui-commands/README.md` `Attachment-carrying submissions` 段；码 `ui-commands/src/client/service.ts:322-324, 340-341, 349-350, 355-356`）。
  - 后端再拦（权威层）：`CommandRuntime.execute` 里三形各自 settle 成 error —— 未声明就收附件 `/${name} does not accept attachments`；没有附件存储 `attachments are unavailable because no attachment store is composed`；图片超限与收据解析失败各抛 `AttachmentError.message`（`packages/interaction/commands/src/index.ts:388-402`）。注释具名口径：`Attachment admission is enforced here, not in the composer`（`index.ts:353-357`）。
  - **校验失败不启动任何写入**：`Validation rejection starts no attachment writes`（`index.ts:355-356`）。
- **选目录**做成**两枚可替换的包**（同一对 slot 的两种填法）：原生 OS 对话框版 `ui-directory-picker-native`（README：`opens the operating system's own chooser on the local machine and reports the single outcome — a picked path, a cancellation, or a failure`；`Choose it when the browser runs on the same machine as the Host`）与自绘浏览版 `ui-directory-picker-browse`（远端/进程内浏览器用）。**切换是组合配置改一行，不是改代码**（`The two surfaces fill the same slots, so switching is a composition change, not a code change`）。⇒ 这一形对我们票 214 第 3 条直接可抄：**Go 侧只保证"拿到路径之后判定链完整"，选法两枚、可换。**

### 1.6 权限与命令的关系（问题 6）

- **命令本身不判权限档**，DSH 的档位是"权限预设"（sandbox mode + approval policy），改它**只有一条写路径，就是 `/permission` 这条命令**：`packages/interaction/permission-presets/src/index.ts:252-256` 注释具名 `The /permission command: the one write path a web client uses`。
- 前端选完档**不是直接调 setter，而是把整行命令回灌**：`ui-permission-presets/src/client/index.ts:121-128` `live.command('/permission ' + preset)`，并且**如果宿主那边根本没有 `/permission` 这条命令就抛错**（`:127` `the host offers no /permission command`）。⇒ **"菜单里有、后端其实不认"这一形被这一句堵死。**
- 严格档下命令会不会被拒：**本轮未读到"按档位拒绝某条命令"的代码**——DSH 的形状是"命令改档位、工具调用受档位管"，不是"档位管命令可见性"。可见性只由 `available()`（会话形态）与组合期挂载（能力）决定。〔已证：没读到 ≠ 没有；见 §6〕
- **有没有"命令能改权限档位"这一形**：**有**，就是 `/permission`。它怎么防被外部内容诱导：
  - `CommandSourceMap` 只有 `user` 一枚变体，注释写明理由：`minimal because every executor caller is a human-facing UI surface dispatching a human-typed line, so the sole variant is user`（`packages/interaction/commands/src/types.ts:68-79`）。⇒ **模型/工具/外部内容根本没有"执行命令"这个入口**，这比"给模型加一条不能调 /permission 的规则"强得多。
  - 命令**不进会话日志的模型面**：`command/run` 是 `Log-only (never model surface)`（`types.ts:96-97`）；`ui-commands/README.md` Model Experience 段：`the command line, the detached result, and every menu and notice rendering stay client-side and never enter the session log`。
  - 内置命令的脸按 `definitionId` 认，改后端描述拿不到别名（§1.3）。

### 1.7 导出／反馈这类"会话出口"（问题 7）

- `/export`：后端注册（`packages/session-query/session-log-export/src/index.ts:78-84`），描述 `Download this Session log as a ZIP archive`；**handler 本身不产文件**，只回 `REQUESTED`，并且**显式拒绝带参数**：`The Web /export command does not accept a path.`（`:81-83`）。真正的 ZIP 走**另一条带鉴权的 HTTP 路由**：`connectionOf(ctx).fetch.register({ path: SESSION_LOG_EXPORT_PATH, methods:['GET','HEAD'], … })`（`:85-90`），压缩级别可配（`config.compressionLevel ?? DEFAULT_…`）。⇒ **"命令只负责请求，产物走一条可鉴权可下载的通道"，两件事不混。**
- `/feedback`：`packages/feedback/command-feedback/src/index.ts:118-125`，描述 `Record feedback about this session`，`input: { hint:'<text>' }`，并且 **`recordInput: false`** —— 因为反馈内容另有权威事件持有，避免在会话日志里重复一份 payload（`index.ts:70-75` 的 `recordInput` 文档＋`types.ts:98-104`）。
- `/compact`：`packages/compaction/command-compact/src/index.ts:101-105`，**无 `input`** ⇒ 前端判定它是"裸 token"命令，带参数不认领（`ui-commands/src/client/service.ts:358-360`）。注册前先 `yield async () => { await Promise.allSettled(active) }` 排空在飞的调用（`:98-100` 注释：组合体是 LIFO 拆除，保证拆除时没有新调用能进来）。
- `/goal` 是唯一声明**收附件**的命令之一：`input: { hint:'[<objective>|clear|edit <objective>|pause|resume]', attachments: true }`（`packages/goal/command-goal/src/index.ts:191-196`）。`/plan` 也收（`packages/plan/plan-mode/src/index.ts:237`），但 `/plan off` 带附件被拒：`Attachments cannot accompany /plan off.`（`:240-242`）。⇒ **子命令级的附件取舍交给 handler 自己判**，这条口径写在 `types.ts:26-29`（`A declaring command's handler … owns every further grammar decision, including rejecting sub-commands that cannot use them`）。

---

## 2. pi-upstream（Pi）——**Step-Code 的地基**，名册由四路汇成一条

⚠ 先说清关系：Step-Code 的 `packages/coding-agent/src/step/slash-commands.ts:1` 第一行注释就是 `Step-only slash command adapters built on Pi's public extension actions.` ⇒ **Pi 是底座、Step-Code 是产品层**。读 Pi 等于读 Step-Code 的下半截，两家不能当两家算。

### 2.1 分类与名册来源（问题 1、3）

- **内置名册＝前端硬编码数组**：`packages/coding-agent/src/core/slash-commands.ts:18-43` 的 `BUILTIN_SLASH_COMMANDS`。
  枚数推导式：`grep -c "{ name:" packages/coding-agent/src/core/slash-commands.ts` ⇒ **24 枚**（settings/model/tree/thinking/scoped-models/export/import/share/bug/copy/name/session/changelog/hotkeys/fork/clone/trust/login/logout/new/compact/resume/reload/quit）。每条只有 `name / description / argumentHint?` 三个字段。
- **名册由四路拼成**，拼装点只有一处：`packages/coding-agent/src/modes/interactive/interactive-mode.ts:683-782` `createBaseAutocompleteProvider()`
  ```
  [...slashCommands(builtin), ...templateCommands(prompt), ...extensionCommands, ...skillCommandList]
  ```
  （`interactive-mode.ts:777`）
- 除内置外，其余三路的**来源标签是类型系统里写死的枚举**：`SlashCommandSource = "extension" | "prompt" | "skill"`（`core/slash-commands.ts:4`）。⇒ **Pi 的"分类"是按来源分，不是按用途分**（与 DSH 的 Add/Commands 两组是两种思路）。
- **技能进名册要过两道闸**（`interactive-mode.ts:759-770`）：
  1. 开关：`if (this.settingsManager.getEnableSkillCommands())` —— 配置关了就**一枚都不进**；
  2. 命名空间：命令名写成 `skill:${skill.name}`（`:762`）⇒ **技能与命令天然不重名**，不需要冲突仲裁。
- **"能力没落地就不进名册"在 Pi 里是三种写法**（比 DSH 更硬）：
  1. 上面的设置开关（`getEnableSkillCommands()`）。
  2. **扩展命令撞内置名 ⇒ 从名册里剔掉**：`interactive-mode.ts:752-755` `.filter((cmd) => !builtinCommandNames.has(cmd.name))`。
  3. **但剔除不是静默**：同一份冲突被 `getBuiltInCommandConflictDiagnostics()` 收集成**警告诊断**（`interactive-mode.ts:668-679`），文案二选一说清后果：
     - `Extension command '/x' conflicts with built-in interactive command. Skipping in autocomplete.`
     - `… Conflicts with built-in interactive command. Available as '/y'.`（有 `invocationName` 改名时可用的另一条路）
     ⇒ **这是我们票 213 AC#2 最该抄的一形：进不了名册要留一句话，而不是悄悄少一行。**
- **注册表本身**（实验线）：`packages/coding-agent/src/experimental/services/slash-commands-provider.ts:24-75`
  - `register()` 重名 throw（`:30`）；`replace()` 允许覆盖，栈式生效（`:34-37`、`:49-63`：只有 `entries[0]` 生效，前面的被 dispose 后依次让位）。
  - **名字文法校验在注册时**：`/^^[a-z0-9][a-z0-9:-]*$/u` 不过就 `throw new TypeError`（`:65-69`）⇒ 允许 `:`，所以 `skill:foo` 这种带命名空间的命令名是**契约允许的形状**。
  - 有 `subscribe(listener)`：名册一变就推给订阅者（`:43-47`），且**订阅时立刻回灌一次当前值**。

### 2.2 `/` 与 `@` 的触发语义（问题 2）

- **触发判定在前端 TUI**：`packages/tui/src/autocomplete.ts`
  - 是不是命令：`prefix.startsWith("/") && beforePrefix.trim() === "" && !prefix.slice(1).includes("/")`（`autocomplete.ts:432`）⇒ **必须行首、且第一根斜杠之后不能再有斜杠**（所以 `/C:/x` 这种绝对路径不会被当命令）。
  - `@` 是**同一个 provider 的另一支**（`CombinedAutocompleteProvider`，`autocomplete.ts:302-308`），文件补全走 `fd`（外部 find-fast 工具）或自扫目录，打分函数 `scoreEntry`：完全等名 100 / 前缀 80 / 文件名包含 50 / 全路径包含 30 / 目录 `+10`（`autocomplete.ts:736-755`），取前 20 条（`:819`）。
- **模糊匹配＝`fuzzyFilter`（子序列）**，且**技能名做两遍匹配**：先按**剥掉 `skill:` 前缀的裸名**匹配，再把没命中的 `skill:` 全名做第二遍（`autocomplete.ts:355-365`）⇒ 用户打 `compact` 也能搜到 `skill:compact`，但内置 `compact` 排在前面。
- **参数怎么传**：命令名之后的整段原样交给 handler（`SlashCommand.getArgumentCompletions(argumentPrefix)`，`autocomplete.ts:257-263`）。**参数补全是每枚命令自己提供的函数**，没有就 `null`。三枚内置命令自带参数补全：`model`（provider/model 模糊，`interactive-mode.ts:691-714`）、`thinking`（档位列表，`:717-729`）、`login`（provider 列表，`:732-743`）。
  - **参数补全与执行是两套判定**：补全用模糊（`createFuzzyAutocompleteItems`），执行用**精确且唯一**：`exactModel()` 里 `matches.length === 1 ? matches[0] : undefined`（`slash-commands-provider.ts:216-225`）⇒ **有歧义＝没匹配上**，然后落到 `ui.select()` 弹选择器，而不是猜一个。
- **解析失败提示**（Pi 的形状，与 DSH 不同）：
  - 参数级：`Unknown model: ${args}`（`:141`）、`Unknown thinking level "${args}". Available levels: ${levels.join(", ")}.`（`:180`）——**抛错、带可用值**。
  - 服务没就绪：`Models service is not ready`（`:138`）——**不静默降级**。
  - 命令级：`dispatch` 找不到就不处理；**没有"未知命令"的统一红字**，未识别的 `/xxx` 会作为普通文本发出去（与 openchamber 同形，见 §4.2）。⚠ 这条只说明"Pi 这么做"，**不能写成"三家都这样"**：DSH 明确拒绝降级（§1.2）。

### 2.3 来源可信度在菜单里怎么说（问题 4，Pi 的最强项）

- **每条非内置命令的描述前面，直接拼一方标签**：`interactive-mode.ts:660-666`
  ```ts
  return description ? `[${sourceTag}] ${description}` : `[${sourceTag}]`;
  ```
  三处调用点：prompt 模板（`:747`）、扩展命令（`:757`）、技能（`:766`）。⇒ **用户在看的那一行就带着"这是谁给的"**，不需要点开详情。
- 标签的数据源是 `SourceInfo`：`{ path, source, scope: "user"|"project"|"temporary", origin: "package"|"top-level", baseDir? }`（`packages/coding-agent/src/core/source-info.ts:3-11`）。⇒ **scope 这一维直接区分"用户自己装的"与"这个仓库带进来的"**，这正是票 215 第 3 条要的"来源可信吗"。
- 每条命令还带 `sourceInfo` 一路传到执行上下文（`RegisteredCommand.sourceInfo`，`packages/coding-agent/src/core/extensions/types.ts:1275-1281`）⇒ **审计能追到是谁注册的**。
- 技能/命令的**生效与否不在菜单里显示**，而是看会话里那一枚工具行的状态（与 DSH 同思路）。〔已证：Pi 的菜单行只有 name/description/argumentHint 三字段，`autocomplete.ts:257-262`〕

### 2.4 会话出口（问题 7）

- `/export`：`Export session (HTML default, or specify path: .html/.jsonl)`（`core/slash-commands.ts:25`）⇒ **默认 HTML，路径后缀决定格式**；HTML 那台机器在 `packages/coding-agent/src/core/export-html/`（`index.ts` + `template.html/css/js` + `ansi-to-html.ts` + `tool-renderer.ts` + `vendor/`）。
- `/import`：`Import and resume a session from a JSONL file`（`slash-commands.ts:26`），用法错时 `this.showError("Usage: /import <path.jsonl>")`（`modes/interactive/interactive-mode.ts:6340`）。
- `/share`：`Share session as a secret GitHub gist`（`slash-commands.ts:27`）⇒ **出口是"发到哪"写死在描述里的**，用户不用猜。
- `/copy`：`Copy last agent message to clipboard`（`:28`）。
- 另有 `/bug`（`Report a bug to the Pi developers`，`:26`），实现在 `core/bug-report.ts` + `core/bug-report-upload.ts`。

---

## 3. Step-Code——Pi 的产品层：只有 6＋1 枚自有命令，但**出口那一块做得最细**

### 3.1 名册与分层（问题 1、3）

- Step 自有的产品命令只有 **6 枚**，推导式：`grep -n 'registerTrackedCommand(' packages/coding-agent/src/step/slash-commands.ts | wc -l` ⇒ **6**（`mcp` / `clear` / `exit` / `theme` / `status` / `feedback`，`step/slash-commands.ts:47,59,74,86,104,116`），加 1 枚插件命令 `registerStepPluginCommand`（`step/slash-commands.ts:26` 引入，`:43` 调用）。
- **它们全部走 Pi 的公开扩展口注册**，不开后门：`step/slash-commands.ts:1` 注释 `Step-only slash command adapters built on Pi's public extension actions.`；`:31-36` 注释把纪律写死了：`The handlers intentionally stay at the public extension boundary. In particular, they do not call InteractiveMode methods or maintain a second session/theme state machine. Pi remains responsible for session replacement, shutdown, selector presentation, persistence, and rendering.`
  ⇒ **这一条对我们极重要**：产品层不加第二台状态机。我们票 213 的"命令名册"如果各模块自己画自己那一份，就是违反这条。
- **分类**：Step 没有自己的分类层，沿用 Pi 的"来源分类"（§2.1）。⇒ 问题 1 在 Step-Code 这一家**给不出独立答案**。

### 3.2 `/mcp` 这一枚是"状态回显"的样板（问题 4）

- `packages/coding-agent/src/step/slash-commands.ts:45-55`：描述 `Show configured MCP servers and loaded tools`，handler 只有一行 `ctx.ui.notify(formatStepMcpStatuses(), "info")`。
- ⚠ **注意它的措辞**：**"配置的服务器"与"已加载的工具"是两件事**，所以这一枚命令的存在本身就承认"配了 ≠ 起来了"。格式化函数在 `packages/coding-agent/src/step/mcp.ts`（该文件同时持有连接池与 `ConnectedServer { name, client, transport, tools, callTimeoutMs }`，`mcp.ts:36-42`）。⇒ **这一形值得抄进票 215 的清单页：一行里同时给"配置里有"与"实际连上没/加载了几枚工具"。**
- 超时是**显式常量**：`MCP_STARTUP_TIMEOUT_SEC = 30`、`MCP_CALL_TIMEOUT_SEC = 300`、每服务器可覆写 `tool_timeout_sec`（`mcp.ts:27-29, 41`）。

### 3.3 插件市场：声明式清单＋"卸了别自己回来"（问题 3、4）

- `packages/coding-agent/src/step/plugins.ts:1-9` 注释是这一层的契约：`marketplace packages are declarative manifests copied into .stepcode/plugins. MCP processes are started by the Step runtime after installation; this module owns discovery and provisioning only.` ⇒ **安装＝拷清单，起进程是另一枚模块的事**，两件事不混。
- **防"预装把用户卸载的东西偷偷装回来"**：`PREINSTALL_MARKER_FILE = ".stepcode-preinstalled"`，注释具名：`Records which built-in plugins have already been auto-installed, so a plugin the user later uninstalls is not silently resurrected on the next launch.`（`plugins.ts:38-42`）⇒ **这正是票 215 AC#2 那一族"状态被悄悄改回"的形，而且它是一条文件名就能解决的。**
- 清单有**大小与命名上限**：`SAFE_NAME = /^[a-z0-9][a-z0-9._-]*$/iu`、`MAX_MANIFEST_BYTES = 512 * 1024`（`plugins.ts:43-44`）。
- 内置市场**不联网取**：`description: "Plugins that ship inside the StepCode binary. Updated with the CLI itself, not fetched."`（`plugins.ts:50`）⇒ 与我们"不许从镜像站取哈希"同一条纪律的另一种落法。

### 3.4 会话出口：`/feedback` 是五家里最完整的一枚（问题 7）

- 目录 `packages/coding-agent/src/step/feedback/` 共 **15 个文件、3101 行**（推导式：`ls packages/coding-agent/src/step/feedback | wc -l` ⇒ 15；`wc -l packages/coding-agent/src/step/feedback/*.ts | tail -1` ⇒ 3101）。
- **产物**：gzip 归档，内含一枚 `manifest` 条目（JSON，`bundle.ts:125-137`）＋会话事件流＋dev log。
- **上限是分层回落的**：`EVENTS_MAX_BYTES = 24 MiB`、`DEV_LOG_MAX_LINES = 2000` / `DEV_LOG_MAX_BYTES = 512 KiB`、`SESSION_HEADER_MAX_BYTES = 1 MiB`（`bundle.ts:19-26`），总上限 `FEEDBACK_BUNDLE_MAX_BYTES`；超了先 `limitRedactedSessionEntry(...)` 缩一次，**再超才放弃**：`return { status: "skipped", reason: "too-large" }`（`bundle.ts:88-97`）。⇒ **拒的时候有一个可枚举的 reason，不是一句"失败"。**
- **先脱敏、再限长**，且截断要留话：`redacted.byteLength <= maxBytes ? ready : tailAtLineBoundary(...)`，并在条目的 `note` 里写明是哪种截断：`omitted after redaction, a single line exceeds ${formatBytes(maxBytes)}` / `tail re-limited after redaction, ${formatBytes(...)} before re-limit`（`bundle.ts:176-185`）。⇒ **"截断了"必须带"怎么截的、还能不能找回来"。**
- **同意页要先把清单摆出来**：`bundle.ts:110` 注释 `File manifest shown before the user consents to sending the conversation.`，每行 `name (N bytes) - note`（`:113`）。
- **同意页防被伪造**（这一形我们完全没有）：`packages/coding-agent/src/step/feedback/consent.ts:73-92` 把**所有 C0/C1/DEL 控制字符在套主题 ANSI 之前**换成可见转义（`\t`、`\r`、`\b`、`\x1b`、`\xNN`），并且**对每一个来自外部的字段都套一遍**：`sessionId`、每个文件名 `file.name`、每条 `file.note`、诊断的展示路径、以及 identity 的 7 枚字段（`consent.ts:38, 45, 61-64`），最后一道 `neutralizeFeedbackConsentText(sections.join("\n\n"))`（`:70`）。
  ⇒ **口径**：任何"要被显示给人看、而内容不是人自己打的"字符串，都要先剥掉控制字符再显示，否则一枚文件名就能改写那页同意书。这条对我们票 214 AC#5、票 215 AC#4 都是直接可落的判据。

---

## 4. minimax-code——**名册字段最完整的一家**（也是"配置 true 但没生效"这一形答得最好的一家）

minimax 是 TUI（`packages/tui`），没有 GUI 的加号按钮；但它的**命令目录**与**技能登记表**是五家里字段最全的，恰好是我们票 213/215 需要的数据形状。

### 4.1 命令名册：一枚 843 行的目录，字段就是设计（问题 1、3）

- 文件：`packages/tui/src/tui/commands/catalog.ts`（843 行）。枚数推导式：`grep -c "^    name: '" packages/tui/src/tui/commands/catalog.ts` ⇒ **41 枚**；其中 `visibleWhen:` **18 枚**、`unavailableReason:` **19 枚**（同一条 `grep -c`）。
- **分类是前端硬编码的 7 值枚举**：`TuiCommandCategory = 'Session' | 'Runtime' | 'Capability' | 'Input' | 'Transcript' | 'Decision' | 'Application'`（`catalog.ts:8-15`）。⇒ 问题 1 的答案：**写死**，后端不参与（minimax 这一层根本没有后端目录，见 §4.1 末）。
- **一枚命令的元数据字段**（`catalog.ts:47-84`），逐个都是我们要的判据形状：
  | 字段 | 口径（原文注释具名） | 我们该不该抄 |
  |---|---|---|
  | `audience: 'user'\|'internal'` | `Internal controls stay reserved for parser compatibility … but never appear in slash discovery. A reserved command without a handler is treated as unrecognized.` | **抄**：这就是"能力没落地就不进名册" |
  | `discoverability: 'primary'\|'contextual'\|'search-only'` | 三档可见度 | 抄 |
  | `readiness: 'immediate'\|'controller'\|'full'` | 装配阶段 | 可缓 |
  | `runAvailability: 'idle'` | `Commands that replace or mutate the active Session are unavailable while a Turn is live.` | **抄**（我们有一模一样的形） |
  | `visibleWhen(context)` | 谓词，context 见下 | **抄** |
  | `unavailableReason: string` | 一句话说人话 | **必抄**（票 213 AC#3） |
  | `argumentHint` | `Its presence also declares that the command accepts non-whitespace text after its token; commands without a hint are exact-token commands.` | **必抄**：语法声明即元数据 |
  | `usage` / `composerTemplate` | 由 `materializeCommand` 派生（`catalog.ts:718-733`） | 抄 |
  | `invocationKind?: 'skill'` | `Runtime-provided Skill invocations carry Agent instructions after their slash token.` | 抄 |
- **可见性上下文是一枚结构体**，不是一堆散参数：`TuiCommandContext { hasSession, hasParentSession, managedTokenPresent, queueEnabled, hasLiveRun, hasPendingInteraction?, queuedCount, canRetry, sideMode? }`（`catalog.ts:18-31`）。⇒ **我们要落 Go，这正好是一张可以照抄成 struct 的表**，每一字段在 Wisp 侧都有对应物（有没有会话／有没有在跑的任务／队列长度）。
- **名册来源**：硬编码目录 `DEFAULT_COMMAND_CATALOG = createTuiCommandCatalog()`（`catalog.ts:775`）＋插件贡献（`TuiCommandContribution { id, order, kind: 'command'|'active-run', source, execute }`，`catalog.ts:106-113`；注册表 `TuiContributionRegistry` 从 `catalog.ts:3` 引入）。
  **重名处理**：`materializeCommand` 里对 name＋全部 aliases 逐个查集合，撞名直接 `throw new Error('Duplicate command name or alias: ' + name)`（`catalog.ts:718-725`）⇒ **连别名都不许撞**。
- ⚠ **口径注意**：minimax 这一层是 TUI 本地目录，**没有 DSH 那种"后端 RPC 给目录"**。它的"后端"是 `packages/agent-modules/*`（枚数推导式 `ls packages/agent-modules | wc -l` ⇒ **12**）。

### 4.2 `/` 与 `@` 的解析：**同一枚解析器同时供"界面判断"和"回车执行"**（问题 2，本票最该抄的一形）

- 解析函数：`matchTuiCommandInput(input, commands)`（`catalog.ts:698-717`）。头注释（`:692-696`）就是设计声明：
  > `Parse the slash token and trailing payload once for both execution and live Composer intent. Exact-token commands reject non-whitespace suffixes, so callers cannot visually classify input differently from the command path that will handle Enter.`
  ⇒ **界面高亮/分类与回车执行共用同一枚函数**，从根上杜绝"看着像命令、按下去不是"。
- 语义细节：`raw = input.trim()` → `token = raw.split(/\s+/u, 1)[0]` → **`/` 与 `@` 两种 sigil 走同一支**（`:704`）→ 名字 `toLocaleLowerCase()` 后与 `name + aliases` 逐个精确比（`:706-711`）→ `args = raw.slice(token.length).trim()` → **`if (args && !source.argumentHint) return undefined`**（`:714`）。
  ⇒ **匹配是精确等名，不是模糊**；模糊只在"搜索"里做（见下）。
- **搜索/模糊匹配**：`filterTuiCommands`（`catalog.ts:788-830`），口径与 DSH 完全不同：
  - 查询串按空白切成 token，`tokens.every(token => searchable.includes(token))`，`searchable` 是 `name + aliases + usage + description + category + shortcut` 拼起来的小写串（`:813-829`）⇒ **多关键词 AND 命中元数据**，不是子序列。
  - **但带 `/` 的单 token 只准撞名字前缀**：`directNameMatches`（`:806-812`），且 `if (slashQuery && tokens.length === 1) return []`（`:812`）⇒ **打了 `/xyz` 而没撞任何命令名前缀，就明确给空列表，不会拿描述来凑数。**
- **解析失败/不可用的三态**（`dispatch`，`catalog.ts:627-663`）：
  - 认不出来 ⇒ `{ status: 'unrecognized' }`（`:629`；没有 handler 的 reserved 命令也落这里，`:661-662`）
  - 认出来但当前不可用 ⇒ `{ status: 'unavailable', reason }`，reason 三选一：
    1. 侧会话白名单外：`resolveSideModeUnavailableReason` ⇒ `${usage} is unavailable in side conversations. Press Ctrl+C to return to the main session first.`（`:760-766`）——**拒的时候连"怎么退出这个状态"一起给**
    2. 有活着的 Turn：`Stop the running turn before using ${usage}.`（`:651-652`）
    3. 否则用命令自带的 `unavailableReason`，兜底句 `${usage} is unavailable right now.`（`:653`）
  - 成功 ⇒ `{ status: 'handled', disposition }`
- **两道闸的分工**（很值得抄）：
  - `resolveTuiCommandVisibility`（`:738-747`）把不可用的**从菜单里直接滤掉**（`if (!isTuiCommandAvailable) return false`），`search-only` 的只在搜索时出现；
  - `dispatch`（`:640-656`）在**手打**那条路上仍然拒并给原因。
  ⇒ **"菜单里看不见"与"打了会被拒"是两码事，两套代码各有各的说法，不会一边放行一边拒绝。**

### 4.3 技能：`winners / losers / diagnostics / metrics` 四件套（问题 4，**票 215 AC#2 的直接答案**）

- 文件：`packages/agent-modules/skills/src/types.ts`（124 行）＋ `registry.ts`（1099 行）。
- **来源是六分枚举**：`SkillSourceKind = 'project' | 'workspace' | 'agent' | 'global' | 'user' | 'builtin'`（`types.ts:1`），每枚 root 带 `priority?: number`、`external?: boolean`、`allowDirectorySymlinksOutsideRoot?: boolean`（`types.ts:3-12`）。
- **每条技能自带出处**：`SkillEntry` 有 `rootId / rootKind / rootScope / rootPriority / sourceExternal / locationUri / skillDir / entryDir?`（`types.ts:24-47`），其中 `entryDir` 的注释具名：`Root-owned directory entry; differs from skillDir for linked skills.` ⇒ **软链进来的技能能一眼看出来。**
- **"配了但被顶掉"是一等公民**：
  ```ts
  export interface SkillLoser { entry: SkillEntry; reason: string; winnerLocationUri: string }   // types.ts:49-53
  export interface SkillViewEntry extends SkillEntry { losers: SkillLoser[] }                     // types.ts:55-57
  export interface SkillSnapshot { version; generatedAt; roots; entries; winners; losers; diagnostics; metrics }  // types.ts:59-68
  ```
  ⇒ **同名技能谁赢了、为什么输、赢的那枚在哪，是快照里的字段，不是靠人复现。** 这一形直接对应票 215 AC#2 的正控（"塞一条配置 true、装配里没人读的假项 ⇒ 判据必红"）。
- **诊断也是字段**：`SkillDiagnostic { level:'warning'|'error', code, message, rootId?, locationUri? }`（`types.ts:14-22`）。
- **还有第三种"没生效"：被上下文预算挤掉**，而且是**可数的**：
  ```ts
  SkillRefreshMetrics { rootsScanned, filesSeen, filesRead, filesReused, diagnostics, entries, winners, losers }   // types.ts:70-79
  SkillRenderMetrics  { total, rendered, compacted, dropped, charsUsed, budgetChars }                              // types.ts:86-93
  ```
  ⇒ **"扫到 30 枚、渲染 22 枚、压成摘要 5 枚、丢掉 3 枚、预算 8000 字用了 7900"** 这一行就是"配置 true 但这次没生效"的完整答案。
- **整表可 dump 成 JSON 供事后核**：`SkillRegistryDump`（`types.ts:100-115`），字段名直接是 snake_case 的 `refresh_metrics` / `render_metrics`。
- ⚠ 三态判定：以上全部〔已证〕是类型与登记表实现；**这些 dump 在 TUI 里画给谁看、有没有画，本轮未读到**（`packages/tui/src/tui/features/inspection/` 未展开）。

### 4.4 附件与"终端注入"防护（问题 5，兼问题 6）

- 文件：`packages/tui/src/tui/features/composer/attachments.ts`（170 行）。
- **类型判定＝按扩展名查表，查不到再猜视频，再兜底 `application/octet-stream`**：`MIME_BY_EXTENSION[extname(fileName).toLowerCase()] ?? inferTuiNativeVideoMimeType(fileName) ?? 'application/octet-stream'`（`attachments.ts:90-93`），`type` 由 `mimeType.startsWith('image/')` 决定（`:95`）。
  ⚠ **注意这与 DSH 不同**：minimax 这里**没有内容嗅探**，纯按扩展名。我们票 214 AC#3 要求"按嗅探结果判、不按扩展名判"，**minimax 这一家不构成支持证据**（我们自己的 `matchesISOBaseMedia` 才是）。
- **不是文件就当场抛，且句子里带被 sanitize 过的路径**：`Cannot attach ${sanitizeTerminalText(normalizedReference)}: path is not a file.`（`attachments.ts:85-87`）。
- **附件清单的展示格式是固定三段**：`N. 文件名 · mime · 人类可读字节数`（`attachments.ts:105-108`）⇒ 票 214 第 2 条要的"文件名、类型、字节数"，这一行就是最小形状。
- **文件名先过终端消毒再显示**：`sanitizeTerminalText`（`packages/tui/src/tui/rendering/terminal-text.ts:15-26`）——**先吃掉 OSC/DCS/SOS 这类终端串（含不完整的前缀与 C1 形态）**，再 `stripVTControlCharacters`，再把 `\r` 换成空格、把 C0/C1/DEL 全删。
  ⇒ 与 Step-Code 的 `neutralizeFeedbackConsentText`（§3.4）是同一族防护，**两家各写了一份、写法不同**（minimax 是"删/替换"，Step 是"转成可见转义"）。对我们：**TUI/WebView 两种宿主各需要哪种，得单独定，不能一句"消毒过了"糊过去。**
- **预览失败不等于附件失败**：中文文案直接写清（`packages/tui/src/tui/features/composer/copy.zh-Hans.ts:5`）`imagePreviewUnavailable: '无法预览 · 附件仍可发送'`。⇒ 票 214 第 2 条"不许只报成功"的镜像形（也不许把"预览不出来"报成"没附上"）。
- **权限与命令的关系（问题 6）**：minimax 的"档位"这一形体现在 `SIDE_MODE_READ_ONLY_COMMANDS`（`catalog.ts:39-50`，9 枚白名单：help/changelog/context/status/usage/export/transcript/copy/parent），注释具名口径：`Side mode follows Codex's read-only command surface. Navigation, mutation, auth, and process controls stay unavailable even when typed directly.` ⇒ **白名单而不是黑名单，且"直接打字也不行"**。另有 `managedTokenPresent` 这一枚 context 字段（`catalog.ts:21`）说明"凭证从哪来"会改变命令可见性。
  ⚠ **本轮未读到 minimax 有"命令改权限档位"这一形**（`agent-modules/permission/` 未展开）。

---

## 5. openchamber——**加号菜单最像 owner 截图的一家**，但它把"名册权威"写成了一条不变量

openchamber 是 GUI（`packages/ui`，React＋CodeMirror），有真正的 `+` 按钮与下拉。

### 5.1 加号里放什么（问题 1、5）

- 载体：`packages/ui/src/components/chat/composer/ui/ComposerAttachmentControls.tsx`。按钮的 `aria-label` 是 `chat.chatInput.actions.addAttachment`（`:81-82`），图标 `attachment-2`（`:94`）。
- **菜单项固定四段，后两段是条件渲染**（`:91-143`）：
  1. `attachFiles`（本地文件）
  2. `linkGithubIssue` / `linkGithubPr`
  3. `linkLinearIssue` —— **仅当 `showLinearPicker && openLinearPicker` 都给了才画**（`:133`）
  4. `attachGuests?.map(...)` —— **扩展（guest）贡献的 attach 项**（`:143`）
  ⇒ 分层是**写死的顺序**，成员是**能力位＋扩展贡献点**（VS Code 式 `contributes.attach`）。
- 桌面画下拉、移动端画底部抽屉：`/** Mobile: open the attachment bottom sheet instead of the dropdown menu. */`（`:35`）。
- **VS Code 与移动端跳过这张列表**：`packages/ui/src/components/chat/composer/DOCUMENTATION.md:125` 末句 `VS Code and mobile skip that list.` ⇒ **同一个产品里，加号里有什么是随宿主变的**，不是随代码版本变的。
- **扩展贡献点三类**：`contributes.attach`（加号项）、`contributes.commands`（斜杠命令）、`contributes.actions`（消息/会话动作），全部走"panel 轨道"或"dialog 轨道"两条路（`DOCUMENTATION.md:125`）。
- **附件不止文件**：`ComposerAttachmentControls` 之外还有 `ComposerContextChips.tsx`（上下文芯片）、`LinkedReferenceRow.tsx`。guest attach 回来的东西**保留一枚不透明 `data`**，注释写明：`The chip keeps the guest's opaque data … so it comes back byte-identical; it is never part of the context text.`（`DOCUMENTATION.md:125`）⇒ **"引用的元数据"与"喂给模型的文本"是两条线，不混。**

### 5.2 `/` 与 `@`：**"名册才是权威，正则只是定位器"**（问题 2，本条是五家里最锋利的一句）

- 一条被写成模块级不变量的话（`DOCUMENTATION.md:150-161`）：
  > `language/` is the single source of truth for composer syntax. … **This is the invariant that matters most in this module.** Before it existed, the `@` rule was written four times with divergent cleanup and the `/` rule three times with different valid character sets, so a token could be painted as a reference and then not resolve as one.
- 这条不变量换来的**具体历史 bug**（`language/prefixTokens.ts:1-15` 头注释）：
  > The send-time skill scanner, for instance, accepted only lowercase names, so a `/My_Skill` token was painted as a command but never collected.
- **判定口径**（`prefixTokens.ts:11-14`）：`Scanning is deliberately generous: it finds every syntactically plausible token and leaves the decision of what exists to the caller, which holds the authoritative set … Membership is the authority; the pattern is only a locator.`
  同一句话在 `DOCUMENTATION.md:167-169` 再写一遍：`**membership in the command, skill or snippet registry is the authority**, not the pattern. An unknown /token stays plain prose.`
  ⇒ **解析失败的用户可见结果＝那一行保持普通文本，不弹提示。** 这与 DSH（§1.2，明确报错、绝不静默降级）**是相反的选择**。⚠ 我们票 213 AC#3 要"拒了要说为什么"，所以 **openchamber 这一形不能当"业界共识"引**；它是"不拒、只是不当命令"。
- 标识符文法集中一处：`TOKEN_NAME = '[A-Za-z0-9][A-Za-z0-9_-]*'`，`/` 与 `#` 共用（`prefixTokens.ts:22, 41-44`），注释：`Kept in one place so / and # cannot drift apart again.`；且**必须在词边界**：`so a/b and #1 inside issue#1 stay ordinary prose`（`:51-53`）。
- **四支选择器、恰好一支激活、优先级写死**：`language/triggers.ts:11-14` ⇒ `command > skill > snippet > mention`。四支各自的规则：
  - 命令面板：`/` 必须在**第 0 列**、**还没打空格**、光标还在命令词内（`triggers.ts:48-62`，注释：`Once a space appears the message is a command invocation, not a search.`）
  - 行内 `/skill`、`#snippet`：光标前最近的那枚 sigil、在词边界上、与光标之间无分隔符（`triggers.ts:66-74`）
  - **shell 模式（`!cmd`）把四支全关掉**：`TriggerContext.inputMode` 注释 `Shell mode (!cmd) disables every picker.`（`triggers.ts:30-31`）
- **参数**：guest 命令的参数是"打了空格之后的原文，逐字交给扩展"（`composer/submit/guestCommands.ts:12-25`：`The match is on the first word after the slash, the rest is the argument string handed to the guest as typed.`）。

### 5.3 名册来源与冲突处理（问题 3）

- **三源并一名册**：`packages/ui/src/components/chat/ChatInput.tsx:783-791`
  ```ts
  const names = new Set<string>([ 'init','review','undo','redo','timeline','compact','fork','btw','summary','workspace-review','plan-feature','craft-goal','schedule-task','catch-up','debug','weigh','explore' ]);
  if (!isMobile && !isVSCodeRuntime()) names.add('handoff-review');
  for (const command of availableCommands) names.add(command.name.toLowerCase());
  for (const skill of availableSkills) names.add(skill.name.toLowerCase());
  ```
  枚数推导式：`ChatInput.tsx:786` 那个字面量数组 **17 枚**（`sed -n '786p' ChatInput.tsx` 逗号计数），再按宿主条件加 1 枚（`handoff-review`，`:788`）。
  ⇒ **内置是硬编码；命令与技能来自 store（`useCommandsStore` / `useSkillsStore`，按工作目录选：`selectCommandsForDirectory(s, currentDirectory)`）**——**"哪一目录下的命令"是名册的一维**，这一维五家里只有 openchamber 有。
- **扩展命令撞内置 ⇒ 直接忽略（不报错）**：`ChatInput.tsx:793-795` 注释 `Extension slash commands. Built-ins, OpenCode commands, and skills are reserved: an extension command with one of those names is ignored.`；`guestCommands.ts:1-6` 补一句 `Planned before local commands so an extension name the composer already uses can never take over.`
  ⚠ **这与 Pi 的诊断形（§2.1）正好相反**：Pi 会出警告，openchamber 静默丢。**两形并存 ⇒ "扩展命令撞名"没有业界标准做法，我们必须自己选一种并写进票面。**
- **"能力没落地就不进名册"的机制**：`requires: CommandRequirement`（`composer/submit/slashCommands.ts:19`，只有 `'session' | 'session-or-draft'` 两值）＋每条命令自带 `errorToastKey`（`:26`，注释 `i18n key for the toast shown when the command fails`）。
  实例：`/summary` 的 `requires: 'session'` 旁边一句注释 `// Summarizing needs a conversation to summarize.`（`slashCommands.ts:66-67`）。
- **命令是数据、不是 else-if 链**：`slashCommands.ts:1-13` 头注释：
  > Most of them do the same thing: render a pair of magic prompts — one the user sees, one the model is instructed with — and send them as a single message. That shape was previously written out nine times as an `else if` chain, so adding a command meant copying twenty lines … **Here the shape is the executor and the commands are data.**
  ⇒ **openchamber 的多数"命令"根本不是动作，而是"给用户看一句话 ＋ 偷偷给模型一段指令"的成对模板**（`visiblePrompt` / `instructionsPrompt`，`:22-24`）。参数用 `buildVariables` 织进两半里（`:31-56`，`/summary rate limiting` 的例子）。
  ⚠ **这一形对我们是个警报**：票 213 明写"命令≠工具，别把 `/compact` 做成模型可自调的工具"；openchamber 恰恰是**把命令实现成"往这次提交里塞一段模型指令"**。这是"面板侧来源给模型下指令"的另一面，**不能照抄**，只能当反面参照。

### 5.4 附件被拒时给用户看什么（问题 5）

- **上限判在服务端，客户端靠错误文本嗅探回一句本地化提示**：`ChatInput.tsx:2205-2206`
  ```ts
  if (normalized.includes('payload too large') || normalized.includes('413') || normalized.includes('entity too large')) {
      toast.error(t('chat.chatInput.toast.attachmentsTooLarge'));
  ```
  同一枚键的德文文案（形状参照）：`Anhänge sind zu groß zum Senden. Bitte versuche, die Anzahl oder Größe der Bilder zu reduzieren.`（`ui/src/lib/i18n/messages/de.ts:2240`）⇒ **告诉你减什么、减哪一枚，不是只说"太大"。**
  ⚠ **这是反面教材**：判据挂在**英文错误串**上，服务端换个措辞这行提示就永远不弹。我们落 Go 应当用**错误码**，不是措辞。
- **正面对照（同仓另一处用码不用串）**：`packages/ui/src/components/sections/extensions/ExtensionsPage.tsx:58` `if (code === 'too-large') return 'settings.extensions.toast.zipTooLarge';`
- **客户端也有一道前置闸，且对 base64 长度也判**：`packages/ui/src/components/chat/markdown/markdownImageAssets.ts:89` `if (blob.size > MAX_MARKDOWN_IMAGE_BYTES) throw new Error('Image is too large')`；`:96` 对 base64 串长度按 `Math.ceil(MAX * 4/3) + 4` 预判 ⇒ **解码前就拦，不先把 20MB 解成 Uint8Array**。
- **另一族"太大"是按行数判的**：`packages/ui/src/lib/contextFileOpenGuard.ts:101` `return t('contextFileOpen.failure.tooLarge', { count: lines })`，德文串 `Datei ist zu groß zum Öffnen (>{count} Zeilen)`（`de.ts:2991`）⇒ **句子里带那个数字本身。**
- **大段粘贴自动转附件**：粘贴 ≥ 约 2000 字符或 25 行 ⇒ 按设置 `largeTextPasteBehavior`（`ask`/`attach`/`inline`）处理，attach 时造一枚内存 `text/plain` 命名 `pasted-context-N.txt`，**走与手选 `.txt` 完全相同的附件管线**（`DOCUMENTATION.md:137-148`）。
  并且 **ask-toast 的动作读的是实时状态**：`Ask-toast actions read live composer/attachment state so typing or other attaches between paste and choice stay consistent.`（`:141-142`）⇒ 防"toast 还挂着、内容已经变了"这一形。
- **文件与引用同一条路**：`Pasted and dropped files share attachFilesWithCitation: every file attaches and is cited in the draft as [name]; images get a generated unique name first, other files keep their own name and are cited only after they attached.`（`DOCUMENTATION.md:131-136`）⇒ **票 214 第 4 条"不许开第二条读文件的路"，这一家就是这么做的。**

### 5.5 权限与命令（问题 6）

- 加号旁边那一枚是**独立的权限按钮**，不是菜单项：`packages/ui/src/components/chat/composer/ui/PermissionAutoAcceptButton.tsx`（与 `FocusModeButton.tsx`、`DraftTargetSelectors.tsx` 并列，推导式 `ls packages/ui/src/components/chat/composer/ui`）。
- ⚠ **本轮未读到 openchamber 有"斜杠命令改权限档位"这一形**，也未读到"严格档下某条命令被拒"的代码。⇒ 问题 6 在 openchamber 这一家**给不出答案**，不编。

---

## 6. 五家横向对照（七个问题 × 五家，一张表收口）

> 空格＝**这一家本轮没读到对应机制**，不是"它没有"，也不是"它一定有"。

| 问题 | DSH | pi-upstream | Step-Code | minimax-code | openchamber |
|---|---|---|---|---|---|
| 1 分类来源 | **两组名硬编码**（Add/Commands），成员来自后端目录 | 按**来源**分（builtin/extension/prompt/skill） | **7 值 category 枚举硬编码** | | **4 段固定顺序硬编码**，后两段条件渲染 |
| 2 谁解析触发 | 前端纯函数 `detectTrigger`；命令名后端 `parseCommand` | 前端 TUI `autocomplete.ts` | 沿用 Pi | **同一枚解析器供界面与回车** | 前端 `language/` 一处，**名册才是权威** |
| 2 失败提示 | **报错，绝不降级成普通消息** | 参数级报错带可用值；命令级不拦 | 沿用 Pi | 三态 `unrecognized/unavailable/handled`＋人话 reason | **不提示，留作普通文本** |
| 3 名册来源 | 后端 RPC 目录＋前端贡献，撞名 throw | 硬编码 24＋扩展＋prompt＋skill | 沿用 Pi＋自有 6 枚 | 硬编码 41＋插件贡献，撞名 throw | 硬编码 17＋命令 store＋技能 store＋guest |
| 3 没能力就不进册 | ①组合期 `ctx.inject` ②运行期能力位 ③逐会话 `available()` | 设置开关＋撞内置名剔除（**出警告诊断**） | 沿用 Pi | `audience:'internal'`＋`visibleWhen` 双重 | 扩展撞名**静默忽略** |
| 4 技能显示什么 | 只有"仅用户可调"前缀；**无开关状态** | **描述前拼 `[来源]` 标签** | `/mcp` 回显"配置了 vs 加载了" | **winners/losers/diagnostics/四组 metrics** | 名册里只取 name |
| 5 选文件 | 原生 `<input type=file>`；目录选择**两枚可换包** | `@` 走 `fd` 或自扫 | 沿用 Pi | 路径直读，**扩展名定 MIME** | 原生选择器＋拖拽＋大段粘贴转附件 |
| 5 上限判在哪 | **两层都判，后端权威**（`input.attachments` 声明位） | | | 只有 `sizeBytes`，无拒绝分支〔未读到〕 | **服务端 413＋客户端嗅错误串** |
| 6 命令改档位 | **有**（`/permission` 是唯一写路径）；命令执行者类型只有 `user` | | | 未读到 | 未读到 |
| 7 出口 | `/export` 只回 REQUESTED，ZIP 走独立鉴权路由 | `/export` HTML/JSONL、`/share` gist | `/feedback` gzip＋manifest＋先脱敏后限长＋同意页消毒 | | 未读到 |

**我们的现量（同表右端）**：

| 维度 | 现量读数 | 三态 |
|---|---|---|
| 命令名册 | `grep -rln --include=*.go -iE "slashcommand\|commandcatalog\|CommandRegistry\|\"/compact\"\|/goal" internal/ cmd/` ⇒ **无输出，exit=1** | 没有〔已证〕 |
| 附件受理器 | `internal/panel/attachments.go` **410 行**；`grep -rn --include=*.go "\.Ingest(\|AttachmentBroker" internal/ cmd/` ⇒ **命中只在 `attachments.go` 与 `attachments_test.go`，生产 0 枚** | **建了但没接**〔已证〕 |
| 输入框载体 | `internal/panel/composer.go:235-257` `ComposerState` 共 **9 个字段**（mode/workspace/attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError/git/currentModel/modelKnown，推导式：`sed -n '235,257p' internal/panel/composer.go \| grep -c "json:"` ⇒ 9）；生产装配点 `internal/panel/pump.go:226` `NewComposerState(mode, workspace, nil, maxAttachment)` —— **附件那一枚实参是字面量 `nil`** | **字段在、值永远为空**〔已证〕 |
| 名册可读性 | `PumpSources`（`internal/panel/pump.go:111-141`）有 Verdicts/Mode/Workspace/Git/Model/Results/Instructions/Tasks 八类读者，**没有 Attachments 读者、没有 Commands 读者** | 没有〔已证〕 |
| "写了不管用"的旗 | `internal/config/unwired.go` **147 行**，`grep -c "path:"` ⇒ **6 枚响亮拒收键**；文件头三条件口径：`Such a key lies: the user writes one line and the program swallows it.` | **有，且比五家都硬**〔已证〕 |

---

## 7. 我这腿没读到什么（具名到包/目录）

1. **DSH `packages/client/ui-plugin-manager` / `ui-settings-plugins` / `ui-settings-plugin-inventory`**：只确认它们存在与命名，**没读内部**（插件的开/关/版本/来源在设置页里到底显示哪几列，未取证）。⇒ 问题 4 的"插件那一半"我只答了"不进加号菜单"，没答"设置页里显示什么"。
2. **DSH `ui-reference`（`@` 源）内部**：只读到它注册了 `@` 源（`packages/client/ui-reference/src/client/index.ts:158`）与 `ReferenceCodec` 契约（`ui-input-trigger/src/types.ts:142-147`），**没读它把引用序列化成什么喂给模型**。
3. **DSH `ui-directory-picker-browse` 的实现**：只读了 `-native` 的 README 与口径，浏览版内部列表/权限判定未读。
4. **minimax `packages/agent-modules/permission/`**：问题 6 在 minimax 一栏的"未读到"就是这个目录没展开（里面有 `tools/bash-safe-first-words.ts` 等，明显是**工具侧**风险词表，不是命令侧）。
5. **minimax `packages/tui/src/tui/features/inspection/` 与 `plugin/manager.ts`**：`SkillRegistryDump` 到底画不画给人看、插件管理器显示哪几列，未取证。
6. **openchamber 的 skills/commands store 的上游**：`useSkillsStore` / `useCommandsStore` 的数据是从 OpenCode 服务端还是本地文件来的，**没往上游追**（只读到 `selectCommandsForDirectory(s, currentDirectory)` 这一层）。
7. **openchamber 服务端上传上限的具体数字**：只读到客户端如何把 413 翻成 toast，**没读到服务端在哪里设的 body limit**（`packages/web/src/api` 内 `grep -iE "MAX_UPLOAD|bodyLimit|413"` ⇒ 无输出）。
8. **Step-Code 的加号/GUI**：Step-Code 只有 `apps/cli`＋`packages/tui`，**没有 GUI 加号菜单**；owner 截图 3（`Seed-Code` 模型＋飞书文档/通讯录二级项）**在五家克隆里都找不到对应代码**（推导式：`grep -rln -iE "飞书\|lark\|feishu" Step-Code/packages minimax-code/packages openchamber/packages pi-upstream/packages` ⇒ 命中的 10 个文件全是 `lark` 作为英文词根/无关命中，无一处是"飞书文档"菜单项）。⇒ **截图 2、3 是别家（疑似 Qoder/TRAE 一类商业产品）的界面，只能当"形状参考"，不能当"某家实现"引用。**
9. **"严格档下命令被拒"这一形**：**五家里只有 minimax 的 side mode 白名单（§4.4）与 DSH 的 `available()` 沾边；没有任何一家把"权限档位"直接接在命令分发上。** 我们票 213 AC#3 想要的东西，**五家都没有现成的**——这是设计缺口，不是我没读到。
10. **五家的 `.git` 一律不存在**（已确认），全文未引提交历史。

---

## 8. 我认为编排者会误读的地方（他之前的哪条可能被推翻）

1. **票 214 现量表里"附件受理器生产零调用者"这句仍然成立，但"输入框载体只有档位＋模型两枚选择器"这句要收紧。**
   实情：`ComposerState` **已经有 4 枚附件字段**（`attachments` / `acceptedAttachmentMimes` / `maxAttachmentBytes` / `attachmentError`，`internal/panel/composer.go:238-241`），而且**生产路径上真的在填**——只是填的是 `AcceptedMIMETypes()` 与上限数字，附件列表那一枚被 `pump.go:226` 传成字面量 `nil`。
   ⇒ **会被误读成"要新加字段"**。不是：**字段与尺都在，缺的是生产者。** 票 214 AC#1 的"改前 0 枚 ⇒ 改后 ≥1 枚"仍然对，但落点应写成"给 `PumpSources` 加一枚 Attachments 读者＋把 broker 接上"，不是"给 ComposerState 加字段"（加了反而会打红那两枚"加字段必两侧同批移动"的尺）。
2. **票 213 现量表"别家怎么分层：DeepSeek 有 `ui-commands`＋`ui-input-trigger` 两枚分包"——这句不完整，容易被读成"分层是后端给的"。**
   实情：DSH 那两组（Add/Commands）的**成员名单是前端硬编码的 8 个字符串**（`presentation.ts:19-22`），后端目录只给 `name/description/input`。**后端从来不说"这条属于哪一组"。**
   ⇒ 若按"分类由后端给"去设计 C17 快照的新键，会白造一层契约面。
3. **owner 那句"比如/可以调出skills，MCPs"——五家里没有一家把 MCP 放进斜杠菜单。**
   实情：MCP 在 DSH 走设置页；在 Step-Code 是 `/mcp` **只读回显**（`step/slash-commands.ts:45-55`）；在 Pi/minimax/openchamber 本轮**未读到任何 `/mcp` 命令**。
   ⇒ 若把这条当需求做，会撞上票 215 的"不接 MCP（D13 取舍）"禁区。**建议把"菜单里有 MCP 的位置"与"接 MCP"分开，只做前者。**
4. **票 213 判据 AC#3"拒了要说为什么"——不能引"业界都这样"来支持。**
   实情：五家分三派。DSH＝报错不降级（§1.2）；minimax＝拒＋人话 reason＋连退出路径一起给（§4.2）；**openchamber＝根本不拒，未知名册的 `/token` 保持普通文本**（§5.2）；Pi＝命令级不拦、参数级报错（§2.2）。
   ⇒ 这条判据是**我们自己的选择**，不是抄来的共识。写票面时别写"参照别家"。
5. **"扩展命令撞内置名"没有标准做法，两形直接矛盾**：Pi 出**警告诊断**（`interactive-mode.ts:668-679`），openchamber **静默忽略**（`ChatInput.tsx:793-795`）。
   ⇒ 若编排者按"别家都拒绝重名"下判据，会漏掉"拒绝之后说不说"这一维。**这一维必须我们拍。**
6. **票 215 AC#2"配置 true 但没生效要能被看出来"——最好的样板不在任何一家的菜单里，在 minimax 的 `SkillSnapshot` 里**（`winners/losers/diagnostics/refresh_metrics/render_metrics`，§4.3）。
   ⇒ 若只去读别家的"菜单代码"，会以为别家也没做、于是把它当高成本项。**它的成本是一枚快照 struct 的字段数，不是菜单交互。**
7. **"命令≠工具"这条禁区（票 213）在 openchamber 有一处看起来像反例，要提防被当成依据**：openchamber 的多数斜杠命令实现为"给用户看一句话＋给模型塞一段指令"（`composer/submit/slashCommands.ts:1-13, 22-24`）。
   ⇒ **这不是"命令做成模型可自调的工具"，但它是"面板侧来源给模型下指令"**，与我们那条"由面板侧来源的 L2 允许"禁令同源。**引用时必须带这一句警告，否则会被读成"别家也是这么给模型发指令的，那我们也可以"。**
8. **`AGENTS.md` §1.2 那条"零 emoji"与本调研无关，但 owner 截图里有 `⚠`（U+26A0）与 `#`。** 按 AGENTS.md 记的仪器射程（`U+2600–U+27BF` 在扫、注释豁免、字符串不豁免），**截图里那个 `⚠` 若照抄进界面字符串会被 `d22scan` 抓到**。⇒ 别把"别家这么画"当成"我们可以这么写"。
9. **票 215 的现量表第一格要推翻重写。** 票面写的是"`internal/config/schema.go` 里 `skill`／`plugin` 字样命中（枚数现读）"，实情是：**`grep -c Skill internal/config/schema.go` ⇒ 0**、`grep -c Plugin` ⇒ **8**。
   ⇒ **技能这一维今天连配置键都没有**，票 215 第 1 格"清单可枚举、每一条都有真装配读者"在技能这一支上不是"没人读"，而是**连写都没地方写**。这会改变本票的工作量判断：**技能那一半要先立配置键（碰 `docs/specs/SPEC-03` 的 section 树，D36 范围），再谈清单；插件那一半才是"键在、读者缺"。** 建议把票 215 拆成 215a（插件：已有键，补读者＋补清单）与 215b（技能：先补配置面），否则一张票里混两种前置。
10. **票 213 现量表第 2 格"输入框现有载体只有档位与模型两枚选择器"这句会被读成"要新加字段"。** 实情见 §0 修正行：`ComposerState` 有 9 枚字段，其中 **4 枚是附件相关**，生产路径 `pump.go:226` 把附件那一枚实参写死成 `nil`。⇒ 落点是"给 `PumpSources` 加读者"，不是"给 `ComposerState` 加字段"。

---

## 9. 变成我们的实现约束（按票 213 / 214 / 215 分开，大白话，不用术语）

> 每条都写成"我们必须怎样，否则会出现什么后果"。括号里是给编排者看的出处锚点，owner 可以整段跳过。

### 9.1 票 213（命令名册与执行通道）

1. **名册里每一条命令，必须同时带着"它现在能不能用"和"不能用的时候说哪句话"。** 不能用的时候不许只变灰。（minimax 的 `unavailableReason`，§4.1、§4.2）
   否则：用户在菜单里点一条没反应，或者点灰掉的条目什么话都没出——就是 `A408` 点名的那一形。
2. **一条命令收不收参数，要写在名册这一行上，不能靠"试了才知道"。** 写了收参数的才允许后面跟文字；没写的，只要多打了一个词就不算这条命令。（minimax `argumentHint` 的口径，§4.1、§4.2）
   否则：会出现"看着像命令、按下去发出去变成普通提问"。
3. **判断"这一行是不是命令"的那段代码，只能有一份，界面上用来高亮、回车时用来执行，都调它。**（minimax §4.2 的设计声明；openchamber 的反面历史 bug：`@` 规则被写了四遍、`/` 被写了三遍，导致"高亮成了引用但解析不成引用"，§5.2）
   否则：同一个 token，输入框里显示成命令、按下去当普通文本发出去——这是我们最容易自己造出来的 bug，因为界面那支和 Go 那支是两拨人。
4. **没打中任何命令的 `/xxx`，必须给一句话，而不是当普通消息发出去。**（DSH 明确"绝不静默降级"，§1.2）
   否则：用户以为触发了什么，实际发出去一段带 `/xxx` 的提问，模型会照着演。
5. **报错的那句话要带上"有哪些是对的"。**（DSH `/permission` 的 `unknown preset "x" (available: a, b, c)`，§1.2；Pi 的 `Unknown thinking level "x". Available levels: …`，§2.2）
   否则：用户只能靠猜重试。
6. **一条命令能不能出现在名册里，取决于"它的真执行者在不在生产路径上"，不取决于配置文件里有没有那一行。** 我们仓里已经有比别家更硬的做法：`internal/config/unwired.go` 那 6 枚"写了不管用就响亮拒收"的键（§6 现量表）。**名册要用同一思路：装配时没有真执行者的命令，注册这一步就该失败，而不是注册了再置灰。**
   否则：菜单里有一行、点了没反应——正是票 213 AC#2 要防的。
7. **命令名字撞了内置命令时，"拒绝"和"要不要出声"必须一次定死。** 别家两种做法直接矛盾（Pi 出声警告、openchamber 静默忽略，§8 第 5 条）。**建议照 Pi：不进名册，但留一条警告，和 `unwired.go` 一个脾气。**
   否则：以后有人加了个和 `/compact` 同名的东西，谁都不知道谁生效了。
8. **命令的执行一律走现成的审计那三档，不留第二条路；并且"谁下的这条命令"这个字段只能是"人"。**（DSH 把命令执行者类型写成只有 `user` 一种变体，理由注释在 `packages/interaction/commands/src/types.ts:68-79`，§1.6）
   否则：模型或网页里的内容有一天能自己下命令，等于把加号菜单变成提权入口。
9. **界面点不动之前，票面终态必须继续写着"界面点不了"。**（票 213 AC#4、票 215 前置段；账 `A410`）
   否则：owner 会以为能用了，然后返工——他今天已经返过一次。

### 9.2 票 214（往这次对话里加东西：文件、附件、引用）

1. **先接线，再加东西。** 附件那台机器（410 行）已经在，`ComposerState` 的四个附件字段也已经在，**缺的只是"有人把文件递进去、有人把结果读出来"**：`internal/panel/pump.go:226` 现在把附件那一枚实参写死成 `nil`，`PumpSources` 里没有读附件的那一枚。
   否则：会误做成"新加一遍字段"，反而打红"加字段必两侧同批移动"那两把尺。
2. **上限的数字从配置读，判断的动作在真正落盘那一层做，界面那一层只做提前拦。** 别家都是两层都判、**后端才是权威**（DSH §1.5 注释原话 `Attachment admission is enforced here, not in the composer`）。
   否则：界面说能传、后端拒收；或者界面拦了、命令行绕过。
3. **被拒时给的话要三样齐全：哪一个文件、为什么、怎么办；而且"为什么"要用错误码，不用错误文字去匹配。**（反面教材：openchamber 靠英文串里有没有 `413` 来决定弹哪句提示，服务端换个措辞提示就永远不弹，§5.4；正面：同仓另一处用 `code === 'too-large'` 映射文案）
   否则：以后我们自己的 Go 改了错误措辞，用户就再也看不到拒收原因，而且没人会发现。
4. **类型判定按内容，不按扩展名；被拒的那一条不许把调用方自称的类型回显成"已经判过了"。** 我们仓里这条**已经写对了**：`internal/panel/attachments.go:87-90` 的注释——`MIME is the SNIFFED type, never the declared one. It is empty unless Stored is true: a refusal must not echo the caller's own claim back as if a verdict had been reached about it.`
   否则：伪装成图片的文件被当图片收下并喂给模型。
5. **附件清单要能一眼看清"第几枚、叫什么、什么类型、多大"，并且文件名要先洗过再显示。** 别家两家各写了一份洗控制字符的代码（minimax `sanitizeTerminalText`，§4.4；Step-Code `neutralizeFeedbackConsentText`，§3.4），**两家都专门对"文件名"这一枚下手**。
   否则：一枚名字里带控制字符的文件能把面板那一行字改成别的话。
6. **"预览不出来"和"没附上"是两句话，不许混。**（minimax 的中文文案就写了这个：`无法预览 · 附件仍可发送`，§4.4）
   否则：用户以为附件丢了，重发一遍，同一份内容进上下文两次。
7. **@引用工作区文件，和拖进来一个文件，必须走同一台受理器。**（openchamber 就是这么做的：粘贴、拖拽共用 `attachFilesWithCitation`，§5.4）
   否则：出现第二条读文件的路，票 214 的禁区直接破。
8. **选文件/选目录的"选法"和"选完之后判什么"要分开。** 别家把选目录做成两枚可以互换的包（原生对话框版／自绘浏览版），**切换只换装配、不改代码**（DSH §1.5）。我们这边：**Go 侧只保证拿到路径之后判定链完整（授权根、改写账户、来源戳），选法由界面那支定。**
   否则：Go 侧写死一种选法，将来接远端或换宿主就得重开一票。
9. **落盘位置必须在它该落的那棵树里，越界那一形必须红。**（票 214 AC#4；票 174/175 那一族）
   否则：附件写到授权根外面，等于用加号开了一个任意写口子。

### 9.3 票 215（技能与插件：列出来、能开关、能看出这次生效没生效）

1. **"生效没生效"不要做成菜单里的一个勾，要做成一张能数出数的表。** 别家做得最好的不是菜单，是它的数据结构（minimax §4.3）：一张快照里同时有 **赢的、输的（为什么输、被谁顶了）、诊断（警告还是错）、扫描了几枚/读了几枚/复用了几枚、渲染了几枚/压掉几枚/丢了几枚、预算用了多少**。
   否则：只会得到一个"看起来是开着的"，还是答不上 owner 那句"我配了呀，怎么没生效"。
2. **"输的那一枚"必须留在表里，不能只留赢的。**（minimax 的 `losers` 带着 `reason` 和"赢的那枚在哪"，§4.3）
   否则：同名技能被顶掉的时候，用户改了半天配置文件，改的其实是被顶掉的那一枚。
3. **被上下文预算挤掉，也是一种"没生效"，必须单独报数。**（minimax 的 `render_metrics`：`total / rendered / compacted / dropped / charsUsed / budgetChars`，§4.3）
   否则：会出现"配置对、装配对、模型这次就是没照着做"，而所有人都去查配置。
4. **每一枚技能/插件要带出处：谁装的、装在哪一层（用户级还是这个项目带的）、优先级多少、是不是外部内容。**（Pi 的 `SourceInfo { source, scope, origin, path }`，§2.3；minimax 的六种来源枚举＋`sourceExternal`，§4.3）
   否则：一个仓库里塞进来的技能和自己装的技能长得一样，等于给外部内容开了一个"看起来是用户自己装的"的面子。
5. **出处要直接显示在那一行字上，不要藏进详情。**（Pi 的做法：描述前面拼 `[来源]`，`interactive-mode.ts:660-666`，§2.3）
   否则：没人会去点详情，出处等于没有。
6. **技能在菜单里被选中，只是往输入框里放一段字面文本；真正加载由后端在下一步决定。**（DSH §1.4：`onPick` 只塞 `/name `，`Determinism lives host-side`）
   否则：界面以为点了就生效，实际后端根本没认。
7. **菜单里那一行的"能不能用"，要按"这个会话现在能不能用"来判，不是按全局配置判一次。**（DSH 的 `available(session)` 每次枚举都重新问一遍，§1.3；`/model` 在子 agent 会话里就是不可用）
   否则：会出现"上一句还能用、这一句点了报错"。
8. **开关动作走现成的配置写路径，并且开关本身要进审计；做不到立刻生效的，就在清单上写明"这一枚要重启"，不许静默假装生效。**（票 215 AC#3 原文；别家没有一家替我们解决了这一条）
   否则：用户关了它还在跑，或者开了它没反应，两种都查不出来。
9. **菜单里的东西一律不许改权限档位、不许改 `allowed_dirs`。**（票 215 AC#4；DSH 的对应做法是把"谁能执行命令"写成只有"人"一种，§1.6，比加规则强）
   否则：一份外部插件说明书里写一句"以后都允许"就有机会提权——这是票 200 那一族。
10. **同意页/清单页要先把"会发出去什么"摆出来，再问用户；摆出来的每个字段都要先洗过控制字符。**（Step-Code 的 `/feedback`：先列清单再要同意、先脱敏再限长、每个外部字段各洗一遍，§3.4）
    否则：一枚精心构造的文件名能把那页确认文字改掉，用户点了"确认"点的不是他读到的那句话。
11. **MCP 只做"菜单里给它留个位置"，不做"接 MCP"。**（票 215 禁区＋D13；owner 那句"调出 MCPs"按这一条落，见 §8 第 3 条）
    否则：顺手把 D13 的取舍推翻了，那是契约级变更，要人工批准。

### 9.4 三票共用的两条

1. **名册、附件状态、技能清单这三样要进快照，都走同一枚规则：读者不存在就不发这个键，读者存在就永远发（哪怕是空）。** 我们仓里已经把这个规矩写在 `PumpSources` 的注释里了（`internal/panel/pump.go:133-141`：`A reader that exists always sends the key … because "this run has spawned nothing" and "this host cannot see subagents" are two different things a page must not render the same.`）——**三票都照这句办，不新造第三种。**
   否则：界面把"没有"和"看不见"画成同一个样子，owner 看到的永远是前者。
2. **新键会打红"加字段必两侧同批移动"那两把尺：红因写进票面，由 owner 带给界面那支，不改尺、不加豁免。**（票 213 第 3 条原文）
   否则：又是一次越界改判据，票 130 那一形会重演。

### 9.5 如果三张票装不下，我建议开的第 4 枚

**「加号菜单的显示层：任何一行字，只要内容不是用户自己打的，都要先洗过控制字符再显示；并且"没电/读不到"和"读到但是空的"必须是两句话。」**

为什么它不属于 213/214/215 任何一枚：
- 213 管命令名册、214 管一次对话的输入、215 管长期挂着的东西——**这一枚管的是三票共同的显示底座**，射程横跨三票的每一个界面字符串。
- 证据强度：五家里有**两家各写了一份**、写法还不一样（Step-Code 把控制字符转成可见转义 `consent.ts:73-92`，minimax 直接删/替换 `terminal-text.ts:15-26`），而且**两家都专门对"文件名"这一枚下手**（§3.4、§4.4）。openchamber 那一族的反面证据也在：靠英文错误串决定弹哪句提示（§5.4）。
- 我们仓里已经有对的一半：`internal/panel/attachments.go:87-90` 那句"MIME 是嗅出来的、被拒不回显调用方自称的类型"，以及 `PumpSources` 注释里那句"没有和看不见不能画成同一个样子"。**缺的是把这两条收成一枚常驻判据，覆盖三票新增的每一行字。**
- 成本判断：**低**（一枚纯函数＋一批判据），收益判断：**高**（不做的话，票 214 的附件文件名、票 215 的插件名与描述、票 213 的命令行回显，三处都会各自发明一套洗法，然后又得合并）。
- 若 owner 只想开三枚：把这一枚降级成 214 的 AC#5 的加强版（"来源戳盖上"那一格扩成"显示前洗过"），**但要在 213/215 票面各写一行"引用 214 AC#5 的洗法"**，否则只有附件受益。
