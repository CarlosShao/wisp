# retriage-194 — 开放票池**再普查**（A–F＋W/M 八桶）＋下一波可派 A 清单（只读普查腿）

- **腿**：`retriage-194`（只读普查腿；第 194 轮波后重普查，接 `retriage-189`）。
- **锚**：派单写 `HEAD = c52f2955a`；本程现读 `git rev-parse --short=9 HEAD` ＝ **`c52f2955a`**，与派单锚**同一枚**（零漂移）。取数时刻 2026-11-06 09:2x–09:4x +08。
- **射程**：`.scratch/wisp/issues/*.md`，排除 `-done` 与 `README.md`。
  **尺**：`ls .scratch/wisp/issues/*.md | grep -v -- -done | grep -v README | wc -l` ⇒ **111 枚**。
  （`retriage-189` 记 115 枚 ⇒ 差的 4 枚＝**103／51／135／266 本轮已 `-done`**，逐枚 `Get-ChildItem` 双名核实，见 §1。）
- **与 retriage-189 的全部差异（本程逐枚核过票面，不是抄表）**：
  1. **结案 4 枚**：`103`（seal seam）· `51`（Tier2 goja）· `135`（128 自证）· `266`（SLO 自锁钉）——文件名均带 `-done`，现读核实；
  2. **出 A 入 F（gated）2 枚**：`346`（AC#0 要 `gh` 取两枚 CI run 锚，结票判据压在 CI 色 ⇒ 编排者按住）· `246`（剩余格修法归 `228`，编排者按住不单派）；
  3. **维持非 A**：`238`＝**M（owner-vetoed）**（票面 `A1558` 行逐字在：A476/Q-72 owner 原话否决、「重开仅认 owner『Q-72 重开』」）· `216`＝**F**（AC#5 归 S6 DEFERRED＋三枚具名残余按住，票面在）。
- **只读自证**：本程**零改**票面／台账／代码／契约；全仓唯一写入＝本文件。
  **零 `git add`／`git commit`／`git push`**。**未跑**任何 `go test`／`scripts/*.sh`／`d22scan`／`gh`（本件是分诊，不是交件自证）。
  盘上现量（派单前必查）：`git status --porcelain -- cmd internal` ＝ **2 行，全是未跟踪测试件**＝`cmd/wisp/panel_resident_323_windows_test.go`（**票 323 在飞**）＋`internal/tools/ticket23_image_channel_test.go`（**票 23 在飞**）。
- **尺（逐条可复跑，引用前先重跑）**：与 `retriage-189 §2` 同一套——
  未勾 AC＝列首 `- [ ]` 计数；Status＝票面最新只追加裁定行；前置票逐枚现读 `-done`；面＝`Packages` 或 AC 点名目录；
  **桶判定优先级：W ＞ M ＞ C ＞ B ＞ D ＞ E ＞ F ＞ A**。

---

## 1. 一句话答案

| 桶 | 枚 | 含义（本程用法） |
|---|---|---|
| **A** | **18** | **本编队现在可做的 Go／测试／文档活**：无 owner 门、无 `frontend/**`／`design/**` 写面、判据不要真机、前置票已闭（或本来就无前置） |
| B | 14 | 写面或**结票判据**在 `frontend/**`／`design/**`，归机主外部 agent（本编队结构性写不到） |
| C | 22 | 待人：改契约文字／花钱或凭据／机主一句话／机主亲拍形 |
| D | 16 | 判据只能在**真机**取：真麦、真声、真托盘、真窗／WebView2、桌面视觉签收、第二块屏 |
| E | 8 | 结票判据**就是某个 CI run 的颜色**（暂时取不到） |
| F | 27 | ⛔ 现在不许开工：前置**票**未闭／触发门未满足／编排者明确按住（含本轮 gated 的 `346`/`216`/`246`）／剩余格已归口别票 |
| W | 4 | 撤票／VOID／已并案（`249` `269` `316` `379`）——单列只为不被当活票 |
| M | 2 | moot／裁"不做、留档"（**`276`** moot · **`238`** owner-vetoed） |
| 合计 | **111** | — |

