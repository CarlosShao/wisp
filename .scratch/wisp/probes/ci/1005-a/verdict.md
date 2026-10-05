# verdict.md — run 37257547572 CI 四问取证（只读腿 ci-attr-1）

交付件＝本文件。取证者＝`ci-attr-1`（只读＋落这一件）。【禁】 全程零 `go` 命令；只用 `gh`／`git`／`grep`／`sed`＋文件读。【禁】 台账／工单／票面 AC 框一字未动；`frontend/**`、`design/**` 未读未引（下文引用它们的唯一途径是 CI 日志里 `internal/panel` 测试自己打出的行，非直读其文件）。

## §0 起手锚与时刻（读数带时刻）

- 起手现量时刻：`2026-10-05 03:12:58Z`（`date -u`）。
- 起手 HEAD：`c6cf66e64849955bf92a4b3356c096ea81c35563`，分支 `dev`。
- 起手远端：`origin = https://github.com/CarlosShao/wisp.git`（GitHub Actions 所在）。`gh auth`＝CarlosShao，scopes 含 `repo`。
- 起手工作树脏度：`git status --short | wc -l`＝**673**（同机三枚写腿的脏增量，不是我的）。我的写面起点＝`.scratch/wisp/probes/ci/1005-a/` **不存在**。
- 终态写面（本腿）＝仅 `.scratch/wisp/probes/ci/1005-a/`（`verdict.md` 提交入库＋`logs/*` 证据留盘不入库），其余 672 项脏量归他人，未被我 staged。
- 取数口径（命令／时刻／来源）：
  - job 台账＋步级 conclusion：`gh run view 37257547572 --json jobs ... --jq`（03:15Z 前后）；来源 run `37257547572`。
  - 各 job 完整日志：`gh run view 37257547572 --job <id> --log`（lint `111597759814` @03:18:47Z；test-windows `111597759830`、test-core `111597759928`、slo-full `111597759876` @03:20:17-20Z）。
  - 全 run `--log-failed`：@03:18:51Z（口径 A，历史少算红名，见 §9）。
  - 红名册核一遍＝完整日志口径（口径 B）；两套都在盘（`logs/`）。
  - D32 两数：`gh run download 37257547572 --name slo-full-report`（@03:22:38Z，落到 `logs/artifact/slo-report.json`，645074 字节）。
  - 上一枚 push run `37249563077`（headSha `21bec8a1`）作差：完整日志缓存 `logs/prior-test-windows-full.log`／`logs/prior-test-core-full.log`。
  - 本批范围：`git rev-list --count 21bec8a1..HEAD`＝**53**（`派单文字"145＋枚"＝编排者口径，可能含更早未单列的 push；本腿只报现量 53`）。

## §1 run 与 job 台账

- run `37257547572`｜event=push｜headSha `c6cf66e6`｜conclusion=**failure**｜created `2026-10-05T02:59:13Z`｜updated `2026-10-05T03:08:54Z`｜url `https://github.com/CarlosShao/wisp/actions/runs/37257547572`。
- 六枚 job **全部 `status=completed`**：

| job | databaseId | conclusion | started→completed | 步数 | 其中 skipped（非 Post）|
|---|---|---|---|---|---|
| slo-smoke | 111597759704 | success | 02:59:15→03:01:33 | 11 | 0 |
| lint | 111597759814 | **failure** | 02:59:15→03:00:55 | 16 | 3（denominator / go vet module / go vet d22scan）|
| test-windows | 111597759830 | **failure** | 02:59:17→03:08:53 | 13 | 0（6 枚带 `if:!cancelled()`，均出结论）|
| lint-frontend | 111597759837 | success | 02:59:15→02:59:37 | 15 | 0 |
| slo-full | 111597759876 | success | 02:59:16→03:02:25 | 9 | 0 |
| test-core | 111597759928 | **failure** | 02:59:15→03:01:09 | 11 | 0 |

