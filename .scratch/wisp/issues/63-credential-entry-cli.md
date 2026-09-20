# 63 — 凭据录入入口：`wisp secret set/get/list/unset`（隐藏输入 → DPAPI）

**Status:** review
**Claimed by:** agent-ticket63-fix2
**Last update:** 2026-09-20T09:40Z
**Blocked by:** 06-secretstore-envs（DPAPI Store 已存在）
**Parallel slots:** ≤1 sub-agent
**Spec refs:** SPEC-03 §3.1（secret refs）, SPEC-06（凭据不落明文）, D22, P13 便携模式, R7

## What to build
owner 2026-09-20 批准（原话：批）。补一个**最小**凭据录入入口，让人不把 key 贴进任何对话框就能把凭据
存进现有 DPAPI 存储：

```
wisp secret set   <name>      # 交互式隐藏输入两次确认；也支持 `--from-stdin`
wisp secret get   <name>      # 默认只回显掩码（RedactSecret），--show 才打印明文
wisp secret list              # 只列 blobID 与创建时间，绝不列内容
wisp secret unset <name>      # 删除 blob，并检查是否仍被 config 引用（有引用则拒绝，--force 覆盖并记审计行）
```

录入后写 `secret:<blobID>` 引用，供 config 的 `[providers.*].api_key_ref` 一类字段使用
（票 05/06 的机制已存在，本票**不新增存储实现**，只加 CLI 外壳）。

## Key constraints
- **绝不回显**：输入走 `golang.org/x/term` 的 `ReadPassword`（Windows 控制台已验证可用）；
  明文**不进 argv**（否则 `Get-CimInstance Win32_Process` 能看到命令行——这是本机可查的泄露面），
  `--from-stdin` 也只读一次、不落地中间文件。
- **绝不进日志**：任何 `slog`/error 字符串里不得出现明文；错误信息只带 blobID。
  复用 `internal/secret.RedactSecret`。
- **DPAPI 用户域**：沿用票 06 的 `protect/unprotect`；便携模式（P13）走既有分支，
  非便携时不得静默降级成明文文件。
- **WISP_ENV 隔离**：dev/test/prod 三套目录与互斥名沿用票 06，`secret set` 必须明确作用在
  当前 `WISP_ENV` 上，并在输出里回显环境名（防止把 dev 的 key 存进 prod 目录还以为成功）。
- 不做 GUI（票 39/40 才是配置界面）；不做云同步；不改 `internal/secret` 的存储格式。
- 禁止裸 `go func(`；禁止 emoji。

## Out of scope
配置文件的 GUI 编辑器（票 39）；密钥轮换仪式（票 56/S8）；把 key 写进 config 明文（永久禁止）。

## Acceptance criteria
- [x] `wisp secret set` 隐藏输入两次确认，落盘为 DPAPI blob，`wisp secret get` 默认掩码显示。
- [x] 明文不出现在：argv、任何日志行、任何错误信息、临时文件（用测试证明，不是靠声明）。
- [x] `list` 只输出 blobID+时间；`unset` 在被 config 引用时拒绝并指出引用位置，`--force` 时留审计行。
- [x] `WISP_ENV` 三档各自独立：同一 name 在 dev/test/prod 是三个不同 blob（测试覆盖）。
- [x] 便携模式与非便携模式两条路径都有测试（沿用票 06 的 seam 注入）。
- [ ] 端到端：`wisp secret set` 存好后，config 里 `api_key_ref = "secret:<id>"` 能被 provider 解析出发请求（用 mockllm，测试里用假 key）。
      —— 未勾原因：`TestSecretEndToEndConfigRefResolvesAtRequestTime` 证明了 ref→`config.LoadFile`→
      `ProviderKeys`→`Authorization` 头真的上线，但发请求的是测试自己的 http client，不是
      `internal/llm` 的 provider；provider 那一段是调用方自证的（README 规则：安全判定不得由调用方选择器决定）。
- [ ] 对抗验收由非实现者执行，报告含与本表 **1:1 的裁决表**（README 规则 6）。
      —— 未勾原因：本票由实现者填写，验收报告尚未产出；这是复核者的格子，不是能自证的格子。

