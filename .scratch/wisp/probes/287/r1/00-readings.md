# 287-r1 — 票 287 `AC#2` 夹具与种坏读数（落地腿自证，非裁决表）

腿：`287-r1`（写腿，独占 `internal/tools/**`）。跑的票：287 的 `AC#2` ∥ `AC#3` ∥ `AC#4`；
`AC#1` 已闭合（只读腿 `287-a1`），本腿未碰那一格、⛔ 未翻任何 `- [ ]` 框、⛔ 未改票名/Status/台账。
起手锚 `400927c4`（`git rev-parse --short HEAD`，2026-10-15 21:16 +08）；交件复量锚 `afcec913`（22:04 +08）。
共享工作树在跑期间漂了 7 发（`41673704`/`c3bf818a`/`aa1427f6`/`76e861a3`/`60481441`/`afcec913` 等），
⇒ 本件每一枚行号都是**本腿现跑**的号，⛔ 抄派单、⛔ 抄票面。

写面（逐枚，越界零）：
`internal/tools/subagent_depth_fixture_287_test.go`（260 行）＋本件＋
`docs/evidence/s1/287-depth-fixture.md`＋`.scratch/wisp/probes/287/r1/**`＋票 287 末尾 1 条 bullet。
⛔ 未碰：`cmd/wisp/**`、`internal/panel/**`、`internal/winsec/**`、三枚冻结件、`frontend/**`、`design/**`（零读零写）。

台面（派单第 4 节那四条，本腿逐条踩过/避开）：

```sh
PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"   # 缺它＝0xc0000135 且没有 --- FAIL
GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\287r1gotmp'            # C 卷＋反斜杠形，先 mkdir -p
mkdir -p /c/Users/swq/AppData/Local/wisp-test-gotmp/287r1gotmp
```
退码一律 `cmd > log 2>&1; rc=$?` 形，⛔ 在 `cmd | tail` 之后取 `PIPESTATUS`。
种坏＝**物理改写＋sha256 还原**，⛔ `-overlay`（对运行时读盘的尺失明）。

---

## 1. 现量锚（`grep -n` 于 `internal/tools/subagent_197.go`，起手 21:16 跑、收工 21:51 复跑同值）

| 问的事 | 我现量到的号 | 那一行逐字 |
|---|---|---|
| 调用者身份那一读 | **271** | `	parentID := TaskID(ctx)` |
| 空 id 的 fail-closed | 272-275 | `	if parentID == "" {` |
| 深度判定入口 | **276** | `	if rec, ok := t.d.Roster.Look(parentID); ok && rec.Kind == TaskKindSubagent {` |
| 拒语两行 | **278-279** | `			"拒绝派生：深度上限是 %d，而派生者 %s 自己就是一枚子代理（父 kind=%s）。"+` / `				"子代理不许再派子代理。", MaxSubagentDepth, parentID, rec.Kind), IsError: true}, nil` |
| 池位判定 | 282 | `	release, ok := t.d.Roster.TryAcquireSubagentSlot(MaxConcurrentSubagents)` |
| 孩子的 `opt.Tools` 被换 | **316** | `	opt.Tools = newSubagentToolProvider(t.d.ParentTools)` |
| 孩子那一行的真发布点 | 332-333 | `	opt.AdmitTask = func(taskID string) func() {` / `		t.d.Roster.PublishSubagent(taskID, parentID, label)` |
| 目录过滤那一半 | 595-597 | `		if info.Name == subagentHiddenTool {` |
| 名字闸那一半 | **604-610** | `	if req.Name == subagentHiddenTool {`（拒语体在 606-607） |
| 深度常量 | 87-91 | `	MaxSubagentDepth = 1` |
| 桥的并发天花板 | `internal/tools/bridge.go:24` | `const MaxToolConcurrency = 4` |

派单给的五个号（271／276／278-279／316／606 一带）与我的现量**逐枚一致** ⇒ 票面旧号（262／267／270／296／573-580）
确已漂 `+9`，成因＝`270-r1` 给同一枚文件加了 20 行。

---

## 2. 夹具的形状（`AC#2` 三件缺一作不出读数）

产码今天到不了 `:276`（`287-a1` 尺③：`:316` 把孩子的 `opt.Tools` 换成 provider，`:604` 的名字闸先出闸）。
⇒ 夹具**故意**把调用者那一程的 provider 摆成**未包层的真桥**（`Tools: x.h.bridge`），
也就是"结构那一半没装配上"的形状——这是 `:276` 唯一可到的形状，也是它必须可测的理由。
深度这一维**一字未动**：`MaxSubagentDepth` 仍是 1、`:316` 的包层仍在、`:604` 的名字闸仍在。

