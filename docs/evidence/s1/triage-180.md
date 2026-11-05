# triage-180 — 全开放票池的**定义性分诊**（A–E 五桶）＋ DONE 票孤儿红盘点

**性质**：只读分诊件（裁决表）。**代码／票面／脚本／台账一字未改，零 `git add` / `git commit` / `git push`；本文件是全仓唯一写入。**
**锚点**：`HEAD = 3ff1b13953e9386ecdbdec1e2e0c6873a5d2ce2f`（分支 `dev`），取数时刻 **2026-11-05 13:0x–13:1x +08**。
**尺（逐条可复跑）**：
- 开放集＝`ls .scratch/wisp/issues/*.md | grep -v -- -done.md | grep -v README` ⇒ **117 枚文件**。
- 未勾＝`grep -cE '^[[:space:]]*-[[:space:]]*\[[[:space:]]\]'`；已勾＝`grep -cE '^[[:space:]]*-[[:space:]]*\[[xX]\]'`（两把尺都按行首锚定，票正文里的示例框不计）。
- 状态＝票面第一行含 `Status` 的那行（原样摘，不解释）；无 `Status` 行＝票面本来就没写（本表标「无 Status 行」，⛔ 不等于不存在）。
- 面＝票面 `Packages` / `写面` / AC 里点名的目录；**判"是不是前端"按目录不按票号**（`issues/README` 规则 7）。

> ⚠ **这是个活树**：本程起手时 `300-…` 已在写件前一轮被收口改名（`-done`），`406` 已被 `git add`（`git status` 现读 `A .scratch/wisp/issues/406-…`）。⇒ **任何引用本表的人先重跑上面三把尺**，别把这张表当常量。

---

## 0. 一句话

**117 枚未 `-done` 文件里，真正"开放"的是 114 枚**（3 枚不是 open：`249` WITHDRAWN、`316` withdrawn、`379` VOID）。
**114 的分桶：A 56 ／ B 12 ／ C 17 ／ D 21 ／ E 8。**
**能由本编队（Go／脚本）动手的是 56 枚，但其中只有 11 枚是"零前置、可当波首动作"的**；其余 45 枚要么被**同包串行**卡住（`internal/panel`／`internal/winsec`／`cmd/wisp`），要么要先裁形，要么只是收口（0–1 格）。
**DONE 票的孤儿红：81 枚全在 `cmd/wisp/**`，`internal/**` 与 `frontend/**` 各 0 枚；其中 52 枚的归属票已经 `-done`＝今天没有主人。**

---

## 1. 桶判据（谁必须动手才能关这一票）

| 桶 | 判据（缺一即不算这一桶） | 反面 |
|---|---|---|
| **A 本编队可做-Go/脚本** | 写面落在 `cmd/**`、`internal/**`、`scripts/**`、`tools/**`、`.github/**`、`docs/evidence`；关票判据可由本机跑出来（含 `-run` 定向、含 Docker Linux 腿）；无 frontend／无真机读数／无真 CI 色 | 只要有**一格**必须人眼看屏、必须推上去看色、必须动契约文字、必须 owner 签字 ⇒ 降到对应桶 |
| **B 前端/design 外部agent** | 关票判据必须落在 `frontend/**` 或 `design/**`（本编队结构性零写面） | Go 半边另有归属票时可标 B（Go 半＝A 的另算） |
| **C owner 改契约文字/花钱/不可逆** | 前置是：改 `PLAN.md`/`docs/specs/**` 的**文字**、花钱（模型/API key/签名服务）、或不可逆动作（删证据台件/推/发布）；或 owner 一句话未答即不许开工 | 「需要人工批准」≠ C：**批准已落**（具名 `A##`）的仍算 A |
| **D 真机/WebView2** | 判据本体是**真桌面/真窗口/真麦克风/真托盘/第二块屏/WebView2 页面/真 .exe** 的读数或人眼签收 | macOS 硬件归此桶（本编队无该平台机器） |
| **E 需真 CI run** | 结案条件写死为「某 step 在 CI 上给出过 success／run id＋job id＋step 号」；本地绿不算 | 本地能证、只是顺手要 CI 的尾巴 ⇒ 仍 A，尾巴注明 |

**一处刻意的取舍**：一只票同时踩两桶时，**按"现在真正卡住它的那一桶"归档**（不是按票的大小）。例：`145` 有前端半，但它的 3 格未勾全在 Go 载体上 ⇒ A；`167` 的 Go 载体能写，但它要交的东西是"屏上看得见" ⇒ B。

---

## 2. 分桶计数

| 桶 | 枚 | 备注 |
|---|---|---|
| **A** | **56** | 内部再分：零前置可立刻起手 **11** ／ 需同包串行 **9** ／ 仪器与普查 **13** ／ 有前置或先裁形 **23** |
| **B** | **12** | 全部要 `frontend/**`；`313`/`312` 还要设置页 |
| **C** | **17** | 花钱或凭据 **4**（53/56/58/60）／ 契约文字 **5**（150/163/168/174/187）／ 具名批准门 **4**（100/138/148/345）／ owner 一句话 **4**（170/177/264/350） |
| **D** | **21** | 球与托盘签收 **9** ／ 真麦 **4** ／ WebView2·真窗口 **4** ／ 真机场景 **3** ／ macOS（无硬件）**1** |
| **E** | **8** | 结案判据就是 CI 色 |
| 不计入 | 3 | `249` `316` `379`（见 §9） |
| **合计** | **114** | ＋3＝117 枚文件 |

---

## 3. 全表（117 枚未 `-done` 文件＝114 枚开放＋3 枚不计；`未/已` ＝ 未勾／已勾 AC 框数）

