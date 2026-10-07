# 票 272 对抗验收裁决件 — 腿 `272-v1`（非实现者）

> 本件是**对抗验收**，不是复述实现件 `272-r2`。实现者是另一枚腿 `272-r2`（已停），
> 那 52 行不是我写的。我的职责是攻它。任何一格我**不**因为它自陈"做了"就判成立。
> 结构：§0 起手锚｜§1 六格逐格判语｜§2 恒真性进攻｜§3 我推翻实现件的哪几句｜
> §4 我推翻编排者的哪几句｜§5 未做完的格｜§6 终态自证。
> ⛔ 本腿零产码改动、零 AC 框翻动。

---

## §0 起手锚（第 1 笔，任何长跑命令之前落盘）

原文读数（本腿 2026-10-07 自跑）：

```
$ date
Wed Oct  7 11:02:56 CST 2026

$ git rev-parse --short HEAD
e18e32da            # 分支 dev

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
（空 —— 六族路径起手全干净，无别人的未提交 Go/脚本/文档改动）

$ git log --oneline -3 -- cmd/wisp/config_reload_223_test.go
8d30a862 验 5g2-v1（非实现者只读腿·第 1 笔）：... [顺手提走那 52 行的编队事故笔]
3d9b8374 probes(232-r2 收尾代提): ...
a16d1ff7 票 231 AC#3＋同批名册 · 常驻钉加一行，cause=newer-build 进两张互斥名册
```

被审对象落点核对（本腿现量，非抄实现件）：

```
$ git show --stat --oneline 8d30a862 -- cmd/wisp/config_reload_223_test.go
 cmd/wisp/config_reload_223_test.go | 52 ++++++++++++++++++++++++++++++++++++++
 1 file changed, 52 insertions(+)

$ git diff --numstat 8d30a862^ 8d30a862 -- cmd/wisp/config_reload_223_test.go
52      0       cmd/wisp/config_reload_223_test.go        # 只增 52／0 删，与票面/实现件一致

$ git diff --numstat 8d30a862 HEAD -- cmd/wisp/config_reload_223_test.go
（空 —— 8d30a862 之后测试文件再无改动，故 HEAD 版 == 8d30a862 版 == 被审版）

wc -l HEAD 版测试文件 = 797
wc -l 8d30a862^（改前尺）版 = 745      # 差 52，对得上
$ git show HEAD:cmd/wisp/config_reload.go | md5sum
5ce441ca5e72b64d18a6c26f1c066882       # 起手时产码基线 md5，§6 还原自证以此为拉
```

结论性锚点：**"改前尺" = `8d30a862^:cmd/wisp/config_reload_223_test.go`（745 行）**；
"被审尺" = `HEAD:cmd/wisp/config_reload_223_test.go`（797 行）。二者差恰为那 52 行。

