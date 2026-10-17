# 票 332 非实现者终裁（腿 `332-v1`）＝那张 `Cf` 普查表本身有没有牙

锚与身体（★本腿所有代码读数只取**具名提交树**，⛔ 读工作树在飞字节）：
起手 `git rev-parse --short HEAD`＝**`2a3b9ff8`**（分支 `dev`，那一分钟别腿又落了两笔，故本腿探针面里 HEAD 逐枚现量、⛔ 拿它当锚）；被裁对象＝`332-a1`，其锚＝`7535d83b`。
★那枚边界的身体：`7535d83b` 与 `eb82a4d5` 的 `internal/panel/display_scrub_216.go` 是**同一枚 blob** `9d94b46ec0a52b9d7f3fe8a8d414d8edf09fdbbd`
＝手抄 15 枚那具；`HEAD` 那一枚已是 `230d3443d7105b1856c557aa9579a7fe8ab077d8`＝`333-r1` 的按范畴那具。
⇒ 本票判的是**前者**，件里逐枚写明读的是哪一枚锚（见下面"锚依赖"一节）。
读数件：`logs/ac0-ac1-rulers-and-setdiff.log`（111 行）／`logs/ac2-carrier-citations-at-7535d83b.log`（159 行）／`logs/ac3-replay-teeth-and-gates.log`（110 行），
绝对路径前缀 `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\332\v1\`。

## 六格判语

**`AC#0` 起手四把尺独立复跑＝成立**（四把都自己跑，⛔ 抄普查腿，⛔ 抄编排者）
① 名册＝**15 枚**（量的：`git show 7535d83b:internal/panel/display_scrub_216.go` 的 `:166-170` 一条跨行 `case`，`grep -oE '0x.. '` 去重＝15；
★普查腿那句"`grep -c 'case 0x'` 回 1"本腿复跑同形，确认那把错尺会读成 1）。② `Cf` 全集＝**170 枚**、`GO_UNICODE_VERSION=17.0.0`（量的：本腿自己一枚**仓外纯 stdlib** 探针
`D:\W332V1EXP\cfprobe\main.go`，`go run` rc=0、stderr 0 字节；★执行面交代：本腿唯一执行过的代码＝`cfprobe`＋`diskprobe` 两枚探针＋在**仓外导出树**里跑的用例，
⛔ import 本仓任何包；探针里那四段函数＝`7535d83b` 的 `:106-136/:164-175/:179-199/:211-220` 逐字节抽取，抽取件 sha256 已入库）。
③ 生产调用点＝**1 枚**＝`internal/panel/composer.go:145`，非测试命中共 8 行＝1 枚调用＋7 行同文件（3 注释＋2 定义＋2 处被洗字段赋值）——与普查表一致。
★本腿自抓一把错尺并保留：第一把用了 `grep -v _test.go$`（`git grep` 的行**不以文件名结尾**⇒ 什么都滤不掉）读成 24，改正为 `_test.go:` 后＝8。
④ 判定层＝**6 枚显名＋整个 `unicode.Mn` 族**；`Mn_TOTAL=2059`、`Mn∩Cf=0`（量的，同一枚探针）⇒ 落在 `Cf` 轴上只有那 6 枚。
取向相反的理由⛔ 由本腿推断，逐字引它自己的注释 `internal/risk/taintmatch.go:32-38`（"…lets an exfiltrator cut a tainted fragment into sub-8-rune pieces with
zero-width characters"，出处＝对抗报告 `16.9#1`/`M-6`）＋它的**用点** `:59`（`normalizeTaint` 拿它做**删除**）⇄ 显示轴 `display_scrub_216.go:126-133` 自陈**转义不删**。
★凡与普查表不一致处＝**无**（四把尺全部同值）。

