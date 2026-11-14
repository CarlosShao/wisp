# 票 448 — 把 446 的 errno 判定从「一次性证明」变成「常驻回归闸」：证据

**票**：`.scratch/wisp/issues/448-the-errno-fix-teeth-live-only-in-probes.md`
**裁决人（本文件作者）＝实现者**（写腿）。对抗验收／缺口审计须由**另一个** agent 做（`SPEC-12 §4.3` #1/#3），本文件不代替它。
**落档 commit**：`git log -1 --format=%H -- docs/evidence/s1/448-errno-teeth-resident.md`（写这份文件时的 HEAD 是 `2d2f6458b`；⛔ 本腿不 push，所以本文件不写死一个尚未被推送的串）
**改动文件**：新增两个测试件（唯一产码改动＝测试面，`internal/tools/system_backend_windows.go` **零字节改动**，sha256 见 §4）
**时刻**：2026-11-14 18:5x–19:1x +08

一句话：**446 的两颗牙原来只在 `probes/` 里；现在它们住在包里了** —— 判定退回英文句子比较会红，调用点绕过 `errnoPresent` 也会红，两条都由 `go test ./internal/tools` 与两个 CI job 自动跑。

---

## 1. 选点：`internal/tools/*_test.go` 而不是 `tools/d22scan` 的静态闸（票要求二选一并写理由）

**决定：落在 `internal/tools/*_test.go`（两个新件）。** 四条理由，按分量排：

1. **这条不变量是包内的。** `errnoPresent` 是 `package tools` 的未导出辅助函数（`system_backend_windows.go:270` 定义），出这个包没有任何调用者。`tools/d22scan` 管的是**全仓横切**的禁止形状（裸 goroutine、PathResolver 旁路、明文密钥…）；为一个包的某个辅助函数加一条全仓 ban，是拿大炮打一只包内蚊子。
2. **往 d22scan 加 ban 是动契约，不是动测试。** 那个工具的 ban 清单是 `PLAN.md` D22 冻结清单的**转写**：`main.go` 头注释逐条列出 1..10，且 `ban #9` 的注释明写「加一棵树都属业主批准的 scope 变更」（`tools/d22scan/main.go:60-71` 附近）。新增 `ban #11` 会把仪器放到契约前面，把一次**人工批准**拖进一张「小活」票里，还会在工具里永久留下一条「对象可以消失」的规则。按 `AGENTS.md` ⓪⓪⓪ 第 4 条，改契约文字是要摆到业主面前的三类之一 —— 本票没有理由背这个包袱。
3. **覆盖面：测试件在两个 CI job 里跑，d22scan 只在一个。** `internal/tools` 在 `scripts/portable-tests.sh` 的 core scope 名单里（`portable-tests.sh:385` 那一行就是它的 import path），所以无 build tag 的那一半在 **test-core（ubuntu）** 与 **test-windows** 都会跑；d22scan 只在 **lint** job。守护这个文件的闸，应当在这两个**会编译这个文件**的 job 里读数。
4. **牙的成本。** d22scan 新 ban 要配 `selftestsamples.go` 的种子样本＋`scan_test.go` 的钉子（这套机制它已经有了，为一条规则再走一遍不便宜）；测试件这边，红直接来自票面点名的 M1/M2 两次变异，加上文件内自带的三段种子字面量（§2 末）。

**反过来，什么命题才该进 d22scan**：若是「全仓任何文件都不许比较 `.Error()` 字符串」，那是横切形状，该进那个工具。本票没有主张那个命题（446 的普查也量到全仓只有 2 处，已经清零）。这个边界写在这里，免得下一条腿把它读成「448 认为 d22scan 不该管这类事」。

---

## 2. 闸门做了什么

两个新件，职责不重叠：

| 文件 | build tag | 内容 | 在哪跑 |
|---|---|---|---|
| `internal/tools/ticket448_errno_gate_test.go` | 无（跨 GOOS） | 读源码形状的三条闸 ＋ 种子反证 | test-core（ubuntu）＋ test-windows ＋ 本机 |
| `internal/tools/ticket448_errno_gate_windows_test.go` | `//go:build windows` | 行为闸：`errnoPresent` 的 10 行判据表 | test-windows ＋ 本机 |

