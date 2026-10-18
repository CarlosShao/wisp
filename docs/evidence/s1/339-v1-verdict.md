# 票 339 终裁腿 339-v1 判语（2026-10-18）

- 本腿＝票 339 非实现者终裁腿（⛔ 实现者、⛔ 编排者）。实现腿＝`339-r1`（甲形拆分）。本件⛔ 翻框⛔ 改票名／`Status:`／`-done`。
- 起手现跑（本腿第一笔）：HEAD＝`98dbf20a35027a9290f52649d8ab2634646cb14e`；`git status --porcelain -- internal/risk`＝空。
  名册尺＝严格前缀逐枚 hash（`票 339 落地腿 339-r1`）＝**10 笔**／57 枚去重路径／`.go` 恰 **2** 枚／越面 **0**
  （子串尺 `339`＝67 笔、`票 339`＝26 笔＝坏尺具名⛔ 用）。票框终量尺现跑：未勾 `^- [ ] **AC#`＝**7**、已勾 `^- [x]`＝**0**＝一枚没翻。
- ★台面＝本腿自建仓外导出树（落点 `D:\tmp\wisp\20261018\339v1\`：`pre`＝e6fd0f50／`post`＝7437ef88／naive／probe／trip／dtl／dtw／dtb／winpre／winpost／realres；
  七棵以上全证非空 764–765 枚 `.go`）。所有决定性读数＝本腿自己重打，⛔ 抄腿／编排者的数；逐字日志在 `.scratch/wisp/probes/339/v1/logs/`。

## 逐格判语〔尺＋射程目录＋量/算〕

**AC#0 起手现跑＝成立。** Linux 整包 `pre` 树 RUN=199/PASS_TOP=114/FAIL_TOP=3/SKIP=1 rc=1（红逐名：`TestAnchorMissReasonNamesCandidateAndMissingAnchor`、
`TestAnchorMissKeepsNonOverridableClassAndNeedL2`、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`）⇄ `post` 树 **195/115/1/1** rc=1（残红只剩 C26，SKIP 同枚在册）
＝与腿报⇄编排者复跑**逐字同数**＝收件（量的，尺＝`^=== RUN`/`^--- PASS`/`^--- FAIL`/`^--- SKIP`，`v1-linux-fullpkg-{pre,post}.log`）。
Windows 原生 `winpre`⇄`winpost` 两棵树 **247/140/1/1** rc=1 两发，名册 247⇄247 `Compare-Object`＝**IDENTICAL**（量的，`v1-win-fullpkg-*.log`＋`roster-*.txt`）。
★同族全仓尺（量的，`census-root-dependency.sh`）：真同族形＝"断言依赖 handle-resolved root 而无平台档"——HEAD 的 `internal/risk/**` 命中 **0 枚活件**
（唯一无档引用 `name_claim_two_arm_216_test.go:79` 是**注释句**，其断言在两台面逐名 `--- PASS`＝本腿四发读数）；正控＝同一把尺打旧锚 e6fd0f50 点名 **2 枚**无档件（334＋216）＝尺灵。
⚠ 具名一枚坏尺：`census-narrowed.sh`（按"设 Windows-only env 而无档"筛）＝0 命中且**结构上抓不到本缺陷**——正控现量 e6fd0f50 的 334 件 `:85,86,87,88` 同发设四枚 env（含 `HOME`）＝票面第 9 行那句诊断作废的同一事实。

**AC#1 三形代价表＝成立（甲＝更正后本票面内唯一合法选形）。**
① 票面"乙＝假主目录按 HOME⇄USERPROFILE 分别设"按字面已作废（量的：`:86` 早设 `HOME`；本腿探针 `fakeHomeEffective=true`，`v1-probe-posix-mechanism.log`）；
真乙＝动 `pathresolver_other.go`＝AC#5 禁区＋head line 5 未定义即停⇒甲合法、且＝盘上先例形（票 336 同形已结案）。
② 甲的分母＝本腿自量（名册 comm）：Linux 少 5 枚 RUN 名（顶层 2＋子用例 3，逐名见 `roster-*-diff`）、多 1 枚（`TestAnchorMissHasNoWorldOnThisPlatform339`）；Windows 名册零变化（逐字 IDENTICAL）。
naive 整档加 tag 打坏 Linux 编译＝本腿复现（量的：`vet: anchor_fallback_teeth_336_other_test.go:73:30: undefined: hostile334`＋`name_claim_two_arm_216_test.go:89:9: undefined: candidateOffAnchorChain334`，
VET_RC=1；落刀前后 tag 行 0→1、还原 match=YES，`v1-jia-naive-buildbreak.log`）。丙＝腿测 staleness 互斥〔本腿未复量，具名〕。

