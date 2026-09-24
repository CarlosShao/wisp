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

```
$ git log --oneline -1
9b29951 docs(119 复判r2,§1): AC#1 三枚条件逐条复算——自造 MUT-119-SEC 拆行两形各红 3 枚(红句点到用例自己种的链接)、自造 MUT-119b-T 唯一红名就是那枚恒真用例;判成立无附条件

$ git show --name-only HEAD
docs/evidence/s1/119-ac1-ac2-r2-acceptance.md
```

⇒ §1 落在 `9b29951`（同上一枚 §0 的 `6e04d1a`）；两枚 commit 的 `--name-only` 都**只有我这一枚文件**。
`internal/winsec/**` 本轮**一个字节都没写**（我是只读验收；另一程要求 `git diff b1010ff..HEAD -- internal/winsec/` 为空，我没有扰动它）。

---

## §2 AC#2 —— "winsec 一字未动"这句今天成不成立：**成立**（闭合方式是措辞对上了读数，不是读数变干净了）〔独立复现〕

派单给的判法我照走：①自己把返修那几枚 commit 逐枚 `--numstat` 量；②把"动了什么"拆成**判定分支／注释／其它**三档，
只有第一档算实质；③如果它这轮把措辞改对了，那也算闭合，但要**明写闭合的方式**。

### 2.1 返修四枚 commit 逐枚量（`git show --numstat --format="" <sha>` 原样）

```
36294c2   15   0	cmd/wisp/secret.go
        259   0	cmd/wisp/secret_dataroot_119b_test.go          ← internal/winsec: 0 枚
33c8acd   76  20	internal/winsec/dataroot_symlink_119_other_test.go
4f19ec6   71   0	.scratch/wisp/issues/119-...-symlinked-data-roots.md   ← 票面，非码
034080c    4   0	cmd/wisp/secret.go
          26   7	internal/proc/envfork.go
          41  14	internal/winsec/winsec_other.go
```

⇒ `internal/winsec/` 里票 119 本轮只碰了两枚文件：`winsec_other.go`（注释）与本票自己的测试文件。
`winsec.go`／`resolve.go`／`winsec_windows.go` **票 119 一枚 commit 都没碰**（现量：
`git log --oneline ce666ea..3a49745 -- winsec.go resolve.go winsec_windows.go` 只有 `a701138`(126)／`4824bb8`(125)／`a45b2e9`(129) 三笔，
逐笔 `--numstat` 全部只落 `resolve.go`＋它们自己的测试件）。

### 2.2 三档拆分（对 `internal/winsec/winsec_other.go`，控制组＝`ce666ea`，被验版＝`3a49745`）

| 档 | 读数 | 用的尺（四把，都不依赖票面那把） |
|---|---|---|
| **第一档：判定分支** | **0 行** | ① `git diff ce666ea..3a49745 -- winsec_other.go \| grep -E '^[+-]' \| grep -v '^+++\|^---' \| grep -v '^[+-][[:space:]]*//' \| wc -l` = **0**<br>② 两版各剥掉整行注释与空行后逐行 diff：`59` 行 vs `59` 行，**diff 无输出（rc=0）**<br>③ 两版 stripped 文件里 `if ancestorIsLink(prefix) {` 同在 `:39`、`func ancestorIsLink(prefix string) bool {` 同在 `:53`<br>④ 被当作"注释"剥掉的行里**没有 directive**：`git diff … \| grep -E '^[+-]' \| grep -c "go:build\|go:generate\|line directive"` = **0**；文件里也**没有块注释**（`grep -c '/\*'` = 0）⇒ "整行 `//` 才算注释"这一刀不会把 `//go:build` 混进注释档 |
| **第二档：注释** | **53 增 / 11 删**（＝64 行，全部是注释行） | `git diff --numstat ce666ea..3a49745 -- winsec_other.go` = `53 11`；分段：`189cb1e` = 21/6、`034080c` = 41/14（单枚 commit 的 41/14 与累计 53/11 不矛盾——后者含对前者已加行的改写） |
| **第三档：其它** | **0** | `winsec_other.go` 里没有既非判定、亦非注释的改动；包内另两枚生产文件本轮 0 笔（见 2.1） |

