# 票 94 独立对抗验收（acceptor-ticket94）

**被验收对象**：`.scratch/wisp/issues/94-winsec-seals-acl-on-an-unresolved-path.md`
（`Status: ready-for-review`，五框自勾）
**代码链**：`dd1e8d3`（判据/checkpoint）→ `7910bcd`（形状落地）→ `20b525d`（AC#3 用例与变异读数）
**验收环境**：本机真机，`go version go1.27.1 windows/amd64`，Windows 11，当前账户
`DESKTOP-LVS7839\swq`（SID `S-1-5-21-1228170099-895614386-1166154857-1001`），icacls = `C:\Windows\System32\icacls.exe`。
本文件由**独立对抗验收代理**撰写；票面与实现者的任何数字都不作为本代理的证据，下面每一条读数都是本代理自己跑的。

**快照纪律（全部在仓外 `/tmp`，仓内无 worktree；工作树生产文件本代理零改动）**：

| 快照目录 | SHA | 用途 |
| --- | --- | --- |
| `/tmp/wisp94acc-94` | `20b525d`（票 94 终态；`diff -rq` 对仓内 `internal/winsec`、`internal/risk` 均**逐文件相同**，故与本票无关的漂移不影响读数） | 门检复跑 + AC#3 变异 + 七项攻击探针 |
| `/tmp/wisp94accbase-94` | `dd1e8d3^`（票 94 动码之前的生产码） | 基线 rc=1 复现 + **成环反证** + 预算测试对照组 |

> 说明：`sh scripts/d22scan.sh` 的 rc=0 读数是在**放入探针之前**、纯净快照上取的；探针包
> （`tmpacc94probe/`、`tmpacc94unwired/`、`tmpacc94order/`、`tmpacc94conv/`、`tmpacc94conv2/`、`tmpacc94b2/`、`tmpacc94f/`）
> 是之后加进快照的，全部只 import `internal/winsec`/`internal/risk`，**不改任何被测文件**，且不存在于仓内。

---

## 进度总览（七项攻击 + 一项对票面论证的反向变异，全部做完；结论在文末 1:1 裁决表）

| # | 内容 | 本代理判定 |
| --- | --- | --- |
| 1 | `ResolvedPath` 是不是唯一铸造口 | **允许残留**（包外三条路全封死=实测；包内伪造**悄悄通过**=实测，票面那句话写过头） |
| 2 | `init()` seam 是不是静默旁路 | **PASS，但确认存在"初始化顺序决定形状"的活证据**（今天无生产码踩到；另发现 `SetPathResolver` 无守卫=见第 1 项） |
| 3 | AC#1 定性重跑 | **PASS（实现者是对的一方，本代理与编排者原断言被推翻）** |
| 4 | AC#3 变异复现 | **PASS（三条红名全复现，SID 级证据本代理自己拿到；但红名里只两条承载判据）** |
| 5 | "内置 verifier 只拒不改写" | **verifier 腿 PASS；C26 腿留一条实测残留 R-a**：过检查却被改写、封的与判的不是同一棵（AC#3 字面仍 PASS，见文末） |
| 5b | 票面那句"未接线即拒会红 20+ 条"（MUTATION-94B） | **为真且保守（实测 49 条）**；附带测出 `internal/winsec` 自身 **0 条红**⇒verifier 正向半句树内无覆盖 |
| 6 | 门检复跑 | **PASS（rc=0 与逐作用域台账 8 行与本票票面逐字相同；ban 文本/allowlist 未动，git 证）** |
| 7 | 两条自认弱处裁决 | **两条都转下张票（必须立案）**：①5b 已把它从"自认"升级为"实测覆盖洞"（verifier 拆掉 winsec 自己 0 条红）；②PROBE F 实测该口能沿 junction 删掉别人的真文件并返回 nil，但 PROBE W 证明今天唯一生产调用方（`removeStray` 的 `WalkDir` 不下降）**够不到** ⇒ 定性为导出 API 地雷；且实现者"要先动 memory"的理由被推翻（修法在 winsec 自己包里） |

---

## 第 1 项：`ResolvedPath` 到底是不是"唯一铸造口"

**先枚举（全仓 grep，`--include=*.go`，从仓根）**：`ResolvedPath` 这个标识符在整个仓库只出现在
`internal/winsec/resolve.go` 与 `internal/winsec/winsec.go` 两个文件里，别处**零命中**；
能产出该类型值的字面量只有三处，全部在 `resolve.go`：

```
internal/winsec/resolve.go:118  return ResolvedPath{}, ...        ← 零值（错误分支）
internal/winsec/resolve.go:121  return ResolvedPath{}, ...        ← 零值（错误分支）
internal/winsec/resolve.go:124  return ResolvedPath{path: p}, nil ← 唯一的带值铸造口
```

消费侧：`grep "ResolvedPath)"` 在全仓**只**匹配到 `privateDirAll(dir ResolvedPath, ...)`
（`winsec.go:159`，**未导出**）与三个方法接收者。⇒ 包外没有任何函数收这个类型。

**构造攻击（四发，全部实测）**：

1. 包外从 `string` 直接转 —— `winsec.ResolvedPath("C:\evil")`：
   `tmpacc94conv\main.go:11:27: cannot convert ... (untyped string constant) to type winsec.ResolvedPath`
   ⇒ **编译不过**。
2. 包外带字段构造 —— `winsec.ResolvedPath{path: "C:\evil"}`：
   `tmpacc94conv2\main.go:10:27: cannot refer to unexported field path in struct literal of type winsec.ResolvedPath`
   ⇒ **编译不过**。
3. 包外零值 —— `var zero winsec.ResolvedPath`：编译过，探针实测
   `String()="" IsZero=true`，且（第 0 段的消费侧枚举）没有任何导出函数吃它
   ⇒ **运行时惰性**，不是旁路。
