# 197-r2 流层证据件 — 每枚子代理一条自己的流，溢出不再"合并着吞"

写码程 `197-r2`（腿 B），派单：
`.scratch/wisp/dispatches/2026-09-28-184x-impl-197r1-entity-and-197r2-stream.md` §腿 B。
裁决与对抗验收由**另一个** agent 做（`SPEC-12 §4.3` #1/#3），本件只是实现者的自述。

---

## 0. 起手锚点与名册（现读，未采信派单里的任何行号）

- `git rev-parse --short HEAD` = **`487ad096`**
- `git status --porcelain` = **79 枚**条目。不属于本程、本程不碰不提交的：
  `.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/flip-*`（7 枚）、
  `design/` 的 11 枚 `D` + 5 枚 `M`、`docs/reports/pending-and-issues.md`、
  `docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md`，以及腿 A 正在写的
  `internal/tools/subagent_197.go`（起手时还未出现，18:4x–18:5x 之间它落进工作树并**语法不通过**，见 §4）。
- 起手时 `go build ./...` **rc=0**（18:4x，`487ad096` 之上的树）。

### 派单读数复核（行号一律现读；漂移具名）

| 派单写的 | 现读（HEAD=`487ad096`，改前） | 现读（本程改后） |
|---|---|---|
| `StreamLog.Append` 约 `:301` | `pump.go:315` | `pump.go:403` |
| `Chunks()` 约 `:339` | `pump.go:353` | `pump.go:444` |
| `DefaultStreamKeys = 32` 约 `:275` | `pump.go:304` | `pump.go:363` |
| 溢出语义"合并不丢"约 `:287-290` | 注释在 `:289-293`，**实现在 `:366-380` `mergeOverflowLocked`** | 该函数已从仓里消失（`grep -rn mergeOverflow internal/ cmd/` = **0 命中**） |

漂移方向是 +14 与 +30 级；派单的数是对更早的树测的。**"约 `:287-290` 是溢出语义"这一条不完整**：
那四行只是**注释**，判决在 `mergeOverflowLocked`，改注释不改函数＝界面照旧说谎。

### 与派单预检不符的一枚（具名上报，不按下不表）

派单写：`DefaultStreamKeys`／`maxKeys` **在全仓测试里零命中 ⇒ 没有既有钉**。
标识符确实零命中，但**行为有钉**：`internal/panel/pump_test.go:165` 的
`TestTheStreamLogMergesInsteadOfDropping`（起手 `go test -count=1 ./internal/panel/` 时它是**绿**的）
逐条断言"越界折成两行、两行合起来仍含每条 delta"。**改溢出语义必然让它红**，所以本程不是"无人替你响"，
而是"有一枚钉正指向要重裁的那一事"。处置见 §3 第 6 条。

---

## 1. 选型：乙（显式标"已截断" + 每键自留头尾），为什么不选甲

**甲＝上限随活跃任务数走。** `StreamLog` 是一枚**文本载体**，它不知道、也无从核实"宿主现在有几枚活跃任务"：
名册在 `internal/tools`（腿 A 的地盘），panel 反向依赖 tools 是装反的依赖方向；
在 pump 里另记一份"活跃数"＝给同一个事实造第二个真相源，而第二枚裁决者正是本项目反对的形状（D22 双角色同一味道）。
更要紧的是：**甲到了上限仍然要回答"这时候两条流怎么办"**——不回答就是同一个谎言推后一轮。

**乙＝到上限时不并键、只削每条流自己的中段，并把削掉多少说出来。** 本票要的每一条不变量
（一枚子代理一行、每行的文字属于且只属于那一枚、越界这件事可见）**都能只在 `StreamLog` 内部表达**，
不需要新依赖、不需要第二份活跃数。附带好处：内存界从"行数"搬到"每行 rune 数"之后，
键集合的增长不再逼着载体撒谎——它逼的是一行自己的完整度，而那一条是**写在行里的**。

"把上限调大就完事"没选，也没做：`DefaultStreamKeys` 仍是 32（`pump.go:363`，值未动），
动的是越界之后发生什么。

**本程不需要腿 A 的代码**：只按 §0 的形状 `subagent:<taskID>` 写与读（`SubagentStreamKeyPrefix` `pump.go:379`、
`SubagentStreamKey` `pump.go:386`：前缀全小写、冒号分隔、`<taskID>` 逐字透传、空 id 返回 `""`）。

---

## 2. 落点（逐跳到 `文件:行号`，均为现读）

`internal/panel/pump.go`

