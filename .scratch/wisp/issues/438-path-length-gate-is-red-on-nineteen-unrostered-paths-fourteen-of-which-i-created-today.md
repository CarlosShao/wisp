# 438 — 票 262 那道路径长度尺**今天真红**：19 枚超预算且未登记，其中 **14 枚是我本轮改名/立票造成的**（名字超过 README 规则 9 的 100 字符帽）

- Status: **ready-for-agent**（写面＝`.scratch/wisp/issues/` 的**文件名**（`git mv`）＋ `scripts/check-path-length-budget.sh` 的 **roster 数据区**；⚠ 先只读取证，判清"改名"还是"登记"，再定改法）
- 来源：**我 2026-11-13 10:1x 现跑** `sh scripts/check-path-length-budget.sh --with-self-test` ⇒ **rc=1 红**（正控三步全 PASS，说明尺本身没坏）
- 关联：票 262（本尺的立案票，已落地）· `.scratch/wisp/issues/README.md:67` **规则 9**（新票名 ≤100 字符）· 台账 `A1813`

## 这是什么（人话）

仓里有一道尺盯**文件路径长度**。理由是：本机 CI runner 上，路径太长会让 `actions/checkout` **整步失败**——不是某个测试红，是**整台机器捡不起代码**。这道尺已经装好了（`ci.yml:176` 的 `Tracked path-length budget (ticket 262)` 步，跑 `scripts/check-path-length-budget.sh --with-self-test`）。

**今天它红了。** 现量逐字：
```
check-path-length-budget.sh: denominator: tracked paths=19508  over-budget=78  covered by roster=59  not in roster=19
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/174-…
…（19 条）
```

**★这 19 枚里有 14 枚是我造成的。** 我今天（2026-11-13）把一批票改名加 `-done` 后缀，而它们的名字本来就长——加 5 个字符就推过了帽。另外几枚是我今天新建的票，**名字直接违反了 README 规则 9 的 100 字符上限**。

**⇒ 这是我的违规，不是产品缺陷。** 但**尺红了就是红了**，必须修好，⛔ 登记了事。

## 现量（逐字，可复跑）

**跑法**（⚠ 必须用 `D:\work\soft\Git\bin\bash.exe`）：
```
D:\work\soft\Git\bin\bash.exe -c "cd '/d/work/workspace/projects plans/Wisp' && sh scripts/check-path-length-budget.sh --with-self-test"
```
**读数**：
```
positive control PASSED          ← 三步正控全过，尺本身没坏
unit=bytes
hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165
check: relative_length + 44 > 165, or a last name component > 100
denominator: tracked paths=19508  over-budget=78  covered by roster=59  not in roster=19
longest=185 chars relative
worst full path on the self-hosted runner=229 chars (hat budget 165, wall open interval (206,217])
bands: over the hat=78  of which in the 122..180 middle=74  past the old debt line(180)=4  in the wall interval=0
```

**19 枚未登记（我现跑逐字，截到 110 字符）**：
```
.scratch/wisp/issues/174-task-output-points-at-a-file-the-model-cannot-read-by-default-spill-artifacts-are-out
.scratch/wisp/issues/194-the-panel-has-two-method-rosters-that-do-not-know-each-other-owner-ruled-align-the-co
.scratch/wisp/issues/246-the-process-that-runs-tasks-is-not-the-process-holding-the-ball-so-no-veto-channel-ha
.scratch/wisp/issues/269-gofmt-step-is-held-red-by-its-own-committed-mutation-sample-and-eats-three-unprotecte
.scratch/wisp/issues/425-resident-shutdown-record-is-never-last-because-ticket-43-replays-ten-steps-in-the-def
.scratch/wisp/issues/427-the-12r2-ast-check-still-looks-for-marksettling-inside-openscope-but-it-moved-into-sc   (×2)
.scratch/wisp/issues/428-the-26p-no-aec-echo-test-sees-zero-episodes-so-the-judge-never-received-the-announcem   (×2)
.scratch/wisp/issues/429-the-384-watchdog-never-sees-the-tree-back-under-the-cap-because-gc-does-not-return-pa   (×2)
.scratch/wisp/issues/430-the-365-stuck-task-probe-never-goes-stuck-because-the-tiny-context-window-never-becom   (×2)
.scratch/wisp/issues/432-the-three-258-hotkey-tests-red-on-first-real-execution-because-the-old-binding-and-li   (×2)
.scratch/wisp/issues/433-the-246-dev-leg-test-is-stably-red-on-its-own-tempdir-cleanup-failing-because-a-handl   (×2)
.scratch/wisp/issues/435-the-33r11-locked-suffix-roster-missed-four-new-names-and-the-question-is-whether-they   (×2)
```

⚠ **为什么有些出现两次**：腿必须自己查清（我猜是「相对路径」与「最后一段名」两条判据各命中一次，但**⛔ 照抄我的猜**——去读脚本 `:121-135` 那段的 roster 语义）。