- 步级 conclusion 逐枚（口径 B，完整日志＋`--json jobs`）：
  - **lint（16）**：Set up/checkout/setup-go＝success；`D22 positive control`(:74)＝**success**；`D22 self-test`(:83)＝success；`D22 seven-ban+emoji`(:108)＝success；`path-length budget`(:136)＝success；`gofmt (gofumpt)`(:168)＝**failure**；`denominator AC#7 form A`(:184)＝**skipped**；`go vet (module)`(:206)＝**skipped**；`go vet (d22scan)`(:209)＝**skipped**；`staticcheck`(:213)＝**failure**；`mockllm module vet`(:289)＝success；`Post setup-go`＝skipped；`Post checkout`/`Complete`＝success。
  - **test-windows（13）**：winsec(:430)／cache(:473)／cgo build smoke(:480)＝success；`cmd/wisp CLI tests`(:487)＝**failure**；`Portable windows tests`(:509)＝**failure**；`PathResolver junction`(:545)＝success；Post cache／Post setup-go＝skipped。
  - **test-core（11）**：Start compose／Probe mock-llm／env fork 断言＝success；`Portable package tests`(:351)＝**failure**；Stop compose(:375, `if: always()`)＝success；Post setup-go＝skipped。
  - **slo-full（9）**：setup／checkout／setup-go／Build wisp.exe／SLO full gate／Upload SLO report／两 Post／Complete＝**全 success，零 skipped**。
- 判：三枚红 job 各带**红步骤**（不是被 skip），红步骤都真产日志；唯一"无日志"的实质格在 lint 的 3 枚 skipped 步（见 §6）。

## §2 lint 为什么红

lint 有**两枚独立红步**，非"一步吃掉"造成的假象（二者都出了结论）：

**(A) `gofmt (gofumpt)`＝failure（链头红，ci.yml:168）**
- 红句逐字（`lint-full.log:788`）：`##[error].scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`；随后 `:789` `##[error]Process completed with exit code 2.`
- 点名文件：**只有 1 枚**＝`.scratch/wisp/probes/185/c1/mut/fs_broken.go`（票 185-c1 的突变种子件，故意让 import 落在声明之后）。步骤脚本 `OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"`（ci.yml:171）里的 `.` **会走遍 `.scratch/**`**，于是盘上一枚被跟踪的坏样本直接把这道门判红。
- **本批引入 vs 预存**：**预存**。`git log --diff-filter=A -- .../fs_broken.go`＝该文件首入库于 commit `4813567e`（`2026-09-28T14:42:58+08:00`），本批范围 `21bec8a1..HEAD` 未碰它（`git log ... -- <file>` 在区间内空）。它今天才**第一次被看见**，是因为本批修掉了它上方那枚链头红（见下），于是 `gofmt` 步从历史 `skipped` 翻成 `failure`——`gh run view` 逐步回看：`37249563077 / 37240161874 / 37166458550 / 37158259050 / 37076437856 / 37021179942` 六枚里 `gofmt` **全是 skipped**（被上方红吃掉），本次是第一枚让它出结论的 run，结论＝红。
- 连带后果：`gofmt` 变成**新的链头红**，把它下方三枚未带 `if:!cancelled()` 的步重新吃回 skipped（`denominator`/`go vet module`/`go vet d22scan`，`--json` 现证＝skipped）。这与 ci.yml:178-204 那段"AC#7 form A 分母步"的自陈相互咬合：分母步本是为绕开"gofumpt 走 .scratch 恒不可满足"而加，但它**不带守卫**，一旦 gofmt 先红它就采不到读数——本腿判这属**门禁形状缺陷，宜立案**（§7）。

**(B) `staticcheck`＝failure（ci.yml:213，带 `if: ${{ !cancelled() }}` 于 :263）**
- 红句：三模块逐 exit=1，`lint-full.log:873/878/881/882`：
  - `module .: exit=1 packages=34 findings=49 toolchain-crash-lines=0`
  - `module tools/d22scan: exit=1 packages=1 findings=3 toolchain-crash-lines=0`
  - `module tools/mockllm: exit=1 packages=1 findings=1 toolchain-crash-lines=0`
  - 汇总 `:882` `staticcheck self-report: version=staticcheck 2026.2.1 (0.8.1) modules=3 packages=36 findings=53 toolchain-crash-lines=0 (step exit is 1)`
