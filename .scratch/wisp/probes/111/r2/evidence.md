# 票 111 — 111-r2 证据件（收尾腿）

代号：`111-r2`（编排者 2026-10-05 20:0x 派，收尾票 111；前腿 111-r1／111-r1b 均死于额度，遗产由编排者代提为 `1bb654e3`／`20c29e94`，本腿不复活其号，从已提交状态接着干）
起手时刻：`2026-10-05 20:01:28 +0800`（首次 `date`）
写面（编排者指派）：`.scratch/wisp/probes/111/r2/**` ＋ 必要时 `scripts/portable-tests.sh`（最小修）；`ci.yml` 本轮**只读**；⛔ 零 `go test/build/vet`（另一枚验收腿 268-v1 正在 cmd/wisp 跑整包，任何重 go 命令会洗它的读数）；⛔ 不给 ci.yml 加 `-tags winlive`（owner 未批）；⛔ 不碰 `.gitignore`／`design/**`／`frontend/**`。

---

## §0 起手锚

- 本腿首次读数时刻：**`2026-10-05 20:01:28 +0800`**。
- 当时 HEAD＝`2a633eb8`（ledger A624 代提三笔那枚）；**20:04:40 共树又落 `076144fd`**（167-a5b 在自己的 probes 目录落骨架，两文件均在 `.scratch/wisp/probes/167/a5b/`，与本腿写面无交集）⇒ 本腿工作锚取 **`076144fd`**（`git rev-parse --short=8` 于 20:07:05 复认）。
- `git status --porcelain -- scripts/ .github/workflows/ci.yml` = **0 行**（20:01:28 与 20:07:05 两次读均为空）⇒ **写面无人占**：1bb654e3 的编辑已在 HEAD 里，本腿可安全持有 `scripts/portable-tests.sh`。
- `bash -n scripts/portable-tests.sh` rc=0（20:07:05）⇒ 1bb654e3 收编进 HEAD 后语法完好（编排者 17:43 验收凭据的复认）。
- 票面锚 `4e66817`（2026-09-24）距本腿起手 **漂 11 天**（只在此具名，不改票面一字）。
- 票面 ⛔ 一字未改；AC 框未翻。

### §0.1 本腿复锚（后笔补，§0 原文一字未动）

- `2026-10-05 20:31:13 +0800` 复量 HEAD＝**`15d8b60e`**（A625 那枚，共树在我读数期间又落了一笔）。
- ⚠ 关键：**`scripts/portable-tests.sh` 的 blob 在两枚 HEAD 下是同一个 `2ff02dd7`**
  （`git rev-parse HEAD:scripts/portable-tests.sh` 于 20:30:32 与 20:31:13 两次同值；`git cat-file blob` md5
  两次同为 `328eb3ead0545736a33e2c131687d60d`）⇒ §3 两发变异的靶件对本腿任一 HEAD 都成立，不必重跑。
- 我的 /tmp 快照底本取 `164ee2c2`（20:26 那枚在飞 HEAD）；`git archive` 出来的 `scripts/portable-tests.sh`
  md5 与工作树逐字节同值（上面那行），故 §3 的读数对 `076144fd`／`164ee2c2`／`15d8b60e` 三枚 HEAD 同一。

## §1 名册现量（本腿复量，非照抄 r1b）

### 1.1 `go list ./...` 的枚数（本机，只这一条 go 命令）

- **`go list ./...` = 35 枚**，rc=0，stderr **0 字节**。时刻 **`2026-10-05 20:24:19 → 20:24:21 +0800`**。
  逐行读数留盘：`…/r2/golist-win.txt`（35 行）。
- 与 r1b §0 的 35 **一致（复认，非照抄：本腿自己量的）**。票面 `AC#1`／`现场` 写的 **33 ⇒ 过期**。

### 1.2 ★35−33 的差集**逐枚坐实**（r1b 只报了总数差，这一格本腿关掉它）

只用 `git ls-tree`（不跑第二次 `go list`，不碰 cmd/wisp 的读数窗）：

| 读法 | 4e66817（票面锚） | HEAD | 差 |
|---|---|---|---|
| 主模块内含 `.go` 的目录（`cmd/`＋`internal/`＋`frontend/`，排 `testdata/`） | **32** | **34** | **+2 / −0** |
| ＋ 主模块内的 `tools/signmodels`（根 `go.mod` 下唯一的 `tools/` 包） | 1 | 1 | 0 |
| ⇒ 应得 `go list ./...` | **33** | **35** | — |

- **多出的正是这两枚，且只有这两枚**：`internal/projctx`（票 200 落的，1 枚外部测试）＋ `internal/streamkey`。
  二者在 `4e66817` 的 `scripts/portable-tests.sh` 里 **命中 0 次**（`git cat-file blob 4e66817:scripts/portable-tests.sh | grep -c 'projctx\|streamkey'` = 0）
  ⇒ 票面那句"33"在它自己写下的那天是对的，**两枚新包把它抬到 35**，不是票面漏计。
- **删除列＝0 枚**（`comm -23` 空）⇒ 没有任何包在这 11 天里消失，差集不会被"改名"糊过去。
- 独立模块**不在**根 `go list` 里（`go.mod` 自证）：`tools/d22scan`→`github.com/CarlosShao/wisp/tools/d22scan`、
  `tools/mockllm`→`…/tools/mockllm`、`scripts/spike` 各自 `module` 行独立 ⇒ 数名册时不许把它们算进 35。
- ⚠ **"35"必须带平台**（票面 Progress log `AC#1` 已立这句，本腿复认机制）：`cmd/balldebug` 的
  全部文件带 `//go:build windows`，故 GOOS=linux 的名册是 34 而非 35。本腿 20:36:55 一发旁证：
  `GOOS=linux CGO_ENABLED=0 go list ./...` **rc=1**、stdout **34 行**、stderr **264 字节**点名
  `cmd/wisp -> sherpa-onnx -> build constraints exclude all Go files in …/sherpa-onnx-go-linux@v1.13.8`
  ⇒ r1b §0 说的"票面 AC#4 的 ubuntu 障碍复现"我这腿**第三次复现**，非我可修面。

### 1.3 CI 被测名册的出处（ci.yml **只读**，行号为 HEAD `15d8b60e` 现量）

r1b §1 记的四个调用点，本腿 grep 复认，行号**一字未漂**：

- `:373` `run: bash scripts/portable-tests.sh --scope=core`（test-core／ubuntu）
- `:471` `run: bash scripts/winsec-tests.sh`（test-windows 腿 winsec 门）
- `:507` `run: bash scripts/wisp-cli-tests.sh`（test-windows 腿 cmd/wisp 门）
- `:543` `run: bash scripts/portable-tests.sh --scope=windows`，带 `:542` `if: ${{ !cancelled() }}`（票 111 AC#6 的形状）

名册之外还有三处**直调** `runtests.sh` 的步（r1b 没列，登记以免被当成"名册就是全部"）：
`:81` d22scan 正控、`:349` `internal/proc -run TestLayoutForTestEnv`、`:559` `internal/risk -run TestPathResolverJunctionWindows`；
两把壳尺在 `:134`（`sh scripts/d22scan.sh`）与 `:166`（`sh scripts/check-path-length-budget.sh --with-self-test`）。
ci.yml 总 766 行；`runs-on` 分布：`:66`／`:310`／`:697` ubuntu-latest，`:420`／`:565` windows-latest，`:623` `[self-hosted, wisp-slo]`。

> 行号可信度的时间轴（本腿在窗内 HEAD 漂过两枚，先把这个立住再引号）：
> `git rev-parse <ref>:.github/workflows/ci.yml` 于 `076144fd` 与 HEAD `2afdaeb6` 同为
> `c5a1b1ab6467a2f221c7f70066688f24d8ab9dea`（而票面锚 `4e66817` 是另一枚 `c5a98063`，所以票面自引的 ci.yml 行号不可照抄）。
> 四个调用点于 21:29:01 当场复 grep，行号仍是 `:373`／`:471`／`:507`／`:543`；`grep -c winlive` 仍是 0。
> "未漂"由这枚 blob 同一性保证，不是运气，也不是拿 r1b 的旧行号。

