# Wisp 缺失功能清单 · 第四版（2026-09-29 v4）

> 调研腿：`v4-survey-1`（只读）。上一版＝`docs/reports/missing-features-2026-09-28-v3.md`（125,852 字节／67 条）。
> 这一版做三件事：**并回**第三版真正没并的那两份调研件、**挑自己上一版的漏**（§5，13 条）、**打还没覆盖的地方**（§4＋§6）。
> ⚠ 起手就纠正派单一句话：那五份调研件**不是"第三版之后落的"**——按文件时刻五份全在第三版之前；第三版并已并了其中三份，**真没并回的只有两份**（加号菜单那份、命令目录那份）。证据与推导见 §5-1。
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

**A-1〔根本没做〕加号点开的那张菜单，我们整个没有。**
① 用户点哪里：输入框左边那枚 `+`，点开是一列"往这次对话里加点东西"的选项（选本地文件、立个目标、进计划模式、提反馈……），或者直接打字 `/` 弹同一张表。
② 别人怎么做：DeepSeek 把这张表分两段、段名与行名**硬编码在前端 8 个字符串**里（`deepseek-harness/packages/client/ui-commands/src/client/presentation.ts:19-22`），段标题走翻译键（`:82-84`）；那枚 `+` 按钮和打字的 `/` **调的是同一张表**（`ui-conversation/src/client/apply.ts:505-516`），而且**这张表的服务不在时按钮直接变灰**、不是画出来再点了报错（`ui-conversation/src/client/skeleton/InputBar.tsx:427`）。
③ 我们现在是什么状态：Go 侧搜"命令名册／命令目录／注册表"三族词 **0 命中**（本轮尺见 §7 丁），所以这张表**连数据源都不存在**。

**A-2〔建了但没接〕往对话里塞文件这一件事：机器装好了，没人按开关。**
① 用户点哪里：加号 →「本地文件」；或者把文件拖进来、或者在文字里打 `@` 指一个工作区里的文件。
② 别人怎么做：清单每一行固定三段「第几枚 · 文件名 · 类型 · 大小」（minimax `packages/tui/src/tui/features/composer/attachments.ts:105-108`）；**"预览不出来"与"没附上"是两句话**（同目录 `copy.zh-Hans.ts:5` 逐字 `无法预览 · 附件仍可发送`）；上限两层都判、**以后端为准**（DeepSeek 注释逐字 `Attachment admission is enforced here, not in the composer`，`packages/interaction/commands/src/index.ts:353-357`）；被拒时的句子三样齐全（哪个文件／为什么／怎么办）。
③ 我们现在是什么状态：受理器**410 行在仓里**，类型判定那半**写得比别家对**（"用的是嗅出来的类型，绝不回显调用方自称的类型"，`internal/panel/attachments.go:87-90`）；**但生产装配点把附件那一枚实参写死成空**（`internal/panel/pump.go:226`），入向那道门明写着空并会报名字（`cmd/wisp/panel_inbound.go:236`＋`internal/panel/composer_dispatch.go:62`）。⇒ **缺的是生产者，不是字段**（输入框载体的 9 枚字段里有 4 枚就是附件用的）。
⚠ 这三行是定性，别把它读成安全事件：
- 现象出现在哪：附件／引用的**显示层**——文件名或扩展名由外部提供，若显示前不洗，一行字可以被一枚文件名改掉（别家为此各写了一份洗法：minimax `rendering/terminal-text.ts:15-26`、Step-Code `step/feedback/consent.ts:73-92`）。
- 有没有本机已被入侵的证据：**没有**。这是**功能缺口**，本腿在本机没读到任何被入侵的痕迹。
- 最坏后果是什么形状：等这条路真接上之后，一枚精心起名的文件能把面板那一行显示成别的话（用户误读），**不是**"现在就在漏"。

**A-3〔根本没做〕粘贴一大段文字，它自己变成附件。**
① 用户点哪里：在输入框里粘一段长文字，界面问一句"要贴进来还是当附件"。
② 别人怎么做：openchamber 阈值是"约 2000 字符或 25 行"，按设置项 `largeTextPasteBehavior`（问我／直接附／原地贴）处理，attach 时造一枚内存里的 `pasted-context-N.txt` 并**与手选 txt 走完全同一条管线**；而且**"问我"那两个按钮读的是当时的实时状态**，防"提示还挂着、内容已经变了"（`openchamber/packages/ui/src/components/chat/composer/DOCUMENTATION.md:137-148`）。
③ 我们现在是什么状态：〔根本没做〕。

### §4.B 会话列表与历史

**B-1〔根本没做〕"把对话退回到某一条"，而且能选退不提文件。**
① 用户点哪里：会话历史里选中某一条，选「回退」。
② 别人怎么做：minimax 把它做成**两枚分开的动作**——`回退会话`（"移除后续消息，工作区文件保持不变"）与 `回退会话和文件`（"并还原可安全恢复的文件改动"），各有说明（`packages/tui/src/tui/features/session-mutation/copy.zh-Hans.ts:28-31`）；退之前先给一张预览：`{ready} 可回退 · {skipped} 已跳过`（`:71,158`）；退到一半失败各说各话：`会话已回退，但部分文件改动未能还原。`（`:160`）。
③ 我们现在是什么状态：〔根本没做〕。

**B-2〔根本没做〕从任意一条历史消息开一枚分支。**
① 用户点哪里：选中一条 →「从这里创建分支」，说明逐字写着"创建子会话，不修改当前会话或文件"（minimax `copy.zh-Hans.ts:24-25`）；确认页把"当前会话／分支点／工作区是否会隔离"三行摆出来（`:136-144`）。Pi 系有 `/fork`、`/clone` 两枚命令（`pi-upstream/packages/coding-agent/src/core/slash-commands.ts:18-43` 名册内）。
③ 我们现在是什么状态：〔根本没做〕。⚠ 别和"子代理"混：**这是给人分的对话岔路，不是给它分的活**。

**B-3〔根本没做〕编辑某一条历史提问、从这里重新生成。**
minimax 连"中途失败怎么办"都分好几种说法：`已保存的编辑操作无法恢复，请使用新的 operation id 再次提交。`、`历史已回退，但修改后的提示词未提交。编辑草稿已保留，按 Enter 将其作为新消息发送。`（`copy.zh-Hans.ts:118-123`）。我们这边〔根本没做〕。

