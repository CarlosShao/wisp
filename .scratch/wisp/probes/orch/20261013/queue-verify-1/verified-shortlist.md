# queue-verify-1 · 已验面可派单（把普查"面标"逐票读真；2026-10-13 22:1x 只读复核）

**腿**：只读复核 `queue-verify-1`｜写面＝本目录 `probes/orch/20261013/queue-verify-1/**`｜⛔ 零 go 命令 · 零产码 · 零 push
**源**：`.scratch/wisp/issues/*.md` 现盘（工作树）＋普查件 `queue-census-1/queue-ready.md`（20:5x–21:2x 收口）
**时点**：2026-10-13 22:1x–23:0x +08；分支 `dev`；未 commit 任何仓内字节
**读数头**：验＝**40 枚**（★近 4／票面有写面行的 23／写面未声明但活跃 10／老存根抽 3）；按派单跳过 `-done`／在飞／明确 blocked-RESERVED

## 0. 尺的射程（口径写死；引用以本节为准）

- **每票读到**：文件名、`wc -l`、头 12–60 行、`grep -n` 的 Status／`写面|只许|只动|只改|不碰|禁区`／`前置|串行|按住|blocked|RESERVED` 行、`^- [2026` 末 bullet、勾框计数（`- [x]`/`- [ ]`）；★近 4 枚＋疑点 12 枚另读整票或定点区间（§4 列名）。
- **真写面**＝票面「写面／禁区／排程」三处**明写的路径**，不是关键词命中；🔶＝票面只给候选或"二选一必须先裁"。
- **面标**只用派单六词：`cmd/wisp`／`internal/panel`／`internal/tools`／`internal/其它`／`docs/evidence`／`probes`（可多标）。
- **没读到**：① 302 枚里只碰 61（40＋无 Status 行组抽检 21 枚头 3 行）② 114 枚 `-done` 不读 ③ 老存根 01–62 只抽 3 枚，其余 40 枚"是否仍在射程"**仍未读** ④ 票内一切行号/指针**一律未复跑**（零 go 命令）⑤ 各票证据件 `docs/evidence/**`、`probes/**` 内容未打开。
- ✎ 凡"已 done／已交件"的凭据＝**同号非 `-done` 文件不存在** 或 票面 Progress log 逐字；未做 git 历史之外的验证。

## 1. 可立即派（按面分组；非面板面在前）

### A. `internal/tools`（现刻刚空出；⚠ 同包六枚互斥，须串行：173／174／237／287／175 只裁腿）
| 票 | 真写面（票面明写） | 面标 | 前置判定 | 末bullet | 判语＋理由 |
|---|---|---|---|---|---|
| 173 | `internal/tools/bridge.go` 那一味＋用例 | internal/tools | ①同包无在飞（现刻空）②`probes/154`、`161/r6` 无人改时跑尺 ③gofumpt 干净 ⇒ 基本满足 | 09-27 20:3x | **可立即派**：写面窄，前置仅剩环境项 |
| 237 | **只** `internal/tools/*_test.go`（判定分支、`bridge.go` 一条不许动） | internal/tools | 没写前置（来源 `235-v1` 已交）⇒ 视为无 | 无 bullet（09-29 立） | **可立即派**：零产码，`-race` 红推断定性 |
| 287 | **只许** `internal/tools/**`（夹具＋用例） | internal/tools | AC#2 方法学前置＝票 282（已 done）⇒ 满足 | 10-09 12:2x | **可立即派**：AC#1 已翻，余三格 |
| 174 | `internal/tools/**`＋`internal/agent/spill.go` 回执文案侧；AC#2b 一行 `cmd/wisp/run.go` | tools＋cmd/wisp | 派单前置＝**无**（测量腿现在就能跑）；AC#2c 等票 177 裁；甲/乙两形不许并发 | 无 bullet（AC#2b 09-27 21:4x） | **可立即派**（测量腿）；接线那行与 177 串行 |
| 175 | 第一腿**只裁不改**；落地面🔶`internal/risk/provenance.go` 名册 **或** `internal/tools/bridge.go` mark 判定形 | internal/tools＋internal/其它 | 票 164 `164-v1` 交件 ✓（`docs/evidence/s1/164-task-output-accept-r1.md` 在盘） | 09-27 20:1x | **可立即派**（只裁那一腿）；两落地面都超已批射程⇒落地需人批 |

