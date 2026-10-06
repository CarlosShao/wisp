# 票 111 — 111-r4 证据件（交票面 AC#1 那张全仓对账表）

代号 `111-r4`（只读普查腿）。本腿唯一交付物＝票面 **AC#1 的全仓对账表**，另附 (a) `go list ./...` 枚数现量＋35−33 差集、
(b) "空分母"那一族的现状、(c) winlive 半边与 cmd/wisp ubuntu 半边的现状登记。
前三格（闸门写进脚本 `1bb654e3`／闸门接进 CI `6c0e3e31`／闸门配自检 `98f62fde`）**本腿不重做**，只引其行号与 blob 凭据。

⛔ 本腿零改动面：`.github/**`、`scripts/**`、`cmd/**`、`internal/**`、票面五枚 `- [ ]` 框、台账、`docs/reports/**`、`docs/evidence/s1/**` 一字未动。
⛔ 本腿零 `go test`／`go build`／`go vet`；`go` 命令只用了 `go list`（含 `go list -f`，不编译测试二进制），给并行腿 `231-r1` 的计时窗让路。
写面＝`.scratch/wisp/probes/111/r4/**` ＋票 111 Progress log 一行。

---

## §0 起手锚（写满）

| 项 | 读数 | 时刻（+0800） |
|---|---|---|
| `git rev-parse HEAD` | `a16d1ff7291fa9ab59f328cbf5f2429184a881f4` | 09:14:54 |
| `git status --porcelain -- .github scripts cmd internal` | **0 行** | 09:14:54 |
| `git status --porcelain` 全仓 | **734 行**（共树在飞；本腿未 add 他人任何路径） | 09:14:54 |
| `go list ./...` 枚数现量（本机 GOOS=windows） | **35**（`| wc -l`＝35，stderr 0 字节） | 09:07:54 |
| `git rev-parse HEAD:scripts/portable-tests.sh` | `2ff02dd7476c8f623a4489932d0d0ff95480e667` | 09:14:54 |
| `git rev-parse HEAD:.github/workflows/ci.yml` | `014861149bde08ee8e4bb6e988f8a605796a856f` | 09:14:54 |
| `bash -n scripts/portable-tests.sh` | rc=**0**（语法自证，只读） | 09:14:54 |
| `bash -n scripts/portable-tests-selftest.sh` | rc=**0** | 09:14:54 |

★**锚点稳定性**（免得被读成"我在 734 行的脏树上量的门"）：上面两枚 blob 在 `1bb654e3`／`15d8b60e`／`6c0e3e31`／`b3d29b0d`／起手 HEAD
**五枚 HEAD 上逐字节同一**（`git rev-parse <ref>:<path>` 逐枚复量＝`2ff02dd7` 与 `014861149bde`）。
⇒ 本件里所有 `scripts/portable-tests.sh` 与 `ci.yml` 的行号对这五枚同时成立，不需要随 HEAD 漂移重量。
`go list ./...` 的**名册**（35 枚包路径）在 `15d8b60e..起手 HEAD` 之间同样未变（该区间 `-- cmd internal frontend tools` 的 `.go` 改动 4 枚，全在 `cmd/wisp` 包内，无新增/删除包目录）。

## §1 AC#1 全仓对账表（`go list ./...` 的每一行都在表里）

**自证完整那把尺**：表体 **35 行** ＝ 起手锚 `go list ./...` 枚数现量 **35**（本机 GOOS=windows，留盘 `…/r4/golist-win.txt`，35 行）。
复算方法（可重跑）：把表 ① 列抠出排序 ＝ `golist-win.txt` 去前缀排序 →
`comm -13` **空**（名册里没有一行没进表）＋ `comm -23` **空**（表里没有一行不在名册上）⇒ 无漏行、无造行、删除列 0。原始 TSV 留盘 `…/r4/census-table.tsv`。

四列口径（逐列给尺，避免被读成凭记忆）：
- **① 包路径**＝`go list ./...` 原行去掉模块前缀 `github.com/CarlosShao/wisp/`。
- **② `_test.go` 枚数**＝`git ls-files <pkg> | grep -c '_test\.go$'`（本腿尺＝跟踪文件、精确目录名；
  与 r1b 的 `git ls-files <pkg>` 前缀尺差在父包会吞子包测试：`internal/agent` 前缀尺 38／精确尺 19，`internal/llm` 前缀尺 21／精确尺 10）。
  同格第二数＝**本平台编译进测试二进制的文件数** `go list -f t/x`（留盘 `…/r4/compiled-tests-win.txt`，GOOS=windows）——它 ≤ 跟踪数（build tag 排除的落在这里差），GUARD A 用这把尺。
- **③ 在不在 CI 某一步**＝逐 scope 给 `scripts/portable-tests.sh`（HEAD blob `2ff02dd7`）里的 **pin 行号**；
  ci.yml 调用点是固定的，列在表下"调用点图例"。★**census 档没有 pin**（`:367` 明写"审计者不认领档"），
  它的角色是**审这四档**、给每行打 `CLAIMED BY`／`NO-SCOPE`，所以这一列对 census 一律记"审计者，非认领"。