| 前置 | 落在哪一行（用例文件里的号） | 怎么钉（断言，⛔ 注释） |
|---|---|---|
| ① 调用者那一行 `Kind == TaskKindSubagent` | 发布点 `:76`（钩子体 `:69-81`），复核断言 `:158-161` | 经 `agent.Loop` 的 `AdmitTask` 钩子调**产码自己的** `Roster.PublishSubagent`（产码同一发布点＝产码 `:333`），再 `Look(callerID)` 断 Kind；不对就 `t.Fatalf` |
| ② 派发经真 loop ⇒ `corr != taskID` | 尺在 `:124-143`（`corrOf287`），主用例调用 `:176` | 读**持久化的** `tool_call` 行：非空、以 task id 为前缀、**且不等于** task id；红句把 `who / taskID / corr` 三枚同屏打印 |
| ③ 请求 `task.spawn` | 尺在 `:103-121`（`spawnLog287`），调用 `:162`/`:210`/`:253` | `res.ToolCalls == 1` 且回执里必须有 `Name == "task.spawn"`，否则 `t.Fatalf` |

夹具本体＝`runCallerLoop287`（`:57-100`）：`Tools: x.h.bridge`（未包层的真桥）＋
`Provider: &fake197Provider{askSpawnFirstTurn: true}`（C5 那枚冻结接缝，第一程发 `task.spawn`、第二程收工）。
三枚用例＝`:155`（主）／`:203`（反形正控）／`:233`（次序尺）；`rowIDs287` 在 `:146`。

三枚用例（还原后同参读数＝`287r1-00-targeted-green.log`，rc=0）：

| 用例 | 守的那一行 |
|---|---|
| `Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop` | `:276` 那一支（拒语 `:278-279`＋身份 `:271`＋第二代不许存在＋`Rejected==false` 证它不是闸门数） |
| `Test287RootCallerWithTheSameShapeStillDerives` | 同一形状的**反形**：`Kind==root` 时那一支**不许**拒（打掉恒真检测器） |
| `Test287DepthRefusalWinsOverThePoolRefusal` | `:276` 在 `:282` **之前**返回（池满时读的必须是深度那句，⛔ 池满那句） |

---

## 3. 分母与活性（⛔ 0 命中＋死仪器＝假绿）

- 改前（起手锚 `400927c4`，未加用例）：`go test ./internal/tools/ -count=1 -v` **rc=0／PASS=239／FAIL=0／SKIP=0**
  ＝与派单引的 `4d90dfcd` 那发 239 同值。日志＝本目录 `287r1-05-test-pre-full.log`；改前那发 `go build ./...`
  的空日志＝`287r1-04-build-pre.log`。
- 改后（同一台面，`287r1-10-test-post-full.log`）：**rc=0／PASS=242／FAIL=0／SKIP=0**。
- 红名册差集：**空**（改前 0 枚、改后 0 枚；新增的 3 枚全 PASS）。
- 不在我面、本腿**没跑**的既有红（具名，⛔ 当"已验"）：`cmd/wisp` 的
  `TestTicket255RosterStillMatchesTheActualReadSites`（归口票 324，324-r1 独占该编译面）、
  `internal/panel` 5 枚。派单 ⛔ 本腿跑 `go test ./cmd/wisp` ⇒ 记为〔待面空补验〕。

---

## 4. 三发种坏（四件套齐；台件＝`run-mutations.sh`，逐发打印还原判等）

统一尺：`go test ./internal/tools/ -count=1 -run 'Test287' -v`。
基线 sha256（种前＝每发还原后）＝`86de44c4b79b35172c5ca4bd42862f3659e1d9268ab992f975c282eaa6b172e3`
（与 `4d90dfcd` 那发记录的基线同串）。`git status --porcelain -- internal` 在每发还原后只剩
`?? internal/tools/subagent_depth_fixture_287_test.go`（本腿新件，收工已 commit）＝**产码零残留**。

| 发 | 改动（内容锚唯一性由脚本自证，命中 1；落地尺寸＝**hunks=1／改行数=1**） | 种后 sha256 | 红的用例 | 还原 |
|---|---|---|---|---|
| M1 摘掉深度那一支 | `ok && rec.Kind == TaskKindSubagent` 前插 `false &&` | `69087120ea40590e…` | 用例 1＋用例 3 | `RESTORED-EQUAL-OK` |
| M2 身份改读 corr | `parentID := TaskID(ctx)` → `CorrelationID(ctx)` | `7c7185cc95d9b098…` | 用例 1＋2＋3 | `RESTORED-EQUAL-OK` |
| M3 那一支放宽成"有名册行就拒" | `ok && rec.Kind == TaskKindSubagent` → `ok` | `5aaba840f472d410…` | 只有用例 2（反形） | `RESTORED-EQUAL-OK` |

