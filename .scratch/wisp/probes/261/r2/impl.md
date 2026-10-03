# 票 261 · 写码腿 `261-r2` 交付件 — 22 枚红归位 ＋ Ⓑ 拒句指路 ＋ Ⓐ firstrun 引导

状态：**已完成**（§0–§8 全部写满；门禁四数齐；22 枚红归零；恒真自查与变异自证双发齐）。

---

## 0. 起手锚

| 项 | 读数 |
|---|---|
| 起手时刻 | `2026-10-03T17:08:27+08:00`（`date -Iseconds` 自取） |
| 起手 HEAD | `cf95c7995ba38a9a8ee100d88c2f719a8ba0b95d`（`git log -1 --format=%H` 自取） |
| 分支 | `dev` |
| 起手写面态 | `git status --porcelain internal/llm cmd/wisp` ＝ **空**（0 行） |
| 第一 commit | 见交件消息（commit 后本节不回填 hash） |
| 票面 | `.scratch/wisp/issues/261-*.md`「编排者翻勾节」读过（17:0x 节，v1 的新钉与 Ⓑ＋Ⓐ 同批裁）；AC 格五格已翻勾的凭据链在票面；⛔ 未碰任何勾框 |
| 证据件 | 本文件＋`probes/261/r2/` 下全部读数（`baseline-before.txt`／`after-cmdwisp.txt`／`proof/`／`mut/`，只建不删） |

**基线复现**（改前存档，`baseline-before.txt`，`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 ./internal/llm/ ./internal/config/ ./cmd/wisp/`）：
`internal/llm` ok 52.110s ＋ `internal/config` ok 2.695s ＋ `cmd/wisp` **FAIL 32.860s，19 枚顶级 FAIL 后 panic 截断**
（panic 逐字 `panic: Fail in goroutine after TestAC1AlwaysBranchDoesNotRevertAHandEditedKey has completed`，
`approval_reply_201_test.go:672` `waitForSentence187`）。与 `255-r5` 两发读数（19 顶级 FAIL／panic 同点）一致。

**fixture 拓扑（本程现量，改前逐文件核过）**：22 枚红**全部汇聚到一处 fixture**＝`newReplyHost`
（`cmd/wisp/approval_reply_201_test.go:119-142`）：`[llm.providers.acme.models.m1]` 条目**缺 `enabled` 键** ⇒
解码 false（`defaults.go` map 分支不执法，票面 §3 的事实）⇒ `261-r1` 的 `resolveEndpoint` 门具名拒绝 ⇒
`wisp run` 装配即死（exit 2）。`newReloadRun223`（223/255-receipt 家族）、`newWidenRun`（226/201-always 家族）、
`approval_seam_201` 两枚、`replyHost.run`（201-reply 五枚）**全部经 `newReplyHost` 起装配**——一处归位，22 枚全回。
派单引的 `255-r5` 名册（22 枚名）逐名核对＝本程基线 panic 前的 19 枚是它的子集（`TestTicket223FailureSentencesAreDistinct`
与 `TestTicket223PermissionDeniedSitsInItsOwnSentence` 在 panic 点之后没跑到），成因**逐字同一**（见 §1）。

---

## 1. 22 枚红逐枚归位表（成因分类 / 改了 fixture 还是断言）

派单给的三形判据，本程逐枚套用后的读数：**a 形＝0 枚**（`cmd/wisp` 没有一枚测试测的是门的拒绝路径——
门的判据全在 `internal/llm` 的两枚文件里，本程已同步）；**b 形＝全部**；**c 形＝1 枚**（次生 panic，见下）。

