# 273-v1 — 裁决件：票面 AC#0（只裁 AC#0 ＋ a2 §5 六枚未决逐枚处置）

本腿＝非实现者只读验收腿。⛔ 不动票面任何框（AC#0–AC#6 一枚未碰）、⛔ 不改 `docs/PLAN.md`／`docs/specs/**`／`internal/**`／`cmd/**`／`scripts/**`／`.github/**` 一字、⛔ 未跑 `go test`／`go build`／`go vet`（同仓 `272-v1` 正在取整包 Go 读数）、⛔ 未读 `design/**`、未进 `D:/wt/fe`。唯一写面＝`.scratch/wisp/probes/273/v1/**`。

---

## §0 起手锚（现跑，逐字）

```
$ date
Wed Oct  7 11:18:22 CST 2026
$ git rev-parse --short HEAD
998d8255
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
（零行＝空）
```

起手时锚点 `998d8255`，那六棵树工作树干净。别人的未提交件在 `design/**`、`.gitignore`、`.scratch/**`——不是本腿的，本腿不还原、不提交、不评价。

本腿全部读数为**静态现跑**：两枚前置腿（`273-a1`／`273-a2`）的枚数一枚没抄；下面每个数都在锚点 `998d8255` 上由本腿自己那把尺产出，尺命令逐条写在段内。

---

## §1 AC#0 判语

### 1.1 第一把尺：`internal/statemachine/table.go` 转移表到底几枚（两把独立尺并报，不合并）

`internal/statemachine/table.go` ＝ **275 行／11,690 字节**（`wc -l -c`，现量）。

| 尺 | 命令 | 读数 | 现量行/号 |
|---|---|---|---|
| 尺A（行级） | `grep -cE 'D43:[[:space:]]*[0-9]+' internal/statemachine/table.go` | **56** | 带行号的表行 56 行 |
| 尺B（唯一号） | `grep -oE 'D43:[[:space:]]*[0-9]+' ... \| sort -u \| wc -l` | **42** | 唯一号 **1..42 连续、无缺号、无 43**（本腿用 `sort -n`＋逐枚对差不报 GAP 复核过；⚠ 抽数字时必须带前缀串一起看，`grep -oE '[0-9]+'` 会把字面量 `D43` 里的 `43` 也抽出来造出一枚假 43） |
| 尺C（权威文字） | `awk 'NR==3057' docs/PLAN.md` 整行读 | **40 条** | 逐字：`20 态（§2）· 40 条转移。**未列出的转移一律非法**（未定义即停，D22 闸门③）。` |

尺C 的旁证（`grep -n '转移' docs/PLAN.md`，本腿现跑）：全 `PLAN.md` 只有三处写条数，**一律 40**——`:3057`、`:1748`「D43 的 20 态 + 40 条转移表」、`:3151`「（40 条转移，20 态，比例正常）」；`:1362`/`:1302`/`:3600` 谈转移表但不写数。**没有一处写 42。**

尺A≠尺B 的原因（本腿自己数出来的，⛔ 不是抄的）：**13 枚唯一号共用多行**，27 行落在这 13 枚号上，另 29 枚号各一行 ⇒ 29+27=56。共用名册（行号现量）：

```
D43:10 → 75 79        D43:12 → 89 90 91      D43:15 → 101 106
D43:17 → 118 123      D43:18 → 128 132       D43:19 → 136 141
D43:23 → 158 163      D43:24 → 168 172       D43:25 → 176 180
D43:32 → 231 236      D43:37 → 264 267       D43:40 → 273 274
D43:41 → 195 204
```

⇒ 编排者背景里那两枚行号（`:195`／`:204` 共用 `D43: 41`）**复核为真、仍在**。但它在票面 §4 末段被写成"这一条"式的孤例；实际有 13 处同号多行（其中 `#12`、`#40` 各跨 3 行），这是 `Row` 结构注释里那条"guard 分支的两个目的地编码为相邻同号行"的既成形状（`table.go:12-13` 逐字："A few rows list two destinations in their To column (guard-branched); those are encoded as adjacent entries sharing one row number."）。

