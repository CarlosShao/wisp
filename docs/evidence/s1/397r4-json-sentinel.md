# 票 397-r2 残余②执行腿 `397-r4` —— 392 仪器面 JSON 安全哨兵（写腿）

日期 2026-11-01 · 腿 `397-r4`（写；写面＝只 `cmd/wisp/resident_voice_bargein_392_windows_test.go` ＋本件）·
起手 HEAD `789c6c6fbe9cd77f65b10398f705d2d5f09444ef`（dev，开工三核过：①锚点在（HEAD＝锚点同枚）②索引净（暂存 0 条、脏区全 `.scratch/**` 一类未跟踪件）③写面无他人脏）·
staging＝A1326 全套（staged 二进制＋三 DLL 贴 binary＋WISP_* 五 env＋manifest 可达＋cwd=包目录＋`unset GOTMPDIR`＋Git Bash runner＋`-test.timeout 30m`）＋**音频错峰**（同行 `speech-base398.test.exe` rig 在飞 ⇒ 等待环轮询至清空，无一发带争用强跑；本腿自身跑批逐发 rig_check 串行）·
工作树全程只写写面两件、突变全在仓外导出树＋overlay、⛔ push/amend/reset/rebase/stash/checkout。
注：派单存档 `.scratch/wisp/probes/orch/20261101/txt/397-r4-dispatch.txt` 盘上不存在（该目录现只到 `397-v2-dispatch.txt`）；本腿按任务面派单全文执行，此缺口如实记。

## 0. 一句话

**JSON 安全哨兵落地、(b) 撞面消除、反向有牙（撞面错误逐字复现）、零回归与门全绿。** `write392JSON` 序列化前经 `voice392JSONSafe` 反射走查，非有限值钳成有限哨兵 `voice392NoSignalDbfs = -999`（dBFS）。单元驱动案（新 `Test392JSONSentinelClampsNonFiniteReadings`，仪器面）rc=0：finite 读数逐位同（byte-equality 断言内嵌）、gate-clean (b) 形（episodes=0 的 rms `-Inf`＋纯静音桶 `-Inf`）→ 工件合法落盘 `rms:-999`；真 rig keepalive 复活案 rc=0（35.22s，WARM REVIVAL MOVED 断言全过）工件全新落盘、正常形 130 键路径新旧逐位同（本跑 (b)=episodes=1 有限值形 ⇒ JSON-SENTINEL 静默＝正常形哨兵零干预活证）；摘哨兵突变（仓外 overlay，writer 整块还原 raw marshal、diff 恰 4 删+1 改）⇒ rc=1 `marshal 392-json-sentinel.json: json: unsupported value: -Inf`（397-v2 撞面错误逐字）、还原核 SHA 逐位同。零回归：392 三案＋26p keyed 四案＋单元案还原核全 rc=0、六族 34/1/0 逐位归位（Test384 单发绿归因）、四数 98/75/0/11 逐位同基线、d22scan 实扫＋导出树 clean rc=0、正控 35/0/0。

## 1. 病灶与修法（选型拍板）

- **病灶（397-v2 §6② 现读复核）**：非有限值进工件面有两处——①`voice392Capture` 无 episode 时 `rms` 保持 `math.Inf(-1)` 哨兵（本腿现读测试件 :214）；②`voice392RmsProfile` 纯静音桶 `dbfs = math.Inf(-1)`（:110）。397-r2 壳修复把 rounds 首次接上 JSON 后，(b) 轮 episodes=0（echo 完全门净＝比常形更干净的健康形）任一字段为 -Inf ⇒ `write392JSON` 的 `json.MarshalIndent` 整件失败 `json: unsupported value: -Inf` ⇒ 工件不落盘、测试尾部红（行为断言全过）。
- **修法（397-v2 建议的 JSON 安全哨兵小改）**：`write392JSON` 在 `json.MarshalIndent` 前调 `voice392JSONSafe(v)`——reflect 走查 map/slice/pointer/interface/导出 struct 字段，凡 `±Inf/NaN` 浮点钳成 `voice392NoSignalDbfs = -999.0`；钳制发生时打一行 `JSON-SENTINEL <name>: clamped N ...` 测试日志（**不进工件、不加键**）。钳制点在 writer＝唯一咽喉：两件工件（keepalive/aec-tuning）＋map 里的原始 control readings 全覆盖，未来任何新字段的非有限值也被同一道闸拦住。
- **为何钳成有限值 -999、不是 null（拍板理由）**：①JSON（RFC 8259）没有 Infinity/NaN，Go 编码器遇一个非有限值整件失败——writer 必须是总函数；②**null 要改 schema**：`voice392RoundOut` 的 float 字段须改 `any`/指针才能表达 null，且若按派单另一选项「null＋哨兵计数」在工件加常驻计数键，就破坏「有限值字段逐位同」性质（正常形工件会多出一个键）——**有限钳制保 schema 零动、正常形逐位同**（本腿用单元 byte-equality 断言钉住）；③**语义诚实**：本面 -Inf 的含义就是「无信号」（无 episode 的数字静音/纯静音桶），-999 dBFS 远低于全套真读数带（最低 ≈-59），读作「无信号」而不是一个测量值，可排序、可被证据表引用；④NaN/±Inf 都不可能从两处具名来源产生（√(平方和)≥0 ⇒ 只有 0→-Inf 一路；episode RmsDbfs 同源），哨兵作为 writer 最后一道防线全捕——同值钳制是因为工件的职责是声明「该读数不存在」而不是猜一个量级，钳成 0 反会被误读成满量程。
- **⛔ 零触碰核（`git diff -U2` 现核）**：+150/−1 全在 reflect import 1 行＋write392JSON 注释与接线＋哨兵三件（const/`voice392JSONSafe`/`voice392JSONSafeWalk`）＋新单元案 `Test392JSONSentinelClampsNonFiniteReadings`；三个既有 Test 函数与全部 Fatalf 断言零 hunk（FINDING/WARM REVIVAL MOVED 语义零字节）；rounds 字段语义零字节（`voice392RoundOut` json tag 逐字未动、`392r1-session-integration.md` 口径零触碰）；thresholds/golden/`internal/**`/D43 表/statemachine/testdata 仓根/frontend/design/票面/台账/specs/PLAN 零触碰。

