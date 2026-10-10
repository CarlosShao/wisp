# 票 301 `AC#3` 只读普查腿 `301-a2` — 两形代价表（只裁形状，⛔ 未下结论，⛔ 未动门）

## 0. 两枚锚（本程一切读数都挂在这两号上）

| | 钟 | HEAD | `git status --porcelain -- internal cmd scripts .github tools` |
|---|---|---|---|
| 起手 | `2026-10-10 13:36:38 +0800` | `86d89478` | **2 行**＝` M internal/audio/wasapi_windows.go` ＋ `?? internal/audio/parse_wave_format_300_windows_test.go` |
| 交回 | `2026-10-10 13:48:16 +0800` | `41329475` | **0 行**（那五枚目录在工作树上干净） |

两号不同 ⇒ 按派单第 1 节那把新规矩，**同一把尺两个号各跑一发、两次读数都留**（⛔ 只留后一发）：

- 尺（逐字可重跑，两个号各跑一次）＝`git ls-tree -r --name-only <锚> | grep '_test\.go$'` 逐枚 `git show <锚>:<path> | head -3 | grep -c 'go:build windows'`
- 读数：`86d89478` ＝ **119 枚** tagged `_test.go`；`41329475` ＝ **120 枚**
- 差集（`diff` 两发原始输出，尺＝`diff <(sort census-raw.md|grep '^TAGGED') <(sort census-raw-second.txt)`）＝**恰好 1 行、只增不减**：`TAGGED internal/audio internal/audio/parse_wave_format_300_windows_test.go`
- 两号之间变了什么（尺＝`git diff --stat 86d89478..41329475 -- scripts .github internal tools cmd`）＝
  `internal/audio/parse_wave_format_300_windows_test.go |157 +` ＋ `internal/audio/wasapi_windows.go |21 ++-`，**其余四枚射程目录 0 字节**
- ★关键：`scripts/portable-tests.sh` 在两号上 **blob 逐字相同**（尺＝`git show <锚>:scripts/portable-tests.sh | md5sum` ⇒ 两枚都是 `328eb3ead0545736a33e2c131687d60d`）。⇒ **本件所有门/清单/ledger 的行号与文字在两枚锚下通用**；`ci.yml` 未进 diff ⇒ 同理。
- 票 300 那枚新用例此刻**已**在 HEAD 上（`41329475`，13:4x 由 `300-r1` 交入）。它**不**在 `86d89478` 的任何名册里；`41329475` 的 audio 行＝**3 枚 tagged 文件／10 枚顶层用例**（前一号＝2／9）。它有没有算进票 300／301 `AC#0` 那两张已裁名册＝**归编排者**（`AC#0` 是在 `772ff880` 上裁的）。

---

## 表① 全仓规模（甲形名册有多大）

**表头**：出自 **`86d89478`**（主名册）与 **`41329475`**（audio 那一行加一枚）；尺＝上节那把 blob 尺；用例尺＝同一批 blob 上 `grep -c '^func Test'`；**射程＝blob，⛔ 工作树**；认领关系尺＝逐枚读 `scripts/portable-tests.sh` blob 的四张 pin（`:164-192` core_pin 28 行／`:193-204` win_pin 10 行／`:205-207` cli_pin 1 行／`:208-210` winsec_pin 1 行）与四张 scope 清单（`:234-247` core／`:248-255` windows／`:256-262` cli／`:263-270` winsec）＋ census 的认领判据本身（`:380` `printf '%s\n' "$pin" | grep -qxF "$p"` ⇒ **认不认识＝看 pin，⛔ 看 scope glob**）。

树里 `_test.go` 总数（blob `86d89478`，尺＝`git ls-tree -r --name-only 86d89478 | grep -c '_test\.go$'`）＝**451 枚**。其中 tagged **119 枚**，落在 **22 枚目录**；但那 22 枚里有 **10 枚在 `.scratch/wisp/probes/**`**（15 枚 tagged 文件／31 枚顶层用例）——**它们⛔ 是枚 Go 模块包**（`go list ./...` 的宇宙到不了那里，见「仅读码推的」第 1 条）。⇒ **甲形名册的真实宇宙＝下表 12 枚包**。

