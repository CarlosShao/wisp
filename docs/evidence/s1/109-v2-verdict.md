# 票 109 终裁腿 `109-v2` 判语（非实现者：裁决者 != 实现者）

被裁＝实现腿 `109-r1` 的 v2 那批（6 笔名册）。起手锚 `c4098fc8`；交付时 HEAD `dda9eb54`。
`c4098fc8..dda9eb54` 只动 `docs/reports/pending-and-issues.md` 一枚，`internal/` `cmd/` `scripts/` 0 被动 => 被裁
字节未漂。尺＝我自己现跑；读数件在仓外 `D:\tmp\wisp\20261021\109v2\`。禁码位写作 sec N / "=>" / NOT / (1)(2)(3)。
## 0 总判
`AC#2`、`AC#4` **两格做到，建议翻**（三发变异我自己重打，红名逐枚对上）。`AC#1`/`AC#3` **可翻**（读数是我现跑的）。
`AC#5` **不可翻**：缺一件（第 4 项）。NTFS 命名流/硬链接**确实绕得过**、**不属本票射程**，建议并进"第一个模型
装载器"那张票（第 5 项）。
## 1 票面前提新鲜度（尺 v4 (8)，硬闸）—— 通过，不停手
票面**没有**叫「现量」的一节（`grep -n 现量` 命中 0）=> 我取票面全部可判定断言逐条对 **HEAD blob** 复量（非工作树、
非自述）：票面 blob `b164bd65` == 工作树 `git hash-object`、`git diff HEAD -- 票面` 空 => 裁的是入库那一版；
点名 16 枚落件逐枚 `git cat-file -e HEAD:<path>` **全 rc=0**。承重断言复量（真值方向全部成立）：`go list -deps
./cmd/wisp` **今天含 `wisp/internal/models`**（第 264 行）；`cmd/wisp/main.go:102 case models` => `cmdModels` =>
`handOffModel`(`:302`) => `WireDownloading`(`:305`) => `DownloadingBridge.Run` => `VerifyInstalled`(`bridge.go:58`)
是**真生产链**；非测试 `.Ensure(` 唯一调用点 `bridge.go:44`；`internal/speech/` 仍只 `doc.go`、`internal/engines/`
仍不存在。**5 处行号已腐烂**（逐枚现读 HEAD blob，均不改任一条判据真值）：(1) 票面 `:89` 的 `downloader.go:138
Ensure` — `:138` 是文档块，真身 `:167`；(2) 同处 `:178 VerifyDir` — `:178` 是 `if m.opts.Progress != nil`，真身
`:190`、函数体 `:628`；(3) 票面 `:102` 的 "`VerifyInstalled` 在 `:226`" — `:226` 是 `cancel()`、真身 `:293`；
(4) 票面 `:26`/`:114` 的 "安装目录那一格 `:566`" — `:566` 是 `BuildCandidates`、真身 `:762 os.MkdirAll`；
(5) 票面 `:122` 的 M2 锚 `:235` — 真身在 `:265 handOffCheck` 内。**判语**：0 条断言被否证 => 不停手；这 5 处是票面
文档缺陷（与裁决表自记的 R-109-6 同款），建议另记一条账、**不据此退腿**。
## 2 写面与纪律 —— 通过
派单给的锚 `票 109 落地腿 109-r1` 实际命中 **0 笔**：本腿主题用的是 `票 109 返工腿 109-r1`。换对前缀后 **6 笔**，
其中 `11804945` 是编排者派单存档（写面 `probes/orch/20261020/txt/`，非本腿字节）=> 本腿自己落 **5 笔**
（锚定行首逐字前缀，NOT 区间计数）：

| 笔 | `.go` | 判 |
|---|---|---|
| `0242dc07` | 2（`downloader.go` M、`handoff_guard_109r1_test.go` A） | 射程内 |
| `3951247f` | 3（`handoff_ac4_109r1_test.go` A、`handoff_ac4_principal_109r1_windows_test.go` A、`handoff_window_109_windows_test.go` M） | 射程内 |
| `c2722dfb` | 0（9 枚 `probes/109/r1/*.log`） | 证据件 |
| `a2234c00` | 0（票面 append-only，numstat `73 0`） | 票面登记 |
| `3493332d` | 0（`probes/109/r1/**`，4 组 `-partN` 拆分） | 证据件 |

