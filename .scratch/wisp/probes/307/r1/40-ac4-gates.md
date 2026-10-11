# 307-r1 · 40-ac4-gates —— 门禁逐把 rc、格式两把尺、逐笔越界

每把尺**自己落一行 `rc=N`**（写法＝`… >> 件 2>&1; rc=$?; echo "rc_xxx=$rc" >> 件`），件＝`logs/gates.txt`（逐字追加，⛔ `cmd | head; echo $?`）。
另两把在 `logs/commit-bounds.txt`（见 `41-boundaries.md`）。

## 1. 三把门禁

| 尺 | 射程 | rc |
|---|---|---|
| `GOFLAGS= go build ./...` | 全仓（`./...`） | **rc_build=0** |
| `go vet ./internal/risk/` | `internal/risk` | **rc_vet=0** |
| `sh scripts/d22scan.sh` | 全仓（自报 `internal/`=229、`cmd/`=39、`internal/tools/`=23、ban#8 `internal/`=528 Go 文件含注释与 `_test.go`、`frontend/`=85、`design/`=39） | **rc_d22scan=0**（`d22scan: clean - no D22 ban violations`，其正控 35 枚 PASS/0 FAIL/0 SKIP） |

⛔ 新增裸 `go func(`、⛔ 新造折法、**⛔ 在 `risk` 包里写 `filepath.Clean`／`filepath.Abs` 做判定**——
尺＝`grep -rn --include=*.go 'filepath.Clean\|filepath.Abs' internal/risk --exclude='*_test.go' | grep -v '^internal/risk/pathresolver.go' | wc -l` ⇒ **0**
（与 `307-a1` `20-landing-shapes.md` 那枚"该包非测试面 0 命中"同读数 ⇒ 本票没写下包内第一例）。

## 2. ★格式名册**两把并排**（各自射程写死）

| 尺 | 射程 | 读数 |
|---|---|---|
| `gofmt -l internal/risk` | **工作树**，目录＝`internal/risk` | **1 枚**＝`internal\risk\provenance.go`；`rc_gofmt_worktree=0` |
| `git archive HEAD internal/risk \| tar -x -C /c/Users/swq/tmp/wisp-307-r1-blob && gofmt -l /c/Users/swq/tmp/wisp-307-r1-blob` | **HEAD blob**，同一射程目录（`internal/risk`）导出到仓外 | **0 枚**（空输出＝有效读数）；`rc_gofmt_blob=0` |
| `$(go env GOPATH)/bin/gofumpt.exe -l internal/risk` | 工作树；尺路径逐字＝`D:\work\base\gopath/bin/gofumpt.exe`（`ls -la` 5,028,352 字节，⛔ 裸 `gofumpt`） | **1 枚**＝同一枚 `internal\risk\provenance.go`；`rc_gofumpt_worktree=0` |
| `$(go env GOPATH)/bin/gofumpt.exe -l /c/Users/swq/tmp/wisp-307-r1-blob` | HEAD blob 导出树 | **0 枚**；`rc_gofumpt_blob=0` |

两把之差**只有那一枚 `provenance.go`**，尺＝`git ls-files --eol internal/risk/provenance.go` ⇒ `i/lf  w/crlf  attr/text eol=lf`
＝派单点名的那 5 枚 `i/lf w/crlf` 假枚之一 ⇒ **⛔ 计入格式债、⛔ 动它**（本腿没碰这枚文件）。
本票写下的两枚文件（`syncdirs.go`、`syncdirs_alias_307_windows_test.go`）在**两把尺下都⛔ 上榜**，
`syncdirs.go` 的 eol 尺＝`i/lf w/lf`（编辑后工作树仍 LF，⛔ 清洗、⛔ 改 `.gitattributes`／`core.autocrlf`）。

## 3. 整包两台面各两发取交集（＝`AC#4` 那格的颜色账）

逐字在 `10-ac1-colours.md` §1：改前 aliased 12 枚红／改后 aliased 7 枚红（甲族）、default 两态都 `FAIL=0`；
**新增红＝0 枚**（两台面、同台面两发取交集，`comm -13` 尺）；**单发红过被交集藏住的那枚＝0 枚**，
但具名交了一枚**发间打印次数差**（`TestSyncFallbackNotDisarmableByWeakRoot` 的 `t.Errorf` 次数），⛔ 用"集合相同"糊过去。

## 4. 台面卫生（起手与收尾各一把）

- 起手（`00-anchor.md`）：`tasklist //FI "IMAGENAME eq wisp.exe"`／`balldebug.exe`／`mockllm.exe` ＝ 三把各 0 枚（`INFO: No tasks are running…`）。
- ⛔ 杀任何进程；⛔ 需要 sherpa DLL（`internal/risk` 档）；⛔ 跑 `cmd/wisp`（本票⛔ 需要它）⇒ 派单那条"照抄 `$PWD` PATH 进 clone＝`0xc0000135` 假绿形"本票**未接触**。
- ⛔ 需要真窗／真账号／真 OneDrive：新夹具全部在 `t.TempDir()` 下自建。

## 5. `AC#4` 的"逐笔越界"与⛔ 区间尺那一坑

本腿**现踩过**这一枚，具名报回：先用派单⛔ 的区间尺跑了一遍，结果把**编排者自己那两笔**数了进来——
共享工作树里 `cb7bb3c..HEAD` 之间落进了 `249753a4`（票面追加节＋A833＋HANDOVER）与 `858ee7f4`（编排者欠读数①②交件），
它们动了 `.scratch/wisp/issues/307-*.md`、`docs/reports/**`、`.scratch/wisp/probes/307/orch/**`。
⇒ 派单"`git show --name-only --format=` **逐笔**量越界（⛔ 区间尺——区间尺会把别家的笔数进来）"这条**被本腿实测验证一次**；
本腿的越界账**只按逐笔名册算**（逐笔输出＋判定＝`logs/commit-bounds.txt`、`41-boundaries.md`）。

## 6. 本腿⛔ 做到的（违例具名）

- ⛔ 勾票面任何一格 ⇒ **本票 `AC#4`/`AC#5` 的框本腿也⛔ 翻**：翻框要写 `.scratch/wisp/issues/307-*.md`，
  而本腿授权写面逐字＝`internal/risk/**` ＋ `.scratch/wisp/probes/307/r1/**`，票面文件⛔ 在其中。
  **这一处是派单内部的两条相撞**（"AC#4/AC#5 你可以勾" ⇄ "写面只有那两面"；且票面 `AC#4` 自己写着 commit "只落自家 `probes/307/<腿名>/**`"）。
  ⇒ 本腿按**窄的那条**执行（⛔ 越界写），把两格的尺全部交在件里，**翻框留给编排者／验收腿**。
- ⛔ push（`git rev-list @{u}..HEAD` 之类⛔ 跑过，本腿⛔ 运行过任何 push/fetch）。
- ⛔ `git add -A`／`.`、⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`、⛔ 在仓内建 worktree。
  仓外导出树 `C:\Users\swq\tmp\wisp-307-r1-tree-main\` **⛔ 是 git 仓**（`tar -x` 落出来的），⛔ 任何 VCS 操作。
- 临时件只建不删，从头就落 `probes/307/r1/logs/`，**仓根 0 枚**；件名⛔ 带 `.out`（尺＝`git show --name-only --format=` 各笔里 `grep -c '\.out$'` ⇒ 0）。
