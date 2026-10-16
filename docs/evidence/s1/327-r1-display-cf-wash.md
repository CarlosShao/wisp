# 327-r1 交件自证：显示边界的洗集从 Cc 扩到裁好的那批 Cf

腿 `327-r1`（落地腿＝写产码）｜票 `.scratch/wisp/issues/327-display-boundary-washes-cc-only-so-cf-invisible-chars-still-pass-through.md`
起手锚＝具名提交 `c315980b`｜2026-10-16｜本件是**实现者自证**，⛔ 裁决（缺口审计／对抗验收归非实现者腿）
逐字读数（台面、命令、四数、两发突变名册、还原哈希）＝`.scratch/wisp/probes/327/r1/logs/readings.txt`

## 1. 落了什么（写面四枚，逐枚提交号见 §5）

`internal/panel/display_scrub_216.go`：新增一枚名单函数 `isDisplayWashCf`，**写入侧 `scrubForDisplay`
的 switch 与判据侧 `needsDisplayWash` 各调这同一枚函数**（AC#2"两把尺不许分家"的形状＝一枚名单两处调用，
不是两份并抄）。洗出来的形沿用既有 `writeControlEscape` 词汇：`U+00AD` 落 `\xad`（`r <= 0xff` 那支），
其余 14 枚落 `\uXXXX` 小写十六进制（`r > 0xff` 那支）⇒ 界面**仍只有一套转义格式**，⛔ 新造第二套。
⛔ 删字符：`abc`＋`U+202E`＋`def` 删完＝`abcdef`，一枚没人起过的名字；转义保留每一个别的字节在原位。

件头"故意不做什么"那一节里，`216-r1` 写的"does not wash Cf"三条已改写为现行事实＋**具名留下裁定没点名的 7 枚**。
⛔ 动 `unicode.Mn`（组合重音是合法可见字形）；⛔ 动 `internal/risk/taintmatch.go` 的 `isIgnorableTaintRune`
（判级输入轴，与显示轴取向相反：判定要"忽略后还认得"，显示要"显出来让人看见"）；名单旁边写了这句差异。

## 2. 名册逐枚对齐（⛔ "一致"两个字，逐枚列）

`needsDisplayWash` 认的集＝`Cc`（`unicode.IsControl`）＋下表 15 枚＋非法 UTF-8 字节。
`scrubForDisplay` 洗的集＝`Cc` ＋下表同一枚 `isDisplayWashCf` 的 15 枚＋非法 UTF-8 字节。两侧**逐枚同一枚来源**：

| # | 码点 | 名 | 洗成的形 | 族 |
|---|---|---|---|---|
| 1 | `U+00AD` | soft hyphen | `\xad` | 隐形排版 |
| 2 | `U+200B` | zero width space | `\u200b` | 隐形排版 |
| 3 | `U+200C` | zero width non-joiner | `\u200c` | 隐形排版 |
| 4 | `U+200D` | zero width joiner | `\u200d` | 隐形排版 |
| 5 | `U+2060` | word joiner | `\u2060` | 隐形排版 |
| 6 | `U+FEFF` | byte order mark / ZWNBSP | `\ufeff` | 隐形排版 |
| 7 | `U+202A` | left-to-right embedding | `\u202a` | 方向控制 |
| 8 | `U+202B` | right-to-left embedding | `\u202b` | 方向控制 |
| 9 | `U+202C` | pop directional formatting | `\u202c` | 方向控制 |
| 10 | `U+202D` | left-to-right override | `\u202d` | 方向控制 |
| 11 | `U+202E` | right-to-left override | `\u202e` | 方向控制 |
| 12 | `U+2066` | left-to-right isolate | `\u2066` | 方向控制 |
| 13 | `U+2067` | right-to-left isolate | `\u2067` | 方向控制 |
| 14 | `U+2068` | first strong isolate | `\u2068` | 方向控制 |
| 15 | `U+2069` | pop directional isolate | `\u2069` | 方向控制 |

九枚方向控制＋六枚隐形排版＝**15 枚**。⚠ **具名交回的一处不一致**：派单头写"**共 22 枚**"，
但它逐枚列出的正是上面这 15 枚（也与票 `AC#1` 的"建议起点＝9 枚方向＋6 枚零宽与软连字符"逐字相同）。
`22` 那个数出自普查件 `census-2` 的**候选**表：候选 22 ＝ 15 ＋裁定**没点名**的 7 枚
（`U+2061`–`U+2064` 数学作用符 4 枚＋`U+061C`＋`U+180E`＋`U+206A`）。
本腿按枚举实现（⛔ 自扩＝派单第一条禁令），并把那 7 枚钉成用例 `TestTheRosterIsClosedAgainstTheRunesTheRulingLeftOut`
（逐枚断言"判据不响、写入侧不改一字"），让"要不要纳入"必须是一纸裁定而不是一次顺手编辑。
**差值请编排者拍**：22 这个标签要不要更正，或那 7 枚要不要纳入（纳入即改名单＝契约级改动）。

## 3. AC#2 成对定向突变（两发、红在不同处）

