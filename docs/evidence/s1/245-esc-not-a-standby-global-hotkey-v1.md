# 245-v1 — 非实现者对抗验收：票 245（球进常驻进程后裸 `Esc` 被抢成全局热键 ⇒ 乙形：稳态不绑、只在 `Confirming` 借、离开即还）

- 腿：`245-v1`（**裁决腿，非实现者**；只读＋临时单点突变台件，⛔ 零产码改动、⛔ 零 AC 框触碰）。
- 被验收的写腿＝`245-r1`（**同一张表的作者**，所以它那份 `docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md` 按 `AGENTS.md §1.5`／`SPEC-12 §4.3` **不是裁决**；本表只把它当〔待攻主张名册〕，⛔ 不引它的读数当凭据）。
- 锚点自取：骨架这一枚的 `git rev-parse --short HEAD` = 见 §6 起手名册那一节（本腿全程用我自己现取的号）。
- 射程：票面 **AC#1／AC#2／AC#3／AC#4／AC#5／AC#6** 六格＋编排者点名的六处攻面。骨架＝第 5 轮内先交（`issues/README` 规则 6＋本仓"表死三次"教训）。

---

## §0 主张名册（从写腿表里逐条抽出来的、需要我独立复现的主张）

| # | 写腿的主张（原文位置） | 我要用的尺（我自己跑，不背它的数） | 本腿读数 |
|---|---|---|---|
| **C-1** | 稳态注册集枚数＝**3**（summon／mute／panel），`cancel` 不在场（`…-r1.md` §0「一句话结论」、§1 AC#1 ①②③④） | `-tags winlive ./cmd/wisp -run TestLive228…` 出厂进程真机读数＋`-tags winlive ./internal/ball -run 'TestLiveHotkey\|TestLiveConfirmingCancelAndEscReturned\|TestBallLiveLifecycle'`＋默认层 `go test ./internal/ball` | **复现成立**。我自己那发出厂真机逐字 `hotkeys live 3/4`（§2.1），包内 4 枚真机 PASS、确定性 5 枚 PASS；E3 另在投递层给出 `keydown_esc=1` |
| **C-2** | 借用态＝**4**，且"借"与"还"两头各有具名读数（§1 AC#2） | 同 C-1＋`escBorrowProbe` 的**语义射程**我另判（攻面 1） | **机制层复现成立、交付层仓内无判据**。借那头 V-M2b 打得红（记账谎报＝三枚点名红），还那头三处 `requireEscReturned` 真机 PASS；"别的程序收得到"只有本腿自造的第二进程能测（§2.3b） |
| **C-3** | 三枚 AC#3 钉**逐枚改了且是收紧**，另加两枚票面未点名的同类计数钉一并改（§1 AC#3① ② ③） | 逐枚读改后断言字面；**把稳态那支反回去看会不会红**（M-v1 系列） | **成立，且确为收紧**。①`!=4`→`!=3`＋新增 `requireIdleRoster`；②`0/4` 条件一字未动＋新增 `4/4` 判红；③注释改准、用例仍绑 `Ctrl+Alt+V`。V-M1 三层 10 枚红（§2.7） |
| **C-4** | `SKIP-LOUD` 计数 **0 枚**＝连注入那支也没躲过去（§1 AC#2 末、§3） | `-tags winlive ./internal/ball` 那发的 `SKIP-LOUD` 计数我自己数 | **复现成立**：我这一发 `--- FAIL` 0 枚、`SKIP-LOUD` 0 枚、`no tests to run` 0 枚（四把我自己数）。⚠ 但 E6 证明这种"零 SKIP"在别的机器上会以另一种形态失效（§2.3b） |
| **C-5** | 三发突变 M1／M2a／M2b 逐发红过、逐枚 `git diff --quiet` 还原（§2 名册） | 我自己**独立重跑** M1（这格唯一有意义的正控），M2 姿势按 §3 攻不动清单判断是否复现 | **不复用它的读数**：本腿自己跑三发 V-M1／V-M2a／V-M2b，逐发红、逐枚 `git cat-file`＋`git diff --quiet` 复认（§2.7）。它的 M3"不需要另开一发"我读完 `hotkey_status_test.go:236-243` 后认同，未重复 |
| **C-6** | 终态名册＝起手名册（差分恰好一枚、只在自己落点）（§6） | 我自己起手的 `git status --porcelain` 逐字存盘，终态对撞 | **成立**：起手 175 行逐字存盘，终态对撞见 §6（本腿落点之外零动） |
| **C-7** | 门禁全绿：build／vet 两档／gofmt／d22scan（含正控）／`cmd/wisp` 整包（§3） | §1 门禁那一列（全部自跑） | **部分不成立**：门禁本腿全复跑（§2.8），build／vet 四档／gofmt／d22scan 全绿；⚠ 但 `cmd/wisp` 整包默认层本腿读到 **1 红 2 绿**、首发用例名已丢 ⇒ "门禁全绿"这句在本腿读数下不成立，欠一次归因（§4 R-7） |
| **C-8** | AC#5 五字段登记已落两处（§4.1＋工单末尾） | 工单末节现读＋台账 A478 现读 | **成立**：三处齐全（工单末＋台账 A478＋`SPEC-12:80` 已由 owner 批「搬」），且代码里零 `DEFERRED(D-` 标记 ⇒ 双向核对未破 |
| **C-9** | 写腿具名交出的**未修**残余 N#1／N#2／N#3／N#4／N#5（§0.1） | 逐枚读那五行指向的代码；N#4（重绑丢借用）我要读形状到底意味着什么 | **五条全复认，但其中 N#4 的归因是错的（V-F3）**：`ball_windows.go:809-812` 那句"没有调用者"对 `cmd/balldebug` 不成立（`:237-243` 轮询桥＋`:633-638` 卡片同进程）。另本腿新增两条它没交的：V-F1 借被拒后记账永久陈旧、V-F2 新钉跨机器假红（§2.3／§2.3b） |

---

## §1 判语汇总（六格 AC ＋ 编排者六处攻面；全部由本腿现跑读数支撑，凭据列只指本腿自己的尺）

