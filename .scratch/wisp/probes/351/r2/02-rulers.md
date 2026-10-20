# 票 351 落地腿 351-r2 件02＝台面配方⇄门禁四数⇄纪律格自证

## 一、台面（公共块 o·仓外导出树，落 `D:\tmp\wisp\20261020\351r2\`）

- 两棵 `git archive bf509bb8 go.mod go.sum internal cmd tools scripts frontend models docs`（八件套＋docs）
  ＋手工三件 `build/*.dll`／`third_party/sherpa-onnx/*.dll`／`frontend/dist`（都⛔ 跟踪）。
- ★单差异证明＝`diff -rq treeB treeL` **只 1 条**：`config_reload_223_test.go differ`；文件集 symdiff **0 行**；
  `run.go`/`config_reload.go`/`winsec.go` 两树 sha256 逐字同。
- 取出后先数 `.go` 证非空：treeB⇄treeL 各 `GO_FILES=777`／`CMDWISP_GO=139`（"0"⛔ 是读数）。
- ★取 DLL 的 PATH＝`cygpath -w` **只转两枚 DLL 目录**＋分号前置（`SHP;DLLP;$PATH`）；
  ★⛔ 整条 PATH `-wp`（公共块 o②：那形把 gcc 指错⇒CGO 掉 0⇒"零枚 `--- FAIL` 的整包死"）。
- ★`GOTMPDIR`/`TMP`/`TEMP`/`TMPDIR` ＝ `cygpath -w` **反斜杠短根**（每发独立目录、先 `mkdir -p`）。
- ⛔ 跑 `./cmd/wisp` 整包（唯一持跑者＝编排者那窗）＝只定向 `-run 'TestTicket223'`；一律 `-count=1 -timeout 30m -v`。
- 四数只从 `^=== RUN`/`^--- PASS`/`^--- FAIL`/`^--- SKIP` 取；退码＝重定向文件＋`echo RC=$?` 紧跟裸命令（⛔ 管进 `tail`）。

## 二、门禁四数（★只在共享工作树·跑前 `unset GOTMPDIR TMPDIR`＝A1027⑤ 定式）

| 门 | 逐字尺 | `rc` |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp/config_reload_223_test.go` | **0**（我的文件不在脏名册＝新字节 LF） |
| d22scan | `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`；cmd/ 143 枚 Go 全数在内＝射程正常） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`） |
| go vet | `go vet ./cmd/wisp/` | **0** |

- 另：`mutA`/`mutB` 两棵去牙树上 `go vet ./cmd/wisp/` 带牙⇄还原各一发，**四次都 rc=0**（见件01）。

## 三、纪律格（AC#4 形状照守）

- ★⛔ 翻框：票 351 框尺终量＝`grep -cE '^- \[[ x]\]'`＝**5 格**，`[x]`=3⇄`[ ]`=2＝**3 勾／2 未勾**＝我一枚没翻（`AC#2`/`AC#3` 仍开，归编排者凭非实现者判语翻）。
- ★产码字节变动＝0 的量证＝`git show -U0 4e946c88` 首字符统计：`48 枚 +`⇄`1 枚 -`，
  逐枚非注释的加行＝`awaitAudit`/`windowCount`×2/`awaitStdout`×2/`if`/`}`×2/`t.Fatalf`/`t.Errorf` ＝**全在 `_test.go` 函数体内**；
  `git show --name-only` 里 `.go`＝**1 枚**、其中**非 `_test.go`＝0 枚**。
- ⛔ 放宽任何既有断言（`:591` 那枚 `Contains(shown,"risk.permission_mode")` 的串与消息逐字保留，只换读源）；
  ⛔ `time.Sleep` 哄时序（新增 0 枚）；⛔ 裸 `t.Skip`（新增 0 枚）；⛔ 新文件⛔ 新 helper（⛔ 新导出名）。
- ⛔ `internal/winsec/**`（只读 blob）、⛔ `internal/agent/approval/**`、⛔ `internal/risk`⇄`internal/config`⇄`frontend`⇄`design`。
- Git＝`add`⇄`commit` 两处逐枚字面 pathspec（`cmd/wisp/config_reload_223_test.go`）；⛔ `add -A`／⛔ 从 index 反推；
  ⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；⛔ push；`-F` 用 Windows 形正文路径。
- 入库正文裸码位 U+00A7＝**0**（写码位名，不落该字符本身）；件面 U+00B7 我也剥掉（⛔ 与收件尺打交道的必要风险）。
- 证据件扩展名⛔ `.out`（全仓 `*.out` 被 ignore）＝本腿只用 `.md`/`.log`/`.txt`/`.tsv`；临时件只建⛔ 删。

## 四、读数与红句的取法（⛔ 事后 grep 配错行）

- `go test -v` 把**断言正文印在 `--- FAIL` 之前** ⇒ 件01 的行号取自 `=== RUN` 段内正文（`:585`/`:592`/`:403`/`:406`），⛔ 从 `FAIL` 行往后配。
- ★判据只认**到达数分子**（`windowCount()-before==1`、`stdout never carried` 这类在场分子），⛔ 拿"总通知数大"当到达＝票 351 实测反例（红发总数＜绿发）。
- ⛔ 机主常驻 `wisp.exe` 抢 `Ctrl+Alt+Q/M/P` 那族热键红归我：本腿射程（`-run TestTicket223`）里**没有**热键族名，红名册逐名只票 351 那枚；⛔ 停他进程。
