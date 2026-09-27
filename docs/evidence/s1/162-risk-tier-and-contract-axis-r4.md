# 票 162 · AC#5（风险档与门控＝L2）＋ AC#6（契约轴）＋ 派单 §3 那一枚〔现读码〕 — 写码位 `162-r4` 交件件

派单＝`.scratch/wisp/dispatches/2026-09-27-173x-impl-162-r4-risk-tier-l2-and-contract-axis.md`（36 行，全文读过）
票面＝`.scratch/wisp/issues/162-add-the-patch-one-region-tool-instead-of-rewriting-whole-files.md`
锚点＝本程 step-0 现量（下面 §1）。**本程没勾任何 AC 框**，勾归编排者。

## 0. 票面两格与三条 `>` 更正（照抄，原句一字未抹）

AC#5（票面 `:42`）：

> **AC#5 风险档与门控**：`fs.edit` 判定必须与 `fs.write` 同族（L1），且**越界路径必须被 `PathResolver` 拒**；⚠ 不许新增豁免、不许动 `allowlist.txt`。

AC#6（票面 `:43`）：

> **AC#6 契约轴**：本票**只许**动 `docs/PLAN.md` 的 `D34` 那一行（**需 owner 单独批准，且批准原文要落台账**）＋ `internal/tools/**` 实现与测试＋自己证据件。`docs/specs/**`、`internal/risk/**`、`thresholds.go`、golden、`frontend/**`、`design/**` 零字节。

票面 09:4x 的三条 `>`（`:44-47`）——本程按 ①② 办、③ 是本程所有行号的读法：

> **① AC#5 的档位过期**：那里写的"与 `fs.write` 同族（**L1**）"已被普查程 §4.2③ 改判、并已按 **L2** 落进冻结文本（理由＝`R8` 覆盖已存在 → L2，与 `fs.write` **覆盖**那一行同级；不是"同族＝同档"）。⇒ **本票实现按 L2 做**，AC#5 剩下的判据（越界路径必须被 `PathResolver` 拒、不许新增豁免、不动 `allowlist.txt`）**原样有效**。若实现程实测认为 L2 判错，**报回、别自己改成 L1**。
> **② AC#6 的 `docs/specs/**` 零字节已被同批放开一枚**：…… `docs/specs/SPEC-07-tools-and-plugins.md` §3 那张表**已由编排者加了一行 `fs.edit`**（现量 `:43`）；`PLAN.md` 那行在 `:2536`。**除这两行外 `docs/specs/**` 与 `docs/PLAN.md` 仍是零字节**，且**两行都已由我落完 ⇒ 实现程再动它们＝越权**。
> ⚠ **行号会腐坏**：本票插了 5 行进 `PLAN.md`，所以票面/表里出现过的 `PLAN.md:NNNN` **一律按"自己现量"读**，不许照抄。

## 1. step-0 五件

可复算命令与本程量到的值：

