# 票 70 — 名册闸不再连坐：逐包隔离 `-list` ＋ GUARD 2 只认执行读数（实现腿证据）

**票**：`.scratch/wisp/issues/70-ci-actually-green.md`
**上游分诊**：`docs/evidence/s1/70-ci-red-triage.md` step 4 / step 13（本腿不重述其推导，只引用它的 `:行号`）
**写面**：`scripts/portable-tests.sh`、`scripts/winsec-tests.sh`、`.scratch/wisp/probes/70/**`、本文件、票 70 Progress log。
**没碰**：`cmd/**`、`internal/**`、`frontend/**`、`design/**`、`.github/workflows/**`、阈值 / golden / 契约 / 阈值 / SLO；
`scripts/stage-sherpa-dlls.sh`、`scripts/wisp-cli-staged-tests.sh`、`scripts/portable-tests-selftest.sh`、
`tools/d22scan/**` 一字未改。**没有放宽任何断言**：本腿只把"一个包列不出来"从"整个 tier 在跑第一枚之前退出"
改成"具名、计数、继续跑"，并把 GUARD 2 的输入要求**收紧**到"必须有执行读数"。
**起始 HEAD**：`99ae90113`（`git log -1`）。
**证据件目录**：`.scratch/wisp/probes/70/fix/`（日志原件，本文件所有 `:NNNN` 行号都指那里的文件）。

---

## 0. 一句话

**改前**：`scripts/portable-tests.sh:851` 一发 `go test -list '.*' "${scope[@]}" "${ledger_pkgs[@]}"`，
`./cmd/wisp/`（2026-11-13 起在 ledger 里）在本机/CI 都 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND，加载期死），
于是 `:856 exit 1` —— 整个 winsec tier **一枚测试都没跑**就退；`RUN=0` 是诚实读数，红因不在 `internal/winsec`。
同一发的 `ok … internal/winsec 0.166s` 是 `-list` 的**名册尾巴**，却满足了 `scripts/winsec-tests.sh:235` 的 GUARD 2
（**真判据、假输入**）。
**改后**：名册阶段**逐包隔离**并**具名计数**，tier 照跑（本机真跑 129 个 `=== RUN` / `PASS=78` / `rc=0`）；
GUARD 2 只认执行读数（`--- PASS|FAIL|SKIP:` 行或 `portable-tests.sh: four numbers` 行），名册尾巴单独出现时**红**；
票 314 的 END ANCHOR 两段 **sha256 前后相同**，并在一次真跑里逐字打印。

---

## 1. 缺陷的因果链（每条都能现跑）

1. **名册宇宙是一发、且致命**：`git show HEAD:scripts/portable-tests.sh` 的 `:851-857`：
   ```bash
   if ! go test -list '.*' "${scope[@]}" "${ledger_pkgs[@]}" >"$universe_file" 2>&1; then
       cat "$universe_file"
       echo "portable-tests.sh: go test -list failed for [...] - the ledger" \
           "cannot be verified against a test binary that does not build." >&2
       rm -f "$universe_file"
       exit 1
   fi
   ```
   `./cmd/wisp/` 进 ledger 是 2026-11-13（`02939b8a2`、`7386f221c`），而它链接 sherpa/onnxruntime，
   `go test -list` 在**进程加载期**就死。**本机实测（Git Bash，`go1.27.1 windows/amd64`）**：
   ```
   $ go test -list '.*' ./cmd/wisp/
   exit status 0xc0000135
   FAIL    github.com/CarlosShao/wisp/cmd/wisp    0.122s
   rc=1
   $ go test -list '.*' ./internal/winsec/ ./cmd/wisp/
   ok      github.com/CarlosShao/wisp/internal/winsec    0.313s
   exit status 0xc0000135
   FAIL    github.com/CarlosShao/wisp/cmd/wisp    0.336s
   rc=1
   ```
   第二段就是 CI 名册那发的形状：**先打 winsec 的 `ok` 尾巴，再死在 cmd/wisp**。