⇒ **与第一轮那 27 行（21/6）相比，"读数"没有变干净，反而更脏**：同一枚文件本轮累计动了 **64 行**。
第一档为 0 这件事，我用比第一轮更狠的尺（②③④）复算成立。

### 2.3 措辞那一半：本轮它把话改口了，且改得对

票面 `:299-308`（commit `034080c` 的 progress log）追加的更正块，**关键句逐字**：

> 一处更正（append-only，不改上面那段的字）："winsec 一字未动"这句按验收方换的尺打折，本轮把话改口
> … 但**票面上"winsec 是否一字未动"那一行的字面不成立**：真实读数到本轮为止是
> `internal/winsec/winsec_other.go` 相对票 119 控制组 `ce666ea` 共 **53 增 / 11 删，全部是注释行** …
> ⇒ 台账只能记"**判定未动，注释动了**"，不许记"一字未动"。

⇒ 它给的 **53/11** 与我 2.2 现量**逐字相同**；它给的"全部是注释行"与我第一档 = 0 相同。
**所以 AC#2 那句字面问题的正确答案现在是"判定未动、注释动了 64 行"，票面自己就是这么写的。**

**闭合的方式要说白**：**是措辞对上了读数**（票面不再主张"一字未动"，改为"判定未动，注释动了"），
**不是读数变干净了**（`winsec_other.go` 动的行数从第一轮的 27 涨到本轮的 64）。
⇒ 台账里这一格**不许**被下一个读者写成"winsec 一字未动"，只能写成"判定分支 0 hunks／注释 53 增 11 删"。

### 2.4 第一轮那句被证伪的全称（`R-119-2`）—— 现在换成了能站住的形状

```
grep -c "so a root handed to this floor names a real tree" internal/winsec/winsec_other.go   = 0（被证伪那句已不在）
```
替它的原文（`winsec_other.go:115-129`，逐字摘录要点）：
"Which roots are resolved today is **a list of this repository's callers, not a property this package can check**"，
把四枚调用方逐枚点名（`proc.TestDataDir`、`proc.DefaultLayout`、`cmd/wisp` 的 `resolveDataDir` 与 `resolveSecretLayout`），
并写"the honest claim is 'every data root this repository ships today is resolved before it reaches this floor',
**never** 'a root handed to this floor names a real tree': the next sealing site has to do the same and **nothing here notices if it forgets**"。
⇒ 与我在 2.5 量到的四枚生产调用点**一一对得上**（不是又一句全称）。

`R-119-10`（层级写坏）也复算了：`grep -c "^//   - "` ⇒ 控制组 `ce666ea`=**2**、票 119 锚点 `189cb1e`=**5**（标题写 "Two costs" 却五项）、
被验版 `3a49745`=**2**，硬链接那段挪成独立一段（`:141-145`）⇒ **它自述的三个数我全部复现**。

### 2.5 第一轮真正让 AC#2 退回的那一半（`②` 只做了一半）—— **已闭合，且有我自己的读数**

四枚读同一份 OS 答案的生产调用点，现量在 `3a49745`：

```
internal/proc/envfork.go:123-125  TestDataDir        ：注入值逐字返回 / os.TempDir() 经 SealableRoot
internal/proc/envfork.go:245      DefaultLayout      ：SealableRoot(os.UserConfigDir 的答案)
cmd/wisp/doctor.go:263            resolveDataDir     ：base = proc.SealableRoot(base)
cmd/wisp/secret.go:148            resolveSecretLayout ：root = proc.SealableRoot(root)   ← 本轮补上的第四枚
```

它是不是"补了一行、没人看着"？—— §1.1 那两发已经否掉了这个读法：
拆掉 `secret.go:148` ⇒ `./cmd/wisp/` 三枚在两形下**各红 3 枚**（E/F，产品原文点到用例自己种的链接），
同一发下旧仪器（winsec+proc 31 条）**全绿**（G）⇒ 这一行**既有产品效果、又有仪器看着**。

