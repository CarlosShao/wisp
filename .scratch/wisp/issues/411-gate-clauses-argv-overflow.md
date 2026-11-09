# 票 411 —— gate-clauses.sh argv 溢出：14 腿里 10 腿假 BAD（编队级活缺陷，仪器程）

- 立案：2026-11-09 09:5x，写腿（修程即本票实现腿）。
- 来源：`.scratch/wisp/probes/154/r20261108/00-argv-overflow.md`（census GUARD D 腿 2026-11-08 23:4x 普查时发现，未立案转本票）＋派单（编排者 2026-11-08）。
- 性质：**共享门禁仪器坏**——不是某票的判据红，是仪器对所有票都跑不出真读数。

## 病（现量）

- 仪器＝`.scratch/wisp/probes/154/gate-clauses.sh`（票 154 AC#1 触发门＋158/161/171/178 历次加腿，14 腿聚合退码）。
- `run()`/`pair()` 内部对 `$@` 未引号 ⇒ bash 在**工作树**上把 `'internal/**/*.go'`／`'cmd/**/*.go'`
  做路径名展开（bash 无 globstar，`**` 退化成一层 `*`）⇒ 摊成逐枚文件字面 pathspec
  **1084 词／39,836 字符** ＞ Windows CreateProcess 32767 命令行上限 ⇒ `git grep` exec 失败
  rc=126 ⇒ 14 腿里 10 腿记 `git-grep-rc-126` 假 BAD，聚合退码恒 10。
