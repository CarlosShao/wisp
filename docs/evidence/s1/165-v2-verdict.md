# 165-v2 判语件 — 票 165 `AC#6`（门铃仪器）非实现者终裁

腿＝`165-v2`（裁决腿·对抗验收）｜实现者＝另一条腿 `165-r3`（⛔ 它自己的盖章不算凭据，本件每一格都带我重跑的尺）
票＝`.scratch/wisp/issues/165-plan-first-then-act-a-brake-the-voice-entry-needs.md`｜起手锚＝具名提交 `44f51fc8`
被裁物＝`internal/agent/plan_present_doorbell_165r3_test.go`（221 行／9378 B／1 顶层／0 子测试／0 `t.Skip`）＋夹具 `internal/agent/testdata/golden/plan-present.sse`（28 行／1684 B）
原始逐发＝`.scratch/wisp/probes/165/v2/logs/`（20 枚）｜尺与颜色全文＝`.scratch/wisp/probes/165/v2/00-readings.md`

## 一、判决（一行）

**`AC#6` 成立，但带三枚具名缺口（G1 最重，G2、G3 次之；另记一枚票面自设的射程 G4）。**
本格⛔ 由本腿翻框——翻框只归编排者，凭据＝本件。

票面 `AC#6` 的判定句现读逐字（`sed -n '51p'` 那一枚 `- [ ]`）：⛔ 有一枚常驻判据管**"请求里带没带这枚工具＋`tool_call` 那一跳真被走过"**。两维我各自复跑过牙：

| AC#6 那一维 | 钉在哪一行 | 我重跑的那一发 | 结果 |
|---|---|---|---|
| 请求里带没带这枚工具 | `:219`（round 1 body `mustContain`） | M-e（`Resident→false` ＋ 描述换成与查询无共同字） | rc=1，红句 `:219: round 1 body: missing "\"name\":\"plan.present\"" in "{\"model\":\"mock-small\",…"` |
| `tool_call` 那一跳真被走过 | `:138` 与 `:143` | f02 去牙（`sed -i '12d'`，28→27）／M-c-B（假计数＋摘牙） | rc=1，`:138: … 0 rows for "plan.present" over 0 tool_call rows`；`:143: no tool_call row was booked … (rows read back = 0)` |
| 行上的字段是真读回来的、⛔ 夹具里抄的 | `:157`／`:199` | M-b2（`RiskL2→RiskL1`）／M-f（调用号改 `call_ppX`） | rc=1，`:157: row.risk_level = "L1", want "L2"`；`:199`＋`:220` 同红 |
| 工具体一枚没跑（⛔ 拿 mock 顶真的） | `:207` `CallCount()==0` | M-g（`RiskL2→RiskL0`＝他 M-C 那一形，我自己重跑） | rc=1，五枚红同批：`:157 …"L0", want "L2"`、`:184 row.decision = "allow", want "reject"`、`:191 row.error_class = "tool", want "permission_denied"`、`:199 …Rejected:false`、`:208: executed tool calls = 1, want 0 (an ungated L2 call must never execute)` ⇒ 此格**已复核**，与他报的五枚逐行相符 |

常驻性我自己证（⛔ 抄他）：`go test ./internal/agent/ -count=1 -timeout 30m -v` 默认跑到那枚用例（`logs/f04-fullpkg-v.txt:136`），整包 RUN=110／PASS_TOP=86／FAIL_TOP=0／SKIP=0；件首行＝`package agent`，⛔ build tag、⛔ `t.Skip`。

## 二、三枚具名缺口

**G1（重）这枚工具"是谁"整对可换，用例不响。** M-a2＝把常量 `planPresentName165r3`（`:38`）与夹具里的 `"name":"plan.present"` **一起**改成 `plan.presentX` ⇒ **rc=0 全绿**。机理：一枚常量同时喂镜像条目（`:67`）、读数靶（`:136`/`:143`）与请求体断言（`:219`），四处自洽即无人值守。票面硬约束（导入环：`internal/tools/bridge.go:13` 已 import `agent`）决定本包⛔ 能 import 真名册，而"镜像 ⇄ 真名册"这对今天**在任何一枚用例里都没对接**——我的尺：`cmd/wisp/plan_assembled_165r2_test.go` 与 `plan_wiring_165r2_test.go` 断的是真 registry／装配出口，`internal/agent` 这枚断的是镜像，两半⛔ 相等。⇒ AC#6 原文那句"**这枚**工具"在仪器层只钉到"名册形状"，⛔ 钉到"名册里那一枚真名"。
旁证（同一维的另一发）：M-a（只改用例期望名）红在 `:138`、M-a3（只改夹具线上名）红在 `:138` 且读数 `over 1 tool_call rows` ⇒ 单侧改动有牙、双侧同步无牙。

