# 95 — 对抗验收裁决表（acceptor-ticket95 · 终版）

- 验收基线：HEAD `a367b79`（实现方 commits：`c0bdc48` 裁定 / `80923a9` 接线+钉子 / `33d86b3` 交件 / `a367b79` 补测；
  本代理 checkpoint1=`1b78f62`）。
- **全部读数的出处树**：`/tmp/wisp-ac95-gate`（`git archive a367b79 | tar -x`，下称〔gate 快照〕）与
  `/tmp/wisp-ac95-mut-agent-ticket95`（同 sha 的第二份快照，变异专用，结束后逐文件 diff 回 gate 快照证还原）。
  真实工作树此刻含邻居未提交 WIP（`internal/winsec/winsec_windows.go` 是 ` M`，票 106 的；`scripts/`、`internal/risk`、
  `internal/tools` 亦有）⇒ 本代理**从未**在真实树里 build/test/变异，每份绿都点名快照。仓根 `git diff --quiet` 不适用
  （别人的 M 不是我造的），改证 `git status --porcelain` 里**没有一条路径属于本代理**。〔独立复现〕
- 改动面独立核对〔独立复现〕：`git show --stat` 四枚 —— 三枚 docs-only；`80923a9` 只动
  `internal/config/migrate.go`、`internal/config/parse.go` + 三枚新测试。`internal/observe/`、`internal/ball/`、
  `cmd/wisp/`、`internal/winsec/`、`internal/secret/` 在本票四枚 commit 内一字未动。与自述一致。
- 档位标注：〔独立复现〕= 本代理亲敲；〔日志＋归档，我抽验〕；〔仅自述，不背书〕。

## AC#1 七类逐类裁定（票面原句：「上面七类逐类裁定"算不算私有数据"，每类给：里面实际会出现什么…泄露后果…封了之后谁会读不到。⚠ 不许整族打包贴同一个结论」）

| # | 落点 | 实现方裁定 | 本代理独立复核 | 档位 |
|---|---|---|---|---|
| ① | 模型下载 staging `downloader.go:224/:287` | 不封＋反向钉子 | **成立**。内容=镜像站公开字节；出门前 `:263-276` 逐文件全量门失败即 `RemoveAll`；泄露后果≈零（hash 钉在离线验签清单里，`doc.go:6-16`）；封的代价=挡掉别的账户/服务实例复用 GB 级缓存，且挡不住真正的问题（见 TOCTOU 段） | 〔独立复现〕 |
| ② | 解包模型文件 `archive.go:41/:73/:78` | 不封＋反向钉子 | **成立**。`ExtractTarBz2` 只物化 regular file/dir、link 整件拒收、成员路径过 `validRelPath`（`:54-67`）⇒ 解包树宽≠任意写入口。票面表②"这些文件之后被当作可信模型读回"**被本代理判为引用了错误前提**：读回侧每次 `Ensure` 都逐文件重验（见三问） | 〔独立复现〕 |
| ③ | 配置写路径 `migrate.go:83`、`parse.go:213` | 要封，已接线 2 处 | **接线属实**：现读 `migrate.go:93` `winsec.PrivateFile(backup, raw, 0o600)`；`parse.go:201-226` `CreateTemp`→`:215` `winsec.SealFile(tmpName)`→写字节→rename。票面 file:line 偏移一行级（83→93、213→215），**语义未腐坏** | 〔独立复现〕 |
| ④ | 凭据迁移（票面引 `secret/migrate.go:154/:168`） | 票面引用已腐坏；票 89 早已接；只剩 `MigratePlaintext` 生产零调用方 | **三点全部复算为真**：现读 `:169` `winsec.PrivateFile(backupPath, raw, 0o600)`、`:174` `winsec.SealFile(backupPath)`（既存备份修复腿）、`:189` `winsec.PrivateFile(tmpPath, out, 0o600)`；全仓 grep `MigratePlaintext` 非测试命中只有定义 `:36/:66/:86` 与两枚注释 ⇒ 生产零调用方属实。**判据本身有毛病的一处**：票面 17:4x② 要求"顺手把装配根一起接"，同张票面又把 `cmd/wisp/` 标为"别人的地界别抢"——两条互斥，实现方选择登记 `R-95-1` 不接，判**合规**（禁改区优先） | 〔独立复现〕 |
| ⑤ | 日志 `observe/logging.go:72/:244` | 交回 owner，不判 | **诚实留白**：复测两处仍是 `MkdirAll(0o755)` / `OpenFile(..., 0o644)`（grep 原文一致）；`schema.go:505` `RedactPaths default:"false"` 属实 ⇒ 默认日志带绝对路径的取证成立。有书面接手方（owner）+ 两路代价段 + 施工约束（不得 `SealDir` 传播）。合法转移，不是破口 | 〔独立复现〕 |
| ⑥ | 球位置 `ball/position.go:73/:77` | 只登记 | 实读 `Put()`：内容 `json.MarshalIndent(显示器设备名→PosEntry)` ⇒ 泄露面只有坐标；`:77` tmp+rename 与配置同形状（rename 保留描述符）⇒ 若接是一行。接手方=票 64/65/68 地界，票面明写"别抢" | 〔独立复现〕 |
| ⑦ | doctor/SLO `doctor.go:248`、`slo_windows.go:559` | 只登记 | 实读：doctor 那处是 `probeWritable` 的 `MkdirAll` + `:252` 写 `"probe"` 探针（可写性测试，不是数据）；SLO 那处 body=`"ready=1\npid=%d\n"`（`:558`）⇒ 泄露面=一个 PID。两类都不是私有数据；接手方=cmd/wisp 主人 | 〔独立复现〕 |