- 点名（逐条 file:line 皆在日志 824-872）：如 `cmd/wisp/approval_reply.go:91`（U1000 未用）、`frontend/embed.go:4 SA9009`、`internal/observe/thresholds.go:42`、`internal/tools/platform_other.go:39` 等 49 条根模块 U1000/ST/S/SA 族。`toolchain-crash-lines=0`＝不是导出格式崩，是真·可解析积压。
- **本批引入 vs 预存**：**预存**。ci.yml:260-262 自陈"This step is EXPECTED RED after this pin ... lint turning green is NOT a criterion of this ticket"，且门的历史定性＝`#65`/票 122 积压族。本批 `tools/**` 未被点名改动，findings 数（53）与上一枚能采到 staticcheck 的 run 属同一族。判：**具名归口票 122**，非本批新增量。【注意】〔仅转述〕台账 `A?` 那句"lint 唯一剩红因是 staticcheck＝积压"与此一致，但本腿按现日志独立复算 53 条得同判。

判语：lint 红＝(A) 预存坏样本 `fs_broken.go` 撞 `gofumpt -l .` 全树扫描（本批只是把它从 skipped 解封成可见红），＋(B) 预存 staticcheck 积压 53 条。二者都非"本批代码新造的红"，但 (A) 暴露一个可立案的门禁形状缺陷。

## §3 test-windows 为什么红（逐名，【禁】 不笼统说"环境"）

红 job 由**两枚红步**构成：`cmd/wisp CLI tests`（包 `github.com/CarlosShao/wisp/cmd/wisp`，`test-windows-full.log:3837` `FAIL ... 376.711s`）＋`Portable windows tests`（其中包 `internal/risk` `:4966` `FAIL ... 17.303s`；同步骤内 `winsec/proc/secret/config/ball/perm/plugin/cmd/llmrecord` 皆 `ok`）。

**先排掉本机历史坑**：全 log `grep -c "0xc0000135|0xc000013a|3221225794"`＝**0**。且 `cmd/wisp` 用例都带真实耗时与 `--- FAIL`（非"0 断言＝没跑"），cgo build smoke 步 success，`slo-full` 里 `wisp doctor` 的 sherpa 1.13.8／onnxruntime 1.28.2.0 均"matches build pin"。⇒ **DLL 未缺，用例确跑**，红是真跑出来的断言红。