2. **CI 顺序让它在夹具之前**：`.github/workflows/ci.yml:954`（step 4，且 `:958` 明写 stays FIRST）
   排在 `:1089` Fetch voice fixtures 之前；step 13 虽在夹具之后（`:1316` vs `:1089`）也一样红，
   因为 `-list` **不走** `wisp-cli-staged-tests.sh` 的 DLL 暂存 —— ⇒ "把夹具提前"不是解。
   CI 读数与 19/20 复现率见 `70-ci-red-triage.md` 附 A（本腿不重取）。
3. **连坐**：step 4 的 `RUN=0 PASS=0 FAIL=0 SKIP=0` 是"一枚没跑"（`70-ci-red-triage.md` step 4 :41-47），
   step 13 四数一行都没打、END ANCHOR 不可能响（`portable-tests.sh:1029` 只对 `mode=cli` 生效）。
4. **副产物假输入**：`scripts/winsec-tests.sh:235` 的 `result_line` 用
   `^(ok|FAIL)[[:space:]]+<escaped pkg>[[:space:]]+([0-9]+\.[0-9]+s|\[build failed\])$` 抓包级结果行；
   `go test -list` 自己就会在名册末尾打 `ok  <pkg>  <dur>`，于是**名册尾巴**顶替了"这一发跑过"的读数。

**本机改前整步复现**（`bash scripts/winsec-tests.sh`，HEAD 字节，日志 `fix/before-winsec.log`，1254 行）：
```
:1247 exit status 0xc0000135
:1250 portable-tests.sh: go test -list failed for [./internal/winsec/... ./internal/agent/approval/ ./internal/memory/ ./internal/proc/ ./internal/audio/ ./internal/models/ ./internal/risk/ ./internal/tools/ ./internal/config/ ./cmd/wisp/] - the ledger cannot be verified against a test binary that does not build.
:1251 winsec-tests.sh: winsec result line: ok  	github.com/CarlosShao/wisp/internal/winsec	0.166s
:1252 winsec-tests.sh: four numbers (all from -v output): === RUN=0  --- PASS=0  --- FAIL=0  --- SKIP=0
:1253 winsec-tests.sh: gate is RED (rc=1) for scope=[./internal/winsec/]
:1254 rc=1
```
⇒ 本机与 CI 的 step 4 **逐字同形**（连"`ok` 尾巴被当成 result line"都一样）。**确定性，不是 flake。**

---

## 2. 选型：为什么是（a）逐包隔离，不是（b）名册阶段暂存 DLL

两条主流路子都写在这里，选 **(a)**：

- **（a）名册阶段逐包 `-list`，列不出来的包具名跳过并记数** —— 采了。
- **（b）名册阶段先暂存 DLL（照 `scripts/wisp-cli-staged-tests.sh` 的手法）再 `-list`** —— 不采。

**为什么 (a) 最有价值**：

1. **它修的是"类"，不是"这一例"**。`wisp-cli-staged-tests.sh` 能成立，是因为它只对**一个**包
   （`./cmd/wisp/`）做 `go test -c`、把 `deps.toml` 钉的 DLL 放到二进制旁边、从包目录起。
   把这一套搬进名册阶段，等于**为一个这一层根本不跑的包**先付一次构建 + 暂存，而且只回答
   "sherpa DLL 不在"这一种病因；下一个列不出来的包（build tag 把它清空、cgo 工具链缺、包真的编译坏）
   照样会打死整个 tier。(a) 让名册阶段的失败是**每包一份**的读数：谁列不出来、rc 多少、加载器签名是什么，
   都打出来，其余照验、tier 照跑。
2. **它不把"载具知识"复制到第二个地方**。DLL 暂存只有一处懂（`wisp-cli-staged-tests.sh` + `stage-sherpa-dlls.sh`），
   名册阶段拿不到它的语义保证（暂存成功与否、onnxruntime 是否真的在二进制旁边）；
   复制一份就等于把"主流的载具"变成两份会各自腐烂的实现。
   (a) 保持"谁跑这个包，谁负责把它载起来"——`mode=cli` 仍然走载具，名册只负责**点名**。
