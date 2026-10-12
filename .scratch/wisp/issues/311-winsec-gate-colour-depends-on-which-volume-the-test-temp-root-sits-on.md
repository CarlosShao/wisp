# 票 311 — `internal/winsec` 那档闸门的**颜色取决于测试临时根落在哪一枚卷**：同一份码，根在 `%LOCALAPPDATA%`＝`PASS=58 FAIL=0`，根在仓所在卷 `D:/tmp`＝`PASS=56 FAIL=2`，而且那两枚红用例的 `t.TempDir()` **自我收尾也失败**（`Access is denied`）

**立票**：2026-10-12 12:2x 编排者（来路＝票 310 落地腿 `310-r1` 撞出并具名报回的一条**派单里没写的代价**＋我自己现跑复现；件 `.scratch/wisp/probes/310/r1/10-impl.md` §2 第 2 条、`probes/orch/20261012/logs/orch-winsec-A-dtmp.txt`／`-B-lappdata.txt`／`orch-winsec-placement.txt`）
**性质**：★**仪器缺陷票（判据对夹具位置敏感）**。⛔ 产码缺陷：`internal/winsec` 的判定逻辑在两次读数里**一个字节都没变**，变的只有 `t.TempDir()` 的落点。也⛔ 票 310 的射程：310 管"残渣落在哪儿、谁收尾"，本票管"**同一道门会因为根换卷而改颜色**"——这是票 310 的修法**造出来**的新变量，所以它⛔ 并回 310、要单独收口。
**为什么要紧**：① 这一档是 CI 里**唯一**直接判本机 DACL 的门（票 110 立的 `winsec-tests.sh`，GUARD 1/2 就是为"它别是枚空仪器"），一门会因环境换色＝**将来任何一次"红"都要先问一句是不是根换了**，而那一句今天⛔ 有仪器能答；② 它会把后续的落地腿推向一枚危险动作——为了消红去**改判据**（本仓明令 ⛔ 为变绿放宽断言）；③ 那两枚红用例的清理失败＝票 310 现量那句"测试正常跑完 ⇒ 0 残留"**只在默认根上成立**，换卷后**正常收尾也漏**（同行证据见下面第 3 条）。

## 现量（每条都带尺；⛔ 引用前先重跑，枚数与颜色都是快照）

- ★**成对两发（我 12:21 自己跑，同一台面＝工作树 HEAD `a0d3b521`＋票 310 已入库的三枚脚本增量，只差那一枚根）**：
  - 根＝`D:/tmp/wisp-test-gotmp/winsec-20261012-122119-659`（靠 `WISP_TEST_GOTMP_ROOT` 覆盖）⇒ 逐字 `portable-tests.sh: four numbers … === RUN=101 --- PASS=56 --- FAIL=2 --- SKIP=0`、`winsec-tests.sh: gate is RED (rc=1)`；红名册逐枚＝`--- FAIL: TestGateJudgesThePrivateSetByResolvedSID (4.80s)`＋两枚子例 `…/runner_descriptor_is_judged_by_what_it_resolves_to (1.82s)`／`…/a_name_spelling_cannot_change_the_verdict (1.19s)`＋`--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (2.01s)`。
  - 根＝`C:/Users/swq/AppData/Local/wisp-test-gotmp/winsec-20261012-122147-814`（出厂默认）⇒ 逐字 `=== RUN=101 --- PASS=58 --- FAIL=0 --- SKIP=0`、`shot_rc=0`。
  - 腿 `310-r1` 的同两发（件 `probes/310/r1/logs/11-winsec.txt`／`12-winsec-shipped.txt`）与我**逐字同数**；⚠ 我另用错过一次拼写：`bash scripts/winsec-tests.sh --scope=winsec` ⇒ 被它自己的 GUARD 1 打回 **rc=2**（具名两行在 `logs/orch-winsec-dtmp-err.txt`），正确形＝`bash scripts/winsec-tests.sh ./internal/winsec/`。
