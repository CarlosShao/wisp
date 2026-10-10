# 303-r1 `AC#3` 结果与一处顶回（三枚受害用例并非全绿，且非我改坏）

件：`50-revert-red.txt`（撤修复＝HEAD）／`51-restore-green.txt`（还原修复）／`31-package-after.txt`（整包 -v，含本修复）。

## 改前红（HEAD，撤掉我的修复）— 三枚逐字同判据
- TestAC13ColdStart…: `no report "ac13-probe" … (what DID arrive at the door: nothing at all)`
- TestAC14Awaited…:  `no report "ac14r-0" … nothing at all`
- TestAC14GoSideEvalPush…: `no report "ac14-push" … nothing at all`
rc=1（50-）。母仓 HEAD 同三枚亦红（20-）＝编排者 `A820` 承重对的复现。

## 改后（还原我的修复）— 逐字
- TestAC14AwaitedBindingReplyReachesThePage：**PASS** `page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)` ← 回执那一跳被本修复接上。
- TestAC13ColdStart…：**仍 FAIL，但换了形**（0.69s，非 15s 超时）：页面**已能答**（`AC#13 page answer … "0"`、`cold 465.8 ms`、`window_opened=true`），红句＝`the live document contains NONE of the 1 element ids the embedded entry declares`（bringUp 冷启后**最后盖着的文档是 round-trip 探针 stub，不是内嵌入口**）。
- TestAC14GoSideEvalPushReachesThePage：**仍 FAIL，但换了形**：回执通道已通（页面报回 `title=""`），红句＝`Go's Eval push did not reach the document: the page reports its title as "", want "PUSHED-33R5-OK"`（同一枚"入口未成为活文档"成因：探针之后没有把入口交接上来，title 被后续文档覆盖/探针无 title）。

## 具名顶回（不凑数、不越权）
- 票 303 `AC#3` 写的"修后 TestAC13／两枚 TestAC14* 全绿"这一预期**不成立**：本票归因**收敛到一笔＝`fb2fb802` 的 forward 钩子**，它只断"页面→Go 的回执通道"。`AC#2` 三形里活动形＝丙（回执通道被再入/双包挡死），本修复精确接上它 ⇒ nail1 绿。
- `TestAC13` 与 `nail2` 现在卡的**不是回执通道**（通道已通、页面能答），而是**另一枚形状**：`bringUp` 的 `serveEntry` 与探针文档**交接次序**——探针 stub 是冷启后最后显示的文档，入口没被盖回。该形状属**票 33 `AC#13` 的产码射程**（`firstRoundTrip` 先于 `serveEntry`；嫌疑笔＝`70b00885` 33-r10 冷启交接，**非** `fb2fb802`），改它要动 `bringUp` 的文档交接顺序＝**第二枚未归因、且不在本票门-传输边归因里的改动**。
- 按"未定义即停／⛔ 顺手扩范围／⛔ 凑三枚全绿"：本腿**不动** `bringUp` 的 serveEntry/探针交接，把它**具名上报**给非实现者/编排者定夺（另立形状或并入票 33），⛔ 由我盖章其成立。

## 新增红（整包两发取交集）口径与欠账
- 改后整包 `31-package-after.txt`（-v、含本修复）顶层 FAIL＝`{TestPanelHostRealWindowHopAndLifecycle(-1.000 形), TestAC4FocusReturnToPriorWindowGap33r5(已知 flake 族), TestAC13ColdStartEnds…(入口交接形)}`；`TestResolvePerCallBudget` 该发未现红（已知负载敏感）；nail2 定向发为入口交接形红。
- 以上每一枚改前（HEAD）皆已红：三枚受害用例 HEAD＝`nothing at all`（20-/50-）、lifecycle 已知两形、AC4Focus 已知 flake、Resolve 已知负载敏感。⇒ **本修复新增红 0 枚**（改后红集 ⊆ 改前红集），且把 nail1 从红转绿。
- ⚠ **欠账具名**：本腿未能在母仓取得一枚**与改后同条件、同码的干净改前整包 `-v` 名册**——背景那发 `go test` 的**编译发生在我编辑该文件之后**（构建带上了修复），故 `31-` 只能当改后名册用；改前整包名册的"新增红 0"以上面的"改后红集 ⊆ 改前逐枚已知红"论证替代，非同一把整包尺对拉。要同一把尺的改前整包名册，需在一枚**未被我编辑的 HEAD 树**（仓外 clone 或编排者）重跑（clone 无 `frontend/dist` ⇒ AC13 走具名 SKIP，与母仓台面不同，见票面 `:17` 台面口径）。