`R-119-3` 那一刀（声明树动不动）今天**两半各有一枚钉**：
逐字返回那一半＝RUN H 唯一红名 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`（红在 `:299`）；
位置契约那一半＝RUN A/B 两形 3/3 绿、RUN E/F 拆行即红。⇒ 原则句不再是一句只对一半的自律话。

### 2.6 一处**必须随本格一起记**的新鲜事实（不是票 119 的账，但会让下一个读者读错）

票面 `:307-308` 记的"三版 md5 相同"这句，今天在 `3a49745` 上**只剩两版成立**：

```
git show 3a49745:internal/winsec/winsec.go         | md5sum = a6144c880de80e43bb1393f3624e7221  = 1499efe 版 ✓
git show 3a49745:internal/winsec/resolve.go        | md5sum = b6876a5efe759f6e17434d1b50a129c3  ≠ 票面记的 7eb8a754eb3db1c8cf42a9dfeaa33074
   （7eb8a754… 在 034080c 那一版上仍成立——我现量它=7eb8a754…，即票面写它时是实话）
git show 3a49745:internal/winsec/winsec_windows.go | md5sum = fb20bca559d1725846dc2366123296b6 （相对它自己的控制组 ce666ea 也是这个值 ⇒ 票 119 没动）
```

⇒ 差异**全部来自票 126/125/129 三笔别人的 commit**（2.1 逐笔点数）。
这正是"归因腐坏"那族：**把带 md5 的事实句写在会被别人改动的文件上，读数还在、归属已经变了**。
本表按 `3a49745` 记账，并按 `A162` 的边界把它处理成"票 119 无关"（`A162` 明写不覆盖 `4824bb8` 之后同一枚文件的任何版本）。
⇒ 票 119 这一格**不因此退回**；但台账里**别把票面那句 md5 三连照抄成对 HEAD 成立的断言**。

### 2.7 AC#2 字面第二问（"该不该由调用方解析成实路径"＋指名两处）

第一轮已核"两处行号与形状属实"（`doctor.go:235`、`envfork.go:98` 都是 `filepath.Join(os.TempDir(), "wisp-test-<pid>")`），本轮不重开。
只补一条**行号腐坏**的读数：今天那条 test 路已经收成 `cmd/wisp/doctor.go:257: return proc.TestDataDir(), nil`
（重复实现被消掉那一处仍在），票面引的 `:235`/`:98` 现在指向别的句子。
⇒ 引这两枚行号时必须连"哪一版"一起引；**语义答案不变**（该解析，且只解析 OS 给的答案）。

⇒ **AC#2：成立。**退回的那半（② 只做了一半）已由 `36294c2`＋`33c8acd` 补完并被我量出牙齿；
"一字未动"那半以措辞对上读数的方式闭合；注释里那句被证伪的全称已换成点名四枚调用方的可核句子。

```
$ git log --oneline -1
fc48fc0 docs(119 复判r2,§2): AC#2 三档拆分——判定分支 0 行(四把尺)、注释 53 增/11 删、其它 0;措辞已改口故闭合方式是"措辞对上读数"而非"读数变干净";另记 resolve.go 的 md5 三连今天在 HEAD 已不成立(来自票126/125/129,非本票)

$ git show --name-only HEAD
docs/evidence/s1/119-ac1-ac2-r2-acceptance.md
```

⇒ §2 落在 `fc48fc0`（前两枚：§0 `6e04d1a`、§1 `9b29951`）。

---

## §3 顺手两问（各一段，不写长文）

### 3.1 你开的 `AC#7` 判据①对不对？—— **对，成立；我给了自己的读数，不替你圆**

判据①说：**两味药必须一起下**（接记名断言 **＋** 换已解析根），"只收紧不换根会在软链形恒红"。
我在被验版 `3a49745` 上把这两枚用例（`:193` 那枚的 `:203` 一腿、`:285` 那枚的 `:308` 一腿）在两形下各跑一遍
（`-count=1 -v`，用例自己把拒因打进日志），只看一句：**底线点名的链接是不是本用例种的那枚**。

