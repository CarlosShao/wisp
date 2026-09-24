# 140 — 静态盘点 r1（票 140 `AC#1` 的「先量面」那一半的**静态前置**；一个测试都没跑）

**取证代理**：`auditor-ticket140-static`（只读盘点程，**非实现者**）
**时刻**：2026-09-24 17:4x–17:5x +08（`date` 现量 `2026-09-24 17:54:52 +0800`）
**码面锚点**：`744d79e781789d368f73a67ac1b485523f9e7b59`（`git rev-parse HEAD` 自量；`git cat-file -t` = `commit`）
　⇒ 本文引用的每一行都取自检过的两个来源之一：① `git show 744d79e:<path>`（钉版，码面以此为准）；
　② 工作树上的 `.github/workflows/ci.yml`（我只读它，且它**不**在两枚在飞代理的地界内）。
**在飞地界回避**：`internal/winsec/**` 与 `cmd/wisp/**` 的未提交改动**一律不算事实**——那两棵目录我全部走 `git show 744d79e:`。
　`internal/winsec/**` 那几枚文件我引用的是**码**（`resolve.go`），不是别的代理可能改过的测试文件；
　`internal/winsec/dataroot_symlink_119_other_test.go` 一枚我在下面点名时同样只报 `git show` 的行。

**本程做了什么 / 没做什么**
- 做了：`git grep`（钉版）、`git show`、`sed`/`cat -n` 读文件、`git ls-tree`、`git rev-parse`、`git cat-file -t`、`date`、`ls`。
- **没做**：`go test` / `go build` / `go vet` / `gofmt` / `docker` / 任何写操作（除本文件与它自己那一枚 commit）。
  票面 `AC#1` 明写「取证方式**必须包含真跑**：同一枚命令在 `WISP_ENV` 未设与设为 `test` 两形各一发，给四数＋名册差集」——
  **那一半本程没有做，也做不了**（简报硬约束：两枚代理正在容器/包内取读数）。本文件是它的**静态前置**，不是它的替代。

**证据档位（逐节标在节标题上）**
- 〔独立复现〕= 本程自己现量的命令＋输出（钉版码面也算，因为我自己在盘上取了那一版）。
- 〔日志＋归档抽验〕= 引用别程读数的原文，**且**我核过它引用的码面在钉版上仍然对得上。
- 〔仅自述不背书〕= 只有别程说过、本程无法用静态手段复算的。

**票面**：`.scratch/wisp/issues/140-job-level-wisp-env-test-makes-the-ask-the-os-hardening-legs-never-get-consulted-on-the-windows-runner.md`
　—— 本程 17:4x 读它时它在盘上是**未跟踪**文件；`git log --diff-filter=A` 现量它由 **`79cfa1d` `docs(A175; 立票 140)`** 入库（在我取数期间，非本程所为）。
**本文件不是那份票面**：不改票面、不改任何既有文件（Rules 第 1 条那半句"每完成一格往票面 append 一条 log"本程**未做**——简报硬约束写死"不许改任何既有文件"，两者冲突时按简报走，请编排者自行决定由谁翻那一格）。

---

## §0　票面核对：问题问得对不对，与票 98／111／123 重不重复　〔独立复现〕

### 0.1 简报与票面的五条前提，逐条现量

| 前提 | 现量命令 | 结果 |
|---|---|---|
| `ci.yml` 有四枚 job 级 `WISP_ENV: test` | `grep -n 'WISP_ENV' .github/workflows/ci.yml` | `:227` `:337` `:482` `:540` ＋ `:250` 一枚**步骤名**。四枚 job 级 = 成立，行号 = 简报给的行号，逐字一致 |
| 那四枚分别是哪四枚 job | `git show 744d79e:.github/workflows/ci.yml` 反查 job 头 | `:224 test-core`(ubuntu-latest) / `:334 test-windows`(windows-latest) / `:479 slo-smoke`(windows-latest) / `:537 slo-full`(**`[self-hosted, wisp-slo]`**) |
| `doctor.go:256-259` 是那个早退形状 | `git show 744d79e:cmd/wisp/doctor.go \| sed -n '240,290p'` | `:256 if env == "test" {` / `:257 return proc.TestDataDir(), nil` / `:258 }` / `:259 base, err := userConfigDir() // %APPDATA%`。成立，行号逐字对 |
| 「四红一绿」里那枚绿的是显式传 `EnvDev` 的第五腿 | `git show 744d79e:cmd/wisp/dataroot_128_test.go \| sed -n '138,144p'` | 在盘上 **`:139` `l, err := resolveSecretLayout(buildinfo.EnvDev)`**。⚠ 台账 `A175` 已经记过：实现程引的是 `:131`，差 8 行（从改前的树引的行号）。**本文按钉版 `:139` 引** |
| 票 128 那枚 helper 真存在且形状如述 | `git show 744d79e:cmd/wisp/dataroot_128_test.go \| sed -n '257,273p'` | `pinEnvThatAsksTheOS128` 在 `:257`，`t.Setenv("WISP_ENV","test")` 在 `:260`、pin 回 dev 在 `:261`，自证 pin 生效在 `:263-264`。成立 |

⇒ **简报没有一条前提不成立**。我不需要报回否定式断言。

### 0.1b 盘上事实：本程工作期间 HEAD 动过，但**码面一字未动**　〔独立复现〕

我把这份文件写完后复查：`git rev-parse HEAD` 已从 `744d79e` 前进到 `d9f008d`（兄弟程在交裁决表）。
**关键复量**：

```
git diff --name-only 744d79e HEAD -- cmd/wisp internal/winsec internal/proc .github/workflows/ci.yml scripts
  → （空）
git merge-base --is-ancestor 744d79e HEAD → YES
```

⇒ 本文引用的**每一个 file:line 在两版上都逐字成立**（`744d79e` 与 `d9f008d`），包括被两枚在飞代理占着的
`cmd/wisp/**` 与 `internal/winsec/**`——它们到这为止**没有产生任何 commit**。

**第二见证（不抵账）**：`d9f008d`（`acceptor-ticket128-ac4-r1` 的裁决表开头）独立重量了我 §0.1 那三条静态前提，
并裁「`:131` 是它取数那版的正确号（`c2fa2e9` 之上加了 8 行注释）」与「`doctor.go:255-258` 那一处才是真 off-by-one，正确是 `256-259`」。
⇒ **与我 §0.1 现量一致**，但按「同格双份表＝勾挂主表、第二见证不抵账」的规矩，我保留我自己那三行命令作为本文的凭据，不把它的表当我的。

### 0.2 问题问得对不对

