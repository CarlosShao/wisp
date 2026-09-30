# 224 — 会话授权「判据返工」对抗验收表 · 腿 `224-v2`

- 被验收对象：六枚 commit `b7ca76b5`（骨架）→ `2a0175ba`（N#1）→ `1c601fab`（N#1＋N#2）→ `ba5db093`（N#3）→ `e8ed5fef`（N#5）→ `43d9096f`（N#6）
- 工单：`.scratch/wisp/issues/224-session-scoped-grant-has-zero-executors-and-no-session-identity.md` §「续做」
- 判语来源（欠账清单）：`docs/evidence/s1/224-session-grant-v1.md`
- 死腿自述：`docs/evidence/s1/224-r2-progress.md`（§1-§5 正文写完；**§6 门禁四数与 §7 判不动仍是「（待填）」＝本腿不代填、不当它已交**）
- 本腿＝**非实现者验收腿 `224-v2`**；产码腿 `224-r2` 已死于 150 轮上限
- 本表作者锚点：起手 `git log -1`＝`ba9e3575`（编排者 A468）
- 骨架落盘时刻 `2026-09-30 13:3x +08`
- ⛔ 票面 `- [ ]`（AC#1-AC#5 与 N#1-N#6）**一枚不碰**，勾框归编排者
- ⛔ `internal/audio/**`＝票 241 地界：零碰、零跑、零归因
- ⛔ `frontend/**`／`design/**`：两层禁（不读、结论不引）
- 射程：只碰 `internal/agent/approval`／`internal/session`／`internal/tools`／`cmd/wisp`

## 口径声明

- 〔我现跑〕＝本腿本机此刻执行并贴回读数；〔读台件〕＝只读代码／台账／commit，未执行。
- 用例枚数口径：`grep -c '^func Test' <file>`（顶层 Test 函数数），不是断言数、不是子用例数。
- 「被扫文件数」≠「违规数」。
- 行号一律现取。
- 编排者已给的起跑基线（**非我跑的，抄自派单**）：`go build ./...` 13:30:19 rc=0；
  `go test ./internal/agent/approval/ ./internal/session/ ./internal/tools/ ./internal/audio/ ./cmd/wisp/ -count=1`
  13:30-13:32:35 五包全 ok（0.367s／1.131s／13.238s／15.574s／129.061s）；`gofumpt -l internal/ cmd/` 空；
  `./tools/d22scan/d22scan.exe -root .` clean rc=0。

---

## §0 逐格 1:1 表（续做五格 → 本表五节）

| 续做格 | 票面要求（要点） | 本表节 | 判语 |
|---|---|---|---|
| N#1 | 「答复⇒落行」常驻判据（M4 必须有指名用例红）＋真 `tool_call.grant_id` 列有用例 | §1 | （待填） |
| N#2 | 真·跨进程重启用例＋处置三处指着不存在测试的注释；判定上限逐字核 | §2 | （待填） |
| N#3 | `run_mode101_test.go` 注释半⇒带条件的事实句、逐字保留原防护 | §3 | （待填） |
| N#5 | glob 方言定案＋常驻钉防「看着有规则其实从不匹配」 | §4 | （待填） |
| N#6 | 反射钉 setter 半射程 | §5 | （待填） |
| （N#4 编排者已裁，不在本腿射程） | — | — | — |

## §1 N#1 「答复 ⇒ 落行」的常驻判据 — （待填）

### M4 复跑（正向：摘掉 `Record`）— （待填）

### 反向一发：`Record` 放回但写错的行 — （待填）

### 真 `tool_call.grant_id` 列的覆盖 — （待填）

## §2 N#2 幻影注释与真·跨进程用例 — （待填）

### 三处注释的处置 — （待填）

### 判定上限逐字核（`run_mode101_test.go` 的 `a new process surface`）— （待填）

## §3 N#3 注释半改写账户 — （待填）

## §4 N#5 glob 方言钉 — （待填）

### (a) 看着有规则其实从不匹配 — （待填）

### (b) 含 `?`／字面量的模式落到哪一支 — （待填）

## §5 N#6 反射钉 setter 半 — （待填）

## §6 门禁读数（本腿自跑，非抄派单）— （待填）

## §7 我攻不动的地方 — （待填）

## §8 留给编排者的台账动作 — （待填）

## §9 突变腿的自证（硬约束 6／终态自证）— （待填）
