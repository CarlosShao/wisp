# 票 88 独立对抗验收 —— `ban #6` 这道门现在到底有没有牙

**验收代理：** `acceptor-ticket88`（独立对抗验收，不写生产码、不 push）
**被验收的票：** `.scratch/wisp/issues/88-arm-ban6-now-that-frontend-exists.md`（`Status: ready-for-review`，六框自勾）
**代码 SHA：** `84e4161`（`git merge-base --is-ancestor` ⇒ 确为 `HEAD` 的祖先）
**接续账：** docs commit `aa5b90c` / `557b458`
**我的仓外纯净快照：** `/tmp/wisp88acc-88`（`git archive 84e4161 | tar -x`，**后缀 `wisp88acc-88` 是我自己的会话名**，
与前任的 `/tmp/wisp88-base`、88b 的 `/tmp/wisp88b-flip`、`/tmp/wisp88b-head` 不同名）
**工作树状态：** 我只读不写；工作树里 `frontend/**`、`internal/{config,tools,risk}`、`cmd/wisp/panel_assets.go`
等**在飞改动全部未被我触碰**（`git status` 起手即录，收工再录，见 §9）。
**所有变异只在 `/tmp` 快照内做**，仓内**没建 worktree、没 checkout、没 stash**。

---

## 0. 结论速览（逐条独立复现，每条我自己的原文读数）

⚠ 下面**章节号 = 任务的检查项号（1~6）**，不是执行时间顺序（我先跑完检查项 1 与 4 就交了第一枚 checkpoint commit）。

| # | 检查项 | 裁决 |
|---|---|---|
| 1 | 台账真数（纯净快照 `sh scripts/d22scan.sh` rc=0 + `ban #6 examined N>0`） | **PASS**（我的 N=**35**，见 §1） |
| 2 | 牙口正向（种 `approval.decide` 进快照 `frontend/` ⇒ rc=1 并点名） | **PASS**（我的种子 `.tsx` ⇒ rc=1 并点名 `:2`，见 §2） |
| 3 | 牙口反向（删 `frontend/` 但 `live:true` ⇒ 致命，不许变绿） | **PASS**（二进制 **rc=2**，`clean` 出现 0 次，见 §3） |
| 4 | "空仪器不算判据"未被削弱（ban 文本/正则/规则原文/allowlist 5 行） | **PASS**，见 §4 |
| 5 | AC#3 五条台账用例（`-count=1`，逐条读断言本体 + 变异） | **PASS**（RUN 28/PASS 28/FAIL 0/SKIP 0，六枚变异全红，见 §5） |
| 6 | AC#1 的"无后缀过滤器"是否真成立 + 钉子变异 | **PASS**（`main.go:742`/`:183`/`:817` 原文，钉子加回白名单即红，见 §6） |

---

## 1. 台账真数：我自己的纯净快照读数（AC#1 / AC#5 的门禁本体）

命令（原文）：
```
rm -rf /tmp/wisp88acc-88 && mkdir -p /tmp/wisp88acc-88
git archive 84e4161 | tar -x -C /tmp/wisp88acc-88        # rc=0
cd /tmp/wisp88acc-88 && sh scripts/d22scan.sh            # 逐字同 CI
```
**我的 rc = 0**。逐作用域**原文**（一字未删，全量输出）：
```
d22scan.sh: positive control - go test ./... (tools/d22scan)
ok  	github.com/CarlosShao/wisp/tools/d22scan	4.050s
d22scan.sh: scan of /tmp/wisp88acc-88
d22scan: examined 207 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/wisp88acc-88
d22scan: scope bans #1-5 internal/      examined 187 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  35 text files
d22scan: scope ban #7 internal/tools/   examined  16 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 internal/         examined 314 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=187, bans #1-5 cmd/=20, ban #6 frontend/=35, ban #7 internal/tools/=16, ban #8 design/=16, ban #8 internal/=314, ban #8 cmd/=25; ban #8 emoji coverage: design/ 16 text files; internal/ 314 Go files, comments and _test.go included; cmd/ 25 Go files, comments and _test.go included
```
**我的 N = 35（> 0）**，独立交叉核对：`find frontend -type f -not -path '*/node_modules/*' | wc -l` = **35**（同一快照）。
⇒ 报告器自报的文件数**与真实文件数一致**，不是"计数器凭空 +1"。

