# 票 126 独立对抗验收 —— `sameTree`/`answerInsideTree` 加卷段（生产判据）

**验收方**：`acceptor-ticket126`（本文件唯一作者）
**日期**：2026-09-22（本会话）
**被验交件**：票 126，实现方 `agent-ticket126`
**基线**：`09b5285`（票 118 收表那发）
**交件 sha**：`a701138`（生产码 `internal/winsec/resolve.go` + 用例 `volume_attribution_126_windows_test.go`）、
`fe93558`（缝守用例改六腿）、`a3ce3a4`（补第七枚 CONTROL）、`bb8393e` / `5e0f63f` / `81b4d5f`（票面三格 + 三格 + 收格，纯文档）
**本文件性质**：1 AC = 1 格，裁决 + 标签；票面一字未改，别人的框未翻。

---

## 0. 会话口径（先声明我怎么做，免得读数无法复算）

- **只读**：本会话在仓库内唯一的写件是 `docs/evidence/s1/126-adversarial-acceptance.md`。
  全程 `git archive <sha> | tar -x -C /tmp/ac126-sq/<名>` 取快照，**没有**在仓库树内建 worktree、
  没有 `checkout`、没有 `--amend`/`reset`/`rebase`/`stash`（A38④）。
- 快照清单：
  | 名 | 内容 | 用途 |
  |---|---|---|
  | `/tmp/ac126-sq/base` | `git archive 09b5285` | 基线整包读数、台账控制组 |
  | `/tmp/ac126-sq/post` | `git archive 81b4d5f` | 交件树整包读数、门禁、变异 |
  | `/tmp/ac126-sq/head` | `git archive bd50c63` | 票 125 叠上之后 126 的判决是否仍成立（本票外溢风险） |
  | `/tmp/ac126-sq/pre2` | `git archive 09b5285` + 交件那枚用例文件 | 复算"改前红" |
- 四数一律从 `-v` 输出量，口径与 `scripts/winsec-tests.sh` 逐字相同：
  `grep -c '^=== RUN'` / `'^--- PASS'` / `'^--- FAIL'` / `'^--- SKIP'`（顶层锚定，子用例缩进不计）。
- 变异先发落地再读名：`grep` 那几字节确认真改掉 → `go build ./internal/winsec/` rc=0 → 才读 `--- FAIL` 名单。
- `GOOS=linux go vet` 只编译不执行；容器真执行另发，加 `MSYS_NO_PATHCONV=1`，挂载后先 `ls -l /src/go.mod` 验非空挂。

**本文件写作状态：骨架。** 下面六格逐格落盘，每格一次。

---

## 1. 首要攻击点（跨格）：它主张"只往严走、不存在把被拒翻成被放行的通路"——我自己找反例

**状态：判定待写**（下面第 1.1 节是我自己穷尽的读数方向，与它的自述无关）

### 1.1 两枚比较在生产里的全部读数点（我自己 grep 全仓，不按它的行号）

`sameTree` 与 `answerInsideTree` 在非测试代码里的调用者，我数出来是**三处、两个方向**。
行号口径我不含糊：票面写的 `:311`/`:325` **在交件树 `81b4d5f` 上逐字成立**（我在那枚快照里自己
`grep -n` 出来就是 311/325）；它们在工作树 HEAD 上是 `:391`/`:405`——**这枚漂移是票 125 后落在同一枚
`resolve.go` 上的 84 行造成的**，不是 126 写错行号。下表用交件树行号：

| 读数点 | 代码 | `true` 被读成 | `false` 被读成 |
|---|---|---|---|
| 缝守第一道包含见证 | `internal/winsec/resolve.go:311` `if answerInsideTree(childAns, parentAns) { return "" }` | **放行**（`""`＝可以装） | 继续往下走 |
| 缝守第二见证 | `internal/winsec/resolve.go:325` `if sameTree(anchorAns, parentAns) \|\| answerInsideTree(anchorAns, parentAns) { return "" }` | **放行** | **拒绝**（走到最后一行的 refusal 文案） |
| 通知归属 | `internal/winsec/winsec_windows.go:111` `return sameTree(n.Path, resolved.String())`（`noticeNamesTree`） | "这枚通知就是关于你那棵树" | "不是"（虚警方向，见 1.2） |

