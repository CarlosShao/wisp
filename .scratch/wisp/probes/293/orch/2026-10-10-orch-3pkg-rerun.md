# 票 293 · 编排者自己那一把三枚包集原尺（`AC#4`，2026-10-10 12:0x）

**性质**：〔仅本机可量〕那一族之外的**普通整包尺**，我 12:0x 自己跑——因为 `293-v1` 具名报了它⛔ 没跑票面那把（它跑的是**四枚并集**，多一枚 `./internal/panel/`），而票面 `AC#4` 写的原尺是**三枚包集**。⇒ 两把数⛔ 互替，本件只交我这一把。
**来路**：非实现者验收腿 `293-v1` 交件后，我据它翻 `AC#2`/`AC#3`/`AC#4`（见票面「编排者裁定（2026-10-10 12:0x）」节与台账 `A801`、停车点 §4.0bb）。

## harness（逐字，可重跑）

```
tasklist 现量：wisp.exe／balldebug.exe／*test.exe ＝ 0 枚（起手闸门）
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ ./internal/ball/ -count=1 -v
```

原始 stdout＋stderr 落盘＝`logs/orch-3pkg-run1.txt`（**361,593** 字节，`go test` 自己的 CRLF 形态）；同一份内容做了一枚 **CR→LF 副本** `logs/orch-3pkg-run1.nl.txt`（`tr '\r' '\n'` 是**逐字节等长替换**，所以它同为 **361,593** 字节——⚠ 我第一稿在这里写"剥 CR 后 355,814 字节"，那枚数我⛔ 没量过就写了＝**假读数**，就地按 `find -printf '%s %p\n'` 的现量改正，并留这句痕）。
⚠ 两份⛔ 不入本 commit（大输出不入库、⛔ 不删，只留盘上；下面所有数都来自这两份，尺附在每条后面）。
计数尺逐字＝`tr '\r' '\n' < logs/orch-3pkg-run1.txt > logs/orch-3pkg-run1.nl.txt`，然后 `grep -c -- '--- PASS'` 等（⛔ 不转形态直接 `grep` 会把整包读成一行，那是一枚假读数）。

## 读数（三数＋两枚名册）

| 尺 | 读数 |
|---|---|
| `grep -c -- '--- PASS'`（**含子测试**） | **507** |
| `grep -c -- '--- FAIL'` | **6** |
| `grep -c -- '--- SKIP'` | **3** |
| `grep -c -- '^--- PASS: TestAC293'`（顶层） | **7**（`--- FAIL: TestAC293` ＝ **0**） |
| 包行 | `FAIL cmd/wisp 468.955s`／**`ok internal/audio 15.977s`**／`FAIL internal/ball 0.213s`；`rc=1` |

`--- SKIP` 三枚逐名＝`TestPanelHostLatencyPercentilesAC2`、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`、`TestLiveWasapiSmoke`（真机/延迟那一族；票面 `AC#5` 本来就⛔ 把它们算进无窗射程）。

**6 枚红逐名**：`TestPanelHostRealWindowHopAndLifecycle`、`TestAC4FocusReturnToPriorWindowGap33r5`、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`、`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`（以上 `cmd/wisp`）＋`TestC21TableColourRowsMatchTokensCSS`（**`internal/ball`**）。

## ★判语（只说这一发撑得到的事）

1. ★**这一发与实现者"改后"那一发逐名同号**：我把它 `293-r1` 的 `logs/post-run-3.txt`／`post-run-4.txt` 名册拉出来对 ⇒ **6 枚同名**、`--- PASS` 同为 **507**。⇒ 票 293 的"新增红 0"从今天起**不再只有腿的数**。
2. ⚠**最后那一枚红的根因⛔ 是本票、⛔ 是产码**，红句逐字（我件 `:2432`）＝
   `tokens_table_test.go:1468: read design/assets/tokens.css: open D:\work\workspace\projects plans\Wisp\design\assets\tokens.css: The system cannot find the path specified. - the CSS leg of this check must never skip`
   ⇒ 工作树那 **805 枚脏**里 `design/assets/tokens.css` 被就地删了（`git status` 显 ` D`）。⛔ 我造成（跑完我复跑 `git status --porcelain -- cmd internal`＝**0 行**）。
   ⚠ 但"盘上少那枚文件就会红"这件事是**另一本账**（C21 族仪器＋项目口径里"token 真相源已换到 `design/doubao/demo/styles.css`"那条），下一轮具名登记；本件⛔ 下结论。
3. ⚠**单发红名册⛔ 是可靠尺**这一条又添一枚现量：`TestAC4FocusReturnToPriorWindowGap33r5` 在我这一发**红**，而 `293-v1` 在同一枚导出树两发之间看见它翻面（`--- PASS` 684↔683）⇒ 这就是我在 `A801` 里登记成**一族两枚**（它＋`Test197Subagent…`）的第二枚的那枚现量。
4. ⛔ **本件证不到的**：改前对照（那把尺需要成对导出树，`293-v1` 已交：改前 677／26 红、改后 684／26 红、对称差两向空——**那是它的数，⛔ 我的**）；`AC#5` 那枚肉眼可见的勾（今天零真窗）。