**cmd/wisp 7 枚（逐名＋归因）：**
1. `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（:1969）——`always_write_no_clobber_226_test.go:91`：盘上 `fs.allowed_dirs=[RUNNER~1/.../002, runneradmin/.../003]`，want 存的一条。归因＝runner 用户目录 **8.3 短名 `RUNNER~1` ↔ 长名 `runneradmin` 双拼写**，非被测逻辑错。
2. `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`（:2202）——同 226 族，`allowed_dirs=['...RUNNER~1.../002','...runneradmin.../003']` 双拼写。同上因。
3. `TestTicket223HandEditedFsLooseningCostsAnL2Card`（:2394）——`config_reload_223_test.go:312/315`："the console did not render the reload card"、"the card does not name the key it is asking about"（拿到的仍是通用"答复监听已接入…"提示）。**当前名册独有、上一枚 push 无此名**（见作差）。【注意】 本批改了 `cmd/wisp/firstrun.go`＋5 枚测试件，此路 reload-card 与 firstrun/run 相关——**唯一"本批可能造成"候选**，但【禁】 禁跑 `go`，不能证因（§8）。
4. `TestTicket223PermissionDeniedSitsInItsOwnSentence`（:2501）——`config_reload_perm_223_windows_test.go:95`："icacls denied nothing, so this case cannot show the permission sentence"。归因＝runner 账户权限下 `icacls` 造不出拒绝（环境耦合，具名）。
5. `TestRunPacketCarriesTheLoadedInstructionFiles`（:2676）——`instructions_200r2_test.go:167`：packet 记的 Path 是 `c:\users\runneradmin\...`，want `C:\Users\RUNNER~1\...`。同 8.3 短名族。
6. `TestPanelHostRealWindowHopAndLifecycle`（:2871）——`panel_host_windows_test.go:660/665`："cold bring-up 3883.2 ms exceeds D32 panel cold budget 1500 ms"（该机 WebView2 冷启 3883ms）。归因＝面板冷启延迟超 D32 预算，性能/机器相关。
7. `TestPanelHostLatencyPercentilesAC2`（:2881）——`panel_host_windows_test.go:1005`："cold P95 3883.234 ms ... exceeds the D32 panel cold budget 1500 ms"（n=1）。同上；且编排者本人在 `A616` 承认这枚分位数行会被并发采样污染（环境敏感，具名）。

**internal/risk 12 枚（逐名＋归因，两类机制）：**
- 8.3 短名族（6 枚）：`TestCanonicalInputGainsNoSecondForm`（`pathresolver_anchor_spelling_windows_test.go:140`，一个拼写得 2 比对形 `runner~1`/`runneradmin`）、`TestPathResolverShortNameAListDenied`（`pathresolver_junction_windows_test.go:135` got `runneradmin` want `RUNNER~1`）、`...UNCAListDenied`（:164）、`...ExtendedLengthPrefixAListDenied`（:182）、`TestClassifyAnchorSpellingIsNotVerdict`（:4550 锚拼写族）、`TestAListWinsWhereBothTablesHit`/`TestBListDefaultDenyAndOverride`（`pathresolver_junction_windows_test.go:292` "B-list single-file override must allow"，同一短名下类判定漂移）。
- SyncDir 无同步根族（5 枚）：`TestSyncFixtureFallbackAndMatch`／`TestSyncFallbackNotDisarmableByWeakRoot`／`TestSyncUnverifiedRootKeepsFallback`／`TestSyncSuspectFallbackIsComponentBounded`／`TestSyncSuspectFallbackWhenUndetectable`——`syncdirs_test.go:402` "under-profile path must be sync-suspect, got Sync:false"、`syncdirs_test.go:494` **"P12 evidence: no sync location detected on this machine"**。归因＝该 runner 无 OneDrive/Dropbox 同步根，sync-suspect 分类落不到期望分支（环境耦合，具名机制，非"泛环境"）。

**作差（本批 vs 预存）：** 上一枚 push `37249563077`（缓存完整日志）红名册 19 枚，本枚 19 枚，交集 18。
- 仅本枚新增：`TestTicket223HandEditedFsLooseningCostsAnL2Card`。
- 上一枚有、本枚已绿：`TestResolvePerCallBudget`（本枚日志 `=== RUN` 后 `--- PASS`，故翻色）。
- `internal/risk`／`internal/panel` 在本批 `21bec8a1..HEAD` **改动文件数＝0**（`git diff --name-only` 现量）。⇒ 12 枚 risk 红**逐枚预存**（同名同因在两枚 run 都在册），cmd/wisp 另 6 枚亦在册＝预存。cmd/wisp 被本批改 6 枚文件（`firstrun.go`＋`approval_reply_201_test.go`／`approval_seam_201_test.go`／`firstrun_257_test.go`／`panel_pump_test.go`／`run_mode101_test.go`／`ticket224_assembly_test.go`——多为 `_test.go`），但改动的测试名与红名不重合；红名里唯一本批相关嫌疑＝#3。

判语：test-windows＝**18/19 预存**（短名 `RUNNER~1`/`runneradmin` 与"无同步根"两类 runner 机器耦合＋WebView2 冷启超预算＋icacls 造不出拒绝），**1 枚（Ticket223HandEdited）本批嫌疑待归口**；零枚命中"缺 DLL 静默不跑"坑。

## §4 test-core 为什么红（ubuntu 侧逐名）

红步＝`Portable package tests`（`scripts/portable-tests.sh` 的 core 范围）。红包＝`github.com/CarlosShao/wisp/internal/panel`（`test-core-full.log:5000` `FAIL ... 1.027s`），步末 `:5428` `##[error]Process completed with exit code 1.`。同步骤 `internal/proc`＝`ok`。【禁】 这是 ubuntu runner，无 windows 短名/同步根因，别套 §3。

