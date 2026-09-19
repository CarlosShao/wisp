# 报告 — T14 dev minisign 密钥登记（C29 / D33-F3 / D41e）

- 日期：2026-09-19
- 登记：docs/reports/pending-and-issues.md [H5]
- 密钥对：minisign 格式（Ed25519），keyid `a3c8794f3fd94fc5`
- 私钥：`E:\work\base\wisp-minisign\wisp-models.key`（**盘外，绝不入仓**；0600；
  自有明文信封格式 `Ed||keyid||seed||pub`，仅 dev——真 minisign 私钥是加密格式，
  S8 生产仪式换用真 minisign/硬件保管并**必须轮换**）
- 公钥：`E:\work\base\wisp-minisign\wisp-models.pub`，且硬编码进
  `internal/buildinfo.MinisignPublicKey`（运行时唯一验证锚点）
- 签名工具：`tools/signmodels`（keygen/sign/verify，stdlib-only）；包装脚本
  `scripts/sign-models.ps1`（私钥路径只经 `-KeyPath` / `WISP_MINISIGN_KEY` 环境变量注入）
- 已签工件：`models/manifest.json` + `models/manifest.json.minisig`（同时提交，
  改其一必须重签；CI 断言 `TestRealManifestInRepoVerifies`）

## 必须遵守

1. 私钥文件、`WISP_MINISIGN_KEY` 值永不入仓/入日志/入诊断包。
2. 生产发布（S8）前必须完成密钥仪式并轮换 buildinfo 公钥——当前 dev 密钥
   仅证明机制，不构成发布签名。
3. 更换密钥 = keygen 重生成 → buildinfo 公钥更新 → manifest 全量重签 → 双提交。

## 兼容性

验证器（`internal/models/minisign.go`）为格式级手写实现：接受 minisign 两种
算法（`Ed` 直签 + `ED` Blake2b 预哈希，后者即 minisign 默认），校验 keyid 绑定
与 trusted-comment 二段签名，与真 minisign CLI 的签名文件互通；本仓签名器产出
`Ed` 直签（真 minisign 默认验证模式可验）。