4. **包内伪造 —— 悄悄通过**。这就是本代理第 4 项变异本身：把 `PrivateDirAll` 的入口写回

   ```go
   abs, err := filepath.Abs(path)   // MUTATION-94
   dir := ResolvedPath{path: abs}   // 同包内直接给未导出字段赋值
   ```

   `go vet ./internal/winsec/` **rc=0**（不是编译失败），测试跑出
   `WRONG TREE SEALED: ...\someone-elses-tree\artifacts went from
   [S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-…-1001] to [S-1-5-18 ×2 …]`。
   ⇒ 类型系统**没有**拦住"绕过铸造口"，拦住它的是 AC#3 那条用例。

   由此 `resolve.go:88-90` 那句 **"filepath.Abs cannot come back without deleting a
   parameter type first" 是错的**：本代理没删任何参数类型，3 行就让它回来了。

5. **另有一条票面没提的**：`SetPathResolver` 是导出的、无 `sync.Once`、无"只许变窄"检查，
   `PathResolverInstalled()` 还把接口原样递出来（读—改—写零成本）。探针 A 在**包外**装一个
   永远返回 `(input, nil)` 的橡皮图章 resolver，然后 `PrivateDirAll(<junction>\artifacts)`：

   ```
   installed by default: risk.c26Pipeline
   PrivateDirAll(...link\artifacts) -> err=<nil>
   victim Everyone grant: before=true after=false      ← 外来读授权被抹掉，且返回成功
   ```

   ⇒ **悄悄通过**（不是编译不过、也不是运行时被拒）。铸造口确实唯一，但它检查的是"装进来的
   那台 resolver"，而"装什么"是全模块可写的。这一条不否决票 94 的 AC（今天没有任何码这么装），
   但"类型上可证明"这个说法的强度必须照实测降级：**它是"链上装了 risk 才可证明"，不是"结构上不可伪造"。**

**判定：允许残留（登记为文末 R-c）。** 修法：`SetPathResolver` 改成 `sync.Once` + 只允许装入"能拒的那一类"
（或干脆把安装口做成 `risk` 侧的 `//go:linkname`/未导出钩子），并把第 4、5 两条写进 `resolve.go` 的注释，
别再说"不可能退回"。

---

## 第 2 项：seam 的 `init()` 是不是一个静默旁路

**先量"谁没 link risk"**（`go list -deps ./<pkg>/ | grep -c internal/risk`）：

| 包 | 是否 link `internal/risk` |
| --- | --- |
| `internal/secret` | **0（不 link）** |
| `internal/memory` | **0（不 link）** |
| `internal/agent` | 1 |
| `cmd/wisp` | 1 |
| `./internal/winsec` 的测试二进制（`-deps -test`） | 1（因为用例 import 了 risk） |

⇒ "未接线"不是假想状态：`internal/secret`、`internal/memory` 两套 suite 和任何只 import
winsec 的生产码，跑的就是内置 verifier。

**实测那条路（探针 D，`/tmp/wisp94acc-94/tmpacc94unwired`，整棵 import 图里没有 risk）**：

```
winsec.PathResolverInstalled() = <nil>          ← 真的没接线（不是测试里 SetPathResolver(nil) 装的）
PrivateDirAll(link\artifacts) err=…refusing to seal…traverses a reparse point at …\data\link   isErrUnresolved=true
PrivateDirAll(link\newdir)    err=同上                                                        isErrUnresolved=true
PrivateFile (link\…\secret.txt) err=同上                                                     isErrUnresolved=true
SealDir     (link\artifacts)  err=同上                                                       isErrUnresolved=true
victim Everyone grant: before=true after=true => untouched (refused)
ok: no directory created through the junction
ok: PrivateFile wrote nothing through the junction
PrivateDirAll(ordinary tree) err=<nil>            ← 没接线也照样能封普通树，所以"未接线即拒"没把世界拒掉
```

⇒ **票面"未接线时一个 junction 输入必须被拒"成立，且四条封存腿（目录/建缺叶/文件/SealDir）都拒**，
不只是 `PrivateDirAll` 一条。

**⚠ 本代理确实找到"靠调用/导入顺序决定安全形状"的证据（探针 E）**：一个 **import winsec 但不 import risk**
的包，在同一个"链上了 risk"的二进制里，**它自己的包初始化阶段看到的是 NIL**：

```
at internal/winsec-only package 'early' init time : NIL -> built-in verifier is what this package's seals get
at main() time (risk linked)                      : risk.c26Pipeline
```

Go 只保证"被依赖的包先初始化"，`early` 与 `risk` 之间没有边 ⇒ 顺序不由规范固定，实测就是
`early` 在前。⇒ **今天不是活缺陷**（`grep -rn "^func init()" --include=*.go internal/ cmd/` 去掉测试后
只有 4 处：三个 llm adapter 注册 + `internal/risk/winsec_c26.go:20`，没有任何包在 init 里封存），
但它把"永远走 C26"这个断言的条件从"链了 risk"细化成"**链了 risk 且不在 risk 之前初始化**"。
票面第 ① 条自认弱处的措辞（"生产二进制永远走 C26"）在这种形状下不完全成立，应补一句。

**判定：PASS（兜底方向正确、拒得响亮）**，附一条必须登记的形状性残留（init 顺序 + `SetPathResolver` 无守卫）。

---

## 第 3 项：AC#1 定性重跑 —— 实现者是对的那一方

**(i) 成环反证（本代理亲测，两重）**

第一重，`go list`：`go list -deps ./internal/risk/` 里**确实**含
`internal/winsec`、`internal/secret`、`internal/observe`、`internal/risk` 自身。
第二重是**编译期判决**——在基线快照 `/tmp/wisp94accbase-94`（票 94 动码之前，`internal/risk` 里还**没有**
`winsec_c26.go`，`git show 7910bcd^:internal/risk/winsec_c26.go` → `fatal: path … exists on disk, but not in`）
给 `internal/winsec/winsec.go` 加一行 `_ "github.com/CarlosShao/wisp/internal/risk"`：