**B-4〔根本没做〕会话的出口：导出／复制／提反馈。**
① 用户点哪里：会话菜单里的「导出」「复制最后一条」「反馈」。
② 别人怎么做：DeepSeek 的导出命令**只回一句"已请求"**，真文件走另一条**带鉴权**的下载路由，而且**明确拒绝带路径参数**（`deepseek-harness/packages/session-query/session-log-export/src/index.ts:78-90`）；Pi 系默认导出 HTML、后缀决定格式，还能一键分享成私密 gist（`pi-upstream/.../slash-commands.ts:25-27`）；Step-Code 的反馈包是五家里最完整的：**先把"会发出去哪些文件"的清单摆出来再要同意**（`Step-Code/packages/coding-agent/src/step/feedback/bundle.ts:110,113`）、**先脱敏再限长、并写明是怎么截的**（`:176-185`）、超出上限先缩一次、再缩不动才放弃且给一个**可枚举的放弃理由**（`:88-97`）。
③ 我们现在是什么状态：〔根本没做〕。第三版整版没写过这一条（§5-9）。

**B-5〔根本没做〕会话列表按"需要你关注"排、带未读点。**
openchamber 的判据逐字 `needsAttention = unseen>0 && (!isSubtask || notifyOnSubtasks)`，**锁屏小组件用的是同一句**（`openchamber/mobile/.../mobileWidgetSnapshot.ts:18-20,84-98`）⇒ 一个判据两处用、不会两处打架。我们〔根本没做〕。

### §4.C 批准、权限与安全门控

**C-1〔根本没做〕卡上没有"就这一次"这个动词。**
① 用户点哪里：批准卡上「仅本次允许」。
② 别人怎么做：DeepSeek 面板只有两枚决定 `allowed-once`／`rejected`（`ui-approval/src/client/contract/slots.ts:66`），并划清界"持久权限策略仍由服务侧拥有"（`ui-approval/README.zh.md:36`）；openchamber 画三枚并给了快捷键（`packages/ui/src/components/chat/PermissionCard.tsx:407-442`）。
③ 我们现在是什么状态：**面板那一侧只能拒绝**，"允许"被结构禁掉（`docs/specs/SPEC-08*.md:156-176` C17 名册逐字）。⚠ 这不是"我们更严"这么简单：**真正该拍的是"允许这个出口该落在哪一层"**（原生卡片？球的长按？语音取消词的反义？），这一问碰 `Q-49`，属 owner 一句话，**我不自行填**。

**C-2〔仅文档〕点"一直允许"时，卡上没告诉你它会存成哪条规则，也没给你挑宽度。**
minimax 把"点一直允许"做成**一次列出几档宽度**（最窄／按第一个词／按前两个参数／按域名／整个工具），并且**第一档永远等于最窄的那档**（兼容律逐字写在 `packages/agent-modules/permission/src/types.ts:96-106`）；openchamber 直接把要存的规则文本印在按钮上（`PermissionCard.tsx:349-358`）。我们只有 `docs/PLAN.md:2028` 手写的三档、且只针对一枚工具。

**C-3〔根本没做〕它拿不准的时候不会回头问你一句。**
① 用户点哪里：它弹一张"请选择／自己写／跳过"的小卡。
② 别人怎么做：DeepSeek 整包 `ui-user-questions`，问题带选项、可多选、可自己写；**"跳过"与"关掉"是两种东西**（跳过＝交一份空答案继续走；关掉＝取消整个等待，两枚不同取消码）`packages/interaction/user-questions/src/types.ts:36-67`＋`ui-user-questions/README.zh.md:32`；题号与已选内容**只放不持久的临时存储、绝不写盘**（`:40`）；**同一时刻只有一枚问题拥有输入框**（`:95`）；计划评审那张卡的批准项**按名字点名、不按位置**（`types.ts:27-29` 逐字"免得界面从顺序里猜结论"）。
③ 我们现在是什么状态：`ask_user` 一族在 Go 侧今天仍 **0 命中**（本轮复量）。⇒ 它只能猜，猜错白烧一轮，**而且烧完不会告诉你它当初在赌**。

**C-4〔根本没做〕拒绝的三种心智，以及"新判据必须有文案"的编译期钉子。**
minimax 给人看的那份本地化（三名 `已允许／需要确认／已拒绝`，`reason-format.ts:27-30`）、**给模型看的那份前缀固定英文** `[permission:denied source=<source>]`（`:324-338`），来源四枚 `user/rule/safety/safety-immune`（`:306`）；并且表头逐字写着"**新增一条判据就必须同时新增它的文案，漏一条编译不过**"（`:7-9`）。我们〔根本没做〕，而这枚"漏一条就编译不过"的钉子**成本极低**（Go 侧等价物＝穷尽 switch）。

**C-5〔根本没做〕现场没有能答的人，就当次判拒绝。**
minimax 有一档专门是这个取值，理由文案逐字 `permission requires user confirmation but no approval channel is available`（`packages/agent-modules/permission/src/permission-core.ts:422-430`）；openchamber 的对应档是**反过来：判定器没给答复就转成等人**（`packages/web/server/lib/permission-auto-accept/modes.js:5-6`）——**两家在这一格方向相反**，所以没有现成答案可抄。我们这边：`internal/risk/mode.go:44-56` 有"零值落最严档"〔已证〕，**但"没有可用的批准通道"这一支 0 命中**（本轮复量）。⇒ 正面撞上我们那枚"深度睡眠时该问谁"与票 197"子代理卡在等批准时算什么"。

**C-6〔取向不同／反向标〕两处我们比他们强，别倒过来写。**
① 可回收删除：minimax 自己写着 Windows 的 `del`/`rd`/`Remove-Item` 三条**还没接改写通路、先用硬拒顶着**（`packages/agent-modules/permission/src/classifier/dangerous-patterns.ts:140-145`）；**我们那条通路已经做完并且会回读回收站记录证明真进了回收站**（`internal/tools/recycle_windows.go:17-28` 逐字 `os.Remove is never reached from this file`）〔已证〕。
② 批准卡的信息量：我们 11 枚字段 vs 他们 5 枚（`internal/panel/approval.go:39-59`），多的是"命中了哪几条规则""没有理由时不许渲染成没风险""这一票来自哪一侧"；应答出口我们 4 条闭集、其中**语音否决**在他们那 52 枚界面包里**没有任何对应物**〔已证：本轮逐名读过那 52 枚名册〕。

### §4.D 模型、密钥与服务商接入

**D-1〔建了但没接〕"每个提供商一行、密钥只写不回显、有没有配上一眼看得颜色"那一屏。**
DeepSeek 的模型设置页：每提供商一行、一次只展开一张卡片（`deepseek-harness/packages/client/ui-settings-models/README.zh.md:32`）；**密钥以只写方式存进引用之下、界面永不回显**，只有确认引用已配置才绿点、确认缺失才红点（`:42`）；密钥字符集校验逐字符限可打印 ASCII、**粘贴进来的 `NAME=value` 环境行会被当格式错误拒掉**（`:78`）；删一行时**只有派生目标完全一致才清凭据**，别的一律保留"因为这枚行无法证明自己拥有它"（`:124`）。我们：单价／窗口／能力位／引用名这些**键全在**（`internal/config/schema.go:352-368,387`），**那一屏没有**。

