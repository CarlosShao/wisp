# 275-a1 — 只读普查腿证据件（票 275：tracked gofumpt 分母里坐着一枚故意种坏的样本）

- 腿：`275-a1`（只读普查腿，⛔ 零写入跟踪文件除本件与 `logs/`，⛔ 不改仪器，⛔ 不选形）
- 票面原文：`.scratch/wisp/issues/275-tracked-gofumpt-denominator-carries-a-deliberately-broken-sample-so-the-guarded-census-step-exits-2.md`（54 行，已整份读完）
- 来源件：`.scratch/wisp/probes/111/r5/ci-guards.md` §4b（`:138-160`，已读）
- 本件状态：**§0 已落地（起手锚笔）**；§1–§6 在锚笔之后逐格填。没测到的格一律写「未做」，⛔ 不写成通过。

---

## §0 起手锚（在任何长跑命令之前落地的一笔 commit）

### 0.1 时刻与锚点

```
$ date
Wed Oct  7 11:47:12 CST 2026

$ git rev-parse --short HEAD
8976ffab

$ git branch --show-current
dev
```

### 0.2 起手工作树状态（票 AC#5 要求的「起手＝终态」比对基准）

```
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
(空输出 —— 这七枚路径起手即干净)
```

⚠ 全仓 `git status --porcelain` 起手并非干净：**751 行**，内容＝别人在飞的改动，**不是我动的、我不提交也不还原**（`M .gitignore`、`M .scratch/wisp/probes/152/my152.py`、`M .scratch/wisp/probes/161/r6/logs/flip-*.txt`、`M .scratch/wisp/probes/268/v1/evidence.md`、`D design/**` 多枚、`M design/doubao/demo/*`、以及一批 `?? .scratch/commit-msg-*.txt`／`?? .scratch/ci-logs/*`／`?? -` 之类的临时件）。
⛔ 我没动过其中任何一枚；§5 终态自查用「我这两枚路径之外的一切保持起手原样」这一条来证。

### 0.3 分母现量（⛔ 不引用票面读数，自己复跑）

```
$ git ls-files '*.go' | wc -l
935

$ git ls-files '.scratch/**/*.go' | wc -l
286
```

⇒ 与票 §2 第 1／2 行（935／286）对上。

### 0.4 gofumpt 路径与版本（现量，⛔ 不假定在 PATH 上）

```
$ go env GOPATH
D:\work\base\gopath

$ ls -la "$(go env GOPATH)/bin/gofumpt.exe"
-rwxr-xr-x 1 swq 197609 5028352 Sep 23 22:23 D:\work\base\gopath/bin/gofumpt.exe

$ "$(go env GOPATH)/bin/gofumpt.exe" -version
v0.12.0 (go1.27.1)
```

⇒ 版本与要求的 **v0.12.0** 对上。`attrib.sh:138-140` 自己也是这样定位的（优先 `.exe`）。

### 0.5 CI 那一步现读逐字

```
$ sed -n '230,240p' .github/workflows/ci.yml
        # they carry readings: scripts/slo-freshness.sh P1 (:160-170) fails the nail
        # if the slo-full job body carries `if:` at all, and :38-45 of that script
        # states that a FAILED sample must not produce an artifact, because P3 ages
        # on "a valid sample happened". Registered in the ticket instead.
        if: ${{ !cancelled() }}
        run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only

      - name: go vet (module)
        run: go vet ./...

      - name: go vet (tools/d22scan module)
```

```
$ wc -l .github/workflows/ci.yml .scratch/wisp/probes/161/r5/attrib.sh
  941 .github/workflows/ci.yml
  440 .scratch/wisp/probes/161/r5/attrib.sh
  1381 total
```

⇒ `if: ${{ !cancelled() }}` 在 `:234`、`run:` 在 `:235`，`ci.yml` 总行 941，与票 §2 第 6 行对上。守卫那一支确在 HEAD 上（`git status --porcelain -- .github/workflows/ci.yml` ＝空 ⇒ 工作树＝HEAD）。

### 0.6 ⚠ 具名回报：票 §2 第 5 行与编排者指令的**行号引用对不上**（内容与行号，二者只有一位是对的）

编排者指令与本票 §2 第 5 行都指向 `attrib.sh:173-178`，并按要求现读了 `sed -n '168,182p'`：

