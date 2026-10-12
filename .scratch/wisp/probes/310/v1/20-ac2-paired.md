# `20-ac2-paired.md` — `310-v1` 对票 310 `AC#2`（反恒真成对两发）的复核＋我自己复跑的那一发

**判语一句**：`AC#2` **成立**——"死亡那一形当场漏一枚、那一步在位时顶层不多、残渣进专用根"这条**被两条独立来路支持**（腿的成对四发＋我此刻自己复跑的两发）；⚠ 但腿件里**一枚具名目录（`TestRoundTripLoadMarshalLoad2236627062`）我从盘上取不到它任何来路**，那一格的"名字"⛔ 成立、"枚数"与"形"**成立**，且正典 §1 那条 `CreationTime` 定式**够硬**（下面 §3 给我自己那把尺的逐字证据）。

台件＝`ac2-rerun-310v1.sh`（复跑）＋`roster-310v1.ps1`（名册尺，计数在 PS Inside 算、逐枚带 `CreationTime`）。读数主件＝`logs/20-ac2-run.txt`（27 行全逐字）、`logs/20-roster-*.txt[.meta]`、`logs/20-a2-new.txt`、`logs/20-b-new.txt`、`logs/20-root-contents.txt`。

---

## 1. ★我自己复跑的那一发（⛔ 抄 `20-paired-readings.md`，一行都⛔ 抄）

台面＝**仓外导出树**（⛔ 在共享工作树里跑"改前那发"）：
```
export_tree=/d/tmp/310v1-tree.UcNNmC
archive_rc=0   go_files_in_tree=682   six_components_present=6 (must be 6)
```
⇒ `git archive HEAD go.mod go.sum internal cmd tools scripts | tar -x -C "$T"` **六件套一枚⛔ 少**，跑之前先数 `.go` 证非空＝**682** 枚（⛔ "取空了所以 0"那一形）。窄用例＝`internal/config` 的 `TestRoundTripLoadMarshalLoad`。
⚠ **我这一法的射程要说清**（否则是Over-claim）：两发都**直跑 `go test`、由我自己控 `GOTMPDIR`**（`env -u GOTMPDIR` ⇄ `GOTMPDIR=<根>`），**测的是"那一步的机制"**；腿那一步在入口件里做的事⛔ 是别的（只是 `export GOTMPDIR`，`10-impl.md` §1 行号已由我在 §AC#1② 逐字节复量）。"步子⛔ 在位"这一侧另有两发**真走入口件**的读数替它兜着：`10-ac1-verdicts.md` §①（5 枚带 `GOTMPDIR=` 的入口发）＋我自己 `30-ac3-run.txt` 那两发 census/winsec（同批根名在盘上、`TempDir RemoveAll cleanup`＝0 行）。

### (a) 撤掉那一步 ⇒ `%TEMP%` 顶层当场多一枚

```
roster a-before rc=0 names=58 meta_count=count=58 err_bytes=0
a2_shot_rc=1
a2_evidence_line=panic: test timed out after 2s
roster a-after rc=0 names=59 meta_count=count=59 err_bytes=0
a2_top_level_new=1 names_in_diff=TestRoundTripLoadMarshalLoad4222770164
```
同行证据（"那一形真发生了"，⛔ 只是枚数）＝`logs/20-shot-a2-test.txt` 第 2～4 行逐字：
```
panic: test timed out after 2s
	running tests:
		TestRoundTripLoadMarshalLoad (0s)
```
＋第 37 行栈帧 `github.com/CarlosShao/wisp/internal/config.TestRoundTripLoadMarshalLoad(0x20665139e600)` ⇒ 帽落在**用例正在跑**的时候（腿的 a1 那发作废就废在帽落在第一次 `t.TempDir()` 之前，本腿先暖了台面：`warm_rc=0 warm_tail=1`）。
★**这枚名字自带 `CreationTime` 佐证、且今天仍在盘上**（两条来路同一时刻）：
```
round=final temp_root=C:\Users\swq\AppData\Local\Temp\ filter=Test* count=59 ...
  roundtrip_name=TestRoundTripLoadMarshalLoad4222770164 created=2026-10-12 12:48:51.544
in_final_roster=1
```
**差集来路**＝`comm -13` 两枚名册件（⛔ 台件自报字段），两枚输入名册一起落盘（`logs/20-roster-a2-before.txt`／`-after.txt`）⇒ **正典 §1 那条定式我照做了**。**我这一程造的目录＝这一枚＋下面 (b) 那一枚在根里的，一枚⛔ 删。**

### (b) 还原（＝那一步在位）⇒ 顶层不多、残渣逐字在专用根里