### B. `internal/其它`（agent）
| 票 | 真写面 | 面标 | 前置判定 | 末bullet | 判语＋理由 |
|---|---|---|---|---|---|
| 179 | `internal/agent/loop.go`（只 `decideRisk` 支＋注释）＋新 `*_test.go`；两行注释在 `cmd/wisp/run.go`、`internal/tools/bridge.go` | internal/其它（＋注释面） | AC#8 已归票 183；AC#5 禁改 `PassThroughUnclassifiedRisk` 默认值 | 无 bullet（09-28 12:2x 收表） | **可立即派**：残 AC#1（零产码 overlay 复跑）＋AC#9＋AC#10 |
| 139 | `internal/agent/compress.go` 成功边＋AC#3 用例（实现已由 `139-r1` 落盘，待非实现者裁） | internal/其它 | 面板可视化退 owner；`DEFERRED(D28-1)` 登记面禁改；AC#1 判据须先改写 | 09-26 09:0x | **可立即派**（收口/验收腿） |

### C. `probes/**`＋CI 新增步骤（无产码面）
| 票 | 真写面 | 面标 | 前置判定 | 末bullet | 判语＋理由 |
|---|---|---|---|---|---|
| 171 | `probes/**` 尺本体＋`.github/workflows/ci.yml` **新增步骤**（动触发表＝契约级） | probes＋CI | ①`internal/tools` 两程交完 ✓（现刻空）②台件 gofumpt 干净 | 09-27 19:1x | **可立即派** |
| 172 | 一枚新尺＋用例（`probes/**`） | probes | ①`internal/tools` 无在飞 ✓②`PLAN.md`/specs 一字不许改 | 09-27 18:2x | **可立即派** |
| 182 | 零产码：只写票面＋证据件＋台件目录（AC#6 写死） | docs/evidence＋probes | 无（AC#7 明写不阻塞 179 链） | 10-08 10:3x（收口节） | **可立即派**：全票只读普查 |
| 236 | `probes/**`＋`internal/tools/*_test.go`＋`ci.yml` 注释＋`scripts/portable-tests.sh` | tools＋probes＋CI | AC#3 前半已交（`236-a4`）；AC#5/6 要**整包窗口** | 10-06 20:0x | **可立即派（限 AC#3 后半／AC#4）**；AC#5/6 需前置 |

### D. `cmd/wisp`（现刻刚空出）
| 票 | 真写面 | 面标 | 前置判定 | 末bullet | 判语＋理由 |
|---|---|---|---|---|---|
| 232 | `cmd/wisp/config_reload*_test.go`（仅测试） | cmd/wisp | 与 231 同文件面 ⇒ **231 已 done**（无同号非 done 件）⇒ 串行约束消解 | 10-06 19:3x | **可立即派**：test-only |
| 288 | **只** `cmd/wisp/panel_inbound_guards_35r3_test.go` 一处对齐（2 行） | cmd/wisp | 无；票面自陈"单独排不值"⇒ 建议随批捎带 | 无 bullet（10-09 立） | **可立即派（捎带）** |
| 233 | `internal/config`（档位标到键）＋`cmd/wisp`（硬编码清单） | cmd/wisp＋internal/config | Blocked by `223-v2` 已交 ⇒ 无阻塞；⚠ 与 232 同撞 config 面 ⇒ 串行 | 09-29 17:5x | **可立即派**（与 232 串行） |
| 180 | `cmd/wisp/panel_host_windows.go`（`:304-305` 尺寸字面量接读者）＋`internal/config`＋快照名册 | cmd/wisp＋internal/config | 09-28"排在票 33 宿主之后"**已被 10-02 现量推翻**（前件满足）；AC#3 禁"改默认值交差" | 10-08 追加节 | **可立即派**；⚠ 与 config_reload 族（232/233）串行 |
| 272 | `cmd/wisp/config_reload*_test.go` 窗口下界＋门禁 | cmd/wisp | 5 勾/2 未勾；AC#7 票面写"本票不选形" | 10-07（编排者翻框节） | **可立即派（限 AC#5 门禁四数）**；AC#7 待裁形 |
| 181 | `internal/panel`＋`cmd/wisp/run.go`（只读那一半已落地） | panel＋cmd/wisp | 残 AC#7＝**非实现者行为突变终裁 `181-v3`**（今天只由实现者跑过） | 无 bullet（09-28 17:x） | **可立即派（验收腿）**：按 SPEC-10 判"勾"必须非实现者 |