起手整包红名册（本腿自跑，`go test -v ./cmd/wisp ./internal/... -count=1`，raw 1.1 MB 留仓外
`/d/tmp/wisp272v1/start-raw.txt`，摘要 `logs/G0-start-rednames.txt`）：
四数 **RUN/PASS/FAIL/SKIP = 2174 / 1472 / 6 / 7**，rc=1。
具名红 6 枚（终态作差以此为基线）：`TestC21TableColourRowsMatchTokensCSS`／`TestApprovalCardViewJSONKeysMatchFrontendTypes`／
`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／
`TestC21DesignTokensFourWayAgree`／`TestResolvePerCallBudget`。
⚠ 与实现件 `272-r2` 起手名册（7 枚）差一枚：`TestAlwaysBranchStoresItsSecondL2Card` 类带载型这发**绿**（本机无 runner 抢 CPU）；
且 `internal/panel` 包 **FAIL 但不打 `--- FAIL` 名**（winlive 面板宿主线程环境性失败，非本票产码所致，`cmd/wisp` 包起手即 `ok`）。
⇒ 本票只动 `cmd/wisp` 测试文件，`cmd/wisp` 无红；其余红/包级失败均**起手就在**，与本票那 52 行无关。

---

## §1 六格逐格判语（每格带本腿自己那份读数 + 原始件路径；框仍由编排者翻，本腿一枚不勾）

判语三值：成立 / 不成立 / 判不动。凡"成立"都指**本腿现量**，非实现件自陈。

### AC#1 — 未改动的尺上 `M` 必须仍 PASS —— **成立**
本腿读数（raw `.scratch/wisp/probes/272/v1/logs/AC1-oldruler-M.txt`）：
`--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.34s)`。
尺 = `8d30a862^` 版（overlay 换回旧 745 行尺），产码 = `M` 形（`git show HEAD:cmd/wisp/config_reload.go` 于 `:127` 后插一行横幅，
复用 `272-r2` 的形定义文件 `mut/M/`，但**读数本腿自跑**）。⇒ 票面立案凭据复到：未修码的尺对 `M` 一声不响。
⚠ 本腿在旧尺上跑的 `M` 的 PASS 时刻(2.34s)与实现件(2.38s)非抄，是另一次自跑。

### AC#2 — 改完 `M` 必须红、红因＝具名断言句 —— **成立（但只达成"红"这一半，见下）**
本腿读数（raw `logs/AC2-headruler-M.txt`）：`--- FAIL (2.15s)`，红句整行（⛔ 未截断）：
```
config_reload_223_test.go:596: the operator stream carried "本次运行不会生效" 1 time(s)
BEFORE this run planted the edit; a copy nobody planted cannot be evidence that this run
will not use the new values. stdout before the plant: ...
```
＝ **arm1 的具名 `t.Errorf`**，非超时（2.15s << reloadCaseBudget 40s）、非 build failed（同尺 `cur` 编译并跑通）。
⇒ AC#2 的**硬要求**（"同一形必须红 + 红因是具名用例"）**成立**。

**但 AC#2 还含一枚手段子句"三形（甲／乙／丙）选一支落地"。本腿对它单独判语：**
- 票面 `:21` 逐字：「⛔ 形不在票面写死；三形是腿交回的候选，落地前自己再判一次」⇒ arm1 这枚**新造的臂不算超射程、不该退回**。
- 让 `M` 变红的是 **arm1（新臂，pre-plant 计数＝0）**，**不是**票面的甲（arm2）。实现件 §5.2 自称"甲 alone 在 M 上 PASS"，
  本腿用 **§2 进攻(i)** 独立复证：把 arm1 阈值 `!=0` 改成 `!=1` 后 **`M` 直接翻回 PASS(2.13s)** ⇒ M 的红**只来自 arm1**，
  arm2/甲 单独对 `M` **零牙**。故"三形选一"里能钉 `M` 的那一支**在票面枚举中不存在**；实现件用枚举外的 arm1 达成红，
  按 `:21` 合法，但"三形选一支"这句**字面不成立**（没有一支三形做到了，是第四形做到的）。

### AC#3 — 零放宽 —— **成立**
- `git diff --numstat 8d30a862^..HEAD -- ...test.go` = **52 插入 / 0 删除**；`grep -c '^-[^-]'` = **0**。
- 断言形状对拉（本腿现量，OLD vs HEAD 同一把 grep）：`t.Errorf` **31 → 33**（恰 +2＝arm1、arm2 各一条），
  `t.Fatal` 32→32、`Contains(out,needle)`/`Contains(why`/`Contains(whole`/`Contains(out,"本次运行不会生效")`/`windowCount` 计数**一字未变**。
  ⇒ 没删任何既有断言（比实现件自报的"41 枚不变"更强：新增只 2 枚 t.Errorf）。
- 六枚 needle 仍在 stdout 窗里：`:636 for _, needle := range []string{"app.autostart","开机自启","需要重启进程","原因：","交给平台层","没有被丢掉"}`
  读的是 `out`（= `awaitStdoutSince232(mark,...)` 窗），拼接串没回来。
- 窗口起点未前移：`mark := r.h.out.String()` 仍在 `:574`（原句未动）；arm1 的 `prePlant` 是**独立新读**（不是把窗起点挪早），
  arm2 复用**同一个** `mark` 再读窗尾（`:668`）⇒ AC#3 禁的"为变绿把窗起点前挪"没发生。

### AC#4 — 不许误伤真绿 —— **成立（对本腿测到的形）**
- `cur`（无 overlay，横幅不含这些字样）：`logs/AC4-headruler-cur.txt` = **PASS 2.09s**。
- `PART`（横幅只含"需要重启进程"、**不含**"本次运行不会生效"那五字）：`logs/AC4-headruler-PART.txt` = **PASS 2.12s**＝边界不误伤。
- 同族回归 `-run TestTicket223`（本腿自跑 `logs/AC4-family-223.txt`）：**10 PASS / 0 FAIL / 0 SKIP**，含被审那枚；
  终态整包里 TestTicket223 全族亦 10 枚全绿（`logs/G-end-rednames.txt`）。兄弟用例一枚没被误伤。
- ⚠ 但"不误伤"的**反面**（该响的没响）由 AC#2 注脚与 §2 进攻(ii) 触及：post-mark 欢迎语形仍漏，见 §2。

### AC#5 — 门禁四数 —— **成立（含一枚具名、已定性为环境 flake 的作差非零）**
本腿自跑（`logs/AC5-d22scan.txt`／`logs/AC5-gates.txt`）：
- `sh scripts/d22scan.sh` **rc=0**，且脚本**先跑种子违规正控**（`runtests.sh OK PASS=35/FAIL=0`，含 `TestBuiltBinaryGoesRedEndToEnd`）
  ⇒ 证明门**能红**，随后的全扫"clean"才不是假绿；被审 `cmd/wisp` 落在 ban #8 覆盖的 `cmd/ 104 Go files（含 _test.go、含注释）`，无违例。
- `gofmt -l` 只喂两枚 `.go`：**clean**；`gofumpt v0.12.0 -l`（GOPATH/bin/gofumpt.exe）：**clean**；`go vet ./cmd/wisp/`：**rc=0**。
- 终态整包名册（`logs/G-end-rednames.txt`）：RUN/PASS/FAIL/SKIP = **2174/1471/7/8**，起手 = **2174/1472/6/7**。
  **逐名作差＝ +1 红**：`TestPanelHostRealWindowHopAndLifecycle`（`cmd/wisp`，真面板宿主窗口用例）。
  本腿单跑它对 clean HEAD ⇒ **PASS 0.85s**，且本腿全程 overlay、工作树 `cmd`/`internal` 逐字节未变（起手与终态跑的是同一份码），
  ⇒ 该红是**带载/环境性 flake**（票 AC#5 已具名同类 `TestTenOpsInOneToolCallGetOneConfirm` 不许为它改断言），**非那 52 行引入**；SKIP 7→8 同属面板宿主抖动族。
  ⇒ 严格"作差＝0"这一格**未做到字面 0**，但增量一枚已定性为环境噪声、与本票产码无关；本腿**未放宽任何断言**。

### AC#6 — 还原自证 —— **成立（本腿全程 overlay，零原地编辑）**
全部突变走 `go test -overlay`，拷贝在仓外 `D:/tmp/wisp272v1/mut/<形>/`；本腿只写 `.scratch/wisp/probes/272/v1/**`。
终态 `git status --porcelain -- cmd internal` 与 §0 起手逐字（空）对照、`config_reload.go` md5 对拉见 §6。

---

## §2 ★恒真性进攻（本仓最硬的规矩：判据换成反形还全绿＝它对这件事不敏感）

本腿造**两发**定向突变攻 arm1/arm2（raw 均在 `logs/`）：

### 进攻(i) — arm1 阈值 `!=0` → `!=1`（问：`M` 的红因到底是不是 arm1？）
尺＝HEAD 测试文件（把 `:595` 那句阈值改一个字符），产码＝`M` 形。
读数 `logs/ATTK1-arm1mut-M.txt`：**`--- PASS (2.13s)`**。
⇒ 结论：**去掉 arm1 的牙，`M` 立刻回绿** ⇒ `M` 的红**唯一由 arm1 造成**（红因归属坐实，不是别的东西在响）。
副作用坐实：**arm2（票面的甲）单独对 `M` 零牙** ⇒ 印证实现件 §5.2"甲 alone 在 M 上 PASS 2.15s"，也印证票面三形里落不了 `M`。

### 进攻(ii) — "横幅里那份出现在 mark 之后"（票面本意那一形的强化版）
产码突变＝**删掉** `reportRestartPending` 的 stdout 句（`:321-326`，让本次改动不再打真重启句）**并在** `reportReload`（`:178`，
跑在 **mark 之后**、每次 `rep!=nil` 只一次）**注入一份**含全部六 needle 的"首 tick 欢迎语"式打印。语义＝
"产码哪天把那句话挪到 mark 之后再打（异步横幅／欢迎语），且真重启路径其实没触发"。
读数 `logs/ATTK2-postmark-welcome.txt`：**`--- PASS (2.13s)`**。
⇒ 结论：**arm1 与 arm2 两支都不响**（pre-plant 份数＝0 → arm1 静；窗内恰 1 份 → arm2 静）。
**票面 `:12` 点名的前向风险——"以后只要有人把那五字复制进横幅／欢迎语，现在这套尺一声都不会响"——本票**没堵死**：
arm1 只堵了**种改动之前**那份（启动横幅），post-mark 单份欢迎语这一形**实测仍假绿**。
⇒ 这比实现件 §2.2/§5.2 更进一步：实现件只在**散文里**承认这支残余臂盖不住，本腿**把它造出来并跑绿**，
证明缺口不是理论残余而是**可复现的假绿**。arm1 注释 `:590-592`（"A copy that arrives between this snapshot and the
plant is inside the window and is caught by arm 2's count"）在**单份**情形下**为假**——arm2 只在 ≥2 份时响。

---

## §3 我推翻／坐实实现件 `272-r2` 的哪几句

1. **坐实**（不是推翻）§4.2 的核心机制：`272-r2` 说"票面 `:10` 那半句『启动横幅天然落在窗口里』复现不了、横幅在 `mark` 之前"——
   本腿 AC#1/AC#2 用 overlay 独立复跑，`MDEL` 在 HEAD 尺 `:623` await 窗段为**空**（`logs/AC4-headruler-MDEL.txt`），
   与实现件 §4.1/§4.2 的读数同向。⇒ 实现件这一推翻**站得住**，予以确认。
2. **推翻（补强）实现件把落地写成"两支臂钉死缺口"的观感**：实现件 §2.2 已诚实登记 post-mark 残余"本腿不硬造"，
   但 §5.2 表格＋§5.3 落地读数的**整体语气**容易让裁决者以为缺口已闭合。本腿把那一形**造出来跑绿**（进攻(ii)），
   证明它**不是不可复现的理论边界而是活的假绿**。⇒ 落地件对 arm1/arm2 覆盖面的表述**偏乐观一档**。
3. **推翻 arm1 注释里一句绝对话**：`:590-592`"…caught by arm 2's count"——单份 post-mark 非改动副本 **arm2 不响**（进攻(ii) 实测）。
4. **小瑕疵（非承重点）**：实现件 §5.1 / 测试文件 `:40-41` 常量注释把两枚既有字面标为 `:622 / :647`，
   现量 HEAD 实为 **`:623`（awaitStdoutSince232）/ `:648`（operator Contains）**，**行号漂 1**。不影响断言，属簿记过期。
5. **红因归属**：实现件称 `M` 红在 `:595`（arm1 的 `if` 行）；本腿现量红在 **`:596`**（arm1 的 `t.Errorf` 行，`:595` 是 `if` 条件行）。
   同一枚 arm1，指向一致，仅行号口径差一行。

---

## §4 我推翻编排者（票面）的哪几句

1. **票面 `:10` 的机制叙述错**："这道窗口的起点取早了——它从'种改动之前'就开始算。**于是那行启动横幅天然落在窗口里**"。
   本腿 + 实现件双向现量：横幅由 `startConfigReload`（`run.go:813`）**同步**写进 `syncWriter`，位置**早于** `mark`（`:574`），
   `TrimPrefix` 把它算到**窗外**；`MDEL`（横幅带全字样＋真句删光）在未改尺上 await 的窗段**为空**、`M` 却 PASS ⇒ 假绿的真正成因
   **不是**"横幅落进窗"，而是**旧尺只做 `Contains`、不数份数也不查 pre-plant**。⇒ **票面把病诊断错了**，导致它给的三形（都建立在"横幅会进窗"的共同前提上）**没有一支能钉 `M`**。
2. **票面 `:27`"三形都要给同一发读数：`M` 必须红"——按票面自己的枚举无法达成**：本腿进攻(i) 证明三形里唯一沾边的甲单独落对 `M` 零牙；
   能钉 `M` 的只有枚举外的 arm1。票面 `:21` 虽写"形不在票面写死"救回了合法性，但 `:27`/AC#2"三形选一支"的**措辞与实际可达性矛盾**，
   该由编排者定案更正（本腿不擅改票面文字）。
3. **票面 `:12` 的价值主张只兑现一半**：票说这票值钱是因为"以后有人把五字复制进**横幅／欢迎语**尺都不响"。
   落地后 **横幅（pre-plant）这一半被 arm1 堵了**，**欢迎语（post-mark）这一半本腿实测仍不响**（进攻(ii)）。
   ⇒ 若票面"复制进欢迎语"是真实前向场景，**本票目标未完全达成**；要钉死它需要一个 plant 与该句之间的**同流因果锚点**，
   而产码里不存在（操作员句在 `config_reload.go:321-326`、`state=applied` 在 `:179` 更晚），造锚点＝改产品文案＝票 `:40` 禁区＋需人工批准。

---

## §5 未做完的格（逐条具名，⛔ 不留空话）

1. **AC#2 的"乙／丙两支候选尺"未由本腿各自 overlay 复跑**：本腿只独立复现了 `M`(AC#1/AC#2)、`cur`/`PART`/`REPEAT`/`MDEL`，
   并用**进攻(i)**（arm1 阈值翻转 → `M` 翻绿）间接钉死"甲/arm2 单独对 `M` 零牙"这一承重点。
   实现件 §5.2 的"乙在 M 上 PASS 2.88s（乙＝无操作）"与"丙在 cur 上 FAIL 42.20s（丙误伤真绿）"两发，本腿**未自跑复现**，
   仅采信其方向（与我读的机制一致：`mark` 已在全部启动输出之后 ⇒ 乙无操作；操作员句早于 `state=applied` ⇒ 丙钉因果序即钉死假绿）。
   ⇒ 若裁决要三形各自读数齐，这一格本腿缺两份自跑。
2. **AC#6 的 porcelain 未逐字等于 §0 起手**：终态 `git status --porcelain -- cmd internal` 多出 `?? internal/winsec/wisp129-dr-15652/`——
   系本腿两次整包 `go test ./internal/...` 期间 winsec 用例留下的**测试产物目录**（非本腿 authored 代码、非跟踪文件改动）。
   按票面"临时件只建不删"（`issues/README` 规则 8），本腿**不删它**、如实登记。`cmd/` 侧 porcelain 为空、跟踪文件零改动（见 §6）。
3. **形定义复用 `272-r2` 的仓外拷贝**（`mut/M`／`MDEL`／`PART`／`REPEAT` 的 `config_reload.go`）：本腿**读数全自跑**，
   但 `M`/`MDEL`/`PART`/`REPEAT` 四形的**产码突变内容**用的是实现件留在 `/d/tmp/wisp272r2/` 的那份形定义（本腿已 `diff` 对拉确认＝票面语义），
   非本腿从票面散文另起炉灶重建。进攻(i)/(ii) 两发的突变文件则是本腿自建。

## §6 终态自证（本腿跑完，一切突变只在 overlay／仓外）

- `git status --porcelain -- cmd internal`：`cmd/` 侧 **0 行**；`internal/` 侧仅 §5-2 那枚测试产物目录。
  ⇒ **跟踪的产码/测试文件一字未回退、未脏改**（GBK 中途炸留脏的情形未发生：本腿从未原地编辑工作树，全程 `-overlay`）。
- `git show HEAD:cmd/wisp/config_reload.go | md5sum` = **`5ce441ca5e72b64d18a6c26f1c066882`** ＝ §0 起手基线 ＝ 工作树文件 md5。产码一字未动。
- `git show HEAD:cmd/wisp/config_reload_223_test.go`（md5 **`2ea71578ba9ca015c535673bc9792d09`**）
  与 `8d30a862` 版逐字节相同（md5 一致）⇒ 本腿审的版 == 落 HEAD 的版；被审的 52 行内容（arm1 `:594-598`、arm2 `:668-671`、
  `restartSentence272` `:43`）与本件 §1/§3 的逐条描述**对得上**（arm1 阈值句现于 `:595`、其 `t.Errorf` 现于 `:596`；
  实现件正文写 `:595`/`:622`/`:647` 系行号口径，本腿现量红点为 `:596`、两枚既有字面为 `:623`/`:648`，见 §3-4）。
- 本腿只写 `.scratch/wisp/probes/272/v1/**`；两笔 commit 均带显式 pathspec（`1ec9fa19` §0、本笔 §1–§6）。
- ⛔ 未 push、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`、未改 SLO/golden/thresholds/PLAN/specs、未给 CI 加 `-tags winlive`、未放宽任何断言。

---

## 一句话总结（给编排者）
六格中 **AC#1/AC#3/AC#4/AC#6 成立、AC#2 硬要求成立但"三形选一支"字面不成立（钉 M 的是枚举外新臂 arm1，
按票 `:21` 合法、不退回）**、**AC#5 成立（静态四数全清，终态作差 +1 枚已定性为面板宿主环境 flake）**。
**最值钱的一条**：本腿进攻(ii) 造出"post-mark 欢迎语单份"那一形，在落地尺上**仍 PASS** ⇒ 票 `:12` 自陈的"复制进欢迎语尺都不响"
这一前向风险**没堵死**（arm1 只堵 pre-plant、arm2 只堵 ≥2 份）；实现件把它当"本腿不硬造"的散文残余，本腿把它**跑成了活的假绿**。
