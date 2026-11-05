# 票 85（85a-2 腿）：把 `gofumpt` 钉成具体版本 —— 只动 `ci.yml` 一处安装行

- 腿：`85-pin`（**实现腿**；裁决者≠实现者，本件**不是**裁决表，只是本腿读数）
- 时间：2026-11-05 11:1x–11:4x +08（`date` 起手 `2026-11-05 11:14` 前后）
- 基线：HEAD `6a18aa8dd41769316471fc71981858b89b49c4f9`
- 写面（硬闸）：`.github/workflows/ci.yml`（**只动 gofumpt 的安装行 + 其注释**）· `.scratch/wisp/issues/85-*.md`（**本腿未动，见 §6-STOP-5**）· `docs/evidence/s1/85-*.md`（本件）
- 本机：go1.27.1 windows/amd64；`/d/work/base/gopath/bin/gofumpt.exe` = `v0.12.0 (go1.27.1)`

---

## 0. 结论速览

| 格 | 判定 | 判据在哪 |
|---|---|---|
| `gofumpt` 从"没钉版本"变成"钉 `v0.12.0`" | **PASS** | §2/§3：`@latest` 命中 2→**0**、无版本安装行 1→**0** |
| 钉的版本对不对 | **PASS（且今日零判定变化）** | §4-B：`go list -m mvdan.cc/gofumpt` = `v0.12.0` = 现在 `@latest` 解析出的那枚；本机二进制同版 |
| 钉的字符串真的能装、能跑 | **PASS** | §4-C：`go install ...@v0.12.0` 到临时 `GOBIN` rc=0，`--version` = `v0.12.0 (go1.27.1)`；该二进制在纯净快照上 `gofumpt -l` **空** |
| YAML 还能解析 / step 体还能跑 | **PASS** | §4-A：`yaml.safe_load` OK（6 jobs / lint 14 步）；抽出的 step body `bash -n` rc=0 |
| AC#1 的"真 run 日志行"那一半 | **STOP（需 push + 托管 run，属编排者）** | §6-STOP-1 |
| AC#2 / AC#3（34/78 条 findings 清算 + 变异回验） | **STOP（本腿明令不做，另立 85b/122）** | §6-STOP-2 |
| AC#4（`PathResolver junction placeholder` 步骤改名） | **STOP（不在本腿写面；且要真 run 结论）** | §6-STOP-3 |
| `ci.yml` 里 staticcheck | 不动（已是 `@2026.2.2`） | §1-B / §4-A |

---

## 1. THREE-CORE 进场核对（现跑，不抄票面）

### A. 票面引用的字符串在不在

- 票 `85-…:15` 写「`ci.yml:72` 写的是 `go install mvdan.cc/gofumpt@latest`——没钉版本」。
  **字符串在**（本腿实测同一行**多了 `@latest`**），**行号早漂**：落刀前在 **`:275`**（落刀后在 **`:293`**）。
  行号从票面 `:72` 起漂 **+203**，再被本腿 **+18**。⇒ 引用请用**步骤名**「`gofmt (gofumpt)`」。
- 票面裁定「禁止 `@latest`」（`:41`）：**本腿落刀前该文件里 `@latest` 命中 2 行**（`101` 注释、`275` 代码），落刀后 **0 行**。

### B. `file:line` 与版本号现状（票面腐烂两处）

1. 票面 `:41`/`:82`/`:166` 通篇写 staticcheck 钉 **`2026.2.1`**；**盘上是 `@2026.2.2`**（本腿实测 `ci.yml:451`，
   连同 `:383-` 的那段注释一起解释 V5 导出格式）。它是**腿 `ci-r1`、commit `702c13f51`（2026-10-26）**改的，
   票面没回填 ⇒ 这是票面腐烂，不是本腿改的。**本腿未动 staticcheck 一个字。**
2. 票面 `:15` 引 `ci.yml:72`、以及 `:31` 引 `ci.yml:136-142`：两处行号均过期（后者票面 `:113` 自己已承认、
   `:193-195` 已改为"此后一律用步骤名引用"）⇒ 本件从其约定，不钉行号。

### C. 别腿有没有动过（写面冲突检查）

- 落刀前 `git status --porcelain -- '*.go'`：**跟踪文件 0 行修改**（只有 `.scratch/**` 未跟踪件）。
- `ci.yml` 最近三笔：`34a630d72`（2026-11-03，票 114）· `8d6e150f8`（2026-10-26，CER）· `702c13f51`（2026-10-26，ci-r1）。**无他腿在飞。**
- ⚠ 共享树里**有他腿的未跟踪遗留**，它们会污染**整树口径**的门（见 §5 两条脚注）：
  `build/mut_rvw.go`（91,003 B，mtime 11-04 10:26；`git ls-files build | wc -l` = **0**）、
  `internal/tools/ticket23_image_channel_test.go`（`??`）。**本腿一律不碰、不删**（临时件只建不删）。

