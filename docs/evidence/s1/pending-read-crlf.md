# pending-read-crlf — 一枚 `gofumpt -l` 残渣的结清（行尾归一，零逻辑字节）

- 腿：`pending-read-crlf`（**实现腿自报件，不是裁决表**；按 `SPEC-12 §4.3`「裁决者≠实现者」，是否结案的判语归非实现者另发）
- 锚：`5a21be431e70d698dbffd6b52ce57178f7337741`（分支 `dev`，**本腿零推送**）
- 时间：2026-11-05 16:38–16:50 +08
- 写面（仅此两枚）：`internal/agent/approval/pending_read.go`、`docs/evidence/s1/pending-read-crlf.md`
- 仪器：`/d/work/base/gopath/bin/gofumpt.exe` —— `v0.12.0 (go1.27.1)`；`go version go1.27.1 windows/amd64`
- 环境：Git Bash；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；**`CGO_ENABLED` 全程未设（⛔ 未用 `CGO_ENABLED=0`）**；`GOFLAGS` 未设

---

## 0. 一句话结论

那枚 `gofumpt -l` 红**不是格式问题，是行尾幻影**：仓内 blob（`i/lf`）本来就是 LF，是本机工作树那 131 个 CR 把它顶红了。
本腿把工作树那 131 个 CR 去掉之后，文件与 `HEAD` blob **逐字节相同**（同一个 git 对象 `a7db9b10…`），
于是 `gofumpt -l` 归零、`d22scan` rc=0、build 与测试 rc=0，而 **`git diff` 必然是 0 增 0 删**——
**这个"0/0"是正确结果，不是漏做**（机理见 §4，这是本件唯一与派单预期不同的地方，具名摆在这里）。

---

## 1. 起手读数（改前，现跑）

```
$ git ls-files --eol -- internal/agent/approval/pending_read.go
i/lf    w/crlf  attr/text eol=lf      	internal/agent/approval/pending_read.go

$ tr -cd '\r' < internal/agent/approval/pending_read.go | wc -c   → 131
$ wc -l < internal/agent/approval/pending_read.go                 → 131
$ wc -c < internal/agent/approval/pending_read.go                 → 6366
$ stat -c '%a %s' …                                               → 644 6366
$ md5sum （工作树 CRLF 形）                                        → acef5b1739e9446cd7118bbffa6a0a3b

$ git show HEAD:internal/agent/approval/pending_read.go | tr -cd '\r' | wc -c  → 0
$ git show HEAD:… | wc -c                                                      → 6235
$ md5sum （HEAD blob，LF 形）                                                   → 219a0465f27c469217227ef7a48b907d

$ git config --get core.autocrlf  → true
$ git status --porcelain -- 该文件 → （空）
```

**算式对得上**：`6366 − 6235 = 131`，而 CR 总数恰 131、LF 总数恰 131、**孤立 CR（不在 `\n` 前）＝ 0**
⇒ 工作树与 blob 的差别**只有 131 个"每个 `\n` 前面多出的 `\r`"**，别的字节一个没动。

## 2. 起手先证伪"是格式问题"

派单前提是"`i/lf w/crlf`"⇒ 必须先排除"LF 形其实也没格式化"这种混合病因：

```
$ tr -d '\r' < <工作树> > /tmp/pr_norm_check.go
$ gofumpt.exe -l /tmp/pr_norm_check.go   → 零行，rc=0        # LF 形本来就干净
$ gofumpt.exe -l <工作树 CRLF 形>         → 列出该文件，rc=0  # 只有 CRLF 形被点
```

⇒ 病因单一：**行尾**。

全仓残渣口径（`git ls-files '*.go'` 减去 `.scratch/**`，共 **1134** 枚）：

```
$ git ls-files -z '*.go' | tr '\0' '\n' | grep -v '^\.scratch/' | tr '\n' '\0' | xargs -0 gofumpt.exe -l
internal\agent\approval\pending_read.go        ← 改前：恰好这一枚
（改后：零行）
```

## 3. 动作（只去行尾，别的字节零改动）

```
$ cp internal/agent/approval/pending_read.go /d/tmp/wisp-crlf-5a21be431/pending_read.go.crlf.bak   # 只建不删
$ tr -d '\r' < internal/agent/approval/pending_read.go > /tmp/pr_lf.bin
$ cp /tmp/pr_lf.bin internal/agent/approval/pending_read.go      # 原地覆写，保留 mode/inode
```

改后读数：

```
CR=0  LF=131  bytes=6235  mode=644（未变）
md5sum(工作树) = 219a0465f27c469217227ef7a48b907d  == md5sum(HEAD blob)
$ cmp 工作树  <(git show HEAD:该文件)  → 无差异
$ git hash-object -- 该文件       → a7db9b104489b4e3b1a6b2e2e619cac8865f04e1
$ git rev-parse HEAD:该文件       → a7db9b104489b4e3b1a6b2e2e619cac8865f04e1   # 同一个对象
$ git ls-files --eol -- 该文件     → i/lf    w/lf    attr/text eol=lf
```

## 4. 与派单预期的唯一偏差（具名）

派单写「`git diff --stat` 应见约等量增删」。**本件实测是 0 增 0 删**，逐条读数：

```
$ git diff --quiet -- internal/agent/approval/pending_read.go ; echo $?   → 0
$ git diff --stat  -- 该文件 | wc -l                                       → 0
$ git diff --cached --stat -- 该文件 | wc -l                               → 0
```