```
$ sed -n '168,182p' .scratch/wisp/probes/161/r5/attrib.sh
        printf 'ignored'
    else
        printf 'untracked'
    fi
}

# ticket_of <repo-relative path> -> the <NN> segment, empty if the path is not a bench path
ticket_of() {
    case $1 in
    .scratch/wisp/probes/[0-9]*/*)
        t=${1#.scratch/wisp/probes/}
        printf '%s' "${t%%/*}"
        ;;
    *) printf '' ;;
    esac
```

⇒ `:168-182` 里**没有**守卫那一支，它是 `classify`/`ticket_of` 的函数体。守卫那支的**现量真行号是 `:338-344`**（`grep -n` 原文，仪器工作树＝HEAD、`git diff --stat HEAD` 空）：

```
$ grep -n "A_RC\|GOFUMPT\|gofumpt\|tracked-only\|exit 2" .scratch/wisp/probes/161/r5/attrib.sh   （抽相关行）
338:    [ "$TRACKED_GO" -gt 0 ] || { echo "attrib.sh: (A) ruler saw 0 tracked .go files - refusing to report 'empty' from a ruler that was handed nothing" >&2; exit 2; }
339:    [ "$QUIET" = 1 ] || echo "== (A) tracked tree: git ls-files -z '*.go' | xargs -0 $GOFUMPT -l   (tracked .go files handed to it: $TRACKED_GO)"
340:    A_RC=0
341:    A_OUT=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?
342:    if [ "$A_RC" -gt 1 ]; then
343:        echo "attrib.sh: (A) gofumpt exited $A_RC on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler" >&2
344:        exit 2
```

⇒ **判**：票 §2 第 5 行抄的**代码内容逐字正确**（`A_RC=$(git ls-files -z '*.go' | xargs -0 "$GOFUMPT" -l 2>&1) || A_RC=$?` → `if [ "$A_RC" -gt 1 ]; then … exit 2`，注释逐字也对得上），**只有行号 `:173-178` 是过期的**。
顺带第二处过期：`attrib.sh:26` 自己的注释写「measured: `git ls-files '.scratch/**/*.go'` = 40」，我现量是 **286**（见 §0.3）——注释层面的旧读数，不影响判据，如实登记。

### 0.7 另外两枚前提的自证（票 §2 第 3／4 行）

```
$ git cat-file -e HEAD:.scratch/wisp/probes/185/c1/mut/fs_broken.go && echo "EXISTS rc=$?"
EXISTS rc=0

$ git ls-files --error-unmatch .scratch/wisp/probes/161/r5/attrib.sh; echo "rc=$?"
.scratch/wisp/probes/161/r5/attrib.sh
rc=0        ⇒ 仪器本体确在跟踪集里（自指一环成立）

$ git status --porcelain -- .scratch/wisp/probes/161/r5/attrib.sh .github/workflows/ci.yml
(空)        ⇒ 二者工作树＝HEAD，AC#1 起手的仪器是**未改动**的
```

### 0.8 另需登记的口径事实（影响 §1 名册的可解释性）

`attrib.sh:341` 的尺读的是 **`git ls-files` 给的名字 + 工作树给的内容**。工作树里另有别人未提交的 `.go` 改动，所以：
「935 枚名册 × 工作树内容」这一把**不等于** CI 干净检出上那一把（111-r5 §4b `:154-155` 已把这条差异写明，我这轮复跑时同样受它约束）。
⇒ §1 我按指令对**这 935 枚名册**逐枚起 gofumpt（读工作树内容），并在 §1 里具名登记哪些枚是「别人未提交的件」。

---

## §0b 起手锚的落档位置（本笔不重复 §0 全文）

> §0 全文（date／HEAD=8976ffab／七枚路径起手干净／935／286／gofumpt v0.12.0 路径／`ci.yml:230-240` 现读／
> `attrib.sh:168-182` 现读＋**行号过期具名回报（真行号 `:338-344`）**／坏样本 EXISTS_AT_HEAD／仪器本体 TRACKED）
> 已在锚笔 commit `019b4045` 里逐字落档。下面直接进三格读数。
> ⚠ `census-2.md` 是本腿一次误路径产生的 §1–§6 暂存副本，内容已原样合入本文件下方，
> 副本按「临时件只建不删」留在原地（⛔ 未提交、无人引用），见 §5.1。

---

## §1 AC#0 —— 935 枚逐枚 gofumpt 退出码名册（复跑腿 §2 第 7 那把尺）

### 1.1 命令原文（长跑脚本，⛔ 退码不经管道）

