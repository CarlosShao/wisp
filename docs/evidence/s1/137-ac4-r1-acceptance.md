# 137 AC#4 —— 终裁表（非实现者）：那 7 处递根点换得对不对，凭据是不是我自己跑出来的

**格子**：票 `.scratch/wisp/issues/137-winsecs-12-posix-rejection-legs-go-green-in-the-symlink-shape-for-the-hosts-own-varlink-not-their-own-planted-link.md`
**AC#4 一格**（票面 `:74-76`）＋ 它 `:173` 那块下面的 `>` 口径更正（`:176-187`，09-24 15:5x 编排者）。
**判据按更正版那两条判（②a／②b），不按 `:173` 原句那句"换根之后 MUT-D 软链形 11 枚转红"判**——那条已作废为凭据。
**agent**＝`acceptor-ticket137-ac4-r1`（**非实现者**；写码那程＝`worker-ticket137-ac4`，它的自述一枚未采信，下面全是本程重走）。
**AC#1／AC#2／AC#3／AC#5 一格未重裁；票面一枚勾未翻；`-done` 未改。**

三档口径（每节自标）：〔独立复现〕＝本程自己在树上跑出来的；〔日志＋归档，抽验〕＝引别程读数但本程重测过它依赖的前提；〔仅自述，不背书〕＝只有别人说过、我没复算也不据它下判的。

---

## §0 锚点自量 ＋ 码面未动证明 ＋ 闸门 ＋ 被验版本盘上身份　〔独立复现〕

### 0.1 锚点自己现取（不读脏工作树冒充被验版本）

```
$ git rev-parse HEAD
a9c4d58f99b134b0ab3253b80469dde92bb2692a
$ git cat-file -t a9c4d58f99b134b0ab3253b80469dde92bb2692a
commit
$ git rev-parse --abbrev-ref HEAD
dev
$ date
Thu Sep 24 15:48:39 CST 2026
```

### 0.2 码面未动证明（派单要求的那一条，必须为空）

```
$ git diff b1010ff..HEAD -- internal/winsec/
（空）
$ echo rc=$?
rc=0
```

⇒ 被验的 `internal/winsec/` 码面从我开工的锚点 `a9c4d58` 回到交付面最后一枚 `b1010ff` **一字未动**，
不需要停手上报。**读数全部走 `git archive` 快照，没有一枚取自工作树**（见 0.5）。

工作树确实不干净，且**与本程无关**：`git status --porcelain` 现量 **18 行**＝`design/**` 那 16 枚 tracked 文件被
owner 那侧挪成未跟踪的 `design/old/`（显示为删除）＋ 两枚未跟踪目录 `design/doubao/`、`design/old/`。
本程**一枚未还原、未提交、未删、没过问**。`git add` 只用显式路径，每枚 commit 前跑 `git diff --cached --name-only`。

### 0.3 争用闸门（现量，逐批各量一次）

- `date` ＝ `15:48:39 CST`（本机 +08）。全程**没有任何一处**用两个时间戳相减算时长。
- `gh run list --branch dev --limit 3`（开工现量＋**每发之前**由 `batch.sh` 再打一行 `GATE-PRE`）：
  三枚全 `completed`（`35967768017` ci 8m38s、`35964449249` ci 9m3s、`35958260537` slo-fresh）⇒ **无 `in_progress`**。
- `docker ps`：只有 `union-proxy`／`clipsync*`／`postgres:15-alpine`／`redis:7-alpine` 这类**长期驻留**服务（Up 4–13 天），
  `docker ps --format '{{.Image}}' | grep -c golang` 逐发＝**0** ⇒ 取读数那些批里没有第二枚 golang 容器。
  ⚠ 这一列只回答"有没有别人在同类容器里跑"，**不主张"宿主零负载"**；那五枚驻留服务本程无法让它们停下。
  ⚠ 同机的 `acceptor-ticket125-r1` 若起容器会打在这一列上：批次 A 六发的 `GATE-PRE` 全部 `golang-containers:0`。
- 仪器可用性**没有采信派单前提**：`docker version --format '{{.Server.Version}}'` ＝ **`29.6.2`**（rc=0），
  `docker images golang:1.27` 现量 `1.31GB / 327MB` ＝ **本机已有，未 `docker pull`**，容器内 `go version go1.27.1 linux/amd64`。

### 0.4 被验版本盘上身份（自己建的归档树，非工作树）

| 树键 | 出自 | 盘上路径 | 现量身份 |
|---|---|---|---|
| `swap`＝**被验那一版** | `git archive b1010ff` | `D:\tmp\wisp137ac4r1-tree-swap` | `go.mod` **883 字节**；`placement_symlink_113_other_test.go` md5 `25660e86ad055e68304839e227713059`；`ancestor_separator_108_other_test.go` md5 `6f6b9f804be1560d92f53fb8d9350eec`；`RAWROOTS=0 SWAPPEDROOTS=9` |
| `base`＝换根之前那一版 | `git archive a02da50` | `D:\tmp\wisp137ac4r1-tree-base` | `go.mod` 883 字节；113 md5 `cb8350cb5e5917a685c10b39d57d5275`；`RAWROOTS=7 SWAPPEDROOTS=2` |

