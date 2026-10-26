# 15d3-r1（写腿）实现证据 —— `sources.tsv` 改成 15d-r1 冻结的 5 列契约

**提交**：`c6ea9f66e`（`sources.tsv` ＋ `sources-extra.tsv` ＋ `README.md`，3 files）／本证据件随后一提交。⛔ push（写腿只 commit）。
**身份**：写腿（实现者）。★只改**自己这边**的形状去适配 15d-r1 的冻结契约，⛔ 反向改契约。
**零改动**：`internal/**`／`cmd/**`／`LICENSE`／100 个 `.wav`＋100 个 `.txt`／台账／票面 —— 一字节未动（读数⑤）。

## 0. 开工前核对（只读）

| 项 | 结果 |
|---|---|
| `git diff --cached --name-status` | **空**（上一轮 15d2-r1 的在飞索引已随它提交清空） |
| `sources.go:26` | `const manifestHeader = "file\tsha256\tspeaker\tcondition\torigin"` |
| `parseRow`（`sources.go:120-137`） | `len(strings.Split(text, "\t")) != 5` ⇒ `ErrManifestRow`；`f[0..4]` 依次落 `File/SHA256/Speaker/Condition/Origin` |
| `validateSource`（`sources.go:139-173`） | 另有：正斜杠、非绝对、`path.Clean` 原样、`.wav` 结尾、`condition ∈ {near_clean,noisy_far}`、**首段目录＝condition**、sha256 64 位小写、`speaker` 非空、**`origin` 非空** |
| `Origin` 的格式约束 | ⛔ **无**。`:169-171` 只判 `== ""`；`:46` 的 `AISHELL-1:test:S0002/BAC009S0002W0122` 是注释里的 `e.g.`，不是校验 |
| `15d-r1-impl.md:68` 冻结段 | 表头逐字 5 列；`speaker`/`origin` **非空**；**清单顺序＝报告顺序** |
| ⇒ 停手①（"解析还要求别的字段"） | ⛔ 未触发：`origin` 必填但**我们有数据**（语料来源） |

## 读数① 改后表头逐字（`cat -A`，取自提交内的 blob）

```
$ git show c6ea9f66e:testdata/asr-baseline/sources.tsv | head -1 | cat -A
file^Isha256^Ispeaker^Icondition^Iorigin$
```

**列数 ＝ 5**（每行 `^I` 计数 4 ⇒ 5 列；101/101 行唯一值 `{5}`）。行尾 `$` 无 `^M`（UTF-8 无 BOM、LF-only、末行有换行；与原件一致，`.gitattributes` 的 `*.tsv text eol=lf` 不变）。
被替换的 10 列表头（备查，逐字）：`file⇥sha256⇥source_utt_id⇥duration_s⇥speaker⇥condition⇥snr_db⇥noise_source⇥text_raw⇥text_norm`。

## 读数② 行数 ＝ 101（1 表头 ＋ 100 数据）

`sources.tsv` 101 行；`near_clean` 60 ＋ `noisy_far` 40 ＝ 100。★**顺序逐行照搬原文件**（60 近场在前、40 带噪在后），满足契约的「清单顺序 ＝ 报告顺序」；`file`/`sha256`/`speaker`/`condition` 四列的值**逐字未改**。

## 读数③ `LoadSet` 改前 vs 改后（`internal/speech/cer` 纯 Go，裸跑，⛔ staging）

同一条命令：`go test ./internal/speech/cer -run 'TestLoadSet|TestManifest|TestBaseline' -count=1 -v`

**改前 ⇒ 读不了（真跑，不是 SKIP）**
```
--- FAIL: TestBaselineCorpusIsCommittedAndTheGateCanRun
  baseline_test.go:145: LoadSet near_clean: …\testdata\asr-baseline\sources.tsv:
  cer: baseline manifest header is not the frozen one: line 1 is
  "file\tsha256\tsource_utt_id\t…\ttext_norm", want "file\tsha256\tspeaker\tcondition\torigin"
FAIL	github.com/CarlosShao/wisp/internal/speech/cer
```
★注意：`15d-r1` 那条「`testdata` 不存在 ⇒ 具名 SKIP」在 `testdata` 落地后**已转成真跑** —— 改前它是**红的**，⛔ 不是 SKIP。

