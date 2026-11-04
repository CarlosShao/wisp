# 票 234 读数表 — gate-rerun-3（编排者亲跑，当前 HEAD caea1b5df，2026-11-04 13:2x-13:4x +08）

> 前置：gate-rerun-1 在旧 HEAD 015f6be1 全做过；gate-rerun-2/3 局部做完 AC#1/M1/M2/M5/AC#3 部分/AC#6。**403 修好 `internal/risk` 测试 import cycle 后**，编排者亲跑补齐全部 12 格与 6 AC 判语。证据目录 `.scratch/wisp/probes/gate-rerun-3/`（及 `gate-rerun-2/` 的 M1/M2/M5 原始 -v 日志）。本票 6 AC 均为「取到当前 HEAD 读数并登记」，非要求全绿——红的如实归因。

## AC#1 起手基线名册（当前 HEAD，绿）
- 包：winsec/risk/config/approval；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" go test -p 1 ./internal/winsec/ ./internal/risk/ ./internal/config/ ./internal/agent/approval/ -count=1 -v`（`gate-rerun-2/ac4-final-v.txt`）
- 四数(顶层)：**rc=0 RUN=726 PASS=455 FAIL=0 SKIP=2**；名册 457 顶层用例名。
- ⚠ 旧 baseline（`gate-rerun-2/ac1-baseline-v.txt` 285 名）是**在 risk 测试 import cycle 坏（[setup failed]）时**取的、缺 risk 全 172 枚 ⇒ **已作废**；本表用 403 修后的干净名册为新基线。SKIP 2=既有 `TestDefaultDeadlineWallClockMeasurement` 一族。

## AC#2 五枚变异（M1..M5）
- **M-1 票92:68**（junction 拒绝改放行）：gate-rerun-2 现 HEAD：overlay build rc=0，winsec 7 red（`TestAC3JunctionInputIsRefusedNotSealed` 等），restore green，worktree 零改。（`gate-rerun-2/m1-*-red-names.txt`）
- **M-2 票104:52**（三方向）：gate-rerun-2 现 HEAD：M2a 4 red / M2b 10 red / M2c 4 red，restore green。（`m2a/b/c-red-names.txt`）
- **M-3 票110:39**（新步自己会红）本表亲跑 /tmp 快照（`m3-summary.txt`）：正控 `bash scripts/winsec-tests.sh` rc=0、四数 RUN=127/PASS=76/FAIL=0；变异A `winsec_windows.go:174 set[ace.trustee]→set[ace.text]` build=0→step **rc=1**、`0xc0000135`=0、≥10 具名红（TestAC118*/TestAC1* 族）；变异B `./internal/config/` **rc=2 GUARD1** 显式失败；worktree 零改。
- **M-4 票113:54**（两形，真 linux 容器 `m4-container.txt`/`m4-inner.sh`）：Docker golang:1.27 linux/amd64（go1.27.1，848 files，假绿护栏过：go.mod yes）；pristine `ok winsec 3.219s`；M4a `ancestorIsLink→false`(:419) build=0→test **rc=1 11 具名 POSIX red**（TestAC1POSIX*/AC118POSIX*/AC2/AC3）；M4b `pathPieces(path)[:1]` 半修 build=0→**同 11 red**；两次 restore identical；仓库零改。
- **M-5 票115:58**（两形）：gate-rerun-2 现 HEAD：M5a/M5b 各 2 red（`TestNoticeAttributionSurvivesAn83AliasOfItsOwnTree` 等），restore green。（`m5a/b-red-names.txt`）