| 格 | 判语 | 凭据（file:line ＋ 我自己跑出来的读数，节号＝出处） |
|---|---|---|
| **AC#1** 稳态名册不含裸 `Esc` | **成立** | 出厂进程真机 `-tags winlive` 逐字 `hotkeys live **3/4**: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not b…`（§2.1 第一行，hwnd=19400042／pid=8392，`--- PASS 3.10s`）；包内稳态 4 枚真机 PASS＋确定性层 5 枚 PASS；交付层第二枚真进程 `keydown_esc=1`（E3）。反回去必红＝V-M1 三层 10 枚红（§2.7）。默认键名未动：`hotkey_windows.go:69` ＋ 钉在 `hotkey_status_test.go:190-192` ⇒ ⛔ 未触 C12／D43 |
| **AC#2** `Confirming` 内借、离开还（两头发） | **附条件成立**（条件两条，都在"射程"，不在"做错"） | **借**那一头：`interaction_live_test.go:205`／`live_windows_test.go:127` 走的是**真 Win32 探针**，V-M2b（谎报形）被它三枚点名打红（§2.7）。**还**那一头：`requireEscReturned`（`hotkey_live_test.go:456-462`）在该枚用例里被调 **三次**（`interaction_live_test.go:222/:250/:273`，键否决／注入 Esc 否决／单击否决各一），三处真机 PASS。**⚠ 条件一（射程）**：判据的"还"量的是**桌面热键表**（`escBorrowProbe`），票面字面那句"别的程序要能收到"在仓里**没有任何用例覆盖**——本腿用第二枚真进程补测了（E1/E2/E3/E4，§2.3b），**补的是验收侧读数，不是仓里的门** ⇒ 要勾这一框的人必须同时认 §4 R-2 那笔欠账。**⚠ 条件二（哪条腿）**：出厂常驻腿今天**既没有卡片也没有借的路径**（`TakeEscForCancel` 的全部生产调用者＝`cmd/balldebug/main.go:635`，本腿 `grep -rn` 现跑，见 §2.3b 上方），所以"借—还"整对在**产品进程**里既未发生也无从发生 ⇒ 这一格证到的是包与调试旁支，⛔ 不许读成"Esc 已经安全了"（乙形那 2–3 秒的借窗仍在，见 §3 与 SPEC-12 §5:80 那行） |
| **AC#3** 三枚钉同批改期望、且是**收紧** | **成立**（附一枚具名缺陷 V-F2：新钉跨机器会假红） | 逐枚"改前 vs 改后"字面对撞在 §2.4：①`live_windows_test.go:74` `!=4`→`!=3` 并新增 `:78 requireIdleRoster`（五件一起判）；②228 那枚 `0/4` 条件**一字未动**、注释改准、**新增 `:163` 让 `4/4` 判红**；③`hotkey_live_test.go:34-38` 注释改成"稳态根本不绑"、用例仍绑 `Ctrl+Alt+V`（`:42`）＝票面那支禁令守住了。**"收紧后仍能红"这一发我实跑了**：V-M1 下红在 `live_windows_test.go:75`，红句逐字「the live ball did not register its three idle hotkeys (summon/mute/panel, cancel is borrowed only during Confirming)」，出厂层同时红在 `resident_ball_live_228_windows_test.go:164` ⛔ 不是装饰。**V-F2**：E6 那一发（外来进程占住裸 Esc）把四枚真机用例**全打硬红**、`SKIP-LOUD` 一枚都不出 ⇒ 在"别的程序已经占了 Esc"这种机器上，正确的实现也过不了（凭据与修法 §4 R-4） |
| **AC#4** 用户能改 `[hotkey] cancel` | **判不了＝阻塞于票 228，且本格按票面形状如实写出、没假绿（合格）** | `cmd/wisp/resident_ball_windows.go:86` 逐字 `Hotkeys:  ball.DefaultHotkeys(),`；`grep -rn "config.Manager|hotkeysFromConfig" cmd/wisp/*.go`（剥 `_test`）现跑＝类型只在 `cmd/wisp/run.go:228` ⇒ 常驻腿整棵不读 `[hotkey]`。⛔ 本腿复认**没有第二条配置通路**被造出来。票面逐字允许这一支（"判不了就写成阻塞于票 228，不许假绿"）⇒ 判**如实**，⛔ 不许翻勾 |
| **AC#5** 残余五字段登记 | **成立**（且比写腿交回时又多落一处） | 三处齐全：工单末尾 `:72-78` 五字段表／台账 `A478`（`docs/reports/pending-and-issues.md:9943`）／**`docs/specs/SPEC-12-roadmap-governance.md:80` 现读已有那行 `| DEFERRED | 确认窗口借用裸 Esc（票 245 乙形残余） |`**（owner 17:5x 一句「搬呗」，由编排者落，commit `dc694430`）。⛔ 代码里没有 `DEFERRED(D-xx)` 标记（`grep -rn "DEFERRED(D-" internal cmd --include=*.go` ＝**0 枚**）⇒ "标记↔表 1:1 双向"那道核对没被打破，这一步判断我复认正确 |
| **AC#6** ①"四枚"文案 ②判据锚的形状 | **①不成立（这句今天仍是半假话，本格不该勾）；②编排者那句自裁我独立判成立，但其完成判据未做 ⇒ 本格也不该勾** | ①：`cmd/wisp/main.go:25` 现读逐字「its four global hot / keys in this same process」，而我**真跑了这一版二进制的 `-h`**（`wisp-clean.exe -h` 输出第 5-6 行同样逐字含那句），同一枚二进制稳态实测 `hotkeys live 3/4` ⇒ **半假话成立、未修**。修的时候有一枚形状钉要保住：`cmd/wisp/leg_dispatch_gate_133_test.go:156` `usageBare133`＝`(?m)^  wisp {2,}\S`。②：三条凭据在 §2.5——`cancelIdleLine`（`hotkey_windows.go:428-449`）不看 `acc.Mods`、**E5 实测**（`[hotkey] cancel = "Ctrl+Alt+V"` → `live=3/4` 且日志逐字 `cancel hotkey left unbound while idle … binding=Ctrl+Alt+V`）、四枚钉全锚在枚数（`hotkey_live_test.go:435`／`live_windows_test.go:74`／`hotkey_status_test.go:181`/`:245`）⇒ 按"有无修饰键"改钉是**必须**的，且要**先改钉后改行为**（§4 R-6） |
| **攻面 1** 那枚新探针到底证到哪一层 | **分两层：机制层成立／用户手感层仓内判据未证**（⇒ AC#2 只能附条件） | 全文 §2.2。**探针自身能分辨**这件事有双向凭据：被测物侧 V-M2b 三枚点名红（§2.7），环境侧 E6 一枚真外来占位者让探针回"不 free"并打出逐字 `bare Esc is claimed on this desktop while the ball is idle`（§2.3b）。⇒ 它**不是**只回"没人占着"的哑尺；但它量的确实只是热键表，**投递层**要第二枚进程才看得见，那一枚本腿造了（E1 `keydown_esc=1`／E2 `-steal` 下 `keydown_esc=0`＋`wm_hotkey=1`），**仓里没有**。更强的形状我给了具体判据（§4 R-2），不空口要求 |
| **攻面 2** 借用与归还的生命周期 | **四形无双注册／无漏还／无反复借；但揪出两条：N#4 的归因是错的（V-F3）、借被拒后记账永久陈旧（V-F1）** | 全文 §2.3。逐形凭据：`ball_windows.go:875` bool 守卫＋`hotkey_windows.go:532` 先撤再绑＋`hotkey_status_test.go:259` 钉双借；`releaseEscWith`（`:553`）只有 `unreg`＋`hotkey_status_test.go:239-243` 钉"还＝不再绑"；`Close()`（`:944`）走 `unregisterAll`（`hotkey_windows.go:507-511` 四枚 id 逐枚撤，注释 `:498-502` 逐字点到 "a takeover id"）⇒ **退出漏还在结构上不可能**（读码＋次序，非实测）。**V-F3**：`ball_windows.go:809-812` 那句「no caller does that today because the config poll and the card live on different legs」**是假话**——`cmd/balldebug/main.go:237-243` 装轮询桥、`:633-638` 驱动卡片，**同一枚进程两样都有** ⇒ 借期中改 `[hotkey]` 就把借用丢掉且无人补借，今天可达（调试旁支；出厂腿仍不可达）。**V-F1**：`TakeEscForCancel` 失败支不设 `escTakenOver`，而 `ReleaseEscAfterSession` 首行就按它早退 ⇒ 报告永久停在 `HotkeyError`、`AllLive()` 从此恒 false、`Problems()` 一直印那行"没注册"，而桌面上其实什么都没绑 |
| **攻面 3** AC#3 那三枚钉是否真收紧 | **成立**（判据见 AC#3 行；反回去会红＝我实跑） | `live_windows_test.go:74` 现逐字 `!got.AllLive() \|\| len(got.Live()) != 3` ＋ `:78`；V-M1 红句逐字（§2.7 第一行第三列）；`0/4` 条件未被放宽、`4/4` 是**新增**一形判红；⛔ 零删断言、⛔ 零词面型新仪器（我另跑一把 `grep -rn "global hot|four global" cmd/wisp/*_test.go` 只命中那枚**错误消息文本**） |
| **攻面 4** 它揪出的那两枚未点名计数钉 | **两枚都改对了；票 64 的判据没有被顺手放宽**（第二枚是**语义反转＋加强**，不是削弱） | 第一枚：`hotkey_live_test.go` boot 支现走 `:149 requireIdleRoster`、rebind 支 `:174` `len(live) == 3 && live[hkSummon].VK == 'R'`（改前逐字 `== 4`），`VK=='R'`／`Rebinds()==1`／注入新旧两键那三段**一字未动**。第二枚：改前逐字 `if rep := b.HotkeyReport(); !rep.IsLive(hkCancel) { t.Errorf("after the B1 return the configured cancel binding is not live…") }` → 现 `requireEscReturned`（bool＋枚数＋三枚点名＋standby 行＋桌面探针五件），旧断言钉的"release 会把配置那枚重新绑回去"**正是本票裁死不许的行为**（`hotkey_windows.go:547-553`），反过来才是对的。⚠ **附带一笔不属于本票的账**：票 64 已勾框 `.scratch/wisp/issues/64-…md:44` 逐字引着「断 `len(live)==4`」，今天码里是 `==3` ⇒ **那是别人票面上的过期指认**，本腿不改框不改票面，只上报（§4 R-5） |
| **攻面 5** AC#6 两半＋编排者那句自裁 | **我独立裁：①"这句是半假话、本格不该勾"＝成立；②"判据要按有无修饰键、我写窄了"＝成立**（两半都不是因为他先写了） | 凭据见 AC#6 行＋§2.5 三枚。⛔ 我不接受"今天零行为差别"就当无事：E5 已经在**读配置那条腿**上量到"用户配的组合键也不绑"，`balldebug` 的日志逐字写着 `cancel hotkey left unbound while idle … binding=Ctrl+Alt+V`——分叉是真的，只是出厂路径今天还到不了 |
| **攻面 6** `winlive` 这些真机用例 CI 到底跑不跑 | **不跑。这道门今天不存在** | `grep -rn "winlive" .github/workflows/ scripts/` ＝**零命中**；四枚 live 文件头逐字 `//go:build windows && winlive`；CI windows 腿只跑 `scripts/portable-tests.sh --scope=windows`（`ci.yml:512`）与 `scripts/wisp-cli-tests.sh`（`ci.yml:476`），两处都不带那枚 tag ⇒ 这些用例**根本不进 CI 的测试二进制**。而且 `gh run list --limit 5` 现跑最近五枚（`36674742682` 05:43Z／`36664144050` 03:23Z／`36643624549` 23:09Z…）**全部早于本票四枚 commit（09:06–09:45Z）** ⇒ 说不出任何 run id＋step 来背书。**残余登记＝§4 R-3；⛔ 不据此判 AC 不成立**（票面没要求 CI 覆盖，见 §2.6 末段） |

---

## §2 逐格凭据与读数（每一发现跑现记；突变逐发单列）

### §2.0 起手闸门与台件姿势

