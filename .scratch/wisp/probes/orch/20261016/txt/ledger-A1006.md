
## A1006（2026-10-16 12:5x +08，编排者）收 216-r1 并翻框三枚＋收 220-v1 终裁（四格成立·AC#5 不可判）＋立两枚新票 327／328＋★台面配方自纠（七件套漏 `models/`＝假红一枚）＋我自抓三笔

### ① 起手锚（⛔ 引用下面任何一句前先现跑）
- HEAD＝`f7d6dbc9`（`165-r3` 落用例那一笔），`git rev-list --count origin/dev..HEAD`＝**111**（本会话开头我量的还是 95，两小时内＋16 笔＝四枚腿在交件）。
- 台账件：`wc -l docs/reports/pending-and-issues.md` 落笔前＝**18781**（`grep -c '^## A1005'`＝1）。
- 放腿闸门读数：12:07＝**67%／53%**（放行 `326-v1`＋`165-r3`）→ 12:1x＝**82%／77%／100%**（⇒ `270-r2` 按住，具名释放条件已落 `470ae1e5`＋`07aa654e`）。
- ⛔ 杀 `wisp.exe` **PID 21892**（机主 10-15 19:03 自起，`A988`）；本波它握热键与单实例锁＝`cmd/wisp` 那六枚红的成因（成对台件 `6c20357b`）。

### ② 收件＝盘上三把尺（完成通知一律当未验证）
- **`216-r1`（落地腿·票 216，写面 `internal/panel/**`）交件并结案到位**：五笔 `b7676daf`（产码 `composer.go`／`display_scrub_216.go`(新)／`pump.go`／`workspace.go`）→ `92e8cc57`（测试件 `display_scrub_216_test.go`(新)）→ `dcc0f7df`（读数件）→ `dd2506f0`（台件）→ `c9ba5254`（票尾 1 行）。逐枚 `git show --name-only`＝**越界 0 枚**；件 `docs/evidence/s1/216-display-scrub-r1.md`＝**176 行／15799 字节**（⛔ 0 字节）；票 216 那笔 `git show --numstat`＝**1 added／0 deleted** ⇒ 它一枚框都没翻、票名与 `Status:` 没动（我已核对）。
- **`220-v1`（非实现者终裁腿·票 220）三笔已落**：`45c4c03e`（判语件，12:07）→ `f30b7817`（整包第一发＋复跑同点断，12:09）→ `4366c1ca`（票尾 bullet，12:18）。件 `docs/evidence/s1/220-v1-verdict.md`＝**165 行／22909 字节**。判语＝**AC#1 成立／AC#2 成立（落点甲，批准＝`A562 §1`）／AC#3 成立／AC#4 成立／AC#5 不可判**；它自己把票框尺交回＝**5 未勾／0 已勾、一枚没翻**（我已复核）。★它交回两枚"归编排者"：① L2 卡片在生产里点不亮行（已立成本票 **328**）；② **台面七件套清单还缺 `models/`**（已改正，见 ⑥①）。
- **另两枚仍在飞**（只落了部分件，⛔ 当已交）：`326-v1`＝`c2e1c07a` 落了那枚常驻判据件 `internal/memory/dao_session_326v1_test.go`；`165-r3`＝`4f247070`＋`f7d6dbc9` 落了夹具与用例两枚。**两枚都还没写票尾／没跑完门禁** ⇒ 收件等通知之外我自己再看盘。

