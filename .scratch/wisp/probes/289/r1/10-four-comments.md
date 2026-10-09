# 票 289 / 腿 289-r1 — 四句注释改前改后逐字 ＋ 两把行中性尺

> 本格只改注释行，零语义改动。改前逐字取自 `git show HEAD:<path> | sed -n '<N>p'`
> （HEAD＝`fd692f7e`，见 `00-anchor-and-baseline.md` §1–§2，逐字命中票面现量表）；
> 改后逐字取自工作树 `sed -n '<N>p'`。每枚件的替换都是**同枚数换行**，所以行号不漂——
> 下面「锚位复量」那一节把四枚原行号重新现跑了一遍。

## 件 1 — `internal/panel/subagent_roster_197.go`（产码件；改 :56-58，3 行 → 3 行）

改前逐字（`:56`／`:57`／`:58`）：

```
	// The entity leg files the deriving call's correlation id, and the agent
	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),
	// so on the production path this reads as the parent's task id.
```

改后逐字（`:56`／`:57`／`:58`）：

```
	// The entity leg files the deriving call's task id (TaskID(ctx),
	// internal/tools/subagent_197.go:262), never its correlation id: since
	// bd124b2a loop.go's callCorr mints one per call (<taskID>#<callID>, :603/:676).
```

为什么这一族事实是现读来的（⛔ 不是照抄票面）：`internal/tools/subagent_197.go:259-262` 现读
逐字 `// The task-level id of the caller: since ticket 242 the correlation id is` / `// minted per call (C18 routes approval replies by it), so it can no longer` /
`// name a task; the roster is keyed by task ids.` / `parentID := TaskID(ctx)`
⇒ `ParentTaskID` 今天**确实**读作父任务 id，但**原因不是 corr==taskID，而是 entities 腿改读了 `TaskID(ctx)`**
（`bd124b2a` 的「连带越格自报」那一跳）。旧注释把因果挂在了 corr==taskID 上，两处都过期。

## 件 2 — `internal/panel/subagent_roster_197_test.go`（改 :466-470，5 行 → 5 行；第 3 条注意事项整段）

改前逐字（`:466`–`:470`）：

```
//  3. The blockedOnApproval join is asserted against a card this file made. On the
//     production path the pairing depends on the loop dispatching with
//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a
//     different pairing would report every row as not blocked, and the only thing
//     that would catch it is the cmd/wisp case, which uses the real gate.
```

改后逐字（`:466`–`:470`）：

```
//  3. The blockedOnApproval join is asserted against a card this file made. On
//     the production path it pairs the row's task id, not the corr: since
//     bd124b2a callCorr mints one corr per call (loop.go:603, used :676),
//     shape <taskID>#<callID>, never == taskID. A host that paired on the
//     corr alone would miss every row; only cmd/wisp's real gate catches it.
```

现读依据：`internal/panel/subagent_roster_197.go:216-220` 对每枚 L1 window **两半都建表**
（`indexWaitingKey(waiting, w.CorrelationID)` 与 `indexWaitingKey(waiting, w.TaskID)`），
行侧查询在 `:235` `BlockedOnApproval: row.TaskID != "" && waiting[row.TaskID]`
⇒ 生产路径上点亮一行靠的是 window 携带的 **task id 那半**；只按 corr 那半配对会漏掉每一行
（corr 是 `<taskID>#<callID>`，名册行不按它建）。

## 件 3 — `internal/panel/subagent_blocked_220_test.go`（改 :155-156，2 行 → 2 行）

改前逐字（`:155`／`:156`）：

```
	// The production pairing (the loop dispatches with CorrelationID == TaskID,
	// internal/agent/loop.go:647) lights through the correlation id alone.
```

改后逐字（`:155`／`:156`）：

```
	// Since bd124b2a the loop's corr is per call (<taskID>#<callID>, loop.go:603,
	// used :676), never == taskID: the corr-only pairing below is this fixture's.
```

按票面 `现量#3` 的告诫本格读法：`:157` 那枚夹具（`{CorrelationID: "child-1", TaskID: ""}`）
**用例本身仍是绿的**，过期的只有那两句叙述 ⇒ ⛔ 本格不是"这里有一枚红"。

## 件 4 — `cmd/wisp/subagent_carrier_197_test.go`（改 :20-23，4 行 → 4 行）

改前逐字（`:20`–`:23`）：

```
// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,
// internal/agent/loop.go:647, which is also what makes the roster's parent link read
// as a task id). The spawn therefore waits for the ROOT's own roster row - proof the
// root loop is live - and only then dispatches through rt.bridge.
```

改后逐字（`:20`–`:23`）：

```
// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
// NOT the loop's: since bd124b2a callCorr mints one corr per call
// (<taskID>#<callID>, internal/agent/loop.go:603, used :676). The spawn waits
// for the ROOT's own roster row - proof it is live - then uses rt.bridge.
```

