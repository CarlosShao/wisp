# `10-roster.md` — 表①：`%TEMP%` 下 `Test*` 目录的**逐枚**归属名册（票 310 `AC#0`，腿 `310-a1`）

锚＝`49153fa8`（产码面读数在此尖上量，见 `00-anchor.md` §1/§9）。本程⛔ 删任何目录，本表**只用于逐枚处置**，⛔ 由它生成"可以批量删"的结论（票面 `:12` 自己写了"⛔ 按名字前缀就直接删"）。

## 1. 尺与分母（三个数一条链，同批给）

| 层 | 尺 | 读数 | 自落 rc |
|---|---|---|---|
| 分母全量 | `powershell (Get-ChildItem -Path $env:TEMP -Directory).Count` | **2754** 枚目录 | `temp_count_rc=0`（件＝`logs/00-gate-310a1.txt`） |
| 本票射程 | `Get-ChildItem -Path $env:TEMP -Directory \| Where-Object { $_.Name -like 'Test*' }` | **50** 枚（`%TEMP%`＝`C:\Users\swq\AppData\Local\Temp`） | 同上 |
| 其中"本仓可证" | 下表逐枚反查（三形同试，见 §2） | **49** 枚 | `rows=50` 尺自 `logs/10-roster-enriched.tsv` |
| 其中"本仓不可证" | 同上 | **1** 枚（`TestProbeStepByStep3843095677`） | 同上 |

支撑尺（全部同批跑，件在 `logs/`）：

- `logs/10-temp-test-dirs.tsv`＝50 行原始读数（名／mtime／内层文件枚数／内层字节数），尺＝`Get-ChildItem … -like 'Test*'` 逐枚 `Get-ChildItem -Recurse -Force -File \| 求和`；`wc -c`=3892 非 0。
- `logs/10-repo-testfuncs.txt`＝本仓测试函数全表，尺＝`git grep -n -E '^func Test[A-Za-z0-9_]*\(' HEAD -- cmd internal tools scripts` → `func_lines=1686`、`distinct_names=1683`；`grep_rc=0`。
- `logs/10-file-firstadd.tsv`＝每个命中文件在本仓的最早提交日（`git log --reverse --format=%ad --date=short -- <file>`），19 枚文件（含 1 行 `-`＝NOMATCH 占位）。
- 汇总尺：`total_bytes=416656144`（≈397 MB）／`total_inner_files=145`／空目录＝**10** 枚／非空＝**40** 枚／mtime 按日分布＝09-21:4、09-23:1、09-25:5、10-01:1、10-02:1、10-03:3、10-05:1、10-07:3、10-08:3、10-09:1、10-11:**29**（`rc=0`，件＝本文件末的行计数尺输出 `logs/` 同批）。

⚠ 本表**没有**递归进目录里数文件以外的东西（⛔ 全量递归列表进对话），只取每枚的聚合数。

## 2. 反查尺的口径（派单那条"去尾数＋git grep"不够用，本腿的改正具名报回）

派单表① 给的尺＝"把去掉尾部随机数字后的名字拿去 `git grep -n "func <名>" HEAD -- cmd internal tools scripts`"。本腿现量：**这一把尺漏 3 形**，理由＝Go 1.27.1 的 `t.TempDir()` 命名机制（provenance＝仓外 GOROOT 源码 `D:\work\base\go\src\testing\testing.go:1584-1633` ＋ `D:\work\base\go\src\os\tempfile.go:22-24,86-104`，本腿逐字读过）：

1. **目录名＝`c.Name()` 截前 64 字节 → `removeSymbolsExcept(pattern, "!#$%&()+,-.=@^_{}~ ")` → `os.MkdirTemp` 追加随机数**。`removeSymbolsExcept` 对不允许的符号是**删掉**（`return -1`），**不是**换成 `_` ⇒ 子测试的 `/` 分隔符直接消失（`TestX/子名` → 目录名 `TestX子名…`）。⇒ "去尾数后拿原名 `git grep` 精确命中"对**子测试里调的 `t.TempDir()`** 必零命中。本腿实测这一形命中 **1** 枚。
2. **合法以数字结尾的用例名会被"去尾数"过剥**：`nextRandom()`＝`strconv.FormatUint(uint32,…)`（1–10 位），而用例名本身可含数字 ⇒ 贪婪剥尾数会把名字的一部分也剥掉。本腿实测命中 **1** 枚（`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` ＋ 随机数 `1963646102`）。
3. **64 字节截断**：用例名（含子测试名）超过 64 字节时目录名只留前 64 字节 ⇒ 反查必须用**前缀**关系（目录 stem 是仓库用例名的前缀）而不是相等关系。本批 50 枚里**没有**落在这一形（0 枚；尺＝下表 `truncated` 列，见 §5 的 `form:` 计数行），但下一批会有——本仓用例名普遍带票号、越长越常见。

