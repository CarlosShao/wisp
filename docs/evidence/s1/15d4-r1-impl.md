# 15d4-r1（写腿）实现证据 —— 生成器补上 5 列＋6 列，重跑不再打回 10 列

**提交**：`3ccb8c19880ebc86f9104a4949259665f219ea05`（2 files：`tools/make_asr_baseline.py` ＋ `README.md`）／本证据件随后一提交。⛔ push（写腿只 commit）。
**身份**：写腿（实现者）。★本腿只修**生成器**（脚本），⛔ 动交付（数据）。
**起因**（`15d3-r1` 遗留①）：脚本仍一次性写 10 列 `sources.tsv`，交付件却是 5 列（`15d-r1` 冻结于 `internal/speech/cer/sources.go:26`）＋ 6 列 `sources-extra.tsv`（`15d3-r1` 加）⇒ **生成器与交付不一致，照 README「Regenerating」重跑会把形状打回 10 列、再撞 `ErrManifestHeader`**。

## 0. 开工前核对（只读）

| 项 | 结果 |
|---|---|
| `git diff --cached --name-status` | **空** ⇒ 未停手（工作树另有他腿的 `M .scratch/**`、`M docs/reports/pending-and-issues.md` 与一批 `??`，**本腿未碰也未 add**） |
| 读 `make_asr_baseline.py` 全文 | 354 行；`:319-330` 一次性拼 10 列：`file,sha256,source_utt_id,duration_s,speaker,condition,snr_db,noise_source,text_raw,text_norm` |
| 两份 tsv 现有形状（只读） | `sources.tsv` 101 行/5 列；`sources-extra.tsv` 101 行/6 列；行序同（`near_clean` 先、`noisy_far` 后，各自 `take()` 轮转序） |
| `D:\tmp\wisp\20261026\15d2\` | **包还在**：`data_aishell.tgz` 15,582,913,665 B、`resource.tgz` 24,813,708 B，解包树 `aishell/{transcript,wav/test}`、`resource/resource/noise/{cafe,car,white}` 也在 ⇒ **走真重跑**，⛔ 走"从仓内重算 tsv"退路 |

环境：`python` = CPython **3.11.15** ＋ numpy **2.4.3**（`py -3` 是 3.12.10/numpy 1.26.4，未用）；`GOTMPDIR` 本就未设（**无 unset 动作**，TMPDIR/TEMP/TMP 未动）；临时件全在 `D:\tmp\wisp\20261026\15d4\`（只建不删）。

## ① ★核心：重跑后两件 tsv 与提交件**逐字节相同**

做法：★**⛔ 没在仓内跑生成器**。空目录 `D:\tmp\wisp\20261026\15d4\rerun1` 重跑改后脚本
（`--aishell D:\tmp\wisp\20261026\15d2\aishell --noise D:\tmp\wisp\20261026\15d2\resource\resource\noise`），
`rc=0`，与仓内件比：不只看 sha256，还 `[System.IO.File]::ReadAllBytes` ＋ `SequenceEqual` 逐字节比。

| 件 | 仓内（＝提交件） | 仓外重跑 | 逐字节相同 |
|---|---|---|---|
| `sources.tsv` | 12,917 B ／ `54b50a499e5b33c0e722f4a1995e02023f3cc25daec6aefa48dc6089f64db537` | 12,917 B ／ **同左** | ★**True** |
| `sources-extra.tsv` | 10,799 B ／ `531677f225b7cc24f5530505fc058458a1bccd356eb96ef028bc27f91b2686dd` | 10,799 B ／ **同左** | ★**True** |

**"已生成"状态上重跑不改字节**（要求②）：对**已跑完的 `rerun1` 树原地再跑一次**，`rc=0`，两件 `before == after`（逐字节比 True）；而 `rerun1` 树在这 202 个件上与仓内目录**逐字节相同**（见④）⇒ 这条即"重跑仓内也不改字节"的等价证明。
⇒ `write_tsv()` 固定列序、LF、无引号无 BOM 无尾随空白，行列序由调用方固定，**重跑幂等**。

## ② 改后脚本产 5 列 ＋ 6 列（表头逐字，取自**提交内的 blob**）

```
$ git show HEAD:testdata/asr-baseline/sources.tsv       | head -1 | cat -A
file^Isha256^Ispeaker^Icondition^Iorigin$
$ git show HEAD:testdata/asr-baseline/sources-extra.tsv | head -1 | cat -A
file^Isource_utt_id^Iduration_s^Isnr_db^Inoise_source^Itext_raw$
$ ... | awk -F"\t" '{print NF}' | sort -u     # sources.tsv ⇒ 5 ;  sources-extra.tsv ⇒ 6
```

列数 ＝ **5** 与 **6**（101/101 行每行字段数唯一值 `{5}`／`{6}`）；行尾 `$` 无 `^M`（UTF-8 无 BOM、LF-only、末行有换行）。
硬校验：`sources.tsv` 首行**逐字等于** `sources.go:26` 的 `manifestHeader`（Go 字面量反转义 `\t` 后 `-ceq` 比较 **True**；重跑件 **True**）。
行数：`sources.tsv` `lines=101 data-rows=100 cols=5`；`sources-extra.tsv` `lines=101 data-rows=100 cols=6`（重跑件同）。`rc=0` 时脚本自报 `rows in sources.tsv: 100 (5 cols) and sources-extra.tsv: 100 (6 cols)`。

## ③ `origin` 两值的生成规则（代码那几行，逐字）

`make_asr_baseline.py:73-76`
```python
# origin column, one value per condition. near_clean ships AISHELL-1 audio
# untouched; noisy_far is AISHELL-1 speech plus additive THCHS-30 noise, and the
# cell names both because the noise is the part that is not AISHELL-1.
ORIGIN = {"near_clean": "aishell-1", "noisy_far": "aishell-1+thchs30-noise"}
```
`make_asr_baseline.py:357-363`（写 `sources.tsv`；`r["cond"]` 只可能是这两个键，缺键即 `KeyError` 崩，⛔ 静默空串）
```python
    write_tsv(
        os.path.join(args.out, "sources.tsv"),
        ["file", "sha256", "speaker", "condition", "origin"],
        [[r["file"], r["sha"], r["speaker"], r["cond"], ORIGIN[r["cond"]]]
         for r in rows],
    )
