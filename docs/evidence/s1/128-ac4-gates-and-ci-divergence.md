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