## Progress log (append-only, newest last)
- [2026-09-20T06:50Z] agent=agent-ticket63-secretcli did=claimed next=set/get/list/unset+leak-tests
- [2026-09-20T07:06Z] agent=agent-ticket63-secretcli did=加 `wisp secret set/get/list/unset`（cmd/wisp/secret.go，热插拔 secretIO+readHidden seam，仅 bool flag）；internal/secret 加命令外壳所需的 Delete/Exists/Blobs/ConfigRefs（存格式零改动，blobPath 仍是唯一路径咽喉）+ctime_windows/other；gofmt/vet/secret count=2/secret -race 全绿 next=cmd/wisp 泄漏测试（argv/日志）
- [2026-09-20T07:41Z] agent=agent-ticket63-secretcli did=cmd/wisp/secret_test.go 13 个测试全绿（掩码/明文 --show、名称校验、确认不一致、无控制台、bool-only flag 源码级 pin、8 种 `--value=<key>` 形态拒绝且不回显、5 条失败路径零明文入日志、unset 引用拒绝+点名字段+--force 审计行、dev/test/prod 三 blob、便携/非便携、端到端 config→request）；修两处真 bug：Go flag 在首个位置参数后停止解析（`set <name> --from-stdin` 原本失败）、Store 写失败错误未点名 ref next=OS 级 argv 取证（PEB CommandLine）
- [2026-09-20T08:50Z] agent=agent-ticket63-fix2 did=resumed after model-connection kill (predecessor committed 722f819, 13 tests green); auditing uncommitted argv-forensics test next=finish OS-level proof then reconcile 7 boxes
- [2026-09-20T09:20Z] agent=agent-ticket63-fix2 did=完成 OS 级 argv 取证。状态报告：前任留下的 `cmd/wisp/secret_argv_windows_test.go` 未跟踪但**可编译、原样即 PASS**（不是半成品），它缺的是"探针本身可信"这一半，不是"测试跑通"那一半。断言对象是**活着的** wisp.exe 子进程的 PEB `_RTL_USER_PROCESS_PARAMETERS.CommandLine`（amd64 偏移 0x20/0x70 + UNICODE_STRING，非 amd64/非 windows 一律 Fatal 不 skip）；此刻明文已经过 stdin 进入子进程并且子进程仍阻塞在读上，所以是在"泄漏会看得见"的时刻取证：OS 可见命令行必须逐 token 等于 `secret set argvprobe --from-stdin`（任何被追加的 token 都 FAIL），整串与中段片段都不含明文；`--value=<key>` 走真二进制=退出码 2、零 blob、不回显。新增正对照 `TestProcessCommandLineProbeDetectsAPlantedValue`：把假 key 当**参数**种进一个活子进程，要求同一个探针把它读出来——探针若是瞎的，上面所有"argv 没有明文"就都是空断言，现在这种情形会 FAIL。修 1 处谎报的失败信息（原 Fatal 声称"最后一次错误由 procCommandLine 报告"，实际并未报告；现按 pollCommandLine 返回真实 last error）。登记遗留：探针只覆盖被观察的直接子进程，Windows 无 execve，未来若 `secret set` 改为再生子进程并在其 argv 里带值，孙进程不在本票视野内；补法=活窗口内 Toolhelp32 枚举子进程并逐个跑同一 PEB 读，判据=任何读到的命令行都不含明文。next=变异检验判据仪器 + 跑门禁
- [2026-09-20T09:40Z] agent=agent-ticket63-fix2 did=变异检验（本票唯一能证明判据仪器可信的手段）+ 门禁全绿 + 勾栏 5/7。①在 runSet 里把明文塞进 slog 属性 → 4 处断言 FAIL；把 secret_test.go 退回**已提交版**（去掉前任未提交的 activeCapture 共享缓冲）后同一变异 **0 FAIL = 假绿**，证明那份未提交改动是承重的，已随 4477f56 一并提交。②在 collectSecret 里把 stdin 明文写成 %TEMP% 一个 scratch 文件 → 新测试 `TestSecretFromStdinWritesNoIntermediateFile` FAIL 并点名该文件（它自带"先种含明文诱饵、要求扫描器找到"的正对照，否则"没有文件含明文"同样空断言）；AC2 的临时文件半边此前无人证明，现在有了。③共享缓冲暴露一个真回归：`unset` 的"forced=true 只属于强制删除"是整缓冲 grep，被同测试前一个子测试的合法 forced=true WARN 污染 → 用 `auditRecordFor` 把断言收窄到本条审计记录，并**加强**为必须出现 `forced=false`（原断言只验"别处没有 forced=true"）；未弱化任何既有断言。④门禁：`gofmt -l` 三个文件为空；`go vet ./cmd/wisp/ ./internal/secret/` RC=0；`go test ./internal/secret/ -count=2` ok 0.628s；`go test ./cmd/wisp/ -count=2` ok 12.039s；`go test ./cmd/wisp/ ./internal/secret/ -race -count=1` ok 7.353s / 1.180s；cmd/wisp 现 18 顶层测试 + 31 子测试。注意：cmd/wisp 的测试二进制链接 sherpa-onnx，跑前必须把 `third_party/sherpa-onnx` 放进 PATH，否则子进程 0xc0000135 起不来（包既有属性，非本票引入，BUILD.md 未记）。⑤AC 勾选：1/2/3/4/5 勾；6 未勾=端到端里发请求的是测试自己的 http client 而非 `internal/llm` 的 provider，provider 段属调用方自证（且现在接 internal/llm 会把本票门禁绑到另一代理在飞的编辑上）；7 未勾=对抗验收要由非实现者产出 1:1 裁决表，不是实现者可自证的格子。next=非实现者按本表 1:1 裁决


