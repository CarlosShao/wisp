# 145 筛查腿（screen）——三格未勾逐格定档：AC#2 / AC#6 / AC#2b

- 日期：2026-11-07 12:4x–12:5x +08；性质＝**只读普查腿，零产码、零选形、零翻框**（勾与 Status 一字未动）。
- 起手读数锚 `d833d8fbe`，收手 HEAD `58eca3738`（共享树自前进，本程引行号以现跑尺为准）。
- 写面：本件＋票 145 末尾 Progress log 一行；两处显式 pathspec、只 commit 不 push。

## 0. 本程现跑过的尺（可复算）

| # | 断言 | 尺 | 读数 |
|---|---|---|---|
| 1 | `Snapshot` 现役形状 | `sed -n '57,92p' internal/panel/composer.go` | 6 枚直接字段（pending/results/composer/generatedAt/instructions/tasks），后两枚 `*Section` 带 `omitempty`，各有票 200/197 注释钉住装配与双向尺后果 |
| 2 | **出向那一跳（Go→页）今天仍不存在** | `grep -rn --include=*.go -E 'PostWebMessage\|EvaluateScript\|CreateWebMessageAsJson' cmd internal` 剔 test | **0 枚**（正控仍是探针 `probes/33/p1/q2/main.go`）；与票 35 r7 腿尾"残差②＝具名死格"同读数 |
| 3 | 泵产物今天只有测试在读 | `(rt \*agentRuntime).lastPanelSnapshot()` 带括号数非 test 调用 | 0 枚非 test 调用者（145-f1 ⓑ 尺复跑同色） |
| 4 | 双向尺今天的脸色 | `go test -count=1 -run 'TestApprovalCardViewJSONKeysMatchFrontendTypes\|TestComposerContractTypesMatchFrontend\|TestComposerCurrentModelTravelsOnlyFromItsReader' ./internal/panel/`（unset GOTMPDIR/TMPDIR/TEMP/TMP＋PATH 前缀 build/sherpa-onnx） | **`TestComposerContractTypesMatchFrontend` 红**：`Go ComposerState emits [files] that interface ComposerState does not declare`（composer_test.go:78）；`TestApprovalCardViewJSONKeysMatchFrontendTypes` 绿；**`TestComposerCurrentModelTravelsOnlyFromItsReader` 绿** |
| 5 | AC#6 的对象用例在盘 | `grep -n 'TestComposerCurrentModelTravelsOnlyFromItsReader' internal/panel/composer_test.go` | `:916`，头注 `:905` 逐字自认 "is ticket 145 AC#6 for the …" |
| 6 | 票 187 射程收编 | `.scratch/wisp/issues/187-*.md` 头＋尾 | Status＝ready-for-agent **未派**（"先只读普查，再落地"）；`Q-64` 已改判＝**做**（不止显示，要能改）；其 AC#4 就是"清单要真、不许造"那一格 |
| 7 | 票 35 现状 | `tail .scratch/wisp/issues/35-panel-bridge-c17.md` | r7 仪器腿已交（入向不泄密尺有牙），**出向跳仍未做**，票面 `:42/:44/:45/:47/:49/:52` 六格未勾 |

---

## 1. AC#2（落地有源那一集）→ 档 **② 归口票 35 为主 ＋ 残角 ③（写权挂 Q-51）＝今天不可派**

**为什么不是 ①（能做）：**

1. **第一阻断不在字段层，在"有没有耳朵／有没有推送"。** 本程复跑：出向整跳生产码 **0 枚**（尺 2）、页面接收器 0（145-f1 ⓓ）、泵 marshal 出的字节非 test 读者 0（尺 3）。145-f1 §2 的判语今天仍成立：**先扩字段＝"更瘦的通路带更肥的字段"，仍画不出来。** 这座的桥就是票 35（C17：Go→页事件推＋`panel.resync`），它未交的那几格（`:44/:47/:49/:52`）正是这段前置。
2. **结构性写权锁未解。** 剩余字段扩张要过 `composer_test.go:48` 一族＋`approval_test.go:105` 两把双向尺 ⇒ 唯一非放水解法＝同一枚 commit 同时含 `internal/panel/**` 与 `frontend/src/lib/panel.ts`；而本编队零 `frontend/**` 写面（项目 AGENTS ①），写权问题＝`Q-51`（台账 `pending-and-issues.md:1086`，**owner 未答**；10-13 `145-a2` 已登记它**同名不同物**——写权未答 vs 序号键已答两件事混在一个号里）。10-07 机主"随便对接"（`A655`）放开的是**读**，票面自己已写明"读≠写"。
3. **对齐欠账已经存在且在变宽。** 尺 4：`TestComposerContractTypesMatchFrontend` 今天红在 `files` 键（Go 侧 `ComposerState` 新增、TS 未声明）。此刻再往 Go 侧加任何键，只会给这枚在册红添新同族条目——按"不许放宽断言换绿"的硬规矩，新键一枚都落不进去。