4 枚 `--- FAIL` 逐名＋归因：
1. `TestApprovalCardViewJSONKeysMatchFrontendTypes`（:4239）——`approval_test.go:129`："Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare"。归因＝Go `Snapshot` 发的 JSON 键多于前端 `PanelSnapshot` 声明（契约漂移）。
2. `TestComposerContractTypesMatchFrontend`（:4396）——`composer_test.go:74`：同上 Snapshot 键，另"Go ComposerState emits [git currentModel modelKnown credentialState credentialKnown] that interface ComposerState does not declare"。同族 Go↔前端接口键漂移。
3. `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`（:4464）——`frontend_hygiene_test.go:216`："a second style source appeared in code the bundle actually ships: ...right-rail.tsx:89 ... `--page: #fafafb;`"，`Add or reference a C21 token in design/assets/tokens.css instead`。归因＝打包可达代码里出现第二处颜色字面量（本应只在生成主题里）。
4. `TestC21DesignTokensFourWayAgree`（:4985）——`tokens_fourway_test.go:532` 连打多行"design/assets/tokens.css declares `--orb-shadow`/`--tint-*-hi|lo` but frontend/src/styles/tokens.generated.css does not carry it"，`:547` **"only 0 of 78 colour rows completed all four legs - the check is not doing its job"**。归因＝tokens.css 与 tokens.generated.css 令牌集不同步，四向一致门 78 行 0 行通过。

**本批 vs 预存**：上一枚 push `37249563077` 的 test-core 完整日志里**同名 4 枚全部 `--- FAIL`**（`prior-test-core-full.log` grep 现证：`FAIL github.com/CarlosShao/wisp/internal/panel 0.853s` ＋ 4 个 `--- FAIL: Test(ApprovalCardView...|C21...|Composer...|PanelColour...)`）。且 `internal/panel`／`design/**`／`frontend/**` 本批改动文件数＝0。⇒ **test-core 红＝预存（契约/令牌漂移），非本批引入**。
【注意】 口径教训复现：本腿**第一次**直接管道 `gh ... --log | grep '--- FAIL'` 抽上枚 test-core 名册返回**空**（少算）；改成先 `--log >` 缓存成盘再 grep 才拿到 4 枚真名册。两套口径都留 `logs/`（§9）。

## §5 【重点】 slo-full 这一枚 success 到底测了什么

来源＝run `37257547572` job `111597759876`（全 success，9 步 0 skipped），＋其上传 artifact `slo-full-report`（ID `11323476336`，19945 字节 zip，解出 `slo-report.json` 645074 字节，@03:22:38Z 下载）。

**它走了"真采样"支，不是"拒采仍绿"支**（这是本题要害）：
- `slo-check.ps1` 采样有效性预检（票 134 AC#4）通过：日志 `:214` `precheck ok - no foreign toolchain/runner process, machine-wide cpu max 45%`。⇒ 机器未被争用，**未**走 `NO CONCLUSION (machine-contended)`＋exit 0 那条（ci.yml:602-621、:674 描述的拒采支）。
- 六态逐一实测：`:217/220/223/226/229/232` `state Sleeping/Armed/Warm/Conversation/PanelOpen/WorkPeak exit=0 pass=True`；`:235` `settle exit=0 pass=True`；正控 `:236-238` `leak fixture self-test (100MB, must FAIL)`→`leak exit=1 flipped_to_fail=True`（证门有牙，能判红）。
- `:239` `report written ... slo-report.json (all_pass=True)`；report 顶层 `all_pass=true`、`sample_errors` 全程 0。

