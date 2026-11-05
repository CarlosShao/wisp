# retriage-186 — 开放票池**再分诊**（A–F 六桶）＋下一波可派 A 清单（只读普查腿）

- **腿**：`retriage-186`（只读普查腿，第 186 轮波后重分诊）。**取数时刻**：2026-11-05 17:0x–17:2x +08。
- **锚**：`HEAD = 248041c2e`（`git rev-parse --short=9 HEAD` 本程现读），分支 `dev`。
- **射程**：`.scratch/wisp/issues/*.md`，排除 `-done` 与 `README.md` ⇒ **116 枚**
  （尺：`ls .scratch/wisp/issues/*.md | grep -v -- -done | grep -v README | wc -l` ⇒ **116**；同一把尺含 `README.md` 时＝117）。
- **只读自证**：本程**零改**票面／台账／代码／脚本／契约；全仓唯一写入＝本文件。
  **零 `git add` / `git commit` / `git push`**（本件落盘未提交，是否提交由编排者定）。
- **尺（逐条可复跑，引用前先重跑）**：
  1. **未勾 AC** ＝ `grep -c '^- \[ \]' <票>`；**已勾** ＝ `grep -c '^- \[x\]' <票>`。
     ⚠ 只算**列首（顶层）**框：缩进框、票正文里的示例框都不计 ⇒ 本表数的是"顶层判据格"，不是"总框数"。
  2. **Status** ＝ 票面第一条含 `Status` 的行（多行时取**最新那条只追加行**的实质判语；只摘不解释）。
  3. **前置** ＝ 票面 `Blocked-by` 点名的票号，逐枚
     `ls .scratch/wisp/issues/<N>-*.md` 现读：命中 `-done` ⇒ **已闭**；命中的仍无 `-done` ⇒ **开**。
     `Blocked-by: nothing` 或只写"无"⇒ **无前置**。
  4. **面（Packages）** ＝ 票面 `Packages` 字段；票面没有该字段的（新票多数没有）直接用 AC 里点名的目录（本表不加"推断"标记，凡是票面给了 `Packages` 的照抄）。
  5. **桶判定优先级**（一票同时命中多桶时，按硬门先后来）：
     **W（撤票）＞ C（人）＞ B（前端）＞ D（真机）＞ E（CI）＞ F（前置票）＞ A（可派）**。

---

## 1. 一句话答案

| 桶 | 枚 | 含义（本程用法） |
|---|---|---|
| **A** | **31** | **本编队现在可派的 Go／测试／文档活**：无 owner 门、无 `frontend/**`／`design/**` 写面、判据不要真机、前置票已闭 |
| B | 14 | 写面在 `frontend/**`／`design/**`，归机主外部 agent（本编队结构性写不到） |
| C | 22 | 待人：改契约文字／花钱或凭据／机主一句话／机主亲拍形 |
| D | 15 | 判据只能在**真机**取：真麦、真扬声器、真托盘、真窗口／WebView2、桌面视觉签收 |
| E | 8 | 结案判据**就是 CI 颜色**（今天机主"暂时不推远程"⇒ 取不到） |
| F | 22 | ⛔ 现在不许开工：前置**票**未闭／自身触发门未满足／编排者明确按住 |
| W | 4 | 撤票／VOID／已并案——**不是派发类**，单列只为不被当活票 |
| 合计 | **116** | — |

**A 的 31 枚内部再分**：

| 子档 | 枚 | 票号 |
|---|---|---|
| ① 在飞／盘上有未提交件（**本轮不再派**，等收件） | 5 | `23` `194` `197` `200` `211` |
| ② 已落盘、只差非实现者翻框（不是写腿） | 1 | `284` |
| ③ ⚑ 票面前提已被盘上读数顶倒（该做的是**作废/改票**，不是照着派腿） | 2 | `176` `180` |
| ④ **干净可派**（本轮真正的候选池） | **23** | 见 §4.1 |

⇒ **A-count ＝ 31；干净可派 ＝ 23。下一波建议派 7 枚（§5），其中若只批一枚＝`98`。**

---

## 2. 六桶的判据（这张表就是分类器的源码）

