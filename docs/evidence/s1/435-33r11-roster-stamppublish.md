# 435 — 33r11 名册漏登 `panelHolder.stampPublishLocked`：判定与交付证据

- 腿：`435-33r11`（写面独占 `cmd/wisp/**`，本波另一条腿 437 只碰 `internal/panel/**`）
- 起手锚：`HEAD=9d94fbadd`（分支 `dev`）
- 落档时间：`2026-11-13 09:2x +08`
- 临时件：`D:\tmp\wisp\20261113\435-33r11\`

---

## 0. AC#0 只读取证（本格零写）

### 0a. 复现（carrier：`D:\work\soft\Git\bin\bash.exe`）

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestLockedSuffixNeverTakesTheLockItself33r11$'"
```

**读数**：`--- FAIL: TestLockedSuffixNeverTakesTheLockItself33r11 (0.07s)`，`rc=1`，
**测试本体耗时 0.07s**，整条 carrier 端到端 **11.29s**（含 `go test -c` 编译与三枚 DLL 暂存）。

**红句逐字（只有 1 条，`panel_locked_naming_33r11_windows_test.go:254`）**：

```
panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11, or renamed
```

**`:245` 那 4 行是 `t.Logf` 普查日志（⛔ 红句）**，逐字：

```
33-r11 census: 4 Locked-suffixed func(s) across 72 production source file(s) of package main
  PanelManager.setPriorFocusLocked in panel_host_windows.go line 928: 0 lock operation(s) in its own body []
  panelHolder.stampPublishLocked in resident_windows.go line 1008: 0 lock operation(s) in its own body []
  waveOutSpeaker.releaseLocked in voice_out_windows.go line 338: 0 lock operation(s) in its own body []
  waveOutSpeaker.unprepareLocked in voice_out_windows.go line 376: 0 lock operation(s) in its own body []
```

⇒ 票面末尾"更正记录"复核**成立**：普查 4 行 / 红句 1 条；`PanelManager.setPriorFocusLocked`（`:87`）、
`waveOutSpeaker.releaseLocked`（`:88`）、`waveOutSpeaker.unprepareLocked`（`:89`）**三枚已在名册**⇒ 绿。
**本票范围确为一枚。**

### 0b. 名册现有登记（`:86` 起，修前 3 条）

`lockedNamingRoster33r11` 是 `map[string]string`，key＝`类型.函数名`，value＝"为什么它自己不持锁"。修前 3 条逐字：

```go
"PanelManager.setPriorFocusLocked": "requires the caller to hold m.mu; body takes no lock",
"waveOutSpeaker.releaseLocked":     "requires the caller to hold s.mu; body takes no lock",
"waveOutSpeaker.unprepareLocked":   "requires the caller to hold s.mu; body takes no lock",
```

格式＝`"<类型>.<方法名>": "requires the caller to hold <锁>; body takes no lock"`，gofumpt 按 key 长度对齐冒号。

### 0c. `stampPublishLocked` 全部调用点，逐个读前后 10 行

`git grep -n 'stampPublishLocked'` 命中产码 **4 处**（定义 1 + 调用 3），另 1 处是测试里的源钉。

| # | 位置 | 调用者是否持锁 | 逐字依据 |
|---|---|---|---|
| 1 | `cmd/wisp/resident_windows.go:985`（`start`） | **持锁** | `:979 h.mu.Lock()` / `:980 defer h.mu.Unlock()` 开在函数首；`:984 h.rp = h.spawn(...)` 后紧接 `:985 h.stampPublishLocked()`，仍在 `defer` 释放之前 |
| 2 | `cmd/wisp/resident_windows.go:1044`（`toggle`） | **持锁** | `:1033 h.mu.Lock()` / `:1034 defer h.mu.Unlock()`；`:1043 h.rp = h.spawn(...)` 后紧接 `:1044 h.stampPublishLocked()` |
| 3 | `cmd/wisp/resident_windows.go:1104`（`show`） | **持锁** | `:1093 h.mu.Lock()` / `:1094 defer h.mu.Unlock()`；`:1103 h.rp = h.spawn(...)` 后紧接 `:1104 h.stampPublishLocked()` |

**三个调用点无一例外都在 `h.mu.Lock()` / `defer h.mu.Unlock()` 的临界区内**，
且每个都紧跟 `h.rp = h.spawn(...)` 那一行——正是"新线程刚换上、必须立刻把 publisher 盖上去"的窗口。
函数自己的 doc（`:1005-1007`）逐字写着 `Called with h.mu held (every spawn site holds it)`。

**函数体内一次加锁都没有**（普查读数 `0 lock operation(s)`），它唯一碰到的锁是
`residentPanel.setSnapshotPublish` 内部的 `rp.mu`（`cmd/wisp/panel_resident_windows.go:555-557`），
那是**面板线程自己的**另一把互斥量，与 `h.mu` 不是同一把，作用是把字段写和 `publishOnShow`（`:565-567`）的读串起来。

