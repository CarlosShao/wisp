# CER 门禁 CI 接线腿 `cerci-r1` —— 落地件（写腿）

> **尺**：开工 `HEAD=1f1b86a5d5`（`dev`），现跑 2026-10-26 18:10–18:2x +08。★八项读数逐条在 §1–§8，全为现跑，⛔ 转述。
> **代码提交**：`8d6e150f8`（2 文件、+84 行：`ci.yml` +56、`scripts/portable-tests.sh` +28）／本证据件随后一提交。⛔ push。
> **写面**：`ci.yml` ＋ `scripts/portable-tests.sh` ＋ 本文件。⛔ `testdata/**`、⛔ `internal/**`、⛔ `cmd/**`、⛔ `frontend/**`／`design/**`、⛔ 票面／台账／`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\cerci\{logs,mut,bin}`。★突变走**仓外件 ＋ `-overlay`** ⇒ 无工作树突变窗口。

## 0. 开工四步（现跑）

| 步 | 读数 |
|---|---|
| ① 索引 | `git diff --cached --name-status` ⇒ **空**（未停手） |
| ② `ci.yml` 全文 | 改前 **1273** 行、6 jobs（`lint`／`test-core`／`test-windows`／`slo-smoke`／`slo-full`／`lint-frontend`）；`test-core` 原有 **7** steps，portable core 那步是 `bash scripts/portable-tests.sh --scope=core`（包列表在脚本里，⛔ 在 yml 里）；改后 **1329** 行、8 steps |
| ③ `runtests.sh:224-235` | `:224` 计 `skipped=$(count '^--- SKIP')`；`:234-239` 是判红分支（§3 逐字） |
| ④ `cer` 现有 skip | `go test ./internal/speech/cer -count=1 -v` ⇒ **rc 0、24 案 24 PASS、0 SKIP**（§2） |

**一句话**：把 `internal/speech/cer` 在 `test-core` 里变成**具名 step**（严格 runner），并同时**认领进 core scope ＋ core_pin**——于是「门禁随 CI 跑」既是 run page 上的一行判决，也是 census/GUARD C/D 能审计的事实。父包 `internal/speech` **不接**（§8）。

## 1. ★★★ 读数①：`ci.yml` 逐行 diff ＋ 接的 job/step

**接的是 `test-core`（`ubuntu-latest`，PR gate）**，位置＝`Portable package tests (core scope…)` 之后、`Stop compose services` 之前（本文件已写「新腿取下一个空位、不越人」）；step 数 7 → 8。

```diff
@@ -675,6 +675,62 @@ jobs:
         if: ${{ !cancelled() }}
         run: bash scripts/portable-tests.sh --scope=core
 
+      # ADDED STEP ONLY (ticket 15 AC#2, leg cerci-r1): no existing step was moved,
+      # edited, deleted, given a new `if:`, or made continue-on-error.
+      # [... 该块前 52 行（ci.yml:678-729）全是理由注释，逐字在提交里；8 小节＝WHAT THIS IS／
+      #      WHAT IT PROVES／WHAT IT DOES NOT PROVE／WHY A SEPARATE STEP／WHY THE STRICT RUNNER／
+      #      WHY `!cancelled()`／WHY UBUNTU AND WHY THE SUBPACKAGE。尺：git show 8d6e150f8 -- .github/workflows/ci.yml]
+      - name: "CER gate over the committed ASR baseline (ticket 15 AC#2; SPEC-10:60-61, SPEC-04:125)"
+        if: ${{ !cancelled() }}
+        run: bash tools/d22scan/runtests.sh ./internal/speech/cer/
+
       - name: Stop compose services
         if: always()
         run: docker compose -f docker/compose.test.yml --profile test down
```

★理由注释里两处**必须读到的边界**（逐字要点）：①「the 6%/15% thresholds … are EXERCISED, not MET」——真 ASR 假设属 `15-e`；②「The PARENT package ./internal/speech/ is NOT wired here」（依据见 §8）。
★`scripts/portable-tests.sh` 的功能性两行（同提交，GUARD C 要求 scope/pin 同步）：`core_pin` 加 `github.com/CarlosShao/wisp/internal/speech/cer`（`internal/session` 与 `internal/statemachine` 之间）、`core)` 的 scope 数组加 `./internal/speech/cer/`；另有 26 行注释写清**为什么认领**／**为什么父包不认领**。

## 2. ★★★ 读数②：`cer` 子包有没有 `--- SKIP` ⇒ **一条都没有**