| 桶 | 判据（必须同时成立） | 典型否决项 |
|---|---|---|
| **A** | ① 活本体是 Go／测试／脚本／文档；② `frontend/**`、`design/**` 零写面；③ 判据能在本机取（测试、脚本、`grep`、定向 `-run`）；④ `Blocked-by` 里点名的票**全部**已 `-done`（或本来就无前置）；⑤ 不欠 owner 一句话／不欠契约文字批准 | 缺一件即出 A |
| **B** | 写面或**结票判据**在 `frontend/**`／`design/**`（面板、设置页、菜单、球视觉的"看得见"那一半） | 即使 Go 半边可做，也不进 A——**结不了票** |
| **C** | 需要**人**：改 `C1–C32`／`D1–D47` 文字、花钱、外部账号、机主亲拍形、机主取舍 | 票面写"归机主／owner"或"人工批准"即归此桶 |
| **D** | 判据要求真麦克风、真扬声器、真托盘、真窗口、真桌面、第二块屏、另一台机器（含 macOS） | 本机无窗/无麦时**不可判** |
| **E** | 结票判据是"某 CI run 的某 step 是 success"或 runner 侧逐字读数 | 本机读数**不许**冒充（票 276 逐字禁） |
| **F** | 前置**票**未闭（票面 `Blocked-by` 里仍有开放票）／自身触发门（`AC#0`）未满足／编排者明确写"按住、默认不排" | 前置闭了就**自动**离开 F |
| **W** | `Status` 含 `WITHDRAWN`／`withdrawn`／`VOID`／已并案 | 文件按规矩**不删不改名**，就地留档 |

> ⚑ 记号：票面 `Status` 的**立案前提**已被本程或同日 `triage-180` 的现量顶倒（§6）。桶位照"今天还能做什么"给，不照票面自称给。

---

## 3. 全表（116 枚；`未/已` ＝ 顶层未勾／已勾 AC 数）

