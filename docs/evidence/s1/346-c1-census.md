# 346-c1 普查件（票 346 第一趟＝只读；★本腿 0 枚 `go test`、0 个 `.go` 字节、⛔ push）

- 腿＝`346-c1`（只读普查）；起手锚现认＝`76fe906997f7f824ad493ce1e3c5b3cef5c122ab`（`git log -1` 本腿第一笔）；分支 `dev`。
- 起手三核＝HEAD ≥ `76fe90699` ✅／索引干净 ✅（`git diff --cached`＝0 枚、`git diff`＝0 枚、`??`＝1202 枚＝别腿散件与 CI 日志，⛔ 本腿造的）／派单存档 ✅ 已读（`probes/orch/20261102/txt/346c1-dispatch.txt`）。
- 临时件＝`D:\tmp\wisp\20261102\346c1\`（只建不删）。
- ★写面自限：入库件**只有本件一枚**。派单把 `.scratch/wisp/probes/**` 划成**只读**，所以 AC#3 那枚台件**全文抄进本件第三节**（★票面 `AC#3` 要的落点 `probes/346/a1/` 本腿⛔ 能落＝边界冲突，具名交回编排者，⛔ 本腿自行放宽）。⛔ 翻框、⛔ 改票面、⛔ 动码。

---

## 一、AC#0 现跑（票面六行逐枚重跑＋★同族尺；所有读数本腿自取）

### (a) 两枚仪器 blob 同一性＝票面那句成立

尺＝`git rev-parse <rev>:<path>`（三锚 × 五枚路径）：

| 路径 | `a42b35af` | `a5080612` | `HEAD 76fe90699` | 判定 |
|---|---|---|---|---|
| ★`internal/winsec/openleaf_322_other_test.go` | `dfaddac8…` | `dfaddac8…` | `dfaddac8…` | **三锚同串** |
| ★`internal/winsec/posix_seal_window_120_other_test.go` | `0b0e5173…` | `0b0e5173…` | `0b0e5173…` | **三锚同串** |
| `internal/winsec/winsec_other.go` | `e039da10…` | `e039da10…` | `e039da10…` | 三锚同串 |
| `internal/winsec/winsec_windows.go` | `63714a83…` | `63714a83…` | `63714a83…` | 三锚同串 |
| ★`internal/winsec/winsec.go` | `67783ae0…` | `67783ae0…` | **`cbafe712…`** | ★两锚同、**HEAD 变** |

⇒ ★票面第 17 行（两枚仪器＋`winsec_other.go` 两锚同串＝⛔ 产码动过）**独立复现成立**。
★**但要给它补边界**：`winsec.go` 在 `a5080612 → HEAD` 之间**变过**（`git rev-list --count a5080612..HEAD -- internal/winsec/`＝**30 笔**，具名产码笔＝票 342-r1／352-r1／353-r1／348-r2／票 400-r1「winsec 342r1 竞态修」）。⇒ **"字节一字未动"只管这两枚仪器文件、⛔ 管 HEAD 之后的整个被测面**。

### (b) 翻脸事实＝本腿自己 `gh run view` 现取（⛔ 抄票面散文）

环境尺：`gh` v2.96.0 可用、已登录、`origin`＝github.com/CarlosShao/wisp（★⛔ 裸推会去 `cnb`）。两枚 run 的 `headSha` 逐枚对上票面锚：

| run | headSha（逐字对锚） | 两枚仪器顶格结果（★本腿现取） |
|---|---|---|
| **39101060753** | `a42b35afad45b38c…` ✅ | `--- PASS: …322r1 (0.46s)`＋`--- PASS: TestPOSIX120SealWindowLeafSwapIsClosed (0.08s)`；另两枚祖先形状同发 PASS |
| **39124122121** | `a5080612bcf02626…` ✅ | ★`--- FAIL: …322r1 (0.63s)`＋★`--- FAIL: TestPOSIX120SealWindowLeafSwapIsClosed (0.11s)`；两枚祖先形状仍 PASS |
| 39071959151 | `adec1db19f3177bb…` ✅ | ★`openleaf_322_other_test.go` 在该锚**不在树里**（`rev-parse adec1db1:<path>` rc=128）；`posix_seal_window…` 同串 `0b0e5173…` |

⇒ ★票面"blob 一字未动却 PASS⇄FAIL"复现；且**同一 run 内两枚同色**＝**同一件事的两个读数**、⛔ 两枚独立随机事件。
★**★票面那把尺在原始件上⛔ 能照抄**：票面第 18 行 `Select-String '^--- (PASS|FAIL)'`（顶格）对 `gh run view --log-failed` 的落盘件命中 **0 枚**——GitHub 日志每行前缀＝`job⇄step⇄时间戳Z `，payload 不在行首。本腿改成**先剥前缀、再在 `^--- ` 上顶格取**才拿到上表（★红名册与分母都用这形）。★另具名一条射程：`--log-failed` 只回**失败步**的日志，所以"从这份 log 数出的步骤枚数"是**下界**、⛔ 该 job 的完整步骤名册。

### (c) 红句①②③＝逐字复现（同一份 log，★③的形状要更正）

- 红句①（`openleaf_322_other_test.go:206`，逐字）：
  `AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 6 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached`
- 红句②（`posix_seal_window_120_other_test.go:370`，逐字）：
  `AC#1 RED instrument is dead shape=leaf-swap: refused=1998, sealed-own=0 - a zero window reading here proves nothing`
- 红句③（同发同段，逐字）：
  `AC#1 shape=leaf-swap: every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate`

★**更正一处（⛔ 改票面原句、只追加此段）**：③**⛔ 是 `t.Errorf`**——现读码 `:239-241` 那一支是 **`t.Logf`**（随弃权同屏打出、但自己⛔ 让测试红）。★真正让甲枚红的是②（`:231-234` 的 `t.Errorf`）。⇒ AC#2 那句"三句红句逐字仍在"必须按 **两句红句（①②）＋一句同屏读数（③）** 去量；★⛔ 去牙那一发按三枚摘＝**把一条读数当判据拆了**。

★**第二处形状（票面⛔ 写的一维）**：甲枚弃权守卫外面还包着一枚读者档——`:107 sealWindowBaseline120() = os.Getenv("WISP_120_BASELINE")=="1"`、`:113 sealWindowAssertClosed120() = !sealWindowBaseline120()` ⇒ 置 1 时②从 `t.Errorf` 降成 `t.Logf`＝**弃权藏进绿**（`ABSTAINED-HIDDEN`）。★乙枚⛔ 任何降级出口（`:205-209` 直 `t.Errorf`）。⇒ **同一判据下两枚的弃权形状⛔ 对称**。★全仓 `internal/winsec` 里这类 env 档尺＝`git grep -n 'Getenv' -- '*_test.go' '*other*.go'`＝命中 3 处，★只有本枚是"弃权降级"档（另两枚＝`WISP_120_ATTEMPTS` 调发数、`seam_probe_root_125` 读 `TMPDIR`）。

### (d) 分母＋其余在册（两枚 run 各自现量）

| 读数（★顶格尺） | run 39101060753（`a42b35af`） | run 39124122121（`a5080612`） |
|---|---|---|
| 顶格 FAIL 条目（job×名） | **14**＝core 9＋windows 5 | **13**＝core 9＋windows 4 |
| 顶格 FAIL 唯一名 | **12** | **11** |
| 本票两枚 | 0（★全绿） | **2**（★全红） |
| 其余唯一名 | 12 | 9 |

★**票面第 21 行"12 枚唯一名＝test-core 9＋test-windows 3"本腿⛔ 对上**：唯一名 **12** 那发是 **39101060753**（★本票两枚**绿**的那发），它的 windows 顶格 FAIL＝**5 枚**（多 `TestAnchorMissKeepsNonOverridableClassAndNeedL2`、`TestAnchorMissReasonNamesCandidateAndMissingAnchor`、`TestTicket223HandEditedFsLooseningCostsAnL2Card`）；★本票两枚**红**那发（39124122121）唯一名＝**11**、windows＝**4**、其余＝**9**（⛔ 10）。⇒ **票面把两发的名册读成一发了**；两发并集＝**14 个唯一名**、逐名归位见第三节。

★**票面第 22 行"其余 10 枚全部在册"现量＝9 枚**（按 39124122121 口径）：`TestC21TableColourRowsMatchTokensCSS`＋`TestC21DesignTokensFourWayAgree`／`TestApprovalCardViewJSONKeysMatchFrontendTypes`＋`TestComposerContractTypesMatchFrontend`／`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`／`TestStreamLog{DroppedNamingLedger,FloodBelowKeyBound}IsNotBounded35r8`／`TestTicket223PermissionDeniedSitsInItsOwnSentence`＋`TestAC14GoSideEvalPushReachesThePage`。**★"⛔ 一枚归新账"这句本腿复现成立**（未在册＝0 枚）。

### (e) lint 那两步＝票面成立

`gofumpt` 步点名 **20 枚** `.go`、★**逐枚全在 `.scratch/wisp/probes/**`**、两枚 run **20⇄20 逐名相同**（作差空）。⇒ 归票 275／276 存量、⛔ 归本票。（★尺要剥 ANSI 前缀 `^[[36;1m`，否则枚数会少读。）

### (f) ★同族尺：全仓还有几枚"到不了窗口就红"（票面 AC#0 明令的那一枚）

尺＝`git grep -n "RED instrument" -- internal/winsec/` 逐名 `file:line`，★再**逐枚读它外面那层 `if` 的守卫**（⛔ 只认字面词——命中里混着"世界造不出来"那族 `Fatalf`，那⛔ 弃权形）。

★现量＝**26 行命中／11 个文件**。★对比 `346-a1` 那趟的 **12 处／8 个文件**⇒ **同族在两趟之间长了一截**（长出的是票 342r2／348／348r2／353 那几波），★且 `openleaf_guard_322_windows_test.go` **今天已⛔ 在这把尺上命中**（`346-a1` 列过它 `:241`；它被**票 369-r1 重写**、blob `CHANGED`）⇒ **★同族名册会自己腐烂、逐趟必须重跑、⛔ 抄上一趟的表**。

★其中**真"到达数==0 就红"那一支＝12 处**（守卫逐枚读码）：

| `file:line` | 外层守卫 | build tag | 编译进哪一档 |
|---|---|---|---|
| ★`openleaf_322_other_test.go:205` | `liveness == 0 \|\| flips == 0` | `!windows` | test-core |
| ★`posix_seal_window_120_other_test.go:231` | `totals.refused == 0 \|\| totals.sealedOwn == 0` | `!windows` | test-core |
| `posix_seal_window_120_other_test.go:223` | `flipErrs > 0`（racer 翻不动＝另一种"量不到"） | `!windows` | test-core |
| `sealhandle_privateroute_120_other_test.go:211` | `flipErrs > 0` | `!windows` | test-core |
| `ancestoropen_342_unix_test.go:57` | `st.provedLinkCycles == 0 \|\| st.pinnedParentLeg == 0` | `unix` | test-core |
| `ancestorunlink_348_unix_test.go:152` | `st.provedLinkCycles == 0 \|\| st.liveness == 0` | `unix` | test-core |
| `ancestoropen_342_windows_test.go:55` | `st.provedReparse == 0 \|\| st.formRefusals == 0` | `windows` | test-windows |
| `ancestorunlink_348_windows_test.go:110` | `st.provedReparse == 0 \|\| st.formRefusals == 0` | `windows` | test-windows |
| `deleteseam_348r2_windows_test.go:125` | `calls != int32(winAncAttempts342r1)`（发数没跑满） | `windows` | test-windows |
| `deleteseam_348r2_windows_test.go:130` | `landed == 0 && refused == 0` | `windows` | test-windows |
| `deleteseam_348r2_windows_test.go:140` | `landed != calls` | `windows` | test-windows |
| `leafrace_342r2_windows_test.go:388` | `flipped == 0` | `windows && wisp342probe` | ★⛔ 进默认构建（票 369 加的 tag） |

★**另 14 行命中⛔ 弃权形**＝第三种形状（`nametopup_353_windows` 的"进了窗但密封没跑"、`ancestorwindow_342r2_windows`、`ancestorunlink_348_*` 那族 `Fatalf`＝世界造不出来／前置不成立）⇒ ★⛔ 本票射程、具名交回编排者做下一票的料。
★**正控**：同一把尺在同目录命中 26 行（⛔ 0）＝尺⛔ 失明。

★**平台射程（票面 AC#0 必答、也是 346-a1 那句的独立复现）**：两枚票面仪器首行逐字 `//go:build !windows` ⇒ **只在 Linux 档编译**；★本机 Windows 那一轴对它们**结构上不存在**＝既⛔ 绿也⛔ 弃权。⇒ 本腿"两平台成对尺"对这枚只能是 **CI `test-core` ⇄ 票面两枚 run 的原始 log**；★⛔ 拿"本机没红"读成绿。

---

## 二、AC#1 三形代价表（★⛔ 选形；每形现量：动几枚文件／顶不顶既有钉／CI 那一步分母变不变）

★**先把三形共同的台面事实量出来**（⛔ 任一形都⛔ 谈成本）：跑这两枚仪器的 CI 步骤＝`test-core` 的 `Portable package tests`＝`run: bash scripts/portable-tests.sh --scope=core`（`ci.yml:676`）；该脚本把 core 档全部包**塞进同一枚 `go test -v -count=1`**（`portable-tests.sh:901` → `tools/d22scan/runtests.sh:211`，★尺＝`grep -n 'go test -v' runtests.sh` 那一行**⛔ 带 `-p`**）。★现量＝两发 `test-core` 各 **27 枚包结果行**、红那发的包行逐字 `FAIL	github.com/CarlosShao/wisp/internal/winsec	1.434s`（log L10882）。★job 墙钟（`gh run view --json jobs`）＝`test-core` **117s ⇄ 115s**、`test-windows` **344s ⇄ 740s**。⇒ **"CI 争用"⛔ 是修辞：这枚竞速仪器与另 26 枚包共享一次 `go test` 调度**。

### 甲＝台面侧安静化

| 必答项 | 本腿现量 |
|---|---|
| **落点（★最便宜的形状已经在仓里、CI⛔ 叫它）** | `portable-tests.sh:398-405` 有一枚建好的 **`--scope=winsec` 档**（票 251 AC#1 留的、`winsec_scope`/`winsec_pin` 逐字在位），★但 `ci.yml` 从来⛔ 叫它——`selftest` 自己逐字写着「`ci.yml:439` never says `--scope=winsec` - it runs `scripts/winsec-tests.sh`，后者递的是**显式包路径**」（`:36-38`）。★而 `winsec-tests.sh` 那一步挂在 **`test-windows`**（`ci.yml:838`）＝**Linux 档⛔ 走过这枚安静入口**。⇒ 甲＝**叫那枚已经存在的档**、⛔ 新造机制。 |
| **动几枚文件** | **1 枚**＝`.github/workflows/ci.yml`（加一步／或把 winsec 从 core 档挪出）。★仪器与门禁脚本**⛔ 一枚⛔ 动**＝甲的全部卖点。 |
| **★代价的落点（只写清、⛔ 拿它当降级理由）** | ①**`.github/**`＝人工批准**（票 346 `AC#5` 逐字禁、`AC#3` 同判）；②**"从 core 档挪走"那形顶红的是别票的钉**：`core_pin` 逐字列着 `github.com/CarlosShao/wisp/internal/winsec`（`portable-tests.sh:322`）＋GUARD C/D（`:688`/`:709`/`:958`）＋票 254 AC#2 钉的"core 与 winsec 两档**共用同一枚 `$winsec_scope` 变量**"（`:353-361`）＋票 251 AC#1 的 `winsec-tier-*` 正反两控（`selftest:474-491`）⇒ ★**挪＝顶票 111／251／254 的钉**；③**"另加一步"那形＝winsec 被跑两遍**（core 那遍照旧红）⇒ ★**分母⛔ 变、时长变**；④**"什么算够安静"这条判据要写死**：现量给料＝仪器自己 `0.08s`/`0.63s`、`internal/winsec` 整包 `1.434s`、`test-core` job `115-117s` ⇒ 安静档的边际成本是**整包重跑量级（≈＋1.4s／＋1 步）**、⛔ 秒级以下的雕花。 |
| **CI 那一步的分母** | 加一步＝**分母不变、多一步结论**；挪出 core＝**★顶红既有钉**（见②）。 |
| **顶不顶票 322・120 的既有钉** | ★⛔ 动断言、⛔ 动红句 ⇒ **⛔ 顶 322/120**；★但**顶 111／251／254 的门禁钉**（挪那形）。 |
| **★甲买到的到底是哪半条** | 甲**⛔ 消除弃权**：弃权那一支是票 322 `AC#1` 逐字要求的（`⛔ 拿"0 命中"当凭据，除非同一发证明窗口真到达过链接态`）⇒ 安静只**降低概率**。★现量给了这个概率的量：红那发到达数 `flips=6 / sealed-own=0`，绿那发 `flips=504 / sealed-own=159`＝**两个数量级**。⇒ ★甲买的是**"能不能常常够得到窗口"**、⛔ "永远够得到"。 |

### 乙＝仪器侧确定性交错

| 必答项 | 本腿现量 |
|---|---|
| **落点** | 仪器自己的 racer 与被测 leg 的接缝。★现读＝**两枚今天⛔ 吃任何接缝**（尺＝`git grep -n applyDescriptor -- openleaf_322_other_test.go posix_seal_window_120_other_test.go`＝**0 命中**），它们靠真 `rename` 线程硬抢。 |
| **动几枚文件** | ★**本仓已有"纯 `_test.go`、0 产码字节"的先例**：票 369-r1（`7d9e5b481`）对**Windows 那两枚同族仪器**做的正是这一形，其 `AC#2` 判语逐字＝**「5 改 1 新、`产码 0 字节`、全 `_test.go`」**（`openleaf_guard_322_windows` `168+/91-`／`ancestorunlink_348_windows` `88+/17-`／`leafrace_342r2_windows` `26+/1-`＋新件 `routeguard_369_test.go`），连跑 10 发＝**`RUN=40 PASS=40 FAIL=0 SKIP=0`**。⇒ ★"乙一定要动产码"这个前提**在本仓⛔ 成立过一次**。 |
| **★票面对乙的前提，本腿现读要收窄（⛔ 本腿改原句、只具名此段）** | 票面第 34 行写「代价＝要动 `internal/winsec` 产码的注入点＝⛔ 本票射程、须另票＋人工批准」。★本腿现量：包自己声明的可替换接缝里有 **`winsec_other.go:50 var applyDescriptor = applyDescriptorPOSIX`**（★`!windows` 侧、★已被 `sealhandle_window_120_other_test.go:170`⇄`private_fail_test.go:17` **在测试里换过**）＝**接缝在 POSIX 侧真存在、真用得上、⛔ 需要新产码**；★且 `winsec_other.go` 三锚逐字同串（`e039da10…`）＝那枚接缝从头到尾⛔ 变过。⇒ **"须另票＋人工批准"⛔ 成立于"必须改产码"这个理由**；★它成立与否要按 `AC#5` 的字面（禁的是那三枚**字节**）重裁：**消费接缝＝写在 `_test.go`、⛔ 动那三枚字节**。★这一条**归编排者裁**。 |
| **★乙的真代价落在哪（本腿现读、★⛔ 好话）** | ⛔ 在"要不要新产码"，在**接缝覆盖⛔ 全**：`applyDescriptor` 是**密封那一步**的接缝（`:288`/`:317`/`:432` 都走它）；而★**open/truncate 那一支⛔ 任何可替换接缝**——`privateFile` 的 POSIX 序列是 `pinPlacement`（`:149`）→ `cur.OpenFile(leaf, O_RDONLY|O_NONBLOCK)`（`:189`）→ `os.SameFile`（`:222`），★逐枚都⛔ 是 var；★包内全部 var 接缝尺＝`grep -n '^var [a-zA-Z_]* = [a-zA-Z_]*$' internal/winsec/*.go`（⛔ `_test.go`）＝**只有 4 枚**：`winsec.go:128 removeRefusedLeaf`（★删除支）、`winsec_other.go:50 applyDescriptor`（★密封支）、`winsec_windows.go:38/:135/:768`。⇒ ★**甲乙两枚要的⛔ 同一枚接缝**：120 那枚（密封窗口）**有现成的可用**；322 那枚（open 窗口）**按盘上现量需要新接缝、或抄票 369 那条⛔ 加接缝的写法**（同步 rename 那一发＋`reparseAt` 证世界进过翻态）。★票 369 选的恰是后者 ⇒ **POSIX 侧有现成可抄的路、且那条路⛔ 动产码字节**。 |
| **CI 那一步的分母** | ⛔ 变（同步、同步骤、同名册）。 |
| **顶不顶票 322・120 的既有钉** | ★**会顶，逐名读**：票 322 `AC#1` 钉「成对读数同发同屏：修前分子 > 0，修后 **0 且 liveness ⛔ 零**」；票 120 `AC#1` 钉「多样本**全报**、⛔ 只报最好的那次、⛔ 用平均抹掉坏尾部」。⇒ 确定性化**⛔ 取消多样本、⛔ 让 `damaged` 变非零**；★但仪器自己的**计数语义**是票 120 的分母：`sealWindowSamples120 = 4`（`:100`）、`sealWindowAttempts120 = 500`（`:101`）、`leafRaceAttempts322r1 = 3000`（`:35`）——★**动这三枚常量＝顶钉**（`4×500 多样本全报`就是那个"4"）。★`WISP_120_ATTEMPTS`（`:116`）那枚 env⛔ 动判据、⛔ 算多样本全报＝它是把分母交给台面，AC#2 禁的正是这一形。 |
| **★本腿⛔ 能答的** | "乙要改哪几行 `_test.go`"＝实现面判断。本腿只给落点＋接缝清单＋先例。 |

### 丙＝判据侧同屏打出可归因读数

| 必答项 | 本腿现量 |
|---|---|
| **★票面对丙的前提，本腿现读＝已成立一大半** | 票面第 34 行要丙「每次运行落 `attempts／flips／refused-at-check／refused-at-act` 那四个数」。★现量＝**两枚今天已经同屏全打**：乙枚 `openleaf_322_other_test.go:192` 逐字 `AC#1 paired reading on one screen: attempts=3000 flips=6 flipErrs=0 damagedWitnesses=0 leafLegRefusals=0 foreignsNarrowed=0 foreignsHoldingPayload=0`＋`:196` 逐腿"哪条腿答的"名册；甲枚 `:365`（逐样本）＋`:370`（合计）逐字带 `attempts / refused-at-check / of-which-refused-at-act / sealed-own-object / other`。⇒ ★**"打数"那一半⛔ 新账**；★丙剩下的只有一件：**把"这一发＝弃权"做成名册里可机读的一行**（今天只有散文红句）。 |
| **动几枚文件** | **0 枚**（数已在）；要把弃权行做成稳定可 grep 的标记＝**2 枚 `_test.go`**（各加一行 `t.Logf`/`t.Errorf` 前缀、★⛔ 动判据⛔ 动那两句红句）。 |
| **CI 那一步的分母** | ⛔ 变。 |
| **顶不顶票 322・120 的既有钉** | ⛔ 顶（★判据零移动＝丙唯一的卖点）。 |
| **★丙买到的／⛔ 买到的** | ★票面自己写了"CI 仍会随机红"⇒ 丙**⛔ 让门别红**（⛔ 与 AC#2 的"⛔ 用变绿换看不见"同判）。它买的是**归因**。⇒ ★**丙与第三节 AC#3 那枚台件是同一条价值的两种付法**：一个落在仪器输出、一个落在仓外台件。★本腿现量：AC#3 那把尺**⛔ 动任何字节**就把两发 run 的 **14 个唯一名全部对回在册票、未在册 0 枚** ⇒ **丙的"可归因"那一半今天已经拿到了**；★丙真正还剩的价值＝**让红名册⛔ 依赖一份仓外 log 就能自证**（台件要拿到 log、仪器内的标记不要）。 |

### 三形并排（★⛔ 选形；这张表交非实现者裁）

| 形 | 动几枚文件 | 落在哪一面 | CI 分母 | 顶哪枚既有钉 | 买到什么 | 本仓先例 |
|---|---|---|---|---|---|---|
| **甲** 台面安静化 | **1**（`.github/workflows/ci.yml`） | 台面 | 加一步＝⛔ 变／挪出 core＝★变 | ★**票 111 GUARD C/D＋票 251 winsec 档＋票 254 共用变量**（挪那形）；⛔ 顶 322/120 | 到达率的两个数量级；★⛔ 消除弃权 | ★`--scope=winsec` **已建好、CI⛔ 叫**（＋`winsec-tests.sh` 在 windows 档已是"单独一步"的先例） |
| **乙** 确定性交错 | **2 枚 `_test.go`**（＋0～1 枚接缝） | 仪器 | ⛔ 变 | ★`sealWindowSamples120/Attempts120`、`leafRaceAttempts322r1` 三枚常量＝顶票 120/322 的分母 | ★**唯一能消除"靠调度运气"那一形的** | ★**票 369-r1**＝同形、Windows 侧、`产码 0 字节`、`40/40×10` |
| **丙** 同屏打数 | **0～2** | 判据输出 | ⛔ 变 | ⛔ 顶（零判据移动） | 归因；★**CI 照旧随机红** | ★两枚仪器**已在打数**（AC#3 尺＝今天已拿到归因那一半） |

---

## 三、AC#3 归因工具（一枚可重跑台件；全文在此，★落点冲突见下）

★它结的是票面那句"编排者侧欠账"＝把"CI 红名册逐名对回在册票"从人肉变成尺。★⛔ 新建 CI 步骤、⛔ 做成常驻断言。
**★落点冲突（具名交回）**：票面 `AC#3` 要它放 `.scratch/wisp/probes/346/a1/`；★本腿派单把 `probes/**` 划成**只读**、且"只 commit 证据件"⇒ ★本腿**⛔ 落盘**、把**全文**抄在这里，等一枚有写面的腿照抄进 `probes/346/a1/`。★⛔ 自行放宽边界。

```bash
#!/usr/bin/env bash
# ci-red-attribution.sh — ticket 346 AC#3: turn a CI red roster into a per-name
# "which ticket owns this" reading, so no leg has to guess whether a red is new.
#
# INPUT   one or more GitHub Actions run ids, OR pre-fetched logs via --log <path>.
# OUTPUT  per unique 顶格 FAIL name: jobs / 判据在册票 / 具名文件在册票 / 散文提及 /
#         台账 / REGISTERED|UNREGISTERED. UNREGISTERED (=0 判据 +0 文件 +0 台账)
#         is the ONLY row a new ticket may be opened for.
#
# Read-only. No CI step (.github/** = human approval), no standing assertion.
# The positive control runs FIRST: if the ruler cannot attribute a name we KNOW is
# registered, it exits 1 instead of reporting a clean roster it never measured.
#
# Usage:
#   bash ci-red-attribution.sh 39101060753 39124122121
#   bash ci-red-attribution.sh --log a.log --log b.log
set -uo pipefail
ISSUES_DIR="${ISSUES_DIR:-.scratch/wisp/issues}"
LEDGER="${LEDGER:-docs/reports/pending-and-issues.md}"
POSCTL="${POSCTL:-TestC26RewriteAccountIsConsumedAtEverySecurityLeg}"
tmpd=$(mktemp -d 2>/dev/null || echo "./.ci-red-attrib.$$"); mkdir -p "$tmpd"
trap 'rm -rf "$tmpd"' EXIT
logs=(); runids=()
while [ $# -gt 0 ]; do case $1 in
  --log) shift; logs+=("${1:-}"); [ $# -gt 0 ] && shift ;;
  -h|--help) sed -n '2,18p' "$0"; exit 0 ;;
  *) runids+=("$1"); shift ;;
esac; done
for r in ${runids[@]+"${runids[@]}"}; do
  dst="$tmpd/run-$r.log"
  if ! gh run view "$r" --log-failed >"$dst" 2>"$tmpd/err-$r"; then
    echo "ci-red-attrib: gh run view $r FAILED -> take the run id off the ticket face," >&2
    echo "  read the roster there and MARK IT 未重取 (a reading, not a green)." >&2
    sed -n '1,3p' "$tmpd/err-$r" >&2; exit 3
  fi
  logs+=("$dst")
done
[ ${#logs[@]} -gt 0 ] || { echo "ci-red-attrib: no run id and no --log" >&2; exit 2; }
# stage 2: strip the GH 'job⇄step⇄timestamp' prefix, THEN anchor ^--- so go-test
# indented SUBTESTS stay out of the denominator (this is where 票面's '^--- FAIL' died).
for L in "${logs[@]}"; do
  [ -f "$L" ] || { echo "ci-red-attrib: unreadable log $L" >&2; continue; }
  awk -F '\t' 'NF>=3 { payload=$3; sub(/^[0-9-]+T[0-9:.]+Z /, "", payload);
    if (payload ~ /^--- (PASS|FAIL|SKIP): [A-Za-z0-9_.\/-]+/) {
      verdict=payload; sub(/^--- /,"",verdict); sub(/:.*/,"",verdict);
      name=payload;    sub(/^--- [A-Z]+: /,"",name); sub(/ .*/,"",name);
      print $1 "\t" verdict "\t" name } }' "$L"
