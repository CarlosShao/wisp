# `20-shapes.md` — 表②：死亡三形对 `%TEMP%` 下 `Test*` 目录的**实测**（票 310 `AC#0`，腿 `310-a1`）

台件＝**仓外导出树**（共享工作树一枚字节都没改）：

- `mktemp -d` → bash 路径 `/tmp/310a1-tree.vMLB73`，Windows 路径 `C:\Users\swq\AppData\Local\Temp\310a1-tree.vMLB73`
- 六件套导出尺：`git archive HEAD go.mod go.sum internal cmd tools scripts | tar -x -C "$T"` → `archive_rc=0`
- 证非空：`find "$T" -name '*.go' | wc -l` → **`go_files=682`**（⛔ 少带 root；本腿没遇到 `walk …: The system cannot find the file specified.` 那形台面故障）
- 六枚 root 逐枚存在：`go.mod=yes go.sum=yes internal=yes cmd=yes tools=yes scripts=yes`

窄用例＝**`internal/config` 的 `TestRoundTripLoadMarshalLoad`**（`git grep` 现量：`internal/config/boundary_test.go:186` 定义，`:188` 调 `sealableTempDir124(t)` → `internal/config/tempdir_resolved_124_test.go:33` 的 `proc.SealableRoot(t.TempDir())`）。
选它的理由：只调 `t.TempDir()`、不起子进程（本包的 `exec.Command` 只在 `c26_seam_posix_125_test.go`／`private_acl_windows_test.go` 里，被我的 `-run` 正则排除）、不碰麦／不碰 webview／不碰 GUI。
`-list` 尺另证：`go test -list 'TestRoundTripLoadMarshalLoad$' ./internal/config` 在本树可列出该用例（见 `logs/20-shot-a.txt` 的 `=== RUN` 行）。

★**本腿一共只跑了这几条 Go 用例**（⛔ 整包 `go test ./...`、⛔ `GOFLAGS= go build ./...`）：
`TestRoundTripLoadMarshalLoad`（形 a／b／c／d／e 全部用它）＋ 甲形实测那两发用的六枚同包用例（正则原文＝`^(TestRoundTripLoadMarshalLoad|TestRoundTripBytes|TestResolvedNeverPersists|TestSaveFileRoundTripsThroughLoad|TestLoadFileMergesDefaultsAndOverrides|TestUnwiredGuardLeavesHonestConfigsAlone)$`）＋ 一枚编译产物 `config310a1.test.exe`（`go test -c -o`，只为本例）。

分母尺（每一形同批跑，件＝`logs/20-before-*.txt`／`logs/20-after-*.txt`，逐枚名＋`comm -13` 作差）：
`powershell 'Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like "Test*" } | ForEach-Object { $_.Name }'`

## 1. 逐形读数（⛔ "0 枚新增"当裸读数——每一形都带"那一形真发生了"的同行证据）

