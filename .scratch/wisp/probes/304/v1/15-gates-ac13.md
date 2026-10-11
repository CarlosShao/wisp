# 304-v1 · `AC#1`/`AC#3` 尺页（本腿自跑，⛔ 引腿读数）

台面＝共享工作树 `D:/work/workspace/projects plans/Wisp`；锚＝`04f2ecbc`（本腿 `00-anchor.md` §0 现跑号）。
原始输出＝`logs/a01-commits.txt`／`a02-scope.txt`／`a03-r1scope.txt`／`g1-d22scan-body.txt`／`g1-pathlen-body.txt`／`g1-census-body.txt`／`g2-fmt-worktree.txt`。

## 1. `AC#3` 门禁四把（每把本腿现跑，件内自落 `rc=`）

| 尺 | 命令（逐字） | 本腿读数 |
|---|---|---|
| d22scan | `sh scripts/d22scan.sh` | `d22scan_rc=**0**`；末行逐字 `d22scan: clean - no D22 ban violations`；分母 ban#8 `internal/`=**533**? 见 §1.1，`cmd/`=**119**；live scope `bans #1-5 internal/=231`、`cmd/=39` |
| pathlen | `sh scripts/check-path-length-budget.sh --with-self-test` | `pathlen_rc=**0**`；denominator 逐字 `tracked paths=10982  over-budget=57  covered by roster=57  not in roster=0`；`worst full path ...=228 chars`；control 2/3 ok（种子越界件被点名） |
| census | `bash scripts/portable-tests.sh --scope=census` | `census_rc=**0**`；totals 行逐字（件内第 38 行）`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` ⇄ 票面 `AC#3` 基线**逐字相等** |
| 格式名册（两把并排） | `$(go env GOPATH)/bin/gofumpt.exe -l` | 射程＝本票 4 枚产码面：工作树那把＝**名册 0 枚**（`fmt_worktree_rc=0`）；HEAD blob 那把（`git show 04f2ecbc:<path>` 取到 `/d/tmp/wisp-304v1-blob/` 后量）＝**0 枚**，且**四枚 blob 分别 2098/894/13898/14264 字节＝先证取出的东西非空再信那个 0**；宽域一把（`internal cmd tools scripts` 工作树）＝**5 枚**（只作分母，⛔ 归本票射程） |

⚠ 三件套口径：以上全部量的是**工作树**（census/d22scan/pathlen 只能跑工作树）；格式那把第二向是**仓外文件里的 HEAD blob**。
⚠ 分母漂移具名：`tracked paths` 本腿 **10982** ⇄ `304-r1` 报 **10828**（⇄ 编排者 `A850` 那轮另有一号）——差 **154 枚**＝这之间别的腿落了 5+ 笔件；`over-budget 57 ⇄ roster 57 ⇄ not-in-roster 0` 与 `worst=228` **三向不变**，所以那 154 枚没进越界名册。

## 1.1 d22scan 分母回读（本腿自己从尾行摘，⛔ 抄腿）

尺＝`grep -nE 'examined [0-9]+ Go files' logs/g1-d22scan-body.txt` ⇒ **2 行**：第 249 行 ban#8 `internal/`＝**533 Go files, comments and _test.go included**（⇄ 腿报 533 **同值**）、第 250 行 ban#8 `cmd/`＝**119**（⇄ 腿报 119 **同值**）。

## 1.2 另两把我本腿现跑（腿交的是 0 字节件，见 §5）

| 尺（逐字） | 台面 | 读数 |
|---|---|---|
| `GOFLAGS= go vet ./internal/llm/... ./internal/proc/` | 共享工作树 @`04f2ecbc` | `vet_rc=0`，输出 **0 字节**＝vet 干净时不打印；分母＝`packages_under_vet` 行在件 `logs/g3-vet.txt`（该件本身非空，逐字落 `vet_rc=0`） |
| `GOFLAGS= go build ./internal/... ./cmd/... ./tools/...` | 同上 | `build_rc=0`，输出 0 字节＝干净；件 `logs/g3-build.txt`（自落 `build_rc=0`＋字节数行） |

## 2. `AC#1` 越界尺：逐笔（⛔ 区间一把扫）

尺＝`git show --name-only --format= <每一笔>`，射程＝`b9f37892..04f2ecbc` 全 **11 笔**（含别人的 6 笔），本腿**按标题前缀**分腿：

