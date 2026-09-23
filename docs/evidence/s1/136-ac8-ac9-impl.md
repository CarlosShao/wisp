# 136 — AC#8（判据②③）＋ AC#9：实现方证据

- 票面：`.scratch/wisp/issues/136-the-zero-sample-fail-closed-guard-underneath-the-freshness-nail-has-no-nail-54-tests-stay-green-without-it.md`
- 实现方代理：`worker-ticket136-ac8-ac9`（只交 AC#8②③ 与 AC#9 两格；AC#8① 归验收方，AC#2..AC#7／AC#10 一字未碰）
- 时间：2026-09-23 22:47 +08（`date -u` 原文 `Wed Sep 23 14:47:12 UTC 2026`，手工 +8；本 shell 里 `TZ=Asia/Shanghai` 不生效）
- **本程未翻任何勾**：AC#8／AC#9 的 `[ ]` 原样留着，终裁归非实现者。

---

## §0 锚点、地界与仪器（开工第一步重量的，不照抄派单）

| 项 | 实测 |
| --- | --- |
| 开工第一次 `git rev-parse --short HEAD` | `1d38206`（全 sha `1d382064b40c6e8de787f37374dfe8e4dd764e30`）— **本程唯一锚点** |
| 分支 | `dev` |
| 开工 `git status --porcelain` | 只有 1 枚未跟踪件：`docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` |
| 派单说"HEAD 大约 1d38206，133 的修方还在推提交" | 我量到 **正是 `1d38206`**，与本程开工那一刻无漂移。派单里那句"你看到的会不同"这条前提**未成立**（照原样登记，不影响判据） |
| 在飞兄弟 `worker-ticket133-ac2-fix` 的地界 | `cmd/wisp/leg_dispatch_gate_133_test.go` ＋它自己的证据文件。**`cmd/wisp/**` 本程整块未碰**（见 §7 暂存清单） |
| 未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` | **未读内容、未提交、未改、未删、未据它开任何一格或改任何判据** |

**仓外纯净树（只建不删，全部带本会话后缀；仓内未建 worktree、未 checkout）**

| 树 | 怎么来的 | 用途 |
| --- | --- | --- |
| `/d/tmp/wisp136ac8-ac9-tree0` | `git archive 1d38206 \| tar -x -C …` | AC#8 改前面（复现 panic 吞读数） |
| `/d/tmp/wisp136ac8-ac9-pristine-tree0/internal/observe/{sampler.go,sampler_test.go}` | 从 tree0 拷出 | 每发变异驱动前的**无条件 restore 源** |
| `/d/tmp/wisp136ac8-ac9-tree1` | `git archive 79ddd49 \| tar -x -C …`（`79ddd49` ＝本程 AC#8② 那枚 commit） | AC#8②③ 与 AC#9 的读数面 |
| `/d/tmp/wisp136ac8-ac9-pristine-tree1/internal/observe/{sampler.go,sampler_test.go}` | 从 tree1 拷出 | 同上 |
| `/d/tmp/wisp136ac8-ac9-mut.py` | 手写驱动器 | 单发：先 restore 再打，锚文本命中数≠1 直接退出（防叠发、防打偏） |
| `/d/tmp/wisp136-acc-r1-tree2/` | 验收方留下的可重跑凭据 | **未动、未删**，本程没用到它（M10 我在自己的树上重造） |

**工具与版本（写明，不假设）**

- `go version go1.27.1 windows/amd64`（`/d/work/base/go/bin/go`）
- `gofmt` ＝同一条 go 工具链里的 `/d/work/base/go/bin/gofmt`
- `gofumpt`：**`/d/work/base/gopath/bin/gofumpt.exe` ＝ `v0.12.0 (go1.27.1)`**（本机现有二进制；PATH 里没有 `gofumpt`，我没有为了跑门禁去重装工具链，CI 那一步是 `@latest`、版本未钉，这条边界记在 §6）
- `tar` `/usr/bin/tar`、`git` `/mingw64/bin/git`

**锚点基线（未变异，tree0，`go test -count=1 -v ./internal/observe/`）**

```
rc=0 / === RUN=56 / --- PASS=56 / --- FAIL=0 / SKIP=0
```

四数一律从 `-v` 量（`-count=2` 那份在 §6，分栏不混装）。

---

## §1 AC#8 改前对照：同一发 M3 在**无守卫**的树上确实吞读数（不是重做判据①，是给②当分母）

判据①（"造一发空 `Samples` 变异 ⇒ 今天必须 panic 且拖走若干枚"）验收方已在 `136-ac1-adversarial-acceptance.md` §5 独立复算成立。我在**自己的锚点**上把同一发变异复跑一次，为的是②的前后对照能出自同一棵树、同一个二进制版本，不跨人借读数。

变异 M3（形状照验收方 §5.1 的"逢读数都丢"）：`internal/observe/sampler.go:290` 的 `} else if m.PrivateWorkingSetBytes <= 0 {` ⇒ `>= 0`。fixture 里的 footprint 全为非负，所以每一发读数都落进丢弃分支 ⇒ `rep.Samples` 处处为空。

**先证落地**

```
$ grep -n 'PrivateWorkingSetBytes >= 0' /d/tmp/wisp136ac8-ac9-tree0/internal/observe/sampler.go
290:		} else if m.PrivateWorkingSetBytes >= 0 {
$ diff -u <pristine> <tree1 的 sampler.go>     # 单行，其余一字未动
-		} else if m.PrivateWorkingSetBytes <= 0 {
+		} else if m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
```

**读数**（原文 `/d/tmp/wisp136ac8-ac9-tree0-M3-prefix-v.txt`）

```
rc=1 / === RUN=52 / --- PASS=47 / --- FAIL=5 / SKIP=0 / ^panic:=1
panic: runtime error: index out of range [0] with length 0 [recovered, repanicked]
```

红名逐名（5 枚，panic 前跑到并红的）：`TestSampleStateAllMetricsAndVerdicts`、`TestSampleStateSleepingDiskWriteGateFails`、`TestSampleStateWorkPeakMemoryIsTargetNotGate`、`TestSampleStateLeakFixtureFlipsRed`、`TestSamplerGoroutineAccountingFollowsRegistry`（它就是 panic 那一枚，`sampler_test.go:308`）。

被拖走（`=== RUN` 名册差集，基线 56 枚 － 该发 52 枚 ＝ 4 枚，逐名）：

| # | 被拖走 | 它本来在名册的位置（改后树实测） |
| --- | --- | --- |
| 1 | `TestLiveRegistryBaselineWithinSleepingGate` | 53 |
| 2 | `TestThresholdTableCoversAllStates` | 54 |
| 3 | `TestSampleStateZeroSampleWindowFailsClosed`（AC#1 那枚钉） | 55 |
| 4 | `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（AC#1 正向对照腿） | 56 |

