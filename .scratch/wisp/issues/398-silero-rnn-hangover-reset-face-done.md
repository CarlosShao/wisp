# 398 — silero RNN hangover：internal/speech 无 session reset 面（392-v3 采纳路径②、非阻塞）

**Status:** **DONE-带残余**（2026-11-01 编排者收口：五格全过（398-v2 七项自跑）；口径差裁定成立=(0,1) 正确口径；残余三条见末节）
- **写面**：`internal/speech/**`（★387-r3 已让出）
- **机理（392-v3 三进程同读）**：★**silero RNN 语境跨空窗携带**（无输入不衰减）、**hangover 能在 post-speech 静音上开火**（探针 silence=1 talker=0）；**冷段吸收之**——但「更硬的清面」＝internal/speech 需 session reset 面（**新 API、契约相邻**）

## 判据（派时细化为逐格）
- [x] **AC#1 reset 面设计**（★对照 detector Finalize 形状；★**⛔ 动 silero 语义**）
- [x] **AC#2 hangover 清除读数**
- [x] **AC#3 反向判据**
- [x] **AC#4 零回归**（★internal/speech 全套 61/0/6＋裸跑 25/0/42）
- [x] **AC#5 纪律格**
## ★★★终裁落地（2026-11-01 19:3x 编排者，凭非实现者 `398-v2`，件 `docs/evidence/s1/398-v2-verdict.md` 116 行、`78349e96c`）

★**五格全过 ⇒ 可收口**（七项全本腿自跑、导出树、`git diff ced345c3c HEAD -- internal/speech/`=空＝实现自落地零再动）。
★**398-v2 七项**：①**形状核 ✓**（Reset() 三锚点现核：binding impl 指针／vad.Reset()→C.SherpaOnnxVoiceActivityDetectorReset sherpa_onnx.go:1787-1788／c-api.h:2235 注释逐字；go.mod pin v1.13.8；Accept/Finalize/Consume/Close/Endpointer 零字节）②**探针复跑 ✓**（PRE-RESET (1,0)／POST-RESET (0,1)／RESET-EQUALS-FRESH 两臂全等——与 398-r1 申报逐字同）③**口径差裁定：成立、实现侧 (0,1) 是正确口径**——silence=0 是主判据、talker=1 是判定通路活性内建正控；**字面 0/0 恰＝reset 清死判定通路＝票面停手③**；post 臂断言硬钉 (0,1)、判据差不是放宽 ④**等价案核 ✓**（per-segment 三段计数落点敏感、kind 序列初版误绿史核过、**改法用突变实测证明咬合**）⑤**反向复刻＋还原核 ✓**（导出树自建突变二进制——主案红＋具名句「the hangover survived the reset」＋PRE 控制臂不动＝单变量归因；还原核四件全过）⑥**零回归 ✓**（staged 69/63/0/6＋裸跑 25/0/44＋cmd KWS 族 2/2＋ban#8 732→734 具名归属 394-r1/400-r1 与本票无关）⑦**五格判 AC#1-#5 全 ✓**
★**残余三条维持**：C++ 源码级核需可通 GitHub 环境（四路解析均被拦）／**owner 接入（Finalize(nil)+Reset() 免冷段）属 392-v3 §6 第 3 步后续裁定**／hangover 衰减扫参未关闭
★**398-v2 注记收档**：基线 staged 首发红＝导出树形状缺口（runtime.Caller walk-up 跨树不自洽）、junction 补齐后 61/0/6 全绿——非基线红
