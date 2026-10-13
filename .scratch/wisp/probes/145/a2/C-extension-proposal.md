# 票 145 · C 节 —— 扩字段提案（宁缺毋造；只建议、⛔ 不拍板）（只读腿 `145-a2`）

- 锚 = dev `8e4f6f83`。本节只给"面 + 源 + 落点 + 与邻居票的切面 + 建议落序"，**不选形、不写码、不裁归属**。
- 硬约束沿用票面 AC#1 ⓑ 与 `A408`：**一个字段若只能靠填常量才通，它不许进落地集**。

## C0 先分家：两条通道，别混

- **设置 / 名册侧**（git / model / credential / tasks / instructions）= 已有落点 **`composer` 段或独立指针段**，加子键**不动顶层键集** => 不触发那三枚字节钉（`pump_test.go:123/:291`、`subagent_roster_197_test.go:214`），只触发双向尺的**随行项**（`composer_test.go:49`、`approval_test.go:105`，属 `Q-51` 族）。
- **工作态侧**（B 表那 14 行）= 今天**没有**落点：`composer` 段装的是输入行，硬塞会把"权限输入"与"运行态"混成一段。**正解是新增顶层段（如 `run`）——而新增顶层键 = 同时撞那三枚字节钉 + 双向尺**（`A560`：两把尺只认 json tag，无 tag 字段全盲）。⇒ 工作态那批字段的落地**结构性等"三枚字节钉的解冻条件"**，不只是等 `Q-51`。

## C1 逐字段提案表

「源状态」三档：**活源-跨 seam**（进程里有、泵读不到）/ **活源-同段**（泵可直接读）/ **无源**。

| # | 名字（建议） | 来源（现读指认） | 落点 | 哪几态用 | 源状态 |
|---|---|---|---|---|---|
| 1 | `run.usage` `{in,out,cached}` | `internal/agent/cost.go:24 Usage` / `loop.go:91/:447-455/:958`；出口已到 `run.go:1125/:1136` | **顶层 `run` 段** | #2(token 计数)、#12(成本) | 活源-跨 seam（活在 loop 局部） |
| 2 | `run.cost` `{micros,currency}` | `internal/agent/cost.go:38 AddUsage` + `:50 PriceMicros`；`config.Price` | 顶层 `run` 段 | #12 | 活源-跨 seam |
| 3 | `run.reasoning` `{text,reasoningMs}` | `run.go:1268 agent.EvReasoningDelta`（**文本已在事件里**）；耗时无源 | 顶层 `run` 段 | #3 | 文本活源-跨 seam；**耗时无源**（不许填常量） |
| 4 | `run.tools[]` `{name,argsDigest,status4,level,correlationId}` | 名/参数/级/ID 有源（`tools.Decision`，见 `panel_pump.go:256 panelArgs`）；**四值运行态无稳定写者**（"started_at 无写者"是未复认否证，`145-c2` 挂着） | 顶层 `run` 段 | #4、#5 | 半活源：静态半有、运行态半**未复认** => 先复认再派 |
| 5 | `run.phase` / `run.status` | **映射表不存在**（`145-r2`：`run.phase` 无映射）；D43 名有表但那是任务态非"界面 phase" | 顶层 `run` 段 | #1 | **无源** |
| 6 | `run.failures[]` `{humanText,errorClass,detail}` | `humanText`/`error_class` 无产码来源（`145-c2`：纯无源） | 顶层 `run` 段 | #10 | **无源** |
| 7 | `pending[].position` | `internal/agent/pending_read.go:44-45 Position`（1-based）、`Queue.Depth()` | **`pending[]` 卡内子键** | #6（队列深度角标） | 活源-同段；**归票 167-r2**（见 C2） |
| 8 | `pending[].remainingMs` | **无源**：`gate.go:257 clock.After` 不给"还剩多少"，`Remaining` 只带静态窗长 | — | #8 | **无源** |
| 9 | `composer.models[]` `{provider,id,display}` | catalog `internal/config/schema.go:348 ModelSpec` ∩ `Enabled` | `composer` 段子键 | 无（14 表不要求） | 半活源：目录有；**「∩enabled」语义未裁**（`a2/census.md` §A1：缺键解码=false）=> 直做会产**假空清单** |
| 10 | `composer.efforts` | `Capabilities.Thinking`（声明位）+ 该适配器实际档位词表 | `composer` 段子键 | 无 | 半活源：声明位有、**实测核实位/适配器词表无导出面** => 缺的是**源不是字段** |
| 11 | `composer.git.*`（工作树/分支清单） | `.git/HEAD`、`.git/worktrees/`、`refs/heads`；**归票 186** | `composer.git` 子格 | 无（工作区/版本维） | 活源（待 186-a1 读数回来定"要不要起外部 git"） |
| 12 | `context.occupancy` `{used,total}` | D32 阈值 + provider 上下文窗；分子/分母口径未裁 | 顶层或 `run` 段 | #12 邻域 | 半活源：**第 0 跳**（loop 读口挂 runtime）未落（`A699`／`167-a2`）=> **归票 167** |