```
package github.com/CarlosShao/wisp/internal/winsec
	imports github.com/CarlosShao/wisp/internal/risk from winsec.go
	imports github.com/CarlosShao/wisp/internal/observe from provenance.go
	imports github.com/CarlosShao/wisp/internal/secret from redact.go
	imports github.com/CarlosShao/wisp/internal/winsec from store.go: import cycle not allowed
```

⇒ **四跳成环，编译器直接拒绝，与票面写的 `risk → observe → secret → winsec` 逐跳一致。**
派单语与票面原文"编排者已核过 `internal/risk` 不 import `internal/winsec` ⇒ 无环"
只对"当时的直接 import"成立；本代理独立确认它是**错的**，实现者推翻得对。
（另注：HEAD 上 `internal/risk` 现在**直接** import `internal/winsec`，就为了装这个 seam ——
所以今天的成环路径其实只剩两跳，任何"让 winsec 反过来 import risk"的想法都会立刻撞墙。）

**(ii) 五个调用点逐个打开（本代理自己的 file:line 读数）**

| 调用点 | 递给 winsec 的路径怎么来的 | 只过词法？ |
| --- | --- | --- |
| `internal/memory/open.go:176` | `abs, err := filepath.Abs(dir)` @ `open.go:169` | ✔ |
| `internal/memory/open.go:188` | `artifactsDir: filepath.Join(abs, artifactsDirName)` @ `open.go:183` | ✔ |
| `internal/memory/open.go:503` | `backupDir: filepath.Join(abs, backupDirName)` @ `open.go:182`（`backupPath` 再 Join 一次） | ✔ |
| `internal/secret/store.go:49` | `s := &Store{dir: filepath.Join(dataDir, "secrets")}` @ `store.go:41` | ✔ |
| `internal/agent/spill.go:111` | `Spiller.dir` ← `filepath.Join(s.dir, name)` @ `spill.go:105` 同一根 | ✔ |

全仓非测试代码里 `risk.Resolve(` **只有 1 处**：`internal/tools/paths.go:56`（另有 `risk.Gate(` 在
`internal/tools/mode.go:98`），两者都不在数据根这条链上 ⇒ 与本票"没有一个落在链上"的断言一致。
而 `internal/memory/open.go` 本身在 allowlist 上（`pathresolver-bypass internal/memory/open.go`，
reason 写着 "migrates to C26 PathResolver in ticket 18"）⇒ 它的"已解析"确实只是**声明**。

**判定：PASS，AC#1 判 (b) 成立。** 四种假绿里的"看图说话"这一项本代理刻意没做：上面每一行都是
`grep -n`/`sed -n` 的原文行号，成环那条是 `go build` 的判决而非 `grep .Imports`。

---

## 第 4 项：AC#3 变异复现

在纯净快照 `/tmp/wisp94acc-94` 里把 `PrivateDirAll` 的入口退回 `filepath.Abs`（同第 1 项攻击 4），
**先证不是编译失败**：`grep -n "MUTATION-94" internal/winsec/winsec.go` → `148:`；`go vet ./internal/winsec/` → **rc=0**。然后本代理自己的读数：

```
$ go test -count=1 -run 'TestAC3JunctionInputIsRefusedNotSealed' -v ./internal/winsec/   TEST_RC=1
    resolve_windows_test.go:92:  PrivateDirAll returned success for a path whose tree lives behind a junction: …\002\data\link\artifacts
    resolve_windows_test.go:92:  WRONG TREE SEALED: …\001\someone-elses-tree\artifacts went from
                                 [S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-…-1001]
                              to [S-1-5-18 S-1-5-18 S-1-5-32-544 S-1-5-32-544 S-1-5-21-…-1001 S-1-5-21-…-1001]
                                 ([NT AUTHORITY\SYSTEM ×2 BUILTIN\Administrators ×2 DESKTOP-LVS7839\swq ×2])
    resolve_windows_test.go:116: PrivateDirAll returned success …（同上，第二次）
    resolve_windows_test.go:116: WRONG TREE SEALED …（同上）
    resolve_windows_test.go:118: the built-in verifier did not name its own refusal
--- FAIL: TestAC3JunctionInputIsRefusedNotSealed (4.07s)
    --- FAIL: …/existing_directory_behind_the_link        (0.10s)
    --- FAIL: …/missing_directory_under_the_link          (0.00s)
    --- FAIL: …/with_only_the_built-in_verifier           (0.08s)
```

⇒ 票面点名的**三条红名逐条复现**，**SID 级证据本代理自己拿到**：旧写法确实把带 `S-1-1-0`(Everyone) 的
外来 DACL 重写成三主体制（Everyone 授权被抹掉），也就是"旧码真的重写过别人目录的 DACL"这句话
**不是修辞，是 icacls 读数**。

**那条证据是不是安全造的（本代理读测试体 + 读自己的运行路径）**：`victim` 来自
`wideParent(t, "someone-elses-tree")`，而 `wideParent` 第一行是 `dir := filepath.Join(t.TempDir(), name)`
（`acl_windows_test.go:157`）⇒ 所谓"别人树"是**Go 测试临时目录里一个被 icacls 显式加了
`*S-1-1-0:(OI)(CI)(RX)` 外来读授权的父目录**，junction 也只建在 `t.TempDir()` 下
（本代理这次运行的真实路径就是 `C:\Users\swq\AppData\Local\Temp\TestAC3…\001\someone-elses-tree\artifacts`），
`os.Mkdir` 的 `victim` 与 `keep-me.txt` 都在里面；测完 `t.TempDir` 自动回收。
⇒ **没有任何一刻真的动过别人/宿主账号的真实树，也没动过仓内 `artifacts`**；本代理自己的探针
（A/D/F）同样全部在 `os.MkdirTemp("", "acc94*")` 里，跑完已 `RemoveAll`（`ls -d $TEMP/acc94*` → 空）。

