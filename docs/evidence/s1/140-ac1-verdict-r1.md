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

---

## 3　③　口径核：三口径分不分得开

**要核的三形（AC#1 的问句只关于第二形）**

- **形一**「这条用例在 CI 上**根本没跑**」——票 98／票 111 那本"有没有腿"的账。
- **形二**「**跑了**，但环境是 `test`，所以那一问**没被咨询**」——**AC#1 的靶**。
- **形三**「**跑了也问了**，但断言**恒真**（或期望由被测函数自己算）」——票 119 `R-119-9`／票面 `R-140-1` 那本账。

### 3.1　被裁文件：它**没有** AC#1 的口径轴，它有一根 AC#2 的轴

它唯一成文的判据口径在 §3 开头（逐字）：
「『**会换**』＝该处的**分支选择**随 job env 消失而改变（不管颜色变不变）；『**不换**』＝分支与颜色都不变；
『**换分支不换色**』＝分支变了、可观测结果不变」。

这三档量的是 **"(a) 之后会怎样"**，不是 **"今天它在不在被咨询"**。
两轴在**第四形**上分道：一枚"今天从不被咨询、而 (a) 之后仍然从不被咨询"的腿，
在这根轴上会**稳稳落在"不换"**，读表人把它当成"这条腿有牙"就是误读。
**它自己知道这件事**——§4.3 末行给那句结论时标的档位是
「静态高置信、**可被 `AC#1` 两形真跑证伪**（"换分支不换色"那族的定义就是看起来一直绿）」，
即它把三形的判定**显式推给了 AC#1 那一发真跑**（§6 N1）。**这一点要记在它账上，是记功不是记过。**

### 3.2　被裁表里**口径塌了**的具体行（逐行点名，不平均）

| 行 | 它写的 | 塌成了什么 | 谁答对了 |
|---|---|---|---|
| §3.1 #2（`slo_windows.go:238-239`） | 「**CI 不换**。**但**：直接跑 `wisp slo` 而不经 `slo-check` 时会**换分支不换色**」 | 把 **CI 那一行**与**本机直跑那一行**并排写进同一格 ⇒ 读表人分不清"`wisp slo` 这条生产腿在 CI 上以哪一形被跑"。现量：CI 里**没有步骤直接跑**它
（`grep -c 'wisp slo' .github/workflows/ci.yml` -> **0**），它只作为 `scripts/slo-check.ps1:325/:344` 的子进程被跑，
而那一进程在 `:111` 已自设 `test` ⇒ **`:238` 那枚兜底在 CI 上永不为真**——这就是形二句子
（"那一问从没被咨询"），被裁文件把它写成了"CI 不换"这种 AC#2 句子 | 两枚文件都没把这句写成一句（兄弟 §3.1 的 job 表里有"零枚 `go test`"那一句，够近） |
| §3.2 全 8 行（测试侧"会换分支 = 0 枚"） | 每行的"为什么不动"理由混着三种：「自己 `t.Setenv` 钉了」／「env 是**字面量入参**」／「`cmd.Env` 后出现者胜」 | 这三条**不是同一件事**：第一条才是 (b) 的形状；第二条是**压根不求值环境**（兄弟 §1.0 的"腿型 B／不是问操作系统"）；第三条是子进程自带。把它们统称"免疫"，就答不出 AC#1 那句"**它有没有去问**" | 兄弟 §1.0 用**反向排除**把第二条单独摘出去（`WISP_TEST_DATA_DIR` 是另一枚变量、`t.TempDir()` 那 421 处命中全是 harness 句柄），§2.2 再逐枚标「不是问操作系统——它是'自己声明落点'」 |
| §5.1 第 5-7 行（`internal/buildinfo/env_test.go`） | 「**算，且是全仓唯一一枚"别包"的**……**不依赖 ambient** ⇒ (a)/(b)/(c) 都不动它」 | 通篇没说这枚包**在 Windows runner 上从不编译**（`internal/buildinfo` 只在 `core_pin`，见 §4.3）⇒ 形一缺席登记 | 兄弟 §3.3 明写「它只在 `test-core`(ubuntu) 的 `--scope=core` 里跑……**`test-windows` 一枚都不支撑**」 |
| §1.3 那 10 枚测试文件（表头"全是写者或断言者，没有一处依赖 ambient"） | 逐处给"判什么" | 表里**没有"这一枚在哪个 runner 编译"这一列**；`cmd/wisp` 整包只在 windows 腿跑（`scripts/wisp-cli-tests.sh:60-68` 非 windows 直接 `exit 2`，本程现量），而这 10 枚里带 `//go:build windows` 的与不带的被混成一栏 | 兄弟 §1.2 每行带 tag（`//go:build !windows ⇒ CI 上零分母，见 §3.4`） |
| §4.3 第 1 行（`ci.yml:250` 那枚步骤名） | 「那一步跑的是 `TestLayoutForTestEnv`……**要么换用例要么改名字**」 | 这一格其实是**三形同时**：形一（名字承诺的那件事 CI 从没做过）＋形二（若它真去咨询 ambient 就会被 job env 定死）＋形三（名字与内容对不上 = 恒真式命名）。它只给了"名实不符"一档 | 兄弟 §2.3 最后一行明标「**是本票的账**：见 §2.4」，并在 §2.4 接到"规范义务没人履行" |
| 全表（形三） | 只有两处擦到：§4.3 引「恒红型装饰（比恒绿更贵）」、§6 N7 说 `envfork_test.go:323`「一枚断言两个可能来源……比 `:192` 更像饵」 | **没有任何一行被标成形三**；票面 `R-140-1`（`:192` 的期望由被测函数 `proc.TestDataDir()` 自己算）**不在被裁文件的任何表里** | 兄弟 §2.1 第 3 行逐字记了同一件事：「期望值由**被测函数自己**算出……它只挡'早退被删掉'，不挡'早退退错地方'」 |

