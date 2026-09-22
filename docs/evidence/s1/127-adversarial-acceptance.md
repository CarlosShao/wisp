# 票 127 独立对抗验收 —— 常驻 GUI 腿那 20 行 install 的第一枚钉（`R-117-A`，票 117 结案的唯一条件）

**验收方**：`acceptor-ticket127`（本文件唯一作者）
**日期**：2026-09-22（本会话）
**被验交件**：票 127，实现方 `agent-ticket127`
**基线**：`6a39820`（票 121 收件、本票开工前的树；实现方选它作控制组，我沿用）
**交件 sha**：`8663a39`（用例 `cmd/wisp/resident_sink_nail_127_windows_test.go` 533 行）、
`c927738`（AC#1 票面）、`938fda1`（AC#2 两处注释 + 票 117「更正一」）、`fd82bf4`（AC#3 票 117「更正二」）、
`9c335d0`（AC#4 票面）、`8ec04f3`（收尾补记，**最终 sha**）
**本文件性质**：1 AC = 1 格，裁决 + 三档标签；**票面一字未改、别人的框未翻**；`R-*` 只给建议归属，
`docs/reports/**` 一字未动（那本账编排者记）。

---

## 一、AC#1 常驻腿的钉 / M4 重量 —— **通过**〔独立复现〕

### 1.1 M4 我自己复算了一遍（原文照贴）

快照 `/tmp/ac127-mut-m4` ＝ `git archive 8ec04f3 | tar -x -C /tmp/ac127-s22` 的**整份副本** + `sed -i '46,65d' cmd/wisp/resident_windows.go`。

**先证落地**（规矩：变异先证那几字节，再谈红名）：

```
diff <(post) <(mut)  →  46,65d45  —— 被删的正是 20 行，逐行是
	// Ticket 117: the resident process is the leg owner actually uses - double
	... (2 行注释) ...
	sink, sinkErr := installLogSink(rt.Layout.DataDir)
	if sinkErr != nil { ... } else { defer sink.close() }      （注释 11 行 + 代码 9 行）
grep -n "installLogSink\|sink.close()\|持久日志未启用\|Ticket 117" cmd/wisp/resident_windows.go → 零命中，rc=1
wc -l cmd/wisp/resident_windows.go → 85 → 65
```

**再读静态门与红名**（本机 go1.27.1，`PATH` 挂了快照里的三枚 dll）：

| 读数 | 实现方自报 | **我复算** |
|---|---|---|
| `go build ./cmd/wisp` | rc=0 | **rc=0** |
| `go vet ./cmd/wisp/` | rc=0 | **rc=0** |
| `go test -count=1 -v ./cmd/wisp/` rc | 1 | **1** |
| `=== RUN` / PASS / FAIL / SKIP | 76 / 73 / 3 / 0 | **76 / 73（顶层 39＋子用例 34）/ 3 / 0** |
| 墙钟 | 72.368s | **74.899s** |

红名**逐字三枚，没有第四枚**：

```
--- FAIL: TestAC1ResidentLegInstallsItsLogListenerOnDisk (13.08s)
--- FAIL: TestAC1ResidentLegOutlivesItsOwnLogFailure (13.25s)
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (13.45s)
```

第一枚失败正文末句就是本票的立项理由（我量到的原文）：

```
resident_sink_nail_127_windows_test.go:363: the resident process wrote no log file under
C:\Users\...\TestAC1ResidentLegInstallsItsLogListenerOnDisk1973562646\002\logs within the budget.
    that is the claim this case exists for: the leg owner actually runs has no listener on disk.
```

**特异性**（"不许是顺带把别的用例拖红"）：同一份变异里 run 腿那六枚 AC#2/AC#3 用例
（`TestAC2SealNoticeLandsInTheRunLegLogFile`、`TestAC2AuditTrailLandsInTheRunLegLogFile`、
`TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`、`TestAC3LogSinkLandsInsideTheEnvDataRoot`、
`TestAC3EmptyDataRootIsARefusalNotAFallback`、`TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole`）
**全部 PASS** ⇒ 红的确实只有被删掉的那条腿。（⚠ 票面 AC#1 写"run 腿那三枚"，实际这条腿有六枚；
方向是"更多枚仍绿"，不是错判，属计数措辞。）

