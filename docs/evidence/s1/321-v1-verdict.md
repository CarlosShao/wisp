# 票 321 终裁判语（321-v1·★非实现者·只读·裁决者≠实现者）

- 腿＝票 321 终裁腿 321-v1｜时刻＝2026-10-20 21:5x +08（下面每处读数都是本腿自己现跑，⛔ 抄任何一枚腿的数）
- 硬闸＝全程只读：⛔ 动任何 `.go` 字节、⛔ 跑任何 `go test`／`go build`（`353-r1`⇄`139-r1`⇄`109-r1` 此刻正在 `internal/winsec`⇄`internal/agent`⇄`internal/models` 落笔）
- 凭据面＝`git cat-file`＋`git grep HEAD`＋`probes/321/r1/**`＋两批交件件＋重放腿仓外 `logs/` 原件＋异腿既有全包底片
- 尺全册＝`.scratch/wisp/probes/321/v1/00-ruler-register.md`（逐把尺名＋命令形＋射程＋量的还是算的）
- 被裁两批一律当未验证＝旧批 `58de465a`（10-15 20:03 落地）⇄`c673f67b`（10-15 20:53 收件）；新批＝今天同腿重放 `e0c8cbf6`→`fd091cd1`（7 笔）＋`1104fea8`（18:45 派单存档）＋`c7f813c7`（19:44 ★编排者代提遗产四枚件）
- ★本腿一枚框没翻、`Status:`⇄票名⇄`-done` 未动、票尾⛔ 补第二条 bullet（额度已被 `64d14a70`⑤ 占掉）。翻框归编排者凭本件裁。

## 第 1 件｜判据件在位⇄未被篡改 ＝ **成立**

| 尺（R1/R1b/R1c，全量的） | 读数 |
|---|---|
| `git cat-file -e HEAD:cmd/wisp/exec_holder_toggle_321_test.go` | rc=0＝在位 |
| 三点 blob 对回 | `58de465a`＝`851936e396458aa2b6d19cea0da8f71d7d9593f7` ⇄ `HEAD`＝同串 ⇄ `git hash-object`(工作树)＝同串；`c673f67b`/`fd091cd1`/`16e2522a` 各笔亦同串 |
| 两种哈希同源（⛔ 当两回事） | git blob SHA-1 `851936e3…` ⇄ sha256 `717b2d2d9b76d8ae9d413e53a14e11f1669b834fcb056af6b884e01febb63738`＝**同一枚字节的两种算法**；重放腿自报的 `717b2d2d…` 与本腿 `851936e3` 逐字对回 |
| 动过该路径的全部笔（尺＝`git log --all -- <path>` 逐笔列名，⛔ 区间归因） | **只有 1 笔**＝`58de465a`；正控＝同尺对不存在路径回 0（尺没瞎） |
| 件维度 | 8671 字节／200 行／CR 枚数 **0**（两把尺：`tr -cd '\015'|wc -c` 与 `od -An -tx1|tr|grep -c '^0d$'` 都是 0；正控＝已知含 1 枚 CR 的件读出 1） |
| 工作树⇄HEAD ⇄index | `git diff --numstat HEAD -- <件>`＝空、`--cached HEAD`＝空；`cmd/wisp/` 未跟踪件 0 行 ⇒ 今天两批合起来⛔ 动过它一个字节 |

⇒ 旧批自报的"200 行整／8671 字节／无 build tag"与件头射程段、两枚用例名 `TestAC321PanelHopDrivesTheHolderToggle`⇄`TestAC321NamedExecutorIsTheToggleRoad` 全在库。**无篡改、无第二枚同形尺**（`git ls-files` 里 `321` 名下只有这一枚 `cmd/wisp/*.go`）。

## 第 2 件｜牙到底有没有 ＝ **成立**（三条旁证并档＋方法学独立判）

1. **红句逐字在库**（尺＝对 blob `grep -F -c`，量的）：`AC#1 RED: this hop is NOT toggle semantics`＝**2**（:158 接线尺、:198 能力尺）／`executor set`＝1／`drives NO executor on`＝1／`want exactly one toggle`＝1／`which lacks`＝1。腿引的四个 `file:line` 154/158/173/198 与本腿逐行对回**全中**。
   ⚠ 具名一枚本腿自己的坏尺：R1c 那发对整串 `drives NO executor on newPanelHolder` 读 0＝**尺射程所致**（源串是 `%s ... drives NO executor on %s`，`newPanelHolder`⇄`[show]` 由占位符运行时插入），R10a 已改锚 `drives NO executor on` 读出 1＝⛔ 据此判腿报假。