| 件 | 命令 | 读数 |
|---|---|---|
| 时间 | `date '+%Y-%m-%d %H:%M %z'` | step-0 于 09-27 17:3x；收尾复量＝`2026-09-27 17:47 +0800` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev`（全程未换，未用 `switch`/`checkout`/`merge`/`rebase`/`reset`/`stash`/`worktree`/`clean`） |
| HEAD | `git rev-parse HEAD` | step-0 落在 `f1b99a70`（派单那一枚）与 `4cee47f2`（同树另一程 `171-r2`）之间；**恢复凭据**＝`git log --format='%h %ci %s' -8`＋`git reflog -6`：`f1b99a70` 17:27:13＝派单、`4cee47f2` 17:33:13＝171-r2。本程第一枚提交 `7a41af51` 的父＝`4cee47f2`。契约轴凭证（§4）用 **`f1b99a70 HEAD`** 这个最宽也最保守的窗口，逐枚点名，非本程的那一枚单独标出 |
| 写面基线 | `git status --porcelain -- internal/tools/ internal/risk/ cmd/wisp/` | 空 |
| 门禁基线 | `go test -count=1 ./internal/tools/` | **本程没能出示这一遍的盘上读数**（见 §7 自报）；改前基线一律引用台账现量：票面 16:5x／17:2x 两节的 `顶层 PASS=101／FAIL=0／SKIP=0、RUN=150（含子测试）、名册 94→101` |

⚠ **口径**：本仓两把尺——`RUN=` 含子测试、`--- PASS` 只数顶层。本程所有"几枚"都按**顶层**报，并在旁边附命令。

## 2. 本程没测什么

- **非 Windows 那一支**：`internal/tools/fs_edit_ac5_sweep_r4_windows_test.go` 带 `//go:build windows`，其它平台**根本不进编译** ⇒ 未经验证，没拿整包绿当它响过。`staging_live_other.go` 自己写着本平台"暂存文件一律不清扫"（fail-closed），那一半不是本程的判语。
- **"路径无法规范化 ⇒ R2 fail-closed"这一腿在 `fs.edit` 侧没钉**，且钉不了：现量证明空白 path 在到达法官之前就被 `bridge.go` 的 `stringValues`（`:1020` TrimSpace）＋`pathArgs`（`:1004` 跳空串）丢掉了，法官看到的是 `rules_hit=[R1]`、`in_allowlist_scope=true`（"没有路径可判"，不是"判成在范围内"）。这一腿的真尺在别处：`internal/risk/assessor_test.go:172`（规则单元）＋`internal/tools/bridge_junction_windows_test.go:232/306/459`（reparse 形状，含 `fs.write` 走 junction 那一发）。本程不另造第二把尺（票面 AC#4b 的原话）。
- 并排读数只钉了 **`fs.write` ↔ `fs.edit`** 这一对（派单点名的两枚）；`fs.trash`/`fs.move`/`fs.delete` 未并排。
- L2 卡片文案／面板侧一律未经真人：本程的"人"是 `gateSpy`。判语只到"走到哪一条通道、卡片上带的 rules_hit/reason/审计行是什么"。
- `tools/d22scan` 名册**没有本程改前的盘上快照**可比（本程没动那块地），只比了四数与台账记录值（§6）。
- M-2 变异只摘 `fs_write.go:288` 那一枚 `reclaimStaging(parent, se)`，**没摘 `:511`** 那一枚（另一条写通路，不是本票的射程）。

## 3. AC#5 — 风险档与门控

### 3.1 档位现量依据（两枚号都是本程自己量的，没照抄派单）

```
git -c core.quotePath=false grep -n "fs\.edit" -- docs/PLAN.md docs/specs/SPEC-07-tools-and-plugins.md
docs/PLAN.md:2536:| **`fs.edit`** | **字面定位替换已存在文件里的一小段** | **L2** | ... **R8 覆盖已存在 → L2**，与上面 `fs.write` 覆盖行**同级** |
docs/specs/SPEC-07-tools-and-plugins.md:43:| **`fs.edit`** | 字面定位替换已存在文件里的一小段 | **L2** | fs.write | S3 | 票 162 新增，镜像 `PLAN.md` D34；... |
```

⇒ 两枚冻结件都写着 **L2**，票面 09:4x 更正①与之一致；**本程没改这两行一个字**（凭证见 §4）。
码侧一致性（前一程已落、本程只复核）：`internal/tools/fs_edit.go:390` `Declared: risk.L2`、`:391` `PathParams: []string{"path"}`、`:392` 无 `Facts` 钩子（R8 的结论写在声明里，不靠每次现算）。**本程一行生产码都没动。**

### 3.2 落点（新增两枚测试文件、四＋一枚用例，全在未修码上跑绿）

`internal/tools/fs_edit_ac5_gate_r4_test.go`（提交 `7a41af51`，gofumpt 形状收尾 `9c00bb33`）：

1. `TestFSEditAndFSWriteShareTheOutOfScopeVerdict` — **并排读数**（派单 §1 硬判据"不许只用一句注释宣告同族"）。
2. `TestFSEditRoutesToApprovalNotTheL1Window` — 钉的是**通道**：范围内 `fs.write` 新文件走 L1 预执行窗（window=1/approval=0），`fs.edit` 走 L2 审批队列（window 不增）。
3. `TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord` — 判定来自声明给法官的那枚 path，不是工具自报。
4. `TestFSEditRefusesAPathC26CannotCanonicalize` — 第二道 fail-closed 腿（工具侧），并按 §2 说明法官侧为什么钉不得。

