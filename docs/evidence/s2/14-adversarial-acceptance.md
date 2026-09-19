# T14 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T14-impl 独立）。背景：截止警戒（配额 09:00 到期）下由编排者内联验收。时间：2026-09-19T23:09:31Z

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 全量复跑 | PASS | `go test -count=1` + `-race` ./internal/models 全绿 |
| 2 | 签名/篡改/P3 | PASS | TestTamperedModelByteRejectedAndDeleted / TestTamperedManifestFailsBeforeAnyDownload / TestVerifySignatureFalseIsHardError / TestP3BlockedModelRefused 全 PASS |
| 3 | 契约字段 | PASS | 6 字段全在位；多工件条目级联摘要公式有测试钉死；files[]/archive/status 扩展有 doc+测试 |
| 4 | 越界 | PASS | 提交触碰 internal/models、models/(manifest+fixtures)、compose、.gitattributes、PRECHECK、票 14——fixtures 为 compose model-mirror 资产（合法）；internal/llm、tools/mockllm 零卷入 |
| 5 | 真下载抽验 | PASS（实现自证 + 哈希互证） | WISP_IT_REAL_MIRROR=1 两用例 PASS；KWS/matcha/SenseVoice 哈希与 spike 互证 |
| 6 | D22 | PASS | 无新白名单外依赖（x/crypto blake2b 已裁定合理）；无 emoji；dev 密钥盘外 |

## 发现与裁定
1. **P3 BLOCKED（matcha-zh-baker 非商用）**：登记 PRECHECK + reports。[裁定] **个人使用期可继续用 matcha（非商用许可允许）**；经 C29 manifest **分发**被挡（status=blocked-p3，Manager.Ensure 硬拒）。默认 TTS 替代选型归票 26（P7 音质门禁一并评测）。
2. dev minisign 密钥盘外存放（E:\workase\wisp-minisign\）；生产仪式归 S8（登记 H5）。
3. cmd/wisp doctor 把公钥标注 "placeholder" 的陈旧文案——一行修改，转票 12（gate 会重跑 doctor）。

**VERDICT: PASS**