### 3.3　一句话裁

- **被裁文件不满足 §3 这一问**：它**没有**"跑没跑／问没问／恒不恒真"这三档口径，
  它有的是 AC#2 那三档，且它在 §4.3 与 §6 N1 里**如实把 AC#1 的判定推给了真跑那一发**。
  => 按 AC#1 的原文（专指第二形）**它不结，且它没自称要结**。
- **兄弟文件满足**：§1.0 定义腿型 A/B（＝形二的判据）、§2.1 的 `两形都绿?` 列（＝形二与形三的交叉）、
  §2.3（形三：一枚断言两个可能来源）、§3.4（形一：文件级零分母）。**三形各有其名、各有凭据。**
- **但静态只能到"候选"为止**：形二的定义是"看起来一直绿"，"一直绿"是**颜色断言**，
  只有票面 `:30` 要的那一发两形真跑能定。=> 这一格不扣两枚文件的分，扣在**"AC#1 整格无法被静态全结"**上。

### 3.4　本节两栏计数增量

真通知回显数：累计 **1**（无新增）。判为注入数：累计 **0**（无新增）。
本节引用被裁文件的句子时全部**按原文引号抄**并给 §号，未把任何兄弟程的结论当已裁读数用
（唯一例外是 §3.2 表末行那句「恒红型装饰」，它是被裁文件**转引** `128-ac4-r1-acceptance.md` §6.2 的，
本程沿用了它的"引别人已裁读数"档位标注，没有升格成本程复现）。

---

## 4　④　CI 侧真凭据（零仪器：只读 `ci.yml` 与三枚脚本）

### 4.1　哪几处**真的**在 job 级设 `WISP_ENV`

`grep -n 'env:' .github/workflows/ci.yml` 现量 **4 处**，逐字节打开：

| job | `env:` 块 | `WISP_ENV: test` 那行 | `runs-on`（现量） |
|---|---|---|---|
| `test-core` | `ci.yml:226` | `ci.yml:227` | `:225 ubuntu-latest` |
| `test-windows` | `ci.yml:336` | `ci.yml:337` | `:335 windows-latest` |
| `slo-smoke` | `ci.yml:481` | `ci.yml:482` | `:480 windows-latest` |
| `slo-full` | `ci.yml:539` | `ci.yml:540` | `:538 [self-hosted, wisp-slo]` |

⇒ **票面 `:9` 那四个行号逐字对**（`:227 :337 :482 :540`），被裁文件 §1.5 的 job 头（`:224/:334/:479/:537`）也逐字对。
**`ci.yml` 里没有任何步骤级 `env:`**（四枚 `env:` 全在 job 键下、缩进 4 空格）；
第二枚 workflow `.github/workflows/slo-fresh.yml` 有一枚**步骤级** `env:`（`:67`，块内只有 `GH_TOKEN`，`:69`），
**它不设 `WISP_ENV`**（`grep -rl 'WISP_ENV' --include='*.yml' .` -> 全仓只命中 `ci.yml` 一枚文件、5 处）。
第 5 处就是 `ci.yml:250` —— **一枚步骤名字符串**，不设任何变量。

### 4.2　哪些步骤**其实不靠它**（自带字面量）

| 步骤 | 它传的字面量 | 为什么 job env 在它身上不被咨询 |
|---|---|---|
| `ci.yml:250`（名字写着 `Environment fork assertion (WISP_ENV=test data dir)`），`run:` 在 `:264` | `bash tools/d22scan/runtests.sh ./internal/proc/ -run TestLayoutForTestEnv`（**`-run` 字面量**） | 那枚用例 `internal/proc/envfork_test.go:75` 在 `:78 t.Setenv(TestDataDirEnv, dir)`（= `WISP_TEST_DATA_DIR`，常量值现量 `internal/proc/envfork.go:36`）**自己声明根**、`:79 LayoutFor(buildinfo.EnvTest, t.TempDir())` **把 env 当字面量入参** ⇒ 全程不求值 `WISP_ENV`。⇒ **步骤名承诺的那件事，这一步从没做**（票面 `:64` 已自陈纠正过这一处；本程独立复现同一读数） |
| `ci.yml:386` `scripts/winsec-tests.sh`；`:422` `scripts/wisp-cli-tests.sh`；`:458`/`:288` `scripts/portable-tests.sh --scope=…`；`:474` `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows` | 包名册 / `-run` 字面量 | 这些包里**没有生产码读这枚变量**：现量 `grep -rl 'wisp/internal/buildinfo' --include='*.go' internal/{secret,config,risk,ball,perm,plugin} cmd/llmrecord` -> **全空**；`internal/proc` 虽 import 它，但 `envfork.go` 里 env 是**入参**，全包 `os.Getenv("WISP_ENV")` **0 处** |
| `scripts/slo-check.ps1:111`（被 `:500`/`:563` 两步调） | `$env:WISP_ENV = 'test'` **字面量** | 它在**同一进程内**先设，`wisp slo` 只是它的子进程（`:325`/`:344` 现量两处 `& $WispExe slo …`）⇒ 覆盖关系，不是依赖关系（与被裁文件 §3.1 #2 同值） |

