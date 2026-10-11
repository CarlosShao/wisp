# 308-r1 · 40-tags-and-gates（`AC#3` 档位与分母／`AC#4` 门禁与越界／`AC#5` 对账）

## `AC#3` 分母与档位：两枚崩点件在 CI 哪一档真跑

| 件 | 头 1 行 tag | test-core (ubuntu) | test-windows (windows-latest) | 本腿能否现量 CI 色 |
|---|---|---|---|---|
| `pathresolver_anchor_spelling_windows_test.go` | `//go:build windows`（尺＝`head -1`） | 不编译、不⛔ 求值（GOOS=linux ⇒ 整件被排除） | 编译＋求值（GOOS=windows ⇒ 整件纳入；`scope=windows` 名册含 `./internal/risk/`） | 否（不 push） |
| `pathresolver_junction_windows_test.go` | `//go:build windows` | 不编译、不⛔ 求值 | 编译＋求值（`TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied` 都进名册） | 否（不 push） |
| `pathresolver_volumeless_spelling_308_test.go`（本腿新增） | 无 build tag（尺＝`head -1`＝`package risk`） | 编译＋求值（helper 的牙证明不依赖 windows 门；`TestVolumeInvariantSpellingsVolumeLess308` 跨 OS） | 编译＋求值（同上） | 否（不 push） |

档位尺（本腿读过的）：
- job 落点＝`test-core` `.github/workflows/ci.yml:419-420` (`runs-on: ubuntu-latest`)、`test-windows` `.github/workflows/ci.yml:533-534` (`runs-on: windows-latest`)。
- 名册落点＝`scripts/portable-tests.sh`：core scope 含 `./internal/risk/...`（行 240，`scope=(… ./internal/risk/... …)` 那段），windows scope 含 `./internal/risk/`（行 251，`windows) scope=(… ./internal/risk/ …)`）。
- `grep -rln 'winlive' internal/risk/*_test.go`＝**0 命中**（⛔ 两枚崩点件带 `&& winlive` 的⛔ 更窄档；⛔ 需拿"winlive 档"⛔ 顶替 test-windows）。
- 另一条 windows 档（`ci.yml:813`）＝`bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`——**只打 Case 1**（`TestPathResolverJunctionWindows`），⛔ 覆盖两枚崩点用例（尺＝`grep -n 'TestPathResolverJunctionWindows\|TestClassifySpellingInvariance\|TestPathResolverUNCAListDenied' …` 在那行只一枚）。⇒ 那步⛔ 是本票的兜底档；本票两枚崩点色的兜底档＝`portable-tests.sh --scope=windows`。

★ 具名欠账（⛔ 由落地腿自勾，交非实现者终裁 `308-v1`／编排者推送那批一起取）：
- **Q-308-r1-A**：test-windows 档在推送后 `portable-tests.sh --scope=windows` 对 `TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied`／`TestVolumeInvariantSpellingsVolumeLess308` 三枚的**实跑 PASS/FAIL 色**——本腿⛔ push、⛔ 跑 CI。本腿只交"哪一档会打到它们"的静态定位（上表＋上段尺），⛔ 用"本机跑过"顶替 CI 色。
- **Q-308-r1-B**：test-core (ubuntu) 档对本腿**新件**（⛔ build tag）的实跑色——helper 的牙（`TestVolumeInvariantSpellingsVolumeLess308`）⛔ ubuntu 也会求值，那发 CI 色⛔ 本腿现量。

先例对齐（票 §规矩）：票 301 包级 guard D（`internal/audio/` 挪进 test-windows 那次，同形"windows 档⛔ ubuntu 求值"）、票 111 那 12 枚"从未求值"（`grep -n 'portable-tests|winlive' scripts/portable-tests.sh` 那批条幅）。本票⛔ 搬票 300/301/302 那三档 windows 档搬运活（票 §边界 ⛔ 条），只交自⛔ 两枚档位读数＋新件档。

---

## `AC#4` 门禁与越界

三把门禁各⛔ 一件一行 `rc=N`（尺＝件本身＋每⛔ 附一份 `.rc.txt`）：

| 门禁件 | rc | 件 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | `logs/gate-build.txt`（空文件，⛔ 命中）、`logs/gate-build.rc.txt`＝`rc=0` |
| `go vet ./internal/risk/` | **0** | `logs/gate-vet.txt`（空）、`logs/gate-vet.rc.txt`＝`rc=0` |
| `sh scripts/d22scan.sh` | **0** | `logs/gate-d22scan.txt` 末行"clean - no D22 ban violations"、`logs/gate-d22scan.rc.txt`＝`rc=0`；⛔ 一枚 emoji 命中（ban #8 scope 内 `internal/`=531, `cmd/`=119 全⛔ 中） |