**`AC#1` 那 155 枚一枚不漏＝成立**（量的，⛔ 抽样）
本腿自算的 170 行表 ⇄ 普查腿 `logs/ac1-cf-170-roster-with-unicode-names.log` 的 170 行表：码位集合**同一枚**（`diff` 空）、
`IN/OUT` 与 `WASHED/RAW` 两列**逐枚同值**（`diff` 行数＝0，两侧归一化表 sha256 同＝`f01c05b59b5f4b7da99693c5539eabfe77c5a4779323aebcc343f9c601126411`）；
本腿的 `WASHED` 列⛔ 按名单推，是把 170 枚逐枚送进**逐字节重打**的 `scrubForDisplay` 量出来的（名单说洗 ⇄ 实际换形，两轴在 `Cf` 上 0 枚不一致；
`BEHAVIOUR_MISMATCH=65` 具名归因＝32＋33 枚 `Cc` 走的是另一臂 `unicode.IsControl`，`mismatchOther=[]`）。
分母自算：`CF_TOTAL(170) − ROSTER_SIZE(15) = CF_OUTSIDE_ROSTER(155)`。族划分本腿按它表里的 A..K 区间**重算**：`familySum=155 familiesDistinct=155 LEFTOVER=[] GHOST=[] BOGUS=[]`。
★缺陷 D-1（逐枚点名）：`census-1-cf-155-roster.md:78`（票尾 bullet 同形）写"方向一支＝B＋D＝2 枚"，它自己表里 B=1（`U+061C`）、D=2（`U+200E`、`U+200F`）⇒ **真值 3 枚**；
A..K 求和 155 与那两列**不受影响**。（旁证：K 族那句"95 枚按名各对应一枚 ASCII 字符"本腿复量＝TAG 前缀 95 枚、非前缀 2 枚，同值。）

**`AC#2` 逐枚可达性＝成立（三枚具名缺陷）**
普查表的射程口径本腿核过并背书：它说的是**枚枚走到 `json.Marshal` 出的快照字节**，⛔ 屏幕像素——它**没有**犯票 329 `AC#0` 那处"把进了快照读成进了人眼"的错。
逐枚重读：本腿对件2 每一个 `file:line` 现取 `git show 7535d83b:<file>` 那一行比它声称的字面，**OK=53／MISMATCH=0**。两枚重点：
**(b) 模型正文／工具 reason**——`results[].text`（`composer.go:59`）与 `pending[].tool|args|reason`（`approval.go:77-83`→`pump.go:241-247`）**根本不经过** `scrubForDisplay`：
唯一的生产调用点只交一枚 `Composer`（`composer.go:145`），而被洗字段只有两枚（`display_scrub_216.go:228` `Workspace.Canonical`／`:234` `Attachments[].Name`），
`NewSnapshot` 的 `Pending`/`Results` 按值直出（`composer.go:138`/`:139`，`return Snapshot{` 在 `:137`）⇒ marshal 在 `pump.go:344`/`:359` 拿到的是未洗的那两串。
**(c) 页面／入向**——门名→`dispatchRaw`→`disp.Handle`→`HandleWorkspaceRequest`→`RequestWorkspaceSwitch`→`workspace.go` 的 `Spelling`/`Canonical` 逐环在位，
那一形**今天真的通**（`RequestWorkspaceSwitch` 非测试调用者 1 枚＝`workspace_handlers.go:125`）；出向时 `Canonical` 只洗名单内 15 枚、`Spelling` 一字不洗。
"无已知载体"那一格＝**字段**不是字符（`Attachments[].Name`）：`Ingest` 非测试读者本腿复量 **0** ⇄ 同尺正控 `json.Marshal` **37** 行。
决定性那把"没有一处按字符筛"本腿复跑＝**18 行**（射程＝`internal`＋`cmd`，滤 `_test.go:`，逐行入库；字符集合枚枚只点 ASCII 标点／`\x00\r\n\t`／shell 元字符），
另两把正控同跑：`pathresolver.go` 字符类拒绝 **0** ⇄ 同正则对 `display_scrub_216.go` 命中 **5**；`exec.Command` 非测试 **5** 行、`internal/panel` **0**（`doctor.go:316` 在场＝正控）。
★缺陷 D-2：件2 的 Git 载体写 `internal/panel/pump.go:281`，那一行是 `if p.src.Git != nil {`，赋值 `composer.Git = p.src.Git()` 在 **`:284`**（载体真在、号早三行）。
★缺陷 D-3：§② 那张载体表的**粒度＝族**（11 行），票面 `AC#2` 要的是**逐枚**；逐枚的证据只有 (a) 那一形（155/155 写真名读回）＋四把生产者侧筛检尺（本腿复跑同值：
`scrub_leaves_raw=155 unicode_is_space=0 taint_ignores=0 unicode_is_control=0 json_escapes=0`）；(b)(c) 两形是"族级载体＋那把 18 行无筛尺"的推理，⛔ 逐枚走过 155 次。
★缺陷 D-5：(b) 那一形**还有一条今天真到人眼、而普查腿没列的出向**＝操作员终端：`cmd/wisp/run.go:1075 consoleSink{out: rt.stdout}` → `run.go:1341`/`:1351 fmt.Fprint(c.out, e.Text)`，
一字不筛 ⇒ 它对自己 (b) 的口径（止于字节）是**报窄了**，⛔ 报多了；本腿⛔ 据此改判它的任何一行。

