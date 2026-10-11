# `304-r1` 00 — 起手锚（现跑）

腿＝`304-r1`（写腿，`AC#1` 落地形＝甲：Job Object 接到两枚 mockllm spawn 点）。
派单＝`.scratch/wisp/probes/orch/20261011/txt/304r1-dispatch.txt`。票面＝`.scratch/wisp/issues/304-external-killed-test-binaries-orphan-mockllm-exe.md`。

## 锚（第 0.1 步，全部现跑）
| 尺 | 命令逐字 | 读数 |
|---|---|---|
| HEAD | `git rev-parse --short HEAD` | `3974ab4d` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` |
| 工作树脏枚数（产码射程） | `git status --porcelain -- internal cmd tools scripts` | **0 行** |
| 全索引暂存枚数 | `git -c core.quotePath=false diff --cached --name-only \| grep -c ''` | **0 枚** |

⇒ 起手索引**干净**：我看到 0 枚别家路径，我没动任何人的东西（派单 §0.2 的"发现大面积"分支未触发）。

## 与派单转述的比对
- 派单 §1 说锚号"不写在派单里"＝这一刻的读数即锚，无冲突可比。
- 票面 `:15` 的两枚 spawn 点（`internal/llm/adaptertest/mockllm.go:65` build＋`:72` start；`internal/llm/openaichat/mockllm_integ_test.go:73` start、`:106` 现有 `t.Cleanup`）已在 `304-a1` 表①复量为"恰两枚"（词面尺 `print-addr`），与派单一致 ⇒ 无冲突。
- 台账 `A850` §4 的裁形（甲）＋边界①②与派单 §2.5/§2.6 逐字一致 ⇒ 无冲突。
- ⚠ 一处**行号快照过期**要具名报：派单 §1 写 `internal/proc/jobscope_windows.go:71 OpenJobScope`／`:110 StartInJob`／`:170 Close`，我落件时现量（见 `10-impl.md`），比对结果只作"我这轮的读数"，⛔ 据此改任何人的原句。

## 本腿写面（派单 §4 末条）
①`.scratch/wisp/probes/304/r1/**` ③票 304 末尾**一条** append-only bullet ＋ `internal/llm/**` 那几枚字节。
⛔ `docs/reports/**`／`docs/evidence/s1/**`／`.github/**`／`ci.yml`／`AGENTS.md`／`docs/PLAN.md`／`docs/specs/**`／`internal/observe/**`／`tools/mockllm/**`／`tools/d22scan/**`／`internal/proc/**`（只读复用，⛔ 改一字）。

## 声明
⛔ 本腿不自翻任何框、⛔ 不自宣布任何 `AC` 成立（判语归非实现者腿 `304-v1`）。⛔ push。
