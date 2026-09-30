# 224 — 会话授权「判据返工」对抗验收表 · 腿 `224-v2`

- 被验收对象：六枚 commit `b7ca76b5`（骨架）→ `2a0175ba`（N#1）→ `1c601fab`（N#1＋N#2）→ `ba5db093`（N#3）→ `e8ed5fef`（N#5）→ `43d9096f`（N#6）
- 工单：`.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md` §「续做」（74 行版）
- 判语来源（欠账清单）：`docs/evidence/s1/224-session-grant-v1.md`
- 死腿自述：`docs/evidence/s1/224-r2-progress.md`（§1-§5 正文写完；**§6 门禁四数与 §7 判不动仍是「（待填）」＝本表不代它填、不把那两节当它已交**；§6 由本腿自己现跑取数、§7 是本腿自己写的）
- 本腿＝**非实现者验收腿 `224-v2`**；产码腿 `224-r2` 死于 150 轮上限，其 commit message 与自述一律按**未验证**处理，逐条由本腿重跑
- 锚点：起手 `git log -1`＝`ba9e3575`（编排者 A468）；本表骨架 commit＝`3dc780e1`
- 本表生成时刻 `2026-09-30 13:5x +08`
- ⛔ 票面 `- [ ]`（AC#1-AC#5 与 N#1-N#6）**一枚未碰**，勾框归编排者
- ⛔ `internal/audio/**`＝票 241 地界：零碰、零跑、零归因（§6 的四包读数不含它，口径见 §6 注）
- ⛔ `frontend/**`／`design/**`：两层禁（不读、结论不引）
- 射程：只碰 `internal/agent/approval`／`internal/session`／`internal/tools`／`cmd/wisp`

## 口径声明

- 〔我现跑〕＝本腿本机此刻执行并贴回读数；〔读台件〕＝只读代码／台账／commit，未执行。
- 用例枚数口径：`grep -c '^func Test' <file>`（顶层 Test 函数数），不是断言数、不是子用例数。
- 「被扫文件数」≠「违规数」（`d22scan` 的 `internal/=223` 是前者）。
- 行号一律本腿现取：工单、`224-v1` 表、`224-r2` 自述里的行号只当线索，不当凭据。
- 起手名册（本腿终态自证的尺）：`git status --porcelain -- internal cmd`＝**空**〔我现跑 13:34:27〕。
- 编排者给的起跑基线（**非我跑的，抄自派单**）：`go build ./...` 13:30:19 rc=0；五包 `-count=1` 全 ok
  （0.367／1.131／13.238／15.574／129.061s）；`gofumpt -l internal/ cmd/` 空；`d22scan` clean rc=0。
  本腿 §6 一律用自己 13:49-13:52 那发读数，不与上面互替。

---

## §0 逐格 1:1 表（续做五格 → 本表五节）

| 续做格 | 票面要求（要点） | 本表节 | 判语 |
|---|---|---|---|
| N#1 | 「答复⇒落行」常驻判据（M4 必须有指名用例红）＋真 `tool_call.grant_id` 列有用例 | §1 | **成立**（M4／M-R1／M5 三发都有指名红，正反两面） |
| N#2 | 真·跨进程重启用例＋处置三处指着不存在测试的注释；判定上限逐字核 | §2 | **成立（附条件）**：用例存在且①②③有牙；上限只到「第二次装配」；两处注释坐标过期；「台架不存在」那半句过强 |
| N#3 | `run_mode101_test.go` 注释半⇒带条件的事实句、逐字保留原防护 | §3 | **成立** |
| N#5 | glob 方言定案＋常驻钉防「看着有规则其实从不匹配」 | §4 | **附条件成立**：定案清楚、五枚钉、M6/M7 有牙；条件＝方言射程只盖斜杠形态，**生产真实形态（Windows 反斜杠）的通配行仍静默不生效且写侧不拒** |
| N#6 | 反射钉 setter 半射程 | §5 | **成立**（M9／M9b 两发：新钉红、旧窄钉照绿） |
| （N#4 编排者已裁＝一次进程存活期） | — | — | 不在本腿射程 |

一句话总账：**`224-v1` 那枚「摘掉 `Record` 四包一枚不红」的死穴已经真的闭上了**（本腿复跑同一发突变，两包四枚指名红）；
新暴露的一枚是 N#5——**方言钉防住的和生产线实际会写的字符串不是一个形态**。

## §1 N#1 「答复 ⇒ 落行」的常驻判据 — **成立**

### 交付物（读台件，枚数口径见上）

| 文件 | 顶层 Test 枚数 | 管哪一格 |
|---|---|---|
| `internal/agent/approval/ticket224_reply_grant_test.go`（326 行） | 4 | 答复路四跳：`Replies.AllowSession`→`Gate.allowSession`→`grants.Record`→审计行 |
| `cmd/wisp/ticket224_assembly_test.go`（574 行） | 3 | 装配根的写／读／跨启动三格 |
| `internal/tools/ticket224_grantid_column_test.go`（185 行） | 1 | 桥写的值经**真 `memory.Store`** 落进 `tool_call.grant_id` 并读回 |

改前绿名册〔我现跑 13:38:18，`-run Ticket224 -v`〕：approval 4 枚 PASS（0.036s）／session 14 枚 PASS（0.703s）／
tools 10 枚 PASS（0.181s）；〔13:38:20〕cmd/wisp `Ticket224|TestTicket101SessionGrantDoesNotCrossRestart`＝4 枚 PASS（ok 10.141s）。
⇒ 本表所有「红」都相对这份名册说话。

