# 票 321 · 321-r1 复验腿 2026-10-20 交件报告（★零新字节＝AC#1 早已在 HEAD 里落地）

- 腿＝票 321 落地腿 321-r1（★本名在 2026-10-15 已跑过一趟并交件＝`58de465a`；编排者 2026-10-20 18:45 的派单
  存档 `1104fea8` 说"防双派四把尺＝前缀 0 笔"⇒ 本趟是**同一枚腿被重放**，不是新缺口。见 §1）。
- 起手尺全部现跑，⛔ 采信票尾 2026-10-15 那条 bullet 的任何读数；下面每个数都是本腿自己今天跑出来的。
- 写面自证＝本腿 `cmd/wisp/**` 产码字节 **0 枚**、测试字节 **0 枚**（判据已在盘上、且经复验有牙＝再钉一枚＝重复件）；
  落盘只有 `.scratch/wisp/probes/321/r1/logs/20261020-recheck-*`（三枚）＋票 321 尾**一条**只追加 bullet。
  ⛔ 动 `internal/winsec/**`⇄`internal/models/**`⇄`internal/agent/**`⇄`cmd/wisp/config_reload_223_test.go`。
- ⛔ 翻框（票 321 终量＝未勾 2／已勾 0，见 §5）、⛔ 改票名/`Status:`、⛔ push。

## 1. 防双派那把尺为什么漏了（★这条是给编排者的收件尺缺陷，不是本票判据）

| 尺 | 逐字命令 | 起手读数（HEAD＝`1104fea8`） |
|---|---|---|
| 锚定前缀名册 | `git log --all --grep='^票 321' --format='%h %ad %s'` | **1 笔**＝`1104fea8`（派单存档自己）⇒ 看着像"0 笔在飞" |
| 非锚定同串 | `git log --all --grep='票 321' --format='%h %ad %s'` | **4 笔**＝`1104fea8`／`41673704`／`58de465a 10-15 20:03 321-r1 落地腿交件（票 321，零产码）…`／`f19bce28` ⇒ **落地笔在内** |
| 腿码任意位 | `git log --all --grep='321-r1' --format='%h %ad %s'` | **7 笔**，其中 `c673f67b 10-15 20:53 收件 321-r1＋120-r2…` ⇒ **收件点也在史上** |
| 判据件在不在 HEAD | `git cat-file -e HEAD:cmd/wisp/exec_holder_toggle_321_test.go` | **EXISTS**（200 行／8671 字节） |
| 盘上⇄HEAD 同字节 | `sha256sum` ⇄ `git show HEAD:… \| sha256sum` | 两枚同为 `717b2d2d9b76d8ae9d413e53a14e11f1669b834fcb056af6b884e01febb63738` |

★两处归因逐枚复核过（笔⑤之后本腿自己把三把尺各跑一遍，件＝仓外 `roster-recheck.sh`，逐字输出在读数件 §J）：
`c673f67b`（收件那笔）**不**含字面"票 321"（`git log -1 --format=%B c673f67b | grep -c '票 321'`＝**0**，
正控＝同件 `grep -c '321-r1'`＝**1**）⇒ 命中它的是**第三把**尺（腿码任意位），⛔ 记在非锚定那把上。
`58de465a`（落地那笔）两把都命中（`grep -c '票 321'`＝**1**）。⚠ 枚数随 HEAD 漂＝公共块 名册尺 那条通则今天再中一次：
本腿交完 5 笔后同一把锚定尺从 **1→6**（自己的 5 笔全含"票 321"）⇒ 只认逐枚 `%s` 起头，⛔ 把某时刻的枚数当通则。

根因＝那笔落地笔的**主题形状**是「`321-r1 落地腿交件（票 321，…）`」（逐字 `%s` 前 40 字符已现量）——
"票 321"在**括号里**、行首是腿码，所以行首锚定尺按构造必然看不见它（锚定本身没错，是公共块 名册尺 那句
"⛔ 假设其后任何字符"反过来也成立：**⛔ 假设主题以"票 NN"起头**）。⇒ 建议编排者把防双派尺补第三把：
`git cat-file -e HEAD:<该票判据件路径>` ＋ `grep -c '321-r1' <票尾>`（＝"活儿是不是已经有人交过"直接量盘上物，
不依赖主题形状）。本腿零成本撞上这条：起手第一次 `cat-file -e` 就撞出来了。

## 2. 判据本体（现量，逐字件名＋尺名＋射程）

