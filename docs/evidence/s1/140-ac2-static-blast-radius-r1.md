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