### ③ 我自己重跑的（⛔ 拿腿交的日志当我的凭据）
- 尺与件：台件 `.scratch/wisp/probes/orch/20261016/txt/replay-216r1-ac2.sh`，日志 `logs/orch-216-ac2-tree-before.txt`／`-tree-after.txt`／`orch-216-ac3-teeth-off.txt`／`orch-216-ac3-restored-green.txt`／`orch-216-panel-head.txt`。
- **AC#2 成对（同台面两枚仓外树，我自己建的树）**：改前锚 `ca0e30fc` 树 **红 rc=1**，红句逐字＝`probe_216_pair_test.go:84: … live control character U+001B at offset 5: "photo\x1b已允许.png"`，同发 `:79` 的读数行带 `raw-ESC-in-go-row=true`；改后树 **绿 rc=0**、同一行读出 `raw-ESC-in-go-row=false`。两枚树各自 `go build ./internal/panel/` rc=0、`.go` 枚数 **574 ⇄ 576**（证非空）。
- **AC#3 去牙**：只在仓外那棵树里把 `scrubForDisplay` 首行插 `return s`（`LANDS: inserted at line 101`）⇒ **恰 3 枚顶层红**（`TestScrubForDisplayEscapesEveryUnseenCharAndTouchesNothingElse` 带 10 枚字符类子用例、`TestTheSnapshotBoundaryWashesBothRosterFieldsAndCopiesTheRows`、`TestAHostileAttachmentNameWalksTheRealBrokerIntoAWashedMenuRow`）；还原后 `git hash-object` 与 `git rev-parse c9ba5254:internal/panel/display_scrub_216.go` **同串** `6ea6d3ca04d7eb07a0ae9f2f8cfc6f4d96b346c0`、还原发 **7 绿／0 红**、树里 `MUTATION` 残留 0 枚。⚠ 枚数口径差别（具名、⛔ 归因）：它写"仍绿两枚"、我这发绿四枚——我的 `-run` 框了那枚测试件里**全部 7 枚顶层用例**，它点名的红是 5 枚里的 3 枚；两把尺射程不同 ⇒ 绿的枚数⛔ 相减、⛔ 谁顶谁。
- **`internal/panel` 整包我自己现跑一发**（共享工作树、`-count=1 -v -timeout 30m`）＝**RUN 246／顶层 151 PASS／5 FAIL／SKIP 0／rc=1（8.510s）**。与它交的两发（225→246 RUN、141/8 → 148/8）分母同形（顶层 156＝151＋5＝148＋8 ⇒ ⛔ 我这两发与它的发相减，台面不同），差的 3 枚红逐枚＝`TestComposerRenderFixtureTellsTheTruth`／`TestGitDimensionHasNoModelCallableTool`／`TestReadGitOnThisRepositoryIsSelfConsistent`＝**正是它具名归因的三枚"导出树缺 `docs/`／`.git`"台面红**，在共享树里绿 ⇒ 台面差别、⛔ 判据差别；我这两发的 5 枚红＝2 枚本机内存仪器（`backpressure_bound_35r8_test.go:309/:336`）＋3 枚 by-design 键钉（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourwayAgree`）＝**既有红，⛔ 算本波引入**。
- **票 328 那三把尺我自己逐条跑过**（不是转述 `220-v1`）：`grep -c 'TaskID' internal/panel/approval.go`＝**0**；名册侧 `internal/panel/subagent_roster_197.go:212-214` 卡片只索引 `CorrelationID` ⇄ `:216-220` L1 窗口索引两枚；认号那把尺 `:277-286` `splitClashSuffix` 要求最后一个 `#` 之后**全是 ASCII 数字**；生产侧 `callCorr` 两分支分别拼 `#call-<index>` 与 `#<providerCallID>` ⇒ **两形都不落在那把尺认的形上**。

### ④ 翻框与结案（裁决者≠实现者）
- **票 216 翻三格**（AC#2／AC#3／AC#4；尺＝机械 sed `.scratch/wisp/probes/orch/20261016/txt/flip-216.sh`，翻前 **4 未勾／1 已勾** → 翻后 **1 未勾／4 已勾**，diff 逐枚只见 `[ ]`→`[x]`、判据文字一字未动）凭据＝上面 ③ 我自己重跑到的那对颜色＋还原哈希。**⛔ 翻 AC#5**（形状缺席那一半我认可，"缺席＝不适用还是＝未闭合"该由非实现者裁 ⇒ 队列 `216-v2`）；**AC#3 我只复现了一半**：正控里"不许因为名字里有'允许'改变判定"那半我只跑到"用例存在且绿"（去牙发里它也绿＝它对这枚牙不敏感），**攻它的那一发没跑** ⇒ 具名写进票尾并归 `216-v2`。**⛔ 改 `-done`**（AC#5 未勾＋残余三枚未落）。
- **票 220 四格待我翻**（凭据＝非实现者腿 `220-v1` 判语，我已复核件真在盘上；`AC#5` 它判**不可判**，理由＝缺"无活实例／前台可控世界"下 `./cmd/wisp -timeout 30m` 的整包终态名册）⇒ 翻框动作放在我把它那两枚"交编排者"事项落账之后，本条只记状态，翻框下一笔 commit。
- **立票两枚**：`327`（显示边界只洗 `Cc`，`Cf` 那批看不见字符今天穿过；六格全未勾）与 `328`（本案见 ③ 末条）；两枚都带 **AC#0 现量普查**当前置、都写了"⛔ 新加 Snapshot 键／新 C17 方法名／新导出 `panel.*` 名 ⇒ 停手上报"。

### ⑤ 我裁的事（⛔ 每件都要盘上可核）
- **`216-r1` 五枚残余逐枚归位**：① 两份洗字符的合一 ⇒ **不并**（要新导出名＝契约面、两侧形状本就不同、会把 `internal/risk` 判定形状拖进射程）→ 队列，要做先落 `A##`；② `Cf` 穿过 ⇒ **另立票 327**（⛔ 塞进 216 已勾格＝改判据）；③ `"attachments"` 一句表两件事＋`attachmentError` 键值恒缺 ⇒ **归口票 313 收尾链**（要新键、撞 313 键钉）；④ `ApprovalCardView.Tool/Args/Reason` 逐字出向 ⇒ **记范围外**；⑤ `cmd/wisp/panel_pump.go` 那半句坍缩 ⇒ **按在 `cmd/wisp` 面空出之后**。
- **`220-v1` 交回的两枚**：① ⇒ **立票 328**（含"今天全绿是因为用例硬写撞名形"这枚缺世界）；② ⇒ **台面配方改正**（见 ⑥①，具名落进我两份在飞腿的派单）。
- **`216-r1` 超预算的两枚原始 `-v` dump（592／547 行）⇒ 裁决＝留着**：我给那条预算管的是 authored 件（它交的 authored 件 176 行、在预算内），原始日志是读数来源；本仓规矩"临时件只建不删"（`issues/README` 规则 8）优先于我的字数偏好。