## 2. 读数（staged A1326 全套）

1. **单元驱动案（仪器面，确定性）**：`Test392JSONSentinelClampsNonFiniteReadings` rc=0——finite 夹具钳制计数=0 且 `json.Marshal(finite)` 与经哨兵后 marshal **逐字节同**（断言内嵌）；gate-clean (b) 夹具（`rms: math.Inf(-1)`＝:214 摘形、`res_profile` 首桶 `math.Inf(-1)`＝:110 摘形、episodes=0/triggers=0/built=true/ready=9000）钳制计数恰 2、marshal 成功、回读 `rms=-999`＋静音桶 `-999`（有限桶 `-58.9` 原样）；端到端走真 `write392JSON` ⇒ `392-json-sentinel.json` 落盘（内容现读：`"rms": -999`、`"res_profile": [-999, -58.9]`、episodes=0）。
2. **真 rig keepalive 复活案**：rc=0（35.22s 七窗全序）——WARM REVIVAL MOVED 断言全过（warm2=episodes=1/triggers=1 形）、`392-keepalive.json` 19:09 全新落盘（写盘前 18:24 旧件在案对照，397-v2 残余⑤「失败发不写盘」形状）。**本跑 (b) 轮＝episodes=1（有限值形）⇒ 日志零 JSON-SENTINEL 行＝正常形哨兵零干预的活证**；(b) episodes=0 形在真 rig 的出现率属材料时序摆动族（397-v2 §6① 同族归因），故本腿判据走确定性单元驱动案（同一 writer、同一 round-out 形状、同一 :214/:110 哨兵值）——真 rig episodes=0 形一旦出现即被同一道闸拦住（JSON-SENTINEL 行在场）。
3. **正常形工件 diff 检查**：新（19:09 本腿）vs 旧（18:24 换哨兵前 HEAD 代码）`392-keepalive.json` 结构 diff＝**130 个键路径集合逐位同**（差异只在跨跑数值＝TTS 材料微差在册口径）；有限值逐位同由单元 byte-equality 断言承载。
4. **还原核（本腿自己的尾读）**：全部实验后单元案复跑 rc=0（sentinel-restore）；写面 SHA 全程 `26b87d23f0a5156be9d1d594792f3151711ed3bb7a56a78af6b83608a4a0f5b0` 未动。

## 3. 反向判据（摘哨兵 ⇒ 撞面复现）

- 突变件＝仓外导出树（`git archive HEAD`＝`789c6c6fb…`，叠加本腿测试件后 SHA `26b87d23…` 与工作树逐位同）副本上的**整块摘除**：writer 的 `safe, clamped := voice392JSONSafe(v)`＋JSON-SENTINEL 日志四行删去、`json.MarshalIndent(safe,…)` 还原 `json.MarshalIndent(v,…)`——**diff 恰 4 删+1 改现核**（首版单行 sed 突变被本腿自抓＝`declared and not used: safe` 编译红，弃用重做，如实记）；构建走 `-overlay`（`mut/overlay-sentinel-removed.json`），二进制不落工作树。
- 读数：突变二进制（构建 rc=0）跑同一单元案 ⇒ **rc=1，`resident_voice_bargein_392_windows_test.go:229: marshal 392-json-sentinel.json: json: unsupported value: -Inf`——397-v2 报的撞面错误逐字**。摘哨兵即复现、装哨兵即消除，闸门有牙。
- 还原核：突变全程 overlay-only；跑后导出树 pristine 副本与工作树 SHA 逐位同（`26b87d23…` 两枚）；工作树零突变。

