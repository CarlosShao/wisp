# 357 — 守卫验的目录必须等于读取者真正 open 的目录（「第一个模型装载器」归口票）

- Status: **DONE（2026-10-22 10:0x 编排者收口＝六格全闭、文件名加 `-done`）** ~~**ready-for-agent**（`AC#0`／`AC#1` 是只读现量，~~
- 立案：2026-10-21 14:3x +08 编排者立（**机主当场批**：「1. 可以」＝为「第一个模型装载器」立票）。
- **为什么现在立**：票 109 的 `AC#6` 票面逐字写着"与「第一个模型装载器」那张票合并收"，但**盘上从来没有那张票**
  （立案时刻读数：`grep -rln '模型装载器'` 只命中票 109 自己；**今天已命中 2 枚**＝109＋357，故该句只可当立案时刻读数引）
  ⇒ 两味（`109 AC#6` 的"守卫验的目录 ≠ 读者开的目录"＋裁腿 `109-v2` 新量到的 NTFS 命名流/硬链接面）**长期悬空**。
  本票就是那个归口。
- 非实现者签名＝**`109-v2` 判语**（`docs/evidence/s1/109-v2-verdict.md`，提交 `5ae2bae4`）：它在 **sec 5(c)** 里判  ← 更正：立案时误写"第 2 条"（`357-a1` 抓到、我复核成立）
  "NTFS 命名流/硬链接该另立，但建议并入『第一个模型装载器』那张票的 AC"。
- 关联：票 109（`AC#6` 迁入本票；票 109 其余五格已 `[x]`）· 票 121（`-done`，`AC#2` 是让该分支**在生产进程里可配落地**的那一步）· 票 119（生产路纪律）· 票 348／342（同族"拒绝路径"面；★更正：**两枚都不是 `-done`**，`Status:` 逐字"已立，未派"——立案时我把 322-v3 那批"已批落地"误记成"已结案"）

## 涉及的两味（人话）

**味一：守卫验的目录 ≠ 读取者开的目录。**
模型下载有个"验哈希"的守卫。今天它验哪个目录，是**跟着 config 走**的——有 `local_override` 时验的是**覆盖目录**，
而 CLI 同一行**仍打印 `store=…models`**。于是可以出现"store 里一个字节都没有，却报与已验签清单一致"。
（验证方法：去掉那行 config 后，同一枚被篡改的 store ⇒ `rc=1` 并点名 `sha256 mismatch`。）
⇒ 要的判据：**守卫验的目录必须严格等于读取者真正 `open` 的那个目录**。

**味二：NTFS 命名流与硬链接绕过哈希面（裁腿 `109-v2` 仓外实测、比原判更狠一格）。**
用真 Go 程序量到：`os.ReadDir` 只列 3 枚、**ADS 不出现**；但 `os.ReadFile("model.onnx:attacker")` **读到 19 B** 且 sha 与主体不同；
`os.Lstat("hard.onnx")` ⇒ `regular=true size=12` 且与主体**同 sha**。
⇒ **硬链接连"多余条目"都不算**（列不出来、也不进 `nlink` 之外的任何账），命名流则"列不出来但读得到"。
⇒ 结论：**本票"再验一次哈希"对它有原理性的无效**——校验的是主体文件，攻击者读的是另一条流。

## 结案判据

