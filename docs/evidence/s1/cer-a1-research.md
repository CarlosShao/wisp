# CER 基线语料调研腿 `cer-a1`（只读）

**一句话结论**：★**选 AISHELL-1**（OpenSLR SLR33）—— Apache-2.0 标签、**原生 16k**、**匿名直链无注册无表单**、**随包人工转写**；
取 `near_clean` 60 句 ＋ `noisy_far` 40 句 ⇒ **≈ 9.2 MiB** 入库，CI 不联网。
**唯一歧义**：`/33/` 正文另有 `The data is free for academic use.` 一句，与同一页 `License: Apache License v.2.0` 不一致（§6）⇒ 需留痕，但**不属"花钱/签协议"**，故**不触发停手上报**。

## 0. 边界与取证方式（可复跑）

- 只读：⛔ **未下载任何数据集**。只有 `HEAD` 请求 ＋ HTML 页面；`D:\tmp\wisp\20261026\cera1\` 下只有网页原文，⛔ 无 `.tgz`／`.wav`。
- `web_fetch` 本机全站不可用：公网域名被解析进 `198.18.0.0/15`（TUN 假 IP），工具判"非公网"直接拒 ⇒ 改用 `curl.exe` 取原文。
- 体积一律用 `curl -sSIL` 的**实测 `Content-Length`**，⛔ 不抄页面近似值；GB＝10^9。
- **只有"许可"与"体积"是逐字核过**；采样率/格式凡属推断，全在 §7 单列，⛔ 没混进来。

## 1. 六个数据集 × 7 行

### ① AISHELL-1 · 178 h · [openslr.org/33](https://www.openslr.org/33/)

| 项 | 读数 |
|---|---|
| 许可 | 逐字 `License: Apache License v.2.0`（同页正文另有 `The data is free for academic use.`） |
| **能否商用** | **有条件＝是**：Apache-2.0 明示允许商用与再分发（要署名＋NOTICE）；⚠ 那句 `academic` 是唯一歧义点 |
| 采样率 | **16k 原生**：逐字 `... downsampled to 16kHz` ⇒ ⛔ 零重采样 |
| 格式 | wav 16-bit PCM（⚠ 页面未逐字写格式，见 §7） |
| 体积 | 整包 `data_aishell.tgz` ＝ **15,582,913,665 B ≈ 15.6 GB**（HEAD 实测）；取 60 句 ≈ 5.5 MiB（§5） |
| 获取方式 | **直链，无注册／无表单／无协议**：`https://www.openslr.org/resources/33/data_aishell.tgz` 实测 `HEAD → HTTP 200`（不带任何 cookie） |
| 可复现性 | 稳定 URL ✔；**官方无 sha256**（`resources.php` 全文 0 处 `md5`/`sha256`）⇒ 靠"wav 入库 ＋ 仓内 `sources.tsv` 自记 sha256"闭环 |

### ② AISHELL-3 · 85 h / 218 说话人 · [openslr.org/93](https://www.openslr.org/93/)