- **起手 HEAD（我自己现取）**：`a43fe78b`（＝编排者 A478 那一枚，17:45:53 +08）。⚠ 派单给我的任何 sha 一律未采用。跑动期间编排者又落了一枚 `dc694430`（SPEC-12 §5 加 DEFERRED 行＝owner 17:5x 一句「搬呗」，台账 A479），**与本腿无冲突**（我没碰 `docs/specs/**`）。本腿骨架 commit＝`2f2d398e`。
- **起手全量名册逐字存盘**：`.scratch/wisp/probes/245/v1/start-status.txt`＝**175 行**（取数 17:47 +08）。其中 `internal/ball`／`cmd/wisp` 计数＝**0 枚**（尺 `grep -c "internal/ball\|cmd/wisp"` ＝ 0 ＝没找到＝那两枚包干净）；`docs/evidence/s1/` 只有别人那枚 ` M 152-subject-death-never-measured-r1-accept-r1.md` 与我自己的骨架（名册第 167 行 `?? docs/evidence/s1/245-...-v1.md`）。
- **起手名册里那枚常红的归因**：名册第 14 行 ` D design/assets/tokens.css`（＋` D design/assets/base.css`／`icons.js`／`theme.js`／` D design/index.html` 一族）＝**别人留在工作树的删除，先于本腿、也先于 245**。核对尺：`git diff --name-only 25ce3ce9 60f888c8 | grep -c "tokens_table"` ＝ **0 枚**，且 `internal/ball/tokens_table_test.go` 末次 commit 是 `ed43bc3d`（票 74）⇒ `internal/ball` 默认层那枚 `TestC21TableColourRowsMatchTokensCSS` 红**不是 245 造成的，也不是本腿造成的**（读数见 §2.7 门禁那一列）。⛔ 本腿**零读**`design/**`／`frontend/**`，只用"文件不存在"这一条错误消息做归因。
- **桌面三连（每一发真机用例前都复量）**：`tasklist //FI "IMAGENAME eq balldebug.exe"`／`eq wisp.exe`／`eq go.exe` 起手＝**0／0／0**（17:47 现跑）。
- **PATH 姿势**：跑 `./cmd/wisp` 一律带 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`；判红绿**只认 `--- FAIL` 行**（不带 PATH 会给 `0xc0000135` 且零枚 `--- FAIL`＝一枚用例都没跑，那不是红也不是绿）。
- **不跑整包 winlive**（两枚 live 套件共用桌面会互洗读数，票 228 的表具名过）；**不跑 `ball-cycle.ps1`**；本腿全部台件在 `.scratch/wisp/probes/245/v1/`（**只建不删**）。
- **还原姿势**：`git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet -- <path>` 复认；⛔ 未用 `checkout`／`stash`／`reset`／`--amend`／`clean`；⛔ 未用 `-overlay`（本包有走 AST／读盘的用例，overlay 照不到）。

### §2.1 AC#1 稳态注册集（我自己跑出来的读数）

**稳态枚数＝3，`cancel` 不在场。** 四层读数全部本腿现跑：

| 层 | 我这发的命令（逐字） | 我自己读到的数 |
|---|---|---|
| **出厂进程真机（最硬的一发）** | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -tags winlive ./cmd/wisp -run 'TestLive228ResidentLegOwnsABallWindowOnTheDesktop' -count=1 -v` | `resident_ball_live_228_windows_test.go:147: AC#1 LIVE: hwnd=19400042 pid=8392 owns the ball window; hotkeys live **3/4**: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not b`（末段被该尺自己 `:145` 那枚 80 字符窗口截断；未截断的全文由 §2.3b-E3 那一发直接从进程控制台读到：`cancel=Esc not bound while idle (cancel is borrowed only during Confirming)`）。`--- PASS (3.10s)`／`ok ... 3.155s` |
| **包内稳态（真机窗口）** | `go test -tags winlive ./internal/ball -run 'TestLiveHotkey\|TestLiveConfirmingCancelAndEscReturned\|TestBallLiveLifecycle' -count=1 -v` | 4 枚 `--- PASS`（RebindEndToEnd 0.60s／OccupiedVsNotAttempted 0.03s／ConfirmingCancelAndEscReturned 0.05s／BallLiveLifecycle 1.90s）、`ok 2.638s`、**`--- FAIL` 0 枚、`SKIP-LOUD` 0 枚、`no tests to run` 0 枚**（这四把我自己数的，不是背写腿的）。稳态那两把尺都在跑：`live_windows_test.go:74` 的 `len(got.Live()) != 3` 与 `hotkey_live_test.go:429 requireIdleRoster` |
| **确定性层（假注册器）** | `go test ./internal/ball -run 'TestRegisterAllLiveSet\|TestDefaultHotkeysIdlePassHoldsNoEsc\|TestCancelBorrowRoundTrip\|TestStandbyIsNotDisabledAndNotAProblem\|TestRegisterAllSplitsFailureFamilies' -count=1 -v` | **5 枚 `--- PASS`**、`ok github.com/CarlosShao/wisp/internal/ball 0.027s`。逐枚逐字读过断言：`:120` 要求 `attempted` 恰好 `[1 2 4]`、`:179` 要求 `reg.bindsEsc()` 为假、`:181` 要求 `len(Live()) != 3` 即红、`:205` 要求稳态对 cancel 的注册尝试次数＝**0** |
| **交付层（第二个真进程收得到 Esc）** | §2.3b 的 E1／E3／E4（本腿自造台件） | 出厂 `wisp.exe`（同一枚源码、稳态 3/4）**挂着的时候**，另一枚真 GUI 进程抢到前台、被注入一发裸 `Esc`：**`keydown_esc=1`**＝收到了。对照发 E4＝把稳态那一支反回去（M1 产码突变）造出的 `wisp-mutant.exe`（自报 `hotkeys live 4/4: … cancel=Esc live`）挂着时同一发：**`keydown_esc=0`**＝被吞掉 |

生产默认键名一字未动：`internal/ball/hotkey_windows.go:69` 现读 `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}`，且 `hotkey_status_test.go:190-192` 那枚反向钉（`DefaultHotkeys().Cancel != "Esc"` 即红）把它钉在"D43 #22 那句话不许换键名"上 ⇒ **没有走甲形、没有触碰 C12 冻结面**（这一条是纪律面，我自己核过码与用例两处）。

### §2.2 那枚新探针到底证到哪一层（攻面 1，分两层报）

**它是什么**：`internal/ball/hotkey_live_test.go:57-81` `escBorrowProbe` —— 在球自己的 hwnd 上、用**另一枚 id**（`spareHKID = 900`，`:30`）向真 Win32 注册 `modNoRepeat + VK_ESCAPE`，注册成功＝"这台桌面的热键表里没有裸 Esc"、返回 1409＝"有人占着"，然后立刻 `UnregisterHotKey`（`:71`）不留东西。跑法合规：它 `b.sta.PostTask` 回窗口属主线程执行（`:60`，注释 `:54-56` 具名"Win32 refuses it anywhere else"）。

**机制层——证到了，而且比"自述"强一档**：
- `RegisterHotKey` 的冲突判定是**按桌面**不按进程：同一 `(mods, VK)` 已被任何窗口占用就返 `ERROR_HOTKEY_ALREADY_REGISTERED`。本仓自己拿这条当尺用了两次——`hotkey_live_test.go:232-235`（同窗另一枚 id 占住 `Ctrl+Alt+U`，球再绑同一枚必须吃 1409）与 `:30` 那句注释 "for squatting a combination in this very process to prove what a real 1409 looks like"。⇒ 探针回"free"读的是**整个桌面的热键表为空位**这一事实，不是 Wisp 自家的记账；这一档是真的。
- 它还配了**反向射程**：`requireEscBorrowed`（`:110-115`）在"报告说借到了、探针却说 free"时**判红**（逐字 `the ball says it borrowed Esc, but the desktop still has it free - the borrow registered on some other key, not on VK_ESCAPE`），`requireIdleRoster`（`:445-451`）在"稳态却有人占着"时判红。两向都有句，所以这不是一把只会被动回"free"的瞎尺。**⚠ 但"探针能否分辨"这一问，仓里从来没有一发真机正控证过**——写腿的 M2b 是把产码改成"谎报已借到"来打红它，那是**被测物**的正控，不是**尺本身**的正控（尺自己的正控需要桌面上真出现一枚外来占位者）。本腿为此造了 E6（`§2.3b`）：用第二枚真进程占住裸 Esc，看这把尺会不会回"不 free"。读数在 §2.3b-E6。

**用户手感层——没证到，而且仓里的判据里根本没有这一层**：
- 票面 AC#2 的字面是「卡片撤下后再按 Esc，**别的程序要能收到**」。探针量的不是"投递"，是"注册表"。从"表里没有那枚"到"焦点窗口收到 `WM_KEYDOWN`"之间还有若干**本票未测**的环：低级键盘钩子（`WH_KEYBOARD_LL`）不占热键表也能吞键；前台/焦点策略；裸 Esc 被别的窗口以非热键方式吃掉。这些**一把注册型探针全照不到**。
- ⇒ **判语（分两层报，不整体放行也不整体判死）**：
  - **机制层：成立**。"离开确实还了"被证到"桌面热键表里那枚裸 Esc 空出来了"这一层，且这层是 Win32 事实而非记账（并有 §2.3b E1/E2 两发本腿独立复现）。
  - **手感层：仓内判据未证**。票面那一半（别的程序收得到）在 `internal/ball` / `cmd/wisp` 里没有任何对应用例；本腿用第二枚真进程补测了（§2.3b E1/E2/E3/E4），**但那是验收侧台件，不进仓就守不住回归**——所以 AC#2 的"还"我判**附条件成立**，条件＝"把 `§2.3b` 那一形钉成 winlive 用例"（见 §4 R-2 的修法与判据），⛔ 不许被读成"别的应用已经能正常用 Esc 了"是由仓里的门证的。
- **要不要更强的形状**：要。我给的形状不是"再注册一次"，而是**第二个真 GUI 进程 + 前台窗口 + `keybd_event` 注一发 Esc + 数自己收到的 `WM_KEYDOWN`**，配一发 `-steal` 正控（同一发 `RegisterHotKey` 由观察者自己做，必须收不到键）。这一形我已实现并跑通（§2.3b），它对"`takeEscWith` 那一发到底吞不吞键"给的是**投递层**答案，比探针高一层。
- ⛔ **我不因为"它已经比自述强"就放行手感层**，也**不因为"证不到终极那一层"就把 AC#2 整格判不成立**：AC#2 要防的是"只测借、不测还"（票面逐字「判据形状＝'还'这一步有具名读数，不是只测'借'」），这一条**写腿做到了**——还的那一头有 `requireEscReturned`（`:456-462`）三处真机调用（`interaction_live_test.go:222/:250/:273`）＋枚数＋状态＋探针四件读数。

### §2.3 生命周期（攻面 2：连两次 `Confirming`、借期中取消、退出路径、N#4）

逐条现读，⛔ 不引写腿表的任何数：