⇒ `acceptor-ticket117` 给的唯一结案判据 **R-117-A**（"删那 20 行必须至少红一条、且红名点到常驻腿"）**成立**。

### 1.2 首要攻击点：这到底是真生产装配根，还是测试自己搭的形状

票 117 的立项理由逐字是"**测试里手搓 sink 不算**"，所以这一格必须先判。链条我逐环走过：

1. 三枚用例的驱动器都是 `buildWispForTest(t)`（`cmd/wisp/secret_argv_windows_test.go:161`）——
   它 `go build -o <tmp>/wisp.exe ./cmd/wisp`，再把 `third_party/sherpa-onnx/*.dll` 三枚拷到 exe 旁边。
   ⇒ 驱动的是**从被测树现编的生产二进制**，不是测试里的函数。
   ⚠ **票面 AC#1 那句"shipped `wisp.exe`"用词不准**：仓根确实躺着一枚 `wisp.exe`，
   而用例**没有**用它——**这是好事**：若真用那枚陈旧产物，删源码 20 行**不可能**让它红，
   那才是假钉。我按事实判：形状是"现编生产二进制"，比票面自述更强（见 `R-127-2`）。
2. `exec.Command(exe)` **零参数**（`resident_sink_nail_127_windows_test.go:155`）⇒ `main()` 走
   `len(args)==0` 那一支（`cmd/wisp/main.go:51-57`）⇒ `attachParentConsole()` ⇒ **`runResident()`** ⇒
   **`installLogSink(rt.Layout.DataDir)`（`resident_windows.go:57`）**。
   `attachParentConsole`（`console_windows.go:33-43`）在 stdout 已是有效句柄时**直接 return**，
   失败路径也只是 `_, _, _ =` 忽略 ⇒ 它不可能绕过 install，也不会 `os.Exit`。
3. 测试**没有**自建 sink、**没有**调 `installLogSink`、**没有**换 `slog.SetDefault`：
   `grep -n "installLogSink\|SetDefault\|InitLog" cmd/wisp/resident_sink_nail_127_windows_test.go` ⇒ **0 命中**。
   它只读盘上那个文件。⇒ 判**真生产装配根**，不是"测试自己搭的形状"。

**`WISP_ENV=test` 那一手是不是安全的**（这条如果被证伪，测试会在别人机器上激活 owner 正在跑的 Wisp）：

- 静态：`internal/proc/envfork.go:78-86` —— test env 的 `Layout` 是 `{DataDir: TestDataDir(), MutexEnabled: false}`；
  `internal/proc/boot_windows.go:93-103` —— `AcquireSingleInstance` **整段在 `if rt.Layout.MutexEnabled` 里**，
  所以 test env 既拿不到互斥、也不可能返回 `ErrAlreadyRunning`；
  而 `SignalExistingInstance` 只在 `runResident` 的 `errors.Is(err, proc.ErrAlreadyRunning)` 分支里被调
  （`resident_windows.go:34-41`）⇒ **没有任何一条路能让 test-env 子进程去激活别人的实例**。
- 真机：我自己手跑了一发这条腿（不经过用例），子进程 stdout 原文：
  `wisp: resident runtime booted (Wisp · test, data dir = C:\...\ac127-probe-data, portable = false,
  job object = on, **single instance = false**)` —— 那半句就是 `rt.Instance != nil` 打出来的，
  即**它确实没有注册会话互斥**。⇒ 这条**证伪不了**，实现方的说法成立。
- 写件边界：`WISP_TEST_DATA_DIR` 是 identity contract（`envfork.go` 的注释与 `TestDataDir()` 逐字回传），
  用例第 383-386 行拿它做**精确相等**比较；再用"盘上只许在 `<data>\logs` 出现 `*.jsonl`"
  那一圈 `jsonlFilesUnder` 兜住。我手跑那一发：`ls /tmp/ac127-probe-data/logs/` ⇒
  **只有 `wisp-20260922-001.jsonl` 一枚**，别处零枚。

### 1.3 "永远绿"的假钉形状：逐项查过，全部不存在

