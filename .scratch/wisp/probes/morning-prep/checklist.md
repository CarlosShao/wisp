# 明早 10:30 桌面签收——我自己看的操作清单（morning-prep-1 盘点腿）

- 取锚时刻：`date` 现读 **2026-09-29 22:37 +0800**；锚点 `git log -1` ＝ **`795ed767`**（`dev`，09-29 22:30）
- ⚠ **写这份清单期间锚点自己往前走了**（22:42 现读＝**`1d2ad737`** `test(235 AC#1): 守池钉子注释的理由换成今天盘上成立的机制…`）＝**`235-r1` 正在陆续 commit**，与本腿无关。**下表的三枚源头在两个锚点上重读一模一样**（`cmd/balldebug 3539d473 09-21 08:52`／`cmd/wisp c44b30c4 09-29 19:41`／`internal/ball 22015301 09-21 11:00`），工作树仍洁净、CI 仍无在跑 ⇒ **本清单的六组结论不受这次前进影响，明早直接照用**。
- 本腿身份：**只读盘点**。零编译、零启动界面、零 commit／push、零碰 `frontend/**`／`design/**`。
- 对照的表：`docs/reports/desktop-signoff-2026-09-30.md`（10,733 字节，台账 `A452` 落地）
- 一句话结论：**两枚可执行文件都过期了 ⇒ 明早第一件事就是重建 `build/balldebug.exe`；而 `-tour` 演不出第 3 件要的五种模样，这条会当场冷场，必须提前备好第二条跑法。**

## 1. 要跑的那个程序在不在、新不新

**判据（派单给的尺，逐字照用）**：源头任一枚比可执行文件新 ⇒ 必须重建。

现读（`stat` 的 Modify 行 ＋ `git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M' -- <路径>`）：

| 可执行文件 | 在不在 | 大小 | mtime（Modify） | 它源头的最近提交 | 谁更新 |
|---|---|---|---|---|---|
| `build/balldebug.exe` | **在** | 7,069,696 B | **2026-09-20 23:03:03 +0800** | `cmd/balldebug` ＝ `3539d473` **09-21 08:52**；`internal/ball` ＝ `22015301` **09-21 11:00** | **两枚源头都比它新** |
| `build/wisp.exe` | **在** | 27,881,312 B | **2026-09-21 15:06:32 +0800** | `cmd/wisp` ＝ `c44b30c4` **09-29 19:41**（"221 甲形＋乙形同一发：注册 task.cancel…"） | **源头比它新 8 天** |

⇒ **结论：明早两枚都必须重建，不可直接用。**
- **`balldebug.exe` 是明早唯一真要跑的那枚**（对照表 17 件里 1–6、17 全部落在球上，球只在调试工具里有新样子）。重建命令＝票 68 `:76` 逐字：`go build -o build/balldebug.exe ./cmd/balldebug`。
- `wisp.exe` 过期但**明早不需要它**——除非他要看"那扇小窗"（见 §2 的坑 #3，那一件 balldebug 给不出来）。真要开面板才需要重建 `cmd/wisp`，而那会把整包依赖拉进来，**不在"只用眼睛看"的射程里**。

三条补充读数（都是为了"别拿错件"）：
- `build/` 里还躺着 **`wisp77.exe`（09-21 13:44）与 `wisp77_partial.exe`（09-21 13:44）**两枚同族旧件 ⇒ 明早**只认 `build/balldebug.exe` 这个全路径**，别 Tab 补全点到别的名字上。
- `git status --porcelain -- cmd/balldebug cmd/wisp internal/ball internal/statemachine` ＝ **空** ⇒ 这四条路径**工作树干净**，重建取到的就是 HEAD 的码，不会把谁的未提交改动烤进那枚 exe。
- 近两日**没有任何提交动过球族**：`git log --since=2026-09-28 -- cmd/balldebug internal/ball docs/evidence/s1/ball-states docs/evidence/s1/62-diff-signoff` ＝ **空**（最后一动是 09-21）。⇒ 源头不会在明早之前再变；**今晚重建与否，明早读到的源码是同一份**。但**exe 本身仍是 09-20/09-21 的**，所以 §1 的"必须重建"不变。

## 2. 跑法清单逐条落定（票 68 `:73-82` 那 8 步 ＋ 明早真要用眼睛看的几条）

### 2a. 票 68 `:73-82` 的 8 步，逐条判"是不是明早项"

出处逐字来自 `.scratch/wisp/issues/68-ball-default-visuals-parity.md` 的行号见括号。