| 号 | 一句标题 | Status（票面原样摘） | 未/已 | 面（Packages） | 桶 |
|---|---|---|---|---|---|
| 21 | Approval gates minimal：L1 阻止窗／L2 原生卡／简易队列／批量 D45-1 | segment-1-landed | 5/0 | `internal/agent/approval`＋`cmd/wisp` | A |
| 23 | system.\*/window.\*/clipboard/media/screen/input.type 全套 | ready-for-agent | 5/0 | `internal/tools`＋`internal/llm`(C7) | A |
| 25 | S3 验收：场景①② ＋ 安全红队全套 | ready-for-agent | 5/0 | 真机脚本面 | D |
| 32 | S4 验收：③听写 ④陪伴 | ready-for-agent | 5/0 | 真机语音面 | D |
| 33 | 面板宿主 C27（WebView2 窗口＋embed 服务） | ready-for-agent | 13/1 | `cmd/wisp`＋`internal/panel` | D |
| 34 | 前端脚手架（React/TS/Tailwind/shadcn/token） | ready-for-agent | 5/0 | `frontend/**` | B |
| 35 | PanelBridge C17：白名单／corr 路由／无状态重同步 | ready-for-agent | 6/4 | `internal/panel` | A |
| 36 | 结果与历史面板（流式／14 态／sanitizer） | ready-for-agent | 6/0 | `frontend/**` | B |
| 37 | 审批 UI（队列视图／L2 卡片面板级） | ready-for-agent | 6/0 | `frontend/**` | B |
| 38 | 命令面板＋任务视图（文本入口／force 规则） | ready-for-agent | 6/0 | `frontend/**`(＋Go task.list) | B |
| 39 | 配置编辑器 GUI（三档徽章／安全锁） | ready-for-agent | 5/0 | `frontend/**` | B |
| 40 | 安全/隐私/成本页（grants／黑名单／记忆／成本） | ready-for-agent | 8/0 | `frontend/**` | B |
| 46 | S6 验收：空闲/KWS 门／settle／电源预演 | ready-for-agent | 5/0 | 真机 | D |
| 47 | TaskScheduler＋PathLock（C20）多任务解锁 | ready-for-agent | 6/0 | `internal/agent/scheduler` | A |
| 48 | ApprovalQueue 全量（C18）FIFO／只看头／fail-closed | ready-for-agent | 5/0 | `internal/agent/approval` | A |
| 49 | 会话授权 D45-2（tool,pattern,session）＋input.type 绑定 | ready-for-agent | 6/0 | `internal/agent/approval`＋`internal/perm` | A |
| 51 | Tier2 goja 运行时（C24 宿主 API／墙钟沙箱／默认关） | ready-for-agent | 6/0 | `internal/plugin` | A |
| 53 | lark-cli 包装插件 e2e（S7"真插件"） | ready-for-agent | 4/0 | 真服务＋凭据 | C |
| 54 | S7 验收：并发红队／生命周期预演／自用一周 | ready-for-agent | 5/0 | 真机＋CI | E |
| 55 | S8 macOS 移植 | ready-for-agent | 5/0 | 无该平台机器 | D |
| 56 | S8 签名＋分发＋P10 命名 | ready-for-agent | 5/0 | 外部账号/花钱 | C |
| 57 | S8 插件 SDK／文档／示例／registry | ready-for-agent | 3/0 | `docs/**`＋`internal/plugin` | A |
| 58 | S8 i18n／英文语音模型 | ready-for-agent | 4/0 | 模型下载＋`frontend` | C |
| 60 | C32 RealtimeEngine（云端 S2S） | ready-for-agent | 7/0 | 服务商 key＝花钱 | C |
| 62 | 液态玻璃球视觉重做（玻璃/吸附/音频驱动） | review | 6/3 | `internal/ball`＋真桌面 | D |
| 64 | 票 07 遗留：热键/交互四项/零定时器/多屏 | review | 1/8 | `internal/ball`＋第二块屏 | D |
| 65 | 球玻璃质感返工（要 owner 参考图） | blocked-on-owner | 6/0 | `design/refs`＋`internal/ball` | D |
| 68 | 默认构建画的就是签收那个球（`prototypeVisuals`） | review | 2/2 | `internal/ball`＋桌面 | D |
| 70 | 让 CI 真成护栏（5 job 全红逐因） | blocked-on-owner | 2/4 | `.github/workflows` | E |
| 71 | 每道门自报工作量（-diff 阈值门/逐作用域自报） | in progress | 6/1 | `cmd/balldebug`＋`tools/d22scan`＋`ci.yml` | A |
| 72 | A 表在 runner 上退化成 B 表（C26 纵深防御） | review | 2/4 | `internal/risk`＋`cmd/wisp` | E |
| 77 | 前端脚手架（vendored react-bits/beautifului） | blocked（owner 收走） | 4/3 | `frontend/**` | B |
| 85 | lint 两把工具从未产生判据（钉版本＋清单） | open（排队） | 5/1 | `ci.yml`＋`staticcheck` | E |
| 98 | `go test ./cmd/wisp/` 在这台机上永远测不到 | open | 5/0 | `scripts/**`＋`docs` | A |
| 100 | AppContainer 落地票（触发门未满足，RESERVED） | RESERVED | 7/0 | `internal/proc`＋真容器 | C |
| 103 | 密封两处旁路（`SetPathResolver` 可被橡皮图章） | rejected-needs-fix | 4/0 | `internal/winsec`＋`internal/memory` | A |
| 112 | winsec 步在 runner 首跑三条红 | ready-for-review | 5/0 | `.github`＋`internal/winsec` | E |
| 122 | staticcheck 清账（85a 下游） | open | 7/0 | 全仓 Go | A |
| 123 | runner 上 4 枚 cmd/wisp 用例假设有真人审批 | open | 6/0 | `cmd/wisp/*_test.go` | A |
| 135 | 票 128 判据无自证腿（双叠假绿） | ready-for-agent | 7/1 | `cmd/wisp`＋`internal/proc` | A |
| 136 | "有报告⇒有数字"承重墙无用例钉 | ready-for-agent | 4/11 | `internal/observe`＋`scripts` | A |
| 138 | 取消任务不杀它自己拉起的子进程 | blocked（Q-55） | 2/2 | `internal/proc` | C |
| 140 | job 级 `WISP_ENV: test` 让 OS 加固腿永不被问 | open（未派） | 4/0 | `ci.yml`＋四包 | A |
| 145 | 面板快照只有 4 字段，14 态中 11 态无输入 | ready-for-agent（未派） | 3/4 | `internal/panel`(快照载体) | A |
| 148 | `cloneDecision` 只护一扇门（replay 出口仍是浅拷贝） | blocked（双批准） | 5/0 | `internal/agent/approval` | C |
| 150 | 4 枚 `DEFERRED` 标记无 §5 登记行 | blocked（Q-55） | 5/0 | `docs/PLAN`＋扫描器 | C |
| 163 | 命令执行要补"常驻会话"（仍长得像一次性） | 立而不派 | 5/1 | `internal/tools`（D34/D46） | C |
| 167 | 面板四枚小出口：占用/排队/停止/草稿 | ready-for-agent·立而不派 | 6/1 | `internal/panel`＋`frontend` | B |
| 168 | 内置浏览器做成可关插件（两件前置） | 立而不派 | 7/0 | `internal/plugin`(C17/C24) | C |
| 169 | 一枚装饰字符把整仓门钉红 | owner/前端会话 | 5/0 | `frontend/**` 单文件 | B |
| 170 | CLI 的 L1 阻止窗无人能否决 | 立而不派 | 6/0 | `internal/agent/approval` | C |
| 174 | `task.output` 许诺的路默认读不回来 | ready-for-agent | 2/7 | `internal/tools`＋`PLAN.md:431` | C |
| 176 | 没有起跑口：`RunAsync`／`TaskRoster.Record` 零调用 | 立而不派（排 175 后） | 4/1 | `cmd/wisp`＋`internal/tools` | A |
| 177 | 盖戳 vs 续读：四支候选取一支 | ready-for-agent（只读腿） | 3/5 | `internal/tools`＋`internal/risk` | C |
| 180 | `[panel] width` 生产零读者 | ready-for-agent | 5/1 | `internal/config`＋`cmd/wisp` | A |
| 182 | 任务监控栏在 145 的十四行表之外 | ready-for-agent（只读） | 6/1 | 只读普查（零产码） | A |
| 185 | 每次成功续读亲手造出下一次阻断者 | 立、按住（等 183-v1） | 7/0 | `internal/tools`＋`internal/risk` | A |
| 187 | 聊天框模型/思考档要能当场改 | ready-for-agent | 9/0 | 面板写入口（契约）＋`frontend` | C |
| 188 | 子代理/后台任务两堆今天无源 | ready-for-agent | 5/1 | `internal/tools`＋`internal/panel` | A |
| 191 | 面板"内置终端"＝执行面进通道 | ready-for-agent | 6/1 | `internal/tools`＋`internal/panel` | A |
| 192 | 浏览器/产物预览没人认领 | ready-for-agent | 5/1 | `internal/panel`＋`frontend` | B |
| 193 | "轻量级"不再约束前端动画 | ready-for-agent | 7/1 | `internal/panel` 门禁＋config | A |
| 194 | 两份面板方法名册互不相认 | ready-for-agent | 7/0 | `internal/panel/bridge.go` | A |
| 195 | 面板改配置缺"按节写入"地基 | blocked（等 33 片 B） | 7/0 | `internal/config`＋`cmd/wisp` | A |
| 197 | 子代理状态＋各自流式工作页 | 在派 | 5/2 | `internal/tools`＋`internal/panel` | A |
| 200 | 项目说明（AGENTS.md 那类）今天没人读 | 在派 | 8/0 | `internal/projctx`＋`internal/agent` | A |
| 201 | 到点该问的没人问：L1 自动跑、L2 自动拒 | 部分结案 | 5/2 | 球/托盘原生面 | D |
| 211 | 8 枚并发子代理要先动 D38d 天花板 | 在派队列 | 5/0 | `internal/tools/bridge.go`＋`budgets.go` | A |
| 213 | 加号菜单"命令"类（斜杠命令目录） | 待派 | 5/0 | `cmd/wisp/commands.go` | A |
| 214 | 加号菜单"往对话里加东西"（附件/@引用） | 待派 | 6/0 | `internal/tools` 受理器＋快照 | A |
| 215 | 加号菜单"技能/插件/MCP 清单" | 待派 | 5/0 | `internal/plugin`＋装配 | A |
| 216 | 菜单显示底座（洗控制字符＋两句话） | 待派 | 1/5 | 显示层纯函数 | A |
| 219 | 审批卡三枚答复按钮＋理由框 | 接单＝做 | 6/0 | 原生卡面＋审计 | D |
| 220 | 名册"卡在等批准"只看得见 L2 | 待派 | 1/4 | `internal/panel` 名册 | A |
| 226 | "以后都别问"把整份内存快照写回 config | 产码已交·六格未裁 | 1/5 | `cmd/wisp`＋`internal/config` | A |
| 228 | 球与托盘不在能干活的那条进程里 | 09-29 部分勾 | 11/5 | `cmd/wisp` 常驻＋`internal/ball` | D |
| 230 | 结案票里 4 枚真残缺 | 待派 | 3/2 | 票面/证据（零码） | A |
| 238 | 新机器常驻首跑：数据根无封条可继承 | 待派 | 4/0 | `internal/proc`＋`internal/winsec` | A |
| 239 | 球被窗口切成方形（三枚数已改） | 待派（0 未勾） | 0/5 | `internal/ball/tokens.go` | A |
| 245 | 裸 `Esc` 成全局热键抢走别处 | 已立未派 | 7/8 | `internal/ball`＋`cmd/wisp` 热键 | A |
| 246 | 会跑任务的腿不在带球的常驻进程里 | 已立未派 | 1/8 | `cmd/wisp` 装配＋`internal/proc` | A |
| 247 | 采集栈零 importer（真麦电平到不了球） | OPEN | 2/8 | `internal/audio`＋`cmd/wisp` | D |
| 248 | 面板没有"设置"这条路 | OPEN | 5/7 | `internal/panel`＋真机链 | D |
| 249 | 常驻进程无 panic 记录面 | **WITHDRAWN** | 7/0 | — | 不计 |
| 256 | 常驻腿吃不到 `[risk]` 两项配置 | 无 Status 行 | 2/2 | `internal/agent/approval` 装配 | A |
| 262 | 全仓没有尺盯"跟踪路径长度" | 无 Status 行 | 2/8 | `scripts`＋`ci.yml` | E |
| 264 | 失败句把本机绝对路径发给服务商 | 无 Status 行 | 4/0 | `internal/tools`（取舍归 owner） | C |
| 266 | `scripts/*.ps1` 门自锁钉可被就地删 | 无 Status 行 | 5/0 | `scripts` | A |
| 269 | gofmt 步被自己仓里坏样本顶成常红 | 已立未派 | 5/0 | `scripts`＋`ci.yml` | A |
| 271 | 票 221 "task 家族几枚"尺误红第三枚 | 已立未派 | 5/0 | `cmd/wisp` 测试 | A |
| 274 | 出货 exe 里有没有页面，无判据 | 无 Status 行 | 3/9 | `scripts/build.ps1` | A |
| 276 | gofumpt 普查被 17 枚台件顶成 exit 1 | 无 Status 行 | 6/0 | `scripts`＋台件归属 | A |
| 284 | `Gate.Replay` 注释没写"零产码调用者" | 无 Status 行 | 3/0 | `internal/agent/approval/gate.go` | A |
| 290 | 门装上了，没有生产路径能拧开 | 5 勾 1 未勾 | 1/5 | 球/托盘＋真机 | D |
| 291 | 电平尺不减直流（安静读到 0.37） | 1 勾 3 未勾 | 3/1 | 真麦 | D |
| 293 | 托盘"静音"勾永不跟门 | 5 勾 1 未勾 | 1/5 | 托盘真机 | D |
| 297 | 安静房间 mean≈0.3677，响/静不可分 | 未开工 | 5/0 | 真麦 | D |
| 301 | 被认领包的 windows 半边 CI 永不执行 | 无 Status 行 | 2/5 | `scripts/portable-tests.sh`＋`ci.yml` | E |
| 302 | 要页面回话的判据躺在托管 CLI 档 | 无 Status 行 | 4/1 | `scripts/portable-tests.sh` | A |
| 305 | 冷启动后"活文档"是谁（Eval 推不进） | 无 Status 行 | 5/1 | WebView2 页面 | D |
| 311 | winsec 门颜色取决于测试临时根在哪卷 | 6 勾 0 未勾 | 0/6 | `internal/winsec`＋runner 色 | E |
| 312 | 设置页自动扫本机已有模型配置 | 无 Status（按住） | 0/0（AC 为正文） | `frontend` 设置页＋Go 扫描 | B |
| 313 | 默认 stepfun＋模型/密钥全在界面配 | 无 Status 行 | 3/0 | `frontend` 设置页＋凭据 | B |
| 314 | cmd/wisp 测试二进制在真窗区整进程死亡 | 无 Status 行 | 5/0 | 真窗崩溃 | D |
| 316 | 232 的 M 形：窗口下界取在 plant 之前 | **withdrawn** | 1/0 | — | 不计 |
| 323 | 票 320 残留：不落盘明文尺零敏感 | 已立未派 | 8/0 | `internal/llm`＋`cmd/wisp/panel_models_320.go` | A |
| 329 | 快照最后一公里 `lastPanelSnapshot()` 零读者 | 已立·329-r1 | 2/7 | `cmd/wisp/run.go` publish 点 | A |
| 330 | 门铃仪器三枚缺口（镜像⇄真名册） | 已立未派 | 5/3 | `cmd/wisp` | A |
| 339 | 票 334 仪器在 Linux 档红（假家目录） | 已立未派 | 1/6 | `internal/risk` 仪器 | A |
| 342 | `privateFile` 祖先换链（两支平台都开） | 已立未派 | 2/5 | `internal/winsec` | A |
| 345 | 迁移-on-读是丢注释的根，无人认领 | 已立未派 | 7/0 | `internal/config`（需具名批准） | C |
| 346 | 两枚竞态窗口仪器 blob 未动却翻脸 | 已立未派 | 7/0 | `internal/winsec`＋`scripts` | A |
| 350 | 有界词表仪器（注释点名的包内符号） | 已立未派 | 7/0 | `internal/risk` 注释 | C |
| 354 | 复验拒删后调用者的名字仍指攻击者对象 | ready-for-agent | 6/1 | `internal/winsec` | A |
| 379 | 三枚生产 `.go` 纯 CRLF | **VOID** | 5/0 | — | 不计 |
| 388 | 测试播报音量衰减（真设备出声） | 待派 | 5/0 | `internal/audio`＋真机耳听 | D |
| 406 | C6 孤儿：POSIX 哨兵不一致（119b） | open | 0/0（新立未成框） | `internal/winsec`＋`internal/risk` | A |

