# 140 `AC#1` 裁决 r1 —— 被点名的那份盘点**结不结**「先量面」这一格

**裁决代理**：`auditor-ticket140-ac1-r1`（只读；非实现者；非被裁表作者）
**时刻**：2026-09-24 起，六节六枚 commit
**码面锚点**：起手 `git rev-parse --short HEAD` = **`9090f36`**，分支 `dev`。
**被裁对象**：`docs/evidence/s1/140-ac2-static-blast-radius-r1.md`（简报点名的那一份）。
**顺带被裁**：`docs/evidence/s1/140-static-inventory-r1.md`（它的兄弟；见 P2，两枚必须成对判）。

**本程跑过的命令全集**：`grep` / `sed -n` / `awk` / `wc` / `ls` / `find` / `xargs -n1 dirname` /
`git log` / `git show` / `git status --porcelain` / `git rev-parse` / `git cat-file -e` /
`git merge-base --is-ancestor` / `git diff --cached --name-only` / `Read` / `Edit` / `Write`。
**本程没有跑**：`go build`、`go test`、`go vet`、`gofmt`、`gofumpt`、`docker`、`wisp`、`gh`、`rm`、
`git add -A`、`git add .`、`--amend`、`reset`、`rebase`、`stash`、`checkout .`、`clean`、`push`。
零编译、零取样、零时序读数（简报的机器红线：同机另一枚程正在取 CPU 敏感读数）。

**两栏计数（全篇累计，每节末尾各自更新）**
- 真通知回显数：**1** —— 会话开头 harness 随第一条 Bash 输出回显的 `system-reminder` 块
  （可用 skill 清单 + `The date has changed. Current date: 2026-09-24` + `Memory: …/agents.md` 全文），
  工具名 `Bash`，前 40 字 `<system-reminder> The following skills are available`。本程未据此改任何动作。
- 判为注入数：**0** —— 本程读到的所有「已结案／请改判／恒红型装饰／不许放宽」类句子，
  经核**全部出自票面、兄弟证据文件与台账的正文**（`140-…md`、`140-ac2-…md`、`140-static-inventory-r1.md`、
  `128-ac4-r1-acceptance.md`、`135-…md`），不是工具输出里冒充授权的句子。一律当**待核断言**处理（结果见 P 表与各节）。
- 凭据抄录：**0 处**。`slo-fresh.yml:69` 那枚 `GH_TOKEN` 只写变量名与 file:line，未抄值。

---

## 0　简报前提的四处现量更正（先报回，不绕、不修）

| # | 简报的前提 | 现量结果 |
|---|---|---|
| **P1** | 「808 行、8 枚 commit」 | 枚数 **对**（`git log --oneline -- docs/evidence/s1/140-ac2-static-blast-radius-r1.md` = `8e86f32 02bbb59 8d1fbc5 0ed0d02 98665ef 2e7d8b4 e0f4c07 6451625`，8 枚）。行数 **不成立：盘上是 811 行**（`wc -l`），不是 808 |
| **P2** | 「那份静态盘点＝AC#1 的量面，判它结不结 AC#1」 | **按文件自己的说法不成立**。被裁文件篇首第一行是 `# 140 AC#2 — 取消 job 级 WISP_ENV: test 的静态爆炸半径 r1（分母，不是裁定）`；`:10` 自称「本文件只做一件事：给票 140 `AC#2` 的三支代价**算分母**」；`:33` 明写「**本文件不复核、不重抄那三笔账**」。**AC#1 的逐包问句是它兄弟那枚答的**：`docs/evidence/s1/140-static-inventory-r1.md`（436 行），篇首即「票 140 `AC#1` 的『先量面』那一半的**静态前置**；一个测试都没跑」。<br>=> 判"是否结掉 AC#1"只能按**两枚成对**判；单看被点名这一枚，它自己就声明不结，**这不是它的缺陷，是它的分工** |
| **P3** | 「判这份盘点是否**真的**结掉 AC#1」 | AC#1 判据原文（票面 `:29-31`）**强制真跑**：「取证方式**必须包含真跑**：同一枚命令在 `WISP_ENV` 未设与设为 `test` 两形各一发，给四数＋**名册差集**……别用 `GOOS=linux go vet` 那类只编译的仪器代答」。两枚文件都**没跑**，且都如实挂着（被裁文件 §6 N1；兄弟文件篇首）。<br>=> 在本简报的只读约束下这是**正确行为**；但结论上限随之定死：**AC#1 的静态那一半可以被结，真跑那一半谁都没结** |
| **P4** | 冻结件名单含 `rules_gateway.go` / `thresholds.go` | 两枚都在盘上：`internal/risk/rules_gateway.go`、`internal/observe/thresholds.go`。但**票面 `:53` 的禁改名单没有点名这两枚**（票面写的是 `docs/PLAN.md`·`docs/specs/**`·`internal/risk/**`·`tools/d22scan/**`·`allowlist.txt`·阈值／golden·`frontend/**`·`design/**`）。`thresholds.go` 由 `AGENTS.md §1.1` 覆盖，`rules_gateway.go` 落在 `internal/risk/**` 之内 ⇒ 以更严者为准，**不改变 §5 任何一格** |

