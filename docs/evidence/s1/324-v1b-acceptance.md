# 票 324 — 非实现者验收表（腿 `324-v1b`，六格逐格判语）

- 腿＝`324-v1b`（**重派**）｜角色＝裁决者，⛔ 实现者（`AGENTS.md` §0 第 3 条／§1）｜写面＝本件＋`.scratch/wisp/probes/324/v1b/**`＋票尾 1 条 bullet
- 起手锚＝`136606724`（放腿那一刻 `git rev-parse HEAD` 现读）｜终态补发＝`8baec646`（这窗口里别腿落了 4 笔：317-v1 读数／322-r1 遗产代提（winsec 产码 `+164 −7`）／A1000 台账／我的读数第一批）
- 实现笔三枚逐枚 `git log -1` 现跑命中：`9c5916e8b`／`e038304b3`／`604814414`；`git rev-parse 9c5916e8b^`＝`aa1427f63`＝改前那发
- ⛔ **前一枚 `324-v1` 的任何"判语"被当凭据**（它 49 秒交回、正文只有 `8880.`、盘上三把尺全空＝假交件、零产出）。
  本表**全部六格都是本腿自己跑的**；读数件＝`probes/324/v1b/logs/00・01・02・03`（00＝成对树与红名册作差，01＝突变台＋AC#2/#4/#5，02＝仪器原始打印与哈希台账，03＝突变脚本与探针件源逐字副本）
- ⚠ 派单那句"写面 1)＝`.scratch/wisp/probes/324/v1/`"与机主直给本腿的"`v1b/`"冲突 ⇒ 本腿落 **`v1b/`**（前一枚的 `v1/` 目录根本不存在＝它没建过，具名登记，⛔ 自行改规矩）

## 台面一句话（每张数都靠它）

三棵**仓外** `git archive` 树（pre=`aa1427f63`／post=`9c5916e8b`／head=`136606724`，先 `find -name '*.go' | wc -l` 证非空＝1073/1073/1087），
`go test -c` 编出的 test exe 与三枚 `third_party/sherpa-onnx/*.dll` 同目录、cwd＝`<tree>/cmd/wisp`、`GOTMPDIR` 先 `mkdir -p`；
全程**定向 `-test.run`**、⛔ 整包（票 314 死亡区）；12 发里**零枚** `0xc0000135`、零枚 `RUN=0`。
pre↔post 在被扫四目录里的差集＝**只有 2 枚件**（尺＋名册）⇒ 两棵树的被扫产码语料逐字节相同。

## 六格判语

