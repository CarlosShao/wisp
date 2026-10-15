## A996（2026-10-15 21:5x +08，编排者）收 `324-r1`（dev 上那枚红已由我亲跑转绿）＋收 `322-a1`（三形代价表）＋**★票 322 裁形**＋我自己那发探针的三个读数

### 一、收 `324-r1`（三笔 `9c5916e8`／`e038304b`／`60481441`，20 枚路径、逐文件 pathspec）
- 尺①②：提交真在、`docs/evidence/s1/324-roster-text-vs-type.md`（130 行）＋14 枚读数件真在盘。
- **尺③＝我自己跑的那一发**（`cmd/wisp` 定向 `-run`，⛔ 整包＝票 314 死亡区；`GOTMPDIR` 用 C 卷并先 `mkdir -p`）：
  `rc=0`，三枚点名全 PASS —— `TestTicket255RosterStillMatchesTheActualReadSites (1.75s)`／`TestTicket324RosterScanRefusesSameNamedNonConfigFieldSite (0.03s)`／`TestTicket324SectionMembersAreReadOffTheSchema (0.00s)`。
  ⇒ **dev 上那枚"起手即红的存量红"今天转绿，且是我量到的**（此前只到 `A991` 的定性）。
- **门禁我自己重跑一遍**（⛔ 抄腿的 rc）：`go build ./...`=0／`gofumpt -l` 两枚 go 件＝0（清单 0 字节）／`go vet ./cmd/wisp/`=0／`sh scripts/d22scan.sh`=0（`internal/` 分母 **568**，比 `A994` 时的 564 又涨 4＝别腿新测试件，各 scope ⛔ 下降）／长度门 `--with-self-test`=0（**GREEN 60=60**）。
- ★**腿报的那枚 `d22scan rc=1` 我已复核＝不再成立**：它当时点是 `287-r1` 在飞的 `internal/tools/subagent_depth_fixture_287_test.go:249` 一枚 emoji；我 21:5x 自己跑＝**0 clean**（同一枚文件已在树里、`internal/=568`）⇒ `287-r1` 中途自己改了。**归因照记**：那不是我修的、也不是存量债。
- **我裁（形）**：**接受"乙′"＝文本匹配＋那一跳必须是段结构体自己声明的导出字段（成员名册用反射现读 `config.Config`）作为这把尺的定型**。理由逐条：① 腿具名报了严格甲的**硬阻塞**——`[voice]` 行今天**有意**收下 `internal/panel/config_handlers.go`（接收者是 `panel.SettingsView`、⛔ 配置类型，出处＝票 313 裁定），严格甲会把它摘掉 ⇒ 那是一次**越界的名册重裁**，不该由一把尺的修法顺手做；② 全形甲要新增 `x/tools/go/packages` 依赖（`go.mod` ⛔ 在该腿写面）或对整条链做类型检查；③ 甲的"方法读被误拒"这一维**腿是量过才写的**（14 枚段上方法 0 枚、`internal/config` 46 枚内嵌 0 枚＝今天零命中）⇒ 属**已量到的残余**、⛔ 凭空担心。撤销口令：**「324 改裁甲」**（⇒ 先要一次票 313 那行的重裁，⛔ 让尺自己决定）。
- **AC 框一枚不翻**：票 324 的实现者＝`324-r1`，判语必须出自**非实现者**（`AGENTS.md` §0 第 3 条）⇒ 排 `324-v1`（面＝`cmd/wisp`，⛔ 与 `321-v1`/`165-r2` 同批抢同一面，按序放）。
- 腿自报两处不合规矩（具名、⛔ 抹平）：① 归因用的两枚仓外快照树 `rm -rf` 过（违规则 8 字面）；② 起手 `-- internal cmd` 的对照值没留档 ⇒ AC#5 只能答到"cmd/ 终态 0 行＋internal/ 只剩别腿在飞件"。另：它在 255 尺里**加严**了两处（成员数地板＋逐枚日志），验收腿必答题＝"这两处只加不减"。

