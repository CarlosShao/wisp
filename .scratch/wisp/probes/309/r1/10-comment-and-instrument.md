# 309-r1 · Ⓑ 注释口径 + `AC#1` 恒真形仪器

## 形状与账（形 Ⓑ，编排者已裁；撤销口令「309 改形 ⓐ」⛔ 我引）
形 Ⓑ 的实话：**对空 id 的实际防御＝零**（`b.assess` 返回路径逐字节不变、行为不变）。注释不写成保证，只写成"债务＋谁该响"。真正会响的只有两枚：本 `AC#1` 仪器＋ `AC#2` 枚数钉。

## Ⓑ 注释落点（`internal/tools/bridge.go`，纯注释，⛔ 改任何行为）
1. **文档块追加**（`assessorFor` 上方，原 `:920-923` 之后）：写明 `taskID == ""` 那一肢是**故意的**、**UNDER-REPORTS 而非拒绝**——空 id 落到 `return b.assess`（⛔ 探测器），`ruleTaint` 在 `if ctx.taint == nil`（`internal/risk/rules_gateway.go`）短路，R4 那一支整条不发生、无错、无痕；其余门照走。
   - 诚实两点、每点此刻可核：① 生产无 caller 触到该肢（现读句，尺＝枚数钉 `Test309ToolRequestProductionConstructionNail` 把 `internal/agent` 之外构造点钉为 0）；② 摘掉那一肢⛔ 等于更严——`scopeMarks`（`internal/risk/provenance.go`）当**另有 scope 带 marks** 时才对未登记 scope 回 `SrcUnboundScope`，所以摘肢会把"静默少报"换成**看情形的误报风暴**，⛔ 换成正确。
   - 谁该响＝点名真写的 `Test309AssessorForEmptyTaskIDDropsTaintDetector`（两向都红），并明写"编译器不守这里，该测试是唯一仪器"。
2. **`:925` 上方具名一行**（函数体内、`if` 之上的独立注释行，⛔ 动布尔式）："`taskID == ""` here is deliberate and under-reports (R4 silently dropped…)；see doc block + `Test309AssessorForEmptyTaskIDDropsTaintDetector`"。
- 尺（①）：`grep -in 'under-report|少报' internal/tools/bridge.go` ⇒ 3 命中（`:931`/`:945`/`:952`）。
- 尺（②）：对我新增行反向 grep 未来时字样（`will |would |later|eventually|接上|以后|将来`）⇒ **命中 0**（`logs/mutation-rigs.log` 外的独立判定见 `40-gates-and-rosters.md`）。

## `AC#1` 仪器设计（`internal/tools/assessorfor_taint_instrument_309_test.go`）
判据锚落在**探测器挂没挂**这件事实上（经 `assessorFor` 拿回的判级器对同一句污点文本的 R4 差异），⛔ 读返回值里的字样、⛔ 新造导出 API／加字段／新造生产可达路径（包内直接驱动 `assessorFor+Assess`）。
- 前置：`prov.Mark(taintedTask, risk.SrcWebFetch, …, "<含 sentinel>")`（`Mark` 对没开过的 scope 也落污点）。
- ① 非空：`b.assessorFor(taintedTask).Assess("notify",{"text":含 sentinel},risk.Facts{})` ⇒ **必须含 R4**。
- ② 空：`b.assessorFor("").Assess(...)` ⇒ **必须⛔ 含 R4**（形 Ⓑ 实防＝零的原样）。
- ③ 两向拉开检查：若空/非空同为 R4 或同不为 R4 ⇒ 判"仪器瞎了"红。

## 突变台件两枚（恒真形必须响，原始读数＝`logs/mutation-rigs.log`；⛔ 让突变留跟踪文件）
按内容锚下手（认那行布尔式／`WithTaintDetector` 调用本身，⛔ 绝对行号），`finally` 还原，`git diff --exit-code` 证逐字节一致：
| 台件 | 突变（内容锚） | 期望 | 实测 |
|---|---|---|---|
| BEFORE | 无 | 绿 | `--- PASS` rc=0 |
| Rig-A（ⓐ恒真形） | 守卫去 `|| taskID == ""` 肢 → `if b.injected || b.prov == nil {` | ② 空 id 竟判出 R4 ⇒ **红** | `--- FAIL` rc=1（② 句命中 + ③ 拉开句） |
| Rig-B（ⓑ没挂） | `WithTaintDetector(b.prov.Detector(taskID))` → `WithTaintDetector(nil)` | ① 非空没判出 R4 ⇒ **红** | `--- FAIL` rc=1（① 句命中） |
- 还原：`RESTORE_BYTE_IDENTICAL=yes`（`git diff --exit-code` 空）＋ `diff -q $ORIG $F` 一致；`git status --porcelain -- internal/tools/bridge.go` 当时空。
- 结论：ⓐ（削弱守卫/摘肢）与ⓑ（探测器没挂）**两向都响**，非只交"改⇒红"。`AC#1` 那格闭合。
