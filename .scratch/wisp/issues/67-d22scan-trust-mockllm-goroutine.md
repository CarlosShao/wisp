# 67 — 让 D22 静态门重新可信：`mockllm.go` 裸 goroutine + emoji 门看不见 Go 字符串

**Status:** blocked-on-ticket:66（AC#1/AC#2 **已交付并经编排者独立验证**；AC#3 要改 `cmd/wisp/providers.go`，与票 66 同包）
**Claimed by:** agent-ticket67（报告已交，**勿重开 AC#1/AC#2**；AC#3 等票 66 收尾后由**新代理接续**，从 Progress log 的 `next=` 起）
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
- [x] **AC#2 门自己必须是可证伪的。** 完成判据：`cd tools/d22scan && go test ./...` 通过
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
- [x] **AC#4 门禁（仅对已交付的 AC#1/AC#2 范围；AC#3 落地后必须重跑）**
      —— 由编排者独立复跑确认：`cd tools/d22scan && go run . -root ../..` →
      `clean`、**`examined 194 production Go files`**、**exit 0**；`go test ./...`（d22scan 自身 seeded-violation）通过；
      `go test -count=2 ./internal/llm/adaptertest/` → **ok 2.900s**；`gofmt -l` 空、`go vet` 干净；
      并核对 `allowlist.txt` **仍是 4 行、没为凑绿加豁免**（裁定 1 的要求）。
      原判据（保留可读）：`gofmt -l` 触及包为空、`go vet ./...`（主模块）与 `go vet ./...`（`tools/d22scan` 模块）、
      `go test -count=2 ./internal/llm/adaptertest/`、`cd tools/d22scan && go test ./...`，
      最后 `cd tools/d22scan && go run . -root ../..` 必须 **0 命中且 exit 0** —— **六条全部逐条复现过**。

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
- [2026-09-20T15:05:00Z] agent=agent-ticket67 did=**AC#2 落地并勾选**，并**挖出第二条比票面更危险的假绿**。
  【新发现，须登记 registry A22】`tools/d22scan` 的 `-root` 默认 `.` ⇒
  **`cd tools/d22scan && go run .`**（不带 `-root`，正是 `docs/evidence/s1/10`、`docs/evidence/s3/19`
  里逐字记录为 "PASS clean" 的那条命令）**扫的是 `tools/d22scan` 自己这棵树**：
  无 `internal/`、无 `cmd/`、无 `design/`、allowlist 缺失也当空表 ⇒ **0 文件被检查、0 命中，却打印
  `d22scan: clean` + exit 0**。票面只说了"从仓根跑会 no-op"（那条至少吐 Go 的模块错误、非零退出）；
  **这条是打印成功的空转**，比仓根误调用更可信也更不可见。历史证据里那些 clean 行是这么来的。
  【修法】`tools/d22scan/main.go` 加 `checkRoot()`：`-root` 必须同时有 `go.mod`、
  `tools/d22scan/allowlist.txt`、`internal/`、`cmd/`，且 `internal/`+`cmd/` 下非测试 `.go` 文件
  ≥ `minProductionGoFiles`（10，同 `internal/observe/nobarego_test.go` 的守卫形状），否则
  **exit 2** 并打出该跑哪条命令；绿/红之前先打印 `examined N production Go files`，
  让"没跑"与"没命中"在 stdout 上可分辨。allowlist 缺失从"静默当空表"变成致命。
  【包装入口】`scripts/d22scan.sh`：从自身位置推导仓根（任何 cwd 可跑），**先** `go test ./...`
  （阳性对照，证明门能变红）**再** `go run . -root "$root"`；`set -eu`，无 `|| true`、无 continue-on-error。
  `.github/workflows/ci.yml` 的 "D22 seven-ban + emoji scan" 步改为 `sh scripts/d22scan.sh`
  （同一 ubuntu/bash job，命令集合与改前逐字等价，未放宽）。调用契约同时写进 `main.go` 包注释与脚本头注释。
  【测量】`sh scripts/d22scan.sh` → `examined 194 production Go files` + clean + exit 0；
  三种误调用现全部非零：无 `-root` → `missing tools/d22scan/allowlist.txt` exit 2；
  `-root ../../../` → `missing go.mod` exit 2；仓根 `go run ./tools/d22scan -root .` → Go 模块错误 exit 1
  （这条在 `main()` 之前就被 Go 拦下，扫描器无法自截，只能靠 wrapper + 文档消除诱因）。
  【端到端种违规】临时建 `internal/llm/adaptertest/zz_ticket67_seed.go`（裸 goroutine + `start.Sub(time.Now())`）
  → wrapper 报 `2 finding(s)`、`examined 195`、exit 1；删除文件 → 回到 clean/194/exit 0。
  单测侧新增 `TestCheckRootRejectsBlindRoots`（4 个盲 root 子用例，含"d22scan 模块自己"这一形状）、
  `TestScanAloneIsNotAFalsifier`（钉住 `Scan()` 对空树返回 0 命中这一事实，防后人把 checkRoot 当冗余删掉）、
  `TestCheckRootAcceptsRealRepo`（同一谓词必须放过真仓，194 文件）。`cd tools/d22scan && go test ./...` 全绿。
  未做：AC#3（碰 `cmd/wisp/providers.go`，票 66 在飞，编排者裁定等其落地）、AC#4 门禁全清单复跑。
  备注：`go run` 会把子进程任意非零退出压成 1，故"有违规(1)"与"用法错(2)"经 `go run` 后不可分辨——
  wrapper 靠 `set -e` 兜住，但若将来有脚本想区分这两态，须直接跑编译出的二进制而非 `go run`。
