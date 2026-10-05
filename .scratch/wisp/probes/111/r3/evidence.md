# 票 111 最后一跳 —— 腿 `111-r3`：把 GUARD D 接进 CI + 给这把闸门配回归覆盖

本文件是写码腿 `111-r3` 的证据件。**只登记本腿自己现量的读数**；行号一律自取，
前腿结论只引 `probes/111/r2/evidence.md` 的 §3／§5.3（派单点名许可的那两节），不重做它的验证。

## §0 起手锚（写满于本腿第一次 commit）

| 项 | 读数 | 时刻（+0800） |
|---|---|---|
| `git rev-parse HEAD` | `862556b38090306e167e732d721a7acfecd3ba8b` | 2026-10-05 22:15 |
| `git status --porcelain -- .github scripts cmd internal` | **0 行** | 22:15 |
| `git status --porcelain` 全仓 | **730 行**（收尾复量 729 行，差的是别人那两小时里落盘的在飞件；本腿没碰他人写面） | 22:15 / 22:31 |
| `grep -c census .github/workflows/ci.yml`（**改前**） | **0** | 22:16 |
| `grep -cE 'GUARD D\|unclaimed\|UNCLAIMED' scripts/portable-tests-selftest.sh`（**改前**） | **0** | 22:16 |
| `grep -c portable-tests-selftest .github/workflows/ci.yml` | **0**（载具自己不在 CI 里，r2 §5.3 那半条洞，本腿不修，见 §5） | 22:34 |
| `grep -cE 'winlive\|-tags\|GOFLAGS' .github/workflows/ci.yml`（**改前**） | **0**（⛔ 收尾必须仍是 0，见 §4） | 22:34 |
| `md5sum scripts/portable-tests.sh`（工作树） | `328eb3ead0545736a33e2c131687d60d` | 22:22 |
| `git cat-file blob HEAD:scripts/portable-tests.sh \| md5sum` | `328eb3ead0545736a33e2c131687d60d`（**同值**） | 22:22 |
| `git rev-parse HEAD:scripts/portable-tests.sh` | `2ff02dd7476c8f623a4489932d0d0ff95480e667` | 22:22 |
| `git rev-parse HEAD:.github/workflows/ci.yml` | `c5a1b1ab6467a2f221c7f70066688f24d8ab9dea` | 22:22 |
| `git rev-parse HEAD:scripts/portable-tests-selftest.sh` | `6c5a5a36b625f9897f37bccfb3cd6f1978840a6f` | 22:22 |

- 靶件同一性（所以 r2 §3 的两发正控对本腿 HEAD 仍成立，不必重跑）：
  `15d8b60e..HEAD` 在 `-- scripts internal cmd .github` 上 `git diff --numstat` = **空**（0 行）⇒
  `scripts/portable-tests.sh` 的 blob 仍是 r2 交付的那枚 `2ff02dd7`，`internal`／`cmd` 下
  `git diff --name-status 15d8b60e..HEAD` = **0 枚**文件 ⇒ 包名册与 GUARD D 的字节本腿一字未变。
- 起手 `bash -n`：`scripts/portable-tests-selftest.sh` rc=**0**、`scripts/portable-tests.sh` rc=**0**（22:27:03）。
- 起手载具全跑（**改前基线**）：`bash scripts/portable-tests-selftest.sh all` 22:27:03→22:29:07 rc=**0**，
  逐字末三行：`portable-tests-selftest.sh: 27 case(s) ran, 0 assertion(s) failed` /
  `carrier scratch kept at /tmp/tmp.tHT8tbamqf...` / `GREEN - every seeded anomaly was refused, and the clean scope passed.`
  留盘 `/tmp/111r3-baseline-selftest.log`（临时件只建不删）。
  ⇒ **27 枚**是 r2 §5.3 说的"27 枚载具用例里没一枚会红"的那 27，本腿复量复认。
- 起手 YAML 解析（python3 `yaml.safe_load`）：`ci.yml` 改前可解析（22:31:38），`test-windows` 8 步。

## §1 格 1 —— census 接进 CI（`.github/workflows/ci.yml`）

**做了什么**：`test-windows`（`runs-on: windows-latest`，`:420`）新增一步，
位置＝**`:509` 起、`run:` 行在 `:596`**，夹在
`cmd/wisp CLI tests`（`:487`／`run:` `:507`）与
`Portable windows tests`（改名前 `:509`／`run:` `:543`，现 `:598`／`:632`）之间。
命令一字：`bash scripts/portable-tests.sh --scope=census`，带 `if: ${{ !cancelled() }}`。

