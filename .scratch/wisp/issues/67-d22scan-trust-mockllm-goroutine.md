# 67 — 让 D22 静态门重新可信：`mockllm.go` 裸 goroutine + emoji 门看不见 Go 字符串

**Status:** in-progress
**Claimed by:** agent-ticket67
**Last update:** 2026-09-20
**Blocked by:** —（包与票 66 不相交：`internal/llm/adaptertest` + `tools/d22scan`；**AC#3 例外，须等票 66 落地**，见下）
**Parallel slots:** ≤1 sub-agent（**不要**碰 `cmd/wisp/`、`internal/proc/`、`internal/observe/` —— 票 66 在飞）
**Spec refs:** D22（七条禁令不可协商）、D38b、PLAN §16、票 12 AC#7
**登记项:** A22（编排者的假绿）、A23（emoji 门覆盖面 + 票 12 带进 3 处字形）

## 为什么要这张票：**CI 的 lint job 从约 05:59Z 起就是红的，而我曾把它读成 clean**

`tools/d22scan` 是**独立 Go module**（`tools/d22scan/go.mod`）。从仓根跑
`go run ./tools/d22scan -root .` 只会打印
`main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan`
——**扫描器从未执行**。**我自己就是这么误判的**，还据此在 registry A13 末尾否定了票 21 代理的一条**如实报告**。

正确调用（与 `.github/workflows/ci.yml` 的 "D22 seven-ban + emoji scan" 步逐字同形）：