⇒ 与验收方 §5.1/§5.2 的读数**逐位相同**（52/47/5 + panic、同一批 4 枚）。

**还原**：`python wisp136ac8-ac9-mut.py … restore` ⇒ `diff -q` 与 pristine **逐字相同**（并另证与仓库里那枚 `sampler.go` 也逐字相同）⇒ 复跑 `rc=0 / 56 / 56 / 0 / 0`。

---

## §2 AC#8 判据②：加长度守卫（只动那枚用例）＋ 同包 `[0]` 直取逐枚判定

### 2.1 改了什么（commit `79ddd49`，净面只有 `internal/observe/sampler_test.go`，`git diff --cached --numstat` ＝ `10 0`）

`TestSamplerGoroutineAccountingFollowsRegistry`（函数仍在 `:297`）里，原先的 `rep.Samples[0]` 直取（改前 `:308`）前面插入长度守卫，改后落在 `:314-317`，直取本身移到 `:318`：

```go
	// Ticket 136 AC#8: length guard in front of the [0] read below. With no
	// guard an empty rep.Samples panicked here, and the panic killed the test
	// binary, so every case scheduled after this one stopped reporting at all
	// (measured: 52 RUN / 47 PASS / 5 FAIL + panic, 4 cases swallowed). A
	// sampler window that holds no sample is a failure of this case, so it
	// says so and dies alone - never t.Skip, never a silent return.
	if len(rep.Samples) == 0 {
		t.Fatalf("sampler recorded 0 samples (sample_errors=%d last_error=%q): there is no reading to check goroutine accounting against",
			rep.SampleErrors, rep.LastSampleError)
	}
```

三条硬约束的对表：

| 约束 | 本程实际 |
| --- | --- |
| 不许改成 `t.Skip`／不许静默返回 ⇒ 要红 | `t.Fatalf`，且消息里带 `sample_errors`／`last_error` 两项自证。文件内 `Skip` 出现次数＝0 |
| 只许动那枚测试；`sampler.go` 的采样器语义一个字不许动 | 本程 AC#8 的 commit 只有 `sampler_test.go` 一枚文件、10 增 0 删；`git diff 1d38206..HEAD -- internal/observe/sampler.go` **无输出**（AC#9 那格同样没动生产码，见 §4） |
| 同包 `[0]` 直取逐枚判、别一把改完 | 见 §2.2：全量普查 17 处，**只改 1 处**，其余 16 处各给理由 |

