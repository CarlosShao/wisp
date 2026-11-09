# 票池定档 R6（pool-audit 重试腿，上一腿失败返 null）

- **腿名**：pool-audit-r6（写腿·docs＋一枚顺手票）。**时刻**：2026-11-09 13:1x–13:5x +08。
- **锚点**：起手 HEAD `71ff0fc87`（dev）；终态 HEAD＝`b3aefc34b`（本腿顺手票 239 残余①）。
- **任务**：90 张未结票分三档（甲＝立即可派／乙＝归口等票／丙＝要人·桌面·不适用）＋按价值排序的「立即可派」清单。
- **方法纪律**：逐票读正文三处（`**Status`／`**Blocked by:`／AC 首格）＋每个 blocker 现跑核 `-done` 后缀（`ls .scratch/wisp/issues | grep '<N>-'`），⛔ 只扫包名、⛔ 信转述。

## 0. 起手三核读数

1. **写面核**：`git status --porcelain` 无已跟踪文件改动（全部为 untracked 临时件）⇒ 共树写面干净。
2. **池核**：`ls .scratch/wisp/issues/*.md | grep -v README | grep -v -- -done | wc -l` ＝ **90**（与派单一致）。
3. **在飞核**：本腿施工窗（13:2x）共树现读 porcelain 出现**别腿写面**：`internal/agent/loop.go`／`internal/agent/sink.go`／`internal/llm/**`（4 枚 adapter/request）已修改＋`internal/llm/thinking_levels.go` 未跟踪 ⇒ **A 腿（cmd/wisp 族写腿）正在飞**。本腿写面（`internal/ball/**`＋docs＋票 239）与其零交集；门禁测试色一律在仓外导出树跑（见 §4）。

## 1. ★派单前提过期清单（现跑核出，⛔ 按票面硬做的那些坑）

| # | 派单/票面说的 | 现跑实况 |
|---|---|---|
| 1 | 「46 卡 41-44」 | **41/42/43/44 全 `-done`，45-diagnostics-guards 也 `-done`** ⇒ 46 的票面 blocker 已全清；但 46 正文要求 real desktop／真触发 pre-mortem ⇒ 丙（真机桌面），不是卡票 |
| 2 | 「53 卡 48/51/52」 | **51、52 已 `-done`** ⇒ 只剩 48 一个 blocker ⇒ 乙（等 48） |
| 3 | 票 23/33 等票面「Blocked by: 20…」「07/12…」 | 20、07、12 全 `-done` ⇒ 23、33 纸面已解锁（甲），但票面 09-19 旧、起手必须复量 |
| 4 | 票 168「前置＝票 50＋51 都没做完」 | **50、51 全 `-done`** ⇒ 前提半过期；剩 C17/C24 白名单人工批准面 ⇒ 丙 |
| 5 | 票 256「等票 265 落地后复判」 | **265 已 `-done`（四格全勾）** ⇒ 256 的复判腿前提已满足 ⇒ 甲 |
| 6 | 票 262 要补「跟踪路径长度尺」 | **`scripts/check-path-length-budget.sh` 已在仓且被 230-r1／216-r2／230-v1 等多腿引用 GREEN 读数** ⇒ 262 要补的那格已被别的票兑现 ⇒ 建议编排者收口（前提过期） |
| 7 | 票 245「已立，未派」 | **AC#1/2/3/5 已落地并验收**（测试日志里「cancel hotkey left unbound while idle…(ticket 245)」就是已落的行为）；剩 AC#4（逐字阻塞于 228）＋AC#6/7/8 ⇒ 甲（AC#7/8 可派） |
| 8 | 票 148「Status: blocked」 | 过期：A1591 已批＋148-r3 已落 AC#1/2/3＋AC#2 已翻（wave105/编排者）；剩 AC#1/AC#3 非实现者复判＋票面三处过期待编排者收 ⇒ 乙 |

## 2. 三档总表（90 张）

图例：blocker 栏「逐字」＝票面原文摘引；「✔done」＝本腿现跑确认该 blocker 已 `-done`。