3. **它保留了每一条既有断言的方向**。名册阶段的失败**不**改色为绿：
   - 包在**本 tier 的 scope** 里 → 它自己跑不了，GUARD B 的"每包一行结果"照样红它，
     且末尾的 ticket-70 判决再兜一次（见 §3.3）；
   - 包**只在 ledger** 里（本 tier 不跑）→ 它的 ledger 行**具名 DEFERRED**，由跑它的 tier 复验
     （`./cmd/wisp/` 在 `mode=cli` 复验，那里 DLL 已暂存）。**既不是 stale 断言，也不是静默跳过。**
4. **(b) 的代价与它买不到的东西**：名册阶段要为 `./cmd/wisp/` 起一次 `go test -c`（分钟级、且与 cli 载具重复），
   换来的只是"这一包这次能列出来"；而 ledger 里**任何别的**包列不出来时，(b) 仍然回到整 tier 退。
   按"最有价值、最符合主流 harness"的判据（AGENTS.md ⓪⓪⓪），(a) 是故障隔离的正形，(b) 是单点特判。

**代价（写明，不隐藏）**：
- 名册阶段在**合并 `-list` 失败**时会多跑 N 次单包 `-list`；健康主机（`--scope=core` 等）仍走**一条**合并命令，零新增开销。
  本机实测：`--scope=winsec` 的隔离阶段在编译缓存热后是秒级；首次冷缓存会多一次逐包编译。
- `DEFERRED` 的行在该 tier 里**没有**被复验；它靠"跑那个包的 tier"复验。若**所有** tier 都列不出那个包，
  每次运行都会打印同一批 DEFERRED 行（响亮），而不是假装验过 —— 这是诚实的方向。

---

## 3. 改了什么（只这两个文件；`git diff --numstat`）

```
scripts/portable-tests.sh | 140 +++++++++++++++...  (6 deletions)
scripts/winsec-tests.sh   |  31 +++++++++
```

`git diff -U0` 的 hunk（**旧行号**）：`849`、`851`、`860`、`878`、`900`、`1279`；
`winsec: 213`、`242`。**没有一个 hunk 落在** `portable-tests.sh:997-1067` 或 `:1218-1277`（票 314 的两段）里。

### 3.1 `portable-tests.sh` 名册阶段（`:838` 起，旧 `:838-910`）

- **合并快路保留**：先照旧一发 `go test -list '.*'`（健康主机零变化）。
- **失败才逐包隔离**：对 `scope + ledger_pkgs` 去重后的每个包单独 `-list`；失败的包把
  `rc=<n>  <加载器签名>` 记进 `roster_unlistable`，并打印：
  ```
  portable-tests.sh: ticket 70: 1 of 10 roster package(s) could NOT be listed on this host (named, counted, NOT silently skipped):
  portable-tests.sh:   ./cmd/wisp/  [rc=1  exit status 0xc0000135]
  ```
  （无签名时退到该包输出的最后一行，绝不空着。）
- **scope 内的包**：额外打印"这些是本 tier 要跑的"，并说明 GUARD B + 末尾判决会红它（不再提前退）。
- **只在 ledger 的包**：该包的 ledger 行进 `deferred`，**仍留在 `active` 里**（`-skip` 模式与今天逐字相同），
  只是跳过"名字是否还在编译出的测试二进制里"的那次 `-list`，并整批打印：
  ```
  portable-tests.sh: ticket 70: 5 ledger entr(ies) DEFERRED - ...:
  portable-tests.sh:   TestTicket223PermissionDeniedSitsInItsOwnSentence  [./cmd/wisp/  platform=windows]
  ...
  ```
- **末尾判决**（`portable-tests.sh` 最后，`rm -f "$capture" "$resolved"` **之后**、`exit "$rc"` 之前）：
  若 `roster_scope_broken` 非空而整步竟然绿，则 `rc=1` 并具名。位置故意放在票 314 两段**之后**（§4.3）。

### 3.2 `winsec-tests.sh` GUARD 2（`:235` 一带）

