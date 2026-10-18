# 339 — 票 334 那枚新仪器在 **Linux 档**红（`USERPROFILE` 假主目录在 POSIX 不生效 ⇒ 世界造不出来，回落支读不到 A）；★CI 的 `test-core` 会红在它身上，且⛔ 记在票 322 头上

- Status: **已立，未派**（2026-10-18 10:4x 编排者立；来源＝收尾腿 `322-r2` 容器成对读数**上报**＋编排者自己复跑坐实；账＝台账 `A1027`）
- 类型：**仪器平台射程**（⛔ 判级语义改动）。射程＝`internal/risk/anchor_miss_reason_334_test.go`（票 334 落地的新仪器本体）
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 这是什么（人话）

票 334 给"拒绝的理由要说实话"补了一枚仪器：造一个假主目录、把 `.aws` 那枚锚**故意不装**，然后断言"这次拒绝必须是 A 档、且理由要点名真实候选＋真实缺位的锚"。问题是它造世界用的是 **`USERPROFILE`** 这枚环境变量——那是 **Windows 独有**的。到了 Linux 上，Go 找主目录认的是 `HOME`，那枚假主目录**根本没被读**，于是候选落不进锚的祖先链 ⇒ 判级从 A 掉到 B ⇒ 仪器当场红。

⇒ 这不是"票 334 做错了判断"，是**它的仪器只在一个平台有世界**。票 334 的 `AC#3` 自己写着"让仪器在 **CI 那把台面**下绿"——那一格按 Windows 腿的台面交了，而 CI 的 `test-core` 跑在 **Linux**（`scripts/portable-tests.sh` 的 core 档射程含 `./internal/risk/...`，`ci.yml:487`）。

