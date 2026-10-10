# 票 293 · `293-v1` · 30 ★跨票动仪器那一格（独立裁，可推翻编排者初判）

对象＝`cmd/wisp/config_readers_255.go` 的两行**行号引用** `resident_ball_windows.go:316` → `:322`
（尺＝`git diff --no-renames 6ef14788..HEAD -- cmd/wisp/config_readers_255.go` ⇒ 恰 **2 枚 `-` ＋ 2 枚 `+`**，
其余一字未动；token `[Hotkeys:  cfg,]` 两处**逐字未变**；`citedRowsFloor = 8` 那枚地板**未动**；
`hotClaimOtherProcess` 前缀**未动** ⇒ ⛔ 断言面没有任何一处被放宽）。

## 1. 本腿自己验的三件事（⛔ 引用编排者的复跑）

**① 行号对得上**：`git show HEAD:cmd/wisp/resident_ball_windows.go | sed -n '316p;322p'` ⇒
`:316` ＝ `	b, err := ball.New(ball.Options{`，`:322` ＝ `		Hotkeys:  cfg,` ⇒ 名册改指的那一行**确实是那个读者**。

**② 这枚改动是不是"不改就红"——本腿做了反证突变体（把复位撤掉）**：
备份 `$TEMP/cr255.bak` → 把两处 `:322` 写回 `:316` → 现跑
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1 -v -run 'TestTicket255RosterEvidenceLines'`
⇒ **FAIL**，红句逐字（`logs/instrument-revert.txt`）：

```
    config_receipt_255_test.go:209: row "hotkey" cites cmd/wisp/resident_ball_windows.go:316 for "Hotkeys:  cfg,", that line now reads "b, err := ball.New(ball.Options{" - the evidence drifted, so re-adjudicate this row (and if the reader really moved, the receipt's claim may have moved with it)
--- FAIL: TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim (0.00s)
```

还原证明＝`sha256sum cmd/wisp/config_readers_255.go` 动手前后**同值 `367d4a43…78cb884`**（`logs/instrument-pre-sha.txt`
＋该 log 末行），`git diff --stat -- cmd/wisp` 空、`git status --porcelain cmd/wisp | wc -l` ＝ 0。
⇒ **那把尺有牙**（它现读被引那一行、比对 token、漂了就红），**这一发是被动复位、⛔ 不是改松保绿**
——**本腿在此确认编排者的初判成立**（不是我复述它：上面两把是本腿跑的）。

**③ 漂移动作是不是不可避免的**：是**部分不可避免**。
`trayMuteState` 那一枚（`resident_audio_windows.go`）它靠"改挂文件尾"把 `:256`／`:257` 那 5 行引用复位了（⛔ 动名册）；
但 `Hotkeys: cfg,` 在 `:322`，而新增的两枚字段**必须在 `residentBall` 的 struct 体内**（实形＝`:174-183` 那一段），
任何 struct 体内插行都必然把它下面那一枚 `ball.New(ball.Options{…})` 整体下移 ⇒ **这一枚躲不掉**，
除非改落点形状（把投影对放进包级变量＝更差的形）。

## 2. 三问逐条答

### ⓐ 本票射程 还是 越界？越界要不要退回？

**判：不是越界，⛔ 退回；但要记两处形状缺陷。**
- 字面形：`AC#3` 那十枚禁列里**没有** `cmd/wisp/**`，本票排程段写的写面＝`cmd/wisp`＋`internal/ball`，
  `config_readers_255.go` **在 `cmd/wisp/` 里** ⇒ 按票面**没有越出写面**。
  尺（本腿 00 件 §3 已跑）＝六枚 commit 的 `git diff-tree -r --name-only --no-renames` 名册里十枚禁列 **0 命中**。
- 实质形：它是**别人家的仪器数据件**（票 255 的收据名册），动它＝动一枚"别的票的验收判据所依赖的事实"，
  这类动作的规矩应当是**先取得授权**，而它是**先动、后在件里具名交给本腿裁**（`10`/`20` 件 §5 申报了）。
  ⇒ 缺陷①：**程序上先动后报**。缓解＝它申报了、动作是最小形（2 行、token 未动）、且本腿已独立证明不改必红。
- 结论：本格 **不退回**。若编排者要立规矩，宜写成一条台账定式（"动他票仪器/名册＝先派单授权，⛔ 事后报"），
  ⛔ 用退回这一票来立（那一来的代价是：要么带红入库、要么把本票射程扩到别人的仪器——两条都不如授权）。

### ⓑ 有没有"既不动行号、也不放宽断言"的第三形？该归谁？

**有，而且它是正解；但归口票 255，⛔ 归本票顺手做。**
- 本腿现量过这枚仪器的形状（`config_receipt_255_test.go:196-211`）：它按**行号**取那一行、再 `Contains(token)`。
  ⇒ 行号是**唯一的脆弱点**：任何插入都打红它，而它红的不是"读者变了"、是"我引用的那格纸挪了位"。
  这正是本仓已经定式化的"行号锚会腐烂"（`A789`：同一枚注释三方三号；`A799` 又记一次）。
- 第三形＝**内容锚**：引用形如 `path#[Hotkeys:  cfg,]`，测试自己 `grep` 那枚 token、
  要求**恰好命中一次**并**报出它当前的行号**（行号降级成"读数"而非"判据"）。
  牙反而更硬：今天这枚只要求"那一行含该 token"（`Contains` ⇒ 行内多余文字放行），
  换成内容锚后可以加"全文件唯一"这一条新断言 ⇒ **⛔ 放宽，是把断言从'纸的位置'挪到'读者的身份'**。
- 为什么⛔ 本票顺手做：改引用格式＝改**收据印出来的那句话**（行号是给用户看的证据串）＋改**测试本体**
  ＋一次要动同族 **22 处行号引用**（本腿尺＝`grep -n 'cmd/wisp/[a-z_]*\.go:[0-9]*' -o cmd/wisp/config_readers_255.go` ⇒ 22 枚命中，
  其中 `resident_audio_windows.go:256/257` 四枚、`run.go` 六枚…）⇒ 那是**票 255 的写面与它自己的再裁定**
  （它自己那句红句就写着 "so re-adjudicate this row"）。
- ⚠ 本腿⛔ 动它（派单第 3 条明令），本件只裁形状。

### ⓒ "自己插 5 行 → 自己造红 → 自己复位 → 自己改别票名册"作为交付质量该怎么记？

**记三条，方向各异，⛔ 一并夸也⛔ 一并打：**
1. **该记功的**：它没有拿"改前也红"把这枚混进去（`20` 件 §5 具名申报），也没有算进环境抖动；
   修法次序是**由轻到重**（① 改挂文件尾复位 5 枚引用 → ② struct 注释 11→6 行 → ③ 才动名册 1 枚引用），
   且第②③步各留了读数。**诚实性这一头本腿复跑核上了**：本腿那枚"撤回复位"的反证突变体（§1②）
   证明它面对的红是真红、它给的红句形状与本报告题一致。
2. **该记缺陷的（实质那条）**：**带红入库了一枚 commit**。尺＝
   `git show 26289b9e:cmd/wisp/resident_ball_windows.go | sed -n '316p'` ⇒ 逐字
   `	// below, so the AC#1 fallback sentence ("配置缺失 ⇒ DefaultHotkeys") is a`（⛔ 不是那个读者，也不是 token 那一行），
   而同枚 commit 的 `config_readers_255.go` 里 `resident_ball_windows.go:316` **仍出现 2 次**
   ⇒ **`26289b9e` 落地时那把仪器是红的**，红着进了跟踪历史，直到 `9dade451` 才复位。
   本仓没有一条明文禁止"中间 commit 带红"，但**片门禁的对象是片**⇒ 后果＝**任何一位在 `26289b9e`..`9dade451` 之间取基线的腿，
   都会多看见一枚与本票行为无关的红**（票 297 那枚只读普查就在这区间里跑过）。
   ⇒ 记法：`AC#2`/`AC#4` 两格**不因此不成立**（终态本腿已复跑逐名作差＝新增红 0），
   但**交付质量项记一枚"分片方式缺陷"**：把"插行"与"复位行号引用"拆成两笔 commit，就等于在历史里挂了一枚红；
   同批落（或先 relocate 再落产码）不会更差。这条宜进台账，⛔ 由本腿翻任何框。
3. **该记口径的**：它的 `10` 件 §4 那枚 `sha256` 还原值今天**不可复算**（本腿同文件现量 `e67b7031…`，见 `10-mutants.md` §5），
   与 `A799` 已抓的"commit 号写错""logs 枚数写错"同类 ⇒ **报告级 sloppiness，第三次**；
   树本身可证干净（`git diff` 空＝更硬的尺），所以这一枚**记号不判失败**，但**它自己的"还原证明"这一件形已经不能用**。