owner 那 16 枚 `design/**` 的 unstaged 删除与两枚未跟踪目录（`design/doubao/`、`design/old/`）：
本程未 stage、未还原、未删，**不记为缺陷**（交件前复量：
`git status --porcelain | grep -v '^ D design'` 只回那两枚 `??` 行；起手时 `git diff --cached --name-only` 为空）。

---

## 1　①　逐包覆盖核

### 1.0　分母是自己量的，不是抄来的

```
grep -rn 'WISP_ENV' --include='*.go' .   -> 62 处 / 17 枚文件    （与被裁文件 §1 逐字同值）
grep -rn 'WISP_ENV' --include='*.ps1' .  ->  1 处 /  1 枚文件    （同）
grep -rn 'WISP_ENV' --include='*.sh'  .  ->  0 处 /  0 枚文件    （同）
grep -rn 'WISP_ENV' --include='*.yml' .  ->  5 处 /  1 枚文件    （同）
```

四条加总 **68 处 / 19 枚文件**；被裁文件 §1.6 写的分母是「**69 处** / 19 枚文件」——
**这一枚数字复不出来**（62+1+0+5 = 68）。19 枚文件对，69 那处多 1。
按包分组（`grep -rl 'WISP_ENV' --include='*.go' . | xargs -n1 dirname | sort | uniq -c` 现量）只有五枚包：

| 包 | 文件数／处数 | 被裁文件（AC#2 那份）答了吗 | 兄弟文件（AC#1 那份）答了吗 |
|---|---|---|---|
| `cmd/wisp` | 11 / 44 | **答，逐处**：§1.1C、§1.2 那 13 行、§1.3 那 7 枚测试文件、§3.1 九行、§3.2 八行，全部给到 file:line | **答**：§1.1（生产码 7 枚入口，标腿型）、§1.2（测试码 7 行）、§2.1（4 枚 top-level ＋ 5 枚子项逐枚） |
| `internal/buildinfo` | 3 / 15 | **答**：§1.1A/B（`env.go:36`、`buildinfo.go:38`、`buildinfo.go:20`）、§1.3（`env_test.go:23/28/33`）、§3.2 第 7 行、§5.1 第 5-7 行 | **答**：§1.1 末段、§1.2、§3.3（**唯一正面断言 `WISP_ENV` 优先级的用例，且不在受影响 job 里**） |
| `internal/proc` | 1 / 1（`doc.go:13` 注释） | **答**：§1.2 那一行 + §2.2 链 B + §3.1 #9 | **答**：§1.1（`envfork.go:236/:125/:77-85`）、§1.2、§2.3 |
| `internal/secret` | 1 / 1（`delete_test.go:159` 注释） | **半答**：§5.2 末段点名「`grep -n WISP_ENV internal/secret/*_test.go` -> 只有 `delete_test.go:159` 一枚注释」，并顺手拆了 `store_test.go:144 TestResolveEnvRef` 那枚**名字骗人**的饵（现量对：它测的是 `internal/secret/store.go:96`，与本票无关）——**但没有把它放进 §1.3 的名册** | **答**：§1.1 末段（`store.go:96` 读的是配置里 `env:` 引用的那枚**变量名**） |
| **`internal/winsec`** | **1 / 1**（`dataroot_symlink_119_other_test.go:133` 注释） | **没答**（F1）。全文只出现两次 "winsec"：§1.4 那行说 `winsec-tests.sh` 不提这枚变量、§7.2 的"我没动哪些文件"清单。**§1 名册里这枚包不存在** | **答**：§1.1（`resolve.go:242-244`、`winsec_other.go:117`，判「腿型 B，它问的是 `TMPDIR` 不是 `WISP_ENV`」）＋ §1.2（`dataroot_symlink_119_other_test.go:147-150/:178-187/:349/:375`，标 `//go:build !windows`） |