| 形状 | 读数 |
|---|---|
| `t.Skip` | 全文 **0** 处（唯一命中是第 32 行注释里那句"无 t.Skip"） |
| `Fatalf` 降级成 `Logf` | `t.Fatalf` **15**、`t.Errorf` **14**、`t.Logf` **3**；三枚 `Logf` 是第 296/471/515 行的**纯报告**（读数、退出状态、RESIDENT LEG RECORDED），不承载任何判定 |
| 断言只查文件存在 | **不是**：`readResidentSink` 先 `observe.CountLogFiles`（拒绝没有管道文件），再逐行 `json.Unmarshal`，然后比 `msg` **逐字相等**、`level=="INFO"`、`dir` **精确相等**、`min_level=="info"`；`sinkHasARecord` 明确"文件在 ≠ 听众记了账"（注释第 320-324 行）并在等待时看**记录**而非文件 |
| 记录可能来自别的写者 | 用例要求同一句**同时**出现在**那个子进程自己的 stderr** 上（第 395 行），并在删除 `install` 后仍跑到事件循环（第 398 行）⇒ 我复算的红名原文里那条 jsonl 归零、控制台也没有该句 |
| 子进程可能活过用例 | `t.Cleanup(leg.stop)` + `stop()` 有界收割（200×50ms）；`CREATE_NEW_PROCESS_GROUP` 只屏蔽 CTRL_C 不屏蔽 CTRL_BREAK，所以第 3 枚能要求真实退出路径 |

**我自己补的三发变异**（都先证那几字节落地、再 `go build` rc=0、才读红名）：

| 变异 | 落地证据 | build/vet | 读数 |
|---|---|---|---|
| **装到别处去**：`installLogSink(rt.Layout.DataDir)` → `... DataDir + "-elsewhere")` | `diff` 只有 57 行那一处 | **rc=0** | 三枚**全红 6/6**（`-count=2`）：听众装了，但装的**不是这条腿的数据根** ⇒ 判据不是"有个 jsonl 就行" |
| **`dir` 属性说谎**：`logsink.go:165` 的 `"dir", dir` → `"dir", logDirName` | `diff` 只有 165 行那一处 | **rc=0** | 第 1 枚**红**、另两枚**绿**，原文 `install record names dir = "logs", want "C:\\...\\logs"` ⇒ `dir` 那格是真的在比内容 |
| **摘掉 `close` 但骗过编译器**：`defer sink.close()` → `_ = sink` | `diff` 只有 64 行那一处 | **rc=0 / rc=0** | 见 1.4——这一发同时**否证**了实现方的一处说法 |

### 1.4 "拆一半那一发"的判决（派单点名要我判的一格）

实现方自报：只删 `defer sink.close()` ⇒ `go build` **rc=1**
（`cmd\wisp\resident_windows.go:57:2: declared and not used: sink`），结论写的是
"这一行'**在不在**'由**编译器**钉着，不需要用例"。

- 它那发我复算为真（编译器确实会在"删干净"时拦下）。
- **但"编译器钉着这一行在不在"这句被我的第三发否证**：把 `defer sink.close()` 换成 `_ = sink`，
  编译器只在乎 `sink` 被**用到**，不在乎 `close()` 有没有被调用 ⇒ `go build` **rc=0**、`go vet` **rc=0**。
  这一发是真实形状（"记得要收尾却忘了收尾"），编译器**不**是防线。
- **防线确实在用例上**：`-count=4 -run TestAC1ResidentLeg` ⇒
  第 1、2 枚 PASS、第 3 枚 **FAIL**，且失败形态是
  `panic: runtime error: index out of range [0] with length 0 [recovered, repanicked]`
  ——`readResidentSink` 返回空切片后第 495 行直接取 `recs[0]`（第 1 枚用例有 370-372 的长度守卫，**第 3 枚没有**）。
  机制：没有 `close()`，而 500ms tick 又没在关停窗口里抢到，**一条记录都没落盘**。
- 判：**这一格达标**（它防的结局被真实造出来时是用例红，不是绿），
  但实现方那句"编译器钉着"是**过度归因**（`R-127-3`），
  红名被 panic 吃掉、且会**截断整包后面的用例**是**仪器形状缺陷**（`R-127-4`）。
  ⚠ 两回事都不构成退回：这条格子的 AC#1 判据是 M4 那一发，它红名逐字、只点常驻腿。
