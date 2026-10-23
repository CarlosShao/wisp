# 票 130 终裁腿 `130-v1` —— `AC#2` 能翻（一条理由要换）｜`AC#1` ⛔（差一次 HEAD 复量）｜产品面 A/B 复核成立

- **被裁**：`130-r1` 五笔 `b19c6612e`／`8f205a6db`／`66e16de00`／`759b10677`／`701448b79`（件 `130-r1-impl.md` 130 行；产码只 1 条语句 `logging.go:149`）。
- **起手锚**：`701448b79`。本件期间共享 HEAD 推进到 `c393740ec`（票 128 收口腿三笔）⇒ `git diff --name-status 701448b79 c393740ec` 不含 `internal/observe`、`logging.go` blob 仍 `7431fcaac` ⇒ 两棵导出树建自 `701448b79`，读数不受影响。
- **写面**：本件＋`D:\tmp\wisp\20261023\130v1\**`；**仓内零 `.go` 字节**、⛔ 未碰他腿 WIP（`cmd/wisp/doctor.go`、`dataroot_128*`、`internal/proc/envfork*`）。
- **尺的声明**：`sha256sum`＝**SHA-256／64 位**（Git Bash，本件自己算）；`git rev-parse` 的 blob id＝**SHA-1／40 位**，两者⛔ 混比。`pre/post` 树＝`git archive` 导出，突变全在仓外。

## ① 九项逐项

### 1) 票面前提新鲜度（对 HEAD blob 复量）
| 票面/件的前提 | 今天 | 尺 |
|---|---|---|
| `FlushEarlyLogRecords` 生产调用者只 `logsink.go:173` 一处 | **仍 1 处**（非测试 `.go` 命中：定义 `logging.go:589`＋注释 3 行＋`logsink.go:173`） | `git grep -n FlushEarlyLogRecords HEAD -- '*.go' \| grep -v _test` |
| `wisp slo` 从不调它（`:272 InitLog`＋`:281 InstallAsDefault`） | **逐字命中**：`slo_windows.go:272`／`:281` | 现读 |
| `InstallAsDefault` 的生产调用者 | **全树 1 处**＝`slo_windows.go:281`（其余是定义/注释） | `git grep -n InstallAsDefault HEAD -- '*.go' \| grep -v _test` |
| `internal/risk/winsec_c26.go:20 init` → `:21 SetPathResolver` | **逐字命中**，`SetPathResolver` 非测试调用者 1 处 | 现读 |
| `logging.go` 还原 sha256 `76ce5f24…b87203` | **逐字相同**；`b19c6612e^` 的 blob＝`13d8180e…db16b12`（件里两枚都对） | `sha256sum`（SHA-256/64） |
| "init 期记录真的掉" | **我自己复现**：pre 二进制 `wisp slo -settle` ⇒ 控制台 2 条／盘上 **0 字节 0 条**（§4） | 仓外 A/B |

⇒ **前提新鲜**；票面标题那把尺今天仍读"丢"。

### 2) `AC#1`（可 grep 的清单）—— **⛔ 按 HEAD 不能翻**；作为 `ac6f31c` 快照实质成立，但缺四件
清单本体＝`docs/evidence/s1/130-ac1-inventory-and-ruling.md`（277 行）。**骨架是真的**：A 组 6 条逐条 `resolve.go:147/151/158/166/176/182` 今天**仍在原位**，每条带消息原文＋"今天去哪了"＋复现命令，不是"可能有"。