**D-2〔根本没做〕第一次打开时的引导。**
DeepSeek 有两个有序弹窗（版本化内测声明＋按条件显示的凭据步骤，`ui-settings-models/README.zh.md:12,36`）；账号那枚包还管"欢迎页必须点开始设置／额度页始终展示／按用途问不同的题"（`ui-settings-account/README.zh.md:60`）；并且**已完成引导的安装不会重新应用默认值**（`:97`）。我们这边：`schema.go` 的 19 个 section 里没有引导账本，Go 侧也没有"第一次跑"这一形。

**D-3〔建了但没接〕各家订阅额度那一屏。**
openchamber 在扩展宿主里做了约 20 家编码订阅的用量窗口（日／周／月／本轮／会话五类，带重置倒计时、余额、消费上限）：`openchamber/packages/vscode/src/quotaProviders.ts:15-52`。我们：**每提供商日／月配额两枚键**（`schema.go:367-368`）＋**日／月预算＋提醒阈值＋超预算动作**（`:544-548`）都在，**但没有那一屏**（数据有处可来）。

**D-4〔根本没做〕网页搜索的提供方设置（密钥／接口地址／一次搜几条）。**
DeepSeek 有一整枚设置页（`ui-settings-web-search/README.zh.md:2,68`）。我们：`schema.go` 的 section 名册里**没有这一节**（本轮逐名读过 19 枚）。⚠ 这一格正好撞 AGENTS §2 那条未定案（`web.search` 实现路径）——**我只登记缺这一节，不定实现路径**。

### §4.E 语音（听与说）

**E-1〔仅文档〕语音那一屏：识别模型逐个列出来、带准确度与速度两条评分、每个能单独下载或删除。**
openchamber 整屏在此（`openchamber/packages/ui/src/components/sections/openchamber/VoiceSettings.tsx:42-66,219-220,383`：四个模型、逐个下载删除）。我们：只有"提供方＋模型名"两枚文本键（`internal/config/schema.go:210-218`），**且真装配还没跑**（票 15／26／41 三条 DEFERRED 在册；`internal/audio/doc.go:27` 逐字 `DEFERRED(playback)`）⇒ 这一格是"规格写了、盘上没有"。

**E-2〔同一条路上别人交完的实测账，不是功能差〕**
他们的识别引擎就是 sherpa-onnx，而**我们自己规划用的也是这一家**（`internal/speech/doc.go:1-2` 逐字规划）。所以值得抄的不是"接云端识别"，是**这些数**：那段模型算得随长度平方涨（60 秒 2.1 秒/+90MB、180 秒 9.3 秒/+490MB、300 秒 21.3 秒/+1.5GB）⇒ 60 秒遇静音就交一段、90 秒硬顶；**振幅峰值低于 300 的静音段直接丢**（不让识别器在静音里"幻听"出字）；只确认连续序号、未确认的自己重传（`openchamber/packages/web/server/lib/dictation/DOCUMENTATION.md:46-57,63-98,99-100`）。

**E-3〔根本没做〕转写稿回到"开始录音那一刻的那份草稿"。**
openchamber 专门有一枚钩子管这件事（`useDictationOrigin.ts:1-10`）：**你在等回字的窗口里切走了会话，文字仍塞回原草稿，不会串台**。我们〔根本没做〕。

### §4.F 工具、插件与扩展

**F-1〔建了但没接〕"这台机器装了哪些插件"那一页（只读清单）。**
DeepSeek 有「内置插件」分区＋一枚只读清单标签页：首次选到那个标签才去取（`ui-settings-plugin-inventory/README.zh.md:28`）；显示时**把包名缩短**（去掉 npm scope 与 `dsh-` 等前缀）但详情与搜索保留完整标识（`:34`）；预设切换器**只改变列表显示什么、不写任何设置**，坏预设带标记并把"为什么坏"显示在原位（`:40`）；没有任何标签页贡献的部署显示分区的空提示（`ui-settings-plugins/README.zh.md:28`）。我们：`schema.go:556-573` 有插件开关／能力白名单／宿主接口白名单等键，**清单那一页没有**。⚠ 反面参照：openchamber 的内置扩展登记表**是空的**（`openchamber/extensions/registry.json:1-3`、README `:17`）——**"看着一堆内置其实零内置"**。

**F-2〔根本没做〕技能这一维：连配置键都没有，更谈不上"生效没生效"。**
① 用户点哪里：设置里一列技能，每条带"谁装的／哪一层／优先级／是不是外部内容"，还要能回答 owner 那句"我配了呀，怎么没生效"。
② 别人怎么做：minimax 的快照里同时有**赢的／输的（为什么输、赢的那枚在哪）／诊断（警告还是错）／扫了几枚·读了几枚·复用了几枚／渲染了几枚·压成摘要几枚·丢掉几枚·预算用了多少字**，整表还能 dump 成一份 JSON 供事后核（`minimax-code/packages/agent-modules/skills/src/types.ts:14-22,49-93,100-115`）；六层来源各带优先级与"外部"标记（`:1-13`）。
③ 我们现在是什么状态：`grep -c Skill internal/config/schema.go` ＝ **0**（同文件 `Plugin` ＝ 8，本轮复量）⇒ **技能连写都没地方写**。所以这一格的正确顺序是**先立配置面**（碰 D36 的 section 树，属契约面），再谈清单。

**F-3〔根本没做〕卸载过名字就不许被下一个继承。**
minimax 把"工具名字"当**永久资源**管：两级分配、落盘＋加锁＋原子替换，三条语义逐字 `Never derives identity from a public name.`／`Removed identities retain their assignments so a new source cannot inherit a name.`／`Disk-backed instances refresh inside an exclusive transaction before allocating.`（`packages/agent-modules/mcp/src/runtime/name-registry.ts:50-54`），长度上限具名（整名 80／服务器段 48，`tool-name.ts:2-4`）。我们：本轮复量 `nameRegistry|tool_name.*unique` **0 命中**；`internal/memory` 里的 4 处 `roster` 全是协程名册，不是名字分配器。⇒ 后果：**你卸了一枚插件、又装一枚新的取了同名 ⇒ 它继承旧插件的授权和缓存**。