**A 的 18 枚内部再分**：

| 子档 | 枚 | 票号 |
|---|---|---|
| ① 在飞（盘上有未跟踪件，**本波不再派**，等收件） | 2 | `23`（`internal/tools` 镜像通道测试件）· `323`（`cmd/wisp` 面板模型测试件） |
| ② 已落盘、只差非实现者翻框（**不是写腿**） | 1 | `194` |
| ③ ⚑ 票面前提已被盘上读数顶倒（该做的是**作废/改票**，不是照着派腿） | 1 | `176` |
| ④ **干净可派**（本轮真正的候选池） | **14** | `98` `123` `136` `140` `177` `182` `185` `197` `200` `211` `228` `245` `248` `271` |

⇒ **A-count ＝ 18（含在飞/收件/待作废 4 枚）；干净可派 ＝ 14。§4 给 prioritized 6 枚；若只批一枚＝`185`。**

---

## 2. 全表（111 枚；`未/已` ＝ 顶层未勾／已勾 AC 数，本程现跑）

| 号 | 桶 | Status（票面摘） | 前置票（现读） | 面（Packages／写面） | 未/已 | 备注 |
|---|---|---|---|---|---|---|
| 21 | D | segment-1-landed；段 2 未做 | `17` 已闭 | `internal/agent/approval`＋球视觉 | 5/0 | 段 2 含"球视觉＝机主签收"，判据要真桌面 |
| 23 | A① | ready-for-agent（11-03 落地大半） | `20` 已闭 | `internal/tools`（system/window/input） | 5/0 | **在飞**：`internal/tools/ticket23_image_channel_test.go` 未跟踪 |
| 25 | F | ready-for-agent | `21`·`23` **开** | 真机脚本面（S3 验收） | 5/0 | 两枚前置都开着 |
| 32 | D | ready-for-agent | 前置已清 | 真机语音面（③听写④陪伴） | 5/0 | 判据要真麦真声 |
| 33 | D | ready-for-agent | 无 | `cmd/wisp`＋`internal/panel`（WebView2） | 13/1 | C27 宿主＝真窗；panel 簇公共前置 |
| 34 | B | ready-for-agent | `01` 已闭 | `frontend/**` | 5/0 | 纯前端脚手架 |
| 35 | F | ready-for-agent | `33`·`34` **开** | `internal/panel`（C17 桥） | 6/4 | panel 簇根闸 |
| 36 | B | ready-for-agent | `35` **开** | `frontend/**` | 6/0 | 结票判据在屏上 |
| 37 | B | ready-for-agent | `21`·`35` **开** | `frontend/**` | 6/0 | 同上 |
| 38 | B | ready-for-agent | `35` **开** | `frontend/**`＋Go `task.list` | 6/0 | 同上 |
| 39 | B | ready-for-agent | `35` **开** | `frontend/**` | 5/0 | 同上 |
| 40 | B | ready-for-agent | `35` **开** | `frontend/**` | 8/0 | 同上 |
| 46 | D | ready-for-agent | `41`–`45` 全闭 | 真机（空闲/KWS/电源） | 5/0 | S6 验收＝真机场景 |
| 47 | F | ready-for-agent | `21` **开** | `internal/agent/scheduler` | 6/0 | 前置开 |
| 48 | F | ready-for-agent | `37`·`47` **开** | `internal/agent/approval`（C18） | 5/0 | 链式前置 |
| 49 | F | ready-for-agent | `48` **开** | `approval`＋`perm` | 6/0 | 链式前置 |
| 53 | C | ready-for-agent | `48`·`51` | lark-cli 真凭据 | 4/0 | 花钱/凭据 |
| 54 | F | ready-for-agent | `47`/`48`/`49`/`53` **开** | 真机＋CI（S7 验收） | 5/0 | 前置全开 |
| 55 | C | ready-for-agent | `54`＋go-ahead | macOS（无硬件） | 5/0 | 机主点头＋无机器 |
| 56 | C | ready-for-agent | `54`＋SignPath | 外部账号/花钱 | 5/0 | 证书与发布账号 |
| 57 | F | ready-for-agent | `54`/`53` **开** | `docs/**`＋`internal/plugin` | 3/0 | 前置开 |
| 58 | C | ready-for-agent | `54`＋选型 | 英文模型＋前端 | 4/0 | 花钱＋选型归人 |
| 60 | C | ready-for-agent | 门槛②未满足 | 服务商实时 key（花钱） | 7/0 | 票面自写 AUTO-DEFER |
| 62 | D | review（从未被验收） | 无 | `internal/ball`＋真桌面 | 6/3 | 玻璃质感＝人眼签收 |
| 64 | D | review | `07` 已闭 | `internal/ball`＋第二块屏 | 1/8 | 热键/多显示器实拖 |
| 65 | C | blocked-on-owner | `62` **开** | `design/refs`＋`internal/ball` | 6/0 | 缺**参考图**＝owner |
| 68 | C | review（AC#2/#3 owner） | 无 | `internal/ball`＋桌面 | 2/2 | 剩两格归 owner 裁 |
| 70 | C | blocked-on-owner | 已解除 | `.github/workflows` | 2/4 | 剩的判据"不在代理手里" |
| 71 | E | in progress | `70` **开**（语义无关） | `cmd/balldebug`＋`ci.yml` | 6/1 | AC#5＝一次**真红** CI 读数 |
| 72 | E | review | 无 | `internal/risk`＋`cmd/wisp` | 2/4 | AC#4＝runner 侧逐字读数 |
| 77 | B | blocked（owner 收走前端） | 已解除 | `frontend/**` | 4/3 | 前端归外部 agent |
| 85 | E | open（排队） | 无 | `ci.yml`＋staticcheck 钉版 | 5/1 | AC#1/#2 要真 run 日志行 |
| 98 | A④ | open（票面自称可结案候选） | 无 | `scripts/**`＋`docs` | 5/0 | AC#3 已落地；剩 docs 半边＋**E 尾**（一次真 run）＋门禁读数 |
| 100 | F | **RESERVED** | 自身 `AC#0` 触发门未满足 | `internal/proc`＋真容器 | 7/0 | 触发门＝不许开工 |
| 112 | E | ready-for-review | 无 | `.github`＋`internal/winsec` | 5/0 | 只在 CI Windows runner 上红 |
| 122 | F | open | `85` **开** | 全仓 Go（staticcheck 清账） | 7/0 | 票面自写"85a 的下游" |
| 123 | A④ | open | `117` 已闭 | `cmd/wisp/*_test.go` | 6/0 | 本机可改可跑；AC#4 要 runner 真读数（**E 尾**）；`cmd/wisp` 槽被 323 占 |
| 136 | A④ | ready-for-agent | 无 | `internal/observe`＋`scripts` | 4/11 | 零样本 fail-closed 守卫**无钉**；54 用例拆守卫仍绿 |
| 138 | C | blocked（`Q-55`） | 无 | `internal/proc` | 2/2 | 等一次人工批准一起落 |
| 140 | A④ | open（未派） | 无 | `ci.yml`＋四包 | 4/0 | job 级 `WISP_ENV=test` 让 OS 加固用例在 runner 从不被咨询；改 `ci.yml` 要登记 |
| 145 | F | review／部分结案 | 归口 `35` **开** | `internal/panel` 快照载体 | 3/4 | panel 簇根闸＝`35` |
| 148 | C | blocked（双批准） | 无 | `internal/agent/approval` | 5/0 | 两枚前置都是人工批准 |
| 150 | C | blocked（`Q-55`） | 无 | `docs/PLAN`＋扫描器 | 5/0 | 三条出路都动冻结面/契约句 |
| 163 | C | 立而不派 | 无 | `internal/tools`（D34/D46） | 5/1 | 扩 `shell.exec` 语义＝契约 |
| 167 | F | ready-for-agent·立而不派 | `145` **开** | `internal/panel`＋`frontend` | 6/1 | 与 `145` 抢快照字段 |
| 168 | F | 立而不派 | `51` 已闭但插件地基未落 | `internal/plugin`（C17/C24） | 7/0 | 要动未定案白名单＝人工批准 |
| 169 | B | owner/前端会话接手件 | 无 | `frontend/src/.../right-rail.tsx` | 5/0 | 本编队结构性写不到；**永不派** |
| 170 | C | 立而不派 | 无 | `internal/agent/approval` | 6/0 | 产品口径归 owner |
| 174 | C | ready-for-agent（余 2 格） | 无 | `internal/tools`＋`PLAN.md:431` | 2/7 | 剩 `AC#4` 契约轴＝`Q-60` |
| 176 | A③⚑ | 立而不派（`175` 已闭） | 无 | `cmd/wisp`＋`internal/tools` | 4/1 | **AC#1 停手条件已触发**（`RunAsync` 产码 2 处、`Roster.Record` 4 处）⇒ 该作废/改票（零产码止损） |
| 177 | A④ | ready-for-agent（只读设计核） | 无 | `internal/tools`＋`internal/risk`（零产码） | 3/5 | `Q-61` 已批；与 185 同族可搭 |
| 180 | F | ready-for-agent（11-04 改判） | `33`/`35` **开** | `internal/config`＋`cmd/wisp` | 5/1 | 编审逐字判 blocked-by 33/35 |
| 182 | A④ | ready-for-agent（只读普查） | 无 | 只读普查（零产码） | 6/1 | **panel gated 簇的归口图**；产出 cascade 立新票 |
| 185 | A④ | 解除按住（`183` 已闭·归属已裁） | 无 | `internal/tools`＋`internal/risk`（测试面） | 7/0 | **主线上**：每次成功续读都造出下一次的 R4 阻断者 |
| 187 | B | ready-for-agent | 无 | 面板写入口＋`frontend` | 9/0 | 屏上那一半归前端 |
| 188 | F | ready-for-agent | 裁定"排 `145-r3` 之后" | `internal/tools`＋`internal/panel` | 5/1 | 编排者按住（`A1493`） |
| 191 | F | ready-for-agent | `33` **开** | `internal/tools`＋`internal/panel` | 6/1 | 票面"排在 `33-r1` 之后" |
| 192 | B | ready-for-agent | `33` **开** | `internal/panel`＋`frontend` | 5/1 | 预览/浏览器＝屏上那一半 |
| 193 | B | ready-for-agent | 无 | 前端门禁＋`config` | 7/1 | 放开的是**前端动画**与卡前端的门 |
| 194 | A② | ready-for-agent（裁决件在盘） | 无 | `internal/panel/bridge.go` | 7/0 | 两份名册对齐（owner 已裁、契约变更已批）⇒ 收件类 |
| 195 | F | blocked（等 `33` 片 B） | `33` **开** | `internal/config`＋`cmd/wisp` | 7/0 | 先地基、后面板 |
| 197 | A④ | 在派（Go 侧 9 枚实体在盘） | 无 | `internal/tools`＋`internal/panel` | 5/2 | 子代理状态/流式页；要 `panel` 槽 |
| 200 | A④ | 在派（owner 放权） | 无 | `internal/projctx`＋`agent`＋`panel` | 8/0 | 项目规矩注入；面宽 |
| 201 | F | 部分结案（5 未勾） | `228`·`246` **开** | 球/托盘原生面 | 5/2 | 剩四格归口"宿主原生入口" |
| 211 | A④ | 在派队列里 | 无 | `bridge.go`＋`budgets.go` | 5/0 | 池常量⇄桥天花板可见化＋硬拒；两枚 `4` 明令不动 |
| 213 | F | 待派 | `33` **开** | `cmd/wisp/commands.go` | 5/0 | 三票共用入向线 |
| 214 | F | 待派 | `33` **开** | `internal/tools`（附件受理） | 6/0 | 同上 |
| 215 | F | 待派 | `33` **开** | `internal/plugin`＋装配 | 5/0 | 同上 |
| 216 | F | **gated**（AC#5 归 S6 DEFERRED） | 无 | `internal/panel` 显示底座 | 1/5 | 票面：⛔ 改名 `-done`、残余具名按住；推 CI 会红一发等机主 |
| 219 | F | 接单＝做 | `201` **开** | 原生卡面＋审计 | 6/0 | 派单条件写死"排 201→223→224 之后" |
| 220 | D | 待派 | 无 | `internal/panel` 名册 | 1/4 | AC#5 终态读数被 **314 真窗死区**挡住（`A1546` STOP） |
| 226 | F | 产码已交·残余转 `219` | `219` **开** | `cmd/wisp`＋`internal/config` | 1/5 | 剩一格已明确转票 |
| 228 | A④ | 部分勾（方向＝甲已裁） | 无 | `cmd/wisp` 常驻＋`internal/proc` | 11/5 | 不是小活；验收要真窗（**D 尾**）；`cmd/wisp` 槽 |
| 230 | C | 待派（AC#3 已处置） | `229` 已闭 | 台账/证据＋注释 | 3/2 | 未闭＝AC#4 **owner/改契约**＋AC#5 `cmd/wisp` 死区 |
| 238 | **M** | **owner-vetoed（A1558）** | — | `internal/proc`＋`internal/winsec` | 4/0 | A476/Q-72 owner 原话否决；重开仅认「Q-72 重开」 |
| 239 | D | 待派（0 未勾） | 无 | `internal/ball/tokens.go` | 0/5 | 缺**非实现者裁决**＋桌面签收 |
| 245 | A④ | 已立未派（排队最前） | 无 | `internal/ball`＋`cmd/wisp` 热键 | 7/8 | 裸 `Esc` 在跑的产品进程里生效；AC#8/9 真机（**D 尾**）；`cmd/wisp` 槽 |
| 246 | **F** | **gated**（编排者按住） | 无 | `cmd/wisp` 装配＋`proc` | 1/8 | 剩格修法归 `228`，不单派 |
| 247 | D | OPEN（8/2） | 残余归 `297` 等 | `internal/audio`＋真麦 | 2/8 | 剩两格要真机 |
| 248 | A④ | OPEN（按住·同包串行） | 无 | `internal/panel`＋`cmd/wisp` | 5/7 | 料齐；AC#4 真机（**D 尾**）＋与 `33` 同包串行 |
| 249 | W | **WITHDRAWN**（10-01 自撤） | — | — | 7/0 | 立案前提错，就地留档 |
| 256 | C | 默认不排 | `255` 已闭 | `internal/agent/approval` 装配 | 2/2 | 要 owner 一句"常驻腿该不该吃配置" |
| 262 | E | 无 Status 行 | 无 | `scripts`＋`ci.yml` | 2/8 | 剩 `AC#5`＝CI self-hosted 侧真跑过 |
| 264 | C | 无 Status 行 | 无 | `internal/tools`（取舍归 owner） | 4/0 | `AC#0` 明令"不派写腿"，等机主一句话 |
| 269 | W | **WITHDRAWN／已并案** | — | — | 5/0 | 与 `236` `AC#6` 同一物理缺陷，`236` 已闭 |
| 271 | A④ | 已立未派 | 无 | `internal/tools` 测试尺 | 5/0 | 先由编排者拍甲/乙/丙（AC#0），再落"合法第三枚不红" |
| 274 | E | 无 Status 行 | AC#0 要 CI 真跑 | `scripts/build.ps1` | 3/9 | 选形与数量压在 CI 那一发上 |
| 276 | **M** | moot/不做（`A1556`） | 无 | `scripts`＋台件名册 | 6/0 | 门已 green（`362` exclude）；本票期望与之相反 |
| 290 | D | 5 勾 1 未勾 | 无 | 球/托盘＋真机 | 1/5 | `AC#3`＝机主在场真机读数 |
| 291 | D | 1 勾 3 未勾 | 无 | 真麦 | 3/1 | `AC#2`＝同一枚真麦读数 |
| 293 | D | 5 勾 1 未勾 | 无 | 托盘＋真机 | 1/5 | `AC#5`＝托盘勾人眼读数 |
| 297 | D | 未开工 | 无 | 真麦 | 5/0 | 六窗读数在盘，判据要真麦 |
| 301 | E | 无 Status 行 | 无 | `scripts/portable-tests.sh`＋`ci.yml` | 2/5 | 剩分档后的 CI 色 |
| 302 | E | 无 Status 行 | 无 | `scripts/portable-tests.sh` | 4/1 | 判据本身＝托管 runner 能否给回执 |
| 305 | D | 无 Status 行 | 无 | WebView2 活文档 | 5/1 | 要真窗＋页面自报 |
| 311 | C | 6 勾 0 未勾 | `Q-311-*` 两问未定 | `internal/winsec`＋runner 色 | 0/6 | 残余待人＋一发要 `cmd/wisp` 色 |
| 312 | F | 无 Status（**按住**） | 机主 nice-to-have | `frontend` 设置页＋Go 扫描 | 0/0 | 票面自写"当前按住不派" |
| 313 | B | 无 Status 行 | 无 | `frontend` 设置页＋凭据 | 3/0 | Go 半边已有 `313-v1` 收件；判据在屏上 |
| 314 | D | 无 Status 行 | 无 | 真窗区崩溃 | 5/0 | 4/4 只在真窗区复现 |
| 316 | W | **withdrawn**（与 272 重复） | — | — | 1/0 | 口令已用，就地留档 |
| 323 | A① | 已立（写腿在飞） | 无 | `internal/llm`＋`cmd/wisp` | 8/0 | **在飞**：`cmd/wisp/panel_resident_323_windows_test.go` 未跟踪 ⇒ 先收件再派 |
| 329 | B | 已立·`329-r1` 已交 | AC#7 不可判 | `cmd/wisp/run.go`＋`frontend` | 2/7 | 剩"最后一公里"是页面侧零读者 |
| 330 | F | 已立未派 | `329` **开** | `cmd/wisp` 门铃 | 5/3 | 票面 `AC#7`："必须排在 329 之后" |
| 342 | C | 已立未派 | AC#1 拍形 | `internal/winsec` | 2/5 | 未拍形不许开工 |
| 345 | C | 已立未派 | AC#5 具名人工批准 | `internal/config` | 7/0 | 开工前须台账具名批准 |
| 346 | **F** | **gated**（AC#0 要 `gh` CI 锚） | 无 | `internal/winsec`＋`scripts` | 7/0 | 两枚竞态仪器弃权报红；第一趟普查可做但结票压 CI 色 ⇒ 编排者按住 |
| 350 | C | 已立未派 | AC#0 待机主拍 | `internal/risk` 注释 | 7/0 | 造它＝人工批准 |
| 354 | C | ready-for-agent | AC#1 拍形 | `internal/winsec` | 6/1 | 票面逐字"未拍形任何实现都⛔ 开工" |
| 379 | W | **VOID** | — | — | 5/0 | 立案前提是假的（`HEAD` blob `CR=0`） |
| 388 | D | 待派 | 无 | `internal/audio`＋真机耳听 | 5/0 | 衰减后可闻度＝真设备/机主确认 |
| 407 | B | open（11-05 新立） | 前端归外部 agent | `frontend/**` | 0/0 | 尚未成 AC 框；判据在屏上 |

