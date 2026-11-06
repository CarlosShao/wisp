# 136-verdict — 票 136 AC#5（refuse-silent 声明）＋ AC#6（schedule 每日 ≥1）非实现者终裁

- 裁定方＝非实现者（本程未写过 `scripts/slo-freshness.sh`、未写过票面 `>` 块、未写过 `136-ac5-ac6-impl.md`；实现程＝`worker-ticket136-ac5-ac6`）。
- 锚点：派发单写 HEAD `edab04ab6`；本程起手现读 HEAD `ac2c1d903`，`edab04ab6` 是其祖先，`git diff --name-only edab04ab6..HEAD` 恰 1 枚＝`docs/reports/pending-and-issues.md`（台账，本程不读不碰）⇒ 判定面上两锚等价。实现程锚 `84b0a7c4d`（其交件自述，本程核 `001b516ce` 的父系存在）。裁定时刻 `2026-11-06T02:51:29Z`（`date -u` 现量）。
- 引用层与判断层分开：§1/§2 每格先列〔独立复现〕的新鲜读数（本程自跑），再给判语。

## §0 门禁独立复跑（本机，本程现跑；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`unset GOTMPDIR TMPDIR TEMP TMP`、未设 `CGO_ENABLED=0`，Git Bash）

| 门 | 命令 | 本程读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | rc=0 |
| test | `go test ./internal/observe -count=1 -timeout 420s` | rc=0，`ok github.com/CarlosShao/wisp/internal/observe 5.362s` |
| gofumpt | `gofumpt -l ./cmd ./internal ./tools` | 无输出，rc=0 |
| d22scan | `sh scripts/d22scan.sh` | rc=0，`clean - no D22 ban violations`（ban #8 internal/=906、cmd/=213，无下降） |
| sh -n | `sh -n scripts/slo-freshness.sh` | rc=0 |

五道门与实现程 §3 的自述读数一致，无一红。

## §1 AC#5 判语：**PASS（无附条件）**

**〔独立复现〕的事实层**（本程现量，不抄交件）：

1. **盘上缺口仍在**：`grep` `tools/d22scan/**` 全量 ⇒ 0 处命中 `paths:`/`paths-ignore:`/workflow 触发过滤键 ⇒ "D22 禁令里少了这一件"在当前 HEAD 仍成立。
2. **P1 确实不查**：`scripts/slo-freshness.sh:194-205` 四条 grep 只锚 `on:`/`push:`/`branches:…dev`/`cron:`，`:214-225` job 体只查 `if:`/`continue-on-error` ⇒ 在 `push:` 子树下加 `paths:` 过滤能通过 P1 全部检查（greps 锚在行首层、不在 `push:` 子树）。
3. **声明已落两处**：`scripts/slo-freshness.sh:175-187`（P1 static 段注释，逐点写明不查什么／为什么 greps 抓不到／mode-6 第四禁当前只由 review 承担／补 ban 需具名解冻）；票面 AC#5 格下 `>` 块（票 `:48-54`，含撤销口令「撤 136 AC#5 不查声明」）。
4. **写面干净**：实现 commit `001b516ce` 的 `git show --numstat` 恰 3 枚文件（票面 24+0、证据 52+0、脚本 14+0），脚本 hunk **纯注释新增、零删除、零语义行**；`internal/`、`tools/d22scan/**`、`.github/**`、阈值/golden 零改动。

**判语**：票面 AC#5 给了两条路——①把禁令补进 `tools/d22scan`（明文"不在本票地界，需编排者具名解冻"）；②"在票面与注释里明写'这一维不查'，不许沉默"。实现程取 ②且两处落笔内容经本程序逐点核实为真（第 1、2 条读数与声明文字逐句相符）。**判定走 ② 即满足本格判据，PASS**。残留（d22scan 补 ban）按票面文字本就不属本格结案条件，登记为待解冻新格，不阻塞。

## §2 AC#6 判语：**PASS（无附条件）**

**〔独立复现〕的事实层**（`gh` 只读 API，本程 2026-11-06 现跑，命令在交件 §2 可复跑）：

1. `gh run list --workflow slo-fresh.yml --event schedule --limit 500` ⇒ **total=142**，首枚 `2026-09-23T11:50:44Z`、末枚 `2026-11-05T20:06:37Z`。
2. 按日分桶：**2026-09-23 至 2026-11-05 共 44/44 天每日 ≥2 枚**（最小值 2：10-05／10-12／11-02…11-05；多数 3–4 枚），0 天低于 2。
3. workflow state＝**active**（`gh api …/actions/workflows` 现量）；默认分支＝**dev**（`gh repo view` 现量），`slo-fresh.yml` 在 dev ⇒ "schedule 只按默认分支评估"的生效面成立。
4. 失败枚抽查：`38561478699`＝`2026-10-14T12:07:32Z` schedule **failure** 现量（＝票 136 AC#3 cap 钉当天如实打红）；10-20…10-24 每日多枚 failure（＝`slo-full-sample-stale` 如实报 D32 真空期，钉在工作）；`42100267869`＝`2026-11-05T20:06:37Z` **success**（现量）。
5. 票面"从未产出一枚 schedule run"的全称前提**被 GitHub 侧现量推翻**：首枚 run（09-23 11:50:44Z）晚于票面取数时刻（09-23 10:17Z）仅 1.5h ⇒ 当时的 `[]` 是"还没出队"，不是配置缺陷。

**判语**：AC#6 的判据是"查明原因＋修好＋**证明每日 ≥1 枚**"。原因已查明（取数时刻太早，非配置问题），修好一栏以"零配置改动即满足"成立（没有可修的东西，票面 ⚠ 并发改写一条因此未触发），"每日 ≥1"由本程第 2 条新鲜读数证明且带 **2 倍余量**。**PASS**。
**一处读数更正（不影响判语，照诚实口径登记）**：交件写"45/45 天"，本程分桶现量 **44/44 天**（8+31+5=44，交件算术差 1）；"每日 ≥2"逐天成立不受影响。

## §3 终裁结论

- **AC#5：PASS。AC#6：PASS。** 两格均无附条件。
- **closeable：yes** —— 编排者可以把票面 AC#5、AC#6 两枚勾 `[ ]`→`[x]` 翻勾（**safe-to-flip boxes**：仅这两枚框；其余 AC#7／AC#10 等格不在本表射程）。
- 不阻塞的登记项（翻勾后仍在册）：① d22scan 补 workflow 触发过滤 ban——需编排者具名解冻 `tools/d22scan/**` 后另派（AC#5 走的第二条路已消除"沉默"，但仪器缺口本身仍在，本表第 1.1 条现量）；② 若要给"4 格 cron 全缺"做显式红——属新格，票面已写明不在 AC#6 判据内。
- 撤销口令（沿用实现程）：**「撤 136 AC#5 不查声明」** ⇒ 删 `scripts/slo-freshness.sh:175-187` 注释块＋票面 `>` 块即复原；AC#6 无配置可撤。
