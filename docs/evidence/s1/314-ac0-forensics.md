# 票 314 AC#0 — `go test -c` 出的 `cmd/wisp` 测试二进制**直跑**：整进程 native 消失**复现**（`0xC0000005`）

- **腿**：314 AC#0 只读取证腿（⛔ 非写腿；⛔ 未翻任何 AC 框；⛔ 未动产码/断言/票面；⛔ 未 kill 球 PID 49348）
- **时刻**：2026-11-15 12:16–13:25 +08（本机 Windows，`go1.27.1 windows/amd64`）
- **票面 AC#0 逐字**：「`go test -c` 编出 cmd/wisp 测试二进制 → **直跑它**（`-test.v -test.skip='…'`），拿**进程退出码**（native 异常码 vs `os.Exit(n)` vs `-test.timeout` panic）＋stderr 全量；⛔ 只看 wrapper 汇总。」
- **读数根**：`D:\tmp\wisp\20261115\w314\`（只建不删；`run1-canonical\` 是第 1 发的原样副本）；摘录副本 `.scratch/wisp/probes/314/ac0/`
- **发射器**（可逐字复跑）：`.scratch/wisp/probes/314/ac0/{run-direct-canonical.ps1, gate-and-run.ps1, desktop-state.ps1}`

---

## 0. 一句话结论

**复现了——2 发同形里 1 发死。** 第 1 发直跑 `cmd/wisp` 的 winlive 测试二进制，跑到第 **554.5 s**、第 **348** 枚顶层用例（`TestAC4PriorFocusSurvivesARefusedPanelSample`，`cmd/wisp/panel_resident_windows_test.go:1203`）时，整进程**不带任何 Go 层信号**消失：**进程退出码 `0xC0000005`（dec `-1073741819`，u32 `3221225477`）＝ native ACCESS_VIOLATION**，**零 `panic:`、零 `fatal error:`、零 `SIGSEGV/signal/exit status` 行、该用例零 verdict、终点锚 `TestZZ315A1` 零命中**；不是 `os.Exit(n)`，也不是 `-test.timeout` panic（只烧了 554.5 s，`-test.timeout=90m` 根本没到期）。第 2 发同形同点把这枚跑成 PASS 并**整包跑完**（`0x00000001`／1389.8 s／669 枚全有 verdict）⇒ 死法是**条件性**的，复现率 **1/2**。

---

## 1. 先把「直跑」的正确发射形钉住（这一步不钉住，读数全是假的）

`go test ./cmd/wisp` 不是"把二进制当普通程序跑"：**cmd/go 把子进程的工作目录设成包目录**，并且载具 `scripts/wisp-cli-staged-tests.sh` 另外做两件事——**把 pin 的 sherpa/onnxruntime DLL 放到二进制旁边**（`:230-247`，`stage-sherpa-dlls.sh:8-21` 逐字：System32 的 `onnxruntime.dll` 1.17.1 会被加载器先搜到），**导出五个 voice fixture 变量**（`:312-331`）。载具自己的发射行逐字（`:492`）：

```bash
(cd cmd/wisp && "$abs_stage/wisp.test.exe" -test.count=1 -test.timeout=90m -test.v "${skip_args[@]}" "$@")
```

**⇒ 本腿的 canonical 形＝ cwd `cmd/wisp` ＋ 三个 pin DLL 拷到二进制旁 ＋ `WISP_{SPIKE_MODELS,SILERO_MODEL,TTS_DIR,KWS_DIR,PUNC_MODEL}` ＋ `WISP_43_REPO_ROOT` ＋ `-test.count=1 -test.timeout=90m -test.v`。**

### 1.1 反面控制（run A）：从仓根直跑 ≠ 同一个测试

同二进制、同参数，只把 cwd 换成仓根、不铺 DLL、不导 env（`run-direct.ps1`）：

| 读数 | 值 |
|---|---|
| 退出码 | `0x00000001`（dec 1，`os.Exit(1)`：测试红，**不是**死法） |
| 耗时 | 1219.7 s（被 `Test320FetchLeg…` 的 396 s 仓根全树 walk 拖长） |
| 顶层 verdict | PASS **534** / FAIL **121** / SKIP **14**（669 枚名册全有 verdict，**无死亡**） |
| `panic:` / `fatal error:` | 0 / 0 |
| 终点锚 | 命中（`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask` PASS） |

**121 枚红里 118 枚逐条可归到"发射形错"**（原文 TSV：`probes/314/ac0/runA-fail-messages.tsv`），另 3 枚单列：

| 类 | 枚数 | 判据（原文摘句） |
|---|---|---|
| cwd：按包相对路径开源文件 | 32 | `open resident_approval_windows.go: The system cannot find the file specified` |
| cwd：目录普查拿到 0 个源文件 | 13 | `no .go files in .: empty instrument`／`扫描集为空：cmd/wisp 下没有 *.go` |
| cwd：`cmd.Dir=../..` 落错父目录 | 25 | `go build ./cmd/wisp/testdata/esclistener: exit status 1`（该测试自设 `cmd.Dir = ../..`，`resident_approval_live_246_windows_test.go:426`） |
| cwd：从仓根走全树 | 1 | `walk D:\…\Wisp: read …\node_modules\…: Incorrect function` |
| cwd：`func main` 不在本目录 | 2 | `func main is not declared in this directory` |
| env：`WISP_SPIKE_MODELS` 未设 | 19 | `install the streaming ASR into <tmp>: WISP_SPIKE_MODELS is not set` |
| env：onnxruntime 被 System32 盖住 | 19 | `onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires` |
| env：`WISP_43_REPO_ROOT` 未设 | 2 | `WISP_43_REPO_ROOT is not set: the staging runner must name the checkout` |
| 驻留球占着全局热键（PID 49348） | 2 | `Hot key is already registered`（`Ctrl+Alt+Q/M/P` 三枚全被占） |
| 桌面前台被别人占（`magpie` 0x5701544） | 3 | `the prior window 0x… did not take the foreground … this is a failed measurement, not a product verdict`（含既知红 `TestAC4FocusReturnToPriorWindowGap33r5`） |

另 3 枚（**不属上表两类，单独具名**，都不在死法射程内）：

| 用例 | 原文摘句 | 判 |
|---|---|---|
| `TestTicket255RealWindowWidthFollowsTheConfig` | `testing.go:1617: TempDir RemoveAll cleanup: unlinkat …\EBWebView\lockfile: The process cannot access the file because it is being used by another process.` | 收尾竞态（Windows 文件锁），**同一枚在 canonical 那发跑到并 PASS** |
| `Test41a2r1OptInGateResolvesTheManifestSpelling` | `the opt-in-off note = "", want the enabled=false sentence` | 与 fixture/构建相关的说明句，未在本发死前跑到 |
| `Test381r2StaleTokenDoesNotKillTheNextAnnouncement` | `MUTANT 381r2-M2 SHAPE: the announcement was killed by a stale interrupt token …` | mutation 自检句，未在本发死前跑到 |

⇒ 反面控制**只用来看"cwd/env 会把 121 枚打成假红"**，⛔ 不拿它当产品判据。

**⇒ 方法论读数（对票 314 直接有用）**：谁要按 AC#0 直跑，**必须**照第 1 节那个形；从仓根裸跑会得到 121 枚假红，把真死法完全淹掉（本腿第一发就是这样，作废重跑）。

---

## 2. 二进制身份（读数属于哪棵树，必须写死）

| 项 | 值 |
|---|---|
| 编译命令 | `go test -c -tags winlive -o D:\tmp\wisp\20261115\w314\wisp-winlive.test.exe ./cmd/wisp`（**仓内零落件**） |
| 编译读数 | rc=0，29.5 s，`63,448,065 B`，SHA256 `1C7F9B324023214601F25B93A24006C7F6FFFBF1A6BFEADAD8896AEEF064029D` |
| `go version -m` | `path github.com/CarlosShao/wisp/cmd/wisp.test`；`build -tags=winlive` |
| 编译时 HEAD | **`bf4adaa7b`**（12:16:45 编出）。⚠ 跑的中途别的腿把 HEAD 推到 `1c19b1953`（`cmd/wisp/run.go` + `agent_profiles_451*.go`，票 451）⇒ **本读数属于 `bf4adaa7b` 的 `cmd/wisp`，不是现在工作树的** |
| 名册 | 该二进制 `-test.list '.*'` = **669** 枚顶层用例，末枚 `TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`（终点锚） |
| DLL 暂存 | `onnxruntime.dll` / `sherpa-onnx-c-api.dll` / `sherpa-onnx-cxx-api.dll` 三枚拷到二进制旁，SHA256 与 `third_party/sherpa-onnx/` **逐枚同串** |

---

## 3. 读数（canonical run，`Suffix=c`）

### 3.1 退出码（真进程退出码，`Start-Process -PassThru` + `$p.ExitCode`，不经 `go test`）

```
EXIT_CODE_DEC=-1073741819
EXIT_CODE_HEX=0xC0000005
EXIT_CODE_U32=3221225477
```

- `0xC0000005` = **STATUS_ACCESS_VIOLATION**（native 异常）。按 MSYS2/bash 的翻译＝`128+11`＝**rc 139**，与票面 2026-10-13 的 `rc=139`、`314c` 的 `rc=139` 同一码位。
- 与 `-1073741819` 逐位同值（`0xC0000005` 当有符号 int32 读）。

### 3.2 四数（**死时**的四数，不是全量）

```
RUN=480  PASS=476  FAIL=3  SKIP=0
```

口径（逐条可核）：`RUN`＝stdout 里全部 `=== RUN ` 行（顶层 348 ＋ 子测试 132）；`PASS/FAIL/SKIP`＝全部 `--- PASS:/FAIL:/SKIP:` 行（顶层 `344/3/0` ＋ 子测试 `132/0/0`）。
**480 − (476+3+0) = 1** ⇒ **恰好 1 枚跑过却没有 verdict**，就是死点那枚。**终点锚 `TestZZ315A1*` 命中 0**（没跑到；分母 669 只走了 348 枚顶层）。

3 枚红（全部与死法无关，均可点名）：

| 用例 | 原文摘句 | 性质 |
|---|---|---|
| `TestSecondLaunchActivationRaisesThePanel319` | `the prior window 0x116F0552 did not take the foreground (the foreground is 0x5701544) … this is a failed measurement, not a product verdict` | 前台被 `magpie`(0x5701544) 占着 |
| `TestSecondLaunchActivationOnAVisiblePanelDoesNotHideIt319` | 同上（同一对面 `0x116F0552` / `0x5701544`） | 同上 |
| `TestAC4FocusReturnToPriorWindowGap33r5` | `AC#2 appendix (head 4ce7cdf81 …)` / `the panel did not take the foreground on Show` | 本机既知红（11-14 干净全量里唯一那枚也是它） |

