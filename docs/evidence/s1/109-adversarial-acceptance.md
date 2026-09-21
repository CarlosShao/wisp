# 票 109 —— 对抗验收裁决表（模型交还点复验 / `VerifyInstalled` + `bridge.go:58` 接线）

验收方：`acceptor-ticket109b`（**接续**：前一名 `acceptor-ticket109` 跑了 22 次工具、**一份证据文件都没落盘**，
只有 `/tmp/ac109-w1/{pre,fix}` 那枚快照留着 ⇒ 本文件是**从零重跑**，非续写其结论）。
时间基准：`date` = **2026-09-21 20:31:04 CST 2026 (Mon)**（本机 Git Bash，非转述）。
仓库：`D:\work\workspace\projects plans\Wisp`，分支 `dev`，**共享工作树很脏**（票 113 在写 `internal/winsec/winsec_other.go`）。
本文件是本次验收**唯一**落盘产物；**一行生产码都没改**（见文末「工作树自证」）。

- 票面：`.scratch/wisp/issues/109-models-have-a-cross-account-write-window-after-ensure.md`（`ready-for-review`，5 格框全**未勾**，本验收不改票面/不勾框/不动台账）
- 锚定 sha：`bdde553`（修前红）+ `f6818f2`（修复）
- 复算/变异一律在**仓外**快照：`git archive <sha> | tar -x -C /tmp/ac109b-qz{,-pre}`（A38④：**仓内不建 worktree、不 checkout**）

## 三档证据口径（逐格标注，不含糊）

- **〔独立复现〕** 我亲手跑过、读数来自我这次执行的输出。
- **〔日志＋归档，我抽验〕** 实现方日志/票面有读数，我只抽查了部分内容或只确认了产物存在。
- **〔仅自述，不背书〕** 我没跑成的原因**逐格写清**。

## 裁决表（本文件是**渐进写**：每验完一格立刻追加，不攒到最后）

| 格 | 主张（攻，不复述） | 我的裁决 | 证据档 |
|---|---|---|---|
| AC#1 | 窗口真实形状（SID 级 `icacls` 前后 + 交还点→读取点 file:line 链 + 下次 `Ensure` 是否重验） | **通过**（三问全量到；"读取点"这一端在本仓不存在，链到 `Dispatch` 即断） | 〔独立复现〕 |
| AC#2 | 收口方向②，**正半边**（窗口关上）＋**反半边**（正常安装/多实例复用不被打断） | **不通过**（两半在包内都成立，但**生产调用者为 0** ⇒ 按第 8 类缺陷判；断点＝`bridge.go:39`） | 〔独立复现〕 |
| AC#3 | `rename`-into-place 序列的 seam 守卫 | **通过**（M2 独立把它跑红 ⇒ 钉咬在"再验"上不是恒真） | 〔独立复现〕 |
| AC#4 | 安装目录那一格的反向钉子，三条界 | **通过，但界线欠覆盖**（第四形滑过：R-109-2 / R-109-3，均已实测） | 〔独立复现〕 |
| AC#5 | 修前红 + M1/M2 + 门禁四数 + gofmt/gofumpt + d22scan | **通过**（修前红逐字复现；四数 RUN66/PASS62/FAIL0/SKIP4 两枚快照各跑一次同读数；SKIP 是既有联网门；d22scan HEAD 复算 rc=0、ban#8 internal/=371 不降） | 〔独立复现〕 |
| 攻1 | **守卫之后还剩不还剩窗** | **还剩一整段**（PROBE1 红；但"引擎装载"今日无生产路径可达成 ⇒ 不据此判 FAIL，立案 R-109-1） | 〔独立复现〕 |
| 攻2 | `bridge.go:58` 是不是真生产调用点 | **不是**（四条仪器；`internal/models` 不在 `cmd/wisp` 依赖图内） | 〔独立复现〕 |
| 攻3 | AC#4 第四形 | **找到三形**：未点名文件/目录（R-109-2）、按账户名过滤（R-109-3）、`(M,DC)`（R-109-5） | 〔独立复现〕 |
| 攻4 | 复验的代价 | 一次粗测：256 MiB ⇒ 1.07–1.22 s（220–252 MB/s），读过的字节数＝pin 的字节数；真实 manifest 全 6 枚 707.8 MB ⇒ 交还复验约 3 s 量级，**与 `Ensure` 已付的那一份同量级＝启动期模型读翻倍**。⚠ 不对阈值作判定 | 〔独立复现〕 |

## 总判

**FAIL-退回。** 一句话理由：票 109 的收口（`VerifyInstalled` + `bridge.go:58`）**在包内确实咬住了它声称要咬的那一段窗**
（M1/M2/修前红/POSIX 全部我亲手复现），但**那条路在生产里调用者为 0**——守卫挂在一个装配根根本到不了的函数上，
且**守卫返回之后到读取者 open 之前还剩同样形状的一整段**（我的探针已红）——两件事同一个根因：**交还点的下游还不存在**，
所以本票的 AC#2 判据（窗口关上）**目前不可能成立**。
按派单口径"调用者为 0 或只有测试调到 ⇒ AC#2 不通过"，我判 FAIL-退回，不写"通过附条件"。

**退回该改什么（最小集合，别把已经量实的东西砸掉）**：AC#1/AC#3/AC#4/AC#5 四格我全判通过，
修前红、变异、门禁、POSIX、SID 级前后读数**都是真的**，实现方自述里没有一条我在场的读数是假的；
需要动的是**AC#2 的落点与口径**（把守卫挪到 `Ensure` 的返回路径上＝攻 4 问的 R-109-4 一起解，或明确把本票范围写成
"预置守卫，待票 4 接线时启用"并把"没人调"做成一条红钉），外加把 R-109-2/R-109-3 两条界线的洞补进 AC#4 的用例。

---

## 追加区（逐格结论，按验证顺序 append-only）

### 攻2（优先级 2）：`bridge.go:58` 是不是真生产调用点 —— **不是。调用者为 0（只有测试调到）。** 〔独立复现〕