问对了**主干**，但票面 `AC#1` 的措辞比事实宽：**「同一枚 job 级变量还在哪些包里把加固腿变成装饰」——现量答案是「不在任何别的包里」，因为除 `cmd/wisp` 外没有任何生产码读这枚变量**（凭据见 §1.4，三条独立 grep）。
所以这张票真正剩下的量，不是"再找几个包"，而是这两件我在盘点里撞到的东西：

1. **一枚从没被任何测试断言过的规范义务**：`SPEC-11 §8` 要求 CI 显式断言 `WISP_ENV=test` 生效（原文见 §2.4），而今天那枚**名字写着这件事的步骤**（`ci.yml:250`）走的码路**根本不读这枚变量**。
2. **三枚非测试门禁其实靠这枚 job 级 env 才不落进 `%APPDATA%`**（`ci.yml:400` / `:497` / `:560`，全部经由 `scripts/build.ps1:161-163` 的 `wisp.exe doctor`）。⇒ 这直接决定 `AC#2` 的 (a) 支**不是一句"去掉就行"**，见 §4.2。

### 0.3 与票 98／111／123 重复吗——引的是它们自己的句子，不引票号

| 被点名的三枚（盘面 Status **全是 open**，见 §0.4） | 它自己票面的原文标题句 | 与本票的关系 |
|---|---|---|
| 票 98 | `# 98 — \`go test ./cmd/wisp/\` 在这台机器上**永远测不到任何东西**（加载期 \`sherpa-onnx-c-api.dll not found\`）⇒ 宿主包的门禁是一架空仪器`；Type 行：`门禁完整性（票 71 / A44① / 票 96 同族：**"没跑过"与"跑过且没问题"长得一样**）` | **不重复**：98 说"进程起不来＝零条用例跑过"。本票说"进程起来了、用例跑了、断言也绿了，但它断言的不是它声称的那条分支"。98 那枚洞**实质**上已被票 111 `AC#4` 补上（我 §3.1 现量到 `ci.yml:422` 那一步真在跑 `--scope=cli`），但 98 的票面**并未勾完**（§0.4） |
| 票 111 | `# 111 — CI 只测 33 个包里的 **20 个**：还有 5 个带测试的包零覆盖（\`internal/ball\`/\`cmd/wisp\`/\`internal/perm\`/\`internal/plugin\`/\`cmd/llmrecord\`）`；判据句：`配置里有一行 ≠ 它给过结论` | **不重复但同族，且是它的下一个粒度**：111 的分母是**包**。我在 §3.4 量到一枚**文件级零分母**的漏（`cmd/wisp/secret_dataroot_119b_test.go` 是 `//go:build !windows`，而 `cmd/wisp` 只在 `test-windows` 跑），`portable-tests.sh` 的 GUARD A 是包粒度的，抓不到它。**那一枚该记在 111/98 的"有没有腿"账上，不是 140 的账** |
| 票 123 | `# 123 — CI 的 runner 上 4 枚 \`cmd/wisp\` 用例全因"**审批超时（1/300 秒未确认）⇒ C18 一律判拒绝**"而红；本机 33 条能全绿，因为本机有人（或有人造的确认腿）` | **不重复**：123 的隐性前提是**交互审批**，不是环境落点分叉；形状还相反（123 是"该绿的在 CI 上红"，本票是"该红的在 CI 上绿"）。票 140 `AC#4` 那句「票 123 那批 CLI 用例……**不许被放宽换绿**——300 秒不是旋钮」是**边界声明**，不是并案 |

⇒ **本票不与那三枚重复，判据成立**。它与票 128 的关系也如票面所述：128 `AC#4` 只管那一包的四数与归因（`docs/evidence/s1/128-ac4-gates-and-ci-divergence.md`，见 §5），本票管范围。

### 0.4 **简报与票面各有一处前提要更正**：那三枚票**在盘面上都没结案**　〔独立复现〕

简报说「和**已结案**的票 98／票 111／票 123 重复吗」；票面 `:5` 的「地界」行也写着「与**已结案**的票 98……同族但不同问题」，
台账 `A175`（`docs/reports/pending-and-issues.md:5498`）是同一句话的出处。现量三枚票面本身：

```
98  -cmd-wisp-tests-never-run-on-this-host.md            -done 后缀=0  Status: open  [x]=0 未勾 AC=5
111 -ci-tests-20-of-33-packages.md                       -done 后缀=0  Status: open  [x]=0 未勾 AC=10
123 -cmd-wisp-cli-tests-assume-a-human-approver-on-ci.md -done 后缀=0  Status: open  [x]=1 未勾 AC=6
（对照：本票 140 自己是 -done 后缀=0 / open / [x]=0 / 未勾 4 —— 与它们同形，所以这个测量本身不含"谁结没结"的判断）
```

`grep '票 98.*结案\|R-98-'` 与 `'票 111 …结案'` 在台账里**没有找到结案登记**，只找到 `A175` 那句把 98 称"已结案"的话；
`:1286` 反而明写「**票 111 的 AC#6/AC#9 既不能结案、也不能被打回**——它欠的是样本，不是红」。

⇒ **怎么读这条，两说，我不敢单方面定案**：
按 `AGENTS.md §1.5`（`-done` 是防重领的唯一键）与票面自身，三枚都还开着；
按编排者的书面判断，98 的那枚洞已经被票 111 `AC#4` 用 `scripts/wisp-cli-tests.sh` 补上（我 §3.1 现量到那一步真的在跑），所以"实质上不再欠读数"这一支**成立**。
⇒ 对本票的**实际影响只有一个**：**§3.4 那枚文件级零分母不能写成"票 111 已结案后的漏网"，得写成"票 111 的 AC 还没勾完、这条正好归它"**。其余判据不受影响。

---

## §1　Q1：全仓「会去问操作系统拿答案」的入口清单（生产码／测试码分栏，按包分组）　〔独立复现〕

### 1.0 先定判据（这是本题能不能答对的关键）

一枚入口属不属于本票，看**两件事**，缺一不算：

- **腿型 A｜环境早退型**：存在一条 `if env == "test" { return … }`，**return 掉了后面那次 OS 读**。
  ⇒ `WISP_ENV=test` 时那次 OS 读**根本不执行**，注入的 seam／真的拒绝分支**从不被咨询** ⇒ **这条腿会变成装饰**。
- **腿型 B｜环境改变落点但不早退型**：OS 读**无条件执行**，env 只决定读来的值往哪放／要不要用。
  ⇒ 两形都咨询闸门，只是答案不同 ⇒ **不是本票的靶**（但仍是"落点被搬走"，写进 §4.2 算 (a) 的代价）。