**改后 ⇒ 能读**
```
--- PASS: TestBaselineCorpusIsCommittedAndTheGateCanRun (0.08s)
    baseline_test.go:145: near_clean: near_clean CER=100.00% gate<=6.00% FAIL (60 utterances, 660 ref chars, 660 errors)
    baseline_test.go:145: noisy_far: noisy_far CER=100.00% gate<=15.00% FAIL (40 utterances, 505 ref chars, 505 errors)
--- PASS: TestManifestFailuresAreNamed / TestLoadSetHappyPath / TestLoadSetFailuresAreNamed
ok  	github.com/CarlosShao/wisp/internal/speech/cer	0.146s
```
⇒ 读出 finding：`LoadSet` 读进 **60 near_clean ＋ 40 noisy_far**；`checkCommittedSet` 里 `VerifyManifest` 逐个 wav 重算 sha256 **全部对上**（对不上会先红）；两集 speaker **无交集**。★那两行 `CER=100.00% … FAIL` 是测试用 `t.Logf` 打的**示范**（"空转写必须 100% 失败"，证明 gate 有牙），⛔ 不是失败；用例本身 PASS。

## 读数④ `sources-extra.tsv`：101 行 × 6 列

列：`file`／`source_utt_id`／`duration_s`／`snr_db`／`noise_source`／`text_raw`（该头行 6 列）。101 行 ＝ 1 表头 ＋ 100 数据，**行序与 `sources.tsv` 完全相同**，按 `file` 一一对应。`near_clean` 行的 `snr_db`/`noise_source` **保留原空值**（⛔ 未编造占位）。
★`text_norm` **不搬**：100/100 行与同名 `.txt` 逐字相同（`TrimSpace` 后 `-ceq` 比较全过）⇒ ⛔ 丢信息，停手②未触发。

## 读数⑤ `.wav`/`.txt` 零改动的证明

```
$ git show --numstat --pretty=format: c6ea9f66e
24	7	testdata/asr-baseline/README.md
101	0	testdata/asr-baseline/sources-extra.tsv
101	101	testdata/asr-baseline/sources.tsv
```
⇒ **提交里只有这 3 个文件**；200 个数据件（100 `.wav` ＋ 100 `.txt`）在 `--numstat` 里**一条都没有**，故其 `-` 行 ＝ **0**（实测筛选计数 ＝ 0）。
更强的一把尺（逐 blob 比对）：`git ls-tree -r c6ea9f66e^ / c6ea9f66e -- testdata/asr-baseline/{near_clean,noisy_far}` ⇒ 两侧各 **200** 条、**逐条 blob hash 全同**（`Compare-Object` 空）⇒ 200 个数据件**字节级未变**。
`git show --name-only --pretty=format: c6ea9f66e | grep -cE '^(internal|cmd)/'` ⇒ **0**。`LICENSE` 未现于提交。

## 读数⑥ `origin` 填了什么、为什么

| condition | `origin` |
|---|---|
| `near_clean`（60） | `aishell-1` |
| `noisy_far`（40） | `aishell-1+thchs30-noise` |

理由：
1. 两种条件的**语句都取自 AISHELL-1**（test split），`noisy_far` 只是**加了 THCHS-30 的噪声**（`cafe`/`car`/`white`，各 14/13/13，见 `sources-extra.tsv` 的 `noise_source`；三件都来自 SLR18 `resource.tgz`，见同目录 `README.md:150-156`）。★所以 `noisy_far` 若也写 `aishell-1`，会把"这不是原始语音、是加噪模拟件"这条**真实来源信息**抹掉，而那正是 `origin` 该承载的东西。
2. ★**小写** `aishell-1`：`sources.go:46` 的 `AISHELL-1:…` 只是注释示例、非校验（`:169` 只判非空）；小写与本文件既有取值风格一致（`near_clean`／`cafe`／`car`／`white`），也与本腿派单逐字一致。
3. ★**⛔ 把 utt id 塞进 `origin`**：`file` 的 stem **就是** AISHELL-1 utt id（`near_clean/BAC009S0764W0121.wav`，100/100 行 `stem == source_utt_id`），再写一遍是冗余；utt id 已完整留在 `sources-extra.tsv` 的 `source_utt_id`。
4. ⛔ 未动 `condition` 里已冻结的目录名（首段＝condition，`validateSource:159` 会查）。

## 读数⑦ 数据一致性：抽 5 句两件对核

