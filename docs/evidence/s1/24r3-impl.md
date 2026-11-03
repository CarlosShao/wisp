# 24r3-impl — 票 24 落地腿 24-r3（零依赖续半：`search.files`＋`search.apps` 枚举喂 `app.launch`；PDF 仍待批）

- agent＝24-r3（实现腿／写腿）· 日期＝2026-11-03 09:3x +08 · **只 commit，不 push**
- 上游＝`docs/evidence/s1/24-v1-verdict.md` §4-②/§4-③ 与台账 `A1439`：24-v1 判 docx＋search.content 半全 PASS，整票不得加 `-done`，剩 AC#1 PDF 两格／AC#3 `search.files`／AC#4 `search.apps`→`app.launch`
- 派单存档＝`.scratch/wisp/probes/orch/20261102/txt/24r3-dispatch.txt` **起手不在盘**（`read` 现读＝not found）；与 24-r1 同款缺件，本腿**不代造**，只如实记
- `go.mod` 新增依赖＝**0**（`path/filepath`·`runtime`·`sort` 与 `golang.org/x/sys/windows/registry` 全为 stdlib 或**仓内已在用**依赖——registry 由 22-r1 的 `launcher_windows.go` 已在 import）；SPEC-01 §96 白名单**一字未动**

## 0. 起手三核＋第 19/20 处（盘上现读）

| 尺 | 读数 |
|---|---|
| HEAD 起点 | 起手 HEAD＝`46dec5c61`（现跑 `git rev-parse --short HEAD~2` 收尾复核）；收尾＝本腿两枚提交 `e1ee3b357`（代码四枚）＋`d59093153`（本件）；24-v1 记录的 `a11e63bee` 已在其之下（44 系他腿在飞），写面不重叠 |
| 写面无他脏（第 19 处） | `git status --porcelain -- internal/tools` 起手＝**空**；收尾（提交后）＝**空** |
| 票面在场（第 20 处） | `.scratch/wisp/issues/24-doc-search-tools.md` **无 `-done`**、Status: ready-for-agent、AC 复选框未勾；前置 `20-host-bridge-fs-tools-done.md` 已 done ⇒ C26 形状可复用 |
| SPEC-07 §3 | `docs/specs/SPEC-07-tools-and-plugins.md:49`（`search.files` L0／`fs.read`／白名单内）＋`:51`（`search.apps` L0／`sysinfo`／只检索不启动）＋`:48`（`app.launch` L1／shell）逐字在场 |
| PLAN D34 | `docs/PLAN.md:2549`/`:2551` 两行与上面同字（本腿现跑 grep 核对） |
| ★`search.apps`／`app.launch` 交付核（题面所命） | **`app.launch` 已由 22-r1 交付**：`internal/tools/launcher_windows.go` 的 `resolveApp`（直路径臂／AUMID 形状臂 `windowsAUMIDPattern`／name 臂→`AppResolver`→App Paths→`exec.LookPath`）＋`startApp`（argv 向量＋C30 Job）＋`launcher_other.go` 同签名件；`docs/evidence/s1/22-v1-verdict.md` AC4＝PASS ⇒ **启动半免做**。但**枚举半未交付**：`launcher_windows.go:102` 现读直写「`search.apps` 检索面尚未接入」，`web.go:113-115` 注释把 `AppResolver` 缝留给「票 23」 ⇒ **AC#4 落枚举半，启动半复用 22 的解析器，不另起第二套解析路径** |
| 未定义即停尺 | `search.files`＝票面 :22「name glob within allowlist」＋D34「按名找文件」⇒ 名片通配有定义；`search.apps`＝:23-24「installed-app enumeration feeding app.launch; enumeration only, never launches」⇒ 枚举＋喂型有定义；**唯 UWP AUMID 的零依赖枚举源无定义**（见 §4，按待批上报，不自造） |

## 1. 改动文件（4 枚全新件；`git add`／`git commit` 两处显式 pathspec 逐枚字面点名）

