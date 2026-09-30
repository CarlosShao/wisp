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

（待填）


## §2 N#2 幻影用例：写用例还是改注释

（待填）

## §3 N#3 `run_mode101_test.go` 注释半的改写账户

（待填）

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
