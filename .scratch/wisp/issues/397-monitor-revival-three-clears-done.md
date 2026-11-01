# 397 — monitor 复活三处窗间清一体落地（392-v3 采纳路径①）

**Status:** **DONE-带残余**（2026-11-01 编排者收口：五格全过（397-v2 重派腿七项自跑）；三条新残余＝默认翻转腿须知；见末节）
- **写面**：`internal/audio/**`（★393-r2/392-r2 已让出）
- **机理（392-v3 矩阵坐实）**：★**APM 持有＋窗间 Reset()（统一器）＋resampler 持有＋窗间 Flush()＋detector 持有＋窗间 Finalize(nil)＋1.2s 静音冷段** ⇒ W2/W3 全愈三进程复现；★**摘任一件即三失败形之一＝联合承重**；★APM 选 Reset 形＝统一器兼省构造、零新 API

## 判据
- [x] **AC#1 三处窗间清一体落地**（★**monitor release()/keepalive 分支**；★**keepalive=false 默认 ⛔ 动——显式开启读数**）
- [x] **AC#2 保温复活读数**：★**第二暖窗重新会判**（对照 392-r2 的 revival-cool 形状）
- [x] **AC#3 反向判据**（★摘任一件 ⇒ 对应失败形复现——三发；还原核）
- [x] **AC#4 零回归**（★internal/audio 全套 94/83/0/11 基线＋392/393 面不退化）
- [x] **AC#5 纪律格**：⛔ 翻框/改票面/台账/push；⛔ AEC3 源树／`cmd/wisp`／thresholds
## ★★★终裁落地（2026-11-01 18:5x 编排者，凭非实现者 `397-v2`（重派腿、397-v1 判死后代跑），件 `docs/evidence/s1/397-v2-verdict.md` 73 行、`d9786ce33`）

★**五格全过 ⇒ 可收口（附三条新具名残余）**。
★**397-v2 七项**：①AC#1 读码核＝**keepalive=false 默认逐字节不动**（旧 :294＝新 :316 逐字同、closeEngines/dispose/readings/drain 零 hunk）；旧形对照＝旧分支 APM 不持有无条件 Close——**申报「替换原 Close/nil 形」属实** ②AC#2＝**断言改钉 WARM REVIVAL MOVED（翻动归属本腿＝A1355 裁定）＋voice392RoundOut 全导出壳 +39 行、392-r1 空壳现读实锤（cold_rounds=[,]、四 round 全 {}）→壳修复闭环**；★**权威读数复跑 rc=0（38.12s 七窗全序）：warm2=episodes=1/triggers=1/engine_warm=true/engine_built=false/ready=9728 逐位同/rms=-25.71@-30** ③AC#3＝**四发 overlay（diff 恰目标块、四 SHA 互异）**：**M3 摘 Finalize+冷段＝判死形两发红**（首发绿＝孤立摆动如实记）＋**M4 摘冷段＝瞬态形红**；M1 摘 Reset＝行为断言全过×3（**复刻形 triggers=0 未迁移**）；M2 摘 Flush＝rc=0 ④AC#4＝**26p 定向 keyed 4/4**；(c) 397 树三发全绿 vs 对照树 **红名归因闭合**（摆动族）；调参族 rc=0；六族＝**34/1/0 逐位归位**＋**`262/640/902` 逐位**；四数＝**RUN=98/PASS=75/FAIL=0/SKIP=11 逐位同基线**（**票面 94/83/0/11＝A1326 时代旧口径、台账 :29600 已勘误**）；./... 三跑红全在 package-main 编译域外（**cmd/wisp 唯一红=TestReadGitOnThisRepositoryIsSelfConsistent＝导出树无 .git 的结构预期红、与 397 零关系**）
★**三条新残余（归默认翻转腿/编排者知悉）**：1. **M3 判死形 staged 链非逐发确定**（1 绿 3 红跨窗）＝真 TTS 摆动族、**多发稳定性归默认翻转腿全链复读** 2. **M1 仪器面 -Inf marshal 脆弱性**：**壳修复把 rounds 首次接上 JSON 后，(b) 轮 episodes=0（完全门净健康形）→ write392JSON 整件 `json: unsupported value: -Inf` → 测试尾部红（行为断言全过）**；**默认翻转腿全链复读必撞面，建议 JSON 安全哨兵小改测试件** 3. **导出树实验须补 gitignored 面**：**build/models/manifest.json 对＋third_party/model-fetch/ 实体**（A1326 卡在导出树同样适用）——本腿六族假红两轮由此；★**另坑：失败发不写 artifact→工件目录残留上一过跑件（392-keepalive-m3.json 内容与其 m3 日志 FAIL 矛盾即此）、证据以日志 PASS/FAIL 行为准**
