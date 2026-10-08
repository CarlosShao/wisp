# 票 35 `:49` Backpressure — 35-r8 落地腿（只写测试件）交付

日期 2026-10-08。分支 `dev`。起手锚 `028f2e24`（派单声称 `1b8f9c4d`，复跑已漂，见 `logs/anchor-git.txt`），本腿第 1 笔 `fcf4af38`。
射程＝票面 `:49` 那格**只剩的一味**：`no unbounded memory (heap cap asserted)`。⛔ 本腿未动 `internal/panel/**` 任何非 test 文件一个字节，⛔ 未碰 `cmd/wisp/**`（只做了两发只读 grep，见 §⑦），⛔ 未翻页面任何框，零 push。

**一句话结论：这一格今天没有闭合，缺的是产码不是仪器。** 仪器补上了（1 枚绿尺 + 1 枚正控 + 5 发仓外突变各自咬不同相），而它量出来的读数是：**洪水在"键数"维度有界、在"每行字节"和"命名账本"两个维度无界**。两枚红尺是诚实交件，不是放宽后的绿。

---

## ① 查重（三种情况分开列，⛔ 未查重就不新立尺）

尺面命令与落点（件 `logs/r8-run-first.txt` 之前的名册扫）：

**情况 A＝已有名册钉，但钉的是行数不是字节数（抵不了这一格）**
- `internal/panel/pump_test.go:165 TestTheStreamLogTruncatesInsteadOfMerging`：`:211` 逐字 `if len(got) > 2*StreamKeyHardCeilingMultiple` ⇒ **对常量比行数**。16000 枚键灌进去它照样绿，因为上界说的是"几行"，不是"多少字节住在内存里"。`:219` 钉的是 `len(DroppedKeys())` 的**条数**。
- `internal/panel/subagent_roster_197_test.go:340`/`:344`/`:417`/`:421`：钉截断**记账**（Truncated/ElidedRunes/DroppedKeys 与 wire 字段一致），同样零内存读数。

**情况 B＝已有内存上界断言：0 枚**（本腿没找到，也没冒充找到）

**情况 C＝根本没有（这一味今天的真实状态）**
- `grep -rn --include='*.go' -E 'HeapAlloc|runtime\.MemStats|MemStats|cap\(' internal/panel` ⇒ **rc=1（零命中）**。整个包从没读过一次运行时内存。⇒ "heap cap asserted" 半句**零尺**，本件就是那把尺。
- 票面 `:27` 那句"每面板 bounded queue"在 `internal/panel` 里**没有对应产码**：`make(chan` 在该包非 test 文件里 0 命中（出向队列住在 `cmd/wisp`，见 §⑦）。

## ② 新件

`internal/panel/backpressure_bound_35r8_test.go` — `wc -l` = **371**（`logs/gate-cr.txt` 同记 CR=0）。新建，无同名件；包内既有 34 枚文件（`ls internal/panel`）里没有任何 `*_35r8*`。
4 枚用例（**A 绿 / B 红 / C 红 / D 正控绿**）：
- `TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8`（绿）
- `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`（**红**）
- `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`（**红**）
- `TestFloodRulerFiresOnAnUnboundedSink35r8`（正控，绿）

洪水是**行为**不是字面：`flood35r8` 真把 keys×perKey 枚 delta 灌进 `StreamLog.Append`，同时在 `NewSnapshotPump(PumpSources{Results, Out})` 上真发布快照，`Out` 往一条**永无人读**的 `page chan []ResultChunk` 做阻塞发（10ms 超时），并且当场要求：`attempts>0`、`Publishes()==attempts`、`refused==attempts`、非阻塞 select 读不到东西 ⇒ "消费者被堵住"是被量的前提，不是注释。判据 `overCap35r8` 的四项（行数／钉住的文本 runes／钉住的命名 runes／活堆字节）**每一项都不含洪水规模**，只有 `retentionCap35r8()` 的上界常量。

现量读数（件 `logs/r8-run-first.txt`，逐字）：
```
A: 1x text=35072 rows=64 names=336 heap+438416B | 4x text=35072 rows=64 names=1536 heap+479120B | cap text=36864 rows=64
B: 1x pinned 1200000 text runes over 3 rows (live heap +1162128 bytes), 4x pinned 4800000 runes over 3 rows (live heap +4833536 bytes) ... Retention grew 4.0x for a 4x flood
C: 1x opened 4000 streams and pins 3936 names (50058 runes, live heap +116176 bytes), 4x opened 16000 and pins 15936 names (211994 runes, live heap +564144 bytes) - rows stayed at the ceiling 64 in both
D: the ruler stayed silent... -> fired as required: rows 100 > the hard ceiling 64; pinned text 800000 runes > the cap 36864 (heap +787152 bytes)
```
A 那行是本腿最想给的数：**400 路与 1600 路洪水（15.4M 输入 runes vs 61.4M）钉住的文本一字不差都是 35072 runes、都是 64 行** ⇒ 扇出维度确有上界，且它是量出来的不是认字符串认出来的。