```
cd tools/d22scan && go run . -root "${WORKSPACE}"
→ internal/llm/adaptertest/mockllm.go:68: [bare-goroutine] bare `go func(` is banned (D22/D38b):
  use observe.Registry.Spawn (named, owner, recover boundary)
→ d22scan: 1 finding(s); D22 bans are not negotiable
→ exit 1
```

引入 commit：**`78b1466`（票 11）**，`tools/d22scan/allowlist.txt` 里**没有**这一条。
CI 该步**无 `continue-on-error`** ⇒ **job 红**。

## 已完成判据（AC）

- [x] **AC#1 `mockllm.go:68` 的裸 goroutine 改走受管出口。** 现状：`internal/llm/adaptertest/mockllm.go:68`
  在 `cmd.Start()` 之后起一个 `go func()` 读 stdout 等 `MOCKLLM_ADDR=`，**无 owner、无 recover**（D22 禁令 #1）。
  改成 `observe.Registry.Spawn`（命名 + owner + recover 边界）。
  ⚠ **不许用 allowlist 豁免把它变绿**——那正是"为了让门闭眼而放宽门"。
  豁免只有我（编排者）能书面给，且必须写明理由；**本票的立场是这条不该豁免**。
  ⚠ 该 helper 服务于测试：**不得**为了让 spawn 好写而把 helper 改成并发不安全或丢掉 `t.Fatal` 语义；
  超时判定**禁止**用墙钟差（D22 禁令之一），沿用仓库里既有的 monotonic 做法。
- [ ] **AC#2 门自己必须是可证伪的。** 完成判据：`cd tools/d22scan && go test ./...` 通过
  （那步就是 seeded-violation 阳性对照），**且**你新加/改动的每条 ban 覆盖面都要有一条"故意种一个违规 → 扫描器必须报"的用例。
  我踩的坑不是"扫描器漏报"而是"我用的调用方式让它没跑"，所以这一格的交付物包括
  **把仓根误调用也变成不可能**：要么在 `scripts/` 加一个包装入口并在 CI 与文档里统一用它，
  要么让它对错误调用**显式报错退出**。
- [ ] **AC#3 emoji 门覆盖面（本票唯一要动 `cmd/wisp/` 的半格，等票 66 落地后再做）。**
  今天 `main.go:71` 的 `emojiRe` 只被 `walkEmoji` 用在 `design/` 与 `frontend/`（`main.go:137-142`），
  而 **`frontend/` 在本 HEAD 不存在** ⇒ 门对 `internal/`+`cmd/` 的 Go 字符串字面量**完全不可见**。
  真命中 3 处、由**票 12 自己的 `cd011b8`** 带入：`cmd/wisp/providers.go:201`（`U+2717`）、`:203`（`U+2713`）、
  `:209` 打到 stdout，另 `:41` 在 help 文本里。
  **两步**：①把那三处字形改成 ASCII 文案（推荐 `PASS`/`FAIL`——Windows 控制台字体不保证有这两个码位，
  这本来就是可读性缺陷，不只是门的问题）；②把 ban #8 的覆盖面扩到 `internal/`+`cmd/` 的
  **字符串字面量**（**注释不算用户可见**：另 5 处命中都在注释里，扩展时要把注释排除掉，
  不许为了少改代码就把注释也报成违规）。
  ⚠ 扩覆盖面**只许从严、不许放松**：如果扩完在别处又挖出命中，**逐个登记**而不是加豁免。
- [ ] **AC#4 门禁**：`gofmt -l` 触及包为空、`go vet ./...`（主模块）与 `go vet ./...`（`tools/d22scan` 模块）、
  `go test -count=2 ./internal/llm/adaptertest/`、`cd tools/d22scan && go test ./...`，
  最后 `cd tools/d22scan && go run . -root ../..` 必须 **0 命中且 exit 0**。逐条贴原始输出。

## 编排者已裁定（不要再来问）

1. **不放宽任何一条 D22 禁令**；`mockllm.go` 走受管出口，不豁免。
2. **票 12 的 AC#7 保持未勾**，直到本票 AC#3 落地；落点措辞由我在票 12 票面改，不归代理。
3. 本票**不动** `internal/observe/`、`internal/proc/`、`cmd/wisp/slo*.go`、`scripts/slo-check.ps1`、`docs/SLO.md`
   （票 66 的领地）；`cmd/wisp/providers.go` 只在 AC#3 的第二步动，且**在票 66 收尾之后**。

## 简报会附加的流程硬规矩（代理必读）

每个 commit 同步本票文件（状态字段 + 勾框 + `- [ISO-UTC] agent=... did=...` 一行）；
**前 15 次工具调用内必须有一个 commit**；`git add` **只准显式路径**，`git add -A`/`git add .` 禁用；
**永不 `git push`**；不复现不勾框；测不过就报 FAIL 并附数字，不许改阈值。

## Progress log（append-only）

- [2026-09-20T14:43:55Z] agent=agent-ticket67 did=认领本票（AC#1+AC#2 范围；AC#3 因票 66 在飞 `cmd/wisp/` 不动）。
  逐条复现票面描述，全部为真：①仓根 `go run ./tools/d22scan -root .` →
  `main module (github.com/CarlosShao/wisp) does not contain package .../tools/d22scan`，exit 1，扫描器从未执行；
  ②`cd tools/d22scan && go run . -root ../../` → `internal/llm/adaptertest/mockllm.go:68: [bare-goroutine] ...`
  + `d22scan: 1 finding(s)`，exit 1；③`cd tools/d22scan && go test ./...` **当前 FAIL**
  （`TestScannerSelfScanOfRealRepoIsGreen` 报同一条命中）⇒ AC#2 第一格在 AC#1 落地前不可能绿。
  已核对 `tools/d22scan/allowlist.txt` **无** mockllm 条目（4 行，均为 pathresolver/mirror-hash）。
  待办：AC#1 走 `observe.Registry.Spawn` + `observe.NewTimeout`（monotonic）；AC#2 种违规阳性对照 + 把误调用变响。
  未做：AC#3（编排者裁定须等票 66）、AC#4 门禁清单。
- [2026-09-20T14:52:00Z] agent=agent-ticket67 did=**AC#1 落地并勾选**。
  `internal/llm/adaptertest/mockllm.go`：读 stdout 等 `MOCKLLM_ADDR=` 的工作改走
  `mockllmRegistry.Spawn("mockllm-stdout-reader", "test", spawnRoot, fn)`（包内私有
  `observe.NewRegistry()`，**不**污染 `observe.Default` 的 resident 花名册计数）；
  超时判定由 `time.After(20s)` 改为 `observe.NewTimeout(mockllmReadyBudget)` 的 monotonic 预算（禁令 #4/D42#9）；
  `t.Cleanup` 提前到等待之前，失败路径同样 Kill+Wait 子进程，并 `spawnRoot.Wait(2s)` 断言 reader 已 drain（漏了就直接 `t.Errorf`）。
  **未加任何 allowlist 条目**（`allowlist.txt` 仍 4 行）；`h.Done()` 分支把 recover 到的 panic 转回 `t.Fatal`，
  不让它退化成 20s 挂死。
  测量：`cd tools/d22scan && go run . -root ../../` → `d22scan: clean`，exit 0（原为 1 命中/exit 1）；
  `gofmt -l internal/llm/adaptertest/` 空；`go vet ./internal/llm/adaptertest/` 过；
  `go test -count=2 ./internal/llm/adaptertest/` ok；`cd tools/d22scan && go test ./...` ok（`TestScannerSelfScanOfRealRepoIsGreen` 由红转绿）。
  变异自证（种完即撤，`grep MUTATION|bogus` 零残留）：①子进程启动即退 → 1.32s **快速失败**
  `mockllm exited before printing its address`；②地址行永不匹配且子进程存活 → 21.30s 失败
  `mockllm did not become ready in 20s`（证明 monotonic 预算真在兜底，不是永久挂死）。
  已知副作用，如实登记：`Spawn` 对花名册外名字无条件 `slog.Warn`，故 mockllm 用例每次多一行
  `WARN goroutine outside the D38 roster (leak symptom) goroutine=mockllm-stdout-reader owner=test`。
  消除它要往 `internal/observe/goroutine.go` 的 `TemporaryNames` 加名——那是票 66 领地，本票不动；
  也不肯借用 `disposal-worker` 之类既有名字骗过漏检器。**留给编排者裁定。**