`date` = 2026-09-21 20:3x（每次读数前现取）。四条独立仪器，全部我亲手跑：

1. `grep -rn "WireDownloading" --include=*.go .` ⇒ 5 处：`bridge.go:32`/`bridge.go:33`（定义与注释）+
   **3 处全在 `_test.go`**（`bridge_test.go:20`、`handoff_reuse_109_test.go:67`、`handoff_window_109_windows_test.go:119`）。
   `DownloadingBridge` 的引用同理只有 `bridge.go` 自身 + `bridge_test.go`。
2. `grep -rn "CarlosShao/wisp/internal/models" --include=*.go .` ⇒ **全仓只有 1 个 importer：`tools/signmodels/main.go`**。
   而 `tools/signmodels` 只用了 `MinisignPublicKey` / `SignPayload` / `VerifyMinisignSignature` / `ParseMinisignPublicKey` /
   `DeriveKeyID`（逐条 grep `models\.` 打印核对）——**签名侧的构建期工具**，一次都没调 `NewManager`/`Ensure`/`VerifyInstalled`/`WireDownloading`
   （`grep -rn "Ensure\|VerifyInstalled\|WireDownloading" tools/ cmd/ --include=*.go` ⇒ **0 行输出**）。
3. `go list -deps ./cmd/wisp`（装配根，`cmd/wisp/run.go` 所在）⇒ 22 个 `wisp/internal/*` 依赖里
   **没有 `wisp/internal/models`**（有 winsec/secret/observe/risk/config/llm/plugin/memory/agent/tools/…/proc，见下方原始清单）。
4. `grep -E "speech|models|ball|statemachine"` 套在 `go list -deps ./cmd/wisp` 上 ⇒ **0 行**（三个候选都无）。

**断点指名**：链不是"断在 `bridge.go` 内部某一行"，而是**从装配根出发根本不存在那条边**——
`internal/models/bridge.go:39 func (b *DownloadingBridge) Run` 是这条链的最高可见入口，
它在 `cmd/` 下的调用者数 = **0**；要落到 `bridge.go:58`，先要有谁调 `:39`。
顺带一条同源的硬事实：**"读取者"在本仓还不存在**——`internal/speech/` 只有 `doc.go` 一个文件
（`ls internal/speech/` ⇒ `doc.go`），`internal/engines/` **目录都不存在**（`ls internal/` 无此项），
`grep -rn "\.onnx\|\.gguf\|sherpa" --include=*.go internal/ cmd/ tools/` 的非测试命中只有
`cmd/wisp/doctor.go` 那几行 sherpa 链接自检和 `internal/config/catalog.go` 的 provider 名字，
**没有任何一处按返回的目录去 open 模型文件**。

⇒ 票面 AC#2 的两半（正半边"窗口关上"、反半边"复用不被打断"）**目前只在包内测试面成立**。
本仓第 8 类缺陷（"测试证明它会工作，生产里没人叫它"）**第五起，形状完全一致**：
新守卫落在一条没有生产调用者的路上。

### 攻1（优先级 1）：守卫之后还剩不还剩窗 —— **还剩一整段，就在 `bridge.go:58` 之上那一层。** 〔独立复现〕

我在 `/tmp/ac109b-qz/internal/models/zz_probe109b_test.go` 写了三枚探针（**仓外快照**，非变异、非门禁），
`go test -count=1 -v -run 'PROBE109B' ./internal/models/` ⇒ rc=1，读数原文：

- **PROBE1 FAIL（探针红＝残留窗成立）**：
  `PROBE1 residual window: the hand-off (bridge.Run) returned nil and announced the model available, then the file was
  swapped and the reader opened 1024 bytes that differ from the 1024 the guard hashed - no code between the guard's
  return at bridge.go:58 and this open notices.`
  形状：走**真实** `WireDownloading(...).Run()` ⇒ `Ensure` → `VerifyInstalled`（守卫，返回 nil）→ `Dispatch(EvDownloadCompleted)`
  → `Run` 返回 nil（**已宣布可用**）→ **此后**换 1 字节 → 读取者 `os.ReadFile(target)` 拿到换过的字节。
- **PROBE1 控制项（同一跑内）**：换完之后**再叫一次** `VerifyInstalled` 当场抓住，错误原文带 pin：
  `model: hand-off re-verification: model kws-fixture installed files no longer match the signed manifest in ...\kws-fixture:
  model.onnx: sha256 mismatch (pins 785b0751fc2c..., got dd7bafd7c0c4...)`
  ⇒ **守卫本身不坏，坏的是"这之后没有任何人再叫它"** ——残留窗是"未接线"问题，不是"算法不成立"问题。
  （`pins 785b0751fc2c` 与实现方自述的那串**逐字对上**，我这里是**我自己跑出来的**，不是照抄。）
- **PROBE3 PASS**：换字节后再叫一次 `Ensure` ⇒ 树被修回 pin，`harm ceiling is one session`
  ——AC#1 第三问（"下一次 `Ensure` 会不会重新验签"）**我自己量到答案是会**，票面这条主张**成立**。

**逐行读码把边界写清（不许含糊成"理论上还剩一点"）**：`bridge.go:39-68` 是这条链的全部可见代码，
`VerifyInstalled` 在 `:58`；`:59-62` 是 `Dispatch` + `return`；`:64` `Dispatch(EvDownloadCompleted, nil)`；`:67` `return b.record(), nil`。
- `b.record()`（`:80-86`）= `machine.State()` + append 到 `b.log`，**不碰磁盘**；
- `internal/statemachine/` 非测试文件里 **`grep '"os"|os.Open|os.Stat|filepath'` ⇒ 0 命中**（`ls` 显示只有 doc.go/events.go/machine.go/states.go/table.go/timeouts.go）
  ⇒ `Dispatch` 是纯内存态转换。
