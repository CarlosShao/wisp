# 70-ac2-testcore-triage —— test-core（ubuntu）step 7 失败逐条分诊（只读腿 70-ac2-triage）

- 腿：`70-ac2-triage`（只读诊断；⛔ 未改任何 .go / scripts / ci.yml；⛔ 未跑 go test；⛔ 未 push）
- 证据源：CI run `43146479546`，test-core job `130945374182`，失败步 = step 7 `Portable package tests`；
  原始日志 `D:\tmp\wisp\20261112\tc-log.zip`（1,843,239 字节纯文本）现摘。
- 现取四数（日志 :12321-:12322，runtests.sh + portable-tests.sh 各一行）：
  `=== RUN=3037  --- PASS=2037  --- FAIL=15  --- SKIP=4`，`go test exited 1`。
- 结论先行：**15 条 FAIL = internal/tools 9 条 + internal/plugin 6 条。A 类（修一处即可移植，测试保留 core scope）11 条；B 类（声明出 core scope，ledger 具名登记）4 条。**
- 唯一"产品侧一处修好七条"的根子：`internal/tools/fs.go:280-288` 的 `joinForListing` 在非分隔符结尾的目录后**硬编码拼 `\`**（`return dir + "\" + name`）。POSIX 上拼出的路径不存在 → `os.ReadDir`/`os.ReadFile` 静默 `continue`（search_content.go:293-296、search_files.go:201-204）→ 子目录进不去、文件读不到、命中 0。这与票 70 早前记录的 `internal/memory` 那条（"反斜杠在 Linux 上是文件名而不是目录分隔符"，`canary-separator.txt`）是**同一族**。

## 分诊表（15 条）