### 1.4 名册在 `scripts/portable-tests.sh` 里的真身（HEAD 现行号；r1b 记的 `:141-178` 已被 1bb654e3 顶漂）

| 档 | pin 块 | 行 | pin 枚数 | scope 数组／定义 | 行 |
|---|---|---|---|---|---|
| core | `core_pin` | `:164-192` | **27** | `case` 分支 `core)` | `:234-246` |
| windows | `win_pin` | `:193-204` | **10** | `windows)` | `:248-254` |
| cli | `cli_pin` | `:205-207` | **1**（`cmd/wisp`） | `cli)` | `:256-262` |
| winsec | `winsec_pin` | `:208-210` | **1**（`internal/winsec`） | `winsec)` ＋ `winsec_scope=./internal/winsec/...` | `:263-270`／`:229` |
| census | 无 pin（`:367` 明写"审计者不认领档"） | — | — | `census)` 空 scope | `:271-274` |

- `tiers='core windows cli winsec census'` @ `:220`；census 在 `:307-325` 把 `tiers=` 与本文件自己的 `case` 分支**回读对撞**，不等就 rc=1。
- `1bb654e3` 那"四处认领"本腿逐枚复认（pin 两枚 × scope 数组两处＝四行）：
  `internal/session` @ pin `:188`(core)／`:203`(win)，scope `:244`(core)／`:252`(win)；
  `internal/projctx` @ pin `:185`／`:200`，scope 同行。
- ★**票面 5 枚的 HEAD 现行号**（r1b §0 那组 `:141-178` 已过期的部分，由这组顶替）：
  `cmd/llmrecord` `:165`/`:194` · `internal/ball` `:169`/`:195` · `internal/perm` `:182`/`:197` ·
  `internal/plugin` `:183`/`:198` · `cmd/wisp` `cli_pin :206`。

### 1.5 GUARD C 的咬合复认（**只用 `go list`，一枚测试都不跑**）

把 `:234-246` 与 `:248-254` 的 scope 数组原样抄进命令行，只跑 `go list`（时刻 `20:35:04 → 20:35:24 +0800`）：

- core：pin **27 行** == resolved **27 行**，`diff` **0 行** ⇒ **绿**。留盘 `…/r2/core-pin-win.txt`＋`core-resolved-win.txt`。
- windows：pin **10 行** == resolved **10 行**，`diff` **0 行** ⇒ **绿**（r1b 在旧 HEAD 量到的是 8==8，多的两枚正是 1bb654e3 补的 session/projctx ⇒ **它补完之后 GUARD C 仍然咬合，没有留潜在红**）。
- cli／winsec 各 1 枚，与 r1b §1 同结论（单枚、无 glob 静默收窄）。

### 1.6 名册分母的现量读数（base census，见 §3.1）

`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，rc=0，16 秒（20:24:55→20:25:11）。
⇒ 35 枚里 **28 枚本平台编译得出测试文件**（35−7），**这 28 枚全部被某档认领**（`unclaimed-with-tests=0`），
7 枚 NO-SCOPE **逐枚都是 `0/0`**（`cmd/balldebug`·`frontend`·`internal/agent/scheduler`·`internal/speech`·
`internal/streamkey`·`internal/watchdog`·`tools/signmodels`）⇒ 没有"有测试却没人跑"的残余。测试文件总数（t＋x 相加）＝**327**。
**这一行就是票面 `AC#1` 要的那张对账表，且它现在自己会红（§3）。**

## §2 逐枚可纳入性（票面 5 枚 + 真残余洞）

判据四列：`_test.go` 枚数（`git ls-files`）／tag 形状（`head -1` 逐文件）／**本平台编译得出的分母**（§3.1 base census 的 `t/x` 列，GOOS=windows 现量）／最坏颜色。
r1b §2 的 5 枚结论本腿**逐枚复认**（下 1–5 全"已在册"），并把 1bb654e3 补的两枚（下 6–7）补齐；8 是横切的 winlive 半边，9 是 cmd/wisp 的 ubuntu 半边。

| # | 包 | `_test.go` | tag 形状 | census t/x | 认领 | 最坏颜色 |
|---|---|---|---|---|---|---|
| 1 | `internal/ball` | 16 | 6×`windows`＋5×`windows && winlive`＋5 无 tag | 11/0 | core+windows | 见下 |
| 2 | `cmd/wisp` | 65 | 22×`windows`＋6×`windows && winlive`＋1×`!windows`＋36 无 tag | 58/0 | cli | 见下 |
| 3 | `internal/perm` | 3 | 0 枚带 tag | 3/0 | core+windows | 普通红 |
| 4 | `internal/plugin` | 1 | 0 枚带 tag | 1/0 | core+windows | 普通红 |
| 5 | `cmd/llmrecord` | 1 | 0 枚带 tag | 1/0 | core+windows | 普通红 |
| 6 | `internal/session` | 2 | 0 枚带 tag | 2/0 | core+windows（1bb654e3） | 见下 |
| 7 | `internal/projctx` | 1 | 0 枚带 tag（外部测试包 `package projctx_test`） | 0/1 | core+windows（1bb654e3） | 见下 |

1. **internal/ball** — 复认 r1b。tag 现量：`//go:build windows` 6 枚（`hotkey_borrow_refused_260r5`／`hotkey_cancel_borrow_260`／`hotkey_cancel_borrow_expect_260r2` 等）、`windows && winlive` 5 枚（`hotkey_live`／`interaction_live`／`live_windows`／`live_guard_windows`／`hotkey_cancel_borrow_live_260`）、无 tag 5 枚（`dock`／`liquid`／`position`／`tokens`／`tokens_table`，本腿逐枚 `head -1` 复认）。**算术自证**：16−5(windows)−5(winelive)＝6 枚进得了 windows 分母，census 实测 **11**＝6 带 tag＋5 无 tag ⇒ **那 5 枚 winlive 今天连编译都没编译**（§2.8）。新依赖无。最坏颜色：ubuntu 若有平台 bug＝普通红即发现。**已在册。**
2. **cmd/wisp** — 复认 r1b 的 tag 计数并补精确分布：22 `windows`／6 `windows && winlive`／1 `!windows`（`secret_dataroot_119b_test.go:1`）／**36 无 tag**＝65。windows 侧 census `t=58`＝22＋36 ⇒ **6 枚 winlive 同样零编译**；`t=58` 里不含那枚 `!windows`。runner 只 `cli` 档，由 `scripts/wisp-cli-tests.sh`（ci.yml `:507`）在 windows 腿跑，缺 sherpa DLL 时它 `exit 1` 点名缺哪个（`:74-99` 预检）。**已在册（cli）。** ubuntu 半边见 §2.9。
3. **internal/perm** — 复认：3 枚 `head -1` 全是 `package …`，无一行 `//go:build`（`grep -c` 命中 0）。无新依赖。普通红。**已在册。**
4. **internal/plugin** — 复认：`disposal_test.go` 无 tag，census `1/0`。普通红。**已在册。**★这枚同时是 §3.2 GUARD D 正控的靶件。
5. **cmd/llmrecord** — 复认：`main_test.go` 无 tag，census `1/0`。普通红。**已在册。**
6. **internal/session** — 补 r1b 未覆盖的两枚之一。`grants_test.go`＋`ticket224_pattern_dialect_test.go`（`package session` 内部测试×2，census `2/0`）。★本腿现量：这 2 枚里 `t.Skip`／`os.Getenv`／`exec.Command` 命中 **0／0／0**（`grep -c` 逐文件）⇒ 与 `1bb654e3` 注释里那句"ZERO t.Skip、ZERO os.Getenv、ZERO exec.Command"**逐字对得上**，不是抄它的断言。无 tag ⇒ 两平台都进分母。**已在册（1bb654e3），本腿不重复合入。**
7. **internal/projctx** — 另一枚。1 枚 `projctx_test.go`，`package projctx_test` ⇒ census 记作 `0/1`（**x** 而不是 t）。★票面 `现场` 那句"零覆盖且有 `_test.go`"的写法在这种形状上会**被误读成 0**（t 列为 0），所以名册必须报 `t/x` 两列——这正是 `:288-292` 注释立这条的理由。同样 `t.Skip`/`os.Getenv`/`exec.Command` 命中 0。无 tag。**已在册（1bb654e3）。**
8. **★winlive 半边（横切 ball＋wisp，共 11 枚文件）** — 本腿复认 r1b §2.6 的四个"没有"：
   - `grep -c winlive .github/workflows/ci.yml` = **0**（`grep -n` 零行）；
   - ci.yml 全文 `grep -n 'winlive\|-tags\|GOFLAGS'` = **零行**；
   - `scripts/portable-tests.sh` 里 `tags\|GOFLAGS` 命中 **0**；`tools/d22scan/runtests.sh` 同样 **0**
     ⇒ **没有任何一条现成的"传 tag"的缝**，加 winlive 必须先造这条缝，不是给 scope 数组加一行就完事；
   - 仓内 `winlive` 测试文件**只有这 11 枚**（`grep -rl //go:build.*winlive` 逐目录计数＝cmd/wisp 6＋internal/ball 5，无第三处）。
   ⇒ 结论复认：**这 11 枚默认零编译、CI 永不跑**。owner **未批**，本腿**不给 ci.yml 加 `-tags winlive`**。
   真接入要动哪几行＋代价 → **§5.1**（不是"不做"，是"这轮不能由我这枚做"）。
