# 票 347 `AC#1` 三形代价表 — 腿 `347-a1`（**只读普查腿**，⛔ 写腿）

- 锚＝`HEAD 385afead9`（2026-10-26 19:19:02 +08）。`git diff --cached --name-status`＝**空**（无前置阻塞）；`git status --porcelain -- internal/risk/`＝**空** ⇒ 本文全部读数取**已提交树**，无脏树混读。
- 本腿射程＝**只出表**。⛔ 未落任何 `.go`、⛔ 未跑任何 `go test`、⛔ 未翻框／改 `Status:`／加 `-done`／改台账、⛔ push。
- 票面 `AC#1` 逐字「**未出表不许动码**」⇒ **本文是落码的前置**，不是落码。下面凡"红/绿"数字，**没有一枚是本腿跑出来的**（分界见 §7）。

## 1. 六条必现量（逐条：读数 ＋ 行号 ＋ 我这把尺）

| # | 必现量 | 读数（逐行） | 我这把尺（痕迹） |
|---|---|---|---|
| 1 | **桩本体** | `internal/risk/pathresolver_other.go:1` `//go:build !windows`；`:5-9` `DEFERRED(macOS/Linux)` 注释块；**`:10`** `func resolveHandle(p string) (string, bool) { return "", false }`；**`:12`** `func reparseComponents(p string) []string { return nil }`（★**两枚**桩，票面只点了第一枚） | `read` 全文 12 行 |
| 2 | **唯一赋值点** | `pathresolver.go:357` `f.root = normPath(final)` —— 落在 `:354` `if final, ok := resolveHandle(cut); ok` → `:356` `if final != ""` → `:365` `return f` 的**成功分支内**；`formsOf` 本体起 `:345` | `git grep -n '\.root = '` 全仓 **5 命中**，`internal/risk` **包内非测试仅此 1 枚**（另 4 枚＝`cmd/wisp`／`internal/observe`／`internal/panel`／`syncdirs_test.go` 的同名异型）＋现读 `:300-409` |
| 3 | **结构不可达的证明** | `blacklist.go:`**`363`** `if cand.root == "" \|\| anchor.root == "" { return false }`（函数签名 `:362`，`:355-361` 是注释；`isUnder` 本体 `:170-171`＝`path == dir \|\| strings.HasPrefix(path, dir+sepStr)` ⇒ `isUnder("","")`**恒真**） | 现读 `:330-389`。★票面写 `:359-361`／`339_other_test.go:20` 也写 `:359` —— **偏 4 行**，`:359` 是注释行 |
| 4 | **五红分母可信度** | 判：**可信**（逐条见 §3） | 读 `339/v1/logs/` 四份落地日志＋逐枚守卫行 |
| 5 | **`portable-tests.sh` 里 `internal/risk` 的档位** | 在 **`core`** 档射程内（scope `:347` `./internal/risk/...`）＋**`windows`** 档（`:359`）；逐名点名在 `core_pin` **`:292`**、`win_pin` **`:309`**；`:7-12` 逐字警告"Go 把 SKIP 记成 `ok`" | 现读脚本 `:1-20`／`:261-315`／`:339-364`。★票面写 `:266`／`:282` **两处都错**（`:266` 是注释首行、`:282` 是 `internal/llm/openaichat`） |
| 6 | **票 55 的 `Blocked by` 与 `AC :28`** | `55-s8-macos-port.md:6-7` `Blocked by: 54-s7-acceptance + self-use week + explicit user go-ahead (S8 starts only pre-open-source, D23)`；`AC :28` 逐字 `PathResolver macOS branch (realpath + lstat) passes red-team port.` | `read` 55 头 45 行 |

## 2. 三形 × 8 列（票面 §前置 的 一／二／三）

