# 票 343 交件报告 B 半（读数册）343-r1

正文（A 半）＝同目录 `report.md`；自证册（C 半）＝`report-sec-c.md`。拆件口径＝公共块"每件 ≤160 行"只管我撰写的
正文件、写不下⇒拆具名分件（`322-r2` 先例）、⛔ 压正文。
★逐字读数本体＝已入库的 `logs/01`–`04` 四枚＋`logs/teeth-sweep/` 原始 stdout；
本件只做点名与对回，⛔ 重复贴那四枚里已有的东西（重复＝两份可能互相说谎）。

## 一 起手三读数（⛔ 用派单里编排者写的号＝我的基线＝我起手的 HEAD）
```
git log -1 --format='%H %s'  -> 29d258f9f221ef12e407308ed18fedae31d56ad1
                                主题＝"第三波 b 派单写档＝补放 343-r1…"
git status --porcelain -- internal/risk                 -> 0 行＝车道干净
git log --all --format='%H' --grep='票 343 落地腿 343-r1' -> 0 枚＝自证我是第一枚
git rev-parse HEAD:internal/risk/blacklist.go           -> c1799cd2…（＝派单那句"未动"）
```
★此后 HEAD 被别腿推进（现 `bf5977f4`）；⛔ 影响我的成对＝四数两侧都来自我锚定的
导出树、⛔ 来自共享树 HEAD。

## 二 票面"现量"三行逐枚重跑（全文＝`logs/02`）
| 票面事实 | 我这把尺（同尺形、我自己跑） | 读数 |
|---|---|---|
| 注释在场 `blacklist.go:193` | `git show <锚>:…` 取第 193 行 | 与票面逐字一致 |
| ★`aEquiv` 全仓只这一处且在注释里 | `git grep -n aEquiv HEAD -- internal/risk` | 1 枚命中＝该行注释自己 |
| 同上 | `git grep -c 'func aEquiv' …` | **rc=1＝0 命中**＝从未有定义 |
| 成因＝引入笔 2026-09-21 | `git log -S aEquiv --date=short … \| tail -1` | `f1033e1e 2026-09-21 fix(72): A 表锚点必须和候选路径走同一条句柄管线…` |

★正控（报 0 命中前先跑一枚已知存在的正控）＝同尺形打 `uncertainAnchorMiss`／`aEq`
都 rc=0 命中在场 ⇒ 那枚 0 是尺在工作、⛔ 尺坏。

## 三 AC#0 同族普查＝我自加的两把尺（949 行全量落 `D:\tmp\…\343r1\logs\`）
尺①`census.exe`＝`go/ast` 抽 `internal/risk/*.go`（⛔ `_test.go`）注释面每枚驼峰标识符、
对撞包内名册。名册第一版只有 func/method/type/包级 var·const ⇒ **自抓盲区**：
`bOverrides` 其实是 `Gate` 的形参（真名、注释引用它属实）；第二版补 param／result／
struct field／`:=` 局部／`var` 局部／range 变量。
尺②`callshape.exe`＝尺①的驼峰过滤器**结构性看不见纯小写 Go 标识符**，故第二把尺专挑
注释里写成**调用形**（`foo(`／`foo()`）的名字＝"那是一枚可调用东西"的断言，对撞包内
func/method 名册（402 枚）。
```
blacklist.go 注释面（尺①第二版）＝77 枚 token 行 / 45 枚不同名字
  包内有定义 21 枚、包内无定义 24 枚 ⇒ 逐枚名册与逐枚归类＝`logs/03`（⛔ 在此重贴）
    那 24 枚归类＝散文强调／环境变量名／注册表值名／Win32 API／契约号／票号／平台名
    ⇒ 逐枚读过：⛔ 一枚是"产码断言"，除 ★aEquiv
  ★尺①同族命中＝aEquiv（驼峰形、与 aEq/aUnder/aAnyForm/aRelates/aRuleReason 同词表、
    全仓 0 枚定义）＝本票靶子
blacklist.go 注释面（尺②）＝调用形断言 1 枚 ⇒ ★同族命中＝fires()
```
⇒ ★★AC#0 答案＝`blacklist.go` 注释面同族 **2 枚**（⛔ 一枚）＝`aEquiv`（本票、已改）＋
`fires()`（基线 307 行＝我改后 310 行）。`fires` 作为标识符的真身在**别包**
`internal/config/unwired.go:54`（struct 字段 `fires func(*Config) bool`）；包内
`func …fires`＝0 枚。引入笔＝`e052b11a` 2026-10-17＝票 334 落码腿自己、就在它把
`aRelates` 拆成两把尺那一句里 ⇒ ★⛔ 只交 `aEquiv` 那一枚就走＝票 334 教训已按。
★第二枚⛔ 我动（票面 AC#1＝"那一行注释"＝193；动 310 行＝扩本票射程＝派单⛔ 那条）
⇒ 具名上报（A 半第 4 节）。同一对尺扫**整包**另得 2 枚同族、都⛔ 我写面＝
`provenance.go:54` `taskScopeID`、`:427` `disposalScope`（全仓 0 枚 IDENT；两枚都从未
作标识符定义过）⇒ 逐枚交编排者另立票。

