# 436 — `TestAC247HandingTheLevelToTheSeamIsNotVisibility` 红：`prototypeVisuals` 出厂值是 `true`，而断言要求 `false`——票 68 把它打开了还是有人误改？

- Status: **ready-for-agent**（写面＝`internal/ball/statevisual.go:110` **那一行**，或 `cmd/wisp/resident_audio_247_windows_test.go:354-356` 的断言——⚠ **先只读取证，判清是"票 68 有意打开"还是"误改"，再定改法**）
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

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴逐字红句：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC247HandingTheLevelToTheSeamIsNotVisibility$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **读票 68 的票面**（`.scratch/wisp/issues/68-*`），**逐字抄出 AC#2 原文**，看它是否要求把 `prototypeVisuals` 默认值翻转成 `true`。
  3. **查票 68 是否已 done**（文件名是否带 `-done`）。
  4. `git log -S 'var prototypeVisuals = true' --oneline -- internal/ball/statevisual.go`，找**哪个提交把它从 false 改成 true**（或它一直是 true），**贴出那个提交的说明**。
  5. 判甲／乙，给一句话结论＋`文件:行`＋AC#2 原文引用。

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲（票 68 已做）**：改 `cmd/wisp/resident_audio_247_windows_test.go:354-356` 的断言，让它反映"默认值已翻转"。⚠ **⛔ 简单删掉这个检查**——它守的是 AC#10（"交接电平"≠"让像素可见"），那个承诺仍然重要。要么改成检查别的证据，要么写明"本格随票 68 AC#2 关闭"。
  - **乙（误改）**：把 `internal/ball/statevisual.go:110` 改回 `false`。
  - **必须配突变控制真跑**：把修好的那处改回去 ⇒ 必须回到红。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 68／247 的关系写清**：本票只处理这一格，⛔ 不翻票 68／247 任何框；若判甲，**必须在证据件里写明"票 68 AC#2 已做，本格的承诺由它接管"**，⛔ 自行改票 68。

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
