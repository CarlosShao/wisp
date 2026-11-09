# 239 — 球被窗口切成方形（09-30 桌面签收当场改的三枚数）：落地记录已齐，**缺一份出自非实现者的裁决**

- **Status**：**待派，小活**（编排者 09-30 11:1x 立，起手锚点＝本票落地的那枚 commit）。
- **来源**：owner 10:4x 当面判词「**其他我还能忍受，就是这个球，它也不是个球形啊，四周边边角角的，有边框啊**」→ 我改；11:0x 前第二判词「**还是不对，整体是个方形，外面边框，包括那个向外扩散的那个波纹，都被外面那个方框给局限了**」→ 我改到根子上；11:1x 前终判「**哎，这不就对了吗，可以这个算你过关吧**」。⛔ **"过关"只覆盖形状这一件**，不覆盖票 65 第 63 行那句"整套演一遍并明确说对了"（那句他今天没说，我不替他判）。
- ⚠ **本票的实现者＝编排者本人**（我写的产码）。所以 **AC#2 之后每一格都必须由另一枚腿裁**，我不给自己判（`AGENTS.md` §0.3）。

## 现象与根因（现量，全部我这一轮自己跑的尺）

| # | 事实 | 读数 | 尺 |
|---|---|---|---|
| 1 | **可见外沿＝整个窗口矩形，四角透明** | 改前 `Speaking` 差分 `box=(120,120 72x72)`、`px>=8=4803`（＝窗口面积 5184 的 **93%**）；`Warm` `box 70x71`、`4731` | `build/balldebug.exe -diff .scratch/wisp/probes/signoff-30/dark-states -diff-states …` |
| 2 | 拿 owner 那张截图逐行量，形状**不是圆** | 中段 **41 行宽度恒为 72**、只有上下各 15 行收窄 ⇒ 圆不可能这样（圆应逐行连续收窄） | python 读 `c717e68a-…b2.png` 行剖面 |
| 3 | **根因**：窗口只给 orb 两侧各 **8px**，而"在听"那圈波纹要走到 **R+21.75** | `RingMarginPx = 8`（`internal/ball/hit.go:23`）对 `R+5px+16px` 的波纹（`internal/ball/renderer_windows.go:455`，那个 `16` 是裸字面量、**不在 C21 表里**） | 现读 |
| 4 | **为什么过去每一次取证都没照到** | 存档差分**全部是在 `Sleeping` 那一形拍的**，而 `Sleeping` 画的是"缩小的静息球"（34.72px）⇒ 它的光晕半径 26 装得进 36 的半窗 ⇒ **只有满尺寸的那些形会露出方形** | `internal/ball/tokens_test.go:255-258` 那段自述＋我现跑两形对照 |

## 我改了什么（三枚数＋一处保险，⛔ 零阈值变更）

1. `internal/ball/hit.go`：**`RingMarginPx` 8 → 22**（44..72 全档能装下波纹的最小值）。
2. `internal/ball/tokens.go`：**`DockTriggerPx` 16 → 24**。⚠ **这一枚是被第 1 枚顶出去的**：窗口停在 work area 内时球离切点本来就差一个边距（22），触发距离若小于边距，"停手自动泊靠"永远不成立 ⇒ 改完第 1 枚之后 `TestBallLiveEdgeDock`／`TestBallLiveEdgeDockHover` **当场两枚红**，抬到 24 之后两枚复绿。**这条耦合（`DockTriggerPx >= RingMarginPx`）今天只有那两枚 live 用例看得见**，我把它写进了常量注释。
3. `internal/ball/renderer_windows.go`：新增 `haloExtents(R)`，把光晕的**填充半径**与**淡出半径**一起夹进窗口内切圆 ⇒ 光晕自己该淡到零，不再靠矩形切。⚠ 对 `Sleeping` 是**恒等变换**（它的 1.9R=33 本来就 ≤ 35），所以存档那三发读数（2098/2103/2120）的口径没被我动过。
4. `docs/evidence/s1/c21-native-tokens.md:166`／`:168`：表里那两格（`环边距 8px`、`距边 16px`）跟着改成 22／24，并把**为什么**写进理由栏。⚠ 这两行自己写的规则就是「ring/glow 不被裁」，8 从未满足它自己那句判据 ⇒ 我按**修 bug**处理，不是按改设计。**`TestC21GeometryRowsMatchCodeConstants` 就是逼表与码同音的那枚仪器**（我拿 8.0 复跑过＝绿，拿 22 不改表＝红，红句逐字 `RingMarginPx = 22, but this row states no UNCLAIMED number equal to it`）。

## 门禁读数（本票落地时我现跑，逐名）

