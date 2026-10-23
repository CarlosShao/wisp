# 票 128 AC#5 落地腿 `128-r1` —— 常驻（GUI）腿的自救文案补齐

**本笔量什么**：`internal/proc` 的 `DefaultLayout` 在 OS 说不出用户配置目录时，那一枚**拒绝句**是否"可被人读懂"（票面 AC#5 ①：变量名／应有的落点／怎么设，且走**同一枚文案单点**）。
**锚点**：`96f1af9c9`（现量 `git rev-parse HEAD`，2026-10-23 10:52 +08 开工／12:1x +08 收工）。
**一句话结论**：常驻腿的拒绝句从 **67 字节、三样自救信息 0 枚**，变成 **490 字节、六枚 marker 全在场**；文案的**唯一写点下沉到 `internal/proc`**（`proc.UserConfigDirUnresolved`）。⚠ **`cmd/wisp` 那半（委派 ＋ 钉子升格）按派单"要动先停手上报"未落笔**，patch 已备并**本机真跑过全绿**（§8）。

---

## 1. 前提新鲜度：票面可现跑断言逐条现跑（4 条里 2 条已腐）

| 票面 | 断言（摘要） | 现量（本机现跑） | 判定 |
| :-- | :-- | :-- | :-- |
| `:12` | `resolveDataDir` 里 `base, err := os.UserConfigDir(); if err != nil { base = "." }` | `cmd/wisp/doctor.go:247-268`：`base, err := userConfigDir()` 失败即 `return "", dataDirUnresolved128(env, err)`，**无 `base = "."` 支** | **已腐**（AC#2 于 09-23 落地，`c2fa2e9` 族）⇒ 派单"让 `resolveDataDir` 不回落 CWD"这一半**已在盘上**，本笔不重做 |
| `:9` | `APPDATA` 未设 ⇒ run 腿把 `wisp-dev\logs` 写进 CWD | 真进程 `rc=2`、CWD 空、句子里点名 `APPDATA`（§3） | **已腐** |
| `:10` | 常驻腿同形拒绝，原话＝`proc: user config dir: %AppData% is not defined` | 真进程 `rc=1`、**逐字如此**（67 字节，§2） | 成立（**本笔改的就是它**） |
| `:40` | `internal/proc/envfork.go:238` 逐字 `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)` | 行号 **238**、逐字相同 | 成立 |

⇒ 本笔落 **AC#5**（票面唯一未勾格）。

## 2. 改前读数（真进程，不是"看起来"）

- 仪器（既有）：`go test -count=1 -run TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 -v ./cmd/wisp` —— 真编译 `wisp.exe`、子进程环境**真删 `APPDATA`**、`WISP_ENV=dev`、四枚 CWD 各一枚空目录。
- 源码＝`HEAD` 版：`git show HEAD:internal/proc/envfork.go`，`sha256sum` = `931fcd3e8376554cfeaef08f010baf5e08c4cad7f68eff55140ae4a1faa75d10`。
- **常驻腿逐字**（67 字节含换行；落盘 `resident-before-line.txt`）：

```
wisp: boot failed: proc: user config dir: %AppData% is not defined
```

- 三样自救信息的**在场计数**（`grep -c`，逐枚）：`APPDATA`→**0**、`XDG_CONFIG_HOME`→**0**、`数据根本应是`→**0**、`修法`→**0**、`当前工作目录`→**0**。
- 同一发里 run／secret／doctor 三腿逐字都带三样（同一枚 **471 字节**的句子）⇒ **票面 `:41`"拒了但读不出自救"逐字复现**，且它是六条腿里唯一读不懂的一条。

## 3. 改后读数（同一条命令、同一形：`APPDATA` 缺、`WISP_ENV=dev`、CWD 各一枚空目录）

- 源码＝改后件（`sha256sum` = `a37c4ecc0e867dad1720405feade0fa8fffa4aab080635a4290dca949651bc3a`）。
- **常驻腿逐字**（490 字节；`resident-line-final.txt`，`sha256sum` = `eddd3f3bced0889c32d9d627b4c18714b69ad7da48d95ad9b608a3373a32f52e`）：

```
wisp: boot failed: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev，而 Wisp 拒绝把它回落到当前工作目录（票 128 AC#1 量到回落会搬家：日志、config.toml、DPAPI 私钥存储与 memory.db 跟着启动目录走，换目录再启动就读到空配置）。修法：Windows 把 APPDATA 设为一个可写目录，Linux/macOS 设 XDG_CONFIG_HOME 或 HOME，然后重试。
```

