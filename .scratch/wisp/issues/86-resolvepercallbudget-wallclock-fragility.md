# 86 — `TestResolvePerCallBudget` 那条 1ms 墙钟预算：并跑必假失败，本轮连单跑也红

**Status:** open（排队）
**Type:** 测试质量（脆弱性），**不是**性能回归
**Blocks:** nothing · **Blocked by:** nothing
**Packages:** `internal/risk/`（那条用例所在的测试文件）。**禁改**：冻结面与 `docs/`。

## 实测（两位代理各自独立量到，数字都在票面上）

- 票 80：`internal/config` + `internal/risk` **两包并跑** ⇒ `--- FAIL: TestResolvePerCallBudget`
  （`1.640779ms > 1ms`）；**单跑同一测试 rc=0**（`0.692` / `0.637 ms/op`）；**单包 `-count=2` rc=0**。
- 票 83：同样并跑红，且**这一轮连单跑也红**（`1.356` / `1.194 ms/op`）。
  它的包级门禁本身是全绿：`go test -count=2 -v ./internal/config/` = `RUN 194 / PASS 106 / FAIL 0 / SKIP 0`。

⇒ 判定：**负载敏感的墙钟断言**，不是断言写反，也不是被测实现变慢。

## 编排者已经定死的边界（照做）

- **那个 `1ms` 预算一个字不许动**（不许"为了绿调大到 3ms"）。本仓明令：不许调阈值、不许重测到运气好的那次。
- 正确方向是**让测量本身稳**：把预算从"绝对墙钟"改成**同机自校准**
  （例如先在同一次运行里测一个基线倍数），或在 `testing.Short()` 下只跑功能断言、不跑墙钟断言
  —— 但**必须仍然有一条用例能抓到"每次解析的开销真的涨了"这件事**，不许把它变成空断言。

## AC

- [ ] **AC#1** 先量出这条用例的**真实分布**：`-count=10` 单跑 / 与 `internal/config` 并跑 / 与全仓并跑，
      三档各给 `ms/op` 的**全部样本**（不许只报平均、不许抹掉坏尾部），并给出机器核数。
- [ ] **AC#2** 按上面的边界改成自校准或 `testing.Short()` 门控，**并保留一条能抓到真回归的断言**
      （变异检验：人为把每次解析多塞一次真实工作 ⇒ 该用例必须红）。
- [ ] **AC#3** 三档并重跑：单跑 rc=0、并跑 rc=0、且**没有任何新增 SKIP**（逐条点名 `--- SKIP`）。
- [ ] **AC#4** 门禁：`gofmt -l` 空、`go vet ./internal/risk/` rc=0、`go test -count=2 ./internal/risk/` rc=0。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + `next=`；
`git commit -q -F - -- <显式路径>` + 带引号 heredoc；禁 `git add -A`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；
**不 push**；不在仓内建 worktree（A38④）；票面 append-only，要改的那行先读再替换；四种假绿逐跑点名。



## Progress log（append-only）

- 2026-09-21 15:5x（编排者，**排队原因，别当成"没人想到"**）：本票**故意现在不派**。
  它测的是**一条 1ms 墙钟预算的分布**，而编队里同时有 2 个写码代理 + 5 个验收/只读代理在跑 `go test` 与 docker
  ⇒ **在吵的机器上量出来的分位数没有意义**，量到"红"还会误判成回归。
  **排程判据（三件事同时成立才派）**：① 编队里 `internal/risk/` 无人写（票 94 已闭、票 82 的验收代理只读快照）；
  ② 没有 docker/`go test -count=` 的并发跑动（`git status --porcelain` 里没有 `internal/risk` 的改动即可初判）；
  ③ 我自己先看一眼当轮 `test-windows`/`test-core` 的 run id，确认这条用例**当时是红还是绿**，
     好让代理有"CI 侧的第三个样本"可比。
  ⚠ AC#1 要求"三档分布"（单跑 / 与 config 并跑 / 与全仓并跑），**多样本全报**是本票的底线：
  拿一次安静的单跑数字去勾框，等于把这条脆弱性重新藏起来。

（空）
