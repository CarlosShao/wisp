# `30-costs.md` — 表③：甲／乙／丙 三形代价表＋票面必答两问（票 310 `AC#0`，腿 `310-a1`）

本件是**代价表**，不是建议书：每形四问逐格答满（会动哪几枚文件／会不会打红既有钉·点名哪一枚／用户侧表现／失败模式＝最坏情况什么留在盘上），⛔ "应该没问题"。
读数来源：`10-roster.md`（名册与分母）、`20-shapes.md`（三形实测＋甲形两变体实测 (d)(e)）、`logs/30-temp-env-census.txt`（必答 (I) 的尺）。

先摆三条**决定代价的机制现量**（provenance＝仓外 GOROOT `go1.27.1`，本腿逐字读过；⛔ 引用编排者转述）：

- M1 `t.TempDir()` 的根＝`os.MkdirTemp(os.Getenv("GOTMPDIR"), pattern)`（`testing.go:1613`）。`GOTMPDIR` 为空 → 落 `os.TempDir()`。⇒ **GOTMPDIR 只搬 `t.TempDir()`／`os.CreateTemp` 这一族**。
- M2 Windows 上 `os.TempDir()`＝Win32 `GetTempPath`（`os/file_windows.go:284-289`），读 `TMP`/`TEMP`，**⛔ 读 `GOTMPDIR`**。⇒ **`TMP`/`TEMP` 重定向会同时搬走 `t.TempDir()` 和所有逐字读 `os.TempDir()` 的测试**。这两枚变量不是一回事，甲形选哪个决定了 (I) 的答案。
- M3 每个 test 对象只在 `%TEMP%` 顶层建**一枚**父目录（名＝用例名截 64 字节＋`nextRandom()`），同用例第 2 次以后的 `t.TempDir()` 是它里面的 `\001`、`\002` 子目录（`testing.go:1623-1633`）。⇒ **顶层枚数＝死掉的 test 对象数，⛔＝临时目录调用数**；清一枚父目录就带走它那批子目录。

---

## 1. 甲＝测试期把 `TMP`/`TEMP`（或 `GOTMPDIR`）指到一枚仓外专用根，收尾按"本仓可证"处置

