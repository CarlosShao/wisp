# 票 278 — 两条「建了但没接」的路径：`AskOnTaskRoot` 与 `Gate.Replay` 是死路还是备用路

**立票**：2026-10-08 19:3x 编排者（来路＝台账 `A720 §3②`／`A732` 普查：两条今天**产码调用点＝0**）
**性质**：裁决 + 记账；**不预设要接线**，先裁"死路／备用路"，再决定要不要守。
**Status:** **done**（2026-10-09 10:1x 编排者结案：三格全勾。取证＝腿 `278-r1` 两件正文（编排者代提 `2bdf385d`）＋**编排者自己复跑那两把调用形状尺**（各 rc=1）；两枚皆判〔备用路、不删〕，理由不同、未压成一格。差的一格转立 **票 284**（给 `Gate.Replay` 注释补"零产码调用者／勿删"，随下一枚真动 `gate.go` 的落地腿同批）；⛔ 编排者裁**不新立**"零调用者"负向钉。台账＝`A754`）

## 现量（引用前先重跑）

- `cmd/wisp/resident_approval_windows.go:697` `AskOnTaskRoot` → 函数体内 `:700` 唯一产码调用 `askConfirmation`；全仓 `.AskOnTaskRoot(` 调用形状尺 **rc=1（0 枚）**（`A720 §3①`）。
- `Gate.Replay` **零产码调用者**（`A732` 逐跳表）。
- `cmd/wisp/resident_windows.go:248` 那句 `This call is the caller` **歧义**（同句先说零调用者再说这一发就是调用者；`A725` 已记两种读法，⛔ 不许改它）。

## 要建什么