```
plain 形（TMPDIR=/tmp/plain119r2）:
  :200  PrivateDirAll("…/001/varlink119/data")         → "…reaches it through the link at …/001/varlink119"     ← 用例自己种的
  :305  PrivateDirAll("…/001/injlink119/harness/picked") → "…the link at …/001/injlink119"                      ← 用例自己种的
软链形（TMPDIR=/varlink/w119tmp，ls -ld 现量 lrwxrwxrwx /varlink -> /realpriv）:
  :200  → "…reaches it through the link at /varlink"      ← 宿主的链接，不是它种的
  :305  → "…reaches it through the link at /varlink"      ← 同上
```

⇒ 软链形里这两枚**永远拿不到"点到本用例那枚链接"的拒因**（底线先撞见宿主的 `/varlink` 就停），
所以**只加记名断言、不换根 ⇒ 未变异也红 ⇒ 恒红**，判据①那半句成立；两味一起下才是可伪的。

**同一发白捡的反向对照**：`:251` 那一腿（`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`，
根已由 `cleanSpelling119` 换成已解析形）在**两形下都点到本用例种的链接**
（plain：`…/001/data/out`；软链：`/realpriv/w119tmp/…/001/data/out`）
⇒ "换过根之后记名断言就有区分力"这一半也被同一次运行钉住；顺带说明你把 `:251` 排除在 `AC#7` 之外是对的（它两形都响、是 137 的对照组）。

⚠ 一句边界（不改变结论）：判据①成立于**它自己指定的修法**（同包 helper `refusalCreditsLink137` 那种"必须等于我种的这枚"的形状）。
理论上还有第三条路——把断言写成"点到 `base` 之下的第一枚链接"就不换根也能伪，但要另造 helper、与 137 的形状不一致，
且把"宿主链接先撞见"这族遮蔽留在原地。**所以在判据自己的口径下，"必须一起下"是对的。**

### 3.2 那三枚用例今天有没有 CI 分母？—— **要看是哪三枚：`cmd/wisp` 那三枚＝没有；`internal/winsec` 那五枚＝有**

⚠ **先纠派单的前提**：派单写"这三枚用例（`dataroot_symlink_119_other_test.go` 里那三枚）"——
**那枚文件里是 5 枚**（`internal/winsec/dataroot_symlink_119_other_test.go:127/:153/:193/:219/:285`）；
票面 `next=` 第 6/7 条与 `36294c2` 说的"本轮三枚新用例"在 **`cmd/wisp/secret_dataroot_119b_test.go`**。两本账分开答。

**（a）`cmd/wisp` 那三枚 ⇒ 没有 CI 分母（票面那句"没有腿"成立）**

| 环节 | 现量出处 |
|---|---|
| 三枚的名字与枚数 | 容器内 `go test -list 'TestAC.*119' ./cmd/wisp/` ⇒ count=3：`TestAC1POSIXSecretRouteSymlinkedHomeBecomesSealable119`、`…SymlinkedXDGConfigHome…`、`TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119` |
| 它们在 `cli` scope、不在 core | `scripts/portable-tests.sh:191-196`（`cli` 分支 `scope=(./cmd/wisp/)`）；core 名单 `scripts/portable-tests.sh:172-181` **无 `cmd/wisp`** |
| `cli` scope 只有一个人调用 | `scripts/wisp-cli-tests.sh:113`（`bash "$portable" --scope=cli`），并且 `scripts/wisp-cli-tests.sh:60` 是硬闸：`if [ "$(go env GOOS)" != windows ]` 就拒跑（`:65-66` 原文 "do not relax this line"） |
| CI 上唯一一步 | `.github/workflows/ci.yml:422` `run: bash scripts/wisp-cli-tests.sh`，属 `test-windows` job（`ci.yml:334` `test-windows:`／`ci.yml:422`；core 那步在 `ci.yml:288`，job 见 `ci.yml:224`） |
| **决定性一发：windows 上这三枚连文件都不存在** | 同一枚快照、同一容器：`GOOS=windows go list -f '{{.TestGoFiles}} {{.XTestGoFiles}}' ./cmd/wisp/` 里 `119b` 命中 **0**，`GOOS=linux` 命中 **1**（文件头 `cmd/wisp/secret_dataroot_119b_test.go:1` = `//go:build !windows`） |

