# 309-r2 · 30 三处注释说实话（job 4）＋ `Options.Assessor` 写点普查（R359-4）

HEAD 面：起手 `202ed4a5`（派单锚），改完之后所有读数打在 `daa2f9be`（14:5x +08，其间只落过 `308-r2` 的件与本腿锚件，`202ed4a5..daa2f9be` 里 `internal/cmd/tools/scripts` 的**非测试源**⛔ 变 ⇒ 分母⛔ 漂）。

## 1. `bridge.go` 那把自指尺（旧 `:879-881`，现量段）

票面/orchestrator §1 那发读数（`internal cmd`，14:1x）＝**69／66／3**，多出的两行＝该注释自己 `:879`＋`assessorFor` 段 `:940`，剥行首 `//` 才＝1。本腿改完之后同一把尺（R359-10，`logs/r359-10-final-rulers.log`）＝**71／69／2**，再剥行首 `//`＝**1**：

| 口径 | 202ed4a5 | 本腿改完（daa2f9be） |
|---|---|---|
| 总行数 | 69 | **71** |
| 文件名以 `_test.go` 结尾 | 66 | **69** |
| 非测试（按**文件名**字段滤） | 3 | **2** |
| 再剥行首 `//` | 1 | **1** |

逐枚对账（`git grep -c` 打 202ed4a5 ⇄ `grep -c` 打工作树，⛔ 两把尺混用）：`toolrequest_nail_309_test.go` 7→10（本腿的词面尺段＋三形控制码），`bridge.go` 2→1（`assessorFor` 那枚 bullet 里的字面量引用改写成⛔ 带花括号的说法）。⇒ 注释里现在写的是 2 与 1，并写明"这把尺会把注释自己算进去"。

★**本腿自抓的第二枚口径洞**（原句⛔ 覆盖它，写进注释了）：`| grep -v _test.go` 滤的是**整行内容**、⛔ 文件名字段 ⇒ 一枚同时带 `ToolRequest{` 与 `_test.go` 两个字样的注释行会被这把尺**自己吞掉**。本腿第一版注释就是这么写的，量到"非测试 1 枚"的**假数**（实际 2 枚）。可重跑的演示＝R359-9（`logs/r359-9-filterdemo.log`，件在 `/tmp` 里，⛔ 落 `probes/**` 免挪别家 gofumpt 分母）：

```
./harness.go:3:// Ruler: grep -rn ToolRequest{ internal cmd --include=*.go minus _test.go yields one point.
./loop.go:4:	req := ToolRequest{TaskID: ""}
A) 文件名字段滤 = 2      B) 整行内容滤 = 1   ← 被吞的正是那行引用尺子的注释
```

处置（按裁语两选一，本腿**两样都给了**）：注释里那句自带"剥行首 `//`"那一支＋"按文件名字段滤"的写法，并把解析器尺（枚数钉）作为首选指出去。

## 2. 钉头 `:9-11` 那组过期数

旧钉头抄 60／59／1（⛔ 带 HEAD、⛔ 带时刻）。现量按 §1 那张表写进钉头，并逐字带 `HEAD 1da192e9，2026-10-11 14:4x +08，工作树面`＋"★这族数⛔ 是常数，引之前先重跑"＋309-r1 那两枚件把 69/66/3 抬上去的因由（编排者 14:1x 复跑过）。`tools`＋`scripts` 加进这把词面尺⛔ 动读数（仍 71，两枚 root 今天 0 命中，尺＝`grep -rn 'ToolRequest{' internal cmd tools scripts`）。

## 3. `bridge.go` 那句 "Every other gate still runs"

裁＝改成说实话的射程句，⛔ 为它造第二枚仪器。落点逐字（内容锚 `Read the next sentence as scope, not as coverage`，现量 `bridge.go:947` 起那一段）：结构为真的那一半**留着并具名**——守卫只决定"哪枚判级器回答那一句 `verdict := b.assessorFor(req.TaskID).Assess(...)`"，退路 `b.assess` 是"minus 探测器"的完整判级器：R5 尺＝`sed -n '211,213p' internal/tools/bridge.go` 逐字
```
		b.assess = risk.NewRiskAssessor().
			WithCanonicalizer(o.Paths).
			WithSensitiveClassifier(b.classifier)
```
（锚 `b.assess = risk.NewRiskAssessor()`＝`:211`、`WithSensitiveClassifier(b.classifier)`＝`:213`、唯一调用点 `verdict := b.assessorFor(req.TaskID)`＝`:297`）。覆盖那一半**明写⛔**：本票只钉"探测器挂没挂"这一件事；其余闸门的覆盖归各条规则自己的牙测（票 236 那两枚真喂空 TaskID 的载体断的是拒、⛔ 断规则）。

## 4. R359-4＝`Options.Assessor` 写点普查（`b.injected` 恒 false 的根据）

尺＝`grep -rn 'Assessor:' internal cmd tools --include=*.go`（件＝`logs/assessor-census.log`，14:4x +08）：

```
internal/tools/bridge.go:132:	injected bool // assess came from Options.Assessor: do not rebuild per call
internal/tools/ticket90_test.go:366:		b := New(Options{Registry: reg, Gate: gate, Modes: &t90Modes{m: m}, Assessor: assessor})
count=2 / 其中测试件 1
```

⇒ 非测试那一行＝字段自己的行尾注释（⛔ 写点），唯一的真写点在 `_test.go` 里 ⇒ **生产写点＝0 枚** ⇒ `injected` 今天恒 false。⚠ 本腿刻意⛔ 在 `bridge.go` 注释里复现 `'Assessor:'` 这个字面量（写成"the Options field named Assessor"），否则这把尺又会像 §1 那样咬到自己——那一处如果将来有人改回带冒号的写法，读数会从 2 变 3，**结论⛔ 变**（生产写点仍 0）。
