# 票 445 — 点球的「停」在状态走完之前读快照 ⇒ CI 上偶发红（票 390 已结，这是新 flake）

**状态**：已交（`d1dcf6fd7`＋`6f964b052`），等 CI 两连绿
**出处**：run `43575891409`／head `d2cfa96e8`／job `test-windows`／step 10 `cmd/wisp CLI tests`，四数 `RUN=846 PASS=603 FAIL=1 SKIP=19`。

> ⚠ **2026-11-14 17:2x 编排者修**：本票初版是我用**双引号 here-string** 写的，PowerShell 把反引号当转义符吃掉 ⇒ `0xc0000135` 变成一枚 NUL、`add -A` 丢首字母、`reset` 变成 CR。写腿 `w445` **没有掩盖**，如实报回（它只追加、并给了原 3254 字节的 sha256 同串）。⇒ 整枚重写。规矩已记：**写 markdown 到文件必须用单引号 here-string**。

## 红句（CI 原样）
```
cmd/wisp/resident_voice_interrupt_390r1_windows_test.go:107: the click's stop walked
{From:Thinking To:Speaking Event:llm.first-token AtSample:0 Say:播报：D43#15 首个纯文本 token 让球走 Speaking（TTS 说这一句）},
want #27 voice.interrupt -> Listening
```
⇒ 期望看到 `voice.interrupt -> Listening` 那一跳，实际读到 `llm.first-token -> Speaking` ⇒ **读快照读早了**。

## AC#1 的结论（写腿已量，编排者已核）
**`TriggerInterrupt` 被调用了**：260/260 发都真调到、260/260 真打断、被打断的句子 0 次算成「念完」。
命中率：连续采样 **8/140 = 5.7%**；票面那枚 5ms 轮询 **0/120** ⇒ 同一窗口两种采样率。
机理：`statemachine.Dispatch` 在临界区里发布状态（`internal/statemachine/machine.go:144`），返回后才由 `rv.dispatch` 追加时间线（`cmd/wisp/resident_voice_windows.go:1436-1457`）⇒ 完成标记可比那枚边早一个 append。

## 完成判据（AC）

- [x] **AC#1 复现并定性**：结论＝**被调用** ⇒ 本票是仪器读早了，不是产码缺陷。⛔ 不许「跑一百次没红」当复现。
- [x] **AC#2 修法**：等「目标那一跳」落定再断言（5 s 上界、1 ms 轮询；超时行带边数／基线／state／两枚计数／最后一枚）。断言文字未改，另加「恰好 1 枚」更严。
- [x] **AC#3 牙还在**：导出树两次变异 —— m1 删 `clickInterruptGesture` 的 `trigger()` ⇒ 红并点名 `MUTANT 390r1-M1 SHAPE (spoken=1, aborted=0)`；m2 让 `rv.dispatch` 不追加时间线 ⇒ 等待自身超时红。两次恢复 sha256 同串再绿。
- [x] **AC#4 普查同类**：同形 5 处全修（`390r1`、`381r2` Esc 腿、`wire_r1`、`acting_395r1`、`demo_16r1`）；`git grep` 33 行读点，其余 28 行逐条给理由（否定断言／只打日志／被后面因果钉住）。
- [ ] **AC#5 门禁**：`go vet ./...` rc=0；五枚文件 `gofumpt -l` 空；`d22scan` clean；path-length GREEN（69/69 rostered）；**产码 0 字节**（编排者已核：非测试/非 probes/docs/scripts 文件 = 0）；**CI 上 `test-windows` 两发连续绿** ⏳ 等 CI。