| 落点 | 行 | 是什么 |
|---|---|---|
| `StreamLog` 段头注释重写 | `:280-307` | 把"MERGES 不 DROPS"改成"TRUNCATES 不 MERGES"，并具名写清为何子代理这一形是界面说谎（D30 同源，只是从显示侧进来） |
| `type StreamLog` | `:309-317` | 新增 `kept map[string]*streamKept` / `truncated bool` / `elided int` / `dropped []string` |
| `type streamKept` | `:328-333` | 每键自己的保留窗：头 `streamTruncateKeepRunes`、滑动尾同样长度、`seen` 累计 |
| `(*streamKept).absorb` | `:336-347` | 增量进尾窗，返回被挤出的 rune 数（＝这条流自己付的代价） |
| `(*streamKept).render` | `:350-357` | `头 + "...[truncated: N runes elided]..." + 尾`；文本只由状态渲染，绝不就地编辑字符串（所以流里出现同名标记也不会被误解析） |
| `DefaultStreamKeys = 32` | `:363` | **值一字未动**，含义从"开始合并的阈值"改为"开始截断的阈值" |
| `StreamKeyHardCeilingMultiple = 2` | `:368` | 硬顶倍数：**上限随键数走的最后一格** |
| `streamTruncateKeepRunes = 256` | `:373` | rune 不是 byte——本项目流的是中文，按字节切会把一个汉字劈成乱码 |
| `SubagentStreamKeyPrefix` / `SubagentStreamKey` | `:379` / `:386` | §0 写死的键形状，常量在本包一处具名 |
| `NewStreamLog` | `:394-399` | 初始化 `kept` |
| `Append` | `:403-423` | 新键：进表 → `enforceBoundLocked`；旧键：先写回表、再（截断态下）`clampLocked` |
| `Close` | `:426-441` | 同旧键/新键两支，未知键仍记一行 done（旧行为保留） |
| `Chunks` | `:444-455` | 未动：仍按到达序、每键一行 |
| `enforceBoundLocked` | `:462-482` | 越界＝置 `truncated` 并**逐键收进自己的头尾**；只有过 `2×maxKeys` 硬顶才**整键丢弃并具名**（`dropped`，`kept` 同步清） |
| `clampLocked` | `:485-514` | 幂等折叠：已有窗口只吸收新增 delta，**不会把自己的头吃掉、也不会叠第二个标记**（这是第一次实现踩到的坑，见 §3 第 4 条） |
| `Truncated` / `ElidedRunes` / `DroppedKeys` | `:517-526` / `:529-539` / `:542-552` | 宿主侧出口。**没有**上到线上契约（见下） |

**契约面刻意没动**：`ResultChunk` 仍是三枚 JSON 键（`composer.go:65-69`），
`Snapshot` 顶层仍四枚键，`frontend/**` 零写面。
现读理由：`TestComposerContractTypesMatchFrontend`（`composer_test.go:47-80`）把
`ResultChunk ↔ frontend/src/lib/panel.ts` 的 `ResultChunkView`（现读 `:122-126`，就 `correlationId/text/done` 三枚）
做双向对账——给 `ResultChunk` 加一枚 `truncated` 就会在那枚**已在册常红**的尺上再添一枚新错因，
而这枚腿改不动 `frontend/**`。所以截断事实本轮落两处：**行内标记**（面板今天就读得到）＋
**Go 侧出口**（宿主/日志读得到）；把它做成界面一等公民是 **197-r3 载体那一程**的活。

新增判据文件：`internal/panel/subagent_stream_197_test.go`
（`:35` 键形状／`:57` N＋1 枚不并流／`:146` 两枚同写互不污染／`:209` 截断那一支的尺／`:275` 硬顶丢弃要具名／
末尾 `WHAT THIS FILE DID NOT MEASURE` 5 条）。

---

## 3. 判据逐条：改前 / 改后读数

1. **起手基线**（`487ad096` 之上、本程产码之前，`go test -count=1 ./internal/panel/`）：
   红 3 枚，逐名 `TestComposerContractTypesMatchFrontend` / `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` /
   `TestC21DesignTokensFourWayAgree`；`TestTheStreamLogMergesInsteadOfDropping` **绿**（＝既有钉，见 §0）。