- **④ 那一步真给过结论的 run id + step 号**＝⚠**语义＝"已推送配置给过的结论"**，读法见 §5。
  已推送 tip＝`origin/dev = c6cf66e6`（2026-10-05 10:58），其 ci.yml blob `c5a1b1ab`、portable-tests.sh blob `5ddb70f8`——
  ★**census 那一步在已推送配置里不存在**（`grep -c scope=census c6cf66e6:ci.yml`＝0；`GUARD D` 亦 0，它由 `1bb654e3`/`6c0e3e31` 引入，都在 tip 之后）⇒ 本列凡涉及 census/`session`/`projctx` 者一律"**待推送后回填**"。

| ① 包 | ② 跟踪 `_test.go` / 本平台编译 t+x | ③ 认领档（pin 行） | ④ 给过结论的 run / step（已推送配置） |
|---|---|---|---|
| cmd/balldebug | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| cmd/llmrecord | 1 / 1/0 | core :165 · windows :194 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| cmd/wisp | 65 / 58/0 | cli :206 | run 37257547572 / test-windows step7（CLI 门）**failure**，own-line `FAIL` |
| frontend | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| internal/agent | 19 / 19/0 | core :166 | run 37257547572 / test-core step7 `ok` |
| internal/agent/approval | 19 / 8/11 | core :167 | run 37257547572 / test-core step7 `ok` |
| internal/agent/scheduler | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| internal/audio | 5 / 5/0 | core :168 | run 37257547572 / test-core step7 `ok` |
| internal/ball | 16 / 11/0 | core :169 · windows :195 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/buildinfo | 1 / 1/0 | core :170 | run 37257547572 / test-core step7 `ok` |
| internal/config | 18 / 16/1 | core :171 · windows :196 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/llm | 10 / 2/8 | core :172 | run 37257547572 / test-core step7 `ok` |
| internal/llm/adaptertest | 1 / 1/0 | core :173 | run 37257547572 / test-core step7 `ok` |
| internal/llm/anthropic | 3 / 3/0 | core :174 | run 37257547572 / test-core step7 `ok` |
| internal/llm/golden | 1 / 1/0 | core :175 | run 37257547572 / test-core step7 `ok` |
| internal/llm/openaichat | 4 / 4/0 | core :176 | run 37257547572 / test-core step7 `ok` |
| internal/llm/openairesponses | 2 / 2/0 | core :177 | run 37257547572 / test-core step7 `ok` |
| internal/memory | 11 / 11/0 | core :178 | run 37257547572 / test-core step7 `ok` |
| internal/models | 17 / 16/0 | core :179 | run 37257547572 / test-core step7 `ok` |
| internal/observe | 15 / 15/0 | core :180 | run 37257547572 / test-core step7 `ok` |
| internal/panel | 20 / 20/0 | core :181 | run 37257547572 / test-core step7 **`FAIL (own line)`**（4 枚 panel 用例红，见 §6.3） |
| internal/perm | 3 / 3/0 | core :182 · windows :197 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/plugin | 1 / 1/0 | core :183 · windows :198 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/proc | 11 / 11/0 | core :184 · windows :199 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/projctx | 1 / 0/1 | core :185 · windows :200 | **待推送后回填**（`1bb654e3` 才进 pin，tip 之后；已推送 core step7 表无此行） |
| internal/risk | 22 / 20/0 | core :186 · windows :201 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 **`FAIL (own line)`** |
| internal/secret | 4 / 4/0 | core :187 · windows :202 | run 37257547572 / test-core step7 `ok`；/ test-windows step8 `ok` |
| internal/session | 2 / 2/0 | core :188 · windows :203 | **待推送后回填**（同 projctx） |
| internal/speech | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| internal/statemachine | 2 / 2/0 | core :189 | run 37257547572 / test-core step7 `ok` |
| internal/streamkey | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| internal/tools | 49 / 46/3 | core :190 | run 37257547572 / test-core step7 `ok` |
| internal/watchdog | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |
| internal/winsec | 27 / 12/8 | core :191 · winsec :209 | run 37257547572 / test-core step7 `ok`（POSIX 半边）；/ test-windows step4（ACL 门）`ok` |
| tools/signmodels | 0 / 0/0 | NO-SCOPE | 待推送后回填（仅 census 点名） |

**调用点图例**（ci.yml HEAD blob `014861149bde`，本腿 grep 现量）：

