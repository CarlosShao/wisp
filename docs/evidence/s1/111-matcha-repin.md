# 票 111 波 8a —— matcha TTS 夹具 pin 重新钉到当前上游真实字节（腿 111-wave8a，2026-11-13）

派单：`.scratch/wisp/probes/orch/20261113/txt/111-wave8a-matcha-repin-dispatch.txt`。
裁定来路（编排者四组实验，派单「为什么」节）：上游 release 资产已被换（换的是
`model-steps-3.onnx` 本体），HF 11 个 revision 全查过无一处等于旧 pin ⇒ 旧 pin 除本机存量外
无处可得；留着一个谁也拿不到的 pin = fetch 步永久红 = 掩盖真问题。本腿执行重钉＋验证。
本机旧夹具（`third_party/spike-models/`，只读未动）实测仍等于旧 pin（见下），是旧模型最后标本。

## 0. 一个现场发现（写面扩展的依据，先报后做）

派单点名 `build/models/manifest.json`。现量：**该路径整个 `build/` 被 .gitignore 忽略**
（`.gitignore:14` `build/`；`git check-ignore -v` 逐字命中；`git ls-files -- build/models/manifest.json` 空）。
CI 链路：被跟踪的真账本是 **`models/manifest.json` ＋ `models/manifest.json.minisig`（C29 签名对，
`git log` 唯一笔 `bfcb230e3`）**；`scripts/build.ps1:259,286-290` 在 CI build 步里把它复制成
`build/models/manifest.json(+.minisig)`，fetch 步（`ci.yml:1089`）读到的就是这份复制件，
cache key `ci.yml:1026` 也 hash 它。cmd/wisp 的 `WISP_MODELS_MANIFEST` 测试注入面
（`resident_voice_15e_windows_test.go:684` 等 9 处）同样指向 `build/models/manifest.json`
并经 `LoadSignedManifest` 验签加载（`internal/models/manifest.go:248`）。

⇒ 只改 build/ 下那份 = git 零 diff = CI 空操作，且改动未签名的字节会让测试注入面的
C29 验签假红。**执行意图（重钉让 CI 冷跑能过）落在真账本上**：`models/manifest.json` 改 pin
＋仓内自带 `tools/signmodels` 重签（dev key 在盘上，`scripts/sign-models.ps1` 即文档化路径；
生产 key 仪式 S8 才发生），**再镜像回 `build/models/manifest.json(+.minisig)`** 保持本地两份一致。
build/ 那份被 gitignore ⇒ commit 只含跟踪的三件：`models/manifest.json`、`models/manifest.json.minisig`、证据件＋票面。
旧账本对已备份（可回滚）：`D:\tmp\wisp\20261113\111-wave8a\backup-old-pin\manifest.json.old`
（sha256 `b97e8af4999f5e819a1d651ba9f9fe13fffe5a2e1c6ebab8d7968ee600c5431c`）＋
`manifest.json.minisig.old`（sha256 `59afa056ed8e5e240a07ec0b172a4cbb96a687b201dc6b14a29e18f8a359b613`）。

## 1. 第 1 步：新 pin 现下现算（★全部本腿自算，非抄派单）

归档现下（2026-11-13，从 manifest 里的 GitHub release URL，
`https://github.com/k2-fsa/sherpa-onnx/releases/download/tts-models/matcha-icefall-zh-baker.tar.bz2`）：

```
DOWNLOADED size=75455267 sha256=ae70096b19c6e283ea351a26b8a718ce6427c8bde8c1da22eb46a5478d39a69b
```

解包（Windows 系统 tar，exit 0，顶层目录恰好 1 个 `matcha-icefall-zh-baker`），18 个成员逐个现算
（相对路径 | size | sha256）：

