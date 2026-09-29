# 外部同类 harness 的"命令层"对照（2026-09-29 第四轮）

> ⚠ **本件由编排者代落**（09-29 11:2x）。原只读腿 `213-c2` 的结论整段交回但**没写出文件**——我把它派成了**没有写工具的代理类型**（同一天第二次犯，模板缺陷已记 `A419` 与用户记忆）。
> **凡编排者抽验过的行标〔抽验〕并给出尺（命令与口径）；未复核的标〔腿报〕，不许当判据引用。**
> ⚠ **枚数一律带口径**：`grep -c 'name:'` 会把嵌套对象一起数进去（**过计**），`^\s*\{\s*name:` 会把多行写法的条目漏掉（**漏计**）——下面每行数都写清用的是哪把尺。
> 被研仓＝`D:\work\AI\open source` 下各家**只读**；本件落在**被研仓之外**。`openchamber`＝**tarball 快照、无 `.git`** ⇒ 它的一切结论只能带 `file:line`，**读不到历史**。

## ① 各家命令数（分母与尺）

| 家 | 命令数 | 尺与口径 |
|---|---|---|
| **deepseek-harness** | **6 枚后端 `/命令`**（＋第 7 枚 `/file` **无后端**，只是"+"里的一个动作） | 〔抽验〕注册表在 `packages/interaction/commands/src/index.ts`（`:285` 一带是 `register(definition)`，名册由插件注册，`sed -n '283,288p'` 见到的是注册函数本体而非清单——**枚数＝〔腿报〕**）；`/file` 无后端＝〔腿报〕`client/ui-conversation/src/client/apply.ts:272` |
| **Step-Code** | **38＝22 builtin ＋ 16 product** | 〔抽验〕**22** 用 `sed -n '19,75p' packages/coding-agent/src/core/slash-commands.ts \| grep -c '"\?name"\?:'` **＝22**（数组起于 `BUILTIN_SLASH_COMMANDS`，`:19`）；16 枚 product＝〔腿报〕（分散在 `step/slash-commands.ts`、`step/plugins.ts:1263`、`features/step.ts:293,346` 等 9 个文件） |
| **pi-upstream** | **24 枚 builtin** | 〔抽验〕`sed -n '19,44p' packages/coding-agent/src/core/slash-commands.ts \| grep -c '"\?name"\?:'` **＝24**（同文件同尺，只是数组更短） |
| **minimax-code** | **52**（41 内联＋11 spread），运行时还把 skills 追加成命令 | ⚠ **口径敏感**：〔抽验〕`grep -cE '^\s*(name\|command):' packages/tui/src/tui/commands/catalog.ts` **＝48**、`grep -c 'name:'` ＝45——**腿报的 41/52 与我的两把尺都不同**（写法混、跨两枚文件）。⇒ **引用时写"约 45-52 枚，取决于把 spread 与 descriptors 算不算"**，不许单引一个数。spread 源＝〔腿报〕`application/command-descriptors.ts:7-48` |
| **openchamber** | **15＝9 prompt 对＋6 动作**（其余命令来自 OpenCode 服务端） | 〔抽验〕**9** ＝`sed -n '61,129p;145,152p' packages/ui/src/components/chat/composer/submit/slashCommands.ts \| grep -cE '^\s*(\{)?\s*(name\|id\|command):'`（**只复现了 9 那组**）；6 枚动作那组**我的尺没数出来**（形状不是 `name:`）⇒〔腿报〕；服务端来源＝〔腿报〕`stores/useCommandsStore.ts:25,45` |
| **`open source\pi`** | **0——而且那枚目录不是编码 harness** | 〔抽验〕顶层只有 `pi.js`／`vllm_manager.py`／`pod_setup.sh`，`README.md` 首行逐字「**GPU Pod Manager**」（部署 LLM 到 GPU pod 用）。⚠ **这不是新发现的错**：台账早已记着"按 `badlogic/pi` 下错、落盘只有 9 枚文件／685 KB，**真 Pi 在 `earendil-works/pi`（原 `badlogic/pi-mono`）⇒ 已重下成 `pi-upstream/`**，错的那 9 枚按'只建不删'原样留着"（见 `pending-and-issues.md` 那一节的"一次我自己下错版本"）。**本轮所有 Pi 系结论一律以 `pi-upstream/` 为对象**（09-28 那份加号菜单调研也是引它），**没有因为这两枚同名目录而串味**。 |

## ② 最要紧的五行：**菜单项背后是实现，还是一句提示词**

| 家／命令 | 背后是什么 | 出处 |
|---|---|---|
| openchamber 的 **9 枚**（summary／plan-feature／craft-goal／schedule-task／catch-up／debug／weigh／explore／workspace-review） | **就是两串提示词**（可见文案＋instructions magic prompt），**零实现** | 〔腿报〕`slashCommands.ts:22-40` |
| minimax `/review` | 常量 `'Please review my uncommitted changes.'` | 〔抽验〕**这句我在两处独立命中**：`minimax-code/packages/tui/src/tui/controller/product/command-flow.ts:142`（`const TUI_REVIEW_PROMPT`）＋`packages/tui/src/headless/invocation.ts:234` |
| Step-Code `/init` | `pi.sendUserMessage(STEP_INIT_PROMPT)`——发提示词，不是干活 | 〔腿报〕`features/step.ts:304` |
| Step-Code `/mcp` | **只把状态回显给用户**（`formatStepMcpStatuses()`），不改任何东西 | 〔腿报〕`step/slash-commands.ts:47-52` |
| deepseek-harness `/permission` | **真写路径**：解析预设 ⇒ 同时写 sandbox＋approval 两枚旋钮 ⇒ 落一条 session 事件 | 〔腿报〕`interaction/permission-presets/src/index.ts:258`、`:264-274`、`:403-418` |