脚本本体＝`.scratch/wisp/probes/275/a1/logs/ac0-roster.sh`（本腿新建，只写 logs 目录），核心循环逐字：

```bash
GF="D:/work/base/gopath/bin/gofumpt.exe"
git ls-files -z '*.go' > "$LOGD/tracked-go.zlist"
while IFS= read -r -d '' f; do
    out=$("$GF" -l -- "$f" 2>&1)
    rc=$?
    printf '%s\t%s\n' "$f" "$rc" >> "$LOGD/gofumpt-exitcodes.tsv"
    ...
done < "$LOGD/tracked-go.zlist"
# batch form (the shape attrib.sh:341 actually runs) - exit code taken WITHOUT a pipe
git ls-files -z '*.go' | xargs -0 "$GF" -l > "$LOGD/batch-gofumpt-l-tracked.txt" 2> "$LOGD/batch-gofumpt-l-stderr.txt"
echo "batch_xargs_rc=$?" > ...
```

遵守了尺子规矩：① 全量输出先落 `gofumpt-exitcodes.tsv`（935 行）再截要引用的行；② 退码一律**去管道**单取
（`rc=$?` 紧跟命令；batch 那把是 `xargs ... > file 2> file` 之后才 `echo $?`）；③ gofumpt 用现量绝对路径 `D:/work/base/gopath/bin/gofumpt.exe`（v0.12.0），⛔ 不假定 PATH。

### 1.2 原始输出（名册统计，⛔ 不是结论是读数）

```
$ cat .scratch/wisp/probes/275/a1/logs/ac0.progress
loop done i=935 Wed Oct  7 11:50:57 CST 2026
ALL DONE Wed Oct  7 11:50:58 CST 2026

$ wc -l < .scratch/wisp/probes/275/a1/logs/gofumpt-exitcodes.tsv
935

$ cut -f2 .scratch/wisp/probes/275/a1/logs/gofumpt-exitcodes.tsv | sort | uniq -c
    934 0
      1 2

$ grep -P '\t[1-9]$' .scratch/wisp/probes/275/a1/logs/gofumpt-exitcodes.tsv
.scratch/wisp/probes/185/c1/mut/fs_broken.go	2
```

非 0 的那**一枚**的 gofumpt 原文报错（`logs/gofumpt-nonzero-detail.txt` 全文）：

```
=== file: .scratch/wisp/probes/185/c1/mut/fs_broken.go
=== rc: 2
=== gofumpt -l combined output (verbatim):
.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations
=== first error line: .scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations
```

整把 batch 形（＝`attrib.sh:341` 真跑的那形）：

```
$ cat .scratch/wisp/probes/275/a1/logs/batch-rc.txt
batch_xargs_rc=123
single_broken_rc=2
935 .scratch/wisp/probes/275/a1/logs/gofumpt-exitcodes.tsv

$ cat .scratch/wisp/probes/275/a1/logs/batch-gofumpt-l-stderr.txt
.scratch\wisp\probes\185\c1\mut\fs_broken.go:4:1: imports must appear before other declarations
```

⇒ 链条逐环对上：单枚 rc=2 ⇒ 批量 xargs rc=123 ⇒ `attrib.sh:342` 的 `[ "$A_RC" -gt 1 ]` 命中 ⇒ `exit 2`。

### 1.3 AC#0 判

**票 §2 第 7 那把尺复跑＝对得上**：935 枚里**恰好 1 枚**退出码非 0，就是 `.scratch/wisp/probes/185/c1/mut/fs_broken.go`（rc=2）。
腿那句「恰好 1 枚」**不作废**，本票 §2 第 7 行可从〔仅腿报〕升为〔本腿 275-a1 已复跑，10-07〕。
⚠ 但同一把尺还读出一件腿没量过的事：**`gofumpt -l` 今天点名 32 枚跟踪文件**（见 §1.4）——"恰好 1 枚非 0"说的是**退码**，不等于"只有 1 枚有问题"。

### 1.4 附：batch 那把 `-l` 的具名清单（32 枚，本腿现量）

`logs/batch-gofumpt-l-tracked.txt`＝32 行；正斜杠归一版＝`logs/batch-gofumpt-l-tracked.fwd.txt`；产码区那 5 枚单列＝`logs/ac2-named-outside-scratch.txt`：

