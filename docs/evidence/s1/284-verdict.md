# 284 裁决：`Gate.Replay` 注释三要素（只读非实现者，gate.Replay 面）

**票**：`.scratch/wisp/issues/284-gate-replays-comment-omits-the-zero-production-caller-claim.md`（唯一一枚，无 `-done` 变体；本裁决时票面三格**仍全部未勾**）。
**角色**：只读非实现者裁决腿（未改一行产码、未改票、未 commit、未 push；**唯一写面＝本文件**）。
**票面锚**：`6589a0625`（任务给定 HEAD）。**实测锚**：`0fa1527f9`（2026-11-05 18:37 +08），18:39 收官复跑时 HEAD 已再漂到 `0af3b17aa`（并发收尾，非本票面变化）。
**漂移无关性（现跑钉住）**：`git rev-parse HEAD:internal/agent/approval/gate.go` ＝ `git rev-parse 556f0483c:internal/agent/approval/gate.go` ＝ `8fdbc8be74c9d2262b645d6563e57aba9290159d`（18:39 在 `0af3b17aa` 复核仍是这一枚）；`git diff --stat 556f0483c HEAD -- internal/agent/approval/` **空** ⇒ 实现提交之后该包零字节变化，票面 HEAD 与实测 HEAD 上的注释是**同一份 blob**。
**实现提交**：`556f0483c`（2026-11-05 16:37:34 +0800，`1 file changed, 9 insertions(+)`，只动 `internal/agent/approval/gate.go`）。

---

## 1. AC#1 —— **PASS**（三要素齐，两条禁忌均未触）

改后注释块 `gate.go:872-884`（逐字，2026-11-05 18:37 现读）：

```
// Replay re-displays a refused or expired L2 request under a fresh
// correlation id (C18 一键重放). It requires the task to be admitted again
// (D47 still applies to a replay), and it is NOT an answer: the new item must
// still be approved with a newly minted grant.
//
// Today it has no production caller, and that is a state to record, not a
// defect to tidy away: do not delete it as dead code. It is the queue half of
// C18's 一键重放, already implemented here - Queue.replay mints the fresh
// correlation id and links the new item back through replayOf. The missing half
// is the answer-side verb, and neither NativeAPI nor PanelAPI declares one, so
// no card, no panel and no console line can reach this method yet. Adding that
// half is a wiring decision with an owner, which is why the caller side reads
// the way it does rather than because somebody forgot to call it.
```

三要素对账（票面 AC#1 要求的〔零产码调用者／非死代码勿删／替 C18 一键重放留的＋缺的一半是接线任务且有 owner〕）：

| # | 要素 | 落在哪一句 | 判 |
|---|---|---|---|
| ① | 今天零产码调用者 | `Today it has no production caller` | 在 |
| ② | 非死代码、勿删 | `not a defect to tidy away: do not delete it as dead code` | 在 |
| ③ | 替 C18 一键重放留的（缺的答侧动词是**有 owner 的接线任务**） | `It is the queue half of C18's 一键重放 … The missing half is the answer-side verb … Adding that half is a wiring decision with an owner` | 在 |

**要素③不是空话，三条现跑都成立**：

- 零产码调用者：`git grep -nE "[A-Za-z0-9_]\.Replay\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'` ⇒ **rc=1（空）**；不剥测试同尺只多出 `internal/agent/approval/queue_test.go:283`（`g.Replay(...)`），与"唯一调用点是测试"一致。
- `Queue.replay` 确会铸新关联 id 并回链：`queue.go:609` 起 `func (q *Queue) replay`，`queue.go:623` `fresh.replayOf = it.Corr`，`queue.go:625` `d.CorrelationID = corr + "#replay"` ⇒ 注释说的"queue half 已实现、铸新 corr、经 replayOf 回链"**逐项属实**。
- "两侧接口面都不声明它"：`ui.go` `NativeAPI` ＝ `Allow`／`AllowSession`／`Reject`，`PanelAPI` ＝ `Reject`／`Head`／`View`，**均无 Replay**；控制台动词 switch（`cmd/wisp/approval_reply.go`）＝ `yes/session/no/veto/panel-no/panel-yes/always/head/view` ＋ `help/quit`，**无 replay** ⇒ "no card, no panel and no console line can reach this method yet" 属实。

**票面两条禁忌**：

| 禁忌 | 现跑 | 判 |
|---|---|---|
| ⛔ 注释里嵌行号 | `sed -n '872,884p' gate.go \| grep -nE '\.go:[0-9]\|:[0-9]{2,3}'` ⇒ **rc=1（空）**；指路只用具名符号 `Queue.replay`／`replayOf`／`NativeAPI`／`PanelAPI` | 未触 |
| ⛔ 写会被自己扫的现在时计数（"0 次"这类） | 全文无**计数式**表述；句首 `Today` 把它框成状态记录。块内出现的数字只有契约号 `L2`／`C18`／`D47`（4 行含数字），**不是可被扫的计数** | 未触 |

> 说明：AC#3 那条禁令射的是"把被数的 token 写进注释"（写进去那一刻就变假话），不是"注释里不许出现任何数字"；契约号不属该射程。

---

## 2. AC#2 —— **FAIL（过程偏差，如实记账）**

票面 AC#2 逐字：⛔ **不许为这一行单独开一条腿去改 `gate.go`**；落地时机＝下一枚**本来就动** `gate.go` 的腿。

事实（三发可核，不替谁圆）：

