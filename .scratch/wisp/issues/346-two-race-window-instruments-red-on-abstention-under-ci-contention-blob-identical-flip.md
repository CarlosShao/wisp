# 346 — 两枚"竞态窗口仪器"在 CI 上**blob 一字未动却 PASS⇄FAIL 翻脸**：它们按设计"到不了窗口就自陈⛔ 作数"⇒ 把**弃权**报成**红**，于是整仓 CI 随机红、后续腿的归因全被污染

- Status: **已立，未派**（2026-10-18 21:3x 编排者立；来源＝编排者对两枚 CI run 的**逐名红名册作差**＋**blob 同一性尺**；账＝台账 `A1041`）
- 类型：**仪器稳性／台面归因**（⛔ 判据改动、⛔ 放宽断言）。射程＝`internal/winsec/openleaf_322_other_test.go` ⇄ `internal/winsec/posix_seal_window_120_other_test.go` 的**弃权那一支**
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 这是什么（人话）

票 322／票 120 各留了一枚仪器：它们要证明"叶子被换成链接那一**瞬间**也被拦住"。证明的办法是真并发地去抢那一发。可一旦抢不到（对方先落地、或者更早那道检查就把请求挡掉了），窗口**根本没到达过**——这时仪器**故意红**，红句逐字写着"这一发的零读数⛔ 是证据"。
⇒ 这个设计**方向是对的**（本仓第 8 条那一族：一把量不到东西的尺不许当凭据）。但它今天有一个真实后果：**CI 随机红**。红的是"今天这台机器太快／太忙、没抢到那一发"，⛔ 是任何一条腿把代码改坏了。
⇒ 后果链：CI 红 ⇒ 每一程都要重新做一次"这枚红是谁带进来的"的归因（今天我已经做了两遍，见「现量」）⇒ 而归因吃的是编排者与裁腿的轮次；更糟的是**下一程很可能归错**（把它们读成"某枚 winsec 字节新造的红"）。

## 现量（编排者 2026-10-18 20:3x–20:5x 自己跑；⛔ 抄腿、引用前先重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| ★两枚仪器**字节完全没变**却换了颜色 | `openleaf_322_other_test.go` blob＝`dfaddac8…`、`posix_seal_window_120_other_test.go` blob＝`0b0e5173…`，在 `a42b35af`⇄`a5080612` 两枚 CI 锚上**逐枚同串**（`git rev-parse <rev>:<path>`） | 两锚 × 两文件＝4 次 `rev-parse` 作对撞；★同时 `winsec.go`／`winsec_other.go` 两锚也**同串**＝⛔ 产码动过 |
| 翻脸事实 | run **39101060753**（锚 `a42b35af`）＝两枚 `--- PASS`；run **39124122121**（锚 `a5080612`）＝两枚 `--- FAIL`；更早的 run 39071959151（锚 `adec1db1`）＝`TestPOSIX120…` **PASS**（`openleaf` 那时还⛔ 在树里＝`cat-file -e` rc=128） | `gh run view <id> --log-failed` 落 `D:\tmp\wisp\20261018\ci<id>.log` 后逐枚 `Select-String '^--- (PASS\|FAIL)'`（原始件在盘、⛔ 入库＝别家 CI 全量日志） |
| 红句①（弃权自陈） | `openleaf_322_other_test.go:206` 逐字 `AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 6 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached` | 同一份 log |
| 红句②③（另一枚的弃权） | `posix_seal_window_120_other_test.go:370` 逐字 `AC#1 RED instrument is dead shape=leaf-swap: refused=1998, sealed-own=0 - a zero window reading here proves nothing` ＋ `AC#1 shape=leaf-swap: every refusal came from the check before the act (refused-at-act=0), so this run does NOT exercise the pinned identity gate` | 同一份 log |
| 分母（这次 CI 一共红几枚） | 12 枚**唯一名**：`test-core` 9 枚＋`test-windows` 3 枚 | `^--- FAIL` 尺（顶格、⛔ 子用例混数） |
| 其余 10 枚**全部在册** | ①票 339 那两枚已**转绿**（★本票因此⛔ 带"339 欠账"那一维）②`TestC21TableColourRowsMatchTokensCSS`＋`TestC21DesignTokensFourWayAgree`＝票 344（卡机主一拍）③`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`＝`frontend/**` 缺席（归机主侧）④`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`＝存量（`A380`／票 200／票 216）⑤`TestStreamLog{DroppedNamingLedger,FloodBelowKeyBound}IsNotBounded35r8`＝票 35r8 在册⑥`TestTicket223PermissionDeniedSitsInItsOwnSentence`／`TestAC14GoSideEvalPushReachesThePage`＝票 303／305／314 那一族 | 逐名对回票面与台账（⛔ 一枚⛔ 归"新账"） |
| lint 那两步 | `lint`＝`failure`，`gofumpt` 步点名 **20 枚**、★**全在 `.scratch/wisp/probes/**`**，与两枚 CI run **逐名相同**（20⇄20 作差空）＝票 275／276 存量、⛔ 归本票 | 同一份 log 逐名册 |

