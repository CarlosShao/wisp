# 票 268 — 常驻（GUI）腿把"配置文件不存在"和"配置文件存在但被值域门拒了"**折成同一枚 provenance**：`defaults (config.toml unreadable)` ⇒ 用户写了个数、系统按 300 s 跑，**且不说**

- Status: **done（2026-10-05 22:5x 编排者翻勾＋改名，凭据见下面第 5 节；原句为「ready-for-agent，但排在其也」，写面＝`cmd/wisp/resident_approval_windows.go` ＋同包 `_test.go`，⛔ 不许与 `267-r2`／`257-r2` 并发——三枚都落 `cmd/wisp`）**
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

- [x] **AC#0（只读普查，第一格）**：现量答四问——① `riskProvenanceUnreadable` 的**全部消费者**（谁读它、到不到用户眼睛，逐枚 `file:line`）；② `config.LoadFile` 返回的错误里，"文件不存在"与"语义校验拒载"**今天能不能机读地分开**（`os.IsNotExist`／`errors.Is`／错误分类 `observe.ClassConfig`——量到不能分开就具名说不能，⛔ 不许靠字符串匹配蒙）；③ 常驻腿有没有**任何一条不过脱敏流水线的用户可见面**（`cmd/wisp/logsink.go` 那台 tee 的 mirror 那一支）；④ `wisp doctor` 那 28 枚检查项（尺＝`grep -n 'pass("\|fail("\|info("' cmd/wisp/doctor.go`）里有没有一处会读配置加载错误，加一条会不会撞既有名册钉。**⛔ 禁跑任何 `go` 命令**（与在飞的写腿抢读数）。**〔11:3x 编排者翻勾：凭据＝非实现者复核腿 `268-a2`（`.scratch/wisp/probes/268/a2/verdict.md`，205 行／26,456 字节／占位 0，首笔 `3f0c4fff`→终笔 `cf3f1ecd`，零 `go` 命令）——四问逐问复跑判定＝①②③复现（③那句"只有 primary 一支过得去"说过死，终端起法下 stderr mirror 转瞬可见）、④两口径复现（调用点 28／去具名 23）而"实印 14–16"不复现（静态枚举为 11 或 13）；★测试读者名册 census 点了四组、它精读为七组（`256_windows_test.go:124/125、149-151、165-167、174-176、247-249、317-319`，我 11:3x 现跑 `grep -n -i provenance` 复认七组成立）⇒ 这直接改 AC#2 的写法：字段断言不止在回落用例里，种子合法值支（`:247`）与 construction-time-only 支（`:317`）也在读这枚字段；★`err` 透传那半句由它独立证真（四支链：`*observe.Error` 无 `MarshalText`/`LogValue`⇒`KindAny`、key `"err"` 对三枚词表 13/5/5 全不中、`redact.go:125` 原样返回、primary 是 `slog.NewJSONHandler` 且 `Detail` 带导出 json tag）并加严一处＝透传**绕开 512 截断**（rule-4 只管 `KindString`）；★它推翻 census 六条，其中一条**剔除**：`:33` 那枚「`DefaultL1Window`=3s 同源瑕疵」属整条误读（fallback 走 `gate.go:146` 零值支、`queue.go:116` 编译到 3s，字面与常量名双双正确，⛔ 不该随本票落地）；U4 那格它按静态尺收口＝**加一行 stdout 撞 0 枚钉**（捕获名册 11 枚文件，断言形状全是 contains 或对别的字符串／文件枚数计数）⇒ 落点①②成立、③只能 info 级（`build.ps1:167-169` 拿 `doctor` 退码当 smoke，加 `fail` 级检查项＝把带内越界那台机器的构建门打红）。⚠ 题面那句「界面上没有任何一处说"你那行被拒了"」按 `A614` §3 读满＝**盘上有、名字错、你手上没有**；票面原话一字不改，按这一条读。账 `A617`〕**
- [x] **AC#1**：被拒那一支必须有**一枚与"文件不存在"不同的具名状态**，且这枚状态**能到达用户的眼睛至少一次**（日志具名 provenance 算最低档；`doctor`／回执算达标）。拿一发"把两枚状态又折回同一枚字符串"的突变证它**必响**。
- [x] **AC#2**：常驻腿回落仍然只能用**编译常量那一枚 300 s**，⛔ 不许因为这格就"顺手"让越界值生效（那是 `Q-77`，待机主）。凭据＝既有钉 `resident_approval_risk_256_windows_test.go:141-144` 保持绿、且新增一发"文件存在但被拒 ⇒ 仍是 300 s 且 provenance 具名"。
- [x] **AC#3**：跑任务那条腿的**退码 2 响亮拒绝不许改**（那是票 267 已经翻勾的形状）；本票只处理常驻腿这一支的沉默。⛔ 两腿不许被"顺手统一"成同一形（一响一静是**有名字的差别**，统一要另立票并先落 `A##`）。
- [x] **AC#4**：禁区——`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/agent/approval/**`／`frontend/**`／`design/**` 零字节；`[risk].l1_window_sec` 与它的钳位一字不动；值域 `[31, 3600]` 一字不动。⛔ **不许新增入向方法名**（那属 C17 白名单＝契约面，先例＝`Q-76` 待机主）。卫生四门读数不扩大，红名集合逐名比对并**带取数时刻**。