## AC#3 门禁 7 枚（G-1..G-7）
- **公共**：全仓 `gofmt -l`/`gofumpt.exe -l` 点名 `internal/agent/approval/pending_read.go` + `internal/agent/tools.go`；经 `gofmt-diff-1/2.txt` 的 `-`/`+` 同文核对＝**CRLF/autocrlf 假红（第17处）**，非真格式问题。（`gate-gofmt.txt`/`gate-gofumpt.txt`）
- **G-1 票92:73**：panel vet rc=0；`sh scripts/d22scan.sh` rc=0 clean；panel `-count=2` **rc=1 FAIL=2**（1 名×2，frontend-coupled 既有红，非本票）。（`ac3-panel-count2-v.txt`/`gate-d22scan.txt`）
- **G-2 票97:55**：approval vet rc=0；`-count=2` **rc=0 RUN=222 PASS=172 FAIL=0 SKIP=2**（SKIP=`TestDefaultDeadlineWallClockMeasurement`×2）。（`ac3-agent_approval-count2-v.txt`）
- **G-3 票104:57**：winsec RUN=254/152/0/0、memory RUN=336/196/0/4、secret RUN=58/42/0/0、approval RUN=222/172/0/2，**全 rc=0**；逐包 vet rc=0；`GOOS=linux go vet` 见 `gate-vet-linux.txt`（winsec/memory/secret ok，tools/risk 红＝跨平台编译缺口非本票）。（`ac3-*-count2-v.txt`）
- **G-4 票105:54**：tools `-count=2` **rc=0 RUN=1068 PASS=770 FAIL=0 SKIP=0**；risk **rc=0 RUN=562 PASS=342 FAIL=0 SKIP=2**（403 修后可跑）。vet：risk 现 rc=0（cycle 已修）；tools rc=1（`system_backend_windows.go:337 unsafe.Pointer`，真红/非本票）。
- **G-5 票110:47**：`bash -n` 三脚本 rc=0（`scripts/winsec-tests.sh`/`portable-tests.sh`/`tools/d22scan/runtests.sh`，`gate-bashn.txt`）；d22scan rc=0 clean；排序 clause＝ci.yml `!cancelled()` 守卫（gate-rerun-1 已逐字读，本 HEAD 同形态）。
- **G-6 票113:56**：容器三包 count2（`g6-container.txt`）：**g6_rc=1 RUN=986 PASS=588 FAIL=4 SKIP=4**；winsec `ok 10.042s`、risk `ok 27.823s`；memory FAIL=2 名（`TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys`、`TestMigration31V5ToV6KeepsOldData`）。⚠ memory-linux 红**非产品范围**（产品 Windows-only；memory 在 Windows host `ac3-*-memory` 全绿 rc=0）；本格分母＝**winsec POSIX 已 ok**。
- **G-7 票115:64**：winsec `-count=2` 见 G-3（RUN=254/152/0/0，SKIP=0 真零）；gofmt/gofumpt 见公共；vet winsec rc=0；d22scan rc=0。

## AC#4 终态名册 vs 基线
- 五变异全走 overlay 或 /tmp 快照，**从不写 worktree**（每发 `git status --porcelain -- internal/ cmd/wisp` 仅余一枚 `.txt` 临时件）⇒ 基线名册与终态名册同 HEAD、逐名一致，difference=0。

## AC#5 10 枚重复分母搬迁
- **已满足**（gate-rerun-1 旧 HEAD 已处置并提交）。现盘逐数核：`⛔ 指向` 标记 **10 枚**＝07:3→64 / 11:1→12 / 63:1→12 / 80:2→230 / 97:1→230 / 115:1→230 / 92:1→114，原句逐字保留、续行未动；删除列非 0（换勾框＝1 删，同 gate-rerun-1 口径说明 A446）。后续票又销了 97/92 更多格（un 更低＝不回归）。

## AC#6 坏尺作废
- `ac6-yardstick.txt`：`grep -rn 'NewGate('` 生产 0 命中、全仓含测试 0 命中 ⇒ 该串无处存在、推不出事＝void；正尺 `risk.Gate(` 定义 `internal/risk/blacklist.go:77`、生产调用 `internal/tools/mode.go:98`（1 枚）、测试调用 2 文件。判性不变（丙）。

## 未证 / 诚实标注（不追，非本票分母）
- G-1 panel rc=1（frontend-coupled 既有红）、G-6 memory-linux 红（Windows-only 产品范围外）、tools vet unsafe.Pointer（真红非本票）、`GOOS=linux go vet` tools/risk 跨平台编译缺口：均如实登记、归口别票或产品范围外。
- M-3 日志因 W 变量内含空格误置重定向，`m3-summary.txt` 为同一真实运行的读数转录（四数与 red 名同源）。
- **全部 12 格均有当前 HEAD 读数，无 UNJUDGED。**
