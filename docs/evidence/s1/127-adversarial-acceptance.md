# 票 127 独立对抗验收 —— 常驻 GUI 腿那 20 行 install 的第一枚钉（`R-117-A`，票 117 结案的唯一条件）

**验收方**：`acceptor-ticket127`（本文件唯一作者）
**日期**：2026-09-22（本会话）
**被验交件**：票 127，实现方 `agent-ticket127`
**基线**：`6a39820`（票 121 收件、本票开工前的树；实现方选它作控制组，我沿用）
**交件 sha**：`8663a39`（用例 `cmd/wisp/resident_sink_nail_127_windows_test.go` 533 行）、
`c927738`（AC#1 票面）、`938fda1`（AC#2 两处注释 + 票 117「更正一」）、`fd82bf4`（AC#3 票 117「更正二」）、
`9c335d0`（AC#4 票面）、`8ec04f3`（收尾补记，**最终 sha**）
**本文件性质**：1 AC = 1 格，裁决 + 三档标签；**票面一字未改、别人的框未翻**；`R-*` 只给建议归属，
`docs/reports/**` 一字未动（那本账编排者记）。

---

## 0. 会话口径（先声明我怎么做，免得读数无法复算）

- **只读**：本会话在仓库内唯一的写件是本文件。全程
  `git archive <sha> | tar -x -C /tmp/ac127-s22/<名>` 取快照，**没有**在仓库树内建 worktree、
  没有 `checkout`/`--amend`/`reset`/`rebase`/`stash`（A38④）。
  `internal/winsec/**` 有另一个写者在飞（票 129）⇒ 本会话**没有一次**读工作树的 winsec 状态、
  也没有对它做过写，一切走快照。
- **本机命令前缀**：`PATH` 必须挂上快照里的 `third_party/sherpa-onnx`，否则 `cmd/wisp` 的测试二进制
  在**进程加载**时就死（`exit status 0xc0000135` = STATUS_DLL_NOT_FOUND，票 98 那一格）。
  ⚠ 这条是**我自己踩到的**：第一次跑 `-run TestAC1ResidentLeg` 得到
  `FAIL github.com/CarlosShao/wisp/cmd/wisp 0.042s` + `rc=1`、`=== RUN` 命中 **0** 条 ——
  四数全零而 rc 非零，是"没跑"不是"跑红"。`git archive` 不含 dll ⇒ 三枚 dll 要自己拷进快照
  （`onnxruntime.dll` / `sherpa-onnx-c-api.dll` / `sherpa-onnx-cxx-api.dll`），
  这与实现方 AC#1/AC#4 两格自报的"票 117 §六 那枚仪器坑照抄"是同一件事，**复算成立**。
- 快照清单（渐进补）：
  | 名 | 内容 | 用途 |
  |---|---|---|
  | `/tmp/ac127-s22`（＝`/tmp/ac127-s22/post`） | `git archive 8ec04f3` + 三枚 dll | 交件树：基线四数、AC#3 分母、门禁、变异 |
  | `/tmp/ac127-s22/mut-m4` | `post` 副本 + 删 `resident_windows.go` 那 20 行 | **M4 重量（本票生死格）** |
- 四数一律从 `-v` 输出量，`-count=2` 才不缓存：`grep -c '^=== RUN'` / `'^--- PASS'` / `'^--- FAIL'` /
  `'^--- SKIP'`（顶层锚定，子用例缩进另计）。
- 变异先发落地再读名：`grep`/`diff` 出被改的那几字节 → `go build ./cmd/wisp` rc=0 → 才读 `--- FAIL` 名单。
- `GOOS=linux go vet` 只编译不执行。

**本文件写作状态：骨架已立（渐进写，每格一次 commit）。**
逐格裁决一览（未落格的写"待"）：

| 格 | 裁决 | 标签 |
|---|---|---|
| AC#1 常驻腿钉子 / M4 重量 | 待 | 待 |
| AC#2 R-117-C 降级是否诚实 | 待 | 待 |
| AC#3 票 117 §一 现状表重算 | 待 | 待 |
| AC#4 门禁 + 四数 + 台账 | 待 | 待 |
| **总判** | 待 | 待 |