| 号 | 桶 | Status（票面摘） | 前置票（现读） | 面（Packages／写面） | 未/已 | 备注（判定依据一句话） |
|---|---|---|---|---|---|---|
| 21 | D | segment-1-landed；段 2 未做 | `17` 已闭 | `internal/agent/approval`＋`cmd/wisp`＋球视觉 | 5/0 | 段 2 含"球视觉＝机主签收"，判据要真桌面 |
| 23 | A | ready-for-agent | `20` 已闭 | `internal/tools`（window／screen／input） | 5/0 | **在飞**：盘上未跟踪 `internal/tools/ticket23_image_channel_test.go`（C7 image 半边） |
| 25 | F | ready-for-agent | `21`·`23` **开** | 真机脚本面 | 5/0 | 两枚前置都还开着 |
| 32 | D | ready-for-agent | `27`/`29`/`30`/`31` 全闭 | 真机语音面（③听写④陪伴） | 5/0 | 前置已清，但判据要真麦真声 |
| 33 | D | ready-for-agent（07/12 已闭） | 无 | `cmd/wisp`＋`internal/panel`（WebView2） | 13/1 | C27 宿主＝真窗；且它是 `35`/`191`/`192`/`195`/`213`/`214`/`215` 的公共前置 |
| 34 | B | ready-for-agent | `01` 已闭 | `frontend/**` | 5/0 | 纯前端脚手架 |
| 35 | F | ready-for-agent | `33`·`34` **开** | `internal/panel`（C17 桥） | 6/4 | 前置两枚全开 |
| 36 | B | ready-for-agent | `35` **开** | `frontend/**` | 6/0 | 结票判据在屏上（前置亦开，按优先级归 B） |
| 37 | B | ready-for-agent | `21`·`35` **开** | `frontend/**` | 6/0 | 同上 |
| 38 | B | ready-for-agent | `35` **开** | `frontend/**`＋Go `task.list` | 6/0 | 同上（Go 半边=工具名册，结票仍在屏上） |
| 39 | B | ready-for-agent | `35` **开**（`05` 已闭） | `frontend/**` | 5/0 | 同上 |
| 40 | B | ready-for-agent | `35` **开**（`29`/`44` 已闭） | `frontend/**` | 8/0 | 同上 |
| 46 | D | ready-for-agent | `41`–`45` 全闭 | 真机（空闲/KWS/电源预算） | 5/0 | 前置已清；判据要真机场景 |
| 47 | F | ready-for-agent | `21` **开** | `internal/agent/scheduler` | 6/0 | 前置开（`10` 已闭） |
| 48 | F | ready-for-agent | `37`·`47` **开** | `internal/agent/approval`（C18） | 5/0 | 链式前置 |
| 49 | F | ready-for-agent | `48` **开** | `internal/agent/approval`＋`perm` | 6/0 | 链式前置（`04` 已闭） |
| 51 | C | ready-for-agent | `17`/`20` 已闭 | `internal/plugin`（C24 宿主 API） | 6/0 | `C24` 初始集属**未定案契约**（`AGENTS §2`） |
| 53 | C | ready-for-agent | `48`·`51` **开** | 真服务凭据（lark-cli） | 4/0 | 花钱/凭据（且前置也开） |
| 54 | F | ready-for-agent | `47`/`48`/`49`/`53` **开** | 真机＋CI | 5/0 | 前置全开 |
| 55 | C | ready-for-agent | `54` **开**＋owner go-ahead | macOS（无该平台机器） | 5/0 | 要机主点头＋无硬件 |
| 56 | C | ready-for-agent | `54` **开**＋SignPath＋go-ahead | 外部账号/花钱 | 5/0 | 证书与发布账号＝花钱 |
| 57 | F | ready-for-agent | `54`/`53` **开** | `docs/**`＋`internal/plugin` | 3/0 | 前置开 |
| 58 | C | ready-for-agent | `54` **开**＋go-ahead＋选型 | 英文模型下载＋前端 | 4/0 | 花钱＋选型归人 |
| 60 | C | ready-for-agent | 前置票全闭，但门槛②未满足 | 服务商实时 key（花钱） | 7/0 | 票面自写"缺 key 则 AUTO-DEFER" |
| 62 | D | review（**从未被验收**） | 无 | `internal/ball`＋真桌面 | 5/3 | 玻璃质感＝人眼签收 |
| 64 | D | review | `07` 已闭 | `internal/ball`＋第二块屏 | 1/8 | 剩 4 框＝热键/多显示器实拖 |
| 65 | C | blocked-on-owner | `62` **开** | `design/refs`＋`internal/ball` | 6/0 | 缺**参考图**＝owner |
| 68 | C | review（AC#2/#3 blocked-on-owner） | 无 | `internal/ball`＋桌面 | 2/2 | 剩两格归 owner 裁 |
| 70 | C | blocked-on-owner | `66`/`68` 已解除 | `.github/workflows` | 2/4 | 剩的判据"不在代理手里" |
| 71 | E | in progress | `70` **开**（语义无关） | `cmd/balldebug`＋`d22scan`＋`ci.yml` | 6/1 | AC#5＝一次**真红**的 CI 读数 |
| 72 | E | review | 无 | `internal/risk`＋`cmd/wisp` | 2/4 | AC#4＝runner 侧逐字读数 |
| 77 | B | blocked（owner 收走前端） | 已解除 | `frontend/**` | 4/3 | 前端归外部 agent |
| 85 | E | open（排队） | 无 | `ci.yml`＋`staticcheck` 钉版 | 5/1 | AC#1/#2 要真 run 日志行 |
| **98** | **A** | open | 无 | `scripts/**`＋`docs` | 5/0 | 零共享编译面；`cmd/wisp` 取证口令（§6.3） |
| 100 | F | **RESERVED** | 自身 `AC#0` 触发门未满足 | `internal/proc`＋真容器 | 7/0 | 触发门＝不许开工 |
| **103** | **A** | rejected-needs-fix（验收退回） | 无 | `internal/winsec`＋`internal/memory` | 4/0 | 缺陷形状已被探针定死，四格全空 |
| 112 | E | ready-for-review | 无 | `.github`＋`internal/winsec` | 5/0 | 只在 CI Windows runner 上红 |
| 122 | F | open | `85` **开** | 全仓 Go（staticcheck 清账） | 7/0 | 票面自写"85a 的下游" |
| 123 | A | open | `117` 已闭 | `cmd/wisp/*_test.go` | 6/0 | 本机可改可跑（AC#4 有 E 尾巴） |
| 135 | A | ready-for-agent | `130` 已闭 | `cmd/wisp`＋`internal/proc` 测试 | 7/1 | 自证腿＋变异，全在本机 |
| 136 | A | ready-for-agent | 无 | `internal/observe`＋`scripts` | 4/11 | 钉与正控本机可跑（AC#6 有 E 尾巴） |
| 138 | C | blocked（`Q-55`） | 无 | `internal/proc` | 2/2 | 等一次人工批准一起落 |
| 140 | A | open（未派） | 无 | `ci.yml`＋四包 | 4/0 | 先量面再裁形，全本机；改 `ci.yml` 要登记 |
| 145 | A | review／部分结案（只一段判据齐） | 无 | `internal/panel` 快照载体 | 3/4 | 续腿可派，但只一段可落 |
| 148 | C | blocked（双批准） | 无 | `internal/agent/approval` | 5/0 | 两枚前置都是人工批准 |
| 150 | C | blocked（`Q-55`） | 无 | `docs/PLAN`＋扫描器 | 5/0 | 三条出路都动冻结面/契约句 |
| 163 | C | 立而不派（更正行称门已开） | 无 | `internal/tools`（D34/D46） | 5/1 | 扩 `shell.exec` 语义＝契约；开工前须编排者确认"门真开" |
| 167 | F | ready-for-agent·立而不派 | `145` **开** | `internal/panel`＋`frontend` | 6/1 | 与 145 抢同一批快照字段 |
| 168 | F | 立而不派 | `51` **开** | `internal/plugin`（C17/C24 白名单） | 7/0 | 插件地基不存在 |
| 169 | B | owner/前端会话接手件 | 无 | `frontend/src/.../right-rail.tsx` | 5/0 | 本编队结构性写不到；**永不派** |
| 170 | C | 立而不派 | 无 | `internal/agent/approval` | 6/0 | "CLI 要不要停下来等答案"＝产品口径 |
| 174 | C | ready-for-agent（余 2 格） | 无 | `internal/tools`＋`PLAN.md:431` | 2/7 | 7 格已勾；剩 `AC#4` 契约轴 |
| 176 | A ⚑ | 立而不派（排 `175` 后；`175` 已闭） | 无 | `cmd/wisp`＋`internal/tools` | 4/1 | **`AC#1` 停手条件已触发**：`RunAsync` 产码 2 处、`Roster.Record` 4 处 ⇒ 该**作废/改票** |
| 177 | A | ready（`Q-61` 已批＝甲，核腿已交） | 无 | `internal/tools`＋`internal/risk` | 3/5 | 写手腿前置已清 |
| 180 | A ⚑ | ready-for-agent | 无 | `internal/config`＋`cmd/wisp` | 5/1 | **"生产零读者"已不成立**：`panel_resident_windows.go:208` 真读 `cfg.Panel.Width` ⇒ 剩的是补落地＋反控 |
| 182 | A | ready-for-agent（只读） | 无 | 只读普查（零产码） | 6/1 | 10 栈逐堆定源＋归口票 |
| 185 | A | 立、按住（等 `183-v1`；`183` 已闭） | 无 | `internal/tools`＋`internal/risk` | 7/0 | 前置已消，全测试面 |
| 187 | B | ready-for-agent | 无 | 面板写入口（契约）＋`frontend` | 9/0 | 屏上那一半归前端 |
| 188 | A | ready-for-agent | 无 | `internal/tools`＋`internal/panel` | 5/1 | owner 原话已定（子代理进产品），只补 Go 侧源 |
| 191 | F | ready-for-agent | `33` **开** | `internal/tools`＋`internal/panel` | 6/1 | 票面"排在 `33-r1` 之后" |
| 192 | B | ready-for-agent | `33` **开** | `internal/panel`＋`frontend` | 5/1 | 预览/浏览器＝屏上那一半归前端（前置 33 亦开，按优先级归 B） |
| 193 | B | ready-for-agent | 无 | 前端门禁＋`config` | 7/1 | 放开的是**前端动画**与卡前端的门 |
| 194 | A | ready-for-agent（**已派**） | 无 | `internal/panel/bridge.go` | 7/0 | 在飞：等收件 |
| 195 | F | blocked（等 `33` 片 B） | `33` **开** | `internal/config`＋`cmd/wisp` | 7/0 | 盘上已有 `195-r1` 方向裁决件＋`internal/secret/**` 在飞 ⇒ **不是干净面** |
| 197 | A | **在派** | 无 | `internal/tools`＋`internal/panel` | 5/2 | 在飞：等收件 |
| 200 | A | **在派** | 无 | `internal/projctx`＋`agent`＋`panel` | 8/0 | 在飞：等收件 |
| 201 | F | 部分结案 | `228`·`246` **开** | 球/托盘原生面 | 5/2 | 剩四格归口"宿主原生入口"＝228/246 |
| 211 | A | **在派队列里** | 无 | `bridge.go`＋`budgets.go` | 5/0 | 在飞：等收件 |
| 213 | F | 待派 | `33` **开**（前置＝入向线） | `cmd/wisp/commands.go` | 5/0 | 三票共用那根入向线 |
| 214 | F | 待派 | `33` **开** | `internal/tools`（附件受理） | 6/0 | 同上 |
| 215 | F | 待派 | `33` **开** | `internal/plugin`＋装配 | 5/0 | 同上 |
| 216 | A | 待派 | 无 | 显示层纯函数（Go） | 1/5 | 5 格已勾，只剩 `AC#5` 一格 |
| 219 | F | 接单＝做 | `201` **开**（→`223`/`224`） | 原生卡面＋审计 | 6/0 | 票面派单条件写死"排在 201→223→224 之后" |
| 220 | A | 待派 | 无 | `internal/panel` 名册 | 1/4 | 只剩 `AC#5` 整包读数（受 314 死区影响） |
| 226 | F | 产码已交·残余转 `219` | `219` **开** | `cmd/wisp`＋`internal/config` | 1/5 | 剩一格已明确转票 |
| 228 | A | 09-29 部分勾（方向＝甲已裁） | 无 | `cmd/wisp` 常驻＋`internal/proc` | 11/5 | 不是小活；验收要真窗（D 尾巴） |
| 230 | F | 待派 | `21` **开**（`AC#3` 要写进 21） | 票面/证据（零产码） | 3/2 | 且 `AC#4` 明文"必须由人点"⇒双重 |
| 238 | A | 待派 | 无 | `internal/proc`＋`internal/winsec` | 4/0 | 封根归属＋常驻尺，全本机 |
| 239 | D | 待派（0 未勾） | 无 | `internal/ball/tokens.go` | 0/5 | 五格已勾仍不结案：缺**非实现者裁决**＋桌面签收 |
| **245** | **A** | 已立未派（**排在队列最前**） | 无 | `internal/ball`＋`cmd/wisp` 热键 | 7/8 | 生产默认值＝裸 `Esc`（今天在跑的产品进程里生效） |
| 246 | A | 已立未派 | `245` 未闭（触发关系） | `cmd/wisp` 装配＋`proc` | 1/8 | 无否决通道执行者；与 245 同包 ⇒ 串行 |
| 247 | D | OPEN 未结案（8/2） | 残余归 `297` 等 | `internal/audio`＋`cmd/wisp`＋真麦 | 2/8 | 剩的两格要真机 |
| 248 | A | OPEN（按住·同包串行） | 无 | `internal/panel`＋`cmd/wisp` | 5/7 | 料齐；与 `33` 同包串行 |
| 249 | W | **WITHDRAWN**（10-01 自撤） | — | — | 7/0 | 立案前提错，就地留档 |
| 256 | C | 默认不排 | `255` 已闭＋owner 无相反意思 | `internal/agent/approval` 装配 | 2/2 | 编排者默认不排；要 owner 一句"常驻腿该不该吃配置" |
| 262 | E | 无 Status 行 | 无 | `scripts`＋`ci.yml` | 2/8 | 剩 `AC#5`＝CI self-hosted 侧真跑过 |
| 264 | C | 无 Status 行 | 无 | `internal/tools`（取舍归 owner） | 4/0 | `AC#0` 明令"不派写腿"，先取机主一句话 |
| **266** | **A** | 无 Status 行 | 无 | `scripts/*.sh`＋`lint` 一步 | 5/0 | `AC#0` 已交（`266-r1`），形状 ①-A 已定 |
| 269 | W | **WITHDRAWN／已并案**（并回票 236，236 已闭） | — | — | 5/0 | 与 236 `AC#6` 同一物理缺陷的第二次立案 |
| 271 | A | 已立未派 | 无 | `internal/tools` 测试尺 | 5/0 | 先派只读腿选形（甲/乙/丙三形代价表已交） |
| 274 | E | 无 Status 行 | `AC#0` 要 CI `lint-frontend` 真跑 | `scripts/build.ps1` | 3/9 | 选形与数量都压在 CI 那一发上 |
| **276** | **A** | 无 Status 行 | 无 | `scripts`＋台件归属名册 | 6/0 | 门今天就是 `exit 1`；判据＝豁免名单＋正控 |
| 284 | A | 无 Status 行 | 无 | `internal/agent/approval/gate.go` | 3/0 | **已落盘**：`gate.go:874-883` 注释三要素齐＋证据件在 ⇒ 只差非实现者翻框 |
| 290 | D | 5 勾 1 未勾 | 无 | 球/托盘＋真机 | 1/5 | `AC#3` ＝机主在场那发真机读数 |
| 291 | D | 1 勾 3 未勾 | 无 | 真麦 | 3/1 | `AC#2` ＝同一枚真麦的读数 |
| 293 | D | 5 勾 1 未勾 | 无 | 托盘＋真机 | 1/5 | `AC#5` ＝托盘勾的人眼读数 |
| 297 | D | 未开工 | 无 | 真麦（六窗读数在盘） | 5/0 | 现象已定实，判据要真麦 |
| 301 | E | 无 Status 行 | 无 | `scripts/portable-tests.sh`＋`ci.yml` | 2/5 | 剩 `AC#4`/`AC#4b`＝分档后的 CI 色 |
| 302 | E | 无 Status 行 | 无 | `scripts/portable-tests.sh` | 4/1 | 判据本身＝托管 runner 能不能给回执 |
| 305 | D | 无 Status 行 | 无 | WebView2 活文档 | 5/1 | 要真窗 + 页面自报 |
| 311 | C | 6 勾 0 未勾 | `Q-311-*` 两问未定 | `internal/winsec`＋runner 色 | 0/6 | 三笔残余具名：两问待人＋一发要 `cmd/wisp` 色 |
| 312 | F | 无 Status（**按住**） | 机主 nice-to-have | `frontend` 设置页＋Go 扫描 | 0/0 | 票面自写"当前按住不派" |
| 313 | B | 无 Status 行 | 无 | `frontend` 设置页＋凭据 | 3/0 | 判据在屏上（Go 半边已有 `313-v1` 收件） |
| 314 | D | 无 Status 行 | 无 | 真窗区崩溃 | 5/0 | 4/4 只在真窗区复现 |
| 316 | W | **withdrawn**（与 272 重复） | — | — | 1/0 | 口令已用，就地留档 |
| 323 | A | 已立未派 | 无 | `internal/llm`＋`cmd/wisp/panel_models_320.go` | 8/0 | 四把尺＋造世界；`AC#6` 有真台件尾巴 |
| 329 | B | 已立·`329-r1` 已交 | `AC#7` 不可判 | `cmd/wisp/run.go`＋`frontend` | 2/7 | 剩的"最后一公里"是页面侧零读者 |
| 330 | F | 已立未派 | `329` **开** | `cmd/wisp` 门铃 | 5/3 | 票面 `AC#7`："必须排在 329 之后" |
| 342 | C | 已立未派 | `AC#1` 拍形（编排者＋机主） | `internal/winsec` | 2/5 | 甲/乙都动产码 ⇒ 未拍形不许开工 |
| 345 | C | 已立未派 | `AC#5` 具名人工批准 | `internal/config` | 7/0 | 开工前须台账 `A##` 具名批准 |
| 346 | A | 已立未派 | 无 | `internal/winsec`＋`scripts` | 7/0 | 第一趟只读普查＋三形代价表，本机可做 |
| 350 | C | 已立未派 | `AC#0` 待机主拍"上/不上" | `internal/risk` 注释 | 7/0 | 造它＝人工批准 |
| 354 | C | ready-for-agent | `AC#1` 拍形（编排者＋机主） | `internal/winsec` | 6/1 | 票面逐字"未拍形任何实现都⛔ 开工" |
| 379 | W | **VOID**（立案前提是假的） | — | — | 5/0 | `HEAD` blob `CR=0`，工作树 CRLF 是 `autocrlf` 显示 |
| 388 | D | 待派 | 无 | `internal/audio`＋真机耳听 | 5/0 | `AC#2`＝衰减后出声可闻度（要真设备/机主确认） |
| 407 | B | open | 前端归外部 agent | `frontend/**`＋`internal/plugin` 只读 | 0/0 | 尚未成 AC 框；判据在屏上 |

