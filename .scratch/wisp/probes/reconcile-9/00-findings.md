# reconcile-9｜只读对账：`docs/reports/pending-and-issues.md` 的 `A730`–`A750` 承重数复跑

腿＝`reconcile-9`（只读）。写面＝本件一枚。⛔ 零产码改动、零票面改动、零台账改动、零 push、零 `go`／`scripts`／`gh`／网络取数。
刀＝`reconcile-8` 那一套，往前挪一节：`A730`（`:14327`）到 `A750`（`:14502`–`:14507`）。

## 0. 尺具名（每一张表都用这些；写法不同会差 1，所以逐条给原文）

- 取数时刻＝2026-10-09 午后（本腿 42 次调用的窗口内）。起手 HEAD＝`ce18b3d6`（标题 `A760 收 285-r1＋reconcile-8…`，`Fri Oct 9 11:01:21 2026 +0800`）。
- 历史面锚定：枚数／行号／名册一律钉到台账自己点名的那枚 ref（`6410a062`／`eb2a217a`／`e7cd0c6e`／`e060cef6`／`50f220b6`／`d844e379`／`0e0ac7b9`／`ecba5bf2^`／`1ccd57b6`／`9ae812d6`／`ccf16baa`／`77150dcc`／`0abe217c`／`369a8f29`／`0ef0de11`／`c39e2853`），⛔ 不拿漂移中的 HEAD 复跑旧句。
- 现状面锚定：产码行号一律 `git show HEAD:<path>`（或 `git grep -n … HEAD -- <pathspec>`）。
- 框尺（同 `reconcile-8`）：`grep -cE '^[[:space:]]*- \[ \]'`＝未勾、`grep -cE '^[[:space:]]*- \[x\]'`＝已勾；件名走 `git ls-tree -r --name-only <ref> -- .scratch/wisp/issues/ | grep -E "issues/<NN>-"` 全路径取（改名前后的旧名同一把尺）。
- **未勾框行号尺**（本腿新增用途，钉 `A740` 的 L 名册）：`git show <ref>:<件> | grep -nE '^[[:space:]]*- \[ \]' | cut -d: -f1`。
- Status 尺：`grep -nE '\*\*Status'`。⛔ **不能写 `\*\*Status\*\*`**——本仓库形是 `**Status:**`，用错会得假 0（本腿第一发就踩，重跑才回真值）。
- 调用形状尺（负向句专用）：`git grep -nE "[A-Za-z0-9_)[:space:]]\.符号\(" HEAD -- '*.go' ':!*_test.go'`；名册尺 `git grep -l <词> <ref> -- .scratch/wisp/issues/`。
- 词面尺（`A736` 的 `go func(`）：`git grep -c "go func(" HEAD -- 'cmd/**' 'internal/**' ':(exclude)*_test.go'`＝**恰 1 枚文件**；不带这组限定器则全仓 **81 枚文件**含 `go func(` ⇒ 台账那句"词面总 1"的射程必须靠限定器撑着，原文没写。
- hash 尺：`git cat-file -t <hex>` 分 commit／blob，`git rev-parse HEAD:<path> | cut -c1-8` 对 blob 等值，`git cat-file -s <blob>` 量字节。
- 台账零改动自证：本腿只读 `sed -n`／`Read`，未对 `docs/**`、`.scratch/wisp/issues/**`、任何 `*.go` 下笔。

判定三形：`对上`／`对不上（真值 X，尺写法 Y）`／`复跑不了（缺什么）`。另有「**对上但射程／措辞不足**」与「**本腿预算内未开尺**」两节，⛔ 不当"没问题"糊掉。

## 1. 汇总表