- **六枚 marker 实测全 PRESENT**：`用户配置目录不可得`／`当前工作目录`／`APPDATA`／`XDG_CONFIG_HOME`／`可写目录`／`数据根本应是`（与 `cmd/wisp/dataroot_128_test.go:92` 的 `rescueMarkers128` 同一份六枚）。
- **"没有当前目录回落"的读数**（逐腿真进程）：`run` rc=**2**、`secret list` rc=**2**、`doctor` rc=**1**、常驻 rc=**1**；四腿的 `assertDirEmpty128` 全过 ⇒ **CWD 一个字节都不多**（AC#2 判据"当前目录一字节不许多"仍正面成立）。
- **五条 CLI 腿文本零变化**：`proc.UserConfigDirUnresolved(dev, "%AppData% is not defined")` 与 `cmd/wisp` 现产句**逐字节相同**（471 字节，两者 `sha256sum` 同为 `4410bb0b95e2a4ef179c1c61180ee4b9469caebf3c5a0d6bf0cf10936fe83c55`，`diff` 空）⇒ 委派是**零文本变化**的改写。
- **唯一转红**（若只落 `internal/proc` 那一半）：`cmd/wisp/dataroot_128_windows_test.go:76`（resident 行仍找旧 marker `user config dir`）——不是缺陷，是**既有钉子按旧句写的**；AC#5 ② 要求把它升格成三样 ⇒ 属 `cmd/wisp` 面（§8）。

## 4. 选型：拒绝启动（AC#2 已裁）＋ 文案单点**下沉**到 `internal/proc`

- 选中：`proc.ErrUserConfigDirUnresolved` ＋ `proc.UserConfigDirUnresolved(env, cause)` 成为**唯一写点**；`DefaultLayout` 用它；`cmd/wisp` 侧改成委派（§8 patch）。
- 为什么它最有价值（主流 harness）：`claude code`／`codex cli`／`gh`／`aider` 这类工具在 `XDG_CONFIG_HOME`／`%APPDATA%` 不可得时**响亮失败并给自救句**，⛔ 回落 CWD；且文案由**拥有那次 OS 读取的那一层**统一产出——`internal/proc` 是常驻腿的读取点（`DefaultLayout`）、`cmd/wisp` 是 CLI 腿的读取点，**只有下沉才能一份模板**（`internal/proc` 不可能 import `cmd/wisp`）。
- 不选另两支：回落"有名字的单一点"会造**第二真相源**（AC#2 已按同理由否决）；"维持 CWD 但写进 `doctor`"把"配置随启动目录漂"留在生产里。
- 不是"改一个字面常量"：句子按 env 三叉递增（dev→`wisp-dev`／prod→`wisp`／**test→`WISP_TEST_DATA_DIR` 或 `%TEMP%\wisp-test-<pid>`**，因为 test 的数据根根本不读这个目录，写 fork 会把人指到本进程永不开的树），且要回答"所有调用者拿到它会怎样"（§7）。

## 5. 新用例逐名 ＋ 为什么今天响

| 用例 | 位置 | 钉什么 |
| :-- | :-- | :-- |
| `TestDefaultLayoutRefusesWithNoUserConfigDir`（**135 既存用例上扩写**） | `internal/proc/envfork_test.go` | ①`errors.Is(err, ErrUserConfigDirUnresolved)`（类身份）②六枚 marker **逐名**③三子句**按人读的顺序**且以"然后重试。"收尾（防"六枚 marker 的袋子"）；仍钉 `DataDir == ""`、`os.UserConfigDir` 仍答得出就 `premise broke` |
| `TestUserConfigDirUnresolvedNamesTheEnvFork`（**新增**） | 同上 | dev／prod／test **三叉各点名**、互不串台（`wisp` 是 `wisp-dev` 前缀，按 `<用户配置目录>\<fork>，` 匹配）、OS 原话不被吞 |
| `TestAC5RefusalSentenceHasOneWordingPoint128`（**新增，在 cmd 面 patch 里**） | `cmd/wisp/dataroot_128_test.go` | 本包产句**逐字等于** `proc.UserConfigDirUnresolved` 的产句；且本包**非测试**源码里不得再有 `数据根本应是`／`修法：Windows 把 APPDATA` 的拷贝 ⇒ 这是"不许造第二份模板"的钉子 |

**为什么今天响**：靶子就是那句 67 字节旧句——六枚 marker 一枚不中、`errors.Is` 也非该类（§6 M-1 实证）。

## 6. 反向判据三读数 ＋ 还原核

