# 115 — 对抗验收七项检查表（`docs/specs/SPEC-10-testing-acceptance.md` 第 8 节逐格）

> **这张表是什么**：票 230 `AC#1`（`230-four-cells-left-unfinished-inside-closed-tickets.md:19`）点名欠的那件东西——
> 票 115 名下头一枚裁决表。出表腿＝`115-v1`（只读、非实现者、未参与票 115 任何一笔落地）。
> **出表时刻** 2026-10-20 16:3x–17:0x +08｜起手锚＝`git log -1` 现认 `c16ff47ec465749f1977efa94a2b024cc8931e27`（分支 `dev`）。
> **本表⛔ 结案票 115、⛔ 翻任何框**：翻勾归编排者凭另一枚非实现者判语（票 115 与票 230 的框本腿一枚没碰，终量见 §6）。

## 0. 七项标题的出处与一处具名偏离

标题逐字抄自现读的 `docs/specs/SPEC-10-testing-acceptance.md` 第 8 节（`:105-119`，本腿 `grep -nE "^#{2,3} "` 现定位，⛔ 凭记忆、⛔ 拿 AGENTS.md 薄索引当权威）。
**一处偏离照实登记**：第 7 项原文写作「**失败预演报告**（§7）」，其中的分节符是码位 **U+00A7**；
派单公共块明写入库正文⛔ 写出该枚字面字符 ⇒ 本表把那一处改写成「第 7 节」，**其余六项标题一字未动**。
这是"逐字抄"与"⛔ 写出被禁字符"相冲时取的口径，请编排者裁这一处是否要另记一笔。

## 1. 先定真分母（⛔ 采任何旧叙述，含派单里编排者自己那条）

尺＝对票 115 票面现跑（文件＝`.scratch/wisp/issues/115-seal-notices-…-done.md`；本腿先验该枚文件**工作树 ⇄ HEAD blob 逐字节相同**，尺＝`git hash-object` ⇄ `git rev-parse HEAD:<path>`，再量的）：

| 把尺 | 读数 | 射程 |
|---|---|---|
| `grep -cE '^- \[ \]'`（带框·未勾） | **2** | AC#4 `:58`、AC#6 `:64` |
| `grep -cE '^- \[x\]'`（带框·已勾） | **3** | AC#2 `:47`、AC#3 `:52`、AC#7 `:68` |
| 尺：行首 `- ` 且含 `**AC#<n>`（**判据格**，不数框） | **7** | AC#1 `:44`、AC#2 `:47`、AC#3 `:52`、AC#4 `:58`、AC#5 `:61`、AC#6 `:64`、AC#7 `:68` |
| 尺：行首 `- ` 且含 `**AC#<n>` 但**不带框** | **2** | AC#1 `:44`、AC#5 `:61`（两枚行首已被换成⛔＝本格不占分母） |

**票面全文还有 2 枚 awk 命中（`:103`、`:307`）＝Progress log 里引用 AC#2/AC#3 的交件叙述，不是判据格**；拿"全文命中 9"当分母会多算两枚，本腿具名剥掉。

### 1.1 三个数彼此不合同——哪两对对不上、差在哪（派单点名要先报的矛盾）

| 叙述 | 它给的数 | 今天盘上真读数 | 差的枚数与原因 |
|---|---|---|---|
| 编排者 15:2x 现量 | 未勾 **2**／已勾 **3**／**AC 标签行 5** | 未勾 2／已勾 3／**判据格 7** | 「未勾 2⇄已勾 3」这两半**对上**；对不上的是**「标签行 5」⇄「判据格 7」＝差 2 枚**＝AC#1（`:44`）与 AC#5（`:61`）。那两枚的行首早已不是 `[ ]`/`[x]` 而是⛔（摘框），所以"数带框行"的尺抓不到它们＝**把"带框行数"当成了"AC 标签行数"** |
| 票 230 `:10`（现量节） | **该票 7 格全部未勾** | 未勾 2／已勾 3／摘框 2 | **差 5 枚**：3 枚（AC#2/AC#3/AC#7）已于 09-29 17:27 由非实现者腿 `done-check-1`（提交 `927ce509`）翻勾；2 枚（AC#1、AC#5）已被摘框（AC#5 于 `980641a9`、AC#1 于 `6d22dd6c`）。逐修订同尺复量：`f5bbccd8`＝未勾 7／已勾 0 ⇒ **那句只在 09-29 16:37 之前成立，是一枚过期读数**，本腿⛔ 采它当分母 |
| 票 230 `:20`（`AC#1` 裁定句） | 3 格已抽过／**剩下 4 格** | 带未勾框 **2** 格＋摘框 **2** 格 | **差 2 枚**：那句按"7−3=4"算，把 AC#1、AC#5 仍留在"剩下"里；但这两枚的框已从票面摘掉（AC#1 的账归票 230 `:19` 本身、AC#5 判性＝非判据是约束）。⇒ "剩下 4 格"与"2 格带未勾框"**说的是两件事**，混用会多算两枚 |

