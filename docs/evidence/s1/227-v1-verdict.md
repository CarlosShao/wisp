# 票 227 终裁腿 `227-v1` 判语（裁决者 NOT 实现者：产码腿 = `227-r1` / `ecebc7a9`，本腿只判不写码）

起手锚 = 本腿 `git log -1` 现认 `9fc2be02`；被裁对象 = `ecebc7a9`（`ecebc7a9..HEAD` 的 `.go` diff = 0 枚，现跑；唯一差件 = `docs/reports/pending-and-issues.md`）。
台面 = 仓外 `D:\tmp\wisp\20261022\227v1\`（`treeA` = HEAD 回滚本腿两枚测试件 / `treeB` = HEAD 原样；两树各 267 MB、`.go` 1154 对 1155、`diff -rq` 只差那两枚，`logs/02-export.txt`）。
全部读数本腿自己现跑：NOT 抄腿、NOT 抄编排者、NOT 采信腿自报与编排者的翻框。`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`、`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`、`-timeout 30m` 逐发在位。
**票面一枚没翻**（本腿无翻框权）：现跑尺 = 判据区 `- [ ] **AC#` = 1 / `- [x]` = 5，与起手 `527b7377` 逐字同值；票面 blob `d94a9e89` 在 `527b7377` / `ecebc7a9^` / `ecebc7a9` / `HEAD` 四锚同串 = 没人碰过。

## sec 1 十项逐项判语（每项另起尺）

**(1) 票面前提新鲜度（硬闸）= 成立，两条前提都在 HEAD blob 上活着。**
尺 = `git show HEAD:<path>` 落盘再量（`logs/01-premises.txt`），NOT 读工作树副本。`nothing polls` **跨行形**逐字复认：`:43` 行尾收 `nothing` YES、`:44` 行首带 `polls:` YES；单行 `grep 'nothing polls'` = **1** 命中，但那 1 枚是 `ecebc7a9` 追加的更正注释（`:46`）里引的那半句，**NOT 原句** —— 原句两行在 HEAD 与 `ecebc7a9^` 两枚 blob 里逐字同形，且该件 `-` 行数 = **0**（`--numstat` 16/0）。
`config_reload.go` 三处在位（`:92` `configReloadPollInterval = time.Second`、`:123` `Spawn("watchdog","config",...)`、`:148` `rt.reloadOnce()`）、`run.go:861 rt.startConfigReload()` 在位、`248095d1` 是 HEAD 祖先 YES。`app.theme` = `hot`（`internal/config/tiers.go:47`），故 `:102`（现 `:118`）那条 `!= "dark"` 与 D36 热档反向 = 成立。

**(2) `AC#6` 判据（真轮询之下手改仍看得见，要一枚常驻用例）= 成立。**
尺一（体内禁物，去注释后逐行）：`always_write_no_clobber_227_test.go` 代码行 124/225，`config.NewManager` / `CheckAndReload` / `ConfirmLocked` / `OnRestartPending` 代码命中 **0**（全部只出现在注释里，逐行标 COMMENT 复核）；正控 = 同尺打 `internal/config/manager_test.go` 17 命中、`cmd/wisp/config_sentences_223r2_test.go` 有 `config.NewManager` 实参 => 尺 NOT 空转。
尺二（真路径）：`startConfigReload` 生产调用点仅 `run.go:861`；`reloadOnce` 仅 `config_reload.go:148` 调、其体 `:157 rt.mgr.CheckAndReload()`；`state=applied` 唯一生产发射点 `config_reload.go:183 reportReload`。用例等的就是这一行 => 走 `run.go:861` -> `:123 Spawn` -> `configReloadTick` -> `:148` -> `:157` 那条链。
尺三（有牙）：本腿自造两发，见 (3) 与 (4)。

