# 436 — `TestAC247HandingTheLevelToTheSeamIsNotVisibility` 红：`prototypeVisuals` 出厂值是 `true`，而断言要求 `false`——票 68 把它打开了还是有人误改？

- Status: **done**（写面＝`internal/ball/statevisual.go:110` **那一行**，或 `cmd/wisp/resident_audio_247_windows_test.go:354-356` 的断言——⚠ **先只读取证，判清是"票 68 有意打开"还是"误改"，再定改法**）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11`）＋台账 `:32512` 存量真红
- 关联：票 247（本测试宿主，AC#10／AC#2）· **票 68（AC#2，翻转那个默认值的票）**· 票 29（L1/L2 记忆面）· 台账 `A1806`

## 这是什么（人话）

有个开关叫 `prototypeVisuals`（原型视觉效果）。**出厂的默认值应该是「关」**。
测试 `TestAC247HandingTheLevelToTheSeamIsNotVisibility` 一开头就检查这件事：
> 若这个开关在**本装配里是「开」**的，说明**有人把「交接电平」和「让像素可见」混为一谈了**——那是 AC#10 禁止的。

现在实测：**它是「开」的**。

## 现量（逐字，可复跑）

**红句**（`cmd/wisp/resident_audio_247_windows_test.go:355`）：
> `prototypeVisuals is ON in this assembly: the shipped default is off, and a leg that says it turned it on has crossed AC#10`

**代码现状**（`internal/ball/statevisual.go:110`，逐字）：
```go
var prototypeVisuals = true
```
⇒ **出厂值就是 `true`。** 而 `:114` 有 `func EnablePrototypeVisuals(on bool) { prototypeVisuals = on }` —— 运行时可改。

**耗时 0.00s**（读一个导出函数，不跑产品）⇒ **稳定，与负载无关。**

## ★关键线索：票 68 AC#2

测试注释 `:350-351` 逐字写着：
> `prettier AC#2 - applyLevel returns at internal/ball/liquid_windows.go:52 before any pixel can move. Flipping that default belongs to ticket 68 AC#2.`

⇒ **"翻转那个默认值"是票 68 AC#2 的事。** 所以有两种可能：
- **甲：票 68 AC#2 已经做了**（有意把它打开），那 **247 的断言就过时了** ⇒ 改断言（并写明"默认值已由票 68 翻转"）。
- **乙：有人误改了**（不是票 68 做的），那 **产码该改回 `false`** ⇒ 改 `:110`。

**⛔ 用"哪个容易"判。必须去查票 68 AC#2 到底做了没有。**

## 要做的（按序）

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC247HandingTheLevelToTheSeamIsNotVisibility$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **读票 68 的票面**（`.scratch/wisp/issues/68-*`），**逐字抄出 AC#2 原文**，看它是否要求把 `prototypeVisuals` 默认值翻转成 `true`。
  3. **查票 68 是否已 done**（文件名是否带 `-done`）。
  4. `git log -S 'var prototypeVisuals = true' --oneline -- internal/ball/statevisual.go`，找**哪个提交把它从 false 改成 true**（或它一直是 true），**贴出那个提交的说明**。
  5. 判甲／乙，给一句话结论＋`文件:行`＋AC#2 原文引用。

- [x] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲（票 68 已做）**：改 `cmd/wisp/resident_audio_247_windows_test.go:354-356` 的断言，让它反映"默认值已翻转"。⚠ **⛔ 简单删掉这个检查**——它守的是 AC#10（"交接电平"≠"让像素可见"），那个承诺仍然重要。要么改成检查别的证据，要么写明"本格随票 68 AC#2 关闭"。
  - **乙（误改）**：把 `internal/ball/statevisual.go:110` 改回 `false`。
  - **必须配突变控制真跑**：把修好的那处改回去 ⇒ 必须回到红。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 68／247 的关系写清**：本票只处理这一格，⛔ 不翻票 68／247 任何框；若判甲，**必须在证据件里写明"票 68 AC#2 已做，本格的承诺由它接管"**，⛔ 自行改票 68。

## 本票**不**解决