- [x] **AC#1 逐条裁**：两枚各自"死路／备用路"，每条给**判据**（谁在什么条件下会需要它；如果不接，删还是留注释）。**（2026-10-09 10:1x 编排者裁，取证＝只读腿 `278-r1`（死于服务端连接中断、两件正文完整、编排者代提 `2bdf385d`），**承重负向句由编排者自己复跑**：`git grep -nE "[A-Za-z0-9_]\.AskOnTaskRoot\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`＝**rc=1 空**；`git grep -nE "[A-Za-z0-9_]\.Replay\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`＝**rc=1 空**（两把都是**整族**结论、剥测试靠路径面不是事后剥）。**① `AskOnTaskRoot`（现量定义 `HEAD:cmd/wisp/resident_approval_windows.go:698`，票面抄的 `:697` 漂 +1；体内唯一产码调用在 `:701`，票面 `:700` 同漂）＝〔备用路，不删〕**：它服务的是"**宿主自己发起一张卡**"那一形（不经模型的工具调用），注释 `:690-693` 已具名点两枚未来来路＝语音链与票 228 的 `config.toml` 接线；今天举卡走的是**工具调用**那一条（bridge→gate→`ballCardUI.Prompt`，行为凭据见票 277 已结案）。**删的代价**＝三枚 `windows && winlive` 真机用例（`resident_approval_live_246_windows_test.go:115/220/311`）当场编不过 ⇒ CI 那道 `go vet -tags winlive`（`ci.yml:655`）立刻红，且票 246 `AC#4` 的接缝作废。**差的形状**＝`residentCard` 在产码里**零构造点**（5 枚命中全在 `resident_approval_windows.go`：注释/常量/类型/两枚签名）⇒ 接它得先有一枚构造点，属功能级、要 owner 一句话。**② `Gate.Replay`（现量定义 `HEAD:internal/agent/approval/gate.go:755`，台账 `:11139` 抄的 `:749` 漂 +6；注释块 `:751-754`）＝〔备用路，不删〕且证据比①更硬**：它的注释**逐字写着它是冻结契约 C18 的实现半件**（"C18 一键重放…requires the task to be admitted again (D47 still applies to a replay), and it is NOT an answer"），`docs/PLAN.md:1368` 的 C18 行原文含"拒绝后任务 root ctx 不取消，可一键重放"⇒ **删它＝改契约面＝人工批准**，不是代码整理。它的唯一下游 `Queue.replay`（`queue.go:598`）与 `replayOf` 字段（`queue.go:69`）同属这条两格死链；答案侧接口面 `NativeAPI`（`ui.go:143-162`）与 `PanelAPI`（`:167-171`）**都不声明 `Replay`**，控制台动词条（`approval_reply.go:566-584`）里也**没有 replay** ⇒ 入口今天结构性不存在。**两枚都判"备用路"，但理由不同、⛔ 不压成一格**（票面禁区）：①缺构造点、②缺入口且有契约背书。**两枚共同判不动的一格**＝"链接器里到底留不留它"要 `go build` 后的符号表读数，取证腿与编排者本轮都没跑 ⇒ 只能裁到"参与编译"这一层（`gate.go` 无 build tag＝全 GOOS 全档编；`resident_approval_windows.go` 双锁 windows）。）**）**
- [x] **AC#2 若判备用路**：写清"守它的最小形"（例：注释写明"未接线，非死代码，勿删"，或一枚白盒断言它仍可编译进产物）——⚠ 本格**只裁形状**，落地另派。**（编排者裁＝**只采"注释具名（形甲）＋台账具名（形丁）"，⛔ 不新立"零调用者"负向钉（形乙）**。理由三条，逐条对着本仓已付过的学费：① 负向钉在真接线那天必须改期望值，那正好是"为变绿动断言"的温床；② 本仓刚做过一轮过期注释普查（`probes/stale-claim-1/`，产码状态级断言 39 枚里 15 枚已过期），"零调用者"这一族**就是腐烂率最高的一类**；③ "没接"这件事今天已有文字面在钉——boot 报告 `resident_windows.go:269-270` 打 `任务来源：%s`，无控制台时取 `taskPostureAbsent`＝"无（任务入口未启用…）"，钉在 posture 比钉在一枚符号的调用计数上抗漂。⇒ 落地差的一格＝**`Gate.Replay` 的注释只写了契约、没写"零产码调用者／非死代码／勿删"**（`AskOnTaskRoot` 那枚 `:690-693` 已到位）⇒ **转立票 284**，⛔ 不单独开腿，随下一枚真动 `internal/agent/approval/gate.go` 的落地腿同批。形丙（白盒断言仍可引用）今天**半有**：`queue_test.go:283` 真调 `Replay`（无 tag，CI 正常跑）、三枚 `AskOnTaskRoot` 用例受 CI 的 `go vet -tags winlive` 编译门保证编得过但不执行 ⇒ ⛔ 不许用同一句话描述两枚。**）**
- [x] **AC#3 与 `resident_windows.go:248` 的关系**：给那句歧义注释一个**不改码**的处置建议（改措辞／留档），⛔ 不许顺手改。**（现量先纠行号：那句 `This call is the caller` 落在 **`:249`**（票面/台账抄的 `:248` 漂 +1；`:248` 装的是前半句 `// the way out - and AskOnTaskRoot / askConfirmation still had zero product`），注释块整段 `:246-259`，行号与短语同一次取数拿到。**裁＝〔留档，⛔ 不改措辞〕**，两种读法分开写：读法甲（`This call` 就近指 `AskOnTaskRoot / askConfirmation` 的调用者）**在 HEAD 上不成立**——编排者那把尺 rc=1（0 枚产码调用者），来源是 `:248` 用过去时 `still had`、`:249` 用现在时 `is the caller` 的时态差；读法乙（指本注释块介绍的 `:260 src := startResidentTaskSource(rt, ra)` 那一发）**成立**且有两个可读条件（`:250-252` 描述的行为只在那条腿成立）。⇒ 这是**措辞**歧义、不是事实错误；改它＝动产码注释，本票禁区明写"不改产码"，且同类"那段注释仍为真、不许当缺陷清掉"的追认在册（`A595`/K3 那一族，⚠ 保护的是 `resident_approval_windows.go` 的头注释、**不是这一句**，别混）。⇒ 处置＝本票 AC#1/AC#2 与台账 `A754` 具名留档；**若将来要改措辞，须单批批准、单独一次改动**。**）**

## 禁区

⛔ 不改产码；⛔ 不许把两枚**压成一格**裁（它们服务于不同的门路）；⛔ 零翻框；⛔ 零 push。