**放行条件（可判）**：①票 35 出向那一跳落地（有推送＋页面耳朵）；②owner 对"谁在同一枚 commit 写 `panel.ts`"给一句（`Q-51`，建议顺手把同名两物拆号）。两者齐之前，AC#2 的写腿**不可派**——这不是"风险大所以不做"，是**通路与写权两个物理前提都不在**。

## 2. AC#6（反向判据：每枚新字段哪枚用例断言真源）→ 档 **① 能做（只读取证程，可派）**

**为什么是 ①：** 验收表 C1 把它写死为"**本放开面内无定义域；与 AC#2 落地集非空互为条件**"。那个条件**今天变了**：AC#2 的落地集已非空（r3 的 `composer.currentModel`/`modelKnown` 两键），且对象用例在盘、本程复跑**绿**（尺 4/5：`TestComposerCurrentModelTravelsOnlyFromItsReader`，`composer_test.go:916`，三臂逐字钉"值只能来自读口"）。⇒ AC#6 第一次**有可答的对象**——结构上从"没有对象的格子"变成"有对象的格子"，可以做了。

**具体方案（如果派）**：
- 形状：只读程，零产码。产出一枚证据件（建议 `docs/evidence/s1/145-ac6-reverse-evidence.md`），对**现役每一枚后加字段/键**逐枚答一句"哪枚用例断言它的值来自真来源"：
  - `composer.currentModel`/`modelKnown` → `TestComposerCurrentModelTravelsOnlyFromItsReader`（三臂＋出口字节钉）；
  - `instructions` 段 → 票 200-r2 那批用例（逐枚具名到 test 名）；`tasks` 段 → 票 197 那批；`files` 键 → 它的来源票的用例（本程未追到具名，交该程补）。
- **答不出的一档怎么办**：候选是 145-c2 点名的恒空两枚——`ComposerState.Attachments`（`pump.go:242` 硬写 nil）与 `attachmentError`（泵从不填）。筛查判：**"答不出"照实写进表，但"从落地集里拿掉"是改码动作，不归取证程**——列成清单交编排者裁（拿掉还是补 reader），取证程不擅自选形。
- 代价：半天以内；零写面冲突（只读 `internal/panel/**`＋既有测试）；不翻任何框，勾仍归非实现者验收表。
- ⚠ 一条纪律要写进派单：禁用"待补"措辞（C1 原判），每枚要么具名用例、要么具名"答不出＋处置建议"。

## 3. AC#2b（models/efforts 两维进快照）→ 档 **② 归口票 187 ＝ 本票不派**

**为什么归 187：**
1. **射程已被具名划走**：`145-r3` 落账逐字"那两维＝票 187 射程：三家档位词表全是包内私有 `var`、全仓无导出面"；票 187 随后立票并把 `Q-64` 改判＝**做**（owner 09-28 要的是"能改"，不止显示），它把"显示"（快照载体）与"写回"（C17 白名单）收进同一张票统一裁决，避免同一片地界两程各做半截。
2. **本票 AC#2b 的两个前提已被现量推翻**（`A393` 采纳）：①票面点 5"复用 `cmd/wisp/models.go:104 cmdModels` 的取数"指错对象——那是签名下载清单（`store.manifest.Models`），不是 `llm.providers.<name>.models.<id>` 配置目录；②"各家真接受的档位词表"量出来＝`openairesponses` 三档／`anthropic` 三档／`openaichat` **一档都不映射**，三张表全是包内私有 `var` ⇒ **这一维缺的是"源"（导出面）不是"字段"**，加字段只会造一枚填不出的键。
3. **"有没有这一维"的权威位也不可直读**：`Capabilities.Thinking`（`schema.go:330`）是人填声明位；核实位在 SQLite `provider_health`，其读 DAO（`internal/memory/dao_providerhealth.go:170,178`）非 test 调用者 0，且发现式录入把 capabilities 全填 false ⇒ **不许把"未探测"画成"不支持"**。这组"声明 vs 核实"两问分开答的形状，正是票 187 AC#4 已经写好的判据。

**票 187 现状态**：ready-for-agent、**未派**（先只读普查再落地）。⇒ AC#2b 的正确动作＝**等 187 派**，本票这一格由编排者在 187 交付后按其表定归宿；若要现在动它，唯一诚实写法是在 187 名下并程，⛔ 不另开重复劳动（机主原话"不要重复劳动哈"）。

---

## 4. 一句话总表

| 格 | 档 | 归口/前提 | 今天可派？ |
|---|---|---|---|
| AC#2 | ②为主（票 35）＋③残角（Q-51 写权） | 票 35 出向跳未做（尺 2/3）；双向尺在册红（尺 4） | **否** |
| AC#6 | ① 能做 | 落地集非空、对象用例在盘且绿（尺 4/5） | **可派**（只读取证程，半天，零产码） |
| AC#2b | ② 归票 187 | r3 具名射程；两前提已被 A393 推翻；187 未派 | **否**（在 187 名下等派） |