⇒ 本腿把一把尺换成**三形同试＋前缀双向**（`logs/10-roster-enriched.tsv` 的 `form` 列）：
① `exact`＝剥尾数后与某枚 `func Test…` **逐字相等**；② `digits_in_name`＝剥尾数过剥、但仓库某用例名**以该 stem 为前缀**（真名更长）；③ `subtest`＝该用例名**是 stem 的前缀**（stem 里多出子标签）。
每枚另外量一把**独立于名字**的第二尺：`mtime` 日 ≥ 命中文件在本仓的最早提交日（`after_firstadd`）——名字是线索、mtime 是第二根棍子。

## 3. 三态判定口径（本腿下的定义，⛔ 含糊）

- **本仓可证**＝三形之一命中 HEAD 里的 `func Test…`（给出文件:行）**且** mtime ≥ 该文件本仓首提交日。本批＝**49** 枚（47 `exact` ＋ 1 `digits_in_name` ＋ 1 `subtest`）。
- **别家（本仓不可证）**＝名字在本仓**三把尺全零命中**：① `git grep` HEAD（0）② 工作树 `git grep`（0）③ 全历史 `git log --all -S<名> -- '*.go'`（0）。本批＝**1** 枚。
- **归属不明**＝本批 **0** 枚，但**这一枚本腿要说清它为什么没落在这里**：派单给的"别家"判据只有"名字零命中"这一条，而"证不了是我们的"≠"证得了是别家的"。`TestProbeStepByStep3843095677` 的 mtime＝2026-09-21 18:40:57，**落在我们的跑测窗口之内**（该日本仓有 **357** 笔提交，`git log --since/--until` 现量，`commits_on_20260921=357`；09-20＝168 笔）⇒ **第二根棍子失效**，只能靠名字。⇒ 本腿按派单口径标它"别家（本仓不可证）"，⛔ 指不出是哪家；若编排者要的是"确证他源"这一档，这一枚应改落**归属不明**（已登记为 `Q-310-a1-1`，见 `90-final.md`）。

## 4. 逐枚名册（50 行，⛔ 只报枚数）

列义：目录名／mtime（本地 `+08`）／内层（文件枚数、字节数，`0 files`＝空目录）／反查命中（形＋用例名＋`文件:行`）／归属三态／理由。
机器生成源＝`logs/10-roster-enriched.tsv`（50 行，`wc -c` 非 0）；下表由该件逐行投影，未手改任何一个字段。