### 2.2 同包 `[0]` 直取普查（`grep -rn '\[0\]' internal/observe/*_test.go`，17 处，逐枚判）

改一枚证明一枚（本轮只改了 `sampler_test.go:308`，读数在 §3）。其余判定如下——"判不改也要给理由"：

| # | 位置 | 直取的东西 | 前置守卫（实测的那一行） | 判定＋理由 |
| --- | --- | --- | --- | --- |
| 1 | `diagnostics_test.go:101` | `def[0].ID` | 同一表达式里 `len(def) != 1 \|\|` 短路 | **不改**：空 slice 时左项已为真，`\|\|` 短路后 `def[0]` 不可达，形状本身就是长度守卫 |
| 2 | `earlylog_130_test.go:127` | `got[0]/[1]/[2]` | 同行 `len(got) != 3 \|\|` 短路 | **不改**：同上 |
| 3 | `earlylog_130_test.go:131` | `recs[1]`、`k[0]` | `recs` 与 `got` 同源（`msgs()` 就是从 `snapshot()` 派生的，实测 `earlylog_130_test.go:62-69`），`:127` 已断 3 枚；`k[0]` 前有 `len(k) != 1 \|\|` | **不改**：索引源长度已被同一测试里的 `t.Fatalf` 钉住，且 `flushed != 3` 在 `:123` 先响 |
| 4 | `earlylog_130_test.go:134-135` | `recs[0]`、`recs[2]` | 同 #3 | **不改**：同因 |
| 5 | `earlylog_130_test.go:172-173` | `recs[0]` | `:169` `if len(recs) != earlyLogMaxRecords+1 { t.Fatalf }` | **不改**：长度守卫在索引之前，且是 Fatal 不是 Error |
| 6 | `earlylog_130_test.go:232` | `sink.snapshot()[0]` | `:230` `if flushed, _ := b.drain(sink); flushed != 1 { t.Fatalf }` | **不改**：`flushed` 就是落进 sink 的记录数，等价于长度守卫；真漂了会红在 `:231` 而不是 panic |
| 7 | `earlylog_130_test.go:241,243` | `outer[0]` | `:237` `if len(outer) != 2 { t.Fatalf }` | **不改**：守卫在前 |
| 8 | `earlylog_130_test.go:254,255,257,258` | `inner[0]`、`inner[1]` | `:250` `if len(inner) != 2 { t.Fatalf }` | **不改**：守卫在前 |
| 9 | `errors_test.go:41` | `classes[0] = "tampered"` | `:16` `if len(classes) != 17 { t.Fatalf }` | **不改**：守卫在前（写索引也一样红在 `:17`） |
| 10 | `errors_test.go:42` | `AllClasses()[0]` | 同一函数、`:16` 已断 17 枚、`AllClasses()` 无随机源 | **不改**：唯一能让它空掉的写法是 `AllClasses()` 返回不定长，而那会先在 `:16` 响 |
| 11 | `goroutine_test.go:152` | `rep.Unknown[0]` | `len(rep.Unknown) == 1 &&` 短路 | **不改** |
| 12 | `goroutine_test.go:158` | `rep.Unknown[0]` | `len(rep.Unknown) != 1 \|\|` 短路 | **不改** |
| 13 | `logging_test.go:183` | `names[0]`、`names[1]` | `:180` `if len(names) != 2 { t.Fatalf }` | **不改** |
| 14 | `logging_test.go:244` | `files[0]` | `:241` `if len(files) == 0 { t.Fatal }` | **不改** |
| 15 | `logging_test.go:248` | `strings.SplitN(…)[0]` | `SplitN` 对任何输入至少返回 1 段 | **不改**：长度天然安全，加守卫是给读者制造噪音 |
| 16 | `sampler_test.go:238` | `got[0]` | `len(got) != 1 \|\|` 短路 | **不改** |
| 17 | **`sampler_test.go:308`（改前）** | **`rep.Samples[0]`** | **无** | **改**（本程唯一一枚，见 §2.1、读数 §3） |

**同族相邻一件，登记但本程不动**（它不是 `[0]` 直取，超出 AC#8 点名的形状，交验收方定性）：`sampler_test.go:31` 的 `fakeTree.ReadTree()` 在 `f.mu` 为空且 `current==nil && err==nil` 时会取 `f.mu[len(f.mu)-1]` ⇒ 同样是一发越界 panic。实测入口面：包内 `&fakeTree{}` 的三处用法（`sampler_test.go:222`／`:226`／`:233`）要么带 `err` 要么在 `SampleState` 的状态校验处先返回错误、`MarkTransition` 压根不读树（`sampler.go:231-238`），当前无路径能踩到 ⇒ 我没为它造变异、也没改它。