另一条**反向排除**（很重要，别把名单灌水）：
`WISP_TEST_DATA_DIR`（`internal/proc/envfork.go:36`）是**另一枚变量**，CI 的四枚 job **都不设它**。
按 `envfork.go:100-106` 的判别（原文："`WISP_TEST_DATA_DIR` is an identity contract … returned verbatim"、
"`os.TempDir()` and `os.UserConfigDir()` are an answer to a question this process asked"）：
凡是**测试自己声明根、再拿声明去比返回值**的用例，它压根没问操作系统 ⇒ 与 `WISP_ENV` 无关，不进名单。
`t.TempDir()` 同理（131 个文件、421 处命中——它们全是 harness 句柄，不是"问操作系统要决策"）。
**再排除一枚容易误纳的**：`cmd/wisp/logsink.go:144-147 installLogSink(dataDir string)` 的 `if dataDir == ""` 是一枚
**"没有根就拒"、不问操作系统**的守卫（原文 `:142-143`："An empty dataDir is a refusal, not a fallback"）。
`cmd/wisp/leg_sink_gate_131_test.go:1315` 把它和 `os.UserConfigDir` 并列提（`var x = installLogSink` vs `var x = os.UserConfigDir`），
那是在讲**门禁怎么认变量初始化形状**，不是在讲同一类落点分叉——别因为那句注释把它算进名单。

### 1.1 生产码

**`cmd/wisp`（全仓**唯一**读环境态 `WISP_ENV` 的包）**

| 入口 | file:line | 腿型 | 说清楚 |
|---|---|---|---|
| `resolveDataDir(env string)` | `cmd/wisp/doctor.go:247`；早退 `:256-258`；OS 读 `:259`；拒 `:261` | **A（早退）** | 本票**唯一**一枚形状。`portable.txt` 分支 `:248-256` 还更早，也是一枚早退（`assertNoPortableOverride128` 就是为它造的，`dataroot_128_test.go:76-83`） |
| `userConfigDir` 那枚 seam 变量 | `cmd/wisp/doctor.go:277` `var userConfigDir = os.UserConfigDir` | — | 生产里没有任何东西重新绑定它；唯一绑它的地方是测试 `dataroot_128_test.go:67-69` |
| `resolveSecretLayout(env)` | `cmd/wisp/secret.go:129`；OS 读 `:130`；`SealableRoot` `:148`；`sessionLayout` `:106` | **B（不早退）** | env 只是 `LayoutFor` 的入参。**这正是票 128 那枚"绿掉的第五腿"绿的原因**：`dataroot_128_test.go:139` 传的是字面量 `buildinfo.EnvDev` |
| `runTextTask` / `cmdModels` / `cmdProviders` / `cmdDoctor` | `run.go:156`（经 `run.go:675 buildEnvString()`）、`models.go:110`、`providers.go:81`、`doctor.go:104-105` | **A 的四个消费者** | 四者都是 `if x.dataDir == ""` 才去解析（`run.go:155`、`models.go:109`、`providers.go:80`）⇒ 测试只要注入 `dataDir` 就与 env 无关 |
| `runResident()` | `cmd/wisp/resident_windows.go:27`（`ResolveEnv`）→ `:35 proc.DefaultLayout(env)` | **B** | 走 `internal/proc`，OS 读无条件。**票面 §"事实"没提它，但它正是"另一形"的现存活例子** |
| `wisp slo` | `cmd/wisp/slo_windows.go:238-240` `if os.Getenv("WISP_ENV") == "" { os.Setenv(...,"test") }`；`:243` `ResolveEnv` | **B + 自带兜底** | ⚠ `AC#2(a)` 的关键：**这枚门禁在 env 未设时自己就设成 `test`**，所以它对 job 级 env 是**幂等**的，不是依赖 |
| `printVersions` | `cmd/wisp/main.go:135` | 只读、只印 | `doctor.go:42-43` 同样把 `EnvString()` 印进 INFO 行 ⇒ **(a) 会改变 `wisp doctor` 的打印内容**，见 §4.2 |

**`internal/buildinfo`** — `env.go:35-40 ResolveEnv()`、`buildinfo.go:37-42 EnvString()`；`buildinfo.go:20 DefaultEnv = "dev"`（"WISP_ENV when the env var is unset"）。
⇒ **腿型：无**（它不问操作系统要路径）。但它是 (a) 的**方向来源**：env 未设 ⇒ 落回 `DefaultEnv`，而 CI 的构建是 `build.ps1:107 -X …DefaultEnv=$Env`，三处都传 `-Env dev`（`ci.yml:400`、`:497`、`:560`）⇒ **未设 ＝ dev，不是 prod**。

**`internal/proc`** — `envfork.go:236 dir, err := os.UserConfigDir()`（在 `DefaultLayout`，**env 是入参**）；`envfork.go:125 filepath.Join(SealableRoot(os.TempDir()), "wisp-test-<pid>")`（在 `TestDataDir`）；`envfork.go:55 LayoutFor` 的 `case buildinfo.EnvTest` 在 `:77-85` 返回 `TestDataDir()`、**丢弃** 注进来的 `userConfigRoot`。
⇒ `LayoutFor`／`DefaultLayout`／`Summarize`（`:291`）**都不读环境态 `WISP_ENV`**。腿型 **B**。`boot_windows.go:64` 校验 env 枚举、`:77` 才 `DefaultLayout`。

**`internal/winsec`** — `resolve.go:242-244 resolverProbeRoot() → resolveProbeRoot(os.TempDir())`，调用点 `:198`、`:340`；`winsec_other.go:117` 的包文档逐字点名三条路（`proc.DefaultLayout` / `cmd/wisp's resolveDataDir` / `resolveSecretLayout`）。
⇒ **不 import buildinfo**，故腿型 **B（对 env 完全无关）**。它问的是 `TMPDIR`，不是 `WISP_ENV`。

**`internal/risk`（冻结件，只列不改）** — `blacklist.go:200-201` `envOr("APPDATA", …)` / `envOr("LOCALAPPDATA", …)`、`:348`；`pathresolver.go:181`；`syncdirs.go:91-92`。
⇒ 读的是 `APPDATA`/`LOCALAPPDATA`/`USERPROFILE`，**与 `WISP_ENV` 无关**，不进本票名单（列出来是为了让名单**可核**，不是为了让它变长）。

**其余生产 env 读（都不早退、也都不看 `WISP_ENV`）**：`internal/models/manifest.go:325`（`WISP_MODELS_MANIFEST`）、`internal/tools/fs_staging.go:101-107`（`user.Current()` → `USERNAME`/`USER`/`LNAME`/`LOGNAME`）、`internal/secret/store.go:96`（配置里 `env:` 引用的那枚**变量名**）。

### 1.3 测试码（会去问操作系统／或自称会问的用例，按包）