静态那一半之所以无 tag，是因为它读的是**源码**（`system_backend_windows.go` 就在 checkout 里），不需要编译 Windows API —— 于是「编译器都碰不到那行」的 ubuntu 上，这条闸照样能红。

### 闸 1 `Test448ErrnoPresentDecidesByCodeNotText`（`ticket448_errno_gate_test.go:398`）

对 `errnoPresent` 的函数体要求三件事，任一不满足即红并点名行号：

- 体内**不许有字符串字面量**（句子比较必然带一个字符串：`==`、`strings.HasPrefix`、`EqualFold`、拼过再比，全是这一条的子集）；
- 体内**不许调 `.Error()`**；
- 体内**必须有** `errors.Is(<...>, syscall.Errno(...))` 形状的码判据（446 修好的那半不许被删）。

### 闸 2 `Test448ErrnoCallSitesGoThroughErrnoPresent`（`:406`）

规则按「errno 有没有被**拿来作决定**」判，不是按「函数里出现没出现 `errnoPresent`」判：

- 一处调用点把 `Proc.Call` 的第三返回值接进命名变量（`r, _, e := procX.Call(...)`，含 `if r, _, e := ...; cond` 的 init 形式），并在任何**决定位置**用到它 —— 决定位置＝任何比较、`if`/`for` 的条件、`switch` 的 tag —— 则那个决定必须由 `errnoPresent(该变量)` 来做，否则红并点名「函数名＋变量名＋行号」。
- 收「所有比较」而不只收 `if` 条件，是为了堵住「先算成布尔再判」那类搬移：`has := e != nil && ...` 里的 `e != nil` 本身就是一处比较。
- **只把 errno 印进消息、从不判它存不存在的调用点放行。** 仓里现有一处正这样：`system_backend_windows.go:680` 的 `LookupPrivilegeValueW`，它拿 `r == 0` 判失败、把 errno 交给 `%v` 打印 —— 那里没有「errno 到底有没有」这个问题可绕。硬要求它走 `errnoPresent` 是把闸门开成误报。**这一条是闸门第一次跑就点出来的第三处绑定**（本票开工前没人点过 `:680`），写了理由才放行，不是漏掉。

### 闸 3 `Test448NoErrorTextIsDecidedOnInTheBackend`（`:421`）

兜底：整个文件里，`.Error()` 的结果不许落进「作决定」的位置 —— 不作为比较的操作数、不交给 `strings.*`、不当 `switch` 的 tag。沿父链判，遇到语句边界就停（那里的文本只是被打印，不是被判）。这条抓的是闸 1/闸 2 之外的写法：`e.Error()[:10] == "The operat"`、`switch e.Error() { case ... }` 这种。

### 种子反证 `Test448TheStaticChecksGoRedOnTheOldShapes`（`:467`）

三段**源码字面量**（不是盘上的文件）：446 修复前的判定体、修复前的调用点、修好后的样本。同一批检查函数在**前两段上必须报红**、在**第三段上必须静默**。这枚反证的作用是：如果哪天这些检查器不再匹配任何东西（正则/AST 形状漂了），它会红 —— 于是「一处都没报」不会读成绿。它也是审计点名的那类洞（"matched nothing = green"）的直接解。

### 行为闸 `Test448ErrnoPresenceIsDecidedByCode`（`ticket448_errno_gate_windows_test.go:55`）

`errnoPresent` 的 10 行判据表（nil、errno 0 裸值、`windows.Errno(0)` 别名、errno 0 带本地化文本、包一层、包两层、errno 5、errno 87、errno 87 包一层、长得像成功句子的真错误）。逐行 `Errorf` 而不是在第一条就 `Fatalf`：一次红要把**所有**读文本的行都点出来，否则证据里只剩第一条，读者看不出牙有几颗。

### 范围（明确不覆盖，免得被读大）

