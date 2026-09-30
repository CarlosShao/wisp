# 245-r1 — 产码腿：球进常驻进程后裸 `Esc` 被抢成全局热键 ⇒ 走乙（稳态不绑 cancel，只在 `Confirming` 借、离开即还）

- 腿：`245-r1`（**产码腿**，不是裁决腿）。工单：`.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md`。
- **本表性质＝实现者自述**。填表的人＝写这批码的人，所以它**不是** `SPEC-12 §4.3` #1/#3 要的那份"非实现者对抗验收"，也**不勾任何 AC 框**（AC 框归编排者，本腿一枚没动票面勾选状态：票内 `^- [ ]`／`^- [x]` 计数见 §6 最后一行）。
- 骨架先交于第 5 轮内（commit `f91feddd`），本节起是逐节填完的版本。硬预算闸门（第 100 轮先写满 §3／§5）：**本轮未触顶**，§3／§5 已随本版写满。

## §0 起手名册与结论

- 起手 HEAD：`25ce3ce9`（编排者 245 票面 AC#3② 更正那一枚）。⚠ 我第一条命令（17:0x）读到的 HEAD 是 `cf03fa31`，**同发再看已是 `25ce3ce9`**——写面在同一分钟里被编排者推进过一枚；本表所有行号与"现读"均以 **`25ce3ce9` 为基准**，终判据见 §6。
- 起手 `git status --porcelain` 全量名册逐字存：`.scratch/wisp/probes/245/r1/start-status.txt`（**173 行**，取数时刻 09-30 17:0x +08）。共享工作树 ⇒ 本腿终态判据＝**等于这一份**，不是"必须为空"。
- 起手时本腿地界（`internal/ball`／`cmd/wisp`）的写面：**0 枚**（尺：`grep -c "internal/ball\|cmd/wisp" start-status.txt` ＝ **0**＝没找到＝那两枚包干净，这是好消息）。起手名册里唯一一枚 `docs/evidence/s1/` 条目是别人的 `152-subject-death-never-measured-r1-accept-r1.md`（` M`），**本腿未碰、未提交**。
- 起跑前的进程尺（派单第 5 节，真机用例开跑前必量）：`tasklist //FI "IMAGENAME eq balldebug.exe"`／`... wisp.exe`／`... go.exe` 三条**都是 "INFO: No tasks are running which match the specified criteria."＝计数 0 枚**（`grep -c` 那把尺零命中返回 rc=1 会断 `&&`，本腿一律用 `;` 串接）。每次跑 winlive 前复量，共量 3 次，全部 0。
- **一句话结论**：**稳态（没有卡片挂着时）注册集枚数＝3**（summon／mute／panel），`cancel` 不在场；**借用态（`Confirming` 期间）＝4**。生产默认键名一字未改（`DefaultHotkeys()` 里 `Cancel: "Esc"` 原样），只把它的作用区间从"从开机起永久绑着"改成"借—还"。
- 本腿三枚 commit：`f91feddd`（骨架）→ `e6d6426f`（乙形码＋三枚钉同批收紧）→ 本表这一枚（见 §6）。

### §0.1 具名欠账一览（本腿现测，无一条来自注释自陈；⛔ 本腿一律未修）

| # | 欠账 | 本腿读到的凭据 | 归口 |
|---|---|---|---|
| N#1 | 常驻腿整棵 `config.toml` 没接 ⇒ `[hotkey] cancel` 用户改了也不生效（AC#4 就卡这） | `cmd/wisp/resident_ball_windows.go:86` 逐字 `Hotkeys:  ball.DefaultHotkeys(),`（无参分支不读配置） | 票 228 后续片 |
| N#2 | `cmd/wisp/main.go:25-26` usage 那句 `its four global hot keys` 在本票之后是半假话（常驻腿稳态三枚，且它没有卡片路径） | 现读 `sed -n '20,35p' cmd/wisp/main.go`；无测试钉这枚串（尺：`grep -rn "global hot" cmd/wisp/*_test.go` 只命中 `resident_ball_live_228_windows_test.go:163` 那枚 4/4 计数） | 票 228 的 AC#7/AC#8 同族文案面（A477 已把 `main.go` 那句判归 228） |
| N#3 | 借 Esc 那 2–3 秒仍吞别处的 Esc（乙形残余） | 见 §4 五字段登记 | 新推迟项（SPEC-12 §5 那一行**要人工/编排者落**，`docs/specs/**` 是本腿禁地） |
| N#4 | `RebindHotkeys`（配置热重载）会丢掉一次在飞的借用：重绑＝回到稳态集，`escTakenOver=false`，没人再借回来 | 现读 `internal/ball/ball_windows.go:822`（`b.escTakenOver = false` 那行原本就在，本腿只补了注释） | 票 228 后续片（卡片路径与配置轮询真接上之后才有意义） |
| N#5 | 用户若把 `cancel` 配成带修饰键的组合（例 `Ctrl+Alt+V`），乙形下稳态同样**不绑**——票面 AC#1 写的是无条件三枚，本腿照办；这条行为变化要有人认 | 现读 `internal/ball/hotkey_windows.go:428-448`（`cancelIdleLine` 对所有可解析的 cancel 绑定一律 standby） | 待人裁（见 §5 J#2） |

