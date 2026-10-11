# 表②——12 枚逐枚机制归因（本腿自己的颜色为凭据）

判定汇总（**枚数＝名册行数＝12**）：**甲 7 枚／乙 5 枚／丙 0 枚（但有 1 枚带〔缺世界〕附注）／丁 0 枚**。

## 1. 本腿的两发颜色（起手就落盘，⛔ 引 `A817` 的读数替它）

| 发 | 命令（逐字） | 退码 | 耗时 | 四数（尺＝`grep -c '^=== RUN'` 等，射程＝整份件） | 件 |
|---|---|---|---|---|---|
| `T01` 本机台面 | `go test ./internal/risk/ -count=1 -v -run '^(…12 枚全名…)$'` | **rc=0**（同发现量） | `real 0m1.213s` | `RUN=12 PASS=12 FAIL=0 SKIP=0` | `logs/t01-local-12.txt` |
| `T02` **runner 形状复现** | `env TMPDIR='C:\Users\swq\tmp\RISKAT~1\RUNNER~1' TEMP=… TMP=… go test ./internal/risk/ -count=1 -v -run '<同 12 枚>'` | ⚠ **本腿这条件里 `rc=` 落空**（我在同一子 shell 里写了 `echo "rc=$rc"` 却⛔ 先 `rc=$?`，那格是空的——纪律自报，见 §5） | `real 0m1.213s` | `RUN=12 PASS=0 **FAIL=12** SKIP=0` ＋件尾逐字 `FAIL github.com/CarlosShao/wisp/internal/risk 0.249s` | `logs/t02-shorttemp-12.txt` |

⇒ **`T01`＝`A817` §2 那句"12 枚本机绿"在本腿手上复现成立**（⛔ 照抄，本腿自己跑的）。
⇒ **`T02`＝本腿造出的形状把 12 枚一次全打红**，与 CI 的 12↔12 枚枚同名 ⇒ **CI 那 12 枚红的成因在这台机器上可复现，⛔ 需要真 runner**。这一发改写了 `A817` §5"归口另待一枚普查腿定射程"的前提出手式：**归因⛔ 欠台面**。

★复现法＝票 115 `AC#3`/`AC#2` 那两件的做法，**换了注入层**：票 115 是在测试体内把"调用方自己的长形字段"换成同一棵树的 8.3 形（`shortFormOf115`，换不动即 `t.Fatal`），件＝`internal/winsec/notice_attribution_115_windows_test.go` ＋票面 `.scratch/wisp/issues/115-…-done.md` §②（`docs/evidence/s1/115-*.md` **不存在**，`find docs/evidence -name '115*'`＝0，与票 230 撞车那事票面已自陈）。
本腿**⛔ 动任何 `_test.go`**（授权面没有这一项），改成**进程内 env 覆盖**：把 `t.TempDir()` 的上游 `TMPDIR` 指到一枚**仓外**目录的 **8.3 别名**上。仓外种植＋别名读数＝`logs/shortname-probe.ps1` 的输出逐字：

```
LONG   = C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1
SHORT  = C:\Users\swq\tmp\RISKAT~1\RUNNER~1
```

⇒ 两级都取得到别名（`riskattrib1`→`RISKAT~1`、`RunnerTempAreaForAttribution1`→`RUNNER~1`），**别名末段与 runner 的 `RUNNER~1` 逐字同名**（纯属巧合，⛔ 是设计出来的；`runneradmin` 的 8.3 形就是 `RUNNER~1`）。
⛔ 改全局／用户级配置、⛔ 写进仓、⛔ 动 git config；临时件只建不删。
✅ 复核了派单给的本机事实：`C:\Users\swq\tmp` 这一支**取得到** 8.3 别名（与"TEMP 没有别名"不矛盾——那是 `%TEMP%` 自己那一段）。

## 2. 甲形（7 枚）——判级输入读错了树：期望侧写的是"这台机器的拼写"

**共同的形**：这 7 枚都把 `t.TempDir()` 的**返回拼写**当成"该树唯一的／真正的拼写"。runner 上 `%TEMP%` 逐字是 `C:\Users\RUNNER~1\AppData\Local\Temp`（＝`runneradmin` 的 8.3 形），而 C26 解析器的答案是**长形** ⇒ 输入拼写 ≠ 解析器答案 ⇒ 红在**比对／前提那一侧**，⛔ 红在判级本身。
★本腿用 `T02` 逐枚量到"错在哪一截"（红句逐字见 `logs/t02-shorttemp-12.txt`，射程＝`=== RUN`↔`--- FAIL` 之间）。

