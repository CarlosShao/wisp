# 435 — `TestLockedSuffixNeverTakesTheLockItself33r11` 红：**只有一枚**新名字 `panelHolder.stampPublishLocked` 未登记进名册

- Status: **done**（写面＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go` 的 **`lockedNamingRoster33r11` 名册表**（`:86` 起）；⚠ 先只读取证，判清"该登记"还是"该改名"，再定改法）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`，`PASS=583 / FAIL=11`）＋台账 `:32512` 存量真红
- 关联：票 33（`…Locked` 命名陷阱的原始记录）· 票 11（33r11 的 r11 修订）· 台账 `A1806`

## 这是什么（人话）

有一条命名规矩：**函数名带 `Locked` 后缀，是在告诉调用者"调用我的时候你已经持有那把锁了"**。
这个规矩的陷阱是：**如果函数自己又去拿同一把锁，就会死锁**（票 33 记的就是这个）。

于是有个检查：**每一个带 `Locked` 后缀的函数，都必须在名册里登记**，说明"它为什么不自己加锁"。
**新出现的名字必须被人审过、写进名册**——这样下一个加名字的人不会悄悄绕过这条规矩。

现在实测：**只有一枚没登记**——`panelHolder.stampPublishLocked`。

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

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 用 carrier 复现，贴**四条红句全文**与耗时：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestLockedSuffixNeverTakesTheLockItself33r11$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. 读 `:86` 起的名册表，**数出现有登记几条**，并抄出 2–3 条现有条目的**写法格式**（你要照那个格式加）。
  3. **对四个函数逐个查调用点**：`git grep -n 'stampPublishLocked\|releaseLocked\|unprepareLocked'`，
     对每个调用点，**读它前后 10 行确认调用者是否真的持着那把锁**（找 `mu.Lock()`／`defer mu.Unlock()` 之类）。
  4. 判甲／乙，给一句话结论＋每个函数的依据（`文件:行`）。

- [x] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲**：在 `:86` 的名册里加四行，**每行写明"为什么它自己不持锁"**（⛔ 空字符串糊过）。
  - **乙**：改产码重命名四个函数 + 所有调用点（**这是大改，若判乙要先在证据件里说明影响面**）。
  - **必须配突变控制真跑**：甲 ⇒ 从名册里删掉一行必须红；乙 ⇒ 把改名还原必须红。
  - ⚠ **仪器**：这枚**读磁盘源码**，`-overlay` **无效**（台账 `A1804`）。**必须用仓外磁盘改写**，并说明你怎么确认仪器有效。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 33 的关系写清**：本票只处理这四个名字，⛔ 不翻票 33／11 任何框；⛔ 不宣称 `…Locked` 全族已审。

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
---

## ⚠ 本票的更正记录（编排者 2026-11-12 23:0x，⛔ 删，留档）

**我初版票面写「四个新名字没登记」，这是错的。** 实测（`D:\tmp\wisp\20261112\435-probe\run.log`）：

- `:245` 那 4 行是 **`t.Logf` 的普查日志**（列出所有带 `Locked` 后缀的函数及其体内锁调用数），**⛔ 红句**。
- **真正的红句只有 1 条**（`:254`）：`panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry`。
- `waveOutSpeaker.releaseLocked`／`unprepareLocked` **在名册里已有条目**（`:88`、`:89`），`PanelManager.setPriorFocusLocked` 也有（`:87`）⇒ **这三枚是绿的**。

**教训（第 9 次同根）**：我把「普查日志的 4 行」当成了「4 条红句」。**这是"读日志只看形状不看级别"的错**——`t.Logf` 与 `t.Errorf` 在日志里长得很像，我该看 `:245` 与 `:254` 的行号差异和 `--- FAIL` 只出现一次这件事。

⇒ **本票的范围缩为一枚：`panelHolder.stampPublishLocked`。** 判定逻辑不变（该登记还是该改名）。


---

## Progress log（只追加，⛔ 翻 AC 框）