抽样规则（可复跑）：`Get-Random -SetSeed 20261026`，`near_clean` 区抽 3 行、`noisy_far` 区抽 2 行，按文件行号升序报 ⇒ 行 **2 / 5 / 6 / 64 / 65**。每行核 6 项，**全 True**：

| 行 | `file` | `sha256`（重算 wav 后比对） | `speaker` | `condition` | `origin` | extra：`source_utt_id`/`duration_s`/`snr_db`/`noise_source` |
|---|---|---|---|---|---|---|
| 2 | `near_clean/BAC009S0764W0121.wav` | `46dbc998…60dce50` ✔ | S0764 | near_clean | `aishell-1` | `BAC009S0764W0121`/4.204/∅/∅ |
| 5 | `near_clean/BAC009S0767W0121.wav` | `cbbc2847…57eb36` ✔ | S0767 | near_clean | `aishell-1` | `BAC009S0767W0121`/4.588/∅/∅ |
| 6 | `near_clean/BAC009S0768W0122.wav` | `738850b3…8ffe7` ✔ | S0768 | near_clean | `aishell-1` | `BAC009S0768W0122`/3.915/∅/∅ |
| 64 | `noisy_far/BAC009S0904W0121.wav` | `1ea9ef96…88fbcd` ✔ | S0904 | noisy_far | `aishell-1+thchs30-noise` | `BAC009S0904W0121`/3.056/11.2916/`white` |
| 65 | `noisy_far/BAC009S0905W0121.wav` | `96191348…a52f5` ✔ | S0905 | noisy_far | `aishell-1+thchs30-noise` | `BAC009S0905W0121`/3.601/10.0039/`cafe` |

每行 6 项判据：①两件 `file` 相同 ②`sha256` 用 `Get-FileHash` 从盘上 wav **重算相等** ③`stem == source_utt_id` ④`speaker == utt_id[6:11]` ⑤`condition == file 首段` ⑥`.txt` 去空格后 == `text_raw`（且全局 100/100 行 `text_norm == .txt`）。★另外 `origin` 与 `noise_source` 两向自洽：`snr_db`/`noise_source` 非空 ⇔ `condition == noisy_far` ⇔ `origin` 带 `+thchs30-noise`。

## 读数⑧ 格式甲乙两尺（★已知 3 条既存假红，⛔ 本腿责任）

| 尺 | 命令 | 读数 | 归属 |
|---|---|---|---|
| 甲 | `gofumpt -l internal cmd tools scripts`（＝`ci.yml` 的 gofmt 步同口径） | **3 条**：`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`cmd/wisp/models.go` | ★**3 条既存**，与 `gofmt -l` 同 3 条；本腿提交 `.go` 命中 **0** ⇒ 非本腿 |
| 乙 | `go vet ./...` | **rc=0，无输出（绿）** | — |

⇒ 两尺的读数与本腿改动**无交集**：`git show --name-only c6ea9f66e` 的 `internal|cmd` 命中 ＝ 0。

## 停手条件复核（三条全未触发）

①`origin` 必填但无数据 ⇒ 未触发（有语料来源，且格式无约束）。②5 列会丢信息 ⇒ 未触发（6 列另存 ＋ `text_norm` 由 `.txt` 承载，100/100 核过）。③改完仍读不了 ⇒ 未触发（`ok`，60＋40 句读入、digest 全对）。

## 遗留（发现但⛔ 本腿修，报给编排者）

1. ★`testdata/asr-baseline/tools/make_asr_baseline.py` **仍生成 10 列 `sources.tsv`**（`tools/` 不在本腿写面）。⇒ 若有人照 `README.md` 的「Regenerating」重跑，会**把形状打回 10 列**、再撞 `ErrManifestHeader`。建议后续腿把生成器改成写两件（`sources.tsv` 5 列 ＋ `sources-extra.tsv` 6 列）并把 `origin` 一并生成。
2. ★`README.md` 的目录树块（第 17 行 `sources.tsv   100 data rows`）**未改**：该句仍为真，且它属"文件清单"不属"列说明" ⇒ 守本腿"只改列说明那几句"的边界；新件已在同一节的「`sources-extra.tsv`」小节里逐列说明。
3. ★`15d2-r1-impl.md` 里对 10 列的描述是**历史记录**，⛔ 回改（只追加不删）。