### 二、★票 322 裁形（来源＝`322-a1` 的三形代价表＋**我自己那一发探针**）
**我自己的读数**（台件＝`/d/tmp/orch322probe/main.go`，容器 `golang:1.27`，`rc=0`；第一发⛔ 编译过＝`os.Is` 应为 `errors.Is`，记我）：
- `case1` 叶子是相对目标链接＋`O_WRONLY|O_CREATE|O_TRUNC|O_NOFOLLOW` ⇒ **`ERR … too many levels of symbolic links`、`isELOOP=true`**，被指的那个真文件 **size=22 不动**；
- `case2` 活性正控（普通文件，同一把 open）⇒ **`OPENED`、size 20→7**（证仪器不是死的、"0 损伤"⛔ 是恒真）；
- `case3` **祖先是链接**（`via -> sub`）⇒ **`OPENED`、目标 size 24→7＝被截断** ⇒ ★`O_NOFOLLOW` **只护最后一个组件、对祖先那一支完全失明**（与腿报同形，我独立复现）；
- 树外见证文件全程 **64→64 不动**。
**我裁（⛔ 摆给机主，按"账目/工程归我、功能问他"的口径这事不到他清单上）**：**分平台两半，⛔ 单一形**：
1. **POSIX＝甲**（`O_NOFOLLOW`，按平台拆常量、`!windows`／`windows` 两文件；编译矩阵腿已实测：unix 全家族可用、`windows/js/plan9` undefined ⇒ 甲形只需两文件）。★**硬约束**：票面与注释**只许写"叶子那一支关着"**，**必须具名写"祖先那一支今天仍开着"**（我自己的 `case3` 就是凭据）＋并进票 120 已落地的逐组件钉去覆盖它——⛔ 任何形式的"窗口已关"。
2. **Windows＝"重解析点判据"形**：以 `FILE_FLAG_OPEN_REPARSE_POINT` 打开**只为判"这个名字是不是链接"**，是 ⇒ **具名拒绝**；⛔ 把它当"不跟随的覆盖开"（腿实测 `CREATE_ALWAYS` 变体会把链接条目换成空普通文件＝另一种破坏）、⛔ 用 POSIX 绿冒充 Windows 绿（AC#3 的〔仅本机可量〕口径保留）。这一半必做——因为**本机洞是量出来的**（`damaged=571/759/780 /3000` 且 `ok=868/1252/1216` 里含 `err=<nil>`＝截断了外来文件还报成功）。
3. **⛔ 乙（两跳）**：腿实测"摘掉第二跳就退化回跟随"（`damaged=669/902/1022` 与基线同量级），而"覆盖写"承诺今天**没有任何一枚钉**（四处调用者依赖它）⇒ 乙的失败模式恰好是本票要防的那件事。**⛔ 丙**：诚实报告形（`damaged=165/135/148/142`≈基线），选了它就⛔ 写"防住"。
4. **落地腿自带的两件硬活**（写进派单，⛔ 事后补）：**AC#2 的覆盖写新钉**（`secret/store.go:79` 同 id 覆盖／`config/migrate.go:93`／`secret/migrate.go:189`／`agent/spill.go:166` 逐枚点名）＋**AC#1 成对读数**（修前分子>0 且旁站分母与活性同屏、修后 0 且 liveness 非零）。
5. 另登记一笔**不归 322 的现量**（腿交回、我读码复认）：`internal/winsec/winsec_windows.go:328` = `func sealHandle(f *os.File) error { return applyDescriptor(f.Name(), false) }` ⇒ **Windows 侧没有票 120 那套"同一对象复验"**＝**票 120 的 Windows 半格**，另开一格排。撤销口令：**「322 改裁乙」**／**「322 改裁丙」**。

### 三、收 `287-r1`（22:1x，三把尺齐＋我自己复跑）＋我这两发**读数层的自抓**
- 尺①②：提交 `e1973191` 真在（22:09），35 枚文件＝**33 枚 `.scratch/wisp/**`（probes 与逐发日志）＋1 枚 `internal/tools/subagent_depth_fixture_287_test.go`＋1 枚 `docs/evidence/s1/287-depth-fixture.md`（13631 字节、盘上真在）**；越界扫（剥掉这四个落点后）＝**零枚**，产码只落在 `internal/tools` 且只有测试件。追加笔另落台账色 `e64294a0`（那条在本节末的推送对账里逐枚点名）。
- **尺③＝我自己跑的那一发**（定向 `-run 'Test287'`，同台面 `GOTMPDIR` 先建）：`rc=0`／`RUN=3 PASS=3 FAIL=0` —— `Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop (0.15s)`／`Test287RootCallerWithTheSameShapeStillDerives (0.11s)`／`Test287DepthRefusalWinsOverThePoolRefusal (0.10s)`。
- **AC 框一枚不翻**：实现者＝`287-r1`，判语必须出自非实现者 ⇒ 排 `287-v1`。⚠ 它交回的三发种坏（M1/M2/M3）红句与还原 sha 我只**读了、没复跑**（复跑要占 `internal/tools` 编译面，归验收腿）；它具名留下的三格欠账（`cmd/wisp` 端无端到端台件／用例 1 池位次序那一发**未证**／M3 整包档那一发撞到的既有挂死）我照原样记账，⛔ 替它销。
- ★**`287-r1` 的 `AC#3` 那张两形代价表我裁**：**维持甲形**（结构闸在前、产码零动、深度上限 1 不动），`AC#3` 本格**不翻**——因为甲形的代价已被这枚夹具付掉（`:276` 的红现在有人钉），而乙形要摘的是**名字闸**（`subagent_197.go:603-612`），它把保证从"不 offering"降级为"reachable but refused"，是**口径变更**、不是补漏。撤销口令：**「287 改裁乙」**（⇒ 同时要给"refused 也算保证"这句话新写判据）。