- `go build ./...` 过；`gofumpt -l` 三枚文件**零输出**；`go vet ./internal/ball/ ./cmd/balldebug/` **空**；`tools/d22scan/d22scan.exe` **clean rc=0**（`ban #8 internal/=462`、`cmd/=63` 是被扫文件数不是违规数）。
- `go test ./internal/ball/ -count=1` ⇒ **只剩 1 枚红**＝`TestC21TableColourRowsMatchTokensCSS`，红句逐字 `read design/assets/tokens.css: … cannot find the path specified`。**这枚红早于本票**（owner 自己把那枚文件从工作目录搬走、版本库里还在 12,942 字节），归界面那一侧，⛔ 不在本票地界。
- `go test ./internal/ball/ -tags winlive -count=1` ⇒ **同样只剩那一枚 CSS 红**；两枚 dock 用例复绿。⚠ 前提＝跑之前 `balldebug.exe` 进程数 **0**（带着球跑 live 用例会因全局热键被占而**多出 8 枚假红**，我 11:1x 踩过一次，读数已作废重跑）。
- `go test ./cmd/wisp/ -count=1` ⇒ `ok 112.312s`。
- ~~改后差分：`Warm box 66x66 / 3359px`＝同半径圆盘面积的 **98%**（改前是方块的 96%）~~ ⇒ **这句作废，我原文照抄在此**：`239-v1` 在 HEAD 现跑同一形是 `66x66 / px>=8 = 5284`，复现不出 3359。归因〔本腿几何推算＋它现验算〕：那两档"改后"表记的是 `edge=72`／分母 `97344`，而 HEAD 窗口 `edge=100` ⇒ 当时挂的是**中间态的 `build/balldebug.exe`**（`.gitignore:14` 的 `build/`，不随提交走、不可再生）。算术本身没错，错在分子口径。
- ⚠⚠ **形状这一件不靠那句 98% 站住，靠下面 AC#1 里 `239-v1` 复跑的那一发**（`239-v1` 在 HEAD 自建二进制上现跑差分：窗口 `edge=100`、分母 `115600`，三形 box 44/66/80 **全部严格小于窗口** ⇒ 方形没了；`Sleeping` 正面命中 2120 落在 2098/2103/2120 那一族）。`Speaking` 那一发被屏幕上别的窗口刷新洗掉了，**没量到**（具名，不假装）。

## 判据（⛔ 框归编排者，产码／验收腿一枚都不许碰）

- [x] **AC#1 复跑我上面那四把尺**（build／gofumpt／vet／d22scan）＋两形差分（`-diff` 指到 `probes` 下的临时目录），逐名对拉：`internal/ball` 的红名册**必须等于**"只有 CSS 那一枚"；多一枚都算本票没做完。
- [x] **AC#2 攻那枚耦合**：把 `DockTriggerPx` 改回 16（只改这一枚、不动 `RingMarginPx`）⇒ **`TestBallLiveEdgeDock` 必须红**；再把它改成 21（＝小于 margin 22 但大于旧值）⇒ 也必须红。 若 21 那一发是绿的，说明"边距与触发距离的耦合"**根本没有仪器看得见**，就具名登记这一事实（别假装测过了）。
- [x] **AC#3 攻 `haloExtents` 的"对 Sleeping 恒等"这句话**：把 `fill` 与 `grad` 两个返回值在 `R*1.9 <= fit` 时强行互换 ⇒ `Sleeping` 的 `px>=8` 读数必须离开 2098/2103/2120 那一族（离开＝这句话有牙；不离开＝这枚断言是装饰，具名写出来）。
- [x] **AC#4 判"窗口放大"这一支的代价有没有人认**：窗口 72→**100** 见方，透明区随之外扩。要判的三件：① 球**停在 work area 内时 orb 离屏幕边最少 22px**（旧值 8px）——这是可见的行为变化，不是纯修 bug；② 贴边露出的百分比（`DockOverlapFrac=0.42` 那一套）我**一字节没动**，但票 62 AC#5 今天已由 `62-v1` 判**不成立**，所以这一支没顶掉任何"已被接受"的数；③ ~~差分像素的**分母**（`of 97344`）不变（那是截图框，不是窗口）~~ ⇒ **这一支我原文写错了，`239-v1` 判不成立并已改账**：分母是"边长²"，边长＝窗口＋2×`RingMarginPx` 富余 ⇒ HEAD 上实测 **`115600`＝(100+2×120)²**，而不是"不变"。**后果不只这一句**：任何百分比类差分引数**不许跨几何互比**（旧档全部是 `edge=72/97344` 那一族的读数）。⚠ 若裁定"22px 贴边间隙不可接受"，那就得回到"把波纹走距 16px 收进 8px 边距"那一支，**并具名写清那会让波纹几乎看不见**。
  > **2026-10-22 11:0x 编排者翻框＝成立（★★机主当场批「22 可以」）**
  > **机主原话**＝「**22可以**」（回复我摆的那条：票面逐字要求"① 22px 贴边间隙**要 owner 单独裁一句**，我不替他判"）。
  > **裁的是什么**：窗口 **72 → 100** 见方后，球停在 work area 内时 orb 离屏幕边**最少 22px**（旧值 **8px**）——★这是**可见的行为变更**、⛔ 纯修 bug。
  > **我上盘核过那个数是真的**：`internal/ball/hit.go:29` 逐字 `const RingMarginPx = 22.0`，`:28` 注释逐字"**22 is the smallest value that fits the waves at every size in 44..72**"；
  > `:57` `return int(orb + 2*RingMarginPx*scale + 0.5)` ⇒ 22 是**波纹能容下的最小值**、⛔ 随手取的余量。
  > **② 那一半**（`DockOverlapFrac=0.42` 一字节没动、票 62 `AC#5` 已判不成立故没顶掉任何"已被接受"的数）＝成立、无需再裁。
  > **③ 那一半**（分母 `of 97344` 不变是原文写错、`239-v1` 已判不成立并改账：分母随几何变、HEAD 实测 **115600＝(100+2×120)²**）＝**已由 `239-v1` 更正、我复核同向**。
  > ★**机主这一句把三个"要裁"里的最后一个补上了** ⇒ 本格按"机主已裁"翻。
  > ⚠ **本票仍**不结案、不加 `-done`**：残余两枚还没修完（①0.25px 富余——满档 72 时 Listening 波纹外沿到「R+21.75」、窗口半只有「R+22」，**余量 0.25 像素**，且那枚走距「16」至今是**裸字面量、不在 C21 名册里**；②走距收进 C21 几何名册并让 `239-v1` 差分脚本变常驻判据）。归口本票。
