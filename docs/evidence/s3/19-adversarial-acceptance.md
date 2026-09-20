# T19 对抗验收报告（T19-adv 独立执行）

> 执行者：T19-adv 对抗验收子代理（实现者为 agent-ticket19-taint，互不隶属；本报告作者未写任何被测代码）。
> 时间：2026-09-20T03:30Z 前后（本地 2026-09-20 上午）。
> 被测工件：`73b6057`（实现）+ `c33d10b`（文档/证据/票据转 review）。
> 环境：go1.27.1 windows/amd64，Windows 11 x64 真机；`export PATH=/d/work/base/go/bin:/e/work/base/msys64/mingw64/bin:$PATH; GOPATH=/d/work/base/gopath`。
> 并行保护：另有两代理在途（internal/agent/* 写票 10；票 02/04/07/17/HANDOVER 文档）。本验收全程
> **只读**它们的文件，构建/测试范围一律 `./internal/risk/`；未执行任何 `git add -A`；未 push。
> 方法论：除复跑既有测试外，另写 7 个**对抗探针测试**（`zz_probe*_test.go`，跑完即删，不入树），
> 逐个证伪实现方的断言。探针结论均可按 §6 的复现配方重跑。

## 1. 全量复跑（gates，实跑非转贴）

| 命令 | 结果 | 输出尾部（原样） |
|---|---|---|
| `go vet ./internal/risk/` | PASS | （无诊断）`VET=CLEAN` |
| `go test ./internal/risk/ -count=2` | PASS | `ok  	github.com/CarlosShao/wisp/internal/risk	3.615s` |
| `go test ./internal/risk/ -count=1 -v` 计数 | PASS | `--- PASS` 行 **60**（含 5 个 subtest），`--- FAIL` 行 **0** |
| `go test ./internal/risk/ -race -count=1` | PASS | `ok  	github.com/CarlosShao/wisp/internal/risk	1.525s` |
| `cd tools/d22scan && go run .` | PASS | `d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/` |
| `go build ./...`（仓根） | 预期失败，**非 T19 缺陷** | 失败点在 `internal/agent/`（票 10 未跟踪 WIP：budgets.go/compress.go/loop.go 等 13 个 `??` 文件）；`./internal/risk/` 单独构建/vet/测试全绿 |
| 提交面（契约纪律） | PASS | `git show 73b6057 --stat`：仅 `.scratch/wisp/issues/19-*.md` + `internal/risk/{provenance,taintmatch,syncdirs,syncdirs_windows,syncdirs_other}.go` + 4 个 `*_test.go`；`git show c33d10b --stat`：票 19 + `docs/PRECHECK.md`（纯追加，+30/-1 仅标题行）+ `docs/evidence/s3/19-*.txt` + `provenance.go`（+8，DEFERRED 头注释）。**`docs/PLAN.md`、`docs/specs/*.md`、`assessor.go`、`pathresolver.go`、他人票据零触碰** |

真机测试名（`-v` 实跑，票 19 相关 30 个测试函数：provenance_test 17 / syncdirs_test 7 / taintmatch_test 6）：

```
--- PASS: TestFourChannelExfilSuite (0.00s)      === RUN  5 个 subtest：web.search_query / notify_text / notify_url / clipboard.write / fs.write_into_sync_dir
--- PASS: TestR4EndToEndViaAssessor / TestSyncWriteNegative / TestDetectorUnknownToolScansEverything
--- PASS: TestScopesNeverInherit / TestDisposalScopeClearsTaints / TestInspectUnknownScopeIsEmptyStore
--- PASS: TestMarkInspectSourceAttribution / TestMarkNegativeNoTaintNoHit / TestMarkEmptyAfterNormalization
--- PASS: TestMarkUnknownSourceToolFailClosed / TestSensitiveSourceSetComplete
--- PASS: TestFragmentThresholdClampedToContractFloor / TestFragmentThresholdStricterAllowed
--- PASS: TestResidualLimitsAreLoggedNotSilent / TestConcurrentMarkInspect / TestLLMParaphraseResidual
--- PASS: TestSyncDefaultLocationProbe / TestSyncDropboxHostDBConfig / TestSyncFixtureFallbackAndMatch
--- PASS: TestSyncSuspectFallbackWhenUndetectable / TestSyncUnresolvablePathFailClosed
--- PASS: TestSyncRootsIntrospection / TestSyncDetectionOnThisMachine
--- PASS: TestNormalizeTaintContractForm / TestFragmentMatchExactBoundary / TestFragmentMatchAcrossTransforms
--- PASS: TestFragmentMatchCJK / TestFragmentIndexShortSourceNeverMatches / TestFragmentHashCollisionCannotFakeHit
```

> 证据文件 `docs/evidence/s3/19-provenance-c25-tests.txt` 与被测树**互证一致**：只出现
> `TestSyncDetectionOnThisMachine` 一个测试名，引用的 `syncdirs_test.go:133` / `:151` 行号与磁盘文件逐行对得上；
> 第 19 行自陈「HKCU Accounts 键只有 LastUpdate，无 UserFolder → 根由 default 探测确认」——诚实，但与 AC#4 的
> 「OneDrive 已配置机器上检出同步盘」合起来读会高估注册表探针（见 §2 行 5、§4 M-2）。

## 2. 实现方 8 条声明逐条证伪结果

| # | 声明 | 裁决 | 证据（本人读过的 file:line / 本人跑过的探针） |
|---|---|---|---|
| 1 | 5 条 AC 全绿；四条通道（自称六条） | **PASS（有保留）** | 21+ 测试确实存在且全绿（§1 复跑）；六通道：`provenance.go:115-119`（web.search/notify/clipboard.write/fs.write）+ `CheckText`（`provenance.go:422`）承接 TTS/HTTP body；`provenance_test.go:151-155` 断 TTS/HTTP 命中。保留：AC#4 的注册表探针**无任何测试**，TTS/HTTP 两条只断「命中」未断 Channel 标签 |
| 2 | 处处 fail-closed：不可解析路径→同步目标；探测不到→profile 下允许目录=sync-suspect | **FAIL（部分）** | 前半成立：`syncdirs.go:130-140`（err / 空 Canonical → `Sync:true`），`provenance.go:377,399`（无路径→按同步扫）；探针 `IsSyncPath("")` 判 sync。**但**：(a) `syncdirs.go:86` `complete = len(s.roots) > 0` 使兜底**全有全无**——任一探针命中即整体关掉（§4 M-3，已实测）；(b) `syncdirs.go:130` 只看 `err`，把 C26 的 `res.Resolved==false`（纯词汇拼写、未核实）当已核实结论 → 8.3 / `\?\` 拼写的**新建写入目标**判「非同步」（§4 B-1，已实测写盘成功）；(c) `provenance.go:355-360` 未打开的 scope 直接「视为无污染」返回不命中——接线拼错即静默放行（§4 M-1） |
| 3 | 成员判定只走票 18 的 `risk.Resolve`，无手写路径比较 | **PASS（有 MINOR 反例）** | 同步判定唯一入口 `Resolve`：`syncdirs.go:95,130`；根侧 `syncdirs.go:102`；比较用票 18 的组件感知 `isUnder`（`blacklist.go:102`，`syncdirs.go:145`）。反例：suspect 兜底用裸 `strings.HasPrefix(cand, s.home)`（`syncdirs.go:149`）——非组件边界，实测 `…\profile` 把兄弟目录 `…\profileevil` 一并判 suspect（方向偏严，MINOR-5） |
| 4 | 归一化 + ≥8 rune 连续窗口；哈希命中必做 substring 复核，碰撞伪造不了命中 | **PASS（生产代码）/ FAIL（测试）** | 归一化三项与 SPEC-06 §5 字面一致：`taintmatch.go:34-50`；索引/二分/复核：`taintmatch.go:94-113,144-149`（`strings.Contains(src, w)` 定生死，碰撞只花时间为真）。**但**：① `TestFragmentHashCollisionCannotFakeHit`（`taintmatch_test.go:87-95`）从未制造碰撞——`yyyyyyyyzzzzzzzz` 与索引零交集，实为「不相交窗口不命中」的重言式；② 探针证明 `taintmatch.go:152-161` 的「碰撞后探邻居」分支不可达（索引按哈希去重，`dups=0`）＝死代码；③ `TestFragmentIndexShortSourceNeverMatches:79` 喂空索引，恒假 |
| 5 | `TestSyncDetectionOnThisMachine` 真机检出 OneDrive | **OVERCLAIM** | 检出者为 `defaultLocationProbe` 的 `os.Stat($HOME\OneDrive)`（`syncdirs.go:196`）——**P12 明文说「不能只靠路径字符串匹配」的那一类**；注册表探针（`syncdirs_windows.go:40-70`，P12 的正解）在本机返回 0 条（证据文件第 19 行自陈），且该测试对探针结果**只 t.Logf 不断言**（`syncdirs_test.go:131-152`，条件式断言全被机器状态短路）。真机探针复跑：`complete=true roots=[{OneDrive C:\Users\swq\OneDrive default}]` |
| 6 | 只留 `DEFERRED(C25-loop-wiring)`，未碰 `internal/agent/` | **PASS** | `provenance.go:46-52`；两次提交文件清单无 `internal/agent/`（§1）；全仓 `NewProvenance(` 仅出现在 `internal/risk/*_test.go`（grep 证实无生产调用点）。保留：marker 点名的票 21/26 **票据正文无一行**承接污染闸门接线（§4 MINOR-6） |
| 7 | 请裁定：`Detector(scopeID)` 适配器是否偷偷弯折冻结契约 | **裁定：未弯折，但有真实代价** | `assessor.go:162-167` 的 `TaintHit(params)` 签名与 `rules_gateway.go:99-114` 的 R4 消费方**均未被两次提交触碰**；适配器 `provenance.go:455-466` 忠实实现接口。经票 19 探针实跑端到端：`reason="R4: 包含来自 web.fetch https://x 的内容" blocked=true`。代价：① 通道标签在缝上丢失（适配器一律 `ChUnknown`，`provenance.go:406`）；② **两条集成路径对同一次调用答案不同**——适配器（工具名 ""→泛扫）比 `Inspect(带工具名)` **更严**，于是「按文档推荐用 Inspect」反而更松（§4 B-2）；③ scope 在构造期烧死、无法重绑，配错即 `provenance.go:355` 静默不命中（§4 M-1） |
| 8 | LLM 改写残留＝文档认可的接受缺口 | **成立：SPEC 背书，非自我豁免** | 见 §5 |

## 3. 五条验收标准逐条

| AC | 裁决 | 判定依据 |
|---|---|---|
| ① 四通道 exfil 套件，逐通道升 L2 且指明来源 | **PASS** | `provenance_test.go:117-157`（五 subtest 逐条断 Channel + `hit.Source()` 含 `search.content`）+ `:161-186` 经**真实 C19** 断 L2/[R4]/SessionOverrideBlocked/「包含来自 doc.read」；反向用例 `:191-208` 钉住「非同步目录写入不算 R4 通道」（与 SPEC-06 §5 原文一致，PRECHECK 第 145-147 行亦已写明）。保留：通道标签只在 Inspect 路径断，缝上丢标签（行 7②） |
| ② 归一化测试：空白/大小写/全角的变换仍命中；<8 不命中 | **PASS（生产）/ 弱（测试）** | 边界真实：`taintmatch_test.go:28-44`（8 命中 / 7 不命中）、`:63-73` CJK、`:46-61` 全角；阈值钳制 `provenance.go:205-211` 有测试（`provenance_test.go:270-296`）。缺陷：**所有 Mark→Inspect 用例喂的是同一个 `marker` 常量**，集成层从未真正走过归一化（本人探针 `P4` 补测证明代码是对的，因此只记测试缺口 MINOR-4） |
| ③ LLM 改写容差测试（记录为限制） | **PASS** | `provenance_test.go:347-360`：断「改写后漏检」+ 断「逐字连续片段仍命中」（后者防空实现把匹配整块写坏）。定性见 §5 |
| ④ 已配置 OneDrive 的机器上检出同步盘（自证）+ fixture 兜底；注册表探针过 P12 评审 | **FAIL** | fixture/兜底/去重都真绿（`syncdirs_test.go:40-122`）；但①注册表探针**零测试 + 本机零贡献**（`syncdirs_windows.go:40-70` 未被任何断言覆盖；证据文件第 19 行）②「OneDriveConsumer/OneDrive 环境变量」这条本机唯一可靠的已配置根线索**代码里根本没读**（`syncdirs.go:184-201` 只 os.Stat 猜默认目录，本人探针实测 `$OneDrive=C:\Users\swq\OneDrive` 在位）③真机测试对探针结果不断言。④更严重：8.3 / `\?\` 拼写的新建写入目标绕过同步判定（§4 B-1） |
| ⑤ 污染作用域：新会话不继承旧污染（DisposalScope） | **PASS（附 1 个 MAJOR）** | `provenance_test.go:224-266`：真 `plugin.DisposalScope` Defer→Dispose 生命周期 + 第二 scope 隔离 + 新 scope 空店。MAJOR 附款：未打开/已关闭 scope 上的 Inspect 一律「无命中」（`provenance.go:355-360`），且该行为被 `TestInspectUnknownScopeIsEmptyStore` 当成正确钉死——默认方向应是拒绝而非放行（§4 M-1） |

## 4. 缺陷清单（B=BLOCKER / M=MAJOR / N=MINOR）

**B-1 同步根成员判定被「未存在路径的词汇回退」绕过（实测可写盘成功）**
`syncdirs.go:130-141` 只把 `Resolve` 的 error 当不可信，忽略 `res.Resolved==false`（`pathresolver.go:86` 对不存在路径走纯词汇回退，`\?\` 前缀与 8.3 短名**都不会**被展开——`stripExtendedPrefix`/句柄展开只在 `pathresolver.go:74-81` 成功分支里）。fs.write 的目标天生「尚不存在」，正中该分支。探针实测（同步根内、文件新建）：
`spell=extlen sync=false r4hit=false` → `EXPLOITABLE: wrote 40 bytes INTO ...\OneDrive\Notes\out2.md`；
`short83=ONEDRI~1 … landed-in-sync-root=true … EXPLOITABLE via 8.3`（junction 与 `\localhost\c$` UNC 方向正确，均 fail-closed 判 sync）。
**一行修法**：`match()` 里 `if err != nil || !res.Resolved || res.Reparse { …sync-suspect 从严… }`，或对候选再做一次 `stripExtendedPrefix`+`GetShortPathNameW` 规范化；并补 8.3/`\?\` 两条真写盘红队用例。

**B-2 具名通道的 `Inspect` 比「未知工具泛扫」更窄＝参数名/类型逃逸**
`provenance.go:389-393` 只 `params[k].(string)` 取表内 key；未知工具分支（`:396-409`）反而扫全部嵌套字符串。探针实测四处逃逸且**同一个调用经冻结适配器就能抓到**：
`web.search{url}`、`web.search{q}`、`notify{title}`（text 无害）、`clipboard.write{payload}`、`notify{text:{body:…}}`、`notify{text:[…]}` 全部 `Inspect` 不命中；`adapter DOES catch it: "doc.read /contracts/acme.pdf"`。fs.write 的 `{"body": …}` 同理（`writeChannelKeys` 只认 content/data，`provenance.go:124`）。D34 参数名权威表（`docs/TOOLS.md`）尚未成文，表内 key 现在是**猜测**。
**修法**：具名通道也走 `collectStrings`，表只用于**贴 Channel 标签**，其余参数以 `ChUnknown` 追加扫描（写载荷仍受同步门约束），方向与 16.9#1「误报安全」一致。

**M-1 未打开/拼错 scope 上的 Inspect 默认放行（缺信息→通过，而非→拒绝）**
`provenance.go:355-361`；`Mark` 落在 scope A、`Inspect` 查 scope B（typo 或 Close 后仍有调用）即恒不命中，只留一行 log。仓规要求「信息缺失必须产生 DENY」。**修法**：引擎内存在任何 taint 而 scopeID 未知时，返回 `Hit{Channel: ChUnknown, SrcTool:"unbound-scope"}` 从严升 L2；并把 `TestInspectUnknownScopeIsEmptyStore` 改成断这个方向。

**M-2 注册表探针（P12 的正解）无任何覆盖，且本机唯一在位的已配置根线索 `%OneDrive%` 未被读取**
`syncdirs_windows.go:25-70` 零测试；`syncdirs.go:184-201` 只 Stat 默认目录。真机事实：`$OneDrive`/`$OneDriveConsumer` 均指向 `C:\Users\swq\OneDrive`，而 `HKCU\...\Accounts` 无 `UserFolder`。**修法**：探针链首加 `envOr("OneDrive","")`/`OneDriveConsumer`/`OneDrivePublic`，并为 `registryProbe`/`oneDriveRoots` 加「可注入 registry reader」的表驱动测试（真注册表用例带 `t.Skip` 亦可）。

**M-3 「探测不到→profile 下目录一律 suspect」是全有全无开关，可被单个诱饵根关掉**
`syncdirs.go:86`（与 `provenance.go:255`）。探针实测：注入一个无关根 → `complete=true` → `…\profile\Documents\exfil.md` 判 `sync=false`（「FALLBACK DISABLED by a single unrelated root」）；本真机同理（default 探针命中即关兜底，`Documents`/`Desktop`/`Downloads` 全判非同步）。即「装不到 OneDrive」这条路**关不掉**兜底（不装只会更严，这一向攻击不成立），但**创建一个空 `~\OneDrive` 目录 / 拥有任一被识别客户端**即可关掉全网状兜底，而 P12 对识别不准的后果写的是「这条外泄通道就是开的」。**修法**：兜底开关改按「是否拿到 **registry/config 级**确认根」（`Source=="registry"||"config"`）判定，仅 `default`/`fixture`/`options` 级证据不撤销兜底；并把 `syncRootC.canonical==false`（`syncdirs.go:59`，现字段**设了没人读**）纳入「该根未证实」计数。

**M-4 `collectStrings` 深度上限静默漏检**
`provenance.go:489-492`（`depth>4` 直接 `nil`，无日志）。探针：6 层嵌套载荷 `P7 depth-5 nesting escaped the generic scan`。工具参数是模型可控 JSON，套壳成本为零。**修法**：命中深度上限时按「扫描被截断」走从严（或整体转 `ChUnknown` 命中 + 日志），与 `MaxScanRunes` 的「记日志」口径一致。

**M-5 片段索引内存/时延与 D32 预算冲突，且主因是一块死字段**
`taintmatch.go:71,112` 的 `winStrs` 写后**从不读**（复核走 `strings.Contains(src, w)`，`:147`）——本人把每个表项篡改成 `"ZZZZZZZZ"` 后匹配照样成功，实证其为死重。实测：单个满额（262144 rune）源 `Mark` 分配 12 MiB、**常驻 8 MiB**，5 个即 44 MiB；`Mark` 一次 52ms；20 个 8k 源之上扫 40k 参数 57ms（响应路径上）。**修法**：删 `winStrs`（顺带删 `:152-161` 死分支），只存排序 `[]uint64`（≈2 MiB/源）；给 `MaxSourceRunes`/每 scope 源数设显式预算并在 PRECHECK/SLO 记数。

**M-6 归一化不处理「零宽/软连字符」等 default-ignorable 字符 → 片段闸门可被隐形字符穿过**
`taintmatch.go:34-50` 只 `IsSpace` 剥离；`U+200B`/`U+00AD` 不是空白，切断连续性。探针实测：把 marker 里的 `-` 换成 `U+200B` 后 `Inspect` **不命中**（`P3 EVADE`），`U+00AD` 同。**这不是** AC#3 被 SPEC 认可的「改写/翻译」残留（那是语义级、16.9#1 明确驳回精确 taint）；这是同一份规范化文本层面的实现缺口，成本 3 行。**修法**：`normalizeTaint` 丢弃 `unicode.Is(unicode.Mn,r)`、`U+00AD`、`U+200B..200D`、`U+2060`、`U+FEFF`。

**N-1 `ProvOptions.AllowedDirs` 是死字段**（`provenance.go:144` 全仓无消费；`syncdirs_test.go:85-101` 传了它却靠 `home` 前缀过测），而 PRECHECK 第 123-124 行与票 19 约束都把「`[fs] allowed_dirs` 位于 profile 下」当规则本体。**修法**：要么真按 allowed_dirs 收口，要么删字段并把文档改成「profile 下任意路径（严于 allowed_dirs）」。
**N-2 `syncRootC.canonical` 设而不读**（`syncdirs.go:59,102`）→「根自身未证实」信息丢失（并入 M-3 修法）。
**N-3 `TestFragmentHashCollisionCannotFakeHit` 名不副实**（不制造碰撞）＋ `taintmatch.go:150-161` 不可达。修法：删死分支，测试改名或按 `hashWindow` 反查构造真碰撞（或注入可替换哈希函数）。
**N-4 集成层无「变形后仍命中」用例**：`provenance_test.go` 全部喂同一 `marker`；本人探针 `P4`（大小写+全角+空白经 Mark/Inspect）通过，说明只是缺测。修法：把 `P4` 那条例子并入 `TestFourChannelExfilSuite`。
**N-5 `syncdirs.go:149` 裸 `HasPrefix` 越界匹配兄弟目录**（`profileevil` 被牵连，方向偏严）。修法：换 `isUnder`。
**N-6 DEFERRED 登记不完整**：`provenance.go:46-52` 点名票 21/26，但两票正文无 C25/taint/R4 任何 AC 行；`syncdirs_other.go:10` `DEFERRED(P12-macos)` 点名片 55，票 55 正文无同步盘/CloudStorage 字样。修法：各补一行 AC（票 26 TTS 播报必须过 `CheckText(ChTTS)`；票 55 macOS 云盘探测判据）。
**N-7 进度日志「21 new tests」实为 30 个测试函数**（少计，全绿，无风险）。

## 5. 两项专门裁定

**(a) LLM 改写残留＝SPEC 背书，不是自我豁免。** SPEC-06 §5 第 78-83 行把判定限定为「≥8 字符**连续片段**（去空白/大小写折叠/全半角统一）」，并**明文解释**为何不精确 taint（「LLM 会改写内容……漏报由 D30 其余层兜（16.9#1 驳回记录）」）；SPEC-06 §12 第 158 行「不做精确 taint tracking（16.9#1 驳回留痕）」；PLAN 第 1537 行 RESERVED 表的后果列直接写「LLM 把敏感内容**改写/翻译**后再外发时可能漏检（由 D30 其余四层兜）」；16.9#1（PLAN 第 3137 行）并禁止把它当「待补完整版」复活。`TestLLMParaphraseResidual` 的注释引用的层号也对得上 SPEC-06 §6（②黑名单/③出站约束/⑤首域名）。→ **接受，无须整改**；AC#3 判 PASS。
**但**：M-6（隐形字符）不属于该背书范围——契约是「规范化之后」的片段级判定，把 `U+200B` 洗掉正是「规范化」应有之义，不得拿 16.9#1 当免检牌。

**(b) 「同步探测 fail-closed」声明：不成立（部分）。** 站得住的三条：空路径/Resolve 报错→sync（`syncdirs.go:130-140`，实测）；未豁免 junction→`ErrReparseDenied`→sync（实测）；「不装 OneDrive」这条路只会让兜底更严（实测 `complete=false` 分支）。站不住的两条：其一，兜底是**全有全无**（`syncdirs.go:86`），一个 default 级/fixture 级甚至诱饵命中即可整体撤销，真机现状即如此（M-3）；其二，**同一根内、拼写不同、文件尚不存在**的写入目标被判「非同步」且实际写盘成功（B-1），这直接推翻 `syncdirs.go:30-33` 与 PRECHECK 第 133-134 行「junction/8.3/UNC 拼写塌缩到同一根」「无法经 C26 解析→从严」的字面承诺。→ 声明 2 判 **FAIL（部分）**。

## 6. 探针复现配方（7 个临时测试，跑完已删，未入树）

`NewProvenance(NoProbe/sync-root)` → `Mark` → `Inspect/CheckText/Detector`，含：P1 通道表窄扫逃逸、P2 嵌套/非字符串载荷、P3 `U+200B`/`U+00AD`、P4 端到端归一化（PASS）、P5 同步命中归因（正例确由根匹配驱动，非 fail-closed 蒙对）、P6 兜底全有全无 + 兄弟前缀、P7 退化输入/深度、P8 短源、P9 经 C19 的适配器（`R4: 包含来自 … 的内容`、`blocked=true`）、P10/P12/P13 拼写矩阵与**真写盘落地核验**（`os.WriteFile(`\\?\…`)` / `dir /x` 取真 8.3 别名）、P14/P16/P17/P18 内存与时延、P15 `winStrs` 死字段篡改证明。
重跑：把上述任一形态写成 `internal/risk/zz_probe_test.go` 后 `go test ./internal/risk/ -run TestProbe -v`。

## 7. 仓规与硬禁令扫描

| 禁令 | 结论 |
|---|---|
| 裸 `go func(`（须 owner+recover） | PASS：新增生产码零 goroutine；唯一 `go func(` 在 `provenance_test.go:325`（测试文件，d22scan 设计上跳过 `_test.go`；带 `wg.Done()`，无 recover——测试内可接受） |
| `filepath.Clean/Abs` 越权（唯一豁免 `internal/risk/pathresolver.go`） | PASS：新增码只用 `filepath.Join`；`tools/d22scan` 全仓 clean（allowlist 未新增条目） |
| 墙钟差判超时 | PASS：`observe.NowWallUTC()` 仅用于 `taintMark.at` 取证时间戳（`provenance.go:311`），无任何差值判时 |
| UI 面 emoji / 硬编码阈值 / 明文凭据 | PASS：新文件无 emoji（perl `\p{Emoji}` 扫描 0 命中）；8 为契约地板且只能调严（`provenance.go:205-211` + 测试）；`MaxSourceRunes/MaxScanRunes` 可配、默认值有日志；无任何凭据落盘。`collectStrings` 深度 4 是硬编码（见 M-4） |
| 冻结面 | PASS：`assessor.go`/`pathresolver.go`/`docs/PLAN.md`/`docs/specs/*` 均不在这两次提交里 |
| 污染碎片不得落库 | PASS：`Hit.Fragment` 在 `internal/risk/` 之外零消费者（grep `\.Fragment` 无外部命中），`provenance.go:189-190,51-52` 明示只进原生卡 |

## 8. VERDICT

**FAIL — 退回修复。** 工程质量与文档纪律整体扎实（30 个测试全绿、race/vet/d22scan 干净、契约零越界、DEFERRED 基本诚实、残留登记有据），但 R4 污染闸门本体存在**两条已被我实测复现的 fail-open 逃逸**（B-1 同步根拼写绕过并真把字节写进同步根；B-2 具名通道参数名/类型逃逸），另有一条与「缺信息必须拒绝」硬规冲突的默认放行（M-1）和一条被 SPEC 背书范围**不覆盖**的规范化缺口（M-6）。AC 判定：①③⑤ PASS，② PASS（附缺测），④ FAIL。

必改（回归前）：**B-1、B-2、M-1、M-3、M-6**；建议同批：M-2（注册表探针至少补可注入测试）、M-4、M-5（删 `winStrs` 顺带解决大半内存）、N-1..N-7。AC#4 若要判过，须给出**注册表或环境变量级**根确认（不接受只靠 `os.Stat` 猜默认目录），并补 8.3/`\?\` 两条真写盘红队用例。

---

## 复验（2026-09-20，第二轮）

> 复验者：独立代理。**未写被测码，未写上面那份 FAIL 判定书**（判定书原文一字未改，历史即历史）。
> 被测工件：`c091ee2`（码+测试）+ `6e2a68f`（文档/证据），基线 `8d9fd95`。
> 环境：go1.27.1 windows/amd64，Windows 11 x64 真机（i7-8750H）；
> `PATH=/d/work/base/go/bin:/e/work/base/msys64/mingw64/bin`，`GOPATH=/d/work/base/gopath`。
> 并行保护：复验期间另有代理在 `internal/agent/`/`internal/llm/`/`tools/mockllm/` 在途并落了 `ddb79e2`
> （票 10 证据，与本票无关，未触碰）；本轮所有构建/测试一律 `./internal/risk/` 范围，未 `git add -A`，未 push。
> 方法：**不复用首轮 7 个探针**，本轮另写 5 个临时探针文件 `internal/risk/zz_reverify19{,b,c,d,e}_test.go`
> （共 23 个探针函数：真 junction/真 hardlink/真写盘落地核验/真注册表与环境变量），
> **跑完即删，未入树**（`git status` 已确认无 `zz_*` 残留，`C:\Users\swq\OneDrive` 与 profile 内零探针残留）。
> 本轮首要攻击对象是**修复自带的新攻击面**：B-1 用「最近已存在祖先 + 词汇尾巴」修好之后，
> reparse 点若长在**尾巴**里而不是祖先里，判定还看不看得见。

### R1 五道闸门（本人实跑，探针删净后复跑一次，非转贴）

| 命令 | 结果 | 输出尾部（原样） |
|---|---|---|
| `go vet ./internal/risk/` | PASS | （无诊断）本人 echo：`VET=CLEAN` |
| `go test ./internal/risk/ -count=2` | PASS | `ok  	github.com/CarlosShao/wisp/internal/risk	4.853s` |
| `go test ./internal/risk/ -race -count=1` | PASS | `ok  	github.com/CarlosShao/wisp/internal/risk	1.623s` |
| `cd tools/d22scan && go run .` | PASS | `d22scan: clean - no D22 ban violations, no emoji in design/ or frontend/` |
| `go test ./internal/risk/ -count=1 -v` | PASS | `top-level PASS: 69  SKIP: 1  FAIL: 0  subtest PASS: 37` / `ok  	github.com/CarlosShao/wisp/internal/risk	1.785s` |
| `GOOS=darwin GOARCH=arm64 go vet ./internal/risk/` + `GOOS=linux go build ./internal/risk/` | PASS | （无诊断）新码可交叉编译（见 R5 N-9：能编不等于能在非 Windows 上按注释所述工作） |

**实跑计数 vs 修复方声明**：日志（N-7 更正条目）称「68 个顶层测试 + 38 个 subtest，1 skip」。
实测是 **70 个顶层测试函数**（69 PASS + 1 SKIP）+ **37 个 subtest**（`=== RUN` 107 行 − 70 顶层 = 37，
无 indented SKIP/FAIL）。仍是同方向的少计（低报 2 个顶层、多报 1 个 subtest），不影响判定，但 N-7 的「更正」本身没更正准（见 R5 N-7 行）。

**冻结面复验（零 diff，实跑非转贴）**：
`git diff --stat 8d9fd95..6e2a68f -- internal/risk/assessor.go internal/risk/pathresolver{,_windows,_other,_budget_norace_test,_junction_windows_test}.go internal/risk/rules_gateway.go docs/PLAN.md docs/specs/`
→ **空输出**。两次提交合计只碰 14 个文件：`internal/risk/` 的 5 个生产文件 + 4 个测试文件 + 票 19/26/55 + `docs/PRECHECK.md` + `docs/evidence/s3/19-adversarial-fix-round.txt`。
`git show --name-only c091ee2 6e2a68f | grep -c internal/agent` → **0**。C26 入口 `Resolve`/`reparseComponents` 与票 18 的红队测试都在冻结面内、未被改动——这一点重要：B-1 的修法是在**消费侧**加锚定，没有去松 resolver。

### R2 头条裁定：**尾巴里的 reparse 点——B-1 没有换帽子回来（判：不是 BLOCKER）**

构造：沙箱 profile 内 `notsync\` 为普通目录，`notsync\sub` 用 `mklink /J` 做成指向同步根的 junction，
写入目标 `notsync\sub\file.md`（**尚不存在**，正是 B-1 的靶形）。同步根以 registry 级注入且目录真实存在，
故 `SyncDetectionComplete()==true` —— **profile 全网状兜底是关着的**，任何「非同步」判定只能来自根成员比较，蒙不出来。

```
Resolve(raw): err=risk: path traverses a reparse point ... resolved=false reparse=false canonical=
reparseComponents(raw) = [...\profile\notsync\sub]
deepestExistingAncestor(c:\...\profile\notsync\sub\file.md) -> anc=c:\...\profile\notsync\sub rest=[file.md]
VERDICT sync=true root=sync-suspect/suspect-fallback why=write target traverses a non-exempted reparse point; fail-closed as sync-suspect
VERDICT deep tail(notsync\sub\a\b\file.md) sync=true why=write target traverses a non-exempted reparse point; fail-closed as sync-suspect
```

真机版（指向操作者真实 `C:\Users\swq\OneDrive`，只做判定不写盘，探针目录跑完已删）：

```
live root OneDrive/env C:\Users\swq\OneDrive   complete=true
REAL ROOT tail-junction verdict sync=true why="write target traverses a non-exempted reparse point; fail-closed as sync-suspect"
```

**为什么看不见尾巴的担心不成立（两道独立防线，均实测）**：
1. `resolveTarget` 第一句 `Resolve(raw, excs)` 里的 `reparseComponents` 走的是**清洗后的整条词汇链**，
   逐个组件 `GetFileAttributes`，尾巴里的 junction 同样被检出 → `ErrReparseDenied` → `match()` 落
   「穿越未豁免 reparse 点」从严分支。祖先接尾巴那步**根本没被执行**（实测 `err` 非空即返回）。
2. 即便走到祖先步，实测 `deepestExistingAncestor` 返回的正是 junction 本身（`anc=...\notsync\sub`）而非其父目录：
   它 Lstat 时给每级补了尾分隔符，而带尾 `\` 的 Lstat 会**跟着 mount point 进去**（Go 对无斜杠的 junction
   Lstat 报 `isdir=false symlink=false`+reparse 属性位——代码注释因此不准，见 R5 N-9）。于是
   `Resolve(anc)` 要么照旧 `ErrReparseDenied`，要么（用户显式豁免该路径时）**穿透 junction 规范化到真根**：

```
Resolve(exempted): err=<nil> resolved=false reparse=true canonical=...\notsync\sub\file.md
deepestExistingAncestor -> anc=...\notsync\sub rest=[file.md]
resolveHandle(junction) = \\?\C:\...\profile\OneDriveRV
VERDICT(exempted junction into sync root) sync=true root=OneDrive/registry why=write target is under a detected sync root
VERDICT(no exception)                    sync=true why=write target traverses a non-exempted reparse point
```
豁免分支拿到的是**根成员归因**（不是 suspect 蒙对），方向比首轮要求的还严。

**尾巴含 `..`（出圈再回圈）**：首轮没问过的形状，我按最坏情况构造了两次。判定确实会**看不见** junction——
`Resolve` 里 `lexCanonical`（票 18 冻结的 Clean）先把 `link\..` 折掉，`reparseComponents` 便扫不到被擦掉的 `link`，
`VERDICT sync=false root=/ why=write target is not under any sync root`。但**落点核验否定了逃逸**：
我把同名文件在两个候选落点都预置好，再走 OS 真写盘，让 `GetFinalPathNameByHandle` 说真话：

```
raw = ...\profile\notsync\link\..\X.txt ; filepath.Clean = ...\profile\notsync\X.txt
  ...\profile\OneDriveRV\X.txt            -> "ORIGINAL"
  ...\profile\OneDriveRV\Notes\X.txt      -> "ORIGINAL"   <- 若 OS 跟着 junction 折 ..，该改的是这里
  ...\profile\notsync\X.txt               -> "MODIFIED-BY-PROBE"
OS final path for raw = \\?\...\profile\notsync\X.txt
two-level: ...\notsync\link\..\..\Y.txt  OS-final=\\?\...\profile\Y.txt  (inside-root=false)
hop\..\Notes\v3.md -> write failed: The system cannot find the path specified.（OS 找的是 notsync\Notes\v3.md）
```
即 **Windows 自己在解析前就把 `.`/`..` 做词汇归一**，落点与 `filepath.Clean` 逐字一致，分类器与真相同源，
没有「判非同步却落进同步根」的裂口；`\\?\`+`..` 那一路 OS 直接 `ERROR_INVALID_NAME` 拒写（哪都不落）。
⇒ **不构成 BLOCKER，B-1 保持关闭**。但这层安全性来自 OS 语义而不是代码里的一句检查，故记一条硬化残留（R5 N-10）。

**尾巴 reparse 的 TOCTOU（判定之后才创建 symlink/junction）**：实测成立但判**超出本票威胁模型**，理由写实：

```
verdict before the junction exists: sync=false root=/
（随后 mklink /J notsync\later -> root）bytes landed in root: true
verdict re-asked after the junction exists: sync=true why=write target traverses a non-exempted reparse point
```
要吃下它必须在 `Inspect` 与 `CreateFile` 之间赢得一个竞态，即攻击者**已是本机进程**；而能建 junction 的本机进程
本就能直接往 `OneDrive\` 写字节，闸门于它无意义。C25 的门禁对象是**模型自选的拼写**（B-1 的原文形状），
而模型分两步（一次 shell `mklink /J`、下一次 `fs.write`）在这套实现里是**关得住的**：写盘前那一刻重新判定即翻成 deny
（上面第三行）。⇒ 记为接受残留，无须整改；建议 PRECHECK 补一行说明「判定与写盘非同一原子操作」的定性。

### R3 逐缺陷处置（closed / not closed / newly broken）

| 缺陷 | 首轮要求 | 修复方声明 | **复验裁定（判据＝本人实跑，非读注释）** |
|---|---|---|---|
| **B-1** | 未存在写入目标的词汇回退绕过同步根，实测写盘成功 | 锚定最近已存在祖先 + 接词汇尾巴 | **CLOSED**。沙箱内 4 种拼写（长名/`\\?\`/OS 8.3/两者叠加）由 `GetShortPathNameW` 取回真别名并门控写盘验 0 字节；我自己另测尾巴 junction、深尾巴、豁免 junction、不可锚定链（`Q:\`、死 UNC）四种形状，全部从严；反向护栏成立（普通新文件写长名/`\\?\`/8.3 → sync=false 且字节真落地，实测 `TestSyncRedTeamGuardPlainNewFileWriteAllowed` 绿）。修法**没有**采用判定书那句 `!res.Resolved 即从严`（那会把每次正常新建写入推去 L2），并给出了理由——这条偏离我认可，因为红队用例证明它同时守住了正例与误报两侧 |
| **B-2** | 具名通道只扫表内 key，改名/嵌套/数组即逃逸 | 表只贴标签，其余一律 `collectStrings` 以 ChUnknown 补扫 | **CLOSED（六条逃逸全实测被抓）**：`web.search{url}`/`{q}`、`notify{title}`(text 无害)、`clipboard.write{payload}`、`notify{text:{body}}`、`notify{text:[…]}`、`fs.write{"body"}` 均命中，且**归因未回归**——表内 key 保留契约标签（query→web.search.query、url/text→notify、text→clipboard.write、content→fs.write.syncdir），其余 `ChUnknown`，`hit.Source()` 仍指向 web.fetch；同形状走冻结适配器 `Detector.TaintHit` 亦抓；`[]byte`/`[]string`/深 map 都抓。写载荷仍受同步门约束：`{path:非同步, content}` 不命中、`{path:同步根, content}` 命中 ChSyncWrite、`{path:D:\plain}` 不命中（实测四条） |
| **M-1** | 未打开/拼错 scope 默认放行，违反「缺信息必 DENY」 | scope 未注册且引擎有污染 → `SrcUnboundScope` 从严 | **CLOSED**。`TestInspectUnknownScopeIsEmptyStore` 被**原地反转**而非删除（名字保留，断言换向）；我自己独立复跑：`Inspect("ghost")` → `SrcUnboundScope`+`ChUnknown`；`CheckText("ghost",ChTTS)` 命中；`Detector("ghost").TaintHit` 命中（三条缝全覆盖）。方向也没有过冲：引擎零污染时 `Inspect("ghost")` 不命中；已注册且为空（OpenScope 后未 Mark / CloseScope 后）不命中——AC#5 的「新会话不继承」仍由 `TestDisposalScopeClearsTaints` 钉住 |
| **M-2** | 注册表探针零覆盖，且 `%OneDrive%` 这条本机唯一线索没被读 | 探针链首加 env 级；注册表逻辑下沉可注入 `registryValueSource` | **CLOSED**。本机实测 `%OneDrive%`/`%OneDriveConsumer%` 均为 `C:\Users\swq\OneDrive`，`SyncRoots()` 给出 `OneDrive/env`、`complete=true`、真机红队归因 root=OneDrive/env——AC#4 要求的「注册表或环境变量级根确认」达成（首轮判 FAIL 的正因消解）。`registryProbeFor` 由 4 个 subtest 的表驱动测试真驱动（Personal/Business 过滤、`ConsumingAccounts` 诱饵、空 UserFolder、大小写去重、Dropbox `Path`、nil 源）。**唯一 t.Skip 诚实**：`TestSyncRegistryProbeLive` 本机返回 0 条（我独立复算同果），skip 理由与 `19-provenance-c25-tests.txt:19` 自陈一致；探针**逻辑**由注入测试覆盖，探针**存在性**由该 skip 测试记录，这个分工是可接受的 |
| **M-3** | 兜底全有全无，一个诱饵根即可关掉 | 仅 env/registry/config 级且 C26-canonical 才撤销 | **CLOSED**。我重跑诱饵攻击：注入单个**真实存在**的无关根，source 取 `default`/`fixture`/`options`/`""` 四档 → `complete=false` 且 `~\profile\Documents\exfil.md` 仍 sync-suspect；`registry`/`config`/`env` 三档 → `complete=true` 且兜底关掉（设计意图）。`canonical` 从此**真被读**：`{registry, 不存在的路径}` → `complete=false`、兜底仍开（N-2 一并闭环）。suspect 网另经组件边界收口（`profileevil` 兄弟目录不受牵连，实测） |
| **M-4** | depth>4 静默 nil | 预算可配（默认 8），超深 → `SrcUnscannedNesting` 从严 | **CLOSED**（12 层嵌套实测 `SrcUnscannedNesting`；3 层内仍是真内容命中 `fs.read`；`MaxParamDepth:0` 回落默认而非「不扫」；空嵌套不误触发截断） |
| **M-5** | 死字段 `winStrs` + 内存/时延撞 D32 | 删死字段/死分支/每候选全量拷贝；**明说不宣称时延改进** | **CLOSED，且声明是诚实的**。我实跑：`BenchmarkMarkFullSource 29309800 ns/op 4319069 B/op 12 allocs/op`、`BenchmarkScanOverTenSources 30019433 ns/op 1574669 B/op 17 allocs/op` —— 与 PRECHECK 记的 27.8ms/4.3MiB、每源 ~2 MiB 哈希集自洽（262144 rune × 8 B/窗＝2.0 MiB，正是删掉 `[]string` 平行表后应有的量级）；`winStrs` 与「碰撞探邻居」分支在磁盘上确已不存在。**未宣称的改进也没有被偷偷宣称**：文档原文写「55.8~58.9ms vs old 57.4~57.8ms，差异落在噪声内，未记为改进」，与我跑出的同量级一致。残留（不影响裁定）：`contains` 每个 mark 重复做一次 `[]rune(paramNorm)`，N 个源即 N 次转换——这解释了「无时延收益」，也留了下一步的空间 |
| **M-6** | 零宽/软连字符切断 ≥8 窗口 | 丢弃 U+00AD/200B/200C/200D/2060/FEFF + `Mn` | **CLOSED**。六个码点 + combining mark 逐一实测（候选侧与**源侧**双向）不再切断 ≥8 窗口；`TestNormalizeTaintDropsIgnorableChars` 用码点构造而非字面隐形字符（编辑器往返安全），并另在归一化空间直打 `contains`。**未过度归一**：无关串、被隐形字符注水的无关串、`MARKER`（7 rune）、`z×200`、`12345678` 全不命中；纯隐形字符源被 `Mark` 拒（返 false）。匹配器没有变成「什么都命中」 |
| N-1 | `AllowedDirs` 死字段 | 删字段，规则定为「profile 下任意路径，严于 allowed_dirs」 | **CLOSED**（`internal/risk/` 内 `AllowedDirs` 0 命中；其余命中都是 `internal/config` 自己的 `fs.allowed_dirs`，与本票无关；PRECHECK 同步改口径） |
| N-2 | `canonical` 设而不读 | 由 `finalize` 消费 | **CLOSED**（见 M-3 行的实测） |
| N-3 | 碰撞测试重言式 + 死分支不可达 | 删死分支；注入 ghost 哈希制造真碰撞 | **CLOSED**（`idx.hashes` 里塞入源中不存在的窗口哈希，断 `contains` 仍 false，再断真片段仍命中——这才是碰撞语义；`TestFragmentIndexShortSourceNeverMatches` 也补了 populated 对照） |
| N-4 | 集成层从未走归一化 | `TestNormalizationEndToEnd` + 四通道套件的变形 subcase | **CLOSED**（同一 Mark→Inspect 链路上跑大小写/全角/空白/隐形字符 7 变体，并带「不相关串不得命中」的对称断言） |
| N-5 | 裸 `HasPrefix` 牵连兄弟目录 | 换 `isUnder` | **CLOSED**（实测 `profileevil` 不被扫进） |
| N-6 | DEFERRED 点名的票 21/26/55 正文无承接 AC 行 | 票 26 / 票 55 各补 AC 行 | **CLOSED**（票 26 补「播报必须过 `CheckText(scope, ChTTS)`」；票 55 补 macOS 云盘探测须达 env/config **级**、明写「裸 `os.Stat` 猜默认目录不撤销兜底」——比我要求的还准）。票 21 那半我未在本轮重开（审批渲染面），首轮该行只作保留说明，不作缺陷 |
| N-7 | 测试少计 | 更正为 68 顶层 + 38 subtest | **NOT FULLY CLOSED（无害）**：实测 70 顶层（69 PASS + 1 SKIP）+ 37 subtest。少计方向与首轮同款，无风险，但「更正」又错了一位（见 R1） |

### R4 五条 AC 复裁（对照首轮）

| AC | 首轮 | **复验** | 依据 |
|---|---|---|---|
| ① 四/六通道逐通道 L2 + 指名来源 | PASS | **PASS（未回归）** | 五 subtest 逐条断 Channel + `hit.Source()`；TTS/HTTP 两条现断 Channel 标签（首轮 §2 行 1 的保留已消）；经真 C19 的 `TestR4EndToEndViaAssessor` 仍 L2/[R4]/SessionOverrideBlocked/「包含来自 doc.read」；B-2 加宽后归因未糊 |
| ② 归一化仍命中，<8 不命中 | PASS（附缺测） | **PASS** | N-4 补齐（集成层变形）+ M-6 双向 + 边界 8/7 仍钉 |
| ③ LLM 改写残留（记录为限制） | PASS | **PASS（不变）** | 首轮 §5(a) 的 SPEC 背书仍成立；本轮另证 M-6 **不属**该背书范围且已闭合 |
| ④ 已配置 OneDrive 机器上检出 + fixture 兜底 + 探针过 P12 | **FAIL** | **PASS** | 首轮列的三硬伤逐条消解：根确认来自 **env 级**（真机实测 root=OneDrive/env）；注册表探针有可注入表测试 + 诚实 skip；8.3 与 `\?\` 两条真写盘红队用例已补（含真机腿，门控 0 字节落地）；「default 探针命中即关兜底」被 M-3 反转为「default 级不关兜底」 |
| ⑤ 作用域不继承（DisposalScope） | PASS（附 M-1） | **PASS（M-1 已闭）** | 反转后的 `TestInspectUnknownScopeIsEmptyStore` 仍保「opened-and-empty 合法不命中」与「Dispose 后引擎无污」两条，AC#5 语义没被从严改坏 |

### R5 本轮新发现（首轮的 FAIL 面从未覆盖）

- **M-7（MAJOR，须修或显式登记后才可 ship）**：`writeChannelKeys` 的同步门**外溢到了非 `fs.write` 的工具形状**。
  `provenance.go:496` 的 `(!syncGateOpen && writeChannelKeys[k])` 出现在「其他工具」分支里，于是任何工具只要参数里
  同时出现一个 `path`/`file`/`filepath`/`dest`/`destination` 字符串（解析到同步根之外）和一个叫 `content`/`data` 的键，
  该载荷就被**静默跳过**。实测（confirmed 根、兜底关闭）：
  `http.post{url,path:<非同步>,data:<污染>}` → **hit=false**；同一形状走**冻结适配器**
  `Detector.TaintHit` → `hit=false src=""`；`notify{text:"ding",path:<非同步>,content:<污染>}` → **hit=false**；
  `web.fetch{dest:<非同步>,data:<污染>}` → **hit=false**。对照：把 `data` 改名 `body` → **hit=true**；把 path 换成同步根 → **hit=true**。
  ⇒ 命中与否由**参数名**决定，正是 B-2 的同一类，只是方向反过来（B-2 修的是「改名就扫」，这里变成「叫 content/data 就不扫」）。
  SPEC-06 §5 里唯一有落点条件的通道是 ⑤（fs.write 落同步根）；⑥（HTTP POST body/URL）**没有**落点条件，
  而票 19 自己在 `provenance.go:22-41` 写的是「every other parameter value is still scanned fail-closed」。
  公平记账：**这不是本轮引入的**（`73b6057` 的未知工具分支就有同一句 `continue`；首轮 B-2 的清单也没覆盖它），
  且今天不可达（channel ⑥ 尚无生产调用点，`NewProvenance(` 全仓仅测试里出现；net 侧推荐入口是 `CheckText(ChHTTP)`，不经此门）。
  **一行收口**：该 skip 只保留在 `def.ch == ChSyncWrite` 那一支（fs.write 自己），其余工具把 `content`/`data` 当普通参数扫。
  条件见 R7 C-1。
- **N-8（MINOR，接受，建议记一行）**：M-7 的镜像——`fs.write` 写**非同步**目录时，载荷键名只要不是 `content`/`data`
  （如 `body`）就被无条件扫 → 普通本地写也升 L2（实测 true）。方向偏严，符合 16.9#1「误报安全」，
  但它让「一次本地写入要不要确认」取决于键名，建议在 PRECHECK 的 fs.write 段补一句，免得票 20/21 接线时当 bug 报。
- **N-9（MINOR，注释与行为不符 + 非 Windows 死分支）**：`syncdirs.go:217-221` 断言「junction 算已存在，所以调用方的
  Resolve 会看见并拒绝」——实测 Go 对无尾斜杠的 junction `Lstat` 报 `isdir=false mode=?rw-rw-rw-`（不带 `ModeSymlink`），
  祖先walk 之所以穿过它靠的是补尾 `\` 后 Lstat 会跟进 mount point。结论仍安全，但**理由写错了**。
  另：`resolveTarget` 的 `pathHandleVerify == false` 分支（`:204-213`，注释称「非 Windows 上词汇形式已是 C26 能给的全部」）
  在 POSIX 上**到不了**：`deepestExistingAncestor` 用 `\` 拼路径去 `Lstat`（`splitPathComponents` 先把 `/` 换成 `\`），
  Linux/macOS 上 `\home\user\…` 永不存在 → `anc=""` → `errTargetUnverified` → **每一次写入都落 sync-suspect**（方向从严，
  非漏洞，但 macOS 票 55 会把「普通本地写全部升 L2」当新 bug）。交叉编译通过不等于该路径正确，建议随票 55 的 AC 行点一句。
- **N-10（MINOR，硬化残留，见 R2）**：尾巴含 `..` 时分类器看不见 junction，今天**与 OS 落点一致**（实测）故不构成逃逸；
  这份一致性来自 Windows 对 `.`/`..` 的前置归一，而非代码里的一句检查。若要把它变成不变式：候选原串里存在被折掉的 `..`
  组件时直接落 sync-suspect（一行，误报面近乎零——正常 `fs.write` 目标不带 `..`）。登记即可，不阻塞。
- **N-11（INFO）**：`rules_gateway.go:96-98` 的 R4 注释仍写「Dormant until ticket 19 wires the TaintDetector」，
  而引擎侧已交付（缺的是循环接线）。属票 17 文件、本票不得触碰，仅登记给票 20/21。
- 另：`Inspect` 在路径**无法核实**（而非确认同步）时也会把命中标成 `ChSyncWrite`（`provenance.go:473-476`）——
  标签偏宽，不改判定方向，INFO。

**新引入的洞？M-7 是本轮口径与实现之间真实不一致的一条（虽然根在旧码）**，其余为注释/计数级。
除 M-7 外，本轮修复未打开新的 fail-open：B-1 的祖先锚定没有把 reparse 挪进盲区（R2 两路实测），
B-2 的加宽没有冲掉归因（R3 实测），M-1 的从严没有冲掉 AC#5 的空店语义（R3 实测）。

### R6 探针复现配方（5 个临时文件 / 23 个探针函数，跑完已删，未入树）

`P-A` 尾巴 junction（沙箱 + **真 OneDrive** 两版）｜`P-A'` 深尾巴｜`P-B` `..`+junction 的**落点**核验（两个候选落点各预置同名文件，看谁被改 + `GetFinalPathNameByHandle` 说真话）｜`P-B'` 四种 `..` 变体 + `\\?\` 腿｜
`P-C` TOCTOU（判后再建 junction）｜`P-D` M-3 诱饵重放（default/fixture/options/""/env + ghost 根）｜`P-E` B-2 六逃逸 + 归因 + 同步门 + M-1 三缝｜
`P-F` M-6 六码点双向 + 反误报｜`P-G` M-4｜`4-WhichFile` OS 真值｜`4-ExemptedJunction` 豁免穿透归因｜`4-FileSymlinkTail` hardlink｜
`5-WriteKeyGate` **M-7**（`{path,data}`/`{text,content,path}`/`{dest,data}` 三形状 × 具名工具/无名适配器）。
重跑：把上述任一形态写成 `internal/risk/zz_reverify_test.go` 后 `go test ./internal/risk/ -run TestRV -v`。
`mklink /J`、`mklink /H` 均无需管理员权限即可成功（实测），这直接决定了 R2 的威胁模型判断。
计数诚实附注：23 个探针里 round-1/2 的 4 个 `..` 形状探针（`TestRVProbeDotDotAfterJunction` +
`TestRV2DotDot{Landing,InsideSandbox,Variant}`）当时**无效**——我用 `filepath.Join` 拼原串，
它自己先把 `..` 折掉了，等于没测到未折叠的原串；结论一律以 round-3 的字符串拼接版（`TestRV3*`）与 round-4 的
同名文件落点版（`TestRV4WhichFileDoesTheOsOpen`）为准。另有两个探针首跑因**我的**前置条件写错而红
（`TestRVProbeB2` 忘了把注入根真建出来 → 兜底开着，判「普通写也命中」；`TestRVProbeM6` 的负例 `marker-qzx` 本是
marker 的前缀 → 真命中不是误报），修正前置后复绿；红的那两条不是被测码的缺陷，记此以免证据被误读。

### R7 VERDICT

**PASS WITH CONDITIONS — 票 19 可继续推进，但下列三条须先落地（都不需要重开 B-1/B-2 的面）。**

首轮必改 **B-1、B-2、M-1、M-3、M-6 全部 CLOSED 并经我独立实测复现其正反两侧**；建议同批的 M-2/M-4/M-5/N-1..N-6 亦 CLOSED；
AC 复裁：①②③⑤ PASS，**④ 由 FAIL 转 PASS**（env 级根确认 + 8.3/`\?\` 真写盘红队用例，正是首轮 §8 开出的判据）。
两条首轮实测复现的 fail-open 逃逸，本轮我**用更强的形状（尾巴 reparse、`..` 出圈、TOCTOU、豁免 junction）都未能重新打开**。

必须落地的条件（回归前，一行码 + 一行文档的量级）：
- **C-1**：**M-7** 二选一——(a) 把 `content`/`data` 的同步门 skip 收回 `def.ch == ChSyncWrite` 一支；或 (b) 在 `provenance.go` 头部 DEFERRED 头与票 22 正文各写一行 AC，明说「channel ⑥ 的 body 若以 `content`/`data` 命名且同 call 带 path，不经 C25」并给出该 AC 的测试名。(a) 更可取：它同时消掉 N-8 的键名依赖。
- **C-2**：PRECHECK 补三行残留登记（不得只留代码注释）：①「同步判定与写盘非同一原子操作，TOCTOU 不属本票威胁模型，理由见 19-adversarial-acceptance §R2」；②「`..` 折叠形状今日与 OS 归一同源故安全，硬化见 N-10」；③「env 级证据可被同用户进程伪造，撤销兜底的信任边界即此」——三者方向均为可接受，登记即免成无声技术债。
- **C-3**：票 19 日志 N-7 行的计数按 **70 顶层（69 PASS + 1 SKIP）+ 37 subtest** 改准；`syncdirs.go:217-221` 的注释按 R5 N-9 改准（连同 `pathHandleVerify==false` 分支在 POSIX 不可达这一句，或把它交给票 55 的 AC 行点名）。

**可随 DEFERRED 延后 ship（不阻塞本票转 done）**：M-5 的每候选 `[]rune` 重复转换（无时延宣称，纯优化）；N-9 的 POSIX 侧收口（票 55 AC 已挂）；N-10 的 `..` 从严硬化；N-11 的 R4 陈旧注释（票 17 文件）；LLM 改写残留（SPEC-06 §5/§12 + 16.9#1 背书，首轮 §5(a) 已裁）；`MaxSourceRunes`/`MaxScanRunes` 截断残留（已记日志）；C25 循环接线（票 10/20/21/22/26，`DEFERRED(C25-loop-wiring)` 在位且我已复核全仓 `NewProvenance(` 无生产调用点）。
另注：**`internal/agent/` 本轮零触碰**（两次提交 `grep -c internal/agent` = 0），本复验全程未读写该目录，构建范围一律 `./internal/risk/`。

---

## ⚠ 事后更正（2026-09-20，编排者，registry A26）：本文件中记为 PASS 的 `d22scan` 那条**是空跑**

本文件里 `cd tools/d22scan && go run .`（**不带 `-root`**）被记成 `clean` / PASS。
实测：该形式**默认扫描 `.`，即扫描器自己那个 module** —— 里面没有 `internal/`、没有 `cmd/`，
allowlist 读不到时被当空 ⇒ **它检查了 0 个生产文件却退出 0**。
⇒ 本文件据此支撑的"裸 goroutine / 墙钟超时 / `filepath.Clean` 越权 / emoji"四项，
**在当次验收中并未被这道门检查过**。原文数字与结论**一律不改写**（历史测量保留可读），
此段为叠加更正。
- 已由票 67（commit `23ebb59`）修好仪器本身：`checkRoot()` 让空范围致命退出，
  输出新增 **`examined N production Go files`** ⇒ **门今后必须自报工作量**，只报 `clean` 不作为证据。
- 当前 HEAD 用正确调用实测：`clean`，exit 0（本条是**事后**为那两个 AC 补上的真证据，不追溯证明当次验收有效）。
