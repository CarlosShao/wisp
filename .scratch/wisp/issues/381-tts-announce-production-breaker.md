# 381 — 票 26 的播报链在生产构造器上仍是死的：打断通道零实现（R26-1）＋生产构造器缺 `provScope`/`newSpeaker`（R26-7）

- **Status**：**待派**（编排者 2026-10-27 20:3x 立；来源＝票 26 终裁 `26-v1` 件 `docs/evidence/s1/26-v1-verdict.md` §10 ＋ `26-r3` 的生产缺陷发现）
- **Blocks:** nothing · **Blocked by:** nothing（★`cmd/wisp` 现无在飞写腿）
- **写面**：`cmd/wisp/**`（★**本票独占 `cmd/wisp`**）
- **来源逐字**：★`26-v1` 判语「**R26-7（新抓）：`buildResidentVoice` 补了 queue 但**仍缺 `provScope`/`newSpeaker`**——生产出声半边对 nil 渲染缝 panic（进程不塌、门不卡、该次播报死）；测试构造器 vs 生产构造器还有第二处差异；修法＝补两字段或两构造器合一」＋「**R26-1 打断通道（承接票 NNN 或具名待开）**」

## 判据
- [ ] **AC#1 生产构造器补全**：★`buildResidentVoice` 补 `provScope`/`newSpeaker` 两字段（★**或两构造器合一、具名选型**）；★**证"测试构造器 vs 生产构造器"差异归零或逐字段具名**
- [ ] **AC#2 nil-渲染缝 panic 消除**：★**生产出声 HOP 3/4（`resident_voice_engines_windows.go:935`）不再 `nil` 函数值调用**（★**当前：registry `recover` 兜进程不塌、`defer` 仍放门、但该次播报死**）
- [ ] **AC#3 打断通道**：★**`EvInterrupt` 从"两处失败分支"接到生产可触发**（★**热键/球单击/否决词至少一条真通道**）；★**播报中打断 ⇒ TTS 停 ≤400ms 读数**（★票 26 `AC1` 的 400ms）
- [ ] **AC#4 端到端**：★**播报进行中打断 ⇒ 真停 ＋ 门放开 ＋ 下一句能进**（★**给逐跳读数**）
- [ ] **AC#5 反向判据**：★**摘掉打断通道 ⇒ 必须红**（★三读数＋还原核）
- [ ] **AC#6 纪律格**：⛔ 翻框、⛔ 改票面／`-done`、⛔ 台账、⛔ push；⛔ 改 `internal/risk`/`internal/speech` 语义；⛔ 加配置键