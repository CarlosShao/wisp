# 285-r2 · 突变自证（三发，全种在盘上，⛔ 零 `-overlay`）

种刀对象只有一个文件：`internal/agent/loop.go` 里的 `callCorr`（唯一调用点在 `executeCalls` 的派发那一跳）。
起手／每发种前／每发还原后的 hash 全部逐字相同：

```
8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go
```

三发还原后各跑一次 `git status --porcelain -- internal cmd` ＝ **空输出**（原文无行）。
⛔ 本腿自己的两枚测试件从未当过种刀对象（没有 MU-t 那一发；判别力靠 MU-1/MU-2/MU-3 三轴互不打扰来给，见 `10-ac4-ruler.md` §3）。

未突变时的对照读数（⛔ 不是"只报改前绿"，是三轴各自的红在下面）：

```
--- PASS: Test285CallCorrIsDistinctPerCall (0.00s)
--- PASS: Test285LoopDispatchesOneCorrPerCall (0.01s)
--- PASS: Test285RosterRowsCarryDistinctCorrPerCall (0.08s)
```

---

## MU-1 · 带 call id 那一支塌成一枚常量后缀（票面点名的第一形）

- 种前 hash：`8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go`
- 种前 `sed -n '603,608p' internal/agent/loop.go` 原文：

```
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return fmt.Sprintf("%s#call-%d", taskID, index)
	}
	return taskID + "#" + callID
}
```

- 种后同一行（`sed -n '603,608p'`）：第 607 行改成 `	return taskID + "#call-285-static"`（其余四行未动）
- 效果：同一任务的两枚调用共用一枚 corr，但**非空／以 task id 为前缀／不等于 task id 三条形状照旧成立**
  ——正是票面说的"逐枚形状尺全绿"那一形。

红句原文（`go test ./internal/agent/ -run 'Test285' -count=1 -v`）：

```
    ticket285_corr_distinct_rulers_test.go:79: two calls of one task share routing key "task-285-corr#call-285-static" (first call and second call): the two asks collapse onto one card
    ticket285_corr_distinct_rulers_test.go:85: routing keys = 1, want 2 distinct ones (one per call of the task)
    ticket285_corr_distinct_rulers_test.go:89: first call and second call carry the same corr "task-285-corr#call-285-static"
--- FAIL: Test285CallCorrIsDistinctPerCall (0.00s)
    ticket285_corr_distinct_rulers_test.go:124: both calls of one task dispatched with corr "811ea7f6-4bae-4518-9e40-f8a53d812ed6#call-285-static": the two asks share one routing key
--- FAIL: Test285LoopDispatchesOneCorrPerCall (0.01s)
FAIL	github.com/CarlosShao/wisp/internal/agent	0.042s
```

红句原文（`go test ./internal/tools/ -run 'Test285RosterRowsCarryDistinctCorrPerCall|Test283IdentityChainThroughTheRealLoop' -count=1 -v`）：

```
--- PASS: Test283IdentityChainThroughTheRealLoop (0.08s)
    ticket285_corr_rows_rulers_test.go:82: 同一任务的两枚调用共用一枚 corr "84007b6a-1089-4dff-b7ce-99885242e488#call-285-static"（task.spawn 与 task.cancel）：两枚 ask 塌成同一个路由键
    ticket285_corr_rows_rulers_test.go:90: corr "84007b6a-1089-4dff-b7ce-99885242e488#call-285-static" 同时是 task.spawn 与 task.cancel 的路由键：名册里两枚调用不可分
    ticket285_corr_rows_rulers_test.go:95: 路由键 = 1 枚，want 2（一枚调用一键）: map[84007b6a-1089-4dff-b7ce-99885242e488#call-285-static:task.cancel]
    ticket285_corr_rows_rulers_test.go:107: task.cancel 的行 corr = "84007b6a-1089-4dff-b7ce-99885242e488#call-285-static", want 后缀带着 "cancel"（各自可路由，不是同一本键）
--- FAIL: Test285RosterRowsCarryDistinctCorrPerCall (0.06s)
FAIL	github.com/CarlosShao/wisp/internal/tools	0.177s
```