## §1 判据逐格表

| 格 | 判语 | 尺与现读（逐字） | 正控／突变 |
|---|---|---|---|
| **AC#1** 稳态名册不含裸 `Esc` | **成立（本机现读，两层都有读数）** | ① 出厂进程面（最硬的一发）：`PATH=... go test -tags winlive ./cmd/wisp -run TestLive228ResidentLegOwnsABallWindowOnTheDesktop` ⇒ `resident_ball_live_228_windows_test.go:147: AC#1 LIVE: hwnd=96077818 pid=19308 owns the ball window; hotkeys live 3/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not b`（末段被该尺自己的 80 字符窗口截断，全串＝`cancel=Esc not bound while idle (cancel is borrowed only during Confirming)`）。② 包内稳态：`internal/ball/live_windows_test.go:74` 的 `len(got.Live()) != 3` ＋ `:78 requireIdleRoster` 真机绿。③ 确定性层：`TestRegisterAllLiveSet`／`TestDefaultHotkeysIdlePassHoldsNoEsc` 绿。④ 桌面面：`requireIdleRoster` 里那枚 `escBorrowProbe`（`:429`）在稳态必须回"free"，绿。**稳态 `Live()` 枚数＝3；借用态＝4。** | **M1**（把稳态那枚分支反回去＝稳态真绑裸 Esc）⇒ 上面每一把都红，读数在 §2 |
| **AC#2** `Confirming` 内借、离开还（两头发） | **成立（真机，两头发各有读数；⚠ 只到"包与调试路径"，常驻腿今天没有卡片路径）** | 借：`internal/ball/interaction_live_test.go:205 requireEscBorrowed` ⇒ 报告里 `cancel` 是 `live` 且 `Acc.VK==0x1B`，**同一发** `escBorrowProbe` 回"不 free"（＝Win32 层面桌面真的拿不到 Esc 了）。还：`:222`／`:250`／`:273` 三处 `requireEscReturned` ⇒ 撤卡后 `Live()` 回到 3、`cancel` 行回 `standby`、**探针回"free"**（＝还出去的是这台机器的 Esc，不是我自己的记账）。尺的姿势复用票 64 那一族（`sendMessage(b, wmHotkey, hkCancel, 0)` → D43 #22 → Acting，`interaction_live_test.go:211-216` 原样保留）。整包读数：`go test -tags winlive ./internal/ball -run 'TestLiveHotkey\|TestLiveConfirmingCancelAndEscReturned\|TestBallLiveLifecycle'` ＝ 4 枚 PASS、**`SKIP-LOUD` 计数 0 枚**（＝这一发连注入那支也没躲过去）。 | **M2a**（`Confirming` 根本不尝试）⇒ 红句 `entering Confirming did not take Esc over (B1)` 与 `ticket 245 RED: Confirming never attempted the Esc borrow (cancel line still standby, live set holds 3 of the 4 slots)`。**M2b**（谎报：记账说借到了、Win32 没注册）⇒ 红句 `ticket 245 RED: the ball says it borrowed Esc, but the desktop still has it free`。两发都在 §2 |
| **AC#3** 三枚钉同批改期望＋收紧后仍能红 | **成立（三枚逐枚改了，另加两枚票面没点名的同类钉一并改；⛔ 零删断言、⛔ 零放宽、⛔ 零词面型新仪器）** | ① `internal/ball/live_windows_test.go:68-70` → 现 `:74-78`：`len(got.Live()) != 4` 改成 **`!= 3`** ＋ 紧跟 `requireIdleRoster`（稳态枚数**且**"cancel 不在场"**且**"裸 Esc 桌面可得"三件一起判＝双向收紧：多一枚红、少一枚也红）。② `cmd/wisp/resident_ball_live_228_windows_test.go`：`:144` 那把只截日志的尺**一字未动**；`:151` 的 `hotkeys live 0/4` 红尺**一字未动**；`:148-150` 那句"four hot keys"注释改成条件事实句（现 `:148-162`）；**新增** `:163` `hotkeys live 4/4` ⇒ 红（这条腿稳态四枚全在场＝裸 Esc 被抢走）。文案 `%d/4` 按票面不改，变值现读：M1 下 `4/4`、修好后 `3/4`（两条逐字读数在 §0 与 §2）。③ `internal/ball/hotkey_live_test.go:33-36` 的 `liveHotkeys()` 注释改成"稳态根本不绑、借期才绑"的带条件事实句；⚠ **用例仍绑 `Ctrl+Alt+V`，一枚都没改回裸 Esc**（票面 AC#3③ 明令）。 | M1 一发同时打红 ①②③三枚（§2 逐发红句）；M2a/M2b 打红收紧后的借/还两支 |
| **AC#4** 用户能改 `[hotkey] cancel` | **阻塞于票 228**（本格照票面写，不许假绿） | 尺＝生产调用链现读：`cmd/wisp/resident_ball_windows.go:86` 逐字 `Hotkeys:  ball.DefaultHotkeys(),` ⇒ 常驻腿**根本不读 `[hotkey]`**，用户改 `config.toml` 到不了 Win32。票 228 的"常驻腿接 `config.toml`"那一格未做（`docs/reports/pending-and-issues.md` A475 末段：`[hotkey]`／`[ball]` 未接）。⛔ 本腿**没有**新增第二条配置通路。 | 不适用（没做＝没有可红的尺） |
| **AC#5** 残余五字段登记 | **成立（登记已落两处；⚠ SPEC-12 §5 那一行本身要由编排者落，本腿禁地）** | 五字段齐全见 §4.1；同一份内容抄进工单末尾新节（`.scratch/wisp/issues/245-...md` 末尾「残余登记（五字段）」）。⛔ 本腿**未**在代码里加 `DEFERRED(D-xx)` 标记——`SPEC-12 §5` 要求标记与登记表 1:1 双向对得上，而那张表在 `docs/specs/SPEC-12-roadmap-governance.md`（本腿禁地），造一枚没人登记的标记＝当场把双向核对打破。归编排者（见 §5 J#1）。 | 不适用（登记类判据，无码可突） |

