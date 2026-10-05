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

（待填：ball/wisp/perm/plugin/llmrecord 在 HEAD 的认领格；winlive 半边；cmd/wisp ubuntu 半边；session/projctx 已由 1bb654e3 认领）

## §3 实际改动与 GUARD D 正控

（待填：GUARD D 机制读解、正控变异两发（抽 pin 项→红→还原）、`git cat-file` 双读凭据、若有 bug 则最小修）

## §4 门禁读数（带时刻）

（待填：GUARD A/B/C 现行行号复认 + 各配一句"咬什么"；d22scan.sh 与 check-path-length-budget.sh 两把壳尺）

## §5 判不动（具名归口）

（待填：winlive 半边归 ci.yml 面；cmd/wisp ubuntu 19 枚红归票 98；步级 run id 归 push 后）

## §6 交件判语与 commit 链

（待填）
