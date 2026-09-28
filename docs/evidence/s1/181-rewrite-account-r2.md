# 181-r2 — 三处新消费者没读"改写账户"：现量结案表

- 程：`181-r2`（写码腿·小面，缺陷修复）｜派单：`.scratch/wisp/dispatches/2026-09-28-165x-impl-181-r2-rewrite-account-consumers.md`
- 起手锚：`git log -1 --format='%h'` = **`e9c216b8`**（本程自取，非派单里的 `4e66817`，与派单正文 `e9c216b8` 本轮未漂）
- 起手时刻：`date "+%Y-%m-%d %H:%M%z"` = **`2026-09-28 15:57+0800`**
- ⛔ 本表不勾任何 AC 框；`internal/risk/**` 一字未动（裁判）。

## ① 起手脏件名册（写面闸门基线）

尺：`git status --porcelain -- internal/ cmd/`

```
(空)
```

⇒ **起手名册 = 空集**，终态必须仍等于空集（派单 §1 的闸门是"等于起手名册"，不是"必须为空"；本程起手这两枚目录下恰好没有别家脏件）。
仓外别家脏件（`.gitignore`／`probes/152/**`／`probes/161/r6/logs/flip-*`／`design/**` 那批未提交删除）本程不动、不提交、不评论、不还原。

## ② AC#α 改写账户到底是什么（现读，非按名字猜）

尺：`grep -rn "Actable\|Rewritten" --include=*.go internal/risk/ | grep -v _test.go`（现跑）

### ②.1 账户由 `risk.Result` 的哪几枚字段表达

全部在 `internal/risk/pathresolver.go`：

| 字段／方法 | 位置 | 它表达什么（原文口径摘要） |
|---|---|---|
| `Canonical string` | `pathresolver.go:58` | 流水线的最终形；**本身不带任何"我是不是被换过"的信息** |
| `Spelling string` | `:70` | 调用方交进来的原拼法，"审批提示里人看到的那一行"，用来把 `Canonical` 上的判定回溯到被问的那个形 |
| `Rewritten bool` | `:76` | **账户的主字段**：展开步（SPEC-06 §4 step 1）有没有替换过东西（`%VAR%`／`$VAR`／前导 `~`），为真即 `Canonical` 可能命名了另一棵树 |
| `Rewrites []string` | `:79` | 是哪几种构造替换的（`"env"`／`"home"`），`Rewritten` 为假时必空 |
| `func (r Result) Actable() (string, error)` | `:93-:99` | **可"行事"的形**：`Rewritten` 为真时返回 `""` + `ErrRewrittenPath`（`:49`），否则才返回 `Canonical` |

`Actable()` 的原文（`pathresolver.go:93-99`）逐字：

```go
func (r Result) Actable() (string, error) {
	if r.Rewritten {
		return "", fmt.Errorf("%w: %s expands to %s (%s); act on the expanded tree only if you asked for it by name",
			ErrRewrittenPath, r.Spelling, r.Canonical, strings.Join(r.Rewrites, "+"))
	}
	return r.Canonical, nil
}
```

口径合同就写在同一文件 `:90-92`：

> Every security consumer of Result outside this file reads the account through this method or through Result.Rewritten; the static criterion in pathresolver_rewrite_account_test.go fails the build-level gate if one stops.

⇒ **"读了账户"在这套规矩里只有两种写法**：调 `res.Actable()`，或显式读 `res.Rewritten`。`Canonical` 单独取用＝票 102 的洞。

### ②.2 两枚已合规最近邻（逐字抄其形，均非本程写面）

**其一 `internal/risk/winsec_c26.go:39-45`**（"既动树又向调用方报成功"那一腿的形）：

```go
func (c26Pipeline) Resolve(input string) (string, error) {
	res, err := Resolve(input, nil)
	if err != nil {
		return "", err
	}
	return res.Actable()
}
```

它的 `ResolveAccounted`（`:53-63`）是同一枚 `Result` 的第二形——把账户**交出去**而不是在这里消费掉：`return "", res.Rewritten, err`（`:60`）／`return path, res.Rewritten, nil`（`:62`），注释 `:50-52` 明说"不是第二次解析、不是第二本账"。

**其二 `internal/risk/syncdirs.go:202-216`＋`:226`**（写判定那一腿，含祖先重解析）：

```go
	canon, err := res.Actable()
	if err != nil {
		return "", err
	}
	if canon == "" {
		return "", errTargetUnverified
	}
```

祖先那跳同样过账户：`anceCanon, err := ares.Actable()`（`:226`），失败注释就是票 102 的口径 `:228`"an ancestor C26 moved is not evidence about this one"。另有 `:143` `case err == nil && res.Rewritten:` 走"降级但仍记录"的形（同步根不再算 canonical 级证据，`syncdirs.go:153` 打日志说明）。

**第三枚参考（在裁判的 leg 名单里）**：`internal/tools/paths.go:82-93` 用 `if res.Rewritten { … }` 显式消费，并把这类根另列成 `RewrittenRoots()`（`:190-194`）供上层报告——即"显示真值＋明说被改写过"这一档**在本仓已有先例形状**。

### ②.3 本程写面今天为什么不合规（现量）

- `cmd/wisp/panel_pump.go:104` `return panel.ReadGit(rt.workspaceView().Canonical)`：**真代码**取 `Canonical`，同文件 `grep -c 'Actable(\|\.Rewritten'` = 0（⇒ 被尺判为未读账户）。
- `cmd/wisp/panel_pump.go:93`、`internal/panel/git.go:132`：这两枚是**注释行**里的 `WorkspaceView.Canonical` 字样（`:92-93`／`:131-133` 的说明句），尺的 `readsCanonicalField`（`pathresolver_rewrite_account_test.go:258`）**不剥注释**，故一并计入。⇒ 这条差别是本程 §⑤（AC#δ）与 §④（AC#γ）的现量对象，**不是把尺改窄的理由**。

## ③ AC#β 三处逐枚（改前／改后／三档后果）

待填（下一节起逐枚落地，三档后果先给、由编排者裁）。
