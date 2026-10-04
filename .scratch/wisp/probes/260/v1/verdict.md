# 票 260 — 非实现者验收腿 `260-v1` 对抗验收表

> 本腿＝`260-v1`，未参与 `260-r1`/`r2`/`r3`/`r4` 任何一程。本表只放本腿自己现跑的读数与判语。
> 起手时刻 `2026-10-04 13:16:10 +08`；起手 HEAD＝`515ca5c5`（`dev`，未推 81 枚）。
> ⛔ 本表不翻任何勾；票 260 的框由编排者翻。⛔ 本表零产码改动。
> 落点＝`.scratch/wisp/probes/260/v1/verdict.md`；读数档同目录 `logs/`。

## 一、起手台本与写面基线

本节记：起手 HEAD / 钟点、`git status --porcelain -- internal cmd tools scripts .github docs` 起手态、
关键六枚文件的起手 md5（`logs/baseline-md5-keyfiles.txt`）、以及本腿全程写面还原证据。

## 二、头一格：`260-r3` 自报的 M4 盲区——本腿亲手重跑那一发突变

本节记：把 `cmd/wisp/resident_approval_windows.go` 的 `residentCancelKeySpelling` 里"读球回执"那一步摘掉之后，
① 是否真的全绿、② 若红则指名用例与红句逐字、③ 若全绿则 AC#1「配置值真被读到」这格的凭据落在谁身上（具名到钉，或判"没人兜"）。

## 三、AC#1（形ⓐ：借用那遍真读 `[hotkey] cancel` 的配置值）

本节记：两档现量（默认档／配置档 `Cancel: "Ctrl+Alt+V"` 形状）——"用户看得见的那句话"各是什么；
以及票 245 那条"稳态不绑、只在确认那两三秒借裸 Esc"的裁定有没有被悄悄改掉（`internal/ball` 三枚回执函数的射程）。

## 四、AC#2（零症状那一格要有声）

本节记：判据要求的"种『借用请求发出但 Win32 没成键』⇒ 指名那一步必须红"到底做没做；
本腿自己现跑这一发种子/突变，判那枚既有尺是不是对"没成键"这一形**不敏感**（判据换成反形还全绿＝瞎尺）。

## 五、AC#3（越界检查：本腿自己拉名册）

本节记：`git show --name-only --format=` 逐枚取四程全部文件（probes 之外），逐名对 `A595`＋`A598 §2` 授权集；
禁区逐名扫描（`docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／
`frontend/**`／`design/**`／三枚冻结件／`scripts/slo-check.ps1`）；`internal/agent/approval → internal/ball` 依赖边是否为空。

## 六、AC#4（条件②文案侧，⛔ 不动逻辑）

本节记：① 默认档零漂移——那十枚句子照旧逐字念 `Esc`；② 配置档念配置里那枚键且不再出现 `Esc`；
以及 `r3` 那句"第十枚停手上报"（`approval.go` 通道名）是真是假——今天那一枚念什么、读起来对不对。

## 七、门禁读数

本节记：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh --with-self-test`、
定向 `go test`（`cmd/wisp`／`internal/agent/approval`／`internal/ball` 逐包 PASS/FAIL/SKIP 终值）、
`gofumpt`/`go vet`（只这三包范围）；红名集合逐名比先例；sherpa PATH 形状坑的本腿现量。

## 八、够格翻哪几枚框

本节记：AC#1／AC#2／AC#3／AC#4 各自"成立／不成立／量不到"的终判、
本腿建议翻哪几枚、哪几枚必须停勾、要不要带条件（判断归本腿，翻勾归编排者）。

## 九、记我 与 判不动／量不到

本节记：本腿自己哪把尺写错／跑歪（逐条留档不当笔误藏）；引任何数字前的自跑声明；
以及判不动／量不到的条目逐条具名＋归口（含 winlive 那批〔未实测〕）。
