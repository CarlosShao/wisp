# 票 76 对抗验收 —— 把 artifacts 路的"隐含安全"变成跑过的不变式（并修掉一个真绕过）

**验收人**：编排者（非实现者；实现方 `agent-ticket76`，3 个 commit：`d9224af` `39420cf` `974c082`）
**验收时间**：2026-09-21 10:30 · **基线**：`afd6353` · **纯净树验收**（`git archive HEAD` 解到 `/tmp`）
**档位**：〔独立复现〕／〔日志＋归档，我抽验〕／〔仅自述，不背书〕

| AC | 裁决 | 判据 | 档位 |
|---|---|---|---|
| **AC#1** 四种恶意组件形状逐个被拒或净化 | **PASS** | 我的变异（下详）能红，即证明"拒/净化"是真的在做事；且 memory 侧断言的是 `errors.Is(err, ErrInvalidArtifactName)` **且不是** `ErrNotFound` ⇒ **拒在碰文件系统之前**，这个区分做对了才有意义。 | 〔独立复现〕 |
| **AC#2** containment 用真目录差分而非字符串比较 | **PASS（本票质量最高的一格）** | 两条**阳性对照**：合法删除必须在差分里恰好显示 `[data/artifacts/real-1.txt]`；手工做的**未净化** join 必须能在 `<root>\CONTROL-escape.txt` 被看见。⇒ 差分仪器本身被证明"能看见逃逸"，否则"没看见逃逸"零信息（A15 那一族的正确姿势）。 | 〔日志＋归档，我抽验〕（对照逻辑我读过，未重跑） |
| **AC#3** 变异检验 | **PASS（我重做了关键那一发）** | 我把 `artifacts.go:167` 的 `strings.TrimRight(name, ". ")` 换成 `name`（守卫永不触发），`sed -n '167p'` 证落盘 ⇒ **`EXIT=1`，`FAIL=2 / PASS=11`**，红的是 `TestDeleteArtifactRejectsTheFourHostileShapes` 与其 `dotdot/bare_and_empty` 子项；还原后 `grep -c TrimRight`=1、`ok 0.294s`。 | 〔独立复现〕 |
| **AC#4** 交出票 20 `:103` 的替换句 | **PASS** | 句子已收到并由**我**落进票面（见下），代理没动别人的票面，这是对的。 | 〔独立复现〕 |
| **AC#5** 门禁 | **PASS（含一处主动披露，加分）** | 全套餐 `-count=2`：**RUN 228 = 2×114 不同测试名**；scoped 8 条：**RUN 44 = 2×22、SKIP 0**。**它主动点名**日志里有 2 条既有的 `--- SKIP`，来自 `internal/memory/concurrent_test.go:186` 的 re-exec 辅助用例（设计上单独跑就 skip，真驱动是 `TestCrashRecoveryKillMidWrite`）——**是"说出来"而不是"过滤掉"**，正是我要的记账方式。 | 〔抽验：我未重跑 -count=2 全套餐〕 |

## 顺带修掉的真绕过（D-76a）—— 与今天 R17 是**同一条不变式**

`DeleteArtifact("....")` 在 Windows 上**就是 `.` 的一种拼写**（解析前会剥掉尾部的点和空格），
旧守卫只看 Go 的路径语义，于是它**真的走到了** `os.Remove(<artifactsDir>\....)`；
当时只是因为目录非空才没出事——**空目录时 artifacts 目录本身会被删掉**。
⇒ 这不是"名字校验不严"，而是**"拼写可以改变安全语义"**——与票 72 的 `RUNNER~1`（8.3 短名改变 A/B 分类）
**同一条不变式**：**判定不得依赖路径的拼写形式**。今天第三例（前两例：票 72 锚点、票 71 的
`test-windows` 真原因）。⇒ 归进 **A45**，并提升为一条**通用检查项**：见到任何"名字/路径比较"，
先问"Windows 会不会把它认成别的东西"（尾部点/空格、8.3、大小写、`\\?\`、trailing separator）。

## 它交回但故意没修的三条，我的处置（不留口头账）

1. **`artifactName` 会把不同 id 折成同一个磁盘名**（`p/q`、`p\q`、`pq` ⇒ 都变 `tool-output-pq.txt`），
   且 **`writeFileExclusive` 根本没用 `O_EXCL`**（名字在撒谎）。
   ⇒ **这是 provenance 破坏，不是外观问题**：模型以为能重读的那个工具结果，字节可能已被另一次调用换掉。
   **已开成票 79**（含"注入性/失败要响"两条路 + 让 `O_EXCL` 名实相符 + 同 id 重试的行为必须被测试钉住）。
2. **没人保证 `artifacts/` 是平的**：`listArtifactsDir` 的 `if e.IsDir() { continue }` 使游离子目录
   对 `ListArtifacts`/`PurgeArtifacts`/500MB 配额**全部隐形**（它的 fixture 里 `nested\` 真的活过了 purge）。
   ⇒ 同进**票 79** Defect 2。"静默忽略"是现状，而**它不能靠默认赢**：对配额隐形的条目就是 500MB 变 2GB 的路径。
3. **两条 artifacts 路仍不过 `risk.Resolve`/C26**（我的裁定，票 76 里已写）：因此若有人在
   `artifactsDir` **内部**放一个 reparse point，裸名 `DeleteArtifact` 会跟着走。
   ⇒ **保留为已登记的残余风险**，理由不变：不接受调用方路径 ⇒ 无调用方可控面；动 C26 是冻结区且零收益。
   **但这条的措辞要在 Q-17 落地后复查**——那时判定层有了可枚举原因，可能就不需要靠"没有入口"来兜。

## 结论：**PASS，5/5** ⇒ 归档 `-done`；票 20 的 `:103` 框同时被我改写并勾选（那句原框是照 API 并不存在的威胁模型写的）