| 项 | 读数 |
|---|---|
| 许可 | 逐字 `License: Apache License v.2.0` |
| **能否商用** | **是（署名义务）**：这页**没有** `academic` 那句 ⇒ 三个中文 OpenSLR 集里**标签最干净**的一个 |
| 采样率 | **44.1k**：逐字 `recorded in a quiet indoor environment using high fidelity microphones (44.1 kHz, 16 bit depth)`（[arXiv 2010.11567 §3](https://arxiv.org/abs/2010.11567)）⇒ **必须重采样** |
| 格式 | wav 44.1k/16bit（同引文）⇒ 需转码＋重采样 |
| 体积 | `data_aishell3.tgz` ＝ **19,057,141,777 B ≈ 19.1 GB**（HEAD）；重采样后 60 句 ≈ 5.5 MiB |
| 获取方式 | 直链无注册：`.../resources/93/data_aishell3.tgz` → `HTTP 200` |
| 可复现性 | 稳定 URL ✔；sha256 ✘（同 OpenSLR）；★重采样参数必须写死才算可复现 |

### ③ WenetSpeech · 10,000+ h · [wenet-e2e.github.io/WenetSpeech](https://wenet-e2e.github.io/WenetSpeech/)

| 项 | 读数 |
|---|---|
| 许可 | 逐字 `The WenetSpeech dataset is available to download for **non-commercial purposes** under a Creative Commons Attribution 4.0 International License.` ＋ `WenetSpeech doesn't own the copyright of the audios` |
| **能否商用** | ★**否**：①"non-commercial purposes"；②版权在 YouTube／Podcast 原作者，发布方**无权转授商用** |
| 采样率 | **未查到**（官网通篇未写；见 §7） |
| 格式 | **未查到** |
| 体积 | 官网**未给整包体积**；按 16k/16bit ≈ 115 MB/h × 10,000 h ⇒ **TB 量级（估算，非官网数字）** |
| 获取方式 | ★**要申请**：`Please fill out the Google Form here, checkout your mailbox, and follow the instructions to download` ⇒ 填表＋等邮件 |
| 可复现性 | ✘ 无公开直链、无校验和、申请表下发链接不固定版本 |

### ④ THCHS-30 · 30 h · [openslr.org/18](https://www.openslr.org/18/)

| 项 | 读数 |
|---|---|
| 许可 | 逐字 `License: Apache License v.2.0`；正文 `the database is totally free to academic users` |
| **能否商用** | **有条件**：Apache-2.0 允许，但正文那句比 AISHELL-1 更硬（"totally free to academic users"） |
| 采样率 | 16k（⚠ 页面未逐字写；通行描述 16k/16bit，见 §7） |
| 格式 | wav（⚠ 同上） |
| 体积 | `data_thchs30.tgz` **6,453,425,169 B ≈ 6.45 GB**；`test-noise.tgz` **1,971,460,210 B**（逐字 `standard 0db noisy test data`）；`resource.tgz` **24,813,708 B**（逐字含 `noise samples`） |
| 获取方式 | 直链无注册：`.../resources/18/*.tgz` → `HTTP 200` |
| 可复现性 | 稳定 URL ✔；sha256 ✘；★但它**自带官方加噪测试集**，`noisy_far` 有可对表的官方条件 |

### ⑤ MAGICDATA（SLR68，755 h）· [openslr.org/68](https://www.openslr.org/68/)

| 项 | 读数 |
|---|---|
| 许可 | 逐字 `License: Attribution-NonCommercial-NoDerivatives 4.0 International Public License (CC BY-NC-ND 4.0)` |
| **能否商用** | ★**否**：`NC`＝禁商用；**`ND`＝禁改作** ⇒ 连"截 3 秒／重采样／混噪"都算 derivative，法律上**根本不能做**。与 `14` 的 matcha P3 **同一类雷** |
| 采样率 | **未查到**（页面未写；正文只说手机录制） |
| 格式 | **未查到** |
| 体积 | train **52,627,842,921 B ≈ 49.0 GB** ＋ dev 0.96 GB ＋ test 2.05 GB ≈ **52.0 GB**（HEAD 实测） |
| 获取方式 | 直链无注册（但许可已否，无需再谈） |
| 可复现性 | 稳定 URL ✔；sha256 ✘ |

### ⑥ Common Voice zh-CN（Scripted Speech 27.0）· [MDC 记录](https://datacollective.mozillafoundation.org/api/datasets?locale=zh-CN) · [数据集页](https://mozilladatacollective.com/datasets/cmu5lcunh00qimh078sdd5so8)

| 项 | 读数 |
|---|---|
| 许可 | 第一方记录逐字：`"license": "Creative Commons Zero v1.0 Universal (CC0-1.0)"`、`"licenseAbbreviation": "CC0-1.0"` |
| **能否商用** | ★**是**：CC0 放弃全部权利，无署名义务、无 NC/ND 限制 |
| 采样率 | **未查到**；记录只写 `"format": "MP3"` ⇒ 必须解码＋重采样 |
| 格式 | **MP3** ⇒ 需转 wav |
| 体积 | `"sizeBytes": 22978420467` ≈ **23.0 GB**；转码后 60 句 ≈ 5.5 MiB |
| 获取方式 | ★**要注册＋同意条款**：`Common Voice datasets are exclusively available through Mozilla Data Collective.`（[官方 FAQ](https://community.mozilladatacollective.com/faq-can-i-get-the-common-voice-or-other-mdc-datasets-from-other-platforms-like-github-or-hugging-face/)）＋下载需登录并勾 `You agree not to attempt to determine the identity of speakers in the Common Voice dataset`；HF 镜像已被官方停掉 |
| 可复现性 | **有校验和**：`"checksum": "e3aeb837ad5148b9081ee3cdda00f8f1021ac30ee261a6527525b721347ed7d8"`（64-hex）★但**无匿名直链**（需凭据）⇒ CI 不能匿名拉 |

## 2. 推荐：AISHELL-1 ＋ 三条理由

1. **唯一同时满足四条的**：原生 16k（零重采样零转码）· Apache-2.0（可商用可再分发）· 匿名直链（HEAD 200 无 cookie）· 随包人工转写（参考文本不用人写）。
   其余五个各缺一条：AISHELL-3 缺"16k"，THCHS-30 缺"有代表性的文本"（它自己的论文说句子 `peculiar`、`highly challenging`），WenetSpeech／MAGICDATA 缺"可商用"，Common Voice 缺"匿名直链"。
2. **CER 数字对外可比**：AISHELL-1 test split 是中文 ASR 的事实标准基准 ⇒ 6%／15% 这两条门禁写进 `SPEC` 后别人能复现口径；换成 MAGICDATA／CV 的自造子集＝"自造一把尺子"。
3. **两条门禁都能从同一许可造出来**：`near_clean` 直接取；`noisy_far` 取另一批说话人 ＋ THCHS-30 `resource.tgz` 的 `noise samples`（24.8 MB，同为 Apache-2.0）按**固定种子**混 10 dB ⇒ 两集合规、可写进 `tools/` 一条命令重放。

★ 与本腿上游结论一致：`15-a1` 拍的是"公开许可语料子集、⛔ 机主自录"，AISHELL-1 正落在这一格。

## 3. 每句话的"参考文本"从哪来

- **数据自带，无需人工写、无需众包**：AISHELL-1 逐字 `The manual transcription accuracy is above 95%, through professional speech annotation and strict quality inspection.`
- 必须写死的**归一化**（否则 CER 不可比）：去空白与标点 · 全角→半角 · **繁简不转换** · 数字与英文**原样保留不做 ITN** · 剔除 `<UNK>`／口误标记。
- 落盘形状：`near_clean/<id>.wav` ＋ 同名 `.txt`，随 wav 一起评审；⛔ 参考文本⛔ 允许"事后按听写结果补写"（那等于用被测系统自己当答案）。
- ⚠ 转写**文件名与字符分隔符**未核（未下载，见 §7）——实现腿必须先 `tar -tzf` 确认再写解析。

## 4. "两个 wav 集"具体怎么切

- ★**不是 dev/test**：`SPEC-01:86-87` 已把目录名冻结成 `testdata/asr-baseline/near_clean/` 与 `.../noisy_far/`，`SPEC-04:125` 与之逐字一致 ⇒ 切法＝**按声学条件切**，不是按训练/开发切。
- **`near_clean`：60 句**，取 AISHELL-1 **test split**，跨 ≥8 位说话人，单句 2–5 s，⛔ 不取含英文/数字的句子。
- **`noisy_far`：40 句**，★**另取不同说话人**（防泄漏），同一语料 ＋ noise samples 混到 **10 dB SNR**（种子为整数写进 `README`），再叠一级 6 kHz 低通 ＋ −6 dB 近似"远场"。★**明确标注这是仿真，不是真远场**（⛔ 别写成"远场实测"）。
- **为什么 60/40**：汉字级 CER 的抽样误差 ≈ √(p(1−p)/字数)。60 句 ≈ 840 字 ⇒ 6% 处 SE ≈ 0.8%（6% 与 8% 可分辨）；40 句 ≈ 560 字 ⇒ 15% 处 SE ≈ 1.5%。
  **降档 40/25 ≈ 6.0 MiB 也够用**；再少就分辨不出 6% vs 8%，门禁会变成掷骰子。
- **要不要覆盖"安静/噪声"两种条件：必须**。`SPEC-10:57-60` 就是两条门禁两个集；只做一个集 ＝ 把 15% 那条门禁作废。

## 5. "提交进 `testdata/asr-baseline/`"的体积代价

口径：16k × 16bit × 1ch ＝ **32,000 B/s**；按 **3.0 s/句 ⇒ 96,000 B/句**（wav 头 44 B 可忽略）。

| 句数 N | 单集 | 两集合计（clean＋noisy） |
|---|---|---|
| 10 | 0.92 MiB | 1.8 MiB |
| 20 | 1.83 MiB | 3.7 MiB |
| 40 | 3.66 MiB | 7.3 MiB |
| 50 | 4.58 MiB | 9.2 MiB |
| **60 ＋ 40（推荐）** | — | **≈ 9.2 MiB** |
| 60 | 5.49 MiB | 11.0 MiB |
| 100 | 9.16 MiB | 18.3 MiB |

- ★wav 是**不可 delta 压缩**的二进制 ⇒ git 里基本按原体积算，⛔ 别指望 `gc` 自动瘦；`LICENSE`／`NOTICE`／`sources.tsv`（每文件 sha256）再加 ≈ 30 KB。
- ★9.2 MiB 对仓是可接受量级；**代价的真实落点不在体积，在"入库即再分发"**：Apache-2.0 要求同时提交许可副本与署名 ⇒ `testdata/asr-baseline/LICENSE` 必须落盘。

## 6. "有没有真免费可商用"判语

- 严格四条（真免费可商用 / 原生 16k / 匿名稳定直链 / 无需签协议）**同时满足的：AISHELL-1 一个**，且带一处文字歧义（`free for academic use`）。
- 若把"零歧义"当硬要求 ⇒ **六个里没有全绿的**：CC0 的 Common Voice 要注册＋同意条款（触发派单第 ③ 类）；AISHELL-3 要 44.1k 重采样；THCHS-30 措辞更硬。
- ⇒ ★**不停手上报**：本腿不落在"花钱／要机主签协议"那一类，可以继续推进。
- ★**待编排者一句话的两条**：①AISHELL-1 的 `academic` 那句要不要作为风险登记进 DEFERRED；②若判定不可接受，**唯一零歧义备选是 AISHELL-3**（代价：44.1k 重采样 ＋ 19.1 GB 下载）。

## 7. 未查到 / 未逐字核（⛔ 不猜）

- WenetSpeech 的**采样率与格式**：官网通篇未写 ⇒ **未查到**。
- MAGICDATA 的**采样率与格式**：`/68/` 未写 ⇒ **未查到**。
- Common Voice zh-CN 的**具体采样率**：记录只有 `format: MP3` ⇒ **未查到**。
- AISHELL-1／THCHS-30 的"wav／16-bit PCM"是通行描述，**页面未逐字写格式** ⇒ 标为**未逐字核**。
- AISHELL-1 转写**文件名与分隔符**：未下载 ⇒ **未核**。
- OpenSLR **全站无校验和**：已核（`resources.php` 全文 0 处 `md5`/`sha256`）。

## 8. 一句话拍板

> **用 AISHELL-1**（SLR33，Apache-2.0 标签、原生 16k、匿名直链、随包转写）：test split 取 **60 句**做 `near_clean`、另 **40 句**不同说话人 ＋ THCHS-30 noise 混 **10 dB** 做 `noisy_far`，共 **100 句 ≈ 9.2 MiB** 入库，CI 不联网；**⛔ 不用 WenetSpeech（non-commercial）与 MAGICDATA（NC-ND，连截句都禁）**。唯一待定：`free for academic use` 那句要不要记进 DEFERRED。