⇒ **结论：守卫是这条交还路径上最后一个碰字节的操作**，它返回之后到"读取者 open"之间**确实有可插入点**（PROBE1 就是插进去的那一刀），
但"读取者"在本仓**还不存在**（见攻2 的 `internal/speech/doc.go`、无 `internal/engines/`）。
所以 AC 声称要防的**那个字面结局**（"引擎装载了被换过的模型"）**今天没有任何生产代码路径能达成**——
**这一格不据此判 FAIL**；但它与本票 AC#1 的结局同类、只隔一层，**必须立案 R-109-1**（见文末），
否则票 4 一接线，这张票的收口立刻退化成"上次恰好没人读"。

### 攻3（优先级 3）：AC#4 那条界线的第四形 —— **找到两形，都滑过去了。** 〔独立复现〕

- **R-109-2（PROBE2 FAIL，两半各红一次）**：`manifest pins exactly these paths for kws-fixture: [model.onnx tokens.txt dict/inner.txt README.md]`，
  然后往安装目录里放一个 manifest 不点名的 `attacker.onnx` ⇒ `VerifyInstalled` **返回 nil**；
  再放**一整个 extra 目录** `sub/w.onnx` ⇒ 同样 nil。原文：
  `PROBE2 slip-through: VerifyInstalled returned nil with an extra unhashed file sitting in the install directory ...
  VerifyDir iterates entry.InstalledFiles() only (downloader.go:557)` +
  `PROBE2 slip-through (subdir form): an entire extra directory inside the install dir passes the hand-off guard`
  ⇒ AC#4 三条界量的是**"谁持有权利"**，一条都没量**"这个目录里有多少字节被钉过"**：
  持有 AC#4 第 1 条**明写允许**的那份继承写权的人，可以在被守卫判为"干净"的目录里放任意多个**从未被任何门算过 sha256** 的文件。
  可达性：`seededStore` 那种形状下就是"另一个本地账户"，写权是票 95 要保留的，**不需要破任何 ACL**。
- **R-109-3（读码 + 现有断言的口径缺口，未做 Windows 实测）**：`aceLinesNaming(t, path, "BUILTIN\Users")`
  （`handoff_window_109_windows_test.go:35-50`）**按账户名过滤 icacls 行**，而 AC#4 的第 2 条界
  （"安装目录自己 DACL 上出现非 `(I)` 的外来 ACE＝破"）说的是"**外来** ACE"。
  两条界只看 `BUILTIN\Users` 一个名字 ⇒ 一块**显式**（非 `(I)`）的 `Everyone` / `NT AUTHORITY\INTERACTIVE` / 别的用户 SID 授权
  落在安装目录上，`aceLinesNaming(..., "BUILTIN\Users")` 根本看不见它，三条界全读"允许"。
  ⚠ 这一形我**没有跑成实测**（要在快照里另写一枚 Windows 注入用例 + `icacls /grant`，本轮时间预算给了 M1/M2 与门禁；
  登记为**读码结论**，档＝〔仅自述（我自己），不背书〕，下一名验收或实现方补测）。

### 变异清单（我独立重跑 M1/M2，**没照抄红名**；每发先证落地）

锚定与"每发重抽"：两枚都从 `git archive f6818f2 | tar -x -C /tmp/ac109b-m1` / `... -C /tmp/ac109b-m2` **重抽**（仓外，A38④），
改 → **同一条 `&&` 链里 `grep -n` 打印被改后整行** → `go build ./internal/models/` rc=0 → 才读结果。

| 枚 | 锚点（file:line） | 落地证明（grep 打出的那一行原文） | build | 红的测试名（`-count=1 -v -run 'TestAC1\|TestAC2\|TestAC3Wide\|TestAC4'`） |
|---|---|---|---|---|
| M1 退回交还点那一行 | `internal/models/bridge.go:58` | `58:\tif verr := error(nil); verr != nil {` | rc=0 | FAIL `TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified`、FAIL `TestAC1TheWriteWindowIsAnInheritedCrossAccountRight`；PASS `TestAC2NormalInstallAndSecondInstanceReuseAreNotBroken`、`TestAC3WideTemp…`、`TestAC4InstallDirWidth…` |
| M2 把再验改成空守卫 | `internal/models/downloader.go:235` | `235:\tif err := m.VerifyDir(entry, dir); false && err != nil {` | rc=0 | 上面两条 AC#1 仍 FAIL，**另加** FAIL `TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor`；AC#2 复用与 AC#4 仍 PASS |

两枚的"该红哪两格 / 不该红哪几格"与实现方自述**逐字对得上**（红名、多咬的那一枚、AC#4 不受影响），
⇒ 这两枚变异我判**〔独立复现〕、成立**，AC#5 的变异那半**通过**。
⚠ M2 有一条我自己踩到的坑要登记：`if err := m.VerifyDir(entry, dir); err != nil {` 这行在 `downloader.go` **出现两次**
（`:161` local_override 分支 / `:235` `VerifyInstalled`），我第一次 `assert count==1` 直接炸了（rc=1、没跑成）⇒
**改成按"所在函数是 `VerifyInstalled`"选锚**才落到位。若谁照自述里"`downloader.go:235`"字面 copy 一行 sed，会打歪或打两下。

### 修前红（AC#5 第一问，锚 `bdde553`）〔独立复现〕

`/tmp/ac109b-pre`（`git archive bdde553`，仓外重抽）⇒ `go build ./internal/models/` rc=0，
`go test -count=1 -v -run 'TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified' ./internal/models/` **rc=1**，断言原文逐字复现：

```
    handoff_window_109_test.go:95: swapped 1 byte of model.onnx inside the post-verification window
    handoff_window_109_test.go:105: AC#1/AC#2 (AC95-R1): the hand-off reported a model available whose file was swapped after Ensure verified it - the span between Ensure's return and the reader's open has no guard. state=FirstRun
--- FAIL: TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified (0.05s)
```

红名与失败语**与票面 Progress log 登记的一致**（含 `state=FirstRun` 这一段），不是转述。

### AC#1（窗口真实形状）——**通过**。〔独立复现〕

我在**纯净快照 `/tmp/ac109b-fmt`**（`git archive f6818f2`，不含我的探针文件）跑
`go test -count=1 -v -run 'TestAC1TheWriteWindow|TestAC4InstallDirWidth|TestAC2NormalInstall' ./internal/models/`，
`icacls` 读数原文（**我自己这次跑的**）：

