# 418 — `internal/ball/orbdebug.go` 没有 build tag，但引用了只在 `_windows.go` 里定义的 `Ball` ⇒ **Linux 侧编译必死**

**Status**: open（2026-11-10 09:1x 编排者建，从 CI 真红 + 本机交叉 vet 双证）

## 现量（⛔ 引用前先重跑）

**1. CI 真红（`test-core` job `129250565302`，跑 `ubuntu-latest`）**：
```
2026-11-09T15:21:21.0973564Z ##[error]internal/ball/orbdebug.go:39:10: undefined: Ball
2026-11-09T15:21:21.0981757Z ##[error]internal/ball/orbdebug.go:43:10: undefined: Ball
2026-11-09T15:21:21.1519468Z ##[error]Process completed with exit code 1.
```

**2. 本机复现（一条命令，零真窗）**：
```
$ GOOS=linux GOARCH=amd64 go vet ./internal/ball/
vet.exe: internal\ball\orbdebug.go:39:10: undefined: Ball
```

**3. 为什么**：
| 事实 | 读数 |
|---|---|
| `orbdebug.go` 的 build tag | **无**（`head -3` 只有 `package ball`） |
| `Ball` 类型定义在哪 | **只在 `internal/ball/ball_windows.go`**（`git grep -ln "^type Ball struct"` ⇒ 1 个文件） |
| 同目录正确样例 | `hotkey_reload.go` 头两行是 `//go:build windows` |
| 谁建的 | 票 65 落地腿，提交 `a2e535674`，**2026-11-07 13:02** |

**4. `orbdebug.go` 的五个函数，性质不同**：
```go
func OrbSurfaceName(attached bool) string              // :18  平台无关
func OrbMappedState(s statemachine.State) string       // :27  平台无关（转 orbStateFor）
func OrbDebugLine(attached bool, s statemachine.State) string  // :32  平台无关
func (b *Ball) DebugOrbAttached() bool                 // :39  ★ 要 Ball
func (b *Ball) DebugOrbState() statemachine.State      // :43  ★ 要 Ball
```

## ★为什么"拆文件"比"加一行 tag"更有价值（铁律 ⓪⓪⓪）

- **ⓐ 给整个文件加 `//go:build windows`**：一行搞定，但**前三个平台无关函数在 Linux 上消失**。它们本身没有任何平台依赖，消失是**为了迁就后两个**而付出的代价 —— 这是把"编译面"当判据，不是把"最有价值"当判据。
- **ⓑ 拆成两文件（主流做法）**：
  - `orbdebug.go`（无 tag）＝ 保留 `OrbSurfaceName` / `OrbMappedState` / `OrbDebugLine`
  - `orbdebug_windows.go`（`//go:build windows`）＝ 两个 `*Ball` 方法
  - ⇒ **平台无关的留在平台无关侧，平台相关的收进 `_windows.go`**，与本包既有命名约定（`ball_windows.go`、`hotkey_reload.go`）一致。
- **★我选 ⓑ**，理由：①与本包既有形状一致（`_windows.go` 后缀就是本仓约定）②不牺牲平台无关函数 ③`balldebug` 那类工具在 Linux 上仍能拿到 `OrbDebugLine` 的字符串（未来跨平台调试面）。
- ⚠ **腿自己判并写明理由**，⛔ 照抄我。若判 ⓐ，必须在票面写清"为什么宁可让三个平台无关函数消失"。

## AC

- [ ] **AC#1 拆文件落地**：`orbdebug.go` 只留三个平台无关函数；新建 `orbdebug_windows.go`（首行 `//go:build windows`）装两个 `*Ball` 方法。两文件都要有说明"为什么在这侧"的注释。
- [ ] **AC#2 正向判据**：`GOOS=linux GOARCH=amd64 go vet ./internal/ball/` **rc=0 且无输出**；`go vet ./internal/ball/`（本机 windows）**rc=0**。两个读数都贴。
- [ ] **AC#3 反向判据**：把 `//go:build windows` 从新文件删掉 ⇒ Linux vet **必须重新变红**（证明 tag 是真起作用的，不是碰巧）；还原 ⇒ 变绿。两个读数都贴。
- [ ] **AC#4 别只修这一个包**：`GOOS=linux GOARCH=amd64 go build ./...` 逐包扫，**列出所有 Linux 编不过的包**并逐枚判"是本票治、还是归已有票"（已知：`internal/panel` 4 枚、`internal/ball` 1 枚在票 230 AC#5 的在册名单里；**本票这一枚不在册** ⇒ 需把新读数追加进 230 的账）。
- [ ] **AC#5 门禁**：`gofumpt`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）空／`go vet ./internal/ball/` rc0／`go build ./internal/ball/` rc0／`go test ./internal/ball/ -count=1` 绿／`d22scan` clean／path-length GREEN。
- [ ] **AC#6 CI 侧验证**：推一发后 `gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs` 拉 `test-core` job 日志，**`undefined: Ball` 必须消失**，并贴那几行。

## 禁区
- ⛔ 动 `ball_windows.go` 的 `Ball` 类型本身（那是既有产码，本票只搬方法）
- ⛔ 删掉那两个 `*Ball` 调试方法（`balldebug` 在用）
- ⛔ 为了变绿给 Linux 造一个假的 `Ball` 类型（那是把闸拆了）
- ⛔ 动 thresholds / golden / allowlist.txt
- ⛔ 碰 `frontend/**`、`design/**`

## 撤销口令
- 「**418 撤回**」⇒ 本票作废。
## 补充：同源测试名（2026-11-10 09:5x，编排者现跑核）

同一根因（`undefined: Ball`）在 CI 还表现为：
```
--- FAIL: TestCrossVetForLinux (24.89s)
    crossvet_test.go:71: go.exe vet ./internal/proc/ ./internal/ball/ for GOOS=linux failed (exit status 1):
    vet.exe: internal\ball\orbdebug.go:39:10: undefined: Ball
```
⇒ **修 418 后 `TestCrossVetForLinux` 应同时转绿**。AC#2 正向判据请把这条也纳入。