### 甲＝立即可派（24 张；每张给可做的具体一步）

| 票 | 一步（具体） |
|---|---|
| **33** panel-host-c27 | **全场网关票**（blocker 07✔done/12✔done）。起手只读复量片结构（09-19 票面旧；33-r1 六环 `composer_dispatch.go` 等后续件散落多票），接已复量确认未落的那一片 ⇒ 解锁 35→36/37/38/39/40→195/180 整条链 |
| **228** ball/tray 宿主 | AC#0 已裁走甲形。起手跑 AC#1 尺（`internal/ball` 生产引用者点数），把审批门＋任务通路接进常驻腿 ⇒ 解锁 201/219/246/123/245-AC#4 |
| **23** system/window 工具族 | blocker 20✔done。起手＝名册差集裁定表（D34/SPEC-07 §3 ↔ 生产注册表，163 AC#1b 同形尺），差集非空不许进实现；然后 A/B 并行槽 |
| **145** 快照载体 | AC#2 写码腿可派（票面⑦「放行条件已满足」），**先解⑥**（读末节）再动 `internal/panel` |
| **187** 模型/思考档可改 | AC#1 起手＝现量读者（`grep -rn "ThinkingIntensity\|roles\." --include=*.go … \| grep -v _test`），逐支定生效档（热改/重载/重启） |
| **188** 任务监视栏源 | AC#2：给 `Record` 补 **D43 逐字对齐**的状态维并进快照（判据钉「值来自真名册，不许硬编码」）；AC#1 已勾勿重做 |
| **194** 面板方法名册对齐 | owner 已批（A381）。AC#1 起手＝18 枚逐枚代价表（现跑，四问：等价实现/落点/撞禁区/门控审计） |
| **191** 面板内嵌终端 | AC#0 已按甲裁（A379）。AC#1 只读设计核：现有工具面哪枚能承载「跑命令＋流式回显」（票 163 常驻会话在不在名册级，现跑） |
| **192** 浏览器/产物预览 | AC#0 已裁（抓结果页）。AC#1 产物预览先拆出来做（本地路径 vs 内容字节、截断带找回指针） |
| **248** 面板设置路由 | 248-v1c 已裁；AC#10 ⓘ 支已拆票 256/265（265 ✔done）。起手复量残余格 ⇒ 接 Go 侧 settings route（界面归外部前端） |
| **213** 斜杠命令目录 | AC#1 起手＝现量注册表，名册可枚举＋每条真执行者；正控「塞假命令必红」（AC#2） |
| **214** 附件接线 | AC#1 受理器生产调用者（`cmd/wisp` 装配路径 grep 点数 0→≥1）⇒ 占 cmd/wisp 写面槽 |
| **215** 技能/插件清单 | AC#1 清单可枚举＋真装配读者；正控「配置 true、没人读 ⇒ 红」（AC#2） |
| **197** 子代理实体 | AC#0 先答「归哪个切片」（只读普查 S3/S5/S7，答不出摆 owner——这一格有丙成分） |
| **323** 判据残余 | 派 **`323-a1` 只读普查「台件该怎么造」**（写面仅 docs/evidence＋probes，零产码）：明文落盘敏感尺/超时夹具/上游错误体/宿主装配/protocol 预检五格各给最小台件方案 |
| **314** cmd/wisp 真窗进程死 | 诊断腿可派：4/4 复现已录、归因候选 4 枚，AC#0/#1 定性（本机真窗可复跑） |
| **302** 仪器放置 | 把三枚 page-reply 判据（AC13 冷启动/两枚 AC14）挪出 CI CLI 档（`cmd/wisp` 测试面写腿槽） |
| **274** 出货 exe 页面钉 | 收紧 `TestPanelBundleShapeSeparatesAnchorFromRealPage_AC12` 的两形放行＋`build.ps1` 第 2 步过期注释（`cmd/wisp` 测试面＋scripts 槽） |
| **238** 数据根创建时序 | 起手复量前提（`internal/proc` 对 winsec 导入数＝0 是否仍真）⇒ 修 `internal/observe/logging.go:72` 创建顺序（写面 internal/observe，不在禁区） |
| **388** 测试播报音量衰减 | 无 blocker。cgo 真浮点衰减（59-r1 rig 已证可行，**⛔ 整型路**）；写面 `cmd/wisp` 测试面或 `internal/audio/**`；需真音频设备（本机在） |
| **245** Esc 热键残余 | AC#7（借被拒记账归位＋正控）／AC#8（"桌面已有占位者"与"没绑上"分两条红句，⛔ 不许改跳过）＝`internal/ball` 判据修理可派；AC#4 逐字阻塞于 228 如实记 |
| **112** winsec CI 门 | ready-for-review ⇒ 派非实现者复判腿（只读；09-21 前提旧，起手复量 winsec 门现色） |
| **193** 轻量级解套 | AC#1 只读列门禁清单可派（逐枚"卡不卡前端/卡在哪一句"，逐字引原句）；**AC#0 owner 未答前不修 panel/tokens** |
| **256** 常驻门配置 | **前提已满足（265 ✔done）**：派非实现者复判腿复判 AC#1/AC#2（Grants 半格＋值正控），凭据复用票 224 r2 仪器 |

