# 24-v2-verdict — 票 24 验收腿 24-v2（非实现者裁决；裁决者≠实现者）

- agent＝24-v2（验收腿／只读＋本证据件）· 日期＝2026-11-03 09:4x +08 · **只 commit 本件，不 push**
- 被验收交付＝`e1ee3b357`（24-r3 实现腿：`search.files`（AC#3）＋`search.apps` 枚举半（AC#4），启动半复用 22-r1 `resolveApp`）＋证据枚 `d59093153`／`dbd82a6b4`
- 上游＝`24-v1-verdict`：docx＋search.content 已 PASS；整票不可结，剩 AC#1 PDF 两格／AC#3／AC#4 ⇒ 本腿裁剩余零依赖格＋收官问题
- 尺的边界（如实）：本腿⛔改码⇒**未做 mutation 反转**；反向判据以「测试逐枚真断言 error／零回显／不列出」＋「保护调用点在 `Execute`/`run` 主干、接缝默认＝真 C26／真 `risk.Classify`／真 22-r1 `resolveApp`」双读代替，逐条列 §5。
- 读数日志＝`D:\tmp\wisp\20261103\24v2\test-scoped.log`＋`test-Test24R3.log`（只建不删）

## 0. 起手三核（盘上现读）

| 尺 | 读数 | 判定 |
|---|---|---|
| HEAD 现盘最新 | `a7dcd2bf1`（编排者 A1441）＞`46dec5c61`；24-r3 三枚（`e1ee3b357`→`d59093153`→`dbd82a6b4`）在其之下 | **PASS** |
| 索引干净 | `git diff --cached --stat`＝**空**；tracked 脏仅 `cmd/wisp/run.go`（worktree 列 M＝**44-r5 在飞写面，本腿零触碰**）；`-- internal/tools`＝**空** | **PASS** |
| 票面在场 | `24-doc-search-tools.md` 无 `-done`、AC 复选框未勾；本腿现读全 41 行，Progress log 无新增勾选（如实半格态） | **PASS** |

## 1. 真跑读数（★只 scope `./internal/tools`，⛔ ./...）

| 命令 | 读数 | 判定 |
|---|---|---|
| `go test ./internal/tools -run 'SearchFiles\|SearchApps\|Search' -count=1 -timeout 20m -v`（题面尺，逐字） | **PASS=39／FAIL=0／SKIP=0**，exit 0，`ok … 0.463s`；＝24-r3 26 枚＋外来 13 枚（r2a `search.content` 族＋`Test22*Search*` 同包他票测试，**也全绿**） | **PASS** |
| 补跑 `go test ./internal/tools -run 'Test24R3' -count=1 -timeout 20m -v` | **PASS=27／FAIL=0／SKIP=0**，exit 0 ⇒ 题面正则漏第 27 枚＝`Test24R3DefaultAppSourcesHasAtLeastThePATHFloor`（`Sources` 不含子串 `Search`，正则不匹配），补跑钉死；与 24-v1 补跑先例同款处置 | **PASS** |

★关键：`SkipsSymlinks` 与 `AUMIDCandidateFeedsTheAUMIDArm` 均**真跑 PASS（SKIP=0）**——本机符号链接可用、Windows AUMID 形状门实跑，非诚实 Skip 兜底。

## 2. 还原核

- `git log --name-only 46dec5c61~1..a7dcd2bf1 -- <票面/台账/SPEC-07/SPEC-12>`＝只有编排者两枚 A 记，**24-r3 三枚零触碰**。
- 四枚新件 `e1ee3b357..HEAD`＝**无后续改动**；盘面 blob＝HEAD **四枚全同**（`search_files.go 8ba614c42e32`／`search_apps.go 9813e0edcbfd`／`search_apps_windows.go c338a00f5175`／`search_files_test.go 3b846afbc1c8`，`git hash-object` 现算）⇒ §1 跑的就是入库的码。
- 保护面 blob 比对（`e1ee3b357^` vs HEAD）：`launcher_windows.go`／`launcher_other.go`／`web_open.go`／`document_text.go`／`search_content.go` **五枚全 SAME**；`git show e1ee3b357 -- go.mod go.sum`＝**空**（零新依赖，`x/sys/windows/registry` 为 22-r1 已在用）。