done >"$tmpd/all.tsv"
awk -F '\t' '$2=="FAIL"{print $1"\t"$3}' "$tmpd/all.tsv" | sort -u >"$tmpd/roster.tsv"
echo "ci-red-attrib: logs=${#logs[@]}  顶格 FAIL (job x name)=$(wc -l <"$tmpd/roster.tsv"|tr -d ' ')"
# stage 3: positive control (the ruler must prove it is not blind).
if [ -s "$tmpd/roster.tsv" ] && grep -q "	$POSCTL$" "$tmpd/roster.tsv"; then
  if grep -rl -- "$POSCTL" "$ISSUES_DIR" >/dev/null 2>&1; then
    echo "ci-red-attrib: POSITIVE CONTROL ok ($POSCTL resolves to a ticket file)"
  else echo "ci-red-attrib: POSITIVE CONTROL FAILED - pool unreadable or name absent." >&2; exit 1; fi
else
  if grep -rl -- "$POSCTL" "$ISSUES_DIR" >/dev/null 2>&1; then
    echo "ci-red-attrib: positive control ok standalone (name absent from this roster; 尺 not blind)"
  else echo "ci-red-attrib: positive control FAILED standalone; exit 1" >&2; exit 1; fi
fi
# stage 4: ONE criteria index (a per-name grep-per-ticket-file loop cost 11k+ process
# spawns and never finished), then per-name lookups against it.
grep -HE '^- \[[ x]\]' "$ISSUES_DIR"/*-*.md >"$tmpd/crit.tsv" 2>/dev/null || : >"$tmpd/crit.tsv"
cut -f2 "$tmpd/roster.tsv" | sort -u >"$tmpd/names.txt"
echo "---- per-name attribution ----"; unreg=0
while IFS= read -r name; do
  [ -n "$name" ] || continue
  jobs=$(awk -F '\t' -v n="$name" '$2==n{printf "%s,", $1}' "$tmpd/roster.tsv" | sed 's/,$//')
  esc=$(printf '%s' "$name" | sed 's/[][\\^$.*+?(){}|]/\\&/g')
  own=$(awk -v re="$esc" '$0 ~ re' "$tmpd/crit.tsv" 2>/dev/null \
        | grep -oE '/[0-9]{2,4}-[^/]*\.md' | sed -n 's#.*/\([0-9]\{2,4\}\)-.*#\1#p' \
        | sort -un | tr '\n' ',' | sed 's/,$//')
  ment=$(grep -rl -- "$name" "$ISSUES_DIR" 2>/dev/null | wc -l | tr -d ' ')
  # FILE-BASED fallback: 票 346 names its instruments by FILE, never by Go func name,
  # so a name-only ruler reports this ticket's own subject as fresh debt (measured).
  file=$(git grep -l -E "^func ${name}(\(| )" -- '*_test.go' 2>/dev/null | head -1)
  base=""; own_file=""
  if [ -n "$file" ]; then base=$(basename "$file")
    own_file=$(grep -rl -- "$base" "$ISSUES_DIR" 2>/dev/null \
      | sed -n 's#.*/\([0-9]\{2,4\}\)-[^/]*\.md#\1#p' | sort -un | tr '\n' ',' | sed 's/,$//')
  fi
  led=$(grep -h -- "$name" "$LEDGER" 2>/dev/null | sed -n 's/.*\([ARQ][0-9]\{2,\}\).*/\1/p' \
        | sort -u | tr '\n' ',' | sed 's/,$//')
  if [ -n "$own" ] || [ -n "$own_file" ] || [ -n "$led" ]; then st=REGISTERED
  else st=UNREGISTERED; unreg=$((unreg+1)); fi
  printf '    name=%s\n      jobs=%s  file=%s\n      判据在册票=%s  具名文件在册票=%s  散文提及=%s枚票  台账=%s\n      => %s\n' \
    "$name" "${jobs:-?}" "${base:-未解析}" "${own:-—}" "${own_file:-—}" "${ment:-0}" "${led:-—}" "$st"