| 步 | 命令（逐字抄） | 明早项？ | 为什么 |
|---|---|---|---|
| 1 | `export PATH="$PWD/third_party/sherpa-onnx:$PATH"`（`:73`） | 前置，**照跑** | 必须在仓库根跑（`$PWD` 拼路径）。本腿现读：`third_party/sherpa-onnx/` **在**（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）。⚠ 但 `grep -rn "sherpa\|onnxruntime" cmd/balldebug/*.go internal/ball/*.go` ＝ **零命中** ⇒ 球这枚不链它，这步对明早是**保险带不是必需**。 |
| 2 | 洁净前置：`tasklist \| grep -iE "wisp.exe\|balldebug.exe"` **必须无输出** ＋ 显示器缩放 **100%**（`:74-75`） | **是，他到场前跑** | 见 §2b 每行的"洁净前置"列与 §6 的那把尺。 |
| 3 | `go build -o build/balldebug.exe ./cmd/balldebug`（`:76`） | **非明早项**（编译） | ⚠ 但它**必须做**——§1 已判 exe 过期。**做成他到场之前**，不在他眼前跑（他会看到命令行）。 |
| 4 | `build/balldebug.exe -diff build/t68-flip -diff-states Sleeping,Listening,Thinking,Acting,Speaking,Warm,Settling -diff-sample 5s -diff-dock right`（`:77`） | **半明早项** | 这是**取证跑不是观看跑**：它每态待 5–6 秒、自己截屏、自己写表。屏幕上确实有球在换态，可以顺带让他看"态在换"；但**判据是数字**（门：`Sleeping` px≥8 ≥ 2098、框落 40–46 家族、`timers=no`、`cpu_all ≤0.5%`、teardown `clean`）。⇒ 明早**只在他愿意等的时候跑**。⚠ 输出目录别照抄 `build/t68-flip`，见 §3 的落点规则。 |
| 5 | 同上一条加 `-frozen` ⇒ 期望 **px≥8 = 0**、框 `0x0`（`:79`） | **是，值得演** | 这就是对照表第 13 件说的"反例很适合演给你看一次"：**切回旧规格，球在白壁纸上"有它和没它"零个像素变化＝物理上根本看不见**。票 68 `:79` 原文还写着"**不是可选项**"（逃生门若坏了没人知道）。 |
| 6 | `build/balldebug.exe -state Sleeping -hold`（Ctrl+C；打印 `timers=false` 与 `closed handles=` ≤600）（`:80`） | **是** | 对照表第 2 件"不动的那颗球"就是这个：一态钉住直到他看完。**要 Ctrl+C 或托盘 Exit 才停**。 |
| 7 | `go test -tags winlive -count=2 -v ./internal/ball/`（`:81`） | **非明早项** | 是测试。那天上午不跑整包（派单 §0／台账 `A452`、`A453` 三令）。 |
| 8 | 翻转后重跑 AC#4 四道门（`:82`） | **非明早项** | `gofmt`／`go vet`／两发 `go test`，全是编译与测试。 |

### 2b. 明早真要用眼睛看的六条（对照表 1–6、17 件各自对哪条跑法）

命令里的 `build\balldebug.exe` 是 Windows 反斜杠写法（他这台是 Windows／Git Bash 两栖；票 62 `:82` 原文用的就是 `build\balldebug.exe -stay`）。

| 对照表第几件 | 跑法（命令） | 屏幕上会产生什么 | 洁净前置 |
|---|---|---|---|
| **第 1 件**（深浅壁纸都清晰，票 62 `:82`） | `build\balldebug.exe -stay` | 一起来就是一颗静止的玻璃球挂在桌面，热键／托盘／点击全活。他先看白壁纸这一面；**换完深色壁纸后这一条要重跑一遍**（球会随桌布重成像） | 步 2 那条 `tasklist` 无输出；缩放 100%；`-stay` 会占住终端，**看完 Ctrl+C 再换壁纸** |
| **第 2 件**（睡着的球完全不动，票 62 `:83`） | `build\balldebug.exe -state Sleeping -hold`（＝步 6） | 一张静止画面，无抖动、无呼吸闪 | 同上 |
| **第 3 件**（五种干活模样一眼可辨，票 62 `:84`） | ⚠ **`-tour` 演不全这一件**（见 2c 坑 #1）。要五态齐：**默认全周期**（不加模式参数）`build\balldebug.exe -cycle-ms 6000`，它按 `allStates` 顺序把 20 态各演 6 秒，`Listening/Thinking/Acting/Warm/Speaking` 全在里面 | 球依次换 20 副面孔，终端同步打 `balldebug: state=<名> handles=<n>` | 同上；**这条会一直跑到周期结束**，中途可 Ctrl+C |
| **第 4 件**（有声音进来液面随音量转，票 62 `:85`） | `build\balldebug.exe -state Speaking -hold -level 0.9` 看完 → Ctrl+C → `build\balldebug.exe -state Speaking -hold -level 0.25` | `-level` 是"合成的音量包络，约 30fps 推给球"（`main.go:106` 帮助原文），**喂进去的是假语音电平**，正合票 62 `:85` "喂已知电平"。两发一比就是"大音量 vs 小音量看起来不一样" | 同上；`-level 0` ＝ 不喂（默认） |
| **第 5 件**（不说话但还在听＝边框"滑"进来，票 62 `:85` 后半句） | `-tour` 的第 4 拍，或直接 `build\balldebug.exe -state Listening -hold` | tour 里这拍逐字写着"STOPPED TALKING, STILL LISTENING … the border glides in around the ball. That arriving border is the transition the owner asked for."（`main.go:521-523`） | 同上 |
| **第 6 件**（拖到左／右／上边压扁贴边、划回完整球，票 62 `:86`） | `-tour` 的四边那几拍（左/右/上/下各压扁一次再弹回，全部走真发布路径 `DebugDock`／`DebugPop`）；要看**单侧钉住**：`build\balldebug.exe -dock right -state Sleeping -hold` | 球沿那一侧挤成一条、只露一小片可点区；tour 最后一拍会**把球留在右边缘不动**，让他的鼠标真的划上去看弹回（`main.go:548-556` 原文 "NOW YOUR TURN, FOR REAL"） | ⚠ **`-dock right` 单独打会掉进默认全周期**（`main.go:263-271` 的 dock 块之后直接落 `switch` 的 `default` 分支＝20 态连演）。**必须和 `-state` 或 `-stay` 一起给。** |
| **第 17 件**（20 种模样挨个演，票 07 `:51`） | 同第 3 件的默认全周期 | 20 态各 6 秒，含出错／没网／卡住 | 见 §4 |
| **主事件**（票 65 `:63`：他第二次看 `-tour` 并说"这次对了"） | `build\balldebug.exe -tour`（每拍默认 5 秒，可用 `-tour-dwell` 调慢） | 一条命令走完：静止 →  summoned 液体流动 → 喂声旋转 → 停声边框滑入 → 回静止 → 四边贴边各压扁弹回 → 留在右边缘让他自己动。终端每拍打 `tour NN. <解说>` ＋ 一行 `timers=/dock edge=/progress=/handles=` | 同上；⚠ **`-tour` 自带喂电平（默认 0.70）**，所以 `:78` 那个 `-level` 在这儿是覆盖项，不给也行 |