**D32 那两个数（空闲内存／CPU）——逐字取自 `slo-report.json` 的 `Sleeping`（深睡/空闲态）out-of-tree 门行：**
- 空闲内存 `tree_private_bytes`＝**`"measured": "4.8MB"`**，`"limit": "<=25MB"`，`"pass": true`。
- 空闲 CPU `cpu_percent_all_core`＝**`"measured": "0.043%"`**，`"limit": "<=0.5%"`，`"pass": true`。
- 口径注：同态另有一行 `8.7MB / 0.043%` 属 **in-tree** 读数，report 自己标 `"note":"recorded, not a gate: ... the observer's own ReadTree cost ... the gate is the out-of-tree row"`（`slo-report.json:922` 类）。⇒ 权威门数＝**4.8MB／0.043%**（in-tree 的 8.7MB 只作记录，非判据）。
- 其余五态阈随态递增（内存 `<=90MB / <=350MB / <=350MB / <=600MB / <=700MB`，CPU `<=2.0% / <=1.0% / <=5.0% / <=10.0% / <=30.0%`），实测皆 `pass=true`。
- settle 块（`slo-report.json` 5747 起）：`mode=settle`、`peak_bytes=13188714`→`final_bytes=9334784`（回落≈8.9MB，`pass=true`）；`leak_fixture.flipped_to_fail=true`；`free_os_memory_requested=true / count=2`（漏样正控取内存）。
- 采样反空转门：末行 `"metric":"sampling","measured":"40 valid / 0 errors","limit":">=1 valid sample and sample_errors == 0","pass":true,"gate":true`，`note:"fully measured: sample_errors=0, 40 valid reads kept of 40 reads taken"`。

**覆盖行门（`settleCoverageRowGates`，票 136 AC#14）这次有没有真把某行判成 fail 位？**
- report 里 `gate/metric/measured/limit/pass` 五元组共 **89 行**（每状态 out-of-tree＋in-tree 各一组，含 `gdi_objects/user_objects/handles/goroutines/threads/disk_write_ops/tcp_connections/sampling` 等）。逐行扫：**89 行 `pass` 全＝true，无一被翻成 fail 位**。
- 但这是**真求值后的 pass**，不是"绿＝什么都没测"：每行都带实测量纲对实阈（如 `handles 200 <600`、`goroutines 1 <=6`、`tcp_connections 0 ==0`、内存/CPU 如上），且 `sampling` 门明记 `40/40 valid、sample_errors=0`；再加独立正控 `leak 100MB must FAIL → flipped_to_fail=true` 证明"能判红"的通路是活的。⇒ **无 `NO CONCLUSION (machine-contained)`＋exit 0 的"绿即空测"形状**（该形状的先例见台账：09-24 同日 21 枚里真取样 11、拒采仍绿 10；本枚落"真取样"侧）。

判语：slo-full 这枚 success＝**一次有效的 D32 真结论**——空闲 **4.8MB／0.043%**（阈 25MB／0.5% 内），六态＋settle＋sampling 门全实测通过，无拒采空转。〔仅转述，未独立复算〕台账 `A609` 称"票 263 结案拿到第一份真结论、这次是第二枚"——本腿只坐实"本枚＝真结论"，序数（是否恰为第二枚）需回扫全部 slo-full 历史绿 run 才能证，本腿未逐一回采，标此存疑。

## §6 七步不再 skipped——判定

"七步"＝被链头红吃掉的那七枚 lint 测试步（`A608` 原写"五步"被只读腿顶回为"七步"并点名漏了 `:136 path-length`；本腿照 ci.yml 文本复认七枚为 `:83 :108 :136 :168 :184 :206 :209`，且与上枚 run 的 skipped 名册两边一致＝**待验断言，本腿独立核**）：

- 上一枚 push `37249563077`：链头红＝`D22 scanner positive control`(:74)＝**failure**，其下**七枚全 skipped**（:83/:108/:136/:168/:184/:206/:209），仅 `staticcheck`(:213,`!cancelled()`)＝failure、`mockllm`(:289)＝success 出了结论。
- 本枚 `37257547572`：:74 已修成 **success**（本批把它解封），于是七枚里——
  - **解封、各出结论（4 枚）**：:83 success、:108 success、:136 success、:168 **failure**（真红、真日志）。
  - **仍 skipped（3 枚）**：:184 denominator、:206 go vet module、:209 go vet d22scan——因 `gofmt`(:168) 新成链头红，而这**三枚未带 `if:!cancelled()`**（ci.yml 无守卫），故被再吃。