| 文件 | 内容 |
|---|---|
| [internal/tools/search_files.go](<internal/tools/search_files.go>) | `search.files` 工具：C26 规范化根＋`InAllowlist` 自闸＋名片通配（`filepath.Match` 语法，大小写不敏感，只作用 basename）＋三重预算（深度 8／ examined 2000／结果 50）＋符号链接不跟随不列出＋A/B 清单在打开与打印**之前**否决＋取消⇒部分结果标 truncated＋`SearchFilesDecl`／`BuiltinSearchFilesEntries`（复用 r2a 的 `SearchDeps`，不新造宿主面） |
| [internal/tools/search_apps.go](<internal/tools/search_apps.go>) | `search.apps` 工具：`AppSource` 枚举止缝（`Label`/`Apps`，只列不启动）＋`InstalledAppCandidate`（含 `AUMID` 位）＋默认源（PATH 目录可执行文件；Windows 另加 App Paths 注册表）＋**去重与上限**＋`SearchAppsDecl`（`sysinfo`／L0／⛔`shell`／⛔`PathParams`）／`BuiltinSearchAppsEntries` |
| [internal/tools/search_apps_windows.go](<internal/tools/search_apps_windows.go>) | Windows 枚举源：只读 HKLM→HKCU 的 `SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths` **子键名**（与 22-r1 `windowsAppPathsResolver` 同一 hive），两 hive 全不可读＝报错（不谎报「机器上没应用」）；键名形状门与 resolver 一致 |
| [internal/tools/search_files_test.go](<internal/tools/search_files_test.go>) | 27 枚定向测试（两族合一文件，与 r2a「docx＋search 同文件」同式）；复用 r2a 缝 `r2aCanon`/`r2aParams`/`r2aCorpus` 与 `mapResolver`，无二进制 fixture |

未触＝`internal/tools/web*.go`／`document_text.go`／`search_content.go`（已交付，只**调用**其 `SearchDeps.resolveSearch`/`classify`/`resultsCap` 等既有方法，一字未改）／`bridge.go`／`launcher_windows.go`／`internal/plugin`／`cmd/wisp`／`internal/agent`·`memory`（44-v2 在验）／thresholds·golden·SPEC-12·`go.mod`。

## 2. AC#3 `search.files` 读数