**七类没有整族打包**：实现方给出 3 种不同判决（封/不封+钉子/交回+登记），本代理逐格复核方向全部一致。
**对实现方推翻台账那条（R-95-4）的独立复算**：`winsec` 封的是当前用户+SYSTEM+Administrators
（`winsec.go:16-19` 文档 + 下方探针 icacls 实测），同用户其它进程照读 ⇒ 台账"封了反而挡住多实例复用"
作为**代价**陈述**不成立**，实现方的更正**对**。〔独立复现〕

### 模型三问（本票最重一格，全文见上表①②，此处给三问直答）

- **问① 验签是不是每个文件都过**：**是**。`VerifyDir`（`downloader.go:519-531`）遍历 `InstalledFiles()`
  （`manifest.go:218-231`：归档件逐成员 `a.Archive.Files`，非归档件逐 `a.Path`），逐条
  `verifyFileHash`＝size 双钉 + 全文件流式 sha256（`:533-557`）。
- **问② 读取时还是只写入时一次**：**读取时**——但准确说法是"**每次交还时**"：`Ensure` 缓存命中 `:170-174` 与
  `local_override` `:161` 都在返回目录前过 `VerifyDir`；全仓非测试调用方只有 `bridge.go:44`。
  写入侧另有 `:263-276` 出门门。⇒ 不存在"只验一次、之后永远信"。
- **问③ 有没有路让未过签名的文件被当可信模型读回（TOCTOU）**：在 `models` 包内**没有**——任何交还都重验。
  **残余窗口是真的**：`Ensure` 返回后目录交给引擎（`doc.go:28` 明言引擎不在本包），**那一刻之后**被换/改的文件
  本包不会再查；且本代理探针实测继承树里外来 SID 拿的是 `(I)(M,DC)`＝**可写**。⇒ "不封"对票面问题
  （保密性：别的账户读到什么）**成立**；对**完整性**存在一条跨账户写窗，封目录可以关它，但那改变多账户复用形态，
  属引擎/语音票的决定 ⇒ 登记 `AC95-R1`，**不推翻本格**，判"通过但有条件"（条件＝钉子必须继续钉住两条前提，M3/M4 已证有牙）。

## AC#2 接线 + 真实读回的 icacls（票面原句：「对判"要封"的每一处：接线 + 一条真实读回的 icacls 证据（照票 89 的形状：主体→SID 白名单），并在同一枚 commit 里给"没接之前它确实宽"的基线读数」）

本代理**自己在仓外临时目录**复现（`go run` 探针，只碰 `%TEMP%` 下自造目录，绝不碰真实数据目录）：

- 父目录 seed `BUILTIN\Users:(OI)(CI)(RX)` 落地 rc=0；另 seed 外来 SID `*S-1-5-21-3623186960-731165060-4091685855-1717338598`
  rc=1332（本机无映射，**但无需 seed 也宽**——该外来 SID 与 `CodexSandboxUsers:(I)(M,DC)` 本来就是这台机器
  `%TEMP%` 祖先链上的**继承授权**，与实现方报的同源同形）。
