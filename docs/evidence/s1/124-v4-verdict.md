# 票 124 终裁腿 124-v4 — 三批收口裁决（裁决者 NOT 实现者）

裁决者 NOT 实现者：一格未翻、一字未改 `.go`、未 push。被裁对象＝三批七笔（part2 sec 6）。
写面只有本件与 `-part2`（仓外台面 `D:\tmp\wisp\20261021\124v4\**`）；临时件只建不删。

> 本件因 160 行帽拆两份：本件 = sec 0..3；`124-v4-verdict-part2.md` = sec 4..7（门禁逐字、红名册逐名、提交与正控、未验证项）。

## 0. 台面与锚点（本件一切读数的出处）

- 开工首读 HEAD `e97cfd9c`，收尾复读 `ef4c9702`；`git diff --name-only e97cfd9c ef4c9702 -- internal cmd go.mod go.sum scripts`
  **零 hunk**（其间只多 `docs/reports/pending-and-issues.md` 30 行）=> 被裁码与我这台秤量的是同一版树。
- 纯净件六台全部仓外 `git archive <sha> | tar -x`：`pre-r1 4ac4549c` / `post-r1 d32b459c` / `pre-r2 e1de7bab` / `post-r2 99f1f9c7` /
  `pre-r3 e4d363f2` / `post-r3 6b712571`。容器 `golang:1.27` / `go1.27.1 linux/amd64` / `GOPROXY=off` / `CGO_ENABLED=1` /
  每形一枚新容器 / `--rm`；每发先打挂载自证 `go.mod 1026` + `md5 0c266f4d…` + `resolve.go b6876a5e…`（与前各批逐字同字）。
- 软链形 `ln -s /realpriv /varlink` + `TMPDIR=/varlink/w124tmp` + `exit 90/97/98`；普通形 `[ ! -e /varlink ] || exit 99` + `/plainroot`。
  全台容器 **26 发 `-v`** ＋宿主 **9 发**，四把形状断言**零命中**，`panic` 与 `test timed out` 各 **0**；
  判据只用**名册集合**（逐名 `comm`），NOT 枚数相减。

## 1. 八项判语与读数（每项另起尺）

### (1) 票面前提新鲜度（硬闸）— 三处腐坏**逐条坐实**

- **饱和尺**：对 HEAD blob（`git grep ... HEAD`，NOT 工作树）量 `winsec.` 与 `TempDir` **同行** => **20 行**；
  其中 18 行是 `winsec.SealableTempDirForTest124(t)` 调用、2 行是注释；`winsec.` 同行 `t.TempDir|os.TempDir` => **0**；
  两行窗口同形 => 18。=> 与 `124-v1` 14:3x 那次**逐数相同**：这把尺已饱和，100% 是"已转换后"的形态，**恒 20**，今天测不出剩余量。
- **AC#6 与三批零交集**：三批 28 枚改动文件里 `internal/winsec/**` **0 枚**、POSIX 标签件 `*_other_test.go` **0 枚**、
  文本含 `absoluteness` **0 枚**（逐枚 `git show HEAD:<file>` 现读）。=> 票面 AC#6 与三批写面**真零交集**。
- **Progress log 缺早期批记录**：全票 `- [20...` 记录**只 3 条**，全是编排者（`124-v1` 收件、摇摆登账、`124-r2` 收件）；
  `d32b459c`／`124-r1`／`6b712571` 在票面命中各 **0**，`124-r3` 仅 1 命中（在 AC#5 翻框注里）。
  => 批 1／批 3 的落地记录**从未写进票面**，与票 171 的 `171-r3` 那族腐坏同形。

### (2) 三批出口是否真按"名册集合" — **三批成立，我全部自己复算**

我自己两形各一枚新容器、`-count=1 -v`，再逐名 `comm`：

| 批 | 包 | only-in-link PRE | only-in-link POST | went-green | newly-red | only-in-plain POST |
|---|---|---|---|---|---|---|
| r1 | `internal/config` | **34** | **0** | **34** | **0** | 0 |
| r2 | `internal/tools` | **38** | **3** | **35** | **0** | 0 |
| r3 | `cmd/wisp` | **37** | **0** | **37** | **0** | 0 |