| 序号 | 测试名 | 包 | 断言 file:line | 根因类型 | A/B | 具体动作 |
|---|---|---|---|---|---|---|
| 1 | Test24R2ASearchContentFindsMatches | internal/tools | document_text_test.go:449 | 原生路径分隔符（joinForListing 硬编码 `\`，fs.go:287 → ReadFile 全失败 → matches=0） | A | 修 `internal/tools/fs.go:280-288`：`return dir + string(os.PathSeparator) + name`（开头后缀判断已兼容 `/`）。⛔ 别用 `filepath.Join`——D22 禁在 PathResolver 之外用 Clean/Abs 做文件系统决策。测试不动 |
| 2 | Test24R2ASearchContentNeverListsBlacklistedFiles | internal/tools | document_text_test.go:496（:493 excluded 同读数） | 同上（sub/c.log 进不去、正文读不到 → matches=0） | A | 同序号 1，一处修七条；测试不动 |
| 3 | Test24R2ASearchContentCapsAndSkips | internal/tools | document_text_test.go:517 | 同上（日志实测 `scanned:4`——只有根目录 4 个文本文件，`notes/` 子目录整个没下去） | A | 同序号 1；测试不动 |
| 4 | Test24R2ASearchContentThroughResolverReportsTruncated | internal/tools | document_text_test.go:611、:617 | 同上（命中 0 → Truncated/header 断言连带红） | A | 同序号 1；测试不动 |
| 5 | Test24R3SearchFilesMatchesNameGlobCaseInsensitive | internal/tools | search_files_test.go:77、:88 | 同上（日志实测列表项是 `/tmp/…/001\B.MD` 混合分隔符；`sub` 的 frame 路径带 `\` → ReadDir 失败 → `sub/c.md` 永不出现，`扫描 5 项` 恰好数不出 c.md） | A | 同序号 1；测试不动 |
| 6 | Test24R3SearchFilesDepthCapBothWays | internal/tools | search_files_test.go:274 | 同上（l1 的 frame 路径带 `\` → ReadDir 失败 → 控制臂 cap 8 也只剩 `扫描 1 项`，深度帽测试测不到帽本身） | A | 同序号 1；测试不动 |
| 7 | Test24R3SearchFilesCancelMarksPartialResult | internal/tools | search_files_test.go:338 | 同上＋连带效应：Classify 收到的是 `/tmp/…\deeper`，`filepath.Base`（POSIX 不认 `\`）取不回 `deeper` → 取消触发器永不命中 → 无"已中止" | A | 同序号 1（join 修好后 Base/取消触发自然恢复）；测试不动 |
| 8 | Test22AppLaunchJobFailureSurfaces | internal/tools | web_launch_test.go:122 | 实现缺口（非 Windows 分支）：`launcher_other.go:72-75` 的 startApp 在 `Jobs.StartInJob` 失败时**裸透传 err**，结果文本无 "Job" 措辞；Windows 侧 `launcher_windows.go:129-134`（startInJobOrBare）有 `进入 Job Object 失败（C30）：%w` | A | 修 `internal/tools/launcher_other.go` startApp 的 Jobs 分支：把 err 包成与 Windows 侧同形的 `进入 Job Object 失败（C30）：%w`（测试用 fakeJobs，跨平台可接线，断言本身可移植） |
| 9 | Test22AUMIDIsWindowsOnlyHere | internal/tools | web_launch_test.go:160 | AUMID 仅 Windows：AUMID 形状判定 `windowsAUMIDPattern` 定义在 `launcher_windows.go:38`（`//go:build windows`，非 Windows 编译不可见）；`launcher_other.go:50-66` 的 resolveApp **没有 aumid 臂** → AUMID 形状落进名称解析 → `找不到应用：…`，而非 startApp :83-84 本已备好的平台拒绝"UWP/AUMID 启动是 Windows 平台能力（D7），本平台不支持" | A | 把 AUMID 形状判定（`windowsAUMIDPattern` 或等价判定）挪到平台无关文件；`launcher_other.go` resolveApp 命中形状 → 返回 `Kind:"aumid"`，:83-84 的拒绝臂原样接住。测试不动（它的本意就是钉"非 Windows 必须答平台拒绝"） |
| 10 | Test52PinComputesAndVerifiesHash | internal/plugin | command_52_test.go:101 | D46 pin 通道 Windows 形状：`internal/plugin/command.go:160-161` 用 `filepath.VolumeName(canonical)==""` 判"drive-less path"，POSIX 绝对路径（/tmp/…）必拒；同通道还有 :98 `commandExeExts` 只收 `.exe/.com`、:157 UNC 判据——整个 exe 稳定性判据（盘符＋可执行镜像扩展名）是 Windows 概念 | B | 在 `scripts/portable-tests.sh:803` ledger 具名登记（platform=linux, class=fixture，name+pkg+reason 五字段）：**D46 的 exe pin 判据按 PLAN.md:2678/SPEC-07 §4 是盘符形状，POSIX 上测试无主体**；`internal/plugin` 同时在 win_pin（portable-tests.sh:396）里，windows 腿照跑，断言一条不丢。不建议 A：为 D7=Windows 的产品功能给 resolveExe 发明 POSIX pin 语义是造契约，不是修移植 |
| 11 | Test52OneByteTamperRefusesAndAlerts | internal/plugin | command_52_test.go:131 | 同序号 10（PinCommandExe 在 resolveExe 第一步就拒） | B | 同序号 10（ledger 四条可共用一条 reason，逐条具名） |
| 12 | Test52ManifestPinMismatchIsRefusedAsASwap | internal/plugin | command_52_test.go:181 | 同序号 10 | B | 同序号 10 |
| 13 | Test52ManifestPinsEveryCommandTool | internal/plugin | command_52_test.go:335 | 同序号 10 | B | 同序号 10 |
| 14 | Test52SchemaViolationsAreRejected（红的是 3 个子用例：below_minimum / above_maximum / not_an_integer） | internal/plugin | command_52_test.go:464 | **其它（平台无关，非移植性问题）**：fixture 与"schema 收窄到模板"契约自相矛盾——`withProps`（:430-438）只换 `def.Parameters`，`ArgvTemplate` 仍是 base 的 `{{to}}/{{text}}`（:44-76），`days` 占位符不存在 → `ScreenCommandParams` 只筛模板占位符（command_argv.go:114-130），`screenValue` 的 integer/min/max 判定（command_argv.go:149-160、:298-299）根本不经手；`required:[days]` 又恰好满足 → `{"days":0}` 被如实接受 | A | 修 fixture（command_52_test.go:450-452 三个子用例）：给 `def.ArgvTemplate` 追加 `"--days", "{{days}}"`（或 withProps 增加模板参数），让 days 真正流过 screenValue；断言一行不动。⚠ 这 3 条与 ubuntu 无关（纯函数路径），windows 腿同码应同红——建议编排者顺带复核 windows 腿最近读数 |
| 15 | Test52CommandKindValidates（红的是 unknown_kind 子用例） | internal/plugin | command_52_test.go:1002 | **其它（平台无关）**：`kind = "shell"` 使 auditKeys（manifest.go:228-233）按 declarative 键表审计，D46 专有键（exe/argv_template/extract/…）逐个报 `tools[0].exe … unknown key`；ParseManifest（manifest.go:178-183）在 audit 阶段就返回，validator.go:206 那个真正报 `path+".kind"` 的分支**不可达** → 错误串里永远没有 ".kind" | A | 推荐实现侧最小修：auditKeys/auditTable 检出 `kind` 非两法值时**优先**追加 `tools[N].kind` 的 finding（manifest.go:228 一处）；备选：改测试期望 field（:978 `".kind"` → `"tools[0]"`）。两者取一由写腿定；推荐前者，因为 C14 的字段级判定本就该指到 kind |