- ★**那两枚用例确实读的是"临时根在哪儿"**（尺＝`sed -n '128,150p;270,300p' internal/winsec/private_set_sid_windows_test.go` 现读）：`:134` `dir := filepath.Join(t.TempDir(), "runner")`＋`plantDescriptor(...)`，`:277` `root := filepath.Join(t.TempDir(), "store")`＋`PrivateDirAll(root, 0o700)`；两枚都以 `privateSetSIDsBySID(t)` 为集合、再问"解析出来的 trustee 在不在集合里／被封之后清掉了谁"。⇒ **继承来的那套 DACL 就是它们的输入**，而 `%TEMP%`（`C:` 用户目录）与 `D:` 卷根给的是**两副不同的继承集合**〔这一句是读码推的，⛔ 定罪——差在哪一枚 SID 归 `AC#0`〕。
- ★**第三种死法：正常跑完也漏，因为清理自己失败**（尺＝`grep -c 'TempDir RemoveAll cleanup'`，同行落盘）：`D:` 根那发＝**4 行**，逐字形如 `testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/.../TestSealNarrowsAndNamesThePrincipalItRemovedBySID3143495113\002\la: Access is denied.`；`%LOCALAPPDATA%` 根那发＝**0 行**（件 `logs/orch-winsec-A-dtmp.txt:140,143,145,152`／`logs/orch-winsec-B-lappdata.txt`）。同批我还量了两枚专用根的条目数：`winsec-20261012-122119-659 entries=4` ⇄ `winsec-20261012-122147-814 entries=0`（件 `logs/orch-temp-after-my-shots.txt`）。⇒ ⚠ 票 310 现量那句"测试正常跑完 ⇒ 0 残留"**⛔ 全域成立**，它只在默认根上成立。
- ⚠ **今天的编队起手名册没被这一族咬到**（尺＝`probes/orch/20261012/gate-310r1-recheck.ps1` 现跑）：`%TEMP%` 顶层 `Test*`＝**58**（我两发 winsec 前后同数，因为两发都走重定向根），`svchost.exe`＝86／`*.exe`＝353 正控同行、家族尺 `*.test.exe`＝**0**、四条精确名各 0、那两发的 stderr **0 字节**。⇒ 换根**没有**把残渣推回 `%TEMP%` 顶层，这一条对票 310 `AC#1` 是**加分**、对本票是"残余在专用根里、没人收尾"那一半仍未闭。
- ★**在册 58 枚的构成我 12:20 现量过（这一条同时更正我自己昨晚的一句话）**：尺＝`Get-ChildItem $env:TEMP -Directory | ? Name -like 'TestRoundTrip*' | Select Name,CreationTime`，件 `logs/orch-roundtrip-times-err.txt`（0 字节）＋对话同批读数。8 枚 `TestRoundTripLoadMarshalLoad*` 的 `CreationTime` **全部是今天**：`10:45:38`／`10:46:42`／`10:47:25`（＝普查腿 `310-a1` 那三发）＋`11:04:41`／`11:05:29`（＝**编排者自己**上午那两发 (b) 形）＋`11:43:24`／`11:49:06`／`11:50:03`（＝`310-r1` 的 a2／a3／a4）。⇒ 三家的数**加得起来**（55→58 那条成立），但 `310-r1` 件里把 a3 那枚的名字写成 `TestRoundTripLoadMarshalLoad2236627062`——**那枚名字在盘上从没存在过**（`grep -rl 2236627062 probes/` 只命中腿自己的三枚件），它那一秒真正活着的是 `…1721168370`。来路＝腿抄了突变台件自报的 `dir=`，而那个字段盯的是"那一瞬间看到的任一枚"（同一形态 `310-a1` 自己也撞过并在 `20-shapes.md` §3.3 具名认过）。⚠ **枚数与"死亡会漏一枚"这条结论⛔ 因此动摇**，只有那一枚**名字**错。