| # | 红名（`255-r5` §4 名册逐名） | 成因 | 归位动作 |
|---|---|---|---|
| 1 | TestAC1AlwaysBranchDoesNotRevertAHandEditedKey（226） | b（newReplyHost 缺键） | fixture 补 `enabled = true`；**c 形次生**：它的 `driveAlways` goroutine 在 run 死后 `t.Errorf`（panic 源）——fixture 归位后 run 正常走完、goroutine 不越界（见下） |
| 2 | TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card（201-always） | b | 同上（同一 fixture） |
| 3 | TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused（201-always） | b | 同上 |
| 4 | TestReplyListenerAllowsAnL2CardFromTheNativeSide（201-reply） | b | 同上 |
| 5 | TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel（201-reply） | b | 同上 |
| 6 | TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject（201-reply） | b | 同上 |
| 7 | TestUnansweredL2CardTimesOutIntoRejectNeverExecution（201-reply） | b | 同上 |
| 8 | TestL1VetoNeedsAChannelTheHostReallyWired（201-reply，2 子测试） | b | 同上 |
| 9 | TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole（201-seam） | b | 同上 |
| 10 | TestNativeHostSeamRefusesAPanelSourcedAllow（201-seam） | b | 同上 |
| 11 | TestTicket255ReceiptOmitsPanelFromTheImmediateSentence（255-receipt） | b | 同上（newReloadRun223 → newReplyHost） |
| 12 | TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere（255-receipt） | b | 同上 |
| 13 | TestTicket255ReceiptStillNamesTheLiveReadSection（255-receipt) | b | 同上 |
| 14 | TestTicket255ReceiptSentenceAssemblyIsFiltered（255-receipt） | b | 同上 |
| 15 | TestTicket223RunArmsTheReloadTick（223） | b | 同上 |
| 16 | TestTicket223HandEditedFsLooseningCostsAnL2Card（223） | b | 同上 |
| 17 | TestTicket223RefusedLooseningKeepsOldValues（223） | b | 同上 |
| 18 | TestTicket223TighteningRaisesNoCard（223） | b | 同上 |
| 19 | TestTicket223ModeLooseningChangesTheRunningModeAfterAllow（223） | b | 同上 |
| 20 | TestTicket223RestartTierSaysItWillNotApply（223） | b | 同上 |
| 21 | TestTicket223FailureSentencesAreDistinct（223，4 子测试） | b | 同上 |
| 22 | TestTicket223PermissionDeniedSitsInItsOwnSentence（223-perm） | b | 同上 |

**每一枚都是"我改的是 setup 不是 assertion"**：22 枚的判据（逐字业务句）一个字未动；动的只有
`newReplyHost` 种的 config.toml 字符串（`[llm.providers.acme.models.m1]` 补一行 `enabled = true` ＋一段具名注释）。

**c 形判定（派单 1(c) 项，approval_reply_201_test.go:672 那枚 panic）**：**属本程半径，已随 fixture 归位消解**。
证据链：panic 点 `waitForSentence187` 的 `t.Errorf` 从 `driveAlways`（approval_always_201_test.go:196）的
goroutine 里调；`driveAlways` 由 `TestAC1AlwaysBranch…` 用 `go driveAlways(t, w, pw, true)` 起出
（always_write_no_clobber_226_test.go:60）；红发里 run 因门拒 exit 2 ⇒ 主测试 `t.Fatalf` ⇒ 测试完结 ⇒
goroutine 里下一个 `t.Errorf`（`waitForSentence187`，由 `driveAlways` 内部对 201-always 判据的等待调用）＝
`panic: Fail in goroutine after … has completed`。**fixture 归位后 run exit 0，goroutine 在测试存活期内收束**——
本程 cmd/wisp 终态整包零 panic（§5 发 A），此判定由运行读数背书，非推断。

**fixture 归位的额外三处（panic 截断解除后才暴露的同形面，防"新增红"）**：`255-r5` 两发整包都死在 panic，
`run_test.go`／`run_mode101_test.go`／`resident_task_source_live_246_windows_test.go` 三个同形缺键 fixture 的家族
（`newRunFixture` 19 处消费者、`t101boot` 五枚、246-live winlive 档）在 `255-r5` 的读数里**从未跑到**。本程
panic 解除后整包会跑得更远，若不归位它们就是"新增红"——逐处补 `enabled = true`（同 b 形，判据未动）：
`run_test.go`（newRunFixture）、`run_mode101_test.go`（t101boot/writeConfig）、
`resident_task_source_live_246_windows_test.go`（prepareResidentHarness）。