结构性事实：两枚函数本次只加了 `if !sameVolume(...) { return false }` 的**提前返回**，
其余分支一字未动。我把 `09b5285..81b4d5f` 里 `resolve.go` 的**全部 14 行删除**逐行列出来核过：
它们是两枚 `fold := func(s string) string { ... ToLower ... }` 闭包（各 6 行）加两枚调用点——
`foldSegment` 是同一枚函数体搬家，**零行为变化**。
⇒ 对任意输入 `(a,b)`，改后为 `true` 必然改前也为 `true`：**逐点单调收紧**，
不存在任何输入让 `false` 变 `true`。这一枚是我自己读 diff 得到的，不是转述。

### 1.2 我要找的"放宽通路"在哪：归属侧的 `false` 会不会被别处当放行读——**找不到**

`grep -rn "noticeNamesTree\|noticesAboutTree"` 扫全仓（含 `cmd/`、`tools/`），非测试命中只有
`winsec_windows.go` 里它自己的定义加三处注释，和 `resolve.go` 一枚注释。
⇒ **这两枚归属助手今天零枚非测试消费者**：`sameTree` 在归属侧的唯一生产读者是
`noticeNamesTree`，而 `noticeNamesTree` 只被 `_test.go` 调用。
所以归属侧连"把 `false` 读成放行"的通路都不存在，收紧它在生产里一个字节都放不宽。

我也把"第三种可能"（某处 `if !sameTree(...) { return "" }` 这种把取反当放行读的写法）数了一遍：
全仓对这两枚函数的取反调用数为 **0**。

顺带钉住我自己的口径，免得这枚读数被反过来用：票 126 的 `Type` 行写"生产缺陷（通知归属判据）"。
按消费者算，归属这一侧今天确实是**只有测试在读的判据**；但它是包内唯一的归属判据
（`winsec_windows.go:116` 明写 every consumer of narrowNotice 都走它），而**缝守那一侧是硬生产**
⇒ 这枚 `Type` 站得住，我不据这一点降格本票。

### 1.3 反向的代价（过严）钉住了吗——钉住了，且我另打三发验它的钉子是双向的

它主张"过严会误伤诚实 resolver、整条 C26 掉回内置底线（run 35595651898 那一族）"，
交件用七枚腿里四枚 `wantRefused:false` 的 CONTROL 加票 112 那枚用例钉。我复算：

- 四枚 CONTROL 在**改前改后都绿**（`refusal=""` 逐腿亲眼读到，一次 `CONTROL RED` 都没触发）；
- `TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent`（票 112 的"诚实 resolver 必须过缝"）
  基线 `PASS` → 交件 `PASS`；
- **我自己补的两发放宽变异**证明这套钉子双向有牙：把 `sameTree` 掰成恒 `true`
  （＝能造出的最大放宽，比 `09b5285` 还宽）红 **9 枚**，其中 **5 枚是既有用例**
  （`TestAC2InheritedNoticeHasANoiseBound`、`TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem`、
  `TestSealReportsThePrincipalsItCleared`、`TestNoticeAttributionKeepsTwoTreesApart`、
  `TestTreeOwnershipProbeAcceptsAn83ShortSpellingOfItsOwnParent`）；
  把 `answerInsideTree` 掰成恒 `true` 红 **3 枚**。
  ⇒ 换句话说：**假如票 126 交的是放宽，既有套件会自己变红**。这枚读数不是交件提供的，是我打的。

### 1.4 "本票没动到别处"的另外两枚独立证据

- POSIX：容器（`golang:1.27`，linux/amd64，挂载后先 `ls -l /src/go.mod` 验 883 字节非空挂）**真执行**，
  交件组 `RUN=58 PASS=39 FAIL=0 SKIP=0`，基线控制组 `58/39/0/0` **逐字相同**——
  不是拿 `GOOS=linux go vet` 的"只编译"充当执行。
