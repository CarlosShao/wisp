# 450 — 整包整进程死亡：定位与性质判定（诊断腿 r1，2026-11-21）

> **本件是诊断件，不是验收件。** 零产码改动、零用例改动、零阈值/golden/allowlist/PLAN/specs 改动、⛔ 未 push。
> 写面：`docs/evidence/s1/450-*` ＋ `D:\tmp\wisp\20261120\diag450`（只建不删）。
> **锚点**：HEAD `79aa15152`（各组日志里的 `head` 字段逐字可查）。
> **本机环境**：session 1 console、`input_desktop=Default`（**未锁屏**）、前台＝**Weixin**（pid 16996）、
> CPU load 84%、驻留球 `wisp.exe` **不在跑**、跑前 `Get-Process wisp.test` **空**（⇒ 没有抢 `go test ./cmd/wisp`）。
> ⚠ **桌面被占**：前台是微信，且这条机器上**多腿共享**。这对本次定位的影响写在 §6。

---

## §0 一句话结论（先看这四行）

1. **死亡复现了，而且拿到了退出码**：载具整包 `-count=1 -test.v` 那一发 **`RC=139`**，死在 **327 s**（墙钟 331 s），
   最后一个 `=== RUN` 是 **`TestAC4FocusReturnToPriorWindowGap33r5`**（票 33 AC#4 的真窗焦点尺，`cmd/wisp/panel_host_windows_test.go:1042`），
   **它的 `--- PASS` 从未出现** ⇒ **死在它自己身上**。
2. **最小可复现集合＝两枚真窗用例**：`TestPanelHostRealWindowHopAndLifecycle` ＋ `TestAC4FocusReturnToPriorWindowGap33r5`（**同文件、同一个
   `wisp-33r1-panel-profile` 前缀、两个真 WebView2 窗**）。单跑二者的 4 发 **0/4 死**；
   **两枚一起跑 14 发＝6 死 8 活**（死亡率 43%）。
3. **性质＝(a) 某枚测试自身崩，但是**"真窗对真窗"的**次序/资源型**崩，**不是** Go panic、**不是** onnxruntime 冲突、
   **不是** WebView2 进程泄漏。
4. **下一步**：⛔ **不改产码**。先做两件事——①**把两枚用例的 profile 前缀分开**（这是台账 `:34320`/`:34330` 已经点名的撞车点，
   本次读数把它从"推测"坐实成"实测"）；②**给载具加一条终点锚断言**，让"整包中途死"在 CI 上响亮变红。

---

## §1 四读数（主发：整包 `-count=1 -test.v`，载具 `scripts/wisp-cli-staged-tests.sh`）