### 2c. 三条会当场卡住的前置（本腿现读盘上确认过的）

1. ⛔ **`-tour` 里没有 `Thinking`／`Acting`／`Warm` 三态**（本腿读了 `runTour` 全体：原型走法只设 `Sleeping→Listening→Speaking→Listening(边框)→Sleeping→四边贴边`；`-frozen` 走法才含 `Armed/Warm/Error` 六态）。而对照表第 3 件要求**五种干活的模样一眼分辨**、台账 `A440` 也把它算成第 1 件。⇒ **只看 `-tour` 会把第 3 件判缺一半**。备好§2b 那条默认全周期当第二条跑法。
2. ⛔ **明早的球是新样子，而仓里那 20 张存档图是旧样子**（见 §4）——他要是拿存档图对现场，一定"对不上"。这句话要**由我先说出口**，别等他问。
3. ⛔ **`balldebug` 开不了"那扇小窗"**：`cmd/balldebug/main.go:198-199` 与 `:610-611` 逐字写着面板那两枚是 **`stub, ticket 33`**（`OnTrayPanel: func() { fmt.Println("tray: open panel (stub, ticket 33)") }`、`hotkey: panel (stub, ticket 33)`）。⇒ 他要问"点开球之后那个聊天窗"，**球这枚给不出来**，得换 `wisp.exe`（而 §1 已判它也过期）。这一条**提前跟他讲在清单里，不临场解释**。

## 3. 深浅壁纸两拍到底要拍什么（票 62 `:82`，深色那一半从未量过）

### 3a. 现有差分取证件的形状（本腿现读，只读目录与 `.txt`／`.md`，未碰 `design/**`）

四个已存在的取证件目录（`ls -d docs/evidence/s1/62-diff-*`）：

| 目录 | 文件枚数 | 里面有什么 | `diff-table.txt` 里 `Sleeping` 那行的关键读数 |
|---|---|---|---|
| `62-diff-baseline` | 10 | 3 态 × 3 张 ＋ 表 | **px≥8 ＝ 0**、px≥3 ＝ 17、max ＝ 4/255、**无像框**（`bbox 0x0`）＝冻结档那颗球物理上不成像 |
| `62-diff-glass` | 10 | 3 态 × 3 张 ＋ 表 | px≥8 ＝ **2098**、像框 **40×38**、`timers=no`、`teardown=clean` |
| `62-diff-border` | 4 | 1 态 × 3 张 ＋ 表 | px≥8 ＝ **2120**、像框 **44×44**（这就是文档里"44"那个数的出处＝包围盒，不是直径） |
| `62-diff-signoff` | 19 | 6 组 × 3 张 ＋ 表（含 `Sleeping/Listening/Speaking` 各一条 `-dock-right` 行） | 最后一枚 `diff-table.txt` 1,474 B |

**文件名规则（工具自己生成，不用手写；`cmd/balldebug/diff_windows.go:10-12` 与 `:462-463` 逐字）**：
```
NN-<State>-alive.png      带球的整屏合成
NN-<State>-dead.png       同区域、球进程已终止
NN-<State>-diff.png       |alive-dead| 放大图（眼睛要看的就是这张）
diff-table.txt            一张 TSV，表头 27 列（逐字）：
state  alive  dead  diff  x  y  edge  timers  cpu_pct_1core  cpu_pct_allcores  privateWorkingSet_MB
privateCommit_MB  workset_MB  handles  gdiObjects  userObjects  px_delta_ge3  px_delta_ge8  px_delta_ge24
px_total  max_delta  mean_delta  bbox_x  bbox_y  bbox_w  bbox_h  teardown
```
带贴边的那一行 label 里多个后缀，实例逐字：`02-Sleeping-dock-right-{alive,dead,diff}.png`。

