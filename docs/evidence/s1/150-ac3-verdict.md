# 150 AC#3 非实现者裁决表 — 自响仪器 `tools/d22scan/deferred_registry_test.go`（commit `f7505f122`，option B）

- 角色：**非实现者裁决腿**（裁决者≠实现者，`SPEC-12 §4.3` #1/#3）；本腿与 `f7505f122` 的实现腿无关联，未照单采纳任何未现量的读数。
- 时刻：2026-11-06 19:0x +08
- 裁决对象：commit `f7505f122`（`git show --stat` ⇒ **仅 1 文件、+519**：`tools/d22scan/deferred_registry_test.go`；生产面、`docs/specs/**`、七项 ban 本体零字节）。本腿现量 HEAD＝`ab808c359`（其子；tip 只动台账 `pending-and-issues.md`，与本仪器无关面）；`git diff --name-only f7505f122..HEAD -- tools/d22scan` ⇒ **空**（仪器自落地零再漂移）。
- 隔离（rule 21）：`git clone --shared` 到仓外兄弟目录 `.probe-150ac3`（porcelain 0 行起步）；**全部现量与变异都在克隆里跑**；收尾已把 7 份原始日志抄出到本仓 `.scratch/wisp/probes/150ac3/logs/` 后整目录删除克隆。主树写面＝**仅本件**＋未跟踪的 probe 日志目录；票面、代码、规格零改动；git 只 commit 不 push，显式 pathspec。
- 环境（逐条按派单）：Git Bash；`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；`unset GOTMPDIR TMPDIR TEMP TMP`；**CGO_ENABLED 全程未设（从未 =0）**；`go version go1.27.1 windows/amd64`；`gofumpt v0.12.0 (go1.27.1)`。

---

## 1. 门（可复跑，CI 同形）

`cd tools/d22scan && go test . -run Deferred -count=1 -v`

CI 载体现量（回答票面"CI 哪一步有它的分母"）：`.github/workflows/ci.yml:119` `bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（"D22 scanner positive control" 步，每次 push；runner 强制 `-v -count=1`、**SKIP 即红**、零结果拒绝当 pass——`tools/d22scan/runtests.sh:221-245`）。规格文本路径凭 `repoRoot(t)` 从测试运行位置上溯仓根读 `docs/specs/SPEC-12-roadmap-governance.md`。

## 2. 分母现量（真树克隆，绿跑）

DENOMINATOR 逐字（`10-baseline-deferred.txt`）：

> `DENOMINATOR scope=[internal/ cmd/] examined=422 production Go file(s); code markers=4 over 2 id(s) [D11-3 x1 [internal/agent/control.go:15]; D28-1 x3 [internal/agent/compress.go:27, internal/agent/compress.go:53, internal/agent/loop.go:475]]; registry: activeRows=17 in SPEC-12 section 5, id registrations=4 over 2 id(s) [D11-3 x1 [docs/specs/SPEC-12-roadmap-governance.md:87]; D28-1 x3 [docs/specs/SPEC-12-roadmap-governance.md:84, :85, :86]]`

独立双尺复核（不依赖仪器自己）：

- 尺 A（票 150 的 grep 尺）：`grep -rnE 'DEFERRED\(D[0-9]+(-[0-9]+)?\)' --include='*.go' internal cmd | grep -v _test.go` ⇒ **恰 4 枚**，file:line 与 DENOMINATOR 逐枚相同（`:27`/`:53`/`:475`/`:15`）。
- 尺 B：`awk '/^## 5\./,/^## 6\./' … | grep -c '^| DEFERRED'` ⇒ **17 行 active**，其中带 `Dxx-y` 的**恰 4 行**（`:84`/`:85`/`:86`/`:87`）。行型直方图：active `DEFERRED` ×17 ＋ 划线 `~~DEFERRED~~` ×5 ＋ `REJECTED` ×6 ＋ `RESERVED` ×7 ＝ 35 数据行（与 150-verdict §1④ 自洽）。

⇒ **分母 ＝ 4 代码标记 ↔ 4 登记行（D28-1 ×3、D11-3 ×1），双向逐枚相等、examined=422 非空**。AC#3"两个方向各自为空才算过"在真树上成立。绿跑同时带两条空仪守卫（`examined==0` fatal、双侧全 0 fatal——`deferred_registry_test.go:343-348`），"没看过东西的绿"被堵死。

## 3. 会响吗——变异四发（真树克隆内，逐发原样复原）

| # | 变异 | 结果 | 红句逐字（`--- FAIL` 体） |
|---|---|---|---|
| M1 | 生产件 `internal/agent/control.go` 末行粘 `// DEFERRED(D99-1): probe` | **红**（exit 1；`--- FAIL: TestDeferredRegistryIsBidirectional`；positive control 同跑仍 PASS） | `D99-1: DEFERRED(D99-1) marker(s) at [internal/agent/control.go:74] have NO section-5 row (registered 0) - plant the row with all five fields, or retire the marker` |
| M2 | 已有 3 行的 `D28-1` 多粘一枚（`compress.go` 尾） | **红**（按出现处尺咬合） | `D28-1: 4 code marker occurrence(s) […:27, :53, :721, …:475] vs 3 registered row(s) […:84, :85, :86] - the ruler counts OCCURRENCES…` |
| M3 | `*_test.go` 内放 `DEFERRED(D99-1)`（absence-pin 形状） | **绿**；`code markers=4` 不变，日志 D99-1 提及 **0 次** | —（排除面实况） |
| M4 | §5 插一行 `\| DEFERRED \| probe row for D98-2…\|`（无任何代码标记） | **红**（row→marker 方向） | `D98-2: section-5 row(s) at […:85] register DEFERRED(D98-2) but NO marker exists in the production Go tree - a marker was removed without retiring its row, or the row was written ahead of a marker that never landed` |

