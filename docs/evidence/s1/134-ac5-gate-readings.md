# 134 AC#5 取证 — 门禁复跑（d22scan 纯净快照 / scope 不降 / YAML 解析器 / 脚本 lint）2026-09-23

被验版本：`b9b2072`（AC#2 + AC#3）+ `e951dfa`（AC#4）+ 本文件所在提交。

## 1. `sh scripts/d22scan.sh` 纯净快照 rc=0

```
$ rm -rf /tmp/d22-134b9b && mkdir -p /tmp/d22-134b9b
$ git archive HEAD | tar -x -C /tmp/d22-134b9b          # HEAD = b9b2072
$ cd /tmp/d22-134b9b && sh scripts/d22scan.sh
runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0
d22scan.sh: scan of /tmp/d22-134b9b
d22scan: examined 224 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/d22-134b9b
d22scan: clean - no D22 ban violations
d22scan-rc=0
```

正向对照（种子红）同批通过：`--- PASS: TestBuiltBinaryGoesRedEndToEnd`（6 条子用例全 PASS）。

## 2. 台账各 scope 不降（与票 99 裁决表那一行逐格比）

| scope | 票 99 基线（`99-adversarial-acceptance.md`） | 本次 `b9b2072` | 判定 |
|---|---|---|---|
| bans #1-5 `internal/` | 197 | **202** | 不降 |
| bans #1-5 `cmd/` | 20 | **22** | 不降 |
| ban #6 `frontend/` | 37 | **40** | 不降 |
| ban #7 `internal/tools/` | 17 | **18** | 不降 |
| ban #8 `design/` | 16 | **16** | 持平 |
| ban #8 `frontend/` | 37 | **40** | 不降 |
| ban #8 `internal/` | 342 | **387** | 不降 |
| ban #8 `cmd/` | 26 | **36** | 不降 |

**8 行 scope 齐全**（与票 99 那一格的计数同为 8），零 emoji 覆盖面（ban #8）含注释与 `_test.go`。
本票**没有新增任何 Go 文件**：改动面是 `scripts/slo-check.ps1`（改）、`.github/workflows/ci.yml`（改）、
`scripts/slo-freshness.sh`（新）、`.github/workflows/slo-fresh.yml`（新）。后两枚都在 ban #8 的声明面
（`design/` + `internal/` + `cmd/`）之外，但内容本身零 emoji、零 `//nolint`，注释也算在内。

## 3. YAML 解析器复核（本仓那套仪器，逐字）

```
PARSER-OK; job count = 6
jobs = lint, test-core, test-windows, slo-smoke, slo-full, lint-frontend
all six present = True (missing: [] extra: [] )
runs-on: ubuntu-latest = 3 / windows-latest = 2 / [self-hosted, wisp-slo] = 1
slo-full runs-on = ['self-hosted', 'wisp-slo'] | env keys = ['MINGW64_ROOT','WISP_ENV'] | steps = 5
slo-full has if/continue-on-error = False False
on: = pull_request: null / push.branches: [main, dev] / schedule: [{cron: 37 19 * * *}] / workflow_dispatch: null
50:  group: ci-${{ github.workflow }}-${{ github.ref }}${{ github.event_name == 'push' && format('-{0}', github.sha) || '' }}
51:  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```

新 workflow 单独解析：`jobs= ['slo-fresh']`、`runs-on= ubuntu-latest`、
`steps= ['actions/checkout@v4', 'slo-full freshness pin (ticket 134 AC#3)', 'Shell lint for the pin']`、
`permissions= {'contents': 'read', 'actions': 'read'}`、`has if= False`、`has continue-on-error= False`。

`ci.yml` 非注释改动行只有两行（其余 33 行是注释）：

```
+  schedule:
+    - cron: '37 19 * * *'
```

## 4. 脚本 lint

| 仪器 | 命令 | 读数 |
|---|---|---|
| `sh -n` | `sh -n scripts/slo-freshness.sh` | **rc=0** |
| `shellcheck` | `command -v shellcheck` | **本机没装**（rc=1）⇒ 这一格**未验证**，等 `.github/workflows/slo-fresh.yml` 的 `Shell lint for the pin` 步（缺 shellcheck 时那一步**硬红**，不静默放行） |
| PowerShell 解析器 | `[Language.Parser]::ParseFile(scripts/slo-check.ps1)` | **PARSE-OK tokens=1663**，错误 0 条；真 CI 里那枚 job 的第 5 步 `conclusion=success` 是第二道证据（`134-ac2-ac3-...md` §5） |
| `bash -n` 对未改脚本 | — | 本票没改任何既有 `.sh`（`portable-tests.sh`/`winsec-tests.sh`/`d22scan.sh` 均 `git status` 干净） |

## 5. AC#2 的 run id 那一格（复述，防止只看本文件的人漏掉）

**run `35810714576` / job `107021435150` / step 5 `SLO full gate (six states + settle + leak)` / 结论 `success`**，
以及第二枚 **run `35810884974` / job `107021957714` / 同一步 / `success`、`all_pass=True`**。
两枚都在 self-hosted `wisp-slo` 上、都含 AC#4 的改动（blob `560186fa…` 逐字节相同）。
`schedule` 那半边**触发次数 0**（cron 只对默认分支求值）⇒ 那一格的结论还没有。
