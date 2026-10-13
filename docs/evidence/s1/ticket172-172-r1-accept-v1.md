# 票 172 裁决表 — 落地腿 `172-r1`（非实现者验收 · 腿 `172-v1`）

- 裁决腿：`172-v1`（非实现者）。起手 `HEAD=def36d0f`，取证途中树呼吸到 `ad2dadd3`（并行腿提交）；每件读数各带时点。分支 `dev`。
- 交件：`8509e52e`（尺三件 `tier_consistency_172_{ruler,machine,cases}_test.go` 193/142/180 行）＋ `8a7ca3a4`（票尾 bullet ＋ `probes/172/r1` 13 件）＋ `569c3d19`（复核日志）。
- 读数落点：`.scratch/wisp/probes/172/v1/logs/`（6 件：`00-head`/`10-baseline`/`20-two-way-red`/`30-denominator-equality`/`31-equality-and-parser-skip`/`40-seam-ac45`，各带 rc；件数 6、每件 ≤73 行，在预算内）。
- 攻击性读数全部仓外副本或 `-overlay`；源件 sha256 起手==收手且 `git status` 对应片空（§5）。⛔ 未动票面原句、⛔ 未翻框、⛔ 未 push。

## 0. 基线（读数）

| 读数 | 结果 | 件 |
|---|---|---|
| 本腿复跑 `-run Test172 -v` | `ok` rc=0，3 PASS；分母行逐字＝`分母 40 枚工具：对得上 9，对不上 0，未落地 31；机器有而契约无 [task.spawn]`（与编排者前置同） | `logs/10-baseline.txt` |
| 独立解析器（本腿自写，不复用交件） | PLAN 与 SPEC-07 **各 37 行 markdown / (工具,档位) 41 条 / 去重 40 枚**；两文本挡位集合折叠后相等；多枚名行恰 3（`reminder.set/list/cancel`、`memory.save/recall`、`task.list/cancel`）；`fs.write` 是唯一双行工具 | `logs/30-denominator-equality.txt` |
| 缝读数 S1（负控） | `WISP172_TIER_DOC=C:/no/such/PLAN.172.md` ⇒ rc=1，`read C:/no/such/PLAN.172.md: ... The system cannot find the path specified.`（缝确在吃这个变量） | `logs/40-seam-ac45.txt` |
| 缝读数 S2（overlay 不可见） | `-overlay` 换 `docs/PLAN.md`（fs.edit L2→L1 副本）⇒ **rc=0、3 PASS**、分母仍 `对得上 9/对不上 0`＝读盘型尺对 overlay 不可见，复现 r1 的结构性事实 | `logs/40-seam-ac45.txt` |
| 缝读数 S3（缺省） | 无 env 时基线满足件内守卫（`cases:117-121`），文本侧缺省＝`machine_test.go:88` 的 `root/docs/PLAN.md` | `logs/10-baseline.txt` |

## 1. 突变台（攻"判据真假"：弄坏必须红）

| # | 自造突变（⛔ 不抄交件形） | 仪器读数 | 还原证 | 件 |
|---|---|---|---|---|
| A1 | **文档说谎·另一枚工具＋反方向**：PLAN 副本 `fs.trash` RiskLevel 格 `**L1**`→`**L2**`（恰 1 行），走 `WISP172_TIER_DOC` | **红** rc=1，红句逐字＝`fs.trash：文档写 L2、机器给 L1（探针：移入回收站 -> L1）—— 冻结文本与门给的档位不一致，两边必须由同一个 ticket 改齐`；`grep -c 文档写`＝**1** 且恰指 fs.trash；分母 `对得上 8，对不上 1` | `docs/PLAN.md` sha256 起手/收手同 `c38d859a…`；副本在仓外 | `logs/20-two-way-red.txt` |
| A2 | **机器说谎**：`-overlay` 改 `fs_write.go:739` `FSMoveDecl` 的 `Declared: risk.L1`→`risk.L2`（恰 1 行） | **红** rc=1，红句逐字＝`fs.move：文档写 L1、机器给 L2（探针：同卷内重命名 -> L2）—— …`；另一枚用例的对照组亦 fatal（`对照组：…红了（[fs.move]）`）＝同发两处牙咬 | 源件 sha256 `e5ff2e6f…` 起手/收手同；`git status -- internal/tools` 空 | `logs/20-two-way-red.txt` |
| A2b | **摘注册**：overlay 删 `BuiltinFSWriteEntries` 里 `{Tool: fsMove…}` 一行（恰删 1 行） | **红** rc=1：`ruler: probe for fs.move but it is not registered`（两枚用机器侧的用例各 fatal）。⚠ 交件 bullet 写"摘掉注册⇒落地名册钉红（9->8）"——实测先撞 reconcile fatal；**名册钉红 9→8 的真实形状是文档行消失时**（A5 现量） | 源件 sha256 同上 | `logs/20-two-way-red.txt` |
| A3 | **等值·多一枚**：副本 `task.output` 格 `**L0**`→`**L0** / **L1**` | **红**：`task.output：文档写 L0/L1、机器给 L0（探针：读一枚不存在的任务行 -> L0）—— …`（"包含"式判据会漏掉这一向） | PLAN sha256 同 | `logs/31-equality-and-parser-skip.txt` |
| A4 | **等值·少一枚**（复现 r1 原形）：副本 `fs.edit` 格 `**L2**`→`**L1**` | **红**：`fs.edit：文档写 L1、机器给 L2（探针：替换授权目录里既有文件的一小段 -> L2）—— …` | PLAN sha256 同 | `logs/31-equality-and-parser-skip.txt` |
| A5 | **解析静默跳行·落地行**：overlay 解析器丢 `fs.trash` 行（+3 行/1 hunk） | **红 3 处**：PLAN 与 SPEC-07 各 `条数 40，现量钉住的是 41`、`去重工具名 39，钉住 40`、落地名册 `…9 枚→8 枚` 且 `机器有而契约无 [fs.trash task.spawn]` | 尺件 sha256 `a65637f2…` 起手/收手同 | `logs/31-equality-and-parser-skip.txt` |
| A6 | **解析静默跳行·未落地行**：overlay 丢 `web.fetch` 行（+3 行/1 hunk） | **红**（计数钉）：`条数 40≠41` ×2 ＋ `去重 39≠40`；冻结用例仍绿＝未落地行不被名册级钉（只在计数里；已具名，见 §4(b)(c)） | 同 A5 | `logs/31-equality-and-parser-skip.txt` |