| 问的形状 | 盘上事实（file:line 逐字读） | 判 |
|---|---|---|
| **连着两次 `Confirming` ⇒ 双重注册？** | `ball_windows.go:873-889` `TakeEscForCancel` 首行 `if b.escTakenOver { return }` ⇒ 第二次同步**根本不碰 Win32**。原语层另有一道：`takeEscWith`（`hotkey_windows.go:531-539`）先 `unreg(hkCancel)` 再 `reg(...)`，所以就算有人绕过 bool，同一 hwnd 同一 id 也不会叠第二份。确定性层把这形钉住了：`hotkey_status_test.go:259` 逐字「Borrow twice without returning … the second attempt must be refused by the registry that already holds it」，断言 `err == nil` 即 `t.Error`。 | **无双注册**（三层：bool 守卫＋先撤再绑＋用例钉） |
| **借期中被用户取消（单击否决那一支）** | `interaction_live_test.go:252-273`：卡片再挂 → `requireEscBorrowed` → 注入单击 → `EvVeto` → `syncBall()` → `ReleaseEscAfterSession`；末尾 `:268` 判 `EscTakenOver()` 必须 false，`:273` `requireEscReturned` 判稳态名册＋桌面探针 free。**否决的两支（键与鼠标）都各自欠一次"还"的读数**，这一形有钉。 | **有覆盖**（真机档；见 §2.4 正控） |
| **还了又借（release 回头再绑配置那枚）** | `releaseEscWith`（`hotkey_windows.go:553`）函数体逐字 `func releaseEscWith(unreg unregisterFn) { unreg(hkCancel) }`——**只有撤销、没有 reg**；文档 `:547-552` 逐字把理由写死（回头再绑＝本票要修的缺陷回魂）。钉：`hotkey_status_test.go:239-243` 「cancel registration attempts = %d after the return, want still 1」。 | **不存在该形状**，且有确定性钉 |
| **退出路径漏还**（D38(e) 第 2 步"热键停听"） | `Close()`（`ball_windows.go:935-969`）在 STA 任务里跑 `unregisterAll(b.hwnd)`（`:944`），`unregisterAllWith`（`hotkey_windows.go:507-511`）对 `hkNames` **四枚 id 逐枚 unreg**、含 `hkCancel`，且注释 `:498-502` 逐字写明"包括我们已经不再跟踪的那一枚（a takeover id）"——**正是为借用的那枚 id 写的**。常驻腿的调用链：`resident_ball_windows.go:151-157`：`func (rb *residentBall) stop()` → `rb.b.Close()`；D38(e) 第 2 步文本在 `internal/proc/shutdown.go:18`。⇒ 即便 `escTakenOver` 仍为 true 时进程离开，那枚注册也随 `unregisterAll`／窗口销毁消失；**结构上留不下"退出后桌面上还占着裸 Esc"**。 | **不可能漏还**（读码＋次序，非实测；实测面见 §2.3b E3） |
| **N#4：重绑丢掉在飞的借用** | `RebindHotkeys`（`ball_windows.go:813-825`）：`unregisterAll` → `registerAll` → `b.escTakenOver = false`（`:822`），**没有任何一支在重绑后把借用来一次**。形状＝卡片还挂着、`cancel` 那枚键已不在场，直到下一次状态同步才可能回来。⚠ **但写腿给这条的理由是错的**：`ball_windows.go:809-812` 逐字「no caller does that today because the config poll and the card live on different legs」——`cmd/balldebug` **同一枚进程里两样都有**：`:237-243` 装 `NewHotkeyReloader` ＋ 1 秒轮询 `mgr.CheckAndReload`，`:633-638 syncMachineToBall` 驱动卡片。⇒ 今天就能走到的形：球进 `Confirming`（借期中）→ 有人在这一秒内改 `[hotkey]` → 重绑丢掉借用 → **卡片挂着却没有 Esc**，且 `relaseEsc` 之后无事发生。出厂常驻腿今天既无卡片也无配置轮询（`resident_ball_windows.go:21-29`），所以**产品路径不可达**；可达的是调试旁支。 | **残余成立、归因错**：账要记（注释那句"没有调用者"是假话），修法（重绑后照 `EscTakenOver` 补借一次，或重绑前记下借期中）属 `internal/ball` 地界＝下一程写腿，本腿未修。见 §4 R-1 |
| **借被拒之后的记账陈旧（本腿新发现 V-F1）** | `TakeEscForCancel` 失败支（`ball_windows.go:878-884`）把 cancel 行写成 `cancelFailedLine(err)`（Status=`HotkeyError`）但**不设 `escTakenOver`**；于是 `ReleaseEscAfterSession`（`:896-899`）首行 `if !b.escTakenOver { return }` **直接返回，那一行永远留在 `HotkeyError`**。后果两笔：① `AllLive()`（`hotkey_windows.go:344-351`）此后一直回 false、`Problems()` 一直印那行「was not registered」，而**桌面上其实什么都没绑**（拒绝＝没借到＝没注册）；② 真机那三枚钉会因此**在"别人占着裸 Esc 的机器"上判红**：`requireIdleRoster`（`hotkey_live_test.go:442-444`）要求 `Status == HotkeyStandby`，`HotkeyError` 直接把 `t.Errorf` 打红，`:449-451` 那枚桌面探针同样会红。⇒ 这**同时推翻写腿表 §5 J#6 那句**「形状上不会（standby＝红、refused＝具名绿）」：`refused` 在 `requireEscBorrowed`（`:99-104`）里确实是具名绿，**但走完之后 `requireEscReturned`／`requireIdleRoster` 那一支是硬红**。 | **本腿判：钉子跨机器不成立**（真机红因＝记账陈旧＋桌面探针射程过宽，不是行为缺陷）。实测复现见 §2.3b E6 |

### §2.3b 真机实验台件（本腿自己造的第二枚进程）

台件＝`.scratch/wisp/probes/245/v1/esclistener/{go.mod,main.go}`（独立模块，根模块 `./...` 编不到它，与票 228 那枚 `twobreak` 同姿势），二进制造到 `$TEMP/wisp245v1/esclistener.exe`（⛔ 不落仓）。它做的事：**真的第二个 GUI 进程**建窗、抢前台、`keybd_event` 注一发裸 `Esc`，然后数**自己这个窗口被投递到的** `WM_KEYDOWN(VK_ESCAPE)`；`-steal` 那一档先替自己 `RegisterHotKey(modNoRepeat, VK_ESCAPE)`（＝与 `takeEscWith` 逐字同一发 Win32 调用，只差注册者是谁）。

| 发 | 命令（逐字） | 读数（逐字，全文存 `.scratch/wisp/probes/245/v1/e*.txt`） | 这一发钉住什么 |
|---|---|---|---|
| **E1** 基线 | `esclistener.exe -watch`（桌面上除 shell 外没有别的窗主） | `hwnd=0x6EE052C foreground=0x6EE052C foreground_is_mine=true steal=false **keydown_esc=1** syskeydown_esc=0 wm_hotkey=0 other_keys=0` | **这把尺在"能收到"那一头不瞎**：前台窗口确实拿到了注入的 Esc |
| **E2** 正控（尺自身的） | `esclistener.exe -steal` | `foreground_is_mine=true steal=true **keydown_esc=0** syskeydown_esc=0 **wm_hotkey=1** other_keys=0` | **一枚全局裸 `Esc` 注册确实把投递吃掉了**：同一发注入变成观察者自己的 `WM_HOTKEY`、前台窗口收不到键。⇒ E1 那个"1"与 E3 那个"1"不是背景噪声；而 `takeEscWith` 干的就是 E2 这一发调用 |
| **E3** 出厂进程稳态 | 手起 `wisp-clean.exe`（`WISP_ENV=test`＋独立数据目录，同一枚 HEAD 源码）＋ `esclistener.exe -watch` | 进程自报逐字：`hotkeys live 3/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, cancel=Esc not bound while idle (cancel is borrowed only during Confirming), panel=Ctrl+Alt+P live`；观察者逐字：`foreground_is_mine=true steal=false **keydown_esc=1** wm_hotkey=0` | **票面 AC#1/AC#2 那一半"别的程序要能收到 Esc"，在出厂进程稳态这一形上今天为真**——这一条是**本腿交付层实测**，仓里的判据没有它（见 §2.2 与 §4 R-2） |
| **E4** 反回去的出厂进程 | M1 产码突变造 `wisp-mutant.exe`（同一发 `-watch`） | 进程自报逐字：`hotkeys live **4/4**: … cancel=Esc live, panel=Ctrl+Alt+P live`；观察者逐字：`foreground_is_mine=true **keydown_esc=0** wm_hotkey=0` | E3 那一发"1"**有对照组**：把稳态那支反回去（＝245 之前的出厂形状）就变成 0。⇒ E3 不是"这台机器本来就不会吞键"，而是"这一版不吞" |
| **E6** 外来占位者 | `esclistener.exe -steal -hold 30000` 挂着 ＋ `go test -tags winlive ./internal/ball -run 'TestLiveHotkey\|TestLiveConfirmingCancelAndEscReturned\|TestBallLiveLifecycle'` | **四枚全硬红**，红句逐字（四处同一句）：`ticket 245 RED: bare Esc is claimed on this desktop while the ball is idle`（`hotkey_live_test.go:149`／`:309`／`interaction_live_test.go:177`／`live_windows_test.go:78`）；`SKIP-LOUD` **0 枚** | 两笔：**① 好的一头**——`escBorrowProbe` 对**真外来进程**确实回"不 free"，这把尺能分辨，不是永远放行的瞎尺（§2.2 机制层因此升一档）。**② 坏的一头（V-F2）**——球完全按 245 落地的形状（稳态谁也没绑）却因**别人的**占位被判红，且 `requireIdleRoster` 跑在任何借用之前 ⇒ 写腿 §5 J#6 那句"refused＝具名绿"在这种机器上**根本走不到**（`requireEscBorrowed` 那支还没执行用例已经 fatal）。⇒ 见 §4 R-4 |
| 收尾 | `taskkill //IM esclistener.exe //F`（台件自带 40 秒看门狗也会自己退） | 复量 `esclistener.exe` 计数＝**0**；随后 winlive trio 复跑 4 枚 PASS／`ok 2.639s` | ⚠ **本腿自己踩到一枚坑并如实记**：E6 之后我按还原流程复跑 winlive trio，第一发是 `FAIL`（0.427s）——红因不是产码，是**我的 holder 进程还挂着占裸 Esc**。这正是票面第 64 行那类"自己的调试进程造 8 枚假红"的形状，换了个占位键重演一遍。**判红绿之前必须数进程**，这一条我写进 §6 终态自证 |

（E5 不在交付层，是 AC#6② 的实测，落在 §2.5。）

### §2.4 AC#3 三枚钉＋两枚未点名钉（攻面 3／4）

**"改前逐字 → 改后逐字"全部由我自己 `git show 25ce3ce9:` 与 `git diff 25ce3ce9 60f888c8` 现取**，不是抄表。