- `internal/tools` 那 **3** 名逐字＝票 123 那族，NOT 其它：`TestL1WriteGoesThroughTheRealBlockWindow`（300.08s）/
  `TestLateVetoRendersTheApprovalLayersAppliedStepsReport`（300.04s）/`TestFSReadOnlyNeverOpensACard`（300.13s）。=> `38 => 3` **成立**。
- `cmd/wisp` `37 => 0` **成立**：`POST-L` 与 `POST-P` 顶层红名册**双向 0 行**（各 68 名，逐名见 part2 sec 5），
  且 37 枚 went-green 在 `POST-L` 的 `PASS` 名册里逐名在场、`SKIP` 名册 0 枚 => 不是"被跳过冒充变绿"。
- 分母与 `124-v1` 对账：tools **38**、`cmd/wisp` **37** 均与它那本**逐名 set-equal**（双向 `comm` 各 0 行）。
- **一处读数更正（我自己的，NOT 谁的自报）**：`internal/config` 我量到 **34**、`124-v1` 记 **33**；差 1 名已定位＝
  `TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent`（已知 POSIX-only 摇摆，我这一发 `PRE-L` 红、`PRE-P` 绿，于是落进 only-in-link）；
  扣掉两枚摇摆名后 **33**，与 `124-v1` 逐数对上 => 差的是摇摆这一次落在哪一侧，不是批 1 多做或少做了一枚。

### (3) AC#3（不许放宽换绿）— **三批成立；决定性一条我自己验**

- 删除侧逐笔 `git show --unified=0` 现读，**全文**只有"根路径怎么来"那一个表达式：
  `d32b459c` **9 行**、`99f1f9c7` **28 行**（`t.TempDir` 族 11 + `tempCanonical(t)` 族 12 + `ArtifactsDir: …` 3 + `openStore(t, t.TempDir())` 1 + 一行两处 1）、
  `6b712571` **15 行**；逐笔 hunk 头＝**9 / 28 / 15**，与删除行数相等（每 hunk 恰删 1 行）。
- 断言扫：对全部 `+/-` 行扫 `Fatal|Errorf|want|expect|assert|Skip|threshold|golden|t.Run|if |return |Error(`——
  **命中 1**，且是批 1 新增注释里的英文词 "assertion"（`manager_test.go` 那 4 行注释）。=> 断言/阈值/golden/判定分支 **0 hunk**。
- **决定性一条**（那些 went-green 在**改前普通形**是否已逐名 PASS）——两支尺都过：
  (a) 整包 `PRE-P` 的 PASS 名册里逐名查：config **34/34**、tools **38/38**、wisp **37/37** 全在册、**0 枚**在 SKIP 名册；
  (b) 我再打一发**只跑这批名字**的隔离台（`-count=3`、普通形、纯净 `pre` 快照）：tools `RUN=132 PASS=114 FAIL=0`、
  wisp `RUN=147 PASS=111 FAIL=0`、config `RUN=132 PASS=99 FAIL=3`（那 3 枚全是同一个摇摆名 `TestNothingWritten...`）。
  => 除那枚已知摇摆外**断言在改前本就成立**，摘掉的只是 harness 形状造成的那次拒；**放宽换绿在这个形状下无从发生**。

### (4) 零新增走法 — **成立，四枚 helper blob 逐批未动**

| helper 文件 | HEAD blob | sha256(前 16) | 该批 `parent:` vs `commit:` |
|---|---|---|---|
| `internal/config/tempdir_resolved_124_test.go:31-34` | `fca65f7df9` | `c07180f0cfc5f92b` | `d32b459c` 两侧同 blob |
| `internal/tools/tempdir_resolved_124_test.go:26-29` | `818687e43e` | `3d949a47ea844aca` | `99f1f9c7` 两侧同 blob |
| `internal/tools/tempdir_resolved_124x_test.go:19-22` | `edf028ede3` | `9c4b7381d47a629d` | `99f1f9c7` 两侧同 blob |
| `cmd/wisp/tempdir_resolved_124_test.go:35-37` | `c818c449e0` | `bac3aac89878b3e3` | `6b712571` 两侧同 blob |

