# T14 对抗验收报告（编排者执行）

> 执行者：orchestrator（实现者 T14-impl 独立）。背景：截止警戒（配额 09:00 到期）下由编排者内联验收。时间：2026-09-19T23:09:31Z

| # | 项 | 裁决 | 证据 |
|---|---|---|---|
| 1 | 全量复跑 | PASS | `go test -count=1` + `-race` ./internal/models 全绿 |
| 2 | 签名/篡改/P3 | PASS | TestTamperedModelByteRejectedAndDeleted / TestTamperedManifestFailsBeforeAnyDownload / TestVerifySignatureFalseIsHardError / TestP3BlockedModelRefused 全 PASS |
| 3 | 契约字段 | PASS | 6 字段全在位；多工件条目级联摘要公式有测试钉死；files[]/archive/status 扩展有 doc+测试 |
| 4 | 越界 | PASS | 提交触碰 internal/models、models/(manifest+fixtures)、compose、.gitattributes、PRECHECK、票 14——fixtures 为 compose model-mirror 资产（合法）；internal/llm、tools/mockllm 零卷入 |
| 5 | 真下载抽验 | PASS（实现自证 + 哈希互证） | WISP_IT_REAL_MIRROR=1 两用例 PASS；KWS/matcha/SenseVoice 哈希与 spike 互证 |
| 6 | D22 | PASS | 无新白名单外依赖（x/crypto blake2b 已裁定合理）；无 emoji；dev 密钥盘外 |

## 发现与裁定
1. **P3 BLOCKED（matcha-zh-baker 非商用）**：登记 PRECHECK + reports。[裁定] **个人使用期可继续用 matcha（非商用许可允许）**；经 C29 manifest **分发**被挡（status=blocked-p3，Manager.Ensure 硬拒）。默认 TTS 替代选型归票 26（P7 音质门禁一并评测）。
2. dev minisign 密钥盘外存放（E:\workase\wisp-minisign\）；生产仪式归 S8（登记 H5）。
3. cmd/wisp doctor 把公钥标注 "placeholder" 的陈旧文案——一行修改，转票 12（gate 会重跑 doctor）。

**VERDICT: PASS**

## Addendum 裁决（2026-09-20，AC 补裁）

> 背景：上表 6 行只裁了 AC#1/AC#6（篡改/签名/P3/契约字段）+ 通用复跑，AC#2/AC#3/AC#4/AC#5
> 从未被裁（票据复核已记为四个空框）。本节由补裁代理实跑实读后补裁；四框均可勾，
> 但 AC#2/AC#4 各带一条明写的残余（不构成扣分项，供票 15/16 与真实链路验收接手）。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#2 续传：传输中途杀掉下载器 → 重启后继续（**字节偏移已验证**）、完成、哈希通过 | **PASS** | `internal/models/downloader_test.go:317 TestResumeContinuesAtByteOffset`。断言链完整且每一步都落在可观测量上：①`srv.setPartial(2_000)` 让服务端只吐 2000 字节、`Attempts:1` 使**首次 `Ensure` 必在中途失败**（`:331-334` 断 err≠nil，即真的"断在半路"而非假成功）；②`:335-342` 断磁盘上 `staging/vad-fixture/vad.onnx` **确实存活且 size==2000**（续传前提）；③`srv.setPartial(0)` 后二次 `Ensure`，`:346-348` 断**服务端实际收到的 HTTP 请求头** `LastRange()=="bytes=2000-"` ——这是"字节偏移已验证"的字面对象，不是内部状态自陈；④`:350-356` 回读安装文件比 `shaOf` 全等；⑤`:357-359` 断 staging 成功后被清。<br>`go test ./internal/models/ -run 'TestResumeContinuesAtByteOffset' -count=2 -v` |

```
=== RUN   TestResumeContinuesAtByteOffset
--- PASS: TestResumeContinuesAtByteOffset (0.02s)
=== RUN   TestResumeContinuesAtByteOffset
--- PASS: TestResumeContinuesAtByteOffset (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/models	0.068s
```
> 残余（不改裁决）：此处"kill"是**同进程内**的中途失败 + 新 `Ensure`，非进程真死再冷启。
> 续传所需状态全部来自磁盘暂存与 Range 头，无进程内残留依赖，故等价性成立；
> SPEC-10 §7 把"模型下载中断"列为**必须真实触发**的失败预演项，该真·冷启复跑归 16/预演阶段。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#3 镜像故障转移：主镜像 404 → 用备用；全程发进度事件；取消不在 staging 外留残件 | **PASS** | 转移+进度：`internal/models/downloader_test.go:369 TestMirrorFailoverAndProgressEvents` —— `primary.setNotFound(true)` 造真 404，`:381-386` 断 **`primary.Hits()!=0`（主镜像确实被试过）且 `fallback.Hits()!=0`（备用确实被用上）**（只断结果会漏掉"根本没试主镜像"的假转移），`:388-390` 断装成内容哈希正确，`:391-396` 断 `PhaseConnecting/Downloading/Verifying/Installing/Done` **五个阶段一个不缺**（缺一 Fatalf）。取消：`:399 TestCancelCleansStagingNoPartialsOutside` —— 用 `release` 通道把服务端卡在 4096 字节之后，`waitForPhase(PhaseDownloading)` 确保真在下文中才 `mgr.Cancel(...)`，`:427-432` 断取消的 `Ensure` 必须返错且错串含 "cancel"，`:434-446` **`filepath.Walk` 整个 DataDir**，任何路径含 `/staging/` 的文件即判错（连 staging 本身也必须清空）。<br>`go test ./internal/models/ -run 'TestMirrorFailoverAndProgressEvents\|TestCancelCleansStagingNoPartialsOutside' -count=2 -v` |

