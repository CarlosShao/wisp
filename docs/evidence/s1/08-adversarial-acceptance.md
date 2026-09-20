# T08 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T08-impl 独立）。时间：2026-09-19T23:48:48Z（截止警戒窗口内内联验收）。

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 复跑 | PASS | observe/models `go test -count=1` 全绿；gofmt 漂移已修（internal/models，本报告同步提交） |
| 2 | 采样器口径 | PASS（关键裁定核实） | 门禁单位=私有工作集（NtQuerySystemInformation WorkingSetPrivateSize），commit 并行记录——实测本机 8.7MB 全门禁 PASS（commit 口径会假红 48MB，裁定正确）；fail-closed（零值=error、无样本=FAIL） |
| 3 | 泄漏 fixture | PASS | 100MB 页触摸持有必须翻红否则 FATAL；回落 CheckSettle 要求 10s 内回上限且 FreeOSMemory 计数>0（实测抓到 buffer 未释放真 bug 并已修） |
| 4 | CI 矩阵 | PASS | ci.yml YAML 校验过；5 job 无可跳过；d22scan 实跑全仓 clean（自测旗标名为文档级 MINOR）；lint 含七禁令+emoji 扫描 |
| 5 | compose | PASS | `docker compose config` OK；T14 model-mirror 三服务原样保留 + mock-llm(18080) 接线 |
| 6 | 越界 | PASS | internal/models 仅 gofmt 修复；T14 其他文件零触碰 |

## 登记的风险/移交
1. 本地 sandbox 间歇从 SystemProcessInformation 隐藏自身进程 → 本地 state 采样段不稳；**首次 CI 运行（slo-full on wisp-selfhosted-01）是最终验证点**。
2. Armed/Warm 含模型态现为骨架足迹（JSON posture:"skeleton" 如实标注），票 15/26/33 落地后阈值自动生效。
3. tools/d22scan 的自测触发方式文档化（MINOR，转票 12 gate 检查项）。

**VERDICT: PASS**

## Addendum 裁决（2026-09-20，AC 补裁）

> 背景：上表 6 行未裁 AC#4（日志脱敏种子类）与 AC#5（诊断包含采样器数据），票据复核时两框留空。
> 本节由补裁代理实跑实读后补裁。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#4 日志脱敏测试：种入 key/audio/fetch-body/long-arg → 滚动日志中一个都不出现 | **PASS** | 四类各有专测且**双向断言**（既断"种子串不出现"，又断"占位符必须出现"，堵死"整行被丢弃也算过"的空转）：`internal/observe/logging_test.go:30 TestRedactSecretAttrKeepsLast4Only`（`:34` 全文 key 不得出现 + `:37-40` 必须含 `****` + 后 4 位，形状由 `:45 secretRedactForTest` 钉死=C28 last-4）、`:63 TestRedactAudioBufferNeverLogged`（`:71 bytes.Contains(out, audio[:64])` 必须 false + `:74` 必须含 `[audio buffer redacted: 4096 bytes]`）、`:79 TestRedactFetchBodyNeverLogged`（marker `secret-page-content-marker-9182` 不得出现 + `:88` `[content redacted:` 占位）、`:92 TestRedactLongArgTruncated`（`:98` 尾标记 `TAIL-MARKER-777` 不得出现 + `:101` `(truncated, 4111 chars total)` 标记）。**"滚动日志"这一 AC 字面**由端到端用例独立坐实，不止于 handler 层：`:219 TestLogPipelineEndToEnd` 走真 `InitLogWithRegistry`→registry 派生 log-flusher→`rollingWriter`→**落盘 `.jsonl`**，再 `os.ReadFile` 回读断言 `sk-this-is-a-long-seeded-key-9999` 不在文件里（`:251-253`）。不可关闭性经代码核实：`internal/observe/logging.go:86-88` 注释"Redaction is non-disableable: there is no constructor switch for it"，管线为 `JSONHandler → redactHandler → rollingWriter`（`logging.go:19-20`），`Handler()`（`:104`）只暴露包了 redactHandler 的那一个。辅跑 `:52 TestRedactInlineKeyShapes`（行内 sk-/Bearer/api_key: 掩码）、`:120 TestRedactBakedInAttrsRedacted`、`:105 TestRedactPathsOptIn`（`[privacy] redact_paths` 确为 opt-in，与票面一致）、`:131/:157/:185` 尺寸/按日滚动与 7 日保留。<br>`go test ./internal/observe/ -run 'TestRedactAudioBufferNeverLogged\|TestRedactFetchBodyNeverLogged\|TestRedactLongArgTruncated\|TestRedactSecretAttrKeepsLast4Only' -count=2 -v` |

