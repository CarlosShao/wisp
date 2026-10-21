# 票 353 — `privateFile` 拒绝路径的删除按【钉目录里的名字】走，而落地句柄带 `FILE_SHARE_DELETE` ⇒ seal 判语之后那枚名字可被腾空＋顶进一枚调用者⛔ 点名的对象＝顶进来的那枚被删（落地的钉锁住祖先、⛔ 锁住名字）

**立票**：2026-10-20 15:1x 编排者（来路＝实现腿 `342-r2` 一次测量确定性量到＋非实现者终裁腿 `342-v3` 独立复现并**签"新立一票"**；判语末节五条判据草案＝`docs/evidence/s1/342-v3-verdict.md`；台账 `A1062`③乙）
**性质**：★守卫造成的**新形状**残差：票 342/348 关死"祖先那一跳"的同时打开了"名字那一跳"。⛔ 属 342/348 判据族（`342-v3` 量到：那两族的分子在这里逐发都是 0＝硬挂进去只会多一枚永红的格）。
**后果（对用户怎么说）**：`wisp` 说"这文件我不动"（拒绝路径）之后，攻击者仍能让程序删掉一枚它⛔ 点名的文件——只要那枚文件在 seal 判语之后被 rename 顶进调用者那个名字里。**跨进程天然拿得到这一跳**（`342-v3` 补打了 `logs/03` 那一发：外部进程握不到钉着的目录、却能腾名顶入）。

## 现量（凭据在裁决腿写面，⛔ 引用前先重跑）
- 机制两枚去牙（`342-v3` 发 A／发 B，仓外导出树、逐字回 pristine `B9DC39AE…`）：摘 `unlinkPinnedLeaf` 那一步 ⇒ damage `3→0`、其余照绿；`openPinnedLeaf` 换按名开 ⇒ `entered 3→0`＝**成对敏感**（⛔ 一把尺两种形状都放行）。
- 落地句柄 `READ|WRITE|DELETE` ⇄ 基线经 `syscall.Open` ⛔ 带 DELETE；删除是 `Root.Remove(base)`＝相对**钉目录里的名字**（`winsec.go` 路由现读：`defer parent.close()` ⇄ `openPinnedLeaf`⇄`unlinkPinnedLeaf` 顺序）。
- liveness 钥匙＝同发 `realSealRanFirst=true`＋`worlds entered=3/3`（⛔ 借 `truncated`／`deleted` 任何一枚分子）。