### 3.3 `panic:` / `fatal error:` 扫描（两个流都扫）

```
stdout: panic:=0  fatal error:=0  (SIGSEGV|signal|exit status)=0
stderr: panic:=0  fatal error:=0  (SIGSEGV|signal|exit status|runtime:)=0
```

⇒ **没有 Go 层任何东西被打印出来**；进程是被操作系统直接带走的（native 未处理异常）。

### 3.4 最后 40 行 stdout（全份：`probes/314/ac0/run1-stdout-tail40.txt`；这里给末尾 12 行）

```
--- PASS: TestPanelThreadIsSTAAndExitsCleanly (1.35s)
=== RUN   TestPanelThreadNameIsNotInResidentRoster
--- PASS: TestPanelThreadNameIsNotInResidentRoster (0.00s)
=== RUN   TestBallPanelGesturesReachThePanelThread
wisp: panel window is up (panel-hotkey, cold 1206.5 ms, hot path 0.0 ms)
wisp: panel window hidden (focus returned to hwnd 0x5701544)
--- PASS: TestBallPanelGesturesReachThePanelThread (1.40s)
=== RUN   TestBallGestureWithoutPanelHostStillRecords
wisp: ball panel-hotkey: no executor was handed … (what is left is recorded by name, never invented)
--- PASS: TestBallGestureWithoutPanelHostStillRecords (0.00s)
=== RUN   TestAC4PriorFocusSurvivesARefusedPanelSample          ← 最后一条 === RUN
wisp: panel window is up (test-harness, cold 1096.2 ms, hot path 0.0 ms)
wisp: panel window hidden (focus returned to hwnd 0x5701544)    ← 最后一行 stdout（12:58:20.461）
```