```
实测落盘：`near_clean/*` 60 行 `aishell-1`；`noisy_far/*` 40 行 `aishell-1+thchs30-noise`。
`text_norm` **不进任何 tsv**（旧 10 列的该列已删；参考文本仍写进同名 `<utt>.txt`，现状不变）。

## ④ 200 个数据件**零改动**的证明

- `git show --numstat`（`3ccb8c198`，逐字两行）：`12 / 1  testdata/asr-baseline/README.md`、`52 / 13  testdata/asr-baseline/tools/make_asr_baseline.py` ⇒ `.wav`/`.txt` 命中 **0**；`.go` 命中 **0**；无第三件。
- `git show --name-only`：`^(internal|cmd)/` 命中 **0**；`\.(wav|txt)$` 命中 **0**。
- `git status --porcelain -- testdata/asr-baseline/{near_clean,noisy_far}`：**空**（200 个数据件工作树一字节未动）；两份 tsv 也**不被列出**（＝未修改）。
- 交叉证（重跑 vs 仓内）：100 `.wav` ＋ 100 `.txt` 逐文件 sha256 比，`compared=200 mismatches=0`。

## ⑤ `TestBaselineCorpusIsCommittedAndTheGateCanRun` 仍 `rc=0`

```
$ go test ./internal/speech/cer -run TestBaselineCorpusIsCommittedAndTheGateCanRun -v -count=1   # rc=0
=== RUN   TestBaselineCorpusIsCommittedAndTheGateCanRun
--- PASS: TestBaselineCorpusIsCommittedAndTheGateCanRun (0.07s)
PASS
ok  	github.com/CarlosShao/wisp/internal/speech/cer	0.102s
```
日志里另有 `baseline_test.go:145` 两条 `CER=100.00% … FAIL` —— 那是 `checkCommittedSet` 故意喂"什么都不听"的 `empty` 假说测**门有牙**，不是失败。
这条 PASS 同时证明：`LoadSet`（⇒ `ParseManifest` 认了冻结表头与 100 行）、`VerifyManifest`（⇒ **盘上 100 个 wav 的 sha256 与清单全对 ⇒ 数据确实没动**）、两集 speaker 不相交。

## ⑥ `README.md` 改了哪几句（逐行 diff，仅「Regenerating」小节）

```
@@ -195,7 +195,8 @@ for s in S0764 ... S0916; do
-# 3. generate (needs python3 + numpy; writes the wavs, the .txt files and sources.tsv)
+# 3. generate (needs python3 + numpy; writes the wavs, the .txt files,
+#    sources.tsv and sources-extra.tsv)
@@ -203,6 +204,16 @@
+Step 3 writes both tables: `sources.tsv` (5 columns, the frozen manifest the CER
+reader parses) and `sources-extra.tsv` (6 columns, the derivation columns that
+reader never sees). Both are written with a fixed column order and a fixed row
+order, so re-running step 3 over an already-generated tree rewrites identical
+bytes: `sources.tsv` stays at sha256
+`54b50a499e5b33c0e722f4a1995e02023f3cc25daec6aefa48dc6089f64db537` and
+`sources-extra.tsv` at
+`531677f225b7cc24f5530505fc058458a1bccd356eb96ef028bc27f91b2686dd`. The
+reference text is not a column in either file: it is the sibling `<utt>.txt`.
```
`git diff --numstat README.md`：`12 1`（+12/−1），两处 hunk 全在 `## Regenerating` 内，⛔ 碰别的小节。

## ⑦ 格式：`gofmt`/`gofumpt` 与本次改动**无交集**（如实报）

本腿改 `.py` ＋ `.md`：`git show --name-only HEAD | grep -c '\.go$'` ⇒ **0**。`gofmt -l` 喂这两个件 ⇒ `rc=2`、`illegal character U+0023 '#'`（`gofmt` 只吃 `.go`，直接拒非 Go 件 ⇒ 射程确实不覆盖本腿）；`gofumpt` **本机未安装**（无输出可报）。Python 侧只 `ast.parse` 语法核（`AST OK 393 lines`）＋两次真跑 `rc=0`；⛔ 跑 `py_compile`（会在 `tools/` 落 `__pycache__`）。

## 硬边界核对（四条）

1. ⛔ 200 个数据件：**未在仓内跑**、`git status` 干净、`numstat` 命中 0（④）。
2. ⛔ 两份 tsv：**一字节未动** —— 不是"改了之后恰好相同"，是 `git status` 根本不列它们，且重跑件与它们逐字节相同（①）。
3. ⛔ `internal/**`／`cmd/**`：命中 0（④），`sources.go` 未碰。
4. ⛔ `LICENSE`／`frontend/**`／`design/**`／`scripts/**`／`.github/**`：未碰（`name-only` 只两件，均在 `testdata/asr-baseline/`）。

## 遗留 / 未定

- ① `15d3-r1` 遗留①（生成器打回 10 列）**本腿关闭**：重跑产出的两件与提交件逐字节相同；未定义即停三个条件**均未触发**（包在、重跑同字节、无需动 tsv）。
- ② ⛔ 未改 `README.md` 第 17 行的目录清单（`sources.tsv   100 data rows, one per .wav`）—— 它在该小节之外，⛔ 越写面；`sources-extra.tsv` 已由 `## sources-extra.tsv` 单独成节讲清，故不算缺口。
