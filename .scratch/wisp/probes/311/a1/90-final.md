# 311-a1 `90-final` — `AC#3` 那一问的答复 ＋ 残留名册 ＋ 越界与自报

本件只收三样：`AC#3` 的答复、本腿留下的每一枚目录、以及"本腿到底动了哪些字节＋自报的越格动作"。
三问的答复正文＝`10-acl-before-after.md` §1–§2（问①）／`20-red-sentence.md` §3（问②）／`30-cleanup-failure.md` §3–§4（问③）。

## 1. `AC#3` 那一问的答复（一句话）

**⛔ 是票 310"死亡三形"之外的第四形——它是第 (a) 形（正常收尾）的一枚条件化破口**：那一形里测试进程一行没死、`go test` 完整跑到结尾（本腿逐字 `gotest_rc=1`、四行的宿主全是 `testing.go:1617` 的 `Cleanup`），而 310 表② 那三形的分类轴是"**进程怎么死**"（正常收尾／`-timeout` panic／外力 `TerminateProcess`）⇒ 这一条在那根轴上⛔ 占新位；它占的是另一根轴（"**残留怎么产生**"），并与 310 已证的"甲形变体＝重定向根治'落在哪'⛔ 治'没人收尾'"直接接续。
⚠ 但它**证伪**了 310 现量里那句直陈："测试正常跑完 ⇒ 0 残留"**⛔ 全域成立**，只在默认那枚根上成立（票 311 现量第 3 条已把这句话说破，本腿给它成因）。

**落点（建议，⛔ 由本腿盖章）**＝**票 310 的追加节**：那句被射程作废的结论写在票 310 的现量里，作废它的范围必须在**同一张票**上追加（票面原句⛔ 回头改，只追加）；本票 311 留"**这一形叫什么＋触发条件＋成因**"（即本件 §2 ＋ `30-cleanup-failure.md`），裁票 310 那一格的人按本件引用即可。⛔ 在两处各写一份判语（票 311 排程节那句"两张票各自收口、⛔ 混批"）。
⚠ 归口时顺手带上这两条读数，它们决定"这一形要不要单独一格"：① 漏的**枚数⛔ 等于行数**（本腿现量＝6 枚锁死叶 ⇄ 4 行 `TempDir RemoveAll cleanup`，因由＝`removeall_at.go:129` 只带第一枚错、`testing.go:1613` 的清理单位是 tempDir）；② 这一形**自带把用例判红的副作用**（4 行 `Errorf` ⇒ `FAIL=2`），⛔ 只是"留个目录"那么轻。
★ 本腿⛔ 实现任何清理器、⛔ 提议清理动作（票面 `AC#3` 那句"⛔ 顺手实现清理器"在本腿一字未违）。

## 2. 触发条件（一句能判的话）

> 当 `t.TempDir()` 所在那枚专用根的**整条父链给本 token 的继承掩码里没有一枚含 0x40**（`FILE_DELETE_CHILD`＝.NET `FileSystemRights.DeleteSubdirectoriesAndFiles`＝icacls 的 `Dc`），**而**夹具在它下面种／封出一枚 protected DACL、其全部 trustee 都⛔ 在本 token 里 ⇒ `t.TempDir()` 的自我收尾必漏：`os.RemoveAll` 的两次 `NtOpenFile` 都被拒、打出 `TempDir RemoveAll cleanup: openfdat …\la: Access is denied.`、整枚 tempDir 连同锁死叶留在盘上、用例被 `Errorf` 判红。