```
$ grep -c '^\.scratch/'  logs/batch-gofumpt-l-tracked.fwd.txt   → 27
$ grep -v '^\.scratch/'  logs/batch-gofumpt-l-tracked.fwd.txt   → 5 枚：
cmd/wisp/models.go
internal/agent/approval/pending_read.go
internal/agent/tools.go
internal/risk/provenance.go
internal/tools/bridge.go
```

⚠ **这 5 枚产码点名是本地工作树的 CRLF 假象，不是格式回归**（现量，⛔ 别当回归报）：

```
$ git config core.autocrlf
true
$ cat .gitattributes | grep go
*.go text eol=lf          ⇒ 仓库自己的规矩就是 .go 必须是 LF
$ wc -c cmd/wisp/models.go /c/Users/swq/tmp/275-a1-head-content/cmd_wisp_models.go
14213 cmd/wisp/models.go                       ← 工作树
13879 .../cmd_wisp_models.go                   ← 同文件的 HEAD blob
$ cmp ... → differ: byte 13, line 1
$ for f in cmd/wisp/models.go internal/agent/tools.go internal/risk/provenance.go \
          internal/agent/approval/pending_read.go internal/tools/bridge.go
  → CR=334 LF=334 / CR=225 LF=225 / CR=131 LF=131 / CR=128 LF=128（逐枚 CR 数==LF 数＝整文件 CRLF）
$ git diff HEAD --name-only -- cmd internal | wc -l
0        ⇒ git 因 text=eol=lf 做了归一，所以 status 看不见这枚差异（"干净"是归一后的干净）
```

⇒ 干净 LF 检出的真实读数见 §3.3（CI 侧只点 **18** 枚，全在 `.scratch/wisp/probes/**`，产码 0 枚）。

---

## §2 AC#1 —— 未改动仪器上的真读数（⛔ 不推理）

### 2.1 命令原文与退码（退码去管道单取）

```
$ sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only > logs/ac1-tracked-only.stdout 2> logs/ac1-tracked-only.stderr; echo "AC1_RC=$?"
AC1_RC=2
```

仪器未改动自证（跑前跑后各一次）：

```
$ git status --porcelain -- .scratch/wisp/probes/161/r5/attrib.sh .github/workflows/ci.yml     （跑前）
(空)
$ git status --porcelain -- .scratch/wisp/probes/161/r5/attrib.sh                              （跑后）
(空)        ⇒ 仪器从头到尾未被本腿改动
```

### 2.2 原始输出（stdout／stderr 全文，＝"末几行"就是这几行，因为尺在 `exit 2` 处断了）

```
=== stdout (verbatim, full)
== (A) tracked tree: git ls-files -z '*.go' | xargs -0 D:\work\base\gopath/bin/gofumpt.exe -l   (tracked .go files handed to it: 935)
=== stderr (verbatim, full)
attrib.sh: (A) gofumpt exited 123 on the tracked set (>=2 means a file would not parse) - the tracked tree is not readable by this ruler
```

⇒ **读数＝rc=2**，与 `attrib.sh:343-344` 那一支逐字对得上；分母自报 935（恒真检查活着：它确实被喂了 935 枚）。
⚠ 注意这一支在 `exit 2` 之前**不打印任何 `-l` 名单**（`:342` 早于 `:346` 的循环），所以 CI 里这一步的红**只有一行"读不下去"、零枚具名**——这是编排者裁甲/丙时要看的一枚形状事实：**今天这枚守卫把 32 枚（本地）／18 枚（CI 形）具名信息全吞在 `exit 2` 后面**。

### 2.3 正控（证明尺命中得了真名，⛔ 不是"永远红"的坏尺）

同一把尺、同一个 `--tracked-only` 形，只在输入集里去掉那一枚坏样本（方法＝§3.1 的仓外临时索引）：

```
$ export GIT_INDEX_FILE=C:/Users/swq/tmp/275-a1-idx-no-broken
$ sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only > logs/ac2-tracked-only-minus-broken.stdout 2> logs/ac2-tracked-only-minus-broken.stderr; echo "AC2_RC=$?"
AC2_RC=1
（stdout 首行）== (A) tracked tree: ... (tracked .go files handed to it: 934)
（stdout 末行）attrib.sh: RED (tracked-only) - tracked-dirty=32 files_dirty=32. ...
（中间 32 行 TRACKED-DIRTY 逐枚具名，原文见 logs/ac2-tracked-only-minus-broken.stdout）
=== stderr (verbatim, full)
(空)        ⇒ "读不下去"那一支不再响
```