```
store root seeded with: [BUILTIN\Users:(OI)(CI)(M)]
BEFORE the hand-off guard - install dir kws-fixture: [BUILTIN\Users:(I)(OI)(CI)(M)]
BEFORE the hand-off guard - installed file model.onnx: [BUILTIN\Users:(I)(M)]
hand-off refused the swap made through that right: model: hand-off re-verification: ... model.onnx: sha256 mismatch (pins 785b0751fc2c..., got dd7bafd7c0c4...)
AFTER the hand-off guard ran - install dir: [BUILTIN\Users:(I)(OI)(CI)(M)]
AFTER the hand-off guard ran - installed file: [BUILTIN\Users:(I)(M)]
```

- 三问逐答：**谁能写**＝`BUILTIN\Users` 的继承 `(I)(OI)(CI)(M)`（父策略，非本对象显式），量到 SID 级；
  **会不会被读**＝这条链的"读取点"在本仓**不存在**（攻2 已给仪器）。
  ⚠ **票面那条 file:line 链我逐行核对之后是"两对一错一对旧"，已就地更正（原来我写"行号对得上"，不准确）**：
  `downloader.go:138 Ensure` ✔（`bdde553` 与 `f6818f2` 都是 138）、`bridge.go:44 b.mgr.Ensure` ✔（两版都是 44）；
  票面写的 **`downloader.go:178 VerifyDir` ✘ —— 两个版本里 `:178` 都是 `if m.busy[id] {`**，缓存命中重验实际在 **`:171`**；
  票面写的 **`bridge.go:51 Dispatch(EvDownloadCompleted)` 在 `bdde553` 是对的，但被 `f6818f2` 自己那 13 行插入挪到了 `:64`** ⇒ 引用随修复腐坏。
  链本身仍然**到 `Dispatch` 就断**（断点见攻2）。
  **下次 `Ensure` 会不会重验**＝**会**，PROBE3 我自己量到"换字节后下一次 `Ensure` 把树修回 pin"⇒ 危害上限"一次会话内"成立。
- "逐字不变"这句：`before/after` 两对读数**逐字相同**，且该用例内建了相等断言（`windows_test.go:138`）⇒ **安装目录 ACL 没被收窄**这条我背书。

### AC#2（收口方向 + 正反两半）——**不通过**（正半边在包内成立，**反半边同样成立，但两半都只在测试面**；按第 8 类缺陷判）

- 方向选②的**理由本身**我认：①要拿掉 `BUILTIN\Users` 那份继承读权＝推翻票 95 判"不封"的全部理由；
  `internal/winsec` 只有 current-user-only 一种形状（`git show f6818f2 --stat` ⇒ **本 commit 未含 `internal/winsec/` 任何文件**，5 个文件里只有 ticket md + `bridge.go` + `downloader.go` + 两枚测试）。
- 正半边：修前红→修后绿，M1/M2 各自红在该打的两格（见上表），**成立**。
- 反半边：`server hits: install=1, second-instance reuse=0`（**我这次跑出来的原文**，`handoff_reuse_109_test.go:59`）
  + 同 store 上第二个 Manager `Ensure`/`VerifyInstalled`/`Run` 全过 ⇒ **在测试面成立**。
- **判不通过的唯一理由**＝攻2：`DownloadingBridge.Run` 的生产调用者为 **0**，`internal/models` 根本不在 `cmd/wisp` 的依赖图里。
  指名断点：`internal/models/bridge.go:39`（`Run` 的定义）之上**没有边**，`bridge.go:58` 因此不可能被生产触发。

### AC#3（`rename` seam 守卫）——**通过（在包内这一层）**。〔独立复现〕

M2 那一发**独立地**把 `TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor` 跑红（见变异表），
⇒ 这枚钉子确实咬在"再验"上而不是咬在描述符上，不是恒真断言。
`stagingDir` 装完即消失也在这次跑里被 `os.Stat` 断言着（`handoff_reuse_109_test.go:88`）。
⚠ 边界同攻1：这条守卫管的是"rename 之后被换 ⇒ 交还点抓"，**不管**"交还点之后再换"。
POSIX 上同一枚用例我也真跑了（见下方 Docker 段）。

### AC#4（安装目录反向钉子，三条界）——**通过，但界线有洞**（两形滑过去，见攻3 的 R-109-2 / R-109-3）。〔独立复现〕

`TestAC4InstallDirWidthIsDeliberateAndHasABreakLine` 我自己跑出 PASS，`AC#4 nail: install dir kws-fixture inherits
[BUILTIN\Users:(I)(OI)(CI)(M)] (parent policy, guard is the re-verify)`；三条界在代码里都是**可执行的**（不是注释）：
第 1 条由 `len(lines)==0 ⇒ Fatal`、第 2 条由 `!Contains(l,"(I)") ⇒ Errorf`、第 3 条由 `/remove` 后仍留显式外来 ACE ⇒ Fatal 承担。
"没接"与"故意不接"可区分这条：`downloader.go:599` 的 `os.MkdirAll(installDir, 0o755)` 旁边现在有测试点名它 ⇒ 成立。
但攻3 两形说明第 2 条界**按账户名过滤**、且三条界**都不看目录里未被点名的文件**。

### AC#5（修前红 / 变异 / 门禁）——**四数与工具门全部通过**。〔独立复现〕

- **修前红**：见上一节，锚 `bdde553` rc=1，断言原文逐字复现。
- **门禁**（口径：分包 `./internal/models/`、带 `-v`、`-count=2` 不缓存；在**纯净快照 `/tmp/ac109b-qz`** 跑的，**不是我改过的工作树**）：
  `go test -count=2 -v ./internal/models/` ⇒ **rc=0，RUN=66 PASS=62 FAIL=0 SKIP=4**。
  自洽核对：`=== RUN` 66 行 ＝ `sort -u` 出来的**33 个不同测试名 × 2** ✔（不是把 PASS 数当 RUN 数）。
  SKIP 四条＝`TestRealDownloadVadThroughPipeline`×2 + `TestRealDownloadPuncArchiveThroughPipeline`×2，
  Skip 条件是**测试自己的第一行**：`manifest_real_test.go:120`/`:147`
  `real-network spot check; set WISP_IT_REAL_MIRROR=1 to run` ⇒ **既有联网跳过，不是拿跳过冒充通过**；
  且 `git show --stat f6818f2` 里**没有** `manifest_real_test.go`，本票没碰它。