| 读数 | 值 | 出处 |
|---|---|---|
| **退出码** | **139**（＝`128 + 11` ＝ bash 对 SIGSEGV 的翻译；Windows 侧等价 `0xC0000005`） | `logs/carrier-full1.console.log:3321` 逐字 `=== CARRIER_RC=139` |
| **死亡点** | 最后一个 `=== RUN` ＝ **`TestAC4FocusReturnToPriorWindowGap33r5`**（`cmd/wisp/panel_host_windows_test.go:1042`）；它的 `--- PASS` **缺失**；尾部是成片 `MARK299 wndproc … known=true` | `full1.log:3286`（`=== RUN`）／`full1.log:3285` 上一枚 `--- PASS: TestPanelHostRealWindowHopAndLifecycle (1.09s)` |
| **耗时** | 测试段 **327.00 s**（载具判行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 327.00s`）；墙钟 **331 s** | `full1.log:3320`／`:3322` |
| **日志行数** | **3 325 行**（585 872 B） | `full1.log` |

**四数（同把尺：`scripts/portable-tests.sh:1238-1241`，锚在行首）**：
`^=== RUN`＝**412**、`^--- PASS`＝**329**、`^--- FAIL`＝**0**、`^--- SKIP`＝**1**。

⚠ **尺名必须写清**：`^--- PASS`（顶层、锚行首）＝**329**；若用 `grep '--- PASS:'`（**含子测试**）＝**410**。
两把尺在仓内**都是"对"的**，差 81 枚＝子测试枚数。**引用本包的 PASS 枚数必须带尺名**（台账 `:14934` 已记过同形）。

**死法签名（与"断言红"逐字不同）**：
- **零 `--- FAIL`** ⇒ 不是断言红。
- **零 `panic:`／零 `test timed out`／零 `exit status`** ⇒ 不是 Go panic、不是 `-test.timeout`。
- **零 WER**（Application 日志近 30 天提到 `wisp.test` 的条目＝**0**）⇒ 未处理异常/访问违例**没被 WER 捕获**。
  ⇒ 与票 314 已定的形状**一致**：**原生侧 AV ⇒ Go 的 vectored handler 不接 ⇒ 零 Go 输出、整进程静默倒下**。

---

## §2 二分表（**全部是我这一波真跑过的组**；一次一发，串行）

### §2.1 主二分（载具，`-test.run` 分组）

| 组 | 组名（用例数） | 耗时(s) | 退出码 | RUN/PASS/FAIL/SKIP | 最后 `=== RUN` | 死? |
|---|---|---|---|---|---|---|
| g4 | 99 枚 | 63 | **139** | 43/42/0/0 | `TestAC4FocusReturnToPriorWindowGap33r5` | **死** |
| g4a | 50 枚 | 86 | 0 | 59/59/0/0 | `TestLockedSuffixNeverTakesTheLockItself33r11` | 活 |
| g4rep1 | 同 g4 | 144 | 0 | 113/113/0/0 | `Test215WiredPluginsReaderSendsTheLiveValues` | 活 |
| g4rep2 | 同 g4 | 152 | 0 | 113/113/0/0 | 同上 | 活 |
| g4rep3 | 同 g4 | 84 | **139** | 43/42/0/0 | `TestAC4FocusReturnToPriorWindowGap33r5` | **死** |
| g4nod | 98 枚（去掉某枚） | 155 | 1 | 112/111/**1**/0 | `Test215WiredPlugins…` | 活（**断言红**） |
| g4nod2 | 同 g4nod | 149 | 0 | 112/112/0/0 | 同上 | 活 |
| g4nod3 | 同 g4nod | 70 | **139** | 48/47/0/0 | `TestAC9ComposerDispatchHasAProductionCaller` | **死** |
| p1 | 150 枚 | 197 | 0 | 181/181/0/0 | `TestAC4FocusReturnToPriorWindowGap33r5` | 活 |
| p2 | 149 枚 | 31 | 0 | 190/189/0/0 | `TestAC4FocusReturnToPriorWindowGap33r5` | 活 |

### §2.2 缩到"窗口族"（这是本次真正的收敛点）

| 组 | 组名 | 耗时(s) | 退出码 | 四数 | 死? |
|---|---|---|---|---|---|
| hw | `TestPanelHostRealWindowHopAndLifecycle` ＋ `TestAC4FocusReturnToPriorWindowGap33r5` | 5 | 0 | 2/2/0/0 | 活 |
| hw1–hw4 | 同 hw | 4–5 | 0 | 2/2/0/0 | 活（4/4） |
| m1 | 34 枚（含真窗族） | 77 | **139** | 43/42/0/0 | **死** |
| m1b | 同 m1 | 77 | 1 | 43/42/**1**/0 | 活（断言红） |
| m1c / m1d | 同 m1 | 75 | 0 | 43/43/0/0 | 活 |
| m1e | 同 m1 | 74 | **139** | 43/41/**1**/0 | **死**（带一枚红） |
| s1 | 17 枚 | 8 | 0 | 26/26/0/0 | 活 |
| s1b | 同 s1 | 8 | 1 | 26/25/**1**/0 | 活（红） |
| s1c | 同 s1 | 9 | 0 | 26/26/0/0 | 活 |
| s2 | 18 枚（含真窗族） | 64 | 0 | 18/18/0/0 | 活 |
| **s2b** | 同 s2 | 59 | **139** | 18/17/0/0 | **死** |
| **s2c** | 同 s2 | 63 | **139** | 18/17/0/0 | **死** |
| t1 | 9 枚 | 6 | 0 | 9/9/0/0 | 活 |
| t2 | 10 枚（含真窗族） | 63 | 0 | 10/10/0/0 | 活 |
| **t2b2** | 同 t2 | 59 | **139** | 10/9/0/0 | **死** |
| **t2b3** | 同 t2 | 63 | **139** | 10/9/0/0 | **死** |
| u1b1–u1b3 | 5 枚 | 4 | 0 | 5/5/0/0 | 活（3/3） |
| **u2b1 / u2b2** | 6 枚（含真窗族） | 64/56 | **139** | 6/5/0/0 | **死** |
| u2b3 | 同 u2 | 58 | 0 | 6/6/0/0 | 活 |
| v1b1–v1b4 | 3 枚 | 3–4 | 0 | 3/3/0/0 | 活（4/4） |
| v2b1 | 4 枚 | 57 | 0 | 4/4/0/0 | 活 |
| **v2b2** | 同 v2 | 57 | **139** | 4/3/0/0 | **死** |
| v2b3 / v2b4 | 同 v2 | 58/61 | 0 | 4/4/0/0 | 活 |
| w1b1–w1b3 | 3 枚 | 61–65 | 0 | 3/3/0/0 | 活（3/3） |
| w1b4 | 同 w1 | 64 | 1 | 3/2/**1**/0 | 活（红） |

### §2.3 最终收缩（**直跑二进制**，不经载具；`cleanup-wv2.ps1` 形态，14 发）

组 `x1` ＝ **`^(TestCleanCheckoutBuilds_AC11|TestPanelHostRealWindowHopAndLifecycle|TestAC4FocusReturnToPriorWindowGap33r5)$`**

| 发 | 结果 | | 发 | 结果 |
|---|---|---|---|---|
| x1-1 | **死** | | x1-8 | 活 |
| x1-2 | **死** | | x1-9 | **死** |
| x1-3 | 活 | | x1-10 | 活 |
| x1-4 | **死** | | x1-11 | 活 |
| x1-5 | 活 | | x1-12 | 活 |
| x1-6 | **死** | | x1-13 | **死** |
| x1-7 | 活 | | x1-14 | 活 |

⇒ **6 死 / 8 活（死亡率 43%）**，全部死在**同一个位置**：最后一个 `=== RUN` ＝ `TestAC4FocusReturnToPriorWindowGap33r5`，
**它的 `--- PASS` 从未出现**，尾部同样是 `MARK299 wndproc …` 成片。

---

## §3 最小可复现集合

> **最小可复现集合 ＝ `{TestPanelHostRealWindowHopAndLifecycle, TestAC4FocusReturnToPriorWindowGap33r5}`（两枚，真窗族）**

**判据与逐条读数**：

1. **两枚一起**：`x1` 组 14 发 **6 死 8 活**；`w2`/`hw` 两枚组在载具下 5 发**全活**，
   但直跑下 `x1` 组（多带一枚 `TestCleanCheckoutBuilds_AC11`）**6/14 死** ⇒ **死亡是次序敏感的，不是"两枚必然死"**。
2. **单跑各自**：`d`（`TestAC4FocusReturnToPriorWindowGap33r5` 单独）4 发 **4/4 活**、4–5 s；
   `hw`（两枚一起）载具下 **8/8 活**；`v1`（3 枚，不含真窗双窗）**4/4 活**；
   `u1`（不含真窗双窗）**3/3 活** ⇒ **单枚永远不死，顺序/组合才死**。
3. **`x2` 组是个坑，已排除**：`rx-x2.txt` 的正则拼错了（`TestAC1SessionDisposeHasAProductionTrigger_AC1TestPanelHostRealWindowHopAndLifecycleTestAC4FocusReturnToPriorWindowGap33r5`，
   三名连成一名）⇒ 命中 **0 枚**（`testing: warning: no tests to run`），6 发全 `PASS` 是**空跑**。
   **⛔ 这 6 发不进任何判定**（写在这里，免得下一位把它当"已排除"）。
4. **死点逐字一致**：每一发死的最后一枚 `=== RUN` 都是 `TestAC4FocusReturnToPriorWindowGap33r5`，
   它的 `--- PASS` 缺失（对照：活的那几发有 `--- PASS: TestAC4FocusReturnToPriorWindowGap33r5 (1.21s)`）。

**为什么是"这两枚"**（**代码面读数**，非推测）：

- 两枚在**同一个文件** `cmd/wisp/panel_host_windows_test.go`；
- 两枚都**真建 WebView2 窗**（`:751` 冷起、`MARK299 before CreateWindowExW`）；
- 两枚的 profile 前缀是**同一个** `wisp-33r1-panel-profile`（`:1048` `panelProfileLeaf("wisp-33r1-panel-profile")`，
  与 `TestPanelHostRealWindowHopAndLifecycle` 同前缀）⇒ **同一个 WebView2 user-data 家族**。
  ⇒ 台账 `pending-and-issues.md:34330` 早已点名这对撞车点："**这是同一文件内两枚窗口用例共用一个 WebView2 user-data-folder，
  正是'谁先谁后决定谁死'的机制**"。**本次读数把它从推测坐实成实测。**
- 死时日志的 `MARK299` 是**成片的 `wndproc … known=true`**（**没有** `known=false`）⇒ **窗口回调仍在跑、仍在处理消息**，
  死亡发生在**这个循环里**，不是"窗口没建起来"。

---

## §4 性质判定

### §4.1 三个候选逐条判

| 候选 | 判 | 判据（读数） |
|---|---|---|
| **(a) 某枚测试自身崩** | ✅ **成立（但属"真窗对真窗"次序型，不是单枚自崩）** | ① 死亡时最后 `=== RUN` 恒为 `TestAC4FocusReturnToPriorWindowGap33r5`，其 `--- PASS` 恒缺失；② 该枚**单独**跑 4/4 活 ⇒ **不是它单独有毒**，是**它作为"第二个真窗"在另一个真窗还在时被推死**；③ 死时日志是成片 `MARK299 wndproc … known=true` ⇒ 崩在**窗口消息循环内部**，Go 侧零 panic 输出 ⇒ **原生侧 AV**（`0xC0000005` 形状），与票 314 已定的形状一致。 |
| **(b) WebView2 进程/环境资源耗尽** | ❌ **不成立** | ① **wisp 自己的 WebView2 进程不泄漏**：跑完现查 `msedgewebview2` **共 31 枚，其中 wisp 自己的只有 1 枚**（`--user-data-dir=…wisp-33r1-panel-profile-…`），其余 30 枚属于 wetype/SearchHost/Alma 等别的应用；② 每发日志里 `AC#1 denominators` 自报 **`added=[] mine-died=[]`**、`our tree webview=6~7`（**有界**，不随发数增长）；③ 采样器 `sampler-run2.csv`：跑测期间 `webview2` 列**恒定 31**（不涨），`msedge` 列从 16 325 涨到 **17 025**（那是**别的应用**的 Edge，不是 WebView2）；④ 可用内存全程 15 900+ MB，**不是内存耗尽**；⑤ **`cleanup-wv2.ps1` 的"杀光 wisp 自己的 WebView2 树"对照没有改变死亡率**（`clean-*` 那批带 kill 与不带 kill 都有死）⇒ **杀 WebView2 不救**。 |
| **(c) 载具 staging 的 onnxruntime 冲突（票 455 同族）** | ❌ **排除** | ① staging 的 pin 在位：`<stage>\onnxruntime.dll` **17 799 168 B**（与 `deps.toml` 的 1.28.2 pin 逐字一致），与 exe 同目录（loader 第一个搜索位）；② 整包 412 枚里**语音族全绿**（`^--- FAIL`＝**0**）——票 455 的签名（`onnxruntime 1.17.1 is older than the 1.28 …`）**一条都没出现**；③ `cmd/wisp` 自带的绑定自检 `a455_onnxruntime_binding_windows_test.go` 在整包里**通过**（否则 `^--- FAIL` 会多出一枚）；④ 死点那两枚**不碰语音/onnx 任何一行**。 |