### 3.3 并排读数（同一枚越界样本、同一座桥、真 C26＋真 C19，`-v` 原样）

命令：`go test -count=1 -v -run 'TestFSEditAndFSWriteShareTheOutOfScopeVerdict' ./internal/tools/`

```
并排读数 fs.write: level=L2 rules_hit=[R1 R2 R8] window=0 approval=1 class="user_rejected" reason="R1: 工具声明为下界（L1）; R2: 目标路径在授权目录之外: C:\...\002\note.txt; R8: 不可逆操作（覆盖已有内容）"
并排读数 fs.write: 回执="L2 审批未通过，已拒绝执行"
并排读数 fs.edit:  level=L2 rules_hit=[R1 R2]    window=0 approval=2 class="user_rejected" reason="R1: 工具声明为下界（L2）; R2: 目标路径在授权目录之外: C:\...\002\note.txt"
并排读数 fs.edit:  回执="L2 审批未通过，已拒绝执行"
并排审计行 fs.write：tools: call kind=refused task=task-1 corr=corr-1 tool=fs.write risk=L2 decision=reject outcome=error rules_hit=[R1 R2 R8] in_allowlist_scope=false reason="…"
并排审计行 fs.edit：tools: call kind=refused task=task-1 corr=corr-1 tool=fs.edit risk=L2 decision=reject outcome=error rules_hit=[R1 R2]    in_allowlist_scope=false reason="…"
```

两边**同为 L2、同含 R2、同走审批通道、同 `kind=refused decision=reject outcome=error in_allowlist_scope=false`、退码形状同为 `IsError`＋`ErrorClass=user_rejected`＋回执文案逐字相同**，且目标字节未动、目录里无暂存件（`noStagingFilesLeft`）。
**唯一差异本程不掩盖**：`fs.write` 多一枚 R8——因为它挂了 `writeFacts` 钩子现场算"覆盖已存在"，`fs.edit` 没挂（R8 已在声明里）。这正是票面 10:2x 编排者裁过的那条（"L2 不按工具名走"），也是"同族≠同规则集"的具体形状，派单 09:4x 更正①那句话（"不是同族＝同档"）在读数上有第二层意思。

### 3.4 越界那发的两向 ＋ 摘掉哪一行会红

- **向①（在授权树之外、可规范化）**：`fs.edit` 判 L2、`rules_hit` 含 R2、审计行 `in_allowlist_scope=false`，人在卡片上按"拒绝"⇒ 回执 `user_rejected`、目标 12 字节原封不动、无暂存件。**今天（未修码）就响**——本程是补断言，不是加功能：`fs.edit` 天然经过 `bridge.Execute` 的 C19 判定（`bridge.go:272` Assess ⇒ `route`），因为它与 `fs.write` 注册在同一座桥上、`PathParams` 都是 `["path"]`。派单 §1 那句"这一格的正解可能是一枚断言而不是新功能"——**实测就是这一形，本程没有为了有活干造改动**。
- **向②（无法规范化）**：见 §2，法官看不见该形状（`rules_hit=[R1]`），响的是工具自己那一行（`fs_edit.go:110` ⇒ 回执 `缺少 path 参数`、`ErrorClass=tool`）。
- **M-1 变异（摘掉 `fs_edit.go:391` 的 `PathParams: []string{"path"}` ⇒ 换成 `nil`）**，`-overlay` 台件＝`.scratch/wisp/probes/162/r4/m1.json`＋`fs_edit_m1.go`（副本内该符号 0 命中，防静默 no-op），原始红句 `.scratch/wisp/probes/162/r4/logs/m1.log`：

```
--- FAIL: TestFSEditAndFSWriteShareTheOutOfScopeVerdict (0.01s)
    fs_edit_ac5_gate_r4_test.go:176: fs.edit out of scope: rules_hit=[R1], want R2 - 越界判定必须由 C26 出，不是由本工具自报
    fs_edit_ac5_gate_r4_test.go:195: 同族判失：R2 必须两边都在，got [R1 R2 R8] / [R1]
    fs_edit_ac5_gate_r4_test.go:202: wanted an audit line for fs.edit carrying [risk=L2 R2 in_allowlist_scope=false], got none in [... tool=fs.edit risk=L2 decision=reject outcome=error rules_hit=[R1] in_allowlist_scope=true ...]
--- FAIL: TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord (0.01s)
    fs_edit_ac5_gate_r4_test.go:326: the card for an out-of-scope edit carried rules_hit=[R1], want R2 - 没有它就没有「越界」这一判
    fs_edit_ac5_gate_r4_test.go:329: the reason shown to the human must name the allowlist verdict, got "R1: 工具声明为下界（L2）"
```