```
go test ./internal/speech/cer -count=1 -v  ⇒ rc 0
=== RUN =24   --- PASS =24   --- FAIL =0   --- SKIP =0
ok  github.com/CarlosShao/wisp/internal/speech/cer  0.152s
```
★`15d-r1:104` 那条具名 SKIP（`TestBaselineCorpusIsCommittedAndTheGateCanRun`）**已不再触发**：语料在仓（205 件、`sources.tsv` 101 行）。★SKIP 分支**仍在** `baseline_test.go:140`，只是 `os.Stat(sources.tsv)` 成功 ⇒ 不 skip；再出现即语料/路径坏了（§3 有 overlay 实测）。`internal/speech/cer/*.go` 全文 `t.Skip` 命中＝**1**（就是这一处）。

## 3. ★★★ 读数③：`runtests.sh:224-235` 的实际行为 ＋ 判「接进去会不会红」

逐字：`:224` 计 `--- SKIP`；`:230-233` 非零 rc 直接传播；`:234-239` **`if [ "$skipped" -ne 0 ]` ⇒ 打「N test(s) SKIPPED and SKIP is not a pass」＋「do not relax this script」⇒ `exit 1`**；`:240-245` 零 PASS 零 FAIL（`-run` 匹配不到）也 `exit 1`。

**判**：接进去**现在不会红**（0 SKIP，§2），且**该红时会红**——两发实测（全在仓外）：

| 发 | 做法 | rc | 逐字 |
|---|---|---|---|
| **A 语料缺失** | `-overlay` 把 `baselineRoot` 指向不存在的 `testdata/asr-baseline-absent-cerci` | **1** | `runtests.sh: 1 test(s) SKIPPED and SKIP is not a pass (ticket 71 AC#3) - … top-level: PASS=23 FAIL=0 SKIP=1, === RUN=24` |
| **B 空分母** | `…/runtests.sh ./internal/speech/cer/ -run TestNoSuchThingExists` | **1** | `go test exited 0 but reported NO top-level result at all` |

★A 之后 `git status --porcelain -- internal/speech/cer` ⇒ **空**（overlay 是编译期替换，不写盘）。

## 4. ★★ 读数④：本地真跑（CI 用的就是**严格 runner** 那条）

| 命令 | rc |
|---|---|
| `go test ./internal/speech/cer -count=1` | **0** |
| `bash tools/d22scan/runtests.sh ./internal/speech/cer/`（＝新 step 逐字命令） | **0**：`runtests.sh: OK - packages=[./internal/speech/cer/] top-level: PASS=24 FAIL=0 SKIP=0, === RUN=24, '[no tests to run]'=0` |
| `bash scripts/portable-tests.sh --scope=core`（认领后整层） | **1**，★但 **cer 绿**：`ok (own line) …/internal/speech/cer`、`TestBaselineCorpusIsCommittedAndTheGateCanRun PASS (0.18s)`；整体红来自**本机既存 7 红**（与 cer／本腿无交集，这几包一字节未动）：`internal/risk` 2、`internal/panel` 4、`internal/winsec` 1 |
| `GOOS=linux CGO_ENABLED=0 go vet ./internal/speech/cer/` | **0** |
| `GOOS=linux CGO_ENABLED=0 go test -c -o …\bin\cer-linux.test ./internal/speech/cer/` | **0**（4,847,276 B）⇒ ubuntu 腿**能编译**该包测试 |

## 5. ★★ 读数⑤：YAML 语法检查

```
python -c "import yaml; d=yaml.safe_load(open('.github/workflows/ci.yml',encoding='utf-8')); …" ⇒ rc 0 "YAML OK"
jobs: ['lint','test-core','test-windows','slo-smoke','slo-full','lint-frontend']
test-core steps: 8   0/1=checkout+setup-go  2=Start compose  3=Probe mock-llm
  4=Environment fork assertion  5=Portable package tests (core scope…)
  6=CER gate over the committed ASR baseline (ticket 15 AC#2; SPEC-10:60-61, SPEC-04:125)   <== 新
  7=Stop compose services
```
★注：PyYAML 按 YAML 1.1 把裸 `on:` 解析成布尔键（`triggers: ['name', True, 'concurrency', 'jobs']`）——**解析器**行为，⛔ 文件缺陷；`True` 键下的 triggers 与改前一致。

## 6. ★★ 读数⑥：⛔ 碰 `testdata/`／`internal/` 的证明（命中 0）