**AC#0 判语（本腿要求锁的写法）**：**代码 56 行／唯一号 42 枚／权威文字 40 条**。三个数一枚都不许抹平；⛔ 本腿未改 `PLAN.md`／`SPEC-08` 一字（改冻结文字＝人工批准，且票面 :61 已由编排者裁为"三层原文都留档、一字不改"）。

`table.go:29` 那行注释——现量逐字：`D43         int   // row number in SPEC-08 §3 (1..40)`——**今天仍与表体矛盾**（表体有唯一号 41、42）。矛盾范围要划清：文件头 `:6-9` 已经把 41/42 显式解释为 SPEC-08 §3 追加行（"the frozen D43 table carries 40 rows; SPEC-08 §3 ... appends two numbered rows ... for 42 rows total. All 42 are implemented and pinned by the row-presence test."），所以过期的是**字段注释那半句 `(1..40)`**，不是整份文件自相矛盾。这个区分是给 AC#3（要不要人工批准）用的，⛔ AC#3 不在本腿射程，本腿只把尺读数交到这儿。

### 1.2 副作用名册的枚数（票面 :19 那枚"我没量"的靶子）

两把独立尺，**同数 50**，且本腿验证过它们逐字节相同：

```
尺D：grep -oE '"[a-z0-9]+\.[a-z0-9-]+"' internal/statemachine/table.go | tr -d '"' | sort -u | wc -l  ＝ 50
尺E：只从 SideEffects 行抽（grep -o 'SideEffects: \[\]string{.*' ... | grep -oE '"[^"]+"' | sort -u）＝ 50
      diff 尺D清单 尺E清单 ⇒ 无差异（IDENTICAL）
```

尺D 的假命中风险本腿自查过：表里非效果的字符串字面量（guard 名 `kws-loaded`／`not-conversation`／`speech-ge-300ms`／`anything not listed is illegal`／`any`／`first`／`conversation` 等）**全部不含 `.`**，落不进尺D 的形；所以 50 不是被 guard 名灌水的数。

另三枚派生读数（本腿现量，直接替 a2 §5-6 补尺）：

- 投递实例（同一名在表里被列出的次数总和）＝ **62**
- 带 `SideEffects` 的字面量行＝ **39**
- **带效果的唯一行号＝ 33 枚**（`1 4 5 7 8 9 11 13 14 15 16 17 19 20 21 22 23 24 26 27 28 29 30 31 32 33 34 35 36 37 38 41 42`）；**零效果的唯一行号＝ 9 枚**（`2 3 6 10 12 18 25 39 40`）

⇒ AC#0 的靶子该这么钉：**名册 50 枚名／62 次投递实例／33 枚唯一号带效果**。"50 枚全投进 no-op"这句话里的"50"是**名册枚数**，不是投递次数（见 §4-4）。

### 1.3 档位边界具名：AC#0 该锁 ⓑ（不锁 ⓐ）

**本腿选 ⓑ**，但 ⓑ 的措辞要按 §2 复量结果改形（B/C 两档语义各有一处要修）。

我这一档的现量口径（全静态，锚在 `Sink` **类型**的构造点与 `Ev*` **常量**的生产者两把尺上）：

| 档 | 含义 | 我现量的枚数 | 具名 |
|---|---|---|---|
| A | 投出去**且被真收件人接住** | **0** | —— |
| A' | 出货二进制里**被投出去**了，但接住它的是包内默认空罐 | **1** | `model.verify-sha256-signature`（D43 行 #37，`table.go:264-265`，事件 `EvDownloadCompleted`） |
| B | 只有 dev 件 `cmd/balldebug` 会把它投出去（静态闭包推演，零发跑过） | **10** | #4 `session.scope-create`/`speech.load-vad-asr`；#13 `audio.discard-buffer`；#29 `session.zero-load-resume`；#31 `session.dispose-scope`/`speech.unload-asr-tts`/`mem.free-os-memory`/`panel.destroy-or-hide`；#32 `mem.rss-verify-10s`；#33 `settling.cancel-fallback` |
| C | 失败分支进了 `Error`，但唯一能投它的那枚要等 `StateError` 的 10s 超时 | **1** | `error.ack`（D43 行 #38，`table.go:269`） |
| D | 非测试代码里**零生产者** | **38** ＝ 50−1−10−1 | —— |

