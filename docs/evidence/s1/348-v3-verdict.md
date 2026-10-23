# 票 348 终裁腿 `348-v3` — `AC#2` 能翻（两平台成对三件，两发皆本腿现跑）

**身份**：终裁腿＝**非实现者裁腿**（裁决者≠实现者）。⛔ 落任何 `.go`、⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`／台账、⛔ push。
**起手锚**＝`7e99ccf81`（`git rev-parse HEAD` 现认，与派单同串）。**被裁对象**＝`348-r2` 两笔（`5ffa553e3`／`c3ae02e29`）＋`348-r1` 两笔（`195ce6d0f`／`e9caddff`）。
**台件**＝自建三棵**仓外**导出树 `D:\tmp\wisp\20261023\348v2\trees\{head,pre,nos}`（`head`＝`git archive HEAD`，17612 件；`pre`＝`head` 只换 `winsec.go` 为基线 blob `67783ae0`、删 `deleteseam_*`；`nos`＝`head` 只换 `winsec.go` 为 `757b16a8`＝有落地无接缝）＋自写仓外台架 4 枚（`probe\`）＋`-overlay` 三份。★**工作树零突变**（`git status --porcelain -- internal/winsec` 起手＝空、收手＝空）⇒ **无突变窗口、窗口期零提交**（`A1159②`）。

## ★0 交件落点：派单写的路径已被占用（具名上报，非静默改）

派单写 `docs/evidence/s1/348-v2-verdict.md`；该名**已是 `342-v3` 的判语件**（`e8d6d23a2` 交件②＋`a7eb956ff` 追加更正，39 行／12113 字节，票 348 尾 bullet 逐字引过）。系列命名是「第 N 次终裁」而非腿名（`348-v1` 由 `342-v2` 交、`348-v2` 由 `342-v3` 交）⇒ 本件落**下一空号** `docs/evidence/s1/348-v3-verdict.md`，**不销毁既有判语**（建新件可逆、覆盖旧件不可逆）。若编排者要原名，`git mv` 一句即得。

## ① 票面前提新鲜度（对 HEAD blob 复量）

| 票面前提 | HEAD 现量 | 判 |
|---|---|---|
| 密封失败支四行含 `_ = os.Remove(resolved)` | 基线 `fde1df02` 命中 **2**，**正在 `:160`／`:165`**（票面行号对基线 blob**逐字成立**）；HEAD 命中 **0**，换成 `removeRefusedLeaf(...)`（`:198`／`:203`，行号因落地而位移） | ★**HEAD 上过期**（落地所致，应然） |
| `os.Remove` 吞错误 | 落地形仍吞（Windows `_ = p.root.Remove(base)`／POSIX `_ = unix.Unlinkat(...)`） | 成立 |
| `oNoFollowPlatform` | `winsec_other.go:29` ＝ `syscall.O_NOFOLLOW`；`winsec_windows.go:335` ＝ `0` | 成立（逐字） |
| `winsec` 只有两枚平台件 | 现 **5** 枚（＋`pinnedleaf_342_{unix,windows,fallback}.go`） | 过期（`348-v1` 已更正） |
| 接缝 | `winsec.go:128` 逐字 `var removeRefusedLeaf = unlinkPinnedLeaf`，`:198`/`:203` 两处调用 | 成立（现读） |

## ② `AC#2` 三件 × 两平台 —— ★**逐格，两发都是本腿现跑**