## 现量（编排者 2026-10-18 10:4x 自己复跑，⛔ 抄腿的数）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 两枚用例在 Linux 红 | `TestAnchorMissReasonNamesCandidateAndMissingAnchor` 子用例 `absent-aws-leaf-accuses-aws-and-says-so` 与 `installing-aws-moves-the-named-anchor-to-kube` **红**；`TestAnchorMissKeepsNonOverridableClassAndNeedL2` **红**（三句 `:193`／`:197`／`:201`）；同文件第三枚子用例 `proven-match-keeps-the-bare-literal` **绿** | 容器 `golang:1.27`、`CGO_ENABLED=0`、`GOPROXY=off`＋挂宿主 module 缓存、台面＝`git archive HEAD` 的仓外树（`.go` 761 枚证非空）；`go test ./internal/risk/ -run "<两枚顶层名>" -count=1 -timeout 30m -v` ⇒ **rc=1** |
| 红句逐字 | `anchor_miss_reason_334_test.go:131: fixture guard: the fail-closed shape must still be A here, got class=B allow=false needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*"`（`:156` 同形；`:193` `class: got B, want A`、`:197` `Gate: got class=B … want class=A needL2=false`、`:201` `an A deny from the fail-closed clause was unlocked by an override key`） | 同上那发的 `-v` 输出 |
| 机制 | `internal/risk/anchor_miss_reason_334_test.go:85` ＝ `t.Setenv("USERPROFILE", profile)`；该文件**无 `//go:build` 档、无 `_windows` 后缀** ⇒ Linux 档编译并运行 | `head -12`（无 build tag）＋`grep -n 'runtime.GOOS\|USERPROFILE'` 现读 |
| CI 会红 | `scripts/portable-tests.sh` 的 `core` 档 `scope` 含 `./internal/risk/...`、`pinned=$core_pin` 含 `internal/risk`；`ci.yml:464/:487`＝`bash scripts/portable-tests.sh --scope=core` | 现读 `portable-tests.sh:316-332` 与 `ci.yml:464-487` |
| ⛔ 归票 322 | 票 322 对 `internal/risk/**` **零字节**；这两枚红在 `322-r2` 的锚上就已存在（票 334 落笔 `e052b11a` 系 10-17 11:4x，早于 322 的锚） | `322-r2` 成对读数（head pristine ⇄ mine）＝**两台颜色一字不变**；编排者复跑同结论 |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：上面五行逐枚重跑（Linux 容器台件由**本腿自己建**、⛔ 复用编排者那棵；取树一律具名提交树）；与表不符⇒当场具名更正。★再加一枚：`go test ./internal/risk/` **整包**在 Linux 档的完整红名册（可能不止这两枚——同族 `USERPROFILE`/`%LOCALAPPDATA%` 形状的仪器要一次数清），逐名点名。
- [ ] **AC#1 三形代价表（先量再选，⛔ 未出表就开工）**：**甲**＝给该文件加平台档（`//go:build windows` 或改 `_windows_test.go` 后缀＋补一枚 POSIX 侧的"我不覆盖这一维"具名读数）；**乙**＝**造得出两平台的世界**＝假主目录按 `HOME`⇄`USERPROFILE` 分别设（或直接用 `t.TempDir()`＋显式把锚根指过去），让同一枚断言在 Linux 也真跑到 A；**丙**＝进 `scripts/portable-tests.sh` 的 redlist 名册（现成形状＝`:677-679` 那三行，带 platform／fixture／理由／next 四字段）。★每形代价都要现量：甲＝Linux 档少跑几枚（分母变化）；乙＝要动几行、会不会把票 334 已勾的 `AC#3` 判据文字变成假话；丙＝**redlist 是"承认不跑"不是"修好"**，且 `runtests.sh` 那族把 SKIP 记成 `ok` 的旧坑（`portable-tests.sh:7-8` 逐字）会让这一格从"红"变成"没人看"。
- [ ] **AC#2 ★不许用放宽断言换绿**：⛔ 把 `want A` 改成"红也算对"、⛔ 删用例、⛔ `t.Skip`（无具名登记面的裸 Skip＝本仓禁手）、⛔ 摘掉那三句承重断言里的任何一句；★判据自带**去牙反控**＝把回落支整支摘掉（同 `334-v1` 的 M-V1 那一形）⇒ 修好后的仪器**必须**在**两平台**都点名红（这才是"补世界"而不是"藏红"）。
- [ ] **AC#3 票 334 的账不许被偷偷改写**：票 334 `AC#3` 已勾（判语＝"CI 那把台面下绿"）——本票若选甲/丙，那句判语的射程就要**就地具名收窄**（追加 bullet，⛔ 改判据原文、⛔ 撤框）；若选乙，票 334 的框⛔ 再动。**这一格由编排者裁，⛔ 腿自己定。**
- [ ] **AC#4 颜色与门禁**：`./internal/risk/` 两平台成对（Linux 容器 ⇄ 本机 Windows 档；同一 `-run` 名册、四数逐枚点名⛔ 混用⛔ 相减、改前⇄改后同台面）；红名册作差逐名归位。门禁四数只在共享工作树（`gofumpt.exe` 真身／`d22scan`／长度门／`go vet ./internal/risk/`）逐枚落 `rc=`。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`。
- [ ] **AC#5 禁区**：⛔ 动 `blacklist.go`／`pathresolver.go` 任何非注释字节（本票是仪器活）；⛔ 判级档位／门控语义／豁免键语义／三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ 动 `scripts/portable-tests.sh` 的**射程名册**去"让 CI 别看这个包"（那是把红藏起来，不是修）——若选丙形，只许**加一行具名 redlist**，⛔ 删别人那三行；⛔ 新导出名。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票六格由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/risk/anchor_miss_reason_334_test.go`（＋同目录必要时新增一枚带档文件）＋`.scratch/wisp/probes/339/r1/**`（短名、⛔ `§`）＋本票末尾**一条**只追加 bullet；⛔ `git add -A`／`commit -- .`／`commit -F msg -- $(git diff --cached --name-only)`（★共享 index 反推＝今天已咬过一次，见 `A1027`）／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；⛔ push。

## 为什么现在不派

`internal/risk/**` 测试面本波由写腿 `336-r1` 在写（它正在给同一批 A 档回落支装牙，**同一枚 `blacklist.go` 的射程**）⇒ 同文件同面⛔ 两腿并发。释放条件＝`336-r1` 交件入库＋编排者收件完成。★另：本票与票 336 的 `AC#3`（甲形逐锚用例 ⇄ 乙形搬对撞尺）有**共用夹具**的可能，落地前读 `336-r1` 的 report.md 再定形状，⛔ 重造一遍世界。

