# 票 272 r2 — 裁决件（写码腿 `272-r2`）

工单：`.scratch/wisp/issues/272-restart-tier-test-window-starts-before-the-plant-so-a-banner-carrying-the-same-words-passes.md`
被审尺：`cmd/wisp/config_reload_223_test.go` 的 `TestTicket223RestartTierSaysItWillNotApply`（`:554` 起）
被审产码（**一字不许动**）：`cmd/wisp/config_reload.go`

生成时刻 `2026-10-07 09:1x +08`。本文件按 AC 逐格长；每一格落地后另起一笔 commit。

---

## 0. 起手锚（在任何改动之前取）

命令（逐字）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v ./cmd/wisp ./internal/...
```

- 起手时刻：`2026-10-07 09:04 +08`，结束 `09:11:53 +08`，`rc=1`（起手就带红，不是本腿造的）。
- 起手时 `git rev-parse HEAD` = `4dab3fbc30b65c29cc1d3a95588315e2734891f0`；
  `git status --porcelain -- cmd internal` = **0 行**（产码与测试文件都是干净的 HEAD 版）。
- 被审两枚文件的起手 blob／md5（终态要逐串对拉）：
  - `cmd/wisp/config_reload.go` blob `d075443c75a6d14d3b737961a2087c70286957e2`／md5 `5ce441ca5e72b64d18a6c26f1c066882`
  - `cmd/wisp/config_reload_223_test.go` blob `ab6c84fb79b09d5403cc519671ce283a326e6942`／md5 `d0f6433c0ee722ecf5f22b98456e7d3b`
- 原始 `-v` 输出 8,915 行／1.1 MB，**不入库**，留在仓外 `/d/tmp/wisp272r2/mut/G0-raw.txt`（临时件只建不删）。
- 入库的摘要件：`logs/G0-start-rednames.txt`（14 行逐名）／`logs/G0-start-failingpackages.txt`（4 行）。

**四数（顶层 `^=== RUN` / `^--- PASS` / `^--- FAIL` / `^--- SKIP`）＝ 2175 / 1471 / 7 / 7。**

起手红名册（7 枚，终态作差必须＝0）：

1. `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` (50.12s) — `cmd/wisp`，带载型（整包并发跑 50s 撞 40s 预算）
2. `TestC21TableColourRowsMatchTokensCSS` (0.00s) — `internal/ball`
3. `TestApprovalCardViewJSONKeysMatchFrontendTypes` (0.01s) — `internal/ball`
4. `TestComposerContractTypesMatchFrontend` (0.01s) — `internal/ball`
5. `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` (0.50s) — `internal/ball`
6. `TestC21DesignTokensFourWayAgree` (0.00s) — `internal/ball`
7. `TestResolvePerCallBudget` (3.70s) — `internal/risk`

起手 SKIP 名册（7 枚）：`TestDefaultDeadlineWallClockMeasurement`／`TestLiveWasapiSmoke`／`TestSubprocessCrashWriter`／
`TestRealDownloadVadThroughPipeline`／`TestRealDownloadPuncArchiveThroughPipeline`／`TestHelperProcess`／`TestSyncRegistryProbeLive`。

`TestTenOpsInOneToolCallGetOneConfirm` 起手这发**绿**（票面 AC#5 具名要求不许为它改断言；本腿不碰）。

---

## 1. 票面六格原文（逐字抄自工单 `:31-36`，AC 框由编排者翻，本腿一枚不动）

- [ ] **AC#1 先复现"今天拦不住"**：在未改动的那枚尺上跑 `M` 形 ⇒ **必须仍 PASS**（＝未修码读数，本票立案的凭据）。抄原始读数，⛔ 不许跳过这一格直接改断言。
- [ ] **AC#2 改完之后同一形必须红**，且红因是**具名那枚用例**、不是超时也不是编译不过。三形（甲／乙／丙）选一支落地，另两支撑不住的理由写进证据件。
- [ ] **AC#3 零放宽**：⛔ 不许删任何既有断言，⛔ 不许把六枚 needle 从 stdout 挪回拼接串，⛔ 不许为了变绿把窗口起点再往前挪（那正是本票要修的形）。
- [ ] **AC#4 不许误伤真绿**：正常那一发（横幅里**没有**这些字样）改完之后仍 PASS；并给一发"横幅只含**部分**字样"的边界读数。
- [ ] **AC#5 门禁四数**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" sh scripts/d22scan.sh` rc=0 且正控先绿；`gofmt -l`／gofumpt v0.12.0 `-l` 只喂 `.go`（⛔ 别把 `.md` 喂给 gofmt，那出 `U+0023` 且 rc=2 是你的工具误用）；`go vet ./cmd/wisp/` rc=0；终态 `go test ./cmd/wisp ./internal/... -count=1` 逐名红名册与起手**作差＝0**（★新增：`TestTenOpsInOneToolCallGetOneConfirm` 属带载型，安静单跑 0.07s 绿 ⇒ 已具名入在册，见 `A647`，⛔ 不许为它改任何断言）。
- [ ] **AC#6 还原自证**：所有突变只走 `go test -overlay`，⛔ 绝不原地编辑共享工作树；终态 `git status --porcelain -- cmd internal`＝起手、被审文件 md5 与 `git show HEAD:` 对拉逐串相同。

