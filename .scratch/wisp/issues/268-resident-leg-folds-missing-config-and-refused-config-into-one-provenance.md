# 票 268 — 常驻（GUI）腿把"配置文件不存在"和"配置文件存在但被值域门拒了"**折成同一枚 provenance**：`defaults (config.toml unreadable)` ⇒ 用户写了个数、系统按 300 s 跑，**且不说**

- Status: **ready-for-agent，但排在其也**（写面＝`cmd/wisp/resident_approval_windows.go` ＋同包 `_test.go`，⛔ 不许与 `267-r2`／`257-r2` 并发——三枚都落 `cmd/wisp`）
- 来源：票 267 的落地链。`267-a2` 普查件 §1.1★／§7.3 抓到、编排者 10-05 09:5x 复读代码确认，落账 `A611` §6。撤销口令「**268 撤**」。
- 关联：票 267（值域门 `[31,3600]` 是这枚洞的**制造者**：门之前越界值会真的流进常驻腿的门，之后被静默折走）· 票 256（常驻腿读 `[risk]` 两枚字段那批钉）· 票 248（回执文案"这两项只作用于跑任务的进程"那一寸）· 票 255（"配了不生效"同族账）

## 这是什么（人话）

球／托盘那条常驻进程也要读 `config.toml` 里"审批卡开多久"这个数。它读的方式和跑任务那条不一样：

- **跑任务的腿**：读不到就**退出来报错**（退码 2，屏幕上写明"配置未就绪"）＝响亮。
- **常驻腿**：读不到就**按内置的 300 秒继续跑**，并在内部记一笔来路 `defaults (config.toml unreadable)`＝安静。

票 267 之前，用户写 99999 或 2 秒，常驻腿**照数跑**（虽然跑得不健康）。票 267 之后，这些数在**加载那一层就被拒** ⇒ 于是常驻腿现在遇到的是第三种情况：**文件明明在、里面明明写了、但被门的规则拒了**——而它把这第三种**当成"文件读不到"**处理。

⇒ 最坏后果的形状：用户在 `config.toml` 里把审批超时写成 20 秒（比如他嫌 300 秒太久），**系统按 300 秒跑，界面上没有任何一处说"你那行被拒了"**。这不是安全问题（300 > 30 ⇒ C18 那句"即将超时"提示仍然武装），**是"说实话"问题**：写下去的数与真正生效的数分家，而且没人承认分家。

## 现量（⚠ 引用前先重跑，别把这几行当常量）

1. 折叠点在 `cmd/wisp/resident_approval_windows.go:457`——`residentRiskGateValues` 走 `config.LoadFile`，**任何**错误（含 band 拒载）都 `return 0, 0, riskProvenanceUnreadable`；`:434` 逐字 `riskProvenanceUnreadable = "defaults (config.toml unreadable)"`。⇒ 唯一那条把来源交给用户的日志在 `:379-382`：`slog.Info("resident gate: [risk] tier taken at construction", "provenance", provenance, "config_path", …, "window_sec_read", …, "confirm_timeout_sec_read", …)`——**被拒与缺失在这一条里长同一枚样**。编排者复跑于 2026-10-05 09:5x。
2. 分不开这件事**今天没有仪器**：`cmd/wisp/resident_approval_risk_256_windows_test.go:141-144` 断的就是"读不到 ⇒ `approval.DefaultApprovalTimeout`（300 s）"，注释逐字 `(300s is the contract default; a fourth number here means the fallback grew a value of its own)`——它守的是"回落值不许自己长出新数"，**不守**"两种读不到不许折叠"。
3. 对照面（跑任务的腿那半，作为本票的"应当长这样"参照）：`cmd/wisp` 的集成用例在 band 之后逐字吐 `wisp run: 配置未就绪（Unconfigured）：config: config.toml: risk.confirm_timeout_sec 20 out of range [31, 3600]`（台件 `.scratch/wisp/probes/267/gate/cmdwisp-HEAD.log`，1,575 行，16 枚 FAIL 全靠这一句归因）。

## 要建什么（⛔ 形不在票面写死，归 AC#0 的普查）