```
date.fst | 59154 | eb8aa079ae3cb81d8f4404992f39d61a0cb990947512b5b8d1e54d1f6980e718
dict/generate_user_dict.py | 545 | 24541017b4b60d8e1ad94a5165f7097d2fd107d4ceecbf13583dd37682fdd2a7
dict/hmm_model.utf8 | 519739 | f17790586ac86dd048c8adffed052c4bd2b28ed0682972c1275e59040c0589a7
dict/idf.utf8 | 5998717 | dbd1e03d72b2263cc8d84a4304ed77677eed9e7deaf43a1a5133bbba9733b535
dict/jieba.dict.utf8 | 5071204 | 3043b77068e09c9904f27cad82f12b6ebe9dbdb5aeff3b25e45ab7f9c1122b55
dict/pos_dict/char_state_tab.utf8 | 327139 | 28b7be1dd7369766a51445af4d42e9a2ba4bf374c13be5bc1ca7721e27271dbb
dict/pos_dict/prob_emit.utf8 | 1687686 | c33c4cb7edf3b3a5947df7209b6e9f267eae1f21335d9e2bd2521ea07105457a
dict/pos_dict/prob_start.utf8 | 4347 | 13623ea0e9300bdb597cb2da28770b7b385d6c0098d66e516083fb01b6bd5d96
dict/pos_dict/prob_trans.utf8 | 124159 | f22363e2307408293d180c6f9f6b5cb75879d52f722f7764fa2d3d0ae2400236
dict/README.md | 683 | 8bcc7c2ed7dff082d9a99c09de0bfd72249bf3b5c4298a3af328da40cbd8cfd1
dict/stop_words.utf8 | 8974 | b788b8a939d2e2fe079abd579ea98f12f9fb84370bfd0dddd81bb9381f7ab42c
dict/user.dict.utf8 | 829194 | 5faabcd0c52adf512de77fd7eee3d556da8c98b50f493b47c66a981b3a24959c
lexicon.txt | 1363705 | 38b886d46aefa50da6322a64d72fd595d5f4fae1051adb160d647541b1e0a4a2
model-steps-3.onnx | 75624188 | 47e82b11454edeeb65adc4e44a24fd54c6715cd4d875cc2a20f75545ecc32d4d
number.fst | 64482 | 743f402181fcfebf76cc2f0546b71fa26476e626fbe4e460fb7b4c3a7a8bd5bd
phone.fst | 88630 | 1ac2b6fa56b1442320c4de7db08353bab8963a2b57f365eebcdd3a2d3562f8d7
README.md | 370 | d47ef70781df15845b6f28c6c071dee9e51b52c36797a8130a620c5c8895929c
tokens.txt | 19591 | 56209b2bf609d5ac1d66ede6dae7bf5254bd3f8aa24c4a6823713d5b884d87ba
```

对照旧 pin：**14 条旧成员逐字节全对**（tokens/lexicon/date/number/phone + dict 全树）；
唯一不符 = `model-steps-3.onnx`（旧 pin `ef7ebdf5987e16a5836136a51d6f3560ca997ffd33d06a40daab5af92b4b86e5`
/ 75624611 → 新 `47e82b11454edeeb65adc4e44a24fd54c6715cd4d875cc2a20f75545ecc32d4d` / 75624188）；
恰好 3 个包内实际存在但 pin 未列的文件 = `README.md`、`dict/README.md`、`dict/generate_user_dict.py`。
交叉铁证：本机 `third_party/spike-models/matcha-icefall-zh-baker/model-steps-3.onnx`（只读哈希）
＝ `ef7ebdf5987e16a5836136a51d6f3560ca997ffd33d06a40daab5af92b4b86e5` / 75624611，**逐字等于旧 pin**
⇒ 本机存量=旧模型、上游现售=新模型，与派单实验一致。本机旧夹具全程只读未动。

### 新旧 pin 对照表（归档级 + 成员级）

| 字段 | 旧 pin | 新 pin（现算） |
|---|---|---|
| 归档 sha256 | `20de2ec034b55562609d6362771c934905dfe11d0f41ec103d593427ad9a7efb` | `ae70096b19c6e283ea351a26b8a718ce6427c8bde8c1da22eb46a5478d39a69b` |
| 归档 size_bytes | 75463442 | 75455267 |
| `model-steps-3.onnx` sha256 | `ef7ebdf5987e16a5836136a51d6f3560ca997ffd33d06a40daab5af92b4b86e5` | `47e82b11454edeeb65adc4e44a24fd54c6715cd4d875cc2a20f75545ecc32d4d` |
| `model-steps-3.onnx` size_bytes | 75624611 | 75624188 |
| 新增成员 `README.md` | （未列） | `d47ef70781df15845b6f28c6c071dee9e51b52c36797a8130a620c5c8895929c` / 370 |
| 新增成员 `dict/README.md` | （未列） | `8bcc7c2ed7dff082d9a99c09de0bfd72249bf3b5c4298a3af328da40cbd8cfd1` / 683 |
| 新增成员 `dict/generate_user_dict.py` | （未列） | `24541017b4b60d8e1ad94a5165f7097d2fd107d4ceecbf13583dd37682fdd2a7` / 545 |
| 其余 14 条成员 | 与 pin 一致 | **一字未动** |