| 包 | 文件 | 它问谁 | 腿型／是否本票靶 |
|---|---|---|---|
| `cmd/wisp` | `dataroot_128_test.go`（4 枚 top-level） | `userConfigDir` seam ＋ `buildinfo.EnvString()` | **A，唯一的靶**；`:281` 那枚已由 `c2fa2e9` 钉，`:151`/`:192` 用字面量 env ⇒ 免疫；`:209` 纯字符串 |
| `cmd/wisp` | `dataroot_128_windows_test.go`（1 top-level × 4 子项） | 真子进程 ＋ 真的摘掉 `APPDATA`（`appDataFreeEnv128` 在 `:91`，摘除循环 `:95-101`，"没得摘就 `t.Fatal`" 的自证在 `:103-105`），`:113 return append(kept, "WISP_ENV=dev")` | **B**（含 `resident` 腿走 `proc.DefaultLayout`）；**自己钉了 dev** ⇒ 免疫。它就是 (b) 形状**改前就存在**的先例 |
| `cmd/wisp` | `secret_dataroot_119b_test.go`（3） | `resolveSecretLayout(EnvDev)` ＋ 真 `os.UserConfigDir()`（`:164`、`:185`） | **B**＋env 字面量 ⇒ 免疫。**⚠ 但它是 `//go:build !windows` ⇒ CI 上零分母，见 §3.4** |
| `cmd/wisp` | `leg_sink_nail_131_windows_test.go:409`、`:526`；`early_log_nail_130_windows_test.go:123`；`resident_sink_nail_127_windows_test.go:199`；`secret_argv_windows_test.go:236/337/371` | 自己 `Setenv`/`cmd.Env` 里带 `WISP_ENV=test` **且**带 `WISP_TEST_DATA_DIR=<声明根>` | **不是问操作系统**——它是"自己声明落点"。免疫，且**是 `AC#2(b)` 说的那枚"自己给自己设前提"的形状**（§4.3 展开） |
| `cmd/wisp` | `tempdir_resolved_124_test.go:35-38 sealableTempDir124` | `t.TempDir()` → `proc.SealableRoot` | 与 env 无关（票 124 的 harness 腿） |
| `cmd/wisp` | `secret_test.go:200/259/965/1011/1186` | 断言输出里有 `WISP_ENV=dev` / `=prod` | **免疫**：印的是 `secretCmd.env` 字段（`secret.go:374/423/459` 全用 `c.env`），不是 `EnvString()`。⚠ 但 `doctor.go:43` 印的**才是** `EnvString()` ⇒ (a) 若有人照抄这批改doctor的断言，会**当场变红**，见 §4.2 |
| `cmd/wisp` | `leg_sink_gate_131_test.go:560 loadMainPackage131`、`leg_dispatch_gate_133_test.go` | AST 静态门禁 | 与 env 无关（`resolveDataDirConsumers128` 靠它算名册，`dataroot_128_test.go:352-379`） |
| `internal/proc` | `envfork_test.go:75 TestLayoutForTestEnv`、`:315`、`:323`；`envfork_mutex_windows_test.go:24`；`boot_windows_test.go:33` | `t.TempDir()` 注入 ＋ 字面量 env；`DefaultLayout` 两枚（`:324`、`:327`）会**真的** `os.UserConfigDir()` | **B／免疫**。`envfork_test.go:323-330` 是同形邻居，见 §2.3 |
| `internal/buildinfo` | `env_test.go:22 TestResolveEnv`（`:23`/`:28`/`:33` 三发 `t.Setenv`） | 环境态本身 | **免疫**（自己钉三形）。⚠ 它是全仓**唯一正面断言 `WISP_ENV` 优先级**的用例，且**不在受影响的 job 里**（§3.3） |
| `internal/winsec` | `dataroot_symlink_119_other_test.go:147-150`、`:178-187`、`:349`、`:375`（`//go:build !windows`）；`seam_probe_root_125_other_test.go`（`!windows`）；`tempdir_resolved_124_windows_test.go` / `_other_test.go` | `TMPDIR`、`HOME`、`XDG_CONFIG_HOME`、`proc.Summarize(EnvProd)` | **B／对 `WISP_ENV` 无关**。跑在 `test-core`(ubuntu) 与 winsec 门禁 |
| `internal/risk` | `pathresolver_anchor_spelling_windows_test.go:40-41`、`pathresolver_junction_windows_test.go:27-28/242-243`、`pathshape_portable_test.go:124-139`、`syncdirs_ancestor_actable_leg_116_test.go` | `APPDATA`/`LOCALAPPDATA`/`USERPROFILE` | 与 env 无关；**冻结件，一字不提改法** |
| `internal/tools` | `fs_test.go`、`ticket90_test.go`、`bridge_junction_windows_test.go` | `HOME`/`TMPDIR`/junction | 与 env 无关 |
| `internal/config` `internal/memory` `internal/llm` `internal/perm` `internal/agent*` | `tempdir_resolved_124_test.go`（7 枚） | `t.TempDir()` 声明根 | 与 env 无关 |

### 1.4 「除 `cmd/wisp` 外没有生产码读这枚变量」的三条凭据　〔独立复现〕

```
git grep -ln 'buildinfo.EnvString\|buildinfo.ResolveEnv' 744d79e -- '*.go' | grep -v 'cmd/wisp'
  → No matches found
git grep -n 'WISP_ENV' 744d79e -- 'internal/**' ':!*_test.go'
  → 只有 internal/buildinfo/{env.go,buildinfo.go}（定义方）与 internal/proc/doc.go:13（注释）
git grep -n 'buildinfo' 744d79e -- 'internal/winsec/*.go' | grep -v _test
  → （空）
```

---

## §2　Q2：逐枚标「哪条用例在 `WISP_ENV=test` 下走早退分支」＋它两形都绿的可能性　〔独立复现（静态）／真跑未做〕

**表头说明**：`早退?` = 这条用例的判据是否经过 `doctor.go:256` 那枚早退；
`两形都绿?` = 未设 env 与 `WISP_ENV=test` 两形**都绿、且绿的原因不是它声称的那条闸门**——即"换分支不换色"，本票最危险的一族。

### 2.1 走早退分支的（全部在 `cmd/wisp/dataroot_128_test.go`，一枚 job：`test-windows` step `:422`）