| 形 | 产码面（逐文件） | 产码量级 | ★契约影响 | ★分母变化（Linux ⇄ Windows，`./internal/risk`） | 下游解锁 | ★风险（逐条；未验即标） | 可逆性 |
|---|---|---|---|---|---|---|---|
| **一**<br>＝AC#1 **甲**<br>落解析器 | ①`pathresolver_other.go`（改：删桩＋删 `DEFERRED` 块；落 `resolveHandle`＝`filepath.EvalSymlinks`＋存在性守卫；**落 `reparseComponents`**＝按 `sepStr` 逐组件 `os.Lstat` 判 `ModeSymlink`、首枚缺失即停）②`pathresolver.go`（改：`normalizeLocalUNC` 加 `filepath.Separator=='\\'` 编译期早退——见风险②；`formsOf` 与 `:117` 调用序**不动**）③`anchor_fallback_teeth_336_other_test.go`（**重写**：tripwire→真牙）④`anchor_fallback_teeth_336_windows_test.go`（改：8 枚牙 roster/夹具提成共享件或去 tag）⑤`anchor_miss_reason_339_other_test.go`（**重写**：两枚 tripwire→reason/class 两枚真牙，**保** `candidateOffAnchorChain334` 名）⑥`syncdirs_other_test.go`（改：删 `TestSyncNoGradeIsConfirmedOnPosix`——其文件头 `:37-38` 自己就写"delete"；拆 `TestSyncMembershipDecidesOffProfilePosix:94-96` 的前提断言）⑦`scripts/portable-tests.sh`（改：`:707` 那行 account-skip 会陈旧——见风险①；`:710` 那行**不**受影响） | 产码 ≈**40-55 行**／1 件 ＋ `pathresolver.go` ≈**3-6 行**；测试 ≈**250-400 行**／3 改＋1 重写；脚本 ≈2 行 ⇒ **合计 7 件、≈300-460 行**（票面"三枚测试件"少算：实为 4 件测试＋1 件脚本） | **零字面改动**即可落地：`C26` 权威行 `PLAN.md:1382` 已把机制收成 `（Win: GetFinalPathNameByHandle）`；`SPEC-06 §4`（`:50-53`）是 Windows 原语流水线、**无** `realpath`／`lstat` 字样；`D7`（`:184-199`）只定平台等级与验证方式、不涉解析语义 ⇒ **三处都不必改**。★**必须动**的只有 `pathresolver_other.go:5` 那枚 `DEFERRED(macOS/Linux)`（结项即摘；`AGENTS.md §1.1` 要求标记⇄§5 1:1）。可选：把 `PLAN.md:2385` 那句"**macOS 侧**"补成"macOS／Linux 侧"——那是**改 C26 邻近正文＝须人工批准**，⛔ 本形的前置 | Linux 顶层 **117**：PASS **115→112**、FAIL **1→5**、SKIP **1→0**（全叶子 191+5＝**196**）；Windows 顶层 **142**：**140／1／1 不变**（改的文件是 `!windows`）。★**转述**：锚 `7437ef88`，距 HEAD **928 提交**；此后两枚**无 tag** 测试件又加 **7 枚**顶层测试（`ticket363_screen_l0_test.go` 3＋`provenance_160_boundary_test.go` 4）⇒ **绝对数⛔ 搬 HEAD** | 票 339 `AC#2` **第一次有世界**（入口是 `Classify`→`formsOf`：`blacklist.go:79`，⛔ 走 `Resolve`）；票 55 `AC :28` 的**产码半边被提前付清**（55 仍被 `Blocked by` 挡着、**不关票**） | ①`internal/tools` **一枚**登记 skip 失去理由（`paths_workspace_test.go:198` 逐字点 `reparseComponents returns nil`）⇒ Linux 新开跑＋`portable-tests.sh:707` GUARD 陈旧红（**未验**；★同行 `:710` 那枚**不**受影响，见 §3 末更正）②`normalizeLocalUNC`（`:249-268`，`:268 return u`）会把 POSIX 上 `//` 起头的绝对路径整条换成 `\` —— 今天**走不到**（`resolveHandle` 恒 false）、甲把它**唤醒**（**我读码得出，票面未写**）③只落 `resolveHandle` 不落 `reparseComponents` ⇒ 软链被"**跟随**"而非"拒绝"⇒ 相对今天**更 fail-open**（`SPEC-06 §4` 明禁 parse-then-continue）④macOS `/var`→`/private/var` 是软链 ⇒ `t.TempDir()` 造的世界整体撞 `ErrReparseDenied`（**未验**）⑤334 那形（hostile 名造 partial）在 POSIX **不复现**：日志 `:33-:41` 现证 class 仍 `B`、`clause=false` ⇒ "搬牙"实为"**重造牙**"⑥A 新可达 ⇒ 任何在 POSIX 上钉 `ClassB` 的断言可能翻（**未验**） | **可单独回退**（单提交；无数据迁移、无契约字节、无 `Status:` 字节）。★前提＝**搬牙与产码同一提交**；拆两笔则只回退产码会留"桩已回、牙已删"的红 |
| **二**<br>＝AC#1 **丙**<br>只钉自陈与口径 | 仅 `pathresolver_other.go:5-9`：把**出处**从"`SPEC-06 §4`"改成 `PLAN.md:2385`、把**射程**从 `non-Windows` 收回"**macOS 侧**"、把"until that ticket lands"改成具名双引用（55＋347）。⛔ 任何 `.go` **可执行字节** | ≈**10-20 行**注释／**1 件**；0 新文件 | **零契约文字改动**；且按 §4.2，**不必新增 `SPEC-12 §5` 行**（那枚标记今天**已有登记面**）⇒ 成本比票面设想的低 | **两平台与今天同**（同一锚：Linux 115／1／1、Windows 140／1／1），**零翻转** | 339 `AC#2` 与 55 `AC :28` **都还是零世界**（解析器仍是桩） | ①"零产码"≠"零 `.go` 足迹"——它仍改 `.go` 里一处注释（票面口径须先认这格）②敞口照旧：`D7` 说 macOS 在矩阵内、而 C26 的 POSIX 半边**仍零实现**③"没人知道它开着"这件事被消掉，但**洞里还是空的** | 可（一处注释） |
| **三**<br>不做 | 无 | **0 件／0 行** | 零改动。代价＝那枚注释的**两处不实**继续挂着（出处错＋射程过宽），且 `SPEC-06 §4` 的 POSIX 半边继续零实现 | **同二**（零变化） | **同二**（零世界） | ①＝二②的敞口 ②"它开着而没人知道"也没被消掉（票面自己列的具名后果） | 无动作可退 |