### 乙＝归口等票（30 张；blocker 逐字＋现核结果）

| 票 | blocker 逐字（＋现核） |
|---|---|
| 25 S3 验收 | 「Blocked by: 21-approval-gates-minimal, 22-web-tools-d30, 23-system-window-input-tools, 24-doc-search-tools」——21/22/24 ✔done，**只剩 23（甲）** |
| 35 panel-bridge | 「Blocked by: 33-panel-host-c27, 34-frontend-scaffold」——33（甲）/34（丙）皆未结 |
| 36 结果历史 | 「Blocked by: 30-result-routing-d10, 35-panel-bridge-c17」——30 ✔done，剩 35 |
| 37 approval-ui-l2 | 「Blocked by: 21-approval-gates-minimal, 35-panel-bridge-c17」——21 ✔done，剩 35 |
| 38 命令面板 | 「Blocked by: 35-panel-bridge-c17」——剩 35 |
| 39 配置编辑器 | 「Blocked by: 05-config-model, 35-panel-bridge-c17」——05 ✔done，剩 35 |
| 40 安全页 | 「Blocked by: 29-memory-l1l2, 35-panel-bridge-c17, 44-costmeter-c23」——29/44 ✔done，剩 35 |
| 48 审批队列全量 | 「Blocked by: 37-approval-ui-l2, 47-task-scheduler-pathlock」——47 ✔done，剩 37 |
| 49 会话授权 | 「Blocked by: 04-sqlite-core, 48-approval-queue-full」——04 ✔done，剩 48 |
| 53 lark-cli e2e | 「Blocked by: 48-approval-queue-full, 51-tier2-goja, 52-d46-command-plugins」——51/52 ✔done（过期），**只剩 48** |
| 54 S7 验收 | 「Blocked by: 47-task-scheduler-pathlock, 48-approval-queue-full, 49-session-grants-d45-2, …」——47 ✔done，剩 48/49 |
| 71 门自报 | 「Blocked by: 70-ci-actually-green（只因为两者都要改 .go，不是语义依赖）」——70（丙）未结 |
| 85 lint 判语 | 「open（**排队**：本票要改 ci.yml 与全仓多个包，必须等写码槽位空出来…）」 |
| 122 staticcheck 清仓 | 「open（…来源=audit-85-preflight…）」——排在 85a 之后，85 未结 |
| 123 CLI 测试待人批 | 「open」——L1/L2 在 CI 上无人能答 ⇒ 归口 228 宿主（201 同源） |
| 148 clone/replay | 票面「blocked（两枚前置都是人工批准）」**已过期**：A1591 批＋148-r3 落地＋AC#2 已翻 ⇒ 剩 AC#1/AC#3 非实现者复判＋编排者收票面三处过期 |
| 167 面板四小出口 | 「**ready-for-agent，但立而不派**（排在票 145 之后——它俩抢同一批快照字段…）」 |
| 174 task.output 路径 | 剩 AC#4「…须 owner 单独批准且批准原文落台账；未拿到批准前本票对这些文件零字节」＋AC#5 结题陈述 |
| 180 面板宽度 | 「**判定=blocked-by 票 33/35，⛔ 派单**…Status 名义 ready 实为 gated」 |
| 195 config.set 地基 | 「**blocked（排在票 33 片 B 交件之后；地基未建之前不许接面板口）**」 |
| 201 审批没人能答 | 「剩余的 AC#1／AC#2／AC#3a／AC#6b 全部归口同一枚前置——**宿主原生入口（球／托盘）在 wisp run 里没接**」＝228 |
| 216 菜单显示底座 | 残余＝①「两份合一」需新导出名＝契约面、③「attachments」拆键（已按规矩停手上报）、⑤ cmd/wisp panel_pump 半句；AC#5/AC#6 框等非实现者凭 216-v3 判语翻 |
| 219 审批卡三按钮 | AC#1「三枚按钮各有真执行者（生产路径…）」⇒ 归口 228 宿主；AC#6 前端半已落等复认 |
| 220 名册 join | 剩 AC#5「**整包终态读数**：PATH=… go test ./cmd/wisp ./internal/… 跑到终态」⇒ 需 cmd/wisp 整包安静窗（与写腿互斥） |
| 230 关闭票内未竟格 | 230-v1 判语「AC#2 的翻勾归编排者凭本判语翻；AC#1 归 115-v1」⇒ 等编排者收 |
| 239 球方形 | **残余①本腿已落（`b3aefc34b`）**；残余②「cmd/wisp 全量没在 HEAD 复跑…归口＝编排者下一枚整包门禁」＋残余 a（真窗差分常驻化）⇒ 等编排者收口 |
| 246 veto 执行体 | 「已立，未派（触发＝验收腿 245-v1…）」——执行体＝球在常驻进程＝228 甲形 ⇒ 同面串行 |
| 297 安静房 0.37 | 「未开工…新排程：292-v1（在飞）→ 296-r1 → 297-a1 → …」——排队 |
| 305 活文档 | 两症状须先分（AC#0）；判据世界＝真窗族，314 未定性 ⇒ 等 314 |
| 313 默认 stepfun | 机主原话「配置面必须是界面」⇒ 归口 248（Go 侧）＋外部前端 |