⇒ **唯一执行 `cli` scope 的 runner 是 windows，而那三枚在 windows 构建里不存在 ⇒ 零分母**（与 `R-119-7`／`R-113-A` 同族）。

**（b）`internal/winsec` 那五枚 ⇒ 有 CI 分母（ubuntu）**

| 环节 | 现量出处 |
|---|---|
| 五枚的名字与枚数 | 容器内 `go test -list 'TestAC.*119' ./internal/winsec/` ⇒ count=5（三枚 `TestAC1POSIX…119` ＋ `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` ＋ `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`） |
| 包在 core scope 里 | `scripts/portable-tests.sh:180`（`… ./cmd/llmrecord/ ./internal/winsec/`）；`:105` 注释写明 "internal/winsec joined the core list for ticket 111 AC#9"；core 钉住的包集 `scripts/portable-tests.sh:149`、`:164-166` |
| core 那步在 ubuntu 上 | `.github/workflows/ci.yml:224-225`（`test-core:` / `runs-on: ubuntu-latest`）→ `ci.yml:288` `run: bash scripts/portable-tests.sh --scope=core` |
| 平台看得见这枚文件 | `GOOS=linux` 命中 1、`GOOS=windows` 命中 0（`internal/winsec/dataroot_symlink_119_other_test.go:1` = `//go:build !windows`） |
| 没有被台账请出射程 | `scripts/portable-tests.sh` 里 `119` 只出现在 `:105/:108/:110` 三行来历注释，NOT IN SCOPE 台账**无 119 行** |

⇒ **但**这台 ubuntu 的 `TMPDIR` 不是软链 ⇒ **只有 plain 那一形有分母**；本轮 E/F/H 所量的软链形今天仍只在容器里。
⇒ 这一条与票 125 结案挂着的那本账是同一笔（"命令面级 POSIX 无门"），别记成两笔。

```
$ git log --oneline -1
97a4e41 docs(119 复判r2,§3): 顺手两问——AC#7 判据①成立(两形自量:软链形那两枚的拒因点到宿主 /varlink、plain 形点到自己种的链接;:251 反向对照同发拿到);CI 分母按文件分开答(cmd/wisp 那三枚零分母,引 portable-tests.sh:191-196 + wisp-cli-tests.sh:60 + ci.yml:422;internal/winsec 那五枚在 ubuntu core 腿有分母,引 :180 与 ci.yml:288)

$ git show --name-only HEAD
docs/evidence/s1/119-ac1-ac2-r2-acceptance.md
```

⇒ §3 落在 `97a4e41`（§4 里我把它的 commit 号漏贴了一次，这一句是补正；正文一字未改）。

---

## §4 逐格总判（只裁两格；一句话理由）

| 格 | 第一轮 | **本轮** | 一句话理由 | 标签 |
|---|---|---|---|---|
| **AC#1** 误伤面做成可重跑用例 | 通过附条件 | **成立（无附条件）** | 三枚条件里两枚硬的（①第二条生产路零用例、②fixture 拿被测函数算⇒恒真）本轮**各自有了我自己重跑出来的读数**（MUT-119-SEC 拆一行⇒两形各红 3 枚、红句点到用例自己种的链接；MUT-119b-T⇒唯一红名就是那枚恒真用例，红在 `:299`）；第三枚是跑者性质，由我自己的挂载自证覆盖 | 〔独立复现〕 |
| **AC#2** 裁 ①/②/③ ＋ "winsec 一字未动" | 退回 | **成立** | 真正让它退回的那半（② 只做了一半、`cmd/wisp/secret.go` 那条生产路没被覆盖）已由 `36294c2` 补上第四枚调用点、并由 E/F/G 三发量出"产品效果＋仪器都看得见"；"一字未动"那半**以措辞对上读数的方式闭合**（票面改成"判定未动、注释动了 53/11"），**不是读数变干净**（动的行数 27→64） | 〔独立复现〕 |