**机理（不是行为省了，是结构上不可能有 diff）**：`.gitattributes` 写 `*.go text eol=lf`，而本机 `core.autocrlf=true`
⇒ git 的 clean filter 在比较／入库时**本就把 CRLF 归一成 LF**。所以：

1. 索引里的 blob **一直是 LF**（`i/lf`，CR=0）——这枚文件**从来没在 blob 层脏过**；
2. 脏的只是**本机检出后的盘上形状**（哪个工具把它重写成 CRLF 的，本腿未归因，也不猜）；
3. 于是把它改回 LF 之后，工作树与索引**逐字节相同** ⇒ `git diff` 空、`git add` 空。

⚠ **这正是本仓在册的那族"CRLF 幻影"**（`A663`／`A674` §3／`A1024` §②「四枚 CRLF 假脏…blob 干净、CI 不咬」同一族）：
`gofumpt -l` 吃的是盘上字节、CI 吃的是检出后的 LF 形。本件把"本机那一发"也弄干净了，**但它在 CI 侧从未红过**——
所以本件不宣称"修好了一枚 CI 红"，只宣称"本机 `gofumpt -l` 残渣归零"。

**本腿没有假造 diff**（⛔ 没有为了让 `--stat` 好看去改任何一个逻辑字节，也没把 CRLF 塞回索引）。
`git status` 改前空、改后空，中间那一次显示 ` M` 是**stat 缓存过期**（size 6366→6235），
`git update-index --refresh -- 该文件` 之后归零（`git diff --quiet` 全程 rc=0）。

## 5. 门禁（本腿现跑，逐条落 rc）

| # | 门 | 命令 | 读数 |
|---|---|---|---|
| G1 | `gofumpt -l <该文件>` 空 | `gofumpt.exe -l internal/agent/approval/pending_read.go` | 零行，**rc=0** |
| G2 | `gofumpt -l internal/agent/…` 空 | 94 枚 `git ls-files 'internal/agent/*.go'` 全量喂一把尺 | 零行，**rc=0** |
| G3 | 全仓非 `.scratch` 残渣 | 1134 枚 tracked `.go` | 零行（改前 1 枚） |
| G4 | `go build ./internal/agent/...` | 见 §0 环境 | **rc=0**（1 s，缓存命中） |
| G5 | `go test ./internal/agent/approval -count=1` | 同上 | **rc=0**，`ok … 0.784s`；`-v` 复量 **PASS 88 / RUN 113 / FAIL 0 / SKIP 1**，`ok 1.043s` |
| G6 | `sh scripts/d22scan.sh` | 同上 | **rc=0**，`clean - no D22 ban violations`（60 s；`bans #1-5 internal/=361 cmd/=57`、`ban #6 frontend/=116`、`ban #7 internal/tools/=48`、`ban #8 design/=39 frontend/=116 internal/=893 cmd/=208`） |

G4/G5 的 `1 s` 对 `0.784 s` 是本机缓存命中，**照记不修饰**（本仓在册：rc=0 可能只是"没干活"，故 G5 另发了 `-v` 那一发看名册）。

## 6. 残余与归属（⛔ 本腿未碰）

1. `.scratch/**` 里仍有 **11 枚** tracked `.go` 的**工作树**是 `w/crlf`（`probes/257/r2/mut/M1–M8/firstrun*.go` 8 枚、
   `probes/287/r1/mutants/stale-crlf/M1–M3*.go` 3 枚）。**一枚都没动**：它们名字自己写着 `stale-crlf` / `mut`，
   是**别的票的故意坏样本**，洗掉＝摘牙。生产面（非 `.scratch`）现在 **`w/crlf` 计数 = 0**。
2. `git update-index --refresh -- <该文件>` 这一发**顺带报出** `internal/panel/git_readonly_exec.go: needs update`。
   那是**共享工作树里别人的在飞件**（本腿写面外）；`--refresh` 只动 stat 缓存、不动内容，本腿**未 add、未 commit、未读其正文**。
3. 本腿**不翻任何 AC 框、不改任何票面、不推**。

## 7. 撤销口令

```
cp /d/tmp/wisp-crlf-5a21be431/pending_read.go.crlf.bak internal/agent/approval/pending_read.go
```
（该备份＝改前 CRLF 形，`6366` 字节、`md5 acef5b1739e9446cd7118bbffa6a0a3b`，只建不删。撤＝把它 `cp` 回去，`gofumpt -l` 会重新点它。）

## 8. 复跑（任何人可逐条重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
GF=/d/work/base/gopath/bin/gofumpt.exe
F=internal/agent/approval/pending_read.go
$GF -l "$F"                                            # 期望：零行
git ls-files --eol -- "$F"                             # 期望：i/lf  w/lf
printf 'CR=%s\n' "$(tr -cd '\r' < "$F" | wc -c)"        # 期望：CR=0
git diff --quiet -- "$F"; echo "diff_rc=$?"            # 期望：diff_rc=0
git rev-parse "HEAD:$F"                                # 期望：a7db9b104489b4e3b1a6b2e2e619cac8865f04e1
git hash-object -- "$F"                                # 期望：同上一行
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
unset GOTMPDIR TMPDIR TEMP TMP
go build ./internal/agent/...                          # 期望：rc=0
go test ./internal/agent/approval -count=1             # 期望：ok / rc=0
sh scripts/d22scan.sh                                  # 期望：clean / rc=0
```