**⚠ 一条诚实的折扣**：三条红名里 `missing_directory_under_the_link` 在变异体下红得"不承载判据"——
它红是因为 `privateDirAll` 的爬升对 junction 报 `is not a directory`（`errors.Is(err, risk.ErrReparseDenied)`
不成立），而不是因为封错了树；票面自己承认这点（"红在偶然原因上…只当第二条腿，绝不当证明"），
所以真正承重的是 `existing_directory_behind_the_link` 与 `with_only_the_built-in_verifier` 两条。
本代理复核这个自陈**属实**：打印 `WRONG TREE SEALED` 的是 `existing_directory_behind_the_link`（`:92`）
与 `with_only_the_built-in_verifier`（`:116`）两条，`missing_directory_under_the_link` 只报
`no C26 refusal: winsec: …\data\link is not a directory`（红，但红在偶然原因上）；
不额外扣账，但登记一句：**红名数量 3 ≠ 判据数量 2**。

**还原证明**：`cp` 回原文件后
`diff -r --brief internal/winsec/ <仓内 internal/winsec/>` → **rc=0（零输出）**；
`git diff --quiet -- internal/winsec/winsec.go`（仓内工作树）→ **rc=0**；
`grep -c MUTATION-94` → **0**。还原后重跑整包：见第 6 项（44 RUN/44 PASS/0 FAIL/0 SKIP）。

**判定：PASS。**

---

## 第 5 项："内置 verifier 只拒不改写"—— 读码 + 实测

**读码**：`builtinVerifier.Resolve`（`resolve.go:144-171`）在两条拒绝分支之后是
`return platformVerifyPlacement(input)`，而 `platformVerifyPlacement`
（`placement_windows.go:32-71`）的三个拒绝分支之外**只有 `return path, nil` 一条出口**，
`path` 就是形参本身，全函数没有 `Clean`/`Abs`/`Join` 重写 ⇒ 结构上不可能改写。

**实测（探针 C，10 种形状，逐条打印是否 byte-identical）**：

```
  in=…\acc94C786796856                err=false  rewrite=NO (byte-identical)
  in=…\acc94C786796856\               err=false  rewrite=NO (byte-identical)      ← 尾分隔符原样通过
  in=…\acc94C786796856\a..b           err=false  rewrite=NO (byte-identical)      ← 组件内的 ".." 不误伤
  in=…\acc94C…\..\acc94C…             err=true   （拒绝：traverses "..")
  in=…\acc94C…\x                      err=true   （尾空格）
  in=…\acc94C…\y.\                    err=true   （尾点）
  in=\\?\C:\Users\…                   err=true   （扩展前缀）
  in=\\localhost\C$\x                 err=true   （UNC）
  in=…\acc94C…\a\..\b                 err=true   （词法 ..）
  in=C:\\Users\\swq\\…                err=true   （双分隔符）
```

⇒ **"只拒不改写"这句对本机 verifier 成立**，票面上列的三类（reparse / 尾点空格 / `\\?\`）都真的拒。

**但本代理找到一条"既过了 verifier 又被改写了形状"的路 —— 不在 verifier 腿上，在装进去的 C26 腿上。**
`internal/risk/winsec_c26.go:37` 交回的是 `res.Canonical`，而 `risk.Resolve` 的第一步是
`expandInput`（`pathresolver.go:100-121`，展开 `%VAR%`/`$VAR`/前导 `~`）——那是为**分类**设计的，
放进**放置**里就成了改写器。探针 B2（C26 已接线，一个目录名字面量含 `%ACC94TARGET%`，NTFS 合法）：

```
caller's spelling : C:\…\acc94B2…\a%ACC94TARGET%b\artifacts
risk.Resolve      : canonical="C:\\…\\acc94B2…\\aelsewhereb\\artifacts" resolved=false
PrivateDirAll err : <nil>
tree the CALLER named, now existing : false
tree C26 REWROTE to, now existing   : true
caller's dir still Everyone-readable: true
```

⇒ **winsec 封了一棵树、返回成功，而调用方写数据用的是另一棵**（`memory/open.go:181-183` 之后
所有路径都是从**自己那份词法 `abs`** Join 出来的，`PrivateDirAll` 只回 `error`，
所以规范形**没有渠道**回到调用方）。这正是票 94 要关的那类缺陷
（"我以为我在封 A，实际在动 B"），只是入口从 reparse 换成了 **env 展开**。
`~` 前导同理（`~\x` 会被改成 home 下的另一棵树）。
频率上它要求数据根名字里带一对能展开的 `%$`，今天不容易撞上，
所以本代理判它是 **FAIL 级残留但不否决本票的 AC#3**（AC#3 字面只考 junction 输入），
**要求立案**：`winsec` 侧要么在"解析后 ≠ 解析前（词法可比）"时**拒绝放置**，
要么把 canonical 交回调用方（后者要动 memory/secret 签名 = 票 18/79 地界）。

**判定：票面那句话 PASS；但"封的与判的不是同一棵"这个问题只关掉了一半（残留登记为文末 R-a）。**

---

## 第 5 项补充：把票面的**设计论证**也当判据测了一次（MUTATION-94B）

票面 Progress log 与"门禁读数"段各写了一句用来证明"为什么必须有内置 verifier"的话：
**"把未接线的默认做成响亮拒绝会让 `internal/secret`/`internal/memory` 两套 suite 红 20+ 条"**。
这句话是可以伪造的（没人会去测），所以本代理在快照里真的做出来了：把 `resolve.go:110-112` 的
`if r == nil { r = builtinVerifier{} }` 换成 `return ResolvedPath{}, fmt.Errorf("%w: no C26 pipeline linked", ErrUnresolvedPath)`，
`go vet ./internal/winsec/` → **rc=0**（不是编译失败），然后：

```
$ go test -count=1 ./internal/secret/ ./internal/memory/       RC=1
FAIL  github.com/CarlosShao/wisp/internal/secret   0.127s
FAIL  github.com/CarlosShao/wisp/internal/memory  30.188s
--- FAIL 行数 : 49（46 个不同顶层测试名 + 3 条子测试），三条抽样原因均为
    "not provably resolved, refusing to seal: no C26 pipeline linked (MUTATION-94B)"