2. **改后同包**（`go test -count=1 -race`，再 `go test -count=1`，两次同一套红名，且**无 DATA RACE**）：
   红 4 枚＝上面 3 枚在册常红 + `TestSliceAAttachesNoHostAndNamesTheOpenWindowHops`。
   第 4 枚**不是本程的**，也不是面板的：它 `composer_dispatch_test.go:427` 扫全仓 native-host 命中，
   命中内容逐字是 `internal/tools/subagent_197.go:0:unparseable: ...:198:31: illegal character U+FF0C '，' (and 10 more errors)`
   ＝腿 A 工作树里的半落状态。本程不修、不 Skip、不进提交（`internal/tools/**` 是禁区）。
   ⚠ **这一枚是暂态，本程没有参与修复**：提交前最后一次跑（19:0x，HEAD 已推进到 `d195cfdf`）
   `go test -count=1 ./internal/panel/` 的红**回到在册那三名**，`TestSliceA…` 随腿 A 那枚文件转可解析而消失；
   §4 的 `go build ./...` rc=1 也是同一原因、同一时刻的暂态。写在这里是让"改前/改后"各带时刻，
   不要让读者以为流层把 native-host 那一格修绿了。
3. **本程 6 枚用例（新 5 ＋ 重写 1）改后全绿**：`TestSubagentStreamKeyIsTheOneSpelling` /
   `TestNPlusOneSubagentsEachGetTheirOwnStream` / `TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther` /
   `TestOverflowTruncatesEachStreamsOwnHeadAndTail` / `TestHardCeilingDropsWholeKeysAndNamesThem` /
   `TestTheStreamLogTruncatesInsteadOfMerging`。
4. **实现过程中真红过一次的那一发**（记进证据，不是装饰）：第一版 `clampLocked` 不幂等，
   每次新键越界都对全体重跑折叠，**把自己的头吃掉并把标记套进标记**；
   读数：`ElidedRunes = 1251, want 1176`、三行都"has no explicit truncation marker"。
   第二发是写回顺序（`clampLocked` 之后又被 `s.chunks[key] = cur` 覆盖原文）：
   `row "subagent:s3" ... 0 truncation markers` / `ElidedRunes = 1206, want 1216`。
   现值：`s1/s2 = 404`、`s3 = 408`、总 `1216`，与 `produced - 2×256` 逐行对得上，且每行标记**恰好 1 枚**。
5. **正控两发（`go test -overlay`，仓内文件未动，突变体在 `$TMPDIR`）**：
   - **M1 关掉截断**（`clampLocked` 阈值改成恒不折叠）⇒ `TestOverflowTruncatesEachStreamsOwnHeadAndTail` **红**，
     逐行 "marker ... appears 0 times, want exactly 1"；
     同时 `TestNPlusOneSubagents...` 仍绿——**尺是分离的**：短流全量保留那一支不该被截断突变连带。
   - **M2 把合并塞回硬顶那一格**（丢弃前先把 oldest 的文本并进下一行）⇒ 两枚红，读数是谎言的原样：
     `row "k46" = "x46x45x44x43...x1x0", want its own delta "x46"`、
     `dropped stream "subagent:t1" was folded into live row "subagent:t4": "only-t4only-t3only-t2only-t1"`。
     ⇒ "N＋1 枚里没有任何一条被并进别人的" 这枚尺**确实能判红**，不是只测"结构体有这个字段"。
6. **动了一枚既有钉，具名**：`pump_test.go:165` 的 `TestTheStreamLogMergesInsteadOfDropping` 更名为
   `TestTheStreamLogTruncatesInsteadOfMerging` 并重写断言。**它是本票要重裁的那一事本身**，
   旧文"two rows, both texts kept, folded"逐条断言的正是"两条流并成一行"，与新语义逻辑互斥、不可能并存。
   不是为变绿放宽：**行数断言从"≤2"收紧成"恰好 3 且逐行等于自己的 delta"**，
   另加"Done 不许从邻行 OR 过来"（旧折叠的副作用）、`Truncated()`、`ElidedRunes()==0`、硬顶、
   以及被删的"every delta still present in the union"一条——那一条的**存在理由就是合并**，随合并一起去。
   与合并无关的旧断言（delta 累积进自己那行、`Close` 收尾、`d` 行 `alphabeta`）原样保留。