### 丙＝要人／桌面／不适用（36 张；具名理由）

| 票 | 理由（逐字/具名） |
|---|---|
| 32 S4 验收 | 真机场景 ③④；票面逐字「**human (user) sign-off required for both subjective scores**」；blocker 27/29/30/31 ✔done（过期但本体要人） |
| 46 S6 验收 | 「SLO runs (sampler, **real desktop**)…Pre-mortems **real-triggered**: suspend/resume…battery…mic exclusive-hold」；blocker 41-45 ✔done 全清（过期）但本体要真机桌面 |
| 55 S8 macos | 「Blocked by: 54-s7-acceptance **+ self-use week + explicit user go-ahead**」 |
| 56 S8 签名分发 | 「Blocked by: 54-s7-acceptance + **SignPath qualification (P8) + user go-ahead**」 |
| 57 S8 插件 SDK | 「Blocked by: 54-s7-acceptance + 53 lark-cli example + **user go-ahead**」 |
| 58 S8 i18n | 「Blocked by: 54-s7-acceptance + **user go-ahead** + English voice model selection」 |
| 60 C32 realtime | 「**Gate condition #2: user provides a realtime provider Key** — if either condition is missing, this ticket AUTO-DEFERS」；条件#1（59/16/28/44）✔done 全满足，#2 等机主 key |
| 34 前端 scaffold | 写面 `frontend/**`＝机主外部 agent 地界（blocker 01 ✔done） |
| 62 球玻璃质感 | review；质感归 65（缺参考图） |
| 64 球缺陷 | review，「剩 4 框各归其主」 |
| 65 玻璃重做 | 「**blocked-on-owner（缺参考图，见 §第一步）**」 |
| 68 球默认视差 | review，「AC#2/AC#3 **blocked-on-owner**（R15 #3/#4/#5）」 |
| 70 CI 真绿 | 「**blocked-on-owner**…剩下每一件事的判据都不在代理手里」 |
| 72 ATBLE | 「AC#4 的 runner 半边与 AC#6 的改名**都不在代理手里**」 |
| 77 前端脚手架 | 「**blocked（09-23，owner 把前端收走了）**」 |
| 100 AppContainer | 「**RESERVED**…AC#0 那条触发门没满足之前，本票不许开工」 |
| 163 常驻 shell 会话 | 「**立而不派**（要扩 shell.exec 的语义＝动 D34/D46 那条线，**人工批准**…）」 |
| 168 内置浏览器插件 | 「**立而不派**（前置＝票 50＋票 51 都没做完…且要动 C17/C24 未定案白名单＝**人工批准**）」——50/51 ✔done 半过期，剩人工批准面 |
| 169 frontend 装饰字 | 「**owner／那一支写着 frontend/** 的会话的接手件——本编队结构性写不到那儿**…不派、不 claim、不加 -done」 |
| 170 CLI L1 无 veto | 「**立而不派**（推荐那一支很可能要碰 internal/agent/approval/**＝冻结面；且…是产品口径）」 |
| 177 戳 vs 续读 | 「现在挡在写手腿前面的…是 **Q-63（owner 一句话…）**」；AC#5/AC#6 属结题/排程格 |
| 247 捕捉栈零进口 | 「AC#6＝不可判…**『无权限／禁用设备』那一形要机主一句话**」；AC#2 已判不成立归 297 |
| 249 panic sink | 「**WITHDRAWN（…自撤；⛔ 任何腿不许据此开工）**」 |
| 262 路径长度门 | **前提过期**：`scripts/check-path-length-budget.sh` 已在仓并被多票引用 GREEN ⇒ 建议编排者收口 |
| 264 错误文本带绝对路径 | 「**⛔ 取舍归机主，本格不派写腿**」 |
| 269 gofmt 样本 | 「**WITHDRAWN／已并案…⛔ 任何腿不许据此开工**」（并票 236 AC#6） |
| 276 gofumpt 普查 | 「本票立票是为了让『不做』有对象、有完成判据，**不是为了排腿**」 |
| 290 门挂没人开 | 「AC#3 未勾——欠的是**机主在场那一发真机读数**」 |
| 291 level 尺 dc | 「AC#2＝**机主在场那发四形真读数**…⛔ 派不出去」；AC#3/4 按在其后 |
| 293 托盘静音勾 | 「盘上只剩 1 枚＝AC#5（**肉眼那枚勾，归合并真机窗口**）⇒ 本票不改 -done」 |
| 311 winsec 门色 | 残余＝Q-311-r1-1 退码语义未定案＋票 310 三方言**盖章人＝编排者**（口令挂着）⇒ 要人 |
| 312 扫本机配置 | 「机主明说**暂时不用开工** ⇒ 本票当前＝**只立不派**」 |
| 316 M 形假绿 | 「**withdrawn（撤票；口令「316 撤」已用）**…与票 272 重复」 |
| 346 竞争窗仪器 | 5 翻/2 需 Linux：两枚仪器首行 `//go:build !windows` ⇒ 本机 Windows 档 ⛔ 编译，需 CI/Linux 环境 |
| 379 CRLF 三件 | 「**VOID**…★立案前提是假的」 |
| 407 K5 插件面板 | 「no ui no spec unowned」——「前端 face 归外部 agent；Go 侧 bridge 若有缺口可另立子格」⇒ 先要归口/规格 |

## 3. 「立即可派」按价值排序（甲 24 张）

| 序 | 票 | 为什么排这里 |
|---|---|---|
| 1 | **33** | 全场网关：解锁 35→36/37/38/39/40→195/180 七张 |
| 2 | **228** | 解锁 201/219/246/123/245-AC#4 五张＋审批终于可答（产品主链） |
| 3 | **23** | 解锁 25（S3 验收）；D34 主表最后一块大空地 |
| 4 | **145** | 解锁 167＋十一态界面输入；票面已具名「可派」 |
| 5 | **188** | 小步：AC#2 状态维进快照，直接喂监视栏 |
| 6 | **187** | 机主点名功能（模型/档位可改） |
| 7 | **194** | owner 已批名册补齐；18 枚逐枚代价表起手 |
| 8 | **248** | 机主「界面配置模型」的 Go 侧正身（313 的前置） |
| 9 | **191** | AC#0 已裁甲支，只读设计核起手 |
| 10 | **192** | AC#0 已裁，产物预览先拆 |
| 11 | **213** | 斜杠目录（composer 族地基） |
| 12 | **215** | 技能/插件清单 |
| 13 | **214** | 附件接线（占 cmd/wisp 槽） |
| 14 | **323** | 只读普查腿，零产码、零写面冲突 |
| 15 | **256** | 复判腿（265 ✔done，前提刚好新满足） |
| 16 | **245** | AC#7/AC#8 internal/ball 判据修理 |
| 17 | **314** | 真窗死因诊断（305/302 的世界稳定性前提） |
| 18 | **302** | CI 常红灯清零（cmd/wisp 测试面槽） |
| 19 | **274** | 出货 exe 页面钉（同槽串行） |
| 20 | **238** | internal/observe 创建时序（独立小票） |
| 21 | **388** | 测试不再吵机主（体验＋cgo 仪已证可行） |
| 22 | **197** | AC#0 归片普查（有 owner 摆题成分） |
| 23 | **193** | AC#1 只读门禁清单（AC#0 未答不修 panel） |
| 24 | **112** | 只读复判腿（09-21 前提旧，起手复量） |

**互斥提醒**：① `cmd/wisp` 族（214/302/274/388＋245-AC#4、220-AC#5、239-残余②的整包读数）同波只许一枚写腿；② `internal/panel` 族（145/167/180/187/191/192/194/248/216 残余）同波只许一枚写腿；③ 本腿顺手票占用过 `internal/ball`（已交件，已释放）。

## 4. 本腿顺手票（§要做第 5 条）：票 239 残余①

- **为什么是它**：甲清单里写面不涉 `cmd/wisp`／`internal/panel`／`internal/agent`／approval／`frontend`／`forks`／`internal/winsec` 的最小可结格＝239 残余①（写面 `internal/ball/**`＋docs）。
- **落地**（`b3aefc34b`，8 文件）：`hit.go` 新增 `WaveTravelPx=16.0`（带 R+21.75/0.25px 推导注释）；`renderer_windows.go` 裸字面量 16 改用常量（值零变化）；C21 几何表新行＋`c21MotionGolden` 补 pin；新常驻尺 `Test239WaveTrainStaysInsideTheWindowWithAMarginFloor`（余量下限 0.25px＋44..160 全档 box 严格小于窗半）。
- **门禁**：共树 gofumpt 空／vet rc0／build rc0；导出树（`git clone --shared` @ `b3aefc34b`）新尺＋C21 双向锁 PASS、d22scan rc=0 clean（ban #8 internal/=967）；**去牙突变**走距 16→17 ⇒ 新尺红（"leaves only -0.750px"），还原复绿。导出树另有 3 枚「内嵌资产 CRLF/sha」红＝**autocrlf 检出伪影**（blob CR=0 / 工作副本 CR>0，控制脚本在 `.scratch/wisp/probes/239/r1-crlf-control.sh`），主树同选器全绿，非本腿缺陷。
- **具名没做**：残余 a「239-v1 差分脚本变常驻判据」的**真窗像素差分半**要独占桌面，本腿不擅自造（未定义即停）；残余 b「cmd/wisp 整包 HEAD 复跑」归编排者整包门禁。票面零框翻（AC 五格早已全勾）、Status/票名未动、⛔ 不加 -done、⛔ push。

## 5. 给编排者的一句话

甲 24 张里最大杠杆是 **33→35 链**与 **228 宿主**；本腿已释放 `internal/ball`，`cmd/wisp` 与 `internal/panel` 槽位均未被本腿占用。撤销口令：「撤 239-r1」。
