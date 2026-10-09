# 票 288 — 票 280 裁完归属之后，那一处"两把尺同声、2 行、已在 CI 名册里"的对齐**成了一张空账**

**立票**：2026-10-09 11:5x 编排者（机主令「及时补票」；来路＝票 280 自己的 `Status` 句：「`panel_inbound_guards_35r3_test.go` 那一处对齐「该改」这件事现在是一张空账……改它要么随下一枚真动 `cmd/wisp/**` 的腿捎带、要么另立落地票」——**本票就是那枚"另立"**）
**性质**：一行注释对齐。**低利害、完全可逆、零功能变化**——⛔ 不上机主清单，归工程内务。

## 现量（引用前先重跑；数字是快照）

- 件＝`cmd/wisp/panel_inbound_guards_35r3_test.go`，差异＝**1 个 hunk／2 行**，形状＝`const` 组里尾随 `//` 注释的**对齐**。
- **两把尺同声**（票 280 `AC#1` 的编排者独立复跑）：裸 `gofmt -d` 14 行 ＝ `gofumpt -d` 14 行，同一个 hunk `@@ -63,8`。⇒ 它不属于"gofumpt 独有加严"那一族，改它**同时** satisfy 两把尺。
- 对照（同批另一枚，⛔ 不在本票射程）：`cmd/wisp/panel_transport_35r2_test.go` ＝ `gofmt` 14 ／ `gofumpt` **101**（gofumpt 独有并组／每行一枚 composite literal／顶层声明补空行三类）⇒ 票 280 判它"留给下一枚真动该文件的票"，**本票不搭这一枚的车**（只挑 gofmt 那 1 处改＝CI 照旧红、要改两遍）。
- ⚠ **CI 那扇门今天已因 `.scratch/**` 里的故意坏样本常红**（票 111／269／276 三处既有账），所以"会进 CI 名册"⛔ 读不成"拦得住"——本票改的是**债本身**，不是让门变绿。

## 要建什么

- [ ] **AC#1 只改那一处对齐**：写面**只许** `cmd/wisp/panel_inbound_guards_35r3_test.go`；尺＝改前后各跑 `gofmt -d <该件>` 两列并排（前＝14 行、后＝**0 行**），且 `gofumpt -d <该件>` 也降到 0（它本就是 14＝同声）。⛔ `-w` 全仓禁用、⛔ 不许批量重排、⛔ 不许碰第二枚文件。
- [ ] **AC#2 证明没伤到测试**：`go test ./cmd/wisp/ -count=1` 改前改后各一次，报 **PASS/FAIL/SKIP 三数＋红名册差集**（⚠ 本机跑 `cmd/wisp` 需要 sherpa/onnxruntime 原生 DLL 的 harness：`GOMODCACHE` 三枚 DLL 拷到测试 exe 旁＋CWD＝包目录；**不带 DLL 时那枚 0xc0000135 是环境红、不是被验物**，基线一律具名写清用了哪支）。⛔ 不许因为"整包本来就红"就不跑改前基线。
- [ ] **AC#3 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／`gofumpt -l` 对动过的件空；`git show --stat` 名册**只含该件**；`frontend/**`／`design/**`／三枚冻结件／`testdata/**`／golden／`thresholds.go`／`allowlist.txt` 零字节；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许顺手改 `panel_transport_35r2_test.go`（那是上面具名"留给下一枚真动它的票"那一枚）；⛔ 不许顺手改任何注释文字／`ci.yml` 的 lint 步／`.scratch/**` 里的坏样本；⛔ 不许把"让 CI 那步变绿"当成本票的判据（今天它拦不住）。

**Status:** **未开工**。排程＝随下一枚真动 `cmd/wisp/**` 的写腿**同批捎带**最省（单独一发要独占 Go 面、不值）；若没有那样的腿在飞，本票可以单独排一小发。⛔ 零翻框、零 push。