**(3) 票面逐字禁手（调大 tick / 缩短 sleep）= 没犯；且新用例 NOT 靠时序侥幸。**
尺一（diff）：`ecebc7a9` 对 `*.go` 的新增/删除行里，`Interval`/`time.Second`/`Millisecond`/`Sleep`/`Ticker`/`Timer` 逐行看，**全是新增注释与被等的那两条 40s 上界**（`time.NewTimer(limit)`、`20ms time.NewTicker`）；`config_reload.go` / `run.go` / `internal/**` 命中 **0** 枚文件；`--numstat` 四件全 `+N 0` = 零删除行 => 间隔与 sleep 都没动。
尺二（本腿自己那一发，`logs/52-mc-*`）：把 `configReloadPollInterval` 由 `time.Second` 改成 `10 * time.Millisecond`（**反向**：禁手是放宽，我把它收紧 100 倍；NOT 碰断言、NOT 摘接缝）=> `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` **当场红**（`always_write_no_clobber_226_test.go:119` 逐字 `got "light", want the dark this run loaded`），而 `TestAC6AHandEditIsStillSeenUnderTheRealPoll` **仍绿**（2.10s）。=> 新用例 NOT 靠 1s 窗口活着；旧用例的绿才靠窗口。
尺三（另一向，`logs/42-ma-*`）：tick 1s 改 10s => 新用例仍绿、只是变慢（12.22s；它等的是产品自己的行，NOT 时长）=> 它也不靠"tick 必须快"。两向都绿 = 与间隔无关。
尺四（稳定性）：新用例单枚连打 3 发 = 3/3 PASS（3.17 / 3.26 / 3.39s）。

**(4) 反向判据（它报 M1/M2 都红在 `:192`）= 落点形状成立，但行号在本腿尺下是 `:198`，具名登记。**
本腿自造一发（`logs/45-mb-*`，落点 = `config_reload.go:123` 的 `Spawn`，与它的 M1 同落点但独立重打）：落刀前命中 1、落刀后 `old:0 new:1`，`go build ./...` rc=0，红句逐字
`always_write_no_clobber_227_test.go:198: the production reload tick never booked "config: HOT-RELOAD state=applied hot=[app]" within 40s: ...`，`--- FAIL: TestAC6AHandEditIsStillSeenUnderTheRealPoll (41.70s)`，rc=1；还原后 `config_reload.go` sha256 = `72762bd5...602c` = HEAD blob 逐字同串，复跑控制发两枚全 PASS（`logs/60-*`）。
行号差：HEAD blob 里那条 `t.Fatalf` 在 **:198**（`:197` 是 `w.failure()` 判断），本腿与 `ecebc7a9` 的 blob 哈希同串 `058328d6...7999`，中间副本（14:24:38）与 HEAD 只差两处注释改指（`227-r1-impl.md` 改 `-part1.md`）=> 仓内不存在 `:192` 那一形。它的 M1/M2 日志（`42-mut-m1-red.txt:15` / `45-mut-m2-red.txt:15`）报 `:192`，**判语 NOT 踩行号**（红句文案、`FAIL` 名、rc、还原哈希四件与本腿同向），此处只作具名登记，NOT 据此退件。
M2 那一发本腿没有重打（它的落点 `:148` 与我的 M-b 同族、读数同句），登记为"文案同向、未逐字复跑"。