判据实现读核（点 3）：`tier172Mismatched` 用 `slices.Equal(v.doc, tier172LevelsOf(v.machine))`，两侧都先排序去重 ⇒ **集合相等**，非包含（A3/A4 两向现量坐实）；机器侧走生产路径 `pathArgs -> factsFor -> assessorFor(...).Assess`，不另列期望表。

## 2. 分母与两向现读（点 2）

- **未落地抽 3 枚**（契约有、机器无，两向现读）：`shell.exec`（PLAN:2568／SPEC-07:74 行在；非测试 Go `"shell.exec"` 命中 **0**）；`web.search`（PLAN:2546／SPEC-07:52 行在；非测试 Go 唯一命中＝`internal/risk/provenance.go:137` 的 taint 来源表，**不是 C4 注册**）；`reminder.list`（PLAN:2562／SPEC-07:68 的 `reminder.set/list/cancel` 行在，positional L1/L0/L1 解析对；非测试 Go 命中 **0**）。
- **`task.spawn` 两向现读**：机器侧 `subagent_197.go:187` `Name() = "task.spawn"`、`:222` `Declared: risk.L0`、`cmd/wisp/run.go:780` 注册 `BuiltinSubagentEntries`；契约侧两处冻结文本 `task.spawn` 字面命中 **各 0**。⇒ 交件的 `tier172Undocumented=[task.spawn]` 属实。
- 全 40 枚名册与 `tier172Landed` 9 枚（`fs.delete/fs.edit/fs.list/fs.move/fs.read/fs.trash/fs.write/task.cancel/task.output`）与基线读数逐枚对上。

## 3. AC#1–AC#5 逐格判语

- **AC#1 成立**。未修前提复核：S2 现量＝读盘尺对 overlay 不可见（所以"改真文本零红"的旧值与红形演示都只能走可指副本缝——这是演示技术的边界，不是尺的边界；真件被改后同一把尺在缺省路径上直接读出红）。修后两向：A1（另一枚工具＋反方向）红句具名 `fs.trash`＋两边档位；A4 复现 r1 的 `fs.edit` 形；件内 `Test172DocCopyThatLiesReddensTheRuler` 自证对照绿＋副本红。**红句格式＝白纸黑字"文档写 X、机器给 Y（探针：…）"**，与 AC#1 要的"指得出哪一枚"一致。
- **AC#2 成立**。A2 现量：`Declared` 摘/改 ⇒ 同一把尺红、红句具名 `fs.move` 与两边档位；A2b 现量：摘注册 ⇒ `probe for fs.move but it is not registered`。不是空心检：机器侧全走生产判定路径（声明、`pathArgs` 抽取、Facts 钩子），判据是集合相等（A3/A4），`fs.write` 的 `[L1 L2]` 覆盖支被 Facts 钩子真探到（基线读数）。
- **AC#3 成立**。分母＝整张表：两处文本各 41 条/40 枚（本腿独立解析），9 对得上、0 对不上、31 未落地、1 枚机器有契约无（`task.spawn`）——与交件 bullet 逐数相同；未落地 31 枚**如实登记不跳过**；"对不上就放宽解析"的坏形未现（解析对 41/40 是硬钉：A5/A6 现量）。
- **AC#4 成立**。`docs/PLAN.md` 与 `docs/specs/**` 一字未动：sha256 起手==收手（PLAN `c38d859a…`、SPEC-07 `a8de888f…`），`git diff -- docs/PLAN.md`／`-- docs/specs` 均空，`git status --porcelain -- docs tools/d22scan` 空；`tools/d22scan/allowlist.txt` 仍 7 行、`fs.edit` 0 命中。尺件只读（唯一 `os.WriteFile` 在 `cases:151`，写 `t.TempDir()`）。
- **AC#5 成立**。`.github/**` 零改动（porcelain 空、`git diff HEAD` 空）、对 `Test172|WISP172|tier_consistency` 零命中；`ci.yml`/`slo-fresh.yml` 未动，未接 CI＝按票面"先出读数"。

