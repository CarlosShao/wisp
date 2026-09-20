# 63 — 凭据录入入口：`wisp secret set/get/list/unset`（隐藏输入 → DPAPI）

**Status:** in-progress
**Claimed by:** agent-ticket63-secretcli
**Last update:** 2026-09-20T06:50Z
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
- [ ] `wisp secret set` 隐藏输入两次确认，落盘为 DPAPI blob，`wisp secret get` 默认掩码显示。
- [ ] 明文不出现在：argv、任何日志行、任何错误信息、临时文件（用测试证明，不是靠声明）。
- [ ] `list` 只输出 blobID+时间；`unset` 在被 config 引用时拒绝并指出引用位置，`--force` 时留审计行。
- [ ] `WISP_ENV` 三档各自独立：同一 name 在 dev/test/prod 是三个不同 blob（测试覆盖）。
- [ ] 便携模式与非便携模式两条路径都有测试（沿用票 06 的 seam 注入）。
- [ ] 端到端：`wisp secret set` 存好后，config 里 `api_key_ref = "secret:<id>"` 能被 provider 解析出发请求（用 mockllm，测试里用假 key）。
- [ ] 对抗验收由非实现者执行，报告含与本表 **1:1 的裁决表**（README 规则 6）。

## Progress log (append-only, newest last)
- [2026-09-20T06:50Z] agent=agent-ticket63-secretcli did=claimed next=set/get/list/unset+leak-tests
- [2026-09-20T07:06Z] agent=agent-ticket63-secretcli did=加 `wisp secret set/get/list/unset`（cmd/wisp/secret.go，热插拔 secretIO+readHidden seam，仅 bool flag）；internal/secret 加命令外壳所需的 Delete/Exists/Blobs/ConfigRefs（存格式零改动，blobPath 仍是唯一路径咽喉）+ctime_windows/other；gofmt/vet/secret count=2/secret -race 全绿 next=cmd/wisp 泄漏测试（argv/日志）
- [2026-09-20T07:41Z] agent=agent-ticket63-secretcli did=cmd/wisp/secret_test.go 13 个测试全绿（掩码/明文 --show、名称校验、确认不一致、无控制台、bool-only flag 源码级 pin、8 种 `--value=<key>` 形态拒绝且不回显、5 条失败路径零明文入日志、unset 引用拒绝+点名字段+--force 审计行、dev/test/prod 三 blob、便携/非便携、端到端 config→request）；修两处真 bug：Go flag 在首个位置参数后停止解析（`set <name> --from-stdin` 原本失败）、Store 写失败错误未点名 ref next=OS 级 argv 取证（PEB CommandLine）