---

## 4. A 清单（56 枚，本编队可做；按"能不能当波首动作"分四档）

### 4.1 零前置、可立刻起手（11）

| 号 | 面 | 未勾 | 起手第一件事 |
|---|---|---|---|
| **98** | `scripts/**`＋`docs` | 5 | 把 `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"` 写成**取数口令**＋一枚永不 skip 的守卫＋台账更正（AC#1–#4） |
| **176** | 票面＋记账 | 4 | **先判"该不该作废"**：AC#1 自写的停手条件已触发（见 §6.1） |
| **182** | 只读普查 | 6 | 10 栈逐堆"有源/无源规格要求/无源无票"+归口票（零产码） |
| **230** | 票面/证据 | 3 | AC#3 已由 `230-ac3` 交；只剩 AC#5 整包读数＋AC#4 明文"本票不接" |
| **284** | `gate.go` 注释 | 3 | 一行注释，**必须搭下一枚真动 `gate.go` 的腿**（不许独占一条腿） |
| **216** | 显示底座纯函数 | 1 | AC#5 错误码分支（其余 5 格已勾） |
| **220** | `internal/panel` 名册 | 1 | AC#5 整包终态读数一发 |
| **226** | `cmd/wisp`＋`internal/config` | 1 | AC#5 撤销路径（其余 5 格已勾、产码已交） |
| **239** | `internal/ball` tokens | 0 | 两个残余：走距 16 入 C21 名册＋`cmd/wisp` 整包复跑 |
| **406** | `internal/winsec`/`risk` | 0 | 新立：先把 AC 做成框，再判"统一哨兵 or 具名披露" |
| **122** | 全仓 Go | 7 | 第一批只清 `SA4000`1＋`SA1019`3＋`SA4006`1（AC#2 明令） |