件＝`cmd/wisp/exec_holder_toggle_321_test.go`（无 build tag＝headless 也编译）。两枚尺：

- `TestAC321PanelHopDrivesTheHolderToggle`＝**接线尺**：AST 读 `runResident` 体（尺＝`runResidentBody258`，
  `resident_hotkey_258_test.go:25`，`parser.ParseFile` 直读盘上 `resident_windows.go`），要求
  `startResidentBall` 恰 1 枚、其参数列里 `withPanelHost` 恰 1 枚、实参闭包对 `newPanelHolder` 绑定的变量
  **恰驱动 1 枚方法且名为 `toggle`**（A1 计数＝`:148`、A2 有无执行器＝`:153`、A3 语义＝`:157`）。
- `TestAC321NamedExecutorIsTheToggleRoad`＝**能力尺**：拿 A 读出的方法名去全包产码（尺＝`productionFiles228`，
  `resident_ball_228_test.go:64`，`os.ReadDir`＋名字排除 `_test.go`）找声明、钉枚数＝1，打开体断言发出 `RequestToggle`。

⇒ **不是词面尺**：把 `show` 的体重命名为 `toggle` 那一形红在能力尺。装配处今天＝`cmd/wisp/resident_windows.go:288-290`
（`rb := startResidentBall(..., withPanelHost(func(via string) bool { return holder.toggle(via) }), ...)`）
＝票面/来源表写的 `:268-270` **已漂 20 行**（票 329/351 各动过这枚文件）⇒ 本腿全程用内容锚，零行号依赖。

## 3. 复验三态（★台面＝仓外导出树两棵，`git archive HEAD` 八件套＋`docs/`＋手工补 3 味）

临时根＝`D:\tmp\wisp\20261020\321r1\`（⛔ 盘符根；只建不删）。`treeB`＝基线、`treeM`＝变异台。
两棵 `.go` 各 **777 枚**（非空读数）。`resident_windows.go` 两棵起步同 sha＝`4f7a10b3…`（＝今天 HEAD blob，
`git show HEAD:cmd/wisp/resident_windows.go | sha256sum` 逐字同串＝pristine 跑前验过）。

| 发 | 尺（`-count=1 -timeout 30m -run <形> ./cmd/wisp/ -v`） | 落刀前锚 | 落刀后锚 | `RUN/PASS_TOP/FAIL_TOP/SKIP` | rc |
|---|---|---|---|---|---|
| 正控① 基线 `treeB` `-run AC321` | 判据两枚 | — | — | **2/2/0/0** | **0** |
| 成对 基线 `treeB` 装配族 | 见 §4 | — | — | **23/18/3/2** | 1（红＝§4 存量，非本腿） |
| **突变甲**（票面指定形）`treeM` `-run AC321` | 判据两枚 | `holder.toggle(via)`＝**1**、`holder.show(via)`＝**0** | `holder.toggle(via)`＝**0**、`holder.show(via)`＝**1**；全文件 diff＝`289c289`＝**加 1 删 1** | **2/0/2/0** | **1** |
| 突变甲 反面对手尺 `-run TestBallPanelGesturesReachThePanelThread` | 票面"反面"那枚 | 同上 | 同上 | **1/1/0/0** | 0（**仍绿＝缺口属实被复认**） |
| 突变甲 族色（＝§4 那一发） | 装配族 23 枚 | 同上 | 同上 | **23/16/5/2** | 1 |
| **突变乙**（反形：闭包体 `return false`，不驱动任何 holder 方法） | 判据两枚 | `holder.toggle(via)`＝1 | `holder.toggle(via)`＝**0**、全文件 diff＝`289c289`＝1/1 | **2/0/2/0** | **1** |
| 还原后（甲、乙各一次）`treeM` | sha＋锚＋`-run AC321` | — | toggle＝1／show＝0 | **2/2/0/0** | **0** |

还原哈希（两轮回同一串，且 == HEAD blob）＝`4f7a10b30d1801325a5fdf95c3c975de1bff29057b98bebf09ebe83145837237`；
落刀后＝甲 `5f32c102…`／乙 `cba71d39…`。备份件＝`D:\tmp\wisp\20261020\321r1\resident_windows.go.treeM.bak`（留底、⛔ 未删）。
共享工作树：跑前跑后 `git status --porcelain -- cmd/wisp/`＝**0 行**；全程 pristine 两侧都验。

★方法学复认（与前腿 bullet 同向、本腿独立撞到）：变异**只能物理改写**＋sha 还原，⛔ `-overlay`——两枚尺运行时
`parser.ParseFile`/`os.ReadDir` 读的是**盘上产码**，overlay 只换编译期输入 ⇒ 用它打这一边尺会**恒绿＝假读数**。
本腿另撞一枚台面坑并如实入账：把整条 PATH 用 `cygpath -wp` 换 Windows 形后**连 bash 自己的 `grep/head/sed` 都找不着**
（rc=127、四数 0/0/0/0，长得像判据响）；正解＝派单逐字那形 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`
（POSIX 冒号形）只对 `go test` 那一发生效即可，跑起来 `rc=0`＋两枚 PASS＝DLL 这条路真打通了（不是台损读数）。