| # | 用例 file:line | 断言是什么 | `test` 形下它实际在比什么 | 早退? | 两形都绿? |
|---|---|---|---|---|---|
| 1 | `:281 TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` × 5 子项（`runTextTask`/`cmdModels`/`cmdProviders`/`cmdDoctor`/`resolveSecretLayout`） | rc≠0 ＋ `out` 含全部 6 枚 `rescueMarkers128`（`:91-98`）＋ `assertDirEmpty128(cwd)`（`:343`） | **`c2fa2e9` 之前**：4 枚子项里 markers 缺失 ⇒ 那四枚**在 CI 上红**、第 5 枚绿；`assertDirEmpty128` **一直绿**，因为落点被 `TestDataDir()` 搬到仓外绝对路径，**"当前目录一字节不许多"是被另一枚回落点满足的，不是被拒绝满足的**（票面 §"事实"末句，我在钉版码上复核成立） | **是** | **修前：4 枚不绿、第 5 枚两形都绿；修后：0 枚**（`:323 pinEnvThatAsksTheOS128` 已钉，且 `:263-264` 正面自证 pin 生效） |
| 2 | `:151 TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128` × 2 子项（`"dev"`、`"prod"`） | `err != nil` ＋ `errors.Is(err, errDataDirUnresolved)` ＋ markers ＋ `filepath.Join("<用户配置目录>","wisp"/"wisp-dev")` ＋ `planted128Cause` ＋ cwd 空 | env 是**字面量**入参（`:153`/`:158`）⇒ 早退分支**根本不参与** | **否** | **否**——它一直走 OS 读。⚠ 但它**不钉 ambient**，靠的是"env 走参数"这枚巧合形状：`resolveDataDir(env string)` 收 string 而非 `buildinfo.Env`，所以将来有人把它改成内部自己 `EnvString()`，这枚用例会**静默变装饰**。这条是 (b) 该防的**下一个**破口，不是今天已有的 |
| 3 | `:192 TestAC2TestDataDirBranchStillResolves128` | `err == nil` ＋ `dir == proc.TestDataDir()`（`:199`） | 它**断言的正是早退本身**（"keep the refusal from swallowing the harness branch"，`:188-191`） | **是（被断言的对象）** | **是（两形必都绿）**——但这是**设计意图**。⚠ 另有一枚独立 weakness：期望值由**被测函数自己**算出（两边同调 `proc.TestDataDir()`），所以 `TestDataDir` 自身漂移时这枚用例**永远看不见**。它只挡"早退被删掉"，不挡"早退退错地方"。**要不要收紧由实现程判，本程不建议改法**（静态不敢定） |
| 4 | `:209 TestAC2RefusalMarkersAreNotAShortenableList128` | marker 列表长度 ≥6 ＋ 每个严格前缀都过不了 `allMarkersPresent128` | 纯字符串仪器自检，无 OS 读 | 否 | 否（不相关） |

### 2.2 票面 `AC#2(b)` 说的"自己给自己设前提"那一族——现量名单（**这些不是靶，是必须先分清的**）

这七处**主动**设 `WISP_ENV=test`（有的还加 `WISP_TEST_DATA_DIR`），它们的断言**不依赖** ambient ⇒ (a)/(b) 两支都**不动它们**；
但如果把它们当"加固腿"去钉，就是**假绿**：

- `cmd/wisp/leg_sink_nail_131_windows_test.go:409-410`、`:526-527`（in-process `t.Setenv` 两枚）
- `cmd/wisp/early_log_nail_130_windows_test.go:123`（子进程 `cmd.Env`）
- `cmd/wisp/resident_sink_nail_127_windows_test.go:199`（子进程）
- `cmd/wisp/secret_argv_windows_test.go:236`、`:337`、`:371`（子进程）
- `cmd/wisp/dataroot_128_windows_test.go:113`（反向：钉 **dev**，并 `:107-112` 把 `WISP_TEST_DATA_DIR` 摘掉）
- `scripts/slo-check.ps1:111-112`（`$env:WISP_ENV='test'` ＋ `WISP_TEST_DATA_DIR`）

判据一句话：**"声明根 + 比声明"＝腿型 B／免疫；"注 OS 失败 + 比是否拒"＝腿型 A／是靶。**

### 2.3 同形邻居：两形都绿、但**与 env 无关**（列出来是怕实现程把它当战果）

| 用例 | file:line | 为什么两形都绿 | 是不是本票的账 |
|---|---|---|---|
| `internal/proc/envfork_test.go:323 TestDefaultLayoutUnknownEnvFails` | `:324 if _, err := DefaultLayout(Env("nope")); err == nil { t.Fatal }` | **一枚断言两个可能来源**：`envfork.go:236-239` OS 读失败，或 `:86-88` 枚举 default 分支。健康机器上永远是后者；OS 读一坏，这枚就**在没检查枚举的情况下绿** | **不是**（env 是字面量、`proc` 不读 ambient）。但它是"换分支不换色"**形状的现成第二枚**，`AC#3` 若要造"摘掉先栽 test 那行"的饵，这枚比 `:192` 更像 |
| `internal/proc/envfork_test.go:327 DefaultLayout(EnvTest) must not error` | 同上 | 逼着问一次 `os.UserConfigDir()`，答案被 `LayoutFor:77-85` **丢掉** | 不是（同上） |
| `ci.yml:250 Environment fork assertion (WISP_ENV=test data dir)` 那一步 | 跑的是 `TestLayoutForTestEnv`（`envfork_test.go:75`） | 那枚用例 `:79` 传 `buildinfo.EnvTest` **字面量**、`:78 t.Setenv(TestDataDirEnv, dir)` **自己声明根** ⇒ **它从不读 `WISP_ENV`** | **是本票的账**：见 §2.4 |

### 2.4 现量撞到的、票面没写的一件：**规范义务没人履行**　〔独立复现〕

`docs/specs/SPEC-11-build-deploy-containerization.md:182-183`（钉版原文，逐字）：

> - 环境分叉在 CI 的可见性：`test-windows`/`test-core` job 显式断言 `WISP_ENV=test` 生效
>   （数据目录在临时路径、互斥未注册）。

`docs/specs/SPEC-03-config-secrets-envs.md:101`：

> - 安装产物默认 `prod`；`go run`/本地调试默认 `dev`（由 `buildinfo` 里的构建期默认值 +
>   环境变量覆盖）；**CI 与测试进程显式 `WISP_ENV=test`**。

⇒ SPEC-11 §8 要的是"**断言这枚变量生效**"。今天**唯一**一枚名字写着这件事的步骤（`ci.yml:250`）**根本不咨询这枚变量**（§2.3），
全仓也没有第二枚用例从 ambient 角度断言它（`internal/buildinfo/env_test.go:22` 是自己 `t.Setenv` 的三形自证，算"优先级"不算"CI 生效"，
且它跑在 `test-core` 的 `:288` 那一 scope 里、同样不看 job env）。
⇒ **`WISP_ENV=test` 今天被读去、被印出来，但从没被断言过。** 这一条无论 `AC#2` 选 (a)(b)(c) 哪一支都要落地，`docs/specs/**` 是禁改件，所以修的地方只能是 `ci.yml`＋测试，不是 spec。

### 2.5 静态答不出、必须真跑的（**我不猜，逐条挂着**）