### 3b. ⚠ 这 27 列里**没有任何一列记录桌布是深还是浅**

四张表的 `x y edge` 只有两种取值（`1600 500 72` 与 `1720 720 72`），全是坐标；深浅壁纸的区别**只能靠目录名承载**。⇒ 这就是"深色那一半从未量过"的形状原因。

### 3c. 明早"怎么算留下证据了"的最小可核对形状

落点与命名（沿用同一族前缀，日期挂在末尾，**新建目录、绝不覆盖上面那四个**）：

```
docs/evidence/s1/62-diff-light-2026-09-30/     ← 浅色桌布那一拍
docs/evidence/s1/62-diff-dark-2026-09-30/      ← 深色桌布那一拍（第一回有它）
```

每枚目录**最少要有**：
1. `01-Sleeping-alive.png` ＋ `01-Sleeping-dead.png` ＋ `01-Sleeping-diff.png` 三张成对（票 62 `:82` 判的是"静止那颗球在两种桌布上都清晰可见"，`Sleeping` 是最不利的一态，抓它就够）；
2. `diff-table.txt`，表头**必须与 §3a 那 27 列逐字相同**（工具自己写，别手改）。

对应命令（就是把 §2a 步 4 的输出目录指到证据位）：
```
build\balldebug.exe -diff docs/evidence/s1/62-diff-light-2026-09-30 -diff-states Sleeping -diff-sample 5s
（换成深色桌布，再跑一遍，只把目录名改成 -dark-2026-09-30）
```

核对尺（明早当场就能用眼睛验的三条，一条一个动作）：
- **两枚目录都在、每枚都有那 3 张 png ＋ 1 张 `diff-table.txt`** ⇒ `ls docs/evidence/s1/62-diff-{light,dark}-2026-09-30/`；
- **深色那一行的 `px_delta_ge8` 不为 0、且 `timers=no`、`teardown=clean`** ⇒ 这一条是"球在深桌布上真的成像了"的数字面，与他的眼睛各判各的、互不替代；
- **两张 `-diff.png` 并排能看出是两颗球**（他自己说的"像／不像"就落在这两张上）。

⚠ 落点规矩：**票 62 `:82` 明文"差分截屏为证，不用自绘合成图"** ⇒ 必须走 `-diff`（真截屏），**不能走 `-shots`**（Go 侧合成）。这一条正是 §4 里那 20 张存档图不能当证据的原因。

## 4. 20 态那一次怎么演（名册核在不在 ＋ 跑法名对齐）

### 4a. 名册今天还在不在、枚数对不对——**在，20 枚，一枚不缺**

`ls docs/evidence/s1/ball-states/` 现读：**20 个文件、全部 `.png`、非-png 枚数＝0**。文件名与另一枚腿报的名册**逐字相同**：

```
01-FirstRun  02-Sleeping  03-Armed  04-Muted  05-Listening  06-Thinking  07-Acting
08-Confirming  09-AwaitingApproval  10-Speaking  11-Settling  12-Warm  13-Conversation
14-Downloading  15-Error  16-NoNetwork  17-Unconfigured  18-WatchdogAlert  19-Queued  20-Stuck
```
- 20 枚 mtime **全是 2026-09-19 21:38**（大小 5,525–13,533 B 不等），目录 mtime 2026-09-19 21:38 ⇒ **自那天起没被改过、没被覆盖过**。
- 入库那一发具名：`git log -1 -- docs/evidence/s1/ball-states/` ＝ **`1e80700a` 09-19 21:40** `feat(ball): fixed-size window + evidence cycle script + 20-state visual captures`。
- 对照表第 17 件列的 20 个态名（第一次运行…卡住）与上面这 20 个文件名**一一对上**（"已就位"＝`Armed`、"没配好"＝`Unconfigured`）。

### 4b. 跑法名对齐：票 07 `:51` 的 `a debug cycle page/window` ＝ 默认全周期，**不等于 `balldebug -tour`**

两枚原文都现读了，**它们是两件事**：