| 包（＝目录） | tagged `_test.go` 枚数 @86d89478 | tagged 顶层用例枚数 | 该包 total `_test.go` | 被哪些 pin 认领（逐枚读 blob） | 那些档今天跑在哪个 job／GOOS | 它的 windows 那一半在 CI 上**可被求值**吗 |
|---|---|---|---|---|---|---|
| `cmd/wisp` | 42 | 168 | 见注 a | `cli_pin`（`:205-207`） | `test-windows`（`ci.yml:516-517`）经 `scripts/wisp-cli-tests.sh:113` → `--scope=cli`（`ci.yml:604`） | **能** |
| `internal/winsec` | 17 | 54 | | `core_pin` ＋ `winsec_pin`（`:208-210`） | `test-windows` 经 `scripts/winsec-tests.sh`（`ci.yml:568`）＋ `test-core` | **能** |
| `internal/ball` | 11 | 45 | | `core_pin` ＋ `win_pin` | `test-windows` `--scope=windows`（`ci.yml:780`）＋ `test-core` | **能** |
| `internal/tools` | **10** | **27** | 53 | **只有 `core_pin`** | **只有 `test-core`（ubuntu，`ci.yml:402-403`／`:470`）** | ⛔ **不能**（与 audio 同形） |
| `internal/risk` | 6 | 26 | | `core_pin` ＋ `win_pin` | 同上两档 | 能 |
| `internal/proc` | 7 | 23 | | `core_pin` ＋ `win_pin` | 同上 | 能 |
| `internal/secret` | 3 | 17 | | `core_pin` ＋ `win_pin` | 同上 | 能 |
| `internal/models` | **3** | **3** | 17 | **只有 `core_pin`** | **只有 `test-core`** | ⛔ **不能** |
| `internal/audio` | 2 →（@41329475：**3**） | 9 →（@41329475：**10**） | 7 → 8 | **只有 `core_pin`** | **只有 `test-core`** | ⛔ **不能**＝本票 `AC#1` 要修的那一枚 |
| `internal/config` | 1 | 3 | | `core_pin` ＋ `win_pin` | 同上 | 能 |
| `internal/agent` | **1** | **1** | 41 | **只有 `core_pin`** | **只有 `test-core`** | ⛔ **不能** |
| `internal/memory` | **1** | **1** | 11 | **只有 `core_pin`** | **只有 `test-core`** | ⛔ **不能** |
| **合计（12 枚模块包）** | **104**（@41329475：**105**） | **377**（@41329475：**378**） | | | | |

注 a：`cmd/wisp` 的 total `_test.go` 枚数本格⛔ 量（只量了 tagged 那一半与那五枚包的总数），⛔ 当已知。

**完全不在任何档里的 tagged 包＝0 枚**（尺＝12 枚包名逐枚对四张 pin `grep -qxF`；`unclaimed` 那一支 `:382-399` 在 windows 档上今天⛔ 收得到它们，因为**它们在 pin 里**）。⇒ 派单表① 里"哪些包完全不在任何档里＝GUARD D 已经看见的整包没人测，⛔ 混进本票的账"这一格**读数＝0 枚**；本票要收的口⛔ 是"没档认领"，**是"认领它的那档跑的平台上求值不到它"**——这两件事在表②③ 里差一整枚形。

那 5 枚"只有 core"的包里的 17 枚 tagged 文件**逐枚具名**（尺＝blob `86d89478`）；41 枚用例名在本件同目录 `table1-roster.md`：

- `internal/agent/spill_acl_windows_test.go`（1 枚）
- `internal/audio/capturelevel_windows_test.go`（1）／`internal/audio/hotplug_test.go`（8）
- `internal/memory/artifacts_junction_tripwire_windows_test.go`（1）
- `internal/models/acl_sid_121_windows_test.go`／`handoff_window_109_windows_test.go`／`no_seal_ruling_windows_test.go`（各 1）
- `internal/tools/`：`bridge_a18_kill_windows_test.go`／`bridge_junction_windows_test.go`／`fs_edit_ac4b_kill_windows_test.go`／`fs_edit_ac5_sweep_r4_windows_test.go`／`fs_staging_windows_test.go`／`paths_shortname_252_probe_test.go`／`paths_shortname_252_r1_test.go`／`paths_twocontainments_252_r2_windows_test.go`／`recycle_windows_test.go`／`task_pointer_authority_ac3_174r4_windows_test.go`（10 枚／27 用例）

### 表① 附：两枚尺陷阱（都现量过）

1. **文件名不带 `_windows` 而内有 tag** 的实据仍是 `internal/audio/hotplug_test.go`（tagged、8 枚用例）；反向也有——`internal/tools/paths_shortname_252_probe_test.go`／`_252_r1_test.go` 两枚同样不带 `_windows` 而内有 tag ⇒ **按文件名分档在两枚包上都漏**（本票现量新增第 2、3 例）。
2. **`go:build windows` 出现在第 3 行以后＝⛔ 一枚真 tag，全仓 5 例**（尺＝逐枚 blob 全文 `grep -n 'go:build windows'` 而 `head -3` 不命中，⛔ 射程目录）：`cmd/wisp/resident_ball_228_test.go:26`、`internal/ball/liquid_test.go:283`、`internal/config/c26_seam_posix_125_test.go:8`、`internal/models/assembly_reachability_121_test.go:22`、`internal/proc/crossvet_test.go:17` ⇒ 逐枚读前 6 行：五枚的 `package` 都在第 1 行（或 `//go:build !windows` 在第 1 行）⇒ 那 5 处全是**注释里的散文提及**，⛔ 是构建约束。⇒ **`head -3` 这把尺是对的；改成全文 grep 会虚增 5 枚**。⚠ 顺带一枚脆弱性：全文 grep 之所以⛔ 把 `c26_seam_posix_125_test.go:1` 的 `//go:build !windows` 也算进去，纯靠模式串里那个空格——**别把 `!windows` 的匹配当理所当然**。

