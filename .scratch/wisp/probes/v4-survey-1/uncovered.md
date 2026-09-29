# 派单料 · 第四版调研腿（`v4-survey-1`）未覆盖面与可开票缺口

> 只给料：**本腿不开票、不写码、不替编排者排程**。
> 所有"我们缺"的判定只按 Go 侧现读（`internal/**`、`cmd/**`）＋工单池文件名；`frontend/**` 与 `design/**` 零读零引。
> 落点枚数＝本轮点到过路径的文件枚数（**是候选落点，不是施工图**）；`撞票`列按 `ls .scratch/wisp/issues/` 里**未带 `-done`** 的 2xx 名册（18 枚：200/201/211/212/213/214/215/216/219/220/221/222/223/224/225/226/227/228）＋代码注释里点名的票号判定，**票面正文本腿没读，派单前须自己核**。

## 第一堆·不动契约就能做（改的是自家实现与自家句子，不碰 `C1–C32`／`D1–D47`／`D43` 转移表／阈值／golden）

**P-1 人按了停止，要把"是人停的"告诉模型；并把停止做成两击。**
- 缺的是什么：我们只有"告诉人"那一半（`internal/agent/sink.go:10-15`、`EvStuck` `:30-32`、刹车调用点 `internal/agent/loop.go:470-483`）；本轮复量"两击／`cancelled by the user`"这一族**0 生产落点**。
- 参照：DeepSeek `packages/client/ui-jobs/README.zh.md:30`（首击上膛＋3 秒内确认击＋**明确告知模型**）；minimax 把 `stopping` 做成正式状态（`packages/tui/src/tui/shell/activity-line.ts:259-260`）。
- 落点：**3–4 枚**＝`internal/agent/sink.go`、`internal/agent/loop.go`、`internal/tools/cancel.go`（或 `internal/tools/task.go`）、＋面板那一节多一枚"已确认／已中止"的读回（`internal/panel/composer.go`）。
- 撞票：**不撞**。⚠ 与 221 相邻但不同物（221 管"名册许诺了没注册的能力"）。判据里必须写死那条陷阱：**抄状态不抄"停止失败"那枚旗标＝造"账上已取消、实际进程还在"的假终态**（minimax 自己的做法见 `agent-modules/background-task/src/manager.ts:120-147`）。

**P-2 每一轮"实际喂给模型什么配置"留一份证据件（不含正文、不含密钥）。**
- 缺的是什么：本轮复量 `thinkingLevel|cacheRetention|maxSerializedInputBytes` 命中的 3 处全是"模型能支持哪些档位"的清单（`internal/config/schema.go:357-359`、`internal/config/validate.go:201`），**没有"这一轮实际用了哪一档"的落盘读数**。
- 参照：minimax `agent-modules/session-report/src/contracts.ts:52-67`＋`llm-call-report-store.ts:47-53,14-15,52,57-75,77-101`（短码去重、4MB 上限、原子替换、压缩后冻结在它所描述的那一轮旁边、反悔删的是上一轮证据不是运行配置）。
- 落点：**2–3 枚**＝`internal/agent/prompt.go`、`internal/agent/spill.go`（或紧邻新增一枚）、存侧一枚（`internal/memory`）。
- 撞票：**不撞**。⚠ 提醒：这条只落盘、**不许实现成受门控的 Tool**（那是既有禁形）。

**P-3 遥测与刹车读数带上"策略是哪一版"（`policy_version` 一维）。**
- 缺的是什么：本轮复量 `policy_version|strategyVersion|strategy_version` **0 命中**；我们出的是 `task_journal` 落库＋事件汇（`internal/agent/loop.go:472-476`）。
- 参照：minimax `packages/local-runtime-v2/src/service/turn-system/runaway-guard/observation.ts:19-56`（白名单投影，现值逐字 `runaway-v1`）；预算侧的 `strategyVersion` 默认值（`agent-modules/context-manager/src/settings.ts:3-11`）。
- 落点：**2 枚**＝`internal/agent/guard.go`、`internal/agent/journal.go`。
- 撞票：**不撞**。这条越晚补越贵（改过一次梯度值，历史数据就再也分不出是哪一版产的）。