**该用例之后再无一行输出**（`stdout-c.log` 共 1325 行，RUN 行在 1323，其后仅 2 行），**没有 `--- PASS/FAIL/SKIP`，也没有 `PANIC`/`FAIL <pkg>` 收尾行**。

### 3.5 stderr 尾部（全量 471,116 B / 2721 行在读数根；这里给末 6 行）

```
time=2026-11-15T12:58:20.517+08:00 level=INFO msg="MARK299 wndproc hwnd=0x7c50ab0 msg=0xc0fa known=true"
time=2026-11-15T12:58:20.520+08:00 level=INFO msg="panel thread ending" why="the library pump returned" window_opened=true
time=2026-11-15T12:58:20.521+08:00 level=INFO msg="panel thread exited cleanly" shows=1 toggles=0 disposals=0
```

⇒ 死前最后一句是**面板线程干净退出**（`disposals=0`），**不是**建窗途中。进程在最后一行 stderr 之后约 **0.5 s 内**消失（最后写盘 12:58:20.522，`Wait-Process` 在 12:58:21 返回）。

### 3.6 采样器（20 s 一档，`probes/314/ac0/run1-sampler.csv`）

| tick | 时刻 | 活着 | CPU(s) | WS(MB) | stdout B | webview2 | 竞争者 |
|---|---|---|---|---|---|---|---|
| 28 | 12:58:16 | 是 | 70.33 | 115.8 | 156,116 | 92 | 0 |
| 29 | 12:58:37 | **否** | — | — | 157,725（+1,609） | **37**（−55） | 0 |