| 枚 | 改前（逐字，来自 `git show 25ce3ce9:<path>`） | 改后（逐字，盘上 HEAD） | 是收紧还是放宽？ |
|---|---|---|---|
| **①`live_windows_test.go:68-70`（票面点名的那枚）** | `if got := b.HotkeyReport(); !got.AllLive() \|\| len(got.Live()) != 4 {` ＋红句「the live ball did not register its four hotkeys」 | `:74` `if got := b.HotkeyReport(); !got.AllLive() \|\| len(got.Live()) != 3 {`，红句改成「the live ball did not register its three idle hotkeys (summon/mute/panel, cancel is borrowed only during Confirming)」，**并紧跟 `:78 requireIdleRoster(t, b)`**（`hotkey_live_test.go:429-452`：cancel 在场即红＋枚数必须 3＋三枚逐一点名 live＋cancel 行必须 standby＋桌面探针必须 free，**五件一起判**） | **收紧，而且是双向**：多一枚红（`!=3`＋`held` 那支）、少一枚也红（`!=3`＋`IsLive` 三枚）、稳态占着裸 Esc 红（探针）。**正控我实跑过**：§2.7/M1 把稳态那支反回去 ⇒ 这一枚**逐字红在 `live_windows_test.go:75`**（红句全文在那发读数里）⇒ 不是装饰 |
| **②`cmd/wisp/resident_ball_live_228_windows_test.go`（票面点名的那枚）** | 票面旧行号 `:151` ＝**现读 `:160`** `if strings.Contains(verdict, "hotkeys live 0/4")` → 旧红句 `t.Errorf("AC#1 RED: the ball came up with none of its four global hot keys registered.")` ＋ 旧注释「A ball whose **four** hot keys all failed to register…」 | `:160` 那枚**条件一字未动**，只有 `:161` 红句里的 "its four global hot keys" 改成 "its global hot keys"（`git diff` 现读：`-`/`+` 两行的差异只在 "four" 这个词）；`:143-147` 那把只截 80 字符当日志的尺**未动**；**新增 `:163-166`**：`if strings.Contains(verdict, "hotkeys live 4/4") { t.Errorf("AC#1 RED (ticket 245): …") }` ＋注释改成带条件事实句（`0/4` 与 `4/4` 两形各归一枚 bug，逐字点名票 64 A1b 与票 245） | **收紧**（新增一形判红、把半假话注释改准），⛔ 零删断言。**正控我实跑过**：M1 下这一枚**逐字红在 `:164`**，同一发里子进程自报 `hotkeys live 4/4: … cancel=Esc live`（§2.7/M1 读数）；`0/4` 那一枚我没复现（它要的是"四枚全注册失败"，与本票不同形，且票 228-v1 已在 M5 证过它是活的尺——⛔ 那枚读数出处是别人的表，我只用它支撑"未复现不等于没钉"） |
| **③`hotkey_live_test.go:32-43`（票面点名的注释那枚）** | 注释逐字「the bare "Esc" cancel default: registering Esc as a GLOBAL hotkey swallows Esc from every other app on the desktop for the whole test run, so the live tests bind a modifier combination instead」 | 现逐字（`:34-38`）「an idle ball no longer registers cancel at all (ticket 245 - the production default is a bare Esc, and RegisterHotKey would take Esc from every other app on the desktop for the whole test run), and while the borrow is live the suite must still be able to tell "our cancel slot" apart from "the key the user configured"」⇒ **"生产默认稳态也绑裸 Esc"那半句已改成"稳态根本不绑"**，并给出第二理由（借期要能分辨） | **合规**：票面 AC#3③ 要的正是"改成稳态不绑、借期才绑的事实句"。⛔ **用例仍然绑 `Ctrl+Alt+V`**（`:42` 现读 `Cancel: "Ctrl+Alt+V"`），一枚都没改回裸 Esc＝票面那句 ⛔ 禁令守住了 |
| **④（票面未点名，写腿自己揪的）`hotkey_live_test.go` boot／rebind 两枚期望** | 改前逐字：`if got := boot…`（`:96-99` 旧形）与 `seen := waitFor(3*time.Second, func() bool { r.Check(); live := b.RegisteredHotkeys(); return len(live) == 4 && live[hkSummon].VK == 'R' })`（`git show 25ce3ce9:…` 现读 `len(live) == 4`） | 现读：`:149` 走 `requireIdleRoster(t, b)`（注释在 `:143-144`；枚数＋standby＋探针三件一起判）、`:174` `return len(live) == 3 && live[hkSummon].VK == 'R'` | **改对了**，且 `:149` 那一枚现在同时能抓"稳态多绑一枚"与"该借的时候没借上"之外还多抓探针。**⚠ 但这一改把票 64 已勾 AC 的引文变成过期指认**：`.scratch/wisp/issues/64-…md:44` 那框逐字写着「`:101-109` 轮询 poll tick 后断 **`len(live)==4`** && live[hkSummon].VK=='R'」——今天码里是 `==3`。**这不是放宽**（`VK=='R'` 那半句一字未动、`Rebinds()==1` 未动、注入新旧两键那段未动），**是别人票面上的引文过期了** ⇒ 归口见 §4 R-5，本腿不改票面、不碰框 |
| **⑤（票面未点名）`interaction_live_test.go` 那枚"归还后配置那枚应重新 live"** | 改前逐字：`if rep := b.HotkeyReport(); !rep.IsLive(hkCancel) { t.Errorf("after the B1 return the configured cancel binding is not live: %+v", rep.Bindings()) }` | 现读 `:244-250`：那三行**换成** `requireEscReturned(t, b)`（＝`EscTakenOver()` 必须 false ＋ `requireIdleRoster` 五件） | **不是放宽，是语义反转＋加强**：票 64 那枚断言的实质是"release 会把配置那枚重新绑回去"，而 245 裁的就是**不许再绑回去**（`hotkey_windows.go:547-553` 逐字「It is never re-bound here, and that is the point」）。⇒ 旧断言在本票裁定下**必须**反过来，否则它钉的是缺陷本身。新尺射程比旧的大：旧只判 `IsLive(hkCancel)`，新判 bool＋枚数＋三枚点名＋standby 行＋桌面探针。**并且同批发加了两处 `requireEscBorrowed`（`:205`／`:257`）与一处 `requireIdleRoster`（`:177`）**——只加不减。⇒ 编排者第 4 问"有没有顺手把票 64 的判据放宽"：**没有** |

⛔ **零词面型新仪器**（票面禁区）：`grep -rn "global hot\|four global" cmd/wisp/*_test.go` 现跑＝只命中 `resident_ball_live_228_windows_test.go:161` 那枚**错误消息文本**，没有任何"扫注释/扫文案里的键名"的门被造出来；判据一律问能力（`Live()` 集合、探针返回、`4/4` 那枚数来自子进程自报的 `len(rep.Live())`）。

### §2.5 AC#4／AC#5／AC#6（攻面 5）

**AC#4 用户能改 `[hotkey] cancel`——判：本格按票面形状如实写成"阻塞于票 228"，合格、不假绿。**
尺（我自己现读，不引它的表）：`cmd/wisp/resident_ball_windows.go:86` 逐字 `Hotkeys:  ball.DefaultHotkeys(),` ⇒ 出厂常驻腿**不读 `[hotkey]`**，用户改 `config.toml` 到不了 Win32；票面 AC#4 自己写了"若 228 那格还没做，本格判不了就写成'阻塞于票 228'，不许假绿"。⛔ 本腿复跑确认**没有第二条配置通路被造出来**：`grep -rn "config.Manager|hotkeysFromConfig" cmd/wisp/*.go`（剥 `_test`）现跑＝命中三行，其中**类型只出现在 `cmd/wisp/run.go:228`**（`mgr *config.Manager`，那是 `wisp run` 那条腿）；常驻腿里建球那一处逐字 `Hotkeys:  ball.DefaultHotkeys(),` ⇒ 配置到不了这枚球。

**AC#5 残余登记——判：成立，且比写腿交回时又多落了一处。**
- 票面末尾五字段：`.scratch/wisp/issues/245-…md:72-78`（`grep -c "残余登记"`＝3 处命中，含节标题）。
- 台账 `A478`（`docs/reports/pending-and-issues.md:9933` 起，`9943` 行逐字写「⛔ 未加 `DEFERRED(D-xx)` 代码标记……**欠账具名**：要真进那张总表，需要 owner 一句『搬』」）。
- **owner 已经说了那一句**：`dc694430`（17:5x）落的 `docs/specs/SPEC-12-roadmap-governance.md:80` 现读逐字有那行 `| DEFERRED | 确认窗口借用裸 Esc（票 245 乙形残余） | … |` ⇒ 五字段这条推迟项**现在三处齐全**。⛔ 本腿只 grep 了那一行做凭据，`docs/specs/**` 一字未改（写面自证见 §6）。
- 双向核对未破：`grep -rn "DEFERRED(D-" internal cmd --include=*.go` 现跑＝**零命中**（＝没找到＝没有"有标记没登记"的形状）。

**AC#6 两半——我独立裁编排者那句自裁，不当它是定论。**

① **"四枚"文案**：`cmd/wisp/main.go:25` 现读逐字仍在 `const usage` 反引号原始串里：「floating ball window, its tray icon and its **four global hot** / keys in this same process」，而我**真跑了这一版二进制的 `-h`**：`PATH=… "$TEMP/wisp245v1/wisp-clean.exe" -h` 输出第 5-6 行逐字含「its four global hot keys in this same process」；同一枚二进制的稳态实测＝**`hotkeys live 3/4`（§2.1／E3）**。⇒ **这句今天确实是半假话**，245-r1 没改它（它把它归票 228 的文案面，见其 N#2；编排者 A478 归票 245）。**本格不该勾＝编排者自己写的那句我复认成立**。修法约束我也核了一枚：`cmd/wisp/leg_dispatch_gate_133_test.go:156` 的 `usageBare133 = regexp.MustCompile("(?m)^  wisp {2,}\\S")` 钉住那一行的**行首形状**，下一程改文案时必须保留 `  wisp   ` 起头（票 228-r1 已经为此留过一句），⛔ 不许为改文案去动那枚钉。
　　**归口分歧（不属我能定）**：写腿说归 228（A477 把 `main.go` 那族句子判归 228），编排者说归 245（"是 245 把枚数改了的"）。我判：**编排者这一句更对**——句子被本轮改动变成假话，按票 228-v1 已确立的姿势"过期指认的普查必须连自己被改动影响到的数字一起重读"，账应落在造成变化那一票（245）名下。但这一处**只是归属，不是缺陷**，两读都不影响任何 AC 的成立性。见 §4 R-5。

② **判据锚该按"有无修饰键"，不是枚数**：**我裁编排者这一句成立**，三条凭据：
 - **形状凭据**：真缺陷的形状写在票面「现量」那一节的第一行与第三行——`RegisterHotKey` 把**裸键**从整个桌面拿走。带修饰键的 `Ctrl+Alt+V` 不是这个形状（它不与任何应用的自然按键冲突），把它一起 standby 掉**不是"更安全"，而是少一枚用户看得见的热键**，与 `ApplyHotkeyDefaults` 注释里在册的方向逐字相反（`hotkey_windows.go:78-80`：「a live hotkey the user can see, never a dead silent one」），而 `cancelIdleLine`（`:428-449`）对**任何**可解析的 cancel 绑定一律 standby——它不看 `acc.Mods`。
 - **实测凭据（E5，本腿现跑）**：`balldebug.exe -config <[hotkey] cancel = "Ctrl+Alt+V">` 打印逐字 `balldebug: hotkeys live=3/4 summon=Ctrl+Alt+Q live` ＋ 1 秒轮询那次重绑后逐字 `msg="cancel hotkey left unbound while idle: Esc is borrowed only during Confirming and handed back at session end (ticket 245)" hotkey=cancel binding=Ctrl+Alt+V` ＋ `msg="ball: hotkeys rebound after config change" … cancel=Ctrl+Alt+V … live=3`。⇒ **"用户配了组合键也不绑"今天已经真实发生在读配置的那条腿（调试旁支）上**，不是纯理论。
 - **钉子形状凭据**：`hotkey_live_test.go:435`（`requireIdleRoster` 里那枚枚数尺，整段 `:429-452`）判的是「`len(rep.Live()) != 3` 即红＋`c.Status != HotkeyStandby` 即 Errorf」，`live_windows_test.go:74` 判 `!= 3`，`hotkey_status_test.go:181/245` 判 `!= 3` ⇒ **四枚钉全锚在枚数上**。按 ② 的裁定改实现后（组合键常绑），这几枚今天绿着的钉会**同时红**，而实现方为了过钉最省的做法就是"继续不绑"——正是编排者说的"一枚今天的绿掩盖明天的岔路"。⇒ ② 的完成判据（钉子改成"稳态注册集里不存在无修饰键的那一枚"＋正控两形各一发）**尚未做**，本格不该勾；且它应当**先改钉、后改行为**（否则行为一改就红一片）。见 §4 R-6。
 - ⚠ ② 同时意味着**今天盘上的实现与最新裁定不一致**（N#5 那一支），但今天**零出厂行为差别**（出厂腿不读 `[hotkey]`＝AC#4 阻塞；默认档只有裸键一支），所以我不把它判成 AC#1 的缺陷，判成 **AC#6② 未完成＋一条要人在 §4 里认的账**。