- 不碰 `TestTicket198*`（票 434 在飞）／`TestLockedSuffix*`（票 435 已立）／`Test313*`（另计）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\436-<腿码>\`。
- **⚠ 本票与真窗族测试（`Test258*`／`TestAC4Focus*`）⛔ 并行**。

## Progress log（append-only）

（每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-11-12 23:11:17 +08] agent=436-247visuals did=**判甲：票 68 AC#2 有意翻转，247 的断言过时**——写面落在**断言侧**（`cmd/wisp/resident_audio_247_windows_test.go`），`internal/ball/statevisual.go:110` **零改动**（两边只动一边）。**AC#0 只读取证**：①carrier（`D:\work\soft\Git\bin\bash.exe`）复现红句逐字 `:355: prototypeVisuals is ON in this assembly: the shipped default is off, and a leg that says it turned it on has crossed AC#10`，FAIL 0.00s；②票 68 AC#2 原文逐字（`68-ball-default-visuals-parity.md:57-65`）＝「`prototypeVisuals` 默认改为 `true`，使默认构建 == owner 签收的样子；`EnablePrototypeVisuals(false)` 与 `balldebug -frozen` 保留为逃生门……必须逐条把断言迁移到新默认下应有的行为……不得为了让测试变绿而删除断言或放宽阈值」⇒ **逐字要求翻 true**；③票 68 **文件名不带 `-done`**（票面 Status 仍 `review`、AC#2 框未勾，等 AC#3 真桌面复测与 owner 目验）**但 AC#2 已在码里落地**，两件事并存是票 68 自己的账；④`git log -S 'var prototypeVisuals = true' --oneline -- internal/ball/statevisual.go` ⇒ **唯一一枚 `4b9efc889`「票 68 落地腿 68-r1」**（CarlosShao，2026-11-01 22:50 +08），说明逐字节录「票 68 AC#2 落地（R15 三答案齐，registry A1378「全按推荐」）：statevisual.go: prototypeVisuals 默认翻 true——默认构建即 owner 签收的样子；EnablePrototypeVisuals(false)/balldebug -frozen 保留为冻结规格逃生门」，台账 `A1378`/`A1379` 旁证「R15 全按推荐＝票 68 解锁」⇒ **有主有据有 owner 批复，非误改**。**AC#1 只修一处＋AC#10 承诺续保**（⛔ 未删 if）：原断言把两件事绑在一个 flag 上——①「交接电平≠让像素可见」（AC#10 真承诺）②「出厂默认是关」（票 68 的样式决定）；票 68 拿走②，①改形留下，改成两条**与默认朝向无关**的读数：**:392-402** 交接一枚电平后 `PrototypeVisualsEnabled()` 必须**未变**（谁在交接路径上顺手开或关＝票 247 AC#10 逐字禁的那一手，当场红）、**:404-417** `EnablePrototypeVisuals(false)` 后同一枚标量**仍须到达目的地**（交接不许依赖开关），`t.Cleanup` 复原模式不泄漏；「像素动不动」那半格在注释 `:351-359` **明写随票 68 AC#2 关闭**并指库侧钉子 `internal/ball/tokens_test.go:158`，本测试不重复钉也不冒充。**突变控制真跑三发**（出仓 overlay，`go test -overlay`，变异体落 `D:\tmp\wisp\20261112\436-247visuals\`，经 carrier 的 `go test -c` 生效）：A 交接路径上顺手 `EnablePrototypeVisuals(true)` ⇒ **FAIL `:416` moved prototypeVisuals off -> on**；B 交接依赖开关 `if out != nil && ball.PrototypeVisualsEnabled()` ⇒ **FAIL `:413` destination saw 1 levels with prototypeVisuals off, want 2`**；C 把断言改回去（`HEAD` 版文件 overlay）⇒ **FAIL `:355` 逐字回到修前红句**（＝"改回去必红"那条尺）。**门禁四读数**：`go build ./cmd/wisp/` rc=0／`go vet ./cmd/wisp/` rc=0／`gofumpt -l` 空／`scripts/d22scan.sh` rc=0 `clean - no D22 ban violations`；整族 `-test.run '^TestAC247'` 全 PASS rc=0，修后本枚日志 `prototypeVisuals_at_assembly=true levels_reaching_the_ball_seam=2 levels_with_the_mode_forced_off=2`。**越界**：本腿改动只有 `cmd/wisp/resident_audio_247_windows_test.go`＋`docs/evidence/s1/436-247-prototype-visuals.md`（commit 显式 pathspec 只含这两枚）；`git diff --name-only` 里其余 12 行**全是进场前就在工作树的别人脏件**（`.scratch/wisp/issues/426-433` 的 D、`frontend/fixtures/composer-states.html`、`probes/291`、`346-a2-dispatch.txt`、以及**票 434 在飞的 `cmd/wisp/firstrun_198_test.go`**），本腿一枚没碰、一枚没进 commit；禁改面零命中、`build/mut_rvw.go` 未删。 evidence=`docs/evidence/s1/436-247-prototype-visuals.md` next=**编排者收件**（尺：提交真在／件真在／自己复跑那枚 carrier 命令）；票 68 的 AC#2 框与票 247 的框**本腿一枚未翻、未改票 68 一字**，⛔ 零 push。
- [2026-11-13 09:1x +08] **编排者：收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。腿交件 `d1f12dd55`（3 files／+222−11；证据件 `docs/evidence/s1/436-247-prototype-visuals.md` 10730B）。**三把尺**：①commit 在；②件 10730B 在盘；③**我自己复跑 PASS 0.03s**。**越界**：只改 `cmd/wisp/resident_audio_247_windows_test.go` 一枚；**`internal/ball/statevisual.go` 逐字未动**（`git show --stat` 证实）。
- **★判定＝甲（断言过时），只改断言。** 依据：**票 68 AC#2 逐字要求把这个默认值翻成 `true`**，且改它的提交 `4b9efc889`（2026-11-01 22:50，票 68 落地腿 68-r1）说明逐字 `票 68 AC#2 落地…statevisual.go: prototypeVisuals 默认翻 true——默认构建即 owner 签收的样子`；台账 `A1378`／`A1379`（`R15 全按推荐＝票 68 解锁`）为证。`internal/ball/tokens_test.go:158` 从库侧钉住新默认。
- **★票 68 AC#2 原文逐字（腿从 `68-ball-default-visuals-parity.md:57-65` 抄出）**：`AC#2 默认值翻转（需桌面复测）：prototypeVisuals 默认改为 true，使默认构建 == owner 签收的样子；EnablePrototypeVisuals(false) 与 balldebug -frozen 保留为"对照旧冻结规格"的逃生门。连带把 internal/ball 里断言旧默认的测试改到位：⚠ 不得为了让测试变绿而删除断言或放宽阈值——必须逐条把断言迁移到"新默认下应有的行为"…若某条测试的存在意义就是"钉住旧的冻结规格"，保留它并显式 EnablePrototypeVisuals(false)，别删。`
- **★票 68 状态的精确读法**：**票未 done**（文件名无 `-done`、Status 仍 `review`、AC#2 框未勾——等 AC#3 桌面复测与机主目视签收），**但 AC#2 本身 2026-11-01 已落码**。⇒ `票 68 未 done` 与 `AC#2 已做` **两者都真**，腿**没翻票 68 任何框**。
- **★承诺如何保住（⛔ 删断言）**：旧断言把**两件事**绑在一个标志上——AC#10 的承诺＋票 68 的样式决定。票 68 拿走了后者；腿把前者留在**不依赖默认值方向**的形状里（`:392-402`／`:404-417`）：①**交接不得移动模式**——任何方向中途翻转都红；②**交接不得依赖模式**——同一个标量在模式强制关闭时仍到达目的地。像素那半写进注释（`:351-359`）明说**由票 68 AC#2 关闭**、并点名库侧钉子，⛔ 假装还守着。
- **★突变控制三发（全经 carrier、仓外 overlay）**：A `EnablePrototypeVisuals(true)` 挂在交接路径 ⇒ **FAIL `:416`**（`the second handoff moved prototypeVisuals off -> on`）；B 汇聚点按模式门控 ⇒ **FAIL `:413`**（`the destination saw 1 levels with prototypeVisuals off, want 2`）；C 断言回退到 HEAD ⇒ **FAIL `:355`**（回到本票开始时的原句）。**三发都真红，证明新断言有牙。**
- **附带**：整个 `TestAC247` 族 PASS。