- [x] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把味一量成成对读数——`VerifyInstalled` 今天验的目录 ⇄ CLI 打印的 `store=…` 那一行 ⇄ 读取者真正 `open` 的目录，三者**逐字**摆出来；有 `local_override` 与无 `local_override` 两种 config 各一发。
  > **2026-10-21 16:3x 编排者翻框（凭普查腿 `357-a1`，我逐条复跑）＝成立，且**落点已改向**：
  > ①正控先跑成立：无 override、store 换成**同长 12 B** 篡改 ⇒ `rc=1` 并逐字点名 `model.onnx: sha256 mismatch (pins 206402cab345..., got 950edc0ad300...)`。
  > ②**无 override：三者同一条**；**有 override：不是同一条**——守卫验 `…\override`／CLI 打 `…\dataB\models`（`models.go:253` 的 `ms.root`）／`Ensure` 真交还 `…\override`。
  > ③**B2 反证守卫没错**：override 被篡改 ⇒ `rc=1`；store 被篡改/不存在 ⇒ 仍 `与已验签清单一致` **rc=0** ⇒ **破裂的是 CLI 的 `store=` 那一行，⛔ `VerifyDir` 的实参**。
  > ④★**根因（我独立核过）**：`internal/models/bridge.go:44` 逐字 `_, err := b.mgr.Ensure(ctx, id, b.onProgress)` ⇒ **交还值被丢掉**；`Ensure` 全仓**唯一产码调用者**就是这一处；
  > `cmd/wisp/models.go:323` 那行 `已交还` 打的是 **`ms.root`**、⛔ 交还值。
  > ⑤⇒ **本格判据改向为「交还值即契约」**：未来任何要读模型字节的调用方**必须拿 `Ensure` 的交还路径、⛔ 自己拼 `store/<root>/<id>`**；并把 `_, err :=` 这处丢掉写成**具名欠账**。
  > ★这是**判据方向**改动、⛔ 契约文字改动（D 表无此条）⇒ 编排者拍；若后续要动 `PLAN.md` 的字，那一步才摆机主。
  > ⚠ 味一今天的**可观测后果只是主诉**（CLI 报"一致"而 store 是篡改的）；"读者 open 错目录"这一后果**可达性 0**（`cmd/`+`internal/` 对 store 字节的 open 两把独立尺都 0 处）。
      ★必须**先跑正控**证明"去掉那行 config 后 store 被篡改会 `rc=1`"（否则你的台子测不出发育不良）。
- [x] **AC#1 现读＋归因（只读）**：把味二量成**可复跑的读数件**（ADS 可读且 sha 不同／硬链接同 sha 且 `Lstat` 报 regular／`os.ReadDir` 两者都列不出），
  > **2026-10-21 16:3x 编排者翻框＝成立，且它比 `109-v2` 多量到一格**
  > 仓外真 Go 程序 ＋ 真 NTFS（⛔ mock）：`ReadDir`/`WalkDir` **列不出 ADS**；`ReadFile` 读到 ADS **19 B** 且 sha 与主体**不同**；`Lstat` 把 ADS 与硬链接**都读作 regular file**；`nlink=2` 且 fileIndex 相同。
  > **产品级复量**（真 `wisp models verify`）：**D1 主体+ADS ⇒ `与已验签清单一致` rc=0（静默通过）**／D2 目录内硬链接 ⇒ 被票 121 未点名判定抓 rc=1／D3 目录外硬链接 ⇒ rc=0／D4 主体换成指向外部字节的硬链接 ⇒ `sha256 mismatch` rc=1。
  > ★**新量到的一格**：目录外种链再从外部名字写入 ⇒ `ReadDir`/`dir /r`/`nlink` 全程无痕，**但哈希面第 4 步抓到**（`got a33238fcb04b…`, rc=1）⇒ **硬链接不是一律绕过；真正静默通过的只有 ADS**。
  > **ADS 耐久性**：`os.Rename`（产品主路 `:771`）**流跟着走**；Go `io.Copy`（`copyFile` `:781` 兜底）**流丢掉** ⇒ ADS 能否进 store 取决于 install 走哪条；安装目录 ACL 逐字 `Authenticated Users:(I)(M)` ⇒ D3 不是纸面形状。
      ★**在仓外用真 Go 程序造三件**（主体／ADS／硬链接），⛔ 用 mock；并判"哪一层该拦"（哈希面？目录项面？读取面？）。
- [x] **AC#2 拍形（编排者＋机主）**：三形各写代价＝〔甲〕只堵味一（守卫目录对齐）；〔乙〕只堵味二（枚举流＋`nlink` 判定）；〔丙〕两者都堵。
  > **2026-10-21 16:3x 编排者拍形＝〔丁〕新增的第四形**（⛔ 甲乙丙任一）：
  > **〔丁〕`wisp models verify` 报"与清单一致"时，若目录里存在"清单没点名的第二数据流"（ADS），必须明说一句**——把"静默通过"（D1）变成"**说了再通过**"。
  > **理由（按 `AGENTS.md` ⓪⓪⓪ 铁律第 2／3 条）**：①**最有价值**——一次改动就把**今天唯一真静默的那味**变成可观测；甲只改一句打印、乙是新平台件换零收益、丙在可达性 0 下买不到价值；
  > ②**复杂买得到真价值**——这是**报告**⛔ **新增拒绝**：⛔ 新造 Windows-only 平台件、⛔ 改哈希面语义、⛔ 让既有绿变红；
  > ③★**⛔ 选乙**：`FindFirstStreamW` 枚举＋`nlink>1` 判定是 Windows-only 新平台件，而它保护的对象**今天不存在** ⇒ 按铁律第 3 条"复杂必须买得到真价值"，这里买不到；
  > ④★**乙那套读数不入废**（普查腿已在本机跑通）⇒ 留作本格"将来形"证据，等真有读取者再启用。
  > **判据（写给写腿）**：`verify` 在 rc=0 的路径上，若枚举到 ADS ⇒ 输出**具名一句**（含流名与"清单未点名"），**rc 仍 0**（⛔ 变红）；若没枚举到 ⇒ 输出与本改前逐字相同（**⛔ 让既有绿变噪**）。
      ★**A 档代价必须写明**：味二要动 `internal/winsec` 加流枚举与 `nlink` 判定＝**Windows-only 新平台件**、且今天**可达性 0**
      （读取者不存在）⇒ 本仓既有口径是"别给没人走的路加第二把锁"。**未拍形任何实现都⛔ 开工。**