```
--- PASS: TestMirrorFailoverAndProgressEvents (0.02s)
--- PASS: TestCancelCleansStagingNoPartialsOutside (0.01s)
--- PASS: TestMirrorFailoverAndProgressEvents (0.02s)
--- PASS: TestCancelCleansStagingNoPartialsOutside (0.02s)
ok  	github.com/CarlosShao/wisp/internal/models	0.127s
```

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#4 `Downloading` 态在球上走一遍（进入/进度/退出），**经状态日志断言** | **PASS** | 断言对象正是 AC 点名的"状态日志"：`internal/models/bridge_test.go:23 TestDownloadingWalkSuccess` —— `:37-40` 取 `bridge.StateLog()` 并断长度恰为 2 且**序列全等** `[Downloading, FirstRun]`（对应 SPEC-08 §3 转移表 **行 #2 `FirstRun`+缺模型→`Downloading`** 与 **行 #37 `Downloading`+完成→`FirstRun`**），`:41-46` 另断 `Ticks()!=0`（进度确有节拍）与 `LastPercent()==100`（进度值真抵达终点），`:48-56` 断行 #37 的副作用 `model.verify-sha256-signature`（C29 校验）确实发射。退出到失败支由 `:60 TestDownloadingWalkFailure` 断 `[Downloading, Error]`（行 #37 失败出口）+ `machine.State()==Error`。**非法进入必须被拒而非强塞**由 `:77 TestDownloadingWalkRejectsIllegalEnter` 钉住（从 `Sleeping` 进入 `Run` 必须报错且 `machine.State()` 仍为 `Sleeping`，D43 #2 白名单不外溢）。桥本身是球侧状态机的正牌入口：`newWalkRig` 用 `WireDownloading(mgr, machine)`（`bridge_test.go:20`）挂上真 `statemachine.Machine`。<br>`go test ./internal/models/ -run 'TestDownloadingWalk' -count=2 -v` |

```
--- PASS: TestDownloadingWalkSuccess (0.02s)
--- PASS: TestDownloadingWalkFailure (0.01s)
--- PASS: TestDownloadingWalkRejectsIllegalEnter (0.01s)
--- PASS: TestDownloadingWalkSuccess (0.01s)
--- PASS: TestDownloadingWalkFailure (0.01s)
--- PASS: TestDownloadingWalkRejectsIllegalEnter (0.01s)
ok  	github.com/CarlosShao/wisp/internal/models	0.101s
```
> 残余（不改裁决）：本框验的是**状态机侧**的 Downloading 走态与进度值；球**窗口**上的
> 环形进度+百分比渲染（SPEC-08 §2.1 `Downloading` 行）只在 `internal/ball` 侧
> `VisualFor` 全覆盖测试与 `winlive` 门控实跑里出现，后者当前不可裁定
> （见 docs/evidence/s1/07-adversarial-acceptance.md §Addendum"发现 A"）。渲染的主观核对并入 H1。

| AC | 裁决 | 证据（file:line）+ 实跑命令 + 输出尾部 |
|---|---|---|
| AC#5 local_override：从本地路径加载模型并做完整性校验，**无网络调用** | **PASS** | `internal/models/downloader_test.go:451 TestLocalOverrideIntegrityNoNetwork`。"无网络调用"用的是**最硬的那类断言**：注入 `HTTPClient: &blocking`（`:460`），其 `RoundTrip` 为 `panicTransport`（`:486-488`，`panic("network call made during local_override (must never happen)")`）——任何一次出网都会直接把测试打崩，而非"事后检查计数器"。`:466-469` 断 `Ensure` 成功且 `:470-472` 断返回目录**就是**本地覆盖目录（未拷贝、未下载）。"integrity check 仍然咬人"另有一例：`:475-482` 把本地 `vad.onnx` 翻一个 bit 后**同一个 mgr 再 `Ensure` 必须失败**，即 override 不是校验旁路。与配置层 `verify_signature=false` 硬错（`:307 TestVerifySignatureFalseIsHardError`，已勾 AC#1）方向一致、互为双保险。<br>`go test ./internal/models/ -run 'TestLocalOverrideIntegrityNoNetwork' -count=2 -v` |

```
--- PASS: TestLocalOverrideIntegrityNoNetwork (0.02s)
--- PASS: TestLocalOverrideIntegrityNoNetwork (0.01s)
PASS
ok  	github.com/CarlosShao/wisp/internal/models	0.079s
```

**补裁小结**：AC#2、AC#3、AC#4、AC#5 均 **PASS**，票据四框改勾；本票 6 框自此全部有逐条书面裁决。
两框改勾时明写残余：AC#2 的真·冷启中断预演（SPEC-10 §7）与 AC#4 的球窗口 Downloading 渲染（并入 H1）。
