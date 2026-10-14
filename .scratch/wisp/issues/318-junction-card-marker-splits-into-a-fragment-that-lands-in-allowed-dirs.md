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
- [2026-10-14 17:02 +08] 写腿 318-r1 落地（实现者；锚 `5893baaa`；产码/交付两笔＝①`internal/agent/approval/{replies.go,ticket318_widening_marker_rulers_test.go,zz_201_a1_junction_marker_windows_test.go}` ②`probes/318/r1/**`＋本 bullet，sha 见编排者收件）；**选案①**＝`WideningRule` 入口**整卡拒收**含 `无法规范化` 标记的 path（5 行守卫＋注释；对齐 `grants.go:260-268`＋`Covering:193-197`＋`scopeVerdictOf:1273` 的 any-path 形；第二裁决＝不做逐 path 跳过——那会让修前拒的混合卡反向开口，M2 为证）；候选②（动 `tools.Decision.Paths`）**未选、未触，无需停报**。读数：**修前红** `logs/10`（三形红句：正/反斜杠碎片 `dir="…(junction"`＋空路径标记 `dir=" (无法规范化: tools: empty path)"`；翻后探针 2b 红）→**修后绿** `logs/23`（真 NTFS junction 链：READING-2 标记仍在卡面／2b 两拼法同拒／4-6 `ok=false` 空 dir rule／**对照 READING-11 无标记越界读仍推真父目录**）；**突变 3 发 overlay**（M1 摘守卫／M3 换判词⇒两钉红；**M2 逐 path 跳过⇒只"整卡拒"钉 `:84` 逐字红**），工作树零突变。门禁全 rc=0（`logs/23/24/42/43/44/53/54`）：approval 全包、grants 先例族 11 枚、tools junction 族 11 枚、cmd/wisp always 两枚（3.70/3.42s，本波唯一涉窗笼跑）、build/vet/gofumpt v0.12.0/d22scan clean/长度门 GREEN 60=60；**session 全包 rc=1＝既有间歇 `TestSessionListByWorkspace`**（台账 `:17402`/`:17760` 在册＋`session` 零 `approval` 依赖＋本腿 `-count=3` 2 红 1 绿）。探针件翻钉：`TestZZ201A1JunctionMarkerBecomesAWideningFragment` → `…CardIsRefusedForWidening`（前缀保留），legs 4-5 按"主题不存在"删、修前读数留 `probes/201/a1/logs/run3-full.txt`。⛔ 未翻框、⛔ 未动 `Status:`/票名、⛔ 未 push、⛔ 未触 `cmd/wisp/**` 产码与 `tools.Decision.Paths`；AC#2 三格逐条（空路径=已量／8.3=具名不属＋反证钉／多路径=已量）见交件 §6。
- [2026-10-14 18:3x +08] agent=318-v1 did=非实现者独立复核（AC#1 三判据／被改 201-a1 件翻面／AC#2 登记／名册零放宽）：仓外 `git archive` 双快照（`8707a7ed^`／`f6fda72d`）＋我自造复现尺（真 `AddAllowedDir` 写盘：修前 rc=1 碎片真写 `allowed_dirs = ['…(junction']`／修后 rc=0 `ok=false` 零写入／摘守卫还原再红再绿）＋自造突变两发（M1 摘守卫 ⇒ `ticket318_…:58/:61`＋zz`:191` 逐字红；M2 逐 path 跳过 ⇒ `:84` 逐字红＋我方混合卡尺独立复到反向开口；两发 build rc=0、仓内零突变、还原 9 PASS）⇒ **AC#1 判据①/②/③ 全成立**。201-a1 件判语＝**诚实翻面（附条件）**：同链同步骤在、2b/leg3 断言翻向非变松、legs 4-5 删＝主题不可存在（修前读数留 `probes/201/a1/logs/run3-full.txt` 未改写）、对照腿加强、`TestZZ201A1` 前缀与零 skip 保留；残余面＝该半跳无修后测试直驱生产消费方（源码现读 `approval_always.go:85-91` 折返，⛔ 未跑整包 cmd/wisp）。**AC#2 三条登记全核到**（空路径=两树已量；8.3=不属本修＋tools 族 11 PASS 我现跑；多路径=已量＋跳过语义反向开口独立复到）。名册 79→83 恰=+4 新增＋1 改名、SKIP 1→1 零新增、全包两树 rc=0、审批包面 `git status` 零改动。件＝`docs/evidence/s1/ticket318-318-v1-accept-v1.md`＋`probes/318/v1/logs/` 7 枚；观察两条（「一直」拒绝文案字面／`replies.go:503` 旧行号）见交件 §0。next=编排者裁定（AC#1/AC#2 勾与收口）；⛔ 未翻框、⛔ 未动 `Status:`、⛔ 未 push。
- [2026-10-14 18:40 +08] agent=318-r1 did=按 318-v1 附条件做最小修复（仅两处头注；零行为改动）commit=59e2b3ed next=编排者翻框
