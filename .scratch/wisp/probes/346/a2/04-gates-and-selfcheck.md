# 346-a2 · 门禁与纪律自证（本腿自己现跑的尺）

## 门禁

| 尺 | 对象 | 读数 |
|---|---|---|
| `go vet ./internal/winsec/`（本机） | 共享工作树 | rc=0，无输出 |
| `go vet ./internal/winsec/`（容器，mut-jia / mut-yi 两棵变异树） | 仓外树 | rc=0, rc=0（变异体编译面干净） |
| `gofumpt -l .scratch/wisp/probes/346/a2` | 本腿写面 | rc=0，输出 0 行 |
| `gofumpt -l internal/winsec`（正控＝尺在跑） | 禁写面 | rc=0，输出 0 行 |
| 本腿是否 build 产码 | — | 本腿零产码改动，vet 只作读数；无 build 面 |

`gofumpt` 全路径＝`D:\work\base\gopath\bin\gofumpt.exe`（PATH 上没有，首跑误报"不识别"已换全路径重跑）。

## 禁区自证（收尾终量）

- `git status --porcelain -- internal/winsec` ＝ **0 条**；`git rev-parse HEAD:internal/winsec/posix_seal_window_120_other_test.go`＝`0b0e5173…`、`…openleaf_322_other_test.go`＝`dfaddac8…`＝票面 AC#0 三锚同串值（blob 一字未动）。
- 真实仓库三句弃权红句逐字仍在（`Select-String` 行号现读）：`openleaf_322_other_test.go:206`、`posix_seal_window_120_other_test.go:225`、`:233`。
- 本腿写面＝`.scratch/wisp/probes/346/a2/*.md` 4 枚（00/01/02/03/04 中的前四枚＋本件）＋`docs/evidence/s1/346-a2-verdict.md`＋票 346 末尾一条 bullet；⛔ 翻 AC 框、⛔ 改 Status/票名/加 -done。
- 写面件自检（现跑）：逐枚非 0 字节（3994/3815/3857/3411/本件）、CR=0（全 LF）、码位 U+00A7＝0、U+00B7＝0（首版各 1 枚 B7 已自抓替换）、`.out` 扩展名 0 枚（probes/346 下现存的 13 枚 `.out` 全在 `a3/`＝别的腿的写面，⛔ 本腿触碰）。
- 临时件＝`D:\tmp\wisp\20261112\346-a2\`（export/mut-jia/mut-yi 三棵树、gocache、logs 全量 raw），只建不删；⛔ 在仓库目录内建 worktree（export 树用 `git clone --shared` 落 D 盘仓外）。
- 变异体落在 `D:\tmp` 仓外树、⛔ `-overlay`、每发前后哈希在 01/02 件；真实仓库零字节触碰。
- git：`add`/`commit` 两处显式手写 pathspec；⛔ amend/reset/rebase/stash/checkout ./clean；⛔ push。