**P-4 中文的字数换算：锚定服务商给的真值、只估尾巴、宁高不低、并标明这数是估的还是真的。**
- 缺的是什么：`ApproxTokens(s) = len(s)/4`（`internal/agent/budgets.go:125-129`），逐条累加、不锚定上一条真值（`internal/agent/compress.go:110-116`）⇒ 按 minimax 注释同向推，**中文低估约 1.3–2.7 倍（这本算术是本腿算的，不是他们原文）**。
- 参照：minimax `agent-modules/context-manager/src/token-estimator.ts:6-11,20-27,31-36,38-45`＋`types.ts:28-31`（来源是值的一部分）。
- 落点：**3 枚**＝`internal/agent/budgets.go`、`internal/agent/compress.go`、`internal/agent/guard.go`（预算那一处消费点）。
- 撞票：**不撞**，但⚠ **绝对不许顺手改阈值**（`thresholds.go`／golden／SLO 一字节都不动）——这条只动**估算函数**与它的来源标注。

**P-5 "现在做不了，但你先做这一步就行"＋"一半成功一半失败"那一族句子。**
- 缺的是什么：我们只有通道级的"XX 取消不可用"四句（`internal/agent/approval/approval.go:106-114`〔已证〕）与孤立的一句诚实拒收（`internal/panel/git.go:82`〔已证〕）；**没有"你该先做哪一步"那一半，也没有分段承认的文案**。
- 参照：minimax `packages/tui/src/tui/features/session-mutation/copy.zh-Hans.ts:23,33,88-89,101,109,113-114,160-163`、`packages/tui/src/tui/commands/catalog.ts:760-766`；openchamber 两行式等待文案 `packages/ui/src/lib/i18n/messages/en.ts:1345-1346`。
- 落点：**3–4 枚**＝`internal/panel/composer_dispatch.go`（把 `ErrNoHandlerAttached` 那句扩成"要先做什么"）、`internal/agent/approval/approval.go`、`internal/panel/git.go`、＋一枚句子常量表。
- 撞票：**可能撞 219／220**（这两枚管的就是批准面与"等批准看不看得见"上的句子）⇒ 派单前必须读票面，别把同一句文案开两遍。
- ⚠ 顺带一条**几乎不要钱的形状**（照抄 minimax 的表头纪律）：**新增一条判据就必须同时新增它的文案，漏一条编译不过**（`reason-format.ts:7-9`）；Go 侧等价物＝穷尽 switch。

**P-6 配置读不懂不许静默：读失败禁用编辑＋说明当前在用哪套＋先备份再修＋"恢复全部默认"不覆盖读不懂的数据。**
- 缺的是什么：本轮未量到任何"读不懂就明说"的处理；我们的强项在另一面（写了不管用就响亮拒收，`internal/config/unwired.go` 147 行、6 枚键〔已证〕），**管的是键名、不管那份文件读不读得懂**。
- 参照：DeepSeek `packages/client/ui-shortcuts/README.zh.md:36`。
- 落点：**2 枚**＝`internal/config/parse.go`、`internal/config/manager.go`（或同族那一枚）。
- 撞票：**挨着 223／226**（223＝热重载零生产调用者；226＝"总是答复"覆写整份配置快照、藏掉手改）⇒ 强烈建议**与这两枚合并判据**，不要另开。

**P-7 会话的最小出口：把这一次对话导出成一个文件（先不做分享）。**
- 缺的是什么：本仓 Go 侧无导出出口；第三版整版没写过这一条（见 `docs/reports/missing-features-2026-09-29-v4.md` §5-9）。
- 参照：DeepSeek 命令只回"已请求"、真文件走带鉴权路由、**明确拒绝带路径参数**（`packages/session-query/session-log-export/src/index.ts:78-90`）；Pi 系默认 HTML、后缀决定格式（`pi-upstream/packages/coding-agent/src/core/slash-commands.ts:25-27`）；**上限分层回落、放弃时给可枚举理由**（`Step-Code/packages/coding-agent/src/step/feedback/bundle.ts:19-26,88-97`）；**先列清单再要同意、每个外部字段各洗一遍控制字符**（`:110,176-185`、`consent.ts:73-92`）。
- 落点：**2–3 枚**＝`internal/agent/journal.go`（数据源）、`internal/memory`（读会话）、＋一枚新文件做导出。
- 撞票：**不撞**。⚠ 若同时想做"提交反馈"，同意页那一套（先摆清单再要同意）建议独立一枚判据，别塞进导出。