### §2.6 CI 那一道门存不存在（攻面 6）

**结论：不存在。票 245 落地的全部真机判据今天只有这台机器跑过。**

- 尺一（tag 在场与否）：`grep -rn "winlive" .github/workflows/ scripts/` ＝ **零命中**（＝没找到＝CI 与全部脚本里没有任何一处带 `-tags winlive`）。
- 尺二（这些文件到底参不参与编译）：四枚 live 文件头逐字 `//go:build windows && winlive`——`internal/ball/live_windows_test.go:1`／`internal/ball/hotkey_live_test.go:1`／`internal/ball/interaction_live_test.go:1`／`cmd/wisp/resident_ball_live_228_windows_test.go:1`。CI 的 windows 腿跑的是 `bash scripts/portable-tests.sh --scope=windows`（`.github/workflows/ci.yml:512`），scope 里 `internal/ball` 在册（`scripts/portable-tests.sh:129/153/187`），但**不带那枚 tag ⇒ 这四枚文件在 CI 上根本不进测试二进制**，`requireIdleRoster`／`requireEscBorrowed`／`requireEscReturned`／新增那枚 `hotkeys live 4/4` 判红**在 CI 上一次也不会执行**。`cmd/wisp` 那一格同理（`.github/workflows/ci.yml:476` ＋ `scripts/wisp-cli-tests.sh`，两处都无 winlive）。
- 尺三（CI 有没有跑过这轮码）：`gh run list --limit 5` 现跑 ⇒ 最近五枚是 `36674742682`（slo-fresh，05:43:47Z）／`36664144050`（ci，03:23:04Z）／`36643624549`（ci，2026-09-29T23:09:03Z）／`36639396577`／`36591499572`。本票四枚 commit 落在 **09:06–09:45Z**（本地 17:06–17:45 +08）⇒ **全部早于本轮码，CI 至今没有跑过 245 的任何一枚 commit**，也说不出任何 run id ＋ step 名来背书。
- ⇒ 按派单第 3 节第 6 条的判据形状（"说不出 run id＋step 就当门不存在"）：**这道门今天不存在**。这与本仓已知形状同族（票 62 AC#9 那条"只在 tag 后的读数 CI 从不执行"）。**残余登记见 §4 R-3**：稳态枚数、借还两头、`4/4` 判红这些**能力型**判据没有任何 CI 面覆盖，回归全靠下一次有人手动跑 winlive。
- ⚠ 这一条**不把 AC#1/AC#2/AC#3 翻成不成立**：票面本来就允许"本机现跑真机读数"当凭据（AC#6① 逐字「凭据要现跑那发 `winlive` 读数」＝本仓承认这一档只在这台机器上取），而 CI 覆盖从来不是本票任何一格的完成判据；它是**本票落地后的持久风险**，不是一次交付缺件。

### §2.7 突变名册（本腿自跑三发，逐发：改的那一行／红的用例名／红句逐字／还原复认）

姿势：每发先 `grep -c` 锚点＝**1**，落盘，跑指名的用例（判红绿只认 `--- FAIL`），跑完 `git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet -- internal/ball cmd/wisp` 复认。⛔ 未用 `-overlay`、⛔ 未用 `checkout`／`stash`／`reset`。台件全文：`.scratch/wisp/probes/245/v1/m1-*.txt`／`m2a-winelive.txt`／`m2b-winelive.txt`。

| 发 | 改的那一行（逐字前 → 逐字后） | 指名的尺 | 我自己读到的红 | 还原复认 |
|---|---|---|---|---|
| **V-M1**（＝派单攻面 3 要的那一发："把稳态那支反回去，让它真绑裸 Esc，会不会红"） | `internal/ball/hotkey_windows.go:457`：`if p.id == hkCancel {` → `if false && p.id == hkCancel { // MUTATION V-M1`（＝稳态注册那支永不走 standby 分支，四枚全绑，回到 245 之前的出厂形状） | 默认层五枚＋球真机 trio＋出厂进程真机一枚 | **三层全红。**默认层 `--- FAIL` 5 枚，红句逐字：`hotkey_status_test.go:120: attempted [1 2 3 4], want exactly the three idle slots [1 2 4] (cancel is standby, ticket 245)`、`hotkey_status_test.go:179: the production defaults bound a bare Esc globally: [{16387 81} {16387 77} {16384 27} {16387 80}]`（`{16384 27}`＝`MOD_NOREPEAT`＋`VK_ESCAPE`＝**那枚被吞掉的键本身**）。球真机层 `--- FAIL` 4 枚，红句逐字：`live_windows_test.go:75: the live ball did not register its three idle hotkeys (summon/mute/panel, cancel is borrowed only during Confirming)` ＋ `hotkey_live_test.go:149`／`:309`／`interaction_live_test.go:177` 三处同一句 `ticket 245 RED: the idle ball holds the cancel hot key`（报告里 `cancel ID:3 … Status:live`）。出厂层 `--- FAIL: TestLive228ResidentLegOwnsABallWindowOnTheDesktop (3.27s)`，红句逐字 `resident_ball_live_228_windows_test.go:164: AC#1 RED (ticket 245): the idle resident ball registered all four hot key slots, so the cancel slot is a desktop-wide hot key while nothing is being confirmed.` ＋同一发子进程自报 `hotkeys live 4/4: … cancel=Esc live`。**交付层另有一发（§2.3b-E4）**：这一版造出的 `wisp-mutant.exe` 挂着时，第二个真进程收到的 `keydown_esc=0` | `RESTORE bytes=21395`；`git diff --quiet -- internal/ball cmd/wisp` rc=**0**；`grep -c "MUTATION V-M1"`＝**0 枚**（＝标记不在盘上）；还原后复跑：默认层五枚 `--- PASS`／`ok 0.027s`，球真机 trio `ok 2.670s`，出厂真机 `ok 3.428s` |
| **V-M2b**（＝"谎报形"：记账说借到了、Win32 从没被问过一次——**这发专门审那把探针**，派单攻面 1） | `internal/ball/ball_windows.go:878`：`if err := takeEsc(b.hwnd); err != nil {` → `if err := error(nil); err != nil { // MUTATION V-M2b`（`escTakenOver` 与 `cancelBorrowedLine()` 照写，**一次 Win32 都不问**） | 球真机 trio | `--- FAIL` **3 枚**，红句逐字：`ticket 245 RED: the ball says it borrowed Esc, but the desktop still has it free - the borrow registered on some other key, not on VK_ESCAPE`（命中点 `hotkey_live_test.go:311`／`interaction_live_test.go:205`／`live_windows_test.go:127`）。⚠ 同发 `interaction_live_test.go:200` 那枚**只查 bool** 的断言是**过**的（`EscTakenOver()` 真的变 true 了）⇒ 这一发证明：**探针不是装饰**，去掉它、只留记账，这一形就溜过去了。`TestLiveHotkeyRebindEndToEnd` 同发仍 `--- PASS`（它本来不借 Esc，射程外，如实报） | `RESTORE bytes=31709`；`git diff --quiet` rc=**0**；`grep -c "MUTATION"`＝**0 枚**；还原后 trio `ok 2.670s` |
| **V-M2a**（＝"根本不尝试借用"形） | `internal/ball/ball_windows.go:875`：`if b.escTakenOver {` → `if b.escTakenOver \|\| true { // MUTATION V-M2a`（进 `TakeEscForCancel` 就早退，不碰 Win32 也不写记账） | 球真机 trio | `--- FAIL` **3 枚**，两种红句都出现，逐字：`interaction_live_test.go:201: entering Confirming did not take Esc over (B1)`（等 bool 那一枚，2.06s 后花完预算）；`hotkey_live_test.go:311`／`live_windows_test.go:127`: `ticket 245 RED: Confirming never attempted the Esc borrow (cancel line still standby, live set holds 3 of the 4 slots)` | 同上：`RESTORE bytes=31709`、`git diff --quiet` rc=0、`grep -c "MUTATION"`＝0 枚、还原后 trio `ok 2.670s` |

⛔ **本腿没有修任何产码**：三发全部还原，终态自证在 §6。⚠ 本腿**没有**复现写腿表里 M3 那一发（它自己判"不需要另开一发"，我读完 `TestCancelBorrowRoundTrip`（`hotkey_status_test.go:236-243`）认同：`releaseEscWith` 只有 `unreg` 一支、`attemptsOf(hkCancel) != 1` 那枚断言直接钉住"还＝再注册一次"，确定性层已覆盖，我未重复。

### §2.8 门禁复跑（本腿自己现跑，⛔ 不背写腿的数字）