- **前**（`os.WriteFile(backup, raw, 0o600)`，接前形状）：`config.toml.bak-1` icacls 原文六条——
  `BUILTIN\Users:(I)(RX)`、`DESKTOP-LVS7839\CodexSandboxUsers:(I)(M,DC)`、`S-1-5-21-…-1717338598:(I)(M,DC)`、
  `NT AUTHORITY\SYSTEM:(I)(F)`、`BUILTIN\Administrators:(I)(F)`、`DESKTOP-LVS7839\swq:(I)(F)`
  ⇒ **`0o600` 一个字节都没落地**，同机别的账户可读可改。
- **后**（`winsec.PrivateFile`）：`sealed.bak` icacls 只剩三条——`NT AUTHORITY\SYSTEM:(F)`、
  `BUILTIN\Administrators:(F)`、`DESKTOP-LVS7839\swq:(F)`，**无任何 `(I)` 继承、无 Users、无外来 SID**。
- **它没说的那件事——封完本进程自己还读得到吗**：三条路全部实测**读得到**——
  (a) seal 后同进程 `os.ReadFile` OK(19B)；(b) seal 后**另起进程** `cmd /c type` OK（＝同用户 tail 的形状）；
  (c) **rename 保留描述符的可疑路**：`CreateTemp`→`SealFile(tmpName)`（空文件先封）→**经已开描述符写入**→
  close→rename→`config.toml` icacls 仍是三条 `(F)` 无继承，`os.ReadFile` OK + 另进程 `type` OK。
  ⇒ "封了挡住自己"不成立；`winsec` 授的就是当前用户 SID。〔独立复现〕
- 基线读数与接线同 commit（`80923a9`）属实：`private_acl_windows_test.go` 的
  `TestAC2BaselineModeIsDecorative` 与两枚接线在同一枚 `--stat` 里。票 89 形状（主体→白名单断言 + t.Logf 原文）
  在用例里齐备。
- **结论：通过。** 位置提醒：票面 `parse.go:200/:213` 与 `migrate.go:83` 行号已漂移一行级（现 `:201/:215`、`:93`），
  引用形状对、坐标旧。〔独立复现〕

## AC#3 反向钉子有牙（票面原句：「判"不封"的每一处要有反向钉子…『没接』与『故意不接』在代码上要能区分」）

- 两枚都是**真用例**（非注释）：`no_seal_ruling_test.go:30-79` 安装后**逐成员翻一字节**⇒`VerifyDir` 必须报错；
  `no_seal_ruling_windows_test.go:28-66` 跑**生产函数** `ExtractTarBz2`，断言每个成员**必须仍带** `BUILTIN\Users`
  继承授权，并顺手验证内容可读且 hash 正确。
- **牙的变异实证（临时探针，仓外快照，非仓库内用例复现）**：
  - **M3**：把 `VerifyDir` 的覆盖面缩到 `[:1]`（build rc=0）⇒ `--- FAIL: TestAC3EveryInstalledFileIsReverifiedAtReadTime`，
    断言原文 `tampering installed file 2/4 (tokens.txt) went undetected: … ticket 95's no-seal ruling for this class is void and it must be re-ruled`。
  - **M4**：给 `ExtractTarBz2` 的每个解包文件加 `winsec.SealFile`（build rc=0）⇒
    `--- FAIL: TestAC3ExtractionIsDeliberatelyNotSealed`，原文 `model.onnx stopped being inherit-wide: something sealed this class. Re-open the ticket 95 AC#1 ruling`。
  - 还原：两发之后 `diff -r internal/`（mut 快照 vs gate 快照）**空**。〔独立复现〕
- "没接 vs 故意不接"可区分：⑤⑥⑦ 未钉是因为**未判/别人地界**（票面与本裁决表都点名归属），不是混同。
- **结论：通过。**

## AC#4 变异（票面原句：「抽一处接线退回 os.WriteFile(0o600) ⇒ 它的 icacls 判据必须红…锚点=承载行为那一行，同链 grep 证落地，还原后 git diff --quiet 证干净；编译失败不算变异」）

