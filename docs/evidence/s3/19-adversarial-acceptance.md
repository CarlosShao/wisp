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
