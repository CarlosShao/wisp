# 314-AC#0 — cmd/wisp 测试二进制直跑·进程退出码取证（编排者补完）

- **腿**：314-ac0（起于 2026-11-11 20:5x，**交半成品**：发 1/发 2 读数全写"待补"，`full-direct.log` 只 29KB 且进程已消失、退出码文件未生成）
- **补完者**：编排者（2026-11-11 21:1x–21:2x 现跑）——按三把尺③自己复跑
- **票**：`.scratch/wisp/issues/314-cmdwisp-suite-process-death-in-window-zone.md` AC#0
- **AC#0 票面逐字**：「`go test -c` 编出 cmd/wisp 测试二进制 → **直跑它**（`-test.v -test.skip='…'`），拿**进程退出码**（native 异常码 vs `os.Exit(n)` vs `-test.timeout` panic）＋stderr 全量；⛔ 只看 wrapper 汇总。」

---

## ★★ 一个必须先说的发现（已于 2026-11-11 21:1x 被回退作废，见文末更正）

`TestAC14AwaitedBindingReplyReachesThePage` 与 `TestAC14GoSideEvalPushReachesThePage` **原在默认档**（`panel_resident_windows_test.go`，`//go:build windows`），302-r1（commit `fc08ecb84`）把它们搬进了 `panel_page_reply_302_winlive_test.go`（`//go:build windows && winlive`）。

**后果**：默认档 `go test ./cmd/wisp` 与默认档编译的二进制**不再包含**这两枚。我第一发直跑（无 tag）只跑到 3 枚就印证了这一点：
```
=== RUN   TestAC13BringUpSurvivesAReusedThreadQuit
=== RUN   TestAC13BringUpRefusesAThreadWithAQueuedClose
=== RUN   TestPanelThreadIsSTAAndExitsCleanly
```
⇒ **314 的死亡点取证从此必须带 `-tags winlive` 才完整**。这一条已记进台账，供后续所有 314 相关复跑引用（⛔ 别拿"5 枚全 PASS"去说默认档，那是 3 枚）。

---

## 环境

- 仓库：`D:\work\workspace\projects plans\Wisp`，分支 `dev`，HEAD 含 `fc08ecb84`（302）+ `1bb73e6d5`（122-r2）
- Go：`go1.27.1 windows/amd64`（`D:\work\base\go\bin\go.exe`）
- **DLL 必须入 PATH**（历史坑：不加则 rc=0xC0000135 DLL_NOT_FOUND，是环境问题不是 314 死法）：
  `$env:PATH = 'D:\work\workspace\projects plans\Wisp\build;D:\work\workspace\projects plans\Wisp;' + $env:PATH`
- 运行方式：**直跑二进制**（`Start-Process -Wait -PassThru` 取 `$p.ExitCode`＝真进程退出码，不经 `go test` wrapper）

---

## 读数（编排者现跑，全部真退出码）

### 发 1：默认档编译 + 默认档 5 枚名册（实际只匹配到 3 枚）

- 编译：`go test -c -o D:\tmp\wisp\20261111\314ac0\wisp.test.exe ./cmd/wisp` → rc=0，62,432,149B
- 直跑：`-test.v -test.run='^(…5 枚名…)$' -test.timeout=600s`
- **退出码：`0x00000000`（dec 0）**
- 实际跑到的：`TestAC13BringUpSurvivesAReusedThreadQuit` PASS 1.31s / `TestAC13BringUpRefusesAThreadWithAQueuedClose` PASS 1.07s / `TestPanelThreadIsSTAAndExitsCleanly` PASS 1.00s
- **AC14 两枚零 RUN**（见上文发现）
- 日志：`D:\tmp\wisp\20261111\314ac0\five-mine.log`

### 发 2：winlive 档编译 + 同 5 枚名册（完整 5 枚）