两枚条件缺一即⛔ 漏（本腿两枚同行证据＝真机量的）：
- 条件一缺（父链给 0x40）＝根 A `%LOCALAPPDATA%…\311a1-A-…`：父链逐枚带 `S-1-5-21-…-1001 … mask=0x1F01FF`（含 0x40）⇒ `cleanup=0`、根内测试残留 `=0` 枚。
- 条件二缺（叶给了本 token）＝根 B 同一发的 `…\TestSealNarrows…1250499198\001`（`PrivateDirAll` 的 `store`，私有集合含本 token）⇒ 那一枚**被正常收尾删掉了**，同发的 `\002`、`\003`（种成 `LA`／`BG` 独占）留下了。
本机那两枚"⛔ 在本 token 里"的 trustee 具名＝`S-1-5-21-1228170099-895614386-1166154857-500`（`DESKTOP-LVS7839\Administrator`，SDDL 的 `LA`）与 `S-1-5-32-546`（`BUILTIN\Guests`，SDDL 的 `BG`）；`S-1-5-18` ⛔ 在 token；`S-1-5-32-544` 在 token 里逐字＝`"Group used for deny only"`。

## 3. 本腿留下的残留目录（逐枚全路径＋`CreationTime`＋哪一发留的）

尺＝`ps1/deep-list.ps1`（逐层 `Get-ChildItem -Directory -Force`，**⛔ 用 `-Recurse`**——`-Recurse` 进不去锁死叶、会把枚数读少：`logs/15-residue-roster-both-roots.txt` 那发就报了 `ps_errors=6`，而它数出的枚数比逐层那把少；逐层尺把每一枚"进不去"具名落盘＝`logs/14-residue-deep-list-both-roots.txt` 的 `CANNOT_ENUMERATE=…` 行，根 B `enumerate_denied_count=7`、根 A `=1`）。第二把尺＝`logs/17`（bash `find -type d`）。读数件＝`logs/31-residue-final-both-roots.txt`（`now=2026-10-12 13:50:27`）。
⛔ 删过任何一枚、⛔ chmod、⛔ takeown、⛔ 动过任何一枚**既有**目录的 ACL。

### 根 B `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594`（换卷那发，`dir_count_excluding_root=20`）

★锁死叶 6 枚（本 token 既进不去也删不掉；根 B 那把尺的 `enumerate_denied_count=7` ＝这 6 枚＋本腿自己的 `plant2-runner`）——**票 311 `AC#3` 那一形的现场**：

| 全路径 | CreationTime | 哪一发／夹具 |
|---|---|---|
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSID3972594964\001\la` | 2026-10-12 13:31:53.203 | `go test`（根 B）`plantedTrusteeSID("D:(A;;FA;;;LA)")` |
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932\001\runner` | 2026-10-12 13:31:53.204 | 同上发 `plantDescriptor(runnerDescriptor)` |
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDa_name_spelling_cannot_2711238538\001\named` | 2026-10-12 13:31:54.750 | 同上发 `plantDescriptor(runnerDescriptor)` |
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDa_name_spelling_cannot_2711238538\001\numeric` | 2026-10-12 13:31:54.751 | 同上发 `plantDescriptor(numeric)` |
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\002\la` | 2026-10-12 13:31:58.659 | 同上发 `administratorSID(t)` |
| `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\003\la` | 2026-10-12 13:31:58.660 | 同上发 `guestsSID(t)`（⛔ 出现在那 4 行里，同一枚 tempDir 只报第一枚错） |

支撑它们的 9 枚 tempDir／序号目录（本 token **能删**，本腿⛔ 删）：
`…\TestGateJudgesThePrivateSetByResolvedSID3972594964`（13:31:53.203）＋其 `\001`（13:31:53.203）；
`…TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932`（13:31:53.203）＋`\001`（13:31:53.204）；
`…TestGateJudgesThePrivateSetByResolvedSIDa_name_spelling_cannot_2711238538`（13:31:54.750）＋`\001`（13:31:54.750）；
`…TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198`（13:31:58.657）＋`\002`（13:31:58.659）＋`\003`（13:31:58.659）。
（那枚 tempDir 下的 `\001`＝`store` **⛔ 在名册里**＝被自我收尾删掉了，见 §2 第二条证据。）

本腿自己的探针 5 枚（⛔ 测试残渣，是 `AC#0` 的取证台；`creation` 见名册）：
`…\311a1-B-…\aclprobe`（13:30:40.913）、`…\plant-runner`（13:40:02.408）、`…\plant-la`（13:40:02.416）、`…\plant2-runner`（13:41:30.695，★锁死叶，`CANNOT_ENUMERATE` 与 6 枚同形）、`…\plant2-la`（13:41:30.700）。

