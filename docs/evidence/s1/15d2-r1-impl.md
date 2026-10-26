# 票 15d2 CER 基线集腿 `15d2-r1` —— 实现记录（写腿）

**一句话**：AISHELL-1（SLR33）test split 截 **60 句**做 `near_clean`、另取 **40 句不同说话人**混 THCHS-30 真噪声到
**10 dB** 做 `noisy_far`，共 **100 句 / 12,971,844 B / 12.37 MiB** 入 `testdata/asr-baseline/`；
`sources.tsv` 100 行逐句 sha256；`LICENSE` ＝ 署名 ＋ Apache-2.0 全文；生成器确定性可重放。

**提交**：`6d120dca19449a33e258b8af7a5037792acc4350`（205 files，⛔ push）。**上游件**全在仓外 `D:\tmp\wisp\20261026\15d2\`，⛔ 入库。

## 读数① 两个集的句数、逐句时长、总字节

| 集 | 句数 | 说话人 | wav 字节 | MiB | 时长 min/mean/max (s) | 合计时长 |
|---|---|---|---|---|---|---|
| `near_clean` | **60** | 8 | 7,654,144 | 7.30 | 2.987 / 3.985 / 4.927 | 239.11 s |
| `noisy_far` | **40** | 8 | 5,317,700 | 5.07 | 3.056 / 4.153 / 4.971 | 166.12 s |
| **合计** | **100** | 16 | **12,971,844** | **12.37** | 2.987 / 4.052 / 4.971 | 405.23 s |

`near_clean` 每说话人句数 8/8/8/8/7/7/7/7（S0764-S0901）；`noisy_far` 每人 5（S0902-S0912）；两集说话人**交集 0 位**（防泄漏）。

**逐句时长**（单位 s；行序**＝**`sources.tsv` 行序，键 `S0764W0121` 即 `BAC009S0764W0121.wav`）：

```text
near_clean (60):
  S0764W0121=4.204 S0765W0121=3.092 S0766W0121=3.174 S0767W0121=4.588 S0768W0122=3.915 S0769W0123=4.137 S0770W0122=4.702 S0901W0121=3.098
  S0764W0122=4.115 S0765W0123=3.302 S0766W0122=3.146 S0767W0124=3.778 S0768W0125=4.317 S0769W0124=4.667 S0770W0124=3.940 S0901W0125=3.338
  S0764W0123=4.000 S0765W0125=3.669 S0766W0123=3.928 S0767W0125=4.239 S0768W0129=3.945 S0769W0125=3.909 S0770W0127=4.626 S0901W0126=4.883
  S0764W0125=4.311 S0765W0127=3.162 S0766W0124=4.519 S0767W0126=4.014 S0768W0130=4.634 S0769W0126=4.304 S0770W0128=4.401 S0901W0127=3.885
  S0764W0126=3.500 S0765W0128=2.987 S0766W0126=3.575 S0767W0127=4.059 S0768W0132=3.813 S0769W0129=3.461 S0770W0134=4.693 S0901W0128=3.773
  S0764W0127=4.453 S0765W0129=4.142 S0766W0127=4.014 S0767W0129=4.633 S0768W0133=3.991 S0769W0131=4.314 S0770W0135=4.927 S0901W0129=3.071
  S0764W0128=3.571 S0765W0130=3.135 S0766W0130=4.224 S0767W0130=3.396 S0768W0135=4.571 S0769W0135=4.272 S0770W0139=4.410 S0901W0131=4.596
  S0764W0129=3.595 S0765W0131=3.705 S0766W0131=3.947 S0767W0132=4.340
noisy_far (40):
  S0902W0122=4.650 S0903W0121=3.841 S0904W0121=3.056 S0905W0121=3.601 S0906W0122=4.299 S0907W0122=4.636 S0908W0122=4.902 S0912W0121=3.965
  S0902W0125=4.098 S0903W0122=4.650 S0904W0124=4.457 S0905W0126=4.889 S0906W0129=4.971 S0907W0123=4.701 S0908W0125=4.149 S0912W0124=4.267
  S0902W0126=3.670 S0903W0123=4.408 S0904W0125=3.156 S0905W0127=4.224 S0906W0131=3.291 S0907W0125=4.852 S0908W0127=3.925 S0912W0125=3.352
  S0902W0127=3.127 S0903W0124=3.769 S0904W0126=3.935 S0905W0128=4.784 S0906W0132=3.800 S0907W0126=4.885 S0908W0135=4.294 S0912W0127=4.000
  S0902W0128=4.152 S0903W0125=4.300 S0904W0127=4.746 S0905W0133=4.594 S0906W0140=3.499 S0907W0127=3.916 S0908W0137=3.366 S0912W0128=4.946