| 目录名 | mtime（本地 +08） | 内层（文件枚数／字节） | 反查命中的本仓用例（形＋用例名＋`文件:行`） | 归属三态 | 可证／不可证的理由（⛔ 空） |
|---|---|---|---|---|---|
| Test221ParentStopsItsOwnChildRowAndStreamSettle3871214707 | 2026-10-08 20:08:20 | 3 files / 208896 bytes | `Test221ParentStopsItsOwnChildRowAndStreamSettle` 形exact（internal/tools/task_cancel_221_legs_test.go:287） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝Test221ParentStopsItsOwnChildRowAndStreamSettle，命中 internal/tools/task_cancel_221_legs_test.go:287；mtime 2026-10-08 ≥ 该文件本仓首提交 2026-09-29 ⇒ 两尺同向 |
| TestAC1ResidentLegOutlivesItsOwnLogFailure1462165283 | 2026-10-09 15:35:46 | 空目录（0 files） | `TestAC1ResidentLegOutlivesItsOwnLogFailure` 形exact（cmd/wisp/resident_sink_nail_127_windows_test.go:520） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestAC1ResidentLegOutlivesItsOwnLogFailure，命中 cmd/wisp/resident_sink_nail_127_windows_test.go:520；mtime 2026-10-09 ≥ 该文件本仓首提交 2026-09-22 ⇒ 两尺同向 |
| TestAC2RealProcessRefusesOnEveryLegWithoutAppData1281963646102 | 2026-10-03 17:04:22 | 1 files / 31092893 bytes | `TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` 形digits_in_name（cmd/wisp/dataroot_128_windows_test.go:49） | 本仓可证 | 形②名内含数字：贪婪剥尾数把用例名自身的数字也剥了（stem=1）；真名以 stem 为前缀＝TestAC2RealProcessRefusesOnEveryLegWithoutAppData128（随机数尾段 1963646102，10 位 ≤ uint32 的 10 位上限），命中 cmd/wisp/dataroot_128_windows_test.go:49；mtime ≥ 首提交 2026-09-23 |
| TestComposedGateBlocksAWriteForTwoSeconds1996183427 | 2026-09-23 18:44:49 | 1 files / 10094592 bytes | `TestComposedGateBlocksAWriteForTwoSeconds` 形exact（cmd/wisp/run_test.go:334） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestComposedGateBlocksAWriteForTwoSeconds，命中 cmd/wisp/run_test.go:334；mtime 2026-09-23 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestConfirmLockedRunsOutsideTheManagerLock1562080786 | 2026-10-02 19:39:40 | 1 files / 198 bytes | `TestConfirmLockedRunsOutsideTheManagerLock` 形exact（internal/config/manager_223_test.go:55） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestConfirmLockedRunsOutsideTheManagerLock，命中 internal/config/manager_223_test.go:55；mtime 2026-10-02 ≥ 该文件本仓首提交 2026-09-29 ⇒ 两尺同向 |
| TestDoubleFailureKeepsTaskContextAndIsResumable481691178 | 2026-10-11 18:47:00 | 1 files / 10094592 bytes | `TestDoubleFailureKeepsTaskContextAndIsResumable` 形exact（internal/llm/fallback_test.go:140） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestDoubleFailureKeepsTaskContextAndIsResumable，命中 internal/llm/fallback_test.go:140；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestFailoverAfterPrimaryExhaustion2419959198 | 2026-10-11 21:46:05 | 2 files / 20188160 bytes | `TestFailoverAfterPrimaryExhaustion` 形exact（internal/llm/fallback_test.go:86） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestFailoverAfterPrimaryExhaustion，命中 internal/llm/fallback_test.go:86；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestFailoverAfterPrimaryExhaustion3970703700 | 2026-10-11 21:46:14 | 2 files / 20188160 bytes | `TestFailoverAfterPrimaryExhaustion` 形exact（internal/llm/fallback_test.go:86） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestFailoverAfterPrimaryExhaustion，命中 internal/llm/fallback_test.go:86；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestFailoverAfterPrimaryExhaustion455177940 | 2026-10-11 18:53:05 | 2 files / 20189184 bytes | `TestFailoverAfterPrimaryExhaustion` 形exact（internal/llm/fallback_test.go:86） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestFailoverAfterPrimaryExhaustion，命中 internal/llm/fallback_test.go:86；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestGateJudgesThePrivateSetByResolvedSID2396759740 | 2026-09-21 18:40:11 | 空目录（0 files） | `TestGateJudgesThePrivateSetByResolvedSID` 形exact（internal/winsec/private_set_sid_windows_test.go:128） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestGateJudgesThePrivateSetByResolvedSID，命中 internal/winsec/private_set_sid_windows_test.go:128；mtime 2026-09-21 ≥ 该文件本仓首提交 2026-09-21 ⇒ 两尺同向 |
| TestGateJudgesThePrivateSetByResolvedSID2459073988 | 2026-09-21 18:40:36 | 空目录（0 files） | `TestGateJudgesThePrivateSetByResolvedSID` 形exact（internal/winsec/private_set_sid_windows_test.go:128） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestGateJudgesThePrivateSetByResolvedSID，命中 internal/winsec/private_set_sid_windows_test.go:128；mtime 2026-09-21 ≥ 该文件本仓首提交 2026-09-21 ⇒ 两尺同向 |
| TestGateJudgesThePrivateSetByResolvedSID2986811054 | 2026-09-21 18:40:05 | 空目录（0 files） | `TestGateJudgesThePrivateSetByResolvedSID` 形exact（internal/winsec/private_set_sid_windows_test.go:128） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestGateJudgesThePrivateSetByResolvedSID，命中 internal/winsec/private_set_sid_windows_test.go:128；mtime 2026-09-21 ≥ 该文件本仓首提交 2026-09-21 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents1101711636 | 2026-10-11 17:52:31 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents1385547188 | 2026-10-11 21:37:17 | 1 files / 10094080 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents1415189317 | 2026-10-11 17:49:25 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents1524086492 | 2026-10-11 16:53:41 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents1904824687 | 2026-10-11 17:56:15 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents2997322679 | 2026-10-11 21:37:03 | 1 files / 10094080 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3039020949 | 2026-10-11 21:34:10 | 1 files / 10094080 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3140925201 | 2026-10-11 18:01:49 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3161225382 | 2026-10-11 17:52:19 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3233152004 | 2026-10-11 17:56:27 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3278041641 | 2026-10-11 21:37:56 | 1 files / 10094080 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3292549456 | 2026-10-11 18:02:05 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents3666715146 | 2026-10-11 18:47:06 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents4011121267 | 2026-10-11 16:54:25 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents696038329 | 2026-10-11 18:05:38 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents718009555 | 2026-10-11 17:49:18 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmGoldenByteIdenticalEvents973367927 | 2026-10-11 18:05:47 | 1 files / 10094592 bytes | `TestMockllmGoldenByteIdenticalEvents` 形exact（internal/llm/openaichat/mockllm_integ_test.go:183） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestMockllmGoldenByteIdenticalEvents，命中 internal/llm/openaichat/mockllm_integ_test.go:183；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-19 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite1562417226 | 2026-10-11 17:56:48 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite1603513437 | 2026-10-11 18:02:32 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite222012253 | 2026-10-11 17:52:40 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite2329949686 | 2026-10-11 18:05:55 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite2868402985 | 2026-10-11 18:02:16 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite3458740087 | 2026-10-11 17:56:35 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite3784017352 | 2026-10-11 17:49:34 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestMockllmSharedFaultSuite959708402 | 2026-10-11 18:06:04 | 1 files / 10094592 bytes | `TestMockllmSharedFaultSuite` 形exact（internal/llm/anthropic/suite_test.go:42） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestMockllmSharedFaultSuite，命中 internal/llm/anthropic/suite_test.go:42；mtime 2026-10-11 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestPanelHostRealWindowHopAndLifecycle537087969 | 2026-10-01 10:04:05 | 70 files / 810867 bytes | `TestPanelHostRealWindowHopAndLifecycle` 形exact（cmd/wisp/panel_host_windows_test.go:614） | 本仓可证 | 形①exact：剥尾随机数（9 位）后与 HEAD 用例名逐字相等＝TestPanelHostRealWindowHopAndLifecycle，命中 cmd/wisp/panel_host_windows_test.go:614；mtime 2026-10-01 ≥ 该文件本仓首提交 2026-10-01 ⇒ 两尺同向 |
| TestProbeStepByStep3843095677 | 2026-09-21 18:40:57 | 空目录（0 files） | `-` 形NOMATCH（-:-） | 别家＝本仓不可证 | 名字三把尺全零命中：HEAD `git grep`＝0／工作树 `git grep`＝0／全历史 `git log --all -S`（*.go）＝0；空目录（0 files）；mtime 2026-09-21 落在我方跑测窗口之内（当日 357 笔提交）⇒ mtime 尺失效，只剩名字一尺；⛔ 指不出是哪家 ⇒ 若口径要"确证他源"则应改落归属不明（`Q-310-a1-1`） |
| TestRollingWriterRetentionSweep3198295670 | 2026-09-25 08:32:26 | 空目录（0 files） | `TestRollingWriterRetentionSweep` 形exact（internal/observe/logging_test.go:188） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestRollingWriterRetentionSweep，命中 internal/observe/logging_test.go:188；mtime 2026-09-25 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestRollingWriterRetentionSweep3638913217 | 2026-09-25 08:34:09 | 空目录（0 files） | `TestRollingWriterRetentionSweep` 形exact（internal/observe/logging_test.go:188） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestRollingWriterRetentionSweep，命中 internal/observe/logging_test.go:188；mtime 2026-09-25 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestRollingWriterRetentionSweep3835089158 | 2026-09-25 08:34:11 | 空目录（0 files） | `TestRollingWriterRetentionSweep` 形exact（internal/observe/logging_test.go:188） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestRollingWriterRetentionSweep，命中 internal/observe/logging_test.go:188；mtime 2026-09-25 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestRollingWriterRetentionSweep3838303947 | 2026-09-25 08:34:24 | 空目录（0 files） | `TestRollingWriterRetentionSweep` 形exact（internal/observe/logging_test.go:188） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestRollingWriterRetentionSweep，命中 internal/observe/logging_test.go:188；mtime 2026-09-25 ≥ 该文件本仓首提交 2026-09-20 ⇒ 两尺同向 |
| TestRunPacketCarriesTheSubagentItsRosterRowFed4174395780 | 2026-10-08 11:35:29 | 9 files / 10337094 bytes | `TestRunPacketCarriesTheSubagentItsRosterRowFed` 形exact（cmd/wisp/subagent_carrier_197_test.go:231） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestRunPacketCarriesTheSubagentItsRosterRowFed，命中 cmd/wisp/subagent_carrier_197_test.go:231；mtime 2026-10-08 ≥ 该文件本仓首提交 2026-09-28 ⇒ 两尺同向 |
| TestSLO144LoopGiveUpSentencesOnRealFiles2237962421 | 2026-09-25 20:36:08 | 空目录（0 files） | `TestSLO144LoopGiveUpSentencesOnRealFiles` 形exact（cmd/wisp/slo_report_144_windows_test.go:221） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestSLO144LoopGiveUpSentencesOnRealFiles，命中 cmd/wisp/slo_report_144_windows_test.go:221；mtime 2026-09-25 ≥ 该文件本仓首提交 2026-09-25 ⇒ 两尺同向 |
| TestTicket101SessionGrantDoesNotCrossRestart2995253497 | 2026-10-07 18:23:33 | 8 files / 10313209 bytes | `TestTicket101SessionGrantDoesNotCrossRestart` 形exact（cmd/wisp/run_mode101_test.go:391） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestTicket101SessionGrantDoesNotCrossRestart，命中 cmd/wisp/run_mode101_test.go:391；mtime 2026-10-07 ≥ 该文件本仓首提交 2026-09-21 ⇒ 两尺同向 |
| TestTicket223PermissionDeniedSitsInItsOwnSentence1318611284 | 2026-10-05 23:25:53 | 9 files / 10340275 bytes | `TestTicket223PermissionDeniedSitsInItsOwnSentence` 形exact（cmd/wisp/config_reload_perm_223_windows_test.go:60） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestTicket223PermissionDeniedSitsInItsOwnSentence，命中 cmd/wisp/config_reload_perm_223_windows_test.go:60；mtime 2026-10-05 ≥ 该文件本仓首提交 2026-09-29 ⇒ 两尺同向 |
| TestTicket223R2FailureSentenceRouting声明旧版_坏表头_A2_3165385141 | 2026-10-03 17:01:11 | 1 files / 35 bytes | `TestTicket223R2FailureSentenceRouting` 形subtest（cmd/wisp/config_sentences_223r2_test.go:39） | 本仓可证（父用例级） | 形③子测试：stem＝父用例名＋子标签（Go 1.27 的 removeSymbolsExcept 把 `/` **删掉**、⛔ 换成 `_`），父名命中 cmd/wisp/config_sentences_223r2_test.go:39；⚠ 子标签逐字拼法本腿未复核（不影响归属，只影响能否再往下指到表行）；mtime ≥ 首提交 2026-09-29 |
| TestTicket223TighteningRaisesNoCard3267158827 | 2026-10-03 17:09:55 | 1 files / 10094592 bytes | `TestTicket223TighteningRaisesNoCard` 形exact（cmd/wisp/config_reload_223_test.go:503） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestTicket223TighteningRaisesNoCard，命中 cmd/wisp/config_reload_223_test.go:503；mtime 2026-10-03 ≥ 该文件本仓首提交 2026-09-29 ⇒ 两尺同向 |
| TestUnansweredL2CardTimesOutIntoRejectNeverExecution2689719599 | 2026-10-08 08:42:21 | 9 files / 10340645 bytes | `TestUnansweredL2CardTimesOutIntoRejectNeverExecution` 形exact（cmd/wisp/approval_reply_201_test.go:434） | 本仓可证 | 形①exact：剥尾随机数（10 位）后与 HEAD 用例名逐字相等＝TestUnansweredL2CardTimesOutIntoRejectNeverExecution，命中 cmd/wisp/approval_reply_201_test.go:434；mtime 2026-10-08 ≥ 该文件本仓首提交 2026-09-29 ⇒ 两尺同向 |

