# 票 122 — AC#6 结案判据：步级 CI 读数（122-ac6 裁决表）

**腿**：122-ac6 · **时间**：2026-11-11 21:39–21:48 +08 · **HEAD**：`dbf5783f7cf7533d3010d739a4edba31d6dc84db`（分支 `dev`）

## 0. 一行结论

**⛔ AC#6 今天关不了：GitHub 全域 TLS 断连（与前腿同签名），CI run id / job id / step 结论取不到——如实报被网络阻塞，零编造。**
**本机侧能做的已做完：纯净快照（`git archive HEAD`）同版 staticcheck 复确认 = windows 100 / linux 54（+1 编译噪音行），与 122-r2 期望逐字一致。**

## 1. 网络状态实测（逐条贴读数）

| # | 探针 | 命令 | 读数 |
|---|---|---|---|
| N1 | gh CLI 在场 | `gh --version` | `gh version 2.96.0 (2026-07-02)` |
| N2 | CI run 列表 | `gh run list --limit 5` | `failed to get runs: Get "https://api.github.com/...": EOF`，exit 1 |
| N3 | api.github.com | `Invoke-WebRequest https://api.github.com` | `FAILED: The SSL connection could not be established` |
| N4 | github.com | `Invoke-WebRequest https://github.com` | `FAILED: The SSL connection could not be established` |
| N5 | 内层错误取证 | 同 N3 取 InnerException | `System.IO.IOException: Received an unexpected EOF or 0 bytes from the transport stream.` |
| N6 | git 远端 | `git ls-remote --heads origin dev` | `fatal: ... TLS connect error: error:0A000126:SSL routines::unexpected eof while reading` |
| N7 | 对照探针（非 GitHub） | `Invoke-WebRequest https://www.baidu.com` | `HTTP 200` |

**定性**：非本机全域断网（N7 对照 200），是 **GitHub 定向 TLS 断连**，内层 `unexpected EOF`（N5）与 `git` 侧 `error:0A000126`（N6）＝前两腿记录的同一签名。⇒ 按票面限制 #1，AC#6 走「被阻塞」路线：**⛔ 编造 run id / job id / 读数，⛔ 推代码**。

## 2. CI 侧目标读数怎么取（网络恢复后可复跑步骤清单，逐条命令）

> 目标锚：**`lint` job 的 staticcheck 步 = `.github/workflows/ci.yml:383`**（job 定义 `ci.yml:95`，`runs-on: ubuntu-latest` `ci.yml:96`）。
> 该步是 lint job **第 12 步**（共 15 步；step 1=`actions/checkout@v4` :105，step 2=`actions/setup-go@v5` :107，step 11=`go vet (tools/d22scan module)` :378，step 12=staticcheck :383，step 13=`mockllm module vet` :489）。步自报行格式在 `ci.yml:486`。

### 步骤 A：找 run（最近一次 dev 分支、含 `lint` job 的 run）

```bash
gh run list --branch dev --limit 10        # 取最顶部 run 的 databaseId（首列）与 headSha
```

### 步骤 B：锁定 `lint` job 的 job id 与 step 结论

```bash
gh run view <RUN_ID> --repo CarlosShao/wisp --json jobs \
  --jq '.jobs[] | select(.name=="lint") | {databaseId, conclusion, steps: [.steps[] | {name, conclusion, number}]}'
```

记三个值：**run id**（=步骤 A 的 databaseId）、**job id**（上面 `databaseId`）、**staticcheck 步的 step 号与 conclusion**（应＝12 / `failure`，见 §3 预期）。

### 步骤 C：取 staticcheck 步原始日志与读数

```bash
gh run view --repo CarlosShao/wisp --job <JOB_ID> --log \
  | grep -E "=== staticcheck module|module \.: exit=|module tools/|self-report:"
```

应见三段 `=== staticcheck module: <mod> ===`，每段尾 `--- module <mod>: exit=… packages=… findings=… toolchain-crash-lines=…`，及收尾自报行：
`staticcheck self-report: version=staticcheck 2026.2.2 (0.8.2) modules=3 packages=37 findings=<读数> toolchain-crash-lines=<读数> (step exit is 1)`（`ci.yml:486` 原文格式）。

### 步骤 D：与期望对照并关票

期望读数（来源＝122-r1 基线 103/57 经 122-r2 真删 3 条 U1000，commit `1bb73e6d51e619929ce167d8503328bb497c51f0`）：