| scope | `run:` 行（HEAD blob `014861149bde`） | 所在 job | 步号：HEAD-yaml 本地序 ／ **已推送配置 API `.steps[].number`（第④列引的那两枚 run）** | 该步 `if:` |
|---|---|---|---|---|
| core | `:373` | `test-core`（ubuntu，`:310`） | yaml 7 ／ **API 7**（同值；本 job 无新增步，两套配置步序一致） | 无 `if`（其下只有带 `always()` 的 `Stop compose` `:376`） |
| winsec 门（`winsec-tests.sh`，非 `--scope=winsec`） | `:471` | `test-windows`（`:420`） | yaml 3 ／ **API 4** | `:436` `!cancelled()` |
| cli | `:507` | `test-windows` | yaml 6 ／ **API 7** | `:506` `!cancelled()` |
| **census（★本票新接）** | `:596` | `test-windows` | yaml 7 ／ **API 无此步**（已推送配置里 census 不存在 ⇒ 待回填；新配置上**预计 API number=8**） | `:595` `!cancelled()` |
| windows | `:632` | `test-windows` | yaml 8 ／ **API 8**（已推送配置；新配置上位移为 9） | `:631` `!cancelled()` |

★**表尾对账**：`go list` 35 枚 ＝ 28 枚本平台编译得出测试文件（35 − 7 枚 `0/0`）＋ 7 枚 `0/0`；
28 枚**全部被某档认领**（`NO-SCOPE` 的 7 枚恰好全是 `0/0` ⇒ `unclaimed-with-tests = 0`，GUARD D 今天不咬）。
测试文件总数：跟踪精确尺 Σ②＝**350**，本平台编译 t+x Σ＝**327**（差 23＝build-tag 排除的 `_test.go`，含 11 枚 winlive，见 §4）。

## §2 附件 (a) `go list ./...` 枚数现量 ＋ 35−33 差集逐枚坐实

- **现量：windows＝35**（09:07:54，rc=0、stderr 0 字节，逐行留盘 `…/r4/golist-win.txt`）。
  ★与 r1b §0（10-05 17:0x）与 r2 §1.1（10-05 20:24:19→21，量到 35）**三腿独立复现同一数**。票面 `现场`／AC#1 写的 **33 ⇒ 过期**。
- **平台必须带**（本腿复量）：`GOOS=linux CGO_ENABLED=0 go list ./...` rc=**1**、stdout **34 行**、stderr **264 字节**
  （`cmd/wisp -> sherpa-onnx -> build constraints exclude all Go files in …/sherpa-onnx-go-linux@v1.13.8`）。
  ⇒ linux 名册 **34**，少的正是 `cmd/balldebug`（三枚 `.go` 全带 `//go:build windows`，`git ls-files cmd/balldebug` 现量）。
  这条是本腿**第四次**复现该障碍（111 原腿 AC#4／r1b §0／r2 §1.2／本腿），非本腿可修面。
- **35 − 33 ＝ 2 枚，逐枚坐实**（不跑第二次 `go list`，只用 `git ls-tree` 对票面锚 `4e66817`＝2026-09-24 09:46 复量）：
  主模块（根 `go.mod`＝`github.com/CarlosShao/wisp`）内含 `.go` 的目录数 **33 → 35**，
  `comm -13` 差集＝**`internal/projctx`**、**`internal/streamkey`**；`comm -23` **空 ⇒ 删除列 0 枚**（11 天里没有任何包消失，差集不会被"改名"糊过去）。

  | 新增枚 | 首枚入库 commit（`git log --diff-filter=A` 现查） | 有无 `_test.go` |
  |---|---|---|
  | `internal/projctx` | `84feec44`（2026-09-28 19:30，票 200） | **有 1 枚**（外部测试包，`go list -f` t/x＝`0/1`）⇒ 真分母；`1bb654e3` 已补进 core+windows |
  | `internal/streamkey` | `335b8d2b`（2026-09-28 21:53，票 197 载体层） | **0 枚**（`git ls-files internal/streamkey` 只有 `streamkey.go` 1 枚，`go list -f` t/x＝`0/0`）⇒ 无分母 |
- 独立模块**不在** 35 里（`git ls-files | grep go.mod` 现量根／`scripts/spike`／`tools/d22scan`／`tools/mockllm` 四枚 module 行独立）⇒ 数名册时不许把它们算进来。
- ★**票面 33 在它自己写下的那天是对的**：不是漏计，是这两枚新包把数抬上去的（r2 §1.2 同判，本腿用 `git ls-tree` 独立复算）。

## §3 附件 (b) "空分母"那一族的现状（GUARD A 与 GUARD D 各咬什么）

**先给族名现量**：本平台（GOOS=windows）`TestGoFiles+XTestGoFiles` 全空＝**7 枚**，逐枚都是 NO-SCOPE：
`cmd/balldebug` · `frontend` · `internal/agent/scheduler` · `internal/speech` · `internal/streamkey` · `internal/watchdog` · `tools/signmodels`。
⇒ **"在清单里但一个测试文件都没有"那一族现在＝0 枚**（票面点名的 `internal/session`/`internal/watchdog` 已分道：
session 有 2 枚测试、已由 `1bb654e3` 认领进 core+windows；watchdog 是空分母且**不在**任何 scope 里 ⇒ 见下表"谁咬它"）。