### 根 A `C:\Users\swq\AppData\Local\wisp-test-gotmp\311a1-A-20261012-133040-1594`（默认那发，`dir_count_excluding_root=5`）

★**测试残渣＝0 枚**：该根内 5 枚目录**全是本腿自己的探针**——`aclprobe`（13:30:40.912）、`plant-runner`（13:40:02.270）、`plant-la`（13:40:02.398）、`plant2-runner`（13:41:30.636，★同一枚 SDDL、同 `dacl_protected_flag=True`、同"进不去"）、`plant2-la`（13:41:30.690）。
⇒ 根 A 那发 `go test` 把 4 枚 tempDir 连同其下同形的锁死叶**删干净了**（`logs/20` 全文 `TempDir RemoveAll cleanup`＝0 行）＝问③那条结论最强的一枚反证。

### `%TEMP%` 顶层（⛔ 本腿的残渣，只是名册尺）

起手 13:27:58＝`temp_top_TestRoster_count=59`（`logs/05`）⇄ 交件前 13:50:27＝**59**（`logs/30`）⇒ 本腿两发**没**往 `%TEMP%` 顶层添任何一枚。⛔ 动过那 59 枚中的任何一枚（含邻居留下的）。家族尺同行：`testexe_family_count=0`，正控 `positive_control_svchost_count=85`／`positive_control_all_exe_count=346`（⛔ 空输出＝0 枚那一形：计数在 PS 内部算好再打印，见 `gate-roster-304.md` §1 那条"假 0"事故）。

## 4. 写面·越界读数（⛔ 区间尺，逐笔）

