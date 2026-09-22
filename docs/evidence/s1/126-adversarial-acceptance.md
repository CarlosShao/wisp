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

**裁决：待定** ｜ **标签：待定**

读数：待写。要复算的三件事：
1. 自造 fake resolver `crossVolumeWitness126` 是否真把 C26 缝守量成"改前放行跨卷候选"（`refusal=""`）；
2. `guardSawAnchor` 哨兵——它在 `a701138` 版存在，交件版已被同票重构换名，我要判新哨兵是否同等强、
   能否挡住"守卫根本没走到第二见证却报绿"这种空转；
3. 六/七条腿里三枚 `wantRefused:false` 的 CONTROL 是否改前改后都绿（不许有 `CONTROL RED`）。

---

## 4. AC#3 —— 修法变异自证 + 同一发不许让任何既有归属用例变绿

**裁决：待定** ｜ **标签：待定**

要复算的三件事：
1. 改前红四枚 / 改后绿（`pre2` 快照独立跑）；
2. 既有 `--- PASS` 名字集合与基线 `diff` **只多本票四枚**（它报 `86/46 → 91/50`）——我用 diff 证，不听自述；
3. `VOL-FALSE` 那发它报红 16 枚，其中既有腿 12 枚 ⇒ 逐名核对这 12 枚是**强度提升的合理红**
   还是**断言被改方向的假红**（后者直接退回）。

---

## 5. AC#4 —— CI 覆盖：补腿 or 如实登记，不许拿"CI 绿"背书

**裁决：待定** ｜ **标签：待定**

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