| 发 | 动作（先证落地、再读红名） | 落地证 | 读数 |
| :-- | :-- | :-- | :-- |
| **M-1** | `return Layout{}, UserConfigDirUnresolved(env, err)` 退回 `fmt.Errorf("proc: user config dir: %w", err)` | `grep` 旧行 1／新行 0；`go build ./internal/proc` **rc=0** | **红 1 枚（逐名）**：`--- FAIL: TestDefaultLayoutRefusesWithNoUserConfigDir`；红行 `envfork_test.go:368`（类身份）、`:391`（六枚 marker 逐名）、`:407`×3（三子句各缺）、`:416`（不以修法结尾）。`TestUserConfigDirUnresolvedNamesTheEnvFork` 仍绿——它量 helper 本身，⚠ 如实登记 |
| **M-2** | 只摘"修法…然后重试。"那一跳（AC#5 ③ 原话） | 同上、`go build` rc=0 | **红 1 枚（逐名）**：同一枚；红行 `:391`（缺 `APPDATA`／`XDG_CONFIG_HOME`／`可写目录`）、`:407`（缺第三子句）、`:416` |
| **还原** | `cp` 回改后件 | `sha256sum`＝`a37c4ecc…bc3a` ＝改后原件 | 复绿：窄射程 **7 枚 PASS／0 FAIL**（`TestLayoutForProdDev`／`TestLayoutForTestEnv`／`TestLayoutForkMatrix`／`TestLayoutForUnknownEnv`／`TestDefaultLayoutUnknownEnvFails`／两枚本笔用例） |

**还原核（先问"谁算的、什么算法、几位"）**：`sha256sum`（Git Bash）＝**SHA-256／64 位十六进制**；`git rev-parse HEAD:internal/proc/envfork.go` 给 **SHA-1／40 位**（`60beb213…`）。两者 ⛔ 互比（尺的射程第 9 例）。

## 7. 调用者影响面（"拿到这个返回值之后会怎样"）

| 调用者 | 拿到新值之后 |
| :-- | :-- |
| `internal/proc/boot_windows.go:120`（`Boot`） | 常驻腿 `wisp: boot failed: <新句>`、`rc=1`，`os.Exit(1)` 在 `installLogSink` **之前** ⇒ CWD 零字节（§3 实测） |
| `cmd/wisp/resident_windows.go:87` | `lerr` 只用于 `ErrAlreadyRunning` 的激活信号、**被忽略** ⇒ 零行为变化 |
| `cmd/wisp` 五条 CLI 腿（`run.go:200`／`models.go:110`／`providers.go:81`／`doctor.go:105`／`secret.go:132`，经 `dataDirUnresolved128`） | 委派后**文本逐字节不变**（§3 的 471 字节同 hash）；`errors.Is(err, errDataDirUnresolved)` 仍真（哨兵别名）；**未委派前**它们拿到的是自己那份拷贝（同文本、第二份模板） |
| `proc.Summarize` | 透传同一错误；**今天无生产调用者**（现量只有 `internal/winsec/…_test.go:185`） |
| 按字符串匹配 `"user config dir"` 的地方 | **只有一处**：`cmd/wisp/dataroot_128_windows_test.go:63` 的 resident marker（实测转红）⇒ 必须同笔升格 |

## 8. 停手上报：AC#5 ① 的"同一枚文案单点"**必须**动 `cmd/wisp`

- AC#5 ① 逐字要求"仍走同一枚 seam 的文案单点（不许再造第二份模板）"。那枚单点今天在 `cmd/wisp/doctor.go:289 dataDirUnresolved128`；`internal/proc` 不能反向 import `cmd/wisp` ⇒ **权威只能下沉**，而"把 `cmd/wisp` 改成委派"是本笔独自无法兑现的一步（派单⛔ `cmd/**`）。
- 我**没有**把 `cmd/wisp` 落笔：只在本机临时验证后**逐字还原**（还原核 `sha256sum`：`doctor.go`＝`a802f6c6…0be`、`dataroot_128_test.go`＝`cbec6e78…b9d`、`dataroot_128_windows_test.go`＝`f8ce6a1a…225`；现量 `git status --porcelain -- cmd internal/proc` **空**）。
- **patch 已备并真跑过**：`D:\tmp\wisp\20261023\128r1\both-halves.patch`（**11 hunks／5 文件／+224 −19**；`proc-half.patch` 6 hunks／2 文件／+153 −4、`cmd-half.patch` 5 hunks／3 文件）。cmd 那半三件：
  1. `cmd/wisp/doctor.go`：`var errDataDirUnresolved = proc.ErrUserConfigDirUnresolved`；`dataDirUnresolved128` 改 `return proc.UserConfigDirUnresolved(buildinfo.Env(env), cause)`；删掉只被它用的 `errors` import。
  2. `cmd/wisp/dataroot_128_windows_test.go:63`：resident 行 markers 改 `rescueMarkers128`（六枚同一份）。
  3. `cmd/wisp/dataroot_128_test.go`：新增 §5 那枚 `TestAC5RefusalSentenceHasOneWordingPoint128`。