**反过来，真正依赖 job env 的只有三枚步骤**：`ci.yml:400`、`:497`、`:560`——
三步都跑 `powershell … scripts/build.ps1 -Env dev`，而 `-Env dev` 那枚字面量**只**喂给构建期 ldflags
（`build.ps1:107 "-X $BuildInfoPkg.DefaultEnv=$Env"`，且 `:23 [ValidateSet('dev','prod')]` **合法值里没有 `test`**），
真正吃 job env 的是紧随其后的 `build.ps1:162 & build\wisp.exe doctor`（`:163` 只看 `$LASTEXITCODE`）。
⇒ 这一条两枚文件都给了，本程复算成立；**要补的是下一节那句"它是门禁不是用例"**。

### 4.3　受影响的包**在不在任何步骤的包名册里**（本程现量名册）

名册尺：`scripts/portable-tests.sh` 的四份 pin（行号现量）
`core_pin :124-150`（25 枚包）、`win_pin :151-160`（8 枚）、`cli_pin :161-163`（1 枚）、`winsec_pin :164-166`（1 枚）。

| 包 | 在 core？ | 在 windows？ | 在 cli？ | 在 winsec 门禁？ | 于是它在 Windows runner 上被编译吗 |
|---|---|---|---|---|---|
| `cmd/wisp`（**唯一有 env 短路生产码的包**） | 否 | 否 | **是**（`:422`） | 否 | **是**（`scripts/wisp-cli-tests.sh:60-68` 非 windows 直接 `exit 2` ⇒ 它也**只在** Windows 腿） |
| **`internal/buildinfo`**（**两处真读者的家**） | **是**（`:130`） | **否** | 否 | 否 | **否 —— 任何 Windows job 都不编译它的测试** |
| `internal/proc` | 是 | 是 | 否 | 否 | 是 |
| `internal/secret` | 是 | 是 | 否 | 否 | 是 |
| **`internal/winsec`**（F1 那枚命中所在的包） | **是**（`:149`） | **否** | 否 | **是**（`:386`，但只跑 `./internal/winsec/` 这一枚包在 windows 上的那半边文件） | 包**是**；**带 `//go:build !windows` 的那批问 OS 的文件不是** |

### 4.4　**登记为 AC#1 的加法**（不是更正，按简报要求单独立目）

**A1｜"从不被咨询"比票面说的更强：受影响的那枚**语义包**根本不在 Windows 名册里。**
票面问的是"job 级 `WISP_ENV=test` 让加固腿从不被咨询"。现量：真正**实现**这枚分叉的包
`internal/buildinfo`（`env.go:36`、`buildinfo.go:38` 两处读者，`env_test.go:22 TestResolveEnv`
是全仓**唯一正面断言 `WISP_ENV` 优先级**的用例）**不在 `win_pin`/`cli_pin`/`winsec_pin` 任何一份里**
⇒ 在 `test-windows`/`slo-smoke`/`slo-full` 三枚带 env 的 Windows job 上，
"**`WISP_ENV` 到底怎么被解析**"这件事**连被编译都没有**，谈不上被咨询。
被裁文件 §5.1 答了"它不依赖 ambient"，**没答**"它在 Windows 腿零分母"；
兄弟文件 §3.3 答了后者（"`test-windows` 一枚都不支撑"）⇒ **这一格是加法，落点在兄弟文件那一侧，不改票面主干。**

**A2｜票面那句"加固腿在 runner 上从不被咨询"，其真身是一枚**门禁步骤**不是一条测试。**
`wisp doctor` 在 CI 上只从 `build.ps1:162` 被跑（三步 `:400/:497/:560`），
它**没有颜色可换**——只有 rc=0/rc≠0，且 `:163` 只看退出码。⇒ 那里"那一问没被咨询"的后果不是"用例假绿"，
而是"**这条生产分支在 CI 上从没被走过**"。票面引用的四红一绿指纹来自 `ci.yml:422` 的 `cmd/wisp` 用例，
**那一格今天的腿是有牙的**（`c2fa2e9` 已钉，本程 §2.1 第 9 行复算过那枚 helper 两面自证）。
⇒ AC#1 的名册必须把这两类**分开列**：门禁侧 3 枚步骤（形二在**生产码**上）与用例侧 0 枚（今天已钉）。

**A3｜`test-core` 那枚 job env（`:227`）在本程口径下是纯赘余。**
它管的 job 里只有 `:264`（不读 env，§4.2）与 `:288`（`--scope=core` 的 25 枚包，除 `internal/buildinfo` 外都不 import buildinfo；
而 buildinfo 的用例自己钉三形）。⇒ 摘掉 `:227` 的静态后果是**零**。
这一条与被裁文件 §4.1 那句「`:227` 那枚按 §3.1 是**纯赘余**，可直删不补」**同值**（本程独立复算，不是转述）。