| 形 | 命令（在导出树里跑） | 跑前 `Test*` | 跑后 `Test*` | 新增枚数 | 新增目录具名名册 | **"那一形真发生了"的证据行（逐字）** |
|---|---|---|---|---|---|---|
| (a) 正常收尾 | `go test -count=1 -run '^TestRoundTripLoadMarshalLoad$' -v ./internal/config` | 50 | 50 | **0**（empty above＝0 枚） | — | `=== RUN   TestRoundTripLoadMarshalLoad` ＋ `--- PASS: TestRoundTripLoadMarshalLoad (0.02s)` ＋ `ok github.com/CarlosShao/wisp/internal/config 0.053s`；`shot_a_rc=0 wall_ms=1185` |
| (b) `go test -timeout` panic | `go test -count=4000 -timeout 2s -run '^TestRoundTripLoadMarshalLoad$' -v ./internal/config` | 50 | 51 | **1** | `TestRoundTripLoadMarshalLoad725496612` | 逐字：`panic: test timed out after 2s`（`logs/20-shot-b.txt:507`）＋ `:508` 的 `	running tests:`；`shot_b_rc=1 wall_ms=3275`；同发 `run_lines=253 pass_lines=252` ⇒ **帽落在这条用例的运行之中**（第 253 发正在跑，前 252 发各自收尾删干净了） |
| (c) 外力 `TerminateProcess` 打死测试二进制（第 1 发） | `go test -c -o config310a1.test.exe ./internal/config` 后由 `logs/20-shape-c-watch.ps1` 同进程起＋盯＋`Stop-Process` | 50→（b 后 51） | 52 | **1** | `TestRoundTripLoadMarshalLoad2289939324` | 同一进程内逐字：`shapec world=CONSTRUCTED dir=TestRoundTripLoadMarshalLoad3749761498 dir_creation=10:46:41.776`／`shapec target_name=config310a1.test.exe pid=7980 ppid=880 this_ps_pid=880`／`shapec cmdline="C:\Users\swq\AppData\Local\Temp\310a1-tree.vMLB73\config310a1.test.exe" -test.run ^TestRoundTripLoadMarshalLoad$ -test.count 4000`／`shapec external_kill_issued at=10:46:42.368`／`shapec after_kill target_pid_still_present=0` |
| (c) 外力杀（第 2 发，改了脚本的差集尺） | 同上 | 52 | 53 | **1** | `TestRoundTripLoadMarshalLoad2836872785` | 逐字：`shapec2 world=CONSTRUCTED dir=TestRoundTripLoadMarshalLoad1350067340 dir_creation=10:47:24.525`／`shapec2 guard_name_is_test_exe=True guard_cmdline_has_export_tree=True guard_parent_is_this_leg=True`／`shapec2 external_kill_issued at=10:47:25.129`／`shapec2 after_kill target_pid_still_present=0`／`shapec2 at_kill_names=TestRoundTripLoadMarshalLoad2289939324,TestRoundTripLoadMarshalLoad2836872785,TestRoundTripLoadMarshalLoad725496612` |
| (d) ★甲形变体实测：`GOTMPDIR` 指仓外根 ＋ (b) 同一形 | `GOTMPDIR=/tmp/310a1-gotmp.i3FTVN go test -count=4000 -timeout 2s -run '^TestRoundTripLoadMarshalLoad$' ./internal/config` | 53 | 53 | **0**（`%TEMP%` 顶层 empty above＝0 枚） | 残渣落在 GOTMPDIR 根里：`/tmp/310a1-gotmp.i3FTVN/TestRoundTripLoadMarshalLoad286226934`（1 枚） | 逐字：`panic: test timed out after 2s`（`logs/20-shot-d.txt`）；`shot_d_rc=1 wall_ms=3254`；`gotmp_root_entries=1` |
| (e) ★甲形变体实测：`TMP`/`TEMP` 指仓外根 ＋ (b) 同一形 | `TMP=… TEMP=… go test -count=4000 -timeout 2s -run '^TestRoundTripLoadMarshalLoad$' ./internal/config` | 53 | 53 | **0**（empty above＝0 枚） | 残渣落在重定向根里：`/tmp/310a1-aredir.cn4yfa/TestRoundTripLoadMarshalLoad3605550926`（1 枚） | 逐字：`panic: test timed out after 2s` ＋ `FAIL github.com/CarlosShao/wisp/internal/config 2.043s`；`shot_e_rc=1` |

⛔ 我这几发里没有任何一发是"帽设在用例跑完之后"：(b)(c)(d)(e) 都同批发了 `=== RUN`／`world=CONSTRUCTED`／`panic:` 三样中的相应两样以上，并且 (b) 直接把"跑到第 253 发才被砍"写在同一张 `logs/20-shot-b.txt` 里。

## 2. ★跑完当场自报：本腿这几发新增的目录（归我、可证＝本仓，⛔ 删）

`Test*` 计数从本腿起手的 **50** 涨到交件的 **53**（`%TEMP%` 顶层），另外在两枚仓外重定向根里各留 **1** 枚。逐枚：

| # | 目录名（或全路径） | 哪一发留下的 | 归属 |
|---|---|---|---|
| 1 | `TestRoundTripLoadMarshalLoad725496612` | (b) `-timeout 2s` panic | 本腿／本仓可证（同用例名在 `internal/config/boundary_test.go:186`） |
| 2 | `TestRoundTripLoadMarshalLoad2289939324` | (c) 第 1 发外力杀 | 本腿／本仓可证（同上） |
| 3 | `TestRoundTripLoadMarshalLoad2836872785` | (c) 第 2 发外力杀 | 本腿／本仓可证（同上） |
| 4 | `/tmp/310a1-gotmp.i3FTVN/TestRoundTripLoadMarshalLoad286226934` | (d) GOTMPDIR 变体 | 本腿／本仓可证；⚠ **不在 `%TEMP%` 顶层**，被 roster 尺看不见 |
| 5 | `/tmp/310a1-aredir.cn4yfa/TestRoundTripLoadMarshalLoad3605550926` | (e) TMP/TEMP 变体 | 本腿／本仓可证；同上 |