## 四 AC#1 完成判据（逐枚输出；全文＝`logs/02`）
```
git diff --numstat 29d258f9..31d4a476 -- internal/risk/blacklist.go  ->  5  2  (1 枚文件)
git diff -U0 同锚  ->  7 枚 +/-（2 删 5 加）、内容不以 // 开头的 = 0 ⇒ 可执行行变化 0
      （那 7 行逐字在 `logs/02` 第 14-20 行）
★第三把尺（比"看起来只是空格"硬＝公共块 ⓔ 那条要求的形状）＝两侧剥掉注释后的
  可执行字节 sha256：base ⇄ delivered = 97796acc2690b59e…0b379bd3、各 265 行
  SAME_CODE = True
git grep -n aEquiv 31d4a476 -- internal/risk  ->  rc=1 / 0 命中（＝AC#1 判据第①条）
新句子里点名的每枚符号都有定义（逐枚现验、各 1 枚）：aEq:271 aUnder:273 aRelates:280
  anchorMatches:292 uncertainAnchorMiss:362 aRuleReason:327
blob: c1799cd2…（起手，未动）-> 62f2dddb…（交件）；★我的新字节＝LF（全文件 CRLF＝0）。
```

## 五 AC#2 守住（命门）＝可执行字节 0 ＋ 禁动件逐枚等值
```
可执行字节变化 = 0（第四节那把剥注释 sha256 尺＝SAME_CODE True）
⛔ 新增 aEquiv 去让注释变真：交件侧 `git grep -c 'func aEquiv'` 仍 0
四枚受保护仪器 base(29d258f9) ⇄ HEAD 逐枚 blob 等值（现跑）：
  anchor_fallback_teeth_336_windows_test.go  a429cdae == a429cdae  SAME=True
  anchor_fallback_teeth_336_other_test.go    f5f888de == f5f888de  SAME=True
  anchor_miss_reason_339_other_test.go       e7b13c39 == e7b13c39  SAME=True
  anchor_miss_reason_334_test.go             5147c189 == 5147c189  SAME=True
笔数尺见 `report-sec-c.md` 尺①（⛔ 在此抄死笔数＝抄死随每补一笔过期）。逐枚
  `git show --name-only` 取并集实测＝24 枚＝1×.go（blacklist.go）＋1×sh＋18×logs
  （4 枚读数正文件＋14 枚原始 stdout）＋1×票 md＋3×report md ⇒ ⊆ 三处写面、
  ★.go 命中恰 1 枚、⛔ 任何别腿件
  + 2×report md ⇒ .go 命中恰 1 枚＝blacklist.go；⛔ 任何别腿件
`git status --porcelain -- internal/risk` 起手／交件前／交件后都＝0 行＝全程独占
```