| 四问 | 甲-G（`GOTMPDIR` 变体） | 甲-T（`TMP`/`TEMP` 变体） |
|---|---|---|
| 会动哪几枚文件 | 跑测入口的环境变量：`scripts/portable-tests.sh`（本机／CI 主入口，`go test` 调用点＝`:617`、`:652`）＋同源另三枚起 `go test` 的件＝`scripts/winsec-tests.sh`、`scripts/wisp-cli-tests.sh`、`scripts/slo-check.ps1`。**零产码、零 `*_test.go`**（全仓 `func TestMain(`＝**0 枚**，现量尺见 `logs/`，所以想"在每个包起手设根"根本没有现成钩子，只能落脚本）。⛔ 动 `.github/workflows/ci.yml`（票面 `AC#1` 硬约束）⇒ **CI 那一侧本票改不到，只有本机变干净** | 同上四枚脚本；⚠ 另外会连带影响 `go` 工具自己的构建临时根与任何读 `os.TempDir()` 的产码路径（`internal/proc/envfork.go:125` 的测试数据根就是从 `os.TempDir()` 派生的）⇒ 想只搬测试而不搬产码路径，就得动 `internal/proc/envfork.go`，**那是产码、且撞票 119／124 已定的 data-root 纪律，本票边界外** |
| 会不会打红既有钉 | 机制上**不动** `os.TempDir()` 的答案（M2）⇒ 那 12 枚逐字读 `os.TempDir()` 的测试件读数不变。**实测**：本次 (d) 一发在 `GOTMPDIR` 下照样 panic、`%TEMP%` 顶层新增 0 枚。⚠ 未跑＝三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）与 4 枚真窗件（`cmd/wisp/panel_host_windows_test.go:620,878`、`cmd/wisp/panel_resident_windows_test.go:61,447,524`、`cmd/wisp/panel_transport_live_35v2_windows_test.go:186`、`cmd/wisp/panel_host_windows_live_test.go:37`，它们把 profile 根写死成 `os.TempDir()+"wisp-33r1-panel-profile"` 一类**固定名**）——⛔ 整包跑是派单 §3.5 的硬规，这几枚本腿没跑＝欠账 `Q-310-a1-3` | 机制上**同时挪走** `os.TempDir()` ⇒ 点名风险（静态尺已数出、本腿⛔ 跑过）：`internal/proc/envfork_test.go:124-139`（现量断言"测试数据根必须落在本机 temp 经 `EvalSymlinks` 之后那棵树里"，换根后仍自洽，但一旦根与 `%TEMP%` 不同卷就会改判据对象）；`internal/winsec/volume_attribution_126_windows_test.go:303` 与 `internal/winsec/absoluteness_attribution_129_windows_test.go:192`（`parent := os.TempDir()` ⇒ **卷归属换根就换答案**）；`internal/tools/paths_shortname_252_probe_test.go:298`（8.3 短名探测，新根若落在已被 8.3 排除的卷上会翻）；`cmd/wisp/secret_test.go:478,496`。**实测那一半是绿的**：`internal/config` 六枚同包用例（`TestRoundTripLoadMarshalLoad`、`TestRoundTripBytes`、`TestResolvedNeverPersists`、`TestSaveFileRoundTripsThroughLoad`、`TestLoadFileMergesDefaultsAndOverrides`、`TestUnwiredGuardLeavesHonestConfigsAlone`）在 `TMP=TEMP=<仓外根>` 下 `run_lines=10 pass=10 fail=0 skip=0`、rc=0（件＝`logs/30-jia-redirect.txt`）——但注意 0 SKIP 不等于零风险：这六枚本来就不断言根在哪 |
| 用户侧表现 | 无（跑测专用，机主看不到变化）。机主唯一能感知的是：`%TEMP%` 里那批 `Test*` 不再新增，磁盘上换成一枚**具名根**里的残渣 | 同上；额外风险＝机主若在别处也依赖 `TMP`（同机其他 Go 项目、任何读 `GetTempPath` 的程序由我们的脚本继承不到，⛔ 传染），本机跑测期间的临时文件全部集中到那枚根 |
| 失败模式（最坏留在盘上） | ① **根建在 `%TEMP%` 里＝把噪声藏起来而不是清掉**：本腿 (d)(e) 两发的根恰恰就是 `mktemp -d` 默认落在 `%TEMP%` 下，结果 roster 尺（**非递归**）读出 `0 新增`、而残渣各 1 枚实实在在躺在 `<根>/TestRoundTripLoadMarshalLoad…`（逐枚名见 `20-shapes.md` §2）。⇒ 任何"甲形生效"的断言都必须**同时**报根内枚数，⛔ 只报 `%TEMP%` 顶层。② 根落 `D:\tmp\` 之类的持久盘 ⇒ 50 枚／397 MB（表① 实测合计）整体搬家，那枚根成为**新的、无人清的堆**，而且从此⛔ 被任何 `%TEMP%` 尺看见。③ 根落点若图省事指到仓内 ⇒ 直接污染共享工作树（票面与本派单都 ⛔）。④ 最大单枚实测 31,092,893 B（表① `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128…`）⇒ 换根⛔ 减小最坏情况，只换位置 |
| 治不治得了"没人收尾" | ⛔ **不治**。(d)(e) 两发都照样留下 1 枚没人删的目录——甲形治"落在哪"，⛔ 治"谁去删" | 同左 |

## 2. 乙＝起手清扫（下次跑测试先扫自己上次留下的、按可证归属删旧的）

| 四问 | 答 |
|---|---|
| 会动哪几枚文件 | 最小落点＝脚本侧一枚清扫步骤，挂在 `scripts/portable-tests.sh` 起手（`go test` 调用点＝`:617`/`:652` 之前），配一枚**腿私有名册件**（上一次自己建了哪几枚，逐枚名）；同源另三枚（`winsec-tests.sh`／`wisp-cli-tests.sh`／`slo-check.ps1`）要不要一起挂＝编排者裁，⛔ 我裁。零产码／零 `*_test.go`（`func TestMain(`＝0 枚 ⇒ 逐包挂根本无处可挂）。**⛔ 动 `.github/workflows/ci.yml`** ⇒ CI 上不成立，只有本机会扫。若把"删"做进产码（比如某枚 helper 的 defer），就要动 `internal/proc/envfork.go`／`internal/winsec/resolve.go` 这一族 data-root 纪律件，撞票 119／124 已定案 ⇒ 边界外 |
| 会不会打红既有钉 | 点名两枚：① `bash scripts/portable-tests.sh --scope=census` 的 totals 行（尺原文在 `scripts/portable-tests.sh:407`，基线＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`）——起手清扫⛔ 改变用例集合，但同一次 commit 里若顺手动了脚本头部注释/账，票 310 `AC#3` 要求那行**逐字未变**；② 三枚冻结件（见上）不读临时目录归属，⛔ 受删动作影响。**真正的风险⛔ 在钉上而在副作用**：删错了会把**别家正在用的**目录带走，那一形打红的是别人的测试、⛔ 是我们的 CI，现场不可复现 |
| 用户侧表现 | 机主第一次跑测试时会看到一批目录被删（今天＝50 枚里的可证那 49 枚，含 397 MB）。若判据错，删到同机其他 Go 项目的运行中数据 ⇒ 不可逆 |
| 失败模式（最坏留在盘上） | 最坏⛔ "留在盘上"，最坏是**删掉了不该删的**：别家某项目正在写的 `TestXxx<随机>` 目录，名字恰好与本仓某条用例同名（表① 的三形尺只证"这名字在我们仓存在"，⛔ 证"这枚目录是我们建的"），mtime 又落在我们跑测窗口里（表① §3 已经出现"mtime 尺整批失效"的形——09-21 那天本仓 357 笔提交，窗口覆盖每一天）。⇒ 乙形把票面 `:33`"删只按名册逐枚、先经非实现者复核"变成脚本里的自动动作，是最容易越界的一形 |
| 治不治得了"没人收尾" | **治**（在"下一次跑测"这个时刻之后才生效，⛔ 当场）：它清的是上一轮的账，本轮死掉的那枚照旧留到下次 |

## 3. 丙＝**不做**，改成把 `%TEMP%` 下 `Test*` 目录枚数登记进编队起手闸门名册

| 四问 | 答 |
|---|---|
| 会动哪几枚文件 | **一枚**：`.scratch/wisp/probes/orch/gate-roster-304.md`（正典件，跟踪件）在 §1 那四条 `tasklist` 之后补一行尺＋一发现量；另加票 310 的 `AC#4` 判据文本。零产码、零脚本、零 `*_test.go`、零 CI |
| 会不会打红既有钉 | ⛔ 零：那枚件是编队闸门的文本正典，⛔ 被任何 Go 测试读；唯一约束是"改后逐字文本＋一发现量"（票面 `AC#4` 原文）与 `SPEC-12 §4.1` 的"契约变更须人工批准"——本件⛔ 契约（C1–C32／D1–D47）射程，属编队文本，仍应由编排者执笔、⛔ 由我落 |
| 用户侧表现 | 机主看到 `Test*` 继续涨：立案 21 → 10-11 44 → 本腿今天 50（起手）／53（含我自己 3 枚）；一天最多 29 枚（10-11）。磁盘占用实测 50 枚＝397 MB、单枚最大 31 MB。**没有任何自动回收** |
| 失败模式（最坏留在盘上） | ① 残留无限增长且**归属长期悬空**（票面"为什么要紧"③原话），每一次"这是谁留下的"都要重跑一遍表① 那种三形尺；② 增长到机主自己动手在 `%TEMP%` 里做一次**不可复核的手工全删**（那一刀会连着别家一起挨）＝票面 ⛔ 的那一形；③ 噪声淹没信号：闸门数到 53 与数到 300 对"今天有没有新死亡"没有分辨力 ⇒ 名册那行必须**同批带 mtime 分布**（本腿这把尺现成：`logs/10-temp-test-dirs.tsv` 的 mtime 列）才有用 |
| 治不治得了"没人收尾" | ⛔ 不治，也⛔ 打算治：丙的立场是"这是已知残留，看得见就行"。它的价值＝**零风险、一小时能落地**，并把"要不要真修"这件事变成一个有数的决策而不是每次现吵 |

---

## 4. 必答 (I)：甲形会不会改掉本仓测试里"临时目录在 `$env:TEMP` 下"这条既有假设？

**先把件数出来（派单要求"⛔ 没数完不许说不影响"）。尺射程＝`git grep -n -E <pat> HEAD -- cmd internal tools scripts`，件＝`logs/30-temp-env-census.txt`（9495 B，逐把自落 `grep_rc`）：**

| 尺 | 命中行 | 文件枚数 | 其中 `*_test.go` | 其中产码 |
|---|---|---|---|---|
| `os\.TempDir()` | 31 行 | 16 | **12** | 4 |
| `"TEMP"`／`"TMP"` 字面量 | 2 行 | 2 | **2** | 0 |
| `%TEMP%`／`$env:TEMP` | 6 行 | 6 | **3** | 3 |
| `GetTempPath` | 1 行 | 1 | 0 | 1 |
| `GOTMPDIR` | **0 行**（empty above＝0 枚，`grep_rc=1`，stderr 与命中数并排落盘） | 0 | 0 | 0 |
| `t\.TempDir\(\)`（分母） | 446 行 | 163 | 163 | 0 |
| 本仓 `*_test.go` 总数（分母） | — | **387** | — | — |
| `func TestMain(` 枚数 | **0**（empty above＝0 枚） | 0 | 0 | 0 |

**答案（分变体，⛔ 一句话"影响／不影响"）：**

1. **甲-G（`GOTMPDIR`）：⛔ 改掉这条既有假设。** 由 M2，Windows 上 `os.TempDir()` 走 `GetTempPath`、⛔ 读 `GOTMPDIR` ⇒ 那 12 枚逐字读 `os.TempDir()` 的测试件**答案不变**；只有 `t.TempDir()` 一族换根。实测支持：(d) 一发 `%TEMP%` 顶层 0 新增、残渣进 GOTMPDIR 根。
2. **甲-T（`TMP`/`TEMP`）：会改。** `os.TempDir()` 的答案**本身**被换掉 ⇒ 上述 12＋2＋3 枚件的读数全部跟着换根。这不等于一定打红，但**"不影响"这一句必须有读数支撑**：本腿只跑了 `internal/config` 那六枚（10 条 RUN 全 PASS、0 FAIL、0 SKIP，`logs/30-jia-redirect.txt`），**剩下 11 枚读 `os.TempDir()` 的测试件一枚都没跑**（⛔ 整包是派单 §3.5 硬规，其中 4 枚是真窗/GUI 件、本腿也⛔ 能起）。⇒ 这一格我只能交"已证的部分＋具名未测清单"＝`Q-310-a1-3`，⛔ 替编排者说"不影响"。
3. 补一条**已证的构造性风险**（⛔ 依赖读数）：那 4 枚真窗件把 profile 根写成 `os.TempDir()` 拼**固定名**（`wisp-33r1-panel-profile` 等，`cmd/wisp/panel_host_windows_test.go:620,878`）。换根后固定名跟着走 ⇒ 若两枚腿一枚走甲、一枚不走，会同时打开**不同**的 profile 根；这一形⛔ 是"临时目录残留"问题，是甲形的副作用，编排者选甲时要一并裁。

## 5. 必答 (II)：乙形"按可证归属删旧"在同机别家用同样命名时怎么保证⛔ 误删？

**它打算用的判据（逐字段说清哪个字段能把两家分开）：**

| 候选字段 | 能不能分开两家 | 本腿实测到的强度 | **什么时候失效**（具名） |
|---|---|---|---|
| 名字（stem 反查本仓 `func Test…`） | ⛔ 单用⛔ 行 | 表① 三形尺做到 49／50 可证——但那是"这名字在我们仓存在"，⛔"这枚目录是我们建的" | ① 别家恰好有同名用例（Go 用例名⛔ 全局唯一，只包内唯一）；② 我们用例改名／删除后，自己的旧目录**当场退化成不可证**；③ 64 字节截断＋子测试粘连（`20-shapes.md` 机制 §1）会让"剥尾数"这一步在长名字上歧义 |
| mtime 窗口 | ⚠ 只能证伪，⛔ 证真 | 本批 49 枚落在命中文件首提交之后、0 枚矛盾——但**那枚不可证的目录 mtime 正好在我们 357 笔提交的那一天里**，窗口尺对它整枚失效 | 两枚腿／别家同时在跑（票面 43 枚那一发就是在两枚验收腿在跑期间量的，混在一堆⛔ 分）；机器时钟改动；`%TEMP%` 被手工清过导致名册与盘不同步 |
| 命令行反查 | ⛔ 对残渣可用 | 只对**活着**的进程有效（`Get-CimInstance Win32_Process` 那条家族尺，票 304 的用法）；目录⛔ 带来源 | 残渣＝进程已死 ⇒ 这条路天然为空；PID 复用（gate-roster §1 第 2 问已写死"⛔ 拿父 PID 当归属"） |
| 仓属根（把归属编码进路径） | ✅ **唯一构造性分得开的判据** | (d)(e) 两发实测证明"根"能把残渣整体隔离到一个我们命名的目录里 | 但那是**甲形**、⛔ 乙形；且根建在 `%TEMP%` 里时判据退化为"藏起来"（本腿踩到，见 §1 失败模式①） |

**⇒ 乙形若要落地，本腿认为唯一⛔ 越界的判据是三条同时成立（缺一⛔ 删）：**

1. **自证名册**：删只删"上一轮**我们自己**逐枚记下"的那批目录名（＝本腿在 `20-shapes.md` 用的 `comm -13` 那把尺的持久化版），⛔ "扫 `%TEMP%` 下所有 `Test*` 然后按名字过滤"；
2. 该枚名字现在仍**可证本仓**（表① 三形尺任一一形命中）；
3. 该枚**不在任何活着的测试二进制名下**（现量家族尺＋命令行含我们那棵树）。

**这条判据什么时候失效（具名，⛔ 含糊）：**
- 名册件丢失／未写成功（脚本被 `TerminateProcess`，正是本票的死亡形！上一轮记名册之前就先死）⇒ 第 1 条无从成立，实现者最自然的"退化"就是改回按前缀扫——那一形票面 `:12` 与 `:33` 双重 ⛔。**这是乙形的内生矛盾：它最需要在死亡后生效，而死亡恰恰会吃掉它的名册。**
- 两枚腿共写一枚名册 ⇒ 互相吞名（A 腿删掉 B 腿刚记下的）。⇒ 名册**必须腿私有**（带腿 id），而腿 id 同日重派（`a1`／`a1b` 这种）会串，得带时间戳。
- 换机／`%TEMP%` 路径变（本例 `C:\Users\swq\AppData\Local\Temp`）⇒ 名册里的绝对路径整体失效。
- 名册记了、盘上那枚已被机主手工删 ⇒ 删动作应报"不存在"而⛔ "已清"，否则账会越记越绿。

---

## 6. 三形并排的一句话（⛔ 替编排者裁）

- 只有**丙**是零产码、零删除风险、且本腿已把尺备好（`10-roster.md` 那把三形尺＋mtime 分布就是名册那一行的内容）。
- **甲**已实测能移出 `%TEMP%` 顶层，但⛔ 治"没人收尾"，且甲-T 会动那条既有假设（甲-G 不会，代价是只覆盖 `t.TempDir()` 一族）。
- **乙**是唯一真正会减少盘上枚数的形，而它的判据在本票的死亡形下会自失（名册被同一刀吃掉）。
- 本腿⛔ 能替编排者答的取舍，全部登记在 `90-final.md` 的 `Q-310-a1-*`。
