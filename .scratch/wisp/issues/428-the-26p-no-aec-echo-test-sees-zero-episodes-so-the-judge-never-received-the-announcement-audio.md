# 428 — `Test26PWithoutAECTheEchoTranscribesVerbatim` 恒红：无 AEC 时麦克风面收到的自己播报音频，回声判定器**零 episode**

- Status: **ready-for-agent**（先只读取证，判定后写面才定；⚠ 预期写面在 `cmd/wisp/**` 的台面侧或 `internal/voice/**`，**由你判清后自己定并写理由**）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）
- 关联：票 26P（回声／AEC 族宿主）· 票 59（rig 来源，红句里点名）· 票 392（AEC 调参）· 台账 `A1801`／`A1803`

## 这是什么（人话）

有个测试要证明：**不开回声消除（AEC）时，麦克风会把自己刚播出去的声音收回来**，
回声判定器应该**看见这些"回声片段"**（episode），然后据此决定不打断自己。
现在实测：**判定器一个片段都没看见**（零 episode）⇒ 红。

## 现量（可逐字复跑）

**逐字红句**（`cmd/wisp/resident_voice_bargein_26p_echo_windows_test.go:216`）：

> `MUTANT 26p-M2 SHAPE: the raw echo opened zero episodes - the judge did not see the announcement audio the mic surface carried (the 59-r1 rig …)`

（红句被日志截断，**你要去源码里抄出 `:216` 那条 `t.Fatalf` 的完整格式串**。）

**稳定性**：本轮整包一次红。⚠ **只跑过一次，⛔ 尚未确认是稳定红还是偶发** —— 你必须先复跑 ≥2 次定性。

## 机制链（现读，供你起步；⛔ 别照抄结论）

测试文件：`cmd/wisp/resident_voice_bargein_26p_echo_windows_test.go`
- `:203` `rig.rv.announce(voice26pAnnounceA)` —— 播报一句
- `:206` `spy := rig.spyAt(1)`／`:207` `<-spy.first` —— 等第一段
- `:208` `stopEcho, _ := rig.inject(echoWav)` —— **注入回声 wav**
- `:209-211` 等播报结束（`ttsAborted >= 1 || ttsSpoken >= 2`）
- `:212` `stopEcho()`／`:213` `rig.waitReleased()`
- `:214` `read := rig.rv.bargein.readings()`
- `:215` `if len(read.Episodes) == 0 {` ⇒ **`:216` Fatal**
- `:218-220` 之后统计 `ep.EnergyPass`

⇒ **失败点是"判定器没开 episode"**，不是"episode 开了但判错"。这两件事的修法完全不同。

## ⚠ 必须自己判清的三个岔路（⛔ 不许我替你定）

1. **是"音频没送到判定器"还是"送到了但判定器没开 episode"？**
   - 前者＝台面/注入问题（`rig.inject` 那条路）；后者＝判定器行为问题。
   - 判据：在 `:214` 之前插探针（仓外 overlay，**但这枚是运行时行为检查，overlay 有效**），
     数一下判定器**收到多少帧**。收到 0 帧 ⇒ 前者；收到很多帧但 `Episodes` 空 ⇒ 后者。
2. **`WISP_SPIKE_MODELS` 等夹具环境是否齐？** 本轮 carrier 已自动解析五个变量（见 `docs/evidence/s1/voice-fixture-env.md`），
   但**这枚测试是否需要 `WISP_26P_*` 一类别名变量**？去读测试的 `TestMain`／`new26pRig` 看它读哪些 env。
3. **是否与票 392 的 AEC 调参有关？** 读 `Test392AecTuningDoubleTalkSwing` 与 `Test392KeepalivePostureFindingsAndColdDefault`
   本轮的状态（**这两枚本轮 PASS**），确认判定器在别的用例里是活的。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. **复跑 ≥2 次**定性（稳定红 or 偶发）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test26PWithoutAECTheEchoTranscribesVerbatim$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **抄出 `:216` 完整格式串**（日志里被截断）。
  3. 读 `new26pRig`／`rig.inject`／`rig.rv.bargein.readings()` 三处，写清**音频从注入到判定器读数的完整路径**。
  4. 用**仓外 overlay 探针**数判定器实际收到的帧数（**这枚是运行时行为检查，overlay 有效**）。
  5. 判清上述三个岔路，给一句话结论＋支撑的 `文件:行` 与逐字读数。

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - 若判「台面/注入问题」⇒ 修台面（测试侧或 rig），⛔ 不许放宽断言。
  - 若判「判定器行为问题」⇒ 修产码。
  - **必须配突变控制真跑**：把修好的那行/那处删掉 ⇒ 断言必须回到红；恢复 ⇒ 绿。
  - ⚠ **仪器选择**：这枚是**运行时行为**检查 ⇒ `go test -overlay`（拷贝放**仓外**）有效。
    （对照：票 427 那枚**读磁盘源码**的静态检查，overlay **无效**，必须磁盘改写。见台账 `A1804`。）

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，那是独立 module，从根模块跑会失败）。四条读数全贴。
  ⚠ 已知残留：`build/mut_rvw.go` 是**未跟踪、被 `.gitignore` 吃掉的残留件**，会让 `go build ./...` 红。**⛔ 不是你的写面，⛔ 不许删**；用 `go build ./cmd/wisp/`。

- [ ] **AC#4 与票 26P／392 的关系写清**：本票只处理这一枚，⛔ 不翻票 26P／392 任何框；若判定涉及 AEC 调参语义，只在证据件写明「需回填」，⛔ 不自行改票。

## 本票**不**解决