**P-8 名册不许许诺没注册的能力：机器闸门＋响亮失败。**
- 现象（三行定性，别读成安全事件）：
  - 出现在哪：派子代理那枚工具的说明文字许诺"可以用 `task.cancel` 单独停它"（`internal/tools/subagent_197.go:196`），而 `task.cancel` 在名册头上逐字 DEFERRED（`internal/tools/task.go:23`，本轮复量）。
  - 本机有没有被入侵的证据：**没有**；这是**自我矛盾的名册**，不是漏洞利用。
  - 最坏后果形状：用户或模型照着说明去使一枚不存在的工具 ⇒ **一次点了没反应／一次报错**，并且"说明文字"这件事从此没人信。
- 参照：DeepSeek 组合期闸门（`packages/plan/plan-mode/src/index.ts:230-231`）＋运行期能力位（`packages/interaction/permission-presets/src/index.ts:279-282`）＋前端逐会话 `available()`（`ui-commands/src/client/contract.ts:88-89`）；minimax `audience:'internal'`＋`visibleWhen` 双重（`packages/tui/src/tui/commands/catalog.ts:47-84`）。
- 落点：**2 枚**＝`internal/tools/task.go`、`internal/tools/subagent_197.go`（或注册表那侧一枚）。
- 撞票：**就是 221**（在飞）⇒ **不该开新票，只该给 221 补一条判据**："装配时没有真执行者的命令，注册这一步就该失败，而不是注册了再置灰"。

**P-9 防转圈降粒度到"每一次动作"，并补两类信号＋两张豁免表。**
- 缺的是什么：本轮复量刹车比的仍是整轮签名（`internal/agent/guard.go:166` 调用、`:234` 定义），**无 ABAB／无错误家族／无按工具名豁免／无运行时总开关／无"轮意图"豁免**。
- 参照：minimax `agent-modules/runaway-guard/src/{contracts.ts:3-9,177-180,guard.ts:28-31,signals.ts:174,199-217,step-view.ts:153-171,388-409}`；豁免表 `.../runaway-guard/tool-policy.ts:23-31,34-62`；轮意图豁免 `.../runaway-guard/extension.ts:21`；被拒动作剔除 `step-view.ts:61-65,116-118`；短码用每轮换密钥打 HMAC、超预算**整条不检测**（`fingerprint.ts:8,21-23`、`step-view.ts:149,160,187`）。
- 落点：**2–3 枚**＝`internal/agent/guard.go`、`internal/agent/loop.go`、（豁免表可新增一枚）。
- 撞票：**不撞**，但⚠ 这一族里"撞线必停 vs 只提醒不停"的**取向不许被这次对标顺手改掉**——要过 `D43` 转移表就是第二堆的事。

## 第二堆·碰契约（须人工批准；本腿一律不排顺序）

**Q-1 子代理深度／并行上限做成设置项。** 现状：编译期常量（`internal/tools/subagent_197.go:73,79`，本轮复量）＋`internal/config/schema.go` 搜 `subagent` 0 命中。参照：DeepSeek 一页两半一次保存、两次写入互相独立、过期草稿报冲突不覆盖、**计口径写在数字旁边**（`ui-settings-subagent/README.zh.md:12,28,32`＋`src/client/subagent-limits-card-controller.ts:10-13,36,48`）。**碰**：D36 配置模型 section 树＋D38/D32 上限归属＋票 211（池不许大于桥天花板）。落点：≥4 枚（`internal/config/schema.go`、`internal/tools/subagent_197.go`、`internal/panel/composer.go`、面板读写侧）。