**本腿裁定的真分母（票 115 票面 @HEAD，量的）**：判据格 **7** 枚＝带框 5（已勾 3＝AC#2/AC#3/AC#7；未勾 2＝AC#4/AC#6）＋摘框 2（AC#1 账在票 230 `:19`；AC#5 属"约束非判据"）。
**已被非实现者抽过的格**＝AC#2／AC#3／AC#7（`done-check-1`，凭据另见 §2 第 5 项与本腿 §5）；AC#1 的"独立复算"这一半**本腿今天已按其字面判据抽过**（§4 必答 a，本腿自取远程日志、未转述任何人）；AC#4／AC#6 **至今只有实现方自述**（票 115 当年实现腿 `agent-ticket115`/`115b`/`115c` 的 Progress log 与 commit 正文＝裁决者≠实现者 ⇒ ⛔ 当凭据）。
★本腿**不**据"票 115 已 `-done`"给任何一格记通过，也**不**据 `:20` 那句去翻其余格的勾。

## 2. 七项检查逐格（每格＝判据原文逐字／盘上凭据〔尺＋射程目录＋量的还是算的〕／结论）

### 一、**stub 扫描**：`todo!()`/`unimplemented!()`/返回 nil 假成功/mock 替代真实实现
- **盘上凭据**：射程＝票 115 落地面 5 枚件（`git show --name-only` 逐枚取自 `391878a8`／`527d3036`／`c6dbbf94` 三枚本票自己的 commit＝`winsec_windows.go`＋`notice_attribution_115_windows_test.go`＋三枚既有通知用例）。尺＝`git show HEAD:<file> | grep -cE`（量的，打 HEAD blob，⛔ 工作树）：`todo!()|unimplemented!()`＝**0**；`\bmock\b|\bstub\b|\bfake\b`＝**0**；`t.Skip`＝**0**（5 枚逐枚）。
  新用例"返回 nil 假成功"这一维反着钉住了：`notice_attribution_115_windows_test.go` 里 `if n.Path == "" { t.Fatalf… }` 在场＝**0 字节通知会被判红**（尺＝grep -cE `n\.Path == ""`＝1，@HEAD）。
- **正控**（同尺打全仓 tracked `internal/**`＋`cmd/**`）：`mock|stub|fake` 词形命中 **173** 枚文件＝尺灵。⚠ `TODO|FIXME` 同尺全仓正控＝**0 命中** ⇒ 该子项在这仓是空集，本腿⛔ 拿"0 命中"当"扫干净"的凭据，只当"票 115 没新增桩"。
- **结论**：**成立**（票 115 落地面无桩、无假成功、无 Skip 顶包）。