## 总数核对（15 vs 票面曾写的 4）

- **本次现取 = 15**（逐条 `--- FAIL:` 顶层名现摘自日志 ：7494-:10647；与派单所附清单逐字一致，无补漏、无多计——`Test52SchemaViolationsAreRejected` / `Test52CommandKindValidates` 的子用例 FAIL 都并入各自顶层名，不另计数）。
- **"4" 的出处**：归档 run-37166458550（`PASS=1010 FAIL=4 SKIP=1`），现仍被 `scripts/portable-tests.sh:815` 的 ledger 行引用。那是**更早一次 run** 的数，不是本次 run 的数；本次 4 → 15 的差额就是上表 15 条的现值。两条数都对，只是指向不同 run。
- **派单提示更正一条**：`command_52_test.go:265: no TEMP in this environment` 与 `:303: native spelling … has no separator to convert` **不属于 15 条 FAIL**——它们是 `--- SKIP`（`Test52AbsoluteTempPathIsRefusedForTheTempReason` 日志 ：10513、`Test52ForwardSlashWindowsPathIsPinnable` :10516），已按 SKIP 计入 SKIP=4，不在本表。
- **SKIP=4 的构成**（日志 ：1663/:1666/:1669/:1672，顶层名）：`TestAppendSessionTurnCrashChild`（reexec 子进程入口）、`Test24R3SearchAppsAUMIDCandidateFeedsTheAUMIDArm`（search_files_test.go:458 平台守卫，skip 原文 "the AUMID shape gate is the Windows launcher"）、`Test52AbsoluteTempPathIsRefusedForTheTempReason`（ubuntu 无 TEMP 环境变量）、`Test52ForwardSlashWindowsPathIsPinnable`（原生拼写无 `\` 可转）。
- ⚠ **波及提示（交给编排者，不属本腿动作）**：这 4 条 SKIP **没有一条**在 `portable-tests.sh:803` 的 ledger 里（本轮 `-skip` 模式 ：12321 只含 8 个 ledger 名）。runtests.sh 对任何顶层 `--- SKIP` 是致命的（ticket 71 AC#3）——**即使 15 条 FAIL 全修绿，step 7 仍会因这 4 条未记账 SKIP 红**。需要按 ：22 机制具名登记（reexec 1 条 / fixture 2 条 / 平台守卫 1 条），或把 TEMP/ForwardSlash 两条改成能力探针形状。本腿只报不改。

## 处置汇总

- **A 类 = 11 条**（序号 1-9、14、15）：测试全部保留 core scope；其中序号 1-7 是**同一处产品缺陷**（`joinForListing` 一处修七条），序号 8-9 是 `launcher_other.go` 两个缺臂，序号 14-15 是 plugin 纯函数层的 fixture/错误形状二选一。
- **B 类 = 4 条**（序号 10-13，即 Test52 pin 家族整族）：从 core scope 具名声明出（ledger platform=linux, class=fixture），理由＝D46 exe pin 判据的 Windows 形状（盘符/VolumeName、.exe/.com 镜像）使 POSIX 无主体；windows 腿（win_pin 已含 internal/plugin）继续覆盖，断言零丢失。
- 处置动作全部落在**未来写腿**的工单面：本腿零代码改动、零 ledger 改动、零 push。