- **2026-11-13 09:2x 写腿 `435-33r11`：判甲，名册加一枚（AC 框一枚未翻、未 push）**
  - **AC#0 只读**：carrier（`D:\work\soft\Git\bin\bash.exe`）复现，`--- FAIL … (0.07s)`、`rc=1`、端到端 11.29s。**红句只有 1 条**（`:254`）：`panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry`；`:245` 那 4 行确认为 `t.Logf` 普查日志。**票面末尾"更正记录"复核成立**。
  - **名册修前 3 条**，格式＝`"<类型>.<方法名>": "requires the caller to hold <锁>; body takes no lock"`（`PanelManager.setPriorFocusLocked`／`waveOutSpeaker.releaseLocked`／`waveOutSpeaker.unprepareLocked`）。
  - **调用点逐个读前后 10 行（3 处产码调用，全部持锁）**：`resident_windows.go:985`（`start`，锁在 `:979-980`）、`:1044`（`toggle`，锁在 `:1033-1034`）、`:1104`（`show`，锁在 `:1093-1094`）——三个调用点都在 `h.mu.Lock()` / `defer h.mu.Unlock()` 临界区内，且都紧跟 `h.rp = h.spawn(...)`。函数 doc `:1005-1007` 逐字：`Called with h.mu held (every spawn site holds it)`。函数体 0 次加锁；它唯一碰的 `rp.mu`（`panel_resident_windows.go:555-557`）是**面板线程自己的另一把锁**，非 `h.mu`。另 1 处命中 `panel_resident_pump_windows_test.go:240` 是票 65 的源钉，⛔ 调用点。
  - **判甲（该登记）**：三个调用点全持 `h.mu`、函数体不碰 `h.mu` ⇒ `Locked` 后缀说的是真话，名册漏登，⛔ 名字骗人。判乙会把这个"spawn 窗口必须持锁"的语义从名字上抹掉。
  - **修法（只修一处）**：`cmd/wisp/panel_locked_naming_33r11_windows_test.go` 名册加 1 行（`:101`）＋ 文档块加 1 段（`:83-93`）写明"为什么它自己不持锁"。**产码零改动**（`resident_windows.go` 一字节未动）⇒ 影响面 0 个调用点。
  - **修后复跑**：`--- PASS … (0.04s)`、`rc=0`、端到端 18.49s；普查仍 `4 Locked-suffixed func(s) across 72 production source file(s)`。
  - **突变控制（★仓外磁盘改写，⛔ overlay）**：仓外模块树 `D:\tmp\wisp\20261113\435-33r11\module-tree\`（24.8 MB／1484 文件，`go test -c ./cmd/wisp/` rc=0，cwd＝`module-tree\cmd\wisp`）。**M0** 基线 PASS／4 命中 72 文件。**M1**（不重编译，二进制 sha256 前后同为 `4B594E99…A349E`）仓外 `cmd/wisp/` 新增一个 `mut435ProbeLocked` ⇒ **FAIL**，普查 `5 Locked-suffixed func(s) across 73 production source file(s)` 并点名 `zz_mut435_probe_windows.go line 7`；删掉该文件同一二进制回 **PASS**。**M2**（不重编译）给 `stampPublishLocked` 体首插 `h.mu.Lock()/defer Unlock` ⇒ **FAIL**，普查该行变 `2 lock operation(s) [h.mu.Lock h.mu.Unlock]`，红句 `:270 …is named as if the caller holds the lock while its own body takes it…`；还原回 PASS。**M3**（名册编译进二进制，故须重编译）仓外删掉本票新增那行 ⇒ **FAIL**，红句回到 `:265 … carries the "Locked" suffix but no roster entry`；还原重编译回 PASS。
  - **仪器有效性怎么确认（三问）**：①**同一二进制两种判定**（M1／M2 二进制 sha256 未变，唯一变量是盘上源码；若读烘焙快照则数字不可能跟着磁盘走）；②**`-overlay` 真的无效**——对照 A 故意写成不可解析的 overlay 源 ⇒ 编译**失败** `syntax error: unexpected {, expected )`（证明 overlay 通路是通的）；对照 B 换成合法 overlay 源 ⇒ 编译**成功**，但同一二进制跑出来普查仍 `4 … across 72` 且 **PASS**（编译器看见第 5 个函数、普查看不见）⇒ 台账 `A1804` 复核成立；③**名册⛔免死金牌**——M2 证明已登记的条目体内一加锁照样红。
  - **越界**：本腿写面＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go` 一枚（＋本票 log ＋证据件）。`internal/panel/**` **零触碰**（437 腿地界）。禁止面（`PLAN.md`／`specs`／`thresholds.go`／golden／`allowlist.txt`／`frontend`／`design`）零命中。`git diff --name-only` 全量输出里另 13 枚是起手即脏的**别人在飞腿**改动（票 426–436 票面、291 探针、346 派单件、前端 fixture），本腿未碰、commit 用显式 pathspec 只提三枚。
  - **门禁四读数**：`go build ./cmd/wisp/` `rc=0`／`go vet ./cmd/wisp/` `rc=0`／`gofumpt -l` **零输出**／`scripts/d22scan.sh` **clean**（`examined 480 production Go files`，正控 `PASS=45 FAIL=0 SKIP=0`，`ok … 53.608s`）。`build/mut_rvw.go` 未删未碰；门禁用 `go build ./cmd/wisp/`，⛔ `go build ./...`。
  - **AC#4 关系**：只处理这一枚名字，**⛔ 不翻**票 33／11 任何框，**⛔ 宣称** `…Locked` 全族已审。
  - **证据件**：`docs/evidence/s1/435-33r11-roster-stamppublish.md`。**撤销口令：「435 回乙」**（把名册那 1 行与文档那段删掉即复原，产码从未改动）。