**`AC#3` 两枚真盘正控＝成立**（两枚都自己重放，⛔ 只重放一枚）
`U+2061` 与 `U+206B`：本腿另起一枚仓外探针（`D:\W332V1EXP\diskprobe\main.go`，边界四段＝`7535d83b` 逐字节抽取）在 **D 卷 NTFS** 上写真名 → `os.ReadDir` 读回（各命中 1）→
`needsDisplayWash=false` 走快路返回同一枚串 → `json.Marshal` 字节里 `e281a1`/`e281ab` **原样在位**；两枚各带 `.log` 与 sha256：
`da1e771280e63fd92d26d37464bfc6cfd6611f8f8e7dab4dda7a6767b3af889b`／`bb7a8ae1915fe249edc6b8b6165788e4b7140337f38f038966f90a352ffdd4db`——
★与普查腿那两个数**逐字符同值**（同一输入形状、同一枚边界、另一台目录树＝独立复现）。同发还带三枚对照：名单内 `U+202E` 被洗成可见转义（marshal 里原始 `e280ae` 消失、
sha `158adf0d…c29bb5` 与普查腿同值）、邻位 `U+206A` 与两枚主体同形、`Cc` 的 `U+0085` 被 `IsControl` 臂洗掉；`C0` 的 `ESC` 在 NTFS **根本写不出名**（"volume label syntax is incorrect"）。
跑成了 ⇒ ⛔ 无"哪一行拦的"可报，⛔ 用"缺世界"当解释。

**`AC#4` 台面与纪律＝成立（半格不可盘核，具名）**
盘上可核的四项本腿自己跑：① 该腿自己的 24 枚提交（尺＝主题前缀 `332-a1 `，⛔ 子串匹配）只落 **18 枚唯一路径**，全在它四处写面内，越格命中 **0**、`.go` 命中 **0**；
② 票面 `AC` 框：未勾 **6**／已勾 **0**（本腿与它都一枚没翻），它那笔票尾提交 numstat＝**2/0**（纯追加）且跨其所有提交对票面**删除行数 0**；
③ 票名／`Status:` 未改（`git diff 7535d83b HEAD` 对该票＝36/0，纯追加）；④ 它四枚探针的 import 逐枚看＝全 stdlib（`crypto/sha256`、`encoding/hex`、`encoding/json`、`fmt`、
`os`、`path/filepath`、`sort`、`strings`、`unicode`、`unicode/utf8`），⛔ 本仓包，绝对路径与"唯一执行过的代码"已写在件里。
★具名不可核的那半格：`⛔ go test／go build／go vet 一律不许跑`——盘上只可能留下**跑了**的痕迹，留不下**没跑**的痕迹；本腿看到的是自报⛔ 反证，故这一半当自报读，⛔ 计入凭据。

