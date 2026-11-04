# 371-r7 —— 与票 372 / 375 / 377 的边界（AC#5'）＋ R6 内层 nil 兜底读数

- 腿 `371-r2p`（写腿）。基线 **HEAD `e5a075586`**，量于 2026-11-04 23:1x–23:4x +08；行号为**本笔落盘后工作树**现取。
- 本文只回答两件事：**谁治哪一层**（§1–§3）与 **R6 改前/改后读数**（§4）。
- 写面：`internal/perm/store.go`（R6 内层兜底）、`internal/perm/store_371_r6_test.go`（新用例）、
  `internal/panel/composer_gate_agreement_371_test.go`（nil 守卫，只改测试）、本文件。
- 未碰：`M4` 产品断言（`cmd/wisp/run_mode101_test.go:240/259` 逐字未动）、票 367／372／375 原句、
  `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden。

---

## 0. 一屏判

1. **372 治「不可问」的读侧**：`s.confirm == nil && configured != ask_every_step`
   ⇒ `PermissionMode()` 把有效档塌到 `ask_every_step`
   （`internal/perm/store.go:274`，返回值 `:269-277`），并落 `MODE-DOWNGRADE`（`:241-247`）。
   它**不碰非 nil 的腿**，也不判「问没问」。
2. **375 治「不可问」的动作侧**：`answer() == false` ⇒ L1 窗口在开窗前就拒（`internal/agent/approval/gate.go:346-350`），
   档位读侧另落 `MODE-UNANSWERABLE`（`internal/agent/approval/reach.go:132/150`）。
   它**只治 L1 那一发**，不判档位切换的那张 L2 卡。
3. **377 治记账面**：`tools: call` 行加 `mode=`（`internal/tools/bridge.go:1215`），值来自**一次**有效档读取
   （`bridge.go:286` → `cmd/wisp/approval_always.go:173` 读 `perm.Store.PermissionMode()`）。它是**读数**，不是门。
4. **371（本票）治「可问而没问」**：门放行后要求腿交回**被答过 allow 的凭据** `Card`
   （`internal/perm/store.go:384`）；零值凭据 ⇒ `ResultCardNotAnswered`（`:389-393`）。
   **R6** 再把「谓词被改坏 ⇒ panic」改成「拒」（`:369-376`）。
5. **不重叠**：三条判据分别是 `confirm == nil` / `answer() == false` / `card.Answered()`，触发面不同、
   落点不同（读侧 / 动作侧 / 档位切换），**任何一条都不能替代另一条**。§2 用「谁不触发」逐条钉。
6. **仍未做**：**R9（〔甲〕扩 M4，让所有变宽出卡）＝机主 2026-09-21 定案，须摆机主，本笔未碰**；
   其余残留见 §3。本票**没有**把 372／375／377 的收益算成自己的。

---

## 1. 逐票边界表

| 票 | 它答的问题（用户视角） | 触发判据（`file:line`） | 它拒 / 记什么 | 它**不做**什么 | 与 371 的关系 |
|---|---|---|---|---|---|
| **372**（读侧降级） | 「这台机器此刻**根本没法问** ⇒ 别让配置里的松档生效」 | `if s.confirm == nil && configured != ask_every_step`（`internal/perm/store.go:274`） | 读侧返回 `ask_every_step`（`:269-277`）；启动时落 `MODE-DOWNGRADE`（`:241-247`），`ResultDowngraded`（`:88-93`） | 不调用任何腿，不出卡，不改写 `config.toml`；**不判**「腿在但从不问」 | 371 的假腿是**非 nil** ⇒ 372 的降级**不触发**（`371-a1p` §⑤ 的同一条边界）。371 也不替代它：没腿时 371 连 GateAsk 都到不了 |
| **375**（动作层 fail-closed） | 「这张卡**没人答得出** ⇒ 那一发别执行」 | `g.answerable()`：`answer == nil` 保持旧行为、`answer() == false` 才算不可答（`internal/agent/approval/reach.go:74-79`） | L1 窗口在开窗前拒：`L1-UNANSWERABLE … decision=reject`（`internal/agent/approval/gate.go:346-350`）；档位读侧 `MODE-UNANSWERABLE`（`reach.go:132-153`） | **不碰 L2**（`gate.go` 的 L2 路径一字未动），不判档位切换的答案，不读 `Card` | 375 判「这一发**能不能**问」；371 判「这一发**问过没有**」。假腿场景里 `answer()` 可以是真（有界面）而腿从不问 ⇒ **375 看不见** |
| **377**（记账面报有效档） | 「事后从日志能看出**当时实际**是哪一档」 | `mode := b.permissionMode()` 一次读（`internal/tools/bridge.go:286`），`tools: call … mode=%s`（`:1215`）；真值来自 `rt.modes.PermissionMode()`（`cmd/wisp/approval_always.go:173`） | 每个 `tools: call` 行尾多一个 `mode=`（有效档，`ask_every_step`／`auto_approve` 实测见 `377-r1` §③④） | 不做门禁、不拒任何切换；不改 `perm` 的 `from/to` 与 `MODE-DOWNGRADE` 落点 | 371-b 的 `asked=` / `card=` 落在**同一条 `MODE-SWITCH`**（`store.go:489-495`），与 377 的 `mode=` 是**不同行的不同字段**；371 没有另造记账面 |
| **371**（本票：门判「问过」） | 「拨进全自动时，要么真弹卡并等到允许、要么不许放行」 | `JudgeCard(to, s.confirm != nil)`（`store.go:348`，定义 `internal/perm/direction.go:57`）→ `GateAsk`；随后 `!card.Answered()`（`store.go:384`） | 零值凭据 ⇒ `ResultCardNotAnswered`（`:389-393`）＋ `ErrCardNotAnswered`；审计行带 `asked=0/1 card="…"`（`:489-495`） | **不扩 M4**（`AE→AH` 仍 0 卡，`run_mode101_test.go:240/259` 未动）；不判「不可问」的读侧/动作侧（那是 372/375） | 门要求的是**「被问过并答了允许」**，不是「`confirm != nil`」。R6 再补一层：连谓词被改坏也**拒而不 panic**（`store.go:369-376`） |

---

## 2. 三条判据不可互相替代（逐条用「谁不触发」钉）

| 形状 | 372 触发？ | 375 触发？ | 371 触发？ | 结论 |
|---|---|---|---|---|
| `AA` ＋ `confirm == nil` | **是**（读侧塌到 `AE`；`Set` 侧 `GateNoChannel` 直接拒，`store.go:349-354`） | 不一定（L1 的 `answer()` 是另一条事实） | 否（到不了 `GateAsk`） | 372 治「压根没腿」 |
| `AA` ＋ **非 nil 但从不问**的腿 | **否**（`s.confirm != nil`） | **否**（`answer()` 可为真；它管 L1 那一发） | **是**（`GateAsk` 放行后 `Card` 为零 ⇒ 拒） | **只有 371 能看见**；这就是 `371-a1p` §⑤ 量出的洞 |
| `AA` ＋ 真腿 ＋ 有人答允许 | 否 | 否 | 是（`card.Answered()` ⇒ `applied`，`asked=1 card=<corr>`） | 正控：尺能读出「问过」 |
| 无答复源的 L1 写（`auto_approve`） | 是（若 `confirm == nil`） | **是**（`answer()==false` ⇒ `L1-UNANSWERABLE`） | 否（与档位切换无关） | 375 治「这一发无人可答」 |

- 三处判据的**值**都不同：372 读「腿存不存在」（`s.confirm == nil`），375 读「此刻有没有人能答」
  （`Gate.answer()`），371 读「这张卡有没有被答过允许」（`Card.Answered()`）。**没有一个能推出另一个。**
- 371 **刻意不**把 372 的降级或 375 的拒绝写成自己的收益：假腿是**非 nil**，372 不触发；假腿场景里
  375 的 `answer()` 完全可能为真，375 也不触发。两者的读数分别见 `docs/evidence/s1/371-a1p-impl.md`
  §⑤（洞）与 `docs/evidence/s1/375-r1-impl.md` §④⑤（不可问侧）。

---

## 3. 本票仍然**没有**做的（残留）

1. **R9（〔甲〕）——扩 M4：所有「变宽」都出卡。** 今天只有切进 `auto_approve` 吃卡：
   注释 `cmd/wisp/run_mode101_test.go:240`、断言 `:259`（`want exactly [auto_approve]`）**逐字未动**；
   `AE→AH` 仍然 0 卡、腿一次都不进（`internal/perm/store_367_test.go:148-159` 的既有牙照旧）。
   这是机主 **2026-09-21 17:1x** 的定案 ⇒ **须摆机主**，本笔**未做、也未自决**。
   （`AE→AH` 变宽的「静默」不是本票的洞：M4 说它本来就不该出卡。）
2. **凭据的信任边界是同进程。** `Card` 由本进程铸造（`internal/perm/store.go:138-144` 的 `THE LIMIT`），
   挡的是**事故性的假腿**，挡不住蓄意调用 `AnsweredCard`。跨进程/跨用户证据不在 371 的射程。
3. **panel 不判「答没答」，这是设计。** 它不出第二张卡（`internal/panel/composer_handlers.go:44-50`），
   只做两件事：同一个 `JudgeCard`（`:176`）与把 `Set` 的拒按**同一 result 词面**记自己一行（`:207/:232`）。
   「两处同判」＝同谓词 ＋ 同词面 ＋ 两层都留行，**不是** panel 也去问一次。
4. **372／375 的读侧与动作侧行为本笔一字未改**：`PermissionMode()` 的降级（`store.go:269-277`）、
   `MODE-DOWNGRADE`、`answerable()`、`UnsilencedMode` 均未动；371 只在 `GateAsk` 里加了 R6 内层兜底。
5. **R8 基线读数**（`AE→AH`／`AA` ＋ 真腿）与 **R1–R5** 的交付在 `371-a1p-impl.md`、`371-r1p-impl.md`
   与 `cmd/wisp`／`internal/panel` 的 371 用例里；本文不重复抄读数。

---

## 4. R6 读数（内层 nil 兜底，panic-before / refuse-after）

### 4.1 红（pre-R6 代码 ＋ 本笔新用例）

命令（`-overlay` 把 `store.go` 换回 `HEAD` blob，仓内产码零变异）：

```bash
git show HEAD:internal/perm/store.go > .scratch/371r2p/mut/store_pre_r6.go
go test -overlay=".scratch/371r2p/mut/overlay-pre-r6.json" ./internal/perm \
  -run 'Test371R6' -count=1 -v -timeout 200s