现读依据：本件 `:172-173` 现跑逐字
`out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{` / `			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),`
⇒ 把同一个值塞给两个字段的是**夹具自己**，"loop 自己用的 id"那一说是过期叙述。
本件行锚 `:173` 在改后仍在 `:173`（见下面锚位复量）。

## 尺 A — 每枚件 `git diff --numstat`（判据：`+N` 与 `-N` 相等）

命令逐字：

```
for p in internal/panel/subagent_roster_197.go internal/panel/subagent_roster_197_test.go \
         internal/panel/subagent_blocked_220_test.go cmd/wisp/subagent_carrier_197_test.go; do
  echo "--- $p ---"; git diff --numstat -- "$p"
done
```

输出原样：

```
--- internal/panel/subagent_roster_197.go ---
3	3	internal/panel/subagent_roster_197.go
--- internal/panel/subagent_roster_197_test.go ---
5	5	internal/panel/subagent_roster_197_test.go
--- internal/panel/subagent_blocked_220_test.go ---
2	2	internal/panel/subagent_blocked_220_test.go
--- cmd/wisp/subagent_carrier_197_test.go ---
4	4	cmd/wisp/subagent_carrier_197_test.go
```

⇒ 四枚件逐枚 `+N == -N`（3/3、5/5、2/2、4/4）。合起来一把（同四枚件一次过）原样：

```
4	4	cmd/wisp/subagent_carrier_197_test.go
2	2	internal/panel/subagent_blocked_220_test.go
3	3	internal/panel/subagent_roster_197.go
5	5	internal/panel/subagent_roster_197_test.go
```

## 尺 B — 非注释行必须＝0

命令逐字（diff 全量落盘＝本目录 `raw-diff.md`）：

```
git diff -U0 -- <这 4 枚件> > .scratch/wisp/probes/289/r1/raw-diff.md
total=$(grep -E '^[+-]' raw-diff.md | grep -vE '^(\+\+\+|---)' | wc -l)
cmt=$(grep -E '^[+-][[:space:]]*//' raw-diff.md | wc -l)
```

输出原样：

```
全部非 +++/--- 的 +/- 行 = 28
匹配 ^[+-][[:space:]]*// 的纯注释行 = 28
非注释行 = 0
```

## 锚位复量（改后工作树，`sed -n '<N>p'`，N 用票面原行号）

```
--- internal/panel/subagent_roster_197.go:57 ---
	// internal/tools/subagent_197.go:262), never its correlation id: since
--- internal/panel/subagent_roster_197_test.go:468 ---
//     bd124b2a callCorr mints one corr per call (loop.go:603, used :676),
--- internal/panel/subagent_blocked_220_test.go:155 ---
	// Since bd124b2a the loop's corr is per call (<taskID>#<callID>, loop.go:603,
--- cmd/wisp/subagent_carrier_197_test.go:20 ---
// bridge with the fixture's OWN ids (TaskID == CorrelationID at :173), which is
--- cmd/wisp/subagent_carrier_197_test.go:173 ---
			TaskID: rootID, CorrelationID: rootID, CallID: fmt.Sprintf("spawn-197-%02d", i),
```

⇒ 四枚原行号仍是那一块注释（同一行的位置没漂），件 4 里被引用的夹具行 `:173` 也没漂
⇒ 票 220／票 197 现量表按行号引 `subagent_roster_197.go`（`:190-210`／`:127`／`:99-104`）的那些锚**一格没被推动**。

## 禁区自证（本腿实际读到的第 5 枚件，⛔ 未动）

`grep -rn "CorrelationID == TaskID\|TaskID == CorrelationID" --include=*.go .` 现跑命中 5 枚**在册**件
＋1 枚 probes 旧快照：

- 在册 4 枚＝本格写面（上面逐枚）。
- 第 5 枚＝`internal/tools/ticket283_corr_identity_rulers_test.go:7`：
  `// 一律 TaskID == CorrelationID（bridge.go 的 orDefault 回落也让缺 corr 的`——
  **本腿一字未动**（票面写了「⛔ 不许碰第 5 枚文件」）。
  现读上下文（`:5-8`、`:14-16`）后具名结论：**那一句不是假话**——它说的是"既有用例的派发 fixture
  一律 TaskID == CorrelationID"，紧接着 `:14-15` 自己就写了真生产链
  `agent.Loop 的 callCorr(taskID, callID, i)（loop.go 每枚调用一个 corr，形如 taskID + "#" + callID）`。
  ⇒ 不需要补票；但记在这里，免得裁决者以为本腿漏了一枚。
- probes 旧快照＝`.scratch/wisp/probes/197/r3b/pre/subagent_roster_197.go:52,173`——
  那是别人票 197 的改前存档件，⛔ 不归本票射程，本腿未动。

`grep -rn "loop\.go:647" --include=*.go .` 改后现跑：在册件里**只剩 probes 那一枚旧快照**命中
⇒ 4 枚件的 `:647` 旧锚已全部换成 `:603`（定义）／`:676`（使用），第 5 枚件不带该锚。
