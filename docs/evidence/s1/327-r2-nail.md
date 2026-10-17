# 327-r2 落地件：把名册基数常钉搬进仓内（票 327 残余 R1）

写腿 `327-r2`（写 `internal/panel/**` 测试面，只追加）｜锚＝`7535d83b`（改前对照）｜
落点提交 `3e1b425d`｜台件＝`.scratch/wisp/probes/327/r2/`（`mutate.sh`／`mutate2.sh`／
`logs/color.log`／`logs/mutations.log`）｜⛔ 翻框、⛔ 改判据、⛔ push。

## 1. 这一枚在闭合什么

`docs/evidence/s1/327-v1-verdict.md` 第 5 节残余 R1：往 `isDisplayWashCf` 的名单里
**静默加一枚 `U+206B`**，整包多红＝**0 枚**（`PASS_TOP 155／FAIL_TOP 7` 与终态逐枚
同名）。今天抓得到它的只有判语腿那份**仓外台件** `zz_v1_exhaustive_test.go`（只装进
临时那台）。本件把这枚钉做进 `internal/panel/display_scrub_327r2_test.go`，成为
**仓内常钉**。既有测试件的排除钉 `TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`
只点 7 枚、**漏了 `U+206B`**——R1 就是从这道缝里漏的。

## 2. 钉了什么（三枚新用例，纯新建文件、既有断言一字未动，numstat=160/0）

- **基数钉** `TestDisplayWashRosterCardinalityIsExactlyTheRuledFifteen`：实际侧
  `productWashedRunes()` 遍历全标量（`0x0000..MaxRune`，跳过 surrogate）只问产码那枚
  `isDisplayWashCf`，**字面量只当期望侧、绝⛔ 当实际侧**。与裁定期望名册**双向成对
  比较**：`extra` 臂（`:72`，多洗＝静默加宽）⇄ `missing` 臂（`:76`，少洗＝静默摘除）
  ＝两枚不同断言行。比较本身＝钉。
- **排除钉** `TestDisplayWashRosterExcludesTheEightNamedRunes`：判语件点名的 8 枚
  （`U+2061-2064`／`U+061C`／`U+180E`／`U+206A`／**`U+206B`**）逐枚断 `washed=false`，
  并验 `needsDisplayWash` 不响、`scrubForDisplay` 不改。`U+206B` 就是缝里那一枚。
- **两侧同集钉** `TestDisplayWashJudgeAndWriterAgreeOnBothFamiliesAndBoundaries`：
  `needsDisplayWash` ⇄ `scrubForDisplay`（含 ESC 强开快路径那一半，使写入侧失明无处躲）
  零分歧。射程写明：`unicode.Cf`（现量 170）＋ `unicode.Cc`（65）＝**全族逐枚**；
  其余为**具名边界抽样**（22 枚），⛔ 把"全标量"偷降成"抽样"却不写明。终态
  `disagreements=0`。

## 3. 颜色尺（仓外导出树，同一台面；四数尺见 logs/color.log）

改前（锚 7535d83b）`RUN 252／PASS_TOP 155／FAIL_TOP 7／SKIP 0`
改后（加本文件）    `RUN 255／PASS_TOP 158／FAIL_TOP 7／SKIP 0`
红名册**作差＝0**（起手 7 枚 ⇄ 交件 7 枚逐枚同名；`TestC21DesignTokensFourWayAgree`
那枚＝导出台缺 `design/` 的台面射程，⛔ 计入判据，且跨每一发不变）。
新增用例只让 `RUN`／`PASS_TOP` 各 +3。

## 4. 成对两发＋摘除对照（都打在仓外导出树字节，逐发 finally 还原；logs/mutations.log）

- 发一（**加** `U+206B`）：`FAIL_TOP 9`（比交件 +2）。基数钉 `:72`（extra 臂）＋`:80`
  红，排除钉 `:91`/`:95`/`:98` 红。⇒ 这正是判语腿 R1 里"整包⛔ 响"那一发，如今响在
  两枚仓内常钉上。
- 发二（摘既有 `U+00AD`）：`FAIL_TOP 10`。基数钉 `:76`（missing 臂，与发一**不同
  断言行**）＋`:80` 红；既有 AC#2 靶子 `display_scrub_216_test.go:368`／`:411` 同红
  （发二独有名册）。
- 发一 ⇄ 发二红名册作差：发二独有 `TestNeedsDisplayWashNamesEveryWashedCfRune`／
  `TestEveryNamedCfRuneIsEscapedAndNeverDeleted`；发一独有 `TestDisplayWashRoster
  ExcludesTheEightNamedRunes`；共有 `TestDisplayWashRosterCardinality` 但**红在不同
  行**。⇒ 成对两发落点不同、⛔ 同名册，达标。
- 摘除对照（把基数钉实际侧改成期望侧的自抄本，保留发二的产码突变）：基数钉翻回
  `PASS`。⇒ 常钉的牙来自读 `isDisplayWashCf`、⛔ 自指（与 R4 点名的 `:400`"钉测试件
  自己那份抄本"正相反），无需就此上裁。
- 逐发还原哈希＝`8050f22d14203bba006cbfabcfb2a748081286d492691109c3a35a7a435ff8fd`。

## 5. 门禁四数（只在共享工作树跑，交件前）

`gofumpt -l internal/panel/display_scrub_327r2_test.go` rc=0（零输出）／
`sh scripts/d22scan.sh` rc=0（clean，ban #8 internal/ 581 枚含注释与 _test.go）／
`sh scripts/check-path-length-budget.sh` rc=0（VERDICT GREEN）／
`go vet ./internal/panel/` rc=0。

## 6. 禁区遵守

⛔ 动判语件 R4 三处装饰断言（本波只追加，`:357`-`:360`／`:400`-`:401`／`:416`-`:417`
一字未碰）；⛔ 扩名单（15 枚外 155 枚 `Cf` 归票 332，未加一枚）；⛔ 动 `unicode.Mn`；
⛔ 动 `internal/risk/taintmatch.go` 的 `isIgnorableTaintRune`；⛔ 动三枚冻结件
（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`ticket90_persist_test.go`）。
写面＝`internal/panel/**`（一枚新测试件）＋`.scratch/wisp/probes/327/r2/**`＋本件＋
票 327 末尾一条只追加 bullet；⛔ 一切其它面零字节、⛔ push。

## 7. 残余状态

R1 的常钉已进仓内：静默加/减一枚名册字符，基数钉现在两头都响、且红句点名到枚。
R2（名册外 7 枚 `Cf` 要不要纳入）仍归票 332 的普查，本波⛔ 代裁。R4 三处装饰断言
本波按令⛔ 动，交编排者另裁。