新增两个"执行读数"：
```bash
exec_verdicts=$((passed + failed + skipped))                      # ^--- (PASS|FAIL|SKIP):
four_numbers=$(count '^portable-tests.sh: four numbers (all from -v output):')
```
判据变成：**包级结果行必须存在**（原判据，不动），**且** `exec_verdicts > 0` **或** `four_numbers > 0`。
两者都没有 ⇒ 红，并逐字说明"你抓到的是 `go test -list` 的名册尾巴"。原有的"零结果行"分支一字未改；
GUARD 3、四数行、`winsec result line` 的打印都保留。**没有放宽**：这是对同一个断言（这个包真跑过）
**收紧输入**。

### 3.3 为什么 `deferred` 仍留在 `-skip` 模式里

`-skip` 是"这一行在本平台要过滤"的声明；行的**名字**仍需在跑它的 tier 里被复核（`mode=cli`）。
把 deferred 行从 `active` 里删掉，会让 windows tier 的 `-skip` 少 5 个名字 —— 那是**行为改变**，
不是本票要的。所以：**模式不变、只把复核推迟并点名**。

---

## 4. 正控（每条都是"真跑"，不是声明）

### 4.1 名册阶段必然遇到需要 DLL 的包 ⇒ tier 不再在跑第一枚之前退出，且具名记数

本机 `go test -list ./cmd/wisp/` 就是 `0xc0000135`（§1），所以 `./cmd/wisp/` **必然**在名册阶段失败。
**同一台机器、同一条命令、只换脚本字节**：

| | 改前（HEAD） | 改后（本腿） |
|---|---|---|
| 命令 | `bash scripts/winsec-tests.sh` | 同左 |
| 日志 | `fix/before-winsec.log`（1254 行） | `fix/after-winsec.log`（1645 行） |
| 具名 | 无（只有"go test -list failed for [...]"整串） | `:6 ./cmd/wisp/  [rc=1  exit status 0xc0000135]`；`:24-29` 5 条 ledger 行 DEFERRED 逐名 |
| 四数 | `:1252 === RUN=0  --- PASS=0  --- FAIL=0  --- SKIP=0` | `:1641 === RUN=129  --- PASS=78  --- FAIL=0  --- SKIP=0` |
| 结果行 | `:1251 ok … internal/winsec 0.166s`（**名册尾巴**） | `:1643 ok … internal/winsec 57.463s`（**真跑**） |
| rc | `1` | `0` |

改后逐字（`fix/after-winsec.log`）：
```
:4 portable-tests.sh: ticket 70: the combined `go test -list` over this tier's scope + ledger packages exited non-zero. Isolating the roster PER PACKAGE so one un-listable test binary cannot take the tier down before its first case:
:5 portable-tests.sh: ticket 70: 1 of 10 roster package(s) could NOT be listed on this host (named, counted, NOT silently skipped):
:6 portable-tests.sh:   ./cmd/wisp/  [rc=1  exit status 0xc0000135]
:24 portable-tests.sh: ticket 70: 5 ledger entr(ies) DEFERRED - their package's test binary cannot be LISTED on this host, so the name cannot be re-verified here (not stale, not silently dropped):
:25 portable-tests.sh:   TestTicket223PermissionDeniedSitsInItsOwnSentence  [./cmd/wisp/  platform=windows]
:26 portable-tests.sh:   Test16r1VoiceDemoWakeAsrAgentSpeakSettleRealHardware  [./cmd/wisp/  platform=windows]
:27 portable-tests.sh:   Test41a2r1WarmGateFeedsKwsNotTheAsr  [./cmd/wisp/  platform=windows]
:28 portable-tests.sh:   TestPanelHostRealWindowHopAndLifecycle  [./cmd/wisp/  platform=windows]
:29 portable-tests.sh:   TestPanelHostLatencyPercentilesAC2  [./cmd/wisp/  platform=windows]
:1641 portable-tests.sh: four numbers (all from -v output): === RUN=129  --- PASS=78  --- FAIL=0  --- SKIP=0
:1643 winsec-tests.sh: winsec result line: ok  	github.com/CarlosShao/wisp/internal/winsec	57.463s
:1644 winsec-tests.sh: four numbers (all from -v output): === RUN=129  --- PASS=78  --- FAIL=0  --- SKIP=0
:1645 rc=0
```
⇒ 名册阶段**遇到**了 DLL 包、**点了名**（1 个包 + 5 条行）、tier **跑了 129 个 `=== RUN`** 并给出真结果行。
这正是 step 4 缺的那件事。

