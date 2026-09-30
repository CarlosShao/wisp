# 224-r2 — 会话授权「判据缺口」返工腿进度表（N#1／N#2／N#3／N#5）

- 工单：`.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md` §「续做」
- 判语来源：`docs/evidence/s1/224-session-grant-v1.md`（非实现者验收腿 `224-v1`）
- 本腿＝**产码腿 `224-r2`**；⛔ 票面 `- [ ]`（AC#1-AC#5 与 N#1-N#6）**一枚不碰**，勾框归编排者
- 射程：只碰 `internal/agent/approval`／`internal/session`／`internal/tools`／`cmd/wisp`
- ⛔ `internal/audio`＝票 241 地界：零碰、零跑、零归因
- ⛔ `frontend/**`／`design/**`：两层禁（不读、结论不引）
- 骨架落盘时刻 `2026-09-30 12:5x +08`

## 口径声明

- 〔我现跑〕＝本腿本机此刻执行并贴回读数；〔读台件〕＝只读代码／台账／commit，未执行。
- 用例枚数口径：`grep -c '^func Test' <file>`（顶层 Test 函数数），不是断言数、不是子用例数。
- 「被扫文件数」≠「违规数」（`tools/d22scan` 的 `internal/=NNN` 是前者）。
- 本腿锚点：起手 `git log -1`＝`3c65467e`（编排者 A467 收 224-v1 那一枚）。
- ⚠ 行号一律现取：票面与 `224-v1` 表的行号在派单里已声明「别信这里的行号」，本表每处 `file:line` 均为本腿现跑 `grep -n` 所得。

## §0 逐格 1:1 表

| 续做格 | 票面要求（要点） | 本表节 | 状态 |
|---|---|---|---|
| N#1 | 「答复⇒落行」常驻判据＋真 `tool_call.grant_id` 列有用例 | §1 | 进行中 |
| N#2 | 真·跨进程重启用例＋处置三处指着不存在测试的注释 | §2 | 待填 |
| N#3 | `run_mode101_test.go` 注释半改写账户 | §3 | 待填 |
| N#4 | 会话时长定案（编排者已裁） | — | 本腿不动 |
| N#5 | glob 方言定案＋常驻钉 | §4 | 待填 |
| N#6 | 反射钉 setter 半射程 | §5 | 待填 |

## §1 N#1 「答复 ⇒ 落行」的常驻判据

### 判据形状

`internal/agent/approval/ticket224_reply_grant_test.go`〔本腿新建，`grep -c '^func Test'`＝4〕，
四枚各管一头，⛔ 不并成一枚「持久化」用例（AC#2 的分裂要求同样适用于答复路）：

| 用例 | 钉住的是哪一发 | 会不会被 M4 打红 |
|---|---|---|
| `TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed` | 答复 ⇒ 卡片自己印的每一条路径各一次 `Record`，且审计行的 `grant_id` ＝ 记录器真返回的 id | 会（两半都红：调用数 0、`grant_id=0`） |
| `TestTicket224ForgedSessionAnswerRecordsNothing` | 「答复」是落行的唯一入口：四枚伪造/失效 nonce 全部拒绝且**零行**，随后真答复仍落一行 | 会（真答复那半红） |
| `TestTicket224SessionAnswerWithoutARecorderSaysTheScopeWasDropped` | 没有记录器的宿主：放行照旧、审计说 `GRANT-DROPPED`，⛔ 不许出现 `GRANT-RECORDED`／`grant_id=` | 不红（这一支 M4 恰好仍走 nil 分支）＝反向钉，防「日后把谎话写进没有记录器那一支」 |
| `TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow` | 两条路径一条写失败：调用仍放行、失败的不上账、成功的带 id | 会（三行都红） |

判据口径：这三枚红的都是**指名用例**（`--- FAIL: TestTicket224...` 逐字打名），不是包级 `[build failed]`，
也不是"整包超时"那种不可归因的红。

### 突变 M4 前后读数

