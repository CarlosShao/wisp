# 分诊 · batch-4（票 188–238，共 20 枚）

- 只读分诊，**未改任何票面/产码**，未跑 git 写命令。起手时间 `2026-11-05 09:2x +08`。
- DONE-skip 核对：本批 20 枚在工作树里**全部只有活件**（无 `-done` 变体），故**零枚 DONE-skip**，逐枚读件。
- 分档口径按派单：T1＝Go/脚本侧今天就能干（非 `frontend/**`、非真机桌面 GUI/音频、非 `cmd/wisp` 整包死区、非纯机主拍板/改契约文字）；T2＝被另一枚票卡住；T3＝要机主/外部（frontend、真机、CI 色、钱、改契约文字、不可逆、纯机主决定、已撤）。
- 每行「阻断者/首动作」列：T1 写**第一枚可动 AC**，T2 写**阻断票号**，T3 写**卡在谁手上**。

| 票 | 档 | 依据(一行) | 阻断者/首动作 |
|---|---|---|---|
| 188 | T2 | 状态维已落字段＋门（`internal/tools/task.go` `TaskOutput.State` 只准填 D43 那 20 枚），但生产写入者仍 0（`cmd/wisp/run.go:1750` 的 `TaskBackfill{}` 不带 `State`），且「送进快照」那半明文归票 145 的载体射程 | **145**（按其 Progress log 末条：`145-r3` 交件后才派 `188-r3` 生产者那一跳） |
| 191 | T2 | AC#0 已由编排者按甲裁结（A379），但甲支要「跑一条命令并流式回显」靠的是**票 163 的常驻终端会话**做工具侧地基，且落地腿与 33 同包串行 | **163**（其自身再等 33：入向那一跳） |
| 192 | T2 | 两层（产物预览／联网浏览）都要在 `internal/panel/**` 新增载体，票面排程写死 `33-r1 → 191 → 192 预览 → 192 联网`，三枚串行 | **191**（其自身再等 33） |
| 193 | T3 | 三处都压在机主手上：AC#0＝他 09-25 亲自按住的 `tokens_fourway_test.go`（未答不许动）、AC#2 唯一那枚「真放开」＝队列 `Q-70`、动画三档**默认值由 owner 定**；界面那半（开关 UI）属 frontend | 机主（「动 tokens 那枚」／Q-70／默认档）＋frontend |
| 194 | T1 | owner 09-28 16:5x 已裁「方法名册按当前规格补，别搞债务了」（A381），代价表 `194-c1` 已交；现量 `internal/panel/bridge.go:42-45` 仍只有 4 枚、交集 0，堆1 那 20 枚里有 7 枚「有货只差名字」；不动 `PLAN.md`／不改 `docs/specs/**`，三条禁区（面板 allow／R20／config.set）具名排除 | 先落堆1：给 `bridge.go` **只追加** A 堆 7 枚 `panel.*` 常量（⛔ 不动 42-45 的顺序与命名），每枚带 AC#4 三枚常驻判据；范围 `internal/panel/**` |
| 195 | T1 | 起跑判据两条已满足（198-r1 早已交件、`panel_config_store.go` 7 处给了非测试调用方）；§1 真缺口现量成立——方向裁决只在 reload 侧（`internal/config/manager.go:323` `planLocked`），`settings.go` 的程序化写入口**没有任何方向判定**（`grep direction\|loosen` 于 `settings.go` ＝ 0 命中） | 在 `internal/config` 程序化写入口**锁外判方向**（放宽→记待决审批＋审计＋立即返回「未写入」；收紧→才进 `mu` 写）＋AC#2 两臂判据；范围 `internal/config/**`（⛔ 不碰面板 `config.set`） |
| 197 | T3 | 剩余 AC 无一格欠新 Go 产码（`197-a2` 现量：实体/状态/流/停止出口都已在树里），真缺的两样都不在本编队：界面那跳＝`frontend/src/lib/panel.ts` 的 `PanelSnapshot` 补 `tasks?:`（Go 侧 `internal/panel/composer.go:91` 已在发），以及五枚未勾格的**非实现者裁决料** | frontend（机主自己带给界面那支）＋非实现者裁决腿 |
| 200 | T3 | 实现已落（`internal/panel/composer.go:74` 发 `instructions` 一节，五状态 `loaded/none_found/off/refused/not_run`），红句逐字 `Go Snapshot emits [instructions] that interface PanelSnapshot does not declare` ⇒ 只能由界面支补一枚对象声明，Go 侧不许少发键/改尺/加豁免 | frontend（机主带过去；撤销口令「这一版先不显示说明清单」） |
| 201 | T2 | 七格里 AC#4／AC#7 已勾，剩余 AC#1／AC#2／AC#3a／AC#6b 票面自己写死「全部归口同一枚前置＝宿主原生入口（球／托盘）在 `wisp run` 里没接」 | **228**（Go 侧答复逻辑今天已全齐，缺的是宿主） |
| 211 | T1 | 甲已落：`internal/tools/subagent_197.go:85` `MaxConcurrentSubagents = 4` ＝ `internal/tools/bridge.go:24` `MaxToolConcurrency = 4`，钉也已在（`Test197SubagentPoolNeverExceedsBridgeCeiling`／`…PoolCapsAtBridgeCeiling`）；未做的只剩 AC#1/#5 的现跑读数与 AC#2 的正控 | 先现跑 AC#1（两常量并排打印）＋AC#5（`git diff` 证 `bridge.go`／`budgets.go` 那两枚 4 一字未动）＋AC#2 正控（人为把池写成 8 ⇒ 判据必红）；范围 `internal/tools/**` 只读＋测试。⚠ 乙（真 8）要机主一句话，不在本格 |
| 213 | T1 | 全仓零命令目录（`commandcatalog\|slashcommand\|CommandRegistry` 于 `internal`＋`cmd` ＝ 0 命中），落点已定＝复用 `internal/agent/control.go` 那枚本仓唯一「名字→本地执行、不走 LLM」的解析器；AC#4 明文给了一个**不依赖入向那一跳**的终态＝「名册可见＋命令行可执行」 | 先做那八条命令的**现量普查**（哪几条 Go 侧真做得了、哪几条要先补能力另立票）→ 再建可枚举名册＋真执行者；范围 `internal/agent/**`（`internal/panel` 只读） |
| 214 | T3 | 受理器 410 行今天生产零调用者，缺的是**装配根那一根线**，而 AC#1 明文要求「那枚调用者必须在 `cmd/wisp` 的装配路径上，不在测试里」；「选文件/@引用」两格由界面支定 | `cmd/wisp` 装配根（与 228/220 同包，整包死区、串行）＋frontend（选择器/@引用） |
| 215 | T3 | 拆成 215a 插件（配置 8 处命中、缺读者）／215b 技能（`grep -c Skill internal/config/schema.go` ＝ 0 ⇒ 要先补 **D36 那棵 section 树＝配置模型契约面，须人工批准**）；且清单要进快照＝新键会打红「加字段必两侧同批移动」那两把尺 | 机主（215b 的 D36 批准，只落 `A##`）＋frontend（新快照键同批声明） |
| 216 | T3 | AC#1/#2/#3/#4/#6 已勾；剩 AC#5 被非实现者裁「不适用」并归口 `SPEC-12 §5` 那条 S6 DEFERRED（`216-v3` 判语），且 `216-r2` 在 **C 卷台面**留了一枚推上去必红的用例（`name_claim_two_arm_216_test.go:105`，成因是 fail-closed 抢在 B 循环前、⛔ 为它改断言）⇒ 属 CI 色＋推送待机主 | 人裁 AC#5 归属（S6 子项／票 225）＋CI 台面色＋推送（只由编排者/机主） |
| 219 | T2 | 派单条件写死：`201（在飞）→ 223 → 224 → 才轮到三枚按钮`（223/224 已结，201 仍未闭），今天只可派「允许侧理由」那一小格；三枚按钮的「本次会话内」还归票 224 的 `GrantScopeSession` | **201**（其后仍有 220 的「该对哪一枚作答」） |
| 220 | T3 | 五格里 AC#1–AC#4 已勾，只剩 AC#5＝「本机无活实例、前台可控」下的 `./cmd/wisp` 整包终态名册——要请机主先退出他 10-15 19:03 自起的常驻 `wisp.exe`（PID 21892，`A988`），⛔ 任何腿不许代关 | 机主（Q-85 待机主窗口，默认不催） |
| 226 | T2 | 六格 5 勾（AC#1/2/3/4/6），AC#5「可撤销」缺的**答复动词**属票 219 的答复语法面（`Manager.SetAllowedDirs` 有实现有用例、生产调用者仍 0）；另「删行被补回」残余已转票 227（已结） | **219**（撤销动词那一格） |
| 228 | T3 | 全票要的是「球＋托盘活在同一条**常驻 GUI 进程**里」：托盘菜单要真机 `winlive`（`228-r3b` 现量：本机开不出可注入点选的托盘世界）、球绘制在 `ui-sta` 线程、`D38(e)` 十步收口；且落地全程在 `cmd/wisp` 同包 | 真机桌面（机主在场）＋`cmd/wisp` 死区（其 L1 那半还要先等 220 定枚举口） |
| 230 | T1 | AC#1（票 115 裁决表）与 AC#2（票 97 注释改成实话）均已由非实现者裁决后翻勾；剩 AC#3 只动**工单文件**（乙支：把票 80 那两格转进本票＋具名作废「接手方＝票 21」那句），AC#4 票面明文「本票不接」只负责摆给机主 | 先落 AC#3 乙支（只改 `.scratch/wisp/issues/**`，不碰码、不动 `-done` 后缀）；AC#5 那发 `./cmd/wisp` 整包终态要等安静车道 |
| 238 | T3 | 机主 09-30 16:5x 已裁「本条不做」（`A476`，原话「锁个屁…」）：甲（`internal/winsec` 新增非递归封根）不批、乙（全树走查）也不做，四枚框全留未勾、不再排落地腿 | **已撤**（withdrawn；撤销口令「`Q-72` 重开」，重开只能他主动说） |

T1-count=5
T1: 194, 195, 211, 213, 230
