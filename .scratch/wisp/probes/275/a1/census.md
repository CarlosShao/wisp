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

## §1 AC#0 退出码名册 — 未做（锚笔之后跑）

## §2 AC#1 未改动仪器的读数 — 未做

## §3 AC#2 归因到那一枚 — 未做

## §4 顺手量到的三形代价 — 未做

## §5 门禁与还原自证 — 未做

## §6 未做完的格逐枚具名 — 本笔时点：§1 §2 §3 §4 §5 全部未做
