# 139-a1 普查证据件 — 「上下文压缩不留痕」这一格在现树上的真实密度

派单代号 `139-a1`（只读普查腿）。票面：
`.scratch/wisp/issues/139-context-compression-leaves-no-trace-zero-log-calls-in-compress-go-zero-readers-of-res-compression.md`。
本篇**只做一件事**：把这一格在**当前工作树**上的读数取到编排者能当场决定派哪一枚落地腿的密度。

---

## 0. 锚点与口径

- 取锚时刻 `2026-10-03 10:01:25 +0800`；`git log -1 --format=%H` = `ccd9543b3d86f3f627d331b3c399cb3d4ee99eec`；分支 `dev`。
- 工作树**非干净**（`git status --porcelain` 有他程在飞的改动，含 `cmd/wisp/config_receipt_255_test.go` 与
  `.scratch/wisp/probes/**` 若干），故本篇所有**产码读数一律取自工作树**、所有**归属读数一律用 `git log` 指到具体 commit**，
  两者不混用。工作树里 `internal/agent/{compress.go,loop.go,budgets.go}` 与 `cmd/wisp/run.go` 当前 **0 枚改动**
  （`git status --porcelain` 全文里无 `internal/` 也无 `cmd/` 的 M 行，除上面那枚 255 的测试文件），所以本格的产码读数可直接当锚点树读。
- 尺：只读工具（`Read` / `Grep` / `find` / `grep` / `sed -n`）。⛔ 本篇**未跑任何一枚 `go` 子命令**
  （`build` / `vet` / `test` / `run` / `list` 全零），原因＝同机三枚写腿在飞，编译会互洗读数。
  因此本篇**不报任何 rc、不报任何用例通过与否**；需要 rc 的判据一律落到 §6。
- **grep 分母口径**（本仓踩过：变异副本长得跟产码一样）：
  `cmd/**` + `internal/**` + `docs/**`，**排除 `.scratch/**`**。这条不是洁癖——
  `./.scratch/wisp/probes/183/a1/mut/loop.go` 是一枚**完整的 loop.go 副本**，
  它在 `res.Compression = rep`（该行副本 `:404`）和 `comp *Compressor`（`:182`）上都命中，
  不排除就会把「生产侧读者」数出虚高的一枚。本篇每条全称否定都按这个分母跑，并带正控。
- 路径写法：报到**包／目录级**＋行号（例 `internal/agent/compress.go:174`），不停在文件名。
- 禁令执行：`frontend/**` 与 `design/**` 未读、未引用；`internal/observe/thresholds.go`、golden、
  `tools/d22scan/allowlist.txt`、`.github/workflows/ci.yml`、`docs/PLAN.md`、`docs/specs/**`、票面 **一字未动**；
  产码／测试 **0 字节改动**（本程唯一写面＝本文件与同目录的提交说明 txt）。

---

## 1. 压缩到底在哪跑

### 1.1 生产侧唯一的入口链（逐枚行号，工作树锚 `ccd9543`）

| 跳 | 位置 | 是什么 |
|---|---|---|
| ① 构造（每枚 Loop 一次） | `internal/agent/loop.go:239` `comp: NewCompressor(b, opt.Summarizer, WithLogger(opt.Logger))` | 压缩机**不是按任务造的**，是 `agent.New` 里按 Loop 造的；`opt.Logger` 在 `:207-209` 已被兜成 `slog.Default()` |
| ② 字段落位 | `internal/agent/loop.go:182` `comp *Compressor` | Loop 的共享字段，跨任务复用 |
| ③ 判该不该压 | `internal/agent/loop.go:393` `if l.comp.Need(hist)` → 实现 `internal/agent/compress.go:119`（阈值比较在 `:120`） | 阈值＝ `internal/agent/budgets.go:107` 的 `HistoryCompressTokens`，由 `:36` `refHistoryTokens = 12000` 按 `:91` `Scale` 等比缩放（`scaleInt` 在 `:114`，只收不放） |
| ④ 真压 | `internal/agent/loop.go:399` `nh, rep, err := l.comp.Compress(withTraceTask(ctx, taskID), hist)` → 实现 `internal/agent/compress.go:174` | `withTraceTask` 定义在 `internal/agent/compress.go:152`，读侧 `:160` |
| ⑤ 失败边 | `internal/agent/loop.go:401` `l.log().Warn("agent: history compression failed", "err", err)` | 唯一一条**在调用方**的痕 |
| ⑥ 成功边 | `internal/agent/loop.go:402` `else if rep.Ran` → `:403` `l.replaceHistory(nh)` → `:404` `res.Compression = rep` | `replaceHistory` 见 `internal/agent/loop.go`（历史换手）；痕本身在压缩机里（下条） |
| ⑦ 痕的产出点 | `internal/agent/compress.go:240` `c.log().Info("agent: history compressed", attrs...)`，由 `:216` `if rep.Ran` 把门 | 属性表在 `:222-231`，task 键在 `:237-239` |