B 档 10 枚是**我自己推的闭包**，不是引 a2：尺＝`cmd/balldebug/main.go` 非测试面只引用 4 枚事件常量（`EvSummon`/`EvVeto`/`EvMuteKey`/`EvInterrupt`，`grep -oE 'statemachine\.Ev[A-Za-z]+'` 现量），机器起手 `StateSleeping`（`main.go:188`），`dispatch()` 恒传 `nil` Facts ⇒ `machine.go:109-111` 把 `facts` 归零值，`KwsLoaded=false`/`SpeechMS=0`/`HasToolCall=false`，于是 #11（要 `speech-ge-300ms`）与 #10 的 `kws-loaded` 支恒不成立；可达链只有 `Sleeping→(#4)→Listening→(#13 EvVeto / #12 EvTimeoutFirstRound 无效果)→Warm→(#29 EvSummon / #31 EvWarmIdle)→Settling→(#33 EvSummon / #32 EvSettleExpired)`；`EvMuteKey` 只在 `Armed`/`Muted` 有行而本链进不了 `Armed`（要 `EvKwsEnabled`，dev 件不投）；`EvInterrupt` 要 `Speaking`，本链进不了（要 `EvFirstToken`，dev 件不投）。

**为什么不锁 ⓐ（"真投＝1／其余全欠"）**——不是因为 ⓐ 数错，是因为ⓐ 会指挥后来人做错事：

1. ⓐ 的"49 枚欠"读起来像**49 枚待实现的收件人**。落地腿照它写，就要给 38 枚**今天没有任何生产者**的效果名写消费者——那些消费者永远不被调用，只能靠**在产码里新造事件投递**或**用测试自己投自己接**来报绿。这正是编队纪律 1.3 禁的"用 mock 代替真的"和对抗验收要退的形状。
2. ⓐ 抹掉了 B/C 之间的**修法差别**：B 档 10 枚在 `cmd/balldebug` 里是真会响的，一旦给包内加默认收件人，dev 件会立刻开始执行这些副作用（`speech.load-vad-asr`、`mem.free-os-memory`、`panel.destroy-or-hide`……）——那是"甲＝装配根注入"和"乙＝包内默认 Sink"两形的**代价差别**（票面 AC#2 的三形表要靠它才选得动）。ⓐ 给不出这层信息。
3. ⓐ 会让 C 档那枚（`error.ack`）被算进"已投出"或"未投出"的任一侧而失真：它既不是响过，也不是永远不响——它是**被 `Close()` 的结构掐死**，一落地就可能翻成"响"。ⓑ 把它单列，后来人才会去查那句 `defer machine.Close()`。

⚠ 锁 ⓑ 的**代价**也要写进票面：ⓑ 的 B/C 两档是**静态推演**（`winlive` 未批、零发跑过），本腿核了机制但**没有**运行时凭据。所以 AC#0 落笔时必须给每档注上"静态／运行时"字样，否则 ⓑ 比 ⓐ 多出来的信息也是**未证的信息**。

---

## §2 承重句复量（"出货二进制里今天只有 1 枚 Effect 真投给收件人"）

### 2.1 尺：按 `Sink` **类型**的构造点锚

```
非测试代码里的机器构造点（grep -rn 'statemachine\.New' cmd internal tools，剔 *_test.go）＝ 2 枚：
  出货那枚 ＝ cmd/wisp/models.go:303
      machine := statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
  dev 件一枚 ＝ cmd/balldebug/main.go:188
      m = statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
两处 Options 里都**没有 `Sink` 字段** ⇒ 走 machine.go:64-67 的包内默认：opts.Sink = func(Effect) {}
```

