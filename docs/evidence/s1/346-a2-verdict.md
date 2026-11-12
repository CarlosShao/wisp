# 346-a2 终裁读数件（票 346 AC#2 去牙反形自证＋AC#4 双平台成对｜只做测量＋证据，产码一字未动）

腿＝`346-a2`；起手锚＝`f3d8281c2f911433051d8929215b8c954da3eeab`（dev）；台面与全部发次＝`.scratch/wisp/probes/346/a2/00..04`（原始 raw 日志在 `D:\tmp\wisp\20261112\346-a2\logs\`，临时件只建不删）。
两枚对象仪器＝票面 AC#0 名册原件：甲 `TestPOSIX120SealWindowLeafSwapIsClosed`（`posix_seal_window_120_other_test.go`）、乙 `TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1`（`openleaf_322_other_test.go`）；起手 blob 与票面三锚同串值逐字同（`0b0e5173…`/`dfaddac8…`）。

## AC#2 判＝**弃权红有牙，直证成立**（原样红⇄去牙假绿，同条件同读数成对）

1. **pristine 弃权红（真实形状复现，红句逐字在场）**：
   - 甲（容器 `-cpu=1`＋仪器自带探针档 `WISP_120_ATTEMPTS=1`，3/3 rc=1）：`AC#1 RED instrument is dead shape=leaf-swap: refused=0, sealed-own=4 - a zero window reading here proves nothing`——racer 在 4 发跑完前一次翻面都没打（refused=0），与 CI 观测（`refused=1998, sealed-own=0`）同一句、同一维。
   - 乙（容器 `-cpu=1`，5 发＝3 红 2 绿）：`AC#1 RED instrument: across 3000 raced attempts the leaf leg ("too many levels of symbolic links") answered 0 and the racer landed 16 flips (0 errors, first ""), so the damaged count above is not evidence about a window that was never reached`——CI 同 run 的形状（`answered 0 … landed 6 flips`）逐字同族；绿发（同条件）＝leafLegRefusals=1、叶腿答上即转 PASS。
2. **去牙反形（只打仓外导出树测试件物理字节，⛔ -overlay）**：甲摘 `:225`/`:233` 两支触发条件、乙摘 `:206` 一支（literal 见 02 件；真损伤红支一字未动；两树 `go vet` rc=0）。
   - 甲变异体同条件 2/2 PASS，总量行逐字同 pristine 红发（`attempts=4 … refused=0 sealed-own=4`），红句换成读数 Logf（`AC#1 INSTRUMENT DEAD reading … refused=0, sealed-own=4`）仍在屏＝**假绿直证**。
   - 乙变异体同条件 5/5 PASS，其中 4 发 `leafLegRefusals=0`＝与 pristine 红发同一未到达形状、却转绿＝**假绿直证**（1 发真绿＋安静档对照真绿＝变异只摘弃权支）。
3. **判据自检**：真实仓库三句红句逐字仍在（`:206`/`:225`/`:233` 现读）；两枚仪器 blob 收尾与起手同串；`internal/winsec` 零脏条目。
   梯次如实记：甲的弃权安静档不出（3 发全绿）、`-cpu=1` 一味也不出（1 发绿），须加 `WISP_120_ATTEMPTS=1` 才 3/3 稳定弃权；乙一味 `-cpu=1` 即达；`-cpu=1`＋容器内 32 空转进程反而反向转绿（负载⇒窗口到达）——"更挤＝更弃权"⛔ 是本台面读数，弃权只认到达数。

## AC#4 判＝**双平台成对全绿、安静档零弃权，红名册空**

1. **全包成对**（同一名册＝全包、`-count=1 -timeout 30m`、跑前 GOTMPDIR/TMPDIR 已移除/容器档不设）：
   - Windows 原生（go1.27.1 windows/amd64）：rc=0，RUN=129 PASS=78 FAIL=0 SKIP=0，`ok … 19.406s`；两枚 `!windows` 仪器在本档 0 命中＝⛔ 编译（结构性缺席，⛔ 读成绿）。
   - Linux 容器（golang:1.27 linux/amd64，`git clone --shared` 导出树）：rc=0，RUN=76 PASS=48 FAIL=0 SKIP=0，`ok … 3.038s`。
   - 四数逐枚点名（Linux 档）：甲 `attempts=2000 refused-at-check=948 refused-at-act=10 sealed-own=916 hit-foreign=0` PASS；乙 `attempts=3000 flips=1272 flipErrs=0 leafLegRefusals=378 damaged=0` PASS；Windows 档该四数不存在（⛔ 混用、⛔ 相减）。
2. **两枚仪器各 3 发独立复跑**（Linux 容器、安静档、串行）：甲 3/3 PASS（refused-at-act=18/13/11、sealed-own=912/822/791 全 REACHED）；乙 3/3 PASS（flips=1334/1192/1268、leafLegRefusals=271/237/359 全 REACHED）⇒ 安静档散布＝3 绿 0 红 0 弃权（每枚）；摆动只在 stress 档现身（AC#2 那 8+5 发逐发在 01/02 件）。
3. **红名册作差**：两发全包顶格 FAIL＝0 枚 ⇒ 归位名册为空；本腿唯一出现过的红＝AC#2 stress 档 pristine 弃权红（票 346 票面点名的仪器形状本身，同条件变异体同读数转绿＝非产码红、非新账）。

## 禁区与纪律

`internal/winsec` 收尾 `git status --porcelain`＝0 条、blob 同串；thresholds/golden/allowlist/ci.yml/.github/PLAN/specs/go.mod/go.sum/frontend/design 零触碰；写面＝`probes/346/a2/*.md` 5 枚＋本件＋票尾一条 bullet（逐枚非 0 字节、CR=0、U+00A7/U+00B7＝0、`.out` 0 枚）；门禁 `go vet ./internal/winsec/` rc=0、`gofumpt -l`（写面＋正控）rc=0；提交只走显式手写 pathspec、⛔ push。

**本腿给下一程的一句话**：AC#2/#4 两格的读数面已齐；选形（甲台面／乙仪器／丙判据）仍归编排者与裁腿，⛔ 由本读数推着走。