- 三处 sha256 与 r2/r3 交件自报的 `3d949a47...` / `9c4b7381...` / `bac3aac8...` **逐字同字**；四枚 `git status --porcelain` **空**、
  工作树字节＝HEAD blob。函数体逐字 `return proc.SealableRoot(t.TempDir())`，即票 119 生产路那一枚叶子
  （`internal/proc/envfork.go:160`；生产调用点 `cmd/wisp/doctor.go:263`、`cmd/wisp/secret.go:148`）。
- NOTE 批 2 有**两枚** helper（`package tools` 与 `package tools_test` 孪生）：孪生不是本腿新写，是 2b-2 批 `dfa3dc44` 就入树的既有对
  （`loop_approval_test.go` 在 `tools_test` 里看不见内包 helper）。两枚函数体**一枚未改**，只多 28 处调用 => `R-125-2` 副本账不增行。

### (5) AC#5（门禁）— **成立；四数分开报，我自己现跑**

宿主 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`：

| 发 | RUN | PASS(顶层) | FAIL(顶层) | SKIP(顶层) | SUB-PASS | SUB-FAIL | SUB-SKIP | rc |
|---|---|---|---|---|---|---|---|---|
| `config -count=2 -v` | **390** | **228** | **0** | **0** | 162 | 0 | 0 | 0 |
| `tools -count=2 -v` | **658** | **492** | **0** | **0** | 166 | 0 | 0 | 0 |
| `config -count=1 -v`（倍数核） | 195 | 114 | 0 | 0 | 81 | 0 | 0 | 0 |
| `tools -count=1 -v`（倍数核） | 329 | 246 | 0 | 0 | 83 | 0 | 0 | 0 |

=> 每个数都是 `-count=1` 那发的**正好 2 倍**（无缓存、无 flake）；与编排者报的 `390/228/0/0`、`658/492/0/0` **逐数相同**（我独立跑到同数）。
- `gofmt -l internal/config internal/tools cmd/wisp` 在**纯净快照**上 **0 行**；`gofumpt -l`（绝对路径 **v0.12.0**）**0 行**。
  NOTE 工作树上这两把尺各报 1 行 `cmd/wisp/models.go`——那是本机 `core.autocrlf=true` 的既有 CRLF（334 对 CRLF / 14213 字节，
  HEAD blob 纯 LF / 13879 字节，`git hash-object` 与 `git rev-parse HEAD:` 同为 `4259a20f`）=> NOT 本批产物，r3 已按同形登记。
- 容器原生（linux）：`go build ./...` rc=**0**；`go vet ./internal/config/ ./internal/tools/` rc=**0**；`go vet ./cmd/wisp/` rc=**0**；
  `gofmt -l` 三包 **0 行**。宿主：`go build ./...` rc=**0**、`go vet` 三包 rc=**0**。
- `sh scripts/d22scan.sh` rc=**0** clean（八 scope `244/42/116/25/39/116/610/145`，`ban #8 design/` 39 含未跟踪 `design/doubao/`）；
  `sh scripts/check-path-length-budget.sh` rc=**0** GREEN（`tracked 17527 / over-budget=60 / covered=60 / not in roster=0`）。
- **`cmd/wisp` 整包 rc=1 按 `A1100` 只作"按名作差 + 逐发登记死点"**：宿主整包 `-count=1 -v` 得
  `RUN_TOP=160 RUN_SUB=58 PASS=158 FAIL=1 SKIP=0 SUBPASS=58 SUBFAIL=0 rc=1`；那 1 枚＝`TestAC14GoSideEvalPushReachesThePage`
  （`panel_resident_windows_test.go:944/946` 逐字同址同文），死点＝`TestPanelThreadIsSTAAndExitsCleanly`（恰 1 枚 RUN 无判决行）；
  该文件 blob 在 `e4d363f2`/`6b712571`/`HEAD` **三处同一串** `f19c77a5`。容器 `-count=2 -v` 软链形＝`342/206/136/0`＋`108/18/0`，
  恰是 `-count=1` 那发的 2 倍（`171/103/68/0`＋`54/9/0`）。