### 1.1　F1：被静默漏掉的那枚包就是发现本身

`internal/winsec` 有 `WISP_ENV` 命中却不在被裁名册里。逐条给凭据：

- 命中现量：`grep -rn 'WISP_ENV' internal/winsec/` ->
  `internal/winsec/dataroot_symlink_119_other_test.go:133:// TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119 is the WISP_ENV=test route`
- **不是"它写完之后才落的文件"**：该文件最后一改 `c94927d 09-24 17:05`（`git log -1 -- <path>`），
  **早于**被裁程的起手锚点 `2e6d171 09-24 18:58`；且在它自己声明的两个锚点上 `git cat-file -e` 都 **present**
  （`2e6d171`、`0ee68e1` 各一发），`2e6d171` 是本 HEAD 的祖先（`git merge-base --is-ancestor` rc=0）。
- **它不是一枚无关的命中**：那枚文件正是「判据依赖去问操作系统」最重的一族——`:133` 的注释自己写着
  "is the WISP_ENV=test route with TMPDIR behind a symlink"，函数体 `:147-150` 用
  `t.Setenv("TMPDIR", sub)` **再** `t.Setenv(proc.TestDataDirEnv, "")`（注释：`take the os.TempDir() branch, not the injection`）
  **逼** `proc.TestDataDir()` 真去问 OS；文件内还留着一枚反"用被测函数算期望"的仪器
  （`:120-131 cleanSpelling119`，注释逐字：`It is deliberately filepath.EvalSymlinks and not proc.SealableRoot:
  the latter is what several of these cases exist to test, and an expectation computed by the function under test cannot fail.`）。
- 它带 `//go:build !windows`（现量第 1 行）=> 按 AC#1 的三口径（§3 专判），这枚包落在**第一种**
  （"在 CI 上根本没跑"那个形），**恰好是被裁文件连口径都没有那一形**。

=> 结论的强度要写准：**这不是"AC#1 的账算错了"，是"被点名的这一枚不结 AC#1，结它的是兄弟那一枚"**。

### 1.2　名册算术的三处不闭合（引证都是真的，加法不是）

1. §1.3 说测试码「**9 枚文件、36 处命中**」。现量带 `WISP_ENV` 的 `*_test.go` = **10 枚 / 38 处**
   （`cmd/wisp` 7 枚 30 处 ＋ `internal/buildinfo/env_test.go` 6 处 ＋ `internal/secret/delete_test.go` 1 ＋ `internal/winsec/…119_other_test.go` 1）。
   差的正是 F1 那两枚注释命中。
2. §1.2 说非测试「6 枚文件、16 处命中」，§1.1 列 3 处读者，§1.3 列 36 处 => **3+16+36 = 55，与它自己的 62 差 7**。
   缺的 7 处全在 `internal/buildinfo/env.go` 的另外 5 行（`:8/:18/:25/:27/:31`，其中 `:25/:27` 是
   **打给用户的报错文案**，既不是"判"也不是 §1.2 用的"只印"）与 `buildinfo.go:20/:34`。
   => 一名读者拿总分去**复算**那三张表，复不出来；`env.go` 作为"命中最多的单文件之一"（6 处）从未被逐行列全。