## 3. 五红分母：我自己判可信不可信（⛔ 复跑，说明为什么）

- **逐名对得上**：`339/v1/logs/v1-posix-resolver-lands-preview.log`（34276 B）里 `--- FAIL:` 恰 **5** 枚，与票面逐字同名：`TestFallbackBranchHasNoWorldOnThisPlatform336`(`:37`)、`TestAnchorMissHasNoWorldOnThisPlatform339`(`:42`)、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`(`:118`)、`TestSyncNoGradeIsConfirmedOnPosix`(`:361`)、`TestSyncMembershipDecidesOffProfilePosix`(`:365`)。
- **四数可自洽分解**：我数该日志得 顶层 `PASS=112`／`FAIL=5`／`SKIP=0`（总 **117**）、含子测试 `PASS=191` ⇒ **196＝191+5**、**112＝117-5**、**5**、**0**。票面的 `196/112/5/0` **就是这四个数**，与日志逐格对得上。
- **每枚红都有对应的守卫行**（我逐行读到，不是转述）：336 `anchor_fallback_teeth_336_other_test.go:75` `if f.root != ""`（报错在 `:76`）；339 `anchor_miss_reason_339_other_test.go:139`／`:142`（报错 `:140`／`:143`）；sync 两枚＝`syncdirs_other_test.go:63` `res.Resolved` 与 `:94` `SyncDetectionComplete()`（报错 `:64`、`:96`、`:100`）。四者的触发条件**都是"root/Resolved 由空变非空"**，与"桩→realpath"一一对应。
- **附加旁证**：`pre`→`post` 两发日志（`v1-linux-fullpkg-pre.log` 118＝114/3/1、`v1-linux-fullpkg-post.log` 117＝115/1/1）证明该锚下 `C26RewriteAccount…` **本来就是红的**（存量 `A380`，非本形引入），且 `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` 的 SKIP 在 preview 里消失 ⇒ 票面"SKIP→PASS"那句也自洽。
- **⇒ 判：可信。但要打两个折**：**(a) 绝对数已陈旧**（锚 `7437ef88` 距 HEAD **928** 提交，且其后新增 7 枚无 tag 顶层测试 ⇒ 今天的分母**不是** 117／142）；**(b) 射程只有 `./internal/risk`** —— `core` 档还跑 `./internal/tools/...`，那里**一枚**以"`reparseComponents returns nil`"为跳过理由的用例（`paths_workspace_test.go:198` 逐字）会在甲落地后**开始跑**，它不在这 5 枚里。⇒ ★**"5 红"是包级读数，不是门禁级读数**。〔★更正：`task_output_pointer_notice_test.go:408` 我原以为同族，现读**不是**——它的 skip 是"造不出真 NTFS junction（`mklink`／`cmd` 不存在）"，与 `reparseComponents` 无关，甲**不会**让它开跑。〕

## 4. 契约影响，逐形逐条（哪一形要改哪个文件的哪句话）

### 4.1 形**一**：**零契约字面改动**（三处都不必改）
- `C26`（权威行 `PLAN.md:1382`）逐字已把机制括注成 `（Win: GetFinalPathNameByHandle）`——"Win"这个限定本身就把非 Windows 机制留白，**落 POSIX 实现不构成改判据**。
- `SPEC-06 §4`（`:47-58`）是 Windows 原语流水线（`GetFinalPathNameByHandle(VOLUME_NAME_DOS)`／`FILE_ATTRIBUTE_REPARSE_POINT`／8.3／UNC），**逐字没有** `realpath`／`lstat`。★票面那句"`SPEC-06 §4` 要求 `realpath`＋`lstat`"**出处错位**：它逐字在 **`PLAN.md:2385`**（D33 的 C26 修复段），且原文限 **"macOS 侧"**。
- `D7`（`PLAN.md:184-199`）只定平台等级（Win 一等／macOS 二等／Linux"不承诺、不进官方支持矩阵"）与验证方式，**不涉解析语义** ⇒ 甲落不落地都不动 D7 一字。
- ★**必须动**：`pathresolver_other.go:5` 的 `DEFERRED(macOS/Linux)`（结项即摘）。**可选**：`PLAN.md:2385` 的"macOS 侧"扩成"macOS／Linux 侧"——那是**改 C26 邻近正文＝人工批准**（`SPEC-12 §4.1`），⛔ 甲的前置。

### 4.2 形**二**：**零契约改动**，且**不必新增 `SPEC-12 §5` 行**
- 票面 `AC#0②` 问"这一枚标记今天有没有登记面"。我量到 **有**：`SPEC-12 §5:66` `| DEFERRED | macOS 平台层 |` ＋ `§5:86` `| RESERVED | Linux 支持 |`；且 `docs/evidence/s1/225-deferred-registry-audit.md:83`（第 18 行）已按"**父项口径**"判过"**§5 有对应行**"。⇒ 二 的成本比票面设想低：只需把注释的**引文改准**（出处＋射程），`§5` 一行不用加。
- ⚠ 但那枚注释现在**两处不实**：把 `PLAN.md:2385` 的 macOS 口径说成"`SPEC-06 §4` requires … for **non-Windows**"——**出处错＋射程被拉宽到 Linux**。二 的真实动作就是修这两处。