- [x] **AC#3 落地＋牙（写腿）**：按拍定那一形做；每处改动配**一枚今天响、不修不响**的用例；
  > **2026-10-22 10:0x 编排者翻框＝成立（〔丁〕形落地；三把尺我现跑）**
  > **新件**：`internal/models/ads_windows.go`（`GetFileInformationByHandleEx`＋`windows.FileStreamInfo`＋手映 **24 B** `FILE_STREAM_INFO`，与 `internal/ball/win32_windows.go:160` 同形）
  > ／`internal/models/ads_other.go`（`!windows` 恒返空、**文件头具名写 POSIX 那枚只是控制腿⛔ 牙**）／`internal/models/ads.go`（`UnnamedStreams` 走 **`resolveInstalledDir`**）。
  > ★**裁决书第 3 条那个坑在码面上被钉住**：`ads.go:46` ＝ `dir = m.resolveInstalledDir(id)`（**override 目录**），⛔ `store.root`；
  > 且`cmd/wisp/models.go:254` 那句 `store=…` **逐字未动**、`:262` 换成 `announceUnnamedStreams357`、**rc 仍 0**。配两枚**成对**用例：`Test357VerifyTailReportsTheVerifiedTreeNotTheStoreRoot` ＋ 配对负例 `Test357StoreRootNamingIsNotWhatTheReportUses`。
  > ★**`ERROR_MORE_DATA` 重试**：起点 4096 B、×4 增长、**上界 16 MiB**、撞界**返具名错误**（⛔ 静默当空）——正是裁决书点名的那条（不加重试会退化成**假阴**）。
  > ★**变异 1**（渲染成空）⇒ 4 枚 CLI 用例红逐名报；★**变异 2**（把调用摘掉、渲染函数留着）⇒ **只有** `Test357VerifyTailIsStillOnThePassingPath` 红、其余五枚全绿
  > ⇒ **这发正是"钩子尺"的价值**：没有它，摘钩后整包全绿而命令退回静默；★**变异 3**（去掉重试）⇒ **`200 planted -> 0 enumerated`**＝**复现了假阴**。
  > ★**三读数**：还原后 `ads.go`／`ads_windows.go` 两处 sha256 三处同串；`models.go` 归一 CRLF 后三处同为 `B0FE3890…603`（本机 autocrlf 检出形态、入库 blob **389 LF／0 CR**）。
  > ★**正例 G 仍绿**：stdout **116 B 逐字节相同**（`SequenceEqual=True`、两枚 sha256 `B6B2D4AA…E45`、CR=0）＋ `Test357VerifyTailWithoutAStreamPrintsNothing`（stdout/stderr 都是空串）。
      ★反向判据＝把你新加的那一跳摘掉 ⇒ 用例必须转红并**逐名报红行**；先证变异落地（`grep -n` ＋ `go build` rc=0）再读数。