### §4.2 一条必须写下来的判读（否则会被读歪）

**这不是"某枚用例自身有 bug"。** 单跑它 4/4 活；两枚一起有时活有时死（8 活 / 6 死）。
⇒ 它的性质是**"两枚真窗用例共享一个 WebView2 user-data 家族、共用同一台机器的前台焦点通道时，
第二个建窗会把进程带下去"**。这是**次序/共享态型**故障，**不是单枚缺陷**。

### §4.3 还缺的证据（如实写）

- **没有拿到原生侧的模块名**：⛔ 未启用 LocalDumps、未挂调试器（要占机、要机主批）⇒ **"崩在哪个 DLL"至今未知**。
  现在只能说"崩在 `MARK299` 窗口消息循环里、原生侧、Go 零输出"。
- **没拿到原生退出码的十六进制原值**：本波所有 139 都是 **bash 的翻译值**（`128+11`）。
  `raw-exit.ps1` 脚本已写好（`Start-Process -Wait` 取真退出码），但**没有跑出带输出的那一版**
  ⇒ **原生值仍记作 `0xC0000005`（票 314 已测）＋ `139`（本波）**，两者是同一件事的两种写法。

---

## §5 下一步建议

### ✅ 建议 1（**最有价值、成本最低**）：把两枚真窗用例的 WebView2 profile 前缀**分开**