## ③ 正控是什么

`TestFloodRulerFiresOnAnUnboundedSink35r8`：仓内常驻、CI 每次都跑。
- 造一枚测试本地的 `unboundedSink35r8`（**同形状、把上界摘掉**：不 clamp、不 drop、不命名），走**同一个 `flood35r8` 夹具、同一个 `overCap35r8` 判据**；
- 断言判据必须回答 **OVER cap**（"种 X 必响"）；
- 断言该 sink 真钉着 `100×50×160 = 800000 runes / 100 rows` ⇒ "响"不可能来自一把不再量东西的尺；
- 再把**同一形状**灌进产码 `StreamLog`，断言它 within cap ⇒ 那声"响"鉴别的是**缺失的上界**，不是夹具自己的体积。
把 `overCap35r8` 中和成 `return ""`、或把 `Measure()` 中和成返回 0，这枚先红（M 表外另由 §⑤ 具名其射程）。

## ④ 突变表（5 发，全部仓外 `D:/tmp/wisp35r8/mut/<M>/pump.go` + `go test -overlay`，驱动 `logs/mutation-driver.txt`，逐发输出 `logs/mut-M*.txt` 各落 `[rc=N]`）

| 发 | 改了哪一形（只改副本） | 哪枚用例由绿转红 | 红句逐字（件可查） | rc |
|---|---|---|---|---|
| M1 | `enforceBoundLocked`：`if len(s.order) <= s.maxKeys {` → `if true {` ⇒ **键数上界整个摘掉** | A（另 B/C/D 同跑红） | `backpressure_bound_35r8_test.go:246: 400 streams x 60 deltas x 160 runes through a 32-key log is not bounded: rows 400 > the hard ceiling 64; pinned text 3840000 runes > the cap 36864; live heap + 3959768 bytes > the cap 1196032 (flooded 3840000 runes in)` | 1 |
| M2 | `clampLocked` 函数头插 `if true { return }` ⇒ **每行 runes 上界摘掉**，键数界留着 | A（D 亦红） | `:246: ... not bounded: pinned text 614400 runes > the cap 36864 (flooded 3840000 runes in)` | 1 |
| M3 | `Publish()` 拒收分支的 `p.count()` 删掉 ⇒ **blocked consumer 那一路静默丢弃而不计数** | A（D 亦红，夹具级） | `:241: blocked consumer: the pump counted 0 publishes for 24 attempts - a refused exit that stops being counted is a silent drop` | 1 |
| M4 | `enforceBoundLocked`：`s.dropped = append(s.dropped, oldest)` 注释掉 ⇒ **计数出口挪出截断路径（丢了但不点名）** | A | `:273: dropped naming = 0 keys, want 1536: every key the log stopped tracking must be named` | 1 |
| M5 | `Append` 头部插 `if true { return }` ⇒ **什么都不存的退化"有界"**（专测两枚红会不会被白判通过） | A；B/C **仍红**（被各自 measurement-sanity `t.Fatalf` 接住） | A：`:259: rows = 0, want exactly the hard ceiling 64: fewer rows than the ceiling means the flood was dropped rather than bounded`；B/C 落点件 `logs/mut-M5.txt` | 1 |
| 还原 | 跑完即弃，仓内零写入 | — | `logs/blob-check.txt`：`pump.go wt=645e1f02…a18 head=645e1f02…a18 same=YES`、`pump_test.go wt=3f84c336…e7a4 head=…e7a4 YES`、`subagent_roster_197_test.go wt=f803179d…3b1d head=…3b1d YES`；`git status --porcelain -- internal/panel cmd/wisp` → 只有 `?? internal/panel/backpressure_bound_35r8_test.go`（rc=0） | 0 |

⛔ 未打过 merge 方向的任何形（派单禁区）：`pump.go:397` 那行逐字 `// Overflow TRUNCATES and never merges (ticket 197 leg B re-cut this rule, which` **本腿只读未动**，也未写过任何"让它合并"的突变。

## ⑤ 恒真面（哪一发怎么坏都照绿，具名）

