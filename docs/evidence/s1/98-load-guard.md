# 98 — AC#3「缺 DLL 必须点名报红」的守卫（实现腿的证据；裁决表仍归验收方）

**票**：`.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md`
**本件主张的 AC**：**AC#3**。这是今天唯一一条**完全落在 `scripts/**` ＋ `docs/**` 里**能做完的格子：
AC#1(iii) 要一次真 CI run（真机/CI 面，本腿无权也无从取 run id）、AC#4 是台账收敛（另做）、
AC#2 与 AC#3 的"能跑起来"那一半已由票 111 交在 `scripts/wisp-cli-tests.sh` 里，剩下的正是 AC#3 的**判据**：
**把 AC#2 的注入拿掉，守卫必须红，并说清是哪条 DLL。**
**Status**：**实现证据，不是裁决表**。票面写死"裁决表 `docs/evidence/s1/98-*.md` 由验收方出" ⇒ 本件只交读数与改动面；
翻勾、结案、以及"AC#3 这一形够不够"由**非实现者**判。
**树**：`ff61b9dbc`（派单 HEAD）。开工前 `git status --porcelain | grep -v '^??'` **空**（无人在飞的 `M`）。
**时间**：2026-11-05 17:30–17:5x +08。
**写面**：`scripts/wisp-cli-tests.sh`（改）＋ 本文件（新增）。**一字未碰**：任何 Go 包、`cmd/wisp/**`、`tools/d22scan/**`、
`.github/workflows/ci.yml`（归票 85）、`thresholds.go`、golden、`internal/risk/**`、`frontend/**`、`design/**`、
`docs/PLAN.md`、`docs/specs/**`、票面（`.scratch/**` 不在本腿写面）。

## 0. 一句话

**改前**：把 AC#2 的 PATH 注入拿掉，这一发是 `rc=1` ＋ 一句 `exit status 0xc0000135` ＋ 结尾一句
"the ledger cannot be verified against a test binary that **does not build**" —— 940 行日志里**没有一个字**说
"哪条 DLL"或"这个平台没测到东西"，而且那句诊断**是错的**（它建得出来，是加载期找不到 DLL）。
**改后**：同一发在**还没有任何测试进程存在**的时点红，逐字点名
`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`（第三枚 `onnxruntime.dll` 本机 System32 里真能找到，所以**不点它**——
守卫是按盘上的事实报，不是把 `deps.toml` 的名单念一遍）。

## 1. 为什么这不是"已经有人做了"（改前实测，不是读码）

仪器：本机 Git Bash（MSYS2），`unset GOTMPDIR TMPDIR TEMP TMP`，`CGO_ENABLED` 保持默认 1（**没有**、也不许设 0）。

变体构造（**只改两行**，其余字节逐字取自已提交的树；两个变体同一手法，所以可比）：
`root=` 钉成仓根（因为副本在 `/d/tmp` 下，脚本靠 `$here/..` 自推根）＋ 删掉 `export PATH="$dll_dir:$PATH"` 这一行。

```
$ git show HEAD:scripts/wisp-cli-tests.sh > /d/tmp/98probe/pre-nopath.sh
$ sed -e 's|^root=.*$|root="/d/work/workspace/projects plans/Wisp"|' \
      -e 's|^export PATH=.*$|:|' pre-nopath.sh > pre-mutant.sh      # HEAD 字节
$ bash pre-mutant.sh ; echo rc=$?
```

**改前读数（pre-mutant，HEAD 字节）**

```
wisp-cli-tests.sh: ticket 310 AC#1: GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\cli-20261105-173707-1559 ...
wisp-cli-tests.sh: scope=./cmd/wisp/ dll_dir=/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx (pinned: onnxruntime.dll sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll )
wisp-cli-tests.sh: handing the verdict to portable-tests.sh --scope=cli
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.226s
...
portable-tests.sh: go test -list failed for [./cmd/wisp/ ./internal/agent/approval/ ... ] - the ledger cannot be verified against a test binary that does not build.
pre_rc=1        log 行数 940
```

计数（同一份日志）：`=== RUN` **0**、`--- PASS` **0**、`--- FAIL` **0** ⇒ **一条测试都没跑**；
`0xc0000135` 出现 **1** 次；"sherpa-onnx-c-api.dll" 只出现在上面那行**信息性**的 pinned 名单里，
**没有任何一行说它没能被加载**，也没有任何一行说"本平台未测"。
⇒ AC#3 的判据在这一版上是**红的但没点名**，而且给了错误归因（"does not build"），正是票 98 说的
"rc=1 的加载崩溃被后人当成环境问题不追"那一形。

