# 票 227 二审裁腿 227-v2 判语（裁决者≠实现者：实现腿＝`227-r1`，一审＝`227-v1`；本腿只裁 `AC#5`⇄`AC#6` 两格）

起手锚＝本腿自己 `git log -1` 现认 **`53b8fc6b`**（2026-10-19 19:33:02）⇄ 发整包时 HEAD＝`d86bf31c`（19:35:03）。★尺＝`git diff --name-only 53b8fc6b..HEAD -- '*.go'`＝**0**、`7c44446a..HEAD -- '*.go'`＝**0** ⇒ 目标字节（`cmd/wisp/**`＋`internal/**`）在编排者那一发（`A1053`，锚 `7c44446a`）与本腿两发之间**一字未动**＝作差可比。
台面＝本腿自己 `git archive d86bf31c`（八件套＋仓根 `docs`＋手工三件：`build/*.dll` 3⇄`third_party/sherpa-onnx/*.dll` 3⇄`frontend/dist`）＝`go_files=775`；★⛔ 复用编排者 `orch227-b` 那棵。GOTMPDIR＝`%LOCALAPPDATA%` 同卷短根＋**反斜杠拼法**（本腿自抓的第二枚台损，见"更正"Ⓒ）。PATH＝pwsh 原生整条 Windows 形＋分号（⛔ POSIX 冒号形＝`A1042` 定式没犯）。全部读数＝本腿自己现跑，⛔ 抄腿、⛔ 抄编排者——编排者 `A1053` 那发只当**待作差的对手名册**用。

## AC#5 整包终态＝**成立**（本腿那发四数＝`RUN=2552 / PASS_TOP=1714 / FAIL_TOP=19 / SKIP=13 rc=1`）

尺＋射程＝票面逐字命令的两半（`./cmd/wisp ./internal/...`，`-count=1 -timeout 30m -v`，仓外导出树，逐名读数件 `.scratch/wisp/probes/227/v2/logs/01-*.md`）。票面要的"到终态＋逐名比红名册"两半都齐：★终态＝本腿这发把 `cmd/wisp` 跑到包尾（块内 `=== RUN=448`、源码定义 356 枚顶层用例实跑 **333** 枚，余 23 枚＝19 枚 `windows && winlive` tag＋3 枚 `!windows`＋1 枚⛔ 真用例＝反引号字面量里的票 133 夹具；⚠ 分母那把尺（`grep '^func Test'`）⛔ 感知 Go 词法、差 1 枚，具名登记＝形状全部可解释、⛔ 缺口）。

**逐名作差（⇄编排者 `A1053` 那发 29 枚，两把尺各自闭合：15＋14＝29、14＋5＝19）**
- **只在编排者侧红＝15 枚**＝junction 族 9（含派单点名的 `Test20r2JunctionRefusalIsNeverAskable`）＋SID 卷形 2＋跨卷 1＋`TestAC4InstallDirWidth…` 1＋`TestAC2SealNotice…` 1＋`TestDirOfBaseOfCut…` 1 ⇒ **本腿逐名 `--- PASS`**（逐字行在 `07-attribute.sh` 输出）。与 `A1053`② 归的台损**逐名相符**。
- **两发同红＝14 枚**＝CI 存量 8（`TestAC14GoSideEvalPush…`／`TestApprovalCardViewJSONKeys…`／`TestC21Design…`／`TestC21Table…`／`TestC26RewriteAccount…`／`TestComposerContract…`／`TestStreamLog*35r8`×2＝⛔ 本票新增、⛔ 顺手修）＋落地续体 2（342 那两枚＝342-v2 已判部分成立的 Windows 半区）＋缺 `.git` 2＋`TestComposerRenderFixtureTellsTheTruth` 1＋`TestResolvePerCallBudget` 1。★`TestResolvePerCallBudget` 按票面自注＋`A426/A432` 定式安静定向复量 3 发＝**3/3 PASS**＝争用型假红复认。
- **只在本腿侧红＝5 枚**＝⛔ 一枚产品缺陷：3 枚 `Test258*`（红句逐字 `Hot key is already registered.`＝本机 `Ctrl+Alt+Q/M/P` 被他程序占住，安静窗定向 3 发 `0P/4F`×3＝稳定非摆动）＋`TestAC246DevLegIgnoresTheTestTaskInjection`（红句逐字 `wisp: another instance is running in this session; activated it; exiting`＋盘上现量 `wisp` **PID 21892** 自 10-15 19:03 活着＝单实例互斥挡住这一发；本腿⛔ 动那枚进程）＋`TestTicket252R1AlignmentAddsNoAuthorization`（＝本腿台损：它取仓根**父段**短名、`…\20261019\227v2` 那层⛔ 别名）。
⇒ **⊖〔台损＋桌面机器态＋CI 存量＋落地续体〕后新增产品红＝零枚。** 票面这一格要的"整包终态⇄逐名名册"两半到此齐＝**AC#5 成立**（翻勾归编排者）。

