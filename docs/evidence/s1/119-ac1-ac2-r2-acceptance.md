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

```
$ git log --oneline -1
6e04d1a docs(119 复判r2,§0): 锚点重取并钉死 3a49745、争用闸门、五枚被验文件三处 md5 一致；登记派单给的 /d/tmp 挂载写法实测静默挂空 rc=0 这发假绿坑

$ git show --name-only HEAD
docs/evidence/s1/119-ac1-ac2-r2-acceptance.md
```

⇒ §0 落在 `6e04d1a`，那枚 commit **只带我这一枚文件**（`git show --name-only` 原样如上）。

---

## §1 AC#1 —— 逐条对第一轮那三枚"条件"：**成立（无附条件）**〔独立复现〕

第一轮表（`119-adversarial-acceptance.md` §一）给 AC#1 判"通过附条件"，三枚条件分别是：
**条件 1**＝第二条生产路（`cmd/wisp/secret.go` 的 `resolveSecretLayout`）用例全盲（它写"这一条是硬的、直接喂给 AC#2 的退回"）；
**条件 2**＝5 枚里 2 枚的 fixture 拿被测函数算期望 ⇒ `MUT-119-R` 读出 56/56 全绿（它写"这条最硬"）；
**条件 3**＝"自证挂载非空"是跑者的性质、实现方那份挂载自证只在票面转述里。

本轮只问一句：**返修有没有把那一枚条件变成读数**。下面每条都给我自己这一轮跑出来的发，不抽验它贴的日志。

### 1.1 条件 1 —— **已变成读数**

**静态那一半（第一轮的两句 `grep` 命中 0，在今天不成立）**：

```
git show 3a49745:cmd/wisp/secret_dataroot_119b_test.go | grep -c "os.Symlink"   = 3     （第一轮：cmd/wisp/*_test.go 全仓 0）
git grep -l resolveSecretLayout 3a49745 -- cmd/wisp    = secret.go + secret_dataroot_119b_test.go + dataroot_128_test.go
四枚生产读数点全部在“解析 OS 答案”那一侧：internal/proc/envfork.go:125、:245、cmd/wisp/doctor.go:263、cmd/wisp/secret.go:148
```

**读数那一半（我自己重跑 MUT-119-SEC 这一发，容器 `golang:1.27`、go1.27.1 linux/amd64、uid=0）**：
命令逐字＝`go test -count=1 -v -run 'TestAC1POSIXSecretRoute|TestAC3POSIXSecretRoute' ./cmd/wisp/`，两形各一遍；
变异＝把 `cmd/wisp/secret.go:148` 那一行**注释掉**（落地原文 `148:	// root = proc.SealableRoot(root) /* MUT-119-SEC-R2 single-line removal */`，
先证"它不再是活语句"（`grep -n "^[[:space:]]*root = proc.SealableRoot(root)"` → 无命中），再 `go vet ./cmd/wisp/` **rc=0**，才读数）。

| 发 | 形状 | rc | RUN | PASS | FAIL | SKIP | `^panic\|^fatal error` | 名册差集 |
|---|---|---|---|---|---|---|---|---|
| A | plain（`TMPDIR=/tmp/plain119r2`） | 0 | 3 | 3 | 0 | 0 | 0 | — |
| B | 软链（`TMPDIR=/varlink/w119tmp`，`ls -ld` 读到 `lrwxrwxrwx … /varlink -> /realpriv`） | 0 | 3 | 3 | 0 | 0 | 0 | A↔B 名集 `diff` rc=0 |
| **E** | plain ＋ **拆那一行** | **1** | 3 | **0** | **3** | 0 | 0 | A↔E 名集逐字相同、**三枚颜色全反** |
| **F** | 软链 ＋ **拆那一行** | **1** | 3 | **0** | **3** | 0 | 0 | E↔F 名集 `diff` rc=0 |

**红在哪一行（派单点名要核的，不是抽验）**：

```
RUN E（plain）:
  secret_dataroot_119b_test.go:171  ← 报在 helper 调用点（assertSecretRouteSeals119b 里 t.Helper() 所致）
      message 出自 helper 内 :122 "AC#1 RED: the secret route's data root is …/homelink119/.config/wisp-dev,
                                    want the tree the OS reads (…/homereal119/.config/wisp-dev)"
      message 出自 helper 内 :127 "AC#1 RED: secret.NewStore refused the dev data root …:
                                    winsec: refusing to seal … reaches it through the link at
                                    …/001/homelink119, which is not the tree this call names"
  secret_dataroot_119b_test.go:192  ← 同一枚 helper，另一条形（XDG_CONFIG_HOME），拒因点到 …/001/xdglink119
  secret_dataroot_119b_test.go:215  ← "AC#3: the route refused to seal its own tree"（第三枚走的是另一条腿）
```

⇒ **红因是产品原文、不是前提破了**：plain 形里被点名的链接就是用例自己种的那枚（`…/001/homelink119`、`…/001/xdglink119`），
而前提自证那两行（`:169`、`:190` 的 `t.Fatalf("premise broke…")`）**一声不响** ⇒ 读的是产品，不是断言写歪。
软链形（F）里产品把拒因打成宿主链接 `/varlink`，但红仍然落在"拼写对不上"那一腿（want `…/homereal119/…`、got `…/homelink119/…`）。

