# 308-r1 · 30-comment-fixes-A843（ⓓ 两枚注释缺陷，归口＝`A843` 末节 ⇒ 票 308）

★ 本节两枚**⛔ 算在 `AC#1` 的账上**（编排者裁语 ⓓ 末段）；`AC#1` 只吃 §10 那三件（helper＋直喂字面量用例＋⛔ Skip）。本腿在**同一枚文件**里顺手把两枚过期注释改成"此刻可核"，改动**只落在注释行**（尺见每节末）。

---

## ① 指⛔ 存在的件（`pathresolver_anchor_spelling_windows_test.go` 头注释里 `See probes/72/r1/20-canonical-form-choice.md` 那行）

现状（改前，尺＝`grep -n 'See probes/72' internal/risk/pathresolver_anchor_spelling_windows_test.go` 在 `HEAD~`）＝指名 `probes/72/r1/20-canonical-form-choice.md`。

盘⛔ 有该件（两把尺坐实，⛔ 我造）：
- `ls .scratch/wisp/probes/72/r1/` ⇒ 四枚 `.md`（`00-anchor`／`10-class-a-expected-side`／`20-class-b-negative-nails`／`30-evidence-readings-gates-mutations`）＋`logs/`＋`run-two-surface.sh`；⛔ `20-canonical-form-choice.md`。
- `find .scratch/wisp/probes -iname '20-canonical-form-choice*'` ⇒ **0 命中**（全仓⛔ 那枚）。

选型实际住在 `10-class-a-expected-side.md` 的 `## #2 TestCanonicalInputGainsNoSecondForm` 一节（尺＝`grep -n '^##' .scratch/wisp/probes/72/r1/10-class-a-expected-side.md` 命中 `## #2` 那条）。

改后（本腿唯一改动＝指针）＝`See probes/72/r1/10-class-a-expected-side.md (section #2, TestCanonicalInputGainsNoSecondForm) for why the`。措辞⛔ 未来时、⛔ "永久保证"；只指盘上可核⛔ 那枚件⛔ 那一节。

`git diff -U0` 尺（只该 hunk）＝`-` 一行、`+` 两行，三行⛔ 皆 `//` 注释行；⛔ 可执行码⛔ 变（`grep -vE '^[+-][[:space:]]*//'` 于该 hunk ⛔ 命中）。

---

## ② 现时假话（同文件头注释 `pathresolver_junction_windows_test.go, … is deliberately left byte-for-byte untouched`）

现状（改前）＝声称 `junction_windows_test.go` "byte-for-byte untouched"——是**⛔ 现在时**。

作废它的 commit（本腿当场复尺，⛔ 抄 `A843` 数字）：
- `git show --numstat --format= f278c57f -- internal/risk/pathresolver_junction_windows_test.go`＝**`36 7`**（＋36 −7）——`f278c57f` 自己就改了那枚件（那是票 72 甲族 #5/#6/#7 的期望侧修复，见同件里 `Ticket 72 甲族 #5/#6/#7` 三段注释）。
- 现 HEAD `1ad48333` 之后⛔ 该件又被本腿改过（守卫落点 ＋12 −1），"⛔ 变"更⛔ 成立。

改后措辞＝改成"此刻可核"⛔ 现在时⛔ 未来时：
> Ticket 308 does not edit that junction file's fixtures here: read its expected values, anchors and build tags as the measurement, not as a frozen guarantee -- commit `f278c57f` (ticket 72 甲族 #5/#6/#7) already changed that file (+36 -7), so "deliberately left byte-for-byte untouched" was a claim that stopped holding the moment that commit landed.

三个"⛔ 此刻可核"⛔ 是：① 指名具体 commit 号 `f278c57f`（可 `git show` 复尺）；② 具名 numstat `+36 -7`（可复尺）；③ 措辞用**过去时**（"was a claim that stopped holding"）描述那⛔ 注释的历史，⛔ 对未来任何时刻作保证、⛔ 用现在时把⛔ 存在的事写成"是"。

`git diff -U0` 尺（只该 hunk）＝`-` 三行、`+` 六行，九行⛔ 皆 `//` 注释行；⛔ 可执行码⛔ 变（`grep -vE '^[+-][[:space:]]*//'` 于该 hunk ⛔ 命中）。

---

## 第三枚过期注释：登记、⛔ 本腿顺手改

票 308 §现量 另有第三条注释缺陷——`TestClassifySpellingInvariance` 头上那句 "every other spelling is unconditional, so the test never quietly stops measuring"（在 `vol==""` 那形下"unconditional"半句为假）。票面对它的处置逐字＝"**登记、⛔ 顺手改对**"，且⛔ `A843` 归口只列 ①② 两枚。⇒ 本腿⛔ 改那句，留作独立待裁项（我的守卫只⛔ 现量地修 ①②；那句要不要随守卫落点一起改写，归 `308-v1`／编排者定，⛔ 由落地腿自勾）。

---

## 归口

- 本节两枚**⛔ 计入 `AC#1` 的账**（裁语 ⓓ：AC#1 只吃 helper＋直喂字面量用例＋⛔ Skip 那三件）。
- 归 `A843` 末节具名归口＝票 308；`308-a1` 已在 `probes/308/a1/20-comment-defects.md` 复尺这两枚。
- 非实现者终裁 `308-v1` 建议尺：① `git show <本腿代码 commit> -- …anchor…` 里这两处 hunk 全是 `//` 行（`grep -vE '^[+-][[:space:]]*//'` 于该 hunk 应无命中）；② 新指针 `probes/72/r1/10-class-a-expected-side.md` 在盘且含 `## #2`；③ 文中 `+36 -7` ⇄ `git show --numstat --format= f278c57f -- …junction…` 逐字同数。
