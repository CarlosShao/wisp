# 票 100 AC#0 触发门复判（2026-11-08 22:4x +08，写腿；判完保持 RESERVED）

## 判据（票 100 AC#0 逐字）

`grep -rn "exec.Command" internal/tools/ internal/plugin/` 出现**非测试**命中 **且被 bridge 调到**。
两半同时成立 ⇒ 门开；任一半不成立 ⇒ 保持 RESERVED。

## 现跑读数（双树一致：共享工作树 HEAD=17fb41455 工作树／clone --shared 导出树 HEAD=8c897dcf5）

### 第一半：grep 非测试命中 —— 成立（且比 9-21 立票时变多了）

非测试命中 **6 处 / 3 文件**（`_test.go` 全排除；导出树同读数）：

- `internal/tools/launcher_other.go:27`（xdg-open）、`:30`（open）、`:71`（直接起 target）
- `internal/tools/launcher_windows.go:113`（直接起 target）、`:116`（explorer.exe shell:AppsFolder）
- `internal/plugin/manifestslot/command.go:454`（`exec.Command(spec.Exe, spec.Argv[1:]...)`，OSCommandSpawner，
  D46 Tier-1 形状，票 52）

注意：立票时（2026-09-21）的读数「全仓非测试 exec.Command 只有 mockllm / SLO 自测 / doctor」**已过期**：
票 22 web 族 launcher 与票 52 Tier-1 CommandRuntime 两个产码族已落码。

### 第二半：「且被 bridge 调到」—— 不成立

- `BuiltinWebEntries`（含 app.launch，`internal/tools/web.go:210` 定义）：**非测试调用方 = 0**
  （全仓 grep 仅定义处与注释命中）。
- `NewCommandRuntime`（`internal/plugin/manifestslot/command.go:92` 定义）：**非测试调用方 = 0**。
- 生产注册面（`cmd/wisp/run.go`）现只注册五族：
  `BuiltinFSEntries(:659)`、`BuiltinTaskEntries(:697)`、`BuiltinSystemEntries(:728)`、
  `BuiltinSubagentEntries(:1206)`、`BuiltinPlanEntries(:1245)` —— web 族与 plugin host 均不在其中；
  `LoadManifest`/`PluginHost` 在 cmd/wisp + internal/agent + internal/plugin 非测试面 grep = 0。

## 裁决

**门未开 ⇒ 本票保持 RESERVED**（未改 Status 行、未翻任何 AC 勾，只加本证据件与 Progress log 一行）。

**门开条件已收敛为一句话**（比 9-21 时更近了）：`cmd/wisp/run.go` 注册 `BuiltinWebEntries`
（app.launch 随之可被调到）**或** plugin host 把 `CommandRuntime` 接进生产路径（任一 Tier-1 插件真被加载执行）
之日，即 AC#0 门开之日 —— 届时按票面「派单简报第一句就是先跑 AC#0 的 grep 证明门开着」派落地腿。