### (6) AC#6 判 NOT 翻 — **判对；且★它要求的缺口在我这程收口前已被票 359 真闭合（我实测）**

- 零交集：见 (1)。AC#6 要动的是 `internal/winsec` 的解析语义分母；三批写面是 `internal/{config,tools}`＋`cmd/wisp` 的用例侧递根点。
- **前提我先在打批 3 的树上量**：把 `sameAbsoluteness` 两条合取摘掉（`grep -c '!sameAbsoluteness('` => 0、加 `var _ = sameAbsoluteness` 守可达性）
  => 容器 POSIX 跑 `./internal/winsec/`：基线 `74/46/0/0`＋`28/0` rc=0，变异 `74/46/0/0`＋`28/0` rc=0 => **零枚红、逐数相同**
  （正控：同一株在宿主 Windows 有 6 枚红、`build failed` 命中 0）=> 票面 AC#6 的前提**在 `6b712571` 上确实成立**。
- ★**但收口前局面变了**：票 `359` 于 20:25–20:32 落地（`7835d67f` 加 `internal/winsec/sameabsoluteness_359_other_test.go`，
  `4ca369d5` 只改票面），我在 20:4x 收尾时现量到它。**我在 359 树上重打同一发变异**（`0cea1927` 快照）：
  基线 `RUN=75 PASS=47 FAIL=0 SUBFAIL=0 rc=0`，摘掉两条合取 => `75/46/1/0 rc=1`，红名**恰 1 枚**＝
  `TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359`。**单点摘**也各有 1 枚红（只摘 `:456` => 1 红；只摘 `:552` => 1 红，同一名）
  => 这一形在 POSIX 侧**已有分母、且两个调用点各有一枚牙**，不再"可以被删掉而无人察觉"。
- 归口更正（重要）：**359 的落地把 AC#6 的活干完了**，但它是**另一张票的另一条腿**，NOT 本票三批的产物：三笔 `git show --name-only`
  只动 `internal/winsec/sameabsoluteness_359_other_test.go` 与 `docs/evidence/s1/359-*`、票面；`internal/winsec/resolve.go` 的 blob
  在 `6b712571`/`7835d67f`/`4ca369d5`/`0cea1927`/`HEAD` **五处同一串** `e0b0e22f` => 生产码一枚未动（合 AC#6 自己的"只动用例侧"地界）。
- 所以本票 AC#6 的处置**仍然判 NOT 翻**（那格写的是本票要做的活，做它的是 359），但**理由从"缺口还在"改成"缺口已由 359 闭合、不在本票名下"**。
  票面行号小偏（不翻格，只记）：`sameAbsoluteness` 一枚定义 `:519`、两枚调用点 `:456`/`:552`；票面点的 `:513-517` 落在注释里，"POSIX 侧"那两行是 `:516-518`。

### (7) 红名册逐名 + 新增红 0 + 豁免三枚 — 逐名见 part2 sec 5

新增红（三包双向名册差集）**0 名**。豁免三枚逐名在 part2 sec 5 且**在 HEAD blob 上有名有址**（`wiring_test.go:112/:182/:291`），
批 2 对该文件**零 hunk**；`cmd/wisp` 我一手量到两形都红**顶层 68 + 子 9**（NOT 28，见 sec 2 的口径更正）。

### (8) 写面与纪律 — **成立；正控/负控都打了**

- 逐笔 `git show --name-only`：见 part2 sec 6。**28 枚路径全是 `*_test.go`**、生产码 **0**、`internal/{winsec,proc,risk}` **0**、
  `thresholds.go`/任何 golden/`allowlist.txt`/`docs/PLAN.md`/`docs/specs/**`/`tools/d22scan/**` **0**；
  `git cat-file -e <commit>^:<file>` 逐枚核**无一枚新文件** => 没有"新写一枚用例"这种扩面。