## 4. 成对四数 ⇄ 逐名作差＝新增红 2 枚、且那 2 枚**正是判据本身**

★先给"落地⇄基线"那一对（本腿的落地字节＝**0 枚新 Go 字节**⇒落地树≡基线树＝同一份 HEAD 字节）：
尺＝`-run 'AC321|Test258AssemblyRootWiresTheChainAndTheBridge|TestAC228ResidentLegIsTheBallHost|TestAC228BallHostAnswersEveryGesture'`
＝基线 `treeB` **5/5/0/0 rc=0**（`b-pair-assembly.txt`）⇄ 变异台还原后 `treeM` **5/5/0/0 rc=0**（`m-pair-assembly-restored.txt`），
逐名作差＝**新增红 0 枚**（两发名册逐字同集）。⇒ 本腿没把任何东西打红，也没拿整包冒充。

下面那对才是变异的判别力量（同一枚正则、同一时刻、两棵树串行）：

- 基线 `treeB`＝`RUN=23 PASS_TOP=18 FAIL_TOP=3 SKIP=2`；红名册＝
  `Test258BridgeMutationNoSrcKeepsOldBinding`／`Test258OccupiedCombinationNamesTheNewValue`／`Test258V1ProbeSummonEditRebindsLiveBall`；
  SKIP 名册＝`Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable`／`…WithOwnerEmptySlots`。
- 突变甲 `treeM`＝`RUN=23 PASS_TOP=16 FAIL_TOP=5 SKIP=2`；红名册＝上面那 3 枚 **＋** `TestAC321PanelHopDrivesTheHolderToggle` ＋ `TestAC321NamedExecutorIsTheToggleRoad`。
- **逐名作差**＝＋2／−0；SKIP 名册与基线**逐字同集**（分母没动）。⇒ 变异只打红"这一跳是不是 toggle"那两枚尺＝判据**专红**、不牵连。
- 那 3 枚存量红＋2 枚 SKIP-LOUD 的归因＝**机主常驻 `wisp.exe`（现量 pid 21892 在场）占着 `Ctrl+Alt+M`／`Ctrl+Alt+P`
  那一族热键**（逐字红句已入 `logs/20261020-recheck-rulers.txt` §F 末段：`mute=Ctrl+Alt+M occupied by another app`、
  `panel=Ctrl+Alt+P occupied by another app`、`live slots after the rebind = 1, want 2`）＝公共块**弃权名册**那一族，
  ⛔ 归因给本票字节（本腿产码字节 0，两棵树除那一行锚之外逐字节同）。
- ⛔ 拿整包 `./cmd/wisp/` 当"全绿"或"我被红"（314／`resident_*` 真窗区整进程死亡族）＝全程定向 `-run`，没编整包分母。

## 5. 门禁四把＋框终量（只在共享工作树；跑前 `unset GOTMPDIR TMPDIR TEMP TMP` 已验＝两串都 `[unset]`）

| 门 | 逐字命令 | rc | 读数 |
|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | **0** | 空输出 |
| 2 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（射程含 `cmd/` 143 枚 Go 文件、注释与 `_test.go` 都扫） |
| 3 | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/exec_holder_toggle_321_test.go` ⇄ `gofmt -l` 并排 | **0** ⇄ **0** | 两侧输出 **0 字节**＝该件已合规（本腿没新 Go 字节可验） |
| 4 | `go vet -tags winlive ./cmd/wisp/`（票面 AC#2 逐字那把） | **0** | 空输出＝0 字节 |
| 附 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN`，分母 17209 枚跟踪路径、over-budget=60＝roster=60、not in roster=0 |