正控（同一把尺在能命中的地方必须命中）：非测试面 `Sink:` 字面量 2 枚，**都是假命中**——`cmd/wisp/providers.go:194 storeHealthSink{...}`（`llm.HealthSink`）、`cmd/wisp/run.go:998 consoleSink{...}`（`agent.Sink`）⇒ 编排者预警的三枚同名骗子里前两枚确为异包同名，第三枚 `leg_sink_nail_131`（日志收件人）是钉件命名，不在 `Sink:` 字面量这把尺的射程里，本腿**取不到**它的命中（不写成"没有"）。尺在测试面命中 2 枚真形：`internal/models/bridge_test.go:18`、`internal/statemachine/table_test.go:298`（都是 `Sink: func(e statemachine.Effect)`／`Sink: func(e Effect)`）⇒ **尺有牙**。

⇒ **按收件人锚：出货二进制里"真投给收件人"的 Effect＝0 枚。**

### 2.2 那"1 枚"到底是谁（逐枚复跑）

出货那台机器只被 `internal/models/bridge.go` 的 4 处 `Dispatch` 喂事件（`bridge.go:40/46/59/64`），事件面只有 3 枚常量：

| 事件 | 命中的表行 | 该行的 `SideEffects` |
|---|---|---|
| `EvModelMissing` | #2，`table.go:46` | **无**（该行零效果） |
| `EvDownloadFailed` | #37 第二行，`table.go:267` | **无** |
| `EvDownloadCompleted` | #37 第一行，`table.go:264-265` | **`model.verify-sha256-signature`**（1 枚） |

⇒ 出货二进制今天**被投出去的 Effect 名只有 1 枚＝`model.verify-sha256-signature`**，而它的收件人是空罐。两枚腿的"1"这一枚本腿独立复跑到了同一个对象（不是抄数）。

**两枚腿都没写、但最要命的一层**：这枚名字所指的活儿今天**真的在做**——`internal/models/bridge.go:57` 直接调 `b.mgr.VerifyInstalled(id)`（票 109 的 hand-off 复核），**不经收件人**。所以"只有 1 枚 Effect 真投给收件人"若被读成"有 1 枚副作用已经有实现"，就是**双重失真**：那 1 枚的投递对象是空罐，而实现是旁路直调。副作用通道对这一枚同样是断的。

### 2.3 失败分支的 10s 那一枚

`error.ack` 挂在唯一号 #38（`table.go:269`，`StateError + EvErrorAck → StateWarm`）。非测试代码**没有任何一处投 `EvErrorAck`**（`grep -rn 'EvErrorAck' cmd internal`：唯一非测试命中是 `events.go:60` 的常量声明本体，正控＝测试面 `machine_test.go:87`、`table_test.go:256/405` 有真投）。⇒ 它唯一的响法是超时机器自己投：`timeouts.go:52` `{state: StateError}: {10 * time.Second, EvErrorAck}`，由 `machine.go:182` 的 `time.AfterFunc` 武装、`onTimeout`（`machine.go:191-203`）回投。

出货腿 `cmd/wisp/models.go:303-305` 建机器、`defer machine.Close()`（`:304`）——`Close()`（`machine.go:169-174`）置 `closed=true` 并 `stopTimerLocked()`。

---

## §3 a2 §5 六枚未决：逐枚处置（⛔ 不代填它的判断，只裁"影响哪格／要不要补尺／在不在射程"）

**1｜B 档 10 枚＝静态可达推演，没有任何一发运行时读数。**
处置：**影响 AC#0 的档位措辞，不需补尺**。AC#0 的原文要求是"要现跑"＝静态现量；运行时读数属 AC#1–AC#6（本腿射程外，且 `winlive` 未批＝本腿无权加）。落地要求：AC#0 里 B 档必须写成"**静态闭包可达，零发跑过**"，⛔ 不许写成"balldebug 响过这 10 枚"。本腿的独立闭包也是 10 枚（§1.3），所以这个数不是腿自己一家的推演——但它仍是**同一类证据**（读码），不因两腿同数而升级。