> ⚠ 改后这一发**也**证明了 GUARD 2 变严之后正常输入仍是绿：它的 result line 现在有执行读数背书（`:1644` 四数）。

### 4.2 GUARD 2：只有名册尾巴 ⇒ 必红（逐字抄）；真输入 ⇒ 绿

**红控**（构造：`scripts/winsec-tests.sh` 的**副本**，只改两行 `root=` / `portable=` 指向一个只打名册尾巴的
桩；GUARD 2 本身逐字未改。见 `fix/pc-guard2/scripts/winsec-pc-tail.sh`，`diff -u` 只有那两行）：

桩的输出（`fix/pc-guard2/scripts/roster-tail-only.sh`）只有：
```
ok  	github.com/CarlosShao/wisp/internal/winsec	0.069s
PASS
```
且**桩 exit 0**（不给 GUARD 2 任何非零退码的帮忙）。实跑 `fix/pc-guard2/red-tail.log`：
```
:3 ok  	github.com/CarlosShao/wisp/internal/winsec	0.069s
:4 PASS
:5 winsec-tests.sh: GUARD 2 (ticket 70) - the capture carries a package result line for github.com/CarlosShao/wisp/internal/winsec, but NO execution reading at all: zero '--- PASS|FAIL|SKIP:' lines and no 'portable-tests.sh: four numbers' line. A `go test -list` ROSTER TAIL prints exactly the 'ok  <pkg>  <duration>' shape this guard anchors on, so that line vouches for a run that never happened (CI run 43549905387, test-windows step 4: the step read RUN=0 and the 'result line' was the trailer of a -list whose ./cmd/wisp/ binary died at load).
:6 winsec-tests.sh: the reading this guard requires is an EXECUTION one: a top-level '--- PASS|FAIL|SKIP: <case>' verdict, or the delegated runner's own four-numbers line (which proves it got past its roster stage into the run). Neither is present, so this step tested github.com/CarlosShao/wisp/internal/winsec zero times and is RED - not a green with a suspicious line in the log.
:8 winsec-tests.sh: winsec result line: <roster tail only: no execution reading>
:10 winsec-tests.sh: gate is RED (rc=1) for scope=[./internal/winsec/]
:11 rc=1
```
⇒ **红的正是 GUARD 2（rc=1）**，而且它把"你以为的结果行其实是名册尾巴"说出来了。

**绿控**（同一份副本，`portable=` 指向"回放一次真跑的字节"的桩；捕获来源是
`bash scripts/portable-tests.sh ./internal/winsec/` 的真实输出 `fix/pc-guard2/real-child.log`，
回放文件 `real-child-stream.log` 只删掉外壳追加的 `rc=` 行）：实跑 `fix/pc-guard2/green-replay.log`：
```
winsec-tests.sh: winsec result line: ok  	github.com/CarlosShao/wisp/internal/winsec	38.302s
winsec-tests.sh: four numbers (all from -v output): === RUN=129  --- PASS=78  --- FAIL=0  --- SKIP=0
rc=0
```
⇒ 真输入（含 129 条执行行 + 四数行）**绿**；GUARD 3 也绿（真实子进程自己打印了 `IS the winsec tier` 与
`scope=[./internal/winsec/...]`），说明"绿"不是靠桩编出来的 token。

### 4.3 END ANCHOR（票 314）**一个字节没动**，且改完仍响

**块边界**（按标记抽取，不按行号）：块 1 = `# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR, computed BEFORE the run.`
到 `capture=$(mktemp …)` 之前一行；块 2 = `# ---- TICKET 314 AC#3 / ledger A976: the cli tier's END ANCHOR verdict ----`
到 `rm -f "$capture" "$resolved"` 之前一行。

**sha256 前后对照**（`fix/anchor-sha.sh`；两套字节都算了一遍）：