1. 每一枚"换色"用例在**修后的 `WISP_ENV=test` 本机形**下的四数与名册差集——票 128 `AC#4` 已给过修前那两形（`R1 202/202/0/0`、`R2 202/192/10/0`、`R3 202/202/0/0`、`R4 202/202/0/0`，见 §5），本程**没复跑**。
2. `ci.yml:250` 那一步在 job env 被摘掉后是否仍绿——静态推断是"仍绿且逐字同字节"（§1.4＋§2.3 的理由），**但这是可跑掉的断言，本程按规矩不下结论**。
3. §4.2 那三枚 build 门禁在 (a) 下的真形：会不会红（`wisp doctor` 的其余检查在 hosted runner 上是否本就全绿），本程不跑，只列出会**新增什么写面**。

---

## §3　Q3：CI 的分母　〔独立复现（读 `ci.yml` 与三枚脚本，钉版）〕

### 3.1 六枚 job 与它们的 test 命令／scope

| job | runner | job env | 跑 test 的步骤 → 实际命令 → 解析到的包 |
|---|---|---|---|
| `lint` | ubuntu | **无** | `tools/d22scan/runtests.sh -C tools/d22scan ./...`(`:81`)、`sh scripts/d22scan.sh`(`:109`)、gofumpt(`:112`)、`go vet` ×2(`:122/:125`)、staticcheck(`:179`)、mockllm vet(`:216`) |
| `test-core` | ubuntu | `:227` | `:264` `runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`；`:288` `bash scripts/portable-tests.sh --scope=core` |
| `test-windows` | windows-latest | `:337` | `:386` `scripts/winsec-tests.sh` → `portable-tests.sh ./internal/winsec/`；`:400` `build.ps1 -Env dev`（**含 `wisp.exe doctor`**）；`:422` `scripts/wisp-cli-tests.sh` → `portable-tests.sh --scope=cli` → `./cmd/wisp/`；`:458` `--scope=windows`；`:474` `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` |
| `slo-smoke` | windows-latest | `:482` | `:497` `build.ps1 -Env dev`（**含 doctor**）；`:500` `slo-check.ps1 -Subset smoke`。**零枚 `go test`** |
| `slo-full` | **`[self-hosted, wisp-slo]`** | `:540` | `:560` `build.ps1 -Env dev`（**含 doctor**）；`:563` `slo-check.ps1 -Subset full`。**零枚 `go test`** |
| `lint-frontend` | ubuntu | 无 | npm 五步，与本票无关 |

`portable-tests.sh` 的四份名册（钉版 `:124-166`：`core_pin` `:124-150`、`win_pin` `:151-160`、`cli_pin` `:161-163`、`winsec_pin` `:164-166`）与 scope 表（`:170-203`）：

- `core`（`:171-183`，scope 数组 `:172-181`）解析出 25 枚包（`core_pin` `:125-149`）：含 `internal/proc`、`internal/buildinfo`（`:130`）、`internal/secret`、`internal/config`、`internal/memory`、`internal/risk`、`internal/winsec`（`:149`）、`cmd/llmrecord`……**不含 `cmd/wisp`**。
- `windows`（`:184-193`，scope 数组 `:185-189`）8 枚：`internal/proc` `internal/secret` `internal/config` `internal/risk` `internal/ball` `internal/perm` `internal/plugin` `cmd/llmrecord`。**不含 `cmd/wisp`、不含 `internal/winsec`**（winsec 有自己那一步 `:386`）。
- `cli`（`:193`，scope 数组只有 `:195` 一行 `./cmd/wisp/`）；注释在 `:192-194` 明写 "Only `scripts/wisp-cli-tests.sh` may call this"；该脚本 `:60-68` 在 `GOOS!=windows` 时 **`exit 2`**（明写 "Not a skip and not a pass"）。

### 3.2 直接回答"有几枚落在受影响的 job 里"

**受影响的 job ＝ 只有 `test-windows` 一枚**（它是**唯一**一处"带 `WISP_ENV=test` 的 job" ∩ "跑会读这枚变量的码"）。

| 第 §2 问那批用例 | 枚数 | 落在带 env 的 job 里？ |
|---|---|---|
| `cmd/wisp/dataroot_128_test.go` 全部 4 枚 top-level ＋ `dataroot_128_test.go:281` 的 5 枚子项 ＋ `dataroot_128_windows_test.go` 的 4 枚子项 | 8 top-level＋9 子项 | **是**（`:422`，全在 `test-windows`）。其中**只有 `:281` 的 4 枚子项曾真中过这枚早退**，现由 `c2fa2e9` 钉住 |
| `cmd/wisp` 其余 13 枚文件的全部用例（`secret_test.go` 14、`run_test.go` 6、`run_mode101_test.go` 5、`logsink_test.go`＋`logsink_windows_test.go` 3+3、`providers_test.go` 4、`leg_sink_gate_131`＋`leg_sink_nail_131`＋`leg_dispatch_gate_133` 1+3+1、`resident_sink_nail_127` 3、`early_log_nail_130` 2、`secret_argv_windows` 4、`tempdir_resolved_124` 0；按 `git show … | grep -c '^func Test'` 现量） | **49** 枚 top-level | 在**带 env 的 job** 里跑，但**腿型 B／免疫**（§2.2） |
| `internal/proc`＋`internal/buildinfo` 的 §2.3 那批 | 5 | `test-core`(带 env，`:227`) 与 `test-windows`(带 env，`:458`) 都跑它们 ⇒ **在受影响的 job 里，但不受这枚变量影响**（`proc` 不读 ambient，§1.4） |
| `internal/winsec` 的 119/124/125 那批 `!windows` 用例 | — | 只在 `test-core`(ubuntu) 的 `:288` 与 winsec 门禁；**与 env 无关** |
| `ci.yml:250` 那一步 | 1 步 | 在带 env 的 job 里，**但从不咨询它**（§2.3/§2.4） |

**净数一句话**：票面问"还在哪些包里变成装饰"⇒ **静态答：0 个别的包，1 个别的包外的门禁面（3 枚 build 步骤，见 §4.2），外加 1 条从没被履行的规范义务（§2.4）**。

### 3.3 一枚"在别的 job 里跑"的例外要点名

`internal/buildinfo/env_test.go:22 TestResolveEnv` 是全仓**唯一正面覆盖 `WISP_ENV` 解析顺序**的用例，它只在 `test-core`(ubuntu) 的 `--scope=core` 里跑
（`core_pin` 有 `github.com/CarlosShao/wisp/internal/buildinfo`，`portable-tests.sh:130`）。
⇒ 若有人以为"`cmd/wisp` 那批钉上 dev 就代表 env 语义在 CI 上有覆盖"，那半句话**只有 `test-core` 支撑，`test-windows` 一枚都不支撑**。