### 2.3 判据②的结案读数：同一发 M3 落在**加了守卫**的树上

树：`/d/tmp/wisp136ac8-ac9-tree1`（`git archive 79ddd49`）。先证落地（同 §1 那一发，逐字同一行、命中数 1、`go build ./...` rc=0）：

```
$ grep -n 'PrivateWorkingSetBytes >= 0' /d/tmp/wisp136ac8-ac9-tree1/internal/observe/sampler.go
290:		} else if m.PrivateWorkingSetBytes >= 0 {
```

`go test -count=1 -v ./internal/observe/`（原文 `/d/tmp/wisp136ac8-ac9-tree1-M3-v.txt`）：

```
rc=1 / === RUN=56 / --- PASS=50 / --- FAIL=6 / SKIP=0 / ^panic:=0
```

- **名册差集**：该发 `=== RUN` 名册与未变异基线名册 `diff` **无输出** ⇒ 没有任何一枚被拖走（改前是 52/56、拖走 4 枚）。
- 红名逐名（6 枚，每枚红在自己身上）：
  1. `TestSampleStateAllMetricsAndVerdicts`
  2. `TestSampleStateSleepingDiskWriteGateFails`
  3. `TestSampleStateWorkPeakMemoryIsTargetNotGate`
  4. `TestSampleStateLeakFixtureFlipsRed`
  5. **`TestSamplerGoroutineAccountingFollowsRegistry` ← 本程改的那枚，红在 `sampler_test.go:315`**
  6. `TestSampleStateTrustworthyWindowNotMarkedUnmeasurable`（AC#1 正向对照腿，红在 `:151`）
- 那枚守卫的原文读数：

```
=== RUN   TestSamplerGoroutineAccountingFollowsRegistry
    sampler_test.go:315: sampler recorded 0 samples (sample_errors=3 last_error="read returned a zero private working set for a live tree"): there is no reading to check goroutine accounting against
--- FAIL: TestSamplerGoroutineAccountingFollowsRegistry (0.03s)
=== RUN   TestLiveRegistryBaselineWithinSleepingGate
--- PASS: TestLiveRegistryBaselineWithinSleepingGate (0.00s)
=== RUN   TestThresholdTableCoversAllStates
```

⚠ **一处措辞边界，如实写，不含糊**：票面②那句"只红这一枚、同包其余用例全部照常跑完"，后半句是**逐字成立**（56 枚全跑到、SKIP=0、`panic=0`）；前半句"只红这一枚"**不成立、也不该成立**——M3 是把丢弃分支改成逢读数都丢，它同时也把 `sampling` 门扇开，所以另外 5 枚依赖"有可信读数"的用例**本来就该红**（改前的同一发在 panic 之前就已经红了 4 枚，§1 有名单）。这一格真正能被判据②钉住的性质是：**"空 Samples"那一发不再由一枚 panic 代答，而是每枚各红各的、名册完整**。我把这条按实测这样报，不把它写成"只红一枚"。

---

## §3 AC#8 判据③：AC#1 那枚钉子的 M3 复跑，现在**整包 `-v`** 就取得到（不再需要定点绕过）

同一次 §2.3 的读数里就有 AC#1 两腿（`/d/tmp/wisp136ac8-ac9-tree1-M3-v.txt` 第 115-119 行原文）：

```
=== RUN   TestSampleStateZeroSampleWindowFailsClosed
--- PASS: TestSampleStateZeroSampleWindowFailsClosed (0.05s)
=== RUN   TestSampleStateTrustworthyWindowNotMarkedUnmeasurable
    sampler_zerosample_136_test.go:151: trustworthy reads must be sampled
--- FAIL: TestSampleStateTrustworthyWindowNotMarkedUnmeasurable (0.06s)
```

⇒ 与验收方用**定点读数**给出的 M3 形状一致（腿 A 绿、腿 B 红在 `sampler_zerosample_136_test.go:151`），差别只在：这次是从**整包 `-v`** 里直接取到的，本程**没有跑过任何 `-run` 定点绕过**来拿这条结论。AC#1 那枚钉子的断言一字未改（`git diff 1d38206..HEAD -- internal/observe/sampler_zerosample_136_test.go` 无输出）。

**还原（三态的第三态）**：restore ⇒ `diff -q` 逐字相同 ⇒ `go test -count=1 -v ./internal/observe/` ＝ `rc=0 / 56 / 56 / 0 / 0`（原文 `/d/tmp/wisp136ac8-ac9-tree1-M3-restore-v.txt`）。