**F-4〔根本没做这一层／要不要开属 owner 一句话〕插件能不能参与决策（钩子点名册）。**
minimax 有 11 个钩子点（`packages/agent-modules/plugin-hooks/src/contracts.ts:2-14`），并且**把"钩子能改哪些东西"写成了一张清单**：能不能替用户批准（他们能，注释逐字 `an explicit allow authorizes the tool without showing the product UI`，`:156-161`）、能不能改工具入参（`:171`）、能不能改"模型看到的工具结果"（允许但**账上必须留原文**，`:187-190`）、能不能改权限规则（加／换／删规则、**换档**、**扩目录授权**，宿主必须整条原子应用，`:137-153`）、能不能停掉正在跑的东西、能不能写终端控制序列（有白名单函数把关）。他们的执行体是**子进程**，Windows 那份取消作业正是我们要的：`taskkill /pid <pid> /T /F` 杀整棵树、隐藏窗口、**只给 500 毫秒排空**、输入 1MB 硬闸、输出超 64KB 判"输出不合法"并杀进程、普通事件预算 15 秒而**收尾事件只给 3 秒**、六种失败结局每枚都带"哪个插件的哪份文件第几个声明"（`plugin-hooks/src/runner.ts:520-529,30-38,418-429,36-37,193-201,57-63`）。
我们：本轮复量那 11 个钩子名 **0 命中**；`internal/tools/fs_write.go:283` 一带确实用 taskkill，但那是"写文件中途杀掉子进程"，**不是插件事件层**，而且那里的注释承认会留下中转文件。⇒ 开了就有合规预处理（"要发往外网的文本先脱敏"），**同时攻击面变大**——这一层是契约级，我不排。
⚠ 三行定性：
- 现象出现在哪：**我们没有这一层**，所以谈不上被谁绕过；列这一条是为了让 owner 看见"别人把红线画在哪儿"。
- 有没有本机被入侵的证据：**没有**，本腿没做任何入侵类验证，也不把"缺少某层"写成事件。
- 最坏后果是什么形状：若将来开了这一层又照抄他们的输出面（钩子可替用户批准／可换档／可扩目录），**最坏是插件说明书里一句话把门控调宽**——那正是我们已有的那条禁形要挡的形状。

**F-5〔别人也没有／只登记〕各家加号里到底放什么，五家没有一个把 MCP 放进去。**
本轮点到行的五家：MCP 在 DeepSeek 走设置页；在 Step-Code 是 `/mcp` **只读回显**（`Step-Code/packages/coding-agent/src/step/slash-commands.ts:45-55`，描述逐字 `Show configured MCP servers and loaded tools`）；Pi／minimax／openchamber 本轮**未读到任何 `/mcp`**。⇒ 若要做，只该做"菜单里给它留个位置"，**不做"接 MCP"**（那是 D13 取舍，人工批准级）。

### §4.G 文件改动、差异与回退

**G-1〔界面地界·不读／Go 侧本轮未量到〕"改动文件卡片＋逐文件对比＋收尾正文里可点击的行内文件引用"。**
DeepSeek 有一整枚包专门管这个（`ui-deliverables`，名册第 8 枚，说明逐字 `改动文件卡片＋逐文件对比 review tab＋交付文件卡片＋收尾正文里可点击的行内文件引用`，`survey-2026-09-28-dsh-ui-packages.md:42`）。⚠ 我只读到它的**说明首行**，没读实现。

**G-2〔仅文档〕"只把这一块改动应用上"。**
openchamber 的宿主代理面里 git 一共 24 种操作，含 `apply-hunk`（逐块应用）与"生成 PR 描述"（`openchamber/packages/vscode/src/bridge-git-*.ts` 一族，逐 case 行号见 `survey-2026-09-28-oc-mobile-vscode-extensions.md:38`）。我们：`internal/panel/git.go:82` 那句实话逐字写着"切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地"〔已证在产码里〕——**连换分支都还没通，谈不到逐块应用**。

### §4.H 终端、命令与进程

**H-1〔根本没做〕面板里那一格真终端（能打开、恢复、插话）。**
DeepSeek 有 `ui-sidebar-terminal` 一整包（"在网页右侧栏打开、恢复和控制交互式 shell 标签页"）。⚠ 这枚包**恰好不在他们官方包地图里**（本轮复核：那张表 50 行、磁盘 52 枚，缺席的两枚是 `ui-settings-account` 与 `ui-sidebar-terminal`，尺 `grep -oE '\[`ui-[a-z-]+/`\]' | sort -u | wc -l`＝50＋两枚 `grep -c`＝0）。我们：刷给界面那包内容仍是**固定 6 个键**（本轮读 `Snapshot` 结构），**没有终端那一块**。⇒ 后果：它跑了一个需要你插话的命令（装依赖问 Y/N、SSH 要口令），你今天**只能看输出不能插手**。

**H-2〔根本没做，但有整套现成算法〕双击打开的程序拿不到你在终端里有的东西。**
① 用户看到的形状：装了某个命令行工具，Wisp 说找不到；你打开终端一敲就有。
② 别人怎么做：桌面程序启动时**主动向系统要一份当前环境清单**，四级回退（先叫一个 PowerShell 问 → 退到另一个 PowerShell → 退到系统目录里的全路径 → 最后用 `cmd` 问一遍），结果用**不可打印字符当分隔符**（防路径里的空格把结果切碎）、隐藏窗口、10 秒超时，还找不到再 `where` 问一句"你到底装在哪"（`openchamber/packages/vscode/src/opencode.ts:529-575`，找不到再 `:447`）。
③ 我们现在是什么状态：本轮在 `internal/**` 搜这一族**没有对位物**。⇒ 我们规划里的模型下载与插件定位都要在"服务／计划任务"这种没有终端环境的语境里找可执行文件，**一定会撞上同一枚坑**。
⚠ 三行定性：**这不是安全问题**——现象是"找不到程序"，本机没有任何被入侵证据，最坏后果是**功能不可用＋工单池里一堆查不出根的幽灵 bug**。

**H-3〔根本没做〕两击式停止，以及"停不下来"要能说出来。**
minimax 把"正在停"做成**正式状态**（`Stopping response`，`activity-line.ts:259-260`）；DeepSeek 是两击控件（首击上膛、3 秒内确认击才真杀，`ui-jobs/README.zh.md:30`）；minimax 后台任务的状态机七枚里 `stopping` 与 `lost` 都是正式状态，**且停止失败照样落 `canceled`、只在事件里带一枚 `stopFailed: true`**（`packages/agent-modules/background-task/src/manager.ts:120-147`）。我们：本轮复量两击与"取消失败"表达**0 生产落点**。⚠ 抄的时候有一条陷阱：**抄状态、不抄那枚旗标＝我们自己造一枚"账上已取消、实际进程还在"的假终态**。