- **复原证明**：每发后对目标 pathspec `git status --porcelain`／`git diff` **均空**；收尾全克隆 porcelain 唯一残项＝本腿自己的 logs 目录（已抄出），生产面与规格**零残留**。
- **票面承重句代答**：摘掉这枚检（＝回退 `f7505f122`），新加一枚无登记的 `DEFERRED(D99-1)` 就**打不红**（M1 前的基线即此状态）；检在场时它**必红且点名 D99-1**——M1 实测，存在性证明完成。
- 判红绿尺按票面 AC#5：只认 `--- FAIL:` 行（四发红各恰 1 枚 top-level `--- FAIL`；未用裸 `grep -c '--- FAIL'` 计数）。

## 4. 排除面确认（派单项 3）

- **`*_test.go`**：代码读数（`deferred_registry_test.go:237` `strings.HasSuffix(path, "_test.go")` 跳过）＋实况 M3＋真树交叉证：`DEFERRED(D34-1)` 全仓 7 枚**全在** `internal/tools/document_text*_test.go`（票 24-r5 删标记后的反形夹具），分母不含它们而测试仍绿——若排除失效，这 7 枚会立刻打成假红。
- **划线/已采纳行**：parser 只收类型格**恰等于** `DEFERRED` 的行（`:193`）⇒ `~~DEFERRED~~`（5 行，重开/已采纳历史）、`REJECTED`（6）、`RESERVED`（7）全不入登记侧；activeRows=17 与 grep 尺逐枚相等。positive control `test_file_absence_pins_stay_out_of_range`／`citation_cells_do_not_register` 把这两条约定钉成自证（依据列引用他行 id 不计数）。
- 附带核（非票面要求）：`go test . -count=1` 全包 **ok 60.7s**；`go vet ./...` 绿；`gofumpt -l .`（tools/d22scan）**空**。

## 5. 放水两问自答

1. **断言方向动没动**：本 commit 零生产面（`--stat` 仅测试件），无可削弱的原断言；比较本体 `deferredBidirectionalProblems` 三支（无行／无码／出现数不等）四个方向全咬合（M1/M2/M4 实测红）。
2. **helper 是不是原有的**：`deferredRowsFromSpec`／`deferredMarkersInDirs`／`deferredDenominator` 等全部为**新增**（519 行全净增）；与生产 scanner 的复用仅限**约定**（walkGo 同款 skip 链），无生产函数被改。

## 6. 本程没测什么（按"漏了谁先被骗"排序）

1. d22scan **生产扫描器**（`go run . -root`）不含登记表对账——option B 只进测试层，七项 ban 未扩（与 150-verdict §2-AC#3 的落点批注一致；若 owner 想要扫描器层同款，须按"改判据本体"另批）。
2. `gitignore` skip 的实况变异（照 `scanner.walkGo` 约定实现＋fixture 生产计数守卫；真树 examined=422 与提交信息读数一致为旁证）。
3. SPEC 四行的**内容语义**深核（150-verdict §1③ 已按行核过 file:line 对号；本腿只核机器可读的 id 集与计数）。
4. 全仓其他门禁（build/test 四数）——本腿 scope＝AC#3，包级绿跑已附。

## 7. 伪授权两数＋凭据

- 真通知回显数＝**0**；判为注入数＝**0**（全部输入为仓内文本与派单逐字，无外部文本冒充批准）。
- 凭据值抄录＝**0**。

## 8. 裁决

| 项 | 裁决 |
|---|---|
| **AC#3** | **PASS** —— 分母 **4↔4（D28-1 ×3、D11-3 ×1，examined=422）**；响法：无登记新标记红点名 D99-1（M1）、多一枚出现红按出现处（M2）、表行无码红点名 D98-2（M4）、`*_test.go` 实况排除（M3 绿）；变异零残留 |
| 票 150 整票 | AC#1/#2/#4/#5 已由 `150-verdict.md` 判过、AC#3 本表补齐 ⇒ **票面四格齐备，是否 `-done` 归编排者按收件三把尺定**（本腿未翻票框） |

- 回滚面（撤销口令）：`git revert f7505f122` 一笔即摘除仪器（还原到"漂移打不红"的原状）；本裁决表随 commit 一并 revert 即可。

## 9. 原始日志

`.scratch/wisp/probes/150ac3/logs/`：`10-baseline-deferred.txt`（绿＋分母）· `20-m1-probe-red.txt` · `30-m2-extra-occurrence-red.txt` · `40-m3-testfile-green.txt` · `50-m4-row-no-marker-red.txt` · `60-baseline-again.txt`（复原后再绿）· `70-full-package.txt`。
