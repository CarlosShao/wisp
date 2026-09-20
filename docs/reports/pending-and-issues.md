# Reports — 问题项 / 阻塞项 / 待人项登记

> 持续更新。条目解决后移入文末"已解决"。格式：[级别] 标题 — 状态 — 阻塞什么 — 下一步。

## 待人工审核（pending-human-review）

- [H1] **悬浮球 20 态视觉签收** — 等用户 — 不阻塞（票 07 已按截图证据收口）— 用户浏览
  `docs/evidence/s1/ball-states/*.png` 与 `design/screens/ball.html` 对照，满意则关闭；
  不满意提修改意见转新工单。
- [H2] **LLM API Key** — 等用户提供 — 阻塞票 09 的黄金 SSE 录制（录制需要真 provider 响应；
  mock 与回放框架不阻塞，已先行）与票 12 的真链路验收 — 用户提供后设 `WISP_LLM_KEY` 环境变量。
- [H3] **P10 命名残余核查**（npm/PyPI/域名/商标）— 等用户/外部查询 — 阻塞票 56 — 非 S1 阻塞项。
- [H4] **`web.search` 实现路径**（抓结果页 vs 搜索 API）— 等用户拍板（PLAN §15#3）— 票 22 以
  SearchProvider 接口 + 抓取默认实现先行，切换不返工。

## 平台事件记录（platform-incidents）

- [P1] **子代理验证码/配额故障（2026-09-19 持续）**：当日 9 次子代理派发被
  "Captcha instance timed out"/"exceed quota limit" 打断（时长 38s–2.9h 不等）。
  应对：票据进度日志 + 断点续传协议（零工作丢失实证：票 05 五棒接力、票 07 三棒接力均无返工）；
  验收类工作由编排者亲自执行（独立性满足）；并发按约定 3→2 回落。
- [P2] **GitHub 直连 TLS 间歇失败**：push 常规重试 ≤5 次可消化；未造成丢失。
- [P3] **`git add -A` 两次吞并行 WIP**（T03/T02 期间）：已按 `git rm --cached` 先例修复；
  此后全仓强制显式路径提交，未再发生。

## 裁定记录（orchestrator rulings）

- **[R1] matcha-zh-baker P3（非商用）**（2026-09-20）：个人使用期可继续用（许可允许）；C29 manifest 分发被挡（blocked-p3 硬拒已实现）。**默认 TTS 替代选型归票 26**（与 P7 音质门禁一并评测；候选： sherpa vits-zh 许可干净系列 / 云 TTS）。
- **[R2] 截止计划**（2026-09-20 07:2x，配额 09:00 到期）：在途 14✅/08（deadline-pause 规则已发）；**票 10（2–4h）与 12 留给新会话**——交接文档 docs/reports/HANDOVER.md 已含新会话行动清单。

## 完成度审计（2026-09-20 08:58，audit-A + audit-B）

- **结论：13 张 done 票无 BLOCKER；实质工作全部有证据支撑**（audit-B 深查 4 张的引用抽验全属实）。
- **[A1] Audit-A MAJOR→已修**：票 02 归档漏改 Status 字段（已修，commit 08f9ff8）。CGO=0 全仓限制为已登记基线（HANDOVER §6），非违规。
- **[B1] Audit-B MAJOR → 转新会话第一组簿记任务**：13 张 done 票的 AC checkbox 未逐条勾选（71 条 `[ ]`）——**实质裁决都在 docs/evidence/ 各验收报告里**，勾选动作纯簿记；票 18 的"perf ≤1ms"AC 为实现自认 DEFERRED（bench 未写），已在票内登记，新会话补 bench 或正式降级。新会话开工前先花 15 分钟做此勾选 + 逐条裁决行。
- **[B2] MINOR×6 全部接受为簿记**：日志时序（02/17/18）、报告路径笔误（17/18）、票 04 双时间戳、票 18 缺 handoff 行、H5/H6 编号、票 07 报告引用文件名笔误（内容在 pending-and-issues.md 未丢失）——新会话顺手清理，不阻塞任何事。
- **[B3] H2/H5 登记补强**：票 16 真机三场景（原 H6 漏登 registry）已由本条补登。

## 阻塞项（blocked）

（当前无硬阻塞。票 09 的黄金录制部分依赖 H2；其余在途/排队票均可推进。）
- [H5] **TTS 模型选型（P3 BLOCKED：matcha-zh-baker 非商用）** — 等用户/编排者拍板 —
  阻塞 S2 的 TTS 环节上船（管线已就绪，条目 `tts-matcha-zh-baker` 已标
  `blocked-p3` 并被下载器拒绝）— 需选定数据许可可商用的中文 TTS（onnx 可转 +
  sherpa-onnx 支持），详见 PRECHECK.md P3 与 docs/reports/2026-09-19-t14-dev-minisign-key.md
  同目录的密钥登记（dev 密钥轮换归 S8）。