## 要建什么（`AC#0` 之前⛔ 动产码、⛔ 动判据、⛔ 改根）

- [ ] **`AC#0` 归因到"差在哪一枚 SID"**（只读）：同一枚夹具在**两枚根**上各跑一次，把封之前／封之后的 DACL 逐 ACE 打出来（`Get-Acl`／仓里现成的 `resolvedTrustees` 口径），具名回答三问——① 两枚根的**继承集合**差哪几枚 trustee（逐枚 SID＋掩码）；② `TestGateJudgesThePrivateSetByResolvedSID` 的红句里那枚"解析出来却⛔ 在集合内"的 SID 是什么；③ 那 4 行 `Access is denied` 发生在哪一步（封之后没人有权删＝**自我锁死**，还是根本身的 `D:` 卷 ACL 就拒 `WRITE_DAC`）。判据＝三问各带逐字读数行＋**成因那半句必须标〔读码推的〕还是〔真机量的〕**；⛔ 由 `AC#0` 改任何字节。
- [ ] **`AC#1` 三形择一处置（⛔ 由 `AC#0` 的读数裁，⛔ 由写腿挑）**：
  - **甲＝让判据⛔ 依赖继承**——夹具自己把基线 DACL 树成**已知的一副**（显式 `SetSecurityDescriptor` 再封），换卷也同一答案。落点在 `internal/winsec/*_test.go` 夹具那一层；⚠ 它会碰到 C21 冻结 token 面吗＝`AC#0` 必答一句（尺＝`grep -n 'tokens_table_test\|SilenceLevelGate' internal/winsec`）。
  - **乙＝承认"根的位置是判据的一部分"并把它钉起来**——门里加一条**响在退码上**的守卫：`WISP_TEST_GOTMP_ROOT` 若指向与 `%TEMP%` 不同卷 ⇒ 具名报错并 rc≠0（今天⛔ 报，今天是"换了卷、门静悄悄变红"）。落点＝票 310 刚动过的那三枚 `scripts/*.sh`。
  - **丙＝不做**——登记为具名已知仪器边界＋把"⛔ 换卷"那一行留在编队正典（**我今天已按丙先行落政**，见 `probes/orch/gate-roster-304.md` §1 新增那一节），代价＝将来任何人把根挪卷都会看到两枚莫名其妙红，而"为什么"只能靠人记。
  判据＝三形各带"`AC#0` 里支持它的那一行逐字"＋"落点文件:行"＋"最坏没防住的形状"。⚠ 本票⛔ 把"甲"当默认；裁形前⛔ 动 `internal/winsec`。
- [ ] **`AC#2` 牙**：判据⛔ 只是"换卷也绿"。ⓐ**成对两发**＝默认根 ⇄ 另一卷根，落地后两发**同数**（甲形）或第二发**具名报错 rc≠0**（乙形）；ⓑ**反形正控**＝把甲形那一步（或乙形那条守卫）摘掉，同两发必须回到 `58/0` ⇄ `56/2` 那对差异——⛔ 摘掉仍全绿＝那一枚步骤对这件事不敏感，⛔ 算凭据（本仓 `A` 系列反复栽过）；ⓒ 两发都同行落"坏世界真被造出来"的证据行（那 4 行 `Access is denied` 或 `panic: test timed out`），⛔ 拿"0 枚"当读数。
- [ ] **`AC#3` 第三种死法具名归口**：`TempDir RemoveAll cleanup: … Access is denied` ⇒ **正常收尾也漏目录**。本格只回答一句：**这算不算票 310 现量里"死亡三形"之外的第四形**，算的话落点写在哪（票 310 追加节 ⇄ 本票），并把它的**触发条件**写成一句能判的话（"根所在卷的继承集合含 X ⇒ 封完之后自己删不掉"）。⛔ 顺手在本票实现清理器。
- [ ] **`AC#4` 门禁＋越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` rc=0；`sh scripts/check-path-length-budget.sh --with-self-test` rc=0；`bash scripts/portable-tests.sh --scope=census` totals 行逐字未变（基线现量＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，跑前先取一次、跑后再取一次作差，⛔ 拿旧数当基线）；`bash scripts/winsec-tests.sh ./internal/winsec/` 改前改后各两发（⛔ `--scope=winsec` 那串拼写，见现量第 1 条末）；格式名册两把并排（工作树 ⇄ `git show HEAD:` blob，写清射程目录＋枚数）；每把门禁件自落一行 `rc=N`（⛔ 0 字节件＝那格没交）；`git show --name-only --format=` 逐笔量越界（⛔ 区间尺）。
- [ ] **`AC#5` 存量与归口（编排者侧的活，⛔ 落地腿）**：`Q-304-a1-H` 那笔逐枚归属今天按名册已能答一半（上面第 4 条那把 `CreationTime` 尺＝12 秒内出全名册）⇒ 由我在台账把该笔**销掉并注明"由票 310 的 `CreationTime` 尺销"**；⛔ 任何腿删任何目录（含专用根里的残渣，含那枚"归属不明"的 `TestProbeStepByStep3843095677`＝票 310 已裁"⛔ 动、永久"）。