票面禁区（工单 `:40` 逐字）：⛔ **不改产品文案**（`cmd/wisp/config_reload.go` 一字不动——那是票 232 已判过的禁区，也是 AC#2 那格"具名让步"的根）；⛔ 不动 golden／`thresholds.go`／`PLAN.md`／`docs/specs/**`；⛔ 给 `ci.yml` 加 `-tags winlive` 一律不许；⛔ 本票只动 `cmd/wisp/*_test.go`，⛔ 不许顺手改产码来"让横幅不冲突"（那是另一件事，且要人工批准）。

票面三形定义（工单 `:23-27` 逐字要点）：

- **甲**＝窗口内计数：`strings.Count(窗口, 那句话) == 1`（多了就是横幅也算进去，少了就是那句话没到）。
- **乙**＝把标记取在 `state=armed` 那行横幅**之后**（起点搬晚，横幅天然落在窗外）。
- **丙**＝要求那句话的**偏移晚于** `HOT-RELOAD state=applied` 那一行（钉"因果顺序"而不是钉存在）。
- ⛔ 三形都要给同一发读数：造出 `M` 形（横幅先带同样字样、之后种改动、重启句照打）⇒ **必须红**。

---

## 2. 本腿判不动／要报回冲突的格子（骨架；读数落地后逐条填实）

### 2.1 ★与票面转述的第一处冲突：起手锚测到的机制与票面 `:10` 的机制不一样（现量在 §4）

票面 `:10` 写："这道窗口的起点取早了——它从'种改动之前'就开始算。**于是那行启动横幅天然落在窗口里**"。
本腿读码位置（不引结论，先给锚）：

- `cmd/wisp/config_reload.go:124-127` 的启动横幅 `fmt.Fprintf(rt.stdout, "wisp run: 配置热加载已接管…")` 在
  `startConfigReload()` 里，而 `startConfigReload()` 的**唯一产码调用点是 `cmd/wisp/run.go:813`**（`assembleRuntime` 内）；
  测试的 `mark := r.h.out.String()`（尺 `:566`）跑在 `run.go:259-260` 的 `onRuntime(rt)` → `r.live` 的 `rtHook` 里，
  ⇒ **横幅在程序顺序上先于 mark 写完**。
- `rt.stdout` 的实体是 `cmd/wisp/approval_reply_201_test.go:64` 的 `syncWriter`（`mu + bytes.Buffer`，`Write` 同步落 buffer，
  **中间没有 io.Pipe／没有拷贝协程**），所以"写完"＝"已经能被 `String()` 读到"。
  ⇒ 只要种改动之前没有别的 stdout 写者，**启动横幅落在 `mark` 里、天然在窗外**。

这条推论如果为真，直接决定 AC#2：**票面给的三形（甲／乙／丙）没有一支能把 `M` 形弄红**——因为 `M` 与"正常那一发"
的差别只在窗外字节里，窗内一模一样。本腿按"原文（＝盘上实测）为准"复跑，读数见 §4；对不上的那一面**具名报回编排者**，
不在本件里替票面说话。

### 2.2 判不动清单（当前）

- **AC#2 的"三形选一支"这一格**：若 §2.1 的机制被实测钉死，则"甲／乙／丙"三形的**共同前提**（横幅会进窗）不成立，
  三形都落不到"把 `M` 变红"。本腿能落地的是与票面同目的的一支（在 `mark` 上钉"种改动之前零份"+ 窗内钉"恰好一份"），
  是否算"三形选一支"要由裁决者定，**不由本腿自封**。
- **AC#5 终态作差＝0 这一格**：起手 7 枚红里 `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 是带载型
  （整包并发 50.12s 撞 40s 预算）。终态若它复绿或别枚带载的复红，**作差就不是 0**，本腿不会为它改任何断言，
  只把两法的红名册与差集逐名列出来交给裁决者。
- **两支臂都盖不住的那一发（本腿不硬造，具名登记为残余）**：产码若把那句挪到 **`mark` 之后、plant 之前**才打
  （异步横幅／首个 tick 打的欢迎语），那么"窗内恰好一份"与"种改动之前零份"两条**同时成立**，尺照样分不出抵账的那一份是谁打的。
  要钉死它需要一个 **plant 与那句之间的同流锚点**，而产码里不存在（操作员句在 `cmd/wisp/config_reload.go:321-326` 写，
  它是 `internal/config/manager.go:201-202` 在 `CheckAndReload` **返回之前**回调的；`HOT-RELOAD state=applied` 那行要等
  `config_reload.go:179` 的 `reportReload` 才写＝**晚于**操作员句），
  造这个锚点要改产品文案＝票面 `:40` 禁区＋要人工批准。本腿**不改产码一字**，把这一形写在这里交给裁决者。