**AC#8 三态齐否**：变异前绿（56/56/0/0）／变异后红（56 RUN、6 FAIL、panic=0、红名逐名）／还原复绿（56/56/0/0）。

---

## §4 AC#9：`CheckSettle` 的零可信样本面装钉（①②③ 逐条对表）

**顺序前提已满足**：AC#8 那枚守卫在 `79ddd49` 就落了，AC#9 的每一次读数都是**整包 `-v`**，本程**没有**为 AC#9 走过任何定点绕过（M11 那一发我两种读数都取了，见 §4.4）。

### 4.1 判据①：新增用例（commit `2f291d0`，净面只有这一枚新文件，153 增 0 删）

`internal/observe/sampler_settle_zerosample_136_test.go`，与 AC#1 那枚钉子同形（一红腿＋一正向对照腿）：

| 腿 | 用例 | 钉住什么 |
| --- | --- | --- |
| A（钉） | `TestCheckSettleZeroTrustworthySamplesFailsClosed` | fixture＝活树（`PIDs:1`）每次读数 `PrivateWorkingSetBytes:0`。①前提腿：读到了树（`:63`）、`len(rep.Samples)==0`（`:66`）；②**归因隔离**：`FreeOSMemoryRequested==true`（`:73`）、`FreeOSMemoryCount>0`（`:76`）、`FinalBytes<=CapBytes`（`:79`）三项**全部满足**，所以 `Pass=false` 只能来自"从未记录过可信读数"这一支（否则会红在这三行而不是钉本身）；③**钉**：`rep.Pass` 必须 false（`:85`）、`BackWithinCapMS` 必须留 `-1` 哨兵（`:88`）；④**"报告要说出自己没测到"**：出线 JSON 里 `samples` 字段在且为空、`pass:false`、`back_within_cap_ms:-1`（`:101-113`） |
| B（正向对照，防恒真） | `TestCheckSettleTrustworthyReadsAreRecorded` | 可信读数（4MB，落在冻结的 Sleeping 25MB 帽之下）必须被记进 `Samples`、`-1` 哨兵必须退场、`FinalBytes` 必须带足迹、窗口必须 `pass=true` |

与已有两枚的分工（不重复记账）：`sampler_test.go:243` 钉 release 计数契约、`:276` 钉"始终够不到帽"，**两枚都没断言过"一次可信读数会不会落进 `Samples`"**——那正是腿 B 的账。腿 A 的判据物（零可信样本 ⇒ 不许 pass ＋ 报告自陈没测到）此前仓内无人认。

**装钉后的基线**（tree2＝`git archive 2f291d0`，未变异，`go test -count=1 -v ./internal/observe/`，原文 `/d/tmp/wisp136ac8-ac9-tree2-baseline-v.txt`）：

```
rc=0 / === RUN=58 / --- PASS=58 / --- FAIL=0 / SKIP=0 / panic=0
```

（58 ＝ 锚点 56 ＋ 本程新增 2 枚。）

### 4.2 判据②：`sampler.go:477` 落 `> 0` ⇒ `>= 0`（M10），新用例转红、红名点到它

**先证落地，再读数**（原文照贴）：

```
$ grep -n 'if err == nil && m.PrivateWorkingSetBytes' /d/tmp/wisp136ac8-ac9-tree2/internal/observe/sampler.go
477:		if err == nil && m.PrivateWorkingSetBytes >= 0 {
$ diff -u <pristine-tree2/sampler.go> <tree2/sampler.go>      # 单行，其余一字未动
-		if err == nil && m.PrivateWorkingSetBytes > 0 {
+		if err == nil && m.PrivateWorkingSetBytes >= 0 {
$ go build ./...   rc=0
```

**读数**（整包 `-v`，原文 `/d/tmp/wisp136ac8-ac9-tree2-M10-v.txt`）：

```
rc=1 / === RUN=58 / --- PASS=57 / --- FAIL=1 / SKIP=0 / panic=0
--- FAIL: TestCheckSettleZeroTrustworthySamplesFailsClosed (0.10s)
```

红名**逐名＝只此一枚**，红点与消息原文：

```
=== RUN   TestCheckSettleZeroTrustworthySamplesFailsClosed
    sampler_settle_zerosample_136_test.go:66: precondition broken: a settle window of zero-footprint reads recorded 5 samples, report=&{TargetState:Sleeping ... CapBytes:26214400 FreeOSMemoryCount:1 FreeOSMemoryRequested:true BackWithinCapMS:20 ElapsedMS:100 FinalBytes:0 Samples:[{... TreePrivateBytes:0 ...}×5] Pass:true}
--- FAIL: TestCheckSettleZeroTrustworthySamplesFailsClosed (0.10s)
=== RUN   TestCheckSettleTrustworthyReadsAreRecorded
--- PASS: TestCheckSettleTrustworthyReadsAreRecorded (0.10s)
```

