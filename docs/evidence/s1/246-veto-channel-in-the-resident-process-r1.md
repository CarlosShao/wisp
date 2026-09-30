# 票 246 — 落地腿 `246-r1` 交付表：最短链「注入 gate → 真卡片 → `Esc` 否决一次并归还 → 退出序列不撒谎」

> ⚠ **这份表的身份**：出自**实现方本人**（`246-r1` 是产码腿）。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，
> `docs/evidence/s1/` 的裁决表**必须出自非实现者**才算判语；本文件全部格子的性质是
> **〔实现方自述＋台件〕**，等另一程对抗验收裁，**本表不构成任何一格 AC 的翻勾依据**（AC 框归编排者，本腿一枚未动）。
>
> 判据出处＝`.scratch/wisp/issues/246-…-no-veto-channel-has-an-executor.md`「判据（本票现行射程）」AC#0..AC#6。
> 前段代价表＝`.scratch/wisp/probes/246/a1/census.md`（206 行）；编排者裁定＝乙形（装配根注入）＋"钩子注册入口不算契约面"，账 `A481`。

---

## §0 主张与起手名册

### §0.1 这一程主张的产品事实（只这一条，别的都不主张）

**双击图标起来的那个常驻进程，今天有了一枚真审批门；一张真待决卡片能在那里挂起来；`Esc` 在 `Confirming`
期间真能否决它一次并归还键；卡片挂着时收到退出信号，D38(e) 第 3 步不再被记成 skipped，卡片以拒绝收口并留下审计。**

⛔ 本表**不主张**：用户看得见卡片（面板今天开不出来，`cmd/wisp/approval_always.go:165` 逐字 this binary links no WebView2 host）、
`Esc` 安全了、四条否决通道通了（只落 `Esc` 一条，AC#6）、常驻进程能跑任务（没有 agent loop、没有 provider）。

### §0.2 起手名册与锚点（闸门＝终态等于起手那一刻的名册，不是"必须为空"）

| 项 | 读数 |
|---|---|
| 取数时刻 | `2026-09-30 18:42:54 +0800`（`date` 现跑） |
| 起手锚点 | `ce1ade1f` — `ledger(A481 收 246-a1：裁走乙＝装配根注入，票 246 判据落成最短链六格…)`（`git log -1` 自取，⛔ 未采用派单里给的任何 sha） |
| 分支 | `dev` |
| 起手名册全量 | `.scratch/wisp/probes/246/r1/roster-start.txt`（`git status --porcelain` 原文，取数即存） |
| `git status --porcelain` 行数 | **178**（含本腿已建的 `.scratch/wisp/probes/246/r1/` 未跟踪件） |
| 起手脏件里"别的写腿／别人的账" | 全部保留原样：`design/**` 删除与改动、`docs/evidence/s1/152-…-accept-r1.md`、`.gitignore`、`.scratch/wisp/probes/**` 各票台件、`part*-*.txt` 等。本腿**一枚不碰、不删、不提交**。 |
| 本腿写面（白名单） | `cmd/wisp/**`（常驻腿＋装配根）＋ `internal/proc/**`（钩子注册入口）＋ 新增测试台件（`internal/ball/**` 或 `cmd/wisp/**`，AC#3）＋ `docs/evidence/s1/246-*.md` ＋ `.scratch/wisp/probes/246/r1/**` ＋ 票 246 面「Progress log」末尾一行 |
| 起手产码面是否干净 | `git status --porcelain internal cmd` 起手现跑＝见 §0.3（⛔ 与票 228 后续片／票 245 AC#6..#9 串行，本腿独占 `cmd/wisp`＋`internal/proc`） |
| 终态闸门算法 | `diff <(git status --porcelain) .scratch/wisp/probes/246/r1/roster-start.txt` 的差集只允许落在本腿白名单写面内；`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／三枚冻结件必须逐字节等于起手 |

### §0.3 起手现跑（写第一枚产码之前）

| 尺 | 命令 | 读数 |
|---|---|---|
| 起手产码面脏件 | `git status --porcelain -- internal cmd` | **0 枚**（现跑，`18:42` 同发） |
| 起手依赖图（AC#1 的"前"） | `GOOS=windows go list -deps ./cmd/wisp` | **276 行**，落盘 `.scratch/wisp/probes/246/r1/deps-cmdwisp-before.txt` |
| 起手逐包 import（"前"） | `GOOS=windows go list -f '{{join .Imports "\n"}}' ./cmd/wisp ./internal/proc ./internal/ball` | **87 行**，落盘 `imports-before.txt` |
| 桌面安静 | `tasklist //FI "IMAGENAME eq balldebug.exe"`／`wisp.exe`／`go.exe` | 三者起手均为 **0 枚**（无同类进程占全局热键） |
| **AC#5 起手红归因** | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1` | **`ok github.com/CarlosShao/wisp/cmd/wisp 143.600s`** — ⛔ **没有红**（`--- FAIL` 0 枚；rc=0；⚠ 不是"根本没跑"那形：那一形会是 `exit status 0xc0000135`＋`0.0xxs`，本次是 143.600s 真跑完）。详见 §1 AC#5 |

