# `185-paged-reread-r1.md` — 票 185 第四支（per-scope 宿主落盘名册）落地腿 `185-r1` 的交付表

- 写手腿代号：`185-r1`｜性质：写码腿（动产码）｜**AC 框一枚没勾**（本票面 7 格全保持未勾）
- 派单：`.scratch/wisp/dispatches/2026-09-28-170x-impl-185-r1-paged-reread-per-scope-roster.md`
- 工单：`.scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-…md`
- 本表按派单 §4 的十二节骨架逐节填；每节自报"这一发在**未修码**上响不响"。

---

## ① 起手锚／起手脏件名册／授权出处

- 起手 HEAD（`git log -1 --format='%h %ad %s'` 本程现跑）：
  **`8ca2291f` Mon Sep 28 16:02:10 2026 +0800 `ledger(A380) 收 33-r1 片 A…`**
  ⇒ 与派单写的"盘上是 `8ca2291f`"**一致，未漂**，不需登记漂移。
- 起手 `git status --porcelain -- internal/ cmd/`（逐枚具名，本程闸门基线）：
  ```
   M cmd/wisp/panel_pump.go
   M internal/panel/git.go
   M internal/panel/git_test.go
  ```
  三枚都是另一枚写腿 `181-r2` 的在飞件 ⇒ **不碰、不提交、不还原、不评论**；本程凡依赖它们的
  端到端腿一律用 `git cat-file blob HEAD:<path>` 钉进自己的 overlay（见 §⑦）。
- 本票的"未修码"具名＝**起手 HEAD `8ca2291f`**（票 183 的修法已在场，尺：
  `grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/taintmatch.go | grep -v _test.go` 有定义有调用；
  本程 §② 的红腿全部在 HEAD 钉死的产码上跑，不是在自己改过的码上跑）。
- 授权出处（台账 `A381`，owner 原话逐字）：
  - 批准①原话：「**好的，那超长输出，我批了，要做的**」⇒ `Q-66` 结案＝**做，走第四支**
    （豁免绑"宿主亲手落盘过的那一枚文件"，不绑"模型报进来的路径字符串"）。
  - 同一句**翻掉票 177 的"per-scope 路径名册＝未批"**（`A353`），并要求把
    `internal/risk/taintmatch.go` 那句 "a per-scope path roster is explicitly NOT approved"
    **改成实话并具名指向 `A381`** ⇒ 已改，逐字见 §⑫ 前一格（结束消息 ⑦）。
  - 撤销口令：**「撤分页」**。
  - ⛔ 三支毒修法不在批准范围内（放宽 R4／参数侧按值放行／把 `fs.read` 摘出 C25 名册／
    设 `Config.PassThroughUnclassifiedRisk`）⇒ **本程一枚都没用，也没判"必须用"**（见 §③ 的形状选择）。
- 起手三枚受保护文字 md5（本程现跑，全部＝基线）：
  - `sed -n '468,473p' internal/risk/provenance.go | md5sum` ＝ `858e45116383caa3e7c1dd4b0924fad1` ✓
  - `sed -n '11,15p' internal/risk/taintmatch.go | md5sum` ＝ `5680ddd18e2d2ec2a85e485b54f4c12e` ✓
  - `sed -n '482,506p' internal/risk/provenance.go | md5sum` ＝ `f89e891e5eee3f3ea2b4f89d921072c4` ✓
  - 另记一枚**仅参考**基线（派单允许随改动变）：`sed -n '100,109p' internal/risk/taintmatch.go | md5sum`
    起手＝`23aea9de4974032a480aa16f52114014`（⚠ 台账 `A363` 记的是 `2bb6d6dd027eeb0446d16dba6e392963`，
    **两者不等＝台账那一格对的是旧锚点或旧口径**；本程把起手值当基线，终态见 §⑨）。

## ② 未修码红：AC#1 两发＋AC#4 三形逐字读数

（待填：判据件写完、在 HEAD 钉死的产码上现跑。）

## ③ 落地形状：名册存在哪、谁写、谁查、活多久

（待填。）

## ④ AC#2 跨 scope／跨任务那两发

（待填。）

## ⑤ AC#3 覆盖面逐枚名册（现跑尺）

（待填。）

## ⑥ AC#5 代价面

（待填：哪一类外来正文不再算证据＋3/10 那把尺的复量与语料组成。）

## ⑦ AC#6 真机两发逐字

（待填，含在册 PATH 前缀两形。）

## ⑧ 变异自证（含"两枚哈希非空"）

（待填。）

## ⑨ 门禁全套

（待填。）

## ⑩ 本程没测的（逐名）

（待填。）

## ⑪ 被拒调用＋零删除自证＋工具调用终值

（待填。）

## ⑫ 对编排者的不服

- 派单 §4 要求"第 ≤5 枚调用内落盘并 commit 第一枚"；本程前 **12** 枚调用全在读派单／工单／
  证据／产码与台账（派单 §开头写"再读工单本体与它的上游证据（必读）"，且 §2 每格都要"先测→再写"
  的未修码读数，量与行号都不许照抄票面）⇒ 骨架表落在第 **17** 枚。**这一格是本程的纪律违规，自报。**
  后果评估：零产码改动、零脏件、无假读数（起手名册在动码前就取全了）。

## ⑬ appendix: machine readings (ASCII index; narrative cells ②-⑫ in the leg reply, this leg ran out of tool calls)

- logs: .scratch/wisp/probes/185/r1/logs/tools-185-landed.txt, risk-185-landed.txt, pkg-tails.txt, panel-fails.txt
- mutants + overlays: .scratch/wisp/probes/185/r1/mut-m1..m4, ov-m1..m4.json (pre/post md5 in the leg reply)
- unfixed-anchor red readings: package risk and tools runs on HEAD production (git diff --stat empty), text in the leg reply