| 节 | 复跑读数 | 对上 | 对不上 | 复跑不了 |
|---|---|---|---|---|
| A730 | 7 | 7 | 0 | 0 |
| A731 | 8 | 7 | 0 | 1 |
| A732 | 6 | 6 | 0 | 0 |
| A733 | 7 | 6 | 1 | 0 |
| A734 | 6 | 5 | 0 | 1 |
| A735 | 3 | 3 | 0 | 0 |
| A736 | 5 | 3 | 0 | 2 |
| A737 | 4 | 4 | 0 | 0 |
| A738 | 5 | 5 | 0 | 0 |
| A739 | 8 | 8 | 0 | 0 |
| A740 | 6 | 6 | 0 | 0 |
| A741 | 4 | 4 | 0 | 0 |
| A742 | 5 | 3 | 2 | 0 |
| A743 | 9 | 9 | 0 | 0 |
| A744 | 5 | 4 | 0 | 1 |
| A745 | 7 | 6 | 1 | 0 |
| A746 | 6 | 4 | 2 | 0 |
| A747 | 8 | 8 | 0 | 0 |
| A748 | 6 | 6 | 0 | 0 |
| A749 | 5 | 5 | 0 | 0 |
| A750 | 9 | 5 | 2 | 2 |
| **合计** | **129** | **114** | **8** | **7** |

## 2. 对不上的 8 笔（按杀伤力排序）

### N1（最要命）`A746` 的 `AskOnTaskRoot`／`Replay` 查点名册少列 1 枚，且这串错名册已随 `ecba5bf2` 的标题进 git 历史

