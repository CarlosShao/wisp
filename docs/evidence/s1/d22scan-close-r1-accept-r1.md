# d22scan 三条件收尾批（`304aeec`..`6a5b321`）· 非实现者对抗验收 r1

**我是谁**：裁决者（非实现者），本轮唯一可写路径＝本件。零代码改动、零 push。
**被验交付**：`docs/evidence/s1/d22scan-gitignore-fix-close-r1.md`（267 行，实现方自陈）＋它的两枚代码差
`tools/d22scan/scan_test.go` +101/-0（`304aeec`）／`tools/d22scan/gitignore.go` +30/-7（`b7c06d2`，自称注释级）／
三枚只动表的追加枚（`cec5e78`/`38f340b`/`6a5b321`）。
**它的上一环**：`docs/evidence/s1/d22scan-gitignore-fix-accept-r1.md` §4（工作单，含那行"只 M1 ⇒ PASS"与末句"从此摘掉 `all` 那一味必红"）。
**本件的靶心**：那枚白盒钉 `holds("frontend/weird/inside.tsx", false)`（`scan_test.go:1632`）是**承重**还是**镜子**，
以及实现程顶回上游验收措辞这件事的**程序**账。

| 用途 | 锚点 | 现量方式 |
|---|---|---|
| 我开始时的 HEAD | `310816f94216af4320b3a8ffa423eea5c4e1028d`（短 `310816f`） | `git rev-parse HEAD` |
| 我量 §0 时的 HEAD | `310816f`（§0 全程未变，末次复查同为 `310816f`） | 每节开头/结尾各查一次 |
| 被测批次 | `304aeec` → `b7c06d2` → `cec5e78` → `38f340b` → `6a5b321` | `git log --oneline 304aeec^..6a5b321` |
| 我读到的台账账 | `A207`/`A214`/`A218`/`A221②`/`A223`（`A223` 当声明读，不当证据读） | `docs/reports/pending-and-issues.md:5969,5989` |
| 我的仓库外工作区（只建不删，全在 `D:\tmp\d22scan-close-r1-accept\`） | `head\`（`tools/d22scan` 的 pristine 拷贝，`git hash-object` 见 §1.1）／`bin\d22-post.exe`／`bin\d22-M1.exe`／`bin\d22-M4a.exe`／`bin\d22-M4aM1.exe`／`mut\M1` `mut\M4a` `mut\M4aM1`／`fuzz\`（两枚我自己写的探针，`zz_overread_probe_test.go` `zz_combo_probe_test.go`）／`probe\seed1`（§2 的真 git 台件）／`seedtree\` `dirprobe\`（§1 的输出对照树）／`step1-accept.log` `gate-accept.log` `vet.log` `status-before.txt` `status-after.txt` `out-*.txt` `dp-*.txt` `real-*.txt` `combo-fuzz.log` | `ls` 现量 |

**并行读数的边界（先说清楚，免得我的数被人当别人的数用）**：另一位实现程此刻正在 `internal/panel/**` 上作业。
本件**没有**跑过整树 `go test ./...`；`go vet ./...` 与 `go build ./...` 我跑了（rc 均 0，见 §0.3），
但那是编译，不是行为读数。§0.2 的 `ban #8 internal/=407` 是**我在 `310816f` 上量的那枚数**，
它的归属是 `internal/panel/l2_grant_boundary_test.go`（盘上 10:09 的那版），**不是本批造的**，我不追它的漂移。

---

## 0. 闸门重跑（全部我现跑，锚点 `310816f`）

### 0.1 step 1（正控，`sh tools/d22scan/runtests.sh -C tools/d22scan ./...`）

```
$ git rev-parse --short HEAD
310816f
$ sh tools/d22scan/runtests.sh -C tools/d22scan ./... ; echo "rc=$?"
rc=0
（末三行，逐字）
PASS
ok  	github.com/CarlosShao/wisp/tools/d22scan	32.050s
runtests.sh: OK - packages=[./...] top-level: PASS=29 FAIL=0 SKIP=0, === RUN=69, '[no tests to run]'=0
```

⇒ 派单里那对数 **PASS=29 / FAIL=0 / SKIP=0 / RUN=69 复现成立**（`28→29`、`68→69` 的"开工基线"那一半见 §6.1，我用另一把尺另量了一次）。
日志：`D:\tmp\d22scan-close-r1-accept\step1-accept.log`。

### 0.2 step 2（真树扫描，`sh scripts/d22scan.sh`）

```
$ sh scripts/d22scan.sh ; echo "rc=$?"
rc=0
（两处 d22scan.sh 标记都在，证明 step 2 确实跑在 step 1 之后）
d22scan.sh: positive control - runtests.sh -C tools/d22scan ./...
d22scan.sh: scan of /d/work/workspace/projects plans/Wisp
（真树读数，逐字 8 行）
d22scan: scope bans #1-5 internal/      examined 203 production Go files
d22scan: scope bans #1-5 cmd/           examined  22 production Go files
d22scan: scope ban #6 frontend/         examined  40 text files
d22scan: scope ban #7 internal/tools/   examined  18 production Go files
d22scan: scope ban #8 design/           examined  32 text files
d22scan: scope ban #8 frontend/         examined  40 text files
d22scan: scope ban #8 internal/         examined 407 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  39 Go files, comments and _test.go included
```

⇒ **八数 `203/22/40/18 · 32/40/407/39` 与派单给的、与交件件 §4.2 的第三列，三处同值；整道门 rc=0。**
日志：`D:\tmp\d22scan-close-r1-accept\gate-accept.log`（step 1 段与 §0.1 同值，`ok 32.050s` 那份）。

### 0.3 卫生三项

```
$ gofmt -l tools/d22scan        -> 空输出
$ gofmt -l .                    -> 空输出（整仓，不止那一枚目录）
$ go vet ./...                  -> rc=0，输出 0 行（D:\tmp\...\vet.log 行数=0）
$ go build ./...                -> rc=0
$ (cd tools/d22scan && go vet ./...) -> rc=0
```

### 0.4 测试跑完之后真树没有多出残渣（这一格同时属于 §2）

```
$ git status --porcelain > status-before.txt   # 20 行
$ go test ./tools/d22scan（经 §0.1 那把尺跑过一遍）
$ git status --porcelain > status-after.txt    # 20 行
$ diff status-before.txt status-after.txt      -> 空输出 ⇒ STATUS_IDENTICAL
```

那 20 行全是 owner 自己的东西：16 枚 `design/**` 未提交删除 ＋ 4 枚未追踪
（`design/doubao/01-ball-states.jpg`、`design/doubao/demo/lib/`、`design/doubao/demo/screenshots/`、`design/old/`）。
⇒ **`frontend/**` 与 `.gitignore` 在本轮任何一次测试运行里都没有被写过**（机制见 §2）。
