# 证据件：gate-clauses.sh argv 溢出修程（票 411，2026-11-09）

实现腿留存。裁决者≠实现者；本件所有读数带可复跑命令。锚点逐节具名。

## 0. 一句话

共享门禁仪器 `.scratch/wisp/probes/154/gate-clauses.sh` 因 `run()`/`pair()` 内 `-- $@` 未引号，
把 glob pattern 摊成 1084 词／39,836 字符的文件清单喂 git grep，超 Windows 32767 命令行上限
⇒ exec rc=126 ⇒ 14 腿里 10 腿假 BAD（恒 rc=10）。修法＝`-- "$@"`（判据零改动）。
修后仪器第一次对全树跑出真读数：rc=9（9 腿真红，逐枚具名，待裁决）。

## 1. 修前读数（锚 9dbd468e3→cda38fd52，2026-11-08 23:4x～11-09 00:0x）

- 命令：`bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **rc=10**；
  `bash -x` 同跑：每腿 `git grep` 后 `meas=git-grep-rc-126`（G1/G1b/G2/G5/G5pos/G5neg/G6/G6neg/G7/G7neg 十腿）；
  G3/G4/G6pos/G7pos 四腿 pathspec 小不溢出，ok。
- stderr 噪音：`line 587: ring: command not found`、`line 590: 改完之后: command not found`
  （双引号 echo 内未转义反引号被当命令替换）。
- 展开规模现量（`probes/154/r20261108/10-measure-expansion.sh`）：
  bash 未引号展开 **1084 词／39,836 字符** ＞ 32767；git 递归 glob 同射程 1209 枚。
- 机制佐证：`bash -x` trace 实锤展开后 exec 失败；小 pathspec 腿（G3 `internal/agent` 等）同脚本同跑正常。

## 2. 选型（三案现跑）

- **甲（落地）**：`-- "$@"`。git grep 原生吃 glob pathspec（修前正控：
  `git grep -nEw 'OpenTask|CloseTask' HEAD -- 'internal/**/*.go' 'cmd/**/*.go' …` ⇒ rc=0 有命中）。
- **乙**：xargs 分批。已实现对照件 `13-equivalence-xargs.sh`，但多进程拼名册、每批重喂排除项，
  且批间 stderr 静默（`2>/dev/null` 吞错）导致读数偏小（G7 29 vs 真值 57）——否。
- **丙**：`--pathspec-from-file`。**`git grep` 无此选项**（现跑 ⇒ `error: unknown option
  'pathspec-from-file=…'`，该选项属 add/commit/restore 家族）——派单对丙的判据不成立，否。

## 3. 修后读数（锚 cb2287a5，2026-11-09 09:4x）

- 命令：`bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **rc=9**；stderr 空；腿数断言
  名册=14 声明=14 记账=14 缺腿=0 空头声明=0。单遍耗时 6m17s（宽射程名册变大，
  `comment_only_lines` 的 bash 行循环是主成本——记为运行时事实，本程不动）。
- 逐腿读数与具名：见票 411「修后真读数」表（G1/G1b/G2/G5/G5pos/G5neg/G6neg/G7/G7neg 九腿 BAD，
  G3/G4/G6/G6pos/G7pos 五腿 ok）。
- 锚 cda38fd52（一小时前）同读数 rc=9 ⇒ 两锚稳定。

## 4. 射程语义变化定量（`probes/154/r20261108/15-glob-equivalence.sh`，git 原生 :(glob) 语义）