---

## 4. A 的 23 枚"干净可派"（§1 ④ 展开）

> 判"干净"的四条：无 owner 门、无 `frontend/**` 写面、判据不要真机、前置票已闭。
> **`cmd/wisp` 是共享编译面：每波只许一枚写腿**（`AGENTS` 运行铁律）。下表标 ✅ 的互相串行。
> ⚠ **不列进下表的三类 A**：① 在飞 `23`（`internal/tools` 有未跟踪台件）／`194`／`197`／`200`／`211`；② 已落盘只差翻框 `284`；③ ⚑ 前提已倒 `176`／`180`（§6）。

| 票 | 面 | 未勾 | 串行键 | 可当波首？ |
|---|---|---|---|---|
| 98 | `scripts/**`＋`docs` | 5 | ❌ | ✅ |
| 103 | `internal/winsec`＋`internal/memory` | 4 | `internal/winsec` | ✅ |
| 123 | `cmd/wisp/*_test.go` | 6 | ✅ `cmd/wisp` | ⏳ 等 `cmd/wisp` 空 |
| 135 | `cmd/wisp`＋`internal/proc` | 7 | ✅ `cmd/wisp` | ⏳ |
| 136 | `internal/observe`＋`scripts` | 4 | ❌ | ✅ |
| 140 | `ci.yml`＋四包 | 4 | ❌（要登记） | ✅ |
| 145 | `internal/panel` 快照载体 | 3 | `internal/panel` | ⏳ 只一段 |
| 177 | `internal/tools`＋`internal/risk` | 3 | `internal/tools` | ✅ |
| 182 | 只读普查（零产码） | 6 | ❌ | ✅ |
| 185 | `internal/tools`＋`internal/risk` | 7 | `internal/tools` | ✅ |
| 188 | `internal/tools`＋`internal/panel` | 5 | `internal/tools` | ✅ |
| 216 | 显示层纯函数 | 1 | ❌ | ✅ |
| 220 | `internal/panel` 名册 | 1 | `internal/panel` | ⏳ 受 314 死区 |
| 228 | `cmd/wisp` 常驻＋`proc` | 11 | ✅ `cmd/wisp` | ⏳ 不是小活 |
| 238 | `internal/proc`＋`internal/winsec` | 4 | `internal/winsec` | ✅ |
| 245 | `internal/ball`＋`cmd/wisp` | 7 | ✅ `cmd/wisp` | ✅ |
| 246 | `cmd/wisp` 装配 | 1 | ✅ `cmd/wisp` | ⏳ 与 245 串行 |
| 248 | `internal/panel`＋`cmd/wisp` | 5 | ✅ `cmd/wisp` | ⏳ 与 33 串行 |
| 266 | `scripts/*.sh`＋`lint` 步 | 5 | ❌ | ✅ |
| 271 | `internal/tools` 测试尺 | 5 | `internal/tools` | ✅（先选形） |
| 276 | `scripts`＋台件名册 | 6 | ❌ | ✅ |
| 323 | `internal/llm`＋`cmd/wisp` | 8 | ✅ `cmd/wisp` | ⏳ |
| 346 | `internal/winsec`＋`scripts` | 7 | `internal/winsec` | ✅ |