**非必需、属 setup 卫生的三处**（消费者不撞门，补齐只为形状一致，具名登记）：
`providers_test.go`（`cmdProviders` 只读目录不点名）、`panel_config_248_test.go`（settings leg 只写不 resolve）、
`secret_test.go:1092`（`config.LoadFile` 只 validate 引用存在性，chain 检查不查 Enabled）。

**没动的**：`firstrun_198_test.go`／`firstrun_198r2_test.go`（首建默认表**零模型条目**，无 enabled 形状可补）、
`config_reload_perm_223_windows_test.go:81`（writeOver 的内联形状**没有 provider 段**，load 本就以 cause=invalid
失败——这正是那枚测试的目的，门不参与）、`logsink_windows_test.go`／`dataroot_128_test.go`（无 config.toml，
Unconfigured 是测试目的）、226 的 `[app] theme` 手编、`TestTicket223PanelInboundSaysHotReloadIsDisabled`
（inbound 只带 `[risk]`）。

---

## 2. Ⓑ 拒句改后逐字 ＋ 断言同步清单

`internal/llm/resolver.go:128-131` 改后逐字（唯一改动行）：

```go
	if !spec.Enabled {
		return Endpoint{}, observe.New(observe.ClassConfig,
			fmt.Sprintf("llm: model %q of provider %q is disabled (enabled=false); re-enable it or remove the entry（缺 enabled 键的条目视为关闭；在条目里写 enabled = true 即可重新启用）", model, provider))
	}
```

派单 Ⓑ 原文要的指路半句＝「缺 `enabled` 键的条目视为关闭；在条目里写 `enabled = true` 即可重新启用」——逐字落进。
旧行（改前逐字）：`fmt.Sprintf("llm: model %q of provider %q is disabled (enabled=false); re-enable it or remove the entry", model, provider))`。

断言同步清单（改断言＝改描述，具名登记；**无一条是放宽**——全部是"新增对指路半句的正向断言"＋ghost 反向对照）：

| 文件:位置 | 改动 | 方向 |
|---|---|---|
| `enabled_gate_261_r1_test.go` TestTicket261R1SelectionRefusesDisabledByExactWording | 新增循环：拒句必须含 `缺 enabled 键的条目视为关闭` 与 `enabled = true` 两个片段 | 正向（收紧） |
| 同上 | ghost 控制新增：ghost 拒句**不得**含 `缺 enabled 键` | 反向（收紧） |
| `enabled_reach_261_test.go` TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled | disabled 拒句同上两片段正向断言 | 正向（收紧） |
| 同上 | nokey 控制补 `else if`：必须含 `enabled = true` | 正向（收紧） |

原有断言（`disabled` 在场、`unknown model` 不在场、class=config、ghost 双向）**全部保留未动**。

---

## 3. Ⓐ firstrun 改后逐字 ＋ "引导形状过门"现跑读数

`cmd/wisp/firstrun.go` `ensureFirstRunConfig` 的收据第 2 段（模型指路句）尾部追加一句（只加不重写），
改后逐字（新部分）：

```
…wisp providers discover 与 probe 读这份文件去问真实端点，不替你写。注意：models.<id> 条目里要写 enabled = true——缺这枚键的条目视为关闭，点名它会在起动时被拒。
```

改动体（Edit 的 new_string 尾部，含一行理由注释「票 261 Ⓐ（A571）：…只加这一句指路，不重写引导」）。

