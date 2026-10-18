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

- [2026-10-18 12:1x +08] **编排者追加更正本票前提（★我写的票面有一处不实；来源＝写腿 `338-r1` 起手 AC#0 顶回＋我自己独立复跑坐实；⛔ 改上文一字）**：本票「现量」表里那句"差的形状＝**只有注释对齐空白**（例：…的尾注释列）"对**两枚**成立＝**不成立**。逐枚实测（blob 取 `git show HEAD:<f>` 落 `D:\tmp\wisp\20261018\orch\338-check`，尺＝`gofumpt v0.12.0 (go1.27.1)` ⇄ `gofmt`）：① `internal/tools/task_output_pointer_notice_errtext_174r7_test.go`＝**成立**（`gofumpt -w` 一发收敛、非空白尺 `tr -d ' \t'` 前后**空输出** rc=0）；② `cmd/wisp/panel_transport_35r2_test.go`＝**不成立**：`gofmt -d`＝**0 字节**（gofmt 侧干净＝脏只属 gofumpt 那一档），而 `gofumpt -d` 要动的是 **token 级**三形＝`type ( … )` 分组（吃掉 5 枚 `type` 词＋新增括号 token）／`[]string` 切片尾逗号／函数间补空行 ⇒ 非空白字节 **62106→62104**、`diff` **31 行**、我的硬尺 **rc=1**；★且**一发 `-w` 不收敛**（第 1 发后仍 `gofumpt -l` 命中 1、第 2 发才净）。⇒ **本票 AC#1 那句"只跑 `gofumpt -w`"与 AC#0 那把"非空白零变化"硬尺在这一枚上互相打脸**，腿按"未定义即停"停手上报＝做对了（⛔ 自行放宽我写的硬约束）。
  **裁＝按甲形放行 token 级改动，边界如下**（凭据与批准＝台账 `A1031`，撤销口令「撤 338 甲形」）：★**票 280（已 `-done`，10-09 10:3x 编排者裁语）逐字已把这一枚判给"下一枚真正动该文件的票"，且明写理由"只挑 gofmt 那 1 处改＝CI 照旧红、要改两遍"** ⇒ 本票就是那一枚"下一枚"，按 280 的原裁**一次改到 gofumpt 收敛**才是它的本意；⛔ 半改（只对齐不动分组）＝把 280 那句"要改两遍"再犯一次。硬尺**只对这枚让路**：`internal/tools/…174r7_test.go` 仍按"非空白零变化"交（AC#0 那把尺对它成立）。
  **让路的补偿（⛔ 空口放行）**：① 改动**只许来自 `gofumpt -w` 本身**（⛔ 手改、⛔ 顺手重排、⛔ `-w` 到别的文件）；② 收敛性自证＝逐发 `gofumpt -l` 直到空输出、每发落 `rc=`；③ ★**颜色必须成对**（AC#3）：`./cmd/wisp/` 改前⇄改后同台面四数＋红名册**逐名**作差＝**必须 0 枚新增**（token 级改动若动了语义，这里会露）；④ `go vet ./cmd/wisp/` 两档（本机 ⇄ `GOOS=linux`）＋`go build ./...` 各 `rc=`；⑤ 该枚交件侧 `gofmt -d` 仍须 **0 字节**（证明没把 gofmt 侧改坏）。
  ★**同时更正我起草票面时漏的一味**：`git grep -l panel_transport_35r2_test HEAD` 命中大量文件＝**全是散文引用**（各腿的门禁读数把这一枚当"既有债、非本腿"点名在册，例：`272-v2` 那句"票 298 AC#5 账"、`135` 票面 `:239`、`280` 票面 `:9-10`），**⛔ 任何仪器按该文件字节做判定**（腿的读者尺与我的复跑同结论：`scripts/**`／`tools/**`／`.github/**` 命中 0）。⇒ 它从 gofumpt 名册里消失是**本票的目的**、⛔ 读成"洗掉了别人的凭据"；但那几处散文引用今天起过期＝**归我具名登记**（⛔ 回改别人已入库的判语件／票面，只在我自己的账上追加）。