⇒ **不撞 `cmd/wisp` 也不撞 `internal/winsec`/`internal/tools` 的"自由面"共 7 枚**：`98` `136` `140` `182` `216` `266` `276`（其中 `136` 也碰 `scripts` ⇒ 与 `266`/`276` 只撞"文档/脚本面"、不撞产码）。
**其余 16 枚必须排队**——这就是"小而自洽票已尽"的同一个事实，也是 `triage-180 §4` 的结论，本程复核未变。

---

## 5. 下一波建议（7 枚；按"买到的价值 ÷ 本程可证的成本"排）

| # | 票 | 面（写面，具名） | 为什么**现在**可派（可证的理由） | 第一件事 |
|---|---|---|---|---|
| **1** | **245** | `cmd/wisp` 热键那一族（`main.go:25` 那句与四枚钉）＋`internal/ball` | 生产默认值＝**裸 `Esc`**，自 09-30（球装进常驻那条腿）起在**正常跑的产品进程**里生效（不再是调试旁支）；票面自写"排在队列最前"；`8 勾/7 未勾`、写面现量 `git status --porcelain cmd/wisp internal/ball` ＝ 0 行；判据全在本机（钉＋真键注入） | 先按 `AC#6` 把**钉改成"组合键也算绑"**（四枚钉今天仍锚在"枚数"），再改行为 |
| **2** | **98** | `scripts/**`＋`docs`（**零共享编译面**） | 它是**其它所有腿的取证前提**：今天任何腿交"某格红/绿"都绕不过 `cmd/wisp`。带 `PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"` 的形态已实测可跑（本程与派单模板都用这条）；`C1` 那 76 枚孤儿红已被 `ae8136261` 修绿 39 枚（本程 `git log` 现读）⇒ 剩下的正是"让它留下"那半 | 把 PATH 口令写进脚本／派单模板＋一枚**永不 skip** 的守卫＋台账更正（`AC#1`–`AC#4`） |
| **3** | **103** | `internal/winsec`＋`internal/memory`（测试面） | `rejected-needs-fix`＝已被验收退回、四格全空、无前置票、无 owner 门；缺陷两形已被探针定死（`SetPathResolver` 可被"橡皮图章"静默改写外来 DACL／`RemoveUnlinked` 沿 junction 删掉外来真文件且返回 `nil`） | 给注册口加"只许装一次且必须是真解析器"的守卫＋`RemoveUnlinked` 的 junction 平台用例（必拒、不删目标） |
| **4** | **266** | `scripts/slo-*.sh`＋`lint` 作业新增一步 | `AC#0` 已由 `266-r1` 交件（形状 ①-A 已定、两枚钉真身行号现读、`gofumpt`/`d22scan` 基线读数齐）；外部牙**落在被检文件之外**＝不动任何阈值；判据本机可判（删钉必响＋还原绿） | 落新尺脚本＋接进 CI 一步（`if: !cancelled()`，先登记 `ci.yml`） |
| **5** | **276** | `scripts`＋`.scratch/wisp/probes/**` 台件归属名册 | 门**今天就是 `exit 1`**（17 枚来自 14 张别家票的探针件）；判据＝"豁免名单＋定向突变正控"，**零产码语义**；不用碰任何阈值/golden | 现跑复认 17 枚名册与 `outside .scratch = 0`，再逐枚归口 |
| **6** | **271** | `internal/tools` 一把计数尺（＋选形） | 只读普查 `271-a1` 已交三形代价表、假红那一发可复现；`AC#2` 逐字禁"放宽成 `<= 3`"⇒ 形状清楚；落地面只一枚测试文件 | 先由编排者拍甲/乙/丙（`AC#0`），再落"合法第三枚不红、误注册必红"的成对判据 |
| **7** | **216** | 显示底座纯函数（Go）＋错误码分支 | `5 勾/1 未勾`——只差 `AC#5`（错误提示按**错误码**分支，不许按英文串里有没有某个数字）；底座函数与其正控都已在盘上 | 只做 `AC#5` 一格＋一枚"按串分支"的反控 |