### 二、**DEFERRED 双向交叉核对**：代码内 `// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点` 标记 ↔ 登记表逐条对账（反向亦查）
- **盘上凭据**：① 票 115 五枚件 `DEFERRED` 命中＝**0**（尺＝`git show | grep -cE`，量的）⇒ 本票没新增推迟项。② 全仓 `DEFERRED(D-<数字>)` 这个形＝**0**（尺＝`git grep -hoE` 去重，射程 `internal cmd tools scripts`）；现存的推迟标记是 `DEFERRED(<主题>)` 形（`git grep` 48 枚文件命中，名册形如 `DEFERRED(macOS/Linux)`／`DEFERRED(D28-1)`／`DEFERRED(C25-loop-wiring)`），**与 SPEC-12 第 5 节登记表用的 `D-xx` 号对不上**。③ **`docs/DEFERRED.md` 不存在**（`git cat-file -e HEAD:docs/DEFERRED.md` ⇒ rc=1、件不在 HEAD）。
- **结论**：**不成立（这台秤本身不存在），且缺两味**——
  缺的凭据不是票 115 少写了标记，而是**判据要的那本登记表锚点件 `docs/DEFERRED.md` 从未产出**（AGENTS.md 第 4 节早已预警这批交付物在仓里不存在，本腿现复量属实），加上**代码侧推迟标记形与表侧 D 号不同构** ⇒ 双向核对无法执行。
  **该由谁补**：这味归**人工批准**（要么补 `docs/DEFERRED.md`、要么改 SPEC-12 第 5 节的标记形），不是任何一枚腿能自造的。本腿⛔ 自行映射、⛔ 据"票 115 面上 0 命中"给这一格记通过。

### 三、**SLO 测量完整性**：是否真的用了进程树口径；有无放宽阈值或跳过采样
- **盘上凭据**：尺＝`git rev-list --count 391878a8^..c6dbbf94 -- <path>`（算的）：`internal/observe/thresholds.go`＝**0**、`internal/agent/testdata/golden`＝**0**、`internal/llm/testdata/golden`＝**0**、`scripts/winsec-tests.sh`＝**0**；`internal/observe/thresholds.go` 的 blob 在 `391878a8^` ⇄ HEAD **同串**（`e2677b11…`，尺＝`git rev-parse`）。三枚冻结件同尺各＝**0**（具名：`internal/panel/tokens_fourway_test.go` 在 `391878a^` **尚不存在**，其"未被票 115 碰过"由 `git rev-list`＝0 与出生时刻两样共同给出）。
- **结论**：**成立（票 115 射程内＝零动阈值/零动 golden/零跳采样）**；⚠ 具名边界：SPEC-10 第 3 节那套进程树 RSS 口径**票 115 一条判据都没引**（票面射程＝`winsec_windows.go` 的通知归属面），所以"是否真用了进程树口径"在**本票**是空射程——这一句只对"票 115 没把 SLO 判据蒙过去"负责，不对"整个 S1 的 SLO 口径可信"负责。

### 四、**契约偏离检查**：D34 工具表 RiskLevel、C19 R1–R9、D43 转移表逐条比对
- **盘上凭据**：尺＝`git rev-list --count 391878a8^..c6dbbf94 -- <path>`：`internal/agent`＝**0**、`internal/config`＝**0**、`internal/tools`＝**0**、`docs/PLAN.md`＝**0**；尺＝三枚本票 commit 的合并 diff 文本 `grep -ciE "RiskLevel|C19|R1–R9|转移表|StateAnswer|threshold|golden|SLO"`＝**1**，逐枚现读＝命中在 `noticeNarrowed is the seam tests swap` 那行注释（词形撞 `narrow`，⛔ 契约名）。
- **结论**：**成立**（票 115 落地面没碰任何契约承载件；D22 ban #2"不许起第二个正规化器"那一枚对本票是**承重**判据，它在第 5 项里另量，本腿现读 `winsec_windows.go` 全文 `filepath.(Clean|Abs)`＝**0**，同尺正控＝全仓 `internal/**` 有 **20** 枚文件命中，尺灵）。