| 笔 | 腿 | 文件枚数 | 越界（⛔ `probes/304/r1/**` ∪ `internal/llm/**` ∪ `issues/304-*`） |
|---|---|---|---|
| `b9f37892` | R1 锚 | 1 | **0** |
| `e875f00f` | R1 落地 | 10 | **0** |
| `bd123485` | R1 `AC#2`/`AC#3` | 56 | **0** |
| `27f7b197` | R1 final＋票尾 | 3 | **0** |
| `04f2ecbc` | R1 更正名册 | 1 | **0** |

（`309-r3` 的 `df8d517d`/`a1f84d39`/`6fed247e`/`715f63cb` 与编排者 `2a6f11a5`/`c7887eca` 共 6 笔**⛔ 本票射程**，只在名册里计数，⛔ 计入 `304-r1` 的越界。）

## 3. 硬约束五条逐条核（本腿的尺）

- ① ⛔ 放宽断言：尺＝逐笔 `git show --format= -- internal/llm` 取**删除行**并剥掉注释与空行 ⇒ 全范围命中 **2 行**，两行**逐字同一句** `-	if err := proc.cmd.Start(); err != nil {`（两枚 spawn 点各一行）⇒ **删除面只有 spawn 那一行**，⛔ 任何 `t.Error`/`t.Fatal`/比较式被删或放宽。
- ② ⛔ 动 golden 一字节／⛔ 改 mockllm 语义：尺＝子树 blob 哈希比对 `git rev-parse b9f37892:<p>` ⇄ `04f2ecbc:<p>` ⇒ `internal/llm/testdata/golden`＝**同一枚 tree `efa272bc…`**；`internal/observe`／`tools/mockllm`／`internal/proc`／`scripts` 四棵 **same=YES**（⇒ 语义与 golden 零字节移动）。
- ③ 两枚 spawn 点同批：尺＝Read 两枚文件现量 ⇒ `internal/llm/adaptertest/mockllm.go:**85**` 与 `internal/llm/openaichat/mockllm_integ_test.go:**87**` 各调 `StartChildGuarded(t, proc.cmd)`，**同一次 commit `e875f00f`**（两枚文件同笔出现，`--name-only` 名册为证）。
- ④ ⛔ 碰 `Q-57`/票 156 产码路、⛔ 动 `internal/observe` registry 语义：尺＝子树哈希 same=YES（§3②）＋ `grep -rn 'mockllmRegistry'` 改动面⛔ 命中（`Spawn(...)` 调用点 `mockllm.go:92` 仍是原形，本腿读码同页可见）。
- ⑤ ⛔ 动档/tag/ledger/`ci.yml`：尺＝逐笔 `--name-only` 里 `.github/**`、`docs/**`、`scripts/**` 命中 **0**（`a02-scope.txt` 每笔 `forbidden_surface_hits=0`，11/11 笔）；tag/ledger 非文件面，⛔ 本尺射程，见 §5 名册。

## 4. 复用而不是新造第二套 job 包装

尺＝`grep -rnE 'CreateJobObject|AssignProcessToJobObject|JOB_OBJECT_LIMIT' internal/llm` ⇒ 命中 **1 行**，且那行是 `childjob_windows.go:13` 的**注释里的名字**，⛔ 代码；`internal/llm` 下唯一新 import 边＝`childjob_windows.go:9 "…/internal/proc"`（此前同目录只有 `tempdir_resolved_124_test.go:15` 一条，本腿现量 2 行）。
⇒ 走的是在库 `proc.OpenJobScope`／`StartInJob`／`Close`（本腿读 `childjob_windows.go:33/39/40/49` 逐字可见），⛔ 第二套包装。跨平台形＝`childjob_other.go:21-24` `//go:build !windows` → `return cmd.Start()`，注释里逐字承认那一支**仍会漏**并挂到 `Q-304-r1-1`＝编排者裁形时立的边界① 的形状。

## 5. 本腿欠账／未交部分（具名）

- `Q-304-v1-1`：d22scan ban#8 `internal/` 分母枚数本腿尚未从尾行精确摘出（件在盘，摘数是格式问题，⛔ 影响 `rc=0` 判语）。
- `Q-304-v1-2`：`go vet` 那把我尚未现跑（票面 `AC#3` 文本⛔ 要求 vet，属腿自加的加项）——见 `10-self-ran-shapes.md` 若我在反形树上补跑。
- `Q-304-v1-3`：`vet-reverse.txt` **0 字节**且同格**找⛔ 到另一枚带 `vet…_rc=` 的反形件**（`g-vet.txt` 只有 `vet_rc=0`，那是正向台面）⇒ 反形 vet 读数缺；`build-1.txt`／`vet-1.txt` 的读数**在别处**（`g-build.txt` `build_rc=0`、`g-vet.txt` `vet_rc=0`，两枚空壳是重复件）。