- [x] **AC#5 归口**：本票动过 `internal/ball/hit.go`／`tokens.go`／`renderer_windows.go` 三枚文件 ⇒ 与票 65／68 同地界。⚠ **票 68 名下那份"出自非实现者的验收表"今天仍是零份**，本票的读数要一并算进它那格的欠账，不许另立一份平行真相。

## 禁区

零阈值变更（`thresholds.go`／golden／`docs/SLO.md` 一字节没动）；⛔ 不许为了变绿放宽任何断言（本票两枚红是我**改出来的**、不是我**改掉的**）；不动 `PLAN.md`／`docs/specs/**`／`allowlist.txt`；不动三枚冻结件；`frontend/**`／`design/**` 零读零写零转述（含 `design/assets/tokens.css` 那枚不在场的文件——**它的红留在原地，不还原、不提交、不删**）；⚠ **与票 222／221／234／235 同撞 `internal/tools`／`cmd/wisp` 的那几枚一律串行**，本票只碰 `internal/ball`，但 live 用例要独占桌面（跑之前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 必须为 0）。

## 翻勾凭据与残余（编排者 09-30 11:5x 落，凭据全出自「239-v1」现跑）

- **AC#1/AC#2/AC#3/AC#5 已勾**：判语依次＝成立（附台件口径偏离）／成立（16 与 21 两发都双红，耦合有仪器看得见）／成立（互换后「px>=8」2120→1694，离族 404 枚，这句话有牙）／成立（附欠账声明）。凭据表：「docs/evidence/s1/239-ball-square-fix-v1.md」126 行／23,354 字节，八节齐、§6/§7 两枚必填节非空。
- **AC#4 不勾**，条件是实的：② 成立；① 22px 贴边间隙**要 owner 单独裁一句**（他只裁过"形状"，这 14px 的行为变更没被裁过，我不替他判）；③ 那一支是我原文写错的，已在上面就地改账（分母随几何变，不是"不变"）。
- **本票不结案、不加「-done」**：残余没修完就不该改名。残余两枚——
  1. **0.25px 富余（239-v1 具名脆弱项）**：满档 72px 时 Listening 波纹外沿走到「R+21.75」，窗口半只有「R+22」，**余量 0.25 像素**；那枚走距「16」至今是裸字面量、不在 C21 名册里。**残缺表现**＝以后任何人把波纹走距往上调（哪怕 +1）或把边距往下调，波纹立刻重新被窗口切边，而全仓没有一枚测试会响。**完成判据**＝把走距收进 C21 几何名册并让「239-v1」那条差分脚本变成常驻判据（满档 72 下 box 严格小于窗口、且余量有下限）。归口：本票。
  2. **「cmd/wisp」全量没在 HEAD 复跑**：本票那发「ok 112.312s」是改前读的〔读台件，未复跑〕。归口＝编排者下一枚整包门禁。