**A4｜`slo-smoke`/`slo-full` 两枚 job 零枚 `go test`**
（现量：`awk 'NR>=479 && NR<=575' .github/workflows/ci.yml | grep -cE 'go test|runtests\.sh|portable-tests'` -> **0**；
该区间内全部 `run:` 只有 4 枚：`:497`/`:500` 与 `:560`/`:563`）。
⇒ 那两枚 job 上的 `WISP_ENV: test` 只服务一件事：`build.ps1:162` 那次 `doctor`。
票面 `:20` 说"副产物才是最该被接住的东西"，**在这两枚 job 上副产物就是唯一产物**。

---

### 4.5　本节净数一句话与两栏计数

**净数**：job 级设值 **4 处**（`ci.yml:226-227/:336-337/:481-482/:539-540`）；真依赖它的步骤 **3 枚**
（`:400/:497/:560`，全部经 `build.ps1:162` 那次 `wisp.exe doctor`）；
传字面量而**不**依赖它的步骤 **6 枚**（`:264/:288/:386/:422/:458/:474`）＋ 自设 env 的 `:500/:563`；
受影响包里在 Windows 名册外的 **2 枚**（`internal/buildinfo` 全包、`internal/winsec` 的 `!windows` 那半边）
⇒ **AC#1 记 A1/A2/A3/A4 四笔加法，票面主干不改。**

真通知回显数：累计 **1**（无新增）。判为注入数：累计 **0**（无新增）。

---

## 5　⑤　`AC#2` 三选一的代价表（**本程不实现、不选支落码**；派单归编排者）

**先报一枚本程量出来的结构事实，它改变整张表的读法：**

> **(a) 支照票面 `:32` 的字面写法（"取消、只有需要它的步骤显式设"）修不了本票自己指的那个缺陷。**
> 现量：CI 上唯一"真需要"这枚 env 的步骤就是 `:400/:497/:560` 那三枚 build 步骤（§4.2），
> 而那三枚步骤需要它的**理由正是**"别让 `wisp doctor` 去碰真实 `%APPDATA%`"。
> ⇒ 照字面做 (a)＝**4 删 3 加**＝运行期行为一字不变＝"那一问在 runner 上仍不被咨询"。
> 真要让 CI 上的生产腿去问 OS，只有 **(a′)＝4 删 0 加**。**(a) 与 (a′) 的代价差一整张表**，
> 被裁文件 §4.1 把两样合写成"净账 4 删 3 加"，没点破这一层。本表把 (a) 拆成 (a)／(a′) 两行。

### 5.1　逐支代价（列＝简报点名的五列）