M4＝验收腿 `224-v1` §2 那一发原样复现：`internal/agent/approval/gate.go:678` 的
`id, err := g.grants.Record(ctx, tool, p)` 摘掉，改成 `id, err := int64(0), error(nil)`（答复照旧放行、
盘上一行不落、`GRANT-RECORDED` 照打）。

**改前（今天全仓无数枚看得见这一发）〔我现跑，起跑基线之后、加用例之前〕**

```
ok  github.com/CarlosShao/wisp/internal/agent/approval   0.329s
ok  github.com/CarlosShao/wisp/internal/session          0.331s
ok  github.com/CarlosShao/wisp/internal/tools           12.193s
ok  github.com/CarlosShao/wisp/cmd/wisp                115.487s   ← 四包一枚不红，与 224-v1 的 0.318/0.332/12.187/156.567 同形
```

**加用例后、同一发 M4 仍在盘上〔我现跑〕**

`internal/agent/approval/ticket224_reply_grant_test.go`（本腿新建，50 枚断言级判据在此文件）：

```
--- FAIL: TestTicket224AllowSessionRecordsEveryPathTheAnsweredCardNamed
--- FAIL: TestTicket224ForgedSessionAnswerRecordsNothing
--- FAIL: TestTicket224RecordingFailureReleasesTheCallButClaimsNoRow
FAIL github.com/CarlosShao/wisp/internal/agent/approval 0.028s
```

⇒ 判据达成：**同样这一发 M4，现在有指名用例红**（三枚），红的第一行逐字为
`recorder saw 0 Record calls, want 2 (the card printed 2 paths, calls=[])`。

**还原后（`git cat-file blob HEAD:internal/agent/approval/gate.go > 同路径`，`wc -c`＝31010→31010，
`grep -c MUTATION-M4-TEMP`＝0）〔我现跑〕**

```
ok  github.com/CarlosShao/wisp/internal/agent/approval   0.325s
go vet ./internal/agent/approval/  rc=0
gofumpt -l internal/agent/approval/ticket224_reply_grant_test.go  （空＝已格式化）
```

### 一处由仪器当场抓出的本腿自身缺陷（先记，不藏）

第一版 fixture 写成 `Grants: f.rec`，而 `f.rec` 是值为 nil 的 `*recGrantRecorder`
⇒ **非 nil 接口装 nil 指针**，`Gate.allowSession` 的 `g.grants == nil` 那一支根本不进。
M4 那发把它照了出来（nil-recorder 用例在突变下打出 `GRANT-RECORDED grant_id=0`，
而不是应有的 `GRANT-DROPPED`）。修法＝只在真有 recorder 时才往接口字段赋值，
理由就写在 `newGrantFixture` 的注释里，与 `cmd/wisp/run.go:511` 那句「typed-nil guard is load-bearing」同源。

### 真 `tool_call.grant_id` 列的覆盖

两处，各管一头：

| 落点 | 钉住什么 | 现跑读数 |
|---|---|---|
| `internal/tools/ticket224_grantid_column_test.go`（1 枚用例） | 桥写的值经**真 `memory.Store`** 落进列、再由第二个句柄按 (task, corr) 读回；断言是**连接**而不是指针：`grant_id` 必须指向一枚真实存在的 `approval_grant` 行（同 tool 同 pattern、未撤销）；对照组＝同库同一次启动里未授权那一发必须 `NULL` | 正常码 `--- PASS`（0.07s）；突变 M5（`bridge.go:1116` 的 `if grantID != 0 {` 改 `if false {`）⇒ **本枚红**，且 `grant_test.go:445` 那枚桩侧同红（说明这一发的牙两半都有）；还原后 `internal/tools` 全绿 12.263s，`bridge.go` 50699→50699 |
| `cmd/wisp/ticket224_assembly_test.go` 第二枚用例 | 装配根铸的 ledger + 装配根的桥 + 真 store：命中授权的那发 `decision=allow_session_grant` 且 `grant_id`＝真行 id，未命中那发 `cards=1`／`grant_id=NULL` | `--- PASS`（2.98s） |