**"附条件 vs 退回"的分界按派单口径走了**：两枚硬条件都是**我造得出来**的（各造了一发红），所以不写附条件；
本轮也没有任何一枚只剩自述 ⇒ 这两格都不给"附条件"章。

**票面六格与 `AC#7` 我一枚勾都没翻、`-done` 没改名**：本文件是唯一写件；
`internal/winsec/**` 一个字节未写（同机另一程要求 `git diff b1010ff..HEAD -- internal/winsec/` 为空，我没扰动它）。

```
$ git log --oneline -1
0dec286 docs(119 复判r2,§4): 逐格总判——AC#1 成立(无附条件)、AC#2 成立;两枚硬条件各由我自造一发读数闭合,"附条件 vs 退回"分界按"造没造出来"走

$ git show --name-only HEAD
docs/evidence/s1/119-ac1-ac2-r2-acceptance.md
```

⇒ §4 落在 `0dec286`。本轮五枚 commit：§0 `6e04d1a`、§1 `9b29951`、§2 `fc48fc0`、§3 `97a4e41`、§4 `0dec286`，
逐枚 `--name-only` 只有本文件（上面五行读数即证）。

---

## §5 我明确**没核**的清单（别拿去当已核）

1. **票面自述的历史读数我一条都没复算**：真二进制 dev 三形 `rc=1→0` 那张表、六包 `-count=2 -v` 的 `602/600/0/2`、
   宿主 Windows 的 `go vet`/`go test`、`gofmt`/`gofumpt` 空输出。
   本轮两格的结论**只建立在我自己跑出来的 A–I 九发之上**；那些历史读数值〔仅自述，不背书〕。
   ⚠ 其中 `gofumpt` 那把尺今天已不可复现：`A162⑤` 现量本机是 **v0.12.0**，票面那两枚 09-22 的读数是 v0.7.0 量的。
2. **AC#3／AC#4／AC#5／AC#6 一枚未裁**（派单令）。票面那四格仍是实现方自勾，本轮**既没确认也没否**；
   第一轮那张表给它们的"通过"不因本表而升级。
3. **AC#7 那一格**：我只判了它**判据①成不成立**（§3.1），未实现、未裁格、票面正文一字未动。
4. **harness 那一族红的条数**（`R-119-8`：83 条 / 整包软链形 17 条）今天到底几条，未复算（不在我两格射程）。
5. **`R-119-5`（POSIX 上"C26 在不在位"零用例）**：我只读到盘上原文
   `internal/winsec/resolve.go:243  return resolveProbeRoot(os.TempDir())`（票 125 落的那一刀），
   **没量它现在还是不是"零用例"，也没跑那条安装 ERROR 的两形对照**
   ⇒ 票面 `next=` 第 2 条"今天仍在"那句**我没核**，谁要用先自己复算（`A161`/`A162` 已改了它的地基）。
6. **`R-119-4`／`R-119-7`**：未碰。`cmd/wisp` 在 ubuntu 上那 19 枚 DPAPI 红本轮也没复算
   （我只跑 `-run` 过滤的三枚；它在 `3a49745` 上还是不是 19 枚，未知）。
7. **macOS 半面**：无 runner，测不了（与第一轮同一条仪器边界）。
8. **`-count=2` 那把"不缓存"尺**：本轮全用 `-count=1`（我要的是"名字＋颜色＋名册差集＋panic 计数"，不是四数门禁），
   所以本轮**不下任何缓存结论**。
9. **`design/**` 那 16 枚未跟踪删除、以及索引里那两枚不是我下的 staged deletion**：按派单令未过问、未还原、未提交、未删。
10. **宿主工作树里除 `internal/winsec/`、`cmd/wisp/`、`internal/proc/` 之外别的包**：我没查脏不脏（与本轮两格无关，
    且我所有读数都取自 `git archive 3a49745` 的纯净快照，不取工作树）。

<!-- 下面两行是 §5 那枚 commit 的读数（提交后补进正文，见文首"提交方式"） -->