（第一趟跑的红句与这一趟逐字同形，只有 loop 自己铸的 task id 换了——那些 id 每趟都不同，
⛔ 当判据。第一趟的突变体是 CRLF 版，见 §5.6；本节的号与 hash 全部来自第二趟。）

### M1 红句逐字（`287r1-M1-disable-branch.log`，rc=1）

```
subagent_depth_fixture_287_test.go:169: subagent_197.go:278 那句「拒绝派生：深度上限是 1」没有出现在 task.spawn 的回执里＝深度判定那一支没跑到（种坏 M1/M2 就是这个读数形）：text="子代理 df68984f-542b-46f4-85ea-e3f0537adb3b（孙代理）已结束，状态 Settling。\n结论：\n孙代理的结论（深度 1 的树里不该出现这一行）"
subagent_depth_fixture_287_test.go:173: subagent_197.go:279 那句 "子代理不许再派子代理。"（AC#1 尺①数出来全部 *_test.go 零命中）没被说出口：text="子代理 df68984f-542b-46f4-85ea-e3f0537adb3b（孙代理）已结束，状态 Settling。\n结论：\n孙代理的结论（深度 1 的树里不该出现这一行）"
subagent_depth_fixture_287_test.go:179: :278 那句里的「派生者 <task id>」抽不出来＝身份这一维没有读数：text="子代理 df68984f-542b-46f4-85ea-e3f0537adb3b（孙代理）已结束，…"
subagent_depth_fixture_287_test.go:185: 父 kind 那一枚 = [], want "subagent"（判定读的不是调用者自己的子代理行）：text="子代理 df68984f-542b-46f4-85ea-e3f0537adb3b（孙代理）已结束，…"
subagent_depth_fixture_287_test.go:189: 名册里出现第二代（1 枚 [df68984f-542b-46f4-85ea-e3f0537adb3b]）＝深度上限 1 被走到并破掉；调用者=5fdfb5db-4e91-40b7-bc75-f94a6bdcfdaa
subagent_depth_fixture_287_test.go:192: parent197 的孩子行 = 2 枚 [5fdfb5db-4e91-40b7-bc75-f94a6bdcfdaa df68984f-542b-46f4-85ea-e3f0537adb3b], want 1（夹具只登记调用者自己那一枚）
subagent_depth_fixture_287_test.go:195: task.spawn 的 outcome = "success", want "error"：被判定拒掉的那一发必须是错误回执，模型才改得动
subagent_depth_fixture_287_test.go:255: 深度判定没有在池位判定之前拦住（池已满，回执是另一句）：text="拒绝派生：同时在跑的子代理已达上限 4 枚（在跑的：37c75927-9ea3-4ea0-905c-76c74644b441）。这是硬拒，不是排队——等哪一枚结束了再派，或者把任务并成一枚子代理。"
subagent_depth_fixture_287_test.go:258: 调用者自己是子代理，却先读到池满那句＝:276 那一支没跑："拒绝派生：同时在跑的子代理已达上限 4 枚（在跑的：37c75927-9ea3-4ea0-905c-76c74644b441）。这是硬拒，不是排队——等哪一枚结束了再派，或者把任务并成一枚子代理。"
```
⇒ M1 之下**第二代真被造出来**（`:189` 那行带着它的 task id），这就是"这一支今天有人守"的后果形读数，
⛔ 只报句子、⛔ 只报字段非空。用例 1 的 `InFlightSubagents` 那一发（`:197-199`）在 M1 下没红——
孩子 join 之后池位本来就回收了；那一发钉的是"深度拒绝必须先于占位"的次序，它的牙是 §4 之外的一条
（把 `:276` 那一支搬到 `:282` 之后才会响），本腿⛔ 为凑读数去造第四发种坏，具名登记为**未证**。

### M2 红句逐字（`287r1-M2-identity-to-corr.log`，rc=1）

用例 1 红在同一批点位（`:169`／`:173`／`:179`／`:185`／`:195`），另外三发是新读数：