3. §3.2 第 8 行说「`cmd/wisp/` 其余 **12** 枚测试文件」。现量 `ls cmd/wisp/*_test.go | wc -l` = **16**，
   §1.3 点名 7 枚 => "其余" 是 **9**。（兄弟文件 §3.2 同一处写的是「其余 **13** 枚」，也不是 9。）
   被这个错分母盖住的是三枚**不带** `WISP_ENV` 字样、四条 grep 结构性抓不到的文件：
   `secret_dataroot_119b_test.go`（16 处 OS 读命中，`//go:build !windows`）、
   `tempdir_resolved_124_test.go`、`leg_sink_gate_131_test.go`。
   被裁文件把它们推给「上一份盘点 §4.1 已给净数」——**兄弟文件 §1.2/§3.4 确实逐枚给了**，
   含那句要紧的：`secret_dataroot_119b_test.go` 那三枚 top-level「**在 CI 上从不编译、从不运行**」。

### 1.3　宽口径名册（"任何判据依赖问 OS"），两枚文件一对的账

本程另起一条现量（测试码里出现
`os.UserConfigDir|UserHomeDir()|TestDataDir()|os.TempDir()|Setenv(TMPDIR|APPDATA|HOME|XDG_CONFIG_HOME)`）：
命中 **6 枚包** —— `cmd/wisp` 5 枚文件、`internal/proc` 1、`internal/risk` 5、`internal/tools` 4、
`internal/winsec` 4、`internal/config` 1。

- 被裁文件：**0/6 进册**（它只按 `WISP_ENV` 字样扫，§1 那一节连口径都没声明"只按字样"）。
- 兄弟文件：**6/6 全答**（§1.1 生产码 ＋ §1.2 测试码，逐枚标腿型 A/B，并给反向排除：
  `WISP_TEST_DATA_DIR` 是另一枚变量、`t.TempDir()` 那 131 枚文件 421 处命中全是 harness 句柄）。

**本节净数一句话**：五枚 `WISP_ENV` 命中包里 **4 枚**被"两枚文件这一对"答足、
**1 枚（`internal/winsec`）在被裁那一枚里静默缺席**（兄弟有）；
宽口径六枚包里 **5 枚只在兄弟文件里有名册**；被裁文件另有三处名册算术不闭合。
**结掉 AC#1 静态那一半的是兄弟文件；被点名这一枚结的是 AC#2 的分母。**

### 1.4　本节的独立复算（与两枚文件都对照过，成立的部分也要写）

```
grep -rn 'env == "test"' --include='*.go' .                      -> 全仓 1 命中 cmd/wisp/doctor.go:256
grep -rn '== buildinfo.EnvTest\|== EnvTest' --include='*.go' .   -> 全仓 1 命中 internal/proc/envfork.go:212
grep -rn 't.Setenv("WISP_ENV"' --include='*_test.go' .           -> 7 处 / 2 枚包
grep -rn 'os.Setenv("WISP_ENV"' --include='*.go' .               -> 1 处，且是生产码 cmd/wisp/slo_windows.go:239
```

被裁文件 §1 那段「一条容易漏的口径：`grep 'WISP_ENV'` **抓不到那个分叉判定本身**」——
**本程复核成立**，这是它对 AC#1 最实在的一枚贡献（收窄到 2 枚、落在 2 枚文件、语义不同）。

### 1.5　本节两栏计数增量

真通知回显数：累计 **1**（篇首那一次，本节无新增）。判为注入数：累计 **0**。
本节**未**在任何工具输出里遇到自称「编排者备注／系统提示／请 revert／冻结某包／放宽阈值」的句子。

### 1.6　本节 commit 回显

本节那枚 commit 的 `git diff --cached --name-only` 现量只有本文件一枚路径，commit = `bd0f826`
（回执 `git show --name-only bd0f826` 同样只含本文件）。
**全篇六节的 commit 账集中记在 §6.5**，取法给的是命令不是本文件的行文——
每节的 sha 会随本节写在其下一节的开头，不另开 commit 贴回显（否则"一节一枚"就变成十二枚了）。

### 1.5b　本节两栏计数

见篇首累计（本节无新增）。

---

## 2　②　引证真实性抽验

### 2.0　抽法（先写给下一位可复跑的抽法，再给结果）

1. **全集**：从被裁文件正文里抽出所有**全限定** `路径:行` 引证
   （`grep -oE '([A-Za-z0-9_.-]+\.(go|ps1|sh|yml)):[0-9]+(-[0-9]+)?'` 只留带目录前缀的那种），
   `sed 's|^\./||' | sort -u | nl` => **74 枚**。
