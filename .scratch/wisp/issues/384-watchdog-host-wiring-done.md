# 384 — watchdog 宿主接线：真 TreeReader/StateFn/OnAlert→D43 映射/四组件注册＋Start（生产调用者 0）

**Status:** **DONE-带残余**（2026-10-31 编排者收口：五格全翻＝384-v1 判能翻；273 处置裁不构成伪完成；残余九条见末节）
- **Blocks:** nothing · **Blocked by:** ★**`cmd/wisp` 每波一枚**——**排 `41a2-r1` 交件之后**；★**OnAlert→D43 映射依赖票 273（状态机无 sink）的真收件人 ⇒ 与 273 同波或其后**
- **写面**：`cmd/wisp/**`
- **来源逐字**：★`42-v1`「生产调用者 0＝真缺口，归宿主接线腿：真 TreeReader、StateFn 接状态机、OnAlert 映射（Armed AlertOverrun→EvWatchdogOverrun/#9、AlertFatal→EvWatchdogFatal/#35、Leak→日志/UI）、四组件注册＋Start」

## 判据
- [x] **AC#1 生产装配**：★真 TreeReader／StateFn／OnAlert 映射（#9/#35 逐字）/四组件注册＋Start（★对照 `42r1-impl.md` 的缝形状）
- [x] **AC#2 端到端**：★over-cap ⇒ 真 unload 趟＋D43#9 留态；3× settle 败 ⇒ #35 `WatchdogAlert`＋one-click-restart（★真机读数）
- [x] **AC#3 反向判据**：★摘 OnAlert 映射 ⇒ D43 事件零派发 ⇒ 红（三读数＋还原核）
- [x] **AC#4 零回归**：★15e/26R3/381 打断案/41a2 KWS 案全绿；`drained` 抖动只报
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ 改 `internal/watchdog`/`internal/session` 语义；⛔ D43 表；⛔ 配置键
## ★★★终裁落地（2026-10-31 16:0x 编排者，凭非实现者 `384-v1`，件 `docs/evidence/s1/384-v1-verdict.md`、`12f5ad324`）

★**五格全翻（8/8 自跑）⇒ 可带残余结案**。
★**四缝逐缝读码为真**：TreeReader＝真 C30 Job；StateFn＝20 态直映 4＋16 fail-closed（**D32③ 反面有测试钉**）；OnAlert＝**全文仅两枚 Dispatch、D43 表各恰一行＝零发明边**；panel-webview＝真注册＋三 skip 臂＋Destroy 可恢复；**三组件不注册＝白名单查表天然跳过（无造假组件）**。
★★**273 sink 处置判"不构成伪完成"**：★**转移是状态机的、投递是 273 的、两件事**——nil⇒显式 no-op 是包冻结契约；273 落甲形本缝零改动。
★**mutation 两发复刻全红**（摘 #9／摘 #35 各 rc=1、红句逐字、#9 effect=1 自证定向）；★零回归 19/19＋Test43 9/9＋255 名册 20/20＋`262/640/902` 逐位同。
★**归因精化**：★**Test256 红（unexpected Options field "Answer"）＝票 375-r1（`b63d56413`）引入 ⇒ 375 族债**（非 381 族——实现件 §7.5 需更正、结论"HEAD 既有"不变）；33r11 两枚＝381 族债 ✓。
★**残余九条具名不阻结案**：三组件无生产对象／273 sink 未结／ASR 冷装载全量串跑抖动＝测试基建债／33r11+256 两族债／VerifySettle 随票 28／panel 真窗口 Destroy 读数属 winlive 未批／staging 未追踪资产经训（junction 零拷贝法）。