- ★**锚在我跑动期间被编排者推进过一笔**（逐字：起手＝`8e3cb370`，commit 前现读＝`HEAD_now=af5179dd`，中间只有那一枚 `派 311-a1…` 的件）；那一笔**⛔ 碰 `internal/winsec`**（尺＝`git diff --name-only 8e3cb370..HEAD -- internal/winsec | wc -l` ⇒ `winsec_diff_lines_8e3cb370..HEAD=0`），工作树里的 `internal/winsec` 也⛔ 脏（`winsec_worktree_dirty=0`）⇒ **本腿两发吃的仍是同一份码**，"⛔ 改一字"那句成立。全读数＝`logs/42-write-surface-before.txt`（`now=2026-10-12 13:56:05`）。
- 本腿写面**只有两件套**＝`.scratch/wisp/probes/311/a1/**` ＋ 票 311 末尾**一条 bullet（`git diff --numstat` 逐字＝`1	0	.scratch/wisp/issues/311-…md`）**。
- 起跑前逐字记录过的共享工作树在飞态＝`M .gitignore`／`M go.mod`／`M design/**`／`M .scratch/wisp/probes/{152,161,242,268}/**` ＋ 一批 `??`（`-`、`.scratch/.scratch/`、`.scratch/ci-logs/**`、`.scratch/commit-msg-*.txt` …）——本腿**⛔ 提交、⛔ revert、⛔ 修改、⛔ 进任何一笔**。⛔ `git add -A`／`git add .`／`git add <目录>`；`add` 与 `commit` 两处都带**显式 pathspec**（枚名逐条列出，见 `logs/42-write-surface-before.txt`／`logs/43-commit-names.txt`）。
- 每笔的越界尺＝`git show --name-only --format=`（逐枚名册，⛔ 区间）＋`git diff --numstat`（票面那枚文件必须 `1 0`）⇒ 读数落 `logs/43-commit-names.txt`。**★披露（第二笔 commit 的形状）：那枚读数只能在第一笔之后才取到 ⇒ 本腿留了两笔。第一笔＝`5552567c`（32 枚件，`32 files changed, 5171 insertions(+)`，越界枚数尺①逐字＝`0`，票面 numstat＝`1	0`）；第二笔带四枚件＝`logs/32-gate-end-postcommit.txt`＋`logs/33-residue-handover-both-roots.txt`（commit 之后的名册／残留现读）＋`logs/43-commit-names.txt`（第一笔的越界读数）＋**本件这一枚更正后的 `90-final.md`**（更正的就是这一段"第二笔到底带了几枚什么"的自述——第一笔里我把它写成"只带 `logs/43` 一枚"，与实际三枚读数件不符）。第二笔⛔ 带新结论、⛔ 碰票面、⛔ push；两处 `add`／`commit` 同样带显式 pathspec。**
- ⛔ push（本腿一次都没 push）、⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`、仓内⛔ 建 worktree。
- 门禁类读数（本腿**只跑只读尺**，⛔ 跑那三枚在飞入口件＋`tools/d22scan/runtests.sh`）：
  - `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ 逐字 `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`，`len_gate_rc=0`（件 `logs/40-length-gate.txt`）。
  - 本腿新建 29 枚件自身的长度尺（`logs/41-my-path-lengths.txt`）⇒ 最长相对名 `72`／最长末段 `39`／绝对前缀 `/d/work/workspace/projects plans/Wisp` ⇒ 全路径 ≤`109`，三档预算（121／100／165）都在里面。
  - ⛔ 跑 `GOFLAGS= go build ./...`、`sh scripts/d22scan.sh`、`bash scripts/portable-tests.sh --scope=census`、`bash scripts/winsec-tests.sh …`：那四枚是**落地腿 `AC#4` 的台面**，且其中三枚此刻归 `310-r2` 在飞（借它跑＝拿别人的在飞态当台面）。本腿跑的是 `GOFLAGS= go test -v -count=1 ./internal/winsec/` 两发（＝票面 `AC#0` 指定要的对照），⛔ 拿它冒充 `AC#4`。
- 件面纪律：⛔ 长名、⛔ 文件名带 `§`、⛔ `.out`、`probes/**` 下⛔ `.go`（本腿 5 枚 `.ps1` 都在 `probes/311/a1/ps1/`，⛔ `.go`）、⛔ 0 字节尺件（每枚自落 `ps_rc=`／`gotest_rc=`／`len_gate_rc=` 行，最小一枚 `logs/00-runids.txt`＝201 B）；"应当为 0"的尺一律 `|| true` 并把**命中数与 stderr 并排**落盘（`logs/20`/`logs/21` 各含全文，`grep -c` 的数写在交付件里）。★读数件⛔ 写进任何专用根里（两枚根内只有 `t.TempDir()` 产物＋本腿探针目录，⛔ 一枚 `.txt`）。

## 5. 自报（⛔ 自己盖章说"没影响"）