```

逐字（截取）：

```
--- FAIL: Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal 0xc0000005 code=0x0 addr=0x0 pc=0x7ff7997bca75]
...
github.com/CarlosShao/wisp/internal/perm.(*Store).Set(...)
	D:/work/workspace/projects plans/Wisp/internal/perm/store.go:364 +0x1b5
...
FAIL	github.com/CarlosShao/wisp/internal/perm	0.040s
FAIL
rc=1
```

⇒ **红行＝`store.go:364` 的 `card, err := s.confirm(ctx, sw)`（HEAD 行号）**：nil 通道被调用 ⇒ panic。

### 4.2 绿（R6 后 ＋ 同一用例）

```
=== RUN   Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking
--- PASS: Test371R6TheGateAskPathRefusesANilChannelInsteadOfPanicking (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/perm	0.044s
rc=0
```

用例同时断言：`ResultRefused`、detail 逐字含「内层 nil 兜底」、**不含** `GateNoChannel` 的
「没有接入确认通道」、腿被调用 0 次、`persisted=ask_every_step / saves=0`、审计恰好一行。

### 4.3 破谓词形（票面 `F4` 的原始动机）

临时探针（仓外，`-overlay` 把 `direction.go` 的 `GateNoChannel` 改成 `GateAsk`）＋ 一枚
nil 通道的 `Set`：**改前 panic、改后拒**。

```
# 改后（本笔产码）
371R6|broken-predicate|err=perm: 从 ask_every_step 切到 auto_approve 需要一次 L2 强确认，
  但确认通道为 nil（内层 nil 兜底：门谓词被改坏时拒而不 panic，未切）
  |result=refused-confirm|persisted=ask_every_step|saves=0
--- PASS: Test371R6ScratchBrokenPredicateWithNilChannelRefuses
```

```
# 改前（同一 broken predicate ＋ HEAD blob 的 store.go）
panic: runtime error: invalid memory address or nil pointer dereference
... internal/perm/store.go:364 ...
rc=1
```

⇒ **谓词被改坏不再能打崩进程**：R6 的判据是「拒」，不是「panic 当见证」。