| 守卫 | 行号（HEAD blob `2ff02dd7`） | 咬什么 | 本腿复算的现量 |
|---|---|---|---|
| **GUARD A** | 测量 `:570-571`，红句 `:574-579`（`EMPTY <pkg>` ⇒ 该 scope 声明了它却本平台编译不出任何测试文件） | **只对"声明在 scope 里的"包**发作（`${scope[@]}`） | 本腿把四档 scope 原样喂给同一条 `go list -f`（core `:235-245`／windows `:249-253`／cli `:260`／winsec `:268`＋`:229`）：**core 档 GOOS=linux 与 GOOS=windows 各一发、windows/cli/winsec 各一发，EMPTY 命中全部 0 行** ⇒ 今天**没有任何一档带着空分母条目** |
| **GUARD D** | 计数 `:396`、标记 `<-UNCLAIMED-HAS-TESTS` `:397`、红句 `:409-421`、`exit 1` `:422`；`0/0` 的豁免分支 `:392-393` | census 期间"**本平台编译得出测试文件、却没有任何命名档认领**"⇒ 拒绝 `exit 0` | `unclaimed-with-tests=0`（表尾对账）⇒ **今天不咬**。它咬不到的方向＝`0/0` 那 7 枚（明写在 `:385-391` 注释里：那是 A 的地盘，给它们进 scope 才是 AC#3 拒绝的假认领） |

★**两把尺的分工（票面 AC#3 那句"响亮"的真身）**：A 管"**名单里有名字但没分母**"（假认领），D 管"**有分母但名单里没名字**"（漏认领）。
两头夹住之后，剩下第三种形状——"没分母也没名字"（那 7 枚 `0/0` NO-SCOPE）——**A/D 都不咬，这是设计而非漏洞**（`:385-388` 逐字："a scope entry for it would be the false claim ticket 111 AC#3 refuses"）。
这 7 枚的现状登记（枚枚可核）：
- `internal/agent/scheduler`（`doc.go` 一行，`DEFERRED(scheduler): implemented by ticket 47`）
- `internal/speech`（`doc.go`，`DEFERRED(engines): implemented by ticket 15 / 26 / 41`）
- `internal/watchdog`（`doc.go`，`DEFERRED(watchdog loop/thresholds): implemented by ticket 42`）
- `frontend`（`embed.go`＋TS 前端，前端测在 `lint-frontend` job 的 `npm run` 步里，**不在 Go 名册的射程**）
- `cmd/balldebug`（3 枚 windows `.go`，0 测试）／`internal/streamkey`（1 枚实现，0 测试）／`tools/signmodels`（`main.go`，0 测试）
⇒ 前三枚是 DEFERRED 桩（实现票落地时**必须**同时放回 scope＋更新 pin，否则 GUARD C 显红＝111 原腿 `next=` 5③ 说好的那道耦合）；
后四枚是"有代码、今天没 Go 测试"的形状，**本票不动**（要动就得给包写测试＝超出票面 AC#1 的普查射程）。

## §4 附件 (c) winlive 半边 ＋ `cmd/wisp` ubuntu 半边现状登记

**(c1) winlive 半边——owner 未批，⛔ 本腿没动 `ci.yml`。**
- 文件名册（`git show :<file>` 逐枚 grep `winlive`，留盘 `…/r4/winelive-files.txt`）＝**11 枚**：
  `cmd/wisp` 6 枚（`panel_geometry_255_winlive_test.go`·`panel_host_windows_live_test.go`·`resident_approval_live_246_windows_test.go`·
  `resident_ball_live_228_windows_test.go`·`resident_hotkey_live_258_windows_test.go`·`resident_task_source_live_246_windows_test.go`）＋
  `internal/ball` 5 枚（`hotkey_cancel_borrow_live_260_test.go`·`hotkey_live_test.go`·`interaction_live_test.go`·`live_guard_windows_test.go`·`live_windows_test.go`）。
- ★**CI 今天零编译**（本腿四发复认）：`grep -cE 'winlive|-tags|GOFLAGS' HEAD:ci.yml`＝**0**；
  `grep -- -tags scripts/portable-tests.sh scripts/winsec-tests.sh scripts/wisp-cli-tests.sh tools/d22scan/runtests.sh`＝**0 命中**（没有任何一条链会带 build tag）；
  旁证＝`go list -f .IgnoredGoFiles` 在 GOOS=windows 下把这 11 枚**逐枚点名**为被忽略文件。
- ⇒ 票面 ② 列的 350（跟踪）− 327（本平台编译）＝**23 枚差**，本腿逐枚拆开（`git ls-files` 精确尺 vs `go list -f t/x` 两把尺相减）：
  **winlive 占 11 枚**（`cmd/wisp` 6 ＋ `internal/ball` 5）＋ 其余 12 枚＝单平台文件
  （`cmd/wisp` 另 1 枚 `//go:build !windows`／`internal/config` 1／`internal/models` 1／`internal/risk` 2／`internal/winsec` 7）。