---

## §1 判语（逐格）

> 每格姿势：判据原文 → 落地形状 → **现跑读数**（文件:行 ＋ 命令输出逐字）→ 自判（成立／不成立／阻塞）＋ 为什么。
> ⛔ 全部是〔实现方自述〕，裁决等另一程。

### AC#1 装配根注入一枚真审批门到常驻腿，零新增包级依赖边

- 判据：由 `cmd/wisp`（装配根）把 gate 递进常驻那条腿；⛔ 不开任何新包级依赖边；尺＝`go list -deps`／逐包 import 清单**落地前后差集逐名相同**。
- 落地形状：（取数中）
- 差集读数：（本轮尚未落地，落地后同尺复跑并逐名对撞）

### AC#2 一张真卡片能在常驻进程里挂起来（⛔ 不许 mock 代真）

- 判据：走 `internal/agent/approval` 既有机制造**真**待决项，状态机真进 `Confirming`；可见证据只到"日志＋状态位＋Win32 那一层"。
- 落地形状：（取数中）

### AC#3 `Confirming` 期间 `Esc` 真能否决一次、完事归还

- 判据：三形读数都要（借到／否决真生效／归还后另一进程真收到 `Esc`）；⚠ 票 245 AC#9 那套"第二个进程观察 `Esc`"的台件**必须先搬进仓**。
- 落地形状：（取数中）

### AC#4 卡片挂着时收到退出信号不许撒谎（D38(e) 第 3 步不再是 skipped）

- 判据：`internal/proc/boot_windows.go` 的 `Shutdown` 今天逐字 `hooks := ShutdownHooks{}` 只填 `CloseJob` ⇒ 第 1–7 步零生产者；把"取消任务根"做成**真注册的钩子**；一发真机：卡片挂起 → 发退出 → 那一步 `StepRecord` 不再是 skipped、待决卡片被判拒绝＋留下审计。⛔ 十步顺序一字不许动。
- 落地形状：（取数中）

### AC#5 起手第一发先归因那枚红

- 判据：`cmd/wisp` 默认层今天有一枚未归因的红（用例名丢了，账 `A480` ④）⇒ 开工前跑默认层、取红名、判"起手即在 vs 本票造成"；⛔ 不许当已知常红略过、⛔ 不许顺手修不属于本票的东西。
- **现跑读数（起手，动任何产码之前）**：

  ```
  PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1
  ok  	github.com/CarlosShao/wisp/cmd/wisp	143.600s
  ```

  取数时刻 `2026-09-30 18:42:54` 起、约 `18:45` 落；锚点 `ce1ade1f`；`--- FAIL` 计数 **0**；`FAIL` 计数 **0**。落盘 `.scratch/wisp/probes/246/r1/ac5-start-default-tier.txt`。
- **归因结论**：**起手这一发没有红可归**。台账 `A480` ④ 那一枚（`245-v1` 自报"1 红 2 绿、用例名丢了"）在本腿起手的锚点、桌面无同类进程、无并发 `go test` 的条件下**未复现**；`245-v1` 自己具名写过它那一发红的环境差（"第 1 发期间我在 17:56 编译过一次台件、CPU 有争用"，表 `245-…-v1.md:184`）。⇒ 本腿判：**既不是"起手即在的稳定红"（这一发不支持），也不是"本票造成"（这一发之前于任何产码改动）**；剩余可能性（时序抖动型 flake）本腿不替它定案，只把复跑样本交回（§3 复跑行）。
- ⛔ 本腿**没有**因此去修任何东西（那枚红不属于本票射程）。

### AC#6 其余三条否决通道不许顺手做