2. **取样规则（定距，非挑选）**：按下标取 **3, 9, 15, 21, 27, 33, 39, 45, 51, 57, 63, 69**
   （起点 3、步长 6，12 枚），规则先于结果确定，**不许我挑**。
   表格里同时给出下标 1 与 74（两端边界），本程也开了，共 **13 枚**读数。
3. 每枚都用 `awk 'NR>=a && NR<=b {print NR"\t"$0}'` 直读**盘上工作树**，逐字节比对被裁文件所引的那句。

### 2.1　12 枚抽样读数（不平均，逐枚报）

| 下标 | 引证 | 被裁文件拿它证什么 | 盘上真内容（摘） | 判 |
|---|---|---|---|---|
| 3 | `.github/workflows/slo-fresh.yml:67-69` | §5.3 S7：第二枚 workflow 里那枚**步骤级** `env:`（`GH_TOKEN`，按纪律只写名不抄值）；同表又指 `:70` 起 `scripts/slo-freshness.sh` | `:67 env:` / `:68 # Read-only: …` / `:69 GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}` / `:70 run: sh scripts/slo-freshness.sh` | **对** |
| 9 | `cmd/wisp/dataroot_128_test.go:76-83` | §2.3：portable.txt 那枚更早早退的测试兜底，并逐字引 `:80` 的报错原文 | `:76 func assertNoPortableOverride128` / `:80 t.Fatalf("premise broke: a portable.txt sits next to the test binary (%s), so resolveDataDir takes the portable branch and never asks the OS for a config dir - these cases would assert nothing", …)` | **对**（引文逐字节） |
| 15 | `cmd/wisp/doctor.go:256` | 全篇的主靶：`if env == "test"` 短路掉 `:259` 那次 OS 读 | `:247 func resolveDataDir(env string)` / `:248-251 portable 分支` / `:256 if env == "test" {` / `:257 return proc.TestDataDir(), nil` / `:259 base, err := userConfigDir() // %APPDATA%` / `:261 拒` / `:263 SealableRoot` / `:264-267 dev 分叉` | **对**（连带 §2.1 那条链的每一格都对） |
| 21 | `cmd/wisp/leg_sink_nail_131_windows_test.go:409` | §3.2 第 5 行：全仓唯一"被测试 in-process 调、且真走 `ResolveEnv()`"的生产入口靠**自己钉**才免疫 | `:408 func TestAC3SecretLegBooksItsAuditRecordsOnDisk(t *testing.T) {` / `:409 t.Setenv("WISP_ENV", "test")` | **对** |
| 27 | `cmd/wisp/resident_windows.go:27` | §2.2 链 B / §3.1 #4：常驻腿走 `ResolveEnv -> proc.Boot` | `:24 func runResident() {` / `:27 env, err := buildinfo.ResolveEnv()` / `:33 rt, err := proc.Boot(env)` | **对**（§4.2 顺手引的 `:24` 也对） |
| 33 | `cmd/wisp/secret.go:289-290` | §1.2：印的是 `c.env` 字段，**不是** `EnvString()` | `:289 return fmt.Sprintf("Active environment: WISP_ENV=%s, data dir=%s, portable=%v\n",` / `:290 c.env, c.dataDir, c.portable)` | **对**（逐字节，含空格） |
| 39 | `cmd/wisp/secret_test.go:200` | §1.3：五处 `strings.Contains` 断言的是字段 ⇒ 与 ambient 无关 | `:200 if !strings.Contains(out, "WISP_ENV=dev") {` | **对** |
| 45 | `cmd/wisp/slo_windows.go:239` | §1.1C：全仓唯一"读来判、然后自己补写"的形状 | `:238 if os.Getenv("WISP_ENV") == "" {` / `:239 _ = os.Setenv("WISP_ENV", "test")` / `:243 env, err := buildinfo.ResolveEnv()` | **对** |
| 51 | `internal/buildinfo/env_test.go:22` | §3.2 第 7 行：相对断言 ⇒ 连 `DefaultEnv` 变了都不红 | `:22 func TestResolveEnv(t *testing.T) {` / `:23/:28/:33` 三形 `t.Setenv` / `:34 want, err := ParseEnv(DefaultEnv)` | **对** |
| 57 | `internal/proc/envfork.go:122` | §5.3 S1：`TestDataDir()` 读的是**另一枚**变量 | `:121 func TestDataDir() string {` / `:122 if dir := os.Getenv(TestDataDirEnv); dir != "" {` / `:125 filepath.Join(SealableRoot(os.TempDir()), "wisp-test-<pid>")`；常量 `envfork.go:36 TestDataDirEnv = "WISP_TEST_DATA_DIR"` | **对**（含常量值另验） |
| 63 | `scripts/build.ps1:107` | §5.3 S6：`(a)` 的下家被 ldflags 钉死成 `dev` | `:107 "-X $BuildInfoPkg.DefaultEnv=$Env",` / `:23 [ValidateSet('dev', 'prod')]` / `:24 [string]$Env = 'dev'`（合法值里**没有** `test` 这一句成立） | **对** |
| 69 | `scripts/slo-check.ps1:111` | §1.4：门禁自己把进程环境设成 `test` ⇒ 对 job env 是覆盖不是依赖 | `:111 $env:WISP_ENV = 'test'` / `:112 $env:WISP_TEST_DATA_DIR = Join-Path $OutDir 'data'`；§5.4 引的 `:114-117` 四行注释也逐字节对 | **对** |