### 1.2 它在哪个循环里

`internal/agent/loop.go:376` `for {` —— **①~⑦ 整条链在 run 的轮次循环体内**（`hist := l.History()` 在 `:392`，
压缩紧跟其后，`buildRequest` 在 `:408`）。⇒ **一枚任务的一次 Run 里这段代码可以被走 N 次**（每轮模型调用前各一次），
不是一个任务一次的形状。这条直接决定 §2.3 那格「算了但没存」的读数。

### 1.3 「活在 run 局部」那一格（本仓旧账，不是新发现，但必须量全）

`res` 是 `run` 的局部变量：`internal/agent/loop.go:369` `res := Result{TaskID: taskID, Status: StatusCompleted}`。
它出 run 只有三条口，全部经 `Result` 这一个结构：

1. `RunAsync` 把它塞进 `t.result`（`internal/agent/loop.go:326`），由 `Wait()` 返回（`:308-314`）；
2. `Run` 直接返回（`:332-334`）；
3. 各收尾分支 `l.finish(...)` / `brakeStuck` / `failWith` 收的也是这枚 `&res`。

生产侧**真的拿到 `Result` 的腿**只有两枚：

- `cmd/wisp/run.go:1099` `bg := loop.RunAsync(ctx, task)` → `:1106` `res := bg.Wait()`，
  消费面是 `cmd/wisp/run.go:1184` `notifyBody(res agent.Result)` 与 `:1195` `costLine(res agent.Result)`；
- `internal/tools/subagent_197.go:344` `bg := child.RunAsync(...)` → `:383` `res := bg.Wait()`，
  消费面是 `:412` `finalize(...)` 与 `:437` `answer(...)`。

⇒ 「泵／快照读不到 run 局部」这一格在 `Compression` 上是**同一形状的第二次命中**：
`res` 到得了宿主、到得了子代理工具，**但它身上那枚 `Compression` 字段在两枚消费面里都没被看一眼**（逐跳见 §4）。

### 1.4 装配面（谁把这格接到盘上）

- `cmd/wisp/run.go:995` `opts := agent.Options{...}` —— 具名读一遍：`Provider` / `Tools` / `Sink` / `Journal: nil` /
  `AdmitTask` / `Registry` / 内嵌 `Config` 全在，**没有 `Logger` 这一项**。⇒ 生产 Loop 的压缩机走的是
  `internal/agent/loop.go:207-209` 的兜底 = `slog.Default()`。
- 常驻腿的面板入口**不装 sink**：`cmd/wisp/panel_inbound.go:49` 明写「it installs no persistent log sink (installLogSink)」，
  且 `grep -rn 'agent\.New(' cmd internal`（排除测试）在生产侧只有 **2 枚命中**：
  `cmd/wisp/run.go:1023`（`wisp run` 腿）与 `internal/tools/subagent_197.go:319`（子代理腿）。
  ⇒ **今天这台机器上会走到压缩的只有 `wisp run` 那条 CLI 腿和它的子代理**，常驻面板腿根本不构造 Loop（尺见 §3.4）。
- 落盘那条管子（痕**能**到盘的证据）：`cmd/wisp/run.go:222` `installLogSink(s.dataDir)` 在
  `:249` `assembleRuntime(s)` **之前**；`assembleRuntime` 内 `:995` 才造 `agent.Options`、`:1023` 才 `agent.New`。
  `installLogSink` 实现在 `cmd/wisp/logsink.go:144`，`:153` `slog.SetDefault(slog.New(teeHandler{...}))`，
  `teeHandler.Handle` 在 `:208-214`（primary＝`p.Handler()`＝带 redact 的 JSONL 文件，mirror＝stderr，`min_level` 写死在 `:82`）。
  ⚠ 一个**时序口径**必须记：`agent.New` 在 `:207-209` 把 `opt.Logger` 兜成的是**那一刻**的 `slog.Default()`，
  压缩机字段 `lg` 于是**持有一个 Logger 指针**（`internal/agent/compress.go:77` / `:87-89` / `:102-107`），
  不是每次调用现取默认值 ⇒ 若哪天真出现「先 `agent.New` 后 `installLogSink`」的装配序，痕会进旧默认（stderr），
  今天的顺序不会（`run.go:222` 早于 `:1023`）。