### 4.2 可做但要**同包串行**（9）——同一波只许一枚写腿

| 号 | 面（串行键） | 未勾 | 备注 |
|---|---|---|---|
| 194 | `internal/panel`（bridge.go 常量＋guard） | 7 | owner 已批（`A381`），是 A；上一轮普查曾误分到 `internal/tools` |
| 197 | `internal/panel`＋`internal/tools` | 5 | 子代理实体，与 188 抢同一批 |
| 188 | `internal/panel`＋`internal/tools` | 5 | 与 197 同批 |
| 195 | `internal/config`＋`cmd/wisp` | 7 | 分节写入地基；`internal/config` 现量 0 命中＝真缺口 |
| 213 | `cmd/wisp` | 5 | `commands.go`＋`commands_213_test.go` 已在盘 |
| 214 | `internal/tools` | 6 | 与 176/163 抢同一批文件 |
| 215 | `internal/plugin` | 5 | 与 51 抢 |
| 51 | `internal/plugin` | 6 | C24 初始集一段属契约（AC 自陈） |
| 191 | `internal/tools`＋`internal/panel` | 6 | `PLAN:2027` 那三句不许 agent 改（AC#5） |

### 4.3 仪器/测试面（不撞共享编译面，可并行多条只读腿）（13）

| 号 | 面 | 未勾 |
|---|---|---|
| 103 | `internal/winsec`＋`internal/memory`（`R-103-1`/`R-103-2`） | 4 |
| 135 | `cmd/wisp`＋`internal/proc` 测试 | 7 |
| 136 | `internal/observe`＋`scripts` | 4 |
| 140 | `ci.yml`＋四包 | 4 |
| 145 | `internal/panel` 快照载体 | 3 |
| 180 | `internal/config`＋`cmd/wisp` 读者 | 5 |
| 185 | `internal/tools`＋`internal/risk` | 7 |
| 211 | `bridge.go`＋`budgets.go` 常量/判据 | 5 |
| 266 | `scripts/*.ps1` | 5 |
| 269 | `scripts`＋`ci.yml` | 5 |
| 271 | `cmd/wisp` 测试 | 5 |
| 302 | `scripts/portable-tests.sh` | 4 |
| 323 | `internal/llm`＋`cmd/wisp`（定向 `-run`） | 8 |

