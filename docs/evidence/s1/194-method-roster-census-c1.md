# `194-c1` — 面板方法名册：只读代价普查（零产码）

- 程：`194-c1`（只读普查）｜派单：`.scratch/wisp/dispatches/2026-09-28-170x-readonly-194-c1-method-roster-census.md`
- 工单：`.scratch/wisp/issues/194-...-owner-ruled-align-the-code-to-the-spec.md`（**本表不勾它的任何 AC 框**）
- 起手时刻（本机 `date "+%Y-%m-%d %H:%M %z"` 现读）：`2026-09-28 16:12 +0800`
- 起手锚点：HEAD＝`8ca2291f`（`git log -1 --format='%h %s'` 现读，与派单标称的 `8ca2291f` **一致，未漂**）
  - ⚠ 记一笔时钟差：派单面写「派单时刻 `2026-09-28 16:5x`」，本程现读 `date` 是 **16:12 +0800**，比派单标称**早 38 分钟**。本表一律以**现读**为准；差因未定性（可能是派单时刻被前推书写），**不当缺陷、不据此推断任何事**。

## §1 起手脏件名册（写面闸门的起手基准）

尺：`git status --porcelain -- internal/ cmd/`（起手现跑）

```
（空——零枚具名脏件）
```

⇒ 起手名册＝**空集**。⚠ 按派单 §1「写面闸门（`A374`）」：闸门是**终态 == 起手名册（逐枚具名差集为空）**，**不是"必须为空"**。
同树在册飞着的写腿：`181-r2`、`185-r1`。终态若多出任何一枚，**具名登记为"不是我的"，不动、不提交、不还原**。
本程自己的写面只有三处（派单 §1）：本表 · `.scratch/wisp/probes/194/c1/**` · 票 194 的 Progress log 追加。

## §2 Q1 两份名册到底有几枚（逐枚现跑）

### 2.1 规格那份（`docs/specs/SPEC-08-ui-ball-panel.md`，表体逐枚具名）

尺：`Read` 现读 `SPEC-08:155-188`（对应派单给的 `sed -n '158,178p'` 窗口）。表体在 **`SPEC-08:161` 的表头**之下，行 **`163`–`174`**。

**逐枚方法名（把每行里 `/` 并列的两枚拆开各算一枚）**：

| # | 方法名 | 出处行 | 表内标的方向 | 需原生侧授权（逐字） |
|---|---|---|---|---|
| 1 | `panel.resync` | `SPEC-08:163` | **Go→前端推送** | `—` |
| 2 | `tasks.list` | `SPEC-08:164` | invoke | `—` |
| 3 | `task.detail` | `SPEC-08:164` | invoke | `—` |
| 4 | `history.query` | `SPEC-08:165` | invoke | `—` |
| 5 | `transcript.get` | `SPEC-08:165` | invoke | `—` |
| 6 | `approval.current` | `SPEC-08:166` | invoke | `—` |
| 7 | `approval.queue` | `SPEC-08:166` | invoke | `—` |
| 8 | `approval.decide` | `SPEC-08:167` | invoke | 「**「allow」拒绝一切面板来源（F2）；仅 `reject` 可面板发起**」 |
| 9 | `config.get` | `SPEC-08:168` | invoke | 「安全节放宽走 L2 重新确认（SPEC-03 §4.2）」 |
| 10 | `config.set` | `SPEC-08:168` | invoke | 同上一字 |
| 11 | `grants.list` | `SPEC-08:169` | invoke | `—` |
| 12 | `grants.revoke` | `SPEC-08:169` | invoke | `—` |
| 13 | `privacy.purge` | `SPEC-08:170` | invoke | 「purge 需确认」 |
| 14 | `privacy.export` | `SPEC-08:170` | invoke | 「purge 需确认」（同行，只管 purge） |
| 15 | `cost.summary` | `SPEC-08:171` | invoke | `—` |
| 16 | `models.list` | `SPEC-08:172` | invoke | `—` |
| 17 | `models.delete` | `SPEC-08:172` | invoke | `—` |
| 18 | `diagnostics.export` | `SPEC-08:173` | invoke | `—` |
| — | **事件推送那一行**：`task.delta`／`tool.chip`／`approval.request`／`ball.state`／`cost.tick`（5 枚，方向 Go→前端） | `SPEC-08:174` | Go→前端 | `—` |