> 本程核对：A 18 ＋ B 14 ＋ C 22 ＋ D 16 ＋ E 8 ＋ F 27 ＋ W 4 ＋ M 2 ＝ **111**，与池清点数一致。

---

## 3. 与 retriage-189 的差异明细（逐处给尺）

| 票 | retriage-189 | 本程 | 尺／依据 |
|---|---|---|---|
| **103** | A②（收件类） | **不在池内（`-done`）** | `103-…-done.md` 双名核实 |
| **51** | A④（契约门开） | **不在池内（`-done`）** | `51-tier2-goja-done.md` |
| **135** | A④ | **不在池内（`-done`）** | `135-…-done.md` |
| **266** | A② | **不在池内（`-done`）** | `266-…-done.md` |
| **346** | A④（E 尾） | **F（gated）** | AC#0 要 `gh` 取两枚 CI run 锚（票面在），结票判据压 CI 色 ⇒ 编排者本轮按住 |
| **246** | A④（§5 #6 推荐） | **F（gated）** | 剩余格修法归 `228`（票面自写"两枚合起来才是那句产品事实"）；编排者本轮按住不单派 |
| **238** | M（`A1558` 更正行） | **M（不变）** | 票面 `A1558` 行现读在：owner-vetoed、勿再派 |
| **216** | F | **F（不变）** | 票面残余具名＋AC#5 归 S6 DEFERRED 现读在 |
| **323** | A②（收件类） | **A①（在飞）** | 盘上未跟踪 `cmd/wisp/panel_resident_323_windows_test.go`（本程现读 `git status --porcelain -- cmd internal`） |
| **23** | A①（在飞） | A①（不变） | `internal/tools/ticket23_image_channel_test.go` 仍未跟踪 |