- [2026-11-09 13:3x +08] agent=`239-r1`（票 239 残余①落地腿，写面＝`internal/ball/**` 两枚产码＋一枚新测试件＋C21 表一枚新行＋本条）交回＝**残余①的"走距收进 C21 几何名册"半已落地、"常驻判据（box 严格小于窗口且余量有下限）"以静态几何尺落地；"239-v1 差分脚本变常驻判据"那一半需要独占桌面的真窗像素差分台面，本腿不擅自造（未定义即停），具名留下**。改动四笔：1) `internal/ball/hit.go` 新增导出常量 `WaveTravelPx = 16.0`（注释带 R+21.75 推导与 0.25px 下限来历）；2) `internal/ball/renderer_windows.go` 波纹绘制点 `float32(t)*16*s` 改用该常量（值零变化，行为零变化）；3) `docs/evidence/s1/c21-native-tokens.md` 几何/动效表新增一行（走距 16px、外沿推导、常驻尺具名），Go 常量列只写 `WaveTravelPx` 一枚（避免与既有行撞单射）；4) `internal/ball/tokens_table_test.go` 的 `c21MotionGolden` 补一枚 pin（双向锁自身规矩：表行命名 tokens.go 之外的常量必须钉值）。新常驻尺＝`Test239WaveTrainStaysInsideTheWindowWithAMarginFloor`（`internal/ball/wave_margin_239_test.go`）：断言 RingMarginPx 减波纹外沿（R+5+16+RingStrokePx/2=R+21.75）不小于 0.25px，且 44..160 全档（96 DPI 基准、半径相对量论据覆盖任意 DPI）box 严格小于窗半；环偏移 5 仍是渲染器裸字面量，测试内以具名局部常量镜像并注明收册属射程外。门禁（共树，13:2x）：gofumpt（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）四枚写面零输出、`go vet ./internal/ball/` rc=0、`go build ./internal/ball/ ./cmd/balldebug/` rc=0；测试色与 d22scan 在 commit 后导出树补跑（见下一条）。残余具名：a) 「239-v1 差分脚本变常驻判据」未做——像素级差分要起真窗（balldebug.exe 独占桌面），本腿不造新台面，判据语义已由静态尺覆盖（外沿/窗半为半径相对量，box 小于窗口在全部档位与 DPI 成立）；b) 残余②「cmd/wisp 全量在 HEAD 复跑」＝归编排者整包门禁，本腿未碰。⛔ 零框翻（本票 AC 五格早已全勾）、⛔ Status/票名未动、⛔ 不加 -done（残余 a/b 未闭）、⛔ push。撤销口令：撤 239-r1。

- [2026-11-09 13:4x +08] agent=`239-r1` 补记（导出树门禁读数，接上一条"见下一条"）：导出树＝`git clone --shared` 于 `b3aefc34b`（仓外 `/tmp/w239r1-export`）。**测试色**：`go test ./internal/ball/ -count=1`（环境按派单：`unset GOTMPDIR TMPDIR TEMP TMP`、PATH 加 build＋sherpa、⛔ CGO_ENABLED=0 未设）⇒ 红 3 枚＝`TestOrbAssetShippedUnmodified`／`TestOrbServingConstants`／`TestJarvisAssetShippedAsBuilt`，红句全为"embedded orb page still contains CRLF"与 sha 不符族；**对照控制坐实为导出树 autocrlf 检出伪红、非本腿缺陷**：主树同选器（含新尺＋C21 双向锁）全绿 rc=0，blob 尺＝`internal/ball/assets/jarvis-orb.html`／`orb-siri.html` 的 `HEAD:` blob CR=**0** 而导出树工作副本 CR=4346/340（`core.autocrlf=true` 两树同值；控制脚本 `.scratch/wisp/probes/239/r1-crlf-control.sh`）；与 174-v3／216-r1 记录过的"导出树检出形红"同族。**新尺与双向锁在导出树独立 PASS**（`-run 'Test239WaveTrain|TestC21GeometryRows'` -v 两枚 PASS）。**d22scan**：导出树 `sh scripts/d22scan.sh` rc=0 clean（ban #8 internal/=967、cmd/=242、frontend/=121、design/=14，全非零）。**去牙突变**：导出树内 `WaveTravelPx` 16→17 ⇒ 新尺红（红句逐字"leaves only -0.750px beyond the wave train outer edge 22.75"）＝票面点名的"走距往上调哪怕 +1"那一形有牙；还原后复绿、导出树 porcelain 空。共树门禁补记：gofumpt 四枚写面零输出／vet rc=0／build rc=0（13:2x 已跑）。⛔ push、⛔ 框、⛔ -done（残余 a「239-v1 差分脚本常驻化」、b「cmd/wisp 整包 HEAD 复跑」未闭）。