## 要建什么（★第一趟只读：⛔ 未出代价表任何一形都⛔ 开工）
- [x] **AC#0 三形代价表**（读腿）：修这一跳的三形各写"动哪几枚文件／会不会改到覆盖写⇄独占语义／C26 入口动⛔ 动／`342` 那批已批落地要⛔ 重述"——**甲**＝删除改按"句柄实际站着的对象"走（＝语义变化候选：拒绝路径的删除从此⛔ 再碰顶进来的名字，但同一改动会不会波及 `AC#2` 那条覆盖写承诺要现读）；**乙**＝保留按名删＋删除前做一次"名字⇄句柄对象"同一性复验（多一跳 TOCTOU，要量窗口）；**丙**＝登记为已知残缺＋常驻仪器把 damage 钉成可见红（⛔ 变绿、⛔ 圆场）。★每形带顶红谁现有的钉。
- [x] **AC#1 编排者裁形**（含"这一味⛔⛔ 是语义变化"那一问＝若甲＝第二道产码语义批准、摆机主；`342-v3` 明写归本票）。**裁于 2026-10-20 19:2x＝机主拍〔乙〕**：保留按名删＋删除前同一性复验（⛔ 改删除语义、⛔ 新接缝、⛔ 重述票 348 已批落地形）。甲＝编排者判⛔ 批（三味硬代价经逐处复跑坐实：`ReOpenFile` 在 `$GOROOT/src/syscall`⇄`x/sys` 两把尺都⛔ 导出＝要新写句柄面原语；改掉票 348 已批落地的删除形＝二次改已批产码；POSIX 无按句柄删的单发原语⇒`AC#3` 成对判据得改口径）。丙＝⛔ 采（零产码代价换一枚常驻红顶 CI 合并道）。⇒ 机主拍同时已认下两味残余＝复验拒了仍留着调用者那枚名字（攻击态可观测行为变化）⇄"复验之后、绑定之前"窗内那一支⛔ 拦（要钩进 `Deleteat` 的 `NtOpenFile` 之内＝平台内部），故 `AC#2` 判据必须逐字承认测不到那支（先例形＝票 322 `AC#1`）。凭据＝`docs/evidence/s1/353-a1-census.md` 第 2／5 节＋台账 `A1066`④＋派单存档 `1431a63b`。
- [x] **AC#2 确定性仪器落地**（写腿，形状按裁决腿草案）：⛔ racer；成对＝落地 `damage>0` ⇄ 基线 `entered=0`（世界进⛔ 去＝具名报、⛔ 冒充绿）；liveness 同发在场；敏感度成对去牙两发各自红；**⛔ 要求"祖先换手"作前提**（那正是它不属于 342/348 的原因）。
      **编排者翻框（2026-10-21 00:2x +08，凭据经我逐处独立复跑）＝★本格那句成对口径按〔乙〕具名更正后再翻**：起草早于 `AC#1` 裁形，
      字面写的"落地 `damage>0` ⇄ 基线 `entered=0`"与机主拍〔乙〕后的真实形**冲突**——现量＝落地侧 `entered=3/3 且 DAMAGE=0`、
      基线侧 `entered=3/3 且 DAMAGE=3`（⛔ 是 `entered=0`）。更正后的判据＝**落地 `damage=0` ⇄ 基线 `damage>0`，两侧 `entered` 同数（3/3）＝
      同一对字节进过同一扇门、只有删除那一步被同一性复验拦住**；基线侧要的是"世界进去了但被删坏了"，⛔ "世界没进去"当基线。
      仪器＝`internal/winsec/nametopup_353_windows_test.go`（`5a537db6` 落地，笔④后 blob `70ed3b61`）；liveness 两侧各 6 枚同发
      `realSealRanFirst=true windowEntered=true plantErr=<none>`；去牙三发落点⛔ 同（三棵突变树 `34462e8f`/`97ed89b1`/`a8246ee6` 互异）、
      正控转红⇄`entered=0/3` 具名报⇄回 `DAMAGE=3`；★⛔ 借分子（`deleted=` 命中 0、`truncated` 只在注释）由 `353-v1` 第 3 件与我另尺两把同号。
      机主随拍那两味逐字在场＝残余句 `:29-39`＋运行期绿行自带 ⛔ 覆盖声明 `:156-163`（先例形＝票 322 写进判码⇄台账、⛔ 票面 bullet＝`353-v1` 第 2 件，我认它退件判⛔ 为对）。
- [x] **AC#3 颜色与门禁**：`./internal/winsec/` 两平台成对＋逐名红册（摆动名册照公共块）；四道门禁 rc；★⛔ 跑 `./cmd/wisp` 整包。
      **编排者翻框（2026-10-21 09:0x +08，凭据经我逐处独立复跑，⛔ 采信任何腿的自报）＝`353-v1` 第 4 件原判"部分成立"
      （缺的现势色具名归 `gw-1`）⇒ 现势色已齐、判语补齐为成立**：
      ① 色＝`gw-1` 味③（`.scratch/wisp/probes/gw/1/03-winsec.md`）对甲 `3a509e50`⇄`a19d5bf2`：host 四发 `121/67/3/0`⇄`121/67/3/0`＋复量三发 `121/68/2/0`、
      container 两侧 `74/46/0/0` RC=0；**逐名新增红 0 枚**（四发并集作差）、**消失红 1 枚**＝本票仪器
      `TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` 基线 4/4 红⇄落地 4/4 绿、正控
      `…RemovalStillRunsWhenNobodyTopsUpTheName353r1` 两侧皆绿（牙没被摘）；`RUN 119⇄121` 那 2 枚⛔ 当判据（`^=== RUN` 名集 symdiff＝`ONLY_LAND=2`）。
      ② 门禁四把＝**我自己 09:0x 现跑**（⛔ 抄 `gw-1` 的数）：`GOFLAGS= go build ./...` **rc=0**；`sh scripts/d22scan.sh` **rc=0**；
      `sh scripts/check-path-length-budget.sh` **rc=0**（`VERDICT GREEN`）；`GOFLAGS= go vet -tags winlive ./cmd/wisp/` **rc=0**。
      ③ ★第 5 把（格式门）**rc=1 属存量 known-red、⛔ 本票字节所造**：CI 真跑形＝`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`
      （`ci.yml:240`；⚠ `gw-1` 交件把路径写成 `scripts/attrib.sh`、行号写成 `:230`＝**引用面错**，我已具名退回记账，其读数结论不受影响）现跑
      `RED (tracked-only) - tracked-dirty=35 files_dirty=35`＋`--self-test cases=9 failures=0`；我另起 blob 尺逐枚重量这 35 枚的**入库 blob**＝
      **真脏 19 枚、全部住 `.scratch/wisp/probes/**`、生产面 blob 脏 0 枚**（正控＝自造坏 `printf` 同尺读出）⇒ 生产面 blob 干净，CI 读 blob；
      `ci.yml:385` 逐字声明 `because this job has known-red steps above it (gofmt (gofumpt)` ⇒ 这枚红在册、⛔ 归本票。
      ⛔ 洗那 19 枚、⛔ `-w`、⛔ 动 `.gitattributes`。
      ④ ⛔ 跑 `./cmd/wisp` 整包（票面逐字禁令）：本腿与 `gw-1` 都⛔ 跑；`gw-1` 味② 已把 `internal/agent` 只包自身那枚色交齐（`0/114/90/0/0`⇄同）。