1. **heap 那一项带 1 MiB 常数余量，弱。** M2 的红句里只有 runes 项开火（614400 runes 的钉住量≈0.6 MB < heap cap 1196032）⇒ **retention 低于约 1.2 MB 的坏形是 runes 项在咬，不是 heap 项**。heap 项单独使用不成立，别把它当那味尺。
2. **正控 D 不是全相尺。** M4（丢了不点名）与 M5（什么都不存）之下 D 都绿：它第三块"同形灌进产码 must within cap"会被"什么都不存"白判通过 ⇒ **反真空的担子全在 A 的 `rows == ceiling` + 每行内容 + `truncated` marker + 三枚计数出口上**，不在正控上。D 只对"判据本身失效/夹具停量"这一相有鉴别力。
3. **出口那一路的内存增长本尺不量。** `flood35r8` 的发布次数被 `maxPublishAttempts = 24` 钉住（不随洪水放大），它证的是"消费者真被堵住 + 拒收仍计数"，⛔ 不证 `Publish()` 阻塞时调用方手里那份快照字节会不会堆起来——那一相今天要靠产码有界队列才有被测物（`internal/panel` 里没有）。
4. **两枚红尺的红句是 `t.Errorf` 打的读数**，"会绿"的唯一条件是产码长出真上界；把洪水调小**不会**让它们绿——A 与 B/C 的 sanity 钉（`t.Fatalf`，逐字 "measurement sanity: ... the harness stopped measuring the flood" / "...the ledger being measured is not the one the log writes"）会先响，M5 已实证。
5. `overCap35r8` 的命名账本项**借用 text cap**（`cap.textRunes`）当额度，不是独立推出的上界；若机主另裁一个命名预算常量，这一项要跟着改，本腿不自造。

## ⑥ 门禁逐条（每件自落 rc；测 rc 那句前面⛔ 无管道）

| 尺 | 命令（逐字） | rc | 件 |
|---|---|---|---|
| vet | `go vet ./internal/panel/` | **0**（件 7 字节＝只有 `[rc=0]`，零输出＝通过） | `logs/gate-vet.txt` |
| test 前 | `go test ./internal/panel/ -count=1`（写入被审件**之前**） | 1（4 枚既有红） | `logs/gate-test-before.txt` |
| test 后 | 同命令（收尾） | 1 | `logs/gate-test-after.txt` |
| 名册作差 | `diff fail-before.txt fail-after.txt` | **1＝差 2 行，且只多了本腿那两枚红**（下表） | `logs/fail-roster-diff.txt` |
| d22scan | `sh scripts/d22scan.sh` | **0**，尾行 `d22scan: clean - no D22 ban violations`（⛔ 未在仓根跑 `go run ./tools/d22scan`） | `logs/gate-d22scan.txt` |
| gofmt | `gofmt -l internal/panel` | **0**（件 7 字节＝零输出） | `logs/gate-gofmt.txt` |
| 行尾符 | `tr -cd '\r' < 件 \| wc -c`（⛔ 未用 `grep -c $'\r'`） | 新件 CR=0，`logs/anchor.md`/driver/after 件 CR=0 | `logs/gate-cr.txt` |
| HEAD blob 对拉 | `git show HEAD:<件> \| tr -cd '\r' \| wc -c` | 见 §④ 还原行（YES 三枚）+ 本件 commit 后复拉 `logs/final-check.txt` | `logs/blob-check.txt` |

新增红具名（**本腿 2 枚，全是故意的红，⛔ 不是放宽来的绿，也没藏**）：`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`、`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`。
其余 4 枚红（`TestApprovalCardViewJSONKeysMatchFrontendTypes`/`TestComposerContractTypesMatchFrontend`/`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`/`TestC21DesignTokensFourWayAgree`）起手就在（`35-v7` verdict §Q4 同名同色），不归本程，本腿一枚没碰。

## ⑦ 这一格今天到底闭合没有／缺的是产码还是仪器

**未闭合。仪器这一味今天补上了；缺的是产码。**

