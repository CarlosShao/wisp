# 票 33 面板宿主 —— 落地腿 `33-r9`：第五形现量 ＋ "被具名拒绝"的出口

**本腿代号**：`33-r9`（写码腿／实现者）。写面：`cmd/wisp/**`（产码最多两枚：`panel_host_windows.go`＋`panel_resident_windows.go`）、
本件、`.scratch/wisp/probes/33/r9/**`。⛔ 票面 `.scratch/wisp/issues/33-panel-host-c27.md` 一字未动。

**起手锚点（同发现取，`2026-10-01 17:25:03 +0800`）**：HEAD＝`426a4d03`（`Thu Oct 1 17:20:32 2026 +0800`，
`docs(evidence): 33-v2 terminal re-run plus the leg's own probe artifacts`），分支 `dev`。
起手 `git status --porcelain`＝**272 行**（台件 `.scratch/wisp/probes/33/r9/00-start-git-status.txt`）＝
本仓常态（`design/**` 是 owner 自己的未提交件），⛔ 不还原、不提交、不删。
起手 `git status --porcelain -- cmd internal`＝**0 行**（同发现量，见 §5）。

**这格缺陷的出处（不是我造的题）**：非实现者终裁表 `docs/evidence/s1/33-panel-host-c27-v2.md`
**§问①**（"拒绝之后还有没有恢复路径？⇒ 没有。今天这形是'永久拒绝'"）与
**§问②**（"那张四形对照表里'上一任'永远是被自己 Destroy 掉的那扇窗……表里没有第五形"）。

---

## 0. 取数口径（钉死）

- 整包/名册命令逐字：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp -count=1 -run 'Panel|AC13|AC14|AC4|BallPanel|BallGesture|ListeningSocket' -v`
- 红绿只认 `--- FAIL`／`--- PASS`／`--- SKIP`；`t.Logf` 的 `file:line:` 前缀不算失败。
- 前人读数一律具名并注"本腿是否复跑"。⛔ 不拿别人的发当自己的凭据。
- `staticcheck`：⛔ 不跑（本机版与 CI 钉版不同，跑了就是假绿）⇒ 那一格标〔未复认〕。
- `go mod tidy`／`go get`：**未跑**（HEAD 上 tidy 退 1，票面 `:267`）。
- `frontend/**`／`design/**`：⛔ 未读、结论未引。
- `internal/**`、`docs/reports/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/SLO.md`、`docs/BUILD.md`、
  `thresholds.go`／golden／任何阈值：⛔ 零改动。

## 1. 判据进度表（骨架，逐格交完改成读数）

| 格 | 要交的凭据 | 状态 | 读数落在哪 |
|---|---|---|---|
| ① 第五形 F1（活 WebView2 窗＋投给它的 WM_CLOSE，跑产品自己的 bringUp） | 检查读到谁的 hwnd／拒绝有没有发生／那扇活窗后来怎样／第二次 Show 的结果，各 3 发 | 待跑 | §2 |
| ① 对照形 C（窗已 Destroy、close 已被自己泵掉） | 同格读数 3 发，证两形在检查眼里可区分 | 待跑 | §2 |
| ① 第五形的收摊观察（线程退出时那扇活窗会怎样） | 逐字读数，只量不改语义 | 待跑 | §2 |
| ② 选型（甲＝拒绝即线程级致命／乙＝stop 后重起一条） | 选一条，另一条写明被否的凭据 | 待定①读数 | §3 |
| ② 会响的钉 | 一枚新用例，钉"第二次 Show 不再走同一条永久拒绝路"＋"拒绝后 `post()` 立刻 false" | 待写 | §3 |
| ② 定向突变 | 中和新加那一处 ⇒ 指名用例当场具名红，逐字红句入表 | 待跑 | §3 |
| 名册复跑 | 派单给的 20 枚逐名对表＋新增枚具名 | 待跑 | §4 |
| 门禁四数 | `go build ./...`／`go vet ./cmd/wisp`／`sh scripts/d22scan.sh`／`gofumpt -l <名下件>`，各 rc 入表 | 待跑 | §5 |
| ⑤ 我可能写错的条目 | 逐条带"如果错了后果" | 待填 | §6 |
| ⑥ 我跑了哪些尺 | 逐条带命令 | 待填 | §7 |
| ⑦ 判不动的地方 | 逐条具名交回编排者 | 待填 | §8 |

## 2. ① 第五形现量（只量读数，不改语义）

（待填：逐字读数）

## 3. ② 把"被具名拒绝"变成有出口的（产码）

（待填：选型＋凭据＋钉＋突变）

## 4. 名册复跑

（待填）

## 5. 门禁读数

（待填）

## 6. ⑤ 我可能写错的条目

（待填）

## 7. ⑥ 我跑了哪些尺

（待填）

## 8. ⑦ 判不动的地方

（待填）