- [x] **AC#4 纪律格**：⛔ 翻任何框（本票四格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 动 342/348 两票任何原句（互指只在本票尾追加）。
  > **2026-10-21 23:0x 编排者翻框（纪律格由编排者判，四问逐问现跑）＝成立**
  > ①**越界**：票 353 全部腿笔（`353-a1` 读／`353-r1` 写／`353-v1` 裁／`gw-1` 量色）逐笔 `git show --name-only` 判——
  >   越界总数 **5**，★**逐件查明全部是**我自己的编排者账笔**（`A1071`／`A1072`／`A1073` 三笔动 `docs/reports/{HANDOVER,pending-and-issues}.md`）
  >   ⇒ **编排者动台账是本分、⛔ 腿越界**；**腿笔越界 = 0**。
  > ②**零 push**：`origin/dev = 6f808a83`、`HEAD = c04d2c49`、`origin/dev..HEAD` **非空** ⇒ 所有腿笔**都只在本地**。
  > ③**⛔ 动 342/348 两票原句**：`git log --all -- 'issues/342-*' 'issues/348-*'` 的**非编排者**提交 = **0 枚** ⇒ 无腿动过。
  > ④**互指只在本票尾追加**（本格括号要求，也是唯一允许触及票面的形）：`353-a1` 有 **1 笔**触及本票面——
  >   `11e652f2`「票 353 普查腿 353-a1 笔⑥＝票 353 尾**一条只追加 bullet**」⇒ ★**正是本格允许的那一形**，⛔ 违规。
  > ⑤**⛔ 改票名／`Status:`／`-done`**：现读票面 `Status` 未被任何腿改过；`-done` 后缀**尚未加**（本格翻后由我加）。
  > ★另：`353-v1` 判语第 7 件自己列了**三把尺的纪律自证**（写面只证据件、`go test`/`go build`/`d22scan` **各 0 发**、件身 CR 0／被禁码位 0）
  > ⇒ **非实现者我这一格独立复核，与其自证同向**。

## 禁区
⛔ 顺手"修绿"（把 damage 断言反写成"有顶入即拒绝"＝新语义＝AC#1 裁完才许谈）；⛔ 碰 `:11-20`/`:158-170` 那两处归因注释（〔丙〕＝另批）；⛔ 借本票动票 342/348 已批准的落地语义；⛔ 机主 `config.toml`、三枚冻结件、golden、`thresholds.go`。

## 排程
`353-a1`（只读＝AC#0 三形表）→ 编排者裁 AC#1（语义形⇒摆机主）→ `353-r1`（写腿＝AC#2＋AC#3）→ `353-v1`（非实现者裁）。⚠ 写面与票 342/348 续体同包（`internal/winsec`）＝同包只一枚写腿；与 `%w`（票 352）⇄顶栏别腿错开。