⇒ 尺**命中得了真名**（32 枚逐枚点名＋各自最后一次 commit），也**不是恒红**（去掉 parse 失败源之后 `exit 2` 那一支就闭嘴了）。正控成立。

---

## §3 AC#2 —— 归因到那一枚（两步只差那一枚文件）

### 3.1 用的是哪一等价形（说清，⛔ 不含糊）

用的是**首选等价形＝仓外临时 git index**，⛔ 不是"名单去掉一枚直接喂 gofumpt"的退化形：

```
$ export GIT_INDEX_FILE=C:/Users/swq/tmp/275-a1-idx-no-broken       （仓外：C:/Users/swq/tmp/，⛔ 不在仓库目录内）
$ git read-tree HEAD
$ git rm --cached -q .scratch/wisp/probes/185/c1/mut/fs_broken.go
$ git ls-files '*.go' | wc -l
934                                     （真索引 935 － 那一枚 ＝ 934，只差这一枚）
$ git ls-files --error-unmatch .scratch/wisp/probes/185/c1/mut/fs_broken.go
error: pathspec '...' did not match any file(s) known to git     rc=1
$ ls -la .scratch/wisp/probes/185/c1/mut/fs_broken.go
-rw-r--r-x ... 11041 Sep 28 14:35 ...   ⇒ 工作树一字未动（文件仍在盘上，只是不进那把尺的输入集）
```

与 CI 那一步的**差在哪**（三条，逐条具名）：
1. **输入集只少那一枚**（934 vs 935）——这正是 AC#2 要求的唯一自变量。
2. `git ls-files` 走临时 index＝**HEAD 的名字集**；CI 走干净检出的 index＝同一套名字。⇒ 这一条**不是差异**（且 `git diff --cached --name-only | wc -l`＝0，说明本地 index==HEAD，两者名字集本就相同）。
3. **内容仍是本地工作树内容**（含 §1.4 那 5 枚 CRLF）；CI 读的是 LF 检出。⇒ 这一条是本地读数和 CI 读数的**真差异**，所以本腿另补一把 §3.3 的 CI 等价尺（仓外 LF 树）把它剔掉。

### 3.2 读数与判

- 本地形：**rc 从 2 变 1**，`exit 2` 那一支（`:342-344`）从此不响，stderr 全空。
- ⇒ **票 §1 的机制判对**：那枚坏样本就是 `rc=2`／「尺读不下去」这一发的**唯一**成因（934 枚里 0 枚 parse 失败，§1.2 名册 934 个 rc=0 已独立证明）。
- ⚠ **但它不转绿**——本腿按 AC#2 的"仍红就扩列真红因"这条老实在此**扩列**：
  - 本地第二发红＝32 枚 `TRACKED-DIRTY`，其中 27 枚在 `.scratch/wisp/probes/**`、5 枚产码（**产码那 5 枚是 CRLF 假象，见 §1.4**）。
  - CI 形第二发红＝**18 枚，全部在 `.scratch/wisp/probes/**`，产码 0 枚**（现量见 §3.3）。
  ⇒ **票 §1 那句"要么把分母口径修对"隐含的"修对就不红"不成立**：把甲（缩分母）／丙（降级）之外的那一枚坏样本处理掉之后，这一步在 CI 仍会因 18 枚 bench 文件而红（rc=1）。**这不是"本票成因判错"，而是"成因只覆盖 rc=2 那一发、没覆盖 rc=1 那一发"**——编排者裁形时⛔ 别把"消掉 rc=2"当成"这一步会变绿"。

### 3.3 补一把 CI 等价尺（LF 干净检出仿真，仓外）