## 禁区

- ⛔ **为变绿放宽任何断言**、⛔ 改 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47、⛔ 动三枚冻结件（`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`）。
- ⛔ 把根**换卷**当"顺手优化"；⛔ 在未裁 `AC#1` 之前改 `internal/winsec/**`；⛔ 动票 304 的 Job Object 面（`internal/proc/jobscope_*`、`internal/llm/adaptertest/childjob_*`、两枚 mockllm spawn 点）；⛔ 动 `ci.yml`／`.github/**`（托管 runner 上 `%TEMP%` 兄弟根可不可写＝另一枚欠账 `Q-310-r1-3/-5`，随 CI 色那一批一起取）。
- ⛔ 删目录（票 310 那条铁规矩在本票一字不动地适用）；⛔ 碰别家项目的 `Test*` 目录；⛔ 碰 `msedgewebview2.exe` 那几棵树。
- ⛔ 新造协程（D38(b) 名册零膨胀）、⛔ 裸 `go func(`（ban #1）。
- git：只 commit ⛔ push；commit 必带**显式 pathspec**（`add` 与 `commit` 两处都带）；⛔ `add -A`／`add <目录>`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；仓内⛔ 建 worktree；临时件**只建不删**（只建在本腿自家 `probes/311/**`）；证据件⛔ 叫 `.out`、⛔ 在 `probes/**` 下建 `.go`、件名⛔ 带 `§`⛔ 长名（相对 ≤121／末段 ≤100／全路径 ≤165，新建前先跑长度门）。

## 排程

`311-a1`（＝`AC#0` 只读归因＋`AC#3` 那一问，⛔ 产码、⛔ 翻框；要跑 `go test` ⇒ 起手先过编队名册那三问）→ 编排者按 `AC#1` 三形裁（若裁丙＝本票就地结案成"具名已知边界"，⛔ 落地腿）→ 落点腿。
⚠ 与票 310 **同源不同题**：310 的 `AC#1` 已入库（三枚 `scripts/*.sh`），本票⛔ 重开那格的判语；本票要的"专用根换卷会改颜色"是 310 那一步**造出来的新变量**，所以两张票各自收口、⛔ 混批。⚠ `internal/winsec` 写面与票 119/137/154 那族同源，派腿前先看编队里有没有别的腿碰 `internal/winsec`。

**Status:** **未开工**（本票 12:2x 立，五格全未勾；现象＝**已证**〔我 12:21 成对两发复跑同数＋那 4 行清理失败逐字在件里〕，成因＝**未定**〔差在哪一枚 SID，归 `AC#0`〕）。⛔ 零翻框、⛔ 零 push、⛔ 零删目录。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
