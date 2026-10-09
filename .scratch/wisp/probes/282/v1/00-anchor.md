# 282-v1 起手锚（只读复核腿）

- 取样时刻：`2026-10-09 10:2x +08`（`date "+%Y-%m-%d %H:%M %z"` 实测）
- 分支：`dev`
- 起手 HEAD：`7789f9532281d9de633eb8a6e2198f2ddb1b7c0f`（`242-v2 三格重裁…`，提交时间 `2026-10-09 10:22:33 +0800`）
- `git status --porcelain` 行数：**758**（只登记，不清、不 stash、不 checkout）
- 本腿射程：AC#2／AC#4／AC#5 三格。AC#1／AC#3 **按住**（撞 `242-v2` 的突变面，串行铁律），见 `20-held-cells.md`。
- 本机 `git config --get core.autocrlf` = `true` ⇒ 工作树带 CRLF 不等于该文件有债（AC#4(b) 两列并排尺的根据）。
- 名册起点（票面转述三枚 sha，已现量复认）：`133b1bfa`／`bd124b2a`／`0abe217c` 三枚**均存在**，作者链同一。
  - `git log --oneline -- internal/tools/cancel.go internal/tools/task.go internal/agent/subagent_197.go -- ':(exclude).scratch'` ⇒ 相关文件里只有 **`bd124b2a`** 一枚是本腿的嫌疑人；`c44b30c4`／`b9fa815b`／`62d7d3b4` 等是更早的别票历史，非本轮对象。⇒ 转述的「三枚」里 **两枚（`133b1bfa`／`0abe217c`）根本不碰这三条路径**，一枚是探针锚、一枚是交付件 md。详见 `10-ac2-ac4-ac5.md` AC#5 名册。
- 先例登记：`.scratch/wisp/probes/282/v1b/`（commit `369a8f29`，另一条腿 `282-v1b` 已交过一份五格判决，含突变）。本腿**不引用其读数为凭据**，全部现量自取；v1b 只作为"此票非首次裁"的背景。
- 写面（本腿唯一允许的三枚新建件）：`.scratch/wisp/probes/282/v1/00-anchor.md`／`10-ac2-ac4-ac5.md`／`20-held-cells.md`。
- 纪律复述：零 `.go` 字节改动；不跑 `go build|test|vet|list|doc`；不跑 `scripts/*.sh`；取数一律 `git show HEAD:<path>`／`git grep … HEAD`；票面 `.scratch/wisp/issues/282-*.md` 一字不动；台账不动；只 commit 不 push。
