# 346-a2 台面与发次总表（每一发 = 一次 docker run 或一次本机 go test）

- 腿＝`346-a2`（票 346 AC#2 去牙反形自证＋AC#4 双平台成对与三发复跑）。
- 起手锚＝`f3d8281c2f911433051d8929215b8c954da3eeab`（HEAD，分支 dev）。
- 台面：
  - 本机 Windows：`go1.27.1 windows/amd64`，工作目录＝共享工作树；跑前移除 `GOTMPDIR`/`TMPDIR`（两者本已为空，现跑确认＝空串）。
  - Linux 容器：`golang:1.27`（`linux/amd64`，宿主 12 核），`docker run --rm -v <tree>:/src -v D:\work\base\gopath\pkg\mod:/gm:ro -v <ROOT>\gocache:/root/.cache/go-build -w /src -e GOPROXY=off -e GOFLAGS=-mod=mod -e CGO_ENABLED=0 -e GOMODCACHE=/gm`。
  - ★容器档 `GOTMPDIR`/`TMPDIR` **不设**（继承 342-v1 定式：⛔ bind 到宿主 Windows 卷），`t.TempDir()` 落容器自有 `/tmp`。
  - 树：`export`＝`git clone --shared` 出的仓外导出树（`D:\tmp\wisp\20261112\346-a2\export`，HEAD 同锚）；
    `mut-jia`/`mut-yi`＝export 的拷贝＋对测试件的物理字节变异（见 02 件）；真实仓库 `internal/winsec` 全程零字节写入（终量 `git status --porcelain -- internal/winsec`＝0 条）。
- 原始 raw 日志全部在 `D:\tmp\wisp\20261112\346-a2\logs\`（临时件只建不删）；本目录只放读数与判语。

## 发次表（逐发 rc 与四数）

| 发 | 树 | 条件 | `-run` | rc | 顶格四数 |
|---|---|---|---|---|---|
| win-native-fullpkg-1 | 共享工作树 | Windows 本机，GOTMPDIR/TMPDIR 已移除 | 全包（无 -run） | 0 | RUN=129 PASS=78 FAIL=0 SKIP=0（`ok … 19.406s`） |
| a-jia-quiet-1..3 | export | 容器安静档 | `^TestPOSIX120SealWindowLeafSwapIsClosed$` | 0,0,0 | 每 RUN=1 PASS=1 FAIL=0 SKIP=0 |
| a-yi-quiet-1..3 | export | 容器安静档 | `^TestAC1POSIXOpen…322r1$` | 0,0,0 | 同上 |
| a-fullpkg-posix-1 | export | 容器安静档 | 全包（无 -run） | 0 | RUN=76 PASS=48 FAIL=0 SKIP=0（`ok … 3.038s`） |
| b-jia-cpu1-full | export | 容器 `-cpu=1`（梯 1，不达） | 甲名 | 0 | PASS |
| b-jia-cpu1-att1 + -2/-3 | export | 容器 `-cpu=1`＋`WISP_120_ATTEMPTS=1`（梯 2） | 甲名 | 1,1,1 | FAIL（弃权红，见 01 件） |
| b-yi-cpu1 + -2..-5 | export | 容器 `-cpu=1`（梯 1 即达） | 乙名 | 1,0,1,1,0 | 3 红 2 绿（散布见 01 件） |
| b-yi-cpu1-load32 | export | 容器 `-cpu=1`＋容器内 32 个空转进程（梯 2，方向反向） | 乙名 | 0 | PASS（leafLegRefusals=9＝窗口到达） |
| c-mutyi-cpu1-1..5 | mut-yi | 容器 `-cpu=1`（与 pristine 同条件） | 乙名 | 0×5 | PASS（其中 4 发 leafLegRefusals=0＝假绿直证，见 02 件） |
| c-mutjia-cpu1-att1-1..2 | mut-jia | 容器 `-cpu=1`＋`WISP_120_ATTEMPTS=1`（同条件） | 甲名 | 0×2 | PASS（读数逐字同 pristine 红发，见 02 件） |
| c-mutyi-quiet-1 / c-mutjia-quiet-1 | mut-yi / mut-jia | 容器安静档（对照） | 乙名 / 甲名 | 0,0 | PASS 且窗口到达（llr=257 / sealed-own=853） |
| mut-jia / mut-yi vet | mut-* | 容器 `go vet ./internal/winsec/` | — | 0,0 | 编译面干净 |

## 梯次注记（诚实记录，⛔ 挑好看的）

- 甲的弃权在安静档不出现（3 发全绿、窗口到达）；`-cpu=1` 一味也不够（1 发绿：refused-at-check=1600、refused-at-act=0、sealed-own=398）；
  **`-cpu=1`＋该仪器自带的探针档 `WISP_120_ATTEMPTS=1`**（文件头自陈 "Scaled by WISP_120_ATTEMPTS for a cheaper probe run"，
  日志行自带 `attempts=1` 自名）3/3 稳定弃权红： racer 在 4 发跑完前根本没被调度 ⇒ refused=0。
- 乙的弃权一味 `-cpu=1` 即达（与 CI 形状同族：CI 观测 flips=6/liveness=0，本台面 5 发里 flips=16..140、leafLegRefusals 0 或 1）＝ racer 挨饿、
  叶腿一次答不上；加负载（load32）反而反向（syscall 变慢 ⇒ P 交出 ⇒ 窗口到达＝绿）＝如实记为"反向梯"。
- 安静档 3 发复跑（AC#4）与 stress 档（AC#2 条件）逐发分开落日志，⛔ 混读。
