# `00-anchor.md` — 票 310 增量程 `310-r2`（实现腿）

派单正文＝`.scratch/wisp/probes/orch/20261012/txt/310r2-dispatch.txt`（本腿逐字读毕并照办）。

## 起手锚（本腿自己现读，⛔ 照派单抄）
- 起手 HEAD（本腿开跑时）＝`73719706`（`dev`），full `737197060662adb6d5235f787ca9cdf8a83ecffd`；件 `logs/00-anchor-read.txt`。
- 派单写的锚＝HEAD `27176c2c`；`git merge-base --is-ancestor 27176c2c HEAD`＝**yes**（其间 `8e3cb370`＝派 310-r2、`af5179dd`／`5552567c`／`134f1d30`／`73719706`＝票 311 三笔落地，均⛔ 碰本腿射程）。
- ★本腿四枚目标文件在 `8e3cb370`（本腿建导出树时的 HEAD）⇄ 现 HEAD `73719706` 之间 **blob 逐枚 SAME**（`winsec 74c3ccf7 / cli 85efcd4d / runtests 932aa9e3 / portable 7bc4e287`），且 `git diff --name-only 8e3cb370 HEAD -- scripts tools`＝**空** ⇒ 本腿"改前那一发"＝货真价实的 HEAD 原文，未被他腿改写。件 `logs/`（见 `30-gates.md` §0）。

## 本腿只做编排者已裁的两笔欠账（⛔ 扩大射程）
- **`Q-310-orch-P`＝名册外第五枚起 `go test` 的入口今天照旧漏**：`scripts/d22scan.sh:51` → `tools/d22scan/runtests.sh:75` 那行 `go test -v -count=1 "$@"` 是真起点。⇒ 在 `runtests.sh` 里补上与 `scripts/winsec-tests.sh:55-65` **同一形**的专用根步子（逐字照形，只换 `<entry>` 名 `d22scan-` 与 echo 前缀），守卫仍带 `-z "${GOTMPDIR:-}"`（委托那发继承上层那枚根，一次运行只一枚根）。
- **`Q-310-orch-Q`＝`LOCALAPPDATA` 为空那一支整块静默**：现形是合取 `&& [ -n "${LOCALAPPDATA:-}" ]`，任一不成立 ⇒ 整块跳过、⛔ `GOTMPDIR=` ⛔ `DEGRADED`。⇒ 只在「`GOOS`＝windows 且 `LOCALAPPDATA` 为空」这一形打**一行**具名告警（前缀同形 `<文件>: ticket 310 AC#1: ...`）；⛔ fail-closed（⛔ 退非 0、⛔ 拒绝跑测试；撤销口令「310 改 fail-closed」），⛔ 动非 windows 那两支。四枚文件同一形。

## 判据的牙（两形各自⛔ 恒真、被现量咬到）
- ①（P）改前 `sh scripts/d22scan.sh` 任何读数⛔ 出现 `GOTMPDIR=`；改后同一入口出现**恰一行** `GOTMPDIR=`，且那一行的根**当场在盘上存在**＋名册尺证它⛔ 落 `%TEMP%` 顶层。→ 读数在 `20-paired-readings.md`。
- ②（Q）改前把 `LOCALAPPDATA` 抽掉那一发**什么行都不打**（完整 stdout/stderr＋`rc=` 落盘，⛔ 用"没输出"当证据）；改后同一台面同一抽法打出那行新告警，`GOTMPDIR` 仍未设、`rc` 与改前同数。→ 读数在 `20-paired-readings.md`。

## 写面（三件套，逐字照办；⛔ 越界、⛔ 顺手清别家）
1. `.scratch/wisp/probes/310/r2/**`（本腿件）。
2. 票 310 末尾**追加一条 bullet**（`git diff --numstat`＝**1 add / 0 del**，⛔ 改任何既有行）。
3. §1 点名的 `scripts/portable-tests.sh`／`scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`／`tools/d22scan/runtests.sh`（本腿实改三枚＝winsec/cli/runtests；portable 是对照、blob 未动）。

## ⛔ 清单（本腿全程）
⛔ 翻任何 `AC` 框 / ⛔ 碰 `Status:` 行 / ⛔ 把票改名 `-done`（归编排者）/ ⛔ 动 `.github/**` 与 `ci.yml` / ⛔ 动任何 `*_test.go` 断言 / ⛔ 动 `internal/**` 生产分级 / ⛔ 动票 304 的 Job Object 面（`internal/proc/**`）/ ⛔ 动票 311 归因 / ⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43／C1–C32／D1–D47；只 commit ⛔ push；`git add` 与 `git commit` 均带**显式 pathspec**；⛔ `add -A`/`.`/整目录扫；⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`；临时件只建不删；⛔ 删 `%TEMP%` 顶层任何一枚（现量 59 枚）；共享工作树里别家未提交（`.gitignore`/`go.mod`/`design/**`/别人 `probes/**`）⛔ 提交、⛔ revert、⛔ 修改、⛔ 进本腿任何一笔。