**(5) "改前读数"那一形（旧用例的绿靠 84-99ms 收口窗口）= 同向，独立复现。**
尺 = 仓外 treeB 内一枚临时探针（`probes/zz227v1_probe_hold_test.go`，不入库），票 226 的同一顺序（手改 -> 一直 -> 受守卫写 -> 收口读），两臂同码：
```
arm-A hold=0      theme="dark"  applied_lines=0  span=50ms                => GREEN（226 断言成立）
arm-B 等 tick 行  theme="light" applied_lines=1  tick_wait=1.003s span=1.027s => RED（got "light", want "dark"）
```
两臂只差"受守卫写落账后还活多久"，**同一个 1s tick**、没有任何 sleep 被缩短或放宽 => 与它报的 dark/light 与 applied_lines 0/1 同向。
一处具名：它三对数字是 `84/95/99ms`，我这一发是 `50ms`（量级同、值不同；它另有一发 `06-probe-226-ordering.txt` 记 `200ms`）=> 该窗口是**区间量**、NOT 定值；它正文与 226 注释里 "measured 58ms" 引的是它 `05` 那发（`58.111ms`），与 part1 sec 1 表里的 `84/95/99` **NOT 同一发**，两处并列读会打架（件里没有 58ms 这个数）。NOT 影响"窗口 < 1 tick"这个结论。
臂 B 我用的判据是**等产品自己的 `state=applied` 行**（有界 timer），NOT `sleep 2.5s` => 读的是"tick 说话之后是什么"，NOT"2.5 秒之后碰巧是什么"。

**(6) 门禁复跑（本腿现跑，rc 与末行逐字）= 五门全 rc=0；格式门禁甲乙两尺分报。**

| 门 | 命令 | rc | 末行逐字 |
|---|---|---|---|
| build | `go build ./...` | 0 | （0 行输出，`logs/10-build.txt` 0 字节） |
| vet | `go vet ./cmd/wisp/` | 0 | （0 行输出，`logs/11-vet.txt` 0 字节） |
| d22scan | `sh scripts/d22scan.sh` | 0 | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=247, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=616, ban #8 cmd/=148; ban #8 emoji coverage: ...` |
| path-length | `sh scripts/check-path-length-budget.sh` | 0 | `check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 61 paths on this tree, recomputed just now in both units`（上一行 `VERDICT GREEN ... roster equals the tree`；分母 **17561** tracked，NOT 17558） |
| gofumpt 甲（**判据**，入库 blob） | `/d/work/base/gopath/bin/gofumpt.exe -l` 对 `git cat-file blob HEAD:<path>` | 0 | 两枚目标件 listing **空**、`-d` 0 字节；全仓 tracked `.go` 1155 枚扫一遍，`.scratch/` 之外红 **0** 枚（含本腿两枚） |
| gofumpt 乙（**只作诊断**，工作树） | 同上对工作树副本 | 无 rc 义 | 红 **4** 枚：`cmd/wisp/models.go`(CR=389) / `internal/agent/approval/pending_read.go`(131) / `internal/agent/tools.go`(225) / `internal/risk/provenance.go`(1115) = 本机 `core.autocrlf=true` 假红，blob 全清 => 按裁定件 `docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md` 判**过** |

两处自纠（具名，件不删只建）：(a) 第一版甲尺用 `gofumpt -l f || echo RED` 判红 = **空尺**（实测：故意未格式化的文件 `gofumpt -l` 仍 **rc=0**、只打路径）=> 改"listing 非空即红"并加正控（正控文件 listing 非空）。(b) 甲尺对**全仓**跑时列出 19 枚 `.scratch/**` 台件，按 `attrib.sh` 自己写的归属规则逐枚归票（`probes/<NN>/**` 且 `issues/<NN>-*.md` 存在）=> **19/19 归属到票 33/163/174/183/185/197/212/220/222/224/231/235/241/253/259**、不归属 0、`.scratch` 之外 0、碰本腿两枚 0（`logs/21-*`）。
另报一条 CI 步的真实读数（不属本腿判据）：`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`（= `ci.yml:240` 那步）本机 **rc=1**、`tracked-dirty=35`，其中 31 枚在 `.scratch/**`、4 枚是上面那四枚 CRLF 假红 => 该步本机不可判，与 `ecebc7a9` 无关（它两枚文件不在那 35 里）。