**预设模板那一格的现量**（派单 Ⓐ 说"预设模板与引导文案要带 enabled"）：firstrun 的**预设**＝
`config.NewDefaults()`（firstrun.go:82），不是模板字符串——默认表里 `[llm]` **零模型条目**（firstrun_198_test.go:215
钉着 `[llm.providers` 在首建产物里出现 0 次），所以**不存在"预设模板缺 enabled"这一格**；形状风险只在引导**教的**
手写条目，Ⓐ 的落点就是收据那一句。这与票面 §3 的事实（手写条目才落 false）与 261-v1 的钉（引导教的形状会被门拒）
一致：v1 钉的是"引导通篇零提 enabled"，Ⓐ 补的就是这句。

**"引导形状过门"现跑读数（恒真自查，硬要求那一发）**：
overlay 注入探针 `probes/261/r2/proof/zz_guided_shape_probe_test.go`（overlay.json 映到
`internal/llm/zz_guided_shape_probe_test.go`；跟踪树零残留），
`PATH=… go test -count=1 -overlay=…/proof/overlay.json -run TestZZ261R2GuidedShape -v ./internal/llm/`：

```
=== RUN   TestZZ261R2GuidedShapePassesTheGate
--- PASS: TestZZ261R2GuidedShapePassesTheGate (0.01s)
```

探针双向：**正向**＝手写条目带 `enabled = true`（引导教的形状）⇒ `LoadFile` 收下＋`ResolveChain` 过门；**负向**＝
同一文件摘掉那一行 ⇒ 拒、且拒句带 `disabled`＋Ⓑ 指路半句两片段。**没有正向过门读数＝Ⓐ 是装饰**——读数在场。

firstrun 收据的既有判据不受影响（本程核过）：`firstrun_198_test.go:98`（`新建默认配置`＋路径，Contains）；
`firstrun_198r2_test.go:383`（mustName 列表 `wisp secret set`／`api_key_ref`／`dpapi:`／`env:`／
`[llm.providers.`／`text_chain`／`roles.chat`——追加句未发明新命令，`wispCommand198r2` 能力检查照旧只挑出
`wisp secret set`）；第二次运行不复述收据的判据（firstrun.go 分支只在文件真缺时走）未动。

---

## 4. 变异自证

**变异件**：`probes/261/r2/mut/resolver-noguidance.go`（HEAD `cf95c799` 的 resolver.go 全文拷贝，唯一差异＝
拒句回到改前逐字——指路半句摘掉），经 `mut/overlay.json`（绝对路径形，同 `261-r1/mut` 定式）注入。

**读数**（`go test -count=1 -overlay=…/mut/overlay.json -run 'TestTicket261R1SelectionRefusesDisabledByExactWording|TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled' -v ./internal/llm/`）：

```
enabled_gate_261_r1_test.go:103: … want the refusal to carry the missing-key guidance "缺 enabled 键的条目视为关闭"
enabled_gate_261_r1_test.go:103: … want the refusal to carry the missing-key guidance "enabled = true"
--- FAIL: TestTicket261R1SelectionRefusesDisabledByExactWording (0.00s)
enabled_reach_261_test.go:201: … "缺 enabled 键的条目视为关闭" / "enabled = true"（两行）
enabled_reach_261_test.go:224: … want the guidance sentence (enabled = true) on the nokey control too
--- FAIL: TestTicket261P1EnumerationSkipsAndSelectionRefusesDisabled (0.00s)
```

**变异 ⇒ 指名断言必红**：4 处红全是本程新加的指路半句断言，其余读数（枚举跳过、两向反转、ghost 对照）
不受扰动——断言对Ⓑ的牙是**定向的**。跟踪树零残留：overlay 全程未触碰 `internal/llm/resolver.go`；
树内 md5 `5fccb0258b414f88de76d92db8138a02`＝带指路半句的入库形；变异后 `go test ./internal/llm/` 全包复绿（36.468s）。

**反形自查（不许留恒真尺）**：若把断言换成"拒句**不含**指路半句"，变异（不含）反而绿、正形（含）红——方向
反转恰好暴露那形是恒假尺；本程的正向断言在变异下发红、在正形下发绿，双向都动过，不是恒真。ghost 反向对照
（ghost 拒句不得带指路半句）保证断言不是"拒句变长就行"的形状。

---

