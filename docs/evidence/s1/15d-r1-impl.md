# 票 15d CER 计算器腿 `15d-r1` —— 落地件（写腿）

> **尺**：`HEAD=37aadd1ec`（`dev`），现跑 2026-10-26 16:41–17:2x +08。★**八项读数逐条在 §1–§8**，全为现跑，⛔ 转述。
> **代码 commit**：`12be96ad3`（8 文件、1776 行；`git show --name-only` 里 `cmd/` 命中 **0**、顶层目录只有 `internal`）。
> **写面**：`internal/speech/cer/**`（新建 8 文件）＋ 本文件。★`cmd/**`、`internal/**` 其余、`scripts/**`、`docs/specs/**`、票面、台账、`testdata/**` **零字节**。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\15dr1\{logs,mut,probe,bin}`。★突变全走仓外件 ＋ `-overlay` ⇒ **无工作树突变窗口**。

## 0. 一句话 ＋ 三个前提

`internal/speech` 下新建子包 `cer`：**rune 级 `(S+D+I)/N`** ＋ **四条写死归一化** ＋ **`sources.tsv` 严格读法** ＋ **6%/15% 双门禁**。纯 Go、无 cgo、零真实语料 ⇒ **裸 `go test ./internal/speech/cer` 即绿**（父包裸跑仍要 sherpa DLL，§7）。

| 开工三步 | 现跑读数 |
|---|---|
| ① 索引 | `git diff --cached --name-status` ⇒ **空**（非空才停手） |
| ② 语料前提 | `ls testdata/asr-baseline/` ⇒ **不存在**；`testdata/` 整个目录也不存在 ⇒ 今天没有语料（§6） |
| ③ 读现有包 | `internal/speech` 11 文件：`doc.go`／`endpoint.go`／`vad.go`／`asr.go`／`transcriber.go` ＋ `endpoint_test`／`fixtures_test`／`vad_test`／`asr_test`／`transcriber_test`／`cgo_boundary_windows_test`（后三个按符号表＋风格抽样，本片是子包、不共享符号）★期间 `15-c` 落地（`a736ad0ea`）⇒ 它的 4 个文件在盘上 |

★**为什么是子包而不是 `package speech`**：`15-c` 在飞时 `asr.go`/`transcriber.go` 未跟踪、同包编译；子包＝**零符号冲突、零编译耦合**，且 CER 不许 cgo（门禁不该为了算分母去加载模型）。

## 1. ★★★ 读数①：字级编辑距离的正确性（ALU 表逐格验）

**手算小例（可逐字复核）**：`今天天气` vs `今天气` ⇒ 参考 4 字、删 1 字：

```
SCORE "今天天气" vs "今天气": refs=4 S=0 D=1 I=0 CER=0.2500
报告行 CER=25.00% (1/4) S=0 D=1 I=0
```

**四类各一组（逐格，含 S/D/I 分裂）**：

| 类 | 参考 | 假设 | S | D | I | CER |
|---|---|---|---|---|---|---|
| 相等 | 今天天气 | 今天天气 | 0 | 0 | 0 | 0.0000 |
| 替换 | 今天天气 | 今天天汽 | 1 | 0 | 0 | 0.2500 |
| 删除 | 今天天气 | 今天气 | 0 | 1 | 0 | 0.2500 |
| 插入 | 今天气 | 今天天气 | 0 | 0 | 1 | 0.3333 |
| 两替换 | 甲乙丙 | 丁戊丙 | 2 | 0 | 0 | 0.6667 |

★★**网格逐格验（225 格）**：字母表 `{甲,乙}`、长度 0–3 ⇒ 15 参考 × 15 假设。每格 **4 条断言**：①DP 的 `Errors()` == **独立写的自顶向下 oracle**（不共享表/滚动行/分裂记账）②`S+D+I == cost` 且 `|n-m| ≤ d ≤ max(n,m)` ③`CER == errors/N`（恒等）④空参考 ⇒ `ErrEmptyReference`（⛔ 0%）。全网格另验 **对称性 225 格** ＋ **三角不等式 3375 个三元组**。

```
cost matrix (rows=ref, cols=hyp) over ["", "甲", "甲甲", "甲甲甲", "甲甲乙", "甲乙", "甲乙甲", "甲乙乙"]
            ""   0   1   2   3   3   2   3   3
           "甲"   1   0   1   2   2   1   2   2
          "甲甲"   2   1   0   1   1   1   1   2
         "甲甲甲"   3   2   1   0   1   2   1   2
        "甲甲乙"   3   2   1   1   0   1   2   1
          "甲乙"   2   1   1   2   1   0   1   1
         "甲乙甲"   3   2   1   1   2   1   0   1
         "甲乙乙"   3   2   2   2   1   1   1   0
