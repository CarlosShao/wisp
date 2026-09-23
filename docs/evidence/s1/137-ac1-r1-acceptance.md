# 137 AC#1 —— 验收方第二程（r1）：把测量方那句"部分成立"独立重走一遍

裁决方：`acceptor-ticket137-ac1-r1`（**非实现者、非测量方**）。本格只裁 **票 137 AC#1 这一格**，
不修任何东西、不翻任何勾。测量方那程是 `worker-ticket137-ac1`，它的证据是
`docs/evidence/s1/137-ac1-teeth-or-not.md`；本件是**独立复现**，不是它的附录。

写作与取证交错进行：先测、再写、最后提交，每裁一节 commit 一次。

## §0 锚点与树（含 winsec 侧差集）

- 开工首读 `git rev-parse --short HEAD` = **`4a0d7a4`**，全 sha `4a0d7a48d03dc6aee1a3e18ceee1d2a99269b9e6`，
  `git cat-file -t` = `commit`。**本程全部读数都在这枚锚上量**。开工时刻 `date -u` = 2026-09-23 15:21:16 UTC ⇒ 本机 23:21 +8。
- 锚点来源纪律：**没有任何一枚 sha 是从别人的报告或工具输出里直接取来当锚用的**。
  测量方自报的 `1d38206` 我只在"核它存不存在、是不是我这一枚的祖先"这两条只读命令里用到：
  `git cat-file -t 1d38206` = `commit`（真实存在），`git merge-base --is-ancestor 1d38206 4a0d7a4` = 真（它的锚在我这枚的祖先线上）。
- **被测面差集（它核过"winsec 侧零 hunk"，这条我重走了、没沿用）**：
  `git log --oneline 1d38206..4a0d7a4 -- internal/winsec/` **输出为空** ⇒ 从它的锚到我的锚，
  `internal/winsec/` 一字未动，两程读的是同一版被测码。
  （同一个区间里飘进来的提交全在 `docs/**`、`.scratch/**`、`internal/observe/**`、`cmd/wisp/**`，与本格无关。）
- 树：**全部**由 `git archive 4a0d7a48d03dc6aee1a3e18ceee1d2a99269b9e6 | tar -x -C /d/tmp/wisp137r2-<名>` 落成，
  仓内**未建 worktree、未 checkout、未 reset**。目录一律只建不删：
  `wisp137r2-tree0`（纯净）／`-tree-muta`／`-tree-mutb`／`-tree-mutab`／`-tree-mutd`，台件与十＋两发日志在 `wisp137r2-io/`。
  ⚠ 我没有把脏工作树当被验版本：工作树里此刻有别人的未提交件（`docs/reports/**` 两枚已改、一枚未跟踪），
  `git status --porcelain internal/winsec/` 输出为空 ⇒ winsec 侧工作树与锚上同版，但我仍然只用归档树取数。
- 生产码零改动：本程对 `internal/winsec/**` **只读**（AC#1 不修东西，收紧断言是 AC#2 的地界）。

## §1 容器与挂载证明

- `docker info` ⇒ ServerVersion **29.6.2**、`linux/x86_64`。宿主是 Windows，而这一族用例带 `//go:build !windows` ⇒
  **宿主没有分母**，本程一片读数都不在宿主取（宿主只跑 `git`／`grep` 这类取证命令）。
- 镜像 `golang:1.27`，容器内 `go version` = **`go1.27.1 linux/amd64`**。
- 挂载写成 `/d/...` 且带 `MSYS_NO_PATHCONV=1`；**每一发**开跑前先 `ls -l /src/go.mod`：

  ```
  -rwxrwxrwx 1 root root 883 Sep 23 15:19 /src/go.mod
  go.mod sha256: d13ba3de2d319f298ed1e598f1702015ca50e646ae293daf183545ce794f40fc
  go version go1.27.1 linux/amd64
  ```

  十二发（十发正式＋两发复跑）里这枚 sha256 与字节数**逐发相同**；脚本里另有硬闸
  `[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97` ⇒ 静默挂空取不到读数。
  ⚠ 顺带把一枚容易读成"分歧"的数说清：仓内工作树与 `git show <锚>:go.mod` 都是 **855** 字节，
  而 `git archive | tar -x` 出来的快照里是 **883** 字节（差 28＝`.gitattributes` 第 1 行 `* text=auto` 配 `core.autocrlf=true` 带来的行尾差，
  28 行 go.mod 每行多一枚 CR）。
  我引用挂载证明时用的是**容器里那个文件的实际字节数**（883），与测量方报的 883 同形，不是同一个量被抄了两遍。
- 被测文件指纹（tree0，容器内 `md5sum` 与宿主侧 `md5sum` 同值）：
  `internal/winsec/winsec_other.go` = `b5056918be4ed13817d236fbcae0f477`、
  `internal/winsec/winsec.go` = `a6144c880de80e43bb1393f3624e7221`
  ⇒ **与测量方 §1 报的两枚 md5 逐字相同**，这是"两程读的是同一版码"的第二条独立凭据（第一条是 §0 的差集为空）。
- 两形硬断言（我的台件命名与测量方**不同**，免得把它的形状当成我的）：
  - 软链形：`mkdir -p /r2priv/w137r2tmp` ＋ `ln -s /r2priv /r2link` ＋ `[ -L /r2link ]` ＋ `readlink` 逐字核为 `/r2priv`
    ＋ `ls -ld` 打印 `lrwxrwxrwx 1 root root 7 ... /r2link -> /r2priv`，`TMPDIR=/r2link/w137r2tmp`；
  - 普通形：断言 `/r2link` **根本不许存在** ＋ `[ ! -L /r2plain ]`，`TMPDIR=/r2plain/w137r2tmp`。
  - 十二发**无一命中 97／98／99**。新鲜容器一枚一发（`--rm`，每发从零起），形状目录不许预存在。
- 模块缓存：`GOPROXY=off` ＋ 我自己的两枚 named volume `wisp137r2-gomod`／`wisp137r2-gobuild`
  （离线复用既有缓存只是省时间，**不是复用别人的读数**；缺件会直接 fail，不存在"网络慢装出来的绿"）。
- 跑法：`go test -count=1 -v ./internal/winsec/`；读数前先 `go build ./...` 与 `go vet ./internal/winsec/`，
  两者任一非 0 ⇒ `exit 95`＝**落地不过证就不取颜色**。十二发全部 `BUILD_RC=0` ＋ `VET_RC=0`。
- 台件（全部在仓外 `D:\tmp\wisp137r2-io\`）：`mutate.py`（变异，逐处锚串唯一性断言）、`make-trees.sh`（建五棵快照）、
  `run.sh`（容器内跑发＋硬闸）、`parse.py`（把 `-v` 日志程序化拆成逐名颜色与名册，非手抄）、`matrix.py`（十发矩阵）。