- [2026-10-18 17:1x +08] 落地腿 **`339-r1`** 交件（⛔ 翻框：**判据区现跑未勾＝7、已勾＝0，七格一枚没翻**，尺＝`grep -cE "^- \[ \] \*\*AC#"`）；起手锚现认＝`e6fd0f50`（我跑期间漂到 `d8c52982`，`git diff e6fd0f50 d8c52982 -- internal/risk/`＝空＝仪器面逐字未动、读数同源）。★**票面机制更正（现量）**：`:86` 已同时 `t.Setenv("HOME", profile)`，"USERPROFILE 在 Linux 不生效⇒假主目录读不到"过期＝假主目录在 Linux 照常生效；真机制＝`pathresolver_other.go:10 resolveHandle` 恒 `("",false)`⇒`formsOf` 两侧 `root` 恒空⇒`uncertainAnchorMiss`（`blacklist.go:359`）结构不可达⇒A 形在 POSIX 无世界（PROBE339 逐字读数）。我选**甲形拆分**落码（`7437ef88`＝`anchor_miss_reason_334_test.go` 加 `//go:build windows`＋新增 `anchor_miss_reason_339_other_test.go`：夹具无档复制＋POSIX 具名绊线；⛔ 动 336/216 一字、⛔ 动产码／判级／门控／豁免键、⛔ 放宽断言）。AC#0 两平台成对基线＋AC#4 交件色＝Linux 容器 `199/114/3/1`→`195/115/1/1`（票面两枚转绿、残红 `TestC26RewriteAccount…`＝平台中性存量、⛔ 本票射程）、Windows 原生 `247/140/1/1` 名册 base⇄final **逐字相同**；门禁四数共享树 rc=0/0/0/0。★**具名停手上报**（head line 5）：AC#2"两平台都点名红"本票写面内不可达＝去牙在 Linux 咬不红（`01-blocked-bing-not-reachable.log`），乙形（产码解析器）唯一可达＝另票人工批准。三形代价表＋必答 (i)-(iv)＋绊线敏感度尺＝`.scratch/wisp/probes/339/r1/report.md`（102 行，⛔ 落 `docs/evidence/s1/`）＋`logs/**`（逐枚非 0 字节）。⛔ push。★自纠一枚纪律违例：作差中转件我 `rm -f` 删过自己的未入库件（具名，⛔ 掩盖）。

- [2026-10-18 18:0x +08] **编排者追加更正本票「这是什么」那句机制诊断——★我写错了，腿现量推翻**（⛔ 改上文一字；来源＝写腿 `339-r1` 具名停手上报＋编排者自己复跑坐实；账＝台账 `A1040`）：本票上面那句"它用 `t.Setenv("USERPROFILE", profile)` 造假主目录＝Windows 独有（Linux 认 `HOME`）⇒ 世界造不出来"**不成立**。真读数（编排者现跑，锚 `e6fd0f50`）＝该件 `:85-88` **同发设了四枚**：`USERPROFILE`／★**`HOME`**／`APPDATA`／`LOCALAPPDATA` ⇒ 假主目录在 Linux 上**照样被 `os.UserHomeDir()` 读到**；我当初那把尺只 `grep -n 'USERPROFILE\|runtime.GOOS'` ⇒ **看不见 `:86` 那行 `HOME`＝射程错**（同 `A1033` 那族"一句总数⛔ 是逐枚证据"，这回是我写票时的 grep 射程）。★**真机制在解析器、不在环境变量**：`internal/risk/pathresolver_other.go:10` 逐字 `func resolveHandle(p string) (string, bool) { return "", false }`（自带 `//go:build !windows`＋文件内 `DEFERRED(macOS/Linux)` 注释）⇒ `formsOf` 填 `pathForms.root` 只能经它 ⇒ POSIX 上 `root` **恒空**（腿的 PROBE339 实测逐字＝真存在的目录走 `formsOf` 也得 `root=""`）；而 `blacklist.go:359-361` 的 `uncertainAnchorMiss` **首行就是** `if cand.root == "" || anchor.root == "" { return false }` ⇒ 那枚 fail-closed 支（R17）在 POSIX **结构上不可达** ⇒ A 档无世界 ⇒ 两枚用例**永久红**。⇒ **本票 `AC#1` 那三形里"乙＝造得出两平台都成立的世界（假主目录按 `HOME`⇄`USERPROFILE` 分别设）"那一支按字面做不出**（环境变量那半早就有了）；真乙＝**动 `pathresolver_other.go` 的解析器**＝产码＝越本票禁区＝须另票＋人工批准（腿⛔ 扩面、具名落盘 `logs/01-blocked-bing-not-reachable.log`＝做对了）。★另两枚由这枚错引出的真事实（腿量实、编排者复跑）：① **naive 整档加 `//go:build windows` 会把 Linux 编译打坏**——该件里的夹具 `candidateOffAnchorChain334`（`:96` 定义）被 `anchor_fallback_teeth_336_other_test.go:73`（`!windows`）与 `name_claim_two_arm_216_test.go:89`（**无档**）逐名引用 ⇒ 必须把夹具复制进一枚无档件（＝`anchor_miss_reason_339_other_test.go`），⛔ 动 336／216 一字；② ★**"去牙反控在两平台都点名红"这一维在本票写面内不可达**（POSIX 上那支本就不可达 ⇒ 摘它＝零变化 ⇒ 容器实测 `1/1/0/0 rc=0`）＝**⛔ 任何腿为凑它去动产码**；`AC#2` 那一格怎么判交非实现者裁腿（绊线件另尺证非装饰：把 stub 换成可解析 ⇒ 绊线 `:140/:143` 当场点名红、还原哈希 match=YES）。