⇒ 病形与验收方 `R-136-1` 的探针读数同族：一份 `Pass:true`、`BackWithinCapMS:20`、五个 `TreePrivateBytes:0` 的"样本"、`FinalBytes:0`——即"从未取到可信读数的 settle 窗口交出一枚绿"。（验收方那次是 `samples=6 back_within_cap_ms=10`，枚数/毫秒差在 ticker 抖动上，形状一致。）腿 B 同发仍绿 ⇒ 红的归因点到腿 A 自己身上。

**改前对照（这一发我在没有 AC#9 钉的树上复跑了一次）**：tree1（＝`git archive 79ddd49`，只有 AC#8 守卫、无 AC#9 钉）落同一发 M10 ⇒ `RUN=56 / PASS=56 / FAIL=0 / SKIP=0 / rc=0`（复跑两次同数，原文 `/d/tmp/wisp136ac8-ac9-tree1-M10-rerun1-v.txt`、`-rerun2-`）。⇒ 验收方"整包 56/56、rc=0"那句我在自己的锚点系上复算成立；现在这枚红**只由本程新装的那一枚仪器认**。（第一次跑这一发时同包多了一枚 `TestNoopTaskReturnsToBaseline` 红，是一枚既有 flake，见 §4.5。）

### 4.3 判据③：还原 ⇒ 复绿

```
$ python wisp136ac8-ac9-mut.py <tree2> <pristine-tree2> restore   # 每发驱动前先无条件 restore
$ diff -q <pristine-tree2/sampler.go> <tree2/sampler.go>          # 无输出＝逐字相同
$ diff -q <tree2/sampler.go> <repo/internal/observe/sampler.go>   # 无输出＝生产码本程一字未动
$ go test -count=1 -v ./internal/observe/
rc=0 / === RUN=58 / --- PASS=58 / --- FAIL=0 / SKIP=0 / panic=0
```

**AC#9 三态齐**：变异前绿（58/58）／变异后红（58 RUN、1 FAIL、红名点名新用例、红点 `:66`）／还原复绿（58/58）。

### 4.4 腿 B 不哑（自加的一发，M11：让 `CheckSettle` 压根不记 `Samples`）

`sampler.go:486` 的 `rep.Samples = append(rep.Samples, sm)` ⇒ `_ = sm`（落地证明：`486:			_ = sm` ＋ `go build ./...` rc=0）：

```
整包：rc=1 / RUN=58 / PASS=57 / FAIL=1 / SKIP=0 / panic=0
--- FAIL: TestCheckSettleTrustworthyReadsAreRecorded   红在 sampler_settle_zerosample_136_test.go:142
TestCheckSettleZeroTrustworthySamplesFailsClosed 同发仍 PASS
```

⇒ 两腿各响各的：M10（放宽判断）只打死腿 A，M11（不记录）只打死腿 B。这条不是判据②要的那发，是防"腿 A 的 `len(Samples)==0` 断言写成恒真"的旁证。还原后 58/58 复绿（`/d/tmp/wisp136ac8-ac9-tree2-M11-restore-v.txt`）。

### 4.5 顺带量到的一枚既有 flake（登记，不归本程修、不影响任何判据）

`TestNoopTaskReturnsToBaseline`（`internal/observe/goroutine_test.go:33`）在本程窗口内命中 **1 次**：

```
--- FAIL: TestNoopTaskReturnsToBaseline (0.00s)
    goroutine_test.go:33: PerTask mid-task = 2, want 3
```

出处：tree1＋M10 的第一发整包读数（`/d/tmp/wisp136ac8-ac9-tree1-M10-nailless-v.txt`）。同一棵树同一发复跑两次＝56/56 rc=0；未变异的 tree1 连跑 5 次＝5/5 全绿（`/d/tmp/wisp136ac8-ac9-tree1-flake1..5.txt`）⇒ 观测频次 1/9。成因读码可辨：`reg.Spawn` 后立刻读 `RosterReport().PerTask`，未获调度则少计。**它不是本程造的**（本程未碰 `goroutine_test.go`，`git diff 1d38206..HEAD -- internal/observe/goroutine_test.go` 无输出），但它是**下一位读数的人会撞到的东西**，且它落在"吞读数判据"的同一族里（一枚仪器偶发把包级 rc 判红），值得单立一格——本程不动它。

---

## §5 门禁读数（票面 AC#7 同款，原文照贴；跑在 `2f291d0` 之后的工作树上，此后本程只再提交 `.md`，不影响这些读数）