**2｜没把"面板/托盘/键盘/语音"入口逐枚追到 OS 事件源。**
处置：**不影响 AC#0，不必补尺**，理由是本腿换了一把更靠得住的尺：非测试代码里引用 `statemachine.Ev*` 常量的**文件枚数＝2**（`cmd/balldebug/main.go`、`internal/models/bridge.go`）——OS 那一跳要投事件，产码里就得出现那枚常量，追不追托盘回调链都改不了 D 档的 38。⚠ 但本腿要如实交底：我用来兜底的"动态造事件"那把尺（`grep -rn 'statemachine\.Event('`）**全仓含测试都零命中**，因此它**没有正控**、本腿**不把它单独当凭据**；D 档 38 的凭据是"常量生产者＝2 个文件"那把。补尺义务落在**选落点的那一格**（AC#2 的三形表），不在这格。

**3｜§1 的 8 枚机器实例没扫"经接口值/`any` 持有"的形状。**
处置：**这一格本腿已替它补上尺**。现量：非测试面 `statemachine.Machine` 出现处**共 6 行／2 文件**（`cmd/balldebug/main.go:178/591/617/630`、`internal/models/bridge.go:24/33`），全部是具体 `*statemachine.Machine` 的字段/参数/局部变量，**未发现** `any`/接口装箱或经第三方接口转手的持有形状 ⇒ a2 那 8 枚实例里，非测试的构造点只有 2 枚、持有者没有隐藏面。**影响 AC#0 的可达性判语：无改判**；此格可销。

**4｜#41/#42 有没有获人工批准，本腿没去台账查 `A##`。**
处置：**不在 AC#0 射程，已由编排者裁掉，不需腿补尺**。票面 :61 原文（本腿读的是票面整行，⛔ 未截断）：10-07 11:0x 编排者已裁并销账（落账 `A655 §7`／`A656 §4`），实质出处＝`docs/PLAN.md:3254`／`:3256`，越界的只有"编号与计数"那层账面，处置＝三层原文留档、`PLAN.md` 与 `SPEC-08` 一字不改。⇒ 本腿不去台账复对（复对台账不是 AC#0 的尺）；但**票面 :61 那句"代码 42 行"是错数**，见 §4-1。

**5｜`error.ack` 归 C 档只靠"Close() 拆定时器"的读码推理，没有时间测量；若 `handOffModel` 在 return 前做够 10s 这一格会变。**
处置：**结论对、理由要换，且它自报的那把"时间尺"不该补**。本腿复核：真正的理由是**构造性**的，不是竞态——`StateError` 只在 walk 的**最后一跳**进入（`bridge.go:46`/`:59` 投 `EvDownloadFailed` 后立即 `return`），`return` 紧接 `defer machine.Close()`（`models.go:304`），这中间只有 `record()` 一个纯内存调用（`bridge.go:76-82`）、**没有任何等待点**；`modelHandoffTimeout = time.Hour`（`models.go:334`）是 `ctx` 上限、不是机器的保命时间。⇒ 就算有人在 return 前耗够 10s，`Error` 也只在那之前的最后一次投递时才存在，时间测量给不出新事实。AC#0 落笔时 C 档要写成"**构造性不可达（Error 只在最后一跳进入，return 即 Close）**"，⛔ 不许写成"偶发／竞态／待测"。

**6｜50 枚名去重有了，但"每枚名分布在几枚行上"的次数分布没单独输出。**
处置：**本腿已现量补齐，此格销**：50 枚名／62 次投递实例／39 行 `SideEffects` 字面量／33 枚唯一号带效果／9 枚唯一号零效果（名册见 §1.2）。AC#0 应把"名册枚数"和"投递实例数"两个词分开用，⛔ 不许再拿 62 当枚数或拿 50 当投递次数。

---

## §4 我推翻票面／编排者的哪几句（具名报回，含给我这轮的背景）

