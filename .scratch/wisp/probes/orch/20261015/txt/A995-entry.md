## A995（2026-10-15 21:5x +08，编排者）★自纠 `A994` 第五节里我自己写坏的那半句（一把窄尺被当成"记忆过期"，实为我的射程错）

`A994` 第六节里我写了两句"按现量改正记忆"，其中关于**枚数**的那半句**错在我这里**，就地具名作废（⛔ 改原文、只追加）：

- 我写的：「名册引的文件＝11 枚、`internal/**`＝**2** 枚（⛔ 记忆写的 3 枚）」，尺＝`grep -oE '"(internal|cmd)/[^"]+\.go"' cmd/wisp/config_readers_255.go | sort -u`。
- 两处错：① **量的是工作树**，而此刻 `324-r1` 正在改这枚文件（盘上现量 `M`）⇒ 取到的是**别人的在飞态**；② 我的正则**只认引号里的裸文件名**，而名册的引用形如 `file.go:123` 且部分写在散文中 ⇒ 射程比记忆里那把窄。
- ★**重放记忆原尺（HEAD blob、逐字可复跑）**：
  `git grep -hoE '(cmd|internal)/[a-z_/0-9]+\.go:[0-9]+' HEAD -- cmd/wisp/config_readers_255.go` ⇒ **27 枚带行号引用、跨 11 枚文件、含 `internal/**` 三枚**
  （`internal/ball/ball_windows.go` ×2、`internal/observe/logging.go` ×1、`internal/audio/gate.go` ×1）。
  ⇒ **记忆的 11／3 两个数今天照样复现**，⛔ 由我"改正"。
- 另一枚顺手读数（**不同射程、⛔ 与上面相减**）：同一 HEAD blob 上把"带不带行号"都算进去的宽尺＝**20 枚不同 `.go` 路径、其中 `internal/**` 7 枚**（多出 `internal/config/manager.go`、`internal/config/tiers.go`、`internal/llm/resolver.go`、`internal/panel/config_handlers.go` 这些**不带行号的裸引用**）。⇒ 引数时必须带"哪把尺"，这就是我记忆里那条"读数必带尺口径"的又一次现撞。

**`A994` 那句实质结论仍然成立、且与这两把尺无关**：`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim` 自票 255 `AC#6`（10-10）起**行号只当读数**——
读码逐字（`cmd/wisp/config_receipt_255_test.go:185-235`）：文档注释自己举的就是我那条旧形（"before AC#6 it did: 313-r2 inserted two lines, the llm row moved :96 -> :98, and this ruler called that 'the evidence drifted' - punishing an edit that never touched the reader"），
代码里漂移那一支现在是 **`t.Logf`（`:228`）**，会 `t.Errorf` 的只有 **token 消失**（`case 0`）与 **token 不唯一**（`default`）。
⇒ 对本波的实际判断不变：**`287-r1` 加测试行打不红 255 那把尺**；真会打红它的是"某枚被引文件里那 27 枚 token 中有谁变得不唯一或不见"。⚠ 这条改的是**我对机制的描述**，⛔ 是"以后引用可以随便加行"——散文中那些**不带行号的裸引用**（宽尺那 9 枚差额）走的可能是别的判据，谁要用这条做决定，先按 `AC#6` 的射程逐枚读一遍。

**教训定式（进记忆）**：拿一把**更窄的新尺**去"推翻"一把**旧尺**的读数之前，先**逐字重放旧尺**（同 `HEAD`、同正则、同射程）。我这次差点把两个都对的数字写成"记忆过期"——这正是本仓"禁令/旧账会被无条件服从"那一族的反向形：**旧记忆错了会害我，旧记忆对我却用窄尺说它错，同样害人。**