**枚数（本程现数）**：
- 表体 `SPEC-08:163-173` 展开后＝**18 枚具名** ⇒ **票 194 面上写的「18 枚」是对的，不需更正**（这一支眼睛数赢了）。
- 加上 `SPEC-08:174` 事件推送那一行的 **5 枚** ⇒ 那张表**具名总面＝23 枚**。
- ⚠ **一处具名更正（不是枚数，是行号）**：票 194 AC#2 与本派单 §2 把 allow 侧禁令引作「`SPEC-08:169`」、把 `config.set` 的 L2 重确认引作「`SPEC-08:170`」。**磁盘现读**：allow 侧禁令在 **`SPEC-08:167`**，`config.set` 那一字在 **`SPEC-08:168`**（`169` 是 `grants.*`，`170` 是 `privacy.*` 的「purge 需确认」）。**两处引用各偏 2 行**；引文本身逐字对得上，只是行号错位。⇒ 影响：写腿若照票面行号去定位会改错行（把 `grants.revoke` 那行当 allow 侧禁令）。**本表只更正行号，不改 `docs/specs/**` 一字（AC#6）。**
- ⚠ 另一处口径要摊：`panel.resync`（第 1 枚）**在 invoke 白名单的表里，但方向标的是「Go→前端推送」**——它不属于「前端调后端」这一族。票 33 前例说的「有名无实现」正是它。**它算不算第 19 枚要补的"门"，是裁定面，本程只标不裁。**

### 2.2 代码那份（现量）

尺：`grep -nE '=\s*"panel\.' internal/panel/bridge.go` ⇒ **4 枚**；`grep -cE` 同尺 ⇒ **4**（与票 194 面一致；⚠ 未跑 `grep -c 'panel\.'`，那把尺会把 `bridge.go:35` 的注释 `"panel.*"` 计进来＝5，`33-r1`／`A380` 已具名报过，本程**不复量那把错尺**）。

| # | 常量名 | 字面值 | 出处行 |
|---|---|---|---|
| 1 | `MethodModeRequest` | `panel.mode.request` | `bridge.go:42` |
| 2 | `MethodWorkspaceRequest` | `panel.workspace.request` | `bridge.go:43` |
| 3 | `MethodAttachmentAdd` | `panel.attachment.add` | `bridge.go:44` |
| 4 | `MethodMessageSend` | `panel.message.send` | `bridge.go:45` |

`knownComposerMethod`（`bridge.go:104-110`）的 switch 就是这四枚常量，无第五枚。

### 2.3 交集

- 规格 18 枚（＋事件 5 枚）∩ 代码 4 枚＝**零枚**，与票 194 面「交集为零」一致（本程现读两侧名册，未见任何一枚同名）。
- ⚠ **反向也要摊**：代码那 4 枚 `panel.*` **在规格那张表里一枚都没有**。所以「代码向规格对齐」不等于「只往代码里加 18 枚」——**现有 4 枚是规格没写的既有门**，动它们（改名／换顺序）会撞 `composer_test.go:502` 那枚钉（票 194 AC#3 已具名警告，出处 `181-r1` 的按磁盘源码数枚数的尺）。本程判定：**4 枚保留，名册是"并集"问题不是"替换"问题** ⇒ 这是给编排者的一枚**待裁形状**（见 §6）。

## §3 Q2 逐枚四问主表

（待填：23 枚逐枚的 ① 等价实现尺现跑 ② 落点／最接近地基 ③ 撞不撞 AC#2 三条禁区 ④ 需不需要新门控与审计）

## §4 Q3 事件推送那一行（5 枚，Go→前端）

（待填）

## §5 Q4 三堆分批建议（不拍板）

（待填）

## §6 Q5 债务清点：`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿到哪一步

（待填）

## §7 本程没测／判错的（逐名）

（待填）

## §8 门禁终态（只读程取数，`A363`：不充当结案凭据）

（待填）

## §9 被拒调用＋零删除自证＋工具调用终值

（待填）

## §10 next＝派写腿之前还缺哪几枚批准

（待填）
