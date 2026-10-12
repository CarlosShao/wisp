
## 10-12 一条：仓里有一枚**按行号取证的既有仪器**，任何"注释加高"都会把它打红 ⇒ 派 `cmd/wisp` 腿之前先现跑那把尺

尺（10-12 现跑，逐字可重放）＝`git grep -ho -E '(cmd|internal)/[a-z_/0-9]+\.go:[0-9]+' HEAD -- cmd/wisp/config_readers_255.go` ⇒ **26 枚行号引用、跨 11 枚文件**（`cmd/wisp/run.go` 6 枚／`cmd/wisp/resident_audio_windows.go` 4 枚／`cmd/wisp/models.go` 2 枚／其余各 1 枚，**含 `internal/ball/ball_windows.go`、`internal/audio/gate.go`、`internal/observe/logging.go`**）。咬人的用例＝`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`（`cmd/wisp/config_receipt_255_test.go`），红句逐字含 **`the evidence drifted, so re-adjudicate this row`**。

实测那一形：落地腿 `295-r1` **只是给两组注释加高**（净 +19 行、语义零变）就把那枚用例打成"新增红"；它的处置形＝**把自己的注释压回原行数（净零增行）＋⛔ 动那两枚 255 的文件＋⛔ 让回执迁就漂移**，并交回"行号复归"的读数 ⇒ 我（编排者）复认＝形对，那条"⛔ 为变绿放宽断言"它守住了。

**Why:** 我原先那条"⛔ 行号锚的突变台件会腐烂"讲的是**我自己新写的**台件；这一枚是**仓里既有**的、而且行号写在**非测试件**里（一张"配置回执"表），所以：①任何写腿往那 11 枚文件里加注释／加行都会撞红；②它红的时候**看起来像产码有缺陷**，完全不像"某人把注释写高了"——我很容易把这种红记成产码账、或让腿去"修"那枚本来正确的回执。
**How to apply:** ①派单写面凡含 `cmd/wisp/**` 或那几枚 `internal/**` ⇒ **起手先跑上面那把尺**，被引文件逐枚具名写进派单："往这枚文件加行＝会红那一枚用例；处置＝只改同行内容、或把增量压回原行数；⛔ 改那两枚 255 的文件、⛔ 让回执迁就漂移"。②我自己收腿／代笔时，见到"新增红 1 枚"先查红句有没有那句 `the evidence drifted`——有 ⇒ 归因到**行数变化**，⛔ 记成产码缺陷。③与"每波只交一枚 `./cmd/wisp/` 色"叠用：两枚腿先后碰那 11 枚文件时，后落地的那枚会把先落地那枚的锚**再打死一次**，所以串行不只是编译面争用，也是**锚位争用**。