| # | 缺什么 | 我的读数（HEAD） |
|---|---|---|
| a | **分母已翻倍** | 主式子 `slog.(Error\|Warn\|Info)(` 今天 **129** 行；件里 **53**。闭包 **27** 包；件里 22，且 `internal/audio`／`internal/ball` 今天**在闭包内**（件 §1.1 的"闭包外"今天为假）⇒ "其余 N 枚都在装之后"这句今天没有分母 |
| b | **总数没核** | 件说"其余 **24** 枚"，逐条只列 **19** 枚（4+4+1+1+1+3+1+3+1）；28＝6+2+1+19 才对得上（尺⑤"自洽≠正确"） |
| c | **16 处行号漂移** | B1 `logging.go:204`→**257**（且该消息今天**两处** 257/379）＝1；四腿装点 57/164/219/276→**109/222/226/339**＝4；`logsink` 144/148/165→**149/153/184**＝3；agent 七枚 373…933→**391/427/650/723/764/804/1016**＋构造点 208→**215**＝8。仍准：A 6/6、B2 `:271`、C1 `:143`、memory 14/14、F1 `:198` |
| d | **漏一个"窗口内"站点** | `internal/winsec/winsec_windows.go:143`（"seal cleared principals"）在 `installLogSink` 的 `InitLog` 内部就发（pipeline 自己封 `logs\` 那次）⇒ 属 **B 窗口**；件里只记在 C 组（别的腿的封印）。我的 post 发里它正是**记录 1**、时间戳早于 install 记录（§4） |

**"每一条能不能造出来"**：A1 **能**（我用 pre 二进制造：控制台在、盘上 0）；A2 **不能**（本机 probe 全过、`probes_passed=2` 我自己读到；件 §6.1 已自认）；A3–A6 生产不可达（`SetPathResolver(nil)`／二次安装非测试调用者 0）＝件里已标"若被造出"。
⇒ **能翻条件**：按 HEAD 把 (a) 复量、把 (b)(c)(d) 改对、A2 标"条件支、无本机读数"。零码，约 1–2h。

### 3) `AC#2`（裁定缓冲策略）—— **能翻**（结论对，**一条理由要换**）
- **(a) 是主流形**：log4j2 `StatusLogger`＝缓冲＋配置完成补发；`slog.SetDefault`／`zap.ReplaceGlobals`／`otel.SetLoggerProvider` 都把"装上"本身当交接时刻 ⇒ 按铁律（最有价值、主流 harness）**这一形成立**。
- **排除三形复核**：① 留在调用者手里＝slo 实测已丢（§4 的 0 条）**成立**；③ 懒补发＝slo 装完一条都不写（post 那 2 条全是补发的，装完之后零记录）**成立**；② 塞进 `InitLog`＝**理由不成立**：我造 M-2 后**票 132 六枚全绿 rc=0**（§6），所以"会花掉票 132 的封印判据"是假话。
- **真代价是另一条**：M-2 会让本笔自己的 `:181`（**带 flush** 的"构造瞬间 0 字节"）红 ⇒ 它花的是本笔新钉，不是 132 的钉。结论仍成立的真理由：**构造 ≠ 装上**（建好但没装的 pipeline 不该把 boot 裁决写进自己的文件）；且 M-2 会让 132 那条 0 字节读数**失去分辨力**（它量的是未 flush 的 `bufio`）——件 §5 自报的跨票仪器问题，我复核**成立**。
- **小缺（登记，不挡翻）**：`InstallAsDefault` 丢弃计数（`_, _`），装上的腿无法 book `early_records`；今天只有 slo 走这条路且它不 book ⇒ 可接受。
- **要改的两处文字**（⛔ 本件不改 `.go`）：`logging.go:486-491` 与 `earlylog_takeover_130r1_test.go:13-20` 的 "spends ticket 132" 句。

### 4) 核心实测：`wisp slo` 真从 0 变 2（我自己造的 A/B）
- 树：`git archive HEAD(701448b79)`→`tree-post`；同一棵只把 `logging.go` 换成 `b19c6612e^` 的 blob→`tree-pre`；`diff -r` **只出这 1 枚**。sha256：post=`76ce5f24…b87203`／pre=`13d8180e…db16b12`。
- 两发 exe 同目录、同三枚 DLL、等长数据根（`leg-slo-E-pre`／`leg-slo-F-fix`、`leg-sec-A-pre`／`leg-sec-B-fix`）。

| 腿 | rc | 控制台 | 盘上 |
|---|---|---|---|
| `slo -settle` **pre** | 0 | 2 | `logs\wisp-20261023-001.jsonl` **0 字节 / 0 条** |
| `slo -settle` **post** | 0 | 2 | **622 字节 / 2 条**，记录 0 是 resolver 裁决 |
| `secret list` **pre** | 0 | 4 | 4 条 / 1311 字节 |
| `secret list` **post** | 0 | 4 | 4 条 / 1312 字节（差 1 是 RFC3339Nano 小数位） |
| `providers` 两发 | 2 | 1 | **无 `logs` 目录** |