| | 动哪几枚文件 | 碰不碰冻结路径 | 要改的用例枚数 | 弄错了的失败模式 |
|---|---|---|---|---|
| **(a)** 字面支：4 枚 job 级 env 取消、只给需要的步骤显式设 | **1 枚** `.github/workflows/ci.yml`：删 `:227/:337/:482/:540`，在 `:400/:497/:560` 各加一枚步骤级 `env:`（净 7 行改动）。`scripts/*.sh` **0 处要跟**（现量 `grep -rn WISP_ENV --include='*.sh' .` -> 0）；`slo-check.ps1:111` 已自设 | **不碰任何冻结件**，但**碰共享件门**：票面 `:52` 逐字「改前先登记交回编排者，不许自行改 workflow」 | **0 枚 Go 用例**（§1/§3.2 复算：测试侧会换分支 0 枚） | **静默无效**：改完 CI 全绿、diff 好看、**票面 `:19` 那句"那条读根本不被执行"仍然成立**。这类"改了个不改变行为的东西然后结案"是**假结**，且没有任何颜色会告诉你 |
| **(a′)** 真问 OS 支：4 枚全删、不补 | 同上 **1 枚**，净 **4 行删** | 不碰冻结件；同样**碰共享件门**（而且比 (a) 更该先登记） | **0 枚 Go 用例**，但**必须**先给 `ci.yml:250` 那一步换一枚真读 env 的用例（否则连"看起来有断言"都没了） | **三枚不确定全在这支**（被裁文件 §3.4 U1/U2/U3，本程复算其机制成立）：`build.ps1:162` 的 `wisp.exe doctor` 从"答 `%TEMP%\wisp-test-<pid>`"换成"问 `%APPDATA%`、答 `%APPDATA%\wisp-dev`"，而 `doctor.go:104-105` 是链上**唯一无守卫**的消费者，紧接 `probeWritable`（现量 `doctor.go:300-309`：`MkdirAll` ＋ 写 `doctor-write-probe.tmp` ＋ `:308 os.Remove`）⇒ **真往配置树写**。<br>**红会吃谁**（现量 `ci.yml:257-261`/`:327-328` 两段的自陈：GitHub 停在第一枚红步骤、且全仓无 `continue-on-error`）：`:400` 红 ⇒ 吃 `:422`（`cmd/wisp` 那格读数＝票 123/128/135 共同落点）＋ `:458` ＋ `:474`；`:497` 红 ⇒ 吃 `:500`（SLO smoke）；`:560` 红 ⇒ 吃 `:563`（SLO full，本机最贵的一发，且跑在 `[self-hosted, wisp-slo]`、跨 run 复用同一份工作树，凭据 `slo-check.ps1:114-117` 原文）。<br>**回退清不掉那半**：追加 commit 能把 yaml 装回去，**装不回已经落进 `%APPDATA%\wisp-dev` 的字节**（U3），且 owner 真配置树清理按被裁文件 §4.1 建议应由 owner 做 |
| **(b)** 保 env，每条受影响用例自己钉（推广 `c2fa2e9` 那枚 helper） | **1-2 枚，全在测试面**：`cmd/wisp/dataroot_128_test.go`（扩 `:352 resolveDataDirConsumers128` 的射程——现量它的判据是 `:370 if fn.calls["resolveDataDir"] {`，按构造收不到"直接求值 ambient"的入口）＋ **1 枚新用例**结 SPEC-11:182-183 那笔未履行的规范义务（现量原文：「`test-windows`/`test-core` job 显式断言 `WISP_ENV=test` 生效（数据目录在临时路径、互斥未注册）」）。**不改 `ci.yml`** ⇒ 不触发共享件门 | **不碰冻结件**。**但撞在飞地界**：票 135 篇上 `Status: ready-for-agent`、`Packages:` 行逐字点名 `cmd/wisp/dataroot_128_test.go`、`cmd/wisp/doctor.go`（只读）、`internal/proc/**` ⇒ **同一枚文件两个程**，要先排串行 | **今天该钉而未钉：0 枚**（两枚文件同值，本程复算）；<br>**要改的仪器：1 枚**（名册尺）；<br>**要新增的用例：1 枚**；<br>**可能被新判据误伤的既有站点：3 枚**（现量 `cmdSecret` 的 in-process 调用点 `leg_sink_nail_131_windows_test.go:425/:457/:534`，它们靠 `:409/:526` 的 `t.Setenv("WISP_ENV","test")` 免疫） | **两支假绿**：**(i)** "自己给自己设前提"——被裁文件 §4.2 现量的形状冲突是真的：那枚 helper 钉的是 **dev**（`dataroot_128_test.go:261`），131 那批要的是 **test**（`:409/:526`）⇒ **照抄必红**，而"为了不变红把判据写成只认 `t.Setenv`"就退化成**恒真**（形三）。**(ii)** 钉而不栽：`pinEnvThatAsksTheOS128` 的牙来自 `:260` 先栽 `test` 再 `:261` pin、`:263/:266-268` 两面自证（本程逐字节复核过）；只学"pin"不学"先栽"，换一台机器就变装饰——这正是 `AC#3` 第二问 |
| **(c)** 维持现状 ＋ 把"这条腿测的不是它声称的东西"做成响亮失败的显式登记 | **0 枚码**（登记面）；<br>落点 `docs/reports/pending-and-issues.md` **在本程与被裁文件的禁改名单内** ⇒ 只能由编排者落；<br>要"响亮"就得动 `ci.yml:250`（改步骤名或换用例）⇒ **重新触发共享件门，与 (a) 同纪律** | **不碰冻结件**（`docs/specs/**` 禁改 ⇒ 规范义务不能靠改 spec 消解，被裁文件与兄弟文件同判） | **0 枚用例改**；但**登记条目要能承载 §4.4 的 A1/A2/A3 三笔加法**，否则登记内容是过期的 | **恒红型装饰**（票面 `:44` 逐字禁 `t.Skip`、禁改期望换绿）。更实际的一条：**(c) 现在登记的那句"加固腿变装饰"在 `cmd/wisp` 里今天已不成立**（唯一的靶已被 `c2fa2e9` 钉；本程 §2.1 第 9 行复算过那枚 helper 两面自证）⇒ 照票面原文登记＝**记错主人**。**仍然成立、可以登记的只有两样**：`ci.yml:250` 名实不符，与 SPEC-11:182-183 从未被履行 |

### 5.2　票面前置约束（`AC#2` 那段 09-24 18:4x 追加）在本表里的落点

票 123 那 4 枚（`TestComposedGateBlocksAWriteForTwoSeconds` ＋ 三枚 `TestTicket101*`）：
本程**不复算**被裁文件 §3.3 那条静态链的结论（它已给完整凭据，且 `AC#1` 要的是颜色，颜色只能真跑），
但**本程独立量到它与本票不共因的那半句站得住**：现量 `run.go:155 if s.dataDir == "" {` / `:156 resolveDataDir(buildEnvString())`
⇒ 注入非空 `dataDir` 时那一块整体跳过；而 `grep -rn 'buildinfo\.' cmd/wisp/run.go` 现量**只命中 `:675` 一枚**（`buildEnvString` 的定义）。
⇒ 前置约束 (i)/(ii)/(iii) 三条在本表里分别对应：(a)/(a′) 不许把它们算成本支的成败读数；(b) 的分母不含它们；
**(c) 的登记清单里不许出现它们的名字**——本表 (c) 那一格列的"仍可登记两样"里没有票 123。

### 5.3　推荐（一枚）与它最强的反论