| # | 用例 | 红点（文件:行） | 本腿读到的逐字红句（截短） | 判定 | 凭据标签 |
|---|---|---|---|---|---|
| 1 | `TestClassifyAnchorSpellingIsNotVerdict` | `pathresolver_anchor_spelling_windows_test.go:57` | `Resolve("…\RISKAT~1\RUNNER~1\…\profile\.ssh\id_testkey") canonical "…\riskattrib1\RunnerTempAreaForAttribution1\…"`: **test premise broken (candidate must be the real path)** | **甲** | **已证**（`T02`；`T01` 绿） |
| 2 | `TestCanonicalInputGainsNoSecondForm` | `…anchor_spelling_windows_test.go:140` | `canonical "…\RISKAT~1\RUNNER~1\…" produced **2 comparison forms** […, …riskattrib1\runnertemparea…], want 1` | **甲** | **已证** |
| 3 | `TestAListWinsWhereBothTablesHit` | `…anchor_spelling_windows_test.go:238` | `B-only control must stay overridable by design: {Allow:false NeedL2:true Class:B …id_*}` ← **⛔ 路径串**（`%+v` 打的是 decision struct） | **甲** | **已证** |
| 4 | `TestPathResolverShortNameAListDenied` | `…junction_windows_test.go:135` | `short name did not expand to long path: got "…\riskattrib1\RunnerTempAreaForAttribution1\…" want "…\RISKAT~1\RUNNER~1\…"` ← 变量名叫 `long`，值却是**短形** | **甲** | **已证** |
| 5 | `TestPathResolverUNCAListDenied` | `…junction_windows_test.go:164` | `UNC spelling "\\?\UNC\localhost\c$\…\RISKAT~1\RUNNER~1\…" normalized to "…\riskattrib1\…", want "…\RISKAT~1\…"` | **甲** | **已证** |
| 6 | `TestPathResolverExtendedLengthPrefixAListDenied` | `…junction_windows_test.go:182` | `\\?\ spelling normalized to "…\riskattrib1\…", want "…\RISKAT~1\…"` | **甲** | **已证** |
| 7 | `TestBListDefaultDenyAndOverride` | `…junction_windows_test.go:292` | `B-list single-file override must allow: {Allow:false NeedL2:true Class:B …\.env*}` ← 同样⛔ 路径串 | **甲** | **已证** |

**枚数核对**：甲 **7** 行。
★**这 7 枚的实质判级在 `T02` 那一发里从没被证伪**：#1/#2/#3/#7 红在前提／上界／正控那一行，#4/#5/#6 红在 `want` 的拼写而非"没展开"（`got` 恰恰是**完全展开的长形**＝解析器干对了）。⇒ **⛔ 把这 7 枚读成"C26 在 runner 上判错"**。
⇒ 这也解释了 `A817` §3 那半句"2 枚红句不含那串"：**不是机制未证的旁证，是那两枚的 `t.Fatalf` 只 `%+v` 打了 decision struct、⛔ 打输入路径**（#3 `:238`／#7 `:292`）。本腿把它们归进甲形并给了逐字红句，`A817` 的"机制未证"这一格**已被本腿填掉**。

**票 115 先例的对位（⛔ 本腿新造规矩，逐字指回）**：票 115 `AC#2` 裁的是**方向 B＝比对按树、⛔ 按拼写**，且它把"要改的是用例的比对面"写死；#4/#5/#6 的 `want` 与 #1 的 premise 正是那种"拿调用方拼写当真值"的比对面。票 115 还留了另外两形本腿用得上：**(a)** 仪器腿要自证种得出（`the instrument measured nothing` ⇒ `t.Fatal`，⛔ Skip）；**(b)** `AC#4` 要求变异（退回原状／`EqualFold` 半修都要红）——本腿⛔ 有变异授权（⛔ 改产码），那一格归裁决者。

## 3. 乙形（5 枚）——真缺陷：sync 归属的比对两边归一层⛔ 对称