| 格 | 判语 | 一句凭据（哪把尺＋射程＋量的还是算的） |
|---|---|---|
| **AC#0 现量复跑** | **成立** | 四枚本腿自己跑：判定入口行号（尺=`grep -n`，工作树）`scanReadFiles:333`／`readsFieldAccess:421`／`configSectionMembers:465`／14 枚那道门`:267`；★同名假命中真身在 `internal/panel/composer_dispatch.go:234`（`:231` 那行 `.Models ` 后是空格、尺⛔ 命中）⇒ **票面行号差 3、实现者的具名更正复认成立**；★真读者 `cmd/wisp/panel_models_320.go:85`（取自 `:81 h.mgr.Config()`）与票面逐字一致；名册两行 `:248`/`:259` ＋ `ModelsSection` 导出成员 4 枚／`LLMSection` 5 枚（⛔ 抄实现者计数，本腿读 `internal/config/schema.go` 现量）。改前那发红句两句**逐字复现**（`rc=1`、`config_receipt_255_test.go:271`、`--- FAIL` 顶层 1 枚），两笔**各自独立具名**（缺真读者＝roster 少 1 枚文件；假命中＝scan 多 1 枚文件） |
| **AC#1 正反两形都有牙**（本票核心） | **成立** | ⛔ `-overlay`、**物理改写＋sha256 还原**（基线 `bf02e902…0e106`＝blob＝导出树＝工作树四处同串，且＝实现者那枚基线 ⇒ 两枚腿的还原尺可对接）。**正控**＝摘掉"那一跳必须是段自己声明的导出字段"那道门 ⇒ `rc=1`、**3 枚顶层用例红**（`TestTicket255RosterStillMatchesTheActualReadSites`／`TestTicket324RosterScanRefusesSameNamedNonConfigFieldSite`／`TestTicket324SectionMembersAreReadOffTheSchema`，红句逐字落 `probes/324/v1b/logs/02`，与实现者 `mutant-A-ruler.txt` 同名同文案）；**反形**＝把成员读数摘死（恒空）⇒ 同一批 3 枚红，**外加地板 14 枚逐枚点名＋7 枚段名册红（agent/audio/hotkey/llm/models/panel/voice）**；两发各自 `restore` 后 `sha_after` 复现基线（`match=YES`）、重编重跑 `rc=0`、14/14 `match=true`。**★本腿自己造的坏世界**（⛔ 用他的夹具件）：三枚新探针件物理落进尺的射程（`tools/probe324/`，仓外树）——P1 同名非配置字段那一形在 HEAD 判定下**没被算成读者**（`[models] match=true`、scan 只 `cmd/wisp/models.go`），摘掉那一跳后**才算**（scan 多出 P1）⇒ 那一跳不是装饰；P3（段上方法／非声明名）**被漏**＝误拒那一维实测会红成"名册不点名也没事"的静默形、P2（镜像真成员名 `CloudASRChain`）**被算成读者**＝误认那一维复现（活例＝`internal/panel/config_handlers.go:609-610`，票 313 明文收着的读者） |
| **AC#2 补 `llm` 那行的真读者** | **成立** | 名册行＝`cmd/wisp/config_readers_255.go:259`（第 2 枚＝`cmd/wisp/panel_models_320.go`）；绿句（HEAD 树）＝`SECTION 7/14 section="llm" … match=true`。`hotClaimConsumed` 那一档⛔ 动（blob 哈希相同，见 AC#4），档由**复认成立**：`:81` 每请求现读活 Manager、与同行已裁的 `panel_config_store.go` 同类，⛔ `snapshot-only`（它没读 boot 拷贝）、⛔ `other-process`（`wisp` 进程内被面板那一跳调用）。★**生产调用者那把尺**（`grep -rn` 六个符号、⛔ `_test.go`、**显式剥定义行**——定义行逐枚点名：`:42`/`:48` 两枚 const、`:53`/`:62`/`:67` 三枚 type、`:73` 方法定义、`:128` 函数定义、`:37`/`:44`/`:50`/`:58`/`:126` 注释行）：定义件之外的生产点＝**1 处构造**（`cmd/wisp/panel_inbound.go:286`）＋**1 处调用**（`internal/panel/composer_dispatch.go:234`，路由名 `bridge.go:75`，`case` 在 `:230`）＋装配链两枚根（`panel_inbound.go:223`→`:147` 命令、`panel_resident_windows.go:228`→`:255` 常驻宿主）。⇒ 接线上身有凭据；`go build`／`go vet` 的 rc=0 **没被当成接线证据**（Go 不为未使用方法报错） |
| **AC#3 14 枚逐枚复跑＋新增红＝0** | **成立** | 起手红名册（pre 树，**量的**）＝`{llm, models}` 两枚段 ⇒ 终态（post 树与 HEAD 树两发各自 `rc=0`、`^--- PASS`（只数顶层）＝3 枚、`SECTION` 行数＝**14**、`match=false`＝**0 枚**）⇒ 集合差：**新增红＝0 枚**（⛔ 剥耗时后作差，本腿比的是段名与文件集）。两枚已知红的去处＝**都"具名已修"**：`llm`→名册补一行（scan 那一侧本来就点名 `panel_models_320.go`，红句自证）；`models`→判定收窄（scan 退出 `internal/panel/composer_dispatch.go`）。⛔ 拿实现者的"cmd/ 终态 0 行"当本腿的成对读数：本腿**自己去掉了那一发的对照**——用 `dropgate` 突变在 HEAD 语料上把**旧形判定的 14 枚 scan 全名册**打出来（`logs/02` 的 `mutA-3`），与新判定逐枚作差＝**只 1 枚文件退出**（`[models]` 少 `composer_dispatch.go`），其余 13 枚 scan 完全相同；`[llm]` 在两形下 scan 相同（⇒ 那一枚红⛔ 是尺改出来的、是名册缺行）。终态另在最新 HEAD `8baec646`（含别腿 winsec 产码代提）补跑一发＝同色同数 ⇒ 新增红仍 0 |
| **AC#4 名册件没被写脏** | **成立** | `rosterScanSkipFiles`＝**一枚没加**（`git show aa1427f63:…` 与工作树两块逐字相同，只有 `config_readers_255.go: true`；且在尺的 diff 里命中 0 行）；`hotRowClaims` 整块 blob↔工作树 **sha256 相同**＝`527e3f7f…5b23ef` ⇒ 判定句一字未扩写（`git diff` 里 `hotRowClaims` 命中行数＝0）；名册件 diff（尺＝`git diff --numstat aa1427f63 9c5916e8`）＝**`+11 −1`**（唯一删除行＝旧 `"LLM"` 行，11 枚新增＝1 枚新行＋上方 10 行注释）。⚠ 具名过期数：实现者证据件 §2 那句"名册件 diff＝＋10 −1"与现量差 **1 枚新增行**（他记的应是净增）；尺件那句"＋203 行"＝`git show --stat` 口径、本腿 `--numstat` 口径＝**`+199 −4`**（两把都对、并列⛔ 互替）⇒ ⛔ 影响判语，只登记 |
| **AC#5 门禁＋面核** | **成立** | 五件本腿自跑自落 `rc=`：`go build ./...` **0**／`go vet ./cmd/wisp/` **0**／`sh scripts/d22scan.sh` **0 clean**（examined 280 枚生产 Go；ban #8 `internal/`＝**571 枚**、`cmd/`＝138 枚 ⇒ 分母比起手 568 涨＝别腿新测试件，⛔ 当红⛔ 归我）／`sh scripts/check-path-length-budget.sh --with-self-test`（**只在仓根**）**0 VERDICT GREEN**（tracked **13470**／over-budget 60／rostered 60／not-in-roster **0**；longest 185 字符＝票 252 那枚件名）／定向 `-test.run` **0**（3 枚顶层 PASS、14 行 SECTION）。格式化两把⛔ 互换且**列出枚数＋逐枚件名**：★**更正派单**——`gofumpt` 本机**在**（`GOBIN` 现量 `dlv.exe`/`gofumpt.exe`/`staticcheck.exe`，绝对路径跑通；派单那句 rc=127 已过期）；`gofmt -l cmd internal tools scripts` 工作树＝**5 枚**（`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`），同一把尺跑**导出树**＝**1 枚**（只剩 `174r7` 那枚）⇒ 差的 4 枚＝`core.autocrlf=true` 的工作树渲染（尺＝python 数字节：`models.go` 工作树 14213 字节／334 枚 CRLF，导出与 blob 都 13879／0 CRLF，且 `git status -- cmd/wisp/models.go` 空＝内容未改）＝A996 那族结账假枚；`gofumpt -l` 导出树四目录＝**2 枚**（`cmd/wisp/panel_transport_35r2_test.go`＋`174r7`）＝A996 已登记的存量旧账（退码仍 0，门不红）。★**本票写面那两枚件**（尺＋名册）在 `gofmt -l` 与 `gofumpt -l` 下、工作树与导出树**四发都列出 0 枚**。面核＝`git status --porcelain -- cmd/wisp` 起手与终态都 **0 行**（全程⛔ 写产码面）；`-- internal cmd` 起手 **7 行**（winsec 6＋session 1）、终态 **1 行**（session 1＝166-r2 在飞）⇒ 差的 6 行＝编排者代提 `8baec646` 把 winsec 收进 HEAD，⛔ 我的写面⛔ 我的字节 |