**备选（同一批可搭，不占 `cmd/wisp`）**：`238`（新机数据根无封条）、`188`（子代理/后台任务的 Go 侧源）、`180 ⚑`（`[panel] width` 补落地＋反控——**按新前提**写，别照"零读者"写）、`176 ⚑`（**作废/改票**，零产码，最便宜）。
**若只批准一枚 ⇒ `98`**：不改契约、不碰 `frontend`、不要真机与 CI 色，产出是一把**别人也能复跑的尺**。

⚠ **派单前必查的三件事**（本程现读，会变）：
1. `git status --porcelain -- cmd internal`（本程现读：`internal/secret/**` **4 改＋3 新**、`internal/tools/ticket23_image_channel_test.go` **未跟踪** ⇒ **确实有腿在飞**，派单前先确认是谁的、写面是否重叠）；
2. `cmd/wisp` 这一波是否已被占（`245`/`246`/`248`/`123`/`135`/`323` 全撞它）；
3. 票面 Status 的**最后一行**（本池里 `163`/`176`/`269`/`316`/`379` 都有"只追加"的更正行顶掉原句）。

---

## 6. 与同日 `triage-180`（13:0x，锚 `3ff1b1395`）的差异（本程改判 5 处，逐处给尺）

| 票 | `triage-180` | 本程 | 尺／依据 |
|---|---|---|---|
| **269** | A（仪器面） | **W（WITHDRAWN／已并案）** | 票面 `:3` 逐字"WITHDRAWN／已并案（10-05 12:3x 编排者自撤，出处＝`A620`）…并案回票 236 `AC#6`"，而 **236 已 `-done`** ⇒ 它是一枚**已闭缺陷的第二次立案**，不是可派的活 |
| **176** | A（先判"该不该作废"） | **A ⚑（前提已倒）** | 本程现跑：`grep -rn "RunAsync(" --include=*.go internal/ cmd/ \| grep -v _test` ⇒ **2 枚产码**（`internal/tools/subagent_197.go:372`、`cmd/wisp/run.go:1730`）；`Roster.Record` 产码 **4 处** ⇒ 票面 `AC#1` 自己写的"若 ≥1 则本票作废"**已触发** |
| **180** | A（"零读者"） | **A ⚑（前提已倒）** | 本程现跑：`grep -rn "cfg.Panel.Width" --include=*.go internal/ cmd/ \| grep -v _test` ⇒ `cmd/wisp/panel_resident_windows.go:208` **真读**，且 `config_readers_255.go:245` 已把它登记成 `hotClaimOtherProcess` |
| **284** | A（3 格，随批改注释） | **A（已落盘，只差翻框）** | 本程现读 `internal/agent/approval/gate.go:874-883` 注释三要素齐（"no production caller"／"not a defect to tidy away: do not delete it"／"the missing half is the answer-side verb…"）＋`docs/evidence/s1/284-gate-replay-zero-caller-comment.md` 在盘 ⇒ **不该再派写腿** |
| **23** | A（可写） | **A（在飞）** | `internal/tools/ticket23_image_channel_test.go` 未跟踪在盘、文件头逐字"Ticket 23 AC#2, the TOOLS half of the C7 image seam" |

