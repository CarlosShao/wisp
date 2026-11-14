# 票 443 — 票 52 那枚「超时杀进程树」仪器，在夹具没报出孙进程 pid 时判红（票 441 同类）

**状态**：待派
**出处**：编排者 2026-11-14 15:4x 读 CI 日志现取（run 43567331249 head 1cc304e0；run 43568581591 head 1f0d04714），不是猜。

## 立案事实（逐字，两个不同 head 同一条）

`
--- FAIL: Test52TimeoutKillsTheProcessTree (0.76s)
    command_52_test.go:748: the fixture never reported a grandchild pid (C:\Users\runneradmin\...)
FAIL	github.com/CarlosShao/wisp/internal/plugin/manifestslot	13.385s
`

- **1cc304e0 红、1f0d04714 红、紧接着的 46352fd3 绿** ⇒ **3 发里红 2 发**，且其中一发（1f0d04714）**只改了 probes 和证据件**，与产码零关系。
- 红句出自**夹具自报**：「the fixture never reported a grandchild pid」＝仪器**没观测到**它要杀的那棵树的孙进程，然后**判红**。这正是票 441 的同一机理：**弃权被读成缺陷**。

## 为什么值得修（比 441 更急）
441 那枚是 Linux core 才跑；**这枚在 	est-windows step 13 每发必跑**，而它 3 发红 2 ⇒ 它是当前 CI 最响的一个假红源。且 	est-windows 是**唯一**跑 ball/proc/risk 的 job ⇒ 它一红，整片 Windows 面都失去绿证。

## 完成判据（AC）

- [x] **AC#1 复现（本机可做）**：command_52_test.go 无 go:build 限制 ⇒ 在本机 Windows 上**必须**复现「fixture never reported a grandchild pid」；给 N 发与命中率。⛔ 不许用「跑一百次没红」当复现。
- [x] **AC#2 区分弃权与真红**：夹具没报出 pid ⇒ **具名弃权行＋数字**（等了多少 ms／重试几次／pid 是否出现过），⛔ 不许硬红也⛔ 静默绿；**而「超时没杀掉树」必须仍然真红**。
- [x] **AC#3 牙还在**：仓外导出树里让 TimeoutKillsTheProcessTree 的杀树逻辑失效（如去掉 taskkill/job-object 那一手）⇒ 修后仪器**真红**并逐字点名；还原后 SHA256 同串、再跑绿。
- [x] **AC#4 普查同类**：全仓 grep 「never reported a … pid」「fixture never reported」以及**所有**「夹具没报数就判红」的自检，给点名表并逐条判是否同形（⛔ 只修这一枚了事）。
- [x] **AC#5 门禁四数**：go vet ./internal/plugin/manifestslot rc=0；gofumpt -l 空；d22scan clean；path-length GREEN；**产码 0 字节**（除非能证明是产码缺陷 —— 若真是产码问题，先报我再动）。

