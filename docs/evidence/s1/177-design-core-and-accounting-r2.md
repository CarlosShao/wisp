# 177 设计核心 ＋ 分账交件表（177-r2 · 零产码 · 写面仅本件）

- 派单：票 177 实现腿（"design core + accounting work"，`Q-61` 已批；锚 `abf2a35a6`）
- 起手五件〔现跑 2026-11-06 12:16 +08〕：`date`＝Fri Nov 6 12:16:14｜分支 `dev`｜HEAD `8879c3d7b1fcf670d89766cde1fb777bb2423bb8`
- **锚漂移自证**：`git diff --stat abf2a35a6..HEAD` ＝ `docs/reports/pending-and-issues.md | 7 +`（编排者台账追加）一枚，`git diff --name-only abf2a35a6..HEAD -- internal/ cmd/ tools/` **零枚** ⇒ 本表引用的产码面与派单锚逐字节同。
- 起手树〔现跑〕：`git status --porcelain -- internal/ cmd/ tools/ docs/` 只含**别腿在飞未跟踪件**——`cmd/wisp/panel_resident_323_windows_test.go`、`panel_resident_assembly_323_test.go`、`ticket123_replyleg{,_other,_windows}_test.go`、`internal/tools/ticket23_image_channel_test.go`、`docs/evidence/s1/{195-r1-direction-verdict,351-ac2-ac3-verdict,378-r1-impl}.md`。本表**一个字节没碰**、**不引它们当读数**；两包测试与 d22scan 在含它们在场时跑过且干净 ⇒ 未触发"红归未跟踪件"的隔离条款（第 21 处），无 `git clone --shared` 之需。
- 本程**零产码**：名下只新增本文件；`internal/**`／`cmd/**`／`tools/**`／`docs/PLAN.md`／`docs/specs/**`／`docs/reports/**`／任何票面 **0 字节**；`docs/reports/pending-and-issues.md` 只 grep 只读（分账引 A 号），没写。
- 每格标〔现跑〕（本程本机命令与读数）／〔读码〕（本程逐字读产码）／〔它现跑〕（上游表读数，注明出处）。

---

## 1. 设计核心（这一票真正裁下来的东西，一页版）

### 1.1 当初的撞车（病灶）

两件都已冻结的事互斥：**外来内容进上下文要盖 C25 戳**（票 175）× **超长产物落盘后模型能自己续读**（`PLAN.md:431` D15(3)、票 164 AC#3）。撞点＝D15(3) 桩逐字嵌宿主落盘路径（`…全文见 %s…`，今读〔读码〕`internal/tools/task.go:562`、同款 `internal/agent/spill.go:190`）：一旦 `task.output` 进盖戳名册，这条路径进了污染索引；模型拿它去 `fs.read`，参数只有 `{"path": …}`，而污染检查对路径键不豁免 ⇒ 命中 ⇒ R4 升 L2 ⇒ 测试具 `NoGate{}` 下逐字 `L2 审批通道尚未接入（票 21），已拒绝执行`（〔读码〕`internal/tools/gate.go:160`）。`175-r1` 真撞红两枚续读腿后**不修、回退、立本票**（台账 `A346`）——没有为了变绿放宽任何断言。

### 1.2 owner 两句拍板（形状边界）

`Q-61`＝**甲·精确豁免**（豁免对象不是目录／不是前缀／不是配置项，是"宿主自己亲手写下过的那几条具体路径"）；`Q-63`＝**甲·不算复活**（删连续窗口 ≠ 逐 token 追踪；`taintmatch.go:11-15` 那句 REJECTED 冻结文字一字未动，本程〔现跑〕复读在场）。台账 `A349`／`A353`。

### 1.3 为什么是甲形，另外几支死在哪（全部有读数，不是偏好）