**(7) 整包读数（主检出不可判）= 主检出死亡点本腿复现；导出树"新增红 0"本腿独立复跑，成立。**
主检出（`logs/80-main-checkout-fullpkg.*`）：`go test -count=1 -timeout 30m -v ./cmd/wisp/` => rc=1、wall=342s、**RUN 225 / 顶层名 167 / PASS 166 / FAIL 0 / SKIP 0**、编译名册 **348**、未跑到 **181** 枚；死亡点 = `=== RUN   TestAC14GoSideEvalPushReachesThePage` 之后**无任何判决行**（该名 result lines = 0），尾停 `MARK299 wndproc ... msg=0x6 known=true` 与 `FAIL github.com/CarlosShao/wisp/cmd/wisp 332.076s`。=> 与编排者复现同点（top RUN=167、判决 0、rc=1）。**主树的"0 枚 FAIL"是"没跑到"、NOT "都绿"**，本腿 NOT 拿它当整包绿。
导出树两本（同一 HEAD `9fc2be02` 的 `git archive`，仓外；`treeA` 回滚两枚测试件作对照、`treeB` 原样）：
```
CONTROL(treeA) rc=1 wall=562s  RUN 472 PASS 321 FAIL 21 SKIP 5   (reached 347 / compiled 347)
AFTER  (treeB) rc=1 wall=562s  RUN 473 PASS 324 FAIL 19 SKIP 5   (reached 348 / compiled 348)
两本最后一条 RUN 都有判决行（TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask）= 都跑到包尾，NOT 截断名册
NEW RED（after 红而 control 不红）= 空
FLIPPED GREEN = TestPanelHostLatencyPercentilesAC2 / TestPanelHostRealWindowHopAndLifecycle
shared red = 19
```
=> **新增红 0 成立**；四数 `473/324/19/5` 与它报的逐字相同（对照侧我读到 `472/321/21/5`，它报 `472/322/20/5`，差在两枚面板窗口族争用型红，见 (8)）。
两本都在包尾 => 这两个数可以当整包读；主检出那本**不可以**。

**(8) 红名册逐名（NOT 只报枚数）= 新增红 0；19 枚共同红逐名，全部在对照树（不含本腿任何字节）上就红。**
```
Test258BridgeMutationNoSrcKeepsOldBinding               TestAC1ResidentLegInstallsItsLogListenerOnDisk
Test258OccupiedCombinationNamesTheNewValue              TestAC1ResidentLegOutlivesItsOwnLogFailure
Test258V1ProbeSummonEditRebindsLiveBall                 TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
TestAC14GoSideEvalPushReachesThePage                    TestAC246ShippedResidentLegWithoutConsoleRefusesItsTaskEntry
TestAC246ShippedResidentProcessOwnsItsCancelStep        TestAC246ShippedResidentLegTakesTheInjectionAndAttemptsThePipeline
TestAC228ResidentLegReportsAndBooksItsBall              TestAC246DevLegIgnoresTheTestTaskInjection
TestAC228ExitRequestDuringBootStillLeavesThroughD38E    TestSecretArgvCarriesNoSecret
TestAC2RealProcessRefusesOnEveryLegWithoutAppData128    TestSecretRealBinaryRefusesValueFlag
TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord
TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
TestCleanCheckoutBuilds_AC11
```
唯一在**我的对照本**里红、而**它的对照本**里不红的两枚 = `TestPanelHostLatencyPercentilesAC2` / `TestPanelHostRealWindowHopAndLifecycle`（反向：它那本红、我这本绿的 = `TestAC4FocusReturnToPriorWindowGap33r5`）=> 三枚全是桌面窗口族**争用型**（同机、两次相邻跑、无本腿字节参与），登记为台面抖动，NOT 算本腿新增、NOT 算它新增。归属：19 枚全部落在 resident / secret / clean-checkout 三族，`always_write_no_clobber_*` 一枚不在册。