⚠ **我的 35 为什么与编排者/88b 的 35 相同**：**因为我 archive 的是同一枚 SHA `84e4161`，树没前移到我量的地方**——
相同数字在这里**不是巧合也不是抄来的**，而是"同一棵树的同一把尺"必然结果。
真正的差异发生在**树前移之后**：88b 自己记录工作树里是 38（票 77 在飞加文件），
所以 §7 我另量了 **`HEAD` 纯净快照**那一档（多样本全报），看的是**我的数会不会随树漂**。

---

## 4. 这三处是本票最可能偷偷让步的地方 —— 逐行证明**没有**

### 4.1 没有任何 ban 的文本/正则被改窄

```
git diff 84e4161^ 84e4161 -- tools/d22scan/main.go | grep -E '^[-+]' \
  | grep -E 'regexp\.MustCompile|banned \(|is banned|panelCheck|approvalPanelRe|MatchString'
⇒ 零命中（grep rc=1）
```
再上一道**字节级**保险（把 9 条 ban 正则所在的声明块整块取出来 sha256）：
```
git show 84e4161^:tools/d22scan/main.go | sed -n '95,112p' | sha256sum   = 65dbe9bd813ab357990f79ffae6055d2c5e433d5af5392bb8259131ded94caa5
git show 84e4161  :tools/d22scan/main.go | sed -n '97,114p' | sha256sum   = 65dbe9bd813ab357990f79ffae6055d2c5e433d5af5392bb8259131ded94caa5
⇒ 完全相同（含 `approvalPanelRe = regexp.MustCompile(\`approval\.decide\`)` @ 84e4161:main.go:105，
   `panelCheck` 的命中语 @ :711）
```
⇒ `approval.decide` 这条**判定文本与正则一个字都没动**，D22 的所有权没被代理碰。

### 4.2 "live 作用域 examine 0 个文件 ⇒ 致命"这条规则原文未动

```
规则文案（`an empty instrument is not a verdict (ticket 71 AC#4)`）在 84e4161^:main.go:913
与 84e4161:main.go:962 都存在，行号漂移=纯插入造成的位移。
判据函数本体哈希：
  git show 84e4161^:tools/d22scan/main.go | sed -n '/^func emptyLiveScope/,/^}/p' | sha256sum
    = 9a0cadd297b51a8e413c096efd06d267449ef0e83fcbe65aa2593daa7fece80c
  git show 84e4161  ... 同上
    = 9a0cadd297b51a8e413c096efd06d267449ef0e83fcbe65aa2593daa7fece80c   ⇒ 逐字节相同
配套的反向守卫 `driftedAbsentScope` 函数体同样逐字节相同
  （两枚均为 bd44c25f2208ecb1e837395655574ea30ec33e7db4f9e865397a870eb78f503e）
```
`git diff 84e4161^ 84e4161 -- tools/d22scan/main.go` 全量 **132 行**我**整块读过**：
除注释、`declaredScopes()` 里 ban #6 那一条（`live: true`、去 `absentOK`、改 `note`）、
`scanScope.absentOK` 的**字段注释**、以及把 `verdict()` 拆成
`verdict() → verdictWithScopes(..., declaredScopes(root))` + 测试专用的 `fixtureVerdict()` 之外，**没有任何逻辑改动**。
⇒ 关键点：`main()` 走的仍是 `verdict()`，它**只**通过一次 `declaredScopes()` 取账；
`fixtureVerdict`（可注入豁免作用域的那个）**在 `_test.go` 之外无人调用**（我用 grep 复核，见 §5.1 第三条）。
⇒ **裁决：PASS**（这一条门没有被削弱，ban #6 现在反过来**归它管**了）。

### 4.3 `allowlist.txt` 仍是 5 行非注释

```
grep -v '^[[:space:]]*#' tools/d22scan/allowlist.txt | grep -c .
  @84e4161 = 5      @HEAD = 5
git diff 84e4161^ HEAD -- tools/d22scan/allowlist.txt   ⇒ 空输出（本票与其后所有 commit 都没碰它）
```
⚠ 顺带一条**不属于本票**的观察（票 96 已由编排者建号，我不在此裁定）：`isTextFile()` 只出现在
`main.go:817`（在 `walkEmoji` 体内，函数边界 `main.go:793`–`main.go:832`），即 ban #8 的
`design/` 分支；`frontend/` **不在 `emojiScopes()` 的三条作用域里** ⇒ ban #8 对面板确实是空的。
这与票 88 的声称**不冲突**（它只管 ban #6），我不因此扣本票的框。

---

## 2. 牙口正向：我自己往快照 `frontend/` 种违规（AC#4 的独立复现）

```
cd /tmp/wisp88acc-88
printf 'export function onAllow() {\n  return approval.decide({id: 42, allow: true});\n}\n' \
  > frontend/src/acceptor-seed-88.tsx
cd tools/d22scan && go run . -root /tmp/wisp88acc-88      ⇒ RC=1
```
**我的原文**（只贴关键三行，全量在 `/tmp/wisp88acc-88-pos.txt`）：
```
d22scan: scope ban #6 frontend/         examined  36 text files
frontend/src/acceptor-seed-88.tsx:2: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
d22scan: 1 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
```
⇒ **PASS**。要点三条：(a) rc≠0；(b) **点名了我的文件与行号**（`:2`，正是 `approval.decide` 那一行）；
(c) 我种的是 **`.tsx`**（票 77 面板的真实语言），**不是**台账里那枚 `.js` 老种子 ⇒ 牙口不依赖某个后缀的巧合。
另外 examined 从 35 变 **36**：计数器与真实文件数一起动，说明它不是"报个固定数"。

## 3. 牙口反向：删掉 `frontend/` 而 `declaredScopes()` 仍 `live:true`

```
cd /tmp/wisp88acc-88 && rm frontend/src/acceptor-seed-88.tsx && rm -rf frontend
cd tools/d22scan && go build -o /tmp/wisp88acc-d22.bin .            # build rc=0（编译成功，变异不算空）
/tmp/wisp88acc-d22.bin -root /tmp/wisp88acc-88                      ⇒ BINARY_RC_TREE_GONE=2
```
**我的原文**：
```
d22scan: scope ban #6 frontend/         examined   0 text files
d22scan: scope ban #6 frontend/ examined 0 files but is declared live in declaredScopes() - an empty instrument is not a verdict (ticket 71 AC#4)
```
计数核对（我数的，不是读票面）：这份输出里 `clean` 出现 **0 次**、`NOT COVERED` 出现 **0 次**
⇒ **绝不许静默变 0 覆盖变绿**这条做到了：它既没绿，也没把"树没了"粉饰成"这块不用扫"。
⚠ 一个**读数坑**要给后来的代理写明：`go run` 会把子进程的 2 折成自己的 1（它额外打印 `exit status 2`），
所以"rc=2"这个判据必须**用编译出来的二进制量**，像我这样；从 `scripts/d22scan.sh` 里看到的是非零（1），
非零这件事本身成立，**但别把 1 当成"只是 finding"那一档**。

票面声称的用例名 `ban_6_tree_gone_while_declared_live_exits_2`：在我的 `-count=1 -v` 全量名单里
**确实存在且 PASS**（见 §5 名单），而**我上面的手工复现独立于那条用例**——它不靠测试自证。
⇒ **PASS**。

---

## 5. AC#3 那 5 条台账用例：`-count=1` 真数 + **逐条读断言本体** + 变异

命令与 tally（快照 `/tmp/wisp88acc-88`，全量输出存 `/tmp/wisp88acc-test.txt`）：
```
cd /tmp/wisp88acc-88/tools/d22scan && go test -count=1 -v ./...     ⇒ RC=0
=== RUN   : 28 行，去重后 28 ⇒ RUN 数 == 不同名数（无重复计数虚高）
--- PASS (顶层)  : 19
    --- PASS(子) : 9         ⇒ 合计 PASS 28，与 RUN 相等
--- FAIL         : 0
--- SKIP         : 0（顶层+子都是 0）
gofmt -l .  ⇒ 空      go vet ./... ⇒ rc=0（在本模块目录里跑）
```
28 个名字（我的原文名单，逐条可点）：`TestScanDetectsAllSeededViolations`、`TestScanCleanRepoIsGreen`、
`TestAllowlistSuppressesOnlyListedPaths`、`TestCheckRootRejectsBlindRoots`(+4 子)、
`TestScanAloneIsNotAFalsifier`、`TestCheckRootAcceptsRealRepo`、`TestScannerSelfScanOfRealRepoIsGreen`、
`TestEmojiBanCoversGoSourcesNotJustDesign`、`TestDeclaredEmojiScopeCannotWalkZeroFiles`、
`TestScopeReportMatchesRealCoverage`、`TestVerdictGreenOnFullyLiveFixture`、`TestVerdictRedOnEmptyBan7Scope`、
`TestVerdictRedOnGoScopeWithOnlyTestFiles`、`TestExemptScopeCannotOutliveItsAbsentTree`、
`TestBan6ScopeIsNotNarrowedByAnExtensionFilter`、`TestVerdictRedOnUndeclaredCounter`、
`TestLedgerCountsMatchAnIndependentWalk`、`TestRealRepoLedgerIsHonest`、`TestBuiltBinaryGoesRedEndToEnd`(+5 子：
`seeded_violation_exits_1` / `empty_live_scope_exits_2` / `armed_ban_6_goes_red_on_the_panel_violation_exits_1` /
`ban_6_tree_gone_while_declared_live_exits_2` / `fully_live_fixture_exits_0`)。

### 5.1 五条被点名用例的**断言本体**逐读（不是测试名）＋ 是否"恒真"

| 用例 | 重述后的断言本体（我读到的） | 恒真？ | 我的牙齿证明 |
|---|---|---|---|
| `TestVerdictGreenOnFullyLiveFixture` | 原 `if !sc.live && sc.count(s)!=0`（"豁免者必须 0 覆盖"）**翻成** `if sc.live && sc.count(s)==0 → Errorf`，外加 `uncoveredScopes(scopes)` 必须为空 | 否（两个方向都可失败） | M6 见下 ⇒ rc=2 红 @ `scan_test.go:477` |
| `TestVerdictRedOnEmptyBan7Scope` | **函数体一字未改**（diff 里没有它）：掏空 `internal/tools/ok.go` ⇒ `examined==0` 前置 + `code!=2` Fatal + 报错必须点名 `ban #7 internal/tools/` | 否 | M4、M6 ⇒ 红 @ `:519`（ban #6 抢在前面变空 ⇒ 点名错位被抓） |
| `TestExemptScopeCannotOutliveItsAbsentTree` | 改成**合成 exempt 作用域**（`synthetic exempt/exempt-fixture-tree`），**双向**断言：树不在 ⇒ `driftedAbsentScope==""` 且 rc=0 且仍打 `[NOT COVERED]`；树出现 ⇒ `==exempt.label`、`fixtureVerdict` rc=2、报错点名该 scope、且 finding 仍留在 stdout。前置还断言真账本里 exempt **为 0 条** | 否 | M3 ⇒ 红 @ `scan_test.go:622`；M7 ⇒ 红（前置那条） |
| `TestRealRepoLedgerIsHonest` | 原"恰好 1 条 uncovered 且必须是 ban #6"、"ban #6 examined **!= 0 就报错"，**翻成** "uncovered 必须为 0 条"、"examined **== 0 就报错**"、"clean 行不许出现 `NOT COVERED`"、"clean 行必须带 ban #6 的真实计数" | 否 | M1（两条都退回旧预期）⇒ 红 @ `:807`、`:813`；M7 ⇒ 红 |
| `TestBuiltBinaryGoesRedEndToEnd` | `want string` → `wantAll []string`（**多枚子串都要在**，是加强不是放宽）；旧 case `exempt scope whose tree appeared exits 2` 被**替换**为两枚新 case：`armed ban 6 goes red on the panel violation exits 1`（wantAll 含 `panel-approval`/`approval.decide`/`ban #6 frontend/`/`frontend/src/app.js`）与 `ban 6 tree gone while declared live exits 2`（wantAll 含 `empty instrument`） | 否 | M2（把后一条退回豁免时代的 setup+预期）⇒ `exit code: want 2, got 1` 红 @ `:952` |

**"删除"这一项我单独核了**：顶层测试函数 18 → **19**，`git diff` 里 `^-func Test` **零命中**
⇒ 没有一条用例被删掉；唯一被换掉的是 e2e 的**一个子 case 名**，它的规则牙齿搬到了
`TestExemptScopeCannotOutliveItsAbsentTree`（合成 scope，双向）—— 而**不是**搬成"没人测了"。

### 5.2 变异记录（全在 `/tmp/wisp88acc-mut`，快照副本。**每枚都编译通过并真的跑出了断言**——
M3 那枚因插入 `return ""` 触发 `go vet` 的 `unreachable code` 告警（**vet 不是编译失败**，测试照常跑出 FAIL），
其余各枚 `gofmt -l` 为空、M2 那枚 `go vet ./...` rc=0。
⚠ 诚实一条：M2 我的**第一次**尝试因自己的转义错误写坏了字符串字面量（`go test` 报 `[setup failed]`），
**那一次作废、不计变异**，重做后才跑出下面的红。"编译失败不算变异"我是照做的。）

| 变异 | 内容 | 结果 |
|---|---|---|
| **M1** | `TestRealRepoLedgerIsHonest` 两条断言**退回旧预期**（`len(unc)!=1\|\|unc[0]!="ban #6 frontend/"`、`examined != 0`） | **FAIL** `:807` + `:813` ⇒ 它今天真的在断言翻牌后的新事实 |
| **M2** | e2e `ban_6_tree_gone_while_declared_live_exits_2` 退回豁免时代预期（seed `frontend/app.js` + 仍要 rc=2） | **FAIL** `:952 exit code: want 2, got 1`；其余 4 子 case 仍 PASS |
| **M3** | 生产码 `driftedAbsentScope()` 直接 `return ""`（反向守卫哑化） | **FAIL** `scan_test.go:622` ⇒ 合成 exempt 那条不是永真 |
| **M4** | 生产码 `emptyLiveScope()` 直接 `return ""`（"空仪器致命"整条规则哑化） | **FAIL** ×3 顶层（`TestVerdictRedOnEmptyBan7Scope`、`TestVerdictRedOnGoScopeWithOnlyTestFiles`、`TestBuiltBinaryGoesRedEndToEnd`）+ ×2 子（`empty_live_scope_exits_2`、`ban_6_tree_gone_while_declared_live_exits_2`） |
| **M6** | 把 `liveFixture()` 里 `frontend/src/panel.tsx` 那枚种子删掉（fixture 回到 ban #6 空跑） | **FAIL** ×2：`TestVerdictGreenOnFullyLiveFixture :477`（rc=2）与 `TestVerdictRedOnEmptyBan7Scope :519`（点名被 ban #6 抢走） |
| **M7** | **把 ban #6 整个退回 `live:false, absentOK:true`**（票 88 之前的状态） | **FAIL** 4 顶层 + 4 子：`TestVerdictGreenOnFullyLiveFixture`、`TestExemptScopeCannotOutliveItsAbsentTree`、`TestRealRepoLedgerIsHonest`、`TestBuiltBinaryGoesRedEndToEnd/{seeded_violation_exits_1, armed_ban_6_…, ban_6_tree_gone_…, fully_live_fixture_exits_0}`。⚠ 连 `seeded_violation_exits_1` 都塌 ⇒ 说明这枚布尔**在台账里是载荷**，不是装饰 |
变异目录收工核对：`diff -rq /tmp/wisp88acc-mut/tools/d22scan /tmp/wisp88acc-88/tools/d22scan` ⇒ **无输出**（`ALL_MUTATIONS_REVERTED_IDENTICAL`）。

⇒ **AC#3 裁决：PASS**（5 条都被重述成"翻牌后的新事实"且**逐条可被打破**，没有一条被我读成恒真）。

---

## 6. AC#1 的"没有后缀过滤器"这句：**读码原文 + 钉子变异**

票面这句话我逐点验，file:line 是我自己在 `84e4161` 快照里量的：
- `tools/d22scan/main.go:183` ⇒ `s.walkText(filepath.Join(root, "frontend"), "panel-approval", s.panelCheck, false)`
  —— 第 4 个参数 `goOnly=false`。
- `tools/d22scan/main.go:728` ⇒ `func (s *scanner) walkText(dir, ban string, check func(string) (string, bool), goOnly bool) error`。
- `tools/d22scan/main.go:742` ⇒ **walkText 里唯一**一句后缀判断：
  `if goOnly && (!strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go")) { return nil }`
  ⇒ `goOnly=false` 时短路，**这句对 ban #6 一个文件都不挡**。
- `tools/d22scan/main.go:749` ⇒ `s.examined[ban]++` 在 `os.ReadFile` 之后、逐行 `check(line)` 之前
  ⇒ "examined 35" 的语义就是"35 个文件被读进内存并被逐行扫"，不是"35 次目录项"。
- `tools/d22scan/main.go:737` ⇒ 目录级跳过只有 `testdata` / `node_modules` / `.git` 三条（对**所有** ban 一致）。
- `tools/d22scan/main.go:817` ⇒ **全仓唯一**一处 `isTextFile(path)` 用在非 go 后缀白名单上的地方，而它在
  `walkEmoji` 体内（`main.go:793`–`main.go:832`），**不在 walkText 路径上**；定义在 `main.go:834`。
⇒ **"没有后缀过滤器"这句成立**。我自己的 N=35 构成（`find frontend -type f -not -path '*/node_modules/*' | sed 's/.*\.//' | sort | uniq -c`）：
16 tsx / 6 json / 3 ts / 3 mjs / 2 css / 1 md / 1 html / 1 go / 1 .gitkeep / 1 .gitignore = **35**，
其中 `.mjs` 与 `.gitkeep`/`.gitignore` 这类**不在 `isTextFile()` 清单**里的文件也在数内 ⇒ 直接印证"没有白名单在缩它"。

**钉子变异（同一枚"为了整齐加白名单"的修法，我加的）**：在 `main.go:742` 之后插入
```go
if !goOnly && !isTextFile(path) { return nil }
```
```
go test -count=1 -run TestBan6ScopeIsNotNarrowedByAnExtensionFilter -v  ⇒ RC=1（FAIL）
  scan_test.go:667: ban #6 examined 5 files, want 7 - walkText's filter and this list disagree, and one of them is narrowing the ban
  scan_test.go:677: ban #6 did not fire in frontend/scripts/dev.mjs …
  scan_test.go:677: ban #6 did not fire in frontend/Procfile …
全套 go test -count=1 ./... ⇒ 2 条顶层红：该钉子 + TestLedgerCountsMatchAnIndependentWalk
```
⇒ 钉子**有牙**，且**不只一枚牙**：`TestLedgerCountsMatchAnIndependentWalk` 在 `scan_test.go` 新增的
`{"ban #6 frontend/", count(..., func(string) bool { return true }), s.examined["panel-approval"]}` 一行
用一次**独立 walk** 交叉核对同一个数（票面 log 第 2 条预告过这条会抓住收窄，我实测真的抓）。
⇒ **AC#1 裁决：PASS**。

---

## 7. 多样本全报（不挑运气那次）+ 与编排者数字的差异说明

| 样本 | 树 | `sh scripts/d22scan.sh` rc | 我的 `ban #6` N | 备注（我的原文） |
|---|---|---|---|---|
| S1 | `84e4161`（本票代码 SHA） | **0** | **35** | `ban #6 frontend/ examined 35 text files`，clean 行带 `ban #6 frontend/=35` |
| S2 | `HEAD`（`5b1d855`，含票 77/90/96 的东西） | **1** | **37** | 红因是 `internal/winsec/winsec.go:126: [pathresolver-bypass]`，**与 ban #6 无关**；`TestScannerSelfScanOfRealRepoIsGreen` + `TestRealRepoLedgerIsHonest` 两条红 |
| S3 | `HEAD` + **只在快照里**给 winsec 补一行 allowlist（沙箱变异，仓内一字未动） | **0** | **37** | `clean … ban #6 frontend/=37` ⇒ 抽掉那笔非本票的红账之后，翻牌后的门在最新树上是绿的且真扫 37 个 |
| S4 | S3 之上往 `frontend/src/` 种 `.tsx` 的 `approval.decide` | **1** | 38 | `frontend/src/AcceptorSeed88.tsx:2: [panel-approval] …` |

**关于"我的 N=35 与编排者的 35 相同"**：因为我 archive 的就是同一枚 SHA，树没前移 ⇒ 数必然相同；
**我并没有少扫**（`find` 独立数出同一枚 35）。差异只出现在树前移之后：**HEAD 是 37**（票 77 又落了
`frontend/scripts/render-l2.tsx` 与 `frontend/fixtures/l2-card-fs-delete.json`）。
88b 记的"工作树 38"是她当时工作树在飞的读数，与我这两档 SHA 读数**不构成矛盾**。

⚠ 一条**仪器坑**（不是票 88 的退化，但 owner 该知道）：`scripts/d22scan.sh` 第一步是**不带 `-count=1`** 的
`go test ./...`。我在 S4 那棵**已种违规**的树里实测：
```
go test -count=1 ./...  ⇒ --- FAIL: TestScannerSelfScanOfRealRepoIsGreen / --- FAIL: TestRealRepoLedgerIsHonest ⇒ FAIL
go test ./...           ⇒ ok  github.com/CarlosShao/wisp/tools/d22scan  (cached)
```
⇒ 本机热缓存下第一步可以**被 cache 端过去**；CI 的 `lint` job **不依赖它**：CI 用的是
`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`，`runtests.sh:75` 是
`go test -v -count=1 "$@"`（并断言顶层 PASS/FAIL 数 > 0、任何 SKIP 致命），且 `lint` job **没有** `actions/cache`
⇒ 每次冷缓存。脚本第二步 `go run . -root` 更是永不缓存。⇒ 我在 §1 报的 rc 取自 S1 那次
（它的 step1 打的是 `4.050s` 真实耗时，不是 `(cached)`），**门禁不吃缓存**这件事我是自己复算过的。

---

## 8. 四种假绿：逐条点名（我这一路）

1. **"pattern 打空造成的 ok"**：我唯一用 `-run` 的三次（M2/M3/M6 的定位跑）都配了 `-v` 并点名 `=== RUN` 行——
   M2 那跑 `=== RUN` **6 行**（父 1 + 子 5，>0），且**目标子用例自己打了 FAIL**；M3 那跑目标测试打了 `scan_test.go:622`；
   M6 那跑两条 `:477`/`:519`。§5 的全量跑**不带 `-run`**。⇒ 无此假绿。
2. **"SKIP 冒充 ok"**：全量 `-count=1 -v` 里 `--- SKIP`（顶层+子）= **0**，`--- PASS` 合计 28 == `=== RUN` 28。⇒ 无。
3. **"缓存冒充新证据"**：`-count=1` 用于全部判定跑（S1/§5/各变异）；`scripts/d22scan.sh` 逐字同 CI 那次
   打的是 `4.050s` 非 `(cached)`。§7 末把 step1 可被 cache 端过去这件事**明写出来**并给出 CI 侧为何不受影响。⇒ 无（且已登记）。
4. **"阈值被偷偷调低/挑数字"**：我没改任何阈值、没加 allowlist（仓内）、没删断言；
   票面的 5 个数字（N=35、5 条红名单、allowlist 5 行、`=== RUN`=28、`ban_6_tree_gone_while_declared_live_exits_2`）
   我**各自独立复现**后才采信；反向的 M7（把布尔退回旧值）是**故意让全套变红**来证明布尔是载荷。⇒ 无。

---

## 9. "别碰工作树"自证 + 现场噪声登记

- 我在仓内**只写过** `docs/evidence/s1/88-adversarial-acceptance.md`（唯一提交物）；
  `git status` 起手即录、每次 commit 前 `git diff --cached --name-only` 只含我这一个路径、
  commit 后 `git log --oneline -1` 确认 HEAD 真动。**所有变异在 `/tmp/wisp88acc-{88,mut,head}` 里做**，
  **仓内没建 worktree、没 checkout、没 stash、没 push、没 `--amend`**。
- 我起手时看到的工作树在飞条目（`M .github/workflows/ci.yml`、`A frontend/fixtures/l2-card-fs-delete.json`、
  `M frontend/package.json`、`A frontend/scripts/render-l2.tsx`、`M cmd/wisp/panel_assets.go`、
  `M internal/config/*`、`M internal/tools/*`、`?? internal/{config/permmode,risk/mode,tools/mode}.go`）
  **由票 77 / 90 的代理随后自行提交**（`14720af`、`1d4f289`），不是我动的。
- 中途 HEAD 前移了两次（`b4435d3` → `1d4f289` → `14720af` → … → `3d716a4` → `5b1d855`），我的 checkpoint
  commit 之上又叠了别人的 commit ⇒ 我的**全部判据都钉在显式 SHA**（`84e4161^`、`84e4161`、当时的 HEAD）上，
  不受工作树漂移影响。
- 我最后一次 `git status` 看到 `tools/d22scan/{main.go,scan_test.go}` **被别的代理在工作树里改动**（应是刚建的票 96），
  `internal/winsec/winsec.go` 也在改（票 94）。⇒ **本验收裁定的是 `84e4161` 的账本**，
  若票 96 随后动了 `walkText`/`declaredScopes`，**它的改动要由票 96 的验收表另算**。
  另登记一笔非我所为：`.scratch/wisp/issues/87-…md` 被别的代理改名为 `-done`。

---

## 10. 裁决表（与票面 AC **1:1**，六行）

| 票面 AC | 我的独立读数 | 裁决 |
|---|---|---|
| **AC#1** 真实读数 N + 过滤器说明 | S1：`ban #6 frontend/ examined 35 text files`，rc=0；`find` 独立数 35，构成 16tsx/6json/3ts/3mjs/2css/1md/1html/1go/1.gitkeep/1.gitignore；`main.go:742` 唯一后缀判断 gated on `goOnly`，`main.go:183` 传 `false`，`isTextFile` 只在 `main.go:817`（walkEmoji 内） | **PASS**（N>0，非空仪器；无收窄） |
| **AC#2** 翻牌落地、ban 文本一字不改 | `main.go:319-320` = `live: true`、无 `absentOK`、`note` 写明翻牌者/批次/理由并保留 D22 那句；ban 正则声明块 sha256 前后相同（`65dbe9bd…`）；diff 中无任何含 `regexp.MustCompile`/`banned (`/`MatchString` 的 +/- 行 | **PASS** |
| **AC#3** 5 条红逐条重新表述并真绿 | `-count=1 -v`：RUN 28（去重 28）/PASS 19+9=28/FAIL 0/SKIP 0；5 条断言本体逐读（§5.1）无恒真；顶层测试函数 18→19、**零删除**；M1/M2/M3/M4/M6/M7 六枚变异全部如实变红 | **PASS** |
| **AC#4** 阳性对照仍能红并点名 `ban #6` | 台账内：`armed_ban_6_goes_red_on_the_panel_violation_exits_1` 在跑（其断言的是子进程 rc=1）；我另在快照 `frontend/src/*.tsx` 独立种违规 ⇒ 二进制 rc=1、点名文件:行号、打 `[panel-approval]` 原文；老种子 `frontend/src/app.js`（`scan_test.go:80`）仍在，**未改后缀绕过** | **PASS** |
| **AC#5** 纯净树逐字同 CI rc=0 + 禁改项核对 | S1 rc=0 且 clean 行带 `ban #6 frontend/=35`、无 `NOT COVERED`；allowlist 非注释 5 行（84e4161 与 HEAD 都是），`git diff 84e4161^ HEAD -- allowlist.txt` 空；`emptyLiveScope`/`driftedAbsentScope` 函数体前后逐字节相同 | **PASS（附一笔明账）**：同一条命令在 **HEAD 的纯净快照里 rc=1**，红因是**票 89 的** `internal/winsec/winsec.go:126 filepath.Abs` + allowlist 无 winsec 条目（`grep -c winsec allowlist.txt` ⇒ 0）。**与本票无关**，88b 已登记、编排者已建票 94 ⇒ 我不因此扣本票的框，但 **push 前必须知道这一点** |
| **AC#6** 真树若命中不改扫描器、登记给编排者 | 我自己两档快照各数一次：`grep -rn --exclude-dir=node_modules --exclude-dir=testdata --exclude-dir=.git 'approval\.decide' frontend/` ⇒ 84e4161 **0 命中**（grep rc=1）、HEAD **0 命中**（rc=1）⇒ 票 77 vendored 的 Approval Card 不需改名；本票确实没为任何人消音（allowlist 一字未动） | **PASS** |

### 最后两句

**1) 票 88 能否改名 `-done`？—— 能。**
六框我全部以**我自己的原文读数**独立复现为 PASS（§10 表），三条最可能让步处（ban 文本、空仪器规则、allowlist 5 行）
我以 sha256 / 逐字节函数体哈希 / 行数三样各验一遍，且用 M1/M2/M3/M4/M6/M7 六枚变异证明这些绿**都可被打破**。
残留一条**不属本票**的前置：AC#5 的"纯净树 rc=0"在**今天的 HEAD** 上因票 89 的 `filepath.Abs` 而红（票 94 在治）；
本票的框按票面文字（`84e4161` 的账本 + 翻牌后的门禁）成立，**`-done` 后缀可以挂**，
但**票 94 未闭之前 push 上去 `lint` 的 D22 两步仍会红**（S2 实测）。