```
subagent_depth_fixture_287_test.go:221: 孩子行 = 0 枚 [], want 1：{CallID:call-197-spawn Name:task.spawn Outcome:success ErrorClass: Text:子代理 a49d0da5-8556-42c4-a5ff-9a65b73538f5（孙代理）已结束，状态 Settling。
subagent_depth_fixture_287_test.go:255: 深度判定没有在池位判定之前拦住（池已满，回执是另一句）：text="拒绝派生：同时在跑的子代理已达上限 4 枚（在跑的：51c37702-5def-49af-8a60-1323d40157c5）。这是硬拒，不是排队——等哪一枚结束了再派，或者把任务并成一枚子代理。"
subagent_depth_fixture_287_test.go:258: 调用者自己是子代理，却先读到池满那句＝:276 那一支没跑："拒绝派生：同时在跑的子代理已达上限 4 枚（在跑的：51c37702-5def-49af-8a60-1323d40157c5）。…"
```
⇒ 身份读错时**连反形那枚正控也红**（孙代理那行挂到了 corr 上，`:221`），这正是票 282/283 那条方法学前置要的：
`corr` 与 `taskID` 两枚值同屏，红句才写得出来。用例 1 的 `:179` 在两种读法下都红（抽不出派生者＝没读数），
`who[1] != res.TaskID` 那一支（`:180-182`）是"派生者读成 corr"的直接牙。

### M3 红句逐字（`287r1-M3-overbroad-any-row.log`，rc=1；只有反形红）

```
subagent_depth_fixture_287_test.go:213: 父 kind=root 的调用者被深度判定拒了（那一支已经宽到认不出 kind）：refusal="拒绝派生：深度上限是 " 出现在 text="拒绝派生：深度上限是 1，而派生者 cd17a648-8d56-4323-a8c5-e7389c9fbf5c 自己就是一枚子代理（父 kind=root）。子代理不许再派子代理。"
subagent_depth_fixture_287_test.go:213: 父 kind=root 的调用者被深度判定拒了（那一支已经宽到认不出 kind）：refusal="子代理不许再派子代理。" 出现在 text="拒绝派生：深度上限是 1，而派生者 cd17a648-8d56-4323-a8c5-e7389c9fbf5c 自己就是一枚子代理（父 kind=root）。子代理不许再派子代理。"
subagent_depth_fixture_287_test.go:217: task.spawn 的 outcome = "error", want "success"：正控必须真的派出来一枚
subagent_depth_fixture_287_test.go:221: 孩子行 = 0 枚 [], want 1：{CallID:call-197-spawn Name:task.spawn Outcome:error ErrorClass:tool Text:拒绝派生：深度上限是 1，而派生者 cd17a648-8d56-4323-a8c5-e7389c9fbf5c 自己就是一枚子代理（父 kind=root）。子代理不许再派子代理。 Spilled:false Artifact: Truncated:false Rejected:false}
```
⇒ 用例 1／3 在 M3 下**照绿**（它们判的是"subagent 调用者该被拒"，M3 恰恰拒得更宽）——
反形那一枚就是为这一形存在的；三枚用例合起来对 M1／M2／M3 都有牙。

---

## 5. 台面自抓（⛔ 藏）

1. **本腿被 d22scan 抓过一次**：首版把 `⊂`（U+2282，落在仪器扫的 `U+2200–U+22FF` 带里、**字符串不豁免**）
   写进了 `t.Errorf` 的文案 ⇒ `sh scripts/d22scan.sh` **rc=1**，两枚自扫用例红：
   `scan_test.go:269: repo HEAD violates: internal/tools/subagent_depth_fixture_287_test.go:249: [emoji] ban #8 ...`。
   处置＝改文案（⛔ 改扫描器、⛔ 加豁免），改后 `rc=0`（`287r1-11-d22scan.log`：`internal/ 568 Go files`）。
2. **M3 的全包档不作读数**：把那一支放宽成"有名册行就拒"之后，`go test ./internal/tools/ -count=1 -v`
   跑到既有 `Test197RowExistsBeforeFirstChildModelCall` 时**挂住**（那枚用例等的是孩子起跑，而 M3 把派生一律拒了）。
   本腿把这趟掐了（`stale-pre-trim-287r1-M3-full.log`，掐到一半：PASS=146／FAIL=5），
   具名登记：`Test236R2TaskCancelRefusesWhenHostGaveNoCallerID`／`Test270R1ThreeRefusalSentencesDoNotSubstitute`／
   `Test270R1WiredBaseOptionsIsNotRefused`／`Test229FinishWatcherIsRosterNamedWhileChildRuns`／
   `Test197SpawnPublishesIdentityRow` 五枚先红。⇒ **这条挂死不是本票的读数、⛔ 算进本票的账**，
   只说明 M3 那形的爆炸半径比定向档大。
