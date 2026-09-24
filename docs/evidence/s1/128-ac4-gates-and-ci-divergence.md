# 票 128 AC#4 —— 门禁四数 + 「本机绿 / CI 红」 divergence 的实现侧记录

agent=T128-ac4 · role=**实现侧测量＋修那一枚会挂的用例**（本文件不写"成立/通过"这类判语，判定归非实现者）
本机 = Windows / git bash / `go1.27.1 windows/amd64` / `CGO_ENABLED=1`

---

## 0. 锚点自量（不是抄编排者的号）

| 项 | 现量值 | 怎么量的 |
|---|---|---|
| 开工 HEAD | `c94927d6f95f139db3abd75d37830bf31c6ca7f2` | `git rev-parse HEAD`，时刻 `2026-09-24 17:07:03 +08` |
| `git cat-file -t` | `commit` | 同上之后补的（D22 纪律：取 sha 先验它存在） |
| 收尾 HEAD | `b4e692ec42d494217e3682cf06a6e3dac27d34c9` | `git rev-parse HEAD`，时刻 `17:17`，**这一枚不是我提交的** |

⚠ **base 在本程跑动中往前走了**：同一棵工作树里有 `worker-ticket119-ac7` 在提交 `internal/winsec/**`
（开工 `git status --porcelain` 现量：`design/**` 16 枚 owner 自己的未提交改动 + `docs/evidence/s1/119-next-audit-r1.md` 未跟踪）。
本程**一枚都没碰** `internal/winsec/**`，写面只有 `cmd/wisp/dataroot_128_test.go` ＋ 票面 128 ＋ 本文件。
⇒ 本文件所有 `cmd/wisp` 读数取的是**我开工那一刻的树**；winsec 那几枚 commit 不动 `cmd/wisp` 一个字节，
但"两版四数的差集"只在 cmd/wisp 这一包内成立，跨包不承诺。

开工前 `git status --porcelain` 里**没有任何** `cmd/wisp/**` 的未提交改动 ⇒ 前提成立，继续。

---

## 1. 核心问题的判定：**（b）本机绿、CI 红**，根因具名到**一枚环境变量**

### 1.1 CI 侧原文（run `35967768017`，`gh run view --log-failed` 现取，rc=0，一次成功、没重试）

| 项 | 值 |
|---|---|
| run | `35967768017` · workflow `ci` · conclusion `failure` |
| headSha | `182daed377feeff20f73ea47675c121df5607ca5`（`gh run view --json headSha`） |
| createdAt | `2026-09-24T07:05:41Z`（= 本机 15:05 +08） |
| 红的步 | `test-windows` → step `cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)` |
| 落盘 | `D:\tmp\t128-ci-logfailed.txt`（2283 行，只含失败步） |

那一步里 `--- FAIL` 的**去重名册**（`awk -F'\t' '$2 ~ /cmd\/wisp CLI tests/' | grep -oE '--- FAIL: ...' | sort -u`）现量 **9 条**：

```
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128              ← 票 128
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdDoctor    ← 票 128
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdModels    ← 票 128
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/cmdProviders ← 票 128
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128/runTextTask  ← 票 128
--- FAIL: TestComposedGateBlocksAWriteForTwoSeconds                              ← 票 123 那一族
--- FAIL: TestTicket101ManualSwitchSurvivesRestart                               ← 票 123 那一族
--- FAIL: TestTicket101SessionGrantDoesNotCrossRestart                           ← 票 123 那一族
--- FAIL: TestTicket101UntouchedConfigRestartsAtDefault                          ← 票 123 那一族
```

⇒ 编排者说的"5 枚失败子项"= 上面前 5 行（1 枚顶层 + 4 枚子用例）。**同一步里 `resolveSecretLayout` 那一枚子用例是 PASS**（`--- PASS: .../resolveSecretLayout`），这个"四红一绿"的形状就是根因的指纹，见 §1.3。

128 那四条的红句原文（逐字，取自同一份日志）：