- `gofmt -l internal/ cmd/ tools/`（快照内）⇒ **0 行输出**、rc=0。
- `"$(go env GOPATH)/bin/gofumpt.exe" --version` ⇒ `v0.7.0 (go1.27.1)`，`gofumpt -l internal/models/ internal/winsec/` ⇒ **0 行输出**、rc=0
  （二进制**在本机存在**，票 92 那类谎报我这没有：路径 `C:\Users\swq\go\bin\gofumpt.exe`，7,506,432 字节，2026-09-21 11:54）。
- `go vet ./internal/models/` ⇒ rc=0。
  ⚠ 一条**与派单口径不同**的读数要如实登记：`GOOS=linux go vet ./internal/models/ ./internal/winsec/` 在**我这台机**上跑成 **rc=0**
  （派单写"本机永远 rc=1"）。它仍然**只做类型检查、不执行任何用例**，所以"POSIX 真跑"另算（下面 Docker 段），这条 rc=0 **不能**当"测过了"。
- `sh scripts/d22scan.sh`（**当前 HEAD 复算**：`git rev-parse HEAD` = **`67cfdc4`**，比开工时的 `9628102` 又前进了 3 枚，
  新那枚 `67cfdc4` 是 `docs(reports)` 单文件、`git log f6818f2..HEAD -- internal/models/` **空** ⇒ 被验代码未动）
  ⇒ **rc=0 clean**，正控那步 `runtests.sh: OK - packages=[./...] PASS=21 FAIL=0 SKIP=0`，台账读数
  `bans #1-5 internal/=202, bans #1-5 cmd/=20, ban #6 frontend/=40, ban #7 internal/tools/=18, ban #8 design/=16, ban #8 frontend/=40, ban #8 internal/=371, ban #8 cmd/=26`
  ——**ban #8 internal/=371 与实现方自述的 371 逐字相同**（它开工时是 366 ⇒ 各 scope **不降**）。

### POSIX 那半（Docker 真跑）——**通过**。〔独立复现〕

`golang:1.27` + `CGO_ENABLED=0`，**tar 走 stdin、容器内落盘后先证明文件在**（刻意绕开 Git Bash `-v "C:\…"` 静默挂空那枚坑；
我第一发**确实踩到了"看不到文件"**：`ls: cannot access '/w/go.mod'` ⇒ 加 `MSYS_NO_PATHCONV=1` 重跑才有下文）：

```
=== FILE PROOF ===
-rw-r--r-- 1 197609 197609 883 Sep 21 12:21 /w/go.mod
-rw-r--r-- 1 197609 197609 4796 internal/models/handoff_reuse_109_test.go
-rw-r--r-- 1 197609 197609 4157 internal/models/handoff_window_109_test.go
Linux a98d70464683 6.6.114.1-microsoft-standard-WSL2 ... x86_64 GNU/Linux
--- PASS: TestAC2NormalInstallAndSecondInstanceReuseAreNotBroken (0.03s)   [server hits: install=1, second-instance reuse=0]
--- PASS: TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor (0.02s)
--- PASS: TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified (0.02s)
ok  github.com/CarlosShao/wisp/internal/models  0.077s
```

### 攻3 补：R-109-3 **升级为独立复现**（我把它跑红了）

`/tmp/ac109b-qz/internal/models/zz_probe109b_windows_test.go` ⇒ `go test -count=1 -v -run TestPROBE109BAC4BreakLineIsAccountScoped`
**rc=1**（探针红＝滑过去成立）。读数原文：

```
install dir descriptor after planting an EXPLICIT Everyone:(OI)(CI)(M):
  <dir> Everyone:(OI)(CI)(M)                                    <-- 非继承、外来、带写权
        BUILTIN\Users:(I)(OI)(CI)(M)
        DESKTOP-LVS7839\CodexSandboxUsers:(I)(OI)(CI)(M,DC)
        S-1-5-21-3623186960-...-1717338598:(I)(OI)(CI)(M,DC)
        NT AUTHORITY\SYSTEM:(I)(OI)(CI)(F) ... BUILTIN\Administrators:(I)(OI)(CI)(F) ... DESKTOP-LVS7839\swq:(I)(OI)(CI)(F)
AC#4 as written sees Users=[BUILTIN\Users:(I)(OI)(CI)(M)] (break line 2 reads OK=true)
non-inherited foreign ACEs actually on the object: [Everyone:(OI)(CI)(M)]
R-109-3 slip-through: AC#4 break line 2 reports OK while this object carries an EXPLICIT non-inherited foreign write grant ...
```

⇒ AC#4 第 2 条界**按一个账户名过滤 icacls 文本**，任何别的主体的**显式**（非 `(I)`）授权对它隐形；
三条界在这块 ACE 存在时**全读"允许"**。（对照我上一格那条"未跑成实测"的登记：这一格**已经跑成**，以本节为准。）
同一条 `icacls` 顺带量到本机现实：安装目录带着**两条继承来的 `(M,DC)`**（`DC`＝可改 DACL）
⇒ "别的账户不只是能写，还能把这一格的 ACL 改写"这一形也在，仍属第 1 条界"继承＝允许"的口径内（登记 R-109-5）。

### 攻4：复验的代价（**一次粗测**，不去对任何阈值）〔独立复现〕

`/tmp/ac109b-qz/internal/models/zz_cost109b_test.go`：造一枚 256 MiB 的真实形状已安装文件（sha256/size 按 pin 配平），
**直接叫生产函数** `Manager.VerifyInstalled`：