---

## 表② 甲形（⛔ 动门：把"哪些包里有 tagged 用例"钉成一枚名册断言）

### ① 能长在哪几处 —— 每一处今天由哪个 job、哪个 GOOS 求值

| 候选落点 | 今天由谁求值（job／GOOS／调用点行号） | 它能不能"看见 windows 那一半"？一句话 |
|---|---|---|
| `scripts/portable-tests.sh` **census 那一支**（`:286-424`） | 唯一调用点 `ci.yml:744`，所在 job＝`test-windows`（键 `:516`），`runs-on: windows-latest`（`:517`），`env:` 只有 `WISP_ENV: test`（`:518-519`），全文件⛔ 任何非注释的 `GOOS=`（尺＝`grep -nE '(env|run|export).*GOOS|GOOS=' ci.yml \| grep -v '#'` → **0 行**）⇒ `goos=$(go env GOOS)`（`:280`）**＝windows** | 能——**而且已经能**：census 正跑在 windows 上，`:364` 的 `TestGoFiles/XTestGoFiles` 在 windows 上把 tagged 文件**算进分母**；它看不见的⛔ 是文件，**是"认领关系"这一层语义**（⇒ 甲形若落这里，代价⛔ 是"让它看得见 windows"，⛔ 是"给它一枚它已经有的视野"） |
| `scripts/portable-tests.sh` **四档 run 那一支**（`:427+`，GUARD C `:548-563`、GUARD A `:565-583`） | 三个 job 各自调：`--scope=core`＝`test-core`（ubuntu，`:470`）；`--scope=windows`＝`test-windows`（`:780`）；`--scope=cli`＝经 `wisp-cli-tests.sh:113`（windows，`ci.yml:604`）；`--scope=winsec`＝经 `winsec-tests.sh`（windows，`ci.yml:568`） | 只能看见**当前这一发的 GOOS**；同一枚断言在 ubuntu 那发上⛔ 看得见 windows tagged（这正是 GUARD A 的 `:565-566` 自我声明："the files that COMPILE INTO THIS PLATFORM's test binary"） |
| `tools/d22scan/runtests.sh` | 三处直调：`lint`（ubuntu，`ci.yml:82` `-C tools/d22scan ./...`）、`test-core`（ubuntu，`:445 ./internal/proc/ -run TestLayoutForTestEnv`）、`test-windows`（`:796 ./internal/risk/ -run TestPathResolverJunctionWindows`） | ⛔ **这枚形长不住**：`runtests.sh` 判的是**一次 `go test` 的输出**（`=== RUN`／`--- PASS`／`--- SKIP`，`:98`→`:102`；编排者已逐行读到），它⛔ 读源码字节 ⇒ 一枚"没有任何档去跑"的用例对它根本⛔ 存在。放在这里＝放一枚只会跟着别人的射程走的断言 |
| 一枚 `_test.go`（现形先例＝`internal/models/assembly_reachability_121_test.go`） | 该文件**无 tag**、在 `core_pin` 的 `internal/models` 里 ⇒ 今天由 **`test-core`／ubuntu／GOOS=linux** 求值（`ci.yml:470`），**同时**它自己在体内对 `windows`／`linux`／`darwin` 三枚 GOOS 各跑一发**元数据查询**（blob `:64-69` 逐字"THE PER-GOOS LEG IS NOT DECORATION"、`:120-123` `graphGOOS`）⇒ **一枚 ubuntu 上的 `_test.go` 是看得见 windows 那一半的**，靠的⛔ 是跑它，**是读字节或 `go list -e`** | 能（前提是断言用"读字节／元数据查询"，⛔ 用"跑"）；⚠ 三枚代价：(1) 它得落在**已被某档认领**的包里才有人跑，`core` 落在 ubuntu ⇒ 断言自己必须在 ubuntu 上也能算；(2) 它在 `internal/models`／`internal/tools` 这类包里的话，**包已经在那 5 枚盲包里**，等于让盲包自己当灯；(3) 任何**新包**都要同一笔补 pin（GUARD C），所以现成包⛔ 是零成本落点 |
| `.github/workflows/ci.yml` 新增一步 | 落 `lint`（ubuntu，`:66`）或 `test-windows`（`:517`）都由编排者选；一把纯 `grep` 步在 ubuntu 上**照样**读得到 windows tagged 字节（⛔ 依赖 GOOS） | 能；⚠ 与本票冲突：`AC#4` 逐字禁 `AC#1` 那笔碰 `ci.yml`（票面 `:23`），所以这枚落点**必须单独一笔**；另 `ci.yml:657-744` 那 88 行注释本身就是 census 的合约，加一步＝再养一段注释 |
| `scripts/portable-tests-selftest.sh`（载体） | `lint`／ubuntu，`ci.yml:395` `bash scripts/portable-tests.sh`… 逐字＝`bash scripts/portable-tests-selftest.sh` | ⛔ **落点，但是强制的连带件**：它用 PATH 假 `go`（`:158`／`:250`／`:335` `PATH="$work/bin:$PATH" FAKEGO_PIN_FILE=...`）、有 GUARD D 的种子（`:282-303`，故意**造一枚⛔ 存在的包** `internal/carrier111unclaimed`，理由逐字写在 `:286-290`）、有 census 名册的两枚案例（`:526-566`）。⇒ **甲形若长进 `portable-tests.sh`，就得在载体里补一枚种子，否则它是一枚⛔ 正控制的门**（载体自己的话：`a guard that exists only in the parent is invisible`，`:309-310`） |