- 纳入 winlive **需要动 `ci.yml`**（加 `-tags winlive`）＋多半要 self-hosted 真机；
  且这 11 枚里 **5 枚文件带共 13 处 `t.Skip`**（本腿逐枚 `git show :<file> | grep -c 't.Skip'`：`resident_hotkey_live_258`=4·
  `hotkey_live`=4·`interaction_live`=2·`live_windows`=2·`live_guard_windows`=1），而 `tools/d22scan/runtests.sh` 把**任何** `--- SKIP` 判 fatal
  （HEAD blob 现行号：计数 `:88`、红句 `:99` "SKIP is not a pass"；r1b §3 记的 `:98`/`:104` 已顶漂一枚）。⇒ **判"要动 ci.yml、owner 未批、本腿不纳入"**，只登记现状。

**(c2) `cmd/wisp` 的 ubuntu 半边——归票 98，本腿不动。**
- 认领现状：`cmd/wisp` 只在 **cli** 档（pin `:206`），唯一调用者 `scripts/wisp-cli-tests.sh`（`ci.yml:507`，test-windows step 7）。
- ★**脚本自己拒跑非 windows**：`scripts/wisp-cli-tests.sh:60` `if [ "$(go env GOOS)" != windows ]; then … exit 2`（`ci.yml` 里 cli 步只挂在 windows 腿）。⇒ ubuntu 腿今天**没有任何一步**跑 cmd/wisp。
- 障碍链（票 98 的地界）：windows 无 DLL ⇒ 加载期 `0xc0000135`、stdout 0 字节（看着像绿，防线＝wisp-cli-tests 的 DLL 预检 `:74-99`）；
  ubuntu `CGO=0` ⇒ 连 `go list` 都过不了（本腿 §2 那发 rc=1 复现）；`CGO=1` 建得出也起得来，但 111 原腿 AC#4 量到 **19/29 枚 FAIL、rc=1、127.4s**。