## 2. 改后读数

**2.1 判据发（post-mutant：本腿字节 + 同一手法删掉注入）**

```
$ bash post-mutant.sh ; echo rc=$?
wisp-cli-tests.sh: ticket 310 AC#1: GOTMPDIR=... （310 那行照旧，未动）
wisp-cli-tests.sh: GUARD (ticket 98 AC#3) - cmd/wisp was NOT tested on this host: the loader cannot reach pinned DLL(s) on PATH: sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll
wisp-cli-tests.sh: the test binary is started from a TEMP directory, so a copy of these files in the current directory does not help: measured on this host, a stray copy in the repo root loads for a direct start of the same binary while go test still dies at load.
wisp-cli-tests.sh: without them the process dies at LOAD time (0xc0000135) before the first test runs: rc=1 with zero '--- FAIL' lines, which reads exactly like a broken package and is why ticket 98's hole was carried around as an 'environment fact' for six briefings.
wisp-cli-tests.sh: fix the staging step (deps.toml -> scripts/fetch-deps.ps1 -> .../third_party/sherpa-onnx) or point PATH at wherever the pinned files really are, then run this step again. Do NOT relax a test, an assertion, a threshold or a golden to get past this line, and do NOT move the tests out of ./cmd/wisp/ - that is the shape ticket 98 AC#2 forbids.
post_final_rc=1
```

日志只有 2 行（310 的 GOTMPDIR ＋ 守卫），**没有第二个测试进程**：守卫在 `bash "$portable"` 之前就 exit 了。
这也保住了票 311 那条**装载顺序**读数（"310 行 -> 311 GUARD -> 没有测试进程"）：新守卫排在 311 那一块**之后**、
在 `export PATH` **之后**、在委托**之前**，跨卷的跑法仍然在 311 就退出，看不到这一行。

**2.2 不删注入的正常发（守卫必须安静，且不许把真红读成加载失败）**

`bash scripts/wisp-cli-tests.sh`（本腿最终字节，17:37:00 起）：

```
portable-tests.sh: four numbers (all from -v output): === RUN=498  --- PASS=358  --- FAIL=41  --- SKIP=1
runtests.sh: go test exited 1 ... top-level: PASS=358 FAIL=41 SKIP=1, === RUN=498, '[no tests to run]'=0
portable-tests.sh: unaccounted SKIP lines ... resident_audio_247_live_windows_test.go:130: AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 ...
FAIL	github.com/CarlosShao/wisp/cmd/wisp	600.175s
script_rc=1
```

- 守卫行出现 **0** 次（`grep -c "ticket 98 AC#3"` = 0）⇒ PATH 扫描在**真装的 DLL 路径**上是绿的，**不误报**。
- 测试**真的跑到了**（498 个 `=== RUN`）⇒ 这一发不是"加载失败被误标"，**第二半守卫也没有误报**
  （它要求"有加载失败签名"且"一条测试结果行都没有"两个条件同时成立）。
- ⚠ **这个 rc=1 不是本腿造成的，也不是守卫造成的**：`ff61b9dbc` 这棵树上 `./cmd/wisp/` 自身有 41 枚 FAIL、
  1 枚 SKIP（真麦克风，`WISP_LIVE_MIC=1` 才跑），最后撞 `go test` 默认 10m 上限（`600.175s` panic）。
  本腿**没有**碰任何测试、任何断言、任何阈值；下面的对照证明这些红**脱离包装脚本也存在**：

```
$ PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -timeout 420s \
    -run 'TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims|Test15eSpeechEndToEndHandsTheTranscriptToTheAgent|TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink' ./cmd/wisp/
rc=1
--- FAIL: TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims (0.00s)
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (7.62s)
--- FAIL: Test15eSpeechEndToEndHandsTheTranscriptToTheAgent (0.02s)
```

⇒ 41 枚红与那枚 SKIP 是**本 HEAD 的在册现状**（票池里另有主人），本件只登记、不追、不修。

**2.3 为什么守卫不写成 `where`，也不写成 `-test.list` 探针（两条都是实测否掉的）**