⚠ **本程要替裁决者先记下的一处形状**：M-1 之下 `fs.edit` **仍是 L2、仍走审批通道、仍被拒**——档位与通道都不会红，红只长在 `rules_hit` 与审计行上（`in_allowlist_scope` 从 `false` 翻成 `true`＝"没东西可判"被写成"判成在范围内"）。所以"档位字符串"不足以当 AC#5 的凭据，本程把规则集与审计行也钉住。**另外两枚用例（通道枚、空白 path 枚）在 M-1 下仍绿**，本程如实分开，没拿"整包红"当读数。

### 3.5 有没有动 `internal/risk/**`

**没有，一字节都没有**（`git -c core.quotePath=false diff --numstat f1b99a70 HEAD -- internal/risk` ⇒ 空，命令与输出见 §4）。本程**没走**"要动那块地才能到 L2 ⇒ 停手上报"那一支：档位由前序程落好的 `FSEditDecl` 声明达成，越界判定复用共用的 C19 R2。未新增豁免、`tools/d22scan/allowlist.txt` 零字节。

## 4. AC#6 — 契约轴凭证（"我没越界"，不是"我改了什么"）

命令（⚠ 两阶段写法，别用 `X..Y` 接管道——本程第一次那么跑出了空输出，改成 `git diff --numstat <A> <B>` 才拿到读数；这条教训记在 §7）：

```
git -c core.quotePath=false diff --numstat f1b99a70 HEAD
834	0	.scratch/wisp/probes/171/r2/logs/gate-pre-fix.txt
278	0	.scratch/wisp/probes/171/r2/r2-run-legs-and-subtraction.sh
398	0	internal/tools/fs_edit_ac5_gate_r4_test.go
78	0	internal/tools/fs_edit_ac5_sweep_r4_windows_test.go
```

逐枚点名：
- 前两枚＝**同树另一程 `171-r2` 的提交 `4cee47f2`**（不是本程的，本程未动、未提交、未评论）。
- 后两枚＝本程的全部写面，**删除列 0**（纯新增文件，零改写）。
- 本程第三枚 `9c00bb33` 是 gofumpt 形状收尾，`--numstat` 为 `8 4`：那 4 枚"删除"是两处 composite literal 的**换行重排**（`gofumpt -d` 全文只有那 4 个 `+`/`-` 行对，见 §6 命令），`cat -A` 复核后**不是内容行被丢**：

```
gofumpt -d internal/tools/fs_edit_ac5_gate_r4_test.go
@@ -138,8 +138,10 @@ 		…got = append(got, reading{"fs.write", …   →  reading{\n\t"fs.write", …
@@ -164,8 +166,10 @@ 		…got = append(got, reading{"fs.edit", …    →  同上形状
```

契约轴各枚号一律 0 行（命令＋空输出）：

```
git -c core.quotePath=false diff --numstat f1b99a70 HEAD -- docs/PLAN.md docs/specs internal/risk internal/panel frontend design .scratch/wisp/probes/154 .scratch/wisp/probes/161 .scratch/wisp/probes/162/v1 tools/d22scan .github thresholds.go
（空＝0 行）
```

⇒ `docs/PLAN.md` 的 D34 那一行（现量 `:2536`）与 `docs/specs/SPEC-07-tools-and-plugins.md` 那一行（现量 `:43`）**本程一字未动**（编排者已落完，动＝越权）。

## 5. 派单 §3 那一枚〔现读码〕的处置＝**钉了**

凭据旧形＝〔现读码 `fs_write.go:288`、`fs.edit` 经 `fs_edit.go:271` 走同一枚〕。本程现量两枚号都还在：`fs_edit.go:271` 确为 `res := t.d.stageAndRename(...)`；`fs_write.go:288` 确为 `reclaimStaging(parent, se)`（⚠ 同文件 `:511` 还有第二枚，不是本票射程）。