## 排程与禁区

- **按在 `267-r2`（种子迁移）之后**：那枚腿正在把 `cmd/wisp/**` 的 12 枚带外种子抬进带内并跑整包；在它交件＋编排者自己复跑 `cmd/wisp` 之前，本票起飞＝读数互相洗。
- ⚠ 本票与 `257-r2`（firstrun 文案那半）同落 `cmd/wisp`，⛔ 三枚串行，不接受"改的是不同文件"这种推理。
- ⚠ 如果 AC#0 量到"被拒那一支要让用户看见就必须给页面加一枚承接位"⇒ **停手上报**（那要动 `frontend/**`，本编队不写前端；只把要的字句写进票面，由机主带去他用的那枚 agent）。

## Progress log (append-only, newest last)

- [2026-10-05T05:04:53Z] agent=268-r1 did=起手：票面＋a2 复核件全文读完，编排者裁形（A＋B，⛔ 不取 C）已接；§0 起手锚逐枚亲读完成（`:434/:457/:454-456/:379-384`、`resident_windows.go:187`、`config_reload.go:396`、`256_windows_test.go:141-144/:560/:563/:310-314` 行号**未漂**）；证据件骨架 `.scratch/wisp/probes/268/r1/evidence.md` 落盘（八节，逐节"未判"占位）；`go.exe` 进程数起跑前＝0。next=整包 `cmd/wisp` 今日绿名册（baseline-cmdwisp.log，约 8.5 分钟），取完再动笔
- [2026-10-05T05:21:38Z] agent=268-r1 did=**整包基线＋落地两枚文件**。基线（动笔前，规矩 64）：`go test -count=1 -v ./cmd/wisp/`＝**RUN 331／PASS 331／FAIL 0／SKIP 0、`ok 466.320s`、rc=0**，起 13:05:56 止 13:13:47，落 `.scratch/wisp/probes/268/r1/baseline-cmdwisp.log`（1,708 行）；编排者点名的四枚带载件今天全绿一枚不跳；起跑前 `go.exe` 进程数＝0。落地＝形 A＋B：`residentRiskGateValues` 用 `errors.Is(err, fs.ErrNotExist)`（⛔ 不靠文案子串）把缺失／被拒分成两支，被拒支新常量 `riskProvenanceRefusedAtLoad = "defaults (config.toml present but refused at load)"` ＋ 照抄 `[hotkey]`（`resident_windows.go:187`）形状的一行 stdout（复用已在手的 `err`，⛔ 不新增第二枚 `[risk]` 读取点）；缺失支的 `slog.Warn` 与 `riskProvenanceUnreadable` **逐字未动**；⛔ doctor 零改动（连 info 级也没加，理由具名在证据件 §2.2）。新用例三枚 `cmd/wisp/resident_approval_risk_268_windows_test.go`（具名／300 s＋stdout 恰一行＋AST 钉"只用哨兵"），现跑 `-run 'TestTicket268|TestTicket256'` 全绿（13:17:08，256 那五枚钉含 `:141-144` 保持绿）。next=门禁四数→突变名册（MUT-1 折回同名必须红）→整包复跑逐名比对→§4-§7 填实
- [2026-10-05T07:35:35Z] agent=268-r1 did=**三发突变自证＋四门读数**。MUT-1（把新常量折回旧字面 `"defaults (config.toml unreadable)"`，13:34:05→13:36:38）⇒ **只 ① 红**，三句原文抄进证据件 §5；MUT-2（摘掉那行 `fmt.Printf`，13:52:04→14:05:08）⇒ ②③ 双红；MUT-3（哨兵旁边加一次文案匹配、行为不变，14:10:16→14:53:26）⇒ **只 ③ 红**＝"不靠文案分类"那枚钉不是装饰。三发还原后 `md5sum cmd/wisp/resident_approval_windows.go`＝`git cat-file blob HEAD:…`＝`69a630bf350daf1032c62b921e9d7044`、`git diff` 0 行；⚠ 一处自曝：MUT-3 第一发还原只回了 if 行、留下注释块，被 md5 尺当场抓到（`8fcd72af…` 不符），15:03:35 以许可形补正。**交付态整包复跑**：RUN 334／PASS 334／FAIL 0／SKIP 0、`ok 401.239s`（止 15:12:06），名册逐名差集＝只多我三枚、一枚未少一枚未红。四门（交付态复跑）：`sh scripts/d22scan.sh` rc=0 clean（15:25:29–15:25:54）、`check-path-length-budget.sh --with-self-test` rc=0 VERDICT GREEN（15:26:32，57／57／0）、`go vet ./cmd/wisp/` rc=0（15:26:37）、`gofumpt -version`＝v0.12.0 自证后 `-l cmd/wisp`＝**只剩预存脏枚 `cmd\wisp\models.go`**（15:26:39）。⚠ 一处分母变化如实登记：d22scan `ban #8 cmd/` 的 Go 文件枚数 103→104＝我新增的那枚 `_test.go`，违规仍 0 枚、红名集合仍空。next=交付笔（证据件 §4-§7 填实＋台件入库），交回编排者裁"具名程度／分档口径"两处差
- [2026-10-05T07:55:46Z] agent=268-r1 did=**交件**：证据件 `.scratch/wisp/probes/268/r1/evidence.md` 八节全满、占位归零；AC#1–AC#4 本腿判语＝四格达标（终判归非实现者腿，AC 框一枚未碰、票面原句一字未改，只在 Progress log 追加）。写面两枚：`cmd/wisp/resident_approval_windows.go`（分叉＋新常量 `riskProvenanceRefusedAtLoad`＋一行 stdout）＋`cmd/wisp/resident_approval_risk_268_windows_test.go`（三枚用例）。⛔ doctor 零改动（连 info 级也没加）。与派单口径不同处具名三条：①新 provenance 取 class 级措辞（不冒充"值域门拒的"，凭 a2 §2②）；②d22scan 文件分母 +1 是否算"读数扩大"；③票面 AC#1 的分档原文与派单"stdout 即达标"的档位差。另 §0 量到票面 `:379-382` 与派单 `:379-384` 两处文本区间不一致（代码没漂，本腿以亲读为准）。next=编排者派非实现者的缺口审计＋对抗验收；`Q-77`、`l1_window_sec` 钳位／`confirm_timeout_sec` 拒载的不一致、票 128 AC#5 自救句三处本腿一枚未动
- [2026-10-05T08:12:01Z] agent=268-r1 did=**更正笔**（不预测自己的号）：§7.2 表第 3 行的交付笔实测回填＝`b93624d4`／16:02:29／`--name-only` 11 枚路径零产码；并补一枚读数——那一行 stdout 在交付态整包里出现 **1 枚**（`grep -c 'wisp: resident \[risk\]' deliver-cmdwisp.log`＝1，行原文带 `out of range [31, 3600]`），基线 **0 枚**。工作树对我那两枚产码路径 `git status --short` 空输出＝交付态与 `9a941965` 逐字节同形（md5 尺已证）。next=交回编排者；⛔ 本腿不再自扩射程，剩下的三格（具名程度／d22scan 分母／AC#1 分档）等裁
- [2026-10-05T08:48:38Z] agent=268-r1 did=**补笔（自曝一处交付件失手）**：上一笔 `2a0e23b4` 的 commit message 只剩一行——我在 `-F` 之前用 shell 重定向把 Write 落的详细 msgfile **覆盖了**。按纪律不 amend／不 reset／不 checkout，只追加：正文补成独立台件 `.scratch/wisp/probes/268/r1/msg-04-body.txt`，证据件 §6.11 具名登记并写死定式"msgfile 只由 Write／Edit 产出，`-F` 前先 `wc -l` 复量"；§7.2 表同时回填第 4 行实测（`2a0e23b4`／16:17:52／3 枚路径零产码）。信息未丢：那两枚读数（stdout 行在整包里出现 1 枚／基线 0 枚；两笔产码 commit 的禁区命中 0 行）本来就在 §1 与 §4.5。next=交回编排者，本腿射程到此为止
- [2026-10-05T09:01:12Z] agent=268-r1 did=**末笔（交件前最后一采）**：门 3／门 4 于 16:55:15／16:55:21 复跑（`go vet ./cmd/wisp/` rc=0、`gofumpt -l cmd/wisp`＝只剩预存脏枚 `cmd\wisp\models.go`）；同刻两枚产码 `md5sum`＝`git cat-file blob HEAD:…` 逐枚同值（`69a630bf…`／`e2405a21…`）、`git status --short -- cmd` 对本腿路径空输出 ⇒ 门 1／门 2 在 15:25:29–15:26:36 的读数即交付态读数。证据件终尺＝320 行／38,309 字节／占位 0；⚠ 一处 Edit 失手自曝：我本想追加这末采，却拿整段当锚点误删了两行 scope 读数，随即以同一段锚补回（`git diff --stat` 净读数＝只多 6 行、删除列 0 行），并再次印证本仓定式"拿整行当锚点，替换串必须逐字含旧行全文"。next=无（本腿射程用完，四格达标判语与三处待裁口径已交回）

