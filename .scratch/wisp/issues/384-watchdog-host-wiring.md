# 384 — watchdog 宿主接线：真 TreeReader/StateFn/OnAlert→D43 映射/四组件注册＋Start（生产调用者 0）

- **Status**：**待派**（编排者 2026-10-31 11:5x 开；来源＝票 42 终裁 `42-v1` 装配缺口＋票 42 残余①；★**对照 28-r1"库＋读数齐、装配等下一波"同形状**）
- **Blocks:** nothing · **Blocked by:** ★**`cmd/wisp` 每波一枚**——**排 `41a2-r1` 交件之后**；★**OnAlert→D43 映射依赖票 273（状态机无 sink）的真收件人 ⇒ 与 273 同波或其后**
- **写面**：`cmd/wisp/**`
- **来源逐字**：★`42-v1`「生产调用者 0＝真缺口，归宿主接线腿：真 TreeReader、StateFn 接状态机、OnAlert 映射（Armed AlertOverrun→EvWatchdogOverrun/#9、AlertFatal→EvWatchdogFatal/#35、Leak→日志/UI）、四组件注册＋Start」

## 判据
- [ ] **AC#1 生产装配**：★真 TreeReader／StateFn／OnAlert 映射（#9/#35 逐字）/四组件注册＋Start（★对照 `42r1-impl.md` 的缝形状）
- [ ] **AC#2 端到端**：★over-cap ⇒ 真 unload 趟＋D43#9 留态；3× settle 败 ⇒ #35 `WatchdogAlert`＋one-click-restart（★真机读数）
- [ ] **AC#3 反向判据**：★摘 OnAlert 映射 ⇒ D43 事件零派发 ⇒ 红（三读数＋还原核）
- [ ] **AC#4 零回归**：★15e/26R3/381 打断案/41a2 KWS 案全绿；`drained` 抖动只报
- [ ] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ 改 `internal/watchdog`/`internal/session` 语义；⛔ D43 表；⛔ 配置键