### ② 它会不会挪动 GUARD C 的分母

- **⛔ 挪**（就甲形本形而言）。尺＝逐行读 GUARD C 全体（`:542-563`）：它比的是 `want=$pinned`（档旁边那张 pin）与 `got=$resolved`（`go list` 把 `scope` 展开的结果），且只在 `[ -n "$pinned" ]` 的 run 支里跑；census 那一支的 `pinned=''`（`:273`）+ `:367` 逐字"the auditor owns no pin (ticket 251 AC#4)"。⇒ **一枚只读文件字节的名册断言既⛔ 动 scope 清单、也⛔ 动 pin ⇒ GUARD C 的 `want/got` 逐字节不变**。
- **会挪的是那枚"顺手修法"**：如果甲形被实现成"把这 5 枚包塞进某档清单"，那**同一笔**必须补该档 pin——**这句要求的原文在 blob `:417-418`**（两行逐字：`:417` `portable-tests.sh: zero coverage for four days after their tests landed. Pull the package`／`:418` `portable-tests.sh: into a named scope and update that tier's pin in the SAME commit, or`）。⚠ **枚报错给编排者（详见第 ⑥ 节）**：**那句原文此刻长在 GUARD D 的退码文本里（整块 `:407-423`），⛔ 长在 GUARD C 的块里**；GUARD C 自己那句同义要求是 `:556-558`（逐字 "Either restore the scope entry or update the pin / in the SAME commit and say why in the CI log."）。票 301 `AC#1` 与本派单都把它标成"GUARD C 逐字要求"——**文字对、归属错**，按 `:542-563` 去找那两句会找⛔ 到。
- 另有两枚**加档才会撞**的自检（甲形若走"新建一枚 windows-tagged 档"这条路就要一起付）：`:307-325` 要求 `tiers=` 与本文件 `case $mode in` 分支**双向同集**，否则 census **直接拒跑**；`:368-378` 要求每一枚档名在 census 循环里有 pin，否则同样拒跑。⇒ "多一枚档"⛔ 是加一行，是**同时**动 `tiers=`、动分支、动 pin 循环，并补载体案例 15/16（`:526-566`）。

### ③ 它今天会让谁红（具名；尺＝表① 那 5 枚包对 `windows|cli|winsec` 三枚**跑在 windows 上的档**作差，blob `86d89478`）

判据取"每枚含 windows-tagged `_test.go` 的模块包，必须被**一枚会跑在 windows 上的档**认领，否则必须出现在豁免名册里"：

- **在 `86d89478` 上：5 枚包红，具名**＝`internal/agent`（1 文件／1 用例）、`internal/audio`（2／9）、`internal/memory`（1／1）、`internal/models`（3／3）、`internal/tools`（10／27）。合计 **17 枚文件／41 枚顶层用例**。
- **`AC#1` 落地之后（audio 进 `windows)` ＋ `win_pin`）：4 枚包红，具名**＝`internal/agent`／`internal/memory`／`internal/models`／`internal/tools`，合计 **15 枚文件／32 枚用例**。
- 若改用例级判据（"每枚 tagged 用例必须被某档求值"）⇒ 同一批具名包，用例数＝**41（@86d89478）／42（@41329475）／32（AC#1 后）**，另需**豁免 3 枚 ledger 里 `platform=windows` 的行**（`:592` `TestHelperProcess|./internal/proc/|windows|reexec`、`:593` `TestLiveWasapiSmoke|./internal/audio/|windows|fixture`、`:596` `TestSyncRegistryProbeLive|./internal/risk/|windows|fixture`），否则 3 枚"按设计从未求值"的用例会被算成洞。
- **⛔ 红的退化形**（必须写进代价表，⛔ 让下一个读的人以为它买了什么）：如果名册断言只要求"`tools/tagged-roster.txt` 必须与磁盘上的名册**同步**"，那么今天 **0 枚红**（新增用例不登记就红＝只在**未来那笔**才咬；今天全绿），它买到的是"下一个加 tagged 用例的人被迫登记"，⛔ 是"这 41 枚被求值"。⇒ **甲形有两枚子形：登记钉（今天 0 红）／认领断言（今天 5 红、AC#1 后 4 红）**，代价差一整枚射程，本格⛔ 替编排者选。