为什么这比 224-v1 §2-③ 说的"零覆盖"多了一寸：`internal/memory/dao_test.go:355` 单独往返过这一列，
`grant_test.go:445` 单独断过桥递出的指针，**两半从来没接起来**——桥到 DAO 之间丢掉这一列，全仓没有一枚用例会红。

### 装配侧两枚（N#1 的写／读两格）

`cmd/wisp/ticket224_assembly_test.go`〔本腿新建，`grep -c '^func Test'`＝3〕：

- `TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun`＝**写**格。答复走 `runSpec.reply`
  那枚 CLI 注入缝（AGENTS.md §1.3 允许的形状，不是 mock 顶替真件），断言盘上的行的
  `pattern` **等于卡片自己印的路径**（不是本测试自己写的串），`scope/session`、`tool` 取自卡片、
  `expires_at > created_at`、`revoked_at IS NULL`，另断审计行的 `grant_id`＝真行 id、
  且**第二个句柄按铸造 id 读得回来**。M4 下这一枚**红**（`0 approval_grant rows ... want 1`，逐字带铸造 id 与卡片印的路径）。
- `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking`＝**读**格（见上一节）。
- `TestTicket224ProductionSessionDoesNotSurviveRestart`＝N#2，见 §2。

⚠ 观测口径：两枚都读 `rt.windowCount()`（**弹卡计数**），⛔ 不读"有没有被拒"——
224-v1 §6-b 第 3 条点名的正是那一形（L1 窗口超时＝放行，"被拒"是那道门给的，不是权限判定给的）。


## §2 N#2 幻影用例：写用例还是改注释

**选了"写用例"，并且同时改了注释。** 理由一句话：那三处注释指的东西（跨启动的铸造者、
测试字面量锁的"贵方向"、票面⑩ 的①②两件控制）**本来就是判据缺的那一格**，把注释删掉只是把
缺口换个写法留在盘上；写出来才让 AC#3 那一格有可跑的东西。而注释也一律改了——因为用例存在之后，
原话"cross-process／needs two boots"仍然**过量声称**，必须带上判定上限。

- 用例：`cmd/wisp/ticket224_assembly_test.go` 的 `TestTicket224ProductionSessionDoesNotSurviveRestart`
  （名字逐字等于三处注释指的那枚）。台架＝`run_mode101_test.go` 的 `t101host`（同目录第二次 `runTextTask`）。
- 票面⑩ 的三件控制，缺一不算判据：
  ① 两次装配铸出两枚不同 id（`id2 == id1` 直接 `t.Fatalf`）；
  ② 用 boot 1 的生产 id 写、再用 boot 2 的生产 id 查 ⇒ **0 行**；
  ③ 同一发 `fs.write` 走**真 `rt.bridge.Execute`**（经 `t101call`）⇒ 第二次**弹卡 ≥ 1**，
     并另断那一行的 `decision != allow_session_grant` 且 `grant_id IS NULL`。
  ⚠ ③ 的极性与 224-v1 §6-b 第 3 条一致：断的是**有没有弹卡**，不是"有没有被拒"（L1 窗口到点＝执行）。
- 反控实测〔我现跑，M1＝`internal/session/session.go:99` 的 `Mint()` 改成派生十六进制串，
  即 224-v1 §2 表那一发〕：

  ```
  --- FAIL: TestTicket224ProductionSessionDoesNotSurviveRestart (0.93s)
      ticket224_assembly_test.go:481: the second assembly minted boot 1's identity
      "sess_776973702d7069642d64657269766564": 「本次会话内」 would then be a permanent pass, ...
  --- FAIL: TestTicket224MintIsRandomAndValid / TestTicket224CoveringDoesNotCrossToANewSession  (internal/session)
  ```

  还原：`git cat-file blob HEAD:internal/session/session.go > 同路径`，`wc -c`＝4174→4174，突变标记计数 0。