- 另：它自报的第三发（两枚 defer 换注册顺序）我复算为 `go build`/`go vet` 干净、
  `-count=4` **4/4 红**、原文与它引的逐字相同（`the D38(e) trail is not in the file the same process closed:
  0 record(s) after the install, all msgs [wisp: persistent log sink installed]`）。
  它给的"~2% flush-tick 折扣"我量不出来（4/4），方向是**它把自己的读数说弱了**——保守一侧，登记不追责。

### 1.5 判

**通过**〔独立复现〕。钉子是**真钉**：驱动现编的生产二进制零参数、读盘上记录、断言比内容不比存在，
删那 20 行 ⇒ 三枚按名字红、run 腿六枚仍绿；`WISP_ENV=test` 的"不会去激活 owner 的 Wisp"静态＋真机两头都立得住。
两条附言（`R-127-3` 归因用词、`R-127-4` 缺长度守卫导致 panic 红）都是**改法一行的仪器条件**，
不改变本格裁决，也**不阻断票 127/票 117 结案**。

---

## 二、AC#2 `R-117-C` 的降级是不是诚实 —— **通过**〔独立复现〕

**先说结论**：三处落点都在、话与事实一致，两条"为什么不补覆盖"的理由都独立复算为真。

| 复算项 | 实现方读数 | 我复算（快照 `8ec04f3` / 控制 `6a39820`） |
|---|---|---|
| `GOOS=linux go vet ./internal/observe/` | rc=0 | **rc=0**（本机 go1.27.1） |
| `GOOS=linux go list -deps ./internal/observe/ \| grep -c "CarlosShao/wisp/cmd/wisp"` | 0 | **0** ⇒ over-claim 的机械证明成立 |
| `GOOS=linux go vet ./cmd/wisp/` | rc=1，三行原文 | **rc=1**，三行原文逐字相同（`imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8`） |
| 同一条命令在 `6a39820` 纯净快照 | 逐字相同 | **逐字相同** ⇒ "既有形状、非票 117/127 引入"成立 |
| CI 腿：跑 `./cmd/wisp/` 的只有 windows | 是 | **是**：`ci.yml` 里 `test-windows`（`runs-on: windows-latest`）第 369 行
  `"cmd/wisp CLI tests ..."` → `bash scripts/wisp-cli-tests.sh` → `portable-tests.sh --scope=cli` →
  `scope=(./cmd/wisp/)`；而 ubuntu 的 `test-core` 第 255 行走 `--scope=core`，那份 16 项清单里**没有** `./cmd/wisp/`
  （我把 `core)` 那一支整段读过）。⇒ "本票三枚 `//go:build windows` 与票 117 三枚 untagged 都只在 Windows 跑"**成立** |

**两处代码注释落点**：
`cmd/wisp/logsink.go` 头 `PLATFORM` 段（第 48-58 行，明写"this file deliberately imports nothing
platform-specific … Until then the claim is unproven, not proven"）、
`cmd/wisp/logsink_test.go` 头 `PLATFORM LEG` 段（第 16-32 行，引了我逐字复算到的那行 sherpa 错误原文、
"measured identical at 6a39820 and at 8663a39"、以及 `go list -deps` 的 0 命中）。
⇒ **两段话里的每一条我都能自己复算**，没有一处留在口头。

**票面落点**：票 117 第 547-575 行「更正一」——§六 原文照引、只作废括号那半句、附同一张四行复算表，
并且**主动把第二半**（那三条 untagged 用例在任何平台的 CI 上只在 Windows 执行）**一起算了**。

**为什么不选"补覆盖"那支**：两条理由复算为硬——
`ci.yml` 是票面点名的禁改件（票 117 第 12 行），且 ubuntu 腿连编译都过不了（我量到 rc=1）。
这一支不是"降级最省事"，是**另一支不存在**；票面上把这句话写成了"两条都是硬约束"，与事实一致。

**判**：**通过**〔独立复现〕。"如实降级"这条选择**没有**藏东西：三处都可 grep、每处都带可复算命令，
且我把命令重跑了一遍。

---