### E. `internal/panel`（⚠ 同族互斥：186／187／189／190 都可能碰 `bridge.go` 读面/白名单 ⇒ 编排者须排串行）
| 票 | 真写面 | 面标 | 前置判定 | 末bullet | 判语＋理由 |
|---|---|---|---|---|---|
| 186 | Go 读面＋请求腿＋快照字段；`internal/panel/workspace.go`（**只追加**）＋`docs/evidence/s1/186-panel-ask.md`；⛔frontend、⛔`internal/tools` 执行面 | **internal/panel**＋docs/evidence（普查标"其它"＝误） | 无硬前置；`181-c1`/`186-a1` 普查已交 | 10-13 18:0x ★近 | **可立即派**（Go 侧）；与 187 同碰 `bridge.go` 白名单⇒串行 |
| 187 | `internal/panel`（`bridge.go` 白名单＋handler）＋契约变更入档 | internal/panel | 与 186 分程、但同碰白名单⇒串行；若新增方法名按票 194/台账先落批准记录 | 无 bullet（09-28 立） | **可立即派**（串行） |
| 189 | `internal/panel/git.go`（**复用票 181 读面，不另起实现**） | internal/panel | 设计核已两轮（`189-a1`）；雷区＝提交/推送不许进面板 | 10-08 10:4x | **可立即派** |
| 190 | `internal/panel` 读面（文件树源） | internal/panel | 原"按住等 `145-r2`"⇒ **145-r2 已交（09-28 18:0x）** ⇒ 解 | 10-08（第二程） | **可立即派** |
| 194 | `internal/panel/bridge.go`＋名册补齐（owner 09-28 点名批准，`A381`） | internal/panel | AC#5：生产调用者仍 0 ⇒ 不许拿"名册补齐"冒充接上；AC#1 只读表先行 | 10-06（更正节） | **可立即派（限 AC#1 逐枚代价表）**；落地格需票 33 的窗 |
| 220 | AC#1：`internal/panel/subagent_roster_197.go` 双查；AC#2🔶`internal/agent/approval`（甲＝新方法名／乙＝L1 进 LiveApprovals） | internal/panel＋internal/其它 | 禁区＝不新增方法名（甲支须先落台账批准记录）；禁改 D43 名 | 无 bullet（09-29 立） | **可立即派（限 AC#1／AC#3）**；AC#2 需人批 |
| 35 | `cmd/wisp/panel_host_windows.go`＋一枚新 `cmd/wisp/*_test.go`（普查腿定：**只 `:42` 前半格＋`:51` 残差**两枚纯测试格） | cmd/wisp＋internal/panel | `:44/:47/:49` 共用前置＝出向那一跳（被票 97 归属、`frontend/**` 只读、票 194 名册挡） | 10-08 10:2x | **可立即派（限两枚纯测试格、零产码）**；其余需前置 |

## 2. 需前置 / 需人批 / 建议不派（逐枚判语）

| 票 | 类 | 判语＋一句理由 |
|---|---|---|
| 314 | 需前置 | 写面＝`probes/314/**`＋必要时 `cmd/wisp/*_test.go`（夹具 profile 名）；⛔ 与 294 任何在飞整包腿**不同波**，AC#1/#2 要静桌面＋整包色 ⇒ 等 294 收口＋编排者约静窗 |
| 248 | 需前置 | 写面＝`internal/panel/bridge.go`＋`cmd/wisp/panel_inbound.go`；顺序钉死：`246-r2`→票 33 宿主段（动构建链）→本票，AC#8 另有两道外部前置（Go→页返回值、请求↔回执配对） |
| 33 | 需前置 | 写面＝`cmd/wisp` 宿主＋`internal/panel`＋构建链；1 勾/13 未勾，AC#12/#13 待非实现者裁，且"用户仍看到空壳页"的根因（页面按 URL 取自己 .js/.css 无注册方）**要另立新票** |
| 275 | 需前置 | 写面＝`attrib.sh`（现盘在 `probes/161/**`，属票 161 仪器）＋CI 输入集；动它需**具名解冻 `A##`**；AC#5 卡"待推送取数"（机主已说不推） |
| 244 | 需前置 | 写面＝构建链（`scripts/build.ps1` 那族）；AC#0 未答前整票不许动构建链；顺序＝**33 → 248（到"能举一发任务"）→ 244 AC#1** |
| 188 | 需前置 | 写面＝`internal/tools/task.go`＋快照字段；**硬前置＝票 33／35（面板宿主与桥）**，票面自己写死 |
| 246 | 建议不派 | 8 勾/1 未勾；残 AC#8（boot 时序 flake）**地界属票 228**，本票只做"归口＋判据"、⛔ 不在本票改 ⇒ 无产码活 |
| 247 | 建议不派 | 8 勾/2 未勾：AC#2 判"不成立"**已归票 297**；AC#6 欠"设备被占"自造形＋"无权限/禁用"要**机主一句话** ⇒ 非编队可解 |
| 238 | 建议不派 | owner 09-30 16:5x 已裁"这条不做、两支都不落、⛔ 不许再摆"（`A476`）⇒ 票面待收口，非待派 |
| 15 | 建议不派 | 老存根（S0 时代英文骨架、无写面行、Progress log 空）；AC 引 `docs/SLO.md` 等仓里不存在的交付物 ⇒ 语义未核，需人重估后重立 |
| 16 | 建议不派 | 老存根（S2 验收骨架，同 15 形）⇒ 同上 |
| 21 | 建议不派 | `Blocked by: 17`（17 已 done）；segment-1 landed、segment-2 未认领，但 bridge/gate 族已由后续票（146/197/201/219/246）长出来 ⇒ 残余语义需人重核 |
| 160 | 需人批 | 契约级、票面写"仍是**等 owner 点头**"；前置②已改判为"票 158 五格成立＋AC#5 残余有票号（169）"，158 仍非 done |
| 193 | 需人批 | AC#0＝`internal/panel/tokens_fourway_test.go` 放开与否——owner 09-25 **亲自按过**，只有他能放；其余格（Go 侧字段＋门禁清单）都建在这枚裁定之后 |
| 196 | 需人批 | 三形（甲映射表／乙 CHECK 约束触 `D35`／丙按住）**待 owner 拍**；票面 15 行、零框 |
| 239 | 需人批 | 4 勾/1 未勾：AC#4 里"窗口 72→100、贴边 22px"是可见行为变化，**要 owner 单独裁一句**；另残留 0.25px 富余与 `cmd/wisp` 全量未复跑 |