## 机主点名留给本腿的那一枚必答（A996 欠的）

**问**：324-r1 在 255 那把尺里加严的两处（成员数地板＋逐枚日志）是不是**只加不减**——它有没有顺手把任何一枚判据从"会红"改成"不红"、⛔ 放宽断言、⛔ 扩大豁免。
**答**：**只加不减＝成立。** 逐枚对拉到改前成对树之后报的是**红名册作差（含文案）**，⛔ 只报枚数：
1. 被删的行＝**4 行**（全名册见 `probes/324/v1b/logs/ruler-diff-minus-lines.txt`）：`if readsFieldAccess(trimmed, field) {`／`func readsFieldAccess(line, field string) bool {`／
   `if next := line[i+len(dotted)]; next >= 'A' && next <= 'Z' {`／`return true` ⇒ 四行**全部属于 `.`Field`.` 那同一支**，新形在那一支上是旧形的**严格子集**（首字母大写那道门还在，后面多挂一道成员名册门）；
   `Config().Field`／`Config.Field` 那**第二支**在 diff 里**零字节改动**（`^-` 行⛔ 含它），且 `blind` 突变那一发实测 `[hotkey]` 的 `scan=[cmd/balldebug/main.go]` **残留**＝第二支⛔ 受成员读数支配 ⇒ 独立复证。
