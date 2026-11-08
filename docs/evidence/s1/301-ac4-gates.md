# 票 301 — AC#4 门禁读数 ＋ AC#4b「显式 pathspec」正向闸门落地

- 腿名＝`301-ac4b`；钟点＝2026-11-08 14:59:16 +0800；基线锚＝`31482f271`（`A1633: wave86 three legs landed`）
- 门禁环境＝`git clone --shared` 导出树 `/d/tmp/wisp/20261108/301-ac4b-tree`（⛔ 整包门禁不在工作树跑），每发跑前 `unset GOTMPDIR TMPDIR TEMP TMP && export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，⛔ `CGO_ENABLED=0` 未设
- 写面＝本证据件＋票 301 Progress log 一行；`scripts/**` 零改动（理由见 §3 结论：门禁无缺口）；零产码；只 commit 不 push
- 三发门禁原始日志存档（仓外）：`/d/tmp/wisp/20261108/301-ac4b-d22scan.log`（27,098 字节）／`301-ac4b-census.log`／`301-ac4b-census-linux.log`／`301-ac4b-gofumpt.log`／`301-ac4b-gofumpt-prod.log`

## 1. 起手三核（尺逐枚现跑，⛔ 信转述）

1. **票面原文仍在**：AC#4＝票 md `:23` 未勾；AC#4b＝票 md `:151` 未勾（本腿现读全文核对）。
2. **引文出处**（行号锚全漂，按裁定节用内容锚现跑）：`grep -n 'portable-tests.sh --scope' .github/workflows/ci.yml` ⇒ `:708 core`／`:1057 census`／`:1093 windows`（census 调用点仍在、仅此一枚）；`sed -n '406,409p' scripts/portable-tests.sh` ⇒ `census)` 档＝`scope=()`＋`pinned=''`。
3. **红/绿非他人在飞件造成**：主工作树 `git status --porcelain` 非 `??` 行＝0；未跟踪在飞件＝`cmd/wisp/panel_resident_323_windows_test.go`／`cmd/wisp/panel_resident_assembly_323_test.go`／`tools/jarvis-export/.verify.js`（另一腿的活），只在工作树、⛔ 进 `--shared` 导出树 ⇒ 三发门禁读数＝纯 HEAD 提交态。

## 2. 门禁一：`sh scripts/d22scan.sh`

- 命令（逐字）：`sh scripts/d22scan.sh`
- **rc=0**
- 末行（逐字）：

  ```
  d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=387, bans #1-5 cmd/=64, ban #6 frontend/=119, ban #7 internal/tools/=48, ban #10 internal/risk/=21, ban #8 design/=14, ban #8 frontend/=119, ban #8 internal/=963, ban #8 cmd/=241; ban #8 emoji coverage: design/ 14 text files; frontend/ 119 text files; internal/ 963 Go files, comments and _test.go included; cmd/ 241 Go files, comments and _test.go included
  ```

- 扫描射程＝`design/ frontend/ internal/ cmd/`（末行自declare）⇒ 本笔只落 `docs/**`＋`.scratch/**`，**post-commit 树的该读数逐字不变**（⛔ 另跑一发假装"交件后"）。

## 3. 门禁二：`bash scripts/portable-tests.sh --scope=census`

- 命令（逐字）：`bash scripts/portable-tests.sh --scope=census`
- **rc=1（GUARD D 红，5 枚包）——★如实报读数**
- census totals 行（逐字）：

  ```
  portable-tests.sh: census totals: packages=41 with-zero-compiled-tests=4 claimed-by-no-scope=9 unclaimed-with-tests=5
  ```

- GUARD D 名册（逐字，5 枚）：

  ```
  portable-tests.sh: GUARD D - 5 package(s) compile a test file for GOOS=windows and
  portable-tests.sh:   NO named scope claims them, so no CI step runs them and no CI step
  portable-tests.sh:   can ever go red over them:
  portable-tests.sh:   github.com/CarlosShao/wisp/cmd/balldebug  tests-compiled-for-windows=1/0
  portable-tests.sh:   github.com/CarlosShao/wisp/internal/agent/scheduler  tests-compiled-for-windows=2/0
  portable-tests.sh:   github.com/CarlosShao/wisp/internal/command  tests-compiled-for-windows=1/0
  portable-tests.sh:   github.com/CarlosShao/wisp/internal/plugin/manifestslot  tests-compiled-for-windows=7/0
  portable-tests.sh:   github.com/CarlosShao/wisp/internal/speech  tests-compiled-for-windows=16/0
  ```

- **STALE 腿＝census 上不存在**（票面更正节在本锚复核成立）：census 分支 `:406` 进入 → totals → GUARD D（红则 `exit 1`）→ `:559 exit 0`，**先于** ledger/stale 循环（`:757` `stale=()`；stale 判红在 `:801` 非零即 `exit 1`＝计退码，只活在 `--scope=windows`/`--scope=core` 档路径）。⇒ 派单里「STALE 只列不判红」在 census 上**无对象可列**；本发 log 里 `stale` 字样命中＝0（grep 现量）。
- **红归属（三核③，逐笔现核）**：
  - **非本腿造成**：导出树＝纯 `31482f271`，本腿交件前对它零改动。
  - **引入 commit 具名**（尺＝`git log --oneline --diff-filter=A -1 -- <file>` 现跑）：`5438d104e`（2026-11-03 票52，`internal/plugin/manifestslot` 7 枚测试件）／`08c4c51e2`（2026-11-05 票71，`cmd/balldebug/diffgate_windows_test.go`）／`84e4cc78d`（2026-11-05 213-r1，`internal/command/roster_test.go`）／`1cc67664b`（2026-11-06 票47，`internal/agent/scheduler` 2 枚测试件）。
  - **`internal/speech` 是脚本自己声明过的已知账**：`scripts/portable-tests.sh:266-277` 逐字含 `./internal/speech/ IS DELIBERATELY NOT CLAIMED, here or in any other tier` ＋ `It stays on the census as UNCLAIMED-HAS-TESTS until such a carrier claims it`（等 carrier 阶 sherpa DLL＋models，硬认领会造永久红）。
  - **票面 10-10 记 census rc=0 属实**：`772ff880` 上 `internal/agent/scheduler` 只有 `doc.go`（桩）、`internal/command`／`internal/plugin/manifestslot`／`internal/speech` 包尚不存在（`git ls-tree --name-only 772ff880 …` 现量）⇒ 红是 11-03…11-06 之间随上面四笔落进来的。
  - **旁证**：票 310 验收件 `docs/evidence/s1/ticket310-310-v1-accept-v1.md`（2026-10-12）里 census 读数＝`packages=35 … unclaimed-with-tests=0`（绿）——绿→红的窗口与本 attribution 吻合。
  - **门禁无缺口**：GUARD D 抓的正是「包里有测试、没有档跑它」，5 枚全是真缺口，⛔ 是仪器漏扫 ⇒ **本票不需要改 `scripts/**`**；修口＝把各包拉进档＋同一笔改 pin（GUARD C 同笔约束），归属＝各引入票／票 111 `AC#12` 族，⛔ 本腿代做（本票 `AC#1` 子句族明文禁顺手改别的档的清单或 pin）。
- **CI 等价读数本机取不到**：`GOOS=linux bash scripts/portable-tests.sh --scope=census` 在本机死于 `go list`（`sherpa-onnx-go-linux` 在本机模块缓存里 build constraints exclude all Go files；日志存档 `301-ac4b-census-linux.log`），脚本自declare census wired on the linux leg ⇒ CI 的 ubuntu census 现色只能看 CI；本腿如实报本机 GOOS=windows 形这一发。

## 4. 门禁三：gofumpt

- 口径＝产码面（此前腿例同形：票 319 `-l` 其 Go 件、票 298 `-l cmd/wisp`）：
  - `/d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools`（导出树根）⇒ **输出 0 行、rc=0（空）**。
  - 全树 `-l .` 另发＝29 行，全部 ∈ `.scratch/wisp/probes/**`（变异探针，刻意坏格式）＋`forks/go-webview2/**`（fork 源码），⛔ 产码面（读数存档 `301-ac4b-gofumpt.log`，如实附）。
- 本腿零 Go 写件 ⇒ 本枚门禁对本笔恒真。

## 5. AC#4b 落地：「显式 pathspec」子句（正向闸门，替代追溯判据）

**一句话形**（给今后每一条写腿的派单逐字带）：

> **每笔 commit 的 pathspec 必须写在 `$( … )` 之外、逐枚点名只落在该腿自家射程内的具体文件；交件时随件交 `git show --name-only --format= <sha>` 的逐枚路径名册＋与派单授权名册的双向差集；编排者收件时自己重跑名册、把读数落进票面，⛔ 以腿的自报为凭。**

三件各自可核（判据＝各留一枚盘上证据：件名＋那把尺逐字；缺任一件＝本格 ⛔ 翻）：

1. **① 逐笔名册交回**：尺＝`git show --name-only --format= <sha>`（**逐笔**；⛔ 区间尺 `git diff A..B`——会把编排者夹在中间的笔算到腿头上，`301-v1` §3-⑥ 已证）。随件交＝名册＋与授权名册**双向差集**。
2. **② pathspec 形状**：写在 `$( … )` 之外；只落腿自家射程（⛔ 票级 glob `probes/<ticket>/**`；编排者补位件落 `probes/orch/<ticket>/`，⛔ 再落进任何在飞腿的射程——`A811④`/`A814` 事故的处置）。
3. **③ 编排者复跑**：编排者收件时把每笔名册自己打一遍并把读数落进本票——本腿交件只供 ③ 的输入，⛔ 自我盖章。

## 6. 本笔自身按 ①② 交件

- **授权名册**（本腿派单逐字）＝`scripts/**`（仅当必需）＋`docs/evidence/s1/301-ac4-gates.md`＋票 301 Progress log 一行。
- **本笔取集**（预先声明，交件后可用 ① 的尺现核）＝2 枚：
  - `.scratch/wisp/issues/301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md`（只追加 Progress log 一行；AC 勾与 Status 零动）
  - `docs/evidence/s1/301-ac4-gates.md`（本件，新建）
- **双向差集＝空**（取集 2 枚全部 ∈ 授权面；授权面里 `scripts/**` 本笔零取＝合法，「仅当必需」且门禁无缺口，零产码合规）。
- **pathspec 形状**（逐字，两处都在 `$( … )` 之外、逐枚点名）：
  - `git add -- "docs/evidence/s1/301-ac4-gates.md" ".scratch/wisp/issues/301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md"`
  - `git commit -m "<message>" -- "docs/evidence/s1/301-ac4-gates.md" ".scratch/wisp/issues/301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md"`
- **③＝编排者收件时现跑**（输入＝本节＋交件汇报里的 sha；sha 在 commit 前不可知，⛔ 腿预写）。
