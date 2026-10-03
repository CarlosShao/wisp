# 145-c2 快照字段普查（只读腿 · 零 Go 命令）

> 派单：编排者直令（无独立 dispatch 文件，代号 `145-c2`）。工单＝`.scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md`。
> 本程只读源码＋grep/find＋git log/show/diff，**一枚 go 命令未跑**；写面只到 `.scratch/wisp/probes/145/c2/**`；禁面一字未动。

## §0 起手锚

- 取数命令（同一发 shell）：`date "+%Y-%m-%d %H:%M:%S%z"` ＋ `git log -1 --format="%h %ci"` ＋ `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`
- 读数（2026-10-03）：
  - `date` = `2026-10-03 09:26:36+0800`
  - HEAD = `5f9ff9d4 2026-10-03 09:25:59 +0800`
  - 脏项计数 = `365`
- 相关目录脏项（本轮 `git status --porcelain -- internal/panel/ cmd/wisp/ internal/config/ internal/agent/approval/ internal/tools/ internal/risk/` 现量）：
  - ` M cmd/wisp/config_reload.go`
  - ` M internal/tools/paths.go`
  - `?? cmd/wisp/config_readers_255.go`
  - `?? internal/config/tiers_app_255r2_test.go`
  - ⇒ **`internal/panel/**`（composer.go/pump.go/approval.go/子件）与 `cmd/wisp/run.go`、`cmd/wisp/panel_pump.go` 此刻均不在脏项里** ⇒ 本程读到的这些文件＝HEAD `5f9ff9d4` 已提交态，非他人未提交中间态。⚠ 在飞的三枚写腿（config/approval/tools+risk）可能随时改脏，复量者请以自己现量为准。

## §1 快照字段名册（三种"没做"分开标）

〔待填〕

## §2 零值冒充真值的名册＋"未知 vs 零"的形状候选

〔待填〕

## §3 状态↔缺字段对照表（十四枚逐行）

〔待填〕

## §4 复认/推翻前人件（具名）

〔待填〕

## §5 我可能写错的条目（自我对抗）

〔待填〕

## §6 量不到的地方（具名）

〔待填〕

## §7 交件判语

〔待填〕
