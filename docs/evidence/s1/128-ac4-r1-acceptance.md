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
