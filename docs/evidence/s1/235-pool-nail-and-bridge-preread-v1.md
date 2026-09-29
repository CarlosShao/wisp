# 票 235 裁决表 —— 235-v1（非实现者验收腿）

- **验收人**：编排者之外的非实现代理（腿号 `235-v1`）。**本表一行实现码都没写**：零产码、零测试码改动，
  所有突变只经 `go test -overlay` 施加，每发之后 `git status --porcelain -- internal cmd` 的读数逐条记在
  `.scratch/wisp/probes/235/v1/readings.md` §5。
- **被验对象**：票
  `.scratch/wisp/issues/235-the-pool-nail-comment-reason-was-falsified-by-ticket-222-and-only-one-of-three-cases-sees-the-real-bridge.md`
  的四格判据，交件腿＝`235-r1`。
- **起手锚点（本腿现跑 `git rev-parse --short HEAD`）**：`98df640a`。
  被验实现腿的五枚 commit＝`795ed767`／`1d2ad737`／`a71f0be9`／`feb88b91`／`9c96f103`；
  改动面（`git diff --name-only 25556879..HEAD -- internal cmd`）＝
  `internal/tools/subagent_197_test.go`＋`internal/tools/subagent_222_test.go`，**零产码**（本腿现跑对上）。
- **行号口径**：票面／台件里的旧行号一律只当线索。下表"票面行号"＝**本腿在锚点 `98df640a` 现读的该 AC 判据所在行**。
- **读数标记**：〔我现跑〕＝本腿自己执行的命令；〔读台件，未复跑〕＝只看 `235-r1` 的现场件。两者不混。

## 0. 逐格判据（本表正文见 §1-§4；骨架先落，逐格裁完一格追加一节）

| 格 | 票面行号（现读） | 判语 | 凭据（本腿读数所在节） |
|---|---|---|---|
| AC#1 改注释、结论与断言零字符变更；M6 抬池仍是有效检测力证明 | `:19` | 待裁 | §1 |
| AC#2 三枚 222 用例的前置读数；M3 下三枚全红且红不来自 30s 护栏 | `:20` | 待裁 | §2 |
| AC#3 两枚自证（① 摘掉前置读数退回 1/3；② 注释改回旧句是否有仪器保护） | `:21` | 待裁 | §3 |
| AC#4 门禁与终态读数（四数只从 `-v` 量／gofumpt／vet／d22scan／逐名双向差集） | `:22` | 待裁 | §4 |

## 1. AC#1（票面 `:19`）——**成立**

**判语**：注释段的活是真的、约束也真守住了；末句"抬池⇒这枚钉子仍红"是一枚**有效的检测力证明**，
"4 枚红"这个报法**不算夸大**（但里面有一枚是仪器自造的阻塞，见下面第 5 条）。

凭据（全部本腿现跑，台件 §0／§2／§4）：

1. **删除列只落在注释行**〔我现跑〕：`git diff -U0 25556879..HEAD -- internal/tools/subagent_197_test.go`
   ＝ `-` 五行全是 `// ` 开头的旧注释、`+` 十三行全是 `// `，函数名行与其下两个 `t.Errorf` 一枚字符都没进 diff
   （`--stat` 同向：13 insertion／5 deletion）。比较方向 `MaxConcurrentSubagents > MaxToolConcurrency` 原样（现读 `:412`）。
2. **假理由真的从盘上消失**〔我现跑〕：`grep -rniE "bridge slot for its child|whole life|bigger pool is a lie" --include=*.go internal cmd`
   ⇒ `subagent_197_test.go` **零命中**（只剩 `internal/agent`／`internal/audio`／`bridge.go` 三处无关件）。
3. **历史句逐字在位**〔我现跑〕：`grep -n` ⇒ 命中 **`:407`**（`which is exactly the 7/8-red shape ticket 197 leg A reported.`），
   同段还把它降级成出处（下一行就写 "Raise MaxConcurrentSubagents and this leg still goes red"）。
4. **新理由的载荷句有行为尺托着，不是故事**〔我现读〕：注释里 "a waiting parent holds none of the four permits"
   ＝`Test222WaitingParentHoldsNoBridgeSlot` 的 `len(h.bridge.sem)` 判据（终态 0.00s PASS，本腿 `-count=3` 那发 18 名全绿）。
   ⚠ 一处**措辞级不准**：注释写 "giveBackWhileWaiting in `subagent_197.go`"，本腿现读＝**调用点**在 `subagent_197.go:365`、
   **定义**在 `bridge.go:694`。指调用点不算假，但严格按盘读应当是 "`subagent_197.go` 里的那次调用"。不影响判据。