### M4 复跑（正向：摘掉 `gate.go:678` 的 `Record`）〔我现跑 13:40〕

改法与 `224-v1` §2 逐字同形：`id, err := g.grants.Record(ctx, tool, p)` → `id, err := int64(0), error(nil)`
（落盘标记 1 枚、`grep -c 'g.grants.Record(ctx, tool, p)'`＝0、`go build` rc=0、31010→31082 字节）。

```
internal/agent/approval  FULL  -count=1 -v：
--- FAIL: TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed
--- FAIL: TestTicket224ForgedSessionAnswerRecordsNothing
--- FAIL: TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow
--- PASS: TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped   ← 设计上不响，见下
FAIL github.com/CarlosShao/wisp/internal/agent/approval 0.329s

cmd/wisp -run Ticket224：
--- FAIL: TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun (0.96s)
    ticket224_assembly_test.go:146: 0 approval_grant rows for the session this run minted
    (sess_de8ddb9824e87f166f6acf53b907db33), want 1 - the card printed [C:\...\remembered-by-a-click.txt].
    The audit still says: [audit] wisp run: SESSION-MINT id=sess_de8ddb98...
（同发 TestTicket101SessionGrantDoesNotCrossRestart／LiveGrant…／ProductionSession… 三枚仍 PASS）

internal/session ok 0.670s ／ internal/tools ok 12.198s   ← 这两包看不见答复路，不是漏判
```

⇒ **与 `224-v1` 的「四包一枚不红」正面对立**：同一发突变现在打红 2 包 4 枚**指名用例**（不是 `[build failed]`、不是超时）。
红的第一行逐字可在源码里核到（`ticket224_reply_grant_test.go:199` 的 `recorder saw %d Record calls, want %d …`），
`224-r2` 自述那句不是修辞。

`TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped` 不响是**设计**：它是反向钉
（`:251`/`:287` 断「没有记录器时不许出现 `GRANT-RECORDED`／`grant_id=`」），M4 恰好也走 nil 分支。
⇒ 本格的红来自另外三枚，不是这四枚齐响；这一条与 `224-r2` §1 自述一致。

### 反向一发：`Record` 放回、但让它写错的行〔我现跑 13:44-13:45，M-R1〕

改法＝`internal/session/grants.go:161` `SessionID: l.id.String()` → `SessionID: "sess_deadbeefdeadbeefdeadbeefdeadbeef"`
（合法形状、固定串；落盘标记 1 枚、`go build` rc=0、15008→15073 字节）。

```
internal/session  -run Ticket224 -v： 8 枚指名 FAIL
  RecordWritesOneRowPerAnswer / CoveringMatchesItsOwnSession / CoveringDoesNotCrossToANewSession /
  ExpiredOrRevokedRowIsNotLive / WildcardRowCoversAConcretePath / DialectIsPathMatchNotRegex /
  PatternWithoutStarIsNeverReinterpreted / MalformedWildcardIsRefusedAtTheDoorAndIgnoredOnDisk
  （MintIsRandomAndValid / TestLiterals / RecordRefusesUnusableSubjects / LedgerFailsClosed /
    NewLedgerRefusesHandWrittenIDs / PatternCoversIsTheOnlyMatcher 六枚 PASS＝射程本不在写侧）
cmd/wisp：3 枚指名 FAIL
  SessionVerbPutsARowOnDisk…（:146 零行）
  LiveGrantStopsTheAssembledBridgeFromAsking（:311 仍弹 1 张卡／:316 decision="allow"／:319 grant_id=<nil>／
    :324 缺 GRANT-USE／:335 ledger rows=[]）
  ProductionSessionDoesNotSurviveRestart（:451 boot 1's own session cannot cover its own answer）
approval ok ／ tools ok ／ TestTicket101 PASS
```

⇒ **写错行也红**，且红在「命中授权就不再弹卡」那一格的行为观测上（不只是「行存没了」）。
顺带把 `224-v1` 的一条担心结掉：`CoveringDoesNotCrossToANewSession` 在这一发里也红 ⇒ 它不是「只断 0 行」的空钉。

**「写错 scope」那一支本腿没打突变**，因为它是〔读台件〕可判的更下层钉子：
`internal/memory/dao_misc.go:18` 逐字 `if g.Scope != GrantScopeSession {` ⇒ 错 scope 在 DAO 当场被拒，
`Record` 返回 error，答复路会走 `GRANT-RECORD-FAILED`（`gate.go:681-688`）而不是谎报。
⇒ 装配用例里那句 `got.Scope` 断言（`ticket224_assembly_test.go:151`）是第二道，不是唯一一道。

### 真 `tool_call.grant_id` 列（M5）〔我现跑 13:46-13:47〕

改法＝`internal/tools/bridge.go:1116` `if grantID != 0 {` → `if false && grantID != 0 {`（标记 1 枚、50699→50749、build rc=0）。

```
internal/tools -run Ticket224 -v：
--- FAIL: TestTicket224GrantedCallBooksAllowSessionGrantWithItsRow      grant_test.go:470 booked grant_id = <nil>, want the covering row's 4242
--- FAIL: TestTicket224GrantIDColumnRoundTripsThroughTheRealStore       ticket224_grantid_column_test.go:131
        "tool_call.grant_id on disk is NULL for a call a live grant covered: the column the frozen DDL
         has carried since SPEC-02 §3 is still never filled in practice"（并 dump 出真行 Decision:allow_session_grant GrantID:<nil>）
cmd/wisp：--- FAIL: TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking（:319）
```