- **slo post 记录 0 逐字**：`{"time":"2026-10-23T12:31:20.2650918+08:00","level":"INFO","msg":"winsec: sealing path resolver installed","resolver":"risk.c26Pipeline","probes_passed":2}`；记录 1＝`winsec: seal cleared principals…`、`path=…\leg-slo-F-fix\logs`、**时间戳早于**任何 listener 记录 ⇒ 封目录 WARN 与 resolver 裁决**双双**回来了（件只说 2 条，没说是哪 2 条；这里逐条点名）。
- **硬约束不动（我核）**：`secret list` **4→4 条不变**；booking `"early_records":2,"early_dropped":0` **两发逐字相同**。

### 5) 反向判据（我自己造，⛔ 只看它的名单）
| 变异 | 落地（grep 证） | build | 红（逐名／逐行） |
|---|---|---|---|
| **M-1** 删 `InstallAsDefault` 里那一跳 | `earlyLogs.drain` 只剩定义 | rc=0 | **4/4**：`:194`（file holds 1 record）／`:256`（boot record 0 times）／`:280`（1/0 want 0/0）／`:317`（not in the file at all） |
| **M-2** 挪进 `InitLogWithRegistry`（`return p, nil` 前） | `drain` 在 `:125`＋`:150` | rc=0 | **1/4**：`:181` "after InitLog (and a flush) the boot file holds **158 bytes**, want 0"；另 3 绿 |
| **M-3** 绕过 redaction（`slog.NewJSONHandler(p.writer,…)`） | `drain` 在 `:149` 换目标 | rc=0 | **1/4**：`:321` "the taken-over record's secret is in the file verbatim" |
- **自纠复核（关键）**：M-2 下把 `:177` 的 `p.FlushNow()` 换成 `error(nil)`（＝裸 `os.Stat`）⇒ 头号用例 **PASS rc=0**；带 flush 才红 ⇒ **它那条"带 flush 的强版"做到了**，裸 stat 量的确实是 `bufio`。
- 原样 post `-count=2`：**162/162/0/0**（§7）；每次还原后 `sha256sum`＝`76ce…b87203`、`diff -q` IDENTICAL。

### 6) 与票 132 共存
- `git diff -U0 b19c6612e^ b19c6612e -- internal/observe/logging.go \| grep -E 'PrivateDirAll\|sealLogFile\|ensureFileLocked'` ⇒ **0 行**（rc=1）。
- post 树 132 六枚：**6 PASS / 0 FAIL rc=0**。
- **M-2 下 132 六枚：6 PASS / 0 FAIL rc=0** ⇒ "塞进 `InitLog` 会花掉 132 的判据"**不成立**；连带复核它那句"132 那条 `:51` 在 M-2 下也绿"——**成立**（`:68` 的 `info.Size()` 不 flush，量的是 bufio）。

### 7) 门禁（甲＝`git archive` 仓外导出树＝判据；乙＝工作树＝只作诊断）
| 门 | 甲 | 乙 |
|---|---|---|
| `go build ./...` | **rc=0** | 未跑（他腿 WIP，非判据） |
| `go vet ./internal/observe/` | **rc=0**（win）；`GOOS=linux` **rc=0** | — |
| `go test -count=2 ./internal/observe/` | **RUN=162 / PASS=162 / FAIL=0 / SKIP=0，rc=0** | — |
| `sh scripts/d22scan.sh` | **rc=0 clean**：248/42/116/25/**14**/116/631/148（design=14：导出树只含 tracked 文件） | rc=0 clean：248/42/116/25/**39**/116/631/148 |
| `sh scripts/check-path-length-budget.sh` | **rc=0 VERDICT GREEN**（worst=229、over=61、wall interval=0） | 同 GREEN |
| `gofmt -l internal/ cmd/` | **0 行** | **3 行**＝`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` —— 全是 **`w/crlf` 的 eol 伪红**（CR 字节 131/225/389；索引 `i/lf`）；本笔两枚 `w/lf` |
| `gofumpt -l . tools/d22scan tools/mockllm`（绝对路径） | 20 行、**全在 `.scratch/`**，rc=0；本笔两枚 **0 行** | 54 行＝51 `.scratch/`＋上面 3 枚 eol 伪红，rc=2（`.scratch` 里故意写坏的 `p_unp.go`） |
- `cmd/wisp` 整包四数**未跑**（既有环境上限 ~301–310s 整进程终止、零 `--- FAIL`）⇒ ⛔ 当判据。