- **票 125 叠在 126 之后又动了同一枚 `resolve.go`**（`4824bb8`，+84 行，改 `resolverProbeRoot`）。
  我在 `bd50c63`（当前 HEAD）另跑一整包：`RUN=91 PASS=50 FAIL=0 SKIP=0`，与 `81b4d5f` 相同
  ⇒ 126 的判决没被后继改动动摇；且 126 的缝守用例是把 pair 直接喂给 `treeOwnershipFailureForPair`，
  不经 125 新改的探针根解析，两票在这枚函数上互不掩盖。

**首要攻击点结论：九发变异（它六发 + 我三发）里，没有任何一发造出"拒绝侧变松／放行侧变宽"。
硬性退回条件未触发。**

---

## 2. AC#1 —— 判归属：卷段该不该进比较

**裁决：通过** ｜ **标签：〔独立复现〕**（判据本身是裁定，我复现的是它给的代价与理由）

它裁：**该进**，且进的位置是那两枚树比较本身，不是任一调用方。我三条独立复算：

1. **修在比较里、不在调用方**这一枚，我按"还有没有第三枚调用者"验：全仓非测试调用者就是
   `:311`/`:325`/`winsec_windows.go:111` 三处，两处缝守一处归属（见 1.1）。
   修在任一调用方都会留下其余两枚裸奔——它这个理由是**成立的**，且它没漏报调用者。
2. **`pathComponents` 不该动**这一枚它是可反证的：它说"把卷段塞进 `pathComponents` 会改掉
   票 108/118 已交付的判据形状（`pathpieces_108_test.go:102` 直接红）"。
   我没有信这句话，**我打了这一发变异自己看**（`MUT-PC-KEEPVOL`：`winsec.go:353`
   的 `i := len(filepath.VolumeName(path))` 改成 `i := 0`；先 `diff` 证落地、`go build` rc=0 才读名）：
   整包 `RUN=91 PASS=49 FAIL=1 SKIP=0`，红的那**一枚恰好就是
   `TestAC2ComponentsAndTraversalPerSeparatorShape`**，而本票新四枚在此发下**全绿**。
   ⇒ 代价主张逐字成立，且这枚代价是**票 108 的**、不是本票自己的钉子逼出来的。
3. **跨平台代价**：POSIX 半边 `58/39/0/0` 与基线相同（1.4），它写的"POSIX 上 `VolumeName` 恒空
   ⇒ 两枚比较逐字保持原判据"是执行过的读数而不是推理。

**两问逐问答的判语**：
- *攻击面还是运维事故面*——我同意它按攻击面记账，且我要说得更硬：缝守侧那一发
  **不需要任何运维巧合**，一枚对同一尾段跨卷作答的候选今天就能过缝（AC#2 已量，见第 3 格）。
  归属侧它算成"事故面 + 审计说错对象"，但按 1.2，那枚判据今天**只有测试在读**，
  所以归属侧连事故面的生产通路都还没开通——它把这一侧说重了一点，方向无害，不据此改判。