★**先说清"复用算不算现跑"**：票面逐字只要「**两平台各一发**（⛔ 一发代两发）」＝**平台覆盖**要求，⛔ 要求"同一腿各跑一次"；且验收的唯一实质条件是**量的对象是落地那支**。判据：POSIX 那一跳的**函数体**在两个被测树与 HEAD 之间**逐字节相同**（`unlinkPinnedLeaf` 的 unix 体＝`_ = unix.Unlinkat(int(p.fd()), base, 0)`＋`runtime.KeepAlive(p.dir)`，`9581a48c` ⇄ HEAD 同串；签名只多了票 353 加的**未使用** `pinned` 形参）⇒ **按"量的对象同一性"判，复用成立**。★但本腿**没让这条判语承重**：我自己把 POSIX 那一发打了（`docker`＋`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读，`TMPDIR`/`GOTMPDIR` 留容器内避免 311 守卫）⇒ **两平台两发都是本腿现跑**，复用只作旁证。

| 平台 | ① 修前分子 > 0（本腿现跑） | ② 修后 0（本腿现跑） | ③ 同一发证窗口到链接态（本腿现跑） | 判 |
|---|---|---|---|---|
| **POSIX**（容器） | 确定性那枚 3/3 红、红句逐字 `the foreign witness is GONE after the refusal-path unlink (…data/someone-elses-tree/leaf-342r1.bin: no such file or directory)`；竞态 `deleted=99 / 83 / 77 of 1500`（`provedLinkCycles=200/286/246`、`pinnedWalkRefusals=33/47/37`、`pinnedParentBeforeOpen=0/0/0`＝基线本来没有父钉） | `deleted=0 / 0 / 0 of 1500`；确定性那枚 3/3 绿、survivor 逐字 `size 75 (was 75), mode -rw-rw-rw- (was -rw-rw-rw-), mtime moved=false` | 同一行日志里 `provedLinkCycles=280 / 294 / 298`、`pinnedWalkRefusals=57 / 69 / 59`、`pinnedParentBeforeOpen=33 / 42 / 31`（⛔ 拿 0 命中当凭据） | **成立** |
| **Windows**（本机） | 真基线字节＋**本包既有** `applyDescriptor` 缝注入：`deleted=9 / 9 / 7 of 400`（删除支真跑 `131/136/111`、`provedReparse=80/92/100`）★另：我自建的 `T2′`（修前整条路由＋接缝）`deleted=400/400/400` | 落地＋接缝：`deleted=0 / 0 / 0`（`seamCalls=400`、`hopLanded=400`、`hopRefused=0`）★旁证：赛跑仪器同三发 `deleted=0`、`provedReparse=94/90/95`、`formBLegRefusals=29/38/42` | 同发 `hopLanded=400`（翻转由 `reparseAt` 当场证明）＝`AC#2` 逐字「同一发证明窗口真到达过链接态」 | **成立** |

★**判语：`AC#2` 能翻**。缺什么：**无缺**。
★**「复用算不算现跑」的判语（单独一句）**：**算**——理由是票面要的是"平台覆盖"＋"量的对象同一性"（已用字节同一性证明），⛔ 是"同一腿各跑一次"；★并且本腿已自己补跑，所以这条判语**即使被⛔ 承认也不影响结论**。

## ③ 「Windows ① 本来就取得到」这条自我更正 → **成立**

本腿自己造的两发（⛔ 采信它的日志）：**不注入密封失败** ⇒ `deleted=0/0/0`、`truncated=7/17/6`、`provedReparse=108/86/82`，而**删除支只跑了 2/1/0 次**（口径＝名册里"a seal refusal on this route"那一行的计数）⇒ `348-r1` 的 `deleted=0` **是仪器形状（删除支几乎没跑）**，⛔"世界不可达"；同一棵真基线字节**注入密封失败** ⇒ `deleted=9/9/7 of 400`。⇒ **①在 git 基线字节上本来就取得到、只是不稳**，接缝买的是**确定性（400/400）与可归因**，⛔"可达"。
★**精度补一条**：`348-r1` 的交件其实已写过"缺的不是窗口，是拒绝路径"（`348-r1-impl.md:23`），而且它自己的 ①′ `mut` 已经打出 `deleted=3,1,1,0,0,0,1,1`（8 发 5 红）⇒ 所以真正被推翻的是**结论"①取不到"**，⛔ 它的机制分析。
★**我的 `T2′` 与它的 `T2` 不是同一棵树**：`T2` 是 `348-r2` 在 `D:\tmp\wisp\20261023\348r2\trees\T2`（22156 字节）；`T2′` 是我从基线 blob 自建（20646 字节，只加接缝＋两处调用改名）。两棵的**形状**相同（修前整条路由＋接缝），读数同为 `400/400/400`，⛔ 说成"我复跑了它那棵树"。

## ④ 「第二道人工批准对'①能不能取到'本来不必要」→ **成立，但射程要写窄**

- **对"①>0"不必要**：我上面 9/9/7 与 `348-r1` 的 ①′ 都证明不加缝也能 >0。
- **对 `AC#2` 三件整体也不必**：①（9/9/7）＋②（赛跑 `deleted=0` ×3）＋③（`provedReparse=94/90/95` 同发）本可凑齐。
- **它真正买到的三样**（这三样不是"可达"）：⑴ ① 从 ~2% 摆动变成 **400/400 确定性**；⑵ ③ 从"赛跑侧证明翻转过"升级为"**删除那一跳的瞬间祖先真站在链接上**"（`hopLanded==seamCalls` 地板）；⑶ 反向判据成为**可常驻、不抖**的钉子（否则仓里那枚钉会 flaky）。⇒ 编排者那句"它实买的是确定性与可归因"**逐字成立**；"对'①能不能取到'本来不必要"**成立**。

