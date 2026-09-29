# Wisp 缺失功能清单 · 第四版（2026-09-29 v4）

> 调研腿：`v4-survey-1`（只读）。上一版＝`docs/reports/missing-features-2026-09-28-v3.md`（125,852 字节／67 条）。
> 这一版做三件事：**并回**第三版之后落地的五份没并回的调研件、**挑自己上一版的漏**、**打还没覆盖的地方**。
> ⛔ 本文件不新造规矩、不改契约；一切以 `docs/PLAN.md` 与 `docs/specs/SPEC-*.md` 为准。

---

## §0 口径与分母（先说数字，再说读了多少）

- 参照仓（本机现读，一律只读，被研仓里零写入）：`D:\work\AI\open source\` 下 6 家
  ＝ `Step-Code` / `deepseek-harness` / `minimax-code` / `openchamber` / `pi` / `pi-upstream`。
- **六家今天都有没有提交历史（这条改了上一版的前提）**：`Step-Code`／`deepseek-harness`／`minimax-code`／`pi`／`pi-upstream`
  **五家都有 `.git`**（HEAD 与各自主张的日期见 §5-2），**只有 `openchamber` 是 tarball 快照、没有 `.git`**
  ⇒ 只有 openchamber 需要"每一条都带 `文件:行`、不引历史"；其余五家两者都可以引（本版仍一律带 `文件:行`，因为那只尺更严）。
- **DeepSeek 有三枚对象，引用必须先说哪一枚**（详见 §5-3）：
  ① 源码克隆 `D:\work\AI\open source\deepseek-harness` ＝ **`0.1.7-rc.2`**（有 `.git`，`package.json` 的 `version` 行）
  ＝ 本版与第三版所有 DeepSeek 行号的出处；
  ② npm 全局 `@deepseek-ai/dsh` ＝ **`0.1.0-rc.6`**；
  ③ `~/.dsh/profiles`（下面只有 `desktop`／`web`／`node_modules` 三项，根目录**没有** `package.json`），
  其中 `node_modules/@deepseek-ai` ＝ **257 枚**（`node_modules` 顶层 276 枚）。派单说的 `0.1.6-alpha.1` 这个号本轮没验出来。
- 本轮**没有新增任何克隆**，也没在被研仓里建过任何文件。
- 我方一侧的「有没有」＝只按 Go 侧现读（`internal/**`、`cmd/**`）＋票面/台账判定；`frontend/**` 与 `design/**` **零读零引零转述**。
- 每条挂状态标记：〔已证＝我们真做了〕／〔建了但没接＝码在、生产没人调〕／〔仅文档＝规格写了、盘上没有〕／〔别人也没有〕。
  标〔已证〕的每条都现跑了"找生产调用者"的 grep，找不到就降级（本版没有一条〔已证〕是靠上一版的说法继承下来的）。
- 分母与实读量：见 §7。

---

## §1 表一 · 逐块逐控件（别人有、我们没有的控件）

> 一屏一块，每块里**逐个控件**列。"别人有"一律点名哪一家＋`文件:行`；空手写不出的写成〔未量〕。
> ⚠ 本表的"我们这边"三列只按 Go 侧现读＋工单/台账判定，界面那两块地界一律标〔界面地界·不读〕。

### 块 1：聊天输入框那一行（＋它旁边的加号菜单）

| 控件（用户看到的名字） | 谁家有（点名＋`文件:行`） | 我们这边状态 |
|---|---|---|
| 加号按钮点开的那张菜单（固定两段：「添加」／「命令」） | DeepSeek 克隆 0.1.7-rc.2：两段的行名硬编码在前端 8 个字符串里 `deepseek-harness/packages/client/ui-commands/src/client/presentation.ts:19-22`（`add: file/goal/plan/feedback`、`commands: compact/permission/model/export`），段标题走本地化键 `:82-84` | 〔根本没做〕本仓 Go 侧搜"名册/命令目录"三族词 **0 命中**（`grep -rln -iE "slashcommand|commandcatalog|CommandRegistry" internal/ cmd/` 空输出） |
| 「＋→选本地文件」 | DeepSeek：`file` 是前端动作项、不是后端命令 `ui-conversation/src/client/apply.ts:270-278`；真选文件用系统自带的选框 `ui-conversation/src/client/skeleton/InputBar.tsx:433-439`，**子会话里直接不给选**（同一行 `disabled={subagent !== null}`） | 〔建了但没接〕受理器在（`internal/panel/attachments.go` 410 行），**产品装配点把附件那一枚实参写死成空**：`internal/panel/pump.go:226`；入向那道门是明写的空＋会报名字（`cmd/wisp/panel_inbound.go:236`，`internal/panel/composer_dispatch.go:62`） |
| 「＋→链接 GitHub 问题／Pull Request」「＋→链接 Linear 问题」 | openchamber（无 `.git`）：四段固定顺序、后两段按能力位条件渲染 `openchamber/packages/ui/src/components/chat/composer/ui/ComposerAttachmentControls.tsx:91-143`（Linear 那一段在 `:133`） | 〔根本没做〕我们这边没有任何"引用外部工单/链接"的入口 |
| 「＋→扩展贡献进来的项」 | openchamber：`attachGuests?.map(...)` 同文件 `:143`；三类贡献点登记在 `composer/DOCUMENTATION.md:125` | 〔根本没做〕 |
| 大段粘贴自动变成附件（三选一：问我／直接附／原地贴） | openchamber：≥ 约 2000 字符或 25 行触发，按设置项 `largeTextPasteBehavior` 走，attach 时造一枚内存里的 `pasted-context-N.txt` 并**与手选 .txt 走同一条管线** `composer/DOCUMENTATION.md:137-148` | 〔根本没做〕 |
| 输入框里的斜杠命令面板（打 `/` 就出候选） | 五家都有：DeepSeek 前端纯函数检测 `ui-input-trigger/src/core/detect.ts:48-76`（含两条"这不算触发"的网址豁免 `:20-30`）；Pi 系要求行首 `pi-upstream/packages/tui/src/autocomplete.ts:432`；minimax 精确等名 `minimax-code/packages/tui/src/tui/commands/catalog.ts:698-717`；openchamber "名册才是权威、正则只是定位器" `packages/ui/src/lib/components/chat/composer/language/prefixTokens.ts:11-14` | 〔根本没做〕 |
| 输入框里的 `@` 引用工作区文件／别的会话 | DeepSeek `@` 用另一套文法（支持带引号跨空白）`detect.ts:6,50-60`；Pi 系文件补全走外部 `fd` 或自扫并打分 `pi-upstream/packages/tui/src/autocomplete.ts:736-755`；openchamber 粘贴与拖拽共用同一枚受理函数 `composer/DOCUMENTATION.md:131-136` | 〔建了但没接〕同"选本地文件"那一行（同一台受理器、同一个空实参） |
| 加号旁边一枚**独立的权限档位按钮**（不藏在菜单里） | openchamber `packages/ui/src/components/chat/composer/ui/PermissionAutoAcceptButton.tsx`（同目录还有 `FocusModeButton.tsx`、`DraftTargetSelectors.tsx`） | 〔界面地界·不读〕Go 侧有档位读写（`internal/perm`、`internal/panel` 的档位处理器），**但没有"每会话自动接受"这一层** |
| 草稿芯片（引用别处来的东西、带一枚不透明数据回传） | openchamber `ComposerContextChips.tsx`＋`LinkedReferenceRow.tsx`；那句"扩展给的不透明数据原样回来、绝不进上下文文本"写在 `composer/DOCUMENTATION.md:125` | 〔根本没做〕 |
| 附件清单那一行「第几枚 · 文件名 · 类型 · 大小」 | minimax 固定三段格式 `minimax-code/packages/tui/src/tui/features/composer/attachments.ts:105-108`；文件名先洗终端控制串再显示 `rendering/terminal-text.ts:15-26` | 〔建了但没接〕我们的受理器有"嗅出来的类型、被拒不回显自称类型"这一半且**写对了**（`internal/panel/attachments.go:87-90`），**清单那一行没人发** |
| 「无法预览 · 附件仍可发送」这一句区分 | minimax 中文文案逐字 `features/composer/copy.zh-Hans.ts:5` | 〔根本没做〕 |
| 发送被打回时那句"太大"的提示（告诉你减什么） | openchamber 德文串逐字 `Anhänge sind zu groß zum Senden. Bitte versuche, die Anzahl oder Größe der Bilder zu reduzieren.`（`packages/ui/src/lib/i18n/messages/de.ts:2240`）⚠ 但它判的是**英文错误串里有没有 413**（`ChatInput.tsx:2205-2206`）＝反面教材；同仓另一处用码不用串（`components/sections/extensions/ExtensionsPage.tsx:58`） | 〔根本没做〕 |

### 块 2：会话与历史（列表／标题／那一页的出口）

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| 会话列表按"需要你关注"排序、带未读点 | openchamber 判据逐字 `needsAttention = unseen>0 && (!isSubtask || notifyOnSubtasks)`，与会话侧栏共用同一判据 `mobile/mobileWidgetSnapshot.ts:18-20,84-98` | 〔根本没做〕 |
| 会话标题可以改名 | Pi 系 `/name` 在名册里 `pi-upstream/packages/coding-agent/src/core/slash-commands.ts:18-43`（24 枚内置，尺 `grep -c "{ name:"`）；minimax 会话操作有整页文案 `packages/tui/src/tui/features/session-mutation/copy.zh-Hans.ts:39-173` | 〔界面地界·不读〕Go 侧本轮未量到改名出口 |
| **回退（把对话退到某一条）＋"要不要连文件一起还原"** | minimax 两枚分开：`回退会话`／`回退会话和文件` 各有说明 `copy.zh-Hans.ts:28-31`；预览行 `{ready} 可回退 · {skipped} 已跳过` `:71,158`；失败还有"会话已回退，但文件还原未完成" `:161` | 〔根本没做〕 |
| 从任意一条历史消息**开分支** | minimax `从这里创建分支`＋"创建子会话，不修改当前会话或文件" `copy.zh-Hans.ts:24-25`；确认页把"当前会话／分支点／工作区是否会隔离"三行摆出来 `:136-144`；Pi 系有 `/fork`、`/clone`（`slash-commands.ts:18-43` 名册内） | 〔根本没做〕 |
| 编辑某一条历史提问并从这里重新生成 | minimax `编辑此输入`＋"替换此提示词，并从这里重新生成" `copy.zh-Hans.ts:26-27`；编辑到一半失败时"草稿已保留，按 Enter 作为新消息发送" `:120-123` | 〔根本没做〕 |
| **会话出口：导出整段对话** | DeepSeek 命令只回"已请求"、真文件走一条带鉴权的下载路由、**明确拒绝带路径参数** `deepseek-harness/packages/session-query/session-log-export/src/index.ts:78-90`；Pi 系默认 HTML、后缀决定格式、另有私密 gist 分享 `pi-upstream/.../slash-commands.ts:25-27` | 〔根本没做〕第三版整版没有这一条（见 §5-9） |
| 复制最后一条回复到剪贴板 | Pi 系 `/copy` `slash-commands.ts:28` | 〔根本没做〕 |
| 「提交反馈／报 bug」那一条完整链 | Step-Code：15 个文件 3,101 行；**先把"会发出去哪些文件"的清单摆出来再要同意** `Step-Code/packages/coding-agent/src/step/feedback/bundle.ts:110,113`；分层回落的上限（超了先缩、再超才放弃并给可枚举的理由）`:19-26,88-97`；**先脱敏再限长且写明怎么截的** `:176-185`；同意页每个外部字段各洗一遍控制字符 `consent.ts:73-92` | 〔根本没做〕 |
| 右侧那一页的第二种看法（逐事件台账＋时间轴） | DeepSeek `ui-trajectory`：7 种行、28 个可选位置、拖选时间区间（`survey-2026-09-28-dsh-ui-packages.md:204-210` 逐条带行） | 〔根本没做〕源头缺：`internal/session/` 今天仍只有 `doc.go`，往外冒的消息仍只有 9 种、无序列无时间（§5-11） |
| 上一次怎么看它、下次点进来还是那个看法 | DeepSeek 渲染前先读持久化的"看法偏好" `ui-conversation/README.zh.md:67` | 〔根本没做〕 |

### 块 3：子代理与后台任务那一块

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| 页头一枚"有几枚孩子在跑"的开关＋可展开的族谱树 | DeepSeek `ui-subagent`：状态点＋模式＋活动＋标题＋**累计花费**＋**已干多久** `SubagentHeaderLineage.tsx:250-270`，`:334` | 〔建了但没接〕名册与每行已进刷给界面那一包（`internal/panel/composer.go` 的 `tasks` 节），**但一行给不出花费与计时**（`internal/tools/task.go:95-140` 六个字段里没有这两枚位置） |
| 悬停 150 毫秒打开／离开 120 毫秒关闭／点一下钉住 | DeepSeek 原文写着这三个数 `ui-subagent/README.zh.md:34` | 〔界面地界·不读〕 |
| 两种落点：点行＝主区搬进去；点行尾箭头＝右侧另开一格 | DeepSeek `ui-subagent/src/client/index.ts:61-69`；同一枚子代理**不会被开成两页**（靠"种类＋内容编号"去重）`ui-dockkit/README.zh.md:59` | 〔根本没做〕我们的寻址仍是"一条流"（`subagent:<编号>`，`internal/panel/pump.go:386`） |
| 后台任务列表：进行中／已结束两组＋可展开的流式输出 | DeepSeek `ui-jobs/README.zh.md:11,28-38`：**没用到一枚都不渲染**（`:26`）、收起即停流（`:34,38`）、复制的是**命令**不是输出（`:34`） | 〔界面地界·不读〕Go 侧：`task.list` 今天仍逐字写着 DEFERRED（`internal/tools/task.go:22`） |
| **两击式停止**（第一击上膛、3 秒内确认击才真停） | DeepSeek `ui-jobs/README.zh.md:30` | 〔根本没做〕Go 侧搜这三族 0 个生产落点 |
| 停止之后**把"是人停的"告诉模型** | DeepSeek 同句注释逐字："模型被明确告知用户停止了它的任务，而不是留给它去猜" `ui-jobs/README.zh.md:30` | 〔根本没做〕我们只有"告诉人"那一半（`internal/agent/sink.go:10-15`） |
| 会话头部"这一步能同时跑几枚工具"的上限读数 | DeepSeek 把它做成设置页 `ui-settings-agent-loop/README.zh.md:28` | 〔仅文档／冻结〕我们的并发上限写死在桥里（`internal/tools/bridge.go` 的 `MaxToolConcurrency`），`internal/config/schema.go` 搜不到对应键 |
| 派子代理时**说明文字里许诺的那枚"单独停它"的工具** | 别家的判据是"能力没落地就不进名册"（DeepSeek 用组合期闸门 `packages/plan/plan-mode/src/index.ts:230-231`；minimax 用 `audience:'internal'`＋`visibleWhen` 双重 `catalog.ts:47-84`） | ⚠ **我们自己反着来**：派子代理的说明文字许诺"可以用 task.cancel 单独停它"（`internal/tools/subagent_197.go:196`），而 `task.cancel` 在名册头上逐字 DEFERRED（`internal/tools/task.go:23`）——见 §5-13，票 221 在飞 |

### 块 4：要批准／要回答那一块

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| 批准卡上「仅本次允许」这一枚按钮 | DeepSeek 面板只有两枚决定：`allowed-once`／`rejected`（`ui-approval/src/client/contract/slots.ts:66`），持久策略划给服务侧（`ui-approval/README.zh.md:36`）；openchamber 画**三枚**（`Allow once`／`Always…`／`Deny`，`packages/ui/src/components/chat/PermissionCard.tsx:407-442`，快捷键 `alt+enter`／`alt+shift+enter`／`alt+backspace`） | 〔根本没做〕面板那一侧只能拒绝，"允许"被结构禁掉（`docs/specs/SPEC-08*.md:156-176` C17 名册那条逐字）。⚠ 这一格碰 `Q-49`，属"要 owner 一句话"，我不自行填 |
| 「一直允许」把**要存成哪条规则**印在按钮上 | openchamber：文案逐字 `Always: {patterns}`，取服务端回的数组前 2 条＋省略号、全量挂到悬停提示 `PermissionCard.tsx:349-358` | 〔仅文档〕 |
| 一次列出几档"宽度"供你挑（最窄／第一个词／前两个参数／按域名／整个工具），**第一档永远是最窄的** | minimax `packages/agent-modules/permission/src/types.ts:151-195`，兼容律逐字写在 `:96-106` | 〔仅文档〕我们只有 `docs/PLAN.md:2028` 手写的三档且只针对一枚工具 |
| 批准卡的键盘规矩：详情聚焦后 Enter 批准／Esc 拒绝；插件挂载期间这两键不许被别的快捷键抢走 | DeepSeek `ui-approval/README.zh.md:21`；且"审批"与"按两下 Esc 停止"在快捷键目录里是**不可改的贡献项** `ui-shortcuts/README.zh.md:77` | 〔仅文档〕我们有"Esc 临时接管必须归还"的规矩但没有把它写成快捷键策略（`docs/specs/SPEC-08*.md:226-227`） |
| 关联详情：这次批准背后那一次工具调用可点开 | DeepSeek 可选注入一枚槽，按调用号渲染 `contract/slots.ts:37-49` | 〔已证〕我们卡上带调用链与"命中了哪几条规则"（`internal/panel/approval.go:39-59` 共 11 枚字段，比他们那 5 枚多）——**这一格是反向标** |
| **agent 反问人一句**（带选项、可多选、可自己写、可跳过、可关掉，且"跳过"与"关掉"是两件事） | DeepSeek 整包 `ui-user-questions`＋服务端类型 `packages/interaction/user-questions/src/types.ts:36-67`；跳过＝交一份空答案、关掉＝以 `ASK_CANCELLED` 拒绝整个等待（`ui-user-questions/README.zh.md:32`）；题号/已选/自写内容只放不持久存储（`:40`）；同一时刻只有一枚问题拥有输入框（`:95`） | 〔根本没做〕`ask_user` 一族在 Go 侧今天仍 **0 命中**（本轮复量） |
| 计划评审那枚批准：批准项**按名字点名**、不按位置 | DeepSeek `user-questions/src/types.ts:27-29` 逐字"按名字而非按位置，免得界面从顺序里猜结论" | 〔根本没做〕 |

### 块 5：命令目录那一块（打 `/` 之后出现的那张表）

| 控件／一栏 | 谁家有 | 我们这边状态 |
|---|---|---|
| 名册每一行都带"现在能不能用"＋"不能用时说哪句话" | minimax `unavailableReason` 出现 **19 次**、`visibleWhen` **18 次**（尺：`grep -c` 于 `packages/tui/src/tui/commands/catalog.ts`，同文件条目名 **41 次**）；三态 `unrecognized/unavailable/handled`＋兜底句 `:627-663` | 〔根本没做〕 |
| **这一行背后是真干活还是一句提示词**（必须分开画） | openchamber 9 枚＝两串提示词、零实现（`packages/ui/src/components/chat/composer/submit/slashCommands.ts:22-40`，转引 `survey-2026-09-29-command-catalog-across-harnesses.md:23`〔腿报〕）；minimax `/review`＝一句常量（`packages/tui/src/tui/controller/product/command-flow.ts:142` 与 `packages/tui/src/headless/invocation.ts:234` 两处独立命中）；DeepSeek `/permission`＝真写路径（`packages/interaction/permission-presets/src/index.ts:258,264-274,403-418`〔腿报〕） | 〔根本没做〕**第三版一栏都没有**（§5-8） |
| 命令收不收参数**写在名册这一行上**，没写的多打一个词就不算这条命令 | minimax `argumentHint` 口径逐字 `catalog.ts:47-84`；判定 `if (args && !source.argumentHint) return undefined` `:714` | 〔根本没做〕 |
| "判断这一行是不是命令"的代码**只能有一份**，高亮与回车共用 | minimax 设计声明 `catalog.ts:692-696`；openchamber 的反面历史 bug：`@` 规则被写过四遍、`/` 三遍，于是"涂成引用却解析不成引用"（`packages/ui/src/lib/components/chat/composer/language/prefixTokens.ts:1-15`、`DOCUMENTATION.md:150-161`） | 〔根本没做〕⚠ 我们界面那支与 Go 那支是两拨人，**这一形最容易自己造出来** |
| 没打中任何命令的 `/xxx` **必须给一句话**，不许当普通消息发出去 | DeepSeek 逐字："命令行绝不静默降级成普通提示" `ui-commands/README.md`；参数级失败带可用值回显 `permission-presets/src/index.ts:269-271`；Pi 系 `Unknown thinking level "x". Available levels: …` `pi-upstream/.../interactive-mode.ts:180` | 〔根本没做〕⚠ 这条**不能写成"业界都这样"**：openchamber 是反例（未知名册的 `/token` 保持普通文本 `prefixTokens.ts:11-14`），Pi 命令级也不拦 |
| 名字撞内置命令时"拒绝"＋"要不要出声"一次定死 | Pi 系出**警告诊断**（`pi-upstream/.../interactive-mode.ts:668-679`）；openchamber **静默忽略**（`packages/ui/src/components/chat/ChatInput.tsx:793-795`）——两形直接矛盾，没有共识 | 〔根本没做〕我们自己拍：建议照 Pi 系（不进名册但留一句警告），与我们那 6 枚"写了不管用就响亮拒收"同脾气（`internal/config/unwired.go`，147 行，尺 `grep -c "path:"`＝6） |
| 每一行带出处标签（内建／项目／插件／技能），**直接显示在那一行字上** | Pi 系描述前拼 `[来源]` `pi-upstream/.../interactive-mode.ts:660-666`；来源结构体 `{path,source,scope,origin}` `core/source-info.ts:3-11`；minimax 六种来源枚举 `agent-modules/skills/src/types.ts:1` | 〔根本没做〕 |
| 菜单里的"能不能用"按**这个会话现在**判，不是全局判一次 | DeepSeek `available(session)` 每次枚举重新问 `ui-commands/src/client/contract.ts:88-89`；实例：`/model` 在子代理会话里不可用 `ui-model-selection/src/client/index.ts:149` | 〔根本没做〕 |
| 「命令只能由人执行」写成类型上只有一种执行者 | DeepSeek `CommandSourceMap` 只有 `user` 一枚变体并写明理由 `packages/interaction/commands/src/types.ts:68-79`；命令行不进模型可见的会话日志 `:96-97` | 〔仅文档〕 |
| 严格档下某条命令被拒 | **本轮逐家点过行的五家里没有一家把"权限档位"接在命令分发上**（`survey-2026-09-28-composer-plus-menu.md:438` 已具名登记"这是设计缺口，不是没读到"） | 〔别人也没有〕 |

### 块 6：语音那一块（听与说）

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| 语音设置页：识别模型逐个列出来、**带准确度／速度两条评分**、每个能单独下载或删除 | openchamber `packages/ui/src/components/sections/openchamber/VoiceSettings.tsx:42-66,219-220,383`（四个模型 Kokoro/Parakeet 族与逐个下载删除） | 〔仅文档〕我们的识别/播报只有"提供方＋模型名"两枚文本键（`internal/config/schema.go:210-218`），**没有那一屏** |
| 点一下开始、再点一下结束的听写入口（可绑快捷键） | openchamber `ComposerDictation.tsx:127,198-210`（键名 `toggle_dictation`） | 〔界面地界·不读〕Go 侧语音三票全 DEFERRED（票 15／26／41；`internal/audio/doc.go:27` 逐字 `DEFERRED(playback)`） |
| 转写稿**绑回"开始录音那一刻的那份草稿"** | openchamber `useDictationOrigin.ts:1-10` | 〔根本没做〕 |
| 逐条消息"念这条"的手动播报 | openchamber `MessageBody.tsx` 挂播报入口、服务端本地合成（`lib/dictation/DOCUMENTATION.md:14-32`） | 〔仅文档〕 |
| 「按住说话」 | **本轮点过行的六家里我都没读到**（openchamber 那份件 `:62` 具名"没有 hold-to-talk 证据"） | 〔别人也没有〕 |
| 「它说完就自动播」 | 同上，未在任何家读到自动播报命中 | 〔别人也没有〕 |

### 块 7：球／常驻视觉件那一块

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| "有几件事在等你"的数字（应用图标角标） | openchamber：角标＝自上次打开以来**不同的提醒条数**（不是会话数），服务端算、三路信号清 `APNS.md:40-66`、`AppDelegate.swift:201-208` | 〔建了但没接〕等批准队列的长度这个数已经现成（`internal/agent/approval/queue.go`），球是天然载体，**但"数上去"那一跳没有** |
| 状态点：同一枚共享组件、三态一名、每行留同宽 | DeepSeek `SubagentHeaderLineage.tsx:334`＋`ui-subagent/README.zh.md:36`（14 像素预留列宽） | 〔界面地界·不读·给界面那支的话〕 |
| 面板里那一格终端（能打开、恢复、控制一个真命令行标签页） | DeepSeek `ui-sidebar-terminal`（**官方包地图漏掉的两枚之一**，`survey-2026-09-28-dsh-ui-packages.md:409`） | 〔根本没做〕本轮复量：刷给界面那包内容仍是固定 6 个键（`internal/panel/composer.go` 的 `Snapshot`），没有终端那一块 |
| 右侧栏的文件树／文档预览／内置浏览器 | DeepSeek `ui-sidebar-files`、`ui-sidebar-documentpreview`、`ui-sidebar-browser`（名册 `survey-2026-09-28-dsh-ui-packages.md:73-75`） | 〔界面地界·不读〕 |
| "在应用内浏览目录选工作区" 与 "用系统自带选框" **两枚可互换** | DeepSeek `ui-directory-picker-browse`／`ui-directory-picker-native`，"切换只换装配、不改代码" 逐字写在 `-native/README` | 〔仅文档〕票 214 落点已定"选法由界面那支定" |
| 每会话一枚"自动接受权限"开关，带单调修订号＋跨面板广播＋串行写队列 | openchamber（无 `.git`）`packages/vscode/src/bridge-permission-auto-accept-runtime.ts:1-103` | 〔根本没做〕 |

### 块 8：不在屏幕前那一块（通知／回得来）

| 控件 | 谁家有 | 我们这边状态 |
|---|---|---|
| 点通知**直达那一条会话**（协议词表 7 种意图，认不出返回空不抛错；冷启动那份意图先存住） | openchamber `mobile/.../deepLinks.ts:23-30,46-53`、`deepLinkNavigation.ts:15-16,133-135` | 〔根本没做〕Go 侧搜协议/深链 0 命中 |
| 批准触发推送时**内容刻意空泛**（固定标题＋会话名，不带模型名/项目/正文） | openchamber `web/server/lib/notifications/APNS.md:16-21` | 〔根本没做〕（我们的"推送"是同进程里 Go 发给网页的一条消息，出不了这台机器） |
| "正在看就别吵"那道门的**判据只盯别的端、从不盯手机自己** | openchamber `HANDOFF.md:108-110`；且他们主动放弃了"问客户端在不在"那道门（`APNS.md:28-38` 注释写着网页容器有竞态） | 〔根本没做〕 |
| 子代理还在跑时，父任务"我做完了"那条通知**不发** | openchamber `notifications/DOCUMENTATION.md:51-53`（查不了才照发） | 〔根本没做〕 |
| 各家订阅额度那一屏（日/周/月/本轮/会话五类窗口＋重置倒计时＋余额＋消费上限） | openchamber `packages/vscode/src/quotaProviders.ts:1-60`（约 20 家，名册逐名见 `survey-2026-09-28-oc-mobile-vscode-extensions.md:37`） | 〔建了但没接〕我们每提供商有日/月配额两枚键（`internal/config/schema.go:367-368`）与日/月预算＋提醒阈值＋超预算动作（`:544-548`），**但没有那一屏** |

---

## §2 表二 · 设置页逐项（别人有／我们有／我们没有，各带证据）

> 我方"有没有"只按 Go 侧配置模型现读：`internal/config/schema.go` 今天共 **173 枚** `toml:"…"` 键、分 **19 个 section**
> （尺：`grep -cE 'toml:"[a-z_0-9]+"' internal/config/schema.go` ＝ 173；section 名册逐名 `:108-133`）。
> ⚠ 关键差别不在枚数：**别家每一枚设置项都在一扇能看见、能改、能"恢复默认"的页上；我们今天 173 枚全靠手改文本文件。**

### 甲·最响的一列差：**"哪一项被我改过"＋"恢复默认"＋"读不懂不许静默"**

| 项 | 谁家有 | 我们 |
|---|---|---|
| 用户覆盖过的字段带**「已覆盖」标签**、旁边一枚**「恢复默认」** | DeepSeek 三页都有，逐字：`ui-settings-agent-loop/README.zh.md:28`、`ui-settings-shell/README.zh.md:28`、`:67`（平台差异：PowerShell 多一枚 `pwshPath`） | 〔根本没做〕我们的键没有"生效值＝覆盖叠默认"这层显示（配置读进来就是一个值，界面无从知道是谁写的） |
| **点保存之前不写入任何东西**；离开页面即丢弃草稿；清空并保存＝重置；填了非数字则**阻止保存并在字段下说明原因** | DeepSeek 同两处 `:28` | 〔不适用／没这形〕我们是直接改文本文件，没有草稿这层；**后果**：改错了没有任何当场拦截 |
| 键位配置**读失败时禁用编辑并说明当前在用哪套**、内容损坏**先备份再修**、被更新版本写的就提示升级、**"恢复全部默认"不覆盖读不懂的数据** | DeepSeek `ui-shortcuts/README.zh.md:36` | 〔根本没做〕 |
| 设置分区导航超出高度时列表独立滚动、"设置"标题固定 | DeepSeek `ui-settings-general/README.zh.md:32` | 〔界面地界·不读〕 |

### 乙·通用／外观

| 设置项 | 谁家有（点名＋行） | 我们有没有 | 证据 |
|---|---|---|---|
| 界面语言 | DeepSeek 有本地化体系（`packages/client/README.i18n.yaml` 在 66 条目名册里）；minimax 每个文案包**中英两份**（`packages/tui/src/tui/features/session-mutation/copy.en.ts` 196 行／`copy.zh-Hans.ts` 173 行）；openchamber 扩展面板 12 语（`extensions/DOCUMENTATION.md:8`） | **有**〔已证〕 | `internal/config/schema.go:138` `app.language` 默认 `zh-CN` |
| 主题（深/浅） | DeepSeek `ui-theme`（主题与正文字号，名册第 47 枚 `survey-2026-09-28-dsh-ui-packages.md:81`）；openchamber 登录链接/打开链接会带当前主题（`ui-settings-account/README.zh.md:74`） | **有**〔已证〕 | `schema.go:140` `app.theme` 默认 `dark` |
| **正文字号** | DeepSeek `ui-theme` 明写"主题与正文字号" | **没有** | 我们只有球的尺寸（`schema.go:165` `ball.size`）与面板宽高缩放（`:532-538`），**没有任何字号键** |
| 开机自启 | 本轮点到行的只有我们自己（不写跨家断言） | **有**〔已证〕 | `schema.go:142` `app.autostart` 默认 false |
| 只允许开一份 | 同上 | **有**〔已证〕 | `schema.go:144` `app.single_instance` |
| 球的透明度／点击穿透／全屏时藏起来 | 别家没有对位物（他们是网页三栏，没有悬浮球） | **有**〔已证〕 | `schema.go:169-173` 三枚键 |
| 快捷键（召唤／静音／取消／开面板） | DeepSeek 除固定三键外**全部可录制／可清除／可恢复默认**，持久化 Web 写 `localStorage['dsh.keybindings.v1']`、桌面写 `userData/keybindings.json`（`ui-shortcuts/README.zh.md:36`＋`packages/client/shortcuts/src/client/storage.ts:6`；固定三键 `src/client/fixed.ts:10-17`）；openchamber 听写有可绑快捷键（`ComposerDictation.tsx:127`） | **有键、没有改键的页**〔建了但没接〕 | `schema.go:180-183` 四枚；热改桥的生产调用者今天**仍只有** `cmd/balldebug/main.go:237`，产品入口 0 处（本轮复量） |
| 快捷键速查窗（列出当前窗口能用的所有命令、按名/别名/键位搜） | DeepSeek `ui-shortcuts` 整包；`Mod+/` 打开 | **没有** | Go 侧无任何"可绑命令目录"类型 |

### 丙·模型与密钥

| 设置项 | 谁家有 | 我们有没有 | 证据 |
|---|---|---|---|
| 每提供商一行、一次只展开一张编辑卡片 | DeepSeek `ui-settings-models/README.zh.md:32` | **没有那一屏** | 我们的提供商是配置块（`schema.go:419` `providers` 映射） |
| **API 密钥以只写方式存入引用之下、界面永不回显**；只有确认引用已配置才用绿点、只有确认缺失才用红点 | DeepSeek `ui-settings-models/README.zh.md:42`、密钥字符集校验 `:78`、删行时只清派生目标 `:124` | **有引用、没有状态显示** | `schema.go:236`（语音）与 `:387`（提供商）都是 `api_key_ref`；本轮在 Go 侧搜"绿/红状态回显"形状未量〔未量〕 |
| 每枚模型的上下文窗口／最大输出／输入类型勾选（文本·图片，**至少留一种**） | DeepSeek `ui-settings-models/README.zh.md:50` | **有**〔已证〕 | `schema.go:354` `context_window`、`:356` `max_output_tokens`、`:352`+`:325-332` 能力位（含 `text/vision/audio_in/audio_out/realtime/thinking/fc/stream`——**比他们那两枚勾更全**） |
| 计费方式／单价（输入·输出·缓存·音频进·音频出） | DeepSeek 未读到；openchamber 有各家用量 | **有**〔已证〕 | `schema.go:361` `billing`、`:338-342` 五档单价 |
| 每提供商日/月配额 | openchamber 约 20 家用量窗口 `packages/vscode/src/quotaProviders.ts:15-52` | **有键、没有那一屏** | `schema.go:367-368` `quota_daily_micro`／`quota_monthly_micro` |
| 每分钟请求数／每分钟字数 | 未点到别家行 | **有**〔已证〕 | `schema.go:397-398` |
| 思考档位 | Pi 系分**"本次"与"设为默认"两种持久化**（`pi-upstream/.../interactive-mode.ts:5014-5040`）；minimax 无命令、只在自带密钥配置里（`packages/config/src/byok-config.ts:97-112`） | **有**〔已证，但没有"本次／长期"那一分法〕 | `schema.go:291` `thinking_intensity` 默认 off、`:359` `thinking_levels` 清单 |
| 首次运行引导（版本化内测声明＋按条件显示的凭据步骤，两个有序弹窗） | DeepSeek `ui-settings-models/README.zh.md:12,36` | **没有** | Go 侧无引导账本；`schema.go` 无对应键 |
| 账号登录/退出／意见反馈（问卷带预填的构建版本·语言·屏幕分辨率）／赠金通知／额度页 | DeepSeek `ui-settings-account/README.zh.md:39,41,51,60,74,81`（这枚包**还不在官方包地图里**，见 `survey-2026-09-28-dsh-ui-packages.md:409`） | **没有**（我们是自带密钥的单机产品，属取向不同） | — |

### 丁·智能体与门控

| 设置项 | 谁家有 | 我们有没有 | 证据 |
|---|---|---|---|
| 权限档位 | DeepSeek 默认两档＋自定义/自动，落会话事件、**重启仍在**（`packages/interaction/permission-presets/src/index.ts:188-199`〔腿报〕）；Step-Code 四档写 `config.toml`、**但 `--cycle` 故意不落盘**（`Step-Code/packages/coding-agent/src/step/features/step.ts:310-318`）；minimax **六档**＋规则源三层（`packages/agent-modules/permission/src/types.ts:13-19,38`）；openchamber 以服务端为准 | **有**〔已证〕且**默认最严** | `schema.go:470` `risk.permission_mode` 默认 `ask_every_step`；零值落最严档（`internal/risk/mode.go:44-56`） |
| "临时切档故意不落盘"这一形 | Step-Code 有（同上行号） | **没有** | 我们只有持久档 |
| 每会话自动接受 | openchamber `bridge-permission-auto-accept-runtime.ts:1-103` | **没有** | — |
| 允许目录 | minimax 钩子可 `addDirectories/removeDirectories`（`agent-modules/plugin-hooks/src/contracts.ts:120-126`） | **有**〔已证〕 | `schema.go:476` `allowed_dirs`（且属锁定段，`:118-120` 注明放宽需 L2 重新确认） |
| 删除动作开关 | 别家未点到行 | **有**〔已证〕 | `schema.go:481` `delete_enabled` 默认 false |
| 命令行开关／白名单／黑名单覆写 | 别家未点到行 | **有**〔已证〕 | `schema.go:455-461` |
| 等批准超时／最严档窗口 | 别家未点到行 | **有**〔已证〕 | `schema.go:450` `confirm_timeout_sec` 300、`:453` `l1_window_sec` |
| 重复刹车的三级阈值 | minimax 一枚 `remindAfterOccurrences`（必须 ≥3，不传＝只观测）`agent-modules/runaway-guard/src/guard.ts:28-31` | **有**〔已证〕 | `schema.go:425` `repeat_thresholds` 默认 `3,5,8` |
| **刹车的总开关（误伤时一键关）** | minimax `enabled` 默认 true、可远端覆盖、**读不到值＝不覆盖而不是打开**、关掉时清空累计（`packages/config/src/runaway-guard-config.ts:10-21`、`agent-extension/src/runaway-guard.ts:46-47,88-92`） | **没有** | 本轮复量：`grep -rniE "loop_guard.*enabled|LoopGuardEnabled"` 仍 **0 命中** |
| 并行工具调用数 | DeepSeek 做成设置页 `ui-settings-agent-loop/README.zh.md:12,28` | **没有**（冻结在桥里） | `internal/tools/bridge.go` 的 `MaxToolConcurrency`；`schema.go` 搜不到对应键 |
| 子代理最大深度／并行上限 | DeepSeek 一页两半一次保存、两次写入互相独立、过期草稿报冲突不覆盖 `ui-settings-subagent/README.zh.md:12,32`，字段与最小值 `src/client/subagent-limits-card-controller.ts:10-13,48`，**计口径写在数字旁边** `README.zh.md:28` | **没有**（编译期常量 4／1） | `internal/tools/subagent_197.go:73,79`；`schema.go` 搜 `subagent` 0 命中 |
| 子代理能用哪些模型（模型路由） | DeepSeek 同页右半 `ui-settings-subagent/README.zh.md:30`（"已保存但当前不可用"分组仍可移除、加载失败报出来且不藏别的提供方、给重试） | **没有** | — |
| 最大轮数／token 预算／单工具超时 | minimax 预算侧反而**没有**（`survey-2026-09-28-minimax-agent-modules.md:112` 逐维对照） | **有**〔已证〕 | `schema.go:431` `max_rounds` 50、`:433` `token_budget` 200000、`:436` `per_tool_timeout_ms` |
| 转向（人在中途插话）开关 | 别家未点到行 | **有**〔已证〕 | `schema.go:440` `steering_enabled` |
| 项目说明总开关 | 别家未点到行 | **有**〔已证〕 | `schema.go:444` `project_instructions_enabled` |
| 单步内并行上限之外：**命令超时／单条输出上限（超出转存临时文件）** | DeepSeek `ui-settings-shell/README.zh.md:12,28` | **半个**：工具超时有键；输出上限是**按 token 等比缩放**的冻结参考值、不是可配字节数 | `schema.go:436`；`internal/agent/budgets.go:22-47`（`refSpillTokens=4000` 一族） |

### 戊·语音／音频／隐私／记忆／面板／成本／插件

| 设置项 | 谁家有 | 我们有没有 | 证据 |
|---|---|---|---|
| 识别提供方与模型；播报提供方/音色/语速 | openchamber 有整屏（`VoiceSettings.tsx:42-66,383`：四模型带准确度/速度评分条、逐个下载删除） | **有键、没有那一屏** | `schema.go:210-218` |
| 唤醒词开关/关键词/灵敏度/**否决词** | 别家没有对位物（他们没有语音否决出口） | **有**〔已证〕 | `schema.go:197-203`（`veto_words` 默认 `取消,停下,别`） |
| 标点／对话模式／回声消除（参考信号）／实时链路／云端识别与播报的降级链 | 别家未逐一点到行 | **有**〔已证〕 | `schema.go:252-265`、`256`、`258` |
| 输入设备／采样率／半双工／麦克风默认静音 | 别家未点到行 | **有**〔已证〕 | `schema.go:271-277` |
| 路径脱敏／诊断上报开关／保留天数／是否留对话／是否留音频 | 未点到别家行 | **有**〔已证〕 | `schema.go:509-517` |
| 记忆层开关／容量上限／第三层保留天数／抽取模型 | 别家未点到行 | **有**〔已证〕 | `schema.go:522-525` |
| 面板开关／宽／高／会话内保活／缩放 | DeepSeek 有"宽度偏好记住一份、不逐会话"（`ui-layout/README.zh.md:89-92`） | **有**〔已证〕 | `schema.go:530-538` |
| 日预算／月预算／提醒阈值／超预算动作 | openchamber 有用量窗口（同上） | **有**〔已证〕 | `schema.go:544-548` |
| 插件开关／能力白名单／插件网络白名单／宿主接口白名单／二层插件开关 | DeepSeek 有**「内置插件」分区＋只读插件清单标签页**（`ui-settings-plugins/README.zh.md:12,28`、`ui-settings-plugin-inventory/README.zh.md:28,34,40,60`：短名显示规则、坏预设带标记、只改列表显示**不写任何设置**）；openchamber 扩展登记表**是空的**（`extensions/registry.json:1-3`） | **有键、没有那一屏** | `schema.go:556-573`（含 `host_api`、`net_allowlist`、`tier2_enabled`；`grep -c Plugin internal/config/schema.go`＝8） |
| **技能：开关／来源／优先级／谁把谁盖掉了** | minimax 六层来源＋优先级＋外部标记＋`losers`＋两级诊断＋两组可数指标（`agent-modules/skills/src/types.ts:1,14-22,49-93,100-115`）；DeepSeek 技能 8 桶＋免重启监听 | **完全没有——连配置键都没有** | 本轮复量：`grep -c Skill internal/config/schema.go` ＝ **0**（同文件 `Plugin` ＝ 8） |
| 模型分发：目录／镜像列表／是否验签／本地覆写 | 别家未点到行 | **有**〔已证〕 | `schema.go:582-590` |
| 日志级别／滚动（大小·天数）／资源采样间隔 | 别家未点到行 | **有**〔已证〕 | `schema.go:602-606` |
| 网页搜索提供方（密钥／接口地址／单次搜索次数上限） | DeepSeek `ui-settings-web-search/README.zh.md:2,68` | **没有这一 section** | `schema.go` 的 19 个 section 名册里没有 websearch；这条正好撞 AGENTS §2 那条未定案（`web.search` 实现路径） |
| **"写了也不管用的键"响亮拒收** | openchamber 有 6 处"文案翻好了 9×14 种语言、4 枚键进了注册表、界面上零控件"（本仓证据件 `docs/evidence/s1/219-approval-reply-surface-c1.md:424` 逐名列出） | **有，且比他们硬**〔已证〕 | `internal/config/unwired.go` 147 行、6 枚拒收键（尺 `grep -c "path:"`＝6），文件头三条件逐字 `Such a key lies: the user writes one line and the program swallows it.` |

---

## §3 表三 · 状态文案（同一件事各家显示成哪几句原话）

> 尺：下面每一句都是**文件里的原文**（英文串照抄英文、中文串照抄中文），不是我的转述。
> 各家路径都在 `D:\work\AI\open source\` 下；`openchamber` 无提交历史，它的每一行只能当"当下快照"引。

| 这件事 | 各家原话（点名＋`文件:行`） | 我们这边现在显示什么 |
|---|---|---|
| **它在动**（正在生成回复） | minimax：`Running`，前面一枚 `◇`（`minimax-code/packages/tui/src/tui/shell/activity-line.ts:241`）；Pi 系：`Thinking...`（`pi-upstream/packages/coding-agent/src/modes/interactive/interactive-mode.ts:455`，那是"把思考过程藏起来时用的默认标签"）；DeepSeek 的子代理行：状态点只有 进行中／完成／静止 三名（`deepseek-harness/packages/client/ui-subagent/src/SubagentHeaderLineage.tsx:334`） | 我们显示的是**状态名**不是短句：20 枚名字逐字 `FirstRun/Sleeping/Armed/Muted/Listening/Thinking/Acting/Speaking/Warm/Conversation/Confirming/AwaitingApproval/Settling/Downloading/Unconfigured/NoNetwork/Error/Queued/Stuck/WatchdogAlert`（`internal/statemachine/states.go:11-31`，上面那行注释逐字 "names exactly as in D43 / SPEC-08"）〔已证〕 |
| **正在压缩上下文** | minimax：`Compacting context`（`activity-line.ts:257`）与 `Compacting now`（`packages/tui/src/tui/transcript/context-visualization.ts:211`）；Pi 系那枚命令**收自定义指令**（`pi-upstream/.../interactive-mode.ts:6858`）；DeepSeek 是"遮蔽、不删除，零参数，只人触发"（`deepseek-harness/packages/compaction/compaction/src/index.ts:143`、`types.ts:115-119`〔腿报〕）；openchamber 走服务端压缩，另有"被钉住的内容不进压缩区"（转引 `survey-2026-09-29-command-catalog-across-harnesses.md:35`） | ⚠ **20 枚状态名里没有"正在压缩"这一名**（本轮逐名读过上面那张名册），界面无从显示这一刻。**这一格是本版新查出来的，第三版没写过** |
| **正在重试** | minimax：`Retrying model request`，或者把模型给的那句话接上（`activity-line.ts:245`）；Pi 系：`Retry failed after ${attempt} attempts: ${finalError}`（`interactive-mode.ts:3649`） | 有重试的**数**（`internal/config/schema.go:306-307` 最大次数与退避毫秒）〔已证＝键在并被读〕，**没有给用户看的那一句**〔仅文档〕 |
| **正在重连** | minimax：`Reconnecting`，后面用 ` · ` 拼上原因（`activity-line.ts:252`） | 有 `NoNetwork` 这枚状态名，**没有"正在重连"这一形** |
| **正在停**（人按了停止） | minimax：`Stopping response`，错误色＋图标 `!`（`activity-line.ts:259-260`）；DeepSeek 的停止控件是**两击**：首击上膛、3 秒内确认击才真杀，行状态先 `stopping` 再进已结束分组、detail 写 `cancelled by the user`（`deepseek-harness/packages/client/ui-jobs/README.zh.md:30`） | 我们的 `Stuck` 说的是"**它自己撞线被刹停**"，不是"**人叫它停**"；两击上膛在 Go 侧 0 个生产落点（本轮复量） |
| **出错了** | minimax：`Runtime error`，或 `Error · <那句话>`（`activity-line.ts:264`）；Pi 系兜底串 `Unknown error occurred`（`interactive-mode.ts:1181`）；openchamber：`Folder access is required.`（`openchamber/packages/ui/src/lib/i18n/messages/en.ts:2313`）、`Camera access is off. Enable it in Settings to scan a QR code.`（同文件 `:150`） | 有 `Error` 状态名＋错误分类（D37）〔已证＝状态名在表上〕；**界面上那句人话本轮没量到**〔界面地界·不读〕 |
| **此刻无事可做** | minimax：活动行**直接不渲染**（`activity-line.ts:141` 逐字 `if (safeWidth === 0 || this.state.phase === 'idle') return []`）；DeepSeek 后台任务那一格**没见到一枚任务就什么都不画**（`ui-jobs/README.zh.md:26` 逐字："普通对话不会为未使用的能力长出控件"） | 我们这颗球**恒有一个状态**（它是常驻视觉件，形状不同，不算缺）；⚠ 面板那一侧"没接上的能力会不会长出控件"这一维**只在数据层管住了**（"没值就不发那个键"，`internal/panel/pump.go:133-141`）〔已证〕 |
| **在等你批准** | openchamber：`Permission required`（`en.ts:619`）＋卡上三枚 `Allow once`／`Always…`／`Deny`（`packages/ui/src/components/chat/PermissionCard.tsx:407-442`）；minimax：给人看的双语三名 `已允许／需要确认／已拒绝`（`minimax-code/packages/agent-modules/permission/src/reason-format.ts:27-30`）；DeepSeek：面板只给 `allowed-once`／`rejected` 两枚（`deepseek-harness/packages/client/ui-approval/src/client/contract/slots.ts:66`） | 有两个状态名（`Confirming`／`AwaitingApproval`）〔已证〕，**但"仅这一次"这个动词在面板上没有出口**（见 §4.C）；四条应答出口不可用时各说实话：`语音取消不可用`／`面板取消不可用（票 37 未接入）`／`悬浮球取消不可用`／`Esc 取消不可用`（`internal/agent/approval/approval.go:106-114`）〔已证在产码里〕 |
| **拒绝了，还要告诉模型** | minimax：给模型那份**前缀固定英文** `[permission:denied source=<source>]`，注释逐字"这样模型和日志检索在任何界面语言下都能解析出来源"（`reason-format.ts:324-338`），来源四枚 `user`／`rule`／`safety`／`safety-immune`（`:306`） | 搜这一族串 **0 命中**〔根本没做〕 |
| **它在转圈** | minimax 给模型的提醒逐字：`[runaway guard] The same tool error family has now occurred ${occurrences} times in a row. Do not retry the same route unchanged. …`（`packages/agent-modules/runaway-guard/src/reminder.ts:72-101`），**每套都挂一段"别把这条存进记忆"的防污染后缀**（`:59-60`） | 撞线时给人的那一条在（`EvStuck`，`internal/agent/sink.go:30-32`＋`internal/agent/loop.go:470-483`）〔已证〕；**按信号分套、说给模型听的话术没有**；防污染后缀这一维本轮**仍没量**（第三版 4.7 也标"没查"，本版维持"没量"） |
| **附件相关的三句** | minimax：`无法预览 · 附件仍可发送`（`packages/tui/src/tui/features/composer/copy.zh-Hans.ts:5`）、`Cannot attach ${洗过的路径}: path is not a file.`（`attachments.ts:85-87`）、清单行 `N. 文件名 · mime · 人类可读字节数`（`:105-108`）；openchamber：`Anhänge sind zu groß zum Senden. Bitte versuche, die Anzahl oder Größe der Bilder zu reduzieren.`（`openchamber/packages/ui/src/lib/i18n/messages/de.ts:2240`） | **三句都没有**：受理器在、没人把结果发出去（`internal/panel/pump.go:226` 仍传空实参）〔建了但没接〕 |
| **在等一件它自己会好的事**（等待态怎么写） | openchamber：`Waiting for the dev server` ＋第二行解释 `It is not accepting connections yet. This page will load as soon as it does.`（`en.ts:1345-1346`）；另有 `Waiting for reviewer`／`Waiting for implementer`（`:1698-1699`） | **没有这种"两行式"等待文案**的先例；最接近的一句反而是我们自己的：`切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主）`（`internal/panel/git.go:82`）〔已证在产码里〕——**这句的形状其实比别家更诚实，但全仓只有这一处** |
| **正在加载／加载失败** | minimax：`正在加载会话历史…`＋`Esc 取消`＋失败时 `Enter 重试 · Esc 关闭`（`packages/tui/src/tui/features/session-mutation/copy.zh-Hans.ts:6-8`）、`无法加载历史消息。`（`:94`）；DeepSeek：加载失败的提供方**被报告出来而不藏起其他提供方**，并给**重试**（`deepseek-harness/packages/client/ui-settings-subagent/README.zh.md:30`） | 〔界面地界·不读〕Go 侧没有"会话历史页"这一物 |
| **这件事现在做不了，但告诉你怎么做才能做** | minimax：`停止当前回复后才能修改会话历史。`（`copy.zh-Hans.ts:23`）、`请先开始或恢复一个会话，再浏览历史。`（`:33`）、`请退出计划模式后再回退此会话。`（`:113`）、`请打开父会话后再执行回退。`（`:114`）；另一族**连退出路径一起给**：`${usage} is unavailable in side conversations. Press Ctrl+C to return to the main session first.`（`packages/tui/src/tui/commands/catalog.ts:760-766`） | **没有这一族句子**〔根本没做〕。我们的"做不了"是通道级实话（上面那四枚"取消不可用"），**没有"你该先做哪一步"那一半** |
| **一半成功一半失败时怎么说** | minimax：`已回退 {turns} · 工作区文件未变`／`已回退 {turns} · 已还原 {files}`（`copy.zh-Hans.ts:88-89`）、`会话已回退，但部分文件改动未能还原。{error}`（`:160`）、`回退已完成，但当前会话视图刷新失败。请重新打开会话并确认状态。`（`:109`） | **没有这种"分段承认"的文案**；我们最接近的做法是"键在、值缺就响亮报"（`internal/config/unwired.go`，147 行、6 枚拒收键）〔已证〕 |
| **改了配置但这次没生效，要不要告诉用户** | Step-Code：明说要先重启——`Restart Step to activate plugin contributions.`（`Step-Code/packages/coding-agent/src/step/plugins.ts:1185`）、`Restart Step to start the plugin's MCP server.`（`:1334`）、兜底 `No marketplaces configured.`（`:1395`）／`No marketplace named '${name}' is configured.`（`:999`）；DeepSeek：**免重启**（技能目录有监听器，`deepseek-harness/packages/skill-filesystem/src/index.ts:249-262`〔腿报〕）；openchamber：MCP **有健康状态**（`packages/ui/src/stores/useMcpStore.ts:34-40`〔腿报〕） | ⚠ **三家分两派，但"改了不生效也不告诉用户"这一形别家都不选**。我们是"靠重启生效"占多数（三档生效级别在 `docs/PLAN.md` 的 D36），**界面上没有任何一处会说出"这一枚要重启"** |
| **还没有会话时**（空态） | Pi 系：`Session cwd not found`（`pi-upstream/.../interactive-mode.ts:2680`）；minimax：`暂无历史用户消息。发送消息后即可使用编辑、回退或分支。`、`暂无可用的已持久化用户提示词。发送消息后请重试 /fork。`（`copy.zh-Hans.ts:46,135`） | 〔界面地界·不读〕⚠ Go 侧已有同形先例，值得点名给界面那支：**名册为空时仍然发那一节**，用空数组＋池数字说明"这次没派过任何孩子"——**"没有"与"看不见"不许画成同一个样子**（`internal/panel/composer.go` 的 `Tasks` 节注释＋`internal/panel/pump.go:133-141`）〔已证〕 |

> ⚠ 本表的"跨家"只在我**逐家点到行**的那些事件上成立。没点到的事件（例如"额度用尽""磁盘满了"）我不写"别家都／都没有"。

---

## §4 条目正文（三段式：①这是个什么东西（用户点哪里能看到）②别人怎么做 ③我们这边现在是什么状态）

### §4.A 聊天输入框与加号菜单
〔待填〕

### §4.B 会话列表与历史
〔待填〕

### §4.C 批准、权限与安全门控
〔待填〕

### §4.D 模型、密钥与服务商接入
〔待填〕

### §4.E 语音（听与说）
〔待填〕

### §4.F 工具、插件与扩展
〔待填〕

### §4.G 文件改动、差异与回退
〔待填〕

### §4.H 终端、命令与进程
〔待填〕

### §4.I 用量、成本与配额
〔待填〕

### §4.J 状态、通知与提醒
〔待填〕

### §4.K 设置页与偏好
〔待填〕

### §4.L 外观、主题与可读性
〔待填〕

### §4.M 斜杠命令与快捷指令
〔待填〕

### §4.N 其它（不好归块的）
〔待填〕

---

## §5 我上一版写错的条目（逐条对抗第三版：哪一行说了什么／今天读数是什么／正确说法）

> 第三版＝`docs/reports/missing-features-2026-09-28-v3.md`，772 行／67 条，本轮**全文读完**（口径：`wc -l` ＝772；本腿 1-260 与 261-772 两段各通读一次）。
> 下面 13 条里，**前 6 条是第三版自己的缺陷**（要改的），7-10 条是"说轻了／范围写歪"，11-13 条是**我今天复核过、仍然成立**的关键读数（列出来是为了让你知道我没偷懒）。

**5-1 〔推翻〕第三版根本没并回那份"聊天框加号菜单"的调研——那一整块今天才是第一次进清单。**
- 第三版说了什么：它末尾那张"原始件→节次"对照表（`missing-features-2026-09-28-v3.md:738-740`）只列了三份原始件（手机壳那家／52 个界面包那家／12 个智能体模块那家）。
- 今天读数：我把第三版全文按词搜了一遍——"**附件**"命中 **0** 次、"**斜杠**"**0** 次、"**加号**"**0** 次、"**同意**"**0** 次、英文 "export" **0** 次（"导出"仅 1 次）（尺：`grep -c` 逐个词，对 `missing-features-2026-09-28-v3.md` 全文）。
- 正确说法：那份 86,486 字节的加号菜单调研（`survey-2026-09-28-composer-plus-menu.md`）**一条都没进第三版**。聊天框那一块（附件、@引用、斜杠命令、技能开关、会话出口、同意页）在本版才是第一次出现——见本版 §4.A／§4.B／§4.M 与表一第一块。⚠ 顺带纠正派单的一句话：**那五份件不是"第三版之后落的"**，按文件时刻它们全在第三版之前（第三版＝09-29 12:55；五份＝09-28 19:52／19:59／20:50／23:30 与 09-29 10:46）；**真没并回的只有两份**（加号菜单那份＋命令目录那份），另三份第三版已经并了。

**5-2 〔推翻〕"那些参考仓都没有提交历史"这句是错的——六家里五家有，只有 openchamber 真没有。**
- 第三版说了什么：第三版正文第 555 行写"那份克隆'没有提交历史'这件事写进派单"；它并的两份原始件更是把这条当自我设限——`survey-2026-09-28-dsh-ui-packages.md:3`「该克隆没有 `.git`」、`survey-2026-09-28-composer-plus-menu.md:439`「五家的 `.git` 一律不存在（已确认）」。
- 今天读数：`Step-Code`／`deepseek-harness`／`minimax-code`／`pi`／`pi-upstream` **五家都有 `.git`**，HEAD 逐名是 `7dd66cb`(2026-09-24)／`477b4f4`(2026-09-24，那一发就是 0.1.7-rc.2 那条发布分支的合并)／`4198174`(2026-09-25)／`812b1f4`(2025-08-05)／`bf8e4b953`(2026-09-26)，五家 `git status --short` 都是干净（行数 0）；只有 `openchamber` 没有。
- 正确说法：**"只能带 `文件:行`、不许引历史"这条限制只对 openchamber 成立**。第三版那些结论本身全是 `文件:行`，**不受影响、不用改**；但以后谁再说"人家没历史所以核不了版本"，那是错的。

**5-3 〔口径错〕"DeepSeek 那家 52 个界面包"——第三版没说这数出自哪一份 DeepSeek。**
- 第三版说了什么：第 7 行「DeepSeek 那家的 **52 个界面分包**」。
- 今天读数：本机其实有**三枚** DeepSeek 对象——① 源码克隆 `D:\work\AI\open source\deepseek-harness\package.json` 的 `version` ＝ **`0.1.7-rc.2`**（有 `.git`，就是这一枚的 `packages/client` 下 `ui-` 目录数 ＝ 52，尺 `ls packages/client | grep -c '^ui-'`）；② npm 全局 `@deepseek-ai/dsh` ＝ **`0.1.0-rc.6`**（尺：`grep -m1 '"version"'` 于 `$APPDATA/npm/node_modules/@deepseek-ai/dsh/package.json`）；③ `~/.dsh/profiles/node_modules/@deepseek-ai` ＝ **257 枚**（尺 `ls | wc -l`；`node_modules` 顶层 276 枚）。派单说的"0.1.6-alpha.1"这个号**我今天没验出来**（`~/.dsh/profiles` 根目录没有 `package.json`，只有 `desktop/`／`web/`／`node_modules/` 三项），所以那半句我不复述。
- 正确说法：本版（和第三版）所有 DeepSeek 的行号都出自 **①源码克隆 0.1.7-rc.2**，**不是任何一套装机版**。引用时必须带上这句。

**5-4 〔数错，而且是分母写歪〕"移动那一面 61 枚命名件"。**
- 第三版说了什么：第 766 行「`packages/ui/src/apps/` **61 枚**移动命名件只全读 6 枚」。
- 今天读数：`openchamber/packages/ui/src/apps/` 目录里**一共只有 52 枚文件**（尺 `find .../apps -type f | wc -l` ＝52），其中文件名带"mobile/widget"的是 **36 枚**；把范围放大到整个 `packages/ui/src`（3,187 枚文件）才是 **62 枚**（尺 `find openchamber/packages/ui/src \( -iname '*mobile*' -o -iname '*widget*' \) -type f | wc -l`）。
- 正确说法：**"61"这个数在任何一把尺下都复现不出来**，而"只全读 6 枚"这个覆盖率是照着那个错分母说的。本版改写为：**整个 `ui/src` 有 62 枚带移动字样的件，其中 `apps/` 目录内 36 枚**。

**5-5 〔说轻了〕"网页到 Go 之间那四道门，只接通了一道"。**
- 第三版说了什么：第 118 行那一格把"〔做了但只在命令行里能跑〕"的例子写成"那四道门只接通了一道"。
- 今天读数：装配点**真的在生产命令里**（`cmd/wisp/panel_inbound.go:232` 构出 `panel.ComposerDispatch`）：档位那一道**装了真处理器**（`:224-230`），另外三道**是明写的空**并且注释里点名各挂哪张工单（`:234-237`：工作区＝186／附件＝92／发消息＝35）。更关键：空**不是沉默**——`internal/panel/composer_dispatch.go:62` 那句逐字是"**该方法在名册内，但本机未接入处理器**"。
- 正确说法：不是"缺一扇门"，是"**四道门都铸好了、三道故意空着并且会报名字**"。这一格的动作是**补处理器**，比"补门"便宜，而且用户侧已经有话可看。

**5-6 〔说过头〕"这是我们唯一比他们强的地方"。**
- 第三版说了什么：第 26 行与第 233 行两处把"我们有 6 枚写了也不管用的键、直接响亮拒收"称为**唯一**一处比对方强。
- 今天读数：第三版自己在第 37 行就写了另一种病（对方有几处文案**根本没走翻译系统、写死英文**），而我本轮在 `internal/config/unwired.go`（147 行、6 枚拒收键，尺 `grep -c "path:"`）复核的那一半仍然成立。
- 正确说法：**两边各有两种病**——他们"文案在、控件不在"和"该翻的没翻"；我们"键在、读者缺"这一面今天也一样存在（附件、名册、技能清单三样都是"字段在、没人生产"）。把"唯一"改成"**这一格**我们比他们硬"。

**5-7 〔范围写歪，本版收窄〕"全仓搜不到'胜者败者名单'的对应物"。**
- 第三版说了什么：第 579 行已经老实写了"只搜了两个目录，0 命中 ⇒ 不能写全仓为 0"。
- 今天读数：我把同一把尺放大到全仓——`grep -rniE "loser|shadow" --include=*.go internal cmd`（排测试）**仍为 0 命中**；这一格现在有了更大的搜索面支持。
- 正确说法：可以从"限定范围内为 0"升级为"**全仓为 0**"，但这仍然是"没做"，不是"做得比他们差"。

**5-8 〔本版新增、第三版整块缺失〕菜单项背后"是真干活还是一句提示词"——第三版一栏都没有。**
- 第三版说了什么：没有。第 118 行那张表把"没有"分成三种，**没有第四种："有一行字，但那行字背后不是功能"**。
- 今天读数（逐家点名）：openchamber 那 9 枚（概要／做计划／立目标／定时／追赶／调试／权衡／探索／工作区复盘）**背后就是两串提示词**（给用户看一句＋给模型塞一段），零实现（`openchamber/packages/ui/src/components/chat/composer/submit/slashCommands.ts:22-40`，转引自 `survey-2026-09-29-command-catalog-across-harnesses.md:23`）；minimax 的 `/review` 是一句常量 `'Please review my uncommitted changes.'`（`minimax-code/packages/tui/src/tui/controller/product/command-flow.ts:142` 与 `packages/tui/src/headless/invocation.ts:234`，两处独立命中，抽验过）；DeepSeek 的 `/permission` 才是**真写路径**（解析预设 ⇒ 同时写沙箱与审批两枚旋钮 ⇒ 落一条会话事件，`deepseek-harness/packages/interaction/permission-presets/src/index.ts:258,264-274,403-418`，〔腿报〕级）。
- 正确说法：**这一栏必须进表一**，否则用户点错预期要我们背——"会改系统状态的"和"只是把一段话递给模型"是两种东西，得在界面上分开。

**5-9 〔本版新增〕第三版从没写过"会话出口"（导出／反馈／分享）。**
- 今天读数：DeepSeek 是 ZIP 走一条带鉴权的独立下载路由、命令本身只回"已请求"并**明确拒绝带路径参数**（`deepseek-harness/packages/session-query/session-log-export/src/index.ts:78-90`）；Pi 系是 HTML／JSONL 导出＋一键分享成私密 gist（`pi-upstream/packages/coding-agent/src/core/slash-commands.ts:25-27`）；minimax 有 `exportCurrentTranscript`（`minimax-code/packages/tui/src/.../feature-flow.ts:278`，〔腿报〕级）；Step-Code 的反馈包是**五家里最完整的**：15 个文件 3,101 行，先列"会发出去哪些文件"的清单再要同意、先脱敏再限长、每个外部字段各洗一遍控制字符（`Step-Code/packages/coding-agent/src/step/feedback/bundle.ts:110,176-185`、`consent.ts:73-92`）。
- 我方状态：〔根本没做〕——本仓 Go 侧搜不到任何"导出这次对话"的出口。
- 正确说法：这是本版**表一"会话这一块"里的第一行新增**，也是命令目录那份件点名"我们零枚"的那一格。

**5-10 〔本版新增〕"改了什么但这次没生效"这件事，别人做成了一张能数出数的表；第三版只写了"要有败者名单"。**
- 今天读数：minimax 的技能快照里同时有**赢的／输的（为什么输、赢的那枚在哪）／诊断（警告还是错）／扫描了几枚·读了几枚·复用了几枚／渲染了几枚·压成摘要几枚·丢掉几枚·预算用了多少字**（`minimax-code/packages/agent-modules/skills/src/types.ts:49-93,100-115`）。
- 我方状态：〔根本没做〕——`internal/config/schema.go` 里"Skill"字样今天**仍是 0 次**（"Plugin"仍 8 次），**技能这一维连配置键都没有**，谈不上清单。
- 正确说法：第三版第 577-580 行那一格只抄了"败者名单"一枚，**漏了另外两组计数**（刷新计数与渲染计数）；owner 那句"我配了呀，怎么没生效"要的恰恰是漏掉的那两组。

**5-11 〔复核仍然成立〕会话那本账本今天还是没有。**
- 第三版说了什么：第 165 行"那个包目录下只有一枚文件……往外冒的消息只有 9 种名字，而且没有序号、没有时间、没有轮次编号"。
- 今天读数：`ls internal/session/` ＝ **只有 `doc.go`**；`internal/agent/sink.go:20-38` 枚举逐名数出来仍是 **9 种**（`text_delta/reasoning_delta/tool_start/tool_end/reminder/stuck/control/error/done`），`Event` 结构体里 `Seq|Timestamp|TurnID` 三枚**仍 0 命中**。
- 正确说法：**这一条不是我的错，是本版确认它还没被补上**——它仍是"点进去那一页／耗时／花费／回放"四件事的共同前置。

**5-12 〔复核仍然成立〕防转圈今天还是只能看见一种转圈。**
- 今天读数：`internal/agent/guard.go` 里比的仍是"把这一轮所有动作拼成一条长字符串"（`turnSignature` 定义 `:234`，调用点 `:166`）；搜"ABAB／错误家族／按工具名豁免"这三族 **0 命中**；`policy_version|strategy_version` 全仓仍 **0 命中**。
- 正确说法：第三版第四节 4.1／4.2／4.3／4.12 那几条**一天过去没有被任何在飞的票改掉**，仍是有效缺口。

**5-13 〔复核仍然成立，但要多写一句〕"名册里许诺了一枚没注册的能力"。**
- 第三版说了什么：**没写过这一形**（全文搜 `task.cancel` ＝ 0 命中）。
- 今天读数：派子代理那枚工具的说明文字**许诺了"可以用 task.cancel 单独停它"**（`internal/tools/subagent_197.go:196`），而 `task.cancel` 在名册头上逐字写着 **DEFERRED**（`internal/tools/task.go:23`；`task.list` 同一条 `:22`）。
- 正确说法：本版把它列进 §4.C，并点名——**命令目录那份件恰好写了一条判据正对着这一形**："能力没落地就不进名册，别家是用机器闸门挡住的"（`survey-2026-09-29-command-catalog-across-harnesses.md:47`）。票 221 在飞，我只登记不越权。

---

## §6 这一版仍未覆盖的地方（具名）

〔待填〕

---

## §7 本轮读数的分母与实读量（口径附命令）

**甲·参照仓分母（`D:\work\AI\open source\`，尺＝`find <家> -type f | wc -l`；代码件另给一把尺）**

| 家 | 全部文件 | `.ts`／`.tsx`（排 `node_modules`／`dist`） | `.go` | 有 `.git`？ | 本轮读到的深度 |
|---|---|---|---|---|---|
| Step-Code | 1,438 | 1,153 | 0 | 有 | 定向读：`step/feedback/{bundle,consent}.ts`、`step/plugins.ts`、`step/slash-commands.ts`、`step/features/step.ts`、`core/slash-commands.ts`（抽点＋逐行读被引那几段） |
| deepseek-harness（克隆 0.1.7-rc.2） | 13,880 | 5,146 | 0 | 有 | `packages/client` 66 条目全名册（52 枚 `ui-*`＋14 枚非 `ui-*`，逐名列出）；本轮**新增读穿 8 枚设置类包的 README**（见下方丙） |
| minimax-code | 4,287 | 3,261 | 0 | 有 | `packages` 15 枚全名册；`agent-modules` 12 枚（第三版已并）；**本轮新增**＝`agent-modules` 之外的名册与分枚数＋`packages/tui` 的 6 枚文案件（`copy.en.ts`／`copy.zh-Hans.ts` 各两族，共 400 行，其中 173 行整枚通读） |
| openchamber（**无 `.git`**） | 5,659 | 2,234 | 0 | **没有** | `packages` 8 枚全名册；`packages/ui/src` ＝ **3,187 枚文件**（568 `.tsx`＋1,418 `.ts`）；`ui/src/apps` ＝ 52 枚（带移动字样 36 枚）；`components/sections` ＝ 21 项名册；`lib/i18n/messages/en.ts` 定向抽句；本轮**未通读**任何一枚巨件 |
| pi | 37 | 0 | 0 | 有 | **不是编码助手**（本轮复核：37 枚文件、0 枚 ts，与 `survey-2026-09-29-command-catalog-across-harnesses.md:17` 的读数同值）⇒ 本版不引它 |
| pi-upstream | 1,969 | 1,597 | 0 | 有 | 定向读：`core/slash-commands.ts`、`modes/interactive/interactive-mode.ts`（抽 6 段带行）、`packages/tui/src/autocomplete.ts` |

**乙·那六份本仓调研件（尺＝`wc -l`）与本轮实读**

| 文件 | 行数 | 本轮实读 |
|---|---|---|
| `missing-features-2026-09-28-v3.md`（上一版） | 772 | **全读**（两段各通读一次） |
| `survey-2026-09-28-composer-plus-menu.md` | 557 | **全读**（`:14-398` 与 `:399-557` 两段） |
| `survey-2026-09-28-dsh-ui-packages.md` | 431 | **全读** |
| `survey-2026-09-28-minimax-agent-modules.md` | 801 | **读 `:36-801`**（前面 `:1-35` 是标题与分母节，只抽过结构与读数，未通读正文） |
| `survey-2026-09-28-oc-mobile-vscode-extensions.md` | 148 | **全读** |
| `survey-2026-09-29-command-catalog-across-harnesses.md` | 61 | **全读** |

⇒ **本腿只"并回"了这些件已经调研过的东西**，没有重复它们的调研；我只在下面三处亲自复核了 `文件:行`：
① 五家 `.git` 有无；② `61 枚` 那个分母（复现不出来，见 §5-4）；③ DeepSeek 三枚对象与版本（见 §5-3）。

**丙·本轮新读的参照面（不是重抄旧件，是打未覆盖面）**

- DeepSeek 8 枚设置类包的说明文件：`ui-settings-agent-loop`、`ui-settings-shell`、`ui-settings-subagent`、`ui-settings-web-search`、
  `ui-settings-models`、`ui-settings-general`、`ui-settings-plugins`、`ui-settings-plugin-inventory`、`ui-settings-account`
  ⇒ 尺＝`Grep` 于 `deepseek-harness/packages/client/ui-settings-*/README.zh.md`，命中的每一条都进了 §2 表二。
  ⚠ 这些是**说明文件**不是实现；`ui-*` 的 `src/**` 本轮一枚没读（第三版腿 2 同样没读，两版都不许当成"读过实现"）。
- minimax 的界面文案双份件：`packages/tui/src/tui/features/session-mutation/{copy.en.ts,copy.zh-Hans.ts,copy.ts}`
  与 `features/composer/{copy.en.ts,copy.ts,copy.zh-Hans.ts}`（尺＝`find packages/tui -name 'copy*.ts'`＝**6 枚**，行数见上）
  ⇒ 这一族是 §3 表三的主要出处，第三版**完全没有引过**（它的 minimax 部分只读了 `agent-modules`）。
- minimax 活动行：`packages/tui/src/tui/shell/activity-line.ts`（状态词表与逐相位文案）。
- `minimax-code/third_party/pi-mono/`（本轮偶然发现：**minimax 仓里带着一份 Pi 的第三方副本**）
  ⇒ 影响：任何"minimax 有 X"的结论都要先确认不是从这份 Pi 副本里读来的。本轮 §3 表三的 minimax 引文**全部出自 `packages/` 下**，未取自 `third_party/`。

**丁·我方一侧本轮的现读（尺全部是 grep/ls/awk，未跑任何编译或测试）**

| 量 | 今天读数 | 尺 |
|---|---|---|
| 配置项总数 | **173 枚键／19 个 section** | `grep -cE 'toml:"[a-z_0-9]+"' internal/config/schema.go`；section 名册 `internal/config/schema.go:108-133` |
| 状态名总数 | **20 枚** | `internal/statemachine/states.go:11-31` 逐名读 |
| 往外冒的消息种类 | **9 种**、`Event` 里无 序号／时间／轮次 | `sed -n '18,45p' internal/agent/sink.go` 逐名数＋`awk 'NR>=44&&NR<=70'` 搜三枚字段名 0 命中 |
| 会话持久层 | 目录下**只有 1 枚文件**（`doc.go`） | `ls internal/session/` |
| 名册行的字段 | 仍是 文本／产物路径／状态／父亲／名字／种类 **六枚**，无花费无计时 | `awk '/^type TaskOutput struct/,/^}/' internal/tools/task.go` |
| 刷给界面那包内容的键 | **6 枚**（待答／结果／输入框／生成时刻／项目说明／名册） | `awk '/^type Snapshot struct/,/^}/' internal/panel/composer.go` |
| 输入框载体的字段 | **9 枚**（含 4 枚附件相关） | `awk '/^type ComposerState struct/,/^}/' internal/panel/composer.go \| grep -c "json:"` |
| 附件受理器的生产调用者 | **0**（只有定义那一行） | `grep -rn "Ingest(" --include=*.go internal cmd \| grep -v _test` |
| 键位热改桥的生产调用者 | **1 处，且在调试工具里** | `grep -rn "NewHotkeyReloader" --include=*.go . \| grep -v _test` |
| 命令名册／斜杠命令目录 | **0 命中** | `grep -rln --include=*.go -iE "slashcommand\|commandcatalog\|CommandRegistry" internal/ cmd/` |
| 技能配置键 | **0**（同文件插件键 8） | `grep -c Skill`／`grep -c Plugin internal/config/schema.go` |
| 刹车粒度 | 仍是整轮签名 | `grep -n "turnSignature" internal/agent/guard.go`（定义 `:234`、调用 `:166`） |
| 刹车总开关 | **0 命中** | `grep -rniE "loop_guard.*enabled\|LoopGuardEnabled" --include=*.go internal/` |
| 策略版本号 | **0 命中** | `grep -rniE "policy_version\|strategyVersion\|strategy_version" --include=*.go internal/` |
| 送达账本／"丢了"状态 | **0 命中** | `grep -rniE "deliveredAt\|delivered_at\|TaskLost" --include=*.go internal/tools internal/memory internal/agent` |
| 反问用户／逐事件台账 | **0 命中** | `grep -rniE "ask_user\|askuser\|user_question\|trajectory" --include=*.go internal/ cmd/` |
| 面板入向四道门 | 装配在产命令里，**1 道有处理器、3 道明写空** | `cmd/wisp/panel_inbound.go:232-240`＋`internal/panel/composer_dispatch.go:62` |