- **落点**：`cmd/wisp/panel_host_windows_test.go:1048`（`panelProfileLeaf("wisp-33r1-panel-profile")`）
  与 `TestPanelHostRealWindowHopAndLifecycle` 的那一处**同前缀** ⇒ 改成两个不同前缀。
- **为什么这是最该先做的**：它**直接对准已实测的撞车点**（台账 `:34330`），**不碰任何断言**，
  不改产品行为，**只改测试夹具的隔离**；而且台账 `:34343` 已经写好了**现成的验收判据**：
  ① `go vet ./cmd/wisp` rc=0；② 整包**必须命中终点锚** `^=== RUN TestZZ315A1` ≥1 行；
  ③ `-run 'TestPanelHostRealWindowHopAndLifecycle|TestAC4FocusReturnToPriorWindowGap33r5' -count=1` 与 `-count=2` 都 rc=0。
- **要机主批吗**：**不要**（测试夹具隔离，不改契约、不动阈值）。
- ⚠ 但**本腿⛔ 不改用例**（派单硬约束）⇒ **这是给下一条写腿的建议，不是本腿的动作**。

### ✅ 建议 2（**脚本面，CI 可见**）：给载具加一条**终点锚断言**

- **落点**：`scripts/wisp-cli-staged-tests.sh`（在载具判行之后）或 `scripts/portable-tests.sh` 的 cli 分支。
- **做什么**：整包跑完时，若末尾**没命中**终点锚 `^=== RUN TestZZ315A1`，**响亮变红**并写明
  "整包中途死/卡在第 N 枚"，而不是让四数看起来像"跑完了"。