$ go test -count=1 ./internal/winsec/
ok  github.com/CarlosShao/wisp/internal/winsec   5.184s      ← 全包 0 条红
```

⇒ **两句同时成立**：① 票面那句设计论证**为真且偏保守**（实测 49 条，不是 20+；本代理不采信"约数"，
按自己数的报）；② 更值钱的是第二行——**把整个内置 verifier 拆掉换成响亮拒绝，`internal/winsec`
自己一条都发现不了**，因为 AC#3 的第三条子用例只断 `errors.Is(err, ErrUnresolvedPath)`，
而"响亮拒绝"和"verifier 拒"共用同一个 sentinel ⇒ 那条腿**钉的是"拒"，不是"有一个只拒不改写的 verifier"**。
verifier 的正向半句（"干净绝对路径未接线时原样通过"）在树内**覆盖 0 条**。这把第 7 项①从"实现者自认弱处"
升级成**实测出的覆盖洞**，也顺带说明为什么它只能在 `secret`/`memory` 里补、在 `winsec` 里补不到。

**还原证明**：`git show 20b525d:internal/winsec/{resolve.go,winsec.go,placement_windows.go,resolve_windows_test.go}`
逐个 `diff -q` 对快照 → **四文件全部 == `20b525d`**，`grep -c MUTATION` → 0。
⚠ 另需记明避免误读：验收期间仓内工作树**已被别票推进**
（先是 `acl_windows_test.go` +107、新增 `migrate_windows_test.go` +147，收尾时 `internal/winsec/{winsec_windows.go,reparse_windows_test.go}`
正被 +243/-4 地改——**都不是本代理**：本会话两枚提交 `git show --name-only` 只含本文件）。
本代理所有读数钉在 `20b525d`/`dd1e8d3^` 两枚 SHA 的快照上，**工作树生产文件零改动**；
顺手对本代理没钉的那份**活工作树**扫了一眼 `filepath.Abs|filepath.Clean` 的**调用**（去掉注释行后）→
**零命中** ⇒ 截至本次验收结束，`internal/winsec` 没有把票 94 的红重新引入。

---

## 第 5 项负向对照（PROBE H/H2）：本代理试图证伪"8.3 短名这一类也被 C26 关掉"，**没证成**

票面与 `winsec.go:144-146`、`resolve.go:11-19` 把"8.3 短名"列为这次一并关掉的四类形状之一。
本代理先拿临时目录做，看着像假话：`dir /X` 报出 `AVERYL~1`，而接线与不接线**两条腿都原样通过**
（`canonical=AVERYL~1`，`expanded=false`）。第二发拿真实系统卷做，结论翻回来：

```
in=C:\PROGRA~1       risk{canonical=C:\Program Files  resolved=true}  winsec="C:\Program Files"
in=C:\Program Files  risk{canonical=C:\Program Files  resolved=true}  winsec="C:\Program Files"
in=C:\Windows\Temp   risk{canonical=C:\Windows\Temp   resolved=true}  winsec="C:\Windows\Temp"
```

⇒ **装了 C26 这条腿确实展开 8.3**（临时目录那次是因为该卷的 8.3 生成是"合成式"的，
`dir /X` 显示的短名不是存盘的别名，OS 本来就把它解析到同一个对象）。
另一半也如实记下：**内置 floor 不展开 8.3**（它只能拒，不能改写），但它不改写就**不会封错树**——
本代理实测"沿短名封存到的对象 == 长名那个对象"（`same object: true`）。
⇒ 这一格**不记残留**，只作为"本代理试图打破而没打破"的负向读数留在表里。

---

## 第 6 项：门检复跑（AC#2/AC#4/AC#5）

**基线复现（`/tmp/wisp94accbase-94` = `git archive dd1e8d3^`）**：`sh scripts/d22scan.sh` → **rc=1**，
命中的正是那条真判据（不是注释字形）：

```
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (0.50s)
    scan_test.go:269: repo HEAD violates: internal/winsec/winsec.go:126:
        [pathresolver-bypass] filepath.Abs outside the C26 PathResolver is banned (D22); …
--- FAIL: TestRealRepoLedgerIsHonest (0.48s)
```

**交件纯净快照（`/tmp/wisp94acc-94` = `git archive 20b525d`）本代理自己的读数**：

```
$ time sh scripts/d22scan.sh      real 0m33.154s     rc=0
d22scan.sh: positive control - go test ./... (tools/d22scan)
ok  github.com/CarlosShao/wisp/tools/d22scan  12.357s
d22scan.sh: scan of /tmp/wisp94acc-94
d22scan: examined 217 production Go files under internal/ and cmd/
d22scan: scope bans #1-5 internal/      examined 197 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  37 text files
d22scan: scope ban #7 internal/tools/   examined  17 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  37 text files
d22scan: scope ban #8 internal/         examined 335 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations
```

⇒ **8 行逐作用域台账与票面"门禁读数"段逐字相同**（`bans #1-5 internal/=197 cmd/=20`、`ban #6 frontend/=37`、
`#7=17`、`#8 16/37/335/25`）。`ban #6` 与 `ban #8` 对 `frontend/` 都是 **37**，两本账一致，无需排除规则解释。
⚠ 一处**不是差错的差别**：本代理的基线是 `dd1e8d3^`（比票面那次的 `15c649f` 晚），台账读到
`internal/=194`、`ban #6=37`，与票面基线的 `190/35` 不同 ⇒ 中间是别票进树的文件数变化，与本票无关；
**`ban #6` 只升未降**这条约束在两次读数里都成立（35→37→37）。