3. **文件压行**：首版 302 行、超派单"每件 ≤260 行"⇒ 压到 260 行。压的是注释与换行，
   ⛔ 未删任何断言、⛔ 未放宽任何判据；因此**早期那批读数的 `文件:行` 已过期**，全部改名
   `stale-pre-trim-*` 留在盘上（临时件只建不删），本件第 4 节引的是**终版**行号。
4. `grep -c '^--- SKIP'` 命中 0 时 rc=1，会把 `&&` 链断掉 ⇒ 本腿一律写成 `; echo rc=$?` 或 `|| true`。
5. 突变台件里躺着 `subagent_197.go` 的**三份拷贝**（`mutants/*.go`）与还原基线（`subagent_197.go.pre-mutation.bak`），
   它们带同一枚拒句 ⇒ 后续任何按句子数的尺**必须** `| grep -v "^./.scratch"`（`287-a1` 第 0 节那条前处理，
   本腿照抄并留档）。
6. **突变体第一趟是 CRLF 版**（python 文本模式在 Windows 写文件会把 LF 翻成 CRLF）：后果＝红句照红、
   还原照等值，但那三枚拷贝与基线的 `diff` 变成"整枚文件重写"（hunks=1 却 570 行全换），**落地尺寸的账是假的**。
   处置＝台件改成 `read_bytes`/`write_bytes` 逐字节替换（`run-mutations.sh` 现形），三枚 CRLF 拷贝改名留档在
   `mutants/stale-crlf/`，第二趟重跑：三发各自 **hunks=1／改行数=1**、种后 sha 与第一趟逐字节一致
   （`69087120…`／`7c7185cc…`／`5aaba840…`），红句点位不变。

## 6. 本腿没做／待人拍板

- ⛔ 未裁 `AC#3`（只交代价表，见 `docs/evidence/s1/287-depth-fixture.md` 第 3 节）、⛔ 未翻任何框、⛔ 未改产码一字。
- ⛔ 未跑 `cmd/wisp`／`internal/panel`／`internal/winsec` 任何尺（面不在我）⇒ 那三面的既有红本腿只引派单数字，具名标"未复跑"。
- 〔待面空补验〕：本夹具的"生产链仍不可达"这一句在 `cmd/wisp` 端**没有端到端台件**证它——
  `287-a1` 的尺③是读码＋行号级证据，本腿未加编译面尺（派单 ⛔ 动 `cmd/wisp`）。
- 待人拍板一条（不属本票，顺手量到就报）：`MaxConcurrentSubagents = 4`（`:85`）与票 287 票面 `:11` 的"并发 8"
  仍是分家的两枚数，`287-a1` §5.5 已登记、归口票 211/324 那一族 ⇒ **本腿一字未动**。

## 7. 本目录件名册（谁是谁的读数）

| 件 | 是什么 |
|---|---|
| `00-readings.md` | 本件（锚＋夹具形状＋四件套摘要＋自抓） |
| `run-mutations.sh` | 三发种坏台件：逐字节替换、内容锚唯一性＋LANDS 自证、每发还原并打印 `RESTORED-EQUAL-OK` |
| `run-form-b-census.sh` | `AC#3` 乙形普查台件（B1 只摘名字闸／B2 连目录过滤一起摘，整包档 `-timeout 180s`） |
| `287r1-M1-disable-branch.log`／`M2-identity-to-corr.log`／`M3-overbroad-any-row.log` | 三发种坏的定向档全文（红句逐字在这里） |
| `287r1-B1-name-gate-off.log`／`B2-name-gate-and-listing-off.log` | 乙形两发的**整包**读数（241／1，红名册逐枚） |
| `287r1-00-targeted-green.log` | 还原后三枚用例同参读数（rc=0） |
| `287r1-04-build-pre.log`／`05-test-pre-full.log` | 改前基线（build 空日志＋239／0／0） |
| `287r1-10-test-post-full.log`／`11-d22scan.log`／`12-path-length.log` | 改后三把门（242／0／0、clean、GREEN） |
| `287r1-20-mutations-run-all.txt`／`21-form-b-census-run.txt`／`22-gates-final.txt` | 上面几批的**驱动输出原文**（sha／hunks／rc 逐条打印） |
| `subagent_197.go.pre-mutation.bak` | 种坏前的产码基线（sha `86de44c4…`，与 `4d90dfcd` 那发的记录同串） |
| `mutants/M*.go`／`mutants/B*.go` | 五枚突变体（LF，与基线逐字节只差 1／1／1／1／2 行） |
| `mutants/stale-crlf/*.go` | 第一趟的 CRLF 版突变体（见 §5.6，留档⛔删） |
| `stale-pre-trim-*.log` | 夹具压到 260 行**之前**的读数（行号已过期，见 §5.3，留档⛔删） |
