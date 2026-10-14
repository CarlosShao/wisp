# 114-v2 探针锚 — 票 114 AC#10「第三装法钉」的非实现者判语（验收腿）

- 身份／时刻：`114-v2`，2026-10-14 16:5x–17:1x +08；起手 HEAD `7e8ac3fa`；被裁＝`114-r6` 产码 `3b6f170f`（`internal/panel/bound_dispatch_door_114_test.go` 530 行）＋交付 `c499adac`。
- 窗口：他腿 `318-r1`（`internal/agent/approval`）在飞（`replies.go` 工作树脏）；`frontend/**` 被界面侧续写（工作树里 `design/assets/*` 被删、fixture 16:17 重生）；本波只读，仓内零突变（5 枚关键件 sha256 起手＝收手）。
- 方法：牙全部种在仓外树 `D:\tmp\wisp114v2\tree`（`git archive HEAD` 真盘＋`mut.py` 逐发还原 pristine 字节；A713：钉读盘，⛔ 不用 overlay 抵账）；`GOTMPDIR=C:\tmp\wisp114v2gotmp`（反斜杠形）；`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`；live 读数为真树。
- 判语：**AC#10 成立**（编排者重裁口径＝"让可达性钉认第三装法"）＋留档句①②（见裁决表 §1/§7）；翻框建议＝可翻。
- 支点（现读）：`cmd/wisp/panel_host_windows.go:80` `panelDispatchBinding = "wispDispatch"` → `:842` `w.Bind(panelDispatchBinding, …)`；`:892` literal 探针门；单参 `cmd/wisp/resident_task_source_windows.go:311` `ws.Bind(run.paths)`；页真读 `frontend/src/lib/panel.ts:350` `const door = window.wispDispatch;`（prose 在 `:320–:346` 的 `/* */` 内、`:329`）。
- 读数：`logs/01` 基线、`logs/02` 树 A/B 对照、`logs/03` 九发牙（8 红 1 绿＋2 发旧尺盲区复现）、`logs/04` 门禁。
- 交件：`docs/evidence/s1/ticket114-114-v2-accept-v1.md`；⛔ 框未翻、未动 `Status:`/票名、零 push。
