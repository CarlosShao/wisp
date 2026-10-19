# 票 351 — `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` 同一份代码三处两绿一红：卡片该不该点名 `risk.permission_mode` 这件事⛔ 有确定性仪器

**立票**：2026-10-19 22:4x 编排者（来路＝`227-v2` 的Ⓓ 段与 `346-a1` 的 (iii) 丙，两条独立；台账 `A1054`/`A1056`/`A1057`）
**性质**：★**仪器缺确定性形状**（⛔ 已知产品缺陷定案）。今天盘上量到的是"同一份代码在三个地方给出两种颜色"，而**判它到底哪一种**所需的那一味⛔ 在场。
**为什么⛔ 塞进票 346**：判据族⛔ 同——票 346 那两枚靠"抢并发窗口"到达，本枚靠 1s tick 与"卡片文案里有没有那枚 key 名"到达；`346-a1` 现量＝未结案票池 grep 这个测试名 **0 命中**＝至今⛔ 人具名它。

## 现量（每条带尺；⛔ 引用前先重跑）

- ★**三处两绿一红**（尺＝`gh run view <ID> --json jobs` 取数字 `databaseId` 后 `--log-failed --job <ID>`；★⛔ 喂作业名＝静默 0 字节＝"看起来⛔ 红"，`A1054`②）：
  - CI `39071959151`／`39101060753`／`39124122121` ⇒ 该枚 **PASS**；
  - CI `39244425666`（sha `6f808a83`）⇒ 该枚 **FAIL (2.22s)**，红句逐字 `config_reload_223_test.go:545: the card does not name risk.permission_mode: wisp run: 答复监听已接入（卡片上给了编号）。…`（字节区间 `win.norm.log` L1829..1844＝B321082..B324162，`346-a1` 现量留档）；
  - 编排者本机整包（`orch227-b/full.log`，19:2x）⇒ 该枚 `--- PASS: … (2.20s)`、红名册精确 grep 0 命中（`A1056`① 我逐字复跑）。
- ★**⛔ 一行代码为它变过**（尺＝`git log --format='%h %s' a5080612..6f808a83 -- cmd/wisp` ⇒ **0 枚**；＋两枚锚 `config_reload_223_test.go` blob 同串 `f1cac524…`）⇒ 红⛔ 来自我方改动；`A1054`② 与 `346-a1` ② 两条独立同向。
- ★**⚠ 我那发本机整包本身⛔ 是整包终态名册**（尺＝`go test -list '.*' ./cmd/wisp/`＝编译名册 333 ⇄ 跑到的顶层名作差＝未跑到 179；`A1056`①）⇒ 本机那一发只能证"跑到时绿"，⛔ 当"整包无红"用；本枚那次恰好跑到了（有结果行），所以这条读数仍可用。
- **断言本体现读**（尺＝`git show HEAD:cmd/wisp/config_reload_223_test.go | sed -n '536,552p'`）：`r.plant(t, "[risk]\n", "[risk]\npermission_mode = \"auto_approve\"\n")` → `card := r.awaitCard(t, configReloadTool)` → `if !strings.Contains(shown, "risk.permission_mode") { t.Errorf(...) }` ⇒ ★它用的是 `t.Errorf`（⛔ `Fatalf`）：这一格红**⛔ 中止**后续 `answer`/`awaitLive` 那几步。
- **卡片文案由谁产**（尺＝`git grep -n '配置热加载已接管' HEAD -- cmd/wisp` ⇒ `cmd/wisp/config_reload.go:129` 那一条 `slog`/notice 文本）⇒ 现量到的那条 stdout 里**有**热加载接管那句话、**没有** `risk.permission_mode` 这枚 key 名＝"卡片上给了编号"那句在场。

## 要判什么（第一趟⛔ 改任何字节）

- [ ] **AC#0 射程普查（只读）**：把"`awaitCard` 之后**只断言卡片文案含某串**"的仪器全数出来（尺＝`git grep -n 'awaitCard' HEAD -- cmd/wisp internal` ⇄ 每枚命中点上下 ±10 行现读它断言什么），逐枚标〔含 key 名断言？／是 `Errorf` 还是 `Fatalf`？／有无到达数分子？〕。★名册必带尺名与射程（含不含 `_test.go`＝本票射程**只含** `_test.go`），⛔ 报裸枚数。
- [ ] **AC#1 定性这一发**＝〔真缺陷（卡片在该路径下确实⛔ 点名被改的 key）／仪器形状（红与绿的差别由**运行学时序**决定，例如 plant 与 awaitCard 之间撞上另一拍 tick、读到的通知⛔ 是这一枚）／CI 环境形〕三选一，**每支都要盘上凭据**＋★"该形状到没到达"的分子（⛔ 拿"总通知数很多"当到达＝本轮实测反例：红发的 `refused` 总数比绿发更大）。
- [ ] **AC#2 给确定性形状（只给形状、⛔ 动手）**：若 AC#1 判"仪器形状"⇒写出该新增/改写的那一枚常驻件的判据（例：plant 前先 settle 到已知集、await 到的卡片必带相关 id／key 名、断言拆成"通知在场"＋"文案点名"两格各自可失败），★⛔ 放宽任何既有断言（⛔ 把 `Contains(key)` 改成"有句话就行"＝把问题扫到地毯下）。
- [ ] **AC#3 若判真缺陷**⇒具名续体形状（卡片该在哪一跳带上被改 key 名＋谁消费），⛔ 本腿⛔ 本票改产码（那是第二道活，归编排者裁后另派）。
- [ ] **AC#4 纪律格**：⛔ 翻任何框（本票四格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`、⛔ push；★台面三味硬要求照公共块（反斜杠 `GOTMPDIR`、整条 `PATH` 一次 `cygpath -wp`、判据只认到达数）。

## 禁区

- ⛔ 动机主那份 `config.toml`（`C:/Users/swq/AppData/Roaming/wisp-dev/config.toml`）＝夹具⛔ 是它；本票要的是"仪器能不能确定性地问到那件事"。
- ⛔ 为变绿删／弱那条 `Contains(shown, "risk.permission_mode")`；⛔ 加 `time.Sleep` 哄时序；⛔ 把它标成"已知摆动"就结案（`A1041` 那族的病＝弃权形仪器被当已钉住）。
- ⛔ 顺手改 `cmd/wisp/config_reload.go` 的卡片文案（那是 AC#3 之后的第二道活）；⛔ 与票 227 `AC#6`／票 346 那批⛔ 混批（同包写腿⛔ 并发）。

## 排程

`351-a1`（只读＝AC#0＋AC#1，凭盘上⇄CI 日志，⛔ 跑任何包）→ 编排者按它的定性裁 AC#2／AC#3 谁先 → 落地腿（只动测试面）→ 非实现者验收翻勾。⚠ 与票 227 那一族（`AC#6` 续欠的"226 写回⛔ Adopt 内存"确定性用例）**同占 `cmd/wisp` 测试面**⇒ 同包只一枚写腿。

**Status:** **未开工**。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-19 22:4x +08] agent=编排者 did=立票 351（来路＝`227-v2` Ⓓ 段＋`346-a1` (iii) 丙；三处两绿一红、a5080612..6f808a83 对 cmd/wisp 0 枚改动、断言用 Errorf 非 Fatalf 三味都现量）next=派 351-a1（只读普查）
