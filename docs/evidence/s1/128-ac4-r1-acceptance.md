# 128 AC#4 r1 — 非实现者裁决表（门禁 / 「本机绿·CI 红」判定 / pin 有没有牙 / 副产物那一发的连带判断）

**本程是谁**：`acceptor-ticket128-ac4-r1`，非实现者验收代理。实现方是另一枚代理（`worker-ticket128-ac4`）。
**本格判据本体**：票面 `.scratch/wisp/issues/128-resolvedatadir-...-memory.md` 的 `AC#4` 那一格（`:25-27`）
+ `AC#2` 裁定段（`:29-39`）。本表只出判语，**`AC#4` 的勾由编排者按本表打**，本程一枚未翻。

> **本文件的提交回显规矩**（写在最前，免得下一位数错）：每裁完一节立刻 commit 一次，
> 只带这一枚路径 ⇒ **任何一节记的「本节 commit 回显」只能由下一节的那枚 commit 落盘**（一列 commit 列不进它自己）。
> 可复算读法：`git log --oneline -- docs/evidence/s1/128-ac4-r1-acceptance.md`。

---

## 0. 我的锚点自量（全部 `git rev-parse` / `git cat-file -t` 现量，不抄派单号）

| 项 | 现量值 | 怎么量的 |
|---|---|---|
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 被验的码 | `c2fa2e98f7ce87312b785a5b67cefe9c0294914c` | `git cat-file -t` = `commit`；`git show --numstat` = **45 增 1 删**，只 `cmd/wisp/dataroot_128_test.go` 一枚路径 |
| 它的父（修前） | `2956897b1a9135b0427821ddf13cf690cc2da4d8` | `git rev-parse c2fa2e9^` |
| 证据件 | `docs/evidence/s1/128-ac4-gates-and-ci-divergence.md`，351 行，四枚 commit `e4a4e9a ef26704 8dac13b 744d79e`（逐枚 `git cat-file -t` = `commit`，`git merge-base --is-ancestor $a HEAD` 全真） | 只当「它声称做了什么」读 |
| 票面进度件 | `c23d825` = **1 增 0 删**，只票面 128 那一枚路径；`AC#4` 那格 `[ ]` 原样在 `:25` | `git show --numstat c23d825` |
| 我读 CI 原文 | `gh run view 35967768017 --log-failed` **rc=0、一次成功、没重试**，421465 字节 -> `D:\tmp\wt-acc128-ac4-r1-ref\gh-selffetch.txt` | 编排者给的 `D:\tmp\t128-ci-logfailed.txt` 我一枚字节都没读 |
| 工具链 | `go version go1.27.1 windows/amd64`；`gofumpt --version` = **`v0.12.0 (go1.27.1)`** | 盘上现跑（派单写的 v0.7.0 是过期值，本表按 v0.12.0 量） |
| 宿主 | windows / Git Bash；`GOOS=linux go vet` 只编译不执行，本表不当它是代跑 | — |

**快照形状**（不读脏工作树）：`git archive <sha> | tar -x -C /d/tmp/wt-acc128-ac4-r1/snap-<name>`，
再把仓外那三枚 sherpa DLL 复制进 `snap-*/third_party/sherpa-onnx/`（DLL 是**未跟踪件**：`git ls-files third_party | wc -l` 现量 **0** ⇒ `git archive` 里必然没有，这一步是必需的，不是抄工作树）。
读数一律 `export PATH="<snap>/third_party/sherpa-onnx:$PATH"`（MSYS 形；`scripts/wisp-cli-tests.sh` 自己的注释就写着 `pwd -W` 那形是坑）。
每棵快照里留了 `.anchor-sha`（`git rev-parse` 的产物）自证是哪版树。

| 快照 | 树锚点（`.anchor-sha` 现量） | 内容 |
|---|---|---|
| `snap-post` | `c2fa2e98f7ce87312b785a5b67cefe9c0294914c` | 被验版，未改 |
| `snap-pre` | `2956897b1a9135b0427821ddf13cf690cc2da4d8` | 修前那版（`c2fa2e9^`） |
| `snap-mut-i-drop-pin` | 自 `snap-post`，摘掉 `dataroot_128_test.go:261` 那一行 pin | 我自己造 |
| `snap-mut-ii-drop-plant` | 自 `snap-post`，摘掉 `:260` 那一行「先栽 test」 | 我自己造 |
| `snap-mut-iv-both-removed` | 自 `snap-post`，两行都摘 | 我自己造 |
| `snap-post-mut-base` | 自 `snap-post`，`doctor.go:261` 的 `return "", dataDirUnresolved128(env, err)` 换成 `base = "."` | 我自己造（AC#3 那枚 M-1 单点回退） |
| `snap-pre-mut-base` | 自 `snap-pre`，同一处同一改 | 我自己造 |

---

## 1. 被验对象身份核对（含编排者点名要我判的那处行号差）

**1.1 `git show --numstat c2fa2e9` 现量 = `45 1 cmd/wisp/dataroot_128_test.go`**，一枚路径，与简报逐字相符。

**1.2 那 1 枚删除行**：`git show c2fa2e9` 里唯一的 `-` 行是 `-\tfailConfigDir128(t)`（EveryLeg 那枚用例的调用点）。
我自己去新 helper 的函数体里找了它，**在**：

```
$ grep -n "^func \|pinEnvThatAsksTheOS128(t)\|failConfigDir128(t)" <c2fa2e9 版 dataroot_128_test.go>
65:func failConfigDir128(t *testing.T) {           <- helper 本体，函数体未改
257:func pinEnvThatAsksTheOS128(t *testing.T) {    <- 新 helper（简报说 :257 起，逐字对上）
262:        failConfigDir128(t)                     <- 那枚被删的行，搬进来了
323:        pinEnvThatAsksTheOS128(t)                <- 原调用点换成新 helper
```

（262/323 两行前面的空格是原文件的制表符缩进，`cat -A` 量过：`-^IfailConfigDir128(t)$`。）

⇒ **判语：「搬进新 helper」成立，不是「放宽断言」。** 尺与读数：
`git show c2fa2e9 | grep -c '^-'` = **2**（其中一枚是 `--- a/...` 文件头，不是内容），
`git show c2fa2e9 | grep -c '^-[^-]'` = **1** ⇒ 真内容删除行只有一枚，就是上面那句；`rescueMarkers128` 六枚 marker、
`if rc == 0 { t.Fatalf }`、`assertDirEmpty128` 三样**逐字未出现在 diff 里**（没被碰）。

**1.3 静态三处（简报说我 17:4x 已核，我自己重量一遍，不接它的数）**

| 主张 | 我的现量 | 判 |
|---|---|---|
| `.github/workflows/ci.yml:337` 给整个 `test-windows` job 设 `WISP_ENV: test` | `grep -n "WISP_ENV" ci.yml` -> `227 / 337 / 482 / 540`；`334: test-windows:` -> `336: env:` -> `337: WISP_ENV: test`；另三枚分别落在 `224: test-core:` / `479: slo-smoke:` / `537: slo-full:` 之下 | **成立**，四枚行号逐字对上，且都是 **job 级** `env:`（不是 step 级） |
| `doctor.go:256-259` 在 `env=="test"` 时先 `return proc.TestDataDir()` | `grep -n "func resolveDataDir" -A 30` -> `247` 函数头、`256 if env == "test" {`、`257 return proc.TestDataDir(), nil`、`258 }`、`259 base, err := userConfigDir()` | **成立**：`userConfigDir` 的读在 `259`，被 `256-258` 那一支短路 |
| `:139` 那枚显式 `buildinfo.EnvDev` | `138 {entry: "resolveSecretLayout", ...}` / `139 l, err := resolveSecretLayout(buildinfo.EnvDev)` | **成立** |

**1.4 编排者要我裁的那处行号差（报告引 `:131`、盘上 `:139`）—— 我的判定：是「改前号」，不是「数错了」**

尺与读数：
- `git log --oneline -- cmd/wisp/dataroot_128_test.go` 现量只有两枚：`c2fa2e9` 与 `4e5d240` ⇒ 报告 §0-§1 写于 `e4a4e9a`（**早于** `c2fa2e9`），那时那枚文件就是 `c2fa2e9^` 版。
- `git cat-file blob c2fa2e9^:cmd/wisp/dataroot_128_test.go | sed -n '129,133p'` 现量：**`131 l, err := resolveSecretLayout(buildinfo.EnvDev)`** —— 逐字就是那一句。
- `c2fa2e9` 的头一个 hunk `@@ -24,6 +24,14 @@` 在它之上插了 **8 行注释** ⇒ 同一句在 HEAD 上是 `131+8 = 139`。同族另一处也自洽：报告写 `failConfigDir128` 在 `:57-62`，HEAD 上是 `:65-70`（同一个 +8）。

⇒ **判语：报告的 `:131` 在它取数那版树上是对的（正确但已腐坏）；缺陷类型 = 引读数没带锚点版本，属「归因/行号腐坏」这一类，不判它数错。** 按本仓纪律这算**本表的一条建议**：那一格该在原地追加一行「HEAD 上为 `:139`（`c2fa2e9` 之上加了 8 行注释）」而不是改写原句。