- [2026-09-20T15:14:00Z] agent=agent-ticket67 did=补 `.gitattributes`：`*.sh text eol=lf`
  （`* text=auto` + 本地 `core.autocrlf=true` 会把 `scripts/d22scan.sh` 检出成 CRLF，
  shebang 与 `set -eu` 在 sh 下即坏；与既有 `*.go`/`*.md` 同形钉死）。
  `git check-attr eol -- scripts/d22scan.sh` → `lf`；`git show HEAD:scripts/d22scan.sh` 的 blob 本就读作 `\n`。
  【AC#4 门禁——本票已完成部分逐条复跑，原始退出码】
  `gofmt -l internal/llm/adaptertest tools/d22scan scripts` → 空；
  `go vet ./...`（主模块，含票 66 在飞的 observe/proc/cmd 改动）→ exit 0；
  `go vet ./...`（`tools/d22scan` 模块）→ exit 0；
  `go test -count=2 ./internal/llm/adaptertest/` → ok 2.644s；
  另跑 `go test -count=2 -race ./internal/llm/adaptertest/` → ok 3.684s（证明改造后的 helper 仍并发安全，未为了 Spawn 好写而变不安全）；
  `cd tools/d22scan && go test -count=2 ./...` → ok 0.487s；
  `cd tools/d22scan && go run . -root ../../` → `examined 194 production Go files` + `d22scan: clean`，**exit 0**（票面期望的 0 命中/exit 0 达成）。
  未做：AC#3（等票 66 收尾后动 `cmd/wisp/providers.go`）。本票范围 AC#1/AC#2 已交完。
- [2026-09-20T15:22:00Z] agent=agent-ticket67 did=**AC#2 交付后自挖出一条门覆盖面缺口，如实登记，未自行放宽/收紧**。
  探针（种完即删，`internal/llm/adaptertest/zz_probe.go`）：禁令 #1 的 AST 规则只匹配
  `*ast.GoStmt` 且 `Call.Fun` 是 **`*ast.FuncLit`** 的情形（`tools/d22scan/main.go` 的 `bare-goroutine` 分支），
  所以 **`go probeReader()`（具名函数、非字面量）扫出来是 `clean` / exit 0 / 195 文件**——
  一个同样无 name、无 owner、无 recover 边界的裸 goroutine，门**看不见**。
  我没有把它一起修掉，理由有两条，都需要你先拍板：
  ①把规则收紧成"任何 GoStmt"会立刻命中 `internal/observe/goroutine.go:281` 的 `go r.run(...)`，
  那是注释里明写的**全仓唯一受 sanctioned 的一条**，必须给它加 allowlist 豁免行——
  而裁定 1 说豁免只能你书面给；且 `internal/observe/` 是票 66 领地，我不碰。
  ②同一条收紧会再命中 `cmd/balldebug/main.go:231/259/440`（`go runHotkeyBridge(...)`、`go feedLevels(...)` ×2），
  即**真命中 3 处**，而 `cmd/balldebug/` 有代理在飞。
  实测计数：`grep -E '^[[:space:]]*go [a-zA-Z_(]' internal cmd`（排除 `_test.go`）= **4 处**，其中 1 处是 sanctioned、3 处待判。
  ⇒ **"d22scan clean" 的准确读法**：现有生产码里**没有以匿名闭包形式起的裸 goroutine**，
  不等于"没有绕过 `Registry.Spawn` 的裸 goroutine"。这条判据目前只对 FuncLit 形状为真。
- [2026-09-20T15:24Z] agent=orchestrator did=**票面对账 + 勾 AC#4（独立复现，非采信自述）+ Status→blocked-on-ticket:66**。
  ①我自己复跑 AC#4 六条：`cd tools/d22scan && go run . -root ../..` → `clean` 且自报
  **`examined 194 production Go files`**、**exit 0**；`go test ./...`（模块内 seeded-violation 阳性对照）通过；
  `go test -count=2 ./internal/llm/adaptertest/` → **ok 2.900s**；`gofmt -l` 空；`go vet` 干净；
  并核对 `allowlist.txt` 仍是 4 行 ⇒ **没有为凑绿加豁免**（裁定 1 兑现）。
  ②读改动认可三处判断（我没要求、它自发做对）：独立命名 registry 而非 `observe.Default`
  （测试协程不该挪动产品 watchdog 读的常驻名册数字）、就绪等待走 `observe.Timeout` 单调钟（避开 D22 禁令 #4）、
  cleanup 先于等待注册所以每条失败路径都收子进程；且它把 `h.Done()` 下的 panic 转成 `t.Fatal` 而不是静默挂 20s。
  ③**代理替我否证了一条我自己说过的话**：它报"CI 的 lint job 还有第二个独立红因（gofumpt 标 69 个文件）"，
  我用 CI 钉的 v0.7.0 本地复跑确认 ⇒ **我 14:59 在 `988c833` 里写的"CI lint job 第一次真绿"是错的**，
  我只验了一个步骤。已登记 **A27**（CI 五 job 全红，实测 `gh run view --json jobs`）并开**票 70** 承接；
  **commit message 无法改（已推），故此行为更正值。**
  ④它另外登记了一条我没让它修的覆盖洞：**ban #1 只匹配闭包字面量**，`go probeReader()` 这类具名调用扫不出
  （生产范围实测 4 处）⇒ 今天的 `clean` 含义是"无闭包字面量裸协程"。归票 70 AC#5，**我未默默扩语义**。
  next=票 66 收尾后：AC#3 三步（改 `providers.go` 三处字形 → 扩 ban #8 覆盖面且**排除注释** → 只许从严）。