- **乙（桩里不写路径／匹配侧按值放行）＝结构性死**：`Inspect` 只看得到候选值——合法续读与外来说到同一条路径，候选值逐字相同；`177-m1` 用 S3b 把按值名册真造出来，同一枚常量两次 `Inspect` 给出同一个放行 ⇒ **外来那一发同时被放走＝洗戳**（〔它现跑〕`177-mutation-shape-a-redlist-m1.md` §5）。
- **丙（每次续读一张 L2 卡）＝不是技术死，是价目表死**：`177-a0` 独立普查（`177-ac0-shape-c-independent-census.md`，2026-11-05 交件）现量四形状——无审批层＝永不咨询当场拒（`internal/agent/loop.go:919-923`，本程〔读码〕在场）；`NoGate`＝当场拒；界面不可达＝当场拒（fail-closed）；界面在没人答＝**挂起到 C18 截止再拒**（31s 可配下限量到墙钟 34.26s，生产默认 300s）。R4 的"不受会话授权覆盖"是**结构性**的：会话授权查表整段包在 `sil.Level == risk.L1` 里（〔读码〕`internal/tools/bridge.go:387/:395`），L2 连查都不查 ⇒ **N 次续读＝N 张卡、N 段等待**，无人值守等于读不回来。丙未被选中；a0 表同时给出"若哪天要走丙，票 164 两枚判据要怎么改形"（＝AC 变更，须 owner）。
- **"宽一档"的偷懒形＝全部有牙**：整段 mark 不入索引（W-1）、按值名册越出本 mark（W-2）、目录前缀（W-3）、**按位置推给全树**（W-4——`177-m1` 量到该形在旧树全仓零红，`177-r1` 造出常驻判据后当场响）、空判定排错序（W-5）、跨 scope 存活（F）。`183` 补 M3（正文不含声明路径时仍须命中——fail-closed 不许只靠注释），`185` 补同 basename／一字符之差／跨 scope 借名／成本面钉住。

### 1.4 落地的最小形状（三条不变式 ＋ 两条被冻结文字逼出来的顺序）

1. **mark 局部窗口排除，零新名册**：不是"维护一张路径列表"，而是"在宿主自产那一枚 mark 的片段索引里，把声明的窗口抠掉"，随 mark／`Scope.Close` 生死；`Provenance`/`scopeReg` 零加收字段。四笔：`069344b6`（`taintmatch.go` 变长 `skip`＋`windowExcluded`，现读 `:146/:159/:187`）→ `149c6d6d`（`MarkWithHostPath`，现读 `provenance.go:484/:487/:563`）→ `3d4c5795`（桥侧每调用 `hostPathBox`，现读 `bridge.go:758-780`）→ `b3c389b5`（`task.go` 指针分支 `box.set(rec.ArtifactPath)`，现读 `:547-563`）。C1 `Result` 零加收、`Loop` 零加收带 taskID 的导出方法、`risk.PathResolver` 之外零 `filepath.Clean|Abs`。
2. **按值规则 ＋ 派生投影**（183/185 两批）：`attachDeclaredPath`（`provenance.go:627`）把"声明路径自身的窗口"按值扣掉——修掉"位置无关的哈希集＋`Contains` 让豁免失效"的真机根因（`183` 判根因＝(d)，变异 M-D1 决定性）；`hostMintedPathFor(scopeID, origin)`（`provenance.go:606`）＝**派生不存储**：由本 scope 自己 marks 的 `declaredPath` 投影出"宿主亲手落盘的那一条"，参数侧永不查表，随 marks 生灭（`ee1e2118`，A381 第四支）。
3. **两条硬顺序不是风格，是冻结文字逼的**：先归一化再定位跨度（原始下标会偏，order 腿钉住：偏移差 9＞窗口 8 时边缘漏进索引）；空判定排在排除之前（`provenance.go:481` "Returns false only when…" 那句契约在 S5 变异下变假，W-5 钉住）。
4. **成对判据（"缺一形＝装饰"的兑现）**：正向＝W-1／J-3／183 referee 三枚"宿主自读干净"；反向＝W-2／J-2／185 外来携带同路径仍命中＋兄弟／近邻仍命中（W-3／183 AC#4／185 同名异目录）。

### 1.5 化解链（四批提交，锚与时刻全查证〔现跑〕）