## §2 突变名册（逐发：改了哪一行／指名用例／红句／还原）

还原一律 `git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet` 复认；⛔ 未用 `git checkout`、未用 `-overlay`。

| 编号 | 位置 | 突变内容（一句话） | 指名要红的用例 | 实际红句（逐字） | 还原复认 |
|---|---|---|---|---|---|
| **M1** | `internal/ball/hotkey_windows.go` `registerAllWith` 的 `if p.id == hkCancel` 分支 | 把稳态那枚分支改成永不进入（`if false && ...`）＝回到 245 之前"开窗即绑四枚（含裸 Esc）" | 默认层 `TestRegisterAllLiveSet`／`TestDefaultHotkeysIdlePassHoldsNoEsc`／`TestCancelBorrowRoundTrip`／`TestStandbyIsNotDisabledAndNotAProblem`／`TestRegisterAllSplitsFailureFamilies`；真机层 `TestBallLiveLifecycle`（钉①）／`TestLiveHotkeyRebindEndToEnd`／`TestLiveHotkeyOccupiedVsNotAttempted`／`TestLiveConfirmingCancelAndEscReturned`；进程层 `TestLive228ResidentLegOwnsABallWindowOnTheDesktop` | 默认层：`hotkey_status_test.go:120: attempted [1 2 3 4], want exactly the three idle slots [1 2 4] (cancel is standby, ticket 245)`；`hotkey_status_test.go:179: the production defaults bound a bare Esc globally: [{16387 81} {16387 77} {16384 27} {16387 80}]`（`{16384 27}`＝MOD_NOREPEAT＋VK_ESCAPE，这就是那枚被吞掉的键）；`hotkey_status_test.go:206: idle pass attempted cancel: [1 2 3 4]`；`:296: cancel standby and mute-disabled share a status (live)`；`:337: cancel while neighbours failed = "live", want standby`。真机层：`live_windows_test.go:75: the live ball did not register its three idle hotkeys (summon/mute/panel, cancel is borrowed only during Confirming)`；`hotkey_live_test.go:149: ticket 245 RED: the idle ball holds the cancel hot key`；`interaction_live_test.go:177: 同一句`。进程层：`resident_ball_live_228_windows_test.go:164: AC#1 RED (ticket 245): the idle resident ball registered all four hot key slots...` ＋ 子进程自己那一行 `hotkeys live 4/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc live, panel=Ctrl+Alt+P live` | `git diff --quiet -- internal/ball/hotkey_windows.go` rc=0；`grep -n "if false && p.id == hkCancel"` 命中 **0 枚**（＝没找到＝突变已不在盘上） |
| **M2a** | `internal/ball/ball_windows.go` `TakeEscForCancel` | 进 `Confirming` 完全不尝试借用（slot 留 standby） | `TestLiveConfirmingCancelAndEscReturned`、`TestBallLiveLifecycle` | `interaction_live_test.go:201: entering Confirming did not take Esc over (B1)`；`live_windows_test.go:127: ticket 245 RED: Confirming never attempted the Esc borrow (cancel line still standby, live set holds 3 of the 4 slots): [...cancel ID:3 Binding:Ctrl+Alt+V Status:not bound while idle (cancel is borrowed only during Confirming)...]` | `git diff --quiet -- internal/ball cmd/wisp` rc=0；`grep -n "MUTATION" internal/ball/ball_windows.go` 命中 **0 枚** |
| **M2b** | 同一枚函数 | **谎报形**：不调 Win32，但把记账与报告改成"已借到"（＝只测自己的 bool 会放过去的那一形） | 同上三枚真机用例 | `hotkey_live_test.go:311`／`interaction_live_test.go:205`／`live_windows_test.go:127` 同一句：`ticket 245 RED: the ball says it borrowed Esc, but the desktop still has it free - the borrow registered on some other key, not on VK_ESCAPE` | 同 M2a |
| **M3**（未跑） | — | 想突"归还时回头再绑配置那枚默认 Esc"（＝缺陷回魂的另一形） | 判**不需要另开一发**：M1 已经覆盖"稳态绑了裸 Esc"这一结局，而 `releaseEscWith` 只有一行 `unreg(hkCancel)`，`TestCancelBorrowRoundTrip` 里 `reg.attemptsOf(hkCancel) != 1` 那枚断言直接把"还＝再注册一次"钉死（默认层，确定性）。 | 不适用 | 不适用 |