**H-4〔根本没做〕"跑完了但还没告诉会话"是一列能查的账；重启后孤儿要有"丢了"这个状态。**
minimax 任务对象上带 `deliveredAt`、查询可只取"跑完了但还没送达"的那些（`agent-modules/background-task/src/types.ts:82,108-110`），落库时**真有一列并为"按会话＋按状态＋按有没有送达"建了复合索引**（`local-runtime/src/background-task/schema.ts:38-45`）；重启后发现宿主进程不在了 ⇒ 落 `lost`＋错误码＋发完成事件（`local-runtime/src/background-task/startup-recovery.ts:88-106`），并钉死"只有那一轮确实没活下来的精确证据才允许结算，**年龄或会话空闲都不算证据**"（`:147-149`）。我们：本轮复量 `deliveredAt|delivered_at|TaskLost` **0 命中**；名册仍是**进程内一张表**（`internal/tools/task.go:165` 注释逐字 `process-local`）。⇒ 后果：崩溃或重启后，后台任务与子代理在账上**要么消失、要么停在"还在跑"的假状态**，你会盯着一个已经死了几小时的"进行中"。

### §4.I 用量、成本与配额

**I-1〔根本没做〕每一轮"实际喂给模型什么配置"要留一份证据。**
minimax 有一枚专门的记录：**不含消息正文、不含密钥**，只存重建这一轮需要哪些入参（系统提示、每枚工具的名字＋说明＋参数形状、模型、服务商、接口、思考档位、输出上限、输入打包上限、缓存保留策略…；`agent-modules/session-report/src/contracts.ts:52-67`），存法是规范化定长文本＋短码去重、单枚 4MB 上限、**原子替换不留半截**、**压缩完成后把这份冻结在它描述的那一轮旁边**、**反悔之后删的是上一轮的证据不是运行配置**（`llm-call-report-store.ts:47-53,14-15,52,57-75,77-101`）。我们：本轮复量 `thinkingLevel|cacheRetention|maxSerializedInputBytes` 命中的 3 处**全是"模型能支持哪些档位"的清单，不是"这一轮实际用了哪一档"的读数**（`internal/config/schema.go:357-359`＋`internal/config/validate.go:201`）。⇒ 后果：出问题时你要回答"当时到底喂给模型什么"，**今天只能靠读代码反推**。

**I-2〔做了但会咬人〕中文的字数换算，我们今天算少了。**
① 用户看到的形状：长中文会话聊到一半，服务商一口把请求打回，报错看着像"窗口不够"。
② 别人怎么做：minimax 把这件事写进源码注释并点名后果——按"每 4 字符算 1 个"会**低估 4-8 倍**（汉字实际每个吃 1-2 个模型字），低估会把"该压缩了"那条线推到右边，**推到兼容接口的服务商直接拒绝**（`agent-modules/context-manager/src/token-estimator.ts:6-11`，逐字点名某型号"输入超过窗口减输出上限就拒绝"）。他们的做法：**往回找最后一条成功且带真用量的消息、信它上报的总数当前缀和、只估它之后的**（`:20-27,137,165`），默认估计算法**故意偏高一点**（注释逐字"宁可早触发"，`:31-36`），**每条计数把来源写进结果里**（服务商给的／自己估的，`types.ts:28-31`），每条消息加 4 字包装开销、图像块按 4800 字计费（`:38-45`）。
③ 我们现在是什么状态：字节数除以 4（`internal/agent/budgets.go:125-129`），逐条累加、**不锚定上一条真值**（`internal/agent/compress.go:110-116`）⇒ 中文每字 3 字节算出约 0.75 个模型字、实际 1-2 个，**低估约 1.3-2.7 倍（这本算术是我本人算的，方向与他们注释同向）**。⚠ 每一条预算阈值都要过这个函数，而我们是中文优先产品。

**I-3〔根本没做〕压缩这件事要有四个读数，而且压之前可以先推掉这一次。**
minimax：四枚观测事件各自带"阶段"（好归因是哪一段崩的）、主动压缩线**永远早于**服务商硬上限线、按窗口比例封顶、输出预算随上下文变胖自动缩小但有地板、**压缩前的检查点回调返回"不"就这次先不压**（`agent-modules/context-manager/src/provider-budget.ts:26-41`、`settings.ts:13-36`、`types.ts:94-127,146-149`）；压缩后的摘要里专门存一个条数，注释逐字"这些保留下来的旧消息用量已经过期，**把它和摘要一起存下来，恢复历史时才有同一条界线**"（`types.ts:9-16`）。我们：只有"按配置窗口等比缩放"（`internal/agent/budgets.go:7-11` 注释逐字 D15 要求）〔已证＝这段注释与这些常数在生产代码里〕，**四枚读数、两条线、压缩前闸门、那条"用量界线"都没有**。⇒ 后果：**两种算法打架时没人能解释为什么这一轮压了／没压**；重启后它会拿一份过期用量继续算钱。

### §4.J 状态、通知与提醒

**J-1〔本版新查／根本没做〕20 枚状态名里没有"正在压缩"这一名。**
别家有这一相位（minimax `Compacting context`，`packages/tui/src/tui/shell/activity-line.ts:256-257`；另有 `Compacting now`，`transcript/context-visualization.ts:211`）。我们：`internal/statemachine/states.go:11-31` 那 20 枚逐名读过一遍，**没有压缩这一名**⇒ 球与面板**说不出"这一刻它在整理历史"**，用户看到的仍是 `Thinking`／`Acting`。⚠ 这一格是**本版相对第三版新查出来的**，而且它是 §4.I 那一串压缩缺口的**可见面**。

**J-2〔根本没做〕"现在做不了，但你先做这一步就行"那一族句子。**
minimax 一整族，每句都带下一步：`停止当前回复后才能修改会话历史。`（`copy.zh-Hans.ts:23`）、`请先开始或恢复一个会话，再浏览历史。`（`:33`）、`请退出计划模式后再回退此会话。`（`:113`）、`请打开父会话后再执行回退。`（`:114`），还有一族**连怎么退出当前状态一起给**（`catalog.ts:760-766`）。我们：只有通道级的"XX 取消不可用"四句（`internal/agent/approval/approval.go:106-114`）〔已证〕，**没有"你该先做哪一步"那一半**。

**J-3〔根本没做〕等待态写成两行：一句"在等什么"、一句"它自己什么时候好"。**
openchamber：`Waiting for the dev server` ＋ `It is not accepting connections yet. This page will load as soon as it does.`（`packages/ui/src/lib/i18n/messages/en.ts:1345-1346`）；另有 `Waiting for reviewer`／`Waiting for implementer`（`:1698-1699`）。我们：**没有这种两行式文案的先例**（最接近的一句反而写得更诚实、但全仓只有这一处：`internal/panel/git.go:82`）。