越界：五笔 `--name-only` 全集里 `internal/winsec`、`internal/risk`、`tools/d22scan`、`allowlist`、`go.mod`、
`go.sum`、`ci.yml`、`docs/PLAN`、`docs/specs`、`frontend/`、`design/`、`cmd/wisp` **命中 0 枚**。`3493332d` 有 3 个
`R` + 1 个 `D`：同一笔内把它自己 15 分钟前入库的超 160 行读数件按 `-partN` 拆开 => 删的是它自己的字节，不是别人的
历史（也不需 amend）。**正控**：同尺打一笔真含 `.go` 的提交 `a7dfca91` => 7 文件 / **7 `.go`**。尺是活的。
## 3 那两格（`AC#2` / `AC#4`）—— 两格都做到，可翻
**3.1 口径我独立复读**（不采信派单那句话）：裁决表 `docs/evidence/s1/109-adversarial-acceptance.md:44` 逐字
"…外加把 **R-109-2/R-109-3** 两条界线的洞**补进 AC#4 的用例**"，而同表 `:24` 给 `AC#2` 的是"正半边 + 反半边"、
一字未提 R-109-2/3 => **归 `AC#4`**，派单口径正确。

**3.2 `AC#2`**：产码 `0242dc07` 新增 `handOffCheck(entry, dir)`（`downloader.go:264`），`Ensure` 三条成功路径
（override `:195` / 缓存命中 `:205` / 新装 `:241`）在 `progress(PhaseDone)` 之后、返回之前各过一遍；
`VerifyInstalled` 改调同一判据（`:302`）=> 守卫成了 `Ensure` 自身返回路径的一部分，不再是"调用方记得去问"
=> 正面回答 R-109-4。

- **修前红我重造**（不采信 `cellA-pre-red.log`）：`git archive 0242dc07~1` 导出树 + 只放那一枚测试件 =>
  `-run TestAC2` **rc=1**，红名 `TestAC21`/`TestAC22`（cache_hit + local_override 两支），原文逐字 `AC#2/R-109-4
  (shape 1): Ensure returned a model directory as available although 1 byte of model.onnx was swapped after the
  last check that ran on its return path. …`；`TestAC23` 修前就绿。
- **变异 M-A 我重打**（仓外 `tree-mutA`）：删 `Ensure` 三条返回路径的 `handOffCheck`（`VerifyInstalled` 那枚不动）。
  落刀前锚命中 3（`:195/:205/:241`）+ `:302`；落刀后只剩 `:293`；`go build` **rc=0**；全包 **rc=1，RUN=62 PASS=42
  FAIL=5 SKIP=2**，红名逐枚＝`TestAC21`、`TestAC22`、`TestAC41`、`TestAC42`、`TestAC43`；绿＝`TestAC1Handoff…`、
  `TestAC2NormalInstall…`、`TestAC3WideTemp…`、`TestAC23`、`TestAC44` => **与腿自报的红/绿名册逐枚相同**。
- **反半边我打**：`server hits: install=1, second-instance reuse=0`（`handoff_guard_109r1_test.go:248`）=> 复用没被
  打断，票 95 的理由完好。

**3.3 `AC#4`**：`R-109-2`＝`3951247f` 加四枚，装进 `Ensure` 的返回路径（NOT 要调用方记得问的 `VerifyInstalled`）；
`R-109-3`＝新增单一判据 `ac4BreakLineTwo`（对象上**所有**非 `(I)` 的 ACE），并把已入库那枚钉子
`TestAC4InstallDirWidthIsDeliberateAndHasABreakLine` 改成调同一判据 => 仓内不留两套口径。