### 连带真值（第 5 处改动，派单字段清单之外，依据如下）

条目级聚合 `sha256`/`size_bytes`（旧 `f2b9e766…36a6a0` / 129347466）在 C29 账本里有**精确定义**
（`internal/models/manifest.go:26-31` 逐字）：多工件条目的条目级 sha256 = **按 Files 顺序对全部工件
字节拼接后的 sha256（"combined digest"）**，size_bytes = 工件总字节和；由 `tools/signmodels` 签发时
重算、`TestManifestCombinedDigestContract` 钉机制、`TestRealManifestAuthoringConsistency:90-98`
机器断言「files[].size_bytes 之和 == 条目 size_bytes」。归档字节换了，聚合值必须跟着换，否则
authoring consistency 断言红。方法自证：用同一算法对 **asr-streaming 条目**（encoder+decoder+tokens
按 Files 序拼接）现算得 `c129b662000dc429d2058048b752081c1cfb39224f9e1ea52e0fa661a7540190`，
与账本里该条目逐字相同 ⇒ 算法与账本定义同轴。新聚合值（新归档 75455267 + vocos 53884024，
vocos 本机现哈希仍逐字等于其 pin `0574a135…0081`，未动）：

```
NEW combined digest = a0bc1fbb7da681ab7a87707bcfa3dafd1665c635cdfaccc4c9d90207eb5d84fc  total=129339291
（= 75455267 + 53884024，机器断言自洽）
```

### 落盘与 C29 重签（真账本）

- `models/manifest.json`：上述 5 处定点替换（归档 sha/size、model-steps-3 sha/size、3 个新成员行、
  条目级聚合 sha/size）；该条目 `status`(blocked-p3)/`license`/`urls`/`id`/`purpose`/`quant` 未动；
  manifest 其余 5 条目一字未动。`py -3 json.load` OK；与 build 材质件逐字节相同（hash 同）。
- 重签：`powershell -NoProfile -ExecutionPolicy Bypass -File scripts\sign-models.ps1` →
  `signmodels: signed ... -> models\manifest.json.minisig (keyid a3c8794f3fd94fc5, runtime-verified)` rc=0；
  `-Verify` 独立复核 rc=0；`go test ./internal/models -count=1` 全包 ok（11.781s，含
  TestRealManifestInRepoVerifies / TestRealManifestAuthoringConsistency / TestP3BlockedModelRefused）。
  改前基线同靶三测 ok（0.049s）。
- 镜像：新 manifest+新 sig 复制回 `build/models/`（材质件，gitignore 不入库）。

## 1b. 条目级 `sha256`/`size_bytes` 的消费点证明（本腿续证，2026-11-13 晚；全带 文件:行号）

### 它是什么（代码原文）

`internal/models/manifest.go:19-31`（Manifest doc 注释，逐字要点）：
> Multi-artifact entry (Files != nil): … entry-level urls[] is the flattened index … (documentation
> field, not the download list); **entry-level sha256 = sha256 over the byte-concatenation of all
> artifacts in Files order ("combined digest") and size_bytes = their total**. Authoring-time
> invariant, recomputed by tools/signmodels and pinned by TestManifestCombinedDigest.

算法本体＝`internal/models/manifest.go:233-242`（`CombinedDigest`：对 Files 顺序的工件字节流式
写同一个 sha256 state——不是逐文件哈希再拼哈希）。

### 谁消费它（全仓枚举，非抽样）

**A. 运行时下载/校验路径＝零消费（多工件 entry 不看这两个字段）：**
- 下载循环 `internal/models/downloader.go:336-360`：`entry.ArtifactCount()`→`entry.Artifact(i)`
  （`manifest.go:211-216`，`Files` 非空时逐个返回 `Files[i]`），逐工件 `a.SHA256/a.SizeBytes` 校验
  （`downloader.go:385` 终验、`:410-411`、`:484-490` 续传、`:546-557` 下载体 size+sha256）。
