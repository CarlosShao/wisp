# 300-v1 — `30` ④ 成对导出树那枚 harness 坑 ＋ ①②③逐格判语 ＋ `AC#2` 开⛔ 开

锚＝`efd85155`。本件是**结论件**；每条读数都指回我自己的原始件（`logs/*.txt`），⛔ 引用编排者或被审腿的数当凭据。

## 1. ④ 三档结论：**无影响**（具名理由）——⛔ 打折；但**记两笔具名缺陷**（其中一笔归编排者核）

### 1.1 我量到的事实（六把尺，全部我自己现跑）

| # | 尺（命令逐字） | 射程 | 读数 |
|---|---|---|---|
| R1 | `git ls-files third_party/sherpa-onnx \| wc -l` | 该目录 | **0** ⇒ tracked＝0 ⇒ `git archive` 的导出树里**必然没有**那三枚 DLL |
| R2 | `ls third_party/sherpa-onnx/*.dll \| wc -l` | 工作树同一目录 | **3**（`onnxruntime.dll` 17799168／`sherpa-onnx-c-api.dll` 4605952／`sherpa-onnx-cxx-api.dll` 259584，尺＝`ls -l`）⇒ 未跟踪 |
| R3 | `grep -rln 'sherpa\|onnxruntime' internal/audio` | `internal/audio/**` | **0 命中**（⚠ 派单举的那个例子这一层**没有东西**，具名见 §1.4） |
| R4 | `grep -rln 'import "C"' --include=*.go .` | 全仓 `.go` | **0** ⇒ 仓自身无 cgo；但依赖有（见 R5） |
| R5 | `grep -rn '#cgo LDFLAGS' $GOMODCACHE/github.com/k2-fsa/sherpa-onnx-go-windows@v1.13.8/build_windows_amd64.go` | 模块缓存 | 逐字 `// #cgo LDFLAGS: -L ${SRCDIR}/lib/x86_64-pc-windows-gnu -lsherpa-onnx-c-api -lonnxruntime` |
| R6 | `objdump -p cmdwisp.test.exe \| grep -i "DLL Name" \| sort -u`（那枚 exe＝`go test -c ./cmd/wisp/` 在**我的导出树**里造的） | 一枚测试二进制的 PE 导入表 | `KERNEL32.dll` / `msvcrt.dll` / **`sherpa-onnx-c-api.dll`** ⇒ `cmd/wisp` 的测试二进制**静态导入那枚 DLL**，载入期就要它 |

原始逐字＝`logs/prod-and-dll.txt`、`logs/sherpa-deps.txt`、`logs/cmdwisp-loadprobe.txt`。

### 1.2 那枚坑**是真的**，而且是**载入级**（我自己造了两形树，⛔ 引用任何人）

同一枚导出树（`/tmp/wisp300-v1/treeA` ＝ `git archive efd85155`，`third_party/sherpa-onnx` **不存在**），
**同一条命令、只用派单写死的那条 harness**，差别只有"那三枚 DLL 在⛔ 在"：

```
形 A（树上无 dll）
  go test ./cmd/wisp/ -count=1 -run '^$' -v
  exit status 0xc0000135
  FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.033s
  FAIL
  rc-pipeline=1
形 B（把 3 枚 dll cp 进 $PWD/third_party/sherpa-onnx 后，同一条命令）
  testing: warning: no tests to run
  PASS
  ok  	github.com/CarlosShao/wisp/cmd/wisp	0.070s [no tests to run]
  rc-pipeline=0
```

- ⇒ **形 A 里 `cmd/wisp` 这一档产出的 `--- PASS`／`--- FAIL` 枚数＝0**（它⛔ 是"红"，它**根本没跑**：没有一枚用例名、只有一行包级 `FAIL`）。这正是编排者台账里那族 `0xc0000135 且⛔ --- FAIL` 的形状，我这一程**第一次拿到可复现的逐字凭据**。
- ⚠ `-run '^$'` 那一发**零用例执行**（只建、只载入），⛔ 跑整包 `cmd/wisp`（约 8 分钟那一发仍归编排者）。我自报这一发：**跑了、rc 逐字在上**。
- 同一棵无 dll 树里我把三枚⛔ 依赖 DLL 的包跑到底（尺＝`grep -c -- '--- PASS'`，⛔ 截断：`grep -c 'panic:'`＝**0**、`=== RUN` 行数＝**322**）：