**2)"今天面板里如果出现 `approval.decide`，CI 会不会红"—— 现在可以对 owner 说"会"。**
依据（三条都是我今天跑的，不是引用）：
(a) **CI 真的调这道门**：`.github/workflows/ci.yml` @HEAD 的 `lint` job（ubuntu）第一步
`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（`runtests.sh:75` 强制 `-count=1` 并断言顶层 PASS/FAIL>0、SKIP 致命），
第二步 `run: sh scripts/d22scan.sh`（`set -eu`，无 `|| true`）；全文件 `continue-on-error` **0 处**，
且 D22 那步被**放在 job 最前**（票 71 的教训：不再被 `go vet` 的红挡成 skipped）。
(b) **门真的看得见面板那棵树**：S1 里 `ban #6 frontend/ examined 35`、HEAD 里 `37`（§7），
且加回后缀白名单的修法会被钉子 + 独立 walk 双双抓红（§6）。
(c) **今天把它种进去就红**：S4 实测——在 HEAD 纯净快照（且把 winsec 那笔红在沙箱里 allowlist 掉、
即**排除一切别的红因**）往 `frontend/src/*.tsx` 放一句 `approval.decide` ⇒ `sh scripts/d22scan.sh` **rc=1**，
打印 `frontend/src/AcceptorSeed88.tsx:2: [panel-approval] …`；
同一棵树上 `go test -count=1 ./...` 亦红两条（self-scan + ledger honest）。
⚠ 边界要说清：这道门的范围**只有 `frontend/`**（`node_modules/`、`testdata/` 被跳过），
所以准确的说法是"**面板树里出现 `approval.decide` ⇒ CI 红**"；
它**不**保证"`approval.decide` 出现在任何别处也红"（别处归 ban #6 之外的账，例如票 96 正在补的 ban #8 覆盖面）。