## 0b. 会话口径补记（骨架里那节的正文，位置在此以免与下面的裁决格混排）

- **只读**：本会话在仓库内唯一的写件是本文件。全程
  `git archive <sha> | tar -x -C /tmp/ac127-s22/<名>` 取快照，**没有**在仓库树内建 worktree、
  没有 `checkout`/`--amend`/`reset`/`rebase`/`stash`（A38④）。
  `internal/winsec/**` 有另一个写者在飞（票 129）⇒ 本会话**没有一次**读工作树的 winsec 状态、
  也没有对它做过写，一切走快照。
- **本机命令前缀**：`PATH` 必须挂上快照里的 `third_party/sherpa-onnx`，否则 `cmd/wisp` 的测试二进制
  在**进程加载**时就死（`exit status 0xc0000135` = STATUS_DLL_NOT_FOUND，票 98 那一格）。
  ⚠ 这条是**我自己踩到的**：第一次跑 `-run TestAC1ResidentLeg` 得到
  `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.042s` + `rc=1`、`=== RUN` 命中 **0** 条 ——
  四数全零而 rc 非零，是"没跑"不是"跑红"。`git archive` 不含 dll ⇒ 三枚 dll 要自己拷进快照
  （`onnxruntime.dll` / `sherpa-onnx-c-api.dll` / `sherpa-onnx-cxx-api.dll`），
  这与实现方 AC#1/AC#4 两格自报的"票 117 §六 那枚仪器坑照抄"是同一件事，**复算成立**。
- 快照清单（**全部在仓外**）：
  | 名 | 内容 | 用途 |
  |---|---|---|
  | `/tmp/ac127-s22` | `git archive 8ec04f3`（最终 sha）+ 三枚 dll | 交件树：基线四数、AC#3 分母与行号、门禁、台账、手跑那一发 |
  | `/tmp/ac127-ctrl` | `git archive 6a39820` + 三枚 dll | 控制组：146/78 四数、linux vet 同形、台账 31 |
  | `/tmp/ac127-mut-m4` | `8ec04f3` 副本 + 删 `resident_windows.go:46-65` 整 20 行 | **M4 重量（本票生死格）** |
  | `/tmp/ac127-mut-elsewhere` | 副本 + 第 57 行数据根改成 `DataDir + "-elsewhere"` | 变异：装了、但装到别的数据根 |
  | `/tmp/ac127-mut-dirlie` | 副本 + `logsink.go:165` 的 `"dir", dir` → `"dir", logDirName` | 变异：install 记录的说辞说谎 |
  | `/tmp/ac127-mut-noclose` | 副本 + 第 64 行 `defer sink.close()` → `_ = sink` | 变异：**骗过编译器的"忘了收尾"** |
  | `/tmp/ac127-mut-order` | 副本 + 两枚 defer 换注册顺序（close 后置） | 复算实现方自报的第三发 |
- 手跑那一发的产物（`/tmp/ac127-probe.exe`、`/tmp/ac127-probe-data`、`/tmp/ac127-probe-{out,err}.txt`）
  跑完 `taskkill` 收掉了子进程；仓内**没有**留下任何探针文件。
- 四数一律从 `-v` 输出量，`-count=2` 才不缓存：`grep -c '^=== RUN'` / `'^--- PASS'` / `'^--- FAIL'` /
  `'^--- SKIP'`（顶层锚定，子用例缩进另计）。
- 变异先发落地再读名：`grep`/`diff` 出被改的那几字节 → `go build ./cmd/wisp` rc=0 → 才读 `--- FAIL` 名单。
- `GOOS=linux go vet` 只编译不执行。

**本文件写作状态：AC#1、AC#2 两格已落盘（渐进写，每格一次 commit）；AC#3、AC#4、结案判在文末续。**
逐格裁决一览（未落格的写"待"）：

| 格 | 裁决 | 标签 |
|---|---|---|
| AC#1 常驻腿钉子 / M4 重量 | **通过**（附两条仪器条件 1.4，不阻断） | 〔独立复现〕 |
| AC#2 R-117-C 降级是否诚实 | **通过**（三处落点与四条理由逐字复算） | 〔独立复现〕 |
| AC#3 票 117 §一 现状表重算 | **通过**（三口径同数 40，逐格点名复算） | 〔独立复现〕 |
| AC#4 门禁 + 四数 + 台账 | **通过**（四数与台账八数我独立量到同数） | 〔独立复现〕 |
| **总判** | 见文末 §七 | 见文末 |