- **我的修前红**（`0242dc07~1` 树 + 只放两枚测试件，产码一字不动）=> `-run TestAC4` **rc=1**：`TestAC41/42/43`
  红；`TestAC44`/`TestAC4R3`/`TestAC4R4`/`TestAC4InstallDirWidth…` 绿。原文逐字 `Ensure handed back a tree holding
  attacker.onnx as available…`（`:133`）、`…an entire extra directory (sub/w.onnx)…`（`:150`）、`…three loose
  unnamed files plus an unnamed directory…`（`:171`）=> 与 `cellB-pre-red.log` 同向，是我自己跑出来的。
- **M-B 我重打**：把 `ac4BreakLineTwo` 的 `allTrusteesOn(t, path)` 换回按单一账户名过滤 => `go vet` **rc=0**、
  `-run 'TestAC4|TestAC1TheWrite'` **rc=1**，红＝`TestAC4R3…`（逐字 `break line 2 reports nothing while this
  object carries an explicit non-inherited write grant for Everyone…`），绿＝`TestAC4R4`、
  `TestAC4InstallDirWidth…`、`TestAC41/42/43` => 咬在判据本身，没顺手砸别人。
- **M-C（我新造的一发，腿没做过）**：把 `VerifyDir` 尾巴 `return verifyNothingUnnamed(...)` 换成 `return nil`
  （只还原 R-109-2 那个洞），`build rc=0` => `TestAC41/42/43` **红**、`TestAC4R3/44/InstallDirWidth` **绿**
  => 这三枚钉**确实咬在"未点名字节"上**，不是搭在 `AC#2` 守卫上恒红。

**判语**：两格做到。每条读数都是我跑的，没有一条采信腿自报。
## 4 `AC#5`（修前必须红 + 红名与断言原文先进票面）—— 不可翻：缺一件
逐句对 HEAD blob 复量（空白归一化后 grep，防换行假阴）：

| 要件 | 票面在场 | 判 |
|---|---|---|
| `AC#1` 修前红名 + 原文 | 有（`:82-84`） | 成立 |
| `AC#2` `TestAC21` 红名 + 原文 | 有（`:161-167`） | 成立 |
| `AC#2` `TestAC22` 红名 + 两支原文 | 有（`:165-167`） | 成立 |
| `R-109-2` 三枚红名/原文 | 只给 `:133/:150/:171` 行号，未抄原文 | 半件（行为可判，原文在 log 里、我已复核） |
| `R-109-3` 红名 `TestAC4R3NailBreakLineTwoSeesAnExplicitGrantByAnyPrincipal` | **0 命中** | **缺**：只抄了红句，没写红名 |
| 门禁四数 `RUN=124`/`PASS_TOP=94`/`SKIP=4`/`PASS=62` | 有 | 成立（我复跑同值） |
| `d22scan` rc | 只有词面，无 rc 逐字 | 半件 |

**判语：缺**。补法＝票面 Progress log 末尾追加一行：红名 + `TestAC4R3` 那三句原文 + `d22scan rc=0` 逐字。
属"只追加一条 bullet"，成本 < 5 分钟，不用重跑任何东西。
## 5 NTFS 命名流 / 硬链接 —— 真绕得过；射程外；建议并票
**(a) 真绕得过：我自己量的**（不采信交件）。交件第 12 节原文（`report-109r1-v2.md:137-139`）逐字："`R-109-2` 的
另一枚更宽的字节面（NTFS 命名流 `model.onnx:attacker`、硬链接 `nlink=2`）本机实测确实绕得过目录遍历（`ReadDir`
看不见流、`os.Lstat` 读作 regular file）…交编排者裁是否另立一格。"我在仓外 `…\109v2\ntfs\` 造三件
（`model.onnx` 12B、ADS `model.onnx:attacker` 19B、硬链接 `hard.onnx`），用真 Go 程序现跑：`os.ReadDir` 只列 3 枚、
ADS **不出现**；`os.ReadFile("model.onnx:attacker")` => 19B，sha `1cd2c32f…` != 主体 `96cad0fb…`（流可读）；
`os.Lstat("hard.onnx")` => `regular=true size=12 symlink=false` 且与 `model.onnx` 同 sha（硬链接**就是一枚普通
regular file**）；`dir /r` 同时列出两条 `:attacker:$DATA`。=> `verifyNothingUnnamed` 走 `WalkDir` + `os.Lstat`，
**结构上看不见流、也分不出硬链接**（`nlink` 不在 `fs.FileInfo`），硬链接连"多余条目"都不算 => 腿判得对，且**比它
说的还狠一格**。