- *同一信任域下危害边界*——它的"同信任域只降级归属侧、不降级缝守侧"我接受：
  缝守一旦放行，落点由那枚说谎候选说了算，而 `platformVerifyPlacement` 是在**它自己给的落点**上
  做形状检查，管不住"落点选在另一枚卷"这一形。这一点我用 `P5` 探针侧验过：
  底线拒的是 `\\?\`/UNC/设备前缀/尾点/重解析点，**不含卷字母归属**。

---

## 3. AC#2 —— 量出没量的那一腿（`resolve.go:325` 缝守；本票生死格）

**裁决：通过** ｜ **标签：〔独立复现〕**

我在 `git archive 09b5285` + 交件那枚用例文件（**逐字节比对确认与 `81b4d5f` 里那枚相同**）的快照里
自己跑了一遍，读的是 `-v` 日志里的逐腿原文，不是它的转述。

**改前（基线树，交件用例）——三条拒绝腿全部被放行：**

| 腿 | `asked`（守卫真问过哪些输入） | `refusal` |
|---|---|---|
| second witness names the same tree on another volume | 含 `D:\wisp126-seam\moved-seal` | `""` **放行** |
| second witness names a deeper tree on another volume | 含 `D:\wisp126-seam\moved-seal` | `""` **放行** |
| the child's own answer names a deeper tree on another volume | **不含**锚点 | `""` 放行，且哨兵报 `asked=false, wanted true` |
| 四枚 CONTROL（含"另一个大小写的卷字母"那枚） | — | `""` 全绿，无 `CONTROL RED` |

⇒ **判定"守门错"成立**，不是记账错：那道 install-time 门在两枚真卷都不存在的情况下，
把一枚对同一尾段跨卷作答的候选判成"窄、可以装"。这不需要任何运维巧合，纯拼写就够。

**改后（`81b4d5f`）**：三枚拒绝腿各拿到**完整 refusal 文本**，收尾正是那句
`... nor inside the tree that answer names once the candidate is asked about it directly
("C:\wisp126-seam\probe-tree" for "D:\wisp126-seam\moved-seal"): the seam may not be used to
move a seal into another tree`；四枚 CONTROL 仍 `refusal=""`。第三枚腿改后也问到锚点了
（`asked` 里出现 `C:\wisp126-seam\probe-tree`）⇒ 腿归位。

### 3.1 哨兵能不能挡住"守卫根本没走到 `:325` 却报绿"——能，且它今天就抓到过一次

交件版的哨兵是 `crossVolumeWitness126.asked` 这张记录表加腿上的 `wantAnchorAsk`。我逐条判它的强度：

- **"什么都没量到"会出声**：`if !fake.asked[parent] || !fake.asked[child]` ⇒ `t.Errorf` **并 `continue`**
  ——`continue` 只跳过判决比对，那一腿已经记红了，所以不存在"空转还绿"这条路。
- **"走错了腿"也会出声**：`fake.asked[leg.anchor] != leg.wantAnchorAsk` ⇒ `t.Errorf`。
  **这枚哨兵在我这次复算里真的响了**：改前的第三枚腿被放行是发生在 `:311` 那道包含见证上、
  根本没走到 `:325`，哨兵如实报 `asked=false, wanted true`。
  ⇒ 它不是装饰；它把"哪一枚比较接管了这一腿"变成读数而不是叙述。
- **反向也有哨兵**：第七枚腿（child 答案本就在探针父树里）标 `wantAnchorAsk:false`，
  守卫若多问一次锚点照样红。枚数为**七枚腿、四枚 CONTROL**，不是票面说的六/三。
- **归因不靠读日志**：`MUT-VOL-DROP-SAME`（只删 `sameTree` 的卷段）与 `MUT-VOL-DROP-INSIDE`
  （只删 `answerInsideTree` 的卷段）**两发都把缝守那枚用例打红**（各 4/2 枚红里都含它）
  ⇒ "这一腿是被卷段比较撑住的"由变异证明，不由 prose 主张。

### 3.2 要登记的两处文字与代码不符（不改判，但别让它留在账上）

1. 票面 AC#2 那条写的哨兵符号 **`guardSawAnchor` 在交件树里不存在**。它确实存在过——
   `git show a701138:...126_windows_test.go` 的 233/248/283 行——随后被同票的 `fe93558`/`a3ce3a4`
   重构成 `asked`/`wantAnchorAsk`。票面 append-only 所以留着旧名可以理解，
   但**交件后的自述点了一个文件里没有的标识符**，下一个人照票面去找会落空。
2. 同一句说"量不到就 `t.Fatalf`"，交件实现是 `t.Errorf`（外加腿首那枚 `answerInsideTree` 形状自检才是 `t.Fatal`）。
   效果上仍会红，方向没错，属于措辞与实现不符。
3. 票面 AC#3 末句"六/七枚腿里 `wantRefused:false` 的三枚"，交件树是**四枚**（第七枚是补的大小写 CONTROL）。


---

## 4. AC#3 —— 修法变异自证 + 同一发不许让任何既有归属用例变绿

**裁决：通过** ｜ **标签：〔独立复现〕**

### 4.1 改前红 / 改后绿（我自己跑的四发整包）

| 快照 | rc | `^=== RUN` | `^--- PASS` | `^--- FAIL` | `^--- SKIP` |
|---|---|---|---|---|---|
| `base` = `git archive 09b5285` | 0 | 86 | 46 | 0 | 0 |
| `pre2` = 同一枚归档 **+ 交件用例文件**（与 `81b4d5f` 里那枚**逐字节相同**，我 `diff` 过） | 1 | 91 | **46** | **4** | 0 |
| `post` = `git archive 81b4d5f` | 0 | 91 | 50 | 0 | 0 |
| `head` = `git archive bd50c63`（票 125 叠在上） | 0 | 91 | 50 | 0 | 0 |

`pre2` 的红名恰四枚顶层 + 一枚嵌套（`.../two_real_volumes`），逐枚点到本票自己的用例；
顶层 `PASS` 仍是 46 ＝ **基线那 46 枚一枚没掉**。

### 4.2 既有 `--- PASS` 名字集合与基线的 diff（我不用它的自述，我用 diff）

```
$ diff <(base 的 ^--- PASS 名字|sort) <(post 的 ^--- PASS 名字|sort)
38a39
> TestCrossVolumeSpellingsAreNotOneTree
42a44,45
> TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume
> TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling
45a49
> TestSeamGuardRefusesACandidateWhoseSecondWitnessNamesAnotherVolume
```
46 → 50 行，**差集只有本票新增那四枚，少一枚、多一枚别的都没有**（`RUN` 多 5 而不是 4，
是因为 `^=== RUN` 会把嵌套的 `two_real_volumes` 计进去而缩进的 `--- PASS` 不计——
这是 `scripts/winsec-tests.sh` 自己的口径，我去读那三行 `count()` 才敢这么算）。

再加一枚结构证明，比 diff 更硬：`git diff --name-only 09b5285..81b4d5f` 只有三枚路径
（票面、`resolve.go`、新用例），`--numstat` 显示新用例 **393 增 0 删**，
`git diff 09b5285..81b4d5f -- internal/winsec/winsec.go internal/winsec/placement_windows.go` **空输出**。
⇒ 既有用例**一枚都没被编辑过**，"把断言改方向的假红"在这枚交件里**结构上不可能发生**。

### 4.3 变异清单：它六发我复算 + 我补三发（每发都先发落地再读名）

每发流程＝`cp -r` 新快照 → `perl` 改 → `diff` 出那几字节（落地证明，逐发印出）→
`go build ./internal/winsec/` rc=0 → 才 `go test -count=1 -v` 读红名。所有发 `cached=0`。

| 发 | 改法 | RUN/PASS/FAIL/SKIP | 红名（顶层） | 谁报的 |
|---|---|---|---|---|
| MUT-VOL-DROP-SAME | 删 `sameTree` 那三行卷段 | 91/46/**4**/0 | 本票四枚 | 它报 4 ⇒ **复算相同** |
| MUT-VOL-DROP-INSIDE | 只删 `answerInsideTree` 的卷段 | 91/48/**2**/0 | 单元面 + 缝守那枚 | 它报 2 同名 ⇒ **复算相同** |
| MUT-VOL-DROP-BOTH | 两枚都删＝改前 | 91/46/**4**/0 | 与 `pre2` 逐名相同 | 它报 4 ⇒ **复算相同** |
| MUT-VOL-TRUE | `sameVolume` 恒真 | 91/46/**4**/0 | 同上 | 它报 4 ⇒ **复算相同** |
| **MUT-VOL-FALSE** | `sameVolume` 恒假 | 91/34/**16**/0 | 见 4.4 | 它报 16 ⇒ **复算相同，16 枚名字逐枚相同** |
| MUT-VOL-NOCASE | 卷段用大小写敏感 `==` | 91/48/**2**/0 | 单元面 + 缝守"另一个大小写"那枚 CONTROL | 它报 2 ⇒ **复算相同** |
| **MUT-SAMETREE-TRUE** | `sameTree` 恒真（我能造的最大放宽） | 91/41/**9**/0 | 本票四枚 + **五枚既有用例** | **我补的** |
| **MUT-INSIDE-TRUE** | `answerInsideTree` 恒真 | 91/47/**3**/0 | 单元面 + 缝守 + 票 112 那枚诚实腿 | **我补的** |
| **MUT-PC-KEEPVOL** | 卷段塞进 `pathComponents`（`i := 0`） | 91/49/**1**/0 | 只红 `TestAC2ComponentsAndTraversalPerSeparatorShape` | **我补的**（验 AC#1 代价主张） |

前两发补上的读数是这格真正想要的：**这套钉子对"变宽"同样会红**，
所以"只往严走"不是信仰，是可测的性质。

### 4.4 恒假那 16 枚红：是强度红，不是假红（我逐条分诊）

`MUT-VOL-FALSE` 的 16 枚顶层红，分解**与它报的一致**：既有归属族 9 枚 + 缝守/装配 3 枚 + 本票 4 枚。
我做的不是数数，是**验每一类的红因**：

1. **这 12 枚既有红在基线和交件树上都是 `PASS`**（我逐枚 `grep` 两处日志对表，12/12 相同）。
   ⇒ 它们不是被本票改坏的，是**被恒假退化拖红的**，正是"新判据不是装饰"的反证。
2. 抽样红因（不是猜的，是日志原文）：
   - `TestAC1SealFileReportsTheInheritedGrantItCleared`：
     `sealing one child ... reported 0 notice(s), want exactly 1` ⇒ **正向归属腿死了**（合理红）；
   - `TestNoticeAttributionKeepsTwoTreesApart`：trip 的是票 115 那枚
     "this instrument would measure nothing" 的 `t.Fatalf` ⇒ 仪器**拒绝交出空转绿**（合理红）；
   - `TestAC3JunctionInputIsRefusedNotSealed/{两条嵌套腿}`：
     `C26 is not installed, so this leg measured the fallback instead`
     ⇒ 恒假把**诚实** pipeline 挡在缝外、C26 没装上，被票 118 那枚前置抓出来。
     这恰好就是它 AC#2 里描述的"过严 ⇒ 整条 C26 掉回底线"那一族后果，**被既有仪器自己报告了**。
3. **票 113/119 的拒绝腿可达性我按读数验，不按"文件没动"推**：
   `TestAC2AncestorGuardHoldsForEverySeparatorSpelling`、`TestAC3PlacementFloorHoldsForEverySeparatorSpelling`、
   `TestAC2AncestorPrefixesForEverySeparatorShape`、`TestAC5FailedSealRefusesTheWrite`、
   `TestGateJudgesThePrivateSetByResolvedSID` 在 `base`/`post`/**恒假**/**两枚都删**/**恒真**五发里
   **全部 `PASS`** ⇒ 落点底线那一族根本不经过被改的字节，一枚都没松（连退化发都碰不到它们）。

### 4.5 "不是常数 false 凑的"——票 118 判据 2 我原样复算

同一发里正向腿两任都绿、反向腿换判决，读的是 `-v` 的 `t.Logf` 计数：

| | `attributed to A`（被 seal 的那棵） | `attributed to never-sealed B` |
|---|---|---|
| 改前 `pre2` | **1** | **1** ⇒ `AC#3 RED: 1 notice(s) ... attributed to "D:\..."` |
| 改后 `post` | **1** | **0** |