```
$ git ls-files -z '*.go' | git -c core.autocrlf=false checkout-index -z --stdin --prefix=C:/Users/swq/tmp/275-a1-ci-sim/
checkout-index rc=0
$ find C:/Users/swq/tmp/275-a1-ci-sim -name '*.go' | wc -l
934                     （临时索引＝HEAD 减那一枚）
$ 逐枚 CR 计数（cmd/wisp/models.go 等 5 枚产码）→ CR=0 全清      ⇒ 这棵树才是 CI 会读到的字节
$ cd C:/Users/swq/tmp/275-a1-ci-sim && find . -name '*.go' -print0 | xargs -0 "$GF" -l > logs/ci-sim-gofumpt-l.txt 2> logs/ci-sim-gofumpt-stderr.txt
$ echo $?（去管道另取）
xargs_rc=0 (123 would mean a parse failure occurred)      ⇒ 无 parse 失败：那一枚确实是 rc=2 的唯一成因
$ wc -l < logs/ci-sim-gofumpt-l.txt
18
$ cat logs/ci-sim-gofumpt-l.txt
.scratch\wisp\probes\163\a1\main.go
.scratch\wisp\probes\174\c2\zz174c2_wiring_pair_windows_test.go
.scratch\wisp\probes\183\r2\mut\task-boxset-off.go
.scratch\wisp\probes\185\r1\mut-m1\hostpath_185.go
.scratch\wisp\probes\197\r1c\pre\subagent_197.go
.scratch\wisp\probes\197\r1c\pre\subagent_197_test.go
.scratch\wisp\probes\212\v1\mut\main-noq9.go
.scratch\wisp\probes\220\r1\prechange\prechange220_probe_test.go
.scratch\wisp\probes\222\v1\mutations\m11-budget-50ms-gated\subagent_222_test.go
.scratch\wisp\probes\224\v2\dialect_probe_test.go
.scratch\wisp\probes\231\v1\mutation\config_reload_branch_after_prefix.go
.scratch\wisp\probes\235\v1\mut\prereadoff-m3_222_test.go
.scratch\wisp\probes\241\v1\posctl\badly_formatted.go
.scratch\wisp\probes\241\v1\probe_v1_readings_test.go
.scratch\wisp\probes\259\r1\gofumpt-negctl\probe.go
.scratch\wisp\probes\33\p1\q1\main.go
.scratch\wisp\probes\33\p1\q2\main.go
.scratch\wisp\probes\33\p1\q3\main.go
$ cat logs/ci-sim-gofumpt-stderr.txt
(空)
```

---

## §4 顺手量到的甲／乙／丙 代价（⛔ 本腿不选形，只把料交出去）

### 4.1 分母与射程的现量

| 量 | 读数 | 出处 |
|---|---|---|
| `git ls-files '*.go'` | **935** | §0.3 / §1.2 |
| 其中 `.scratch/**` | **286**（30.6%，全部住在 `.scratch/wisp/probes/**`：`git ls-files '.scratch/wisp/probes/*.go'`＝286） | §0.3 / 本笔 |
| 其中 `*/mut/*.go`（甲的窄形射程） | **72** | `git ls-files '*/mut/*.go' \| wc -l` |
| 排除 `.scratch/**` 后的分母 | **649**（＝少看 286 枚） | `git ls-files '*.go' \| grep -vc '^\.scratch/'` |
| parse 不可读的枚数 | **1**（`.scratch/wisp/probes/185/c1/mut/fs_broken.go`，rc=2） | §1.2 名册 |
| CI 形（LF）今天被 `-l` 点名的枚数 | **18**，⛔ 产码 0 枚、全部 `.scratch/wisp/probes/**` | §3.3 |
| 本地形今天被点名的枚数 | 32（27 枚 `.scratch/**` ＋ 5 枚产码 CRLF 假象） | §1.4 / §2.3 |

### 4.2 甲（改仪器分母口径）的代价——现量数

- **缩分母会瞎掉多少真回归**：⛔ 不是"0"。**今天 CI 形被 `-l` 点名的 18 枚全数落在 `.scratch/wisp/probes/**`** ⇒ 分母排除 `.scratch/**` 之后，这 18 枚**一枚都不再被这把尺看**（＝题面要的现量数：**18**）。
- 但这 18 枚里**产码回归 0 枚**（`cmd/internal/tools/scripts` 区在 CI 形上是干净的）⇒ 甲瞎掉的是 **bench 侧**的可见性，不是产码侧。
- 甲瞎掉的 18 枚里至少**两枚是别的票亲手做的格式控件**（本腿读了文件头，逐字）：
  - `.scratch/wisp/probes/241/v1/posctl/badly_formatted.go`（票 241 的**正控**）：`func  Bad( ) int {`／`return   1`——故意不格式化；
  - `.scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go`（票 259 的 **gofumpt 负控**目录）。
  ⇒ 走甲＝这两枚控件从此对这把尺隐形；票 259 目录名直接叫 `gofumpt-negctl`，这条代价要按"控件被门瞎掉"记，不能只按"bench 垃圾"记。