**(9) 写面与纪律 = 一笔四件、`.go` 恰 2 枚、越界 0；正控回验命中 2。**
`git show --name-only ecebc7a9` 逐字：`cmd/wisp/always_write_no_clobber_226_test.go`(M) / `cmd/wisp/always_write_no_clobber_227_test.go`(A) / `docs/evidence/s1/227-r1-impl-part1.md`(A) / `docs/evidence/s1/227-r1-impl-part2.md`(A)；`.go` = 2；落在这四件之外的 = **0**；`internal/**` / `run.go` / `config_reload.go` / `scripts/**` / `.github/**` / `.scratch/wisp/**` / `PLAN.md` / `docs/specs/**` 命中 **0**。生产码零字节改动复认：`config_reload.go` 与 `writeguard.go` 的 `HEAD` blob 对工作树 `git hash-object` 逐枚同串。
正控（同尺打一笔**含 `.go`** 的真提交）：同一把尺打 `ecebc7a9` 回验 `.go` = **2**（尺 NOT 空转）；负控 = 同尺打纯文档笔 `9fc2be02` = **0**。严格前缀尺（`^票 227 落地腿 227-r1`）= **9 笔**（本笔 + 8 枚 10-18 旧笔），逐笔 `--name-only` 越界 **0**、`.go` 合计 4 枚次（`b72d7391` 的 `writeguard.go`、`3b8c114e`/`dd0e489d` 的 `writeguard_handedit_227_test.go`、`ecebc7a9` 的两枚）=> 与它自报的"一笔"不同口径：它只报本轮那一笔，历史八笔属上一轮 227-r1，两者不矛盾，登记清楚。
本腿自己的写面：只 `docs/evidence/s1/227-v1-verdict.md`（本件，覆盖 10-18 那枚同名旧件 = 同题换人重判，旧件内容仍在 `088c7171` 里）+ `D:\tmp\wisp\20261022\227v1\**`（仓外）；NOT 落任何 `.go`、NOT 改 `internal/**` / `cmd/**` / `scripts/**` / `.github/**` / `.scratch/wisp/**` / `docs/reports/**` / `PLAN.md` / `docs/specs/**`、NOT 翻框、NOT push。

**(10) `226-v2` sec 6 第 8 条点的那枚"226 写回不 Adopt 内存"确定性用例该不该算 `AC#6` = 不算，属另一格。**
依据三条，各自现跑：
(a) **票面 `AC#6` 逐字只点名一跳**：「真轮询之下，手改仍看得见」—— 要求一 = 那句 `nothing polls` 改准、要求二 = 给这一跳补一枚常驻用例。**"写回不 Adopt 内存"是反向的**（写回之后内存**NOT 许**变），NOT "手改仍看得见"。
(b) **那一跳今天在 `cmd/wisp` 车道仍有牙、且不在本腿射程**：`always_write_no_clobber_226_test.go:118-120` 的 `!= "dark"` 断言仍在（`ecebc7a9` 只追加注释、断言 0 字节改动）；包内裸 Manager 车道另有 `internal/config/writeguard_226_test.go:247-249`（`must not adopt a key it never confirmed`）与 `:284-286`（`must not apply a foreign hand edit to memory`）两枚常驻用例。
(c) **本腿自造一发证明"它不是 `AC#6` 的缺失件"**（`logs/56-md-*`）：把 `writeguard.go` 的 `diverged := diffKeyPaths(...)` 之后插一行无条件 `m.statOwnWrite()`（= `226-v2` 的 a3/rerun1 同形）=> `go build ./...` rc=0、落刀前后锚点命中 1 对 1 且标记 1、`TestAC6AHandEditIsStillSeenUnderTheRealPoll` **当场红**（`always_write_no_clobber_227_test.go:198` 逐字同句，41.70s，rc=1），而 226 那枚仍 PASS；还原后 sha256 = `efe12661...47bd` = HEAD blob 逐字同串、标记 0。
结论：新用例**已经**在真轮询下间接钉住了"写回 NOT 许把文件内容认领成自己的"（无条件认领会把 tick 饿死、红句正是那句 audit 等待），但**票面 `AC#6` 字面点的 NOT 是它** => **该枚算另一格**：建议另立一票（或并入票 226 的残余格）。理由 = 它要断的是"写回后内存与文件的分叉"这一形状，与 `AC#6` 的"轮询把文件改回内存"方向相反，塞进 `AC#6` 会重新撑开一格（与编排者 `9fc2be02` 第 5 条同向）。