**AC#2 ★＝部分成立（三选一，凭据全部本腿重打）。**
Windows 半区＝成立：未突变正控 **22/3/0/0 rc=0** ⇄ 落点A（`aRelates` 尾行摘回落支）**22/0/3/0 rc=1** 三枚点名红 ⇄ 落点B（`aRuleReason` 短路恒真、判级不动）**22/1/2/0 rc=1**
（＝成对两发落点不同；两枚 `want A` 承重整句在落点B 仍绿是**结构正确**——它钉的是判级，落点B 不动判级；红的两枚点名 reason 维）。量的，`v1-deteeth-*.log`＋`v1-posctl-*.log`。
两发 pristine＝`blacklist.go c1799cd2…` 跑前验⇄还原后 match=YES 双次；落刀前旧锚命中＝1、落刀后＝0、新锚＝1（两枚都验）。
POSIX 半区＝**本票写面内结构不可达**（见必答 vi 本腿实测三锚）⇒ 摘牙零变化＝去牙后 **1/1/0/0 rc=0** 与未突变正控一字不变（量的）。这一维＝缺世界已登记（腿具名停手上报在票尾＋行 39＋台账 A1040；后续载体＝必答 (iv) 新票），⛔ 用动产码换它＝越 AC#5。
★绊线不是装饰（必答 ii）：突变 `resolveHandle`→可解析（导出树物理字节）⇒ `anchor_miss_reason_339_other_test.go:140`／`:143` 点名红、**RC=1**；
哈希 `P0=37c42e94…` ⇄ `MUT=48f26a7c…` ⇄ `R=37c42e94… match=YES`；未突变同树 `1/1/0/0 rc=0`＝不恒红。量的，`v1-tripwire-sensitivity.log`。
⛔ 放宽断言换绿＝没发生：`want A`／`ClassA` 承重句逐字在场、两枚仪器件 `t.Skip` 命中 **0**（尺⑥，量的）。
⇒ 本格"两平台都点名红"按票面原句＝不可满足；按本格意图（⛔ 藏红、⛔ 放宽）＝甲形交出的东西成立。判**部分成立**：Windows 有牙成立＋POSIX 半区缺世界已登记。

**AC#3 票 334 的账＝射程要收窄，收窄句如下（本腿⛔ 动票 334 一字，句交编排者落笔）。**
现读原文（票 334 判据区 AC#3 格逐字）：`让 internal/risk/name_claim_two_arm_216_test.go 在**CI 那把台面**（GOTMPDIR 按 %LOCALAPPDATA% 组、TMP/TEMP 不动，照 scripts/portable-tests.sh:152-159 那个形）下绿`——
该句主语（216 仪器）在两平台至今为真（量的：本腿 `--- PASS`×台面 Linux⇄Windows）。要收窄的是 `334-r1` 票尾那句"链内那形在两棵树上都可判"⇄334 件头那句 `decidable on every machine`：
腿已把头句改成 `every WINDOWS machine`＝**本腿独立裁＝算说实话、⛔ 算偷改别票**（票 334 的 issues 文件被 339-r1 被动枚数＝0，尺④现跑；.go 那件在票 339 AC#6 声明写面内）。
可直接粘贴的收窄 bullet（建议编排者追加到票 334 末尾，⛔ 撤框⛔ 改判据原文）：
> [2026-10-18 22:0x +08] **射程收窄（据票 339 终裁腿 `339-v1` 判语 `docs/evidence/s1/339-v1-verdict.md` 由编排者追加；⛔ 本句＝撤 AC#3 那格）**：AC#3 那句"CI 那把台面下绿"的主语（`name_claim_two_arm_216_test.go`）两平台至今为真（`339-v1` 现量 Linux⇄Windows 仓外台面逐名 `--- PASS`）。需要点名收窄的是 `334-r1` 票尾"链内那形在两棵树上都可判"一句＝对 POSIX 不真：`pathresolver_other.go:10` 的 `resolveHandle` 恒 `("",false)` ⇒ `formsOf` 的 `pathForms.root` 在 POSIX 恒空 ⇒ `blacklist.go:359-361` 的 `uncertainAnchorMiss` 首行守卫恒 `return false` ⇒ 334 钉的 fail-closed A 在 POSIX 无世界。自 `7437ef88` 起该仪器带 `//go:build windows`，POSIX 半区由 `anchor_miss_reason_339_other_test.go` 具名点名（绊线，`339-v1` 复打：突变⇒`:140/:143` 点名红 rc=1、还原哈希 match=YES；未突变 1/1/0/0 rc=0）。