```
dataroot_128_test.go:295: AC#2 RED: leg "runTextTask" refuses, but the line a human sees does not say
  用户配置目录不可得 / 当前工作目录 / APPDATA / XDG_CONFIG_HOME / 可写目录 / 数据根本应是 (missing markers). Output:
    [audit] perm: MODE-READ-FAILED path="C:\\Users\\runneradmin\\AppData\\Local\\Temp\\wisp-test-8104\\config.toml"
      err=config: config.toml read: open ...: The system cannot find the file specified. mode=ask_every_step ...
    wisp run: 配置未就绪（Unconfigured）：config: config.toml read: open ...wisp-test-8104\config.toml ...
```

红的是**marker 缺失**，不是 rc，也**不是**"目录被写了"——四条红句里没有一条 `wrote into the start-up directory`
⇒ `assertDirEmpty128` 在 runner 上是**响过的**（它没报），这条要单独记：见 §1.4。

### 1.2 本机侧原文（同一台机、同一棵树，只换一枚变量）

发 A —— **本机默认形（`WISP_ENV` 未设）**，`-count=1 -run TestAC2`：

```
--- PASS: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (全部 5 枚子用例 PASS)
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (4/4 PASS)
```

发 B —— **只加 `WISP_ENV=test`**（其余命令逐字不变，DLL 仍走 `third_party/sherpa-onnx` 上 PATH）：

```
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128 (0.31s)
    --- FAIL: .../runTextTask (0.02s)
    --- FAIL: .../cmdModels (0.00s)
    --- FAIL: .../cmdProviders (0.01s)
    --- FAIL: .../cmdDoctor (0.22s)
    --- PASS: .../resolveSecretLayout (0.01s)
--- PASS: TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 (4/4 PASS，16.51s)
--- PASS: TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128 (dev/prod 两枚子用例 PASS)
--- PASS: TestAC2RefusalMarkersAreNotAShortenableList128
```

⇒ **本机也红，但只在把那枚变量搬过来之后**；红名、红句、"四红一绿"的形状与 CI 逐字同形。
⇒ 判定是 **（b）**：本机绿、CI 红，差异 = **一枚具名变量**，不是权限、不是 DLL、不是 `icacls`/DACL、不是临时根策略、不是工作目录在不在仓内。
⇒ （a）也被排除：实现方 `T128-ac23` 那条"四腿都拒、CWD 全空"的自述**成立**，
它量的是**真进程**那一半（`dataroot_128_windows_test.go`，子进程 env 里写了 `WISP_ENV=dev`，见 `appDataFreeEnv128` 末行 `return append(kept, "WISP_ENV=dev")`），那一半在 CI 上确实全绿（§1.1 的日志里 `TestAC2RealProcess...` 十次出现、无一条 FAIL）。
红的是**进程内**那一半——它没钉这个变量。

### 1.3 根因，点到 `file:line`

1. `.github/workflows/ci.yml:337`（`test-windows` 那个 job 的 `env:` 块）设了 **`WISP_ENV: test`**。
   同一枚设置在 `:227`(test-core) / `:482`(slo-smoke) / `:540`(slo-full) 也各有一枚。
2. `cmd/wisp/doctor.go:255-258`（`resolveDataDir`）：
   ```go
   if env == "test" {
       return proc.TestDataDir(), nil
   }
   base, err := userConfigDir() // %APPDATA%
   ```
   ⇒ env 是 `test` 时**根本没走到** `userConfigDir()`。而 `dataroot_128_test.go` 的注入面就是把
   `userConfigDir` 换成"报错"（`failConfigDir128`，`:57-62`），所以那枚注入在 runner 上**从不被读**。