### 4.4 其余 A（23）——有前置票或需先裁形

| 号 | 面 | 未勾 | 前置/备注 |
|---|---|---|---|
| 21 | `internal/agent/approval` | 5 | 上一腿死于轮帽（`A1537`），先复判段 1 边界 |
| 23 | `internal/tools`＋C7 Image | 5 | 盘上已有未提交 `ticket23_image_channel_test.go` |
| 35 | `internal/panel` | 6 | 3 枚 cmd/wisp pump 红属 C1 |
| 47 | `internal/agent/scheduler` | 6 | PathLock 现量只在 doc/状态表 |
| 48 | `internal/agent/approval/queue.go` | 5 | `Depth`/300s 已在盘＝多为收口 |
| 49 | `internal/agent/approval`＋`perm` | 6 | `AllowSession` 已在盘（224）＝多为收口 |
| 57 | `docs`＋`internal/plugin` | 3 | 被 54 挡 |
| 71 | `cmd/balldebug`＋`d22scan`＋`ci.yml` | 6 | AC#5 需 CI 色（E 尾巴） |
| 123 | `cmd/wisp/*_test.go` | 6 | AC#4 需 runner 读数 |
| 193 | `internal/panel` 门禁＋config | 7 | AC#0 一枚 owner 亲自按过的测试 |
| 200 | `projctx`＋`agent`＋`panel` | 8 | 在派 |
| 238 | `internal/proc`＋`winsec` | 4 | 封根 |
| 245 | `internal/ball`＋`cmd/wisp` | 7 | AC#8/9 需第二进程真 Esc（D 尾巴） |
| 246 | `cmd/wisp` 装配＋`proc` | 1 | 只剩 boot flake 归属 |
| 256 | `internal/agent/approval` 装配 | 2 | AC#1 要具名 `A##`（编排者侧） |
| 274 | `scripts/build.ps1` | 3 | 三形代价表未选形（AC#2 明文"本票不选形"） |
| 276 | `scripts`＋台件归属 | 6 | 动 17 枚台件＝动证据引用面，先归口 |
| 329 | `cmd/wisp/run.go` publish 点 | 2 | 写腿已派（形甲） |
| 330 | `cmd/wisp` 门铃 | 5 | 必须排在 329 之后（共享编译面） |
| 339 | `internal/risk` 仪器 | 1 | 只剩 AC#2"不许放宽断言" |
| 342 | `internal/winsec` | 2 | 盘上已有 `ancestoropen_342_*` 两枚台件 |
| 346 | `internal/winsec`＋`scripts` | 7 | 要与 339 分面 |
| 354 | `internal/winsec` | 6 | `AC#1` 拍形已过 |

> A 桶的**真波首候选**只有 12 枚（4.1）。其余 44 枚必须按「同包串行」或「先裁形」排队——这就是上一轮 `A1514` 得出"小而自洽票已尽"的同一个事实。

---

## 5. B／C／D／E 清单

### B 前端/design 外部 agent（12）
`34`（前端脚手架）· `36`（结果/历史面板）· `37`（审批 UI）· `38`（命令面板＋任务视图）· `39`（配置编辑器 GUI）· `40`（安全/隐私/成本页）· `77`（前端脚手架·被 owner 收走）· `167`（四枚小出口：要看得见）· `169`（装饰字符把门钉红，单文件 `right-rail.tsx`）· `192`（浏览器/产物预览）· `312`（设置页自动扫本机配置）· `313`（默认 stepfun＋界面配模型/密钥）

