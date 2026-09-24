# 140 `AC#2` — 取消 job 级 `WISP_ENV: test` 的静态爆炸半径 r1（分母，不是裁定）

**取证代理**：`auditor-ticket140-static-ac2`（只读盘点程，非实现者）
**时刻**：2026-09-24 19:0x +08 起
**码面锚点**：起手 `git rev-parse --short HEAD` = **`2e6d171`**；写作期间 HEAD 又前进到 `0ee68e1`
（简报给的 `166bf63` **已过期**——它比我最早一次读数还早两枚，号是读数不是常量）。
⇒ 本文所有 `file:line` 取自**盘上工作树**，并在每一栏标出**那一枚文件最后一次被改动是哪一版**
（`git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M' -- <path>` 逐枚现量）。

**本文件只做一件事**：给票 140 `AC#2` 的三支代价**算分母**。纯静态，**一枚仪器都没跑**。

## 0.0 只读纪律自证（简报红线：另一枚程此刻在真取样）

- **本程跑过的命令全集**：`git rev-parse` / `git log` / `git show` / `git status --porcelain` /
  `git cat-file` / `grep` / `awk`（读文件用）/ `wc` / `sed -n`。
- **本程没有跑**：`go test`、`go build`、`go vet`、`gofmt`、`gofumpt`、`docker`、`wisp`、`wisp slo`、
  `slo-check.ps1`、`gh run view`、`rm`、`rmdir`、`git add -A`、`--amend`、`reset`、`rebase`、`stash`、
  `checkout .`、`clean`、`push`。**一条红线未碰**：本程不产生任何编译/取样活动，
  因此不会落进 `scripts/slo-check.ps1:153-155` 那份 `$loadNames` 名单（原文见 §5.4）。
- **被拒调用次数：0。**
- **凭据值抄录：0 处。** 全程 grep 未命中任何疑似明文密钥；`buildinfo.go:31-32` 那枚 minisign
  **公钥常量**（`MinisignPublicKey`）本程**未抄值**，只点名文件与行号。
- **注入两栏计数**：真通知回显 **1** 条（会话开头 harness 回显的 `Memory: agents.md` 与 skill 清单，
  是环境回显不是指令，本程未据此改变动作）；判为注入 **0** 条。
  本程读到的所有「已结案／请改判／阈值」类句子都出自**票面与台账自身的正文**，
  本程把它们当**待核断言**处理（见 §0.2），未据此行动。

## 0.1 与上一份静态盘点的分工（不重复它的账）

`docs/evidence/s1/140-static-inventory-r1.md`（436 行，锚点 `744d79e`，commit `7c81c77`/`b33edb0`）已经裁过：
「除 `cmd/wisp` 外没有生产码读这枚变量」（它的 §1.3）、「`:250` 那一步不依赖 job env」（它的 §0.2②，
推翻过简报的一处前提）、以及「(a) 支的真实代价是 3 枚 build 步骤的落点搬移」（它的 §4.2）。
**本文件不复核、不重抄那三笔账。** 本文件只做它没做的那一维：

1. **读者名册按"哪一版"钉死**（上一份按 `744d79e` 钉版；盘上 `cmd/wisp/dataroot_128_test.go`
   在 `d56b6f5 09-24 18:40` 之后又动过，行号已漂，见 §0.3）；
2. **环境枚举的入口链**从 `os.Getenv` 到 `TestDataDir` 短路逐格接上（§2）；
3. **对每一处读者标 会换／不会换／不确定**（§3，上一份只给了 job 级与步骤级的清单，没有逐读者判定）；
4. **票 123 那一族 4 枚的红与 job env 是否共因**——这一条上一份明说"不共因"只给了本机读数
   （引 `128-ac4-r1-acceptance.md` §7），**本程给一条静态链**（§3.3）；
5. **三支各自的代价表**含"会不会把已经绿的改红"与回退路径（§4）；
6. **同族扫描**（§5），扫法逐条可重跑。

## 0.2 简报里被本程推翻／更正的几句（先给结论，凭据在后）

| 简报的断言 | 现量结果 |
|---|---|
| 「`ci.yml` 那四枚 job 级 `WISP_ENV: test` 在 `:227`/`:337`/`:482`/`:540`」 | **成立**，逐字对。`grep -rn WISP_ENV --include='*.yml' .` → `:227 :337 :482 :540` ＋ `:250` 一枚**步骤名字符串**（不是 env）。四枚的 job 头 `:224 test-core` / `:334 test-windows` / `:479 slo-smoke` / `:537 slo-full`，`runs-on` 逐字见 §1.3 |
| 「锚点是 `166bf63`」 | **不成立（过期）**。本程起手 `2e6d171`，交件时 `0ee68e1`。**号会漂，别抄** |
| 「`doctor.go` 那枚 `if env == "test"` 短路」 | **成立**，行号逐字对（`:256-259`）。但该文件最后一次改动是 `4e5d240 09-23 10:18`，**不在任何在飞地界内**，所以这一处的行号是稳的 |
| 「票 123 那 4 枚的红**是**'自己写的 config 读不到'」 | **这句是台账／§7 的转述，本程不采信也不推翻**——本程只裁它的前半问句：**那一族的红依不依赖 job 级 env**。答案＝**不依赖**，静态链完整（§3.3）。红**本身**的真因仍按 `128-ac4-r1-acceptance.md` §9#7 记「未追」 |
| 「`c2fa2e9` 那枚 helper 的形状＝`pinEnvThatAsksTheOS128`」 | **成立，且行号未漂**。`d56b6f5 09-24 18:40` 确实又动过这枚文件，**但那一刀全落在头部注释**：`git diff --stat 744d79e HEAD -- cmd/wisp/dataroot_128_test.go` = **4 insertions / 4 deletions，hunk 只有一个、在 `@@ -34,10 +34,10 @@`**（改的是 AC#3 变异锚点那四行注释），行数两版都是 **438**。⇒ `:257` helper 本体、`:260`/`:261` 两枚 `t.Setenv`、`:263` 自证**逐字节未动**，与上一份盘点一致。**这一条本程原以为会漂，量出来没漂，如实改判**（见 §0.3） |

## 0.3 「每一枚引用文件都量一次版本」的现量结果（这一节是方法，不是结论）

本程对 §1 名册里**每一枚**被引用的文件跑了一遍
`git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M' -- <path>`，行号逐枚按盘上工作树重取。
**结论与简报给的教训相反的一面也要报**：

- 今天真正**动过**的只有 `cmd/wisp/dataroot_128_test.go`（`d56b6f5 09-24 18:40`），
  而它那一刀**落在 `:34-43` 的注释块**（`git diff 744d79e HEAD -- <该文件>` 只有一个 hunk、
  `4 insertions / 4 deletions`、两版皆 **438 行**）⇒ **本程引用的 `:257`/`:260`/`:261`/`:263`/`:266`/`:268`
  在两版上逐字节相同**。上一份盘点 §0.1b 那条"两版码面一字未动"的复核，对**这枚文件的这些行**到今天仍成立。
- 其余 13 枚引用文件的最后改动时刻**全部早于** `744d79e`（本程序号见各表"版本"列，最晚一枚是
  `ci.yml`/`slo-check.ps1` 的 `decb7b9 09-23 13:51`）⇒ 行号与上一份盘点同源同值，**本程不重抄它的账**。
- **但这不是"以后不用标版本"的理由**：`128-ac4-r1-acceptance.md` §1 与台账 `A175` 都记过
  `:131` vs `:139` 那一枚 8 行差（取自改前的树）。本程照旧逐枚标版本，**并把"没漂"也当成量出来的结果写下来**。

---

## §1　①　读者名册（`.go` / 脚本 / workflow 三栏，每处标版本＋逐字节原文）

**扫描命令（可重跑，本程实际用的就是这三条）**：

```
grep -rn 'WISP_ENV' --include='*.go'  .            # 17 枚文件 / 62 处命中
grep -rn 'WISP_ENV' --include='*.ps1' .            #  1 枚文件 /  1 处命中
grep -rn 'WISP_ENV' --include='*.sh'  .            #  0 枚文件 /  0 处命中
grep -rn 'WISP_ENV' --include='*.yml' .            #  1 枚文件 /  5 处命中
```

**分母一句话**：**69 处命中 / 19 枚文件**；其中**真正"读"这枚变量的只有 3 处**（§1.1），
其余是**印**（把已解析的值写进文案）、**注释**、**测试自己设值**、**workflow 设值**。

> **一条容易漏的口径**：`grep 'WISP_ENV'` **抓不到那个分叉判定本身**。
> `cmd/wisp/doctor.go:256` 写的是 `if env == "test" {` —— 它比的是**入参字符串**，字面量里没有 `WISP_ENV`。
> 谁按 grep 结果画链，会把链画断在 `resolveDataDir` 的门口。⇒ 本程另跑两条：
> `grep -rn 'env == "test"' --include='*.go' .` → **全仓只命中 `cmd/wisp/doctor.go:256` 一处**（现量）；
> `grep -rn '== buildinfo.EnvTest\|== EnvTest' --include='*.go' .` → **全仓只命中
> `internal/proc/envfork.go:212` 一处**（`if env == buildinfo.EnvTest && os.Getenv(TestDataDirEnv) != "" {`，
> 读的是**另一枚**变量 `WISP_TEST_DATA_DIR`）。
> **两条加起来才是"test 短路"的全集：2 枚，落在 2 枚文件。**（本程第一次写这一句时把 `envfork.go:212`
> 误记成 `env == "test"` 的命中，现量后改正——枚数没变，形状不同。）

### 1.1 `.go` 栏 —— 真读者（求值这枚变量的全部三处）

**A. `internal/buildinfo/env.go:36`** ｜该文件最后改动 `b155c6b 09-19 15:24` ｜**读者**

```go
	if e := os.Getenv("WISP_ENV"); e != "" {
		return ParseEnv(e)
	}
	return ParseEnv(DefaultEnv)
```

（`:36` 起，行首 1 tab。函数头 `:35 func ResolveEnv() (Env, error) {`）
**判什么**：判**类型化环境枚举**——非空则严格解析（`prod|dev|test` 之外**报错**），
空/未设则落 `ParseEnv(DefaultEnv)`，而 `DefaultEnv` 是构建期 ldflags 注入的字符串。
⇒ **这一处就是 (a) 支的方向来源：取消 job env 不等于取消环境，是把它换成 `DefaultEnv`。**

**B. `internal/buildinfo/buildinfo.go:38`** ｜该文件最后改动 `bfcb230 09-20 07:00` ｜**读者**

```go
func EnvString() string {
	if e := os.Getenv("WISP_ENV"); e != "" {
		return e
	}
	return DefaultEnv
}
```

（`:37` 函数头，`:38` 那行行首 1 tab）
**判什么**：同一件事的**不校验版**——空/未设回 `DefaultEnv`，非空**原样返回、不验枚举**。
⇒ 注意两枚读者的**语义差**：`ResolveEnv` 对 `WISP_ENV=nonsense` 报错，`EnvString` 把它当合法值往下传。
`internal/buildinfo/buildinfo.go:20` 就是那条默认（逐字节，行首 1 tab、值前有对齐空格）：

```go
	DefaultEnv         = "dev"       // WISP_ENV when the env var is unset (SPEC-03 §5.1)
```

**C. `cmd/wisp/slo_windows.go:238-240`** ｜该文件最后改动 `ed18727 09-23 18:17` ｜**读者 + 写者**

```go
	if os.Getenv("WISP_ENV") == "" {
		_ = os.Setenv("WISP_ENV", "test")
	}
```