1. **我跑了两发真测试，改的是这台机的盘上状态**：`go test -v -count=1 ./internal/winsec/` 各一发（`gotest_rc=0`／`1`），它们**真的** `icacls /grant`、`SetNamedSecurityInfo`、封目录、建删目录——只是全部落在我自己新建的两枚根里（根 A 里 4 枚 tempDir 已被测试自己删掉；根 B 里留下 §3 那 15 枚）。这些射击⛔ 动任何既有对象的 ACL，但**⛔ 等于"对台面零影响"**：枚数、`%TEMP%` 名册 59⇄59、家族尺 0⇄0 是我能给出的边界，⛔ 超出这个边界替任何人担保"没影响"。
2. **我在两枚根里各造了 5 枚探针目录，其中 2 枚（`plant2-runner`）是我自己造出来的锁死叶、本机删不掉**（本 token⛔ 能列举）。它们**⛔ 是测试残渣、⛔ 是票 310/311 现量名册里的东西**，后续任何一次"根内条目数"作差都会把它们算进去——落册在 §3，⛔ 删、⛔ chmod。⛔ 拿它们冒充任何一形残留。
3. **我用了两枚"代用尺"而不⛔ 是 wisp 的码路**（本机没有别的工具能取"根 A 的封后逐 ACE"这一格）：`ps1/plant-probe.ps1`＝PowerShell `Set-Acl`（Win32 `SetNamedSecurityInfo`）。其一（带 `P` 的 `D:PAI…`）与真码残留**逐字同形**、可用；其二（`D:(A;;FA;;;LA)` 那枚，不带 `P`）⛔ 复现 `PROTECTED_DACL_SECURITY_INFORMATION`，落出 4/9 ACE 的"仍带继承"形 ⇒ **那一枚读数只算半根棍子**，`10-acl-before-after.md` §3 与 `30-cleanup-failure.md` §5 都按这个口径标注。第一次跑那枚脚本还失败过一回（枚举名打错：`SecurityDescriptorSections` ⛔ 存在，正确名＝`AccessControlSections`），失败读数留在 `logs/19-plant-probe-both-roots.txt`/`logs/19b-*.txt`，那 4 枚 `plant-*`（⛔ `plant2-*`）目录是**那次失败的副产品**＝默认继承 ACL、未被种死——已落册、⛔ 删。
4. **`Get-Acl`/`icacls` 之外我⛔ 用任何特权工具**：⛔ `Set-Acl` 作用于既有目录、⛔ `takeown`、⛔ `icacls /reset`、⛔ `chmod`。⛔ 杀进程（本腿⛔ 起常驻进程；`go test` 的子进程自己退完了——正控读数＝家族尺 0 枚）。
5. **我⛔ 翻框、⛔ 碰 `Status:`、⛔ 改名 `-done`、⛔ 产码、⛔ 动判据**：`internal/winsec/**` 与 `scripts/**`、`.github/**`、`docs/**` 一字未改（唯一动过的票面文件＝票 311 **末尾追加一条 bullet**）。我没有裁 `AC#1` 三形（⛔ 本腿射程），只在 `20-red-sentence.md` §4 与 §1 落了"甲形够不着红因"这**一枚读数＋一句成因标注**，供编排者裁形时⛔ 把"判据换了答案"当既成事实。
6. **一处我⛔ 交的东西**：⛔ 演示"同一枚锁死叶在根 A 上删得掉"（那要真删一次，⛔ 票面 `AC#0` ⛔ 删任何东西）⇒ 该结论在 `10-acl-before-after.md` §3 里是**间接**的（靠根 A 残留＝0 枚 ＋ 两枚 `plant2-runner` 同形），标的是〔读码推的〕⇄〔真机量的〕混合，⛔ 当成"已演示"。
7. 一次自曝的读数以"我先前少读了一处"为形态：`logs/15` 那把 `-Recurse` 尺我一开始当成"某枚 `003\la` 不存在"，重测（bash `ls -1` ＋ 逐层 `Get-ChildItem`）才知它**在**、只是 `-Recurse` 进不去——已把这条写成尺的告警（§3 第一行 + `logs/14`），⛔ 让第一次那个错读数留在结论链里。
8. 留给后续程的一条已知坑（**⛔ 本腿的读数，是转抄**）：`bash scripts/winsec-tests.sh --scope=winsec` ⇒ 被它自己的 GUARD 1 打回 `rc=2`（正确形 `bash scripts/winsec-tests.sh ./internal/winsec/`）；出处＝票 311 `:12` 与编队正典 §1，本腿⛔ 用那枚脚本、故⛔ 复跑那条 `rc=2`。