```
$ git show --name-only --pretty=format: 8d6e150f8
.github/workflows/ci.yml
scripts/portable-tests.sh
$ … | grep -cE '^testdata/' ⇒ 0    $ … | grep -cE '^internal/' ⇒ 0
$ … | grep -cE '^cmd/' ⇒ 0         $ … | grep -cE '^(frontend|design)/' ⇒ 0
$ git status --porcelain -- testdata internal   ⇒ 空（提交后现跑）
```
★200 个数据件（100 wav ＋ 100 txt）与 `internal/speech/**` 8 文件**一字节未动**；§3 两次突变全在 `D:\tmp\…\cerci\mut\`（仓外）。

## 7. ★★ 读数⑦：⛔ `continue-on-error`／`if: false`（新增 0）

```
$ git diff 1f1b86a5d5 8d6e150f8 -- .github/workflows/ci.yml | grep '^+'  ⇒ 56 行（含注释）
$ … | grep -E '^\+\s*continue-on-error:' ⇒ 0    $ … | grep -E '^\+\s*if:\s*false' ⇒ 0
$ … | grep -E '^\+\s*if:'                ⇒ 1（新 step 的 `${{ !cancelled() }}`；本 job 每个非 setup step 都带）
```
★全 range（到 `f20432bd7`）同一把**锚定**尺 ⇒ 0／0。★★但**不加锚**的全 range grep 会命中 1／3 行——那 1 行是本证据件里被引用的 grep 命令本身、另 3 行是注释/标题里的散文；⛔ 任何 YAML 键。
★新增 84 行**全 ASCII**（最大码位 U+007D）⇒ ⛔ emoji。

## 8. ★★ 读数⑧：父包要不要 staging ⇒ **要，故只接 `cer` 子包**

| 判据 | 现跑读数 |
|---|---|
| 裸跑父包 | `go test ./internal/speech -count=1` ⇒ **rc 1**、`exit status 0xc0000135`（三枚 sherpa DLL 不在 loader 路径，**一条测试都没跑起来**） |
| 模型 | `third_party/spike-models/`（`silero_vad.onnx`、`asr-{encoder,decoder}.int8.onnx`…）是 **fetch 目录**：`git ls-files third_party` ⇒ **0 件** ⇒ 干净 checkout 上没有 |
| 具名 SKIP | 父包 `*_test.go` 共 **7** 处 `t.Skip*`（`asr_test.go:65,78`／`cgo_boundary_windows_test.go:125,158,226`／`engines_test.go:248`／`vad_test.go:37`），全部在 session 建不起来时触发 |
| 与判红条的关系 | `runtests.sh:234-239`：**任何 `--- SKIP` 判红** ⇒ 无 staging 载体时认领父包＝**永久红**；加 `-skip`／删 skip＝放宽判据（⛔ 边界 1） |
| 子包 | `cer`：0 cgo、0 `//go:build`、0 具名 SKIP（缺语料分支除外）；`GOOS=linux` 交叉编译 rc 0（§4） |

**结论**：父包**必须 staging**（`go test -c` ＋ 三枚 DLL 同目录 ＋ `WISP_SILERO_MODEL`／`WISP_SPIKE_MODELS`），且模型不在仓 ⇒ hosted ubuntu 跑不了 ⇒ **只接 `cer` 子包**，父包如实留在 census 上（下表），⛔ 用 skip／构建标签掩盖。

**认领生效的读数（census 前/后，同一工作树）**：

| 尺 | 前（`1f1b86a5d5`） | 后（`8d6e150f8`） |
|---|---|---|
| `…/internal/speech` | `7/0 NO-SCOPE <-UNCLAIMED-HAS-TESTS` | 同左（**如实留下**） |
| `…/internal/speech/cer` | `4/0 NO-SCOPE <-UNCLAIMED-HAS-TESTS` | **`4/0 core`** |
| `census totals: packages=36 … unclaimed-with-tests=` | **2**，rc **1** | **1**，rc **1**（红的那 1 个是父包） |

★GUARD C 通过（对不上会 `exit 1` 于 `:657-670`；本次跑到 strict runner 且 `ok (own line)` 名单列出 cer）。★`sh scripts/d22scan.sh` ⇒ rc **0** `clean - no D22 ban violations`。★`bash scripts/portable-tests-selftest.sh`（lint job 的载体自测，从 `core_pin` 现抽行数）⇒ rc **0**。

## 9. 硬边界 ＋ 停手核对

- **四条边界**：①⛔ `continue-on-error`／`if: false`／`--force`／放宽断言 ⇒ **0**（§7）；②⛔ `testdata/**` ⇒ **0 字节**（§6）；③⛔ `internal/speech/**` ⇒ **0 字节**（§6）；④⛔ `frontend/**`／`design/**` ⇒ 0。
- **停手四类**：①`cer` 子包 SKIP ⇒ 没有，未触发；②须动 `internal/**` ⇒ 不必，未触发；③须动 `testdata/**` ⇒ 不必，未触发；④须 `continue-on-error` ⇒ 不必，未触发。
- **未做（明写）**：真 ASR 假设的 6%/15% 实测（要模型文件，属 `15-e`）；`SPEC-11:128` 的 `asr-cer` **自托管** job（nightly ＋ 模型变更触发）——本腿接的是**纯 Go 计算器 ＋ 已提交语料**，hosted ubuntu 每次 PR 都跑；自托管那半仍是缺口。
- **留给编排者的唯一红**：父包 `internal/speech` 的 census 行（§8）——要一个像 `scripts/wisp-cli-tests.sh` 那样**先 stage DLL ＋ 模型**的载体才能认领，需要模型分发链（`D26`）就位。
