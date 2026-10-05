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

## §4 门禁读数（带时刻）

（待填：GUARD A/B/C 现行行号复认 + 各配一句"咬什么"；d22scan.sh 与 check-path-length-budget.sh 两把壳尺）

## §5 判不动（具名归口）

（待填：winlive 半边归 ci.yml 面；cmd/wisp ubuntu 19 枚红归票 98；步级 run id 归 push 后）

## §6 交件判语与 commit 链

（待填）