### 五、**D22 禁止清单七条新增禁令静态扫描**：裸 `go func(` / 绕过 `PathResolver` 的 `filepath.Clean|Abs` / 明文 Key / wall clock 超时差值 / 镜像取哈希 / 面板批准 L2 / 内部工件写成受门控 Tool
- **本腿边界声明（照实写，⛔ 藏着）**：`scripts/d22scan.sh` 内含 `sh tools/d22scan/runtests.sh`（＝`go test`）与 `go run .`（尺＝`grep -n` 该脚本 `:51`／`:54`）⇒ **本腿跑它就等于跑 `go test`，撞派单第 16 行的硬边界**，故本腿**没跑**，凭据取两样：静态七尺逐条＋CI 自己的仪器步读数。
- **静态凭据**（射程＝票 115 五枚件 @HEAD blob，尺＝`git show | grep -cE`，量的）：七条逐条＝**0／0／0／0／0／0／0**（裸 `go func(`／`filepath.Clean|Abs`／明文 Key 词形／`time.Now().Sub|Add`／镜像取哈希词形／面板侧 L2 词形／artifacts 写成受门控 Tool 词形）；`recover()`＝0（与"无裸 go func"自洽）。
  emoji（ban #8）另尺：`git show | perl -CSD` 打仪器字符类（U+1F000–1FAFF／U+2200–22FF／U+2600–27BF／U+2B00–2BFF／U+FE0F）＝**5 枚件各 0**；正控同尺打 `AGENTS.md`＝**12** 枚命中＝尺灵。
  ⚠ 具名一枚**本腿自己发现的既有非 ASCII**（⛔ 不算票 115 的账、⛔ 不在仪器射程）：`winsec_windows.go:186`／`:724` 各有一枚 **U+2461**（带圈数字②，`od` 逐字节现认），落在 AGENTS.md 写明的"刻意留的空隙"（U+2460–U+24FF）＝`tools/d22scan/main.go:186 emojiRe` 确实不扫这一段。
- **仪器级凭据（本腿 gh api 自取，未转述任何腿）**：run `35616790753`（head `d1056c0`；本腿尺 `git merge-base --is-ancestor` 现验**含** `c6dbbf94`／`527d3036`／`391878a8` 三枚）的 **lint job `106389417800`** 步名册逐字：step4「D22 scanner positive control (seeded red)」＝**success**、step5「D22 seven-ban + emoji scan (tools/d22scan)」＝**success**、step6「gofmt (gofumpt)」＝**success**、step7／8「go vet」＝**success**；该日志判决行逐字（本腿取的原件 59618 字节／618 行，sha256 `916cc89f0573bc07…`）：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=202, bans #1-5 cmd/=21, ban #6 frontend/=40, ban #7 internal/tools/=18, ban #8 design/=16, ban #8 frontend/=40, ban #8 internal/=374, ban #8 cmd/=29`。
  ⇒ 正控制红得起来（step4 success＝seeded violation 真打红）＋台账各 scope 有活口 ⇒ 这枚"clean"不是瞎的。
- **结论**：**成立**（静态七条零命中＋CI 仪器步在含票 115 末笔的 head 上 success 且带正控）。⚠ 残余具名＝**本腿没能自己跑那把仪器**（边界所限），要一发"本腿自跑 d22scan"的读数⇒ 得等两枚写腿退出、由编排者解禁 `go run`/`go test` 另派。

### 六、**场景↔切片双向核对**：每场景每能力能在切片 done 判据找到；每判据能追到决策或场景；任一侧落空 = 方案缺口，回报而非自行补实现
- **盘上凭据**：正向（判据→决策）可追：票 115 票面射程段逐字点名 **D22 ban #2** 与 **C26**（`:11`、`:38`），机制锚在 `docs/specs/SPEC-06-security-gatekeeping.md:52`「→ 展开 8.3 短名 → 规范化 UNC」（尺＝`grep -nE "8\.3|短名|pathForms"`，量的）。
  反向（场景→本票）**落空**：`docs/specs/SPEC-00-product-overview.md` 第 4 节与 `docs/PLAN.md:3118-3123` 那张 D44(b) 矩阵里，`winsec`／通知／8.3／115 命中＝**0**（尺＝`grep -nE … | wc -l`，射程＝SPEC-00 第 4 节正文＋矩阵行）；承载体 `docs/slices/` **不存在**（`ls -d docs/slices` ⇒ 无；SPEC-12 第 6 节列它属"S1 时按此生成骨架"）。
- **结论**：**部分成立，缺哪一味＝切片级承载体不在盘上**：票 115 是一枚**环境相关缺陷＋旧债**票（票面 `Type:` 逐字），它的判据追得到 D22/C26/SPEC-06，但**追不到任何一条用户场景**，而 SPEC-10 第 6 项的字面对象是"每场景每能力 ↔ 切片 done 判据"——**这张映射表（`docs/slices/`＋`docs/STATE_MACHINE.md`／`docs/TOOLS.md` 同批件）从未产出**。⇒ 该由谁补：归**编排者按 SPEC-12 第 4.3 节派一枚切片级（S1）核对腿**，本票这格按字面只能记"反向落空已回报、⛔ 自行补实现"（判据原文自己就是这么写的：任一侧落空＝回报而非补）。