**★我这两发读数层的自抓（都带凭据）**
1. **`gofumpt -l` 的"哪一层脏"必须三层并排量**（blob／HEAD 导出树／工作树）。本轮读数：工作树 **6 枚**、导出树 **2 枚**。差的 4 枚（`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`cmd/wisp/models.go`）判据＝**`blobbytes + CR行数 == wtbytes`** 且导出树里 gofumpt 不认脏 ⇒ 是 `core.autocrlf=true` 的**结账形状假枚**（本机 22:1x 现量，六枚全中）。留着的 2 枚（`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`／`cmd/wisp/panel_transport_35r2_test.go`）**两树字节完全相同**（`blobbytes == wtbytes`，blob 里本身就带 CRLF）⇒ **HEAD 内容真脏**，且都是**旧票的测试件（174-r7／35-r2）**、⛔ 本波造的。⇒ 结论：**dev 上没有因为这 6 枚而新起的红**，推送不拦；这 2 枚内容脏登记为存量（`gofumpt -l` 退码本身为 0，⛔ 门不会因此红）。
2. **我第一发的 CRLF 尺是坏的**：同一枚文件、同一条 `grep -c $'\r'`，在脚本里报 **0**、在三层并排那一发报 **114**（且那一发的 `cmp` 与字节数互相对得上＝自洽）。⇒ 定式：**一把尺给的数与另一把独立尺矛盾时，先假定坏的那把是新的那把**，把它留下的坏读数（"archive CRLF=0"）⛔ 用、⛔ 删件（`/d/tmp/orchalign.sh` 原样留档）。
3. **容器面两处读数归"台面"不是"红"**：① 容器里 `go build ./...` 报 `rc=1`（`sherpa-onnx-go-linux: build constraints exclude all Go files`，因 `CGO_ENABLED=0`）⇒ 同一棵 HEAD 导出树**宿主 `go build ./...` rc=0** 反证；② 长度门在导出树 `rc=2` 的红句逐字是 `not inside a git work tree - refusing to guess a denominator`＝**射程拒绝**（它只数 tracked 路径），在仓根跑＝`rc=0`、`GREEN 60=60`（`over-budget=60`、`roster=60`、`tracked=13361`）。
4. **winsec 整包色沿用规则**：我在 `6eb3d2e8` 量到 `rc=0 / RUN=62 / PASS_TOP=40 / FAIL=0 / SKIP=0`（含 `120-r3` 那七枚新用例逐枚点名 PASS）；`6eb3d2e8` 之后落的那几笔**逐枚 diff 点名**见本节末——只有碰到 `internal/winsec/**` 才需要重跑容器，⛔ 拿"看起来无关"当理由沿用。

### 四、编队／推送（本条落账时刻）
- 已收并复跑过：`120-v1`／`322-a1`／`324-r1`／`287-r1`。**在飞两枚**：`120-r3`（`internal/winsec` 测试件，已落四笔 `76e861a3`／`3403c496`／`8beaed39`／`6eb3d2e8`，**证据件与读数目录仍未出现** ⇒ 按未交对待；22:06 它还在写 `gates.sh`＝活腿）∥ `165-r2`（`cmd/wisp` 装配半格，22:03 放，写面独占该面）。
- 待派（已存档、⛔ 算在飞）：`321-v1`／`322-r1`（面＝`internal/winsec`，⛔ 与 `120-r3` 撞面）／`323-a1`（要跑 `cmd/wisp` 定向 `-run`，⛔ 与 `165-r2` 撞面）。
- **推送判据（我自己立、自己守）**：⛔ 把可能带红的状态推上去。判据＝① 当前 HEAD 的 `internal/winsec` 整包 `FAIL=0`（`6eb3d2e8` 上量到，其后至 HEAD 的差集里 `internal/winsec/**` 命中 0 枚 ⇒ 沿用成立）；② 三门禁在当前 HEAD 面上 `rc=0`（`d22scan`／长度门／宿主 `go build ./...`）；③ 推前验 fast-forward（`git rev-list --count HEAD..origin/dev` 必须 0）。本地领先远程笔数与执行结果写在本节末。⛔ 裸推（`branch.dev.remote=cnb`）⇒ 只 `git push origin dev`。

- 待派（已存档、⛔ 算在飞）：`321-v1`／`165-r2`／`322-r1`（面＝`internal/winsec`，⛔ 与 `120-r3` 撞面，按序）。

### 五、你需要做什么
没有。