**同一夹具、同一真链、同一发突变：票 283 的尺二 PASS，本腿的乙 FAIL**
⇒ 作差那一发是**新增的承重件**，不是把既有覆盖换个名字再写一遍。

- 还原：`sed -n '607p'` ＝ `	return taskID + "#" + callID`；
  还原后 hash `8eb37e9f59e563b58325d2ebc898b843`（与起手逐字等值），`git diff --stat -- internal/agent/loop.go` 空，
  `git status --porcelain -- internal cmd` 空。

## MU-2 · 回落支的 `index` 抹成常量（票面点名的第二发／反向那一支）

- 种前 hash：`8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go`
- 种后 `callCorr` 全文（`sed -n '603,608p'`）：

```
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return fmt.Sprintf("%s#call-0", taskID)
	}
	return taskID + "#" + callID
}
```

  （第 605 行：`%s#call-%d` + `index` ⇒ `%s#call-0`、参数被抹掉；其余四行未动）

红句原文（`go test ./internal/agent/ -run 'Test285' -count=1 -v`）：

```
    ticket285_corr_distinct_rulers_test.go:79: two calls of one task share routing key "task-285-corr#call-0" (first id-less call and second id-less call): the two asks collapse onto one card
    ticket285_corr_distinct_rulers_test.go:79: two calls of one task share routing key "task-285-corr#call-0" (second id-less call and third id-less call): the two asks collapse onto one card
    ticket285_corr_distinct_rulers_test.go:85: routing keys = 1, want 3 distinct ones (one per call of the task)
    ticket285_corr_distinct_rulers_test.go:89: first id-less call and second id-less call carry the same corr "task-285-corr#call-0"
--- FAIL: Test285CallCorrIsDistinctPerCall (0.00s)
--- PASS: Test285LoopDispatchesOneCorrPerCall (0.01s)
FAIL	github.com/CarlosShao/wisp/internal/agent	0.048s
```

对照读数（同一发突变下，⛔ 全仓其余尺都不响——这就是"回落支今天零尺"的直接证据）：

```
--- PASS: Test283TaskIDAccessorIsTheTaskIDNotTheCorr (0.00s)
--- PASS: Test283IdentityChainThroughTheRealLoop (0.11s)
--- PASS: Test285RosterRowsCarryDistinctCorrPerCall (0.10s)
ok  	github.com/CarlosShao/wisp/internal/tools	0.268s
```

外加现量：`grep -rn '#call-' internal/ --include=*_test.go` ＝ **空输出**（连本腿的件都不写这个字面串，
因为全仓 SSE 夹具没有一枚不带 call id）⇒ 这一支的形状尺**只能**是包内直接调用 `callCorr`，
这就是落点甲不可替代的理由（`10-ac4-ruler.md` §1）。

- 还原：`sed -n '605p'` ＝ `		return fmt.Sprintf("%s#call-%d", taskID, index)`；
  还原后 hash `8eb37e9f59e563b58325d2ebc898b843`（与起手逐字等值），`git status --porcelain -- internal cmd` 空。

## MU-3 · 保住互不相等、摘掉 task 前缀（第三轴，防"作差即万能"）

- 种前 hash：`8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go`
- 种后 `sed -n '607p'` ＝ `	return "foreign-" + callID`

红句原文（`internal/agent`）：

```
    ticket285_corr_distinct_rulers_test.go:76: first call corr = "foreign-call-285-a", want the task id "task-285-corr" as prefix
    ticket285_corr_distinct_rulers_test.go:76: second call corr = "foreign-call-285-b", want the task id "task-285-corr" as prefix
--- FAIL: Test285CallCorrIsDistinctPerCall (0.00s)
    ticket285_corr_distinct_rulers_test.go:133: call call_p2 corr = "foreign-call_p2", want the task id "f507c735-a6b2-4653-8b20-1d926f625add" as prefix
    ticket285_corr_distinct_rulers_test.go:137: call call_p2 corr = "foreign-call_p2" routes to "foreign-call_p2", want its own call id
    ticket285_corr_distinct_rulers_test.go:133: call call_p1 corr = "foreign-call_p1", want the task id "f507c735-a6b2-4653-8b20-1d926f625add" as prefix
    ticket285_corr_distinct_rulers_test.go:137: call call_p1 corr = "foreign-call_p1" routes to "foreign-call_p1", want its own call id
--- FAIL: Test285LoopDispatchesOneCorrPerCall (0.01s)
```