- 判据：单击球 veto／KWS 语音否决词／面板拒绝本票一律不实现；完成判据＝票面「Progress log」具名写"只落 `Esc` 一条，另三条各有归口"。
- 形状：（落地后复述，具名三枚通道的 loaded 读数见 §1 末「通道名册」）

---

## §2 突变与正控名册

> 姿势：每发先 `grep -c` 锚点＝1，落盘，跑指名的用例（判红绿只认 `--- FAIL`），跑完
> `git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet -- <写面>` 复认。⛔ 不用 `-overlay`。

| 编号 | 要攻的那一格 | 突变内容 | 期望红 | 实跑读数 | 还原复认 |
|---|---|---|---|---|---|
| M1 | （待取数，见 §2.1） | | | | |

### §2.1 名册说明

（本腿的突变名册与正控在 §1 各格形状落地后逐枚执行；⛔ 未执行的不得写成已执行。）

---

## §3 门禁读数

> 每行＝命令、时刻、rc、逐字摘要。缺档要具名写"该档未读"，不填空。

| 尺 | 命令 | 时刻 | rc | 读数 |
|---|---|---|---|---|
| 默认层（起手，产码前） | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1` | 18:42:54 起 | 0 | `ok … 143.600s`；`--- FAIL` 0 枚 |
| 默认层（起手复跑第 2 发，`-v`） | `… go test ./cmd/wisp -count=1 -v` | 18:5x 起 | （进行中） | 落盘 `ac5-repeat2-v.txt`（含 `=== RUN` 名册，供收尾对撞） |

（本节后续按取数逐行补：`go build ./...`／`go vet ./...`／`go vet -tags winlive ./...`／`gofmt -l`／`scripts/d22scan.sh`／指到包的真机 winlive 读数／依赖边差集复跑。）

---

## §4 残余与归口（本腿做了的、没做的、该归谁的）

| # | 事项 | 归口 |
|---|---|---|
| R-1 | D38(e) 第 1／2／4／5／6／7 步今天仍零生产者（本票只吃第 3 步） | `246-a1` 的 D-1，编排者已分：step 4 语音链、step 5 票 15/28、step 6 票 33/35、step 1/2 scheduler/hotkey 全量＝别的落地面 |
| R-2 | 常驻腿整棵 `config.toml` 未接、托盘「退出」无执行者、球位置不持久、`ball.New` 失败路径不清理 | 票 228 后续片（`246-a1` D-4／D-5，账 `A481`） |
| R-3 | 过期注释两处（`cmd/wisp/main.go:8`／`:24` 那族句子） | 票 228 AC#8／票 245 AC#6①（账 `A477`／`A478`） |
| R-4 | 面板拒绝 `ChannelPanel` 执行者（面板侧要什么） | ⛔ 本编队零读零转述 `frontend/**`／`design/**`；只写进票面与台账，由 owner 自己带给他的前端 agent |
| R-5 | KWS 否决词 `ChannelKWS` | 票 41（`approval.go:104` 的 DEFERRED 行） |

---

## §5 判不动的地方（本腿不上报就没人裁的）

1. **"常驻进程是不是任务管线的永久之家"**：`246-a1` §6 具名判不动（D2 正文 `PLAN.md:73` 未逐字核）。本票 AC 只要求"注入 gate ＋ 一张真卡片 ＋ Esc 否决一次"，本腿按**最短链**做，⛔ 不擅自把 agent loop／provider／config 搬进常驻腿；如果产品其实想让 veto **跨进程** reach 到 `wisp run`，本腿的形状要重议——**这条本腿也判不动**。
2. **"卡片由谁在什么产品动作上举起"**：本票射程内只有 host-initiated 一发（真 gate、真队列项／真 L1 窗口）。真正"任务跑出来的一张卡"要有任务源（票 228 config ＋ 语音链）。本腿**不主张**"用户今天会看到卡"。
3. `proc.Shutdown` 钩子注册入口的**契约归属**：编排者已在 `A481` 裁"不算契约面"（凭据 `boot_windows.go:148-150` 注释），本腿照裁执行；若裁决腿认为这枚裁定越界，退回它即可，本腿不自判。
4. （后续取数中若再冒出判不动的，逐条追加，⛔ 不空写"暂无"。）

---

## §6 收尾名册对撞

（终态 `git status --porcelain` 与 `roster-start.txt` 的对撞、逐枚 commit 号、写面是否清空、临时件落盘清单——收尾现跑后逐字填，⛔ 不许在跑之前预填读数。）
