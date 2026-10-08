# 182-c2 · CI 步级 `if:` 现量与"本票要不要新增 CI 步骤" — 锚 HEAD `a7e9b2e7`

## 0. 结论先给
**票 182 的六枚未勾框里，没有一枚需要新增 CI 步骤**（本票是只读普查票，票面 `:29` 逐字禁 `internal/**`；`AC#6` 那条"契约轴"要的是**别碰**，不是**多跑一步**）。⇒ 派单点名的那把尺**不适用本程**；仍按要求现量一次，读数如下（供编排者引别的票时用，⛔ 不许被读成"本票动过 CI"）。

## 1. 结构尺现量（⛔ 不用词频尺；派单第 129 条：同一枚数两种尺写法差 1）

| 问的是什么 | 尺原文（结构形状） | 现量 | rc |
|---|---|---|---|
| **步级 `if:` 枚数** | `grep -nE '^        if: ' .github/workflows/ci.yml \| wc -l`（8 空格＝step 内） | **13** | 0 |
| 其中 `!cancelled()` | `grep -nE '^        if: \$\{\{ !cancelled\(\) \}\}' ci.yml` 逐名 12 行：`234 294 331 387 522 560 570 592 643 732 768 784` | **12** | 0 |
| 其中 `always()` | `:462` `if: always()` ＝"Stop compose services"那一枚清理步（`docker compose … down`，`ci.yml:455-462` 上下文可核） | **1** | 0 |
| **job 级 `if:` 枚数** | `grep -nE '^      if: \|^    if: ' ci.yml \| wc -l` | **0** | 0 |
| `runs-on` 枚数（＝job 数） | `grep -nE '^    runs-on: ' ci.yml \| wc -l` | **6**（`:66 lint`／`:396 test-core`／`:506 test-windows`／`:791 slo-smoke`／`:849 slo-full`（self-hosted, wisp-slo）／`:923 lint-frontend`） | 0 |
| 对照：**词频尺**（派单警告那一把） | `grep -c 'cancelled' .github/workflows/ci.yml` | 那把数的是**含该词的行**（含注释／含 `ci.yml:277` 那种解释句），⛔ 不是"有几处守卫"——本件不用它作枚数 | — |

## 2. ⚠ 派单那句"新步骤必须带步级 `if: !cancelled()`"要和文件里另一条硬规矩**分两格写**（本腿现量到的 nuance）
- **要**：`ci.yml:277` 逐字"`if: ${{ !cancelled() }}` is the one shape"⇒ 步级守卫只有这一形；缺它的新步在真实 run 里会是 `[skipped]`＝**从未执行**（这条纪律的出处＝编排者自己记的第 109 条）。
- **不要**：同一个文件在**job 级**逐字禁 `if:`——`ci.yml:34-35` 逐字"No `if:`, no continue-on-error, no skip flag anywhere in this block（D22 mode-6: 'A skippable job is a job that will one day be skipped'）"、`ci.yml:920-921` 同一句写在 `lint-frontend` 头上、`ci.yml:630`/`:884`/`:34` 各有一处同形禁令。
⇒ **正确写法＝新步带 8 空格步级 `if: ${{ !cancelled() }}`；新 job ⛔ 不许带任何 job 级 `if:`／`continue-on-error`／skip flag。**⛔ 不许把派单那句话读成"整份 yaml 想要 `if:` 就行"，那会撞上 mode-6 那道扫描。

## 3. 与本票那一族读数有关的一条旁证（不展开，只指认，⛔ 不重跑）
`memory` 与 `docs/reports/pending-and-issues.md`（`A693`／`A691` 一带）记着：`test-windows` 里 `go vet` 与 `gofmt denominator` 两枚曾在真实 run 中是 **skipped＝从未求值**，而 `winlive` 步只在 HEAD 上从未进过 run。⇒ 任何"这道门在 CI 有牙"的断言，判据只能取**该步在某发 run 出现过 success/failure**，⛔ 不许把"它在 yaml 里"读成"它跑过"（第 109 条原文）。本腿 ⛔ 零 push、⛔ 无取 run 权限，故**不引用任何一发 run 的颜色**，只把这条纪律落在文件里。