**推荐：(b) 为主，并把 (c) 收窄成两笔具名登记，(a)/(a′) 整体挂起等读数。**
具体形状：① 扩 `resolveDataDirConsumers128` 的射程（把"直接求值 ambient env 的生产入口"纳入名册；
现量这些入口是 `cmdSecret`（`secret.go:184`）、`runResident`（`resident_windows.go:27`）、`cmdSlo`（`slo_windows.go:243`）
三枚，`cmdDoctor`（`doctor.go:43/:104`）与 `printVersions`（`main.go:135`）已在册或不判）；
② 新增**一枚**用例去履行 SPEC-11:182-183（正面断言"来自环境的 `WISP_ENV` 分叉真的生效"），
它同时把 §4.4 的 A1（`internal/buildinfo` 在 Windows 腿零分母）顶掉一半；
③ (c) 只登记 `ci.yml:250` 名实不符 ＋ SPEC-11 那笔义务，不登记"加固腿变装饰"；
④ (a′) 的 U1/U2/U3 由实现程用**一发本机两形 `wisp.exe doctor`** 与**一发 `ls build\portable.txt`** 定掉，再回来重开 (a) 这一格。

**为什么不是别的三支**：(a) 字面支静默无效（5.1 第一行）；(a′) 的失败模式是全表唯一会**写进 owner 真配置树**、
且在 self-hosted 那台上**跨 run 复用工作树**的一支，而它的三处不确定**没有一处能静态定**；
(c) 单开则它的登记命题已经过期一半。

**反论（对推荐最强的那一条，如实写）**：
**(b) 修不了这张票标题里写的那件事。** 票 140 的标题是"job 级 `WISP_ENV: test` 让加固腿在 Windows runner 上从不被咨询"；
(b) 之后，CI 上那三枚 build 步骤的 `wisp doctor` **仍然**答 `test`、**仍然**不问 OS（§4.2），
被咨询的只是**测试自己**。**换句话说 (b) 买的是"仪器不再漏人"，卖的是"生产分支在 CI 上仍不被走"。**
加上它与票 135 撞同一枚文件（`dataroot_128_test.go`），排不上串行就是两个程在同一处名册尺上互相覆盖。
如果编排者认为本票的靶是**后者**（生产腿），那真正该起的支是 (a′)，
而 (a′) 的前置是 U1/U2 有读数——**这两条都不该由本程（裁决者）拍，故本程只把它们摆到台面上。**

---

## 6　⑥　总裁

### 6.1　`AC#1`：**附条件成立**

**一句话总结论**：票 140 的**主干判断在码面上成立、而且比票面写的更强**（§4.4 A1），
但「先量面」这一格**只能被两份静态盘点合起来结掉它的静态那一半**，
票面 `:30` 强制的**两形真跑那一半谁都没做**（本程按简报红线也不能做）⇒ **附条件**。

### 6.2　逐格一句话理由

| 格 | 判 | 一句话理由 |
|---|---|---|
| ① 逐包覆盖核 | **成对结、单枚不结** | 五枚 `WISP_ENV` 命中包里 4 枚答足，`internal/winsec` 在**被裁那枚**的名册里静默缺席（兄弟文件有）＝**F1**；另三处名册算术不闭合（68≠69、3+16+36=55≠62、"其余 12 枚"实为 9）⇒ 被裁文件**自己声明它是 AC#2 的分母**（篇首＋`:10`＋`:33`），拿它结 AC#1 是**指错文件** |
| ② 引证真实性抽验 | **过** | 定距抽 12 枚＋边界 2 枚＋为写其余各节另开 11 组 = **25 枚打开，24 对 / 1 偏（`ci.yml:546-549` 那句原文在 `:545`，差 1 行语义不塌）/ 0 不存在 / 0 枚路径捏造** |
| ③ 口径核 | **被裁那枚不满足；兄弟那枚满足** | 被裁文件**没有**"跑没跑／问没问／恒不恒真"这根轴（它的是 AC#2 的 会换/不换/换分支不换色），且在 §4.3/§6 N1 **如实把 AC#1 推给真跑**；兄弟文件三形各有其名（§1.0 腿型 A/B、§2.1 `两形都绿?`、§2.3 恒真形、§3.4 文件级零分母）。塌口逐行已点名（§3.2 表六行） |
| ④ CI 侧真凭据 | **成立并加固** | job 级设值现量正是 4 处（`ci.yml:226-227/:336-337/:481-482/:539-540`，票面 `:9` 行号逐字对）；真依赖它的只有 3 枚 build 步骤（`:400/:497/:560`，全经 `build.ps1:162`）；6 枚步骤靠字面量、2 枚自设 env ⇒ **A1–A4 四笔 AC#1 加法**，其中最要紧的是**语义包 `internal/buildinfo` 不在任何 Windows 名册里**（⇒"从不被咨询"在那里是"连编译都没有"） |
| ⑤ 代价表 | **可裁，但 (a) 要先拆** | (a) 照字面做＝运行期一字不变的静默无效支；真修只有 (a′)＝4 删 0 加，而它的失败模式是全表唯一会**写进 owner 真配置树**的；(b) 分母最小（0 枚待钉＋1 枚尺＋1 枚新用例）但撞票 135 同一枚文件；(c) 的登记命题在 `cmd/wisp` 已过期一半。**推荐 (b) 为主＋(c) 收窄两笔具名登记＋(a)/(a′) 挂起等读数**，反论已如实写（(b) 修不了票标题那件事） |