**allowlist 与 ban 文本（git 证，不看票面）**：
`grep -vc '^\s*#' tools/d22scan/allowlist.txt` → **5 行非注释**；
`git diff --name-only 7910bcd^ HEAD -- tools/d22scan/allowlist.txt` → **空**；
`git log --oneline dd1e8d3^..20b525d -- tools/d22scan/` → 只有 `5e8f87b`（**票 96** 给 ban #8 的
`emojiScopes()` **补上** `frontend/`，是加严）⇒ 票 94 的两枚 commit 只动 6 个文件
（`git show --stat 7910bcd`/`20b525d`：`internal/risk/winsec_c26.go`、`internal/winsec/{resolve.go,
placement_windows.go, resolve_windows_test.go, winsec.go, winsec_other.go}`），
`tools/d22scan/**` 一字未动，`winsec.go` 也没进 allowlist（它现在根本不含 `filepath.Abs`）。

**AC#5 门禁（同一纯净快照）**：
`gofmt -l internal/winsec internal/risk` → **0 行**；`gofumpt -l` 同 scope → **0 行**；
`go vet ./internal/winsec/` → **rc=0**；
`go test -count=2 -v ./internal/winsec/` → **TEST_RC=0**，`=== RUN` **44**、`--- PASS` **44**、
`--- FAIL` **0**、`--- SKIP` **0**，不同名测试 **22** 个（44 = 2×22 ✔）。
⇒ 票面"44/0/0 + 12 顶层 + 10 子测试"被独立复现。

**连带回归（本代理在快照里跑的，票面没要求本代理复算但值得钉）**：
`go test -count=1 ./internal/secret/ ./internal/memory/ ./internal/agent/ ./internal/observe/` →
**四包全 ok**（1.326s / 14.289s / 3.322s / 1.460s）。⚠ 同一次 `./internal/risk/` **红了一条**：
`TestResolvePerCallBudget … 1.508534 ms/op (budget 1.000 ms, 702 samples)`。
本代理没有接受票面"并发争用、静默后三次全 ok"的归因，自己取了**多样本**（静默、逐条全报）：

| 树 | 静默单跑样本 | 结果 |
| --- | --- | --- |
| `20b525d`（票 94） | 7 次 | **6 ok / 1 FAIL**（那条 FAIL 就是 1.508534 ms/op 那次） |
| `dd1e8d3^`（基线） | 4 次 | 4 ok / 0 FAIL |
| 与 winsec `-count=2` 并发 | 1 次 | FAIL 1.509 ms/op |

⇒ 结论：这条**不是票 94 的回归**（票 94 没碰 `pathresolver*.go`，git 已证），但它也**不是纯并发争用**——
预算 1.000 ms 与实测 1.0~1.5 ms 只差一丝，静默也会偶红。票面写"静默后三次全 ok"
属于**多样本没全报**（见下面假绿第 4 项），应改成"这台机器上这条预算贴着边界，红/绿不可复现"，
并交给仪器账（票 88 已经立过"会回放的自检脚本"那一类）。

**判定：PASS（AC#2/AC#4/AC#5）**：另补一条"本代理试图把门弄窄但没找到"的负向读数——
`git diff dd1e8d3^ HEAD -- tools/d22scan/scan_test.go` 的**删除行只有 5 条，且无一条含
`abs`/`pathresolver`** ⇒ 那条票 96 的改动没有顺带削掉任何断言。

---

## 四种假绿逐条点名（本票范围）

1. **"没跑"冒充"绿"** —— 未发现。本代理所有 rc 都单独打印（`TEST_RC`/`BUILD_RC`/`RC=`），
   并且 `-count=2` 的 44 RUN 与 22 个不同名做了乘法核对；`d22scan.sh` 自带正向对照
   （`go test ./...` 先跑 seeded-violation 自扫），且本代理**同一脚本在基线快照上拿到了 rc=1**
   ⇒ "clean" 不是门瞎了。
2. **"改了判据让它绿"** —— 未发现。`tools/d22scan/**`、`allowlist.txt`、`pathresolver*.go` 用 git 证未动；
   AC#3 的三条子用例断言的是"拒 + 外来 DACL 一字不变"，变异能红，说明判据没被写歪。
   ⚠ 但有一条**同类倾向**：`SetPathResolver` 无守卫（第 1 项攻击 5）+ AC#3 里
   `installed := winsec.PathResolverInstalled(); … t.Cleanup(restore)` 这种"测试自己换开关"的形状，
   将来谁把生产接线挪走而测试仍绿，是可能的。
3. **"看图说话"（没测就写）** —— 票面 AC#1 的 (b) 本代理逐行复算过，**全部属实**；
   "旧写法**确实**重写了别人的 DACL"这句，票面是拿测试打印当证据的，本代理**独立重跑变异**才确认；
   票面"未接线即拒 ⇒ `secret`/`memory` 会红 20+ 条"这一条本代理**也真去做了那个变异**
   （MUTATION-94B，见第 5 项补充）：**方向为真、量级保守**（实测 49 条），
   同时暴露出"`internal/winsec` 自己 0 条红"这个覆盖洞 ⇒ 这条**不再记为未证**。
   本代理唯一仍未证的票面断言是 `cmd/wisp/doctor.go:222-245`/`internal/proc/envfork.go:55-136`
   那两条**来源链**（`os.UserConfigDir()`/`PortableDataDirName` 全只过 `filepath.Join`）——
   它们只影响 AC#1 的"来源有多宽"，不影响"链上零解析"的判决（那由 5 个调用点 + 全仓 `risk.Resolve(` 计数支撑）。
4. **"多样本挑运气"** —— **命中一条**：`TestResolvePerCallBudget` 票面报"静默三次全 ok"，
   本代理静默 7 样本 6 ok/1 FAIL（数字见第 6 项）。这条不改 AC 的绿，但它把一条**边界性 flake**
   写成了"争用所致、已排除"，归因过头。

---

## 第 7 项：两条自认弱处裁决

