# 245-r1 — 产码腿：球进常驻进程后裸 `Esc` 被抢成全局热键 ⇒ 走乙（稳态不绑 cancel，只在 `Confirming` 借、离开还）

- 腿：`245-r1`（产码腿，不是裁决腿）。工单：`.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md`。
- 本表由**实现者本人**填写 ⇒ 它的性质是〔实现者自述〕，**不是**非实现者的对抗验收裁决。缺口审计与对抗验收必须由另一枚 agent 做（`SPEC-12 §4.3` #1/#3、D22 双角色）。⛔ 本表任何一格都不构成"AC 已勾"。
- 骨架先交（09-30 17:0x，第 5 轮内），之后逐节填。硬预算闸门：跑到第 100 轮先写满 §3／§5 并 commit。

## §0 起手名册与结论占位

- 起手 HEAD：`25ce3ce9`（编排者 245 票面 AC#3② 更正那一枚）；起手 `git status --porcelain` 全量名册已逐字抄进 `.scratch/wisp/probes/245/r1/start-status.txt`（本腿取数时刻 09-30 17:0x +08；共享工作树，终态判据＝**等于这一份**，不是"必须为空"）。
- 起手时 `internal/ball`＋`cmd/wisp` 的写面：见 §0.1 现读。
- **结论（逐格判语）**：见 §1。本骨架阶段一律"（填写中）"，§1 逐格填完后本节给一句话总结。

### §0.1 起手写面（本腿自己的尺）

（待填：`git status --porcelain internal/ball cmd/wisp docs/evidence/s1` 的起手读数）

## §1 判据逐格表（AC#1／AC#2／AC#3／AC#4／AC#5）

| 格 | 判语（成立／不成立／阻塞于票 228／未做） | 尺与现读 | 正控（ mutation 名 ） |
|---|---|---|---|
| AC#1 稳态名册不含裸 `Esc` | （待填） | | |
| AC#2 `Confirming` 内借、离开还（两头都要读数） | （待填） | | |
| AC#3 三枚钉同批改期望＋收紧后仍能红 | （待填） | | |
| AC#4 用户能改 `[hotkey] cancel` | （待填） | | |
| AC#5 残余五字段登记 | （待填） | | |

## §2 突变名册（逐发：改了哪一行／指名用例／红句／还原）

（待填：名册表头 `编号 / 位置 / 突变内容 / 期望被红的那枚用例 / 实际红句 / 还原方式 / 复认读数`）

## §3 门禁读数（硬预算闸门：跑到第 100 轮先把这一节写满并 commit）

- `go build ./...`：（待填）
- `go vet`：（待填）
- `go test ./internal/ball/...`（默认档，无 `winlive` tag）：（待填）
- `go test -tags winlive ./internal/ball/...`（真机档）：（待填）
- `go test ./cmd/wisp/...`（需 `PATH=$PWD/third_party/sherpa-onnx:$PATH`）：（待填）
- `tools/d22scan`：（待填）
- 起跑前的进程尺（`balldebug.exe`／`wisp.exe` 计数都必须＝0）：（待填）

## §4 残余（乙形那半条代价与新增的具名欠账）

（待填）

## §5 判不动的地方（硬预算闸门：与 §3 同批写满）

（待填：真判不动就写"判不动＋为什么"，不许留空）

## §6 名册核对（终态＝起手名册）

（待填：终态 `git status --porcelain` 与 §0.1 的差分逐条归因）