---

## 表③ 乙形（GUARD D 从包级扩到用例级）的最小诚实版

### ① census 那一跑今天用的平台（尺＝逐行读 `ci.yml` blob 的 job 边界与 env）

- `test-windows:` 键＝`ci.yml:516`；`runs-on: windows-latest`＝`:517`；`env:`＝`:518-519`，**只有一枚 `WISP_ENV: test`**；`steps` 起 `:520`；`- uses: actions/checkout@v4`＝`:521`。
- census 那一步：名字＝`:657`，`shell: bash`＝`:742`，`if: ${{ !cancelled() }}`＝`:743`，`run: bash scripts/portable-tests.sh --scope=census`＝`:744`（**GUARD D 唯一调用点**；`:658` 那两行注释自己就这么声明）。
- 全 `ci.yml` **没有任何非注释的 `GOOS=`**（尺见表② 第一行）⇒ `portable-tests.sh:280` 的 `goos=$(go env GOOS)` 在这一发上＝**windows**。
- ⇒ **一句话留给读者**：GUARD D 今天**已经**跑在一枚能看见 windows tagged 的平台上了；它的盲区⛔ 是 GOOS，**是"认领"这枚关系只有包级、且⛔ 携带"该档跑在哪台 OS 上"这枚输入**（尺＝`grep -nE 'ubuntu|windows-latest|runs-on|test-core|test-windows' scripts/portable-tests.sh` ⇒ **12 处命中全部落在注释与 ledger 理由串里**（代码行 0 处；`core→test-core(ubuntu)` 写在 `:104`、`windows→test-windows` 写在 `:128`，`cli`／`winsec` 两档的 runner 连注释里都⛔ 成文）⇒ **乙形的第一笔硬代价＝把 tier→runner-OS 从散文变成数据**）。

### ② 在那个平台（GOOS=windows）上，今天有多少枚 windows-tagged 顶层用例落进"没人认领"

这一格**取决于"认领"怎么定义**，两枚读数都交（⛔ 只交一枚会让人以为门扩了射程其实没扩）：

- **定义 A＝沿用今天的认领关系**（`:380` "任一 pin 里有这枚 import path"），只把粒度下到用例 ⇒ **0 枚包、0 枚用例红**：表① 那 12 枚 tagged 包**枚枚都在某张 pin 里**（尺＝逐枚 `grep -qxF` 对四张 pin）。⇒ **"只扩粒度、⛔ 动语义"的乙形是一枚恒绿的门**，买到 0 枚可见性。
- **定义 B＝认领必须来自"会跑在该用例 GOOS 上的档"** ⇒ 落进"没人认领"的＝**5 枚包／17 枚文件／41 枚用例**（@`86d89478`）、**5／18／42**（@`41329475`，audio 多一枚）、**4／15／32**（`AC#1` 落地后）。具名见表① 与表②③。
- 需要的输入清单（乙形⛔ 跑 `go test` 也能凑齐的前四枚）：(1) tagged 文件名册（字节尺，表① 那把）；(2) 每枚 tagged 文件里的顶层用例名（`^func Test` 尺）；(3) **tier→runner-OS 映射**（今天只⛔ 是注释，见 ①）；(4) **ledger 名册**（3 枚 `platform=windows` 行必须被读成"已登记豁免"，否则误报）；(5) 若还想买"用例真的编得进二进制"这层，就要 `go test -list`——**census 那一支今天逐字声明⛔ 跑它**（`ci.yml:700` 注释："runs `go env GOOS` and `go list` only, never `go test`, never a build"），而载体那枚假 `go` 要跟着扩（`FAKEGO_LIST_RC`／`FAKEGO_LIST_DROP` 之类旗子已存在于 `:403`／`:415`，`-list` 那一支⛔ 读到它支持，⇒ 本格⛔ 裁）。

### ③ 与既有裁决撞不撞车（三处具名）