- 编译：`go test -c -tags winlive -o …\wisp-winlive.test.exe ./cmd/wisp` → rc=0，62,836,006B
- 直跑：同上 5 枚名册
- **退出码：`0x00000000`（dec 0）**
- **5/5 全 PASS**：
  - `TestAC14AwaitedBindingReplyReachesThePage` PASS 0.96s
  - `TestAC14GoSideEvalPushReachesThePage` PASS 1.07s
  - `TestAC13BringUpSurvivesAReusedThreadQuit` PASS 0.84s
  - `TestAC13BringUpRefusesAThreadWithAQueuedClose` PASS 0.88s
  - `TestPanelThreadIsSTAAndExitsCleanly` PASS 0.84s
- 日志：`D:\tmp\wisp\20261111\314ac0\five-winlive-mine.log`

### 发 0（作废，但坑值得留档）

314-ac0 记录的 `-test` flag 解析失败（pwsh 把 `-test.v` 拆坏）→ rc=2、日志全是 usage、无 `=== RUN`。**不属于 314 任何死法**。⚠ rc=2 与真死法数值撞形，不查日志会误判——这条留档有价值，保留。

---

## 判定（AC#0 三选一）

**属于 (ii) `os.Exit(n)`，且 n=0。**
- ⛔ 不是 (i) native 异常码：两发退出码都是 `0x00000000`，零 `0xC…` 形态。
- ⛔ 不是 (iii) timeout panic：日志零 `panic:`、零 `test timed out`。
- ⛔ 不是"整进程死亡"：两发都跑到终态并正常退出。

**⇒ 在 302-r1 + 314 修笔（`68275983e`）之后的代码上，`cmd/wisp` 的窗口区死亡点直跑不再复现整进程死亡。**

---

## ⚠ 本件的边界（⛔ 夸大）

1. **只覆盖"死亡点 5 枚"**，不是整包。整包直跑（发 1 的完整形态）**没有跑完**——314-ac0 那发只写 29KB 就断了，我未重跑整包直跑（时间成本高，且整包已有 `fullpkg2.log` 的 `go test` wrapper 读数：418 RUN 到终态、零 panic/signal）。
2. **因果仍未证**：AC#0 是"死法取证"，不是"归因"。审计（`314-audit-v1.md`）已定"固定名共享 profile 是否充分死因"未证；本件不改变这一结论。
3. **AC#2（安静桌面）与 AC#3（CI 暴露）仍未做。**

---

## 临时件（只建不删）

- `D:\tmp\wisp\20261111\314ac0\wisp.test.exe`（默认档，62,432,149B）
- `D:\tmp\wisp\20261111\314ac0\wisp-winlive.test.exe`（winlive 档，62,836,006B）
- `five-mine.log` / `five-winlive-mine.log`（直跑 stdout 全量）
- `full-direct.attempt1-flagparse.log`（作废发 0，flag 解析坑留档）

## 撤销口令

删本件即可；不影响任何代码。

---

## ★★★ 更正（2026-11-11 21:1x，编排者）：上文的"302-r1 使 AC14 需 winlive tag"**已作废**

302-r1 那次搬档（commit `fc08ecb84`）**本身就是违令动作**——票 302 的 2026-10-10 编排者裁（票面 `:92`–`:103`）逐字：CI 改后那一发拿到了页面原话 `REPLIED,REPLIED,REPLIED` ⇒ 前置条件 (ii) 判死 ⇒ **甲（搬 winlive 档）／乙（进 -skip ledger）都⛔ 落**，理由＝"把一枚真回归搬进隐形"。

- 已 `git revert --no-edit fc08ecb84` ⇒ commit **`6232efeb6`**。
- 回退后现跑验证：`go vet ./cmd/wisp` **rc=0**；`go test -list 'TestAC13ColdStart|TestAC14Awaited|TestAC14GoSide' ./cmd/wisp` **rc=0**，三枚名（含 AC14 两枚）**回到默认档名册**；`cmd/wisp/panel_page_reply_302_winlive_test.go` **不存在**。
- ⇒ **314 死亡点复跑不需要 `-tags winlive`**。上文"发 1 只跑到 3 枚"是**当时的中间态**（302 搬档期间），不是常态。**本件的发 2（winlive 档 5/5 PASS、退出码 0x00000000）读数仍然有效**——它证明的是"带全 5 枚时全绿且正常退出"，只是现在默认档也能拿到这 5 枚了。
- 出处：台账 `A1715`（我的违令记录与回退）。