## 5. 门禁读数

| # | 门禁 | 命令 | 读数 |
|---|---|---|---|
| 1 | 三包整包（改前基线） | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 ./internal/llm/ ./internal/config/ ./cmd/wisp/` | `llm` ok 52.110s／`config` ok 2.695s／`cmd/wisp` FAIL（19 枚顶级＋panic 截断）＝`baseline-before.txt` |
| 5 | cmd/wisp 整包（改后）×3 | `PATH=… go test -count=1 ./cmd/wisp/`（发 A 非 verbose／发 B 非 verbose／发 C verbose） | 发 A：19/22 归零＋余 1 枚 C18 窗挤压红＋ball 挂死撞包超时（600s 截断）；发 B：21/22 归零（247.451s）；发 C：**rc=0 全绿**（217.767s，FAIL 0）。三发红句集合里**门拒句 0 命中**（grep `disabled (enabled` 全零） |
| 6 | `internal/llm` 整包（改后） | `go test -count=1 ./internal/llm/` | ok 36.468s（含 Ⓑ 新断言） |
| 7 | `internal/config` 整包（改后） | `go test -count=1 ./internal/config/` | ok 0.951s |
| 8 | `go vet` | `go vet ./internal/llm/ ./cmd/wisp/` | 过（exit 0） |
| 9 | d22scan | `sh scripts/d22scan.sh` | clean（`d22scan.txt`：no D22 ban violations; bans #1-5 cmd/=37…） |
| 10 | gofumpt | `"$(go env GOPATH)/bin/gofumpt" -l` 对本程写过的每枚 .go（11 枚） | 空输出（无格式偏差） |
| 11 | solo ×3 always 家族 | `-count=3 -run 'TestAlwaysBranch…|TestAC1AlwaysBranch…'` | 2 绿 1 红（41.44s，C18 窗挤压形，§6.6） |

**整包红名集合改前／改后**：
- **改前**（`255-r5` §4 名册 22 枚顶级，本程基线 panic 前重现 19 枚）：TestAC1AlwaysBranchDoesNotRevertAHandEditedKey／
  TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card／TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused／
  TestReplyListenerAllowsAnL2CardFromTheNativeSide／TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel／
  TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject／TestUnansweredL2CardTimesOutIntoRejectNeverExecution／
  TestL1VetoNeedsAChannelTheHostReallyWired（＋2 子测试）／TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole／
  TestNativeHostSeamRefusesAPanelSourcedAllow／TestTicket255ReceiptOmitsPanelFromTheImmediateSentence／
  TestTicket255ReceiptOmitsASectionWithNoReaderAnywhere／TestTicket255ReceiptStillNamesTheLiveReadSection／
  TestTicket255ReceiptSentenceAssemblyIsFiltered／TestTicket223RunArmsTheReloadTick／
  TestTicket223HandEditedFsLooseningCostsAnL2Card／TestTicket223RefusedLooseningKeepsOldValues／
  TestTicket223TighteningRaisesNoCard／TestTicket223ModeLooseningChangesTheRunningModeAfterAllow／
  TestTicket223RestartTierSaysItWillNotApply／TestTicket223FailureSentencesAreDistinct（＋4 子测试）／
  TestTicket223PermissionDeniedSitsInItsOwnSentence。
- **改后（三发整包实测分布，逐发如实）**：
  - 发 A（`after-cmdwisp.txt`，非 verbose）：**19/22 归零**（红名集合里其余 21 枚顶级名全部消失）；
    余 1 枚 `TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 41.22s 红（红句 `approval_always_201_test.go:159`
    「审批超时（40 秒未确认），C18 一律判拒绝」＝C18 契约的产品行为句，**非门拒句**）；随后包在
    `TestAC246ChannelNeedsBothWindowAndExecutor` 的 ball STA `waitStarted` 挂 ~7 分钟后撞 10m 包超时（见 §6.7）。
  - 发 B（`after-cmdwisp-2.txt`）：**21/22 归零**，同一枚红 41.29s；包跑完（247.451s），ball 挂死未复现。
  - 发 C（`after-cmdwisp-3-verbose.txt`，verbose）：**rc=0 全绿，FAIL 0，217.767s**——22 枚全部归零、零新增红。
  - **solo ×3**（三枚 always 家族，`-count=3`）：2 绿 1 红（第 2 发 41.44s，同一枚同句）——该枚是**独立于整包的
    既有间歇形**，见 §6.6。
- **终点判据的如实读法**：门造出的 22 红（fixture 缺 `enabled` 键 ⇒ exit 2）在 fixture 归位后**已灭**——三发里
  出现过的红句**没有一句是门拒句**（grep `disabled (enabled` 在三发输出中 0 命中），也没有一枚 22 名之外的
  新增红。余下那枚是**前于本门就存在的间歇形**（§6.6 的先例链），三发 2 绿 1 红的分布如实登记，不写成"一次全绿"。
- `TestPanelHostRealWindowHopAndLifecycle`（255-r5 换的 pid 同一尺，⛔ 不许动）：三发整包**均未红**。

## 6. 判不动（具名）

1. **票面 AC 勾框**：归编排者（⛔ 派单）。
2. **Ⓒ `defaults.go` map 分支默认执法**：派单禁区（另格），本程未动一字；"缺键视为关闭"的指路半句因此
   仍是**行为真话**（缺键落 false 的语义在 Ⓒ 落地前不变）。
3. **`internal/panel` 2 枚 colour/token 已知常红**：非本程包；票 261 AC#4 已按已知读。
4. **255-r5 §5.4 的 winlive 几何读数**（真窗尺寸）：本程不碰，仍留原登记。
5. **`checkChainElement` 是否前移查 Enabled**（票 261 翻勾节"残留记账"第三条）：归编排者后裁，本程未动
   `internal/config/**`。
6. **`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` 的 C18 窗挤压间歇红**（三发整包 2 红 1 绿＋
   solo×3 2 绿 1 红，红句恒为 `:159`「审批超时（40 秒未确认）」41.2-41.4s）：**前于 261-r1 的门的既有形**——
   先例＝`probes/255/r3/final-verbose-2.txt` 13:30:49 发（**pre-gate 时代**，261-r1 门 16:19 才落库）：
   同一枚、同句、41.14s，且是那发唯一 cmd 侧 FAIL；该测试用 `newWidenRun` 把 C18 窗（40s）让给
   `driveAlways` 的三段有界等待（首卡句 30s 上限＋15s widen 卡上限＋25s WIDEN 审计上限），整包重载下
   等待贴上限即越过窗。solo 1.15s 绿＝非确定形；本程三发里它的红与门无涉（红句不是拒句、fixture 已带
   `enabled = true`）。**定性归编排者**：是"把 C18 窗与等待预算的比例"另裁，还是按既有间歇形挂账——
   本程 ⛔ 不动它的判据与等待常数（l2Wait=40s 是派给测试的输入，改它=改判据输入，超本程射程）。
7. **发 A 的 ball STA 挂死**（`TestAC246ChannelNeedsBothWindowAndExecutor` → `startResidentBall` →
   `waitStarted` 挂 ~7 分钟撞 10m 包超时）：发 B/C 未复现（包 247s／218s 跑完）、solo 0.29s 绿。栈顶＝
   `internal/ball/d2d_windows.go:108` comCall nil deref（observe 已 recover，`goroutine panic recovered`
   读数在 `after-cmdwisp.txt`），随后 STA 线程 `started` 信号未落 ⇒ 永等。**与 llm/config 零交联**
   （本程改动不触 ball/D2D/observe），但"谁该给 `waitStarted` 上界"归 ball 族（票 246 线），具名留编排者。
8. **票 261 翻勾节"残留记账"的"写侧无重开路径"**：v1 已入账，本程无新读数。

## 7. 推翻清单（谁写的哪句、我的读数支不支持）

| # | 原话 | 谁写的 | 我的读数 | 判语 |
|---|---|---|---|---|
| 1 | "22 枚红的 fixture 在哪些文件（`255-r5` 的 impl 与 verbose log 有名册），先 `grep -rn enabled` 判每一枚的成因" | 派单 §1 | `255-r5` 的名册只有**红名**，没有 fixture 文件名（其 verbose log 345 行，死在 panic 点）；fixture 定位靠本程现量：`grep -rn "models\."` ＋逐文件读装配链 | **修正（定位法）**：名册够判"谁红"，不够判"红在哪一行 setup"；实际拓扑＝22 枚红 ⇒ 1 处 fixture（newReplyHost），不是 22 处 |
| 2 | （隐含）"22 枚红可能是 a 形（测的就是门的拒绝路径）" | 派单 §1(a) | `cmd/wisp` 零枚测试断言 llm 拒句文案（grep `llm: model`／`disabled (enabled` 全零命中）；门的判据全在 internal/llm | **修正**：a 形在 cmd/wisp 不存在；归位面纯 b 形＋1 枚 c 形次生 |
| 3 | "panic 那枚要判是不是门拒次生；不是 ⇒ 具名留给编排者" | 派单 §1(c) | 是次生（`driveAlways` goroutine 越界，fixture 归位后 run exit 0、panic 消失，整包零 panic 背书） | **支持"是次生"分支**；已消解，未留给编排者 |
| 4 | "Ⓐ firstrun 引导**预设模板**要带 `Enabled: true`" | 派单 Ⓐ / 票面翻勾节 | 预设＝`NewDefaults()` 导出表，零模型条目（firstrun_198_test.go:215 钉着）；不存在带模型的"预设模板" | **修正（落点）**：Ⓐ 的实际落点＝引导收据的指路句（模板格里没有可写的键）；形状过门由探针双向证明 |
| 5 | "`newRunFixture`（run_test.go:95）同形缺键"会不会红 | 派单未点名；本程自查 | 255-r5 两发都没跑到它（panic 截断）；panic 解除后会跑——本程先手补齐，终态**零新增红**兑现 | **补充**：派单名册的"22 枚"是 panic 截断后的可见面，不是全量风险面 |
| 6 | "capability 检查（wispCommand198r2）会把新句里的命令当 invent" | （自查担忧，非派单） | 追加句未出现 `wisp <词>` 形状（只有 `models.<id>`／`enabled = true`），198r2 判据不受影响 | 支持 Ⓐ 文案形态可落 |
| 7 | （隐含）"22 枚归位后 cmd/wisp 应稳定全绿" | 派单 §3 终点判据 | 三发 2 绿 1 红＋solo×3 里 1 红；红枚的先例链在 255-r3 13:30（pre-gate）同句同枚 | **部分修正**：门造的 22 红确实全灭（三发门拒句 0 命中）；但"整包应稳定全绿"的前提被这枚前于门的 C18 窗挤压形顶住（§6.6），不是本程改动引入，如实登记不粉饰 |

## 8. 量不到的格子

1. **cmd/wisp 整包的机器态漂移**：包 4-10 分钟一发，读数是逐发快照；C18 窗挤压形与 ball STA 形都是
   负载敏感的，共享树上有别腿跑时分布会漂——本程三发＋solo×3 的分布如实登记，不外推"稳定绿"。
2. **firstrun 收据追加句的"用户能读懂"那一维**：只能量"句子在场＋形状过门"，读不量"用户真照做"——那是
   票 257 ⓒ 形的真机面（winlive），非本程序程。
3. **Ⓐ 对 `wisp run` 首建→手写→过门的全链真机走查**：探针在 `internal/llm` 层把"引导教的形状"直接喂给
   LoadFile+Resolver（同管线），但 `wisp run` 的 argv→exit code 全链没有本机真窗那一发（首建默认表零模型，
   首建路径本身不撞门；撞门只发生在用户照引导手写之后）。探针的负向半发已覆盖"手写缺键被拒"的运行面。