| 名字 | 原文逐字 | 现读实现 |
|---|---|---|
| `a debug cycle page/window`（票 07 `:51-52`） | 「**Visual: 20 states rendered in a debug cycle page/window; human screenshot review vs `design/screens/ball.html`**（colors/opacity/sizes per SPEC-08 §2.1）」——本格是 `- [ ]` 未勾，`:53` 追加的丁类原话「**待人项：需要一次桌面独占跑**（见台账 `A440` 第④节第 1 条）」 | `cmd/balldebug/main.go:311-332` 的 **`default:` 分支**，注释逐字「**Full 20-state visual cycle（`scripts/dev/ball-cycle.ps1` screenshots each step by polling the window title marker）**」：`for _, s := range allStates` 逐态 `SetState` ＋ 打 `balldebug: state=<名>`，`-cycle-ms` 控制每态停留。**跑法＝`build\balldebug.exe`（不加模式参数）** |
| `balldebug -tour`（票 65 `:63`） | 「**owner 第二次看 `-tour`，并明确说"这次对了"**——本票的完成判据是人的判决，不是测试绿。在此之前 SPEC-08 §2 的 INTERIM 标注**不得摘除**」 | `main.go:111-112` 帮助逐字「guided walk for owner sign-off, in one run: liquid flow -> voice-driven rotation -> idle border -> static Sleeping -> docked on all four edges -> popped back -> yours to play with」＝**7 拍向导，不是 20 态** |

⇒ **结论（明早怎么演）**：**两枚都要跑，各管各的件**。
- **第 17 件（20 态挨个看一眼）＝默认全周期**：`build\balldebug.exe -cycle-ms 6000`（默认 2000 太快，6 秒他能说完一句）。
- **主事件（他说"这次对了"那一次）＝`-tour`**（票 65 `:63` 点名的就是 `-tour` 这个名字，不是"cycle page"；台账 `R13` 记的他上一次的原话也是「owner 实机跑 `-tour` 后原话…」，见 `docs/reports/pending-and-issues.md:299`）。
- ⚠ 派单里"后来的等价跑法是 `balldebug -tour`"这句，**盘上证据不支持等价**：tour 连 `Thinking/Acting/Warm` 都不演（§2c 坑 #1）。

### 4c. ⚠⚠ 演之前必须知道的两条（一条是"对不上"，一条是"会覆盖档案"）

1. **这 20 张存档图是"旧样子"的图**，明早现场是新样子，**必然对不上**。具名出处：`docs/evidence/s1/62-visual-spec-draft.md:57` 逐字——「`ball-states/*.png`（20 张，`1e80700`，2026-09-19）| **冻结档**、且是 `-shots` 的 **Go 侧合成图** | 名义上 20 态，但票 62 AC#1 明文"**不用自绘合成图**" ⇒ **不能当 AC#1/#3 的证据**，只能当 §2.1 旧文本的示意图」。同一句在 `:152` 又数了一遍（13 个态**一条像素数字都没有**，那四个 `62-diff-*` 目录合起来只覆盖 7 个态）。
   ⚠ **一处证据自相矛盾，明早别当成结论讲**：`:57` 说这 20 张是 `-shots` 的 Go 侧合成图，但 `scripts/dev/ball-cycle.ps1` 的头注释与代码说的是**真截屏**（`$g.CopyFromScreen(...)`，并 `FindWindow("WispBallWindow","Wisp")` 找窗口，输出位就是 `docs/evidence/s1/ball-states/<NN>-<State>.png`）。两种跑法**产出的文件名一模一样**，本腿只读分不出当年用的是哪一枚。**要定这一条得跑一次才知道 ⇒ 记进 §7。**
