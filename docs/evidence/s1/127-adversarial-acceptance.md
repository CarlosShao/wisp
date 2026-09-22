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
| AC#3 票 117 §一 现状表重算 | 待 | 待 |
| AC#4 门禁 + 四数 + 台账 | 待 | 待 |
| **总判** | 待 | 待 |