- 三笔码提交的票面主题前缀**逐字合规**；七笔里被禁码位 **0**（扫 `U+00A7/2713/2264/2260/2192/21D2/26D4/26A0/2460-24FF/1F000-1FAFF/FE0F`，
  正控样本 3 命中、交付件 11 份各 0 命中）；七笔**一枚都没上远端**（`git branch -r --contains` 全空）=> "只 commit NOT push"成立。

## 2. 逐格结论（七格已翻：我全部追认；AC#6：判 NOT 翻对）

| 格 | 已翻 | 我的判语 |
|---|---|---|
| AC#1（分母做成可重跑读数） | 是 | **追认**。读数仍在（`124-ac1-denominator-readings.md`），30 包分母与两形拆分可复跑 |
| AC#2（接票 119 纪律） | 是 | **追认，但票面口径要改**（见下） |
| AC#2b（只转"可转"那批） | 是 | **追认**。三批 131 枚落地，我量到 34/38/37 三本分母逐名对上 |
| AC#3（不许放宽换绿） | 是 | **追认**。删除侧 52 行全是根表达式；决定性一条两支尺都过 |
| AC#4（变异自证） | 是 | **追认**。我自己重打三发变异，本批目标名**逐名全回归**：config 34/34、tools 38/38、wisp 37/37 |
| AC#5（门禁） | 是 | **追认**。四数我自己跑到同数，倍数核过，容器/宿主双 GOOS 都干净 |
| AC#6（winsec 绝对性分母） | **未翻** | **NOT 翻正确**；前提我在打批 3 的树上复现（POSIX 零红 + Windows 正控 6 红），拆 359 是该拆。★收口前 359 已落地并真闭合该缺口（我实测 1 红），但那是**另一票的另一条腿**，不归本票三批 |

**两处口径更正，必须随结案一起落账（票面 append-only 由编排者写，我不改）：**

1. **票面 AC#2/AC#2b 的"软链形红名数 ＝ 拒绝腿那一档的枚数"按字面 NOT 可达**：AC#2a 账上"拒绝腿"是 **0**，
   而落地后 `internal/tools` 仍留 **3** 枚红。这三枚证据链是硬的（票 123 具名豁免、`wiring_test.go` 零 hunk、改前改后逐名同一、
   失败文案是 `审批超时（300 秒未确认），C18 一律判拒绝`）=> 豁免判得对，**但判据须改写成名册集合口径**，NOT "枚数等于某档"。
2. **"`cmd/wisp` 两形都红 28 枚"是低估**：我一手量到**顶层 68 + 子测试 9**（`POST-P` 与 `POST-L` 双向 0 行差），
   与 `124-v1` 14:3x 那次量的 68+9 **逐数相同**，与 r3 交件 sec 8 的名册 set-equal。
   => 豁免名册 `3 + 29 + 4` 里那个"29"应改写成 `1（panel，本腿未复测）+ 68（cmd/wisp 顶层）`，子测试 9 另列。

## 3. 判决：票 124 能不能结案

**分两层，别混（混了就是 `A1099`(2) 那族"把批了当结了"的镜像）：**

- **就"三批工作"：能结案。** 三批出口三批全过（`34=>0`、`38=>3`、`37=>0`，新增红 **0**，名册集合双向 0 行）；
  AC#1..AC#5 七格我全部追认；三发变异我自己重打、目标名逐名全回归；门禁我跑到与自报同数。
- **就"票面八格全勾"：NOT 能自动结 —— 缺"AC#6 归属"的落账那一行。** 票面 AC#6 仍 `[ ]`（我逐字核过 HEAD blob：7 勾 1 未勾），
  它的活已由 `359` 做完并结票（359 侧五格全闭、文件名已带 `-done`）。=> **结案前提**：编排者把"AC#6 移交 359 且 359 已闭、
  本票不再由自己闭合这一格"写进本票 Progress log 与台账。`A1104` 已有"拆出"的意思，但**票面正文没有一行说"本票按 7 格结"**，
  也没写"AC#6 的关闭条件＝359 落地"——而 359 今天**已经落地**，所以这一格今天具备翻的条件，只差那一行落账。**这一行不由我写。**
- 另建议同批落账 sec 2 那两条口径更正，否则下一位拿"红名数 = 0"这把尺复算必然读到"未达"。