| 读数 | 期望值 | 说明 |
|---|---|---|
| findings（GOOS=linux，root module） | **54** | CI native 形状（ubuntu-latest）；**+1 条编译噪音行**（sherpa-onnx-go-linux 构建约束排除，`grep -cE` 正则不计它，r2/r1 已登记同一环境形状） |
| findings（GOOS=windows 侧） | **100** | ⚠ CI lint job 无 windows 侧 staticcheck（全文件单点 `ci.yml:383`）；windows 读数只能由**同版本交叉分析**（§3 本机读数）承载，CI 步自报行只会报 linux 侧 |
| d22scan / mockllm 两 module | 各 **0** | r1 已实测 0 |
| toolchain-crash-lines | **0** | 出现非 0 ⇒ export-data 版本错位，读数作废 |
| 步结论 | **failure**（exit 1） | 预期红，门禁绿不属票 122 判据（`ci.yml:430-432` 注） |

四件齐（run id / job id / step 号＋结论 / linux findings=54+噪音 0 crash）⇒ AC#6 可关。

## 3. 本机复确认（纯净快照，同版同形状，可复跑）

- **快照**：`git archive HEAD | tar -x -C D:\tmp\wisp\20261111\122ac6\snap`（`dbf5783f7`）。
  - 1722 = 1722（repo vs 快照 `.go` 文件数，等）；tar rc=1 的 16 枚失败件全在 `.scratch/wisp/probes/**`（Windows tar 文件名怪癖），代码面零损失。
  - 删码已在快照内现核：`shell_exec.go` 含 `waitOutcome` 定义行 :65 / 构造 :128，**无 `out`/`errs` 字段**——`1bb73e6d5` 真删生效，且快照不含任何工作树残渣（`build\mut_rvw.go` 不入 `git archive`）。
- **版本**：`staticcheck.exe 2026.2.2 (0.8.2)`（`D:\work\base\gopath\bin\`，与 CI 钉的 `@2026.2.2` `ci.yml:465` 同版）。
- **形状**：与 CI 步同款三 module loop（`.` / `tools/d22scan` / `tools/mockllm`），findings 计数用 CI 同一正则 `\((SA|ST|S[0-9]|QF|U)[0-9]+\)$`（`ci.yml:477`），crash 行用同一串 `export data version`（`ci.yml:478`）。

| 读数 | GOOS=windows | GOOS=linux（CI native 形状） |
|---|---|---|
| root module `.` | **100 条**，exit 1，0 噪音行 | **54 条**＋**1 条编译噪音行**（`sherpa-onnx-go-linux@v1.13.8` 构建约束排除），exit 1 |
| tools/d22scan | 0 条，exit 0 | 0 条，exit 0 |
| tools/mockllm | 0 条，exit 0 | 0 条，exit 0 |
| 合计 findings | **100** | **54** |
| toolchain-crash-lines | **0** | **0** |

- **对照 122-r2 期望（windows 100 / linux 54＋噪音）**：**两 GOOS 全部成立，逐字一致**。
- 原始输出存档：`D:\tmp\wisp\20261111\122ac6\ac6-windows-_.txt`（100 行全 finding）、`ac6-linux-_.txt`（55 行 = 54 finding＋1 噪音）、`ac6-{windows,linux}-tools_{d22scan,mockllm}.txt`（0 字节＝空集）、`version-windows.txt`。

## 4. 越界检查

- tracked 写面全集＝**1 枚**：本表 `docs/evidence/s1/122-ac6-ci-reading.md`（`git status` 中 `docs/` 下仅此一枚本腿新件）。
- 未动：任何 `.go` / `.yml` / `scripts/**` / thresholds / golden / `PLAN.md` / `docs/specs`——零字节；⛔ `add -A` / amend / reset / rebase / stash / checkout . / clean；⛔ 删文件；⛔ push。快照与读数只落在票面授权的 `D:\tmp\wisp\20261111\122ac6\`。
- 本腿零代码改动：diff 为空，读数是复确认不是改码。

## 5. 回执

`DONE: 122-ac6 网络不通 被阻塞（GitHub TLS 全域断连，CI run/job/step 读数取不到，零编造） 本机复确认 100/54 成立（纯净快照 dbf5783f7 + 同版 2026.2.2 + CI 同款 loop：windows 100/0噪音，linux 54+1噪音，d22scan/mockllm 各0，crash 0）`
