# 票 409 — `internal/speech` 语音族测试被环境红：根因＝`C:\Windows\System32\onnxruntime.dll`（1.17.1 旧版）遮挡 PATH 上的新版（1.28.2）；修法＝测试二进制与三枚钉值 DLL 同目录的 staging carrier

**立票**：2026-11-08，Wisp 编队写腿（编排者派单「治 System32 旧 onnxruntime.dll 遮挡——它一口气弄红 30 枚语音测试」）
**性质**：环境遮蔽票（测试台面稳性），非代码缺陷。与票 98（`cmd/wisp` 缺 DLL 载入死）同族不同层：98 是一级导入找不到（PATH 前置可救），本票是一级导入找到了、其**二级依赖** `onnxruntime.dll` 被 System32 抢绑（PATH 前置救不了）。

## 根因（现测，逐件核过）

- `C:\Windows\System32\onnxruntime.dll`：1.17.260311-1434.1.os-germanium.6bcd20d，10,572,960 B（OS 自带 ML 组件）。
- `third_party/sherpa-onnx/onnxruntime.dll`：1.28.2，17,799,168 B（`deps.toml` 钉值，`fetch-deps.ps1` 落盘）。
- Windows 加载器标准搜索序：**exe 所在目录 → System32 → … → PATH（最后）**。`internal/speech` cgo 链接 `sherpa-onnx-c-api.dll`，其二级依赖 `onnxruntime.dll` 在**进程载入时**解析——System32 旧版赢。
- 实测（2026-11-08，HEAD `eec09a028`，裸 `go test ./internal/speech/ -count=1`）：**RUN=106／PASS=42／FAIL=0／SKIP=47**，rc=0 假绿；47 枚 skip 理由 100% 为 `checkOnnxruntime` 的遮蔽文案（`onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires…`）。整包口径下 `tools/d22scan/runtests.sh` 把 SKIP 记 fatal（票 71 AC#3）⇒ 整包 931s 那轮「语音族 30 红」的机制＝遮蔽→skip→计红，**不是断言红**。

## 修法选型（⓪⓪⓪：主流 harness 形状，理据写在证据件）

**选 staging：测试二进制与三枚钉值 DLL 同目录**。加载器搜索序第一步＝exe 所在目录，是唯一在载入时、System32 之前生效的杠杆；也是 Windows 标准部署形状（DLL 与 exe 同目录），与产线 `scripts/build.ps1`（`build\wisp.exe` + 三枚 DLL colocated）同构，且 `internal/speech/fixtures_test.go` 的 `sileroModelPath` 逐字预留了该形状（「run from a staging directory that holds the sherpa DLLs」）。被否两条：`TestMain`/`init` 显式 `LoadLibrary` 太晚（cgo 导入树在进程入口前解析完）；PATH 前置只救一级导入（本票基线 PATH 已前置仍 47 skip，实锤）。

## 落地（本腿）

- `scripts/stage-sherpa-dlls.sh`（新）：从 `deps.toml` `[sherpa-onnx.dll."…"]` 钉值读 DLL 名（与 `wisp-cli-tests.sh` 同一合同，不复制名单），GUARD 校验后拷到目标目录；setup 失败响红不静默。
- `scripts/speech-tests.sh`（新）：carrier——`go test -c` 产出测试二进制 → stage 三枚 DLL 到 exe 旁 → 包目录（`go test` 同款 cwd 合同）运行 → runtests 同形四数。`WISP_SPEECH_STAGE_DIR` 可改 stage 根，默认 `.scratch/wisp/stage-speech/`。
- `docs/evidence/s1/sherpa-dll-shadowing.md`：裁决表（根因/被否两条/修法/验收读数）。
- 零 Go 产品码、零断言改动、零 thresholds/golden 改动、未触碰 System32。

## AC

- [ ] AC#1 carrier 真跑：`bash scripts/speech-tests.sh` 到终态，遮蔽型 skip（文案含 onnxruntime 1.17.1）归零，模型背书测试以真跑（`--- PASS`）落地。
- [ ] AC#2 载入来源验证：`cgo_boundary_windows_test.go` 的 `EnumProcessModules` 断言以真跑通过（进程内 onnxruntime 来自 stage 目录新版）。
- [ ] AC#3 门禁：gofumpt 空／vet rc0／build rc0／d22scan clean／改动文件 LF。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-11-08 21:02:38 +0800] agent=Wisp 写腿(票409验收接手腿) did=上一腿交件未commit未跑完验收——本腿现核三件套(两枚DLL实测System32=1.17.1/10572960B vs third_party=1.28.2/17799168B)+自己复跑carrier到终态: RUN=106/PASS=83/FAIL=0/SKIP=6 rc=0 全程147s(修前42/0/47)——作差PASS+41/SKIP-41, 41枚遮蔽skip全转真跑PASS, log内1.17.1提及=0; 残留6 skip逐名=标点模型缺件(WISP_PUNC_MODEL, 非遮蔽); 载入验证TestSingleOnnxruntimeInstance PASS 0.26s(EnumProcessModules断言onnxruntime来自stage目录1.28.2); 门禁(clone --shared导出树75ff9c4eb)gofumpt改动件空(全树28枚均.scratch/probes+forks存量)/vet rc0/build rc0/d22scan clean rc=0; 如实=遮蔽是本机独有(CI只挂internal/speech/cer纯Go子包, cgo父包不在CI射程), 裸go test本机仍skip(加载器物理), 可复跑命令=bash scripts/speech-tests.sh; 实测读数已填docs/evidence/s1/sherpa-dll-shadowing.md §4; AC框/Status零触碰 next=编排者收件(只commit未push); 若要整包口径吃掉39红, 需整包runner把internal/speech路由过carrier——编排者拍
- [2026-11-08 21:58:30 +0800] agent=Wisp 写腿(speech-runner / wave99) did=整包runner路由落地: portable-tests.sh新增speech tier(speech_pin+tiers名册+case分支+census pin循环, GUARD C/A对speech照咬), run步mode=speech不走strict runner改路由scripts/speech-tests.sh(stdout经tee进同一capture, speech-tests.sh的-v流由>"log"改tee+PIPESTATUS), GUARD B用speech形(-v分母非零, 残留skip不当fatal=carrier exit code已携带断言判决); 导出树(clone --shared, 资产已拷)验收: --scope=speech⇒RUN=106/PASS=83/FAIL=0/SKIP=6 rc=0(与carrier单跑逐字同), 残留6 skip=标点模型缺件非遮蔽, 39假红→0; census⇒internal/speech行=UNCLAIMED-HAS-TESTS→CLAIMED BY speech(census自检rc0; GUARD D另剩4枚与speech无关存量洞: balldebug/agent/scheduler/command/plugin/manifestslot, 归后续票); portable-tests-selftest.sh首跑1断言红=unknown-scope-still-hard的known:钉过期, 补speech后重跑32case/0失败GREEN; 门禁(导出树)gofumpt scripts/命中0/vet rc0/build rc0/d22scan clean/三脚本CRLF行0; 证据件=docs/evidence/s1/speech-runner-routing.md; 写面=scripts/portable-tests.sh+portable-tests-selftest.sh+speech-tests.sh+证据件+本行, 零Go产码零断言零workflow; AC框/Status零触碰 next=编排者收件(只commit未push); CI侧 speech tier 未接线(主线scope认领走本腿本地档, ci.yml动=另判)