**1.5 但我另外量到一处真的数错了，与上面那一处要分开记**

报告 `§1.3` 第 2 条写 `cmd/wisp/doctor.go:255-258`，并把那四行原文贴了出来。
`git cat-file blob e4a4e9a:cmd/wisp/doctor.go | sed -n '253,260p'` 现量：`256` 才是 `if env == "test" {`；
且 `git diff --numstat e4a4e9a HEAD -- cmd/wisp/doctor.go` **空输出**（doctor.go 自报告写下之后一字未动）。
⇒ **`doctor.go:255-258` 在写下的那一刻就是错的（off by one），HEAD 上正确的是 `:256-259`**；
同一枚 off-by-one 也出现在 `e4a4e9a` 的 commit message 正文（「doctor.go:255 提前 return」）。
那一枚 commit message 改不了（已落地），**记在这里当更正**；报告正文那一格本程**未动**（不是我的文件）。
判定影响：**零** —— 它贴出的四行原文是对的、结论是对的，只是行号差一。

### 1.6 §0-§1 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
d9f008d evidence(128 AC#4 r1 §0-§1): 非实现者裁决表开头——锚点自量（被验码 c2fa2e9 的 45增1删、七棵 git archive 仓外快照＋未跟踪 DLL 要手拷）＋身份核对：那枚删除行 cat -A 现量是 -<TAB>failConfigDir128(t)$、新 helper :257 体内 :262 找回来 ⇒「搬进新 helper」成立不是放宽断言；三处静态（ci.yml 四枚 job 级 env 行号、doctor.go 256-259 短路、:139 显式 EnvDev）我自己重量全成立；行号差裁定＝报告引的 :131 是它取数那版的正确号（c2fa2e9 之上加了 8 行注释），缺陷类型是「引读数没带锚点」不是数错；另量出 doctor.go:255-258 那一处是真 off-by-one（doctor.go 自 e4a4e9a 一字未动，正确是 256-259），对结论零影响

docs/evidence/s1/128-ac4-r1-acceptance.md
```

（`git show --name-only HEAD` 只列出这一枚路径；那次 commit 前 `git diff --cached --name-only` 也只有这一枚。）

---

## 2. 主张一「本机绿 / CI 红，根因＝一枚具名变量」—— 我在自己的快照上独立复现了

**主张原文（要判的就是这一句）**：`.github/workflows/ci.yml:337` 给整个 `test-windows` job 设 `WISP_ENV: test`
-> `cmd/wisp/doctor.go:256-258` 在 `env=="test"` 时**先** `return proc.TestDataDir()`、根本不读 `userConfigDir`
-> `dataroot_128_test.go` 注入的那枚 seam 在 runner 上从不被咨询 -> 四条腿红、第五枚（显式传 `EnvDev` 的 `resolveSecretLayout`）绿。

**2.1 先核「我复现的是不是同一段码」——按 blob 号核，不按票号核**

`gh run view 35967768017 --json headSha,conclusion,createdAt,event` 现量：
`headSha=182daed377feeff20f73ea47675c121df5607ca5`（`git cat-file -t` = `commit`）、`conclusion=failure`、
`createdAt=2026-09-24T07:05:41Z`、`event=push`；`git merge-base --is-ancestor 182daed c2fa2e9` -> **yes**。

| 文件 | CI 那版的 blob | 我的 `snap-pre`（=c2fa2e9^） | 我的 `snap-post`（=c2fa2e9） | 判 |
|---|---|---|---|---|
| `cmd/wisp/dataroot_128_test.go` | `686b3fd` | `686b3fd` | `cf40db4` | **与修前逐字节同** |
| `cmd/wisp/doctor.go` | `c8bdccc` | `c8bdccc` | `c8bdccc` | 三版同（本程未动生产码，这一条也顺手证了） |
| `.github/workflows/ci.yml` | `c5a9806` | `c5a9806` | `c5a9806` | 三版同（ci.yml 未被碰，与 §1 的静态读数一致） |

⇒ 下面所有本机读数取的是**runner 跑过的那一枚测试文件的字节**，不是"差不多的另一版"。

**2.2 CI 侧原文（我自己 `gh` 取的那份；编排者给的台件 `D:\tmp\t128-ci-logfailed.txt` 我一枚字节都没读）**

`test-windows` -> step `cmd/wisp CLI tests (...)` 里 CI 自己打的四数行（`portable-tests.sh` 打的，非我数）：

```
=== RUN=101  --- PASS=49  --- FAIL=5  --- SKIP=0
```

那一步的 `--- FAIL` 原文（含缩进子用例）**9 行**，前 5 行是票 128、后 4 行是票 123 那一族
（`TestComposedGateBlocksAWriteForTwoSeconds` + 三枚 `TestTicket101*`）；四红一绿的**绿**那一行也在这一步里：

```
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128                     <- 顶层
    --- FAIL: .../runTextTask   --- FAIL: .../cmdModels   --- FAIL: .../cmdProviders   --- FAIL: .../cmdDoctor
    --- PASS: .../resolveSecretLayout                                              <- 反形：这一枚显式传 EnvDev，不读环境
```

四条红句的断言点与 marker 清单，CI 原文与我的复现**逐字同**（见 2.3 的差分尺）。
另两件事在 CI 原文里直接可见，对本格值钱：
- 红句里**没有一条** `wrote into the start-up directory` -> `assertDirEmpty128` 在 runner 上没响（= §6 那发副产物的原始凭据）；
- 同一步的日志把落点印成了 `dir=C:\Users\runneradmin\AppData\Local\Temp\wisp-test-8104\logs`，
  即 `proc.TestDataDir()` —— **仓外绝对路径**，不是当前目录。

**2.3 我的复现（P1）：`snap-pre` + `WISP_ENV=test`，只跑这一枚用例**

```
cd /d/tmp/wt-acc128-ac4-r1/snap-pre && export PATH="/d/tmp/wt-acc128-ac4-r1/snap-pre/third_party/sherpa-onnx:$PATH" \
  && WISP_ENV=test go test -count=1 -v -run 'TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128' ./cmd/wisp/
rc=1   4 x --- FAIL + 1 x --- PASS（resolveSecretLayout），断言行 dataroot_128_test.go:295
```

**与 CI 日志同形到什么程度（尺，不是形容词）**：把两份文本各自砍掉"每发都会变的字节"
（`(<num>s)` 时长、`time=...Z` 时间戳），只保留 `=== RUN` / `--- PASS|FAIL|SKIP` 行与含 `AC#2 RED` 的断言行，
各得 **16 行**，`diff` -> **rc=0，零输出**（文件 `N-ci.txt` 16 行 / `N-p1.txt` 16 行）。
⇒ **红名、红的深度（顶层 + 4 子）、断言点 `dataroot_128_test.go:295`、六枚 marker 的原文顺序，逐字同形。**
这不是"看起来一样"：16/16 行等字。

**2.4 修后同形（P2）：`snap-post` + `WISP_ENV=test`，同一条命令**

```
rc=0   --- PASS: TestAC2EveryLeg...128 + 5 枚子用例全 PASS（共 6 行结果）
```

**2.5 结论与档位**

| 判据 | 判语 | 档位 |
|---|---|---|
| 「本机绿、CI 红」这件事成不成立 | **成立**。同一条命令、同一批字节，唯一变量是 `WISP_ENV`：未设（P 的前一发／G6）-> 全绿；`=test`（P1/G5）-> 只有这一枚用例红，且红形与 CI 逐字同 | **独立复现** |
| 根因是否具名到那一枚变量 | **成立**。短路点我自己读到（`doctor.go:256-258` 在 `:259` 的 `userConfigDir()` 之前）；反形自己量到（第五枚子用例在两形下都 PASS） | **独立复现**（静态 + 行为两侧都我亲自） |
| 「四红一绿是指纹」这句 | **成立**，CI 与我两边都是 4 红 1 绿，绿的正好是显式传 `EnvDev` 那一枚 | **独立复现** |
| 报告 §1.2「本机也红，但只在把那枚变量搬过来之后」 | **成立**（我的 P1 就在这台机上） | **独立复现** |
| 它的落盘 `D:\tmp\t128-ci-logfailed.txt`（2283 行） | 我**没读**，也不需要用：同一条 `gh` 命令我自己取到 421465 字节、一次成功没重试 | 仅自述（不影响判定） |

### 2.7 §2 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样；标题被终端截了一刀，全名见上一条命令）

```
d3dd1cc evidence(128 AC#4 r1 §2): 主张一「本机绿/CI 红，根因＝ci.yml job 级 WISP_ENV=test」我独立复现——先按 blob 号核身份（CI run 35967768017 的 headSha 182daed 三枚相关文件 blob 与我的 snap-pre 逐字节同：686b3fd/c8bdccc/c5a9806），再在仓外快照上跑：WISP_ENV=test 下 4 红 1 绿、断言点 dataroot_128_test.go:295；把 CI 原文与我的日志各砍成 16 行结构骨架，diff rc=0 零输出＝逐字同形；snap-post 同形 rc=0 全绿；CI 侧 gh 自取 421465 字节一次成功，编排者那份台件一枚字节未读

docs/evidence/s1/128-ac4-r1-acceptance.md
```

（那次 commit 前 `git diff --cached --name-only` 现量只有这一枚路径。）

---

## 3. 「pin 有牙」这句话本身 —— 我自己造了三发，**没有默认它成立**

被裁的主张（报告 §3/§4.2）：`pinEnvThatAsksTheOS128` **先栽 `WISP_ENV=test`、再 pin 回 dev**，
「栽在前 => 摘掉 pin 本机就红（`dataroot_128_test.go:323`）」，因而这枚 pin 在**任何平台**都咬得住。

三发变异都在我自己的 `snap-post` 拷贝上、都**只动测试文件**（生产码一枚字节未动）：

| 发 | 我摘掉的行（HEAD 版行号） | 本机形（`WISP_ENV` 未设） | CI 形（`WISP_ENV=test`） |
|---|---|---|---|
| **(i) 摘 pin** | `:261 t.Setenv("WISP_ENV", string(buildinfo.EnvDev))` | **rc=1 红** | **rc=1 红** |
| **(ii) 摘「先栽」** | `:260 t.Setenv("WISP_ENV", "test")` | **rc=0 全绿**（6 行结果全 PASS） | **rc=0 全绿**（6 行结果全 PASS） |
| **(iv) 两行都摘**（我另加的一发） | `:260` + `:261` | **rc=0 全绿**（6 行结果全 PASS） | **rc=1 红**，但红的**不是** CI 原来那 4 红 1 绿，而是 `:321 premise broke: the pin did not hold - WISP_ENV resolves to "test" (ambient on entry was "test", present=true)` -> 顶层一枚 FAIL、5 枚子用例一行都不出现 |

**(i) 的红句（两形各一条，原样）：**

```
    dataroot_128_test.go:322: premise broke: the pin did not hold - WISP_ENV resolves to "test"
      (ambient on entry was "", present=false), ...      <- 本机形
      (ambient on entry was "test", present=true), ...   <- CI 形
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (顶层，5 枚子用例一行都没出现)
```

⇒ **MUT-4 复现成立**，红的正是 helper 里那句自证（不是 CI 原来那句含混的 marker 缺失）。
**一处行号细节要如实记**：我量到的是 `:322`，报告写的是 `:323`。**两个数都对，框不同**：
`t.Fatalf` 在 `t.Helper()` 标记的函数里 -> Go 报的是**调用点**那枚文件行；我把 `:261` 摘掉之后，
调用点 `pinEnvThatAsksTheOS128(t)` 自己从 `:323` 上移一行变成 `:322`。它引的是**未变异文件**的 `:323`
（读者在 HEAD 上翻到的就是那一行），我引的是**我这棵变异树**的 `:322`。这一处不判它错。

**(ii) 的读数就是本格最该量的一句话**：把「先栽 `test`」那一行摘掉，**在 `WISP_ENV` 未设的机器上它照样绿**，
连 CI 形也绿（因为剩下那行 pin 单独就够盖住 ambient）。⇒ **「helper 在自己家里变成装饰」这一句，量出来是真的**：
没有栽在前那一步，`:263-265` 那枚 `buildinfo.EnvString() != EnvDev` 自证在 dev 主机上**恒真**，
它既不证明 pin 存在、也不证明 pin 有效。**这不是 c2fa2e9 的缺陷**（栽在前正是为了避免这个形状），
而是它那句「栽在前 ⇒ pin 有牙」的**反面被我这发正面量出来了**。

**(iv) 把这两发合起来才是答案**：两行都摘掉 -> **本机全绿（零信号）**，CI 形才红。
（红句由 premise 那句代劳，不再是我 §2.2 抄的那条含混的 marker 缺失 —— 这一点我原先写成了"逐字红回 CI 那一发"，
读日志后**改口**：`M-IV-both-removed-testenv.log` 现量只有 1 行 `--- FAIL`、断言点 `:321`，5 枚子用例根本没跑。）
⇒ 「栽在前」买到的东西，用一条读数说完：
**它把「这枚 pin 有没有生效」从一枚只在没人盯的 runner 上才响的主张，变成一枚在写码的人自己机器上就响的主张。**
没有它，(i) 那一发（摘 pin）在本机零信号，只有 (iv) 那种"两行一起消失"才会漏到 CI ——
而这正是票 140 立起来的那族形状（加固腿在 runner 上从不被咨询）。**判语：这一枚前提成立，且方向是加强（多红、不放宽、零 SKIP）。**

| 判据 | 判语 | 档位 |
|---|---|---|
| 摘掉 pin -> 本机即红 | **成立**（两形都红，红句自证那句） | **独立复现**（我自己造的变异） |
| 摘掉「先栽」-> 本机是否照样绿（= helper 成装饰） | **是，两形都绿**；栽在前的价值由 (iv) 那发正面量出 | **独立复现** |
| 「pin 有牙」这句话本身 | **成立**，但**必须连同 (iv) 一起读**才不是修辞：单看 (ii)，pin 自己是无牙的，牙是"栽在前 + pin"这一对给的 | 独立复现 + 我的补充判语 |
| 它的 MUT-4 落地三证（`grep -n` 命中 `:260,261`、`go build` rc=0、`go vet` rc=0 再读数） | 我这发用 `go test` 直接编译并跑到断言（编译不过就拿不到红句）-> **等价可核**；它引的 `dataroot_128_test.go:260,261` 与 HEAD 行号逐字对上 | 日志＋归档抽验 |

### 3.5 §3 那枚 commit 的回显

```
e220bf5 evidence(128 AC#4 r1 §3): 「pin 有牙」我自己造三发变异量（都只动测试文件、生产码零改动）——(i) 摘 pin：两形都红，红句是 premise broke: WISP_ENV resolves to "test"，行号 :322 vs 报告的 :323 属同一调用点不同框（t.Helper 把 Fatalf 报到调用点，摘掉上面一行它自己上移），不判它错；(ii) 摘「先栽 test」：两形都全绿＝helper 在自己家里确实成装饰（这不是缺陷，是栽在前的理由被反面量出来）；(iv) 我另加：两行都摘 → 本机零信号、只有 CI 形红（红在 premise 那句、5 枚子用例不出现；这条我先写成"逐字红回 CI"、读日志后改口并留了改口记录）⇒ 判语：前提成立且方向是加强

docs/evidence/s1/128-ac4-r1-acceptance.md
```

---

## 4. 门禁读数复核：`go test -count=2 -v ./cmd/wisp/`，我量了**四发**（两形 × 修前修后）

命令逐字（`<snap>` 是 §0 那两棵仓外快照；`PATH` 用 MSYS 形，与 `scripts/wisp-cli-tests.sh` 的自证一致）：

```
export PATH="<snap>/third_party/sherpa-onnx:$PATH"
WISP_ENV=test go test -count=2 -v ./cmd/wisp/     # CI 同步形
env -u WISP_ENV PATH="<snap>/third_party/sherpa-onnx:$PATH" go test -count=2 -v ./cmd/wisp/   # 本机默认形
```

四数只从 `-v` 数（`^=== RUN` / `--- PASS` / `--- FAIL` / `--- SKIP`，**顶层与缩进子用例分开数再相加**，
避免"哪个口径"糊掉）；`panic` 用**整文件不分大小写**扫，不是只扫 `^panic:`：

| 我的发 | 树 | env 形 | 包级 rc | `=== RUN` | `--- PASS`（顶层+子） | `--- FAIL`（顶层+子） | `--- SKIP` | `panic` 命中 | 去重名 | 对它的哪一发 |
|---|---|---|---|---|---|---|---|---|---|---|
| **G6** | `snap-pre` | 未设 | 0（`ok`） | **202** | **108+94=202** | 0 | **0** | **0** | 101 | = 它 **R1** 202/202/0/0 |
| **G5** | `snap-pre` | `test` | **1**（`FAIL`） | **202** | **106+86=192** | **2+8=10** | **0** | **0** | 101 | = 它 **R2** 202/192/**10**/0 |
| **G1** | `snap-post` | `test` | 0（`ok`） | **202** | **202**（108+94） | **0** | **0** | **0** | 101 | = 它 **R3** 202/202/0/0 |
| **G2** | `snap-post` | 未设 | 0（`ok`） | **202** | **202**（108+94） | **0** | **0** | **0** | 101 | = 它 **R4** 202/202/0/0 |

⇒ **四发全部与它报的数逐字一致，我没有量到不一致。** 档位：独立复现（四发全是我自己跑的）。

**名册差集（本格要求"差集要给"，不是只给四个数）**：

| 比 | 怎么算的 | 结果 |
|---|---|---|
| G1 ↔ G2（修后两形之间，AC#4 要的那"至少两形"） | 各取 `--- (PASS\|FAIL\|SKIP): ` 后的名字、去空白、`sort -u`、`comm -3` | **两侧都 101 名，差集空** -> `WISP_ENV` 不改变跑哪些用例，只改变其中一枚的判定路径 |
| G6 ↔ G1、G5 ↔ G2 之间（修前 vs 修后） | 同上 | **差集空**：`c2fa2e9` 没加也没删任何用例名。**四发名册两两比对（6 对）`comm -3` 全空**，各 101 名 |
| 尺的独立版 | `git diff c2fa2e9^..c2fa2e9 \| grep -c '^+func Test'` = **0**；`grep -c '^-func Test'` = **0**（`c94927d..c2fa2e9` 那对也是 0/0） | 与名册差集**同向**，两条尺互相不抵账 |
| G5 里那 10 行 FAIL 是谁 | 逐名 `uniq -c` | `TestAC2EveryLeg...128` 顶层 2 + `/cmdDoctor`、`/cmdModels`、`/cmdProviders`、`/runTextTask` 各 2 = **10**，5 枚名字一条不多一条不少 |
| G5 里真进程那一半 | `grep -cE '^[ -]*--- (PASS\|FAIL): TestAC2RealProcess[A-Za-z0-9/-]*128'` | **10 次出现、全 PASS**（顶层 2 + 四枚子用例各 2），含 `resident` 那枚 -> 「实现方那句真进程读数在 CI 形下不塌」在**这台机上**也复现了（与 §6 的 CI 原文互为两版） |
| 它 §2 那句「跨日基线 196 -> 202 的 6 行差未逐名核」 | 我**同样核不了**：09-23 那发的名册文件不在，本程也没假装核过 | **如实登记为未核**（本格不需要它成立：本格四发都是同日同法） |

**门禁档位**：四数 + 名册 + panic = **独立复现**；它那条"跨日 196→202 归因到 09-24 批次"= **仅自述**（我没核，也不替它背书）。

### 4.1 §4 那枚 commit 的回显

```
87339ba evidence(128 AC#4 r1 §4): 门禁读数我自己量四发（snap-pre/snap-post × WISP_ENV=test/未设，各 -count=2 -v ./cmd/wisp/）——G6 202/202/0/0、G5 202/192/10/0、G1 202/202/0/0、G2 202/202/0/0，四发 panic 整文件不分大小写命中 0、去重名册各 101 名且两两 6 对 comm -3 全空 ⇒ 与它报的 R1–R4 逐字一致，我没量到不一致；G5 那 10 行 FAIL 逐名 uniq 就是 EveryLeg128 顶层+4 子各 2 遍；真进程那一半在 G5 里 10 次出现全 PASS（含 resident）；+func Test/-func Test 现量 0/0；跨日 196→202 那 6 行我也核不了，如实登记未核不替它背书

docs/evidence/s1/128-ac4-r1-acceptance.md
```

---

## 5. 单点回退：只撤生产码那一处，问「哪条用例变得不响」（这发是本格的因果闭环，不是门禁的重复）

只撤 45 行里的 0 行、只撤生产码里的 **1 行**：`cmd/wisp/doctor.go:261`
`return "", dataDirUnresolved128(env, err)` -> `base = "."`（= AC#3 的 M-1 形状，票面 `:23` 要的就是这一发）。
两棵修后/修前树各造一次（`snap-post-mut-base` / `snap-pre-mut-base`），整包跑 `-count=1 -v ./cmd/wisp/`：

| 发 | 树 | env 形 | `=== RUN` | PASS | FAIL（顶层+子） | SKIP | panic | 红名清单 |
|---|---|---|---|---|---|---|---|---|
| **G3** | **修后树** + `base="."` | 未设 | 96 | 51+38 | **3+4=7** | 0 | 0 | `TestAC2ResolveDataDir...128`(+`/dev`,`/prod`)、`TestAC2EveryLeg...128`(顶层)、`TestAC2RealProcess...128`(+`/run`,`/doctor`) |
| **G4** | **修后树** + `base="."` | `test` | 96 | 51+38 | **3+4=7** | 0 | 0 | **与 G3 逐名同** |
| 对照 | 修前树 + `base="."` | 未设 | 只跑了这一枚用例（`-run`，整包**未**跑） | 1 子 PASS | 顶层+4 腿 FAIL | 0 | 0 | 四条腿在 `:295` 报 marker 缺失、其中 **3 条**另在 `:299` 报 `wrote into the start-up directory`（= AC#3 §55 记的那一形） |
| 对照 | 修前树 + `base="."` | `test` | 同上 | 1 子 PASS | 顶层+4 腿 FAIL | 0 | 0 | **与"修前树未变异 + `WISP_ENV=test`"（=P1）结构骨架 16 行逐字同**，见 §6 第二句 |

**答「删掉/改回哪一条用例会变红」这一问（AC#4 本格要能答）**：
- **修后树上，把生产码退回 `base="."` -> CI 形（`WISP_ENV=test`）红**（G4 rc=1，7 枚 FAIL、0 SKIP）。
  这一句就是"CI 现在真的看着这条用例"的**因果**凭据，不是"我把它跑绿了"的**颜色**凭据。
- 红法**换了形状**：`TestAC2EveryLeg...128` 现在死在 `dataroot_128_test.go:323 premise broke: ... resolveDataDir returned ("wisp-dev", <nil>) instead of errDataDirUnresolved`
  —— 顶层 `t.Fatalf` 在 `t.Run` 循环**之前**，所以**5 枚腿子用例不出现**（`=== RUN` 因此 96 = 101 − 5）。
  这条要说清，别让人以为是"少跑了、所以更松"：它是**更紧**（修前那发是"腿各自报 marker 缺失"这种含混红，修后是"前提直接不成立"这种指名红）。
- 但**粒度确实掉了一格**：修前那发能逐腿指出是谁回落（`runTextTask`/`cmdModels`/…各印一条 `wrote into the start-up directory`），
  修后这一发里 `wrote into the start-up directory` 只剩 **2 次**，都出自 `dataroot_128_windows_test.go:79`（真进程那两枚腿）。
  ⇒ **整包层面信号没丢**（那一句仍然印、`TestAC2ResolveDataDir...{dev,prod}` 仍然红），丢的只是 `EveryLeg128` 内部的"哪条腿"信息。**判：可接受，登记不判退回。**
- **一处既有文案的过度承诺要顺手点名（不是本枚 commit 造的）**：文件头 `:35-40` 的 `AC#3 MUTATION ANCHOR` 说
  改回 `base="."` 时 `TestAC2RefusalMarkersAreNotAShortenableList128` 也会红 —— 我两发（G3/G4）里它**都没红**
  （它只查 `rescueMarkers128` 的长度下限与前缀行为，与 `resolveDataDir` 无关）。那段注释是 `4e5d240` 写的、
  `c2fa2e9` 只把它当上下文行读过 ⇒ **不算进本格的账**，但**建议编排者开一枚一字改的追件**（或并入票 135 那两格），
  免得下一位按那句找红名找不到、以为门禁哑了。

**档位**：独立复现（G3/G4 是我自己造、自己跑、自己逐名数的）。

### 5.1 §5 那枚 commit 的回显

```
96953695e24b40c97c38d9bd0e5cbbcfff66fc6c evidence(128 AC#4 r1 §5-§6): 单点回退＋副产物两判。…（全文见 git log）

docs/evidence/s1/128-ac4-r1-acceptance.md
```

⇒ **只列这一枚路径**（那枚 commit 前 `git diff --cached --name-only` 里出现过兄弟程 staged 的
`docs/evidence/s1/140-static-inventory-r1.md`，我按带 pathspec 的 `git commit -- <我的路径>` 提交，
事后双量：我这枚 `--name-only` 只有本文件；兄弟程随后那枚 `29d8938` 的 `--name-only` 只有它自己的文件。
**两边的归属都没被卷走**，登记在这里是因为这条尺在共享树里只有事后核才算数。）

---

## 6. 副产物那一发的连带判断（本格最值钱的一条）—— 两句都要明确裁

**先把它主张的那件事实独立量一遍，判真假**：「CI 上那四条红里 `assertDirEmpty128` 没报，
因为落点被搬去了仓外绝对路径 -> 修之前"当前目录一字节不许多"在 runner 上是**被另一枚回落点满足的，不是被拒绝满足的**」。

| 我的尺 | 读数 |
|---|---|
| P1（`snap-pre` + `WISP_ENV=test`，整份日志）里目录类句子命中数：`grep -c "start-up directory\\|did not write\\|dirEmpty\\|not empty" P1-pre-cishape.log` | **0** —— 四条红全是 `:295` 的 marker 缺失，目录那一格一条没响 |
| 同发的日志里落点被印成什么 | `msg="wisp: persistent log sink installed" dir=C:\Users\<me>\AppData\Local\Temp\wisp-test-<pid>\logs`（我这台机）、CI 原文是 `...\runneradmin\...\wisp-test-8104\logs` -> **两边都是仓外绝对路径** |
| CI 自己那一步的原文（我 `gh` 取的那份） | 同上，且 `[PASS] data dir writable (test) C:\Users\runneradmin\AppData\Local\Temp\wisp-test-8104` |
| 目录断言在真进程那一半是不是也做 | **做**：`cmd/wisp/dataroot_128_windows_test.go:79 assertDirEmpty128(t, cwd, "real leg "+leg.name)`，且它跑的是**真删 APPDATA 的子进程** |

⇒ **主张为真**，档位：独立复现（本机 + CI 原文两版都我亲自取到）。下面是两句要裁的话。

### 6.1 ① 这会不会回头动摇票 128 `AC#2` 那格已判成立的裁决？—— **不会。那一半是真实的。**

`AC#2` 的判据本体（票面 `:38-39`）是「同形下**三条腿都拒、拒绝原因可被人读懂**，并且**当前目录里一个字节都不许多出来**」。
实现方 `T128-ac23` 交回那句"run rc=2 / secret rc=2 / doctor rc=1 / 常驻 rc=1、四枚 CWD 全空"，量的是**真进程 + `WISP_ENV=dev`** 那一半。
**那一半在我手上有三条独立凭据，没有一条依赖进程内那枚 seam**：

1. **CI 自己的日志**（run `35967768017`，修前那版码）：`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128`
   顶层 + `run` / `secret-list` / `doctor` / `resident` **五枚全 PASS**，且日志逐条印着
   `APPDATA unset, WISP_ENV=dev, cwd=C:\Users\RUNNER~1\...\001 -> rc=2`（四枚子用例各一条），拒绝文案含全部自救句。
   ⇒ runner 上真发生过"拒绝 + CWD 空"，这不是本机专属。
2. **我这台机**：G5（修前树、`WISP_ENV=test`、整包 count=2）里那枚用例**10 次出现、0 枚红**；
   G1/G2（修后两形）里同样 10 次全 PASS。⇒ ambient 那枚变量动不到它，因为 `appDataFreeEnv128` 末行
   （`dataroot_128_windows_test.go:113 return append(kept, "WISP_ENV=dev")`）**给子进程钉了 dev** —— 它的凭据是自足的。
3. **它有没有牙**：G3/G4（把生产码退回 `base="."`）里 `TestAC2RealProcess...{,/run,/doctor}` **照样红**，
   红句 `wrote into the start-up directory` 出自上面那枚 `:79`。⇒ 这一半不是"恒绿的装饰"，它咬得住。

⇒ **判语：`AC#2` 那格不动。** 本枚（AC#4）暴露的是**第二见证**（进程内那五条腿）在 runner 上曾经过弱，
不是**已入账的那一枚凭据**为弱。顺带把 `AC#3` 也过一遍尺（因为它的红名清单与这一枚同源）：
票面 `:23` 那句"改回 `base="."` -> 判定用例必须红、红名点到它"是在 **dev 形**下量的，那一发我复到了（§5 对照行第一发）；
而 AC#4 之后它**两形都红**（G3/G4）⇒ **AC#2/AC#3 两格都只会因为本枚被加强，不会因为本枚被追回。**
唯一要落账的一句限定（写进台账、不改写历史）：**修前那版的"三腿都拒"在 runner 上只有真进程那一半在作证**，
读到 `AC#2` 那段的人不该以为"进程内那五条腿也在 runner 上作证"。这条限定**不改变档位**（那一格本来就是真进程凭据撑的）。

### 6.2 ②「CI 会看着这条用例」这句话在此之前是不是一句装饰？—— **是装饰，而且是恒红型装饰（比恒绿更贵）**

要把两件事分开，不然这句会被读成"CI 没跑它"：
- **CI 有没有跑到它**：跑到了。CI 那一步自己打的四数 `=== RUN=101 --- PASS=49 --- FAIL=5 --- SKIP=0` 里就含它，名字、缩进深度、断言行号 `:295` 全在。
- **CI 有没有"看着"它（=它的判定路径能不能区分被裁的两种行为）**：**没有**。这一发的凭据是我自己造的两份日志：

| 比 | 同一棵 `snap-pre` 树、同一条命令、同一个 `WISP_ENV=test` | 结果 |
|---|---|---|
| (a) 生产码**未动** | `P1-pre-cishape.log` | 顶层 FAIL + 4 腿 FAIL + `resolveSecretLayout` PASS |
| (b) 生产码退回 `base="."`（= AC#2 明令禁止、AC#1 实测过的搬家行为） | `M-PRE-basefallback-testenv.log` | **同上** |
| 两份的**结构骨架**（`=== RUN` / `--- PASS\|FAIL` / 含 `AC#2 RED` 的断言行） | 各 16 行 | `diff` **rc=0，零输出** |
| 两份的**逐字节**差 | — | 只差时间戳、`wisp-test-<pid>`、`go-build<hash>`、时长；断言层**一字不差**（`diff` 输出 10 段，全在这四类字节上） |

⇒ **修之前，这条用例在 runner 上对"拒绝 / 回落到 CWD"这一对它是被裁来区分的区别，输出的是一模一样的字。**
它的红还是**固定红**（那四条腿在 runner 上无论如何都红），所以它同时**吃掉了自己的信号量**：
下一位真伤（比如有人把 `base="."` 放回去）在这条用例上**不会引起任何颜色变化**。
这正是本仓记过的那族形状——**"这条用例在这个 runner 上没有分母"**，以及**恒真判据**的近亲（恒红判据）。

**「删掉哪条用例会变红」这一问的答法（两版都给，因为它就是这句判据的正反两面）**：
- **修前**：把 `doctor.go:261` 那枚 `return "", dataDirUnresolved128(env, err)` 删掉/退回 `base="."`
  -> **CI 形下没有任何一条用例改变颜色**（(a) 与 (b) 逐字同形）。⇒ 那一格当时在 runner 上答不出这一问。
- **修后**：同一处同一撤 -> CI 形下 **7 枚 FAIL**（`TestAC2ResolveDataDir...128{,/dev,/prod}`、
  `TestAC2EveryLeg...128`、`TestAC2RealProcess...128{,/run,/doctor}`），0 SKIP、0 panic（= 我的 G4）。
  ⇒ 这一问现在答得出，且答法是"三条腿里两半都响"（进程内那枚 + 真进程那两枚）。

**判语（这两句一起构成我对本格的终判依据）**：`c2fa2e9` 买的不是"把 CI 那四条红洗绿"，
而是**给这条用例在 runner 上装回了分母**；副产物那句登记我判它**为真、且是本次交付里信息量最大的一条**，
它不追回任何已入账的格子，但它把"门禁颜色变绿"这一件事从判据里**降级**了：
以后引用本格，凭据应当是 §5/§6 这两发（红因换轨 + 由红转绿），不是四数。

---

## 7. 其余三套门禁与四把尺（AC#4 那一格逐项，全部我自己真跑）

| 门禁 | 我在哪棵树上跑的 | 命令原文 | rc | 读数 |
|---|---|---|---|---|
| `gofmt` | `snap-post`（`c2fa2e9` 纯净快照） | `gofmt -l cmd/wisp/` | **0** | **空输出**（0 枚文件被点名） |
| `gofumpt` | 同上 | `"$(go env GOPATH)/bin/gofumpt.exe" --version` -> **`v0.12.0 (go1.27.1)`**；`… -l cmd/wisp/` | **0** | **空输出**。⇒ **盘上现量就是 v0.12.0**，票面/派单写的 v0.7.0 是过期值；它 §2 那处更正我复到了 |
| `go vet`（本机 GOOS=windows） | 同上 | `go vet ./cmd/wisp/` | **0** | 空输出（0 字节，`wc -c` 现量） |
| `go vet`（GOOS=linux，**只编译不执行**） | `snap-post` 与 `snap-pre` 各一发 | `GOOS=linux go vet ./cmd/wisp/` | 1 / 1 | 两版**逐字相同**：`build constraints exclude all Go files in …sherpa-onnx-go-linux@v1.13.8`（`diff` rc=0）=> **既有宿主交叉形状，不是本枚造的**；本表**不拿它当代跑** |
| `sh scripts/d22scan.sh`（纯净快照） | `snap-post`（`git archive` 出来的，工作树里 owner 那 16 枚 `design/**` 未提交删除进不了它） | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；**正控制真跑了**：`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` |
| 台账各 scope 不降 | 同上（d22scan 自打的 examined 数） | `grep -oE "examined +[0-9]+ …"` | — | `bans #1-5 internal/=203`、`bans #1-5 cmd/=22`、`ban #6 frontend/=40`、`ban #7 internal/tools/=18`、`ban #8 design/=16`、`ban #8 frontend/=40`、`ban #8 internal/=404`、`ban #8 cmd/=39` —— **与它 §5 那张表八枚数逐字同**，且都 >= 台账 `:4192` 记的口径（404 vs 397、39 vs 38） |

**票 123 那批不许被放宽（四条尺，我各量一遍）**：

| 尺 | 读数 |
|---|---|
| 本程 diff 是否含 `run_test.go` / `run_mode101_test.go` / `internal/agent/approval/**` | **共享树里 `HEAD` 是动的**，所以两端钉死再数：`git diff --name-only c94927d..c23d825`（它 §6 那把尺用的就是这条区间，末端取它的交件枚）现量 **8 枚路径**（3 枚 `.scratch/wisp/issues/*` + 3 枚 `docs/evidence/s1/*` + `cmd/wisp/dataroot_128_test.go` + 2 枚 `docs/reports/*`，其中 119/HANDOVER/台账是兄弟程在这段区间里落的），对它点名的三类名字 `grep -Ec` = **0**；被验那一枚单独再数：`git show --name-only --format="" c2fa2e9` = **只有 `cmd/wisp/dataroot_128_test.go`**。（它 §6 写"现量 6 枚"，那是在它自己那一刻的 HEAD 上数的，量级差是兄弟程的 commit 挤进同一段区间，**不是它漏了东西**。） |
| `internal/agent/approval/queue.go` 那个 300 秒 | 现量 `107: DefaultApprovalTimeout = 300 * time.Second`（**一字节未动**，且该文件不在上面那 8 枚路径里） |
| 那四枚用例在修后两形里的颜色 | G1（CI 形）与 G2（本机形）里各 **PASS=2 / FAIL=0**（count=2 => 各两遍）⇒ 它们没被改成 Skip、没被放宽（`t.Skip` 在 `dataroot_128_test.go` / `_windows_test.go` 现量 **0 / 0**） |
| 它们的 runner 红与本枚变量共不共因 | **在这台机上不共因**：`WISP_ENV=test` 整包跑（G1）那四枚全绿；CI 原文里它们的红句是 `run_test.go:378: an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认）…` 与 `run_mode101_test.go:309`（配置里写着 `confirm_timeout_sec = 1` 却被默认值取代）=> 红因是**那四枚自己读不到自己写的 config**，与本票那枚 seam 无关。**但**：CI 那四条红的真实根因**不在本格判据里**，我只登记现象 + 一句要往下追的话（见 §9） |

**owner 真实数据目录**（本格要求"两次数都记"）：开工时与全部读数跑完后各一次，
`ls -A %APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 均 **0 条目**（两枚目录都存在、都空）=> 本程未写入。
（本程新建的落点只有 `D:\tmp\wt-acc128-ac4-r1*` 那些仓外快照/日志与 `Temp\wisp-test-<pid>`，都在仓外。）
> **§7 那一行的一处口径要更正（追加，不改写）**：我实际取的两次是**读数中段 `18:0x`（与 `d22scan` 同批）**与
> **收尾 `18:21:59 +08`**，两条都是 **0 条目**；"开工时"那一次我**没有单独取**（开工时我先取的是 §0 那批身份读数）。
> 影响：结论（本程未写入 owner 真实目录）不变，因为它跑的是 `t.TempDir()`/仓外根，且两枚目录在两枚时刻皆空。

### 7.1 补一发：把 **CI 那一步本身**逐字搬进快照跑（AC#4 那句"四数"之外的同形凭据）

`AC#4` 写的是 `go test -count=2 -v ./cmd/wisp/`（我 §4 已经量了）；但 runner 真正执行的是
`bash scripts/wisp-cli-tests.sh` -> `portable-tests.sh --scope=cli` -> `tools/d22scan/runtests.sh ./cmd/wisp/ -count=1 -skip <台账造的 pattern>`
（SKIP 致命、零 PASS 零 FAIL 致命、包必须打出自己的顶层结果行）。两棵树、同一个 `WISP_ENV=test`：

| 发 | 树 | 步级 rc | CI 自己打的四数 |
|---|---|---|---|
| **S1** | `snap-post`（被验版） | **0** | `=== RUN=101  --- PASS=54  --- FAIL=0  --- SKIP=0` + `portable-tests.sh:   ok (own line)  github.com/CarlosShao/wisp/cmd/wisp` |
| **S2** | `snap-pre`（`c2fa2e9^`） | **1** | `=== RUN=101  --- PASS=53  --- FAIL=1  --- SKIP=0` + `FAIL (own line)` + `strict runner exited 1 for scope=[./cmd/wisp/]` |

⇒ **CI 那一步（含它自己的三道 guard 与 -skip 台账）在修前红、修后绿，同一枚 `TestAC2EveryLeg...128` 是唯一差别。**
这条比"四数"更硬，因为它跑的是 runner 用的那台仪器本身，不是我对它的模仿。
（**口径要算清，别把 54 与 101 当成同一件事**：Go 1.27 的 `-v` 里 `=== RUN` 不分深度都落在第 0 列、
`--- PASS|FAIL` 只有顶层落在第 0 列、子用例缩进 4 空格 —— 与它 §2 开头那句口径声明一致。
于是 **runner 那一步的顶层结果总数 = 49 PASS + 5 FAIL = 54**，与我 S1 的 `PASS=54`、S2 的 `53+1=54` **同一枚分母**；
差的只是颜色分配：CI 那 5 枚顶层 FAIL = 票 128 那 1 枚 + 票 123 那一族 4 枚，我这边那 4 枚全绿 -> 只剩 1 枚。
⇒ 不是"我换了台更强的仪器"，是**同一枚仪器、同一枚分母、颜色只因这枚 commit 而变**。）
档位：**独立复现**。⇒ **本格的终判不依赖"推送后去看 CI"**（那是编排者的 next=①，属于加强件不是必需件）。

### 7.2 §7 与 §7.1 那两枚 commit 的回显

```
9c8eb3f evidence(128 AC#4 r1 §7): 其余三套门禁与四把尺全部我自己真跑…（全文见 git log）
docs/evidence/s1/128-ac4-r1-acceptance.md

f0d023e evidence(128 AC#4 r1 §7.1): 补一发 runner 仪器本身的逐字同形…（全文见 git log）
docs/evidence/s1/128-ac4-r1-acceptance.md
```

（两枚提交前 `git diff --cached --name-only` 各现量 1 枚路径＝本文件；§5-§6 那枚之前出现过兄弟程 staged 的
`140-static-inventory-r1.md`，事后双量确认两边归属都没被卷走，见 §5.1。）

---

## 8. 总判

判据本体是票面 `:25-27` 那一格，加它这轮主张的因果判定。逐条对：

| # | 判据（票面原文，`:25-27`） | 我的凭据（本节只指路，不重述读数） | 档位 | 判语 |
|---|---|---|---|---|
| 1 | `go test -count=2 -v ./cmd/wisp/` 四数 | §4 表 G1/G2（修后两形）+ G5/G6（修前两形），四数、panic、101 名两两差集全为我的现量 | 独立复现 | **成立**（与它 R1–R4 逐字一致，我没量到任何不一致） |
| 2 | `gofmt`/`gofumpt` 真跑 | §7 前两行：`gofmt -l cmd/wisp/` 空、`gofumpt --version` = `v0.12.0 (go1.27.1)`、`-l` 空 | 独立复现 | **成立**（票面写的 v0.7.0 确为过期值，它的更正是对的） |
| 3 | `go vet` | §7 第 3/4 行：本机 `./cmd/wisp/` rc=0 零字节；`GOOS=linux go vet` 修前/修后**逐字同一行**（`diff` rc=0）| 独立复现（本机半）/ **仅自述**（容器 `CGO_ENABLED=1` 那两发，我没起容器，见 §9#1） | **成立（按包 scope）**，附一条口径：linux-cgo 那半我没复 |
| 4 | `sh scripts/d22scan.sh` 纯净快照 rc=0 + 台账各 scope 不降 | §7 第 5/6 行：`git archive` 树上 rc=0 `clean`、正控制真跑（`PASS=21 FAIL=0 SKIP=0 / === RUN=31`）、八 scope 203/22/40/18/16/40/404/39，与台账 `:4192` 的 397/38 口径**只增不减** | 独立复现 | **成立** |
| 5 | ⚠ 票 123 那批（`审批超时（1/300 秒未确认），C18 一律判拒绝`）不许被放宽换绿 | §7 那四把尺：三类名字在交付区间 `grep -Ec`=0、`queue.go:107` 现量 `300 * time.Second` 一字节未动、那四枚在 G1/G2 各 `PASS=2/FAIL=0`、两枚 128 测试文件 `t.Skip` 现量 0/0 | 独立复现 | **成立：未放宽、未换 Skip、未动常量** |
| 6 | （交付主张）「本机绿 / CI 红，根因＝`ci.yml:337` 那枚 job 级 `WISP_ENV: test`」 | §1 静态三处 + §2.1 blob 号核身份 + §2.3/§2.4 复现（P1 与 CI 骨架 16 行 `diff` rc=0；P2 同形全绿）+ §7.1 runner 仪器本身 S1/S2 一步红绿分明 | 独立复现 | **成立** |
| 7 | （交付主张）只补前提、没动期望值 | §1.2：真内容删除行**只有 1 枚**且被 helper 在 `:262` 找回；`rescueMarkers128`/`rc==0`/`assertDirEmpty128` 未出现在 diff 里 | 独立复现 | **成立**（不是"用放宽断言换绿"这一形） |
| 8 | （交付主张）"pin 有牙" | §3 三发变异（(i) 两形都红、(ii) 两形都绿、(iv) 本机零信号/CI 才响） | 独立复现 | **成立**，但判语限定为"栽在前 + pin 这一对才有牙"（单味不承重） |
| 9 | （交付主张）副产物登记 + 本格最值钱那两句 | §6（① 不追回 `AC#2`，凭据三条；② "CI 会看着这条用例"在此之前是**恒红型装饰**，"删掉哪条用例会变红"两版答语都在 §5/§6） | 独立复现 | **成立**，且登记为**本次交付信息量最大的一条** |

### 总判语

> **`AC#4` 这一格：成立，无附条件。** 判据五项全部我亲自复到，且**没有一项靠实现方的自述**。
> 可以翻 `AC#4` 的勾（**勾由编排者打，本程一枚未翻、票名未改、票面 128 那枚文件一枚字节未写**）。
> 交付里最强的两发不是"四数"，是 **§7.1 的 S1/S2**（runner 用的那台仪器本身修前红、修后绿）与 **§6.2 那张 16 行同形表**
> （修前那版对"拒绝 vs 回落 CWD"输出的字一模一样）。以后引用本格，凭据应当绑这两发，不是绑颜色。

**三件往前带的话（都不构成本格退回，也都不需要我这一程动手）**：
1. **两处行号引用**：`§1.4` 那处（报告 `:131` -> HEAD `:139`）与 `§1.5` 那处（`doctor.go:255-258` -> 正确是 `256-259`，
   这一处写下来时就差一）。建议**原地追加**一行更正，不改写原句（append-only 纪律）。
2. **文件头 `AC#3 MUTATION ANCHOR`（`:35-40`）过度承诺**：它说改回 `base="."` 时
   `TestAC2RefusalMarkersAreNotAShortenableList128` 也会红，我 G3/G4 两发它都没红（那枚只量 marker 清单长度与前缀，与
   `resolveDataDir` 无关）。那段是 `4e5d240` 写的既有文，**不算进本枚 commit 的账**，但它会让下一位找红名找不到 ->
   建议开一枚一字改的追件（或并入票 135 那两格）。
3. **票 123 那一族在 CI 上那 4 枚红，与 128 不共因**（我这台机 `WISP_ENV=test` 整包跑它们全绿，见 G1；
   CI 原文里它们的红句是"自己写的 config 读不到"：`run_test.go:378` 拿到默认 300 秒、`run_mode101_test.go:309` 同形）。
   ⇒ 裁 `ci.yml:227/337/482/540` 那四枚 job 级 `WISP_ENV: test` 时**别把它们一起结掉**：撤 env 未必能让那 4 枚变绿，
   而那 4 枚的真因不在本票地界。这条是给票 140 与编排者的线索，本程**未追**。

**承重判定**（尺：摘掉任意一味，是否**存在一发变异从此打不红**。每一味都真摘真跑，不靠推理；下面五发的命令同 §3，都是 `-count=1 -v -run 'TestAC2EveryLeg...128' ./cmd/wisp/` 在各自快照里，PATH 挂了那三枚 DLL）：

| 味 | 摘掉它之后我量的那一发 | 读数 | 承重？ |
|---|---|---|---|
| `:261` 那枚 pin | §3 (i)：两形（未设 / `test`） | 都 **rc=1**，红在 `premise broke: the pin did not hold` | **承重**（摘掉立刻响） |
| `:260` 那枚"先栽 test" | §3 (ii)：两形 | 都 **rc=0 全绿**；但连 `:261` 一起摘（§3 (iv)）-> 本机绿、CI 形红 | **单独不承重，与 `:261` 成对才承重**：它承的是"摘 pin 这件事在写码的人自己机器上也响"，不是"拒绝行为对不对" |
| `:262 failConfigDir128(t)`（那枚被搬进 helper 的行） | **W3**：摘掉它（seam 不再被重绑），两形 | 都 **rc=1**，红在 `:322 premise broke: with WISP_ENV="dev" and a failed OS read, resolveDataDir returned ...` | **承重，但被 `:266-272` 兜住**（摘了它还有下一味响） |
| `:266-272` 那两枚正面自证（`errDataDirUnresolved` + 空串） | **W1 vs W2**：把生产码 `doctor.go:297` 的 `%w` 身份换成一枚**同文案、不同 identity** 的 `errors.New`（=> 六枚 marker 一字未变、`errors.Is` 分类断掉）；W1 保留自证、W2 摘掉自证 | W1（保留）**rc=1**，红在 `:323 premise broke: ... instead of errDataDirUnresolved`；W2（摘掉）**rc=0、6 行结果全 PASS** | **承重** —— 摘掉它，这一发变异（分类被破坏而文案完好）在**任何平台都打不红**，腿级断言看不见它 |
| 同一味，换一发变异 | **M-V-real**：摘掉 `:266-272` 之后再退 `base="."`（`WISP_ENV=test`） | **rc=1**，4 枚腿各自红 + `start-up directory` 3 次（比保留自证时**更细**） | 对 `base="."` 这一发**不承重**（有没有都红，只是红法不同）——⇒ §5 那句"红法换形状"就是这一味造成的 |

⇒ **一句话总结这五行**：`c2fa2e9` 往 helper 里放的是**三味互相兜底的门**（pin / seam 重绑 / 正面自证），
摘任意一味都有另一味把对应的变异接住；**唯一摘了就没人接的，是"先栽 `test`"那一步对摘 pin 这件事的本机可见性** ——
而那一条恰好是 §3 (iv) 量出来、票 140 立起来的那族形状。**这一格我给它的修法判"没有可摘的装饰"。**

### 8.1 §8 那枚 commit 的回显

```
27a880a evidence(128 AC#4 r1 §8): 总判——AC#4 那一格成立、无附条件…（全文见 git log）
docs/evidence/s1/128-ac4-r1-acceptance.md
```

---

## 9. 我明确没核的清单（不核的就是没核的，别按已核读）

| # | 没核的那半 | 为什么没核 | 谁还能核 |
|---|---|---|---|
| 1 | 容器里 `golang:1.27` + `CGO_ENABLED=1` 的 `go vet ./cmd/wisp/`（它 §5 报 rc=0 两发） | 兄弟程 `acceptor-ticket119-ac7-r1` 此刻正在同一台机上取 winsec 读数，容器 cgo 编译会抢 CPU；我只做了宿主 `GOOS=linux go vet` 的**修前/修后差分**（同一行，`diff` rc=0） | 编排者：容器空闲时按它 §5 那行挂载形状复跑 |
| 2 | `go vet ./...`（整树）与任何整树 `go test ./...` | 简报红线：按包 scope，整树会吃到别人的东西 | 已在 §7 按包做完 |
| 3 | 它 §10 的"四枚 commit 账 / 九枚临时件"逐枚清单 | 我只抽验了我要用的那枚码 commit（`c2fa2e9` 的 numstat + `--name-only`）与它四枚证据 commit 的存在性（`cat-file -t` 全 `commit`、`is-ancestor HEAD` 全真） | 不需要：本格不依赖它 |
| 4 | 跨日基线 196 -> 202 那 6 行差 | 09-23 那发的名册文件已不在（它自己也登记同一条） | 无人（历史读数丢了就是丢了） |
| 5 | `memory.db`/`wisp.db` 真实落点、DPAPI blob 真写入 CWD 树、共享目录当 CWD 的 ACL、prod 环境整机形 | 这四条是 `AC#1` 就没量的（票面 `:52` 的"未验证四条"），本程一枚未碰 | 票 135 / 后续 |
| 6 | `-race` | 票面 `AC#4` 没要求，它没取，我同样未取 | 需要时按包 scope 取 |
| 7 | 票 123 那 4 枚在 CI 上"读不到自己写的 config"的**真因** | 超出本格判据；我只登记了"本机同形不共因"这条读数（G1 里它们 PASS=2/FAIL=0）与 CI 原文两句红名（`run_test.go:378`、`run_mode101_test.go:309`） | 票 123 / 票 140 的下一位 |
| 8 | `c2fa2e9` **推送之后** CI 那一步的颜色 | 本程不 push（纪律），所以核不到 | 编排者：推完读那一步，凭据应绑 §7.1 的 S1/S2 形状 |
| 9 | 常驻腿（`resident`）在 `ambient WISP_ENV=test` 下的行为 | 它只在子进程 `WISP_ENV=dev` 那一形下被量过（§6.1 第 1 条）。那枚腿的文案仍不含自救句（它 §7.5 登记的 `R-128-3`/`R-128-4` 在票 135 地界，本程一枚未动） | 票 135 |
| 10 | `resolveSecretLayout` 那枚第五腿在"两行都摘 + 生产码退回"叠加形下的行为 | 超出本格判据（它显式传 `EnvDev`，我 §3/§5 的每一发里它都是绿，没造组合） | 无人必要 |

---

## 10. 临时件路径（**只建不删**）＋ 注入两栏计数

**10.1 本程新建的一切，全部在仓外；读数件全部留着**（收尾时刻 `2026-09-24 18:15 +08` 现量）。
**一条纪律偏离要如实登记**：我对**自己刚拷出、还没取过数**的中间变异树执行过 `rm -rf` 后重建
（`snap-pre-mut-base` 1 次——那刻它内容还是 `snap-post` 的拷贝、未取任何读数；`snap-mut-v-drop-selfcheck` 2 次——两次都是我的 python `assert` 下标数错、变异根本没落地，见 §10.2 自伤第 4 条）。
⇒ **零枚日志/名册/CI 原文被删**，最终态 **11 棵快照全在盘上**（`ls -d snap-* | wc -l` 现量 = 11），可重跑性不受影响；但"临时件只建不删"这五个字，本程**没做到全字**，按偏离登记不按达标写。

| 路径 | 是什么 | 量出来的数 |
|---|---|---|
| `D:\tmp\wt-acc128-ac4-r1\snap-post` | 被验版纯净快照（`git archive c2fa2e9`），`.anchor-sha` 内写死 `c2fa2e98f7ce…` | `go.mod` 883 字节、`third_party/sherpa-onnx/*.dll` 3 枚（未跟踪件手拷） |
| `…\snap-pre` | `git archive c2fa2e9^` | 同上（`anchor-sha=2956897…`） |
| `…\snap-mut-i-drop-pin` / `snap-mut-ii-drop-plant` / `snap-mut-iv-both-removed` | §3 三发变异树 | 各 1 行 / 1 行 / 2 行摘除，均只动 `dataroot_128_test.go` |
| `…\snap-post-mut-base` / `snap-pre-mut-base` | §5 单点回退树（`doctor.go:261` -> `base = "."`） | 各 1 行 |
| `…\snap-mut-v-drop-selfcheck`、`snap-W1-selfcheck-kept`、`snap-W2-selfcheck-dropped`、`snap-W3-no-seam-rebind` | §8 承重那五发 | 见 §8 表 |
| `…\*.log` + `…\*.txt` + `…\*.names` | 全部原始读数：P1/P2、M-*（§8 那几发变异）、G1–G7（整包，含交回前复跑那发）、S1/S2（CI 同步形）、D-scan-post、VET-*、N-*/K-* 骨架与名册、HEAD-/WT-lines.txt（结构修复那次的内容守恒尺） | **收尾现量**（`ls -1 \| wc -l`）：目录条目 **55** = **11 棵快照树 + 32 枚 `.log` + 12 枚骨架/名册件**，总体积 **441 MB**（`du -sh`） |
| `D:\tmp\wt-acc128-ac4-r1-ref\` | `gh-selffetch.txt`（421465 字节，我自己取的 CI 原文）+ `gh-selffetch.err`（0 字节）+ `blob_test.go`（`c2fa2e9` 版测试文件原件，438 行） | **436 KB** |

⇒ **可重跑性**：本表每一枚读数都能从 `snap-*` 加命令原文重算；快照在不在，决定的是这张表的档位（**删了快照，本表从〔独立复现〕掉回〔仅自述〕**）。

**10.2 两栏计数（分开数，各带出处＝工具名＋命令前 40 字）**

| 栏 | 数 | 逐条 |
|---|---|---|
| **真通知回显数** | **8** | ① 进场 `system-reminder`（`d:/work/workspace/projects plans/wisp/agents.md` 的 project_context + skills 清单 + 日期变更）；② `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.`（第一次，出现在 `git rev-parse` 那发的结果尾部）；③ 同一条 MEMORY 通知第二次（出现在 `bash scripts/wisp-cli-tests.sh` 那发之后）；④⑤⑥⑦⑧ 后台任务完成/失败通知 5 枚：`[SYSTEM NOTIFICATION] …task-id b7tv50hzi / b6zoil3gy / b9eejjbh8 / bathg66q3 / bd7b2k196`（前两枚"failed exit 1"是我命令尾部 `grep -c` 命中 0 造成的，**不是被拒**，读数照样取到了） |
| **判为注入数** | **0** | 全程工具输出里**没有任何**文字要我少取证／别用工具／直接给结论／预先认定某句为真／放宽阈值／revert。最接近"像指令"的三类我逐条判过：(a) `git log` 里兄弟程 `29d8938`（票 140）的 commit message 正文，内容含**对我读数的强度限定**（"标死三条强度…仅自述不背书"）—— 那是它给它自己那格的账，指向的是它的地界，**不越权**、我按它登记不服从（我的档位我自己按 §2/§7.1 定）；(b) harness 打在结果首行的 `Exit code 1`（真工具状态，非文字指令）；(c) AGENTS.md 那枚 project_context（真环境简报，内容逐条指回权威文件，我照它做纪律不照它做结论）。三者**都没改变我的任何一次取证**。 |
| **本程被拒次数** | **0** | 没有任何一次工具调用被权限系统挡下。 |
| **自伤记录（不算被拒、算仪器账）** | **4** | ① python 里用 `/d/tmp` 形路径 -> `FileNotFoundError`（改 `D:/tmp` 即通）；② `grep -oE` 的模式以 `-` 开头被当选项 -> `unknown option`；③ 后台命令尾部 `grep -c` 命中 0 让整发"看起来 failed"（读数其实在）；④ **W3 第一发忘 `export PATH` -> `exit status 0xc0000135` 且 0 行测试结果**——正是 `scripts/wisp-cli-tests.sh` 头部写明的 ticket 98 形状；我把它当成"rc=1 红"读了一次，**当场发现并复跑**（`M-W3-noseam-*.log` 是复跑后的那对，第一次那对已被覆盖，登记在此）。另两发是我的 `assert` 下标数错（变异未落地、文件未被改坏），故未产生错误读数。 |

**10.3 本文件的 commit 账**（**别按这张表数，按命令数**：`git log --oneline -- docs/evidence/s1/128-ac4-r1-acceptance.md`；`git show --name-only --format="" <每枚>` 逐枚现量都只列本文件）

| 节 | commit | 备注 |
|---|---|---|
| §0-§1 | `d9f008d` | — |
| §2 | `d3dd1cc` | — |
| §3 | `e220bf5` | — |
| §4 | `87339ba` | — |
| §5-§6 | `9695369` | **这枚之前 staged 里出现过兄弟程的 `140-static-inventory-r1.md`**：我带 pathspec 提交，事后双量两边归属都没被卷走（见 §5.1） |
| §7 | `9c8eb3f` | — |
| §7.1 | `f0d023e` | — |
| §8 | `27a880a` | — |
| §9-§10 | `d484edf` | — |
| §7.1 口径更正 | `9fcb2ed` | 更正"两版数不同口径"那句 |
| 结构修复（本次这一枚） | 由 `git log --oneline -1` 现量 | **修的是我自己造成的两件事**：(a) 逐节 Edit 的锚点选在已挪过位置的文本上 -> §6 一度落到 §10 之后（现已按 0-10 顺序排回，内容零增删）；(b) 一枚 Edit 的 `new_string` 被我发到一半就截断 -> §10.3 那张账表被换成残句（现已补回并顺手写全） |

**10.5 交回前复跑一发（纪律：状态断言会过期；也为了排掉"G1/G2 当时是两枚后台并发跑的"这一枚方法论阴影）**

| 项 | 时刻（`date` 现量） | 读数 |
|---|---|---|
| **G7** = G1 同法同树**单独**跑：`WISP_ENV=test PATH="<snap-post>/third_party/sherpa-onnx:$PATH" go test -count=2 -v ./cmd/wisp/` | `18:27:57` 起、`18:29:59` 止 | 包级 **rc=0**；`=== RUN=202`、`--- PASS=108+94=202`、`--- FAIL=0`、`--- SKIP=0`、`panic` 命中 **0**（日志 `G7-post-cishape-recheck.log`） |
| 被验那枚用例逐名 | 同发 | 两遍都是 `--- PASS: TestAC2EveryLeg...128` + 四枚子用例 + `resolveSecretLayout` 全 PASS，`FAIL/SKIP` 各 0 |
| owner 两棵真实目录（`%APPDATA%\wisp`、`wisp-dev`） | `18:21:59` | 均 **0 条目** |
| 被验那枚码是否仍在历史里（本程未 push） | `18:21` 前后 | `git cat-file -t c2fa2e9` = `commit` |

⇒ **G7 与 G1 逐数相同** => §4 那两发的并发跑没有互相污染；本表所有档位不变。
⇒ 三处引用别人程的地方（`§9#1` 容器 vet、`§6.1` CI 原文、`§7.2` 兄弟程 staged 事件）取的都是**已落盘的字节**
（commit 与我那份 `gh-selffetch.txt`），不随时间腐坏。

**交回的话**：`AC#4` 判**成立、无附条件**；勾由编排者按本表打；本程未翻勾、未改票名、未写票面 128 那枚文件、未 push、未碰 `ci.yml`／生产码／`internal/winsec/**`／`internal/risk/**`／阈值／golden／`frontend/**`／`design/**`／`docs/reports/**`／别人的证据件。