**G2（中）镜像条目的 `Resident` 是装饰。** M-b（`true→false`）**rc=0 绿**。机理具名：`internal/agent/inject.go:44-49` 把目录切成"驻留⇒必注入／非驻留⇒BM25 候选"，`bm25Select`（`inject.go:172-224`）只要 `score>0` 即入选；本用例目录**只有一枚候选**，而查询"帮我把这个项目收拾一下"与描述"先交一份「我要干什么」的清单并停下等点头"共用 unigram `我`/`一`/`下` ⇒ 非驻留照样被注入 ⇒ `:219` 看不见档位在换。⇒ 按派单口径，这一格写**"缺世界"**（单候选目录造不出"驻留⇄检索"的对立面），⛔ 写"判据成立"。缓解（具名，⛔ 抵掉缺口）：真名册的 `Resident`/`Declared` 由 `cmd/wisp/plan_assembled_165r2_test.go:153-158` 钉着，缺的是镜像这一半。对照：`Declared` **有牙**——M-b2（`L2→L1`）**唯一红**在 `:157`（`loop.go:818-826` 把 L1 与 L2 走同一支：`AdmitTask==nil` 一律拒 ⇒ 换到 L1 只改表上那一列、⛔ 改判决），M-g（`L2→L0`）才换支、五枚同红 ⇒ 牙的形状与产码支一致。

**G3（轻，描述面）`:131` 那行注释把 `THE JUDGMENT` 标在一枚对牙不敏感的门上。** 凭据＝M-c-A（`n := sinkCount165r3(...)`→`n := 1`、牙完好）**绿 rc=0**＝那一发⛔ 携带信息；M-c-B（同一处常量计数＋摘牙）**红 rc=1 在 `:143`**＝假计数顶不掉缺牙。真牙在 `:143`（`rowOfTool165r3` 独立再遍历一遍行）＋`:219`/`:220`/`:199`。`:78-82` 的注释自己承认了这一点，所以本件记**措辞会误导下一位读者**，⛔ 记恒真。

**G4（射程，票面自设，⛔ 记这条腿）** 本格量的是**线路**（有人按⇒名字与参数落进表；不按⇒表空），⛔ 量"生产里今天会不会有人按"——那一格需要 AC#6 自己⛔ 采纳的丙形（真 provider 凭据＋真窗）。落件后分母他自己复跑（夹具面 `plan.present` 0→1、其余 10 枚真注册名仍 0）＝本腿独立核过形：`git grep -F -l 'plan.present' 44f51fc8 -- '*.sse' '*/testdata/*'` 命中**只**这枚新夹具。⇒ "整族零读数"只挪走一格、缺口没关掉，那是普查腿的原话，⛔ 本票 AC#6 的射程。

## 三、编排者那道必答题的独立答复

问：**"`n` 那一枚门自身对牙不敏感，这算不算这枚仪器的缺陷？"**
答：**不算"这枚仪器是恒真"，算"一枚冗余门＋一句错位注释"（G3）。** 分开算的证据链在我自己手上：绿的那发（M-c-A）⛔ 携带信息、⛔ 当凭据；红的那发（M-c-B）红在 `:143` 而⛔ 在 `:138`＝摘掉牙之后仍然响，只是换一枚牙响。⇒ "整枚用例恒真"这一形**不成立**；同时"每一个非零门都单独有牙"这一形也**不成立**，本件据 G3 具名。
另按派单 3(b) 的口径：`Declared` 有牙（M-b2）、`Resident` 无牙（M-b）⇒ 那两枚字段**⛔ 同为装饰**这句要收窄成"一枚装饰、一枚有牙"。

## 四、我自己重跑的那几发 ⇄ 实现者那一发并排

