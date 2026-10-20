# 230-v1 · 00 尺册（每把尺：尺名｜命令形｜射程｜量还是算）
# 裁决腿＝票 230 终裁腿 230-v1（非实现者·只读·⛔ go build／go test／⛔ 动 .go 字节）
# 本机 core.autocrlf=true ⇒ 一切 CR／码位／剥注释尺一律对入库 blob 量，⛔ 对工作树量。

| 尺 | 命令形 | 射程 | 量/算 | 落档 |
|---|---|---|---|---|
| V1-a 名册 | `git log --grep='^票 230 落地腿 230-r1' --format='%h %ad %s'` | HEAD 起全史 | 量 | 01 |
| V1-A 可执行字节 | `git show -U0`＋awk 判 `^\s*//` | 笔①⇄笔② 入库 diff | 量 | 01 |
| V1-B 剥注释正文 | `git cat-file blob <rev>:<path>`｜`grep -vE '^[[:space:]]*//'`｜sha256⇄diff | `c16ff47e`⇄`c73df208`⇄`ca0ad6bd`⇄HEAD | 量 | 01 |
| V1-B2 blob 同值 | `git rev-parse <rev>:<path>` ⇄ `git hash-object`（仓外副本） | 同上⇄全 8 笔 | 量 | 01 |
| V1-C 真值名册 | `git cat-file blob HEAD:queue.go`＋`git grep … HEAD`＋逐枚判注释性 | queue.go⇄全仓 | 量 | 02 |
| V1-D enclosing | awk 追 `^func` 推宿主函数 | approval 包 | 算（对入库字节） | 02 |
| V1-E 格式两把 | `gofmt -l` ⇄ `gofumpt.exe -l`（仓外副本＝HEAD blob） | queue.go blob | 量 | 03 |
| V1-E2 正控 | 同两把对故意写坏的 `badctl.go` | 仓外 | 量 | 03 |
| V1-F misbound | `git grep 'denialMisbound'`＋`spend(` 逐枚实参 | approval 包 | 量 | 02 |
| V1-G2 锚正反控 | `git grep 'only reader of q.alias' <rev> -- <file>` | HEAD⇄`c16ff47e` | 量 | 03 |
| V1-P 拼接诚实尺 | 剥行首 `//` 后 `tr '\n' ' '` 再 grep 同一枚锚 | 两侧 blob | 量 | 03 |
| V1-I 同族射程梯 | `git grep -E '<腿的原样正则>'` 逐射程（queue.go⇄包⇄internal+cmd⇄*.go⇄全仓） | 五档 | 量 | 03 |
| V1-O 词界修正 | 同族正则加 `(^|[^[:alnum:]])` 剔 `con\|sole reader` 噪声 | internal+cmd | 量 | 03 |
| V1-K/U 先例⇄单变量 | 逐字重跑腿 `sh/17` 那把尺；`git diff --name-only 3eeaa2e2 c73df208` | issues 池⇄两枚 rev | 量 | 03 |
| V1-T 长度门 | `sh scripts/check-path-length-budget.sh` 裸形（纯 `git ls-files`＋awk，⛔ 调 go） | 共享树全 tracked | 量 | 03 |
| V1-Q 对质腿的尺 | 只读腿自己落档的 `sh/01-metric-exec-bytes.sh` 正文⇄`logs/metric-exec-bytes.txt` | 仓外临时根 | 量（旁证射程用，⛔ 采信其数） | 03 |
| V1-J/V1-R 写面⇄纪律 | `git show --format='' --name-only` 逐枚＋numstat＋票 230 前 5 行 diff | 全 8 笔⇄票 230 | 量 | 03 |
| 码位带一把 | python 按 `tools/d22scan` 声明的 band 表扫 queue.go blob | queue.go | 算 | 03 |

⛔ 本腿复跑的四把（硬约束：三枚写腿在 winsec/models/cmd-wisp 落笔，跑任何包都可能读到半保存态假红）＝
`GOFLAGS= go build ./...`、`sh scripts/d22scan.sh`、`go vet ./internal/agent/...`、定向色 `./internal/agent/...`。
⇒ 这四把的 rc 在判语 第 4 节／第 5.4 节 里明列为"以腿报为据、本腿⛔ 复跑"，⛔ 当本腿读数。

读法提醒（本腿自己踩到两次，具名登记，⛔ 下位再踩）：
- `cut -c`／`awk substr` 在 MSYS 下**按字节**截中文＝造非法 UTF-8（公共块 A1042③ 同族）。本腿第一发
  `v1-g-discipline.out.txt` 就因此三行不可解，改用 python 按字符截才复跑成 `v1-g2`。
- `git show --format=… --name-only` ⇄ `-s`／`--check` 同用＝`fatal: … cannot be used together`，
  名册尺会**静默数到 0**（本腿 V1-I 第一发就中了这条，改用 `--format='' --name-only` 才对，见 V1-J）。
