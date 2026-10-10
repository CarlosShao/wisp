# 298-r1 — AC#2 换行符假枚的标注写法 ＋ AC#3 "gofmt 那格在量谁"的票面措辞

## AC#2 假枚：本格的判据＝"写清楚怎么让它不误导人"，⛔ 清洗

### 本腿做了什么、⛔ 没做什么

- ⛔ 未 `git checkout -- cmd/wisp/models.go`（票面禁区＋AGENTS.md「禁 `checkout .` 族」；共享工作树里会吞别人的活）
- ⛔ 未删文件重写、⛔ 未 `gofmt -w`／`gofumpt -w` 它、⛔ 未动它的任何一个字节（`git diff --numstat` 除那两枚 `_test.go` 外**无任何 `.go` 路径**）
- ⛔ 未动 `.gitattributes`、⛔ 未动任何 git 配置（`core.autocrlf` 是机主的配置，AGENTS.md「NEVER update the git config」逐字压着）

### ★"要不要真换成本机 LF"这一问——按票面具名交给编排者裁，本腿不裁

票面 `AC#2` 与编排者更正节都写着这一问归编排者，且更正节里编排者已给本轮判断＝**本轮不换**
（理由逐字引票面：机主的 `core.autocrlf=true` 是他的 git 配置；逐文件 `.gitattributes` 改动属新射程、不在本票欠账里）。
本腿只把代价与收益重述一遍供裁：代价＝一次全内容重写（334 行全行改动、`git status` 会出一枚大 diff）；
收益＝名册少 5 枚假账。⛔ 本腿不替编排者执行。

### 今后可复用的标注写法（照票面现量节那三把尺的形式，逐字可重跑）

**规则：任何 `gofmt -l`／`gofumpt -l` 的名册，命中一枚件时先问一句"它是 blob 脏还是本机换行脏"，
判据＝下面三把尺同跑；三把同时给出 `i/lf w/crlf` ＋ `CR＝行数` ＋ `blob 那一把不命中` ⇒ 在名册里标成
`<路径>（换行符假枚：i/lf w/crlf，CR=<n>，blob 尺不命中）`，⛔ 计入格式债枚数。**

```
# 尺 1（库内是什么）  git ls-files --eol <路径>          ⇒ 期望 i/lf w/crlf attr/text eol=lf
# 尺 2（工作树 CR 数） tr -cd '\r' < <路径> | wc -c        ⇒ 与行数相等即"整文件 CRLF"
# 尺 3（blob 脏不脏）  git show HEAD:<路径> > $T/x.go; gofmt -l $T  ⇒ 不命中即⛔ 格式债
```

本腿 2026-10-10 16:1x–16:2x 现跑的同形读数（⛔ 引用前先重跑）：

| 尺（逐字） | 读数 |
|---|---|
| `git config core.autocrlf` | `true` |
| `gofmt -l cmd/wisp`（工作树） | 改前 3 枚 → 改后 **1 枚**＝`cmd\wisp\models.go` |
| `git show HEAD:<三枚>` 落仓外目录后 `gofmt -l <目录>` | 改前 **2 枚**（无 `models.go`）→ 改后 **0 枚** |
| `git ls-files --eol cmd/wisp/models.go` | `i/lf w/crlf attr/text eol=lf` |
| `tr -cd '\r' \| wc -c`（工作树） | `models.go` CR=**334**＝其行数；那两枚 `_test.go` CR=0 |
| `git ls-files --eol cmd internal tools \| grep -c "w/crlf"` | **5** |
| `git status --porcelain -- cmd/wisp/models.go` | 空（内容归一化后与 blob 相同） |

### 假枚名册（本票射程外，只标注；＝编排者更正节那 5 枚，本腿逐枚验过 `i/lf w/crlf`）

