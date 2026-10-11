# 309-v1 · 00-anchor（非实现者终裁腿的锚件／写面声明）

腿＝`309-v1`（裁决者，⛔ 实现者；实现腿＝`309-r1`）。
裁决对象＝票 `.scratch/wisp/issues/309-assessorfor-empty-taskid-limb-silently-drops-r4-detector.md` 的 `AC#1`–`AC#5`（`AC#0` 编排者已裁，⛔ 重开）。

## 台面
- HEAD 起手＝`0d2260eb`（branch `dev`）；量完仍＝`0d2260eb`（本腿⛔ push／⛔ 翻框／⛔ 动产码与别人的测试）。
- 被裁的三笔（`git log --oneline` 现读确认）＝
  `1cd57ec1`（锚件）→ `16e86bae`（`bridge.go` Ⓑ 注释＋两枚 `_test.go`）→ 票尾 bullet 笔＝`b2dd7fb2`（＋其后的读数补交 `8ffa61c6`／`8bbffc50`）。
- `16e86bae..HEAD` 之间 **⛔ 任何 `.go` 变更**（尺＝`git log --format='' --name-only 16e86bae..HEAD -- '*.go'` ＝空）⇒ 我今天在 HEAD 上量的产码＝那腿交件时的产码。
- `internal/tools/bridge.go` 自 `16e86bae` 起**无人再动**（尺＝`git log --oneline 16e86bae..HEAD -- internal/tools/bridge.go` ＝空），
  且 工作树 ⇄ HEAD blob 同一枚 sha256＝`cb86eaa7053a0f00758d642c59a8f22b881bb606abe954e9a8c03ee2d4291167`。

## 本腿自己复跑过的承重读数（⛔ 一件⛔ 接受实现腿的数）
| # | 尺 | 射程 | 单位 | 本腿读数 |
|---|---|---|---|---|
| 1 | `go test -run Test309… -v`（无突变正控） | `internal/tools` | 用例 | 两枚全 PASS，rc=0 |
| 2 | 四形突变 Rig（A 摘肢／B 探测器改 nil／C 整枚守卫 off／D 肢取反） | `internal/tools/bridge.go` 工作树（瞬态，trap 还原） | 红线条号 | 四形**全红**，还原 `cmp` 逐字节一致 |
| 3 | 两形突变 Rig（E 摘 `b.prov == nil`／F 摘 `b.injected`） | 同上 | 红线条号 | 两形**全绿**＝仪器对这两肢瞎 |
| 4 | 枚数钉三形 Rig（1 字面量／2 声明＋赋值／3 别名） | `bridge.go`（被钉真扫的一枚文件） | constructions | 1 枚→红；2、3→绿＝钉对这两种形状瞎 |
| 5 | `find cmd internal -name '*.go' -not -name '*_test.go'`＝269 减 `internal/agent`＝24 | 仓外导出树外的同一把尺 | 枚 | **245**＝钉打印的分母（对得上） |
| 6 | `git diff -w --numstat 16e86bae^ 16e86bae -- bridge.go` | 单文件 | 行 | **32 3**（⛔ 我派单里写的"expect 0"） |
| 7 | 逐行性质尺：`git show -U0 … \| grep -E '^[+-]'` 滤 `^//` | 单文件 | 枚 | 变更 35 行，**非注释行＝0** |
| 8 | 代码身份尺：两版 blob 各剥全行注释＋空行后 `cmp` | 单文件 | 字节 | **IDENTICAL**（768⇄768 代码行） |
| 9 | `GOFLAGS= go build ./...`／`go vet ./internal/tools/`／`sh scripts/d22scan.sh` | 全仓／单包／全仓 | rc | 0／0／0（各自落行，见 `gates-rc.txt`） |
| 10 | `gofmt -l internal/tools` 工作树 ⇄ HEAD blob 导出树；`gofumpt -l` 同两档 | 两档＋宽域 `internal cmd` | 枚 | 0↔0；宽域工作树 4 枚（全＝`w/crlf` 假枚，⛔ 309 名下）、宽域 blob 0 枚；gofumpt blob＝1 枚＝票 298 `AC#5` 那笔既有债 |
| 11 | 导出树同名作差（`1ad48333`＝`16e86bae^` ⇄ HEAD） | `internal/tools` 整包，两棵各自 `mktemp -d` | 用例名集合 | **纯新增 2 枚／删除 0 枚／SKIP 0**；两枚 236 空 id 载体两棵树**都 PASS** |
| 12 | `grep -rn 'ToolRequest{' internal cmd --include=*.go \| grep -v _test.go` | 全仓 `internal cmd` | 行 | **3**（评论自称"the single … point"＝1）⇒ 新写的尺**第一发重跑就对不上** |

## 写面（三张，⛔ 其余）
1. `.scratch/wisp/probes/309/v1/**`＝本腿证据件（只 `.md`/`.txt`/`.log`；Rig 脚本存成 `.txt` 交 `sh` 跑，⛔ 在 `probes/**` 下建任何 `.go`）。
2. `docs/evidence/s1/309-v1-verdict.md`＝裁决表（AGENTS §1.5 指定落点）。
3. 票 309 末尾 **1 条** append-only bullet（1 加／0 删）。

## 瞬态突变的还原纪律
所有突变只打在 `internal/tools/bridge.go`（派单点名的许可面），每次跑前 `cp` 到 `logs/bridge.go.pristine`，
还原写在 `trap … EXIT` 里，每次还原后打 `cmp` ＋ `sha256sum` 两行；收尾 `git status --porcelain -- internal cmd`＝**0 行**。
尺与原始读数＝`logs/mut-rig.txt`／`mut-rig2.txt`／`nail-rig.txt`（脚本本体也在盘上，可逐字重跑）。