## AC#6 时间炸弹那一格＝**部分成立**（要求②本腿量成"有牙"＝闭合；要求①逐字在场、逐字为假＝仍未做）

读数件＝`.scratch/wisp/probes/227/v2/logs/02-*.md`（四发一台＋逐字红句＋还原哈希）。
- **要求①（那句 `nothing polls`）＝可判、判"未做"，且★假话半坐实到装配层**：`HEAD:cmd/wisp/always_write_no_clobber_226_test.go:43-44` 两行逐字在场（★跨行形＝`nothing` 收 `:43` 尾、`polls:` 起 `:44` 首，本腿⛔ 重犯 `227-v1` 那把单行 `tr` 尺）；blob `5fde7930…` ⇄ 起手 `feac8f77` 同串＝至今没人动。生产侧 1s tick 三处逐字在位（`config_reload.go:92`＋`:123`＋`run.go:861`，`runTextTask:249` 走到），★本腿控制发里该宿主的 `[audit] config: HOT-RELOAD state=armed tick=1s … goroutine=watchdog owner=config` 整行被红句 dump 打出来＝装配真发生 ⇒ **那句今天为假＝复认**。⛔ 本腿改它（产码注释⛔ 我写面）＝具名交编排者或续腿落笔。
- **要求②（常驻用例＝"真轮询下手改仍看得见"）＝★有牙**：去牙两发落点不同——MU-1 摘 `rt.reloadHandle = observe.Default.Spawn(… rt.configReloadTick)`（tick ⛔ 起）⇄ MU-2 把 `rep, err := rt.mgr.CheckAndReload()` 摘成 `var rep/err`（tick 醒着、从不重读）⇒ **`TestTicket223RunArmsTheReloadTick` 两发都红**（红句逐字 `config_reload_223_test.go:354: audit never carried "config: HOT-RELOAD state=applied" within 40s`）；未突变控制发＝`2/2/0/0 rc=0` 两枚全绿＝⛔ 恒红正控过。⇒ 票面 09-29 更正后那句"绿的原因是没人钉"——**"真轮询看得见"这一跳今天有人钉、摘接缝就响**＝要求②成立。
- **★本腿另送一枚硬读数＝时间炸弹的机制被第三种形状坐实**：MU-3⛔ 动断言⛔ 摘接缝、只把 `configReloadPollInterval` 由 `time.Second` 改成 `10 * time.Millisecond`（同形状只加快）⇒ `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` 当场红在 `always_write_no_clobber_226_test.go:103: the guarded write applied the hand-edited theme to the running process (got "light", want the dark this run loaded)`（复量×3 同句、rc=1）。机制＝`app.theme` 确是热档（`internal/config/tiers.go:47`＋`manager.go:357-359 planApp` 热并），`:102` 那条期望与热档相反；它的绿**依赖"1s tick 没在收口前赶到"这一运行学事实**＝⛔ 形状钉（223 那枚钉的是"看得见"、⛔ 这一跳）。⇒ 本格两半合起来＝**要求②闭合（常驻用例有牙）＋要求①未落（那句假话逐字在场）** ⇒ 本腿裁＝**部分成立**，⛔ 由本腿翻框；两样欠的都是同一枚续腿的活：①把 `:43-44` 那句改准（＝只动字面）＋②给"tick 必赶得上"那一跳配一枚确定性形状的用例（写面＝`cmd/wisp/**`＝共享编译面）。★⛔ 把本格读成"成立"：`227-v1` 那枚腿与编排者收 bullet 都写"要求①可判且判未做"＝本腿复认同一形。

## 两枚未定性（派单具名交回，本腿定性＝**都归台损**，根因一枚）

尺＝成对两发（同树同 `-run`、只差 GOTMPDIR 拼法）：★反斜杠短根＝`2/2/0/0 rc=0` **两枚全绿** ⇄ 正斜杠短根＝`2/0/2/0 rc=1` **两枚全红**（`logs/03-*.md` 表①）。
- ① `TestDirOfBaseOfCutOnThePlatformSeparator`＝**台损**（红句逐字已读全：`want` 串里混着 `C:/…` 正斜杠段⇄`dirOf` 的归一化形＝必然不同串；`os.MkdirTemp` 原样带 GOTMPDIR 前缀）。⚠ 同时具名一枚**真缺陷形状**交编排者决定要不要立票：那枚用例把 raw `t.TempDir()` 当 canonical 入参，断言随宿主 env 拼法改色；同包 `helpers_test.go:83 tempCanonical` 已有现成归一化形。
- ② `Test20r2JunctionRefusalIsNeverAskable`＝**台损**（`mklink /J` 在 target 为 `C:/…` 形时被 cmd 开关解析吃掉＝红句逐字 `Invalid switch - "Users"` 与编排者长 D 卷那发**同串**）。