done <"$tmpd/names.txt"
echo "---- summary ----"
printf 'unique red names=%s  UNREGISTERED=%s\n' "$(wc -l <"$tmpd/names.txt"|tr -d ' ')" "$unreg"
[ "$unreg" -eq 0 ] || printf 'NOTE: only UNREGISTERED rows may open a new ticket.\n'
```

### ★两枚 run 现跑它（★输出逐字是本腿量出来的）

起手读数行＝`ci-red-attrib: logs=2  顶格 FAIL (job x name)=16`（14 个唯一名、两枚重叠名各算两条）、`POSITIVE CONTROL ok`、`unique red names=14  UNREGISTERED=0`。

★本票那两枚的逐字输出：
```
name=TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1
      jobs=test-core  file=openleaf_322_other_test.go
      判据在册票=—  具名文件在册票=322,346,369  散文提及=2枚票  台账=—
      => REGISTERED
name=TestPOSIX120SealWindowLeafSwapIsClosed
      jobs=test-core  file=posix_seal_window_120_other_test.go
      判据在册票=—  具名文件在册票=120,322,342,346,369  散文提及=2枚票  台账=—
      => REGISTERED
```

其余 12 个唯一名（★两发并集，逐名归位）：

| 红名 | jobs | 具名文件在册票 | 台账 | 归位 |
|---|---|---|---|---|
| `TestAnchorMissKeepsNonOverridableClassAndNeedL2` | core | — | `A380` | 票 339／336／322 族（★只在 39101060753 红） |
| `TestAnchorMissReasonNamesCandidateAndMissingAnchor` | core | — | `A380` | 票 339 族（★只在 39101060753 红） |
| `TestTicket223HandEditedFsLooseningCostsAnL2Card` | win | 111,223,227,231…351 | A510,A512,A515,A518 | 票 303／★票 351 那一族（只在 39101060753 红） |
| `TestC21TableColourRowsMatchTokensCSS` | core,win | — | A1037,A380,… | 票 344（★两档都红＝同一条 CSS/表对撞） |
| `TestC21DesignTokensFourWayAgree` | core | — | A208,A237,… | 票 344 |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | core | — | A380,A426,… | `frontend/**` 缺席＝归机主侧 |
| `TestComposerContractTypesMatchFrontend` | core | — | A376,A383,… | 同上 |
| `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | core,win | 93,102,108,116,181 | `A380` | 存量（★＝正控那枚） |
| `TestStreamLogDroppedNamingLedgerIsNotBounded35r8` | core | — | `A716` | 票 35r8 在册 |
| `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8` | core | — | `A716` | 同上 |
| `TestTicket223PermissionDeniedSitsInItsOwnSentence` | win | 231 | `A875` | 票 303／305／314 那一族 |
| `TestAC14GoSideEvalPushReachesThePage` | win | — | A508,A607,A779,A875,A1079,A1100 | 同上 |

### ★这把尺自己的两条射程（写进件、⛔ 假装它是全称）

1. ★**"判据在册票"那一列今天对全部 14 枚都回 `—`**——因为票面把红名写在**散文与现量表**里、⛔ 写进 `^- [ ] AC#` 那一行⇒ **那格今天⛔ 灵**；今天真正承重的是**"具名文件在册"＋台账**那两列。★`git grep -c` 那形⛔ 能替代这格。
2. ★它只回答"这枚红⛔ 有票管吗"，**⛔ 回答"该不该算这枚票的新红"**——后者＝两发红名册的**差集**＋`rev-parse` 那把 blob 尺（见 (a)/(d)）。★第一版尺没分"判据格⇄散文提及"时，一枚红被 13–25 枚票"提及"＝归因等于没归因；★第一版按"枚×枚 grep 全票池"跑＝**11k+ 进程、240s⛔ 完**，所以 stage 4 改成**一次建 `crit.tsv` 索引**。★这两条都是本腿踩出来的，抄件的腿⛔ 再踩一遍。

---

## 四、★四组新读数（票面与上一趟都⛔ 有的、直接影响裁形）

1. ★**票面两枚今天的色＝连 8 发全绿**（尺＝`gh run view --log --job <test-core databaseId>` 逐发顶格取；41372451347→41379424566，含 HEAD 那发 **41376302416**）：`…322r1=G / LeafSwap=G / 两枚祖先=G` ×8。★绿那发的到达数**非零**＝`flips=504`、`leafLegRefusals=43`、`refused-at-act=10`、`sealed-own=159` ⇒ ★**真绿、⛔ 读者档藏出来的绿**。★两把尺各自证这一句：① `:242-244` 那行「AC#1 reader mode (WISP_120_BASELINE=1)…」只要档置位就**无条件打**，★两发 log（绿⇄红）里 grep 该串各＝**0 命中**⇒ **CI 从头到尾⛔ 置这枚档**；② 红那发②那行 `t.Errorf` **真响了**＝置位时它只会是 `t.Logf`⇒ 红发同样⛔ 在读读者档。
2. ★**但⛔ 由此读成"本票已自愈"**：那 8 发的锚**全在 `winsec.go` 已变之后**（`cbafe712…`、30 笔 `internal/winsec` commit、其中票 400-r1 逐字叫「winsec 342r1 **竞态修**」）⇒ ★**它们⛔ 覆盖票面两锚、⛔ 证明"弃权被台面安静化治好"**；★本腿⛔ 跑过任何 `go test`、★**runner 负载维＝未知**（⛔ 一个负载读数都没取）。⇒ 结案判据仍是票面 `AC#4`（两平台成对＋各 ≥3 发独立复跑）。★**这一条是给裁形的：甲的"分母值不值"⛔ 在今天的色上**。
3. ★**票 369 已把这族的"乙"落地过一次**（`-done`），其判语逐字给了 POSIX／Windows 分野：★「票 346 的射程只到 POSIX」（`369-…-done.md:16-22`：两枚首行 `//go:build !windows` ⇒ Windows runner 上⛔ 编译 ⇒ 今天红在 Windows 的那两枚⛔ 在 346 射程、当时**无主**）＋★它用 `//go:build windows && wisp342probe` 把一枚胜率探针**移出默认构建**、编排者逐字判「build tag 比搬文件好…我认这个偏离」。⇒ ★**两件事进裁形料**：① ★同族一 half 的乙**已经付过、付得起**；② ★"用 tag 让仪器⛔ 进 CI 门"这一手**在本仓有裁过的先例**——但★它对票 346 两枚＝**票 112 `AC#2` 逐字禁的形**（⛔ 加 `//go:build`／⛔ 加 `t.Skip` 让它别响）＋本票 `AC#2` 同禁 ⇒ ★**⛔ 把它当甲／乙／丙之外的第四形端出来**，除非另票＋人工批准。
4. ★`346-a1` 那趟说的常驻件 **`scripts/abstention-census.sh` 今天盘上⛔ 存在**（尺＝`git ls-files scripts/` 全列＋`glob **/abstention-census*`＝**0 命中**）＝它当年是"只给形状"、⛔ 落盘。⇒ ★若还要那枚件，需要一枚有写面的腿；★本件第三节的尺今天已能做 AC#3 要的**归因**那一半。

---

## 五、写面自证尺（本腿自己现跑）

- 入库件＝**1 枚**＝本件。★⛔ 新落 `probes/346/**`（派单＝只读）；⛔ 动 `.go`／票面／台账；⛔ push。
- 被禁字面字符＝码位 **U+00A7**：本件命中 **0**（★与间隔号 U+00B7 两码位分量、⛔ 混）；CR 字节＝**0**（全 LF）；`.out` 扩展名枚数＝**0**；零字节枚数＝**0**。
- ⛔ 产码尺（本腿复跑）＝`git diff --cached --name-only`＝**0**、`git diff --name-only`＝**0**；`internal/winsec/**` 经本腿的手＝**只读**（`read`／`git grep`／`sed -n`，★0 次写入、★0 次 `go test`）。
- ★票面⛔ 动尺＝`git status --porcelain -- .scratch/wisp/issues/`＝**0 枚**；`Status:`／⛔ 翻框复量＝判据区 **open 7／closed 0**（与本腿起手时一字不差）。
- 临时件只建不删＝`D:\tmp\wisp\20261102\346c1\`：两枚原始 CI log（**2358566B** ⇄ **2795862B**）＋`spread/core-<run>.log` 8 份＋尺三枚（`roster.ps1`／`family.ps1`／`ci-red-attribution.sh`）。★原始 CI 全量日志⛔ 入库（别家 CI 的字节）。