- 安装终门 `downloader.go:377-390` 与缓存/override 重验 `VerifyDir`
  （`downloader.go:646-673`）都走 `entry.InstalledFiles()`（`manifest.go:218-231`）＝**只含成员级 pin**。
- 条目级 `e.SHA256`/`e.SizeBytes` 在 `internal/models/**` 的全部读点（rg 枚举）：
  `manifest.go:145/151`（ParseManifest 结构校验：64hex/非负）、`manifest.go:207`（ImplicitFile，
  仅 `Files==nil` 的单工件条目把条目级值当工件 pin 用）、`manifest_real_test.go:137`（VAD 真网下载
  抽检，单工件条目）。**没有任何运行时代码对多工件 entry 读条目级 sha256/size_bytes。**

**B. 作者一致性断言（真正盯着它的机器）：** `internal/models/manifest_real_test.go:83-98`
`TestRealManifestAuthoringConsistency`：`len(Files)==1` ⇒ 条目级必须等于 `Files[0]`（:85-89）；
`len(Files)>1` ⇒ **`Σ files[].size_bytes == entry.size_bytes`（:90-98）**。多工件 entry 的条目级
sha256 的算法级钉子＝`internal/models/manifest_test.go:86-100`（`TestManifestCombinedDigestContract`
钉 `CombinedDigest` 本身）＋`manifest.go:31`「recomputed by tools/signmodels」（作者时不变式，
`tools/signmodels/main.go` 的 sign 路径 :115-147 只签文件字节、重算属作者工序）。

**C. 展示面：** `cmd/wisp/models.go:221-222`（`wisp models list` 打印 `e.SizeBytes`，仅显示）。

**D. CI 夹具脚本：** `scripts/fetch-voice-fixtures.ps1` 对 `files`/`archive` 目标只用成员级 pin
（`$fe.sha256` :176/206/263、`$fe.archive.files` :247-248）；条目级值只在单文件目标（VAD）
被用（`$model.sha256` :178）。

### 结论（三情形判定：**(a)**，且新值经独立重算命中）

条目级值**有精确定义、有机器断言**（B 组），不是"无消费点"也不是"仅单文件用"。本腿用
`manifest.go:236-242` 的同一算法对**两条 entry 各自独立新下全工件**重算：

```
TTS  matcha+concat(vocos):  matcha=75455267/ae70096b19c6e283ea351a26b8a718ce6427c8bde8c1da22eb46a5478d39a69b
                            vocos =53884024/0574a135aa1db2de6e181050db2ec528496cacd4a4701fc5d7faf9f9804c0081
                            combineddigest=a0bc1fbb7da681ab7a87707bcfa3dafd1665c635cdfaccc4c9d90207eb5d84fc
                            == 盘上 models/manifest.json:173 逐字命中；size 75455267+53884024=129339291 == :174 ✓
ASR  encoder+decoder+tokens 拼接: combineddigest=c129b662000dc429d2058048b752081c1cfb39224f9e1ea52e0fa661a7540190
                            == 账本 :73 逐字命中（算法与账本同轴的第二例证，非 matcha 自证）
```

⇒ 上一腿把条目级 `sha256` f2b9e766…→`a0bc1fbb…`、`size_bytes` 129347466→`129339291` **改得对**
（作者不变式要求归档字节换了聚合值必换，`manifest_real_test.go:90-98` 会盯 size 和）；**保留，
不再改，不再重签**（签名对已覆盖现值，`sign-models.ps1 -Verify` rc=0）。

## 2. 第 2 步：冷跑 fetch（空临时树，从零下载）

