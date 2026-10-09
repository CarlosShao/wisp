# 票 289 / 腿 289-r1 — 四句注释改前改后逐字 ＋ 两把行中性尺

> 本格只改注释行。下面每枚件：改前逐字 → 改后逐字 → 两把尺的原样命令与输出。
> 待填：改后逐字与两把尺输出（改完追加，⛔ 不改已写下的改前逐字）。

## 件 1 — `internal/panel/subagent_roster_197.go:57`

改前（`git show HEAD:… | sed -n '57p'` 现跑，见 `00-anchor-and-baseline.md` §2）：

```
	// loop dispatches with CorrelationID == TaskID (internal/agent/loop.go:647),
```

改后：待填。

## 件 2 — `internal/panel/subagent_roster_197_test.go:468`

改前：

```
//     CorrelationID == TaskID (internal/agent/loop.go:647); a host that used a
```

改后：待填。

## 件 3 — `internal/panel/subagent_blocked_220_test.go:155`

改前：

```
	// The production pairing (the loop dispatches with CorrelationID == TaskID,
```

改后：待填。

## 件 4 — `cmd/wisp/subagent_carrier_197_test.go:20`

改前：

```
// bridge with the ids the agent loop itself uses (TaskID == CorrelationID,
```

改后：待填。

## 尺 A — 每枚件 `git diff --numstat`（判据：`+N` 与 `-N` 相等）

命令：待填。输出：待填。

## 尺 B — 非注释行必须＝0

命令：`git diff -U0 -- <这 4 枚件>`，取所有非 `+++`/`---` 的行，比对 `^[+-][[:space:]]*//`。
输出：待填。