```
installed file on disk: 268435456 bytes; manifest pins: 268435456 bytes across 1 file(s)
VerifyInstalled call#1 (page cache partly warm): 1.1753855s  =>  228 MB/s over 256 MiB
VerifyInstalled call#2 (warm):                   1.2177493s  =>  220 MB/s
VerifyInstalled call#3 (warm):                   1.0668591s  =>  252 MB/s
VerifyDir (Ensure 缓存命中分支每次付的那一份) sample 1: 1.1438545s / sample 2: 1.1521295s
```

- **读过的字节数**：一次调用＝**逐字节读完该模型全部 pin 文件**（`downloader.go:557` 遍历 `InstalledFiles()`、
  `:583` `io.Copy(h, f)`），实测被读 268,435,456 B 与 pin 的 `size_bytes` **完全相等** ⇒ 没有"只读个头"的侥幸。
- **真实 manifest 外推**（`models/manifest.json` 按 `InstalledFiles()` 同语义算，我自己解析）：
  kws 5.0 MB／vad 0.6 MB／asr-streaming 237.2 MB／asr-offline 239.5 MB／punc 79.7 MB／tts 145.7 MB ⇒
  **全 6 枚合计 707.8 MB**；按上面 220–250 MB/s ⇒ **单次交还复验 ≈ 0.02s(kws) … 1.1s(两个 asr) … 0.7s(tts)**，
  全部模型各交还一次 ≈ **2.9–3.2 s 的额外纯读**（每个模型的交还只付自己那一份，不是一起付）。
  ⇒ 量级判断：**这份成本与 `Ensure` 缓存命中分支已经在付的那一份**（实测 `VerifyDir` 1.14–1.15 s／256 MiB）**同量级**，
  所以本票把"启动期模型校验"从 ~707.8 MB 读到 **~2×707.8 MB**。owner 关心的是"翻倍"这件事，不是单次几秒。
- ⚠ **噪声来源如实标注**：这台机 20:3x–20:5x 有 ≥3 个代理在编译/跑测，且 Windows 上我无法丢弃页缓存
  （需要特权）⇒ 三发都是"偏热"读数，冷启动只会更慢不会更快；**本报告不对 0.5%/25MB 那两条阈值作任何判定**，
  也没因数字难看调过任何阈值。

### "有没有配置能绕过它"（攻2 的第三问，逐条答）

- **没有开关能关掉守卫**：`NewManager` 在 `VerifySignature=false` 时**直接拒绝构造**（`downloader.go:104-107`，C29 双保险），
  `LocalOverride` 分支 `Ensure:161` 也走 `VerifyDir`，`VerifyInstalled:232` 同样认 override ⇒ 这两条**没有旁路**。
- **但守卫是"调用点自带"的，不是" choke point 强制"的**：`grep .Ensure(` ⇒
  生产码里唯一一处是 `bridge.go:44`，测试里 25 处；`VerifyInstalled` 生产码里唯一一处是 `bridge.go:58`。
  ⇒ **任何绕过 bridge 直接用 `Manager.Ensure` 的调用者（面板、引擎、CLI 子命令都算）拿到的仍是"只有旧验签"的交还**，
  守卫不会跟着走。登记 **R-109-4**（守卫该放在 `Ensure` 的返回路径上、或把 `Ensure` 收成非导出/加一层强制，
  那是实现方的选择，不是我这格判据）。

### 四种假绿，逐条点名（我这轮的自查）

1. **跳过冒充通过**：无。SKIP 四条我核到了 `Skip` 条件那一行原文（`manifest_real_test.go:120/:147`，环境变量门），
   且证明本票没碰那个文件（`git show --stat f6818f2` 的 5 个文件里没有它）。
2. **断言恒真**：无。M1/M2 两发变异各自把该红的用例跑红（若断言恒真，变异后仍会绿）。
   ⚠ 但我抓到**我自己**写过一枚恒真形状的探针（R-109-3 第一发里那个 `(!Contains(...) == false)` 的糊涂条件让探测器永远返回空集 ⇒
   探针"绿"是假的）；我发现后重写、这次它红了，以重写后的读数为准。**登记在此，不掩盖第一发的失真。**
3. **跑错对象/错文件**：无。所有复算在 `git archive f6818f2` / `bdde553` 的**仓外**快照；
   门禁跑了两枚快照（`/tmp/ac109b-qz` 与不含我方任何文件的 `/tmp/ac109b-fmt`）⇒ 同一读数 RUN66/PASS62/FAIL0/SKIP4，
   排除"是我的探针文件混进了门禁"这一种可能。
4. **门禁压根没跑**：无。`GATE rc=0` / `rc=0` 与 `ok github.com/CarlosShao/wisp/internal/models 5.372s`（及 2.572s 那一发）都是原始尾行。

### R-109-x（交回立案）