**① 内置 verifier 腿的 CI 覆盖** —— 现状：全仓 `ErrUnresolvedPath` 只出现在 `internal/winsec/` 内
（`grep -rn "ErrUnresolvedPath" --include=*.go internal/ | grep -v internal/winsec` → **No matches**）
⇒ 这条腿的覆盖**只有** `resolve_windows_test.go:112-120` 那一腿，而那是**运行时把开关拔了**，
不等于"二进制没 link risk"这一真实状态（本代理的 PROBE D 才是，但它在临时目录里、没进树）。
**并且本代理实测：把 verifier 整个换成响亮拒绝，`internal/winsec` 0 条红（MUTATION-94B，第 5 项补充）
⇒ 它现在钉住的是"拒"，不是"存在一个只拒不改写的 verifier"，正向半句树内覆盖 = 0 条。**
**裁决：一条没覆盖满的 AC 缺口（登记为文末 R-d），不阻塞本票 `-done`，但必须转下张票**：
在 `internal/secret` 或 `internal/memory`
（两者经实测**不 link risk**）加两条用例：
(i) 数据根拼成 junction ⇒ `PrivateDirAll` 以 `winsec.ErrUnresolvedPath` 拒；
(ii) 干净的绝对路径未接线时**能封成**且 `String()` 与输入**逐字节相同**（这才是钉住"只拒不改写"的那一半）。
理由：行为本身已被本代理证对（PROBE D 四腿全拒、没建目录、外来 DACL 不变；PROBE C 十种形状 byte-identical），
缺的只是"别人以后能重跑"。

**② `RemoveUnlinked` 刻意不解析** —— **本代理实测它比票面写的更糟，但可达性也被本代理自己限定住了。**
探针 F（不 link risk 的二进制）：路径 `<data>\link\artifacts\keep-me.txt`，其中 `link` 是 junction，
末组件是**真文件不是链接**：

```
RemoveUnlinked(path reached THROUGH the junction) err = <nil>
foreign file still present afterwards: false      ← 别人树里的文件被删掉了
junction itself still present: true
```

⇒ 它的守卫只有"**末组件**是链接就别跟"，**没有任何一条"路径中途穿过链接就不许动手"**的检查。
票面说"要收这个口得先动 `internal/memory`（票 18/79 的地界）"——**这句被推翻**：最少只需在动手前对
**除末组件外**的祖先链跑一次现成的 `platformVerifyPlacement`/已装的 C26（正是 `placement_windows.go`
里那段逐级 `isReparsePoint` 的循环，去掉最后一级），**一行生产码地界都不越**。

**但本代理没有停在这句上，而是去查了唯一生产调用方到底够不够得到（探针 W）**：
`internal/memory/artifacts.go:153` 的 `removeStray` 用 `filepath.WalkDir` 量自己的子树，再
`filepath.Join(s.artifactsDir, rel)` 重建路径去删。实测那份 walk **看见 junction 时报的是
`isdir=false type=?---------`，并且不下降**：

```
. (isdir=true)   junklink (isdir=false type=?---------)   stray (isdir=true)   stray\ours.txt (isdir=false)
walk descended into the junction: false
foreign victim still present: true
```

⇒ 生产码**今天构造不出**"中途穿过链接的 `full`"：链接要么被 walk 当成一个文件条目交回（那时
`RemoveUnlinked` 正好走"只解链不进入"的设计分支），要么根本进不了 `rel`。
所以这一格的准确定性是：**导出 API 上的地雷（下一个调用方一次 `Join` 就能踩，且失败方向是删别人的东西、返回 nil），
不是今天可达的活缺陷。**
**裁决：转下张票**——但**必须立案**，且立案时把票面"要先动 memory"那句错账改掉（修法在 `internal/winsec` 自己包里，
红名就用本代理这条 PROBE F）。

---

## 最终裁决表（与票面 AC 1:1）

| AC | 票面声称 | 本代理独立读数（原文摘要） | 裁决 |
| --- | --- | --- | --- |
| **AC#1** | 判 (b)：四个调用点只过 `filepath.Join/Abs`；并推翻"winsec import risk 无环" | `go build` 在基线快照给出四跳 `import cycle not allowed`；5 个调用点 file:line 逐条复算为词法；非测试 `risk.Resolve(` 全仓仅 `internal/tools/paths.go:56` | **PASS** |
| **AC#2** | 纯净快照 `git archive 20b525d` d22scan rc=0，台账 `ban #6` 不减；allowlist 未动 | rc=0（33.2s），台账 8 行逐字复现（`internal/=197`、`#6 frontend/=37`、`#7=17`、`#8 16/37/335/25`）；`grep -vc '^\s*#' allowlist.txt`=5，`git diff --name-only 7910bcd^ HEAD -- allowlist.txt` 空 | **PASS** |
| **AC#3** | junction 输入"要么拒要么封到解析后那棵"；变异退回 `filepath.Abs` 三条红名 + SID 证据 | 变异 `vet rc=0`、三条子用例全红、`WRONG TREE SEALED … from [S-1-1-0 …] to […×2]` 本代理自己拿到；未接线二进制四腿 `ErrUnresolvedPath` 拒且外来 DACL 不变；还原 `diff -r` rc=0。**折扣**：红名 3 条里 1 条红在偶然原因；另发现 C26 腿的 `%VAR%` 改写会"封 A 判 B"（PROBE B2，err=nil），此路 AC#3 未覆盖 | **PASS（附一条须立案的同缺陷类残留，见第 5 项）** |
| **AC#4** | "收尾前必须跑 `sh scripts/d22scan.sh`"写进票面并实测跑得起来（给时长与 rc） | 固定动作在本票 `## 门禁（AC#4）` 段（第 85-102 行，`## Rules` 之前）；本代理自测时长 **33.2s / rc=0**（纯净快照）与 **基线 rc=1** ⇒ 脚本真的会红也真的能绿；
`git diff --name-only 7910bcd^ HEAD -- .github/workflows/ci.yml` → **空**（未抢票 85 的地界） | **PASS** |
| **AC#5** | gofmt/gofumpt 0 行、vet rc=0、`-count=2 ./internal/winsec/` rc=0 并逐条点名 SKIP/FAIL | 0 行 / 0 行 / rc=0 / **44 RUN·44 PASS·0 FAIL·0 SKIP（22 不同名）**；四包连带回归 ok；`internal/risk` 预算测试 6ok/1FAIL（非本票回归，归因被票面写过满） | **PASS（附仪器账一条）** |