## ⑤ `T4` 承重归因 → **成立，并经本腿第五臂加固**

`348-r2` 的 T4＝保留钉、只把删除体换回 `os.Remove(filepath.Join(p.dir, base))` ⇒ `deleted=400/400/400`。★**但 T4 一口气换掉了两枚机制**（`unlinkPinnedLeaf` 的身份复验＋fd-相对 act），单发**归因不出**是哪一枚承重。我造了第五臂把两枚分开（`-overlay`，仓外）：
- **armA**（复验＋fd-相对全去，纯按名删）⇒ `deleted=400/400/400`、`hopLanded=400`；
- **armB**（**保留** `p.root.Lstat(base)`＋`os.SameFile(pinned, now)` 复验，只把 `p.root.Remove` 换成按名 `os.Remove`）⇒ **仍然 `deleted=400/400/400`、`hopLanded=400`**。
⇒ 复验**不救**：它本身也是 fd-相对地问（`root.Lstat(base)` 问的是被钉那个对象，与 `pinned` 永远自洽），所以按名落刀照样被搬走。**真正承重的是"fd-相对"那一半**——`348-r2` 的结论成立，且现在有分离测量支撑（⛔ 只是推断）。

## ⑥ 接缝是不是"最小面" → **是**（逐行核过，⛔ 顺手重构）

`git diff 5ffa553e3^:…/winsec.go 5ffa553e3:…/winsec.go` ＝ **`23 2`**，拆开＝**21 行注释 ＋ 1 行 `var` ＋ 2 处调用改名**；`openPinnedLeaf`/`unlinkPinnedLeaf` 的签名一字未动（`pinned` 形参是票 353 的，⛔ 本腿）；⛔ 新导出名、⛔ 新参数、⛔ 改 `privateFile` 返回语义。**整笔提交只有 3 件**（`A` 交件、`A` 新测试件、`M winsec.go`），⛔ 碰任何别的文件（`git show --name-status` 逐枚见 ⑨）。

## ⑦ 门禁 rc（四数分开报；★跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`）

★★**两枚锚必须分开报**：起手锚 `7e99ccf81`（build／vet／gofumpt／整包三发都在它上面跑）⇄ **收手时 HEAD＝`47987d32`**（跑完我又回头复跑了 build／vet／d22scan／pathlen）。**`internal/winsec` 在两锚之间逐字节未变**（`git diff --stat 7e99ccf81..47987d32 -- internal/winsec` 为空）⇒ **winsec 的读数与 `348-r2` 两笔文件全部照原样结转**。

| 门 | rc | 读数 | 锚 |
|---|---|---|---|
| `go build ./...` | **0** | — | 两锚都 0 |
| `go vet ./internal/winsec/` | **0** | — | 两锚都 0 |
| `GOOS=linux` / `GOOS=darwin` 各 build＋vet | **0/0/0/0** | 四发全 0 | `7e99ccf81` |
| `gofumpt -l internal/winsec`（**甲＝blob 树**／判据） | **0** | 0 行；`[\\/]` 归属核命中 **0**（＝真没列文件，⛔"列了被过滤"） | `7e99ccf81` |
| 同上（**乙＝工作树**／只诊断） | **0** | 0 行 | `7e99ccf81` |
| `sh scripts/check-path-length-budget.sh` | ★**0（锚A）／1（锚B）** | 锚A：`VERDICT GREEN`、`tracked=17612`、`over-budget=61`、`in roster=61`、`not in roster=0`；锚B：**`VERDICT RED`**、`over-budget=62` 而 roster 仍 **61** | 见下 |
| `sh scripts/d22scan.sh` | **0** | 锚A：正控 `PASS=35 FAIL=0 SKIP=0`、`RUN=77`、真扫 290 枚产码体 `clean`；锚B：正控同号、真扫 **`clean - no D22 ban violations`** | 两锚都 0 |
| 整包 `go test -count=1 -v ./internal/winsec/` | **1** | 三发四数（RUN/PASS/FAIL/SKIP）＝**`123/121/2/0`**、**`123/120/3/0`**、**`123/121/2/0`**；每发 `PASS+FAIL+SKIP==RUN` 逐发对上（⛔ 相减） | `7e99ccf81` |