- 三枚现量确认尺（交本件前又跑了一次）：`total=2763 test=53`，`mine=TestRoundTripLoadMarshalLoad2289939324 created=10:46:42`／`…2836872785 created=10:47:25`／`…725496612 created=10:45:38`（件＝`logs/` 同批；`rc=0`）。
- ⚠ 上面 #4 #5 说明**"残渣不见了"和"残渣不见了"是两回事**：我那两枚重定向根本身就建在 `%TEMP%` 底下（`mktemp -d` 的默认落点），所以表③ 的甲形若照这个形状落，`Get-ChildItem $env:TEMP -Directory` 这把**非递归**尺会读出 `0 新增`——**看着像修好了，其实是把残渣挪到尺的射程外**。⇒ 甲形的根必须**显式落在 `%TEMP%` 之外**（例如 `D:\tmp\…`），否则它只是把噪声藏起来。这条是本腿实测出来的，不是推理。

## 3. 三形小结（本票的机制结论）

1. **(a) 正常收尾＝0 枚新增**，且这条用例每跑一发都会把自己那枚删掉（同发跑了 1 遍；(b) 那发跑了 252 遍、`%TEMP%` 顶层仍然只多了 1 枚＝正在跑的那一发）⇒ **`t.TempDir` 的 cleanup 在正常路径上是可靠的**，本票的账只记在"外力致死"这一族上。
2. **(b) `go test -timeout` panic＝漏 1 枚**：`panic: test timed out` 走的是 runtime 直接杀进程路径，`t.Cleanup` 里那个 `removeAll` 根本不执行。
3. **(c) 外力 `TerminateProcess`＝漏 1 枚**，两发各漏 1 枚，且两发都同批发到了 `target_pid_still_present=0`（＝被杀的那枚测试二进制当场没了，父侧没人替它收尾）。
   ⚠ 本腿在 (c) 第 1 发的台件里犯了一个**读数尺的错，具名记下**：脚本的 `verdict` 拿"盯到的那枚目录"去判存在性，而迭代只有约 10 ms 一发，等 `Stop-Process` 落下时那一枚早已被它自己的迭代收尾删掉、真正留下的是**当时正在跑的另一发迭代的目录**。⇒ 第 1、2 发脚本自报的 `dir_still_on_disk=0`／`verdict=NO-LEAK` **都读反了**，真相由外部的 `comm -13`（跑前／跑后名册作差）给出：**各漏 1 枚**（见 §1 表与 §2 名册）。这不是台面故障（六件套齐全、`go_files=682`），是本腿的判据建错了对象，改后的尺（快照＋差集）在同一发的 `at_kill_names` 里也对上了。
4. **(d)(e) 甲形的两个变体都能把残渣移出 `%TEMP%` 顶层**，且都不改变"漏"这件事本身：换根之后照样有一枚目录留在盘上没人删——**甲形治的是"落在哪"，不治"没人收尾"**。票 310 真正要的"收尾"要么落在乙（下次跑测先扫自己上次留下的），要么就得有第四个形（给测试二进制装 Job Object／`KILL_ON_JOB_CLOSE` 一类的宿主侧回收——那是票 304 的进程面射程，本票边界里明写 ⛔ 塞进去）。

## 4. 自落 rc 与件清单（每件 `wc -c` 非 0）

```
logs/20-shot-a.txt            PASS 全貌      rc 行＝shot_a_rc=0
logs/20-shot-b.txt            803 行含 panic rc 行＝shot_b_rc=1
logs/20-shot-c.txt            c 第 1 发      rc 行＝shapec_rc=1（脚本判语错，见 §3.3）
logs/20-shot-c2.txt           c 第 2 发      rc 行＝shapec2_rc=1（同上）
logs/20-shot-d.txt            GOTMPDIR 变体  rc 行＝shot_d_rc=1
logs/20-shot-e.txt            TMP/TEMP 变体  rc 行＝shot_e_rc=1
logs/20-build-c.txt           go test -c     rc 行＝build_rc=0 exe_exists=yes
logs/20-before-*.txt / 20-after-*.txt / 20-new-*.txt   名册作差尺各自行数自落
logs/20-shape-c-watch.ps1     台件本体（ASCII，⛔ emoji ⇒ 跟踪件不会被 d22scan 的字符串段判违规）
```

引用 provenance：表内所有命令与读数＝**本腿在仓外导出树里现跑**（⛔ 引用编排者快照）；用例定义行＝**HEAD blob**（`git grep -n … HEAD -- internal/config`）。