## C2 与三枚邻居票的切面（谁管哪一维，⛔ 别各加各的）

- **186（git 维）**：管"工作区/工作树/分支"这一维（读 + 切换请求）。与 14 行工作态**零重叠**（是版本上下文，非运行态）。
  - 切面：`GitView` 已在快照（票 181 只读半）；186 扩的是"工作树清单 + 分支清单 + 切换请求"。**切换请求 = C17 契约面**（`186 AC#2`）。=> 与 145 无键冲突，但**同一 `composer.git` 子格**，两条腿不许各加各的子键。
- **167（占用条/排队/停止/草稿）**：Status 明写"**与 145 抢同一批快照字段**"。
  - 占用条 = #12 `context.occupancy`（**167 领地**，非 145）。
  - 排队序号 = #7 `pending[].position`（**167-r2 领地**；`A699` 已裁"先答 Q-51 卡面键"——见 D 节 Q-51 同名不同物的说明）。
  - 停止 = `tasks[]` cancellable 位 + C17 第五枚入向方法（`Q-76`）。
  - => 145 与 167 **不许同时改 `pending[]` 的键集**，必须同批（票面 09-28 已记"两条腿不许各加各的"）。
- **313（模型清单 / 思考档位）**：= 票 145 `AC#2b` 的同一件事（#9/#10 两维）。
  - 313 的"界面配置面"是**写回**（改 provider/model/key）=> 必然新增 C17 入向方法（旧账 `Q-64`、新账 `Q-76`）=> **默认不做**。
  - 145 这侧只做"**只读显示**"。=> 切面：**313 管写、145-AC#2b 管读**；两维的**源**今天都缺（#9 的 enabled 语义、#10 的词表导出面），先解源再谈字段。

## C3 建议落序（**只建议**，⛔ 不拍板）

1. **第 0 步：先接出向那一跳（Go -> 页推 + 页面接收器）。** 理由 = `145-f1` 判语（`A657`）：阻断是"更瘦的通路带更肥的字段"，字段先落地仍画不出。归属 = 票 35 出向那一跳（**⛔ 未裁**：`A657 §8` ⓐ只读给清单 / ⓑ等票 274 / ⓒ机主点名）。
2. **第 1 步：只在"同段子键"上加、且源已活的那几枚**（不触发三枚字节钉）：`pending[].position`（与 167-r2 同批）、`composer.git.*`（与 186 同批）。代价最小、不撞顶层键钉。
3. **第 2 步（按住）：工作态顶层段（`run`）**。#1/#2/#3 有活源但**跨 seam**，需要"把 loop/cost 读口挂上 runtime"那一跳；且新增顶层键**必须同时处理那三枚字节钉的解冻条件**（票面记的是"两枚"，本仓实为**三枚**——见 `census.md` §E6）。
4. **第 3 步（按住）：`AC#2b`（#9/#10）** —— 缺的是源不是字段；`enabled` 语义与适配器词表导出面都未定。
5. **⛔ 永不进落地集**：#5 `phase`、#6 `failures[].humanText`、#8 `remainingMs` —— 今天**无源**，按"C0/AC#1 ⓑ 宁缺毋造"只能留空或另开票。

（以上 5 条是**建议**；排序权在编排者，字段名/落点/段名皆待裁。）