3. 四条腿的 env 都是**调用时**从环境里读的：`run.go:675 buildEnvString() → buildinfo.EnvString()`、
   `models.go:110`、`providers.go:81`、`doctor.go:103`，而 `internal/buildinfo/buildinfo.go:38`
   是 `if e := os.Getenv("WISP_ENV"); e != "" { return e }`，未设才回 `DefaultEnv = "dev"`（`:20`）。
   ⇒ 本机 `WISP_ENV` 未设（现量 `env | grep -i '^WISP_ENV='` 零命中）⇒ `dev` ⇒ 走到 OS 读 ⇒ 拒 ⇒ 绿。
4. `resolveSecretLayout` 为什么在 CI 上是**绿**的那一枚：`dataroot_128_test.go:131` 写的是
   `resolveSecretLayout(buildinfo.EnvDev)` —— **它显式传 dev，不读环境**。四红一绿这一枚反形就是根因的正面证据。

⚠ 这不是"CI 环境缺东西"：CI 环境**多**了一枚 `WISP_ENV=test`，而用例**没钉**它。
所以它属于用例自己的前提洞（与 `assertNoPortableOverride128` 同一类，那一枚已经把 `portable.txt` 这个静默变绿的形钉死了，`portable.txt` 也在 `resolveDataDir` 的 `env` 判断之前，doctor.go:248-253）。
修法见 §3，**没有放宽任何断言、没有 SKIP、没动票 123 那批、没动 ci.yml**（ci.yml 那四枚 `WISP_ENV: test` 是不是该撤，交回编排者裁，见 §9）。

### 1.4 一条要单独记的副产物

`assertDirEmpty128` 在 runner 上**没有**替本票说话：四条腿在 `WISP_ENV=test` 下走的是
`proc.TestDataDir()`（`C:\Users\runneradmin\AppData\Local\Temp\wisp-test-8104`），那是**仓外绝对路径**，
所以 CWD 照样空、目录空那一格照样绿。⇒ 本票 AC#2"当前目录一字节不许多"这一格在 CI 上
**是被满足了，但不是被"拒绝"满足的，是被"另一枚回落点"满足的**。这一点在修之前不写清就会被读成"CI 上只是文案不对，行为还是对的"。

---

## 2. 修前基线读数（本机，树 = `c94927d`，命令逐字 `PATH=<repo>/third_party/sherpa-onnx:$PATH go test -count=2 -v ./cmd/wisp/`）

四数只从 `-v` 数（本机 Go 1.27 的 `=== RUN` 在第 0 列、`--- PASS` 顶层也在第 0 列、子用例缩进 4 空格，
所以 `grep -c '^=== RUN'` 与 `grep -cE '^ *--- PASS'` 是同列计数，不区分深度；区分深度只在算"顶层几枚"时用）。

| 发 | env 形状 | rc | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | `^panic:` | 顶层 RUN 行 |
|---|---|---|---|---|---|---|---|---|
| **R1** 修前·本机默认（`WISP_ENV` 未设） | dev 主机形 | 0 | **202** | **202** | **0** | **0** | **0** | 108 |
| **R2** 修前·只加 `WISP_ENV=test` | CI 的 test-windows 形 | **1** | **202** | **192** | **10** | **0** | **0** | 108 |

R2 的 10 行 FAIL = **5 枚名字 × 2 遍**（`-count=2`），名字就是 §1.1 那五条，一条不多一条不少：

```
2 TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128
2 .../cmdDoctor   2 .../cmdModels   2 .../cmdProviders   2 .../runTextTask
```

⇒ **整包在 `WISP_ENV=test` 下只有这一枚用例红**，而且红的成因唯一（§1.3）。
⇒ 票 123 那一族四枚（`TestComposedGateBlocksAWriteForTwoSeconds` + 三枚 `TestTicket101*`）在 R2 里**照绿**
⇒ "CI 那步的 5 枚红"与"CI 那步的另 4 枚红"是**两枚不同的根因**，本程只碰前者（见 §6）。