- 台账原文（`:14470`）：``AskOnTaskRoot`/`Replay` 命中 `11/170/228/246/256/258/259-done/87` ⇒ 246 是母票、本票承接它的**两条未接路径**``。
- 本腿尺：`git grep -lE 'AskOnTaskRoot|Replay' ecba5bf2^ -- .scratch/wisp/issues/`（`ecba5bf2` 就是"补票四张"那一笔，10-08 19:44；`^` 即查重当时的盘面）。
- 真值：**9 枚**＝`11 170 228 246 256 258 259 87 97`。少列的那枚是 **`97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md`**（命中处 `:8`，逐字含 `` `Replay` 非答案``）。分词复核：`AskOnTaskRoot` 3 枚（`246/256/258`）、`Replay` 7 枚（`11/170/228/246/259/87/97`），并集 9 枚。
- 同一串数还写进了 commit 标题（`git log` 里 `ecba5bf2` 那句"AskOnTaskRoot/Replay 命中 11/170/228/246/256/258/259d/87"）⇒ 历史面不可回改，只能在台账追加。
- 影响：这条是票 278 的**查重凭据**。后续程拿它当"Replay 这个词只有这 8 张票碰过"的名单用，就会漏掉 97——而 97 的主题恰好是"注释 invents a caller"，与 278 要裁的"建了但没接"同族。裁决定向（票要立）不受影响，**名单本身错**，必须改成 9 枚。

### N2 `A746` 的"举卡／行为凭据 命中 `244/245/246/256`"多计 1 枚：`245` 两个词都零命中

- 台账原文（`:14470`）：`"举卡/行为凭据"命中 `244/245/246/256` ⇒ 取新票并在票里写明与 246 的边界`。
- 本腿尺：`git grep -lE '举卡|行为凭据' ecba5bf2^ -- .scratch/wisp/issues/` ＋逐枚 `git grep -ohE '举卡|行为凭据' … -- .scratch/wisp/issues/<NN>-*`。
- 真值：**3 枚**＝`244`（举卡）／`246`（举卡＋行为凭据）／`256`（举卡）。`245-the-bare-esc-cancel-default-is-a-global-hotkey-…` 里 `举卡` 0 命中、`行为凭据` 0 命中（只有泛用的"凭据"二字，`git grep -cE '举卡|行为凭据|凭据'` 才回它）。
- 影响：与 N1 同一句、同一凭据面。若照 4 枚名单去核对票 277 的边界，会多写一条"与 245 的边界"——那一枚根本不存在。

### N3 `A745` 的"决定行为读点 29 处"分文件名册三处虚增，名册合计 34 与自称的 29 不自洽

- 台账原文（`:14463`）：`决定行为读点 **29 处**（`tools/bridge.go` 5／`approval/gate.go` **10**／`queue.go` 3／`replies.go` 5／`subagent_197.go` 1／`task.go` 1／`subagent_roster_197.go` **2**／`resident_approval_windows.go` **4**／`run.go` 1／`models.go` 1／`report.go` 1）`。
- 本腿尺（钉在腿件自己的 D 表上）：`git show HEAD:.scratch/wisp/probes/242/corrcensus1/20-types-and-readers.md | awk '/^## 2\./,/^## 3\./' | grep -oE '^\| D[0-9]+ \| [a-zA-Z0-9_/.]+' | awk '{print $4}' | sed 's/:[0-9].*//' | sort | uniq -c`。
- 真值：总数 **29 对**（D 表恰 29 行 ✓）。分文件＝`bridge.go` 5／`approval/gate.go` **7**／`queue.go` 3／`replies.go` 5／`tools/subagent_197.go` 1／`tools/task.go` 1／`panel/subagent_roster_197.go` **1**／`resident_approval_windows.go` **3**／`run.go` 1／`memory/models.go` 1／`report.go` 1 ＝ 29。台账那三枚（10／2／4）加起来 34，且**不是**"读点＋审计点"的合并形（审计表 A 段的 `gate.go` 只有 2 枚、`resident` 有 6 枚，拼不出 10／4）。
- 影响：这条是 `A745` 那一记"产码默认应改"裁定的**依据面**。后续程若照 `gate.go 10 处` 去排落地面／估改动量，会多做 3 处；更糟的是它会以为"分文件能加回 29"从而怀疑腿件。裁定本身不必翻（29 那枚总数对），**分文件名册要就地打旧**。

### N4 `A742` 的尺②原文复跑回 20 而不是 0（结论对、尺缺一个过滤器）

- 台账原文（`:14439`）：`②累计 `+/-` 行过滤后**非注释行＝0**（`grep -vcE '^[-+]\s*//'`）✓＝全注释改动`。
- 本腿按原文逐字复跑：`git diff -U0 a1b19873^..e5b46742 -- cmd internal | grep -E '^[-+]' | grep -vcE '^[-+]\s*//'` ＝ **20**（不带 `-U0` 也是 20）。
- 那 20 行逐枚读出来**全是** `--- a/<file>` 与 `+++ b/<file>`（10 枚文件 × 2）⇒ 真结论（**没有一行非注释改动**）站得住，但台账写的这把尺**少了排除 diff 文件头的那一步**（正确写法：`| grep -vE '^(---|\+\+\+) '` 再计数＝0）。
- 影响：中。后续程照抄这把尺会得到 20，最可能的误读是"land-1 里夹了 20 处非注释改动"，反过来怀疑那 13 枚注释落地的干净度——而它其实是干净的。

### N5 `A742`／`A733` 的 `config_readers_255.go:109` 今天已是 `:110`（正是 A742 自己在同一节里描述过的那族"顶位"）

- 台账原文（`:14441`，`A742`，写于 land-1 落地**之后**）：`照贴 +3 行会把 `cmd/wisp/run.go:991`（`[cfg := rt.cfg]`）漂到 `:994`，而它被 `config_readers_255.go:109` 名册逐字引`；同一句在 `A733`（`:14366`，写于 land-1 **之前**）也出现。
- 本腿尺：`git show <ref>:cmd/wisp/config_readers_255.go | grep -n 'run\.go:991'`。
- 真值：`:109` 只在 **`e7cd0c6e`（10-08 18:52，land-1 之前）** 成立；**`e5b46742`（19:21，land-1 末笔）起就是 `:110`**（HEAD 亦然，三处命中 `:58/:107/:110`）。原因＝land-1 那批自己也改了 `config_readers_255.go`（它在 `a1b19873^..e5b46742` 的改动文件单里）＋1 行。`run.go:991` 那枚锚本身**今天仍逐字对**（`cfg := rt.cfg` ✓，`:1008` = `AdmitTask: rt.admitTask,` ✓）。
- 影响：票 279（P39 与 255 名册联动改）的派单如果照 `:109` 取那行，会拿到"值已换进内存"那句而非名册条目——而这一族正是 `A742` 自己在同一节里给 `doc.go:27` 修过的形状（那一处我复跑**对**：`doc.go:27` 今天逐字写 `:370`，且 `resident_approval_windows.go:370` 确是 `ra.gate = approval.New(approval.Options{` ✓）。**引行号前先重跑**这一条对 `:109` 又灵了一次。

### N6 `A750` 的"只此一枚文件"是 pathspec 造出来的射程，同区间全树是 5 枚

- 台账原文（`:14505`）：``git diff --numstat c39e2853 HEAD -- docs/reports/HANDOVER.md`＝**`26 1`**（只此一枚文件）✓`。
- 本腿尺：带 pathspec 复跑 `git diff --numstat c39e2853 0ef0de11 -- docs/reports/HANDOVER.md`＝`26 1` ✓（`wc -l` 亦 **1913 → 1938**＝+25 ✓ 与 26−1 自洽）；**剥掉 pathspec** 复跑 `git diff --name-only c39e2853 0ef0de11`＝**5 枚**（`probes/281/r1/00-anchor.md`、`probes/283/r1/00-anchor.md`、`probes/parking-3/{00-anchor,01-section}.md`、`HANDOVER.md`）；只看 `parking-3` 那三笔自身（`4e9a2df6`／`ef56d0d8`／`0ef0de11`）＝**3 枚**（HANDOVER ＋ 它自己的两枚探针件）。
- 影响：中。"只此一枚文件"被读成"那一波只动了停车点"就会漏掉同波入库的探针件（`reconcile-8` 的 F5 是同一族：数的射程被换、方向没错）。数 `26 1` 没错。

### N7 `A750` 的新节标题写成"逐字"但少了两个加粗标记与整句尾巴

- 台账原文（`:14505`）：`新节＝`## 4.0af 停车点（2026-10-09 08:4x+0800 版；新会话从这一节起读）``。
- 本腿尺：`git show 0ef0de11:docs/reports/HANDOVER.md | grep -nE '^## 4\.0af'`（同 `reconcile-8` §5 第 20 条那族的"剥 `**`"形）。
- 真值（`:1915`）：`## 4.0af 停车点（**2026-10-09 08:4x+0800 版；新会话从这一节起读**，`4.0ae`（10-08 18:3x）及更早只作追溯）`——剥了两个 `**`、并整句截掉了"及更早只作追溯"那一支。同节的另外几枚**全对**：`grep -cE '^## 4\.0a[f]'`＝**1** ✓、H2 `^## ` 计数 **39→40** ✓、节头名册 `aa:87／ab:167／ad:261／ac:287／ae:1884／af:1915` **六枚逐枚对上且稳定到 HEAD** ✓（`ad` 排在 `ac` 之前**不是笔误**，文件里真就是这个顺序）、`4.0ae` 只改标题行、正文一字未动 ✓（全区间只 1 行删除，正是那一行标题）。
- 影响：轻（措辞级），但台账用的是"新节＝"这种逐字口吻 ⇒ 与 `reconcile-8` F7 同形登记。

### N8 `A733` 的"13 处判语：11 处对上"与腿件自己的分布不符（真值 10，且台账合计成 14）

- 台账原文（`:14364`）：`13 处判语：**11 处对上／P10 对不上／P02 一处差／P39 有条件**（P09、P12 两处轻差）`。
- 本腿尺：`git show e7cd0c6e:.scratch/wisp/probes/comment-fix-check-1/01-check.md` 的判语列 `grep -oE '\*\*(基本对上|有条件对上|对不上|对上)\*\*' | sort | uniq -c` ＋该件 `:133-134` 的总判原文。
- 真值：**对上 8／基本对上 2／对不上 1（P10）／有条件对上 1（P39）** ＝ 12 行，第 13 行是 P02（该件总判称"带红"）。腿件自己的总判逐字＝`13 处里 **10 处clean对上**（P01/P03/P04/P06/P07/P08/P11/P13＋P09/P12"基本对上"各带一处具名轻差）` ⇒ "对上"那一档**最多 10**（8 纯＋2 基本）。台账写 11，又把 P10／P02／P39 单列 ⇒ 合计 14 ≠ 13。件长 **135 行** ✓ 对。
- 影响：轻——三处真问题（P39／P02／P10）的处置台账都逐条写对了。但枚数被后续程引用时会自相矛盾（11＋3＝14）。

## 3. 复跑不了的 7 处（具名缺哪把尺）

| 条目 | 缺的读数 | 本腿为什么拿不到 |
|---|---|---|
| `A731` 四发突变的红句集与 `M4 rc=0 全绿` | `go test . -run '…'` 的 rc／PASS 枚数 | 题目禁 `go test`。本腿只核了静态面：`workspace.go:120`（`res.Actable()` 闸）／`:122` 被它挡／`:91`＋`:141` 两处 `Rewritten: res.Rewritten` 填充点／`:144-146` 死枝 **逐枚在位**；`RequestWorkspaceSwitch` **非测试零调用者** ✓（调用形状尺只回 `cmd/wisp/instructions_200r2_test.go:124` 与 `panel_pump_test.go:357` 两枚 `_test`）；`16901acb` 对 `internal/risk/` 空 diff ✓；evidence 件 **12,089 字节** ✓；票面 94→**101**／2→**1**／5→**6** ✓（`b4602f13` → HEAD 实测）；"只追加 §7" ✓（`:96` 真是 §7 节头）。⇒ 只有"那四发真跑过"无法独立复现 |
| `A734` 三发突变红句＋复绿 12/12 | 同上（`go test`） | 静态面全对上了：`approval.go:558` 逐字 `func (s *grantStore) spend(nonce, bind string) grantDenial` ✓；"改前那句已撤"✓（`git grep -n "cannot be replayed" HEAD -- internal cmd`＝**零命中**）；两枚尺件在库 ✓；票面 59→**69**／5→**1**／1→**5** ✓（`0af6594d`＝59/5/1、`ccf16baa`＝69/1/5）；`subagent_selfapproval_197_test.go:109` 当时逐字 `TaskID: taskID, CorrelationID: taskID,` ✓、末碰 `a818df46`＝**09-30 12:05** ✓ |
| `A736` ① d22scan 分母（bans #1-5 `internal/=228 cmd/=38`、#6 `f/=85`、#7 `tools/=23`、#8 `design/=39 f/=85 internal/=516 cmd/=113`） | `sh scripts/d22scan.sh` 的输出 | 题目禁跑 `scripts/*.sh`、禁 `go`。台账自己已写"⚠ 引这类数必带 ref"，本腿只能确认它**附了 ref**（`@50b34971` ✓ 该笔存在）⇒ 分母原样采〔腿报〕 |
| `A736` ④⑤ 四路径 porcelain 空＋未推 **427→431** | 18:58／19:0x 两个时点的盘上态 | 时点读数、盘上无历史尺可回打。本腿只核了两枚来历 hash 的日期：`5e8748b3`＝10-03 21:07 ✓、`3a343bc7`＝10-08 09:37 ✓（与"6 枚既有＋1 枚今天已提交"同向）；③`go func(` 词面 **1 枚** ✓（唯一位点 `cmd/wisp/testdata/esclistener/main.go` ✓，**但尺必须带 `cmd/** internal/** :(exclude)*_test.go` 三件限定器**，否则全仓 81 枚文件） |
| `A744` 两发突变的 rc=1 红句＋基线 12 PASS | `go test` | 禁 `go test`。静态面：`approval.go:565` 确是 `if !equalSecret(v, nonce) {` ✓、`:564` 确是 `for v, stored := range s.values {` ✓（`A738` 说的 `564-567` 那一圈）、`queue.go:390` 确是 `if it.state != statePending {` ✓；红句行 `ticket259_denial_rulers_test.go:139/:208/:259` 逐枚是 `t.Fatalf`／`t.Fatal` 行 ✓；还原态两枚 hash＝**`67fb1468`/`66fec7ae` 恰为 HEAD 上 `approval.go`/`queue.go` 的 blob** ✓（且两文件末碰仍在 10-05，"种前＝HEAD"成立）；`internal/agent/approval → 仓根` 差三级 ✓（`../../` 确实少一层） |
| `A750` 的 porcelain 空 | 08:4x 时点读数 | 同上，时点无尺 |
| `A750` 的 未推 **484** | 08:4x 的未推计数 | 需要当时与 `origin/dev` 的差；本腿不许联网、且该时点的远端 ref 状态不可回放 ⇒ 采〔腿报〕 |

## 4. 对上但射程／措辞不足（后续程别当已验）

1. **`A730` 的 `composer_dispatch.go:50`**：在 `6410a062` 上逐字含 `it has NO production caller yet` ✓（台账当时"与现 HEAD 逐字对上"是真的）；**今天整句已被 land-1 改掉**（`git grep "NO production caller yet" HEAD -- internal/panel/composer_dispatch.go`＝0 命中）。同节 `cmd/wisp/firstrun.go:11` 的 `… no production path ever did was CALL that pair once …` ✓ 今日仍在 `:11`。⇒ 这两枚锚命运相反，引用时必须带 ref。
2. **`A732` 的 `AskOnTaskRoot:697`**：真身 **`cmd/wisp/resident_approval_windows.go:697`**（@`eb2a217a` 逐字 `func (ra *residentApproval) AskOnTaskRoot(...)` ✓），台账**没写文件**；今天该定义在 `:698`（与 `A753` 记的 +1 漂一致）。`main.go:66`＝`runResident()` ✓、`resident_windows.go:260`＝`src := startResidentTaskSource(rt, ra)` ✓、`bridge.go:352/:443/:458` ✓、`.ui.Prompt(` **恰 2 枚**（`gate.go:294`／`:536`，调用形状尺）✓ 全部逐字对上。
3. **`A735` 的 `resident_approval_windows.go:885`／`:880`**：@`e060cef6` 逐字对上 ✓（`:885`＝`fmt.Printf("wisp: 卡片挂起：…")`、`:880`＝`slog.Info("approval: 常驻进程显示一张确认卡片",`）；今日各 +1（`A753` 已具名）。`01-recipe.md`＝**83 行** ✓、`bridge.go:424`＝`switch sil.Level {` ✓。`build/wisp.exe` mtime 一条**不复跑**（`reconcile-8` §5 第 17 条已判）。
4. **`A737` 的"写面只有 `subagent_selfapproval_197_test.go`，12 增/5 删"**：12／5 ✓ 精确；但 `77150dcc` 那笔**还含** `.scratch/wisp/probes/259/r4/10-split.md`（＋39 行）⇒ "写面"若指产码面则对，若读成"这笔 commit 只动一枚文件"则错（`git show --name-only 77150dcc`＝2 枚）。
5. **`A745` 的 `loop_approval_test.go:213-214`**：本腿**未按 `0e0ac7b9` 现读那一行**（预算内未开尺）；已知该句在 `A748`/`A749` 之后重写成 `:204`／`:221` 两处（本腿对 `:221` 现读＝`t.Errorf("correlation_id = %q, want a non-empty per-call id routing back to task %q (C18)"` ✓、`:204`＝`t.Fatalf("tool_call rows = %d…")` ✓ 形状）。⇒ 台账那句"恰一处把等值写成断言"的**行号只在该 ref 上成立**。
6. **`A745` 的 `其它产码赋值点 22 组（新造 12＋誊抄 10）`**：照抄腿件表头 ✓（`30-impact-pins-producers.md` §4 两档各写"（12 组）"／"（10 组）"）；但**腿件自己的行编号只到 21**（誊抄段为 13–21 共 9 行）⇒ 台账与源件在"10 组"这一档上自数不一致，**这条要编排者自己裁**（要么按行号 21、要么认定某行计两枚）。
7. **`A748` 的 `task.go:708`**：真身 **`internal/tools/task.go:708`**＝`caller := TaskID(ctx)` ✓（台账未具名目录；`internal/agent/task.go` 不存在）。同节 `loop.go:603`／`:676`／`cancel.go:67`／`corr_percall_242_test.go:82`／`subagent_197.go:262` **逐枚对上** ✓；三笔（`133b1bfa`／`bd124b2a`／`0abe217c`）逐笔查 `docs/reports|issues/`＝**零命中** ✓"台账/票面零动"属实。
8. **`A749` 的五枚 hash**：`87ac0624`（`tools/cancel.go`）／`0d53277d`（`agent/loop.go`）／`8c9266a7`（`tools/subagent_197.go`）／`3ec8d498`（`tools/task.go`）／`d46124de`（`tools/bridge.go`）＝**逐枚等 HEAD blob** ✓；且它报告把 197 写成 `internal/agent/` 的更正 ✓ 成立。`run.go:1003`＝`Journal: nil,` ✓、`loop.go:369`＝`j := newTaskJournal(...)` ✓、`subagent_197_test.go:302`＝`t.Fatalf("派生失败：%s", res.Text)` ✓（红句形状）。"三处 blob CR=0"本腿未开尺（见 §5）。
9. **`A750` 的"四把尺"**：句里其实列了**五发**读数（`wc -l`／numstat／`grep -c`／porcelain／节头名册），名册那发台账自己标了"它逐个数过"＝腿报 ⇒ "四把"这个自数与清单差 1（不影响任何结论）。

## 5. 本腿预算内没开尺的读数（具名，⛔ 不算已验）

`A734` 的"15 笔提交对禁列零命中"；`A737` 的 `bridge.go:376`（`TaskID` 只吃 `req.TaskID`）；`A739` 表一里 57 枚"存在"的**逐枚**路径归属（本腿只跑了 55 枚 hex 的存在性：commit／blob 全在，唯 `76b694e` **MISSING**＝与台账"那枚 `A722` 当场自纠过"一致 ✓、`37703959747` MISSING＝run id ✓）；`A742` 的"每处贴后 roster 尺 rc=0"（需 `go test`）；`A743` 的 `219 前置序 201→220`；`A745` 的 `loop_approval_test.go:213-214`；`A747` 的 `02` 件原文首 80 字；`A748`/`A749` 的票 282／283 框数（本腿末笔实测 **281＝4 框 @`52501e22`** ✓、**282＝5 框 @`369a8f29`** ✓，283 未测）；`A749` 的"三处 blob CR=0"（未具名是哪三枚文件 ⇒ 尺开不了，缺文件名）。

## 6. 硬的那一块全对（可放心引的名册）

- **`A740` 整节一枚都不错**——本腿把两张名册**全池重跑**了：工单 **263** 枚（含 README 共 **264** 个 `.md`）✓、`-done` **95** ✓；**甲类 10 张**用"无 `-done` 且未勾＝0"现算＝**恰那 10 枚**（`125/196/211/213/214/215/216/219/258/261`），逐枚已勾数 `125`=4／`258`=4／`261`=5、其余 7 张＝0 ✓；**乙类 8 张的 14 枚未勾框行号逐枚全中**——`07:L51`／`92:L68,L73,L79`／`97:L55`／`104:L52,L57`／`105:L56`／`110:L39,L47`／`113:L54,L56`／`115:L58,L64`（尺见 §0，@`50f220b6`）；README 的 `L15`／`L30-31`／`L189-191`（票 62 那形、2026-09-20）／`L201` 四段**逐字对上** ✓。
- **`A743` 的收口那笔全对**：尺 `125` **425 行**／`258` **78**／`261` **95**（未勾全 0）✓、`-done` **95→98** ✓（@`d844e379`）、`196`＝**14 行**／`219`＝**100 行**零框 ✓、`125` 原 Status 确是 `open`（`:3`）✓、`258`/`261` 原无 Status、现各在 `:3` ✓、`grep -nE '125-|258-|261-'` README＝**零命中** ✓。★更正 `A740` 的写法是**诚实的**：`A740:14419/14422` 至今逐字写着"全勾"，而 `A743` 明说"⛔ 不回改、就地打旧于此"——**这一处不是 `A751` 那一形**（本腿专门按 `reconcile-8` 的教训去验了"声称已更正"的每一句）。
- **`A739` 那节的三张表全对**：`66 枚／存在 57／不存在 9`＋`52 枚＝带目录 28＋裸名 20＋片段 4`＋`28/28`＋`HEAD=0 但工作树在 3 枚` 逐枚与 `01-sweep.md` 同文 ✓；`bridge.go` 裸名**＝9 枚** ✓（台账"它纠我"那笔对）；`02-commit2-readings-correction.md`＝**3,641 字节** ✓ 且 @`1ccd57b6` **确未跟踪** ✓ → 现 HEAD **已入库** ✓（"本程我代提"兑现）；`181` 的 101/1/6 ✓、`246` 的 **105 行** ✓；票 259 的 **"69→79 是偏数、真值 77"**——§12 那句错数**至今逐字在场**（`:77`），而 §13（`:81`）已具名打旧 ✓ ⇒ 与 `A751` 那笔的区别是**这一处旧句被指认了**，不是漏了；`ci.yml:521`＝`- uses: actions/checkout@v4` ✓、`:469`＝`if: ${{ !cancelled() }}` ✓ 两枚在 `1ccd57b6` 与 HEAD 上**都没漂**。
- **`A747` 全对**：`07`＝**147 行／1 未勾／2 已勾／Status `done`** ✓，`:51` 那格原文确为"20 states rendered in a debug cycle page/window; human screenshot review" ✓（"真未做＋不许代签"有字面凭据）；`105` Status 逐字 `**rejected-needs-fix**（2026-09-21 20:5x …）` ✓；`110` Status＝`ready-for-review` ⇒ "相抵" ✓；`115` 在 `docs/evidence/s1/` 的 `115-` 名下证据表 **find=0** ✓；合计 **14 枚未勾框＝漏勾 10＋判不动 3＋真未做 1**（与 `A740` 那 14 枚同尺同数）✓；`probes/<NN>/` 8 票**全空**（逐枚 `git ls-tree`＝0）✓；票 281＝**4 框** ✓。
- **`A741` 全对**：票 242＝**63 行／未勾 3／已勾 0** ✓；三格原文确在 `:18`（AC#1 真·跨卡）／`:21`（AC#2 出向读面不得带 grant）／`:28`（AC#3 本格不设产码要求）✓；`loop.go:654` 当时逐字 `TaskID: taskID, CorrelationID: taskID, CallID: p.call.ID,` ✓（@`9ae812d6`）；`subagent_selfapproval_197_test.go:116` 已拆（`CorrelationID: corrID`）✓＋`:396`/`:463` 两枚 `-corr-1`/`-corr-2` ✓。
- **`A738` 全对**：`approval.go:564-567`／`queue.go:390`／`workspace.go:120`+`:144-146`／红句 `:139/:208/:259`／两枚 hash ✓。
- **`A746` 的其余部分**：查重当时最大票号＝**276** ✓、`P39`/`run.go:991` 在 `issues/` **零命中** ✓、新四张票 `277/278/279/280` **各 3 框、已勾 0** ✓、那三笔 leg commit **没碰台账与票面** ✓（N1／N2 只打在"命中名册"那两发上）。
- **`A750` 的 HANDOVER 那半边**：`1913+25=1938` ✓、`26 1` ✓、`4.0af` 唯一 ✓、六枚节头行号逐枚 ✓、`H2 39→40` ✓、`4.0ae` 只改标题行 ✓、三笔 `4e9a2df6`/`ef56d0d8`/`0ef0de11` 存在 ✓。

## 7. 一句话结论

`A730`–`A750` 这 21 节里，本腿复跑 **129 处**可独立复现的读数：**对上 114**、**对不上 8**、**复跑不了 7**（`go test` 红句／`d22scan` 分母／时点 porcelain 与未推计数，均已给静态等价复算并逐条具名缺口）。八笔里最要命的是 **N1／N2 同一句**（`A746:14470` 的两张查点名册：`Replay` 族少列 `97`、"举卡"族多计 `245`，且同串数已写进 `ecba5bf2` 的 commit 标题）；其次是 **N3**（`A745:14463` 的 29 处分文件名册三处虚增，合计 34 与自身 29 不自洽）；再次是 **N4／N5**（`A742:14439` 那把"非注释行＝0"的尺原文复跑回 **20**——全是 diff 文件头；同节 `config_readers_255.go:109` 在 land-1 之后已是 **`:110`**，正是该节自己描述的"顶位"那一族）。**`A739`／`A740`／`A741`／`A743`／`A747`／`A748` 六节一枚不错**，其中 `A740` 的两张名册（含 14 枚未勾框的行号）与 `A743` 的收口尺是本腿复跑里最硬的一块；`A743` 对 `A740` 的"★更正"经查是**诚实的就地打旧**（`A740` 的旧句仍在、且被指名打旧），与 `A751` 那一形不同，可以照它做。判决类只登记一句：`A745` 誊抄段"10 组 vs 行号只到 21"这一档**要编排者自己裁**，本腿不动。