★★**`pathlen` 锚B 那枚红归谁（用"能区分两棵树"的尺）**：唯一新增越界件＝`.scratch/wisp/issues/365-two-taskLogState-…-on-the-other.md`（relative **126**、末名 **105**、full **170** ＞ hat 165）⇒ `git log -1 -- <该路径>` ＝ **`8beec66b1`（编排者 A1166：立票 365）** ⇒ **属票 365 的立票件、⛔ 票 348**；票 348／`348-r2` 两笔在该锚下仍 **0** 条越界。★**这枚也说明"门禁绿"必须连锚一起写**：同一把尺在同一小时内由 GREEN 转 RED，红因⛔ 本票。
★★**我自己抓的一处尺的错（`A1168` 第 12 例同一个坑）**：第一遍我跑 `sh scripts/check-path-length-budget.sh | tail -4`，`$?` 取到的是 **`tail` 的 rc=0**，读成 `pathlen_rc=0 GREEN`，而脚本真 rc 是 **1**。⇒ **不带管重跑**才是真值（本表已按重跑值写）；`d22scan` 同理**不带管**取 rc。

`d22scan` 归因（★两锚都说）：锚A 时编排者预告的"今天会 rc=1"**已不复现**——phantom-citation 当时**只在工作树**（`git show HEAD:…/schema.go` 那行是 `);` ⇄ `Get-Content` 那行含 `docs/evidence/s1/196-r1-impl.md`）；★**收手时它已随票 196 的落地笔（`6e19f4a43`）进 HEAD，且被引件 `196-r1-impl.md` 现已跟踪**（`git ls-files` 命中 1）⇒ 幻影两个方向都消掉，锚B 仍 **0**。⇒ 该红自始至终属**票 196**、⛔ 票 348；`348-r2` 当时报的 `rc=1`／`PASS=33 FAIL=2` 在其后**自然转绿**（它的归因口径是对的）。

## ⑧ 新增红 0 —— **我的核**（⛔ 采信它的名册）

三发红名逐枚（上表）：`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（＋子测试 `symlink-dir`，`AC#1 RED: 3 of 400` / `1 of 400`＝**票 342 的 `truncated` 残留**）、`TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（票 342-r2）、`TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`（★**票 348 的叶子腿**，但红的是**仪器地板**：`AC#2 RED instrument: no leaf leg answered (guard=0, landing's re-check=0) while reparseAt proved 91 flipped leaf cycles`，`deleted=0`，⛔ 损伤红）。
★**那一枚 348 红的归因我自己做了 A/B**：同一条用例在 `head`（有接缝）**5/5 绿**、在 `nos`（`757b16a8`＝有落地无接缝）**4/5 绿＋1/5 红**、红句与 `head` 那发**同一句** ⇒ **接缝前的既有 flake，⛔ `348-r2` 的新增红**。⇒ **新增红 0 成立**，但★**名册「红名全是票 342」这句的射程不够**（有 1/3 发出现票 348 的**仪器地板**红）。

## ⑨ 三条残余 —— 逐条判

①**派单预期要改**（Windows ① 本来就取得到）→ **对**（见③）；但★更准的写法是「① 在基线字节上取得到、**不稳**（我 9/9/7）；接缝买确定性，⛔ 可达」。
②**`d22scan` 要 196 的腿落件才会绿** → **半对，今天已过期**：真正让它是红的不是"文件缺"而是**工作树的 `internal/memory/schema.go` 引了它**；文件现已落盘（未跟踪）⇒ **今天 rc=0**。⇒ 派单/台账若还写"等 196 落件才绿"应改为「等 196 那枚**跟踪**进仓」。
③**`deleteseam` 的平台断言**（地板 `hopLanded == seamCalls`，现 `400/400`）→ **对，且这条地板是这枚仪器最值钱的一行**：将来若"父目录内换名"不再落地，它会**变红**＝**仪器丢了世界**的报警，⛔"删除变安全了"。★我补一句射程：这条地板只能证明"翻转在我们的 `toFlip()` 这一形下站住了"，⛔ 证明"所有祖先变换都被拒/都被放过"。

## ★不声称的（照实）

