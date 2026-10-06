# pool-validity-3b / 批3b —— 零勾开放票有效性普查（号段 150–199 的零勾开放票，18 枚）

> 只读普查腿 `pool-validity-3b`。**刻意做小**：前两枚同职腿（`pool-validity-1`／`pool-validity-2`）都死于"扫全池"，
> 本腿只切 150–199 这一段里"零勾"的那一片。
> 起手 HEAD `4b15869d`（2026-10-06 10:12 +0800），分支 `dev`。⚠ 换 HEAD 要重量。
> 本腿**零 go 命令**（并行腿 `231-v1` 此刻在 `cmd/wisp` 取整包终态）；枚数一律 `git ls-tree`／`git ls-files`。
> 工具全集＝`git log`／`git grep`／`git show`／`git ls-tree`／`ls`／`grep`／`wc`／Read。
> 票面零改动：零翻勾、零改名、零 `Status:` 改写、零产码改动、零台账写入、`docs/**` 一字不碰。
> 唯一写入路径＝`.scratch/wisp/probes/pool-validity/3b/**`。
> ⛔ `pool-validity/1/**` 只读不采信；`pool-validity/2/batch1.md`·`batch2.md` 只学表形与判据密度，
> **不照抄它的口径**；⛔ 不动 `pool-validity/2/batch3.md`·`batch4.md`（前一枚死腿的空骨架）。
> ⛔ 零读 `frontend/**`、`design/**`、`.gitignore`、`cmd/wisp/**` 内容面、`scripts/**`、`.github/**`、
> `docs/reports/HANDOVER.md`、`docs/evidence/s1/**` 内容面，
> 以及并行腿半成品 `probes/231/v1`、`probes/269/a1`、`probes/111/a5b`、`probes/pool-validity/3a`。

## §0 分母：本腿自己现量的尺

在 `.scratch/wisp/issues/` 目录内逐字跑的两条（命令原文）：

```
ls [0-9]*.md | grep -v -- '-done' | awk -F- '$1+0>=150 && $1+0<=199' | wc -l   -> 37（号段内开放票）
ls [0-9]*.md | grep -v -- '-done' | awk -F- '$1+0>=150 && $1+0<=199' \
  | while read f; do [ "$(grep -c '^- \[x\]' "$f")" -eq 0 ] && echo "$f"; done | wc -l
                                                                               -> 18（本腿分母）
```

- **本腿名册＝18 枚**，与编排者给的 18 对上。
- ⚠ 已知尺漏点（承 `pool-validity/2/summary.md` §4 第 5 条）：勾选框尺只认顶格 `^- \[x\]`，
  **缩进框读不到** ⇒ 18 这个分母可能偏大（含了其实有勾、只是缩进的票）。本腿逐枚遇到时在备注列具名。

## §0b 四档尺（本腿自己的判据；不照抄 pv2）

| 档 | 判据 | 复法要求（必须留命令原文＋命中 file:line 或 0 命中） |
|---|---|---|
| **仍成立** | 票面点名的那枚缺陷今天还在树上 | 把它 §现场／AC 里**最硬的一条**自己跑尺：`git --no-pager grep -n '<符号或那句字面>' HEAD -- internal/ cmd/` |
| **已失效／已被别人做掉** | 缺陷不在了 | ⛔ **硬门**：一枚 commit 号＋`git log --name-only` 命中，**或**今树 0 命中；拿不出不许写这档 |
| **差翻勾** | 缺陷不在了，凭据是票内注记／台账 `A##`／非实现者裁决表 | 凭据种类写进复法列；**不等于**本腿裁它可结案 |
| **量不到** | 判据落运行期行为／真开窗／**界面侧已委托**（`skipped=frontend(owner-delegated)`） | 归口写清缺哪一行读数 |

★ `上次动过` 列尺＝`git log -1 --format='%h %ad' --date=format:'%m-%d' -- <票文件>`，18 枚全部本腿自己现跑。
★ 尺面统一为**提交面**（`git --no-pager grep -n <pat> HEAD -- <path>`），工作树有别的腿在飞，用 HEAD 才不被脏面洗。
★ "那个函数还在"≠"缺陷仍成立"的充分凭据 —— 本腿每枚都要求把**票面那句缺陷**本身落到行号。

## §1 名册与判定（18 枚）

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 150 | `150-four-deferred-markers-...md` | c6c6a499 09-26 | 未判 | 未判 | 未判 |
| 159 | `159-record-level-cleanup-c1-c5-...md` | d0b72775 09-26 | 未判 | 未判 | 未判 |
| 165 | `165-plan-first-then-act-...md` | a651fe80 09-27 | 未判 | 未判 | 未判 |
| 166 | `166-no-session-layer-list-search-...md` | f1ce2f2b 09-26 | 未判 | 未判 | 未判 |
| 168 | `168-built-in-browser-as-a-disable-able-plugin-...md` | a651fe80 09-27 | 未判 | 未判 | 未判 |
| 169 | `169-a-single-decorative-glyph-in-frontend-...md` | 82476e85 09-26 | 未判 | 未判 | 未判 |
| 170 | `170-the-cli-legs-l1-block-window-...md` | 27f798ec 09-26 | 未判 | 未判 | 未判 |
| 172 | `172-no-ruler-links-the-frozen-tier-...md` | 9b5f9e9d 09-27 | 未判 | 未判 | 未判 |
| 173 | `173-the-audit-line-prints-in-allowlist-scope-true-...md` | d0c164a1 09-27 | 未判 | 未判 | 未判 |
| 178 | `178-pair-census-ruler-g6-negative-leg-...md` | 386716a6 09-28 | 未判 | 未判 | 未判 |
| 182 | `182-the-task-monitor-rail-is-outside-ticket-145s-...md` | fccaf3e3 09-28 | 未判 | 未判 | 未判 |
| 185 | `185-every-successful-reread-mints-the-blocker-...md` | 1f22f1aa 09-28 | 未判 | 未判 | 未判 |
| 186 | `186-the-composer-area-should-detect-git-...md` | 38747655 09-28 | 未判 | 未判 | 未判 |
| 187 | `187-the-model-and-thinking-tier-...md` | 03c01d71 09-28 | 未判 | 未判 | 未判 |
| 189 | `189-the-review-stack-needs-uncommitted-count-...md` | fccaf3e3 09-28 | 未判 | 未判 | 未判 |
| 194 | `194-the-panel-has-two-method-rosters-...md` | 74524fe5 09-28 | 未判 | 未判 | 未判 |
| 195 | `195-config-set-needs-a-per-section-write-foundation-...md` | f25f9572 10-02 | 未判 | 未判 | 未判 |
| 196 | `196-two-task-state-vocabularies-already-coexist-...md` | fccaf3e3 09-28 | 未判 | 未判 | 未判 |

## §2 本批计数

| 档 | 枚数 | 票号 |
|---|---|---|
| 仍成立 | 未判 | — |
| 已失效／已被别人做掉 | 未判 | — |
| 差翻勾 | 未判 | — |
| 量不到 | 未判 | — |
| 合计 | **18** | ✓ 与 §0 分母对上 |

## §3 界面侧委托（本腿按编排者要求单列的一栏）

名册 18 枚里带 `skipped=frontend(owner-delegated)` 标记的读数（本腿自己 `grep -n 'skipped'` 跑出来的）：
待逐枚补完。命中同类标记者判〔量不到〕并写明"界面侧已委托、本编队永不写"。
