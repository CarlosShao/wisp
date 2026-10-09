# 287-a1 — 票 287 三把尺复跑读数（只读普查腿）

**腿**：`287-a1`（只读，非编排者，不裁任何格、不翻任何 `- [ ]` 框）
**跑的票**：票 287 `AC#1`（唯一归我的那一格；`AC#2`/`AC#3`/`AC#4` 不归我，未碰）
**时间**：2026-10-09 12:1x +08
**锚点**：开工时 `git rev-parse --short HEAD` = `fc359af8`，收工复量时已漂到 `9d00ceb2`（共享工作树，别的腿在提）。
  ⇒ 本件里**每一枚行号都在 `9d00ceb2` 上重新逐字复量过一遍**，两枚锚点下 `internal/tools/subagent_197.go` 的关键行**逐字一致**（见尺③第 0 步）。

**本腿的读面**：`.go` / `.sh` / `.ps1` / `.yml` 零写、`frontend/**` 与 `design/**` 零读零引、零 `go build` / `go vet` / `go test`、零 push。
写点只有本件一枚。

---

## 0. 读数环境的一处前处理（不做这一步，尺①会读到假数）

`.scratch/**` 里躺着**历史探针的产码副本**，其中包括票 222 的**变异件**
`.scratch/.scratch/wisp/probes/222/v1/mutations/asis/subagent_197.go` —— 它带着同一枚拒句，
且它的行号（`:257-258`）与真产码（`:269-270`）**不一样**。
所以下面每一条 grep 都挂了 `grep -v "^./.scratch"`，枚数是**剔掉 `.scratch` 之后**的枚数。

命令形状（统一）：

```sh
grep -rn "<模式>" --include=*.go . | grep -v "^\./\.scratch"
```

---

## 1. 尺① —— 拒句在 `*_test.go` 里命中几枚

### 1.1 先在产码里找"真正的拒句原文"（不照抄票面措辞）

派单提醒票面那句「子代理不许再派子代理」**可能已过期**。现扫四条语义变体，逐条给命中：

```sh
for pat in "子代理不许再派子代理" "子代理不能再派生子代理" "task.spawn 也不在你的工具目录里" "拒绝执行：子代理"; do
  t=$(grep -rn "$pat" --include=*_test.go . | grep -v "^\./\.scratch" | wc -l)
  p=$(grep -rn "$pat" --include=*.go . | grep -v "_test\.go" | grep -v "^\./\.scratch" | wc -l)
  echo "[$pat] test_lines=$t prod_lines=$p"
done
```

真实读数（逐字照抄输出）：

```
PATTERN[子代理不许再派子代理] test_files=0 test_lines=0 prod_lines=1
PATTERN[子代理不能再派生子代理] test_files=0 test_lines=0 prod_lines=2
PATTERN[task.spawn 也不在你的工具目录里] test_files=0 test_lines=0 prod_lines=1
PATTERN[拒绝执行：子代理] test_files=0 test_lines=0 prod_lines=1
```

命中的产码位置逐字（`sed -n` 抽的，不是引用）：

| 行 | 逐字原文 |
|---|---|
| `internal/tools/subagent_197.go:270` | `				"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil` |
| `internal/tools/subagent_197.go:207` | `		"子代理不能再派生子代理（深度 %d），同时在跑的子代理上限 %d 枚",` |
| `internal/tools/subagent_197.go:576` | `			Text: fmt.Sprintf("拒绝执行：子代理不能再派生子代理（深度上限 %d），"+` |
| `internal/tools/subagent_197.go:577` | `				"task.spawn 也不在你的工具目录里。", MaxSubagentDepth),` |

⇒ **票面措辞没有过期**：「子代理不许再派子代理」在产码里**逐字存在**，就在 `:270`，
也就是**深度判定那一支的自己肚子里**（`:267-271`）。它另外还是 `:207`（Description 文本）与 `:576-577`（provider 硬拒）两枚**同义不同字**的句子。

### 1.2 判定

**成立**。票面 ① 说该拒句在 `*_test.go` 命中 **0** —— 现量 `test_lines=0`，四条变体全部 0 命中，
非测试件命中 1 / 2 / 1 / 1。复现。

### 1.3 但这把尺的**证据力**要在件里说清楚（不是推翻票面，是划清它能证什么、不能证什么）

