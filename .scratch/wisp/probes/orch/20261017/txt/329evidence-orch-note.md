
---

## 8. 编排者代提注（2026-10-17 10:1x +08，⛔ 抹上面任何一字、⛔ 替它圆）

`329-r1` 于 09:58 撞 **150 轮次帽**（`Reached the maximum turn limit`，171 次工具调用、53 分钟），本节由编排者追加，把"哪些是盘上已有的、哪些是它没来得及做的"分清：

1. **★两处 `待填` 属实、我没代填**：§表第 68／69 行（改前·skip3 ⇄ 改后·skip3 的四数）与第 83 行"桌面卫生尺：待填"仍是空的。⇒ **本件的成对颜色那一格⛔ 当凭据**，交下一枚非实现者腿自己取数（台账 `A1012`）。
2. **★一处具名不实（件名不存在＝幽灵引用）**：§1 指"全文在 `probes/329/r1/logs/commits-roster.log`"，该件**盘上没有**（`ls` 只回 `ac4-roster-nails-diff.log`）。⛔ 我删这句、只在此具名登记；名册本身可由 `git show --name-only` 现跑复得，编排者已逐枚复跑＝七笔全部只落 `cmd/wisp/**`＋`.scratch/wisp/probes/329/**`（⛔ 票面/证据件在任何一笔里）。
3. **改前／改后那两发整包的真实死法（编排者自己数过，非本腿读数）**：尺＝对 `logs/full-*.log` 逐枚 `grep -c`，⛔ 与"名册"混用。
   - `full-before-anchor-feab45e3.log`＝`RUN 166／PASS_TOP 115／FAIL_TOP 2／SKIP 0`、rc=1（红名里两枚具名＝`TestAC2SealNoticeLandsInTheRunLegLogFile`、`TestCleanCheckoutBuilds_AC11`）；
   - `full-before-anchor-feab45e3-run2.log`＝`RUN 203／PASS_TOP 148／FAIL_TOP 0／SKIP 0`、rc=1；
   - `full-after-tip-1d1261e0-skipset2.log`＝`RUN 206／PASS_TOP 151／FAIL_TOP 0／SKIP 0`、rc=1；
   - ★后两发**包回 `FAIL`＋耗时（204.198s／196.850s）却零枚 `FAIL:`、零枚 `ok`**，日志尾都停在 `MARK299 wndproc …` 之后 ⇒ 形状＝**整进程在真窗区死亡**（本仓已登记的同一形状＝台账里"CI 两发早停＝票 314 的 CI 暴露"那一节；⛔ 我在此引行号，它会腐烂），⛔ 读成"329 的字节把包打红"，也⛔ 读成"全绿"。**⇒ 名册作差这条路今天取不到**，改前改后各自死在同一处。
4. **编排者自己复跑过的那把尺（唯一我背书的一把）**：`grep -rn lastPanelSnapshot --include='*.go' internal cmd tools | grep -v _test` 现量 5 行＝定义 `cmd/wisp/panel_pump.go:422/429` ＋ 注释 3 枚 ＋ **生产读者 1 枚 `cmd/wisp/panel_snapshot_deliver_329.go:95`**（起手那一把尺＝0 命中，见 `.scratch/wisp/probes/329/r1/logs/ruler-lastPanelSnapshot-before-A.log`）。⇒ **AC#2 那一跳的"有非测试读者"这半成立**；"**字节逐位相等**"那半我只在盘上数到用例在场（`cmd/wisp/panel_snapshot_deliver_329_test.go:76`＝`TestBookedSnapshotReachesHostBitForBit329`），**没复跑它的成对三发**（去牙／反控／摘读者）⇒ 归 `329-v1`。
5. **摘读者那一发（本腿 mutation-3）盘上读数逐字在**：`logs/mutation-3-reader-removal.log` 里那枚用例红在三行（`:141`／`:164`／`:167`，`host=""`、"received 0 pushes"）⇒ **它确实是"摘掉生产读者就响"的**，这是我读过原文、没自己跑的那一发。
