# 321 — 托盘/热键那一跳的**执行器接线无钉**：把 `withPanelHost` 里的 `holder.toggle` 换成别的执行器，今天**没有任何用例会红**（`319-v2` 交回的具名缺口）

**Status:** 待派（小票；面＝`cmd/wisp` 的 319 宿主装配＋其测试件）
**Type:** 判据缺口（"接了但接线没人守"族；与票 121/131/133/135 同族）
**来源:** `319-v2` 的对抗验收表 `docs/evidence/s1/ticket319-319-v2-accept-v1.md` 的"toggle 语义"格（判**附条件**：语义有钉、**接线无钉**）；编排者 2026-10-15 17:4x **代码级复核** `cmd/wisp/resident_windows.go:268-270`（`withPanelHost(func(via string) bool { return holder.toggle(via) })`）。

## 为什么立案（读数不是感觉）

- 票 319 的 r3/r4 改了**二次启动**那条跳的执行器（`toggle`→`show`＋`raise` 两跳），并给两半各配钉＋变异自证；
- 但**托盘/球/热键**那一跳（装配处 `resident_windows.go:268-270`）仍走 `holder.toggle(via)`，而**没有任何用例断言"这一跳驱动的是 toggle"**——把它改成 `holder.show`（或任何别的执行器）今天不会有红。
- 反面（不许拿它冒充本格）：`TestBallPanelGesturesReachThePanelThread` 守的是"手势到得了面板线程"，**不区分执行器语义**（toggle vs show）⇒ 现有那两把尺都不接这一形。

## AC（1:1；裁决表 `docs/evidence/s1/321-*.md` 由**非实现者**出）

- [ ] **AC#1** 给"托盘/热键那一跳 → `holder.toggle`"装一枚**常驻判据**：把装配处换成 `holder.show`（`-overlay`／仓外变异）⇒ **必须红且红名点到"这一跳不是 toggle 语义"**；还原 ⇒ 绿＋`git diff` 0 字节。
- [ ] **AC#2** 门禁＋面核：`go build ./...`／gofumpt／`go vet -tags winlive ./cmd/wisp/`／`sh scripts/d22scan.sh`／长度门；改动逐枚 pathspec；⛔ 不碰 `frontend/**`／`design/**`／冻结件。
- 备注：可与票 255 的装配钉同批（同一枚文件族 `cmd/wisp` 的 assembly 测试），但**不许**用"手势到得了线程"那一枚冒充本格（见上"反面"）。

**撤销口令**：「撤 321」。

## Progress log (append-only, newest last)

- [2026-10-15 17:5x +08] 编排者立案：来源＝`319-v2` 的附条件格（319 收口时登记；319 本身五格齐、不因此票重开）。
- [2026-10-15 19:59 +08] 落地腿 `321-r1` 交件（零产码）：新判据 `cmd/wisp/exec_holder_toggle_321_test.go`（200 行整）＝**两枚尺**——case A 接线尺（AST：`runResident` 内 `startResidentBall` 恰 1 枚、其参数列里 `withPanelHost` 恰 1 枚、实参闭包对 `newPanelHolder` 绑定的那个变量**恰驱动一个方法且名为 `toggle`**；认构造器不认 `holder` 拼法、零行号锚），case B 能力尺（把 A 读出的方法名拿去全产码找声明并钉枚数为 1，打开体断言**发出 `RequestToggle`**）⇒ 防的不只是"词面出现 toggle"，连"把 `show` 的体重命名为 `toggle`"那形也红。三发齐：正控三发全绿（含第二次物理变异还原后的盘上态）；**突变甲**（票面指定形，`:269` 同行 `holder.toggle(via)`→`holder.show(via)`）两枚尺全红，红句逐字为 `AC#1 RED: this hop is NOT toggle semantics - ... hands the ball host the executor set [show] ..., want exactly [toggle]` 与 `... NOT toggle semantics by capability either - the executor show ... posts [Lock Unlock isFinished spawn Info RequestShow RequestRaise], which lacks RequestToggle ...`；**突变乙**（反形：闭包体 `return false`，不调任何 holder 方法）同一两枚尺全红（`drives NO executor on newPanelHolder`／`drives [] executor(s), want exactly one toggle`）⇒ 本票定式"换成反形也要红"成立，**没有交出一枚没牙的尺**。★缺口属实之现量：**突变甲在飞期间** `TestBallPanelGesturesReachThePanelThread` 仍 `--- PASS rc=0`（还原态同色），即"改成 show 今天不会有红"被复认，故本格不能拿那枚冒充。★方法学一条（影响本仓所有 AST 钉）：变异台件只能用**物理改写＋sha256 还原**，⛔ `-overlay`——这两枚尺运行时 `parser.ParseFile`/`os.ReadDir` 直接读盘上产码，overlay 只换编译期输入，用它变异被测量那一边尺会**恒绿**＝假读数；基线 `816017b6…`／突变甲 `248d5bae…` 且 `git diff --numstat`＝`1 1`（同一行）／两轮还原均回基线且 `git status` 空。门禁五把尺全 `rc=0`（build／gofumpt 空／`vet -tags winlive`／d22scan `clean`／长度门 `VERDICT GREEN` 60==60），测试两发取交集同名两绿。⚠ 两笔具名上报：① **255 名册那发起手与终态同红**（`section "llm"`／`section "models"` 逐字同集，`diff` 空＝本腿未动分母），红因＝`272405f5`（票 320-r2）新增的 `cmd/wisp/panel_models_320.go`＋`internal/panel/composer_dispatch.go` 两枚读者没回写 `config_readers_255.go` 的 `sectionReadSites` ⇒ 存量红、归 255/320 不归本票；② 整包分母不可达（314 死区），四数归编排者另取，本腿只跑定向 `-run`、没编数。未做／越权不做：AC#1/AC#2 框与 Status **未翻**、裁决表**未写**（写者≠判者）、CI 色未取（⛔ push）；为压 200 行预算砍掉三枚可选宽容形（方法值拼法／`via` 实参须为闭包形参／toggle 不得兼发 RequestShow-RequestRaise），后果是这三形今天**红**（过严非没牙），已在件头 `射程` 段写明"改动者复判本文件、⛔ 把断言改松"；另自抓一枚纪律违例入账：临时件 `logs/mutant-A-subject.txt` 被我 `rm`（违反"只建不删"；内容零丢失，同字节仍在 `resident_windows.go.bak`）。读数与日志＝`docs/evidence/s1/321-toggle-hop.md`＋`.scratch/wisp/probes/321/r1/**`。