## sec 2 `AC#6` 能翻 / NOT 翻 + 缺什么

**判语：能翻（本腿判"成立"），但翻框权 NOT 在本腿** —— 派单逐字「NOT 翻任何框」，故本件只交判语，翻不翻由编排者凭此件做。
判"成立"的三条硬据：(i) 常驻性成立（代码行 0 枚 `NewManager`/`CheckAndReload`，正控回验尺 NOT 空转）；(ii) 真轮询路径成立（唯一生产链 + 唯一 `state=applied` 发射点，本腿 M-b 摘 `Spawn` 即红）；(iii) NOT 靠时序成立（10ms 反向收紧仍绿、10s 放宽仍绿、单枚 3/3 PASS）。
**缺的东西（都不阻断 `AC#6`，但具名交回）**：
- (i) `226_test.go:55` 与 `227_test.go:26` 两处注释写 "measured 58ms"，而 `227-r1-impl-part1.md` sec 1 表里是 `84/95/99ms`，`58.111ms` 在它 `logs/05-probe-226-two-shapes.txt` —— 同一窗口的两发读数被并列引用、件里无 58 这个数。**属文案**，建议续腿只动字面。
- (ii) `part2` sec 5 表里 path-length 的分母写 `17558 tracked`，本腿现跑 = **17561**（差 3，随时间自然增长）；`part1` sec 4 引的行号 `:192` 与 HEAD blob 的 `:198` 差 6（本腿 M-b 实测 `:198`）。**属读数过期**，NOT 判据错。
- (iii) `226_test.go` 那句原注释仍是假话（`nothing polls` 逐字在场、blob `38b32a55`），`ecebc7a9` 用**追加更正**处理、原句 NOT 抹 —— 与票面要求一的**字面**口径有缝（票面写"改准"，它选"追加更正、原句保留"）。本腿判：**实质满足**（更正紧贴原句、逐字说明从 `248095d1` 起为假、并指到新用例），**字面**是否算"改准"请编排者拍；NOT 据此退件。
- (iv) `part2` sec 7 第 3 条它自报"226 写回不 Adopt 内存那枚没做、若仍要属另一格"= 与本件 (10) 同结论，**NOT 欠账**。

## sec 3 "票面禁手有没有犯"的明确判决（含我自己那一发）

**判决：没犯。** 三把尺各自独立：
- **间隔**：`ecebc7a9` 的 `.go` diff 里 `configReloadPollInterval` / `time.Second` / `Millisecond` 的新增或删除**只出现在注释与两处等待上界**，生产码 `config_reload.go` 未被本笔触碰（文件集里根本没有它）。
- **sleep**：新增件里没有 `time.Sleep`；两处等待是 `time.NewTimer(40s)` + `20ms time.NewTicker` 的**有界上界**（失败上界，NOT 成功条件），注释与代码一致。
- **我自己的那一发（反向证伪）**：tick 收紧到 **10ms** => 旧用例红、新用例绿（`logs/52-mc-*`）。若新用例是靠"收口窗口比 tick 短"活着，收紧 tick 只会让它更容易红；它没红 => 它 NOT 是时序侥幸。另一向放宽到 **10s** => 仍绿（`logs/42-ma-*`）。**两向都绿 = 与 tick 间隔无关。**

## sec 4 "改前读数"那一形我独立复核的结果