| 尺 | 我的读数 |
|---|---|
| `go build ./...` | **rc=0，零输出**（`.scratch/wisp/probes/245/v1/gates-static.txt`） |
| `go vet ./internal/ball/` | rc=0 零输出 |
| `go vet -tags winlive ./internal/ball/` | rc=0 零输出 |
| `go vet ./cmd/wisp/` | rc=0 零输出 |
| `go vet -tags winlive ./cmd/wisp/` | rc=0 零输出（这一档写腿表没交，本腿补上） |
| `gofmt -l` 对它落的八枚文件（`hotkey_windows.go`／`ball_windows.go`／`hotkey_status_test.go`／`hotkey_live_test.go`／`live_windows_test.go`／`interaction_live_test.go`／`resident_ball_windows.go`／`resident_ball_live_228_windows_test.go`） | **零命中＝没找到＝干净** |
| `bash scripts/d22scan.sh` | 第一步**正控先吃读数**：`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ⇒ 门是活的；真扫：`d22scan: examined 253 production Go files under internal/ and cmd/`，`ban #8 internal/ 473 枚`／`cmd/ 69 枚`，末行 `d22scan: clean - no D22 ban violations`，rc=0 |
| `go test ./internal/ball -count=1`（默认层整包） | **1 枚红**：`--- FAIL: TestC21TableColourRowsMatchTokensCSS` ＋ `tokens_table_test.go:1468: read design/assets/tokens.css: … The system cannot find the path specified.` ⇒ **归因：起手即在、非本票**（起手名册第 14 行 ` D design/assets/tokens.css`；`git diff --name-only 25ce3ce9 60f888c8` 里 `tokens_table_test.go` 命中 **0 枚**，该文件末次 commit 是 `ed43bc3d`＝票 74）。本票那族默认层用例逐枚绿（§2.1 那一发 5 枚 PASS） |
| `PATH=… go test ./cmd/wisp -count=1`（整包默认层） | ⚠ **三发里第一发红、后两发绿，且第一发的红名丢了**：第 1 发（不带 `-v`，我只 `tail -20`）＝`FAIL github.com/CarlosShao/wisp/cmd/wisp 160.620s`、**用例名未落进我留的那 20 行**；第 2 发（`-v` 全量）＝`ok … 162.117s`、`--- FAIL` **0 枚**；第 3 发（`-v` 全量）＝`ok … 163.001s`、`--- FAIL` **0 枚**、`=== RUN` 201 枚。⇒ **如实记：1 红 2 绿、红因未定案，不能报成"全绿"**。已排除的一项：这三发之间我没有并发过别的 `go test`（d22/内部 ball 那发在 17:51 就结束了），但第 1 发期间我在 17:56 编译过一次台件、CPU 有争用；`cmd/wisp` 这一包大量用 `waitFor` 秒级预算（票 228/127 那族真起子进程），时序型抖动是本仓已知形状。**这一枚欠一次归因**，进 §4 R-7 |
| `-tags winlive ./internal/ball`（指名 trio，还原前后各一发） | 还原前基线：4 枚 `--- PASS`、`ok 2.638s`、`SKIP-LOUD` 0 枚、`no tests to run` 0 枚；三发突变还原后：`ok 2.670s` |
| `-tags winlive ./cmd/wisp -run TestLive228…` | 基线 `--- PASS (3.10s)`／`ok 3.155s`（读数＝§2.1 第一行）；V-M1 下 `--- FAIL`；还原后 `ok 3.428s` |
| ⛔ 未跑 | 整包 `-tags winlive ./...`（两枚 live 套件共用桌面会互洗读数）；`ball-cycle.ps1`；`gofumpt`（本机无该工具＝与 CI 那一档差一级，见 §5）；`staticcheck`；`GOOS=linux go vet ./cmd/wisp`（写腿报的是 sherpa cgo 构建约束、票 78 地界，本腿**未复跑该尺**，不替它背书也不推翻） |

---

## §3 攻不动的清单（写了就是没攻，不硬判；也写下"它的哪些主张我攻不动＝复现不了"）

1. **`Confirming` 借期那一 2–3 秒里的投递，本腿没测到**。E2／E4 证的是"**一枚全局裸 `Esc` 注册确实把键从焦点窗口拿走**"，注册者分别是观察者自己和 M1 突变后的出厂进程；**球真在借期的那一发**我没造出来——出厂常驻腿没有卡片路径（§2.3b 上方那串 `grep`），而 `cmd/balldebug -state Confirming` 那一支逐字只做 `b.SetState(s)`（`cmd/balldebug/main.go:285-289`），**不调 `TakeEscForCancel`**（那只在 `syncMachineToBall`（`:630-639`）里由 `dispatch` 驱动）。⇒ 借期吞键今天的凭据是**读码链**：`takeEscWith`（`hotkey_windows.go:531-539`）与 E2 是同一发 `RegisterHotKey(modNoRepeat, VK_ESCAPE)`，差别只在注册者。⛔ 我不把它写成实测，也不许别人把它写成实测。
2. **`escBorrowProbe` 照不到的那一形攻不动**：低级键盘钩子（`WH_KEYBOARD_LL`）不占热键表也能吞键。本腿没有造钩子程序，所以"探针回 free 而用户仍收不到 Esc"这一形**未被排除**（现实中谁会这么干是另一件事）。⇒ 这一条留在 AC#2 的"附条件"里，不升成不成立。
3. **真实第三方占位者复现不了**：本机稳态下裸 Esc 无人占（探针回 free、E1 收得到）。E6 用的是我自己造的进程当外来占位者，**不是**某个真的 IME／截屏软件。⇒ E6 那两笔读数（探针能分辨／新钉跨机器假红）在"形状"层成立，在"这台机器之外还有多少人会遇到"层我没有数。
4. **写腿表里那三发 M1／M2a／M2b 的原始读数我一枚都不采信**（同一人所作），但**我用等价姿势自己跑过**（V-M1／V-M2a／V-M2b，§2.7）⇒ 这一条不是攻不动，是"重跑过"。真正没重跑的是它的 **M3**（我读完 `hotkey_status_test.go:236-243` 与 `releaseEscWith`（`hotkey_windows.go:553`）后**认同它"不需要另开一发"**）与 **`0/4` 那一枚**（要的是"四枚全注册失败"，非本票形状，我未复现）。
5. **CI 那台机器有没有可建窗口的桌面**：判不动（本票四枚 commit 之后一次 run 都没有，见 §2.6 尺三）。
6. **`gofumpt` 那一档**：本机无该工具，且装它要跨网络拉（历史代价）⇒ 与 CI 差一级，如实报（§5 第 12 条）。
7. **票 64 剩下那几枚未勾框（多显示器实拖等）**与本票无关，本腿一枚未碰、未裁。
8. **`design/**`／`frontend/**` 两层禁读**：我只用 `tokens_table_test.go:1468` 那条"文件不存在"的错误消息做归因，**没读过那目录里任何一个文件**，也不转述其内容。

---

## §4 留给编排者的（要人拍的／要归口的，一条一枚；⛔ 本腿一枚未改别人的地界）

| 号 | 一句话 | 凭据 | 我建议的归口与形状（**不是判语，等你拍**） |
|---|---|---|---|
| **R-1** | **N#4 是真的，但它给的理由是假的**：`RebindHotkeys` 丢掉在飞的借用这一形，今天**在 `cmd/balldebug` 里可达**（同进程既有卡片驱动又有 1 秒配置轮询），而注释逐字写着"没有调用者、因为轮询与卡片不在同一条腿" | `internal/ball/ball_windows.go:809-812`（那句假话）／`cmd/balldebug/main.go:237-243`（轮询桥）／`:630-639`（卡片）／`ball_windows.go:822`（`escTakenOver = false` 无补借）。出厂侧不可达我已复认：`grep -rn "NewHotkeyReloader" cmd/wisp/*.go` ＝**零命中** | 归**票 228 后续片**（卡片接进常驻腿时一并处理）或直接立一枚小票。最小修法两式：重绑后若 `EscTakenOver()` 为真就补借一次；或重绑前记借状态、重绑后照记恢复。**本腿不修**。至少 `:809-812` 那三行注释今天必须改成带条件事实句（假注释比没注释更坏） |
| **R-2** | **AC#2 的"别的程序要能收到"在仓里没有任何判据**。我造的第二进程能测，但它不在仓里 ⇒ 这一半的回归今天无人守 | §2.2 与 §2.3b 全部；台件 `.scratch/wisp/probes/245/v1/esclistener/`（只建不删，可直接搬进 `internal/ball` 当 winlive 辅助） | 建议立**下一程写腿**的活：把"第二枚真进程收到／收不到 `WM_KEYDOWN(VK_ESCAPE)`"钉成 `-tags winlive` 用例，**判据两形各一发**（稳态必收到＝正；自占裸 Esc 必收不到＝反），⛔ 不许只钉正向。要人拍的点：这类用例需要一个能建窗口**且允许抢前台**的桌面，跑法与本仓 `requireQuietBallDesktop` 冲突时算谁的 |
| **R-3** | **票 245 的全部真机判据没有 CI 覆盖**，而这轮落地的恰恰只有真机判据（稳态枚数、借还两头、`4/4` 判红） | §2.6 三把尺（`grep winlive` 零命中／tag 逐字／最近五枚 run 全早于本轮码） | 与票 62 AC#9 同族。两条路：①承认"这一族只在这台机器上跑"，在票面／派单里把它写成**每次改热键稳态必手跑**的硬规矩；②立票给 winlive 找一档有桌面的 runner。⛔ 本腿不造门、不改 `ci.yml` |
| **R-4** | **V-F2：新钉在"别的程序已经占了 Esc"的机器上会把正确实现判红**，且因此让 `SKIP-LOUD` 那一支永远走不到（写腿 §5 J#6 的说法被实测推翻） | §2.3b-E6 读数：四枚硬红、红句逐字 `ticket 245 RED: bare Esc is claimed on this desktop while the ball is idle`（`hotkey_live_test.go:149`／`:309`／`interaction_live_test.go:177`／`live_windows_test.go:78`），`SKIP-LOUD` 0 枚；根因 `hotkey_live_test.go:445-451`（探针不分辨"谁在占"）＋调用次序（`requireIdleRoster` 在任何借用之前） | 归下一程写腿：`requireIdleRoster` 里把"探针不 free"按球自己的报告分两支撑——`rep.IsLive(hkCancel)` 为真＝**球占着＝红**；为假＝**别人占着＝`SKIP-LOUD` 具名**（与本文件 `:99-104` 已有姿势同形）。这不是放宽：⛔ 判"球占着"那一支一字不动 |
| **R-5** | **别人票面上的过期指认**：票 64 已勾框里逐字引着「断 `len(live)==4`」，今天码里是 `==3`；另 `cmd/balldebug` 那张进度日志里 `hotkeys live=4/4` 那句也是旧读数 | `.scratch/wisp/issues/64-ball-defects-hotkey-interactive.md:44`（已勾 `[x]` 那框正文）对照 `internal/ball/hotkey_live_test.go:174`；`cmd/balldebug/main.go:224` 现跑打 `live=3/4`（E5 那一发逐字读到的） | 归口＝**编排者**（框与票面都是你的）。⛔ 本腿没动任何框、没动 64 的票面。建议：在 64 名下追加一行"该引文由票 245 改成 `==3`，判据实质（改配置→新键生效→旧键失效）未变"，别让它读起来像有人把票 64 放宽了 |
| **R-6** | **AC#6② 的裁定会改变 AC#1 的字面**（"稳态只能是三枚" → "稳态不得有裸键"）。今天零出厂行为差别，但**盘上实现与最新裁定已不一致**（带修饰键的 `cancel` 也不绑），且钉子全锚在枚数 | §2.5 三枚凭据；E5 实测逐字 `msg="cancel hotkey left unbound while idle … " binding=Ctrl+Alt+V`＋`live=3`；钉：`hotkey_live_test.go:435`／`live_windows_test.go:74`／`hotkey_status_test.go:181`/`:245` | 要拍的点只有一个：**认不认"带修饰键的 cancel 常绑"这句**。认，则**顺序必须是先改钉（判有无修饰键＋正控两形）再改行为**，且 AC#1 票面那句话要重写成带条件式；不认，则回退 A478 里 ② 那半句并写明理由。⛔ 本腿不动实现、不动钉子 |
| **R-7** | **`cmd/wisp` 整包默认层本腿读到 1 红 2 绿，首发用例名已丢**（我只留了 `tail -20`）——"门禁全绿"这句在我自己读数下不成立 | §2.8 那一行读数（160.620s 红／162.117s 绿／163.001s 绿，`--- FAIL` 后两发各 0 枚） | 归**下一程**（谁再动 `cmd/wisp` 或本票后续片）：整包重跑并**全量留档**，把那枚红归因成"时序抖动"或"真缺陷"。⛔ 我不把它算成 245 的缺陷：后两发同样在 HEAD 上跑绿、且本腿三发突变还原后 `internal/ball` 两档与 `cmd/wisp` live 档全绿。教训自记：**整包复跑一律 `>` 全量落盘，不许 `tail`** |
| **R-8** | **AC#6① 那句 `wisp -h` 正文"four global hot keys"仍未修**（半假话，且是打给用户看的第一句） | `cmd/wisp/main.go:25`＋本腿真跑的 `-h` 输出（§2.5）；形状钉 `cmd/wisp/leg_dispatch_gate_133_test.go:156` | 归口分歧：写腿说票 228，编排者说票 245。**我裁编排者对**（造成枚数变化的是本票）。改法约束：保留 `  wisp   ` 行首、改成"稳态三枚＋`Confirming` 借第四枚"的带条件事实句、⛔ 不新增词面型仪器 |