2. ⛔ **`scripts/dev/ball-cycle.ps1` 默认会把那 20 张档案覆盖掉**：脚本里 `$OutDir` 缺省＝`Join-Path $repo "docs\evidence\s1\ball-states"`，开头 `New-Item -Force`。⇒ 明早**若要跑它必须带 `-OutDir` 指到新目录**（例：`-OutDir docs\evidence\s1\ball-states-2026-09-30`）。**台账"临时件只建不删"，档案被覆盖等于删。** 另外它还硬编码 `D:\work\base\go\bin\go.exe` 与 `E:\work\base\msys64\mingw64\bin`——**它会自己 `go build` 到 `%TEMP%\wisp-balldebug.exe`**（不是 `build\` 那枚），这一点正好和"那天不跑编译"打架 ⇒ **默认改用 §4b 那条手工命令，别用这枚脚本**，除非他真要看 20 张新档案。
3. ⚠ 明早**只用眼睛演 20 态**不需要脚本；脚本是为了**留档**。留档与否，等他到场前我判断（§7）。

## 5. 四件"看不了"的清单（他问起来，盘上可核的一句话证据在哪）

对照表 `docs/reports/desktop-signoff-2026-09-30.md:40-43` 那 4 件，逐件给一句**能点开的出处**（全部来自已有读数，**本腿没有为了证明去枚举 `design/**`**）。

### 第 7 件：那两张"球该长这样"的参考图，项目里从来没有过

**一句话证据（台账 `R13` 第②条，逐字）**：
> 「**一条必须承认的事实：图1/图2 从未落到磁盘。** 我逐个核对本会话 11 张 PNG 附件（商标查询 4 + 托盘 1 + AnySearch 文档 1 + 代理面板 1 + 本轮 tour 3 + 小裁图 1），无一为玻璃参考图 ⇒ 票 62 的**四个代理全部是按文字盲做配色**。」

出处：`docs/reports/pending-and-issues.md:305-307`（这一整条是台账 `R13`，标题行在 `:298`；同一份普查票面也留了一份：`.scratch/wisp/issues/62-liquid-glass-ball-visuals.md:167`）。
要求那两张图入库的判据原文＝票 65 `:57`「**参考图已入库（`design/refs/ball-glass-*.png`）**，且 `65-glass-reference-analysis.md` 给出逐条量化属性」——`- [ ]` 未勾，且那句"所以'赝品'这个判决在信息缺失下是必然的"把责任判给了缺图，不是做图的。

**他若问"为什么当时没存"**：台账同段已经替他答了并写了补救方向——「票 65 的第一步因此不是写码，而是**向 owner 取回参考图并入库 `design/refs/`**」（`pending-and-issues.md:309`）。

### 第 8 件：改完的球和参考图并排出图

判据原文＝票 65 `:60`「**并排差分产物在档：当前实现 / 返工后 / 参考图，同一区域同一底色**」（`- [ ]` 未勾）。**盘上不需要别的证据：它的前置就是第 7 件**，第 7 件的红句一引，这一件自动成立。明早**只登记"等他给图"，不演**。

### 第 9 件：那条"唯一样式真相源"的路径本机不存在

**红句逐字（这是最能点开的一份）**：
> `tokens_fourway_test.go:441: read design/assets/tokens.css: open …\design\assets\tokens.css: The system cannot find the path specified.`

出处（三处独立记录同一句，任选其一给他看，前两处最硬）：
- `docs/evidence/s1/145-snapshot-fields-landed-r1-accept-r1.md:344`（红句原文一行）
- `docs/evidence/s1/221-task-cancel-v1.md:168` — 「**在册常红 5 枚对上：panel 4＋ball 1，失败正文本腿现读逐字＝`read design/assets/tokens.css: …`**」
- `docs/evidence/s1/181-186-git-detection-census-c1.md:149` — **具名归因**那份：「`design/assets/` 被**别的会话**删了且未 staged——`git status --porcelain -- design/` 现量 4 枚 ` D`（`base.css`／`icons.js`／`theme.js`／`tokens.css`）。⇒ **不是本程造的**」

⚠⚠ **明早这句话要说全，否则会误导他**：对照表 `:42` 写的是"**此刻在你这台机器上不存在**"，而**盘上更完整的读数**是 `docs/reports/frontend-session-log.md:25` 那条 **X3**，逐字：「token 生成器的输入源已被 owner 挪出原位…该路径**盘上已不存在**（`git status` 显示 ` D design/assets/tokens.css`，**HEAD 里仍有 12942 字节**）。owner 把它**移动**成了未跟踪的 `design/old/assets/tokens.css`，两枚文件 md5 相同：`141c570f063606cf663629c717098f40`」。
⇒ **一句话给他**：「**东西没丢，是他在 9 月把它搬走了；库里的原件还在，只有本机这一份不在原位。**」这条**不是**"项目丢了真相源"。
另一句相关的定性（`docs/evidence/s1/193-frontend-gates-census-a1.md:76` T11）：「**不是**。文件不见了，与阈值无关」⇒ 明早不许被听成"阈值/测试被放宽"。

与今晚 CI 四枚报红的对应（对照表 `:75` 那句话的实底，本腿现读）：
`.scratch/wisp/probes/ci-red/ci-red-1.md:47` 列出那 4 枚逐名——`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／**`TestC21DesignTokensFourWayAgree`**；`:206` 逐字——「**这 4 枚在本机同样红**（09-29 20:04 那发逐名 `--- FAIL`）⇒ **明早真机签收时它们一样成立**；CI 那两发红没有额外改变明早能看到的东西」，同一行还给了结论尺：「**别把签收判成通过／不通过**，而不是 CI 挡住了签收」。

### 第 10 件：副屏

判据原文＝票 64 `:83-84`，逐字：
> 「- [ ] 多显示器：真拖到第二屏验证并留证；**本机无第二屏则保持未勾并写明所需硬件**。本机物理单屏（`\\.\DISPLAY4` 3440x1440），**硬件缺席**，未做任何双屏 mock。」

而且这条"未勾"**已经被非实现者裁过**：票 64 `:87-90` 逐字「对抗验收由非实现者执行…2026-09-20 由编排者执行：`docs/evidence/s1/64-adversarial-acceptance.md`（9 行 = 票面 9 框）。结论 **0 BLOCKER / 0 MAJOR / 3 MINOR**；**四个未勾框全部判"正确未勾"**」。
⇒ 明早他问起，一句到位：**「这台机器只有一块屏，这条当年就判'正确未勾'，没人拿假双屏糊弄。」**（对照表 `:43` 那句"加分项"的实底就在这行。）
补一条票 64 `:29` 的账：「**A4 多显示器实拖：detach→primary 与持久化/恢复有测，但'拖到第二块屏'从未真执行**」＝"代码侧测过、手没拖过"，别说成"完全没测"。

## 6. 明早的环境风险 ＋ 到场前的三条自查尺

### 6a. CI（今晚实测 ＋ 定时器的钟点都算过了）

- **今晚现读**：`gh run list --status in_progress` ＝ **空**（22:3x，rc=0，不是报错）。名册里最近几发全 `completed`：`slo-fresh` 16s／`ci` 那发 **1h40m23s**（09-28 23:53Z 起）／另一发 **2h11m30s**（09-27）。
- `slo-full` 确实走本机：`.github/workflows/ci.yml:590-591` 逐字 `slo-full:` ＋ `runs-on: [self-hosted, wisp-slo]`。
- **两枚定时器的钟点（换算到 +0800）**：
  - `ci.yml:36-37` `cron: '37 19 * * *'`（UTC）＝**北京 09-30 03:37 起一整包**。按历史最长 2h12m 推，**约 05:50 结束**，离 10:30 有 4.6 小时余量 ⇒ 正常情况不构成风险，但**明早仍要现查一次**（它可能因为我今晚的推送提前/叠加起发）。
  - `slo-fresh.yml:46-51` `cron: '23 */6 * * *'`（UTC）＝北京 **08:23／14:23** 等；实测每次只 **12–18 秒** ⇒ 08:23 那一发离 10:30 远，可忽略。