**J-4〔根本没做〕"提醒"该是一条链＋一份配置，不是一段写死在环路里的代码。**
minimax：30 枚提醒块挂在一条链上，可按名字开、按框架开、还有一档"要紧的"（关掉全量提醒时**要紧的照跑，但要紧的也走同一份白名单、不是特权通行**）（`agent-modules/system-reminder/src/registry.ts:25-73`）；**每个名字有自己的频控**（一份区间表＋冷却＋间隔轮数），语义钉死"没配＝全跑（老行为）／配了空表＝全关（安全兜底）／配了名字＝只跑这些"（`types.ts:378-395,409-422`）；记忆提醒是**指数退避**（第 10 轮第一次提、不理就 20、再 40 封顶、照做就保持节奏，那份按会话的状态**还带容量上限**）（`evolution.ts:21-29,44-52`）；同一份内容备了**浓、淡两版**（`blocks.ts:89`／`:206`）；**密钥那一枚只报变量名、绝不报值**，三种语义分清（没接线就闭嘴／库是空的也闭嘴／非空才渲染，且只在第一轮或名字集合变了时渲染）（`types.ts:360-377`）。我们：只有一种提醒（重复调用那一条）直接注入环路，**按名字频控／冷却／退避的表 0 命中**（本轮复量）。

**J-5〔建了但没接〕"有几件事在等你"这个数，载体与数据源都现成。**
openchamber 的应用图标角标＝自上次打开以来**不同的提醒条数**（不是会话数），服务端算、三路信号清（可见性心跳＋打开会话＋发消息都算"人在用"）（`APNS.md:40-66`、`AppDelegate.swift:201-208`）。我们：**等批准队列的长度已经是一个现成的数**，球是天然载体，**但没有"数上去"那一跳**。
⚠ 同时抄四条端无关规则（第三版一(7) 已列，本版继续）：① 别拿"问客户端你还在不在"当"要不要吵你"的门（他们用注释逐字写着放弃了那道有竞态的门）；② 断线重连后**先对齐权威态**（`sessionActivityWatcher.ts:34-52` 逐字"凡不在册的一律 idle，包括本进程断流前以为 busy 的"）；③ 冷启动那份"要点开哪一条"的意图先存住（`deepLinkNavigation.ts:15-16`）；④ **子代理还在跑时父任务的完成通知不许发**（`notifications/DOCUMENTATION.md:51-53`）。

### §4.K 设置页与偏好

（这一整块的逐枚对照在 §2 表二，这里只留三条最要紧的叙述。）

**K-1〔根本没做〕三件配套：「已覆盖」标签＋「恢复默认」＋"读不懂不许静默"。**
DeepSeek 每一枚设置字段都显示"生效值＝用户覆盖叠在默认之上"、被覆盖的带**「已覆盖」标签**、旁边一枚**「恢复默认**」、**点保存之前不写入任何东西**、清空并保存等于重置、填了非数字**阻止保存并在字段下说明原因**（`ui-settings-agent-loop/README.zh.md:28`、`ui-settings-shell/README.zh.md:28`）。键位那面更狠：**读失败就禁用编辑并说明当前在用哪套**、损坏时**先备份再修**、被更新版本写的就提示升级、**"恢复全部默认"不覆盖读不懂的数据**（`ui-shortcuts/README.zh.md:36`）。我们 173 枚键**全靠手改文本文件**，这三件一件都没有。

**K-2〔编译期常量／碰配置表要 owner 一句话〕两个上限做不成设置项。**
DeepSeek 把"能派多深／能同时活几枚"做成一页两半一次保存、两次写入互相独立、过期草稿**报冲突要求放弃而不覆盖较新的那份**、字段带最小值与安全整数校验，**并把计口径写在数字旁边**（"容量统计同一主脑下所有层级还活着的子代理，主脑自己不计入；深度让位于工具自己的上限"）（`ui-settings-subagent/README.zh.md:12,28,32`＋`src/client/subagent-limits-card-controller.ts:10-13,36,48`）。我们仍是 `internal/tools/subagent_197.go:73,79` 两枚常量（本轮复量），而那个 4 是从**冻结的并发天花板**借来的 ⇒ 把它做成设置项会开一条"从名册侧把上限调到桥之上"的路，**属契约级**。

**K-3〔没有〕正文字号。**
DeepSeek 的主题页管"主题**与正文字号**"（名册第 47 枚，`survey-2026-09-28-dsh-ui-packages.md:81`）。我们：有球尺寸、面板宽高与缩放（`internal/config/schema.go:165,532-538`）〔已证＝键在〕，**没有任何字号键**——对看不清字的人是硬缺口。

### §4.L 外观、几何与可读性

**L-1〔必须自己拍／抄不来〕谁让位、按什么顺序让位。**
DeepSeek 把所有数值摆在明处（中栏最小 400、左栏 264-420、默认 280、收起 56、窗口窄到 1024 自动收起、右栏最小 300、默认占 45%、上限 70%，`ui-layout/src/client/columns.ts:11-29` 八枚具名常量），**挤压次序逐字写着**："为了给中栏留 400，先把右侧面板缩到 300，再报告空间不足让占用方自己关闭，最后才压缩中栏"；变宽**不自动重新展开**；每一行 chrome **自己声明自己是可拖区**，于是这行的空白处能拖、行里的控件照样能点（`ui-layout/README.zh.md:30,36,38,54`）。我们：Go 侧不产这类数；⚠ **我们有面板＋球两块面、比他们多一层，而他们只写了两级次序、没写"贴边面板 vs 悬浮球"**——这一格是必须我们自己拍的空格。

**L-2〔给界面那支的话，不是我们的活〕状态图标：一枚共享组件、三态一名、每行留同宽。**
DeepSeek 全程用具名矢量图标、没有一枚字符图标；状态点是共享组件、只有 进行中／完成／静止 三态；每行为状态图标预留同一 14 像素列宽（`SubagentHeaderLineage.tsx:167-188,334`＋`ui-subagent/README.zh.md:36`）。我们球那侧 20 种状态各有**颜色**字（`states.go`＋球渲染），**"同一形状、同一位置、不只靠颜色"这一维值得点名给他们**，但不由我拍。

**L-3〔抄法警告〕抄他们那"每秒一跳"的计时之前先看这条。**
他们的耗时是"已结束累计＋（现在 − 本轮起点）"，"现在"由 1 秒定时器推（`SubagentHeaderLineage.tsx:205-211,77-91`）——**这是墙钟差**，而我们禁的是"用墙钟差实现超时"（`AGENTS.md` §1.2）。可用的替代是他们自己的分层：**存储侧只放事件起点／终点／累计，那一秒一跳只在显示那一层**（`projection-types.ts:22-37`）。⇒ 写手照原形抄会白跑一整腿。

### §4.M 斜杠命令与快捷指令