### 七、**失败预演报告**（第 7 节）与**本切片缺口审计**（每切片完成必须重跑，审计者 ≠ 实现者——四轮打磨仍未穷尽是方案自己写下的结论）
- **盘上凭据**：尺＝`find docs -iname '*premortem*' -o -iname '*pre-mortem*'`＝**0**；`grep -rl "失败预演" docs/reports docs/evidence` ⇒ 只有 **0 枚**在 `docs/reports`（`wc -l`＝0），命中的 8 枚全在 `docs/evidence/**` 且**逐枚现读都属别的票**（136／156／158／161／219／257／265／s2/14）——**票 115 名下 0 枚**（正控＝同尺打出那些枚枚非 0，尺灵）。
  本切片缺口审计这一半：票 115 名下裁决表本腿起手现量 `find docs/evidence -name '115*'`＝**0**（⇒ 就是本件，出表前零枚）。
- **结论**：**不成立（缺整件），缺哪一味具名**——① **失败预演报告**（SPEC-10 第 7 节那八项＋D42 十二项）在盘上零枚，且它的字面对象是"每切片完成后跑"，票 115 不是切片 ⇒ **该由编排者派一枚 S1 级预演腿**，⛔ 由本票代付、⛔ 本腿现造（造出来就是假凭据）。② "审计者 ≠ 实现者"这一半本腿**已就位**＝本件就是头一枚非实现者署名、落在票 115 名下的缺口审计。⇒ 这一格记**不成立**，不记"票已 `-done` 所以做过"。

## 3. 未勾的两格（AC#4／AC#6）缺哪一味＋该由谁补

- **AC#4（变异）**：判据要"把修法退回原状 ⇒ AC#3 新用例必须红；再一发只 `EqualFold` 的半修 ⇒ 也要红"。
  盘上有的：两枚用例在 `notice_attribution_115_windows_test.go:191`／`:240`（@HEAD，尺＝`grep -nE "^func "`）；归属函数在 `winsec_windows.go:106`／`:120`。
  **缺的那一味＝非实现者自己打这两发并读红名的一次记录**。盘上现成的两发都**不是**打在这两枚用例上的等价物：`docs/evidence/s1/126-adversarial-acceptance.md:99-104`（验收方自打，具名"这枚读数不是交件提供的，是我打的"）——`MUT-SAMETREE-TRUE` 把 `sameTree` 掰成恒真红 **9** 枚、红名册里**含 `TestNoticeAttributionKeepsTwoTreesApart`**；`docs/evidence/s1/129-adversarial-acceptance.md:111-123` 的 `MUT-ASKED-REFUSES` 三态读数也真打了归属面（改前 rc=0＝假绿量成读数／改后 rc=1／还原复绿）。
  ⇒ 这两样**支持**"这组钉子双向有牙"，但 AC#4 字面要的是那**两枚具体形状**（退回原状＝按拼写 `==`；半修＝只 `EqualFold`）在**今天 HEAD** 上的复跑读数——`126`／`129` 那两发放的不是这两形。⛔ 为凑齐而圆判。
  **该由谁补**：一枚能跑 `go test` 的突变腿（等 `352-r1`／`230-r1`／`351-r2` 退出、编排者给编译面窗口），载具＝仓外 `git archive` 快照，两发落点必须不同、逐发还原并报还原哈希。
- **AC#6（门禁）**：判据要"按包 `-count=2 -v` 四数逐条点名＋`gofmt -l`＋`gofumpt -l`＋`go vet`＋`sh scripts/d22scan.sh` 纯净快照 rc=0、台账各 scope 不降"。
  本腿**取到**的（⛔ 编译）：① `gofumpt -l` 逐枚＝**0 标脏**，载具＝仓外导出树（`git archive c16ff47e` 的 HEAD blob，尺＝`/d/work/base/gopath/bin/gofumpt.exe -l` 跑前先 `unset GOTMPDIR TMPDIR`；正控＝同尺打一枚已知脏形命中 1、打干净形 0）；② CI 仪器步读数（见第 5 项：d22scan 正控＋真扫＋gofumpt＋两枚 `go vet` 全 success）。
  本腿**取不到**的：`go test -count=2 -v` 的包内四数、以及票 115 当年那批 `RUN/PASS/FAIL/SKIP` 在今天 HEAD 上的复量——它们只在票面 Progress log（署名 `agent-ticket115`／`115b`／`115c`＝**实现方自述**，裁决者≠实现者 ⇒ ⛔ 当凭据）。
  **该由谁补**：同 AC#4，一枚有编译面窗口的门禁腿（`internal/winsec` 与 `cmd/wisp` 此刻都在别人写面上）。