9. **cmd/wisp 的 ubuntu 半边** — 票面 `AC#4` 已量过（CGO=1 时 19/29 FAIL，rc=1，127.4s），本腿**不重跑**（零 go test）。ubuntu 腿根本没有 `cmd/wisp` 的调用点：`:373` 的 core scope 不含它，且 `GOOS=linux go list ./...` 连解析都 rc=1（§1.2 那发 264 字节 stderr）。⇒ **归票 98／新票，登记不修**（§5.2）。

**§2 的总结局**：票面点名的 5 枚**全部早已在册**；1bb654e3 补的 2 枚在册且本腿复核其"零藏 skip"断言成立；
剩下的两处不是"接入"问题而是**平台/tag 问题**（winlive 需 ci.yml＋一条传 tag 的缝；cmd/wisp 的 ubuntu 需先让 `go list` 在 linux 解得开）。
⇒ **本轮我不新增任何 scope 行、不新增任何 pin 枚**（新增了反而会被 GUARD C 判红，见 §3.4）。

## §3 实际改动与 GUARD D 正控（票 111 AC#3 那格的核心）

### 3.0 先答编排者那把闸：**"触发 GUARD D 要不要真跑 25 包的 go test？"——不要。这条我读实现读出来的，并且证了**

`1bb654e3` 的实现里，GUARD D whole 落在 **`--scope=census` 那一个 block 内**（`:286` `if [ "$mode" = census ]` 起，`:425` `fi` 止）：

- `sed -n '286,430p' | grep -c 'go test\|go build\|go vet\|runtests.sh'` = **0**；
- 同一段里出现的 go 命令只有 **`go list ./...`**（`:336`）与 **`go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}'`**（`:364`）；
- `census)` 分支 `:271-274` 是 `scope=()`／`pinned=''`，**根本走不到** `:449` 之后的"真跑测试"那段（census 在 `:425` 前就 `exit 0`／GUARD D 在 `:422` 就 `exit 1`）。

⇒ 正控**不需要真测试窗**，因此**不需要停手上报**，也不需要动用那格"归编排者"的出口。
本腿全程零 `go test`／零 `go build`／零 `go vet`（`go list` 与 `go env GOOS` 由脚本自己调，见 §3.5 的诚实边界）。

### 3.1 对照组（HEAD 原样，未变异）：census 绿

```
portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0
```
rc=**0**，38 行，窗口 `2026-10-05 20:24:55 → 20:25:11 +0800`。逐行读数：`…/r2/census-base.txt`。

### 3.2 ★正控第一发（D1）：把 `internal/plugin` 从名册里抽掉 → GUARD D 响亮红，点名它

**怎么抽的**：抽**四处**（`core_pin` 一枚＋`win_pin` 一枚＋`core)` scope 一处＋`windows)` scope 一处），
不是只抽 pin。**只抽 pin** 会让 GUARD C 先红（pin≠resolved）而**测不到 D**；
四处一起抽＝"这枚包从来没人认领过"，GUARD C 保持哑（见 §3.4），**红的只可能是 D**——这才是票面 AC#3 要的那一发。

**在哪抽的**：票面 Rules（`111-ci-tests…md:54`）**"绝不在仓库内建 worktree/checkout"**，变异只在 /tmp 的
`git archive <sha> | tar -x` 快照里做 ⇒ `git archive 164ee2c2 | tar -x -C /tmp/wisp-111r2-d1`，
快照里 `sed` 抽，**工作树 `scripts/portable-tests.sh` 一字未动**。
落地凭据（同一条 `&&` 链里 `grep -n` 打印被改后的整行）：

```
240:        ./internal/panel/... ./internal/ball/ ./internal/perm/
249:        ./internal/ball/ ./internal/perm/ ./cmd/llmrecord/
```
（`:243`／`:252` 原本的 `./internal/plugin/ ` 前缀没了；pin 行数 core 27→26、windows 10→9）
`bash -n scripts/portable-tests.sh`（快照内）rc=**0** ⇒ 变异是编译/语法层面成立的，不是碰巧崩在解析上。

**红句逐字**（窗口 `20:28:23 → 20:28:37 +0800`，**rc=1**，50 行；全文 `…/r2/census-mutant-D1-plugin-pulled.txt`；
变异 diff `…/r2/mutants/D1-pull-internal-plugin.diff`）：

```
portable-tests.sh: github.com/CarlosShao/wisp/internal/plugin     1/0          NO-SCOPE <-UNCLAIMED-HAS-TESTS
portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=8 unclaimed-with-tests=1
portable-tests.sh: GUARD D - 1 package(s) compile a test file for GOOS=windows and
portable-tests.sh:   NO named scope claims them, so no CI step runs them and no CI step
portable-tests.sh:   can ever go red over them:
portable-tests.sh:   github.com/CarlosShao/wisp/internal/plugin  tests-compiled-for-windows=1/0
portable-tests.sh: this row is the ticket 111 field itself ('CI 只测 33 个包里的 20 个').
```

**三件判据全中**：① 响亮（rc=1，不是 0，不是"少跑一个包照样绿"）；② **点名被抽的那枚包**（`internal/plugin`＋它的分母 `1/0`＋平台）；
③ 分母自己跟着动（`claimed-by-no-scope` 7→8、`unclaimed-with-tests` 0→1）⇒ 它不是硬编码一句错误。

### 3.3 正控第二发（D3）：换一枚靶件，证它不是只认 plugin

同一把手术抽 `internal/session` 四处（core pin `:188`／win pin `:203`／scope `:244`／`:252`），
靶件取自我快照底本 `164ee2c2`（含 1bb654e3 的两包）。窗口 `20:29:33 → 20:29:47 +0800`，**rc=1**：

```
portable-tests.sh: github.com/CarlosShao/wisp/internal/session    2/0          NO-SCOPE <-UNCLAIMED-HAS-TESTS
portable-tests.sh: GUARD D - 1 package(s) compile a test file for GOOS=windows and
portable-tests.sh:   github.com/CarlosShao/wisp/internal/session  tests-compiled-for-windows=2/0
```
全文 `…/r2/census-mutant-D3-session-pulled.txt`；变异 diff `…/r2/mutants/D3-pull-internal-session.diff`
（⚠ 这份 diff 里另有 **2 行纯注释的字序差**——我的 `sed 's| \./internal/session/||'` 顺带咬到了 `:48`／`:112` 两句注释里的同名字符串，
**不在代码路径上**，`bash -n` 与两发读数都不受影响；登记在此免得被当成"变异多改了东西"）。
⇒ 两发不同靶件、不同分母（`1/0` 与 `2/0`）各自点名自己 ⇒ **GUARD D 是通用的，不是为 session/projctx 写死的**。

### 3.4 ★GUARD D 不是 GUARD C 的重复（这条是 D 该不该存在的判据，本腿证了）