2. **落刀前后锚命中⇄还原哈希＝本腿独立复现**（尺 R21，全在**仓外副本**上做物理改写，探针件⛔ 叫 `.go`；⛔ 碰仓内字节、⛔ 跑测试）：

| 发 | 本腿复现 | 腿自报 | 判 |
|---|---|---|---|
| 旧批基线 `58de465a` blob | `816017b68c495b6c…f9d90d`（锚行 269、前 1/0 后 0/1、变动 2 行＝1 加 1 删） | 同串 | **逐字复现** |
| 旧批突变甲 | `248d5baebb9a651f…c90d57e` | 同串（全 64 位） | **逐字复现** |
| 新批 pristine `HEAD` blob | `4f7a10b30d180132…837237`（锚行 **289**＝行号漂 20、命中 1/0） | 同串 | **逐字复现** |
| 新批突变甲⇄突变乙 | `5f32c10272aed98a…` ⇄ `cba71d39cca3dc20…` | `5f32c102…` ⇄ `cba71d39…` | **前缀对回**（腿只给 8 位） |
| 还原链 | 旧批 `.bak@HEAD` sha256＝`816017b6…` 且 `cmp` ⇄ `58de465a` blob＝**SAME**；`sha-baseline.txt`⇄`sha-restored.txt` 同串 | 同 | 成立 |