全部在 `/tmp/wisp-ac95-mut-agent-ticket95`（仓外快照，A38④）做；**每发先 `go build` rc=0，同链 `grep -n` 打印被改后整行**；
仓库内从未改过 ⇒ 真实树 `git status --porcelain` 无我产生的源码条目（别人的 ` M` 一字未动，见开头）。

- **M1（migrate.go 退回旧形状）**：import 行 `:10` 退回 `"os"`、`:93` 退回
  `if err := os.WriteFile(backup, raw, 0o600); err != nil {`（grep 原文打印）；BUILD_RC=0；
  `go test -v -run TestAC2`：RUN=3 PASS=2 **FAIL=1**——
  `--- FAIL: TestAC2MigrationBackupLandsPrivate`，原文 `config.toml.bak-1 is not private, these principals hold grants: [BUILTIN\Users DESKTOP-LVS7839\CodexSandboxUsers S-1-5-21-3623186960-731165060-4091685855-1717338598]`。
- **M2（parse.go 那条腿单独退回）**：`SealFile` 块退回 `os.Chmod(tmpName, 0o600)`；BUILD_RC=0；
  结果 RUN=3 PASS=1 **FAIL=2**：`TestAC2SaveFileLandsPrivate`（`private_acl_windows_test.go:130`）**红**
  ＋ `TestAC2MigrationBackupLandsPrivate`（`:90`，它同时断言最终 config.toml 私有）**连带红**。
  基线用例 `TestAC2BaselineModeIsDecorative` 两发都绿（**它本来就断言宽，绿=噪声里的对照组，不算红在守卫**）。
- **裁定"只有一处红"**：实现方 M1 那轮另一条腿绿**不是那条腿没锚定**——是分层正确（M1 没碰它的承载行）。
  本代理补做 M2 证明：退回 parse.go 腿 ⇒ 它自己的用例立刻红。**两腿各有锚**。红全部红在守卫
  （icacls 判据 `assertPrivate`），无编译失败充数、无加载噪声（M2 第一发因注释符笔误编译失败，**作废未计**，
  重发后才是正式读数）。
- **结论：通过。**（价值声明：以上红名都是**临时探针复现**；仓库内用例在 HEAD 上全绿，见 AC#5。）

## AC#5 门禁（票面原句：「按包 scope + d22scan 纯净快照 rc=0…gofmt -l/gofumpt -l 空、go vet rc=0、go test -count=2 各包 rc=0 并逐条点名 SKIP/FAIL」）

`go test -v -count=2` 于〔gate 快照〕，四包 **rc=0**，四数逐包点名（全部 `-v` 数出，`=== RUN` 含子用例，
TOP=top-level 无缩进行，SUB=`=== RUN` 带 `/` 的子用例行）：

| 包 | rc | `=== RUN`(全) | PASS(全) | FAIL | SKIP | 其中 TOP / SUB |
|---|---|---|---|---|---|---|
| `./internal/config/` | 0 | 200 | 200 | 0 | 0 | 112 / 88 |
| `./internal/models/` | 0 | 56 | 52 | 0 | **4** | 52+4 / 0（无子用例） |
| `./internal/secret/` | 0 | 58 | 58 | 0 | 0 | 42 / 16 |
| `./internal/observe/` | 0 | 94 | 94 | 0 | 0 | 94 / 0 |

- 与实现方四数**完全一致**。SKIP 逐条点名（2 条唯一 × `-count=2` 两遍＝4，均 `-v`）：
  `TestRealDownloadVadThroughPipeline`、`TestRealDownloadPuncArchiveThroughPipeline`，
  跳过原因原文 `manifest_real_test.go:147: real-network spot check; set WISP_IT_REAL_MIRROR=1 to run` ⇒ 环境闸，非回归。
- **models 计数"不一致"的裁定：判据本身有毛病。** 「PASS+SKIP(4) > RUN(56)」算术不成立——52+4=56 **恰好等于** RUN，
  自洽。真实易混点在**两种 grep 口径**：`^--- PASS`（只数 top）对 config 会得 112 而非 200（88 条子用例 PASS 行带缩进）；
  models 无子用例所以两口径重合。实现方报的是"含子用例"口径且四包口径统一 ⇒ 无缺陷，登记为读数方法差异。
