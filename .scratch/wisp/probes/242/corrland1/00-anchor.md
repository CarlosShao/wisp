# 242-corrland-1 · 00 起手锚

- `date`：`Thu Oct  8 19:31:22 CST 2026`（起手一发）
- `git log -1 --format='%H %ad %s'`：`bbbc18fa23c663a79634cf67f18cc44e64b2711a Thu Oct 8 19:31:01 2026 +0800 A745 落账：收 242-corrcensus-1（9ae812d6/0e0ac7b9 四枚件）correlation 语义与读者普查。…在飞 0 枚；零翻框零 push`
- `git status --porcelain | head -20`（只登记，⛔ 与本腿无关的改动一律不碰）：
  ` M .gitignore`
  ` M .scratch/wisp/probes/152/my152.py`
  ` M .scratch/wisp/probes/161/r6/logs/flip-*.txt`（9 枚）
  ` M .scratch/wisp/probes/242/r3/logs/probe-routed.txt`
  ` M .scratch/wisp/probes/268/v1/evidence.md`
  ` D design/assets/base.css`、` D design/assets/icons.js`、` D design/assets/theme.js`、` D design/assets/tokens.css`
  ` M design/doubao/README.md`、` M design/doubao/demo/app.js`、` M design/doubao/demo/index.html`、` M design/doubao/demo/styles.css`
  （以上 design/** 一栏＝票 242 票面 §"排程"所记"别人在工作树里删了 design/assets/**"，非本腿地界。）

## 本腿任务（按 A745 裁定与派单逐字）

- 产码：`internal/agent/loop.go:654` 那一跳从 `CorrelationID: taskID` 改为**每次工具调用铸一枚独立 corr**。
- 同批重判：`loop.go:363` 注释（"引用与原文不符"候选）＋`internal/tools/loop_approval_test.go:213-214` 唯一等值断言（具名重判，⛔ 不放宽强度）。
- 新用例：同一任务并发两问 ⇒ 两枚 corr 不同且各自可路由（`internal/agent` 包内）。
- 三枚测试构造点（pointer_183:127／pointer_185:120／ticket175r2:95）逐枚重判。
- 写面：产码（loop.go 那一带）＋相关测试件＋本目录 `.md`。⛔ 不 push、零翻框、不改台账。