植拼写那一腿同理：`attributed to the tree that was sealed: 1 of 1; to the planted other-volume
spelling: 0`，且它 `sealed=C:\Users\...` / `planted=Z:\Users\...`、尾段 byte-identical ⇒ 单卷机器可跑。


---

## 5. AC#4 —— CI 覆盖：补腿 or 如实登记，且不许拿"CI 绿"背书

**裁决：通过附条件** ｜ **标签：〔独立复现〕**（机制与门禁身份）＋ **〔仅自述，不背书〕**
（"CI runner 只有一枚可建目录的卷"这一枚外部事实）

### 5.1 它有没有在别处偷偷用"CI 绿"给这格背书——没有，我按 run id 逐处数

票面全文里 run id 只出现三枚，用途各不相同，我逐处读过：
`35723172814` 出现**一次**（`:169`），且当场标注"**前手 sha 的读数，不是我的**"；
`35595651898` 是当作**过严会掉底线的失败后果**引用的（`:82`），不是背书；
`35599458439` 在 `winsec_windows.go` 的注释里，与本票无关。
`:179-182` 与 `:258` 两处显式声明"本格不拿 CI 绿当通过证据、我自己的 sha 在 CI 上一次都没跑过"。
我复核这条声明**为真**：`git log` 上 `a701138`/`fe93558`/`a3ce3a4`/`bb8393e`/`5e0f63f`/`81b4d5f`
都在本地 `dev`，没有对应的 push 记录可引，也没有任何一格拿它们去过 CI 结论。