1. 落地方式就是**一条专为票 284 单开的腿**：`556f0483c`，提交信息自带自认——`284 gate.Replay 注释补零产码调用者三要素(死腿遗产代提,AC#1实质;注:AC#2本应随批勿单腿=登记纪律偏差)`。
2. 该腿的交付件 `docs/evidence/s1/284-gate-replay-zero-caller-comment.md` §5 也自认："从字面上看它**正是 `AC#2` 想避免的那条腿**"，并给出"AC#2 保持不勾、票不 `-done`"的处置建议。
3. 代价可量化：为 9 行注释花掉一次 `gate.go` 的编译面（build 全量＋`approval` 包测试＋vet＋d22 全量）。
4. 无"随批"掩护可辩：`git diff --stat 556f0483c^ 556f0483c` 只含 `gate.go`，同一提交没有别的功能改动 ⇒ 不是"别人本来就要动 `gate.go`，我捎带一行"。

⇒ 判定：**AC#2 未满足，属调度/登记纪律偏差，不是产码缺陷**。它不改写 AC#1/AC#3 的实质，但**票不能按"三 AC 全绿"收官**：AC#2 应保持不勾，由编排者在台账（`A##`）记一笔偏差，再决定是否以"AC#1/AC#3 实质达标＋AC#2 过程偏差已记"结票。

---

## 3. AC#3 —— **PASS**（两发门禁 rc=0，实测记录）

环境：Git Bash；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；**未设 `CGO_ENABLED=0`**。

| 门 | 命令 | 读数 | 取数时刻／HEAD |
|---|---|---|---|
| vet | `go vet ./internal/agent/approval/` | **rc=0**（无输出） | 18:34:19 @ `6a46140ea`；18:36:37 复跑仍 **rc=0** @ `38d05fcf4` |
| d22 | `sh scripts/d22scan.sh` | **rc=0**；`d22scan: clean - no D22 ban violations`；`examined 420 production Go files under internal/ and cmd/`；正控 `PASS=35 FAIL=0 SKIP=0`，`=== RUN=77`；`ban #8 internal/ examined 903 Go files` | 18:34:19 @ `6a46140ea`（rc=0）；18:39:31 复跑 **rc=0** @ `0af3b17aa` |
| gofumpt（本文件，非 AC#3 要求，附读） | `gofumpt -l internal/agent/approval/gate.go` | **空** | 18:36:37 |

**"纯净树"这一前提的实测口径（必须写清，别含糊）**：

- `git status --porcelain=v1 -uno | wc -l` ⇒ **0**（两发都 0）：**无任何已跟踪文件的未提交改动**，即票面所指的"纯净树"成立。
- `git status --porcelain=v1 | grep -c "^??"` ⇒ **1405** 枚 untracked（全是 `.scratch/**` 等临时件；本仓规矩"临时件只建不删"）⇒ 树**不是字面意义的全干净**，但这些路径**不在 `d22scan` 的任何 scope 里**（scope＝`internal/`、`cmd/`、`frontend/`、`design/`）⇒ 对 AC#3 读数零影响。
- 器具有非零射程（420 枚产码 Go、ban#8 internal 903 枚含注释与 `_test.go`）⇒ 这枚绿**不是空仪器扫零文件**的假绿。
- 附注（并发诚实项）：18:31 读票阶段盘上曾有 3 枚已跟踪文件为 `M`（`internal/risk/blacklist.go`／`internal/tools/bridge.go`／`internal/tools/mode.go`，属共享工作树里另一条写腿的 WIP）；到 18:34／18:36 两发门禁时已回落为 0（HEAD 同时自 `6a46140ea` 漂到 `0fa1527f9`）。⇒ AC#3 的读数一律"时刻＋HEAD"双标，不主张"任何时刻树都干净"。

**来路条款（与 AC#2 同一笔偏差，不一笔带过）**：票面 AC#3 还写着"由**那枚随批腿**交，本票不预先勾"。现实是**没有随批腿**——交件的是专为票 284 单开的 `556f0483c`。⇒ 两发读数本身**合格且被我独立复跑**，但"由随批腿交"这一来路**未兑现**；它与 §2 的 AC#2 是同一处过程偏差，不可拆开记成"AC#3 全绿"。

---

## 4. 两处附带发现（不改变 AC 判定，但必须留档）

1. **实现者交付件里一处读数不实**：`284-gate-replay-zero-caller-comment.md` §2 写"全文**没有一个数字**（`grep -c '[0-9]'` 对该块＝0）"。现跑 `sed -n '872,884p' gate.go | grep -c '[0-9]'` ⇒ **4**（`L2`／`C18`（两处）／`D47`）。**读数不实属实**，但**不影响 AC#1/AC#3 判定**——禁令射的是计数式表述，契约号不在射程（见 §1 说明）。建议编排者把该句改成"块内无计数式数字，只有契约号"。
2. **交付证据件从未提交**：`git status --porcelain docs/evidence/s1/284-gate-replay-zero-caller-comment.md` ⇒ `??`；`git log --oneline --` 该路径 ⇒ **空**；`556f0483c` 的 stat 只有 `gate.go`。⇒ 收件三把尺里的"**提交真在**"对**产码**成立（`8fdbc8be7` 在 HEAD），对**交付件**不成立（文件在盘、未入库）。是否补提交由编排者定（本裁决腿只读，不 commit）。

---

## 5. 一行可结判决

**284 可结（AC#1 PASS＋AC#3 PASS，注释三要素逐字齐、两句禁忌未触、vet/d22 双 rc=0）——但 AC#2 是由单独一条腿落地的，属明写的过程偏差，须由编排者记台账后按"AC#1/AC#3 实质达标＋AC#2 不勾"收官，⛔ 不许对外称"三 AC 全绿"。**

**撤销口令**：删 `docs/evidence/s1/284-verdict.md`，一句话复原——本裁决对产码与票面一字未动，`556f0483c` 的注释原样仍在。