- ★**已推送配置上的第一个 CI 读数（本腿读 CI，不是本机冒充）**：run `37257547572`（head `c6cf66e6`＝远程 tip）/ job `111597759830 test-windows` /
  **step 7 `cmd/wisp CLI tests` conclusion＝`failure`**；该步四数 `=== RUN=323 PASS=219 FAIL=7 SKIP=1`，
  own-line 表打 `FAIL (own line) github.com/CarlosShao/wisp/cmd/wisp`。
  ★7 枚 FAIL 逐名逐文件（名取自日志，文件用 `git grep -l "func <名>(" HEAD` 现查，⛔ 不读工作树在飞版本）：
  `always_write_no_clobber_226_test.go`(TestAC1AlwaysBranchDoesNotRevertAHandEditedKey) ·
  `approval_always_201_test.go`(TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card) ·
  `config_reload_223_test.go`(TestTicket223HandEditedFsLooseningCostsAnL2Card) ·
  `config_reload_perm_223_windows_test.go`(TestTicket223PermissionDeniedSitsInItsOwnSentence) ·
  `instructions_200r2_test.go`(TestRunPacketCarriesTheLoadedInstructionFiles) ·
  `panel_host_windows_test.go`(TestPanelHostRealWindowHopAndLifecycle ＋ TestPanelHostLatencyPercentilesAC2)；
  1 枚 SKIP＝`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`panel_resident_windows_test.go`）。
  ⇒ 7 枚红**全是 cmd/wisp 包自己的用例**，且**不是**"缺 DLL 的加载期症状"（那条会打 `GUARD - pinned DLL(s) absent` 并 `exit 1` 在跑之前，见 `wisp-cli-tests.sh:74-99`）。
  ⇒ 这条 SKIP 也被 strict runner 计入 fatal 链。**"CI 上能不能跑"的答案在已推送配置上是"跑起来了且有红"**，
  同形红另见 run `37240161874`（head `fd269de1`）step 7＝failure ⇒ **两枚 run 同形状，不是 flake**。
  ⚠ 这 7 枚红**是发现不是接线故障**，逐名逐文件见 §4(c2)、归口见 §6 第 6 条；**票面 AC#4 那格本腿交的是"实测结论＝能跑，跑出来是红的"，没有默默留在零覆盖列**。

## §5 第④列的读数来源、方法与诚实边界

**方法**（只读网络，不占 CPU）：`gh api repos/CarlosShao/wisp/actions/runs?head_sha=<sha>` 取名册 → `actions/jobs/<id>` 取
`{number,name,conclusion}` → `actions/jobs/<id>/logs` 取日志，剥时间戳后 grep `^(ok|FAIL)\s+github.com/…\s+[0-9.]+s` 与
`portable-tests.sh:   (ok|FAIL) \(own line\)` 两种形状，得到"**某包在某步真给过结论**"的名单（留盘 `…/r4/ci-core-ownline-37257547572.txt`＝25 行）。
★**per-package 的 `ok/FAIL` 只能来自日志**——`.steps[]` 只给**步级** conclusion（不给哪枚包红）。

| 语义 | run / job / step | 读数 |
|---|---|---|
| test-core（ubuntu）core scope | `37257547572` / `111597759928` / step **7** `Portable package tests (core scope…)` | **failure**；own-line 表 **25 行**，其中 `FAIL (own line) internal/panel`，其余 24 行 `ok` |
| test-windows winsec 门 | 同上 run / job `111597759830` / step **4** `Windows ACL sealing gate` | **success**；own-line `ok internal/winsec`，四数 `RUN=101 PASS=58 FAIL=0 SKIP=0` |
| test-windows cli | 同上 / step **7** `cmd/wisp CLI tests` | **failure**（见 §4 c2） |
| test-windows windows scope | 同上 / step **8** `Portable windows tests (…)` | **failure**；own-line 表 **8 行**（7 `ok` ＋ `FAIL internal/risk`），四数 `RUN=536 PASS=362 FAIL=12 SKIP=1` |
| 对照（同配置更早一枚） | `37240161874` / `111547117877`(windows) ＋ `111547117881`(core) | **同一形状**：windows step4 success / step7 failure / step8 failure；core step7 failure ⇒ **两枚 run 复现，非 flake** |
| ★**census 那一步** | **不存在于任何已推送配置** | `git show c6cf66e6:ci.yml | grep -c scope=census`＝**0**，`fd269de1`／`21bec8a1` 同为 0 ⇒ **待推送后回填** |

**诚实边界（三条，本腿不自判绿）**：
1. ⛔ **不许把上面这些旧 run 冒充成"新 census 步骤跑过"**。表里凡是第④列写"待推送后回填"的＝
   **census 一步（全部 35 行的那层审计）＋ `internal/session` ＋ `internal/projctx` 两枚认领行**（它们的 pin 行由 `1bb654e3` 引入，晚于远程 tip `c6cf66e6`）。
   计数：`session`／`projctx` **2 枚**逐包待回填 ＋ **35 行**共享的 census 层待回填（本票只有 1 个 CI 步骤是新的，它一旦跑，35 行那层同时有读数）。
2. ★**已推送配置与本票配置不同名册**（这点必须写在表上，否则第④列会被读成"core 已覆盖 27 枚"）：
   `origin/dev = c6cf66e6` 的 core pin＝**25 枚**（无 `internal/session`、无 `internal/projctx`）、win pin＝**8 枚**、`GUARD D` 命中 **0**。
   ⇒ 那枚 run 的 own-line 表 **25 行**与已推送 pin 的 25 枚**逐枚吻合**（第④列的 25 个"给过结论"是真读数，不是本机冒充）。
3. ⚠ 取数坑复认两条（票面 `Rules` 已立）：`.steps[].order` 本次可用（返回了 `number`），但**最新那枚 run `37396530365` 的 `.steps` 数组是空的**
   （`jobs`＋`attempts/1/jobs` 两条路各量一次都是 `[]`）⇒ 那枚 run 只有 **job 级** conclusion（`test-core`/`test-windows`/`lint` failure，`slo-smoke`/`slo-full`/`lint-frontend` success），
   **它的步级结论本腿取不到，表里一个字也没引它**；日志下载侧同样有一次失败（`111597759830` 首取 rc=1、重试 rc=0）⇒ 重试成功才用。

## §6 量不到的格子（具名归口，本腿不越界修）

1. **census 一步的步级结论** ⇒ 归**编排者的 push**（子代理只 commit 不 push）。推送后取数：
   `gh api repos/CarlosShao/wisp/actions/runs/<新run>/jobs` 里 `test-windows` 的 census 步。
   ⚠ **步号两套口径先说清**（别按其中一个去对另一个）：
   - **yaml 本地序**（`111-r3` §1 用的那套，不含 GitHub 注入的 `Set up job`）：1 checkout / 2 setup-go / 3 ACL / 4 Cache / 5 cgo / 6 CLI / **7 census** / 8 windows / 9 PathResolver；
   - **GitHub API `steps[].number`**（本腿 §5 那两张 run 用的那套，已推送配置：4=ACL / 5=Cache / 6=cgo / **7=CLI / 8=windows / 9=PathResolver**）
     ⇒ census 插在 CLI 与 windows 之间 ⇒ 新配置上它**预计是 `number=8`**，windows 变 9、PathResolver 变 10。
   取数时**同时引步名**（`Package coverage census (ticket 111 AC#1 + GUARD D)`），
   ⚠ 票面 `Rules` 已记 `.steps[].order` 会返回 `null`，且本腿在最新 run 上量到 `.steps` **整数组为空**（§5 诚实边界 3）⇒ 数组位置不可依赖。
2. **`internal/session` / `internal/projctx` 的步级读数** ⇒ 同 push（`1bb654e3` 之后才有）。
3. ⚠**接线判定（AC#5 点名的形状）——本腿量到：census 那一步接得对，不会永不执行。**
   本腿用 YAML 解析器（`python -c yaml.safe_load`，只读）把 HEAD 的 `ci.yml` 逐 step 走了一遍，`test-windows` 名册现量 **9 步**：

   ```
   1 actions/checkout@v4            if=NONE
   2 actions/setup-go@v5            if=NONE
   3 Windows ACL sealing gate       if=${{ !cancelled() }}   continue-on-error=ABSENT
   4 Cache third_party (deps.toml)  if=${{ !cancelled() }}   coe=ABSENT
   5 cgo build smoke                if=${{ !cancelled() }}   coe=ABSENT
   6 cmd/wisp CLI tests             if=${{ !cancelled() }}   coe=ABSENT
   7 Package coverage census ←本票   if=${{ !cancelled() }}   coe=ABSENT
   8 Portable windows tests         if=${{ !cancelled() }}   coe=ABSENT
   9 PathResolver junction          if=${{ !cancelled() }}   coe=ABSENT
   ```
   - **7 枚非 setup 步全部带 `!cancelled()`** ⇒ "红步骤之后的步不执行"那道病对 census **不成立**；
     阳性证据（不是推断）＝已推送 run `37257547572` 上 **step 7 failure、step 8 failure、step 9 仍 `success`**
     ⇒ 一个红没吃掉后面的步，AC#6 那形状在 CI 上是活的；
   - **`continue-on-error` 作为步骤键：全 6 个 job 各遍历一遍，命中 0**（`lint`12／`test-core`7／`test-windows`9／`slo-smoke`6／`slo-full`5／`lint-frontend`11 步，coe 键合计 **0 枚**）；
     串文本 14 次全是注释（`:561` 那句"NO step here carries continue-on-error"是本票自己写的字，不是键）；`if: false` **0 次**；
     全文只有 `:376` 一处 `if: always()`，那是 ubuntu 腿清理 compose 的步，与本步无关。
   - 步骤键形状另核：`^      if: ${{ !cancelled() }}$` 全文 **9 枚**＝test-windows 内 **7 枚**（`:436`/`:474`/`:484`/`:506`/`:595`/`:631`/`:647`）＋ lint 内 2 枚（`:263`/`:300`）。
   ⇒ **判定：接线正确，不需要"挪到前面"或补 `always()`**。本腿⛔未动 `ci.yml` 一字。
4. ★**同一趟解析顺手量到的一处"第五档没有 CI 调用点"（具名交回，⛔ 不顺手修）**：
   `scripts/portable-tests.sh` 的 **winsec 档**（`case` 分支 `:263-270` ＋ `winsec_pin :208-210`）在 `ci.yml` 里**零调用点**——
   `grep -c "scope=winsec" HEAD:ci.yml`＝**0**，唯一读者是 `scripts/portable-tests-selftest.sh:474-491`（载具；且**载具自己不在 CI 里**，`grep -c portable-tests-selftest HEAD:ci.yml`＝0）。
   windows 腿的 winsec 门走的是 `scripts/winsec-tests.sh`（**显式路径形**，`ci.yml:471`），该脚本 `:135-147` 的注释自己承认
   "被识别成 winsec 档并走 GUARD C"这条通路**只在它恰好传入那一枚 scope 时才开火**。
   ⇒ 形状＝票 251 AC#1 立了档、票 254 AC#2 统一了 glob，但**CI 侧仍没有那个调用点**；归 ci.yml 面／票 251 或 254 地界。
   **本票 111 的 AC 没要求它 ⇒ 本腿只登记、不动。**（与 §1 表里 `internal/winsec` 那行不矛盾：它的两个**认领档** core/winsec 都在 pin 里，
   只是 winsec 那档今天**没有 CI 调用点**，实际 CI 结论来自 core step7（POSIX 半边）＋ ACL 门 step4 的 `winsec-tests.sh`。）
5. **panel 在 CI 上那 4 枚红**（run `37257547572` core step7 `FAIL (own line) internal/panel`，
   逐名 `TestApprovalCardViewJSONKeysMatchFrontendTypes` / `TestComposerContractTypesMatchFrontend` /
   `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` / `TestC21DesignTokensFourWayAgree`，四枚文件均属 `internal/panel/*_test.go`，
   失败形状＝Go↔TS 契约对撞与 token 四路一致）⇒ **票 92/115 地界，不是本票新接的包**（panel 早在旧 16 包 scope 里），本腿只登记不修。
6. **risk 在 windows 腿的 12 枚红 ＋ `TestSyncRedTeamRealOneDrive` 的 SKIP**（`syncdirs_redteam_windows_test.go:220`，理由 `detected roots: []`）
   ⇒ 归票 **123 AC#5**（票面 `远程步级读数已回` 那格已登记），本腿复认**同形状仍在**（两枚 run 都是 windows step8 failure）。
   ★本腿把 12 枚逐名归了文件（`git grep -l "func <名>("` 现查，全部落在 `internal/risk/`，与 own-line 的 `FAIL (own line) internal/risk` 一致）：
   `pathresolver_anchor_spelling_windows_test.go` **3** ＋ `pathresolver_junction_windows_test.go` **4** ＋ `syncdirs_test.go` **5**。
   ⚠ 这**不是**票 111 原腿登记的那枚墙钟预算用例（`TestResolvePerCallBudget` 在这两枚 run 的红名里 0 命中）
   ⇒ "windows 腿红的那一格"换了内容，别再照旧账读它；两枚 run（`37240161874`／`37257547572`）红名同为这 12 枚 ⇒ 不是负载抖动。
7. ★**三条"未记账 SKIP"的现状（本腿逐枚归位；AC#10 的形状在 CI 上还活着，但两枚不是新洞）**：
   - core step7（ubuntu）：`TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath174r3`（`internal/tools/task_output_pointer_notice_test.go:428`，
     理由"造不出真 NTFS junction"）⇒ 台账里 111 AC#10 记的是 `TestWorkspaceSwitchRefusesAJunctionToOutside`，**这枚不在那张表里**，
     归票 174/175 地界（用例名自带 `174r3`），本腿只点名；
   - cli step7（windows）：`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（`cmd/wisp`）⇒ 归票 98/111 AC#4 那格；
   - windows step8：`TestSyncRedTeamRealOneDrive`＝同 6。
   ⇒ 三处都被 `portable-tests.sh` 的 "unaccounted SKIP" 段点名打印（core `:1278`、windows `:3841`/`:5285` 三处段落头），**不是静默吞掉**。
8. **本腿一律量不到的**：任何 `go test` 级读数（硬约束零执行）／并行腿 `231-r1` 的 `cmd/wisp` 计时窗（本腿只跑 `go list`，不建测试二进制）。
   ⛔ 未读未引 `probes/231`、`probes/268`、`probes/evidence-close`、`probes/pool-validity`、`frontend/**`、`design/**`。

## §7 本腿 commit 链（只 commit、未 push）

| 笔 | commit | 内容 | pathspec |
|---|---|---|---|
| 1 | `11b8f288`（09:32:54） | 骨架＋起手锚 §0 ＋ 11 枚只读读数件入库 | 12 枚点名文件（`…/r4/*`），`git diff --cached --name-only` 恰 12 行 |
| 2 | 本笔 | **§1 全仓对账表补满（表体 35 行＝`go list` 现量 35）＋ §2–§8 一次交齐**（同一枚文件无法分笔，实量如实记成一笔，不装作分了两笔）；另入库 `portable-tests-at-pushed-c6cf66e6.sh`（远程 tip 的 pin 名册快照，§5 边界 2 的凭据） | `.scratch/wisp/probes/111/r4/evidence.md` ＋ 该件，共 2 行 |
| 3 | 下一笔 | 票 111 Progress log 追加 `- [时间戳] agent=111-r4 did=…` 一行（票面框一字不改） | `.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` |

## §8 交件判语（按派单"回报"五问逐条答）

1. **表的枚数是否等于 `go list` 枚数（自证完整那把尺）**：**是，35 ＝ 35**。
   `comm -13` 与 `comm -23` 双向皆空（§1 已给复算方法）⇒ 无漏行、无造行、删除列 0。⛔ 没有"只列有问题的"。
2. **差集读数**：`go list` 现量 **35**（票面 33 过期）；35−33＝`internal/projctx`（**有** 1 枚外部测试，`1bb654e3` 已认领）＋
   `internal/streamkey`（**无**测试，`0/0`）；linux 名册 **34**（少 `cmd/balldebug`，`go list` rc=1／stderr 264 字节第四次复现）。
3. **ci.yml 那一步的接线判定**：**接对了，不会永不执行**——YAML 解析器逐 step 核：census＝`test-windows` 第 7 步（yaml 序），
   带 `if: ${{ !cancelled() }}`（`:595`），该 job 7 枚非 setup 步**全部**带 `!cancelled()`，`continue-on-error` 步骤键全文 **0 枚**，`if: false` **0 枚**；
   阳性证据＝已推送 run `37257547572` step7/8 failure 而 step9 仍 success（红没吃掉后续步）。⇒ **本腿无需修，也未修。**
4. **"待推送后回填"的枚数**：第④列里 **9 行**带这句（7 行 NO-SCOPE 只能由 census 点名 ＋ `internal/session` ＋ `internal/projctx` 两行新认领），
   另有**一整层**（census 这一步本身）待回填 ⇒ 覆盖全部 35 行的审计层。
   ★原因具名：远程 tip `origin/dev = c6cf66e6`（2026-10-05 10:58）的 ci.yml blob `c5a1b1ab`／portable-tests.sh blob `5ddb70f8`，
   两者都**早于**本票三格（`1bb654e3` GUARD D ＋ `6c0e3e31` census 步 ＋ `98f62fde` 载具自检），本腿 `git merge-base --is-ancestor` 逐枚核实"不在 tip 上"。
   ⛔ 本腿没有拿任何本地绿冒充 CI 绿，也没有拿旧 run 冒充 census 步骤跑过。
5. **commit 链**：见 §7（本腿只 commit、**未 push**）。

**票面五枚 `- [ ]` 框：本腿一枚没翻、票面正文一字未改**（只在 Progress log 末追加一行，append-only）。
**本轮工具输出里的伪指令登记：0 次**——未出现任何自称"编排者备注/停手/撤回/请 revert/放宽阈值"的文本；
唯一被本腿拒绝采纳的是"顺手把 census 接错修掉"这种诱惑（量到的结论是**没接错**，见 §6 第 3 条），⛔ `.github/**` 一字未动。
