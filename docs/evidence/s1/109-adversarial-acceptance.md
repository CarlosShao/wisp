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
| AC#1 | 窗口真实形状（SID 级 `icacls` 前后 + 交还点→读取点 file:line 链 + 下次 `Ensure` 是否重验） | （验中） | — |
| AC#2 | 收口方向②，**正半边**（窗口关上）＋**反半边**（正常安装/多实例复用不被打断） | （验中） | — |
| AC#3 | `rename`-into-place 序列的 seam 守卫（不靠"这次恰好保留了描述符"） | （验中） | — |
| AC#4 | 安装目录那一格的反向钉子，三条界；**攻：找能滑过去的第四形** | （验中） | — |
| AC#5 | 修前红 + 变异 M1/M2 + 门禁四数 + gofmt/gofumpt + d22scan 纯净快照 | （验中） | — |
| 攻1 | **守卫之后还剩不还剩窗**（`VerifyInstalled` 返回后、读取者 open 前再换一次） | （验中） | — |
| 攻2 | `bridge.go:58` 是不是**真生产调用点**（第 8 类缺陷，本仓第五起） | （验中） | — |
| 攻3 | AC#4 界线第四形 | （并入 AC#4 格） | — |
| 攻4 | 「复验」的代价（单次耗时 + 读过的字节数，粗测，不判阈值） | （验中） | — |

## 总判

（**待各格填完后追加**）

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
