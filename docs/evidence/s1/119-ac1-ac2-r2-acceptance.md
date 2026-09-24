# 票 119 — AC#1／AC#2 第二轮回判（`acceptor-ticket119-ac1-ac2-r2`，非实现者）

被验：`.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md`
**只裁两格**：AC#1（第一轮＝通过附条件）、AC#2（第一轮＝退回）。
AC#3／AC#4／AC#5／AC#6 与**新落的 AC#7 一枚都不裁**；票面一枚勾没翻、`-done` 没改名。
第一轮表＝`docs/evidence/s1/119-adversarial-acceptance.md`（那张表 `:4-5` 明写票面六格是**实现方自勾**）。

**本节目的与读法**：第一轮判 AC#2 退回的直接理由有三处（§二(b) 措辞、§二(c) ② 只做了一半、注释里那句被证伪的全称），
本轮逐条问"返修有没有把它变成读数"，而不是"返修有没有自述做了"。

**提交方式（一次性说明，免得下面每节都解释）**：派单要求"每裁完一节 commit 一次，并把那次
`git log --oneline -1` ＋ `git show --name-only HEAD` 原样贴进该节正文"。
共树禁 `--amend`，所以第 N 节的这两行是在第 N+1 枚 commit 里进账的；
本文件每节末尾那两行的**内容**是它上一枚 commit 的读数，**归属**由每节自己的 commit 号给出。

---

## §0 锚点自量 · 闸门 · 被验版本的盘上身份

### 0.1 锚点（自己现取，不采信任何转述）

```
开工第一步  date = 2026-09-24 15:57:58 +0800 ，git rev-parse HEAD = d55c00d0b4ff71af029075853a34b063a506ba14
                                                     git cat-file -t  = commit ，分支 = dev
取快照时 HEAD 已被同机另一程推进 ⇒ 本表被验版本**重取并钉死**为
  VERIFIED_SHA = 3a497457cfc5ea4564749cbbf80b620cdc210b71   （git cat-file -t = commit，现量）
```

⇒ **下面所有容器读数、所有 `git diff`／`--numstat` 都按 `3a49745` 这一版量**；
读数期间工作树里那三枚 `internal/winsec/**`、`cmd/wisp/**`、`internal/proc/**` 之外别人的改动与我无关，
`git status --porcelain internal/winsec/ cmd/wisp/ internal/proc/` 现量**空**（0 行）。

### 0.2 争用闸门（开工时与每批取数前各查一次）

```
gh run list --branch dev --limit 3   →  三枚全部 completed（两枚 ci failure、一枚 slo-fresh success），无 in_progress
docker ps                            →  只有 union-proxy / clipsync* / minio / postgres / redis 这些长期件，
                                        没有第二枚 golang:* 容器（＝同机那枚 137 验收程此刻不在容器里取数）
docker version --format {{.Server.Version}} → 29.6.2（rc=0，镜像 golang:1.27 本机已有，未 pull）
```

### 0.3 被验版本的盘上身份（工作树 ≠ 被验版；三处必须一致）

`git archive 3a49745 | tar -x -C /d/tmp/wisp119r2-head`（469 枚 `.go`），然后同一枚文件三处比 md5：

| 文件 | `git show 3a49745:` | 快照（容器内读） | 宿主工作树 |
|---|---|---|---|
| `internal/winsec/winsec_other.go` | `b5056918be4ed13817d236fbcae0f477` | 同 | 同 |
| `internal/winsec/dataroot_symlink_119_other_test.go` | `cec75e819aef89584cde7796878ef78a` | 同 | 同 |
| `cmd/wisp/secret_dataroot_119b_test.go` | `c2f8238b6c419754e8c211aeac9d5033` | 同 | 同 |
| `cmd/wisp/secret.go` | `8e556140125905b88d2f4e0b99c34f19` | 同 | 同 |
| `internal/proc/envfork.go` | `6272c8633f126ac6a69e5838811249e4` | 同 | 同 |

⇒ 五枚全部 SAME〔独立复现〕。**这五枚就是本轮返修的全部落点**（见 §2.1 的 numstat）。

### 0.4 本轮造出来的一发仪器坑（具名登记，方向是**假绿**）

派单给的挂载写法在 Git Bash 下**逐字照抄会静默挂空**：

```
docker run --rm -v /d/tmp/wisp119r2-head:/wisp golang:1.27 sh -c 'ls -l /wisp/go.mod; ...'
  → ls: cannot access '/wisp/go.mod': No such file or directory
  → find /wisp -name '*.go' | wc -l  = 0            ← 但整条命令 rc=0
docker run --rm -v "D:/tmp/wisp119r2-head:/wisp" ... → go.mod 883 B、469 枚 .go、md5 对上
```

⇒ 处置：**所有后续读数一律用 `D:/tmp/...` 形式**，并且每发容器内先 `ls -l /wisp/go.mod` ＋ md5 三枚被验文件，
读到空就当仪器失败、不读任何 rc。**这与第一轮表 §〇 那发"`/varlink` 被 `mkdir -p` 建成真目录"是同一族**：
挂载为空时"没有用例可跑"会被读成"没有用例失败"。

### 0.5 容器与被复用的卷（只建不删）

- 容器：`golang:1.27`（go1.27.1 linux/amd64，`id` 现量 `uid=0(root)`）。
- 模块缓存：宿主 `D:/work/base/gopath/pkg/mod` 以 **`:ro`** 挂进一发一次性容器，`cp -a` 到**我自己**的卷
  `wisp119r2-gomodcache`（派单规矩：复用别人的东西只能 `:ro` 源 → `cp -a`，**不写源卷**）；
  构建缓存另建 `wisp119r2-gocache`。宿主那枚源卷全程只读，未写一个字节。

<!-- 下面两行是 §0 这枚 commit 的读数（提交后补进正文，见文首"提交方式"） -->
