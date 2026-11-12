# 435 — `TestLockedSuffixNeverTakesTheLockItself33r11` 红：四个新出现的 `…Locked` 函数未登记进名册（且各自体内零加锁操作）

- Status: **ready-for-agent**（写面＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go` 的 **`lockedNamingRoster33r11` 名册表**（`:86` 起）；⚠ 先只读取证，判清"该登记"还是"该改名"，再定改法）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11`）＋台账 `:32512` 存量真红
- 关联：票 33（`…Locked` 命名陷阱的原始记录）· 票 11（33r11 的 r11 修订）· 台账 `A1806`

## 这是什么（人话）

有一条命名规矩：**函数名带 `Locked` 后缀，是在告诉调用者"调用我的时候你已经持有那把锁了"**。
这个规矩的陷阱是：**如果函数自己又去拿同一把锁，就会死锁**（票 33 记的就是这个）。

于是有个检查：**每一个带 `Locked` 后缀的函数，都必须在名册里登记**，说明"它为什么不自己加锁"。
**新出现的名字必须被人审过、写进名册**——这样下一个加名字的人不会悄悄绕过这条规矩。

现在实测：**四个新名字没登记**。

## 现量（逐字，可复跑）

**红句**（`cmd/wisp/panel_locked_naming_33r11_windows_test.go:254`，出现 4 次）：
```
panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry: a new name that …
waveOutSpeaker.releaseLocked (voice_out_windows.go line 338) carries the "Locked" suffix but no roster entry: …
waveOutSpeaker.unprepareLocked (voice_out_windows.go line 376) carries the "Locked" suffix but no roster entry: …
```
**且日志同时打出**（`panel_locked_naming_33r11_windows_test.go:245`）：
```
panelHolder.stampPublishLocked in resident_windows.go line 1008: 0 lock operation(s) in its own body []
waveOutSpeaker.releaseLocked in voice_out_windows.go line 338: 0 lock operation(s) in its own body []
waveOutSpeaker.unprepareLocked in voice_out_windows.go line 376: 0 lock operation(s) in its own body []
```
⇒ **这四个函数体内确实一次加锁都没有**（`:257-259` 那条"体内有锁"的红**没有触发**）。
**所以它们⛔ 死锁陷阱，只是"没登记"。**

**耗时 0.03s**（纯 AST 静态检查）⇒ **稳定，与负载无关。**

## 机制链（逐字，⛔ 凭印象）

- 测试：`cmd/wisp/panel_locked_naming_33r11_windows_test.go`
  - `:86` `lockedNamingRoster33r11` —— **名册表定义处**（map，key 是 `类型.函数名`）
  - `:249-256` 遍历所有带 `Locked` 后缀的函数：
    ```go
    seen := make(map[string]bool, len(hits))
    for _, hit := range hits {
        seen[hit.qualified] = true
        why, rostered := lockedNamingRoster33r11[hit.qualified]
        if !rostered {
            t.Errorf("%s (%s line %d) carries the %q suffix but no roster entry: a new name that claims the caller holds a mutex must be reviewed and named in lockedNamingRoster33r11", …)
        }
        if len(hit.lockCalls) > 0 {
            t.Errorf("%s (%s line %d) is named as if the caller holds the lock while its own body takes it (%v) - that is the trap ticket 33 recorded: …")
        }
    }
    ```
- ⇒ **这是 AST 静态检查**（扫全包源码找 `Locked` 后缀 + 函数体内的锁调用），⛔ 运行时行为。
- **⚠ 仪器选择**：`go test -overlay` **对静态检查无效**（读磁盘源码）。**必须用仓外磁盘改写。**

## ⚠ 必须自己判清的核心岔路（⛔ 不许我替你定）

**"这四个名字该登记、还是该改名"——两种答案，修法不同：**

- **甲：该登记** ⇒ 它们**确实**是"调用者持锁"语义（名字对），只是名册漏了。⇒ 在 `:86` 的名册里加四行，**每行写清"为什么它自己不持锁"**（这是名册的存在意义）。
  **判据**：去读这四个函数的**调用点**，确认调用者是否真的持着那把锁。
- **乙：该改名** ⇒ 它们的名字在骗人（调用者并不持锁），后缀该去掉。⇒ **那是改产码**（重命名四个函数 + 所有调用点），⛔ 改名册。
  **判据**：若调用点并不持锁，那 `Locked` 后缀是**错的**，登记进名册等于**给一个错误的命名盖章**——那更糟。

**★这是本票的关键判断：名册是"审过的名单"，⛔ "消红工具"。**
若这四个名字确实该带 `Locked` ⇒ 登记（甲）；若不该 ⇒ 改名（乙），**⛔ 用登记把它们糊过去**。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴**四条红句全文**与耗时：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestLockedSuffixNeverTakesTheLockItself33r11$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. 读 `:86` 起的名册表，**数出现有登记几条**，并抄出 2–3 条现有条目的**写法格式**（你要照那个格式加）。
  3. **对四个函数逐个查调用点**：`git grep -n 'stampPublishLocked\|releaseLocked\|unprepareLocked'`，
     对每个调用点，**读它前后 10 行确认调用者是否真的持着那把锁**（找 `mu.Lock()`／`defer mu.Unlock()` 之类）。
  4. 判甲／乙，给一句话结论＋每个函数的依据（`文件:行`）。

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：在 `:86` 的名册里加四行，**每行写明"为什么它自己不持锁"**（⛔ 空字符串糊过）。
  - **乙**：改产码重命名四个函数 + 所有调用点（**这是大改，若判乙要先在证据件里说明影响面**）。
  - **必须配突变控制真跑**：甲 ⇒ 从名册里删掉一行必须红；乙 ⇒ 把改名还原必须红。
  - ⚠ **仪器**：这枚**读磁盘源码**，`-overlay` **无效**（台账 `A1804`）。**必须用仓外磁盘改写**，并说明你怎么确认仪器有效。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 33 的关系写清**：本票只处理这四个名字，⛔ 不翻票 33／11 任何框；⛔ 不宣称 `…Locked` 全族已审。

## 本票**不**解决

- 不碰 `TestTicket198*`（票 434 在飞）／`Test313*`／`TestAC247*`（另计）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\435-<腿码>\`。
- **⚠ 本票与真窗族测试（`Test258*`／`TestAC4Focus*`）⛔ 并行**。