**AC#4 颜色与门禁＝成立。** 四数见上（量的，⛔ 混用⛔ 相减，同一 `-run .` 名册两平台四发）；`TestResolvePerCallBudget` 四台全 `PASS`（1.76/2.14/2.30/1.39s）＝本腿未遇争用红、⛔ 追加复量，残红 C26＝平台中性在册存量（A380／票 200／216）⛔ 记本票。
★门禁四数只在共享工作树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`（现跑）：`gofumpt -l`（射程＝本腿写面，`.go` 枚数 0＝正控空集）**rc=0**／`sh scripts/d22scan.sh` **clean rc=0**／`sh scripts/check-path-length-budget.sh` **GREEN rc=0**（分母 16084、roster 60＝60）／`go vet ./internal/risk/` **rc=0**（`logs/v1-gates-shared-tree.txt`）。⚠ 具名本腿一枚坏尺＝第一发把 gofumpt 打到 `.md` 判语件上 ⇒ rc=2 `illegal character U+0023`＝尺用错对象、⛔ 读成红⛔ 计入判据。全仓 gofumpt 脏列 50 枚＝存量对照、无一枚在本腿写面、⛔ 顺手重排别人件。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`（本腿零次）。

**AC#5 禁区＝成立。** 尺①（量的，`git log e6fd0f50..HEAD -- <逐枚具名>`）：`blacklist.go`／`pathresolver.go`／`pathresolver_other.go`／336 两件／216 件／三枚冻结件／`internal/observe/thresholds.go`／`scripts/portable-tests.sh`／`ci.yml`／`go.mod`／`go.sum`／`.gitignore`／`docs/PLAN.md`／`frontend`／`design`／`cmd/wisp`／`internal/panel` 共 21 枚＝**被动全 0**。
尺②：三枚生产码件 pre⇄head blob 逐枚同串（`c1799cd2`／`15349e28`／`37c42e94`）。`pathresolver_other.go`＝真乙落点，本腿只指认、突变只打仓外树且逐发还原（哈希三处对上）。

**AC#6 纪律格＝成立带一枚已自纠违例（违例在腿、不在本判语）。** 落码笔 `7437ef88` 只碰声明两枚（尺⑤）；票 339 票面 numstat＝4/0＝**纯追加**；票尾仅一条 bullet；无改名（`--follow` 只 M）、无翻框（7/0 现跑）、`Status:` 行逐字仍"已立，未派"、⛔ push。
它对票 334 **文件**的 +23/−1＝一枚 build tag＋注释＝**在票 339 AC#6 声明写面内、⛔ 算越面**（见 AC#3 段裁语）。
违例＝`rm -f` 三处删自造未入库中转件（规则 8 临时件只建不删）＝腿已具名自纠＋专笔更正第一稿不实句（`72c186c3`）；本腿抽查决定性日志逐枚在库非 0 字节＝**登记⛔ 退件**，落 `A##` 账归编排者。
本腿同样具名一枚自己的坏尺：`tripwire-sensitivity.sh` 的 `new_anchor_hits_after` 第一 grep 把替换文本写成带 `\/` 的字面串⇒打出 0；同发兜底 `grep -cF 'MUT339V1'`＝1 与"旧锚 1→0"配对成立＝读数有效、坏行⛔ 当凭据。

## 必答题