---

## 2. RED（落刀前，`HEAD 6a18aa8dd`；可逐字重跑）

```
$ git show HEAD:.github/workflows/ci.yml | grep -c '@latest'
2
$ git show HEAD:.github/workflows/ci.yml | grep -E '^[[:space:]]+go install ' | grep -vc '@'
1
$ grep -nE '^[[:space:]]+go install ' <(git show HEAD:.github/workflows/ci.yml)
275:          go install mvdan.cc/gofumpt@latest      <- 无版本（@latest）
433:          go install honnef.co/go/tools/cmd/staticcheck@2026.2.2
```

⇒ 两把工具里**只有 gofumpt 还悬空**；装上它的是 `:275`，而**跑它的有两处**：
`:294` 的 `gofumpt -l internal cmd tools scripts frontend`，以及
`gofmt (gofumpt) - the tracked set is the denominator (ticket 161 AC#7 form A)` 那一步所调
`.scratch/wisp/probes/161/r5/attrib.sh`（其 `:152-154` 从 `$GOPATH_BIN/bin/gofumpt` 取二进制，
即**本步 `go install` 装出来的同一枚**）。⇒ **一处安装、两处判据**，所以这一行是这把门的唯一版本源。

---

## 3. 落刀（`ci.yml`，`git diff --numstat` = **20 增 / 2 删**）

1. `:293` `go install mvdan.cc/gofumpt@v0.12.0`（原 `@latest`）。
2. `:101` 的注释里那枚安装串同步成钉版（否则文件里留一句假话，也留一个 `@latest` 命中）。
3. `- name: gofmt (gofumpt)` 名下加**一段 18 行注释**说明：为什么钉、为什么是 `v0.12.0`、
   装过一次的读数、以及"提升这个 pin = 一次带读数的代码改动，不是一次上游事件"。
   **零行为改动**：`if:` 未动、`run:` 的判断逻辑未动、没有新增 `if:`、没有 `continue-on-error`、没有放宽任何断言。

---

## 4. GREEN（落刀后，现跑）

### A. 形状仍是合法的

```
$ python -c 'import yaml; d=yaml.safe_load(open(".github/workflows/ci.yml",encoding="utf-8")); ...'
yaml.ok jobs= ['lint', 'lint-frontend', 'slo-full', 'slo-smoke', 'test-core', 'test-windows']   lint steps= 14
$ # 用 PyYAML 把 `gofmt (gofumpt)` 的 run 体抽出来（174 B）后：
$ bash -n /c/tmp/85-step-body.sh ; echo rc=$?
rc=0
```

抽出的 run 体逐字（7 行）：

```
go install mvdan.cc/gofumpt@v0.12.0
OUT="$(gofumpt -l internal cmd tools scripts frontend)"
if [ -n "$OUT" ]; then
  echo "gofumpt would reformat:"
  echo "$OUT"
  exit 1
fi
```

### B. 钉的版本 = 今日 `@latest` 解析出的版本（⇒ 今日判定变化 **0**）

```
$ go list -m mvdan.cc/gofumpt        # 未钉版时它给什么
mvdan.cc/gofumpt v0.12.0
$ go list -m mvdan.cc/gofumpt@v0.12.0
mvdan.cc/gofumpt v0.12.0
$ /d/work/base/gopath/bin/gofumpt.exe --version
v0.12.0 (go1.27.1)          # 与本仓 go.mod 的 toolchain go1.27.1 同线
$ go list -m -versions mvdan.cc/gofumpt | tail -1
… v0.10.0 v0.11.0 v0.12.0   # v0.12.0 是当前最新；钉它不是"降级"，是把今天的行为写死
```

### C. 钉的字符串本身可装、可跑（不是只改字）

```
$ GOBIN=/tmp/85-gofumpt-bin go install mvdan.cc/gofumpt@v0.12.0 ; echo rc=$?
rc=0
$ /tmp/85-gofumpt-bin/gofumpt --version
v0.12.0 (go1.27.1)
$ cd /c/tmp/85-clean   # git archive HEAD 抽出的 internal cmd tools scripts frontend，跟踪 .go = 1112 枚
$ /tmp/85-gofumpt-bin/gofumpt -l internal cmd tools scripts frontend
（空）  rc=0
```

### D. 计数（RED→GREEN 的机读判据）