## 3. 抽核 3 枚复核记录（对普查自动命中逐项现读）

1. **314**（普查：无 Status 行→人工改 A ✎；[面:cmd/wisp]；末 10-13 17:22:56）。现读：头 3 行确无 `Status` 行；写面行确为 `:39`（逐字同）；**修正**＝票尾已多一条 21:23 的 bullet（`294-r1b` 机主在场对照 5/5 死），末 bullet 应为 **21:23**；面标正确（写面含 `cmd/wisp/*_test.go`），但真正的"跑"要整包窗口＋静桌面 ⇒ 普查的"可派"偏乐观。
2. **186**（普查：[面:其它]、★近 10-13 18:0x、x0/e10）。现读：★近属实（末 bullet 10-13 18:0x `186-a1`）；**面标纠错**＝真写面含 `internal/panel/workspace.go`（只追加）＋`docs/evidence/s1/` ⇒ 应标 `internal/panel`，不是"其它"；勾数 e9 未勾（普查 e10 差 1，我数 `- [ ]`＝9）。
3. **232**（普查：[面:其它] S:cmd/wisp、★近 10-06、x5/e3）。现读：真写面＝`cmd/wisp/config_reload*_test.go` **仅测试**（禁区逐字"不改产品文案/不新增导出名/frontend 零读零写"）；勾数 5 勾/**1 未勾**（AC#4）；**前置消解**＝"与票 231 同文件面"的 231 已 done（无同号非 done 件）⇒ 普查"面冲突"不再成立。

（另二枚顺带复核：**220** 普查 [面:其它] ⇒ 真面含 `internal/panel/subagent_roster_197.go`＋`internal/agent/approval`，也是面标漏报；**174** 普查 [面:tools] ⇒ 真面还含 `internal/agent/spill.go`＋`cmd/wisp/run.go` 一行。）

## 4. 交件自报：没读到的行 / 是猜的

- **未复跑任何行号/尺/测试**（本腿零 go 命令）：票内 `:LNN` 指针、"零调用者/零命中"这类断言一律按票面原话引用，**未经我复算**。
- **只抽不读**：老存根 01–62 段只抽 15／16／21 三枚，其余 40 枚"是否仍在射程"**没读到**；无 Status 行组 21 枚只读头 3 行（确认全是 10-02–10-13 活跃票；312／313 票面明写"只立不派"）⇒ 整组**未逐票读**，要点将请另派。
- **未打开**：任何 `docs/evidence/s1/**`、`.scratch/wisp/probes/**` 证据件正文（§1/§2 里引它们只作"文件在盘/票面自陈"级凭据）。
- **"已 done／已交件"** 的判定＝同号非 `-done` 文件不存在 ＋ 票面 Progress log；未跑 `git log` 逐枚核对交件 commit。
- **猜的**：🔶 标记的候选写面（175 两形、220 AC#2 甲乙、187 AC#2 通道）＝票面给的候选集，未裁前不当作真写面；14 枚票的"末 bullet"栏写"无 bullet（立票日）"＝票内无 `^- [2026` 行，日期取自票面首段立票时刻，**非逐字 bullet**。