1. **GUARD A 那把尺**（`scripts/portable-tests.sh`，注释块起 `:565`、判红块 `:570-583`）：`:565-566` 逐字＝"GUARD A: the loud empty denominator. `.TestGoFiles`+`.XTestGoFiles` are the / files that COMPILE INTO THIS PLATFORM's test binary"，`:570-571` 的 `denom=$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{else}}EMPTY ...` ⇒ **它判的是"档里有一枚包在本平台零分母"**。撞车形状＝如果有人把"让 4 枚盲包开口"实现成**把它们塞进 `core)`（ubuntu 档）**，只有当该包在本平台零分母时才红；现量**不会**红：那 5 枚盲包的**无 tag 文件枚数**＝agent 40／audio 5／memory 10／models 14／tools 43（尺＝blob `86d89478` 下 `git ls-tree -r --name-only <锚> <dir>/ | grep -c '_test\.go$'` 减 tagged 枚数）⇒ 每枚都⛔ 是零分母。**真正的撞车⛔ 在这格**：乙形若把"某档只跑 windows"变成分档输入，`core` 那档在 ubuntu 上会对这些包**永久⛔ 求值**，而 GUARD A 的语义（"entry 必须在本平台有分母"）⛔ 管这件事 ⇒ 两枚门会对同一件事给出相反的颜色读数（A 绿、乙红）。本格⛔ 裁该不该调和。
2. **ledger 的 `fixture`／`reexec` 两类理由**（`:585-601`，11 行＝6 fixture／3 opt-in／2 reexec；platform 分布 4 any／4 linux／3 windows）：`reexec` 那两行（`:591` memory、`:592` proc）的逐字理由都写着"child-process entry point ... **not a test of its own**"；`fixture` 那枚 windows 行（`:593` audio）的理由写着"hosted runner has no audio endpoint"。⇒ 用例级 GUARD D **必须把 ledger 读成豁免名册**，否则今天多红 **3 枚**（具名：`TestHelperProcess`／`TestLiveWasapiSmoke`／`TestSyncRegistryProbeLive`）；而更硬的一处＝编排者 13:3x 自己已裁"包进档后第 9 枚走的是**被 `-skip` 排除＝从未求值**、⛔ 是自己 `t.Skip` 出 `--- SKIP`"（`:680`／`:686`／`:688`／`:694`）⇒ **"从未求值"在本仓同时是（一）件要被乙形追捕的事、（二）一件 ledger 明文批准的状态**；乙形不新造第三态就⛔ 自洽。本格⛔ 裁那枚区分怎么写。
3. **票 111 那句"CI 只测 33 个包里的 20 个"**（`portable-tests.sh:413` 逐字："this row is the ticket 111 field itself ('CI 只测 33 个包里的 20 个')"，另 `:292` 同一句的半行注释与 `:406` 那行 totals 输出）：那句的**单位是包**。乙形下到用例级后，`--scope=census` 打印的那张表与那枚 totals 行（`:406` `packages=$n with-zero-compiled-tests=$empty claimed-by-no-scope=$noscope unclaimed-with-tests=$guardd`）要么改单位、要么成**混单位读数**；而 `:406` 那行形状是被载体逐字判的（`:531` `has 'census totals: packages='`）⇒ **改它＝同时改载体案例 15/16（`:526-566`）**。`ci.yml:688` 那句"GOOS=windows resolves all 35"同属包单位，也要一起记账。⚠ 35 那枚数在本件里的对应现量＝**36** 枚含 `.go` 的非隐藏目录减 `cmd/wisp/testdata/esclistener` 一枚＝35（见「仅读码推的」第 1 条，**这⛔ 是现量、是推断**）。

**乙形另有一枚⛔ 在任何票面上的代价（本格新量，具名）**：census 的宇宙由 `go list ./...`（`:336`）给出，而 `.scratch/wisp/probes/**` 那 10 枚目录／15 枚 tagged 文件／31 枚用例**⛔** 在这个宇宙里。⇒ 谁要是把乙形实现成"自己走目录树"（grep 法），名册会**立刻多出 15 枚文件**、并把别人的工单探针变成 CI 分母。**两形共用这一枚射程陷阱**：断言的射程目录必须逐字写死，⛔ 交给"当前目录下的 `*_test.go`"。

---

## 仅读码推的（⛔ 盘上现量，越读越⛔ 定）