还原后复跑（每一发都回到绿）：
- `go test -tags winlive ./internal/ball -run 'TestLiveHotkey|TestLiveConfirmingCancelAndEscReturned|TestBallLiveLifecycle'` ⇒ `ok github.com/CarlosShao/wisp/internal/ball 2.707s`，`--- FAIL` **0 枚**，`SKIP-LOUD` **0 枚**。
- `PATH=... go test -tags winlive ./cmd/wisp -run TestLive228...` ⇒ `--- PASS: TestLive228ResidentLegOwnsABallWindowOnTheDesktop (3.96s)`，`ok ... 4.010s`，读数回到 `hotkeys live 3/4`。

## §3 门禁读数（本机 09-30 17:2x–17:3x，全部现跑一次，逐字）

| 尺 | 读数 |
|---|---|
| `go build ./...` | rc=0，无输出（存：`.scratch/wisp/probes/245/r1/build.txt`＝空文件＝零诊断） |
| `go vet ./internal/ball/` | rc=0，无输出 |
| `go vet -tags winlive ./internal/ball/` | rc=0，无输出 |
| `go vet ./cmd/wisp/` | rc=0，无输出 |
| `GOOS=linux go vet ./internal/ball/` | rc=0，无输出（本腿只动 windows-tagged 文件，没搬符号） |
| `GOOS=linux go vet ./cmd/wisp/` | **不吃读数**：报 `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in ...@v1.13.8`——这条与本票无关（`cmd/wisp` 的 sherpa cgo 上游约束，票 78 地界），本腿未修、如实报 |
| `gofmt -l internal/ball/ cmd/wisp/` | 修前 1 命中（`internal/ball/hotkey_status_test.go`），`gofmt -w` 后**零命中**＝没找到＝干净。⚠ **`gofumpt` 本机没有**（`which gofumpt` 空，CI 那一档跑的是 `go install mvdan.cc/gofumpt@latest` ＋ `gofumpt -l`，本腿不擅自跨网络装工具）⇒ 门禁面比 CI 松一档，见 §5 J#3 |
| `bash scripts/d22scan.sh` | 种子正控：`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`；真扫：`d22scan: examined 253 production Go files under internal/ and cmd/`，末行 `d22scan: clean - no D22 ban violations`（存：`.scratch/wisp/probes/245/r1/d22scan.txt`） |
| `go test ./internal/ball/ -count=1`（默认层） | **1 枚红**：`--- FAIL: TestC21TableColourRowsMatchTokensCSS` ＋ `tokens_table_test.go:1468: read design/assets/tokens.css: open ...: The system cannot find the path specified.` ⇒ **归因：起手名册第 14 行已带 ` D design/assets/tokens.css`**（别人留下的工作树删除，先于本腿；本腿零读 `design/**`、`frontend/**`，也没提交任何那条目）。⚠ 除此之外 `--- FAIL` 计数＝**1**，`grep -c "^--- FAIL"` 现跑＝1；本票那族用例（`TestRegisterAll*`／`TestCancel*`／`TestStandby*`／`TestHotkeyStatusString`／`TestDefaultHotkeys*`／`TestApplyHotkeyDefaults`）**逐枚 PASS**（读数存 `.scratch/wisp/probes/245/r1/`） |
| `go test -tags winlive ./internal/ball -run 'TestLiveHotkey\|TestLiveConfirmingCancelAndEscReturned\|TestBallLiveLifecycle'` | 4 枚 PASS／0 枚 SKIP-LOUD／2.691s（修前）＋ 2.707s（突变还原后复跑） |
| `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1` | `ok github.com/CarlosShao/wisp/cmd/wisp 148.421s`＝全绿（含票 228 那族默认层用例与票 127 那枚钉，一枚未放宽）。⛔ 本腿**没有**不带 sherpa PATH 跑它——那条形会造 `exit status 0xc0000135` ＋ `0.0xxs` 且**没有** `--- FAIL`＝用例根本没跑，不是绿也不是红（派单第 5 节坑，本机历史代价） |
| `PATH=... go test -tags winlive ./cmd/wisp -run TestLive228...` | 修后 PASS，`hotkeys live 3/4`；M1 下 FAIL，`hotkeys live 4/4` |
| 未跑的整包 winlive | ⛔ 本腿只跑了指名那几枚真机用例，**没有**跑 `-tags winlive ./...`（票 228 的表具名过：两枚 live 套件共用桌面会互相洗读数，`internal/ball` 的 `requireQuietBallDesktop` 与 `cmd/wisp` 造窗那条不能并发） |