布置：`D:\tmp\wisp\20261113\111-wave8a\cold\` 为全新空树，内含 `scripts/fetch-voice-fixtures.ps1`（仓内
原样复制）＋ `build/models/manifest.json(+.minisig)`（仓内 `models/` 真账本的原样复制，材质一致：两处
`Get-FileHash` 逐字相同——脚本读 `build\models\manifest.json`，见脚本 :63）。脚本以 `$repo=父目录` 定位
仓库根 ⇒ 下载全部落该临时树，本机 `third_party/` 全程只读未动（未发生写）。跑法＝CI 同款解释器
`powershell -NoProfile -ExecutionPolicy Bypass -File …`（CI 步 `ci.yml:1116` 同款）。2026-11-13 19:57–20:1x，
约 20 分钟（hf-mirror 对每个 ASR 工件先拒 308 再落 huggingface 候选；matcha/vocos/punc/kws 走 GitHub）。

**fetch-voice-fixtures 输出逐行（后台作业流式回收，行序即产出序）：**

```
fetch-voice-fixtures: vad-silero: downloading -> D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\silero_vad.onnx
fetch-voice-fixtures: vad-silero: ok (downloaded, sha256 verified)
fetch-voice-fixtures:   downloading asr-encoder.int8.onnx
fetch-voice-fixtures:   url failed: https://hf-mirror.com/csukuangfj/sherpa-onnx-streaming-paraformer-bilingual-zh-en/resolve/main/encoder.int8.onnx (The remote server returned an error: (308) Permanent Redirect.)
fetch-voice-fixtures:   asr-encoder.int8.onnx: ok (downloaded, sha256 verified)
fetch-voice-fixtures:   downloading asr-decoder.int8.onnx
fetch-voice-fixtures:   url failed: https://hf-mirror.com/csukuangfj/sherpa-onnx-streaming-paraformer-bilingual-zh-en/resolve/main/decoder.int8.onnx (The remote server returned an error: (308) Permanent Redirect.)
fetch-voice-fixtures:   asr-decoder.int8.onnx: ok (downloaded, sha256 verified)
fetch-voice-fixtures:   downloading asr-tokens.txt
fetch-voice-fixtures:   url failed: https://hf-mirror.com/csukuangfj/sherpa-onnx-streaming-paraformer-bilingual-zh-en/resolve/main/tokens.txt (The remote server returned an error: (308) Permanent Redirect.)
fetch-voice-fixtures:   asr-tokens.txt: ok (downloaded, sha256 verified)
fetch-voice-fixtures: tts-matcha-zh-baker: downloading archive (75455267 bytes pinned)
fetch-voice-fixtures: tts-matcha-zh-baker: ok (extracted, all pinned members sha256 verified)
fetch-voice-fixtures: tts-matcha-zh-baker: downloading -> D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\vocos-22khz-univ.onnx
fetch-voice-fixtures: tts-matcha-zh-baker: ok (downloaded, sha256 verified)
fetch-voice-fixtures: punc-ct-transformer-zh-en-vocab272727: downloading archive (64717756 bytes pinned)
fetch-voice-fixtures: punc-ct-transformer-zh-en-vocab272727: ok (extracted, all pinned members sha256 verified)
fetch-voice-fixtures: kws-zipformer-wenetspeech-3.3M-2024-01-01: downloading archive (32654866 bytes pinned)
fetch-voice-fixtures: kws-zipformer-wenetspeech-3.3M-2024-01-01: ok (extracted, all pinned members sha256 verified)
fetch-voice-fixtures: all five targets verified:
fetch-voice-fixtures:   WISP_SPIKE_MODELS=D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models
fetch-voice-fixtures:   WISP_SILERO_MODEL=D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\silero_vad.onnx
fetch-voice-fixtures:   WISP_TTS_DIR=D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\matcha-icefall-zh-baker
fetch-voice-fixtures:   WISP_KWS_DIR=D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01
fetch-voice-fixtures:   WISP_PUNC_MODEL=D:\tmp\wisp\20261113\111-wave8a\cold\third_party\model-fetch\x-punc\sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8\model.int8.onnx
```

**exit code＝0**（包装行 `COLD-RUN-EXIT=0`，脚本自身 `exit 0`）。判定：**vad / asr 三件 / tts-matcha
归档＋vocos / punc / kws 全部 `ok`，五目标（六行判定）齐，从零下载无一处用到本机旧夹具**。
说明：`url failed … 308` 三行是脚本对 hf-mirror 候选的如实记录（Windows PowerShell 5.1 的
`Invoke-WebRequest` 不跟随 308，转下一个候选 huggingface——CI 同解释器同行为，`ci.yml:1116`），
不构成目标失败；无一次重试即全过。

## 3. 第 3 步：新夹具喂 carrier 跑 TTS 用例

夹具＝**冷跑新下的树**（第 2 步产物，非本机旧夹具、非上一腿的手拼树）：
`D:\tmp\wisp\20261113\111-wave8a\cold\third_party\spike-models\`，五个环境变量全部指向它
（`WISP_SPIKE_MODELS`/`WISP_SILERO_MODEL`/`WISP_TTS_DIR`/`WISP_KWS_DIR`/`WISP_PUNC_MODEL`），
carrier（`scripts/wisp-cli-staged-tests.sh`）日志自证 `voice fixture environment resolved (caller values win)`。
命令＝派单原文 `CGO_ENABLED=1 sh scripts/wisp-cli-staged-tests.sh -test.run "^(Test382WR2|Test41a2r1|Test26P|Test390r1|Test394R1)"`
（经一层 bash 包装脚本 `D:\tmp\wisp\20261113\111-wave8a\run-tts.sh` 执行，使正则参数只在 bash 内引号化；
首次直接从 pwsh 传参被原生参数跳变吞掉正则，载体报 `flag provided but not defined: -test` exit 2——
是调用方引号层错误，非载体缺陷，如实记）。

**四数（逐字，尺＝`.scratch/wisp/stage-cli/wisp-cli-tests.log` 顶格行计数）**：

```
RUN=19  PASS=14  FAIL=0  SKIP=5
```

- 包级结果行：`PASS` ＋ `ok  github.com/CarlosShao/wisp/cmd/wisp	153.00s`；二进制 rc=0，`TTS-CARRIER-EXIT=0`。
- SKIP=5 全部是 `Test26P*`（EchoNeverTriggers/WithoutAEC/MarginRetest/BargeinStopsWithin400ms/FeedOnline），
  跳因＝测试自报 `the AEC3 binding is not in this build (-tags webrtcap + the deps.toml CGO flags)`——
  build-tag 门控的集成读数缺绑定，**与新模型无关、非红、期望值零改动**（⛔ 未动任何断言）。
- 19 枚构成：41a2r1×5 全 PASS（含两条真唤醒端到端）＋ 26P×5 全 SKIP ＋ 390r1×4 全 PASS
  （含 400ms unstick 真设备读数 n=5 min=33ms max=45ms）＋ 394R1×2 全 PASS ＋ 382WR2×3 全 PASS。
- ★判定：**新 matcha 权重（`47e82b11…`）没有打红 voice 家族**——TTS 真出声由新模型完成
  （`announcement spoken runes=30 chunks=4 samples=140671 rate=22050` 等多行；
  `Test382WR2TTSCleanSentencePlaysThroughTheGate` samples=75378 PASS）。
  顺带实证：41a2r1 端到端把 `wisp models ensure` 打到**新签 manifest** 上
  （`model hand-off: effect d43=37 name=model.verify-sha256-signature at=1` → `hand-off OK`）——
  运行时 C29 链在新 pin 上自证通过。
- 已知噪声（非失败）：`lexicon.cc:ConvertTokensToIds:93 Unknown token: shei2`（`lexicon.txt` 新旧
  归档逐字节相同，本告警非本重钉引入）。

## 门禁

- `py -3 -c "import json;json.load(open('models/manifest.json'))"` → `JSON-OK`，rc=0。
- `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/sign-models.ps1 -Verify` →
  `signmodels: OK D:\work\workspace\projects plans\Wisp\models\manifest.json (keyid a3c8794f3fd94fc5)`，rc=0。
- `sh scripts/d22scan.sh` → `d22scan: clean - no D22 ban violations`（examined 480 production Go files，
  45/45 包级 PASS），rc=0。

## 变更面与撤销口令

- 跟踪写面恰 4 件：`models/manifest.json`、`models/manifest.json.minisig`（均上一腿已改好，
  本腿零再改、零重签）、`docs/evidence/s1/111-matcha-repin.md`（本腿补 §1b/§2/§3/门禁/本节）、
  票 111 Progress log（append-only 一条）。显式 pathspec commit，⛔ push。
- 本机 `third_party/spike-models/` 全程只读未动；临时件全部落 `D:\tmp\wisp\20261113\111-wave8a\`
  （`cold\`＝冷跑树与夹具、`combined\`＝重算用独立下载件、`run-tts.sh`＝调用包装）。
- 撤销口令：回「撤 111-wave8a repin」⇒ 用 `D:\tmp\wisp\20261113\111-wave8a\backup-old-pin\`
  的 `manifest.json.old`＋`manifest.json.minisig.old` 复原两件并 revert 本笔 commit。
