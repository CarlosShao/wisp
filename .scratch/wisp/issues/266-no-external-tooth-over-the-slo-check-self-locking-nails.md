# 票 266 — `scripts/*.ps1` 里的**门自锁钉**可以被就地删掉而本机**零枚外牙**会发现：唯一兜得住的是 `slo-freshness.sh` 的 P3 aging，而那枚牙按"天"钝（票 263 的两枚钉实测删掉即 rc=0 全绿）

**立票时刻**：2026-10-04 13:2x +08，锚点 HEAD `515ca5c5`（`dev`）
**来路**：非实现者验收腿 `263-v1` 的 `.scratch/wisp/probes/263/v1/verdict.md` §⑤/§⑥（突变名册 `delete-nails` 一发）＋票面 AC#4 那句"绕得过就具名说绕得过，别当牙"⇒ 它把"绕得过"写满了，但**补牙不在 263 射程**（263 只修"测不到"）。台账 `A599`。

## 现量（编排者 2026-10-04 13:2x 现跑，⚠ 引用前先重跑，别把这几行当常量）

1. **`tools/d22scan` 的 scope 名册里没有 `scripts/`**：尺＝`grep -c '"scripts' tools/d22scan/main.go` ＝ **0**；它的 scope 只列 `internal/`＋`cmd/`＋`frontend/`＋`internal/tools/`＋`design/`（`main.go:5`、`:24`、`:39`、`:406`）。⇒ `scripts/**` 里任何 PowerShell 文件的形状**不在任何一把既有扫描器射程内**。
2. **删钉那一发（`263-v1` 实测，本票不引用它的读数当凭据、只当来路）**：把 `scripts/slo-check.ps1` 里两枚钉的调用删掉 ⇒ **rc=0、全绿、屏上连 `shape nail ok` 两行都没有**，本机没有任何外部尺变红。
3. **既有外牙只有一枚，且按天钝**：`scripts/slo-freshness.sh` 的 **P3 sampling**—— newest `slo-full-report` **产物**年龄超 `SLO_FULL_SAMPLE_MAX_AGE_DAYS` 才响（`:60-63` 逐字"This is the probe that stays sharp when every run is a machine-contended no-conclusion run"）。⇒ 它看的是**产物**不是**钉的存在**：删钉后只要还出一份结论它就安静；真要它响，得先让 D32 连续**好几天**出不了结论。
4. 两枚钉今天的真身（行号取数时刻＝13:2x，派单前重跑）：词面钉＝`scripts/slo-check.ps1:249`（`PSCommandPath` 空）／`:256`（命中即 `Fail-InstrumentBroken`）／`:258`（ok 句），能力钉的探测在 `:223` 一形；两枚都在**任何测量之前**跑（`:286-287` 注释具名）。

## 为什么会发生（一句机制，⛔ 别当"PowerShell 的怪癖"）

一把**自锁**的钉（"我这个文件里不许出现那个形状"）天生能被"把钉删了"这个动作绕过——这不是缺陷，是这类仪器的定义。**真正的缺陷是：这个仓没有任何**外部**东西知道那枚钉应该存在。** 于是"加固"与"拆加固"在盘上留下的差异只有那两行 `ok` 句，而它红不红只由被删的那枚文件自己决定。
⇒ 净结果：D32 那道门今天有了两枚牙，但**牙的存在本身没牙**。下一任删掉它时，不会有东西响，直到 SLO 连续几天出不了结论（第 3 条那枚钝牙）。

## 要建什么

- **外牙**（形状由落地腿自己现读后选，⛔ 不许照抄本栏）：让"某枚门期望它有 N 枚自锁钉、且钉必须真跑"这件事由一枚**外部件**记数——外部件与被检文件分属两枚文件，否则"删检查也删声明"＝同义反复。候选形（⚠ 三条本票都没跑过作用面，属〔待验〕，⛔ 不许据此排除任何一条）：① CI/本机一枚独立尺读 `slo-check.ps1` 里两枚钉的**调用点**枚数；② 让 `slo-check.ps1` 在成功路径上**必须**打出那两行 `ok`，由**外层**步（不是它自己）断这两行在场；③ 把 `scripts/*.ps1` 纳入 `tools/d22scan` 的 scope 并补一条 ban（⚠ 这会改那把门的分母，动它要先具名解冻，见 `A207`/`5d463bb` 那族先例）。
- **正控**（既有定式：负向尺必配"种 X 必响"）：交件含两发——① 把两枚钉的调用删掉 ⇒ **新的外牙必须响**，红句逐字；② 还原 ⇒ 绿。⛔ 不许只证"加了东西会红"，要证**删了东西会红**（这才是本票的对象）。
- **诚实档**：外牙拿不到被检文件／读不到产物时，要**具名**报"无法判定"，⛔ 不许安静返回绿（先例：票 263 的 `Fail-InstrumentBroken`、票 134 AC#4 的 validity precheck）。
- ⛔ **不许**为了变绿动任何**阈值/golden/`internal/observe/thresholds.go` 一字节**；⛔ 不许把任何一步从"必须出结论"改成"跳过"；⛔ 不许 `Set-StrictMode` 放宽；⛔ **step 层的绿可以等于"什么都没测"**——本票不许把这一形做得更糟。

## 判据（AC 框由编排者翻，产码腿一枚都不许碰）

- [ ] **AC#0**：落地腿自己现跑复认 §现量 1/2/3（那把 scope 尺、P3 的射程那句），并现读两枚钉今天的真身行号（不许照抄 §现量 4）。
- [ ] **AC#1**：外牙的形状**落在被检文件之外**，具名 `file:line`；并答"它凭什么发现不了『删检查也删声明』"——如果它也能被同一动作绕过，具名说绕过，别当牙。
- [ ] **AC#2**：正控两发读数齐（删钉必响逐字／还原绿），并附还原证明（`git status --porcelain -- scripts tools .github` 在交件时为空）。
- [ ] **AC#3**：⛔ 四条偷懒形状逐条自查并在交件里点名"我没做哪一种"（动阈值／放宽严格模式／把某步改成跳过／用 `try/catch` 吞）。
- [ ] **AC#4**：卫生门读数不扩大（`tools/d22scan`、票面长度的门 `scripts/check-path-length-budget.sh --with-self-test`、票 212 ban #9）；红名集合逐名比对并**带取数时刻**；⚠ 若走候选③动了 d22scan 的 scope ⇒ 分母变化要具名报给编排者，不许自行翻任何豁免。

## 排程与禁区

- **按住原因（写面互斥）**：本票写面＝`scripts/**`（＋若走候选②：`.github/workflows/ci.yml`）。⛔ 与 **`262-v1`**（读 `scripts/check-path-length-budget.sh`＋`ci.yml`）同目录一律串行；先派 `262-v1`，它交完再由我核，才许动本票。
- ⛔ 不动 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）。
- ⛔ 本票**不碰** `internal/agent/**`、`cmd/wisp/**`、`internal/ball/**`（票 260 那批在飞/待验收，包级互斥）。
- ⛔ 不在本票里"顺手"给 `scripts/**` 加别的规矩——一枚外牙就是一枚外牙。