### G1 `gofmt -l` 整包

```
$ gofmt -l internal/observe/
（无输出）                      rc=0
```

### G2 `gofumpt -l`（版本写明）

```
$ /d/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)
$ /d/work/base/gopath/bin/gofumpt.exe -l internal/observe/
（无输出）                      rc=0
```

⚠ 边界：CI 那一步是 `go install mvdan.cc/gofumpt@latest`（`ci.yml:113`，票 122 的 AC#7 已把这版本未钉落成单独一格），**本程没有为了跑门禁去重装/升级工具链**（PATH 里根本没有 `gofumpt`，只有上面那个二进制），所以这条读数**只在 v0.12.0 下成立**，不等于 CI 那一步的读数。

### G3 `go vet` 双 GOOS（逐错误行归因，不整树 rc=1 就甩给工具链）

| 发 | 命令 | rc | 读数与归因 |
| --- | --- | --- | --- |
| 3a | `go vet ./internal/observe/` | 0 | 零输出 |
| 3b | `GOOS=linux go vet ./internal/observe/` | 0 | 零输出 ⇒ **本程的改动面在 linux 交叉下清白** |
| 3c | `go vet ./...`（原生全仓） | 0 | 零输出 |
| 3d | `GOOS=linux go vet ./...`（交叉全仓） | **1** | 全文 3 行、**只有 1 条诊断**：`package github.com/CarlosShao/wisp/cmd/wisp` → `imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx` → `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8`。落在 `cmd/wisp` 的 cgo 依赖加载上，**既不算破口也不算清白**（它停在那一格，后面的包这一发没答） |
| 3e | 从 3d 的名单里剔掉 `cmd/wisp` 再跑（31 枚包） | **1** | 全文 1 行、**只有 1 条新诊断**：`package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in …\cmd\balldebug`（宿主 GOOS=windows 下该包的 linux 构建约束排除全部文件）。⇒ 逐错误行归因完成，仍是 `cmd/` 那一格 |
| 3f | 剔掉两枚 `cmd/` 包后剩下的 **30 枚**（`internal/**`＋`frontend`＋`tools/signmodels`） | **0** | 零输出 ⇒ **30/30 非 cmd 包在 linux 交叉下清白**（原文 `/d/tmp/wisp136ac8-ac9-vet-linux-nocmds.txt` 空、名单 `/d/tmp/wisp136ac8-ac9-pkgs-nocmd.txt` 30 行） |

⚠ 派单里提到的"交叉那一发停在 cgo 包加载"＝这里的 3d/3e 两枚，形状对上了；但派单若指的是 `internal/risk/pathresolver`，**本程实测未读到**（3f 含 `internal/risk`，零输出）。这条按"前提与实测的差异"报回，不据此对别的树下结论。

### G4 `go test -count=2 -v ./internal/observe/` 四数（改前／改后各一份）

| | 树 | rc | `=== RUN` | `--- PASS` | `--- FAIL` | SKIP |
| --- | --- | --- | --- | --- | --- | --- |
| 改前 | tree0（`git archive 1d38206`，锚点） | 0 | **112** | **112** | 0 | **0** |
| 改后 | 本程工作树（AC#8 守卫＋AC#9 钉都在） | 0 | **116** | **116** | 0 | **0** |

逐名账（两发各自 `=== RUN` 去重后 `diff`，原文 `/d/tmp/wisp136ac8-ac9-names-before.txt` / `-after.txt`）：

```
2a3
> TestCheckSettleTrustworthyReadsAreRecorded
3a5
> TestCheckSettleZeroTrustworthySamplesFailsClosed
```

⇒ 112→116 的 ＋4 **全部**是"本程新增 2 枚 × 2 轮"，无既有枚改名、无消失、无转 SKIP（SKIP 两发都是 0）。AC#8 那枚改动不新增名（改的是既有用例），这一点单独说明，免得被当成漏账。

### G5 `sh scripts/d22scan.sh`（票面 AC#7 同款，顺手取了改前分母）

| | 树 | rc | `ban #8 internal/` | `ban #8 cmd/` | `ban #8 frontend/` | `ban #8 design/` |
| --- | --- | --- | --- | --- | --- | --- |
| 改前 | tree0 @1d38206 | 0 | **401** | 39 | 40 | 16 |
| 改后 | 本程工作树 | 0 | **402** | 39 | 43 | 16 |