2. **地板**（`:275-279`）是一枚**新增**的 `t.Errorf`（`len(configSectionMembers(t, field)) < 2`，对 `configSectionsForScan()` 的全部 14 枚段一视同仁、⛔ 按段名点名放行、用 `Errorf` ⛔ `Fatalf` ⇒ ⛔ 早退、⛔ 不改 `got`/`want`）⇒ 形状上只会**多**红；它的牙＝`blind` 那一发**14 枚逐枚点名红**（本腿⛔ 触发＝HEAD 两发地板 0 行红，那是今天的 schema 满足它、⛔ 等于它没牙）。
3. **逐枚日志**（`sectionLogged := 0`／`sectionLogged++`／`t.Logf("SECTION %d/%d …")`）＝`t.Logf` 在 Go 仪器里⛔ 置失败位；`match=%v` 调的是与判据**同一个** `slicesEqualSorted(got, want)`（同形复算、⛔ 新口径）；`++` 在 `continue` 那支之前 ⇒ 计数⛔ 少报（实测两发都是满 14 行 `1/14…14/14`）。
4. 那 3 枚邻居用例（`TestTicket255Split…`／`HotRowRosterCovers…`／`RosterEvidenceLines…`）读码⛔ 经过 `scanReadFiles`/`readsFieldAccess`，成对两树上**逐枚同色**（`^--- PASS`＝3 枚 vs 3 枚、两发 `rc=0`）＝第二重保险。
5. 唯一的**语义收窄**＝`readsFieldAccess` 那一支（效果＝scan 少 1 枚文件）：那是**编排者 A996 裁过的形**（乙′，撤销口令「324 改裁甲」在案），属"该不该换形"那一格、⛔ 属"顺手加严的两处"。本腿⛔ 认为乙′ 判错：**⛔ 造不出"哪一枚真字段／哪一行会让它判错"的世界**——本腿造的那枚 P2（镜像 `CloudASRChain`）落在"**过度报告＝会红**"那一侧（尺会喊名册去点名一枚非读者，方向是过严⛔ 过松），而 P3 那一形在生产里今天**零命中**（尺＝`internal/config/**` 非测试件里带接收者的方法 38 枚，接收者只有 `*Manager` 35／`*Config` 2／`*reloadPlan` 1 ⇒ 挂在 14 枚段结构体上的方法 **0 枚**；匿名内嵌字段 **0 枚**，17 枚非测试件）。

