# 票 358 终裁腿 358-v1 - 判语件 part 3/3（非实现者，裁决者 NOT 实现者）

前件＝`docs/evidence/s1/358-v1-verdict-part2.md`（part 2/3：sec 3 `AC#2`、sec 4 `AC#5`、sec 5 票 355 钉子）。

## 6 门禁复跑（本腿现跑，逐把 rc ＋末行逐字）

### 6.1 六把门
| 门 | 命令 | rc | 末行 / 读数逐字 |
|---|---|---|---|
| G1 | `GOFLAGS= go build ./...` | **0** | 无输出 |
| G2 | `go vet ./internal/agent/ ./cmd/wisp/` | **0** | 无输出 |
| G3 | `go test -count=1 ./internal/agent/` | **0** | `ok  github.com/CarlosShao/wisp/internal/agent  4.905s` |
| G4 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; ... ban #8 internal/=610, ban #8 cmd/=145; ...` |
| G5 | `sh scripts/check-path-length-budget.sh` | **0** | `... unit cross-check: the byte view and the character view select the SAME 60 paths on this tree, recomputed just now in both units` |
| G6 | `gofumpt -l`（绝对路径）对 5 枚涉改件 | **0** | 列表**空** |

G3 四数分开报（本腿现跑，未相减）：顶层 `RUN 102` / `PASS 102` / `FAIL 0` / `SKIP 0`
（RUN 与终态名册两向皆空）。G4 分母 `cmd/=145` 本腿独立复算：`git ls-tree -r --name-only HEAD cmd/`
里 `.go` 共 **146**，去掉 `cmd/wisp/testdata/esclistener/main.go` 后 **145**（d22scan ban#8 跳 testdata），
pre-358 同尺＝**144** => **144=>145 成立**（新增那枚 `compression_trace_358_test.go`）。
`gofumpt` 用绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`；未用 `CGO_ENABLED=0`。

### 6.2 全族同发（355＋358）
`go test -count=1 -run 'CompressionTrace' -v ./cmd/wisp/` => **rc=0**，
`ok github.com/CarlosShao/wisp/cmd/wisp 0.138s`。顶层 RUN 6 / PASS 6 / FAIL 0 / SKIP 0；
子测试 PASS 4 / FAIL 0。逐名：`ThreeFlavoursReachTheClosingLine` PASS、
`SilentWhenNothingCompressed` PASS、`ReachesStdout` PASS、`CompactedLineIsTheOne355Shipped` PASS、
`AttemptedSaysItFoldedNothing` PASS（两子）、`SkippedStaysSilent` PASS（两子）。

### 6.3 整包 `cmd/wisp` - **rc=1，但本腿判它 NOT 能当判据，也 NOT 是 358 的伤**
本腿现跑两发（**同尺、带 `-v`、同一时段**），两发都抬 `-timeout 30m`：

```
本树  : go test -count=1 -timeout 30m -v ./cmd/wisp/   => rc=1
        顶层 RUN 159 / PASS 158 / FAIL 0 / SKIP 0 ；包末行 `FAIL  github.com/CarlosShao/wisp/cmd/wisp  344.846s`
        RUN 有名无终态者恰 1 枚：TestAC14GoSideEvalPushReachesThePage
pre-358: 同命令 + -overlay 把 run.go / 355 test 换回 5ca2472e blob、删 358 test  => rc=1
        顶层 RUN 152 / PASS 151 / FAIL 0 / SKIP 0 ；包末行 `FAIL  github.com/CarlosShao/wisp/cmd/wisp  342.309s`
        RUN 有名无终态者恰 1 枚：TestAC13BringUpSurvivesAReusedThreadQuit
```
=> **两发 `^--- FAIL` 都是 0**，死点都是“进程被吃掉一枚”、且**两发死点不同名**（印证 `A1100` 那族会漂）。
本腿**按名作差**（join 两发终态名册）：
```
NEW RED (pre 非 FAIL、post FAIL):   (空)
RECOVERED (pre FAIL、post 非 FAIL): (空)
ONLY IN PRE  (名消失):              (空)
ONLY IN POST (名新增):  7 枚，全部 PASS
   TestAC13BringUpRefusesAThreadWithAQueuedClose / TestAC13BringUpSurvivesAReusedThreadQuit /
   TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever /
   TestAC14AwaitedBindingReplyReachesThePage /
   TestCompressionTraceAttemptedSaysItFoldedNothing /
   TestCompressionTraceCompactedLineIsTheOne355Shipped /
   TestCompressionTraceSkippedStaysSilent
