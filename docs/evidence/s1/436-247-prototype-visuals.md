# 436 — `prototypeVisuals` 出厂 `true`：票 68 AC#2 做的，还是有人误改？

- 票：`.scratch/wisp/issues/436-the-247-test-finds-prototype-visuals-on-while-the-shipped-default-was-supposed-to-be-off.md`
- 腿：`436-247visuals`（写）
- 判：**甲——票 68 AC#2 有意翻转，247 的断言过时**
- 写面：`cmd/wisp/resident_audio_247_windows_test.go`（**只动断言这一侧**；`internal/ball/statevisual.go:110` 一字节未改）

## 1. 红句复现（AC#0，只读）

尺（必须经 carrier，`D:\work\soft\Git\bin\bash.exe`；PATH 上的 `bash` 是 WSL 转发器、跑零枚测试）：

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC247HandingTheLevelToTheSeamIsNotVisibility$'"
```

修前逐字：

```
=== RUN   TestAC247HandingTheLevelToTheSeamIsNotVisibility
    resident_audio_247_windows_test.go:355: prototypeVisuals is ON in this assembly: the shipped default is off, and a leg that says it turned it on has crossed AC#10
--- FAIL: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.00s)
```

## 2. 票 68 AC#2 原文（逐字，`68-ball-default-visuals-parity.md:57-65`）

> - [ ] **AC#2 默认值翻转（需桌面复测）**：`prototypeVisuals` 默认改为 `true`，
>   使**默认构建 == owner 签收的样子**；`EnablePrototypeVisuals(false)` 与 `balldebug -frozen` 保留为"对照旧冻结规格"的逃生门。
>   连带把 `internal/ball` 里**断言旧默认**的测试改到位：⚠ **不得**为了让测试变绿而删除断言或放宽阈值——
>   必须**逐条**把断言迁移到"新默认下应有的行为"，并在 log 里列出"这条原来断什么、现在断什么、为什么等价"。
>   若某条测试的存在意义就是"钉住旧的冻结规格"，**保留它**并显式 `EnablePrototypeVisuals(false)`，别删。

⇒ AC#2 **逐字要求**把 `prototypeVisuals` 默认改成 `true`，并要求把"断言旧默认"的测试**逐条迁移**（⛔ 删断言／放宽阈值）。这正是本票做的事。

### 票 68 是否 done

- **文件名不带 `-done`**：`.scratch/wisp/issues/68-ball-default-visuals-parity.md` ⇒ **票面未结**。
- 但 **AC#2 那一格已经在代码里落地了**：票面 Status 仍是 `review`、AC#2 框**未勾**（还等 AC#3 的真桌面复测与 owner 目验签收），而翻转本身由 `68-r1` 腿于 2026-11-01 交付并入库。
- ⚠ 本票**不翻票 68 任何框**、不改票 68 一字（票面 AC#4 要求）。"票 68 未 done"与"AC#2 已落地"两件事并存，是票 68 自己的账。

## 3. 哪个提交把它改成 true，它说了什么

```
$ git log -S 'var prototypeVisuals = true' --oneline -- internal/ball/statevisual.go
4b9efc889 票 68 落地腿 68-r1
```

`4b9efc8897284474caf994f7225b5c359e689084`，CarlosShao，`Sun Nov 1 22:50:00 2026 +0800`。说明逐字（节录翻转那一行）：

> 票 68 AC#2 落地（R15 三答案齐，registry A1378「全按推荐」）：
>
> - statevisual.go: prototypeVisuals 默认翻 true——默认构建即 owner 签收的样子；
>   EnablePrototypeVisuals(false)/balldebug -frozen 保留为冻结规格逃生门。

同一个 commit 还按票 68 AC#2 的要求做了"测试逐条迁移"：`tokens_test.go` 四锚点迁新默认读数、旧冻结数字挂进显式 `EnablePrototypeVisuals(false)` 块（**零删除断言、零放宽**）、票面点名的地雷 `TestBallLiveAudioLiquidGate` 补 `EnablePrototypeVisuals(false)`。

旁证（台账，`docs/reports/pending-and-issues.md:30042` `A1378`、`:30060` `A1379`）：

> 3. ★**R15 全按推荐**＝**Settling alpha 0.33 接受／停靠露出 82% 接受／命中半径 21px→26px 修正**——**票 68 解锁**（prototypeVisuals 翻转的 R15 材料齐）

> - ★★★**AC#2 翻转落地**：★**statevisual.go `var prototypeVisuals = true`**（注释按 A1378「R15 全按推荐」重写）——**默认构建从现在起画的就是 owner 签收的玻璃球**

⇒ 翻转**有主、有据、有 owner 批复**（A1378 逐字「R15 全按推荐」），**不是误改**。判**甲**。

## 4. 判定与写面

**判甲**：票 68 AC#2 已于 `4b9efc889` 有意把出厂值翻成 `true`（`internal/ball/statevisual.go:110`），且 `internal/ball/tokens_test.go:158-160` 从库侧钉住了新默认：

```go
if !PrototypeVisualsEnabled() {
	t.Fatal("prototypeVisuals must default to ON after ticket 68 AC#2")
}
```

⇒ `cmd/wisp/resident_audio_247_windows_test.go` 修前（`HEAD` 版 `:354-356`）那条"出厂是关"的断言**过时**。**改断言，不动产码。**

## 5. AC#10 那个承诺怎么保住的（⛔ 不是删掉 if）

原断言其实把**两件事**绑在一个 flag 上：①"交接电平 ≠ 让像素可见"（AC#10 的真承诺）②"出厂默认是关"（票 68 的样式决定）。票 68 拿走的是②，①必须留下。

改后的测试守两条**与默认值朝向无关**的读数（`cmd/wisp/resident_audio_247_windows_test.go:392-417`）：

1. **交接电平不许移动可见性开关**（`:392-402`）：先记 `modeAtAssembly := ball.PrototypeVisualsEnabled()`（`:383`），交接一枚电平后要求 `PrototypeVisualsEnabled()` **未变**。谁在交接路径上顺手 `EnablePrototypeVisuals(true)`（票 247 AC#10 逐字禁的那一手）或反手关掉，都当场红。
2. **交接不许依赖这个开关**（`:404-417`）：`EnablePrototypeVisuals(false)` 后**同一枚标量仍须到达目的地**。若哪天有人把交接做成"只有像素被允许动才投递"，红。

⇒ 承诺守的是**形状**而不是**朝向**：默认值以后再怎么被票 68／后续样式裁定改，这两条都还站着。而"像素到底动不动"这半格**明写在测试注释里随票 68 AC#2 关闭**（`:351-359`），并注明库侧钉子 `tokens_test.go:158` 的位置——本测试不重复钉、也不冒充钉住。

测试自己用 `t.Cleanup` 把模式复原（`:378-379`），所以"强制关掉"那条读数不会泄漏给同包其他用例。

## 6. 突变控制（真跑，出仓 overlay）

仪器：`go test -overlay`，变异体落 `D:\tmp\wisp\20261112\436-247visuals\`（仓外）。经 carrier 的 `go test -c` 生效（`GOFLAGS=-overlay=...`）。

| # | 变异 | 读数 |
|---|---|---|
| A | 交接路径上顺手 `ball.EnablePrototypeVisuals(true)`（＝票 247 AC#10 逐字禁的那一手） | **FAIL** `resident_audio_247_windows_test.go:416: the second handoff moved prototypeVisuals off -> on` |
| B | 交接依赖开关：`if out != nil && ball.PrototypeVisualsEnabled() { out(level) }` | **FAIL** `resident_audio_247_windows_test.go:413: the destination saw 1 levels with prototypeVisuals off, want 2` |
| C | 把断言改回去（＝用 `HEAD` 版测试文件覆盖回去） | **FAIL** `resident_audio_247_windows_test.go:355: prototypeVisuals is ON in this assembly: ...` |

C 就是"把修好的那处改回去 ⇒ 必红"那条尺：复原后逐字回到修前的红句。A／B 是新增的两条读数各自的牙——**不是装饰**。

## 7. 越界（AC#2）

`git diff --name-only`（本腿 commit 前现跑）：

```
.scratch/wisp/issues/426-warm-idle-mic-close-counted-once-against-two-because-the-second-defence-never-runs.md
.scratch/wisp/issues/427-the-12r2-ast-check-still-looks-for-marksettling-inside-openscope-but-it-moved-into-scopeoptions-hooks-logf.md
.scratch/wisp/issues/428-the-26p-no-aec-echo-test-sees-zero-episodes-so-the-judge-never-received-the-announcement-audio.md
.scratch/wisp/issues/429-the-384-watchdog-never-sees-the-tree-back-under-the-cap-because-gc-does-not-return-pages-on-its-own.md
.scratch/wisp/issues/430-the-365-stuck-task-probe-never-goes-stuck-because-the-tiny-context-window-never-becomes-an-overdraw.md
.scratch/wisp/issues/431-the-two-ac1-resident-leg-tests-red-on-a-hardcoded-boot-budget-not-on-the-product.md
.scratch/wisp/issues/432-the-three-258-hotkey-tests-red-on-first-real-execution-because-the-old-binding-and-live-slot-count-premises-moved.md
.scratch/wisp/issues/433-the-246-dev-leg-test-is-stably-red-on-its-own-tempdir-cleanup-failing-because-a-handle-outlives-the-child.md
.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go
.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt
cmd/wisp/firstrun_198_test.go
cmd/wisp/resident_audio_247_windows_test.go
frontend/fixtures/composer-states.html
```

⚠ **上面只有 `cmd/wisp/resident_audio_247_windows_test.go` 是本腿的**。其余全是**别人在飞/既有**的脏件，本腿一枚没碰、一枚没进 commit：

- `cmd/wisp/firstrun_198_test.go` 是**票 434 的在飞改动**（`git diff` 显示它在改 `ensureFirstRunConfig` 的 AST 调用集，与本票无关）。
- `.scratch/wisp/issues/426-433` 的 `D`（删除）与 `frontend/fixtures/composer-states.html` 的 `M`、`.scratch/wisp/probes/291/**`、`346-a2-dispatch.txt` 都是**进场前就在工作树里的**，本腿只读未写。
- 本腿 commit 用**显式 pathspec**，只含自己那两枚文件（测试＋本证据件）。

禁改面（`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`）**零命中**。`build/mut_rvw.go` 未删。

## 8. 门禁四读数（AC#3）

```
$ go build ./cmd/wisp/
rc=0

$ go vet ./cmd/wisp/
rc=0

$ "$(go env GOPATH)/bin/gofumpt" -l cmd/wisp/resident_audio_247_windows_test.go
（空输出）

$ D:\work\soft\Git\bin\bash.exe -c "scripts/d22scan.sh"
d22scan: clean - no D22 ban violations; ...
rc=0
```

（d22scan 走 `scripts/d22scan.sh`，⛔ 不是 `go run ./tools/d22scan`。）

整族复跑（`-test.run '^TestAC247'`）：`TestAC247*` 全 PASS（含本枚），`PASS` rc=0。
修后本枚逐字读数：

```
resident_audio_247_windows_test.go:419: AC#10 reading: prototypeVisuals_at_assembly=true levels_reaching_the_ball_seam=2 levels_with_the_mode_forced_off=2 (a number arriving is not a pixel moving, whichever way the default points; the pixel half is closed by ticket 68 AC#2, commit 4b9efc889)
```

## 9. 与票 68／247 的关系（AC#4）

- **票 68 AC#2 已落地，本格"像素动不动"那半随它关闭**：默认值由 `4b9efc889` 翻成 `true`，库侧钉子 `internal/ball/tokens_test.go:158`。本票**不翻票 68 任何框**、不改票 68 一字（AC#2 框仍归票 68 自己在 AC#3 真桌面复测与 owner 目验签收后勾）。
- **票 247 AC#10 的承诺由本票改形续保**：原来靠"默认是关"证"看不见"，现在靠"交接不移动开关"＋"交接不依赖开关"证"这是两件事"。票 247 的框**一枚未翻**。
- ⛔ 本票没做：`internal/ball/statevisual.go:110` 零改动（判甲，写面在断言侧，**两边只动一边**）。