`177-r1` 甲形落地＋八枚常驻判据（09-27 23:2x，四笔＋判据件 `a19fa2e7`）→ `175-r2` 名册九枚＋桥级五枚判据（`962ba9d1` 08:52:58／`58db7400` 08:53:12，盖戳真生效）→ `183` 真机接缝按值规则（`0662a35a`；端到端第一发翻绿）→ `176-r1` 起跑口（`41e8873bd`／`06dae8316` 09:23，`TaskBackfill.Backfill` → `Roster.Record`，`internal/tools/task_backfill.go:141`）→ `185` 派生投影修"第二次续读"（`ee1e2118` 16:43:15；常驻件 `internal/risk/pointer_185_test.go` ＋链判据 `internal/tools/pointer_185_reread_chain_test.go`）。`315`/`179` 是同一根管子更后面的两扇门（子代理 scope／声明 L0），不是排他的化解者。

---

## 2. AC#5 分账（结题动作：四枚是同一根管子上的四处，批一枚不通其余）

> 本表就是票 177 AC#5 的载体。⚠ 逐行都答"它今天是什么状态"，**不用"续读终于通了"顶掉任何一枚**。台账 `A1279` ② 已做过 175 那份三态分账并注明"**177 AC#5 是它的更宽版（多 `Q-60` 一栏）、§3 可喂**"——本节按现行树把宽版补齐并全部现读重锚。

| # | 邻居 | 今天状态 | 现量凭据（本程〔现跑〕/〔读码〕，行号按 HEAD `8879c3d7b1`） | 还欠什么 |
|---|---|---|---|---|
| ① | **票 175（盖戳）** | **已完**（票面 5/5 勾、名 `-done`） | C25 名册九枚含 `SrcTaskOutput`（〔读码〕`provenance.go:92-102`）；桥级正向 J-1 `TestTaskOutputStampedOnRealBridge175r2` **PASS 0.00s**〔现跑〕；其 AC#5 三态分账已由 `close3-r1` 做成（〔台账 A1279 ②〕） | 无 |
| ② | **票 174（授权根）** | **两处半**（`Q-60` 丙射程内的都落了，能力一格没动） | 授权根仍只来自 `[fs] allowed_dirs`（〔读码〕`run.go:608-609`）；`artifacts` 目录未入根（〔读码〕`run.go:1625`／`:1779`）⇒ **默认仍读不到**；判定者接生产已落（〔读码〕`run.go:653` `TaskDeps{Roster: rt.tasks, Paths: rt.paths}`）；回执诚实化已落（〔读码〕`spill.go:24/:43/:192-200` `pointerNotice`）。票面 7/9 勾 | **AC#4（契约轴：`PLAN.md:431` 那句文字改不改＝须 owner 单独批）＋ AC#5（分账格，本表 §2 ② 行即可喂）** |
| ③ | **`Q-59`（偏移参数）** | **DEFERRED（owner 已批乙）**，报警器在岗且安静 | 五字段 DEFERRED 行在 `PLAN.md:1538`（A1279 已纠正旧记 `:1532` 的漂移；本程〔现跑〕复读该行在场）；`fs.read` 仍只 `{path, max_bytes}`（〔读码〕`fs.go:112-115`）、`max_bytes` 只能改小、超 256 KiB 静默截断 ⇒ **后半截仍读不回**；canary `TestPointerPast256KiBIsNotFullyReadable` **PASS 0.13s**〔现跑〕＝"没人偷偷加偏移"，属正常 | 哪天真加偏移须 owner 另批（`Q-59` 甲）；届时 canary 当场转红是设计好的闹钟 |
| ④ | **`Q-60`（授权根拍板）** | **owner 已批丙**（能力一格不动，回执说实话） | 默认空 allowlist＝每次读根外路径都要一张被批的 L2 卡（`SPEC-03:35`）；L2 不查会话授权（〔读码〕`bridge.go:395`）、不聚合 ⇒ 丙形价目表由 `177-a0` §4 钉死。177 的豁免**没有**把 artifacts 放进授权根，也没有给任何目录开洞（AC#4 零放宽，四批交付全过 d22scan） | 若哪天要放开那棵含 `secrets\` 的树＝安全边界变更，**须 owner 说话**（甲／乙两支都还在桌上） |
| ⑤ | **票 183＋185（第五席，票面四枚之外、诚实具名）** | **已落、常驻判据全绿** | 183 修"第一发能读回"、185 修"第二次续读不再被上一次的读回正文挡住"（本程〔现跑〕：`TestPointer185SecondRereadOfAHostArtifactIsClean` PASS、反控 `…NothingDeclaredIt` PASS、链判据 `TestPointer185LongRereadChainNeverRemintsTheBlocker` PASS 0.11s、`TestPointer185ChainDoesNotLaunderAForeignCarrier` PASS 0.13s）；证据表 4 枚在盘（`185-{paged-reread-r1,r2-ac1-ac7-readings,r3-head-criterion,reread-owner-census-c1}.md`） | 证据表 §②–⑪ 的填齐程度归 185 自己的收口（a1 普查具名过；非本票义务） |

