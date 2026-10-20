# 351-a1 件03 = `AC#2` 确定性形状（★只给形状、本腿⛔ 动一字节）＋`AC#3` 触发判定

尺＝件01 的名册＋件02 的定性；"现成那把尺在哪"＝`git cat-file blob HEAD:cmd/wisp/config_reload_223_test.go` 行号。
★本节每一枚都是**给落地腿的判据形状**，⛔ 我立落地形状（票面第一趟只普查）。

## 一、要补的那一味（由件02 的机制倒推，⛔ 新造需求）

件02 定到的形状＝**ledger 立面先可见、stdout 立面后落字**（产码顺序由 `run.go:1423→:1431→:1441` 保证）。
所以确定性⛔ 是"把 `:544` 那句放宽"，而＝**让每一格读它自己那一面、并且先等那一面到齐**。三格形状：

| 格 | 判据形状（写死到行） | 现成的尺在哪 | 需⛔ 需产码接缝 |
|---|---|---|---|
| **G1 settle** | plant **前**：`awaitAudit(t,"config: HOT-RELOAD state=armed")` ＋ ★基线分子 `before := r.rt.windowCount()`（须＝0）＋ ★操作员侧横幅也等到（`awaitStdout(t,"配置热加载已接管")`） | 三味全在 `cmd/wisp` 本文件：`:117 awaitAudit`／`:148 awaitStdout`／`:929`(run.go)`windowCount`；先例＝`:340`＋`:344`（同族 `:393`、`:753`） | ★⛔ 需（⛔ 新导出名、⛔ 新 seam） |
| **G2 通知在场** | `card := r.awaitCard(t, configReloadTool)`（⛔ 变）→ ★断 `r.rt.windowCount()-before == 1`（**分子**）→ ★再 `out := r.awaitStdout(t, "[确认 L2 "+configReloadTool+"]")`，即"卡片那一行确实进了 stdout"单独成一格，红句自述"没渲染／没等到" | `:252 awaitCard`＋`run.go:929`＋`:148`；先例＝分子差值形 `ticket224_assembly_test.go:270/:272`、在场⇄点名两格并列形＝**同文件 `:402`⇄`:405`**（★本票射程内已有先例，⛔ 新形状） | ★⛔ 需 |
| **G3 文案点名** | 在 **G2 那条 await 返回的窗口**上断 `strings.Contains(out, "risk.permission_mode")`；★串⛔ 变、★`t.Errorf`⇄`Fatalf` 取舍由落地腿按票面裁（本腿只⛔ 允许把串删掉或改成"有句话就行"） | 就是 `:544` 那枚断言本体，只把它的**读源**从 `:543` 的单读换成 G2 的 await 返回值 | ★⛔ 需 |

⇒ **三格各自可失败**＝G2 红而 G3 绿 ＝ "卡片到了、文案⛔ 点名"（那才叫真缺陷）；G3 红而 G2 绿 ＝ "文案缺那枚串"；
G2 红而 G3 未取 ＝件02 量到的那一形。**本票那枚红落在最后这一类**＝现在的仪器分不出这三类。
★⛔ 用 `time.Sleep` 哄（票面禁区逐字写着）；★⛔ 把 `:544` 换成"stdout 里有任何一句人话"。

## 二、落哪一枚文件（★只指位置、⛔ 动手）

★**⛔ 新增文件**＝就在 `cmd/wisp/config_reload_223_test.go` 的 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`
（尺 A＝`:535`，体＝`:536-556`）里补 G1/G2/G3；同族那枚 `:401`（甲1）＝**同一失同步形状**、
`config_receipt_255_test.go` 的 `:633/:701/:743/:797` **已经是**改好的形状＝落地腿照搬现成尺、⛔ 新造 helper。
⚠ 同包写腿互斥（票面排程段自己写着：与票 227 `AC#6` 续欠那枚同占 `cmd/wisp` 测试面）⇒ 编排者派单时只放一枚。

★**射程提醒（件01 丁族现量）**：裸种的 plant 有 **4** 枚（`:466`/`:506`/`:541`/`:599`）。
G1 那条 settle **⛔ 只是本枚的补丁**＝它是全族的形状问题；落地腿若顺手补别枚＝撞票 227 那一族的格子，
⇒ 建议编排者裁"本枚补 G1⇄G3，其余三枚另格"，⛔ 我裁（本腿只具名形状）。

## 三、`AC#3` 触发判定＝**未触发**（件02 判语⛔ 真缺陷）

⇒ ⛔ 续体＝⛔ 动 `cmd/wisp/config_reload.go` 的卡片文案（那是第二道活、⛔ 本腿本票一字节）。
★但把"若后续裁腿改判真缺陷"那一支的形状先具名写死，免得再普查一趟：

- 卡片**现在**就带被改 key 名，那一跳＝`config_reload.go:283-292` 的 `Decision.Reason`（`fmt.Sprintf(...会放宽 [%s]：%v..., section, keys)`）
  → `internal/agent/approval/gate.go:606` `Reason: d.Reason` → `ui.go:35` → `run.go:1441` 打印；消费＝人眼＋`panel_pump_test.go:232`⇄`:443` 那种整段解析。
- ★若"缺陷"指的是**ledger 侧⛔ 携带 key 名**（即 `awaitCard` 返回的那枚结构体问⛔ 出文案）＝
  `internal/agent/approval/replies.go:70` `ReplyCard` 的字段集只有
  `CorrelationID/TaskID/Tool/Level/Grant/Paths/Window/Deadline`＝**⛔ `Reason`、⛔ `Keys`**；
  `config_reload_223_test.go:399-400` 的注释已经具名记录了这一点（「the ledger's ReplyCard has no reason field,
  so the naming is asserted off the console」）⇒ 那一支**需新产码接缝**（新导出字段⇄改 `Prompt→ReplyCard` 的搬运形状）
  ＝⛔ 本票⛔ 本腿射程，须机主／编排者批"第二道接缝"后再派（先例＝票 342 的 `342-r2` 那一枚具名上报）。
- ★判据正控（落地腿该带的那一发突变，⛔ 现在跑）＝在**仓外导出树**把 `config_reload.go:287` 的 `%v`（keys 那一味）删掉：
  按第三节 G3 的形状应**只有 G3 红、G2 绿**；若两格一起红＝仪器还没拆干净。成对第二发＝删掉 `run.go:1441` 整行
  ⇒ 应**G2⇄G3 都红**（证明 G2 真的在看卡片、⛔ 只看横幅）。牙的验证按公共块：落刀前量锚命中、落刀后验新锚＋旧锚归零。