`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／
`internal/risk/provenance.go`／`internal/tools/bridge.go`
尺＝`git ls-files --eol` 那五枚路径 ⇒ 五枚全部逐字 `i/lf    w/crlf  attr/text eol=lf`；
`git ls-files --eol cmd internal tools | grep -c "w/crlf"` ⇒ 5（⇒ 该三目录里工作树 CRLF 的就这 5 枚，⛔ 更多）；
`gofumpt -l cmd/wisp internal tools`（工作树）⇒ **7** 枚＝上面 5 枚假枚＋本票那 2 枚真债（与票面更正节的 7 对得上）。
⇒ 本腿⛔ 碰这 5 枚里的任何一枚（`internal/**` 根本不在本票写面上）。

---

## AC#3 "gofmt 那格在量谁"＝钉成可复用判据（本腿只追加到票 298 票面，⛔ 不改票 292 原句）

### 一句话判据 ＋ 逐字尺（本腿追加进票面的就是这一段，两处措辞逐字相同）

> **碰 Go 面的票在写 `gofmt -l`（或 `gofumpt -l`）那一格时，必须同时写明三件事：
> ① 射程目录（`cmd/wisp` 还是 `cmd/wisp internal tools` 还是全仓），② 量的是 **HEAD blob** 还是**工作树**，
> ③ 交件时把**两把的枚数并排**给出，并具名指出差额里哪些是本机 `core.autocrlf` 造成的换行符假枚（标注法＝AC#2 那三把尺）。**

尺（逐字，两把并排那把）：

```
gofmt -l cmd/wisp                                    # 工作树那一把（射程要照票改写）
B=$(mktemp -d); mkdir -p "$B/blob"
for p in $(git ls-files cmd/wisp); do git show "HEAD:$p" > "$B/blob/$(basename $p)"; done
gofmt -l "$B/blob"                                    # HEAD blob 那一把，落点⛔ 不落仓内
```

### 为什么必须这么写（四档全是本腿现跑，⛔ 转抄；原件＝`logs/ac3-four-rulers.txt`／`logs/ac3-wide-ruler-fixed.txt`／`logs/ac3-blob-wide-ruler.txt`）

| 尺（逐字） | 射程 | 量的对象 | 改前命中 | 改后命中（本腿现跑） |
|---|---|---|---|---|
| `gofmt -l cmd/wisp` | `cmd/wisp` | 工作树 | 3 | **1**（剩＝`models.go` 假枚） |
| `gofmt -l cmd/wisp internal tools` | 三目录 | 工作树 | 7 | **5**（＝AC#2 那 5 枚假枚，逐枚 `i/lf w/crlf`） |
| `git show HEAD:<path>` 落仓外再 `gofmt -l` | `cmd/wisp` 全部 114 枚 `.go` | **HEAD blob** | 2（票面现量，只取三枚） | **0**（`logs/ac3-four-rulers.txt` 末节：整目录展开后一把尺给空） |
| `git archive HEAD cmd/wisp internal tools \| tar -x -C <仓外>; gofmt -l <三目录>` | 三目录 658 枚 `.go` | **HEAD blob** | 2（票面更正节读数） | **0**（同树 `gofumpt -l` 给 **1**＝35r2 的 gofumpt-only 残留，见下） |
| `git ls-files '*.go' \| xargs gofmt -l` | 全仓 tracked | 工作树 | 29（编排者更正节） | **27**＝22 枚 `.scratch/**` 故意脏拷贝＋5 枚假枚，真债 **0** |

⚠ **最后那一把还带一个口径坑，本腿撞到了，具名留档**：同一命令实测 **`rc=123` 而不是 0**，
根因＝`.scratch/wisp/probes/185/c1/mut/fs_broken.go:4:1: imports must appear before other declarations`
（一枚**故意做坏**的变异夹具让 `gofmt` 报错，`xargs` 因此退 123）。
⇒ 复跑这一把时必须把 stderr 与 stdout 分开（`2>/dev/null` 取名册、`2>&1 >/dev/null` 看报错），
否则 stderr 那行会被 `wc -l` 计进名册（本腿第一次就这么多数了 1 枚：28 vs 真值 27）。
⇒ 这条与 `AC#3` 那句"射程目录必写"同属**那把尺的口径**，⛔ 不是代码差。

⇒ 四档之差最大 27 枚，**全是尺的口径差，不是代码差**。不写"射程目录＋blob 还是工作树"，
下一位会在任意一档上误报一次——票 292 的 `AC#4` 那句「`gofmt -l` 空」就是这一形（工作树那一把在 HEAD 上永远满足不了）。
⇒ **本票 `AC#3` 的落点＝票 298 票面追加 ＋ 今后派单**；"搬进下一批发单模板"是编排者的活，⛔ 不由腿做（票面逐字写着）。

### 顺手把"gofumpt 那格在量谁"也钉一句（本腿只量⛔ 不修）

`gofmt -w` 洗完之后，blob 那一把 **`gofmt` 已空**，但**同一棵树**上 `gofumpt -l` 还给 **1 枚**＝
`cmd/wisp/panel_transport_35r2_test.go`。它要的是 `type ( … )` 成组、复合字面量换行、相邻 `func` 之间补空行——
**这些是 gofumpt 的加严规则，`gofmt -w` 管不到，⛔ 本票射程（票面 `AC#1` 只许 `gofmt -w`）内不可修**。
尺＝`$(go env GOPATH)/bin/gofumpt.exe -d cmd/wisp/panel_transport_35r2_test.go`（原件＝`logs/post-gofumpt-d-35r2.txt`，93 行／5 hunk）。
**这是 HEAD 上既有债、⛔ 本腿新增**：改前 blob 同尺给 5 hunk／101 行，改后 5 hunk／93 行，
逐段对拉（`logs/gofumpt-residue-pre-vs-post.diff.txt`）差额只有两处＝文件路径头一行 ＋
本腿那一枚 `:804` 对齐 hunk 的正文（gofumpt 先前也报它、现在不报了）⇒ **gofumpt 的抱怨只减不增、增 0 枚**。
⇒ 归 `AC#4` 的"gofumpt 只许减少⛔ 新增"＝3 枚→2 枚，达成；**残留那一枚要不要交给 gofumpt 族另立一票＝编排者裁**。