⇒ `224-v1` §2-③「真 store 的 `tool_call.grant_id` 列今天零用例」这一格**闭合**，且牙在「桥→DAO→盘上列→第二个句柄读回」的连接上：
桩侧（`grant_test.go:470`）与真 store 侧（`:131`）同红，装配侧第三处也红。

### N#1 判语

**成立**。正反两面都有常驻指名判据（M4 摘掉⇒红；M-R1 写错行⇒红；M5 丢列⇒红），
`224-v1` §2 那条「答复⇒落行零判据」的欠账已在盘上闭合。

## §2 N#2 幻影注释与真·跨进程用例 — **成立（附条件）**

### 幻影是否还在

〔我现跑 13:38:20〕`go test ./cmd/wisp/ -count=1 -run TestTicket224ProductionSessionDoesNotSurviveRestart -v`
⇒ `--- PASS: TestTicket224ProductionSessionDoesNotSurviveRestart (2.99s)`／`ok 10.141s`。
⇒ **不再是 `no tests to run`**：那枚被三处注释指着的用例现在真实存在（`cmd/wisp/ticket224_assembly_test.go:406`）。
本腿另在 M1／M-R1／M5 三发突变下都看见它红，所以它不是「存在但恒真」。

### 三处注释处置得对不对

票面 N#2 给了两个可选项（写出真用例／把注释改成说实话），交件**两件都做了**，且改后的注释逐字带判定上限：
现读〔13:35〕`internal/session/grants_test.go:20`／`:101-106`／`:147-152` 三处均已改成
「该用例现在存在＋它声称得了的那一半＋⛔ 它的第二次是同一测试进程里的第二次装配」，
`internal/tools/grant_test.go:6-9` 第四处同样改成指真文件。⇒ **不留指向不存在之物的指针**这一条做到了。

