# 45-v1 验收裁决｜票 45 诊断守卫（45-r1 交付件，非实现者复核）

裁决腿：45-v1（只读验收；裁决者≠实现者）。起手 HEAD `247eb436e`，含被裁 commit `39cf67751`（`git merge-base --is-ancestor` exit 0）。
被裁交付：`internal/observe/guards_45.go`（实测 614 行＝题面符）、`guards_45_test.go`（772 行）、`diagnostics.go`（629 行，本次 +405＝stat 符）。
派单存档：`.scratch/wisp/probes/orch/20261103/txt/45v1-dispatch.txt`。

## 起手三核（读数）

| 核 | 尺 | 读数 |
|---|---|---|
| 包全量真跑 | `go test ./internal/observe -count=1 -timeout 12m` | **`ok 5.150s`，FAIL=0**（收尾另跑一遍 `ok 4.319s`，两遍一致） |
| observe 无他腿脏 | `git status --porcelain -- internal/observe` | 起手**空**；本腿全部动作后**复归空**（⛔残留） |
| 本腿只读纪律 | 全程未 push；未碰票面／台账／`internal/tools`／`cmd/wisp`／`go.mod` | 符合（反向判据为受控临时摘⇒逐字节还原，见下） |

## 逐格 AC（对票面 `:35-41`，全部本腿真跑）

| AC | 点名尺（本腿实跑命令） | 读数 | 判定 |
|---|---|---|---|
| **AC#1** 包内容＝预览所述；密钥／转写／音频默认不在；opt-in 逐条＋逐条预览 | `go test ./internal/observe -run 'TestAC1BundleContents' -count=1 -v` | **1 PASS／0 FAIL（0.05s）**。用例内部断言本腿逐行核过：默认包无 `transcripts/`·`watchdog/`·`audio*` 前缀条目（`:642`）；opt-in 转写逐条预览＋正文＝脱敏形状（`:689-724`）；七枚畸形名两路径皆红（`:730`） | **PASS** |
| **AC#2** 磁盘闸在阈值跳、L3+日志停、活任务不受影响、腾出即恢复 | `-run 'TestAC2DiskGuard' -v` | **3 PASS／0 FAIL（0.041s）**：跳闸拒三分类＋`essential` 恒放行单列；警告文案；恢复即判；探针失败＝fail-closed 拒写不冒充 | **PASS** |
| **AC#2 附**（票面 `:29` 压力期日志不丢最后一条错误行） | `-run 'TestGatedLogs' -v` | **2 PASS／0 FAIL（0.042s）**：held>0／lost=0、恢复先重放、尾满丢最旧计数、ENOSPC 回尾 | **PASS** |
| **AC#3** WAL 持续写下有界＋动作点名 | `-run 'TestAC3Wal' -v` | **1 PASS／0 FAIL（0.03s）**：越界首读点名 TRUNCATE、半界报 AUTO、末读不衰减、无 -WAL＝健康、父目录不可达记 `ReadingFail` 不算 0 字节 | **PASS** |
| **AC#4** Focus-Assist 抑制时 badge+task list 仍更新 | `git grep -iln 'FocusAssist' -- internal cmd` | **0 命中（grep exit 1）**——全仓（不止 observe）无人落。45-r1 自报「未落地、上报编排者另派写腿（B 槽）」**与盘上实况一致**：自报诚实，但这条 AC **仍未完成** | **未落地（不计 45-r1 虚假交付；计票 45 未闭环）** |
| **AC#5** artifacts 500MB LRU 超额驱逐＋逐枚记日志 | `-run 'TestAC5Artifacts' -v` | **1 PASS／0 FAIL（0.01s）**：最旧先逐、件真没了、`BytesFreed`／`TotalAfter≤quota`、`quota<=0` 取默认不关闸、目录绕过不记 Refused | **PASS** |

## 反向判据（本腿亲手摘，⛔照抄 45-r1 的报告）＋还原核

**第一枚（题面点名的守卫摘除）**：`guards_45.go:189` `low := free < g.threshold` ⇒ `low := false`。
- 摘后点名尺 `-run 'TestAC2DiskGuard|TestGatedLogs'`：**3 FAIL**，头句逐字 `guard did not stop write class "log" at 199MiB free (threshold is 200MiB)`、`after the recheck window the guard must look again and trip`、`the guarded write reached the file while the guard was tripped`——与 45-r1 报告引的头三句**逐字对得上**。
- 再全量摘跑 `-count=1 -v`：**恰好 4 枚用例红**（三枚点名外 `TestAC1BundleContents` 因包预览含 guard 读数**连坐红**）＝45-r1 所称「4 枚用例红」复现吻合。
- 还原：`Copy-Item` 自 `D:\tmp\wisp\20261103\45v1\guards_45.go.orig` ⇒ SHA256 `0C18590E1B2C398D3527E0A486AF91CD4F60029E7E24065BC1BBB0089FA151B3` 与摘前**逐字节同**；`git diff --stat` 空；同批用例 `ok 0.132s` 回绿。