### 3.4 顺带量到、但**不该记在本票账上**的一件（文件级零分母）　〔独立复现〕

`cmd/wisp/secret_dataroot_119b_test.go:1` 是 `//go:build !windows`（3 枚 top-level，含票 119 `R-119-1` 的正腿），
而 `cmd/wisp` 的测试二进制**只**在 `test-windows` 通过 `--scope=cli` 跑（§3.1，`wisp-cli-tests.sh:60-68` 在非 windows 直接 `exit 2`）。
⇒ **那三枚用例在 CI 上从不编译、从不运行**。`portable-tests.sh` 的 GUARD A（`:300-312`，仪器命令在 `:300-301`）
**自己写明它是包粒度**的（`:305` 原文：`compile NO test file at all for GOOS=$goos`，靠 `go list -f '{{if or .TestGoFiles .XTestGoFiles}}'`），
一个包里只要还剩一枚 windows 形测试文件，`//go:build !windows` 的那一整枚文件就可以零分母而不响。
⇒ 这是票 111／票 98 那本"**有没有腿**"账的**新粒度**（文件级），不是票 140 的"走错分支"账。**建议编排者另立账行**，别并进本票的 AC。

---

## §4　Q4：最小改动面　〔独立复现（静态）／结论受 §2.5 三条真跑限制〕

### 4.1 (b) 支：每条用例自己钉环境（推广 `pinEnvThatAsksTheOS128` 的形状）

**现量到的、今天真正要钉但还没钉的枚数：Go 测试面 0 枚。** 依据：

- 唯一走早退的用例族 `dataroot_128_test.go:281` **已经钉了**（`c2fa2e9`，helper 在 `:257`，调用在 `:323`）；
- 同文件另两枚（`:151`、`:192`）env 走**字面量参数**，钉了反而是错的（`:192` 断言的正是早退）；
- 除 §2.1／§2.3 点名的三枚（`dataroot_128_test.go`、`dataroot_128_windows_test.go`、`secret_dataroot_119b_test.go`）外，
  `cmd/wisp` 还剩 **13 枚**测试文件（全目录 16 枚，`git ls-tree cmd/wisp/ | grep -c '_test\.go$'` 现量），共 **49 枚** top-level 用例（§3.2 同一条命令现量）——
  它们不是注入 `dataDir`／`env` 字段，就是自己 `Setenv`，**已经处于 (b) 要求的形状**。

⇒ **所以 (b) 的真实改动面不在 Go 里，在三枚 shell 门禁上**（§4.2）；再加上两条**防回归**的形状，我按"要动几枚文件、几处"给全：

| 动作 | 文件数 | 处数 | 说明 |
|---|---|---|---|
| 现有靶的钉 | **0** | **0** | `c2fa2e9` 已交（`dataroot_128_test.go:257-273`＋`:323`） |
| 让"新消费者也必须钉"变成响亮失败（`AC#2(c)` 与 (b) 可合成一支） | **1** | **1** | `dataroot_128_test.go:352-379 resolveDataDirConsumers128` 已经用 AST 名册逼平表；缺的是"名册里每一行都得先过 `pinEnvThatAsksTheOS128`"这一条判据——**加在那枚函数旁边最省** |
| 补上 §2.4 的规范义务（ambient 生效的正面断言） | **1–2** | **1 枚新用例（＋ `ci.yml` 一步的名字／引用）** | 必须是"注 `test` 与不注两形都跑"的自包含形状，**不能**再往 `TestLayoutForTestEnv` 上挂（它按设计不看 ambient） |
| (a) 的替身：给 build.ps1 那一步显式设 env（若选 (a) 才需要） | **1–3** | **1–3** | 见 §4.2 |

⇒ **(b) 单独开出去几乎是一枚空枪**。这句话对"这张票的实现程该怎么排"是决定性的，但我给它的证据强度只有静态，
所以 §4.4 保留了"要真跑才知道"的口子。

### 4.2 (a) 支：取消四枚 job 级 env ⇒ **哪些步骤会拿不到它现在依赖的 `test` 语义**（点名）

**先纠正简报里的一个指向**：简报举例说「`:250` 那一步"Environment fork assertion (WISP_ENV=test data dir)"显然就是靠它」——
**这一步不靠它**（凭据：`internal/proc/envfork_test.go:75` 那枚用例的 `:79` 传字面量 `EnvTest`、`:78` 自己 `t.Setenv(WISP_TEST_DATA_DIR)`；§1.4 的三条 grep 证明 `internal/proc` 不读 ambient）。
**真正靠它的是另外三枚，而且它们都不是 `go test`**：

| 会变的步骤 | file:line | 现在靠 env 做什么 | 摘掉 env 之后 |
|---|---|---|---|
| "cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)" | `ci.yml:395`/`:400` | `scripts/build.ps1:161-163`：`& build\wisp.exe doctor`；`if ($LASTEXITCODE -ne 0) { Fail }`。`doctor` → `doctor.go:104 EnvString()` → `resolveDataDir` | 未设 ⇒ `buildinfo.go:20`/`build.ps1:107` 的 `DefaultEnv=dev`（三处都是 `-Env dev`）⇒ `resolveDataDir("dev")` 走 `doctor.go:259` 去问 `%APPDATA%` ⇒ `probeWritable`（`doctor.go:300-309`）**`MkdirAll` 并写 `doctor-write-probe.tmp` 到 `%APPDATA%\wisp-dev`**；今天写的是 `%TEMP%\wisp-test-<pid>\…`。**不是"变红"，是"新增一棵仓外真实配置树的写面"** |
| "Build wisp.exe"（slo-smoke） | `ci.yml:496`/`:497` | 同上 | 同上（hosted runner，每台新机器） |
| "Build wisp.exe (deps cached on the runner)"（slo-full） | `ci.yml:559`/`:560` | 同上 | ⚠ **最重的一枚**：`slo-full` 的 `runs-on: [self-hosted, wisp-slo]`（`ci.yml:538`），`scripts/slo-check.ps1:115` 的注释自陈那台**跨 run 复用同一个工作目录**（原文：`The runner reuses E:\work\base\actions-runner\_work\wisp\wisp across runs`）。⇒ 摘 env 会把 `wisp doctor` 的写面从 `%TEMP%` 挪到**那台机器真实的 `%APPDATA%\wisp-dev`**，且**跨 run 残留** |