## ★四处与派单／台账矛盾的事实（具名上报，⛔ 自行改判、⛔ 回改已发出的正文）

- Ⓐ **派单第 14 行⇄`A1053`③ 那句"两枚同名不同包（`internal/panel` 同名那枚已绿）"盘上⛔ 成立**：`git grep -n '^func Test20r2JunctionRefusalIsNeverAskable' HEAD` **只命中 1 枚**＝`internal/tools/refusal_20r2_windows_test.go:33`；`git grep -c 'Junction' HEAD -- internal/panel`＝0、`-- cmd/wisp`＝0。旁证＝编排者自己 `attrib/final-attr.log:14` 那枚红紧跟的包行＝`FAIL …/internal/tools`。⇒ 那两枚其实是**同一枚**在长/短两种拼法下的两发；两处描述的方向（"这枚要单独定性"）⇄结论（台损）⛔ 因此动摇。
- Ⓑ **`A1053` 的 29 枚归因在枚数上⛔ 自洽闭合**：末行写"CI 存量 8＋落地续体 2＋台损 17＋未定性 2＝29"，而 `window-facts.txt` 第 5-6 行的台损细分**已把 `dirOf` 计入 17**、⑴ 又把 `20r2` 计入 junction 9 ⇒ ③那两枚**已在 17 之内**＝同两枚被数两次；两份文件之间还互相差一枚（`window-facts` 有 dirOf⛔ AC2Seal，`A1053`② 有 ⑺AC2Seal⛔ dirOf）。本腿的独立闭合＝上面那把 **15⇄14⇄5** 尺。⛔ 改它任何一格结论，只请编排者自己重对归属。
- Ⓒ **junction 那族的根因⛔ 系于"长／D 卷"，系于 `GOTMPDIR` 的拼法**（本腿短 C 卷＋正斜杠照样红，甲形全绿）⇒ 归类（台损）⛔ 变，可复现的修法要改写＝定式加一条：**`GOTMPDIR` 与 PATH 同批守 `cygpath -w`／反斜杠形**（这条对下一枚腿是硬要求，⛔ 只写"短临时根"）。
- Ⓓ **`A1053`⑹ 对 `TestComposerRenderFixtureTellsTheTruth` 的归类要更正**＝本腿拷了 `frontend/dist`（5 枚文件）在两种树上都仍红；红句缺的是 render fixture 的三块具名状态，`composer_test.go:855` 自己写 regenerate 靠 `npm run render:composer` ⇒ 仍属台面形、⛔ 产品缺陷，但补的那一味完全不同。

## 门禁四数（共享工作树，先 `unset GOTMPDIR TMPDIR`）

`gofumpt -l`＝**⛔ 对象**（本腿写面＝`.md`/`.sh`/`.ps1`，全史 `.go` 命中 0，具名登记、⛔ 编数）／`sh scripts/d22scan.sh`＝**rc=0 clean**／`sh scripts/check-path-length-budget.sh`＝**rc=0 VERDICT GREEN**／`go vet ./cmd/wisp/`＝**rc=0**。★票框终量尺（本腿一枚⛔ 翻）＝判据区 `- [ ] **AC#`＝**2**／`- [x]`＝**4**。

## 本腿自报三枚（⛔ 圆场）

① 第一发整包（`full-A` 2537/1706/27/13）＝**本腿自造台损**：树命名 `tree` ⇒ 整条路径每段≤8 字符⛔ 别名 ⇒ 多打红 8 枚 `internal/tools` 票 252 探针；改名重打＝`full-B`（决定性那发）。`full-A` 的四数⛔ 作 AC#5 读数用。② `grep -c $'\r'` 那把尺本腿打坏一次（`$'\r'` 被吞成空模式＝命中恰等于行数，一度让本腿以为入库字节违反 LF 条）；`tr -cd "\r" | wc -c` 重打＝**6 枚已入库 blob CR 字节全 0**、工作树同 0 ⇒ 缺陷⛔ 成立、⛔ 撤，坏尺本身具名留下。③ `09-targeted.sh` A 段四形 mklink 矩阵里本腿⛔ 带引号手工拼串，甲行（两侧反斜杠）也被 cmd 切开＝那一行作废并就地具名；Ⓐ/② 的结论只取与编排者红句逐字同串那一格⇄Go 侧 `exec.Command` 传参的配对台两发。

## 交件名册

判语＝本件。读数散件 3 枚＝`.scratch/wisp/probes/227/v2/logs/01-ac5-four-numbers-and-diff.md`／`02-ac6-comment-and-teeth.md`／`03-unqualified-and-gates.md`。台面仪器 14 枚＝`.scratch/wisp/probes/227/v2/scripts/01…14`（含被本腿具名作废的那两把坏尺）。原始 log 1.5MB 级⛔ 入库＝留仓外 `D:\tmp\wisp\20261019\227v2\`（只建不删）。票尾只追加 1 条 bullet。