| 发 | 我的颜色（命令＋rc＋红句逐字） | 他报的颜色 | 差在哪（具名） |
|---|---|---|---|
| 正向 | `go test ./internal/agent/ -run '^TestAC\|^TestPlanPresent' -count=1 -timeout 30m` ⇒ rc=0、`--- PASS: TestPlanPresentCallBooksTheToolCallRow165r3 (0.07s)` | `--- PASS … (0.13s)` rc=0 | 仅耗时；同枚、同 rc |
| 去牙 | `sed -i '12d' …plan-present.sse`（28→27、`"name":"plan.present"` 1→0）⇒ rc=1、`:138: doorbell sink reading = 0 rows for "plan.present" over 0 tool_call rows: the replay never reached internal/agent/loop.go:637-640` | 同枚红、同读数 `over 0 tool_call rows` | ⛔ 差 |
| 复算他三发 | M-c-A（常量计数、牙完好）rc＝**0**；M-c-B（同一处常量计数＋摘牙）rc＝1 红在 `:143`；M-g（`L2→L0`）rc＝1 五枚红 `:157`／`:184`／`:191`／`:199`／`:208` | 他的 M-A 绿、M-B 红在 `:143`、M-C 五行同红 | ⛔ 差（逐行相符⇒本腿由此⛔ 抄任何一发） |
| 台面 | `git archive 44f51fc8` **八件套**，`.go`＝747＝同射程跟踪件（diff 0 行） | `git archive c9ba5254` **七件套**，报 744／746 | 锚不同＋清单七枚⇒八枚（编排者 A1006 自纠）；我用同一把尺在 `c9ba5254` 现量 **745**、在 `f7d6dbc98` 现量 **747** ⇒ 他的分母口径各**少 1 枚**，⛔ 影响任何一发颜色 |
| 还原 | 10 发突变逐枚当场还原，末次复量 `.sse 7d26fa889b59fc3a2b39ff18883ee331f8ff1c06ddb45a5a23bc6820af9079be`／`.go 00c32049994c441c8927066bc4d4099c1b2977fc20ee9e602bb98f3670d37b75`＝两串**逐位回基线**（我现量的基线串＝同一串，⛔ 抄他） | 同两串前缀 | ⛔ 差 |
| 工作树自证 | `git status --porcelain -- internal cmd` ＝ **0 行**（本腿在编译面零字节） | — | 本腿⛔ 动 `internal/**`／`cmd/**` 任何跟踪字节 |

⚠ 一处本腿口径自纠（记在这里，⛔ 改已落的票尾那一句）：票尾那条 bullet 写"定向突变九枚（他⛔ 做过的形）"——口径偏松。
准确形＝**十发突变**，其中**他没做过的 7 枚**（M-a／M-a2／M-a3／M-b／M-b2／M-e／M-f），**复算他的 3 枚**（M-c-A＝M-A、M-c-B＝M-B、M-g＝M-C，派单第 3(c) 项点名要分开算的那两发＋本腿补的一发）。
派单只许票尾**一条**bullet ⇒ 本腿⛔ 追加第二条、⛔ 改写已落那一行；更正以此段＋读数件 §3 为准。

## 五、实现者交回的两枚断口：只答"复核属实／不属实"（归属⛔ 本腿裁）

- **断口② 复核＝属实。** 尺：`sed -n '123p' cmd/wisp/plan_assembled_165r2_test.go` 逐字＝`list, err := rt.bridge.Tools(context.Background())`（名册出口）；同文件对 `requestBodies|LastRequest|sent["tools"]|adaptertest` 命中 **0**；`ask165r2`（`:67-69`）拿 `rt.bridge.Execute(...)` 直接喂工具名（`CallID=corr`），⛔ provider 流 ⇒ 那枚行⛔ 由 `loop.go:640` 写。另一枚 `plan_wiring_165r2_test.go` 对 request body 断言命中 **0**（同尺）。
- **断口③ 复核＝一半属实、一半⛔ 属实。** 严格那句"产码里 **0 枚**写入"**⛔ 成立**：`internal/llm/probe.go:91` `ToolChoice: &ToolChoice{Mode: ToolChoiceRequired}` 是一枚**产码**写入（非 `_test`，走 `wisp probe` CLI）。他真正说中的那一半**属实**：装配／agent 支 0 枚写入（尺＝`git grep -n ToolChoice 44f51fc8 -- internal/agent cmd/wisp` ⇒ **0 行**），且 `tools/mockllm/chat.go:172-185` 的合成按位取 `req.Tools[0].Function.Name`、⛔ 按名 ⇒ 从装配面按不到。**建议措辞收窄**为"agent／装配支零写入"，⛔ 裸写"产码 0 枚"。