**同向。** 读数逐字见 (5)：`hold=0` => `theme="dark"` / `applied_lines=0` / span `50ms` / 绿；等 tick 行 => `theme="light"` / `applied_lines=1` / `tick_wait=1.003s` / span `1.027s` / 会在 `:102`（现 `:118`）转红。
两臂同码、同一个生产 1s tick、只差"受守卫写落账后再活多久"；臂 B 用**有界等待产品自己的 `state=applied` 行**代替 `sleep 2.5s`，故读数 NOT "睡够了碰巧"，而是"tick 说话之后必然"。
与它报的三对数字**方向、量级、applied_lines 0/1 全部一致**；数值差（`50ms` 对 `84/95/99ms`）见 (5) 的具名，属区间量、NOT 矛盾。

## sec 5 导出树整包读数我复跑的结果

**成立：新增红 0。** 两本都在包尾（最后一条 RUN 有判决行、`reached = compiled` = 347/348），四数 `CONTROL 472/321/21/5` 对 `AFTER 473/324/19/5`、`NEW RED` 空、`FLIPPED GREEN` 两枚面板窗口族、`shared red 19`（逐名见 (8)）。
与它报的 `472/322/20/5` 对 `473/324/19/5` 只在对照侧的 `322/20` 与我的 `321/21` 差两枚（`TestPanelHostLatencyPercentilesAC2` / `TestPanelHostRealWindowHopAndLifecycle`），**新增红与翻绿集合 NOT 因此变化**，且那两枚在它自己的对照本里也时红时绿（争用型）。=> **"新增红 0"这个结论我独立复跑成立。**
另：主检出不可判已复现（见 (7)），本腿**没有**拿主树"0 枚 FAIL"当整包绿。

## sec 6 门禁 rc 与末行逐字

见 (6) 的表（五门 rc 全 **0**；末行逐字照抄，NOT 转述）。格式门禁按裁定件**甲乙两尺都报**：甲（入库 blob，判据）= 两枚目标件空 listing、`.scratch` 之外 0 枚；乙（工作树，诊断）= 4 枚 CRLF 假红、blob 全清。
`attrib.sh --tracked-only`（CI 那步）本机 rc=1 / 35 枚（31 在 `.scratch` + 4 枚 CRLF 假红）= 与 `ecebc7a9` 无关，登记为"该步本机不可判"。

## sec 7 红名册逐名

见 (8)（19 枚逐名 + 三枚争用型面板红的归属）。**新增红 = 0 枚**；`always_write_no_clobber_226/227` 两枚在 `AFTER` 本里分别 `PASS (2.48s)` / `PASS (3.12s)`。

## sec 8 提交号 + 逐枚 `.go` / 越界 + 正控

- 被裁提交：**`ecebc7a9`**（`票 227 落地腿 227-r1：真轮询之下手改仍看得见（AC#6）补一枚常驻用例`，2026-10-22 14:42:14 +0800），一笔四件，`.go` **2** 枚（`always_write_no_clobber_226_test.go` M +16/-0、`always_write_no_clobber_227_test.go` A 225 行），交付件两枚（138 + 59 行），**越界 0**、生产码零字节。
- 严格前缀名册：`^票 227 落地腿 227-r1` = 9 笔，逐笔越界 **0**、`.go` 合计 4 枚次。
- 正控：同尺打 `ecebc7a9` 回验 `.go` = **2**（尺活）；负控：同尺打 `9fc2be02` = **0**。
- 本腿提交：`git log -- docs/evidence/s1/227-v1-verdict.md`（本件入库笔；件内 NOT 写自身哈希，写死在件里必然过期）。本腿 NOT push。

## sec 9 未定义即停 / 交编排者拍

1. **"改准那句话"的字面口径**（追加更正 对 就地改写原句）—— 见 sec 2 (iii)，本腿判实质满足、字面请拍。
2. **`226-v2` sec 6 第 8 条那枚该落哪** —— 见 (10)，本腿判"另一格 / 另立票"，与编排者 `9fc2be02` 第 5 条同向。
3. **`attrib.sh --tracked-only` 本机 35 枚红** —— NOT 属本票，但它是 CI 的格式门禁步骤；31 枚 `.scratch` 台件 + 4 枚 CRLF 假红，建议另立票或并入格式门禁裁定件的后续。