5. **M6（抬池⇒红）本腿自己重造、自己读数**〔我现跑，⛔ 工作树零突变字节〕：overlay 只把产码常量
   `MaxConcurrentSubagents = 4` 改成 `8`（一行），终态其余一字未动 ⇒ `=== RUN` 218／全名 218／`--- FAIL` **4**，
   与交件报的红名册**逐枚同名同秒**：三枚 197 钉各 **0.00s**（正文 `:413` "池 8 大于桥的 D38d 天花板 4…名册却说它们在跑"、
   `:452` 同族）＋ `Test222SpawnConclusionArrivesThroughRealBridgeChildren` **3.00s**。
   - **有效的那三枚**：`Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)` 正是注释那句 "Raise MaxConcurrentSubagents and
     **this leg** still goes red" 的正面证明 ⇒ 牙还在，这条支撑成立。
   - **3.00s 那枚的性质**〔我现读正文，不是信注②转述〕：它红在 `subagent_222_test.go:488`（"结论…出现在 0 枚父任务回复里"），
     四枚父任务收到的是 `工具 task.spawn 超时（3000ms），已协作式中止`，而 `3000ms` ＝该台件自己的
     `h222PreFixBudget = 3 * time.Second`（`:81`，喂给 `DefaultTimeout :272`），gate 又要求 8 枚孩子同时进桥、桥只有 4 枚许可
     ⇒ **确证是测试构造的阻塞**，与 `222-v1` 注② 同判。
   - **"报成 4 枚红"是否夸大＝不夸大**，三条理由：① 票面现量 #3 自己就是这个数（"4 枚红＝三枚 197 钉＋一枚 222 用例"）；
     ② 交件同一段里**具名限定**了它不当"池 8 有害"的证据（不是把四枚并列当四条理由用）；
     ③ 代码注释正文里被写进去的只有 "this leg still goes red" 那一枚 0.00s 钉，四枚的数只出现在"顺带指出处"的句子里。
     ⚠ 唯一可挑的一处**读者风险**（不影响本格判语）：注释末三行把 "four legs red" 写进了源码头，
     下一位读者如果不去点 `docs/evidence/s1/222-…v1.md` 的注②，可能把这四枚读成四条"池 8 有害"的证据。


## 2. AC#2（票面 `:20`）——**成立**（含一条**具名偏离**：交付机制与票面点名的补法不是同一个，而本腿判这条偏离是**必要**的）

**判语**：票面这格的**硬判据**（同一枚 M3 下三枚全红、且红不许来自 30s 挂死护栏）本腿独立复现成立，
并且比要求的形状更强一档；票面**建议的那句字面补法**（"在 `await` 之前先读一次 `h.probe.runs`"）没有被照做，
但照做会假红，所以这一处不是脱靶、是纠偏。

1. **改前 1/3 我复现了**〔我现跑〕：本腿自己从 `1d2ad737`（AC#2 之前那版）拉出 222 件、只把
   `ParentTools`／`Tools` 两行改指 `&fake197Dir{}`（旧行号 `:278`／`:288` ⇒ 票面现量 #4 的行号本腿独立撞上），
   `-overlay` 跑终态同包 ⇒ `=== RUN` 218／全名 218／`--- FAIL` **1**：
   只有 `Test222SpawnConclusionArrivesThroughRealBridgeChildren (30.00s)`，正文
   `subagent_222_test.go:374: …只等到 0/4 枚，护栏到点`，另两枚 **0.00s 绿**。⇒ 与票面现量 #5／`222-v1` 注① 同形。
2. **改后 3/3 全红、全 0.00s，而且是因果红不是护栏红**〔我现跑＋现读正文〕：同一枚 M3 叠在终态件上
   ⇒ 三枚同名用例各 **0.00s** 红，三条正文分别落在 **`:444`／`:536`／`:604`**（＝三处前置读数的调用点），
   红句同源：`桥上的到达信号一枚都没有，而第一枚父任务已经拿到了结论（真桥上累计执行 0 次／0 次／6 次…）`。
   ⛔ 那枚 30s 护栏自己的文案（`前置读数两头都没到…护栏到点` 与 `只等到 %d/%d 枚`）在我这一发日志里**一次都没出现**
   ⇒ "红不来自挂死护栏"这一句是按正文验的，不是按秒数猜的。
   第三枚里的 "6 次" 也确实只能是它自己那 6 枚宿主调用（`oversubscribed := h222CeilingLiteral + 2`，现读 `:560`），措辞没装成孩子在跑。