| # | 用例 | 红点 | 本腿读到的逐字红句 | 判定 | 凭据标签 |
|---|---|---|---|---|---|
| 8 | `TestSyncFixtureFallbackAndMatch` | `syncdirs_test.go:162` | `write under fixture root must be sync` | **乙** | **已证**（`T02` 红／`T01` 绿） |
| 9 | `TestSyncFallbackNotDisarmableByWeakRoot` | `syncdirs_test.go:213` | `source "default"/"fixture"/"options"/"": under-profile fallback must stay armed`（**四种 Source 同一形，枚枚红**） | **乙** | **已证** |
| 10 | `TestSyncUnverifiedRootKeepsFallback` | `syncdirs_test.go:233` | `suspect fallback must stay armed` | **乙** | **已证** |
| 11 | `TestSyncSuspectFallbackIsComponentBounded` | `syncdirs_test.go:355` | `under-profile must be suspect` | **乙** | **已证** |
| 12 | `TestSyncSuspectFallbackWhenUndetectable` | `syncdirs_test.go:402` | `under-profile path must be sync-suspect, got {Sync:false Root:{Provider: Path: Source:} Why:write target is not under any sync root}` | **乙** | **已证** |

**枚数核对**：乙 **5** 行。甲 7 ＋ 乙 5 ＝ **12＝名册行数**。

**错在哪一截（这五枚＝同一处，一次记一次）**：
- `internal/risk/syncdirs.go:88`（同形复制在 `internal/risk/provenance.go:345`）：`home: normPath(orDefault(o.HomeDir, userHomeDir()))`——`home` 只过**词法**归一（`blacklist.go:133 normPath`＝统一小写＋换分隔符，`blacklist.go:169 isUnder`＝归一后的**裸前缀比较**），**从不过 C26 解析器**。
- `internal/risk/syncdirs.go:408` 起：`cand, err := s.resolveTarget(rawPath)`——候选侧**过解析器**（＝8.3 会被展开成长形）。
- `internal/risk/syncdirs.go:429`：`if !s.complete && s.home != "" && isUnder(cand, s.home)` ⇒ **一边展开、一边不展开** ⇒ 判级取决于"profile 是用哪一种拼写交到手上"。根成员那一支同形（`:425 isUnder(cand, e.canon)`）。
- ⇒ 红句 #12 的 `{Sync:false …}` 是本腿**手上的 fail-open 读数**：`Sync:true` 在这条链上是**受限**的那一侧（代码自陈的 fail-closed 分支返回的就是 `Sync:true`），所以归属丢失＝**该被当外逃通道看的写被放行**。

**今天有没有生产调用者（这把尺＝`grep -rn --include=*.go "IsSyncPath(" . | grep -v _test.go`，射程＝仓根全树，命中 21 行、剥掉 `.scratch/wisp/probes/**` 那 18 行变异快照后＝2 行）**：
- `internal/risk/provenance.go:928` → `if p.IsSyncPath(t).Sync`——**真生产调用点（在 `Inspect` 的 R4 通道判定里）**；
- `internal/risk/syncdirs.go:187`——定义本身。
⇒ **接了**，⛔"建了没接"。另：`SyncDetectionComplete(` 非测试面命中＝**0 枚调用者**（只有 `syncdirs.go:458` 定义）⇒ 这五枚里那些 `SyncDetectionComplete()` 断言是**只有测试在看的面**（本腿没见它们红，但归口时值得知道）。

**生产可达性的诚实边界〔仅读码推〕**：`%USERPROFILE%` 在生产里通常是长形，候选路径以 8.3 形交进来时会被解析成长形 ⇒ 与长形 `home` 相合、**⛔ 触发**；要触发需要 `home`（或注入的根）那一侧是短形而候选是长形——注册表／客户端配置／用户手填的根**正是这种"别人给的拼写"**，而候选侧永远被解析器规整过。本腿⛔ 有生产实例读数和变异授权，所以这一格**⛔ 写"已在撞"**，只写"形状上今天也撞得到，且撞上是 fail-open"。
★派单给的既有事实（生产码不调 `GetShortPathNameW`）与本腿一致：缺陷⛔ 在"枚举别名"那一侧，而在**比较两侧归一层不对称**这一侧。

## 4. 丙／缺世界／丁