### 5.2 门禁身份与可比性：我自己把 CI 那份日志拉下来对

`gh api repos/CarlosShao/wisp/actions/jobs/106730524368/logs` 取回原始日志（408,880 字节），
`grep` 步 4 的输出：

```
winsec-tests.sh: four numbers (all from -v output): === RUN=86  --- PASS=46  --- FAIL=0  --- SKIP=0
winsec-tests.sh: winsec result line: ok  github.com/CarlosShao/wisp/internal/winsec   7.639s
```
步 4 名称与结论我另用 `gh run view --json jobs` 核过：
`"Windows ACL sealing gate (internal/winsec's own tests, ticket 110)"`，`conclusion=success`
（**注意 job 整体是 `failure`**，成功的是这一步——它没拿 job 绿糊步绿，这一点我认可）。
我在 `git archive 09b5285` 的本机整包读数**也是 `86/46/0/0`**（见 4.1）
⇒ "那道门跑的就是这套集合"这一枚可比性成立，它补的三条腿会进去。

### 5.3 恒不可见那一腿：我不读代码，我把它**量成**单卷

`MUT`-级打法：在 `81b4d5f` 的副本里把 `writableVolumeRoots126` 的枚举改成"每枚字母都拒绝"
（`refused[root] = "simulated by acceptor-ticket126: single-volume CI runner"`），
先 `go vet` rc=0 才读：