- 缺仪器那半句**已结**：包内第一次有了读 `runtime.MemStats` 的尺，第一次把"灌 N 千条之后那个数还在不在涨"量出来（A：15.4M vs 61.4M 输入 ⇒ 钉住量 35072 runes 一字不差；B：1.2M vs 4.8M runes 输入 ⇒ 钉住量 4.0x 同步涨到 4.83 MB）。5 发突变各自咬到不同相（§④）。
- 缺**产码**那半句：`pump.go` 的上界只挂在**键数**一维——`Append` 只在 `s.truncated` 为真时才 `clampLocked`，而 `truncated` 只由 `enforceBoundLocked` 在键数过 `maxKeys` 时置位。⇒ 3 条流各灌 2000-rune×800 段（一次长回答，最常见的形状）**没有任何字节上界**；过界之后 `dropped` 账本又按事件数线性长。要绿必须二选一或并行动产码：(a) 每行 runes 在键数之外也设界（丢失照 `ElidedRunes` 记账点名），(b) 命名账本设预算（计数 + 样本名 / 环形保留），(c) 出口在消费者堵住时拒收洪水。**本腿一律不代裁、不动 `pump.go` 一字。**
- ⛔ 不许用"行数上界已有尺"抵这一格：那是 `pump_test.go:211` 对常量的比法，与 AC 的 "heap cap asserted" 不是同一味。
- 只读复认（⛔ 未写、未跑该包任何测试，按派单约束整包 `cmd/wisp` 不跑）：`cmd/wisp/panel_resident_windows.go:148 tasks: make(chan func(), 16)`、`:130 failedPost atomic.Int64`、`:369 rp.failedPost.Add(1)` 三枚**行号今天未漂**；`cmd/wisp` 内 `make(chan func()` 的测试命中只有 `panel_host_windows_test.go:70`（**夹具装配，不是上界断言**），`failedPost` 的测试命中 **0** ⇒ 普查腿"零测试引用"只对 `failedPost` 成立，对 chan 字面量需更正为"只出现在一枚夹具里"。

## ⑧ 我没量到的（⛔ 不装成量过）

1. **`cmd/wisp` 那侧队列的字节维度**：`chan func(), 16` 是**任务数**上界，本腿未量它实际钉住多少字节，也没测它的红/绿（写权禁、整包不跑，按约束归编排者补）。
2. **物理内存（RSS）**：本尺读 `HeapAlloc`（活对象），`debug.FreeOSMemory`/`Sys`/`RSS` 一未用 ⇒ "进程实际占多少物理页"没有读数。被 clamp 掉的字符串变垃圾后仍在 Go arena 里时，本尺看不见。
3. **`internal/panel` 里除 `StreamLog` 之外是否还有别的累加点**：本腿只按名册词（bounded/queue/make(chan）扫了 `pump.go`，`composer*.go`/`bridge.go` 的 map 累加面未逐一量。若还有第二个洪水蓄水池，本件看不见。
4. **真页面/WebView2 那一跳**（需 `winlive` 面）：出向推载荷今天在本仓仍无被测物（`35-v7` verdict §Q3 同判），故"堵住的是真消费者"只用无读者的 channel 建模，未用真窗。
5. **上游**：`internal/agent`/`internal/tools` 的洪水产生速率与背压（射程外）。
6. **票面 `:27` "每面板 bounded queue"** 的兑现度：本腿只证明 `internal/panel` 没有队列，未证明 `cmd/wisp` 那枚 16 深队列是不是票面要的那一枚。

## 顶回/更正派单（逐条具名，均为本腿现量）

1. **HEAD 号漂**：派单 `1b8f9c4d` → 起手复跑 `028f2e24`（`logs/anchor-git.txt`）。本腿按复跑走。
2. **`:49`/`:392` 行号锚今天未漂**：判据行仍在 `:49`、〔10-08 11:3x 编排者注〕在 `:51`、普查表 `:49` 那一行在 `:392`、`pump.go:397` 逐字命中、三枚计数出口逐枚复认仍在 `:632`(`Truncated`)/`:644`(`ElidedRunes`)/`:657`(`DroppedKeys`)。⇒ 这一格派单给的锚点本腿**全部复认成功，无需更正**；但普查那句"每面板推队列没写"需按 §⑦ 最后一条补一个限定（chan 字面量在 `cmd/wisp` 有 1 枚夹具引用）。
3. **派单 §3 的假设"上界可能在 cmd/wisp 那侧"成立，且比它写的更糟**：`internal/panel` 不是"没有队列"这么简单——它有一个**按键数**的上界，而那恰恰不是内存上界（B/C 两枚红就是这个区别）。如果编排者要把这一格记成"部分在"，账要写成"**键数在、字节不在**"，⛔ 别写成"有界队列在、只缺测试"。
4. **本腿自报一处初稿缺陷已修**：第一次 Write 把 §⑦ 的 ⛔ 字符写进了注释（d22scan 注释豁免、但本仓零 emoji 规范的射程是 U+2600–U+27BF，`⛔`=U+26D4 在仪器带内）。已改为 ASCII 文字后 d22scan rc=0；现量件里⛔不再出现（`grep -c '⛔' internal/panel/backpressure_bound_35r8_test.go` → 0）。