只一件事：**把"没有配置文件"与"配置文件里的值被规则拒了"这两种状态分开，并让被拒那一支对用户可见**。可见到什么程度＝AC#0 量完再定形（候选：日志里换一枚具名 provenance＋`wisp doctor` 里加一条检查项／常驻腿启动时把拒载句打到那唯一一条不过脱敏流水线的 stderr mirror／面板回执里加一句——第三支要先查页面有没有承接位，**没有就不许造**）。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [ ] **AC#0（只读普查，第一格）**：现量答四问——① `riskProvenanceUnreadable` 的**全部消费者**（谁读它、到不到用户眼睛，逐枚 `file:line`）；② `config.LoadFile` 返回的错误里，"文件不存在"与"语义校验拒载"**今天能不能机读地分开**（`os.IsNotExist`／`errors.Is`／错误分类 `observe.ClassConfig`——量到不能分开就具名说不能，⛔ 不许靠字符串匹配蒙）；③ 常驻腿有没有**任何一条不过脱敏流水线的用户可见面**（`cmd/wisp/logsink.go` 那台 tee 的 mirror 那一支）；④ `wisp doctor` 那 28 枚检查项（尺＝`grep -n 'pass("\|fail("\|info("' cmd/wisp/doctor.go`）里有没有一处会读配置加载错误，加一条会不会撞既有名册钉。**⛔ 禁跑任何 `go` 命令**（与在飞的写腿抢读数）。**〔11:3x 编排者翻勾：凭据＝非实现者复核腿 `268-a2`（`.scratch/wisp/probes/268/a2/verdict.md`，205 行／26,456 字节／占位 0，首笔 `3f0c4fff`→终笔 `cf3f1ecd`，零 `go` 命令）——四问逐问复跑判定＝①②③复现（③那句"只有 primary 一支过得去"说过死，终端起法下 stderr mirror 转瞬可见）、④两口径复现（调用点 28／去具名 23）而"实印 14–16"不复现（静态枚举为 11 或 13）；★测试读者名册 census 点了四组、它精读为七组（`256_windows_test.go:124/125、149-151、165-167、174-176、247-249、317-319`，我 11:3x 现跑 `grep -n -i provenance` 复认七组成立）⇒ 这直接改 AC#2 的写法：字段断言不止在回落用例里，种子合法值支（`:247`）与 construction-time-only 支（`:317`）也在读这枚字段；★`err` 透传那半句由它独立证真（四支链：`*observe.Error` 无 `MarshalText`/`LogValue`⇒`KindAny`、key `"err"` 对三枚词表 13/5/5 全不中、`redact.go:125` 原样返回、primary 是 `slog.NewJSONHandler` 且 `Detail` 带导出 json tag）并加严一处＝透传**绕开 512 截断**（rule-4 只管 `KindString`）；★它推翻 census 六条，其中一条**剔除**：`:33` 那枚「`DefaultL1Window`=3s 同源瑕疵」属整条误读（fallback 走 `gate.go:146` 零值支、`queue.go:116` 编译到 3s，字面与常量名双双正确，⛔ 不该随本票落地）；U4 那格它按静态尺收口＝**加一行 stdout 撞 0 枚钉**（捕获名册 11 枚文件，断言形状全是 contains 或对别的字符串／文件枚数计数）⇒ 落点①②成立、③只能 info 级（`build.ps1:167-169` 拿 `doctor` 退码当 smoke，加 `fail` 级检查项＝把带内越界那台机器的构建门打红）。⚠ 题面那句「界面上没有任何一处说"你那行被拒了"」按 `A614` §3 读满＝**盘上有、名字错、你手上没有**；票面原话一字不改，按这一条读。账 `A617`〕**
- [ ] **AC#1**：被拒那一支必须有**一枚与"文件不存在"不同的具名状态**，且这枚状态**能到达用户的眼睛至少一次**（日志具名 provenance 算最低档；`doctor`／回执算达标）。拿一发"把两枚状态又折回同一枚字符串"的突变证它**必响**。
- [ ] **AC#2**：常驻腿回落仍然只能用**编译常量那一枚 300 s**，⛔ 不许因为这格就"顺手"让越界值生效（那是 `Q-77`，待机主）。凭据＝既有钉 `resident_approval_risk_256_windows_test.go:141-144` 保持绿、且新增一发"文件存在但被拒 ⇒ 仍是 300 s 且 provenance 具名"。
- [ ] **AC#3**：跑任务那条腿的**退码 2 响亮拒绝不许改**（那是票 267 已经翻勾的形状）；本票只处理常驻腿这一支的沉默。⛔ 两腿不许被"顺手统一"成同一形（一响一静是**有名字的差别**，统一要另立票并先落 `A##`）。
- [ ] **AC#4**：禁区——`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/agent/approval/**`／`frontend/**`／`design/**` 零字节；`[risk].l1_window_sec` 与它的钳位一字不动；值域 `[31, 3600]` 一字不动。⛔ **不许新增入向方法名**（那属 C17 白名单＝契约面，先例＝`Q-76` 待机主）。卫生四门读数不扩大，红名集合逐名比对并**带取数时刻**。

## 排程与禁区

- **按在 `267-r2`（种子迁移）之后**：那枚腿正在把 `cmd/wisp/**` 的 12 枚带外种子抬进带内并跑整包；在它交件＋编排者自己复跑 `cmd/wisp` 之前，本票起飞＝读数互相洗。
- ⚠ 本票与 `257-r2`（firstrun 文案那半）同落 `cmd/wisp`，⛔ 三枚串行，不接受"改的是不同文件"这种推理。
- ⚠ 如果 AC#0 量到"被拒那一支要让用户看见就必须给页面加一枚承接位"⇒ **停手上报**（那要动 `frontend/**`，本编队不写前端；只把要的字句写进票面，由机主带去他用的那枚 agent）。

## Progress log (append-only, newest last)

- [2026-10-05T05:04:53Z] agent=268-r1 did=起手：票面＋a2 复核件全文读完，编排者裁形（A＋B，⛔ 不取 C）已接；§0 起手锚逐枚亲读完成（`:434/:457/:454-456/:379-384`、`resident_windows.go:187`、`config_reload.go:396`、`256_windows_test.go:141-144/:560/:563/:310-314` 行号**未漂**）；证据件骨架 `.scratch/wisp/probes/268/r1/evidence.md` 落盘（八节，逐节"未判"占位）；`go.exe` 进程数起跑前＝0。next=整包 `cmd/wisp` 今日绿名册（baseline-cmdwisp.log，约 8.5 分钟），取完再动笔