## 边界
写面：internal/plugin/manifestslot/** 的 *_test.go（＋只在 AC#5 允许且先报的情况下动产码）、docs/evidence/s1/443-*、.scratch/wisp/probes/443/**、本票 Progress log。
⛔ 放宽既有断言；⛔ 加 	.Skip（untests.sh 把 core 档顶层 SKIP 判致命 —— 用「具名行＋数字」代替，照票 441 的做法）；⛔ dd -A／--amend／eset／ebase／stash／checkout .；⛔push。

## 参照
票 441 的落法（internal/winsec/posix_seal_window_120_other_test.go，AC#1 ABSTAIN(441:<shape>:raced-quiet) attempts=... samples=...）：**把两件事拆成两个载体** —— 「窗口到达过」与「密封还能封住」分开证，弃权可见且有数字，不再被判成缺陷。

## 撤销口令
「票 443 撤」＝删本票＋回退其写腿 commit（新 commit 反改，不 reset/amend）。

---

## Progress log

- **2026-11-14 15:4x 编排者**：本票立。依据＝CI run `43567331249`(head `e1cc304e0`) 与 `43568581591`(head `1f0d04714`) 均红，`a46352fd3` 绿 ⇒ **3 发红 2**；红句 `command_52_test.go:748: the fixture never reported a grandchild pid`。
- **2026-11-14 15:5x 编排者**：机主要重启 DSH，**写腿被中途叫停**（不是我判它错）。⇒ 工作树里留下 **2 枚未提交的测试件改动**：
  `internal/plugin/manifestslot/command_52_main_test.go`、`internal/plugin/manifestslot/command_52_test.go`
  **接手腿须知**：这两枚是**被中断的半成品**，**既不可采信也不可当成已修**。开工第一件事＝读一遍 `git diff -- internal/plugin/manifestslot/` 判断：能用就接着做并说明，不能用就**显式 `git checkout --` 还原后在 Progress log 记一行**（⛔ 不许默默覆盖，⛔ 也不许默认它们是对的）。
  **产码 0 字节**（`internal/plugin/manifestslot` 的非测试文件零改动），`cmd/wisp` 干净（0 行）。
- **2026-11-14 16:1x–17:0x 写腿（接手）**：**上一条腿的 2 枚半成品判读＝可用，接着做并在其上补齐**（`git diff -- internal/plugin/manifestslot/` 逐行读过：两载体骨架、`awaitPID443`/`judgeTreeDead443`/控制运行、`fixtureDelay443` 旋钮都在，方向与本票 AC#2/#3 一致，未发现要放宽断言或改产码的地方）。本腿在其上做的：①`clearPIDFile443` → `clearPIDFile`（138 也要用）；②给 138 补 `clearPIDFile`（同句同夹具目录，陈旧 pid ⇒ 假绿/假红；它的硬红保留——120s 预算/10s 等待/无杀树在飞）；③把旋钮的复现值钉成 1500 并写明 3000 会踩控制预算（实测第 1 发 `RED(443:timeout-tree:control-no-pid)`）；④证据件 `docs/evidence/s1/443-instrument-abstention.md`。
  **读数**：复现＝自然 N=50（25+25）**3 发弃权 = 6%**，旋钮 1500ms N=5 **5/5**、3000ms N=20 **20/20**；弃权行 `AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) budget-ms/elapsed-ms/waited-ms/polls + pid-seen=no + last-poll-error`；真红 token `AC#2 RED(443:timeout-tree:tree-survived)`。
  **牙**：仓外导出树把 job-object 树杀降级成单进程杀（`7EFD1619…`）⇒ `RED(…tree-survived) control run: grandchild pid 58852 is still running`，FAIL；还原回 pristine `B597F102…`（**同串**）⇒ PASS。★「把 kill 整枚摘掉」那条错突变也记在证据 §3.1：它会让 `proc.Wait()` 挂到 10 分钟超时 panic，不是点名红。
  **门禁**：vet rc=0、gofumpt 本包空、d22scan clean（正控 PASS=45）、产码 0 字节；**path-length 仍 RED 且不在本票写面**——缺 roster 的是 `.scratch/wisp/issues/274-…-done.md`（文件名 102 字符 > 100 帽，HEAD 已提交、非本腿产物）；`gofumpt` 全量另有 1 枚非本腿 `internal/models/acl_sid_121_test.go`（CRLF）。两条都已在证据 §5 点名交回。
  **同族普查**：`git grep -E 'never reported|never appeared|did not appear|…' -- '*.go' ':(exclude).scratch'` ＝ 62 行（`*_test.go` 50 行），属「夹具没报数就判红且缺席可能是调度」的只有本片与 138 两枚；`resident_hotkey_live_258` 形似（live 档、12s 有界轮询）已在证据 §4 交回下一票。
  **落地**：commit `d78a0ae5c`（显式 pathspec；⛔ 未 push）。落地后工作树对本腿路径干净，**在这串字节上复跑**：无旋钮 PASS、`WISP52_443_FIXTURE_DELAY_MS=1500` PASS（先弃权行再由控制接住）、`Test138` PASS（`probes/443/14-post-commit-rerun.txt`）。
  ★ 本票 `probes/443/logs/prefix-delay700-clean-01..05.log` 是**上一条被中断的腿**留的旧码直接复现（5/5 FAIL、逐字 `the fixture never reported a grandchild pid`，行号 765 是它当时半成品的行号）；本腿逐条打开核对后才在证据 §1.4 引用，`probes/443/r1/`、`repro/` 同样是它留下的台面，原样保留未改。
---

## 编排者收件（2026-11-14 16:5x，非实现者，只追加）：5/5 翻，本票结

**三把尺（我自己量）**
- **尺1 提交真在**：`d78a0ae5c`（主体）＋`310b38cf4`（落地复跑记录）两枚 `cat-file` rc=0。
- **尺2 件真在**：`docs/evidence/s1/443-instrument-abstention.md` **19537 B**；弃权 token 现读在 `internal/plugin/manifestslot/command_52_test.go:811`（`t.Logf("AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) %s pid-seen=no …`），`:738` 是它的注释说明。
- **尺3 我自己复跑**：`go test ./internal/plugin/manifestslot -run 'Test52|Test138' -count=1 -v` ⇒ **RUN 21 / PASS 21 / FAIL 0 / SKIP 0**，`ok … 14.280s`，RC=0（日志 `D:\tmp\wisp\20261114\lead\443-rerun.log`）。关键那枚 `Test52TimeoutKillsTheProcessTree` **PASS 1.30s** 且走的是**真路**：`command_52_test.go:805: AC#2 raced run: grandchild pid 9812 is gone, so the kill reached the tree`；同跑 `Test138TaskCancelKillsTheProcessTree` PASS 1.57s。
- **产码 0 字节我自己核**：`git show --name-only d78a0ae5c` 里**非 `_test.go`、非 `probes/`/`docs/` 的文件 = 0 枚**。

**逐格**
- **AC#1** ✓ 自然 N=50 ⇒ 弃权 3 发（6%）；放大到 1500ms N=5 ⇒ **5/5**、3000ms N=20 ⇒ **20/20**。另有上一腿留下的 5 发 `prefix-delay700-clean-*`：旧断言 **5/5 FAIL** 且逐字同 CI（`the fixture never reported a grandchild pid`），腿逐条打开核过才引用。
  ⚠ **我这发没复现出弃权**（走的是真路）⇒ 弃权侧读数是腿的，不是我的。
- **AC#2** ✓ 弃权＝具名行＋四数（budget/elapsed/waited/polls ＋ `pid-seen=no` ＋ `last-poll-error`），**不红也不静默**；真红＝`AC#2 RED(443:timeout-tree:tree-survived)` 点名活 pid。**两条我都现读过 token。**
- **AC#3** ✓ 仓外导出树把 job-object 树杀降级成单进程杀 ⇒ `RED(443:timeout-tree:tree-survived) control run: grandchild pid 58852 is still running`，FAIL；还原后 pristine `B597F102…` **同串** ⇒ PASS。★它还记了一条**错突变**：把 kill 整枚摘掉**不红**，会让 `proc.Wait()` 挂到 10 分钟 panic（直子进程只有 kill 才退）——这条比我要求的更有用。
- **AC#4** ✓ 普查 62 行（测试件 50 行），判「缺席可能是调度」的只有本片与 **138** 两枚（都已各自处理）；`resident_hotkey_live_258` 形似 ⇒ **已具名转到票 258**（见下）。
- **AC#5** ✓ `go vet` rc=0、`gofumpt -l` 本包 0 文件、`d22scan` clean（正控 PASS=45）、**产码 0 字节**、三条定向测试串行全绿。
  **path-length 那一项**：腿跑时是 **RED**，但红源是**我自己**结 274 时文件名顶到 102 字符（roster 外第 70 个）⇒ 我把 274 改名到 69 字符后**现跑 `VERDICT GREEN`（over 69 / roster 69 / not-in-roster 0）**。规矩记在台账 **A1866**。

**我自己的失误，记在此不隐**：推台账那两笔时，**443 的两枚 commit 跟着一起上了 origin**（我核之前就推了）。核完是绿的、无实际损失，但动作违反「核过之后再推」。改法：`git push` 前先 `git log --oneline origin/dev..HEAD` 看清单里有没有别人的腿。

**并行带来的额外收获**：这条腿还如实交回两个**不在它写面**的红：`internal/models/acl_sid_121_test.go` 的 gofumpt（CRLF→LF，历史件）。已记，不吸收。