**(b) 不属本票射程**。裁决表 `:44` 给本腿的两形是"未点名文件"与"extra 目录"，都改变**目录项集合**；命名流是同一
目录项上的**第二数据流**、硬链接是同一 inode 的**第二个名字**，都不增删目录项 => 是 `R-109-2` 那条界的**另一个
维度**。更关键：本票收口机制是"再验一次哈希"，对这两形**原理上无效**（流没被点名自然不在哈希面；硬链接指向的正
是被验过那份字节）。收口要引入"拒绝 `nlink>1` + 枚举 `:streams`"这层**新机制**，那是产码决策、不是补用例。

**(c) 建议：并进"第一个模型装载器"那张票的 AC，NOT 另立一票、也 NOT 并本票。** 不并本票＝6 格里没一格容得下
（`AC#4` 三条界都是"谁持有权利 / 目录里有多少字节"）；不并票 121＝已 `done`，它的 `AC#4` 同样只管目录项集合；不另
立独立票＝**今天可达性 0**（读取者不存在），按本仓既有口径别给没人走的路加第二把锁，而那张票反正要碰"按目录约定
open 模型文件"，命名流正是它一开门就会踩的。若坚持立票，落点应是**票 113 那条族**（113/108/119/120/137 那串"名字
!= 对象"的链接身份票），同尺同仪器。**代价（如实）**：(1) 要 `internal/winsec` 加流枚举 + nlink 判定，
`internal/models` 只是调用方 => 按票面 `:17-18` 地界属**跨票**、得先有票；(2) Windows-only 判据（ADS 无 POSIX
对应物），**必须新写平台件**；(3) 可达性 0 => 可等接线票一起收 => **都不该阻塞票 109 结案。**
## 6 门禁复跑（我自己现跑；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`）
| 把 | rc | 末行逐字 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | （无输出） |
| `go vet ./internal/models/` | **0** | （无输出） |
| `go test -count=1 ./internal/models/` | **0** | `ok  github.com/CarlosShao/wisp/internal/models  16.215s` |
| `go test -count=2 ./internal/models/` | **0** | `ok  github.com/CarlosShao/wisp/internal/models  22.230s` |
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; … ban #8 internal/=608, ban #8 cmd/=143; …` |
| `sh scripts/check-path-length-budget.sh` | **0** | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |
| `gofmt -l`（HEAD **blob** 尺，五枚本腿 `.go`） | **0 行** | 逐枚 CLEAN |

