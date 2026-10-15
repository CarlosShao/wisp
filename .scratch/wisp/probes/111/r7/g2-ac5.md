# 111-r7 / g2 — 票 111 AC#5（门禁四数 + 新步步序）逐格现读

**本格判（第 0 步现读）**：**半做**——"新步排在会失败的步骤之后"那半**已有 CI 侧正证**（AC#6 的
`if: ${{ !cancelled() }}` 形状＋下面 §2 的同 run 结论），而 `bash -n` / `sh scripts/d22scan.sh`
两把尺此前只到 `111-r5` 一拍（且验收件点名"没有非实现者的复跑件"）⇒ 本腿在当前 HEAD 重取全部四数，
⛔ 未改 `scripts/**`／`.github/**` 一字节。

## §1 四数（全在 `logs/g2-gates.txt` / `g2-d22scan-*.txt`）

1. **`bash -n`**：`scripts/*.sh` 7 枚 + `tools/d22scan/runtests.sh` 共 **8 枚逐枚 rc=0**
   （含本票改动面 `portable-tests.sh`／`winsec-tests.sh`／`portable-tests-selftest.sh`／`wisp-cli-tests.sh`）。
2. **`sh scripts/d22scan.sh` 纯净快照**（`git archive HEAD | tar -x -C D:/tmp/wisp111r7/snap-ac3`，仓外）：
   **rc=0**，逐字末行 `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=237,
   bans #1-5 cmd/=40, ban #6 frontend/=115, ban #7 internal/tools/=24, ban #8 design/=30,
   ban #8 frontend/=115, ban #8 internal/=558, ban #8 cmd/=131`（全文 `logs/g2-d22scan-snapshot.txt`，
   252 行，含它自报的 `gitignore rules NOT APPLIED - git cannot be consulted`＝快照无 `.git` 的响亮方向，rc 仍 0）。
3. **同尺工作树一发**：**rc=0**，同一批数（`design/=39`、`frontend/=115`…），tail 在
   `logs/g2-d22scan-worktree-tail.txt`。
4. **YAML 合法性**：`python yaml.safe_load` rc=0，各 job 声明步数（lint 13 / test-core 7 /
   test-windows 10 / slo-smoke 6 / slo-full 5 / lint-frontend 11）——`test-windows` 含新步后 10 枚。

**"各 scope 不降"对照**（逐 scope 三世代）：
| scope | r2 10-05（树） | r5 10-06（双发同值） | 本腿 10-15 快照 / 树 | 判 |
|---|---|---|---|---|
| #1-5 internal/ | 228 | 228 | **237 / 237** | 涨 |
| #1-5 cmd/ | 38 | 38 | **40 / 40** | 涨 |
| #6 frontend/ | 85 | 85 | **115 / 115** | 涨 |
| #7 internal/tools/ | 23 | 23 | **24 / 24** | 涨 |
| #8 design/ | 39(树) | 30 | **30 / 39** | 两源各自持平（快照30=30、树39=39） |
| #8 frontend/ | 85 | 85 | **115 / 115** | 涨 |
| #8 internal/ | 512 | 514 | **558 / 558** | 涨 |
| #8 cmd/ | 104 | 104 | **131 / 131** | 涨 |
⇒ **无一下降**。"树 39 / 快照 30"这对差是未跟踪件方向（r5b 记过同一对），不是覆盖缩水。

## §2 步序（AC#5 那条"新步不得被前面的红吃掉"）

- 三个新步在盘上**全部带** `if: ${{ !cancelled() }}`（HEAD `2ad0aa86`，逐字行号）：
  `:342` selftest 步名／`:394` if／`:395` run；`:623` winlive 步名／`:671` if／`:672` run；
  `:674` census 步名／`:760` if／`:761` run；核对面 `:796` if／`:797` run（windows portable）。
  步级 `if:` 现量 **36**（35 枚 `!cancelled()` + `:490` 那枚 `always()`＝test-core 的 Stop compose）；
  作业级 `if:` 0、`continue-on-error` 0。全在 `logs/g2-ci-step-order.txt`。
- **CI 侧正证（run `38676514264`，已推送 head `39b0a246`）**——同 run 里"前面红、后面仍有结论"：
  - `lint` job `116231259522`：step 8 `gofmt (gofumpt)` = **failure**、step 9 = **failure**、step 12 `staticcheck` = **failure**；
    而 step 13 `mockllm module vet` = **success**、**step 14（本票 selftest 新步）= success**。
    ⇒ 本票新步**排在会失败的步骤之后**，且**真的答了话**（这正是 AC#5 那句"本仓实测过这道门因此从未执行"的反面）。
  - `test-windows` job `116231259514`：step 7 `cmd/wisp CLI tests` = **failure**；
    step 8 winlive = success、step 9 census = success、step 10 portable = failure、step 11 junction = success
    ⇒ step 8/9 两个新步在上一枚红步之后**各自有结论**。
  - 修前形状（票面记录，非本腿重跑）：`35600043583` step4 failure → step5-8 全 skipped；
    `111-c1` 普查里 `go vet (module)`／`gofmt denominator` 曾 **12/12 skipped**。

## §3 本格判语

- 四数齐、无一下降；步序用"同 run 结论表"钉住（而不是靠 yaml 里有 `if:` 就宣布）。
  ⇒ 建议非实现者复跑 `bash -n` 与 `sh scripts/d22scan.sh`（快照）后可翻；⛔ 本腿不翻框。
- 具名边界：d22scan 本腿未在 **CI** 里跑（CI 的 lint 步只跑 `tools/d22scan` 三枚自带步）；
  本格的"纯净快照"= 仓外 `git archive`，与本票一贯口径一致。