```
go test ./internal/panel/ ./internal/audio/ ./internal/ball/ -count=1 -v -timeout 240s
grep -c -- '--- PASS'  ⇒ 313          grep -c '^--- PASS' ⇒ 231          grep -c -- '--- FAIL' ⇒ 8
ok   github.com/CarlosShao/wisp/internal/audio  16.362s
ok   github.com/CarlosShao/wisp/internal/ball    0.193s
FAIL github.com/CarlosShao/wisp/internal/panel    6.564s     （真红，⛔ 载入问题；它跑完了）
```

⇒ **`internal/audio` 与那三枚 DLL 无关**（派单 harness 那条 `PATH` 对它 inert；a1r 的 `cp` 是多余的保险，⛔ 污染它的读数）。

### 1.3 对 `293-v1` 那两组的裁决：**⛔ 打折**（凭据＝它自己的红名册逐枚定位）

尺＝把 `293-v1` 件 `logs/prepost-compare.txt` 的 POST 红名册里我抽的 **8 枚**逐名 `grep -rl "func <名>" --include=*_test.go cmd internal`：

```
TestAC2RealProcessRefusesOnEveryLegWithoutAppData128 => cmd/wisp/dataroot_128_windows_test.go
TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord  => cmd/wisp/early_log_nail_130_windows_test.go
TestAC246ShippedResidentProcessOwnsItsCancelStep     => cmd/wisp/resident_approval_246_windows_test.go
TestAC228ResidentLegReportsAndBooksItsBall           => cmd/wisp/resident_ball_228_windows_test.go
TestAC1ResidentLegInstallsItsLogListenerOnDisk       => cmd/wisp/resident_sink_nail_127_windows_test.go
TestSecretArgvCarriesNoSecret                        => cmd/wisp/secret_argv_windows_test.go
TestCleanCheckoutBuilds_AC11                         => cmd/wisp/panel_host_gate_test.go
TestPanelHostRealWindowHopAndLifecycle               => cmd/wisp/panel_host_windows_test.go
```

**8/8 全在 `cmd/wisp/`** ⇒ 它那两棵导出树里 `cmd/wisp` **确实跑起来了**（有逐枚用例名、有真实耗时 2.6s／5.05s／20.06s）
⇒ 它的树**必然是 DLL 齐的**；否则这一档会像我的形 A 那样只剩一行包级 `FAIL`、枚数为 0。
第二把尺同向：`TestAC293` 本体在 `cmd/wisp/resident_tray_mute_293_windows_test.go`，而它交回 `post1/post2 AC293 top-level PASS=7` ⇒ **同一枚包载入过**。
第三把尺同向：我这一发的 3 包总数 **313** 对上它的 **677**，缺的 **364**（≈54%）只能由 `cmd/wisp` 那一档补齐——与"它在里面跑了"一致。

⇒ **结论＝无影响**：`677 ↔ 684` 那对作差（＋7）与两枚方向的对称差（`comm -3` 皆空）**⛔ 被这枚坑污染**，⛔ 打折，也⛔ 反过来加强它。

### 1.4 但要记两笔具名缺陷（⛔ 改判语，只记账）

1. **`293-v1` ⛔ 命名这枚前置条件**（我把它整目录 grep 过：`sherpa`／`.dll`／`0xc0000135`／`[build failed]` 只命中它自己写的 harness 那一条 `PATH`，见 `logs/293-and-treeA.txt`）。
   它的件 §4.1 逐字写"同一台机器、同一条 PATH、同一无 `.git`／无 `build/` 的环境 handicap"——**漏了唯一一枚会静默决定 54% 分母存在的 handicap**。
   ⇒ 实际后果（我实测）：**下一位照它的逐字配方建一棵新导出树，得到的是 313 而⛔ 677**，且红名册会凭空少一整档。这是**复现性缺陷**，⛔ 读数错误。⇒ 建议：成对导出树那套定式（`PLAN`/票面口径里那句"改前基线永远取得到"）应补一句"**导出树必须显式声明 dll 就位那一发怎么做的**"（这条归编排者落，⛔ 我改任何定式件）。
2. **我⛔ 能说明它那三枚 DLL 是**怎么**进树的**（tracked＝0，`git archive` 带不进未跟踪件）。可能性：`$TEMP/pre293`/`post293` 被别的腿共用过、或它 `cp` 了没写。
   ⇒ **这一笔归编排者核**（只有你能查那两棵树的当前形态与 `date`）。我这把尺只能证明"它们当时在"，⛔ 证明"是谁放的"。
   ⚠ 这一笔⛔ 动摇 §1.3 的"无影响"：它动摇的是**记录完整性**，⛔ 是那些读数。