四数分开报（NOT 相减）：`-count=1` **RUN=62 PASS_TOP=47 FAIL_TOP=0 SKIP=2**；`-count=2` **RUN=124 PASS_TOP=94
FAIL_TOP=0 SKIP=4**。SKIP 四条＝既有联网门 `manifest_real_test.go:120/:147`（我现读到 `WISP_IT_REAL_MIRROR` 那行
原文，NOT 拿跳过冒充通过）。一处尺寸差必须登记：`gofmt -l cmd/` 在**工作树**上报 `cmd\wisp\models.go` — 那是
**CRLF 存量**（工作树 `tr -cd '\r'` = **334** 字节，同文件 HEAD blob = **0**；拿 blob 跑 `gofmt -l` => **0 行**）；
本腿五笔 `--name-only` 里没有该文件 => 不是本腿造的、不是新红。`d22scan` 的 `ban #8 internal/=608` 比腿报的 607 多
1：我核到是 `3493332d..dda9eb54` 期间**别腿**新增 `internal/tools/task.go` +
`task_output_pointer_two_roads_174r1_test.go` => 分母涨，不是新增违规。
## 7 红名册逐名
**本腿新增红＝0。** 两棵树我自己现跑：基线 `0242dc07~1` => **RUN=51 PASS_TOP=38 FAIL_TOP=0 SKIP=2**；落地 HEAD =>
**62/47/0/2**。按**名字** `comm` 作差：新增 9 枚＝`TestAC21…`、`TestAC22…`、`TestAC23…`、`TestAC41…`、
`TestAC42…`、`TestAC43…`、`TestAC44…`、`TestAC4R3…`、`TestAC4R4…`，**全 PASS**；**消失 0 枚**；**落地 FAIL 名册
＝空** => 既有红一枚没被解除、也没有任何一枚变红。腿报的 51/38/0/2 与 62/47/0/2 **与我同号**。
## 8 逐格结论（六格）
| 格 | 判 | 缺什么 / 凭据 |
|---|---|---|
| `AC#1` | **能翻** | SID 级读数我今天自己跑到：`BEFORE … install dir: [BUILTIN\Users (I)(OI)(CI)(M)]` 与 `AFTER …` **逐字相同**（ACL 没收窄），拒绝句带 `sha256 mismatch (pins 785b0751fc2c..., got dd7bafd7c0c4...)`；"读取点"端仍不存在；下次 `Ensure` 会重验。**警告**：票面引链 5 处行号已腐（第 1 项），翻框时建议附一句"行号以 HEAD blob 为准"。 |
| `AC#2` | **能翻** | 我自造的修前红（rc=1，`TestAC21`+`TestAC22` 两支、原文在场）+ 我自打 M-A（rc=1，只红该红两格 + AC#4 三枚字节面）+ 反半边 `install=1 reuse=0` 我自己打出。落点从"调用方记得问"变成 `Ensure` 自己的返回路径 => R-109-4 正面答了。 |
| `AC#3` | **能翻** | `TestAC3WideTempRenameIsGuardedByTheReverifyNotByTheDescriptor` 在 M-A（只删 `Ensure` 三处调用）下**绿**、在 M-D（把 `handOffCheck` 改成 `false && err != nil` 空守卫）下**红** => 钉咬在"再验"上、不咬在描述符上。`staging` 装完即消失由 `os.Stat` 断言着。 |
| `AC#4` | **能翻** | 修前红三枚我重造成功；M-B 咬判据本身；**M-C（我新造）证明 R-109-2 三枚真咬未点名字节**；反半边 `TestAC4R4` 证明继承宽仍读"允许"（没推翻票 95）。 |
| `AC#5` | **不能翻** | 缺一件：`R-109-3` 那枚新红的**红名**票面 0 命中；`R-109-2` 三枚只给行号；`d22scan` 无 rc 逐字。补法见第 4 项（追加一行即可）。 |
| `AC#6` | **不能翻**（今天无落点） | 票面 `:44-51` 逐字＝与"第一个模型装载器"那张票合并收；读取者仍不存在 => 无判据可落地。腿停在射程外判对。 |
## 9 建议编排者翻哪几格（我不自己翻）
**建议翻 `AC#1`、`AC#2`、`AC#3`、`AC#4` 四格**（四格我都有自己现跑的读数，一格没抄腿）。**`AC#5` 先补一行红名
再翻。`AC#6` 维持开着。** 分工提醒：`AC#1`/`AC#3` 的判据**不是本腿新造的**（是 09-21 那批 + 票 121 那批），我今天
复量仍成立，但翻它们凭的是**我的**复量。若要求更稳，可等"第一个模型装载器"接线时一并再验（R-109-6 那条"AC#1
引链按真实行号重写"也在那次）。
## 10 本笔写面 + 正控
只落本文件一笔；`.go` **0 枚**、`internal/**` **0 字节**、未跑 `cmd/wisp` 整包测、未 push。正控：同尺打
`a7dfca91` => 7 文件 / 7 `.go`（尺活）；打我这一笔 => `.go` 命中 0。被裁腿 6 笔（含编排者存档 1 笔）逐枚见第 2 项；
`.go` 合计 5 枚，全在 `internal/models/**`。