**⇒ "七步不再被 skipped"这条只在"链头从 :74 挪走"意义下部分兑现：4/7 现已产出裁决（含把 gofmt 的真红读出来），但 3/7 仍 skipped、仍无日志。** 台账 `A616`（line 12044）那句"六枚 job 全部 status=completed、没有一枚被前一步的红吃掉"在 **job 级为真**（6 job 皆 completed），在 **step 级不真**——lint 三枚实质步本枚仍 `skipped`。这是进步（读不到→读得到且红）而非终态：`go vet`×2 与分母步今天的 verdict 仍采不到（见 §8）。

## §7 本批 vs 预存的分账

**预存（应具名归口，不新立案）：**
- lint `staticcheck` 53 findings → 票 122／`#65` 积压族（ci.yml:260 自陈 expected red）。
- lint `gofmt` 红于 `fs_broken.go`：坏样本预存（09-28 `4813567e`），本批仅解封可见——**根因预存、可见性本批**。
- test-core `internal/panel` 4 枚（键漂移/颜色字面量/令牌四向 0-of-78）→ 前端令牌与 C21 契约治理面（票 77／C21 线）；上枚同红、本批未碰 panel/design/frontend。
- test-windows 18/19（risk 短名＋无同步根、cmd/wisp 短名＋icacls＋WebView2 冷启超预算）→ runner 机器耦合，两枚 push 同册。

**本批嫌疑（应立案/点名，非"泛环境"）：**
- **test-windows `TestTicket223HandEditedFsLooseningCostsAnL2Card`＝唯一当前独有红**（上枚无、本枚有），且本批改了 `cmd/wisp/firstrun.go`＋相邻测试；reload-card 未渲染的机制与 firstrun/run 相关。归口＝**cmd/wisp / 票 223 落地腿**，须一次能跑 `go test` 的名册腿坐实是回归还是带载偶发（本腿禁 `go`，见 §8）。
- **门禁形状缺陷（lint）＝建议新立案**：`gofumpt -l .` 扫全树把 `.scratch` 里被跟踪的突变种子件当红，且它一旦红就把无守卫的 :184/:206/:209 再吃成 skipped；与"AC#7 form A 分母步本意绕开此结构性不可满足"自相矛盾（分母步不带 `!cancelled()`＝解封后照样采不到）。具名：ci.yml:168（走 `.` 的 gofmt）＋:184/:206/:209（无守卫）。这是 `A608/A609` 那笔"七步"账的续集，值得独立成票。

**该立案的红名清单（给编排者）：** ① `gofmt`+分母/go vet×2 的门禁形状（解封后仍 3 步采不到，且被一枚 `.scratch` 种子件长期顶红）；② `TestTicket223HandEditedFsLooseningCostsAnL2Card` 需具名判定回归-or-偶发。其余预存红只具名归口、不重复立案。

## §8 量不到（本腿能力边界，诚实登记）

- 【禁】 `go` 全禁 ⇒ 不能本地复跑任何包取逐名册，也不能力证 §7②那枚是"回归"还是"带载偶发"（只能靠 run-vs-run 名册差集＋file:line 机制推断）。
- lint 三枚 skipped 步（:184 分母 / :206 go vet module / :209 go vet d22scan）**无日志**⇒今天它们的真/假/崩**采不到**（这正是"skipped 不产日志＝读数永久采不到"那格）。只能确定"本枚未给出裁决"，不能确定其内容。
- D32 结论的**序数**（"第二枚真结论"）未独立复算：需回扫全部历史 slo-full 绿 run 的 `slo-report.json` 是否存在，本腿只验了本枚为真结论。〔仅转述 `A609`〕
- `slo-check.ps1` 的 `settleCoverageRowGates` 源码未直读（属 `scripts/**` 可读，但本腿按最小取证只取了 report 产出的 89 行门读数来判"有没有真把某行判成 fail 位"，未逐行回推函数体）。
- 起手 HEAD 上另两枚写腿的包名册读数与 `frontend/dist` 内容（是否只有 `.gitkeep`）未取——`§4` 令牌四向 `0 of 78` 的**根因**（是 design/frontend 漂移还是 CI 未 build dist）本腿只报机制、未追至归属；`frontend/**`、`design/**` 依约不读。

## §9 自我对抗（对本判语的反向攻击）