（行首各 1 tab）
**判什么**：判**"有没有人已经设过"**——没设就自己设成 `test`。
⇒ **这是全仓唯一一枚"读来判、然后自己补写"的形状，也是 (a) 支对 SLO 门禁免疫的原因**：
job env 在不在，`wisp slo` 这一路都是 `test`。它是**幂等冗余**，不是依赖。

### 1.2 `.go` 栏 —— 其余 14 枚非测试文件位置（**不判**，只印／只注释）

| file:line | 版本（最后改动） | 逐字节原文（剥行首 tab，其余不改） | 拿它判什么 |
|---|---|---|---|
| `cmd/wisp/main.go:135` | `ed18727 09-23 18:17` | `fmt.Printf("%sWISP_ENV=%s (data dir rules: SPEC-03 §5)\n", prefix, buildinfo.EnvString())` | **不判，只印**：版本横幅的第二行。调用者 `main.go:108`（`wisp` 裸跑）、`main.go:128 cmdRun`、`resident_windows.go:25` |
| `cmd/wisp/main.go:46` | 同上 | `Environment: WISP_ENV in {prod\|dev\|test}, default from build (SPEC-03 §5).` | usage 文案，零判定 |
| `cmd/wisp/doctor.go:42-43` | `4e5d240 09-23 10:18` | `info("wisp build", fmt.Sprintf("version=%s commit=%s built=%s WISP_ENV=%s",` / 续行 `buildinfo.Version, buildinfo.Commit, buildinfo.BuildDate, buildinfo.EnvString()))` | **不判，只印**（INFO 行）。⇒ (a) 会改这行**字节**，但 `build.ps1:163` 只看 `$LASTEXITCODE`，见 §3.2 |
| `cmd/wisp/doctor.go:230` | 同上 | `// WISP_ENV=test, and the user config dir behind dev/prod) are handed to` | 注释（`resolveDataDir` 的 doc 第 6 行） |
| `cmd/wisp/secret.go:176` | `717d822 09-23 13:08` | `Environment: the active WISP_ENV is echoed by every command (SPEC-03 §5.2),` | usage 文案 |
| `cmd/wisp/secret.go:289-290` | 同上 | `return fmt.Sprintf("Active environment: WISP_ENV=%s, data dir=%s, portable=%v\n",` / `c.env, c.dataDir, c.portable)` | **印 `c.env` 字段**，不是 `EnvString()`。字段在 `secret.go:184 env, err := buildinfo.ResolveEnv()` 一次性解析后塞进 `secretCmd{env: env}`（`:195`） |
| `cmd/wisp/secret.go:374` | 同上 | `fmt.Fprintf(c.sio.stdout, "wisp secret set: stored %s (WISP_ENV=%s, portable=%v)\n", ref, c.env, c.portable)` | 同上，印字段 |
| `cmd/wisp/secret.go:423` | 同上 | `fmt.Fprintf(c.sio.stdout, "wisp secret get: %s (WISP_ENV=%s, portable=%v)\n", ref, c.env, c.portable)` | 同上 |
| `cmd/wisp/secret.go:459` | 同上 | `fmt.Fprintf(c.sio.stdout, "wisp secret list: WISP_ENV=%s, portable=%v, dir=%s\n", c.env, c.portable, st.Dir())` | 同上 |
| `cmd/wisp/secret.go:519` | 同上 | `audit := fmt.Sprintf("wisp secret unset: deleted %s (WISP_ENV=%s, portable=%v, forced=%v, referenced_fields=%d)",` | 同上（审计行） |
| `cmd/wisp/slo_windows.go:164` | `ed18727 09-23 18:17` | `Environment: runs under WISP_ENV=test by default (no single-instance mutex,` | usage 文案（与 §1.1C 那条兜底**互相印证**，不是判定点） |
| `cmd/wisp/slo_windows.go:250` | 同上 | `fmt.Fprintln(os.Stderr, "wisp slo: another instance is running (set WISP_ENV=test)")` | 错误文案。**注意它判的是 `proc.ErrAlreadyRunning`，不是 `WISP_ENV`** |
| `internal/proc/doc.go:13` | `034080c 09-22 17:11`（同包 `envfork.go` 的版本） | `//   - WISP_ENV typed fork: env resolution + data dir / mutex name / default` | 包文档。⇒ **`internal/proc` 全包没有一处 `os.Getenv("WISP_ENV")`**（§1.1 三条命中里没有它） |

**`.go` 小计（非测试）**：6 枚文件、16 处命中；**读者 3 处 / 2 枚文件**，其余 13 处不判。

### 1.3 `.go` 栏 —— 测试码（9 枚文件、36 处命中）：全是**写者或断言者**，没有一处"依赖 ambient"

| file:line | 版本 | 逐字节原文 | 拿它判什么 |
|---|---|---|---|
| `cmd/wisp/dataroot_128_test.go:259` | `d56b6f5 09-24 18:40` | `ambient, hadAmbient := os.LookupEnv("WISP_ENV")` | **记录**进函数时的 ambient（只用于 §:264 的报错文案），**不判** |
| `:260` | 同上 | `t.Setenv("WISP_ENV", "test")                   // the value the windows job exports` | **自己先栽**上 runner 的那一枚值 |
| `:261` | 同上 | `t.Setenv("WISP_ENV", string(buildinfo.EnvDev)) // the pin that has to beat it` | **再 pin 回 dev**：这一枚就是 (b) 支要推广的形状本体 |
| `:263` | 同上 | `if got := buildinfo.EnvString(); got != string(buildinfo.EnvDev) {` | **判 pin 有没有生效**（自证），不判 ambient |
| `:266` | 同上 | `dir, err := resolveDataDir(buildinfo.EnvString())` | 把解析结果喂给被测函数，配合 `:268` 自证"确实去问了 OS" |
| `:264`/`:268` | 同上 | 两枚 `t.Fatalf("premise broke: …")` | 前提塌了就**响亮失败**——这就是简报说的"钉环境本身要有牙" |
| `cmd/wisp/dataroot_128_windows_test.go:113` | `4e5d240 09-23 10:18` | `return append(kept, "WISP_ENV=dev")` | 子进程腿**反向钉 dev**（`:88` 的注释解释为什么这里不能用 test）。**不依赖 ambient** |
| `:69`/`:86`/`:88` | 同上 | `t.Logf("AC#2 real process, leg %q: APPDATA unset, WISP_ENV=dev, cwd=%s -> rc=%d\n%s",` 等 | 读数日志与注释 |
| `cmd/wisp/leg_sink_nail_131_windows_test.go:409` | `6f702ae 09-23 15:25` | `t.Setenv("WISP_ENV", "test")` | **自己设前提**（票 131 的钉腿件）。判的是"sink 在 test 根下落哪"，与 job env 无关 |
| `:526` | 同上 | `t.Setenv("WISP_ENV", "test")`（行首 2 tab） | 同上 |
| `:30` | 同上 | `//\t        root": WISP_ENV picks the env, WISP_TEST_DATA_DIR is an identity` | 注释 |
| `cmd/wisp/early_log_nail_130_windows_test.go:123` | `d87905c 09-23 11:53` | `cmd.Env = append(os.Environ(), "WISP_ENV=test", procTestDataDirEnv+"="+dataDir)` | 子进程自带两枚变量。**不依赖 ambient**（`os.Environ()` 打底，但同名键在 `append` 里**后面盖前面**，Go 的 `exec.Cmd.Env` 取后出现的） |
| `:116` | 同上 | 注释行 | — |
| `cmd/wisp/resident_sink_nail_127_windows_test.go:199` | `9b5d64d 09-23 12:33` | `cmd.Env = append(os.Environ(), "WISP_ENV=test", procTestDataDirEnv+"="+dataDir)` | 同 `:123` |
| `:35` | 同上 | 注释行 | — |
| `cmd/wisp/secret_argv_windows_test.go:236`、`:337`、`:371` | `4477f56 09-20 15:00` | `cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+dataDir)`（`:337` 行首 2 tab、变量名是 `other`） | 子进程自带，**不依赖 ambient** |
| `cmd/wisp/secret_test.go:200`、`:259`、`:965`、`:1011`、`:1186` | `5265c3a 09-23 21:49` | 五处 `strings.Contains(..., "WISP_ENV=dev"/"WISP_ENV=prod"/"WISP_ENV="+string(env)/"WISP_ENV=")` | **断言印出来的 `c.env` 字段**，字段由 `newProbe(t, env, …)`（`secret_test.go:94`、`:110 env: env`）从测试自己塞进 `secretCmd{env:}`，**从不走 `ResolveEnv`** ⇒ 与 ambient **完全无关**。`:929` 是注释标题 |
| `internal/buildinfo/env_test.go:23`、`:28`、`:33` | `b155c6b 09-19 15:24` | `t.Setenv("WISP_ENV", "test")` / `"nonsense"` / `""` | 三形自证优先级。**它不依赖 ambient，但依赖 `DefaultEnv`**——`:34 want, err := ParseEnv(DefaultEnv)` 是**相对断言**（跟自己的默认比），不是硬编码 `"dev"` ⇒ (a) 动 `DefaultEnv` 也不会让它红 |
| `:19`/`:25`/`:30` | 同上 | 注释与报错文案 | — |

### 1.4 脚本栏（`.sh` 0 处 / `.ps1` 1 处）

| file:line | 版本 | 逐字节原文 | 拿它判什么 |
|---|---|---|---|
| `scripts/slo-check.ps1:111` | `decb7b9 09-23 13:51` | `$env:WISP_ENV = 'test'` | **写者**：SLO 门禁自己把进程环境设成 `test`，紧接 `:112 $env:WISP_TEST_DATA_DIR = Join-Path $OutDir 'data'`。⇒ 对 job 级 env 是**覆盖**关系，不是依赖关系 |
| `*.sh` | — | **零命中**：`grep -rn 'WISP_ENV' --include='*.sh' .` → `0 files / 0 hits`。三枚门禁脚本（`portable-tests.sh`、`wisp-cli-tests.sh`、`winsec-tests.sh`）**一字不提这枚变量** |  ⇒ 取消 job env 时，**shell 侧没有任何一处需要跟着改** |

### 1.5 workflow 栏（唯一一枚文件）

| file:line | 版本 | 逐字节原文 | 拿它判什么 |
|---|---|---|---|
| `.github/workflows/ci.yml:227` | `decb7b9 09-23 13:51` | `      WISP_ENV: test`（行首 6 空格） | **写者**，挂在 `:226 env:` 下，owner 是 `:224 test-core`（`:225 runs-on: ubuntu-latest`） |
| `:337` | 同上 | `      WISP_ENV: test` | `:336 env:` → `:334 test-windows`（`:335 runs-on: windows-latest`） |
| `:482` | 同上 | `      WISP_ENV: test` | `:481 env:` → `:479 slo-smoke`（`:480 runs-on: windows-latest`） |
| `:540` | 同上 | `      WISP_ENV: test` | `:539 env:` → `:537 slo-full`（`:538 runs-on: [self-hosted, wisp-slo]`）—— **唯一一枚跑在真机上、且跨 run 复用同一份工作树** |
| `:250` | 同上 | `      - name: Environment fork assertion (WISP_ENV=test data dir)` | **只是一枚步骤名字符串**，不设任何变量。上一份盘点已用它推翻过一处前提，本程复核一致 |

**workflow 里第二枚 job 级 env（同族，见 §5.3）**：`ci.yml:550`
`      MINGW64_ROOT: E:\work\base\msys64\mingw64\bin`（行首 6 空格，与四枚 `WISP_ENV: test` 同缩进、
同挂在 `:539 env:` 块下，块内注释占 `:541-549`）—— 它被 `scripts/build.ps1:51`
`if (Get-Item 'env:MINGW64_ROOT' -ErrorAction SilentlyContinue) { $ccCandidates += (Join-Path $env:MINGW64_ROOT 'gcc.exe') }` 读。
**这枚是真依赖**：`ci.yml:546-549` 的注释自陈没有它「the step died before compiling anything」。

