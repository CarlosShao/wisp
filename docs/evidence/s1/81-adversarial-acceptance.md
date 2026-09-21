# 票 81 — 对抗验收裁决表（编排者亲自复跑，2026-09-21 15:1x）

**被验对象**：两份 containment 判据的夹具与**阳性对照**是 Windows 形状 ⇒ ubuntu 上必红（A52③ 的 2 条）。
**被验 commit**：`522efec`（夹具改平台无关 + 保留字面 `\` 并钉不对称）、`4683c34`（Linux 首轮逼出的修正：
期望集合改派生 + 反-vacuous 卫兵）、`1ff5366`/`6a77116`（票面账）。
**判据的真正仪器不是我的笔记本机，是 CI**——那两条红本来就是在 ubuntu 上观测的（票面 ①② 的来源）。

## 决定性证据：run **`35562680354`**（headSha `e5e5eb7`，status `completed` ⇒ 真样本）

| 项 | 读数 |
|---|---|
| `test-core` 的 `--- FAIL` 全量去重 | **8 条**（修复前 47 → 上一版 10 → 本版 8） |
| 那 8 条是谁 | `TestWriteGate*` ×4 + `TestSync*` ×4 = **`internal/risk` 的 sync-root 一家族**（票 55/票 82 的账） |
| 本票那 2 条 | **`TestSpillContainmentByDirectoryListing`、`TestArtifactsContainmentByDirectoryListing` 已不在 FAIL 名单里** |

⇒ **AC#1/AC#2 在 CI 上闭**：不是我"本地跑过了"，是**同一台 ubuntu、同一条命令**不再报它们。

## 我自己复跑的数（Windows 本机，与代理同一过滤器）

```
go test -count=2 -v -run 'Containment|LiteralBackslash|HostileShapes|StaysUnderDataDir|RejectsTheFour' \
  ./internal/agent ./internal/memory
rc=0  ·  === RUN 50  ·  --- PASS 50  ·  --- FAIL 0  ·  --- SKIP 0  ·  "no tests to run" 出现 0 次
```
⇒ 与代理自报 `50/50/0/0（两侧同数）`**逐字一致**。它的 AC#4 两包全量：两侧各 `RUN 284 = 2 × 142`、
2 条 SKIP 都是既有的 `TestSubprocessCrashWriter`（同名同数，不是新增），`no tests to run` 五份日志 0 命中。

## 逐 AC 裁决

| AC | 裁决 | 证据档 |
|---|---|---|
| **AC#1** 平台无关表达 + **保留**字面 `\` 那一例并钉住两侧不对称；**不许加 build tag 藏起来** | **PASS** | 〔日志＋归档，我抽验〕两份文件 `grep -c 'go:build'` = 0（我核过），不对称对照命名成 `separator_backslash_literal` / `dotdot_backslash_literal` 两条显式子例 |
| **AC#2** 两侧各一次红→绿（MUT-A 去 `/` 转义 ⇒ FAIL 5；MUT-B 短路 `validArtifactName` ⇒ FAIL 16），且**变异只发生在仓外快照** | **PASS**（Linux 侧的最终裁定用 CI 替掉了本地 docker） | 〔独立复现我的侧：Windows 50/50；CI：那 2 条从 FAIL 名单消失〕 |
| **AC#3** 全仓扫同族并给**命中清单**（不许只说"扫了没问题"） | **PASS** | 〔抽验〕广谱 389 行/70 文件 → 门挡掉 13 文件 46 行 → 未挡 57 文件 343 行；**窄谱决定性仪器 11 行 → 真用字面 `\` 构造物理路径只有 5 行**，其中 2 行是本票保留的对照、3 行属票 75 的 `filepath.Separator` 分支（不是缺陷）。**它把"大头是正则转义和只当数据用的 Windows 路径"逐处点名**，这就是判据要的形态 |
| **AC#4** 门禁两侧同数 + 四种假绿点名 | **PASS** | 〔独立复现〕见上面的 50/50 |

## 我要落进账里的三条**通用**收获（不属本票但由它产出）

1. **期望集合要由规则派生，不要硬编码字面名**。它把 purge 的期望从"点名 4 个文件"改成
   "按 target 落在 artifacts 树里"派生——**点名的期望值等于把测试作者的机器形状当成规格**，
   而这正是本票要修的病。并且留了**两条反-vacuous 卫兵**防止派生逻辑退化成"什么都算对"。
   ⇒ 固化成判据：凡"列举目录后比对集合"的用例照此办。
2. **命名债我判"不改名、改注释"**：`TestDelete{Artifact,PrivacyItem}RejectsTheFourHostileShapes` 现在跑 6 个子测试，
   函数名仍叫 "Four"。代理没改是对的——那名字被 `76-adversarial-acceptance.md`、本 registry、票 20/76/81 **逐字引用**，
   **改名 = 替别人重写验收账**。处置：下次碰该文件时补一行注释说明 "Four" 指哪四种，多出的两例是不对称对照。
3. **它登记了一处潜伏的空仪器**：`internal/risk/rules_test.go:198` 用字面 `C:\Program Files\Git\bin\git.exe`
   断言折叠行为，今天两侧同字节同结论；**一旦折叠语义收窄成"只在 Windows 折"，这条会在 ubuntu 上静默空转**。
   ⇒ 已追加成 **票 82 的 AC#1b**（要两侧各自的命中数/读数，空转就地改派生期望值），
   并顺手给票 82 加了 **AC#2b**：分层必须用构建约束表达，因为 `*_windows_test.go` 这个**文件名什么门都不施加**，
   只改名会让 `go test ./internal/risk/` 打 `ok [no test files]` ——看着全绿、整包判据消失。

## 结论

**票 81 转 `-done`。** 全程零 build tag、零调阈值、变异只在仓外快照、工作树一次没脏（当场 `git diff --quiet` 证）。
⚠ 一处它如实报的自己的弯路：第一次变异把 `\` 塞进 sed 替换写成非法字面量 ⇒ **编译失败，它作废重跑了**
（本仓"编译失败不算行为变异"这条规矩被代理自己用上，是第二次）。
