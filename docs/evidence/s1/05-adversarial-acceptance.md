# T05 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者为 T05-impl/T05-resume3/T05-resume4 三棒接力；编排者独立）。
> 背景：T05-adv 子代理被平台验证码打断（22 分钟），按"不阻塞"指令由编排者亲自验收。
> 时间：2026-09-19T13:56:25Z

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 全量复跑 | PASS | `go test -count=1 ./internal/config/...` ok（47 顶层/115 RUN） |
| 2 | 行号实证 | PASS | TestLoadFileUnknownKeyErrorNamesLine / LineNumbersAreOneBased 独立可见 PASS |
| 3 | 存储分界 | PASS | schema.go 中 api_key/health/probe/latency 命中均为顶部**边界说明注释**，无真实字段 |
| 4 | 越界 | PASS | 实现提交仅 internal/config + go.mod/go.sum + 票 05 |
| 5 | catalog/迁移/三档 | PASS（信任实现测试 + 抽查文件） | catalog.go 链校验、migrate.go dry-check/备份、manager.go 方向引擎代码审读 |

注：受平台故障限制，本票验收为"复跑 + 代码审读"级，未逐项独立重构实验（对照 T04 深度）。
风险低（纯逻辑层、47 测试覆盖密）。**VERDICT: PASS**

## Addendum 裁决（2026-09-20，AC 补裁）

> 背景：原报告 5 行审计表未覆盖 AC#5/AC#6，票据复核时两框留空。本节由补裁代理实跑实读后补裁。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#5 硬编码只读字段写入报错 | **PASS** | 生产侧 `internal/config/validate.go:68`（`audio.half_duplex`）、`:78`/`:82`（`privacy.keep_transcript`/`keep_audio`）、`:92`（`models.verify_signature`）四处各自返回字面量错误串，例 `validate.go:70`：`"config.toml: audio.half_duplex is hard-coded true (read-only); writing false is rejected"`。`validate` 确在加载管线上（`internal/config/loader.go:62`，`LoadFile` 内），故"写入文件"路径亦被拦。测试侧 `internal/config/validate_test.go:129 TestValidateHardcodedReadOnly` 对四个字段逐一断言 ①必须报错 ②错误串必须点名 key path（`strings.Contains(err.Error(), tc.path)`），非空断言。<br>`go test ./internal/config/ -run 'TestValidateHardcodedReadOnly' -count=2 -v` |

```
--- PASS: TestValidateHardcodedReadOnly (0.00s)
    --- PASS: TestValidateHardcodedReadOnly/audio.half_duplex=false (0.00s)
    --- PASS: TestValidateHardcodedReadOnly/privacy.keep_transcript=true (0.00s)
    --- PASS: TestValidateHardcodedReadOnly/privacy.keep_audio=true (0.00s)
    --- PASS: TestValidateHardcodedReadOnly/models.verify_signature=false (0.00s)
=== RUN   TestValidateHardcodedReadOnly
--- PASS: TestValidateHardcodedReadOnly (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/config	0.030s
```

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#6 往返 load→marshal→load 结构全等 | **PASS** | `internal/config/boundary_test.go:188 TestRoundTripLoadMarshalLoad`：`SaveFile → LoadFile → SaveFile → LoadFile`，**两次** `reflect.DeepEqual`（`:204` 第二轮 vs 第一轮；`:210` 落盘前原结构 vs 首次载入结构——即" marshal 不增不减字段"这一 GUI 双真相源_guard_ 的正断言）。夹具 `boundary_test.go:74 fullConfig` 对 SPEC-03 §3 全部 17 段（app/ball/hotkey/session/voice/audio/llm/agent/risk/fs/net/privacy/memory/panel/cost/plugins/models/observe）均写入**非默认值**，含 plugins 裸映射、compat、price、quota、chains——不是"全默认值恒等"的空转。辅以 `boundary_test.go:216 TestRoundTripBytes`（`MarshalCanonical→decodeStrict` 恒等，不经文件）。<br>`go test ./internal/config/ -run 'TestRoundTripLoadMarshalLoad|TestRoundTripBytes' -count=2 -v` |

```
=== RUN   TestRoundTripLoadMarshalLoad
--- PASS: TestRoundTripLoadMarshalLoad (0.01s)
=== RUN   TestRoundTripBytes
--- PASS: TestRoundTripBytes (0.00s)
=== RUN   TestRoundTripLoadMarshalLoad
--- PASS: TestRoundTripLoadMarshalLoad (0.01s)
=== RUN   TestRoundTripBytes
--- PASS: TestRoundTripBytes (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/config	0.049s
```

**补裁小结**：AC#5、AC#6 均 PASS，票据对应两框改勾。SPEC-03 §3 对 `keep_transcript`/`keep_audio` 的
原文（"硬编码 false，只读显示，写入 true 报错"）与测试断言方向一致。