⇒ 不是"卡住后被外部杀"：CPU 与日志都在推进，下一档就已经没了；55 个 `msedgewebview2` 随进程一起走。

---

## 4. 归因（AC#0 要的三选一）

| 候选 | 判 | 依据 |
|---|---|---|
| **(i) native 异常码** | **是** | 退出码 `0xC0000005` ＝ STATUS_ACCESS_VIOLATION；零 `panic:`/`fatal error:`；零 verdict；终点锚零命中 |
| (ii) `os.Exit(n)` | 否 | 若走 `os.Exit(1)`，测试框架会先打 `--- FAIL: <该枚>` + `FAIL <pkg> <秒>s`；两者都没有。且 run A（真走 `os.Exit(1)`）的退出码是 `0x00000001`，形状完全不同 |
| (iii) `-test.timeout` panic | 否 | 只烧 **554.5 s**，`-test.timeout=90m`（5400 s）未到期；零 `panic: test timed out`、零 goroutine dump |
| 「静默退出」 | 否（就码位而言） | 进程不是"安静地正常退出"（那是 rc=0）；它是**不带 Go 层信号的 native 崩溃**——"静默"只指**没有 panic 输出**这一面 |

**归因一句话**：`TestAC4PriorFocusSurvivesARefusedPanelSample` 收尾段（面板线程已干净退出之后、verdict 写出之前）整进程 native ACCESS_VIOLATION（`0xC0000005` / rc 139），不是 Go panic、不是超时、不是 `os.Exit`。

### 4.1 具名上游（票面要的"死在哪一枚"）

- 最后一条 `=== RUN`：**`TestAC4PriorFocusSurvivesARefusedPanelSample`**（`cmd/wisp/panel_resident_windows_test.go:1203`，该文件最后一枚用例，函数体到 `:1252`）
- 该枚**没有** verdict；其后没有任何 `=== RUN`
- 时序细节：`Show` 之后 `mgr.Hide()` 打了 `panel window hidden`，stderr 记 `panel thread exited cleanly`；**死点在"测试体+cleanup 已跑完、testing 还没写出裁决行"的那一段**（Go 的 `tRunner` 是先跑函数、再跑 cleanup、最后才写 `--- PASS`）

### 4.2 ⚠ 与票面 2026-10-13 旧签名的同与异（⛔ 不能混说）

- **同**：无 Go panic、无 `--- FAIL`、无 `exit status` 行、整进程消失；都在"真窗测试区"；退出码族相同（旧账 `rc=139` ≡ 今天的 `0xC0000005`）。
- **异（新读数）**：票面旧账的死点在 **`CreateWindowExW` 返回到首几条 wndproc 之间**；**今天这发死在面板线程"三个 disposals=0 的干净退出"之后**，是**收尾段**而不是建窗段。⇒ 同一死法（native 裸崩），**位置相邻但不同**；因此今天这发**不能**直接拿去给旧账的死点背书，只能说"该族在 `bf4adaa7b` 上仍在，且落点漂到 AC4 收尾"。

### 4.3 故障模块：本机**没留档**（诚实边界）

- `Get-WinEvent -LogName Application` 近 45 min：**0 条**与本二进制相关的 `Application Error`；
- `%LOCALAPPDATA%\CrashDumps`：**空**；
- `WER\ReportArchive` / `ReportQueue`：本时刻**无**新档。
- 盘上唯一的 `wisp.test.exe` WER 旧档：**2026-10-22 03:11:47Z**，`EventType=APPCRASH`、`Sig[6] Exception Code=c0000005`、`Sig[3] Fault Module Name=wisp.test.exe`、`Exception Offset=0x95393`、`AppPath=C:\Users\swq\go-build330552677\b001\wisp.test.exe` ⇒ **仅历史旁证**（旧构建），⛔ 不能当今天这发的模块归因。**要有模块名，得在崩溃路径上开 LocalDumps/WER 或挂调试器——那是另一发实验。**

---

## 5. 并发与桌面边界（读数干净到什么程度）