- **落地口令（给编排者）**：`git apply /d/tmp/wisp/20261023/128r1/both-halves.patch` ⇒ `go vet ./internal/proc ./cmd/wisp` ⇒ `go test -count=2 -v ./internal/proc` ⇒ `go test -count=2 -v -run '^(TestAC2[A-Z]|TestAC5|TestDefaultLayout|TestUserConfigDir|TestLayout)' ./cmd/wisp` ⇒ 拆两笔 commit（`internal/proc` 一笔、`cmd/wisp` 一笔），撤销口令「撤 128 AC#5」。

## 9. 门禁（rc 一律 `cmd > out 2>&1; rc=$?`，⛔ 走管道；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 尺 | 命令 | rc／读数 |
| :-- | :-- | :-- |
| `go vet` | `go vet ./internal/proc ./cmd/wisp`（两半都在时） | **rc=0**（`vet-both.txt`） |
| 本包测试 | `go test -count=2 -v ./internal/proc` | **rc=0／RUN=100／PASS=98／FAIL=0／SKIP=2**；两枚 SKIP＝`TestHelperProcess`（`jobscope_windows_test.go:87` "helper process mode not set"，既存） |
| 本票族（cmd 面，两半都在时） | `go test -count=2 -v -run '^(TestAC2[A-Z]|TestAC5|TestDefaultLayout|TestUserConfigDir|TestLayout)' ./cmd/wisp` | **rc=0／RUN=72／PASS=72／FAIL=0／SKIP=0**；含四枚真进程腿（`run`／`secret-list`／`doctor`／**`resident`** 全 PASS）与 `TestAC5RefusalSentenceHasOneWordingPoint128` **×2 PASS** |
| 格式甲 | `gofmt -l <改后两枚>` | **输出 0 行／rc=0**（判归属用 `[\\/]`） |
| 格式乙 | `/d/work/base/gopath/bin/gofumpt.exe -l <改后两枚>`（盘上 **v0.12.0**） | **输出 0 行／rc=0**；⚠ 记我一枚：第一版 test 文件被 `gofmt` 判红（`for _, tc := range []struct{ env buildinfo.Env; fork string }{` 单行写法），已展开成多行 |
| 整包 `cmd/wisp` | `go test -count=2 -v ./cmd/wisp` | **本机跑不完（既有、非本笔）**：`HEAD` 原样与改后**两发都在 ~301-310 s 整进程终止**、零 `--- FAIL` 行、无 panic／无 `test timed out`；再加 `-skip 'TestAC14'` 只把终止点挪到 `TestAC13`（`baseline-skip-ac14.sh.log`，310.752 s）⇒ **本机环境的时间上限**，⛔ 与本笔无关 |
| `GOOS=linux go vet` | 未跑（本笔零平台分支；`internal/proc` 的 `envfork.go` 无 build tag） | — |

## 10. 红名册逐名（取主检出）＋ 新增红

- 整包两发**RUN 名册逐名同序**：`HEAD` 原样 225 枚；改后 226 枚（多出的唯一一枚＝cmd 面 patch 里的 `TestAC5RefusalSentenceHasOneWordingPoint128`）。两发的**命名红名册都为空**（`^\s*--- FAIL` 零命中）。
- 本票族（§9 第三行）**72 枚全 PASS、0 FAIL**，`internal/proc` **98 PASS、0 FAIL**。
- ⇒ **新增红 = 0**（同尺、同树、只差本笔那几枚文件）。
- ⚠ **如实登记**：这不是"整包绿"，是**本机整包跑不完**；判据只能是"两发同点终止、名册差集为空"。

## 11. 提交

- **`internal/proc` 那半 ⛔ 未提交**：它与 `cmd/wisp` 那半必须**同笔**才同时满足 AC#5 ① 与"新增红=0"（只落 proc 半 ⇒ `dataroot_128_windows_test.go:76` 一枚新红，§3 已实测）。派单⛔ `cmd/**` ⇒ 按"要动先停手上报"，交 §8 的 patch 由编排者裁。
- 本笔唯一提交：本件 `docs/evidence/s1/128-r1-impl.md`（写面内）；提交号见件尾。