1. ⛔ 声称 POSIX 的 ①②③ 与 `348-r1`/`342-v2` 逐位同号——我的竞态分子是 **99/83/77 of 1500**（它的 76/62），同量级、⛔ 同号；这类赛跑读数**不该**期望逐位复现。
2. ⛔ 声称 `H2c`/`H2f` 那对矛盾已归因（`348-r2` 照实标"⛔ 归因、⛔ 当判据"）；**处置正确**——把互相矛盾的两发**排除在判据之外**并写出来，比硬选一个可信。★本腿也⛔ 拿它当判据。
3. ⛔ 声称 Windows 的 ① 在本仓**常驻门禁**意义上已钉住：`deleteseam` 的 T2/T4 形态是**仓外**读数，仓里常驻的只有 T3 那一发（绿）。要常驻"①必红"的反向钉子，得另立实现笔。
4. ⛔ 声称任何 `.go` 被我改过：本腿零 `.go` 字节（`git status --porcelain -- internal/winsec` 收手为空）。
5. ⛔ 声称 `348-r2` 的提交正文行数（102/140）已改正——`c3ae02e29` 只在件里更正，提交正文仍是旧数（⛔ amend，**以 `c3ae02e29` 的 `wc -l` 真值 144/151 为准**）。

## 提交与写面

- 本腿提交：**`d1ca3e643`**（第一笔，`A` 本件）→ **`48554984c`**（第二笔，把第一笔的号补进本件）→ **本笔**（第三笔，门禁两锚复跑＋我自己那处尺的错的更正）。三笔都只碰 `docs/evidence/s1/348-v3-verdict.md` 一枚。
- 逐枚点名：`docs/evidence/s1/348-v3-verdict.md`（`A`）。索引核验 `git diff --cached --name-status` 现跑＝**只此一行**，⛔ 出现本笔不该有的 `D`/`A`。
- 正控：同一把尺打 `5ffa553e3`（`A docs/…/348-r2-impl.md`／`A internal/winsec/deleteseam_348r2_windows_test.go`／`M internal/winsec/winsec.go`）与 `195ce6d0f`／`e9caddff`（两枚 `.go`）⇒ 尺**认 `.go`**，⛔ 空白尺。
- 越界：**0**（`internal/**` 之外零字节；`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**` 一字未动）。
- hash 口径（核前先问"谁算的／什么算法／几位"）：`winsec.go` `git hash-object`＝**SHA-1/40**＝`fc4e0d696ac482951902a0b82dcf9429284bc33b`、`sha256sum`＝**SHA-256/64**＝`6bed1b5bfae9bf37f7b5312e7040a290a874246406f5eb4756a1b42e75a6e52e`；`deleteseam_348r2_windows_test.go`＝`3272fdace3b8ffb21ea77768c62332bc9f7d449b`／`41d1121d3dd41759afc1b1cbdafc4d9c50999d9ff55266130f20ab5e11eeb88f` —— **四枚与 `348-r2` 交件逐字同串**（⛔ 混比）。
- `.go` 行数/CR 尺（★逐枚，尺＝`wc -l`／`tr -cd "\r" | wc -c`）：`internal/winsec/winsec.go`＝**467 行／23066 字节／CR=0**；`internal/winsec/deleteseam_348r2_windows_test.go`＝**151 行／7778 字节／CR=0**（后者与 `c3ae02e29` 更正后的 151 对上）。

## 结语

`AC#2` **能翻**：两平台三件逐格成立，且**两发都是本腿现跑**（POSIX 容器／Windows 本机），"复用"判语另有字节同一性支撑但已不承重。
`348-r2` 的两条推翻**都成立**（① 本就取得到／根因射程不够），`T4` 承重归因**成立且被第五臂加固**，接缝**确是最小面**；门禁在**锚A 全绿**（build／vet／两平台 build＋vet／gofumpt 甲乙／pathlen／d22scan 全 0），**锚B 除 `pathlen` 一枚红外仍全绿**，而那枚红 **归因票 365 的立票件（`8beec66b1`）、⛔ 本票**；`d22scan` 的 rc=1 已自然转绿、归因票 196；新增红 **0**（★名册射程要收窄：有票 348 的**仪器地板** flake，已 A/B 证明是接缝前既有）。
★**我自己的一处尺的错也已入件**：`pathlen | tail -4` 的 `$?` 是 `tail` 的 rc（读成 0 GREEN），不带管重跑才是真值 1（`A1168` 第 12 例同一个坑）。
三条残余：①对（写法要改成"取得到但不稳"）②半对且今天过期（要等的是**跟踪进仓**）③对（那是"仪器丢世界"报警）。
