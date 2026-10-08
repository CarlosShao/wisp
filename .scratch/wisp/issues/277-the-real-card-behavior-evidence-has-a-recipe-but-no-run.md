# 票 277 — 「常驻进程会不会真举出一张审批卡」的行为凭据：配方齐了，那一发没有跑

**立票**：2026-10-08 19:3x 编排者（来路＝台账 `A725` 举卡链量清／`A732` 裁定／`A735` 配方收档）
**性质**：`AC#7`（票 246）的勾**只覆盖到「起一条任务管线」**（台账 `A732` 已裁）；本票承接**另一半**——**真举出一张卡**的**行为凭据**。⚠ 本条**不是契约问题**、⛔ 不需要 owner 拍板。

## 现量（引用前先重跑）

- 链的静态形状全接：`cmd/wisp/main.go:66`→`resident_windows.go:260` 起任务源→装配→`run.go` 注门/交 bridge→`internal/tools/bridge.go:424/:443/:458`（有效级 L1 无会话授权 或 L2 才问）→`internal/agent/approval/gate.go:294/:536` 两处 `ui.Prompt`→常驻侧 `ballCardUI.Prompt`（`resident_approval_windows.go:864` 一带）。
- 卡出现的可观测信号真身：`resident_approval_windows.go:885` `fmt.Printf("wisp: 卡片挂起：%s %s（编号 %s）
", p.Level, p.Tool, p.CorrelationID)`＋`:880` 审计行 `approval: 常驻进程显示一张确认卡片`。
- 配方（可照做）：`.scratch/wisp/probes/card-proof-prep-1/01-recipe.md`（A 道＝winlive 台件注入文本，今天跑得动；B 道＝真控制台 `task` 动词，今天跑不动：`%APPDATA%\wisp` 无凭据＋`build/wisp.exe` 旧于 HEAD）。

## 要建什么

- [ ] **AC#1 真机跑 A 道**：照配方 A 道跑一发，逐字收"卡行"（`:885` 那句）与审计行；**判红绿按配方 §4**（绿＝卡行＋`esc_borrowed=true`/`orb_state=Confirming` 同 corr；红＝入口/管线被拒或超时无卡；判不了＝Esc 被占/无桌面）。⚠ 本票**不要求**面板在场（该腿无面板，观测面＝stdout＋审计＋winlive `t.Logf`）。
- [ ] **AC#2 B 道的前置记账**：真控制台那一支要 owner 的凭据＋重出 exe ⇒ 在条件齐之前，⛔ 不派；本格只要求把"缺谁/缺什么"逐字登记（配方 §5 已写，复跑核）。
- [ ] **AC#3 与票 246 的边界**：⛔ 本票**不动**票 246 任何框；结论回写进本票。

## 禁区

⛔ 不许把"静态链全接"当行为凭据；⛔ 不许为让 A 道变绿放宽任何断言；⛔ 真窗类发帖按 `A718 §3` 的凭据强度限制（老 exe/旧 dist 不可作数）；⛔ 三枚冻结件一字不动；⛔ 零 push。
