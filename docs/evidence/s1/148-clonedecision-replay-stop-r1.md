# 148 — 实现腿 `148-r1` STOP 交件（未定义即停：两枚前置都卡在人工批准，产码零字节）

- 腿：`148-r1`（实现腿 · clone-decision 兄弟出口）。派单原文条件：**"judge whether the 'what the panel should receive is unwritten' part is genuinely an UNDEFINED spec — if so (未定义即停), STOP with ZERO production code and a named reason + what owner/spec decision is needed."**
- 判语：**STOP**。命名的阻塞器（两枚，彼此独立，都是票面自己立的"不许实现方自决、不许编排者替 owner 推出来"的人工批准项）：
  1. **Ⓑ（AC#3）"面板该收到什么"确属 UNDEFINED 规格**——票 148 面上逐字：「往外送的字段范围**今天没有任何一份冻结文本明确说过**"该送哪些"（这一格是空的，**别读成"已批准全送"**）」。解锁＝owner 回「**批 148 Ⓑ**」＋二选一定方向：**甲**＝具名字段白名单＋一枚会红的对账（且必须**连带处置 `Params` 再往下一层共用**——value 里的 `[]any`／嵌套 map 仍共享，`reflect.DeepEqual` 照不到，否则等于把票 146 那扇门重开一条缝）；**乙**＝明写"全字段照送"＋一句理由。⚠ 这一支改的是**票 146 `AC#2` 那句契约文字**（面板载荷口径），属 `AGENTS §0.2` 人工批准面。
  2. **Ⓐ（AC#1）`queue.go` 的禁改面未获批准**——修法本身在代码层是良定义的（见下"现量"），但票 148 面 AC#4 逐字：「**未被批准前 `queue.go`／`approval.go` 也零字节**（批准只放开 Ⓐ 那一枚函数那一行，不是放开整枚文件）」，且批准的唯一形状是 owner 回「**批 148 Ⓐ**」。本腿派单里**没有**这两个批准短语 ⇒ 不动。
- AC#2（不依赖批准的那格）本腿**未做、不冒领**：它要的是 `cmd/wisp` 车道上的 WorkPeak 实跑取数（先证 `LiveApprovals` 在那 5 秒里被调用 ≥1 次再取前后读数），与派单"DO NOT run cmd/wisp full live suite"直接冲突，且票面已排程"等 147 交件之后一起派"。一格未结，票 146 `AC#3` 那句"字面那支从未复算"照票面要求**原样保留**。

## 起手三核（全部本腿现跑，HEAD＝`86fb145eed4ee8ce0912c5468d3208fc9fe3bf8d`，2026-11-06 19:18:10 +0800；⚠ 比派单写的 `bca0990f8+` 更新，行号继续漂移——票面早已警告"下一程一律现量，别抄这里的行号"）

| 票面引的形状 | 票面行号 | 本腿现跑行号 | 是否仍在 |
|---|---|---|---|
| `func (q *Queue) replay(corr string) (tools.Decision, string, error)` | `queue.go:478` | **`internal/agent/approval/queue.go:609`** | ✅ 在 |
| 其内 `d := it.Dec`（浅拷贝原形，未走 `cloneDecision`） | `queue.go:493` | **`queue.go:624`** | ✅ 在（`:= it\.Dec` 全仓非测试 3 处命中中，`:406`/`:414` 是 `.Dec.Tool` 标量读、只有 `:624` 是整结构按值交出） |
| `viewLocked` 投影只真拷 `Paths` 一枚 | `queue.go:453-456` | **`queue.go:581-590`**（`Tool`/`LevelString()`/`Reason` 标量＋`Paths` 一枚 `append(nil,…)`） | ✅ 在（范围不同、非漏拷——按票面"别按第三枚 bug 处理"只登记） |
| 护住的门＝`cloneDecision`（8 枚引用槽位逐枚分配） | （票 146） | **`internal/agent/approval/pending_read.go:66-77`**，唯一使用点＝同件 `:126`（`LiveApprovals`） | ✅ 在 |
| `LiveApprovals` 唯一产码读侧＝`cmd/wisp/panel_pump.go` | 票面 | **`cmd/wisp/panel_pump.go:66`**（其余命中全是 `_test.go`） | ✅ 在 |
| （本腿顺手量到的新一味）`Gate.Replay`＝`queue.replay` 的唯一调用者 | — | **`internal/agent/approval/gate.go:885`（调用在 `:889`）**，其 doc 注释逐字自陈 **"Today it has no production caller"** | ℹ️ 即：`replay` 这扇兄弟出口今天**没有产码调用方**，浅拷贝危害是**潜伏形不是在跑形**——这解释了"为什么没在产线上炸"，但不解除票面（C18 一键重放的队列半边是**故意留着的实现**，注释明说"do not delete it as dead code"） |

## 如果两枚批准落了，修法长什么样（已核过落点，**本腿一字未写**）

- `queue.go:624` 的 `d := it.Dec` → `d := cloneDecision(it.Dec)`（同包函数、零新依赖、行为除 8 枚引用槽位真拷外逐字段恒等）；**`Gate.Replay`（`gate.go:893` 的 `d.TaskID = taskID`）是标量写、不在危害面内，不动**。
- 配检照票 146 `AC#2` 的规格：`ticket146_liveapprovals_backing_test.go` 那把反射尺（`queueSlotPaths`，8 枚槽位）已有同族先例可照抄形状——对 `Gate.Replay`/`q.replay` 的返回值就地改每一个引用槽位、断言 `history` 里原始记录不变；承重判据照本仓操作定义＝摘掉修法这一发从绿变红。**这套测试票面 AC#1 已写明是"不许只改实现不配检"的硬要求。**

## 本腿纪律读数

- 产码（`internal/agent/approval/**`、`cmd/wisp/**`、其余任何 `.go`）：**零字节**（未 clone、未 build、未 test——无产码 diff 即无门禁对象）。
- `docs/reports/pending-and-issues.md`：零字节（铁律：绝不触碰）。
- `frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**`：零字节。
- 本件＋票面 Progress log 追加一条＝本腿全部写入，显式 pathspec 提交，未 push。