7. **两枚"快照四枚字节键"钉：一字未动，且有证据**：
   `git diff -U0 -- internal/panel/pump_test.go` 的 hunk 只落在 `:22`（同名用例的目录注释行，2 行换 2 行）
   与 `:165-205`（被重写的用例本体）；`pump.go:111-124` 那一枚仍在 `pump_test.go:110-125`（断言行 **`:124`**，未漂），
   另一枚从 `:270-276` 漂到 **`:288-292`**（漂 +16，纯粹是上面用例变长），
   `git diff` 里 `four PanelSnapshot declares` / `composer,generatedAt,pending,results` 出现次数 = **0**（＝两枚钉的文本没进过 diff）。
   ⚠ 连带一处过期引用：`pump.go:37` 那段（票 145 的历史记录）写的是 `pump_test.go (:111-124, :270-276)`，
   第二处现在应是 `:288-292`。**本程不改那一段**（它是"不许动的钉"的说明书，改它比留它过期更可疑），只在此具名报出。

---

## 4. 门禁逐包读数

| 门禁 | 读数 |
|---|---|
| `go build ./...` | 改前（18:4x）**rc=0**；改后（18:5x）**rc=1**，唯一错误块逐字 `# github.com/CarlosShao/wisp/internal/tools` + `internal\tools\subagent_197.go:198:31: invalid character U+FF0C '，' in identifier`（另有 `:198:74`／`:198:87`／`:198:112`／`:304:35`／`:305:34`）——**腿 A 写面，非本程**；本程自己那一包 `go vet ./internal/panel/` **rc=0** |
| `go vet ./internal/panel/` | **rc=0**，无输出 |
| `go test -count=1 ./internal/panel/` | FAIL，红 4 枚逐名见 §3#2；`-race` 同一套、无 DATA RACE（面板包 15.2s / 1.3s） |
| `gofumpt -l` 本程 3 枚文件 | **空**（`$GOPATH/bin/gofumpt.exe`；裸 `gofumpt` 不在 PATH，这一枚具名报出） |
| `sh scripts/d22scan.sh` | 第 1 步（种子自测）**FAIL**：`--- FAIL: TestScannerSelfScanOfRealRepoIsGreen`、`--- FAIL: TestRealRepoLedgerIsHonest`；第 2 步真实扫描的**唯一 finding** 逐字 `internal/tools/subagent_197.go:1: [unparseable] ...illegal character U+FF0C '，'` ⇒ **`internal/panel/**` 零 finding**（含注释与 `_test.go` 的 ban #8 扫了 444 枚 Go 文件）。判读＝"逐名"，且这两枚红是腿 A 的半落文件造成的仪器目盲，不是面板侧违规 |
| `internal/observe/nobarego_test.go` | 在场且扫不到本程：它只扫非 `_test.go`（现读 `:39-41`），本程产码里**没有** `go func(`；测试里那三枚 goroutine 带 owner/recover 形状（`subagent_stream_197_test.go:154-168`） |

禁区核对：`internal/tools/**`、`cmd/wisp/**` **一枚文件都没开过**（连读也只读了派单点名要复算的 `internal/panel/**`）；
`Snapshot` 顶层键集未加键；`internal/panel/tokens_fourway_test.go` **未打开过**；`frontend/**`／`design/**` 零写；
`PLAN.md`／`docs/specs/**` 零写；无删除命令、无 `--amend/reset/rebase/stash/checkout ./restore/clean/worktree`；
提交只 `git add -- <显式路径>` ＋ `git commit` 带显式 pathspec，**未 push**。

---

## 5. 不修会怎样 / 修好什么不会坏（具名）

**不修会怎样**：`DefaultStreamKeys=32` 的旧规则下，第 33 枚并发流一到，
载体就把最老的两行**拼成一行**并把两行的 `Done` **OR** 到一起。
今天的生产者（每 `wisp run` 一枚 task 键）到不了 32，所以这枚雷**从未被点过**——
票 197 第 1 层把 8 枚子代理＋根共享一个池装进去之后，它就是可达路径：
owner 点某枚子代理，看到的却是**别人那句工作的续集**，而且包里没有任何一处写着"这里并过"。
界面在说谎，且说的是**出处**的谎（和 D30 防的是同一枚病，只是从显示侧进）。

**修好什么不会坏**（逐条有尺）：
- 快照**四枚顶层键**与 `ResultChunk` 的**三枚键**一字未动 → 两枚字节级钉（§3#7）与
  `TestComposerContractTypesMatchFrontend` 的错因集合都没新增；N+1 用例里还正面断言了
  顶层 4 枚、每行 3 枚（`subagent_stream_197_test.go:116-141`）。
- **短流全额无损**：越界之后仍不削 ≤512 rune 的行（`ElidedRunes()==0` 两处断言）。
- **32 以内行为完全不变**（`Truncated()` 在 3 键/32 键界内为 false，两枚用例正面钉）。
- `Chunks()` 到达序、`Close` 对未知键仍记 done、`Marshal/Publish` 路径——都没动。
- 旧钉里与合并无关的三件事（delta 累积、Close 收尾、d 行全文）原样绿。

