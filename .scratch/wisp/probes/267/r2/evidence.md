# 票 267 落地腿 `267-r2` — 证据件（种子迁移：`cmd/wisp` 带外 `confirm_timeout_sec` 抬进带内）

任务：票 267 的 ⓐ 值域（`[31, 3600]`，`internal/config/validate.go`，已进 HEAD）把 `cmd/wisp` 一批集成用例的种子打到带外 ⇒ 装配根拒载 ⇒ 16 枚 FAIL。本腿执行编排者 §7.1 裁的出口**甲**＝把带外种子抬进带内，并顺手拿票面 AC#2 那半格（C18"即将超时"提示的正向读数）。
具名解冻：台账 `A611` §1，撤销口令「267 撤」。

## §0 起手锚

取数时刻／分支／自取 HEAD／并发腿／本腿写面归属的 `git status --porcelain` 读数。

## §1 名册自取尺

本腿自己拉的尺（根＝`cmd internal tools docs .scratch scripts`）数出的带外 `confirm_timeout_sec` 喂值点名册，以及与编排者派单那 12 枚的差集逐枚说明。

## §2 逐枚改动

每一枚种子的 before/after 逐字，带 `file:line` 与取值理由（31 或 40）。

## §3 配套件：`panel_pump_test.go` 的 ctx 时长

种子从 2 s 抬到 31 s 之后必须连带抬的那一枚 `context.WithTimeout`，以及为什么它是配套件而非优化。

## §4 ★C18"即将超时"提示正向读数的落点与突变自证

现读的链（`gate.go:528` 的 `warn` → `g.ui.Update` → 消费侧可读面）、新增断言的逐字凭据、两发突变的红句，或"够不着"的具名归口。

## §5 门禁五把尺读数

整包 `cmd/wisp`／`internal/config`／`gofumpt`／`d22scan`／`path-length-budget`／`go vet`，逐把带取数时刻。

## §6 红名作差

16 枚逐名判决（只减不增），任何新增红的红句逐字并标〔未归因，归编排者〕。

## §7 判不动／量不到

具名归口，⛔ 不写"应该没问题"。

## §8 污染面与提交名册

逐笔 `git show --name-only`＋票面 AC 框零改动自证＋写面越界尺。

## §9 我写错的读数

自我对抗：名册尺／观察对象搬家／行号漂移／注释与字面量不一致／突变是否真有牙。

## §10 交件判语