```

## 读数② `sources.tsv` 的结构 ＋ 3 行原文

1 行表头 ＋ **100 行数据**，制表符分隔，10 列。**前 4 列即票面要求的四项**：`file`（文件名）／`sha256`／
`source_utt_id`（来源句 id）／`duration_s`（时长，3 位小数、由 wav 头实测）；其余 6 列 `speaker`／`condition`／
`snr_db`（文件内实测，`near_clean` 为空）／`noise_source`（`cafe`／`car`／`white`）／`text_raw`（上游原文）／
`text_norm`（归一化参考文本，与同名 `.txt` 逐字相同）只为可审计，**非必需**。下面 3 行原文制表符逐字保留、⛔ 截断：

```text
file	sha256	source_utt_id	duration_s	speaker	condition	snr_db	noise_source	text_raw	text_norm
near_clean/BAC009S0764W0121.wav	46dbc998c9d1d48111267c40741dd3200f2e5bcf4075f8c4c97f4451160dce50	BAC009S0764W0121	4.204	S0764	near_clean			甚 至 出 现 交 易 几 乎 停 滞 的 情 况	甚至出现交易几乎停滞的情况
near_clean/BAC009S0765W0121.wav	60cca1290fb34656b085a6fa079a7cbc3c6ee3b78bad157748bc43322803c003	BAC009S0765W0121	3.092	S0765	near_clean			一 线 城 市 出 现 日 光 盘	一线城市出现日光盘
near_clean/BAC009S0766W0121.wav	3448e857c9e907de282c385d6f71be6f60508fd7b5881092df4c7a181c8e61de	BAC009S0766W0121	3.174	S0766	near_clean			实 现 数 字 化 整 合 营 销	实现数字化整合营销
```

## 读数③ `LICENSE` 头 5 行

`testdata/asr-baseline/LICENSE` ＝ **署名块 ＋ Apache-2.0 全文**，15,060 B / 275 行；正文自第 75 行起逐字取自 `https://www.apache.org/licenses/LICENSE-2.0.txt`（11,358 B，sha256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`）。

```text
1: AISHELL-1 ASR baseline subset -- LICENSE and attribution
2: =======================================================
3: 
4: Everything under this directory is a derived subset of two Apache-2.0 licensed
5: speech corpora published by OpenSLR. The full text of the Apache License,
```

## 读数④ `noisy_far` 的 SNR 实测

混噪目标逐句 **10.000 dB**（按功率比反解增益）。下面的数字**不是**目标值，而是**从交付的 wav 反测**：
`resid ＝ wav/0.501187 − 低通(上游干净句)`，`SNR ＝ 10log10(P(低通s)/P(resid))`。

| 噪声 | 句数 | min dB | mean dB | max dB |
|---|---|---|---|---|
| `cafe` | 14 | 9.964 | 9.994 | 10.011 |
| `car` | 13 | 9.952 | 9.989 | 9.998 |
| `white` | 13 | 11.114 | 11.271 | 11.299 |
| **全部** | **40** | **9.952** | **10.408** | **11.299** |

★ **白噪 11.3 而非 10.0 是滤波器性质，⛔ 混噪错误**：6 kHz 低通削掉宽带白噪能量的比例大于语音（语音能量集中在 4 kHz 以下）
⇒ 白噪那 13 句过完链后抬高约 1.3 dB；`cafe`／`car` 能量本就在低频 ⇒ 仍在 10.0。两种口径都在表里，没有藏。

## 读数⑤ `near_clean` 真的是 test split

**出处＝上游自己的路径前缀**，⛔ 靠猜号段：`data_aishell.tgz` 里**没有** `wav/test/` 目录，只有 400 枚 `wav/S####.tar.gz`；
**每枚内层 tar 的首段路径就是 split**（`train/`／`dev/`／`test/`）。逐枚只读内层首 8 KB（⛔ 解 15.6 GB）实测：

| split | 说话人数 | 说话人 |
|---|---|---|
| dev | 40 | S0724-S0763 |
| **test** | **20** | S0764-S0770, S0901-S0908, S0912-S0916 |
| train | 340 | 其余 |

⇒ `near_clean` 取的 S0764, S0765, S0766, S0767, S0768, S0769, S0770, S0901 **全部落在 test**；`noisy_far` 取的 S0902…S0912 同样在 test 内且与前者**不相交**。
test split 实测 **7,176 句 / 20 说话人 / 1,155,801,150 B**（＝通行口径的 AISHELL-1 测试集）；`near_clean` 是其中 60 句。

## 读数⑥ 总字节数 ＋ 与 9.2 MiB 的偏差

**100 句 wav ＝ 12,971,844 B ＝ 12.37 MiB ＝ 12.97 MB**；整目录（含 tsv／LICENSE／README／生成器）＝ 13,038,268 B ＝ 12.43 MiB。
音频合计 405.23 s ⇒ **32,011 B/s**，与 16k×16bit×1ch ＝ 32,000 B/s 的**口径相符**（多出的 ~11 B/s 即每文件 44 B 的 wav 头）⇒ **字节口径没出错**。

⚠ **比 cer-a1 §5 的 9.2 MiB 高 34%**，原因已定位，⛔ 不是口径错：cer-a1 的成本模型按 **3.0 s/句**；
AISHELL-1 test 实测均长 **5.0 s**（1,155,801,150 / 7,176），在票面要求的 **2-5 s** 窗口内仍均长 **3.99 s**（`noisy_far` 4.15 s）⇒ 60+40 必然落在 12.37 MiB。
⛔ 没有为贴数字而缩时长窗口——那等于按字节数挑语料。要回到 ~9.2 MiB 就把窗口收到 ~[2,4] s：`tools/make_asr_baseline.py` 的
`DUR_MIN`／`DUR_MAX` 两个常量改完重跑即可，但那是**改口径**，留给编排者拍。