### ⑥ 我自抓（记我不记腿）
1. **★台面配方漏了一枚：七件套其实要八件**。我今天派给三枚腿的清单一直写 `go.mod go.sum internal cmd tools scripts frontend`，而 `internal/models/manifest.go:317 ResolveManifestPath` 运行时读 `models/manifest.json`（跟踪件 2 枚：`manifest.json`＋`manifest.json.minisig`）⇒ 少了 `models/` 就长出 `committed models/manifest.json not found` 那三枚红（`220-v1` 的 `./internal/...` 一发里 `internal/models` 三枚红正是此形，它具名交回）。⇒ 从现在起清单一律写**八件套 `go.mod go.sum internal cmd tools scripts frontend models`**，并已在两枚在飞腿的派单尾各追加一节更正（⛔ 改它们已读过的段落，只追加）。★教训同形第二遍：`frontend`（10-16 上午）与 `models`（今天）都是**"仪器自己声明的 root 少一枚就造一枚假红"**那一族，我此前只补了前者。
2. **派单缺一条"⛔ 全树 `go test ./...`"**：`216-r1` 跑了 `./...`（进程名册里现量到），那一形会把 `cmd/wisp` 整包重跑一遍、又不带 `-timeout` ⇒ `panic: test timed out after 10m0s` 会被读成"包红了"。腿在飞时我无法插话 ⇒ 预处置已落 `7d19e56e`：它那枚"整包红"我一律先按台面故障复核。两枚新腿的派单已逐字带上该条款。
3. **"取改前锚"我给的是方法没给号**：`216-r1` 第一版按 `HEAD~2` 取改前树，而我自己的两笔提交正夹在中间 ⇒ 差点把已含产码的树读成"改前"。它自己靠两树文件清单作差抓到并改用 `git rev-parse b7676daf^`。⇒ 今后派单里那一节一律写成"**锚＝具名提交号，⛔ `HEAD~n`**"，本条已同步进票 216 票尾与我这两份派单的更正节。
4. **一处引用差点写错**：我原打算在票 216 票尾引"`-done`＝零残余证明（README `L201`）"，落笔前 `sed` 逐行读了 README 第 13–17／189–201 行 ⇒ 机械规矩实为 **README:15"all boxes checked 才改名"**，而 `README:201` 恰恰是在**警告**"别因为看到 `-done` 就以为它零残余"。⇒ 票面那句已按读回的原文写；⛔ 我凭记忆引条款。
5. 存量（⛔ 本波账）：工作树 `go.mod` 仍脏（第 6 天，mtime `2026-10-10 22:28:50`，两枚腿都具名"⛔ 本腿碰、归编排者处置）；跟踪脏件与未跟踪件名册仍大（`.scratch/wisp/probes/**` 各腿自家台面在长）；`internal`／`cmd` 产码面此刻只有 `165-r3` 在飞的两枚未跟踪新件（⛔ 我代提）。

### ⑦ 编队与队列
- **在飞（12:5x）＝3 枚**：`326-v1`（裁决，`internal/memory`＋`internal/session`，已落一枚判据件）∥ `165-r3`（落地，`internal/agent/**`，已落夹具＋用例两笔）∥ `220-v1`（判语已交三笔，等它自己收尾的通知）。**已交并收＝`216-r1`**。写腿只剩 `165-r3` 一枚，`internal/panel/**` 与 `internal/tools/**` 都空着。
- **待放（具名释放条件＝CPU 连发读数 ≤70%，机主 10-13 定的门槛；12:1x 实测 82/77/100 ⇒ 按住）**：`270-r2`（派单 `07aa654e`，写面 `internal/tools/subagent_197.go` 整行注释面）∥ `328-a1`（只读普查）∥ `327-a1`（只读普查）∥ `216-v2`（裁决：AC#5 那一格＋AC#3 未攻的那半）。⛔ 把"写了派单"当"放了腿"——这四枚一次启动都没发生过。
- **排程依据**：写面互不重叠（`internal/tools` ⇄ `internal/panel` ⇄ 只读面）；`cmd/wisp` 整包那一格自 `220-v1` 交完起**空着**，下一枚需要整包终态的腿（票 220 `AC#5`／票 314 那一族）可接。
- ⛔ 推送：本波还有两枚在飞、ahead 已 **111** 枚 ⇒ 收完再推（只 `origin/dev`、⛔ force、推前核 `git rev-list --count HEAD..<远端 tip>`＝0）；推送会自启本机 `slo-full` runner 抢 CPU。