新枚 `internal/tools/fs_edit_ac5_sweep_r4_windows_test.go`（`//go:build windows`、文件内零 `t.Skip`；提交 `27f0cb8a`）：

```
go test -count=1 -v -run 'TestFSEditWriteReclaimsADeadWritersOrphan' ./internal/tools/
--- PASS: TestFSEditWriteReclaimsADeadWritersOrphan (0.46s)
    读数：fs.edit 一次成功写盘带走死进程孤儿 .wisp-tmp-f725d3bd-38528-4242；台账="清扫 C:\Users\swq\AppData\Local\Temp\…\001 里本程序留下的历史暂存文件（上次写盘被中断的残口）：清扫 1 个，跳过 0 个；…"
```

断言三件：(a) `plantOrphan` 种下的**可归因、创建者已死**孤儿在 `fs.edit` 那次成功写盘后 `os.Stat` 必为 NotExist；(b) 该次写盘自己不留暂存件（`stagingFiles==0`）；(c) 台账里出现"清扫"那一步（外部可读，不是注释）。
没另造尺：helper 全复用票 73 那族（`deadPID`/`plantOrphan`/`stagingFiles`），只把触发清扫的写盘换成 `fs.edit`；"只扫自己的、活进程不动"那一半仍由 `TestSweepReclaimsOnlyItsOwnStagingFiles`（`fs_staging_windows_test.go:131`）钉着，可复算断言就是它 `:159` 那句 `if _, err := os.Stat(orphan); !os.IsNotExist(err)`。

**M-2 变异**（摘 `fs_write.go:288` 那一行；台件 `m2.json`/`fs_write_m2.go`/`logs/m2.log`）⇒ 当场红：

```
--- FAIL: TestFSEditWriteReclaimsADeadWritersOrphan (0.11s)
    fs_edit_ac5_sweep_r4_windows_test.go:60: fs.edit 的写盘没扫掉死进程留下的孤儿 …\001\.wisp-tmp-84410d2f-68776-4242（err=<nil>）⇒ 清扫器根本没跑在这一侧
```

⚠ **持续性断言的时限**（162-v2 抓过 r3 的那处）：本枚说的是**用例窗内**那次写盘之后孤儿消失；它本来也活不到窗尾——`t.TempDir()` 在清理阶段会带走整个目录。本程**零删除命令**，`.wisp-tmp-*` 与两份变异副本都只建不删，落点在 `.scratch/wisp/probes/162/r4/**`（本程结束时留在盘上，直到有人回收；本程不会带走）。

## 6. 门禁四数＋名册差集＋gofumpt 现跑（逐包 scope，禁全仓 `./...`）

```
go test -count=1 -v ./internal/tools/            → rc=0
  grep -c '^=== RUN   Test'      = 155      （含子测试那把尺）
  grep -cE '^--- PASS'           = 106      （顶层）
  grep -cE '^--- FAIL'           = 0 ; grep -cE '^--- SKIP' = 0 ; grep -c 'panic' = 0
  原始件＝.scratch/wisp/probes/162/r4/logs/tools-verbose-final.txt（改前一遍＝tools-verbose-after.txt，同 106/0/0/0）
名册两向 comm（声明层，改前树 vs 本遍跑）：
  git -c core.quotePath=false grep -h -E '^func Test' f1b99a70 -- internal/tools/ | sed -E 's/^func Test([A-Za-z0-9_]+).*/Test\1/' | sort   → 101 枚
  grep -E '^--- (PASS|FAIL|SKIP)' …/tools-verbose-final.txt | sed -E 's/^--- [A-Z]+: ([^ ]+).*/\1/' | sort                              → 106 枚
  comm -23 = 缺 0 枚 ; comm -13 = 多 5 枚（逐枚：TestFSEditAndFSWriteShareTheOutOfScopeVerdict / TestFSEditOutOfScopeIsJudgedFromTheDeclaredPathNotTheToolWord / TestFSEditRefusesAPathC26CannotCanonicalize / TestFSEditRoutesToApprovalNotTheL1Window / TestFSEditWriteReclaimsADeadWritersOrphan）
  ⚠ 没有用例 panic，故同包读数没被吞（panic=0）。106 = 台账 101 ＋ 本程 5，与 16:5x／17:2x 两节的顶层口径对得上。
sh scripts/d22scan.sh                            → rc=0
bash tools/d22scan/runtests.sh -C tools/d22scan ./... → rc=0，四数＝PASS=34 FAIL=0 SKIP=0 RUN=76 '[no tests to run]'=0
  （与票面 16:5x／17:2x 记录的基线 34/76 逐枚相等；原始件＝.scratch/wisp/probes/162/r4/logs/runtests.txt；本程未动 tools/d22scan，名册改前快照缺，见 §2）
export PATH="$PATH:$(go env GOPATH)/bin" ; gofumpt --version → v0.12.0 (go1.27.1)
gofumpt -l（甲形＝已跟踪 .go ＋ 本程新增两枚）    → 0 行
  已跟踪 .go 现量＝`git ls-files '*.go' | wc -l` = **570 枚**（⚠ 台账里 561／566／554 都是更早时刻的真值，本程 570＝现量，含本程两枚）
```