- 甲的窄形（只排 `**/mut/**`，72 枚）现量：能消掉 `rc=2` 那一支（坏样本路径含 `/mut/`），但**这一步在 CI 仍红**——18 枚里只有 3 枚（`183/r2/mut/task-boxset-off.go`、`212/v1/mut/main-noq9.go`、`235/v1/mut/prereadoff-m3_222_test.go`）在 `/mut/` 下，其余 15 枚照点。⇒ 窄形**换来的是"红得具名"，不是"绿"**。
- ⚠ 交付条件提醒（本腿不裁，只复述票面与本腿读数）：甲要动 `attrib.sh`＝票 161 的跟踪仪器，按票 §3/§5 需先落**具名解冻 `A##`**，⛔ 腿不许自行解冻。

### 4.3 乙（动那枚样本）的代价——现量补充

- 该样本是**唯一**让 `rc=2` 那一支响的文件（§1.2：934 枚全 rc=0）⇒ 把它撤出跟踪集或改成"能解析"，确实能消掉这一发；但同 §3.2：**这一步在 CI 仍红 18 枚**，乙也不换来绿。
- 另现量一枚本腿没做的代价边界：票 185 的凭据链本腿⛔ 未审（那是 185 的活），只登记"它在 HEAD 上、它是票 185 的凭据"这一条票面说法本腿不复用（未验证）。

### 4.4 丙（读不下去降级成具名清单）的代价——现量补充

- 今天 `exit 2` 那一支在打印任何名单**之前**就退了（`:342` 早于 `:346`，见 §2.2 stdout 只有 1 行）⇒ **丙拿到的信息量比想象中多**：`-l` 的 18/32 枚具名清单本来就存在于同一次调用里，只是被守卫吞了；丙若只做"逐枚列不可解析文件名＋照常打印分母＋WARN"，日志会**新增 18 枚（CI 形）具名行**——这是"降级"顺带解掉的一枚旧遮蔽，编排者可当作丙的一枚正收益记。
- 丙的反面现量：走丙 ⇒ 这一步今天**以 rc=0 收尾而日志里躺着 18 枚未格式化跟踪文件**（CI 形），且票 §3 要求的"种一枚新坏文件⇒名册必须出现它的具名行"定向判据**目前不存在**（本腿在 `attrib.sh` 里 `grep -n 'self-test\|negative-control'` 只见 `--self-test` 模式与 `161/r5/negative-control/bad-sample.go`（未跟踪），⛔ 没有针对 `rc=2` 那一支的定向控件）。

---

## §5 门禁与还原自证（本腿这一支）

### 5.1 本腿写了什么（终态自查）

- 本腿只写这两枚路径：`.scratch/wisp/probes/275/a1/census.md`、`.scratch/wisp/probes/275/a1/logs/`（外加仓外临时件 `C:/Users/swq/tmp/…`：commit message、临时索引 `275-a1-idx-no-broken`、LF 仿真树 `275-a1-ci-sim/`、HEAD blob 抽取目录 `275-a1-head-content/`——⛔ 全在仓外，仓内零残留）。
- ⚠ 具名登记一枚本腿造的杂件：`.scratch/wisp/probes/275/a1/census-2.md`＝本腿一次误把 §1–§6 写到旁边路径的**暂存副本**，正文已原样合入本文件；副本按「临时件只建不删」**留在原地、⛔ 不入 commit**（本腿两笔 commit 的 pathspec 里都没有它）。它会让 `git status` 在 `.scratch/wisp/probes/275/a1/` 下多一枚 `??`，那是本腿留下的，不是别人在飞的。
- 本腿在 `logs/` 里建的件（全部只增不删）：`ac0-roster.sh`、`tracked-go.zlist`、`tracked-go.zlist.bytes`、`tracked-go.count`、`gofumpt-exitcodes.tsv`（935 行名册）、`gofumpt-nonzero-detail.txt`、`ac0.progress`、`batch-gofumpt-l-tracked.txt`（32 行）、`batch-gofumpt-l-tracked.fwd.txt`、`batch-gofumpt-l-stderr.txt`、`batch-rc.txt`、`single-broken-file-stdout.txt`、`single-broken-file-stderr.txt`、`ac1-tracked-only.stdout`、`ac1-tracked-only.stderr`、`ac1-rc.txt`、`ac2-tracked-only-minus-broken.stdout`、`ac2-tracked-only-minus-broken.stderr`、`ac2-rc.txt`、`ac2-named-outside-scratch.txt`、`ci-sim-gofumpt-l.txt`（18 行）、`ci-sim-gofumpt-stderr.txt`、`ci-sim-rc.txt`。
- 未被本腿改动的东西（逐枚自证，命令原文）：

