# 318 — junction 卡面的标记串被 `dirOf` 切成碎片：“半截目录”真写进 `[fs] allowed_dirs`（`201-a1` 活体复现、逐字）

- **Status**：**ready-for-agent**（编排者 2026-10-14 2026-10-14 16:38 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）。
- 来路＝`20-r2` 具名上报＋探针腿 `201-a1` **活体复现**（`64bf9b29`；件 `probes/201/a1/00-verdict.md`＋`logs/run3-full.txt`）。
- 关联：票 201（`WideningRule`）· 票 21（L2 通道）· `session/grants.go:264`（既有“拒标记”口径先例）· `bridge.go:1273`（另一处特判）· `tools.Decision.Paths` 语义面。

## 现量（`201-a1` 真跑：真 NTFS junction → 真 bridge → 真 Gate/Queue/Replies → 真 WideningRule → 真 AddAllowedDir）

1. junction 拒卡面把标记串送进 `ReplyCard.Paths`：`"…/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions)"`。
2. `WideningRule`（`replies.go:503`）经 `dirOf`（`batch.go:91` `LastIndexAny("/\\")`）在**标记里的 `/`** 处切开 ⇒ `dir` 尾部＝半截 `(junction`、`ok=true`；两种分隔符拼法同碎。
3. **写口实量**：`AddAllowedDir` 两道闸（空／含 `..`）都过 ⇒ config.toml 真写成 `allowed_dirs = ['…victim.txt (无法规范化: …(junction']`；第二张 L2 卡 `Paths[0]`＝同一碎片逐字；对照（无标记越界读）推出真父目录。
4. 唯一不认标记的消费方＝放宽链；`bridge.go:1273` 与会话授权都特判。

**最坏后果**：连答 `always`＋第二张卡 allow ⇒ `allowed_dirs` 多一条非目录持久垃圾行（惰性行、非直接越权；最后半跳下游未实量、登记）。

## 要建什么

- [ ] **AC#1 修**（两候选，落地腿现读后裁并具名理由；**② 若触 `tools.Decision.Paths` 语义面 ⇒ 停手上报**）：
  ① `WideningRule` 入口**拒收含 `无法规范化` 标记的 path**（对齐 `session/grants.go:264` 既有口径）——预期首选；
  ② 让标记不进 `Paths`（动 `tools.Decision.Paths` 语义，须先报）。
  判据＝**能力形钉**：真跑 junction 拒卡面流 ⇒ 放宽链**不再产出碎片**（修前该形红、修后绿）＋**对照不回归**（无标记越界读仍推真父目录）＋`approval`/`session`/`bridge` 相关族 rc=0＋既有断言零放宽。
- [ ] **AC#2 边界登记**：空路径标记／8.3 短名／多路径 `ok=false` 分支——逐条量或具名"不属本修"。

## 禁区

⛔ 不动三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`／`.github/**`／`scripts/**`；⛔ `frontend/**`／`design/**`；⛔ 放宽任何既有断言；⛔ push。

## Progress log (append-only, newest last)

- [2026-10-14 16:38 +08] 编排者立票（来路＝`201-a1` 复现；现量四条；两候选；②触语义面须停报）。