### 4.3 形**三**：零改动；代价＝上面那两处不实**继续挂着**，且 CI 每天在 Linux 跑这个包时，"这一支不可达"没有任何具名登记面。

## 5. 票面缺陷（我读到的字面，供重裁时用）

1. **`AC#1` 的 甲/乙/丙 与 §前置 的 一/二/三 不是同一组三分**：`一＝甲`、`二＝丙`，而 `三＝不做` 在 `AC#1` 里**没有席位**；`AC#1` 的 **乙**（"把 `f.root` 换成另一条不依赖句柄的解析路径：词法链＋逐组件 `Lstat`"）在 §前置 三候选里**也没有席位**。⇒ 本表按编排者点名的 **一/二/三** 出（§2），乙的单列读数在下面，**⛔ 当成第四形**。
2. **乙 的代价（我读码推的，⛔ 量过）**：改 `formsOf`（`pathresolver.go:345-375`）加一条词法＋逐组件 `Lstat` 的 root 计算 ≈**25-35 行／1 件**，`resolveHandle` **留原样** ⇒ `Resolve` 仍 `Resolved=false` ⇒ `syncdirs.go:140/198/222` 三条不受影响 ⇒ **红只有 2 枚**（336／339 两枚 tripwire），sync 两枚**仍绿**。但代价是：root 变**词法根**（非真实根）⇒ `uncertainAnchorMiss` 的 tie-break 落回拼写比较＝正是 C26 要关的那枚洞；且 `formsOf` 与 `Resolve` 成了**两条解析路**⇒ 违反 `SPEC-06 §4:55`"白名单与黑名单都必须经 C26（同一个函数）"。⇒ **乙买不到真值。**
3. **三处行号陈旧**（本文已给现量）：`blacklist.go:359→:363`、`anchor_miss_reason_334_test.go:85-88→:107-110`（四枚 `t.Setenv`）、`portable-tests.sh:266/:282→:292/:309`。
4. **"三枚受它影响的测试件"少算**：实为 4 件测试＋`portable-tests.sh` 的 2 行 GUARD（见 §2 形一风险①）。