⇒ **给我们的那条判据**：别家**大量**菜单项是提示词包装。这不是"别家偷工"，是**两种东西要分开画**——**会改系统状态的**（DSH `/permission` 那种）与**只是把一段话递给模型**（openchamber 那 9 枚）。我们的名册若不分这一栏，用户点错了预期就要我们背。

## ③ 六枚重点命令：各家做法差异（票面可直接引）

| 重点 | 差异 | 出处 |
|---|---|---|
| **上下文压缩** | DSH：**遮蔽、不删除**；`/compact` **零参数**；**只人触发**〔腿报 `compaction/src/index.ts:143`、`types.ts:115-119`〕。pi／Step-Code：**收自定义指令**〔`interactive-mode.ts:6858`／`input-dispatch.ts:376`〕。openchamber：走服务端 `compactSession`，另有**"pin 出压缩区"**（被钉住的内容不进压缩区）〔`client.ts:1256`、i18n `en.ts:2364`〕 | ⚠ 对我们的含义：压缩必须能回答"**这次被藏掉了什么、还能不能找回来**"——"遮蔽不删除＋可 pin"是两家独立选同一形状，值得抄 |
| **goal／plan 这类"模式"** | DSH `goal`＝**域对象**（phase／revision）〔`command-goal:144-174`〕；Step-Code `plan`＝**状态机＋用 `tool_call` 硬拦写操作**〔`step-plan.ts:155-193`〕；openchamber `craft-goal`＝**提示词** | ⇒ 三种形状差得很远：**"计划模式"要真拦写**，只发提示词等于没门 |
| **权限档位** | DSH 默认 **2 档**＋custom／auto，落 session 事件、**重启仍在**〔`:188-199`〕；Step-Code **4 档**，写 `config.toml`，但 `--cycle` **故意不持久**〔`step.ts:310-318`〕；minimax **6 档**、规则源分 global／agent／session〔`permission/src/types.ts:13-36`〕；openchamber **以服务端为准**〔`permissionStore.ts:101`〕 | ⇒ 我们已定 R20＝三档／默认最严／**档位持久化**（owner 推翻过我的推荐）。**"临时切档故意不落盘"那一形（Step-Code）我们没做过**，值得进票 187 的对照栏 |
| **思考档位** | pi 分**"本次／设为默认"两种持久化**〔`interactive-mode.ts:5014-5040`〕；minimax **无命令**，只在 BYOK 配置里〔`config/src/byok-config.ts:97-112`〕 | ⇒ **正好对上票 219 那三枚按钮的形状**："本次／本次会话内／长期"在别家的思考档上是同一套三分法 |
| **导出** | DSH＝ZIP over HTTP；pi＝HTML／JSONL＋`/share`（gist）〔`slash-commands.ts:25-27`〕；minimax＝`exportCurrentTranscript`〔`feature-flow.ts:278`〕 | ⇒ 我们零枚；票 213 的候选 |
| **"这次生效了没有"的来源清单** | DSH skills＝**8 桶＋rank＋watcher 免重启**〔`skill-filesystem/src/index.ts:36-40,249-262`、`skill/src/index.ts:39`〕；Step-Code 装／卸／MCP 一律**"Restart Step"**〔`step/plugins.ts:643,1129,1185`〕；openchamber MCP **有健康状态**〔`useMcpStore.ts:34-40`〕 | ⇒ **三家分成两派**：免重启（DSH）vs 明说要先重启（Step-Code）。**"改了不生效也不告诉用户"这一形别家都不选**——这条直接进票 215 的判据 |

## ④ 该抄的形状（四条）

1. **同一枚解析器同时供"显示"与"执行"**（别家：minimax `catalog.ts:692-717`〔腿报〕）⇒ 我们若菜单显示是一套代码、命令执行是另一套，两边必然漂移。
2. **"没有 argumentHint 就拒绝带后缀"** ⇒ 未声明参数的命令不接受尾巴，别把 `/model gpt` 静默当 `/model`。
3. **名册每一项必带 `source／origin`** ⇒ 用户能看出这枚是内建的、项目的、还是插件带的（对上票 215 的"逐条列出来源"）。
4. **能力没落地就不进册**（DSH 用 `input.attachments` 闸门，`commands:391`〔腿报〕）⇒ **这条最硬**：我们今天恰好反着——`task.spawn` 的说明文字许诺了一枚**没注册**的 `task.cancel`（**票 221**）。**同一枚病，别家用机器闸门挡住了。**

## ⑤ 不该抄（只给形状级理由）

- openchamber 的**服务端 HTTP 命令 CRUD ＋ tunnel／relay**——它是 web＋electron＋手机＋VS Code **多端**产品；我们是本机单进程。
- DSH 的 `definitionId` 品牌＋**Cordis 注入**——它把命令当**插件挂载点**设计；我们没有那个插件容器。
- Step-Code 的 `registerFlag`＋**`tool_call` 拦截**——宿主是 pi 的扩展 API，我们没有那一层。
- minimax 的 **side-mode 白名单**——它有子会话投影；我们的子代理是独立任务、不投影成会话。

## ⑥ 本轮仍然没看到／没证的（诚实栏）

- **各家命令的"执行体逐枚"没读完**：只挖了 ②③ 点名的十几枚；minimax 那 45-52 枚里约 **35 枚我只看了名字没看后端**。
- **DSH 的 6 枚分母是我的尺够不到的**（它的名册由插件运行时注册，静态数不出来）——**引用时写"至少 6，实际按装的插件变"**。
- **skills 追加成命令**那一形（minimax）只在腿报里，我没复现。
- `openchamber` 无 `.git` ⇒ 全部结论只有当下快照的 `file:line`，**不能声称"他们一直这么做"**。