现量（改后）：

| 尺 | 读数 |
|---|---|
| `grep -c census .github/workflows/ci.yml` | 0 → **13** |
| `grep -cE 'winlive\|-tags\|GOFLAGS'` | 0 → **0**（⛔ 禁面守住，见 §4 那发更正） |
| `grep -c portable-tests-selftest` | 0 → **0**（载具名一字不写，理由见下"为什么注释里不点载具文件名"） |
| 逐步 `continue-on-error` 键（yaml 遍历全部 6 job） | **0**（字符串出现 13→14 次，多的那一枚是本腿注释里那句"NO step here carries `continue-on-error`"，非步骤键；HEAD 基线同为注释 13 枚） |
| `git diff --numstat .github/workflows/ci.yml` | **89 增 / 0 删** |
| 文件行数 | 766 → **855** |
| 改前后相邻步的行号（工作树现量） | `test-windows:` `:419`、`runs-on: windows-latest` `:420`、ACL 门禁 `run:` `:471`、`cmd/wisp CLI tests` `run:` `:507` ＝ **未漂**；新步名 `:509`／`run:` `:596`；`Portable windows tests` 名 `:598`（原 `:509`）／`run:` `:632`（原 `:543`） |
| 逐步名册（yaml 解析，9 步） | 1 checkout / 2 setup-go / 3 ACL / 4 Cache third_party / 5 cgo smoke / 6 CLI tests / **7 Package coverage census (ticket 111 AC#1 + GUARD D)** / 8 Portable windows tests / 9 PathResolver junction；非 setup 的 7 枚全部 `if: ${{ !cancelled() }}` |

**为什么在这一步、这个位置（派单点名要写明的那句）**：

1. **腿＝windows**。census 在 GOOS=linux 必红且是"拒绝"不是"发现"：`go list ./...` rc=1
   （`sherpa-onnx-go-linux` 排除全部文件）⇒ 命中 `scripts/portable-tests.sh:337-351` 的拒绝分支 ⇒ `exit 1`。
   本腿没有重跑这条（零 `go test/build/vet`，`go env` 也一次没手打），出处＝r2 §1.2 的现量
   （2026-10-05 20:36:55，stdout 34 行／stderr 264 字节）＋脚本自己 `:326-332` 的注释复认。
   ⇒ ubuntu 腿接它＝一枚永久红且原因不是产品 bug，**只能挂 windows 腿**（r2 §5.3 硬事实 ②，本腿复认机制）。
2. **步骤不受上面红的影响**：本腿新步与 `test-windows` 其余 6 枚非 setup 步同形带 `if: ${{ !cancelled() }}`
   （yaml 逐步核过），所以"红步骤之后的步骤根本不执行"那条病（票 111 AC#6 的
   run `35599458439`：step4 failure ⇒ step5–8 全 skipped）在这里不成立 —— 位置是关于**顺序**，不是躲红。
3. **顺序的取舍**：census 只跑 `go env GOOS` + `go list`，**不跑测试、不建任何东西、不需要 third_party／DLL**，
   所以它放在两条它不消费其产物的步（cache／cgo／CLI）之后；放在 `Portable windows tests` **之前**
   ＝它审计的那份名册先答，windows 腿一红，先看到的是"哪枚包失去了档"而不是分母自己缩短。
   ⛔ **没有**插到 `Windows ACL sealing gate` 前面：票 110 AC#5 把 ACL 钉死为 job 内第一步是有出处的
   （`ci.yml:432-435`），census 红不该推迟 ACL 证据。
4. 无步骤被移动／编辑／删除（`git diff --numstat` 删除列 **0**）。

**为什么注释里不点载具文件名**：`grep -c portable-tests-selftest ci.yml` = 0 是 r2 §5.3 用来
判"载具不在 CI 里"的那把尺，本腿若把它写进 ci.yml 就把这把尺洗成 1 了；同理第一版注释写了
`-tags winlive`／`GOFLAGS`／`-skip` 三个词，把 `grep -cE 'winlive|-tags|GOFLAGS'` 从 0 抬到 1×3
—— 那三枚现量就是 §4 里那次更正的来源（措辞改为不含字面 token 的等价陈述，事实与出处一字不丢）。

## §2 格 2 —— GUARD D 的回归覆盖（`scripts/portable-tests-selftest.sh`）

改前三把尺（本腿现量）：
- `grep -cE 'GUARD D|unclaimed|UNCLAIMED'` = **0** ⇒ 零覆盖（r2 §5.3 点名的洞）。
- `scripts/testdata/portable-tests/go:50-56`：`go list -f` 那一支对**任何** fmt 都只认
  `FAKEGO_GUARD_A_SEED`，否则 `exit 0` 打印**空**⇒ census 的逐包探针
  `go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' <pkg>`（`portable-tests.sh:364`）**结构上无法被种**，
  所以"GUARD D 零覆盖"不只是"没人写用例"，是"写不了用例"。
- 载具 invocation 现量：`grep -cE '^\s+(run|run_with|run_chain) '` = **23**，`grep -c '^if selected '` = **22**，
  `ran` 自增加 3 枚手写发射 ⇒ 基线报 **27** 次（上面 §0 那发全跑）。

**改了什么**：
1. `scripts/testdata/portable-tests/go` —— `list -f` 那一支先认 census 的 fmt，按包答计数：
   新增 `FAKEGO_CENSUS_COUNTS`（`<importpath> <t/x>` 表）与 `FAKEGO_CENSUS_DEFAULT`（默认 `0/0`）；
   非该 fmt 仍走 GUARD A 的老形状（**未改语义**），无包参数时 `exit 99` 响亮拒绝而不是编一个空答案。
   契约注释同步进文件头。
2. `scripts/portable-tests-selftest.sh` —— 种子 + 两枚用例：
   - 种子 `unclaimed_path=github.com/CarlosShao/wisp/internal/carrier111unclaimed`（**发明的路径**，
     不借真包：借真包会让 GUARD C 在档位跑上红，本用例就变成在判别的东西）。
     universe＝`core_pin (27) + 该枚 = 28 行，claimed by no tier`，与 counts 表一起写在文件顶部
     （单枚用例也要能单独跑，`set -u` 不许它撞未定义名）。
   - **case 23 `census-unclaimed-package-goes-red`**（一次选中跑 3 发，同一天空三种形状）：
     (a) 计数 `1/0` ⇒ rc=1 + `GUARD D - 1 package(s) compile a test file for GOOS=` +
     `NO named scope claims them, so no CI step runs them` + `<-UNCLAIMED-HAS-TESTS` +
     `<path>  tests-compiled-for-<goos>=1/0` + `census totals: packages=28 .*unclaimed-with-tests=1`；
     (b) 计数 `?/?`（读不到）⇒ 仍 rc=1（default-deny 那条语义被钉住）；
     (c) **同一天空、计数回默认 `0/0`** ⇒ rc=0、`GUARD D -` 不出现、`<-NO-TESTS` 在、
     `unclaimed-with-tests=0` ⇒ 这一发是防"因"的：证明 (a)(b) 红的是**种下去的计数**，
     不是"这枚包没人认领"这件事本身（否则 GUARD D 死了 (a) 也会红，本腿的接入主张就不可证伪）。
   - **case 24 `guard-d-slice-is-what-bites`**（★牙）：把 `portable-tests.sh:393` 那枚
     `            0/0) ;;` 换成 `            *) ;;`（**shadow 副本**，工作树不碰），
     同一 universe 同一 counts 表 ⇒ **mutant rc=0 且 `GUARD D -` 消失、`unclaimed-with-tests=0`**，
     而真字节同一天空 **rc=1 点名**；`touched -ne 1` 时用例自己 FAIL 并打印
     "the mutant seed took N line(s) ... re-read the script instead of trusting the carrier"，
     不交 vacuous green（照 case 16/18/22 的既有形状）。
   - 锚唯一性现量：`grep -c '^            0/0) ;;$' scripts/portable-tests.sh` = **1**，
     `grep -c 'guardd=\$((guardd + 1))'` = **1**（在 `:396`）。

**读数**（全部带时刻，台件留盘 `/tmp/111r3-case23.log`／`case24.log`／`full-after-edit.log`）：
- 23:04:17→23:04:2x `bash scripts/portable-tests-selftest.sh census-unclaimed-package-goes-red` rc=**0**，
  `3 case(s) ran, 0 assertion(s) failed`（三发逐枚 ok 见 §2 表下）。
- 23:05:00 `... guard-d-slice-is-what-bites` rc=**0**，`2 case(s) ran, 0 assertion(s) failed`。
- 23:05:40→23:08:15 **全跑** rc=**0**，`32 case(s) ran, 0 assertion(s) failed`，
  `GREEN - every seeded anomaly was refused, and the clean scope passed.` ⇒ 27 → **32**（+5 次发射＝23 的 3 发 + 24 的 2 发）。
- 改后 `grep -cE 'GUARD D|unclaimed|UNCLAIMED' scripts/portable-tests-selftest.sh` = **待 §3 填**。

## §3 ★工作树"改坏必红"那一窗（派单格 2 的第二半，边界具名）

窗口＝**同一条 bash 调用内**完成：md5 起手 → 改坏 → `grep -n` 打印被改后整行 → 载具单发 ×2 →
`git cat-file blob HEAD:scripts/portable-tests.sh > scripts/portable-tests.sh` 还原 → md5 双读。
窗口外任何时刻工作树的 `portable-tests.sh` 都等于 HEAD。读数与红句逐字见 §3.1／§3.2。

- 窗口边界（起止时刻）：`______`（§3 填）
- 净变化：`md5sum` 与 `git cat-file blob HEAD:... | md5sum` 双读同值 = `______`
- `git status --porcelain -- scripts` 还原后：只应剩本腿的两枚文件（selftest + testdata go），`portable-tests.sh` **不出现**。

## §4 门禁四把尺（带时刻＋HEAD）

| 尺 | 读数 | 时刻 |
|---|---|---|
| `bash -n scripts/portable-tests-selftest.sh` | rc=**0** | 23:04:1x |
| `bash -n scripts/testdata/portable-tests/go` | rc=**0** | 22:5x（shim 单测） |
| `bash -n scripts/portable-tests.sh` | rc=**0**（本腿净变化 0） | 22:27:03 |
| `sh scripts/portable-tests-selftest.sh`（全跑） | 见 §5 | |
| `sh scripts/d22scan.sh` | 待填 | |
| `sh scripts/check-path-length-budget.sh --with-self-test` | 待填 | |
| ci.yml 禁面 `grep -cE 'winlive\|-tags\|GOFLAGS'` | **0** | 23:0x |

⚠ **本腿自查并就地更正一处**：第一版 ci.yml 注释为了让"我这步没加什么"可核，
把 `-tags winlive`／`GOFLAGS`／`-skip` 三个词写成了字面量 ⇒ `grep -cE 'winlive|-tags|GOFLAGS' .github/workflows/ci.yml`
当场从 0 变 3 枚命中（各 1）。那把尺是本票 AC#1／票 111 r2 §5.1 用来判"winlive 零接入"的仪器，
写注释的人不许把它洗软。已改写为不含这些 token 的等价句子（`no build tag of any kind is passed here,
no environment-level tag override, no skip flag by name`），复量回到 **0 / 0 / 0**；
事实与出处（11 枚 winlive 文件／5 枚自带 13 处 `t.Skip`／`SKIP-LOUD: this host reports no ball window`
@ `cmd/wisp/resident_hotkey_live_258_windows_test.go:55`／`tools/d22scan/runtests.sh:98-102` 判 SKIP fatal）一字未减。

## §5 判不动的地方（具名归口）

- **载具自己不在 CI**：`grep -c portable-tests-selftest .github/workflows/ci.yml` = 0，本腿**不改**。
  理由具名：把载具接进 CI 是"每 run 起 32 次 `bash portable-tests.sh` 假跑"的另行决定，
  派单写面只给 census 一步，越界即自造契约变更（`AGENTS.md §0` 第 2 句）。⇒ 归编排者。
- **winlive 半边**：owner 未批 ⇒ 本腿零接入（见 §4 的更正与 ci.yml 注释"WHAT THIS STEP DOES NOT ADD"）。
- **GOOS=linux 的名册**：census 在 ubuntu 是"拒绝"不是"发现"（§1 理由 1）⇒ 归票 98 的地界。
- **真实 census 全量读数（packages=35/unclaimed-with-tests=0）**：本腿**没重跑**，
  依据＝§0 的靶件同一性（blob `2ff02dd7` 与 r2 同一枚，`internal`/`cmd` 在 `15d8b60e..HEAD` 零枚文件变化）
  ⇒ r2 §3.1 的读数对本腿 HEAD 直接成立；且本腿写面里没有一枚 `.go` 文件，名册不可能被本腿改动。
- **票面 AC 框**：一枚未翻（`.scratch/wisp/issues/111-*.md` 的 `- [ ]` 全数不变），见 §6 自证。
- **步级读数**（那一步真给过结论的 run id + step 号）：本腿**未 push** ⇒ 这一格交的是"欠"，
  取数照票面 `:49-50`：`.steps[].order` 返回 `null`，步号按数组位置数并同引 step 名。

## §6 交件自证与 commit 链

- 起手 porcelain：scoped（`.github scripts cmd internal`）**0 行**／全仓 **730 行**（22:15）。
- 收尾 porcelain：见下（commit 链之后复量）。
- `git diff --numstat` 逐笔：见各笔。
- commit 链：`______`
