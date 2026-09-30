# 239-v1 — 球被自己的窗口切成方形（RingMarginPx 8→22 / DockTriggerPx 16→24 / haloExtents）：出自非实现者的 1:1 裁决表

- 本表由**验收腿 239-v1** 写。本票实现者＝编排者本人（票面 `:5`），裁决者≠实现者（`AGENTS.md` §0.3）。
- 被裁改动锚点：HEAD `0589fd9cc3c48eab85950475d25505f447c0a425`（本腿 `git log -1` 现取，11:24）；差分 `d5a59d66..0589fd9c` 实到产码＝`internal/ball/hit.go`／`tokens.go`／`renderer_windows.go` ＋ `docs/evidence/s1/c21-native-tokens.md:166/:168` ＋ 票面与探针件 ＋ 台账 A462。
- ⛔ **本表不给"整体通过"**；工单所有 `- [ ]` 复选框一枚都不勾（勾框归编排者，票面 `:31`）。
- 凭据口径：**〔我现跑〕**＝本腿当场跑尺并记录读数；**〔读台件，未复跑〕**＝引票面／A462 的读数，本腿未复现。
- 突变一律走 `go test -overlay` / `go build -overlay`，产物只落 `.scratch/wisp/probes/239/`；工作树 `internal/` `cmd/` 收尾必须干净。

## 0. 逐格判据（1:1 表；每格裁完追加对应小节并 commit）

| AC | 判据（票面摘要） | 判语 | 凭据（口径标注） |
|---|---|---|---|
| AC#1 | 复跑四把尺（build／gofumpt／vet／d22scan）＋两形差分＋`internal/ball` 红名册必须＝"只有 CSS 那一枚" | **成立**（附一处**台件口径偏离**，具名见 §1） | 四把尺＋两通测试＋差分全〔我现跑〕；唯 A462 那句"改后 Warm 66x66/3359＝圆盘 98%"〔读台件，未复跑，且本腿判它出自旧几何台件〕 |
| AC#2 | 耦合突变：`DockTriggerPx`→16 与 →21（`RingMarginPx` 保持 22）⇒ `TestBallLiveEdgeDock`／`…Hover` 必须红；21 若绿则具名登记"耦合无仪器" | （待裁） | （待补） |
| AC#3 | `haloExtents` 恒等性突变：`R*1.9 <= fit` 时强行互换 fill/grad ⇒ Sleeping 的 `px>=8` 必须离开 2098/2103/2120 那一族 | （待裁） | （待补） |
| AC#4 | "窗口放大"代价三件：① work area 内停手时 orb 离边最少 22px（旧 8）；② `DockOverlapFrac` 一字节没动；③ 差分分母 `of 97344` 不变 | （待裁） | （待补） |
| AC#5 | 归口：本票动过 `internal/ball` 三枚文件⇒与票 65/68 同地界；读数须并入票 68 名下欠账，不另立平行真相 | （待裁） | （待补） |

## 1. AC#1 — **成立**（附一处具名口径偏离，见本节末段）

**判据（票面 `:33`）**：复跑四把尺（build／gofumpt／vet／d22scan）＋两形差分（`-diff` 指临时目录）＋`internal/ball` 红名册必须＝"只有 CSS 那一枚"。

**本腿现跑读数（〔我现跑〕，窗口 11:27:54–11:34:23）**：

| 尺 | 读数 |
|---|---|
| `go build ./...` | rc=0（11:27:54） |
| `gofumpt -l internal/ball/hit.go tokens.go renderer_windows.go` | 零输出，rc=0（11:27:54） |
| `go vet ./internal/ball/ ./cmd/balldebug/` | 空，rc=0（11:28:08） |
| `./tools/d22scan/d22scan.exe`（独立 module 直跑） | clean rc=0（11:28:09）；`ban #8 internal/=462`、`cmd/=63` 是**被扫文件数**不是违规数，与票面逐字同 |
| `go test ./internal/ball/ -count=1`（Sherpa PATH 带上） | rc=1；`--- FAIL` 名册＝**仅 `TestC21TableColourRowsMatchTokensCSS`** 一枚，红句逐字 `read design/assets/tokens.css: … cannot find the path specified`（11:33:39，包 0.088s） |
| `go test ./internal/ball/ -tags winlive -count=1` | 跑前/跑后 `tasklist //FI "IMAGENAME eq balldebug.exe"` 读数均 **0**（桌面独占自证）；rc=1；红名册仍**只有 CSS 那一枚**；两枚 dock 用例不在红名册＝绿（11:34:11–23，包 11.047s）。产物 `probes/239/test-ball-winelive-r1.txt`，SKIP 行 0 |
| 差分（自建 `probes/239/bin/balldebug-v1.exe`＝`go build ./cmd/balldebug` @ HEAD） | **edge=100**（`WindowEdgePx(56,96)=56+2*22=100`，hit.go:54-58）；Sleeping `px>=8=2120`/box 44x44、Warm `5284`/66x66、Speaking `5544`/80x80；分母 **115600**=340²=(100+2*120)²。产物 `probes/239/diff-v1/diff-table.txt`（11:30:56–11:32:04） |

**形状口径一句话**：改前 Speaking 差分 box `72x72`＝整窗（旧 edge 72）、`px>=8=4803` 占窗口面积 5184 的 **93%**〔读台件（票面 `:11`），未复跑〕；本腿 HEAD 现跑三形 box 全部**严格小于** edge=100（Speaking 80x80＝窗口面积 10000 的 64%，`px>=8=5544` 占分母 115600 的 **4.8%**）⇒ "外沿＝窗口矩形"这一形在 HEAD 上没了。边距 22 的最小性推导〔我现读码〕：波纹外沿 `R+5*s+t*16*s+RingStrokePx/2`，t=1 时＝`R+21.75`（renderer_windows.go:450-457，`RingStrokePx=1.5` tokens.go:450）⇒ 边距最小整数 22。

**Sleeping 恒等复现**：本腿现跑 `px>=3=2172 / px>=8=2120 / max=171 / box 44x44`，与台件 `62-diff-border`、`62-diff-signoff` 两档同名次逐字同族（`px>=24` 漂 7 枚：1394 vs 1387＝读数的 0.5%；`mean` 差是**分母差**所致，1.134*97344/115600=0.954 与本腿 0.953 吻合）⇒ 家族 2098（62-diff-glass，其分母 63504）/2103（12-ball-walk）/2120 中，本腿**正面命中 2120**。

**具名偏离（不影响判语、影响一句台件的采信）**：票面 `:29`／A462 的"改后 `Warm 66x66 / 3359px`＝同半径圆盘面积 98%"**在 HEAD 二进制上复现不出**——本腿现跑 Warm 同为 66x66 但 `px>=8=5284`。原因〔读台件＋几何推算，我现验算〕：`signoff-30/fixed-states`、`fixed2-states` 两档"改后"表 `edge=72`、`px_total=97344`，而 HEAD 窗口 edge=100（本腿实测）⇒ 那两档**不是 HEAD 几何的读数**（推断＝当时挂的是中间态的 `build/balldebug.exe`，该产物在 `.gitignore:14` 的 `build/` 里、不随提交走）；"3359/π·33²=98.2%" 算术成立但分子口径错。**票面自己要求"复跑"，本腿复跑结果支持结论、不支持那句 98%。**

## 2. AC#2 — （待裁）

## 3. AC#3 — （待裁）

## 4. AC#4 — （待裁）

## 5. AC#5 — （待裁）

## 6. 没做完／留给编排者（⛔ 不许为空）

（待补）

## 7. 我攻不动的地方（⛔ 不许为空；末节）

（待补）