## 4. 两条"未覆盖"按票面射程判 ＋ 具名边界

1. **D34 备注/动作格的升级条款**（如 `fs.move`"跨盘 = 复制+删 → L2"、`window.manage` 关闭=L2）：**不构成本票 AC 缺口**（AC#3 字面要的是"带档位的工具行进制分母并逐枚给读数"——该行已在分母、读数已给；解析器读的是 RiskLevel 格，升级条款在备注格，交件已如实登记不假装覆盖）。**但**它意味着"文档写 L2 的那条腿 ⇄ 机器给 L2"这一半没被这把尺看着 ⇒ **残余风险、建议编排者排后续票**（跨卷探针需要第二卷/跨卷场景），不驳本票。
2. **`fs.write` 覆盖支"哪个机制抬的（R8 vs R1 下界）"**：**射程外**（票面"本票不解决"第 1 条＝不裁档位对错；归因账在 `162-v3 §4①`）。判：不属本票。

具名边界（不算"没牙"，供排票与将来接 CI 时知情）：
- (a) **未落地行无机器侧可对**：两处冻结文本**同步**改一枚未落地工具的档位 ⇒ 本尺零红（结构性；AC#3 只要求进分母与读数，未要求对没有机器的行逐枚钉档位）。
- (b) **解析器不读备注格**（同 §4.1）。
- (c) **A6 的薄牙**：未落地行丢行只由 `41/40` 计数钉抓住；若同一提交把行数改掉**并同步改** `tier172TierEntries`/`tier172DocNames` 常量，就不响（"现量钉住"的固有形；落地行与 `undocumented` 另有名册级名钉）。

## 5. 还原证 / 环境（自证）

- 全部突变在仓外 `C:/Users/swq/AppData/Local/Temp/172-v1/mut/`（副本＋overlay JSON）；**源件 sha256 起手==收手**：`fs_write.go e5ff2e6f…`、`fs_edit.go 7efedfff…`、尺三件 `a65637f2…`/`b7d9e1e4…`/`cffe296b…`、PLAN `c38d859a…`、SPEC-07 `a8de888f…`；`git status --porcelain -- internal/tools docs/PLAN.md docs/specs .github` 空。
- 跑前 `tasklist`：无 `go.exe`/`.test`（仅 cargo×2）；`GOTMPDIR=C:/Users/swq/AppData/Local/Temp/172-v1/gotmp`（C: 卷）；leg 间重包错峰执行。
- 本腿 6 件读数件 + 本表 1 枚 + 票尾 1 bullet＝写面全部。

## 6. 哪格没牙 / 没跑（具名）

- **没牙**：无。五格判据在本腿 4 类自造破坏（文档两向、机器两形、解析静默跳行两发）下**均红**，且红句都指得出具名对象与两边读数。
- **没跑**（如实登记）：① 未复跑整包 `go test ./internal/tools/...`／`go build ./...`／`gofumpt`／`d22scan`／长度门（r1 有读数；编排者前置只跑三枚测试；不属本腿派单的攻坚面）；② 未复跑 r1 的"未修树"模拟（前提由 S2 结构事实＋r1 `logs/23` 承载）；③ 未做"同时删注册与机器探针"的协调型突变（A2b 只删注册一行）。

## 7. 翻框建议

- 建议翻：**AC#1、AC#2、AC#3、AC#4、AC#5**（均成立）。
- 两条未覆盖的落法：① 不驳本票、建议编排者是排后续票还是登记残余；② 射程外。
- 一条报告精确化（不涉 AC）：r1 bullet 里"摘注册⇒落地名册钉红 9→8"的机制描述与实测不同（实测先撞 `probe for fs.move but it is not registered`；9→8 是文档行消失时的形状，A5 现量）。
- ⛔ 本表未动任何 AC 框、未动 Status、未动票面原句。