- 复跑命令（修前）：`bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ rc=10；
  `bash -x` 同跑，trace 里每腿 `git grep` 后 `meas=git-grep-rc-126`。
- 历史对照：10-14 票 178 收口时同展开 496 枚路径（未过限）能跑（rc=2）；其后树长肥过限，
  仪器从「能跑」变「全红」。

## 影响面（现跑 grep 取证）

- 票 **195 AC#5**（`.scratch/wisp/issues/195-config-set-...md:28`）逐字引用本仪器读数当验收判据。
- 票 154（AC#1/AC#3 本体）、158、161、171、178 全族的读数与回归件
  （`probes/161/r6/flip-declaration.sh`）都压在这把尺上。
- `.scratch/wisp/probes/**` 里 **52 份**探针/证据件引用 `gate-clauses.sh`。

## 修法（选型与理由）

- **甲（落地）＝`-- $@` 改 `-- "$@"`**（`run()` 与 `pair()` 各一处）：引号保住 pattern 本体
  （两枚短参数），由 git 在 `$A` 那棵树上自己展开。git grep 原生吃 glob pathspec（修前正控已证
  rc=0 有命中）。argv 从此固定两枚短参数，**树再长肥也不会再溢出**。
- 乙（xargs 分批）否：多进程拼名册、每批要重喂排除项，实测分批模拟还有静默失败
  （`probes/154/r20261108/13-equivalence-xargs.sh` 的 stderr 被 2>/dev/null 吞掉导致读数偏小）。
- **丙（--pathspec-from-file）否：`git grep` 根本没有这个选项**（现跑 `git grep
  --pathspec-from-file=…` ⇒ `error: unknown option`；它是 add/commit/restore 家族的选项）。
  派单倾向丙的判据（git 原生支持）对 git grep 不成立，故落甲。
- 顺带修（同件纯噪音，不改读数）：脚本 587/590 两行 echo 的未转义反引号（`` `ring 1` ``／
  `` `改完之后…` ``）被 bash 当命令替换执行（stderr `ring: command not found`）⇒ 转义 `\``，
  输出字节不变。
- **判据一字未动**：14 腿名册、want/want_n 逐枚、pair 四堆算术、聚合退码，全部原样。

## 射程语义变化（须裁，本腿不改判据）

修后射程从「bash 两层展开（工作树、1084 枚）」变「git 递归 glob（$A 树、1209 枚）」：

1. 这同时收掉票 178 随单发现⑤前半条：名册**文件清单**从此与内容同源（都来自 $A），
   合成树里 commit 过的新文件这把尺看得见了。
2. 盲区定量：旧展开结构性看不见第 3 层及以下目录，HEAD 树上现量 **9 枚文件**带腿词表命中
   （8 枚 .go＝`internal/agent/approval/{batch,cardrender,gate,ui}.go`＋
   `internal/plugin/manifestslot/{command,install,lifecycle,tier2}.go`；另 1 枚 .html 是
   普查词表误入，仪器 glob 只吃 .go）。普查件＝`probes/154/r20261108/15-glob-equivalence.sh`
   （git 原生 `:(glob)` 魔法做旧射程精确等价形对照）。
3. ⇒ 修后读数＝这把尺**史上第一次**对全树的真实读数，与 09-27/10-14 登记的基线不可直接比
   （基线全部是在「两层盲区射程」下登记的）。

## 修后真读数（锚 cb2287a5，2026-11-09 09:4x，rc=9）

聚合退码 **9**；腿数断言 14/14 相符；基线过期 0。逐腿：

| 腿 | 声明 | 基线 | 实测 | 判 | 具名 |
|---|---|---|---|---|---|
| G1 | quiet | 0 | 1 | BAD | `internal/plugin/manifestslot/tier2.go:46` `ex.Execute(ctx, plugin.Tier2ToolRequest{`（词表子串命中 Tier2ToolRequest——是 G1 想抓的形状还是 pattern 需要 -w，**归裁决者**） |
| G1b | quiet | 0 | 4 | BAD | tier2.go:46、`internal/plugin/tier2_gate.go:210`、`tier2_loader.go:195`、`internal/tools/subagent_197.go:616`（Tier2/Goja 执行面，10-14 后新落） |
| G2 | ring | 1 | 2 | BAD | `cmd/wisp/run.go:1514`（OpenTask）＋`:1520`（CloseTask）＝resident 钩子新落的一对 |
| G3 | quiet | 0 | 0 | ok | — |
| G4 | quiet | 0 | 0 | ok | — |
| G5 | ring | 1 | 2 | BAD | 在册 `cmd/wisp/panel_assets.go` ＋ **新** `cmd/wisp/resident_present_windows.go`（开方调用点=1，同文件无 Close——**真红候选，归裁决者**） |
| G5pos | quiet | 0 | 1 | BAD | 同上 resident_present_windows.go（正控样本被第二枚活形状污染，机制正确） |
| G5neg | ring | 8 | 14 | BAD | +6＝resident_present＋394/382/403 家族测试＋provenance_syncdirs_other_test（负控按设计响，枚数漂移） |
| G6 | quiet(piles) | 0 | 0 | ok | — |
| G6pos | quiet(rev) | 0 | 0 | ok | — |
| G6neg | quiet(piles) | 0 | 1 | BAD | ②只开＝`internal/tools/web_bridge_test.go`（测试文件只开不关，10-14 后新落） |
| G7 | ring(rev) | 1 | 6 | BAD | 反向 +5＝`resident_voice_bargein_windows.go`、`resident_voice_windows.go`、`internal/plugin/goja_runtime.go`、`manifestslot/command.go`、`manifestslot/lifecycle.go`（拿 scope 未同文件挂收尾——真漏还是委托，**归裁决者**） |
| G7pos | quiet(rev) | 0 | 0 | ok | — |
| G7neg | ring(rev) | 2 | 7 | BAD | 同 G7 家族漂移 |

- 完整逐行读数存档：`probes/154/r20261108/`（gate-before/after/time/mutation/revert 各份）。
- 本腿**不改任何 want/want_n/声明**：重登基线＝改判据周边的裁决动作（非实现者裁，
  178-v1 判语先例）；上表 9 腿真红全部具名交 154/158/161/171/178/195 家族裁决。

## 防作弊正控（导出合成树，git clone --shared，只建不删）

- 注入：clone2 里 commit `a512cd7a9`（＝锚 cb2287a5＋`internal/tools/mutation_411_openonly.go`
  一枚只 OpenTask 不 CloseTask 的生产形状）⇒ 门**红**：G6 ②只开 0→1（BAD）、G6neg 1→2、
  G2 +2（含 interface 类型行），聚合退码 9→**10**。
- 还原：同 clone 里 `git rm`＋commit ⇒ 门回 rc=**9**（＝本票修后的真基线读数，见文末追加段）。
- 途中一并实锤一条工具链事实：`git restore --worktree` 不回填 index，checkout 失败后
  `git commit` 会把近空 index 提成新根树（首版克隆 ls-tree=1 事故，已废弃重克隆）。

## 验收（AC）

- [x] **AC#1 修后仪器能对全树跑出真读数**：`bash .scratch/wisp/probes/154/gate-clauses.sh`
  ⇒ rc 从假 10 变真读数（锚 cb2287a5＝9），10 腿假 `git-grep-rc-126` 归零，stderr 零噪音，
  腿数断言 14/14 相符。
- [x] **AC#2 判据一字未动**：LEGS_EXPECT/want/want_n/pair 四堆/聚合逐枚原样
  （`git diff` 只见 `-- "$@"` 两处＋转义两处＋注释）。
- [x] **AC#3 正控两发**：变异红（rc=10，G6 具名）/还原回真基线（rc=9）。
- [ ] **AC#4 裁决者≠实现者**：9 腿真红与射程语义变化（本文件「射程语义变化」节）交
  154/158/161/171/178/195 家族裁决者复核；重登基线（若有）须由裁决方按 171 r3 规矩同格现量。
- [x] **AC#5 回归件复跑**：`flip-declaration.sh` 全电池跑毕（锚 cb2287a5，约 75 分钟 8 发）＝
  **RED rc=1，finding 路径**（工具判语逐字：「A declaration that stopped matching today is a
  finding for the next slice, not a knob to turn」）。分解：① 每一发翻转都真实移动聚合、
  还原发与基发两读一致（rc=9）＝机制完好；② stale 发 3 FAIL＝G5neg 实测 14 ＞ 抬到 9 的
  基线 ⇒ 走 BAD 不走 SHR——与 9 腿声明过期是同一件事，正是本票交裁决的 finding 本体，
  不是仪器坏。日志＝`probes/154/r20261108/30-flip-final.txt`。
- [ ] **AC#6 受影响票的 AC 文字刷新**：195 AC#5 的「BAD 腿名册只 G6neg」一句按裁决结果改写
  （归票主，本腿不碰别票票面）。

## Progress log

- 2026-11-09 09:5x 修程落地：`-- "$@"` 两处＋反引号转义两处＋头注 411 块；判据零改动；
  真读数 rc=9（14 腿逐枚具名见证据件）；正控红/还原对照成立；只 commit 不 push。
- 2026-11-09 11:0x flip 电池终局补录：RED rc=1（finding 路径，工具判语逐字在案）；
  翻转移动聚合＋还原两读一致＝机制完好；stale 3 FAIL＝基线过期 finding 本体。
  AC#5 收勾；AC#4/AC#6 仍留给裁决方与 195 票主。
- 2026-11-09 11:1x **AC#6 落地**（编排者派 195 侧写腿执行，195 票面只追加不改写）：195 AC#5 追加具名更新一行（修后真读数 rc=9、BAD 名册九腿＝G1/G1b/G2/G5/G5pos/G5neg/G6neg/G7/G7neg、溢出期原句「只 G6neg」作废）＋195 Progress log 一行；九腿具名命中零涉 `internal/config/**`（scoped 复量六词表全 0）⇒ 195 射程内 0 红、判语「AC#5 具备翻框条件，勾选留裁决者」存证据件 `docs/evidence/s1/195-ac5-refresh.md`（新建）。本票 AC#6 勾选留裁决方（本腿非 411 实现腿，不翻框）。