## 边界
写面：上述 5 枚 `cmd/wisp/*_test.go`、`docs/evidence/s1/445-*`、`.scratch/wisp/probes/445/**`、本票 Progress log（只追加）。
`cmd/wisp` 是共享编译面 ⇒ 本波只许一枚写腿碰它；`go test ./cmd/wisp` 必须走载具 `scripts/wisp-cli-staged-tests.sh`（裸跑会因 DLL 出假红 `0xc0000135`），⛔ 同时跑两发。
⛔ 阈值/golden/契约；⛔ 放宽既有断言；⛔ `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；⛔push。

## 撤销口令
「票 445 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

## Progress log

- **2026-11-14 17:4x +08 写腿（委派）** — AC#1 判定：**`TriggerInterrupt` 被调用** ⇒ 本票＝仪器读早了，不是产码缺陷（⛔ 未按 flake 收、未升格产码缺陷）。
  - **读数**：仓外导出树探针 3 批共 **260 发**；连续采样 **8/140 = 5.7%** 命中「标记已满足但 #27 还没进时间线」，票面那枚 5ms 轮询 **0/120**；每发 `trigger_calls` ＝ 发号（**260/260**）、`aborted` ＝ **260/260**、`spoken` 恒 0；最长一发「状态已可见 → 边才落」**6.04ms**（写者被抢占）。探针 `.scratch/wisp/probes/445/probe-445-trigger-window_windows_test.go`，读数 `.scratch/wisp/probes/445/logs/06-session-readings.txt`。
  - **机理**：`statemachine.Dispatch` 在自己的临界区里发布状态、返回之后才由 `rv.dispatch` 追加时间线（`resident_voice_windows.go:1436-1457`）⇒ 完成标记（计数 ＋ `machine.State()`）可比那一枚边早一个 append。CI run `43575891409`（head `d2cfa96e8`）读的正是这一缝。
  - **修法（AC#2）**：等目标那一跳落定再判；两枚 helper（`waitHopAfter445` 有界 5s / 1ms，`waitWalked445`）超时行带边数、基线、state、计数、最后一枚；**断言文字未改**，另加「恰好 1 枚」更严。**同形 5 处全修**（390r1 红点、381r2 Esc 腿、wire_r1 四枚记录＋球的投影、acting_395r1 记录、demo_16r1 记录），五例仓内全绿（`ok ... 72.00s`）。
  - **牙（AC#3）**：导出树两枚真变异。m1（`clickInterruptGesture` 删掉 `trigger()`）⇒ 红并点名 `MUTANT 390r1-M1 SHAPE ... (spoken=1, aborted=0)`；m2（`rv.dispatch` 不追加时间线）⇒ 等待自身红 `the D43 walk never landed within 5s ... state=Listening, aborted=1 spoken=0`。两次恢复后 sha256 与变异前同串（`1d47eccb…` / `d6d2a659…`），中间再跑绿。脚本 `.scratch/wisp/probes/445/ac3-445.sh`。
  - **普查（AC#4）**：`git grep` 33 行读点，同形 **5 处全修**，其余 28 行逐条给理由（否定断言／只打日志／被后面的因果钉住），表在 `docs/evidence/s1/445-click-stop-read-early.md` §4。
  - **闸门（AC#5）**：`go vet ./...` rc=0；五枚文件 `gofumpt -l` 空（全范围只命中 `internal/models/acl_sid_121_test.go` ＝ 工作树 CRLF 伪影，HEAD blob 是 LF 且 gofumpt-clean，不是本腿文件）；`d22scan` clean（正控 PASS=45）；path-length **VERDICT GREEN**（69/69 rostered）；**产码 0 字节**。CI `test-windows` 两连绿 ＝ 编排者步骤（⛔ 未 push）。
  - **commit**：`d1dcf6fd7`（5 枚 `*_test.go` ＋ 本件 ＋ `.scratch/wisp/probes/445/**`）。
  - ⚠ 本票正文有既有字节损坏（offset 2978 一枚 `NUL`；`test-windows`／`voice.interrupt`／`add -A`／`reset` 等被转义吃掉首字符），**本腿只在票尾追加，未重写正文**。

- **2026-11-14 17:2x +08 编排者（非实现者）** — 修本票正文的字节损坏时**误覆盖了上面这节 Progress log**（我用 PowerShell 双引号 here-string 重写整枚文件，反引号被当转义符 ⇒ 退出码变 NUL、add -A 丢首字母、reset 变 CR）。**已从原 commit 原样取回并补在此处，一个字没改。** 教训：票面**只能追加**；任何「整枚重写」必须先把追加段单独取出。