**第二枚（本腿自加，防断言自造洞）**：摘 `CheckWalFile` 越界分支（`> boundBytes` ⇒ `&& false`）⇒ `TestAC3Wal` 红 2 句（点名 TRUNCATE 断言＋"no decay" 断言）；摘 `EnforceArtifactsQuota` 驱逐循环（`range files` ⇒ `range files[:0]`）⇒ `TestAC5Artifacts` 红 3 句（`before=300 after=300`、`evicted ""`、`BytesFreed=0`）。
- 还原：同哈希 `0C18590E…`；`git status --porcelain -- internal/observe` **空**；`-run 'TestAC3Wal|TestAC5Artifacts'` 回绿 `ok 0.065s`；全量终跑 `ok 4.319s` FAIL=0。
- ⇒ **五格 AC 的断言全部有牙齿**：绿不是断言自己造的洞。

## 纪律核（`git show 39cf67751`）

- stat＝**5 件 1847 插入**：`guards_45.go`(614, 新)＋`guards_45_test.go`(772, 新)＋`diagnostics.go`(+405, 改)＋本腿证据件＋派单存档——**写面＝internal/observe＋证据＋存档，与题面逐枚符**；⛔cmd／⛔internal/tools／⛔internal/agent·memory／⛔plugin／⛔thresholds·golden·go.mod·go.sum（全部不在 stat）。
- `go vet ./internal/observe` ⇒ exit 0 静默。
- 禁形：`go func` 0 命中；`filepath.Clean|Abs` 0 命中（只 `Join`/`Dir`）；新三件无 emoji；`git log 39cf67751..HEAD -- internal/observe` 空 ⇒ 交付后无他腿动过本包。
- 200MB/500MB/8MiB 系票面自给的写入准入数，不在任何 SLO／`thresholds.go` 门文件里（stat 不含即证）。

## 本腿独立看到、须记账在两处（非 45-r1 之过，是票 45 的剩余面）

1. **装配缝未接**（45-r1 已自报）：`NewDiskGuard`／`NewDiskGatedWriter`／`EnforceArtifactsQuota`／`CheckWalFile`／`PreviewDiagnosticsBundle` 在 `internal/observe` **之外零调用点**（`git grep … -- ':!internal/observe' ':!docs' ':!.scratch'` ⇒ exit 1）。守卫在包内是真的，运行时还没人按开关——与 44-r5「配额闸装进 cmd/wisp 派发链」同一分法，属下一枚装配写腿。
2. **AC#4 全仓未落地**（45-r1 已自报且如实上报 B 槽）：写面在通知／球／面板侧，本波无人认领。

## 结论：是否同意结 45

**不同意结票 45。**
- 对 **45-r1 交付件本身**：五格（AC#1/#2/#2附/#3/#5）**全 PASS**、反向判据双枚摘还原全数吻合、纪律核与题面逐字符、两处留账自报与盘上实况一致（诚实）⇒ **该交付本身验收通过，建议编排者照单收下（45-r1 可结）**。
- 但**票 45 整体未闭环**：票面五条 AC 中 **AC#4 全仓 0 命中**、且守卫族**装配缝零调用点**（票的 AC 尺是产品行为面「fixture 跳闸／恢复」「驱逐有日志」在包内已达，但"stop L3/log/artifacts writes"要真装在运行链上才算票面 `:21-23` 落尽）⇒ **须另派（一）Focus-Assist／双通道降级 B 槽写腿、（二）cmd 侧装配腿（含 GetDiskFreeSpaceEx 探针注进 `NewDiskGuard`＋`InitLog` 包 `NewDiskGatedWriter`），两腿交付并过非实现者验收后，方可结 45**。
- ⛔ 本裁决「不同意结票」不是打回 45-r1：它是本票迄今唯一可验收的实交付，质量核（有牙齿的断言＋fail-closed＋诚实留账）本腿全部复现。

（45-v1 完；全程只读＋受控摘还原；零残留；⛔push。）
