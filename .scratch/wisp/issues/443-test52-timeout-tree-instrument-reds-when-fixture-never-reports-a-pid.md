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

- [ ] **AC#1 复现（本机可做）**：command_52_test.go 无 go:build 限制 ⇒ 在本机 Windows 上**必须**复现「fixture never reported a grandchild pid」；给 N 发与命中率。⛔ 不许用「跑一百次没红」当复现。
- [ ] **AC#2 区分弃权与真红**：夹具没报出 pid ⇒ **具名弃权行＋数字**（等了多少 ms／重试几次／pid 是否出现过），⛔ 不许硬红也⛔ 静默绿；**而「超时没杀掉树」必须仍然真红**。
- [ ] **AC#3 牙还在**：仓外导出树里让 TimeoutKillsTheProcessTree 的杀树逻辑失效（如去掉 taskkill/job-object 那一手）⇒ 修后仪器**真红**并逐字点名；还原后 SHA256 同串、再跑绿。
- [ ] **AC#4 普查同类**：全仓 grep 「never reported a … pid」「fixture never reported」以及**所有**「夹具没报数就判红」的自检，给点名表并逐条判是否同形（⛔ 只修这一枚了事）。
- [ ] **AC#5 门禁四数**：go vet ./internal/plugin/manifestslot rc=0；gofumpt -l 空；d22scan clean；path-length GREEN；**产码 0 字节**（除非能证明是产码缺陷 —— 若真是产码问题，先报我再动）。

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