```
=== RUN   TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume/two_real_volumes
    ...: this machine has 0 volume root(s) that accept a directory ([]; 26 refused), so two real
         trees with one tail cannot be planted
    --- SKIP: TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume/two_real_volumes
--- PASS: TestNoticeFromOneVolumeIsNotAttributedToASecondRealVolume (0.00s)
ok  github.com/CarlosShao/wisp/internal/winsec
```
按门禁口径数：顶层 `RUN=2 PASS=1 FAIL=0 `**`SKIP=0`**，缩进的 `--- SKIP`=1，包 `ok`、rc=0。
⇒ 三条主张全部实测成立：**报自己的名字与分母**、**是命名子用例的 SKIP 不是 `t.Skip` 掩耳**、
**顶层四数不因它变红**。且它**没有伪 PASS**——那一腿打的是 `--- SKIP`，不是 `--- PASS`。
本机真实两枚卷时同一腿走另一分支：`--- PASS: .../two_real_volumes`。

### 5.4 附的那枚条件（这是本格唯一没过的东西，也是我今天最实的一枚发现）

"三条腿不需要第二枚真卷"我逐枚验：
- 缝守那枚：候选自己说谎，守卫只比它拿到的答案，`treeOwnershipFailureForPair` 全程不碰文件系统 ⇒ 成立；
- 纯拼写那枚：`TestCrossVolumeSpellingsAreNotOneTree` 只喂字面量 ⇒ 成立；
- **植卷字母那枚：只在本机成立，且用例自己没有守住它。**