3. ⚠ 派单那句"`internal/audio` 里带 `sherpa`／`onnxruntime` 字样的那条"**在 `internal/audio` 里⛔ 存在**（我的 R3＝0 命中）。我改用内容锚另找（`cmd/wisp/dataroot_128_windows_test.go` 头顶逐字 "the sherpa/onnxruntime DLLs - without that the child dies at load time with **0xC0000135** and zero output"；`scripts/portable-tests.sh:129` 逐字 "cli cmd/wisp, whose test binary needs the sherpa DLLs"），⛔ 照例子空跑。**具名报回这一条派单过期处。**

## 2. ①②③ 逐格判语（每条指回我自己的件）

| 格 | 判 | 我的凭据（⛔ 引用编排者的复跑） |
|---|---|---|
| **① `AC#0` 权威凭据有没有牙** | **成立**：偏移＝**24**；`wasapi_windows.go` 是错的那一枚，`wavinjector.go` 与权威同形 | `10` 件：我自己 `sed` 抽的 `WAVEFORMATEX`／`WAVEFORMATEXTENSIBLE`／`guiddef.h GUID` ＋ 头文件两处逐字 `/* Format.cbSize = 22 */` ＋ 我自己算的偏移表（18＋2＋4＋16＝40，`SubFormat` 起 24）；26 那一形会让 `cbSize`＝24、总长＝42 ⇒ 与盘上原文**直接矛盾** |
| **①b 那句 union 之争谁对** | **腿对、票面 `:13` 后半句错** | `10` 件 §5：`Samples` union 三成员逐字皆 `WORD`（2 字节），`} Samples;` 之后 `DWORD dwChannelMask;` 是结构体**直接成员**；`grep -c "ValidSamples" mmreg.h`＝**0** ⇒ 标准结构里⛔ `nValidSamples`/`wValidSamples` 这枚名字。**前半句（validBits⛔ 独立字段）对。** |
| **② 台件是不是循环论证** | **⛔ 循环 ⇒ 成立（有牙）** | `20` 件：夹具 10 处写入逐处落在我自算的权威表上（⛔ 一处用 `@18/@22/@26`）；我在看读数**之前**用"权威＋夹具"手算出锚上四形该读 `0/0/0/3`，实测**逐字命中**（若夹具来自产码假设，锚上会全绿）；我自己重跑 **4/4 FAIL、`rc-gotest=1`、无 `0xc0000135`** |
| **②b 四形是不是"两形⛔ 都放行"** | **⛔ 都放行**（套件层）；**但单形层有三枚是废尺** | `20` 件 §4：五枚偏移扫描 **20/22/24/26/28** ⇒ 全绿**只在 24**；26 全红；`S1` 单形在 **20 与 24** 都绿、`S4` 单形在 **22/24/28** 都绿 ⇒ 判别力实际扛在 **S2＋S3** 上 |
| **③ `AC#2` 入库那枚要不要保留 S4** | **要保留，但⛔ 因为腿给的那条理由**；"少了它就变成一把只能单向响的尺"这一句**不成立** | 见下面 §3 |

## 3. ③ 单独裁（本格⛔ 写用例，只裁形状）

- **判：入库版少了 S4 ⛔ 会退化成单向尺。** 我的五枚偏移扫描直接回答它：剩下的 `S1+S2+S3` 子集在 **20 红**（S2/S3 读到 3）、**22 红**、**26 红**、**28 红**、**只在 24 全绿** ⇒ 两形（24／26）里它只放行一形 ⇒ 方向性⛔ 依赖 S4。
- **腿那句"票面'两枚都读出 0'只对良构面成立"——成立**（我复现：S4 在锚上读到的是 **3**）。但它是**读数层面的更正**，⛔ 仪器必要性的论证；腿把它当"建议保留 S4"的理由给，**这条理由我⛔ 接受**。
- **我给的保留理由两条（各自独立）**：
  1. S4 是四形里**唯一一形**让锚上那枚错位读取产出一枚**看似合法的 `FLOAT(3)`** 而⛔ 是 0 ⇒ 它把票面「为什么要紧」那句的形状（"浮点被当 int16 切"的反向：**从错字节位置读出对的数字**）钉在盘上；S1/S2/S3 只会证明"读到空"。
  2. 成本＝**零**：正确偏移下它是绿的（我实测），⛔ 造恒红跟踪测试（那才是本仓定式真正禁的东西）。
