# 95 — 对抗验收裁决表（acceptor-ticket95）

- 验收基线：HEAD `a367b79`（实现方 commits：`c0bdc48` 裁定 / `80923a9` 接线+钉子 / `33d86b3` 交件 / `a367b79` 补测）。
- 改动面独立核对〔独立复现〕：`git show --stat` 四枚 commit —— `c0bdc48`/`33d86b3`/`a367b79` 只动票面文件；
  `80923a9` 只动 `internal/config/migrate.go`、`internal/config/parse.go` 与三枚新测试
  （`internal/config/private_acl_windows_test.go`、`internal/models/no_seal_ruling_test.go`、
  `internal/models/no_seal_ruling_windows_test.go`）。`internal/observe/`、`internal/ball/`、`cmd/wisp/`、
  `internal/winsec/`、`internal/secret/` 在本票四枚 commit 内一字未动。与自述一致。
- 纯净快照：`/tmp/wisp-ac95-gate`（`git archive a367b79 | tar -x`），门禁读数出处。
- 档位标注：〔独立复现〕= 本代理亲敲；〔日志＋归档，我抽验〕；〔仅自述，不背书〕。

## AC#1 七类逐类裁定（票面原句：「上面七类逐类裁定"算不算私有数据"…不许整族打包」）

| # | 落点 | 实现方裁定 | 本代理独立复核 | 档位 |
|---|---|---|---|---|
| ① | 模型下载 staging（`downloader.go:224/:287`、安装目录 `:566`） | 不封＋反向钉子 | 不封成立（内容公开，清单离线验签 `doc.go:6-16`，出门全量门 `:263-276`）；⚠残留 TOCTOU 见"模型三问" | [待填] |
| ② | 解包模型文件（`archive.go:41/:73/:78`） | 不封＋反向钉子 | `ExtractTarBz2` 只物化 regular file/dir、拒 link、验 rel path（`archive.go:54-67`） | [待填] |
| ③ | 配置写路径（`migrate.go:83`、`parse.go:200/:213`） | 要封，已接线 | 实读：`migrate.go:93` `winsec.PrivateFile(backup, raw, 0o600)`；`parse.go:201-218` `CreateTemp`→`winsec.SealFile(tmpName)`→写→rename | [待填] |
| ④ | 凭据迁移（票面引 `secret/migrate.go:154/:168`） | 票面引用腐坏，票 89 早已接；`MigratePlaintext` 生产零调用方 | grep 复核：`MigratePlaintext` 生产代码零调用方属实（只有定义 `:66/:86` 与测试/注释） | [待填] |
| ⑤ | 日志（`observe/logging.go:72/:244`） | 交回 owner，不判 | `git show --stat` 证实 `internal/observe/` 一字未动；留白有书面接手方（owner）与两路代价段 ⇒ 合法转移 | [待填] |
| ⑥ | 球位置（`ball/position.go:73/:77`） | 只登记（别人地界） | 同上，`internal/ball/` 零改动 | [待填] |
| ⑦ | doctor/SLO（`cmd/wisp/doctor.go:248`、`slo_windows.go:559`） | 只登记 | `cmd/wisp/` 零改动 | [待填] |

### 模型"不封"三问（本票最重一格）——独立读码结论〔独立复现〕

- **问① 验签是不是每个文件都过**：是。`VerifyDir`（`downloader.go:519-531`）遍历 `InstalledFiles()`
  （`manifest.go:218-231`：归档件走 `a.Archive.Files` 逐成员，非归档件走 `a.Path/SHA256/SizeBytes`），
  逐条 `verifyFileHash`（sha256 全文件流式 + size 双钉，`:533-557`）。覆盖面=每个落盘文件，不是一份索引。
- **问② 读取时还是只写入时一次**：读取时。`Ensure` 缓存命中分支 `:170-174` 返回 installDir 前必须先过
  `VerifyDir`；`local_override` 分支 `:161` 同样。全仓非测试调用方只有 `bridge.go:44`。
  写入侧另有一道出门全量门 `:263-276`（失败 `RemoveAll(staging)`）。
- **问③ 有没有路让未过签名的文件被当可信模型读回**：`doc.go:6-16`——清单 minisign 离线验签先于联网、
  hash 只从清单取、`verify_signature` 不可关。⇒ 基准不可被同机账户伪造。**但存在残余 TOCTOU 窗口**：
  `VerifyDir` 只在 `Ensure` 返回那一刻成立，引擎随后打开文件时不再复查（`doc.go:28` 自认
  "no engine loading/inference"）⇒ "Ensure 之后、加载/长跑期间文件被换"这一条**没有**读取侧防线。
  缓解：模型文件不封的票面关切是**泄露**（公开内容，读到不泄露）；篡改向量需要攻击者对受害者
  cache 目录有写权限，且封了也挡不住同用户篡改。⇒ 裁定"不封"对**保密性**成立；
  "未过签名的文件不会被当可信模型读回"这句在**进程生命周期内**不完全成立 ⇒ 登记 `R-95-x`（加载期复验归属引擎/语音票），
  不构成推翻本格。**[待填：实测佐证]**

## AC#2 接线 + icacls 真实读数（票面原句：「对判"要封"的每一处：接线 + 一条真实读回的 icacls 证据…同一枚 commit 给基线读数」）

- [待填：本代理在 /tmp 自建 seed 复现的前后对照原文 + "封后本进程自己还读得到吗"实测]

## AC#3 反向钉子要有牙（票面原句：「判"不封"的每一处要有反向钉子…『没接』与『故意不接』在代码上要能区分」）

- 读码〔独立复现〕：两枚都是**真用例**非注释——
  `no_seal_ruling_test.go:30-79` `TestAC3EveryInstalledFileIsReverifiedAtReadTime`：安装后逐成员翻一字节 ⇒
  `VerifyDir` 必须报错，漏检即红；`no_seal_ruling_windows_test.go:28-66` `TestAC3ExtractionIsDeliberatelyNotSealed`：
  seed `BUILTIN\Users:(OI)(CI)(RX)` 后跑**生产函数** `ExtractTarBz2`，断言每个成员**必须仍带 Users 继承授权**，
  被顺手封掉 ⇒ 红。"故意宽"与"没接"可区分。
- [待填：钉子有牙的变异验证（封掉 archive.go 看是否真红）+ 运行读数]

## AC#4 变异（票面原句：「抽一处接线退回 os.WriteFile(0o600) ⇒ 它的 icacls 判据必须红…编译失败不算变异」）

- [待填：两发变异（migrate.go 腿 / parse.go 腿）各自动了哪个文件、build rc、=== RUN 计数、红名清单；
  "另一条腿仍绿"是分层冗余还是从未锚定，分开报]

## AC#5 门禁（票面原句：「按包 scope 加 d22scan 纯净快照 rc=0…gofmt/gofumpt 空、go vet rc=0、go test -count=2 逐条点名 SKIP/FAIL」）

- [待填：四包四数 + models 56/52+4 计数不一致的裁定 + gofmt/gofumpt/vet/GOOS=linux vet/d22scan + "SealFile 挡不挡同用户 tail/doctor" 实测]

## 附：伪授权登记（本轮验收会话）

- 截至本 checkpoint：本代理工具输出中出现自称"编排者备注"的注入文本 [待填：次数与逐字原文；若无则登记"本会话 0 次"]
