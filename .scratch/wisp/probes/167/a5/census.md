# 167-a5 — 撞钉预检（只读普查腿）：`cmd/wisp`＋`internal/panel` 今天所有"计数型"钉名册，逐枚裁"落地腿 167-r2 一旦新增字段／新增那一跳，哪些绿用例变红"

> 运行类型＝**只读**：⛔ 零 `go` 命令（`go test`／`go build`／`go vet`／`go env` 一律未跑，因写腿 `268-r1` 正在 `cmd/wisp` 取整包终态）。
> 本件的"红/不红"一律是**射程判断**（读尺体本体得出），颜色属〔预测〕，逐枚归口给编排者跑。
> 前一腿 `167-a3` §4 已做过**占用那一枚**的撞钉射程；本腿**独立复测**（⛔ 不抄它的读数），且把面从"占用"扩到**四枚输出全量＋`cmd/wisp` 侧**。
> 派单里每一句行号都当**待验断言**处理：复认成立才写，漂了具名报回。

## 0. 起手锚（逐字读数，同一发命令取）

| 项 | 读数 |
|---|---|
| 时刻 | `2026-10-05 14:00:02 +0800` |
| `git rev-parse --short HEAD` | `94380ab3` |
| 分支 | `dev` |
| 五根 porcelain（`git status --porcelain -- cmd internal tools scripts docs \| wc -l`） | **1 行**，逐字 ` M cmd/wisp/resident_approval_windows.go` ⇒ **别人的脏，本腿只报不改**（本腿对该根贡献 0 行；见 §7 自证） |
| 尺面枚数（起手即量） | `cmd/wisp` **99** 枚 `.go`／其中 `*_test.go` **65** 枚；`internal/panel` **34** 枚 `.go`／其中 `*_test.go` **20** 枚 |
| 本腿写面 | 仅 `.scratch/wisp/probes/167/a5/census.md`（＋同目录临时 msg 件）；⛔ 未动任何产码／测试／工单／台账；⛔ 未碰任何票面 `- [ ]` 框 |
| 阅读禁令自证 | grep/find 根一律显式写 `cmd internal tools scripts docs .scratch` 或其子集；⛔ 从未拿 `.` 当根；`frontend/**` 与 `design/**` **未读一字**（本件凡涉及页面侧只引 Go 尺体的字面路径，不引内容） |

## 1. 问一（本腿主体）— 今天 `cmd/wisp` 与 `internal/panel` 所有"计数型"钉

未判。

## 2. 问二 — 四枚输出各自的现有断言现状（占用／序号／停止／草稿＋崩溃自救）

未判。

## 3. 问三 — `internal/panel` 三枚冻结件的射程（`tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`frontend_hygiene` 一族）

未判。

## 4. 问四 — C17 入向方法名名册钉：现量枚数＋"第五枚入向方法名"的代价

未判。

## 5. 问五 — 分位数／真窗那一族的脆弱点与"该钉哪两行"

未判。

## 6. 问六 — 给落地腿的一页清单

未判。

## 7. 量不到／判不动／我写错的读数（自我对抗）／⛔ 未动产码自证

未判。
