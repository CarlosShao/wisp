# 票 111 — 111-r1b 证据件（骨架 · 先写未判）

代号：`111-r1b`（编排者 2026-10-05 派，替代 `111-r1` 的重派腿）
骨架落盘时刻：`2026-10-05 16:4x +08`
写面（硬边界）：**只许改 `scripts/` 下的测试范围表 `scope=(...)` 数组**；⛔ 绝不碰 `.github/workflows/ci.yml`；⛔ 零 `go test/build/vet/env`（仅 `go list` 允许）。

## §0 起手锚与现读

- HEAD 锚：`9761d630`（`2026-10-05 15:32:16 +08`）。⚠ 票面 AC#0 的锚是 `4e66817`（`2026-09-24`），相差 11 天，票面所有 `file:line` 与"现场"数当**待验断言**。
- `go list ./...`：**本机 GOOS=windows = 35 枚**（`wc -l`=35）。票面 AC#0 写 **33** ⇒ **漂**（`cmd/balldebug`/`frontend`/`internal/agent/approval`/`internal/llm/adaptertest`/`internal/projctx` 等新增/改名所致，逐条列 §1）。
- 5 枚目标包的 `_test.go`（`git ls-files <pkg> | grep -c '_test\.go$'`，本机复跑）：
  - `internal/ball` = **16**（票面写 11 ⇒ 漂）
  - `cmd/wisp` = **65**（票面写 5 ⇒ 严重漂；236 枚里 cmd/wisp 独占 65）
  - `internal/perm` = **3**（票面写 2 ⇒ 漂）
  - `internal/plugin` = **1**（票面写 1 ⇒ 复现）
  - `cmd/llmrecord` = **1**（票面写 1 ⇒ 复现）
- ★★ **重大漂移**：票面"另有 5 枚带测试的包**零覆盖**"这一句，在 HEAD `9761d630` 上**不成立**。
  `bash scripts/portable-tests.sh --scope=census` 现读（GOOS=windows）：
  `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7` ⇒ **凡本平台能编译出测试文件的包，全部被某 scope 认领**；
  5 枚目标包逐条：`internal/ball`→`core,windows`；`internal/perm`→`core,windows`；`internal/plugin`→`core,windows`；
  `cmd/llmrecord`→`core,windows`；`cmd/wisp`→`cli`。**scope 表里都已有名字。** 详见 §1。

## §1 CI 名册真身

未判（待填：ci.yml 触发点行 + scripts scope 数组行 + winsec-tests/wisp-cli-tests 路由，逐条带行号）

## §2 五枚逐枚可纳入性

未判（待填：每枚 build tag 原文行 + runner + 新依赖 + 最坏颜色 + winlive 0 命中复跑）

## §3 实际改动与正控

未判（待填：本轮实际写了 scripts/ 哪几行；若无新增，具名说为什么；正控＝零 PASS 即 fatal 的那把真身）

## §4 壳尺读数（带时刻）

未判（待填：`sh scripts/d22scan.sh` + `sh scripts/check-path-length-budget.sh --with-self-test`，各带时刻与"落在谁的窗口里"）

## §5 判不动／量不到（具名归口）

未判（待填：需动 ci.yml 才能纳入的枚＝归 `111-r1`；远程步级结论＝欠 push，归编排者；等）

## §6 交件判语与 commit 链

未判（待填：逐枚 commit 号＋时刻；五枚逐枚判语；与编排者给的数不同处清单）