**Q-2 "允许"这个出口落在哪一层（面板／球长按／语音反义）。** 现状：面板只能拒（`docs/specs/SPEC-08*.md:156-176` C17 名册逐字），"仅本次允许"在界面上没有出口（本轮复量 `ask_user` 等仍 0 命中）。参照：DeepSeek 面板两枚决定＋持久策略划归服务侧（`ui-approval/src/client/contract/slots.ts:66`、`README.zh.md:36`）；openchamber 三枚＋把要存的规则文本印在按钮上（`packages/ui/src/components/chat/PermissionCard.tsx:349-358,407-442`）；minimax 宽度候选枚举＋**第一档永远是最窄的**兼容律（`agent-modules/permission/src/types.ts:96-106,151-195`）。**碰**：C17 方法白名单（未定稿、阻塞中）＋`Q-49`＋票 219／224。

**Q-3 "没有能答的人就当次判拒绝"这一支。** 现状：本轮复量"没有批准的通道"0 命中；别家在这一格**方向相反**（minimax 判拒绝 `permission-core.ts:422-430`；openchamber 转成等人 `packages/web/server/lib/permission-auto-accept/modes.js:5-6`）⇒ **没有现成答案可抄**。碰：D2 深度睡眠拓扑＋D4 门控＋票 201／220＋票 197 子代理那一半。

**Q-4 "这一条是谁发起的"做成一等业务数据（13 枚来源词表）。** 现状：本轮复量 `ConversationSource|MessageSource` 0 命中，最接近的是内容来源（`internal/risk/provenance.go`）。参照：minimax `agent-modules/conversation-contract/src/index.ts:1-14,16-29,357`。碰：D39 契约深化（把只有名字的契约变成规格）＋票 224（会话身份）＋D12 定时。落点：≥3 枚（事件汇、名册、投递入参）。

**Q-5 20 枚状态名里补"正在压缩"这一名（以及重连、正在停）。** 现状：本轮逐名读过 `internal/statemachine/states.go:11-31`，**无这三名**；别家有（minimax `activity-line.ts:252,256-260`）。**碰：D43 权威转移表（冻结那张表）** ⇒ 人工批准级，不是票级。⚠ 本腿只报"界面上说不出这一刻"这一现象，不提议怎么改那张表。

**Q-6 面板里那一格真终端。** 现状：刷给界面那包内容仍是固定 6 枚键（本轮读 `internal/panel/composer.go` 的 `Snapshot`）。参照：DeepSeek `ui-sidebar-terminal` 一整包，**且恰好不在他们官方包地图里**（本轮复核：那张表 50 行、磁盘 52 枚）。碰：C17 白名单＋D34 内置工具权威表＋D29 视觉架构（多一格里布局）。

**Q-7 插件钩子层（11 个钩子点＋"能改哪些东西"的红线清单）。** 现状：本轮复量那 11 个钩子名 0 命中；全仓 `hook` 一词的非测试命中全是系统语义（键盘钩子／关机钩子／桥回调）。参照：minimax `agent-modules/plugin-hooks/src/{contracts.ts:2-14,155-191,runner.ts:520-529,30-38}`；openchamber 的信任模型（`extensions/DOCUMENTATION.md:11-19`，且**内置登记表是空的** `registry.json:1-3`）。碰：D19 插件信任＋D46 Tier-1＋D3 插件运行时 ⇒ 契约级，本腿只把别人的红线清单摆出来。

**Q-8 要不要把二次判定搬到云端。** 现状：我们把这类双轴判定挂在 RESERVED（`docs/PLAN.md:1534` 逐字理由）。参照：minimax 真搬到了服务端并写明隐私代价（`classifier/cloud-classify-client.ts:5-6,26-39`、`cloud-gateway.ts:10-21,83-84`、`conversation-renderer.ts:73-99,110-113`、`http-cloud-gateway-client.ts:11-16,165-172,139-149`、`cloud-gateway.ts:110-118`）。碰：D4／D33／延迟预算＋**新开一条隐私账**（两条账不许互相冒充）。