⚠ 但**注释坐标漂移**两处（同一枚文件里残留的过期指针，轻形）〔我现跑 `grep -n`〕：
1. `cmd/wisp/ticket224_assembly_test.go:23` 写「Three comments in internal/session/grants_test.go (:20/:94/:134 in the ticket's anchor, **now :20/:94/:134 here as well**)」、
   `:388` 又写一遍 `:20, :94 and :134`；**实际落在 `:20`／`:101`／`:147`**（差 7 与 13 行——正是它自己那批编辑把注释撑开的）。
   ⇒ 「now …here as well」这半句今天是假话（工单锚点里的 `:94/:134` 本来就不动，是它自己动的）。
2. `:24-25` 说 `internal/tools/grant_test.go` 也「point AT THIS NAME」；改后那处指的是**文件名**
   （`ticket224_assembly_test.go`），不再指用例名 ⇒ 指名字枚数从 4 降到 3，头部描述仍是 4。
⇒ 判**不致命**（被指对象真实存在、可跑、有牙），但要按「注释不许过期」记一笔，见 §8 第 4 行。

### 票面⑩ 的三件控制与牙齿

| 控制 | 盘上断言（现读 file:line） | 本腿实测 |
|---|---|---|
| ① 两次装配铸出两枚不同 id | `:484-487` `id2 == id1` 直接 Fatalf | **M1 下红**：`:485` 逐字 `the second assembly minted boot 1's identity "sess_776973702d7069642d64657269766564"`〔13:45:32〕 |
| ② boot1 生产 id 写、boot2 生产 id 查 ⇒ 0 行 | `:495-499`＋`:523-526`（并 `:502-505` 断旧行没被删、留给审计） | M-R1 下红（`:451` fixture 自证半先响） |
| ③ 同一发 `fs.write` 走真 `rt.bridge.Execute` ⇒ 第二次**弹卡 ≥ 1** | `:510-512` 读 `rt.windowCount()` 差值；`:532-538` 另断 `decision != allow_session_grant` 且 `grant_id IS NULL` | 极性与 `224-v1` §6-b 第 3 条一致（测弹卡、不测「被拒」）✓ |
| 附带的「贵方向锁」 | `:436-439`／`:491-493` 两枚 boots 的 id 各自比对 `session-before-restart`／`session-after-restart` | 现读为真 ⇒ 派单增量第 2 条的两方向都齐了 |

**同发 M1 下四枚旧钉照绿**〔我现跑 13:45:32〕：`--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (3.02s)`。
⇒ `224-v1` §3「四枚全绿而永久通行证实际生效是可能的」由**盘上的新用例**抓住，不再只靠验收腿的临时探针。

### 判定上限逐字核（这是票面点名的那一发）

交件自述（`cmd/wisp/ticket224_assembly_test.go:34-49`）引 `run_mode101_test.go:72-76`＋`:143-158`，
并逐字写「a second call into the composition root **INSIDE THE SAME TEST PROCESS**」＋
「It does NOT prove 'a second OS process does'」＋「Do not cite this file as having settled AC#3's OS-process wording」。

本腿逐字复认〔现读＋现跑〕：
- `grep -n "a new process surface" cmd/wisp/run_mode101_test.go` ⇒ **真身在 `:73-74`**（`t101host` 的类型文档），
  不在派单与 `224-v1` §7 所引的 `:145-147`（那三行是 `start` 的函数体）。⇒ **交件改用的坐标是对的，前人引的那套是错的**（§8 第 5 行更正）。
- `:143-158 func (h *t101host) start` 体内是 `code := runTextTask(runSpec{...})` 直调，`runTextTask` 定义在 `cmd/wisp/run.go:155` ⇒
  **同进程第二次装配**，没有任何 `exec.Command`。⇒ 用例把自己的射程**只声称到它真测到的那一半**，这一条**核过＝真话**。

⚠ 但同一节里那句 **"…that harness does not exist yet"（真两枚 OS 进程的台架还不存在）过强**〔我现跑现读 13:48〕：
`cmd/wisp` 今天就有真·第二枚 OS 进程的原语——`secret_argv_windows_test.go:161 buildWispForTest`（编真 `wisp.exe`、
按 `scripts/build.ps1` 同法摆 DLL），并被 `resident_sink_nail_127_windows_test.go:197`、
`early_log_nail_130_windows_test.go:121`、`dataroot_128_windows_test.go:128` 用 `exec.Command(exe, …)`＋
`WISP_ENV=test`＋data-dir 环境变量起第二枚进程；而 `t101` 那套路数把 mockllm 的 URL 写进 `config.toml`
（`run_mode101_test.go:107-134`），第二枚进程读同一目录即拿到同一 stub。
⇒ 诚实的上限写法是「**没人为授权路接这根线**」，不是「这根线造不出来」。
这一字之差有后果：它决定 AC#3「进程重启」那一半是**硬墙**还是**欠做**（本腿判：欠做，可续）。

### N#2 判语

**成立（附条件）**。条件三条：① 判定上限只到「第二次装配」，票面 AC#3 的 OS 进程那一半仍无仪器（且**可**做）；
② 交件头两处注释坐标过期；③ 「台架不存在」那半句需按上面更正为「没接线」。
三条都不改变「幻影已消灭、①②③ 三件控制已在盘上且有牙」这个主判。

## §3 N#3 `run_mode101_test.go` 注释半改写账户 — **成立**

### 机械核「断言半一个字没动」〔我现跑〕

`git show ba5db093 --numstat -- cmd/wisp/run_mode101_test.go`＝**26 / 2**；
删除的那 2 行逐字＝旧绝对句 `// Refused with the grant live IN ITS OWN SESSION, too: the assembly hands`／
`// the bridge a mode, never a grant source (Options.Confirmations nil).`；
新增行里**非注释行数＝0**（`git show | grep '^+' | grep -v '^+++' | grep -vE '^\+[[:space:]]*//' | wc -l` ⇒ `0`）。
⇒ 行为不可能不同，这条不需要靠复跑背书（复跑也做了，见下）。

### 新句是不是「带条件的事实句」

现读 `cmd/wisp/run_mode101_test.go:427-455`：不是绝对句，是三段条件式陈述——
1. **原防护逐字保留**：首行 `Refused with the grant live IN ITS OWN SESSION, too` 仍在，并当场把「它的会话」定义成
   fixture 字面量 `session-before-restart`、逐字写明「which is NOT the session this boot minted」；
   它原本要防的东西在断言半里一字未动（`:461-465` `if !firstErr { t.Fatalf("boot 1 already let the B-tier .env write through …") }`）。
2. **换了原因的条件句**：`This sentence used to read "…never a grant source…" Ticket 224 made that false` ＋
   `the conditional that is actually doing the refusing is this one: tools.GrantSource answers only for the identity it holds …`
   ＋ `Options.Confirmations stays nil either way`。
3. **极性陷阱原地留了路标**：`firstErr` 观察到的是**审批卡超时**、不是权限判定说「不」，
   日后若改期望「必须继续测有没有弹卡而不是有没有被拒」——这正是 `224-v1` §6-b 第 3 条点名的假绿形状。

事实核〔我现跑 13:48〕：`cmd/wisp/run.go:414` 铸、`:529 Grants: grantWrite`、`:653 Grants: grantRead` ⇒ 三处引用为真；
`grep -rn "Confirmations:" --include=*.go internal cmd` ⇒ **1 命中且在 `internal/tools/ticket90_test.go:524`** ⇒
「`Options.Confirmations` 仍是 nil」为真。
⚠ 一处轻过期：第 2 段点名 `TestTicket224ProductionSessionDoesNotSurviveRestart` 作贵方向凭据却没在原地重述它的上限
（上限写在该用例自己文件头）⇒ 不构成假话（它没声称那是 OS 进程级），但「引用不带上限」是本仓登记过的形状，记一笔。

### 两发读数（09-29 派单口径第 3 条）

- 改前／现状＝`--- PASS: TestTicket101SessionGrantDoesNotCrossRestart`：本腿 13:38 基线 3.19s、13:45 M1 下 3.02s、13:44 M-R1 下 3.05s，三发一致绿。
- 全量覆盖＝§6 那发 `ok github.com/CarlosShao/wisp/cmd/wisp 120.780s`〔13:51:48〕。
- `gofumpt -l` 对该文件：空（§6 的 `internal/ cmd/` 全表为空）。

### N#3 判语

**成立**。绝对句已换成带条件的事实句、事实三条均可现跑复认、原防护（B 档 `.env` 在无授权时仍被拒）逐字保留、
断言半机械可证未动。

## §4 N#5 glob 方言定案 — **附条件成立**

### 定案与落点（读台件）

`internal/session/grants.go`：`patternCovers`（`:322-337`）＝方言唯一住所（`*`-含⇒`path.Match`，否则逐字相等，`err`⇒false fail-closed）；
`storeablePattern`（`:280-289`）＝写侧**只**拒「含 `*` 且 `path.Match` 判 `ErrBadPattern`」那一类。
常驻钉＝`internal/session/ticket224_pattern_dialect_test.go`，`grep -c '^func Test'`＝**5**，改前全绿〔13:38〕。

### 两发突变的牙〔我现跑〕

- **M6**（方言退回逐字相等：`path.Match(pattern, p)` → `pattern == p, error(nil)`；15008→15061、标记 1、build rc=0）⇒
  **4 枚指名红**：`WildcardRowCoversAConcretePath`（`:75` 逐字「…the pattern dialect is back to exact equality,
  so every wildcard row on disk grants nothing while the audit view shows a rule」）／
  `DialectIsPathMatchNotRegex`（`:184`×5）／`MalformedWildcard…`（`:290` 控制半）／`PatternCoversIsTheOnlyMatcher`（`:346`×3）。
- **M7**（拆掉「含 `*` 才解释」那道闸）⇒ **2 枚指名红**：`PatternWithoutStarIsNeverReinterpreted`
  （`:224`×3：`note[s].txt` 覆盖了 `notes.txt`、`a?c.txt` 覆盖了 `abc.txt`、`[a-c].txt` 覆盖了 `b.txt`）／
  `PatternCoversIsTheOnlyMatcher`（`:346`）。
⇒ `224-r2` §4 自述的 M6＝4 枚、M7＝2 枚**复跑复现**（M8 本腿未打，见 §7）。

### 攻击 (a)：「看着有规则其实从不匹配」——钉防住了哪一半

**这一半被 `224-v1`→本腿复现，但只防住了「斜杠形态」那半边。** 证据链：

1. 生产 canonical 形态是**反斜杠**〔读台件＋旁证〕：`internal/risk/pathresolver.go:204-213` `lexCanonical = filepath.Abs + filepath.Clean`；
   `Resolve`（`:105-137`）命中已存在文件时取 `GetFinalPathNameByHandle`（`C:\…`），不存在时回落词法串。
   分隔符与大小写折叠只在 `tools/paths.go:294 foldPath`（ToLower＋unifySeparators）里做，**授权匹配链不调它**。
   旁证〔我现跑〕：M4/M-R1 下装配用例打出的卡片路径逐字 `C:\Users\swq\AppData\Local\Temp\…\remembered-by-a-click.txt`；
   M5 下真 store 行的 `ArgsJSON` 也是反斜杠。
2. 钉自己喂的是**斜杠**串〔现读〕：`ticket224_pattern_dialect_test.go:60` 用 `filepath.ToSlash(...)` 造 dir，
   `DialectIsPathMatchNotRegex`／`PatternCoversIsTheOnlyMatcher` 全是 `/work/...` 字面量。
   ⇒ 五枚钉没有一枚拿「生产真形态的通配行 × 生产真形态的路径」比过。
3. 本腿探针实测〔我现跑 13:43:45，台件 `.scratch/wisp/probes/224/v2/dialect_probe_test.go`＋`overlay.json`，
   经 `-overlay` 注入虚拟文件 `internal/session/zz_probe_224v2_test.go`，**跟踪文件零改动**〕：

| 形 | `patternCovers` | `storeablePattern` | 真 `Record`+`Covering`（同库同会话） |
|---|---|---|---|
| A `C:\work\notes\*` vs `C:\work\notes\todo.txt` | **false** | **不拒**（nil） | Record 收下（id=1），Covering **false、grant_id=0** |
| B `C:\work\notes/*` vs 反斜杠子文件 | false | 不拒 | 同上（id=2，false） |
| E `filepath.Join(dir,"*")`（本仓自己那形） | false | 不拒 | 同上（id=3，false） |
| C `dir/*` vs 斜杠子文件 | **true** | 不拒 | — |
| F 逐字反斜杠行 vs 同一串 | true | — | — |
| G `C:\Work\new.txt` vs `c:\work\new.txt` | **false** | — | — |
| H `C:/work/new.txt` vs `C:\work\new.txt` | **false** | — | — |
| 控制 | 斜杠行覆盖斜杠子文件＝true（⇒ 上面不是「匹配器整个死了」） | | |

⇒ **结论：票面 N#5 那句「防止『看着有规则其实从不匹配』」只兑现了一半。**
钉防住的是「匹配器把通配行当字面量」（M6 会红）；**没防住**「一枚合法、写得像规则、`path.Match` 编得过、
但对生产真形态路径永远命不中的通配行」——而且写侧照收、盘上留一行、审计视图里它就是一条规则。
更要紧的是：唯一碰到反斜杠那枚断言（`:353-356`）把这个分歧钉成了**期望行为**（注释逐字「fail-closed」），
⇒ 这一形在盘上被写成设计、不是缺陷，后人不会再碰它。

⛔ **不许把这一条夸成安全事故**，本腿按三行给：
- 现象出现在哪：`internal/session/grants.go` 的 `patternCovers` 与 `storeablePattern` 之间（形态不匹配无人拦）。
- 有没有本机被入侵的证据：无。方向是 **fail-closed**（永远多问一次，不是少问一次）。
- 最坏后果是什么形状：**人手／别处 DAO 写进 `approval_grant` 的通配规则静默不生效**（「写了等于没写」的原形），
  以及操作者按审计视图以为「有这条规则」。生产答复路只会写**具体**路径（Windows 文件名不含 `*`，D7 把射程钉在 Windows），
  所以今天**没有**自动产生的通配行会因此静默。⇒ 未定性为缺陷之前，它是一枚**未写进票面的射程边界**。

### 攻击 (b)：含 `?`／字面量的模式落到哪一支、那一支有没有仪器

- **不含 `*`** 的 `?`／`[ab]`／`]` ⇒ 落**逐字相等**支（`grants.go:326-328` 的星闸）。**有仪器**：
  `PatternWithoutStarIsNeverReinterpreted`（M7 下三形当场红）＋`PatternCoversIsTheOnlyMatcher:342`
  （`{"/work/a?c.txt","/work/abc.txt", false}`）。
- **含 `*` 又含 `?`/`[ab]`** ⇒ 落 `path.Match` 支、`?` 与字符类**被解释**。
  **有仪器**：`DialectIsPathMatchNotRegex` 的三组分歧（`/w?rk/*`、`/work/[ab]*.txt`、`/work/.*`）。
- ⇒ (b) 两支**都有仪器管**，本腿判这一发**没有洞**；只具名一处轻的：混形行里作者本意是字面 `?` 时方言会把它当量词，
  而 `?` 在 Windows 文件名里非法 ⇒ 只有手写行撞得上，无仪器提示（与 (a) 同一族，轻）。

### N#5 判语

**附条件成立**。定案明确、落点单一、五枚常驻钉、两发突变（M6/M7）都有指名红。
条件＝**方言射程与生产 canonical 形态不同形**，导致通配档在 Windows 真实路径上静默不生效、写侧不拒、
且该分歧被一枚钉写成期望 ⇒ 「看着有规则其实从不匹配」只堵了一半。要收这一半必须**换比对口径**
（比对前统一 `foldPath`）或**写侧拒反斜杠通配行＋把射程写进票面**——前者碰 `tools/paths.go` 与 SPEC-02 §3 的语义，属契约级，
本腿不替 owner 挑（`224-r2` §4 末自己也已声明「可原样降格〔待批〕」，本腿附议）。

## §5 N#6 反射钉 setter 半射程 — **成立**

落点＝**新增**一枚 `internal/tools/ticket224_setter_scope_test.go:41 TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter`
（`grep -c '^func Test'`＝1），⛔ 旧那枚 `TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority` 本体一字未改
（⇒ 不会被读成「放宽前人判据」）。三半内容现读为：只许一枚 grant 形状字段且类型必须正是 `GrantSource`（`:44-77`）／
`*Bridge` 上不许有 grant×动作方法、也不许有任何导出 `Set*`（`:79-97`）／接口仍只一枚方法（`:99-112`）。

两发突变〔我现跑〕：

- **M9**＝追加 `func (b *Bridge) SetGrants(g GrantSource) { b.grants = g }`（50699→50779、标记 1、build rc=0）：
  ```
  --- PASS: TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority      ← 旧窄钉不响（224-v1 说的"放过"复现）
  --- FAIL: TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter
      :87 (*Bridge).SetGrants names an act on a grant: ...
      :94 (*Bridge).SetGrants is an exported setter on the enforcement layer: ...
  ```
- **M9b**＝第二枚 **interface 型** grant 字段 `grantsAlt GrantSource`（正是 `Kind()==Func` 那支放过的一形）：
  ```
  --- PASS: TestTicket224GrantSeamCarriesNoCallerSuppliedAuthority      ← 照旧不响
  --- FAIL: TestTicket224BridgeGrantChannelIsOneSeamWithNoSetter
      :54 Bridge carries 2 grant-shaped fields, want exactly 1: [grants(tools.GrantSource, kind interface)
          grantsAlt(tools.GrantSource, kind interface)]. ...
  ```
两次都用 `git cat-file blob HEAD:internal/tools/bridge.go > 同路径` 还原，`wc -c` 均回到 **50699**、标记计数 0。

⇒ **N#6 成立**：`224-v1` §1/§6-c-6 登记的「setter 半只罚 `Kind()==Func`」那条漏判，现在有两形各一枚指名红堵着。

⚠ 射程注记（要台账，不是指控）：`:93-96` 那半条是**词面型**禁令（`*Bridge` 上任何导出 `Set*` 方法一律红），
日后一枚与授权无关的合法 setter（如给桥加 `SetWorkspace`）也会被它打红。
⇒ 「拦名字 vs 拦能力」的口径选择归编排者／owner（本仓先例：词面型负向尺与新功能互斥时立成待人裁定）。
今天 `*Bridge` 无导出 `Set*`（该用例绿），不构成现状冲突。

## §6 门禁读数（本腿自跑，13:49:40 → 13:51:48）〔我现跑〕

台件：`.scratch/wisp/probes/224/v2/gate.sh` 与 `gate.txt`（临时件只建不删）。
起跑前 `tasklist //FI "IMAGENAME eq balldebug.exe"`＝**0 枚**〔13:35〕；
PATH 已带 Sherpa（`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`）。

| 尺 | 读数 | 时刻 |
|---|---|---|
| `go build ./...` | **rc=0** | 13:49:40→13:49:43 |
| `go vet ./internal/agent/approval/ ./internal/session/ ./internal/tools/ ./cmd/wisp/` | **rc=0** | 13:49:43→13:49:44 |
| `./tools/d22scan/d22scan.exe -root .` | **clean，rc=0**（口径：`bans #1-5 internal/=223`、`cmd/=29`、`ban #8 internal/=473`、`cmd/=65` 全是**被扫文件数、不是违规数**） | 13:49:44→13:49:45 |
| `go test ./internal/agent/approval/ ./internal/session/ ./internal/tools/ ./cmd/wisp/ -count=1` | **四包全 ok**：0.342s／0.805s／12.907s／**120.780s**，`TEST_RC=0` | 13:49:45→13:51:48 |
| `gofumpt -l internal/ cmd/` | **空**（rc=0）。⚠ 第一发在 gate.sh 里是 `command not found`（rc=127＝仪器没跑），本腿补 GOPATH/bin 后重跑；**正控**：同一把尺对 `.scratch/wisp/probes/224/v2/posctrl.go` 故意写歪的文件**列得出来** ⇒ 尺是活的 | 13:51:56 |

四枚包口径：本腿**未**把 `internal/audio` 放进自己的门禁（票 241 地界）。
⚠ 但 `go test ./cmd/wisp/` 会把那枚包链进测试二进制 ⇒ `120.780s` 那枚只当「cmd/wisp 用例没被 224 的突变洗掉」用，
**不当「全仓此刻绿」用**（与 `224-v1` §7 末同口径）。与编排者 13:30-13:32:35 那发（五包、cmd/wisp 129.061s）**不互替**。

工作树终态：见 §9 末（起手名册＝空，终态＝空）。

## §7 我攻不动的地方

- ⚛⚛ **真·两枚 OS 进程的授权路**：本腿没有接这根线的时间预算（`buildWispForTest` 编一次 exe 加两发 `wisp run` 冷启动
  ≈ 每发几十秒，且要 `//go:build windows` 与进程 env 布置）。⇒ 判语上限见 §2：**AC#3「进程重启」那一半今天仍无仪器**，
  但它是**欠做**不是**硬墙**（原语就在同一枚包 `cmd/wisp`）。这条同时打在交件那枚新用例、`224-v1` 那枚旧钉与票面⑩ 的措辞上。
- ⚛ **「台架不存在」与「没接线」的区别我没法用读数消掉**：我没有造第二枚进程，所以对「第二次装配铸新 key」这一事实之外的一切，
  本腿只能给文字判语（§2 那条更正属〔现读现跑的原语清单〕，不是那条链路的执行证据）。
- ⚛ **M8（解释器换成 regexp）本腿未打**：`224-r2` 自述它红 3 枚。本腿已用 M6＋M7 两发证明这五枚钉不是装饰，
  且 M8 的落点（把 `path.Match` 换成 `regexp`）需要引入新 import 并重排 `grants.go` 的导入块，突变面比本票需要的更大；
  方言「是哪一枚」已由 `DialectIsPathMatchNotRegex` 的分歧表钉住 ⇒ 不补判语。
- ⛔ **面板（WebView2）侧**：`PanelAPI` 会不会长出会话档入口，只能靠 `internal/panel/l2_grant_boundary_test.go`
  （**冻结件，只读不跑**）的词表判断，本腿不改它、不探它的边界。
- ⛔ **`internal/audio` 在飞（票 241）**：零碰、零跑、零归因。若 224 的某枚读数与那枚包耦合，本腿判不动，按此条具名而不归因。
- ⚛ **随机源故障分支**：`session.Mint()` 的 `crypto/rand` 失败路径没有注入口（不改产码就打不红），
  「铸不出来时只关掉这一档」那一支本腿只能读、不能打。
- ⚛ **`GRANT-DROPPED` 的「卡片主题在答复前已离开队列」那一支**（`gate.go:674-677` 的 `tool == ""` 分支）：
  要造它得让 queue item 在 `sessionSubject` 与 `allowScoped` 之间消失，本腿没有不动产码的注入口。
  ⇒ `224-v1` §7 同一支，本腿复看仍未攻。
- ⚛ **本腿没打的一支（具名，不是"判不动"）**：`Covering` 里 `g.Tool != tool` 那枚工具项逐字比对若被改成大小写不敏感
  （放宽方向），会不会有钉红——本腿没测。它与 §4 的大小写射程（探针 G 行）同族，留给续腿或 N#5 的收口一起做。
- ⚛ **`internal/observe/thresholds.go`／golden／`docs/SLO.md`／`allowlist.txt` 零碰**，故本表对这些无任何判语。

## §8 留给编排者的台账动作（本腿只登记、不代裁、不翻框）

1. **票面 `- [ ]` 一枚未碰**。建议（非本腿权限）：N#1／N#3／N#6 可翻；N#2 附条件入账；N#5 附条件入账。
   AC#2 的「答复⇒落行」那一格经 N#1 已补上判据，但 AC#3 的 OS 进程那一半仍缺（见 7）⇒ **AC#3 不因本表升格**。
2. **新欠账建议立 `A##`/`Q##`（本腿判属契约级，要 owner 拍）：通配模式的分隔符形态**。
   事实：生产 canonical＝反斜杠（`internal/risk/pathresolver.go:204-213`），`patternCovers` 只解释斜杠形
   （`internal/session/grants.go:322-337`），且 `storeablePattern` 不拒反斜杠通配行〔探针：A/B/E 三形写侧全收、读侧全不命〕。
   两条收法：ⓐ 比对前统一走 `tools/paths.go:294 foldPath` 那一类（**碰 SPEC-02 §3 `pattern` 的语义**）；
   ⓑ 写侧拒「含 `*` 且非斜杠形」的行＋把「通配档只保证斜杠形态」写进票面与注释（**不碰契约、只把射程说清**）。
   ⛔ 本腿不挑（这不是实现细节）。
3. **同族轻量账（不必 owner，但要立号）**：`patternCovers` 不分折大小写（探针 G：`C:\Work\new.txt` 不覆盖 `c:\work\new.txt`），
   而未存在的目标走词法兜底、大小写由拼写决定 ⇒ 同一目标可能有多种 canonical 串；今天无仪器（既没钉成设计、也没钉成缺陷）。
4. **注释坐标过期两处**（`cmd/wisp/ticket224_assembly_test.go:23`、`:24-25`、`:388`）：
   实际三处在 `internal/session/grants_test.go:20/:101/:147`（不是 `:20/:94/:134`）；
   `internal/tools/grant_test.go:6-9` 现在指**文件名**、不再指用例名 ⇒ 「four comments point AT THIS NAME」已降为三处。
   指向的对象真实存在，所以这不是幻影复发，是轻形。
5. **更正一条前人引错的路标（台账只追加不删）**：`224-v1` §7 与本票派单都写「`run_mode101_test.go:145-147` 那句
   `a new process surface`」——〔我现跑 `grep -n`〕该句真身在 **`:73-74`**（`t101host` 的类型文档），`:143-158` 是 `start` 的函数体。
   交件 `ticket224_assembly_test.go:34-36` 改用的 `:72-76`＋`:143-158` 才是对的坐标。⇒ 判语不变（仍是同一测试进程），
   但那枚引用以后照交件写。
6. **更正一句过强的上限（同一处文字）**：`ticket224_assembly_test.go:47-49` 的「that harness does not exist yet」应改成
   「真两枚进程的台架原语已在同包（`buildWispForTest`＋`exec.Command`），只是没人为授权路接上」。
   理由与坐标见 §2。
7. **N#6 的词面型禁令**：`ticket224_setter_scope_test.go:93-96` 禁 `*Bridge` 上任何导出 `Set*`；
   日后与授权无关的合法 setter 会被打红 ⇒ 是否收窄成能力型，待人裁（现状无冲突）。
8. **死腿自述的尾部两节仍是「（待填）」**：`docs/evidence/s1/224-r2-progress.md` §6/§7。
   本腿**未代填**（§6 是自己 13:49-13:52 那一发、§7 是本腿自己的攻不动清单）。
   ⇒ 那份文件要么标注〔§6/§7 未交，由 224-v2 独立取数〕，要么由编排者就地补一行指针；不许读成「两枚表都有尾部自证」。
9. **本表用的探针台件**：`.scratch/wisp/probes/224/v2/`（`dialect_probe_test.go`＋`overlay.json` 可复跑；
   `gate.sh`／`gate.txt`／`m1-cmdwisp.txt`／`m7-session.txt`／`m9-tools.txt`／`m9b-tools.txt`／`posctrl.go` 是读数留档）。
   仓内跟踪文件零改动。

## §9 突变腿的自证（硬约束 6）

| 突变 | 打的是哪一枚产码 | 突变态 | 还原方式 | 还原后 `wc -c` | 标记计数 | 复绿 |
|---|---|---|---|---|---|---|
| **M4** 摘掉 `gate.go:678` 的 `Record` | `internal/agent/approval/gate.go` | 31010→31082 | `git cat-file blob HEAD:…>同路径` | **31010** | 0 | `ok approval 0.325s` |
| **M-R1** `Record` 写错 session id | `internal/session/grants.go` | 15008→15073 | 同上 | **15008** | 0 | 后续 M6 同文件复跑绿 |
| **M1** `Mint()` 改派生串 | `internal/session/session.go` | 4174→4228 | 同上 | **4174** | 0 | §6 全绿 |
| **M6** 方言退回逐字相等 | `internal/session/grants.go` | 15008→15061 | 同上 | **15008** | 0 | §6 全绿 |
| **M7** 拆掉「含 `*` 才解释」闸 | `internal/session/grants.go` | 突变态字节数**本腿未记**（只记了标记＝1、build rc=0） | 同上 | **15008** | 0 | §6 全绿 |
| **M5** 桥不写 `tool_call.grant_id` | `internal/tools/bridge.go` | 50699→50749 | 同上 | **50699** | 0 | §6 全绿 |
| **M9** 加 `(*Bridge).SetGrants` | `internal/tools/bridge.go` | 50699→50779 | 同上 | **50699** | 0 | §6 全绿 |
| **M9b** 加第二枚 interface 型 grant 字段 | `internal/tools/bridge.go` | 突变态字节数**本腿未记**（标记＝1、build rc=0） | 同上 | **50699** | 0 | §6 全绿 |

- **每一发都先证「突变真落地」再跑**：`grep -c` 标记命中数、旧串计数归零、`go build` rc=0 三件齐才算数；
  没落地的发一律作废不采（本腿没有作废发，八发全部落地）。
- ⛔ **没有任何跟踪文件被写成 0 字节**：每发还原后立刻 `wc -c` 对回基线（上表列全）。
- 探针一律走 `-overlay` 注入虚拟文件（`internal/session/zz_probe_224v2_test.go`），**不动跟踪文件**；
  探针读的是内存里的真值表与真 SQLite（`t.TempDir()` 下新建库），没有任何一处去读那枚被替换的路径
  ⇒ 不踩本仓「`-overlay` 对读盘的测试结构性失明」那枚坑。
- 全程未碰：`internal/audio/**`（零跑零读零归因）、三枚冻结件（`internal/panel/tokens_fourway_test.go`／
  `internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）、
  `internal/observe/thresholds.go`、golden、`docs/SLO.md`、`allowlist.txt`、`docs/PLAN.md`、`docs/specs/**`、
  `frontend/**`、`design/**`；**零阈值变更、零断言放宽**（本腿唯一动的产码全是用完即还原的突变）。
- **Git 纪律**：只 commit、未 push；每枚 commit 都带显式 pathspec（无 `add -A`/`.`）；无 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
  别人那几枚脏文件（`.gitignore`、`design/**`、`probes/**`）零碰、零还原、零提交。
- **终态自证（一把尺，口径＝起手名册，不是「必须为空」）**：
  `git status --porcelain -- internal cmd` 起手〔13:34:27〕＝**空**，本腿八发还原后各自复验均＝**空**
  ⇒ 与起手名册相同（`-- internal cmd` 这两枚包内本就没有别人的在飞件；包外的别人脏件本腿一律不列、不动）。
