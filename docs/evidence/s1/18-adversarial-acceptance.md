# T18 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T18-impl 独立）。时间：2026-09-20T00:27:30Z。

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | C26 管线全序 | PASS | pathresolver.go：env/~ 展开 → 绝对化+Clean（唯一豁免点）→ GetFinalPathNameByHandle（天然展开 8.3）→ REPARSE_POINT 逐组件检测默认拒 → UNC 规范化（3 种拼写） |
| 2 | 红队四连（真 OS 产物） | PASS | 真 mklink /J junction 拒；真 GetShortPathNameW 8.3 拒（卷不支持 skip 注明合法）；UNC×3 拼写拒；`\?\` 前缀拒；evil-twin 兄弟 junction 拒；大小写/./.. 混写拒 |
| 3 | A/B 黑名单 | PASS | A 档 9 类锚点全 deny 且 Gate 无 override；B 档 6 规则默认 L2 + 单文件豁免出审计日志 |
| 4 | reparse 豁免 | PASS | 按具体路径、大小写不敏感、无前缀泄漏 |
| 5 | D22 | PASS | tools/d22scan 全仓 clean（allowlist 1 行豁免经裁定：T14/票 18 迁移用例）；真 junction 测试（非 mock 字符串） |
| 6 | 测试 | PASS | 10 条全绿（全包合流后）+ race 绿 |

**剩余**：暖缓存 bench（现每次新开句柄，量级达标）、~user 展开、macOS realpath（DEFERRED）——均登记票据。
**VERDICT: PASS**