⇒ `internal/` 401→402 ＝ 本程新增的那一枚 `_test.go`，**各 scope 无一下降**；两形都 `d22scan: clean - no D22 ban violations`。`frontend/` 的 40→43 **不是兄弟在飞的树**：`git diff --name-only 1d38206..HEAD -- frontend/` **无输出**，那 3 枚是工作树里未跟踪的构建产物 `frontend/dist/index.html`＋`frontend/dist/assets/index-*.css`＋`index-*.js`（我用两棵树的 `find` 差集逐名核过）。

---

## §6 我未做的档（诚实清单，别当已验）

1. **linux 全仓两态**没跑。本程的 linux 面只有宿主的 `GOOS=linux go vet`（G3），**没有** `golang:1.27` 容器里的 `go test`。AC#8/AC#9 的读数全是 windows/amd64 宿主的。
2. **CI run id 未取**：本程只 commit 未 push；"门禁上有一枚红过的 run"是 AC#7 的账，不在本格。
3. **端到端 `wisp slo`（AC#10）一字未碰**：`cmd/wisp/**` 对本程是整块禁区。
4. **AC#2..AC#7 一字未碰**，也没替它们做任何判定。
5. ⚠ **生产码零改动 ⇒ 我没有给 `SettleReport` 新增一枚像 `sampling` 那样的"没测到"门行**。票面 AC#9① 那句"报告要说出自己没测到"，我钉的是**现有哨兵**（空 `samples` ＋ `back_within_cap_ms=-1` ＋ `pass=false`）而不是新造字段。若验收方裁"必须像 `SampleState` 那样在报告里产出一枚显式 gate 行"，那需要动 `sampler.go:477-501`（生产码），本程按派单（"大概率不需要动生产码；若需要，先报回"）**没有动、在此报回**，等裁决。我的实测是：pristine 下这面确实 fail-closed（§4.1 基线 58/58、腿 A 全绿），M10 只有新装的这一枚认（§4.2、含改前 56/56 rc=0 对照）。
6. `-race`、`-shuffle` 未测。AC#8 的"名册差集"是**当前执行顺序**下的账（同验收方 §8.5 的边界）。
7. M10 族其余写法（`!= 0`、条件反写、把 `err == nil` 那半句摘掉）未打——派单只要那一发。M11 是自加的正向对照自校，不在判据上。
8. **AC#1 的 M1/M2/M4/M6/M7/M8 六发未复算**（不是本格账，也没推翻任何结论）。
9. `TestNoopTaskReturnsToBaseline` 那枚既有 flake（§4.5）未修、未造变异，只登记了 1/9 的观测频次。
10. `gofumpt` 未用 CI 的 `@latest` 跑（G2 的边界），`staticcheck`／其他 linter 未跑。
11. 临时件**只建不删**：tree0/tree1/tree2、两份 pristine 目录、驱动器与全部读数文件都在 `/d/tmp/`（`wisp136ac8-ac9-*`）。`/d/tmp/wisp136-acc-r1-tree2/`（验收方凭据）未动、未删。
12. 翻勾：一枚未翻。AC#8／AC#9 的 `[ ]` 原样（终裁归非实现者）。

---

## §7 本程 commit 清单（每次暂存只有我自己的路径）

| sha | 文件 | 净面 |
| --- | --- | --- |
| `79ddd49` | `internal/observe/sampler_test.go` | 10 增 0 删（AC#8② 那枚守卫） |
| `36443f2` | `docs/evidence/s1/136-ac8-ac9-impl.md` | 本文件 §0-§3 |
| `595abd3` | `.scratch/wisp/issues/136-….md` | 票面 log 追加 18 行（AC#8，append-only，勾未翻） |
| `2f291d0` | `internal/observe/sampler_settle_zerosample_136_test.go` | 153 增 0 删（AC#9① 两腿） |
| （本枚起） | `docs/evidence/s1/136-ac8-ac9-impl.md` ＋票面 log | §4-§7 ＋ AC#9 log 一条 |

每枚 `git add -- <显式路径>`、`git diff --cached --name-only` 只出现上述路径；`git commit -q -F - -- <显式路径>`；**只 commit 未 push**；未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；仓内未建 worktree。工作树里那枚未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` 全程未读、未提交、未改、未删。

**next=** 派**非实现者**按票面 AC#8②③／AC#9①②③ 终裁并翻勾（可重算凭据：M3 与 M10 各一发、两棵纯净树、`/d/tmp/wisp136ac8-ac9-*` 读数文件只建不删）。另请裁两件本程主动留下的事：①§6.5——AC#9① 那句"报告要说出自己没测到"钉在**现有哨兵**上够不够，还是要动生产码给 `SettleReport` 补一枚 `sampling` 同形门行；②§4.5——`TestNoopTaskReturnsToBaseline` 这枚既有 flake 是否单立一格。AC#10 仍排 `cmd/wisp` 空出来之后。


