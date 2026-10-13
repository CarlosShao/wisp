# 197-a2 §B — 已有落地件盘点（只读；锚 HEAD `8e4f6f83`）

## B1｜票 197 全链已入库的产码/证据件（逐笔）

| 笔／号 | 是什么 | 尺 |
|---|---|---|
| `b9fa815b` | **实体层**入库（票面 §「要 owner 带给界面那支的一跳」第一条逐字点名） | 票面第 85 行 |
| `335b8d2b` | **载体层**入库（同上点名） | 票面第 85 行 |
| `a818df46` | `197-r1`＝`cmd/wisp/subagent_selfapproval_197_test.go`（710 行，**产码零改动**）＝AC#5 那枚正控 | 票面 §「197-r1 进度」第 107–112 行 |
| `237e64f4` | `197-r4`＝`cmd/wisp/subagent_blocked_197_test.go`＝AC#4 那枚生产路径判据 | 票面第 41 行 |
| `94ea50f7` | `197-r2`＝流层溢出语义重裁（见 B2） | `git show --stat 94ea50f7` |
| `937357c3` | `197-r2` 证据件 §7 计数更正（+1/−1） | `git show --stat 937357c3` |

## B2｜`197-r2` 那笔改了什么（`git show --stat 94ea50f7`）

```
docs/evidence/s1/197-subagent-stream-r2.md   | 225 ++++
internal/panel/pump.go                       | 229 ++++------   (净改)
internal/panel/pump_test.go                  |  66 ++--
internal/panel/subagent_stream_197_test.go   | 357 +++++++++
4 files changed, 823 insertions(+), 54 deletions(-)
```

- **改的语义（commit body 逐字）**：StreamLog 到 `maxKeys` 原把最老两行拼成一行并把 Done OR 到一起；重裁为**乙**＝键集合永不合并，越界只削**每条流自己的中段**（头尾各 256 rune），行文本里显式写出削掉多少；只有过 `2×maxKeys` 硬顶才整键停跟踪并具名进 `DroppedKeys`。`DefaultStreamKeys` 仍 32（**没靠调大上限把谎言推后**）。
- **契约面刻意不动**：`Snapshot` 顶层键、`ResultChunk` 三枚 JSON 键一字未改；截断事实落「行内标记＋Go 侧出口」两处。
- **判据**：`internal/panel/subagent_stream_197_test.go` 五枚（`TestSubagentStreamKeyIsTheOneSpelling:35`／`TestNPlusOneSubagentsEachGetTheirOwnStream:57`／`TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther:146`／`TestOverflowTruncatesEachStreamsOwnHeadAndTail:209`／`TestHardCeilingDropsWholeKeysAndNamesThem:275`）；`pump_test.go` 同名旧钉按新语义重写。
- **正控**：M1 关截断⇒截断尺红；M2 塞回合并⇒一行多枚直接判红。
- ⚠ 本笔**只补 AC#3 的 Go 侧判据**，界面一等公民（把「少了 N 段」画出来）归 `197-r3`／owner。
- 尺：`git show --stat 94ea50f7`；`grep -n "^func Test" internal/panel/subagent_stream_197_test.go`。

## B3｜`197-c2` 普查被「取消／blocked」的边界是什么

- **件**：`.scratch/wisp/probes/197/r5c/census.md`（13011 字节；⚠ **由编排者代落**——派单把只读腿派成了**没有写工具的类型**，件里每行由编排者现树重跑，标「〔腿报，未复核〕」的格子不许当判据）。
- **射程＝两格**：标题逐字「票 197 残余判据现读普查（**取消语义** ＋ "**卡在等批准**"那一格）」⇒ 只覆盖 **AC#6** 与 **AC#4 的后半**（读得对不对）。
- **边界内结论**：
  - ① `task.spawn` 执行体**无裸 `go func(`**；孩子经 `agent.New→RunAsync`（recover 边界在 `observe/goroutine.go:287-330`）；收尾 watcher 名＝`subagent-finish-<id>`（**不在 D38 名册**）。
  - ② 父取消不级联（`subagent_197.go:327`）＋「写给模型看」已落地；**但孩子零停止出口**（`TaskRoster.Cancel` 当时生产零调用者、`task.cancel` 未注册）⇒ **拆出票 221**。
  - ③ 名册 `blockedOnApproval` 是**布尔、不是 D43 态**；D43 的 `Confirming`／`AwaitingApproval` 生产零使用。
  - ④ 同格**第二枚洞**：join 用 `row.TaskID` 而卡键是 `corr` ⇒ 同一任务第二张卡静默读成「没卡住」；L1 短窗口住 `Gate.windows`、**全仓无枚举口** ⇒ **拆出票 220**。
  - ⑤ AC#5 当时＝〔仅文档〕（结构上不能，但**零尺会红**）⇒ 与 `197-r1`／`197-v1` 相接。
- **⛔ 丢弃的一格（合规）**：原腿交回内容有**一行来自本编队禁读目录**，编排者整行丢弃、不转述；派单模板已据此补「结论里也不许出现那两棵目录」。
- ⚠ **编号口径差**：census §⑥ 那张差异表把 AC 排成「AC#1 实体／AC#2 每枚自己的流／AC#3 载体」，与**票面**（AC#1 实体层有生产调用者／AC#2 名册与状态同源／AC#3 每子代理一条流）**不同号**；引用该表须按内容对齐，别按号。

## B4｜`197-v1` 裁决表在说哪几格

- **件**：`docs/evidence/s1/197-ac5-selfapproval-v1.md`（52,134 字节／§0–§6 齐）；验收腿 `197-v1`（**非实现者**）。
- **只覆盖一格＝`AC#5`**（票面第 44 行那格），裁定＝**成立，附四条具名条件**。
- 关键两发（票面第 47–49 行逐字）：**M6**（`queue.go:376 allowScoped` 不校验令牌⇒十条具名红，其中被测文件 `:284`／`:291` 打出「只被面板路线递过令牌的卡真的执行了、文件真的写出来了」）与 **D3**（宿主 `AnswerAllow` 改投 `AnswerReject`⇒红**只在** `:277`／`:280` 正控两行、六发拒因一字未变⇒不是「整条路都是死的」的假绿）。
- **四条附条件（都是尺的缺陷，`197-v1` §0.2）**：C1 第七枚载具 `spawnText` 恒空（宣称了没测）；C2 `:254` 六枚用例共用 `{ErrBadGrant, ErrPanelAllow}` 两值集合（D1 原地换拒因全绿）；C3 脚注 `M1` 指 `ui.go` 的 `PanelItem`（改了全绿）；C4 `:466` 名为跨卡实为同名重放。
- **残余去处（票面第 56–57 行逐字）**：C1／C2／C3 归「下一程写腿 `197-r3`」；C4 与出向读面零尺**单开票 242**（`242` 文件名带 `-done`）。
- **C1／C2 今天仍未修**（本腿现跑）：`cmd/wisp/subagent_selfapproval_197_test.go` 的 `spawnText` 只有 `:156` 声明＋`:315` 读、**全文零赋值**；`:248-261` 那枚两值集合 `if !errors.Is(...ErrBadGrant) && !errors.Is(...ErrPanelAllow)` 逐字仍在。
- 尺：`head -60 docs/evidence/s1/197-ac5-selfapproval-v1.md`；`grep -n "spawnText" cmd/wisp/subagent_selfapproval_197_test.go`。