## 5. 计数自落（本表自己的尺，逐条 rc）

```
rows_total=50            (grep -c '^| Test' -> rc=0)
form_exact=47            (awk -F'\t' '$5=="exact"' logs/10-roster-enriched.tsv | grep -c '' -> rc=0)
form_digits_in_name=1    (同上 $5=="truncated64" 这一枚经 §2 复核后改名)
form_subtest=1
form_never_seen_in_repo=1  ($5=="NOMATCH")
cons_after_firstadd=49   (第二根棍子：mtime 日 ≥ 命中文件本仓首提交日)
cons_BEFORE_firstadd=0   (empty above＝0 枚 ⇒ 没有任何一枚出现"目录比用例还早"的矛盾)
attribution_本仓可证=49
attribution_别家_本仓不可证=1
attribution_归属不明=0    (empty above＝0 枚)
sum_bytes=416656144  sum_inner_files=145  empty_dirs=10  nonempty_dirs=40
rc=0
```

- ⚠ 上面这几把计数尺凡命中 0 都写了"empty above＝0 枚"，⛔ 空输出当读数。
- 名册逐枚的机器源＝`logs/10-roster-enriched.tsv`（50 行 × 10 列，`wc -c` 非 0），表里⛔ 手改字段；唯一的文字更正＝那一枚 `form=truncated64` 的标签按 §2 复核结果改叫 `digits_in_name`。