**它自述的"三段读数（红→绿→拆行再红）"在我这里塌成两段**，且两段我都自己复算了：
"只放新用例不放生产改动"（它的第一段）与"把那一行拆掉"（它的第三段）是**同一枚树形**，我的 E/F 就是它；
第二段（补上那一行⇒两形 3/3 绿）＝我的 A/B。**没有一段是抽验它贴的日志。**

**旧仪器今天仍然瞎（这条差异才是"仪器先于改码"的复算）**：同一发拆行变异下
`-run 'TestAC[1-4]POSIX' ./internal/winsec/ ./internal/proc/`＝**RUN=31 PASS=31 FAIL=0 SKIP=0 rc=0 PANIC=0**（RUN G）
⇒ 第一轮那句"产品行为差一个 rc，判据仪器一条都不差"**今天反过来了**：差的那个 rc 现在有枚用例咬得住，旧的那把尺照旧看不见它。

### 1.2 条件 2 —— **已变成读数**

**静态**：两枚 119 测试文件里 `SealableRoot` 共 7 处命中，**逐行看全在注释里**
（`dataroot_symlink_119_other_test.go:18/41/43/109/224/283` ＋ `secret_dataroot_119b_test.go:24`）；
期望值改成问文件系统（`cleanSpelling119` `:112-119` 用的是 `filepath.EvalSymlinks`），声明值用字面拼写。

**读数**：我自己重跑第一轮那发 `MUT-119-R`（本轮票面叫 `MUT-119b-T`），落地
`internal/proc/envfork.go:123: return SealableRoot(dir) /* MUT-119b-T-R2: also resolve the injected root */`，
`go vet ./internal/proc/ ./internal/winsec/` rc=0，然后
`-count=1 -v -run 'TestAC[1-4]POSIX' ./internal/winsec/ ./internal/proc/`（plain 形）＝

```
### SUMMARY H-mutT-plain rc=1 RUN=31 PASS=30 FAIL=1 SKIP=0 PANIC=0
--- FAIL: TestAC2POSIXInjectedTestDataDirStandsAsDeclared119
    dataroot_symlink_119_other_test.go:299: AC#2 RED: an injected data root was rewritten:
      got "…/001/injreal119/harness/picked", want the declared "…/001/injlink119/harness/picked"
```

⇒ **唯一那一枚红名就是第一轮说"抹掉纪律也照样绿"的那枚**，红在"逐字返回"那一腿（`:299`）。
另加一发控制：同一枚变异下 `./cmd/wisp/` 那三枚**仍 3/3 绿**（RUN I，rc=0、PANIC=0）
⇒ 两枚读同一份 OS 答案的生产调用者**各自独立被钉**，一枚退回去不会连带。

⚠ **口径必须连版本一起引**（我自己这条规矩写进正文，免得下一位拿它当装饰）：
第一轮"`MUT-119-R` ⇒ 56/56 全绿"那发是在 **`8c8aad3`** 那版树上量的（第一轮表 §一条件 2／§五表）。
今天同一包集匹配 `TestAC[1-4]POSIX` 的已经有 **31 条**（票 125/137 又各加了若干枚），
所以我**不是**拿 31 去比 56，我问的是"**这一发在 3a49745 上今天响不响**"——响了、且只响这一枚。

### 1.3 条件 3 —— 性质未变，也不构成剩余缺口

它是**跑者的性质**（"自证挂载非空"不是一枚 Go 用例能断的东西），所以它**不由实现方的返修闭合**：
本轮 119 仍**没有任何一份日志入 `docs/evidence/s1/`**（现量：那目录里以 `119-` 开头的只有第一轮那张表＋本文件）。
我的处置与第一轮同一条：AC#1 的结论**建立在我自己每一发的挂载自证上**（每发容器内先 `ls -l …/go.mod` ＋五枚被验文件 md5，
日志留在 `D:/tmp/wisp119r2-results/`），因此这一枚**不需要**它补交日志才成立；
但它自述的那些历史读数（真二进制 dev 三形 rc=1→0、门禁 `-count=2` 四数）本轮我**没有复算** ⇒ 在本表里只值〔日志＋归档，抽验〕，见 §5。

### 1.4 顺带把 AC#1 的字面再走一遍（"断今天的真实结局"）

票面 AC#1 要的是"把误伤面做成**可重跑的用例**、断**今天的真实结局**（红就是红）"。
今天（`3a49745`）的真实结局＝**这 8 枚（winsec 5 ＋ cmd/wisp 3）在两形下全绿、0 SKIP**（我的 RUN A/B/C/D 逐名可查）；
harness 那一族红（`R-119-8`）**不在本格的清除射程**（第一轮 §二(a) 已裁过一次，本轮不重开）。

⇒ **AC#1：成立，无附条件。**

<!-- 下面两行是 §2 那枚 commit 的读数（提交后补进正文，见文首"提交方式"） -->