## 5. 编排者翻勾记录（2026-10-05 22:5x +08，锚 HEAD `e8df2643`）

**凭据＝非实现者验收腿 `268-v2`**（`.scratch/wisp/probes/268/v1/evidence.md` **459 行／62,585 字节**，六笔 `689bc241`→`ba2482c9`→`862556b3`→`0d207e14`→`f73558c0`→`e8df2643`；实现腿＝`268-r1`，两腿无交集，D22 双角色成立）。编排者另独立复跑过三把尺：①名册差集（顶层 `--- PASS` 238 枚 vs 基线 235 枚，多的恰为本票三枚新用例、少的 0 枚）；②两把壳尺 21:30 rc=0（d22scan clean、path-length VERDICT GREEN、tracked=6026／over=57／roster=57／not-in-roster=0）；③**安全尺两采**（22:45：被验两文件 `md5sum` 与 `git cat-file blob HEAD:` 逐枚同值＝`69a630bf…`／`e2405a21…`，scoped porcelain 空 ⇒ **三发自重突变全部还原、盘上无残留针**）。

| 格 | 判语 | 凭据节 |
|---|---|---|
| AC#0 | 成立（11:3x 已翻，凭据＝`268-a2` verdict.md） | 见本票上面注记 |
| AC#1 | **成立（达标档）** | v2 §1.1：独立具名 `…present but refused at load`（`:456`，与 `:446` 逐字不同名）＋两面到达眼睛（日志具名 4 次＋stdout 行首尺＝1）＋★MUT-1 由 v2 亲手折回同名 ⇒ 用例① 红、红句逐字在 §2 |
| AC#2 | **成立** | v2 §1.2：新增行时长字面 0 枚、回落走 `return 0, 0, <provenance>`；票 256 那五枚钉**定向复跑逐枚点名 `--- PASS`**（台件 `v2-targeted-256.log` rc=0）；三面零放宽（交付笔对 256 仪器与 `internal/` 零字节／`[31,3600]` 逐字在／未压种子未 Skip）；★`Q-77` 仍关着，没被这格顺手放开 |
| AC#3 | **成立** | v1/v1b §5.1／§5.4 静态面：`cmd/wisp/run.go` 零字节、退码 2 那一支一字未改；两腿未"顺手统一"（一响一静是有名字的差别） |
| AC#4 | **成立** | v1/v1b §5.5＋v2 §7.2：产码笔对禁区文件名命中 0；`[risk].l1_window_sec` 钳位与值域 `[31, 3600]` 一字未动；生产 diff 新增导出方法名 0 枚（⛔ 未新增 C17 入向面）；卫生四门读数不扩大（红名集合逐名比对，`cmd/` Go 文件 103→104 的分母增长见下裁第 2 条） |