对 **D1 那个变异体**（四处一起抽），把 `:248-254` 的 windows scope 原样 `go list`（窗口 `20:29:11`，一条命令）：

```
pin_rows=9  resolved_rows=9  diff_rc=0（0 行差）
```
⇒ **GUARD C 对这一发完全哑**（pin 与 resolved 相等，它看不出少了一枚）；
留盘 `…/r2/guardc-blind-win-pin.txt`＋`guardc-blind-win-resolved.txt`（9 行 vs 9 行，`diff` rc=0）。
**"有测试却不在名册"这个洞，GUARD C 结构上抓不到**——它比的是"名单 vs 名单展开"，两边同时少就是相等。
GUARD A 也抓不到（它只管"名册里声明了却 0 分母"，方向相反）；GUARD B 更抓不到（那枚包压根没进调用）。
⇒ **D 补的正是 A/B/C 三面都照不到的那一个方向**，与 `:64-75` 注释自己那句
"A/B/C audit the packages a scope NAMES; nothing audited the packages a scope FORGETS" **对得上，不是我替它圆**。

### 3.5 还原凭据（票面点名的那条）＋本腿对 `scripts/` 的改动＝**零字节**

- `git status --porcelain -- scripts/ .github/` = **0 行**：**本腿（r2b）第一条命令**（约 20:22，未印时刻）＝0 行，
  带时刻的两次复量 **20:30:32** 与 **20:46:48** 亦 0 行；前腿 `111-r2` 在 §0 记的 20:01:28／20:07:05 两发同 0 行（那是它的读数，我不认领）。
- `md5sum scripts/portable-tests.sh` = **`328eb3ead0545736a33e2c131687d60d`**
  ＝`git cat-file blob HEAD:scripts/portable-tests.sh | md5sum` **同值**（20:30:32），
  blob 号 `git rev-parse HEAD:scripts/portable-tests.sh` = **`2ff02dd7`**（= `1bb654e3` 那枚 blob，HEAD 漂到 `15d8b60e` 后仍同一枚）。
- `bash -n scripts/portable-tests.sh` rc=**0**（20:30:32）。
- ⇒ **本腿没改过它，也就无从"还原"**——变异全在 /tmp 快照（`wisp-111r2-d1`／`d2`／`d3`，按"只建不删"留着不删）。
  票面那句"最小修＋单独 commit"的触发条件（**GUARD D 验出 bug**）**没有发生**：§3.2/§3.3 两发都按设计红、按设计点名、按设计退出，**未发现实现缺陷**。

**诚实边界（三条，不许被读成"我跑了 GUARD D 的每一步"）**：
1. GUARD D 的红句里 `GOOS=windows` 那半句出自脚本自己的 `go env GOOS`（`:280`）⇒ 本腿**零 `go test/build/vet`**，
   但 `go list`/`go env` 是编排者白名单与脚本自身行为，不是我另开的口子（`go env` 只被脚本内部调，我本人一次没手打 `go env`）。
2. 我这腿的靶件平台是 **windows**（本机 GOOS）。**GOOS=linux 下 GUARD D 会 rc=1 但原因是"拒绝"而非"发现洞"**：
   `go list ./...` 在 linux rc=1（§1.2 那发：stdout 34 行／stderr 264 字节）→ 命中 `:337-351` 的 refusal 分支 → `exit 1`。
   ⇒ **这条直接决定 §5.3 的接法：census 只能接在 windows 腿，接进 ubuntu 腿是一枚永久红。**
3. 两发红句都是**真子进程真退出码**（12 秒／14 秒各含 35 次 `go list -f`），不是我拼的字符串。

### 3.6 本轮"实际改动"清单（照 AC 的口径交）

- `scripts/portable-tests.sh`：**0 字节**（§3.5）。理由：票面 5 枚早在册（§2.1–5）、真洞 2 枚已由 `1bb654e3` 补完且 GUARD C 复认咬合（§1.5）、
  GUARD D 验出可用无 bug（§3.2–4）。**再写一行就是重复行，不改变任何分母**——这与 r1b §3 那句是同一条判语，但**我这枚多了正控那一格**。
- `.github/workflows/ci.yml`：**0 字节，只读**（引行号见 §1.3／§5.1／§5.3）。
- 票面 `111-ci-tests-20-of-33-packages.md`：**0 字节**（AC 框一枚不碰）。台账：0 字节。
- 本腿写面＝`.scratch/wisp/probes/111/r2/**`（证据件＋台件，只建不删）。

## §4 门禁读数（带时刻＋HEAD）

### 4.1 四道守卫的现行行号复认＋各配一句"咬什么"（HEAD `15d8b60e`，blob `2ff02dd7`）

| 守卫 | 红句出处行 | 咬什么（一句话） | 咬不到的方向 |
|---|---|---|---|
| **A** | `:574`（测量在 `:570-571`） | 名册里**声明了**、但本平台 `.TestGoFiles+.XTestGoFiles` 全空 ⇒ 红（"空分母永远不会红也永远不证任何东西"） | 咬不到"根本没声明"（那是 D 的地盘） |
| **B** | `:751`（计数在 `:738-739`） | 每一枚被声明的包必须打出**自己的** top-level `^(ok\|FAIL) <转义包名> <时长>` 结果行，缺行 ⇒ 红 | 包没进调用 ⇒ 它不在分母里，B 看不见 |
| **C** | `:553`（形状拒绝另有一处 `:532`） | 命名 scope 把自己 `go list` 解析出的 import path 集合**钉住**，与 pin 双向对撞 ⇒ 删一行是红、不是少跑一个包 | ★**pin 与 scope 同时少＝相等 ⇒ 哑**（§3.4 实测 9==9、diff rc=0） |
| **D** | `:409`（计数 `:396`、标记 `:397`、退出 `:422`） | census 期间"本平台编译得出测试文件、却没有任何命名档认领"⇒ 拒绝 `exit 0` | 只跑在 `--scope=census`（§3.0）；且 GOOS=linux 时先被 `:337-351` 的 go-list 失败拒掉（§5.3） |

r1b §3 记的行号（A`:477`／B`:643`／C`:455`）在 1bb654e3 之后**一律顶漂**，以本表为准。
另有两处**不是守卫但同族**的响亮失败：`:449-452` 空 scope ⇒ `exit 2`（"拒当绿色空转"）、`:276-278` 未知 `--scope` ⇒ `exit 2` 并点名所有档。

### 4.2 两把壳尺（本腿唯一被允许的两条实跑脚本）