```
=> **新增红 0 枚**（两向皆空）。那 4 枚 `TestAC13`/`TestAC14` 是“pre 发死在它们之前、post 发跑到了它们”
（全 PASS），NOT 新增红；那 3 枚 358 用例是本票新增（全 PASS）。整包 **362 枚编译的顶层用例只有
159 枚跑到**（`grep -hcE '^func Test' cmd/wisp/*_test.go` 合计），未跑到那批**本腿无法用整包这一发判色，
属票 314 射程**。

**自查一次污染（记在这里，NOT 藏）**：本腿第一发整包跑在 18:21-18:35，那一发 **rc=1 且 11 枚真红**，
末行还跑到 `TestZZ315A1*`。复量发现：那 11 枚红句全是
`parse run.go: run.go:1356:2: expected '}', found 'case' (and 1 more errors)` -
**是本腿自己的变异 B 落地窗口**（那一发跑着时我正在改 `run.go`，而这族用例在运行时**从盘上 parse `run.go`**）。
=> **那一发作废、NOT 当读数**；6.3 是**全部变异还原后**重跑的干净发。
教训（建议登台账）：`cmd/wisp` 里有一族用例**运行时读盘上的 `run.go`**，故“边跑整包边改 `run.go`”会自污染，
变异读数必须与整包门禁**串行、NOT 并行**。另：该次并发腿同时改了 9 枚 `cmd/wisp/*_test.go`
（`t.TempDir()`=>`sealableTempDir124(t)`），本腿的干净发与 pre-358 基线**都用 `-overlay` 把这 9 枚
钉回 HEAD blob**，故两发同尺、差分只反映 358。

## 7 红名册（逐名，NOT 只报枚数）

**新增红＝0 枚。** 今天整包**两发都无 `^--- FAIL`**，故没有需要逐名判归属的新红。
SKIP 3 枚（NOT 红，登记在此）：`Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable`、
`Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots`、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`
（两树同 SKIP）。被死点吃掉、**无法判色**的具名：本树 `TestAC14GoSideEvalPushReachesThePage`、
pre-358 `TestAC13BringUpSurvivesAReusedThreadQuit` - 两枚都是台账在册真窗族
（`A1079`/`A1100`/票 314），**非新增**。

## 8 写面与纪律（逐笔 `.go` / 越界 ＋正控）

派单白名单（`358r1-dispatch.txt` 写面节＋`358-ruling-nail-355.txt` 边界节）＝
`internal/agent/{compress,loop}.go`、`cmd/wisp/**`、`docs/evidence/s1/358-r1-*.md`。

| 提交 | 件 | `.go` | 越界 |
|---|---|---|---|
| `05f76ad6` | 5（`compress.go`/`loop.go`/两新 test/impl.md） | **4** | **2**（见下） |
| `3ed3fdce` | 1（impl.md） | 0 | 0 |
| `43f44f59` | 3（`run.go`/355 test/358 test） | **3** | 0 |
| `a8ae09ce`/`10214eec`/`45f0b72f`/`cfa8bff2` | 各 1（part2.md） | 0 | 0 |

**那 2 枚“越界”本腿判它不是越界**：`compress_state_358_test.go`、`compress_trigger_358_test.go`
严格按派单字面（白名单只列 `compress.go`/`loop.go`）＝名单外；**但**派单同一句要求“每一处行为改动配
一枚今天响的用例”＋`AC#5` 要求“真值/估算两形各一枚用例” => **不新起测试件就交不出这格**，
且新件在 `internal/agent/**`（本波该腿独占面）。=> 判：**派单白名单漏列测试件**（派单缺陷），
实现腿做法**合意图**；本腿**NOT 据此判作废**，但**登记这处白名单不足**。

**正控/负控（证明尺不是恒 0）**：同尺打 `b676092d`（票 355 落地腿）=> `.go`=**3**；
打 `3ed3fdce`（纯交件）=> `.go`=**0**。两向都可读 => “`.go` 4 / `.go` 3 / 越界 0”不是瞎尺读数。
**新件唯一性**：`git ls-tree -r --name-only HEAD | grep -c` 逐枚 => 三枚新件各 **1**。
**未 push**：起手 `git status -sb` 为 `## dev...cnb/dev [ahead 1002]`，收线时仍只涨本地、远端未动。
**被禁码位自查**：两份交付件 `grep -P` 六段全部 **0 命中**；七枚涉改 `.go` 里 `run.go` 有 6 行命中
（`:207`/`:285`/`:554`/`:1167`/`:1169`/`:1171`）- 本腿复量 **pre-358 同文件同为 6 行**（全在既有注释里，
是 `(1)`/`(2)`/`(3)`/`定案(2)` 这类**原本就在的**带圈字符，NOT 358 新增），故 NOT 本票引入；
`d22scan` rc=0 亦未报。

## 9 `AC#4` 处置建议 - **该现在摆机主，但只摆“文字”、NOT 摆“要不要做”**

`AC#4` 要求 D15 那段文字同步。本腿现读：`docs/PLAN.md` 与 `docs/specs/**` 在 `5ca2472e..HEAD`
**零字节**（`git diff --stat` 空），即实现腿**确实没碰契约文字**（合纪律）。而实现已把触发基数换成
provider 真值，D15 今天仍逐字写“会话历史 > 12000 token => 对最老轮次压缩” => **文字与已落地的判据已分叉**。

**建议：现在就摆，且按“改契约文字”那一类摆（铁律第 4 条三类之一），NOT 问“要不要改”。**
理由：(1) 分叉已存在，不摆就是让规格与仪器继续分叉（`A201` 第二点那族旧病）；(2) 改动面小 - 只需把
“12000 token”改成“按所配 `context_window` 等比缩放后的 `HistoryCompressTokens`，分子优先取 provider
回报的 `Usage.InputTokens`、无回报时回落本地估算并在痕中标明”；(3) 阈值本身**一字节未动**，
故这不是“改预算”，是“把已经落地的判据写进文字”。
**NOT 本腿代拟契约文字**（那是机主批的事）；建议由编排者把改前/改后逐字对照摆机主，一步过。

## 10 逐格结论（能翻 / NOT 翻 ＋缺什么）

| 格 | 结论 | 缺什么 |
|---|---|---|
| `AC#0` | **能翻** | 无。三态读数＋正负控＋“改前形(1)(3)同形”本腿实测复现齐备 |
| `AC#2` | **能翻** | 无。三发变异逐发转红、逐发还原核 sha256、compacted 逐字节独立核对 |
| `AC#5` | **能翻** | 无（票面三条判据逐条满足）。附带一枚**残余面**记在 sec 4.3，属后续加固、NOT 挡本格 |
| `AC#4` | **NOT 翻**（属机主批的契约文字，归编排者摆） | 缺“机主批”这一步；不缺任何产物。见 sec 9 建议 |
| （票面前提一处数错） | 不影响任何一格 | sec 1：票面 `:180-183` 引的是 `b23c7f77` 坐标；实现腿报的 `:343-347` 是注释不是 `break` |

## 11 本腿读数的落点（可复跑）

- 快照与探针：`D:\tmp\wisp\20261021\358v1\`（`m0\` 三枚 M0 快照、`pre358\` 四枚 pre-358 blob、
  `headpin\` 9 枚并发腿改过的 test 的 HEAD blob、`zz358v1a/b/c_probe_test.go` 三枚探针、`ov-*.json` 五份 overlay）。
- 整包两发全量 `-v` 输出：`cmdwisp-clean-v.txt`（本树）、`cmdwisp-pre358-v.txt`（基线）、
  `cmdwisp-whole-v.txt`（**作废那一发**，见 sec 6.3 污染记录，留档不删）。
- 门禁与全族读数：现场读数已逐字贴在 sec 6，未另存。