## 机制链（逐字，⛔ 凭印象）

- **尺**：`scripts/check-path-length-budget.sh`
  - `:120-127` 「THE ROSTER AND ITS GUARDS (ticket 262 AC#4)」——**豁免是文件里的数据，⛔ 注释**
  - **两条守卫**：**A**＝超预算且不在 roster ⇒ 红；**B**＝roster 里有但今天不再超预算 ⇒ 也红（**漂移对**：A 抓"树长过名册"，B 抓"改名把名册落下"）
- **CI 步**：`.github/workflows/ci.yml:176` `Tracked path-length budget (ticket 262)`，`:207` 跑 `sh scripts/check-path-length-budget.sh --with-self-test`
- **规矩**：`.scratch/wisp/issues/README.md:67` 逐字 `新立票一律 NN-短 slug.md ≤100 字符；因果句子写在票面第一行标题里，不要塞进文件名。`

## ⚠ 必须自己判清的核心岔路（⛔ 不许我替你定）

**"改名 vs 登记"——两种答案，修法相反：**

- **甲：改名（把名字缩短到 ≤100）** ⇒ 这是 **README 规则 9 要求的正道**。⚠ **但 `git mv` 会破坏历史引用**：台账与证据件里按名引用这些票的地方会失联。
  **判据**：先数一下有多少处按名引用（`git grep -l '425-resident-shutdown'` 之类）。**若引用很少 ⇒ 改名；很多 ⇒ 要权衡。**
- **乙：登记进 roster** ⇒ **⛔ 消红工具**。roster 的语义是"**这些是历史债、我认了**"（`:120-127`），**⛔ "我新造的长名字也认了"**。**若用乙处理我今天新造的票，等于把我的违规写进债表。**
- **★很可能混合**：**旧债（174/194/246/269）** 也许该登记（它们是历史遗留），**我今天造的（425–435）** 该改名。

**⛔ 拍脑袋。逐枚看它是"历史债"还是"我今天造的"。**

## ★必须先做的判断（⛔ 跳过）

1. **`git log --diff-filter=R --oneline -- '.scratch/wisp/issues/'`**：查这些名字**什么时候变成这么长的**。
2. 对**每一枚**，判它属于：
   - **(a) 我今天改名加 `-done` 造成的**（425/427/428/429/430/432/433/435 等）
   - **(b) 我今天新建票名就超标**（434/436/437/438）
   - **(c) 历史债**（174/194/246/269）
3. **数一下按名引用的处数**（决定改名代价）。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. 跑上面那条命令，贴**完整 19 条**（⛔ 我截断的版本）。
  2. **查清"为什么有些出现两次"**（读脚本 `:121-135`，逐字引用判据）。
  3. `git log --diff-filter=R` 查这些名字的来历。
  4. **逐枚分类**（a）／（b）／（c），给 `文件:行` 或 commit。
  5. 数**按名引用的处数**（`git grep -c` 之类），贴读数。
  6. 给出**你的修法方案**（改名／登记／混合），并说明为什么。

- [ ] **AC#1 按 AC#0 的方案修**：
  - **改名** ⇒ `git mv`（⚠ **显式单枚**，⛔ 通配符批量），并**同步修掉按名引用的地方**（若有）。
  - **登记** ⇒ 在 `scripts/check-path-length-budget.sh` 的 roster 数据区加行，**每行写清理由**（⛔ 空理由）。
  - **必须配突变控制真跑**：改名/登记后跑尺 ⇒ **rc=0**；**故意把一枚改回长名 ⇒ 必须 rc=1 并点名它**。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。
  ⚠ **本票会大量改文件名**——`git diff --name-only` 会长，**必须逐枚列出旧名→新名**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l internal cmd tools scripts frontend` 零输出 · **`scripts/d22scan.sh`**。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 防复发写清**：**为什么我（编排者）今天会造出 14 个超长名？** 是 README 规则 9 没人执行，还是 `-done` 后缀没被规则 9 覆盖？**给出可执行的防复发措施**（例如：规则 9 改成"**含 `-done` 后仍 ≤100**"，或加一道立票时的自检）。
  ⛔ 改 `README.md`（那是票池规矩，属契约文字面 ⇒ **写建议，⛔ 自己改**）。

## 本票**不**解决

- 不改 `README.md` 规则 9 的文字（属契约面，只写建议）。
- 不动那 59 枚**已登记**的旧债（它们是历史，⛔ 我今天的责任）。

## 派单约束

- 先只读（AC#0 零写），判定后再动。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- **⚠ `git mv` 必须单枚显式**，⛔ 通配符/循环批量（一旦错一个就不可逆）。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\438-<腿码>\`。
- **⚠ 本票与真窗族测试（`Test258*`／`TestAC4Focus*`）⛔ 并行**。