- ⚠ **同时必须写进入库用例的一句**（我给落地腿的输入，⛔ 我写用例）：**S4 是四枚里最弱的一枚**（放行 22/24/28 三形）⇒ ⛔ 只留它、⛔ 把它单独当"决定性那一发"；`S1`/`S2` 同样⛔ 单独成立。**这套件的强度只在四枚合起来。**

## 4. `AC#2` **开**（可派落地腿）——一句结论 ＋ 四条派单必须逐字带的输入

**开。** ①②③三格我全部独立复现到"权威＝24、台件有牙、改前必红凭据齐"，硬门（"⛔ 是 26 才停手"）未触发；
票面 `AC#2` 的判据形状（改成与权威同形那枚 ＋ 把头顶注释改成实话）⛔ 依赖任何我⛔ 交的数据。

落地腿必须逐字带着这四条（都来自我自己的读数，⛔ 转述）：

1. **头顶注释里三个数**全是坏的（逐字现读：`// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].`）
   ⇒ 权威形＝`Samples@18（2 字节 union）, dwChannelMask@20, SubFormat@24`；`bits@14`、`cbSize@16`、`tag@0`、`channels@2`、`rate@4` 这五枚**本来就对**（我逐字段核过），⛔ 顺手改它们。
   ⚠ **⛔ 只改 `26` 一个数**＝把两枚假话留在原地（这是腿给的第一条输入，我**独立证实**）。
2. **偏移那一枚**（内容锚逐字 `sub := *(*windows.GUID)(unsafe.Add(p, 26))`）改成 `24`；我实测这一改让台件 4/4 绿（`20` 件 §4）。⛔ 别的形状（20/22/28）我一枚一枚扫过，**全红** ⇒ "改成 24"是唯一解，⛔ 试错空间。
3. **入库用例的形状**：四形齐进（含 S4，附 §3 那句"最弱一形"的注释），⛔ 断"被调用过"，断的是解析出的 `tag` 值。
4. **⛔ 出射程**：`internal/audio/**` ＋ `probes/300/**`；⛔ 动 `cmd/wisp`（票 255 名册按行号引产码行，先例 `A799`/`A801`）；
   ⛔ 动 `internal/audio/level.go` 的 `SineLevelTolerance`、`internal/ball/liquid.go` 的 `SilenceLevelGate`、`internal/observe/thresholds.go`；
   ⛔ 用真设备／真麦克风（`AC#3` 那发 `GetMixFormat` 仍归编排者，**这仍是我这一程⛔ 能答的唯一硬缺口**：本机混音格式到底⛔ 是 extensible＋float ⇒ "这枚缺陷今天在生产里响不响"**仍未定**，`AC#2` 的正确性⛔ 依赖它，但**结案**依赖它）。

## 5. 我这一程的门禁（每把自己落一行 `rc`）

```
rc-d22scan=0        sh scripts/d22scan.sh                     （原始件 logs/hygiene-d22scan.txt，末行逐字 "d22scan: clean"）
rc-build=0          GOFLAGS= go build ./...                    （logs/hygiene-build.txt 全文 11 字节＝我追加的 rc 行，⛔ 0 字节那格）
scope-lines=0       git status --porcelain -- internal cmd     （logs/hygiene-scope.txt；我动过的路径＝只有 .scratch/wisp/probes/300/v1/**）
gitignore-line8=*.out（尺＝sed -n '8p' .gitignore；本程⛔ 一枚 .out 件，正文全 .md、原始输出全 logs/*.txt）
```

**自报我跑过什么**（⛔ 藏着当没跑）：`go version`／`go env GOOS GOARCH GOFLAGS`（派单允许、须自报）；
`go test ./internal/audio/`（导出树，5 发：改前 1 ＋ 突变 4）；`go test ./internal/panel/ ./internal/audio/ ./internal/ball/`（导出树 1 发）；
`go test -c ./cmd/wisp/` ＋ 两形载入探针（**零用例执行**）；`GOFLAGS= go build ./...`；`sh scripts/d22scan.sh`；
`objdump -p`／`sha1sum`／`find`／`grep`／`sed`／`wc` 只读尺若干；`git add`/`commit`（全部带显式 pathspec）。
⛔ 跑：整包 `cmd/wisp` 测试、`go list`、`gofumpt`/格式名册（`AC#4` 落点归落地腿）、真设备那一发。