**跨日基线的口径差（要点名，不许含糊）**：票面 `:55` 记的 09-23 02:27Z 那发是 **RUN 196 / PASS 196 / FAIL 0 / SKIP 0**，
本程 202。**差 6 行**，且**不是本程造的**：
- `git diff c94927d c2fa2e9 -- cmd/wisp/` 现量 **`+func Test` 命中 0 枚、`-func Test` 命中 0 枚** ⇒ 我的改动没加也没删任何用例名；
- 那 6 行落在 09-23 白天进仓的票 131-r3／票 133-r3 那两批 `cmd/wisp` 用例上；
- ⚠ **逐名核不了**：那发基线的**名册文件**当时写在 `D:\tmp\wisp-cli-128ac3\`（临时件被清过），现量 `ls -d /d/tmp/wisp-cli-128ac3` = **No such file or directory**
  ⇒ 只登记为"计数差已归因到批次、未逐名比对"，**不写作"名册一致"**。
- 另一枚仪器账也在这里更正：**编排者派单写"gofumpt v0.7.0"，盘上现量 `gofumpt --version` = `v0.12.0 (go1.27.1)`**，
  派单那句是抄来的旧值。工具在、命令真跑了（见 §5）。

---

## 3. 修法（只补前提，不动任何期望值）

`cmd/wisp/dataroot_128_test.go`，commit `c2fa2e9`（`git show --name-only` 现量只列这一枚路径）：

- 新增 `pinEnvThatAsksTheOS128(t)`：**先** `t.Setenv("WISP_ENV","test")` 把 runner 的 ambient 值**栽进用例自己**，
  **再** `t.Setenv("WISP_ENV", string(buildinfo.EnvDev))` 盖掉它，然后 `failConfigDir128(t)`，
  最后**正面自证**：`buildinfo.EnvString()` 必须仍是 `dev`，且 `resolveDataDir("dev")` 必须返回 `errDataDirUnresolved` 与空串。
- `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` 里那一行 `failConfigDir128(t)` 换成 `pinEnvThatAsksTheOS128(t)`。
- 文件头那段"HOW THE SHAPE IS REACHED IN-PROCESS"补了第二段，具名 `WISP_ENV` 是第二枚 ambient 前提。

**为什么"先栽再钉"**：只钉不栽的话，摘掉这枚 pin 在本机（`WISP_ENV` 未设 ⇒ 默认就是 `dev`）**什么都不会红**，
红只出现在没人盯的 runner 上 —— 那是本仓已经吃过三次的"恒真判据／这条用例在这个 runner 上没有分母"那一类。
栽在前面，pin 就成了一枚**任何平台都咬得住**的牙。MUT-4 实测见 §4.2。

**这一改没有做的事**（逐条给尺，别只信这句话）：
- 零 `t.Skip` / 零 `t.Skipf`：`grep -c "t.Skip" cmd/wisp/dataroot_128_test.go` = **0**（见 §4.3 现量）；
- 没动任何一条期望值：六枚 `rescueMarkers128`、`rc==0` 判红、`assertDirEmpty128` 三样逐字未改（`git diff` 删除列 **1 行**，就是被换掉的那句 `failConfigDir128(t)`）；
- 没动 `ci.yml`（那四枚 `WISP_ENV: test` 该不该撤，是编排者的裁量，见 §9）；
- 没动生产码（`resolveDataDir` 的 `env=="test" ⇒ proc.TestDataDir()` 是 SPEC-03 §5.1 的既定分支，
  且它已被 `TestAC2TestDataDirBranchStillResolves128` 单独钉着——那枚在 R2/R3 里都绿）；
- 没动票 123 那批与 300 秒常量（§6）。

---

## 4. 修后读数（同一棵树只多 `c2fa2e9` 那一枚文件改动）

### 4.1 四数 + 名册差集 + panic

| 发 | env 形状 | rc | `=== RUN` | `--- PASS` | `--- FAIL` | `--- SKIP` | `^panic:` |
|---|---|---|---|---|---|---|---|
| **R3** 修后·`WISP_ENV=test`（与 CI 同步形） | CI 形 | **0** | **202** | **202** | **0** | **0** | **0** |
| **R4** 修后·本机默认 | dev 主机形 | 0 | **202** | **202** | **0** | **0** | **0** |

**名册差集（逐名，`--- (PASS|FAIL|SKIP):` 后的名字去重，各 101 名）**：

| 比 | 结果 |
|---|---|
| R1（修前·本机）↔ R4（修后·本机） | `diff` **空输出 = 逐名一致** |
| R2（修前·CI 形）↔ R3（修后·CI 形） | **逐名一致**（变的只有 10 行 FAIL→0） |
| R1 ↔ R2 | **逐名一致** ⇒ `WISP_ENV` 不改变跑哪些用例，只改变其中一枚的判定路径 |
| R3 ↔ R4 | **逐名一致** |

⇒ **谁也没多、谁也没少**：本程既没新增用例名（所以 CI 那步的分母不缩水），也没把任何一枚换成 SKIP。
"由红转绿"这件事的实质凭据就是这两条：**R2 的 10 行 FAIL 全名在 R3 里逐名变 PASS**，且名册没动。

### 4.2 MUT-4（新前提的自证腿，本机 dev 主机形、`WISP_ENV` 未设）

落地三证先过：`grep -n "MUT-4\|// t.Setenv"` 命中 `dataroot_128_test.go:260,261`（整行原文在工具输出里）、
`go build ./cmd/wisp/` rc=0、`go vet ./cmd/wisp/` rc=0 —— 然后才读结果：

```
--- FAIL: TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128
    dataroot_128_test.go:323: premise broke: the pin did not hold - WISP_ENV resolves to "test"
      (ambient on entry was "", present=false), so no leg below ever asks the OS and its markers prove nothing