"0 命中" ≠ "孩子派孙这件事没有任何仪器"。测试里**确实**有一枚在管这件事，只是它**不按句子断言**：

```
internal/tools/subagent_197_test.go:661:	if !out.IsError || !strings.Contains(out.Text, "深度") {
```

`Test197ChildCannotDeriveSubagent`（`:631`）断言的是 **provider 那一发**（`:640` 拿 `newSubagentToolProvider(h.dir)`、`:654` 直接 `childDir.Execute` 发 `task.spawn`）
返回的文本里**含 "深度" 这两个字**——它既没钉 `:270` 那句，也没钉 `:576` 那句整句。
所以：**尺① 的字面读数成立，但它只排除"`:270`/`:576` 这两句句子没被钉"，排除不了"拒绝行为有仪器"**；
真正的"没仪器"结论要靠尺③的**调用形状**来支撑（下面尺③给）。

另外 `:647`、`:677`、`:895`、`:430` 与 `task_cancel_221_legs_test.go:266` 也提到"深度"，逐条读过，都不是钉拒句的：
`:647` 是"孩子目录里不该有 task.spawn"，`:677` 是"名册里不许出现第二代"，`:895` 是常量漂不漂，
`:430` / `:266` 断言的是 **Description 文本**含 `深度 %d`，不是拒答文本。

---

## 2. 尺② —— `TaskKindSubagent` 逐枚读法

### 2.1 全部命中

```sh
grep -rn "TaskKindSubagent" --include=*.go . | grep -v "^\./\.scratch"
```

真实读数：**18 行**（剔 `.scratch` 后）。同一枚命令下的三枚计数：

```
TOTAL_LINES=18
TEST_LINES=9      # grep "_test.go:"
PROD_LINES=9      # grep -v "_test.go:"
```

分成三堆：

- **产码 9 行**：`cmd/wisp/panel_pump.go:161`、`internal/panel/subagent_roster_197.go:64`（注释）、
  `internal/tools/subagent_197.go:55`（注释）、`:58`（常量定义 `TaskKindSubagent = "subagent"`）、
  `:267`（**深度判定本身**）、`:429`（孩子收口落行）、
  `internal/tools/task.go:137`（注释）、`:344`（`PublishSubagent` 落行）、`:410`（`RunningSubagentIDs` 筛行）。
- **测试 9 行**（下表；按测试文件分布：`cmd/wisp/subagent_carrier_197_test.go` 4、`internal/tools/subagent_197_test.go` 2、`internal/tools/task_cancel_221_legs_test.go` 1、`internal/tools/ticket283_corr_identity_rulers_test.go` 2）。
- 上面两堆没有重叠：18 = 9 产码行 + 9 测试行（其中 `subagent_197.go:55`、`:58`、`task.go:137`、`internal/panel/subagent_roster_197.go:64` 是注释与常量定义，不参与"断言什么"的问答）。

### 2.2 测试命中逐枚：它断言的是**孩子那一行**，还是把**调用者那一行**摆成 subagent？

（"断言位点"按同一枚 `if` + 它的 `t.Errorf` 归一发，因为 `t.Errorf` 行只是那枚断言的报错文本。）

| # | 文件:行 | 逐字原文（sed -n 抽） | 现答 |
|---|---|---|---|
| 1 | `cmd/wisp/subagent_carrier_197_test.go:218` | `		if row.Kind == tools.TaskKindSubagent && row.ParentTaskID == rootID {` | **孩子那一行**（`childRowAfterJoin197` 在名册里挑 root 下面的孩子） |
| 2 | `cmd/wisp/subagent_carrier_197_test.go:291` | `	if child.Kind != tools.TaskKindSubagent {` | **孩子那一行**（`child := carrierRow(t, sect, childID)`） |
| 2' | `cmd/wisp/subagent_carrier_197_test.go:292` | `		t.Errorf("kind = %q, want %q as the spawner filed it", child.Kind, tools.TaskKindSubagent)` | 同上那枚的报错文本 |
| 3 | `cmd/wisp/subagent_carrier_197_test.go:570` | `		if row.Kind != tools.TaskKindSubagent {` | **孩子那一行**（`for _, row := range sect.Rows` 的筛子循环，`continue` 掉非孩子） |
| 4 | `internal/tools/subagent_197_test.go:305` | `	if row.Out.Kind != TaskKindSubagent {` | **孩子那一行**（`row := h.childRow(t)`，即 `Descendants(parent197)`） |
| 4' | `internal/tools/subagent_197_test.go:306` | `		t.Errorf("Kind = %q, want %q", row.Out.Kind, TaskKindSubagent)` | 同上那枚的报错文本 |
| 5 | `internal/tools/task_cancel_221_legs_test.go:333` | `	if rec.Out.ParentTaskID != parent197 || rec.Out.Kind != TaskKindSubagent \|\| rec.Out.Label != label197 {` | **孩子那一行**（`rec := x.h.mustLook(t, childID)`） |
| 6 | `internal/tools/ticket283_corr_identity_rulers_test.go:197` | `	if rec.Kind != TaskKindSubagent {` | **孩子那一行**（`rec, ok := x.h.roster.Look(childID)`，上面 `:189` 逐字可查） |
| 6' | `internal/tools/ticket283_corr_identity_rulers_test.go:198` | `		t.Errorf("Kind = %q, want %q", rec.Kind, TaskKindSubagent)` | 同上那枚的报错文本 |