⇒ A 由 189 的 25（更正后实为 24）变为 **18**：−4（结案）−2（346/246 gated）−1（238 已在更正中出 A）＋1（323 仍留在 A 内、只是 ②→①；对总数无增减）。
其余 97 枚桶位逐枚与本程现读票面复核，无改判。

---

## 4. prioritized 可派 A 清单（6 枚；按"买到的价值 ÷ 本程可证的成本"排）

> 串行硬约束（本程现读）：`cmd/wisp` 槽被 **323 在飞**占（未跟踪测试件在包里）⇒ `123`/`228`/`245`/`248` 全部让行；
> `internal/tools` 被票 **23 在飞**的测试件占着同一包 ⇒ `185`/`271`/`177` 串行、不与 23 同波。

| # | 票 | 面（写面，具名） | 未勾 | 为什么**现在**可派（可证的理由） | 第一件事 |
|---|---|---|---|---|---|
| **1** | **185** | `internal/tools/**`＋`internal/risk/**`（测试面，零产码语义） | 7 | 前置 **183 已 `-done`**、按住已解除（票面 09-28 归属裁定在）；缺陷在**主线上**：`task.output`→`fs.read` 成功→同一路径再发即命中 R4 ⇒ **每一次成功续读都亲手造出下一次的阻断者**；无 owner/前端/真机门；与 23 同包 ⇒ 排在 23 收件之后即开 | 把"第二发为什么被拒"落成 `internal/tools/` 常驻红判据（HEAD 先红），再答"这枚 C25 声明该由谁给"（桥侧 vs 落盘层） |
| **2** | **182** | 只读普查（零产码；写点只在本票面＋`docs/evidence/s1/182-*`） | 6 | 与任何腿零冲突；它是 **145/187/188/192/195/213/214/215/248 那一整簇 gated 票的归口图**——不先出这张图，panel 簇只能一直卡在"谁欠谁"；普查件 `182-census.md` 已在盘、AC#1 已勾 | AC#2 十堆逐堆三档定性（有源／无源但规格要求／规格外）＋AC#3 逐堆指认归属票（允许 cascade 立新票） |
| **3** | **136** | `internal/observe`＋`scripts` | 4 | 承重墙**无钉**：拆掉 `sampler.go` 零样本 fail-closed 守卫，整包 54 用例仍绿（票 134 验收方判高）；无前置、无 owner/前端/真机门；不撞 `cmd/wisp`/`internal/tools` | AC#5（拒沉默：补禁令或明写"不查"）＋AC#6（`slo-fresh` schedule 从未响）＋门禁读数 |
| **4** | **271** | `internal/tools` 测试尺 | 5 | 全池最小的写腿之一：票 221 的"task 家族计数"尺对合法第三枚工具**假红**；前置全闭、无 owner 门；与 185 同包串行 | AC#0 由编排者拍甲/乙/丙（哪支形状），然后落"合法第三枚不红、误注册必红" |
| **5** | **98** | `scripts/**`＋`docs` | 5 | 票面自称"可结案候选"；AC#3 已由 `4a682f698` 落地；剩的是 **docs 半边（AC#4 台账收敛）＋门禁读数**——不占任何 Go 写面；**E 尾**（AC#1(iii) 一次真 run）具名留档不冒充 | AC#4 台账收敛＋AC#5 门禁读数；AC#1(iii) 的 CI 真读数具名挂 E 尾等窗口 |
| **6** | **140** | `ci.yml`＋四包 | 4 | 家族账（`R-128-*`）的具名根因：job 级 `WISP_ENV: test` 让"要先问操作系统才生效"的加固用例在 runner **从不被咨询**；先量面再裁形，全本机可做；改 `ci.yml` 要登记（有先例规矩） | 只读量面：四包里哪些用例真吃 `WISP_ENV`，出裁形清单再落 |