---

## 三、AC#3 票 117「更正二」那张现状表 —— **通过**〔独立复现〕

### 3.1 分母我自己跑了一遍（派单点名要的那一行）

口径行照票面逐字跑在快照 `/tmp/ac127-s22`（`8ec04f3`）：

```
grep -rnE "slog\.(Warn|Error)\(" --include=*.go internal cmd | grep -v _test.go \
  | grep -vE '^[^:]+:[0-9]+:[[:space:]]*//'          ⇒  40
```

**40 枚，分布逐包与票面一字不差**：
`ball 17 / audio 5 / winsec 4 / observe 4 / proc 3 / plugin 2 / config 2 / statemachine 1 / secret 1 / cmd/wisp 1`。

⚠ **这一格我额外量了一件事：这把尺子为什么今天可信**。那条管道里 `grep -v _test.go` 比的是**整行文本**
（路径在前缀里，所以能过滤测试文件），但**同时**会误伤"正文提到 `_test.go` 的非测试行"。我把三种口径并排跑：

| 口径 | 读数 |
|---|---|
| A 票面原样（文本过滤） | **40** |
| B 换成严格路径过滤 `^[^:]+_test\.go:[0-9]+:` | **40** |
| C 去掉"注释行"那一层过滤 | **40** |

⇒ 三个口径**同数**，且 A 的集合里正文含 `_test.go` 的非注释命中是 **0**。
所以 40 不是靠一把歪尺子凑出来的；但**下一个人换包数时 A 与 B 就可能分歧**（`R-127-5`，建议口径行改写成路径过滤版）。

### 3.2 派单点名的每一格我都逐行打开看过（不是抄票面）

| 票面主张 | 我的复算 |
|---|---|
| 漏项 5 枚坐实 | `statemachine/machine.go:139` = `slog.Error("state machine rejected transition"…` ✓；`plugin/disposal.go:208` = `slog.Error("disposal_incomplete: Defer called after Dispose"` ✓、`:332` = `slog.Error("disposal step failed"…` ✓，且它**在 `for _, f := range res.Failed` 循环体内**（325-333 行整段读过）⇒ "一次关停可产出 N 条"为真；`observe/goroutine.go:271` = `slog.Warn("goroutine outside the D38 roster…"` ✓（在 `Spawn` 的 `cat == CategoryUnknown` 分支里 ⇒ 逐 spawn 一条）；`observe/logging.go:204` + `:307` 两句同为 `observe: log retention sweep failed` ✓ |
| audio **5**（不是 4）、ball **17**（不是 14），并给逐文件分解 | ✓ 逐文件复算：ball = `ball_windows.go` **7**、`hotkey_windows.go` **6**、`hotkey_reload.go` **3**、`renderer_windows.go` **1**；audio = `audio.go` **2**、`mmdevice_windows.go` **2**、`wasapimic_windows.go` **1**。逐行坐标也对：`audio.go:136/:152`、`mmdevice_windows.go:60/:65`、`wasapimic_windows.go:112`、`ball_windows.go:393`(`slow drawFrame`)/`:417`(`slow ULW frame`) |
| `run.go:347/:421` → 实为 `:355/:429`（+8） | ✓ `run.go:355` = `rt.auditf("perm: MODE-READ origin=startup…`、`:429` = `rt.auditf("perm: MODE-READ-FAILED…`，与票面标的两枚事件名逐字对上 |
| `secret.go:471` WARN / `:473`、`:326` INFO → **漂 +19** ⇒ `:490` WARN、`:492` INFO、`:345` INFO | ✓（**这一格我自己差点读错，写下来免得下一个也错**：`sed -n '490p;492p;345p'` 按**文件升序**输出，不按参数顺序；按升序还原 ⇒ `:490` WARN、`:492` `slog.Info(audit)`、`:345` `slog.Info("wisp secret: stored dpapi blob"…` ⇒ 三条**级别标注全对**，且相对票 117 交件树 `ce666ea` **恰好统一 +19**：326→345、471→490、473→492） |
| 行号未漂的照实登记（不"整张表全错"） | ✓ `proc/shutdown.go:130` INFO（**就是本票钉子读的那句** `shutdown step skipped (module not present)`）、`:145/:147` ERROR、`:172` WARN；`winsec/resolve.go:146` INFO / `:150/:157/:165` ERROR / `:175` Debug / `:181` INFO；`risk/winsec_c26.go:20` 是 `func init()`、`:21` 是 `winsec.SetPathResolver(c26Pipeline{})`；`secret/migrate.go:198` WARN ⇒ 全部复算为真 |
| 「腿不是两条是三条」 | ✓ `grep -rn "installLogSink(" --include=*.go cmd internal \| grep -v _test.go` ⇒ `run.go:164`、`resident_windows.go:57`、`models.go:276` 三处，**没有第四处**；`wisp secret` 那条腿**仍无听众**（`cmd/wisp/secret.go` 内 `installLogSink` 0 命中）✓ |
| 「ball/audio 今天都不在 wisp.exe 的任何一条腿上」 | ✓ `internal/ball` 非测试 importer 只有 `cmd/balldebug/main.go`；`internal/audio` 非测试 importer **0**；`cmd/wisp/**` 非测试文件里两个包各 **0** 命中（只有 `notify_windows.go:13` 一句注释提到它）⇒ C 段"下一次测量的起点是 22 枚"的**前提**（还没进二进制）成立 |
| `statemachine` / `plugin` 的可及性判语 | ✓ `statemachine` 非测试 importer 含 `cmd/wisp/models.go`、`internal/models/bridge.go`、`internal/ball/*`（+`cmd/balldebug`）⇒ run/models 腿可达、常驻腿今天不可达；`plugin` 唯一非测试 importer = `internal/memory/retention.go` ⇒ 走记忆保留的 run 腿可达 |