**M-1〔分母先摆出来〕各家到底有多少枚命令（每行数都带尺）。**
Step-Code **38 枚＝22 内置＋16 产品**（内置尺：`sed -n '19,75p' packages/coding-agent/src/core/slash-commands.ts | grep -c '"\?name"\?:`＝22；16 枚那半是〔腿报〕）；pi-upstream **24 枚内置**（同文件同尺、数组更短）；minimax **约 45-52 枚，取决于把"拆开写在别处的"和"由描述表自动生成的"算不算**（两把尺分别给 48 与 45、件里报 41/52 ⇒ **引用时不许单引一个数**）；openchamber **15＝9 提示词对＋6 动作**（9 那组我复现过、6 那组〔腿报〕）；DeepSeek **至少 6 枚后端命令＋第 7 枚 `/file` 没有后端**（名册由插件运行时注册、静态数不出来）。**`D:\work\AI\open source\pi` 那 37 枚文件不是编码助手**（README 首行是 GPU Pod 管理器）⇒ 本版所有 Pi 系结论的对象都是 `pi-upstream`。

**M-2〔根本没做〕名册、判"是不是命令"的代码、以及"能不能用"的三样。**
① 同一枚解析器**同时**供"界面高亮"与"回车执行"（minimax 设计声明逐字 `Exact-token commands reject non-whitespace suffixes, so callers cannot visually classify input differently from the command path that will handle Enter.`，`packages/tui/src/tui/commands/catalog.ts:692-717`）；openchamber 的反面历史：`@` 规则被写过四遍、`/` 三遍，于是"涂成引用却解析不成引用"（`language/prefixTokens.ts:1-15`）。
② 没声明收参数的命令**不接受尾巴**（`catalog.ts:714`）。
③ 每条命令自带"现在能不能用"＋"不能用时说哪句话"（同文件 `visibleWhen` 18 枚、`unavailableReason` 19 枚、三态 `unrecognized/unavailable/handled` `:627-663`）；"菜单里看不见"与"手打会被拒"是**两套代码各有各的说法**，不会一边放行一边拒收（`:738-747`＋`:640-656`）。
④ 名册里每一项**必带出处**（Pi 系描述前拼 `[来源]`，`interactive-mode.ts:660-666`）。
我们：〔根本没做〕全部四条，且**界面那支与 Go 那支是两拨人** ⇒ ①那一形是我们最容易自己造出来的 bug。

**M-3〔已证的错处／票 221 在飞〕我们的名册说明文字许诺了一枚没注册的能力。**
派子代理那枚工具的说明逐字写着"可以用 `task.cancel` 单独停它"（`internal/tools/subagent_197.go:196`），而 `task.cancel` 在名册头上逐字 **DEFERRED**（`internal/tools/task.go:23`，本轮复量）。别家用**机器闸门**挡住同一枚病：DeepSeek 靠"能力不在那一行就不存在"（`packages/plan/plan-mode/src/index.ts:230-231`），minimax 靠 `audience:'internal'`＋`visibleWhen` 双重（`catalog.ts:47-84`）。⇒ 判据该写成"装配时没有真执行者的命令，注册这一步就该失败，而不是注册了再置灰"（转引 `survey-2026-09-29-command-catalog-across-harnesses.md:47`）。

**M-4〔取向不同／别家也不一致〕没打中命令时到底说不说。**
DeepSeek：报错，**绝不静默降级成普通提示**（`ui-commands/README.md` 与 `ui-commands/src/client/service.ts:391-393`）；minimax：拒＋人话理由＋连退出路径一起给；Pi 系：命令级不拦、参数级报错带可用值（`interactive-mode.ts:141,180`）；openchamber：**不提示，那一行保持普通文本**（`prefixTokens.ts:11-14` 逐字 `Membership is the authority; the pattern is only a locator.`）。⇒ **"拒了要说为什么"是我们自己的选择，不是业界共识**，写票时不许引"别家都这样"。

### §4.N 其它（不好归块的）

**N-1〔根本没做〕"这一条是谁发起的"必须是一等业务数据。**
minimax 有一份 13 枚来源的词表（网页接口／定时／任务／后台任务／团队／线程目标／问卷／通讯／代码评审／问候／三个聊天软件频道），它进了**每一条投递的入参**（`agent-modules/conversation-contract/src/index.ts:1-14,357`）；会话"当前什么状态"（5 值）与"这是哪一类会话"（6 值）**分成正交两轴**（`:16-29`）。我们：最接近的东西管的是**内容**从哪来（`internal/risk/provenance.go`），**不是触发**从哪来（本轮复量 `ConversationSource|MessageSource` 0 命中）。⇒ 后果：**定时唤起的那一次和你在球上按的那一次，在它眼里一模一样**；出了事分不清是哪条路进来的。

**N-2〔仅文档／登记先例，不替你拍〕挂起的"等人回答"该续上还是作废。**
我们待定案清单里明写着"等批准时系统挂起（倾向作废判拒绝）"。minimax 的做法是**给一条独立的窄入口**，注释逐字 `Resumes a persisted user-input wait. This is a narrow internal admission seam: it never joins the ordinary Queue, and requestId is its durable replay identity.`（`conversation-contract/src/index.ts:368-387`）——两条路都满足"不复活陈旧的批准请求"，但他们多给了"挂起的用户输入可以被合法续上"这一格。**我们的定案不改，这只是一条外部先例；碰它属"没定义就停"。**

**N-3〔根本没做，而且几乎不要钱〕目标文本用摘要钉住，"账过期了"分三种原因。**
minimax：目标语句用 sha256 钉住（`agent-modules/goal/src/objective-digest.ts:4-5`），过期原因是三枚分开的名字——目标没了／纪元换了／目标文本改了（`store-port.ts:89`），**不共用一个"过期了"**。我们：本轮未量到同形物；后果是"目标被改过一次"和"目标被删了"混成同一条，回来续跑时分不清该接着算还是该重算。

**N-4〔取向不同＋缺一格〕定时交给系统这件事，我们的方案里少讨论了一维。**
minimax 有应用内调度器（我们不接，D12 是**取向**不是缺口），但有一枚我们完全没有的形状：恢复时可以选择**"隔离恢复"**——从别的机器或别的副本恢复来的那份档期表，**先只登记供查看和供用户修改、不执行**（`agent-modules/cron/src/registry.ts:42-48,114`，注释逐字 `A quarantined clone still registers tasks for inspection and later user mutations`）。⇒ 一旦我们把定时交给系统，"从备份恢复来的定时任务立刻开跑"就是必然发生的事，而**我们的方案里连这一维都没出现过**。
⚠ 别把他们写成"已解决"：他们**错过的档期不补跑**（本轮件里那条：搜索"逾期／补偿／错过的／漏发"0 命中，重启登记时"上次跑是什么时候"一律写空）。

**N-5〔取向不同／不接〕把 git 工作树当成"会话的一种运行位置"。**
minimax 的词表允许"当前／新工作树／已有工作树"作会话去处，还带分支与母仓路径（`conversation-contract/src/index.ts:79-85`）。**这和我们"绝不在仓库目录内建工作树或检出"直接冲突**——要么明确写"不接"，要么单独议，不许在实现里顺手接进来。