- 不碰 `Test26P*` 另外四枚（`EchoNeverTriggersAndNeverTranscribes`／`MarginRetestAndGateReverse`／`BargeinStopsWithin400ms`／`FeedOnlineMatchesWindowReplay`）——**本轮它们的状态你要一并记录**（若有红，写清是同一根因还是不同）。
- 不碰 `Test384*`／`Test365*`。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面；**产码与台面二选一，⛔ 两边都改**（否则突变控制证明不了任何事）。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push（推送由编排者做）。
- ⛔ 翻任何 AC 框（翻框由编排者做，非实现者）。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\428-<腿码>\`。

## Progress log

### 2026-11-12 20:26–20:5x +08 — 腿 `428-26p`（⛔ 未翻任何 AC 框）

**判定＝台面问题（"音频没送到判定器"），根因是判定器从未 armed，不是判定器行为。**
一句话：commit `e83494cf5`（2026-11-11 180-census leg K）给 `newResidentVoice` 加了
`rv.bargein.setAECEnabled(plan.AECEnabled)`（`resident_voice_windows.go:591`），而 `newVoice26pRig` 手搓的
`voicePlan` 从不设这一格（`resident_voice_bargein_26p_windows_test.go:139-143`）⇒ 零值 `false` 让
`arm()` 的第一个 case 直接拒绝（`resident_voice_bargein_windows.go:392-393`），整族判定器再没起来过。

- **稳定性＝稳定红，3/3**：22.86s / 24.42s / 23.81s，同句同读数。⛔ 不是偶发。
- **★决定性读数（仓外 overlay 探针插在 `:214` 之前）**：`armed=false`、
  `reason="[voice] aec.enabled=false：AEC 插话检测不装配（…P15/A1307…）"`、**`frames_drained=0`**、
  **`processed_samples=0`**、`episodes=0`、`aec_bypass=true`；同时 `tts samples=281536`（播报确实出声）。
  ⇒ **判定器实际收到 0 帧** ⇒ 按票面判据＝前者（台面/注入），后者（判定器没开 episode）**排除**。
- **同一探针的对照发**：唯一差别是 plan 补 `AECEnabled: true` ⇒ **`frames_drained=200`、
  `processed_samples=102080`、`episodes=1`**、`lag=501`、逐字转写
  `"我已经完成明天下午三点的会议改到四点请准时参加谢谢"`（`distinctive=true`）。
  日志行 `voice: barge-in monitor armed for the announcement window` 修复前**零次**、修复后两次。
  ⇒ 断点在 **arm**，不在 `rig.inject`。
- **三个岔路的另外两条**：(2) **没有 `WISP_26P_*` 别名变量**（`git grep` 零命中）；真正缺的不是 env，
  是 `voicePlan.AECEnabled` 这个**编译期零值**字段 —— 所以"查 env 齐不齐"这个方向查不到它。
  (3) **与票 392 的 AEC 调参语义无关，但同根**：392 三枚**也用 `newVoice26pRig`**，
  只因 `requireAEC3()` 先 Skip 把同一缺陷掩盖了；且**票面"本轮 Test392* 两枚 PASS"与本程现读的
  round-88 日志不符**（实测是 `--- SKIP`，原因是 `-tags webrtcap` 缺绑定），请编排者核。
- **修形：只改台面**（产码零字节）。`newVoice26pRig` 的 plan 补 `AECEnabled: true` ＋ 六行注释说明
  为什么手搓 plan 必须说出生产默认值（沿用票 423 对 `PunctEnabled` 的同一条裁定，
  `resident_voice_wire_r1_windows_test.go:286-288`）。**断言一个字节未动**（`:215-217` 原样）。
- **修复后 carrier 两发**：`M2-VERBATIM episodes=1 passed=0 triggers=0 transcript="…"`，
  26.94s / 24.37s ⇒ **`:216` 红句消失**。⚠ 但两发都是 **SKIP 而非 PASS**：
  测试继续走到**还回 AEC 的那半**（`:249-250` 的 `rig2.requireAEC3()`）Skip，
  本机 `third_party/webrtc-audio-processing` 源码树 MISSING ⇒ 本机证明的是 M2 前半
  （判定器活、200 帧进、开 episode、逐字复现），**restore 半本机不可测**（具名记残余）。
- **突变控制两发（全经 `go test -overlay`，拷贝在仓外 `D:\tmp\wisp\20261112\428-26p\mut\`）**：
  - **发①** 删掉修好的那行（`true`⇒`false`）⇒ **红**，逐字复现原红句，21.75s。
  - **发②** 恢复 ⇒ 红句消失、`episodes=1`、逐字转写命中，27.26s。
- **门禁**：`go build ./cmd/wisp/` rc=0 · `go vet ./cmd/wisp/` rc=0 · `gofumpt -l` 两尺零输出 ·
  `scripts/d22scan.sh` rc=0（`clean - no D22 ban violations`，`runtests.sh: OK PASS=45 FAIL=0 SKIP=0`）。
- **证据件**：`docs/evidence/s1/428-26p-echo-episodes.md`。
- **越界**：`git diff --name-only` 本腿只有 `cmd/wisp/resident_voice_bargein_26p_windows_test.go`
  与证据件；其余五行是他腿（426/427 收件改名、429 在飞、既存脏改）。禁面全未动，`build/mut_rvw.go` 未删。
- **★需回填（⛔ 本程未改票）**：修好的是**共享 rig**，26P 另外四枚与 392 三枚也一并解封，
  但它们在本机仍 SKIP ⇒ **其真实状态仍未被测**，需带 `-tags webrtcap` 的 staging 才能判。