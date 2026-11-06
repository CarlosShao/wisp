# 136-ac5-ac6-impl — 票 136 AC#5（refuse-silent）＋ AC#6（schedule 到点没响）交付表

- 实现 agent＝`worker-ticket136-ac5-ac6`；只 commit 不 push。
- 锚点：起手 HEAD `84b0a7c4d`（派单锚，现量一致）；本程写面＝`scripts/slo-freshness.sh`（仅注释）＋ 本表 ＋ 票面两枚 `>` 块。
- `internal/observe/**`、`tools/d22scan/**`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、阈值/golden 零改动。

## §1 AC#5：`paths:` / `paths-ignore:` 这一维明写"不查"

- 票面给的两条路：①把禁令补进 `tools/d22scan`（不在本票地界，需编排者具名解冻 ⇒ 本程未动）；②在票面与注释里明写"这一维不查"。本程取 ②，拒绝沉默。
- 现量核对（2026-11-06，HEAD `84b0a7c4d`）：
  - `tools/d22scan/main.go` 头注的 ban 名册是 #1–#9（max-clock / mirror-hash / panel-approval / internal-artifact-tool / emoji / phantom-citation 等），**没有任何一条**读 workflow YAML 的触发过滤键 ⇒ AC#5 说的"禁令里少了这一件"在当前盘上仍成立。
  - `scripts/slo-freshness.sh` P1 的四条 grep（`on:`/`push:`/branches main-dev/`cron:`）与 job 体的 `if:`/`continue-on-error` 检查，全部**不读** `paths:`/`paths-ignore:`，且锚在触发块整体而非 `push:` 子树 ⇒ 在 `push:` 下加 `paths:` 过滤能通过 P1 全部检查，slo-full 的自动触发会被静默收窄到"只碰指定文件的提交"。
- 落点：`scripts/slo-freshness.sh` P1 static 段注释（commit 见 §4），逐点写明：不查什么、为什么 greps 抓不到、D22 mode-6 第四禁（"no path filter"，见 `slo-fresh.yml` 头注）当前只由 review 承担、补进 d22scan 需具名解冻。**同一句话也落在票面 AC#5 格下的 `>` 块。**
- 未做档：把这一件真正补成 d22scan 的 ban（或 allowlist 例外形状）＝须编排者具名解冻 `tools/d22scan/**` 后另派；本程没有替编排者扩权。

## §2 AC#6：schedule "从未产出一枚 run" 的前提已被 GitHub 侧现量推翻

票面 09-23 18:1x（UTC 10:17）的读数 `gh run list --event schedule ⇒ []` 是**取数时刻太早**，不是配置缺陷；同日 **11:50:44Z** 首枚 schedule run 已产出。本程现量（2026-11-06）：

- 根因核：
  - workflow `slo-fresh` state＝**active**（`gh api …/actions/workflows`）；
  - 仓库默认分支＝**dev**（`gh repo view`），`slo-fresh.yml` 在 dev 上（`git log origin/dev -- .github/workflows/slo-fresh.yml` 顶＝`44ab500f0`）⇒ "schedule 只按默认分支的配置评估"这条规则的生效面成立，无阻塞；
  - 文件首次入库＝`b9b20723f`（commit 时刻 2026-09-23T02:41:50Z；票面写的 03:41:11Z 与 `git show -s --format=%cd` 不符，只记差异，不影响结论）；
  - **06:23Z 那格未响的原因**：GitHub 对新注册 schedule 的首次评估不保证落在下一个 cron 格点（其调度器按自己的负载节奏出队，官方文档明说 schedule run 会延迟、高负载时段会丢格）。10:17Z 取数时首枚 run 尚未出队 ⇒ 当时读到的 `[]` 是"还没出队"，**同日 11:50:44Z 出队后此疑问关闭**——不成立的是票面"从未产出"这个全称判断，不是 workflow 配置。
- 现量读数（命令可复跑）：
  - `gh run list --workflow slo-fresh.yml --event schedule --limit 500` ⇒ 共 **142 枚**，首枚 `2026-09-23T11:50:44Z`；
  - 按日分桶：**2026-09-23 至 2026-11-05 共 45/45 天，每日 ≥2 枚**（最少的 10-05＝2、11-02..11-05 各 2，多数 3–4 枚）⇒ "每日 ≥1 枚"成立且有 2 倍余量；
  - nominal cron `23 */6 * * *`（00/06/12/18 Z）对实际出队时刻普遍延迟数十分钟到数小时（GitHub 侧已知行为），**频率与覆盖不受影响**。
- 失败枚（72 success / 28 failure，最近 100 枚内）逐因归档：
  - `38561478699`（10-14 12:07Z）＝`slo-full-scan-incomplete`，即票 136 AC#3 的 cap 钉当天如实打红（total 350 > cap 300），同日已由 `:-20` 收口提交修复——**fail-closed 按设计响**；
  - 10-18…10-24 每日 3–4 枚＝`slo-full-sample-stale`（"last VALID SLO SAMPLE is 4–10 day(s) old (> 3) … that was 2026-10-14T00:57:21Z run 38484004521"）——**钉在工作**：它如实报出 10-14 之后本机约 10 天没产出有效样本（D32 未被验证的真空期），而不是被绿 job 骗过；10-25 起恢复 success；
  - **2026-10-25 至 2026-11-05 连续 12 天全 success**，最近一枚 `42100267869`（11-05 20:06Z, 22s）。
- 结论：**AC#6 要的"修好后证明每日 ≥1 枚"已由 GitHub 侧现状满足，本程零配置改动**（票面 ⚠"改共享触发/并发配置要按事件类型分别核"一条因此未触发——没有可改的东西）。若编排者认为还需补强（如对 4 格 cron 全缺做显式红），属新格，不在 AC#6 判据内。

## §3 门禁读数（本机，2026-11-06）

按派单四道门原样复跑（`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`unset GOTMPDIR TMPDIR TEMP TMP`、未设 `CGO_ENABLED=0`）：

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | rc=0（原文见 §4 commit） |
| test | `go test ./internal/observe -count=1 -timeout 420s` | rc=0，`ok` |
| gofumpt | `gofumpt -l <changed>` | 无输出 |
| d22scan | `sh scripts/d22scan.sh` | rc=0 |

另跑：`sh -n scripts/slo-freshness.sh` rc=0（本程对脚本的改动是**纯注释**，无语义行变化；shellcheck 属 CI 步，本机未装则由 CI 承担——纯注释改动不产生新 lint 面）。

## §4 本程 commit 与残留

- commit：`scripts/slo-freshness.sh`（P1 段 AC#5 声明）＋ 本表 ＋ 票面两枚 `>` 块与一条 progress log。
- 未做档：①d22scan 补 ban（需解冻）；②AC#5/AC#6 的勾**未自翻**（归编排者据非实现者复核翻）；③AC#7/AC#10 不属本程。
- 凭据：本表零凭据、零 token；GitHub 读数均为只读 API。
