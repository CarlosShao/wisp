# 423 — 25 枚 config 字段零生产读者，但归属票已全 done：缺口没随票走，票池今天无票认领

- **Status**：ready-for-agent（编排者 2026-11-10 21:0x 立，从票 180 AC#4 只读普查腿的真发现反推）
- 来路：票 180 普查腿 `6eda2be3b`（`docs/... 见 .scratch/wisp/probes/orch/20261110/180-census.md`，102 行）。母体 70 枚带 `default:` 字段：33 有读者／7 守卫拦截／**30 零读者（5 挂未结票，25 归口悬空）**。
- ★**25 枚悬空的病**：这些字段的"归属票"（28/26/27/61/13/10/45/29/08/42/66/43 等）都已 `-done`，但字段至今零生产读者——**缺口没随票走**，今天票池没有任何一张票认领它们。改 config.toml 这些字段＝白改（没生效）。

## 现量（⛔ 引用前先重跑 180-census.md §AC#4）
- 分母＝`grep -c 'default:'` 母体 70 枚；25 枚悬空清单见 census 表。
- 零读者尺：`grep -rn "\.<字段名>"`（A1667 警告：票面/票断言会过期，起手先重跑）。

## AC
- [ ] **AC#1 逐枚坐实**：25 枚逐枚现跑 `grep -rn "\.<字段名>" --include=*.go` 出"确零生产读者"（排除 _test.go 与外设撞名），出表。
- [ ] **AC#2 分类处置**：逐枚判"该补读者 / 该删字段 / 该挂台账 DEFERRED"三档，写明理由（腿自判 ⛔ 照抄）。
- [ ] **AC#3 落地**：对判"该补读者"的，接真读者（照票 65/33-r3 既有模式）；对"该删/该 DEFERRED"的不动码，只出建议交编排者裁。
- [ ] **AC#4 门禁**：gofumpt 空／go vet 相关包 rc0／go build rc0／go test 相关包绿／d22scan clean／path-length GREEN。
- [ ] **AC#5 禁碰**：⛔ 改 PLAN/specs 契约文字（删字段若涉 D36 须机主批，只出建议）；⛔ frontend/design。

## 撤销口令
- 「**423 撤回**」。

## Progress log（423-followup 收口腿 · 2026-11-11 10:2x +08）
- 上一腿（423-wire）把真活写在盘上但未 commit、未写测试即断线；本腿收口：读懂 6 件改动、补测试、跑门禁、提交。原 6 件的语义一字未动（`resident_voice_engines_windows.go` 仅把 openScope 内联的 ScopeOptions 字面量抽成 `scopeOptions()` 方法供喂线断言——纯等价重构，行为不变）。
- **补测试**（新件 `cmd/wisp/config_wire_423_windows_test.go`，三例）：
  - `Test423AgentTailReachesAgentConfig`：config.toml 种非默认值（max_rounds=9 / token_budget=12345 / repeat_thresholds=[2,4,6]）⇒ 经装配根真跑 ⇒ 读 `rt.loopOpt.Config` 三值全数到达 agent.Config。
  - `Test423SessionKeysReachTheScopeOptions`：[session] 三键（warm 41 / conversation_idle 44 / settling 7）⇒ boot LoadFile ⇒ plan ⇒ `scopeOptions()` ⇒ 真 `session.NewScope`（录制钟）armed 44s 会话窗。
  - `Test423PunctuationSwitchGatesTheStage`：四半——① OFF 契约 `OpenPunctuation(enabled=false)`=nil,nil（伪模型路径自足，任何机器可跑）② boot 默认 true 进 plan ③ loader 把 plan 字段传进 enabled 实参 ④ acquire 闸门：关掉的 epoch 恰好 ColdLoads==1（只装 ASR）。ON 控（真模型开真级）须活 DLL：本机 System32 onnxruntime 1.17.1 遮蔽 ⇒ 诚实 SKIP。
- **反向 7 组**（摘喂线⇒红、字节级还原⇒绿，读数存 `.scratch/wisp/probes/423/followup/logs/`）：M1 agent 尾（RED：三值全 0/空）、M2a conv_idle（RED：0s→30s）、M2b settle（RED：0s）、M2c punct 喂线（RED：boot plan false）、M3a ScopeOptions 两枚（RED）、M3b loader 硬编码 true（RED：伪模型被真开）、M4 acquire 摘闸（RED：ColdLoads=2）——7/7 成对，正控/终控 GREEN。
- **门禁读数**：gofumpt 空｜go vet ./cmd/wisp rc0｜go build ./cmd/wisp rc0｜go test -run Test423 -count=1 两例 PASS＋一例 SKIP（OFF 半全跑）｜255 roster 三例 PASS（本票新增 cite token 全部唯一命中）｜d22scan clean rc0（462 生产 Go 文件）｜path-length GREEN。
- **额外回归**：TestTicket324×2、Test41a2r1GuardRows PASS。`Test15eSpeechEndToEnd` 与 `Test41a2r1OptInGate` 在无 staging 的裸机上是红，属 15e/41 系文件头写明的 staging 依赖（System32 旧 DLL 遮蔽；OptInGate 补 WISP_SILERO_MODEL 即绿），与本腿改动无关——本腿 diff 未触碰两例的任何被测路径。
- **普查账诚实说明**：census §2 的 25 枚悬空里，本票接走 **6 枚**（session.settling_sec、session.conversation_idle_sec、voice.punctuation、agent.loop_guard.repeat_thresholds、agent.max_rounds、agent.token_budget）；其余 **19 枚仍悬空**（single_instance、veto_words、tts.provider/speed、aec 两枚、conversation_mode、audio 两枚、privacy 三枚、memory 三枚、roll 两枚、observe 两枚），AC#2 的三档处置留给编排者裁。
- commit：`cmd/wisp` 六件＋新测试件＋本票面追加（两处显式 pathspec，只 commit 不 push）。