### 6.3　派给实现程**之前**必须先量掉的六件事（每件给一发仪器与谁能按）

| # | 要先量什么 | 那一发是什么 | 谁能量 / 前置 |
|---|---|---|---|
| **M1** | **AC#1 本体**：`WISP_ENV` 未设 vs 设为 `test` 两形各一发，给四数＋**名册差集**（哪些**换色**、哪些**换分支不换色**）＋ panic 计数 | `go test -count=2 -v ./cmd/wisp/` 两形（票面 `:30` 逐字要求；`:31` 明令**不许**用 `GOOS=linux go vet` 那类只编译的仪器代答） | 实现程。**前置＝与 `cmd/wisp/**` 上在飞的终裁程错开**（票面 `:77`：它按包 scope 读数，同包并行会造出假颜色）**且等 SLO 取样窗口空出来**（本简报的机器红线即此） |
| **M2** | **U1**：(a′) 之后 `ci.yml:400/:497/:560` 那三步红不红 | 本机一发 `wisp.exe doctor` 的两形对照（覆盖 hosted 那半）；self-hosted 那半**只有 CI 一发覆盖** | 本机＝实现程常规；CI＝**要编排者推**（本程只 commit 不 push） |
| **M3** | **U2**：`build\portable.txt` 在 wisp-slo 那台上存不存在（它比 `:256` 更早 return，见 `doctor.go:248-255`） | 那台机上 `ls build\portable.txt` 一发即定；`grep -n portable scripts/build.ps1` 本程现量 **0 命中** ⇒ doctor 前没有任何东西检查它 | 编排者／owner（登机要授权）；或在 `ci.yml` 加一枚只印不写的诊断步＝**动共享件** |
| **M4** | **A1 的后果**：把 `internal/buildinfo` 纳入 Windows 视野后是什么颜色 | 一发 `runtests.sh ./internal/buildinfo/ -run TestResolveEnv`（与 M1 同窗口，零额外风险） | 实现程。**注意 这是"加法"不是"更正"**，别顺手当 (b) 的分母 |
| **M5** | **`R-140-1`**：`dataroot_128_test.go:192` 的期望由被测函数自己算（`dir == proc.TestDataDir()`） | 一票一发：把早退目标换成另一枚真实目录 ⇒ **必须红**（票面 `:74-75` 给的完成判据） | 本票**不顺手修**（改它＝动别人的对照），但**必须出现在终裁表的可查清单里**（票面 `:74`）。本程现量：兄弟文件 §2.1 第 3 行已逐字记同一件事 ⇒ **不是新账，是催办** |
| **M6** | **F1**：`internal/winsec` 那枚命中补进名册（含它是 `//go:build !windows` ⇒ 在 Windows 腿零分母这一栏） | **零仪器**：一栏表格的补写 | 任何一程都能做，**不该成为 M1 的阻塞** |

### 6.4　派单建议（一句话，不越权）

**140 可以派**，但 scope 要收窄成 **"M1＋M2（＋M3 一发登机）读数 ⇒ 回来重开 `AC#2` 那一格 ⇒ 再落码"**，
且与**票 135 串行**（同文件 `cmd/wisp/dataroot_128_test.go`；票面 `:52` 的共享件门与 `:77` 的错窗约束是同一枚程要一次接的两道前置）。
**本程不替编排者选支、不改任何码、不改 `ci.yml`、不翻票面任何一格。**

### 6.5　本文件的 commit 账（按命令现量，不按本节行文）

```
git log --oneline -- docs/evidence/s1/140-ac1-verdict-r1.md
```

| 节 | commit |
|---|---|
| 篇首＋§0 前提更正＋§1 逐包覆盖核 | `bd0f826` |
| §2 引证抽验 | `e731a7a` |
| §3 口径核 | `7d13450` |
| §4 CI 侧凭据（含 A1–A4 加法） | `6aad697` |
| §5 三选一代价表 | `c4d54c6` |
| §6 总裁（本节） | 写作时 HEAD 上一枚 = `c4d54c6` |

**共树交错的现样本程也撞上了**：本程五枚 commit 之间夹进了兄弟程的 `3f46bf1`
（`docs(136 AC#14 翻勾, A187, 4.0q 追加)`），起手时 `git log -12` 里另有 `a4b5deb`/`92dd40f` 等同树程的落点
⇒ **"按 `git log -6` 数本文件枚数会数错"这一句在本程不是引来的教条，是本程量到的事实**。

**自指余留（照仓里的老规矩如实挂）**：本表前五行的 sha 是在**写本行之前**量出来的，
第 6 行只能写"上一枚"而不能写自己 ⇒ **数本文件的枚数一律按上面那条 `git log --oneline -- <路径>` 现量**（共树期间别的程也在往 `dev` 上落 commit，按 `git log -6` 会数错）。