- `where sherpa-onnx-c-api.dll` 在**钉的目录不在 PATH 上**时仍然回了
  `D:\work\workspace\projects plans\Wisp\sherpa-onnx-c-api.dll`：`where` **也搜当前目录**。
  仓根确实躺着一枚同名 DLL（`git ls-files sherpa-onnx-c-api.dll` -> *did not match any file(s)* ⇒ 未跟踪的野文件），
  所以从仓根直接起一个 `go test -c` 出来的二进制**能加载**——这正是"探针绿、门却是瞎的"那形。
- 同一个二进制换到临时目录起，加载期就红，且**加载器自己点名**：
  `error while loading shared libraries: sherpa-onnx-c-api.dll: cannot open shared object file: No such file or directory`（rc=127）。
  `go test` 恰恰是把测试二进制放到**临时目录**、以**包目录**为 cwd 起的 ⇒ 仓根那枚野 DLL 救不了它
  （审计读数：`go test ./cmd/wisp/` 在无注入 PATH 上 = `exit status 0xc0000135`、rc=1、`=== RUN`/`--- FAIL` 各 0）。
- `-test.list` **的颜色取决于站在哪儿**：同一枚二进制、同一条 `-test.list` 命令，从**仓根**起是
  **rc=0** 并打出 **482** 行顶层用例名（这个数本身也是 AC#1(ii) 想要的分母），从 `/d/tmp` 起是 **rc=127**
  加载失败。⇒ 它当分母可以，当守卫不行：一个颜色随 cwd 翻转的检查，正是本票要杀掉的那形；
  何况它还要先付一次 `go test -c` 的编译。

## 3. 改了什么（`scripts/wisp-cli-tests.sh`，+91/−1）

1. **PATH 可达性守卫（主守卫，确定性）**：`export PATH="$dll_dir:$PATH"` 之后，对 `deps.toml` 里
   `[sherpa-onnx.dll."<name>"]` 抽出的**每一枚** DLL，自己遍历 `PATH` 的每一项（反斜杠归一成 `/`）看文件在不在；
   有一枚不可达就红并**逐名列出**，`exit 1`。不读 cwd、不读 `where`、没有 `|| true`、没有环境变量开关、不可 skip。
2. **加载失败兜底（第二半）**：委托 `portable-tests.sh` 的输出**照旧透传并保留其退码**（`tee` 一份日志），
   只有当"退码非 0 ＋ 日志里有加载失败签名 ＋ 日志里**一条** `--- PASS|FAIL|SKIP` / `=== RUN` / `ok` 都没有"
   三条同时成立时，才补一行**具名的 setup 失败**（列出 `deps.toml` 的钉名与日志路径），然后**用原退码**退出。
   真红（有 `--- FAIL`）不会被它改色，`OK` 更不会。
3. 一处**不是**放宽：日志落盘在 `GOTMPDIR`（310 已设）下，路径打印在告警里供附票。

## 4. 本件**没有**建立的事（宁缺毋造）

- **AC#1(iii) 的 CI 步级结论**：没取 run id。那要一次真 run（真机/CI 面），本腿不越界；按票面判据"说不出 run id 就当它不存在"。
- **AC#4 的台账收敛**：没做（另有其格）。本件只报"哪些读数对得上"，不改 `docs/reports/**`。
- **ubuntu 腿的重新普查**：`ae8136261`（180-orphan，2026-11-05 13:48）已按 `98-attribution` §2 C1 把 **39 枚** POSIX 红转绿
  （其提交信息逐字：`修前 PASS=149 FAIL=79 -> 修后 PASS=188 FAIL=40`）。⇒ `98-attribution.md` 里那张"81 枚/76 枚 C1"
  的表**已过期**，本腿**没有**重跑容器普查，**因此本件不引用 81 当现读数**，也不主张剩余 40 枚的逐枚归属。
- **裸 `go test ./cmd/wisp/` 仍然不点名**：要让**裸命令**也响亮，只能在 `cmd/wisp` 里加 `TestMain`/init 钩子
  = 动共享编译面 Go 包 = 本腿写面之外（也正是票 98 AC#2 禁止的"把被检对象挪出去"的近邻）。⇒ 只登记：
  从今往后**凡跑 `cmd/wisp` 走 `scripts/wisp-cli-tests.sh`**，裸命令的红**不是裁决**。