| 格 | 落点 | 钉它的测试 |
|---|---|---|
| 名片通配（正向） | `compileNameGlob`＋`filepath.Match(glob, ToLower(e.Name()))`；`a.md`/`B.MD`/`sub/c.md` 命中、`notes.txt` 不命中 | `Test24R3SearchFilesMatchesNameGlobCaseInsensitive` |
| 中文名片子集（票面场景「找上周那个合同」） | `*合同*` 子集通配 | `Test24R3SearchFilesChineseNameSubstring` |
| 只匹配名字（⛔路径通配） | 含 `/` `\` `..` 的 pattern 直接拒 | `Test24R3SearchFilesRejectsSeparatorPattern` |
| 坏输入 fail-closed | 空 pattern＝「全盘列举」拒；`[a-z` 语法坏**报错**而非「匹配 0 个」；根不是目录拒 | `…EmptyPatternIsNotAListing`／`…RejectsMalformedGlob`／`…NeedsADirectoryRoot` |
| C26 越界拒 | 复用 `SearchDeps.resolveSearch`：无 resolver 拒→`Canonicalize`→`??` 拒→`InAllowlist` 不在授权根即拒，且**不回任何结果行** | `…RefusesOutsideAllowlist`（零 `SECRET` 回显）／`…FailClosedWithoutResolver` |
| 三重预算 | `MaxDepth`（默认 8）／`MaxFiles`＝examined（默认 2000）／`MaxResults`（默认 50）；撞顶⇒`Truncated`＋头部写「达到上限，truncated」 | `…DepthCapBothWays`（含 cap 8 控制组，证明红的是预算不是走偏）／`…EntryCapTruncates`／`…ResultCapTruncates` |
| 不跟随符号链接 | `e.Type()&fs.ModeSymlink` ⇒ 不打开、不列出，只给匿名计数「符号链接跳过 N 项」 | `…SkipsSymlinks`（本机无权限时诚实 Skip） |
| 敏感清单不泄元数据 | 每个候选在打印**之前**过 `risk.Classify`（可注入接缝，默认＝真分类器）；根也被否决 | `…NeverListsBlacklistedFiles`／`…BlacklistedRootIsRefused` |
| 取消⇒部分结果标 truncated | 显式栈遍历＋每帧查 `ctx.Err()`；`aborted` 非空⇒`Truncated` | `…CancelMarksPartialResult` |
| D34 行形状 | `fs.read`／L0／`PathParams=["path"]`（`pattern` **故意不进** PathParams：它只作用 basename，不是路径）／Resident／KindBuiltin | `…DeclMatchesD34Row`／`…RegistersInRegistry` |
| 只列文件 | 目录被遍历但不作命中（D34 行＝「按名找文件」） | 同正向测试内断言 `d.md` 目录不出现 |

## 3. AC#4 `search.apps`→`app.launch` 读数（启动半核对＋枚举半落地）

- **启动半＝22-r1 已交付，本腿不重做**：`resolveApp` 三臂（直路径／AUMID 形状门／name→`AppResolver`→App Paths→`LookPath`）＋`startApp`（argv 向量、D33 ⛔字符串拼接、入 C30 Job 与未接线的诚实姿态）现读在场，`22-v1-verdict.md` AC4 已判 PASS。
- **枚举半＝本腿落地**，且喂型是**结构性的**而不是许诺：每个候选名都过 `webShared.resolveApp`（＝`app.launch` 自己那套解析器），只有解析成功的才列出，并直接列出它给出的 `kind`（`path`／`aumid`）与目标串；解析不到的候选**计数后丢弃**（头部「未能解析 N」），所以清单不会宣传 `app.launch` 起不来的应用。
  - 测试：`Test24R3SearchAppsEnumeratesFixtureAndFeedsLaunch`＝题面 AC#4 一枚到位：fixture 源里的 `fixtureapp` 经注入 `mapResolver` 枚举出、列出的目标串＝解析器给的同一绝对路径、`ghostapp` 不出现且被计数。
  - **只检索不启动**有双尺：① 结果行必含「未启动任何应用」，且启动姿态文字（`已纳入 Job`／`未纳入`／`启动载体`）出现即判失败（出现＝真启动了）；② `Decl` 只声明 `sysinfo`、⛔`shell` ⇒ 一旦谁在这个工具里加启动调用，桥按 SPEC-07 §2 硬拒（`Test24R3SearchAppsDeclIsSysinfoL0NoShell` 钉住名册）。
- 其余读数：`query` 是字面、大小写不敏感子集（⛔正则，与 `search.content` 同语义，`…QueryIsLiteralSubstring` 用 `note.*` 反证）；跨源大小写去重（`…DedupeAcrossSources`）；`max_results` 撞顶⇒`Truncated`（`…ResultCapTruncates`）；**所有枚举源报错＝IsError**，绝不叙述成「这台机器没装应用」（`…AllSourcesErrorIsNotAnEmptyMachine`）；坏 JSON 参数⇒错误结果（`…RejectsBadParams`）；默认源非空且必含 PATH 地板（`…DefaultAppSourcesHasAtLeastThePATHFloor`）。

## 4. UWP AUMID 枚举半：**未落，按依赖/裁定上报**（⛔自造）

题面把「解析 name/path/UWP AUMID 喂 app.launch」列为 AC#4 的一种可能形状。现盘事实：

- **消费端已就绪**：`app.launch` 有 AUMID 臂（`windowsAUMIDPattern`＋`explorer.exe shell:AppsFolder\!AUMID`，姿态文字诚实）；本腿的枚举接口也已备好——`InstalledAppCandidate.AUMID` ＋ `kind=aumid` 列出形 ＋ 测试 `Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm`（Windows 真跑，非 Windows 诚实 Skip）。**只要拿到一个权威 AUMID 源，接线即完成，无需改形状。**
- **零依赖的"枚举源"不存在**，两条主流路都越出本腿写面：① `Get-AppxPackage`＝从一枚 L0／`sysinfo` 工具里 spawn PowerShell 去起子进程——正是本工具「只检索不启动」要否决的副作用形状，且能力面要 `shell`（改声明＝契约面）；② Appx 打包 API／`IShellFolder(shell:AppsFolder)` COM＝**新增第三方/系统绑定依赖**，白名单变更须人工批准（`AGENTS.md` §1.1、SPEC-01 §96）。③ 解析 `.lnk` 的 PS property store 取 `AppUserModelID`＝无文档保证的二进制自造（且 `.lnk` 路径喂 `app.launch` 的 `exec.Command` 起不来＝违 AC#4 的"喂得动"）。
- ⇒ 依「未定义即停＋引第三方＝停手上报」：本腿**不选任一自造路**，把这一半记为待裁（建议编排者按铁律「最有价值＋主流 harness 方向」拍：要么批准 `sysinfo` 枚举止步走 spawn/COM，要么把 AUMID 枚举并入 D34 表的 `system.*`/接线票一并定），代码侧接口已就位。

## 5. 仍待批／仍待派（不在本腿写面）

1. **PDF 两格（AC#1）维持待批**：pdfcpu 不在 `go.mod`、不在 SPEC-01 §96 白名单⇒白名单变更须人工批准（`A1435` 已摆机主）；`24-r2a` 的 stub（`errors.Is(ErrUnsupported)`＋DEFERRED 文案）一字未动，`.pdf` 仍是明确报错而非静默。
2. **`DEFERRED(D34-1)` 五字段行仍未追加**（SPEC-12 §5 是契约文档⇒须编排者/机主具名落表，`24-v1-verdict` §4-③ 的结案前置条件未闭）。
3. **bridge/组合根接线一枚**：本腿两枚 entry（`BuiltinSearchFilesEntries`／`BuiltinSearchAppsEntries`）与 r2a 一样**未接 `bridge.go`**（派单⛔碰 `web*.go`，`bridge.go` 属他腿/接线票写面）；per-tool timeout 来自 config、>4k spill 落 artifacts 两格仍属宿主层（`24-v1-verdict` §3.5-④ 两观察随该票验）。
4. **过期文案上报**：`launcher_windows.go:102` 的错误文字「`search.apps` 检索面尚未接入」在本腿落地后**过期**（该件＝22-r1 已交付面，本腿⛔改）。这条文案只在该文件 name 臂未命中时出现；改它请随接线票或 22 的续票点名。

## 6. 形状禁令自查（新件四枚）

- 裸 `go func(`＝**0**（本腿不起任何 goroutine）；超时＝ctx deadline，⛔墙钟差。
- `risk.PathResolver` 之外的 `filepath.Clean|Abs`＝**0**：`search_files.go`/`search_apps.go` 只用 `filepath.Match`（通配语法）与 `filepath.SplitList`/`Ext`（字符串形状，不做文件系统决策）；根一律 `Canonicalize`＋`InAllowlist`，子项用既有 `joinForListing` 拼在**已规范**的父路径上。
- 明文密钥＝0；L2「允许」不由面板侧来源＝不涉及；受门控 Tool 形状＝本腿两枚都是只读 L0。
- emoji： ripgrep `[\x{2190}-\x{2BFF}\x{FE0F}\x{1F000}-\x{1FAFF}]` 命中＝**0**（含箭头段，注释与字符串都不放过）。
- `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`）四枚＝**归空**。

## 7. 真跑读数（只 scope `./internal/tools`；日志 `D:\tmp\wisp\20261102\24r3\test-24r3.log`，只建不删）

| 尺 | 读数 |
|---|---|
| `go build ./internal/tools`（GOOS=windows） | **零输出＝通过**（go1.27.1 windows/amd64） |
| `GOOS=linux go build ./internal/tools` | **零输出＝通过** ⇒ `resolveApp`/`launchTarget` 跨平台可用，`platformAppSources` 在非 Windows 保持 nil 只走 PATH 地板 |
| `go test ./internal/tools -run 'Test24R3' -count=1 -timeout 20m -v` | **PASS 27／FAIL 0**（SKIP 0，唯一候选是 …SkipsSymlinks：本机符号链接权限不足即诚实 Skip，不假绿），`ok github.com/CarlosShao/wisp/internal/tools 0.307s`，exit 0 |
| 还原核 | 代码件提交＝`e1ee3b357`（4 files changed／1334 insertions／**0 deletions**，name-status 全 `A`）；收尾 `git status --porcelain -- internal/tools`＝**空** ⇒ §7 跑的就是入库的码；`document_text.go`/`search_content.go`/`web*.go`/`launcher_windows.go` blob 未变（本腿零改动，只调用） |
| 反向三问 | ①越界：根不在授权目录⇒IsError 且零回显（`…RefusesOutsideAllowlist`）；②坏输入：空/分隔符/坏 glob/非目录根⇒IsError（四枚）；③泄露：A 级文件名与正文都不出现（`…NeverListsBlacklistedFiles`）⇒全绿。⛔ 未做摘保护 mutation 反转（越权改码），以「保护调用点在 `Execute`/`run` 主干＋接缝默认＝真 `risk.Classify`／真 `resolveApp`」代替 |

## 8. Git 纪律

只 commit 不 push；`git add`／`git commit` 两处显式 pathspec 逐枚字面点名（代码四枚＝`e1ee3b357`，本件＝收尾提交）；无 `add -A`／amend／reset／rebase／stash／checkout ．／clean。票面 `24-doc-search-tools.md` **不加 `-done`**：AC#1 PDF 两格（§5-1）与 bridge 接线（§5-3）未落，AC#4 的 UWP AUMID 枚举半待裁（§4）。