**(i) AC#2 ＝ 部分成立。** 凭据逐字见 AC#2 格（三发突变＋两发正控＋四数逐枚、pristine/还原哈希双验）。⛔ 成立＝绊线等价满足原句（原句要"两平台都点名红"＝盘上不可满足）；⛔ 不成立（绊线复打会红＝非装饰）。
**(ii) 绊线件是活的。** 本腿重打红句逐字：`anchor_miss_reason_339_other_test.go:140: tripwire: this platform now resolves a candidate root (…) - pathresolver_other.go has grown a real resolver…` ⇄ `:143: tripwire: this platform now resolves an anchor root (…)`；`--- FAIL: TestAnchorMissHasNoWorldOnThisPlatform339`、`RC=1`；还原哈希 `R=37c42e9471101c8eedacfad6b7c539a5d4bc4c46 match=YES`（＝起手 blob 同串）。复现＝成立，不退格。
**(iii) 收窄句**＝AC#3 格那条 blockquote（可直接粘贴，编排者落笔；本腿⛔ 动票 334）。
**(iv) 真乙该立票＝该。** 票面形状（立票权＝编排者，本腿⛔ 落笔）：
射程＝`internal/risk/pathresolver_other.go`（`resolveHandle` 落 realpath＋lstat symlink reject＝SPEC-06 第 4 节原文，⚠ 与票 55 `S8 macOS port` 的 `:28` 那枚 AC 同体——建议**并入或显式引用票 55**，⛔ 两票同面并发）＋`internal/risk` 三枚无档/`!windows` 测试件搬牙。
判据候选 AC#0＝本腿量死的分母（量的，`v1-posix-resolver-lands-preview.log`：stub→realpath 最小形后 Linux 整包 **196/112/5/0 rc=1**，红逐名 5＝`TestFallbackBranchHasNoWorldOnThisPlatform336`(`:76`)／`TestAnchorMissHasNoWorldOnThisPlatform339`(`:140/:143`)／`TestSyncNoGradeIsConfirmedOnPosix`(`syncdirs_other_test.go:64`)／`TestSyncMembershipDecidesOffProfilePosix`(`:95/:100`)／C26＝在册存量）；SKIP 那枚 `TestC26RewrittenSyncRootDoesNotDisarmSuspectNet` 翻 PASS＝守卫自解。
⛔ 顶红票 336 那 8 枚牙？＝不顶（牙在 `windows` 档、Windows 台面绿，本腿落点B 复打到它们点名红＝活着）；新票要**搬**不是删：336/339 两枚绊线红＝设计好的"搬我"信号，`name_claim_two_arm_216_test.go:89` 引用的 `candidateOffAnchorChain334` 合件时必须保名（339_other `:36-43` 注释已具名该合并条件）。
★一枚⛔ 可承诺的事（本腿算的、带盘上疑点）：POSIX 上 ESC/U+202E 是**合法文件名字符**⇒ 334 那形靠"Windows 拒绝解析"造 partial 的机制在新台面未必复现——新票 AC#1 必须先量"POSIX 世界点不点得燃 A"，⛔ 假定"搬牙即绿"。
**(v) 同族普查（internal/risk 之外）**：真形尺（root 依赖）在全仓命中＝代码面 **0 枚**（`formsOf/anchorForms/uncertainAnchorMiss` 在 `internal/risk` 之外只有 `.scratch` 探针日志抄本＝非编译件，尺丁逐枚点名）；
环境变量形普查＝`internal/tools/**` 五枚无档件设 Windows-only env（`fs_test.go:159`／`grant_test.go:351`／`refusal_20r2_test.go:143`／`ticket90_test.go:388`＋另两枚带 `windows` 后缀）＝**每枚都同发设 `HOME`**（量的，`v1-census-windows-env.txt` 正控 A：USERPROFILE 在场 23 件）⇒ 无"假主目录只有 Windows 有世界"残留。生产码侧 Windows-only env 读取（尺丙）＝0 枚 `Getenv` 形命中（解析器经 `os.UserHomeDir` 一族，⛔ 本票射程）。
**(vi) 编排者行 39 的机制更正＝认（本腿三锚全部现跑坐实，⛔ 读散文）：**
① `internal/risk/pathresolver_other.go:10` 逐字 `func resolveHandle(p string) (string, bool) { return "", false }`＋文件首行 `//go:build !windows`（现读）；
② `formsOf`（`pathresolver.go:345`）里 `f.root` 唯一赋值点在 `:354-357` 的 `if final, ok := resolveHandle(cut); ok` 分支内（`git grep '\.root = '`＝包内非测试赋值仅此一处，量的）⇒ POSIX 上 root 恒空；本腿探针实测 `resolveHandle(真存在的目录)`＝`("",false)`、`formsOf(真目录)`＝`real="" root=""`（`v1-probe-posix-mechanism.log` 逐字）；
③ `blacklist.go:359-361` 首行守卫现读 `if cand.root == "" || anchor.root == "" { return false }`；探针里 `guard-expr=true`＝返回发生在首行（若有人只删守卫不落地解析器，`isUnder("","")` 恒真会把 POSIX 反向全量 fail-closed 成 A＝⚠ 新票须钉"⛔ 单独摘守卫"，此句为更正的补强、非推翻）。
④ 票面第 9 行"USERPROFILE 在 Linux 不生效"＝作废复现：e6fd0f50 的 334 件 `:85-88` 同发设四枚 env（含 `HOME`，量的）＋本腿探针 `fakeHomeEffective=true`。⇒ 编排者更正**成立，无需追加更正**；腿的 :86 与行 39 的对读没有"偷懒"。

## 交件自证（三把尺）
① 本腿逐笔 `git show --name-only` ⊆ 三处写面（`docs/evidence/s1/339-v1-verdict.md`＋`.scratch/wisp/probes/339/v1/**`＋票 339 尾一条只追加）；`.go` 命中＝**0**
（正控＝同一把尺打 `339-r1` 的 10 笔并集⇒`.go` 命中 **2**，起手侧尺①现跑）。② 本件绝对路径＋`wc -c`＋行数＋占位词尺命中 0（正控＝同尺打 `docs/reports/HANDOVER.md` 非 0）——现跑读数见 `logs/v1-final-selfcheck.txt`。
③ 决定性逐字（起手侧⇄交件侧、带 `file:line`）＝红句／哈希／四数／门禁 rc 全在本文各格＋`logs/**`。票框终量尺现跑＝未勾 7／已勾 0＝本腿一枚没翻；⛔ push。