- [2026-11-13 09:3x +08] **编排者：收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。腿交件 `7cf184068`（3 files／+218−0；证据件 `docs/evidence/s1/435-33r11-roster-stamppublish.md` 12576B）。**三把尺**：①commit 在；②件 12576B 在盘；③**我自己复跑 PASS 0.06s**，且普查仍报 `4 Locked-suffixed func(s) across 72 production source file(s)` ⇒ **⛔ 靠缩小扫描面消红**。**越界**：只改名册测试一枚，**产码零改动**（`--stat` 逐字证实）。
- **★判定＝甲（该登记），只改名册一处。** 理由：**三个产码调用点无一例外都在 `h.mu` 临界区内**（`:985` 的 `start` 在 `:979 Lock`／`:980 defer Unlock` 内；`:1044` 的 `toggle` 在 `:1033`／`:1034` 内；`:1104` 的 `show` 在 `:1093`／`:1094` 内），且都紧跟 `h.rp = h.spawn(...)`；函数 doc `:1005-1007` 逐字 `Called with h.mu held (every spawn site holds it)`。⇒ **`Locked` 后缀说的是真话，名册是漏登**；判乙改名反而会把"spawn 窗口必须持锁"这层语义从名字上抹掉。
- **★登记条目写得很准确（值得记）**：`"panelHolder.stampPublishLocked": "requires the caller to hold h.mu; body takes no lock on h.mu (rp.setSnapshotPublish takes the panel thread's own rp.mu)"` —— **连"体内碰到的那把 `rp.mu` 是面板线程自己的另一把锁"都注明了**（`panel_resident_windows.go:555-557`），⛔ 糊一句"takes no lock"。
- **★仪器有效性三问（腿做得最硬的一次，值得记）**：仓外模块树（24.8MB／1484 文件）**⛔ 全程未用 overlay**。①**同一二进制两种判定**——M1（仓外新增一个 `…Locked` 探针函数）与 M2（给函数体首插 `h.mu.Lock()/Unlock`）**都不重编译**（sha256 前后同为 `4B594E99…A349E`），读数却分别变成 `FAIL 5 … across 73`（并点名 `zz_mut435_probe_windows.go line 7`）与 `FAIL 2 lock operation(s)` ⇒ **唯一变量是盘上源码，真读盘**。②**`-overlay` 确实无效（正反两向对照）**——对照 A：故意不可解析的 overlay 源让编译**失败** `syntax error`（证明通路是活的）；对照 B：合法 overlay 源编译**成功**，但同一二进制跑出来普查仍 `4 … across 72` 且 **PASS**（**编译器看见第 5 个函数、普查看不见**）⇒ 台账 `A1804` **复核成立**。③**名册不是免死金牌**——M2 证明已登记条目体内一加锁照样红。
- **★M3（删名册那行）必须重编译**（名册编译进二进制）⇒ **FAIL**，红句回到 `:265 … no roster entry`。**四发突变全真红。**
- **门禁**：build rc=0 · vet rc=0 · gofumpt 零输出 · `scripts/d22scan.sh` clean（480 文件，`PASS=45 FAIL=0 SKIP=0`，53.608s）。
- **撤销口令：「435 回乙」**——删掉名册那 1 行与文档那段即复原，**产码从未改动**。