**枚数现量**：测试**行** 9 枚（票面 ② 写的是"5 枚测试命中"）；按**断言位点**归并是 **6 枚**。
⇒ **无论按哪种数法都对不上票面的 5**。派单已经预告"枚数会漂"，这一枚就是漂了。

**逐枚读法结论**：**6 枚全部断言"孩子那一行"的 Kind，没有一枚把调用者那一行摆成 subagent。**
⇒ 票面 ② 的**实质成立、枚数不成立**（部分成立）。

### 2.3 支撑上面那句的两条额外现量（不是推理，是查出来的）

调用者行在现有夹具里**永远是 root**，这是"没有一枚能走到 `:267`"的直接原因：

- `internal/tools/subagent_197_test.go:163` 逐字 `	h.roster.MarkRoot(parent197, "根任务")`
- `internal/tools/subagent_222_test.go:261` 逐字 `	h.roster.MarkRoot(parent222, "根任务 222")`
- `internal/tools/task.go:356`（`MarkRoot` 本体）逐字 `	r.Record(taskID, TaskOutput{Label: label, Kind: TaskKindRoot})`

而测试里仅有的三枚 `PublishSubagent`（能把行落成 subagent 的那枚方法）摆的都不是调用者：

- `internal/tools/subagent_197_test.go:739` 逐字 `	h.roster.PublishSubagent("sibling-197", parent197, "隔壁那枚")` —— 兄弟
- `internal/tools/task_cancel_221_legs_test.go:298` 逐字 `	x.h.roster.PublishSubagent("sibling-221", parent197, "隔壁那枚")` —— 兄弟
- `internal/tools/task_cancel_221_legs_test.go:521` 逐字 `	x.h.roster.PublishSubagent("kid-221", parent197, "只登记不起环的孩子")` —— 孩子，只登记

---

## 3. 尺③ —— 调用形状尺：`:267` 到不到得了

### 第 0 步：先钉"工作树 == HEAD"，否则行号尺无意义

```sh
git status --porcelain -- internal cmd        # 读数：空（无输出）
git diff --quiet HEAD -- internal/tools/subagent_197.go && echo "IDENTICAL" || echo "DIFFERS"
# 读数：IDENTICAL（收工时 12:1x 在 9d00ceb2 上又跑了一遍同一枚命令，读数同）
```

### 第 1 步：深度判定到底在哪几行（票面标题写 `:267-271`）

```sh
git show HEAD:internal/tools/subagent_197.go | sed -n '259,271p'   # 逐字，制表符照盘上
```

```
259		// The task-level id of the caller: since ticket 242 the correlation id is
260		// minted per call (C18 routes approval replies by it), so it can no longer
261		// name a task; the roster is keyed by task ids.
262		parentID := TaskID(ctx)
263		if parentID == "" {
264			return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：不知道父任务是谁，" +
265				"就没法登记父子关系，也没法把结论盖戳进父任务的 C25 作用域）", IsError: true}, nil
266		}
267		if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {
268			return Result{Text: fmt.Sprintf(
269				"拒绝派生：深度上限是 %d，而派生者 %s 自己就是一枚子代理（父 kind=%s）。"+
270					"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil
271		}
```