3. **重放腿仓外原件本腿逐枚重量**（尺 R7：`grep -c '^=== RUN'`⇄`^--- PASS`⇄`^--- FAIL`⇄`^--- SKIP`，四数不混用；原件 41 枚 `.txt` 仍在 `D:\tmp\wisp\20261020\321r1\logs\`）：`b-family`＝23/18/3/2 ⇄ `mA-family`＝23/16/5/2 ⇄ `b-pos1`＝2/2/0/0 ⇄ `m-pos1`＝2/0/2/0 ⇄ `mB-pos1`＝2/0/2/0 ⇄ `m-restored`＝2/2/0/0 ⇄ 装配族成对 5/5/0/0⇄5/5/0/0 ⇄ `b-255-roster`＝1/1/0/0 ⇄ `m-neg-ruler`＝1/1/0/0＝**与腿报逐枚同号**；`.meta` 载具行记明 `TREE/PAT/MODE/GOTMPDIR/PATH_HEAD`，`runner.sh` 里那发是 `go test -count=1 -timeout 30m -run` ＋ POSIX 冒号 `PATH`（sherpa/build 在前）＝公共块定式。**本腿⛔ 复跑**（缺的那一味＝一枚同刻异腿重打，见第 3 件）。
4. **★方法学那条＝本腿独立判（不是抄腿的）＝成立**：`runResidentBody258`（`resident_hotkey_258_test.go:25`）＝`parser.ParseFile(fset, filepath.Join(".","resident_windows.go"), nil, …)`，第三参 `nil` ⇒ **测试运行时自己开盘上的产码**；`productionFiles228`（`resident_ball_228_test.go:64`）＝`os.ReadDir(".")`＋逐件 `parser.ParseFile(…, nil, …)` ⇒ 同一形状。`-overlay` 只换 **go 命令编译期**的输入、⛔ 改测试进程 `open` 到的内容 ⇒ 用它变异被测量那一边，这枚尺**恒绿＝假读数**。⇒ 两批"禁 overlay、只用物理改写＋sha 还原"的定式**正确、必守**，且适用范围＝本仓**所有**读盘的 AST 钉（本腿另现量 `cmd/wisp` 里 26 枚件用 `parser.ParseFile`、16 枚用 `os.ReadDir`，同族尺一并按此办）。
5. 反面尺复认（缺口属实）：`TestBallPanelGesturesReachThePanelThread` 在库 `cmd/wisp/panel_resident_windows_test.go:1009`＝旧批证据件引的号**今天仍对**；其突变甲下仍 `--- PASS rc=0`（原件 `m-neg-ruler.txt`＝1/1/0/0）⇒"改成 show 今天不会有红"成立、本格⛔ 拿它冒充。

## 第 3 件｜成对逐名作差可不可信 ＝ **部分成立**（缺的那一味具名＋归谁）

- 可信的一半（尺 R7c，量的＋算的）：本腿从**原件**逐名重算红册⇄SKIP 册＝基线红 3 枚（`Test258BridgeMutationNoSrcKeepsOldBinding`／`Test258OccupiedCombinationNamesTheNewValue`／`Test258V1ProbeSummonEditRebindsLiveBall`）；突变红 5 枚＝那 3 枚**＋判据两枚**；`comm -13`＝**恰 2 枚且正是判据自己**、`comm -23`＝**空**（没把任何一枚打绿）、SKIP 名册 `diff`＝空（分母没动）。正控＝突变册含 `TestAC321PanelHopDrivesTheHolderToggle`（1）、基线册含它（0）。
- 台面同一性亦成立：两棵导出树今天 138 枚 `cmd/wisp/*.go` `cmp` **全 SAME**（＝还原彻底）、`treeB`⇄`treeM` 的 `resident_windows.go` sha 与 `HEAD` blob 三串相同；`mA-diff.txt` 原件＝`289c289` 单行加删。
- 那 3 枚存量红＋2 枚 SKIP 的归因本腿复认＝机主常驻 `wisp.exe` **pid 21892 此刻仍在场**（尺＝`tasklist`，量的），底片里逐字 `mute=Ctrl+Alt+M occupied by another app`／`panel=Ctrl+Alt+P occupied by another app`／`live slots after the rebind = 1, want 2`＝公共块**弃权名册**那一族 ⇒ ⛔ 记进本票字节。
- **缺的那一味**＝"落地⇄基线"那对（`5/5/0/0`⇄`5/5/0/0`）是**退化对**＝本腿零 Go 字节 ⇒ 两棵树≡同一份 HEAD 字节，它只证"⛔ 把别枚打红"、**⛔ 证牙**；判别力量全在"基线⇄变异"那一对。而变异侧的底片＝**同一枚腿的两趟**（10-15 与今天），异腿只给了**绿侧**（`orch-322v1-ac4-cmdwisp-full.txt` 10-16、`surface-c.log` 10-17、329 两枚 10-17、322 两枚 10-18：`--- PASS: TestAC321…`×2／FAIL 0，四数与名册本腿逐枚复量对上）。
- ⇒ 判：**作差的算术⇄底片可信**；**作"AC#1 有牙"的凭据只到"部分"**，欠的一味＝**一枚异腿（⛔ 321 身份）在可跑色的窗口里重打同一对变异＋逐名作差**。该由谁补＝编排者给窗口（同一时刻 `cmd/wisp` 只许一枚跑色、且要 winsec⇄agent⇄models 三无在飞）。本腿⛔ 圆判。

## 第 4 件｜`AC#2` 门禁面核 ＝ **部分成立**（只能算"曾经绿"，⛔ 今天判据绿）

- 读数核过（尺 R7f/R4，量的）：新批五把 `rc=0` 的底片件在盘且字节数与自报逐枚同号（`gate-build.txt`＝0、`gate-gofumpt.txt`＝0、`gate-gofmt.txt`＝0、`gate-vet-winline.txt`＝0、`gate-d22scan.txt`＝23774、`gate-length.txt`＝2369）；`d22scan` 末行逐字 `clean - no D22 ban violations`；长度门末行 `VERDICT GREEN … roster equals the tree`、`over the hat=60 roster entries=60`；旧批 `gates.txt` 同五把 rc=0。
- `vet` 那把的**正控形状对**：第三棵仓外树带坏 `fmt.Printf("%d","not-an-int")` ⇒ `rc=1`＋逐字诊断 `cmd\wisp\zz321r1vetprobe_windows.go:11:38 …`（原件 106 字节在盘）＋另一把坏包名 `rc=1` ⇒ 那枚 `rc=0` 是**验过的读数**。
- 射程漂移本腿解释了、⛔ 当假读数：两批 `d22scan` 的 `cmd/` 136⇄143＝件真在长（本腿现量 tracked `cmd/**.go`＝144、其中 `cmd/wisp`＝139、`internal/**.go`＝607）。
- ★分档（本腿独立量的窗口事实，尺 R8/R16/R18e/R24）：腿跑门时自陈 `HEAD=4e0ac985`（19:13，祖先尺 YES）；判据件⇄被测量件⇄两枚提取器宿主在 `4e0ac985`⇄今天 HEAD **blob 全同串**，`git diff --name-only 4e0ac985 HEAD -- cmd/wisp/`＝**0 枚** ⇒ **判据自己那面今天字节未动**。但 `go build ./...`／`d22scan`／`vet -tags winlive ./cmd/wisp/` 的射程含整包编译闭包：`cmd/wisp` 有 1 枚件直 import `internal/winsec`（`secret_dataroot_119b_test.go`），并经 `internal/agent/spill.go`（prod→winsec）传递吃进去；而 `internal/winsec/**` 在 19:13 之后被 **4 笔**动过（`3a509e50` 20:19、`fa395b33` 20:36、`110ed00b` 20:41、`ffc3c62a` 21:17）＋此刻仍有在飞字节（`winsec.go` 20:31:37、`pinnedleaf_342_{fallback,unix,windows}.go` 20:28–20:31，★本腿只 stat、⛔ 读内容）⇒ 编排者 20:0x 那发 build/d22scan/vet rc=0 对本裁**只算"曾经绿"**。
- **缺的那一味**＝一枚在**干净窗口**（`internal/winsec`⇄`internal/agent`⇄`internal/models`⇄`cmd/wisp` 四面无在飞、且门禁四数只在共享树跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）对**今天 HEAD** 重跑那五把并报 `rc=`。该由谁补＝编排者给窗口（⛔ 归本腿、本腿硬闸⛔ 跑）。⇒ `AC#2` 本格判**部分成立**：字面要求的五把尺两批都交过且带正控，缺"今天这一发"的时效凭据。

## 第 5 件｜重放这件事本身该不该退件 ＝ **⛔ 退件；错在编排者的放腿尺，⛔ 在实现腿**

- 放腿那一刻那把锚定尺**枚数为真、集合为空**（本腿独立量，尺 R16b/R23g）：`git log 1104fea8 --grep='^票 321 落地腿 321-r1'`＝**1**＝派单存档自己；严格形 `git log 1104fea8^ --format='%s' | grep -c '^票 321 落地腿 321-r1'`＝**0**；18:45 前 60 分钟动过 `cmd/wisp` 的笔＝**0**；`1104fea8` 正文逐字"防双派四尺全过＝前缀 0 笔、cmd/wisp 60 分被动 0、工作树脏 0…"。根因坐实＝旧批落地笔主题起头逐字 `321-r1 落地腿交件（票 321，零产码）…`＝"票 321" 在**括号内** ⇒ 行首锚定**按构造**看不见它（本腿另量：非锚定 `票 321`＝12 枚、腿码任意位 `321-r1`＝14 枚，两枚旧笔在内＝`58de465a` 两把都中、`c673f67b` 只中腿码那把）。
- 重放腿的处置**正确**：起手 `git cat-file -e` 撞出已落地 ⇒ ⛔ 钉第二枚同形尺（那会造出两枚守同一格、把"哪一枚承重"变糊）、改道交复验＋把盲区具名上报（`e0c8cbf6`／`ed7d183b`／`a09674fc`／`64d14a70`）。⇒ ⛔ 退件；要退的是**编排者那封 18:45 的派单**（放腿前⛔ 跑判据件维度那把尺）。
- ★具名两处**编排者自报⇄盘上不符**（本腿⛔ 替它圆）：
  1. 我派单第 16 行引"台账 `A1068`②"——**盘上无此节**：`grep -c '^## A1068' docs/reports/pending-and-issues.md`＝**0**，全仓 `A1068` 只命中 2 枚件（＝我这封派单＋`predispatch-rulers.sh` 的注释），而 `A1067`⇄`A1069`⇄`A1070` 三节俱在（正控：同尺 `A1069`＝5 枚文件）⇒ **该笔账尚未入库**，请编排者补登（本腿无写台账的额度）。
  2. 派单第 8 行"今天同腿重放 **8 笔** `e0c8cbf6`→`fd091cd1`"——那枚区间里是 **7 笔**（笔①…⑦），第 8 枚锚定命中是 `1104fea8` 派单存档＝**口径写歪**（名册见下）。
  3. 本轮新坐实那句"`internal/agent`⇄`internal/models`⇄`cmd/wisp` **都** import `internal/winsec`"——**只对两枚**：`internal/models` 直 import winsec **处数＝0**（尺 R22a：逐件 import 块抽取，全包只 import `buildinfo/not/observe/statemachine`，且那三个下游引 winsec 处数各＝0；包内 7 枚 `internal/winsec` 命中**全是注释/引文**，见 R23 逐行）。`internal/agent`＝1 处（`spill.go`）、`cmd/wisp`＝1 处（`secret_dataroot_119b_test.go`）＋经 agent 传递 ⇒ **禁跑三包的闸⛔ 因此放松**（models 那发的真风险另有其面），但**那句归因本身⛔ 成立**，请编排者改正形。
- **`c7f813c7` 具名剔离处置**＝该笔主题起头逐字 `编排者代提 321-r1 遗产四枚件（…）` ⇒ 锚定 `^票 321 落地腿 321-r1`＝**0**、非锚定 `票 321`＝**0**、腿码任意位 `321-r1`＝**1** ⇒ **⛔ 计入重放腿名册**（编排者故意⛔ 用腿前缀＝保名册诚实，本腿认这个形）；但它替这票落了四枚读数/尺件（`…build-final.sh` 3394B/44 行、`…final-head.txt` 4255B/66 行、`…tree-equal.sh` 2228B/34 行、`…tree-equal-v1-broken.sh` 1352B/22 行，全对入库 blob 量、CR＝0、⛔ 0 字节、⛔ `.out`）⇒ 本判语按**两栏**记账：重放腿＝8 笔（7 交件＋1 派单存档）、编排者＝1 笔代提（内容属腿的读数面、归属⛔ 属腿）。
- **"复验批"能不能当凭据**＝**能当第 1、2 件的凭据**（件在位＋哈希复现＋原件四数重量＋异腿绿侧旁证），**⛔ 能当第 3 件红侧⇄第 4 件"今天绿"的凭据**（须一枚异腿／干净窗口的趟）。
- **要不要具名登"10-15 已落地"**＝**要**，且只能由编排者落（本腿⛔ 动票面⇄台账）：票面 `Status:` 现量仍 `待派`（＝今天第二颗脑袋还会重放）、建议凭本件在票面具名"判据＝`58de465a` 落地／`c673f67b` 收件；今天 8 笔＝复验"，并把 `predispatch-rulers.sh` 的③（判据件维度，脚本自己写着"不传路径＝⛔ 不许据此放腿"）设为**放腿前硬闸**——那把尺 20:35 才入库、18:45 那趟根本还没它＝单靠"建议补尺"防不住下一颗脑袋。

## 本腿自报⇄盘上（⛔ 替谁圆，包括替自己）

- 本腿自抓两把坏尺（已就地更正、旧件保留⛔ 删）：① R3b 第 h 段用 `$(git cat-file blob …)` 取字节数＝命令替换吃掉行尾换行 ⇒ 四枚遗产件各少报 1 字节／1 行，R4 改文件重定向后＝3394/44、4255/66、2228/34、1352/22；② R13 v1 用 `cut -c1-180` 按字节截中文＝非法 UTF-8（`iconv` 报 cannot convert），R13b 改 `LC_ALL=C.UTF-8`＋`awk substr` 后 `ICONV_RC=0`；另 R2 v1／R10b 两发 `git grep` 参数序（`-e`／`-F` 落在 rev 之后）报 fatal，按"报 0 命中前先跑正控"当场识破、未产生任何假读数。
- 腿⇄盘上不符＝上面第 5 件那三枚（台账 `A1068` 缺节、"8 笔"口径、models 直 import winsec）；两批腿各自的读数本腿逐枚重量**无一处不符**（红句行号、四数、六串哈希、字节/行/CR、名册锚定形、框量）。

## 票 321 框终量⇄写面自证

- 框终量（尺＝对票面 HEAD blob `grep -c` 行首，量的）＝**未勾 2／已勾 0**；正控 `^- \[` 全形＝**5**＝2 枚 AC 框＋3 条 Progress bullet（10-15 两条＋10-20 一条）；`Status:` 行仍 `待派`、票名⛔ `-done`；票面 blob 10437 字节／25 行／CR＝0 ⇄ 工作树同字节。⇒ **本腿一枚没翻**（与派单第 5 行的现量同号）。
- 写面＝只两枚：`docs/evidence/s1/321-v1-verdict.md`（本件）＋`.scratch/wisp/probes/321/v1/**`；逐枚 `git show --name-only` 自证 `.go` **0 枚**（件册⇄提交号＝本交件 send_message 逐枚列）。`internal/winsec/**`⇄票 115⇄230 面在飞字节＝本腿⛔ 读其内容、⛔ 卷进任何一笔（`git add` 全程逐枚字面 pathspec）。