1. **`go list ./...` 到⛔ 了 `.scratch/**` 与 `**/testdata/**`**：依据＝Go 对隐藏目录／`testdata` 的模式排除语义（读码），旁证＝`ci.yml:688` 记的"GOOS=windows resolves all 35"与我逐枚数的 36 枚非隐藏含 `.go` 目录（`git ls-tree -r --name-only 86d89478 | grep '\.go$' | grep -v '^\.scratch/' | grep -v '^scripts/spike/' | grep -v '^tools/d22scan/' | grep -v '^tools/mockllm/' | xargs -n1 dirname | sort -u` ＝ 36，其中 `cmd/wisp/testdata/esclistener` 是 testdata）。⇒ **⛔ 跑过 `go list`**（本票⛔ 编译面），所以"35＝宇宙"是**推断**。
2. **census 今天在 `86d89478`／GOOS=windows 上会打印 `unclaimed-with-tests=0`**：推自"表① 的 12 枚包枚枚在 pin 里"（盘上现量的名册与 pin）＋ `:382-399` 只在 `where` 为空时才计数。**⛔ 实跑过 `--scope=census`**（要跑＝归编排者那一发；另 `:336-351` 还有一枚"`go list` 退非 0 就整支拒跑"的前置，我⛔ 能⛔ 跑就证）。
3. **`internal/agent`／`memory`／`models`／`tools` 的 tagged 用例在 CI 上"从未被求值"**：文件字节层⛔ 是 windows tagged（现量），"core 档跑在 ubuntu"是现量（`ci.yml:402-403`／`:470`），"ubuntu 上 tagged 文件⛔ 进分母"⛔ 是**读码**＋`:565-566` 那句自我声明；⛔ 枚⛔ 跑过 ubuntu 那发的名册作差。
4. **载体那枚假 `go` 支持⛔ 支持 `go test -list`**：我只读到 `FAKEGO_LIST_RC`／`FAKEGO_LIST_DROP`（`:403`／`:415`）**读⛔ 到**它们判的是 `go list` 还是 `go test -list`；⇒ 乙形若要 `-list` 输入，代价枚数⛔ 知。
5. **表② 落点行为"纯 grep 步在 ubuntu 上也读得到 windows tagged"⛔ 实测**：推理自"文件字节⛔ 依赖 GOOS"；真跑＝归编排者。

## 盘上现量（每一行都读到那一行才敢说）

- `portable-tests.sh`（blob `86d89478`＝`41329475`，md5 相同，767 行）：`:155-280` 四张 pin 与四张 scope 清单逐枚读过（core 28 行 pin／windows 10 行、cli 1 行、winsec 1 行；`tiers='core windows cli winsec census'` 在 `:220`；`winsec_scope` 在 `:228-229`）；`:280` `goos=$(go env GOOS)`；`:286-424` census 全体（`:307-325` 名册自比、`:336` `go list ./...`、`:364` counts、`:380` 认领判据、`:382-399` NO-SCOPE 分支、`:406` totals、`:407-423` GUARD D 文本、`:413` 票 111 那句、`:417-418` "同一笔"那两句）；`:542-563` GUARD C；`:565-583` GUARD A；`:585-601` ledger 11 行；`:616`／`:651` `go test -list`；`:639-641` 平台过滤；`:666-674` stale；`:677-688` skip_pattern；`:694` 实调只喂 `"${scope[@]}"`。
- `ci.yml`（blob `86d89478`，已落盘 `probes/301/a2/ci-blob-verbatim.md`）：job 键 `:65 lint`／`:402 test-core`／`:516 test-windows`／`:801 slo-smoke`／`:862 slo-full`／`:936 lint-frontend`；`runs-on` 逐行 `:66`／`:403`＝ubuntu-latest、`:517`／`:802`＝windows-latest、`:863`＝`[self-hosted, wisp-slo]`；`env:` 块 `:404-405`／`:518-519`／`:803-804` 均只有 `WISP_ENV: test`；调用点 `:82`／`:395`／`:445`／`:470`／`:568`／`:604`／`:744`／`:780`／`:796`。
- `scripts/wisp-cli-tests.sh:113` `bash "$portable" --scope=cli`；`scripts/winsec-tests.sh`（`:39` 引 portable、`:88` 自报平台、`:158-160` GUARD 3 判子步⛔ 审计了那枚显式路径）。
- `internal/models/assembly_reachability_121_test.go` blob `:60-125` 逐行读到（per-GOOS 硬判据 `:64-69`；`-e` 的取舍 `:71-91`；`modulePrefix` `:110`；`capabilityPackages` 两枚 `:115-118`；`graphGOOS` 三枚 `:120-123`）。
- 表① 全部枚数（文件／用例／包）＝blob 层逐枚跑出来的，两枚锚各一发，差集已具名。
- 表① 注：`internal/audio` 在 `41329475` 上 8 枚 `_test.go`＝3 tagged（`capturelevel_windows_test.go` 1 用例／`hotplug_test.go` 8／`parse_wave_format_300_windows_test.go` 1 枚用例 `TestParseWaveFormatSubFormatOffset300`）；该新用例体内 `grep -nE 'WISP_LIVE_MIC|testing\.Short|t\.Skip|exec\.Command'`＝**0 命中**（读到那一行为止），**但⛔ 据此判它的前置条件**（本票 `A807`③ 那条定式管着这一读）。

---

## 归编排者（我⛔ 能裁／⛔ 能证的格，具名）