---

## §5 判不动与没做的（硬预算闸门：这一节必须写满；本腿第 100 轮帽**未触顶**，交件前已写满）

**判不动（要外推才有的结论，我不给）**

1. **"用户现在可以放心用 Esc 了"——不给**。三层都在：①乙形那 2–3 秒的借用窗**仍在**（`SPEC-12:80` 那行 DEFERRED 写的就是它，V-M1/E2/E4 证明这一发的确吞键）；②出厂常驻腿今天没有卡片路径，所以"借"在生产里还没发生，等票 228 后续片把卡片接进来那天，本表的 AC#2 判语要**重裁**；③低级钩子那一形我没排除（§3 第 2 条）。⇒ 本表能说的最大一句是：**稳态（没有卡片挂着时）出厂进程不再从别的程序手里拿走裸 `Esc`，这一条我在交付层实测到了（E3）**。
2. **E6 那枚假红在别的真机器上会不会发生**：形状成立、发生率我没有数（§3 第 3 条）。
3. **CI 那台 runner 有没有桌面**：判不动（§3 第 5 条）。
4. **`0xc0000135` 那类环境性无效读数**：我一律只认 `--- FAIL` 判红绿，`build failed` 与零 `--- FAIL` 的 rc≠0 都不算证据；本表全部读数符合这一姿势（三发突变的红都有 `--- FAIL` 行）。

**没做的（列全，不遮掩）**

5. **没跑整包 `-tags winlive ./...`**（两枚 live 套件共用桌面会互洗读数）；只跑指名的 `internal/ball` trio（4 枚）与 `cmd/wisp` 的 `TestLive228…`（1 枚），也**没有**跑 `internal/ball` 里其余 live 用例（`TestLiveMuteHotkeyEndToEnd`／`TestLiveSleepingZeroTimerHandles`／`TestLiveNeverStealsFocus`／位置持久化那族）——它们不属本票射程，但因此**本票落地的三枚真机改动我只在指名四枚上取过数**。
6. **没复现 `0/4` 那一枚**（四枚全注册失败形，非本票形状）；**没复现写腿的 M3**（读完断言后认同它"不需要另开一发"）。
7. **`cmd/wisp` 整包默认层那一发红因未定案**（首发用例名已丢，见 §4 R-7），本腿**没有**为归因再跑第四发。
8. **没测 `Confirming` 借期那一窗的投递**（造不出便宜驱动者：`balldebug -state` 那一支不调 `TakeEscForCancel`，`cmd/balldebug/main.go:285-289`；出厂腿无卡片路径）。§3 第 1 条给了替代链与其折扣。
9. **没测低级键盘钩子那一形**、**没测多显示器／跨会话／UAC 提升桌面**下热键表的差异。
10. **零性能、零句柄、零 RSS 判定**：本票六格没有一条是性能判据，`live_windows_test.go` 自带那枚 <600 句柄门我照跑未另评。
11. **没读 `design/**` 与 `frontend/**`**（两层禁读），只用 `tokens_table_test.go:1468` 的错误消息把那枚常红归因成"起手即在、非本票"。
12. **`gofumpt`／`staticcheck`／`slo-*`／`GOOS=linux go vet ./cmd/wisp` 一档未跑**：本机无 gofumpt；其余不属本票射程。⇒ 门禁面比 CI 松的那一级如实报，不谎称对齐 CI。
13. **没碰任何 AC 框、没写台账、没改 `PLAN.md`／`docs/specs/**`／SLO／阈值／golden／allowlist／三枚冻结件**；没修任何产码（三发突变全部还原，§2.7 逐枚 `git diff --quiet` rc=0）；终态名册对撞见 §6。
14. **AC#4 之外没有第二处"阻塞"需要我裁**；`SPEC-12 §5` 那一行本腿只 grep 了一行做凭据（读一行不算改动，写面自证在 §6），**没读那张表的其他内容**，⛔ 也未据其内容做任何判定。

---

## §6 名册核对（终态＝起手那一刻的名册，共享树里不写"必须为空"）

- **起手**（17:47 +08，锚点 `a43fe78b`）：`git status --porcelain` 全量 **175 行**逐字存 `.scratch/wisp/probes/245/v1/start-status.txt`；其中 `internal/ball`＋`cmd/wisp`＝**0 枚**（＝没找到＝干净）；`docs/evidence/s1/` 只有别人的 ` M 152-subject-death-never-measured-r1-accept-r1.md`（起手即在、本腿未碰、未提交）与我自己的骨架那行 `??`。
- **终态**（18:17 +08，最后一次量）：同一条尺跑两次取同一份（`.scratch/wisp/probes/245/v1/end-status.txt`，**176 行**）。
- **差分（`diff start end` 逐字，恰好两行改动）**：
  - `167d168` ＋ `31a32`：同一枚文件的两级变化——`?? docs/evidence/s1/245-esc-not-a-standby-global-hotkey-v1.md`（起手，未跟踪）→ ` M docs/evidence/s1/245-esc-not-a-standby-global-hotkey-v1.md`（终态，本腿已提交两枚、正在写第三批）。⇒ **本腿自己的落点**。
  - `128a130`：`?? .scratch/wisp/probes/orchestrator/msg-a479.txt` ⇒ **编排者在飞行中新增的台件**（他 A479 那条 commit 的信息件），**本腿没碰、没提交、没删**。⛔ 别人留下的 172 枚条目一枚未动。
- **产码面**：`git diff --quiet -- internal cmd` ＝ **rc=0**（三发突变全部还原，逐枚 `RESTORE bytes=…`＋`grep -c "MUTATION"`＝0 枚在 §2.7）；`git status --porcelain -- internal/ball cmd/wisp` ＝ **0 枚**。
- **禁地写面**：`docs/specs/**`／`docs/PLAN.md`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`frontend/**`／`design/**` ＝ **零改动**（终态名册里没有任何一枚这些路径的 ` M`／`A`／`D` 由本腿产生；`design/**` 那批 ` D` 是起手名册里别人留下的原样）。
- **票面**：`.scratch/wisp/issues/245-…md` 终态＝**未修改**（`git status --porcelain -- .scratch/wisp/issues/245-…md` 零命中）；AC 框计数尺（我自己跑）：`grep -c "^- \[ \]"`＝**6 枚未勾**、`grep -c "^- \[x\]"`＝**0 枚已勾**，与起手一致 ⇒ **本腿一枚框都没动**（AC 框归编排者）。
- **本腿落盘清单**（`issues/README` 规则 8：只建不删，全部在 `.scratch/wisp/probes/245/v1/`）：`msg-skeleton.txt`／`msg-fill1.txt`／`msg-fill2.txt`／`start-status.txt`／`end-status.txt`／`gates-static.txt`／`gate-d22-and-ball.txt`／`gate-cmdwisp.txt`／`gate-cmdwisp-full.txt`／`gate-cmdwisp-rerun.txt`／`winelive-ball-baseline.txt`／`winelive-cmdwisp-baseline.txt`／`m1-default-tier.txt`／`m1-winelive-ball.txt`／`m1-winelive-cmdwisp.txt`／`m2a-winelive.txt`／`m2b-winelive.txt`／`e1-baseline-watch.txt`／`e2-steal.txt`／`e3-listener-vs-clean.txt`／`e4-listener-vs-mutant.txt`／`e6-steal-holder.txt`／`e6-winelive-vs-foreign-owner.txt`／`esclistener/{go.mod,main.go}`。⛔ 未删除任何文件（含别人的条目）；`$TEMP/wisp245v1/` 下的 exe／数据目录在仓外，同样不清理。
- **本腿 commit（三枚，全部显式 pathspec，只 commit 未 push）**：`2f2d398e`（骨架，17:4x）→ `6f5595fd`（§0／§1／§2.1／§2.3b／§2.4／§2.5／§2.7／§2.8）→ 本枚（§3／§4／§5／§6 ＋ §2.0／§2.2／§2.3／§2.6 早批）。⛔ 未 `push`；未用 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
- **交件前自查（占位符尺）**：对本文跑那四枚占位词的 grep（「待」＋「填」连写那一枚、「填写／中」那一枚、三个英文字母那枚、以及骨架期我自己用的那枚两字占位词）⇒ **零命中**（＝没找到＝没有留空节；这一行故意不把它们连写，否则尺会命中它自己）。逐节自查同时做过：§1 十二行判语无一行是占位、§5 十四枚条目全部写满。