## 六 AC#4 成对四数（尺逐字、⛔ 混用⛔ 相减；逐字全文＝`logs/01`＋14 枚原始 stdout）
每发＝`go test ./internal/risk/ -count=1 -timeout 30m -v`，各在自己那棵导出树（八件套
＋docs＋`build/*.dll`＋`third_party/sherpa-onnx/*.dll`＋`frontend/dist`、各 766 枚 `.go`
证非空）。Linux 侧配方逐字抄前一腿那件＝我自己那份 `sh/linux-343.sh`，★⛔ 复用腿的树
＝我自己重导两棵。尺＝`RUN=` `grep -c '^=== RUN'`（含子测试）⇄ `PASS_TOP=` `'^--- PASS'`
（只顶层）⇄ `FAIL_TOP=` `'^--- FAIL'` ⇄ `SKIP=` `'^--- SKIP'`。
```
WINDOWS 基线 RUN=247 PASS_TOP=140 FAIL_TOP=1 SKIP=1 rc=1 ⇄ 交件 同值（247/140/1/1）
LINUX   基线 RUN=195 PASS_TOP=115 FAIL_TOP=1 SKIP=1 rc=1 ⇄ 交件 同值（195/115/1/1）
红名册作差 0、RUN 名集作差 0（247⇄247／195⇄195）；SKIP 两侧同枚（WIN
  TestSyncRegistryProbeLive／LINUX …SyncRootDoesNotDisarmSuspectNet）
台面故障尺（walk／build failed／no test files／setup failed／creating work dir／
  0xc0000135）＝四发全 none ⇒ ⛔ 台损、四数可读。
四发唯一那枚红＝TestC26RewriteAccountIsConsumedAtEverySecurityLeg＝★存量、⛔ 我的新增红
  红句逐字点名 internal/panel 四枚读者（composer.go:127、display_scrub_216.go:26/37/218）
  ＝静态源码审计、读不到本包注释一字；前一腿 339-r1 报里同一枚同名红已登记为存量⇒跨腿对得上
TestResolvePerCallBudget（在册争用件）＝四发全 PASS（1.20／1.19／2.07／1.73 s）
  ⇒ ⛔ 红 ⇒ "安静复量 3 发"那条未触发、⛔ 计入新增红。
跨腿交叉核对＝我这 Windows 四数逐字等于 339-r1 报里那句 247/140/1/1；Linux 我这侧
  195/115/1/1 ⇄ 它交件侧同值 ⇒ 分母跨腿稳定＝我改注释⛔ 动了分母。
```

## 七 ★票 336 那 8 枚牙（AC#4 的保险发；台件＝现成 `run-ac2-sweep.sh` 一字⛔ 动）
env＝`ANCHOR=31d4a476 TREE=D:/tmp/wisp/20261018/343r1/treeDel`（★我自己仓外重导那棵）、
pristine 期望哈希＝我的交件 blob `62f2dddb`。逐腿全文＋红名册＋还原哈希＝`logs/01`
第七节那段与 `logs/teeth-sweep/` 原始 stdout；这里只给结论尺：
```
未突变正控 TEETH rc=0 RUN=17 PASS_TOP=1 FAIL_TOP=0（8 行 operand 全 class=A、
  matched=false clause=true）；LEG rc=0 RUN=10 PASS_TOP=3 FAIL_TOP=0
★8/8 逐枚各自把那行带红＝1a git-credentials｜1b wisp-config-under｜1c 那两枚（case 是
  aUnder AND aEq、两半都挂在支上）｜1d ssh｜1e aws｜1f kube｜1g dpapi-protect｜
  1h local-credentials（1e/1f 的 LEG 也红＝票 334 同名仪器在册、⛔ 我新增）
落刀命中尺（⛔ 空操作伪装成一次成功复跑）＝每腿 needle_before=1 → needle_after=0、
  branch_removals_on_disk=1、changed_lines=2 ⇒ 每发都真落了刀
还原＝11 次全 match=YES（start／after-pristine／每腿一次／final）；
  shared_worktree_dirty_internal_cmd=0＝⛔ 在共享树留字节
★⇒ 我这枚注释字节没顶红任何一枚牙（AC#2 的正向证据）
```

## 八 门禁四数（只在共享工作树；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`；全文＝`logs/04`）
```
GATE1 gofumpt.exe -l internal/risk/blacklist.go  rc=0、列出 0 枚＝我的文件干净；新字节＝LF
      ★具名登记 4 枚存量假脏（⛔ 我动、⛔ 顺手重排）＝internal/agent/approval/
      pending_read.go、internal/agent/tools.go、internal/risk/provenance.go、
      cmd/wisp/models.go；我的文件⛔ 在其中
GATE2 sh scripts/d22scan.sh  rc=0 clean（先正控 runtests.sh、再实扫；八段分母逐枚在
      logs/04：241／42／116／25／39／116／592／143）
GATE3 sh scripts/check-path-length-budget.sh  rc=0 VERDICT GREEN（16258 跟踪路径、
      over-budget 60、roster 覆盖 60、名单外 0）
GATE4 go vet ./internal/risk/  rc=0
```