1. **`AC#3` 落⛔ 落、甲⛔ 还是乙⛔ 还是都⛔**——派单逐字写着⛔ 我下结论；上表只交形状与枚数。
2. **两枚"非跑不可才能裁"的格**：
   - census 在 windows 上的**真读数**（`unclaimed-with-tests` 到底几枚、totals 行原文）＝缺 `bash scripts/portable-tests.sh --scope=census` 一发（GOOS=windows 那台上）。我⛔ 跑（派单⛔ 编译面）。
   - `--scope=core` 在 ubuntu 上会不会因"把盲包塞进 windows 档"而变颜色＝缺 `--scope=core` 与 `--scope=windows` 各一发的前后名册作差（同一枚包⛔ 在飞时那发会吞 `300-r1` 的面，票面 `AC#1` 排程已钉住次序）。
   - `internal/audio/parse_wave_format_300_windows_test.go`（13:4x 才进 HEAD）算⛔ 算 `AC#0` 那张已裁名册的第 10 枚用例。
3. **甲形两枚子形（登记钉／认领断言）选一**：今天 0 红与今天 5 红买的是两件事，⛔ 我替裁。
4. **ledger 豁免语义**（乙形要不要把 `fixture`／`reexec` 读成豁免、豁免⛔ 求值的用例算几枚）＝与"SKIP⛔ 是 pass"那条既有裁决正面相遇，⛔ 我调。
5. **`ci.yml` 一步 vs `portable-tests.sh` census 一支**：前者要单独一笔（`AC#4` 禁 `AC#1` 那笔碰 yml），后者要连带改载体种子；两笔次序＝编排者的排程。

## 必答：我⛔ 同意的派单预设（四处，全部带尺）

1. **`:417-418` 的归属**：派单与票面都说那是"**GUARD C 逐字要求**"。**盘上那两句长在 GUARD D 的退码文本里**（整块 `:407-423`，`:409` 那句 `GUARD D - $guardd package(s)...` 就是块首）。GUARD C 自己的同义句是 `:556-558`。文字逐字对、**块名⛔ 对**——照派单去 `:542-563` 找那两句会找⛔ 到，然后很可能有人"顺手补一句"进 GUARD C＝重复立法。
2. **"认领关系是现成的，只⛔ 差粒度"这枚预设⛔ 成立**。派单表③② 的写法（"有多少枚 windows-tagged 用例落进'没人认领'"）暗含"照 census 的 unclaimed 那把尺问一遍就行"。盘上：census 的 `unclaimed`（`:382-399`）判⛔ 的是"任一 pin⛔ 提这枚包"，而**表① 那 12 枚 tagged 包枚枚在 pin 里**⇒ 按现定义**用例级 GUARD D 今天 0 红**。盲区的成因⛔ 是"没档认领"，**是"认领它的档跑在另一枚 OS 上"**——而那枚信息（tier→runner-OS）在 `portable-tests.sh` 里**只有散文（`:104`／`:128`），⛔ 有数据**（12 处命中全在注释与 ledger 串里，代码 0 处）。⇒ 乙形的真实射程⛔ 是"包级→用例级"，**是"包级→用例级＋档级 OS 语义"**，比派单给的框架大一枚。
3. **"表① 的产出之一是'哪些包完全不在任何档里'（GUARD D 已经看见的那一类）"⛔ 有产出**：读数＝**0 枚**。本格真正要收的账⛔ 在那一列，**在"只被 ubuntu 档认领"那一列**（5 枚→AC#1 后 4 枚）。如果按派单那一列去找，会带着一枚空表回来并写"没有盲区"。
4. **⛔ 定"tagged 名册＝`head -3` 就够，⛔ 担心第 3 行以后"**——盘上有 5 枚第 3 行以后出现 `go:build windows` 的文件（逐枚读到前 6 行：全是散文提及，其中 `internal/config/c26_seam_posix_125_test.go:1` 其实⛔ 是 `//go:build !windows`）。⇒ `head -3` 是⛔⛔ 两向都对；但**有人把尺"改进"成全文 grep 就会虚增 5 枚**，而那 5 枚里有 1 枚方向恰好相反（`!windows`）。这一条是替派单**加分**，⛔ 顶人。

## Progress log

- [2026-10-10 13:36:38 +0800] agent=301-a2 did=起手锚 blob 层落盘并单独一笔 commit（HEAD 86d89478，porcelain 2 行全在 internal/audio） next=表①
- [2026-10-10 13:48:16 +0800] agent=301-a2 did=交回锚（HEAD 41329475，porcelain 0 行）＋同一把尺第二发（119→120 枚，差集 1 枚＝audio 那枚 300-r1 新件；portable-tests.sh md5 两号相同）；三张表＋两态＋归编排者＋必答写成 next=编排者裁 AC#3；⛔ 零 push