## ★归因边界（先钉住，⛔ 让任何腿把它读大）

- ⛔ **"CI 随机红"⛔ 等于"该把它们改成恒绿"**。这两枚的弃权那一支是票 322 `AC#1` **逐字要求的**（"⛔ 拿 0 命中当凭据，除非同一发证明窗口真到达过链接态"）。⇒ 任何"把它改成 SKIP／改成绿／把重试次数调到能过就算完"＝**摘牙**＝本仓禁手。
- ⛔ **动 `design/**`／`frontend/**`**（机主外部 agent，本编队零写面）。
- ⛔ 与票 342 混记：342 管"祖先那一支还开着"（**真敞口**）；本票管"叶子那一支的仪器抢不到窗口"（**测量稳性**）。两件事的修法、代价、结案判据都不同。

## 判据（逐枚可机读；★本票第一趟**只许只读普查**，⛔ 未出表⛔ 动码）

- [ ] **AC#0 起手现跑**：上面六行逐枚重跑（两枚 blob 同一性要**你自己**在两枚 CI 锚上对撞；`gh run view <id> --log-failed` 自己现取、⛔ 引我这段散文）。★再加一枚同族尺：**全仓还有几枚"到不了窗口就红"的弃权形仪器**（尺＝逐枚读 `t.Errorf` 前的守卫条件、把"窗口未到达"那一支逐枚点名；命中逐名 `file:line`，0 命中带正控）。⛔ 只交这两枚就走＝票 334／343 那一族"抓到一枚就当同族只有这一枚"的教训。
- [ ] **AC#1 三形代价表（⛔ 未出表任何一形都不许开工、⛔ 选形）**：**甲**＝**台面侧**——让 CI 那一步在**安静档**下跑（例：把这枚仪器单独一次 `-run`、`-count=N` 有界重投，或把它挪到不与整包并发的步骤），代价＝CI 时长＋"什么算够安静"这条判据要写死；**乙**＝**仪器侧**——把抢窗口那一发从"随机竞速"改成**确定性交错**（自己持有 racer 与被测 leg 的时钟／屏障接缝，⛔ 依赖调度运气），代价＝要动 `internal/winsec` 产码的注入点＝⛔ 本票射程、须另票＋人工批准；**丙**＝**判据侧**——弃权**照旧红**但**同屏打出可归因的读数**（每次运行落"attempts／flips／refused-at-check／refused-at-act"那四个数＋一句"这一发弃权＝⛔ 证据"，让 CI 红**自带名册**），代价＝零判据移动、但 CI 仍会随机红。★每形现量：动几枚文件／会不会顶红票 322・120 的既有钉（逐名读断言）／CI 那一步的**分母**变不变。
- [ ] **AC#2 ★⛔ 用"变绿"换"看不见"**：⛔ 放宽或删掉那三句弃权红句、⛔ `t.Skip`（无具名登记面的裸 Skip＝本仓禁手）、⛔ 把重试拉到"总能过"、⛔ 给 CI 加 `continue-on-error`／撤步骤。⇒ 判据＝**交件侧那三句红句逐字仍在**（只在"真到不了窗口"时响）。★自证＝去牙反形：把弃权守卫摘掉的那一发（**只打仓外导出树物理字节**、⛔ `-overlay`、每发前后各验 pristine 哈希）⇒ 仪器必须从"弃权红"变成"**假绿**"，且这一维**能被点名**（⛔ 一把两种形状都放行的尺）。
- [ ] **AC#3 归因工具（编排者侧欠账，本票要把它结掉）**：给"CI 红名册逐名对回在册票"这件事留一枚**可重跑的台件**（输入＝两枚 CI run id，输出＝逐名 × 〔在册票号／未在册〕），放 `.scratch/wisp/probes/346/a1/`；⛔ 新建 CI 步骤（那属 `.github/**`＝人工批准）、⛔ 把它做成常驻断言。
- [ ] **AC#4 颜色与门禁**：`./internal/winsec/` 两平台成对（Linux 容器 ⇄ 本机 Windows 档；同一 `-run` 名册、`-count=1 -timeout 30m`、四数逐枚点名⛔ 混用⛔ 相减）＋★**两枚弃权仪器各自至少 3 发独立复跑**（同一台面、逐名点名红⇄绿⇄弃权，散布如实报）；红名册作差逐名归位。门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR`（票 311 AC#1 shape B GUARD，带 env 时 rc=2 射程拒绝⛔ 读成红）。⛔ `./cmd/wisp` 整包（车道按波次排）、⛔ `go test ./...`。
- [ ] **AC#5 禁区**：⛔ 动 `winsec.go`／`winsec_other.go`／`winsec_windows.go` 任何字节（本票第一趟是普查；真乙形另票＋人工批准）；⛔ 动票 113／116／119／120／**322** 的既有钉去"让它别响"；三枚冻结件／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`.github/**`／`ci.yml`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ `frontend/**`／`design/**`；★⛔ 把票 322 的 `AC#1` 判据文字回头改（它已 `-done`，要更正只追加）。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票六格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`.scratch/wisp/probes/346/a1/**`（短名、⛔ U+00A7 那枚字符、逐枚非 0 字节、⛔ `.out` 扩展名）＋票 346 末尾**一条**只追加 bullet；`add` 与 `commit` 两处 pathspec 逐枚手写字面路径（⛔ 从共享 index 反推＝`A1027`／`A1029`／`A1040` 已咬三枚腿）；中文正文⛔ 走 PowerShell 内联（`A1030`）；⛔ push。