- 判定上限（**逐字写在用例文件头部**，不在本表代它声称）：`run_mode101_test.go:72-76` 那句
  `"Restart" in every case below is a second runTextTask over the same dir: a new process surface`
  ＋ `:143-158` 的 `start` ⇒ 第二次装配**仍是同一个测试进程**。
  ⛔ 因此今天全仓没有任何仪器能区分"重启"的两种含义，这条上限同时打在交件那枚旧钉、
  本腿这枚新用例与票面⑩ 的措辞上；要真两枚 OS 进程得另建台架，本票没有那条时间预算。
- 三处注释的处置（票面锚点 `:20`／`:94`／`:134` 的三处，本腿现跑 `grep -n` 落在
  `internal/session/grants_test.go:23`／`:99`／`:145`；外加第四处 `internal/tools/grant_test.go:6`）：
  逐枚改成**"该用例现在存在＋它声称得了的那一半"**，⛔ 不留"指着不存在的东西"的指针，也不留过量声称。
  `git diff --numstat` 删除列＝5／1，与替换行数相等（没有吞行）。

## §3 N#3 `run_mode101_test.go` 注释半的改写账户

**改前**（本腿现跑 `grep -n` 于改动前，逐字两行）：

```
427		// Refused with the grant live IN ITS OWN SESSION, too: the assembly hands
428		// the bridge a mode, never a grant source (Options.Confirmations nil).
```

⇒ 后半句今天是假话：`cmd/wisp/run.go:653` 现在写着 `Grants: grantRead`（写侧 `:529`，铸造点 `:414`）。

**改后**（同处，`cmd/wisp/run_mode101_test.go:427-452`，带条件的事实句）——三段各有分工：

1. 保留原句要防的东西，逐字仍在第一行：`Refused with the grant live IN ITS OWN SESSION, too`，
   并把"它的会话"当场定义为**fixture 的字面量**，不是本次启动铸出的那枚；
2. 说明为什么这一发仍然被拒（原因换了）：`GrantSource` 只回答它自己持有的身份，
   行上的 key 是一枚生产铸造永远产不出的串，并点名两方向锁的两枚用例
   （`internal/session/grants_test.go` 的 `TestTicket224TestLiteralsAreOutsideTheMintedShape`
   ＋本腿新写的 `TestTicket224ProductionSessionDoesNotSurviveRestart`）；`Options.Confirmations` 仍是 nil；
3. ⛔ 不留绝对句、且把 224-v1 §6-b 第 3 条的极性陷阱写进注释：`firstErr` 观察到的是**审批卡超时**，
   不是权限判定说"不"；日后若改成"带活授权也要被拒"，必须继续测**有没有弹卡**而不是"有没有被拒"。

**断言半（票面锚点 `:434-441`，注释块插进去之后现在落在 `:458-465`）一个字没动**〔我现跑尺〕：`git diff --numstat` ＝ `26 / 2`，
且 `git diff -U0` 里剥掉注释行后**没有任何非注释变更行** ⇒ 行为不可能不同。

**两发读数（09-29 派单口径第 3 条）**：

- 改前＝该枚在起跑基线与本腿 M4 复跑里都是绿的（`go test ./cmd/wisp/ -count=1` ok 115.487s 覆盖它，
  编排者给的起跑读数 ok 125.521s 同样覆盖）；
- 改后＝`go test ./cmd/wisp/ -count=1 -run TestTicket101SessionGrantDoesNotCrossRestart -v`
  ⇒ `--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (3.17s)`／`ok cmd/wisp 3.218s`〔我现跑〕。
- `gofumpt -l` 对该文件：空。

## §4 N#5 glob 方言定案

### 选了哪一种、为什么

（待填）

### 那枚常驻钉防住的是哪一种失败

（待填）

## §5 N#6 的处置（收／不收）

（待填）

## §6 门禁四数与工作树

（待填）

## §7 判不动的地方

（待填）