## §4 残余（乙形那半条代价＋本腿新增的具名欠账）

### §4.1 AC#5 的五字段登记（逐字同稿已抄进工单末尾）

| 字段 | 内容 |
|---|---|
| **类型／项** | DEFERRED：`Confirming` 借 Esc 那 2–3 秒仍吞别处的 Esc（票 245 乙形残余） |
| **为什么现在不做（依据）** | 乙形的裁定就是"稳态不绑、借期照旧"（票 245「裁定」三条理由）；把借期也消掉要靠"卡片在场时才临时接管＋用低级钩子而非 `RegisterHotKey`"那一形，属新机制、本票没裁；且 `SPEC-06` 的 L1 窗口本身只有 2–3 秒，`internal/statemachine/timeouts.go:48` 现读 `{state: StateConfirming}: {3 * time.Second, EvConfirmExpired}`＝窗口上限 3 秒 |
| **完成判据** | 卡片挂着期间按 Esc：Wisp 否决生效（D43 #22）**且**同一发里另一个程序的 Esc 也仍能收到（例：焦点在记事本时 Esc 关"未保存"对话框不被吞）＝借 Esc 不再是桌面独占；做不到就要给出"另一条否决通路＋不接管全局 Esc"的实现并过 `Confirming` 那族真机用例 |
| **当前残缺表现** | 卡片等批的那 2–3 秒内，别的应用收不到 Esc（记事本关窗、资源管理器改名、对话框取消、游戏退全屏那一段会失灵）；反向也一样：想否决 Wisp 时若焦点在别的应用并按下 Esc，那发 Esc 归 Wisp |
| **前置／落哪个模块·切片·谁做·触发条件** | 落 `internal/ball`（热键借用那对函数）＋ `internal/agent/approval` 的卡片生命周期；切片＝S6 之后的球／门控联动面（票 228 后续片把卡片接进常驻腿之后才有真用户）；谁做＝先由编排者立一枚后续票（本腿不许自己开票），需要 owner 定形（低级钩子 vs 只在面板内取消）；触发条件＝常驻腿出现真卡片（票 228 AC#2/AC#3 落地）**或**owner 实机撞到"那两三秒按 Esc 没反应"任一条 |