## 六、派单第 9 行那一问的独立裁定（会决定下一波怎么派料）

我现读 `165r3-dispatch.txt:9` 逐字＝"★现在**零枚用例**把一次 `plan.present` 的 `tool_call` 打到 `internal/agent/loop.go:637-640` 那两行 sink 上"——**已带 `plan.present` 限定词**，与票面 `AC#6` 断口① 原文同形。
⇒ **就 line 9 而言⛔ 需要更正**；r3 的 bullet 把那句引成了无限定形（"那两行 sink 今天零枚用例打到它"），那一形**确实**现量不真。他自己补的读数我复跑**属实**：锚 `4366c1ca1`（＝`4f2470703^`）上 `withJournal(store)` 8 处／5 件（夹具分布 `tool-then-text` 3／`slow-tool` 2／`six-tools` 1／`disconnect-mid-toolcall` 1／`max-tokens-toolcall` 1），golden 面 `"name"` 只有 **echo 18／sleep 7**，`plan.present` 在 `.sse`＋`testdata` 面 **0 件**。
⇒ 下一波派料与 AC#6 原文一致的那一形＝**"那两行走得过（8 处用例），但零枚用例带名册里的真工具名打到它"**；⛔ 裸写"零枚用例打到那两行"（那一句现量假），也⛔ 裸写"表里读不出真名"当成本格已关掉的东西（那是 G1 的另一半）。

## 七、门禁与写面自证（逐枚 rc）

工作树（⛔ 部分导出树）：`sh scripts/d22scan.sh` rc=**0**（ban #8 射程 design/=39、frontend/=116、internal/=577、cmd/=139，`clean - no D22 ban violations`）｜`sh scripts/check-path-length-budget.sh` rc=**0**（VERDICT GREEN，跟踪 14330／超预算 60／在册 60／未在册 0）｜`go vet ./internal/agent/` rc=**0**｜`gofumpt -l`／`-d` 两枚被裁字节 rc=**0**／**0**（零命中／零差异）。
⛔ 计入判据的台面形状（本腿⛔ 踩，具名给下一位）：同一把 d22scan 指到七件套树会 `exit status 2`＋`ban #8 scope design/ examined 0 files … never leave a scope pretending to scan`＝GUARD 射程拒答，长得像红。
实现者写面尺（逐枚 `git show --name-only` 十枚）：路径全集 ⊂ `{internal/agent/**, .scratch/wisp/probes/165/r3/**, docs/evidence/s1/165-doorbell-instrument-r3.md, 票 165}` ⇒ 越界 **0 枚**；`c3f55651` numstat `2/0`（只追加）；票框现量 5 勾／1 未勾（AC#6 仍空，他一枚没翻）；logs 现量 16 枚。
本腿名册与枚数⛔ 写死在这里（自指会过期，故交尺）＝`git log --oneline 44f51fc8..HEAD` 里署名 `165-v2` 的逐枚提交，配 `git show --name-only`；本腿写面只 `docs/evidence/s1/165-v2-verdict.md`、`.scratch/wisp/probes/165/v2/**`、票 165 末尾**一条**只追加 bullet。⛔ 翻框、⛔ 改票名／`Status:`、⛔ `-done`、⛔ push。

## 八、如果只能补一件事

**补 G1 那枚"镜像 ⇄ 真名册"的对接钉**：在装配面（`cmd/wisp/**`，本票写面外，需另派一枚写腿）加一条判据，把 `internal/tools/plan_present.go` 的真工具名与 `internal/agent/testdata/golden/plan-present.sse` 里线上那枚名字**要求相等**，正控＝只改夹具名／只改镜像名各该红、两边同改而真名不动**也必须红**。
理由：G2 有 `cmd/wisp/plan_assembled_165r2_test.go:153-158` 兜着真名册那一半、G3 只是注释；唯 G1 能让这枚仪器被"两头一起改名"静静掏空，而 AC#6 那句写的正是"**这枚**工具"。次选（⛔ 替代 G1）：把 `:131` 那句 `THE JUDGMENT` 改成"冗余非零门；牙在 `:143`"。