**第四条命中 `cmd/wisp/panel_resident_pump_windows_test.go:240`** 是 `strings.Contains(boot, "h.stampPublishLocked()")`
的源钉（票 65 的 resync-on-show 钉），⛔ 调用点，改名会连带红它。

### 0d. 判定：**甲（该登记）**

**一句话理由**：三个调用点全部在 `h.mu` 的临界区内，函数体不碰 `h.mu`，所以 `Locked` 后缀说的是真话——
它是"调用者已持 `h.mu`，本函数不得再拿 `h.mu`"的契约；名册是漏登，不是名字在骗人。

⇒ **⛔ 判乙**。若判乙改名，等于把一个语义正确的名字去掉后缀，反而让"这三个 spawn 窗口必须持锁"这件事从名字上消失。

---

## 1. AC#1 修复（只修这一处）

### 1a. 改动

`cmd/wisp/panel_locked_naming_33r11_windows_test.go`：

1. 名册表加 1 行（`:101`）：
   ```go
   "panelHolder.stampPublishLocked":   "requires the caller to hold h.mu; body takes no lock on h.mu (rp.setSnapshotPublish takes the panel thread's own rp.mu)",
   ```
2. 名册上方的文档块加 1 段（`:83-93`），写明**为什么它自己不持锁**：三个调用点位置、
   `h.rp = h.spawn` 的窗口关系、以及它唯一碰到的 `rp.mu` 是另一把锁。

**未改任何产码**（`resident_windows.go` 一字节未动）⇒ 影响面＝0 个调用点。

### 1b. 修后复跑（carrier 同一条命令）

**读数**：`--- PASS: TestLockedSuffixNeverTakesTheLockItself33r11 (0.04s)`，`rc=0`，
测试本体 **0.04s**，carrier 端到端 **18.49s**。普查仍报 `4 Locked-suffixed func(s) across 72 production source file(s)`。

### 1c. 突变控制（★ 仓外磁盘改写，⛔ `-overlay`）

**为什么要绕这一圈**：这枚检查用 `os.ReadDir(".")` + `os.ReadFile` + `go/parser` **读磁盘源码**，
但名册 `lockedNamingRoster33r11` 是**编译进二进制的 Go 变量**。两个面必须分开证。