**桶位口径差异（不是改判，是分法不同）**：`triage-180` 把"有前置票／需先裁形"的 23 枚仍留在 A；本程按派单口径把**前置票未闭**的 23 枚移入 **F**、把**判据要真机**的 16 枚移入 **D**、把**结案判据是 CI 色**的 8 枚移入 **E**。
⇒ 所以 `A 56 ⇄ A 31` 是**分桶口径**之差，不是"票变多了/活变少了"：两边"干净可当波首"的枚数是 **11（triage-180 §4.1）⇄ 23（本程，含只读普查与一格收口）**，量级同阶。

**同族提醒（沿用 `triage-180 §6.5`，本程复核未变）**：`33 34 36 37 38 39 40 46 47 48 49 51 53 54 55 56 57 58 60 25 32` 这批 09-19 的 S 切片票是**模板票**（`Claimed by: —`、Progress log 零条），`blocked by` 里点名的票大多今天已 `-done`。**直接派写手＝按两万年前的图纸施工**；先做一轮"票面对拉"比派写腿值钱。

---

## 7. 不计入派发池的 4 枚（W）

| 号 | 票面原样 | 理由 |
|---|---|---|
| 249 | `Status: OPEN … → WITHDRAWN（10-01 14:2x 编排者自撤）` | 立案前提（见"撤回理由"节）被自纠为错 |
| 269 | `WITHDRAWN／已并案（10-05 12:3x 编排者自撤，出处＝台账 A620）` | 与 236 `AC#6` 同一物理缺陷，236 已闭 ⇒ 撤票口令见票面，本表不转抄 |
| 316 | `withdrawn（口令「316 撤」已用）` | 与票 272 重复 |
| 379 | `VOID（10-24 18:3x 编排者自纠）` | `HEAD` blob `CR=0`；工作树 CRLF 是 `autocrlf` 的显示 |

> 按 `issues/README` 规则，**文件名不删不改**；本表只是不让它们回到派发队列。

---

## 8. 本程**没有**做的事（宁缺毋造）

- **没有跑**任何 `go test`／`scripts/*.sh`／`d22scan`／长度门：本件是分诊，不是交件自证。§5 里"可跑/已实测"那几句都注明了出处（本程现跑的是 `grep`／`git`／`ls`／`sed` 四类只读命令）。
- **没有**逐枚复跑票面里的读数（行号、枚数、台件名册都当**票面快照**引用；`triage-180` 已警告"这是个活树"）。
- **没有取**任何 CI run id／job id：E 桶 8 枚的"真 CI 色"全是**票面自写的结案条件**，不是我量到的读数。
- **没有裁定**任何一格的 PASS/FAIL：翻框权在编排者与非实现者裁决腿（`SPEC-12 §4.3`／`issues/README` 规则 6）。
- **没有验证** B 桶在 `frontend/**` 的实现程度（本编队零读零写面，只按票面判据归类）。
- **没有**改票名、`Status:`、任何 AC 框、台账；**零 commit、零 push**。