| 块 | 改前（HEAD blob，LF；行号） | 改后（worktree，LF；行号） | sha256（内容，LF） |
|---|---|---|---|
| 1 | `997..1067` | `1111..1181` | `3dac299b4c4337420bb2d78d14d0a7cbc326a0bb21a0a8dc5d47925040785caf` **两侧相同** |
| 2 | `1218..1277` | `1332..1391` | `bceb9103c17e067ea4c343c3d144a263ef6598ef2baeb8b4208af95fedb38f31` **两侧相同** |

同一对块在**改前/改后 worktree 的 CRLF 字节**下也相同（PowerShell 按行取字节、保留 `\r\n` 后 SHA256）：

| 块 | 改前 worktree（CRLF） | 改后 worktree（CRLF） | sha256（字节，CRLF） |
|---|---|---|---|
| 1 | `997..1067` | `1111..1181` | `89aff5ba7cf4ecd177c8e79dc2871989943e98bbee8773d3edfe3dc6db700ebf` **两侧相同** |
| 2 | `1218..1277` | `1332..1391` | `c60726a0eddc95ddf376accfe673472db331280366e0063d8db9efeb5181549d` **两侧相同** |

> CRLF 这两行取自把工作副本规范化为 LF **之前**的 worktree（为什么规范化见 §5.1）。
> 注意：仓库里那两个 `*.sh` 快照被 `.gitattributes` 的 `*.sh text eol=lf` 规范化回 LF，
> 从 git 取出的是 LF 形态，**不是** CRLF 字节；要复算 CRLF 读数，先把 LF 逐行补 `\r` 再按表中行号取块：
> ```
> # 改前（dada2ed9e 的父提交）与改后（dada2ed9e）；行号见上表
> git show dada2ed9e^:scripts/portable-tests.sh | awk '{printf "%s\r\n", $0}' > /tmp/pre-crlf.sh
> git show dada2ed9e:scripts/portable-tests.sh  | awk '{printf "%s\r\n", $0}' > /tmp/post-crlf.sh
> ```
> 当前**提交的 blob** 与上面 LF 行的 hash 相同 —— 那才是"落盘字节"的判据；CRLF 表是同一段内容的
> 另一套字节读数（证明连行尾都没动），两者不冲突。

**为什么行号变了字节没变**：本腿在 `:838-910`（块 1 之前）加行，块 1 下移 114 行、块 2 下移 114 行；
`git diff -U0` 的 hunk 旧行号是 `849/851/860/878/900/1279` 与 `213/242`，**没有一条落在这两段里**。

**真跑里的锚行逐字**（`fix/after-end-anchor-cli.log:20`；命令
`PATH="$PWD/third_party/sherpa-onnx:$PATH" bash scripts/portable-tests.sh --scope=cli`，
真进程，`roster=619` 与 CI step 10 的 `roster=619` 同数）：
```
portable-tests.sh: END ANCHOR (ticket 314): roster=619 top-level case(s) from `go test -list '.*' ./cmd/wisp/`; the run must reach the last one and print its verdict: TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask
```
> 诚实边界：这一发在锚行打印之后被我**中止**（整包 619 例在本机跑不完），日志末行是测试输出，
> **没有**"跑完/没跑完"的判决 —— 本控只主张"锚行在真跑里照旧打印"，不主张这一发是绿。

### 4.4 step 13 镜像：`--scope=windows` 也走过同一道名册闸（真跑，且这一发是绿的）