## 读数⑦ ⛔ 碰 `internal/`／`cmd/` 的证明

`git show --name-only 6d120dca19449a33e258b8af7a5037792acc4350` ⇒ 命中 `internal/` 与 `cmd/` **0 条**（可复跑）：

```text
$ git show --name-only --pretty=format: 6d120dca19449a33e258b8af7a5037792acc4350 | wc -l
205
$ git show --name-only --pretty=format: 6d120dca19449a33e258b8af7a5037792acc4350 | grep -cE '^(internal|cmd)/'
0
```

全量 205 条**全部**在 `testdata/asr-baseline/` 之下：`near_clean/*.wav` 60 ＋ `near_clean/*.txt` 60 ＋ `noisy_far/*.wav` 40
＋ `noisy_far/*.txt` 40 ＋ `sources.tsv` ＋ `LICENSE` ＋ `README.md` ＋ `.gitattributes` ＋ `tools/make_asr_baseline.py` ＝ **205**。
`frontend/` `design/` `scripts/` `.github/` `docs/specs/` `docs/audit/` `docs/PLAN.md` 命中**均 0**；`internal/speech/**`
（并存腿 `15d-r1` 的写面）**一个字节未读未写**。

## 读数⑧ 仓里没有 15.6 GB 原始包

```text
$ du -sh testdata/asr-baseline
13M     testdata/asr-baseline
$ git ls-files testdata/asr-baseline | wc -l
205
$ git ls-files | grep -c '\.tgz$'
0
```

原始件**只建不删**地躺在仓外：`D:\tmp\wisp\20261026\15d2\data_aishell.tgz` ＝ **15,582,913,665 B**（`curl -sI` 的 `Content-Length` 与下载实测**同值**，
`curl -L` rc 0、无重试）。单文件最大 **159,116 B**（`BAC009S0906W0129.wav`）⇒ 远低于 GitHub 10 MiB 单文件上限。

## 复现与刻度

`tools/make_asr_baseline.py`：`SEED=20261026`，纯 numpy ＋ stdlib（**⛔ ffmpeg**、⛔ 外部滤波器）；每处随机只由 seed 与
排序后的上游 id 决定 ⇒ 无墙钟、无未播种 RNG、无字典序遍历。**重跑实测**：另跑一遍到临时目录，`sources.tsv` 与本提交
**逐字节同 sha256**（100 句全同）、wav 总字节同值 ⇒ 交付物可由件本身重放。

两枚上游包 sha256 **首次固定**（OpenSLR 全站无校验和，故由本腿自记，已写进 `README.md`）：

| 上游包 | B | sha256 |
|---|---|---|
| `SLR33/data_aishell.tgz` | 15,582,913,665 | `a4a0313cde0a933e0e01a451f77de0a23d6c942f4694af5bb7f40b9dc38143fe` |
| `SLR18/resource.tgz`（噪声源） | 24,813,708 | `5f10b11a86930d159250c979b2dccb26ac2527775b2aa628df08513af32d85ee` |

噪声件 3 枚（`resource/noise/{cafe,car,white}/*.wav`，16 kHz／mono／16-bit／300.0 s、**零重采样**）sha256 亦逐枚记在 `README.md`。
⛔ 下 THCHS-30 整包 6.45 GB —— 只用 24.8 MB 的 `resource.tgz`（SLR18 页面即单列该件，可单独下，已核 200）。

## 未做 / 待编排者拍

1. **9.2 → 12.37 MiB**：按票面 2-5 s 窗口取语料，未为贴数字缩窗口（见读数⑥）。若要求回到 ~9.2 MiB，改 `DUR_MAX=4.0`
   重跑即可 —— 属**改口径**，不由本腿自裁。
2. `cer-a1` §6 那句 `The data is free for academic use.` 与 Apache-2.0 标签的歧义：本腿**只留痕不裁定**，已逐字写进
   `LICENSE` 的 `Known upstream ambiguity` 段。
3. `.txt` 参考文本随 wav 同目录落盘（`cer-a1` §3 的落盘形状），`sources.tsv` 的 `text_norm` 与之逐字相同 ⇒ CER 计算器
   （并存腿 `15d-r1`）读哪一份都一致：**两份由同一生成器一次写出**，⛔ 人工二次录入。
4. ⛔ 未跑 `go test`／`d22scan`：本片写面不含 Go，且 `internal/speech` 有并存腿在飞。
5. 归一化加了**无损闸**：`text_norm` 必须等于原文去空格且每 token 全 CJK；会掉字的行**整句弃用**而非截短。全库 141,600 行里
   有 7 行含全角拉丁（`ｔ恤` `ａ级` `ｋ歌` `ｃ罗` `ｂ超`）⇒ 按此闸排除；交付 100 行 `text_raw` 实测**零**非 CJK 字符，⛔ 掉字。