- ⚠ **今晚推送的账**：`ci.yml:12 on: push` ⇒ 我推一发，本机 self-hosted 就起一整包。台账 `A453` 写着「⛔ **最迟明天 10:30 签收之前我必推**」，而对照表 `:72` 与 `A452`／`A453` 状态行都写着「**那天前后……不推送**」。**这两句在盘上互相矛盾**⇒ 已记进 §7 第 1 条，**这是一个必须今晚定的决定，不是明早现场能补救的**。
- ⚠ **实测过的坑（具名）**：`A453` 逐字——推送起的那一整包「与 `235-r1` 抢 CPU，**今天已实测造过一枚争用假红**」。同一条口径对明早的含义：**签收期间任何 CI 在跑，都会让球掉帧／差分数字失真**，而他看到的"糊"就可能是这个。

### 6b. 还有什么会占住屏幕或机器

- **今晚现读：没有任何球进程挂着**。`tasklist | grep -iE "wisp.exe|balldebug.exe"` ＝ **零输出**（grep 返回码 1）；`tasklist //FI "IMAGENAME eq wisp.exe"` 与 `... eq balldebug.exe` 两发都逐字回「**INFO: No tasks are running which match the specified criteria.**」。
- ⚠ 今晚**没有在跑的编译**（`tasklist //FI "IMAGENAME eq go.exe"` ＝ 同上"无匹配"），**但台账说 `235-r1` 在飞**（`A453` 状态行「**在飞＝1 枚（`235-r1`）**」，起手 22:2x，正跑测试与突变）。⇒ **"今晚这一把尺读为零"不等于"它已交件"**，只能等它的交件通知。它交件之前我不能编译，这与派单一致。
- **有没有别的会话在改球**：`git log --since=2026-09-28 -- cmd/balldebug internal/ball docs/evidence/s1/ball-states` ＝ **空**（最后一动 09-21）；工作树那四条路径 `git status --porcelain` ＝ **空** ⇒ **球这一族今天没人动，明早不会有人在我眼皮底下改它**。
- ⚠ **别的会话确实在动这台机器上的文件**：工作树里有 16 枚 `design/**` 的未 staged 删除（`git status` 起手就记着，**不是我造的、我不还原不提交不删**，与 §5 第 9 件同一件事实）＋ `235`/`152`/`161` 的 probes 文件在改。⇒ **明早别把它们当"我动了界面那一侧"解释**。
- **球会把终端占住**：`-stay`／`-hold`／`-tour` 都以 `waitForExit` 收尾（托盘 Exit 或 Ctrl+C），**跑之前告诉他我要按 Ctrl+C**，别让他以为卡死了。
- **球会出现在哪儿**：`ball-cycle.ps1` 的历史做法是把光标挪到 `(4,4)`、把球放在主屏工作区右下内侧（`WorkingArea.Right-420 / Bottom-330`），理由逐字写着「**park the cursor out of the way so no third-party bubble overlaps the shot**（**input-method bars live there**）」⇒ **明早手动把球放离输入法状态栏远一点**，否则"看不清"是他那条 IME 白条而不是球。

### 6c. 三条"不派开窗代理／不跑整包／不推送"落到能自查的一把尺

他到场前（建议 10:05 前后）跑这一把，**七行全读一次就够**：

```bash
date "+%F %T %z"                                              # ① 时刻留痕
gh run list --status in_progress                              # ② 必须空 —— "不推送"的后果之一由它兜底
tasklist | grep -iE "wisp.exe|balldebug.exe"                  # ③ 必须零输出 —— 票 68 :74 的洁净前置
tasklist //FI "IMAGENAME eq go.exe"                        # ④ 必须只有 1 行 "INFO: No tasks ..." —— "不跑整包"的现场证明
reg query "HKCU\Control Panel\Desktop\WindowMetrics" //v AppliedDPI   # ⑤ 必须 0x60（＝96 DPI＝缩放 100%）
git log -1 --format='%h %ad' --date=format:'%m-%d %H:%M'       # ⑥ HEAD 前进了是正常的（235-r1 还在 commit，22:42 已到 `1d2ad737`）；要核的是下一行那条：
git log -1 --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- cmd/balldebug internal/ball   # ⑥b 必须仍是 09-21（3539d473／22015301）＝没人动过球
```