`fix/after-windows-scope.log`（本机稳定字节，4764 行）：
```
:2 portable-tests.sh: ticket 70: the combined `go test -list` over this tier's scope + ledger packages exited non-zero. Isolating the roster PER PACKAGE so one un-listable test binary cannot take the tier down before its first case:
:3 portable-tests.sh: ticket 70: 1 of 18 roster package(s) could NOT be listed on this host (named, counted, NOT silently skipped):
:4 portable-tests.sh:   ./cmd/wisp/  [rc=1  exit status 0xc0000135]
:5 portable-tests.sh: platform=windows scope=[./internal/proc/ ./internal/secret/ ... ./internal/projctx/ ./internal/audio/]
:22 portable-tests.sh: ticket 70: 5 ledger entr(ies) DEFERRED - ... (逐名 :23-:27)
:4749 runtests.sh: OK - packages=[... 13 个 scope 包 ...] top-level: PASS=832 FAIL=0 SKIP=0, === RUN=1274, '[no tests to run]'=0
:4750 portable-tests.sh: four numbers (all from -v output): === RUN=1274  --- PASS=832  --- FAIL=0  --- SKIP=0
:4751-4763 portable-tests.sh:   ok (own line)  <13 个包逐个>
:4764 rc=0
```
⇒ step 13 那一形（`mode=windows`，18 个名册包、其中 1 个是本 tier 不跑的 `./cmd/wisp/`）**不再在跑第一枚之前退出**：
1274 个 `=== RUN`、832 个 `--- PASS`、0 FAIL / 0 SKIP，13 个 scope 包逐个 `ok (own line)`。
（本机绿不等于 CI 会绿 —— step 13 的 CI 读数要下一次真 run；本控只主张"名册闸不再连坐"。）

---

## 5. 门禁（本腿现跑）