未跑的尺（派单点名禁止，不是遗漏）：`probes/154/gate-clauses.sh`、`probes/161/r6/flip-declaration.sh`——那一族此刻正被 `171-r2` 改，本程不读别人的半成品。

## 7. 被拒／没成功的调用（发生在取数之前还是之后）

- **被权限系统拒绝：0 次。**
- 本程自己失手三处，逐条＋时点：
  1. `git commit -F - -- internal/tools/<新文件>` 一步式对**未跟踪新文件**报 `pathspec … did not match any file(s) known to git`（取数之前，尚未提交任何东西）⇒ 改成同一命令行里 `git add -- <显式路径> && git commit -q -F - -- <同一路径>`（显式 pathspec、未用 `-A`/`.`/`-a`），随后的提交全部成功。**这条是本程对派单 §5 那条规矩的第一次实测：一步式 `git commit -- path` 只覆盖已跟踪文件。**
  2. 变异台件 M-1 的第一版 python 断言写错（`assert head.count(line)==1`——把 `fs_write.go` 里 `FSWriteDecl` 的缩进当成同一文件的上文）⇒ 断言直接拦下、**没写出变异副本、那次没跑任何判语**（取数之前）。改正为"tail 恰 1 命中＋head 0 命中"后才跑，红句见 §3.4。
  3. M-2 第一版 `assert s.count(line)==1` 失败（真值 2：`fs_write.go:288` 与 `:511`）⇒ 也是**取数之前**被拦，改成按行号只摘 `:288` 后才是 §5 那条红。
  4. `git diff --numstat A..B | cat -A | sed | head` 跑出**空输出**（原因未查明）⇒ 本程**没有把那个空当成"零改动"的凭据**，改用 `git diff --numstat A B` 重取，两份读数都留在本文件 §4。
  5. step-0 那一遍 `go test -count=1 ./internal/tools/` 的**输出没落盘**，本程不引记忆里的数当基线 ⇒ 改前基线一律引台账（16:5x／17:2x），本程自己的数一律现跑现贴（§6）。这是本程的缺陷，不是编排者的。
- 一次编辑告警（`Edit` 之后提示文件态可能已变）：发生在取数之前，随后的 `sed`/`gofmt` 复核确认改动已落地，无判语影响。

## 8. 删除命令次数＝**0**

全程未跑 `rm`/`del`/`Remove-Item`/`git clean`/`git checkout .`/`restore`；未删任何盘上件。产出的临时件（`.scratch/wisp/probes/162/r4/**`：`m1.json`、`fs_edit_m1.go`、`logs/m1.log`、`m2.json`、`fs_write_m2.go`、`logs/m2.log`、`logs/tools-verbose-after.txt`、`logs/tools-verbose-final.txt`、`logs/runtests.txt`、`logs/gofumpt-after.txt`、`logs/gofumpt-after2.txt`、`_names-r4.txt`、`_names-step0.txt`、`_tracked-go.txt`）**只建不删，本程结束时留在盘上**（本程不会在窗尾带走它们，因为它们不在 `t.TempDir()` 下）。

## 9. 伪授权两栏