**`AC#5` 交件自证三把尺＝成立（一枚格式缺陷）**
① 逐枚提交号＋`git show --name-only`：它的名册自称 23 笔／18 枚路径并**声明不含生成名册的那一笔**；本腿改正尺独立复跑＝24 笔／18 枚路径 ⇒ **差一枚＝它自己声明的那笔自指**，两侧同值。
② 件路径＋`wc -c`＋行数＋**空槽尺**（＝票面那把分词形尺；本件⛔ 逐字重打它的模式串，免得自触）：两枚普查件本腿现量＝`9136 字节／79 行`、`15586 字节／137 行`，与它件里写的数**逐枚同值**；分词形空槽尺在两件 **0 命中**，
在证据件命中 **4 行**＝该尺自己的名字与它引用的模式串（它已在 §七 具名自触）。③ `AC#1` 一枚不漏（本腿 `AC#1` 节）、`AC#2` 每行带 `file:line` 或写"无已知载体"（本腿 `AC#2` 节，缺陷 D-2/D-3/D-5 具名）、
`AC#3` 两枚读数（本腿 `AC#3` 节）。★缺陷 D-4：证据件 §二 那张表把**本件自身**的字节数写成一句话（"以最后一笔落盘后现跑 `wc` 为准"）而非数字，
与"件绝对路径＋`wc -c`＋行数"的字面要求不完全对齐（现量＝15989 字节／157 行）。末条"⛔ 结论／建议"：尺＝`grep -cE '建议|应该扩|该扩|不该扩|推荐'` 在 `census-1`=**0**、`census-2`=**0**、证据件=1（那句"本件通篇只有读数与尺"的自我声明）。

## 牙（⛔ 只打在仓外导出树的物理字节上）

导出台＝`git archive eb82a4d5` 八件套＋仓根 `docs/`＋手工三件（`build/*.dll`、`third_party/sherpa-onnx/*.dll`、`frontend/dist`）→ `D:\W332V1EXP\tree`
（D 卷、父目录 8.3 别名 `D:\W332V1~1`；`.go` 枚数＝**753**、`internal/panel`＝**44**，"0"⛔ 是读数）。
成对两发**落点不同、方向相反**：`M1` 把 `0x2061` 悄悄**加进**名册（那正是普查表标 `OUT/RAW` 的一枚），`M2` 把 `0xFEFF` 悄悄**从名册删掉**（那正是它标 `IN/WASHED` 的一枚）。
改前（未突变，`-count=1 -timeout 30m`）：整包 `RUN=255 PASS_TOP=158 FAIL_TOP=7 SKIP=0 rc=1`、scoped 五枚 `RUN=7 PASS_TOP=7 FAIL_TOP=0 SKIP=0 rc=0`；
★临时根成对两发（台账 `A1013` 那形）＝`TMP/TEMP`+`GOTMPDIR` 指 D 卷 ⇄ 不指（C 卷默认）：两发四数**逐枚相同**（255/158/7/0）、失败名集合相同 ⇒ 本腿没有任何读数靠临时根翻身。
那 7 枚既有红具名归**台面**：4 枚要 `frontend`/`design` 全量、2 枚是 `35r8` 活堆仪器、1 枚（`TestReadGitOnThisRepositoryIsSelfConsistent`）要 `.git` 目录——⛔ 一枚是显示边界的用例。
改后：`M1` scoped `7/4/3/0` rc=1，红句逐字含 `display_scrub_327r2_test.go:72: the product washes Cf runes the ruling never named: [U+2061] (ruled=15 product=16)`、
`:91/:95/:98` 三句 `U+2061 is off the ruled list but …`、`display_scrub_216_test.go:448/:451`；整包 `255/155/10/0`，新红恰 3 枚。
`M2` scoped `7/4/3/0` rc=1，红句含 `nail 1: product-washed=14, ruled-list=15`、`a ruled rune is no longer washed by the product: [U+FEFF]`、
`display_scrub_216_test.go:368/:407/:411`；整包 `255/155/10/0`，新红恰 3 枚（与 `M1` 名集合**不同**＝两发互为盲证）。
★还原：两发各在 `finally` 里写回原始字节并逐字节比对（`matches_pristine=True` 两发皆印）；交件侧现量 `git hash-object tree/internal/panel/display_scrub_216.go`＝
`9d94b46ec0a52b9d7f3fe8a8d414d8edf09fdbbd`＝`7535d83b`/`eb82a4d5` 那枚 blob；全部还原后整包再跑＝`255/158/7/0`，失败名集合与改前**相同**。
门禁四数只在共享工作树跑（本腿产码**零字节**）：`gofumpt.exe -l .scratch/wisp/probes/332/v1` rc=**0**；`sh scripts/d22scan.sh` rc=**0**（clean，`internal`=584 枚 Go 文件在射程内）；
`sh scripts/check-path-length-budget.sh` rc=**0**（VERDICT GREEN，over-budget 60／roster 60）；`go vet`＝**不适用**（本腿无自己的 Go 包可 vet；对 `./internal/panel/` 或 `./cmd/wisp/` 跑＝编译 `333-r1`/`228-r3` 在飞字节，按派单⛔，具名跳过⛔ 当红）。