### 8) 新增红 0 与那条 flake 的归因
- `internal/observe`：pre／post 各 **77 RUN / 0 FAIL**，名册 77 行**逐字相同** ⇒ 新增红 **0**。
- `cmd/wisp` sink 家族 8 枚：post **8/8 绿 rc=0**；pre **7/8 rc=1**（红 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`，行 `resident_sink_nail_127_windows_test.go:498`）—— **与件"HEAD 8/8、本笔 7/8"方向相反**。
- 单跑三发：post `rc=1,1,0`；pre `rc=0,0,1` ⇒ 两棵树**都红都绿**（红时文件 7 条、绿时 11/18/20 条：子进程没跑到 event loop 就被收）。
- **归因**：这枚跑 `wisp resident`，它走 `installLogSink`（`resident_windows.go:109`），**根本不执行本笔那条语句**（`InstallAsDefault` 全树只有 `slo:281` 调）⇒ 结构上不可能由本笔造 ⇒ **既存环境时序 flake，新增红 0**。但件里"HEAD 只会更红／HEAD 单跑 3/3 红"两个方向性读数**我复现不出**，应删。

### 9) 它交回的三条残余
1. **`logsink` 不走 `InstallAsDefault`**，仍靠 `FlushEarlyLogRecords`：**接受**。两条路互斥有钉（`TestInstallThenFlushDeliversOnce130R1` 我跑绿）；风险是"两条都走"的腿会把 booking 读成 `0/0`（"没有可交的"）——今天无此腿、doc 已写清；落点在 `cmd/**`，非本笔。
2. **无 pipeline 的腿仍只有控制台**：两发都量到 `providers` rc=2／console=1／无 `logs` 目录 ⇒ 属票 131 枚举域，**不并**。
3. **补发先于 `SetDefault` 的微窗口**：读码成立（`drain` 先置 `closed=true` 再 `SetDefault`，此间别的 goroutine 的记录只到旧默认＝tee 的 mirror）。**未造竞态读数**；引导期只有安装者一个 goroutine ⇒ 接受。

## ② `AC#1`／`AC#2` 逐格
| 格 | 判 | 缺什么／条件 |
|---|---|---|
| **`AC#1`** 会掉记录清单 | **⛔ 不能翻（按 HEAD）**；作为 `ac6f31c` 快照**实质成立** | ①分母复量（129 行／27 包，件是 53／22）②"其余 24"改 19 ③16 处行号刷新（B1 204→257 且两处、四腿装点、logsink 三处、agent 七处＋构造点 208→215）④补 B 窗口的 `winsec_windows.go:143`；A2 标"条件支、无本机读数" |
| **`AC#2`** 缓冲策略裁定 | **能翻** | 结论 (a-with-mirror)＋交接点 `InstallAsDefault` 都成立；**必须把"塞进 `InitLog` 会花掉票 132 判据"这句换掉**（M-2 下 132 六枚全绿），换成"构造≠装上＋M-2 会让 132 的 0 字节读数失去分辨力"；两处 `.go` 注释随之改（另派写腿） |

## ③ 提交号／逐笔越界
| 笔 | 号 | 写面 | 越界 |
|---|---|---|---|
| 1 | `b19c6612e` | `docs/evidence/s1/130-r1-impl.md`（A）＋`internal/observe/earlylog_takeover_130r1_test.go`（A）＋`internal/observe/logging.go`（M） | **0** |
| 2–5 | `8f205a6db`／`66e16de00`／`759b10677`／`701448b79` | 各只改 A 件一枚 | **0** |
- 五笔 `.scratch/wisp/**` 零命中（票面未被改）；`logging.go` 相对 `b19c6612e^`＝**55 增 / 2 删**（`--numstat`；`--stat` 的 57 是增删合计），无 `D`、无写面外 `A`。
- **本件**：第一笔 `c3ca063d0` 只落本件；第二笔只落本行（写面仍只本件一行）。索引核：`git diff --cached --name-status` 里另有一笔**他腿已入索引的改名**（票 132 `…-done.md`，`R075`）⇒ 本件两笔都用**显式 pathspec 提交**，未吞、未动它。
