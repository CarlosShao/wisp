# 表③——归口建议（只给形状，⛔ 动手、⛔ 新造票号）

★两条硬约束本表都守：**⛔ 塞进票 302**（`A817` §5 原话——那一族的修法在"档"，这 12 枚的修法在"判据读的是哪棵树"，混一票必然顺手改错那半）；**同一物理缺陷只在一个链上记一次**（先例＝`A817`／票 302 那处的记账纪律）。

本腿⛔ 建议立任何新票：**两族各自已有 open 票名下**（下表点名的两枚票面现在都**没有** `-done` 后缀，盘上尺＝`ls .scratch/wisp/issues/ | grep -E "^72|^252"`＝各 1 枚）。立不立、扩不扩射程由编排者裁。

## 1. 甲族 7 枚 → **该挂**（挂到票 72，它名下的用例、它自己欠的那半边）

| 要 | 内容 |
|---|---|
| 归口 | **挂**：`.scratch/wisp/issues/72-atble-classification-runner.md`（Status 现量＝`review`，票面逐字「**AC#4 的 runner 半边与 AC#6 的改名都不在代理手里**」） |
| 枚 | 名册 #1–#7（`TestClassifyAnchorSpellingIsNotVerdict`／`TestCanonicalInputGainsNoSecondForm`／`TestAListWinsWhereBothTablesHit`／`TestPathResolverShortNameAListDenied`／`TestPathResolverUNCAListDenied`／`TestPathResolverExtendedLengthPrefixAListDenied`／`TestBListDefaultDenyAndOverride`）＝**7 枚＝表② §2 行数** |
| 凭据 | `logs/t02-shorttemp-12.txt`：仓外 8.3 形状一发把 7 枚全打红，逐字红句与 CI 那 5＋2 枚同点同形 |
| 为什么是挂⛔ 立 | 那 7 枚用例**就是票 72 写的**（`pathresolver_anchor_spelling_windows_test.go:12` 逐字「Ticket 72: the A-tier verdict must not depend on the spelling a path arrived in (R17 invariant)」），而它 `AC#4` 欠的正是"runner 半边"。另立一票＝同一枚物理缺陷在两条链上各记一次 |
| 修法**射程**（只给形状） | 只动**比对面／前提面**：`anchor_spelling:57`（premise 那句"candidate must be the real path"要改成"要么相等、要么是同一棵树的另一种拼法"）、`:139`（"want 1 form"要按"这棵树有没有第二种合法拼法"来写，⛔ 按本机卷设置）、`:238` 与 `junction:292`（override 键的两侧要同层归一，见下表"顺带"一栏）、`junction:135/:164/:182` 的 `want`（变量名 `long` 装的必须是解析器答案，⛔ `t.TempDir()` 的返回拼写）。**判据极性、期望数、阈值一字节⛔ 动** |
| 先例做法 | 票 115 `AC#2` 方向 B（按树⛔ 按拼写）＋ `AC#3` 的仪器腿（种不出即 `t.Fatal`）。本腿的注入层不同：本腿⛔ 授权改 `_test.go`，只用 env `TMPDIR` 在仓外复现 |

## 2. 乙族 5 枚 → **该挂**（挂到票 252，那是同一物理缺陷的第一枚调用点）

| 要 | 内容 |
|---|---|
| 归口 | **挂**：`.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-…-r2-l2.md`（标题逐字「**同一条路径的两种拼法在允许根判定里分家：短名（8.3）侧留短、长名侧折长 ⇒ 允许根内的新建操作被判成越界**」） |
| 枚 | 名册 #8–#12（`TestSyncFixtureFallbackAndMatch`／`TestSyncFallbackNotDisarmableByWeakRoot`／`TestSyncUnverifiedRootKeepsFallback`／`TestSyncSuspectFallbackIsComponentBounded`／`TestSyncSuspectFallbackWhenUndetectable`）＝**5 枚＝表② §3 行数** |
| 凭据 | `logs/t02-shorttemp-12.txt` 五枚红点逐字＝`syncdirs_test.go:162/:213/:233/:355/:402`（与 `A817` §3 那五个行号**本腿独立复跑对上**）；其中 `:402` 逐字 `{Sync:false … Why:write target is not under any sync root}` ＝**fail-open 读数** |
| 为什么是挂⛔ 立 | 票 252 的票面 `:1`–`:12` 已经在写"一侧折长、一侧留短 ⇒ 两次包含不可能同时成立"，只是它的调用点是 `internal/tools/paths.go InAllowlist`。**同一缺陷在 sync 归属网上的第二枚调用点** ⇒ 挂进去＝把它的射程从 allowlist 扩到 `syncdirs`，⛔ 让世界上出现两张管同一件事的票 |
| 修法**射程**（只给形状） | `internal/risk/syncdirs.go:88`（`home: normPath(...)`——`home` 侧从不过 C26）与同形复制 `provenance.go:345`；`syncdirs.go:429 isUnder(cand, s.home)`／`:425 isUnder(cand, e.canon)` 两侧**必须同层归一**（候选侧 `:408 resolveTarget` 已折长）。参照体就在同一包里：`pathresolver.go:345 formsOf`／`:397 anchorForms` 是"两边都过解析器再比"的既有做法（票 115 票面点过的 `pathForms` 先例） |
| 今天有没有生产调用者 | **有**：`internal/risk/provenance.go:928 if p.IsSyncPath(t).Sync`（尺＝`grep -rn --include=*.go "IsSyncPath(" . \| grep -v _test.go`，全树 21 命中、剥掉 `.scratch/wisp/probes/**` 的 18 枚变异快照＝2 行：1 枚调用＋1 枚定义）。⛔"建了没接"这一枚**接上了** |
| ⚠ 挂的时候必须带上两条 | ① 票 252 的定性句「⚠ **这一枚算产品行为**，不是纯仪器」对 sync 这一侧**同样成立且更硬**（allowlist 侧偏响＝多弹一张卡；sync 侧偏松＝外逃通道不被标）。② `SyncDetectionComplete(` 非测试面**0 枚调用者**（只有定义 `syncdirs.go:458`）——那一面⛔ 有人吃，别顺手把它算进影响面 |