| 读数 | 落刀前 | 落刀后 |
|---|---|---|
| `grep -c '@latest' .github/workflows/ci.yml` | **2** | **0** |
| `go install` 行总数 / 其中**无** `@` 的 | 2 / **1** | 2 / **0** |
| `grep -n '@staticcheck\|staticcheck@'` | `@2026.2.2` | `@2026.2.2`（未动） |

---

## 5. 门禁读数（本腿现跑；两处共享树脚注**不是**本腿造成的）

| 门 | 读数 | 判定 |
|---|---|---|
| `go build ./...` | **rc=1** —— 红在**未跟踪**的 `build/mut_rvw.go`（`undefined: voiceTtsFiles` 等 11 条）。`git ls-files build \| wc -l` = **0** ⇒ 该目录不跟踪、CI 上不存在，是**他腿遗留的变异探针** | 见脚注① |
| `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（11.9 s） | **PASS** |
| `go test -count=1 -timeout 200s ./internal/... ./tools/...` | **rc=1** —— **29 包 `ok` / 3 包 `FAIL` / 4 条用例红**：`internal/tools::Test229SpawnNamesAreRosterShaped`、`internal/panel::TestComposerContractTypesMatchFrontend` ＋ `TestC21DesignTokensFourWayAgree`、`internal/risk::TestC26RewriteAccountIsConsumedAtEverySecurityLeg`。四条**都在 `git archive HEAD` 纯净快照上逐条复现**（见脚注②） | 见脚注② |
| `gofumpt -l internal cmd tools scripts frontend`（**纯净快照**，1112 枚跟踪 `.go`） | **空**，rc=0 | **PASS** |
| 同上（**共享工作树**） | 1 行：`internal/agent/approval/pending_read.go` —— 该文件工作树是 **CRLF**（`file` 逐字 `with CRLF line terminators`）、`git status` **0 行**、`git show HEAD:` 的同一文件喂 gofumpt **干净**（rc=0）⇒ **换行符工作树假阳性，不是格式问题** | 见脚注③ |
| `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；runtests `PASS=35 FAIL=0 SKIP=0`，`=== RUN=77`；八格 scope 分母 **落刀前 `356/56/116/48/39/116/876/203`、提交后复跑 `357/56/116/48/39/116/878/203`**（⚠ 分母是**工作树**读数，随别腿新增未跟踪件漂；两次判定都 `clean`） | **PASS** |
| `go vet ./...` | **未跑**（会撞脚注①同一枚未跟踪文件）；本腿零 `.go` 改动，无被改包可点名 | 未跑·如实记 |

**脚注①（`go build ./...`）**：红点是 `build/mut_rvw.go`，91,003 B、mtime `11-04 10:26`、**未跟踪**（`git ls-files build` = 0）。
票 230-ac3 的件里已把同一枚文件记成"别腿遗留"（`docs/evidence/s1/230-ac3-handoff-disposition.md`）。**不删**（临时件只建不删）。

**脚注②（4 条基线红）**：`./internal/... ./tools/...` 全跑 = 29 `ok` / 3 `FAIL` / 4 条用例红。四条**都不是本腿的**，两把独立尺：

1. **纯净快照复现**（`git archive HEAD | tar -x -C /c/tmp/85-head`，**整树跟踪字节，不含本腿任何改动**，
   也不含共享树里任何未跟踪遗留）——同一条命令同一批用例，逐条 `--- FAIL`：
   `TestComposerContractTypesMatchFrontend` · `TestC21DesignTokensFourWayAgree` ·
   `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` · `Test229SpawnNamesAreRosterShaped`。
   逐字读数：`FAIL internal/panel`、`FAIL internal/risk`、`FAIL internal/tools`，`SNAPSHOT_TEST_RC=1`。
2. **被红包没有一个读 `ci.yml`**：`grep -rln 'ci\.yml' internal/panel internal/risk internal/tools` = **0 命中**
   ⇒ 本腿唯一的改动字节（`.github/workflows/ci.yml`）结构上不可能被这四条用例看到。

各自的红因（逐字摘，**都不是格式/本腿**）：
- `internal/panel::TestComposerContractTypesMatchFrontend`：`composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare`
  ⇒ Go 侧与 `frontend` 的 `ComposerState` 接口漂了（`git status -- frontend` **0 行**，是 HEAD 上的真漂移）。