## 两枚必答

**① 那句"已腐烂的读数"今天准确的射程（给尺）**
腐烂的那句＝`327-a1` 的"快照字节今天没有进 WebView 的生产路径／`lastPanelSnapshot` 生产读者 0 枚／最后一公里没接"。
今天（本腿现跑，两个锚都跑了）：非测试命中 **5 行**＝定义 `panel_pump.go:429`＋3 行注释＋**1 枚真读者** `cmd/wisp/panel_snapshot_deliver_329.go:95`（`7535d83b` 同形，枚数未变）。
链＝`run.go:744 Out: rt.bookPanelSnapshot` → `panel_pump.go:447 publishPanelSnapshot` → `:463 pushPanelSnapshotToHost()` → `deliver_329.go:95` → `:105 h.evalPanelJS(...)` →
`:65-66` 用 `:57` 那枚常量 `panelSnapshotGlobal = "window.wispOnPanelSnapshot"` 逐字包裹 marshal 字节 → `panel_host_windows.go:331/:336 w.Dispatch(func(){ w.Eval(js) })`。
★准确的射程＝**marshal 字节今天已进 WebView2 的 `Eval` 调用串；页面接收器仍缺位 ⇒ 像素那一半⛔ 成立**。尺＝
`git grep -c wispOnPanelSnapshot <锚> -- frontend | wc -l` → `7535d83b`＝**0 枚文件**、`HEAD`＝**0 枚文件**；正控＝同一把尺去掉目录限定 `7535d83b`＝**16 行** / `HEAD`＝**56 行**
⇒ 引这把尺**必须带射程**（分母会随记录长）。旁证＝那件自己的注释 `:28-31`（"window.wispOnPanelSnapshot does not exist yet…"）。
本腿另加一条只报读数的：`wisp run` 的终端那支（`run.go:1075`→`:1341`/`:1351`）**不过任何边界**，是 `Cf` 今天真到人眼的一条形——它⛔ 改变上面那把尺的读数，只是把 D-5 摆出来。