1. **票面 :61 那句"AC#0 从此必须两把尺并报（代码 42 行／权威文字 40 条）"——"代码 42 行"是错数。** 本腿现量：代码**侧带 `D43` 号的表行＝56 行**，42 是**唯一号枚数**。两把尺并报必须写成"代码 56 行／唯一号 42 枚／权威文字 40 条"。⚠ 这条不是让谁去改 `PLAN.md`，是要求 AC#0 落笔时用对量纲；票面框本腿一枚没碰。
2. **票面 :65 的承重句"出货二进制里今天只有 1 枚 Effect 真投给收件人"——按它自己的锚（`Sink` 类型构造点）读，答案是 0 枚，不是 1 枚。** "1"是"**被投出去**"的枚数，接住它的是 `machine.go:64-67` 的包内默认空罐。票面 :30 的 AC#1 尺（"有没有任何一枚 `statemachine.New` 传了 `Sink`"⇒ 期望 0）与本腿读数一致，⇒ 承重句和 AC#1 只在"**投出去 ≠ 被接住**"这个区分成立时共存。建议编排者把这半句改成"出货二进制里今天只有 1 枚 Effect **被投出去**，而它投进的是包内默认空罐；真收件人 0 枚"。（框与票面文字本腿未动，⛔ 等编排者自己改。）
3. **同一句还有个更深的失真（两枚腿都没写）**：那 1 枚 `model.verify-sha256-signature` 指的活儿**今天真在做**，但是在 `bridge.go:57` 直调 `VerifyInstalled`，**不经收件人**。⇒ 后来人照"有 1 枚已经投出去了"去接线，会做出**双份校验**或误判"这一枚已落地"。
4. **票面 :3/:19 转述普查腿的"表上 50 枚副作用名全投进 no-op"——"全投进"三字与我的复量矛盾**：今天连"投进"都只有 **1** 枚（A' 档），其余 **49 枚根本没被投出去**（B/C/D 三档）。票面已经正确地把这句降级成 AC#0 的靶子，但 AC#0 的写法要把"名册 50"与"投递 1"分开钉住；⛔ "50 枚全投进空罐"是错的。
5. **编排者转述里 C 档的因由"但 `defer machine.Close()` 把 10s 定时器拆了"——拆定时器不是决定性的因**（见 §3-5）；决定性的是 `Error` 只在最后一跳进入、`return` 即 `Close`。另：a2 那条"出货腿失败分支**真进** `Error`"配上"所以 `error.ack` 可达"读起来像失败分支会投东西——现量 `EvDownloadFailed` 命中的 `table.go:267` 那一行**零副作用**，失败分支本身**一枚都不投**。
6. **`docs/PLAN.md:3057` 与 `table.go:29` 两条背景复核为真**（逐字："20 态（§2）· 40 条转移。"；`D43 int // row number in SPEC-08 §3 (1..40)`），`:195`/`:204` 共用 `D43: 41` 也复核为真；`build.ps1` 只 build `./cmd/wisp` 一按指令未重复。本腿不评 AC#1–AC#6。

---

## §5 没做完的格（照实报）

- **AC#1–AC#6：本腿一枚没评**（射程外，指令要求等 `cmd/wisp` 写面空出后的落地腿）。§2 的构造点读数与 AC#1 的尺**重叠但不代它裁**。
- **票面框：0 枚勾选、0 处修改。** 本腿对票面是只读。
- B 档 10 枚＝静态闭包推演，**零发运行时**（`winlive` 未批、`go test`/`build`/`vet` 一枚未跑，避免与 `272-v1` 抢 CPU/红名册）。
- 台账 `A655`/`A656` 的原文本腿**未复对**（§3-4 的理由）；若编排者要"批准记录"的独立腿读数，那是另一格。
- a2 `logs/` 目录里的原始输出本腿**未读**（只读它的 §5 正文）；本腿每个数都是自己重跑的。
- `table_test.go` 里那枚"row-presence test"（声称把 42 枚全钉住）本腿**未打开验证**（属 AC#1/AC#4 的射程，且要跑 Go 才有意义）。
- 本腿**没有**为 `table.go:29` 的 `(1..40)` 与表体 41/42 的矛盾提修法——那是 AC#3／人工批准射程。

---

## §6 终态自证

（终态自查与起手同一把尺，逐字对照见下。本腿唯一写面＝`.scratch/wisp/probes/273/v1/**`；`design/**`、`.gitignore`、`.scratch/**` 里别人在飞的件未动、未还原、未提交。）