★这把 `go vet` 尺的**正控**（rc=0 得先证明这尺能红，否则那枚 0 是没验过的读数）：第三棵仓外树 `treeV`＝`treeB` 的副本
＋一枚故意带坏 printf 的探针件（`fmt.Printf("%d", "not-an-int")`；⛔ 落共享树、⛔ 落成对那两棵⇒§3 的 pristine 主张不受影响）
⇒ 同一把尺 `rc=1` ＋ 逐字诊断 `cmd\wisp\zz321r1vetprobe_windows.go:11:38: fmt.Printf format %d has arg "not-an-int" of wrong type string`。
另补一把同族正控：`go vet -tags winlive ./cmd/wisp-does-not-exist-321r1/`＝`RC=1`（`directory not found`）＝包路径写错时尺会响、⛔ 静默 0。

票 321 框终量（尺＝对票面数行首：`grep -c '^- \[ \]'`＝**2**、`grep -c '^- \[x\]'`＝**0**；
正控＝`grep -c '^- \['` 回 **4** ＝那 2 枚 AC 框＋Progress log 两条日期 bullet，两把尺相加对得上＝尺没瞎）
⇒ **未勾 2 枚（AC#1／AC#2）／已勾 0 枚**＝本腿一枚没翻。`**Status:** 待派`⇄票名＝未动。

## 6. 本腿据实上报的三件（⛔ 自行假设、⛔ 越面修）

1. **★重复派单**：AC#1 判据早在 `58de465a`（10-15 20:03）落地、`c673f67b`（10-15 20:53）已被收件；今天 18:45
   的派单按"前缀 0 笔"放腿＝**同一枚腿重放**。本腿据此⛔ 新钉第二枚同形判据（那会造出两枚守同一格的尺、
   并把"哪一枚承重"变糊），改交复验三态＋名册证据。裁决权在编排者／非实现者：票面 AC#1 的两格**该不该翻**
   不由本腿盖章。
2. **前腿具名上报的那枚存量红今天已消**：`TestTicket255RosterStillMatchesTheActualReadSites` 本腿现跑
   `RUN=1 PASS_TOP=1 FAIL_TOP=0 SKIP=0 rc=0`（`section "llm"` 名册已含 `cmd/wisp/panel_models_320.go`、
   `section "models"` 已含 `internal/panel/composer_dispatch.go`）⇒ 由 `9c5916e8`（10-15 21:42，票 324-r1）修掉；
   前腿那条"归 255/320 不归本票"的账可结。
3. **★票面"零产码就钉不住"那一问：不需要产码，本腿没撞上停手条件**。判据读的是 AST 形状＋能力，
   装配处今天仍是 `withPanelHost(func(via string) bool { return holder.toggle(via) })`；本腿把行号漂了 20 行
   （`:268-270`→`:288-290`）这一事实也一并验掉＝尺⛔ 依赖行号。前腿自陈的三枚"过严形"（方法值拼法
   `withPanelHost(holder.toggle)`／`via` 实参须为闭包形参／toggle 兼发 `RequestShow`＋`RequestRaise`）今天**仍红**
   ＝过严、不是没牙；改装配者复判本文件、⛔ 把断言改松（那三形归非实现者裁，本腿⛔ 动判据字节）。

## 7. 交件件册（本腿新落五枚＋票尾一条；`wc -c`≠0、⛔ `.out`、⛔ 0 字节）

- `.scratch/wisp/probes/321/r1/logs/20261020-recheck-report.md`＝本件（判语⇄三态表⇄上报）。
- `.scratch/wisp/probes/321/r1/logs/20261020-recheck-rulers.txt`＝锚／哈希／红句／四数／门禁 `rc=`／正控 的**逐字**输出。
- `.scratch/wisp/probes/321/r1/logs/20261020-recheck-build-rulers.sh`⇄`…-build-rulers2.sh`＝上面那把尺的**生成器**
  （后缀⛔ `.go`；把 §3 那枚 `cut -c` 截中文＝非法 UTF-8 的坑与正解 `LC_ALL=C.utf8`＋`gawk substr` 写死在注释里）。
- `.scratch/wisp/probes/321/r1/logs/20261020-recheck-runner.sh`＝今天那把颜色尺的定式件（同字节来自仓外 `run.sh`；
  注释里写死两味坑＝DLL PATH 缺⇒`0xc0000135`＋四数 0/0/0/0、整条 PATH 转 Windows 形⇒bash 自己找⛔ 到 grep）。
- `.scratch/wisp/issues/321-…-rewired.md`＝尾**一条**只追加 bullet（numstat 删除行＝0，见 rulers 件 §H 之后的补量）。
- 前腿 15 枚件逐枚**未动**（同目录共存，本腿新件一律带 `20261020-` 前缀＝唯一键不撞）。