```
$ git status --porcelain -- .scratch/wisp/probes/161/r5/attrib.sh .github/workflows/ci.yml
(空)
$ git status --porcelain -- cmd internal scripts tools .github docs frontend
(空)                       ← 与起手（§0.2）逐字相同
$ git diff --stat HEAD -- .scratch/wisp/probes/161/r5/attrib.sh
(空)
$ ls .scratch/wisp/probes/185/c1/mut/fs_broken.go
存在                      ← 坏样本一字未动（本腿只把它从"临时索引的输入集"里剔过）
```

- 起手 751 行别人的未提交改动（`design/**`、`.gitignore`、`.scratch/**` logs、`?? -` 等）：**本腿没提交、没还原**，原样留在盘上（见 §0.2）。
- ⛔ 本腿没 push、没 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`restore`；两笔 commit 都带显式 pathspec 字面量。

### 5.2 本腿**没跑**的门禁（具名，⛔ 不写成通过）

| 门禁 | 票面归属 | 本腿状态 | 原因 |
|---|---|---|---|
| `sh scripts/d22scan.sh` rc=0（含先跑正控） | 票 275 **AC#5** | **未跑** | AC#5 不在本腿三格（AC#0/1/2）范围内；且它要起 Go 构建，⛔ 与另一枚写腿 `274-r1` 的构建/门禁读数抢 CPU（指令 §3 明禁抢整包） |
| `go vet ./cmd/wisp/` rc=0 | 票 275 **AC#5** | **未跑** | 同上（整包 Go 读数归 274-r1） |
| 票 161 自检载具 `portable-tests-selftest` 仍绿＋正控仍红 | 票 275 **AC#5** | **未跑** | 触发条件是"本票动 `attrib.sh`"；本腿⛔ 未动仪器，故本条待选形落地后再验 |
| CI 侧"这一步真跑过且颜色如何" | 票 275 **AC#5** 末条 | **未取数** | 必须**推送后**取数，机主 10-07 已说"暂时不推远程"⇒ 〔待推送取数〕，⛔ 本腿一切读数都是本机读数，不得冒充 CI 读数 |
| `go test ./...` | — | **故意不跑** | 指令 §3 明禁（洗 274-r1 的数） |

---

## §6 未做完的格（逐枚具名，⛔ 不留空话）

1. **票 275 AC#3 零误伤**（语法正确但未格式化的 tracked 样本经 `-overlay`／临时索引仍被看见）——⛔ **本腿未做**。本腿只做了它的"上半"：AC#2 那次正控（§2.3）已证明这把尺命中得了**已存在**的 32/18 枚真名，但**没有另造一枚新样本**去证明"选形落地后仍看得见"。归后续腿。
2. **票 275 AC#4 恒真性进攻**（指名哪一发突变让落地判据响）——⛔ **本腿未做**（要选形之后才有判据可攻；本腿只在 §4.4 登记了"丙 缺定向控件"这一枚事实）。
3. **票 275 AC#5 门禁**（d22scan／go vet／票 161 载具／CI 侧颜色）——⛔ **本腿未做**，理由逐条见 §5.2。
4. **§2 那把"935 枚逐枚退码名册"的第二把（HEAD blob 内容形）**——⛔ **本腿未做**：本腿的名册读的是**工作树字节**（`gofumpt -l <path>`），⛔ 没做"934 枚逐枚 HEAD 内容"的名册（只在 §1.4 对 5 枚产码抽了 HEAD blob、§3.3 对整棵树做了 `checkout-index` 批量形）。⇒ 若编排者要"逐枚 HEAD 内容退码名册"，那是新的一把尺，本腿没交。
5. **票 185 的凭据链**（那枚样本到底被哪些判据依赖）——⛔ **本腿未审**，只在 §4.3 登记为未验证。
6. **`git log` 具名 `f6b79ab0` 把守卫送进 HEAD** 这一条——⛔ **本腿未复跑**（票 §2 第 6 行的 commit 归属那一环，本腿只验了"守卫现在在 HEAD 上"＝`ci.yml` 工作树==HEAD＋`:234` 现读，没验它是哪一笔 commit 带进来的）。
7. **甲/乙/丙 三形的选择**——⛔ **本腿不选**（指令 §0/§3 明令，选形归编排者）。
