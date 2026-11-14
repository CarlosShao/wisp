# 票 446 — 「这个 errno 到底有没有」是拿**英文句子**判的：`"The operation completed successfully."`

**状态**：待派（**一行级修法**，但是真 bug）
**来路**：写腿 `w216` 在普查「用户可见句子由英文错误文本选」时点名交回的**不在 216 射程内的同形**，它没动、只报告 —— 这是正确的做法。

## 现量（写腿给的，落地前自己复核）

`internal/tools/system_backend_windows.go:286` 与 `:426`，用**英文全串**比对一个 errno 的 `Error()` 结果来判断「有没有 errno」：

```
"The operation completed successfully."
```

`syscall.Errno(0)` 的 `Error()` 恰好就是这句英文。⇒ 在**非英文 Windows** 上（系统错误表被本地化）这句话不是英文，比较**永远不成立** ⇒ 判「有 errno」⇒ 走错分支。

**这不是文案问题，是行为问题**：判存在性应该用 `errors.Is(e, syscall.Errno(0))`，⛔ 不该拿句子。

## 完成判据（AC）

- [x] **AC#1 现量复核**：`git grep -n 'operation completed successfully'` 给出全部命中（含测试件），逐条判是「判存在性」还是「只是打印」。⛔ 只改报上来的两处就完事。
- [x] **AC#2 改成码判**：存在性判定换成 `errors.Is(e, syscall.Errno(0))`（或该 errno 的正确零值判据），**两处都改**；若同一文件还有同形，一并改并在证据里点名。
- [x] **AC#3 能变红**：一发**能红的检查** —— 把修好的判定临时改回字符串比较、并让该分支必须被走到（构造 errno=0 的路径）⇒ 必须红且点名；还原后 SHA256 同串、再绿。⛔ 只写「改了」。
- [x] **AC#4 门禁**：`go vet ./internal/tools` rc=0；`gofumpt -l` 本包空（⚠ 本机可能没装 gofumpt，如实记，⛔ 拿 gofmt 冒充）；`d22scan` clean；path-length GREEN；**产码改动只许这两处（或 AC#1 点名的同形）**。
- [ ] **AC#5 CI**：`test-core` 与 `test-windows` 各一发含本修复的**完成**运行（给 run 号 ＋ `git merge-base --is-ancestor <fix> <head>` ＋ 完成时刻）。

## 边界
写面：`internal/tools/system_backend_windows.go`（及 AC#1 点名的同形）、`docs/evidence/s1/446-*`、`.scratch/wisp/probes/446/**`、本票 Progress log（只追加）。
⛔ `cmd/wisp/**`、`internal/panel/**`、`internal/llm/**`（别的腿的写面，派单前现查）；⛔ 阈值/golden/契约；⛔ 放宽断言；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。

## 撤销口令
「票 446 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

## Progress log（只追加）

### 2026-11-14 18:3x +08 — 写腿（agent-30），基线 `447929bfe`

- **AC#1 全量 132 处**（`git grep -n 'operation completed successfully' HEAD`）：`*.go` **2** 处（`internal/tools/system_backend_windows.go:286` / `:426`，**两处都是判存在性**）、`*_test.go` **0** 处、非产码 **130** 处（probe/CI 日志 121、已落档证据 5、票面 4，全是打印/散文）。同形族 `Err(or)?().*(==|!=)\s*"` 全仓 `*.go` 只有这 2 处 ⇒ **M=2，改了 2，无漏改**。
- **AC#2**：判定收成一个函数 `errnoPresent(e error) bool { return e != nil && !errors.Is(e, syscall.Errno(0)) }`，两处调用点（`:301` / `:441`）都改走它；新增 import `errors`、`syscall`。句子在 `internal/tools/*.go` 里 **0 命中**。
- **AC#3 能变红，两向**：探针 `probes/446/r1/probe446_errno_presence_windows_test.go`（`go test -overlay`，不在 `internal/tools/` 留测试件），10 格其中 5 格在 errno=0 路上。A 修复后 PASS；B 判定改回句子比较 ⇒ **FAIL 并点名** `errno-zero-localized-text`（`操作成功完成。`）；C 还原 ⇒ **sha256 同串** `8CAD96…38845` 且 PASS。调用点回退方向由 `gate-446-static.ps1` 抓：修复后 GREEN（0 句子命中 / 2 调用点），`:441` 回退 ⇒ **RED 三行点名 line 441**。
- **AC#4 门禁**：`go build`/`go vet ./internal/tools` rc=0；`d22scan` clean rc=0；path-length **GREEN** rc=0（20122 tracked / 69 over-budget / 69 rostered / 0 not in roster）；`gofumpt` **本机没装、离线装不上 ⇒ 这一格留空**，`gofmt -l internal/tools` 空输出只作旁证、不冒充。测试：`go test ./internal/tools -count=1` rc=0。
- **AC#5 CI**：**留白**（本腿无 CI 权限，不声称绿）。
- **前提核查（票面那半句要改口径）**：`Errno.Error()` 的英文来自 `syscall_windows.go:169` **显式索要 en-US**，失败才在 `:171` 退回系统默认语言；本机界面语言 0x0409 ⇒ 旧比较**在本机碰巧成立**，bug 在一台拿不到 en-US 表的机器上才现形。本地化文本已量出：`FormatMessage(errno=0, zh-CN 0x0804) = "操作成功完成。"`。**真机复现那一支本腿做不到**（需要没有 en-US 表的机器），已在证据里单列。
- **证据**：`docs/evidence/s1/446-errno-presence.md`；器具与日志 `.scratch/wisp/probes/446/r1/`（`32-ac3-sequence.log`、`04-locale-reading.log`、`50-d22scan.log`、`51-path-length.log`）。
---

## 编排者收件（2026-11-14 18:3x，非实现者）：AC#1–#4 翻，AC#5 等 CI

- **尺1**：`be0f599f3`（改）＋`3dedc7796`（证据更正）两枚 `cat-file` rc=0。
- **尺2**：`docs/evidence/s1/446-errno-presence.md` **13554 B**。
- **尺3 我自己复跑**：`go test ./internal/tools` ⇒ **`ok … 41.459s`，RC=0**；`sh scripts/d22scan.sh` ⇒ **RC=0，clean**；`gofumpt -l internal/tools` ⇒ **空输出＋rc=0**（腿说本机没装、我把这格补上了）。
- **我自己 grep**：`internal/**`＋`cmd/**` 的 `*.go` 里那句英文 **0 命中**；`errnoPresent` 定义在 `:270`、两处调用 `:301`/`:441`。
- **腿纠正了我票面前提的一半（照记）**：`Errno.Error()` 那句英文**不是 Go 静态表**，是 Go 在 `syscall_windows.go:169` **显式索要 en-US**、失败才退回系统语言 ⇒ **本机（en-US）旧写法碰巧成立**，bug 只在拿不到 en-US 表的机器上现形。它量出 `FormatMessage(errno=0, zh-CN) = "操作成功完成。"`，并**明确不声称真机复现过**。⇒ 修它的理由因此改写为：判据不该是一句话（换语言/换构建就变），且 `errors.Is` 会**拆包**而字符串比较只看得见最外层。
- **AC#5 留白**：等我推后 CI 的 `test-core`/`test-windows` 各一发**含 `be0f599f3` 且完成**的运行（给 run 号＋`merge-base`＋完成时刻）。
