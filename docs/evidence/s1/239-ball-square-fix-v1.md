# 239-v1 — 球被自己的窗口切成方形（RingMarginPx 8→22 / DockTriggerPx 16→24 / haloExtents）：出自非实现者的 1:1 裁决表

- 本表由**验收腿 239-v1** 写。本票实现者＝编排者本人（票面 `:5`），裁决者≠实现者（`AGENTS.md` §0.3）。
- 被裁改动锚点：HEAD `0589fd9cc3c48eab85950475d25505f447c0a425`（本腿 `git log -1` 现取，11:24）；差分 `d5a59d66..0589fd9c` 实到产码＝`internal/ball/hit.go`／`tokens.go`／`renderer_windows.go` ＋ `docs/evidence/s1/c21-native-tokens.md:166/:168` ＋ 票面与探针件 ＋ 台账 A462。
- ⛔ **本表不给"整体通过"**；工单所有 `- [ ]` 复选框一枚都不勾（勾框归编排者，票面 `:31`）。
- 凭据口径：**〔我现跑〕**＝本腿当场跑尺并记录读数；**〔读台件，未复跑〕**＝引票面／A462 的读数，本腿未复现。
- 突变一律走 `go test -overlay` / `go build -overlay`，产物只落 `.scratch/wisp/probes/239/`；工作树 `internal/` `cmd/` 收尾必须干净。

## 0. 逐格判据（1:1 表；每格裁完追加对应小节并 commit）

| AC | 判据（票面摘要） | 判语 | 凭据（口径标注） |
|---|---|---|---|
| AC#1 | 复跑四把尺（build／gofumpt／vet／d22scan）＋两形差分＋`internal/ball` 红名册必须＝"只有 CSS 那一枚" | （待裁） | （待补） |
| AC#2 | 耦合突变：`DockTriggerPx`→16 与 →21（`RingMarginPx` 保持 22）⇒ `TestBallLiveEdgeDock`／`…Hover` 必须红；21 若绿则具名登记"耦合无仪器" | （待裁） | （待补） |
| AC#3 | `haloExtents` 恒等性突变：`R*1.9 <= fit` 时强行互换 fill/grad ⇒ Sleeping 的 `px>=8` 必须离开 2098/2103/2120 那一族 | （待裁） | （待补） |
| AC#4 | "窗口放大"代价三件：① work area 内停手时 orb 离边最少 22px（旧 8）；② `DockOverlapFrac` 一字节没动；③ 差分分母 `of 97344` 不变 | （待裁） | （待补） |
| AC#5 | 归口：本票动过 `internal/ball` 三枚文件⇒与票 65/68 同地界；读数须并入票 68 名下欠账，不另立平行真相 | （待裁） | （待补） |

## 1. AC#1 — （待裁）

## 2. AC#2 — （待裁）

## 3. AC#3 — （待裁）

## 4. AC#4 — （待裁）

## 5. AC#5 — （待裁）

## 6. 没做完／留给编排者（⛔ 不许为空）

（待补）

## 7. 我攻不动的地方（⛔ 不许为空；末节）

（待补）