- **本程收到的"要动契约"类指示：0 条。** 出处＝派单 §2 与票面 09:4x `>` ②（`git -c core.quotePath=false grep -n` 查派单文件 36 行＋票面 `:44-47`）：两行冻结文本"已由编排者落完 ⇒ 实现程再动它们＝越权"⇒ 本程对 `docs/PLAN.md`/`docs/specs/**` 一字节未动（§4 凭证）。
- **本程把自己读过的话当"批准"用过：0 次。** 派单里每一条前提本程都现量：`fs_edit.go:271`／`fs_write.go:288` 相符（§5）；"L2 已落进两枚冻结件"相符（§3.1，现量 `:2536`/`:43`）；**"越界路径必须被 `risk.PathResolver` 拒"这一条不完全相符**——`PathCanonicalizer.Canonicalize`（`paths.go`）不做范围拒绝，只做规范化＋reparse 拒绝，越界判定实际由 C19 的 R2（`internal/risk/rules_gateway.go:45` 读 `InAllowlist`）达成并停在审批通道。本程按**实测形状**写判据、没为了对上派单那句话去改判据或改测试（§3.3／§3.4）。

## 10. 凭据值零抄录

本程未读、未贴、未提交任何 API 密钥／令牌／DPAPI 密文；出现的字符串只有测试自造的夹具路径与 `"ALPH"`／`"BRAVO"` 一类内容。台账里那条"明文 API 密钥"禁令本程无需绕过任何东西（`sh scripts/d22scan.sh` rc=0 即含该扫项，未新增豁免）。

## 11. `next=`

- **本程交回三枚提交**：`7a41af51`（AC#5 四枚用例）／`27f0cb8a`（派单 §3 清扫链一枚，升成用例钉）／`9c00bb33`（gofumpt 甲形归零）。逐枚 `git log --oneline -1`＋`show --name-only` 的原样输出在本文件 §3.2／§5／§6 对应小节与提交信息内。
- **AC#5／AC#6 两格一律未勾**，等非实现者表（`162-v3`）。裁决者建议先看三处：① §3.3 里 `fs.write` 多一枚 R8 而 `fs.edit` 没有——那是"同族"的合法差异还是缺口，归它判；② §3.4 M-1 之下档位与通道**都不红**，红只在规则集与审计行——这条判断（"档位字符串不足以当凭据"）值得复算；③ §2 里"空白 path 到不了法官"这一形状（审计行写成 `in_allowlist_scope=true`）是不是另一枚该记账的缺口，本程判它**不归本票**（尺在 `assessor_test.go:172` 与 `bridge_junction_windows_test.go`），可复核。
- 票面 Progress log 一节由本程追加（最后一步，写前 `git status --porcelain -- 票面` 为空，已核）。

### 11.1 本件自身的提交（原样输出，取数＝`git log --oneline -1` 与 `git -c core.quotePath=false show --name-only HEAD`）

```
b542bc1c 162-r4 AC#6: 契约轴凭证＋AC#5／§3 交件件落盘（docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md）

docs/evidence/s1/162-risk-tier-and-contract-axis-r4.md
```

同一批的另外三枚（本文件 §3.2／§5／§6 引用的就是它们）：

```
7a41af51 162-r4 AC#5: fs.edit 的越界判定钉在 C26 的 R2 上，并与 fs.write 并排读数
internal/tools/fs_edit_ac5_gate_r4_test.go
27f0cb8a 162-r4 派单 §3: 把「fs.edit 的写盘会不会带走死进程留下的暂存残件」从〔现读码〕升成用例钉
internal/tools/fs_edit_ac5_sweep_r4_windows_test.go
9c00bb33 162-r4 AC#5 收尾：gofumpt 甲形归零（只动 composite literal 的换行形状，零判语改动）
internal/tools/fs_edit_ac5_gate_r4_test.go
```

父链核对：`git log --format='%h parent=%p' -1 7a41af51` ⇒ `7a41af51 parent=4cee47f2`（同树另一程 `171-r2` 的提交，本程未动它）。
本小节之后的那一枚提交（票面 Progress log ＋本小节同批）的号本文件里预填不了——**它落号之后由裁决者或编排者复算，本程不引未跑出来的读数**。