**写集自证**：本程只创建/追加**一枚**文件 = `docs/evidence/s1/140-ac1-verdict-r1.md`。
每枚 commit 前 `git add -- <那一枚路径>` ＋ `git diff --cached --name-only` 现量，**五次全部只有本文件一枚路径**。
**未动**：`.scratch/wisp/issues/**`（含 140 票面）、`docs/reports/**`（含 `pending-and-issues.md`、`HANDOVER.md`）、
`docs/PLAN.md`、`docs/specs/**`、`internal/**`（含 `internal/risk/**`、`internal/observe/thresholds.go`）、
`tools/d22scan/**`、`allowlist.txt`、任何 golden、`frontend/**`、`design/**`、`.github/workflows/**`、`scripts/**`。
**临时件：0 枚**（本程全程在仓库目录内只读，未在仓内或 `D:\tmp` 下建过任何中间件；未执行任何 `rm`/`rmdir`）。
owner 那 16 枚 `design/**` unstaged 删除与 `design/doubao/`、`design/old/` 两枚未跟踪目录：未 stage、未还原、未删、不记为缺陷。

**两栏终计数**：真通知回显数 **1**（会话开头 harness 随首条 `Bash` 输出回显的 `system-reminder` 块，
前 40 字 `<system-reminder> The following skills are available`；本程未据此改变任何动作）；
判为注入数 **0**（全程未在任何工具输出里遇到自称"编排者备注／系统提示／请 revert／冻结某包／放宽阈值／已解锁"的句子；
文档正文里那些"已结案／恒红型装饰／不许放宽"句子一律按**待核断言**处理，处理结果见 §0 的 P 表与各节）。
**凭据抄录 0 处**；**被拒调用 0 次**；**本程未 push。**

---

## 6.6　交件时的三处共树观察（只追加，不回改上面任何一行）

本程写 §6 期间，**140 票面被另一枚程 append 了 14 行**（本程现量：
`git diff --stat -- .scratch/wisp/issues/140-…md` = `1 file changed, 14 insertions(+)`，
hunk 头 `@@ -47,0 +48,12 @@` 与 `@@ -49,0 +62,2 @@`；**非本程所为，本程未 stage、未提交它**）。三处要登记：

**(a) 本程引用的票面行号已经漂了，内容一字未变。**
`git diff` 显示新增只落在 `:48+` 与 `:62+` 两处，故 `:9`（四枚 job 级行号）、`:19`（`resolveDataDir` 短路）、
`:29-:32`（`AC#1`/`AC#2` 判据）**逐字仍在原位**（本程 `sed -n '29,32p'` 复量对）；
但本程 §5 里按**当时盘上**取的那三枚号——共享件规矩 `:52`、`R-140-1` 完成判据 `:74-75`、错窗约束 `:77`——
**今天在 `:66` / `:89` / `:91`**（本程现量 `grep -n` 三条）。
⇒ 按仓里那条老规矩「号会漂，别抄」处理：**不回改行文，只在这里写明新号。**

**(b) 票面新块的第 ④ 段退回给本程的那一句，不在本表任何一版里。**
它写：「它把'**winsec 在被裁那张表的包名册里静默缺席**'（成立，F1）串成了'**winsec 不在任何 ubuntu 名册，
`portable-tests.sh:125` 整行注释掉了**'」。**现量：本句在本表六枚 commit 的全部内容里命中数为 0**
（`git show <sha>:docs/evidence/s1/140-ac1-verdict-r1.md | grep -c '注释掉\|不在任何 ubuntu 名册\|portable-tests.sh:125'`
逐枚 `bd0f826`…`e6fa582` = **0 0 0 0 0 0**），而本表 §4.3 第 296 行写的恰恰是**相反的那句**——
「`internal/winsec` 在 core **是**（`:149`）」。
⇒ 编排者那句盘上更正（`core_pin` 是 `:124-150`、`:149` 就是 `internal/winsec`、`:125` 是 `cmd/llmrecord`）
**与本表一致，不需要本表改任何东西**；被退回的那句**不是本表的句子**。
这一格按"转述与盘上不符"登记在**票面那一侧**，与本表的 F1 判定无关；**本程不据此改判任何一格。**

**(c) 起手锚两枚读数不一致，本程不复现。**
票面新块记本程「起手锚 `5c6f824`」；本表篇首记的是本程自己 `git rev-parse --short HEAD` 的 **`9090f36`**。
`5c6f824` 不在本程任何一次读数里（本程未见过该号），**登记，不抹任何一侧**——按 (a) 同一条规矩，
锚点是读数不是常量，两枚号都可能只是取数时刻不同。

**(d) 票面新块第 ③ 段把 `AC#1` 判据改写成"三轴"，本表对得上两轴半。**
**(1) ambient 取值来源** ⇒ 本表 §1/§2 与 §4.2 逐处给了（三处真读者、2 枚短路、5 枚字面量入口）；
**(2) 是否只读纯函数** ⇒ 本表 §4.2 那行给了（`LayoutFor` 是纯函数、`DefaultLayout` 才问 OS）；
**(3) 在 `ci.yml` 哪个 step、哪个 runner 有分母** ⇒ 本表 §4.3 给到**包粒度**，
**没给到用例粒度**（那需要票面新块自己说的"把 step 脚本搬进 `D:\tmp` 真跑"那一发，或 `go test -list` 级别的读数）。
⇒ **这一处是本表真实的缺口，不是转述错误**，写在派单必读里：下一发真跑程按**用例**填第三轴。







