# 383 — `CheckAndReload` 失败路径无条件收养 stat：一次瞬时 sharing violation 把该次配置写永久隐藏（R-42r1-1）

- **Status**：**实施中**（编排者 2026-10-31 11:5x 开；来源＝票 42 终裁 `42-v1` §6 ＋ `42-r1` 并发轮询实证；**`config-r1` 已在飞执行本票**）
- **Blocks:** nothing · **Blocked by:** nothing（`internal/config` 空闲）
- **写面**：`internal/config/**`（★**`cmd/wisp` ⛔——41a2-r1 独占**）
- **机理（`42-v1` 读码复认）**：★`manager.go:160-168` 失败支**无条件收养 stat**（`seenValid=true`）⇒ ★**与原子 rename 竞态的一次瞬时 sharing violation 读失败 ⇒ 该次写永久隐藏**
- ★**修法（`42-v1` 判"今天可行"）**：★**错误链已可分**——`loader.go:81` read 支（原始 `os.ReadFile` err，可 `errors.Is/As`）vs `:102/:118` parse 支 ⇒ ★**解析/校验类失败维持收养（防循环重读）；IO 类失败 ⛔ 收养（保留旧 stat ⇒ 下 tick 重试，成功生效、失败每 tick 一条 Warn 有界自愈，223 已有 `describeReloadFailure` 出线）**

## 判据
- [ ] **AC#1 分类落地**：★read 支 IO 错误 ⇒ ⛔ 收养；★parse 支 ⇒ 收养保持（★`cfg125-r1` 的解析类收养语义零回归）
- [ ] **AC#2 IO 失败一发**：★模拟瞬时 sharing violation（★共享冲突打开文件把手等可信形；★⛔ sleep+hope）⇒ **下次 tick 重试成功 ⇒ 该次写不被隐藏**
- [ ] **AC#3 解析失败回归一发**：★收养保持（防循环）
- [ ] **AC#4 反向判据**：★摘分类（回到一律收养）⇒ IO 发必红（三读数＋还原核）
- [ ] **AC#5 零回归**：★233 热档（`Test233HotTierAgentReachesTheNextTask`）还绿／`cfg125` 解析类语义零回归／`42-r1` 的 `config_poll_test` 不退化
- [ ] **AC#6 纪律格**：⛔ 翻框/改票面/台账/push；⛔ thresholds/golden；⛔ `cmd/wisp`／`internal/watchdog`（42 刚收口）