## 3. 逐格裁决 — AC#3 `search.files`（`search_files.go` 262 行）

| 格 | 读数（源码＋测试双读） | 判定 |
|---|---|---|
| 名片通配 | `compileNameGlob`：TrimSpace＋`filepath.Match` 冻结语法（`*`/`?`/`[chars]`），**只作用 basename**（`:226` `ToLower(e.Name())`）；大小写不敏感；⛔正则⛔`**`路径通配＝不发明；`MatchesNameGlobCaseInsensitive`（5 文件三命中，目录 `d.md` 不作命中＝D34「按名找文件」）＋`ChineseNameSubstring`（`*合同*`＝票面场景 :3135「找上周那个合同」现字）真跑绿 | **PASS** |
| 分隔符 pattern 拒 | 含 `/` `\` 的 pattern **编译期即拒**（`:154`），`../etc/passwd` 形状三式断言 `IsError`＝名片≠路径的结构性证明 | **PASS** |
| 坏输入 fail-closed | 空 pattern＝「全盘列举」拒（`:151`）；`[a-z` 坏语法**报错**而非「匹配 0 个」（`:160` 用 128 长探针名强制 `ErrBadPattern` 早现——这是给坏输入上牙齿的真细节）；非目录根拒（`:99`）；无 resolver 拒。四枚测试真跑绿 | **PASS** |
| C26 越界拒 | **复用 r2a `SearchDeps.resolveSearch`**（`search_content.go:140` 原方法，一字未改）：无 resolver→`Canonicalize`→`??`→`InAllowlist` 不在授权根即拒；`RefusesOutsideAllowlist` 断 `IsError`＋文案含 `fail-closed`＋**零结果行回显**（corpus 含 SECRET 不外泄）；`Decl.PathParams=["path"]` 让 C19 对同一规范形判 R2/R3（`pattern` 故意不进＝只作用 basename） | **PASS** |
| 三重预算 | `MaxDepth`(8)／`MaxFiles`＝examined(2000)／`MaxResults`(50)＋args 可降不可升（`:104` `a.MaxResults < limit`）；三式分钉：`DepthCapBothWays` 含 **cap 8 控制组**（证明红的是预算不是走偏——好牙齿）；`EntryCapTruncates`／`ResultCapTruncates` 断 `Truncated`＋头部「达到上限，truncated」 | **PASS** |
| 符号链接不跟随 | `e.Type()&fs.ModeSymlink`⇒不打开、不列出、只匿名计数（`:222-224`）；ReadDir 的 DirEntry 来自 lstat 语义⇒**指向目录的链接也不钻**；`SkipsSymlinks` 本机真跑 PASS（非 Skip） | **PASS** |
| 敏感清单不泄元数据 | 每候选在**打印之前**过 `classify`（`:212`，接缝默认＝真 `risk.Classify`），A/B 命中⇒排除＋匿名计数；根也被否决（`:92`，`BlacklistedRootIsRefused`）；`NeverListsBlacklistedFiles` 断 SECRET 名零出现＋可见文件仍在 | **PASS** |
| 取消⇒部分结果标 truncated | 显式栈遍历（无递归栈爆）＋每帧查 `ctx.Err()`；取消⇒`aborted`＋`Truncated`＋「已中止」，**不是整调用失败也不是假装扫完**（票面 :27）；`CancelMarksPartialResult` 经 classify 缝注入 cancel 真跑绿 | **PASS** |
| D34 行形状 | `SearchFilesDecl`＝`fs.read`／L0／Resident／KindBuiltin，与 `PLAN.md:2549`＋`SPEC-07 :49` **本腿现读逐字对得上**；`RegistersInRegistry` 过 C1 校验 | **PASS** |

## 4. 逐格裁决 — AC#4 `search.apps` 枚举半（`search_apps.go` 355＋`search_apps_windows.go` 94）

| 格 | 读数 | 判定 |
|---|---|---|
| ★复用 22-r1 `resolveApp`，未另起第二套 | `search_apps.go:261` 调 `shared.resolveApp`＝`webShared` 方法本体（`launcher_windows.go:79`／`launcher_other.go:50`，**blob SAME＝一字未改**）；`app.launch` 自己走同一方法（`web_open.go:365`）⇒ 全 package 只有这一套解析器。`webShared` 以最小字段面构造（deps.AppResolver＋logf），现读 `resolveApp` 主体只触这两处＋平台缺省 ⇒ 构造合法非 hack | **PASS** |
| 喂型结构性（AC#4 的「喂得动」） | 候选**只有过 `resolveApp` 解析成功才列出**，列出的 `kind`（path/aumid）与目标串＝解析器原样输出＝`app.launch` 消费的形状；解析不到的**计数丢弃**（「未能解析 N」）⇒ 清单不宣传起不来的应用。`EnumeratesFixtureAndFeedsLaunch` 一枚到位：fixture 枚举行＝resolver 给的同一绝对路径、`ghostapp` 不出现且被计数 | **PASS** |
| 只检索不启动（三重锁） | ①源码 grep：`startApp`/`exec.Command`/`ShellExecute`/`cmd.Start` 在三件内**零调用**（唯一命中＝注释自述）；②结果行必含「未启动任何应用」，启动姿态文字（`已纳入 Job`/`未纳入`/`启动载体`）出现即判失败＝测试反向钉；③`Decl` 只声明 `sysinfo`／⛔`shell`⇒日后谁加启动调用，桥按 SPEC-07 §2 硬拒（`DeclIsSysinfoL0NoShell` 钉名册） | **PASS** |
| Windows 枚举源 | 只读 HKLM→HKCU `…\App Paths` **子键名**（`search_apps_windows.go:42/:53`）＝与 22-r1 `windowsAppPathsResolver` **同一 hive**（`launcher_windows.go:158` 本腿现读比对）⇒「列出的必可被 app.launch 解析」是构造性成立；**注册表值数据不读不信任**（路径由解析器按启动侧规则给）＝不再开第二个路径决策点；键名形状门与 resolver 一致；**两 hive 全不可读＝报错**，不谎报「机器没装应用」 | **PASS** |
| 缺省源地板 | `DefaultAppSources`＝平台源（Windows 注册表）＋PATH 地板（每平台）；`DefaultAppSourcesHasAtLeastThePATHFloor` 钉「缺省枚举集永非空」（空集＝静默零结果的真隐患，有牙） | **PASS** |
| 预算/去重/查询语义 | `MaxResults`(50)＋`MaxCandidates`(2000/源)＋跨源大小写去重；query＝字面大小写不敏感子集（与 `search.content` 同语义），`note.*` 反证非正则；撞顶⇒`Truncated`＋头部文字 | **PASS** |
| 全源报错≠空机器 | `len(apps)==0 && errored==len(sources)`⇒`IsError`「检索失败：所有枚举源均报错」（`:306-310`）——不叙述成「这台机器没装应用」；部分源失败＝匿名计数入头部＋继续（诚实混合态） | **PASS** |
| D34 行形状 | `sysinfo`／L0／⛔PathParams（不读用户路径）／Resident／KindBuiltin＝`PLAN.md:2551`＋`SPEC-07 :51`「只检索不启动」现读逐字对得上；注册过 C1 校验 | **PASS** |

### 4.1 UWP AUMID 枚举源未落＝如实上报，不判 FAIL

AC#4 字面＝「enumerates **fixture-installed** app; feeding result to app.launch works」⇒ fixture 枚＋喂型已闭环（§4 首两行）。AUMID 的**零依赖权威枚举源**在盘上无定义：三条候选路各自越界（spawn PowerShell＝L0 sysinfo 工具起子进程、正是「只检索不启动」要否决的形状且需改能力声明＝契约面；Appx/COM API＝新增依赖须人工批准；`.lnk` property store＝无文档保证的自造且喂 `exec.Command` 起不来＝违「喂得动」）。24-r3 依「未定义即停＋⛔自造」停手上报（impl §4），且**接口已就位**：`InstalledAppCandidate.AUMID`＋`kind=aumid` 列出形＋`AUMIDCandidateFeedsTheAUMIDArm`（Windows 真跑 PASS）——拿到批准的源即接线，零形状改动。**此判定＝PASS（随待裁上报），非缺格。**

## 5. 反向核（题面点名三问，全部真跑＋读码双证）

| 反向问 | 牙齿 | 读数 |
|---|---|---|
| 越界拒 | `RefusesOutsideAllowlist`＝根不在授权⇒`IsError`＋零回显；`FailClosedWithoutResolver`＝无 C26 不走路 | 真跑 PASS；调用点＝`Execute` 主干第一道（`resolveSearch`），非测试脚手架 |
| 符号链接不跟 | `SkipsSymlinks`＝不列出＋不进入＋匿名计数 | 真跑 **PASS（SKIP=0，本机链接可用）** |
| 坏输入 fail-closed | 空/分隔符/坏 glob/非目录根/坏 JSON（apps 侧）＝五枚各断 `IsError`＋点名文案；`ghostapp`＝解析不到不列出 | 全真跑 PASS |

⛔ 未做摘保护 mutation 反转（越权改码）；以上述「主干调用点＋默认真实现」代替，与 24-v1 同尺。

## 6. 纪律（`git show` 逐枚）

| 尺 | 读数 | 判定 |
|---|---|---|
| stat 与题面符 | `e1ee3b357`＝**4 files／1334 insertions／0 deletions／name-status 全 `A`**＝`search_apps.go`(355)＋`search_apps_windows.go`(94)＋`search_files.go`(262)＋`search_files_test.go`(623)；A1441 题面「5 新件」＝代码四枚＋本腿补的收尾证据枚（`d59093153` 90 行＋`dbd82a6b4` 补 HEAD 读数 1 行）——逐符 | **PASS** |
| 写面 | 代码全部落 `internal/tools` 搜索族新件；⛔件（`web*.go`／`document_text.go`／`search_content.go`）＋22 系 launcher 两件 **blob SAME**；`bridge.go` 未接线＝与 impl §5-3 自述一致（grep 现证）；票面/台账/SPEC **零触碰**（§2）；`go.mod/go.sum` 零改 | **PASS** |
| 形状禁令 | 裸 `go func(`＝0；`risk.PathResolver` 外 `filepath.Clean\|Abs`＝0（唯一命中＝注释自述，实际只用 `Match`/`SplitList`/`Ext` 字符串形状）；emoji＝题面全段带 `[\x{2190}-\x{2BFF}\x{FE0F}\x{1F000}-\x{1FAFF}]`（含箭头段，比仪器宽）rg **四件 0 命中**；`gofumpt -l` 四件归空；`go vet` exit 0；`go build` windows＋linux 双 GOOS exit 0（跨平台件 `platformAppSources` 非 Windows 保持 nil 只走 PATH 地板，实核） | **PASS** |
| 证据自述与盘符 | impl §7「27 枚全绿」＝本腿补跑 27/27 复现；§0 起手 HEAD（`46dec5c61`）＝`dbd82a6b4` 补齐后与提交序相符；§1 派单存档「起手不在盘」＝本腿现查 `probes/orch/**/24r3*` **仍为空**（缺存档，编排者面处置，见 §8-②） | **PASS** |

## 7. ★收官裁决（题面三问）

**① `search.files`／`search.apps` 是否 PASS＝是。** §3＋§4 逐格 PASS：39＋27 真跑 FAIL=0、还原核四枚全同、越界/链接/坏输入三反向有牙、喂型＝同一 `resolveApp` 未起第二套、只检索不启动三重锁、零新依赖、写面纪律无违。AUMID 枚举源未落＝合规停手上报（§4.1），不占 AC 字面。

**② 24 剩余是否仅 PDF＝基本成立，精确说＝「两枚待批＋一枚非 AC 面」。**
- **AC#1 PDF 两格**：维持待批——pdfcpu 白名单属「改契约文字」（AGENTS.md §1.1＋SPEC-01 §96），A1435 摆机主、本腿现读台账**无批准痕迹**；stub 三处 `DEFERRED(D34-1)` 标记未动（`.pdf` 仍是明确报错非静默）。
- **`§5 DEFERRED(D34-1)` 登记**：本腿现读 `SPEC-12` grep＝**`D34-1`/`pdfcpu` 零命中**（xlsx/OCR 具名行在 :78，两式标记的 1:1 双向尺对 `D34-1` 仍开环）⇒ **未落，契约文字待机主/编排者具名追加**。这是 **AC#5「registry cross-check passes」的结案前置**：追加前 AC#5 复选框不能勾（marker 在场✓、cross-check 未闭✗）。
- **bridge/组合根接线**：未接（现证 `bridge.go` 零命中，impl §5-3 自认）——**不属票 24 的 AC 五格任何一格**，随接线续票走（24-v1 §4-③ 既议），如实列出、不算「24 剩余 AC」，但⛔忘账。

**③ 24 能否在「PDF 待批」下结＝「零依赖 AC 格」全 PASS 可结；整票 −done 仍不可，且不取决于代码。** 明确处置建议：
1. **本腿同意结**＝24 的全部零依赖格：AC#1 的 docx 半（24-v1 已结）＋AC#2（24-v1 已结）＋**AC#3／AC#4（本腿结）**。
2. **整票收官的卡点只剩两张待批纸**：pdfcpu 白名单（机主）＋§5 `DEFERRED(D34-1)` 五字段行（契约文字，机主/编排者具名）。**建议 §5 行先行**——它不依赖 pdfcpu 批不批：marker 已在码里（r2a 起），1:1 双向尺现开环，先落行即 AC#5 可勾、开环即闭；逐字候选行 24-v1 §4-③-1 已给（五字段齐）。
3. 落行后若 pdfcpu 仍待批 ⇒ 整票按 v1 既案处理：**开「24b：PDF 三支 AC」续票挂白名单**，票 24 本体 AC 除 #1 外全勾；#1 是 PDF＋docx 复合格，**⛔全勾即⛔整票 −done**——除非机主拍板把 #1 拆格改写票面（那是票面文字，属机主/编排者面，非本裁决权）。
4. 接线一枚随 bridge 续票（含 24-v1 §3.5-④ 两观察）。

## 8. 小观察（不判 FAIL，随对应面处置）

① `pattern:"*"` 合法＝有界列举，与 `find -name '*'` 同形（预算＋黑名单＋授权根三兜都在），非漏洞；② **派单存档缺件复发**——24-r3 dispatch 未存档（impl 如实记、本腿现查仍空），「先存档再发」执行面漏，归编排者面；③ `launcher_windows.go:102` 文案「search.apps 检索面尚未接入」在 24-r3 后**过期**（impl §5-4 已报；22 系件⛔本腿与 r3 皆不可改），随接线票改字；④ apps 排序按名大小写敏感（纯展示序，无正确性影响）；⑤ `filesScan` 命中数恰等于上限时最后一个命中也报「达到上限 truncated」＝过报方向，安全侧。

## 9. 结论

- 逐格：**§3 search.files＝PASS｜§4 search.apps 枚举半＝PASS（AUMID 源＝合规待裁上报）｜§5 反向三核＝有牙真跑｜§6 纪律＝PASS｜§2 还原核＝PASS｜§1 真跑 39+27 FAIL=0**。
- 收官：**同意结 24 的零依赖部分（AC#2/3/4＋AC#1 docx 半）＝票 24 四格中三枚整格闭合**。
- 整票处置：**不同意 −done**。卡点＝①pdfcpu 白名单（机主待批）②SPEC-12 §5 `DEFERRED(D34-1)` 行（契约文字待具名追加，**AC#5 前置，建议先行落行**）③接线随续票。24 剩余**确实只剩 PDF 一格**（＋上述两张纸），无零依赖格再欠——⛔再以「等 PDF」挂任何零依赖工作。