- **独占闸**：`gate-and-run.ps1` 等到 `wisp.test.exe`/`wisp-winlive.test.exe` **连续 60 s 不存在**才发车（票 314 §3 / 台账 `A1860`：两发整包并发会作废读数）。
- **起跑前**：竞争者 0；**起跑后**采样逐档记竞争者数：12:49–12:58 期间仅 **12:51:29 一档为 1**（另有一条腿的单枚定向跑，<20 s），其余全 0。⇒ 这发**不是**干净到"零竞争者"，但**不是**两发整包并发；`0xC0000005` 出现在 12:58:20，距那档 7 分钟。
- **桌面**（候选 (a) 的代理变量，`desktop-state.ps1`）：`input_desktop=Default`（**未锁屏**）、session 1 `LogonUI=0`、前台 `magpie`(pid 57932)、CPU 71%。⇒ 锁屏桌面**不在场**，既不能坐实也不能证伪候选 (a)。
- **驻留球**：`wisp.exe` PID **49348** 全程在跑（未动它）；它占着 `Ctrl+Alt+Q/M/P` 三枚全局热键 ⇒ 上面 2 枚热键红是它的直接后果（票 314 与票 450 都记过这条）。
- **跑了两次同形的说明**：第 1 发（`Suffix=c`）就是上面的死亡读数；第 2 发（`Suffix=d`）是同形复验，读数见 §6。

---

## 6. 同形复验（`Suffix=d`）：**没复现**，跑完了

同一二进制、同一 cwd、同一 env、同一桌面（`Default`）、闸后独占（起跑前竞争者 0；全程采样竞争者 0）：

| 项 | 值 |
|---|---|
| 起止 / 耗时 | 13:01:12 → 13:24:22，**1389.8 s**（`-test.timeout=90m` 未到期） |
| **退出码** | **`0x00000001`（dec 1）**——`os.Exit(1)`：测试红，**不是** native 异常 |
| 四数（全行口径） | **RUN=892 / PASS=861 / FAIL=11 / SKIP=20**（顶层 `669/638/11/20`，子测试 RUN=PASS=223） |
| 终点锚 | **命中 1 次**（`--- PASS: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask (0.25s)`）⇒ **669 枚顶层用例全有 verdict，零 unverdicted** |
| `panic:` / `fatal error:` | 0 / 0（两流） |
| 死点那枚 | `TestAC4PriorFocusSurvivesARefusedPanelSample` **PASS 1.30 s**（与第 1 发同点、同二进制、同桌面） |

11 枚红**全部是本机既知噪声族**，逐枚可点：

| 族 | 用例 | 依据 |
|---|---|---|
| 前台被 `magpie` 占（3） | `TestSecondLaunchActivationRaisesThePanel319`、`TestSecondLaunchActivationOnAVisiblePanelDoesNotHideIt319`、`TestAC4FocusReturnToPriorWindowGap33r5` | 与第 1 发同句（`did not take the foreground … failed measurement`） |
| 驻留球（PID 49348）占全局热键（8） | `TestLive246*`×4、`TestLive228*`×2、`TestLive258ResidentOccupiedCombinationNamesTheNewValue`、`TestAC246DevLegIgnoresTheTestTaskInjection`、`TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero` | 原文逐字：`another program on this machine already owns mute=…, panel=Ctrl+Alt+P (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED)`／`SKIP-LOUD: … the premise this case judges is missing, not the product` |

**⇒ 复现率：2 发里 1 发死。** 死法是**条件性**的（与票面"5 死 4 活"的同一条性质一致），⛔ 不是"每个同形发必死"，也 ⛔ 不是"这一枚用例必红"——第 2 发把同一枚跑成 PASS，随后整包跑完。

---

## 7. 本件**不做**什么（边界）

- ⛔ 不改产码/测试/断言/票面/阈值/golden；本腿在仓内只新增本件与 `probes/314/ac0/**`（`git status` 里 `cmd/**`、`internal/**` 零改动）。
- ⛔ 不归因到具体模块/行：没有 dump/WER，只有退出码与死点；**"哪一行访问违例"这件事本件给不出**。
- ⛔ 不当 AC#1/#2 的读数：只跑了"同形一发/两发"，没做 fresh-profile A/B，也没跑锁屏桌面。
- ⛔ 不判断"是否 299/408 回归"：本件只交死法取证。

## 8. 撤销口令

删 `docs/evidence/s1/314-ac0-forensics.md` 与 `.scratch/wisp/probes/314/ac0/**` 即复原；读数根 `D:\tmp\wisp\20261115\w314\` 可原样留着复跑。