### 3.3 它自己新增那两条（不在派单里）

1. **`:204` 与 `:307` 必须分两格记**——机制我读码复算：`:204` 在 `InitLog` 里（`rollingWriter` 起手 sweep 失败），
   而换默认 logger 是**回到 `installLogSink` 之后**才做的（`logsink.go:148` `slog.SetDefault(teeHandler{…})`），
   `internal/observe` 自己的 `InstallAsDefault()`（`logging.go:107`）这条腿没人调 ⇒
   **它永远进不了它正在抱怨的那本文件**。`:307` 在周期 sweep 里、装在换默认之后 ⇒ 会进文件。判**为真**。
2. **听众装不上那一瞬听众自己是哑的**——我把自己手跑那一发拿来验**同族形状**：
   子进程 stderr 两行是
   `2026-09-22 22:53:39 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2`
   （**stock 默认 logger 的时间格式**，没有 `time=` 前缀）与
   `time=2026-09-22T22:53:39.654+08:00 level=INFO msg="wisp: persistent log sink installed"…`（tee mirror 的格式），
   而盘上那本 jsonl **只有第二条**（1 条记录）⇒ **"换默认之前出声的记录只在 stderr"我在常驻腿上独立复现**。
   这一发同时把 **R-125-3 的机制**量清了（见 §五）。

**判**：**通过**〔独立复现〕。这张表现在是本仓**唯一一把我对过三口径的分母尺**；两条新账都成立。
**我没有动 `internal/observe`**：派单已把那一格登记为 `A99③`/待立票 130，本会话对该包**零写**。

---

## 四、AC#4 门禁 / 四数 / 墙钟账 / 台账 —— **通过**〔独立复现〕

### 4.1 四数（两枚快照、`-count=2`、`-v` 量）

| 包 | 组 | rc | `=== RUN`（全部 / 顶层） | PASS（全部 / 顶层） | FAIL | SKIP | 墙钟 |
|---|---|---|---|---|---|---|---|
| `./cmd/wisp/` | 控制 `6a39820` | 0 | **146** / 78 | **146** / 78 | **0** | **0** | 72.879s（自报 67.681s） |
| `./cmd/wisp/` | **交件 `8ec04f3`** | **0** | **152** / **84** | **152** / **84** | **0** | **0** | 92.567s（自报 89.758s） |
| `./internal/observe/` | 交件 | 0 | **94** / 94 | **94** / 94 | 0 | 0 | 2.296s（自报 2.316s） |