> B 桶里**Go 半边能由本编队先做**的：`38`（task.list/cancel 工具）、`167`（快照字段）、`192`（入向方法）、`312/313`（扫描与写配置）。但**没有一枚能因此结票**——结票判据都在屏上。

### C owner 改契约文字／花钱／不可逆（17）
| 号 | 卡在哪一类 | owner 的一句话是什么 |
|---|---|---|
| 53 | 花钱/凭据 | lark-cli 真服务凭据；无凭据 e2e 不可做 |
| 56 | 花钱/外部账号 | SignPath／证书／winget 发布账号 |
| 58 | 花钱 | 英文 ASR/TTS/KWS 模型选型与下载（+frontend） |
| 60 | 花钱 | realtime 服务商 Key（票面自写"缺则 AUTO-DEFER"） |
| 100 | RESERVED | AC#0 触发门未满足 ⇒ 本票不许开工 |
| 138 | owner 批准 | `Q-55` 那一次批准一起落 |
| 148 | 双批准 | 两枚前置都是人工批准（Ⓐ 拷贝形状／Ⓑ 往外送什么） |
| 150 | 契约文字 | `Q-55` 三选一；三条出路都动冻结面或契约句 |
| 163 | 契约文字 | 扩 `shell.exec` 语义＝动 `D34`/`D46` |
| 168 | 契约文字 | `C17`/`C24` 两张白名单未定案 |
| 170 | 产品口径 | "CLI 要不要停下来等答案"不是我该拍的 |
| 174 | 契约文字 | `PLAN.md:431`"按需再读"那句（AC#4） |
| 177 | 裁形 | `Q-61` 四支候选写腿（只读核仍是 A） |
| 187 | 契约文字 | 面板→Go 写入口要写成契约变更入档 |
| 264 | owner 取舍 | 失败句里的绝对路径**抹不抹**（AC#0 明令不派写腿） |
| 345 | 具名批准 | AC#5 开工前须台账 `A##` 具名批准（文件/行/理由/边界/撤销口令） |
| 350 | owner 拍"上/不上" | AC#0 摆给机主的那一问 |

### D 真机/WebView2（21）
`25`（S3 场景①②在真 OS 上跑）· `32`（S4 听写/陪伴）· `33`（WebView2 窗口；Go 侧不冻结）· `46`（S6 电源/空闲预算）· `55`（macOS，无机器）· `62`（球视觉签收）· `64`（第二块屏）· `65`（owner 参考图＋人眼"这次对了"）· `68`（桌面复测）· `201`（球/托盘真能答一张卡）· `219`（三枚按钮在真卡上）· `228`（球/托盘进常驻进程）· `247`（真麦电平）· `248`（面板→设置→凭据→模型全链真机）· `290`（默认档双击启动读数）· `291`（真麦该读多少）· `293`（托盘勾人眼）· `297`（安静房间六窗）· `305`（WebView2 活文档）· `314`（真窗区进程死亡）· `388`（真设备出声／耳听）

### E 需真 CI run（8）
`54`（S7 红队＋自用一周）· `70`（5 job 逐因）· `72`（两侧逐字同命令）· `85`（step 日志行）· `112`（AC#5＝该步在 CI 给出过 success）· `262`（AC#5ⓑ 本机 self-hosted 侧真跑过）· `301`（`--scope=census` 分档后的色）· `311`（`Q-311-r1b-1` 要 `cmd/wisp` 的色）

---

## 6. 红标：票面**前提**已与盘上不符（本程现量，不是推理）

> 这一节是本次分诊**最值钱的产出**：下面几枚票的"立案前提"在 HEAD 上已经不成立或已变形。⛔ 我没有替它们结案（裁决者≠实现者），只把尺和读数摆出来。

### 6.1 `176` — **它自己写的停手条件已经触发**

| 票面主张 | HEAD 现量（`grep` 原命令，排测试） | 结论 |
|---|---|---|
| `RunAsync` 生产**零**调用点 | `grep -rn "RunAsync(" --include=*.go internal/ cmd/ \| grep -v _test` ⇒ **2 枚**：`internal/tools/subagent_197.go:372`、`cmd/wisp/run.go:1729` | ≥1 |
| `TaskRoster.Record` 生产**零**写者 | 同上 ⇒ **4 枚**：`internal/tools/task.go:343`、`:356`、`task_backfill.go:141`、`subagent_197.go:462` | ≥1 |

票面 `AC#1` 逐字：「⚠ 若任一枚变成"≥1"，**本票作废、回来改票 164/175 的条件**，不许留着。」
⇒ **本票应作废（或按 AC#1 改写 164/175 的条件），不是再派一条落地腿。**

### 6.2 `180` — "生产零读者"已不成立

`grep -rn "Panel.Width" --include=*.go internal/ cmd/ | grep -v _test` ⇒ `cmd/wisp/panel_resident_windows.go:208` 真读 `cfg.Panel.Width`，且 `cmd/wisp/config_readers_255.go:245` 已把它登记成 `hotClaimOtherProcess`。
⇒ 剩的 5 格是"补落地＋反向判据"，**不是"零读者"那个立票理由**。

### 6.3 `98` — "在这台机器上永远测不到"只对一半

本程现跑（**唯一一枚我跑过的测试**，只为验证 C2 的当前色）：

```bash
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go test -count=1 -run "TestTicket255RosterStillMatchesTheActualReadSites" -v ./cmd/wisp/
# --- PASS: TestTicket255RosterStillMatchesTheActualReadSites (1.45s)   14/14 section match=true
# ok  github.com/CarlosShao/wisp/cmd/wisp  1.563s
```

⇒ `cmd/wisp` **能跑**，前提是**那枚 PATH 前缀**。票面"加载期缺 DLL"是对的，"永远测不到"是**没带前缀**的读法。这就是 §4.1 把 `98` 排在第一位的原因：它把"81 枚红可见"变成一条一行命令的事。