### §4.2 其余具名残余（不藏）

1. **常驻腿没有借的路径**：`cmd/wisp` 今天无状态机、无审批门（`resident_ball_windows.go:21-29` 逐字 `no task pipeline, no microphone, no approval gate`）⇒ 生产里 `TakeEscForCancel` 的调用者**只有** `cmd/balldebug/main.go:635`（调试旁支）。所以本票落地后的事实是"常驻腿稳态 0 枚裸 Esc、也从不借"；AC#2 那两头发是在包内真机窗口＋状态机桩上证的，**不是**在出厂进程的卡片上证的（那格还不存在）。⛔ 不许由此读成"Esc 已经安全了"。
2. **N#5**：用户把 `cancel` 配成带修饰键的组合时，乙形下稳态也不绑它（票面 AC#1 是无条件三枚）。今天这条只影响包内 API 与 `balldebug`（`liveHotkeys()` 那族用例已按新语义改判），出厂路径拿不到配置（N#1）。待人裁，见 §5 J#2。
3. **N#4**：配置热重载（`RebindHotkeys`）会丢掉一次在飞的借用，且没人补借。既有形状，本腿只把注释写清楚，未改行为。
4. **N#2**：`cmd/wisp/main.go:25-26` 那句 `its four global hot keys` 现在成半假话。归票 228 的文案面（A477 已把 `main.go` 那族句子判归 228），本腿不跨票改。
5. `cmd/balldebug/main.go:224` 的 `hotkeys live=%d/4` 分母仍是"槽位数 4"，与 `cmd/wisp` 那枚 `%d/4` 同口径（票面 AC#3② 说文案不改、只核变值）；调试旁支稳态现在打 3/4。

## §5 判不动的地方（判不动就写判不动＋为什么）