⚠ 309-r1 在飞面：本发 `go build ./...` 全绿、⛔ 命中 `internal/tools/**` 红句。若下轮命中，本腿按硬约束"具名记那发红句、稍后重试一发"处理；本发⛔ 需要。

★ 格式名册两把并排、各写射程（既有债，⛔ 我修、⛔ 我因它判红；⛔ 与票 §"⛔ 动 SLO 阈值／golden／`thresholds.go`／`allowlist.txt`／D43／C1–C32／D1–D47／三枚冻结件" ⛔ 相⛔ ）：

| 尺 | 射程 | 命中数 | 逐枚 | 性质 |
|---|---|---|---|---|
| `gofmt -l internal/risk` | 工作树窄尺 | **1** | `internal/risk/provenance.go` | **CRLF 假枚**（`git ls-files --eol`＝`i/lf w/crlf attr/text eol=lf`；blob 尺 `git cat-file blob ":…provenance.go" \| gofmt -l /dev/stdin`⛔ 命中＝blob 侧⛔ 脏）。⛔ 我引入，⛔ 我修 |
| `git ls-files '*.go' \| xargs gofmt -l` | 全仓 tracked 宽尺 | **27** | `.scratch/**` 23 枚（变异拷贝，故意脏）＋非 `.scratch` 4 枚（`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`） | 同 23 枚⛔ 因（变异拷贝⛔ 干净⛔ 是它的职责）；4 枚非-`.scratch`⛔ 同形 CRLF 假枚（⚠ 编排者裁语说"5枚非-`.scratch`"，本腿现量＝4 枚，⛔ 差一枚＝`internal/tools/bridge.go`；那枚⛔ 我在飞面上，属 309-r1；本腿⛔ 追、⛔ 修，⛔ 因⛔ 判红） |
| `gofmt -l` ⛔ 我三枚（`pathresolver_anchor_spelling…`／`pathresolver_junction…`／`pathresolver_volumeless_spelling_308…`） | 三枚⛔ 面窄尺 | **0** | — | 本⛔ 件⛔ 干净（尺＝`gofmt -l <三枚文件>` 空输出、`rc=0`；件＝`logs/gofmt-my-three.txt`） |

`git show --name-only --format=` **逐笔**量越界＝见 `90-final.md` §1；越界枚数应＝0（写面只有 `internal/risk/*_test.go` × 3 ＋ `probes/308/r1/**` ＋票 308 末尾一枚 bullet）。

---

## `AC#5` 对账（⛔ 翻票 72 与票 307 任何一格）

- **⛔ 票 72 任何⛔**：本票⛔ 已改那两枚件，但**改的是不同件事**——票 72 管"期望侧写死了 temp 拼法"（甲族 #1/#4/#5/#6/#7 那批），本票管"崩点让整包零求值"（`vol[:1]` 无守卫那两行）。尺＝`git log --oneline -- internal/risk/pathresolver_junction_windows_test.go` 上 `f278c57f`（票 72 甲族 #5/#6/#7，+36 −7）与本腿 commit 是两个独立事件；本腿⛔ 回退⛔ 翻 `f278c57f` 任何一行。
- **⛔ 扩票 307**：票 307 已 `-done`（`.scratch/wisp/issues/307-*-done.md`），残余归 `A839` §3。本票⛔ 借道⛔ 一枚 307 的框。尺＝`ls .scratch/wisp/issues/307*` 只一枚 `-done`。
- **⛔ 把票 300／301／302 三档 windows 档搬运活并进本票**：票 §边界逐字禁；本腿只交自⛔ 两枚档位读数（上表）＋新件档，⛔ 动票 300/301/302 任何一枚。
- **⛔ 动产码**：`internal/risk/pathresolver.go` 与 `internal/risk/pathresolver_windows.go` 逐字节⛔ 变（尺＝`git diff --numstat -- internal/risk/pathresolver.go internal/risk/pathresolver_windows.go`＝0 行；写面声明见 `00-anchor.md`）。
- **⛔ 动票面正文**：本腿只在票 308 末尾追加**一枚** bullet（`git diff --numstat` 应＝1 加 0 删，尺＝commit 后 `git show --numstat --format=`）。