**增量自洽**：`152-146 = 6` ＝ 3 枚新用例 × `-count=2`；顶层 `84-78 = 6` 同解；
子用例两侧都是 **68** ⇒ "本票没加子用例"是**能核的**（我独立量到）。
我读 PASS 的口径与它一致：`^--- PASS` **84** ＋缩进 `^    --- PASS` **68** ＝ 152；
票面把"全部/顶层"两栏并列的写法**没有骗人**，沿用。墙钟差 4-5s 是本机负载，不构成读数分歧。

### 4.2 格式与静态门（都真跑在快照里）

- `gofmt -l cmd/wisp/` ⇒ **空**；`"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/ internal/observe/` ⇒ **空**；
  `gofumpt --version` ⇒ **`v0.7.0 (go1.27.1)`**（与票面引的版本串逐字相同）⇒
  "写'未跑'必须引命令原文"那一格**不适用**，两把都真跑了；票面点名的 `logsink_windows_test.go` 未被点名。
- `go vet ./cmd/wisp/ ./internal/observe/` ⇒ **rc=0**。
- `sh scripts/d22scan.sh`（交件快照）⇒ **rc=0 / clean**，正向控制原文
  `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`
  ——**逐字**与它自报相同。

### 4.3 台账八 scope：交件 vs 同 sha 控制组

| scope | 控制 `6a39820` | 交件 `8ec04f3` | 判 |
|---|---|---|---|
| bans #1-5 `internal/` | 202 | **202** | 持平 |
| bans #1-5 `cmd/` | 22 | **22** | 持平 |
| ban #6 `frontend/` | 40 | **40** | 持平 |
| ban #7 `internal/tools/` | 18 | **18** | 持平 |
| ban #8 `design/` | 16 | **16** | 持平 |
| ban #8 `frontend/` | 40 | **40** | 持平 |
| ban #8 `internal/` | 382 | **382** | 持平（它引的 382 我逐字复算为真） |
| **ban #8 `cmd/`** | **31** | **32** | **+1，来源我自己定位到一枚文件** |

**+1 的来源我不采信它的说法**：两棵树 `find cmd -name '*.go'` 求差 ⇒
**唯一差异**是 `cmd/wisp/resident_sink_nail_127_windows_test.go`（31 → 32 枚）⇒ 它的归因**逐字成立**。
⚠ 票面 AC#4 自写的基线 **29** 是票 117 时期的数，同 sha 控制组今天是 **31**；
它按"同 sha 控制组为准"处理并把矛盾自己翻出来登记，是**对的**（这正是本仓要的形状）。
零 scope 下降；`tools/d22scan/**`、`allowlist.txt`、任何阈值/golden 我复算**一字未动**（见 4.4）。

### 4.4 "没顺手动掉票 123 那批 CLI 用例"的判据

`git diff --numstat 6a39820..fd82bf4 -- cmd/wisp/` ⇒
`11 0 logsink.go`、`18 0 logsink_test.go`、`533 0 resident_sink_nail_127_windows_test.go`
= **3 files changed, 562 insertions(+), 0 deletions**（逐字复算为真）⇒ 纯插入，**没有一条既有用例被改或被删**。
`git diff --name-only fd82bf4..9c335d0` ＝ 2 枚 `.md`、`9c335d0..8ec04f3` ＝ 1 枚 `.md`
⇒ 最终 sha 的**码面**与 `fd82bf4` 相同，所以 4.1 那四数可以直接记在它的交件码面上（这条推理我核过）。
票 117 票面：`git diff --numstat 6a39820 8ec04f3 -- …117….md` ⇒ **72 / 0**（零删除）⇒
"§六 原文一字不删、只在其后作废"**是真的**；127 票面那 4 处删除**全部**是 `- [ ] **AC#N**` → `- [x]`（自己那四格），
append-only 未破、别人的框一枚没翻。

**判**：**通过**〔独立复现〕。四数、台账、格式门、静态门、diff 口径五件事逐条重量，无一项需要背它的数。