```
b_root=/c/Users/swq/AppData/Local/wisp-test-gotmp/v1-20261012-124853-148 mkdir_rc=0
roster b-before rc=0 names=59 ...  roster b-after rc=0 names=59 ...
b_shot_rc=1   b_evidence_line=panic: test timed out after 2s
b_top_level_new=0 empty above = 0 枚 means the OS temp root gained nothing
b_root_entries=1 b_root_names=TestRoundTripLoadMarshalLoad3712892429
  /c/Users/swq/AppData/Local/wisp-test-gotmp/v1-20261012-124853-148/TestRoundTripLoadMarshalLoad3712892429
```
⇒ **同一形、同一枚用例、只差那一枚变量**：顶层 **0 新增**（差集件空＋`empty above＝0 枚`），残渣**逐字路径**在根里。根选在 `%LOCALAPPDATA%`（＝与 `%TEMP%` **同卷同 ACL 族**，⛔ 换卷，正典 §1 那一节；我也⛔ 为了"仓外好看"把根挪卷）。

## 2. 对腿交四发的复核（作废那两发指名为什么作废得对/不对）

| 发 | 腿写什么 | 我这把尺读到什么 | 判 |
|---|---|---|---|
| a1〔作废〕 | 冷台面首试：`run_lines=0`＋帽落在第一次 `t.TempDir()` 之前 | — | **作废得对**（坏世界没被造出来就当 0 用＝假绿那一形；它⛔ 拿那 0 当任何一边读数、件里留了 `shot_rc=1`） |
| a2 | 55→56，`TestRoundTripLoadMarshalLoad2461500944` | 盘上现量：`created=2026-10-12 11:43:24.391`，**今天仍在 58/59 枚名册里**（`logs/20-roster-a-before.txt.meta` 第 7 行） | **✓ 名、数、盘上三头对** |
| a3 | 56→57，`TestRoundTripLoadMarshalLoad2236627062`（"由它自报的 `test_star_dirs_now=57` 与外部门册作差对齐"） | ★**那枚名字在盘上任何名册里零命中**（下 §3） | **枚数/形 ✓、名字 ✗** |
| a4 | 57→58，`TestRoundTripLoadMarshalLoad826293033`，差集件`10c-topnew.txt` | 盘上现量 `created=2026-10-12 11:50:03.330` ✓；`10c-topnew.txt` 39 字节＝**恰好一枚名** ✓ | **✓ 这一发用的是差集尺，做得对** |
| b1 | 顶层 56→56、残渣 `…2086954487` 在 `D:/tmp/wisp-test-gotmp/probe-114728-128` | 我 `ls -1` 那枚根：逐字**只有** `TestRoundTripLoadMarshalLoad2086954487`（件＝`logs/20-root-contents.txt`） | **✓ 逐字路径复现**（⚠ 那一枚根在 `D:` 卷＝第一稿的换卷形，本腿⛔ 采用其根形） |
| b2〔作废〕 | 首试 58→58 但 `root_entries=0` ⇒ 帽落在两发迭代之间 | 同批的 `logs/10b2-root.txt`＝**0 字节**、`logs/10b2-new.txt`＝**0 字节** | **作废得对**（⛔ 拿"顶层没多"当"那一步挡住了残渣"，因为它同时读到根里也⛔ 有东西＝那形压根没落在用例里）；⚠ 这两枚 0 字节件**没被任何交付件按名点名**（`grep -c -F`＝0，见 `30-ac3-gates.md` §3 的 ⓐ/ⓑ 分类） |
| b3 | 58→58、残渣 `…1761995879` 在出厂根 `portable-20261012-115401-2013` | 我 `ls -1` 那枚根：逐字**只有** `TestRoundTripLoadMarshalLoad1761995879` ✓；`logs/10b2-r1-new.txt`＝0 字节且行内写了 `empty above＝0 枚` ✓ | **✓ 三头对** |
| b4 | 顶层 58→58、残渣 `…489708636` 在 `D:/tmp/.../shapec-115000` | 我 `ls -1` 那枚根：逐字**只有** `TestRoundTripLoadMarshalLoad489708636` ✓ | **✓**（⚠ 同一枚根目录里还躺着腿的**读数件** `shapec-115000-c2_redirected.run.txt`（0 字节）＋`….err`（35458 字节）＝**把测试输出写进了专用根**；我起手用 `Get-ChildItem -Directory` 数出 6 枚子目录、`ls -1` 数出 8 条目，这个 6⇄8 的差就是这两枚件。⛔ 判它违规（票面⛔ 禁写进根、只禁删），但**读数件⛔ 该住在被测程序的临时根里**——这条写进我欠账 `Q-310-v1-7`） |