**② 本票与票 333 是不是同一件事的两个面**
是**同一枚边界的两个面**，⛔ 两件事：`332` 裁"这张普查表有没有牙"，`333`（立案提交 `ae710c49`）用的是这张表给的料去改那枚函数；本腿⛔ 裁 `333` 对不对（归 `333-v1`）。
**哪几格依赖哪一枚锚的字节**：`AC#0`①、`AC#1` 的 `IN/WASHED` 两列与 155 那个分母、`AC#2` 的"名单外原样穿过"这一**谓词**、`AC#3` 的两枚正控——四者读的都是 `9d94b46e` 那具手抄名册身体；
`AC#0`②（只依赖工具链表版本 `17.0.0`）、③（只依赖 `composer.go:145` 在位）、④（读的是 `internal/risk/taintmatch.go`，`7535d83b`⇄`HEAD` 同一枚 blob `07bc9815`、`333-r1` 一字未碰）、
以及 `AC#4`／`AC#5`（纪律与自证）**都⛔ 依赖那枚名册的形状**。
**如果边界改成按范畴判定**（＝`HEAD` 现在那具 `230d3443`，本腿现读 `isDisplayWashCf` 为 `return unicode.Is(unicode.Cf, r)`，被洗字段仍是那两枚、行号漂移为 `:218`/`:224`），
要**重取凭据形状**的是这几格：`AC#0`①那把"数名册枚数"的尺本身失效 ⇒ 换成"对全标量跑 `isDisplayWashCf` 的洗集枚数"（正控＝`U+2061` 从 `RAW` 变 `WASHED`）；
`AC#1` 的 `OUT` 集合变空集（155→0），那张 170 行表的**存在**还在、**分族与两列**要重打；`AC#2` 三形链的 `file:line` 载体**不必重取**（字段与链没动），
只有"穿过／被洗"那一谓词整列翻转，且 (b) 那一形的 `results`/`pending`/`git.*`/`instructions.files[].path` 仍**不过**边界⇒ 那几行不变；
`AC#3` 两枚正控**必须整发重放**——新的 marshal 字节会带可见转义，本腿这两个 sha（`da1e7712…`／`bb7a8ae1…`）在新身体下**⛔ 可再用作凭据**。

## 交件自证（三把尺）

① 逐枚提交号＋`git show --name-only`（探针三面已入库）：`14d143ab`＝`logs/ac0-ac1-rulers-and-setdiff.log`；`358fae4c`＝`logs/ac2-carrier-citations-at-7535d83b.log`；
`469839a9`＝`logs/ac3-replay-teeth-and-gates.log`；`4c0d20af`＝同两枚件的**更正追加**一笔（AC#2 加三枚尺、AC#3 把那 7 枚既有红的归因改正，原猜测句留在件里⛔ 抹）；
逐枚展开＝`git show --name-only 14d143ab 358fae4c 469839a9 4c0d20af`（尺＝主题前缀 `332-v1`，⛔ 子串匹配）。本腿写面命中＝
`docs/evidence/s1/332-v1-verdict.md`（**本件＝判语件，它自己那一笔的号⛔ 写在自己体内**——尺＝交件时点现跑 `git log -2 --format='%h %s'`，票尾 bullet 紧随其后）＋`.scratch/wisp/probes/332/v1/**`＋票 332 末尾**一条**只追加 bullet；
其它一切路径＝零字节（尺＝下面那行现跑的 `git show --name-only` 全集）。
② 证据件绝对路径＋`wc -c`＋行数＋**空槽尺**：见上面三枚 `logs/` 与本件；那把分词形尺（五枚词形，本件⛔ 逐字写它，模式原文在票面 `AC#5` 与 `332-a1` 证据件 §二）在本腿四面＝**0 命中**（尺不瞎的正控＝同一把尺在普查腿证据件命中 4 行，那 4 行是它自己的尺名与它引用的模式串）。
③ 决定性那把尺的逐字输出：`AC#1` 的对合（`diff_rows=0`＋两侧归一化表同一枚 sha256）与 `AC#3` 的两枚 sha 已在件里逐字入库；
成对两发红句逐字＋还原哈希＋改前改后四数在"牙"一节与其 `.log`。本腿**一枚框没翻**（交件侧现跑：票 332 未勾仍 **6**＝`AC#0`…`AC#5`）、
⛔ 改判据文字／票名／`Status:`／`-done`、⛔ 动任何断言（`internal/**`、`cmd/**`、`thresholds.go`、golden、`allowlist.txt`、`go.mod`/`go.sum`/`.gitignore`/`ci.yml` 零字节）、
⛔ 产码一字（本腿只交 `.md`／`.log`）、⛔ push、⛔ 在共享工作树跑 `go test`/`go build`/`go vet`、⛔ 杀 PID 21892。