**备选（同批可搭，不占上表串行键）**：`177`（零产码设计核＋分账，Q-61 已批；与 185 同包串行）·
`200`／`197`（两张功能大票，面宽、要 `internal/panel`／`cmd/wisp` 槽，等 323 收件）·
`176`（零产码"作废/改票"止损，全池最便宜的纸面活）· `194`（收件类：只差非实现者翻框，不占写腿）·
`211`（`bridge.go`＋`budgets.go`，可见化＋硬拒，两枚 `4` 明令不动）。

**若只批准一枚 ⇒ `185`**：不改契约、不碰 `frontend`、不要真机与 CI 色，修的是**主线上每一次成功续读都自造阻断者**的活缺陷，产出是一张别人也能复跑的红/绿判据。

⚠ **派单前必查的三件事**（本程现读，会变）：
1. `git status --porcelain -- cmd internal`（本程＝2 行未跟踪：`cmd/wisp/panel_resident_323_windows_test.go`＋`internal/tools/ticket23_image_channel_test.go` ⇒ 323/23 在飞，两包的槽都占着）；
2. 票面 `Status` 的**最后一行**（只追加更正行顶掉原句是本仓常态）；
3. 前置票逐枚现读 `-done`（本表"前置票（现读）"列之外的新点名，派单时现查）。

---

## 5. 本程**没有**做的事（宁缺毋造）

- **没有跑**任何 `go test`／`scripts/*.sh`／`d22scan`／长度门／`gh`：本件是分诊，不是交件自证。
- **没有**逐枚复跑票面里的读数（行号、枚数、台件名册都当**票面/账快照**引用；这是一棵**活树**）。
- **没有取**任何 CI run id：E 桶 8 枚的"真 CI 色"全是票面自写的结案条件，不是本程量到的读数。
- **没有裁定**任何一格的 PASS/FAIL：翻框权在编排者与非实现者裁决腿（`SPEC-12 §4.3`／`issues/README` 规则 6）。
- **没有**改票名、`Status:`、任何 AC 框、台账；**零 `git add`、零 commit、零 push**；全仓唯一写入＝本文件。
- 谁要推翻某格桶位，请按 §头 的尺重跑那一格，不要按本表的口径反推票面。