- 其余 AC#1/#3/#4/#6 都在本腿射程内，不需要人拍板。

---

## 3. 计划（写给自己，不是承诺）

1. AC#1：未改动的尺上跑 `cur`／`M`／`MDEL`（＝票面 `MK` 那一形：横幅带同样字样 **且** 操作员重启句删掉）。
   `M` 必须 PASS（票面字面要求）；`MDEL` 才是"横幅能不能抵账"的判据——它红＝抵不了，它绿＝票面机制成立。
2. AC#2：落一支尺，`M` 必须红且红因是具名用例的**断言句**（不是 40s 超时、不是 build failed）。
3. AC#3/#4：`diff` 自证零删除；`cur` 仍绿；`PART`（横幅只含部分字样）给一发边界读数。
4. AC#5/#6：门禁四数＋还原自证。

---

## 4. 读数台账（ newest last；每格落地后追加，原话不抹 ）

### 4.1 AC#1 —— 未改动的尺（`config_reload_223_test.go` blob `ab6c84fb`）上的五发读数

命令（逐字）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 -v -run 'TestTicket223RestartTierSaysItWillNotApply' [-overlay <json>] ./cmd/wisp/`
五份原始读数：`logs/AC1-oldruler-{cur,M,MDEL,PART,REPEAT}.txt`（19/19/34/19/19 行，共 11.5 KB，全部本腿自跑，**没抄 r1 的**）。
产码形全部走 `go test -overlay`，拷贝在仓外 `/d/tmp/wisp272r2/mut/<形>/cmd/wisp/config_reload.go`，工作树未动（`git status --porcelain -- cmd internal` 仍 0 行）。

| 形 | 改动行数（对 `git show HEAD:cmd/wisp/config_reload.go` 的 diff） | 语义 | 未改动的尺读数 |
|---|---|---|---|
| `cur` | 0（无 overlay） | 基线 | **PASS 2.77s**（件内 `--- PASS` 行 2.38s／包 2.771s） |
| `M` | **1 行插入**（`127a128`，与 r1 那份 M 逐字节同文） | 启动横幅追加一句同样字样的操作员句，重启句照打 | **PASS 2.38s＝票面 AC#1 要的那发，本腿复到了** |
| `MDEL` | 1 插 + 6 删（＝票面 `MK` 那一形：横幅带全字样 **且** 操作员重启句删光） | 横幅能不能单独抵账 | **FAIL 42.08s**，红在 `:592` `stdout never carried "本次运行不会生效" AFTER the plant within 40s`，红句里 **`window since the mark:` 那一段是空的**，`上一版横幅` 只出现在 `full stdout:` 第 3 行 |
| `PART` | 1 行插入 | 横幅只含部分字样（有「重启进程」，**没有**「本次运行不会生效」） | **PASS 2.48s** |
| `REPEAT` | 6 行插入（`:321-326` 操作员句 Fprintf 原样重复一遍，仍是种改动之后打的） | 窗内出现两份真句 | **PASS 2.22s＝今天窗内计数不钉，重复打印不响** |

### 4.2 ★AC#1 的复验结论：票面 `:10` 那半句机制**复现不了**，但缺口本身换了个形状还在

- **复现不了的那半**：`M` 形 PASS 不是"横幅冒充了重启句"。`MDEL` 这一发（横幅带全字样＋真句删光）在未改动的尺上**红在 40s 超时**，
  且红句里的窗口段是**空的** ⇒ 启动横幅那一份**根本没进窗**。机制与本腿 §2.1 的读码一致：
  `startConfigReload()`（`run.go:813`，`assembleRuntime` 内）同步写 `syncWriter`（`approval_reply_201_test.go:64`，无 pipe 无拷贝协程），
  而 `mark` 在其后的 `onRuntime`→`rtHook`（尺 `:566`）里取 ⇒ **种改动之前打的一切都在 `mark` 前缀里，`TrimPrefix` 天然把它们算到窗外**。
  ⇒ 票面"窗口起点取早了所以横幅天然落在窗口里"这半句，按本腿现量应读作：**窗口的下界确实早于 plant，但那个位置今天恰好也晚于全部启动输出**。
- **还在的那半（＝本票真正的价值所在）**：`REPEAT` 形今天**一声不响**（PASS 2.22s）——窗内两份真句与一份真句在这把尺上不可区分；
  而"横幅能不能抵账"只钉在一条**程序顺序**上：产码哪天把那句改到 `mark` 之后再打（异步打印／首 tick 打印／欢迎语晚一步），
  `MDEL` 那种抵账立刻就能绿。⇒ 本腿落地的两支臂（见 §5）钉的是"种改动之前零份"＋"窗内恰好一份"，
  **不**是票面三形里任何一支的原样，理由与量到的证据在 §5.2。