| # | 事项 | 为什么判不动 |
|---|---|---|
| **J#1** | `SPEC-12 §5` 那张推迟表里**要不要／由谁**新增 §4.1 那一行 | 表在 `docs/specs/SPEC-12-roadmap-governance.md`＝本腿禁地（"⛔ 不改 `docs/specs/**` 一字"）。同时 `issues/README` 硬约束要求代码里的 `DEFERRED(D-xx)` 标记与该表 **1:1 双向**对得上 ⇒ 我自己既不能加表、也不该加标记（加了就是"有标记没登记"）。⇒ 只能由编排者/人工落这一行，本腿把五字段内容备齐在 §4.1 与工单末尾 |
| **J#2** | 带修饰键的 `cancel` 绑定在稳态该不该绑 | 票面 AC#1 逐字只给了"稳态只能是三枚（summon／mute／panel），`cancel` 不得在场"，没区分裸键与组合键；裁定理由三条讲的是"别在稳态抢 Esc"。**两种读法都能自圆**：无条件三枚（本腿照字面做的）／只在裸键时不绑（会让 AC#3① 那枚钉的"稳态三枚"在读法二下不成立）。⇒ 属"方案没覆盖的情况＝停下来问"（D22 闸门③），本腿按票面字面取严的一支做，并把这个分叉具名交给编排者，⛔ 不自作改判据 |
| **J#3** | `gofumpt` 那一档门禁本机不可用 | 本机 `which gofumpt` 空；装它要 `go install mvdan.cc/gofumpt@latest`（跨网络拉工具），历史上这样执行掉过一次不该执行的东西 ⇒ 本腿不擅自装。只跑了 `gofmt -l`（更松的一档，零命中）。真值要等 CI 或编排者本地跑 |
| **J#4** | "别的程序真能收到 Esc"这一句能不能证到 | 本腿能证的是**两件事**：稳态/归还后 Win32 报告裸 Esc 无人占（探针用另一枚 id 注册成功），以及借期它被占（同一把探针 1409）。**"另一个进程此刻收到了那发 Esc"**需要第二个真程序＋真按键，本腿没有造（不装键盘钩子、不起第二个 GUI 进程）。⇒ 表里 AC#2 的判语只写到"还这一步有具名读数（探针＋报告＋枚数）"，⛔ 不写成"别的应用已经能正常用 Esc 了" |
| **J#5** | `staticcheck` / `test-windows` / `slo-*` 等 CI 档 | 本机没跑（工具在场与否未查全）；本腿只交上面 §3 那批尺的读数，跨档一致性归 CI 与非实现者裁决 |
| **J#6** | 桌面若被一枚外部程序占着裸 Esc，本票的三枚真机用例会走哪一支 | 本机现测：桌面稳态探针回"free"＝这台机器没人占 ⇒ 那条 `SKIP-LOUD`（"attempted and refused"）分支**在本机从未被真走到过**，只有默认层的假注册器（`TestCancelBorrowFailureIsAProblemLine`）走过。判不动"那支在别的机器上会不会把一枚本该红的用例放成绿"——形状上不会（standby＝红、refused＝具名绿），但没实机证据 |

## §6 名册核对（终态＝起手名册）

- 起手：`.scratch/wisp/probes/245/r1/start-status.txt` ＝ **173 行**；其中 `internal/ball`＋`cmd/wisp`＝**0 枚**。
- 终态尺：`git status --porcelain -- internal/ball cmd/wisp` ＝ **0 枚**（本腿在这两枚包里的改动全部进了 commit）。
- 终态全量名册与起手名册的**差分**（本腿视角，逐条归因；⛔ 本腿没提交任何别人留下的脏条目）：
  - **本腿新增并提交的条目**：`docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md`（本表）、`.scratch/wisp/probes/245/r1/start-status.txt`（起手名册逐字）、`internal/ball/hotkey_windows.go`、`internal/ball/ball_windows.go`、`internal/ball/hotkey_status_test.go`、`internal/ball/hotkey_live_test.go`、`internal/ball/live_windows_test.go`、`internal/ball/interaction_live_test.go`、`cmd/wisp/resident_ball_windows.go`、`cmd/wisp/resident_ball_live_228_windows_test.go`。
  - **本腿新造、只留在盘上不提交的临时件**（`issues/README` 规则 8"只建不删"，全部在 `.scratch/wisp/probes/245/r1/`）：`msg-skeleton.txt`、`msg-code.txt`、`build.txt`、`test-ball-default.txt`、`test-cmdwisp-default.txt`、`d22scan.txt`、`winelive-hotkey-1.txt`、`m1-cmdwisp-livelive.txt`。
  - **别人留下的条目**：起手名册里除本腿地界以外的全部条目（`.gitignore`、`design/**`、`.scratch/wisp/probes/**` 其余票目录、`docs/evidence/s1/152-...`、根目录那 10 枚 `part*.txt` 等）**原样不动、未 staged、未提交**。
- 票面 AC 框计数（本腿一枚未勾，尺：`grep -c "^- \[ \]"`／`grep -c "^- \[x\]"` 于工单文件）：**未勾 5 枚／已勾 0 枚**，与起手一致。⚠ 本腿在这份工单里唯一的写入＝末尾追加「残余登记（五字段）」一节（不含任何勾选状态变化），见下一枚 commit 的 pathspec。
- 本腿只 commit、⛔ 未 push；commit 全部带显式 pathspec；未用 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内零删除（含别人的条目）。
