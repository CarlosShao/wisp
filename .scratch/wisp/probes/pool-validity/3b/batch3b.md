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

### 第 1 组（150／159／165）—— 已判，已 commit

| 票号 | 票文件 | 上次动过 | 判定 | 复法（跑的那把尺） | 读数 |
|---|---|---|---|---|---|
| 150 | `150-four-deferred-markers-have-no-row-in-the-spec-12-registry-so-the-1-to-1-rule-is-violated-today-and-no-instrument-scans-it.md` | c6c6a499 09-26 | **仍成立** | ①代码侧标记尺：`git --no-pager grep -n -E 'DEFERRED\(D[0-9]+(-[0-9]+)?\)' HEAD -- internal cmd tools`；②登记表侧：`git --no-pager show HEAD:docs/specs/SPEC-12-roadmap-governance.md \| sed -n '/^## 5/,/^## 6/p' \| grep -c '\|'` 与同段 `\| grep -c 'D28'`；③仪器侧：`git --no-pager grep -in 'deferred' HEAD -- tools/d22scan`；④子号定义：`git --no-pager grep -n -E '\bD[0-9]+-[0-9]+\b' HEAD -- docs/PLAN.md docs/specs`；⑤台账：`grep -n 'Q-55' docs/reports/pending-and-issues.md \| tail -3` | ①**4 枚／2 种**原样在：`internal/agent/compress.go:27`、`compress.go:53`、`loop.go:401`（三枚 `D28-1`）＋`internal/agent/control.go:15`（一枚 `D11-3`）。⚠ 票面 `:11` 写的行号是 `:49`／`:394`，今树已漂到 `:53`／`:401`——**种类数 2、出现次数 4 这两枚数没变**，变的只是行号，不算翻案。②§5 段内含 `\|` 的行＝**30**（票面 09-26 现量＝**29**）⇒ 分母**自己又漂了一行**：新增那枚是"确认窗口借用裸 `Esc`（票 245 乙形残余）"那行 DEFERRED。`D28` 在该段命中 **0 次**（"No matches found"）＝**票面那条"表上查无此人"今天照旧**。③`tools/d22scan` 里 `deferred`（不分大小写）**0 命中** ⇒ "没有任何仪器扫它"这半句仍成立，本腿额外确认全 `tools/` 只有 `scripts/portable-tests.sh` 与 `internal/config/doc.go` 出现 `DEFERRED` 字样，**无一是 marker↔登记表的双向差集检**。④子号写法今树仍只有 `D45-1`／`D45-2` 两例（`PLAN.md:2028`、`:3111`、`:3115`、`SPEC-05:136`、`SPEC-12:23`）⇒ "`D28-1`／`D11-3` 这种写法在规格里没有任何文本定义"仍成立。⑤`Q-55` 在台账最新几枚 `next=` 里仍列在**未答**那一档 ⇒ 票 `Status: blocked` 的闸门没开。**备注**：本腿**没有**读 `docs/specs/**` 的新内容面去判"表里有没有行"以外的东西，只取 §5 那一段计数——票面 AC#1 要的正是"把行数尺钉成一句可复算的话"，本腿这一把给出 30，**且顺手证明了它那把尺本身会继续漂**（29→30 的原因是后来票 245 补了一行，不是数法变了）。⚠ 弱：AC#1 那格"该收成一个数"这件事仍没人在做，本腿的 30 只是第三枚分母，不构成翻案凭据。 |
| 159 | `159-record-level-cleanup-c1-c5-left-by-ticket-156-acceptance.md` | d0b72775 09-26 | **仍成立** | ①追加物在不在：`git --no-pager show HEAD:docs/evidence/s1/156-exited-asks-os-r1.md \| grep -n '收口程追加\|非 r1 原文'`；②那枚文件被谁最后动过：`git log --format='%h %ad %s' --date=format:'%m-%d' -- docs/evidence/s1/156-exited-asks-os-r1.md`；③AC#2 目标文件同尺：`git --no-pager grep -n '票 159' HEAD -- cmd/wisp/slo_report_144_windows_test.go`；④分母尺：`git ls-tree -r --name-only HEAD -- docs/evidence/s1 \| grep '/156- \| /159-'`；⑤台账：`grep -n '票 159' docs/reports/pending-and-issues.md \| tail -6` | 五笔条件（C1–C5）**一笔产物都没落**：①③两把尺**均 0 命中**——r1 证据件里既没有"编排者／收口程追加，非 r1 原文"那一节，`slo_report_144_windows_test.go`（今树 938 行）里也没有任何 159 的追记。②那枚证据件最后两笔是 `2262c2b6`／`2ce77e18`（都是 09-26 **票 156 自己的** r2/r4 程），此后**没人再写过它** ⇒ C3 点名的错形原样躺在盘上：`156-exited-asks-os-r1.md:283` 逐字仍是"改后 cmd/wisp **RUN=144  PASS=84**"挂在 gate-post 名下（票 159 说这发的真值是 143/83、+5 属 `gate-r3-pre`），**票面点名的那处"表与盘不符"今天还不符**。④台账里"票 159"**只出现在排队话术里**（`A337`／`A340`／`A341`／`A345`／`A346`／`A348` 六枚 `next=` 逐字写"票 159 与'全树门禁＋推送'**在最后**"）⇒ 立而不派＝从没起手，不是被谁顺手做掉。**边界（必须说）**：⑴AC#4⑵那半条要求核 `frontend/embed.go:19` 的 `//go:embed all:dist`，本腿禁读 `frontend/**` ⇒ 那一格本腿**量不到**，不影响整票〔仍成立〕（C1/C2/C3/C5 四笔的产物都 0 命中，凭据与那一格无关）。⑵票 156 的三枚验收/实现件仍在树上（`156-close-and-load-bearing-r1-accept-r1.md` 就是票面点名的来路），本腿**只按文件名与行号读，未据内容翻任何勾**。⚠ 弱：票面 `:68` 自己写了"本票两枚数字若与你现量不符以你为准"——本腿没有 go 命令，**没有**复算 143/83 那个日志读数，本腿的量到之处是"追加更正**没落**"这件事，不是"那对数字对不对"。 |
| 165 | `165-plan-first-then-act-a-brake-the-voice-entry-needs.md` | a651fe80 09-27 | **仍成立** | ①已批工具名在生产码里有有没有：`git --no-pager grep -n 'plan\.present' HEAD -- internal cmd tools`；②注册工具全集名册：`git --no-pager grep -n 'func (.*) Name() string { return' HEAD -- internal/tools`；③宽形尺：`git --no-pager grep -inE '"plan\.\|"propose\|planPresent\|PlanTool' HEAD -- internal cmd tools`；④冻结文本那格：`git --no-pager grep -n 'plan\.present' HEAD -- docs/PLAN.md docs/specs`；⑤commit 坐实：`git log -1 --format='%h %ad %s' --date=format:'%m-%d' 824bb9c` ＋ `git show --name-only --format='' 824bb9c` | ①**0 命中**；②`internal/tools/` 的 `Name()` 全集今树＝`fs.read`/`fs.list`/`fs.edit`/`fs.write`/`fs.trash`/`fs.move`/`fs.delete`/`task.output`/`task.cancel`/`task.spawn`（外加两枚 `_test` 夹具）⇒ **名册里没有 plan/propose 那一形**，票面"我们表上零枚"从 09-26 那句"来源"起，到今天**码侧一格没变**；③宽形尺整棵生产面只命中 1 处，且那处是 `internal/tools/bridge_test.go`/`git` 相关测试里读 `docs/PLAN.md` 的变量名 `plan`，**与"先出方案"这档无关**；④⑤⛔ **本腿明确不把它降进〔已失效〕**：owner 批下来的是**契约文本**那一半——`824bb9c0`（09-27，标题逐字"票 165 契约文本落地（owner 批准＝甲路，台账 A322）"）的 `--name-only` **只有两枚文件**＝`docs/PLAN.md`＋`docs/specs/SPEC-07-tools-and-plugins.md`（今树落在 `PLAN.md:2567`、`SPEC-07:73`，两行都逐字写着"票 165（2026-09-27 owner 批准新增，走**甲路**）…**不新增状态**"）⇒ 那枚 commit 是**本票的前置条件被满足**，不是本票的交付；票面 `Status 更正`自己就是这么读的（"可做的是把这枚工具按已批准的甲路**实现出来**"）。**所以今天成立的是后半截缺陷**：D34 表已把 `plan.present` 标 **L2／S3**，实现体**一行没有**＝一个已进冻结名册、零实现、零登记缺口的工具。**派单前置**那两格本腿判到：owner 那句已由 `A322` 给了 ⇒ 门开；票 161 是否结案本腿不判（越出分母）。AC#1"未修码读数"要跑真链路才能答"今天能不能只读出方案"⇒ 那一格〔量不到〕，但整票〔仍成立〕由①②③三把尺独立坐实。⚠ 弱：本腿没读 `cmd/wisp/**` 内容面（并行腿 `231-v1` 地界），若组合根里有未过桥的 plan 草稿，①③那两把尺在 `internal/tools` 之外**仍能命中**（本腿的尺是整 `internal cmd tools` 三棵），所以这一条弱得不致命。 |
| 166 | `166-no-session-layer-list-search-rename-archive-fork-...md` | f1ce2f2b 09-26 | 未判 | 未判 | 未判 |
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