⇒ **正控两发齐不齐**（派单口径）：**齐**——(a) 侧两形各一枚具名＋同行"那一形真发生了"的行，(b) 侧三发都同时给出"顶层 0 枚＋根里逐字路径"；两发⛔ 成立的（a1、b2）**自标〔作废〕并留了非零 rc**，这正是派单要的形状。

## 3. ★专攻"名字"那一形（本腿自己两把尺，⛔ 用编排者的件）

尺 1＝`grep -rl '2236627062' .scratch/wisp/probes/` 逐字命中 8 枚：`310/r1/20-paired-readings.md`、`310/r1/90-final.md`、`310/r1/logs/10c-watch.txt`（**突变台件自报的 `dir=` 字段**）、`310/v1/00-anchor.md`（我转述它）、`orch/20261012/txt/**`（派单与台账转述）——**⛔ 一枚名册件命中**。名册件＝`10c-roster-start.txt`（57 行）／`10c-roster-end.txt`（58 行）／`10b2-{before,after}.txt`／`10b2-r1-{before,after}.txt`／`11-roster-winsec-before.txt`／`12-roster-{before,after}.txt`，**逐枚含 `…1721168370`、逐枚⛔ 含 `…2236627062`**（尺＝`grep -rl` 两次，件＝本段）。
尺 2＝我自家 `CreationTime` 尺（`logs/20-roster-a-before.txt.meta` 第 7～9 行逐字）：
```
  roundtrip_name=TestRoundTripLoadMarshalLoad2461500944 created=2026-10-12 11:43:24.391
  roundtrip_name=TestRoundTripLoadMarshalLoad1721168370 created=2026-10-12 11:49:06.507
  roundtrip_name=TestRoundTripLoadMarshalLoad826293033 created=2026-10-12 11:50:03.330
```
腿件 `20-paired-readings.md` §1 a3 行的动作时刻＝逐字 `external_kill_issued at=11:49:06.508` ⇄ **同一毫秒活着的那枚是 `…1721168370`**（created `11:49:06.507`）。
⇒ **判**：a3 那一发的"漏"确实发生（名册 56→57 那一枚就是它，只是名字叫 `…1721168370`）；腿把**台件自报的 `dir=` 字段当成了漏的那枚目录名**——这正是正典 §1 点名的那一形（同一秒 `t.TempDir()` 会建了又删）。⚠ 加重一条：同一形状 `310-a1` 上午撞过、还为此换了差集尺（`20-shapes.md` §3.3），`310-r1` 复用了它的台件原形**却把它的旧自报名一起抄了**。

**本格要的那句裁决**："**具名来路错、但枚数与'死亡会漏一枚'的结论仍被两条独立来路支持**"这一形 ⇒ **`AC#2` 算成立**。理由：① 这一格裁的是"反恒真"＝"撤掉那一步当场多一枚、还原⛔ 多"，两半都被三条独立来路打中（腿 a2/a4/b3/b4 的差集名册＋我此刻的 (a)(b) 两发＋编排者 11:1x/12:3x 自己的两发）；② 错的是**一枚名字**，⛔ 是判据；③ 但**"具名凭据"这一维按正典就是⛔ 成立**——所以判语必须带这句尾巴：**凡"某某目录是我这一发漏的"那一格，腿件里 a3 那处⛔ 算交过**，更正归编排者落账（`A858` §3 已落），本腿⛔ 动票面原句。

**判"今后这类凭据必须自带 `CreationTime` 佐证"那条够不够硬**（正典 §1 我照的是这两行）：
```
⚠ 枚数⛔ 是证据，名字才是；而名字只有"作差"那种来路算数 …
⇒ 定式＝起进程前的名册 ⇄ 杀完之后的名册作差（`comm -13`）＋逐枚带 `CreationTime`；⛔ 引用任何"某某目录是我这一发漏的"之前，先用 `Get-ChildItem … | Select Name,CreationTime` 证它今天还在盘上。裁决腿⛔ 按交件里的名册盖章（本仓 10-12 两枚腿各栽过一次，见 `A857`／`A858`）。
```
**判＝那一行够硬，⛔ 需要再加半句**：它已经把"差集＋`CreationTime`＋今天仍在盘上"三件都写成必要条件，本腿照做即可自查（我这一发就是先做错成"名册"再被它拦住的形状——不，我这一发一次就做对了：差集件与两枚输入名册一起落盘）。⚠ 唯一缺口＝它管⛔ 住"**台件自报的 `dir=` 字段被写进交付件**"这一形（那字段⛔ 是名字来路、只该当动作行）；建议编排者把那半句补成具名禁令：`突变台件自报的 dir= 字段只能当"那一形发生了"的动作行，⛔ 当漏掉的那枚目录名`（`Q-310-v1-6`，落⛔ 归我，那是正典面）。