- 三条静态闸的范围＝**`internal/tools/system_backend_windows.go` 这一个文件**。包里别的文件里若新长出一处 errno 判定，本闸看不见 —— 本票写面就到这里，不扩。
- 把 errno 交给**另一个函数**或存进字段、事后才判，AST 级闸门跟不到（需要数据流）。这类绕行仍会被闸 3 抓住「只要它比较了文本」，而 446 要防的正是文本比较。

---

## 3. 两个方向各红过一次（AC#2）—— 含「行为闸为什么看不出来」

三次变异，都在写腿现场跑，日志在 `.scratch/wisp/probes/448/`，每次跑完都从 pristine 还原并核对 sha256（§4）。

### M1 判定函数本体退回英文句子比较

变异＝把 `errnoPresent` 体改回 `e != nil && e.Error() != "The operation completed successfully."`，并按 `git revert` 那半个 hunk 的**真实形状**同时删掉随之空转的 `errors`／`syscall` 两条 import（只改函数体不删 import 会先在编译期红在 "imported and not used"，那读起来是编译器在报，不是闸门在报 —— 这种形状也记在 `10-M1-...log` 的第一次读数里）。

日志 `10-M1-judgement-reverted-to-sentence-RED.log`，`go test ./internal/tools -run Test448` **rc=1**：

- 闸 1 点名 3 处：`system_backend_windows.go:269` 两句（`.Error()`、字符串字面量）＋ `:268`（码判据没了）；
- 闸 3 点名 1 处：`:269` 错误文本被拿来做比较；
- 行为闸 **4/10 行红**：`errno-zero-localized-text`、`errno-zero-wrapped`、`errno-zero-wrapped-twice`、`lookalike-text-is-not-errno-zero`。**这 4 行正是与语言无关的 4 颗牙**（§5）。

**一处必须被读到的读数**：同一个变异下，`errno-zero-raw` 与 `errno-zero-via-windows-alias` 这两行在**本机是绿的** —— 因为本机拿得到 en-US 消息表，旧比较那句英文恰好对得上（446 已追到 `syscall_windows.go:169-171`）。⇒ **拿这两行当牙的闸门是假闸**；本闸门的牙是那 4 行。

### M2 两处调用点退回内联句子比较（判定函数本身不动）

变异＝`realEnumWindows` 与 `sendInputs` 两处 `errnoPresent(e)` 改回 `e != nil && e.Error() != "..."`。日志 `20-M2-callsites-bypass-errnoPresent-RED.log`，rc=1：

- 闸 2 点名**两处**：`system_backend_windows.go:297`（`realEnumWindows`）、`:439`（`sendInputs`）；
- 闸 3 红（这两行都拿文本作了比较）；
- 闸 1 **绿**（判定函数还是对的）；
- **行为闸 10/10 全绿。**

**这一格是整张票的理由**：M2 对行为测试是**不可见的** —— 行为闸一直在测一个正确但没人用的 `errnoPresent`。所以方向②只能靠静态一条抓；票面 AC#2 那句「静态闸可，但要在证据里说明它为什么能抓到」，答案就是上面这段读数：它看的不是「判定函数对不对」，而是「调用点把 errno 交给谁作决定」。

### M2b 调用点内联 `errors.Is` 绕过（无文本参与）

变异＝只把 `sendInputs` 那处改成 `if e != nil && !errors.Is(e, syscall.Errno(0))`。日志 `21-M2b-callsite-inline-errorsIs-bypass-RED.log`，rc=1：

- 闸 2 红，点名 `:439`；
- 闸 1 绿、**闸 3 绿**（一个字面量都没出现，兜底那条确实看不见）、行为闸 10/10 绿。

⇒ 闸 2 抓的是**结构**（谁作决定），不是「有没有出现那句英文」；加这一次变异就是为了让这个区别有读数，免得下一条腿把闸 2 读成一句英文的黑名单。

---

## 4. 446 的判定逻辑零字节改动（还原判据）