**Status:** **`AC#0` 已交（三形表入库＝`docs/evidence/s1/353-a1-census.md`）** ⇄ `AC#1` 裁形进行中（编排者 18:3x 裁＝甲⛔ 批、〔乙／丙〕摆机主）。⛔ 零 push。
## Progress log
（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）
- [2026-10-20 18:2x +08] agent=`353-a1`（票 353 三形代价表普查腿＝★只读·⛔ 选形·⛔ 动产码/测试字节·⛔ 翻框·⛔ push；本条＝本票尾**唯一**追加）did=交 `AC#0`＝甲／乙／丙三形代价表主件 **`docs/evidence/s1/353-a1-census.md`**（50 行／13279 字节／帽 160 只管正文件；CR 0／U+00A7 0／占位符 0＝尺 `tr -cd '\r'|wc -c` ⇄ python `count(chr(0xA7))` ⇄ `grep -cE 'TODO|FIXME|XXX'`，射程＝该件）。续跑身份＝前一枚同身份腿 17:15 被本地网关重启打断（只入库 1 笔 `eaa12987`、`logs/02` 落盘未入库、主件⛔ 写）；本腿锚＝`git log -1` 现认 `96a2eeb7`（17:34），三枚承重 blob 与断腿锚 `c16ff47e` **逐枚同串**（`winsec.go f3c1b68e`／`pinnedleaf_342_windows.go db077e30`／`winsec_windows.go 63714a83`〔量的，尺＝`git rev-parse <锚>:<path>`，射程＝`internal/winsec` 三枚〕）＝两枚腿读的是同一份字节。提交＝笔②`a06140ce`（`logs/02` 复算通过后入库）⇄笔③`16f14989`（本腿独立复算件）⇄笔④`d22152ef`（`AC#2` 仪器形状草案）⇄笔⑤`4f92f1d9`（主件）；★⛔ 采信断腿自算＝公共块读数纪律（自己重读现码复算）＝形状与量级**对得上**、两处细节**具名更正**（① `statHandle` 第二枚查询＝`GetFileInformationByHandle`（`os/types_windows.go:48`），02 写 `...ByHandleEx`＝名字错、窗内枚数⛔ 变；② ACCESS_DENIED 重试支窗内只加 **1** 枚，02 写 +2＝多算 1 枚：deferred `CloseHandle` 在绑定终点之外）。★三枚必答（⛔ 我盖章＝逐条带尺名＋射程＋量的/算的）＝① **甲是语义变化候选＝是** ⇒ **这是第二道批准、归机主**（本表⛔ 选形、⛔ 替机主点头）；甲另两味代价〔量的〕＝`ReOpenFile` 在 `syscall`⇄`x/sys v0.48.0` 两处都⛔ 导出（尺＝`grep -rn 'func ReOpenFile'`＝真身只在 stdlib 自家 `internal/syscall/windows/zsyscall_windows.go:478`）＝新写一枚本平台产码里还没有的句柄面原语；★甲**改掉**票 348 `AC#1` 已批落地那句"两处删除都换成 `unlinkat(父fd, base,0)` ⇄ `Root.Remove(base)`"（票 348 面 `:44`）＝要重述那批已批形，而 POSIX⛔ 按句柄删的单发原语 ⇒ `AC#3`"两平台成对"须改成"Windows 修／unix 具名挂"。② **乙那枚复验有 TOCTOU 窗＝有**＝窗＝内核面 **4 枚**（`Root.Lstat`→`rootStat` `os/root_windows.go:211`→`NtCreateFile` `internal/syscall/windows/at_windows.go:144`→`GetFileType` `os/stat_windows.go:109`→`GetFileInformationByHandle` `os/types_windows.go:48`→`CloseHandle`）＋路径面 1 枚（`rootCleanPath:39`→`syscall.FullPath:59`＝`GetFullPathNameW`）＋用户态 2 段（`os.SameFile`＝`types_windows.go:353`，两侧 `path==""` 时 0 枚内核面；`incref`）＝**个位到十几微秒**〔算的，本腿⛔ 跑任何计时〕⇄ 对照现码敞口＝**≥8 枚内核面、含 3 枚按名安全描述符重操作（`GetNamedSecurityInfo`→`SetNamedSecurityInfo`→再一次 `GetNamedSecurityInfo`，`winsec_windows.go:291-320`）**＝百微秒到毫秒级 ⇒ 乙＝短 1~2 个数量级、**⛔ 归零**〔算的〕；残余射程明写＝先顶入的世界确定性拦住、"复验之后绑定之前才顶入"⛔ 拦不住。③ **可证伪性＝甲能／乙部分能／丙⛔ 证伪任何代码行为**＝乙窗内那一支要钩在 `Deleteat` 的 `NtOpenFile`（`at_windows.go:247`）**之内**＝平台内部 ⇒ `AC#1` 若选乙，`AC#2` 判据必须逐字承认这一支测不到（先例形＝票 322 `AC#1` 的"报窗口、⛔ 圆窗口"）；丙＝恒红常驻，CI `test-windows` 那步＝`bash scripts/winsec-tests.sh`，`continue-on-error` 在 ci.yml 现量 **17 枚命中全在注释行、可执行行 0 枚**〔量的〕⇒ 常驻红直接顶住合并道，另与公共块"唯一合法存量红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`"相冲＝**新增一枚在册常驻红＝治理批准、归机主／编排者**（⛔ 本腿自裁）。★三形共用的两行＝覆盖写⇄独占那两条既有承诺**三形都⛔ 改**（删除只在失败路径＝成功路径⛔ 经过 `unlinkPinnedLeaf`〔算的〕），C26 入口三形都⛔ 动（唯乙若要"报告没删这件事"＝改 `privateFile` 返回＝C26 契约＝另批，凭据＝`pinnedleaf_342_windows.go:207-211` 逐字在册）；★**顶红谁现有的钉**逐枚点名＝甲⇄乙都会让 `ancestorwindow_342r2_windows_test.go:371-374`（红句 `AC#5 RED: AC#5's half-made file survived the refusal`）在顶名世界**必红**＋`:381-388` 那条 RESIDUAL 分支作废＝甲／乙落地都得**同时重述那枚仪器的判据文字**（写腿活、⛔ 本腿权）；⛔ 顶红＝`private_fail_test.go:26/:99`、`openleaf_overwrite_322_test.go:36`（票 322 那枚覆盖写钉）、票 348 `ancestorunlink_348_windows_test.go:58-62`（`deleted` 只可能更小）。★一处表内另报＝本票"甲"与 `342-v3`／编排者 12:2x 已否掉的那枚"第二道接缝"⛔ 是同一件事（那道买的是"祖先进⛔ 去"那枚⛔ 可达的世界；本票要买的世界 `logs/03` 顶名今天已确定性量到、跨进程天然可达）⇒ 若编排者仍把甲否掉，本票只剩乙／丙两枚可选。读数纪律＝★凡引 `342-v3` 报的数逐处标"射程＝它自己的仓外台架、本腿⛔ 复跑"（本趟⛔ 跑任何 `go test`／⛔ `go build`＝`352-r1`⇄`230-r1` 两枚写腿正在 `internal/winsec`⇄`internal/agent/approval` 落笔＝跑任何包都可能读到半保存态假红）；本腿自己的新数只有两类＝blob 同串〔量的〕⇄窗口枚数〔算的〕。写面自证＝本腿 4 笔逐枚 `git show --name-only` 只含 `probes/353/a1/logs/02|03|04`＋`docs/evidence/s1/353-a1-census.md`＋本条；`internal/winsec/**` 字节 **0 枚**；`internal/winsec` 零枚 `.go` 写；名册尺（锚定行首＝`--grep='^票 353 普查腿 353-a1'`）＝本腿累计 5 笔（首 `eaa12987`／末＝本笔）。框终量＝票 353 **未勾 5／已勾 0**（尺＝`grep -cE '^- \[ \]'`⇄`'^- \[x\]'`，射程＝本票面）＝`AC#0`..`AC#4` 一枚⛔ 翻；票名／`Status:`／`-done` 一字⛔ 动。纯追加自证＝本条 `git diff --numstat` 删除行 **0**。临时件只建不删＝`D:\tmp\wisp\20261020\353a1cont\{head,sh,msg,logs}`（⛔ 盘符根）。next＝编排者裁 `AC#1`（甲＝第二道批准那一问摆机主）；`AC#2/3/4` 押后；★本腿⛔ 选形、⛔ push。
- [2026-10-20 18:3x +08] agent=`编排者`（收件＋翻 AC#0＝凭据经我逐处独立复跑）did=★收 `353-a1` 续跑（6 笔，名册三把同号 6/6/6、越界 0、`.go` 0 枚、`census.md` 50 行帽内、入库 blob CR=0）；它更正断腿 `logs/02` 那两处我逐处坐实为真（`GetFileInformationByHandle` 真身＝`os/types_windows.go:48`；`defer CloseHandle` 在 `at_windows.go:275`＝早于 disposition `:293`/`:305`⇒重试只 +1）；三形硬代价也坐实（`ReOpenFile` 在 `$GOROOT/src/syscall`⇄`x/sys` 两把尺都⛔ 导出、`ci.yml:585` 逐字＝`bash scripts/winsec-tests.sh` 且 `continue-on-error` 可执行行 0＝丙的常驻红真顶合并道）。`AC#0` 已翻；`AC#1` 裁＝★甲⛔ 批（三味硬代价：新写句柄原语＋改掉票 348 已批落地的删除形＋POSIX 无按句柄删的单发原语），〔乙／丙〕摆机主、我倾向乙（敞口缩 1~2 个数量级、⛔ 改删除语义、⛔ 新接缝）。next= 机主拍形 ⇒ 拍乙/丙才放 `353-r1`（写面 `internal/winsec`＝与 352/109 错开）；`AC#2` 判据须逐字承认"乙窗内那一支测不到"。