**Q-9 MCP 在菜单里给个位置（但不接）。** 现状：本轮点到行的五家**没有一家把 MCP 放进加号菜单**，只有 Step-Code 的 `/mcp` 只读回显（`packages/coding-agent/src/step/slash-commands.ts:45-55`）。碰：D13 取舍（**只做"留位置"就不碰**；一旦接就推翻 D13）。事实层必须补一句：`minimax-code/packages/agent-modules/mcp/` 是一份**可读的运行时**（13 枚件，含连接池 911 行＋持久名字治理 258 行）⇒ 以后**不能再写"无处可抄"**。

**Q-10 字号／主题新键、`websearch` 那一节、深链协议注册、"恢复来的定时档期先只登记不执行"。**
都是"界面名词上有、配置面里没有"的格：字号（DeepSeek `ui-theme` 明写"主题与正文字号"；我们只有球尺寸与面板缩放 `schema.go:165,532-538`）；
`websearch` 三字段（DeepSeek `ui-settings-web-search/README.zh.md:2,68`；我们 19 个 section 名册里没有，⚠ 且这条撞 AGENTS §2 未定案的 `web.search` 实现路径——**本腿不定路径**）；
点通知直达会话（openchamber `mobile/.../deepLinks.ts:23-30,46-53`＋冷启动先存住 `deepLinkNavigation.ts:15-16`；我们 0 命中）；
隔离恢复（minimax `agent-modules/cron/src/registry.ts:42-48,114`）。
碰：D36 section 树＋D41 部署打包＋D7 平台范围＋D12 定时 ⇒ 逐枚人工批准。

**Q-11 球／托盘上"几件事在等你"那个数。** 数据源现成（等批准队列长度），**但载体不在跑任务的那个进程里**（票 228 的标题就是这个）⇒ 这条**必须排在 228 之后**，本腿不排。参照：openchamber 角标＝不同的提醒条数、三路信号清（`APNS.md:40-66`、`AppDelegate.swift:201-208`）；反面参照：做 VS Code 扩展的那家**一枚状态栏条目都不做**（`grep -rn createStatusBarItem src/*.ts` 0 命中）⇒ "常驻视觉件该摆什么"，连做了四端的人都在克制。

## 本轮读不到／判不了（不许当成"没有"）

1. **Step-Code 与 pi-upstream 的设置项名册**：本轮没读两家的设置页正文 ⇒ 表二里这两家的列是空格，**空格＝没读到**。
2. **各家命令的执行体逐枚**：minimax 那 45–52 枚里约 35 枚只看了名字没看后端（继承 `survey-2026-09-29-command-catalog-across-harnesses.md:58`，本轮未补）。
3. **DeepSeek 的 6 枚后端命令分母**：名册由插件运行时注册，静态数不出来 ⇒ 引用只能写"至少 6，实际按装的插件变"（同件 `:59`）。
4. **openchamber 的服务端字段表**：无 `.git`＋快照里缺依赖 ⇒ 凡"服务端到底回什么字段"都只能算客户端反推。
5. **`D:\work\AI\open source\pi` 不是编码 harness**（37 枚文件／0 枚 `.ts`，本轮复核）；**`minimax-code/third_party/pi-mono/` 是 minimax 仓里的一份 Pi 副本**，本轮未读 ⇒ 任何"minimax 有 X"要先确认出处不在 `third_party/`。
6. **DSH `0.1.6-alpha.1` 那个版本号本轮没验出来**：`~/.dsh/profiles` 根目录没有 `package.json`，只有 `desktop`／`web`／`node_modules` 三项（`node_modules/@deepseek-ai` ＝257 枚）⇒ 派单里那半句我不复述。
7. **提醒防污染后缀在我方有没有等价物**：两轮都没量（第三版 4.7 同样标"没查"）。
8. **票面正文**：本腿只看了工单池文件名与代码注释里的票号，**没读任何一枚票面**；"撞不撞在飞的票"因此是**初判**，编排者派单前须自己核票面。