3. **比票面要求的形状更窄一档也照样红**〔我现跑，本腿加测〕：票面 M3 要两行同时改，本腿把两行拆开各改一枚——
   - 只改 `:312 Tools:` 一行 ⇒ 三枚 222 用例**全绿**（`--- FAIL` 0）⇒ **那行今天仍是死行**，
     覆写点本腿现读在 `subagent_197.go:293`（`opt.Tools = newSubagentToolProvider(t.d.ParentTools)`；
     ⚠ 票面写的 `:284`／`newSubagentToolChain` 已经漂了号也漂了名，结论不变）。
   - 只改 `:302 ParentTools:` 一行 ⇒ 三枚**全红、全 0.00s**，红句仍是前置读数那句。
     ⇒ 新加的尺盯的是**载荷行**本身，不需要靠"两行一起改"才红。这一条比 AC#2 的字面判据更硬。
4. **具名偏离：交付的不是票面点名的那把尺**〔我现读＋我现跑〕：AC#2 正文写的是"在 `await` 之前先读一次 `h.probe.runs`"，
   落地的是另一件东西：`probe222.arrived`（缓冲 16 的非阻塞投签，`:134`／`:172-176`）＋
   `h.awaitChildOnBridge222`（`:391-407` 的三头 `select`：到达／父任务返回／护栏），`runs` 只被印进失败正文、不当判据。
   本腿判这条偏离**该走**，两条实证理由：
   ① 直接断言 `runs > 0` 会**假红**——三处调用点都在 `close(h.release)` 之后一刻（`:438→:444`、`:531→:536`、`:603→:604`），
     孩子那时有实率仍卡在**第一枚模型调用**里（`child222.Stream` 的 `release` 门），计数器此刻大概率还是 0；
     而交付的形状在同一位置 **0.00s 返回**（终态 `-count=1` 全绿、`-count=3` 18 名全绿），说明它等的不是计数而是**顺序**。
   ② 要把 `runs` 读成判据就得"读到它有值为止"＝拿墙钟时间差做超时（D22 禁形），而到达先于父返回是**因果**
     （父结论来自子、子的结论来自它的工具调用从桥上回来），不需要钟。
5. **反控那发的 drain 顺序无漏洞**〔我现读＋我现跑〕：`h.drainBridgeArrivals()` 在 `:602`、`close(h.release)` 在 `:603`、
   前置读数在 `:604` ⇒ drain 之后不可能先落孩子签；而 `wg.Wait()`（`:592`）之后 6 枚宿主投签必已缓冲
   （`Execute` 里投签在返回之前，`:172-176`）。两枚 leg1／leg2 里宿主一次都不碰探针（`launchParents :327-345` 只发 `task.spawn`）
   ⇒ 那两枚不需要 drain 是**结构事实**，不是漏做。
   缓冲上限那句也核过：本文件最多 6（宿主）＋4（孩子）＝10 枚投签 < 16 ⇒ 非阻塞投签永远不会变成谁在等谁。
6. **本腿额外攻的两处，都没攻动**〔我现读〕：
   - `awaitChildOnBridge222` 可能取走一枚父任务结论再回投（`:397-398`）会打乱顺序 ⇒ 三枚 leg 下游对 `results`
     的用法全是**按内容计数或只把下标印进消息**（`:468-487`、`:538-543`、`:606-610`），没有一处按下标对应孩子，
     且回投后进 channel 的分母仍是 n ⇒ 不伤判据。
   - 前置读数把 `h.parents` 当"先到先红"的判据，会不会在**合法**形状下假红（孩子一枚工具都不跑的那种 spawn）
     ⇒ 三枚 leg 的孩子 fixture 都是 turn1 必发 `probe222` 工具调用（`child222.Stream :216-243`，`EvToolCallStart` 在 `:231`），
     读数绑的是**本用例自己的 fixture**，不是生产承诺 ⇒ 假红面被 fixture 堵住。


## 3. AC#3 —— 待裁

## 4. AC#4 —— 待裁

## 5. 我攻不动的地方

（逐格裁完后回填，不许为空）

## 6. 没做完／留给编排者

（逐格裁完后回填，不许为空）