- [2026-10-18 12:5x +08] **写腿 `338-r1` 交件 bullet（只追加，⛔ 翻框／⛔ 改票名／⛔ 改 `Status:`）**：起手锚＝我自己 `git log -1` 现认 `d19288a8`（11:53，dev）。落笔序列＝`0b1d1f0b`（AC#0 尺＋logs）→ `87521cb4`（AC#1 tools 枚归位）→ `3d26240c`（AC#2 改前四道门）→ `93d0f39c`（AC#1b panel 枚归位＝你 12:1x 票面甲形裁语，凭据 `A1031`）→ `ee75c23f`（AC#2 交件侧门＋条件④）→ `fa5f008f`＋`b84dedd0`（AC#3 两包成对＋归位尺）→ `787dbfb9`（report.md 115 行）。★写面越界尺逐笔 0 行（`git show --name-only --format= | grep -v -E 三处写面`）；⛔ push。
  **判据读数（现跑，尺全在 `.scratch/wisp/probes/338/r1/logs/**`，共 34 枚）**：tools 枚非空白尺＝`RULER_RC=0`＋空输出、非空白字节 4966→4966、numstat 5/5、`i/lf w/lf`、CR 0 行、复量 0 行、二发幂等 YES；panel 枚按甲形＝两发 `-w` 收敛（round1 `rc_w=0`/l=1 → round2 `rc_w=0`/l=0）、交件侧 `gofmt -d`=0 字节、numstat 34/20、`i/lf w/lf`、CR 0 行、非空白尺记录 `rc=1`（62106→62107）。名册成对＝改前 6 行 ⇄ 交件 **4 行**、少的恰好本票两枚、⛔ 多动（`comm -13` 空输出）。四道门 rc 全 0；`GOOS=linux go vet ./cmd/wisp/`＝rc=1 逐字 `build constraints exclude all Go files in …sherpa-onnx-go-linux@v1.13.8`＝改前锚同配方复跑同 rc 同红句＝存量。颜色＝`./internal/tools/` 323/235/7/0 ⇄ 同；`./cmd/wisp/` 定向（我那枚文件 9 枚顶层测试）9/9/0/0 ⇄ 9/9/0/0；整包 207/150/2/0 ⇄ 203/146/2/0＝FAIL 名册两侧同两枚＝**本票新增红 0 枚**（整包名册作差 only_before=4＝驻留族整进程死亡的截断差，具名写在两枚末行＝`TestAC14AwaitedBindingReplyReachesThePage` ⇄ `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`，见 `logs/ac3-full2.txt`；⛔ 杀本机常驻 `wisp.exe` PID 21892）。
  **★AC#0 那处与票面不符你已裁并收下（甲形）；本腿再具名两枚：① `gofumpt -l internal cmd` 今天现量 7 行＝多出的 `internal/config/writeguard_handedit_227_test.go` 当时是 `??` 未跟踪（在 `87521cb4` 那笔前后被他腿提交），按跟踪集过滤才是票面那 6 枚＝分母尺具名；② 我自己两把尺的错已在 `logs/ac3-tgt2.txt`／`logs/run-ac3-tgt2.sh` 具名改正＝`-run` 串以 `/` 开头会让 Go 切出空首段＝匹配一切（v1 定向读数作废不当凭据），以及 v1 那对树差 3 枚（两锚间夹别腿提交）非单差异对。**
  **另四枚 CRLF 假脏（`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`／`cmd/wisp/models.go`）＝票面禁区逐字未动，仍脏、⛔ 本腿射程**（尺＝wt sha256 三点对读 `logs/ac1-ac2.txt` [2]⇄[8]⇄`logs/ac1b.txt` [7] 逐枚同哈希）。
  ★**四格一枚没翻**：本票判据格的框交件侧现量＝`grep -c "^- \[ \]"`＝**5** 枚（`:23`–`:27`＝AC#0～AC#4）、`grep -c "^- \[x\]"`＝**0** ⇒ 全未勾、我一枚未动（⛔ 翻框由非实现者裁）；另具名＝派单与我起手读票都写作"四格"，盘上现量是 **5 枚框**（AC#4 也在内）＝计数口径以盘上这把尺为准，我⛔ 改票面一字。