### 6.4 `239`／`311`／`312`／`406` — 零未勾但仍开放

四枚票的 `grep -c '^[[:space:]]*- \[ \]'` ＝ **0**，但都**没有** `-done`：`239`（2 残余，票面自写"不结案"）、`311`（3 残余）、`312`（AC 是正文，无框）、`406`（新立，未成框）。
⇒ 它们不是"可做的活"，是**收口/立案动作**；放进 A 桶只是因为动作本身不需 owner/真机/CI。

### 6.5 21 枚 S 切片票是 09-19 的**模板票**，从未被认领

`33 34 36 37 38 39 40 46 47 48 49 51 53 54 55 56 57 58 60 25 32`（+`23 21` 部分）全部 `Claimed by: —`、`Last update: 2026-09-19`、Progress log **零条**，票身是 `What to build / Key constraints / AC` 模板，而 `blocked by` 里点名的票**大多今天已 `-done`**。
⇒ 它们的 AC 是按 09-19 那棵**还不存在的仓**写的：直接派写手＝按一张两万年前的图纸施工。**先做一轮"票面对拉"（半小时级）比派写腿值钱**——这条与 §6.1 同族：本仓最大的浪费不是"没做完"，是"按过期前提做完"。

---

## 7. DONE 票的孤儿红盘点（源：`docs/evidence/s1/98-attribution.md`）

### 7.1 落在哪个包（用户问的那一刀）

| 包 | 枚 | 说明 |
|---|---|---|
| **`cmd/wisp/**`** | **81** | 全部是 `cmd/wisp/*_test.go` 里的用例；81/81 逐行都是 `cmd/wisp/…` |
| `internal/**` | **0** | 该件 81 行里**零枚**落在 `internal/` |
| `frontend/**` | **0** | 同上，零枚 |
| 合计 | **81** | `PASS=143 / FAIL=81 / SKIP=0`，顶层用例 224 |

⇒ **「孤儿红在哪个包」的答案是唯一的一格：`cmd/wisp` 这一个包。** 这也解释了为什么票 `98`（让这个包在这台机上跑得起来）与票 `123`（runner 上 4 枚用例假设有真人审批）一直是所有红的共同上游。

### 7.2 按归属票 × 票的现状（81 枚全覆盖；枚数相加＝81）

| 归属票 | 枚 | 票现状 | 类 | 孤儿？ |
|---|---|---|---|---|
| 201 | 9 | **open** | C1 | 有主 |
| 223 | 8 | `-done` | C1 | **孤儿** |
| 63 | 8 | `-done` | C1 | **孤儿** |
| 12 | 6 | `-done` | C1 | **孤儿** |
| 255 | 5 | `-done` | 4×C1＋1×**C2** | **孤儿**（C2 已绿，见 7.3） |
| 101 | 5 | `-done` | C1 | **孤儿** |
| 197 | 5 | **open** | C1 | 有主 |
| 233 | 3 | `-done` | C1 | **孤儿** |
| 224 | 3 | `-done` | C1 | **孤儿** |
| 315 | 3 | `-done` | C1 | **孤儿** |
| 35 | 3 | **open** | C1 | 有主 |
| 200 | 2 | **open** | C1 | 有主 |
| 248 | 2 | **open** | C1 | 有主 |
| 165 | 2 | **open** | C1 | 有主 |
| 370 | 2 | `-done` | C1 | **孤儿** |
| 174 | 2 | **open** | C1 | 有主 |
| 151 | 2 | `-done` | C1 | **孤儿** |
| 365 | 2 | `-done` | C1 | **孤儿** |
| 226 | 1 | **open** | C1 | 有主 |
| 227 | 1 | `-done` | C1 | **孤儿** |
| 33 | 1 | **open** | **C5（假红）** | 有主（不计红） |
| 133 | 1 | `-done` | C3 | **孤儿** |
| 131 | 1 | `-done` | C4（故意红） | **孤儿**（有据） |
| 220 | 1 | **open** | C1 | 有主 |
| 329 | 1 | **open** | C1 | 有主 |
| 119 | 1 | `-done` | **C6** | **孤儿**（→ 已由票 406 认领） |
| 378 | 1 | `-done` | C1 | **孤儿** |

**合计**：归 `-done` 票的 **52 枚**（＝孤儿）／归开放票的 **29 枚**（含 1 枚假红 C5）。
**六类根因**（该件 §2 原样）：C1 非 Windows 上 DPAPI 桩 **76** ／ C2 名册漂移 **1** ／ C3 门把"没编进来"读成"没有" **1** ／ C4 故意的红 **1** ／ C5 不是红（仪器没 `.git`）**1** ／ C6 POSIX 哨兵不同名 **1**。

### 7.3 该件之后盘上发生的三笔（本程现读，决定"孤儿"的今天数目）

| 笔 | 事实 | 影响 |
|---|---|---|
| `A1539`（11-05 13:0x） | **255 的 C2 修绿** | 我现跑复认：`TestTicket255RosterStillMatchesTheActualReadSites` **PASS**、14/14 `match=true` ⇒ 孤儿红 **52 → 51** |
| `A1539` | 133 的 C3 **披露化**（转成具名披露、不判红） | 该枚不再以"红"计 ⇒ 孤儿红 **51 → 50** |
| `A1539` + 票 `406` | 119 的 C6 撞出**孤儿**并新立票 406（本程 `git status` 现读 `406` 已 `A`） | C6 **有主了**（但票 406 自己还没有 AC 框） |

⇒ **今天的账面**：孤儿红 **50 枚**（52 − C2 − C3），全部在 `cmd/wisp`；其中 **76 枚 C1 那一族**的修法两条路（加 `//go:build windows`／改名 `_windows_test.go`，或把夹具换 `env:` ref）**都还要用例主人点头**（该件 §7 逐字：两条路代价不同，"本件只负责点出是谁的"）。
⚠ 该件的读数钉在 `f4ad34ee1`（第二发 `a980734f0`），**不是 HEAD 3ff1b1395**；引用前按它 §6 的 `git archive` 口令重跑。