- **票面标题 `:267-271` vs 我现量：一致。** 判定入口行 `:267`，块尾 `:271`，拒语体 `:268-270`。
- **票面 `:8` 说 `parentID := TaskID(ctx)` 在 `:262` vs 我现量：一致**（`:262` 逐字 `	parentID := TaskID(ctx)`；`TaskID` 本体在 `internal/tools/cancel.go:67` 逐字 `func TaskID(ctx context.Context) string {`）。
- 这段的宿主函数：`internal/tools/subagent_197.go:238` 逐字 `func (t subagentSpawn) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {`

### 第 2 步：它的入口是谁调的

`subagentSpawn.Execute` 是 C4 工具面，生产上只由桥发到（`internal/tools/bridge.go`），
而桥只被环路的 dispatch 调：

```sh
grep -rn "\.Execute(ctx" --include=*.go internal/agent | grep -v "_test.go"
```

读数里那一发：`internal/agent/loop.go:760` 逐字 `		return l.opt.Tools.Execute(ctx, req)`（dispatch 的超时分支另在 `:764`，同一枚 `l.opt.Tools`）。

⇒ **能不能到 `:267`，取决于孩子的 `l.opt.Tools` 是哪一枚 provider。**

### 第 3 步：孩子那枚 provider 是什么（票面 ③ 猜在 `subagentToolProvider.Execute`，`:573-580` 一带）

```sh
git show HEAD:internal/tools/subagent_197.go | awk 'NR==296||NR==553||NR>=573&&NR<=582{printf "%d: %s\n", NR, $0}'
```

```
296: 	opt.Tools = newSubagentToolProvider(t.d.ParentTools)
553: const subagentHiddenTool = "task.spawn"
573: func (p *subagentToolProvider) Execute(ctx context.Context, req agent.ToolRequest) (agent.ToolOutcome, error) {
574: 	if req.Name == subagentHiddenTool {
575: 		return agent.ToolOutcome{
576: 			Text: fmt.Sprintf("拒绝执行：子代理不能再派生子代理（深度上限 %d），"+
577: 				"task.spawn 也不在你的工具目录里。", MaxSubagentDepth),
578: 			IsError: true,
579: 		}, nil
580: 	}
581: 	return p.inner.Execute(ctx, req)
582: }
```

- **票面 ③ 的"`:573-580` 一带" vs 我现量：一致**，且更精确：函数签名 `:573`、名字闸 `:574`、拒语 `:576-577`、`IsError` `:578`、返回 `:579`、块尾 `:580`、**只有出了这个 if 才走 `:581` 的 `p.inner.Execute`**。
- 孩子的目录还**看不见**这个名字：`internal/tools/subagent_197.go:565` 逐字 `		if info.Name == subagentHiddenTool {`（`:566` `continue`）。
- 装配点：`:296` 把孩子的 `opt.Tools` 换成 `newSubagentToolProvider(t.d.ParentTools)`；生产上 `ParentTools` 就是宿主那枚桥（`cmd/wisp/run.go:782` 逐字 `		ParentTools: rt.bridge,`）。

⇒ **孩子的 `task.spawn` 在 `:574` 就被 `:575-579` 硬拒，`:581` 不执行 ⇒ 发不到桥 ⇒ 进不了 `subagentSpawn.Execute` ⇒ 到不了 `:267`。** 票面 ③ **成立**。

### 第 4 步：反向也堵上（"有没有别的调用者本来就是 subagent"）

生产上名册里 `Kind == subagent` 的行**只可能是孩子**，逐字：

- `internal/tools/task.go:344` `		ParentTaskID: parentTaskID, Label: label, Kind: TaskKindSubagent,`
- `internal/tools/subagent_197.go:313` `		t.d.Roster.PublishSubagent(taskID, parentID, label)`
- `internal/tools/subagent_197.go:348` `	t.d.Roster.PublishSubagent(bg.ID, parentID, label)`
- `internal/tools/subagent_197.go:429` `		ParentTaskID: parentID, Label: label, Kind: TaskKindSubagent,`

宿主自己那一行**一律 root**：`internal/tools/task.go:356`
`	r.Record(taskID, TaskOutput{Label: label, Kind: TaskKindRoot})`，发起点 `cmd/wisp/run.go:1105`
`	rt.tasks.MarkRoot(bg.ID, task)`。
⇒ 唯一能触发 `:267` 的调用者（一枚 Kind=subagent 的行）**在生产链上拿不到 `task.spawn`**。

现量测试侧的全部 `task.spawn` 发起点，逐枚看它把谁的 id 摆成调用者：