| 号 | 内容 | 属哪个包地界 | 阻不阻塞本票结案 |
|---|---|---|---|
| R-109-1 | 守卫返回之后到读取者 open 之前**同样形状的窗还在**（PROBE1 实测红）。今天不成灾**只因为读取者不存在**；`internal/speech/`=doc.go、无 `internal/engines/` ⇒ **票 4 接线之日即复活** | 现在 `internal/models/`；接线后归 **Phase 4  engines/`internal/speech`** 地界 | **不阻塞**（但必须与"接线票"成对，别单独立着） |
| R-109-2 | `VerifyDir` 只遍历 `entry.InstalledFiles()`（`downloader.go:557`）⇒ 安装目录里**未被 manifest 点名的文件/子目录**从未被任何门算过哈希（PROBE2 实测红，两形） | `internal/models/`（manifest 契约＋交还点） | **不阻塞**，但 AC#4 的"宽到什么程度算破"在字节面上其实**没量**，建议本票补一枚用例再结 |
| R-109-3 | AC#4 第 2 条界**按账户名过滤 icacls 文本** ⇒ 显式非 `(I)` 的 `Everyone`/别的主体授权对它隐形（PROBE 实测红：`Everyone:(OI)(CI)(M)` 存在时三条界全读 OK）。对照 `internal/winsec` 自己写的"Names are never part of a judgment here / 判据用 SID" ⇒ **同一仓两套口径** | 判据实现落 `internal/models/` 测试；**"该用什么口径"属 `internal/winsec/` 地界（票 106/108/112 在飞）⇒ 我只登记不跨界改** | **不阻塞**，但与票 89/104 那条族同源，建议并入下一张 winsec 判据票 |
| R-109-4 | 守卫挂在 `bridge.Run` 这个**旁路**上而非 `Ensure` 的返回路径 ⇒ 任何直接用 `Manager.Ensure` 的调用者拿不到守卫（`Ensure` 是导出的、测试里 25 处调用 vs `VerifyInstalled` 生产 1 处） | `internal/models/` | **阻塞**：与 AC#2 的落点问题是同一件事，退回时要一起解 |
| R-109-5 | 本机实测：安装目录带着**两条继承来的 `(M,DC)`**（`CodexSandboxUsers` 与一条裸 SID）⇒ 同机另一账户不仅可写，还可**改这一格的 ACL**；AC#4 三条界把 `(M,DC)` 归入"继承＝允许"，未讨论"能改 DACL"这一层 | `internal/models/` 的 AC#4 ＋ `internal/winsec/` 的 DACL 判据 | **不阻塞**，登记给 owner 的"封不封"选择题当实测代价用 |
| R-109-6 | 票面 Progress log 里 AC#1 要求引的那条 file:line 链**有一处是错的、一处已被修复自己挪走**：`downloader.go:178 VerifyDir` 在 `bdde553`/`f6818f2` **两版都是 `if m.busy[id] {`**（真身在 `:171`）；`bridge.go:51 Dispatch(EvDownloadCompleted)` 只在 `bdde553` 成立，`f6818f2` 插了 13 行后是 `:64` | 票面文档（`.scratch/wisp/issues/109-*.md`，**我不改票面**） | **不阻塞结案**，但 AC#1 的"引链"这半格要按真实行号重写一遍再结 |

### 我亲手跑过的命令（全部原文，读数见各格）

```
date                                                          # 每节开头现取，本文件时间基准 20:31:04 -> 收尾 20:5x
git status --short / git rev-parse --abbrev-ref HEAD / git log --oneline -12
git archive f6818f2 | tar -x -C /tmp/ac109b-qz ; git archive bdde553 | tar -x -C /tmp/ac109b-qz-pre
git archive f6818f2 | tar -x -C /tmp/ac109b-m1 ; ... -C /tmp/ac109b-m2 ; git archive HEAD | tar -x -C /tmp/ac109b-head ; git archive f6818f2 | tar -x -C /tmp/ac109b-fmt
git show --stat f6818f2 ; git show f6818f2 -- internal/models/bridge.go internal/models/downloader.go
grep -rn "WireDownloading|--include=*.go" ; grep -rn "CarlosShao/wisp/internal/models" ; grep -rn "Ensure\|VerifyInstalled\|WireDownloading" tools/ cmd/
go list -deps ./cmd/wisp | grep -n "wisp/internal" ; go list -deps ./internal/speech ; ls internal/speech/ ; ls internal/
cd /tmp/ac109b-qz && go test -count=2 -v ./internal/models/            # rc=0 RUN=66 PASS=62 FAIL=0 SKIP=4
cd /tmp/ac109b-fmt && go test -count=2 -v ./internal/models/           # 同上读数（无我方文件的对照跑）
cd /tmp/ac109b-pre && go build ./internal/models/ && go test -count=1 -v -run 'TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified' ./internal/models/   # rc=1 修前红
M1: python 改 bridge.go:58 -> grep -n "verr := error(nil); verr != nil" -> go build rc=0 -> go test -count=1 -v -run 'TestAC1|TestAC2|TestAC3Wide|TestAC4'   # rc=1，两枚 AC#1 红
M2: 同上，锚 downloader.go:235（按"所在函数是 VerifyInstalled"选锚）                                     # rc=1，另加 AC#3 红
cd /tmp/ac109b-qz && go test -count=1 -v -run 'PROBE109B' ./internal/models/                            # rc=1：PROBE1/PROBE2 红、PROBE3 绿
cd /tmp/ac109b-qz && go test -count=1 -v -run 'TestPROBE109BAC4BreakLineIsAccountScoped' ./internal/models/  # rc=1（R-109-3）
cd /tmp/ac109b-fmt && go test -count=1 -v -run 'TestAC1TheWriteWindow|TestAC4InstallDirWidth|TestAC2NormalInstall' ./internal/models/   # icacls 前后 + server hits
cd /tmp/ac109b-qz && go test -count=1 -v -run 'TestPROBE109BCostOfOneVerifyInstalled' ./internal/models/                                 # 攻4 粗测
cd /tmp/ac109b-fmt && gofmt -l internal/ cmd/ tools/ ; "$(go env GOPATH)/bin/gofumpt.exe" --version ; gofumpt -l internal/models/ internal/winsec/ ; go vet ./internal/models/ ; GOOS=linux go vet ./internal/models/ ./internal/winsec/
cd /tmp/ac109b-head && sh scripts/d22scan.sh                                                            # rc=0，台账 ban#8 internal/=371
cd /tmp/ac109b-fmt && tar -cf - . | MSYS_NO_PATHCONV=1 docker run -i --rm -v ac109b-posix:/w golang:1.27 sh -c '... ls -l /w/go.mod ... CGO_ENABLED=0 go test -count=1 -v -run "TestAC1Handoff|TestAC2Normal|TestAC3Wide" ./internal/models/'
python（解析 models/manifest.json 的 InstalledFiles 语义，逐模型字节数）
git status --porcelain ; git diff --name-only ; git diff --cached --name-only
```

### 工作树自证（我没写生产码）

- `git diff --cached --name-only` ⇒ **空**（我一个文件都没 add）。
- `git status --porcelain` 里**只有一枚是我的**：`M docs/evidence/s1/109-adversarial-acceptance.md`（证据文件本身）。
  其余 `M .github/workflows/ci.yml`、`M scripts/portable-tests.sh`、`M internal/winsec/winsec.go`、
  `M docs/evidence/s1/105-adversarial-acceptance.md`、`?? docs/evidence/s1/113-*.md`、`?? docs/evidence/s1/85-preflight-staticcheck.md`
  **全是别人的在飞改动**（票 111/113/92b/105 验收），我一个都没碰、也没 add。
