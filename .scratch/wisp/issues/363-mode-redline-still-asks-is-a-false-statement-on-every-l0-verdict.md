# 363 — 生产日志每一发 L0 都在说假话：`MODE-REDLINE … still asks` 而实际放行

- Status: **ready-for-agent**
- 立案：2026-10-22 23:2x +08 编排者立
- ★**来路（非实现者签名）**：`179-r1`（写腿，件 `docs/evidence/s1/179-r1-impl.md`、4 笔 `32c8eeac`／`18e3f670`／`0aa30ede`／`8956ece7`）§7① 具名上报
- ★**编排者已现跑核实四条链**（见下）

## 这是什么（人话）

**生产日志在每一发 L0（最低风险级）上都打印"仍然会问"，而实际上 L0 直接放行、一张确认卡都不弹。**

⇒ ★**这是向操作员说假话**：他读到"还会问你"，于是**以为系统会拦一下**，实际**已经执行完了**。

## ★四条链（编排者现跑核实，逐条带 `file:line`）

1. **`internal/tools/bridge.go:311`** 逐字 `sil := mode.Screen(verdict)` ⇒ **L0 的判定被交给 `mode.Screen`**。
2. **`internal/risk/mode.go:163-165`** 逐字：
   ```
   if d.Level != L1 && d.Level != L2 {
       return Silenced{Level: d.Level, Kept: "无法识别的风险级，fail-closed 保留询问", Mode: m...}
   }
   ```
   ⇒ ★**L0 ≠ L1 且 L0 ≠ L2 ⇒ L0 落进这一支**，被当成"**无法识别的风险级**"。
3. **`internal/tools/bridge.go:324`** 逐字：`case sil.Kept != "": b.log("tools: MODE-REDLINE mode=%s still asks tool=%s rules=%v because=%s", …)`
   ⇒ ★**只要 `sil.Kept` 非空就打印 "still asks"**、**⛔ 区分那是不是真的会问**。
4. **`internal/tools/bridge.go:352`** 逐字注释 `// --- routing: L0 pass, L1 window, L2 approval, Deny never reaches a gate ---` ＋ `b.route(ctx, &dec, sil, grantID)`
   ⇒ ★**`route` 对 L0 直接放行** ⇒ **0 张卡**。

⇒ ★★**链条闭合**：L0 → 被 `Screen` 误判为"无法识别的风险级" → 日志说 "still asks" → 实际 `route` 放行、**0 卡**。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把这条假话量成可复跑读数——
      ①**逐条复核上述四条链**（每条带 `file:line` 现跑复认，⛔ 抄本票）；
      ②★**造一发真读数**：跑一发 L0 工具调用，**逐字贴出那行 `MODE-REDLINE … still asks`**，并**同一发里证明确认卡数 ＝ 0**（★**这就是"假话"的直接证据**）；
      ③**正控**：造一发 **L1 或 L2**（那两档 `Screen` 认得）⇒ ★**它的日志必须 ⛔ 说 "still asks"**（或说的是**真的**会问）⇒ **证明这条链只在 L0 上假**；
      ④**票面前提新鲜度**：行号漂不漂要现跑。
- [ ] **AC#1 裁形（编排者拍）**：三形各写代价＝〔甲〕**只改日志文案**（把 "still asks" 改成能说实话的措辞，例如带上 `sil.Kept` 的原文、或对 L0 单独一支）；
      〔乙〕**`mode.Screen` 认 L0**（让它返回 `Kept=""`、不进"无法识别"支）；〔丙〕**两者都做**。★**未裁形任何实现都⛔ 开工**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数（变异红／还原绿／原样绿）逐字贴 ＋ **还原核**。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet` 动过的包／`go test -count=1`（**四数分开报、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
      ★★**`cmd/wisp` 整包四数 ⛔ 当判据**（未判决恒 1、红名册逐发漂移、截断点后整族 `RUN=0`）⇒ 只按名作差＋逐发登记死点。
- [ ] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 179/183/185 任何原句。

## ★写面（初判，裁形后可能收窄）

- ★**最可能落在** `internal/tools/bridge.go`（日志）与 `internal/risk/mode.go`（`Screen` 的 L0 支）——★**两处都在别家地界** ⇒ 落地时**必须一次只碰一处**、且**与那两包的其他腿串行**。
- ⛔ `docs/PLAN.md`／`docs/specs/**`（契约文字须人工批准）、`thresholds.go`／golden／`allowlist.txt`／`tools/d22scan/**`／`frontend/**`／`design/**`。

## ★两条硬提醒（⛔ 踩）

1. ★**⛔ 为了让日志"好看"而放宽任何判定**：本票要的是**日志说实话**、⛔ 把 L0 改成会问。
2. ★**⛔ 只改文案就交差**——若 `mode.Screen` 把 L0 判成"无法识别的风险级"这件事本身是错的（★**L0 是认得的级**），那它**在别处也会假**（不止日志）⇒ 裁形时要判这一点。

## 排程

`363-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `363-r1`（写腿＝`AC#2`）→ `363-v1`（非实现者裁）。

## Progress log

- [2026-10-22 23:2x +08] agent=`编排者` did=立案（来路＝`179-r1` §7①；★四条链编排者已现跑核实：`bridge.go:311`／`mode.go:163-165`／`bridge.go:324`／`bridge.go:352`）next=派 `363-a1` 只读腿。