## 6. 这张表能用来干什么／不能干什么

- **能**：给下一步**按名册逐枚**处置用（每一枚都带 `文件:行` 与 mtime 两把独立证据），以及给"这 50 枚到底是 1 枚用例留的还是 40 枚用例留的"这个问题提供分布。
- **不能**：① ⛔ 由它生成"按 `Test*` 前缀批量删"的结论——票面 `:12` 与本腿 §3 那枚零命中目录都说明名字前缀不等于归属；② ⛔ 把"49 枚本仓可证"读成"49 枚可以删"：可证归属只是删的**必要条件**，票面 `:33` 还要求先经非实现者复核；③ ⛔ 拿"命中 `文件:行`"当"就是那条用例跑的"——名字能证的是**这个仓库里有这条用例**，不能证"这一枚目录是哪次运行、哪条腿留下的"（本批没有任何一件运行日志能把 mtime 对上具体一发）。
- 第二根棍子的读数：`cons_BEFORE_firstadd=0`（empty above＝0 枚）⇒ 本批**没有**一枚出现"目录 mtime 早于所命中用例文件在本仓的首提交日"这种矛盾；这一形一旦在下一批出现，那一枚就该从"本仓可证"降级重判（改名／搬文件会让本仓首提交日晚于目录 mtime，所以这条尺只证伪、不证真）。
- 分布补充（同表可核）：50 枚里 **29** 枚 mtime 在 2026-10-11、**4** 枚在 09-21、**5** 枚在 09-25 ⇒ 残留的主体是"这两天编队跑测"的产物，⛔ 均匀老化；最大三枚＝`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128…`（31 092 893 B）、`TestFailoverAfterPrimaryExhaustion…455177940`／`…2419959198`（各 ≈20 MB）⇒ **单枚能到 30 MB 量级**，这是表③ 代价尺的分母之一。
