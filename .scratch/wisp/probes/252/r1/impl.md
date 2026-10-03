# 票 252 · 实现腿 `252-r1` · 证据件（AC#2 修法 + AC#3 载具）

本腿只做 **AC#2 与 AC#3**。AC#4／AC#5／AC#6 一枚未做，具名原因写在 §6 甲-1。
判语归非实现者；本件只交读数，并且每一枚读数都自带它当时的 HEAD。

---

## §0 起手锚

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-03 09:13:00 +08:00`（`date -Iseconds` 逐字 `2026-10-03T09:13:00+08:00`） |
| `git log -1 --format=%h` | `3736f0dd` |
| `git log -1 --format=%H` | `3736f0dd5de1787fc401ce52138f629c90b14900` |
| `git log -1 --format=%s` | `docs(probes/258-a2): skeleton + section 0 anchor (r...` |
| 全仓 porcelain 行数 | `398` |
| 本腿写面（`git status --porcelain -- cmd/wisp internal/tools internal/risk`） | `0`（起手即空） |
| 分支 | `dev` |
| `internal/tools/paths.go` 起手版本 | `8e100950`（该文件最后一次提交），工作区 md5 `6e169604f35f756a5e50f8b21f215d2a`＝备份 `.scratch/wisp/probes/252/r1/backup/paths.go.orig` |

⚠ 起手锚之后 HEAD 因别的腿持续推进（本腿自己的三枚提交：`bbbb3593`→`34d89191`→`2b1a3071`；交件时仓库 HEAD 已是别人的 `598cc093`）。
过程中 `cmd/wisp` 由干净转为 4 枚脏件（`config_reload.go`／`config_reload_223_test.go`／`config_readers_255.go`／`config_receipt_255_test.go`），
`internal/config` 出现一枚他人新增的 `tiers_app_255_257-261_test.go`——**都不是本腿写的**，本腿写面全程保持空。

本腿提交链与当时的 HEAD：

| 小格 | commit | 内容 |
|---|---|---|
| 骨架＋§0 | `bbbb3593` | 本件 §0 |
| AC#3 载具（改前必红） | `34d89191`（`test(tools/paths): ticket 252 AC#3 carrier, red before the fix`） | 新判据文件＋修前逐字读数 |
| AC#2 修法（改后必绿） | `2b1a3071` | `internal/tools/paths.go`＋修后读数三份 |

---

## §1 现量

### 1-甲 AC#1 那三行逐字复证（本腿自己跑，不抄票面）

复证方式有两发，互相咬合：
(i) 本腿的新判据在**未修的生产码**上直接跑（`.scratch/wisp/probes/252/r1/pre-fix-run.log`）；
(ii) 修完之后种一发 M1 突变（把两侧同形那一步拿掉）再跑 252-p1 探针（`.scratch/wisp/probes/252/r1/mutation-M1.log`），
其形状＝修前形状，于是 AC#1 那三枚正控与 Q1 那一行都由本腿在同一台机器上重取。

**(1) 分岔那一行（roots 长名／asked 短名＋叶子缺失）**，`pre-fix-run.log` 逐字（A1 枚，与本仓 `q252p1-never-created.txt` 那一发在 M1 下发里逐字一致）：

```
Roots()                        : ["d:\\work\\workspace\\projects plans\\wisp"]
foldPath(canonical)            : "d:\\work\\worksp~1\\projec~1\\wisp\\q252r1-new-note.md"
LEG1 rootsContain(roots,folded): false
resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252r1-new-note.md" ok=true
foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252r1-new-note.md"
LEG2 rootsContain(roots,foldRF): true
InAllowlist(canonical)         : false
FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
```

⇒ **复证为真：先失配的是词法那一段（现 `paths.go:196`），`resolvedForm` 那一段（现 `:208-209`）反而宽容。**
票面 §12「这两次包含就不可能同时成立」被 `252-p1` 推翻、本腿再复证一次：**一次成立一次不成立，两段对同一条物理路径给相反答复。**

**(2) AC#1 那三枚正控（M1 发＝修前形状，同一台机器）**，`mutation-M1.log` 逐字：

```
CASE P1 roots=LONG asked=LONG+missing-leaf    InAllowlist(canonical): true
CASE P2 roots=LONG asked=SHORT of an EXISTING dir   InAllowlist(canonical): true
CASE P3 roots=SHORT asked=SHORT of an EXISTING dir  InAllowlist(canonical): true
```

同发里 `Q3b roots=SHORT asked=SHORT+missing-leaf` 仍是 `InAllowlist=false`、仍响在第一段
⇒ **复证：两侧"配置"同形不是正控，因为存在性把 roots 折成长名、缺失叶子留原拼法。**

**(3) 分岔变量是存在性不是拼法**——本腿在同一台机器上再加一发直接读数（§4 表 A2/A3 行）：
roots 配短名时 `Roots()` 仍是 `["d:\\work\\workspace\\projects plans\\wisp"]`（长），
asked 的短名只在**叶子缺失**时才活到比较界面（`internal/risk/pathresolver.go:138` 的 `res.Canonical = p`）。

### 1-乙 本腿改的那两行，前后对照

| | 修前（HEAD `3736f0dd`，文件 `8e100950` 版） | 修后（`2b1a3071`） |
|---|---|---|
| `internal/tools/paths.go:127` | `res, err := p.resolve(raw)` 之后直接 `return res.Canonical, nil`——**C26 没解析出来时，调用方的拼法原样进比较界面** | 同一位置起为 252 的说明块（`:127-146`），`:144-146` 是 `if res.Resolved { return res.Canonical, nil }` |
| `internal/tools/paths.go:148` | （不存在） | `return sameFormOfUnresolved(res.Canonical), nil`——**只有 C26 未能确认的那一发走补齐** |
| `internal/tools/paths.go:151-186` | （不存在） | 新函数 `sameFormOfUnresolved`：平台常量守卫＋存在性守卫＋UNC/扩展前缀守卫，最后调**已有的** `resolvedForm`（`:303`） |
| `internal/tools/paths.go:189-219`（`InAllowlist`） | 两次包含＋`!ok` 那一支 | **逐字节未动**（见 §2 的 md5 证据） |

同一发输入的两行前后差（A1 枚）：

```
修前 Canonicalize  = D:\work\WORKSP~1\PROJEC~1\Wisp\q252r1-new-note.md   -> InAllowlist false（响在 :196）
修后 Canonicalize  = D:\work\workspace\projects plans\Wisp\q252r1-new-note.md -> InAllowlist true
```

---

## §2 修法（AC#2）

**落点：`internal/tools/paths.go:119-149`（`Canonicalize` 尾部）＋ `:151-186`（`sameFormOfUnresolved`）。**
即票面裁定那句「**进 `InAllowlist` 之前那一步**补齐成两侧同形」的字面位置：`Canonicalize` 是 `InAllowlist` 前一步，
生产里三枚判定方（`internal/risk/rules_gateway.go:45`、`internal/tools/bridge.go:1145`、`internal/tools/task.go:838`）
都是先 `Canonicalize` 再 `InAllowlist`，所以补齐后的形就是判定实际吃到的形。

**为什么这里就够**：C26 的句柄查询（展开 8.3 的那一步）只对存在的路径回答（`internal/risk/pathresolver.go:127-138`），
所以 `allowed_dirs` 的根（必然存在）进册时已被折成长名，而"要新建的文件"却带着调用方原拼法进比较界面。
把缺失叶子那一发也送到**同一条已存在的折叠**上过一遍（存在的部分取 OS 自己的拼法、缺失的尾巴逐字重接），
两侧就同形了，第一段与第二段于是各自作差同一个形。

**没有新造规范化器**（票面硬禁）：
- 折叠本身一枚都没新写——`sameFormOfUnresolved` 里唯一的展开动作是调用本文件已有的 `resolvedForm`（`paths.go:303`，即 `252-p1` 量到"给缺失叶子返回长名且 `ok=true`"那一枚），没有复制粘贴它的走链逻辑；
- 也没有引入 `GetShortPathNameW`／手工打 `~1`／字符串近似：`internal/risk/pathresolver.go:343-344` 立着一条本仓规矩「using the C26 handle pipeline only: no string heuristics, no case or short-name "alignment"」，本修法完全遵守；
- 平台分支是编译期常量（`filepath.Separator != '\\'`），与本文件 `unifySeparators`（`:349` 附近）同一形状，POSIX 侧逐字节返回原串；
- 前置守卫只在 `res.Resolved == false` **且** `os.Lstat` 报 `ErrNotExist` 时才动字符串：已存在的路径（含 ACL 挡住句柄查询那一发）保持 C26 原答复，本票射程只有"新建"。
- 同一个形状本仓已写过一遍：`internal/risk/syncdirs.go:189-196`（票 124 的 `resolveTarget`，"resolve the nearest EXISTING ancestor and re-attach the lexical remainder"），本修法沿用其形，未再发明。

**两枚硬禁的落地证据（不是口头）**：

1. 「绝不许把两次包含合成一次」＋「绝不许把解析不到＝未授权改软」：`InAllowlist` 本体**一字节未改**。
   `git diff -U0 -- internal/tools/paths.go` 只有一枚 hunk：`@@ -127 +127,57 @@ func (p *PathCanonicalizer) Canonicalize(...)`；
   并且把修前 133-163 行与修后 189-219 行（同一段体）分别 md5：两边都是 **`e4aac2d53fe6bca6f43e5505b03e4202`**
   ⇒ 两次包含（现 `:196` 与 `:209`）与 `!ok` 那一支原样保留，只是整体位移 56 行。
2. 「不许顺手做成硬边界／不许放宽任何一层的拒绝」：`Roots()` 的内容与构造路径（`paths.go:93`）没动；
   判级侧（`internal/risk/rules_gateway.go`、R2→L2）零改动；§4 的 N1–N7 七枚"必绿且不许变红"的负控逐枚给了同机读数。

**已知残余不对称（不掩盖）**：补齐发生在 `Canonicalize`，不在 `InAllowlist` 内。
因此**绕过 `Canonicalize` 直接把短名串塞给 `InAllowlist`** 的调用（本仓只有测试这么干，生产判定方都是先 `Canonicalize`）仍会响在第一段。
为什么不许把它挪进 `InAllowlist`：那样第一段就得吃 `resolvedForm` 的形，而 `resolvedForm` 的形正是第二段的形——
两次包含就会**塌成一次**（这正是票面硬禁那一枚），且"词法在根外、解析后落回根内"那一发会被新授权（放宽）。
本腿选择保住两段的独立性，把不对称留在"必须经 C26 入口"这条既有契约上（`paths.go:14-18` 的自陈：判定、执行、上报同一条树）。

---

## §3 门禁四数真实读数

| 门禁 | 命令 | 读数 |
|---|---|---|
| 全仓构建 | `GOFLAGS= go build ./...` | `BUILD_EXIT=0`，无输出（`.scratch/wisp/probes/252/r1/gate-build.log`） |
| 两包测试（只跑本腿写面） | `GOFLAGS= go test ./internal/tools/ ./internal/risk/ -count=1` | 终态：`ok github.com/CarlosShao/wisp/internal/tools 14.232s` ／ `ok github.com/CarlosShao/wisp/internal/risk 4.273s`，`--- FAIL` 枚数 `0`（`.scratch/wisp/probes/252/r1/final-packages.log`）；起手基线同为两枚 `ok`（`baseline-gotest.txt`） |
| 新判据逐枚 | `go test -count=1 -v -run 'TestTicket252R1' ./internal/tools/` | 终态 `--- PASS` **18** 枚／`--- FAIL` **0** 枚／`SKIP` **0** 枚（`final-r1-run.log`） |
| 252-p1 探针复跑 | `go test -count=1 -v -run 'TestTicket252P1' ./internal/tools/` | `--- PASS` **9**／`--- FAIL` **0**／SKIP **0**（`post-fix-p1-probe.log`） |
| gofumpt | `$(go env GOPATH)/bin/gofumpt.exe -l internal/tools internal/risk` | 输出为空（首跑曾报 `paths_shortname_252_r1_test.go` 一枚，已 `-w` 修掉后复跑为空） |
| d22scan | `./tools/d22scan/d22scan.exe` | `D22_EXIT=0`；末行逐字：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=227, bans #1-5 cmd/=37, ban #6 frontend/=85, ban #7 internal/tools/=23, ...`（`gate-d22scan-final.log`；改动前另有一发 `gate-d22scan.log` 同判语） |

⚠ d22scan 输出里有一枚 `skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`
——那是仪器自己的 git-ignore 台账行（本腿未新增任何被忽略文件，改动前后同串），**不是本腿的测试 SKIP**；
按派单"SKIP 算红"的口径，本腿把它的逐字读数摊在这里而不是替它解释。
⛔ 本机 staticcheck 产假绿，未跑。

---

## §4 改前必红／改后必绿同机对照表 ＋ 变异自证表

### 4-甲 同一台机器、同一载具、同一枚枚名

| 用例（`internal/tools/paths_shortname_252_r1_test.go`） | 修前（`pre-fix-run.log`） | 修后（`post-fix-run.log`／终态 `final-r1-run.log`） | 判语 |
|---|---|---|---|
| `TestTicket252R1ShortSpellingOfNewFileIsAuthorized/A1_roots=LONG_asked=SHORT+missing-leaf` | `--- FAIL`（`:98` `POSITIVE CONTROL RED ... was refused`） | `--- PASS` | AC#2 主格，AC#1 那一行 |
| `.../A3_roots=SHORT_asked=SHORT+missing-leaf` | `--- FAIL`（同上，`:98`） | `--- PASS` | 两侧配置同形仍须同答 |
| `.../A2_roots=SHORT_asked=LONG+missing-leaf` | `--- PASS` | `--- PASS` | 反方向正控，不许被打红 |
| `.../A4_roots=LONG_asked=LONG+missing-leaf` | `--- PASS` | `--- PASS` | 同形正控，不许被打红 |
| `TestTicket252R1BothSpellingsAnswerTheSame/{Wisp,docs,tools}` 三枚 | `--- FAIL`（`:131` `AC#2 RED: InAllowlist gives two answers to one physical path`） | `--- PASS` ×3 | AC#2 那句直译：一种物理路径两种拼法一个答复 |
| `TestTicket252R1AlignmentAddsNoAuthorization/{N1..N7}` 七枚 | `--- PASS` ×7 | `--- PASS` ×7 | 负控：修前修后都必须拒绝（放宽即红） |
| `TestTicket252R1AlignmentAddsNoAuthorization`（母格尾判） | `--- FAIL`（`:252`） | `--- PASS` | 短名新文件在自身根下必须被授权 |
| `TestTicket252R1CanonicalStillNamesOneTree` | `--- FAIL`（`:270` 两串不等） | `--- PASS` | 规范化后同名同树、叶子名逐字保住 |
| 合计 | `--- FAIL` 9／`--- PASS` 9／SKIP 0 | `--- FAIL` 0／`--- PASS` 18／SKIP 0 | ⛔ 本件 0 枚 `t.Skip`；未挂〔仅 CI 可量〕 |

票面 AC#1 那三种正控**没有被本修法打红**：`TestTicket252P1AllowlistLegs` 的 `P1/P2/P3` 与
`TestTicket252R1.../A2/A4`、以及 `internal/tools/paths_ticket107b_probes_test.go` 的 A/B/C 三枚探针在两包终态跑里全部 `--- PASS`。

### 4-乙 变异自证（每发前记 md5、跑完立刻还原、两发之间量面）

备份：`git show HEAD:internal/tools/paths.go > .scratch/wisp/probes/252/r1/backup/paths.go.orig`，md5 `6e169604f35f756a5e50f8b21f215d2a`（＝起手工作区）。
还原一律 `git cat-file blob HEAD:internal/tools/paths.go > internal/tools/paths.go`（⛔ 未删除任何文件）。

| 号 | 种的形 | 改动（file:line） | 指名新用例读数 | 正控（还原后） |
|---|---|---|---|---|
| **M1** | **两侧同形那一步被拿掉**：`sameFormOfUnresolved` 末尾 `return rf` 改成 `return canonical`（`paths.go:185`；`rf`/`ok`/`resolvedForm` 仍被上面那行 `if` 使用 ⇒ 编译通过，非"编译失败当红"） | 一处 | 整包 `./internal/tools/` 跑：**FAIL 9／PASS 248／SKIP 0**（`mutation-M1.log`）。红的恰是依赖补齐的 9 枚：`A1`、`A3`、`BothSpellings/{Wisp,docs,tools}`、母格尾判（`:252` `AC#2 RED: the short spelling ... was still refused`）、`CanonicalStillNamesOneTree`；**N1–N7 七枚全 `--- PASS`** ⇒ 这一发只杀掉同形、没有牵动拒绝面。同发里 252-p1 探针仍 9 枚 `--- PASS`（它是读数件，不持期望） | `git cat-file` 还原后 md5 `7a86da7aeba8420639491cd252a47a0c`＝提交版；`go test -run 'TestTicket252R1'` → `ok`；`git status --porcelain -- internal/tools internal/risk` 为空 |
| **M2** | **`!ok` 那一支被改软**（正面自证"我没把它改软、且仪器会抓到改软"）：`paths.go:209` 的 `if !ok \|\| !rootsContain(...)` 改成 `if ok && !rootsContain(...)` | 一处 | 整包跑：**FAIL 3／PASS 254／SKIP 0**（`mutation-M2.log`）。红在：`.../N7_inside_the_root_but_unresolvable`（`:234` `AC#2 RED (widening): an unresolvable path inside the root was authorized (leg1=true leg2=false ok=false rf="")`）＋母格尾判＋**既有件 `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`**（`paths_ticket107b_probes_test.go:180` `AC#3 RED: ... leaves the allowed root ... through the link ...`） ⇒ 这台机器上"解析不到＝未授权"那一支确实有牙，本腿保留它不是形式上的 | 还原后同上：md5 一致、两包 `ok`、写面 porcelain 空 |

两发之间的量面读数：M1 还原后与 M2 跑完后各量一次 `git status --porcelain -- internal/tools internal/risk` ＝ **空**
（同一时刻 `cmd/wisp` 有他人 4 枚脏件、`internal/config` 有他人 1 枚新增，见 §0 末段——都不是本腿写的，本腿一枚没跑它们的测试）。

---

## §5 我可能写错的条目（自我对抗）

- **甲-1｜落点位移风险**：票面裁的是「进 `InAllowlist` **之前**那一步」，本腿把它落在 `Canonicalize` 尾部。
  如果编排者那句其实指"`InAllowlist` 函数体内的第一段之前"，那本腿的形与裁法**不同**（§2 末段给了为什么另一种落点会塌成一次包含）。
  ⚠ 这一条请验收腿按原话判定，不要替本腿圆场。
- **甲-2｜"生产都先 `Canonicalize`"这句本腿只数了三枚调用方**：`rules_gateway.go:45`／`bridge.go:1145`／`task.go:838`。
  这三枚是 `InAllowlist` 在非测试码里的全部调用点（`grep -rn InAllowlist cmd internal tools`），
  但**"每一枚调用点上游都真的过了 `Canonicalize`"本腿只逐字读了 `rules_gateway` 与 `task.go` 的注释，`bridge.go` 那一段的 `p` 来源没有逐字追到底**。
  若 `bridge.go:1145` 的 `p` 有可能是未经 `Canonicalize` 的原串，甲-1 的残余不对称在 `bridge` 侧就是活的。
- **甲-3｜POSIX 逐字节不变这句是**推断**（编译期常量分支 + 未跑 Linux 测试），不是读数。
  派单禁全仓 `go test ./...`，本腿只跑了 Windows 两包。
- **甲-4｜`res.Resolved == false` 的覆盖面**：本腿假定 Windows 上它等价于"链上有不存在的分量"。
  但 ACL 挡住句柄查询、离线卷这类"存在却没解析出来"也会落进这一发——本腿用 `os.Lstat` 的 `ErrNotExist` 再挡了一层，
  **`Lstat` 报 `EACCES`／`ERROR_ACCESS_DENIED` 的那一发本腿没有载具可量**（造不出"存在但 Lstat 被拒"的真形），只能断言它落回原串（守卫是 `!errors.Is(err, os.ErrNotExist)`）。
- **甲-5｜存在性竞态**：补齐读的是"问的这一刻链上存不存在"。`Canonicalize` 与真正 `open` 之间文件被人建出来／删掉，形状可能不同。
  本腿没有引入新窗口（今天的 C26 同样有这一步），但也没有读数证明它不变差；票 252 射程外，登记不结。
- **甲-6｜`GetFinalPathNameByHandle` 的大小写副作用**：补齐后的形取 OS 的磁盘拼法。
  本腿的正控只证到"两形折成同一枚、判定同答"，**没有断言大小写与配置文件里写的那串逐字相同**（`Roots()` 那枚守卫断的是 `foldPath(long)`，折叠过的大小写）。
  CI 上 `C:\Users\RUNNER~1` 那族如果磁盘大小写与 `os.TempDir()` 给的不一致，`Canonicalize` 的输出会变（这本来也是票 107/72 那族"根与锚必须同形"的处理方向）。
- **甲-7｜`errors.Is(err, os.ErrNotExist)` 对 Windows 错误码的覆盖**：本腿依赖 Go 把 `ERROR_PATH_NOT_FOUND`/`ERROR_FILE_NOT_FOUND` 归到 `ErrNotExist`。
  载具侧 N1/N2/N5（父目录在根外、兄弟树）实测走的就是这条路（它们修前修后同答），但
  **`ERROR_INVALID_NAME` 那一类本腿只在 N7 上量到"resolvedForm 拒"，没量到"Lstat 归类"**。
- **甲-8｜我把 `N6`（不存在的卷）当成 `!ok` 的证据之一**：后来逐字看 `measure252` 的打印才确认 `N6` 是先响在第一段（词法），
  真正 isolate `!ok` 的是 `N7`。§4 表里本腿已按 N7 说话；若有人把 N6 当 `!ok` 的读数，那是本件行文容易带偏的一处，这里先自纠。
- **甲-9｜票面 §23 那枚 301 秒 witness 的归口**：本腿**没有**任何读数字据把它结清（它在 `cmd/wisp`，写面被占）。
  ⛔ 不要把本件当成那一格的答案。

---

## §6 判不动的地方

**甲-1｜AC#4／AC#5／AC#6 未做（具名）**：三格的作用面都在 `cmd/wisp`
（AC#4＝`TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking`／`cmd/wisp/run_test.go:378`／`cmd/wisp/run_mode101_test.go:506`；
AC#5＝`cmd/wisp/ticket224_assembly_test.go` 的 `grantLinesOf`；AC#6＝要一发真 `wisp run` 的 `in_allowlist_scope=` 审计行）。
派单硬禁：`cmd/wisp` 此刻由别的腿在写，本腿连该包的测试一枚都没跑。
⇒ **AC#4/5/6 未做，因 `cmd/wisp` 写面被占**；本腿也没做那三格的任何"顺手"版本。

**甲-2｜AC#2 的效果在 owner 真机上撞不撞得到**：本腿的载具是**同一形状的本机复刻**（真短名＋真缺失叶子），
"他今天到底撞不撞"要真机 `[fs] allowed_dirs` 配置＋真 `wisp run` 审计行 ⇒ **量不到**（同 AC#6，属 `cmd/wisp` 面）。

**甲-3｜CI runner 那一族（`C:\Users\RUNNER~1\...`）在本机的复现度**：本机 `%TEMP%` 侧无可用别名
（`252-p1` 已 census：`TEMP/TMP/USERPROFILE/APPDATA/LOCALAPPDATA/HOME/os.TempDir()/os.UserHomeDir()` 全 `same-as-long`），
本腿因此用**仓库根**（有真别名）当临时根复刻形状。CI 上"以 temp 为根"那一发的**逐字读数：量不到**。

**甲-4｜Linux/POSIX 两包全绿**：禁全仓跑＋本机无 POSIX 环境 ⇒ **量不到**（见 §5 甲-3）。

**甲-5｜"补齐之后 R2→L2 那一跳真的不再发生"**：判定链在 `internal/risk`（本腿两包内），
但**审批卡真的少弹一枚**要看 `cmd/wisp`／面板侧的 witness ⇒ 量不到。本腿只证到 `InAllowlist` 布尔值这一层。

**甲-6｜`d22scan` 的 frontend/dist git-ignore 那一行**：不是本腿产生的，本腿也无从让它不出现 ⇒ 逐字摊在 §3，不解释成通过。

**甲-7｜票 252 §25 那句"其余 21 枚共红里有多少同因"**：本腿未做归因（一枚 `cmd/wisp` 测试没跑），
⛔ 也不许把本件的绿读成那 21 枚有了答案。

---

## §7 交件判语（本腿只交读数，判语归非实现者）

1. **做了**：AC#2（修法在 `internal/tools/paths.go:119-186`，`InAllowlist` 本体零改动）＋ AC#3（载具 `internal/tools/paths_shortname_252_r1_test.go`，四枚 `TestTicket252R1*`，本机真短名复刻，0 枚 `t.Skip`，未挂〔仅 CI 可量〕）。
2. **没做（具名）**：AC#4／AC#5／AC#6——`cmd/wisp` 写面被别的腿占，本腿未跑该包一枚测试、未写该包一字节。
3. **AC 勾选框一枚未碰**（票面 §29-34 六枚框与 §4 的复现块原样，翻勾归编排者凭非实现者验收表）。
4. **未 push**：本腿只 commit；`git remote` 未动过。commit 全部带**显式 pathspec、写在命令上**
   （`git add -- 明列路径 && git commit -F 消息文件 -- 明列路径` 同一条命令，中间不停顿），
   ⛔ 未用 `add -A`／`.`／`-a`、`commit -a`、`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／merge／worktree。
5. **未删任何文件**：突变还原全部走 `git cat-file blob HEAD:路径 > 同一路径`；备份与日志只建不删，落在 `.scratch/wisp/probes/252/r1/`（含 `backup/paths.go.orig`）。
6. **一字未动**：`docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚点名既有判据件／`.github/workflows/ci.yml`；
   `frontend/**` 与 `design/**` 未读、未引、未转述（`d22scan` 输出里那两枚 scope 行是仪器自陈，不是本腿读过的内容）。
7. **写面终态**：`git status --porcelain -- internal/tools internal/risk` ＝ 空；改动只有 `internal/tools/paths.go`（1 枚 hunk）与新增测试 1 件；
   证据与日志落在 `.scratch/wisp/probes/252/r1/`。⚠ `cmd/wisp`／`internal/config` 的脏件是别人的活，本腿没碰。
8. **门禁四数（终态）**：`go build ./...` 退出 0；两包 `ok`／`--- FAIL` 0；`gofumpt -l` 空；`d22scan` `clean`（退出 0）；新判据 18 `--- PASS`／0 FAIL／0 SKIP；`252-p1` 探针 9 PASS 复跑不退。
9. **推翻编排者派单之处**（本腿自己撞到的、与票面/裁定文字不同者）：
   (a) 票面 §30 与 §83 说"asked 与 roots 走同一条已存在的折叠"，本腿实测这条折叠**不能同时喂给第一段与第二段**（喂了就会塌成一次包含，且新增一次"词法在根外但解析落回根内"的授权）⇒ 本腿把折叠只加在 `Canonicalize` 这一步，并把不对称写在 §2 末段与 §5 甲-1，请裁定这是否算遵裁。
   (b) 其余行号位移：票面/裁定引的 `paths.go:133/140/152/153` 在本腿提交后为 `:189/:196/:208/:209`（内容逐字节相同，见 §2 的 md5）。