## 为什么现在不派

`internal/winsec/**` 此刻由本波只读普查腿 `342-a1` 持有（读色＋出表），⛔ 同刻两枚腿在同一包抢读数；★更硬的＝本票的**甲／丙两形都要动 CI 那一步或仪器本身**，而票 342 的三形表会**改变同一批仪器将来该怎么写**⇒ 排 `342-a1` 交件之后，并与票 342 的落地腿一起排写面。

**撤销口令**：「撤 346」（⇒ 本票退回未派；两枚仪器一字未动、CI 颜色不受影响）。

## Progress log (append-only, newest last)

- [2026-10-18 21:3x +08] 编排者立票：两枚 CI run 逐名红名册作差（12 枚唯一名＝10 枚在册＋本票 2 枚）＋blob 同一性尺（两锚四枚路径逐串同＝⛔ 字节带起来的新红）＋红句逐字进票面。★选形⛔ 在我（本票第一趟只许只读普查）。
- [2026-10-19 22:2x +08] 普查腿 `346-a1` 第一趟交件（★只读、⛔ 本腿 0 枚 `go test`、⛔ 动产码／⛔ 翻框／⛔ push）：普查件＝`docs/evidence/s1/346-a1-census.md`＋`.scratch/wisp/probes/346/a1/logs/00..04`（6 枚正文件，逐枚非 0 字节、CR=0、U+00A7 命中 0、`.out` 0 枚）。★四件必答的落点＝(i) 两枚弃权形状拆开列表（`logs/01`⇄`logs/02`，每格带"哪一发日志＋字节区间"，无凭据的格逐枚写"未知"）；(ii) 结论＝**三味里对本票两枚只有两味成立**（仪器形状为主、争用只钉住"同 run 同色"半条、负载维＝未知；**临时根拼法那味⛔ 属本票射程**＝两枚 `//go:build !windows`、造拼法红的那味在 Linux 档⛔ 存在）＋只给形状＝新增常驻件 `scripts/abstention-census.sh`（判据＝逐枚仪器打"到达数"三值判定 `REACHED`⇄`ABSTAINED-RED`⇄`ABSTAINED-HIDDEN`，弃权名册现量 **12 处／8 个文件**逐名 `file:line` 已列，★⛔ 只交票面那两枚）；(iii) 续票判定＝甲归本票 `AC#1` 丙形⛔ 新立、乙归票 311＋`227-v2` Ⓒ 更正＋本腿 (iv) 配方⛔ 新立（★现量：票 311 全文⛔ 覆盖"拼法"这一维＝尺 0 命中）、丙**该另立一枚 `cmd/wisp` 侧摆动票**（⛔ 塞进本票）；(iv) 台面配方原文在普查件 (iv) 节，含★诚实边界（①那味对 `internal/risk` 4 枚在 CI 可证有效、对 `internal/tools` 5 枚⛔ 在 CI Windows 档射程内＝未知；两枚仪器⛔ 有任何配方能让它们"不弃权"＝弃权那一支是票 322 `AC#1` 逐字要求的）。★本腿抓到两处**票面／编队前提要收窄**（⛔ 本腿改票面原句，只追加此段）：① blob 尺三锚现量＝两枚仪器 `dfaddac8…`/`0b0e5173…` **三锚同串**成立，但被测产码 `winsec.go` 在 `a5080612`→`6f808a83` 之间**变过**（`67783ae0`→`07300bff`）⇒ 票面那句"⛔ 字节带起来的新红"只可成立于**仪器自身两枚文件**；② 本机那一发 `orch227-b/full.log` 最后一枚包结果行＝L10587 `FAIL … internal/winsec`、其后⛔ 包行 ⇒ ⛔ 整包终态名册（与编排者 `0ca0966a` 那句独立同向）⇒ 本票射程只认读码那条 tag 尺。