---

## 两句必须直答的结论

### 1) 票 94 能否挂 `-done`？—— **能挂，但是条件式挂：四条残留先立案、两句过头文案先降级。**

先说清楚判据从哪来往哪走：五枚 AC 本代理**逐条独立复现且全 PASS**（上表），
门检 rc=0 的 8 行台账、`-count=2` 的 44/44/0/0、变异三条红名与 SID 级证据、AC#1 对编排者断言的推翻，
没有一条是靠票面自报的数 ⇒ **这不是票 89 那种"账不实"的 `returned-for-fix`**。
本代理**中途改过一次判**：起初把第 7 项②（`RemoveUnlinked`）当成"本票地界内没做完"从而压住 `-done`，
随后自己做了探针 W，证明今天唯一生产调用方的 `filepath.WalkDir` **不下降进 junction**、
构造不出那个删除路径 ⇒ 那条的准确定性是**导出 API 地雷**而非活缺陷，压 `-done` 不成比例。

挂 `-done` 的四个条件（都只要求在账上，不要求在本票内改实现）：

| 残留 | 本代理实测锚点 | 该谁收 |
| --- | --- | --- |
| R-a `expandInput` 把放置改写到另一棵树（含 `%VAR%`/前导 `~` 的数据根），封一棵、返回 nil、调用方写另一棵 | PROBE B2（三行打印：canonical 变了、caller 那棵 `exists=false`、err=nil） | **新票（安全类，优先级不低）**；AC#3 字面它 PASS（封的确实是 resolver 给的那棵），没覆盖的是"canonical 与调用方拼写指向两个对象"这种情形 |
| R-b `RemoveUnlinked` 不校验祖先链，可沿 junction 删别人树里的真文件并返回 nil | PROBE F（`err=<nil>`、`foreign file … present: false`）；可达性由 PROBE W 限定为"今天不可达" | **新票**，立案时改掉票面"要先动 memory"那句（修法在 `internal/winsec` 自己包里） |
| R-c `SetPathResolver` 无 `sync.Once`、无"只许变窄"检查；包外装橡皮图章后 `PrivateDirAll` 静默重写外来 DACL | PROBE A（`err=<nil>`、victim 的 `S-1-1-0` 由有到无） | 小改（本票地界即可）或并入 R-b 那张票 |
| R-d 内置 verifier 的**正向半句**（干净绝对路径原样通过）树内覆盖 0 条 | MUTATION-94B：把 verifier 整个换成响亮拒绝，`go test ./internal/winsec/` **ok**，而 `secret`+`memory` 红 **49 条** | **转下张票**（两条用例，写在第 7 项①） |

外加一条文案账：`resolve.go:88-90` 的 "filepath.Abs cannot come back without deleting a parameter
type first" 与本票 Progress log 的"生产二进制永远走 C26"两句都**比实测强**（前者被本代理的变异直接打穿、
后者被 PROBE E 的 init 顺序打穿），实现者补一条 append-only 的 Progress log 把话说准即可，不改判据。

⚠ **不要把票 89 的账混进这一格**：票 89 现在是 `returned-for-fix`
（`.scratch/wisp/issues/89-*.md:3` 自记"AC#2 一格 FAIL，其余五格 PASS"，锚点是"删掉 propagatePrivate
后 34 条全绿"与"唯一那份明文产物没封"）。那是**另一张票的独立欠账**，既不因票 94 落地而消除，
也不该算到上面 R-a..R-d 任何一条上；两票现在都在改 `internal/winsec/`，**并树时先后顺序要编排者定**。

### 2) "封私有数据时用的路径现在是解析过的"—— **能有条件地说；按票面那句原话说是不行的。**

**能说的部分（本代理实测）**：不管进程有没有链上 `internal/risk`
（实测 `internal/secret`、`internal/memory` **不链**，`internal/agent`、`cmd/wisp` **链**），
五个封存入口 `PrivateDirAll`/`SealDir`/`SealFile`/`PrivateFile`/`PrivateFileExclusive`
现在都会在动手**之前**先过 C26 或内置 verifier：一个中途带 junction 的输入
**四条腿全部拒、不在链接下建目录、不写文件、外来 DACL 一字不变，而同一棵树不带链接时照样封得成**
（PROBE D + AC#3 变异对照）。这条覆盖面比"把 `winsec.go:126` 那一行换掉"宽，是票 94 的真实增量。

**必须同时说的三个限定（都是实测，不是假想）**：
(i) **只保"封存"这一族**：`RemoveUnlinked` 仍不解析（R-b），今天够不到，但它是同一个包里导出的口；
(ii) **"解析过"不等于"封的就是调用方要写的那棵"**：装了 C26 时 `expandInput` 会改写含 `%VAR%`/前导 `~`
的拼写，winsec 封改写后那棵并返回 nil，而 canonical 没有任何渠道回到调用方
（`memory/open.go:181-183` 之后全用自己的词法 `abs`）⇒ 见 R-a；
(iii) **它是"装了才解析"的保证，不是结构保证**：同模块任何包能把 seam 换成橡皮图章（R-c），
包内 3 行就能伪造 `ResolvedPath`（本代理的变异），所以"下次不可能退回"这种话不能对 owner 讲。

⇒ 给 owner 的准确句子是：
**"决定'给哪棵树封权限'的那一步，今天不会再拿一条只做词法拼接的路径去动别人的 DACL 了——
带链接的形状一律响亮拒绝，这一条我们能在真机上调出来给你看；
但'解析过的路径'与'调用方真正写进去的那棵树'是不是同一棵，还差 `%VAR%` 展开那一格（R-a），
删除那条腿也还留着（R-b），这两格已经立案，没算在'已完成'里。"**