## 3. 仪器形状三条 → **该登记**（登记为"修法必须带的三条"，⛔ 独立一票）

这三条⛔ 是新的红、也⛔ 是新的票；它们是"谁来替这一格管"的账。归口＝并进上面两票的落地判据，**由裁决者（⛔ 实现者）那把变异尺管**（先例＝票 115 `AC#4` 的 M1/M2/M3 三发）。

1. **〔缺世界〕#2 `TestCanonicalInputGainsNoSecondForm` 的本机绿是"没东西可测"的绿**：本机 `%TEMP%` 那一段没别名 ⇒ "只有一种拼写"是**世界的性质、⛔ 守卫的性质**；它要的"第二种合法拼法"必须在夹具里**种出来**（票 115 形：种不出即 `t.Fatal`）。替它管的格＝票 72 落地票的变异那一格。
2. **逃逸口：#4 `TestPathResolverShortNameAListDenied`（`junction:129`）在卷不生短名时 `t.Skip`**——12 枚里唯一会"安静地不再测量"的那一枚。本腿两发都⛔ 走过它（`SKIP=0`）。⇒ 修法⛔ 把它扩成 Skip；要动就照票 115 定式动成 Fatal。**归编排者／机主裁**（改 Skip 语义＝改判据形状）。
3. **#11 `TestSyncSuspectFallbackIsComponentBounded` 的标题那件事（N-5 同名前缀兄弟目录）在 runner 那一发上根本没被执行**：`:355` 先 `t.Fatal`，`:357` 永远跑⛔ 到 ⇒ CI 报的⛔ 它写的那件事。⇒ 乙族落地后必须**单独确认那一句重新被执行**（凭据＝它自己在 `--- PASS` 列里，⛔ 拿整包 `ok` 顶）。
   ★**丙＝0 枚**：本腿⛔ 变异授权（⛔ 改产码／⛔ 改 `_test.go`），"判据换成反形它也不响"那一把尺本腿**没打**——那是**欠的读数**，归裁决者，⛔ 由本表的枚数假装闭掉。
   ★**丁＝0 枚**：12 枚全在本机判死，⛔ 要真 OneDrive／真账号／真 profile 的欠账。（同族里唯一要台面的是 `TestSyncRedTeamRealOneDrive`，它⛔ 在这 12 枚内、CI 上是 `--- SKIP`。）

## 4. 编译面差一条（⛔ 是归口，是给修法的边界）

`syncdirs_test.go` **无 build tag** ⇒ 那 5 枚在 POSIX／`--scope=core` 也编译也求值（另 7 枚在 `//go:build windows` 文件里）。
⇒ 修法⛔ 能长成"给这 5 枚加 windows tag"：那会把它们从 core 档分母里一起摘掉，撞上票 301（`301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md`，**guard D 是包级**）。这条写进乙票的"⛔ 这么做"栏。

## 5. 反向那 1 枚（派单要我复核并判"它是不是同一族"的那枚）

`TestC21TableColourRowsMatchTokensCSS`＝本机绿／runner 红的**反方向**那 1 枚（`A817` §2）。

- 归属现量：`internal/ball/tokens_table_test.go:1465` ⇒ **⛔ 在 `internal/risk`**，⛔ 在我车道（我的 Go 面只许 `./internal/risk/...`）⇒ **本腿⛔ 跑它**，只读判族。
- 判据面：`tokens_table_test.go:1384 data, err := os.ReadFile(path)`，常量 `c21TokensCSSPath = "design/assets/tokens.css"`（`:1364`）、另一处 `:142 os.ReadFile(filepath.Join(root, …, c21TablePath))` ⇒ 它读的是**工作树上的资产字节**。
- 本机工作树现量（只读）：`git status --porcelain -- design/assets/tokens.css`＝` D`，且 `git status --porcelain -- design/assets/` 报 `warning: could not open directory 'design/assets/': No such file or directory` ⇒ 整支 `design/assets/` 被别家腿在共享工作树里删掉了（⛔ 碰、⛔ 抱怨，只记为"这枚在本机也测不出真话"）。
- **判语：它⛔ 属于这 12 枚那一族（不是同一成因、不是同一棵树）。** 这 12 枚的成因＝**同一棵树的两种拼法**；那枚的成因＝**判据读的是工作树资产面**（CI 上是 checkout 里的 `design/**` 内容，本机现在是"目录都没了"）。它该归的是**资产／放置那一族**（派单点的先例＝`frontend/dist` 那族），**既⛔ 挂票 72、也⛔ 挂票 252**。CI 侧那枚到底红在哪一格的成因本腿**⛔ 判**（要跑 `internal/ball`、⛔ 在我车道）⇒ **归编排者或机主那一发**。