边界两枚一并报：下标 1 `.github/workflows/ci.yml:227`（对，行首 6 空格，`:226 env:`/`:224 test-core`/`:225 runs-on: ubuntu-latest` 三格同对）；
下标 74 `scripts/wisp-cli-tests.sh:109`（对，`export PATH="$dll_dir:$PATH"`，且 §5.3 S4 所指 `:103-105` 两形自陈确在该区间内）。

### 2.2　抽验净数

- **抽样 12 枚 ＋ 边界 2 枚 ＝ 14 枚打开，全部 对。零枚"不存在"，零枚路径捏造。**
- 为写 §1/§4/§5 另开 **11 枚**未抽到的引证，逐枚：
  `doctor.go:37-39`（`info(… critical: false)`）**对**、`doctor.go:49-53`（`fail("gcc (build-time)")`）**对**、
  `doctor.go:300-309`（`probeWritable` 三行引文）**对**、`doctor.go:308`（`return os.Remove(probe)`）**对**、
  `dataroot_128_test.go:257/259/260/261/263/264/266/268`（那枚 helper 全文）**对**、
  `dataroot_128_test.go:281/323`（5 枚腿 ＋ 调用点）**对**、`dataroot_128_test.go:288-292`（四行注释引文）**对**、
  `dataroot_128_test.go:352/370`（`resolveDataDirConsumers128` 与 `if fn.calls["resolveDataDir"] {`）**对**、
  `envfork.go:55/57-65/66-76/77-85/86-88/235/236/238/245`、`boot_windows.go:64/77`、
  `run.go:155-156`、`models.go:109-110`、`secret.go:129-132`、`build.ps1:115-119/162/163` 全 **对**；
  两处否定式现量（`grep -rn 'wisp secret' .github/workflows/ci.yml` -> 0 命中、
  `grep -rn 'data dir writable\|data dir rules'` 在测试/脚本/.workflow 里 0 命中）**对**（与被裁文件 §3.1 #3/#7 同值）。
- **一枚 偏（本程唯一一处行号偏差，写在 §4.2）**：`ci.yml:546-549` 被引作"不设 `MINGW64_ROOT` 就 `the step died before
  compiling anything`"的出处，那句原文其实在 **`:545`**（`:546-549` 讲的是 "build.ps1's documented knob / Verified locally…"）。
  差 1 行，**语义不塌**（同段 `:541-549` 在 §1.5 里另一处引用是对的）。

**本节净数一句话**：**25 枚引证打开、24 对 1 偏（差 1 行、不影响任何判定）、0 枚不存在 ⇒ 这份名册没有编造痕迹。**
（本项不是加分项：它是"这份账能不能被下一位照着复核"的及格线，它过了。）

### 2.3　本节两栏计数增量

真通知回显数：累计 **1**（无新增）。判为注入数：累计 **0**（无新增）。
本节 commit = 见 §6.5 账（写本节时 HEAD 上一枚是 §1 的 `bd0f826`）。