| 文件:行 | 调用者 id 逐字 | 那一行是什么 |
|---|---|---|
| `internal/tools/subagent_197_test.go:217` | `		TaskID: parent197, CorrelationID: parent197, CallID: "call-197-spawn-" + label,` | root（`:163` MarkRoot） |
| `internal/tools/subagent_222_test.go:332` | `				TaskID: parent222, CorrelationID: parent222, CallID: fmt.Sprintf("call-222-parent-%d", i),` | root（`:261` MarkRoot） |
| `internal/tools/failclosed_236_teeth_test.go:288` | `		TaskID: "", CorrelationID: "", CallID: "call-236r2-spawn-no-parent",` | 空 id ⇒ 撞 `:263` 那一支，不是 `:267` |
| `cmd/wisp/subagent_carrier_197_test.go:173` | `			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),` | root |
| `internal/tools/subagent_197_test.go:655` | `		TaskID: "child-197", Name: "task.spawn",` | 这枚**是** child 身份，但它发的是 `childDir.Execute`（`:654`，provider），**在 `:574` 就被拦**，进不了 `:267` |

⇒ **零枚测试把"调用者那一行"摆成 subagent 再发到真 spawn 工具**；`:655` 是唯一以孩子身份发 `task.spawn` 的一发，但它落在 provider 那一层。
`:267-271` 既无仪器也无后果 —— **票面前提复现，AC#1 要求的三把尺全部回来。**

### 第 5 步：产码自己怎么解释这两半（只报，不裁，这条给 AC#3 用）

`internal/tools/subagent_197.go:87-90` 逐字：

```
	// MaxSubagentDepth is how deep the tree may go: 1, i.e. a subagent may not
	// derive a subagent. Enforced structurally (the child's tool directory has no
	// task.spawn in it) and again by the roster check in Execute, because
	// "reachable but refused" is a weaker guarantee than "not offered".
```

`cmd/wisp/run.go:774-775` 逐字：`		// (the spawner filters task.spawn out of it, so a child never even sees the` / `		// name - depth 1 by structure, not by counting), the same C25 engine, the`

⇒ 产码注释**把 `:267` 明写成"第二道"（again by the roster check in Execute）**。这条与本腿尺③的读数同向，
但"是不是设计如此"是票面 `AC#3` 的问题，**本腿不裁**。

---

## 4. 顺手量的一枚名册：深度／代数在产码里的落点（只报名册，⛔ 不裁夹具挂哪）

尺＝`grep -rn "深度\|代数\|generation\|Generation\|depth\|Depth" --include=*.go internal/tools internal/agent | grep -v "_test.go:"`。

### 4.1 真·子代理树深度（AC#2 用得上的那一堆）

| 文件:行 | 一句逐字引文 |
|---|---|
| `internal/tools/subagent_197.go:41` | `// Subagent identity, pool and depth constants. They live in this one file on` |
| `internal/tools/subagent_197.go:87` | `	// MaxSubagentDepth is how deep the tree may go: 1, i.e. a subagent may not` |
| `internal/tools/subagent_197.go:91` | `	MaxSubagentDepth = 1` |
| `internal/tools/subagent_197.go:207` | `		"子代理不能再派生子代理（深度 %d），同时在跑的子代理上限 %d 枚",` |
| `internal/tools/subagent_197.go:269` | `			"拒绝派生：深度上限是 %d，而派生者 %s 自己就是一枚子代理（父 kind=%s）。"+` |
| `internal/tools/subagent_197.go:270` | `				"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil` |
| `internal/tools/subagent_197.go:544` | `// subagentToolProvider is the structural half of MaxSubagentDepth: the child's` |
| `internal/tools/subagent_197.go:576` | `			Text: fmt.Sprintf("拒绝执行：子代理不能再派生子代理（深度上限 %d），"+` |
| `internal/tools/task.go:137-138` | `	// Kind is TaskKindRoot or TaskKindSubagent (see subagent_197.go), and the` / `	// spawner's depth check reads it: a subagent may not derive another.` |
| `internal/tools/task.go:635` | `		"子代理停兄弟、停自己都一律被拒（深度上限 %d 的树里，别人才不是你的孩子）；"+` |
| `internal/tools/task.go:638` | `		MaxSubagentDepth)` |

枚数＝**11 处**（`subagent_197.go` 8 处 + `task.go` 3 处）。