红句原文（`internal/tools`）：

```
    ticket283_corr_identity_rulers_test.go:223: tool task.spawn 的 correlation_id = "foreign-call-283-spawn", want 非空 per-call id（前缀 = task id "7de56d96-8265-44ad-b787-1531f0300215"，且不等于它）
    ticket283_corr_identity_rulers_test.go:223: tool task.cancel 的 correlation_id = "foreign-call-283-cancel", want 非空 per-call id（前缀 = task id "7de56d96-8265-44ad-b787-1531f0300215"，且不等于它）
--- FAIL: Test283IdentityChainThroughTheRealLoop (0.07s)
    ticket285_corr_rows_rulers_test.go:76: tool task.spawn 的行 corr = "foreign-call-283-spawn", want 以 task id "45a4841b-924d-4200-9e36-79b0a4494cd4" 为前缀
    ticket285_corr_rows_rulers_test.go:76: tool task.cancel 的行 corr = "foreign-call-283-cancel", want 以 task id "45a4841b-924d-4200-9e36-79b0a4494cd4" 为前缀
--- FAIL: Test285RosterRowsCarryDistinctCorrPerCall (0.06s)
```

三轴互不打扰的现量：MU-3 下**作差那一发（甲 `:88`／乙 `:81`）静默**（两枚 corr 仍然互不相等），
红只落在前缀与"后缀＝自己的 call id"那一发；MU-1 下**前缀那一发静默**。
⇒ 每根轴都是独立承重的，这把尺不是一句"不同就行"的恒真句。
另外 MU-3 也让票 283 的尺二变红 ⇒ 那一形今天已有尺（本腿不占功）。

- 还原：`sed -n '607p'` ＝ `	return taskID + "#" + callID`；
  还原后 hash `8eb37e9f59e563b58325d2ebc898b843`（与起手逐字等值），`git status --porcelain -- internal cmd` 空。

---

## 收工复核（三发全还原之后）

```
$ md5sum internal/agent/loop.go
8eb37e9f59e563b58325d2ebc898b843 *internal/agent/loop.go     # 与起手逐字等值
$ git status --porcelain -- internal cmd
                                                              # 空输出
$ git diff --name-status ce18b3d..HEAD -- '*.go'
A	internal/agent/ticket285_corr_distinct_rulers_test.go
A	internal/tools/ticket285_corr_rows_rulers_test.go
$ gofmt -l <两枚新件>        # 空
$ gofumpt -l <两枚新件>      # 空（D:\work\base\gopath\bin\gofumpt.exe，只读 -l）
$ ./tools/d22scan/d22scan.exe -root .
d22scan: clean - no D22 ban violations                        # rc=0
```

整包前后（`-count=1`，`-v` 数 `^--- (PASS|FAIL|SKIP)`，`rc` 取不带 `-v` 那一次）：

| 包 | 起手 | 收工 | rc |
|---|---|---|---|
| `internal/agent` | 82 PASS / 1 FAIL / 0 SKIP | 84 PASS / **1 FAIL** / 0 SKIP | 起手 1 ／ 收工 1 |
| `internal/tools` | 208 PASS / 0 FAIL / 0 SKIP | 209 PASS / 0 FAIL / 0 SKIP | 起手 0 ／ 收工 0 |

`internal/agent` 收工那枚 FAIL 与起手逐字同一枚（`TestGoldenSingleToolCall`，`loop_golden_test.go:70`），
本腿既没让它变绿（不属 AC#4 射程、⛔ 未放宽任何断言），也没新增红。