## 没裁完的格（具名＋缺哪一发读数）

1. **⛔ AC 框一枚没翻**（框归编排者，本腿只交判语）。
2. **同文件那 8 枚 stdout／开窗族顶层用例（`TestTicket255ReceiptOmits*`／`StillNamesTheLiveReadSection`／`SentenceAssemblyIsFiltered`）＝证据不足，没裁**：
   缺的那一发读数＝**带窗台件的 `wisp run` 实跑**（票 314 死亡区里整包⛔ 能跑，本腿⛔ 跑整包）；它们⛔ 经过 `scanReadFiles` ⇒ 与 AC#1/#3 那两格的凭据⛔ 相交，但票面 AC#5"只 `-run` 定向"这一格**只有本腿⛔ 交的那一发**——**待派给带窗台件的腿**。
3. **甲形（AST＋接收者类型）与 `[voice]` 那枚镜像读者的重裁＝移交编排者**（撤销口令「324 改裁甲」在案；本腿的 P2 读数给出了这一形的实测代价：`[voice]` 会要求名册点名一枚⛔ 配置读者的件）。
4. **`wisp` 进程真发一次 `config.models.list` 的运行时回执＝没裁**（缺窗台件；本腿只交接线上身，见 AC#2 那格最后一句）。
5. **CI 色＝⛔ 本腿的射程**（本腿⛔ push、⛔ 碰 `.github/**`／`ci.yml`，零字节）。

## 台面故障登记（全部⛔ 计入判据、⛔ 删件）

1. `d22scan` 我第一发 **`rc=2`＝我自己带进去的 `GOTMPDIR=D:/tmp/…`** 撞上 `runtests.sh` 的票 311 AC#1 shape B 跨卷守卫（红句逐字"GUARD - the test temp root names ANOTHER VOLUME than %LOCALAPPDATA% … Refusing to run is the fix"）⇒ 那是**射程拒绝**；`unset GOTMPDIR` 后重跑＝`rc=0 clean`。⚠ 同发守卫还点名"winsec 同一件在 D 卷根上读作 `PASS=56 FAIL=2`"＝别腿那面的颜色随卷变，⛔ 归我账。
2. 突变脚本第一版 `open(...)` 未指 `encoding` ⇒ `UnicodeDecodeError: 'gbk' codec can't decode byte 0x8d` 在**读**阶段退出、盘上⛔ 没改 ⇒ 我随后那一发 `rc=0 / 3 枚 PASS` 是**假读数**（`mutA-2-dropgate.log` 留着、具名）。修脚本后 `sha_after` 落上、真红三枚。定式＝**突变台件的脚本自己也要一发正控：跑完必须回读 `sha_after`，⛔ 看测试色就认定突变落了盘。**
3. `printf '----- …'` 开头那个 `--` 被 bash 当选项（`printf: --: invalid option`）⇒ 三枚分隔标记没进件（正文没丢）。
4. 实现者那枚 `d22scan rc=1` 归因（未跟踪的 `internal/tools/subagent_depth_fixture_287_test.go:249 [emoji] ban #8`）**本腿不复现**：该件已 tracked（`e1973191b`）且 blob 与盘上 emoji-band **0 命中**（尺＝Python 逐字符扫 `U+1F000–1FAFF`／`U+2200–22FF`／`U+2600–27BF`／`U+2B00–2BFF`／`U+FE0F`）⇒ 那枚红是 287-r1 在飞态、已由它自己入库前收掉；⛔ 记在 324-r1 账上。
5. 派单那句"gofumpt ⛔ 装（rc=127）"＝**过期台面读数**（22:4x 那次），本腿现量 GOBIN 三枚 exe、绝对路径跑通。