---

## 6. 我没测到什么（≥3 枚，具体）

1. **没有真子代理喂这条流**：`internal/panel` 里没有 197-r1 的任何代码，本程只按 §0 的**形状**喂与读。
   ⇒ 腿 A 若把键写成别的（`subagent/`、大写、加了序号后缀），这里照样"每键一行"，**不会有任何一枚尺替它响**；
   跨腿一致性要在集成那一程或对抗验收里核。
   ⚠ **这一条已经现形**（现读，18:5x，腿 A 写面）：`internal/tools/subagent_197.go:91` 自己定义了
   **同名第二枚** `SubagentStreamKey`，并在 `:377` 用它喂 `t.d.Stream(...)`。仓里现在有两枚同名 helper、
   两包各一份，**没有一处共用**。派单 §0 写"键形状写死"，写死的是**形状**不是**实现落点**，
   所以本程不判腿 A 违规、也不去动它（`internal/tools/**` 是禁区）；但它正是这条"没测到什么"要防的形状漂移，
   具名交给裁决者：要么 197-r3 把 helper 收到一处，要么由验收判它重复。
2. **截断在界面上不可见**：没有任何用例断言面板**显示**"已截断/谁被丢了"。行内 marker 是文本，
   Go 侧 `Truncated/ElidedRunes/DroppedKeys` 三枚出口**今天没有任何生产调用者**（现读：新出口在仓里除本程测试外零引用）。
   载体那一层（197-r3）落地之前，一行被截断的文字在页面上和全文**长得一样**。
3. **没测内存字节数**：界是"每行 rune 数 × 行数"，没有一枚尺 marshal 一枚饱和的日志再量包体大小；
   若将来出现"行数在上限内、每行都在头尾地板"的增长，本程的尺判不出红。
4. **没接 §0 的并发上限 8**：`StreamLog` 无从知道宿主认为几枚任务活跃（甲就是这个原因没选），
   所以"子代理数超过 32 枚是否还能发生"这件事在本包里没有被任何断言覆盖。
5. **`Append("")` 仍被接受**（旧行为，刻意不改）：非空文本 + 空键会开一行**没人认领**的记录；
   `SubagentStreamKey("")` 返回 `""` 挡住了子代理那一形，但直接 `Append("")` 的老路径没有尺，
   也没测 `Close` 造成的空行与硬顶丢弃顺序的相互作用。
6. **`-race` 只覆盖本程那枚并发用例的形状**：三枚写者各写各的键，**没有**测"两枚 goroutine 写同一枚键"
   （真扇出里同一 taskID 不应并发写，但这条约定在载体侧无尺）。

---

## 7. 工具调用预算

派单预算 **30 枚**，本程实际 **50 枚**（含写本件、本行更正与提交），**超 20 枚**。花在哪（可核）：
骨架现读 13 枚（派单／锚点＋名册／`pump.go` 两段／`pump_test.go` 两段／`ResultChunk`＋契约尺定位／
基线跑／`panel.ts` 现读／`nobarego`＋`d22scan` 射程现读）；产码与修订 9 枚（含 1 枚 harness 报错重试：
`Edit` 对 `pump.go` 报 `File does not exist`，`ls` 复核文件在、重试即成功——**非权限拒绝**）；
门禁与归因诊断 8 枚（其中 3 枚花在"第 4 枚红到底是谁的"——`TestSliceA…` 与 `d22scan` 两枚自测的错因都指向腿 A 的
`internal/tools/subagent_197.go`，本程必须自己跑到逐字读数才敢归因）；
**实现自身的两发红 2 枚**（§3#4：不幂等折叠、写回顺序，都是跑出来才发现的，不是预判）；
正控突变 2 枚（`-overlay` M1/M2，仓内文件未动）；行号／钉完整性／跨腿 helper 复核 4 枚；
证据件 1 枚 ＋ 其落点更正 4 枚（首稿三处 `文件:行号` 是我按记忆写的、现读后逐条改准，这一笔是**我自己的账**）；
提交 1 枚。
**未为此放宽任何判据或断言**：两枚"快照四枚字节键"钉未动（§3#7），
3 枚在册常红逐名报、不当绿、不修、不 Skip，第 4 枚红逐名归因到腿 A 的写面。

被拒调用枚数（权限层拒绝）：**0**。