- [2026-10-19 22:5x +08] **编排者收件 `346-a1`（票 346 第一趟只读普查；件 `docs/evidence/s1/346-a1-census.md` 103 行／13664 字节；★零枚 `go test`＝凭盘上⇄CI 日志）⇒⛔ 翻框（本票 7 格一枚未动，判语归后续非实现者裁）**：三把尺我自己跑＝名册两把锚 6 笔逐枚祖先 Y、10 枚次越界 0／`.go` 0（它自称零产码成立）；件尺＝普查件＋散件 5 枚逐枚非 0、CR=0、裸码位 0、`.out` 0；框尺终量 open=7 x=0、票尾 numstat 1/0 纯追加。★它三处"要你裁的前提"我都复跑：⑴ **两枚仪器首行＝`//go:build !windows`**（尺＝`git show HEAD:internal/winsec/posix_seal_window_120_other_test.go`⇄`openleaf_322_other_test.go` 各取首行逐字相符）⇒ 本机与 CI 的 Windows 档**⛔ 编译＝既⛔ 绿也⛔ 弃权**＝★我 `A1053` 把它们记成"摆动名册"那句**口径不完整**（弃权那一支是票 322 `AC#1` 逐字要求的＝任何台面配方都消⛔ 了）；⑵ **(d) 它独立同向我"整包⛔ 跑到包尾"**（用我 log 尾两枚包行作差自证、⛔ 采信我的报＝这形我要给它记功）；⑶ **(c) 它说"票 311 全文⛔ 覆盖拼法"＝两把尺射程不同**：我 grep 全文命中 3 行（`Status:`＋两条 Progress log 散文）、它的尺问的是**判据格**＝我另跑 `grep -cE '^- \[[ x]\].*(cygpath|正斜杠|反斜杠|拼法)'`＝**0** ⇒ ★它那句在判据层成立、我那句在全文层成立 ⇒ 裁＝"拼法"那味的落点＝**公共块（已进块 k 正文，grep 自证在场）＋本票尾这条只具名 bullet**，⛔ 起新票（同它 (iii) 乙）。★它交回的 (iv) 台面配方**改宽了我发给全编队的口径**三条：① `GOTMPDIR`⇄`TMPDIR` 一律 `cygpath -w` **反斜杠形**（我原说"短临时根"＝⛔ 够；正斜杠 ⇒ `mklink /J` 的 target 被 cmd 吃掉＝`Invalid switch`＝一片假红，与长短⇄卷无关＝与 `227-v2` 的Ⓒ 独立同向）②整条 `PATH` 用 `cygpath -wp` 一次转（⛔ 半冒号半分号＝"零枚 `--- FAIL` 的整包死"）③★★**判据只认到达数**、⛔ 拿"总数大"当到达——实测**红发的总数 1998 ＞ 绿发的 1768＝方向与直觉相反**。★它另摊出同族 12 处／8 个文件的弃权形仪器＋两处 `t.Skipf` 形弃权＝第三种形状，都⛔ 本票射程、具名给我做了下一票的料。★(iii) 丙那条"该另立一枚 `cmd/wisp` 侧摆动票"我照办＝**已立票 351**（`TestTicket223ModeLoosening…` 同码三处两绿一红、判据族⛔ 同本票、未结案池 grep 该名 0 命中）。⛔ push。