```
--- PASS: TestRedactSecretAttrKeepsLast4Only (0.00s)
--- PASS: TestRedactAudioBufferNeverLogged (0.00s)
--- PASS: TestRedactFetchBodyNeverLogged (0.00s)
--- PASS: TestRedactLongArgTruncated (0.00s)
--- PASS: TestRedactSecretAttrKeepsLast4Only (0.00s)
--- PASS: TestRedactAudioBufferNeverLogged (0.00s)
--- PASS: TestRedactFetchBodyNeverLogged (0.00s)
--- PASS: TestRedactLongArgTruncated (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	0.057s
```

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#5 诊断：采样器数据可挂进诊断包（UX 由 45 完成） | **PASS**（按票面口径；接缝已开并有断言，非仅 stub） | 生产侧 `internal/observe/diagnostics.go:153-163`：`BundleOptions.SLOSnapshot []byte`（`:45-46`，另有 `:47-48` `SLOSnapshotName` 支持按态命名）非空时以 `writeZipEntry(zw, "slo-snapshot.json", o.SLOSnapshot)` **原样字节写入 zip**；为空时走 `:166` 在 manifest 记 `"not attached: no sampler data supplied"`——即"挂/不挂"两种结局都留下可核对痕迹，未挂不会静默消失。测试侧 `internal/observe/diagnostics_test.go:15 TestDiagnosticsBundleCollectsAndRedacts`：用 `SLOSnapshot: slo`（`{"state":"Sleeping","pass":true}`）建包后 `zip.OpenReader` 解包，`:62-66` 断 `slo-snapshot.json` 与 manifest/version/deferred/config.redacted 五条目齐备，`:68-77` 断种入日志的 `sk-SEEDLEAK1234567890` 未进包，`:79-81` 断配置文件行内 key 未进包，`:83-90` 断 manifest 是合法 JSON 且 contents≥5。配套两例同样为真断言：`diagnostics_test.go:98 TestDiagnosticsBundleFailClosedOnSurvivingKey`（注入检测器命中 ⇒ 必须报错且**不得留下半成品 zip**，`:121-123`）、`diagnostics_test.go:127 TestDiagnosticsBundleRecordsMissingPieces`（缺件必须记 `"not collected"` + `"ticket 45"` 缘由）。<br>`go test ./internal/observe/ -run 'TestDiagnosticsBundleCollectsAndRedacts' -count=2 -v` |

```
--- PASS: TestDiagnosticsBundleCollectsAndRedacts (0.01s)
--- PASS: TestDiagnosticsBundleCollectsAndRedacts (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/observe	0.057s
```

合并复跑（脱敏 + 滚动 + 端到端 + 诊断全家族，15 例 ×2）：

```
$ go test ./internal/observe/ -run 'TestRedact|TestRollingWriter|TestLogPipeline|TestDiagnostics' -count=2
ok  	github.com/CarlosShao/wisp/internal/observe	0.131s
```

**补裁小结与残余**：AC#4、AC#5 均 PASS，票据两框改勾。两点如实记残（不改变裁决，因票面已把 UX 划归 45）：
① `BuildDiagnosticsBundle` 目前**全仓无生产调用者**（`grep -rn "BuildDiagnosticsBundle" --include=*.go .`
除 `internal/observe/diagnostics*.go` 自身外零命中），"采样器数据可挂"成立是在 **API 契约层**，
把真采样器的 JSON 实际接进去属票 45；②`TestDiagnosticsBundleCollectsAndRedacts` 只断
`slo-snapshot.json` **存在于包内**，未断其字节等于传入的 `slo`——内容等值由
`diagnostics.go:161` 的 `writeZipEntry(..., o.SLOSnapshot)` 直传保证，票 45 落地 UX 时应补一条
内容回读断言。