⚠ **本程把 A/B 的地基钉在 `a02da50` ＝ `9c0f546^`（换根那枚码的真正父提交）**，不是实现方用的 `182daed`；
两条理由：① `git log --oneline 182daed..9c0f546` 显示中间还有一枚 `a02da50`（docs），
② `git diff 182daed..a02da50 -- internal/winsec/` **为空** ⇒ 两版在 winsec 上逐字节相同，
本程的 A/B 与实现方的 A/B 是**同一对码**，只是本程用了更严格的那个父。

两树差**只有那 7 处**（文件系统侧独立复算，不引 git）：`diff -ru` 两树的 `internal/winsec/` ⇒ 变动行 **42 行**，
其中 `^[+-]\s+root := ` 命中的正好 **7 对（−7 raw／+7 已解析）**，其余 28 行是两枚文件头的说明块。

### 0.5 本程台件与三把身份闸（读数之前先过闸）

- 台件：`D:\tmp\wisp137ac4r1-rig\{shot.sh,batch.sh,apply_mutd.py,revert_spot.py,parse.py}`；日志 `logs/`。
- 快照一律 `git archive <sha> | tar -x -C /d/tmp/wisp137ac4r1-<名>` ＋ `-v /d/tmp/...:/src`（**不走 `docker run -v C:\` 那条静默挂空且 rc=0 的假绿坑**），
  容器内 `wc -c /src/go.mod` ＋ `[ -s ]` 不成立 ⇒ `exit 97`，**不发测**。
- 形状硬断言（`gate 98`，**本程自己的路径名**，不复用别家）：软链形 `ln -s /r1priv /r1link` ＋ `[ -L ]` ＋ `readlink` 逐字核 ＋ `TMPDIR=/r1link/w137ac4r1tmp`；
  普通形**断言 `/r1link` 根本不许存在** ＋ `/r1plain` 不是软链 ＋ `TMPDIR=/r1plain/w137ac4r1tmp`。
- 落地先证再取颜色（`gate 95`）：`go build ./...` ＋ `go vet ./internal/winsec/` 任一非 0 ⇒ **不取颜色**（第一发就这样被自己挡了一次，见 0.6）。
- 树身份闸（每发容器内现打）：`RAWROOTS`／`SWAPPEDROOTS`／`MUTD_MARKERS`／`CREDIT_BRANCHES`（AC#2 那把尺活着＝3）。
- 读数口径：`go test -count=1 -v ./internal/winsec/`，四数由 `parse.py` 从 `-v` 日志**程序化生成**（**没有一枚名是手抄**），
  `-count=N` 没有出现过 ⇒ 不会翻倍；每发另打 `grep -cE '^(panic|fatal error)'` 与**逐名 FAIL/SKIP 名册**＋**名册差集**。
- 模块缓存：`GOPROXY=off` ＋ `GOMODCACHE=/gomod`。**本程自己的卷** `wisp137ac4r1-gomod`／`wisp137ac4r1-gobuild`；
  源＝别人的 `wisp137ac4-gomod` 以 **`:ro`** 挂进 `/from` 一次性 `cp -a`（原文：`726M /from` → `CP_RC=0` → `726M /to`），**源卷没被写过一枚字节**。

### 0.6 一次仪器自拒（登记，不抹；发生在取读数之前）

第一发 `smoke-swap-plain` **`DOCKER_RUN_RC=95`**：容器里空 `GOMODCACHE` ⇒ `go build ./...` 撞
`module lookup disabled by GOPROXY=off`（`internal/secret/configrefs.go:10:2` 等 5 行），`BUILD_GATE failed`。
**那一发零颜色入账**（`gate 95` 在取颜色之前 `exit`），只是**本程台件缺料**、不是树错、也不是被验码错。
处置＝按规矩以 `:ro` 源 `cp -a` 建自己的模块缓存卷，复跑同一发拿到 `TEST_RC=0` ＋ `52/30/0/0＋22/0/0`，
**与本程之外的那一发 `r3` 同数**，然后才开批。

### 0.7 本程第一枚 commit 的账（只能在提交之后现量，故本节末尾回填）

`git log --oneline -1` 与 `git show --name-only HEAD` 的原样输出见本节下方的 0.8——
它们必须是**提交之后**的读数，这是共树里"提交账只能在提交后现量"的固有循环，本程按 `A155` 那枚形状留痕。

（占位：0.8 由下一枚 commit 回填。）