**N-6〔别人也没有／只登记〕"每日汇总／今天它替我做了哪些决定"这一屏。**
本轮逐家点到的只有 minimax：那枚最容易误会的模块**不是日报**，是"用户点提交反馈时把一棵会话目录里的文件打包成清单"（`agent-modules/session-report/package.json:5` 逐字 `Session diagnostics and feedback artifact collection.`）；仓里那枚"每日摘要"设置写的是**记忆文件**、默认关、而且被它引用的实现文件**在这份副本里根本不存在**。⇒ 第三版那一节维持；**难点不在画那张表，在有没有"每一轮结束时的一枚结构化小结"**（他们至少三份现成原料：每轮防跑飞信号的次数与首末位置、每次判定的双语理由、后台任务的完成与送达记录）。**其余几家本轮未查，不许写"六家都没有"。**

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

> 规矩：**读不到就写"读不到＋为什么"**，并且把顺序交代清楚。下面每一行都是"本轮没读／没量"的东西，不是推测。

**甲·参照仓侧（按本轮的读序，越靠前＝越先想做但没做完）**

1. **派单第 2 优先级点名要的"Step-Code 与 pi／pi-upstream 的设置页与状态文案"——这是本轮最大的未完成项。**
   本轮实际做到的是：pi-upstream 抽了 6 段人话串（`interactive-mode.ts:455,1181,141,180,2680,3649`）、Step-Code 抽了 `plugins.ts`／`slash-commands.ts`／`features/step.ts` 里被引的那几行；
   **没做到的是**：两家的**设置项名册逐枚**（Pi 系 `/settings` 背后那台设置管理器有哪几页几行、Step-Code 有没有设置页、页上逐枚叫什么），本轮**没读** ⇒ 所以 §2 表二里这两家的列大部分是空格，**空格＝没读到，不等于"它没有"**。
2. **DeepSeek 那 52 枚界面分包：本轮新读了 9 枚设置类包的说明件，但 `src/**` 一枚未读**；其余 43 枚仍只有说明首行。
   `ui-*` 之外那 14 条目（`connection`／`store`／`modules`／`web`／`locale`／`file-upload`／`shortcuts` 服务侧／`hmr`／`resources`／`AGENTS.md`／两份 README／`README.i18n.yaml`／`tsdown.client.ts`）**只列名未读**。
   本轮点名进表一却未读实现的：`ui-deliverables`（改动文件卡片那枚包）、`ui-attachment`、`ui-goal`、`ui-plan`、`ui-reference`、`ui-schedule`、`ui-open-in-app`、`ui-workspace`、`ui-agent-preset`、`ui-message-feedback`、`ui-sidebar-*` 五枚、`ui-theme`。
3. **minimax `agent-modules` 之外的两枚巨包一枚未读**：`local-runtime-v2`（913 枚 `.ts`）与 `local-runtime`（629 枚）。
   继承第三版的未读清单也仍然未读：`permission` 那 18 枚（其中最该补的是直接对 Windows 的 `windows-native-delete.ts` 162 行与 `written-files-registry.ts` 105 行）、`context-manager` 的压缩主逻辑 528 行、`plugin-hooks` 的 coordinator 1,289 行与 runner 其余约 2,000 行、`system-reminder` 的 `blocks.ts` 907 行正文与 `providers.ts` 854 行正文。
4. **openchamber `packages/ui/src`（本轮量到 3,187 枚文件／568 `.tsx`／1,418 `.ts`）仍是最大的一块空白**：
   本轮只定向抽了 `lib/i18n/messages/en.ts` 里被引的 8 句与 `components/sections` 的**目录名册（21 项）**；
   **未读**：批准卡与那一页的正文（`PermissionCard.tsx` 只读到被引段落）、`ChatView`、设置屏 21 项各自的正文（含语音那一屏只读到抽出来的三段行号）、`web/server/lib/guests/*.js` 实现体、`packages/sdk` 正文、`packages/electron`、`web/server` 其余全部；
   继承腿 1 的未读巨件：`gitService.ts`(4,073)、`opencodeConfig.ts`(3,192)、`SessionEditorPanelProvider.ts`(708)、`webviewHtml.ts`(410)、`webview/api/` 21 枚；
   带移动字样的 62 枚件本轮**一枚未新读**（只把分母复量了一遍，见 §5-4）。
   ⚠ openchamber 无 `.git` ⇒ 它那一切结论只能当"当下快照"，**不许写成"他们一直这么做"**。
5. **六家的 `.git` 历史本轮一次没用**（纪律只禁执行二进制与写被研仓，`git log` 于他家是允许的）。⇒ 下一步若要回答"这条规矩他们是从哪一版开始有的／是不是回退过"，**这五家都能查**。
6. `D:\work\AI\open source\pi` 已复核仍**不是编码 harness**（37 枚文件／0 枚 `.ts`），本轮不引它；`minimax-code/third_party/pi-mono/` 那一份 Pi 副本**本轮未读**，任何"minimax 有 X"的结论都要先确认不是从这份副本读来的。

**乙·我方侧（本轮没量／量不到）**

1. **提醒文案会不会被模型写进记忆**（别家给每套提醒都挂了"别把这条存进记忆"的防污染后缀）——**这一维两轮都没量**（第三版 4.7 标"没查"，本版维持"没量"）。
2. **"每一轮结构化摘要"到底有没有同类物**——本轮只复量了两枚具体版本号名字为 0，**没做全仓反向取证** ⇒ 引用请按弱版本："缺的是可归因的策略版本号＋逐信号计数"，**不是**"我们一枚摘要都没有"。
3. **界面那一面一律未读**：`frontend/**` 与 `design/**` 零读零引零转述 ⇒ 表一/表二/表三里凡"那一屏有没有""拖拽几何由谁声明""状态点画不画"都标了〔界面地界·不读〕。这些格**不是缺口，是本腿射程外**。
4. **各家服务端字段表读不到**（openchamber 无历史＋快照里缺依赖）⇒ 凡涉及"服务端到底回什么字段"的结论都只能算客户端反推。
5. **票面本轮没逐枚读**：§6 与 B 份文件里"撞不撞在飞的票"是按**工单池文件名**（`ls .scratch/wisp/issues/`，其中未带 `-done` 的 2xx 共 18 枚：200/201/211/212/213/214/215/216/219/220/221/222/223/224/225/226/227/228）＋代码注释里点名的票号判定的，**没有读票面正文**。⇒ 编排者派单前须自己核票面落点，我不替他排程。
6. 本轮**没跑任何编译、测试、门禁、二进制**（机器上有别的腿在跑门），所以所有我方读数都是文本级的（`grep`／`ls`／`awk`），**"能跑"这件事本版一条都没声明**。

---

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