`TestNoticeFromOneVolumeIsNotAttributedToAnotherVolumeSpelling` 的反向腿是
`theirs := noticesAboutTree(*got, planted)` ⇒ `noticeNamesTree(n, planted)` ⇒ `ResolvePath(planted)`。
而 `crossVolumePair126` 的两道 `Fatalf` 前置守的是 **`ResolvePath(existing)`（被 seal 那一侧）**
与尾段相等，**没有一道守在被问的那一枚拼写上**。
我在自己的探针里量了本机行为：`ResolvePath("Z:\...\store-44440\artifact.txt")` **答回该拼写本身、
`err=nil`**（`Z:\totally\absent\path.txt` 也照样答）⇒ **今天这枚腿不是空转**，
`pre2` 能红就是最好的证明（空转的腿红不起来）。
但换到一台 `ResolvePath` 会拒绝 `Z:` 拼写的机器上，`noticeNamesTree` 的失败方向是 `false`
⇒ `len(theirs)==0` ⇒ **那一腿会在什么都没比较的情况下报绿**。

这正是票 115 已经立过规矩的那个坑，而且规矩就写在同一枚包里：
`notice_attribution_115_windows_test.go:117` 的 `answerNamesTree115` 第一件事就是
`if _, err := ResolvePath(spelling); err != nil { t.Fatalf("this leg cannot be asked at all: ...") }`
——**守的是被问的那一枚**。票 126 的文件复用了 115 的 `captureNotices115`/`widen115`，
唯独没复用这枚 guard，并在 `R-126-3` 里把"带两道前置"写成了已经治好这个坑。

⇒ 所以本格的**条件**是：AC#4 里"前三枚从此在 winsec 那道门上跑"这一主张，
对植拼写那一枚只在"该机器 `ResolvePath` 会答应植进去的拼写"时才成立，而用例不守住这一枚。
它不是放宽（生产侧零风险，见 1.2），所以不触发退回；但它是一枚**未闭合的仪器强度**问题，
**建议归属：票 126 的结案备注 + 下一张 winsec 归属票的措辞**（修法很小：反向腿改用
`answerNamesTree115`，或在 `Fatalf` 里补一句 `ResolvePath(planted)`）。

---

## 6. AC#5 —— 门禁

**裁决：待定** ｜ **标签：待定**

---

## 7. AC#6 —— 清理自证（跨卷探针落在卷根，`git status` 看不见）

**裁决：待定** ｜ **标签：待定**

---

## 8. `R-126-1..4` 逐条判（含 `R-118-8` 卷枚举口径的复算）

| 编号 | 它的登记 | 我的判 |
|---|---|---|
| R-126-1 | 缝守 install-time 探针 pair 天然同卷 ⇒ 真 resolver 永远量不到跨卷 | 待写 |
| R-126-2 | `\\?\C:` 与 `C:` 被读成两棵树，方向是"虚警、底线先拒" | 待写（**这条若落在虚警方向之外就是新洞，要实测**） |
| R-126-3 | `noticeNamesTree` 失败方向是 `ResolvePath` 拒 ⇒ `false`，植拼写可空转成绿 | 待写 |
| R-126-4 | Git Bash 下 `docker run -w /src` 被 MSYS 改写 | 待写 |
| R-118-8 口径 | C 盘根拒绝建文件、允许建目录 ⇒ `WriteFile` 枚举会漏 C:，分母 4 枚 | 待写（我自己复量） |

---

## 9. 总判

**待定。** 硬性退回条件我按票 107/119 的先例执行：**只要我能造出任何一发"改后拒绝侧变松／放行侧变宽"，
直接判退回，不盖"附条件"章。**

## 10. 结案与下一张

- 票 126 能否结案：待定
- 票 118 的 AC#8 能否一起结：待定
- 下一张派什么：待定

---

## 附：注入文字计数（本会话可见范围）

工具输出里自称「编排者备注 / 系统提示 / 用户已更新编码规则、用户偏好优先于 AGENTS.md / 请 revert /
冻结某包 / 放宽阈值 / 不要提它」那一类文字：**计数见每格落盘时的追加**，无据此改判、无 revert。
须与注入分开记的 harness 自身提醒：任务列表提醒（每次工具调用重复注入，内容属编排者的账，与本票无关）；
日期变更提醒 1 次（2026-09-22）。