```
★**并列最小解定序**：替换优先、删次之、插最后 ⇒ 分裂可复现；**总数在任何定序下都等于 Levenshtein 距离**，门禁读总数。

## 2. ★★★ 读数②：归一化四条各一发（＋两条「必须算错」）

| 规则 | 输入 → `Normalize` | 判据 |
|---|---|---|
| 全角→半角 | `ＡＢＣ１２３` → `ABC123`；`２０２６年，今天天气` vs `2026年今天天气` ⇒ **CER 0** | `TestRuleTwoFoldsFullWidth` |
| 去空白标点 | `今天，天气！` → `今天天气`；` hello,\tworld! ` → `helloworld` | `TestRuleOneDropsWhitespaceAndPunctuation` |
| **繁简不转换**（必须算错） | `汉语` vs `漢語` ⇒ **S=2 / CER 1.0000**（⛔ 0） | `TestRuleThreeNeverConvertsSimplifiedAndTraditional` |
| **数字不做 ITN**（必须算错） | `2026` vs `二零二六` ⇒ **S=4 / CER 1.0000**；反向 `二零二六` vs `2026` 亦 **1.0000**（分母 4） | `TestRuleFourDigitsAreNotSpelledOut` |

★**第 5 条（`cer-a1:100` 也写了，本片并入）**：`<UNK>`／`<sil>`／`<noise>` 类标记**整枚剔除**（`今天<UNK>天气` → `今天天气`）；未闭合的 `<UNK天气` **不猜**，字母保留（`今天UNK天气`）。★**大小写不折叠**（规则 4 说原样；`near_clean` 按 `cer-a1:107` 不含英文/数字 ⇒ 动不了门禁）。
★**复合一发**：`２０２６年，今天<UNK>天气！` vs `2026年今天天气` ⇒ 两边都归一化成 9 字、**CER 0** ⇒ 四条同时在场也成立。

## 3. ★★ 读数③：`sources.tsv` 缺失/空/坏，**全是具名错误**

★★**本片替基线集腿冻结的契约**（`testdata/asr-baseline/sources.tsv`，⛔ 推测）：表头**逐字** `file⇥sha256⇥speaker⇥condition⇥origin`；`file` 相对 root、正斜杠、以 `.wav` 结尾、**首段必须等于 `condition`**（`SPEC-01:86-87` 冻结的目录名）；`sha256` **64 位小写** hex；`condition ∈ {near_clean,noisy_far}`；`speaker`/`origin` 非空；**参考文本＝同名 `.txt`**（`cer-a1:101` 的落盘形状）；**清单顺序＝报告顺序**。

| 情形 | 具名错误（现跑逐字，路径省略） |
|---|---|
| **缺失** | `cer: baseline manifest is missing: open …\sources.tsv: The system cannot find the file specified. (the committed baseline set is a separate leg; an absent set must never be read as an empty, passing one)` ★同时命中 `ErrManifestMissing` **与** `fs.ErrNotExist` |
| 空文件 | `cer: baseline manifest has no data rows: no header line "file\tsha256…" and no data` |
| 仅表头 | `cer: baseline manifest has no data rows: header "file\tsha256…" then nothing` |
| 表头不对 | `cer: baseline manifest header is not the frozen one: line 1 is "file\tspeaker", want "file\tsha256…"` |
| 该 condition 0 行 | `cer: no utterances to score: …\sources.tsv lists no "near_clean" rows` |
| 缺 wav | `cer: baseline audio file is missing: …\near_clean\0001.wav (manifest lists near_clean/0001.wav)` |
| 缺 `.txt` | `cer: baseline reference text is missing: …\near_clean\0001.txt (the reference is the sibling .txt of …)` |
| `.txt` 空 | `cer: reference text is empty after normalization: …\near_clean\0001.txt` |

★另有 12 类**坏法**各一发（列数≠5／大写 hash／63 位 hash／`condition` 与目录不符／反斜杠／`../`／绝对路径／非 `.wav`／未知 condition／空 speaker／空 origin／重复行），分别 ⇒ `ErrManifestRow`／`ErrManifestDuplicate`。★★**⛔ 任何一条返回 `(空集, nil)`** —— 这就是「静默当 0 句」的封堵口。

## 4. ★★ 读数④：反向判据 —— 三读数 ＋ 还原核

三发变异全在**仓外**（`mut\{A,B,C}`）＋ `-overlay` 换入编译：

| 发 | 变异 | rc | 红的断言（逐字节选） |
|---|---|---|---|
| **A** | `Normalize` ⇒ **identity**（归一化关掉） | **1** | **7 条红**：`norm_test.go:21`「`Normalize("今天，天气！")="今天，天气！", want "今天天气"`」／`:34` 全角／`:107` 标记／`cer_test.go:236`「`normalised-equal pair scored CER=33.33% (2/6), want 0`」／`cer_test.go:362` 标点-only 参考不再报 `ErrEmptyReference`／`sources_test.go:175`「`scoring the loaded reference: CER=20.00% (1/5)`」 |
| **B** | 门禁恒过（`v.Pass = true`） | **1** | **3 条红**：`cer_test.go:273`「`near_clean 601/10000: Pass=true, want false`」／`cer_test.go:334`「`30% passed the 6% gate`」／`baseline_test.go:165`「`empty transcript must fail with CER 100%, got … PASS`」 |
| **C** | **缺失 ⇒ 空集**（`ReadManifest` 返回 `[]Source{},nil` ＋ 删 `LoadSet` 的 0 行守卫） | **1** | **4 条红**：`sources_test.go:129`「`ReadManifest on an empty directory: got <nil>`」／`:236`「`no rows for the condition: got <nil>`」／★`baseline_test.go:138`「`absent baseline set: LoadSet returned <nil>, want ErrManifestMissing`」（**连 skip 那条也变红**）／`:176`「`absent corpus: LoadSet returned 0 utterances and <nil>`」 |
| — | **原样绿**（变异前，无 overlay） | **0** | **24 案：23 PASS / 1 具名 SKIP / 0 FAIL** |
| **D** | **还原核**（最后一发 **C 之后**，无 overlay 重编重跑） | **0** | 23 PASS / 1 SKIP；★盘上变异签名 grep（`v.Pass = true`／`normalizeUnused`／`return []Source{}, nil`）＝**0**；8 文件 sha256 与提交一致（`cer.go 25D91A114BD4`、`gate.go 48C2161DEE07`、`norm.go 662F8A5EC8E0`、`sources.go 003700411925`…） |

★**还原核时机**：在**最后一发变异（C）之后**单独跑，不在 A/B 之间 —— A/B 之间跑只能证明 overlay 生效，不能证明树没被写过。

## 5. ★★ 读数⑤：正控 —— 参考＝假设 ⇒ CER ＝ 0（⛔ 差一字符）

`今天天气`／`人工智能`／`甲乙丙`／`ABC123` 四枚 `ref == hyp` ⇒ **`S=D=I=0` 且 `CER == 0`（float 精确等于 0，不是 ≈0）**；归一化后相等（`今天，天气！` vs `今天 天气`）亦 **0**。
★**同一发里的反向正控**：`今天天气` vs `今天天汽` ⇒ `S=1`、**CER 0.25 ≠ 0** —— 防的是「`==` 坏成永远 0」。★变异 A 恰好把这条打成红（§4），说明正控不是摆设。

## 6. ★★★ 读数⑥：`testdata` 不存在时（今天就是这样）⛔ 崩、⛔ 静默绿

- **具名 SKIP（逐字）**：`--- SKIP: TestBaselineCorpusIsCommittedAndTheGateCanRun` —「`baseline corpus absent (…\testdata\asr-baseline\sources.tsv): the CER gate is NOT verified by this run. The committed sets are a separate leg - AISHELL-1, 60 near_clean + 40 noisy_far, see docs/evidence/s1/cer-a1-research.md. SKIP, not pass: a green over 0 utterances would be a lie.`」
- **具名报错**：同一条测试在 skip **之前**先断言 `LoadSet(...) == ErrManifestMissing`（§3）—— 缺语料**先**得到一个具名错误，**然后**才 skip。
- **★★★「不会变成永远绿」的三重证明**：①skip 之前那条断言：接缝若改成「缺失⇒空集」，skip 立刻变**红**（变异 C 实测 `baseline_test.go:138`）；②`TestAbsentCorpusIsANamedErrorNotAnEmptyPass` **无条件跑**（今天就在跑），断言缺失是 `ErrManifestMissing` ＋ `fs.ErrNotExist`，且**任何零分母的报告/集合都拿不到 `Pass=true`**；③`Set{}.Score`／`Report{}.Verdict`／`Aggregate{}.Judge` 三条路全部具名报错（`ErrNoUtterances`）。
- ★★**给下游的硬事实**：`tools/d22scan/runtests.sh:224-235` 把**任何 `--- SKIP` 判红**。⇒ `15-d` 的 CI 接线一旦把本包交给 `runtests.sh`，这条 skip **会立刻变红** —— 这是**故意的**：基线集必须在接线之前落地，⛔ 不是把 skip 删掉。

## 7. ★ 读数⑦：四数（`internal/speech` 改前／改后）

| 尺 | 命令 | **改前** | **改后** |
|---|---|---|---|
| 裸跑 | `go test -count=1 ./internal/speech` | rc **1**、`exit status 0xc0000135`（sherpa DLL 缺失，`15b` §9.1 既存） | rc **1**、**逐字同因** ⇒ 未引入新破法；★**新包裸跑 rc 0** |
| staging（DLL 就位） | `go test -c` ＋ 拷贝 DLL ＋ `speech.test.exe -test.v=true` | **28 PASS / 0 FAIL / 0 SKIP**，rc **0** | **28 PASS / 0 FAIL / 0 SKIP**，rc **0**（**零回归**） |
| 新包 | `go test -count=1 ./internal/speech/cer` | — | rc **0**，**24 案：23 PASS / 1 具名 SKIP** |
| 静态 | `go vet ./internal/speech/...`；`sh scripts/d22scan.sh` | — | vet rc **0**；D22 rc **0** `clean - no D22 ban violations`（`internal/` ban#8 扫 **664** 个 `.go`） |

★**门禁 rc**：本片**没有 CLI**（`scripts/**` 不在写面）⇒ rc 取测试进程 —— **原样 0**、变异 B（门禁恒过）**1**、变异 A／C **1**。★`asr-cer` 的 CI job（`SPEC-11:128`）与「把本包交给 `runtests.sh`」属 `15-d` 的**接线那半**，⛔ 本片未做。

## 8. ★ 读数⑧：格式甲乙两尺（⚠ 3 条既存假红，⛔ 本片责任）

| 尺 | 命令 | 输出 |
|---|---|---|
| 甲 `gofmt` | `gofmt -l internal cmd tools scripts frontend` | **3 行**：`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`、`cmd\wisp\models.go` |
| 乙 `gofumpt`（`/d/work/base/gopath/bin/gofumpt.exe`） | `gofumpt -l internal cmd tools scripts frontend` | **同上 3 行，逐字相同** |
| 归属 | 对两尺输出按 `[\\/]` 判 `internal[\\/]speech` | **0 命中** ⇒ 3 条全是既存 CRLF/`autocrlf` 假红，与本片无关 |
| 本片 | `gofmt -l internal/speech`；`gofumpt -l internal/speech` | **两尺都空** |

## 9. 边界自陈 ＋ 未触发停手的核对

- ⛔ **下载／网络**：**0 次**（无 `curl`/`wget`/`http`/`os/exec`；`go.mod` 未变）。⛔ **真实 wav**：**0**（夹具是合成文本对 ＋ `t.TempDir()` 里的非音频字节）。⛔ **`testdata/**`**：0 字节。
- ⛔ **`cmd/**`**：`git show --name-only 12be96ad3` 里 `cmd/` 命中 **0**；`git diff --name-only 12be96ad3~1 12be96ad3` 中 `internal/speech/cer/` 之外的路径 **0 条**。
- 未碰 `15-a`／`15-c`／`15-e`；`15-c` 的 4 个文件（`a736ad0ea` 已提交）不在我的 add 列表。⛔ 未加配置键、⛔ 未动 `DEFERRED(engines)`。
- **未做（明写）**：CI 接线（`scripts/**` 不在写面）· 基线集与真实音频（另一条腿）· 真 ASR 假设与端到端门禁（`15-c`/`15-e`）。
- ★**停手条件核对**：①`SPEC-10:57-62` 只写「近场清洁 ≤6% · 带噪/远场 ≤15%」与「两个 wav 集必须提交进 `testdata/asr-baseline/`」，**未写公式** ⇒ `(S+D+I)/N` 与契约**无冲突**，未触发停手①；②无需动 `testdata/**` ⇒ 未触发②；③无需语料即可验计算器 ⇒ 未触发③。