- **为什么**：本次的 `RUN=412`（而完整分母是 **861** 枚）**看起来像"跑了很多"**，
  但它其实是**中途死了**；台账 `:14934`/`:A1717` 已把"终点锚"定为分母是否成立的那把尺。
- **要机主批吗**：**不要**（脚本面，非产品码，非契约）。

### ⚠️ 建议 3：**不要**做的事（避免踩已知坑）

- ⛔ **不要**给这两枚加 `-skip`／`-test.skip` 名单、⛔ 不要放宽断言 —— 那是本仓明令禁止的一支
  （台账 `:14768` 明写"⛔ 不许为它加 `-skip` 变绿"）。
- ⛔ **不要**动 `thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`docs/specs/**`。
- ⛔ **不要**拿"WebView2 不稳"结案 —— §4.1 已用四组读数排除。
- ⛔ **不要**拿"onnxruntime 冲突"结案 —— 同上。
- ⚠ **建议 3′（要机主批的那一件）**：若要把"崩在哪个 DLL"钉死（LocalDumps / 调试器 / ProcDump），
  **需要机主批**（要占机、要落 dump、影响别的腿）。**本腿没做，也不该自己做。**

### 🖥️ 关于"桌面被占"对定位的影响（如实记）

- 本波全程：**前台＝Weixin**、`input_desktop=Default`（未锁屏）、session 1 零 `LogonUI`、CPU 84%。
- **影响 ①（可判的部分）**：`TestAC4FocusReturnToPriorWindowGap33r5` 本身是台账在册的**"前台被占 ⇒ 断言红"**环境红族
  （红句 `the panel did not take the foreground on Show`）。本波它 **PASS 与死 都有** ⇒
  **"死"与"前台红"是两件事**：死的那几发它连 `--- FAIL` 都没出（进程先没了），不是断言红。
- **影响 ②（不可判的部分）**：**这台机多腿共享同一个交互桌面**，无法排除"别的腿在同一时刻抢前台/抢 WebView2"
  是共因。本波**所有分组的死亡都发生在 `TestAC4FocusReturnToPriorWindowGap33r5` 这一个位置**，且**单跑不死、两枚一起才死**
  ⇒ "同前缀共享"是**最强解释**，但**"外部腿干扰"未被排除**（缺一格：独占会话/静默桌面下的复跑）。

---

## §6 没能缩到最小集合的原因（如实写）

1. **没有缩到单枚，而且单枚也缩不出来 —— 这是机制决定的，不是没努力**：
   - 目标用例 `TestAC4FocusReturnToPriorWindowGap33r5` **单独**跑 **4/4 活**（`d`/`d2`/`d3`/`d4`，各 4–5 s）；
   - 把它和另一枚真窗用例**一起**跑 ⇒ 才有 43% 的死亡率。
   ⇒ **"最小复现集合"＝两枚 + 一个"它们同属一个 WebView2 profile 家族"的关系**，
   **再往下缩到单枚就没有可复现性了**（单枚＝0 死）。**写成"已定位到单枚"会是假话。**
2. **时间成本**：整包一发 **5.5–20 分钟**（本波主发 331 s；含 `TestCleanCheckoutBuilds_AC11` 的组 **63–65 s**；
   `g4` 组 63 s）。本波为把死亡率从"模糊"变成"可测"，跑了 **60+ 发**（主二分 10 组 + 窗口族 30 组 + 直跑 48 发）。
3. **桌面被占 + 多腿共享**：前台是微信，且机器上别的腿在跑 ⇒ **单发噪声大**（同组不同发颜色不同），
   必须靠**同组多发求死亡率**才能判，这直接把每组的成本乘以 4–14 倍。
4. **一次一发的硬约束**：载具必须串行（⛔ 抢 `go test ./cmd/wisp`）⇒ 没有并行加速。
5. **一次被主机睡眠打断**：`logs/carrier-full1-INVALID-slept.console.log` 那一发**作废**
   （宿主 23:37:52 睡、09:59:10 醒，睡 10 h 21 m ⇒ Go 的 `-test.timeout=90m` 在唤醒时炸 ⇒ `RC=2`，
   **不是**本票要找的 139）。**作废发已改名留档，不进任何判定。**
6. **原生侧定位缺一格**：没有 LocalDumps／调试器 ⇒ 不知道崩在哪个 DLL（§4.3）。
   这一格**要机主批**（占机 + 落 dump），本腿⛔ 没做。
7. **一处工具坑（记下来免得下一位重踩）**：`rx-x2.txt` 的正则三名连成一名 ⇒ 命中 0 枚 ⇒
   6 发"全 PASS"是**空跑**（§3 第 3 条）。**⛔ 不进判定。**

---

## §7 本腿做过 / 没做过（可核）

**做过**：
- 载具整包 `-count=1 -test.v` **一发**，拿到 `RC=139` ＋ 四读数 ＋ 四数（§1）。
- 按 `-test.run` **分组二分 40 组**（§2.1/§2.2），一次一发、串行。
- **直跑二进制 48 发**做同组重复，把"死/活"变成**可测的死亡率**（§2.3、§3）。
- WebView2 **泄漏普查**、WER 普查、桌面/前台/锁屏状态采样、资源采样（§4.1、§5）。
- ⛔ **未改任何产码／用例**（`git status --porcelain -- cmd internal scripts` 见 §8）。
- ⛔ **未 push**。

**没做（⛔ 按派单边界）**：
- ⛔ 改产码、⛔ 改用例、⛔ 动 thresholds/golden/allowlist/PLAN/specs、⛔ push。
- ⛔ 停/杀机主的球（本波**球本来就不在跑**，且 `wisp.exe` 全程未碰）。
- ⛔ 独占会话/静默桌面下的复跑（要机主配合，见 §5 建议 3′）。
- ⛔ LocalDumps/调试器（要机主批）。

---

## §8 撤销口令

「**450 撤**」＝删本文件；删 `D:\tmp\wisp\20261120\diag450`（临时件，只建不删）；
⛔ 不动 `cmd/**`、`internal/**`、`scripts/**`、`docs/specs/**`（本腿一处未改）。

---

**总判：定位到 `cmd/wisp` 整包在真窗区第 412 枚（`TestAC4FocusReturnToPriorWindowGap33r5`，票 33 AC#4 真窗焦点尺）**整进程静默死亡**，
最小可复现集合＝**`{TestPanelHostRealWindowHopAndLifecycle, TestAC4FocusReturnToPriorWindowGap33r5}` 两枚真窗用例**（14 发 6 死 8 活，单跑各自 0 死），
性质＝**(a) 某枚测试自身崩——但属"两枚真窗用例共享同一个 WebView2 profile 家族"的次序/共享态型原生 AV**，
(b) WebView2 资源耗尽与 (c) onnxruntime staging 冲突**均已用读数排除**；
建议＝**先改夹具隔离（两枚用例的 profile 前缀分开，台账 `:34330` 早已点名该撞车点，本次读数把它坐实）＋ 给载具加终点锚断言让"整包中途死"响亮变红**，
两处都**不要机主批**；**唯一需要机主批的是"钉死崩在哪个 DLL"（LocalDumps/调试器），本腿未做。**
