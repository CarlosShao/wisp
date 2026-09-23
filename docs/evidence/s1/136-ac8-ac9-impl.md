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