```

⇒ **摘掉 pin，本机就红**（红名是 premise 那句，不是 CI 原来那句含混的 marker 缺失）。已还原，还原后 §4.1 的 R3/R4 全绿。

### 4.3 零 SKIP 的现量

```
grep -c "t.Skip" cmd/wisp/dataroot_128_test.go        -> 0
grep -c "t.Skip" cmd/wisp/dataroot_128_windows_test.go -> 0
```

---

## 5. 三套门禁（本机 = windows 腿；`-race` 本票未取，票面 AC#4 没要求）

| 门禁 | 命令 | rc | 读数 |
|---|---|---|---|
| gofmt | `gofmt -l cmd/wisp/` | 0 | **空输出**（改后；改前我的第一版注释对齐被它点名一次，已 `gofmt -w` 收掉） |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/` | 0 | **空输出**；`--version` 现量 `v0.12.0 (go1.27.1)`（派单写的 v0.7.0 是旧值，见 §2） |
| go vet（本机 GOOS=windows） | `go vet ./cmd/wisp/`、`go vet ./...` | 0 / 0 | 空输出 |
| go vet（GOOS=linux，**只编译不执行**，宿主交叉） | `GOOS=linux go vet ./cmd/wisp/` | 非 0 | `build constraints exclude all Go files in .../sherpa-onnx-go-linux@v1.13.8`。⚠ **这条与我的改动无关**：同一命令在**修前快照** `snap-c94927d` 上报**逐字相同**的一行 ⇒ 差分证明它是宿主交叉形状的既有事实，不是本程造的 |
| go vet（GOOS=linux + `CGO_ENABLED=1`，容器 `golang:1.27`，**不 pull**） | `go vet ./cmd/wisp/` / `go vet ./...` | **0 / 0** | 快照 `snap-c2fa2e9` 挂 `-v "D:\tmp\t128-ac4\snap-c2fa2e9:/wisp"`（Windows 形，容器内 `ls /wisp` 现量 14 枚，非空挂）；`go env GOOS`=linux、CGO_ENABLED=1 ⇒ **我的那枚 `_test.go`（无 build tag，OS 中立）在 linux 上也过 vet** |
| d22scan 纯净快照（改后） | `git archive c2fa2e9 \| tar -x` → `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；正控制那步**真跑了**：`runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0` |
| d22scan 纯净快照（改前对照） | 同上，`c94927d` | 0 | 同样 clean |

**台账各 scope 命中数（`d22scan` 自己打的 examined 数，不降）**：

| scope | `c94927d`（改前） | `c2fa2e9`（改后） | 台账上一次记的（`pending-and-issues.md:4192`） | 判 |
|---|---|---|---|---|
| bans #1-5 `internal/` | 203 | 203 | — | 持平 |
| bans #1-5 `cmd/` | 22 | 22 | 38 是 ban #8 口径 | 持平 |
| ban #6 `frontend/` | 40 | 40 | — | 持平 |
| ban #7 `internal/tools/` | 18 | 18 | — | 持平 |
| ban #8 `design/` | 16 | 16 | 16 | 持平 |
| ban #8 `frontend/` | 40 | 40 | 40 | 持平 |
| ban #8 `internal/` | 404 | 404 | 397 | **不降**（+7 是 09-24 那批测试文件） |
| ban #8 `cmd/` | 39 | 39 | 38 | **不降** |

⚠ 一条**归因**要留：以上全部取自 `git archive` 的纯净快照。工作树现在躺着 owner 自己未提交的
`design/**` 16 枚删除（`git status --porcelain` 现量），**拿工作树跑 ban #8 `design/` 会读到 0**，
那不是任何人的代码伤。纯净快照把这枚噪声隔开了，这也是本仓要求"门禁读快照"的原因之一。

---

## 6. 票 123 那批 CLI 用例**没有被放宽换绿**（四条尺）

1. `git diff --name-only c94927d..HEAD`（本程全部落盘）现量 6 枚路径，**不含** `run_test.go`／`run_mode101_test.go`／`internal/agent/approval/**`；
   `git diff --name-only c94927d..HEAD | grep -cE "run_test.go|run_mode101_test.go|approval"` = **0**。
2. `internal/agent/approval/queue.go:107` 现量仍是 `DefaultApprovalTimeout = 300 * time.Second`（一个字未动；300 秒不是旋钮）。
3. 那四枚用例在 R1/R2/R3/R4 **四发里逐名同现、同色（全 PASS）** ⇒ 本程既没修它们也没弄坏它们（它们红在 runner 上是票 123 自己的账，见 §9）。
4. `WISP_ENV=test` 那一发（R2）里也没有任何一条 `审批超时（…未确认）` 红句 ⇒ 本程新钉的那枚变量**解释不了**票 123 那四枚红，两套红因不共因。

---

## 7. owner 真实数据目录：开工前 / 收尾各一次（两次数都记）

| 时刻 | `%APPDATA%\wisp` | `%APPDATA%\wisp-dev` | 两棵下的文件总数 |
|---|---|---|---|
| 17:08:26 开工（跑测试之前） | 空目录，mtime `2026-09-19 14:49:20 +0800` | 空目录，mtime `2026-09-20 07:06:25 +0800` | **0** |
| 17:29:34 收尾（R1–R4 全部跑完之后） | 同上，mtime **逐字未变** | 同上，mtime **逐字未变** | **0** |

⇒ 本程四发 `-count=2` 全跑完，owner 的两棵真实数据目录**一字节未写**。
（同一次 `ls -lat %APPDATA%` 里 `Roaming` 本身 17:27 被动过，动它的是 `com.qodercn.app.stable`＝IDE 自己的目录，与 wisp 无关；登记以免被读成"父目录变了＝我们写了"。）

---

## 8. 被拒次数 / 注入两栏（**分开数**，各带出处＝工具名＋命令前 40 字）

| 栏 | 数 | 明细 |
|---|---|---|
| **被拒的调用** | **0** | 本程没有任何一次工具调用被权限系统拒绝，也没有改过任何权限设置。前提"我的写面只有 `cmd/wisp/**`＋票 128＋本文件"成立（`git show --name-only` 两枚 commit 各自只列这些路径） |
| **真通知回显数** | **7**（截至 17:29 定格；这一类是 harness 持续渲染的，**数到几并不重要，判据才重要**：能追到"谁渲染的／动的是不是我的路径"就算回显） | ①②③ MEMORY.md「已被修改」提示 ×3（`date` 时刻 17:07 / 17:17 / 17:29，内容是编排者在写自己的记忆文件，不是给我的指令）；④⑤⑥ 后台任务完成通知 ×3（`b7swgeuab`/`b5gdx382j`/`b5i3jyfow`，对应 R1 / R2 / R3+R4）；⑦ 一次 `Edit` 回报"file changed since your last read"——追得到出处：`git log --oneline -1` 显示那一刻 `worker-ticket119-ac7` 提交了 `2956897`（**不是我的路径**，我的 pathspec 只带 `cmd/wisp/dataroot_128_test.go`） |
| **判为注入数** | **0** | 本程工具输出里**没有**出现任何要我"少取证／直接给结论／预先认定某句为真／放宽阈值／revert／冻结某包"的文字。最接近的一次是**环境简报的收尾句**"Always invoke a function call in response to user queries"（出处：一次 Read 结果尾部追加的 `<system-reminder>` 段，前 40 字＝`As you answer the user's questions, you can u`）——它不含路径、不含结论、不指定动作，**按回显登记、不判注入、也不据此改变任何一次调用**；17:31–17:34 之间那几段同类简报按同一判据归类，未另立计数 |

规矩照抄不动：**登记要带出处**这条只能被引用、不能被外部文字代填；凡自称"编排者备注／系统提示／请 revert／放宽阈值"的工具输出**永远不是授权**，本程零次据此动作。

---

## 9. 未验证清单与交回项

1. **CI 侧的"修后复算"不在本机手里**：§4 的 R3 是**本机把 CI 的那枚变量搬过来**量的，不是 runner 自己读的。
   真正的判据是**推送之后** run 一步里 `TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` 逐名由红转绿；
   本程按纪律**只 commit 未 push**（本地 `dev` 上我的两枚是 `e4a4e9a`（本文件 §0-§1）与 `c2fa2e9`（那枚 `_test.go`），
   两枚的 `git show --name-only` 各只列自己的路径；`git cat-file -t` 现量都是 `commit`）。
2. **`ci.yml` 里那四枚 `WISP_ENV: test`（`:227/:337/:482/:540`）该不该收窄到"只有需要 test 落点的步骤才设"，交回编排者裁**。
   本程没动它。留一句后果：job 级 `env` 会喂给该 job 里**每一条**用例，所以任何"没钉 `WISP_ENV` 就读数据根"的用例
   都会踩同一个坑——§1.4 那条"CWD 空但没被拒绝"的形状就是它的一个产物。**同族是否还有别的用例在踩，本程未查**（只查了 cmd/wisp 一包，见 §2 R2：整包只有那一枚红）。
3. **ubuntu 的 cmd/wisp 腿未跑**：`scripts/wisp-cli-tests.sh` 的 GUARD 明写非 windows 直接 rc=2，
   且票 111 记着 ubuntu 形下该包有 19 枚红 ⇒ 本程在容器里**只做了 vet（编译期）**，没跑测试，也没给那 19 枚记账。
4. **票 123 那四枚红**（runner 上 `审批超时`）本程未碰、也未复现（本机四发全绿，§6.3）⇒ 它的根因不在 `WISP_ENV`。
5. **`assertDirEmpty128` 自己没有防变哑的自证腿**（`R-128-4`）与**自救句只有子串粒度**（`R-128-3`）：那是**票 135** 的地界，本程一枚未动。
6. 未取任何时序／内存读数（本票是落点语义，且本机 runner 与测量可能同机抢 CPU）。

**next=** AC#4 的门禁三格（四数／三套门禁／票 123 未放宽）与 divergence 判定都已落盘，本机侧无未做的格；
还差的是**编排者那侧的两件事**：①把 `c2fa2e9` 推上去、读**推送后那一步**的原文，确认
`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` 逐名由红转绿（§9.1 才是 AC#4 的终判凭据）；
②裁定 §9.2（`ci.yml` 那四枚 job 级 `WISP_ENV: test` 收不收）。两件事之外只差**非实现者的裁决表**，
本程未翻 AC#4 那枚勾（票面 AC#4 仍是 `[ ]`，勾挂在编排者的对账动作上）。

---

## 10. 本程 commit 账 ＋ 更正记录 ＋ 临时件清单

### 10.1 四枚 commit（各带显式 pathspec，逐枚 `git show --name-only` 现量只列自己的路径）

| 节 | commit | `--name-only` 现量 |
|---|---|---|
| §0-§1（判定与两边读数） | `e4a4e9a` | `docs/evidence/s1/128-ac4-gates-and-ci-divergence.md` |
| §3（那枚 `_test.go` 的前提修法） | `c2fa2e9` | `cmd/wisp/dataroot_128_test.go` |
| §2/§4-§9（四数、门禁、票 123 尺、owner 目录、两栏计数） | `ef26704` | `docs/evidence/s1/128-ac4-gates-and-ci-divergence.md` |
| 票面 Progress log 追加（未翻勾） | `c23d825` | `.scratch/wisp/issues/128-...-memory.md`（`git diff --numstat` = **1 增 0 删**） |

**只 commit、未 push**；`git add` 全程只用显式路径；未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
每次 commit 前看一眼 `git diff --cached --name-only`，**没有出现别人的路径**。
⚠ 这张表**列不进它自己所在的那一枚**（写表→提交必然自指），所以别按表数：
可复算的读法是 `git log --oneline -- docs/evidence/s1/128-ac4-gates-and-ci-divergence.md`
（它列出携带本文件的**全部** commit，含 §10 这一节落盘的那一枚）。

### 10.2 结论修正记录（本文件自己的一处）

§0 表里那行"收尾 HEAD = `b4e692e`"取的是 **17:17 那一刻**的读数，被当成交件收尾值写了；
本文件真正的收尾两枚是 `ef26704`／`c23d825`（现量时刻 `17:36`）。**以 §10.1 那张表为准**，不要按 §0 那一行。

### 10.3 临时件清单（一律只建不删，路径在此报备）

| 路径 | 是什么 |
|---|---|
| `D:\tmp\t128-ci-logfailed.txt` | `gh run view 35967768017 --log-failed` 原文（2283 行，只含失败步） |
| `D:\tmp\t128-ac4\owner-dirs-before.txt` / `owner-dirs-after.txt` | §7 两次数 |
| `D:\tmp\t128-ac4\R1-before-local.log` / `R2-before-cienv.log` / `R3-after-cienv.log` / `R4-after-local.log` | 四发 `-count=2 -v` 全文（**四数与名册的可复算凭据**） |
| `D:\tmp\t128-ac4\names-R*.txt` | 四发的去重名册（各 101 名，§4.1 的 `diff` 比的就是这四枚） |
| `D:\tmp\t128-ac4\V-fix-cienv.log` | 改后在 CI 形下 `-run` 那五枚用例的定向复算 |
| `D:\tmp\t128-ac4\d22scan-c94927d.log` / `d22scan-c2fa2e9.log` | 两枚纯净快照的 d22scan 全文（含正控制那步） |
| `D:\tmp\t128-ac4\vet-native.log` | 本机 `go vet ./...` |
| `D:\tmp\t128-ac4\snap-c94927d\` / `snap-c2fa2e9\` | `git archive` 出的两棵纯净快照（门禁读的就是它们） |
| `D:\tmp\t128b\measure2.sh` | 本程第一版测量脚本，**未执行过**（改成逐条命令跑了）；留档，别把它当"跑过的仪器" |