- 开工前 `internal/tools/system_backend_windows.go` sha256＝`8CAD968D67AE4F3D4795CDBDB136A6A28E48AB93B44A76FEF716D61026738845`（与 `docs/evidence/s1/446-errno-presence.md` 记的**同一个串**），`git diff` 空。
- 变异前先把原文件 **Copy-Item** 到 `.scratch/wisp/probes/448/system_backend_windows.go.pristine`（不用 `git checkout`／`stash`，避开共享工作树里吞别人活的禁手）。
- 每次变异后 `Copy-Item` 还原，逐次核对：M1 后 `30-post-restore-GREEN.log`、M2/M2b 后 `31-post-restore-GREEN-M2.log`，两处 sha256 都回到 `8CAD96…845`，`git diff --name-only` 空、rc=0，`go test -run Test448` rc=0。
- 本票落档的 diff 里**没有** `system_backend_windows.go`（`git show --stat` 可核）。

---

## 5. AC#3 环境无关性（任何语言的 Windows 上都成立）

闸门**不读** `GetUserDefaultUILanguage`／`GetSystemDefaultLangID`，也不用任何真机消息表：输入全是 `syscall.Errno(0)` 的字面值、`windows.Errno(0)` 别名、以及两个**常量样例**（`"操作成功完成。"` 与 `"The operation completed successfully."`）。

「与语言无关的牙」不是散文，是结构字段 `localeIndependentTooth`，门里断言**至少 3 处**（现为 4 处），所以后人裁表时不能把这几行悄悄删掉。本机读数（`60-package-tests-full.log` 里行为闸逐行）：

| 行 | 旧写法（446 前）在本机 | 旧写法在任何语言的 Windows 上 | 牙 |
|---|---|---|---|
| `nil` | 对 | 对 | 否 |
| `errno-zero-raw` | **对（碰巧）** | 错 | 否 |
| `errno-zero-via-windows-alias` | **对（碰巧）** | 错 | 否 |
| `errno-zero-localized-text` | 错 | 错 | **是** |
| `errno-zero-wrapped` | 错 | 错 | **是** |
| `errno-zero-wrapped-twice` | 错 | 错 | **是** |
| `errno-5` / `errno-87` / `errno-87-wrapped` | 对 | 对 | 否（防「真错误被吞」的控制行） |
| `lookalike-text-is-not-errno-zero` | 错 | 错 | **是** |

M1 的实测红点＝这 4 行，与表一致。

---

## 6. 门禁（AC#4，全部现跑；日志 `.scratch/wisp/probes/448/`）

| 闸 | 读数 | 日志 |
|---|---|---|
| `go test ./internal/tools -count=1` | **rc=0，PASS=431 FAIL=0 SKIP=1**；包行 `ok … 34.556s` | `60-package-tests.log` ＋ `60-package-tests-full.log` |
| `go vet ./internal/tools` | **rc=0** | `60-package-tests.log` |
| `GOOS=linux go vet ./internal/tools` | **rc=0**（无 build tag 的那半个在 ubuntu 上也编译，跨 OS 的说法有读数） | `61-linux-compile.log` |
| `gofumpt -l`（两个新件） | **空输出，rc=0** | 见 §7 |
| `d22scan`（`cd tools/d22scan && go run . -root ../../`） | **rc=0，clean**；`ban #8 internal/ examined 1016 Go files`（含新件） | `50-d22scan.log` |
| path-length budget | **VERDICT GREEN rc=0**；tracked=20177，over-budget=69 全在 roster | `51-path-length.log` |

**那一枚 SKIP 不是本票的**：`Test23LiveClipboardRoundTripThroughRealOS` 要 `WISP_T23_LIVE_CLIPBOARD=1`（真剪贴板往返是操作者点名才跑的腿），本票开工前就在跳。本票新增的 7 个测试条目（4 个顶层 + 3 个子测试）**0 SKIP**。

---

## 7. AC#5 留白（CI）

⛔ 本腿不 push（票面与全局铁律）。`test-core`／`test-windows` 两发含本 commit 且**完成**的运行，要等编排者推送后取：需要 run 号＋`merge-base`＋完成时刻三项。**本文件不声称它们已经跑过。**

---

## 撤销口令

「票 448 撤」＝删本票＋回退其写腿 commit（新 commit 反改，⛔ reset／⛔ amend），并删除两个新测试件与 `docs/evidence/s1/448-errno-teeth-resident.md`；`internal/tools/system_backend_windows.go` 不需要动（本票没改它一个字节）。
