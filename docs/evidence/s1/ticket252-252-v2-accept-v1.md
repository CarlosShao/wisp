# 票 252 · AC#4／AC#5 收线 · 非实现者验收表（腿 `252-v2`）

被验收件：产码 `d3cccfda`（`cmd/wisp/ticket224_assembly_test.go`，+10/−1）＋交付 `c8a360e8`（`probes/252/r3/` 三件）＋票尾 `ade60808`。
本腿＝验收腿（非实现者＝与 `252-r3` 不同的角色）；⛔ 不翻框、⛔ 不动 `Status:`/票名；判语只写成立/附条件/不成立。全部读数＝本腿自跑自取（仓外两棵封存树），⛔ 未只引 r3 自述。

## §0 锚与对象（细节在 `probes/252/v2/00-anchor.md`）

- 起手 `2026-10-14T12:24:54+08:00`，HEAD `1d19b60b`；`go1.27.1 windows/amd64`；收件前 HEAD 近 `a1cf63b2`。
- 改前树＝本腿自己 `git archive 2b1a3071^`（`62ca6c34`）到 `C:\tmp\wisp252v2\tree-before`；交付树＝本腿自己 `git archive d3cccfda` 到 `C:\tmp\wisp252v2\tree-delivery`；`git diff d3cccfda 1d19b60b -- cmd internal tools scripts`＝0 行。
- 短名根＝本腿自建 `C:\tmp\wisp252v2\short-root-252v2`，`dir /x` 现读 `SHORT-~1`；每发 `GOTMPDIR` unset、`TMP=TEMP=SHORT-~1`、PATH 前置 sherpa（MSYS 形）。
- 受扰：279-r1 于 12:31/12:33 写 `cmd/wisp/run.go`／`config_readers_255.go`（12:38:46 现读 porcelain 2 行；12:44 由其 `a1cf63b2` 落件）⇒ AC#4 的绿一律改在封存交付树复跑（12:39:48 起）；三枚 witness 件全程与交付 blob 同哈希。

## §1 AC#4 判语：**成立**

判据＝三枚 witness 各配"改前必红、改后必绿"同机对照＋明说 301 秒归账。本腿复读（同机/同壳/同 PATH/同短名根；起止与 rc 逐发在 01 件）：

- W1 改前 `FAIL (43.48s)` rc=1（`:311` 命中授权仍弹 1 卡、`:327` 文件未写出）／改后 `PASS (5.76s)` rc=0。
- W2 改前 `FAIL (303.07s)` rc=1（`run_test.go:378` 逐字 `an unvetoed L1 window means EXECUTE, got: 审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝`）／改后 `PASS (5.15s)` rc=0。
- W3 改前 `FAIL (6.55s)` rc=1（`run_mode101_test.go:506` 逐字含 `审批超时（1 秒未确认）`）／改后 `PASS (83.11s)` rc=0（其载具两枚 40s 窗口，非挂死）。
- 与 r3 同形（r3→v2：43.52→43.48／303.23→303.07／6.20→6.55；5.55→5.76／5.78→5.15／83.47→83.11）。
- ★301 秒归账：票尾与 01 件原句逐字在（引文在 01 件 §301）；因果本腿自核成立——超时句只出自 `internal/agent/approval/queue.go:479`，进队列＝非 L1（`bridge.go:351` 路由注释 `L0 pass, L1 window, L2 approval`）；L2 的因＝短/长拼法分岔 R2（`rules_gateway.go:49`；`2b1a3071` 自陈；overlay 红 dump 逐字 `risk=L2 … in_allowlist_scope=false … kind=approval-timeout`）；封存交付树 5.15s 绿＝两树唯一变量是折叠修复。⇒ 账归票 252，不留两笔。
- 具名边界：W2 自己的日志不含审计行（其断言路径不 dump）；W2 侧因果按"红句（队列入队）＋同形分类直读（W1 发 overlay 红）＋封存绿"三层给，不是单发直读。

## §2 AC#5 判语：**成立**

- `d3cccfda` 单文件单 hunk `+10/−1`：只给 `grantLinesOf` 滤网加一枚针 `tools: call ` 与注释；4 处调用点全是消息参数 ⇒ 零断言放宽、产码语义零动（`--numstat`＝`10 1` 单文件；非测试文件 0）。
- 改前 grep 计数 0（W1 红发全日志 `tools: call`＝0／`in_allowlist_scope`＝0；dump 只有 SESSION-MINT＋GRANT-RECORD 两行）。
- 本腿自造同形红（overlay 把交付件叠进改前树＋短名；`FAIL (83.11s)` rc=1）：dump 逐字含 `[audit] tools: call … risk=L2 … in_allowlist_scope=false … rules_hit=[R1 R2] grant_id=0 reason="… R2: 目标路径在授权目录之外: …SHORT-~1…"`（2 枚：ask/hit）；计数 2/2。
- 还原证：工作树件 `366da8f5…` 发前发后同；`tree-before` 件 `8ac227a8…` 发前发后同（overlay 零写盘）。
- 具名边界：这枚针是**仪器**（让下次同因红一眼可读），不是新断言；AC#5 判据原文要求的就是"dump 那行等级"，已满足。

## §3 门禁抽核（本腿现跑）

- build rc=0（12:42:18→12:42:22，0 B）；d22scan rc=0 `clean - no D22 ban violations`（12:42:29→12:42:56）；长度门 rc=0 `VERDICT GREEN`／`over-budget=59 covered by roster=59 not in roster=0`（12:42:57→12:42:58）。门禁跑时 279-r1 改动在场（12:44 已落件 `a1cf63b2`）；三门禁全绿不受其影响。

## §4 翻框建议（判语归编排者；⛔ 本腿未翻）

- **AC#4 建议翻勾**。留档句：`252-v2（非实现者）：三 witness 同机复读改前红/交付绿（43.48/303.07/6.55 ↔ 5.76/5.15/83.11，rc 逐发 1→0）＋301 秒账归本票句逐字在、因果自核成立（超时句只出自 approval/queue.go:479；进队＝非 L1＝R2/L2 之卡）；翻。`
- **AC#5 建议翻勾**。留档句：`252-v2（非实现者）：d3cccfda 单 hunk +10/−1 只加 dump 滤针（tools: call ），4 调用点全为消息参数、零断言放宽；改前 grep 0/0 → 自造同形红 dump 含 risk=L2/in_allowlist_scope=false（2/2）；overlay 零写盘 sha256 两件同；翻。`
- 其余格（AC#1/AC#2/AC#3/AC#6）本腿未碰、无判语。

## §5 本腿没做／未定义即停核

- ⛔ 未跑整包；⛔ 未 push（只 commit、显式 pathspec）；⛔ 未翻框、未改 `Status:`/票名；零 Go 文件改动。停手条件核：短名根拿得到（自建＋`dir /x` 读得 `SHORT-~1`）；改前红复得出（三枚全红）；针的 diff 与"只加滤针"相符（单文件/单 hunk/+10/−1）⇒ 无停手项。
- 原始日志与 overlay 件留盘位置具名在 `probes/252/v2/00-anchor.md`；交件＝本件＋`probes/252/v2/`（00／01／02＋`logs/10`／`logs/20`）＋票 252 尾 1 条 bullet。
