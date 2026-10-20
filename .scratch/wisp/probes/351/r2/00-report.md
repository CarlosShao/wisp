# 票 351 落地腿 351-r2 交件件＝三格确定性形状落地⇄AC#3 敏感度自证

★写面＝只有 `cmd/wisp/config_reload_223_test.go` 那枚**已存在的函数体**
`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（⛔ 新文件／⛔ 新 helper／⛔ 任何产码字节／⛔ 翻框／⛔ push）。

## 〇、起手锚与可比性凭据

- 本腿第一笔 `git log -1` 现认起手锚＝`bf509bb837aa1328617320608be38e4bca5a5d35`（`dev`）。
- ★改前两发起手现跑 `git rev-parse HEAD:internal/winsec/winsec.go`＝`07300bfffc3e5b14c017029d6f84fb5a7ede3db8`
  ＝两棵导出树（treeB⇄treeL）内 `winsec.go` 的 sha256 **逐字同**＝`9086f2bbbf531823ed5c4a615391785f5af749542df15c046b5fb4f2da6f2331`。
  ⇒ 我的所有读数对同一枚 winsec 字节可比、⛔ 被 `352-r1` 的 `%w` 落地污染。
- ⚠ 交件时 `HEAD:winsec` 已漂到 `f3c1b68e…`＝`352-r1` 笔②`ecffb93b` 在我起手之后落了 `internal/winsec`（**别腿的合法落地**）。
  ⇒ 我这枚红句**从未点名 `internal/winsec/**`**（全程绿或票 351 那枚运行学时序红），⛔ 我据此判自己失败、⛔ 我动它一字节；
  我的 `git show --name-only`＝1 枚文件＝只我那枚 `_test.go`（越界 0／产码 `.go` 0）。

## 一、落进去的三格形状（各自可失败）

| 格 | 落点（落地后行号） | 判据形状 | 现成的尺（⛔ 新接缝） |
|---|---|---|---|
| **G1 settle** | `:561`/`:562`/`:566` | plant **前** `awaitAudit(state=armed)` ＋ `before := windowCount()`（应 0，非 0 即 `t.Fatalf` 前提挪了）＋ `awaitStdout("配置热加载已接管")` | `:117 awaitAudit`／`:148 awaitStdout`／`run.go:929 windowCount` |
| **G2 通知在场** | `:582`/`:585` | `windowCount()-before == 1`（**分子**，`t.Errorf`）＋ `shown := awaitStdout("[确认 L2 "+configReloadTool+"]")` **独立成一格**（⛔ 与文案格合并） | `ticket224_assembly:270/:272` 差值形／同文件 `:402⇄:405` 并列形 |
| **G3 文案点名** | `:591` | 在 **G2 返回的那条流**上 `Contains(shown,"risk.permission_mode")`；★串⛔ 改、只把读源从 `:543` 单读换成 await 返回值 | 就是 `:544` 那枚断言本体换读源 |

⇒ 三格全⛔ 需产码接缝（awaitCard/windowCount/awaitStdout 都是现成尺）⇒ 没撞上"要产码才带得上 key"那形（`342-r2` 那枚真停下上报的形在本腿**未成立**）⇒ 据实落完、⛔ 停手上报触发。

## 二、成对读数（四数只从 `^=== RUN`/`^--- PASS`/`^--- FAIL`/`^--- SKIP` 取，逐名作差）

| 发 | 树 | RC | RUN | PASS_TOP | FAIL_TOP | SKIP | 顶层红名 |
|---|---|---|---|---|---|---|---|
| B1 | treeB（基线字节） | 1 | 23 | 9 | **1** | 0 | `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` |
| B2 | treeB（同字节复跑） | 0 | 23 | 10 | 0 | 0 | — |
| L1 | treeL（落地字节） | 0 | 23 | 10 | 0 | 0 | — |
| L2 | treeL（落地字节复跑） | 0 | 23 | 10 | 0 | 0 | — |

- ★**基线一红一绿＝票 351 AC#1 的运行学时序形状在本台面再现**：同一份 `bf509bb8` 基线字节，B1 撞到那扇窗红、B2 没撞绿。
  B1 红句逐字（`go test` 正文印在 `--- FAIL` 之前）＝
  `config_reload_223_test.go:545: the card does not name risk.permission_mode: wisp run: 答复监听已接入（卡片上给了编号）。…`
  ＝票面第一节具名的同一枚红（`:545`⇄同串）。⇒ 尺有牙、⛔ 恒绿。
- ★**落地两发都绿**＝G1 的 settle 把那扇"ledger 先可见、stdout 后落字"的窗关上（G3 读的是 await 到的流、⛔ 单读）。
- 逐名作差＝基线⇄落地只 `…ModeLoosening…AfterAllow` 一枚由"偶红"变"稳绿"；其余 22 枚 `=== RUN` 名集两树**逐字同**、⛔ 枚数当"我改了 N 枚测试"（公共块：枚数⛔ 是单枚字节的函数）。