- `internal/panel::TestC21DesignTokensFourWayAgree`：`tokens_fourway_test.go:726: docs/evidence/s1/c21-native-tokens.md:58: the table names Palette.GlassRing, but internal/ball/tokens.go holds no such colour literal in DarkPalette()`（同形 5 条）⇒ 四方表与 `internal/ball/tokens.go` 不一致。
- `internal/risk::TestC26RewriteAccountIsConsumedAtEverySecurityLeg`：`pathresolver_rewrite_account_test.go:200: cmd/wisp/panel_pump.go:283 reads Result.Canonical without reading the rewrite account`（4 条）⇒ 票 102 的洞被新消费者重开。
- `internal/tools::Test229SpawnNamesAreRosterShaped`：`internal/memory/l1_warm.go:221: Spawn 首参不是字符串字面量（l1ExtractGoroutine）且没有复审记录`；
  该行在 HEAD 上逐字就是 `reg.Spawn(l1ExtractGoroutine, l1ExtractOwner, nil, …)`，而
  `grep -c l1ExtractGoroutine internal/tools/spawn_roster_census_229_test.go` = **0**（名册没登记）。

⚠ 本腿零 `.go` 改动 ⇒"碰过的包"是**空集**；`cmd/wisp` **刻意没跑**（仓规：`go test ./cmd/wisp` 同一时刻只许一枚写腿跑，
且本腿不碰它；另见 `scripts/wisp-cli-tests.sh:7-16` 的 DLL 载入前提）。

**脚注③（CRLF）**：三把尺同向——`git status` 空、HEAD 版喂 gofumpt 空、`file` 说 CRLF ⇒ 只有工作树副本的换行符差异。
按 CI 的真口径（checkout 出来的 LF 跟踪字节）这把尺**是空的**，纯净快照已实测为空。

---

## 6. STOP / 本腿没做的格（具名、具因、具归口）

- **STOP-1 · AC#1 的后半句（真 run 日志行）**：要 `##[group]Run go install mvdan.cc/gofumpt@v0.12.0`
  与 `go: downloading mvdan.cc/gofumpt v0.12.0` 这类**托管 run** 的行。**本腿不 push**（铁律），⇒ 由编排者 push 后取。
  预期读数：`gofmt (gofumpt)` 步的日志首行是钉版串；该步**仍可能红**，但红因只能是真实格式问题，不再可能是"上游今天换了版"。
- **STOP-2 · AC#2 / AC#3（34 条 linux / 78 条 windows findings 的清算 + 变异回验）**：**明令不在本腿**
  派单里（"**不要**做全仓 findings 清算"），且跨 12+ 个包 —— 属**另一条腿/另一张票**（票面 `:95` 的 **85b**、`:219` 的 **票 122**）。
  本腿**零** findings 处置，**未加** `//lint:ignore`、**未**排除包、**未**改 `continue-on-error`。
- **STOP-3 · AC#4（步骤改名）**：`test-windows` 里 `PathResolver junction placeholder (real cases tickets 18/20)`
  名字里的 "placeholder" 已过期，改名本身在写面内（`ci.yml`），但 ① 本腿派单写死"**只钉 gofumpt**"，
  ② AC#4 还要求"给一次真 run 的该步结论"，同样要 push。⇒ 留给下一腿或同一枚 push 之后。
- **STOP-4 · `lint` 仍是红的**：这是票面 `:101-106` 的**裁定①**批准的结果（"85a 落地之后 `lint` 继续红——
  从读不懂的红变成读得懂的红"）。**本腿不让 CI 变绿，也不该让它变绿**（那要靠清账）。
- **STOP-5 · 票面本身未动**：本腿写面不含 `.scratch/wisp/issues/85-*.md` ⇒ Status/Progress log/勾框**留给收件一步**。
  收件时要补的一条读数：`@latest` 2→0、无版本 `go install` 1→0（§2/§4-D）。

- **STOP-6 · 没加 `_test.go`（写面允许，但本腿判定不该加）**：想要的守卫是"`ci.yml` 里每一行 `go install` 都必须带 `@版本`"，
  今天靠**一条可复跑的单行尺**（§4-D）就够。要变成常驻测试，得先有一个**读仓库根 `.github/` 的测试落点**——
  本仓现成的读 `ci.yml` 的用例全在 `cmd/wisp`（`dataroot_128_test.go` / `logsink_test.go`），
  而 `cmd/wisp` **整包在本腿禁改面内**；`tools/d22scan/**` 同样禁改。⇒ 为一行 YAML 新开一个包不成比例，
  **留给"要不要给 CI 版本 pin 立常驻仪器"这一格由编排者裁**（本腿不擅自立）。

---

## 7. 撤销口令

一句话复原：**把 `:293` 的 `gofumpt@v0.12.0` 改回无版本串（`:101` 注释同步），即回到落刀前。**
按仓规以追加提交更正（不 `--amend`、不 `reset`）：`git revert <本腿的 ci.yml commit>`；
若只想局部回退而不留提交，改这一行即可，**不要**用 `git checkout .` / `reset`（共树 A34）。