## 4. 必答（派单点名的独立复算，本腿 gh api 自取，⛔ 转述任何人）

**(a) AC#1 那一句"独立复算上面四条红＋复算上一枚 run 里它是 PASS"**——本腿自己取数坐实了后半句、并逐字取出四枚红点：

| 项 | 本腿现量（尺＝`gh api repos/CarlosShao/wisp/actions/jobs/<id>/logs` ＋ `grep`） | 票面自陈 ⇄ 本腿对拉 |
|---|---|---|
| 上一枚 run `35595651898`／job `106319703680` | 日志第 **206** 行逐字 `--- PASS: TestSealReportsThePrincipalsItCleared (0.03s)`；第 **204** 行＝`ERROR winsec: refusing to install a path resolver …`；该轮 winsec 四数（第 1359 行）＝`RUN=71 PASS=32 FAIL=3 SKIP=0` | 票面 `:112-122` 写"第 206 行 PASS／第 204 行 ERROR"⇒ **逐字对上**；字节 214419／行 1378 ⇄ 票面 `:111` 同数 ✓ |
| 本轮 run `35599458439`／job `106331840177` | 第 **204** 行＝`INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2`；顶层 FAIL 逐名 **4** 枚＝`TestAC1SealFileReportsTheInheritedGrantItCleared`、`TestAC2InheritedNoticeHasANoiseBound`（连同 leg 3 子项）、`TestSealReportsThePrincipalsItCleared`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`；四数（第 1458 行）＝`RUN=80 PASS=36 FAIL=4 SKIP=0` | 字节 234680／行 1477 ⇄ 票面 `:126` 同数 ✓；红点 1／2 逐字（`:219`／`:273-274`）与票面 `:130-132` 同串 |
| 正向形状四判据 | `HTTP/1.1 200 OK`（两枚 job 各一发，stderr 0 字节）；首行真时间戳（`Current runner version: '2.337.0'`）；`##[group]Run bash scripts/winsec-tests.sh` 两枚都在**第 187 行**；字节数⇄行数与票面逐字相同 | ⇒ 读到的东西算读到 |
| 转绿链（本腿另取两发） | job `106361546809`（run `35608530583`，head `a13c6dc`）四数 `82/40/2/0`＋两枚 `TestNoticeAttribution*` 逐字 `--- FAIL`；job `106389418059`（run `35616790753`，head `d1056c0`）四数 `82/42/0/0`＋同两枚逐字 `--- PASS`；三句 Fatal（`the instrument planted nothing`／`not attributable through spelling`／`the instrument measured nothing`）在两枚里各 **0** 次 | 本腿⛔ 采表：原件落仓外 `D:\tmp\wisp\20261020\115v1\remote\`，sha256 见 `logs/04`、`logs/07` |
| ⚠ 具名一枚**读数更正**（本腿现量 ⇄ 票面漂移） | 票面 `:110-111` 说 run `35595651898` 的第 187 行是 `##[group]`、第 **206** 行是 PASS——都对上；但票面 `:141`／`:253` 那枚"本机 `RUN=82 PASS=42`"**不是**这两发里的任何一发（本腿取的两发＝`71/32/3/0` 与 `80/36/4/0`） | ⇒ 那句是**本机**读数，⛔ 与远程混用；`82/42/0/0` 只出现在结案 run `35616790753` |

