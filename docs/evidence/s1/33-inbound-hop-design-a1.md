# `33-a1` — 面板入向那一跳（网页点一下 → Go 收到）：地基普查与代价表

- 程：`33-a1`（只读设计核·零产码）｜工单：`.scratch/wisp/issues/33-panel-host-c27.md`
- 派单：`.scratch/wisp/dispatches/2026-09-28-150x-readonly-33-a1-inbound-hop-design-core.md`
- 起手时刻：`2026-09-28 14:55 +0800`｜起手 HEAD：`ee0ef8ec`（派单锚 `1f22f1aa`，**已漂**，见 §1）
- 性质：**只读** ⇒ AC 框一枚未勾、`internal/**`／`cmd/**`／`go.mod` 零字节不动。

---

## 1. 起手锚＋写面闸门

（待填：起手 `git status --porcelain -- internal/ cmd/ go.mod` 读数、锚点偏差登记）

---

## 2. Q1 规格原文：这一跳被定成谁的责任

（待填：`SPEC-08`／`PLAN.md` D29／C27／C17 逐字原句 + 现跑行号；`C24`／`C17` 白名单今天定到哪一步）

---

## 3. Q2 依赖面与可复用地基

（待填：`grep -c webview go.mod` 现量；`internal/ball/sta_windows.go`／`internal/proc/**`／`internal/winsec/**` 逐枚 `grep -n "^func "` 现量；新写几枚／复用哪几枚）

---

## 4. Q3 最小闭环逐环表

（待填：从 `postMessage` 到 `HandleModeRequest` 被真调用，逐环 = 环名 / 现在有什么 / 缺什么 / 最近邻指到现有行）

---

## 5. Q4 冷拉起代价

（待填：`AGENTS §2` 那条待定项今天能给什么现量）

---

## 6. Q5 与 owner 那四枚按钮的依赖关系表

（待填：票 186／187／92+R20／114 逐枚）

---

## 7. Q6 能不能不引 WebView2 先建最小入向

（待填：结论 + 落点 + 判据形状，或具名"必须真宿主"）

---

## 8. 推荐拆法与代价（不拍板）

（待填）

---

## 9. 本程没测什么（逐名）

（待填）

---

## 10. 门禁终态

（待填：`d22scan.sh` rc 与 `ban #8 internal/` examined、`gate-clauses.sh` 红腿名册、`go test ./internal/panel/ ./internal/config/`；2 枚已知红照实记）

---

## 11. 被拒调用＋零删除自证＋工具调用终值

（待填）

---

## 12. next：派写腿之前还缺什么、哪几枚要人先批准

（待填）