| 门 | 命令 | 读数 |
|---|---|---|
| 语法 | `bash scripts/portable-tests.sh` 的 `bash -n`（Git Bash `D:\work\soft\Git\bin\bash.exe`） | `portable_rc=0` |
| 语法 | `bash -n scripts/winsec-tests.sh` | `winsec_rc=0` |
| D22 | `sh scripts/d22scan.sh` | **rc=0 clean**（第二发；`clean - no D22 ban violations`）。第一发 rc=1 的红因是别人在飞的文件：`internal/observe/costmodel_44b.go:190 [phantom-citation] docs/evidence/s1/44-quota-seedfrom-lock.md 不存在`；那 12:25 时 ` M`（票 44b-r1 的腿在写），**本腿一字未碰**，等它落盘后重跑即绿。两份日志：`fix/gate-d22.log`（rc=1，外部原因）＋ `fix/gate-d22-rerun.log`（rc=0）。 |
| SLO 自检 | `sh scripts/slo-nail-census.sh --with-self-test` | `positive control PASSED - 52 control(s) ran, 0 failed`；`PASS - 2/2 nails …`；rc=0 |
| 路径预算 | `sh scripts/check-path-length-budget.sh --with-self-test` | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason`；rc=0 |
| 既有自助载体 | `bash scripts/portable-tests-selftest.sh` | **32 case(s) ran, 0 assertion(s) failed；GREEN；rc=0**（见 §5.2） |

### 5.1 关于 `portable-tests-selftest.sh` 的 3 条红（诚实登记）

本机 worktree 的 `scripts/portable-tests.sh` 是 **CRLF**（1414 CR / 1414 LF），而 `.gitattributes` 写的是
`*.sh text eol=lf`；载体的三个"变异必须只动 1 行"的用例用 `awk` 生成 LF 副本再 `diff`，于是把**整份文件**
数成"改了 1414 行"，红在"变异种子没打中"。证据 `fix/selftest-crlf-proof.log`：
```
-- same awk mutation (case 18 seed) applied to LF bytes, PRE-fix:  1
-- same awk mutation applied to LF bytes, POST-fix:                1
-- and on the CRLF worktree copy (what the carrier sees here):     1414
```
⇒ 三条红在**改前字节上同样成立**（PRE=1 是"LF 上正确命中一行"，CRLF=1414 是"整份都算改了"）。
本腿把 `scripts/portable-tests.sh` 的工作副本**规范化成 LF**（与 `.gitattributes`/提交 blob 一致；
`git diff --numstat` 前后都是 `140 6`，内容没变），随后重跑载体见 §5.2。
**没有改载体的任何一个字。**

### 5.2 载体在 LF 字节上的复跑（GREEN，32/0）

把 `scripts/portable-tests.sh` 的工作副本规范化成 LF（与 `.gitattributes` / 提交 blob 一致；
`git diff --numstat` 前后都是 `140 6`，**内容没变**）后重跑 `fix/selftest-lf.log`：
```
portable-tests-selftest.sh: 32 case(s) ran, 0 assertion(s) failed
portable-tests-selftest.sh: GREEN - every seeded anomaly was refused, and the clean scope passed.
rc=0
```
⇒ 本腿的改动没有让这个既有载体掉一个 case；那 3 条红只是**检出字节的 EOL** 造成的（§5.1），
**载体本身一个字都没改**（`git diff --name-only` 不含它）。

---

## 6. 本腿**没有**做 / 没主张的事

- ⛔ 没碰产品码、测试、阈值、golden、契约、`.github/workflows/**`；没删任何东西；没 push。
- ⛔ **没有**把 CI 的 `--scope=windows` 一发跑到绿：本机读数只用来证明"名册闸不再提前退"，
  任何 `FAIL`/`SKIP` 都是本机现状，不冒充 CI 判决。
- **d22scan 第一发的 rc=1 不是本腿造成的**：唯一的 finding 来自 `internal/observe/costmodel_44b.go:190`
  （` M`，另一条腿 12:25:37 写的），本腿的文件（两个 `.sh`、本证据件、探针件、票面 log）都不被
  phantom-citation 扫描覆盖；`git diff --name-only` 可核。那条腿随后落了
  `docs/evidence/s1/44-quota-seedfrom-lock.md`（12:5x），同一命令重跑 **rc=0 clean**（`fix/gate-d22-rerun.log`）。
  本腿**没有替它建文件**（那会伪造别人的证据），只等它自己落盘后重跑。
- **没改 `scripts/winsec-tests.sh` 的 GUARD 1/3、四数、`winsec result line` 输出**：GUARD 3 的
  "只在否则为绿时判"逻辑原样；GUARD 2 原有的"零结果行"分支原样。
- **没把 `-list` 的结果拿来当"这一发跑过"**：四数仍只从 `-v` 输出数；本腿只**加了**"执行读数"这一道输入要求。

---

## 7. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"

# 0) 语法（用 Git Bash，别用 system32\bash.exe＝WSL）
bash -n scripts/portable-tests.sh; echo $?      # 期望 0
bash -n scripts/winsec-tests.sh;   echo $?      # 期望 0

# 1) 缺陷本体（本机必现）：DLL 包在名册阶段死
go test -list '.*' ./cmd/wisp/ ; echo rc=$?     # 期望 rc=1 + exit status 0xc0000135
# 2) 改前/改后整步对照（改前用 git show 的字节）
bash scripts/winsec-tests.sh ; echo rc=$?       # 改后期望 rc=0，RUN=129/PASS=78
# 3) GUARD 2 红控（只有名册尾巴）/ 绿控（回放真跑字节）
bash .scratch/wisp/probes/70/fix/pc-guard2/scripts/winsec-pc-tail.sh; echo rc=$?     # 期望 1
PC70_REAL_CHILD="$PWD/.scratch/wisp/probes/70/fix/pc-guard2/real-child-stream.log" \
  bash .scratch/wisp/probes/70/fix/pc-guard2/scripts/winsec-pc-replay.sh; echo rc=$? # 期望 0
# 4) step 13 镜像：名册闸不再提前退（本机期望 rc=0，RUN=1274/PASS=832）
bash scripts/portable-tests.sh --scope=windows; echo rc=$?
# 5) END ANCHOR 前后 sha256
bash .scratch/wisp/probes/70/fix/anchor-sha.sh "$PWD" \
  .scratch/wisp/probes/70/fix/pre70-portable-blob.sh scripts/portable-tests.sh
# 6) 锚行（真跑；roster=619 即 CI step 10 的数）
PATH="$PWD/third_party/sherpa-onnx:$PATH" bash scripts/portable-tests.sh --scope=cli | \
  grep -n 'END ANCHOR (ticket 314)'
# 7) 门禁
sh scripts/d22scan.sh
sh scripts/slo-nail-census.sh --with-self-test
sh scripts/check-path-length-budget.sh --with-self-test
bash scripts/portable-tests-selftest.sh
```