### 7.5 本程写件期间盘上**正在发生**的事（⛔ 不是我的活；写在这里是为了不让下一位按我这儿的数去开工）

`git status` 从 13:05 起变了三次。13:16 现读：

```
 M cmd/wisp/approval_reply_201_test.go
 M cmd/wisp/run_mode101_test.go
 M cmd/wisp/zz315a1probe_test.go
 M internal/agent/reminder_timer.go
 M internal/proc/treemetrics_windows.go
?? cmd/wisp/secret_fixture_180_test.go
```

`cmd/wisp/secret_fixture_180_test.go` 的文件头逐字写着：**"票 98 C1（orphan-red repair, ticket 180 leg）… The red this closes: docs/evidence/s1/98-attribution.md §2 C1 (76 of that census' 81 ubuntu reds)"**，且明写**不**加 `//go:build windows`／**不**改名（那会把用例移出 `--scope=cli` 的分母，正是票 98 AC#2 禁的捷径）。
⇒ **§7 的 76 枚 C1 此刻正被另一条腿按"portable seam"改**（已改 3 枚用例＋1 枚共享夹具）。**我这儿的 52/50 是 13:1x 的快照**，不是终值。⇒ 谁要引用孤儿红数量，先 `git log --oneline -3` 并重跑 `98-attribution` §6 的口令。

---

## 8. 前五枚候选腿（按"买到的价值"排序）

| # | 票 | 桶 | 面 | 买到什么 | 为什么排这个位次 |
|---|---|---|---|---|---|
| **1** | **98 残余** | A | `scripts/**`＋`docs`（零共享编译面） | 把 `cmd/wisp` 的 **81 枚红**从"永远看不见"变成"一行命令看得见"（本程已实测 1.45s PASS 一枚） | 它是**其它所有腿的取证前提**：今天任何腿交"某格红了/绿了"都绕不过这个包。成本＝一个脚本＋一段文档＋一枚守卫。★**13:16 现读：C1 那 76 枚已被另一条腿（`secret_fixture_180_test.go`，自称 ticket 180 leg）开始改**⇒ 本票剩下的正是"让它留下"的那半：口令进 `AGENTS`/派单模板＋一枚永不 skip 的守卫＋台账更正（§6.3、§7.5） |
| **2** | **176 作废/改票** | A | 票面＋`docs/evidence`（零产码） | 止住一条**注定白跑**的落地腿；同时把 164／175／154 三张票"功能在册、没人喂"的记账改正 | 票面自己写了停手条件，而条件已触发（§6.1）。**最便宜的一枚**，不做的代价是下一波派错腿 |
| **3** | **195 分节写入地基** | A | `internal/config`＋`cmd/wisp` | 面板"能改配置"的唯一地基；`internal/config` 现量 **0 命中**＝真缺口 | owner 头号产品要求（设置页）的**唯一 Go 侧前置**；解锁 248/313 的 Go 半。比 187/194 都更靠根 |
| **4** | **342 祖先换链** | A | `internal/winsec`（两支平台都开） | 关掉"`privateFile` open 沿被换掉的祖先走"；盘上**已有** `ancestoropen_342_{windows,unix}_test.go` 两枚台件可复用 | 安全面里**半条腿已铺**的一枚：台件在、只缺收口与成对读数。与 354（同一性复验）同族，可作一波 |
| **5** | **220＋226（各 1 格）＋284 随批** | A | `panel` / `config` / `approval`（三面互斥） | 三张票的**结案**，同波并行不撞面 | 单格成本最低、确定性最高；`284` 的票面自己要求"随下一枚真动 `gate.go` 的腿同批"，正好搭 226 |

> 若**只批准一枚**：派 **98**。它不改任何契约、不碰 `frontend`、不需要真机与 CI 色，产出是一把**别人也能复跑的尺**。

---

## 9. 不计入 114 的三枚（逐枚给理由，⛔ 悄悄丢掉）

| 号 | 票面原样 | 处置理由 |
|---|---|---|
| 249 | `Status: OPEN … → WITHDRAWN（10-01 14:2x 编排者自撤）` | Status 含 WITHDRAWN ⇒ 按要求排除 |
| 316 | `Status：withdrawn（撤票；口令「316 撤」已用）`（与票 272 重复） | 同上 |
| 379 | `Status: VOID（10-24 18:3x 编排者自纠：立案前提是假的——HEAD blob 的 CR=0）` | 不是 withdrawn 字面，但**前提已被自纠为假**（`autocrlf` 显示），与 316 同性 ⇒ 不计入开放池；文件名按规矩不删 |

---

## 10. 本程**没有**建立的事（宁缺毋造）

- **没有跑** `go test ./cmd/wisp/` 整包（只跑了 1 枚定向用例，为复核 C2 的当前色）；81 枚红在 HEAD 上的**分布**仍引 `98-attribution.md` 的 `f4ad34ee1` 读数，**未在 HEAD 逐枚复跑**。
- **没有跑** `sh scripts/d22scan.sh`、`bash scripts/portable-tests.sh`、`probes/154/gate-clauses.sh` 三把门禁——本件是分诊，不是交件自证。
- **没有裁定**任何一枚票的 PASS/FAIL；§6 的那几枚只给"票面前提 ⇄ 盘上读数"的对拉，**结案权在编排者与非实现者裁决腿**。
- **没有验证** B 桶 12 枚在 `frontend/**` 的当前实现程度（本编队零读零写面，只按票面判据归类）。
- **没有取** CI run id / job id：E 桶 8 枚的"真 CI 色"全部是**票面自写的结案条件**，不是我量到的读数。
- 盘上现状：13:05 起手的四行未跟踪件（`cmd/wisp/zz_probe_temp_test.go.txt`、`docs/evidence/s1/{351-ac2-ac3-verdict,378-r1-impl}.md`、`internal/tools/ticket23_image_channel_test.go`）＋13:16 起的六行（见 §7.5）**都不是我写的**；本程只写本文件，零 `git` 写操作。