- **41 FAIL / 1 SKIP / 600s 超时**：本 HEAD 的现状，本腿一个字没改，只登记（见 2.2 的对照）。

## 5. 门禁（本腿现跑，不用别人的数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| gofumpt（受检树） | `gofumpt -l cmd internal tools scripts` | **空，rc=0** |
| 形式检查 | `bash -n scripts/wisp-cli-tests.sh` | **rc=0** |
| 票 98 主判据 | 见 §2.1 | 改前**不点名**、改后**点名 2 枚**；两发 rc=1 |
| 包测试 | `scripts/wisp-cli-tests.sh` -> `runtests.sh`（scope=cli） | `PASS=358 FAIL=41 SKIP=1`，`=== RUN=498`，**rc=1**（红因见 2.2，非本腿） |
| D22 扫描 | `sh scripts/d22scan.sh` | **rc=0**：正控 `runtests.sh: OK ... PASS=35 FAIL=0 SKIP=0`；扫描 `clean - no D22 ban violations`（`internal/=363`、`cmd/=57`、`frontend/=116`、`design/=39`、`internal/tools/=48`） |

⚠ **D22 那一格有一发中间读数，写在这里省下一位的误判**：17:37 第一次跑是 **rc=1**，唯二红是
`TestScannerSelfScanOfRealRepoIsGreen` / `TestRealRepoLedgerIsHonest`，理由是
`internal/config/overrides.go:43` 引用了 `docs/evidence/s1/230-ac4-config-half.md`（当时盘上没有）。
**那不是本腿、也不是 HEAD 的病**：这两枚都**未跟踪**（`git ls-files` 对二者皆空，`git status` 是 `??`），
是另一条腿在共享树里**在飞**的文件；17:51 它们落盘（`230-ac4-config-half.md` 13,410 字节）后同一命令复跑即
**rc=0**。⇒ 共享树里 D22 的读数要**在别的腿落地后重取**，别把别人的半成品当成自己的红。

⚠ `gofumpt -l .`（整仓）会打出几十个 `.scratch/**` 探针文件与几个语法坏的负控样本——那些是**未跟踪**的临时件，
不是本腿产物；受检树的三目录读数才是干净的空。

## 6. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
unset GOTMPDIR TMPDIR TEMP TMP                 # 派单要求

# 1) 判据两发（构造同 §1；只改 root 与注入两行）
R="/d/work/workspace/projects plans/Wisp"
git show HEAD:scripts/wisp-cli-tests.sh | \
  sed -e "s|^root=.*$|root=\"$R\"|" -e 's|^export PATH=.*$|:|' > /d/tmp/pre-mutant.sh
sed  -e "s|^root=.*$|root=\"$R\"|" -e 's|^export PATH=.*$|:|' scripts/wisp-cli-tests.sh > /d/tmp/post-mutant.sh
bash /d/tmp/pre-mutant.sh  ; echo "pre_rc=$?"    # 期望：exit status 0xc0000135，无点名，0 条测试
bash /d/tmp/post-mutant.sh ; echo "post_rc=$?"   # 期望：GUARD (ticket 98 AC#3) 点名 2 枚 DLL，rc=1

# 2) 正常发（守卫必须安静；测试会真的跑到，红因见 §2.2）
bash scripts/wisp-cli-tests.sh ; echo "script_rc=$?"

# 3) 为什么不能用 where / -test.list 当守卫（重点是 cwd 会翻转颜色）
go test -c -o /d/tmp/wisp-cli.test ./cmd/wisp/            # 编译本身 rc=0，坑在起进程
cd "/d/work/workspace/projects plans/Wisp" && PATH="/c/Windows/System32:/c/Windows" \
    /d/tmp/wisp-cli.test -test.list '.*' | wc -l          # 从仓根起 => rc=0，482 行（仓根那枚野 DLL 救了它）
cd /d/tmp && PATH="/c/Windows/System32:/c/Windows" \
    /d/tmp/wisp-cli.test -test.run 'TestSecretSetGetListRoundTrip'
#   换到临时 cwd 起一条真的用例 => rc=127，加载器逐字点名 sherpa-onnx-c-api.dll（同 §2.3）
cd "/d/work/workspace/projects plans/Wisp" && PATH="/c/Windows/System32:/c/Windows" \
    where sherpa-onnx-c-api.dll     # 钉的目录不在 PATH 上，仍回仓根那枚 => where 不能当守卫
```