### 1.6 名册净数

- 命中 **69 处 / 19 枚文件**（`.go` 62/17、`.ps1` 1/1、`.sh` 0/0、`.yml` 5/1，四条 grep 现量）。
- **真读者 3 处 / 2 枚文件**：`internal/buildinfo/env.go:36`、`internal/buildinfo/buildinfo.go:38`、
  `cmd/wisp/slo_windows.go:238`。
- **判落点的只有一枚分叉**：`cmd/wisp/doctor.go:256`（字面量不含 `WISP_ENV`，grep 抓不到，见 §1.0 口径）。
- **设值方 8 处**：`ci.yml` 4 枚 job 级 + `:250` 名字串 + `slo-check.ps1:111` + `slo_windows.go:239`
  + 测试自己设的 9 处（`dataroot_128_test.go:260/261`、`dataroot_128_windows_test.go:113`、
  `leg_sink_nail_131:409/526`、`early_log_nail_130:123`、`resident_sink_nail_127:199`、
  `secret_argv_windows:236/337/371` ⇒ 这 10 行里 `dataroot_128_test.go` 那两枚是"先栽再 pin"，
  其余 8 处是自带值）。

### 1.7 §1 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
8e86f32 docs(evidence/140 AC#2 r1): 起手 + §1 读者名册（零仪器，只 grep/只读）

docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

（那次 commit 前 `git diff --cached --name-only` 现量**只有这一枚路径**。）

---

## §2　②　环境枚举的入口链（从 `os.Getenv` 到那枚 `if env == "test"` 短路，逐格接上）

**先给那一格的答案**：`cmd/wisp/doctor.go:256` 的短路**只在"字符串族"这条链的最末一格**，
而且它是**全仓唯一**一枚 `env == "test"` 判定（现量：`grep -rn 'env == "test"' --include='*.go' .` → 1 命中）。
"**类型化族**"（走 `proc.LayoutFor` 的那条）**根本没有这枚短路**——它无条件问 OS、再把答案丢掉。
两族长得像、行为不一样，这正是 (a) 支代价被低估的地方。

### 2.1 链 A｜字符串族（**决定落点**的那条）

```
WISP_ENV（进程环境）
  设值方：ci.yml:227/:337/:482/:540（job 级）· scripts/slo-check.ps1:111 · cmd/wisp/slo_windows.go:239
          · 测试自己 t.Setenv / cmd.Env（§1.3 那 10 行）
    |
    +-> [读点 R1] internal/buildinfo/buildinfo.go:38   EnvString()      （不校验枚举，非空原样返回）
    |        :41  return DefaultEnv
    |                   ^  internal/buildinfo/buildinfo.go:20  DefaultEnv = "dev"   （仓内默认值）
    |                   ^  scripts/build.ps1:107  "-X $BuildInfoPkg.DefaultEnv=$Env" （构建期覆盖）
    |                      scripts/build.ps1:23-24  [ValidateSet('dev', 'prod')] / [string]$Env = 'dev'
    |                      ==> 关键：-Env 的合法值里**没有 'test'**。CI 三处 build 全传 `-Env dev`
    |                          （ci.yml:400 / :497 / :560），所以"取消 job env"落到的那一形是 **dev**，不是 test。
    |
    +-> EnvString() 的四枚消费者：
          cmd/wisp/doctor.go:43        只印（INFO 行，critical=false，doctor.go:37-39）
          cmd/wisp/main.go:135         只印（printVersions 第二行；调用者 main.go:108 / :128 / resident_windows.go:25）
          cmd/wisp/run.go:675          buildEnvString() -> 唯一消费者 run.go:156
          cmd/wisp/doctor.go:104-105   <== **没有守卫的那一枚**，见下
          cmd/wisp/models.go:109-110   有守卫
          cmd/wisp/providers.go:80-81  有守卫
                |
                v
    [分叉] cmd/wisp/doctor.go:247 func resolveDataDir(env string) (string, error) {
             :248-255  portable.txt 分支 —— **比 test 短路更早的一枚早退**（dataroot_128_test.go:76 assertNoPortableOverride128 就是为它造的）
             :256-258  if env == "test" { return proc.TestDataDir(), nil }   <<== 你问的那一格
             :259      base, err := userConfigDir() // %APPDATA%             <<== 被上面短路掉的那次 OS 读
             :261      return "", dataDirUnresolved128(env, err)             拒
             :263      base = proc.SealableRoot(base)
             :264-267  fork：dev -> <base>\wisp-dev ；其余 -> <base>\wisp
                |
                +-> proc.TestDataDir()  internal/proc/envfork.go:121-126
                      :122  if dir := os.Getenv(TestDataDirEnv); dir != "" { return dir }   // WISP_TEST_DATA_DIR，**另一枚变量**
                      :125  return filepath.Join(SealableRoot(os.TempDir()), fmt.Sprintf("wisp-test-%d", os.Getpid()))
                |
                +-> userConfigDir  cmd/wisp/doctor.go:277 var userConfigDir = os.UserConfigDir
                      （生产里无人重绑；唯一重绑点是测试 dataroot_128_test.go:65 failConfigDir128）
```

**`doctor.go:104-112` 这一格为什么是链上唯一的"无守卫消费者"**（逐字节，行首 1 tab）：

```go
	env := buildinfo.EnvString()
	dir, dirErr := resolveDataDir(env)
	if dirErr != nil {
		fail("data dir resolvable ("+env+")", dirErr.Error())
	} else if err := probeWritable(dir); err != nil {
		fail("data dir writable ("+env+")", dir+": "+err.Error())
	} else {
		pass("data dir writable ("+env+")", dir)
	}
```

对比 `run.go:155`、`models.go:109`、`providers.go:80` 三处**同一个动作都被 `if …dataDir == ""` 包着**：
`if s.dataDir == "" {`（`run.go:155`）、`if io_.dataDir == "" {`（`models.go:109`、`providers.go:80`）。
⇒ **`doctor` 是唯一一枚"无论调用方注入什么都去解析落点"的生产入口**，
而 `fail()` 在 `doctor.go:34-36` 里 `critical: true` ⇒ `cmdDoctor()` 返回 false ⇒ `main.go:88-89 os.Exit(1)`
⇒ `scripts/build.ps1:163 if ($LASTEXITCODE -ne 0) { Fail 'wisp doctor reported FAIL.' }`。
**这一条就是 §3/§4 里 (a) 支的全部风险来源。**

`probeWritable`（`doctor.go:300-309`，行首 0 tab 起）确认会**真的建目录 + 写文件 + 删探针**：

```go
func probeWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	probe := filepath.Join(dir, "doctor-write-probe.tmp")
```

### 2.2 链 B｜类型化族（**不决定**落点分叉，因为它没那枚短路）

```
WISP_ENV -> [读点 R2] internal/buildinfo/env.go:35-40 ResolveEnv()   （严格校验：nonsense 直接报错）
  +-> cmd/wisp/secret.go:184        -> :189 resolveSecretLayout(env)
  |        :130 root, err := userConfigDir()      <== **无条件 OS 读，没有 test 早退**
  |        :132 return … dataDirUnresolved128     拒
  |        :148 root = proc.SealableRoot(root)
  |        :155 return sessionLayout(env, root, exeDir)
  |               :106 func sessionLayout -> :107 proc.LayoutFor(env, configRoot)
  |                      （sessionLayout 的 :99-101 注释自陈 configRoot/exeDir 是**参数**，
  |                        "so a test can point dev/test/prod at three temp dirs"）
  +-> cmd/wisp/resident_windows.go:27 -> :33 proc.Boot(env)
  |        internal/proc/boot_windows.go:64 枚举止检 -> :77 layout, err := DefaultLayout(env)
  +-> cmd/wisp/slo_windows.go:238-240 [读点 R3] -> :243 ResolveEnv() -> :248 proc.Boot(env)

[链 B 的公共尾] internal/proc/envfork.go:235 func DefaultLayout(env buildinfo.Env) (Layout, error) {
        :236   dir, err := os.UserConfigDir()            <== OS 读**无条件执行**
        :238   return Layout{}, fmt.Errorf("proc: user config dir: %w", err)
        :245   l, err := LayoutFor(env, SealableRoot(dir))
        :77-85 case buildinfo.EnvTest:  DataDir: TestDataDir(), MutexEnabled: false   <== 注进来的 root 被**丢掉**
        :66-76 case buildinfo.EnvDev:   <root>\wisp-dev + `Local\wisp-dev-single-instance` + MutexEnabled: true
        :57-65 case buildinfo.EnvProd:  <root>\wisp     + `Local\wisp-single-instance`  + MutexEnabled: true
        :86-88 default:                 unknown env 报错
```

⇒ **一句话给下一位**：想找"env 把 OS 读短路掉"，只有 `doctor.go:256` 那一枚；
`proc` 那侧的 `test` 分叉（`envfork.go:77`）**问了 OS 再丢答案**，属于"换落点不换腿"，
它**不会**因为 (a) 变成装饰——它两形都真问 OS。这是两族最要紧的差别。

### 2.3 链上还有第三枚"早退"，比 `test` 更早（别在裁定时漏了它）

`cmd/wisp/doctor.go:248-255` 的 `portable.txt` 分支在 `:256` **之前** return。
⇒ 若 self-hosted 那台机器上 `build\wisp.exe` 旁边落过一枚 `portable.txt`，
`doctor` 连 `test` 短路都不会走到，直接答 `<exeDir>\data-dev`。
`cmd/wisp/dataroot_128_test.go:76-83 assertNoPortableOverride128` 在**测试**里替这件事兜底
（`:80` 的报错原文：`premise broke: a portable.txt sits next to the test binary (%s), so resolveDataDir takes the portable branch and never asks the OS for a config dir - these cases would assert nothing`），
但 **`build.ps1:162` 那次 `wisp.exe doctor` 前没有任何东西检查 `build\` 里有没有 `portable.txt`**（现量：`grep -n 'portable' scripts/build.ps1` → 0 命中）。
⇒ 这一条是 (a) 支的**第二个不可静态定的余量**，见 §3.4。

### 2.4 §2 的枚数小结

| 环节 | 枚数 | 名字 |
|---|---|---|
| 读 `WISP_ENV` 的生产码 | **3** | `buildinfo/env.go:36`、`buildinfo/buildinfo.go:38`、`slo_windows.go:238` |
| `env == "test"` 型短路 | **2**（且**分属两族、语义不同**） | `doctor.go:256`（字符串族，**短路 OS 读**）、`proc/envfork.go:212`（类型化族，读的是 `WISP_TEST_DATA_DIR`，**不短路 OS 读**） |
| 被短路掉的 OS 读 | 1 | `doctor.go:259 userConfigDir()` |
| 无守卫的 `EnvString()` 消费者（CI 上真会走到） | **1** | `doctor.go:104-105` |
| 有守卫的（测试注入即免疫） | 3 | `run.go:155/156`、`models.go:109/110`、`providers.go:80/81` |
| 更早的一枚旁路早退 | 1 | `doctor.go:248-255` portable.txt |

### 2.5 §2 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
02bbb59 docs(evidence/140 AC#2 r1): §2 环境枚举的入口链（两条族，逐格接上）

docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

（那次 commit 前 `git diff --cached --name-only` 现量**只有这一枚路径**。）

---

## §3　③　"取消 job 级 env" 对 §1 每一处判定的影响面（会换 / 不会换 / 不确定）

**判定口径**：「**会换**」= 该处的**分支选择**随 job env 消失而改变（不管颜色变不变）；
「**不换**」= 分支与颜色都不变；「**换分支不换色**」= 分支变了、可观测结果不变
——票 140 说的这一族才最危险，本程把它单列。

### 3.1 逐处判定（生产码 3 枚读点 + 6 处判定点 + CI 步骤）

| # | 判定点 file:line（版本同 §1） | 现在（job env = `test`）走哪支 | (a) 之后（env 未设 -> `DefaultEnv=dev`）走哪支 | 判 |
|---|---|---|---|---|
| 1 | **`cmd/wisp/doctor.go:256`（经 `doctor.go:104-105`）** | `return proc.TestDataDir(), nil`（`:257`），**`:259` 的 OS 读不执行** | 落到 `:259 userConfigDir()` -> `:263 SealableRoot` -> `:264-267` 答 `%APPDATA%\wisp-dev` -> `:108 probeWritable` **`MkdirAll` + 写 `doctor-write-probe.tmp` + 删** | **会换**（唯一一处 CI 上真被触发的） |
| 2 | `cmd/wisp/slo_windows.go:238-239` | 条件 `os.Getenv("WISP_ENV") == ""` 为 **false** -> **不执行** `:239` | CI 上**仍为 false**——`scripts/slo-check.ps1:111 $env:WISP_ENV = 'test'` 在**同一进程内先设**，子进程 `wisp.exe slo` 继承它（`:563`/`:500` 都是从 slo-check 里起的） | **CI 不换**。**但**：直接跑 `wisp slo` 而不经 slo-check（本机手测形）时会**换分支不换色**（`:239` 自己补 `test`，`:243` 结果仍是 `test`） |
| 3 | `cmd/wisp/secret.go:184 ResolveEnv` -> `:129 resolveSecretLayout` | env=`test` -> `LayoutFor` 的 `:81-85` case -> `DataDir: TestDataDir()` | env=`dev` -> `LayoutFor` 的 `:66-76` case -> `%APPDATA%\wisp-dev` + **`MutexEnabled: true`** | **CI 不换**：CI 上没有任何步骤跑 `wisp secret`（`grep -n 'wisp secret' .github/workflows/ci.yml` -> 0 命中）。**且**受影响的测试自己钉（§3.2 第 5 行）。注意 若实现程改成"顺手加一步 `wisp secret` 冒烟"，这一处立刻从"不换"变"会换" |
| 4 | `cmd/wisp/resident_windows.go:27 ResolveEnv` -> `:33 proc.Boot` | `test` -> `boot_windows.go:77 DefaultLayout` -> 注册**零枚** mutex（`envfork.go:84 MutexEnabled: false`） | `dev` -> `Local\wisp-dev-single-instance`（`envfork.go:70`）**会真注册** | **CI 不换**：CI 无步骤跑裸 `wisp`（常驻腿只在测试里以**子进程 + 自带 `WISP_ENV=test`** 起：`resident_sink_nail_127_windows_test.go:199`；`dataroot_128_windows_test.go:113` 反向钉 dev） |
| 5 | `cmd/wisp/run.go:155-156`（`buildEnvString()`） | 守卫 `:155 if s.dataDir == ""` **为 false**（测试注入） -> 整块跳过 | 同左，**守卫先挡住，env 根本没被求值** | **不换** |
| 6 | `cmd/wisp/models.go:109-110` / `providers.go:80-81` | 同 #5 | 同左 | **不换** |
| 7 | `cmd/wisp/doctor.go:43` 与 `main.go:135`（**只印**） | 印 `WISP_ENV=test` | 印 `WISP_ENV=dev` | **不换分支**，但**输出字节变**。⇒ 现量：全仓**没有**任何测试/脚本断言这两行的内容（`grep -rn 'data dir writable' --include='*.go' --include='*.ps1' --include='*.sh' --include='*.yml' .` 只命中 `doctor.go:109/:111` 两枚**生产者**；`grep -rn 'data dir rules'` 与 `'wisp build'` 在测试里 0 命中）⇒ **文案变无人连坐** |
| 8 | `internal/buildinfo/env.go:36` / `buildinfo.go:38` 本身 | 返回 `test` | 返回 `DefaultEnv` | 两枚都**不是决策点**，只是求值。**注意语义差不等于 (a) 的风险**：`ResolveEnv` 对非法值报错、`EnvString` 原样放过——(a) 让二者都走 default 分支，取值合法，**不报错** |
| 9 | `internal/proc/*`（`LayoutFor`/`DefaultLayout`/`Summarize`/`TestDataDir`/`boot_windows.go:64`） | — | — | **不换**：`internal/proc` 全包没有一处 `os.Getenv("WISP_ENV")`（§1.1 三条读者里没有它；上一份盘点 §0.2② 已裁同一件事，本程不重抄，只标"与它一致"） |

**CI 步骤侧的净账（把 #1 摊到步骤上）**：

| job | env 行 | 该 job 里**真正依赖**它的步骤 | 判 |
|---|---|---|---|
| `test-core`（`:224`，ubuntu） | `:227` | **零枚**。`:264` 跑 `runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`、`:288` `--scope=core` ⇒ 全是 `go test`，#9 说不换 | **不换** |
| `test-windows`（`:334`，windows-latest） | `:337` | `:400`（`build.ps1 -Env dev` -> `build.ps1:162 wisp.exe doctor`） | **会换**（#1） |
| `slo-smoke`（`:479`，windows-latest） | `:482` | `:497`（同 `build.ps1`） | **会换**（#1） |
| `slo-full`（`:537`，**`[self-hosted, wisp-slo]`**） | `:540` | `:560`（同 `build.ps1`） | **会换**（#1，**最重**：真机、跨 run 复用工作树） |
| `lint`（`:65`）/ `lint-frontend`（`:611`） | **无 job env** | — | 本来就不设，(a) 与它们无关 |

⇒ **"会换分支"的枚数：判定点 1 枚（`doctor.go:256`），CI 触发面 3 枚步骤（`:400` / `:497` / `:560`）。
另有 1 枚（`slo_windows.go:238`）在 CI 上不换、在本机直跑形上是"换分支不换色"的现成活例。**

### 3.2 测试侧：会因 (a) 换分支的枚数 = **0**（逐处给凭据）

| # | 测试 | 为什么不动 |
|---|---|---|
| 1 | `cmd/wisp/dataroot_128_test.go:281` 那族 5 条腿 | 它**自己**在 `:323 pinEnvThatAsksTheOS128(t)` 里 `:260` 栽 `test` → `:261` pin 回 `dev` → `:263` 自证。**job env 被 `t.Setenv` 覆盖两次，(a) 摘的是 ambient，摘不到 `t.Setenv`** |
| 2 | 同文件 `:151` / `:192` / `:209` | env 是**字面量入参**（`:158 resolveDataDir(env)`、`:195 resolveDataDir("test")`），压根不求值环境 |
| 3 | `cmd/wisp/dataroot_128_windows_test.go:113 return append(kept, "WISP_ENV=dev")` | 子进程 `cmd.Env` 自带，且**反向**钉 dev。`:108` 还把 `WISP_TEST_DATA_DIR` 摘掉 |
| 4 | `secret_argv_windows_test.go:236/337/371`、`early_log_nail_130:123`、`resident_sink_nail_127:199` | 子进程 `cmd.Env = append(os.Environ(), "WISP_ENV=test", …)`；同名键在 `exec.Cmd.Env` 里**后出现者胜** ⇒ ambient 被盖 |
| 5 | **`leg_sink_nail_131_windows_test.go:425/:457/:534` 三处 in-process `cmdSecret(...)`** | **这是 §3.1 #3 唯一的真实暴露面**——`cmdSecret` 会走 `secret.go:184 ResolveEnv()`。现量：`TestAC3SecretLegBooksItsAuditRecordsOnDisk` 在 **`:409 t.Setenv("WISP_ENV", "test")`**、`:524` 那个 `t.Run("secret")` 子项在 **`:526 t.Setenv("WISP_ENV", "test")`** ⇒ **两形都是 `test`，不换**。注意 但它是**全仓唯一一枚"依赖 job env 的值恰好等于自己钉的值"才成立的 in-process 生产入口调用**——这一枚的免疫是**靠它自己钉**，不是靠"生产入口没被调用" |
| 6 | `cmd/wisp/secret_test.go:932 TestSecretSameNameUnderThreeEnvsIsThreeBlobs`（`:935 t.Setenv("WISP_TEST_DATA_DIR", testData)` 却**没钉 `WISP_ENV`**） | 这是简报让我专门找的那一形。**现量结论：不换**——它**不经** `cmdSecret`/`ResolveEnv`，三枚 env 是**字面量**喂进去的：`:942 layout, err := sessionLayout(env, configRoot, "")`（`env` 来自 `:937 envs := []buildinfo.Env{EnvDev, EnvTest, EnvProd}`），落点再显式回填 `:961 newProbe(t, env, dirs[env], …)`。`:935` 那枚 `t.Setenv` 只为让 `LayoutFor` 的 `:83 TestDataDir()` 落在**声明过的**目录上（否则是 `%TEMP%\wisp-test-<pid>`），**与 `WISP_ENV` 无关** |
| 7 | `internal/buildinfo/env_test.go:22` | `:23/:28/:33` 三形自钉；`:34 want, err := ParseEnv(DefaultEnv)` 是**相对断言**（跟自己的默认比），不硬编码 `"dev"` ⇒ 连 `DefaultEnv` 变了都不红 |
| 8 | `cmd/wisp/` 其余 12 枚测试文件 | 上一份盘点 §4.1 已给"要么注入 `dataDir`/`env` 字段、要么自己 `Setenv`"的净数（**0 枚待钉**），本程**不重抄**；本程独立复核的是 §3.1 #5/#6 两枚**具体**暴露面 |

⇒ **§1.3 那 36 处测试命中里，(a) 之后会换分支的：0 枚。**

### 3.3 **简报点名要我专核的那一条**：票 123 那 4 枚的红，依不依赖 job 级 env？

**先厘清我核的是哪个问句**（票面 `AC#2` 前置约束 (i)/(ii)/(iii) 要的就是这一条）：
不是"它们为什么红"（那是 `128-ac4-r1-acceptance.md` §9#7 登记的**未追**项，本程不追），
而是"**取消 job env 会不会改变它们的颜色**"。这一条**静态链讲得清**，所以本程给结论、不推给读数。

**链（逐字节，版本：`run_test.go` = `cd011b8 09-20 20:42`，`run_mode101_test.go` = `7286297 09-21 16:51`）**：

1. 四枚红名（引 `128-ac4-r1-acceptance.md` §2.2 的 CI 原文清单，档位＝**引别人已裁读数**，非本程复现）：
   `TestComposedGateBlocksAWriteForTwoSeconds` + 三枚 `TestTicket101*`。
2. 前两枚的驱动路径：`run_test.go:118`

   ```go
   	return runTextTask(runSpec{
   ```
   参数里 `:122` 是 `		dataDir:   f.dir,`，而 `f.dir` 在 `:79` 是 `	f.dir = t.TempDir()`；
   config 写在**同一枚目录**：`:80` `	cfgPath := filepath.Join(f.dir, configFileName)`、`:101` `os.WriteFile(cfgPath, …)`。
3. 后三枚同形：`run_mode101_test.go:146-150` `code := runTextTask(runSpec{` … `		dataDir:     h.dir,`，
   `h.dir` 在 `:89` `h := &t101host{t: t, dir: t.TempDir(), base: srv.Base}`，config 在 `:135` `os.WriteFile(h.configPath(), …)`。
4. `cmd/wisp/run.go:155` 的守卫（逐字节，行首 1 tab）：`	if s.dataDir == "" {`
   ⇒ 注入非空 `dataDir` 时**整块跳过**，`:156 resolveDataDir(buildEnvString())` **一次都不执行**。
5. `runTextTask` 路径上**没有第二个 env 读点**（现量：`grep -n 'buildinfo\.' cmd/wisp/run.go` → **只命中 `:675`**
   那一行 `func buildEnvString() string { return buildinfo.EnvString() }`，而它只被 `:156` 调用）。
6. 四枚用例调的是 `runTextTask`，**不是** `main.go:127 cmdRun` ⇒ `main.go:128 printVersions`（唯一另一处 `EnvString()` 消费者）
   也不在路径上。

**⇒ 静态结论：票 123 那 4 枚的断言路径上没有任何一格求值 `WISP_ENV`。
(a) 支取消 job env 既不救它们、也不害它们——它们与这枚变量**不共因**，这一条**不需要读数就能定**。**
⇒ 直接落到票面前置约束的三条后果：**(i)** 别把"取消后它们仍红"记成 140 的失败（成立）；
**(ii)** 选 (b) 时它们不在本票分母里（成立）；**(iii)** 选 (c) 时登记清单里不许出现它们的名字（成立）。

**本程未定的那两半（不猜，如实挂着）**：
- 它们**真因**是什么：`128-ac4-r1-acceptance.md` §9#7 登记为「超出本格判据……票 123 / 票 140 的下一位」，本程未追。
  顺带给下一位一条**静态旁证**（不是结论）：`run_test.go:81-100` 写的那份 config **没有 `[risk]` 段**
  （对比 `run_mode101_test.go:110` `riskLines := "[risk]\nl1_window_sec = 1\nconfirm_timeout_sec = 1\n"`），
  ⇒ 这两族的**超时来源本来就不是同一枚**，把它们当一族可能就该拆。本程**不据此下判**。
- **同包共享态**这条耦合：`cmd/wisp` 是**一枚测试二进制**，`doctor.go:277 var userConfigDir = os.UserConfigDir`
  是**包级 var**、`t.Setenv` 有自动复原时序。(a) 经由"别的用例先跑了、留下状态"波及这 4 枚，
  **静态排不掉**——能排掉的只有"两形真跑的名册差集"（即 `AC#1` 要的那一发）。**这一条要读数才能定，本程未定。**

### 3.4 三处**不确定**（各自写清为什么不确定、谁能定）

| # | 不确定的那一处 | 为什么静态定不了 | 谁能定、要什么授权 |
|---|---|---|---|
| U1 | `doctor.go:104-105` 换支之后，那三枚 build 步骤（`ci.yml:400/:497/:560`）**颜色变不变** | 静态只能证"**落点从 `%TEMP%` 搬到 `%APPDATA%\wisp-dev`**"。红不红取决于**该 runner 上 `%APPDATA%` 是否可得且可写**——hosted 上大概率可得（大概率仍绿），**self-hosted 那台（`ci.yml:538`）无从静态推断** | 实现程：一发 `wisp.exe doctor` 的两形对照（本机即可，**不需要 CI**）。注意 但本程按红线**没跑**，也不建议在本机对 self-hosted 那台做写试验 |
| U2 | **`portable.txt` 会不会已经躺在那台机器的 `build\` 旁边**（`doctor.go:248-255` 的更早早退） | 那是**机器状态**不是仓库状态。`grep -n 'portable' scripts/build.ps1` 现量 **0 命中** ⇒ 没有任何东西在 `doctor` 前检查/清除它；测试侧的兜底（`dataroot_128_test.go:76-83`）只保**测试二进制旁边**，不保 `build\wisp.exe` 旁边 | 实现程/编排者：在 wisp-slo 那台 `ls build\portable.txt` 一发即定。本程**未查**（那是仓外真机） |
| U3 | (a) 之后 `%APPDATA%\wisp-dev` 的**跨 run 残留**会不会污染别的读数 | `ci.yml:560` 那枚 job 跑在真机且**跨 run 复用工作树**（凭据是 `scripts/slo-check.ps1:115` 的自陈注释，原文见 §5.4）。残留本身**不是红**（`probeWritable` 自己 `:308 os.Remove(probe)`），但同树里**攒下来的 `config.toml` / `secrets/`** 会不会被后面的读数当成"已有状态"读，静态不知 | 票 135／SLO 那本账（`128-ac4-r1-acceptance.md` §9#5 已把"真实落点/DPAPI 真写入"类四条件挂给票 135）。**要授权**：那是 owner 的真配置树，取数前需编排者批准 |

### 3.5 本节净数一句话

**会换分支：判定点 1 枚（`cmd/wisp/doctor.go:256`）× CI 触发面 3 枚步骤（`ci.yml:400`/`:497`/`:560`）；
外加 1 枚（`slo_windows.go:238`）在 CI 上不换、本机直跑形上"换分支不换色"。
测试侧会换分支：0 枚。不确定：3 处（U1/U2/U3，全部集中在 self-hosted 那台机器的状态，不在仓库里）。**

### 3.6 §3 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
8d1fbc5 docs(evidence/140 AC#2 r1): §3 取消 job 级 env 的逐处影响面（会换/不换/不确定）

docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

（那次 commit 前 `git diff --cached --name-only` 现量**只有这一枚路径**。）

---

## §4　④　三支各自的代价表

**要裁的原句（票面 `:32` 逐字节）**：

```
- [ ] **AC#2** 裁"收哪一支"并写代价，三选一：**(a)** 把那四枚 job 级 `WISP_ENV: test` 取消、改成**只有需要它的步骤**显式设；**(b)** 保持 env，但**每条受影响用例自己钉环境**（票 128 `c2fa2e9` 那枚 helper 的形状推广）；**(c)** 维持现状＋把"这条腿在 CI 上测的不是它声称的东西"做成**响亮失败**的显式登记。
```

### 4.0 先立一把"红会吃掉什么"的尺（本仓自己的原文，不是平台常识）

`ci.yml:257-261`（版本 `decb7b9 09-23 13:51`，逐字节）已经把这利害写死在仓里：

```
        # MOVED ABOVE the portable step (ticket 93, tightening, no step removed
        # and no threshold changed): it sat below a step that bare `go test` lets
        # go red, and GitHub stops at the first failed step - so on any run where
        # the portable step was red, this gate had literally never produced a
        # verdict, which is the A44-1 shape ticket 71 already fixed once in the
```

`ci.yml:327-328`（同一版，逐字节）：

```
  # WHY THIS IS NOT A DOOR REMOVED (the two forbidden shapes, both absent): no
  # step was deleted, and NO step here carries `continue-on-error`, whose entire
```

⇒ 两条合起来：**这一步红 ⇒ 后面所有步骤不产出 ⇒ 不产日志 ⇒ 那一格的读数这一发永久采不到，且没有
`continue-on-error` 可以兜**。⇒ 任何让 `ci.yml:400`（在 `test-windows` 里、位置在 `:422` **之上**）变红的改法，
**吃掉的是 `:422` 那格 `cmd/wisp` CLI 读数本身**——那一格正是票 128（已裁）／票 123（open）／票 135
（`Status: ready-for-agent`，`Packages` 行点名 `cmd/wisp/dataroot_128_test.go` 与 `cmd/wisp/doctor.go`）
三本账共同的落点。**这就是简报说的"整包红连坐后续步骤"的机制凭据。**

### 4.1 支 (a)｜取消四枚 job 级 env、改成只有需要的步骤显式设

| 项 | 内容 |
|---|---|
| **要动哪几枚文件** | **1 枚：`.github/workflows/ci.yml`**。形状：删 `:227`/`:337`/`:482`/`:540` 四枚，**并在三枚 build 步骤上各加一枚步骤级 `env:`**（`:395-400` cgo build smoke、`:496-497` Build wisp.exe、`:559-560` Build wisp.exe deps cached）。⇒ 净账 **4 删 3 加**（`:227` 那枚按 §3.1 是**纯赘余**，`test-core` 里没有一个步骤读它，可直删不补）。<br>**Go 面：0 枚**（§3.2 已给"测试侧会换分支 0 枚"）。<br>**脚本面：0 枚**（§1.4 现量 `*.sh` 零命中；`slo-check.ps1:111` 自己设）。 |
| **共享件纪律** | **`.github/workflows/ci.yml` 属共享件，改前要交回编排者——票面 `:52` 逐字节**：「`.github/workflows/ci.yml` 属**共享件**：本票 `AC#2` 若选 (a)，**改前先登记交回编排者**，不许自行改 workflow。」⇒ 本程（只读）**不提议、不改**；这一行是给实现程的**前置门**。 |
| **会不会把已经绿的改红** | **有这一支独一份的风险面**。静态能定的：`ci.yml:400`/`:497`/`:560` 三步的 `wisp.exe doctor`（`build.ps1:162`）从"答 `%TEMP%\wisp-test-<pid>`"换成"问 `%APPDATA%`、答 `%APPDATA%\wisp-dev`、`probeWritable` 真 `MkdirAll`+写+删（`doctor.go:300-309`）"。<br>**红不红＝§3.4 的 U1（本程未定）**。<br>**一旦红**：`:400` 红吃 `:422`/`:458`/`:474`（见 §4.0）；`:497` 红吃 `:500`（SLO smoke 门禁）；`:560` 红吃 `:563`（**SLO full，六态 + settle + leak，是这台机上最贵的一发读数**）。<br>**副作用不是红而是写面**：落点搬进真实配置树；`slo-full` 跑在 `[self-hosted, wisp-slo]`（`ci.yml:538`）且跨 run 复用工作树（`slo-check.ps1:115` 自陈，原文见 §5.4）。 |
| **还有两枚静态排不掉的余量** | **U2**：`doctor.go:248-255` 的 `portable.txt` 分支比 `:256` 更早 return，而 `grep -n 'portable' scripts/build.ps1` **0 命中** ⇒ (a) 之后到底走 `dev` 还是走 `<exeDir>\data-dev`，取决于那台机器 `build\` 旁的**文件状态**。<br>**U3**：残留树会不会被后面的读数当既有状态。 |
| **回退路径** | **一步**：`ci.yml` 只有 1 枚文件、改动是 4 删 3 加，**单独成一枚 commit** ⇒ 回退＝追加一枚把 `:227`/`:337`/`:482`/`:540` 加回去的 commit（**按仓纪禁 `--amend`/`reset`，已推送历史不改写，要更正就追加**）。<br>**但代码回不去的那半**：如果那一发已经**写进了真实 `%APPDATA%\wisp-dev`（或 self-hosted 那台的）**，回退 commit **不清理**它——清理属仓外真机操作，**要单独授权**，且本程建议由 owner 而非 agent 做（票 135 的 `Packages` 行也把 `doctor.go` 标成"只读其判定，不改语义"）。 |

### 4.2 支 (b)｜保持 env，每条受影响用例自己钉（推广 `pinEnvThatAsksTheOS128` 的形状）

| 项 | 内容 |
|---|---|
| **今天真正"该钉而未钉"的枚数** | **0 枚**——这一条上一份盘点 §4.1 已经裁过并给了凭据，**本程不重抄、不复算**。本程只补它没算的那一维（下一行）。 |
| **本程新给的一条：推广时的射程缺口** | 现成的名册尺 `cmd/wisp/dataroot_128_test.go:352 resolveDataDirConsumers128` 收的是**"调用了 `resolveDataDir` 的函数"**（判据在 `:370` 逐字节：`		if fn.calls["resolveDataDir"] {`）。<br>⇒ **它按构造收不到那三枚"直接求值 ambient env、不经 `resolveDataDir`"的生产入口**：`secret.go:182 cmdSecret`（`:184 ResolveEnv()`）、`resident_windows.go:24 runResident`（`:27`）、`cmdSlo`（`slo_windows.go:243`，另带 `:238` 自设兜底）。<br>该文件自己**知道**这件事——`:288-292`（逐字节）：`	// resolveSecretLayout is not a caller of resolveDataDir - it is the second,` / `	// independent reader of the same OS answer, and AC#1's third leg. It is named` / `	// here rather than left out, so the set below is checked for equality rather` / `	// than containment.` / `	want["resolveSecretLayout"] = true`<br>⇒ 它是**手写补进名册的一枚**，不是名册扫出来的。<br>**后果（给裁定用）**：上一份盘点 §4.1 提议的那条"名册里每一行都得先过 `pinEnvThatAsksTheOS128`"判据，**若照 `:352` 的射程写，会漏掉 `cmdSecret`**——而 `cmdSecret` 恰是全仓唯一一枚**被测试 in-process 调用、且真的走 `ResolveEnv()`** 的生产入口（`leg_sink_nail_131_windows_test.go:425/:457/:534`，见 §3.2 第 5 行）。它的免疫目前**来自 `:409`/`:526` 那两枚 `t.Setenv`**，不来自任何门禁。<br>⇒ (b) 的**真实分母**不是 0，是「0 枚待钉 + **1 枚名册尺需要扩射程**」。这一条本程只量面、**不建议改法**。 |
| **要动哪几枚文件** | `cmd/wisp/dataroot_128_test.go`（名册尺那半）＋**可能需要 1 枚新用例**补 ambient 生效断言（规范义务那半上一份盘点已登记，本程不重抄）。**都不碰 `ci.yml`** ⇒ 不触发共享件门。注意 但 `dataroot_128_test.go` **正在票 135 的 `Packages` 名单里** ⇒ (b) 与 135 **同文件相撞**，要先排串行。 |
| **会不会把已经绿的改红** | 对**既有**步骤：不改红（0 枚待钉，动的是测试自身仪器）。<br>对**新写的名册判据**：**有把既有绿改红的一条真路径**——若判据写成"名册里每个 entry 必须显式过 pin helper"，`cmdSecret` 那三处过的是 `t.Setenv("WISP_ENV","test")` 而**不是** `pinEnvThatAsksTheOS128`，按字面匹配就会**当场把 `leg_sink_nail_131` 判红**。这是"新仪器咬错人"，不是既有回归。<br>另一条形状冲突（本程现量、票面没写）：`pinEnvThatAsksTheOS128` 钉的是 **dev**（`dataroot_128_test.go:261`），而 131 那批要的是 **test**（`leg_sink_nail_131:409`/`:526`）⇒ **同一枚 helper 不能直接推广**，照抄必红。 |
| **票 123 那 4 枚** | **不在 (b) 的分母里**（票面前置约束 (ii) 要这一句）。凭据：§3.3 的静态链——它们的断言路径上根本没有 env 读点，"每条用例自己钉环境"对它们是**空操作**。 |
| **回退路径** | 一步：(b) 的改动全在测试文件、自成 commit ⇒ 追加一枚反向 commit 即回。无仓外写面，**无残留**。 |

### 4.3 支 (c)｜维持现状＋响亮失败登记

| 项 | 内容 |
|---|---|
| **要动哪几枚文件** | 登记面：**0 枚码**。台账 `docs/reports/pending-and-issues.md` 与本程禁改名单里的 `docs/reports/**` 重合 ⇒ **落点须由编排者定，不由 agent 自落**。<br>若要"响亮失败"落地成码：`ci.yml:250` 那枚步骤名对不上它跑的东西（现量：该步骤 `:264` 跑 `bash tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`，而那枚用例不看 ambient）——**要么换用例要么改名字**，两处都比现状响。改名字＝动 `ci.yml` ⇒ **重新触发共享件门**，与 (a) 同纪律。<br>外加 ambient 生效的正面断言（无论选哪一支都要落地；`docs/specs/**` 禁改 ⇒ 只能落在 `ci.yml` + 测试）。 |
| **会不会把已经绿的改红** | **登记本身：0 枚改红**（不动码）。<br>**风险在"响亮失败"的形状**：票面 `:44` 逐字节写着「**不许**用 `t.Skip`、不许改期望值换绿」⇒ 若实现成 skip 或翻转断言，就是把恒绿型装饰换成恒红型装饰，两样都白。<br>注意 现量补一条：那一步**已经有恒红型判决在前**——`128-ac4-r1-acceptance.md` §6.2 的标题句（档位＝**引别人已裁读数**）：「**是装饰，而且是恒红型装饰（比恒绿更贵）**」。⇒ (c) 若落在那一步，必须**同时**处理名字问题，否则登记与响亮各说一半。 |
| **票 123 那 4 枚** | **登记清单里不许出现它们的名字**（票面前置约束 (iii)）。支持该判据的凭据是同一条静态链（§3.3）：既然它们不读这枚变量，把它们登记进"因这枚变量而测的不是它声称的东西"就是**记错主人**。 |
| **回退路径** | 一步：登记是文本、改名是 1 枚 `ci.yml` 行。 |
| **本程量到的一条"这支可能在登记一件已经不成立的事"** | (c) 的登记命题是"这条腿在 CI 上测的不是它声称的东西"。但 §2/§3 的量面显示：`test-windows` 的 `:422` 那一格里 `cmd/wisp` 的加固腿**今天有牙**（`c2fa2e9` 已钉，`dataroot_128_test.go:263`/`:266-268` 两面自证），而**除它之外全仓再无第二枚走短路的用例**（§3.2 = 0 枚）。⇒ **(c) 现在能登记的只剩两样**：ambient 生效从未被断言那笔规范义务，与 `ci.yml:250` 那枚名字。**不再是"加固腿变装饰"**。这一句与上一份盘点 §4.1「(b) 单独开出去几乎是一枚空枪」是同一件事的两面；档位＝静态高置信、**可被 `AC#1` 两形真跑证伪**（"换分支不换色"那族的定义就是看起来一直绿）。 |

### 4.4 三支并排（只给本程量得出的那几列）

| | (a) 取消四枚 | (b) 每条自己钉 | (c) 现状＋登记 |
|---|---|---|---|
| 动的文件数 | **1**（`ci.yml`，4 删 3 加） | **1–2**（测试，含 1 枚新用例） | **0 码** ＋ 1 处可选 `ci.yml` 步骤名 |
| 碰不碰共享件门 | **碰**（票面 `:52` 明写先交回编排者） | 不碰 | 只在"改步骤名"那一支碰 |
| 会换分支的判定点 | **1**（`doctor.go:256`）＋ 1 枚换分支不换色（`slo_windows.go:238`，本机形） | 0 | 0 |
| 新增仓外写面 | **有**：真实 `%APPDATA%\wisp-dev`；self-hosted 那台跨 run | 无 | 无 |
| 改红既有步骤的可能 | **有**（U1 未定；红则连坐 `:422` / `:500` / `:563` 三格读数） | 无（除非新名册判据咬错人，见 §4.2） | 无 |
| 回退 | 1 枚追加 commit，**但仓外残留不清理** | 1 枚追加 commit，无残留 | 1 枚追加 commit |
| 与在飞票撞哪一枚文件 | `ci.yml` 共享件；`slo-full` 那格是 SLO 读数的宿主 | **`dataroot_128_test.go` 在票 135 的 `Packages` 名单里** | 与 135 不撞文件 |

### 4.5 §4 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
0ed0d02 docs(evidence/140 AC#2 r1): §4 三支代价表 + §3 回显

docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

（那次 commit 前 `git diff --cached --name-only` 现量**只有这一枚路径**。）

---

## §5　⑤　同族扫描：还有谁在用"环境把分支选走"这种形状

**票面 `next=` 点名要的那一条是"除 `cmd/wisp` 里为 128 钉的那枚 helper 之外还有谁"。
扫法全部写成可重跑的命令；**扫到 0 枚的地方也写命令与路径**。**

### 5.1 扫法一：其它包的 `t.Setenv("WISP_ENV", …)`

**命令**：`grep -rn 't.Setenv("WISP_ENV"' --include='*_test.go' .`
**现量：7 处命中 / 2 枚包。**

| 枚 | file:line（版本见 §1.3） | 钉的是 | 算不算同族 |
|---|---|---|---|
| 1-2 | `cmd/wisp/dataroot_128_test.go:260`、`:261` | 先 `test` 再 `dev` | **就是那枚 helper 本体**，票面已点名，不计新账 |
| 3-4 | `cmd/wisp/leg_sink_nail_131_windows_test.go:409`、`:526` | `test` | **算**。它们是"用 `t.Setenv` 把 Go 分支选走"的**既有第二例**，与 128 那枚**同形不同人**：钉 `test` 而不是 `dev`，所以 §4.2 那条"推广时照抄必红"的形状冲突就是它们造出来的 |
| 5-7 | `internal/buildinfo/env_test.go:23`、`:28`、`:33` | `test` / `nonsense` / `""` | **算，且是全仓唯一一枚"别包"的**。它是三形优先级自证（`:34 want, err := ParseEnv(DefaultEnv)` 是**相对**断言），**不依赖 ambient** ⇒ (a)/(b)/(c) 都不动它 |

**`cmd/wisp` 之外还有谁 `Setenv` 这枚变量**：`grep -rn 'os.Setenv("WISP_ENV"' --include='*.go' .` →
**唯一命中 `cmd/wisp/slo_windows.go:239`（生产码，不是测试）**。⇒ 没有第三枚包。

### 5.2 扫法二：依赖 `WISP_ENV` **默认值**的测试（这一形最容易在 (a) 后悄悄变义）

**命令**：`grep -rn 'EnvString\|ResolveEnv\|DefaultEnv' --include='*_test.go' .`
**现量：命中 17 行 / 3 枚文件**（`cmd/wisp/dataroot_128_test.go` 4 行：`:250`/`:263`/`:266`/`:268`；
`internal/buildinfo/env_test.go` 11 行：`:19`/`:20`/`:22`/`:24`/`:25`/`:29`/`:30`/`:34`/`:36`/`:38`/`:39`；
`internal/secret/store_test.go` 2 行：`:142`/`:144`——三枚文件的命中数**加总正好 17**，这条加法尺顺手核过）。逐条判：

| 处 | 判 |
|---|---|
| `internal/buildinfo/env_test.go:34` | **唯一一枚真的"依赖默认值"的**，但它比的是 `ParseEnv(DefaultEnv)` 本身（相对）⇒ **不硬编码 `dev`，免疫**。(a) 之后**默认值不变**（`DefaultEnv` 是 ldflags 注入的，见 §5.4 S6），所以这里连"变义"都不会发生 |
| `cmd/wisp/dataroot_128_test.go:263` | `if got := buildinfo.EnvString(); got != string(buildinfo.EnvDev) {` —— 它读的是**自己 `:261` 刚 pin 的值**，不是默认值 |
| `internal/secret/store_test.go:144 TestResolveEnvRef` | **名字骗人**：它测的是 `internal/secret/store.go:96 v, ok := os.LookupEnv(value)`——**配置里 `env:` 引用的那枚任意变量名**，与本票这枚无关（凭据：`grep -n WISP_ENV internal/secret/*_test.go` → 只有 `delete_test.go:159` 一枚**注释**） |
| 其余 | 注释/报错文案 |

⇒ **"依赖 ambient 默认值才成立"的测试：全仓 0 枚。** 这一条是 (a) 支最干净的一面。

### 5.3 扫法三：workflow / 脚本里"job/step 级 env 影响 Go 分支"的**其它例证**（票面要的那条主账）

**命令**：`grep -n 'env:' .github/workflows/ci.yml` → **只有 4 处、全是 job 级**（`:226 :336 :481 :539`），
**ci.yml 里一枚步骤级 `env:` 都没有**；另一枚 workflow
`grep -n 'env:' .github/workflows/slo-fresh.yml` → **1 处步骤级（`:67`）**。
再用 §2 的读者名单反查"这枚变量会不会被 Go 读到"，得**同族 6 枚**：

| # | 设值处（file:line + 版本） | 逐字节原文 | 谁读它、读来判什么 | 与 (a) 的关系 |
|---|---|---|---|---|
| **S1** | `scripts/slo-check.ps1:111-112`（`decb7b9 09-23 13:51`） | `$env:WISP_ENV = 'test'` / `$env:WISP_TEST_DATA_DIR = Join-Path $OutDir 'data'` | `internal/buildinfo/buildinfo.go:38` 与 `internal/proc/envfork.go:122` ⇒ 判**落点 + 枚举** | **最贴的同族，且方向相反**：它是"门禁自己设 env"，正是 §3.1 #2 说 `slo-check` 对 job env 免疫的原因 |
| **S2** | `scripts/build.ps1:116`（`bfcb230 09-20 07:00`） | `$env:CC = $cc` | `cmd/wisp/doctor.go:312 cc := os.Getenv("CC")`（`gccVersion()`，`:316` 拿它起 `exec.Command(cc, "--version")`）⇒ **判 `doctor` 的一枚 `critical` 检查**（`doctor.go:49-53`：`fail("gcc (build-time)", …)`） | **新账，且正落在 (a) 的风险步骤里**。⇒ `ci.yml:400/:497/:560` 那枚 `wisp doctor` **本来就有一枚靠脚本级 env 才成立的 critical 检查**。`grep -rn 't.Setenv("CC"' .` → **0 命中** ⇒ 没有任何测试兜它。这一条不改变 (a) 的分支判定，但**改变"这一步为什么可能红"的归因面** |
| **S3** | `.github/workflows/ci.yml:550`（`decb7b9`） | `      MINGW64_ROOT: E:\work\base\msys64\mingw64\bin`（行首 6 空格，挂在 `slo-full` 的 `:539 env:` 块里） | `scripts/build.ps1:51` PowerShell 分支 ⇒ 判 **gcc 发现路径**（`:59` 找不到就 `Fail`） | **形状完全同族**（job 级 env、选走一条分支、只在 self-hosted 那枚 job 上需要）。`:546-549` 的注释自陈不设它「the step died before compiling anything」⇒ **它是"job 级 env 是必需件而非装饰"的现成反例**，裁定 (a) 时别把它一起摘 |
| **S4** | `scripts/wisp-cli-tests.sh:109`（`8fe5c7c 09-21 21:05`） | `export PATH="$dll_dir:$PATH"` | 不是 Go 读的，是 **PE loader** 读的 ⇒ 判**测试二进制起不起得来**（`:103-105` 自陈两形：`pwd -W` 形 → `exit status 0xc0000135`＝票 98 的症状；shell 自己的形 → `PASS=33 FAIL=0 SKIP=0`） | 同族里的**极端例**：它选走的不是分支而是**整条腿的存在**。归票 98/111 的"有没有腿"账 |
| **S5** | `scripts/build.ps1:115/117/118/119` | `$env:CGO_ENABLED = '1'` / `$env:GOOS = 'windows'` / `$env:GOARCH = 'amd64'` / `if (-not $env:GOPROXY) { $env:GOPROXY = 'https://goproxy.cn,direct' }` | 构建期工具链，不选 Go **运行**分支 | 边缘同族（"脚本报 env 决定产物形状"）。注意 `:119` 是**条件设值**，与 `slo_windows.go:238` 同形 |
| **S6** | `scripts/build.ps1:107`（配 `:23-24`） | `"-X $BuildInfoPkg.DefaultEnv=$Env",` | `internal/buildinfo/buildinfo.go:20` 的 `DefaultEnv` ⇒ 判**env 未设时落到哪** | **不是 env 变量，是"构建期把默认分支选走"**。`:23-24 [ValidateSet('dev', 'prod')]` / `[string]$Env = 'dev'` ⇒ **`test` 不是合法构建默认**。⇒ (a) 的下家永远是 `dev`，这一条把 §3/§4 全部结论钉死 |
| **S7** | `.github/workflows/slo-fresh.yml:67-69`（该文件最后改动 `44ab500 09-23 14:04`） | 步骤级 `env:` 块里的 `GH_TOKEN:`（引用 secrets 的表达式，**本程按纪律只写变量名与 file:line，不抄任何值**） | `scripts/slo-freshness.sh`（`slo-fresh.yml:70` 起）判能否向 GitHub 查询 ⇒ **凭据可得性** | 形状同族（step 级 env 决定脚本分支），**与本票无关**，列出来只为让"还有谁"这一问的名册可核 |

### 5.4 两条支撑性引用（逐字节，给上面的判定当底）

`scripts/slo-check.ps1:114-117`（§4.1 的 U3 凭据；证明那台机器**跨 run 复用同一棵树**）：

```
# --- clear stale numbers ---------------------------------------------------
# The runner reuses E:\work\base\actions-runner\_work\wisp\wisp across runs, so
# build\slo can still hold the previous run's JSON when this run refuses to
# sample. "No numbers" has to mean no numbers on disk, not just no new ones.
```

`internal/proc/envfork.go:121-126`（S1 里那枚 `WISP_TEST_DATA_DIR` 的读者，也是链 A 的尾）：

```go
func TestDataDir() string {
	if dir := os.Getenv(TestDataDirEnv); dir != "" {
		return dir
	}
	return filepath.Join(SealableRoot(os.TempDir()), fmt.Sprintf("wisp-test-%d", os.Getpid()))
}
```

### 5.5 同族净数一句话

- 除那枚 helper 之外，**`t.Setenv("WISP_ENV")` 的同族：3 处**（`leg_sink_nail_131:409/:526` 钉 `test`、
  `internal/buildinfo/env_test.go:23/:28/:33` 三形自证 ⇒ 按**枚**是 5 处、按"另一例形状"是 **3 枚站点**：
  131 那两枚 + buildinfo 那一组）。注意 别把 §1.3 里那 6 处子进程 `cmd.Env` 算进来——那是"自带值"，不是"用 env 选分支"。
- **依赖 `WISP_ENV` 默认值的测试：0 枚。**
- **workflow/脚本里"环境把分支选走"的其它例证：7 枚（S1–S7）**，其中
  **S2（`CC` -> `doctor.go:312` -> critical 检查）与 S3（job 级 `MINGW64_ROOT`）是本次新扫出来的**，
  上一份盘点没记；S6（`DefaultEnv` 走 ldflags、且合法值里没有 `test`）是把 (a) 的下家钉死的那一枚。
- **扫法可重跑**：本节每条命令都已写在表格上方的粗体行里，路径全集＝仓库根（`.`），
  含 `.github/workflows/` 下**两枚** workflow（`ci.yml`、`slo-fresh.yml`——第二枚本程此前未见有人点名扫过）。

### 5.6 §5 那枚 commit 的回显（`git log --oneline -1` + `git show --name-only HEAD`，原样）

```
98665ef docs(evidence/140 AC#2 r1): §5 同族扫描（扫法可重跑）+ §4 回显

docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

（那次 commit 前 `git diff --cached --name-only` 现量**只有这一枚路径**。）

---

## §6　本程没查的档（不查的就是没查的，别按已查读）

| # | 没查的那半 | 为什么没查 | 谁能查 | 要查需要什么授权 |
|---|---|---|---|---|
| N1 | **`AC#1` 要的那一发真跑**（同一枚命令 `WISP_ENV` 未设 / 设为 `test` 两形各一发，四数 + 名册差集 + panic 计数） | 简报硬红线：另一枚程（`worker-ticket136-ac14`）此刻在本机做 `wisp slo -settle` 真取样，而本仓对"同一时刻既取样又编译"的处置是**整段判无效**（名单 `scripts/slo-check.ps1:153-155`，含 `go.exe`/`compile.exe`/`link.exe`/`gcc.exe`/`wisp.exe`）⇒ 本程**一枚仪器都没跑** | 票 140 的实现程 | 与 `cmd/wisp/**` 上在飞的终裁程**错开**（票面 `:77` 明写"按包 scope 读数，同包并行会造出假颜色"）；且要等 SLO 取样窗口空出来 |
| N2 | **§3.4 U1**：`(a)` 之后那三枚 build 步骤（`ci.yml:400`/`:497`/`:560`）红不红 | 需要真跑 `wisp.exe doctor`；且 self-hosted 那台的 `%APPDATA%` 可得性**不是仓库状态** | 实现程（本机一发两形即可覆盖 hosted 那一半）；self-hosted 那一半只有 CI 一发覆盖 | 本机发＝常规；CI 发＝**要编排者推**（本程只 commit 不 push） |
| N3 | **§3.4 U2**：`build\portable.txt` 在 wisp-slo 那台上存不存在（`doctor.go:248-255` 的更早早退） | 属**仓外真机文件状态**。本程全程未离开仓库目录 | 编排者／owner：`ls build\portable.txt` 一发即定 | 登进那台机＝要授权；或在 `ci.yml` 里加一枚只印不写的诊断步（＝动共享件，走票面 `:52` 那道门） |
| N4 | **§3.4 U3**：真实 `%APPDATA%\wisp-dev` 的跨 run 残留会不会污染 SLO 读数 | 本程未读**任何** owner 配置树，也未跑任何会写它的东西 | 票 135／SLO 那本账（`128-ac4-r1-acceptance.md` §9#5 已把"真实落点 / DPAPI 真写入"四条件挂给票 135） | 那是 owner 的真配置树 ⇒ **取数前需编排者批准**，且清理建议由 owner 做，不由 agent 做 |
| N5 | **CI 侧颜色**：run `35967768017` 的原文、以及 `c2fa2e9` 推送后 `:422` 那一格的颜色 | 简报明写**不许 `gh run view` 取日志**；`128-ac4-r1-acceptance.md` §9#8 也已把同一件事登记为"本程不 push 所以核不到，编排者：推完读那一步" | 编排者 | 推送权。凭据应绑 §7.1 的 S1/S2 形状（**引该文的结论时档位＝已裁读数，不是本程复现**） |
| N6 | **票 123 那 4 枚为什么在 CI 上红**（真因） | 不在本票判据内，且票面前置约束明写"别把票 123 那 4 枚一起结掉"。本程只裁了**它依不依赖 job env**（§3.3：不依赖，静态链完整），**没有**追它的红 | 另立票（前置约束 (ii)/(iii) 要求另立）。`128-ac4-r1-acceptance.md` §9#7 已登记同一件事 | 需一张新票；§3.3 末那条"两族的超时来源本来不同枚"的旁证可当起手，**但本程不担保它成立** |
| N7 | **`AC#3` 的两枚变异**（摘 `:261` 的 pin / 摘"先栽 test"那行 `:260`）在 runner 上响不响 | 是仪器验证，属实现程的活，且必须真跑 | 实现程／终裁程 | 同 N1 的排程约束。注意 本程给一条**静态旁证**帮它选形：`AC#3` 若问"这枚 helper 在别的机器上是不是变成装饰"，`internal/proc/envfork_test.go:323 TestDefaultLayoutUnknownEnvFails` 那枚"一枚断言两个可能来源"的形状比 `:192` 更像饵（此判据出自上一份盘点 §2.3，本程复核码面仍逐字成立，档位＝引别程＋本程静态复算） |
| N8 | **`AC#4` 的门禁四套**（`-count=2 -v` 四数、`gofmt -l`、`gofumpt -l`、`go vet` 宿主＋容器双 GOOS、`sh scripts/d22scan.sh`） | 全部是**编译/执行**仪器，N1 同因 | 实现程 | 同 N1。gofumpt 版本按票面 `:45`「盘上现量 v0.12.0，票面/派单里写 v0.7.0 的都是过期值」 |
| N9 | **`internal/observe/**` 里有没有新增读者**（AC#14 正在写那棵目录） | 本程**扫了**（`grep -c WISP_ENV -r internal/observe/` → **零命中**），但本程读到的是**工作树当前态**，不是钉版；那棵目录最后一次改动已经落到 `aef82f5 09-24 19:08`，**与本程取数几乎同时** ⇒ 名册对**这一枚**包有新鲜度风险 | 终裁表按提交时刻重跑一遍 §1.0 那四条 grep 即可 | 无需新授权（只读）。⇒ 本程在 §0.0 已声明："本文所有 file:line 取自盘上工作树"，**这既是它的方法也是它的边界** |
| N10 | **容器/Linux 腿**（票面 `:54` 那套 `golang:1.27` 挂载形状） | 简报禁 `docker` | 实现程（只在 Linux 腿真需要时） | 票面 `:54` 的挂载自证要求（Windows 形 `-v "D:\tmp\...:/wisp"` ＋ 容器内 `ls -l /wisp/go.mod` 非空；`/d/tmp` 那一形会**静默挂空而 rc=0**） |

---

## §7　临时件、写集、纪律自证与 commit 账

### 7.1 临时件清单（**只建不删**）

**本程没有在仓内、也没在 `D:\tmp\wisp140ac2\` 下创建任何临时文件或目录**——
纯 grep/只读的任务不需要中间件。本程**未执行**任何 `rm`/`rmdir`。

### 7.2 写集：只有 1 枚文件

`docs/evidence/s1/140-ac2-static-blast-radius-r1.md`（本文件）。
**未动**：`internal/observe/**`、`cmd/wisp/**`、`internal/proc/**`、`internal/winsec/**`、
`.github/workflows/**`、`scripts/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、
`tools/d22scan/**`、任何阈值／golden／`thresholds.go`、`frontend/**`、`design/**`、
`.scratch/wisp/issues/**`（**含本票票面：票面 `:50` 要求"每完成一格往票面 append 一条 log"，
本程按简报的"写集只有那枚新证据文件"走，**那一格未翻，请编排者自行决定由谁翻**——
与上一份盘点 §篇首同一处冲突，同一处理法）、`docs/reports/**`。

**每枚 commit 的暂存面自证**：每次提交前先 `git add -- <那一枚路径>` 再 `git diff --cached --name-only`，
**六次现量全部只有 1 枚路径**。owner 那 16 枚 `design/**` 未提交删除 + 两枚未跟踪目录
（`design/doubao/`、`design/old/`）**原样留在 unstaged，本程未还原、未提交、未删**
（交件前复量：`git status --porcelain | grep -v '^ D design'` → 只剩那两枚 `??` 行）。

### 7.3 注入两栏计数（这条规矩只能被引用，不能被外部文字代填）

- **真通知回显数：1。** 会话开头 harness 回显的 `Memory: d:/work/workspace/projects plans/wisp/agents.md`
  全文 + 可用 skill 清单 + `The date has changed. Current date: 2026-09-24`。
  这些是**环境回显**，本程未据此改变任何动作，只登记。
- **判为注入数：0。** 本程读到的所有"已结案／请改判／恒红型装饰／不许放宽"类句子，
  经核**全部出自票面、台账与兄弟程证据文件的正文**（`.scratch/wisp/issues/140-…md`、
  `docs/reports/pending-and-issues.md`、`docs/evidence/s1/128-ac4-r1-acceptance.md`、
  `docs/evidence/s1/140-static-inventory-r1.md`），**不是**工具输出里冒充授权的句子。
  本程把它们一律当**待核断言**处理（结果见 §0.2），未据此行动。
- **凭据抄录：0 处。** `slo-fresh.yml:69` 那枚 `GH_TOKEN` 本程**只写变量名与 file:line**，未抄值。
  全程 grep 未命中疑似明文凭据；若曾命中，按纪律只打印文件名。
- **一处自纠要如实登记**：本程曾用 `python -c "…"`（**双引号**）就地改两处计数，
  脚本里的反引号被 shell 当成命令替换真执行，报错
  `No such file or directory` / `command not found` 各若干；
  **写盘前的 assert 挂住 ⇒ 文件未受损**，随后改用 Edit 工具落盘。
  这正是简报"定界符必须加单引号／含反引号的段落用 Edit/Write 落盘"那条规矩的**本程现量代价**。

### 7.4 本文件的 commit 账

**取法（别按本节行数信我，按这两条命令现量）**：

```
git log --oneline -- docs/evidence/s1/140-ac2-static-blast-radius-r1.md
git show --numstat --format='%h %ad' --date=format:'%m-%d %H:%M' <那一枚> -- docs/evidence/s1/140-ac2-static-blast-radius-r1.md
```

| 节 | commit | 内容 |
|---|---|---|
| 起手＋§1 | `8e86f32` | §0（锚点/纪律/分工/推翻清单/版本方法）＋ §1 读者名册 |
| §2 | `02bbb59` | 环境枚举入口链（两族） |
| §3 | `8d1fbc5` | 逐处影响面 + 票 123 专核 + 三处不确定 |
| §4 | `0ed0d02` | 三支代价表 + 尺（含 §3 回显） |
| §5 | `98665ef` | 同族扫描 S1–S7（含 §4 回显） |
| §6＋§7 | 本表所在的那一枚 | 没查的档 + 纪律自证 |

注意 **自指余留**：本表最后一行是在它自己那枚 commit **之前**写的 ⇒ "§6/§7 落在第 6 枚"这句话
要靠第 6 枚本身才成立。⇒ **数本文件的 commit 一律按上面那条 `git log --oneline -- <路径>` 现量。**

### 7.5 交付净数（给 `AC#2` 当分母的那一句）

- **① 读者**：全仓 `WISP_ENV` 命中 **69 处 / 19 枚文件**；**真读者 3 处 / 2 枚文件**
  （`buildinfo/env.go:36`、`buildinfo/buildinfo.go:38`、`slo_windows.go:238`）；
  判落点的分叉 **1 处**（`doctor.go:256`，grep 单查 `WISP_ENV` 抓不到）。
- **② 链**：两族。**字符串族**（`EnvString` -> `doctor.go:247 resolveDataDir`）有那枚短路；
  **类型化族**（`ResolveEnv` -> `proc.LayoutFor`）**无条件问 OS 再把答案丢掉**（`envfork.go:236` -> `:245` -> `:77-85`）。
  链上唯一"CI 真会走到且无守卫"的消费者是 `doctor.go:104-105`。
- **③ 会换分支**：生产判定点 **1 枚** × CI 触发面 **3 枚步骤**（`ci.yml:400`/`:497`/`:560`，
  全部经由 `build.ps1:162` 那次 `wisp.exe doctor`）；测试侧 **0 枚**；
  另 1 枚"换分支不换色"（`slo_windows.go:238`，CI 上因 `slo-check.ps1:111` 先设而不换）。
  **票 123 那 4 枚与 job env 不共因——这一条静态链完整，不需要读数**（§3.3），真因仍未追（N6）。
- **④ 三支**：(a) 1 枚文件 4 删 3 加、**碰共享件门**、有把 `:400` 改红从而吃掉 `:422`/`:500`/`:563` 三格读数的路径、
  回退清不掉仓外残留；(b) Go 面 0 枚待钉但**名册尺有射程缺口**（`dataroot_128_test.go:370` 只认 `resolveDataDir` 的调用者，
  收不到 `cmdSecret`），且 `pinEnvThatAsksTheOS128` 钉 dev、131 那批要 test ⇒ **不能直接推广**；
  (c) 0 枚改红，但**它要登记的那件事在 `cmd/wisp` 里今天已经不成立**（唯一的靶已被 `c2fa2e9` 钉住）。
- **⑤ 同族**：`t.Setenv("WISP_ENV")` 7 处 / 2 枚包（除那枚 helper 外，新账是 131 那 2 处 + buildinfo 那 3 处）；
  **依赖默认值的测试 0 枚**；workflow/脚本"环境把分支选走"的其它例证 **7 枚 S1–S7**，
  其中 **S2（`CC` -> `doctor` 的一枚 critical 检查）与 S3（job 级 `MINGW64_ROOT`）是本程新扫出来的**。

### 7.6 收尾两处要如实报回的（不是本票的账，但都是量出来的）

**(1) 一处仪器事实，按 `AGENTS.md` 自己的规矩上报，不当授权用、不改任何码。**
`AGENTS.md §1.2` 把 ban #8 的范围写成「`U+2190–U+2BFF`、`U+1F300–U+1FAFF`、`U+FE0F`」。
现量那枚仪器本体（`tools/d22scan/main.go`，冻结件，本程只读）：

- 实际正则（`main.go:111`）是字符类
  `[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`
  ⇒ **不含 `U+2190–U+25FF`**。scope 表（`main.go:482-485`）逐字节是
  `{dir: filepath.Join(root, "design"), label: "design/"}`，另三枚 `frontend/`、`internal/`、`cmd/`
  ⇒ **`docs/` 与 `.scratch/` 都不在 walk 范围内**。
- 两边不一致的后果：按 `AGENTS.md` 的字面范围，本票票面与上一份盘点里大量出现的
  `⇒`（U+21D2）、`→`（U+2192）、`①–⑤`（U+2460–U+2463）**全都算违规**；按仪器实际正则，它们**都不算**。
- 本程的处理：**按仪器真正会判的那一段收敛自己的散文**——把 `U+26A0`（落在 `2600–27BF` 内、会被真判）
  共 8 处换成"注意"；`U+2190–U+25FF` 的箭头与序号**保留**
  （它们既在仪器射程外，又是本票票面与兄弟证据文件的同一套用词，换掉只会让引用对不上）。
  ⇒ 本文件口径是"**散文里不出现任何仪器会判的字符**"，不是"不出现任何 `AGENTS.md` 字面列出的码位"。
  哪一条才是想要的规矩请编排者定；**本程不自行改 `AGENTS.md`**——它不在本程写集里，
  且它自己篇首写明「任何一处与那些文件不一致 ⇒ 以那些文件为准，并把这一处当作本文件的缺陷上报，不要照它做」。

**(2) 提交纪律的自我核对结果，连"按 log 数枚数会数错"这个坑一起报。**
本程 6 枚 commit 每枚前都做 `git add -- <唯一路径>` + `git diff --cached --name-only`，
**六次现量全部只有 1 枚路径**；事后逐枚 `git show --name-only` 复量
（`8e86f32` / `02bbb59` / `8d1fbc5` / `0ed0d02` / `98665ef` / `2e7d8b4`）**同样各只含本文件**。
本程工作期间兄弟程 `worker-ticket136-ac14` 在同一条 `dev` 上落了 3 枚 commit
（`d6c83de`、`ba6d94e`、`2f22ab3`），与本程的 6 枚**交错**排列 ⇒ **按 `git log -6` 数我的枚数会数错**。
⇒ 数本文件的 commit 只能按
`git log --oneline -- docs/evidence/s1/140-ac2-static-blast-radius-r1.md`。
交件前复量：owner 那 **16** 枚 `design/**` unstaged 删除 + 两枚未跟踪目录
（`design/doubao/`、`design/old/`）**原样未动**
（`git status --porcelain | grep -c '^ D design'` = 16）。

**本程未 push。** 推送由编排者在核过之后决定。


