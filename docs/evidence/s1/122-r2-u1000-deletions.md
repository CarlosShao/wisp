# 票 122 — staticcheck 清账腿 122-r2 裁决表（AC#2 落地：U1000 判删候选 3 条真删＋AC#1 遗留证据件补 commit）

- 腿名：122-r2 · 日期：2026-11-11 · 写面：本腿（实现者）；验收方另派
- 开工 HEAD：`9c8e84098`（A1713）；分支 `dev`；上一腿＝`06a67b40f`（122-r1，零 .go 改动的底稿腿）
- **仪器版本（AC#7 格式要求，逐字）**：`staticcheck.exe 2026.2.2 (0.8.2)`、`gofumpt v0.12.0 (go1.27.1)`
  （本机 GOPATH bin：`D:\work\base\gopath\bin\`；与 122-r1 同版同机，读数可比）
- 跑法（与 122-r1 同形）：staticcheck 与纯净 build/vet 读数跑在 `git archive HEAD` 导出的
  **仓外纯净快照**（`%TEMP%\wisp122r2\snap`，tracked `.go` 计数 1723 == `git ls-files '*.go'` 1723），
  删后读数＝快照 + 本腿两枚改动件同步进去再跑；`go build`/`go vet`/测试/gofumpt/d22scan 也在**主工作树**跑
  （`build\mut_rvw.go` 是 gitignore 残渣、非 HEAD 内容，只影响主工作树 build/vet，见 §3）。
- **派单路径笔误登记**：派单写 `internal/agent/tools/**`，实际树与 122-r1 证据件均为 **`internal/tools/`**
  （`internal/agent` 存在但无 `tools/` 子目录）。目标文件按 122-r1 判删表的 file:line 无歧义对应，本腿写面按实际路径 `internal/tools/**` 执行。

---

## 1. 删了哪几条（逐条：锚点 + 删前零引用证据 + 删除形状）

| # | 锚点 | 删除内容 | 删前零引用证据（本腿现核，全仓 grep） |
|---|---|---|---|
| 1 | `internal/tools/shell_exec.go:68` | `waitOutcome.out`（`strings.Builder` 字段） | `waitOutcome` 全仓 5 处命中全在本文件（:61 注释、:65 定义、:127 chan、:130 构造、:170 `exitCodeOf` 形参）；构造处 `:130` 字面量只填 `state`/`err` 两键，`out`/`errs` 零写零读 |
| 2 | `internal/tools/shell_exec.go:69` | `waitOutcome.errs`（同上） | 同上；实际输出走 `:116` 的 `outBuf/errBuf`（`guardedBytes`），经 `:136` `shellResultSections(outBuf.String(), errBuf.String())` 进入结果文本，与 `waitOutcome` 无交集 |
| 3 | `internal/tools/shell_session_test.go:31` | `fakeSessionShell.closeErr`（`error` 字段） | `closeErr` 全仓 10 处命中中，属 `fakeSessionShell` 的只有定义行 `:31`；其余＝`bridge.go:1058/1064/1067` 局部变量、`fs_edit_ac34_test.go:558-560` 局部变量、`internal/memory/open.go:112/396/398` 另一类型（`store`）的字段——全是同名异物 |

三名的唯一命中均为定义行本身，无任何读写引用；与 122-r1 判删表（`122-r1-staticcheck-backlog.md:114-116`）逐条对应。

**删除形状（`git diff` 全文，+5/-8）**：
- `shell_exec.go`：删 `:68-69` 两行，`waitOutcome` 剩 `state`/`err` 两字段（该 struct 的注释文档本就只承诺这两个通道）。
- `shell_session_test.go`：删 `:31` 一行，其余 +5/-5 为 struct 字段对齐收窄（`closeErr` 删除后最长字段名变短，gofumpt 对齐要求）。

## 2. 删前/删后 staticcheck 对照（两 GOOS，快照同版）

原始输出：`.scratch/wisp/probes/122/r2/before-windows.txt`、`before-linux.txt`、`after-windows.txt`、`after-linux.txt`（逐行可复跑）。

| 读数 | GOOS=windows | GOOS=linux（root module） |
|---|---|---|
| 删前（HEAD `9c8e84098` 快照，本腿现跑） | **103 条 finding** | **58 行 = 57 条 finding + 1 条编译噪音**（sherpa-onnx-go-linux 构建约束排除，r1 已登记的环境形状） |
| 122-r1 基线（`e31937a90`） | 103 | 57 + 同一噪音行 |
| **删后（快照+本腿改动）** | **100** | **55** |
| before→after 差集 | 恰好 3 条，全为目标行 | 恰好 3 条，全为目标行 |

- HEAD `e31937a90`→`9c8e84098` 之间只进了 `cmd/wisp` 两枚测试件改动（`302` 波），本腿开工时把删前发现集合与 r1 原始输出逐条比对（归一斜杠后）：**windows/linux 差集均为空**——r1 基线未腐烂。
- 两 GOOS U1000 去重并集：删前 147 条（真 finding）→ 删后 144 条，差集恰为目标 3 条（`union-gone.txt`）。
- 预期 CI 读数（推送后）：**GOOS=windows 100 条 / GOOS=linux 54 条 + 1 条编译噪音行**（CI 侧另有 d22scan/mockllm 两 module，期望各 0 条，r1 已实测 0）。

## 3. 门禁读数（逐条贴 rc）

| 门禁 | 命令 | rc | 读数 |
|---|---|---|---|
| go build | `go build ./...`（主工作树） | **1** | 红因全部＝`build\mut_rvw.go`（gitignore 残渣，`.gitignore:14`；`git ls-files build/` 为空）。**非 HEAD 内容、预先存在**：122-r1 证据件 :8-9 已登记同一残渣。**纯净快照（HEAD）build rc=0**；快照+本腿改动 build rc=0（`gate-build-snap-after.txt`） |
| go vet | `go vet ./...`（主工作树） | **1** | 同上唯一红因 `build\mut_rvw.go`。**快照（HEAD）vet rc=0**；快照+本腿改动 vet rc=0（`gate-vet-snap-after.txt`） |
| go test（受影响包） | `go test -count=1 ./internal/tools/` | **0** | `ok github.com/CarlosShao/wisp/internal/tools 30.692s` |
| go test -v 细账 | `go test -count=1 -v ./internal/tools/` | **0** | top-RUN=572，PASS=422 顶层 + 149 子测试 = **571**，FAIL=0，SKIP=1（572=571+1 对账平） |
| d22scan | `D:\work\soft\Git\bin\sh.exe scripts/d22scan.sh` | **0** | 正例控制 PASS=45/FAIL=0；scan clean。scope 对 r1 不降：bans#1-5 internal/=402（持平）、cmd/=74（持平）、ban#6 frontend/=123（122↑，他腿新代码）、ban#7 internal/tools/=51（持平）、ban#10 internal/risk/=21（持平）、ban#8 internal/=1003（持平）、cmd/=273（122↑）、design/=39、frontend/=123（自然增长） |
| gofumpt | `gofumpt.exe -l internal\tools` | **0** | **0 行**（v0.12.0）；`gofmt -l internal\tools` 同样 0 行 |

## 4. 越界检查

- 本腿 tracked 写面全集＝2 枚：`internal/tools/shell_exec.go`、`internal/tools/shell_session_test.go`（`git status --porcelain -- internal/ docs/` 仅此两枚 `M`；docs 下的 `??` 全是他腿遗留未跟踪件，零接触）。
- 未动：`cmd/wisp/**`（314-ac0 占用）、`internal/panel/**`、`internal/risk/**`（冻结面）、thresholds/golden/PLAN.md/docs/specs/frontend/design——零字节。checker 开关、`//lint:ignore`、lint 步位置：一字未动（本腿 diff 只有 3 个字段删除＋对齐，无任何抑制形状）。
- 禁令遵守：零 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；commit 两笔均显式 pathspec；只 commit 不 push。

## 5. 交付物与 next=

- 本腿 commit（拟，两笔显式 pathspec）：
  1. **r1 遗留证据件补交**（派单令）：`.scratch/wisp/probes/122/r1/ac1-windows.txt`、`u1000-windows.csv`、`u1000-linux.csv`、`u1000-union.txt`（4 枚；r1 的 commit 只含 ac1-*-head 两枚＋ac5 log＋本表 4 枚，此 4 枚在工作树未跟踪）。
  2. **AC#2 落地**：`internal/tools/shell_exec.go`、`internal/tools/shell_session_test.go`、本表、`.scratch/wisp/probes/122/r2/`（18 枚探针/门禁原始输出＋2 枚本腿 commit message 存档，共 20 枚）。
- 原始门禁输出全部落 `.scratch/wisp/probes/122/r2/`（gate-*.txt、before/after-*.txt、diff-*、union-gone.txt、version.txt）。
- **next=编排者推送后取 CI `lint` job staticcheck 步读数**（期望 windows 100 / linux 54+噪音行，须同报 staticcheck 版本）关 AC#6；linux 分母归因照 r1 §AC#1 第 4 条不变。
- 后续批候选仍以 122-r1 判删表为底（本腿清掉 3 条后剩「留 119／冻 2／越界-cmd 8」），归属裁定权在编排者。

## 6. 回执

`DONE: 122-r2 删3条 build:主树rc1(残渣build/mut_rvw.go,非HEAD,快照rc0)/vet:同前/test:0(rc0,572 RUN,571 PASS,0 FAIL,1 SKIP)/d22scan:0(clean) 越界=0`