1. 攻"七步解封 4/7"：会不会 :136 path-length 的 success 是假绿（本批把 262 的门放宽了）？→ 本腿查 `git diff --name-only 21bec8a1..HEAD -- scripts tools`＋ci.yml:136 未加 `if:`／`continue-on-error`，:108 七禁＋emoji 步亦 success，无放宽证据；判 :136 success 为真。
2. 攻"红名册 19 枚完整"：口径 A（`--log-failed`）本会少算红名——本腿用**完整 `--log`**（口径 B）逐包抽 `--- FAIL`＋`FAIL\tgithub` 双向对齐，prior 与 current 各 19 枚；且 §4 亲历"直管道抽前枚 test-core 返回空、缓存后得 4 枚"，据此把口径 A 判为不可信、结论一律以口径 B 为准。
3. 攻"cmd/wisp 全预存"：本腿只证"18/19 名在两枚 push 同册"，未证"同因"；已把 #3 `Ticket223HandEdited` 单列为本批嫌疑而非归入预存，未犯"笼统环境"错。risk 12 枚虽同名，本腿仍逐枚给了短名/无同步根的 file:line 机制，不以"windows 环境"一句带过。
4. 攻"slo-full 绿＝真测"：反向假说是"拒采仍绿"。本腿以 `precheck ok cpu max 45%`＋`40 valid/0 errors`＋`leak flipped_to_fail=True`＋`all_pass=true`＋report 落地(artifact 存在) 五点共同排除拒采支；其中 report artifact 的**存在本身**即"产了数才上传"（ci.yml:616-619 逻辑），拒采支不产 report。判真测。
5. 攻"staticcheck 53＝本批"：若本批新增了未用符号会推高 findings。本腿见 `tools/**` 未改、根模块红名多为历史 `_test.go` U1000 与既有 SA 族，且 ci.yml 自陈 expected red；判预存。【注意】 但"53"是**本枚实测**，与 ci.yml:231 那句"34/78"旧读数不同，故不引旧数、只引本日志数。

## §10 判语

1. **lint 红＝两枚真红步**：`gofmt` 因 `.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1`（预存被跟踪突变件，本批解封其可见性）而红，并新成链头把 `denominator`/`go vet`×2 吃回 skipped；`staticcheck` 因预存 53 条积压红（`toolchain-crash-lines=0`，具名归口票 122）。**本批未新造 lint 的代码红**，但暴露一处应立案的门禁形状缺陷（§7）。
2. **test-windows 红＝19 枚**：cmd/wisp 7 ＋ internal/risk 12，**18/19 预存**（短名 `RUNNER~1`↔`runneradmin`、无同步根、icacls 造不出拒绝、WebView2 冷启 3883ms 超 D32 1500ms），**1 枚 `TestTicket223HandEditedFsLooseningCostsAnL2Card` 为本批嫌疑待归口**；`0xc0000135` 类"缺 DLL 静默不跑"坑**未触发**（用例皆带耗时与断言）。
3. **test-core 红＝4 枚 internal/panel**（Go↔前端契约键漂移、颜色字面量、令牌四向 `0 of 78`），与上一枚 push **逐枚同名＝预存**，非本批引入；非 windows 因（ubuntu 侧独立成立）。
4. **slo-full success＝有效 D32 真结论**：空闲 **4.8MB／0.043%**（阈 <=25MB／<=0.5%），六态＋settle＋`sampling 40/40`＋`all_pass=true`，正控 `leak→flipped_to_fail=true` 证门有牙；**未走 `NO CONCLUSION`＋exit 0 空转支**。序数（第二枚）＝〔仅转述〕，未独立回采全史。
5. **本批两验事项**：①"七步不再 skipped"＝**部分兑现**（链头从 :74 挪走、4/7 解封；gofmt 新成链头使 3/7 仍 skipped，job 级全 completed 但 step 级有残）；②"D32 第二枚真求值"＝**兑现为真结论**（本枚确为真采样、真判据）。

—— ci-attr-1，取数时刻 2026-10-05 03:12–03:24Z，run 37257547572（六 job id 见 §1）；本腿只 commit `verdict.md` 一件。
