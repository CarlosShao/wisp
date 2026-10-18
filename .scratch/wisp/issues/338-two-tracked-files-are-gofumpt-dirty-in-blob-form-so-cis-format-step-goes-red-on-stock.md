# 338 — 两枚跟踪件在 blob 态就被 `gofumpt` 判真脏（注释对齐形）＝CI 那两步从推上去那天起会红

- Status: **已立，未派**（2026-10-18 10:3x 编排者立；来源＝新任编排者接手自检现跑（`A1024` 第③节）；账＝本条）
- 类型：**格式归位**（⛔ 语义改动、⛔ 判据改动）。射程＝`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`＋`cmd/wisp/panel_transport_35r2_test.go`
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 这是什么（人话）

仓库有一道"格式门"（`gofumpt`），CI 里跑两次（`ci.yml:172` 与 `:189`，分母＝跟踪集）。今天有**两枚文件即使从提交里取出来也是脏的**——不是谁工作树里被某个程序改成了 CRLF 那种假脏，是**真脏**：注释没对齐。⇒ 只要往 `origin/dev` 推，那两步就会红，而红的是**存量**、不是任何一条腿新造的。本票就是把这两枚归位，让门说实话。

## 现量（编排者 2026-10-18 09:4x 现跑；引用前先重跑）

| 事实 | 读数 | 尺 |
|---|---|---|
| 全量脏名册＝6 枚 | `gofumpt -l internal cmd` 现量 **6 行** | `/d/work/base/gopath/bin/gofumpt.exe -l internal cmd` |
| 其中 4 枚＝**CRLF 假脏** | `internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`cmd/wisp/models.go`：`git ls-files --eol`＝`i/lf w/crlf`，把 **HEAD blob 抽到仓外 LF 文件**复跑 gofumpt＝**CLEAN** | `git show HEAD:<f>` → `D:\tmp\wisp\20261018\orch\blobcheck\<f>` → `gofumpt -l` |
| ★其中 2 枚＝**blob 态真脏** | `internal/tools/task_output_pointer_notice_errtext_174r7_test.go`（末笔 `67a7ab41` 10-14 21:0x）、`cmd/wisp/panel_transport_35r2_test.go`（末笔 `18d438e` 10-10 16:1x）：blob 抽出来仍 **DIRTY** | 同上那一把尺（只换 rev→blob） |
| 差的形状 | `gofmt -d` 现读＝**只有注释列的对齐空白**（例：`canonicalizeFailsHere,` 后面那串空格与下一行 `"tools:", "risk:",` 的尾注释列）；⛔ 任何 token／断言／字符串字面量变化 | `gofmt -d <file>` |
| CI 会不会红 | **会**——`ci.yml:172` 那道 `gofumpt -l . tools/d22scan tools/mockllm` 非空即 `exit 1`；`:189` 那道同尺同分母 | 读 `ci.yml:172-208` |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：把上面四行逐枚重跑（blob 一律 `git show HEAD:<path>` 取到**仓外** LF 文件再量，⛔ 在工作树量 CRLF 形）；与表不符⇒当场具名更正。★再加一枚：`gofumpt -d` 的**逐字 diff** 进件，并证它**只动空白**（尺＝`diff <(tr -d ' \t' <before>) <(tr -d ' \t' <after>)` 必须空输出＝非空白字节零变化）。
- [ ] **AC#1 归位两枚**：只跑 `/d/work/base/gopath/bin/gofumpt.exe -w <那两枚>`（或等价手工），⛔ 动其它 4 枚 CRLF 件、⛔ 动任何断言／用例名／字符串字面量／注释文字内容；新字节必须 **LF**（⛔ 让 `core.autocrlf` 把归位变成换行符改写；尺＝`git ls-files --eol` 交件侧仍 `i/lf`）。
- [ ] **AC#2 门变绿要成对**：改前⇄改后各跑一发 `gofumpt -l internal cmd`（同一台面），名册**逐名作差**＝必须恰好少了这两枚、⛔ 只报枚数；`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`／`go vet ./internal/tools/`／`go vet ./cmd/wisp/` 逐枚落 `rc=`。
- [ ] **AC#3 颜色**：`./internal/tools/` 与 `./cmd/wisp/` 两包**仓外导出树**成对（改前⇄改后、同一台面、`-count=1 -timeout 30m`、四数逐枚点名⛔ 混用⛔ 相减）；红名册作差逐名归位（`cmd/wisp` 整包有已知真窗区整进程死亡＝具名归位、⛔ 当本票新增）。★**`cmd/wisp/**`＝共享编译面**：派单时该面本波零写腿，落地那波仍须**每波一枚写腿**＋`go test ./cmd/wisp` 同一时刻只跑一个。
- [ ] **AC#4 纪律格**：⛔ 翻框（本票四格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝那两枚跟踪件＋`.scratch/wisp/probes/338/r1/**`＋本票末尾一条只追加 bullet；⛔ `git add -A`／`commit -- .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ push。

## 禁区

⛔ 顺手把 4 枚 CRLF 假脏"一起修了"（那是另一件事：谁把工作树写成 CRLF 的机制未归到底，见 `A1024` 第③节ⓑ；动它们＝把未归到底的东西混进一笔格式提交）；⛔ 为让门绿去改 `ci.yml` 或放宽任何断言；⛔ 动 `thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`go.mod`／`go.sum`／`.gitignore`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ `frontend/**`／`design/**`。

**撤销口令**：「撤 338」（⇒ 两枚退回本票前形状，本票退回未派）。

## 为什么现在不派

本波六条腿在飞（3 写：`336-r1`／`166-r1`／`322-r2`；总代理帽 ≤8）。本票动 `cmd/wisp/**`＝共享编译面，且第二波已排 `329-r1`∥`330-r1`（同面单车道串行）⇒ **排在第二波之后**，或等 `cmd/wisp` 那面空出来再派。⛔ 与任何碰 `cmd/wisp` 的腿同波。