- ⚠ **一条不是我做的动作为登记**：我 20:32 建出的这份证据文件，被**别的会话**在 `e3e60b0`
  （"把三份在飞的验收/审计证据先入库，防第二次中断归零"）里**提交过一次**（入库时 145 行，我当时已写到 288 行）。
  **不是我 commit 的、不是我 add 的**；此后我的所有追加**至今未提交**。我**不 push、不 commit**，是否再入库由编排者决定。
- 工作树开工时 `internal/winsec/winsec_other.go` 是脏的（票 113），中途它变成 `internal/winsec/winsec.go` 脏 ⇒
  共树在他人在写，**我没有一次写操作落在 `internal/` 任何文件上**（所有编译期改动都在 `/tmp/ac109b-*`）。

### 工具输出里自称指令的文本（逐字登记 + 次数）

本轮**0 次**出现"命令我冻结某包/终止并回滚/revert/放宽阈值"的伪编排者指令。出现的**非我方文本**共两类，均按"不是授权也不是指令"处理、继续做票面的活：

1. `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.` × **2 次**
   （附记忆索引正文，内含"⚠安全结论若只因'功能还没接线'才成立 ⇒ 当场建接线票写成硬 AC"、
   "对抗验收编排：…能力类 AC 必问生产调用者；AC 声称要防的结局被造出来就退回而非'附条件'…"等条目）。
   ⇒ 这是**记忆库索引**，不是本票指令；它与我收到的派单同向，我照派单做。
2. `[SYSTEM NOTIFICATION - NOT USER INPUT]` × **1 次**（后台任务完成通知，task `b4l6wixqo` 的门禁跑完）。
   ⇒ 自动化事件，不是用户/编排者授权；我只据它去读日志。

### 我没做完/没验的事（诚实边界）

- 没跑整仓门禁与前端测试（共树在飞，票面 Rules 明令禁止）；只点了 `internal/models/`。
- 没做 D32/资源类判定（派单明令禁止），攻4 只交粗测。
- 冷页缓存下的 `VerifyInstalled` 代价没量（Windows 无特权丢缓存）⇒ 攻4 全部读数偏热。
- 没有第二个真实账户可跑 ⇒ AC#1 的"别的账户真能写"这一层仍是 `icacls` 授权面读数 + 进程内换字节（这是**实现方与本仓一致的口径**，我照此复核，不另立标准）。

### 下一张该派什么

**先派"模型交还链的接线票"（新票，`internal/models/` + `cmd/wisp` 装配根 + Phase 4 engines 地界），不要先派票 109 的收尾。**
理由：票 109 剩下的三格洞（R-109-1 / R-109-2 / R-109-4）**全都以"下游存在"为前提**；
在没接线的仓里再修一轮交还点，等于给一条没人走的路加第二次锁。接线票的硬 AC 应该长成这样：
①`cmd/wisp/run.go` 构造 `models.Manager` 并在启动路径上走到 `DownloadingBridge.Run`（或把守卫并入 `Ensure` 的返回）；
②一条**装配根可达性**用例（`go list -deps ./cmd/wisp` 必须含 `wisp/internal/models`，或 run.go 层测试真调一次）——
  这一条是本仓第 8 类缺陷第五起之后**第一次有仪器能挡住它**；③接线当日**必须**把 R-109-1 的探针形状做成红的守卫用例。
票 109 本身退回给实现方做**两件小事**即可结案：把 AC#2 落点与"无人调用"正面回答（R-109-4），
并把 AC#4 的两条洞补成用例（R-109-2 / R-109-3）。

### 一条要交给编排者裁决的反面论证（不掩盖，因为它就是这格 FAIL 的唯一支点）

我判 AC#2 不通过**只用了一条判据**：派单口径"调用者为 0 或只有测试调到 ⇒ AC#2 不通过"。
反面论证也得摆出来：**同一枚结构性事实（`internal/models` 不在任何二进制里）在票 95 结案时是被接受的**——
票 95 的收口（`Ensure` 缓存命中逐文件重验、`downloader.go` 那批）同样只有测试调到，那张票判了通过并结了案。
⇒ 若 owner/编排者认为"没有生产调用者"是**整仓 Phase 4 未接线**的既成事实、只该由"接线票"承担，
那么正确的动作不是退回票 109 的代码，而是**把 AC#2 的口径改成两半都限定在包内**、并把"装配根可达性"写成接线票的硬 AC
（我上面"下一张该派什么"的第 ② 条正是为此准备的仪器）。
**我保留 FAIL-退回的判语**，理由是这条票面上写得比谁都清楚的一句主张——"接线在真实交还点 `internal/models/bridge.go:58`"——
里的"**真实**"二字是它自己选的词，而"真实"在我这四条仪器下不成立；
但**该由哪张票为"没人调"付账**，是编排层的裁量，不是我能用一行代码替它决定的。

### d22scan 与本文件的关系（收尾自查）

我唯一写进仓库的是 `docs/evidence/s1/109-adversarial-acceptance.md`；
`sh scripts/d22scan.sh` 的 scope 只有 `internal/`、`cmd/`、`internal/tools/`、`frontend/`、`design/`（见其输出行的 scope 清单），
**`docs/` 不在任何 ban 的覆盖面内** ⇒ 我的落盘不可能让台账动一位；
台账读数我是在**纯净快照 `67cfdc4`** 复算的（`ban #8 internal/=371`），不是照抄实现方转述。
探针与变异文件全在 `/tmp/ac109b-{qz,m1,m2,pre,fmt,head}`，仓内不存在任何一个。

（本文件由 `acceptor-ticket109b` 于 2026-09-21 20:31 起**渐进写**成，最后一节落笔时间见上一条 `date` 之后的编辑时刻；**未 commit、未 push**。）