| 腿词表 | 旧射程（:(glob) 两层＝bash 展开精确等价形） | 新射程（:(glob) 递归） | Δ 来源 |
|---|---|---|---|
| G1 ToolRequest\\{ | 2 | 3 | manifestslot/tier2.go（盲区） |
| G1b \\.Execute\\( | 6 | 7 | 同上 |
| G2 OpenTask\\|CloseTask | 11 | 11 | 同 |
| G5 OpenScope\\|CloseScope\\|Close() | 228 | 251 | 盲区 4 枚 manifestslot/approval 文件 |
| G6 同 G2 词表 | 11 | 11 | 同 |
| G7 Defer(Named)?\\|(New)?DisposalScope | 57 | 72 | 盲区 manifestslot 生命周期族 |

- 盲区＝第 3 层及以下目录，HEAD 树上带词表命中的文件 **9 枚**：
  `internal/agent/approval/{batch,cardrender,gate,ui}.go`、
  `internal/plugin/manifestslot/{command,install,lifecycle,tier2}.go`（8 .go）＋
  `internal/ball/assets/jarvis-orb.html`（普查词表误入，仪器 glob 只吃 .go，不算）。
- 旧射程还有第二重坏：清单来自**工作树**（旧清单里有 2 枚 HEAD 树上不存在的
  `cmd/wisp/panel_resident_323*_test.go`＝别腿在飞的未提交件；HEAD 树上 369 枚两层以下的
  .go 它全部看不见）。修后清单与内容同源（都来自 $A）＝收掉票 178 随单发现⑤前半条。
- 结论：修后读数＝仪器史上第一次全树真读数；09-27/10-14 登记的 want/want_n 全部是在
  盲区射程下现量的 ⇒ 不可直接比，重登基线归裁决方（171 r3「谁改口径同格核下来」的规矩）。
- 口径注：`11-deep-hits.sh`（锚 9dbd468e3，2026-11-08 23:5x）曾量得盲区命中 0 枚；其后
  树推进（tier2/manifestslot 落码），`15-glob-equivalence.sh`（锚 cb2287a5）为权威读数
  ＝盲区 9 枚带命中。两件并存，以 15 为准。

## 5. 防作弊正控（导出合成树 `.scratch/wisp/probes/154/r20261108/clone2`，只建不删）

1. 变异：clone2 commit `a512cd7a9`（锚 cb2287a5＋`internal/tools/mutation_411_openonly.go`，
   只 OpenTask 不 CloseTask 的生产形状）⇒ 门红：**G6 ok→BAD（②只开=1，具名该文件）**、
   G6neg 1→2、G2 2→4（含 interface 类型行），聚合退码 9→**10**。
2. 还原：同 clone `git rm`＋commit ⇒ **rc=9**，14 腿逐枚与主仓基线一致（`/tmp/gate-revert.txt`）。
3. 途中实锤：`git restore --worktree` 不回填 index ⇒ checkout 失败的 clone 里 `git commit`
   会把近空 index 提成新根树（首版克隆 ls-tree=1 事故，废弃；重克隆用
   `git -c core.longpaths=true clone --shared`）。

## 6. 回归件（probes/161/r6/flip-declaration.sh）

- 电池 09:33 起跑，**宽射程下单发门 ≈6.5 分钟 ⇒ 全程 40+ 分钟**（运行时事实，本程不动）。
- 中期读数（前 6 发：基线＋flip1-5，存档 `.scratch/wisp/probes/154/r20261108/30-flip-interim.txt`）：
  ① 基线发按设计转为 finding——其恢复基线断言（rc=0）在 9 腿声明与实测不符下不再成立，
  而那正是它要报的事（其文件头「a finding, NOT a knob」），不是仪器坏；
  ② 每一发翻转都真实移动聚合：G5neg ring→quiet 后由「枚数漂移 BAD」变「空/非空 BAD」；
  G7 ring→quiet 后仍 BAD（枚数味 6>基线1 顶住）——翻声明退码跟着变、且枚数味独立在线，
  机制完好；方向与预期一致。
- 终局还原发跑毕后由下一程在票 411 AC#5 补一行终读（本腿不干等，不占共享树）。

## 7. 追记（2026-11-09 10:1x）

- 门禁：gofumpt `scripts/` 命中 0；改动件 CR=0（LF）；d22scan clean rc=0
  （examination 面：bans#1-5 internal=387＋cmd=64；ban#8 internal=966＋cmd=244 全数过）。
