# 独立审计：编排者本会话自翻框结案的 4 枚票（180 / 238 / 439 / 440）

- 审计者：**独立审计者（非实现者）**。派单 `.scratch/wisp/probes/orch/20261113/txt/audit-round91-4closures-dispatch.txt`（87 行，已从头读到尾）。
- 性质：**只读**。⛔ 未改任何已跟踪文件、⛔ 未翻任何 AC 框、⛔ 未 commit 别人的东西。本件是唯一新建件。
- 日期：2026-11-13 12:39–13:0x +08。起手 HEAD `8a8673fe8`；共享工作树当时有 3 处**他人**未提交改动
  （`.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`、
  `.scratch/wisp/probes/orch/20261112/txt/346-a2-dispatch.txt`、`frontend/fixtures/composer-states.html`），
  全部在本审计射程之外，一字未碰。
- 工具：bash 全程 `D:\work\soft\Git\bin\bash.exe`；Go `D:\work\base\go\bin\go.exe`；
  staticcheck `D:\work\base\gopath\bin\staticcheck.exe`（2026.2.2 / 0.8.2，与 CI 同版）；
  gh `C:\Program Files\GitHub CLI\gh.exe`。
- 台件：`D:\tmp\wisp\20261113\audit91\`（只建不删）。

---

## 0. 一句话结论（先给）

**4 枚里我只能签 1 枚（439，附 1 处注记）。180 / 238 / 440 三枚的「结案记录」不能按现状签**——
180 的 AC 框记录与票面自己的 11-09 节互相矛盾；238 的 Status 还停在「待派」；
440 的合并本身成立，但它给 CI 的 staticcheck 新增了 2 条 U1000，而证据件写的是「门禁全绿」。
**这 4 枚的实质工作我逐条独立复核，大部分成立**（读者链、A476 否决、7 枚零读者、ACL 合并、成对判据、门禁读数）；
问题集中在**结案记录与框面**，以及一条**派单前提里的事实错误**（CI 从未绿过，见 §4）。

---

## 1. 四枚 × 六件事 矩阵

图例：✅过 ｜ ⚠️缺陷/应注明 ｜ ❌站不住

| 票 | 1 提交真在·内容一致·无改名残留 | 2 HEAD 无重复票号 | 3 证据件在·支撑主张 | 4 AC 框翻得有依据吗 | 5 复跑门禁 | 6 CI 红是否与本票有关 |
|---|---|---|---|---|---|---|
| **180** | ✅ `5da3848e2`＋`bf13237ed` 真在；`bf13237ed` 是**真 rename（R094）**，HEAD 只 1 条路径 | ✅（仓库级，见 §2） | ⚠️ 件在（22885B）支撑读者链与 30→7 普查；但**不支撑框面**（AC#5/#6 文字说满足、框未翻） | ❌ **AC#2/#3 判「不适用」与票面 11-09 节「应翻勾」直接矛盾**；AC#5/#6 框文不一致 | ✅ d22scan clean；路径长度 69/69/0 GREEN、正控 PASSED、HEAD 重放同值 | ✅无关（两枚提交零 `.go` 改动） |
| **238** | ✅ `e05a389e7`＋`31a262a98`＋`673da11df` 真在；A/D 两枚分提交，HEAD 只 1 条路径 | ✅ | ✅ 件在（95 行），STOP 主张有 A476/A1558 支撑（我复核 `:9904`/`:31279`） | ⚠️ **Status 仍是「待派」**，与 `-done` 文件名和结案节矛盾；结案注残留未解析占位符 `{{PASSWORD_8qdy7ovy}}` | ✅ 零实现，门禁读数＝基线 | ✅无关（零 `.go` 改动） |
| **439** | ✅ `3cdc57faa` 真在（2 files +176），HEAD 只 1 条路径 | ✅ | ✅ 件在（172 行），7 枚零读者逐枚有 file:line；#6 精化（`cloudtts.go:140`）复核成立 | ⚠️ **AC#2 在「零甲 ⇒ 无对象」下被打勾，AC 行本身零注记**（注记只在 Progress log/结案节） | ✅ 我复跑 `./internal/config/ ./internal/panel/` 全绿（2.893s / 17.886s） | ✅无关（零 `.go` 改动） |
| **440** | ✅ `b7bc01482` 真在（6 files +414/−92），HEAD 只 1 条路径 | ✅ | ⚠️ 件在（159 行），主体支撑；但**两个数字不复现**（33 处调用→实测 31；35 个文件→实测 42 个 `.go`/503 tracked），且「门禁全绿」不含 staticcheck | ⚠️ AC#1 判甲**理由成立**（我复核依赖方向、`:(` 5→4、成对判据 PASS）；但 **AC#4 门禁读数不全**（缺 staticcheck，而 CI 因本提交 64→66） | ✅ 我复跑 `./internal/models/ ./internal/winsec/` 全绿（11.307s / 34.705s）；gofumpt 空；vet rc=0；d22scan clean | ❌ **有关**：本提交新增 2 条 U1000（`type aceLine`、`var isNumericSID`）；`aceLine` 在**两种 GOOS 下都是死声明** |

---

## 2. 逐件取证（命令逐字＋读数）

### 2.1 第 1 件：提交真在、内容一致、无「只 add 没 delete」的改名

```
$ for s in 5da3848e2 bf13237ed e05a389e7 31a262a98 673da11df 3cdc57faa b7bc01482; do git cat-file -t $s; done
commit   （7/7 全部 commit）
```

`git show --name-status --format=''` 逐枚：

| 提交 | name-status | 与声称 |
|---|---|---|
| `5da3848e2` | `M` 180 票面（旧名）＋`A` 4 件（probes 3＋证据件 1） | ✅ 腿交付 |
| `bf13237ed` | `R094` 180 旧名→`-done` ＋ `A` 439 新票 | ✅ **真 rename**（相似度 94%），HEAD 只 1 条路径 |
| `e05a389e7` | `A` 238 `-done`（91 行） | ⚠️ 只 add 新名（旧名留给下一枚） |
| `31a262a98` | `D` 238 旧名（66 行） | ✅ 补删，**终态 HEAD 只 1 条路径** |
| `673da11df` | `M` 238 `-done`（+2）＋`A` 证据件 | ✅ STOP 腿 |
| `3cdc57faa` | `M` 439 票面（+4）＋`A` 证据件（+172） | ✅ |
| `b7bc01482` | `M` 440 票面＋`A` 证据件＋`M` models 两件＋`M` winsec 一件＋`A` `internal/winsec/aclparse.go` | ✅ 六件，与声称一致 |

**改名残留检查**（派单点名的尺）：

```
$ for n in 180 238 439 440; do git ls-tree -r --name-only HEAD | grep -c "issues/$n-"; done
180: 1   238: 1   439: 1   440: 1
```

⇒ **4 枚各 1 条，无旧名残留**。238 的「A 与 D 分在两枚 commit」在 `e05a389e7` 那一瞬 HEAD 曾同时有两条，
`31a262a98` 补删后终态干净——**这是派单已知的 A1820/A1822 形状，但终态不复现**，判过。

### 2.2 第 2 件：HEAD 里到底有没有重复票号

```
$ git ls-tree -r --name-only HEAD | grep '^.scratch/wisp/issues/[0-9]' | sed 's/.*issues\///' | cut -d- -f1 | sort | uniq -d
（空输出）
$ git ls-tree -r --name-only HEAD | grep -c '^.scratch/wisp/issues/[0-9]'
426
```

⇒ **无重复票号**。本会话修过两次的 17 枚重复路径问题，今天复验**不复现**。

### 2.3 第 3 件：证据件逐件读，哪一段支撑哪个主张

| 证据件 | 字节/行 | 支撑的主张 | 撑不住的主张 |
|---|---|---|---|
| `docs/evidence/s1/180-panel-width-reader.md` | 22885B / 270 行 | §1① 读者链 4 跳（`:231`→`windowOptions`→建窗）、§1② 规格（`PLAN.md:2736/2753`）、§1③ 归属票 255、§2 普查 70 枚 R48/G8/Z7 与 7 枚逐字具名、§3 门禁五条 | **不支撑 AC#5/AC#6 的「已满足」框面**（件里确有门禁读数，但票面框是 `[ ]`）；件 §1 末自陈「真窗 winlive 未跑」——结案也未宣称跑过，一致 |
| `docs/evidence/s1/238-data-root-seal.md` | 14296B / 95 行 | §1 A476/A1558/A1655 三连否决、§2 现量表（proc 零 MkdirAll、logging.go:99 创建即封）、§3 逐格回应、§6 门禁五读数 | 无（STOP 件，零实现，主张与读数自洽） |
| `docs/evidence/s1/439-seven-config-keys.md` | 17306B / 172 行 | §1 七枚逐枚复核（粗尺/真读/差集）、§2 逐枚归因甲乙丙丁、§3 三件摆机主、§5 边界声明 | **不支撑 AC#2 的 `[x]`**（件里明写「零甲 ⇒ 无判据可跑」） |
| `docs/evidence/s1/440-acl-parser-drift.md` | 14541B / 159 行 | §0 射程差表＋真实 dump 逐行对照、§0.3 结构体字段、§0.4 五份 `:(` 、§AC#1 判甲理由、§AC#2 成对判据两发＋变异实证、§AC#5 边界 | **不支撑「Final gates green at the committed state」**（那句只列 models/winsec/gofumpt/d22scan，**不含 staticcheck**，而 CI 因本提交 64→66）；两个数字不复现（见 §5.3） |

### 2.4 第 5 件：复跑门禁（本机现跑）

```
$ sh scripts/d22scan.sh
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=404, bans #1-5 cmd/=77,
  ban #6 frontend/=123, ban #7 internal/tools/=52, ban #10 internal/risk/=21,
  ban #8 design/=39, ban #8 frontend/=123, ban #8 internal/=1005, ban #8 cmd/=276
d22scan rc=0
```

路径长度那把尺（**含 HEAD 重放**）：

```
$ sh scripts/check-path-length-budget.sh --with-self-test          # 正常 index
check-path-length-budget.sh: positive control PASSED
check-path-length-budget.sh: denominator read: 19510 tracked paths
check-path-length-budget.sh: denominator: tracked paths=19510  over-budget=69  covered by roster=69  not in roster=0
check-path-length-budget.sh: bands: over the hat=69 ... roster entries=69
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
normal-index rc=0

$ export GIT_INDEX_FILE=/tmp/aud91 && rm -f /tmp/aud91 && git read-tree HEAD && sh scripts/check-path-length-budget.sh --with-self-test
check-path-length-budget.sh: positive control PASSED
check-path-length-budget.sh: denominator read: 19510 tracked paths
check-path-length-budget.sh: denominator: tracked paths=19510  over-budget=69  covered by roster=69  not in roster=0
check-path-length-budget.sh: VERDICT GREEN ...
head-replay rc=0
```

⇒ **两读同值**（19510 / 69 / 69 / 0 / GREEN，正控 PASSED）。这条过。（`A1822` 的教训已守住。）

---

## 3. ★第 4 件：AC 框翻得有没有依据（逐格结论）

### 3.1 票 180：**AC#2/#3 判「不适用」站不住**

独立读票面原文（`-done` 件 `:22-27`）：

- **AC#2** 原文：`正向落地（若 AC#1 判"规格要求生效"）：把那一跳接上，并给一发常驻判据…`
  ⇒ **它的适用条件是「AC#1 判规格要求生效」**。AC#1 的判词（`:22` 逐字）正是「**规格要它生效**＝`PLAN.md:2743` 逐字 `hot`」。
  **条件成立 ⇒ AC#2 适用**。而它要的两件事今天都在盘上：接线（`panel_resident_windows.go:231`，我复核）＋常驻判据
  （`TestTicket255AssemblyRootGeometrySourceReachesTheWindowOptions`，517→733，我复核文件在场）。
  更关键的是：**票面自己的 11-09 复验节（`:254`）逐字写着**
  「AC#2 正向落地 … 已在盘（`908c6652c`）⇒ **该由非实现者复判后翻勾**，⛔ 再派腿重做」。
  ⇒ **结案把票面自己指定的「翻勾」改成了「不适用」，直接冲突。**
- **AC#3** 原文：`反向判据（与 AC#2 成对，缺一形＝装饰）：不许用"把默认值改掉"来交差…`
  11-09 节（`:255`）逐字：`判据前提变了（读者在）⇒ 该重写成「删掉 :231 那行喂线必须变红」`。
  **腿恰恰跑了这一发**（证据件 §1 AC#3 的 M3：`:231` 改 `return 1,1` ⇒ `panel_geometry_255_test.go:152` 语义红，
  字节级还原 sha256 `C4E79861…F4EF19` 前后一致）。⇒ **重写后的 AC#3 已满足**，结案却判「不适用」、框留 `[ ]`。
- **AC#5/AC#6**：结案节（`:333-334`）写「AC#5：契约轴无改动」「AC#6：门禁已复跑，两包全绿」，
  但框仍是 `- [ ]`（`:26`、`:27`）。**框面与结案文字互相矛盾。**

**判定：应改判。** AC#2 应翻勾并注记「由票 255 AC#4 落地、本票非实现者复核」；AC#3 应翻勾并注记「已重写为摘线突变、M3 实测红→还原绿」；
AC#5/AC#6 应翻勾（或删去结案文字里的「已满足」）。**这不是假绿**（框没勾），但**结案记录是错的**。

### 3.2 票 439：**AC#2 在零甲下打勾 ⇒ 应注明**

- AC#2 原文（`:40`）：`成对判据（对每一枚判"甲"的）…`。本票归因 **0 甲**（我复核 7 枚逐枚读数，与票面无一处矛盾），
  ⇒ **没有可跑对象**，成对判据确实没跑。
- Progress log 与结案节都如实写了「零甲 ⇒ 无成对判据可跑（如实记，非漏跑）」，**但 AC#2 那一行本身零注记**。
  一个只扫框面的读者会读成「成对判据已跑通」。
- 本仓先例（180 的 AC#2/#3、238 的 AC#1-4）在「无对象」时一律**留框不勾**。
- **判定：应注明**（在 AC#2 行尾写「零甲 ⇒ 无对象，Vacuous；反向半边自动满足＝default 零改动」）。
  不是硬违规，但框面与证据不符。

### 3.3 票 440：**AC#1 判甲成立；AC#4 门禁读数不全**

- **AC#1 判甲**：我独立复核了三条硬机制——
  ①依赖方向：`git grep 'CarlosShao/wisp/internal/models' -- 'internal/winsec/*'` ⇒ **0 命中**（反向才有，仅 `acl_sid_121_test.go:46` 一条 test import）；
  ②`:(` 解析器 **5 份 → 合并后 4 份**（`internal/agent/spill_acl_windows_test.go:51`、`internal/config/private_acl_windows_test.go:166`、`internal/observe/logprivate_132_windows_test.go:110`、`internal/winsec/aclparse.go:76`）；
  ③成对判据 `TestAC440TheSharedParserKeepsTheRightsPayload` **我定向复跑 PASS**（`ok internal/models 0.060s`）。
  ⇒ **判甲理由（主流 harness 的「测试辅助上抽」＋价值只有合并买得到）站得住。**
- **AC#2 成对判据**：两发都在（models 3→4 测试函数、winsec 5→5，**零断言删除**；变异实证在件里）。
- **AC#4 门禁**：我复跑 `go test -count=1 ./internal/models/ ./internal/winsec/` ⇒ **`ok models 11.307s` / `ok winsec 34.705s`，rc=0**；
  `gofumpt -l` 空；`go vet` rc=0；d22scan clean。**但 staticcheck 这一格不在票面 AC#4 里、也不在证据件「门禁全绿」里**——
  而它恰恰是 CI 的红格之一，且**因本提交变红 2 条**（见 §4.3）。

---

## 4. ★第 6 件：CI 红的独立复算

### 4.1 读数（`gh run list` / GitHub API，2026-11-13 12:5x 现跑）

```
$ gh run list --limit 100 --json conclusion,workflowName
  ci           failure = 95
  ci           (空/排队中) = 3
  slo-fresh    success = 2
⇒ 最近 100 跑：95 红 / 2 绿 / 3 未完成。
```

**★那 2 枚绿不是 CI。** 它们的工作流名是 **`slo-fresh`**（`.github/workflows/slo-fresh.yml`，ticket 134 AC#3 的排程「防静默死亡」钉，
不是门禁）：

```
43297215531 | success | 2026-11-12T20:08:02Z | 255224a56 | slo-fresh
43175176199 | success | 2026-11-12T07:25:37Z | 6cb5b9f01 | slo-fresh
```

**★`ci` 工作流本身：1124 跑，0 绿。**

```
$ gh api --paginate '.../actions/workflows/362354689/runs?per_page=100' --jq '.workflow_runs[] | .conclusion' | sort | uniq -c
  failure   = 1069
  cancelled = 55
  success   = 0
$ gh api '.../actions/workflows/362354689/runs?per_page=1' --jq '.total_count'
  1124
$ gh run list --workflow ci --branch main --limit 20   ⇒  []（main 无 ci 跑）
分支分布：dev = 1124（100%）
最老一跑：2026-09-19T23:44:42Z，head f088ce3fc，conclusion=failure
```

⇒ **`ci` 从 2026-09-19 建仓跑起来到现在，一次都没绿过。**
派单/A1827 的前提「最近 100 个 CI run 里 94 红 / 2 绿，最后一次绿是 2026-11-12 20:08（`255224a56`）」
**把 `slo-fresh` 排程钉的绿当成了 CI 的绿**；`255224a56` 自己的 `ci` 跑（run `43243081268`，11-12 15:00 UTC）**是 failure**。
「红从 11-12 晚上才开始」也不成立——**红是常态，不是新事件**。

### 4.2 这个红与这 4 枚结案有无关系

- **180 / 238 / 439**：三枚的提交**零 `.go` 改动**（180 两枚只动票面/证据/probes；238 三枚只动票面/证据；439 只动票面/证据）。
  ⇒ **结构上不可能引入新的 CI 红**。
- **440**：**有关，且是净增**。逐字读数：

```
11-12 的 ci 跑（run 43243081268）：staticcheck self-report: ... findings=64 ...
11-13 的 ci 跑（run 43366115702，head b7bc01482）：staticcheck self-report: ... findings=66 ...
11-13 的 ci 跑（run 43366312122，head b404fc651）：staticcheck self-report: ... findings=66 ...
```

对 64→66 的差集（逐行 diff，去掉时间戳/步骤名差异）＝**恰好 2 条，全在 `internal/models/acl_sid_121_test.go`**：

```
##[error]internal/models/acl_sid_121_test.go:53:6: type aceLine is unused (U1000)
##[error]internal/models/acl_sid_121_test.go:62:5: var isNumericSID is unused (U1000)
```

本机独立复算（staticcheck 2026.2.2，与 CI 同版）：

```
$ GOOS=linux   staticcheck ./internal/models/ ./internal/winsec/ | grep -E 'aceLine|isNumericSID'
internal\models\acl_sid_121_test.go:53:6: type aceLine is unused (U1000)
internal\models\acl_sid_121_test.go:62:5: var isNumericSID is unused (U1000)

$ GOOS=windows staticcheck ./internal/models/ ./internal/winsec/ | grep -E 'aceLine|isNumericSID'
internal\models\acl_sid_121_test.go:53:6: type aceLine is unused (U1000)
（isNumericSID 不再出现——它只在 windows 独占测试件 :53 被用，属 GOOS 盲区）
```

⇒ **`aceLine` 在两种 GOOS 下都是死声明**（合并前它喂 `parseACELines` 的返回类型 `[]aceLine`；
合并后 `parseACELines = winsec.ParseACELines` 返回 `[]winsec.ACELine`，别名再没人提）；
`isNumericSID` 是 GOOS 盲区（windows 形状不报）。**证据件「Final gates green at the committed state」漏掉了 staticcheck 这一格。**

- **测试红**：11-12 与 11-13 的 `--- FAIL` 逐字同名同数（test-core 4 枚 `Test52*` ＋ test-windows 2 枚
  `TestReadFailureDoesNotAdoptStat…`/`Test52AbsoluteTempPath…`），**与 4 枚结案无关**，是陈旧红。
- **lint 的另一格**：`SLO nail census (ticket 266)` 在 11-12 与 11-13 都红（`positive control FAILED - 32 control(s) ran, 14 failed`），
  也与本批结案无关（零 `scripts/` 改动）。

### 4.3 结论

**这个红不是这 4 枚结案造成的**（红早于它们，且 3 枚零 `.go` 改动）；
**但 440 给 staticcheck 净增了 2 条 U1000**（64→66），其中 1 条是两种 GOOS 下都该删的死别名。
派单里「CI 是红的」成立，「94 红/2 绿、最后绿 255224a56」**不成立**。

---

## 5. ★A1773 的 GOOS 裁定今天还成立吗：抽 3 枚复验

**裁定原文**（台账 `A1773`，`:35681` 起；派单写的 `:35686` 是块内第 6 行，近似对）：
> 「47 U1000 里 44 条是"linux 构建下真 unused、windows 构建下有主"的 GOOS 盲区，**删了会把 windows 功能搞坏**」
> 「真正主流做法＝staticcheck 按 GOOS 分分母」

**判据**：同一符号在 `GOOS=linux` 的 staticcheck 里被报 U1000、在 `GOOS=windows` 的 staticcheck 里**不被报**，
且能在 `_windows.go` / `_windows_test.go` 里指到消费者。

| # | 符号 | `git grep` 指到的 windows 消费者 | GOOS=linux | GOOS=windows |
|---|---|---|---|---|
| 1 | `internal/ball/jarvisRingInner`（`jarvis.go:17`） | `jarvis.go:120` 所在函数族被 `jarvis_windows.go` 消费（传递死） | **报 U1000** | **不报** |
| 2 | `internal/audio/selectCaptureDevice`（`device.go:77`） | `wasapimic_windows.go:293` ＋ `input_device_180k_windows_test.go:45/64/76/95` | **报 U1000** | **不报** |
| 3 | `internal/tools/argInt`（`system_bridge.go:274`） | `system_backend_windows.go` 12 处（`:889/:906/:986/:1073-1076/:1177/:1289/:1315-1318/:1406`） | **报 U1000** | **不报** |

命令逐字（可复跑）：

```
$ GOOS=linux   staticcheck ./internal/ball/   | grep -F jarvisRingInner
internal\ball\jarvis.go:17:2: const jarvisRingInner is unused (U1000)
$ GOOS=windows staticcheck ./internal/ball/   | grep -F jarvisRingInner
（空）
$ GOOS=linux   staticcheck ./internal/audio/  | grep -F selectCaptureDevice
internal\audio\device.go:77:6: func selectCaptureDevice is unused (U1000)
$ GOOS=windows staticcheck ./internal/audio/  | grep -F selectCaptureDevice
（空）
$ GOOS=linux   staticcheck ./internal/tools/  | grep -F argInt
internal\tools\system_bridge.go:274:6: func argInt is unused (U1000)
$ GOOS=windows staticcheck ./internal/tools/  | grep -F argInt
（空）
```

⇒ **3/3 复验成立：A1773 的 GOOS 裁定今天仍然成立**，抽到的 3 枚确实「linux 报死、windows 有主」，
删它们会摘掉 windows 功能。**这条裁定我背书。**（顺带：440 新增的 `isNumericSID` 也落在这条裁定的射程里——
它正是「linux 报、windows 不报」；但 `aceLine` 不是，它两种形状都死。）

---

## 6. ★缺陷清单（按严重度排序）

### S1【前提性·非 4 枚结案引入】CI 的「最后绿」是误植
- **事实**：`ci` 工作流 **1124 跑 0 绿**（1069 failure / 55 cancelled，全在 `dev`，最老 2026-09-19）。
  派单与 `A1827` 说的「最近 100 个 CI run 里 94 红 / 2 绿，最后一次绿 2026-11-12 20:08（`255224a56`）」，
  那 2 绿是 **`slo-fresh`** 排程钉（run `43297215531`/`43175176199`），`255224a56` 的 `ci` 跑（run `43243081268`）**是 failure**。
- **为什么重要**：它把「红从 11-12 晚上开始、是本轮引入的」这个错误印象写进了台账，
  会引导下一位编排者去追一个不存在的「昨天还好好的」。
- **建议落点**：`A1827` 就地追加更正（不改历史读数），并把「CI 全史 0 绿」写成一句事实。

### S2【440】合并给 CI staticcheck 净增 2 条 U1000，其中 1 条两种 GOOS 都死
- **事实**：`b7bc01482` 后 CI staticcheck **64→66**；新增 `type aceLine`（`acl_sid_121_test.go:53`）与
  `var isNumericSID`（`:62`）。本机双 GOOS 复算：`aceLine` **linux/windows 都报**（真死）；`isNumericSID` 只在 linux 报（GOOS 盲区）。
- **为什么重要**：证据件写「Final gates green at the committed state」并**故意只列 4 项**，读起来像全绿；
  CI 的 staticcheck 这一格因这次合并**更红了**。合并本身是好活，但**交付件必须自报这件回归**。
- **建议落点**：删 `type aceLine = winsec.ACELine`（无人引用）＋在证据件补记 2 条新 U1000 与归属
  （`isNumericSID` 归 A1773 的 GOOS 口径；`aceLine` 是本次合并的垃圾）。

### S3【180】AC 框记录与票面 11-09 节互相矛盾
- **事实**：①11-09 节（`:254`）逐字「AC#2 … **该由非实现者复判后翻勾**」，结案却判「不适用」；
  ②11-09 节（`:255`）逐字「AC#3 … 该重写成『删掉 `:231` 那行喂线必须变红』」，**腿已照此跑过 M3（红→还原绿）**，结案仍判「不适用」、框留 `[ ]`；
  ③结案文字说 AC#5/AC#6 已满足，框却是 `[ ]`（`:26`/`:27`）。
- **为什么重要**：框面是机器可读的验收状态，结案文字是人读的；两者打架时，下一位无法判断这枚票到底收了没有。
- **建议落点**：按 §3.1 翻 AC#2/#3（带注记）＋翻 AC#5/#6，或把结案文字的「已满足」删掉。

### S4【238】Status 仍是「待派」，与 `-done` 和结案节矛盾
- **事实**：`-done` 件 `:3` 逐字 `**Status**：**待派**（编排者 09-29 23:5x 立…）`；结案节（`:71`）写「按机主裁定结案：本条不做」；
  结案节的撤销口令（`:92`）自己写着「把本节与 **Status** 改回去」——**说明 Status 本该改，但 `e05a389e7` 的 diff 里 Status 行是原样搬过来的**。
- **为什么重要**：`A1821` 的教训正是「派单前读 Status 行」；一枚 owner 已否决（A476）的票挂着「待派」，
  就是下一轮误派的种子（本程已经误派过 3 次：`A1823`、`A1774`、`A1655`）。
- **附带**：结案注 `:83` 与证据件 `:5` 逐字残留未解析占位符 `Q-72`（应为 `Q-72 重开`）。
- **建议落点**：Status 改成「done（owner 裁定不做，A476）」之类，并把占位符改回 `Q-72 重开`。

### S5【439】AC#2 在「零甲 ⇒ 无对象」下被打勾，AC 行零注记
- **事实**：AC#2 原文限定「对每一枚判甲的」，本票 0 甲；Progress log/结案节写了「如实记，非漏跑」，但 AC 行本身无注记。
- **为什么重要**：只扫框面的读者会读成「成对判据已跑」。本仓先例（180 AC#2/#3、238 AC#1-4）无对象时留框不勾。
- **建议落点**：AC#2 行尾补「零甲 ⇒ 无对象；反向半边自动满足（default 零改动）」。

### S6【440】证据件两个数字不复现
- `aclSIDs`「33 处调用零改动」→ 实测 **31 处调用**（`b7bc01482` 时 `git grep -c 'aclSIDs(' -- internal/winsec` ＝ 32 行，减掉 `func aclSIDs` 定义 ＝ **31**）。
- 「`icacls` 出现在 35 个文件里」→ 实测 **`git grep -l -i icacls -- '*.go'` ＝ 42 个 `.go`**；全 tracked 文件 ＝ **503**。
- **为什么重要**：都是「印象当计数」的同族错误（`A1817` 记过：「我看了票面，报了 14，实测只有 1」）。
- **建议落点**：改成实测值，或写明计数口径（如「42 个 `.go` 文件提及 `icacls`」）。

---

## 7. 签字结论

| 票 | 实质工作 | 结案记录 | 签字 |
|---|---|---|---|
| **180** | ✅ 读者链、30→7 普查我独立复核成立 | ❌ AC#2/#3 判「不适用」与票面 11-09 节冲突；AC#5/#6 框文不一致 | **不能按现状签**（改框即可签） |
| **238** | ✅ A476 否决我复核成立（`:9904` 原话逐字、`:31279` A1558 勿再派） | ❌ Status 仍「待派」；结案注残留 `Q-72` 占位符 | **不能按现状签**（改一行 Status 即可签） |
| **439** | ✅ 7 枚零读者逐枚复核无矛盾；#6 精化成立 | ⚠️ AC#2 打勾无注记 | **可签**，附 1 处注记 |
| **440** | ✅ 判甲成立、成对判据 PASS、双包测试全绿 | ❌ 新增 2 条 U1000（含 1 条双 GOOS 死别名）；「门禁全绿」漏 staticcheck；两处数字不复现 | **不能按现状签**（删死别名＋补记 2 条 U1000 即可签） |

**一句话**：**4 枚里我只能签 439（附注记）；180、238、440 的「结案记录」不能签，
但三枚的实质工作我认——它们缺的都是小修：180 补/改 4 个框，238 改 1 行 Status，440 删 1 行死别名＋补 1 句读数。**

**撤销口令：「审计回退」**——本件是唯一新建件，删掉它即复原；4 枚票的框面与文字本审计一字未动。

---

## 附：本审计跑过的关键命令（可逐字复跑）

```
# 1 提交与改名
git cat-file -t 5da3848e2 bf13237ed e05a389e7 31a262a98 673da11df 3cdc57faa b7bc01482
git show --name-status --format='' <sha>
git ls-tree -r --name-only HEAD | grep -c "issues/<n>-"          # n=180/238/439/440 → 1
# 2 重复票号
git ls-tree -r --name-only HEAD | grep '^.scratch/wisp/issues/[0-9]' | sed 's/.*issues\///' | cut -d- -f1 | sort | uniq -d
# 3 证据件
ls -la docs/evidence/s1/{180-panel-width-reader,238-data-root-seal,439-seven-config-keys,440-acl-parser-drift}.md
# 4 AC 框
git show HEAD:<票面> | grep -nE '^\- \[[ x]\] \*\*AC'
# 5 门禁
sh scripts/d22scan.sh
sh scripts/check-path-length-budget.sh --with-self-test
export GIT_INDEX_FILE=/tmp/aud91 && rm -f /tmp/aud91 && git read-tree HEAD && sh scripts/check-path-length-budget.sh --with-self-test
go test -count=1 ./internal/models/ ./internal/winsec/
go test -count=1 ./internal/config/ ./internal/panel/
# 6 CI 与 GOOS
gh run list --limit 100 --json conclusion,workflowName
gh api --paginate 'repos/CarlosShao/wisp/actions/workflows/362354689/runs?per_page=100' --jq '.workflow_runs[] | .conclusion' | sort | uniq -c
GOOS=linux   staticcheck ./internal/models/ ./internal/winsec/ | grep -E 'aceLine|isNumericSID'
GOOS=windows staticcheck ./internal/models/ ./internal/winsec/ | grep -E 'aceLine|isNumericSID'
GOOS=linux   staticcheck ./internal/ball/  | grep -F jarvisRingInner
GOOS=windows staticcheck ./internal/ball/  | grep -F jarvisRingInner
GOOS=linux   staticcheck ./internal/audio/ | grep -F selectCaptureDevice
GOOS=windows staticcheck ./internal/audio/ | grep -F selectCaptureDevice
GOOS=linux   staticcheck ./internal/tools/ | grep -F argInt
GOOS=windows staticcheck ./internal/tools/ | grep -F argInt
```
