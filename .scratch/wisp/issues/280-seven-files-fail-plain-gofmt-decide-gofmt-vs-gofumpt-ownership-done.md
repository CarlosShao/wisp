# 票 280 — 7 枚文件不过裸 `gofmt`：逐枚判归属（机器 `gofmt` vs CI 那步 `gofumpt`）

**立票**：2026-10-08 19:3x 编排者（来路＝台账 `A736` 收工门禁快照：`gofmt -l cmd internal` ＝ 7 命中）
**性质**：读数为先；**不预设要改**（改它们＝"顺手改"，且可能与 CI 那步不是同一把尺）。
**更正（2026-10-09 10:3x 编排者，题面不改字）**：标题里那枚"7 枚"是**工作树**读数；HEAD 上只有 2 枚真未格式化，另 5 枚＝CRLF artifact（现量与尺见下面 §现量 第一条）。件名 slug 不动（`-done` 名是防重领唯一键，二次改名会让引用漂），要读结论请读 §现量 与三格 ✅ 注。

## 现量（引用前重跑；⚠ 数字是快照）

- ⚠⚠ **本节"7 枚"这枚数在 HEAD 上只成立 2 枚**（2026-10-09 10:3x 编排者现量打旧；来路＝腿 `280-r1b` 顶回票面前提，编排者独立复跑同一把尺确认）：`gofmt -l` 对**工作树**给 7 枚，对**`git show HEAD:<path>` 落地成 LF 的 blob** 只给 **2 枚**（`cmd/wisp/panel_inbound_guards_35r3_test.go`、`cmd/wisp/panel_transport_35r2_test.go`，各 14 行 diff＝各 1 个 hunk）。其余 5 枚（`models.go`／`pending_read.go`／`tools.go`／`provenance.go`／`bridge.go`）在 HEAD blob 上 `gofmt -d`＝**0 行**，被点红的原因是**行尾**（现量：工作树 CR 字节＝行数 334/131/225/1115/1284；HEAD blob CR＝0；`core.autocrlf=true`；`HEAD:.gitattributes` 含 `*.go text eol=lf`）⇒ **是尺的射程错，不是文件的债**。下面那句原样保留、不删。
- 7 枚（＝工作树射程；HEAD 射程只有 2 枚，见上）：`cmd/wisp/models.go`／`panel_inbound_guards_35r3_test.go`／`panel_transport_35r2_test.go`（10-08 今天那枚）／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`internal/tools/bridge.go`。
- 来历（`A736` 逐枚 `git log -1`）：6 枚既有（10-03～10-07），1 枚今天已提交、树净。
- ⚠ CI 的 lint 步量的是 **gofumpt**，与裸 `gofmt` **不是同一把尺**（`A721` 一族）。

## 要建什么

- [x] **AC#1 逐枚判**：7 枚各给"差异属于 `gofmt` 风格还是 `gofumpt` 加严"（本机若有 gofumpt 二进制就跑，没有就具名写"判不动＋缺什么"）。
      ✅ 编排者 2026-10-09 10:3x 翻勾（凭据＝腿件 `.scratch/wisp/probes/280/r1b/10-census.md:17-53` 逐枚表＋hunk 原文；gofumpt 二进制本腿自己定位到了＝`D:\\work\\base\\gopath\\bin\\gofumpt.exe` `v0.12.0 (go1.27.1)`，**不在 PATH**，零 `go install`／零下载）。编排者独立复跑同一把尺（`git show HEAD:<path>` 落地成 LF 后逐枚 `-d` 行数）＝`inbound_guards` **gofmt 14 / gofumpt 14**（同声、同 hunk `@@ -63,8`＝const 组尾随 `//` 注释**对齐**＝纯 gofmt 风格）；`transport_35r2` **gofmt 14 / gofumpt 101**（gofumpt 独有三类：连续单行 `type` 并组 `@@ -400/-421`、composite literal 每行一枚 `@@ -698`、顶层声明间补空行 ⇒ **混合且裸 gofmt 低估它**）；其余 5 枚 **0/0**。⇒ 三形齐了：`gofmt 风格`／`gofumpt 加严`／**第三种＝两把尺都无关的行尾 artifact**（票面原来没给它名字，本程补上）。
- [x] **AC#2 与 CI 的关系**：逐枚答"它在 CI 那步会不会被点红"，⛔ 不许用"应该会"填空。
      ✅ 编排者 2026-10-09 10:3x 翻勾，但**判据按腿的更正收窄一档**（腿件 `:55-73`／`:86`）：CI 那两步（`ci.yml:172` 步名→`:173` `if: !cancelled()`→`:175` `go install mvdan.cc/gofumpt@latest`→`:176` `gofumpt -l . tools/d22scan tools/mockllm`；`:189`→`:239` 守卫→`:240` `attrib.sh --tracked-only`＝tracked 全集分母）**本程零 run 观测**，且那扇门今天已因 `.scratch/**` 里的故意坏样本常红（票 111 `:310`/`:394`、票 269、票 276 三处既有账）⇒ 所以格子里那句"会不会被点红"**不能**读成"拦得住"，只有名册口径可答：**会进 OUT／分母名册＝2 枚**（两枚 `_test.go`，tracked、住 `cmd/**`，两步都吃），**不会进＝5 枚**（HEAD blob 上 gofumpt 0 发声＋ubuntu 检出按 `.gitattributes` 给 LF ⇒ CRLF 那档在 runner 上根本不存在），**判不动＝0 枚**；唯一残余不确定＝CI 用 `@latest` 未钉版本（票 124 `:251` 早具名同一处）而本机尺是 `v0.12.0` ⇒ "会进"是**下限**。⚠ 顺带顶回编排者上一程：`280-r1` 那把"gofumpt -l 这 7 枚＝7 枚全中"吃的是**工作树字节**，换 HEAD blob 只剩 2 枚（腿件 `:85` 具名，账入 `A756`）。
- [x] **AC#3 处置建议只写形状**：改／不改／留给下一枚真正动该文件的票；⛔ 本票不落地任何格式改动。
      ✅ 编排者 2026-10-09 10:3x 翻勾（凭据＝腿件 `:75-80`，三档形都给了且**零落地**）：5 枚行尾 artifact＝**不改**（去"修"它只能改写行尾＝无的放矢）；`panel_inbound_guards_35r3_test.go`＝**改**（一处对齐、2 行、两把尺同声、已在 CI 名册里）——但**不在本票做**，本票禁区写着"不落地任何格式改动"，落地另派；`panel_transport_35r2_test.go`＝**留给下一枚真正动该文件的票**（只挑 gofmt 那 1 处改＝CI 照旧红、要改两遍）。同堵墙但不在本票射程（具名不改）：`:176` 的 scope 含 `.scratch/**`＝票 269／275／276 的既有账。

## 禁区

⛔ 不许 `-w`；⛔ 不许批量重排；⛔ 三枚冻结件一字不动；⛔ 零 push。

**Status:** **done**（2026-10-09 10:3x 编排者按 README 规则 4 收口：三格全勾，凭据见各格 ✅ 注；普查程＝腿 `280-r1b`（前程 `280-r1` 死于服务端断流、遗产已代提），翻勾者＝编排者＝非实现者，且两把尺编排者本人独立复跑逐枚对上）。⚠ 收口不等于零残余——**本票只裁归属、一字节格式都没落地**：`panel_inbound_guards_35r3_test.go` 那一处对齐「该改」这件事现在是一张空账，⛔ 不许在本票名下顺手改，改它要么随下一枚真动 `cmd/wisp/**` 的腿捎带、要么另立落地票。