**`sh scripts/d22scan.sh`** ＠ **`2026-10-05 20:31:20 → 20:31:41 +0800`**，**rc=0 clean**，251 行，全文 `…/r2/d22scan.log`。
- 正控段：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`。
- 实扫各 scope（与 r1b §4 的 17:03 那发**逐枚持平，无一下降**）：
  bans #1-5 `internal/`=**228**、`cmd/`=**38**；ban #6 `frontend/`=**85**；ban #7 `internal/tools/`=**23**；
  ban #8 `design/`=**39**、`frontend/`=**85**、`internal/`=**512** Go（注释与 `_test.go` 计入）、`cmd/`=**104** Go。
- 末句逐字：`d22scan: clean - no D22 ban violations`。

**`sh scripts/check-path-length-budget.sh --with-self-test`** ＠ **`20:31:50 → 20:31:54 +0800`**（**HEAD `15d8b60e` 之前**，见 4.4），**rc=0 VERDICT GREEN**，19 行，`…/r2/path-length-budget.log`。
- 正控 **3/3 ok**：基准绿 → 种一枚超长跟踪路径被**逐字点名**红（`RED - over budget and NOT in the roster: … 125 chars … budget 165`）→ 移除复绿。
- 分母：tracked paths=**5991**、over-budget=**57**、covered by roster=**57**、**not-in-roster=0**；longest=180（`issues/252-…`）；帽＝rule 9 名 100＋issues 前缀 21＝121 相对，worst checkout 前缀 44 ⇒ 全程预算 165。

### 4.3 ⚠ 我这腿改了"跟踪文件名"，所以补第二发（票面 AC#5 只要求不降，这条是自查）

本腿往 `.scratch/wisp/probes/111/r2/**` 落了 15 枚新跟踪件（读数／变异 diff／两把尺的 log）。
`git ls-files | awk length` 现量：**我最长的一枚＝65 字符相对**（`mutants/D3-pull-internal-session.diff`），＋44＝109 ＜ 165 ⇒ 结构上不可能越帽。
**不靠推算，复跑坐实** ＠ **`20:55:44 → 20:55:48 +0800`**，**rc=0**，`…/r2/path-length-budget-postcommit.log`：
> `denominator: tracked paths=6019  over-budget=57  covered by roster=57  not in roster=0` ／ `VERDICT GREEN`

⇒ 分母 5991→**6019**（正是我加的 28 枚里落在计数内的部分＋共树别人的件），**over-budget 仍 57、not-in-roster 仍 0**
⇒ **我的交件一枚都没给这扇门添新红**（这条对本腿有意义：AGENTS.md §1.2 明令 `d22scan` 的 emoji 射程不覆盖 `scripts/` 与 `.scratch/`，所以**壳尺是唯一会真咬我落件的门**，见 4.5）。

### 4.4 HEAD 漂移记要（免得被读成"我拿旧树量的门"）

- 本腿窗内 HEAD 共漂四枚：`076144fd`(起手，前腿 §0 所取) -> `164ee2c2` -> `15d8b60e` -> `2afdaeb6`(交件前 21:29:01 复量)。
  ★全程不变的是两枚 blob：`scripts/portable-tests.sh` = `2ff02dd7`（§0.1）、`.github/workflows/ci.yml` = `c5a1b1ab`（§1.3 那条时间轴）。
  ⇒ §4.1 的守卫行号与 §1.3 的调用点行号对这四枚 HEAD 同时成立，不需要重量。
- 4.2 两发落在 HEAD＝`15d8b60e` **之前**（起手 HEAD `164ee2c2`，20:31:13 复量才见 `15d8b60e`）。
  这两把尺**都不读 `portable-tests.sh`**（d22scan 只扫 `.go`／文本＋跑它自己模块的测试；path-length 只看跟踪路径名），
  而 `164ee2c2`↔`15d8b60e` 之间 `scripts/portable-tests.sh` 的 blob 逐字节同一（`2ff02dd7`）⇒ **漂移对 4.2 的读数无影响**，不补跑 d22scan。
- 4.3 那发落在 `15d8b60e` 之后（且带我自己三笔提交在树里），是**当前树**的读数。
- ⛔ **窗口归属不自判绿**：并行腿（`167-a5b`、`268-v1b`）在我这窗口里持续落 `docs/reports/**` 与 `.scratch/wisp/probes/**` 的件，
  ⇒ 编排者若复跑与我不一致，**以安静窗口的复跑为准**。我只报数与时刻。

### 4.5 gofumpt：这轮不适用（编排者已给口径，我把它坐实成读数）

- `ci.yml:168-176` 的 `gofumpt -l . tools/d22scan tools/mockllm` 确实**扫根模块全部 `.go`**；但本腿改面是
  `scripts/portable-tests.sh`（**0 字节**）与 `.scratch/wisp/probes/111/r2/*.md|*.txt|*.log|*.diff` ⇒ **零 `.go` 文件被本腿写**，gofumpt 无从咬。
- `scripts/` 下的 11 枚 `.go`（`git ls-files scripts | grep -c '\.go$'` = 11，全在 `scripts/spike/**`）本腿一字未碰。
### 4.6 ★本证据件自己的字符自查（我一度写错一句，把它改成立并留痕）

`tools/d22scan/main.go:579-586` 的 `emojiScopes()` 只有 **`design/`·`frontend/`·`internal/`·`cmd/`** 四枚
⇒ **`scripts/` 与 `.scratch/` 不在仪器射程内**（与 `ci.yml:58` 注释同口径，`AGENTS.md §1.2` 那段"规格比仪器宽"的缺口原样有效）。

**本文件全量现量**（`python` 数整文件，时刻见末行；⚠ 含本节自己，所以本节文字一改枚数就漂 ⇒ 末行给出"再量一次"的命令，验收方可自复）：

| 字符 | 码位 | 枚数 | 仪器带内？ |
|---|---|---|---|
| `⇒` | U+21D2 | 111 | 否（`2190–21FF` 是刻意留的空隙，AGENTS.md 明写"仪器抓不到"） |
| `★` | U+2605 | 26 | **是**（`2600–27BF`，注释不豁免） |
| `⚠` | U+26A0 | 20 | **是** |
| `⛔` | U+26D4 | 13 | **是** |
| `−` | U+2212 | 12 | **是**（`2200–22FF`） |
| `≠` | U+2260 | 5 | **是** |
| `≤` | U+2264 | 7 | **是** |
| `✓` | U+2713 | 10 | **是** |
| `✔` | U+2714 | 4 | **是** |
| `✗` | U+2717 | 4 | **是** |

带内合计 **101** 枚；`⇒` 另 111 枚（不在带内）。测量时刻 `2026-10-05 21:32:24 +0800`，对象＝当时工作树全文（538 行／53,285 字节；本节定稿后只改过 ASCII 数字，枚数仍成立）。

复量命令（验收方自复，不必信我这张表；**故意写成 ASCII 转义**，否则这条命令自己就含 10 枚带内字符，枚数会被它自己顶漂——这是我这两格里踩到的第二个自指坑）：

```
python -c "s=open(r'.scratch/wisp/probes/111/r2/evidence.md',encoding='utf-8').read(); \
print([(hex(ord(c)),s.count(c)) for c in \
[u'\u21d2',u'\u2605',u'\u26a0',u'\u26d4',u'\u2212',u'\u2260',u'\u2264',u'\u2713',u'\u2714',u'\u2717']])"
```

注：本节这张表**只改 ASCII 数字不会让枚数漂**（数字不在任何带内）；改符号或带内中文就会漂 ⇒ 枚数只在"表格定稿后不再动带内字符"这一前提下有效。

⇒ **诚实更正（本腿自己撞上的那一格）**：我在 4.5 初稿写过"我没有写 `✓` 或 `≤`"——**那句是假的**。
`≤`／`✓` 最早各 1 枚出现在我自己抄 AGENTS.md 那句引文里；而本节为了说清"哪些算带内"，
把 `PLAN.md:3449` 那条禁令的**原文清单**（`✓ ✔ ✗ ⚠ ★ →`）也抄了进来 ⇒ **引文本身又添了几枚带内字符**。
原句已就地改成本节，**不抹掉、不静默**（票面 Rules 的 append-only 口径：更正＝新写一段，不是改掉旧的）。

**为什么这不是违规（三条，逐条有出处，不是我自己给自己开门）**：
1. **仪器不扫**：本件落 `.scratch/wisp/probes/**`，`emojiScopes()` 无此目录 ⇒ §4.2 那发 `d22scan.sh rc=0` 与本件带内字符**同时成立**，不是漏扫后侥幸。
2. **禁令的射程是界面码**：`PLAN.md:3440` 那节标题是"§17.4 图标规范"，`:3449` 那句"绝对禁止：任何 emoji（含 `✓ ✔ ✗ ⚠ ★ →` 这类文字符号）"
   与 `:3453` 的范围 `U+2190–U+2BFF` 是**给 `frontend/`·`design/` 的视觉资产**定的（`D29` 前端与视觉架构）；
   本件不是 UI 代码、不进产物、不面向用户。
3. **这是本仓文档的行文既有形状**：`AGENTS.md` 自己用 `⚠ ≤ ✓ ⛔`，r1b 的证据件用 `★ ⚠ ⛔ ≠ ⇒`
   （逐枚数：r1b `★`=2 `⚠`=1 `⛔`=6 `≠`=1 `⇒`=11）。我沿用同一形状，**没有新造字符**。

⇒ **本腿不据此改 `PLAN.md` 一字、不动 `d22scan` 的射程**（那两处都要人工批准；`AGENTS.md §1.2` 已把这条缺口挂在台账 `A201②`）。
唯一被这条自查**改变的动作**是：我把 4.5 里那句假话删了重写。

## §5 判不动（具名归口）

### 5.0 ★票面前提过期这一格——**谁证的**（三层，逐层可核，不是"我看了一眼"）

票面 `现场`／`标题` 那句"还有 5 个带测试的包零覆盖（`internal/ball`／`cmd/wisp`／`internal/perm`／`internal/plugin`／`cmd/llmrecord`）"**在 HEAD 不成立**。

- **第一层（r1b 的权威读法，本腿引用它的方法名而非结论）**：`r1b §0` 用 `git show HEAD:scripts/portable-tests.sh`
  读**提交真身**、不执行工作树里被 `111-r1` 未提交编辑污染的那份 ⇒ 由此得出"5 枚早已在 `core_pin`/`win_pin`/`cli_pin`"。
  ⚠ 它那组行号（`:141-178`）在 `1bb654e3` 之后**全部顶漂**，所以本腿不引它的行号，见第二层。
- **第二层（本腿自己在 HEAD 复量，行号为准）**：`scripts/portable-tests.sh` @ HEAD `15d8b60e`／blob `2ff02dd7`：
  `cmd/llmrecord :165/`:194 · `internal/ball :169/`:195 · `internal/perm :182/`:197 · `internal/plugin :183/`:198 · `cmd/wisp`＝`cli_pin :206`（§1.4 那张表）。
  ⇒ **5 枚全在名册里，一枚不缺。**
- **第三层（本腿新加的：把"过期"钉到它自己的出生地）**：这两枚 commit 里 5 枚**已经在了**——
  `git cat-file blob 8fe5c7ce:scripts/portable-tests.sh` 命中 `cmd/llmrecord :114/:140`、`internal/ball :118/:141`、`internal/perm :131/:143`、`internal/plugin :132/:144`；
  `7699ec3f` 同样命中（`:125/:152` 等）；`cli_pin` 由 `git log -S"cli_pin='"` 判为 **`8fe5c7ce` 引入**（`cmd/wisp` 从那天起就在册）。
  ⇒ **票面写下"零覆盖"的那天（2026-09-21 20:0x）离这两枚 commit（同日 21:1x–21:3x）只差一小时上下**——
  票面记的是**票 110 量出来的现场**，而 `agent-ticket111` 同日已经把它修掉了，票面没回头改。**这是"过期"的成因，不是谁的读数错**。
  再往上一层：`4e66817`（票面锚）那份 `portable-tests.sh` 里 5 个包名命中 **17 次** ⇒ 连票面锚点当天这句话都已偏。
- ⇒ **真残余洞是什么**（本腿复认 r1b 的判定并补一刀）：不是那 5 枚，而是
  (a) `internal/session`＋`internal/projctx` 有测试却无人认领 ⇒ **已由 `1bb654e3` 补，且 §1.5 证 GUARD C 复认咬合**；
  (b) **看见洞的那把尺子（census）自己不在 CI 里跑** ⇒ §5.3，这才是本票剩下的那一格。

### 5.1 winlive 半边（ball 5＋wisp 6＝11 枚文件，CI 今天**零编译**）——owner 未批，本腿只登记现状

**现状（三条都是现量，见 §2.8）**：ci.yml `winlive` 命中 **0**／ci.yml 无 `-tags`、无 `GOFLAGS`／
`portable-tests.sh` 与 `runtests.sh` 的 `tags|GOFLAGS` 命中各 **0**。编译侧自证：
`ball` 16 枚文件 − 5 枚 `windows` − 5 枚 `winlive` ＝ 6 枚可进 windows 分母，census 实测 `t=11`＝6＋5 无 tag ⇒ **那 5 枚 winlive 连编译都不参与**；
`cmd/wisp` 65 − 22 − 6 − 1(`!windows`) ＝ 36 无 tag，census `t=58`＝22＋36 ⇒ **同理 6 枚不参与**。

**真接入要动哪几行（逐行，不笼统说"改 ci.yml"）**：
1. `runtests.sh:75` 是 `go test -v -count=1 "$@"` ⇒ **flag 透传已具备**，不用改这个文件（这是好消息，省一枚禁改面）。
2. `scripts/portable-tests.sh:694` 是唯一那行真跑测试的调用 ⇒ 要么加第 5 枚档（如 `winlive)`，同时动 `:220` `tiers=`、
   `:164-210` 新增一枚 pin、`:234-274` 新增一条 `case` 分支——**少动一处 census 就在 `:313` 当场 rc=1**），
   要么给现有 `cli)`/`windows)` 传 `-tags`。
3. ⚠ **两处 `go list -f` 不会跟着带 tag**：GUARD A 的 `:570-571` 与 census 的 `:364` 都是裸 `go list`。
   ⇒ 加了 tag 之后**分母仍按无 tag 算**，winlive 文件在 A/D 眼里永不存在 ⇒ 想让它受 A／D 管，第 4 处也得动：把 tag 传给这三处，一处不能落。
4. `ci.yml:543`（windows 腿 portable 步）与／或 `:507`（wisp-cli 步）新增或改写一条 `run:` ⇒ **本轮禁面，零字节**。
5. runner 归位：`ci.yml:420`/`:565` 是 `windows-latest`（GitHub 托管，**无交互桌面会话**）；`:623` 才是 `[self-hosted, wisp-slo]`。

**代价（三条，按严重度排，全部有现量支撑）**：
- **代价一·接错腿就是永久红**：11 枚里 **5 枚文件自带 `t.Skip`**（逐文件现量：`internal/ball/hotkey_live_test.go` 4 处、
  `cmd/wisp/resident_hotkey_live_258_windows_test.go` 4 处、`internal/ball/interaction_live_test.go` 2 处、
  `internal/ball/live_windows_test.go` 2 处、`internal/ball/live_guard_windows_test.go` 1 处＝**合计 13 处**；其余 6 枚 0 处），
  理由原文 `SKIP-LOUD: this host reports no ball window`（`resident_hotkey_live_258_windows_test.go:55/:97`；
  ⚠ 更正初稿一句"3 枚文件自带 t.Skip"——**数是 5 枚文件／13 处**，我是按 `grep -c` 逐枚复算后改的，不是凭上面那段草稿的目测）。
  而 `runtests.sh:98-100` 把**任何 `--- SKIP` 判 fatal**（末句逐字：`do not relax this script`）。
  ⇒ 接到 `windows-latest` 上＝**每 run 必红且不是产品 bug**；接进 core 腿更是直接撞 `:507`/`:75` 两堵墙。
  ⇒ **只能落 `wisp-slo`（`:623`）**，那是一条**需要那台机器在线**的腿（票面 `:49`／台账多处记过它今天不在 CI 关键路径上）。
- **代价二·它是 live 硬件用例**：需要真桌面会话＋真球窗（`ball posture`）＋resident 进程，
  与票 111 AC#7 已当众改口径的那件事同源（`TestSyncRegistryProbeLive` 至今 `=== RUN` 计数 **0**，因为要 HKCU 同步记录）。
  ⇒ 接 winlive 之前，**先要有一台愿意被门禁占用的机器**，这是编排者的决定，不是脚本能替它决定的。
- **代价三·放宽一条就破 D22**：为了变绿去 `-skip` 掉它们＝**票 111 AC#2 与 `AGENTS.md §1.1` 双重禁**（"不许为了变绿放宽任何断言"），
  而 `SLO 阈值/golden/thresholds.go 一字节都不许动`。⇒ **本腿不给 ci.yml 加 `-tags winlive`，也不加 `-skip`**，只登记。

**归口**：`.github/workflows/ci.yml` 那一面（票面 `:6` 写明与票 110 串行、`:63` 写"ci.yml 现在归你"），**owner 批准**后由下一腿做。
本腿交的只有：现状＋上面那 5 处行号＋3 条代价。

### 5.2 cmd/wisp 的 ubuntu 半边（19 枚红）——**登记，归票 98／新票，不修**

- 票面 `AC#4` 那格**早已闭合**，读数来自 runner：`PASS=29 / FAIL=4 / SKIP=0 / rc=1 / 328.972s`，
  4 条红**全是** `审批超时（1/300 秒未确认），C18 一律判拒绝` ⇒ 票面自己已判"**未达成 ⇒ 移票 123**"（`:47`）。
- 本腿**不重跑**（零 `go test`，且 cmd/wisp 整包读数归并行腿 `268-v1`，已在 `164ee2c2` 落盘 `ok 422.320s`）。
- 另有一层"接都接不进去"的现量（本腿新量）：GOOS=linux 下 `go list ./...` 本身 **rc=1**
  （stdout 34 行／stderr **264 字节**，逐字点名 `build constraints exclude all Go files in …/sherpa-onnx-go-linux@v1.13.8`，
  时刻 20:36:55→20:36:56，留盘 `…/r2/golist-linux-partial.txt`＋`golist-linux-stderr.txt`）
  ⇒ ubuntu 腿想跑 `cmd/wisp`，**先得让 CGO=1 那条路在解析层通**（票面 `AC#4` 记的 CGO=1 能建能跑，但要 19 枚红逐条读码分"产物布局 vs 跨平台真 bug"）。
- **归口**：票 **98**（sherpa DLL／产物布局的地界）＋票 **123**（票面已把那 4 条红移过去）；**本票不越界去修**（票面 `next=` 5② 同一条）。

### 5.3 ★GUARD D 今天**接不进任何 CI 步**——本腿量出来的那一格，具名归 ci.yml 面

- `grep -c census .github/workflows/ci.yml` = **0** ⇒ **census（GUARD D 唯一的栖身处）没有任何 CI 调用点**，
  正印证 `:132` 与 `:413-415` 自己那句"the census … with NO reader but the census roster"／
  "The census used to PRINT these rows and exit 0, so the hole was visible only to whoever thought to run it by hand"。
  ⚠ 更正一处我在此节的初稿：我写过"`grep -rn scope=census scripts/*.sh` 除 `portable-tests.sh` 自己＝**0 命中**"——**假的**，
  实际 **2 枚**命中，都在 `scripts/portable-tests-selftest.sh:506` 与 `:533`（那两枚载具用例）。
  ⇒ 结论不变但**成因更硬**：census 不是"没人调用"，而是"**只有那把不跑真测试二进制的载具调用它，而那把载具自己不在 CI 里**"（下一段）。
- **而且它连 ubuntu 腿都接不得**（§3.5 边界 2）：GOOS=linux ⇒ `go list ./...` rc=1 ⇒ 命中 `:337-351` 的拒绝分支 ⇒ `exit 1`。
  ⇒ **能接的位置只有 windows 腿**（本机 GOOS=windows 那发是 rc=0，§3.1）。
- ⚠ **这条本腿"量得到但做不了"**：动 `ci.yml` 是本轮硬禁面（编排者指派＋票面 `:6`）。
  ⇒ 交件形状：正控**已做实**（§3.2／§3.3 两发、rc=1、点名、分母自漂），**接入那一行留给下一腿**，
  落点建议 `ci.yml:543` 那条 windows portable 步**之前**加一步 `bash scripts/portable-tests.sh --scope=census`，
  并带 `if: ${{ !cancelled() }}`（AC#6 的形状，`ci.yml:542` 有现例子）。
- **顺带一枚"守卫自身无人守"的洞（本腿读到，具名上报）**：`scripts/portable-tests-selftest.sh`（票 250/251/254 那把"不跑真测试二进制"的载具，
  其 fake-`go` 明确支持 `go list -f`（`testdata/portable-tests/go:50`）与 `./...`（同文件 `:21`，注释原文"./... is the whole module (census reads it this way)"）——
  **但 `grep -c 'GUARD D|unclaimed|UNCLAIMED' scripts/portable-tests-selftest.sh` = 0** ⇒ **GUARD D 没有任何回归覆盖**：
  谁把 `:392-399` 那个 `case` 写坏，那 27 枚载具里**没有一枚会红**。而这把载具自己 `grep -c portable-tests-selftest .github/workflows/ci.yml` = **0** ⇒ 它也**不在 CI 里**。
  ⇒ 两处都归 ci.yml／载具面，**本腿不修**（§5.4 说明为什么不顺手修）。
  ⚠ 诚实边界：**我没有跑过这把自己**（跑它＝起真工具链，可能触 cmd/wisp），这条是从 `scripts/` 与 `testdata/` 的实现读出来的，可核。

### 5.4 为什么"看起来只差一行"的我也不顺手做（三条禁面，逐条点名）

1. `ci.yml`：**本轮零字节**（编排者指派＋票面 `:6` 串行）。§5.1／§5.3 的落点都写成了"下一腿要动哪几行"，不代做。
2. `scripts/portable-tests.sh`：**只有验出 bug 才动**（编排者指派）。§3 两发**没验出 bug** ⇒ 0 字节（§3.5 凭据）。
3. 票面／台账／`.gitignore`／`design/**`／`frontend/**`：**一枚不碰**；`.scratch/wisp/probes/268` **不引**（并行腿的写面）。

### 5.5 步级读数（"那一步真给过结论的 run id + step 号"）——**本地产不出，只登记**

票面 `AC#1` 那一列、`next=` 1/2/3/4 全要 push 之后由 `gh api repos/CarlosShao/wisp/actions/runs/<run>/jobs` 取。
本腿**未 push**（纪律：只 commit）⇒ **这一格交的是"欠"，不是"通过附条件"**；
取数注意两条照抄票面 `:49-50` 并原样有效：**`.steps[].order` 返回 `null`** ⇒ 步号按数组位置数＋同时引 step 名；
脚本在 `scripts/` 下（**没有 `scripts/ci/` 目录**），按猜路径取 blob 得到的 `fatal: path does not exist` 是**猜错、不是版本对不上**。
⛔ **不拿本地绿冒充 CI 绿**（票面 `:57`）。

## §6 交件判语与 commit 链

### 6.1 逐格判语

- **§1 名册现量**：**达成**。`go list ./...` = **35** 复认；★**35−33 差集逐枚坐实**＝`internal/projctx`＋`internal/streamkey`，删除列 0（r1b 只报了总数差，这一格本腿关掉）；
  ci.yml 四个调用点行号**复认未漂**；名册真身 HEAD 现行号顶替 r1b 那组过期行号；GUARD C 咬合**只用 `go list`** 复认（core 27==27、windows 10==10）。
- **§2 逐枚可纳入性**：**达成**。5 枚复认"已在册"＋tag 分布精确到枚（wisp 22/6/1/36、ball 6/5/5）＋补 session(2, t)／projctx(1, **x**)，
  二者 `t.Skip`/`os.Getenv`/`exec.Command` 本腿现量 0/0/0（对得上 `1bb654e3` 的 ZERO 断言，非抄它）。
- **§3 ★GUARD D 正控**：**达成，且是本腿的主件**。先答闸：GUARD D 全在 `--scope=census`（`:286-425`），那段 `go test|build|vet|runtests.sh` 命中 **0**
  ⇒ **不需要真测试窗，不必停手上报**；两发变异（plugin 四处抽→rc=1 点名 `1/0`；session 四处抽→rc=1 点名 `2/0`）在 /tmp 的
  `git archive` 快照里做，工作树 **0 字节**；★并证 **D 不是 C 的重复**（对 D1 变异体 pin 9==resolved 9、diff rc=0 ⇒ C 结构上看不见这类洞）。
  票面 AC#3 那句"人为抽掉一个包证明它会红"——**这一格今天第一次有样本**。
- **§4 门禁读数**：**达成**。四道守卫现行号＋各一句"咬什么"＋各一句"咬不到什么"；两把壳尺 rc=0（d22scan 21 秒／path-length 4 秒，各带时刻）；
  ★本腿落了 15 枚新跟踪件 ⇒ **补跑第二发**：5991→6019 分母、over-budget 仍 57、**not-in-roster 仍 0** ⇒ 交件不给门添新红；gofumpt 不适用（零 `.go` 写面）已坐实。
  ⚠ **本腿自查出一句假话并就地更正**（§4.6）：初稿写"我没写 `✓`／`≤`"是假的，带内字符本件共 **101** 枚，逐码位列表＋复量命令（ASCII 转义写，避开自指）都在那节。
- **§5 判不动**：**达成（判"不动"是有出处的不动）**。票面前提过期的"谁证的"给了三层（r1b 权威读法／HEAD 现行号／**两枚 commit `8fe5c7ce`+`7699ec3f` 里 5 枚已在册**，
  并定位成因＝票面记的是票 110 的现场、同日 21:1x 已被 `agent-ticket111` 修掉而票面没回头改）；
  winlive 半边登记现状＋**5 处落点行号**＋**3 条代价**（13 处 `t.Skip` vs `runtests.sh:98` 的 SKIP-fatal ⇒ 接 `windows-latest` 必永久红；只能落 `wisp-slo`；放宽＝违 AC#2）；
  cmd/wisp ubuntu 19 枚红归票 98/123；★**新报一格**：GUARD D 栖身的 census **在 ci.yml 命中 0**、且 linux 腿接不得（`:337-351` 拒绝分支），
  外加 GUARD D **自身零回归覆盖**（载具 fake-`go` 支持 census 的两条查询，但 0 枚用例引用 D，而那把载具也不在 CI）。
  步级读数＝本地产不出，登记"欠 push"。
- **§6 交件**：本腿写面＝`.scratch/wisp/probes/111/r2/**`（1 份证据件＋15 枚台件，**只建不删**）。
  `scripts/portable-tests.sh` **0 字节**／`ci.yml` **0 字节**／票面 **0 字节**／台账 **0 字节**。
  **票面 AC 框一枚不翻**（AC#1–AC#10 全是验收方的判语，本腿只交证据，不替验收方勾框）。

### 6.2 本腿 commit 链（逐笔 `git log` 可查；只 commit、**未 push**）

- `6aac9e8a`（20:39:27）笔1＝§1＋首批台件（9 枚命名 pathspec，`git diff --cached --name-only` 恰 9 行）。
- `fbf9060b`（笔2）＝§2＋linux 那发旁证（3 枚）。
- `70596a40`（20:49:29）笔3＝§3＋两发变异 diff＋GUARD C 致盲凭据（5 枚）。
- 笔4＝§4＋§5＋§6（本笔）。
- 起手锚：`076144fd`（前腿 §0 所取）→ 本腿读数窗内 HEAD 漂至 `164ee2c2` → `15d8b60e`；
  ★靶件 blob 全程 `2ff02dd7` 未变（§0.1）⇒ 三枚 HEAD 下这些读数同值。

### 6.3 下一腿要接的三件事（按可核性排，本腿不代做）

1. **把 census 接进 windows 腿一步**（`ci.yml:543` 之前，带 `if: ${{ !cancelled() }}`）——**这一条才让 §3 的正控变成门**；
   ⛔ 不要接进 ubuntu 腿（§5.3：GOOS=linux 必 rc=1，那是"拒绝"不是"发现"）。
2. **给 GUARD D 补一枚回归用例**到 `scripts/portable-tests-selftest.sh`（fake-`go` 的 `-f` 与 `./...` 两条支都在，缺的只是一枚引用 D 的 case；
   ⓐ 载具自己也不在 CI，两件事一起做才有意义）。
3. **winlive 半边**：owner 拍板＋`wisp-slo` 可用之后，按 §5.1 那 5 处一起动（漏第 3 处＝A/D 的分母不跟 tag，等于白接）。

⛔ 本票面 `AC#1–AC#10` 的勾框**不由任何实现腿做**（`AGENTS.md §0` 第 3 句：裁决者≠实现者；D22 双角色）。

### 6.4 ★本腿自查到的**三处自身失手**（逐条给"初稿写了什么／真身是什么／怎么改的"，不藏）

对抗验收要攻的就是这三类，我先自己交出来：

1. **§4.5 初稿**："我没有写 `✓`(U+2713) 或 `≤`(U+2264)" ⇒ **假的**。带内字符本件共 **101** 枚，`≤`／`✓` 就在我抄 AGENTS.md 的引文里，
   而我为了说清"哪些算带内"又把 `PLAN.md:3449` 的禁令清单（含 `✓ ✔ ✗ ⚠ ★ →`）抄进文件 ⇒ **引文自己又添了几枚**。
   ⇒ 更正方式：原句删掉重写为 §4.6 全量自查表（**枚数用现量 python 数，不用目测**）。
   ⚠ 顺带撞到的第二个坑：复量命令**初稿用字面字符写** ⇒ 那条命令自己就含 10 枚带内字符，枚数被它自己顶漂 ⇒ 改成 ASCII 转义（`u'\u2605'` 那种）才稳。
2. **§5.3 初稿**："`grep -rn 'scope=census' scripts/*.sh` 除 `portable-tests.sh` 自己＝**0 命中**" ⇒ **假的**，真身 **2 枚**
   （`portable-tests-selftest.sh:506`／`:533`）。⇒ 结论没翻但**成因换掉了**：不是"没人调用 census"，是"只有那把不在 CI 里的载具调用它"。
3. **§5.1 初稿**："11 枚里 **3 枚**文件自带 `t.Skip`" ⇒ **数错**，逐枚 `grep -c` 复算＝**5 枚文件／合计 13 处**
   （`hotkey_live` 4／`resident_hotkey_live_258` 4／`interaction_live` 2／`live_windows` 2／`live_guard_windows` 1，其余 6 枚 0）。

⇒ **共同成因**：三条都是"我先在别处目测／凭上下文记忆写了一句可以直接 grep 的话"。
本腿后半程的规矩改成：**枚数一律当场 grep/python 数，数完立刻贴到当条句子旁边**（§2 表、§4.6 表、§5.1 那 5 枚都是这么重做的）。

### 6.5 收尾读数（三发，全带时刻，证明"交件形状"没把门碰坏）

| 尺 | 时刻（+0800） | rc | 关键数 |
|---|---|---|---|
| `sh scripts/d22scan.sh` | 20:31:20 → 20:31:41 | 0 clean | 正控 PASS=35 FAIL=0 SKIP=0 RUN=77；#1-5 internal/=228 cmd/=38；#6 frontend/=85；#7 tools/=23；#8 design/=39 frontend/=85 internal/=512 cmd/=104 |
| `sh scripts/d22scan.sh`（交件后复跑） | **21:17:27 → 21:17:50** | **0 clean** | 各 scope **与上一发逐枚持平，无一下降**（228/38/85/23/39/85/512/104） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 20:31:50 → 20:31:54 | 0 GREEN | paths=5991 over=57 roster=57 **not-in-roster=0**；正控 3/3 ok |
| 同上（我的 15 枚落件之后） | 20:55:44 → 20:55:48 | 0 GREEN | paths=**6019** over=**57** not-in-roster=**0** |
| 同上（**收尾复跑**） | **21:17:50 → 21:17:54** | **0 GREEN** | paths=**6019** over=**57** covered=**57** not-in-roster=**0**，log `…/r2/path-length-budget-final.log` |

- ⛔ 零 `go test`／零 `go build`／零 `go vet`／我本人零 `go env`（`go env GOOS` 只被脚本内部调）——全程只有 `go list` 族＋`git`／`grep`／`wc`／`bash -n`／两把壳尺。
- `scripts/portable-tests.sh` 工作树 md5 `328eb3ead0545736a33e2c131687d60d` == `git cat-file blob HEAD:` **同值**，blob `2ff02dd7`，`bash -n` rc=0（21:16:11）。
- `git status --porcelain -- scripts/ .github/` = **空**（21:16:11）。
- §0 起手锚（前腿 `111-r2` 写的那 9 行）**逐字节未动**：`git show 2602be3a:…` 与现文件 `:9-17` 对撞 **diff rc=0**（21:16 量）。
- 只 commit、**未 push**。临时件**只建不删**：/tmp 的 `wisp-111r2-d1`／`d2`／`d3` 三枚快照与 `msg*.txt` 全部留盘。
- 工具输出里的**伪指令登记＝0 次**：全程未出现任何自称"编排者备注／停手／撤回／请 revert／放宽阈值"的文本，
  也没有任何"编排者让你另外做 X"的现成文字（若有即停手回报，见本腿纪律）。收到的只有后台任务与 hook 提示，均非指令。