### 4.2 `internal/agent/**` 子代理深度落点枚数＝**0**

现量：`grep -n "subagent\|Subagent\|depth\|Depth" internal/agent/loop.go` ⇒ **空输出**；
`grep -rn "Subagent\|subagent" --include=*.go internal/agent | grep -v "_test.go"` ⇒ **空输出**。
⇒ 深度这套机器**整枚住在 `internal/tools/**`**，`internal/agent/**` 一点没有。

### 4.3 必须剔掉的**同名异物**（这些"depth"是审批队列的徽标深度，与树深度无关）

`internal/agent/approval/` 下 15 处：`gate.go:167`、`gate.go:535`、`pending_read.go:54`、
`queue.go:74`、`queue.go:146`、`queue.go:147`、`queue.go:148`（`func (q *Queue) Depth() int {`）、`queue.go:181`、
`queue.go:282`、`queue.go:577`、`queue.go:592`、`replies.go:304`、`replies.go:311`、`ui.go:41`、`ui.go:56`。
外加 `internal/agent/prompt.go:384` `// 2 and 3 are deliberately redundant defense-in-depth, and they are NOT`（讲的是 token 预算）。
⇒ 挂夹具时**别顺着这些名字摸进 approval 包**。

### 4.4 测试侧的深度落点（AC#2 要避让／要并置的既有钉子）

`internal/tools/subagent_197_test.go:430`、`:647`、`:661`、`:677`、`:895`；`internal/tools/task_cancel_221_legs_test.go:266`。共 6 枚。

---

## 5. 复跑后与票面**不符**的地方（本票最有价值的一栏，逐条具名）

1. **尺② 的枚数不符**：票面 ② 说 `TaskKindSubagent` 的测试命中是 **5 枚**；我现量 **9 个测试行 / 6 个断言位点**。
   两种数法都不是 5。⇒ 派单预告过"枚数会漂"，这条就是漂了；**实质（"全部断言孩子那行"）不变**。
2. **尺① 的字面成立但证据力比票面说法窄**：`_test.go` 命中 0 复现了，可是 `subagent_197_test.go:661`
   确实在断言 provider 那发的拒答文本（按 `"深度"` 子串）。⇒ "0 命中"证的是"**句子**没被钉"，
   不证"拒绝行为没仪器"；票面把 ① 当"不可达"的证据之一，**这一格里真正承重的是尺③**，不是尺①。
   ⛔ 这不是推翻票面（前提仍由尺③独立成立），但票面 ① 的**措辞**（"三把尺证明它今天不可达"）里，①只证明了"句子无钉"。
3. **票面措辞没有过期**（派单留的疑点，现答）：「子代理不许再派子代理」逐字在 `internal/tools/subagent_197.go:270`。
4. **锚点在我跑期间漂了**：开工 `fc359af8`、收工 `9d00ceb2`；关键行两锚点逐字一致，`git status --porcelain -- internal cmd` 全程空 ⇒ 读数未被脏工作树污染。
5. **票面 `:11` 的"并发 8"与现产码不符**（顺手读到，不在我的尺里，⛔ 不裁）：
   `internal/tools/subagent_197.go:85` 逐字 `	MaxConcurrentSubagents = 4`，
   `internal/tools/bridge.go:24` 逐字 `const MaxToolConcurrency = 4`；
   且 `subagent_197_test.go:446-447` 明写"8"正是当年绕过桥的那枚旧腿留下的错数。
   ⇒ 若 `A399/A401` 的原文真写着"并发 8"，那是**契约文字与常量已分家**的另一枚缺口，归属不在本票，报给编排者定。

---

## 6. 本腿没做、也不该由本腿做的事

- ⛔ `AC#2` 的夹具：一行未造、`go test` 未跑（编排者此刻占着 Go 编译/测试面）。
- ⛔ `AC#3` 的代价表、⛔ `AC#4` 的门禁：未动。
- ⛔ 未翻票 287 任何 `- [ ]` 框；⛔ 未裁任何格。
- ⛔ 零 `frontend/**`、零 `design/**` 读引。

**待人拍板的一条（未定义即停，本腿不自行假设）**：尺② 的"枚数"该按**行**（10）还是按**断言位点**（6）记，
票面用的是"5 枚"这种没定义计数单位的写法。⇒ 请编排者定一枚口径，再决定 `AC#1` 那格翻不翻。