| 发 | 摘掉哪一支 | 突变件 sha256（前缀） | RUN/PASS/FAIL/SKIP | 比终态多出的红 |
|---|---|---|---|---|
| 一 | 只摘写入侧：`scrubForDisplay` 的 `case isDisplayWashCf(r):` 7 行支（判据留着） | `e3bb5a1c…` | 252／150／12／0 | `TestTheWriterBranchWashesCfEvenWhenTheFastPathAlreadyFired`（断言行 `:352`＋`:360`）、`TestEveryNamedCfRuneIsEscapedAndNeverDeleted`（`:407`）、`TestABidiNameWalksTheRealBoundaryIntoAWashedRow`（`:477`）、`TestTheBidiClaimInANameChangesNoVerdict`（`:534`） |
| 二 | 只摘判据侧：`needsDisplayWash` 里 `|| isDisplayWashCf(r)` 一枚 term（写入侧留着） | `579e02d5…` | 252／150／12／0 | `TestNeedsDisplayWashNamesEveryWashedCfRune`（断言行 `:368`，逐枚 15 行）、`TestEveryNamedCfRune…`、`TestABidiNameWalks…`、`TestTheBidiClaim…` |

名册作差：**只在发一红**＝`TestTheWriterBranch…`；**只在发二红**＝`TestNeedsDisplayWashNames…` ⇒ 两发红在
**不同枚用例、不同枚断言行**（发一红在"洗出来必须是可见转义"那一格，发二红在"快路径必须先认得要洗"那一格），
⛔ 一发"整支去牙"冒充两发、⛔ 同名册。两枚靶子用例的设计前提也写在件里：发一的输入另带一枚 ESC，
使判据无论如何都响，所以它**只可能因写入侧去牙而红**；发二只调 `needsDisplayWash`，**只可能因判据侧去牙而红**。
还原：两发做完各自把工作树件拷回导出台，`sha256sum` 逐字回到 `8050f22d14203bba006cbfabcfb2a748081286d492691109c3a35a7a435ff8fd`
（＝改后终态件头哈希，两发同串）。突变**只在仓外导出台副本**（`D:\tmp\327r1\m1`／`\m2`）做，工作树与 `internal/risk/**` 零字节。

## 4. 颜色与门禁

台面＝仓外导出树八件套（`.go` 枚数 747，取出即数）；`-count=1 -timeout 30m`；⛔ 整包 `./cmd/wisp`、⛔ 全树。
`./internal/panel` 改前 `RUN 246／PASS 148／FAIL 8／SKIP 0` ⇄ 改后 `RUN 252／PASS 154／FAIL 8／SKIP 0`；
红名册 `diff` 作差＝**0**（8 枚起手红逐枚同名；其中 `TestC21DesignTokensFourWayAgree` 的红句是
`open …\bench\design\assets\tokens.css: cannot find the path`＝八件套里没有 `design/` 的台面射程，
⛔ 计入判据，具名于 readings.txt [2]）。RUN/PASS 各＋6＝本腿新增的 6 枚顶层用例全绿。
门禁四数（只在工作树跑）：`gofumpt -l`（两枚件）rc=0 零输出／`sh scripts/d22scan.sh` rc=0 `clean - no D22 ban violations`／
`sh scripts/check-path-length-budget.sh --with-self-test` rc=0 `VERDICT GREEN`（tracked=14393 over-budget=60 covered=60 not-in-roster=0）／
`go vet ./internal/panel/` rc=0。界面文案与用例字符串零 emoji（新增串只含 ASCII 转义与汉字）。

## 5. 提交与残余（交回编排者）

提交：`c2f9310b`（产码件＋msg 存档）／`0eaf5a6a`（用例件＋msg 存档）／本件所在提交（读＋证据）／票尾 bullet。
`git status --porcelain -- internal cmd` 交件时为空（见交件报告）。

1. **票 `AC#4` 的"攻它的那一发"本腿未做**：那发要求把判定改成"名字里含'允许'就放行"再看用例红，
   射程＝`internal/risk/**`／`internal/perm/**`＝本腿禁区（派单原文"⛔ 为过这格去动 perm／risk 任何判定代码"）。
   本腿交的是同形的**正向半格**：真 `risk.NewRiskAssessor()`，三名（`photo.png`／`photo<U+202E>已允许.png`／
   洗过的）判级与命中规则逐字同形＝`L2/[R1 R8]`，且界面拿到的是 `\u202e` 那形。攻它那一发请派给有 risk 写面的腿。
2. **名册外的 7 枚 `Cf` 照旧穿过**（§2 末），已钉为用例；要不要纳入＝一纸裁定。
3. **派单头"共 22 枚"与逐枚枚举（15 枚）不符**，本腿按枚举实现并具名上报，⛔ 自行更正票面或派单。
4. 普查件 §① 的 7 枚"出向且没过这道边界"的字段（`pending[].tool/args/reason`、`results[].text`、
   `composer.workspace.reason`、`composer.attachmentError`、`tasks.rows[].label`、`instructions.files[].path`、`git.*`）
   本票射程外（票 `AC#1` 名册＝`Workspace.Canonical`＋`Attachments[].Name` 两枚），⛔ 顺手扩，
   其中 `composer.workspace.spelling` 是 `composer.go` 自己声明的"用户本人的串、逐字出向"，故那枚可见用例把
   `Spelling` 放干净、把 `U+202E` 放在用户没打的 `Canonical` 上（＝普查的入口 A 形状）。
5. 最后一公里（快照字节→页面）仍断在 Go 出向之后（普查件 §③ 环 8 的读数），本票只把**出向字节**洗干净。