**对 `268-r1` 交回的六处差异，编排者裁定**：①**provenance 取 class 级措辞＝追认**（不读 error prose 是 `internal/config` 那层的既有纪律，细分需跨包造 marker＝契约面，另票）；②**d22scan 分母 103→104＝判"未扩大"**（AC#4 那句"红名集合逐名比对"约束的是违规名集合，不是被扫文件枚数；本票新增一枚 `_test.go` 必然抬分母，违规 0／红名 0）；③**AC#1 分档＝按派单口径**（日志具名＋stdout 行＝达标；票面原文把 doctor／回执列为达标档，此处按"终端起法"这一现实档裁，票面原句不改）；④硬钉射程＝追认 v1b §5.1 复认；⑤票面现量段写 `:379-382` 而亲读为 `:379-384`＝**票面一字未改**，以 v2 §0.4 的现读为准，属行号漂移不是产码改动；⑥触发面窄于 a2（既有 14 枚 `[risk]` 种子里带外 confirm＝0 枚）＝追认，这是交付整包 334 枚全绿无噪声的成因解释。

**残余与归口（⛔ 不由本票顺手做）**：
- ★**双击（`-H=windowsgui`）起法下那行 stdout 无处可去**——那一面只剩日志具名＝票面**最低档**。⇒ 归票 244 那一族（同一枚"黑控制台"根因），本票不碰构建面。
- **"被哪一条 loader 规则拒"分不开**（语法／未知键／迁移／值域四因今天不可机读分开）——要一名一因须给 `internal/config` 造 marker 类型＝跨包契约面。⇒ 待人项（`Q-77` 同族另说）。
- **跑任务腿与陪聊腿的回落一致性**不在本票射程（v2 §1.2 末条具名"量不到"）。
- `wisp doctor` 那档：加 `fail` 级检查项会把带越界值那台机器的构建门打红（268-a2 U4 现量），本票未动 doctor。
- ⛔ 页面承接位一格未动（`frontend/**` 归 owner 委托的会话）。

**结案动作**：本票五格全勾 ⇒ 改名 `-done`（防重领唯一键）。腿号归口：产码＝`268-r1`（`9a941965`），撞钉料＝`268-a2` census.md，AC#0 复核＝`268-a2` verdict.md（⚠ 同目录代号撞车已在 A626 记我），终裁＝`268-v1`／`268-v1b`／`268-v2` 三腿接力（前两枚死于服务故障、第三枚撞 150 轮帽但六笔全落）。