## 4. 零回归

- **392 finding 族全案**（staged tagged 二进制，逐发 rig_check 错峰）：`Test392KeepalivePostureFindingsAndColdDefault` rc=0（35.22s，§2）；`Test392DisposeReleasesTheEngines` rc=0（PASS=1，EnginesHot 翻转读数在场）；`Test392AecTuningDoubleTalkSwing` rc=0（PASS=1，4 设置×4 发＋还原核全序）。WARM REVIVAL MOVED / FINDING 语义红句零出现＝零行为漂移。
- **26p keyed 四案**（tagged）：`Test26PBargeinStopsWithin400ms`＋`Test26PEchoNeverTriggersAndNeverTranscribes`＋`Test26PWithoutAECTheEchoTranscribesVerbatim`＋`Test26PFeedOnlineMatchesWindowReplay` rc=0，**PASS=4 FAIL=0**。
- **六族 staged 合跑**（no-tag 二进制＋三 DLL＋五 env＋cwd=包目录）：`-test.run 'Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384'` ⇒ **PASS=34 FAIL=1 SKIP=0 逐位归位**（392-r1/397-v2 基线 34/1/0）；唯一红 `Test384ResidentWatchdogEndToEnd`＝合跑负载时窗既有形状，**本腿单发复跑 rc=0（1.24s PASS）归因核**；`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同基线**（六族日志在案）。
- **四数**（无 tag、`./internal/audio/... -count=1 -v`、静窗）：**rc=0，RUN=98/PASS=75/FAIL=0/SKIP=11 顶层口径与 397-v2 基线逐位同**（`internal/audio ok 17.732s`＋`aecspike ok 0.309s`）。⚠ 本腿首跑 70/57/0/1 rc=1＝**本腿 staging 缺口如实记**：aecspike 测试二进制旁无 onnxruntime.dll ⇒ `exit status 0xc0000135`（DLL_NOT_FOUND）；按 397v2 处方（quiet-batch.sh:71）`PATH="$REPO/third_party/sherpa-onnx:$PATH"` 前缀＋净 WISP env 补跑后逐位归位。

## 5. 门与尺

- gofumpt `-l`（`D:\work\base\gopath\bin\gofumpt.exe`，cmd/wisp＋internal/audio）＝空（rc=0）；vet 两形（无 tag／`-tags webrtcap`）rc=0；build 两形 rc=0。
- d22scan（Git Bash `scripts/d22scan.sh`）：**正控（种子违规）PASS=35 FAIL=0 SKIP=0（RUN=77）与 392-r1 基线逐位同**；实扫 **rc=0 clean**（bans #1-5 internal/=288、cmd/=52；ban#8 cmd/=184、internal/=734、frontend/=116、design/=39——分母随树漂移如实记）。
- d22scan 导出树（`tree/tools/d22scan` `go run . -root tree`）：**rc=0 clean**（同分母族；design/=14＝gitignored 面导出差，397-v2 §6⑤ 同形伪影）。
- 还原核尾读：跑后工作树写面 SHA `26b87d23f0a5156be9d1d594792f3151711ed3bb7a56a78af6b83608a4a0f5b0` 与导出树 pristine 副本逐位同；`git status --porcelain -- cmd/wisp/resident_voice_bargein_392_windows_test.go` 仅本腿一件。

## 6. 写面申报与边界

- 仓内（commit pathspec 逐枚）：`cmd/wisp/resident_voice_bargein_392_windows_test.go`（+150/−1，构成见 §1 零触碰核）＋本件 `docs/evidence/s1/397r4-json-sentinel.md`。
- 仓外 `D:\tmp\wisp\20261101\397r4\**`（bin/cases/logs/mut/out/tree/tree.tar/build.sh/run-staged.sh/run-resume.sh/fixup.sh/evidence-draft.md，只建不删）。
- ⛔ 越界核：行为断言/FINDING MOVED/WARM REVIVAL MOVED 语义零字节；rounds 字段语义零字节；thresholds/golden 零触碰；`internal/**`/D43 表/statemachine 零触碰；testdata 仓根/frontend/design（机主外部 agent）/票面/台账/specs/PLAN 零触碰；`voice392Out` 常量一字未动（工件仍落 392-r1 的既定目录，快照另存本腿 out/）；⛔ push（推送归编排者）。
- 未定义即停三条款全未触发：-Inf 来源在 writer 可钳范围内（:214/:110 两处，§1）；哨兵无须动行为断言即消除撞面；(b) 形经单元驱动案在本地确定性复现（§2①/§3）。
