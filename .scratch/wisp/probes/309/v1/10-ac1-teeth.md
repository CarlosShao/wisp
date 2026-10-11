# 309-v1 · `AC#1` 恒真形仪器有没有牙（`internal/tools/assessorfor_taint_instrument_309_test.go`）

判据（票面逐字）＝"把 `taskID == ""` 那一肢改成'空 id 也挂探测器'或整枚摘掉守卫时，指名用例必须红；反过来探测器没挂上时也要红。⛔ 只交'改⇒红'不交'换反形也红'＝那格没交。"

## 0. 正控（⛔ 突变）
尺＝`GOFLAGS= go test ./internal/tools/ -run Test309AssessorForEmptyTaskIDDropsTaintDetector -v`，射程＝工作树 `internal/tools`，台＝母仓。
件＝`logs/pristine-309.log`／`logs/mut-0-pristine.log`＝`--- PASS`／`rc=0`。
⇒ 两枚前提都成立：`prov.Mark` 那枚 `t.Fatalf` 前置没响（污点真落了），`fsBridgeWith` 给的 `prov` 非 nil（`helpers_test.go:32`，注释逐字 "The C25 engine is always wired"）。

## 1. 四形突变，逐形"落地自证＋红线条逐字＋还原"
Rig＝`logs/mut-rig.txt`（存成 `.txt` 交 `sh` 跑）；原始读数＝`logs/mut-a.log`/`mut-b.log`/`mut-c.log`/`mut-d.log`。
每形都：改前 `grep -n` 打那一行**自证落到位**（⛔ 空转的突变＝白拿的凭据），跑完 `trap` 还原并打 `cmp`＋`sha256sum`。

| 形 | 打的哪一行 | LANDS 自证（改后那一行逐字） | 红线条逐字（本腿自己跑出来的） | rc |
|---|---|---|---|---|
| **A** 摘 ` \|\| taskID == ""` 肢 | `bridge.go:954`＝算出被测值那一行本身 | `954:	if b.injected || b.prov == nil {` | `assessorfor_taint_instrument_309_test.go:64: 空 task id 竟判出了 R4（rules=[R4]）…` ＋ `:70: 两向没有拉开（空/非空同为 [R4] …）` | 1 |
| **B** 探测器不挂（`WithTaintDetector(nil)`） | `bridge.go:960`＝attach 那一支 | `960:		WithTaintDetector(nil)` | `:54: 非空 task id 的判级器没判出 R4（level=L0 rules=[]）…` ＋ `:70: 两向没有拉开（空/非空同为 [] …）` | 1 |
| **C** 整枚守卫 off（`if false`） | `bridge.go:954` | `954:	if false {` | `:64:` ＋ `:70:`（同 A） | 1 |
| **D** 肢取反（`taskID != ""`） | `bridge.go:954` | `954:	if b.injected || b.prov == nil || taskID != "" {` | `:54:` ＋ `:64:` 双红（`:70` 此形静默，见 §3） | 1 |

还原证据（四形之后各打一组，逐字相同）＝
`RESTORE cmp: byte-identical OK` / `cb86eaa7053a0f00758d642c59a8f22b881bb606abe954e9a8c03ee2d4291167 *internal/tools/bridge.go`（⇄ pristine 同一枚），
收尾 `git status --porcelain -- internal/tools/bridge.go` ＝ **0 行**。

⇒ 判据两半都交齐且**本腿复跑兑现**：摘肢⇒空 id 腿红；不挂⇒非空腿红；"换反形也红"＝形 C 兑现。

## 2. 恒真/恒假三问，逐问带读数（⛔ 形容词）
- **(a) 有没有⛔ 永不触发的比较？** 没有。两枚主断言各有"绿世界"和"红世界"的实测：`:53` 的"必须判出 R4"在正控里为真、在形 B/D 里为假；`:63` 的"必须判不出"在正控里为真、在形 A/C/D 里为假。⇒ ⛔ 永不红也⛔ 永不绿。
- **(b) 有没有对空集合的循环＝白送绿？** 唯一循环＝`hasRule309` 走 `d.RulesHit`。形 B 打印出该切片**逐字 `rules=[]`**（＝探测器没了就一枚规则都没有），此时 `:53` **红** ⇒ 空集那支⛔ 换绿；`:63` 那一支吃空集是绿的，但它被 `:69` 的"两向没拉开"与 `:53` 夹住（形 B 一发里 `:54`＋`:70` 同时红＝实测）。
- **(c) missing world（守卫那一支要的第二个值，夹具⛔ 造过）？——有，两枚，具名：**
  1. **`b.prov == nil`／`b.injected` 这两肢没有让它们分叉的世界**：夹具恒 `prov != nil`、`injected == false`。本腿把这两肢各摘一枚去跑（Rig `logs/mut-rig2.txt`，`logs/mut-e.log`/`mut-f.log`），**两发都 rc=0＝仪器全绿＝对这两肢瞎**。票面 `AC#1` 标题逐字"给**三枚**前提装上仪器"，判据只点名 `taskID` 那一肢；`b.injected` 那肢票面已另记（"归装配那族另裁"），`b.prov == nil` 那肢**⛔ 另记也⛔ 有仪器**。⇒ 欠账 `Q-309-v1-prov-nil-limb-blind`。
  2. **"`return b.assess` 那一支除了 R4 其余门照走"这句话没有夹具**：注释逐字 "Every other gate still runs; only R4 is dropped"。结构上为真（`bridge.go:211-213`＝`b.assess` 仍带 canonicalizer＋classifier，只⛔ detector），但**盘上⛔ 任何一枚断言**为"空 id 那一发仍报得出 R2/R3"负责——测试只断言 R4 一枚。⇒ 欠账 `Q-309-v1-other-gates-uninstrumented`。

## 3. `:69` 那枚"两向没拉开"的断言⛔ 单独响过（实测，⛔ 形容词）
四形读数并排：A（`:64`＋`:70`）／B（`:54`＋`:70`）／C（`:64`＋`:70`）／D（`:54`＋`:64`，`:70` 静默）。
`:70` 触发的**充要条件**＝两腿同真或同假，而那两种各自已让 `:63` 或 `:53` 先红 ⇒ **`:69` 是 `:53`＋`:63` 的逻辑推论，⛔ 可能独立报警**＝冗余保险，⛔ 是假凭据（承重仍在前两枚）。写进裁决表的口径缺陷，⛔ 阻格。

## 4. 那枚仪器⛔ 走过 `Bridge.Execute`
只在包内直调 `assessorFor(...).Assess(...)`（票面 ⛔ 新造生产可达路径＝守住；`git diff -w` 证 `:297` 那枚真调用点一字节没动）。