- **今晚的 ⑤ 实测＝`AppliedDPI REG_DWORD 0x60`＝96＝缩放 100%** ⇒ 与票 68 `:75`（"A.2 全部数字都是 96 DPI 下测的，否则不可比"，另见 A28/A29 的 DPI 不对称）**当场可比**。这条尺派单没给形状，本腿现场找到的，注意 Git Bash 里 `/v` 要写成 `//v`，否则逐字报「ERROR: Invalid syntax」。
- **"不派开窗代理"这把尺盘上没有自动读数**（没有注册表／台账记着"当前派了哪几枚开窗代理"）⇒ 只能靠**我自己不派**这一条，**加上** ②④ 两把间接尺（开窗代理必然伴随 CI 或编译）。**这一条要在 §7 挂着**：它不是一把能被打勾的尺。
- 判"干净"的口径：**②③④⑤ 四行都如期望 ⇒ 可开始**；任一不符 ⇒ **先说"等 X 分钟"再开始，不要边跑边解释**。
- ⚠ **还有一件明早可能缺档的**（对照表第 16 件自己写了"我明天之前补派"）：本腿现读 `ls docs/evidence/s1/ | grep -E '^62|^68'` 里**没有 `62-adversarial-acceptance.md`、`68` 名下零枚证据件** ⇒ **票 62 的 1:1 裁决表今晚仍不存在**。台账 `A453` 把它排在 `235-r1` 之后（「与 `235-r1` **互斥**，因为两枚都要跑突变」）。⇒ **若他问"这份是谁判的"，老实答"那张表还没出，是我该派而还没派的活"，别推给实现方。**

## 7. 没盘完／明早我自己现场定的

（本节⛔不许为空 — 实际有 7 条）

1. ⚠⚠ **"今晚推不推送"今晚就得定，明早定不了**：台账 `A453` 对我自己写死了「⛔ 最迟明天 10:30 签收之前我必推」，而对照表 `:72`／`A452` 状态行写的是「那天前后……**不推送**」。**盘上两句矛盾，我没有权限替 owner 裁**（推不推是编排者的动作，但"那天不推送"是给签收定的规矩）。⇒ 现场处置预案：**若 10:05 那把尺的第②行读到非空，就明说"CI 正在收尾，等它绿/等它停"，并把开始时刻整体后移**；绝不在 CI 在跑时开始演。
2. **那 20 张存档图到底是 `-shots`（Go 侧合成）还是真截屏**——`62-visual-spec-draft.md:57` 与 `scripts/dev/ball-cycle.ps1` 的头注释互相矛盾，本腿只读分不出（两枚跑法产出**同名文件**）。⇒ 只在"他问这图是怎么来的"时才需要答；**答法：说它 09-19 入库、是旧规格的样子、不做真截屏/合成的定性**。
3. **深色桌布用哪一张壁纸**——盘上没有任何记录。明早现场让他自己挑一张深色，我不动他的桌面设置（改设置算动他的机器）。
4. ⚠ **第 3 件"五种干活模样"明早到底用哪条跑法**：§2b 给了默认全周期，但**每态 6 秒够不够他说一句"这俩我分不出"**没量过（历史 dwell 是 2000ms，票 68 的 `-diff-dwell` 默认 6s）。⇒ 现场定 `-cycle-ms`，**先给 6000，他嫌慢再降**，第一次没有基线可参照。
5. **`-frozen` 反例要不要真的演**（对照表第 13 件说"很适合演给你看一次"，但它是**数字判据**，且明早**不跑整包**）：演一次只多占 1–2 分钟，但它会让"新旧两版"在同一块屏幕上同场，**他有可能把两版混成一谈**。⇒ 现场决定；若演，**演完立刻重跑一次新版 `-tour` 收尾**，让他最后看到的是新样子。
6. **"那扇小窗"（聊天面板）明早开不开**：`balldebug` 那两枚是 `stub, ticket 33`（§2c 坑 #3）。开面板要另建 `wisp.exe`（§1 判它也过期），而且票 68 的 AC#2 那三项数（`#3/#4/#5`）**还没裁**（票 68 `:69` 逐字「⛔ AC#2 被 R15 三项挡住，未答之前不得翻转」）⇒ 也就是说**出厂默认现在还是旧样子，日常启动画的仍是小点**（对照表第 12 件原话）。**明早只看不翻这一件，且绝不把面板说成已签收。**
7. **明早需不需要留 20 张新档案**：留＝跑 `ball-cycle.ps1`，而它会**① 自己 `go build`（与"那天不跑编译"冲突）、② 默认覆盖 `docs/evidence/s1/ball-states/` 那 20 张 09-19 的档案**。⇒ 本腿倾向**明早不跑它、只演**，留档另排一枚腿；但这是编排者的决定，**现场定**。若真跑：**必须带 `-OutDir` 指新目录**。