- [x] **AC#4 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1` 动过的包（四数分开报、⛔ 相减）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
  > **2026-10-22 10:0x 编排者翻框（六把门我自己现跑）＝成立**
  > `GOFLAGS= go build ./...` **rc=0**／`go vet ./internal/models/ ./cmd/wisp/` **rc=0**／`go test -count=1 ./internal/models/` **rc=0**（`ok`）／
  > `sh scripts/d22scan.sh` **rc=0 clean**／★**格式甲尺（入库 blob 树）0 枚**。
  > ★**`Test357*` 六枚我单独同发跑：全 PASS**（`VerifyTailSaysOneSentenceAndKeepsRC0`／`VerifyTailWithoutAStreamPrintsNothing`／`VerifyTailReportsTheVerifiedTreeNotTheStoreRoot`／`StreamNoticeIsTheOneSentence`／`VerifyTailIsStillOnThePassingPath`／`StoreRootNamingIsNotWhatTheReportUses`）。
  > ⚠ **`cmd/wisp` 整包 `rc=1` ⛔ 本票的伤**：红名 **1 枚** `TestAC14GoSideEvalPushReachesThePage`＝**既有红**（台账 `:20471` 早就具名）＋**死点 1 枚**（`A1100` 同族）⇒ **新增红 0**。
  > ★★**path-length 门禁一开始 RED，根因是**我**、⛔ 它**：我改名 `181`／`275` 加 `-done` 后**超预算**，而名册里写的是**旧名**（`275` 根本没条目）。
  > 我已修名册（`181` 旧名→新名 1 处 ＋ 补 `275` 条目）⇒ 现跑 `over-budget=61 / covered by roster=61 / not in roster=0`、**VERDICT GREEN**。
  > ★这是**同一坑第 2 次**（`A1076` 记过第 1 次）⇒ 定式：**改名加 `-done` 后必须同笔把名册一起改**，且**立刻现跑门禁**（⛔ 等别人报）。
      ★`cmd/wisp` 整包必抬 `-timeout 30m`（台账 `A1079`）；⛔ `CGO_ENABLED=0`（本仓必然编不过）。
- [x] **AC#5 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 109 任何原句（互指只在本票尾追加）。
  > **2026-10-22 10:0x 编排者翻框（纪律格由编排者判，四问逐问现跑）＝成立**
  > ①**越界**：票 357 全部腿笔（`357-a1` 普查／`357-r1` 落地含停手件与交付件）逐笔 `git show --name-only` 判——
  >   `3843d651`／`aeea4e5b`／`876aa377`／`fc9ffa95`／`47a631be` **越界合计 0**；其中 `aeea4e5b` 是**唯一带 `.go` 的一笔（8 枚全 `.go`）**。
  > ②**零 push**：`origin/dev..HEAD` ＝ **346**（非空）⇒ 所有腿笔**只在本地**。
  > ③**⛔ 改票面**：`git log --all -- 'issues/357-*.md'` 的**非编排者**提交 = **0 枚** ⇒ 票面只我动过。
  > ④**⛔ 改台账**：同上，非编排者提交 = **0 枚**。
  > ★**两条停手都是正当的**：第一发撞"派单前提事实上错了"（`FindFirstStreamW` 在 x/sys v0.48.0 **不存在**，我独立核实）⇒ **停手对**；
  > 续跑发**零手写 syscall**（用 `GetFileInformationByHandleEx`）、**未动 `internal/winsec`**、**未加依赖** ⇒ 三条停手条款均未触发。
  > ★**它另主动登记一处自己的错**：`cmd/wisp` 新用例件初版真被 `gofumpt` 咬一格（`"…tell "+` 应写成 `"…tell +`），已改、改后 0 命中 ⇒ **自报不掩**。

## 硬约束

- ⛔ 重述票 118／120-r3 已批的"link 只 unlink、⛔ 跟随"（两票均 `-done`）。
- ⛔ 碰 `docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden、`allowlist.txt`、`tools/d22scan/**`、`frontend/**`、`design/**`。
- ★**Windows-only 面**：味二的所有读数必须在真 Windows 上取（ADS 是 NTFS 特性）；POSIX 侧只要求**不退化**。
- ★**写面互斥**：`internal/winsec/**` 同包只许一枚写腿；`cmd/wisp/**` 同波只许一枚。

## 排程

`357-a1`（只读＝`AC#0`＋`AC#1` 两味现读）→ 编排者＋机主拍 `AC#2` 三形 → `357-r1`（写腿＝`AC#3`）→ `357-v1`（非实现者裁）。

**Status 追加（2026-10-21 14:3x）**：号位＝编排者定（`354`/`355`/`356` 已占，`357` 现量空号）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 14:3x +08] agent=`编排者` did=立案（机主当场批「1. 可以」；非实现者签名＝`109-v2` 判语第 2 条；把票 109 `AC#6` 与 NTFS 味一并挂进 `AC#0`/`AC#1`）next=派 `357-a1` 只读出两味现读；`AC#2` 三形代价表回来后摆机主拍形。