**不受 (a) 影响的（别虚报代价）**：
- `ci.yml:250`（§2.3 已证）、`:288`、`:458`、`:422`、`:386`、`:474`——`go test` 面的码**没有一处读 ambient `WISP_ENV` 之外**（读的那批是 `cmd/wisp`，其用例要么已钉要么注入）。
- `ci.yml:499`/`:562` 两枚 SLO 门禁：`slo-check.ps1:111-112` **自己设** `$env:WISP_ENV='test'` ＋ `WISP_TEST_DATA_DIR`；`cmd/wisp/slo_windows.go:238-240` 再兜一层。⇒ job env 对这两枚是**幂等冗余**。
- `cmd/wisp/*_windows_test.go` 里那六处子进程腿：`cmd.Env = append(os.Environ(), "WISP_ENV=test", …)` 自带（§2.2）。
- `internal/secret`／`internal/winsec`／`internal/risk` 全部：`grep` 证零命中。

**所以"谁读去这枚 env"（简报要求点名）的完整答案**：
生产读点 **3 处**——`internal/buildinfo/env.go:36`、`internal/buildinfo/buildinfo.go:38`、`cmd/wisp/slo_windows.go:238`；
决定落点的**只有一条链**——`buildinfo.EnvString()/ResolveEnv()` → `cmd/wisp/doctor.go:247 resolveDataDir` → `:256-258`（早退）／`:259`（问 OS）；
以及**一枚非测试消费者**——`build.ps1:162` 的 `wisp.exe doctor`。**说 (a) 零风险是错的**，风险不是红，是 `:400/:497/:560` 三枚步骤会开始往真实 `%APPDATA%` 里建目录、写探针文件，
其中 `:560` 落在 self-hosted 那台上。**能不能只在那三步显式设 `WISP_ENV: test`、把四枚 job 级摘掉？静态判断：可以且只需 3 处**（`ci.yml:400`、`:497`、`:560` 各自加 `env:`），
但这仍是共享件改动 ⇒ 票面 Rules 明写「`AC#2` 若选 (a)，**改前先登记交回编排者**，不许自行改 workflow」。**本程只量面，不提议**。

### 4.3 (c) 支：维持现状＋做成响亮失败——静态能落到哪几处

现成的、**今天就差一条判据**的两个位置：
1. `ci.yml:250` 那枚步骤名对不上它跑的东西（§2.4）——要么换用例，要么改名字，**两者都比现在响**。
2. `dataroot_128_test.go:263-264` 那枚"pin 没生效就 `t.Fatalf`"的形状可以推广成"名册里每一腿进表前都得自证问过了 OS"（§4.1 第二行）。

### 4.4 我给编排者的三档判断（**不给顺手的结论**）

- **静态高置信**：受影响的 job 只有 `test-windows`；受影响的包只有 `cmd/wisp`；(a) 的真实代价是 3 枚 build 步骤的落点搬移，不是 `:250`。
- **静态可推、仍须真跑背书**：(b) 在 Go 测试面上是**零枚新改动**（因为 `c2fa2e9` 恰好已经把唯一的靶钉住了）。**这一条我最不确定**，因为"换分支不换色"那一族的定义就是"看起来一直绿"——只有两形差集能证伪。
- **本程完全答不了**：`AC#3` 的两枚变异（摘 `:261` 的 pin / 摘 `:260` 的"先栽"）在 runner 上响不响；票 123 那 4 枚审批超时用例是否被任何候选改法波及；`:400/:497/:560` 在 hosted 与 self-hosted 上 `wisp doctor` 的其余检查本来绿不绿。

---

## §5　引用别程读数时的锚点核对　〔日志＋归档抽验〕

| 我引的读数 | 出自 | 它是在**哪一版树**上量的 | 我核到了什么 |
|---|---|---|---|
| 四数 `R1 202/202/0/0`、`R2 202/192/10/0`、`R3 202/202/0/0`、`R4 202/202/0/0`、panic 全 0、名册 101 枚逐名一致、MUT-4 摘 pin 本机就红 | `docs/evidence/s1/128-ac4-gates-and-ci-divergence.md`（commit `ef26704`），其 §0 收尾 HEAD | 该文 §0 自称的 HEAD；**我没有独立复算那 202 枚，本程一个测试都没跑** | 我只复核了它引用的**码面**：`c2fa2e9 test(wisp/128 AC#4): 进程内那五条腿钉住 WISP_ENV —— 先栽 runner 的 test 再 pin 回 dev` 真在 `git log`（`git log --oneline` 现量），且 `dataroot_128_test.go:260-261` 两行逐字与该 commit message 相符（钉版 `744d79e`） |
| run `35967768017` ＝ CI 上四红一绿那一发 | 同上＋`ci.yml:337` | 该 run 的 head sha **不在我这份证据里**，我没有把它绑到某个 sha ⇒ **引用它当判据时必须连 sha 一起写**，否则下次颜色不变就是归因腐坏 | 未核（要 `gh`，本程按只读盘点未取） |
| 「四枚 job 级 env」＋「`:131` 那枚绿腿」 | `docs/reports/pending-and-issues.md` **A175**（`:5489-5499`） | 台账自述 09-24 17:5x，锚点即 `744d79e` | 四枚 env 行号**我复核一致**；`:131` 那一处**我复核为 `:139`**（钉版），与台账自记的"差 8 行"一致 ⇒ 台账那条更正**成立** |
| `SPEC-11 §8` / `SPEC-03 §5.1` 两句 | `docs/specs/SPEC-11…:182-183`、`SPEC-03…:101` | 钉版 `744d79e`（`docs/specs/**` 是禁改件，钉版即权威） | 逐字引于 §2.4，未改写 |

**〔仅自述不背书〕**：票 128 `AC#4` 那句「修之前"当前目录一字节都不许多"是被另一枚回落点满足的」——**它的实质我在码上复核成立**
（`resolveDataDir` 早退到 `TestDataDir()` 仓外绝对路径 ⇒ `assertDirEmpty128(cwd)` 必然空），
但它附带的"CI 那一发真的这样绿过"我只在日志里看到，**本程未复算**。

---

## §6　临时件清单（只建不删）

本程**没有在仓内或 `D:\tmp` 下创建任何临时文件／目录**。
唯一一个由工具替我落盘的旁产物（在仓库之外，未删）：

- `C:\Users\swq\.qoder-cn\tmp\D--work-workspace-projects-plans-Wisp\tool-outputs\session-a91d07b8-f536-47da-a2b4-a297d49f5f4f\b5js0lz8f.output`
  —— 一条 `git grep` 输出超限被 harness 转存的文件，**内容是我的第一次全量 grep，无凭据值**。

## §7　commit 账（本程自己的 git 写操作）

见本文件末尾追加的那一节（§8）。**取法**：`git log --oneline -- docs/evidence/s1/140-static-inventory-r1.md`。

## §8　本文件的 commit 账　〔独立复现〕

（本节在写完全文后追加，逐枚 `git show --name-only` 现量。）