- **丙＝0 枚**。本腿⛔ 有变异授权（⛔ 改产码／⛔ 改 `_test.go`），所以"判据换成反形它也不响"那一把尺**本腿没打**——这一格是**欠的读数**，⛔ 由名册枚数假装闭掉（归裁决者，先例＝票 115 `AC#4` 的 M1/M2/M3 三发）。
- 但本腿能给的两条**仪器形状**结论（⛔ 是"它今天红了"替它背书）：
  1. 〔缺世界〕**#2 `TestCanonicalInputGainsNoSecondForm` 在 `T01` 那发的绿是"没东西可测"的绿**：`formsOf` 只有在某段祖先的 8.3 别名／解析答案与输入拼写**不同**时才产出第二种形，而这台机器的 `%TEMP%` 那一段没有别名 ⇒ 本机那一发里"只有一种拼写"这件事**是世界的性质、⛔ 是守卫的性质**。它的名目（"爆炸半径上界"）在**任何一台 profile 有别名的机器上都会自己红掉**（＝`T02`、＝CI）。⇒ 修法形状：把它要的"第二种形"**在夹具里种出来**（票 115 那形：种不出即 `t.Fatal`），⛔ 让绿依赖卷设置。**这一条＝派单点名的"缺世界"那一形，本腿在 12 枚里只认出这一枚属它。**
  2. 〔仅读码推〕**#11 `TestSyncSuspectFallbackIsComponentBounded` 的主题（N-5 同名前缀兄弟目录⛔ 被扫进）在 runner 那一发上根本没被执行**：`:355` 的 `t.Fatal` 先响，`:357` 那句永远跑⛔ 到 ⇒ CI 的这枚红**报的⛔ 它标题写的那件事**。它的牙＝`isUnder` 用的是 `dir+sepStr` 而非裸 `HasPrefix`（`blacklist.go:169-171`，读码可见换成就红），⛔ 本腿跑过。
  3. 〔仅读码推〕#9 的判别轴是 `Source`（四种档），而 `T02` 那一发四种档**同一形全红** ⇒ 那一发上它分辨⛔ 出档级、只分辨得出拼写。
- **另记一枚仪器逃逸口**：#4 `TestPathResolverShortNameAListDenied:129` 在卷不生短名时 `t.Skip`（票 18 授权）。⇒ 它是这 12 枚里唯一一枚**会安静地不再测量**的；票 115 的定式是"种不出即 `t.Fatal`、⛔ Skip"。本腿两发里它都没跳（`SKIP=0`），所以那条 Skip 眼下没被走过——归口时要写进去，⛔ 顺手把它改成跳。
- **丁＝0 枚**：本腿⛔ 需要真 OneDrive／真账号／真 profile 的欠账——**12 枚全在这台机器上被判死**。（同族里唯一要台面的是 `TestSyncRedTeamRealOneDrive`，它⛔ 在这 12 枚内，CI 上是 `--- SKIP`。）

## 5. 纪律自报（违过的，具名）

1. **`T02` 那一发的 `rc=` 落空**：我写的是 `… > 件 2>&1; echo "rc=$rc"`，**漏了 `rc=$?`** 那一拍 ⇒ 那格空的。退码证据改由件尾逐字 `FAIL` ＋ `--- FAIL=12 / PASS=0` 承担；`T01` 的 `rc=0` 是同发现量、⛔ 受影响。**定式：下次那两拍必须写全。**
2. **第 1 笔 commit 我试过 `git commit -- <未跟踪路径>`**（退码 1、pathspec 未匹配）⇒ 改成显式两枚路径 `git add` 再 commit；⛔ `add -A`／`.`。⚠ 另一处：`00-anchor.md` 里我先把 raw log 写成 `logs/gate-00.txt`，而件里那句"落笔前现量"引用的名字一开始对不上（已按盘上真名 `gate-00.txt` 走），⛔ 删、⛔ 改名。
3. ⛔ 改任何 `_test.go`／产码、⛔ 新建 `.go`、⛔ 翻任何 `AC` 框、⛔ 判"本格闭合"、⛔ push、⛔ 跑 `./cmd/wisp/`、⛔ `d22scan`／`gofumpt`／`gofmt -l`、⛔ 放宽断言、⛔ 造 Skip。临时件全在 `C:/Users/swq/tmp/riskattrib1/`（仓外，只建不删）与本目录 `logs/`。