- `gofmt -l`（四包目录）**空**。`gofumpt`：本代理 PATH 查无 ⇒ **如实登记"未跑"**（台账 MEMORY 提示它在 GOPATH/bin，
  本会话 `command -v gofumpt` → NOT FOUND，未越权去找）。
- `go vet ./internal/config/ ./internal/models/ ./internal/secret/ ./internal/observe/` **rc=0**；
  `GOOS=linux go vet` 同四包 **rc=0**（台账坑"整仓 `GOOS=linux go vet ./...` 恒红"在按包 scope 下不触发）。
- `sh scripts/d22scan.sh`（与 CI 逐字同形，gate 快照）：**rc=0**，`clean - no D22 ban violations`，
  台账八行 `bans #1-5 internal/=197 cmd/=20`、`ban #6 frontend/=37`、`ban #7 internal/tools/=17`、
  `ban #8 design/=16 frontend/=37 internal/=345 cmd/=26` ⇒ 与实现方逐字一致；`allowlist.txt` 非注释 **5 行**未变短未变长。
- `go test ./cmd/wisp/` 本机既有红（缺 dll，加载期 `0xc0000135`）：**未跑、未追**（票面 AC#5 原话授权"不要去追"）。〔独立复现〕
- **"SealFile 挡不挡同用户 tail / `wisp doctor`"实测答复（拿给 owner 的那句）**：不挡——同用户**另起进程**
  `cmd /c type` 密封文件成功（tail 形状），rename 路线密封后同用户重开读回成功；`wisp doctor` 本体因本机既有红
  未跑，但它的读与 tail 同形（同用户句柄）。**被关掉的只有别的非管理员本地账户 + 继承来的外来 SID。**〔独立复现〕
- **结论：通过**（gofumpt 未跑一项如实登记）。

## 总判

**AC#1 通过但有条件；AC#2/AC#3/AC#4/AC#5 通过。票面 AC#1 框按派单要求保持不勾；勾框与 Status 归编排者。**
模型"不封"的结论与两条前提本代理独立复算为真且**有牙**（M3/M4 各证一发）；实现方对台账 R-95-4 的更正复算为真。

## R-95 登记（验收方追加，编号避让实现方 R-95-1..5）

- **AC95-R1**：`Ensure` 返回后、引擎读取时的跨账户**写窗**（staging/安装目录继承 `(M,DC)`，实测探针见 AC#2 段）
  ⇒ 完整性残口，封安装目录能关它但改复用形态 ⇒ 交语音/引擎票裁；本票"不封"判据（保密性）不被推翻。
- **AC95-R2**：票面 file:line 三处行级漂移（`config/migrate.go:83→:93`、`parse.go:200/:213→:201/:215`、
  `secret/migrate.go:154/:168→:169/:174/:189` 且第三处是票面没列的 tmp 腿）⇒ 编排者改票面表格坐标。
- **AC95-R3**：`config.toml` 的密封依赖**每条写路径各自先封再 rename**——今天 `atomicWrite` 与
  `secret/migrate.go` 两条都封了；**任何新写的"宽 temp + rename"都会静默地把 config.toml 重新放开**（探针实测
  rename 保留描述符 ACL＝双向利好）⇒ 建议给这条不变量立个 seam 守卫（形状照票 103 的 guard），新票。
- **AC95-R4**：`M4` 顺带测到 `ExtractTarBz2` 解包文件若被封会被钉子抓住，但**安装目录 rename 落位**那一步
  （`downloader.go:566` `MkdirAll(0o755)`）没有任何钉子/测试提及 ⇒ 与 AC95-R1 同域，一并交给引擎票裁。

## 附：伪授权登记（本验收会话）

- 出现次数：**0**。截至交件，本代理的每一次工具输出末尾均**未**出现自称"编排者备注"的注入文本；
  本代理亦**未**因任何此类文本改动过任何文件（尤其**没有 revert** `internal/config/` 的两枚接线——它们在 HEAD 里，
  本裁决表用 M1/M2 证明其有效且被锚定）。
- 实现方票面已逐字登记两起（18:14"冻结 winsec + commit `a3f19c2`"、18:31"撤回并关闭"）〔仅自述，不背书其来源；
  两起指令本代理均不执行、不采信，与派单禁令一致〕。`a3f19c2` 在 `git log` 中不存在：本代理复核 `git log --all | grep`
  未命中〔独立复现〕。