## 6. 我对"一"的看法：**判它对**，但票面把它的射程写小了

- **方向对**：`D7:198` 把 macOS 列在支持矩阵内（二等＝CI 交叉编译＋README 标注），`PLAN.md:2385` 逐字承诺"macOS 侧：`realpath`＋`lstat` 拒绝 symlink" ⇒ 这是**已承诺而未实现**，不是"远期移植项"。按铁律②（选"最有价值／最像主流 harness"），一个解析入口只该有**一条**解析路——甲把 `resolveHandle` 这条缝补上，乙/二/三都留着"两套口径"或"零口径"。
- **但票面把甲写小了，三处必须补进射程**（否则落地即 fail-open 或自伤）：
  1. **`reparseComponents` 也得落**（`:12` 那枚桩票面没点）。只落 `resolveHandle`＝软链被"跟随"而非"拒绝"⇒ 相对今天**更** fail-open，正撞 `SPEC-06 §4` 的"默认拒绝"。
  2. **`normalizeLocalUNC` 的 POSIX 唤醒**（`pathresolver.go:249-268`）：今天 `resolveHandle` 恒 false ⇒ 这枚函数在 POSIX 上收不到非空 `final`；甲一落，`//` 起头的绝对路径会被整条换成 `\`。**这是甲的第 2 个产码点，票面零提及。**
  3. **"搬牙"实为"重造牙"**：336 那 8 枚牙全靠 `hostile334`（U+001B／U+202E）"Windows 拒绝解析"造 partial；POSIX 上这两个是**合法文件名字符**（日志 `:33-:41` 现证 class 仍 `B`、`clause=false`），partial 制造器必须换成 **dangling symlink／000 目录**。⇒ AC#3 的"改其射程"低估了工作量，且**⛔ 承诺"搬牙即绿"**（票面 §前置 一 的 ⚠ 已经预警，我这次把它量成了**确定**）。
- **⇒ 结论**：选**一**成立；但正式落码前请把上面三条写进落地腿的射程，并把分母口径从"`./internal/risk` 5 红"改成"**包内 5 红 ＋ `internal/tools` 1 枚开跑 ＋ `portable-tests.sh:707` 1 行 GUARD**"。

## 7. "我核过" 与 "我没核" 的分界（⛔ 混）

**我核过（每条都有上表的痕迹）**：桩两枚本体与首行 tag；`f.root` 唯一赋值点（全仓 `git grep`＋现读分支）；`blacklist.go:363` 守卫与 `isUnder` 两行体（⇒ `isUnder("","")` 恒真）；`Classify`→`formsOf` 的入口（`:79`）；`portable-tests.sh` 档位与 `:707`／`:710` 两行；票 55 的 `Blocked by` 与 `:28`；`D7:184-199`；`PLAN.md:1382`／`:2385`；`SPEC-06 §4:47-58`；`SPEC-12 §5:66`／`:86`；`225` 审计 `:83`；四份 339 日志的**四数分解**与逐名红；五枚红各自的**守卫行**；HEAD／index／脏树；`normalizeLocalUNC:249-268`；`reparseComponents` 的跨包消费者（`paths.go:117` 之外只有 `paths_workspace_test.go:198` 一枚 skip 理由，另一枚 `:408` 现读**不属此族**）；928 提交漂移与 2 枚新测试件。
**我没核（⛔ 当已核用）**：⛔ **未跑任何 `go test`** ⇒ `196/112/5/0` 与两平台四数**全是转述**，且**是 928 提交前那棵树的**；⛔ 未在 HEAD 量四数（故本文⛔ 给"今天的分母"）；⛔ 未验 macOS `/var` 软链那条；⛔ 未验 `internal/tools` 那枚开跑后是红是绿；⛔ 乙形"2 红"是**读码推的**；⛔ 未读 `blacklist.go` 全部 A 规则（只读到 `aUnder/aEq` 那层结构）；⛔ 未核 `A380` 台账原文（"存量红"三字是从票面＋339-v1 转述的）；⛔ 未验"搬牙重造后 POSIX 上 8 枚牙真能逐个变红"。