**一句话**：四枚（加第五席）今天各自"通到哪一格"清清楚楚——**盖戳真生效、豁免承重、授权根没放宽、偏移没加（按乙挂账带闹钟）、续读链两发都过（桥级/包级常驻判据为凭）**。任何一枚的"通"都不证明其余枚的"通"。

---

## 3. AC#6 排程约束（继承账，非本票可翻——只把读数摆齐）

- 票面原文：本票落地那天，**票 176 不得早于本票合入**。
- 时序账〔现跑〕：177 落地批最后一枚产码提交 `b3c389b5` ＝ **09-27 23:26:19**；176-r1 起跑口两枚生产提交 `41e8873bd`／`06dae8316` ＝ **09-28 09:23** ⇒ **顺序成立，裕量约 10 小时**，"不盖戳就能读回外部内容"那一段从未先通。
- 附记：A1279 ⑥ 判"AC#6 已过期、归属请编排者裁"——过期 ≠ 被违反；本格翻勾与否由编排者裁，本表只供时序读数。

---

## 4. 本程现树读数（判据全绿；〔现跑〕2026-11-06，命令可复跑）

- 环境：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`（永不 CGO_ENABLED=0）。
- **tools 级七枚**（164 两枚续读腿＋175 三枚 J 腿＋183 两枚）：`go test -count=1 -v -run 'TestLongOutputPointerRecoversEveryByte|TestPointerPast256KiBIsNotFullyReadable|TestTaskOutputStampedOnRealBridge175r2|TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2|TestHostMintedPointerRereadStaysCleanOnRealBridge175r2|TestPointer183BackfilledArtifactRereadStaysClean|TestPointer183SiblingArtifactInTheSameDirectoryStillHits' ./internal/tools/` ⇒ **7 PASS／0 FAIL**。
- **risk 级 26 枚**（ShapeA 八枚＋183 族八枚＋185 族十枚）：`go test -count=1 -v -run 'TestShapeA|TestPointer183|TestPointer185' ./internal/risk/` ⇒ **26 PASS／0 FAIL**。
- **185 链判据两枚**：`go test -count=1 -v -run 'TestPointer185LongRereadChainNeverRemintsTheBlocker|TestPointer185ChainDoesNotLaunderAForeignCarrier' ./internal/tools/` ⇒ **2 PASS**。
- **两包整跑**：`go test -count=1 ./internal/risk/ ./internal/tools/` ⇒ `ok internal/risk 13.946s`／`ok internal/tools 26.853s`（含别腿未跟踪件在场，仍绿）。
- **d22scan**：`sh scripts/d22scan.sh` ⇒ **rc=0**，逐字 `clean - no D22 ban violations`（live scopes：bans #1-5 internal/=363、cmd/=58、ban #7 internal/tools/=48、ban #8 internal/=907、cmd/=216）。
- **冻结文字〔现跑复读在场〕**：`taintmatch.go:11-15`（逐字 `Exact per-token taint tracking is REJECTED (16.9#1) and must not be resurrected here.`）；`provenance.go:481`（`Returns false only when …`）。本程零字节、自然零触碰。

## 5. 本表**不是**什么（别把上面的绿读过头）

1. **不是"真机 CLI 端到端今天复绿"的读数**：本程全部读数是桥级/包级常驻判据。真机第一发读回的最后读数停在 09-28（`183` e2e）；post-185 的 CLI e2e＝179 AC#8 那一发，未在本程序跑（`cmd/wisp` 整包在本机有既知环境红与非环境红，见 `177-ac0-shape-c-independent-census.md` §6.1，本程不拿它当门）。
2. **不是 AC#0 的翻勾凭据**：丙形普查已真交件（`177-ac0-shape-c-independent-census.md`），"零读数"状态已终结，但**翻勾是编排者的裁定**，本表只引用。
3. **不是任何票面的写面**：本程没动 `177-*.md` 票面一个字节（含没翻任何框）。
4. **不是产码评审**：四批落地各自带非实现者验收（175-r2 五判据、183-v2 对抗验收、176-r1 四判据等）；本表只在**现行 HEAD** 上复核判据仍绿＋重锚行号。

## 6. 行号重锚账（本表全部现读；引用前别照抄旧表行号）

| 旧表说 | 现行 HEAD 实测〔读码/现跑〕 |
|---|---|
| 桩 `…全文见 %s…` 在 `task.go:245`（m1/r1 时代） | **`task.go:562`**（174 桩 notice 占位与 188 家族改动后漂移）；同款 `spill.go:190` |
| 指针分支/`box.set` 在 `task.go:254-267` | **`task.go:547-563`**（`:558-560` `box.set`） |
| `case RiskL1, RiskL2` 在 `loop.go:916-924` | **`loop.go:919`** 起，语义未变 |
| `sensitiveSourceTools` 八枚 | **九枚**（`provenance.go:99-102`，`SrcTaskOutput` 在册） |
| 授权根 `run.go:327-331` | **`run.go:608-609`**；`Paths: rt.paths` 接线在 **`run.go:653`** |
| `fs.read` 参数 `fs.go:113-115` | **`fs.go:112-115`**（结构体上移一行，语义未变） |

## 7. 伪授权两栏（各带出处）

| 栏 | 内容 | 出处 | 本程态度 |
|---|---|---|---|
| 甲：转述在场、我未逐字复核台账原文 | "`Q-61`＝甲、`Q-63`＝甲已 owner 批" | 派单＋票面 09-27 21:5x 追加；台账 grep 现量 `A349`（8068 行）／`A353`（8116 行"四枚 Q 全批"）逐字含"批乙/丙/甲"清单 | 派单明写"`Q-61` already approved"且台账条目在场；我只把批准当**形状边界**（目录/前缀/配置级零造、冻结文字零改），未据此动任何产码 |
| 乙：在场而互相矛盾 | 派单称"设计核心＋分账"待做；票面 AC#0 那格写"丙形零读数、不许被当成已裁" | 派单 vs 票面 `:34` | 按**更新事实**处理：AC#0 普查 2026-11-05 已交件（零读数状态终结），但翻勾归编排者；本表 §5.2 具名，不冒充"丙已裁" |

## 8. 交件账

- 被拒／失败调用：零枚被权限拒；过程失败 1 次——`git cat-file -t` 一次传两枚提交号报用法错（重跑为逐枚循环后成功，无副作用）；`go test -run` 两次模式不匹配（`[no tests to run]`，读文件内真实用例名后改精确 `-run` 成功），均为取数修正、零树影响。
- 删除命令：**零枚**（没跑 `rm`／`git clean`／`restore`／`checkout .`）；没建 worktree；**push 零枚**；`git add -A`／`amend`／`reset`／`rebase`／`stash` 零枚；提交带显式 pathspec（只有本文件）。
- 凭据值抄录：零枚。

**next＝** ① 编排者核本表后裁 **177 AC#5 翻勾**（本表 §2 即其载体）与 **AC#6**（时序读数在 §3，过期≠违反）；② **177 AC#0** 翻勾凭据＝a0 普查表（2026-11-05 交件），同一轮顺手裁；③ 仍未闭合、与 177 无依赖的邻居欠账＝**174 AC#4**（契约轴，须 owner）、**179 AC#8** 的 post-185 CLI e2e（给 177 AC#3 的端到端条件补最后读数）、185 证据表收口；④ 台账 A 号登记与推送归编排者。