**仓外树**：`D:\tmp\wisp\20261113\435-33r11\module-tree\`
（`go.mod`/`go.sum`/`deps.toml` + `cmd/` + `internal/` + `frontend/`（去 `node_modules`）+ `forks/`，
24.8 MB / 1484 文件；三枚 DLL 与 `wisp.test.exe` 同目录；cwd＝`module-tree\cmd\wisp`，与 carrier 同约定）。
该树内 `go test -c ./cmd/wisp/` 成功（24.6s，`rc=0`）。

| 突变 | 改什么 | 是否重编译 | 读数 | 说明 |
|---|---|---|---|---|
| **M0 基线** | 不改 | — | **PASS**，`rc=0`，4 命中 / 72 文件 | 仓外树与仓内一致 |
| **M1** | 仓外树 `cmd/wisp/` 新增 `zz_mut435_probe_windows.go`，里面一个 `func (h *panelHolder) mut435ProbeLocked() {}` | **⛔ 不重编译**（二进制 sha256 前后同为 `4B594E99…A349E`） | **FAIL**，`rc=1`，普查 `5 Locked-suffixed func(s) across 73 production source file(s)`，红句 `panelHolder.mut435ProbeLocked (zz_mut435_probe_windows.go line 7) carries the "Locked" suffix but no roster entry` | ★**仪器有效性主证**：同一个二进制、磁盘多一个文件，普查立刻从 4/72 变 5/73 并点名新文件。删掉该文件后同一二进制**回 PASS** |
| **M2** | 仓外树 `resident_windows.go` 里给 `stampPublishLocked` 体首插 `h.mu.Lock()` / `defer h.mu.Unlock()` | **⛔ 不重编译**（同一二进制） | **FAIL**，`rc=1`，该行普查变 `2 lock operation(s) in its own body [h.mu.Lock h.mu.Unlock]`，红句 `:270 … is named as if the caller holds the lock while its own body takes it ([h.mu.Lock h.mu.Unlock]) - that is the trap ticket 33 recorded …` | ★**名册⛔免死金牌**：已登记的条目体内一旦真加锁，仍然红。还原后同一二进制回 PASS |
| **M3** | 仓外树名册表**删掉**本票新增的那 1 行 | **必须重编译**（名册是编译进去的） | **FAIL**，`rc=1`，红句回到 `:265 panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry` | ★**本票的突变齿**：删登记必红。还原该行重编译后回 PASS |

**★ 仪器有效性的确证（三问三答）**：

1. **它真的读磁盘吗？** ⇒ M1 用**同一个二进制**（sha256 未变）证明：磁盘上多一个 `_windows.go`，
   普查分母 `72→73`、命中 `4→5`、并点名 `zz_mut435_probe_windows.go line 7`。若它读的是二进制里的快照，
   这两组数字**不可能**跟着磁盘变。
2. **`-overlay` 真的无效吗？** ⇒ 两条对照：
   - **对照 A（overlay 确实到达编译器）**：把 overlay 源写成故意不可解析的 `func broken435( {`，
     `go test -c -overlay=…` **编译失败** `syntax error: unexpected {, expected )` ⇒ overlay 通路是通的，不是没生效。
   - **对照 B（到达编译器但普查看不见）**：overlay 源换成合法的 `func (h *panelHolder) ovl435ProbeLocked() {}`，
     编译**成功**（`rc=0`），用该二进制跑测试：普查仍报 **`4 Locked-suffixed func(s) across 72 production source file(s)`**、
     **PASS** ⇒ 编译器看见了第 5 个函数，普查**看不见**。这正是台账 `A1804` 说的"静态检查读磁盘，overlay 无效"。
3. **突变齿会不会假绿？** ⇒ M3 删登记后**重编译**仍红，说明名册条目是承重的，不是装饰。

---

## 2. AC#2 越界检查

`git diff --name-only` 逐行（**我只动最后那一行**，其余 13 行是起手就已脏的存量改动，非本腿所为）：

```
.scratch/wisp/issues/426-…-never-runs.md                        （存量脏，非本腿）
.scratch/wisp/issues/427-…-hooks-logf.md                        （存量脏，非本腿）
.scratch/wisp/issues/428-…-announcement-audio.md                （存量脏，非本腿）
.scratch/wisp/issues/429-…-return-pages-on-its-own.md           （存量脏，非本腿）
.scratch/wisp/issues/430-…-never-becomes-an-overdraw.md         （存量脏，非本腿）
.scratch/wisp/issues/431-…-boot-budget-not-on-the-product.md    （存量脏，非本腿）
.scratch/wisp/issues/432-…-live-slot-count-premises-moved.md    （存量脏，非本腿）
.scratch/wisp/issues/433-…-handle-outlives-the-child.md         （存量脏，非本腿）
.scratch/wisp/issues/434-…-runconfigset-now-creates-the-config-too.md（存量脏，非本腿）
.scratch/wisp/issues/436-…-supposed-to-be-off.md                （存量脏，非本腿）
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go（存量脏，非本腿）
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt      （存量脏，非本腿）
cmd/wisp/panel_locked_naming_33r11_windows_test.go              ← ★本腿唯一改动
frontend/fixtures/composer-states.html                          （存量脏，非本腿）
```

**判定**：无 `docs/PLAN.md`、无 `docs/specs/**`、无 `internal/observe/thresholds.go`、无 golden、
无 `tools/d22scan/allowlist.txt`、无 `design/**`。`internal/panel/**` **零触碰**（437 腿的地界）。
`cmd/wisp/` 内只有一枚文件改动，与"本波独占 cmd/wisp"一致。

---

## 3. AC#3 门禁四读数

| # | 命令 | 读数 |
|---|---|---|
| 1 | `go build ./cmd/wisp/` | **`rc=0`**，零输出 |
| 2 | `go vet ./cmd/wisp/` | **`rc=0`**，零输出 |
| 3 | `gofumpt -l cmd/wisp/panel_locked_naming_33r11_windows_test.go` | **零输出**（`rc=0`） |
| 4 | `scripts/d22scan.sh`（经 `D:\work\soft\Git\bin\bash.exe`） | **`rc=0`**，`d22scan: clean - no D22 ban violations`；`examined 480 production Go files`；自检 `top-level: PASS=45 FAIL=0 SKIP=0`，`ok github.com/CarlosShao/wisp/tools/d22scan 53.608s` |

`build/mut_rvw.go` **未删未动**；门禁用的是 `go build ./cmd/wisp/`，⛔ `go build ./...`。

---

## 4. AC#4 与票 33 / 11 的关系

- 本票**只处理一枚名字**（`panelHolder.stampPublishLocked` 入名册）。
- **⛔ 不翻**票 33、票 11 的任何 AC 框；**⛔ 宣称** `…Locked` 全族已审。
- 名册的含义保持票面原样：**它是"审过的名单"，⛔ "消红工具"**。
  M2 那条突变正是这句话的机器证明：名字上了名册，体内一加锁**照样红**。

---

## 5. 一条留给下一个人的观察（非本票范围，未动）

`panel_resident_pump_windows_test.go:240` 用 `strings.Contains(boot, "h.stampPublishLocked()")` 做源钉。
将来若这个名字真的该改（本票判它**不该**改），那枚钉会一起红——改名的人需要同时处理它。
本票未触碰该文件。