**(b) R-115-1（`noticeNamesTree`/`noticesAboutTree` 无生产读取者）今天还在不在**：**在**。
尺＝`git grep -nE 'noticeNamesTree|noticesAboutTree' HEAD -- 'internal/**' 'cmd/**' 'tools/**'` 剔 `_test.go` ⇒ **7** 枚命中逐枚现读全是**定义／注释/包内互引**（`winsec_windows.go:106`／`:120` 定义、`:123` 是 `noticesAboutTree` 体内调 `noticeNamesTree`、`:74`／`:114`／`:316` 与 `resolve.go:428` 是注释）；测试侧调用＝**30** 枚命中散在 **7** 枚 `*_test.go`（含票 118／126／129 的后继用例）。
⇒ 这一格的账**已由两枚非实现者裁决表各自独立复算并记为"仍未接上"**：`docs/evidence/s1/121-preflight-deps-reachability.md:366`／`:372`／`:399` 与 `126-adversarial-acceptance.md:76-79`、`129-adversarial-acceptance.md:101-102`。票 117（`security-warnings-have-no-listener-in-production`）与票 121 都带 `-done` 而**能力面仍未接线** ⇒ 这笔账现在挂在票 117／121 那条链上，⛔ 由票 115 结案替它说话。

**(c) 本表与票 230 `AC#1` 的关系**：票 230 `AC#1` 原句"那 7 格必须由非实现者出一张 `docs/evidence/s1/115-*.md`（七项检查按 SPEC-10 第 8 节），表出来之后才谈补勾"——**本件交付那枚表**，故 `AC#1` 里"表仍未交付"这一半（票 230 `:20` 末段现量 `find docs/evidence -name '115*'`＝0）到此**有物可销**；
但"才谈补勾"这半句归**编排者**，且本腿按 §1.1 的裁定：`AC#2/AC#3/AC#7` 那三枚已翻的勾**不撤**、`AC#4/AC#6` 两枚未勾的框**本腿一枚不翻**（凭据缺的就是 §3 那两味）。

## 5. 本表⛔ 采为凭据的东西（写明，免得后人误引）

- 票 115 Progress log 里 `agent-ticket115`／`115b`／`115c` 的全部自述读数（含 `82/42/0/0`、五发变异表、`R-115-2` 的登记）＝**实现方自述**，本腿只当线索与"该复算哪一样"的清单。
- `git show -s --format=%B f5bbccd8` 那段结案语＝编排者的，**不是**裁决表的判语（票 115 AC#7 那格的 09-29 附注已把这枚 id 落过两处独立盘上件，本腿 §4(a) 又自取了第三发）。
- "票 115 带 `-done`"＝收口动作的产物，⛔ 零残余证明（`issues/README` 口径，票 230 `:19` 逐字也这么写）。
- 票 230 `:10` 那行"7 格全部未勾"⇄`:20` 那句"剩下 4 格"＝**两枚过期/异口径叙述**，本腿§1.1 已各自具名，⛔ 采任何一枚当分母。

## 6. 两票框终量尺（本腿一枚没翻）

| 票 | 未勾 `^- [ ]` | 已勾 `^- [x]` | 判据格（行首 `- ` 含 `**AC#`） | 摘框（不带框） |
|---|---|---|---|---|
| 115（`:115-*-done.md`） | 2 | 3 | 7 | 2 |
| 230（`:230-…-closed-tickets.md`） | 5 | 0 | 5 | 0 |

尺＝`grep -cE` 锚定行首；本腿对该两枚文件**只追加票尾 bullet**，框行一字未动（票尾 numstat 删除列＝0，见各笔提交正文）。

## 7. 本件自证（尺名随读数）

- 起手 `find docs/evidence -name '115*'`＝**0**（正控同尺打 `104*`＝1）⇄ 交件侧＝**1**（就是本件）。
- 逐枚尺底片：`.scratch/wisp/probes/115/v1/logs/`＝`00-start-rulers`／`01-denominator`／`02-remote-recompute`／`03-redpoints-verbatim`／`04-transgreen-selffetch`／`05-d22scan-credential`／`06-ci-d22scan-on-115-bytes`／`07-lint-d22scan-ledger`／`08-ci-scan-verdict-lines`／`09-closing-rulers`／`11-attribution-and-controls`。
- 远程日志原件**落仓外** `D:\tmp\wisp\20261020\115v1\remote\`（临时件只建不删；⛔ 盘符根）。
- 本腿全程 `go test`／`go build`／`go vet`／`go run` 执行次数＝**0**；`.go` 字节改动＝**0**；push＝**0**。
